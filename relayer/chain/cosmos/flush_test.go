package cosmos

import (
	"encoding/hex"
	"io"
	"log"
	"testing"

	"github.com/cosmos/gogoproto/proto"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"

	abcitypes "github.com/cometbft/cometbft/abci/types"

	"relayer/services"
)

// commitmentFor is the hash the chain stores for a packet, and therefore what a
// flush pass compares a candidate transaction against. Built through the same
// ibc-go function the chain uses, not restated here, so the test cannot agree
// with the code about a value the chain would disagree with.
func commitmentFor(packet channeltypesv2.Packet) flushCommitment {
	return flushCommitment{sequence: packet.Sequence, data: channeltypesv2.CommitPacket(packet)}
}

func encodedSendAttributes(packet channeltypesv2.Packet) []abcitypes.EventAttribute {
	raw, err := proto.Marshal(&packet)
	if err != nil {
		panic(err)
	}
	return []abcitypes.EventAttribute{
		{Key: "packet_source_client", Value: packet.SourceClient},
		{Key: "packet_sequence", Value: "7"},
		{Key: "encoded_packet_hex", Value: hex.EncodeToString(raw)},
	}
}

// TestPacketFromSendEvent: the enumeration path rebuilds a packet from the same
// encoded attribute the subscriber consumes, so one wire format serves both.
// Decoding a malformed one must report failure rather than hand a zero packet
// to the relay loop, which would relay sequence 0 on an empty client id.
func TestPacketFromSendEvent(t *testing.T) {
	t.Parallel()

	want := channeltypesv2.Packet{
		Sequence:          7,
		SourceClient:      "08-wasm-0",
		DestinationClient: "client-0",
		TimeoutTimestamp:  1_700_000_000,
	}

	t.Run("decodes the encoded packet", func(t *testing.T) {
		t.Parallel()
		got, ok := packetFromSendEvent(encodedSendAttributes(want))
		if !ok {
			t.Fatal("a well-formed send_packet event must decode")
		}
		if got.Sequence != want.Sequence || got.SourceClient != want.SourceClient ||
			got.DestinationClient != want.DestinationClient {
			t.Fatalf("decoded %+v, want %+v", got, want)
		}
	})

	t.Run("reports failure rather than returning a zero packet", func(t *testing.T) {
		t.Parallel()
		for _, tc := range []struct {
			name       string
			attributes []abcitypes.EventAttribute
		}{
			{"no encoded attribute", []abcitypes.EventAttribute{{Key: "packet_sequence", Value: "7"}}},
			{"not hex", []abcitypes.EventAttribute{{Key: "encoded_packet_hex", Value: "zzzz"}}},
			{"hex but not a packet", []abcitypes.EventAttribute{{Key: "encoded_packet_hex", Value: "ffffffff"}}},
			{"no attributes at all", nil},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				if _, ok := packetFromSendEvent(tc.attributes); ok {
					t.Fatal("reported success on an undecodable event; the relay loop would then relay sequence 0")
				}
			})
		}
	})
}

// TestMustMarshalPacket: the packet was just decoded from the chain's own event,
// so a failure here means it did not round-trip. Returning nil keeps a malformed
// event from taking the process down.
func TestMustMarshalPacket(t *testing.T) {
	t.Parallel()

	packet := channeltypesv2.Packet{Sequence: 7, SourceClient: "08-wasm-0", DestinationClient: "client-0"}
	raw := mustMarshalPacket(packet)
	if len(raw) == 0 {
		t.Fatal("a valid packet must re-encode")
	}
	var round channeltypesv2.Packet
	if err := proto.Unmarshal(raw, &round); err != nil {
		t.Fatalf("re-encoded packet does not decode: %v", err)
	}
	if round.Sequence != packet.Sequence || round.SourceClient != packet.SourceClient {
		t.Fatalf("round trip changed the packet: %+v -> %+v", packet, round)
	}
}

