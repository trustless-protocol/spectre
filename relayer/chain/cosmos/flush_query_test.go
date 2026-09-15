package cosmos

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtbytes "github.com/cometbft/cometbft/libs/bytes"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/gogoproto/proto"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"

	"relayer/services"
)

// flushRPC serves the two calls a flush pass makes -- the paginated commitment
// query and the transaction search -- against a synthetic chain, and counts the
// ABCI round-trips so a test can assert what a pass COSTS and not only what it
// returns.
type flushRPC struct {
	t *testing.T

	// commitments is the outstanding set, ascending, exactly as the store holds
	// it: keyed by sequence, valued by the commitment hash.
	commitments []channeltypesv2.PacketState
	// txs are the indexed send_packet transactions, oldest first, which is the
	// order "asc" hands back.
	txs []*coretypes.ResultTx

	abciCalls int
	txCalls   int
}

type flushRPCRequest struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type flushRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
}

func (m *flushRPC) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req flushRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		m.t.Errorf("decode RPC request: %v", err)
		return
	}
	switch req.Method {
	case "abci_query":
		m.abciCalls++
		m.writeResult(w, req.ID, &coretypes.ResultABCIQuery{
			Response: abci.ResponseQuery{Value: m.serveCommitments(req.Params)},
		})
	case "tx_search":
		m.txCalls++
		// per_page is honoured, because it is load-bearing: a search that asks for
		// one result gets only the OLDEST transaction at that sequence, which on a
		// replayed sequence is the stale one. A stub that ignored it would let the
		// commitment check look sufficient on its own.
		//
		// `page` is honoured for the same reason, and it is what the second
		// review round turned up: a stub that always answered page 1 would let a
		// single-page lookup look complete while the live send sat at result 11.
		m.writeResult(w, req.ID, &coretypes.ResultTxSearch{
			Txs: m.pageOf(req.Params), TotalCount: len(m.txs),
		})
	default:
		m.t.Errorf("unexpected RPC method %q", req.Method)
	}
}

// serveCommitments answers one page the way the SDK paginator does: start the
// iterator at Pagination.Key, return Limit entries, and hand back the key of the
// next one. That is the behaviour the window query depends on, so faking it
// loosely would make the test agree with a store that does not exist.
func (m *flushRPC) serveCommitments(params json.RawMessage) []byte {
	var decoded struct {
		Data cmtbytes.HexBytes `json:"data"`
	}
	if err := cmtjson.Unmarshal(params, &decoded); err != nil {
		m.t.Fatalf("decode abci_query params: %v", err)
	}
	var request channeltypesv2.QueryPacketCommitmentsRequest
	if err := proto.Unmarshal(decoded.Data, &request); err != nil {
		m.t.Fatalf("decode commitments request: %v", err)
	}

	start := 0
	if key := request.Pagination.GetKey(); len(key) == 8 {
		from := binary.BigEndian.Uint64(key)
		for start < len(m.commitments) && m.commitments[start].Sequence < from {
			start++
		}
	}
	limit := int(request.Pagination.GetLimit())
	if limit <= 0 {
		limit = 100
	}

	response := channeltypesv2.QueryPacketCommitmentsResponse{Pagination: &query.PageResponse{}}
	for i := start; i < len(m.commitments) && len(response.Commitments) < limit; i++ {
		commitment := m.commitments[i]
		response.Commitments = append(response.Commitments, &commitment)
	}
	if next := start + len(response.Commitments); next < len(m.commitments) {
		response.Pagination.NextKey = sequenceKey(m.commitments[next].Sequence)
	}
	raw, err := proto.Marshal(&response)
	if err != nil {
		m.t.Fatalf("marshal commitments response: %v", err)
	}
	return raw
}

// pageOf slices the indexed transactions the way CometBFT paginates them:
// one-based page, per_page entries, and an empty slice past the end.
func (m *flushRPC) pageOf(params json.RawMessage) []*coretypes.ResultTx {
	perPage := m.perPage(params)
	page := m.page(params)
	start := (page - 1) * perPage
	if start >= len(m.txs) {
		return nil
	}
	return m.txs[start:min(start+perPage, len(m.txs))]
}

