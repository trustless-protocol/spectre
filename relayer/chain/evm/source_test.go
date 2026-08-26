package evm

import (
	"testing"
	"time"

	"relayer/chain"
	"relayer/services"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func ethPacket(t services.EthPacketType) services.EthPacket {
	p := services.EthPacket{
		Type:        t,
		Packet:      &channeltypesv2.Packet{Sequence: 1, SourceClient: "eth-router-0", DestinationClient: "08-wasm-0"},
		BlockNumber: 42,
	}
	// Only a WriteAck carries acknowledgement bytes; without them the mapping skips.
	if t == services.EthWriteAck {
		p.AckBytes = [][]byte{{0x01}}
	}
	return p
}

// TestEthPacketToEvent_TypeMapping mirrors cosmos TestCosmosPacketToEvent: it
// pins which queued EthPacket types become relay events and which are skipped
// (EthTimeout is scanner-owned; EthAck is a terminal settle, not relayed here).
func TestEthPacketToEvent_TypeMapping(t *testing.T) {
	tests := []struct {
		name     string
		packet   services.EthPacket
		wantOK   bool
		wantType chain.EventType
	}{
		{
			name:     "send relayed as recv",
			packet:   ethPacket(services.EthSend),
			wantOK:   true,
			wantType: chain.SendPacket,
		},
		{
			name:     "write-ack relayed as ack",
			packet:   ethPacket(services.EthWriteAck),
			wantOK:   true,
			wantType: chain.AckPacket,
		},
		{
			// EthTimeout is refunded on ETH by the async scanner, never relayed here.
			name:   "eth timeout skipped (scanner-owned)",
			packet: ethPacket(services.EthTimeout),
			wantOK: false,
		},
		{
			// EthAck settles a pending ETH-origin send; it is not a relayable event.
			name:   "eth ack skipped (terminal settle)",
			packet: ethPacket(services.EthAck),
			wantOK: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e, ok := ethPacketToEvent(tc.packet)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if ok && e.Type != tc.wantType {
				t.Fatalf("type = %d, want %d", e.Type, tc.wantType)
			}
		})
	}
}

// TestEthPacketToEvent_AckBytesCarried guards that an ack event carries its
// acknowledgement bytes through to the event (needed to build MsgAcknowledgement).
func TestEthPacketToEvent_AckBytesCarried(t *testing.T) {
	e, ok := ethPacketToEvent(ethPacket(services.EthWriteAck))
	if !ok {
		t.Fatal("write-ack with bytes must be relayed")
	}
	if len(e.AckBytes) != 1 || len(e.AckBytes[0]) != 1 || e.AckBytes[0][0] != 0x01 {
		t.Fatalf("ack bytes not carried: %v", e.AckBytes)
	}
}

// TestEthPacketToEvent_EmptyAckSkipped guards the no-ack-bytes skip: a write-ack
// without acknowledgement bytes cannot be built into a MsgAcknowledgement.
func TestEthPacketToEvent_EmptyAckSkipped(t *testing.T) {
	p := services.EthPacket{
		Type:   services.EthWriteAck,
		Packet: &channeltypesv2.Packet{Sequence: 1, SourceClient: "eth-router-0", DestinationClient: "08-wasm-0"},
		// AckBytes deliberately empty
	}
	if _, ok := ethPacketToEvent(p); ok {
		t.Fatal("write-ack with no acknowledgement bytes must be skipped")
	}
}

// TestEthPacketToEvent_NilPacket guards the encode-failure skip.
func TestEthPacketToEvent_NilPacket(t *testing.T) {
	if _, ok := ethPacketToEvent(services.EthPacket{Type: services.EthSend}); ok {
		t.Fatal("nil packet must be skipped")
	}
}

