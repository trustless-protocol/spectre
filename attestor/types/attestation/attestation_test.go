package attestation

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	attestorpb "attestor/types/attestor"
)

// This pins the protobuf numbers as well as the explicit conversion. A change
// to either representation is a signing-wire-format change and must be made
// deliberately across the Go and CosmWasm implementations.
func TestProtoRunModesMatchSigningBytes(t *testing.T) {
	tests := []struct {
		proto attestorpb.RunMode
		want  RunMode
	}{
		{proto: attestorpb.RunMode_RUN_MODE_UNSAFE, want: RunModeUnsafe},
		{proto: attestorpb.RunMode_RUN_MODE_SAFE, want: RunModeSafe},
		{proto: attestorpb.RunMode_RUN_MODE_FINALIZED, want: RunModeFinalized},
	}
	for _, test := range tests {
		got, err := RunModeFromProto(test.proto)
		if err != nil {
			t.Fatalf("RunModeFromProto(%v): %v", test.proto, err)
		}
		if got != test.want {
			t.Errorf("RunModeFromProto(%v) = %d, want %d", test.proto, got, test.want)
		}
		if RunMode(test.proto) != test.want {
			t.Errorf("protobuf %v number = %d, want signing byte %d", test.proto, test.proto, test.want)
		}
	}
}

func TestSignerBindsEveryBlockIdentityField(t *testing.T) {
	signer, err := NewSigner(8453, "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20")
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	stateRoot := make([]byte, HashLength)
	blockHash := make([]byte, HashLength)
	stateRoot[31] = 1
	blockHash[31] = 2
	signature, err := signer.Sign(RunModeSafe, 123, stateRoot, blockHash)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	const expectedPublicKey = "79b5562e8fe654f94078b112e8a98ba7901f853ae695bed7e0e3910bad049664"
	const expectedSignature = "22cc6cdee3b8e865946ade58d2ed25efdc5cff7a918970c0bc03aee03b32ebbd6a2d253606f37eecd0b6aea651f898c44abeb8b49349e2c4c8b47bc43e1fe609"
	if got := hex.EncodeToString(signer.PublicKey()); got != expectedPublicKey {
		t.Fatalf("public key = %s, want %s", got, expectedPublicKey)
	}
	if got := hex.EncodeToString(signature); got != expectedSignature {
		t.Fatalf("signature = %s, want %s", got, expectedSignature)
	}
	message, err := SigningBytes(8453, RunModeSafe, 123, stateRoot, blockHash)
	if err != nil {
		t.Fatalf("SigningBytes: %v", err)
	}
	if !ed25519.Verify(signer.PublicKey(), message, signature) {
		t.Fatal("valid block identity signature did not verify")
	}
	for _, bad := range [][]byte{
		mustSigningBytes(t, 8454, RunModeSafe, 123, stateRoot, blockHash),
		mustSigningBytes(t, 8453, RunModeUnsafe, 123, stateRoot, blockHash),
		mustSigningBytes(t, 8453, RunModeSafe, 124, stateRoot, blockHash),
		mustSigningBytes(t, 8453, RunModeSafe, 123, append([]byte(nil), blockHash...), blockHash),
	} {
		if ed25519.Verify(signer.PublicKey(), bad, signature) {
			t.Fatal("signature verified for a different L2 block identity")
		}
	}
}

func mustSigningBytes(t *testing.T, chainID uint64, runMode RunMode, blockNumber uint64, stateRoot, blockHash []byte) []byte {
	t.Helper()
	message, err := SigningBytes(chainID, runMode, blockNumber, stateRoot, blockHash)
	if err != nil {
		t.Fatal(err)
	}
	return message
}
