package server_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"attestor/optimism/opstack"
	"attestor/optimism/server"
	attestorpb "attestor/types/attestor"
)

// stubReplica satisfies opstack's unexported replica interface structurally, so
// an external test package can still inject one.
type stubReplica struct {
	status      opstack.SyncStatus
	commitments map[uint64]opstack.L2Commitment
	err         error
}

func (s stubReplica) SyncStatus(context.Context) (opstack.SyncStatus, error) {
	return s.status, s.err
}

func (s stubReplica) OutputAtBlock(_ context.Context, block uint64) ([32]byte, error) {
	c, ok := s.commitments[block]
	if !ok {
		return [32]byte{}, errNotFound
	}
	return c.OutputRoot, nil
}

func (s stubReplica) CommitmentAt(_ context.Context, block uint64) (opstack.L2Commitment, error) {
	c, ok := s.commitments[block]
	if !ok {
		return opstack.L2Commitment{}, errNotFound
	}
	return c, nil
}

var errNotFound = errNotFoundType{}

type errNotFoundType struct{}

func (errNotFoundType) Error() string { return "block not derived" }

func newVerifyAttestor(t *testing.T, srcChain string, rep stubReplica) *opstack.OpStackAttestor {
	t.Helper()
	store, err := opstack.LoadStore(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	cfg := opstack.Config{SrcChain: srcChain, AttestationHead: opstack.HeadFinalized}
	return opstack.New(cfg, nil, rep, store, nil, nil, nil)
}

func verifyFixture(t *testing.T) (*server.Server, opstack.L2Commitment) {
	t.Helper()
	commitment := opstack.L2Commitment{
		BlockNumber: 500,
		BlockHash:   common.HexToHash("0xaa11"),
		StateRoot:   common.HexToHash("0xbb22"),
	}
	rep := stubReplica{
		status:      opstack.SyncStatus{UnsafeL2: 900, SafeL2: 700, FinalizedL2: 600},
		commitments: map[uint64]opstack.L2Commitment{500: commitment},
	}
	a := newVerifyAttestor(t, "op-test", rep)
	return server.New(map[string]*opstack.OpStackAttestor{"op-test": a}), commitment
}

// The reason the RPC exists: the relayer asks whether the block it is about to
// package is the one this attestor's replica holds at that height.
func TestVerifyStateRootAcceptsTheReplicaBlock(t *testing.T) {
	srv, commitment := verifyFixture(t)
	c := attestorpb.NewAttestorServiceClient(dial(t, srv))

	resp, err := c.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		BlockNumber:       500,
		ExpectedStateRoot: commitment.StateRoot.Bytes(),
		ExpectedBlockHash: commitment.BlockHash.Bytes(),
		RunMode:           attestorpb.RunMode_RUN_MODE_FINALIZED,
	})
	if err != nil {
		t.Fatalf("VerifyStateRoot: %v", err)
	}
	if !resp.GetValid() {
		t.Fatal("valid = false for the replica's own block")
	}
	if resp.GetBlockNumber() != 500 {
		t.Errorf("block number = %d, want 500", resp.GetBlockNumber())
	}
	if common.BytesToHash(resp.GetStateRoot()) != commitment.StateRoot {
		t.Errorf("state root echo = %x", resp.GetStateRoot())
	}
}

// A different block at an approved height is exactly the reorg/wrong-endpoint
// case the binding exists to catch.
func TestVerifyStateRootRejectsADifferentBlock(t *testing.T) {
	srv, commitment := verifyFixture(t)
	c := attestorpb.NewAttestorServiceClient(dial(t, srv))

	for _, tc := range []struct {
		name  string
		state []byte
		hash  []byte
	}{
		{"wrong state root", common.HexToHash("0xdead").Bytes(), commitment.BlockHash.Bytes()},
		{"wrong block hash", commitment.StateRoot.Bytes(), common.HexToHash("0xdead").Bytes()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := c.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
				BlockNumber:       500,
				ExpectedStateRoot: tc.state,
				ExpectedBlockHash: tc.hash,
				RunMode:           attestorpb.RunMode_RUN_MODE_FINALIZED,
			})
			if err != nil {
				t.Fatalf("VerifyStateRoot: %v", err)
			}
			if resp.GetValid() {
				t.Fatal("valid = true for a block the replica does not hold")
			}
		})
	}
}

