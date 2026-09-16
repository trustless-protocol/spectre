package relay

import (
	"context"
	"errors"
	"testing"
	"time"

	"relayer/chain"
)

// Status is the read side of a running path, so the two things that make it
// usable during an incident are that it reports the oldest wait first and that
// it answers at all while the path is busy.
func TestModuleStatus(t *testing.T) {
	t.Run("reports the waiting set oldest first", func(t *testing.T) {
		m := NewModule("test", "client-0", &mockSource{latest: 100, relayable: 90}, &mockDest{}, &mockBuilder{})

		// Distinct ages, inserted youngest first so a passing test cannot be map
		// iteration order.
		base := time.Now()
		m.waits.now = func() time.Time { return base }
		m.waits.observe(chain.Event{Type: chain.AckPacket, Sequence: 2, Height: 8})
		m.waits.now = func() time.Time { return base.Add(-2 * time.Minute) }
		m.waits.observe(chain.Event{Type: chain.SendPacket, Sequence: 1, Height: 7})
		m.waits.now = func() time.Time { return base.Add(time.Minute) }

		got := m.Status(context.Background())
		if len(got.Waiting) != 2 {
			t.Fatalf("waiting = %+v, want both packets", got.Waiting)
		}
		if got.Waiting[0].Sequence != 1 {
			t.Fatalf("first entry is seq=%d; the oldest wait must lead, because it is the one most likely to be stuck",
				got.Waiting[0].Sequence)
		}
		if got.Waiting[0].Age <= got.Waiting[1].Age {
			t.Fatalf("ages = %s, %s: not sorted oldest first", got.Waiting[0].Age, got.Waiting[1].Age)
		}
		if got.Waiting[0].AgeText == "" {
			t.Error("no human-readable age; the raw nanoseconds are for machines, not the operator reading this")
		}
	})

	t.Run("reports the heights the operator is asking for", func(t *testing.T) {
		m := NewModule("test", "client-0", &mockSource{latest: 100, relayable: 90}, &mockDest{}, &mockBuilder{})
		m.lastHeight = 88

		got := m.Status(context.Background())
		if got.TrustedHeight != 88 {
			t.Errorf("trusted height = %d, want 88", got.TrustedHeight)
		}
		if got.RelayableHeight != 90 {
			t.Errorf("relayable height = %d, want 90: this is the number that answers 'why has nothing moved'", got.RelayableHeight)
		}
		if got.Path != "test" || got.ClientID != "client-0" {
			t.Errorf("identity = %s/%s, want the module's", got.Path, got.ClientID)
		}
	})

	// The one RPC in the snapshot must not be able to fail the snapshot: an
	// unreachable source is exactly when someone runs this.
	t.Run("a failed probe is reported, not fatal", func(t *testing.T) {
		src := &mockSource{latest: 100, relayableErr: errors.New("source rpc down")}
		m := NewModule("test", "client-0", src, &mockDest{}, &mockBuilder{})
		m.lastHeight = 88

		got := m.Status(context.Background())
		if got.RelayableHeightErr == "" {
			t.Fatal("a failed probe produced no error text; the field would read as height 0")
		}
		if got.TrustedHeight != 88 {
			t.Errorf("trusted height = %d, want the in-memory value to survive a failed probe", got.TrustedHeight)
		}
	})

	// updateClientTo holds m.mu across the builder, which on the Cosmos->EVM path
	// is across proof generation. A status command that blocks for the length of a
	// Groth16 proof is useless at the moment someone runs it, so the busy lock is
	// reported as an answer instead of waited on.
	t.Run("answers while a client update holds the lock", func(t *testing.T) {
		m := NewModule("test", "client-0", &mockSource{latest: 100, relayable: 90}, &mockDest{}, &mockBuilder{})
		m.mu.Lock()
		defer m.mu.Unlock()

		done := make(chan Status, 1)
		go func() { done <- m.Status(context.Background()) }()
		select {
		case got := <-done:
			if !got.UpdateInFlight {
				t.Fatal("the snapshot did not say an update was in flight, so a stale trusted height reads as current")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Status blocked on the module lock: it would hang for the length of a proof")
		}
	})
}
