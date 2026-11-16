package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"relayer/chain"
)

func TestDropBackoffDoublesAndIsBounded(t *testing.T) {
	want := []time.Duration{
		10 * time.Minute, 20 * time.Minute, 40 * time.Minute, 80 * time.Minute,
		160 * time.Minute, 320 * time.Minute, 6 * time.Hour, 6 * time.Hour,
	}
	for i, hold := range want {
		if got := dropBackoff(i + 1); got != hold {
			t.Errorf("dropBackoff(%d) = %s, want %s", i+1, got, hold)
		}
	}
	// A count far past the cap must not overflow into a short or negative hold.
	if got := dropBackoff(1 << 20); got != dropBackoffMax {
		t.Errorf("dropBackoff(1<<20) = %s, want the cap %s", got, dropBackoffMax)
	}
}

func TestDropLedgerHoldsAPacketUntilItsBackoffExpires(t *testing.T) {
	ledger, err := NewDropLedger("")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Unix(1_700_000_000, 0)
	packet := []byte("ack-7")

	if _, held := ledger.Held(chain.AckPacket, packet, now); held {
		t.Fatal("a packet that was never dropped must not be held")
	}
	entry, err := ledger.Record(chain.AckPacket, 7, packet, "app callback reverted", now)
	if err != nil {
		t.Fatal(err)
	}
	if entry.Attempts != 1 || !entry.NextRetryAt.Equal(now.Add(dropBackoffMin)) {
		t.Fatalf("first drop: attempts=%d next=%s, want 1 and now+%s", entry.Attempts, entry.NextRetryAt, dropBackoffMin)
	}

	if until, held := ledger.Held(chain.AckPacket, packet, now.Add(dropBackoffMin-time.Second)); !held || !until.Equal(entry.NextRetryAt) {
		t.Fatalf("inside the backoff: held=%v until=%s", held, until)
	}
	if _, held := ledger.Held(chain.AckPacket, packet, now.Add(dropBackoffMin)); held {
		t.Fatal("the hold must end when the backoff expires, or a dropped ack is never retried")
	}
	// The key is the event type and the packet: the receive of the same bytes is
	// a different obligation.
	if _, held := ledger.Held(chain.SendPacket, packet, now); held {
		t.Fatal("a drop of the ack must not hold the receive of the same packet")
	}

	second, err := ledger.Record(chain.AckPacket, 7, packet, "app callback reverted", now.Add(dropBackoffMin))
	if err != nil {
		t.Fatal(err)
	}
	if second.Attempts != 2 || second.NextRetryAt.Sub(second.LastDroppedAt) != 2*dropBackoffMin {
		t.Fatalf("second drop: attempts=%d hold=%s, want 2 and %s", second.Attempts, second.NextRetryAt.Sub(second.LastDroppedAt), 2*dropBackoffMin)
	}
	if !second.FirstDroppedAt.Equal(now) {
		t.Fatalf("first-dropped time moved to %s, want it to stay %s", second.FirstDroppedAt, now)
	}
}

func TestDropLedgerClearAndDue(t *testing.T) {
	ledger, _ := NewDropLedger("")
	now := time.Unix(1_700_000_000, 0)
	for seq := uint64(1); seq <= 3; seq++ {
		if _, err := ledger.Record(chain.AckPacket, seq, []byte(fmt.Sprintf("ack-%d", seq)), "reverted", now.Add(time.Duration(seq)*time.Minute)); err != nil {
			t.Fatal(err)
		}
	}
	if due := ledger.Due(now.Add(11 * time.Minute)); len(due) != 1 || due[0].Sequence != 1 {
		t.Fatalf("due = %+v, want only seq 1", due)
	}
	if removed, err := ledger.Clear(chain.AckPacket, []byte("ack-2")); err != nil || !removed {
		t.Fatalf("clear: removed=%v err=%v", removed, err)
	}
	if removed, _ := ledger.Clear(chain.AckPacket, []byte("ack-2")); removed {
		t.Fatal("clearing twice must report nothing removed")
	}
	snapshot := ledger.Snapshot()
	if len(snapshot) != 2 || snapshot[0].Sequence != 1 || snapshot[1].Sequence != 3 {
		t.Fatalf("snapshot = %+v, want seq 1 then 3 (oldest drop first)", snapshot)
	}
}

