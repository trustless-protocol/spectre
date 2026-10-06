package witness

import "testing"

func TestPadWithDummies(t *testing.T) {
	real := []ValidatorSignature{{
		Signature:   append(make([]byte, 63), 0xAB),
		PublicKey:   append(make([]byte, 31), 0xCD),
		SignedBytes: []byte("vote"),
		Index:       7,
		Active:      true,
	}}

	padded, err := PadWithDummies(real, 4)
	if err != nil {
		t.Fatalf("pad: %v", err)
	}
	if len(padded) != 4 {
		t.Fatalf("expected 4 slots, got %d", len(padded))
	}
	if !padded[0].Active || padded[0].Index != 7 {
		t.Errorf("slot 0 should be the real signer, got %+v", padded[0])
	}
	seen := map[string]bool{string(padded[0].PublicKey): true}
	for i := 1; i < 4; i++ {
		if padded[i].Active {
			t.Errorf("slot %d should be inactive padding", i)
		}
		// Distinctness keeps the ECIP divisor off the doubling branch.
		if key := string(padded[i].PublicKey); seen[key] {
			t.Errorf("duplicate pubkey at slot %d", i)
		} else {
			seen[key] = true
		}
		if len(padded[i].SignedBytes) != DummyMsgLen {
			t.Errorf("slot %d dummy message length %d, want %d", i, len(padded[i].SignedBytes), DummyMsgLen)
		}
	}

	// Determinism: same (n, slot) must reproduce identical padding across runs.
	again, err := PadWithDummies(real, 4)
	if err != nil {
		t.Fatalf("pad again: %v", err)
	}
	for i := 1; i < 4; i++ {
		if string(again[i].PublicKey) != string(padded[i].PublicKey) ||
			string(again[i].Signature) != string(padded[i].Signature) {
			t.Errorf("slot %d padding not deterministic", i)
		}
	}

	if _, err := PadWithDummies(make([]ValidatorSignature, 5), 4); err == nil {
		t.Fatal("expected too-small bucket error")
	}
}

// DummyMsgLen is the circuit's SHA-512 MinMsgLen, so it is baked into the R1CS
// and every published verifying key. Changing it needs re-setup + redeploy of
// every Groth16Verifier_N{N} + setBucket — this pins it so that is never silent.
func TestDummyMsgLenPinnedToCircuit(t *testing.T) {
	if DummyMsgLen != 18 {
		t.Fatalf("DummyMsgLen = %d, want 18: changing it changes the circuit R1CS", DummyMsgLen)
	}
}
