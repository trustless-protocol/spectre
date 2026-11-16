package gnark

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"testing"

	"relayer/prover/witness"

	"filippo.io/edwards25519"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
)

// dummySeed mirrors witness.dummySeedFmt. The test checks the key it derives
// against witness.PadWithDummies, so a change to the format fails here loudly
// instead of making the collision below silently disappear.
func dummySeed(n, slot int) [32]byte {
	return sha256.Sum256(fmt.Appendf(nil, "spectre-dummy-seed|bucket=%d|slot=%d", n, slot))
}

// signWithDummyNonce signs msg with priv but with the nonce the padding slot
// (n, slot) uses, so R equals that dummy's R. The dummy's seed is public, so any
// validator can do this.
func signWithDummyNonce(t *testing.T, priv ed25519.PrivateKey, msg []byte, n, slot int) []byte {
	t.Helper()
	seed := dummySeed(n, slot)
	dh := sha512.Sum512(seed[:])
	rh := sha512.New()
	rh.Write(dh[32:])
	rh.Write(witness.DummyMsgBytes(uint16(n), uint16(slot)))
	r, err := edwards25519.NewScalar().SetUniformBytes(rh.Sum(nil))
	if err != nil {
		t.Fatalf("nonce: %v", err)
	}
	h := sha512.Sum512(priv.Seed())
	a, err := edwards25519.NewScalar().SetBytesWithClamping(h[:32])
	if err != nil {
		t.Fatalf("clamp: %v", err)
	}
	R := new(edwards25519.Point).ScalarBaseMult(r)
	kh := sha512.New()
	kh.Write(R.Bytes())
	kh.Write(priv.Public().(ed25519.PublicKey))
	kh.Write(msg)
	k, err := edwards25519.NewScalar().SetUniformBytes(kh.Sum(nil))
	if err != nil {
		t.Fatalf("challenge: %v", err)
	}
	return append(R.Bytes(), new(edwards25519.Scalar).MultiplyAdd(k, a, r).Bytes()...)
}

// collidingSigner returns an active signer over vote that repeats a point of
// padding slot (n, slot): its key (same A) or its nonce (same R).
func collidingSigner(t *testing.T, sameR bool, vote []byte, n, slot int) witness.ValidatorSignature {
	t.Helper()
	if !sameR {
		seed := dummySeed(n, slot)
		priv := ed25519.NewKeyFromSeed(seed[:])
		return witness.ValidatorSignature{
			Signature: ed25519.Sign(priv, vote), PublicKey: priv.Public().(ed25519.PublicKey),
			SignedBytes: vote, Active: true,
		}
	}
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return witness.ValidatorSignature{
		Signature: signWithDummyNonce(t, priv, vote, n, slot), PublicKey: priv.Public().(ed25519.PublicKey),
		SignedBytes: vote, Active: true,
	}
}

// TestBatchCircuit_SolvesWhenASignerCollidesWithADummy pins #579 item 4 on the
// current circuit. The padding keys are derived from public seeds, so a validator
// can put a point of padding slot j into its own active slot: register that
// dummy's key (same A twice in the batch), or sign with that dummy's nonce (same R
// twice). On the ecip-gnark build the #579 audit ran, the repeated A made the
// batch unsatisfiable, so every proof over that quorum failed. Since #582 bumped
// ecip-gnark to v0.1.0-rc.4 (with the gnark fork to fork.5) the batch solves, and
// witness.PadWithDummies does not avoid the collision. If a later ecip-gnark
// rejects repeated points again, this fails, and PadWithDummies must re-derive a
// colliding dummy key.
//
// The cases cover the smallest bucket and a 14-signer bucket 16, where two slots
// are padding and the colliding validator sits among the others rather than last.
func TestBatchCircuit_SolvesWhenASignerCollidesWithADummy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}

	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0x50 + i)
	}
	vote := soundnessVote(blockHash)

	cases := []struct {
		name    string
		n       int  // bucket
		signers int  // active slots; slots signers..n-1 are padding
		at      int  // active slot of the colliding validator
		slot    int  // padding slot it collides with
		sameR   bool // repeat the dummy's R instead of its A
	}{
		{name: "bucket 4, 3 signers: same A as padding slot 3", n: 4, signers: 3, at: 2, slot: 3},
		{name: "bucket 4, 3 signers: same R as padding slot 3", n: 4, signers: 3, at: 2, slot: 3, sameR: true},
		{name: "bucket 16, 14 signers: same A as padding slot 14", n: 16, signers: 14, at: 6, slot: 14},
		{name: "bucket 16, 14 signers: same A as padding slot 15", n: 16, signers: 14, at: 6, slot: 15},
		{name: "bucket 16, 14 signers: same R as padding slot 15", n: 16, signers: 14, at: 6, slot: 15, sameR: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			active := make([]witness.ValidatorSignature, tc.signers)
			for i := range active {
				if i == tc.at {
					active[i] = collidingSigner(t, tc.sameR, vote, tc.n, tc.slot)
					continue
				}
				pub, priv, err := ed25519.GenerateKey(rand.Reader)
				if err != nil {
					t.Fatalf("generate key: %v", err)
				}
				active[i] = witness.ValidatorSignature{Signature: ed25519.Sign(priv, vote), PublicKey: pub, SignedBytes: vote, Active: true}
			}
			if !ed25519.Verify(active[tc.at].PublicKey, vote, active[tc.at].Signature) {
				t.Fatal("the colliding signer's signature does not verify")
			}

			padded, err := witness.PadWithDummies(active, tc.n)
			if err != nil {
				t.Fatalf("witness.PadWithDummies: %v", err)
			}
			signer, dummy := padded[tc.at], padded[tc.slot]
			repeated := bytes.Equal(signer.PublicKey, dummy.PublicKey)
			if tc.sameR {
				repeated = bytes.Equal(signer.Signature[:32], dummy.Signature[:32])
			}
			if !signer.Active || dummy.Active || !repeated {
				t.Fatal("no collision: the dummy derivation no longer matches dummySeed, so the test proves nothing")
			}

			hash, err := witness.ComputeWitnessHash(padded)
			if err != nil {
				t.Fatalf("witness.ComputeWitnessHash: %v", err)
			}
			assignment, err := buildBatchAssignment(padded, hash)
			if err != nil {
				t.Fatalf("buildBatchAssignment: %v", err)
			}
			if err := test.IsSolved(NewBatchCircuit(tc.n), assignment, ecc.BN254.ScalarField()); err != nil {
				t.Fatalf("a batch with a point repeated between an active slot and padding is unsatisfiable: %v", err)
			}
		})
	}
}
