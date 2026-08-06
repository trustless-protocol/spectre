package l2rollup

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"relayer/chain"

	contractICS26Router "relayer/bindings/ICS26Router"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	idleTestClientID = "arb-client-0"
	idleTestDestID   = "08-wasm-1"
	idleTestHead     = uint64(100)
)

// idleL2Node is a JSON-RPC stub for a demand-driven rollup that has gone quiet:
// it serves one SendPacket log on the first eth_getLogs call and nothing after,
// while eth_getBlockByNumber reports the same head forever. Arbitrum Nitro seals a
// block per transaction, so "no traffic" really does mean "the head stops" — this
// is the steady state, not an edge case.
type idleL2Node struct {
	router ethcommon.Address

	// logBlock is where the single SendPacket log sits; zero means idleTestHead.
	logBlock uint64

	mu sync.Mutex
	// head is the block eth_getBlockByNumber reports; zero means idleTestHead. It is
	// settable so a test can drive an unsafe-head reorg — the head moving BACKWARDS —
	// rather than only a head that stops advancing.
	head      uint64
	logCalls  int
	headCalls int
	// failHeads makes the first N head reads answer with a JSON-RPC error, standing
	// in for the rate limit a public endpoint returns during a restart.
	failHeads int
	// logRanges records the [fromBlock,toBlock] of every eth_getLogs, so a test can
	// assert the cursor neither drifted nor rewound across a reorg.
	logRanges [][2]uint64
}

func (n *idleL2Node) setHead(h uint64) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.head = h
}

func (n *idleL2Node) scannedRanges() [][2]uint64 {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([][2]uint64(nil), n.logRanges...)
}

func (n *idleL2Node) start(t *testing.T) *ethclient.Client {
	t.Helper()
	logBlock := n.logBlock
	if logBlock == 0 {
		logBlock = idleTestHead
	}
	sendLog := sendPacketLogJSON(t, n.router, logBlock)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     json.RawMessage   `json:"id"`
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		var result string
		switch req.Method {
		case "eth_getBlockByNumber":
			n.mu.Lock()
			n.headCalls++
			calls := n.headCalls
			failing := n.failHeads
			head := n.head
			n.mu.Unlock()
			if calls <= failing {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"error":{"code":-32005,`+
					`"message":"You reached Public endpoint rate limit"}}`, req.ID)
				return
			}
			if head == 0 {
				head = idleTestHead
			}
			result = idleHeaderJSON(head)
		case "eth_getLogs":
			from, to := idleFilterRange(t, req.Params)
			n.mu.Lock()
			n.logCalls++
			first := n.logCalls == 1
			n.logRanges = append(n.logRanges, [2]uint64{from, to})
			n.mu.Unlock()
			if first {
				result = "[" + sendLog + "]" // the only packet this L2 ever emits
			} else {
				result = "[]"
			}
		case "eth_chainId":
			result = `"0x1"`
		default:
			result = "null"
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":%s}`, req.ID, result)
	}))
	t.Cleanup(srv.Close)

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("dial stub L2: %v", err)
	}
	t.Cleanup(client.Close)
	return client
}

func (n *idleL2Node) counts() (logs, heads int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.logCalls, n.headCalls
}

// idleFilterRange pulls the block range out of an eth_getLogs filter so a test can
// assert which range was scanned, not merely that a scan happened.
func idleFilterRange(t *testing.T, params []json.RawMessage) (uint64, uint64) {
	t.Helper()
	if len(params) == 0 {
		t.Fatalf("eth_getLogs called with no params")
	}
	var filter struct {
		FromBlock string `json:"fromBlock"`
		ToBlock   string `json:"toBlock"`
	}
	if err := json.Unmarshal(params[0], &filter); err != nil {
		t.Fatalf("decode eth_getLogs filter: %v", err)
	}
	parse := func(field, raw string) uint64 {
		v, err := strconv.ParseUint(strings.TrimPrefix(raw, "0x"), 16, 64)
		if err != nil {
			t.Fatalf("parse %s %q: %v", field, raw, err)
		}
		return v
	}
	return parse("fromBlock", filter.FromBlock), parse("toBlock", filter.ToBlock)
}

