package gnark

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"strings"
	"testing"

	"0x5ea000000/ecip-gnark/utils"

	"filippo.io/edwards25519"
	curve_bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bn254 "github.com/consensys/gnark/backend/groth16/bn254"
)

// Carried from the pre-rewrite suite: coverage over the Ed25519 -> gnark
// pre-computation path and proof serialization, both part of the carried
// circuit/assignment code.

func TestPreComputationWithRealSignature(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	msg := []byte("pre-computation test message")
	sig := ed25519.Sign(priv, msg)

	R := sig[:32]
	S, err := edwards25519.NewScalar().SetCanonicalBytes(sig[32:])
	if err != nil {
		t.Fatalf("parse S: %v", err)
	}
	aX, aY, err := utils.DecompressPoint([]byte(pub))
	if err != nil {
		t.Fatalf("decompress pub: %v", err)
	}
	if aX == nil || aY == nil {
		t.Fatal("decompressed pub nil")
	}
	rX, rY, err := utils.DecompressPoint(R)
	if err != nil {
		t.Fatalf("decompress R: %v", err)
	}
	if rX == nil || rY == nil {
		t.Fatal("decompressed R nil")
	}
	if utils.ScalarToBigInt(S).Sign() <= 0 {
		t.Error("S big.Int should be positive")
	}
	hasher := sha512.New()
	hasher.Write(R)
	hasher.Write([]byte(pub))
	hasher.Write(msg)
	if len(hasher.Sum(nil)) != 64 {
		t.Fatal("sha512 output not 64 bytes")
	}
}

func TestProofToBigInts_NilProof(t *testing.T) {
	_, _, _, err := ProofToBigInts(nil)
	if err == nil || !strings.Contains(err.Error(), "expected BN254 proof") {
		t.Fatalf("want 'expected BN254 proof' error, got %v", err)
	}
}

func TestProofToBigInts_WrongType(t *testing.T) {
	_, _, _, err := ProofToBigInts(&mockGnarkProof{})
	if err == nil || !strings.Contains(err.Error(), "expected BN254 proof") {
		t.Fatalf("want error, got %v", err)
	}
}

func TestProofToBigInts_ValidBN254Proof(t *testing.T) {
	p := &groth16_bn254.Proof{}
	p.Commitments = make([]curve_bn254.G1Affine, 1)
	proofInts, commitments, commitmentPok, err := ProofToBigInts(p)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	for i, v := range proofInts {
		if v == nil {
			t.Errorf("proof[%d] nil", i)
		}
	}
	for i, v := range commitments {
		if v == nil {
			t.Errorf("commitments[%d] nil", i)
		}
	}
	for i, v := range commitmentPok {
		if v == nil {
			t.Errorf("commitmentPok[%d] nil", i)
		}
	}
}

type mockGnarkProof struct {
	groth16.Proof
}
