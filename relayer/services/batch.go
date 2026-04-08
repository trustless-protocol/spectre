package services

import (
	"sync"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

type PacketType int

const (
	Send     PacketType = iota // Cosmos→ETH: RecvPacket on ETH
	Ack                        // ETH→Cosmos: AckPacket on ETH (triggered by Cosmos acknowledge_packet)
	Timeout                    // Timeout
	WriteAck                   // Cosmos→ETH: ETH wrote ack → submit MsgAcknowledgement to Cosmos
)

type Packet struct {
	PacketType PacketType
	Packet     *channeltypesv2.Packet
	AckBytes   [][]byte // populated for WriteAck packets (one entry per payload)
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
	b.mtx.Unlock()
}

func (b *BatchBuilder) ClearBatch() {
	b.timestamp = time.Now()
	b.packets = []Packet{}
}

func (b *BatchBuilder) CheckBatch(config BatchConfig, ch chan<- BatchPackets) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if len(b.packets) < int(config.BatchSize) && time.Now().After(b.timestamp.Add(config.BatchPeriods)) {
		ch <- BatchPackets{
			Packets: b.packets,
		}

		b.ClearBatch()
	} else if len(b.packets) > int(config.BatchSize) {
		ch <- BatchPackets{
			Packets: b.packets,
		}

		b.ClearBatch()
	}
}
