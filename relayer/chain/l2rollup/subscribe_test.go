package l2rollup

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"relayer/chain"

	contractICS26Router "relayer/bindings/ICS26Router"

	"github.com/ethereum/go-ethereum/ethclient"
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

func TestLookbackBlocks(t *testing.T) {
	t.Parallel()

	window := 512 * time.Second
	for _, tc := range []struct {
		chain     string
		blockTime time.Duration
		want      uint64
	}{
		{chain: "OP Stack", blockTime: 2 * time.Second, want: 256},
		{chain: "Base", blockTime: 2 * time.Second, want: 256},
		{chain: "Arbitrum Nitro", blockTime: 250 * time.Millisecond, want: 2048},
		{chain: "Ethereum L1", blockTime: 12 * time.Second, want: 43},
	} {
		t.Run(tc.chain, func(t *testing.T) {
			t.Parallel()
			if got := lookbackBlocks(window, tc.blockTime); got != tc.want {
				t.Fatalf("lookbackBlocks(%s, %s) = %d, want %d", window, tc.blockTime, got, tc.want)
			}
		})
	}

	t.Run("Arbitrum and OP differ for the same window", func(t *testing.T) {
		t.Parallel()
		op := lookbackBlocks(window, 2*time.Second)
		arb := lookbackBlocks(window, 250*time.Millisecond)
		if op == arb {
			t.Fatalf("both chains got %d blocks; a block count that ignores block time is the bug B2 removes", op)
		}
		if arb <= op {
			t.Fatalf("Arbitrum lookback %d <= OP lookback %d; the faster chain needs MORE blocks for the same time", arb, op)
		}
	})

	t.Run("rounds up so the window is never short", func(t *testing.T) {
		t.Parallel()
		// 10s at 3s per block is 3.33 blocks; 3 would cover only 9 seconds.
		if got := lookbackBlocks(10*time.Second, 3*time.Second); got != 4 {
			t.Fatalf("lookbackBlocks(10s, 3s) = %d, want 4 (rounded up)", got)
		}
	})

	t.Run("an unusable block time yields no lookback rather than dividing by zero", func(t *testing.T) {
		t.Parallel()
		for _, blockTime := range []time.Duration{0, -time.Second} {
			if got := lookbackBlocks(window, blockTime); got != 0 {
				t.Fatalf("lookbackBlocks(%s, %s) = %d, want 0", window, blockTime, got)
			}
		}
	})
}

// TestStartupWindow pins the window as a SAFETY property. It is the entire
// crash-recovery window for the L2 path -- the cursor is not persisted, so a
// packet emitted before the window is never rescanned and its escrow stays
// locked. Shrinking it must not reach main silently.

func TestStartupWindow(t *testing.T) {
	t.Parallel()

	t.Run("the default reproduces the window this code already had", func(t *testing.T) {
		t.Parallel()
		// 256 blocks at OP's 2s block time, now stated in time so every chain gets
		// the same amount of it. If this changes, re-derive it and say why.
		if got := startupWindow(0); got != 512*time.Second {
			t.Fatalf("startupWindow(0) = %s, want 8m32s", got)
		}
	})

	// The window is a downtime allowance, and head_kind must not widen it. The
	// scan is anchored at head(head_kind) -- Source.head reads SafeBlockNumber or
	// FinalizedBlockNumber -- so both ends of "head minus window" sit in the same
	// delayed view and the finality lag is already absorbed. Adding it again would
	// make every finalized run pay for a case that is not a crash.
	t.Run("head kind does not widen the window", func(t *testing.T) {
		t.Parallel()
		if got := startupWindow(0); got != defaultStartupWindow {
			t.Fatalf("startupWindow(0) = %s, want %s", got, defaultStartupWindow)
		}
	})

	t.Run("an operator override replaces the baseline", func(t *testing.T) {
		t.Parallel()
		if got := startupWindow(2 * time.Hour); got != 2*time.Hour {
			t.Fatalf("startupWindow(2h) = %s, want 2h", got)
		}
	})
}

// TestStartupWindowFromEnv covers the fallback, not the size of the window:
// an absent or unparseable override must leave the default in place rather than
// turning a typo into a chain that starts at the head with no recovery.

func TestStartupWindowFromEnv(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		raw  string
		want time.Duration
	}{
		{name: "empty uses the default", raw: "", want: 0},
		{name: "valid duration overrides", raw: "90m", want: 90 * time.Minute},
		{name: "seconds are a duration too", raw: "45s", want: 45 * time.Second},
		{name: "a bare block count is no longer valid", raw: "1200", want: 0},
		{name: "unparseable uses the default", raw: "not-a-duration", want: 0},
		{name: "negative uses the default", raw: "-5m", want: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := startupWindowFromEnv(tc.raw); got != tc.want {
				t.Fatalf("startupWindowFromEnv(%q) = %s, want %s", tc.raw, got, tc.want)
			}
		})
	}
}

// The retired key changed units as well as name -- block count to duration --
// so an old value cannot be carried forward even in principle. Silently ignoring
// it leaves the operator on the default window and nothing says so, which is the
// event-loss shape the window exists to prevent.

func TestRejectRetiredLookbackEnv(t *testing.T) {
	t.Run("unset is fine", func(t *testing.T) {
		t.Setenv(l2StartupLookbackBlocksEnv, "")
		os.Unsetenv(l2StartupLookbackBlocksEnv)
		if err := RejectRetiredLookbackEnv(); err != nil {
			t.Fatalf("an unset key must not fail startup: %v", err)
		}
	})

	// Even an empty value is a deliberate export, and it still means the operator
	// believes that key does something.
	for _, raw := range []string{"256", "1024", ""} {
		t.Run("set to "+raw+" is refused", func(t *testing.T) {
			t.Setenv(l2StartupLookbackBlocksEnv, raw)
			err := RejectRetiredLookbackEnv()
			if err == nil {
				t.Fatal("a retired key that sizes nothing must fail startup, not be ignored")
			}
			// The message has to carry the operator from what they set to what
			// they should set; naming only the dead key leaves them guessing.
			for _, want := range []string{l2StartupLookbackBlocksEnv, l2StartupLookbackEnv, "DURATION"} {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("error %q does not mention %q", err, want)
				}
			}
		})
	}
}

// measureBlockTime runs on the subscribe goroutine before the cursor is seeded,
// so an endpoint that accepts the connection and never answers must not hold the
// L2 direction open indefinitely. A ws:// endpoint has no transport timeout at
// all, so this deadline is the only bound.

func TestMeasureBlockTimeIsBounded(t *testing.T) {
	// Order matters: defers run LIFO, and httptest.Server.Close waits for its
	// handlers. Registering Close first means the handler is released before
	// Close runs -- the other order deadlocks the test rather than failing it.
	blocked := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-blocked // accept, then never answer
	}))
	defer server.Close()
	defer close(blocked)

	client, err := ethclient.Dial(server.URL)
	if err != nil {
		t.Fatalf("dial stub: %v", err)
	}
	defer client.Close()

	source := &Source{eth: client}
	start := time.Now()
	// The subscribe lifetime context: nothing here cancels it, so only the
	// per-call deadline can end this.
	if _, err := source.measureBlockTime(context.Background(), 1_000); err == nil {
		t.Fatal("a hung endpoint must surface as an error, not a completed measurement")
	}
	if elapsed := time.Since(start); elapsed > blockTimeSampleTimeout+5*time.Second {
		t.Fatalf("measureBlockTime took %s; it must be bounded by %s", elapsed, blockTimeSampleTimeout)
	}
}
