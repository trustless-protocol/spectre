package attestation

import (
	"crypto/ed25519"
	"encoding/hex"
	"strings"
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

func TestParseRunMode(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  RunMode
		valid bool
	}{
		{value: "unsafe", want: RunModeUnsafe, valid: true},
		{value: "safe", want: RunModeSafe, valid: true},
		{value: "finalized", want: RunModeFinalized, valid: true},
		{value: "", valid: false},
		{value: "confirmed", valid: false},
	} {
		t.Run(tc.value, func(t *testing.T) {
			got, err := ParseRunMode(tc.value)
			if tc.valid {
				if err != nil || got != tc.want {
					t.Fatalf("ParseRunMode(%q) = (%d, %v), want (%d, nil)", tc.value, got, err, tc.want)
				}
				return
			}
			if err == nil {
				t.Fatalf("ParseRunMode(%q) succeeded", tc.value)
			}
		})
	}
}

func TestNewSignerValidatesKeyMaterialAndEnvironment(t *testing.T) {
	const seed = "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	t.Setenv("ATTESTATION_TEST_SEED", seed)
	t.Setenv("ATTESTATION_TEST_MISSING", "")
	if _, err := NewSigner(10, "env:ATTESTATION_TEST_SEED"); err != nil {
		t.Fatalf("NewSigner from environment: %v", err)
	}

	for _, tc := range []struct {
		name    string
		chainID uint64
		key     string
		want    string
	}{
		{name: "zero chain ID", chainID: 0, key: seed, want: "chain ID"},
		{name: "invalid hex", chainID: 10, key: "not-hex", want: "decode"},
		{name: "wrong seed length", chainID: 10, key: "01", want: "32-byte"},
		{name: "missing environment key", chainID: 10, key: "env:ATTESTATION_TEST_MISSING", want: "unset"},
		{name: "empty environment variable name", chainID: 10, key: "env:", want: "empty"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewSigner(tc.chainID, tc.key)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("NewSigner(%d, %q) error = %v, want %q", tc.chainID, tc.key, err, tc.want)
			}
		})
	}
}

func TestSigningBytesRejectsInvalidIdentity(t *testing.T) {
	validHash := make([]byte, HashLength)
	for _, tc := range []struct {
		name      string
		chainID   uint64
		runMode   RunMode
		stateRoot []byte
		blockHash []byte
	}{
		{name: "zero chain ID", runMode: RunModeUnsafe, stateRoot: validHash, blockHash: validHash},
		{name: "invalid run mode", chainID: 10, runMode: 0, stateRoot: validHash, blockHash: validHash},
		{name: "short state root", chainID: 10, runMode: RunModeSafe, stateRoot: validHash[:HashLength-1], blockHash: validHash},
		{name: "short block hash", chainID: 10, runMode: RunModeFinalized, stateRoot: validHash, blockHash: validHash[:HashLength-1]},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := SigningBytes(tc.chainID, tc.runMode, 1, tc.stateRoot, tc.blockHash); err == nil {
				t.Fatal("SigningBytes succeeded for invalid identity")
			}
		})
	}
}
