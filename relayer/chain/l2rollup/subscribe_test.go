package l2rollup

import (
	"math/big"
	"testing"

	"relayer/chain"

	contractICS26Router "relayer/bindings/ICS26Router"
)

func TestKeepIndexed(t *testing.T) {
	events := []chain.Event{
		{Height: 1}, {Height: 2}, {Height: 3},
	}
	cases := []struct {
		name    string
		indices []int
		want    []uint64
	}{
		{"none re-queued", nil, nil},
		{"middle re-queued", []int{1}, []uint64{2}},
		{"all re-queued preserves order", []int{0, 1, 2}, []uint64{1, 2, 3}},
		{"out-of-range ignored", []int{2, 9, -1}, []uint64{3}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := keepIndexed(events, tc.indices)
			if len(got) != len(tc.want) {
				t.Fatalf("kept %d events, want %d", len(got), len(tc.want))
			}
			for i, h := range tc.want {
				if got[i].Height != h {
					t.Fatalf("kept[%d].Height = %d, want %d", i, got[i].Height, h)
				}
			}
		})
	}
}

// TestL2AckToEvent_EmptyAckSkipped guards the no-ack skip: a WriteAcknowledgement
// with no acknowledgement bytes cannot become a MsgAcknowledgement, so it must be
// dropped rather than relayed.
func TestL2AckToEvent_EmptyAckSkipped(t *testing.T) {
	ev := &contractICS26Router.ContractICS26RouterWriteAcknowledgement{
		Acknowledgements: nil,
	}
	if _, ok := l2AckToEvent(ev, "08-wasm-3"); ok {
		t.Fatal("empty acknowledgement must be skipped")
	}
}

func TestL2AckToEvent_FiltersStaleCosmosClient(t *testing.T) {
	ev := &contractICS26Router.ContractICS26RouterWriteAcknowledgement{
		Sequence:         big.NewInt(1),
		Packet:           contractICS26Router.IICS26RouterMsgsPacket{SourceClient: "08-wasm-1", DestClient: "arb-client-0"},
		Acknowledgements: [][]byte{[]byte(`{"result":"AQ=="}`)},
	}
	if _, ok := l2AckToEvent(ev, "08-wasm-3"); ok {
		t.Fatal("ack for stale Cosmos wasm client must be skipped")
	}
}

func TestL2AckToEvent_AllowsConfiguredCosmosClient(t *testing.T) {
	ev := &contractICS26Router.ContractICS26RouterWriteAcknowledgement{
		Sequence: big.NewInt(2),
		Packet: contractICS26Router.IICS26RouterMsgsPacket{
			SourceClient: "08-wasm-3",
			DestClient:   "arb-client-0",
			Payloads: []contractICS26Router.IICS26RouterMsgsPayload{{
				SourcePort: "transfer",
				DestPort:   "transfer",
				Version:    "ics20-1",
				Encoding:   "application/x-solidity-abi",
				Value:      []byte("packet"),
			}},
		},
		Acknowledgements: [][]byte{[]byte(`{"result":"AQ=="}`)},
	}
	got, ok := l2AckToEvent(ev, "08-wasm-3")
	if !ok {
		t.Fatal("ack for configured Cosmos wasm client must be emitted")
	}
	if got.Type != chain.AckPacket || got.Height != 0 || len(got.AckBytes) != 1 {
		t.Fatalf("unexpected event: %+v", got)
	}
}

// TestScanTerminalLogs_SettlesWithoutRelaying pins the rule that makes terminal
// events safe to read: they close a packet's lifecycle, so they must settle the
// tracker and produce NO relay work. Emitting one as a chain.Event would relay a
// packet with nothing left to do.
func TestScanTerminalLogs_SettlesWithoutRelaying(t *testing.T) {
	t.Parallel()

	t.Run("a source with no hook ignores them", func(t *testing.T) {
		t.Parallel()
		// The behaviour every L2 source had before terminal events were read: no
		// tracker to settle, so nothing to do and nothing to fail on.
		s := &Source{}
		settled, err := s.scanTerminalLogs(nil, nil, nil)
		if err != nil {
			t.Fatalf("a nil settle hook must be a no-op, got %v", err)
		}
		if settled != nil {
			t.Fatalf("a source with no tracker settles nothing, got %v", settled)
		}
	})

	t.Run("the hook receives the marshaled packet", func(t *testing.T) {
		t.Parallel()
		var got [][]byte
		s := (&Source{}).WithSettleHook(func(raw []byte) { got = append(got, raw) })
		if s.settle == nil {
			t.Fatal("WithSettleHook did not install the hook")
		}
		s.settle([]byte("packet"))
		if len(got) != 1 || string(got[0]) != "packet" {
			t.Fatalf("hook got %v, want one marshaled packet", got)
		}
	})
}

