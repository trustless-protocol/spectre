package witness

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
)

// TestDummyKeySlotFindsEveryPaddingKey checks DummyKeySlot against the keys
// PadWithDummies actually writes, for every slot of every supported bucket, so
// the lookup and the padding cannot drift apart.
func TestDummyKeySlotFindsEveryPaddingKey(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	one := []ValidatorSignature{{Signature: ed25519.Sign(priv, []byte("vote")), PublicKey: pub, SignedBytes: []byte("vote"), Active: true}}
	for _, n := range SupportedBuckets {
		padded, err := PadWithDummies(one, n)
		if err != nil {
			t.Fatalf("PadWithDummies(%d): %v", n, err)
		}
		for slot := 1; slot < n; slot++ {
			bucket, gotSlot, ok := DummyKeySlot(padded[slot].PublicKey)
			if !ok || bucket != n || gotSlot != slot {
				t.Fatalf("DummyKeySlot(padding key of bucket %d slot %d) = (%d, %d, %v)", n, slot, bucket, gotSlot, ok)
			}
		}
	}
	// Slot 0 is never padded (a batch has at least one signer), but its key is
	// just as public, so it is recognised too.
	if bucket, slot, ok := DummyKeySlot(dummyKey(16, 0).Public().(ed25519.PublicKey)); !ok || bucket != 16 || slot != 0 {
		t.Fatalf("DummyKeySlot(bucket 16 slot 0 key) = (%d, %d, %v)", bucket, slot, ok)
	}
}

func TestDummyKeySlotIgnoresOtherKeys(t *testing.T) {
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	if _, _, ok := DummyKeySlot(pub); ok {
		t.Fatal("a fresh key was taken for a padding key")
	}
	if _, _, ok := DummyKeySlot(nil); ok {
		t.Fatal("an empty key was taken for a padding key")
	}
}
