package evm

import (
	"testing"

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
