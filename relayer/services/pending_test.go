package services

import (
	"sync"
	"testing"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func TestPendingTracker_AddRemove(t *testing.T) {
	pt := NewPendingPacketTracker()
	p := channeltypesv2.Packet{SourceClient: "src-0", Sequence: 1}

	pt.Add(p, 100)
	if pt.Len() != 1 {
		t.Fatalf("expected len 1, got %d", pt.Len())
	}

	pt.Remove("src-0", 1)
	if pt.Len() != 0 {
		t.Fatalf("expected len 0 after remove, got %d", pt.Len())
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
				pt.Remove("src-0", seq)
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
