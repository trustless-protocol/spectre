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
	// NotBefore, when set, holds the packet out of Check* flushes until the time
	// passes — a re-queue backoff for a packet that is not yet relayable (source
	// state lag / finality) so the loop polls it less aggressively than the batch
	// period. Zero (the legacy default) means "ready now". waitAttempts drives the
	// exponential growth of that backoff.
	NotBefore    time.Time
	waitAttempts int
}

type EthPacket struct {
	Type        EthPacketType
	Packet      *channeltypesv2.Packet
	AckBytes    [][]byte
	BlockNumber uint64
	// NotBefore / waitAttempts — see CosmosPacket.
	NotBefore    time.Time
	waitAttempts int
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

// Waiting-backoff bounds: a packet re-queued because it is not yet relayable
// (source state / finality lag) is held out of flushes with an exponentially
// growing delay, so a quick precondition (Cosmos AppHash H+2, ~1 block) clears
// fast while a long one (ETH beacon finality, minutes) polls sparsely instead of
// every batch period. Bounded so latency after the precondition clears stays low.
const (
	waitingBackoffBase = 3 * time.Second
	waitingBackoffMax  = 15 * time.Second
)

// nextWaitingBackoff returns the delay for the given (already-incremented)
// attempt count: base, 2×base, 4×base, … capped at max.
func nextWaitingBackoff(attempt int) time.Duration {
	d := waitingBackoffBase << (attempt - 1)
	if d > waitingBackoffMax || d <= 0 { // <=0 guards shift overflow
		d = waitingBackoffMax
	}
	return d
}

// RequeueCosmosWaiting re-queues a packet that is valid but NOT YET relayable
// (its source state is not yet provable). It sets a growing NotBefore backoff so
// the loop does not re-attempt (and re-log, and re-hit the RPC) every batch
// period while waiting, then prepends the packet ahead of newer arrivals so a
// re-queued packet is retried before newer work. This is the only re-queue path:
// the adapter module drops permanently-failed packets at the source (the timeout
// scanner refunds them), so there is no retry-budget / dead-letter classification.
func (b *BatchBuilder) RequeueCosmosWaiting(packets []CosmosPacket) {
	if len(packets) == 0 {
		return
	}
	now := time.Now()
	for i := range packets {
		packets[i].waitAttempts++
		packets[i].NotBefore = now.Add(nextWaitingBackoff(packets[i].waitAttempts))
	}
	b.cosmosMtx.Lock()
	b.cosmosPackets = append(packets, b.cosmosPackets...)
	remaining := len(b.cosmosPackets)
	b.cosmosMtx.Unlock()
	log.Printf("[BatchBuilder] re-queued %d cosmos packet(s) (waiting) for retry (queue now %d)", len(packets), remaining)
}

// RequeueEthWaiting — see RequeueCosmosWaiting.
func (b *BatchBuilder) RequeueEthWaiting(packets []EthPacket) {
	if len(packets) == 0 {
		return
	}
	now := time.Now()
	for i := range packets {
		packets[i].waitAttempts++
		packets[i].NotBefore = now.Add(nextWaitingBackoff(packets[i].waitAttempts))
	}
	b.ethMtx.Lock()
	b.ethPackets = append(packets, b.ethPackets...)
	remaining := len(b.ethPackets)
	b.ethMtx.Unlock()
	log.Printf("[BatchBuilder] re-queued %d eth packet(s) (waiting) for retry (queue now %d)", len(packets), remaining)
}

func (b *BatchBuilder) CheckCosmos(config BatchConfig, ch chan<- CosmosBatch) {
	b.cosmosMtx.Lock()

	if len(b.cosmosPackets) == 0 {
		b.cosmosMtx.Unlock()
		return
	}

	// Hold back packets still under their waiting backoff (see CheckEth). Zero
	// NotBefore (legacy default) is always ready, so legacy behavior is unchanged.
	now := time.Now()
	ready, waiting := partitionReadyCosmos(b.cosmosPackets, now)
	if len(ready) == 0 {
		b.cosmosMtx.Unlock()
		return
	}

	flush := false
	reason := ""

	if len(ready) >= int(config.BatchSize) {
		flush = true
		reason = fmt.Sprintf("size limit reached (%d >= %d)", len(ready), config.BatchSize)
	} else if now.After(b.cosmosTimestamp.Add(config.BatchPeriods)) {
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
	if chunkSize <= 0 || chunkSize > len(ready) {
		chunkSize = len(ready)
	}
	chunk := make([]CosmosPacket, chunkSize)
	copy(chunk, ready[:chunkSize])
	b.cosmosPackets = append(ready[chunkSize:], waiting...)
	b.cosmosTimestamp = now
	log.Printf("[BatchBuilder] Flushing cosmos batch: %d packets (%s, queue remaining: %d)",
		len(chunk), reason, len(b.cosmosPackets))
	batch := CosmosBatch{Packets: chunk}
	b.cosmosMtx.Unlock()

	ch <- batch
}

// partitionReadyCosmos — see partitionReadyEth.
func partitionReadyCosmos(packets []CosmosPacket, now time.Time) (ready, waiting []CosmosPacket) {
	for _, p := range packets {
		if p.NotBefore.IsZero() || !now.Before(p.NotBefore) {
			ready = append(ready, p)
		} else {
			waiting = append(waiting, p)
		}
	}
	return ready, waiting
}

func (b *BatchBuilder) CheckEth(config BatchConfig, ch chan<- EthBatch) {
	b.ethMtx.Lock()

	if len(b.ethPackets) == 0 {
		b.ethMtx.Unlock()
		return
	}

	// Hold back packets still under their waiting backoff (not yet relayable): only
	// the ready ones are eligible to flush. If none are ready this is a quiet wait
	// — no flush, no log, no downstream RPC. Packets with a zero NotBefore (the
	// legacy default) are always ready, so legacy behavior is unchanged.
	now := time.Now()
	ready, waiting := partitionReadyEth(b.ethPackets, now)
	if len(ready) == 0 {
		b.ethMtx.Unlock()
		return
	}

	flush := false
	reason := ""

	if len(ready) >= int(config.BatchSize) {
		flush = true
		reason = fmt.Sprintf("size limit reached (%d >= %d)", len(ready), config.BatchSize)
	} else if now.After(b.ethTimestamp.Add(config.BatchPeriods)) {
		flush = true
		reason = fmt.Sprintf("time limit reached (%v elapsed)", time.Since(b.ethTimestamp).Round(time.Millisecond))
	}

	if !flush {
		b.ethMtx.Unlock()
		return
	}

	chunkSize := int(config.BatchSize)
	if chunkSize <= 0 || chunkSize > len(ready) {
		chunkSize = len(ready)
	}
	chunk := make([]EthPacket, chunkSize)
	copy(chunk, ready[:chunkSize])
	// Un-flushed ready packets plus the still-waiting ones stay queued.
	b.ethPackets = append(ready[chunkSize:], waiting...)
	b.ethTimestamp = now
	log.Printf("[BatchBuilder] Flushing eth batch: %d packets (%s, queue remaining: %d)",
		len(chunk), reason, len(b.ethPackets))
	batch := EthBatch{Packets: chunk}
	b.ethMtx.Unlock()

	ch <- batch
}

// partitionReadyEth splits packets into those eligible to flush now (zero or
// elapsed NotBefore) and those still under a waiting backoff, preserving order
// within each group.
func partitionReadyEth(packets []EthPacket, now time.Time) (ready, waiting []EthPacket) {
	for _, p := range packets {
		if p.NotBefore.IsZero() || !now.Before(p.NotBefore) {
			ready = append(ready, p)
		} else {
			waiting = append(waiting, p)
		}
	}
	return ready, waiting
}
