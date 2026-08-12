package services

import (
	"context"
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func cancelTestConfig() BatchConfig {
	// BatchSize 1 with a zero period flushes on the first packet, so the test
	// reaches the handoff immediately.
	return BatchConfig{BatchSize: 1, BatchPeriods: 0}
}

func cosmosPacketSeq(seq uint64) CosmosPacket {
	return CosmosPacket{
		Type:   CosmosSend,
		Packet: &channeltypesv2.Packet{Sequence: seq, SourceClient: "08-wasm-0", DestinationClient: "client-0"},
	}
}

func ethPacketSeq(seq uint64) EthPacket {
	return EthPacket{
		Type:   EthSend,
		Packet: &channeltypesv2.Packet{Sequence: seq, SourceClient: "client-0", DestinationClient: "08-wasm-0"},
	}
}

// TestCheckCosmosReturnsAndRestoresOnCancel pins the property the fix exists
// for: the chunk is sliced off the queue BEFORE the handoff, so a send that
// never completes must not take those packets with it.
//
// The channel is deliberately unbuffered with no receiver, which is what a
// shutting-down consumer looks like. Without the ctx guard this blocks forever
// and the packet is gone from the queue — the test hangs and then fails on the
// 5s deadline.
func TestCheckCosmosReturnsAndRestoresOnCancel(t *testing.T) {
	b := NewBatchBuilder()
	b.AddCosmos(cosmosPacketSeq(1))

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan CosmosBatch) // no receiver, no buffer

	done := make(chan struct{})
	go func() {
		b.CheckCosmos(ctx, cancelTestConfig(), ch)
		close(done)
	}()

	// Let CheckCosmos reach the blocked send, then shut down.
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("CheckCosmos did not return after cancellation: the send is unguarded")
	}

	b.cosmosMtx.Lock()
	got := len(b.cosmosPackets)
	var seq uint64
	if got == 1 {
		seq = b.cosmosPackets[0].Packet.Sequence
	}
	b.cosmosMtx.Unlock()

	if got != 1 {
		t.Fatalf("packet count after cancelled flush = %d, want 1 (the chunk must go back on the queue)", got)
	}
	if seq != 1 {
		t.Fatalf("restored packet sequence = %d, want 1", seq)
	}
}

// TestCheckEthReturnsAndRestoresOnCancel — the ETH mirror of the above. The two
// paths are edited in lockstep by convention; this keeps that honest.
func TestCheckEthReturnsAndRestoresOnCancel(t *testing.T) {
	b := NewBatchBuilder()
	b.AddEth(ethPacketSeq(7))

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan EthBatch)

	done := make(chan struct{})
	go func() {
		b.CheckEth(ctx, cancelTestConfig(), ch)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("CheckEth did not return after cancellation: the send is unguarded")
	}

	b.ethMtx.Lock()
	got := len(b.ethPackets)
	var seq uint64
	if got == 1 {
		seq = b.ethPackets[0].Packet.Sequence
	}
	b.ethMtx.Unlock()

	if got != 1 {
		t.Fatalf("packet count after cancelled flush = %d, want 1 (the chunk must go back on the queue)", got)
	}
	if seq != 7 {
		t.Fatalf("restored packet sequence = %d, want 7", seq)
	}
}

// TestCheckRestoresChunkAheadOfNewerArrivals: a restored chunk must go back at
// the HEAD, so shutdown does not silently reorder the queue relative to packets
// that arrived while the flush was blocked.
func TestCheckRestoresChunkAheadOfNewerArrivals(t *testing.T) {
	b := NewBatchBuilder()
	b.AddCosmos(cosmosPacketSeq(1))

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan CosmosBatch)

	done := make(chan struct{})
	go func() {
		b.CheckCosmos(ctx, cancelTestConfig(), ch)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	b.AddCosmos(cosmosPacketSeq(2)) // arrives while the flush is parked
	cancel()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("CheckCosmos did not return after cancellation")
	}

	b.cosmosMtx.Lock()
	defer b.cosmosMtx.Unlock()
	if len(b.cosmosPackets) != 2 {
		t.Fatalf("packet count = %d, want 2", len(b.cosmosPackets))
	}
	if got := b.cosmosPackets[0].Packet.Sequence; got != 1 {
		t.Fatalf("head sequence = %d, want 1 (the restored chunk must precede newer arrivals)", got)
	}
}

