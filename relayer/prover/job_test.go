package prover

import "testing"

func TestNewProofID_DeterministicAndNamespaced(t *testing.T) {
	vk1 := [32]byte{1}
	vk2 := [32]byte{2}
	w1 := [32]byte{0xAA}
	w2 := [32]byte{0xBB}

	a, b := NewProofID(vk1, w1), NewProofID(vk1, w1)
	if a != b {
		t.Fatal("same inputs must give the same id")
	}
	if NewProofID(vk1, w1) == NewProofID(vk1, w2) {
		t.Fatal("different witness must give a different id")
	}
	// The vk namespace: same witness under a re-setup circuit is a NEW id, so
	// a stale proof from the old vk can never be returned for it.
	if NewProofID(vk1, w1) == NewProofID(vk2, w1) {
		t.Fatal("different vk must give a different id")
	}
}
