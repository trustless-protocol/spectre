package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func (p CosmosPacketType) String() string {
	switch p {
	case CosmosSend:
		return "Send"
	case CosmosAck:
		return "Ack"
	case CosmosTimeout:
		return "Timeout"
	default:
		return fmt.Sprintf("Unknown(%d)", int(p))
	}
}

func (p EthPacketType) String() string {
	switch p {
	case EthSend:
		return "Send"
	case EthWriteAck:
		return "WriteAck"
	case EthAck:
		return "Ack"
	case EthTimeout:
		return "Timeout"
	default:
		return fmt.Sprintf("Unknown(%d)", int(p))
	}
}

type CosmosPacketType int

const (
	CosmosSend CosmosPacketType = iota
	CosmosAck
	CosmosTimeout
)

type EthPacketType int

const (
	EthSend EthPacketType = iota
	EthWriteAck
	EthAck
	EthTimeout
)

type CosmosPacket struct {
	Type        CosmosPacketType
	Packet      *channeltypesv2.Packet
	AckBytes    [][]byte
	BlockNumber uint64
	// Retries counts how many times this packet's batch failed to submit and
	// was re-queued. Bounded by maxPacketRetries so a permanently-broken packet
	// (corrupt proof, etc.) can't starve the queue forever (issue #80).
	Retries int
}

type EthPacket struct {
	Type        EthPacketType
	Packet      *channeltypesv2.Packet
	AckBytes    [][]byte
	BlockNumber uint64
	// Retries — see CosmosPacket.Retries.
	Retries int
}

type CosmosBatch struct {
	Packets []CosmosPacket
}

type EthBatch struct {
	Packets []EthPacket
}

type BatchBuilder struct {
	cosmosMtx         sync.Mutex
	ethMtx            sync.Mutex
	cosmosTimestamp   time.Time
	ethTimestamp      time.Time
	cosmosPackets     []CosmosPacket
	ethPackets        []EthPacket
	PendingTracker    *PendingPacketTracker
	EthPendingTracker *PendingPacketTracker

	// Dead-lettered packets: those that hit maxPacketRetries on PERMANENT
	// (deterministic revert) failures. Kept rather than silently dropped so
	// they can be inspected / surfaced as metrics later (issue #80 review).
	// Guarded by the matching direction mutex.
	deadLetterCosmos []CosmosPacket
	deadLetterEth    []EthPacket
}

// DeadLetterCounts returns how many packets have been dead-lettered per
// direction. Used for observability / alerting.
func (b *BatchBuilder) DeadLetterCounts() (cosmos, eth int) {
	b.cosmosMtx.Lock()
	cosmos = len(b.deadLetterCosmos)
	b.cosmosMtx.Unlock()
	b.ethMtx.Lock()
	eth = len(b.deadLetterEth)
	b.ethMtx.Unlock()
	return cosmos, eth
}

func NewBatchBuilder() *BatchBuilder {
	now := time.Now()
	return &BatchBuilder{
		cosmosTimestamp:   now,
		ethTimestamp:      now,
		cosmosPackets:     []CosmosPacket{},
		ethPackets:        []EthPacket{},
		PendingTracker:    NewPendingPacketTracker(),
		EthPendingTracker: NewPendingPacketTracker(),
	}
}

func (b *BatchBuilder) AddCosmos(packet CosmosPacket) {
	b.cosmosMtx.Lock()
	b.cosmosPackets = append(b.cosmosPackets, packet)
	count := len(b.cosmosPackets)
	b.cosmosMtx.Unlock()
	log.Printf("[BatchBuilder] Inserted cosmos packet: type=%s seq=%d (batch size: %d)",
		packet.Type, packet.Packet.Sequence, count)
}

func (b *BatchBuilder) AddEth(packet EthPacket) {
	b.ethMtx.Lock()
	b.ethPackets = append(b.ethPackets, packet)
	count := len(b.ethPackets)
	b.ethMtx.Unlock()
	log.Printf("[BatchBuilder] Inserted eth packet: type=%s seq=%d (batch size: %d)",
		packet.Type, packet.Packet.Sequence, count)
}

func (b *BatchBuilder) ClearCosmos() {
	b.cosmosTimestamp = time.Now()
	b.cosmosPackets = []CosmosPacket{}
}

func (b *BatchBuilder) ClearEth() {
	b.ethTimestamp = time.Now()
	b.ethPackets = []EthPacket{}
}

// maxPacketRetries caps how many times a packet that failed with a PERMANENT
// (deterministic on-chain revert) error is re-queued before being dead-lettered.
// Only permanent failures consume this budget — transient infrastructure
// failures (RPC, beacon finality, build, broadcast timeout) re-queue without
// counting, so a valid packet is never lost just because the infra was briefly
// down (issue #80 review). The cap exists purely to stop a genuine poison packet
// (corrupt proof, perpetually-reverting inner call) from blocking the queue.
const maxPacketRetries = 5

// RequeueCosmosTransient re-queues after an infrastructure failure without
// touching the retry budget: the packet is valid and will succeed once the
// transient cause clears.
func (b *BatchBuilder) RequeueCosmosTransient(packets []CosmosPacket) {
	b.requeueCosmos(packets, false)
}

// RequeueCosmosPermanent re-queues after a deterministic on-chain revert,
// consuming the retry budget; once it exceeds maxPacketRetries the packet is
// dead-lettered (kept for inspection + a loud log) rather than silently dropped.
func (b *BatchBuilder) RequeueCosmosPermanent(packets []CosmosPacket) {
	b.requeueCosmos(packets, true)
}

