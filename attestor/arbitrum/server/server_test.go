package server

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"attestor/arbitrum"
	attestorpb "attestor/types/attestor"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestVerifyStateRoot(t *testing.T) {
	header := &types.Header{
		Number: big.NewInt(123),
		Root:   common.HexToHash("0x1234"),
	}
	verifier, err := NewAttestorServer(headerReaderFunc(func(
		_ context.Context,
		number *big.Int,
	) (*types.Header, error) {
		if number == nil {
			return &types.Header{Number: big.NewInt(200)}, nil
		}
		if number.Uint64() != 123 {
			t.Fatalf("requested block: got %d want 123", number.Uint64())
		}
		return types.CopyHeader(header), nil
	}))
	if err != nil {
		t.Fatalf("create verifier: %v", err)
	}

	response, err := verifier.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		BlockNumber:       123,
		ExpectedStateRoot: header.Root.Bytes(),
		ExpectedBlockHash: header.Hash().Bytes(),
		RunMode:           attestorpb.RunMode_RUN_MODE_UNSAFE,
	})
	if err != nil {
		t.Fatalf("verify state root: %v", err)
	}
	if !response.GetValid() {
		t.Fatal("matching Nitro commitment was rejected")
	}
	if common.BytesToHash(response.GetStateRoot()) != header.Root ||
		common.BytesToHash(response.GetBlockHash()) != header.Hash() {
		t.Fatalf("canonical response mismatch: %+v", response)
	}

	response, err = verifier.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		BlockNumber:       123,
		ExpectedStateRoot: common.HexToHash("0x5678").Bytes(),
		RunMode:           attestorpb.RunMode_RUN_MODE_UNSAFE,
	})
	if err != nil {
		t.Fatalf("verify mismatched state root: %v", err)
	}
	if response.GetValid() {
		t.Fatal("mismatched state root was accepted")
	}
}

