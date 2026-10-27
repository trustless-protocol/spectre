package witness

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestDigestPublicInputs(t *testing.T) {
	hashBytes, err := hex.DecodeString("3bbb741e27fe3429a6de507cacefd4f3a4b478995eae56e546c869910c0a6b8d")
	if err != nil {
		t.Fatalf("decode hash: %v", err)
	}
	var hash [WitnessDigestBytes]byte
	copy(hash[:], hashBytes)

	got := DigestPublicInputs(hash)
	if got[0].Text(16) != "3bbb741e27fe3429a6de507cacefd4f3" {
		t.Fatalf("limb0 = %s, want %s", got[0].Text(16), "3bbb741e27fe3429a6de507cacefd4f3")
	}
	if got[1].Text(16) != "a4b478995eae56e546c869910c0a6b8d" {
		t.Fatalf("limb1 = %s, want %s", got[1].Text(16), "a4b478995eae56e546c869910c0a6b8d")
	}
}

func TestDigestPublicInputsRoundTrip(t *testing.T) {
	hash := [WitnessDigestBytes]byte{
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
		0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00,
	}

	limbs := DigestPublicInputs(hash)
	rebuilt := append(limbs[0].FillBytes(make([]byte, WitnessDigestFieldBytes)), limbs[1].FillBytes(make([]byte, WitnessDigestFieldBytes))...)
	if !bytes.Equal(rebuilt, hash[:]) {
		t.Fatalf("round-trip mismatch: got %x want %x", rebuilt, hash[:])
	}
}

// TestCommonPrefixFields_Round pins the committed round bytes: the sfixed64
// payload exactly as it sits on the wire (little-endian), all zero when the
// round field is absent — the value the circuit and SignatureVerifier hash.
func TestCommonPrefixFields_Round(t *testing.T) {
	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(i + 1)
	}
	cases := []struct {
		name  string
		round uint32
		want  [RoundLen]byte
	}{
		{"absent", 0, [RoundLen]byte{}},
		{"one", 1, [RoundLen]byte{1}},
		{"multi-byte", 0x01020304, [RoundLen]byte{4, 3, 2, 1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vote := crossCheckVote(5, tc.round, blockHash, 10)
			_, round, hash, err := CommonPrefixFields([]ValidatorSignature{{SignedBytes: vote, Active: true}})
			if err != nil {
				t.Fatalf("CommonPrefixFields: %v", err)
			}
			if round != tc.want {
				t.Errorf("round = %x, want %x", round, tc.want)
			}
			if hash != blockHash {
				t.Errorf("block hash = %x, want %x", hash, blockHash)
			}
		})
	}
}

// TestComputeWitnessHash_BindsRoundValue is the Go half of the cross-round fix:
// two batches that differ only in the vote round must commit to different
// digests. Before, only round presence was hashed, so rounds 1 and 2 collided.
func TestComputeWitnessHash_BindsRoundValue(t *testing.T) {
	var blockHash [32]byte
	hashFor := func(round uint32) [32]byte {
		sigs := make([]ValidatorSignature, 4)
		for i := range sigs {
			pub := make([]byte, 32)
			pub[0] = byte(i + 1)
			sigs[i] = ValidatorSignature{PublicKey: pub, SignedBytes: crossCheckVote(5, round, blockHash, 10), Active: true}
		}
		h, err := ComputeWitnessHash(sigs)
		if err != nil {
			t.Fatalf("ComputeWitnessHash(round=%d): %v", round, err)
		}
		return h
	}
	if hashFor(1) == hashFor(2) {
		t.Fatal("witness hash is identical for rounds 1 and 2")
	}
	if hashFor(0) == hashFor(1) {
		t.Fatal("witness hash is identical for rounds 0 and 1")
	}
}

// TestCommonPrefixFields_RejectsZeroRoundField: proto3 omits round == 0, so a
// round field carrying zero is non-canonical. The circuit derives presence from
// the value and would reject it; fail at derivation rather than prove time.
func TestCommonPrefixFields_RejectsZeroRoundField(t *testing.T) {
	var blockHash [32]byte
	vote := crossCheckVote(5, 1, blockHash, 10)
	// body starts at index 1 (1-byte varint); round payload at body offset 12.
	for k := 0; k < RoundLen; k++ {
		vote[1+RoundBodyOff+k] = 0
	}
	if _, _, _, err := CommonPrefixFields([]ValidatorSignature{{SignedBytes: vote, Active: true}}); err == nil {
		t.Fatal("expected an error for a round field carrying zero")
	}
}
