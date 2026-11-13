package gnark

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"strings"
	"testing"

	"github.com/decentrio/ecip-gnark/utils"

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
	p.Commitments = make([]curve_bn254.G1Affine, commitmentCount)
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

// The order is the one gnark's MarshalSolidity writes and the contract reads: C0.x, C0.y,
// C1.x, C1.y, so the two commitments cannot be swapped or collapsed into one.
func TestProofToBigInts_CommitmentsKeepBothPointsInOrder(t *testing.T) {
	p := &groth16_bn254.Proof{}
	p.Commitments = make([]curve_bn254.G1Affine, commitmentCount)
	for i := range p.Commitments {
		p.Commitments[i].X.SetUint64(uint64(10*(i+1) + 1))
		p.Commitments[i].Y.SetUint64(uint64(10*(i+1) + 2))
	}
	p.CommitmentPok.X.SetUint64(91)
	p.CommitmentPok.Y.SetUint64(92)

	_, commitments, commitmentPok, err := ProofToBigInts(p)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	for i, want := range []int64{11, 12, 21, 22} {
		if commitments[i].Int64() != want {
			t.Errorf("commitments[%d] = %v, want %d", i, commitments[i], want)
		}
	}
	if commitmentPok[0].Int64() != 91 || commitmentPok[1].Int64() != 92 {
		t.Errorf("commitmentPok = %v, want [91 92]", commitmentPok)
	}
}

// A circuit change that adds or drops a commitment must fail here, not ship a proof the
// verifier contract reads with the wrong layout.
func TestProofToBigInts_RejectsAWrongCommitmentCount(t *testing.T) {
	for _, n := range []int{0, 1, commitmentCount + 1} {
		p := &groth16_bn254.Proof{}
		p.Commitments = make([]curve_bn254.G1Affine, n)
		_, _, _, err := ProofToBigInts(p)
		if err == nil || !strings.Contains(err.Error(), "commitments") {
			t.Errorf("%d commitments: want a commitment-count error, got %v", n, err)
		}
	}
}

type mockGnarkProof struct {
	groth16.Proof
}