func TestAttestedFrontierRPCs(t *testing.T) {
	runtime, err := arbitrum.NewRuntimeState(headerReaderFunc(func(
		context.Context,
		*big.Int,
	) (*types.Header, error) {
		return &types.Header{Number: big.NewInt(1)}, nil
	}))
	if err != nil {
		t.Fatalf("create runtime: %v", err)
	}
	store, err := arbitrum.NewAttestedRootStore("arbitrum-one", 1)
	if err != nil {
		t.Fatalf("create feed: %v", err)
	}
	for _, item := range []struct {
		height      uint64
		id          byte
		provisional bool
	}{
		{height: 90, id: 1, provisional: false},
		{height: 100, id: 2, provisional: true},
	} {
		blockHash := common.BigToHash(new(big.Int).SetUint64(item.height))
		proposal := arbitrum.ProposedAssertion{
			AssertionHash: common.Hash{31: item.id},
			L2BlockHash:   blockHash,
		}
		if err := store.RecordProposal(proposal); err != nil {
			t.Fatalf("record proposal: %v", err)
		}
		if err := store.RecordAssertionAttestation(
			proposal,
			arbitrum.BlockCommitment{
				BlockNumber: item.height,
				BlockHash:   blockHash,
				StateRoot:   common.Hash{30: item.id},
			},
			time.Unix(1_700_000_000+int64(item.height), 0),
			item.provisional,
		); err != nil {
			t.Fatalf("record attestation: %v", err)
		}
	}
	service, err := NewAttestorServerWithRuntimeAndFeeds(
		runtime,
		map[string]AttestedRootReader{"arbitrum-one": store},
	)
	if err != nil {
		t.Fatalf("create service: %v", err)
	}

	info, err := service.Info(context.Background(), &attestorpb.InfoRequest{})
	if err != nil {
		t.Fatalf("query info: %v", err)
	}
	if len(info.GetChains()) != 1 {
		t.Fatalf("info response: %+v", info)
	}
	chain := info.GetChains()[0]
	if chain.GetSrcChain() != "arbitrum-one" ||
		chain.GetAttestedUpTo().GetL2BlockNumber() != 90 ||
		chain.GetAttestedUpToProvisional().GetL2BlockNumber() != 100 {
		t.Fatalf("info chain: %+v", chain)
	}
	if chain.GetAttestationHead() != "" {
		t.Fatalf("attestation_head = %q, want unset", chain.GetAttestationHead())
	}

	confirmed, err := service.AttestedUpTo(
		context.Background(),
		&attestorpb.AttestedUpToRequest{SrcChain: "arbitrum-one"},
	)
	if err != nil {
		t.Fatalf("query confirmed frontier: %v", err)
	}
	if !confirmed.GetFound() || confirmed.GetRoot().GetL2BlockNumber() != 90 {
		t.Fatalf("confirmed frontier: %+v", confirmed)
	}

	provisional, err := service.AttestedUpTo(
		context.Background(),
		&attestorpb.AttestedUpToRequest{
			SrcChain:           "arbitrum-one",
			IncludeProvisional: true,
		},
	)
	if err != nil {
		t.Fatalf("query provisional frontier: %v", err)
	}
	if !provisional.GetFound() ||
		provisional.GetRoot().GetL2BlockNumber() != 100 ||
		!provisional.GetRoot().GetProvisional() ||
		common.BytesToHash(provisional.GetRoot().GetAssertionHash()) != (common.Hash{31: 2}) {
		t.Fatalf("provisional frontier: %+v", provisional)
	}

	atOrBelow, err := service.AttestedRootAtOrBelow(
		context.Background(),
		&attestorpb.AttestedRootAtOrBelowRequest{
			SrcChain:           "arbitrum-one",
			L2BlockNumber:      95,
			IncludeProvisional: true,
		},
	)
	if err != nil {
		t.Fatalf("query frontier at or below: %v", err)
	}
	if !atOrBelow.GetFound() || atOrBelow.GetRoot().GetL2BlockNumber() != 90 {
		t.Fatalf("frontier at or below: %+v", atOrBelow)
	}
}

func TestAttestedFrontierRejectsUnknownSourceChain(t *testing.T) {
	runtime, err := arbitrum.NewRuntimeState(headerReaderFunc(func(
		context.Context,
		*big.Int,
	) (*types.Header, error) {
		return &types.Header{Number: big.NewInt(1)}, nil
	}))
	if err != nil {
		t.Fatalf("create runtime: %v", err)
	}
	service, err := NewAttestorServerWithRuntime(runtime)
	if err != nil {
		t.Fatalf("create service: %v", err)
	}
	_, err = service.AttestedUpTo(
		context.Background(),
		&attestorpb.AttestedUpToRequest{SrcChain: "unknown"},
	)
	if status.Code(err) != codes.NotFound {
		t.Fatalf("unknown source status: got %v want %v", status.Code(err), codes.NotFound)
	}
}

func TestAttestedRootProtoUsesExclusiveProvenance(t *testing.T) {
	game := attestedRootToProto(arbitrum.AttestedRoot{
		Source:    arbitrum.SourceGame,
		GameIndex: 0,
	})
	if _, ok := game.GetProvenance().(*attestorpb.AttestedRoot_GameIndex); !ok {
		t.Fatalf("game index zero has no explicit oneof presence: %T", game.GetProvenance())
	}

	assertionHash := common.HexToHash("0x1234")
	assertion := attestedRootToProto(arbitrum.AttestedRoot{
		Source:        arbitrum.SourceAssertion,
		GameIndex:     42,
		AssertionHash: assertionHash,
	})
	if _, ok := assertion.GetProvenance().(*attestorpb.AttestedRoot_AssertionHash); !ok {
		t.Fatalf("assertion has wrong provenance type: %T", assertion.GetProvenance())
	}
	if common.BytesToHash(assertion.GetAssertionHash()) != assertionHash {
		t.Fatalf("assertion provenance: got %x want %s", assertion.GetAssertionHash(), assertionHash)
	}
	if assertion.GetGameIndex() != 0 {
		t.Fatalf("assertion unexpectedly exposed game index %d", assertion.GetGameIndex())
	}

}