// sendPacketLogJSON ABI-encodes a real SendPacket log so the production filterer
// decodes it — a hand-rolled log would only prove the stub matches itself.
func sendPacketLogJSON(t *testing.T, router ethcommon.Address, block uint64) string {
	t.Helper()
	parsed, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("parse ICS26Router ABI: %v", err)
	}
	event, ok := parsed.Events["SendPacket"]
	if !ok {
		t.Fatal("ICS26Router ABI has no SendPacket event")
	}
	packet := IICS26RouterMsgsPacketFor(idleTestClientID, idleTestDestID)
	data, err := event.Inputs.NonIndexed().Pack(packet)
	if err != nil {
		t.Fatalf("pack SendPacket data: %v", err)
	}
	// Indexed string params are stored as the keccak of their contents.
	clientTopic := crypto.Keccak256Hash([]byte(idleTestClientID))
	seqTopic := ethcommon.BigToHash(big.NewInt(1))
	return fmt.Sprintf(`{
		"address":%q,"topics":[%q,%q,%q],"data":"0x%x",
		"blockNumber":"0x%x","transactionHash":%q,"transactionIndex":"0x0",
		"blockHash":%q,"logIndex":"0x0","removed":false
	}`, router.Hex(), event.ID.Hex(), clientTopic.Hex(), seqTopic.Hex(), data,
		block, ethcommon.Hash{}.Hex(), ethcommon.Hash{}.Hex())
}

// IICS26RouterMsgsPacketFor builds the minimal packet the SendPacket decoder accepts.
func IICS26RouterMsgsPacketFor(sourceClient, destClient string) contractICS26Router.IICS26RouterMsgsPacket {
	return contractICS26Router.IICS26RouterMsgsPacket{
		Sequence:         1,
		SourceClient:     sourceClient,
		DestClient:       destClient,
		TimeoutTimestamp: uint64(time.Now().Add(time.Hour).Unix()),
		Payloads: []contractICS26Router.IICS26RouterMsgsPayload{{
			SourcePort: "transfer", DestPort: "transfer",
			Version: "ics20-2", Encoding: "application/x-solidity-abi",
			Value: []byte{0x01},
		}},
	}
}

// idleHeaderJSON is the minimum header go-ethereum decodes into types.Header.
func idleHeaderJSON(number uint64) string {
	zero := ethcommon.Hash{}.Hex()
	return fmt.Sprintf(`{
		"parentHash":%q,"sha3Uncles":%q,"miner":"0x0000000000000000000000000000000000000000",
		"stateRoot":%q,"transactionsRoot":%q,"receiptsRoot":%q,
		"logsBloom":"0x%0512x","difficulty":"0x1","number":"0x%x",
		"gasLimit":"0x1","gasUsed":"0x0","timestamp":"0x1","extraData":"0x",
		"mixHash":%q,"nonce":"0x0000000000000000","baseFeePerGas":"0x1","hash":%q
	}`, zero, zero, zero, zero, zero, 0, number, zero, zero)
}

