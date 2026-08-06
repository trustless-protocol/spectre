package l2rollup

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	attestorpb "attestor/types/attestor"

	"github.com/ethereum/go-ethereum/core/types"
)

// testL2Header is a minimal execution header; only Root and the derived Hash matter
// to the binding check.
func testL2Header(t *testing.T) *types.Header {
	t.Helper()
	return &types.Header{
		Number:     big.NewInt(4096),
		Root:       [32]byte{0xaa, 0xbb},
		ParentHash: [32]byte{0x01},
		GasLimit:   30_000_000,
		GasUsed:    21_000,
	}
}

// The binding this whole helper exists for: the attestor gates HOW FAR the relayer
// may go, so without a check on WHAT sits at that height an L2 RPC that reorged past
// the frontier hands back a replacement block and the client accepts it. A refusal
// must abort the build rather than package unattested state.
func TestBindToAttestation_RejectsABlockTheAttestorDoesNotRecognise(t *testing.T) {
	at := &fakeAttestor{verifyValid: false}
	b := &attestedHeaderBuilder{attestor: at, runMode: attestorpb.RunMode_RUN_MODE_SAFE, name: "l2-opstack"}

	err := b.bindToAttestation(context.Background(), 4096, testL2Header(t))
	if err == nil {
		t.Fatal("bindToAttestation accepted a block the attestor rejected, want an error")
	}
	if !strings.Contains(err.Error(), "does not recognise L2 block 4096") {
		t.Errorf("error = %v, want it to name the refused height", err)
	}
	if !at.verifyCalled {
		t.Error("VerifyStateRoot was never called")
	}
}

// The attestor must be asked about exactly the block being packaged, at the finality
// the source gated the height on — a check against a different height or a laxer head
// would pass headers the source would never have approved.
func TestBindToAttestation_SendsTheBlockIdentityAndRunMode(t *testing.T) {
	at := &fakeAttestor{verifyValid: true}
	b := &attestedHeaderBuilder{attestor: at, runMode: attestorpb.RunMode_RUN_MODE_FINALIZED, name: "l2-arbitrum"}
	header := testL2Header(t)

	if err := b.bindToAttestation(context.Background(), 4096, header); err != nil {
		t.Fatalf("bindToAttestation: %v", err)
	}
	if at.gotVerifyHeight != 4096 {
		t.Errorf("height = %d, want 4096", at.gotVerifyHeight)
	}
	if got, want := at.gotVerifyRoot, header.Root.Bytes(); string(got) != string(want) {
		t.Errorf("state root = %x, want %x", got, want)
	}
	if got, want := at.gotVerifyHash, header.Hash().Bytes(); string(got) != string(want) {
		t.Errorf("block hash = %x, want %x", got, want)
	}
	if at.gotVerifyMode != attestorpb.RunMode_RUN_MODE_FINALIZED {
		t.Errorf("run mode = %v, want RUN_MODE_FINALIZED", at.gotVerifyMode)
	}
}

// A transport failure must not be read as approval: the answer is unknown, so the
// build has to fail and be retried rather than proceed unbound.
func TestBindToAttestation_PropagatesTransportFailure(t *testing.T) {
	at := &fakeAttestor{verifyValid: true, verifyErr: errors.New("connection refused")}
	b := &attestedHeaderBuilder{attestor: at, runMode: attestorpb.RunMode_RUN_MODE_SAFE, name: "l2-base"}

	err := b.bindToAttestation(context.Background(), 4096, testL2Header(t))
	if err == nil {
		t.Fatal("bindToAttestation swallowed a transport failure, want an error")
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("error = %v, want the transport cause wrapped", err)
	}
}

// An attestor binary older than VerifyStateRoot answers Unimplemented. Relayer and
// attestor ship separately, so treating that as a refusal would fail every header
// build against such a deployment; it must degrade to the pre-binding behaviour
// instead — this test is what stops a future "tighten the error handling" change
// from taking a version-skewed deployment offline without warning.
func TestBindToAttestation_UnsupportedAttestorDegradesInsteadOfFailing(t *testing.T) {
	at := &fakeAttestor{verifyErr: fmt.Errorf("%w: rpc error", ErrVerifyStateRootUnsupported)}
	b := &attestedHeaderBuilder{attestor: at, runMode: attestorpb.RunMode_RUN_MODE_SAFE, name: "l2-opstack"}

	if err := b.bindToAttestation(context.Background(), 4096, testL2Header(t)); err != nil {
		t.Fatalf("an attestor without VerifyStateRoot must not fail the build, got %v", err)
	}
	if !at.verifyCalled {
		t.Error("VerifyStateRoot was never attempted")
	}
}

// Skip-finality mode has no attestor at all — RelayableHeight degrades to the raw L2
// head there, so there is no attestation to bind against and the check must not fail
// the build.
func TestBindToAttestation_NoAttestorIsNotAFailure(t *testing.T) {
	b := &attestedHeaderBuilder{attestor: nil, name: "l2-opstack"}

	if err := b.bindToAttestation(context.Background(), 4096, testL2Header(t)); err != nil {
		t.Fatalf("bindToAttestation with no attestor: %v", err)
	}
}

// The run mode handed to the attestor must track the configured head kind, or the
// binding is judged at a different finality than the height gate was.
func TestHeadKindRunMode(t *testing.T) {
	for _, tc := range []struct {
		kind HeadKind
		want attestorpb.RunMode
	}{
		{Unsafe, attestorpb.RunMode_RUN_MODE_UNSAFE},
		{Safe, attestorpb.RunMode_RUN_MODE_SAFE},
		{Finalized, attestorpb.RunMode_RUN_MODE_FINALIZED},
	} {
		if got := tc.kind.RunMode(); got != tc.want {
			t.Errorf("HeadKind(%s).RunMode() = %v, want %v", tc.kind, got, tc.want)
		}
	}
}
