package cosmos

import (
	"testing"

	"relayer/chain"
	"relayer/services"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

const routerClientID = "eth-client-0" // this Cosmos chain's ETH-side router id

func cosmosPacket(t services.CosmosPacketType, src, dst string) services.CosmosPacket {
	p := services.CosmosPacket{
		Type:        t,
		Packet:      &channeltypesv2.Packet{Sequence: 1, SourceClient: src, DestinationClient: dst},
		BlockNumber: 42,
	}
	// Acks carry acknowledgement bytes; without them the mapping skips the event.
	if t == services.CosmosAck {
		p.AckBytes = [][]byte{{0x01}}
	}
	return p
}

func TestCosmosPacketToEvent_TimeoutFilter(t *testing.T) {
	tests := []struct {
		name     string
		packet   services.CosmosPacket
		wantOK   bool
		wantType chain.EventType
	}{
		{
			// Cosmos-originated: destination_client == our router id → timeout fires
			// on Cosmos, ETH holds no commitment → must NOT relay to ETH.
			name:   "cosmos-originated timeout dropped",
			packet: cosmosPacket(services.CosmosTimeout, "08-wasm-0", routerClientID),
			wantOK: false,
		},
		{
			// ETH-originated: destination_client != our router id → ETH still owns a
			// commitment that the timeout must delete + refund → relay it.
			name:     "eth-originated timeout relayed",
			packet:   cosmosPacket(services.CosmosTimeout, routerClientID, "08-wasm-0"),
			wantOK:   true,
			wantType: chain.TimeoutPacket,
		},
		{
			name:     "send always relayed",
			packet:   cosmosPacket(services.CosmosSend, routerClientID, "08-wasm-0"),
			wantOK:   true,
			wantType: chain.SendPacket,
		},
		{
			name:     "ack always relayed",
			packet:   cosmosPacket(services.CosmosAck, "08-wasm-0", routerClientID),
			wantOK:   true,
			wantType: chain.AckPacket,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e, ok := cosmosPacketToEvent(tc.packet, routerClientID)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if ok && e.Type != tc.wantType {
				t.Fatalf("type = %d, want %d", e.Type, tc.wantType)
			}
		})
	}
}

// TestCosmosPacketToEvent_NilPacket guards the encode-failure skip.
func TestCosmosPacketToEvent_NilPacket(t *testing.T) {
	if _, ok := cosmosPacketToEvent(services.CosmosPacket{Type: services.CosmosSend}, routerClientID); ok {
		t.Fatal("nil packet must be skipped")
	}
}

// TestCosmosPacketToEvent_EmptyAckSkipped guards the no-ack-bytes skip: an ack
// without acknowledgement bytes cannot be relayed and must not be emitted.
func TestCosmosPacketToEvent_EmptyAckSkipped(t *testing.T) {
	p := services.CosmosPacket{
		Type:   services.CosmosAck,
		Packet: &channeltypesv2.Packet{Sequence: 1, SourceClient: "08-wasm-0", DestinationClient: routerClientID},
		// AckBytes deliberately empty
	}
	if _, ok := cosmosPacketToEvent(p, routerClientID); ok {
		t.Fatal("ack with no acknowledgement bytes must be skipped")
	}
}