// TestSubscribeRetriesPendingWhileL2IsIdle is the regression test for a stall that
// silently stranded packets — in practice, acknowledgements.
//
// After a successful scan the cursor moves to head+1, so on a chain that has stopped
// producing blocks every later tick sees head < from. The loop used to `continue` on
// that condition, skipping the handler call entirely — and the handler call is the
// only thing that re-offers `pending`, the events the handler asked to retry because
// they were not yet relayable.
//
// The acknowledgement is the worst case: it is normally the last event a packet
// produces on the L2, so nothing further is written, no new block ever arrives, and
// the re-queued ack is never retried. The packet stays committed on Cosmos and its
// escrow stays locked, with a single "waiting" line in the log and then silence.
// Observed on a Cosmos↔Arbitrum devnet: the ack sat for 15 minutes while the attestor
// frontier passed its height, and only a relayer restart (which resets the cursor)
// delivered it.
//
// The handler here always re-queues, mimicking a packet whose attested height has not
// been reached. The test asserts the handler keeps being called even though the head
// never moves.
func TestSubscribeRetriesPendingWhileL2IsIdle(t *testing.T) {
	restore := l2SubscribeInterval
	l2SubscribeInterval = 20 * time.Millisecond
	t.Cleanup(func() { l2SubscribeInterval = restore })

	router := ethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	node := &idleL2Node{router: router}
	src := &Source{
		eth:                node.start(t),
		headKind:           Unsafe,
		l2ClientID:         idleTestClientID,
		cosmosWasmClientID: idleTestDestID,
		router:             router,
	}

	var (
		mu       sync.Mutex
		calls    int
		lastSize int
	)
	handler := func(_ context.Context, batch []chain.Event) []int {
		mu.Lock()
		calls++
		lastSize = len(batch)
		mu.Unlock()
		requeue := make([]int, len(batch)) // never relayable in this test
		for i := range batch {
			requeue[i] = i
		}
		return requeue
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- src.Subscribe(ctx, handler) }()

	// Give the loop many poll periods. The head never advances after the first scan,
	// so a loop that gates the handler on new blocks calls it exactly once.
	time.Sleep(600 * time.Millisecond)
	cancel()
	<-done

	mu.Lock()
	gotCalls, gotSize := calls, lastSize
	mu.Unlock()

	if gotCalls < 2 {
		t.Fatalf("handler called %d time(s) on an idle L2: the re-queued packet is never retried, "+
			"so an acknowledgement waiting on finality is stranded until the relayer restarts", gotCalls)
	}
	if gotSize != 1 {
		t.Fatalf("last batch carried %d event(s), want the 1 re-queued packet", gotSize)
	}
	logs, heads := node.counts()
	if heads < 2 {
		t.Fatalf("head polled %d time(s); the loop should keep polling", heads)
	}
	if logs == 0 {
		t.Fatal("expected the first tick to scan for logs")
	}
}

// TestSubscribeSurvivesUnsafeHeadReorg covers the other way `head < from` happens:
// with head_kind=unsafe the head can move BACKWARDS when the sequencer reorgs, not
// merely stop advancing. The two share a branch, but the recovery behaviour is
// distinct and worth pinning:
//
//   - while the head is behind the cursor, no scan is issued at all (an eth_getLogs
//     with fromBlock > toBlock is a malformed request, not an empty result), yet the
//     pending buffer must still be re-offered;
//   - when the head climbs back, scanning must resume at exactly the cursor the loop
//     already reached — not rewound (that is the documented forward-only limitation)
//     and not skipped past (which would silently drop the blocks in between).
func TestSubscribeSurvivesUnsafeHeadReorg(t *testing.T) {
	restore := l2SubscribeInterval
	l2SubscribeInterval = 20 * time.Millisecond
	t.Cleanup(func() { l2SubscribeInterval = restore })

	// Above l2StartupLookback so the initial cursor is a real number (head-256),
	// which is what makes "did the cursor move?" observable at all.
	const (
		startHead = uint64(1000)
		reorgHead = uint64(900) // sequencer drops 100 blocks
		afterHead = uint64(1010)
	)
	wantCursor := startHead + 1 // where the loop must resume after the reorg

	router := ethcommon.HexToAddress("0x2222222222222222222222222222222222222222")
	node := &idleL2Node{router: router, logBlock: startHead}
	node.setHead(startHead)
	src := &Source{
		eth:                node.start(t),
		headKind:           Unsafe,
		l2ClientID:         idleTestClientID,
		cosmosWasmClientID: idleTestDestID,
		router:             router,
	}

	var (
		mu            sync.Mutex
		calls         int
		callsAtReorg  int
		lastBatchSize int
	)
	handler := func(_ context.Context, batch []chain.Event) []int {
		mu.Lock()
		calls++
		lastBatchSize = len(batch)
		mu.Unlock()
		requeue := make([]int, len(batch)) // never relayable in this test
		for i := range batch {
			requeue[i] = i
		}
		return requeue
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- src.Subscribe(ctx, handler) }()

	// Phase 1: the head is ahead, the first scan picks the packet up.
	time.Sleep(200 * time.Millisecond)

	// Phase 2: reorg — the head falls behind the cursor.
	mu.Lock()
	callsAtReorg = calls
	mu.Unlock()
	node.setHead(reorgHead)
	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	callsDuringReorg := calls
	mu.Unlock()
	if callsDuringReorg <= callsAtReorg {
		t.Fatalf("handler stopped being called during the reorg (%d -> %d): the pending packet "+
			"is stranded exactly as it was on an idle chain", callsAtReorg, callsDuringReorg)
	}

	// Phase 3: the head climbs back past the cursor.
	node.setHead(afterHead)
	time.Sleep(300 * time.Millisecond)
	cancel()
	<-done

	mu.Lock()
	gotSize := lastBatchSize
	mu.Unlock()
	if gotSize != 1 {
		t.Fatalf("last batch carried %d event(s), want the 1 re-queued packet", gotSize)
	}

	ranges := node.scannedRanges()
	if len(ranges) < 2 {
		t.Fatalf("scanned ranges %v: expected a scan before the reorg and another after recovery", ranges)
	}
	// No scan may be issued while the head is behind the cursor.
	for _, r := range ranges {
		if r[0] > r[1] {
			t.Errorf("scanned an inverted range [%d,%d] — the head was behind the cursor and no "+
				"scan should have been issued", r[0], r[1])
		}
	}
	// The recovery scan resumes at the cursor the loop had already reached.
	last := ranges[len(ranges)-1]
	if last[0] != wantCursor {
		t.Errorf("recovery scan started at block %d, want %d: the cursor %s",
			last[0], wantCursor,
			map[bool]string{true: "rewound over already-scanned blocks", false: "skipped blocks"}[last[0] < wantCursor])
	}
	if last[1] != afterHead {
		t.Errorf("recovery scan ended at block %d, want the new head %d", last[1], afterHead)
	}
}