// Cosmos stores a send commitment under the packet's LOCAL source client, which
// for a Cosmos->EVM packet is the EVM light client on Cosmos (EVMOnCosmos), not
// the id by which the EVM router knows Cosmos (CosmosOnEVM). Querying the latter
// enumerates an empty set, so the backstop silently finds nothing on exactly the
// path it exists for.
func TestFlushQueriesTheLocalSourceClient(t *testing.T) {
	source := &Source{
		ids:    services.ClientIDs{CosmosOnEVM: "eth-side-id", EVMOnCosmos: "08-wasm-7"},
		logger: log.New(io.Discard, "", 0),
	}
	if got := source.commitmentQueryClientID(); got != "08-wasm-7" {
		t.Fatalf("commitment query client id = %q, want the local source client %q", got, "08-wasm-7")
	}
}

// TxSearch matches TRANSACTIONS; the event list it hands back is the whole
// transaction's, unfiltered. A transaction carrying several sends must not
// answer a request for sequence N with the first send it happens to contain.
func TestSendEventForSequencePicksTheMatchingEvent(t *testing.T) {
	events := []abcitypes.Event{
		{Type: "send_packet", Attributes: encodedSendAttributes(channeltypesv2.Packet{
			Sequence: 41, SourceClient: "08-wasm-7", DestinationClient: "eth-side-id",
		})},
		{Type: "send_packet", Attributes: encodedSendAttributes(channeltypesv2.Packet{
			Sequence: 42, SourceClient: "08-wasm-7", DestinationClient: "eth-side-id",
		})},
	}

	wanted := channeltypesv2.Packet{Sequence: 42, SourceClient: "08-wasm-7", DestinationClient: "eth-side-id"}
	got, found := matchSendEvent(events, "08-wasm-7", commitmentFor(wanted))
	if !found {
		t.Fatal("the requested sequence is present in the transaction and must be found")
	}
	if got.Sequence != 42 {
		t.Fatalf("returned sequence %d for a request for 42; the earlier send in the same tx was taken", got.Sequence)
	}

	absent := channeltypesv2.Packet{Sequence: 43, SourceClient: "08-wasm-7", DestinationClient: "eth-side-id"}
	if _, found := matchSendEvent(events, "08-wasm-7", commitmentFor(absent)); found {
		t.Fatal("a sequence absent from the transaction must report not-found, not the nearest send")
	}
}

// Sequences are per-client, so seq=N exists independently on every client. The
// query names the client, but that only decides which TRANSACTIONS come back --
// one carrying sends for two clients matches it and offers both events. Taking
// the wrong client's send relays the wrong packet to the wrong destination and
// leaves the sequence that was actually asked about outstanding.
func TestSendEventForSequenceRejectsAnotherClientsSameSequence(t *testing.T) {
	events := []abcitypes.Event{
		{Type: "send_packet", Attributes: encodedSendAttributes(channeltypesv2.Packet{
			Sequence: 42, SourceClient: "08-wasm-OTHER", DestinationClient: "other-destination",
		})},
		{Type: "send_packet", Attributes: encodedSendAttributes(channeltypesv2.Packet{
			Sequence: 42, SourceClient: "08-wasm-7", DestinationClient: "eth-side-id",
		})},
	}

	wanted := channeltypesv2.Packet{Sequence: 42, SourceClient: "08-wasm-7", DestinationClient: "eth-side-id"}
	got, found := matchSendEvent(events, "08-wasm-7", commitmentFor(wanted))
	if !found {
		t.Fatal("the requested client's send is present and must be found")
	}
	if got.SourceClient != "08-wasm-7" || got.DestinationClient != "eth-side-id" {
		t.Fatalf("returned src=%s dst=%s; another client's seq=42 was taken, so the requested packet stays outstanding",
			got.SourceClient, got.DestinationClient)
	}

	// And a client with no send in this transaction must not borrow one.
	if _, found := matchSendEvent(events, "08-wasm-ABSENT", commitmentFor(wanted)); found {
		t.Fatal("a client absent from the transaction must report not-found")
	}
}