func TestVerifyStateRootRejectsMalformedCommitments(t *testing.T) {
	verifier, err := NewAttestorServer(headerReaderFunc(func(
		context.Context,
		*big.Int,
	) (*types.Header, error) {
		t.Fatal("malformed request reached Nitro")
		return nil, nil
	}))
	if err != nil {
		t.Fatalf("create verifier: %v", err)
	}

	_, err = verifier.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		ExpectedStateRoot: []byte{1},
		RunMode:           attestorpb.RunMode_RUN_MODE_UNSAFE,
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("malformed root status: got %v want %v", status.Code(err), codes.InvalidArgument)
	}

	_, err = verifier.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		ExpectedStateRoot: make([]byte, common.HashLength),
		ExpectedBlockHash: []byte{1},
		RunMode:           attestorpb.RunMode_RUN_MODE_UNSAFE,
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("malformed block hash status: got %v want %v", status.Code(err), codes.InvalidArgument)
	}
}

func TestVerifyStateRootReportsMissingNitroBlock(t *testing.T) {
	verifier, err := NewAttestorServer(headerReaderFunc(func(
		_ context.Context,
		number *big.Int,
	) (*types.Header, error) {
		if number == nil {
			return &types.Header{Number: big.NewInt(1_000)}, nil
		}
		return nil, ethereum.NotFound
	}))
	if err != nil {
		t.Fatalf("create verifier: %v", err)
	}

	_, err = verifier.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		BlockNumber:       999,
		ExpectedStateRoot: make([]byte, common.HashLength),
		RunMode:           attestorpb.RunMode_RUN_MODE_UNSAFE,
	})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("missing block status: got %v want %v", status.Code(err), codes.NotFound)
	}
}

func TestNewAttestorServerRejectsMissingReader(t *testing.T) {
	if _, err := NewAttestorServer(nil); err == nil {
		t.Fatal("expected missing Nitro reader error")
	}
}

func TestVerifyStateRootEnforcesSafeHead(t *testing.T) {
	verifier, err := NewAttestorServer(headerReaderFunc(func(
		_ context.Context,
		number *big.Int,
	) (*types.Header, error) {
		if number.Int64() >= 0 {
			t.Fatal("block above safe head reached canonical block lookup")
		}
		return &types.Header{Number: big.NewInt(100)}, nil
	}))
	if err != nil {
		t.Fatalf("create safe verifier: %v", err)
	}

	_, err = verifier.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		BlockNumber:       101,
		ExpectedStateRoot: make([]byte, common.HashLength),
		RunMode:           attestorpb.RunMode_RUN_MODE_SAFE,
	})
	if status.Code(err) != codes.FailedPrecondition {
		t.Fatalf("block above safe head status: got %v want %v", status.Code(err), codes.FailedPrecondition)
	}
}

