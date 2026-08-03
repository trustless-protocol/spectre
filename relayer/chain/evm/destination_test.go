package evm

import (
	"context"
	"testing"

	contractICS26Router "relayer/bindings/ICS26Router"
	spectreContract "relayer/bindings/SpectreClient"
	updateclientContract "relayer/bindings/UpdateClient"
	"relayer/chain"
	"relayer/chain/codec"
	"relayer/services"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

type recordingTxHandler struct {
	services.TransactionHandler
	ethBatches [][]any
}

func (h *recordingTxHandler) SendEthTxBatch(_ context.Context, _ services.Context, msgs []any) error {
	h.ethBatches = append(h.ethBatches, msgs)
	return nil
}

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
	err := d.UpdateClient(context.Background(), "", chain.ClientUpdate{Payloads: [][]byte{[]byte("not-a-gob-cosmos-update")}})
	if err == nil {
		t.Fatal("undecodable client update payload must error")
	}
}

func TestUpdateClient_RequiresExactlyOnePayload(t *testing.T) {
	d := &Destination{}
	if err := d.UpdateClient(context.Background(), "", chain.ClientUpdate{}); err == nil {
		t.Fatal("empty payload list must fail")
	}
	if err := d.UpdateClient(context.Background(), "", chain.ClientUpdate{Payloads: [][]byte{[]byte("one"), []byte("two")}}); err == nil {
		t.Fatal("multiple payloads must fail")
	}
}

func TestSupportsUpdatePacketFolding_RequiresRouterProofSubmitter(t *testing.T) {
	const (
		router = "0x0000000000000000000000000000000000000001"
		other  = "0x0000000000000000000000000000000000000002"
	)
	svcCtx := services.NewCtx(nil, nil)
	svcCtx.SetAddresses(router, "", "", "", "", other)
	d := &Destination{svcCtx: svcCtx}
	if d.SupportsUpdatePacketFolding() {
		t.Fatal("different role manager must disable folding")
	}
	svcCtx.SetAddresses(router, "", "", "", "", router)
	d.svcCtx = svcCtx
	if !d.SupportsUpdatePacketFolding() {
		t.Fatal("router proof submitter must enable folding")
	}
}

func TestRelayWithUpdate_SubmitsOneOrderedBatch(t *testing.T) {
	const router = "0x0000000000000000000000000000000000000001"
	svcCtx := services.NewCtx(nil, nil)
	svcCtx.SetAddresses(router, "", "", "", "", router)
	handler := &recordingTxHandler{}
	d := &Destination{worker: services.NewWorker(handler, nil), svcCtx: svcCtx}

	payload, err := codec.EncodeCosmosUpdate(
		int(services.ApplicationUpdate),
		updateclientContract.ISpectreClientMsgsMsgUpdateApplicationState{},
		spectreContract.IICS07TendermintMsgsValidatorSet{},
	)
	if err != nil {
		t.Fatalf("encode update: %v", err)
	}
	err = d.RelayWithUpdate(context.Background(), "client-0", chain.ClientUpdate{
		Height: 80, Payloads: [][]byte{payload},
	}, []chain.RelayPacket{{
		Type: chain.SendPacket, Packet: mustMarshalPacket(t), Proof: []byte("proof"), Height: 80,
	}})
	if err != nil {
		t.Fatalf("folded relay: %v", err)
	}
	if len(handler.ethBatches) != 1 {
		t.Fatalf("batch submissions = %d, want 1", len(handler.ethBatches))
	}
	batch := handler.ethBatches[0]
	if len(batch) != 2 {
		t.Fatalf("inner message count = %d, want update + packet", len(batch))
	}
	if _, ok := batch[0].(services.CosmosClientUpdateBuildResult); !ok {
		t.Fatalf("first inner message = %T, want CosmosClientUpdateBuildResult", batch[0])
	}
	if _, ok := batch[1].(contractICS26Router.IICS26RouterMsgsMsgRecvPacket); !ok {
		t.Fatalf("second inner message = %T, want MsgRecvPacket", batch[1])
	}
}
