package l2rollup

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

// burstyHeaderServer answers eth_getBlockByNumber with timestamps describing a
// demand-driven chain: a long idle gap, then a recent burst. Anything at or
// above burstFrom is one second apart; the block below it sits idleSeconds
// earlier, which is the shape Nitro takes after an hour of no demand.
func burstyHeaderServer(t *testing.T, burstFrom, idleSeconds uint64) *ethclient.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
			Params []any           `json:"params"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		w.Header().Set("Content-Type", "application/json")
		if req.Method != "eth_getBlockByNumber" {
			_, _ = fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"error":{"code":-32601,"message":"%s"}}`, req.ID, req.Method)
			return
		}
		var number uint64
		if hex, ok := req.Params[0].(string); ok {
			_, _ = fmt.Sscanf(hex, "0x%x", &number)
		}
		// The burst is one second per block; everything below it is pushed back
		// by the idle span, so a span crossing the boundary averages badly.
		ts := number
		if number < burstFrom {
			ts = number - idleSeconds
		}
		_, _ = fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":{`+
			`"number":"0x%x","hash":"0x%064x","parentHash":"0x%064x","timestamp":"0x%x",`+
			`"sha3Uncles":"0x%064x","miner":"0x%040x","stateRoot":"0x%064x",`+
			`"transactionsRoot":"0x%064x","receiptsRoot":"0x%064x","logsBloom":"0x%0512x",`+
			`"difficulty":"0x0","gasLimit":"0x0","gasUsed":"0x0","extraData":"0x",`+
			`"mixHash":"0x%064x","nonce":"0x0000000000000000"}}`,
			req.ID, number, number, number-1, ts, 0, 0, 0, 0, 0, 0, 0)
	}))
	t.Cleanup(server.Close)

	client, err := ethclient.Dial(server.URL)
	if err != nil {
		t.Fatalf("dial stub: %v", err)
	}
	t.Cleanup(client.Close)
	return client
}

// Found in review: two timestamps cannot tell a slow chain from a fast chain
// that was idle. After an hour of no demand the newest blockTimeSamples span
// that hour, the average reads minutes per block, and the startup window
// collapses to a couple of blocks -- while the packets that matter were emitted
// seconds ago during the burst. Nothing persists this cursor, so anything
// outside the window is lost for good.
func TestStartupLookbackSurvivesAnIdleSpanBeforeABurst(t *testing.T) {
	const (
		head        = uint64(10_000)
		idleSeconds = uint64(3600) // an hour of no demand, inside the sample span
	)
	// The burst starts INSIDE the sample span, so the span straddles the gap.
	source := &Source{eth: burstyHeaderServer(t, head-blockTimeSamples/2, idleSeconds)}

	blockTime, err := source.measureBlockTime(context.Background(), head)
	if err != nil {
		t.Fatalf("measureBlockTime: %v", err)
	}
	// The raw measurement is expected to be wrong -- that is the premise. If it
	// is not, this test is not exercising the case it was written for.
	if blockTime <= maxMeasuredBlockTime {
		t.Fatalf("measured %s, which is already within the %s cap; the stub no longer "+
			"produces an idle-inflated average and this test proves nothing",
			blockTime, maxMeasuredBlockTime)
	}
	t.Logf("raw two-sample average = %s (inflated by the idle span)", blockTime)

	got := source.startupLookback(context.Background(), head)
	want := lookbackBlocks(defaultStartupWindow, maxMeasuredBlockTime)
	if got != want {
		t.Fatalf("startupLookback = %d blocks, want %d: an idle span before a burst sized the "+
			"window from a %s average, so a packet emitted seconds ago but %d blocks back falls "+
			"outside the startup scan -- and this cursor is not persisted, so it is never rediscovered",
			got, want, blockTime, got+1)
	}
	t.Logf("capped window = %d blocks", got)
}
