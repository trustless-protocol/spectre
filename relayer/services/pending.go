package services

import (
	"sync"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

type packetKey struct {
	sourceClient string
	sequence     uint64
}

type pendingPacketInfo struct {
	Packet      channeltypesv2.Packet
	ObservedAt  time.Time
	BlockNumber uint64
}

type PendingPacketTracker struct {
	mtx     sync.Mutex
	packets map[packetKey]pendingPacketInfo
}

func NewPendingPacketTracker() *PendingPacketTracker {
	return &PendingPacketTracker{
		packets: make(map[packetKey]pendingPacketInfo),
	}
}

func (t *PendingPacketTracker) Add(packet channeltypesv2.Packet, blockNumber uint64) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	key := packetKey{sourceClient: packet.SourceClient, sequence: packet.Sequence}
	t.packets[key] = pendingPacketInfo{
		Packet:      packet,
		ObservedAt:  time.Now(),
		BlockNumber: blockNumber,
	}
}

func (t *PendingPacketTracker) Remove(sourceClient string, sequence uint64) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	key := packetKey{sourceClient: sourceClient, sequence: sequence}
	delete(t.packets, key)
}

func (t *PendingPacketTracker) GetAll() []pendingPacketInfo {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	result := make([]pendingPacketInfo, 0, len(t.packets))
	for _, info := range t.packets {
		result = append(result, info)
	}
	return result
}

func (t *PendingPacketTracker) Len() int {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return len(t.packets)
}
