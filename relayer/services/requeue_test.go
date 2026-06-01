package services

import "testing"

// Issue #80 (+ review): failed-batch packets must be re-queued (not dropped),
// prepended ahead of newer arrivals. Only PERMANENT (deterministic revert)
// failures consume the retry budget and eventually dead-letter; TRANSIENT
// (infra) failures re-queue indefinitely so a valid packet is never lost.

func TestRequeueCosmosTransient_NoBudgetConsumed(t *testing.T) {
	bb := NewBatchBuilder()
	p := makeCosmosPacket(1, CosmosSend)
	p.Retries = 0

	// Many transient re-queues must never increment Retries nor dead-letter.
	for range maxPacketRetries * 3 {
		bb.cosmosMtx.Lock()
		bb.cosmosPackets = nil
		bb.cosmosMtx.Unlock()
		bb.RequeueCosmosTransient([]CosmosPacket{p})
		bb.cosmosMtx.Lock()
		p = bb.cosmosPackets[0]
		bb.cosmosMtx.Unlock()
	}
	if p.Retries != 0 {
		t.Fatalf("transient re-queue must not increment Retries, got %d", p.Retries)
	}
	if c, _ := bb.DeadLetterCounts(); c != 0 {
		t.Fatalf("transient re-queue must never dead-letter, got %d", c)
	}
}

func TestRequeueCosmosPermanent_PrependAndIncrement(t *testing.T) {
	bb := NewBatchBuilder()
	bb.AddCosmos(makeCosmosPacket(10, CosmosSend)) // newer arrival

	bb.RequeueCosmosPermanent([]CosmosPacket{
		makeCosmosPacket(1, CosmosSend),
		makeCosmosPacket(2, CosmosSend),
	})

	bb.cosmosMtx.Lock()
	defer bb.cosmosMtx.Unlock()
	if len(bb.cosmosPackets) != 3 {
		t.Fatalf("expected 3 packets, got %d", len(bb.cosmosPackets))
	}
	wantSeq := []uint64{1, 2, 10} // re-queued prepend, in order, ahead of seq=10
	for i, want := range wantSeq {
		if got := bb.cosmosPackets[i].Packet.Sequence; got != want {
			t.Fatalf("position %d: want seq %d, got %d", i, want, got)
		}
	}
	if bb.cosmosPackets[0].Retries != 1 || bb.cosmosPackets[1].Retries != 1 {
		t.Fatalf("permanent re-queue should increment to Retries=1")
	}
	if bb.cosmosPackets[2].Retries != 0 {
		t.Fatalf("untouched packet should keep Retries=0")
	}
}

func TestRequeueCosmosPermanent_DeadLetterAtCap(t *testing.T) {
	bb := NewBatchBuilder()
	atCap := makeCosmosPacket(1, CosmosSend)
	atCap.Retries = maxPacketRetries // next permanent requeue pushes over the cap
	belowCap := makeCosmosPacket(2, CosmosSend)

	bb.RequeueCosmosPermanent([]CosmosPacket{atCap, belowCap})

	bb.cosmosMtx.Lock()
	defer bb.cosmosMtx.Unlock()
	if len(bb.cosmosPackets) != 1 || bb.cosmosPackets[0].Packet.Sequence != 2 {
		t.Fatalf("expected only seq=2 kept, got %+v", bb.cosmosPackets)
	}
	if len(bb.deadLetterCosmos) != 1 || bb.deadLetterCosmos[0].Packet.Sequence != 1 {
		t.Fatalf("expected seq=1 dead-lettered, got %+v", bb.deadLetterCosmos)
	}
}

func TestRequeueEthTransient_NoBudgetConsumed(t *testing.T) {
	bb := NewBatchBuilder()
	p := makeEthPacket(1, EthSend)

	for range maxPacketRetries * 3 {
		bb.ethMtx.Lock()
		bb.ethPackets = nil
		bb.ethMtx.Unlock()
		bb.RequeueEthTransient([]EthPacket{p})
		bb.ethMtx.Lock()
		p = bb.ethPackets[0]
		bb.ethMtx.Unlock()
	}
	if p.Retries != 0 {
		t.Fatalf("transient re-queue must not increment Retries, got %d", p.Retries)
	}
	if _, e := bb.DeadLetterCounts(); e != 0 {
		t.Fatalf("transient re-queue must never dead-letter, got %d", e)
	}
}

func TestRequeueEthPermanent_DeadLetterAtCap(t *testing.T) {
	bb := NewBatchBuilder()
	atCap := makeEthPacket(1, EthSend)
	atCap.Retries = maxPacketRetries

	bb.RequeueEthPermanent([]EthPacket{atCap})

	bb.ethMtx.Lock()
	defer bb.ethMtx.Unlock()
	if len(bb.ethPackets) != 0 {
		t.Fatalf("expected at-cap packet removed from queue, got %d", len(bb.ethPackets))
	}
	if len(bb.deadLetterEth) != 1 {
		t.Fatalf("expected 1 dead-lettered eth packet, got %d", len(bb.deadLetterEth))
	}
}

func TestRequeue_Empty(t *testing.T) {
	bb := NewBatchBuilder()
	bb.RequeueCosmosTransient(nil)
	bb.RequeueCosmosPermanent(nil)
	bb.RequeueEthTransient([]EthPacket{})
	bb.RequeueEthPermanent([]EthPacket{})
	if len(bb.cosmosPackets) != 0 || len(bb.ethPackets) != 0 {
		t.Fatal("empty requeue should be a no-op")
	}
}

// TestRequeueCosmosPermanent_AccumulateThenDeadLetter verifies a packet
// survives exactly maxPacketRetries permanent failures, then is dead-lettered.
func TestRequeueCosmosPermanent_AccumulateThenDeadLetter(t *testing.T) {
	bb := NewBatchBuilder()
	p := makeCosmosPacket(1, CosmosSend)

	for i := range maxPacketRetries {
		bb.cosmosMtx.Lock()
		bb.cosmosPackets = nil
		bb.cosmosMtx.Unlock()
		bb.RequeueCosmosPermanent([]CosmosPacket{p})
		bb.cosmosMtx.Lock()
		if len(bb.cosmosPackets) != 1 {
			bb.cosmosMtx.Unlock()
			t.Fatalf("cycle %d: packet should still be queued (Retries=%d <= cap)", i, p.Retries+1)
		}
		p = bb.cosmosPackets[0]
		bb.cosmosMtx.Unlock()
	}

	// p.Retries == maxPacketRetries now; one more permanent failure dead-letters.
	bb.cosmosMtx.Lock()
	bb.cosmosPackets = nil
	bb.cosmosMtx.Unlock()
	bb.RequeueCosmosPermanent([]CosmosPacket{p})

	if c, _ := bb.DeadLetterCounts(); c != 1 {
		t.Fatalf("expected 1 dead-lettered after exceeding cap, got %d", c)
	}
}
