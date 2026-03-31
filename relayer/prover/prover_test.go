package prover

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"testing"

	"0x5ea000000/ecip-gnark/utils"

	curve_bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"filippo.io/edwards25519"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bn254 "github.com/consensys/gnark/backend/groth16/bn254"
)

// TestProveSignature_InvalidSigLength verifies that ProveSignature rejects
// signatures that are not exactly 64 bytes.
func TestProveSignature_InvalidSigLength(t *testing.T) {
	p := &EcipProver{}

	tests := []struct {
		name   string
		sigLen int
	}{
		{"empty signature", 0},
		{"too short", 32},
		{"too long", 128},
		{"off by one short", 63},
		{"off by one long", 65},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sig := make([]byte, tc.sigLen)
			pub := make([]byte, 32)
			msg := []byte("test message")

			_, _, _, err := p.GenerateProof(sig, pub, msg)
			if err == nil {
				t.Fatalf("expected error for sig length %d, got nil", tc.sigLen)
			}

			expected := "invalid signature length"
			if !containsSubstring(err.Error(), expected) {
				t.Errorf("expected error containing %q, got %q", expected, err.Error())
			}
		})
	}
}

// TestProveSignature_InvalidPubLength verifies that ProveSignature rejects
// public keys that are not exactly 32 bytes.
func TestProveSignature_InvalidPubLength(t *testing.T) {
	p := &EcipProver{}

	tests := []struct {
		name   string
		pubLen int
	}{
		{"empty public key", 0},
		{"too short", 16},
		{"too long", 64},
		{"off by one short", 31},
		{"off by one long", 33},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sig := make([]byte, 64)
			pub := make([]byte, tc.pubLen)
			msg := []byte("test message")

			_, _, _, err := p.GenerateProof(sig, pub, msg)
			if err == nil {
				t.Fatalf("expected error for pub length %d, got nil", tc.pubLen)
			}

			expected := "invalid public key length"
			if !containsSubstring(err.Error(), expected) {
				t.Errorf("expected error containing %q, got %q", expected, err.Error())
			}
		})
	}
}

// TestProofToBigInts_NilProof verifies that ProofToBigInts returns an error
// when given a nil proof.
func TestProofToBigInts_NilProof(t *testing.T) {
	_, _, _, err := ProofToBigInts(nil)
	if err == nil {
		t.Fatal("expected error for nil proof, got nil")
	}

	expected := "expected BN254 proof"
	if !containsSubstring(err.Error(), expected) {
		t.Errorf("expected error containing %q, got %q", expected, err.Error())
	}
}

// TestProofToBigInts_WrongType verifies that ProofToBigInts returns an error
// when given a proof that is not a BN254 proof type.
func TestProofToBigInts_WrongType(t *testing.T) {
	// Use a mock proof that implements groth16.Proof but is not *groth16_bn254.Proof
	mockProof := &mockGnarkProof{}

	_, _, _, err := ProofToBigInts(mockProof)
	if err == nil {
		t.Fatal("expected error for wrong proof type, got nil")
	}

	expected := "expected BN254 proof"
	if !containsSubstring(err.Error(), expected) {
		t.Errorf("expected error containing %q, got %q", expected, err.Error())
	}
}

