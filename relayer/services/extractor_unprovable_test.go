package services

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"strings"
	"testing"

	relayerclient "relayer/client"
	"relayer/prover/witness"

	cmted25519 "github.com/cometbft/cometbft/crypto/ed25519"

	"filippo.io/edwards25519"
)

// zeroNonceSign signs msg with a fresh key and nonce r = 0: R is the identity
// and S = k*a. crypto/ed25519 accepts it; the proof circuit does not (#579).
func zeroNonceSign(t *testing.T, msg []byte) (pub, sig []byte) {
	t.Helper()
	seed := make([]byte, ed25519.SeedSize)
	if _, err := rand.Read(seed); err != nil {
		t.Fatalf("rand: %v", err)
	}
	h := sha512.Sum512(seed)
	a, err := edwards25519.NewScalar().SetBytesWithClamping(h[:32])
	if err != nil {
		t.Fatalf("clamp: %v", err)
	}
	A := new(edwards25519.Point).ScalarBaseMult(a)
	R := edwards25519.NewIdentityPoint()
	kh := sha512.New()
	kh.Write(R.Bytes())
	kh.Write(A.Bytes())
	kh.Write(msg)
	k, err := edwards25519.NewScalar().SetUniformBytes(kh.Sum(nil))
	if err != nil {
		t.Fatalf("challenge: %v", err)
	}
	pub = A.Bytes()
	sig = append(R.Bytes(), new(edwards25519.Scalar).Multiply(k, a).Bytes()...)
	if !ed25519.Verify(pub, msg, sig) {
		t.Fatal("crypto/ed25519 rejects the r = 0 signature; the test would prove nothing")
	}
	return pub, sig
}

// rekeyWithZeroNonce replaces validator idx's key with a fresh one and its
// commit signature with an r = 0 signature over the same vote.
func rekeyWithZeroNonce(t *testing.T, lb *relayerclient.LightBlock, chainID string, idx int) cmted25519.PubKey {
	t.Helper()
	commit := lb.SignedHeader.Commit
	pub, sig := zeroNonceSign(t, commit.VoteSignBytes(chainID, int32(idx)))
	cmtPub := cmted25519.PubKey(pub)
	lb.ValSet.Validators[idx].PubKey = cmtPub
	lb.ValSet.Validators[idx].Address = cmtPub.Address()
	commit.Signatures[idx].ValidatorAddress = cmtPub.Address()
	commit.Signatures[idx].Signature = sig
	return cmtPub
}

// freshReportedKeys gives the test an empty set in place of the package-level
// one *target points at, and restores that set when the test ends. The tests
// that use it assert on the first report's WARN level, which depends on what
// earlier tests in the process reported. They are not parallel, so no other
// test calls ExtractValidatorSignatures while the variable is swapped.
func freshReportedKeys(t *testing.T, target **reportedKeys) {
	t.Helper()
	previous := *target
	*target = &reportedKeys{keys: make(map[string]struct{})}
	t.Cleanup(func() { *target = previous })
}

// TestExtractValidatorSignatures_SkipsSignerTheCircuitRejects is the #579 item 3
// regression. The highest-power validator signs with nonce r = 0: the signature
// verifies, so the greedy selection used to take it first, and every proof over
// the quorum failed. It must be skipped, reported once at WARN, and the quorum
// built from the other validators.
func TestExtractValidatorSignatures_SkipsSignerTheCircuitRejects(t *testing.T) {
	freshReportedKeys(t, &unprovableSigners)
	chainID := "test-chain"
	// Sorted by power, the 30 is index 0. Without it 70 of 100 still sign.
	lb := makeSignedCommit(t, chainID, []int64{30, 25, 25, 20}, []bool{true, true, true, true})
	bad := rekeyWithZeroNonce(t, lb, chainID, 0)

	var got *ExtractorResult
	var err error
	logs := captureLog(t, func() { got, err = ExtractValidatorSignatures(lb, chainID, 4) })
	if err != nil {
		t.Fatalf("ExtractValidatorSignatures: %v", err)
	}
	for _, s := range append(append([]witness.ValidatorSignature{}, got.Candidates...), got.Signatures...) {
		if s.Index == 0 {
			t.Fatal("the validator the circuit rejects was kept")
		}
	}
	if len(got.Signatures) != 3 {
		t.Fatalf("selected %d signers, want the other 3", len(got.Signatures))
	}

	wantWarn := fmt.Sprintf("WARN [ExtractValidatorSignatures] skipping validator %s (index 0, power 30, pubkey %X) at height 100",
		bad.Address(), bad.Bytes())
	if !strings.Contains(logs, wantWarn) || !strings.Contains(logs, "signature R: point has x = 0") {
		t.Fatalf("missing warning, want %q in:\n%s", wantWarn, logs)
	}

	// The validator signs every block; the next extraction says nothing new.
	logs = captureLog(t, func() { _, err = ExtractValidatorSignatures(lb, chainID, 4) })
	if err != nil {
		t.Fatalf("second ExtractValidatorSignatures: %v", err)
	}
	if !strings.HasPrefix(logs, "DEBUG [ExtractValidatorSignatures] skipping validator "+bad.Address().String()) {
		t.Fatalf("a repeat must log at DEBUG, got:\n%s", logs)
	}
}

// TestExtractValidatorSignatures_NoQuorumWithoutTheRejectedSigner: when the
// remaining signers fall short, extraction fails instead of selecting the
// signature it knows cannot be proven.
func TestExtractValidatorSignatures_NoQuorumWithoutTheRejectedSigner(t *testing.T) {
	chainID := "test-chain"
	// Without the 40 only 60 of 100 sign, below the 67 quorum.
	lb := makeSignedCommit(t, chainID, []int64{40, 30, 20, 10}, []bool{true, true, true, true})
	rekeyWithZeroNonce(t, lb, chainID, 0)

	var err error
	captureLog(t, func() { _, err = ExtractValidatorSignatures(lb, chainID, 4) })
	if err == nil || !strings.Contains(err.Error(), "insufficient voting power") {
		t.Fatalf("err = %v, want insufficient voting power", err)
	}
}

