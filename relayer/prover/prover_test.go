package prover

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

// TestGenerateProof_InputValidation exercises the cheap preconditions that
// run before we touch any circuit artifacts. It lets us verify the batch API
// surface without needing compiled r1cs/pk/vk on disk.
func TestGenerateProof_InputValidation(t *testing.T) {
	p := &EcipProver{byBucket: map[int]*bucketArtifacts{}}

	t.Run("empty inputs", func(t *testing.T) {
		_, _, _, _, _, err := p.GenerateProof(nil)
		if err == nil || !strings.Contains(err.Error(), "no signatures") {
			t.Fatalf("want 'no signatures' error, got %v", err)
		}
	})

	t.Run("exceeds max bucket", func(t *testing.T) {
		n := MaxBucket() + 1
		sigs := make([]ValidatorSignature, n)
		_, _, _, _, _, err := p.GenerateProof(sigs)
		if err == nil || !strings.Contains(err.Error(), "largest bucket") {
			t.Fatalf("want over-max error, got %v", err)
		}
	})
}

func TestGenerateMisbehaviourProof_InputValidation(t *testing.T) {
	p := &EcipProver{byBucket: map[int]*bucketArtifacts{}}

	_, err := p.GenerateMisbehaviourProof(nil, []ValidatorSignature{{Active: true}})
	if err == nil || !strings.Contains(err.Error(), "header1 proof: no signatures") {
		t.Fatalf("want header1 no-signatures error, got %v", err)
	}

	_, err = p.GenerateMisbehaviourProof([]ValidatorSignature{{Active: true}}, nil)
	if err == nil || !strings.Contains(err.Error(), "header2 proof: no signatures") {
		t.Fatalf("want header2 no-signatures error, got %v", err)
	}
}

func TestSmallestBucketGEQ(t *testing.T) {
	if len(Buckets) == 0 {
		t.Fatal("Buckets must not be empty")
	}

	for i, bucket := range Buckets {
		for n := 1; n <= bucket; n++ {
			if i > 0 && n <= Buckets[i-1] {
				continue
			}
			got, err := SmallestBucketGEQ(n)
			if err != nil {
				t.Fatalf("n=%d: unexpected error %v", n, err)
			}
			if got != bucket {
				t.Fatalf("n=%d: got bucket %d, want %d", n, got, bucket)
			}
		}
	}

	_, err := SmallestBucketGEQ(MaxBucket() + 1)
	if err == nil {
		t.Fatalf("n=%d: expected error", MaxBucket()+1)
	}
}

func TestPadWithDummies_FillsTrailingSlotsWithDistinctDummies(t *testing.T) {
	real := []ValidatorSignature{{
		Signature: append(make([]byte, 63), 0xAB),
		PublicKey: append(make([]byte, 31), 0xCD),
		Index:     7,
		Power:     100,
		Active:    true,
	}}
	dummies := generateDummySlots(4)
	padded, err := padWithDummies(real, dummies)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(padded) != 4 {
		t.Fatalf("expected 4 slots, got %d", len(padded))
	}
	if !padded[0].Active || padded[0].Index != 7 {
		t.Errorf("slot 0 should be the real signer, got %+v", padded[0])
	}
	for i := 1; i < 4; i++ {
		if padded[i].Active {
			t.Errorf("slot %d should be inactive padding", i)
		}
	}
	// Distinctness: every padding pubkey differs.
	seen := map[string]bool{string(padded[0].PublicKey): true}
	for i := 1; i < 4; i++ {
		key := string(padded[i].PublicKey)
		if seen[key] {
			t.Errorf("duplicate pubkey at slot %d", i)
		}
		seen[key] = true
	}
}

func TestPadWithDummies_TooBig(t *testing.T) {
	_, err := padWithDummies(make([]ValidatorSignature, 5), generateDummySlots(4))
	if err == nil {
		t.Fatal("expected too-small bucket error")
	}
}

// TestPreComputationWithRealSignature keeps coverage over the Ed25519 → gnark
// field conversions we still rely on inside buildBatchAssignment.
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
	proof, commitments, commitmentPok, err := ProofToBigInts(p)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	for i, v := range proof {
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