func (b *BatchBuilder) requeueCosmos(packets []CosmosPacket, permanent bool) {
	if len(packets) == 0 {
		return
	}
	kept := make([]CosmosPacket, 0, len(packets))
	var dead []CosmosPacket
	for _, p := range packets {
		if permanent {
			p.Retries++
			if p.Retries > maxPacketRetries {
				log.Printf("[BatchBuilder][DEAD-LETTER] cosmos packet type=%s seq=%d dead-lettered after %d permanent failures",
					p.Type, p.Packet.Sequence, maxPacketRetries)
				dead = append(dead, p)
				continue
			}
		}
		kept = append(kept, p)
	}
	b.cosmosMtx.Lock()
	if len(dead) > 0 {
		b.deadLetterCosmos = append(b.deadLetterCosmos, dead...)
	}
	if len(kept) > 0 {
		b.cosmosPackets = append(kept, b.cosmosPackets...)
	}
	remaining := len(b.cosmosPackets)
	b.cosmosMtx.Unlock()
	if len(kept) > 0 {
		kind := "transient"
		if permanent {
			kind = "permanent"
		}
		log.Printf("[BatchBuilder] re-queued %d cosmos packet(s) (%s) for retry (queue now %d)", len(kept), kind, remaining)
	}
}

// RequeueEthTransient — see RequeueCosmosTransient.
func (b *BatchBuilder) RequeueEthTransient(packets []EthPacket) {
	b.requeueEth(packets, false)
}

// RequeueEthPermanent — see RequeueCosmosPermanent.
func (b *BatchBuilder) RequeueEthPermanent(packets []EthPacket) {
	b.requeueEth(packets, true)
}

func (b *BatchBuilder) requeueEth(packets []EthPacket, permanent bool) {
	if len(packets) == 0 {
		return
	}
	kept := make([]EthPacket, 0, len(packets))
	var dead []EthPacket
	for _, p := range packets {
		if permanent {
			p.Retries++
			if p.Retries > maxPacketRetries {
				log.Printf("[BatchBuilder][DEAD-LETTER] eth packet type=%s seq=%d dead-lettered after %d permanent failures",
					p.Type, p.Packet.Sequence, maxPacketRetries)
				dead = append(dead, p)
				continue
			}
		}
		kept = append(kept, p)
	}
	b.ethMtx.Lock()
	if len(dead) > 0 {
		b.deadLetterEth = append(b.deadLetterEth, dead...)
	}
	if len(kept) > 0 {
		b.ethPackets = append(kept, b.ethPackets...)
	}
	remaining := len(b.ethPackets)
	b.ethMtx.Unlock()
	if len(kept) > 0 {
		kind := "transient"
		if permanent {
			kind = "permanent"
		}
		log.Printf("[BatchBuilder] re-queued %d eth packet(s) (%s) for retry (queue now %d)", len(kept), kind, remaining)
	}
}

func (b *BatchBuilder) CheckCosmos(config BatchConfig, ch chan<- CosmosBatch) {
	b.cosmosMtx.Lock()

	if len(b.cosmosPackets) == 0 {
		b.cosmosMtx.Unlock()
		return
	}

	flush := false
	reason := ""

	if len(b.cosmosPackets) >= int(config.BatchSize) {
		flush = true
		reason = fmt.Sprintf("size limit reached (%d >= %d)", len(b.cosmosPackets), config.BatchSize)
	} else if time.Now().After(b.cosmosTimestamp.Add(config.BatchPeriods)) {
		flush = true
		reason = fmt.Sprintf("time limit reached (%v elapsed)", time.Since(b.cosmosTimestamp).Round(time.Millisecond))
	}

	if !flush {
		b.cosmosMtx.Unlock()
		return
	}

	// Slice up to BatchSize so a backlog (queue >> BatchSize) is flushed in
	// fixed-size chunks instead of one oversized multicall that would blow
	// past the EVM 128KB tx size limit.
	chunkSize := int(config.BatchSize)
	if chunkSize <= 0 || chunkSize > len(b.cosmosPackets) {
		chunkSize = len(b.cosmosPackets)
	}
	chunk := b.cosmosPackets[:chunkSize]
	b.cosmosPackets = b.cosmosPackets[chunkSize:]
	b.cosmosTimestamp = time.Now()
	log.Printf("[BatchBuilder] Flushing cosmos batch: %d packets (%s, queue remaining: %d)",
		len(chunk), reason, len(b.cosmosPackets))
	batch := CosmosBatch{Packets: chunk}
	b.cosmosMtx.Unlock()

	ch <- batch
}

func (b *BatchBuilder) CheckEth(config BatchConfig, ch chan<- EthBatch) {
	b.ethMtx.Lock()

	if len(b.ethPackets) == 0 {
		b.ethMtx.Unlock()
		return
	}

	flush := false
	reason := ""

	if len(b.ethPackets) >= int(config.BatchSize) {
		flush = true
		reason = fmt.Sprintf("size limit reached (%d >= %d)", len(b.ethPackets), config.BatchSize)
	} else if time.Now().After(b.ethTimestamp.Add(config.BatchPeriods)) {
		flush = true
		reason = fmt.Sprintf("time limit reached (%v elapsed)", time.Since(b.ethTimestamp).Round(time.Millisecond))
	}

	if !flush {
		b.ethMtx.Unlock()
		return
	}

	chunkSize := int(config.BatchSize)
	if chunkSize <= 0 || chunkSize > len(b.ethPackets) {
		chunkSize = len(b.ethPackets)
	}
	chunk := b.ethPackets[:chunkSize]
	b.ethPackets = b.ethPackets[chunkSize:]
	b.ethTimestamp = time.Now()
	log.Printf("[BatchBuilder] Flushing eth batch: %d packets (%s, queue remaining: %d)",
		len(chunk), reason, len(b.ethPackets))
	batch := EthBatch{Packets: chunk}
	b.ethMtx.Unlock()

	ch <- batch
}