func TestReportedKeys_FirstReportOncePerKey(t *testing.T) {
	r := &reportedKeys{keys: make(map[string]struct{})}
	if !r.firstReport([]byte("a")) {
		t.Fatal("first sighting of a not reported")
	}
	if r.firstReport([]byte("a")) {
		t.Fatal("second sighting of a reported again")
	}
	if !r.firstReport([]byte("b")) {
		t.Fatal("first sighting of b not reported")
	}
}

func TestReportedKeys_ClearsAtTheBound(t *testing.T) {
	r := &reportedKeys{keys: make(map[string]struct{})}
	for i := 0; i < maxReportedKeys; i++ {
		r.firstReport([]byte(fmt.Sprint(i)))
	}
	// The set is full: the next new key clears it, so key 0 is new again.
	r.firstReport([]byte("overflow"))
	if len(r.keys) != 1 {
		t.Fatalf("set holds %d keys after reaching the bound, want 1", len(r.keys))
	}
	if !r.firstReport([]byte("0")) {
		t.Fatal("a key from before the clear was not reported again")
	}
}

// rekeyWithDummyKey gives validator idx the public padding key of (bucket, slot)
// and signs its vote with it. The seed format mirrors witness.dummySeedFmt; the
// test checks witness.DummyKeySlot recognises the key, so a format change fails
// loudly instead of silently testing nothing.
func rekeyWithDummyKey(t *testing.T, lb *relayerclient.LightBlock, chainID string, idx, bucket, slot int) cmted25519.PubKey {
	t.Helper()
	seed := sha256.Sum256(fmt.Appendf(nil, "spectre-dummy-seed|bucket=%d|slot=%d", bucket, slot))
	priv := ed25519.NewKeyFromSeed(seed[:])
	cmtPub := cmted25519.PubKey(priv.Public().(ed25519.PublicKey))
	if b, s, ok := witness.DummyKeySlot(cmtPub.Bytes()); !ok || b != bucket || s != slot {
		t.Fatalf("witness.DummyKeySlot does not recognise the bucket %d slot %d key: (%d, %d, %v)", bucket, slot, b, s, ok)
	}
	commit := lb.SignedHeader.Commit
	lb.ValSet.Validators[idx].PubKey = cmtPub
	lb.ValSet.Validators[idx].Address = cmtPub.Address()
	commit.Signatures[idx].ValidatorAddress = cmtPub.Address()
	commit.Signatures[idx].Signature = ed25519.Sign(priv, commit.VoteSignBytes(chainID, int32(idx)))
	return cmtPub
}

// TestExtractValidatorSignatures_WarnsAboutADummyKeySigner covers #579 item 4.
// A validator keyed with a public padding key still proves, so it stays selected,
// but the relayer names it once at WARN and at DEBUG after that.
func TestExtractValidatorSignatures_WarnsAboutADummyKeySigner(t *testing.T) {
	// The padding key is fixed, so under -count=N or after another test that
	// reports it, the package-level set would already hold it and log DEBUG.
	freshReportedKeys(t, &dummyKeySigners)
	chainID := "test-chain"
	lb := makeSignedCommit(t, chainID, []int64{30, 25, 25, 20}, []bool{true, true, true, true})
	bad := rekeyWithDummyKey(t, lb, chainID, 0, 16, 15)

	var got *ExtractorResult
	var err error
	logs := captureLog(t, func() { got, err = ExtractValidatorSignatures(lb, chainID, 4) })
	if err != nil {
		t.Fatalf("ExtractValidatorSignatures: %v", err)
	}
	selected := false
	for _, s := range got.Signatures {
		selected = selected || s.Index == 0
	}
	if !selected {
		t.Fatal("the validator with the padding key was dropped; its signature still proves and must stay selectable")
	}
	wantWarn := fmt.Sprintf("WARN [ExtractValidatorSignatures] validator %s (index 0, power 30, pubkey %X) at height 100 "+
		"uses the public padding key of bucket 16 slot 15", bad.Address(), bad.Bytes())
	if !strings.Contains(logs, wantWarn) {
		t.Fatalf("missing warning, want %q in:\n%s", wantWarn, logs)
	}

	logs = captureLog(t, func() { _, err = ExtractValidatorSignatures(lb, chainID, 4) })
	if err != nil {
		t.Fatalf("second ExtractValidatorSignatures: %v", err)
	}
	if !strings.HasPrefix(logs, "DEBUG [ExtractValidatorSignatures] validator "+bad.Address().String()) {
		t.Fatalf("a repeat must log at DEBUG, got:\n%s", logs)
	}
}

// TestExtractValidatorSignatures_NoWarningForHonestKeys: an ordinary commit logs
// nothing about padding keys or unprovable signers.
func TestExtractValidatorSignatures_NoWarningForHonestKeys(t *testing.T) {
	chainID := "test-chain"
	lb := makeSignedCommit(t, chainID, []int64{30, 25, 25, 20}, []bool{true, true, true, true})
	var err error
	logs := captureLog(t, func() { _, err = ExtractValidatorSignatures(lb, chainID, 4) })
	if err != nil {
		t.Fatalf("ExtractValidatorSignatures: %v", err)
	}
	if strings.Contains(logs, "[ExtractValidatorSignatures]") {
		t.Fatalf("honest validators were reported:\n%s", logs)
	}
}