// TestProofToBigInts_ValidBN254Proof verifies that ProofToBigInts successfully
// converts a zero-valued BN254 proof without error and returns non-nil big.Ints.
func TestProofToBigInts_ValidBN254Proof(t *testing.T) {
	p := &groth16_bn254.Proof{}
	// Add one commitment point so Commitments[0] access doesn't panic
	p.Commitments = make([]curve_bn254.G1Affine, 1)

	proof, commitments, commitmentPok, err := ProofToBigInts(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, v := range proof {
		if v == nil {
			t.Errorf("proof[%d] is nil", i)
		}
	}
	for i, v := range commitments {
		if v == nil {
			t.Errorf("commitments[%d] is nil", i)
		}
	}
	for i, v := range commitmentPok {
		if v == nil {
			t.Errorf("commitmentPok[%d] is nil", i)
		}
	}
}

// TestSHA512HashScalar verifies that SHA512(R || A || msg) produces a 64-byte
// digest that can be successfully reduced to an Ed25519 scalar via
// edwards25519.NewScalar().SetUniformBytes().
func TestSHA512HashScalar(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	msg := []byte("test message for SHA512 hash scalar verification")
	sig := ed25519.Sign(priv, msg)

	R := sig[:32]
	A := []byte(pub)

	hasher := sha512.New()
	hasher.Write(R)
	hasher.Write(A)
	hasher.Write(msg)
	sum := hasher.Sum(nil)

	if len(sum) != 64 {
		t.Fatalf("SHA512 digest length = %d, expected 64", len(sum))
	}

	scalar, scalarErr := edwards25519.NewScalar().SetUniformBytes(sum)
	if scalarErr != nil {
		t.Fatalf("SetUniformBytes failed: %v", scalarErr)
	}

	scalarBytes := scalar.Bytes()
	if len(scalarBytes) != 32 {
		t.Fatalf("scalar byte length = %d, expected 32", len(scalarBytes))
	}

	// The scalar must be non-zero for a valid signature
	allZero := true
	for _, b := range scalarBytes {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Error("hash scalar is all zeros, expected non-zero")
	}
}

// TestPreComputationWithRealSignature verifies the pre-computation steps in
// ProveSignature using a real Ed25519 keypair and signature. This exercises
// point decompression, scalar extraction, and hash computation without
// requiring R1CS files for full proof generation.
func TestPreComputationWithRealSignature(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	msg := []byte("pre-computation test message")
	sig := ed25519.Sign(priv, msg)

	// Step 1: Extract R and S from signature
	R := sig[:32]
	S, sErr := edwards25519.NewScalar().SetCanonicalBytes(sig[32:])
	if sErr != nil {
		t.Fatalf("failed to parse S scalar: %v", sErr)
	}

	// Step 2: Decompress public key to Weierstrass coordinates
	aX, aY, decompErr := utils.DecompressPoint([]byte(pub))
	if decompErr != nil {
		t.Fatalf("failed to decompress public key: %v", decompErr)
	}
	if aX == nil || aY == nil {
		t.Fatal("decompressed public key coordinates are nil")
	}
	if aX.Sign() == 0 && aY.Sign() == 0 {
		t.Error("decompressed public key is the point at infinity")
	}

	// Step 3: Decompress R point to Weierstrass coordinates
	rX, rY, decompErr := utils.DecompressPoint(R)
	if decompErr != nil {
		t.Fatalf("failed to decompress R point: %v", decompErr)
	}
	if rX == nil || rY == nil {
		t.Fatal("decompressed R point coordinates are nil")
	}

	// Step 4: Convert S to big.Int
	s := utils.ScalarToBigInt(S)
	if s == nil {
		t.Fatal("ScalarToBigInt returned nil")
	}
	if s.Sign() <= 0 {
		t.Error("S scalar as big.Int should be positive")
	}

	// Step 5: Compute H = SHA512(R || A || msg) and reduce to scalar
	hasher := sha512.New()
	hasher.Write(R)
	hasher.Write([]byte(pub))
	hasher.Write(msg)
	sum := hasher.Sum(nil)

	H, hErr := edwards25519.NewScalar().SetUniformBytes(sum)
	if hErr != nil {
		t.Fatalf("failed to set hash scalar: %v", hErr)
	}

	h := utils.ScalarToBigInt(H)
	if h == nil {
		t.Fatal("hash ScalarToBigInt returned nil")
	}
	if h.Sign() <= 0 {
		t.Error("hash scalar as big.Int should be positive")
	}
}

// TestSHA512HashDeterminism verifies that the SHA512(R || A || msg)
// computation is deterministic: the same inputs always produce the same scalar.
func TestSHA512HashDeterminism(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	msg := []byte("determinism check")
	sig := ed25519.Sign(priv, msg)
	R := sig[:32]

	computeScalar := func() []byte {
		hasher := sha512.New()
		hasher.Write(R)
		hasher.Write([]byte(pub))
		hasher.Write(msg)
		sum := hasher.Sum(nil)

		scalar, scalarErr := edwards25519.NewScalar().SetUniformBytes(sum)
		if scalarErr != nil {
			t.Fatalf("SetUniformBytes failed: %v", scalarErr)
		}
		return scalar.Bytes()
	}

	first := computeScalar()
	second := computeScalar()

	if len(first) != len(second) {
		t.Fatalf("scalar lengths differ: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("scalars differ at byte %d: %02x vs %02x", i, first[i], second[i])
		}
	}
}

// mockGnarkProof is a minimal implementation of groth16.Proof for testing
// type assertion failures in ProofToBigInts.
type mockGnarkProof struct {
	groth16.Proof
}

// containsSubstring reports whether s contains substr.
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
