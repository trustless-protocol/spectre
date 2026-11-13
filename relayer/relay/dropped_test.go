package relay

import (
	"context"
	"strings"
	"testing"
	"time"

	"relayer/chain"
	"relayer/services"
)

// newLedgerModule wires a module to an in-memory drop ledger and a clock the
// test moves.
func newLedgerModule(t *testing.T, src chain.Source, dst chain.Destination, builder chain.Builder) (*Module, *time.Time) {
	t.Helper()
	ledger, err := services.NewDropLedger("")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Unix(1_700_000_000, 0)
	m := NewModule("cosmos->eth", "client-0", src, dst, builder, WithDropLedger(ledger))
	m.clock = func() time.Time { return now }
	return m, &now
}

// The audit's engine PoC: one acknowledgement whose callback always reverts, and
// five flush passes with the source head moving. Before the ledger every pass
// advanced the client and re-submitted the ack (5 updates, 5 reverting
// submissions). A drop now lasts until its backoff expires (#575 M2).
func TestDroppedPacketIsNotResubmittedEveryFlushPass(t *testing.T) {
	ack := chain.Event{Type: chain.AckPacket, Sequence: 7, Height: 40, Raw: []byte("poison-ack")}
	src := &listerSource{mockSource: mockSource{latest: 100, relayable: 100}, candidate: []chain.Event{ack}}
	dst := &mockDest{poison: map[string]bool{"poison-ack": true}}
	builder := &mockBuilder{}
	m, now := newLedgerModule(t, src, dst, builder)

	out := captureLog(func() {
		for pass := 0; pass < 5; pass++ {
			src.relayable = 100 + uint64(pass) // the head moved since the last pass
			m.flushOnce(context.Background(), src)
		}
	})

	if dst.relayCalls != 1 {
		t.Fatalf("submitted the poison ack %d times across 5 flush passes, want 1", dst.relayCalls)
	}
	if builder.built != 1 {
		t.Fatalf("built %d client updates across 5 flush passes, want 1: a held packet must not advance the client", builder.built)
	}
	if !strings.Contains(out, "DROP packet (ack seq=7)") || !strings.Contains(out, "held until") {
		t.Fatalf("the drop must say when the packet is retried:\n%s", out)
	}

	// The hold is a backoff, not a verdict: a dropped ack has no other exit.
	*now = now.Add(10 * time.Minute)
	src.relayable = 110
	captureLog(func() { m.flushOnce(context.Background(), src) })
	if dst.relayCalls != 2 {
		t.Fatalf("after the backoff the ack was submitted %d times in total, want 2", dst.relayCalls)
	}

	// Each further drop lengthens the hold.
	*now = now.Add(10 * time.Minute)
	captureLog(func() { m.flushOnce(context.Background(), src) })
	if dst.relayCalls != 2 {
		t.Fatalf("the second hold is 20 min; 10 min in, the ack was submitted again (%d calls)", dst.relayCalls)
	}
}

// The live path and the startup rescan reach handleBatch without the flush, so
// the hold has to apply there too, before the client update.
func TestHeldPacketCostsNoClientUpdateInABatch(t *testing.T) {
	ack := chain.Event{Type: chain.AckPacket, Sequence: 7, Height: 40, Raw: []byte("poison-ack")}
	src := &mockSource{latest: 100, relayable: 100}
	dst := &mockDest{poison: map[string]bool{"poison-ack": true}}
	builder := &mockBuilder{}
	m, _ := newLedgerModule(t, src, dst, builder)

	captureLog(func() { m.handleBatch(context.Background(), []chain.Event{ack}) })
	if builder.built != 1 || dst.relayCalls != 1 {
		t.Fatalf("first attempt: %d update(s), %d submission(s), want 1 and 1", builder.built, dst.relayCalls)
	}

	src.relayable = 120
	var requeue []int
	captureLog(func() { requeue = m.handleBatch(context.Background(), []chain.Event{ack}) })
	if builder.built != 1 || dst.relayCalls != 1 {
		t.Fatalf("held packet: %d update(s), %d submission(s) in total, want still 1 and 1", builder.built, dst.relayCalls)
	}
	// A held packet left in a source's pending buffer would pin its recovery cursor.
	if len(requeue) != 0 {
		t.Fatalf("held packet was re-queued (%v); the flush backstop is what offers it again", requeue)
	}
}

func TestHeldPacketDoesNotBlockItsSiblings(t *testing.T) {
	held := chain.Event{Type: chain.AckPacket, Sequence: 7, Height: 40, Raw: []byte("poison-ack")}
	fresh := chain.Event{Type: chain.AckPacket, Sequence: 8, Height: 41, Raw: []byte("good-ack")}
	src := &mockSource{latest: 100, relayable: 100}
	dst := &mockDest{poison: map[string]bool{"poison-ack": true}}
	m, _ := newLedgerModule(t, src, dst, &mockBuilder{})

	captureLog(func() { m.handleBatch(context.Background(), []chain.Event{held}) })
	captureLog(func() { m.handleBatch(context.Background(), []chain.Event{held, fresh}) })

	if len(dst.relayed) != 1 || string(dst.relayed[0].Packet) != "good-ack" {
		t.Fatalf("relayed %+v, want only the fresh ack", dst.relayed)
	}
}

