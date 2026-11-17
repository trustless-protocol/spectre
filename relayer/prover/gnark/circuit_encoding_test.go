package gnark

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"testing"

	"relayer/prover/witness"

	"filippo.io/edwards25519"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
)

// lowOrderSig returns a signature crypto/ed25519.Verify accepts under a public
// key A of low order (or the identity): R = [s]B, S = s verifies whenever
// [k]A = O, which happens for one k in ord(A). It retries s until it does.
func lowOrderSig(t *testing.T, pub, msg []byte) []byte {
	t.Helper()
	for i := 0; i < 1024; i++ {
		var buf [64]byte
		if _, err := rand.Read(buf[:]); err != nil {
			t.Fatalf("rand: %v", err)
		}
		s, err := edwards25519.NewScalar().SetUniformBytes(buf[:])
		if err != nil {
			t.Fatalf("scalar: %v", err)
		}
		sig := append(new(edwards25519.Point).ScalarBaseMult(s).Bytes(), s.Bytes()...)
		if ed25519.Verify(pub, msg, sig) {
			return sig
		}
	}
	t.Fatalf("no signature crypto/ed25519 accepts under %x", pub)
	return nil
}

// zeroNonceSig signs msg with a fresh key and nonce r = 0, so R is the identity:
// S = k*a, and [S]B = [k]A = R + [k]A.
func zeroNonceSig(t *testing.T, msg []byte) (pub, sig []byte) {
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
	return pub, sig
}

func compressed(t *testing.T, first, fill, last byte) []byte {
	t.Helper()
	out := make([]byte, 32)
	for i := range out {
		out[i] = fill
	}
	out[0], out[31] = first, last
	return out
}

// TestCheckCircuitEncodingMatchesCircuit is the #579 item 3 regression. Every
// case is a signature crypto/ed25519.Verify accepts, signed by an active slot
// beside one honest signer, and the case asserts that witness.CheckCircuitEncoding
// and the real BatchCircuit solver agree on it. Before the check existed, the
// relayer selected the rejected shapes like any other signature, and every proof
// over that quorum failed in the prover.
func TestCheckCircuitEncodingMatchesCircuit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}
	const n = 4

	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0x40 + i)
	}
	vote := soundnessVote(blockHash)

	order8, err := hex.DecodeString("26e8958fc2b227b045c3f489f2ef98f0d5dfac05d3c63339b13802886d53fc05")
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	p, err := new(edwards25519.Point).SetBytes(order8)
	if err != nil {
		t.Fatalf("order-8 point does not decode: %v", err)
	}
	p4 := new(edwards25519.Point).Add(p, p)
	p4.Add(p4, p4)
	p8 := new(edwards25519.Point).Add(p4, p4)
	if p4.Equal(edwards25519.NewIdentityPoint()) == 1 || p8.Equal(edwards25519.NewIdentityPoint()) != 1 {
		t.Fatal("the order-8 vector is not of order 8")
	}

	honestPub, honestPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	zeroPub, zeroSig := zeroNonceSig(t, vote)

	type signedCase struct {
		name     string
		pub, sig []byte
	}
	lowOrder := func(name string, pub []byte) signedCase {
		return signedCase{name, pub, lowOrderSig(t, pub, vote)}
	}
	cases := []signedCase{
		{"honest", honestPub, ed25519.Sign(honestPriv, vote)},
		{"R identity (nonce r = 0)", zeroPub, zeroSig},
		lowOrder("A identity", compressed(t, 0x01, 0x00, 0x00)),
		lowOrder("A order-2 point (0, -1)", compressed(t, 0xec, 0xff, 0x7f)),
		lowOrder("A y = p + 1 (identity, non-canonical)", compressed(t, 0xee, 0xff, 0x7f)),
		lowOrder("A y = p (order 4, non-canonical)", compressed(t, 0xed, 0xff, 0x7f)),
		lowOrder("A order-4 point y = 0", compressed(t, 0x00, 0x00, 0x00)),
		lowOrder("A order-8 point", order8),
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !ed25519.Verify(tc.pub, vote, tc.sig) {
				t.Fatal("crypto/ed25519 rejects the case; it must accept it for the test to mean anything")
			}
			checkErr := witness.CheckCircuitEncoding(tc.pub, tc.sig)

			otherPub, otherPriv, err := ed25519.GenerateKey(rand.Reader)
			if err != nil {
				t.Fatalf("generate key: %v", err)
			}
			padded, err := witness.PadWithDummies([]witness.ValidatorSignature{
				{Signature: ed25519.Sign(otherPriv, vote), PublicKey: otherPub, SignedBytes: vote, Active: true},
				{Signature: tc.sig, PublicKey: tc.pub, SignedBytes: vote, Active: true},
			}, n)
			if err != nil {
				t.Fatalf("witness.PadWithDummies: %v", err)
			}
			hash, err := witness.ComputeWitnessHash(padded)
			if err != nil {
				t.Fatalf("witness.ComputeWitnessHash: %v", err)
			}
			assignment, err := buildBatchAssignment(padded, hash)
			if err != nil {
				t.Fatalf("buildBatchAssignment: %v", err)
			}
			solveErr := test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField())

			switch {
			case checkErr == nil && solveErr != nil:
				t.Fatalf("CheckCircuitEncoding accepts the signature but the circuit rejects it: %v", solveErr)
			case checkErr != nil && solveErr == nil:
				t.Fatalf("CheckCircuitEncoding rejects the signature (%v) but the circuit proves it", checkErr)
			}
		})
	}
}
