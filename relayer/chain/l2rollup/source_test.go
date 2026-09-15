package l2rollup

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"attestor/types/attestation"
	attestorpb "attestor/types/attestor"
	"relayer/chain"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// fakeAttestor is a hand-rolled AttestorClient recording the includeProvisional
// flag the source passes, so the head-kind → provisional mapping is testable.
type fakeAttestor struct {
	root  *attestorpb.AttestedRoot
	found bool
	err   error

	verifyBlock   <-chan struct{}
	frontierBlock <-chan struct{}

	gotSrcChain      string
	gotProvisional   bool
	gotAtOrBelowArg  uint64
	atOrBelowCalled  bool
	attestedUpCalled bool

	// verifyValid is what VerifyStateRoot answers; verifyErr overrides it.
	verifyValid       bool
	verifyErr         error
	verifySigner      attestation.Signer
	gotVerifyHeight   uint64
	gotVerifyRoot     []byte
	gotVerifyHash     []byte
	gotVerifyMode     attestorpb.RunMode
	gotVerifySrcChain string
	gotVerifyRouter   [20]byte
	gotVerifySetHash  [32]byte
	verifyCalled      bool
}

func (f *fakeAttestor) VerifyStateRoot(ctx context.Context, request VerificationRequest) (SignedVerdict, error) {
	if f.verifyBlock != nil {
		select {
		case <-ctx.Done():
			return SignedVerdict{}, ctx.Err()
		case <-f.verifyBlock:
		}
	}
	f.verifyCalled = true
	f.gotVerifySrcChain = request.SrcChain
	f.gotVerifyHeight = request.BlockNumber
	f.gotVerifyRoot = append([]byte(nil), request.StateRoot...)
	f.gotVerifyHash = append([]byte(nil), request.BlockHash...)
	f.gotVerifyMode = request.RunMode
	f.gotVerifyRouter = request.L2Router
	f.gotVerifySetHash = request.AttestorSetHash
	if f.verifyErr != nil {
		return SignedVerdict{}, f.verifyErr
	}
	verdict := SignedVerdict{
		Valid:       f.verifyValid,
		BlockNumber: request.BlockNumber,
		BlockHash:   append([]byte(nil), request.BlockHash...),
		StateRoot:   append([]byte(nil), request.StateRoot...),
	}
	if f.verifyValid && f.verifySigner.Configured() {
		signature, err := f.verifySigner.Sign(request.L2Router[:], request.AttestorSetHash[:], request.BlockNumber, request.BlockHash, request.StateRoot)
		if err != nil {
			return SignedVerdict{}, err
		}
		verdict.Signature = signature
	}
	return verdict, nil
}

func (f *fakeAttestor) AttestedUpTo(ctx context.Context, srcChain string, includeProvisional bool) (*attestorpb.AttestedRoot, bool, error) {
	if f.frontierBlock != nil {
		select {
		case <-ctx.Done():
			return nil, false, ctx.Err()
		case <-f.frontierBlock:
		}
	}
	f.attestedUpCalled = true
	f.gotSrcChain = srcChain
	f.gotProvisional = includeProvisional
	return f.root, f.found, f.err
}

func (f *fakeAttestor) AttestedRootAtOrBelow(_ context.Context, srcChain string, l2BlockNumber uint64, includeProvisional bool) (*attestorpb.AttestedRoot, bool, error) {
	f.atOrBelowCalled = true
	f.gotSrcChain = srcChain
	f.gotAtOrBelowArg = l2BlockNumber
	f.gotProvisional = includeProvisional
	return f.root, f.found, f.err
}

func newAttestedSource(kind HeadKind, at *fakeAttestor) *Source {
	return NewSource(chain.ChainType("opstack"), nil, kind, "07-tendermint-0", "08-wasm-0", "op-sepolia", ethcommon.Address{}, at, kind != Finalized)
}

func TestRelayableHeight_UsesAttestedFrontier(t *testing.T) {
	at := &fakeAttestor{root: &attestorpb.AttestedRoot{L2BlockNumber: 100}, found: true}
	got, err := newAttestedSource(Safe, at).RelayableHeight(context.Background())
	if err != nil {
		t.Fatalf("relayable height: %v", err)
	}
	if got != 100 {
		t.Fatalf("relayable height = %d, want 100", got)
	}
	if !at.attestedUpCalled || at.gotSrcChain != "op-sepolia" {
		t.Fatalf("attestor not queried with src chain: %+v", at)
	}
}

