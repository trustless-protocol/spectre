package services

import (
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func trackedPacket(seq uint64) channeltypesv2.Packet {
	return channeltypesv2.Packet{SourceClient: "client-0", Sequence: seq}
}

// TestNextTimeoutBackoffGrowsAndCaps pins the schedule: 1m, 2m, 4m … capped at
// 30m. Without it a deterministically-failing timeout was re-attempted on the raw
// 30s scan cadence, and every attempt costs a Groth16 proof and a transaction.
func TestNextTimeoutBackoffGrowsAndCaps(t *testing.T) {
	want := []time.Duration{
		1 * time.Minute, 2 * time.Minute, 4 * time.Minute, 8 * time.Minute,
		16 * time.Minute, 30 * time.Minute, 30 * time.Minute,
	}
	for i, w := range want {
		if got := nextTimeoutBackoff(i + 1); got != w {
			t.Fatalf("attempt %d: backoff = %s, want %s", i+1, got, w)
		}
	}
	// A very large attempt must still clamp rather than overflow the shift.
	if got := nextTimeoutBackoff(64); got != timeoutRetryBackoffMax {
		t.Fatalf("attempt 64: backoff = %s, want the %s cap (shift overflow guard)", got, timeoutRetryBackoffMax)
	}
}

// A failed timeout submission must hold the packet out of the next scans, so the
// scanner stops re-proving it every 30 seconds.
func TestRecordTimeoutFailureHoldsPacketOutOfScan(t *testing.T) {
	tr := NewPendingPacketTracker()
	tr.Add(trackedPacket(1), 100)

	now := time.Now()
	if got := len(tr.GetDue(now)); got != 1 {
		t.Fatalf("a fresh packet must be due, got %d", got)
	}

	if dead := tr.recordTimeoutFailureForTest("client-0", 1, now); dead {
		t.Fatal("one failure must not dead-letter")
	}
	if got := len(tr.GetDue(now)); got != 0 {
		t.Fatalf("packet under backoff must not be re-scanned, got %d due", got)
	}
	// It comes back once the backoff elapses.
	if got := len(tr.GetDue(now.Add(2 * time.Minute))); got != 1 {
		t.Fatalf("packet must be due again after its backoff, got %d", got)
	}
}

// After maxTimeoutAttempts the packet stops being retried and is kept for an
// operator. It must leave the scan set — otherwise the cap achieves nothing.
func TestRecordTimeoutFailureDeadLettersAtCap(t *testing.T) {
	tr := NewPendingPacketTracker()
	tr.Add(trackedPacket(1), 100)

	now := time.Now()
	for i := 1; i <= maxTimeoutAttempts; i++ {
		if dead := tr.recordTimeoutFailureForTest("client-0", 1, now); dead {
			t.Fatalf("attempt %d dead-lettered early (cap is %d)", i, maxTimeoutAttempts)
		}
	}
	if dead := tr.recordTimeoutFailureForTest("client-0", 1, now); !dead {
		t.Fatalf("attempt %d must dead-letter", maxTimeoutAttempts+1)
	}
	if tr.Len() != 0 {
		t.Fatalf("a dead-lettered packet must leave the scan set, %d still tracked", tr.Len())
	}
	if got := tr.DeadLetteredTimeouts(); got != 1 {
		t.Fatalf("DeadLetteredTimeouts = %d, want 1 — escrowed funds must stay visible", got)
	}
}

// Recording against a packet that is no longer tracked (it succeeded, or was
// purged) must be a no-op rather than resurrecting it.
func TestRecordTimeoutFailureIgnoresUnknownPacket(t *testing.T) {
	tr := NewPendingPacketTracker()
	if dead := tr.recordTimeoutFailureForTest("client-0", 99, time.Now()); dead {
		t.Fatal("unknown packet must not be reported as dead-lettered")
	}
	if tr.Len() != 0 || tr.DeadLetteredTimeouts() != 0 {
		t.Fatal("unknown packet must not be added to any set")
	}
}

// TestDeadLettersRemainVisibleUntilRetryRetentionExpires ensures a read-only
// inspection does not clear the gauge while the timeout remains suppressed.
func TestDeadLettersRemainVisibleUntilRetryRetentionExpires(t *testing.T) {
	tracker, packet := trackerWith(t, 1)
	deadAt := time.Now()
	for i := 0; i <= maxTimeoutAttempts; i++ {
		tracker.recordTimeoutFailureForTest(packet.SourceClient, packet.Sequence, deadAt)
	}
	if got := tracker.DeadLetteredTimeouts(); got != 1 {
		t.Fatalf("expected one dead-lettered timeout, got %d", got)
	}
	packets := tracker.DeadLetteredTimeoutPackets()
	if len(packets) != 1 {
		t.Fatalf("dead-lettered packets=%d, want 1", len(packets))
	}
	if tracker.deadLettered[0].DeadLetteredAt.IsZero() {
		t.Fatal("DeadLetteredAt must be stamped so the age is observable")
	}
	if got := tracker.DeadLetteredTimeouts(); got != 1 {
		t.Fatalf("inspection must not clear the timeout dead-letter gauge, got %d", got)
	}
}

// OldestDeadLetteredTimeout reports zero with nothing dead-lettered and a real
// age once something is — a count alone cannot distinguish a failure nobody has
// looked at from one that just happened.
func TestOldestDeadLetterReportsAge(t *testing.T) {
	tracker := NewPendingPacketTracker()
	if age := tracker.OldestDeadLetteredTimeout(time.Now()); age != 0 {
		t.Fatalf("empty queue must report zero age, got %s", age)
	}

	packet := trackedPacket(1)
	tracker.Add(packet, 1)
	deadAt := time.Now()
	for i := 0; i <= maxTimeoutAttempts; i++ {
		tracker.recordTimeoutFailureForTest(packet.SourceClient, packet.Sequence, deadAt)
	}
	if age := tracker.OldestDeadLetteredTimeout(deadAt.Add(time.Hour)); age < time.Hour {
		t.Fatalf("timeout dead-letter age = %s, want at least an hour", age)
	}
}
