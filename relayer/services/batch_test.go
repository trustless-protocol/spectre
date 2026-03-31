package services

import (
	"sync"
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func makePacket(seqNum uint64, pktType PacketType) Packet {
	return Packet{
		PacketType: pktType,
		Packet: &channeltypesv2.Packet{
			Sequence: seqNum,
		},
	}
}

func TestNewBatchBuilder(t *testing.T) {
	bb := NewBatchBuilder()
	if bb == nil {
		t.Fatal("expected non-nil BatchBuilder")
	}
	if len(bb.packets) != 0 {
		t.Fatalf("expected empty packets, got %d", len(bb.packets))
	}
	if bb.timestamp.IsZero() {
		t.Fatal("expected timestamp to be set")
	}
}

func TestInsertPacket_Single(t *testing.T) {
	bb := NewBatchBuilder()
	pkt := makePacket(1, Send)

	bb.InsertPacket(pkt)

	bb.mtx.Lock()
	defer bb.mtx.Unlock()

	if len(bb.packets) != 1 {
		t.Fatalf("expected 1 packet, got %d", len(bb.packets))
	}
	if bb.packets[0].Packet.Sequence != 1 {
		t.Fatalf("expected sequence 1, got %d", bb.packets[0].Packet.Sequence)
	}
	if bb.packets[0].PacketType != Send {
		t.Fatalf("expected PacketType Send, got %d", bb.packets[0].PacketType)
	}
}

func TestInsertPacket_Multiple(t *testing.T) {
	bb := NewBatchBuilder()
	count := 5
	for i := 0; i < count; i++ {
		bb.InsertPacket(makePacket(uint64(i+1), Send))
	}

	bb.mtx.Lock()
	defer bb.mtx.Unlock()

	if len(bb.packets) != count {
		t.Fatalf("expected %d packets, got %d", count, len(bb.packets))
	}
	for i := 0; i < count; i++ {
		if bb.packets[i].Packet.Sequence != uint64(i+1) {
			t.Fatalf("packet %d: expected sequence %d, got %d", i, i+1, bb.packets[i].Packet.Sequence)
		}
	}
}

func TestInsertPacket_Concurrent(t *testing.T) {
	bb := NewBatchBuilder()
	goroutines := 50
	packetsPerGoroutine := 20
	total := goroutines * packetsPerGoroutine

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(offset int) {
			defer wg.Done()
			for i := 0; i < packetsPerGoroutine; i++ {
				bb.InsertPacket(makePacket(uint64(offset*packetsPerGoroutine+i), Send))
			}
		}(g)
	}

	wg.Wait()

	bb.mtx.Lock()
	defer bb.mtx.Unlock()

	if len(bb.packets) != total {
		t.Fatalf("expected %d packets, got %d", total, len(bb.packets))
	}
}

func TestClearBatch(t *testing.T) {
	bb := NewBatchBuilder()
	for i := 0; i < 5; i++ {
		bb.InsertPacket(makePacket(uint64(i+1), Send))
	}

	bb.mtx.Lock()
	if len(bb.packets) != 5 {
		bb.mtx.Unlock()
		t.Fatalf("expected 5 packets before clear, got %d", len(bb.packets))
	}
	bb.mtx.Unlock()

	bb.ClearBatch()

	bb.mtx.Lock()
	defer bb.mtx.Unlock()

	if len(bb.packets) != 0 {
		t.Fatalf("expected 0 packets after clear, got %d", len(bb.packets))
	}
}

func TestCheckBatch_ExceedsBatchSize(t *testing.T) {
	bb := NewBatchBuilder()
	config := BatchConfig{
		BatchSize:    3,
		BatchPeriods: time.Minute * 10, // long period so only size triggers
	}

	// Insert more than BatchSize packets (> 3, so need 4+)
	for i := 0; i < 5; i++ {
		bb.InsertPacket(makePacket(uint64(i+1), Send))
	}

	ch := make(chan BatchPackets, 1)
	bb.CheckBatch(config, ch)

	select {
	case batch := <-ch:
		if len(batch.Packets) != 5 {
			t.Fatalf("expected 5 packets in batch, got %d", len(batch.Packets))
		}
	default:
		t.Fatal("expected batch to be sent to channel, but channel was empty")
	}

	// Verify packets were cleared
	bb.mtx.Lock()
	defer bb.mtx.Unlock()
	if len(bb.packets) != 0 {
		t.Fatalf("expected packets to be cleared after check, got %d", len(bb.packets))
	}
}

func TestCheckBatch_PastBatchPeriods(t *testing.T) {
	bb := NewBatchBuilder()
	// Set timestamp in the past so the period has elapsed
	bb.timestamp = time.Now().Add(-time.Minute * 5)

	config := BatchConfig{
		BatchSize:    10, // well above our packet count
		BatchPeriods: time.Second * 1,
	}

	bb.InsertPacket(makePacket(1, Send))
	bb.InsertPacket(makePacket(2, Send))

	ch := make(chan BatchPackets, 1)
	bb.CheckBatch(config, ch)

	select {
	case batch := <-ch:
		if len(batch.Packets) != 2 {
			t.Fatalf("expected 2 packets in batch, got %d", len(batch.Packets))
		}
	default:
		t.Fatal("expected batch to be sent to channel when period elapsed, but channel was empty")
	}

	bb.mtx.Lock()
	defer bb.mtx.Unlock()
	if len(bb.packets) != 0 {
		t.Fatalf("expected packets to be cleared after check, got %d", len(bb.packets))
	}
}

func TestCheckBatch_BelowSizeAndBeforePeriod(t *testing.T) {
	bb := NewBatchBuilder()
	// Timestamp is now (just created), so period hasn't elapsed
	config := BatchConfig{
		BatchSize:    10,
		BatchPeriods: time.Minute * 10,
	}

	bb.InsertPacket(makePacket(1, Send))
	bb.InsertPacket(makePacket(2, Send))

	ch := make(chan BatchPackets, 1)
	bb.CheckBatch(config, ch)

	select {
	case <-ch:
		t.Fatal("expected no batch to be sent, but channel had data")
	default:
		// correct: nothing sent
	}

	// Verify packets are still present (not cleared)
	bb.mtx.Lock()
	defer bb.mtx.Unlock()
	if len(bb.packets) != 2 {
		t.Fatalf("expected packets to remain, got %d", len(bb.packets))
	}
}
