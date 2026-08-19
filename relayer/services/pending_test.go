package services

import (
	"reflect"
	"sync"
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func TestPendingTracker_AddRemove(t *testing.T) {
	pt := NewPendingPacketTracker()
	p := channeltypesv2.Packet{SourceClient: "src-0", Sequence: 1}

	pt.Add(p, 100)
	if pt.Len() != 1 {
		t.Fatalf("expected len 1, got %d", pt.Len())
	}

	pt.removeCurrentForTest("src-0", 1)
	if pt.Len() != 0 {
		t.Fatalf("expected len 0 after remove, got %d", pt.Len())
	}
}

func TestPendingTracker_AddPreservesTimeoutRetryState(t *testing.T) {
	pt := NewPendingPacketTracker()
	p := channeltypesv2.Packet{SourceClient: "src-0", Sequence: 1, TimeoutTimestamp: 1}
	pt.Add(p, 100)
	now := time.Now()
	pt.recordTimeoutFailureForTest(p.SourceClient, p.Sequence, now)
	pt.deferTimeoutRetryForTest(p.SourceClient, p.Sequence, now)

	// Mirrors the relay module re-offering a requeued event before every attempt.
	pt.Add(p, 200)

	got := pt.GetAll()
	if len(got) != 1 {
		t.Fatalf("tracked packets = %d, want 1", len(got))
	}
	if got[0].TimeoutAttempts != 1 || got[0].Deferrals != 1 {
		t.Fatalf("retry state after duplicate Add = attempts:%d deferrals:%d, want 1:1",
			got[0].TimeoutAttempts, got[0].Deferrals)
	}
	if !got[0].NotBefore.After(now) {
		t.Fatalf("duplicate Add cleared retry deadline: %s", got[0].NotBefore)
	}
	if got[0].BlockNumber != 100 {
		t.Fatalf("duplicate Add replaced block number = %d, want original 100", got[0].BlockNumber)
	}
}

func TestPendingTracker_AddTreatsEquivalentProtobufRepresentationsAsSamePacket(t *testing.T) {
	pt := NewPendingPacketTracker()
	nilPayloads := channeltypesv2.Packet{
		SourceClient: "src-0", DestinationClient: "dst-0", Sequence: 1,
		TimeoutTimestamp: 100, Payloads: nil,
	}
	pt.Add(nilPayloads, 100)
	now := time.Now()
	pt.recordTimeoutFailureForTest(nilPayloads.SourceClient, nilPayloads.Sequence, now)

	emptyPayloads := nilPayloads
	emptyPayloads.Payloads = []channeltypesv2.Payload{}
	pt.Add(emptyPayloads, 200)

	got := pt.GetAll()
	if len(got) != 1 || got[0].TimeoutAttempts != 1 {
		t.Fatalf("equivalent protobuf replay reset retry state: %#v", got)
	}
	if got[0].BlockNumber != 100 {
		t.Fatalf("equivalent protobuf replay replaced block number = %d, want 100", got[0].BlockNumber)
	}
}

func TestPendingTracker_DeadLetterTombstoneBlocksReplayUntilRetentionExpires(t *testing.T) {
	pt := NewPendingPacketTracker()
	p := channeltypesv2.Packet{SourceClient: "src-0", Sequence: 1, TimeoutTimestamp: 1}
	pt.Add(p, 100)
	now := time.Now()
	for i := 0; i <= maxTimeoutAttempts; i++ {
		pt.recordTimeoutFailureForTest(p.SourceClient, p.Sequence, now)
	}
	if pt.Len() != 0 || pt.DeadLetteredTimeouts() != 1 {
		t.Fatalf("dead-letter state = active:%d dead:%d, want 0:1", pt.Len(), pt.DeadLetteredTimeouts())
	}

	// A replayed source event must not silently receive a fresh retry budget.
	pt.Add(p, 200)
	if pt.Len() != 0 {
		t.Fatal("replayed event revived a dead-lettered timeout")
	}

	pt.mtx.Lock()
	pt.deadLettered[0].DeadLetteredAt = now.Add(-deadLetterRetention)
	pt.mtx.Unlock()
	pt.PurgeStaleWithoutTimeout(time.Hour)
	got := pt.GetAll()
	if len(got) != 1 || got[0].TimeoutAttempts != 0 || got[0].Deferrals != 0 || !got[0].NotBefore.IsZero() {
		t.Fatalf("renewed retry state = %#v, want fresh active packet", got)
	}
	if pt.DeadLetteredTimeouts() != 0 {
		t.Fatal("expired dead letter remained in the gauge")
	}
}

func TestPendingTracker_AddReplacesRenewedPacketWhenSequenceIsReused(t *testing.T) {
	pt := NewPendingPacketTracker()
	old := channeltypesv2.Packet{SourceClient: "src-0", DestinationClient: "dst-0", Sequence: 7, TimeoutTimestamp: 1000}
	pt.Add(old, 100)
	now := time.Now()
	for i := 0; i <= maxTimeoutAttempts; i++ {
		pt.recordTimeoutFailureForTest(old.SourceClient, old.Sequence, now)
	}

	pt.mtx.Lock()
	pt.deadLettered[0].DeadLetteredAt = now.Add(-deadLetterRetention)
	pt.mtx.Unlock()
	pt.PurgeStaleWithoutTimeout(time.Hour)

	// A counterparty client can be re-created while reusing its router client ID,
	// restarting packet sequences. The new packet must replace the expired
	// tombstone's renewed packet rather than being silently ignored as a replay.
	new := channeltypesv2.Packet{SourceClient: "src-0", DestinationClient: "dst-1", Sequence: 7, TimeoutTimestamp: 999999}
	pt.Add(new, 500)
	got := pt.GetAll()
	if len(got) != 1 {
		t.Fatalf("tracked packets = %d, want 1", len(got))
	}
	if !reflect.DeepEqual(got[0].Packet, new) || got[0].BlockNumber != 500 {
		t.Fatalf("tracked packet = %#v at block %d, want new packet at block 500", got[0].Packet, got[0].BlockNumber)
	}
}

func TestPendingTracker_StaleScanCannotMutateReplacement(t *testing.T) {
	pt := NewPendingPacketTracker()
	old := channeltypesv2.Packet{SourceClient: "src-0", DestinationClient: "dst-old", Sequence: 7, TimeoutTimestamp: 1000}
	pt.Add(old, 100)
	stale := pt.GetDue(time.Now())[0]

	newPacket := channeltypesv2.Packet{SourceClient: "src-0", DestinationClient: "dst-new", Sequence: 7, TimeoutTimestamp: 999999}
	pt.Add(newPacket, 500)

	// Every mutation a timeout scan can apply must be conditional on the packet
	// identity it scanned, not merely the reused lookup key. The same applies to
	// a delayed terminal source event carrying the old packet.
	pt.RemovePacketIfCurrent(old)
	pt.DeferTimeoutRetryIfCurrent(stale, time.Now())
	pt.RecordTimeoutFailureIfCurrent(stale, time.Now())
	pt.ClearDeferralsIfCurrent(stale)
	pt.RemoveIfCurrent(stale)

	got := pt.GetAll()
	if len(got) != 1 {
		t.Fatalf("stale scan removed replacement: tracked=%d", len(got))
	}
	if got[0].identity != identifyPacket(newPacket) || got[0].BlockNumber != 500 {
		t.Fatalf("tracked packet = %#v at block %d, want replacement at block 500", got[0].Packet, got[0].BlockNumber)
	}
	if got[0].TimeoutAttempts != 0 || got[0].Deferrals != 0 || !got[0].NotBefore.IsZero() {
		t.Fatalf("stale scan mutated replacement retry state: %#v", got[0])
	}
}

func TestPendingTracker_GetAllSnapshot(t *testing.T) {
	pt := NewPendingPacketTracker()
	p1 := channeltypesv2.Packet{SourceClient: "src-0", Sequence: 1}
	p2 := channeltypesv2.Packet{SourceClient: "src-0", Sequence: 2}

	pt.Add(p1, 100)
	pt.Add(p2, 200)

	all := pt.GetAll()
	if len(all) != 2 {
		t.Fatalf("expected 2 items, got %d", len(all))
	}

	// Mutating the returned slice must not affect internal state
	all = append(all, pendingPacketInfo{})
	if pt.Len() != 2 {
		t.Fatalf("expected internal len 2 after external append, got %d", pt.Len())
	}
}

func TestPendingTracker_Concurrent(t *testing.T) {
	pt := NewPendingPacketTracker()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(seq uint64) {
			defer wg.Done()
			p := channeltypesv2.Packet{SourceClient: "src-0", Sequence: seq}
			pt.Add(p, seq)
			if seq%2 == 0 {
				pt.removeCurrentForTest("src-0", seq)
			}
		}(uint64(i))
	}
	wg.Wait()

	// All even sequences should have been removed
	for i := 0; i < 50; i++ {
		exists := false
		for _, info := range pt.GetAll() {
			if info.Packet.Sequence == uint64(i) {
				exists = true
				break
			}
		}
		if i%2 == 0 && exists {
			t.Fatalf("seq %d should have been removed", i)
		}
		if i%2 == 1 && !exists {
			t.Fatalf("seq %d should still exist", i)
		}
	}
}

func TestPendingTracker_PurgeStaleWithoutTimeout(t *testing.T) {
	pt := NewPendingPacketTracker()
	noTimeout := channeltypesv2.Packet{SourceClient: "src-0", Sequence: 1}
	withTimeout := channeltypesv2.Packet{SourceClient: "src-0", Sequence: 2, TimeoutTimestamp: uint64(time.Now().Add(2 * time.Hour).Unix())}

	pt.Add(noTimeout, 100)
	pt.Add(withTimeout, 200)

	pt.mtx.Lock()
	for key, info := range pt.packets {
		info.ObservedAt = time.Now().Add(-2 * time.Hour)
		pt.packets[key] = info
	}
	pt.mtx.Unlock()

	pt.PurgeStaleWithoutTimeout(1 * time.Hour)

	all := pt.GetAll()
	if len(all) != 1 {
		t.Fatalf("expected len 1 after targeted purge, got %d", len(all))
	}
	if all[0].Packet.Sequence != 2 {
		t.Fatalf("expected packet with timeout to remain, got seq=%d", all[0].Packet.Sequence)
	}
}