// A held send must still be tracked: the timeout scanner reads only the tracker,
// and a send it does not know about is never refunded.
func TestHeldSendIsStillTrackedForItsTimeout(t *testing.T) {
	send := chain.Event{Type: chain.SendPacket, Sequence: 3, Height: 40, Raw: []byte("expired-send")}
	src := &mockSource{latest: 100, relayable: 100, permanentMembershipOn: map[string]bool{"expired-send": true}}
	dst := &mockDest{}
	m, _ := newLedgerModule(t, src, dst, &mockBuilder{})
	tracked := 0
	m.track = func([]byte, uint64) bool { tracked++; return true }

	captureLog(func() { m.handleBatch(context.Background(), []chain.Event{send}) })
	captureLog(func() { m.handleBatch(context.Background(), []chain.Event{send}) })

	if src.membershipCalls != 1 {
		t.Fatalf("built %d proofs for a send whose proof failed permanently, want 1", src.membershipCalls)
	}
	if tracked != 2 {
		t.Fatalf("tracked the send %d times, want 2: the hold must not skip the timeout record", tracked)
	}
}

func TestDeliveryAndSettlementClearTheLedger(t *testing.T) {
	ack := chain.Event{Type: chain.AckPacket, Sequence: 7, Height: 40, Raw: []byte("ack")}

	t.Run("a retry that lands", func(t *testing.T) {
		src := &mockSource{latest: 100, relayable: 100}
		dst := &mockDest{poison: map[string]bool{"ack": true}}
		m, now := newLedgerModule(t, src, dst, &mockBuilder{})

		captureLog(func() { m.handleBatch(context.Background(), []chain.Event{ack}) })
		if got := m.Status(context.Background()).Dropped; len(got) != 1 || got[0].Sequence != 7 || got[0].Type != "ack" || got[0].Drops != 1 {
			t.Fatalf("status.Dropped = %+v, want the stranded ack", got)
		}

		dst.poison = nil // whatever made the callback revert was fixed
		*now = now.Add(10 * time.Minute)
		captureLog(func() { m.handleBatch(context.Background(), []chain.Event{ack}) })
		if len(dst.relayed) != 1 {
			t.Fatalf("relayed %d packets after the backoff, want 1", len(dst.relayed))
		}
		if got := m.Status(context.Background()).Dropped; len(got) != 0 {
			t.Fatalf("status.Dropped = %+v after delivery, want none", got)
		}
	})

	// Another relayer delivered it. The source stops listing a settled packet, so
	// nothing would offer it to the batch again; the flush asks on its own.
	t.Run("settled without this process", func(t *testing.T) {
		send := chain.Event{Type: chain.SendPacket, Sequence: 9, Height: 40, Raw: []byte("send")}
		src := &listerSource{mockSource: mockSource{latest: 100, relayable: 100}}
		dst := &mockDest{poison: map[string]bool{"send": true}}
		m, now := newLedgerModule(t, src, dst, &mockBuilder{})

		captureLog(func() { m.handleBatch(context.Background(), []chain.Event{send}) })
		if got := m.Status(context.Background()).Dropped; len(got) != 1 {
			t.Fatalf("status.Dropped = %+v, want the dropped send", got)
		}

		dst.hasReceipt = true
		captureLog(func() { m.flushOnce(context.Background(), src) })
		if got := m.Status(context.Background()).Dropped; len(got) != 1 {
			t.Fatalf("inside the backoff the entry is not asked about; status.Dropped = %+v", got)
		}
		*now = now.Add(10 * time.Minute)
		captureLog(func() { m.flushOnce(context.Background(), src) })
		if got := m.Status(context.Background()).Dropped; len(got) != 0 {
			t.Fatalf("status.Dropped = %+v, want none: the destination holds a receipt", got)
		}
	})
}

// Without a ledger the module behaves as it did: this is what mocks and the
// paths that wire none get.
func TestNoLedgerKeepsTheOldBehaviour(t *testing.T) {
	ack := chain.Event{Type: chain.AckPacket, Sequence: 7, Height: 40, Raw: []byte("poison-ack")}
	src := &mockSource{latest: 100, relayable: 100}
	dst := &mockDest{poison: map[string]bool{"poison-ack": true}}
	m := NewModule("cosmos->eth", "client-0", src, dst, &mockBuilder{})

	captureLog(func() { m.handleBatch(context.Background(), []chain.Event{ack}) })
	captureLog(func() { m.handleBatch(context.Background(), []chain.Event{ack}) })
	if dst.relayCalls != 2 {
		t.Fatalf("submitted %d times, want 2", dst.relayCalls)
	}
	if got := m.Status(context.Background()).Dropped; got != nil {
		t.Fatalf("status.Dropped = %+v, want nil", got)
	}
}
