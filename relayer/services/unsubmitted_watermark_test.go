package services

import (
	"context"
	"testing"
	"time"
)

func TestLowestUnsubmittedIncludesQueueAndInFlight(t *testing.T) {
	bb := NewBatchBuilder()
	a := makeCosmosPacket(1, CosmosSend)
	a.BlockNumber = 300
	b := makeCosmosPacket(2, CosmosSend)
	b.BlockNumber = 200
	bb.AddCosmos(a)
	bb.AddCosmos(b)

	ch := make(chan CosmosBatch, 1)
	bb.CheckCosmos(context.Background(), BatchConfig{BatchSize: 1, BatchPeriods: 0}, ch)
	batch := <-ch
	if got := bb.LowestUnsubmittedCosmosHeight(); got != 200 {
		t.Fatalf("lowest while queued+in-flight = %d, want 200", got)
	}
	bb.ReleaseCosmosInFlight(batch)
	if got := bb.LowestUnsubmittedCosmosHeight(); got != 200 {
		t.Fatalf("lowest after release = %d, want queued height 200", got)
	}
}

func TestEqualInFlightLowsAreCountedSeparately(t *testing.T) {
	bb := NewBatchBuilder()
	for seq := uint64(1); seq <= 2; seq++ {
		packet := makeCosmosPacket(seq, CosmosSend)
		packet.BlockNumber = 50
		bb.AddCosmos(packet)
	}
	ch := make(chan CosmosBatch, 2)
	bb.CheckCosmos(context.Background(), BatchConfig{BatchSize: 1, BatchPeriods: 0}, ch)
	bb.CheckCosmos(context.Background(), BatchConfig{BatchSize: 1, BatchPeriods: 0}, ch)
	first, second := <-ch, <-ch
	bb.ReleaseCosmosInFlight(first)
	if got := bb.LowestUnsubmittedCosmosHeight(); got != 50 {
		t.Fatalf("one equal-low chunk still in flight, got %d", got)
	}
	bb.ReleaseCosmosInFlight(second)
	if got := bb.LowestUnsubmittedCosmosHeight(); got != 0 {
		t.Fatalf("all chunks released, got %d", got)
	}
}

func TestRequeueBeforeReleaseKeepsWatermark(t *testing.T) {
	bb := NewBatchBuilder()
	packet := makeEthPacket(1, EthSend)
	packet.BlockNumber = 400
	bb.AddEth(packet)
	ch := make(chan EthBatch, 1)
	bb.CheckEth(context.Background(), BatchConfig{BatchSize: 1, BatchPeriods: 0}, ch)
	batch := <-ch
	bb.RequeueEthWaiting(batch.Packets)
	bb.ReleaseEthInFlight(batch)
	if got := bb.LowestUnsubmittedEthHeight(); got != 400 {
		t.Fatalf("requeued packet watermark = %d, want 400", got)
	}
}

func TestCancelledHandoffRestoresQueueAndClearsInFlight(t *testing.T) {
	bb := NewBatchBuilder()
	packet := makeCosmosPacket(1, CosmosSend)
	packet.BlockNumber = 77
	bb.AddCosmos(packet)
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan CosmosBatch)
	done := make(chan struct{})
	go func() {
		bb.CheckCosmos(ctx, BatchConfig{BatchSize: 1, BatchPeriods: 0}, ch)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("cancelled handoff did not return")
	}
	if got := bb.LowestUnsubmittedCosmosHeight(); got != 77 {
		t.Fatalf("restored watermark = %d, want 77", got)
	}
	bb.cosmosMtx.Lock()
	inFlight := len(bb.cosmosInFlight)
	queued := len(bb.cosmosPackets)
	bb.cosmosMtx.Unlock()
	if inFlight != 0 || queued != 1 {
		t.Fatalf("after cancel: inFlight=%d queued=%d", inFlight, queued)
	}
}