func TestVerifyStateRootEnforcesFinalizedHead(t *testing.T) {
	header := &types.Header{
		Number: big.NewInt(100),
		Root:   common.HexToHash("0x100"),
	}
	verifier, err := NewAttestorServer(headerReaderFunc(func(
		_ context.Context,
		number *big.Int,
	) (*types.Header, error) {
		switch {
		case number.Int64() == int64(rpc.FinalizedBlockNumber):
			return types.CopyHeader(header), nil
		case number.Sign() >= 0 && number.Uint64() == 100:
			return types.CopyHeader(header), nil
		default:
			t.Fatalf("unexpected finalized-mode selector: %s", number)
			return nil, nil
		}
	}))
	if err != nil {
		t.Fatalf("create finalized verifier: %v", err)
	}

	response, err := verifier.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		BlockNumber:       100,
		ExpectedStateRoot: header.Root.Bytes(),
		ExpectedBlockHash: header.Hash().Bytes(),
		RunMode:           attestorpb.RunMode_RUN_MODE_FINALIZED,
	})
	if err != nil {
		t.Fatalf("verify finalized state root: %v", err)
	}
	if !response.GetValid() {
		t.Fatal("matching finalized commitment was rejected")
	}

	_, err = verifier.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		BlockNumber:       101,
		ExpectedStateRoot: make([]byte, common.HashLength),
		RunMode:           attestorpb.RunMode_RUN_MODE_FINALIZED,
	})
	if status.Code(err) != codes.FailedPrecondition {
		t.Fatalf("block above finalized head status: got %v want %v", status.Code(err), codes.FailedPrecondition)
	}
}

func TestVerifyStateRootRequiresRunMode(t *testing.T) {
	verifier, err := NewAttestorServer(headerReaderFunc(func(
		context.Context,
		*big.Int,
	) (*types.Header, error) {
		t.Fatal("request without a run mode reached Nitro")
		return nil, nil
	}))
	if err != nil {
		t.Fatalf("create verifier: %v", err)
	}

	_, err = verifier.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
		ExpectedStateRoot: make([]byte, common.HashLength),
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("unspecified run mode status: got %v want %v", status.Code(err), codes.InvalidArgument)
	}
}

func TestVerifyStateRootSelectsRunModePerRequest(t *testing.T) {
	header := &types.Header{
		Number: big.NewInt(100),
		Root:   common.HexToHash("0x100"),
	}
	selected := make([]arbitrum.RunMode, 0, 3)
	verifier, err := NewAttestorServer(headerReaderFunc(func(
		_ context.Context,
		number *big.Int,
	) (*types.Header, error) {
		switch {
		case number == nil:
			selected = append(selected, arbitrum.RunModeUnsafe)
		case number.Int64() == int64(rpc.SafeBlockNumber):
			selected = append(selected, arbitrum.RunModeSafe)
		case number.Int64() == int64(rpc.FinalizedBlockNumber):
			selected = append(selected, arbitrum.RunModeFinalized)
		case number.Uint64() == 100:
			return types.CopyHeader(header), nil
		default:
			t.Fatalf("unexpected block selector: %s", number)
		}
		return types.CopyHeader(header), nil
	}))
	if err != nil {
		t.Fatalf("create verifier: %v", err)
	}

	for _, mode := range []attestorpb.RunMode{
		attestorpb.RunMode_RUN_MODE_UNSAFE,
		attestorpb.RunMode_RUN_MODE_SAFE,
		attestorpb.RunMode_RUN_MODE_FINALIZED,
	} {
		response, verifyErr := verifier.VerifyStateRoot(context.Background(), &attestorpb.VerifyStateRootRequest{
			BlockNumber:       100,
			ExpectedStateRoot: header.Root.Bytes(),
			RunMode:           mode,
		})
		if verifyErr != nil {
			t.Fatalf("verify with %s: %v", mode, verifyErr)
		}
		if !response.GetValid() {
			t.Fatalf("matching commitment rejected for %s", mode)
		}
	}

	want := []arbitrum.RunMode{
		arbitrum.RunModeUnsafe,
		arbitrum.RunModeSafe,
		arbitrum.RunModeFinalized,
	}
	if len(selected) != len(want) {
		t.Fatalf("selected heads: got %v want %v", selected, want)
	}
	for index := range want {
		if selected[index] != want[index] {
			t.Fatalf("selected heads: got %v want %v", selected, want)
		}
	}
}

type headerReaderFunc func(context.Context, *big.Int) (*types.Header, error)

func (f headerReaderFunc) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	if f == nil {
		return nil, errors.New("nil header reader")
	}
	return f(ctx, number)
}
