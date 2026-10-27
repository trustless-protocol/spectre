package prover

import (
	"crypto/ed25519"
	"fmt"

	"relayer/prover/witness"
)

// ErrInvalidInput reports a request rejected by validation at submit time.
// The input is wrong, not the system: submitting the same input again will
// fail the same way.
type ErrInvalidInput struct {
	Reason string
	Slot   int // offending slot, -1 when the problem is not slot-specific
}

func (e *ErrInvalidInput) Error() string {
	if e.Slot >= 0 {
		return fmt.Sprintf("invalid input at slot %d: %s", e.Slot, e.Reason)
	}
	return "invalid input: " + e.Reason
}

// validateSignatures runs every cheap check at submit time, so a caller never
// waits out a prove to learn its input was broken. Anything that passes here
// can only fail later for expensive reasons (prove / local verify), which the
// error catalog classifies as internal.
func validateSignatures(sigs []witness.ValidatorSignature) error {
	if len(sigs) == 0 {
		return &ErrInvalidInput{Reason: "no signatures provided", Slot: -1}
	}
	active := 0
	for i, s := range sigs {
		if len(s.Signature) != 64 {
			return &ErrInvalidInput{Reason: "signature must be 64 bytes", Slot: i}
		}
		if len(s.PublicKey) != 32 {
			return &ErrInvalidInput{Reason: "public key must be 32 bytes", Slot: i}
		}
		if len(s.SignedBytes) == 0 {
			return &ErrInvalidInput{Reason: "signed bytes empty", Slot: i}
		}
		if len(s.SignedBytes) > witness.MaxMsgLen {
			return &ErrInvalidInput{Reason: "signed bytes exceed witness.MaxMsgLen", Slot: i}
		}
		if !s.Active {
			continue
		}
		active++
		// A full stdlib Ed25519 verify (~50µs/slot) subsumes point
		// decompression and canonicity checks, and catches a bad signature
		// here instead of as a cryptic in-circuit failure minutes later.
		if !ed25519.Verify(ed25519.PublicKey(s.PublicKey), s.SignedBytes, s.Signature) {
			return &ErrInvalidInput{Reason: "signature does not verify", Slot: i}
		}
	}
	if active == 0 {
		return &ErrInvalidInput{Reason: "no active signer", Slot: -1}
	}
	return nil
}
