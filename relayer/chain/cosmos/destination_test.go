package cosmos

import (
	"context"
	"encoding/json"
	"testing"

	"relayer/chain"
	relayerclient "relayer/client"
	"relayer/services"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

type atomicTestTxHandler struct{ services.TransactionHandler }

func (*atomicTestTxHandler) SendCosmosTxBatchAtomic(context.Context, services.Context, []any) error {
	return nil
}

func TestSupportsUpdatePacketFolding_RequiresAtomicSender(t *testing.T) {
	withoutAtomic := &Destination{worker: services.NewWorker(&struct{ services.TransactionHandler }{}, nil)}
	if withoutAtomic.SupportsUpdatePacketFolding() {
		t.Fatal("handler without no-split submission must not advertise folding")
	}
	withAtomic := &Destination{worker: services.NewWorker(&atomicTestTxHandler{}, nil)}
	if !withAtomic.SupportsUpdatePacketFolding() {
		t.Fatal("handler with no-split submission must advertise folding")
	}
}

func mustBeaconHeader(t *testing.T, signatureSlot string) []byte {
	t.Helper()
	payload, err := json.Marshal(relayerclient.EthereumHeader{
		ConsensusUpdate: relayerclient.LightClientUpdate{
			SignatureSlot: signatureSlot,
			FinalizedHeader: relayerclient.LightClientHeader{
				Beacon: relayerclient.BeaconBlockHeader{Slot: signatureSlot},
			},
		},
	})
	if err != nil {
		t.Fatalf("marshal header: %v", err)
	}
	return payload
}

func TestHighestBeaconSignatureSlot_UsesLatestHeader(t *testing.T) {
	slot, err := highestBeaconSignatureSlot([][]byte{
		mustBeaconHeader(t, "123"),
		mustBeaconHeader(t, "130"),
		mustBeaconHeader(t, "127"),
	})
	if err != nil {
		t.Fatalf("highest slot: %v", err)
	}
	if slot != 130 {
		t.Fatalf("slot = %d, want 130", slot)
	}
}

func TestHighestBeaconSignatureSlot_RejectsMalformedPayload(t *testing.T) {
	if _, err := highestBeaconSignatureSlot([][]byte{[]byte("not-json")}); err == nil {
		t.Fatal("malformed header must fail")
	}
}

func TestHighestBeaconSignatureSlot_RejectsInvalidSlot(t *testing.T) {
	if _, err := highestBeaconSignatureSlot([][]byte{mustBeaconHeader(t, "not-a-slot")}); err == nil {
		t.Fatal("invalid signature slot must fail")
	}
}

func TestHighestBeaconFinalizedSlot_UsesLatestHeader(t *testing.T) {
	slot, err := highestBeaconFinalizedSlot([][]byte{
		mustBeaconHeader(t, "123"),
		mustBeaconHeader(t, "130"),
		mustBeaconHeader(t, "127"),
	})
	if err != nil {
		t.Fatalf("highest finalized slot: %v", err)
	}
	if slot != 130 {
		t.Fatalf("finalized slot = %d, want 130", slot)
	}
}

func TestBuildPacketMessages_UsesFoldedUpdateTargetSlot(t *testing.T) {
	pkt := channeltypesv2.Packet{Sequence: 7, SourceClient: "eth-router-0", DestinationClient: "08-wasm-0"}
	raw, err := pkt.Marshal()
	if err != nil {
		t.Fatalf("marshal packet: %v", err)
	}
	wantHeight := clienttypes.Height{RevisionNumber: 0, RevisionHeight: 130}
	msgs, err := buildPacketMessages("cosmos1signer", wantHeight, []chain.RelayPacket{{
		Type: chain.SendPacket, Packet: raw, Proof: []byte("proof"), Height: 999,
	}})
	if err != nil {
		t.Fatalf("build packet messages: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("message count = %d, want 1", len(msgs))
	}
	msg, ok := msgs[0].(*channeltypesv2.MsgRecvPacket)
	if !ok {
		t.Fatalf("message type = %T, want *MsgRecvPacket", msgs[0])
	}
	if msg.ProofHeight != wantHeight {
		t.Fatalf("proof height = %v, want %v", msg.ProofHeight, wantHeight)
	}
}