func TestRelayableHeight_UsesQuorumWhenFirstAttestorHangs(t *testing.T) {
	blocked := make(chan struct{})
	attestors := []SigningAttestor{
		{Index: 0, Client: &fakeAttestor{frontierBlock: blocked}},
		{Index: 1, Client: &fakeAttestor{root: &attestorpb.AttestedRoot{L2BlockNumber: 100}, found: true}},
		{Index: 2, Client: &fakeAttestor{root: &attestorpb.AttestedRoot{L2BlockNumber: 90}, found: true}},
	}
	frontier, err := NewQuorumAttestationFrontier(attestors, 2)
	if err != nil {
		t.Fatal(err)
	}
	source := NewSource(chain.ChainType("opstack"), nil, Safe, "07-tendermint-0", "08-wasm-0", "op-sepolia", ethcommon.Address{}, frontier, true)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	got, err := source.RelayableHeight(ctx)
	if err != nil {
		t.Fatalf("healthy frontier quorum blocked by first endpoint: %v", err)
	}
	if got != 90 {
		t.Fatalf("relayable height = %d, want quorum minimum 90", got)
	}
}

func TestQuorumFrontier_WaitsWhenReachableMembersHaveNotAttested(t *testing.T) {
	frontier, err := NewQuorumAttestationFrontier([]SigningAttestor{
		{Index: 0, Client: &fakeAttestor{found: false}},
		{Index: 1, Client: &fakeAttestor{found: false}},
		{Index: 2, Client: &fakeAttestor{err: fmt.Errorf("%w: unavailable", ErrAttestorUnavailable)}},
	}, 2)
	if err != nil {
		t.Fatal(err)
	}
	root, found, err := frontier.AttestedUpTo(context.Background(), "op", true)
	if err != nil || found || root != nil {
		t.Fatalf("frontier = (%v, %v, %v), want normal waiting state", root, found, err)
	}
}

func TestQuorumFrontier_ClassifiesImpossiblePermanentQuorum(t *testing.T) {
	frontier, err := NewQuorumAttestationFrontier([]SigningAttestor{
		{Index: 0, Client: &fakeAttestor{err: fmt.Errorf("%w: bad route", ErrAttestorUnknownRoute)}},
		{Index: 1, Client: &fakeAttestor{found: false}},
	}, 2)
	if err != nil {
		t.Fatal(err)
	}
	_, _, queryErr := frontier.AttestedUpTo(context.Background(), "op", true)
	if queryErr == nil || !errors.Is(queryErr, ErrAttestorUnknownRoute) || !chain.IsPermanent(queryErr) {
		t.Fatalf("error = %v, want permanent impossible quorum", queryErr)
	}
}

// TestRelayableHeight_NotAttestedYet: before the first attestation the source must
// report 0 (nothing relayable) — never a raw head — so the module waits.
func TestRelayableHeight_NotAttestedYet(t *testing.T) {
	at := &fakeAttestor{found: false}
	got, err := newAttestedSource(Finalized, at).RelayableHeight(context.Background())
	if err != nil {
		t.Fatalf("relayable height: %v", err)
	}
	if got != 0 {
		t.Fatalf("relayable height = %d, want 0 (not yet attested)", got)
	}
}

func TestRelayableHeight_ErrorPropagates(t *testing.T) {
	at := &fakeAttestor{err: errors.New("attestor down")}
	_, err := newAttestedSource(Safe, at).RelayableHeight(context.Background())
	if err == nil {
		t.Fatal("expected attestor error to propagate")
	}
}

// TestRelayableHeight_ProvisionalMapping locks the head-kind → include_provisional
// contract: only Finalized demands non-provisional (finalized) roots.
func TestRelayableHeight_ProvisionalMapping(t *testing.T) {
	cases := map[HeadKind]bool{
		Unsafe:    true,
		Safe:      true,
		Finalized: false,
	}
	for kind, wantProvisional := range cases {
		at := &fakeAttestor{root: &attestorpb.AttestedRoot{L2BlockNumber: 1}, found: true}
		if _, err := newAttestedSource(kind, at).RelayableHeight(context.Background()); err != nil {
			t.Fatalf("kind %d: %v", kind, err)
		}
		if at.gotProvisional != wantProvisional {
			t.Fatalf("kind %d: include_provisional = %v, want %v", kind, at.gotProvisional, wantProvisional)
		}
	}
}
