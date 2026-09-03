// Package attestation defines the stable signed L2-block attestation format.
package attestation

import (
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	attestorpb "attestor/types/attestor"
)

const (
	// Domain separates these signatures from every other Ed25519 use. It is
	// terminated because the fields that follow are fixed-width binary values.
	// This is a deployed wire-format identifier: do not rename it with a
	// Fast-IBC/Spectre branding change, because doing so invalidates existing
	// client profiles and their signatures.
	domain = "fast-ibc/l2-attestation/v2\x00"

	HashLength      = 32
	SignatureLength = ed25519.SignatureSize
)

// RunMode binds an attestation to the finality level at which the signer
// checked the block. Its byte values are a stable part of the signed format.
type RunMode uint8

const (
	RunModeUnsafe    RunMode = 1
	RunModeSafe      RunMode = 2
	RunModeFinalized RunMode = 3
)

// ParseRunMode parses the JSON/profile spelling of a signing finality level.
func ParseRunMode(value string) (RunMode, error) {
	switch value {
	case "unsafe":
		return RunModeUnsafe, nil
	case "safe":
		return RunModeSafe, nil
	case "finalized":
		return RunModeFinalized, nil
	default:
		return 0, fmt.Errorf("attestation run mode must be unsafe, safe, or finalized, got %q", value)
	}
}

func (m RunMode) valid() bool {
	return m == RunModeUnsafe || m == RunModeSafe || m == RunModeFinalized
}

// RunModeFromProto explicitly maps the transport enum to the byte included in
// the signature. Keep this mapping rather than casting: protobuf enum numbers
// are not part of the attestation signing format.
func RunModeFromProto(mode attestorpb.RunMode) (RunMode, error) {
	switch mode {
	case attestorpb.RunMode_RUN_MODE_UNSAFE:
		return RunModeUnsafe, nil
	case attestorpb.RunMode_RUN_MODE_SAFE:
		return RunModeSafe, nil
	case attestorpb.RunMode_RUN_MODE_FINALIZED:
		return RunModeFinalized, nil
	default:
		return 0, fmt.Errorf("unsupported protobuf attestation run mode %d", mode)
	}
}

// Signer owns the key used to attest canonical L2 block identities for one L2.
type Signer struct {
	chainID uint64
	key     ed25519.PrivateKey
}

// Configured reports whether the signer has usable key material. It lets
// transport adapters fail closed instead of accidentally returning a positive
// but unsigned verdict when composition omitted a per-chain signer.
func (s Signer) Configured() bool {
	return s.chainID != 0 && len(s.key) == ed25519.PrivateKeySize
}

// NewSigner parses a 32-byte Ed25519 seed encoded as 0x-prefixed or bare hex.
func NewSigner(chainID uint64, privateKeyHex string) (Signer, error) {
	if chainID == 0 {
		return Signer{}, fmt.Errorf("L2 chain ID must be non-zero")
	}
	privateKeyHex, err := resolvePrivateKey(privateKeyHex)
	if err != nil {
		return Signer{}, err
	}
	raw, err := hex.DecodeString(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return Signer{}, fmt.Errorf("decode attestation signing key: %w", err)
	}
	if len(raw) != ed25519.SeedSize {
		return Signer{}, fmt.Errorf("attestation signing key must be a %d-byte Ed25519 seed, got %d", ed25519.SeedSize, len(raw))
	}
	return Signer{chainID: chainID, key: ed25519.NewKeyFromSeed(raw)}, nil
}

func resolvePrivateKey(value string) (string, error) {
	if !strings.HasPrefix(value, "env:") {
		return value, nil
	}
	name := strings.TrimPrefix(value, "env:")
	if name == "" {
		return "", fmt.Errorf("attestation signing-key environment variable is empty")
	}
	key, ok := os.LookupEnv(name)
	if !ok || key == "" {
		return "", fmt.Errorf("attestation signing-key environment variable %q is unset", name)
	}
	return key, nil
}

// PublicKey returns the configured signer's 32-byte Ed25519 public key.
func (s Signer) PublicKey() []byte {
	return append([]byte(nil), s.key.Public().(ed25519.PublicKey)...)
}

// Sign signs one canonical L2 block identity at a specific finality level.
func (s Signer) Sign(runMode RunMode, blockNumber uint64, stateRoot, blockHash []byte) ([]byte, error) {
	message, err := SigningBytes(s.chainID, runMode, blockNumber, stateRoot, blockHash)
	if err != nil {
		return nil, err
	}
	return ed25519.Sign(s.key, message), nil
}

// SigningBytes is the exact, cross-language byte encoding signed by the attestor:
// domain || l2_chain_id (uint64 big endian) || run_mode (uint8) ||
// l2_block_number (uint64 big endian) || state_root (32 bytes) || block_hash
// (32 bytes).
func SigningBytes(chainID uint64, runMode RunMode, blockNumber uint64, stateRoot, blockHash []byte) ([]byte, error) {
	if chainID == 0 {
		return nil, fmt.Errorf("L2 chain ID must be non-zero")
	}
	if !runMode.valid() {
		return nil, fmt.Errorf("invalid attestation run mode %d", runMode)
	}
	if len(stateRoot) != HashLength || len(blockHash) != HashLength {
		return nil, fmt.Errorf("state root and block hash must both be %d bytes", HashLength)
	}
	message := make([]byte, len(domain)+8+1+8+HashLength+HashLength)
	offset := copy(message, domain)
	binary.BigEndian.PutUint64(message[offset:], chainID)
	offset += 8
	message[offset] = byte(runMode)
	offset++
	binary.BigEndian.PutUint64(message[offset:], blockNumber)
	offset += 8
	offset += copy(message[offset:], stateRoot)
	copy(message[offset:], blockHash)
	return message, nil
}
