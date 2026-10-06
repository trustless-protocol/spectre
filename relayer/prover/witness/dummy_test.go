package witness

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"testing"
)

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

// Layout: dummyMagic zero-padded to dummyMagicWidth || bucket(2 BE) || slot(2 BE).
func TestDummyMsgBytesLayout(t *testing.T) {
	cases := []struct{ bucket, slot uint16 }{{4, 0}, {4, 3}, {128, 127}, {0xFFFF, 0xFFFE}}
	for _, c := range cases {
		msg := DummyMsgBytes(c.bucket, c.slot)
		if len(msg) != DummyMsgLen {
			t.Fatalf("(%d,%d): len %d, want %d", c.bucket, c.slot, len(msg), DummyMsgLen)
		}
		if !bytes.HasPrefix(msg, []byte(dummyMagic)) {
			t.Errorf("(%d,%d): missing magic prefix: %x", c.bucket, c.slot, msg)
		}
		if pad := msg[len(dummyMagic):dummyMagicWidth]; !bytes.Equal(pad, make([]byte, len(pad))) {
			t.Errorf("(%d,%d): magic padding not zero: %x", c.bucket, c.slot, pad)
		}
		if got := binary.BigEndian.Uint16(msg[dummyMagicWidth:]); got != c.bucket {
			t.Errorf("(%d,%d): bucket field %d", c.bucket, c.slot, got)
		}
		if got := binary.BigEndian.Uint16(msg[dummyMagicWidth+2:]); got != c.slot {
			t.Errorf("(%d,%d): slot field %d", c.bucket, c.slot, got)
		}
	}
}

// Every padding slot of every supported bucket signs a distinct message.
func TestDummyMsgBytesDistinct(t *testing.T) {
	seen := map[string][2]int{}
	for _, n := range []int{4, 8, 16, 32, 64, 128} {
		for i := 0; i < n; i++ {
			key := string(DummyMsgBytes(uint16(n), uint16(i)))
			if prev, dup := seen[key]; dup {
				t.Fatalf("(%d,%d) collides with (%d,%d)", n, i, prev[0], prev[1])
			}
			seen[key] = [2]int{n, i}
		}
	}
}

// Every dummy slot must carry a valid Ed25519 signature over DummyMsgBytes(n, slot):
// the circuit verifies every slot, padding included. Pubkeys stay distinct across
// the whole bucket so the ECIP divisor never hits the doubling branch.
func TestPadWithDummiesSignaturesVerify(t *testing.T) {
	for _, n := range []int{4, 8, 16, 32, 64, 128} {
		padded, err := PadWithDummies(nil, n)
		if err != nil {
			t.Fatalf("n=%d: %v", n, err)
		}
		seenPub := map[string]bool{}
		seenR := map[string]bool{}
		for i, s := range padded {
			if s.Active || s.Power != 0 {
				t.Errorf("n=%d slot %d: padding must be inactive with zero power, got %+v", n, i, s)
			}
			if !bytes.Equal(s.SignedBytes, DummyMsgBytes(uint16(n), uint16(i))) {
				t.Errorf("n=%d slot %d: signed bytes %x", n, i, s.SignedBytes)
			}
			if !ed25519.Verify(s.PublicKey, s.SignedBytes, s.Signature) {
				t.Errorf("n=%d slot %d: dummy signature does not verify", n, i)
			}
			if seenPub[string(s.PublicKey)] || seenR[string(s.Signature[:32])] {
				t.Errorf("n=%d slot %d: duplicate pubkey or R point", n, i)
			}
			seenPub[string(s.PublicKey)] = true
			seenR[string(s.Signature[:32])] = true
		}
	}
}

func TestPadWithDummiesEdges(t *testing.T) {
	full := make([]ValidatorSignature, 4)
	for i := range full {
		full[i] = ValidatorSignature{Index: i, Active: true}
	}
	padded, err := PadWithDummies(full, 4)
	if err != nil {
		t.Fatalf("exact fit: %v", err)
	}
	for i, s := range padded {
		if !s.Active || s.Index != i {
			t.Errorf("exact fit slot %d altered: %+v", i, s)
		}
	}

	if padded, err := PadWithDummies(nil, 0); err != nil || len(padded) != 0 {
		t.Errorf("n=0: got %d slots, err %v", len(padded), err)
	}
	if _, err := PadWithDummies(nil, 0x10000); err == nil {
		t.Error("expected 16-bit slot encoding error for n=0x10000")
	}
}