// The block hash is optional; omitting it must still verify the state root
// rather than accepting anything.
func TestVerifyStateRootAcceptsAnOmittedBlockHash(t *testing.T) {
	srv, commitment := verifyFixture(t)
	c := attestorpb.NewAttestorServiceClient(dial(t, srv))

	resp, err := c.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		BlockNumber:       500,
		ExpectedStateRoot: commitment.StateRoot.Bytes(),
		RunMode:           attestorpb.RunMode_RUN_MODE_FINALIZED,
	})
	if err != nil {
		t.Fatalf("VerifyStateRoot: %v", err)
	}
	if !resp.GetValid() {
		t.Fatal("valid = false when only the state root was supplied")
	}
}

// "Above the head" must not answer valid=false: the caller retries that, but
// treats a false as a divergence and stops.
func TestVerifyStateRootSeparatesNotYetDerivedFromWrong(t *testing.T) {
	srv, commitment := verifyFixture(t)
	c := attestorpb.NewAttestorServiceClient(dial(t, srv))

	_, err := c.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		BlockNumber:       800, // above the finalized head of 600
		ExpectedStateRoot: commitment.StateRoot.Bytes(),
		RunMode:           attestorpb.RunMode_RUN_MODE_FINALIZED,
	})
	if got := status.Code(err); got != codes.FailedPrecondition {
		t.Fatalf("code = %v, want FailedPrecondition (got err %v)", got, err)
	}
}

// The run mode selects which head bounds the answer, so the same block can be
// above the finalized head and below the unsafe one.
func TestVerifyStateRootHonoursRunMode(t *testing.T) {
	commitment := opstack.L2Commitment{
		BlockNumber: 800,
		BlockHash:   common.HexToHash("0xcc33"),
		StateRoot:   common.HexToHash("0xdd44"),
	}
	rep := stubReplica{
		status:      opstack.SyncStatus{UnsafeL2: 900, SafeL2: 700, FinalizedL2: 600},
		commitments: map[uint64]opstack.L2Commitment{800: commitment},
	}
	a := newVerifyAttestor(t, "op-test", rep)
	c := attestorpb.NewAttestorServiceClient(dial(t, server.New(map[string]*opstack.OpStackAttestor{"op-test": a})))

	resp, err := c.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		BlockNumber:       800,
		ExpectedStateRoot: commitment.StateRoot.Bytes(),
		RunMode:           attestorpb.RunMode_RUN_MODE_UNSAFE,
	})
	if err != nil {
		t.Fatalf("VerifyStateRoot at unsafe: %v", err)
	}
	if !resp.GetValid() {
		t.Fatal("valid = false for a block below the unsafe head")
	}
}

func TestVerifyStateRootRejectsMalformedRequests(t *testing.T) {
	srv, commitment := verifyFixture(t)
	c := attestorpb.NewAttestorServiceClient(dial(t, srv))

	for _, tc := range []struct {
		name string
		req  *attestorpb.VerifyStateRootRequest
	}{
		{"short state root", &attestorpb.VerifyStateRootRequest{
			BlockNumber: 500, ExpectedStateRoot: []byte{1, 2, 3},
			RunMode: attestorpb.RunMode_RUN_MODE_FINALIZED,
		}},
		{"short block hash", &attestorpb.VerifyStateRootRequest{
			BlockNumber: 500, ExpectedStateRoot: commitment.StateRoot.Bytes(),
			ExpectedBlockHash: []byte{9}, RunMode: attestorpb.RunMode_RUN_MODE_FINALIZED,
		}},
		{"unspecified run mode", &attestorpb.VerifyStateRootRequest{
			BlockNumber: 500, ExpectedStateRoot: commitment.StateRoot.Bytes(),
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := c.VerifyStateRoot(context.Background(), tc.req); status.Code(err) != codes.InvalidArgument {
				t.Fatalf("code = %v, want InvalidArgument", status.Code(err))
			}
		})
	}
}

// VerifyStateRootRequest carries no src_chain, so a multi-chain daemon cannot
// tell which replica is meant. Guessing would answer valid=false for a good
// header, which reads as a divergence — refuse instead.
func TestVerifyStateRootRefusesWhenSeveralChainsAreServed(t *testing.T) {
	rep := stubReplica{status: opstack.SyncStatus{FinalizedL2: 600}}
	srv := server.New(map[string]*opstack.OpStackAttestor{
		"op-test":   newVerifyAttestor(t, "op-test", rep),
		"base-test": newVerifyAttestor(t, "base-test", rep),
	})
	c := attestorpb.NewAttestorServiceClient(dial(t, srv))

	_, err := c.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		BlockNumber:       500,
		ExpectedStateRoot: common.HexToHash("0xbb22").Bytes(),
		RunMode:           attestorpb.RunMode_RUN_MODE_FINALIZED,
	})
	if got := status.Code(err); got != codes.FailedPrecondition {
		t.Fatalf("code = %v, want FailedPrecondition", got)
	}
}