// page reads the one-based page number, defaulting the way CometBFT does.
func (m *flushRPC) page(params json.RawMessage) int {
	var decoded struct {
		Page *string `json:"page"`
	}
	if err := cmtjson.Unmarshal(params, &decoded); err != nil {
		m.t.Fatalf("decode tx_search params: %v", err)
	}
	if decoded.Page == nil {
		return 1
	}
	n, err := strconv.Atoi(*decoded.Page)
	if err != nil {
		m.t.Fatalf("page %q: %v", *decoded.Page, err)
	}
	return n
}

// perPage reads the page size out of the tx_search parameters, defaulting the
// way CometBFT does when the caller does not ask.
func (m *flushRPC) perPage(params json.RawMessage) int {
	var decoded struct {
		PerPage *string `json:"per_page"`
	}
	if err := cmtjson.Unmarshal(params, &decoded); err != nil {
		m.t.Fatalf("decode tx_search params: %v", err)
	}
	if decoded.PerPage == nil {
		return 30
	}
	n, err := strconv.Atoi(*decoded.PerPage)
	if err != nil {
		m.t.Fatalf("per_page %q: %v", *decoded.PerPage, err)
	}
	return n
}

func (m *flushRPC) writeResult(w http.ResponseWriter, id json.RawMessage, result any) {
	raw, err := cmtjson.Marshal(result)
	if err != nil {
		m.t.Fatalf("marshal RPC result: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(flushRPCResponse{JSONRPC: "2.0", ID: id, Result: raw}); err != nil {
		m.t.Errorf("encode RPC response: %v", err)
	}
}

func newFlushSource(t *testing.T, rpc *flushRPC) *Source {
	t.Helper()
	server := httptest.NewServer(rpc)
	t.Cleanup(server.Close)
	client, err := rpchttp.New(server.URL, "/websocket")
	if err != nil {
		t.Fatalf("rpc client: %v", err)
	}
	return &Source{
		cosmos: services.CosmosEndpoint{Client: client},
		ids:    services.ClientIDs{CosmosOnEVM: "eth-side-id", EVMOnCosmos: "08-wasm-7"},
		logger: log.New(io.Discard, "", 0),
	}
}

func syntheticCommitments(n int) []channeltypesv2.PacketState {
	out := make([]channeltypesv2.PacketState, 0, n)
	for i := 1; i <= n; i++ {
		packet := channeltypesv2.Packet{
			Sequence: uint64(i), SourceClient: "08-wasm-7", DestinationClient: "eth-side-id",
		}
		out = append(out, channeltypesv2.PacketState{
			ClientId: "08-wasm-7", Sequence: packet.Sequence, Data: channeltypesv2.CommitPacket(packet),
		})
	}
	return out
}

// Reported by @DongLieu on #452 (P2). The pass used to page the WHOLE commitment
// set before capping it, so its query cost was proportional to the backlog: a
// chain with a large outstanding set spent hundreds of serial ABCI round-trips
// and could burn the 60-second pass timeout without relaying anything. The cap
// that was advertised as bounding the work only ever bounded the TxSearch half.
//
// The bound asserted here is the one that matters operationally: the number of
// round-trips must depend on the WINDOW, not on how much is outstanding.
func TestOutstandingCommitmentsCostDoesNotGrowWithTheBacklog(t *testing.T) {
	const backlog = 50_000
	rpc := &flushRPC{t: t, commitments: syntheticCommitments(backlog)}
	source := newFlushSource(t, rpc)

	window, err := source.outstandingCommitments(context.Background())
	if err != nil {
		t.Fatalf("outstanding commitments: %v", err)
	}
	if len(window) != flushSequenceCap {
		t.Fatalf("window = %d commitments, want the cap %d", len(window), flushSequenceCap)
	}

	wantCalls := flushSequenceCap / flushCommitmentPageSize
	if rpc.abciCalls > wantCalls {
		t.Fatalf("a %d-commitment backlog cost %d ABCI queries, want at most %d; "+
			"the query is paging the whole set again and the pass timeout bounds the backstop, not the cap",
			backlog, rpc.abciCalls, wantCalls)
	}
}

// The rotation has to survive moving into the query: passes must advance, and
// the pass that runs out of tail must wrap rather than stop.
func TestOutstandingCommitmentsRotatesAndWraps(t *testing.T) {
	rpc := &flushRPC{t: t, commitments: syntheticCommitments(flushSequenceCap * 2)}
	source := newFlushSource(t, rpc)

	take := func() []flushCommitment {
		window, err := source.outstandingCommitments(context.Background())
		if err != nil {
			t.Fatalf("outstanding commitments: %v", err)
		}
		source.advanceFlushCursor(window)
		return window
	}

	first, second, third := take(), take(), take()

	if first[0].sequence != 1 {
		t.Fatalf("the first pass must start at the oldest, got %d", first[0].sequence)
	}
	if second[0].sequence != first[len(first)-1].sequence+1 {
		t.Fatalf("second pass starts at %d, want %d; later packets would never be considered",
			second[0].sequence, first[len(first)-1].sequence+1)
	}
	if third[0].sequence != 1 {
		t.Fatalf("the pass after the end must wrap to the oldest, got %d", third[0].sequence)
	}
}

// A set that fits inside one window must come back whole on every pass, and the
// cursor must not skip its own head next time.
func TestOutstandingCommitmentsKeepsAShortSetWhole(t *testing.T) {
	rpc := &flushRPC{t: t, commitments: syntheticCommitments(3)}
	source := newFlushSource(t, rpc)

	for pass := 0; pass < 3; pass++ {
		window, err := source.outstandingCommitments(context.Background())
		if err != nil {
			t.Fatalf("pass %d: %v", pass, err)
		}
		if len(window) != 3 || window[0].sequence != 1 {
			t.Fatalf("pass %d returned %d commitments starting at %d, want all 3 starting at 1",
				pass, len(window), window[0].sequence)
		}
		source.advanceFlushCursor(window)
	}
}

// Reported by @DongLieu on #452 (P1). A sequence names a POSITION, not a packet:
// a client migration keeps the client id and restarts sequences, so the index
// holds an old send at the same (client, sequence). Ordering by "asc" hands back
// that old transaction first, and matching on the sequence alone would rebuild a
// packet whose lifecycle closed long ago -- while the packet actually
// outstanding stayed undiscovered, identically, on every pass.
//
// The commitment is what separates them, and it is the same hash the chain would
// check a proof against.
func TestSendEventForCommitmentPicksThePacketTheCommitmentNames(t *testing.T) {
	stale := channeltypesv2.Packet{
		Sequence: 42, SourceClient: "08-wasm-7", DestinationClient: "eth-side-id",
		Payloads: []channeltypesv2.Payload{{SourcePort: "transfer", DestinationPort: "transfer", Value: []byte("old")}},
	}
	live := channeltypesv2.Packet{
		Sequence: 42, SourceClient: "08-wasm-7", DestinationClient: "eth-side-id",
		Payloads: []channeltypesv2.Payload{{SourcePort: "transfer", DestinationPort: "transfer", Value: []byte("new")}},
	}

	rpc := &flushRPC{t: t, txs: []*coretypes.ResultTx{
		// Oldest first, which is what "asc" returns and what a single-result
		// search would have taken.
		txWithSend(11, stale),
		txWithSend(97, live),
	}}
	source := newFlushSource(t, rpc)

	event, found, err := source.sendEventForCommitment(context.Background(), commitmentFor(live))
	if err != nil {
		t.Fatalf("send event: %v", err)
	}
	if !found {
		t.Fatal("the outstanding packet is indexed and must be found")
	}
	if event.Height != 97 {
		t.Fatalf("resolved the send at height %d, want 97: the replayed sequence's OLD transaction was taken, "+
			"so a settled packet is relayed and the outstanding one is never recovered", event.Height)
	}

	// And a commitment naming a packet nothing indexed must report not-found
	// rather than borrowing the nearest send at that sequence.
	absent := live
	absent.Payloads[0].Value = []byte("never sent")
	if _, found, err := source.sendEventForCommitment(context.Background(), commitmentFor(absent)); err != nil || found {
		t.Fatalf("found=%v err=%v; a commitment with no matching send must report not-found", found, err)
	}
}

// Found in review: raising per_page from 1 to 10 defeated a stale entry that
// came back FIRST, but the lookup still read only page one. On the very nodes
// that cross-match -- psql, or kv entries in the legacy format -- the stale
// entries can FILL that page and leave the live send at result 11. The lookup
// then reported not-found, the cursor moved on, and the next pass asked the
// identical question and got the identical answer: outstanding forever, while
// every pass says there is nothing to do.
func TestSendEventForCommitmentPagesPastAFullPageOfStaleMatches(t *testing.T) {
	live := channeltypesv2.Packet{
		Sequence: 42, SourceClient: "08-wasm-7", DestinationClient: "eth-side-id",
		Payloads: []channeltypesv2.Payload{{SourcePort: "transfer", DestinationPort: "transfer", Value: []byte("new")}},
	}

	// Exactly one full page of cross-matches, oldest first, then the packet the
	// commitment actually names.
	txs := make([]*coretypes.ResultTx, 0, flushTxSearchResults+1)
	for i := range flushTxSearchResults {
		stale := live
		stale.Payloads = []channeltypesv2.Payload{{
			SourcePort: "transfer", DestinationPort: "transfer",
			Value: []byte(fmt.Sprintf("stale-%d", i)),
		}}
		txs = append(txs, txWithSend(int64(11+i), stale))
	}
	txs = append(txs, txWithSend(97, live))

	rpc := &flushRPC{t: t, txs: txs}
	source := newFlushSource(t, rpc)

	event, found, err := source.sendEventForCommitment(context.Background(), commitmentFor(live))
	if err != nil {
		t.Fatalf("send event: %v", err)
	}
	if !found {
		t.Fatalf("the outstanding packet is indexed at result %d and must be found; a lookup that "+
			"reads one page leaves it outstanding forever while every pass reports no work",
			len(txs))
	}
	if event.Height != 97 {
		t.Fatalf("resolved height %d, want 97", event.Height)
	}
	if rpc.txCalls < 2 {
		t.Fatalf("tx_search was called %d time(s); the match is on page two, so a single call "+
			"cannot have found it honestly", rpc.txCalls)
	}
}

// The other end of the same loop: a cross-matching index must not be able to
// spin it. The pass is bounded per sequence, and reaching the bound reports
// rather than looking like a clean not-found.
func TestSendEventForCommitmentStopsAtTheScanBound(t *testing.T) {
	live := channeltypesv2.Packet{
		Sequence: 42, SourceClient: "08-wasm-7", DestinationClient: "eth-side-id",
		Payloads: []channeltypesv2.Payload{{SourcePort: "transfer", DestinationPort: "transfer", Value: []byte("new")}},
	}

	// More stale entries than the bound allows, and the live one past all of
	// them, so the lookup must give up rather than page on forever.
	total := flushTxSearchMaxResults + flushTxSearchResults*3
	txs := make([]*coretypes.ResultTx, 0, total+1)
	for i := range total {
		stale := live
		stale.Payloads = []channeltypesv2.Payload{{
			SourcePort: "transfer", DestinationPort: "transfer",
			Value: []byte(fmt.Sprintf("stale-%d", i)),
		}}
		txs = append(txs, txWithSend(int64(11+i), stale))
	}
	txs = append(txs, txWithSend(9_000, live))

	rpc := &flushRPC{t: t, txs: txs}
	source := newFlushSource(t, rpc)

	if _, found, err := source.sendEventForCommitment(context.Background(), commitmentFor(live)); err != nil || found {
		t.Fatalf("found=%v err=%v; past the scan bound the lookup must give up", found, err)
	}
	maxCalls := flushTxSearchMaxResults/flushTxSearchResults + 1
	if rpc.txCalls > maxCalls {
		t.Fatalf("tx_search was called %d times, want at most %d: one sequence must not be able "+
			"to page without bound, because flushSequenceCap of them run per pass",
			rpc.txCalls, maxCalls)
	}
}

func txWithSend(height int64, packet channeltypesv2.Packet) *coretypes.ResultTx {
	return &coretypes.ResultTx{
		Height: height,
		Tx:     cmttypes.Tx{},
		TxResult: abci.ExecTxResult{
			Events: []abci.Event{{Type: "send_packet", Attributes: encodedSendAttributes(packet)}},
		},
	}
}
