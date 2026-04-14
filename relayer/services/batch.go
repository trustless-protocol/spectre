package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func (p PacketType) String() string {
	switch p {
	case Send:
		return "Send"
	case Ack:
		return "Ack"
	case Timeout:
		return "Timeout"
	case WriteAck:
		return "WriteAck"
	default:
		return fmt.Sprintf("Unknown(%d)", int(p))
	}
}

type PacketType int

const (
	Send     PacketType = iota // Cosmos→ETH: RecvPacket on ETH
	Ack                        // Cosmos→ETH: AckPacket on ETH
	Timeout                    // Timeout
	WriteAck                   // Cosmos→ETH: ETH wrote ack → submit MsgAcknowledgement to Cosmos
)

type Packet struct {
	PacketType  PacketType
	Packet      *channeltypesv2.Packet
	AckBytes    [][]byte
	BlockNumber uint64 // ETH block where event was emitted (for beacon finality checks)
}

type BatchPackets struct {
	Packets []Packet
}

type BatchBuilder struct {
	mtx       sync.Mutex
	timestamp time.Time
	packets   []Packet
}

func NewBatchBuilder() *BatchBuilder {
	return &BatchBuilder{
		timestamp: time.Now(),
		packets:   []Packet{},
	}
}

func (b *BatchBuilder) InsertPacket(packet Packet) {
	b.mtx.Lock()
	b.packets = append(b.packets, packet)
	count := len(b.packets)
	b.mtx.Unlock()
	log.Printf("[BatchBuilder] Inserted packet: type=%s seq=%d (batch size: %d)",
		packet.PacketType, packet.Packet.Sequence, count)
}

func (b *BatchBuilder) ClearBatch() {
	b.timestamp = time.Now()
	b.packets = []Packet{}
}

func (b *BatchBuilder) CheckBatch(config BatchConfig, ch chan<- BatchPackets) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if len(b.packets) == 0 {
		return
	}

	flush := false
	reason := ""

	if len(b.packets) >= int(config.BatchSize) {
		flush = true
		reason = fmt.Sprintf("size limit reached (%d >= %d)", len(b.packets), config.BatchSize)
	} else if time.Now().After(b.timestamp.Add(config.BatchPeriods)) {
		flush = true
		reason = fmt.Sprintf("time limit reached (%v elapsed)", time.Since(b.timestamp).Round(time.Millisecond))
	}

	if flush {
		log.Printf("[BatchBuilder] Flushing batch: %d packets (%s)", len(b.packets), reason)
		ch <- BatchPackets{
			Packets: b.packets,
		}
		b.ClearBatch()
	}
}
