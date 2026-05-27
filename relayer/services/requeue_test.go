package services

import "testing"

// Issue #80: packets whose batch submission fails must be re-queued (not
// dropped), prepended ahead of newer arrivals, and dropped once they exceed
// maxPacketRetries.

func TestRequeueCosmos_PrependAndIncrement(t *testing.T) {
	bb := NewBatchBuilder()
	bb.AddCosmos(makeCosmosPacket(10, CosmosSend)) // already queued, newer arrival

	bb.RequeueCosmos([]CosmosPacket{
		makeCosmosPacket(1, CosmosSend),
		makeCosmosPacket(2, CosmosSend),
	})

	bb.cosmosMtx.Lock()
	defer bb.cosmosMtx.Unlock()
	if len(bb.cosmosPackets) != 3 {
		t.Fatalf("expected 3 packets, got %d", len(bb.cosmosPackets))
	}
	// Re-queued packets prepend, preserving their order, ahead of seq=10.
	wantSeq := []uint64{1, 2, 10}
	for i, want := range wantSeq {
		if got := bb.cosmosPackets[i].Packet.Sequence; got != want {
			t.Fatalf("position %d: want seq %d, got %d", i, want, got)
		}
	}
	if bb.cosmosPackets[0].Retries != 1 || bb.cosmosPackets[1].Retries != 1 {
		t.Fatalf("re-queued packets should have Retries=1, got %d/%d",
			bb.cosmosPackets[0].Retries, bb.cosmosPackets[1].Retries)
	}
	if bb.cosmosPackets[2].Retries != 0 {
		t.Fatalf("untouched packet should keep Retries=0, got %d", bb.cosmosPackets[2].Retries)
	}
}

func TestRequeueCosmos_DropAtCap(t *testing.T) {
	bb := NewBatchBuilder()
	atCap := makeCosmosPacket(1, CosmosSend)
	atCap.Retries = maxPacketRetries // next requeue pushes it over the cap
	belowCap := makeCosmosPacket(2, CosmosSend)

	bb.RequeueCosmos([]CosmosPacket{atCap, belowCap})

	bb.cosmosMtx.Lock()
	defer bb.cosmosMtx.Unlock()
	if len(bb.cosmosPackets) != 1 {
		t.Fatalf("expected 1 packet kept (the at-cap one dropped), got %d", len(bb.cosmosPackets))
	}
	if bb.cosmosPackets[0].Packet.Sequence != 2 {
		t.Fatalf("expected seq=2 kept, got %d", bb.cosmosPackets[0].Packet.Sequence)
	}
}

func TestRequeueEth_PrependAndIncrement(t *testing.T) {
	bb := NewBatchBuilder()
	bb.AddEth(makeEthPacket(10, EthSend))

	bb.RequeueEth([]EthPacket{
		makeEthPacket(1, EthSend),
		makeEthPacket(2, EthWriteAck),
	})

	bb.ethMtx.Lock()
	defer bb.ethMtx.Unlock()
	if len(bb.ethPackets) != 3 {
		t.Fatalf("expected 3 packets, got %d", len(bb.ethPackets))
	}
	wantSeq := []uint64{1, 2, 10}
	for i, want := range wantSeq {
		if got := bb.ethPackets[i].Packet.Sequence; got != want {
			t.Fatalf("position %d: want seq %d, got %d", i, want, got)
		}
	}
	if bb.ethPackets[0].Retries != 1 || bb.ethPackets[1].Retries != 1 {
		t.Fatalf("re-queued packets should have Retries=1")
	}
}

func TestRequeueEth_DropAtCap(t *testing.T) {
	bb := NewBatchBuilder()
	atCap := makeEthPacket(1, EthSend)
	atCap.Retries = maxPacketRetries

	bb.RequeueEth([]EthPacket{atCap})

	bb.ethMtx.Lock()
	defer bb.ethMtx.Unlock()
	if len(bb.ethPackets) != 0 {
		t.Fatalf("expected at-cap packet dropped, got %d packets", len(bb.ethPackets))
	}
}

func TestRequeue_Empty(t *testing.T) {
	bb := NewBatchBuilder()
	bb.RequeueCosmos(nil)
	bb.RequeueEth([]EthPacket{})
	if len(bb.cosmosPackets) != 0 || len(bb.ethPackets) != 0 {
		t.Fatal("empty requeue should be a no-op")
	}
}

// TestRequeueCosmos_RetriesAccumulate verifies a packet dropped only after
// maxPacketRetries successive failures, not before.
func TestRequeueCosmos_RetriesAccumulate(t *testing.T) {
	bb := NewBatchBuilder()
	p := makeCosmosPacket(1, CosmosSend)

	// Simulate the relay loop: each failed flush pulls the packet, fails, and
	// re-queues it. It should survive maxPacketRetries cycles then be dropped.
	for i := range maxPacketRetries {
		bb.cosmosMtx.Lock()
		got := len(bb.cosmosPackets)
		bb.cosmosMtx.Unlock()
		if got != 1 && i > 0 {
			t.Fatalf("cycle %d: expected packet still queued, got %d", i, got)
		}

		bb.cosmosMtx.Lock()
		bb.cosmosPackets = nil // simulate CheckCosmos slicing it off
		bb.cosmosMtx.Unlock()

		bb.RequeueCosmos([]CosmosPacket{p})

		bb.cosmosMtx.Lock()
		p = bb.cosmosPackets[0] // carry forward the incremented retry count
		bb.cosmosMtx.Unlock()
	}

	// p.Retries is now maxPacketRetries; one more requeue drops it.
	bb.cosmosMtx.Lock()
	bb.cosmosPackets = nil
	bb.cosmosMtx.Unlock()
	bb.RequeueCosmos([]CosmosPacket{p})

	bb.cosmosMtx.Lock()
	defer bb.cosmosMtx.Unlock()
	if len(bb.cosmosPackets) != 0 {
		t.Fatalf("packet should be dropped after exceeding maxPacketRetries, got %d", len(bb.cosmosPackets))
	}
}
