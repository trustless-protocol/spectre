package prover

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
)

// makeCanonicalVote builds a canonical-vote-shaped message in the #199 layout:
//
//	uvarint(bodyLen) || Type(0x08 0x02) || Height(0x11 + 8B) || [Round(0x19 + 8B)]
//	  || BlockID(0x22 0x48 0x0a 0x20 hash[32]) || trailing padding
//
// trailLen drives the total body length so callers can deterministically pick a
// 1-byte (bodyLen < 128) or 2-byte leading varint. The Type|Height prefix head
// and the 32-byte block hash content stay fixed regardless of varint width, so
// votes built with different trailLen still share the committed common prefix.
func makeCanonicalVote(roundPresent bool, trailLen int, blockHash [32]byte) []byte {
	body := []byte{0x08, 0x02} // Type = precommit (field 1)
	body = append(body, 0x11)  // Height tag (field 2, sfixed64)
	body = append(body, 1, 2, 3, 4, 5, 6, 7, 8)
	if roundPresent {
		body = append(body, 0x19) // Round tag (field 3, fixed64)
		body = append(body, 0, 0, 0, 0, 0, 0, 0, 1)
	}
	body = append(body, 0x22, 0x48, 0x0a, 0x20) // BlockID tag/len + inner hash tag/len
	body = append(body, blockHash[:]...)
	body = append(body, bytes.Repeat([]byte{0x2a}, trailLen)...) // timestamp/chainID stand-in
	prefix := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(prefix, uint64(len(body)))
	return append(prefix[:n], body...)
}

func wantPrefixHead() [prefixHeadLen]byte {
	return [prefixHeadLen]byte{0x08, 0x02, 0x11, 1, 2, 3, 4, 5, 6, 7, 8}
}

func TestCommonPrefixFields_VarintWidths(t *testing.T) {
	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0xA0 + i)
	}
	wantHead := wantPrefixHead()

	cases := []struct {
		name         string
		roundPresent bool
		trailLen     int
		wantTwoByte  bool
	}{
		{"1-byte varint, no round", false, 10, false}, // bodyLen = 47 + 10 = 57
		{"2-byte varint, no round", false, 90, true},  // bodyLen = 47 + 90 = 137
		{"1-byte varint, round", true, 10, false},     // bodyLen = 56 + 10 = 66
		{"2-byte varint, round", true, 90, true},      // bodyLen = 56 + 90 = 146
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vote := makeCanonicalVote(tc.roundPresent, tc.trailLen, blockHash)
			if gotTwoByte := vote[0] >= 0x80; gotTwoByte != tc.wantTwoByte {
				t.Fatalf("varint width: got 2-byte=%v, want %v (sb[0]=%#x)", gotTwoByte, tc.wantTwoByte, vote[0])
			}
			head, hash, round, err := commonPrefixFields([]ValidatorSignature{{SignedBytes: vote, Active: true}})
			if err != nil {
				t.Fatalf("commonPrefixFields: %v", err)
			}
			if head != wantHead {
				t.Errorf("prefix head = %x, want %x", head, wantHead)
			}
			if hash != blockHash {
				t.Errorf("block hash = %x, want %x", hash, blockHash)
			}
			if round != tc.roundPresent {
				t.Errorf("roundPresent = %v, want %v", round, tc.roundPresent)
			}
		})
	}
}

// TestCommonPrefixFields_SkipsInactive confirms derivation comes from the first
// ACTIVE signer (padding/dummy slots are skipped), even when an inactive slot
// with a different varint width leads.
func TestCommonPrefixFields_SkipsInactive(t *testing.T) {
	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(i)
	}
	dummy := makeCanonicalVote(false, 90, [32]byte{}) // unrelated 2-byte vote, inactive
	real := makeCanonicalVote(false, 10, blockHash)   // 1-byte vote, active
	head, hash, _, err := commonPrefixFields([]ValidatorSignature{
		{SignedBytes: dummy, Active: false},
		{SignedBytes: real, Active: true},
	})
	if err != nil {
		t.Fatalf("commonPrefixFields: %v", err)
	}
	if want := wantPrefixHead(); head != want {
		t.Errorf("prefix head = %x, want %x", head, want)
	}
	if hash != blockHash {
		t.Errorf("block hash = %x, want %x", hash, blockHash)
	}
}

// TestBatchCircuit_MixedVarintWidths solves the circuit (no trusted setup) over
// a bucket whose active signers mix 1-byte and 2-byte leading varints over the
// SAME block — the case the original 2-byte-only offset assumption broke on.
// Padding slots are filled with deterministic dummies, exactly as the prover does.
func TestBatchCircuit_MixedVarintWidths(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}
	const n = 4
	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0x10 + i)
	}
	// Two active signers over the same block, different varint widths.
	votes := [][]byte{
		makeCanonicalVote(false, 10, blockHash), // 1-byte varint
		makeCanonicalVote(false, 90, blockHash), // 2-byte varint
	}
	sigs := make([]ValidatorSignature, 0, n)
	for _, vote := range votes {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		sigs = append(sigs, ValidatorSignature{
			Signature:   ed25519.Sign(priv, vote),
			PublicKey:   pub,
			SignedBytes: vote,
			Active:      true,
		})
	}

	padded, err := padWithDummies(sigs, generateDummySlots(n))
	if err != nil {
		t.Fatalf("padWithDummies: %v", err)
	}
	hash, err := ComputeWitnessHash(padded)
	if err != nil {
		t.Fatalf("ComputeWitnessHash: %v", err)
	}
	assignment, err := buildBatchAssignment(padded, hash)
	if err != nil {
		t.Fatalf("buildBatchAssignment: %v", err)
	}
	if err := test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField()); err != nil {
		t.Fatalf("circuit not satisfied for mixed varint widths: %v", err)
	}
}
