package services

import (
	"errors"
	"sync"
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func makeCosmosPacket(seq uint64, pktType CosmosPacketType) CosmosPacket {
	return CosmosPacket{
		Type: pktType,
		Packet: &channeltypesv2.Packet{
			Sequence: seq,
		},
	}
}

func makeEthPacket(seq uint64, pktType EthPacketType) EthPacket {
	return EthPacket{
		Type: pktType,
		Packet: &channeltypesv2.Packet{
			Sequence: seq,
		},
	}
}

func TestNewBatchBuilder(t *testing.T) {
	bb := NewBatchBuilder()
	if bb == nil {
		t.Fatal("expected non-nil BatchBuilder")
	}
	if len(bb.cosmosPackets) != 0 {
		t.Fatalf("expected empty cosmos packets, got %d", len(bb.cosmosPackets))
	}
	if len(bb.ethPackets) != 0 {
		t.Fatalf("expected empty eth packets, got %d", len(bb.ethPackets))
	}
	if bb.cosmosTimestamp.IsZero() {
		t.Fatal("expected cosmos timestamp to be set")
	}
	if bb.ethTimestamp.IsZero() {
		t.Fatal("expected eth timestamp to be set")
	}
}

func TestAddCosmos(t *testing.T) {
	bb := NewBatchBuilder()
	bb.AddCosmos(makeCosmosPacket(1, CosmosSend))

	bb.cosmosMtx.Lock()
	defer bb.cosmosMtx.Unlock()

	if len(bb.cosmosPackets) != 1 {
		t.Fatalf("expected 1 cosmos packet, got %d", len(bb.cosmosPackets))
	}
	if bb.cosmosPackets[0].Packet.Sequence != 1 {
		t.Fatalf("expected sequence 1, got %d", bb.cosmosPackets[0].Packet.Sequence)
	}
	if bb.cosmosPackets[0].Type != CosmosSend {
		t.Fatalf("expected CosmosSend, got %d", bb.cosmosPackets[0].Type)
	}
}

func TestAddEth(t *testing.T) {
	bb := NewBatchBuilder()
	bb.AddEth(makeEthPacket(1, EthSend))

	bb.ethMtx.Lock()
	defer bb.ethMtx.Unlock()

	if len(bb.ethPackets) != 1 {
		t.Fatalf("expected 1 eth packet, got %d", len(bb.ethPackets))
	}
	if bb.ethPackets[0].Packet.Sequence != 1 {
		t.Fatalf("expected sequence 1, got %d", bb.ethPackets[0].Packet.Sequence)
	}
	if bb.ethPackets[0].Type != EthSend {
		t.Fatalf("expected EthSend, got %d", bb.ethPackets[0].Type)
	}
}

func TestShouldTimeoutEthSend(t *testing.T) {
	expired := makeEthPacket(1, EthSend)
	expired.Packet.TimeoutTimestamp = uint64(time.Now().Add(-time.Second).Unix())

	if !shouldTimeoutEthSend(expired, errors.New("receive packet verification failed: timeout elapsed")) {
		t.Fatal("expected timeout fallback for expired packet with timeout error")
	}

	active := makeEthPacket(2, EthSend)
	active.Packet.TimeoutTimestamp = uint64(time.Now().Add(time.Hour).Unix())
	if shouldTimeoutEthSend(active, errors.New("receive packet verification failed: timeout elapsed")) {
		t.Fatal("did not expect timeout fallback before timeout timestamp")
	}

	if shouldTimeoutEthSend(expired, errors.New("unrelated submit failure")) {
		t.Fatal("did not expect timeout fallback for unrelated errors")
	}
}

func TestAddCosmosMultiple(t *testing.T) {
	bb := NewBatchBuilder()
	count := 5
	for i := 0; i < count; i++ {
		bb.AddCosmos(makeCosmosPacket(uint64(i+1), CosmosSend))
	}

	bb.cosmosMtx.Lock()
	defer bb.cosmosMtx.Unlock()

	if len(bb.cosmosPackets) != count {
		t.Fatalf("expected %d cosmos packets, got %d", count, len(bb.cosmosPackets))
	}
	for i := 0; i < count; i++ {
		if bb.cosmosPackets[i].Packet.Sequence != uint64(i+1) {
			t.Fatalf("packet %d: expected sequence %d, got %d", i, i+1, bb.cosmosPackets[i].Packet.Sequence)
		}
	}
}

func TestAddEthMultiple(t *testing.T) {
	bb := NewBatchBuilder()
	count := 5
	for i := 0; i < count; i++ {
		bb.AddEth(makeEthPacket(uint64(i+1), EthSend))
	}

	bb.ethMtx.Lock()
	defer bb.ethMtx.Unlock()

	if len(bb.ethPackets) != count {
		t.Fatalf("expected %d eth packets, got %d", count, len(bb.ethPackets))
	}
	for i := 0; i < count; i++ {
		if bb.ethPackets[i].Packet.Sequence != uint64(i+1) {
			t.Fatalf("packet %d: expected sequence %d, got %d", i, i+1, bb.ethPackets[i].Packet.Sequence)
		}
	}
}

func TestAddCosmosConcurrent(t *testing.T) {
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
				bb.AddCosmos(makeCosmosPacket(uint64(offset*packetsPerGoroutine+i), CosmosSend))
			}
		}(g)
	}

	wg.Wait()

	bb.cosmosMtx.Lock()
	defer bb.cosmosMtx.Unlock()

	if len(bb.cosmosPackets) != total {
		t.Fatalf("expected %d cosmos packets, got %d", total, len(bb.cosmosPackets))
	}
}

