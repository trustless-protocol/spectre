package adapter

import (
	"context"
	"crypto/ed25519"
	"errors"
	"math/big"
	"testing"

	"attestor/arbitrum"
	"attestor/core"
	"attestor/types/attestation"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type contractFeed struct{}

func (contractFeed) HighestAttested(bool) (arbitrum.AttestedRoot, bool) {
	return arbitrum.AttestedRoot{}, false
}

func (contractFeed) HighestAttestedAtOrBelow(uint64, bool) (arbitrum.AttestedRoot, bool) {
	return arbitrum.AttestedRoot{}, false
}

type contractHeaderReader struct {
	head      *types.Header
	candidate *types.Header
	err       error
}

func (r contractHeaderReader) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	if r.err != nil {
		return nil, r.err
	}
	switch {
	case number == nil:
		return types.CopyHeader(r.head), nil
	case number.Int64() == int64(rpc.SafeBlockNumber):
		return types.CopyHeader(r.head), nil
	case number.Int64() == int64(rpc.FinalizedBlockNumber):
		return types.CopyHeader(r.head), nil
	case number.Uint64() == r.candidate.Number.Uint64():
		return types.CopyHeader(r.candidate), nil
	default:
		return nil, errors.New("unexpected Nitro height")
	}
}

func newContractAdapter(t *testing.T, reader arbitrum.NitroHeaderReader, signer attestation.Signer) *Adapter {
	t.Helper()
	runtime, err := arbitrum.NewRuntimeState(reader)
	if err != nil {
		t.Fatalf("NewRuntimeState: %v", err)
	}
	bridge, err := New(Options{
		SrcChain:        "arbitrum-contract",
		Runtime:         runtime,
		Feed:            contractFeed{},
		AttestationHead: arbitrum.RunModeFinalized,
		Signer:          signer,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return bridge
}

func TestCorePortContract(t *testing.T) {
	head := &types.Header{Number: big.NewInt(10), Root: common.HexToHash("0x10")}
	candidate := &types.Header{Number: big.NewInt(7), Root: common.HexToHash("0x70")}
	reader := contractHeaderReader{head: head, candidate: candidate}
	signer, err := attestation.NewSigner(42161, "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20")
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	bridge := newContractAdapter(t, reader, signer)

	if _, found, err := bridge.AttestedUpTo(context.Background(), core.AttestationPolicy{}); err != nil || found {
		t.Fatalf("empty durable feed = (_, %t, %v), want not found", found, err)
	}
	if _, found, err := bridge.AttestedRootAtOrBelow(context.Background(), 7, core.AttestationPolicy{IncludeProvisional: true}); err != nil || found {
		t.Fatalf("empty bounded durable feed = (_, %t, %v), want not found", found, err)
	}
	initialStatus, err := bridge.Status(context.Background())
	if err != nil || initialStatus.Ready || initialStatus.SrcChain != "arbitrum-contract" {
		t.Fatalf("initial status = (%+v, %v), want unready local snapshot", initialStatus, err)
	}

	var stateRoot [32]byte
	copy(stateRoot[:], candidate.Root[:])
	blockHash := [32]byte(candidate.Hash())
	request := core.BlockIdentityRequest{
		BlockNumber:       7,
		ExpectedStateRoot: stateRoot,
		ExpectedBlockHash: &blockHash,
		RunMode:           core.RunModeFinalized,
	}
	for _, tc := range []struct {
		name      string
		request   core.RunMode
		signedFor attestation.RunMode
	}{
		{name: "unsafe", request: core.RunModeUnsafe, signedFor: attestation.RunModeUnsafe},
		{name: "safe", request: core.RunModeSafe, signedFor: attestation.RunModeSafe},
		{name: "finalized", request: core.RunModeFinalized, signedFor: attestation.RunModeFinalized},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request.RunMode = tc.request
			verdict, err := bridge.VerifyStateRoot(context.Background(), request)
			if err != nil || !verdict.Valid || len(verdict.Signature) != ed25519.SignatureSize {
				t.Fatalf("matching verdict = (%+v, %v), want signed valid", verdict, err)
			}
			message, err := attestation.SigningBytes(42161, tc.signedFor, 7, candidate.Root[:], candidate.Hash().Bytes())
			if err != nil {
				t.Fatalf("SigningBytes: %v", err)
			}
			if !ed25519.Verify(signer.PublicKey(), message, verdict.Signature) {
				t.Fatal("positive verdict signature does not bind the requested run mode")
			}
		})
	}

	request.RunMode = core.RunModeFinalized
	request.ExpectedStateRoot = [32]byte{1}
	verdict, err := bridge.VerifyStateRoot(context.Background(), request)
	if err != nil || verdict.Valid || len(verdict.Signature) != 0 {
		t.Fatalf("mismatch verdict = (%+v, %v), want unsigned valid=false", verdict, err)
	}
	request.RunMode = core.RunMode("unknown")
	if _, err := bridge.VerifyStateRoot(context.Background(), request); core.ErrorKindOf(err) != core.ErrorInvalidArgument {
		t.Fatalf("invalid run mode error kind = %v, want InvalidArgument (err %v)", core.ErrorKindOf(err), err)
	}
	request.RunMode = core.RunModeFinalized
	request.BlockNumber = 11
	if _, err := bridge.VerifyStateRoot(context.Background(), request); core.ErrorKindOf(err) != core.ErrorFailedPrecondition {
		t.Fatalf("above-head error kind = %v, want FailedPrecondition (err %v)", core.ErrorKindOf(err), err)
	}

	// A missing signer is a CONFIGURATION error and never clears on its own; the
	// above-head case a few lines up is the one that does. They shared a kind
	// until now, which made the consumer retry a misconfigured daemon forever.
	unsigned := newContractAdapter(t, reader, attestation.Signer{})
	request.BlockNumber = 7
	request.ExpectedStateRoot = stateRoot
	if _, err := unsigned.VerifyStateRoot(context.Background(), request); core.ErrorKindOf(err) != core.ErrorUnimplemented {
		t.Fatalf("unsigned positive error kind = %v, want Unimplemented (err %v)", core.ErrorKindOf(err), err)
	}

	unavailable := newContractAdapter(t, contractHeaderReader{err: errors.New("Nitro unavailable")}, signer)
	if _, err := unavailable.VerifyStateRoot(context.Background(), request); core.ErrorKindOf(err) != core.ErrorUnavailable {
		t.Fatalf("Nitro error kind = %v, want Unavailable (err %v)", core.ErrorKindOf(err), err)
	}
}

func TestNewRejectsInvalidComposition(t *testing.T) {
	if _, err := New(Options{}); err == nil {
		t.Fatal("New accepted empty adapter options")
	}
	runtime, err := arbitrum.NewRuntimeState(contractHeaderReader{})
	if err != nil {
		t.Fatalf("NewRuntimeState: %v", err)
	}
	if _, err := New(Options{SrcChain: "arbitrum-contract", Runtime: runtime, AttestationHead: arbitrum.RunModeFinalized}); err == nil {
		t.Fatal("New accepted nil feed")
	}
}
