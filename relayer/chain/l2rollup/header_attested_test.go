package l2rollup

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"attestor/types/attestation"
	attestorpb "attestor/types/attestor"

	"relayer/chain"

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

func signedVerifier(t *testing.T, header *types.Header) (AttestationVerifier, []byte) {
	t.Helper()
	signer, err := attestation.NewSigner(8453, "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20")
	if err != nil {
		t.Fatalf("create test signer: %v", err)
	}
	verifier, err := NewAttestationVerifier(8453, "finalized", signer.PublicKey())
	if err != nil {
		t.Fatalf("create test verifier: %v", err)
	}
	signature, err := signer.Sign(attestation.RunModeFinalized, header.Number.Uint64(), header.Root[:], header.Hash().Bytes())
	if err != nil {
		t.Fatalf("sign test header: %v", err)
	}
	return verifier, signature
}

// The binding this whole helper exists for: the attestor gates HOW FAR the relayer
// may go, so without a check on WHAT sits at that height an L2 RPC that reorged past
// the frontier hands back a replacement block and the client accepts it. A refusal
// must abort the build rather than package unattested state.
func TestBindToAttestation_RejectsABlockTheAttestorDoesNotRecognise(t *testing.T) {
	at := &fakeAttestor{verifyValid: false}
	b := &attestedHeaderBuilder{attestor: at, runMode: attestorpb.RunMode_RUN_MODE_SAFE, name: "l2-opstack"}

	_, err := b.bindToAttestation(context.Background(), 4096, testL2Header(t))
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
	header := testL2Header(t)
	verifier, signature := signedVerifier(t, header)
	at := &fakeAttestor{verifyValid: true, verifySignature: signature}
	b := &attestedHeaderBuilder{attestor: at, verifier: verifier, srcChain: "arbdev", runMode: attestorpb.RunMode_RUN_MODE_FINALIZED, name: "l2-arbitrum"}

	if _, err := b.bindToAttestation(context.Background(), 4096, header); err != nil {
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
	// A daemon serving several chains verifies against the wrong replica without
	// this, and a wrong replica answers valid=false for a good header.
	if at.gotVerifySrcChain != "arbdev" {
		t.Errorf("src_chain = %q, want the source's own key", at.gotVerifySrcChain)
	}
}

// A transport failure must not be read as approval: the answer is unknown, so the
// build has to fail and be retried rather than proceed unbound.
func TestBindToAttestation_PropagatesTransportFailure(t *testing.T) {
	at := &fakeAttestor{verifyValid: true, verifySignature: make([]byte, 64), verifyErr: errors.New("connection refused")}
	b := &attestedHeaderBuilder{attestor: at, runMode: attestorpb.RunMode_RUN_MODE_SAFE, name: "l2-base"}

	_, err := b.bindToAttestation(context.Background(), 4096, testL2Header(t))
	if err == nil {
		t.Fatal("bindToAttestation swallowed a transport failure, want an error")
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("error = %v, want the transport cause wrapped", err)
	}
}

// An attestor without the signing RPC cannot provide a header the light client
// accepts, so a version skew must fail closed instead of silently relaying bare
// state.
func TestBindToAttestation_UnsupportedAttestorFailsClosed(t *testing.T) {
	at := &fakeAttestor{verifyErr: fmt.Errorf("%w: rpc error", ErrVerifyStateRootUnsupported)}
	b := &attestedHeaderBuilder{attestor: at, runMode: attestorpb.RunMode_RUN_MODE_SAFE, name: "l2-opstack"}

	if _, err := b.bindToAttestation(context.Background(), 4096, testL2Header(t)); err == nil {
		t.Fatal("an attestor without the signing RPC must fail the build")
	}
	if !at.verifyCalled {
		t.Error("VerifyStateRoot was never attempted")
	}
}

func TestBindToAttestation_NoAttestorFailsClosed(t *testing.T) {
	b := &attestedHeaderBuilder{attestor: nil, name: "l2-opstack"}

	if _, err := b.bindToAttestation(context.Background(), 4096, testL2Header(t)); err == nil {
		t.Fatal("bindToAttestation without an attestor must fail")
	}
}

func TestBindToAttestation_RejectsMissingSignature(t *testing.T) {
	at := &fakeAttestor{verifyValid: true}
	b := &attestedHeaderBuilder{attestor: at, runMode: attestorpb.RunMode_RUN_MODE_SAFE, name: "l2-opstack"}

	if _, err := b.bindToAttestation(context.Background(), 4096, testL2Header(t)); err == nil {
		t.Fatal("missing signature must fail")
	}
}

func TestBindToAttestation_RejectsASignatureFromTheWrongAttestor(t *testing.T) {
	header := testL2Header(t)
	verifier, _ := signedVerifier(t, header)
	at := &fakeAttestor{verifyValid: true, verifySignature: make([]byte, 64)}
	b := &attestedHeaderBuilder{attestor: at, verifier: verifier, runMode: attestorpb.RunMode_RUN_MODE_SAFE, name: "l2-opstack"}

	_, err := b.bindToAttestation(context.Background(), 4096, header)
	if err == nil {
		t.Fatal("signature made by the wrong key must fail")
	}
	// Failing is not enough: an unclassified error defaults to TRANSIENT, so the
	// path would retry a key that can never verify, on every flush, forever.
	if !errors.Is(err, ErrAttestorSignature) || !chain.IsPermanent(err) {
		t.Fatalf("err = %v; want ErrAttestorSignature classified permanent", err)
	}
}

// The signing key may serve another relay path at a lower finality level. That
// signature must not be reusable by this client, whose profile pins finalized.
func TestBindToAttestation_RejectsASignatureFromAnotherFinalityLevel(t *testing.T) {
	header := testL2Header(t)
	signer, err := attestation.NewSigner(8453, "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20")
	if err != nil {
		t.Fatalf("create test signer: %v", err)
	}
	verifier, err := NewAttestationVerifier(8453, "finalized", signer.PublicKey())
	if err != nil {
		t.Fatalf("create test verifier: %v", err)
	}
	signature, err := signer.Sign(attestation.RunModeUnsafe, header.Number.Uint64(), header.Root[:], header.Hash().Bytes())
	if err != nil {
		t.Fatalf("sign unsafe header: %v", err)
	}
	at := &fakeAttestor{verifyValid: true, verifySignature: signature}
	b := &attestedHeaderBuilder{attestor: at, verifier: verifier, runMode: attestorpb.RunMode_RUN_MODE_FINALIZED, name: "l2-opstack"}

	_, bindErr := b.bindToAttestation(context.Background(), 4096, header)
	if bindErr == nil {
		t.Fatal("signature issued at unsafe finality was accepted by a finalized verifier")
	}
	if !errors.Is(bindErr, ErrAttestorSignature) || !chain.IsPermanent(bindErr) {
		t.Fatalf("err = %v; want ErrAttestorSignature classified permanent", bindErr)
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
