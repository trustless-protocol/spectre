package evm

import (
	"context"
	"testing"

	"relayer/chain"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// mustMarshalPacket returns a proto-marshaled minimal packet for RelayPacket.Packet.
func mustMarshalPacket(t *testing.T) []byte {
	t.Helper()
	pkt := channeltypesv2.Packet{Sequence: 7, SourceClient: "eth-router-0", DestinationClient: "08-wasm-0"}
	raw, err := pkt.Marshal()
	if err != nil {
		t.Fatalf("marshal packet: %v", err)
	}
	return raw
}

// These exercise RelayPackets' input validation — the branches that return before
// SendEthTxBatch, so no worker is needed. A zero-value Destination is sufficient.

// TestRelayPackets_AckMissingBytes: an ack packet with no acknowledgement bytes
// cannot build a MsgAckPacket and must error, not silently drop or panic.
func TestRelayPackets_AckMissingBytes(t *testing.T) {
	d := &Destination{}
	err := d.RelayPackets(context.Background(), []chain.RelayPacket{{
		Type:   chain.AckPacket,
		Packet: mustMarshalPacket(t),
		// AckBytes deliberately empty
	}})
	if err == nil {
		t.Fatal("ack packet without acknowledgement bytes must error")
	}
}

// TestRelayPackets_UnknownType guards the default branch (an event type the EVM
// destination does not build a message for).
func TestRelayPackets_UnknownType(t *testing.T) {
	d := &Destination{}
	err := d.RelayPackets(context.Background(), []chain.RelayPacket{{
		Type:   chain.EventType(99),
		Packet: mustMarshalPacket(t),
	}})
	if err == nil {
		t.Fatal("unknown packet type must error")
	}
}

// TestRelayPackets_MalformedPacket guards the packet-decode failure path.
func TestRelayPackets_MalformedPacket(t *testing.T) {
	d := &Destination{}
	// A byte string that is not a valid proto packet.
	err := d.RelayPackets(context.Background(), []chain.RelayPacket{{
		Type:   chain.SendPacket,
		Packet: []byte{0xff, 0xff, 0xff, 0xff},
	}})
	if err == nil {
		t.Fatal("malformed packet must error")
	}
}

// TestRelayPackets_Empty: no packets is a no-op (returns nil without touching the
// worker), matching the module's "nothing to relay" contract.
func TestRelayPackets_Empty(t *testing.T) {
	d := &Destination{}
	if err := d.RelayPackets(context.Background(), nil); err != nil {
		t.Fatalf("empty relay must be a no-op, got %v", err)
	}
}

// TestUpdateClient_BadPayload guards the codec-decode failure path (returns before
// the eth tx submit, so no worker needed).
func TestUpdateClient_BadPayload(t *testing.T) {
	d := &Destination{}
	err := d.UpdateClient(context.Background(), "", chain.ClientUpdate{Payload: []byte("not-a-gob-cosmos-update")})
	if err == nil {
		t.Fatal("undecodable client update payload must error")
	}
}