// TestCheckCosmosStillFlushesWhenTheConsumerIsReady guards against the fix
// turning every flush into a cancellation: a live consumer must still receive.
func TestCheckCosmosStillFlushesWhenTheConsumerIsReady(t *testing.T) {
	b := NewBatchBuilder()
	b.AddCosmos(cosmosPacketSeq(3))

	ch := make(chan CosmosBatch, 1)
	b.CheckCosmos(context.Background(), cancelTestConfig(), ch)

	select {
	case batch := <-ch:
		if len(batch.Packets) != 1 || batch.Packets[0].Packet.Sequence != 3 {
			t.Fatalf("unexpected batch %+v", batch)
		}
	default:
		t.Fatal("no batch delivered to a ready consumer")
	}

	b.cosmosMtx.Lock()
	defer b.cosmosMtx.Unlock()
	if len(b.cosmosPackets) != 0 {
		t.Fatalf("queue = %d after a successful flush, want 0", len(b.cosmosPackets))
	}
}

// TestCancelledFlushStrandsNoPacketAtProductionCapacity is the case the first
// round of tests missed, reported by @neitdung.
//
// Those tests used an unbuffered channel written by hand, so they only ever
// exercised the ctx.Done() arm. Production used a 16-slot buffer, where the
// send ARM WINS: it is ready, so select takes it, the batch lands in the buffer,
// and the consumer — selecting on the same cancellation — returns without ever
// receiving it. The chunk was already sliced off the queue, so it existed
// nowhere. Guarding the send with ctx.Done() does not help at all in that shape.
//
// The channel here is built exactly as the source bridges build it, from the
// same constant, so setting BatchHandoffCapacity back above zero fails this test
// rather than silently reopening the hole.
func TestCancelledFlushStrandsNoPacketAtProductionCapacity(t *testing.T) {
	t.Parallel()

	t.Run("cosmos", func(t *testing.T) {
		t.Parallel()

		b := NewBatchBuilder()
		b.AddCosmos(cosmosPacketSeq(1))

		ctx, cancel := context.WithCancel(context.Background())
		// Exactly what chain/cosmos/source.go builds. No receiver: a consumer
		// that has already returned on shutdown.
		ch := make(chan CosmosBatch, BatchHandoffCapacity)

		done := make(chan struct{})
		go func() {
			b.CheckCosmos(ctx, cancelTestConfig(), ch)
			close(done)
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("CheckCosmos did not return after cancellation")
		}

		// The packet must be back in the builder -- not sitting in a channel
		// buffer that nothing will ever read.
		b.cosmosMtx.Lock()
		got := len(b.cosmosPackets)
		b.cosmosMtx.Unlock()
		if got != 1 {
			t.Fatalf("queue holds %d packet(s), want 1; the batch was stranded in the handoff "+
				"(BatchHandoffCapacity = %d must be 0)", got, BatchHandoffCapacity)
		}
	})

	t.Run("eth", func(t *testing.T) {
		t.Parallel()

		b := NewBatchBuilder()
		b.AddEth(ethPacketSeq(7))

		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan EthBatch, BatchHandoffCapacity)

		done := make(chan struct{})
		go func() {
			b.CheckEth(ctx, cancelTestConfig(), ch)
			close(done)
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("CheckEth did not return after cancellation")
		}

		b.ethMtx.Lock()
		got := len(b.ethPackets)
		b.ethMtx.Unlock()
		if got != 1 {
			t.Fatalf("queue holds %d packet(s), want 1; the batch was stranded in the handoff "+
				"(BatchHandoffCapacity = %d must be 0)", got, BatchHandoffCapacity)
		}
	})
}
