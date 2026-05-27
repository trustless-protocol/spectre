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
	cosmosMtx       sync.Mutex
	ethMtx          sync.Mutex
	cosmosTimestamp time.Time
	ethTimestamp    time.Time
	cosmosPackets   []CosmosPacket
	ethPackets      []EthPacket
	PendingTracker  *PendingPacketTracker
}

func NewBatchBuilder() *BatchBuilder {
	now := time.Now()
	return &BatchBuilder{
		cosmosTimestamp: now,
		ethTimestamp:    now,
		cosmosPackets:   []CosmosPacket{},
		ethPackets:      []EthPacket{},
		PendingTracker:  NewPendingPacketTracker(),
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

// maxPacketRetries caps how many times a packet whose batch failed to submit is
// re-queued before being dropped. Without a cap, a permanently-invalid packet
// (corrupt proof, perpetually-reverting inner call) prepended to the queue would
// block every later packet forever. 5 covers transient faults (RPC blips,
// insufficient sync-committee participation that recovers within a few epochs)
// while still giving up on a genuine poison packet.
const maxPacketRetries = 5

// RequeueCosmos prepends packets back onto the Cosmos queue after a failed batch
// submission so the next flush retries them, preserving their order relative to
// each other and ahead of newer arrivals (issue #80). Packets that have already
// failed maxPacketRetries times are dropped with a log line instead of looping
// forever.
func (b *BatchBuilder) RequeueCosmos(packets []CosmosPacket) {
	if len(packets) == 0 {
		return
	}
	kept := make([]CosmosPacket, 0, len(packets))
	for _, p := range packets {
		p.Retries++
		if p.Retries > maxPacketRetries {
			log.Printf("[BatchBuilder] dropping cosmos packet type=%s seq=%d after %d failed attempts",
				p.Type, p.Packet.Sequence, maxPacketRetries)
			continue
		}
		kept = append(kept, p)
	}
	if len(kept) == 0 {
		return
	}
	b.cosmosMtx.Lock()
	b.cosmosPackets = append(kept, b.cosmosPackets...)
	remaining := len(b.cosmosPackets)
	b.cosmosMtx.Unlock()
	log.Printf("[BatchBuilder] re-queued %d cosmos packet(s) for retry (queue now %d)", len(kept), remaining)
}

// RequeueEth — see RequeueCosmos. Mirror for the ETH→Cosmos direction.
func (b *BatchBuilder) RequeueEth(packets []EthPacket) {
	if len(packets) == 0 {
		return
	}
	kept := make([]EthPacket, 0, len(packets))
	for _, p := range packets {
		p.Retries++
		if p.Retries > maxPacketRetries {
			log.Printf("[BatchBuilder] dropping eth packet type=%s seq=%d after %d failed attempts",
				p.Type, p.Packet.Sequence, maxPacketRetries)
			continue
		}
		kept = append(kept, p)
	}
	if len(kept) == 0 {
		return
	}
	b.ethMtx.Lock()
	b.ethPackets = append(kept, b.ethPackets...)
	remaining := len(b.ethPackets)
	b.ethMtx.Unlock()
	log.Printf("[BatchBuilder] re-queued %d eth packet(s) for retry (queue now %d)", len(kept), remaining)
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