// Mirror of chain/cosmos TestEventsWithOrigins_StaysIndexAligned. The handler
// reports failures as indices into the events slice and Subscribe re-queues
// orig[idx], so a packet dropped from one slice and not the other shifts every
// later index: the relayer re-queues a packet that already succeeded and forgets
// the one that failed.
//
// The ETH side drops packets for one more reason than Cosmos does — a terminal
// EthAck/EthTimeout settles the pending tracker and is not relayed — so those are
// interleaved here too.
func TestEventsWithOrigins_StaysIndexAligned(t *testing.T) {
	seq := func(p services.EthPacket, n uint64) services.EthPacket {
		p.Packet.Sequence = n
		return p
	}
	ackNoBytes := func(n uint64) services.EthPacket {
		p := ethPacket(services.EthWriteAck)
		p.AckBytes = nil // no ack bytes: cannot build a MsgAcknowledgement, skipped
		p.Packet.Sequence = n
		return p
	}

	packets := []services.EthPacket{
		seq(ethPacket(services.EthAck), 1),     // terminal: settled, not relayed
		seq(ethPacket(services.EthSend), 2),    //
		ackNoBytes(3),                          // skipped, between
		seq(ethPacket(services.EthTimeout), 4), // terminal: settled, not relayed
		seq(ethPacket(services.EthSend), 5),    //
		seq(ethPacket(services.EthWriteAck), 6),
		ackNoBytes(7), // skipped, last
	}

	var settled []uint64
	events, orig := eventsWithOrigins(packets, func(p channeltypesv2.Packet) {
		settled = append(settled, p.Sequence)
	})

	if len(events) != len(orig) {
		t.Fatalf("slices out of step: %d events vs %d origins", len(events), len(orig))
	}
	for i := range events {
		if orig[i].Packet == nil {
			t.Fatalf("origin %d has no packet", i)
		}
		if events[i].Sequence != orig[i].Packet.Sequence {
			t.Fatalf("index %d maps event seq=%d to origin seq=%d",
				i, events[i].Sequence, orig[i].Packet.Sequence)
		}
	}
	want := []uint64{2, 5, 6}
	if len(events) != len(want) {
		t.Fatalf("want %d relayable packets, got %d", len(want), len(events))
	}
	for i, w := range want {
		if events[i].Sequence != w {
			t.Fatalf("events[%d].Sequence = %d, want %d (dropped packets must leave BOTH slices)",
				i, events[i].Sequence, w)
		}
	}
	// A terminal event must still reach the tracker even though it is not relayed:
	// dropping it silently leaves the timeout scanner chasing a settled packet.
	wantSettled := []uint64{1, 4}
	if len(settled) != len(wantSettled) {
		t.Fatalf("settled %v, want %v", settled, wantSettled)
	}
	for i, w := range wantSettled {
		if settled[i] != w {
			t.Fatalf("settled %v, want %v", settled, wantSettled)
		}
	}
}

// An ETH→Cosmos send past its timeout can never be received on Cosmos, so the
// relay path must hand it to the timeout scanner instead of proving it. Getting
// the direction wrong is expensive both ways: a live packet dropped as dead, or a
// dead packet re-proven forever. The boundary second is included because the chain
// itself treats an equal timestamp as expired.
func TestSendPacketExpired(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)

	tests := []struct {
		name    string
		timeout uint64
		want    bool
	}{
		{"no timeout set never expires", 0, false},
		{"timeout well in the future", uint64(now.Unix()) + 600, false},
		{"one second left", uint64(now.Unix()) + 1, false},
		{"timeout is exactly now", uint64(now.Unix()), true},
		{"one second past", uint64(now.Unix()) - 1, true},
		{"long past", uint64(now.Unix()) - 86400, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt := channeltypesv2.Packet{Sequence: 1, TimeoutTimestamp: tt.timeout}
			if got := sendPacketExpired(pkt, now); got != tt.want {
				t.Fatalf("sendPacketExpired(timeout=%d, now=%d) = %v, want %v",
					tt.timeout, now.Unix(), got, tt.want)
			}
		})
	}
}