// A scan window that contains BOTH a send and its later terminal event is the
// ordinary case after a restart -- the window is sized to cover the downtime.
// Settling only the tracker leaves the send in the batch, so handleBatch relays
// it and TRACKS IT AGAIN, after settle removed it and after the cursor moved
// past the terminal log that would have cleared it. Nothing settles it a second
// time.
func TestDropSettled(t *testing.T) {
	t.Parallel()

	// Raw is the identity: both the send event and the terminal log reach it
	// through the same EthPacketToCosmosPacket + proto.Marshal.
	event := func(seq uint64, client, raw string) chain.Event {
		return chain.Event{Type: chain.SendPacket, Sequence: seq, ClientID: client, Height: 100, Raw: []byte(raw)}
	}
	settledSet := func(raws ...string) map[settledKey]struct{} {
		set := map[settledKey]struct{}{}
		for _, raw := range raws {
			set[settledKey(raw)] = struct{}{}
		}
		return set
	}

	send := event(5, "08-wasm-7", "packet-5")
	other := event(6, "08-wasm-7", "packet-6")
	elsewhere := event(5, "08-wasm-9", "packet-5-other-client")

	t.Run("nothing settled leaves the batch alone", func(t *testing.T) {
		t.Parallel()
		batch := []chain.Event{send, other}
		if got := dropSettled(batch, nil); len(got) != 2 {
			t.Fatalf("kept %d of 2 with an empty settled set", len(got))
		}
	})

	t.Run("a settled packet is removed", func(t *testing.T) {
		t.Parallel()
		got := dropSettled([]chain.Event{send, other}, settledSet("packet-5"))
		if len(got) != 1 || got[0].Sequence != 6 {
			t.Fatalf("kept %+v, want only seq=6", got)
		}
	})

	// A bare sequence would drop an unrelated packet that happens to share a
	// number on another client.
	t.Run("the same sequence on another client survives", func(t *testing.T) {
		t.Parallel()
		got := dropSettled([]chain.Event{send, elsewhere}, settledSet("packet-5"))
		if len(got) != 1 || got[0].ClientID != "08-wasm-9" {
			t.Fatalf("kept %+v, want only the other client's packet", got)
		}
	})

	// And the case client+sequence cannot express at all: a migration keeps the
	// client id and restarts sequences from 1, so seq=5 names two different
	// packets on the SAME client. A stale terminal event for the old one must
	// not drop the new one -- it would never be relayed and nothing would notice.
	t.Run("a replayed sequence on the same client survives", func(t *testing.T) {
		t.Parallel()
		replayed := event(5, "08-wasm-7", "packet-5-after-migration")
		got := dropSettled([]chain.Event{replayed}, settledSet("packet-5"))
		if len(got) != 1 {
			t.Fatalf("kept %+v; the pre-migration seq=5 settled, not this one", got)
		}
	})

	// Both halves of the batch are filtered, not just the fresh ones: a send
	// already in the retry queue can be settled by another relayer.
	t.Run("a requeued packet is removed too", func(t *testing.T) {
		t.Parallel()
		pending := []chain.Event{send}
		fresh := []chain.Event{other}
		got := dropSettled(append(pending[:len(pending):len(pending)], fresh...), settledSet("packet-5"))
		if len(got) != 1 || got[0].Sequence != 6 {
			t.Fatalf("kept %+v; a settled packet in the RETRY queue must go too", got)
		}
	})
}

// TestSettledIdentitySurvivesASequenceReplay drives BOTH ends through production
// code: settleTerminal builds the settled key, l2SendToEvent builds the batch
// event. That is the point -- a test that constructs the key itself can only
// document the choice, not hold it, because changing the key type changes the
// test with it.
//
// A client migration keeps the client id and restarts sequences from 1, so seq=N
// names two different packets on one client. The terminal event for the old one
// must not drop the new one from the batch: it would never be relayed, and the
// cursor has already moved past the log that would have settled it.
func TestSettledIdentitySurvivesASequenceReplay(t *testing.T) {
	t.Parallel()

	const client = "arb-client-0"
	source := &Source{settle: func([]byte) {}}

	old := IICS26RouterMsgsPacketFor("08-wasm-1", client)
	old.TimeoutTimestamp = 1000
	old.Payloads[0].Value = []byte("pre-migration")

	replayed := IICS26RouterMsgsPacketFor("08-wasm-1", client)
	replayed.TimeoutTimestamp = 2000
	replayed.Payloads[0].Value = []byte("post-migration")

	// Both carry sequence 1 -- the migration restarted the numbering.
	settled := map[settledKey]struct{}{}
	source.settleTerminal("AckPacket", old, big.NewInt(1), settled)

	event, ok := l2SendToEvent(&contractICS26Router.ContractICS26RouterSendPacket{
		Sequence: big.NewInt(1), Packet: replayed,
	}, client)
	if !ok {
		t.Fatal("the replayed send must decode; the fixture is wrong otherwise")
	}

	if got := dropSettled([]chain.Event{event}, settled); len(got) != 1 {
		t.Fatal("the pre-migration seq=1 settled, not this one; dropping it loses a packet that was never relayed")
	}

	// And the packet the terminal event really names is still dropped.
	sameEvent, ok := l2SendToEvent(&contractICS26Router.ContractICS26RouterSendPacket{
		Sequence: big.NewInt(1), Packet: old,
	}, client)
	if !ok {
		t.Fatal("the original send must decode")
	}
	if got := dropSettled([]chain.Event{sameEvent}, settled); len(got) != 0 {
		t.Fatalf("kept %+v; the settled packet's own identity must still match", got)
	}
}
