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
	shared := SharedBlockData{
		BlockIDHash: make([]byte, 32),
		PartSetHash: make([]byte, 32),
		ChainID:     "test",
	}

	t.Run("empty inputs", func(t *testing.T) {
		_, _, _, _, err := p.GenerateProof(shared, nil)
		if err == nil || !strings.Contains(err.Error(), "no signatures") {
			t.Fatalf("want 'no signatures' error, got %v", err)
		}
	})

	t.Run("bad blockIDHash", func(t *testing.T) {
		bad := shared
		bad.BlockIDHash = make([]byte, 31)
		_, _, _, _, err := p.GenerateProof(bad, []ValidatorSignature{{Signature: make([]byte, 64), PublicKey: make([]byte, 32)}})
		if err == nil || !strings.Contains(err.Error(), "BlockIDHash") {
			t.Fatalf("want BlockIDHash error, got %v", err)
		}
	})

	t.Run("exceeds max bucket", func(t *testing.T) {
		n := MaxBucket() + 1
		sigs := make([]ValidatorSignature, n)
		_, _, _, _, err := p.GenerateProof(shared, sigs)
		if err == nil || !strings.Contains(err.Error(), "largest bucket") {
			t.Fatalf("want over-max error, got %v", err)
		}
	})
}

func TestSmallestBucketGEQ(t *testing.T) {
	cases := []struct {
		in   int
		want int
		err  bool
	}{
		{1, 4, false},
		{4, 4, false},
		{5, 8, false},
		{17, 32, false},
		{128, 128, false},
		{129, 0, true},
	}
	for _, tc := range cases {
		got, err := SmallestBucketGEQ(tc.in)
		if tc.err {
			if err == nil {
				t.Errorf("in=%d: expected error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("in=%d: unexpected error %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("in=%d: got bucket %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestPadSigsToBucket_FillsWithSlotZero(t *testing.T) {
	sigs := []ValidatorSignature{{
		Signature:        append(make([]byte, 63), 0xAB),
		PublicKey:        append(make([]byte, 31), 0xCD),
		Index:            7,
		Power:            100,
		TimestampSeconds: 42,
		TimestampNanos:   7,
	}}
	padded, err := padSigsToBucket(sigs, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(padded) != 4 {
		t.Fatalf("expected 4 slots, got %d", len(padded))
	}
	for i, s := range padded {
		if s.Index != 7 || s.Power != 100 || s.TimestampSeconds != 42 {
			t.Errorf("slot[%d] not a copy of slot 0: %+v", i, s)
		}
	}
}

func TestPadSigsToBucket_TooBig(t *testing.T) {
	_, err := padSigsToBucket(make([]ValidatorSignature, 5), 4)
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
