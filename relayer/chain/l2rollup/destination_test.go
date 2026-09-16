package l2rollup

import (
	"context"
	"strings"
	"testing"

	"relayer/chain"
	"relayer/chain/wasmclient"
	"relayer/services"

	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// TestBuildWasmUpdateClient verifies the L2 header is wrapped in a wasm
// ClientMessage inside a MsgUpdateClient — the exact wrapper Dũng's client (and
// the existing beacon client) consumes.
func TestBuildWasmUpdateClient(t *testing.T) {
	l2Header := []byte(`{"height":123,"state_root":"0xabc"}`)
	msg, err := wasmclient.BuildUpdateClient("cosmos1signer", "l2-op-0", l2Header)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if msg.ClientId != "l2-op-0" || msg.Signer != "cosmos1signer" {
		t.Fatalf("wrong envelope: client=%q signer=%q", msg.ClientId, msg.Signer)
	}
	// The ClientMessage Any must decode back to a wasm ClientMessage carrying the
	// L2 header verbatim.
	cm, ok := msg.ClientMessage.GetCachedValue().(*ibcwasmtypes.ClientMessage)
	if !ok {
		t.Fatalf("client message is not a wasm ClientMessage: %T", msg.ClientMessage.GetCachedValue())
	}
	if string(cm.Data) != string(l2Header) {
		t.Fatalf("L2 header not preserved: got %q", cm.Data)
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

// countingTxHandler records which submission path was taken and what it carried.
// Both methods, because the whole claim of folding is that ONE of them replaces
// two calls to the other.
type countingTxHandler struct {
	services.TransactionHandler
	batches [][]any // SendCosmosTxBatch — the unfolded path
	atomic  [][]any // SendCosmosTxBatchAtomic — the folded path
	err     error
}

func (h *countingTxHandler) CosmosSignerAddress() (string, error) { return "cosmos1signer", nil }

func (h *countingTxHandler) SendCosmosTxBatch(_ context.Context, _ services.CosmosEndpoint, msgs []any) error {
	h.batches = append(h.batches, msgs)
	return h.err
}

func (h *countingTxHandler) SendCosmosTxBatchAtomic(_ context.Context, _ services.CosmosEndpoint, msgs []any) error {
	h.atomic = append(h.atomic, msgs)
	return h.err
}

func (h *countingTxHandler) transactions() int { return len(h.batches) + len(h.atomic) }

func l2Packet(t *testing.T, seq, height uint64) chain.RelayPacket {
	t.Helper()
	pkt := channeltypesv2.Packet{Sequence: seq, SourceClient: "l2-0", DestinationClient: "08-wasm-3"}
	raw, err := pkt.Marshal()
	if err != nil {
		t.Fatalf("marshal packet: %v", err)
	}
	return chain.RelayPacket{Type: chain.SendPacket, Sequence: seq, Height: height, Packet: raw, Proof: []byte{0x01}}
}

// F1: this leg cost two transactions per batch while the other three cost one,
// on the leg with the highest latency of the four. The assertion is the COUNT,
// not that the folded path works -- a folded path that still submitted twice
// would pass every other test in this file.
func TestRelayWithUpdate_CostsOneTransactionInsteadOfTwo(t *testing.T) {
	const height = 4210
	packets := []chain.RelayPacket{l2Packet(t, 1, height), l2Packet(t, 2, height)}
	update := chain.ClientUpdate{Height: height, Payloads: [][]byte{{0xaa}}}

	// Unfolded, the update travels ALONE -- so whatever carries the packets has to
	// be a second transaction. Measured rather than assumed: the unfolded packet
	// path reads its height off the chain and cannot run without a node, so the
	// count for it is established by what this submission does NOT contain.
	unfolded := &countingTxHandler{}
	d := &Destination{worker: services.NewWorker(unfolded, nil), clientID: "08-wasm-3"}
	if err := d.UpdateClient(context.Background(), "08-wasm-3", update); err != nil {
		t.Fatalf("update: %v", err)
	}
	if unfolded.transactions() != 1 {
		t.Fatalf("UpdateClient sent %d transactions, want 1", unfolded.transactions())
	}
	if carried := len(unfolded.batches[0]) + len(unfolded.atomic); carried != 1 {
		t.Fatalf("the unfolded update carries %d messages; if it already carried the packets there would be nothing to fold", carried)
	}

	folded := &countingTxHandler{}
	d = &Destination{worker: services.NewWorker(folded, nil), clientID: "08-wasm-3"}
	if !d.SupportsUpdatePacketFolding() {
		t.Fatal("an atomic handler must advertise folding, or the engine never takes this path")
	}
	if err := d.RelayWithUpdate(context.Background(), "08-wasm-3", update, packets); err != nil {
		t.Fatalf("folded relay: %v", err)
	}

	if folded.transactions() != 1 {
		t.Fatalf("folded relay sent %d transactions, want 1", folded.transactions())
	}
	if len(folded.atomic) != 1 {
		t.Fatalf("folded relay used the splittable path; a split would land the update without its packets")
	}
	if got := len(folded.atomic[0]); got != 1+len(packets) {
		t.Fatalf("the folded batch carries %d messages, want the update plus %d packets", got, len(packets))
	}
	if _, ok := folded.atomic[0][0].(*clienttypes.MsgUpdateClient); !ok {
		t.Fatalf("the update is not first (%T); the packet proofs verify against the state it installs",
			folded.atomic[0][0])
	}
}

// The expensive half of F1. Unfolded, the proof height is READ from the chain and
// is true by construction; folded, the update is not on chain yet, so the two
// numbers are chosen in different places with nothing tying them together. The
// ETH path shipped this exact mistake (#443) and the chain's complaint named no
// height at all.
func TestRelayWithUpdate_RefusesProofsBuiltAtAnotherHeight(t *testing.T) {
	handler := &countingTxHandler{}
	d := &Destination{worker: services.NewWorker(handler, nil), clientID: "08-wasm-3"}

	err := d.RelayWithUpdate(context.Background(), "08-wasm-3",
		chain.ClientUpdate{Height: 4210, Payloads: [][]byte{{0xaa}}},
		[]chain.RelayPacket{l2Packet(t, 1, 4210), l2Packet(t, 2, 4209)}, // the second is stale
	)
	if err == nil {
		t.Fatal("a proof built at another height was accepted; the chain rejects it with an error naming no height")
	}
	if !strings.Contains(err.Error(), "4209") || !strings.Contains(err.Error(), "4210") {
		t.Errorf("error = %v; it must name both heights, since neither appears in the chain's own complaint", err)
	}
	// Before the transaction, not after: a batch that cannot verify is not worth
	// the fee, and an atomic batch takes the client update down with it.
	if handler.transactions() != 0 {
		t.Fatalf("sent %d transactions for a batch that cannot verify", handler.transactions())
	}
}

func TestRelayWithUpdate_RequiresAnAtomicHandler(t *testing.T) {
	// A splittable handler would land the update and lose the packets, or land
	// some packets against a client the same failed batch never advanced.
	d := &Destination{worker: services.NewWorker(&struct{ services.TransactionHandler }{}, nil), clientID: "08-wasm-3"}
	if d.SupportsUpdatePacketFolding() {
		t.Fatal("a handler without atomic submission advertised folding")
	}
	err := d.RelayWithUpdate(context.Background(), "08-wasm-3",
		chain.ClientUpdate{Height: 1, Payloads: [][]byte{{0xaa}}}, nil)
	if err == nil {
		t.Fatal("folded relay proceeded on a handler that cannot submit atomically")
	}
}

// The invariant folding rests on, asserted on the VALUE rather than the shape.
// The batch can carry the right number of messages, in the right order, with the
// update first -- and still declare a proof height the proofs were not built at,
// which is the whole failure #443 was. Counting messages does not see it.
func TestRelayWithUpdate_ProvesPacketsAtTheHeightTheUpdateInstalls(t *testing.T) {
	const height = 4210
	handler := &countingTxHandler{}
	d := &Destination{worker: services.NewWorker(handler, nil), clientID: "08-wasm-3"}

	packets := []chain.RelayPacket{l2Packet(t, 1, height), l2Packet(t, 2, height)}
	err := d.RelayWithUpdate(context.Background(), "08-wasm-3",
		chain.ClientUpdate{Height: height, Payloads: [][]byte{{0xaa}}}, packets)
	if err != nil {
		t.Fatalf("folded relay: %v", err)
	}
	if len(handler.atomic) != 1 {
		t.Fatalf("submissions = %d, want 1", len(handler.atomic))
	}

	proven := 0
	for _, msg := range handler.atomic[0] {
		recv, ok := msg.(*channeltypesv2.MsgRecvPacket)
		if !ok {
			continue // the update
		}
		proven++
		if recv.ProofHeight.RevisionHeight != height {
			t.Errorf("seq=%d declares proof height %d, but the update installs %d; the chain rejects this with an error naming neither",
				recv.Packet.Sequence, recv.ProofHeight.RevisionHeight, height)
		}
		if recv.ProofHeight.RevisionNumber != 0 {
			t.Errorf("seq=%d declares revision %d; an L2 height has no revision",
				recv.Packet.Sequence, recv.ProofHeight.RevisionNumber)
		}
	}
	if proven != len(packets) {
		t.Fatalf("checked %d packet messages, want %d", proven, len(packets))
	}
}
