package prover

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"strings"
	"testing"

	"filippo.io/edwards25519"
)

// TestValidateSignatures_RejectsWhatTheCircuitCannotProve is the #579 item 3
// check at the prover's door: a signature crypto/ed25519 accepts but the circuit
// rejects fails at submit time as ErrInvalidInput, before a prove is spent on it.
// The identity public key is one: [k]A = O, so R = [s]B, S = s verifies for any
// message.
func TestValidateSignatures_RejectsWhatTheCircuitCannotProve(t *testing.T) {
	sigs := MakeTestSigs(t, 2)

	var buf [64]byte
	if _, err := rand.Read(buf[:]); err != nil {
		t.Fatalf("rand: %v", err)
	}
	s, err := edwards25519.NewScalar().SetUniformBytes(buf[:])
	if err != nil {
		t.Fatalf("scalar: %v", err)
	}
	identity := edwards25519.NewIdentityPoint().Bytes()
	sigs[1].PublicKey = identity
	sigs[1].Signature = append(new(edwards25519.Point).ScalarBaseMult(s).Bytes(), s.Bytes()...)
	if !ed25519.Verify(identity, sigs[1].SignedBytes, sigs[1].Signature) {
		t.Fatal("crypto/ed25519 rejects the identity-key signature; the test would prove nothing")
	}

	err = validateSignatures(sigs)
	var invalid *ErrInvalidInput
	if !errors.As(err, &invalid) || invalid.Slot != 1 || !strings.Contains(invalid.Reason, "the circuit cannot prove") {
		t.Fatalf("validateSignatures = %v, want ErrInvalidInput at slot 1 naming the circuit", err)
	}

	if err := validateSignatures(MakeTestSigs(t, 2)); err != nil {
		t.Fatalf("honest signatures rejected: %v", err)
	}
}