// A restart must not reset the hold: the audit's point was that restarting did
// not stop the re-submission.
func TestDropLedgerSurvivesARestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "dropped-cosmos-to-eth.json")
	now := time.Now()

	ledger, err := NewDropLedger(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ledger.Record(chain.AckPacket, 7, []byte("ack-7"), "reverted", now); err != nil {
		t.Fatal(err)
	}
	if _, err := ledger.Record(chain.AckPacket, 7, []byte("ack-7"), "reverted", now); err != nil {
		t.Fatal(err)
	}

	reloaded, err := NewDropLedger(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, held := reloaded.Held(chain.AckPacket, []byte("ack-7"), now.Add(dropBackoffMin)); !held {
		t.Fatal("the reloaded ledger lost the hold; a restart would re-submit the dropped packet")
	}
	got := reloaded.Snapshot()
	if len(got) != 1 || got[0].Attempts != 2 || got[0].Reason != "reverted" || string(got[0].Packet) != "ack-7" {
		t.Fatalf("reloaded entry = %+v", got)
	}

	if _, err := reloaded.Clear(chain.AckPacket, []byte("ack-7")); err != nil {
		t.Fatal(err)
	}
	empty, err := NewDropLedger(path)
	if err != nil {
		t.Fatal(err)
	}
	if n := len(empty.Snapshot()); n != 0 {
		t.Fatalf("a cleared entry came back after a reload (%d entries)", n)
	}
}

func TestDropLedgerForgetsEntriesNothingOfferedAgain(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dropped.json")
	ledger, _ := NewDropLedger(path)
	stale := time.Now().Add(-dropLedgerMaxAge - time.Hour)
	if _, err := ledger.Record(chain.AckPacket, 1, []byte("stale"), "reverted", stale); err != nil {
		t.Fatal(err)
	}
	if _, err := ledger.Record(chain.AckPacket, 2, []byte("fresh"), "reverted", time.Now()); err != nil {
		t.Fatal(err)
	}
	reloaded, err := NewDropLedger(path)
	if err != nil {
		t.Fatal(err)
	}
	got := reloaded.Snapshot()
	if len(got) != 1 || got[0].Sequence != 2 {
		t.Fatalf("reloaded = %+v, want only the fresh entry", got)
	}
}

func TestDropLedgerIsBounded(t *testing.T) {
	ledger, _ := NewDropLedger("")
	now := time.Unix(1_700_000_000, 0)
	for i := 0; i <= dropLedgerMaxEntries; i++ {
		if _, err := ledger.Record(chain.AckPacket, uint64(i), []byte(fmt.Sprintf("ack-%d", i)), "reverted", now.Add(time.Duration(i)*time.Second)); err != nil {
			t.Fatal(err)
		}
	}
	snapshot := ledger.Snapshot()
	if len(snapshot) != dropLedgerMaxEntries {
		t.Fatalf("ledger holds %d entries, want the cap %d", len(snapshot), dropLedgerMaxEntries)
	}
	if snapshot[0].Sequence != 1 {
		t.Fatalf("oldest remaining entry is seq %d, want 1: the oldest drop is the one evicted", snapshot[0].Sequence)
	}
}

// A failed write must not turn a drop back into a retry on every pass.
func TestDropLedgerKeepsTheHoldWhenTheWriteFails(t *testing.T) {
	ledger, _ := NewDropLedger(filepath.Join(t.TempDir(), "dropped.json"))
	ledger.write = func(string, []byte) error { return errors.New("disk full") }
	now := time.Unix(1_700_000_000, 0)

	if _, err := ledger.Record(chain.AckPacket, 7, []byte("ack-7"), "reverted", now); err == nil {
		t.Fatal("a failed write must be reported")
	}
	if _, held := ledger.Held(chain.AckPacket, []byte("ack-7"), now.Add(time.Minute)); !held {
		t.Fatal("the hold was lost with the failed write")
	}
}

func TestDropLedgerRejectsACorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dropped.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewDropLedger(path); err == nil {
		t.Fatal("a corrupt ledger must be reported, not silently replaced")
	}
}
