package l2rollup

import (
	"context"
	"errors"
	"testing"

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

	gotSrcChain      string
	gotProvisional   bool
	gotAtOrBelowArg  uint64
	atOrBelowCalled  bool
	attestedUpCalled bool
}

func (f *fakeAttestor) AttestedUpTo(_ context.Context, srcChain string, includeProvisional bool) (*attestorpb.AttestedRoot, bool, error) {
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
