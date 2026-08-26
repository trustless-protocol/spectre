package cosmos

import (
	"context"
	"encoding/json"
	"testing"

	"relayer/chain"
	relayerclient "relayer/client"
	"relayer/services"

	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

type atomicTestTxHandler struct{ services.TransactionHandler }

func (*atomicTestTxHandler) SendCosmosTxBatchAtomic(context.Context, services.CosmosEndpoint, []any) error {
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

// --- chain.Destination contract: the guards that run before any RPC ---
//
// These mirror the set chain/evm/destination_test.go already has. The Cosmos side
// had tests for its helper functions but none for the interface methods
// themselves, so every early-return below — the ones that decide whether a
// transaction is attempted at all — was unguarded.

// cancelledCtx returns a context that is already done, standing in for shutdown.
func cancelledCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

// Beacon headers must be wrapped one-for-one and stay in order: the wasm client
// applies them sequentially, and a dropped or reordered header leaves it at a slot
// nobody chose. Order is the property worth pinning here — length alone would pass
// a batch that shuffled.
//
// Note there is no "rejects a bad payload" case: BuildUpdateClient wraps whatever
// bytes it is given, and its only error branch (NewAnyWithValue on a nil message)
// cannot be reached from here, because the message it wraps is always non-nil.
// Validation of the header contents happens on-chain, in the wasm client.
func TestBuildBeaconUpdateMessages_WrapsEveryPayloadInOrder(t *testing.T) {
	payloads := [][]byte{[]byte("header-a"), []byte("header-b"), []byte("header-c")}

	msgs, err := buildBeaconUpdateMessages("signer-1", "client-0", payloads)
	if err != nil {
		t.Fatalf("wrap payloads: %v", err)
	}
	if len(msgs) != len(payloads) {
		t.Fatalf("got %d messages for %d payloads", len(msgs), len(payloads))
	}
	for i, m := range msgs {
		msg, ok := m.(*clienttypes.MsgUpdateClient)
		if !ok {
			t.Fatalf("message %d is %T, want *MsgUpdateClient", i, m)
		}
		if msg.ClientId != "client-0" || msg.Signer != "signer-1" {
			t.Fatalf("message %d carries client=%q signer=%q", i, msg.ClientId, msg.Signer)
		}
		var wrapped ibcwasmtypes.ClientMessage
		if err := wrapped.Unmarshal(msg.ClientMessage.Value); err != nil {
			t.Fatalf("message %d: unwrap client message: %v", i, err)
		}
		if string(wrapped.Data) != string(payloads[i]) {
			t.Fatalf("message %d carries payload %q, want %q (order must be preserved)",
				i, wrapped.Data, payloads[i])
		}
	}
}

// mustMarshalCosmosPacket returns proto bytes for a minimal ETH-origin packet.
func mustMarshalCosmosPacket(t *testing.T) []byte {
	t.Helper()
	pkt := channeltypesv2.Packet{Sequence: 7, SourceClient: "eth-router-0", DestinationClient: "08-wasm-0"}
	raw, err := pkt.Marshal()
	if err != nil {
		t.Fatalf("marshal packet: %v", err)
	}
	return raw
}

// The chain.Destination contract on the Cosmos side: the guards that run before
// any RPC, and decide whether a transaction is attempted at all.
//
// The name is TestDestination, not TestCosmosDestination: E9 requires the prefix
// to be a production identifier, and the type here is Destination. chain/evm has
// a Destination too, so the two test names collide across packages -- which the
// plan calls correct, not a clash: same unit, same name.
func TestDestination(t *testing.T) {
	// Nothing may be submitted once the run context is cancelled. A transaction
	// started during shutdown is broadcast with no one left to observe the result:
	// the process exits, and the tracker never learns whether it landed.
	t.Run("refuses work after shutdown", func(t *testing.T) {
		d := &Destination{}
		packets := []chain.RelayPacket{{Type: chain.SendPacket, Packet: mustMarshalCosmosPacket(t)}}

		if err := d.RelayPackets(cancelledCtx(), packets); err == nil {
			t.Error("RelayPackets started a packet tx during shutdown")
		}
		if err := d.UpdateClient(cancelledCtx(), "", chain.ClientUpdate{Payloads: [][]byte{{1}}}); err == nil {
			t.Error("UpdateClient started a client-update tx during shutdown")
		}
		if err := d.RelayWithUpdate(cancelledCtx(), "", chain.ClientUpdate{Payloads: [][]byte{{1}}}, packets); err == nil {
			t.Error("RelayWithUpdate started a folded tx during shutdown")
		}
	})

	// An empty packet list is a no-op, not an error: the module can hand over a batch
	// that emptied out after filtering, and turning that into an error would fail a
	// relay cycle that had nothing wrong with it.
	t.Run("treats an empty packet list as a no-op", func(t *testing.T) {
		d := &Destination{}
		if err := d.RelayPackets(context.Background(), nil); err != nil {
			t.Fatalf("empty packet list must be a no-op, got %v", err)
		}
	})

	// An empty payload list means "client already current" on this side, and must be
	// a silent success — the beacon builder returns one whenever the on-chain client
	// already covers the target slot.
	//
	// This deliberately DIFFERS from chain/evm, where UpdateClient requires exactly
	// one payload and errors otherwise. The two are not inconsistent: a beacon update
	// is an ordered batch of headers that may legitimately be empty, while the groth16
	// path submits exactly one proof per update. Anyone harmonising the two signatures
	// should break this test and read this comment first.
	t.Run("treats no payloads as client-already-current", func(t *testing.T) {
		d := &Destination{}
		if err := d.UpdateClient(context.Background(), "client-0", chain.ClientUpdate{}); err != nil {
			t.Fatalf("empty payload list means the client is current, got %v", err)
		}
	})

	// Folding several messages into one transaction only works if the handler can
	// submit them atomically. Without that guarantee a split would land the client
	// update in one block and the packets in another, against a client state the
	// packets' proof height no longer matches — so this must be refused up front
	// rather than discovered after the first half is already on chain.
	t.Run("refuses folding without an atomic handler", func(t *testing.T) {
		update := chain.ClientUpdate{Payloads: [][]byte{mustBeaconHeader(t, "100")}}
		packets := []chain.RelayPacket{{Type: chain.SendPacket, Packet: mustMarshalCosmosPacket(t)}}

		tests := []struct {
			name string
			dest *Destination
		}{
			{"no worker at all", &Destination{}},
			{"worker without a transaction handler", &Destination{worker: services.NewWorker(nil, nil)}},
			{
				"handler that cannot submit atomically",
				&Destination{worker: services.NewWorker(&struct{ services.TransactionHandler }{}, nil)},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := tt.dest.RelayWithUpdate(context.Background(), "client-0", update, packets); err == nil {
					t.Fatal("folded relay proceeded without an atomic sender")
				}
			})
		}
	})

	// A folded relay with no client update is a programming error, not an empty
	// batch: the whole point of the call is to install the update the packets prove
	// against. Falling through would submit packets at a proof height the destination
	// client has not reached.
	t.Run("refuses folding with no client update", func(t *testing.T) {
		d := &Destination{worker: services.NewWorker(&atomicTestTxHandler{}, nil)}
		err := d.RelayWithUpdate(context.Background(), "client-0", chain.ClientUpdate{}, nil)
		if err == nil {
			t.Fatal("folded relay accepted an empty client update")
		}
	})

	// HasPacketReceipt is deliberately unimplemented on this side. It must report that
	// loudly rather than answering "no receipt", which a caller would read as "safe to
	// time this packet out" — refunding a packet that was in fact delivered.
	t.Run("reports that HasPacketReceipt is not wired", func(t *testing.T) {
		d := &Destination{}
		ok, err := d.HasPacketReceipt(context.Background(), nil)
		if err == nil {
			t.Fatal("an unimplemented receipt check must error, not answer")
		}
		if ok {
			t.Fatal("an unimplemented receipt check must never report a receipt")
		}
	})
}