// A transient RPC error on the FIRST head read used to be fatal while the identical
// error one tick later was logged and retried: the read sat before the loop and its
// error propagated out of Subscribe, up through the relay module, into logger.Fatal.
// A rate-limit blip during a restart therefore stopped the process from coming up at
// all, instead of it starting and recovering on the next tick.
//
// The stub rejects the first two head reads the way a public endpoint rate-limits,
// then serves normally. Subscribe must survive that and still deliver the packet.
func TestSubscribeSurvivesTransientHeadErrorAtStartup(t *testing.T) {
	restore := l2SubscribeInterval
	l2SubscribeInterval = 20 * time.Millisecond
	t.Cleanup(func() { l2SubscribeInterval = restore })

	// A head well above the lookback, so the seeded cursor is a distinctive number
	// rather than the 0 that "head below the window" and "never seeded" share.
	const startHead = uint64(1000)

	router := ethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	node := &idleL2Node{router: router, failHeads: 2, logBlock: startHead}
	node.setHead(startHead)
	src := &Source{
		eth:                node.start(t),
		headKind:           Unsafe,
		l2ClientID:         idleTestClientID,
		cosmosWasmClientID: idleTestDestID,
		router:             router,
	}

	delivered := make(chan int, 1)
	handler := func(_ context.Context, batch []chain.Event) []int {
		select {
		case delivered <- len(batch):
		default:
		}
		return nil // relayed
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	errCh := make(chan error, 1)
	go func() { errCh <- src.Subscribe(ctx, handler) }()

	select {
	case n := <-delivered:
		if n == 0 {
			t.Fatal("handler ran with an empty batch")
		}
	case err := <-errCh:
		t.Fatalf("Subscribe returned before delivering anything: %v", err)
	case <-ctx.Done():
		t.Fatal("packet never delivered after the startup head errors cleared")
	}

	// The cursor must be seeded from the head that finally answered, not from zero:
	// a scan starting at block 0 would hammer the endpoint that just rate-limited us.
	cancel()
	<-errCh
	ranges := node.scannedRanges()
	if len(ranges) == 0 {
		t.Fatal("no eth_getLogs range recorded")
	}
	if wantFrom := startHead - l2StartupLookback; ranges[0][0] != wantFrom {
		t.Fatalf("first scan started at %d, want %d (head - lookback)", ranges[0][0], wantFrom)
	}
}
