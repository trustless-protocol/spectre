package l2rollup

import (
	"crypto/ed25519"
	"fmt"

	"attestor/types/attestation"
	attestorpb "attestor/types/attestor"
)

// AttestationVerifier verifies the attestor key pinned in a Cosmos L2-client
// profile. The relayer checks a response before forwarding it; the wasm client
// independently repeats the same check before accepting the update on-chain.
type AttestationVerifier struct {
	chainID   uint64
	runMode   attestation.RunMode
	publicKey ed25519.PublicKey
}

// NewAttestationVerifier constructs a verifier for one L2 and validates the
// fixed-size profile key before any relay loop starts.
func NewAttestationVerifier(chainID uint64, attestationHead string, publicKey []byte) (AttestationVerifier, error) {
	if chainID == 0 {
		return AttestationVerifier{}, fmt.Errorf("attestation L2 chain ID must be non-zero")
	}
	runMode, err := attestation.ParseRunMode(attestationHead)
	if err != nil {
		return AttestationVerifier{}, fmt.Errorf("invalid attestation head: %w", err)
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return AttestationVerifier{}, fmt.Errorf("attestor public key must be %d bytes, got %d", ed25519.PublicKeySize, len(publicKey))
	}
	return AttestationVerifier{
		chainID:   chainID,
		runMode:   runMode,
		publicKey: append(ed25519.PublicKey(nil), publicKey...),
	}, nil
}

// MatchesRunMode reports whether the relayer's requested replica head is the
// finality level pinned into the Cosmos client's immutable profile.
func (v AttestationVerifier) MatchesRunMode(mode attestorpb.RunMode) bool {
	runMode, err := attestation.RunModeFromProto(mode)
	return err == nil && runMode == v.runMode
}

// Verify confirms the signature over the exact L2 block identity the relayer
// fetched. This is the Go counterpart of the l2-client CosmWasm verification.
func (v AttestationVerifier) Verify(blockNumber uint64, stateRoot, blockHash, signature []byte) error {
	if len(v.publicKey) != ed25519.PublicKeySize {
		return fmt.Errorf("attestor public key is not configured")
	}
	if len(signature) != attestation.SignatureLength {
		return fmt.Errorf("attestor signature must be %d bytes, got %d", attestation.SignatureLength, len(signature))
	}
	message, err := attestation.SigningBytes(v.chainID, v.runMode, blockNumber, stateRoot, blockHash)
	if err != nil {
		return fmt.Errorf("encode signed L2 block identity: %w", err)
	}
	if !ed25519.Verify(v.publicKey, message, signature) {
		return fmt.Errorf("attestor signature does not match the configured L2 block identity")
	}
	return nil
}