func TestAddEthConcurrent(t *testing.T) {
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
				bb.AddEth(makeEthPacket(uint64(offset*packetsPerGoroutine+i), EthSend))
			}
		}(g)
	}

	wg.Wait()

	bb.ethMtx.Lock()
	defer bb.ethMtx.Unlock()

	if len(bb.ethPackets) != total {
		t.Fatalf("expected %d eth packets, got %d", total, len(bb.ethPackets))
	}
}

func TestClearCosmos(t *testing.T) {
	bb := NewBatchBuilder()
	for i := 0; i < 5; i++ {
		bb.AddCosmos(makeCosmosPacket(uint64(i+1), CosmosSend))
	}

	bb.ClearCosmos()

	bb.cosmosMtx.Lock()
	defer bb.cosmosMtx.Unlock()

	if len(bb.cosmosPackets) != 0 {
		t.Fatalf("expected 0 cosmos packets after clear, got %d", len(bb.cosmosPackets))
	}
}

func TestClearEth(t *testing.T) {
	bb := NewBatchBuilder()
	for i := 0; i < 5; i++ {
		bb.AddEth(makeEthPacket(uint64(i+1), EthSend))
	}

	bb.ClearEth()

	bb.ethMtx.Lock()
	defer bb.ethMtx.Unlock()

	if len(bb.ethPackets) != 0 {
		t.Fatalf("expected 0 eth packets after clear, got %d", len(bb.ethPackets))
	}
}

func TestCheckCosmos_ExceedsBatchSize(t *testing.T) {
	bb := NewBatchBuilder()
	config := BatchConfig{
		BatchSize:    3,
		BatchPeriods: time.Minute * 10,
	}

	for i := 0; i < 5; i++ {
		bb.AddCosmos(makeCosmosPacket(uint64(i+1), CosmosSend))
	}

	ch := make(chan CosmosBatch, 1)
	bb.CheckCosmos(config, ch)

	select {
	case batch := <-ch:
		if len(batch.Packets) != 5 {
			t.Fatalf("expected 5 packets in batch, got %d", len(batch.Packets))
		}
	default:
		t.Fatal("expected cosmos batch to be sent")
	}

	bb.cosmosMtx.Lock()
	defer bb.cosmosMtx.Unlock()
	if len(bb.cosmosPackets) != 0 {
		t.Fatalf("expected cosmos packets to be cleared, got %d", len(bb.cosmosPackets))
	}
}

func TestCheckEth_ExceedsBatchSize(t *testing.T) {
	bb := NewBatchBuilder()
	config := BatchConfig{
		BatchSize:    3,
		BatchPeriods: time.Minute * 10,
	}

	for i := 0; i < 5; i++ {
		bb.AddEth(makeEthPacket(uint64(i+1), EthSend))
	}

	ch := make(chan EthBatch, 1)
	bb.CheckEth(config, ch)

	select {
	case batch := <-ch:
		if len(batch.Packets) != 5 {
			t.Fatalf("expected 5 packets in batch, got %d", len(batch.Packets))
		}
	default:
		t.Fatal("expected eth batch to be sent")
	}

	bb.ethMtx.Lock()
	defer bb.ethMtx.Unlock()
	if len(bb.ethPackets) != 0 {
		t.Fatalf("expected eth packets to be cleared, got %d", len(bb.ethPackets))
	}
}

func TestCheckEth_PastBatchPeriods(t *testing.T) {
	bb := NewBatchBuilder()
	bb.ethTimestamp = time.Now().Add(-time.Minute * 5)

	config := BatchConfig{
		BatchSize:    10,
		BatchPeriods: time.Second,
	}

	bb.AddEth(makeEthPacket(1, EthSend))
	bb.AddEth(makeEthPacket(2, EthWriteAck))

	ch := make(chan EthBatch, 1)
	bb.CheckEth(config, ch)

	select {
	case batch := <-ch:
		if len(batch.Packets) != 2 {
			t.Fatalf("expected 2 packets in batch, got %d", len(batch.Packets))
		}
	default:
		t.Fatal("expected eth batch to be sent")
	}

	bb.ethMtx.Lock()
	defer bb.ethMtx.Unlock()
	if len(bb.ethPackets) != 0 {
		t.Fatalf("expected eth packets to be cleared, got %d", len(bb.ethPackets))
	}
}

func TestCheckEth_BelowSizeAndBeforePeriod(t *testing.T) {
	bb := NewBatchBuilder()
	config := BatchConfig{
		BatchSize:    10,
		BatchPeriods: time.Minute * 10,
	}

	bb.AddEth(makeEthPacket(1, EthSend))
	bb.AddEth(makeEthPacket(2, EthWriteAck))

	ch := make(chan EthBatch, 1)
	bb.CheckEth(config, ch)

	select {
	case <-ch:
		t.Fatal("expected no eth batch to be sent")
	default:
	}

	bb.ethMtx.Lock()
	defer bb.ethMtx.Unlock()
	if len(bb.ethPackets) != 2 {
		t.Fatalf("expected eth packets to remain, got %d", len(bb.ethPackets))
	}
}

func TestCheckCosmos_BelowSizeAndBeforePeriod(t *testing.T) {
	bb := NewBatchBuilder()
	config := BatchConfig{
		BatchSize:    10,
		BatchPeriods: time.Minute * 10,
	}

	bb.AddCosmos(makeCosmosPacket(1, CosmosSend))
	bb.AddCosmos(makeCosmosPacket(2, CosmosAck))

	ch := make(chan CosmosBatch, 1)
	bb.CheckCosmos(config, ch)

	select {
	case <-ch:
		t.Fatal("expected no cosmos batch to be sent")
	default:
	}

	bb.cosmosMtx.Lock()
	defer bb.cosmosMtx.Unlock()
	if len(bb.cosmosPackets) != 2 {
		t.Fatalf("expected cosmos packets to remain, got %d", len(bb.cosmosPackets))
	}
}
