package services

import (
	"context"
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func waitingEthPkt(seq uint64) EthPacket {
	return EthPacket{Type: EthSend, Packet: &channeltypesv2.Packet{Sequence: seq, SourceClient: "c"}}
}

// TestNextWaitingBackoff verifies the exponential growth + cap (3s,6s,12s,15s…).
func TestNextWaitingBackoff(t *testing.T) {
	want := []time.Duration{3 * time.Second, 6 * time.Second, 12 * time.Second, 15 * time.Second, 15 * time.Second}
	for i, w := range want {
		if got := nextWaitingBackoff(i + 1); got != w {
			t.Fatalf("attempt %d: got %s want %s", i+1, got, w)
		}
	}
}

// TestCheckEth_WaitingHeldBack verifies a packet under its NotBefore backoff is
// NOT flushed (quiet wait), and IS flushed once the backoff elapses.
func TestCheckEth_WaitingHeldBack(t *testing.T) {
	b := NewBatchBuilder()
	cfg := BatchConfig{BatchSize: 5, BatchPeriods: 0} // period 0 => time-based flush always ready

	// A waiting packet (NotBefore in the future) must be held back: no flush.
	b.RequeueEthWaiting([]EthPacket{waitingEthPkt(1)})
	ch := make(chan EthBatch, 1)
	b.CheckEth(context.Background(), cfg, ch)
	select {
	case <-ch:
		t.Fatal("waiting packet must not flush before its backoff elapses")
	default:
	}

	// Force the backoff to have elapsed and re-check: now it flushes.
	b.ethMtx.Lock()
	for i := range b.ethPackets {
		b.ethPackets[i].NotBefore = time.Now().Add(-time.Second)
	}
	b.ethMtx.Unlock()
	b.CheckEth(context.Background(), cfg, ch)
	select {
	case batch := <-ch:
		if len(batch.Packets) != 1 || batch.Packets[0].Packet.Sequence != 1 {
			t.Fatalf("expected the ready packet, got %+v", batch.Packets)
		}
	default:
		t.Fatal("ready packet must flush once its backoff elapsed")
	}
}

// TestRequeueEthWaiting_GrowsBackoff verifies repeated waiting re-queues grow the
// per-packet delay (so a long finality wait polls ever more sparsely).
func TestRequeueEthWaiting_GrowsBackoff(t *testing.T) {
	b := NewBatchBuilder()
	p := waitingEthPkt(1)
	b.RequeueEthWaiting([]EthPacket{p})
	b.ethMtx.Lock()
	first := b.ethPackets[0].NotBefore
	attempts1 := b.ethPackets[0].waitAttempts
	pk := b.ethPackets[0]
	b.ethPackets = nil
	b.ethMtx.Unlock()
	if attempts1 != 1 {
		t.Fatalf("first waiting re-queue attempts=%d, want 1", attempts1)
	}

	b.RequeueEthWaiting([]EthPacket{pk})
	b.ethMtx.Lock()
	second := b.ethPackets[0].NotBefore
	attempts2 := b.ethPackets[0].waitAttempts
	b.ethMtx.Unlock()
	if attempts2 != 2 {
		t.Fatalf("second waiting re-queue attempts=%d, want 2", attempts2)
	}
	if !second.After(first) {
		t.Fatalf("backoff must grow: second NotBefore %s not after first %s", second, first)
	}
}
