package prover

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
