package services

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func TestPersistentPendingPacketTrackerRehydratesActiveAndTombstones(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pending.json")
	tracker, err := NewPersistentPendingPacketTracker(path)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Truncate(time.Nanosecond)
	active := channeltypesv2.Packet{SourceClient: "src-0", DestinationClient: "dst-0", Sequence: 1, TimeoutTimestamp: 100}
	dead := channeltypesv2.Packet{SourceClient: "src-0", DestinationClient: "dst-0", Sequence: 2, TimeoutTimestamp: 200}
	tracker.Add(active, 101)
	tracker.recordTimeoutFailureForTest(active.SourceClient, active.Sequence, now)
	tracker.deferTimeoutRetryForTest(active.SourceClient, active.Sequence, now)
	tracker.Add(dead, 202)
	for i := 0; i <= maxTimeoutAttempts; i++ {
		tracker.recordTimeoutFailureForTest(dead.SourceClient, dead.Sequence, now)
	}

	restored, err := NewPersistentPendingPacketTracker(path)
	if err != nil {
		t.Fatal(err)
	}
	got := restored.GetAll()
	if len(got) != 1 {
		t.Fatalf("restored active packets = %d, want 1", len(got))
	}
	if got[0].TimeoutAttempts != 1 || got[0].Deferrals != 1 || got[0].BlockNumber != 101 || got[0].NotBefore.IsZero() {
		t.Fatalf("restored active retry state = %#v", got[0])
	}
	if restored.DeadLetteredTimeouts() != 1 {
		t.Fatalf("restored tombstones = %d, want 1", restored.DeadLetteredTimeouts())
	}

	// Startup/gap recovery may replay this event immediately. The rehydrated
	// tombstone must block it instead of granting a fresh retry budget.
	restored.Add(dead, 303)
	if restored.Len() != 1 || restored.DeadLetteredTimeouts() != 1 {
		t.Fatalf("replayed dead letter after restore: active=%d dead=%d", restored.Len(), restored.DeadLetteredTimeouts())
	}
}

func TestNewPersistentBatchBuilderRestoresBeforeReturn(t *testing.T) {
	dir := t.TempDir()
	first, err := NewPersistentBatchBuilder(dir)
	if err != nil {
		t.Fatal(err)
	}
	packet := channeltypesv2.Packet{SourceClient: "l2-client", DestinationClient: "cosmos-client", Sequence: 9, TimeoutTimestamp: 100}
	first.L2PendingTracker.Add(packet, 99)
	first.L2PendingTracker.recordTimeoutFailureForTest(packet.SourceClient, packet.Sequence, time.Now())

	restarted, err := NewPersistentBatchBuilder(dir)
	if err != nil {
		t.Fatal(err)
	}
	got := restarted.L2PendingTracker.GetAll()
	if len(got) != 1 || got[0].TimeoutAttempts != 1 || got[0].BlockNumber != 99 {
		t.Fatalf("batch builder did not restore L2 tracker before return: %#v", got)
	}
	if restarted.cosmosInFlight == nil || restarted.ethInFlight == nil {
		t.Fatal("persistent batch builder left in-flight recovery maps nil")
	}
}

func TestPersistentPendingPacketTrackerRollsBackFailedMutation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pending.json")
	tracker, err := NewPersistentPendingPacketTracker(path)
	if err != nil {
		t.Fatal(err)
	}
	packet := channeltypesv2.Packet{SourceClient: "src-0", DestinationClient: "dst-0", Sequence: 1, TimeoutTimestamp: 100}
	tracker.Add(packet, 10)
	tracker.writeState = func(string, []byte) error { return errors.New("disk full") }

	if dead := tracker.recordTimeoutFailureForTest(packet.SourceClient, packet.Sequence, time.Now()); dead {
		t.Fatal("failed persistence reported a committed dead letter")
	}
	got := tracker.GetAll()
	if len(got) != 1 || got[0].TimeoutAttempts != 0 || !got[0].NotBefore.IsZero() {
		t.Fatalf("failed persistence changed in-memory retry state: %#v", got)
	}

	restored, err := NewPersistentPendingPacketTracker(path)
	if err != nil {
		t.Fatal(err)
	}
	persisted := restored.GetAll()
	if len(persisted) != 1 || persisted[0].TimeoutAttempts != 0 || !persisted[0].NotBefore.IsZero() {
		t.Fatalf("failed mutation leaked to durable retry state: %#v", persisted)
	}
}

func TestPersistentPendingPacketTrackerAddReportsFailedMutation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pending.json")
	tracker, err := NewPersistentPendingPacketTracker(path)
	if err != nil {
		t.Fatal(err)
	}
	tracker.writeState = func(string, []byte) error { return errors.New("disk full") }
	packet := channeltypesv2.Packet{
		SourceClient: "src-0", DestinationClient: "dst-0",
		Sequence: 1, TimeoutTimestamp: 100,
	}

	if tracker.Add(packet, 10) {
		t.Fatal("Add reported success after its durable mutation was rolled back")
	}
	if tracker.Len() != 0 {
		t.Fatalf("failed Add left %d packet(s) in memory, want 0", tracker.Len())
	}

	restored, err := NewPersistentPendingPacketTracker(path)
	if err != nil {
		t.Fatal(err)
	}
	if restored.Len() != 0 {
		t.Fatalf("failed Add left %d packet(s) on disk, want 0", restored.Len())
	}
}

func TestDeferTimeoutRetriesIfCurrentPersistsBatchOnce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pending.json")
	tracker, err := NewPersistentPendingPacketTracker(path)
	if err != nil {
		t.Fatal(err)
	}
	for sequence := uint64(1); sequence <= 5; sequence++ {
		tracker.Add(channeltypesv2.Packet{
			SourceClient: "src-0", DestinationClient: "dst-0",
			Sequence: sequence, TimeoutTimestamp: 100,
		}, sequence)
	}

	writes := 0
	tracker.writeState = func(string, []byte) error {
		writes++
		return nil
	}
	counts := tracker.DeferTimeoutRetriesIfCurrent(tracker.GetAll(), time.Now())
	if writes != 1 {
		t.Fatalf("persist writes = %d, want 1 for the full deferral batch", writes)
	}
	for i, count := range counts {
		if count != 1 {
			t.Fatalf("deferral count[%d] = %d, want 1", i, count)
		}
	}
}
