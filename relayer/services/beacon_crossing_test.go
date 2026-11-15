package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	relayerclient "relayer/client"
	"relayer/rpcmock"

	abci "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
)

// Geometry for these tests: 8 slots per epoch, 4 epochs per period, so a
// sync-committee period is 32 slots.
func crossingTestState(trustedSlot uint64) *relayerclient.EthereumClientState {
	return &relayerclient.EthereumClientState{
		LatestSlot: trustedSlot, LatestExecutionBlockNumber: 456,
		SlotsPerEpoch: 8, EpochsPerSyncCommitteePeriod: 4, SecondsPerSlot: 1, SyncCommitteeSize: 8,
	}
}

const testCommittee = `{"pubkeys":[],"aggregate_pubkey":"0x00"}`

// updateJSON is one light_client/updates entry.
func updateJSON(finalized, signature, block uint64, withCommittee bool) string {
	committee := ""
	if withCommittee {
		committee = `,"next_sync_committee":` + testCommittee
	}
	return fmt.Sprintf(`{"data":{"finalized_header":{"beacon":{"slot":"%d"},"execution":{"block_number":"%d","timestamp":"%d"}},"signature_slot":"%d"%s}}`,
		finalized, block, finalized, signature, committee)
}

func finalityJSON(attested, finalized, signature, block uint64) string {
	return fmt.Sprintf(`{"data":{"attested_header":{"beacon":{"slot":"%d"}},"finalized_header":{"beacon":{"slot":"%d"},"execution":{"block_number":"%d","timestamp":"%d"}},"signature_slot":"%d","sync_aggregate":{"sync_committee_bits":"0xff"}}}`,
		attested, finalized, block, finalized, signature)
}

// crossingBeacon serves the four endpoints the builder reads and records the
// slots the bootstrap fallback asks a block root for.
type crossingBeacon struct {
	finality  string
	updates   func(start, count string) (string, int) // body, status
	bootstrap bool                                    // serve a bootstrap for any root
	// checkpoints, when set, are the only slots whose roots get a bootstrap, as
	// on a node that keeps bootstraps for finalized checkpoints alone.
	checkpoints map[string]bool

	mu         sync.Mutex
	rootsAsked []string
}

// rootOf is the block root the fake answers for a slot.
func rootOf(slot string) string { return fmt.Sprintf("0x%064s", slot) }

func (b *crossingBeacon) start(t *testing.T) string {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/eth/v1/beacon/light_client/finality_update":
			fmt.Fprint(w, b.finality)
		case path == "/eth/v1/beacon/light_client/updates":
			body, status := b.updates(r.URL.Query().Get("start_period"), r.URL.Query().Get("count"))
			w.WriteHeader(status)
			fmt.Fprint(w, body)
		case strings.HasPrefix(path, "/eth/v1/beacon/blocks/") && strings.HasSuffix(path, "/root"):
			slot := strings.TrimSuffix(strings.TrimPrefix(path, "/eth/v1/beacon/blocks/"), "/root")
			b.mu.Lock()
			b.rootsAsked = append(b.rootsAsked, slot)
			b.mu.Unlock()
			fmt.Fprintf(w, `{"data":{"root":"%s"}}`, rootOf(slot))
		case strings.HasPrefix(path, "/eth/v1/beacon/light_client/bootstrap/"):
			root := strings.TrimPrefix(path, "/eth/v1/beacon/light_client/bootstrap/")
			servable := b.bootstrap
			for slot := range b.checkpoints {
				servable = servable || root == rootOf(slot)
			}
			if !servable {
				// Lighthouse's answer for a block that is not a finalized checkpoint.
				http.Error(w, "NOT_FOUND: Sync committee branch for block root "+root+
					" not found. This typically occurs when the block is not a finalized checkpoint.", http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, `{"data":{"current_sync_committee":%s}}`, testCommittee)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

func (b *crossingBeacon) asked() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]string(nil), b.rootsAsked...)
}

func decodeHeaders(t *testing.T, headers [][]byte) []relayerclient.EthereumHeader {
	t.Helper()
	out := make([]relayerclient.EthereumHeader, len(headers))
	for i, raw := range headers {
		if err := json.Unmarshal(raw, &out[i]); err != nil {
			t.Fatalf("decode header %d: %v", i, err)
		}
	}
	return out
}

func buildCrossing(t *testing.T, beaconURL string, state *relayerclient.EthereumClientState, finalizedSlot uint64, finality string) (ethUpdateBuild, error) {
	t.Helper()
	var update relayerclient.LightClientFinalityUpdate
	var envelope struct {
		Data relayerclient.LightClientFinalityUpdate `json:"data"`
	}
	if err := json.Unmarshal([]byte(finality), &envelope); err != nil {
		t.Fatal(err)
	}
	update = envelope.Data
	trustedPeriod := state.ComputeSyncCommitteePeriodAtSlot(state.LatestSlot)
	targetPeriod := state.ComputeSyncCommitteePeriodAtSlot(finalizedSlot)
	return (&Worker{}).buildEthClientUpdateHeadersWithPeriodCrossing(context.Background(), beaconURL, nil, state,
		state.LatestSlot, trustedPeriod, targetPeriod, &update, finalizedSlot)
}

// #575: the builder used light_client/updates in the order the node answered
// and never checked that each header steps one period. A node answering out of
// order made it cross from period 0 straight to 2, which the client rejects.
func TestBuildEthClientUpdateHeaders_SortsUpdatesAndStepsOnePeriod(t *testing.T) {
	finality := finalityJSON(81, 80, 82, 900)
	beacon := &crossingBeacon{
		finality: finality,
		updates: func(string, string) (string, int) {
			// Newest first: the order a non-conforming node might answer in.
			return "[" + updateJSON(70, 72, 800, false) + "," + updateJSON(40, 42, 700, true) + "," + updateJSON(16, 18, 456, true) + "]", http.StatusOK
		},
	}
	build, err := buildCrossing(t, beacon.start(t), crossingTestState(16), 80, finality)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	headers := decodeHeaders(t, build.headers)
	var got []string
	for _, h := range headers {
		got = append(got, fmt.Sprintf("%d->%s", h.TrustedSlot, h.ConsensusUpdate.FinalizedHeader.Beacon.Slot))
	}
	if want := "16->40 40->70 70->80"; strings.Join(got, " ") != want {
		t.Fatalf("headers = %v, want %s: one period per step, in slot order", got, want)
	}
}

func TestBuildEthClientUpdateHeaders_RefusesAPeriodGap(t *testing.T) {
	finality := finalityJSON(81, 80, 82, 900)
	beacon := &crossingBeacon{
		finality: finality,
		updates: func(string, string) (string, int) {
			// The period 1 update is missing.
			return "[" + updateJSON(16, 18, 456, true) + "," + updateJSON(70, 72, 800, true) + "]", http.StatusOK
		},
	}
	_, err := buildCrossing(t, beacon.start(t), crossingTestState(16), 80, finality)
	if err == nil || !strings.Contains(err.Error(), "skip from period 0 to 2") {
		t.Fatalf("err = %v, want the missing period named", err)
	}
}

// A crossing header must be signed in the period it finalizes in: the client
// verifies it with the committee it stores as next.
func TestBuildEthClientUpdateHeaders_RefusesASignatureInALaterPeriod(t *testing.T) {
	finality := finalityJSON(81, 80, 82, 900)
	beacon := &crossingBeacon{
		finality: finality,
		updates: func(string, string) (string, int) {
			return "[" + updateJSON(16, 18, 456, true) + "," + updateJSON(60, 64, 700, true) + "]", http.StatusOK
		},
	}
	_, err := buildCrossing(t, beacon.start(t), crossingTestState(16), 80, finality)
	if err == nil || !strings.Contains(err.Error(), "signed at slot 64 (period 2)") {
		t.Fatalf("err = %v, want the signature period refused", err)
	}
}

// The bootstrap fallback reads the committee of the block it is asked for. At a
// boundary the attested slot is still in the old period and the signature slot
// in the new one, and the client checks the signature slot's committee.
func TestBuildEthClientUpdateHeaders_BootstrapFallbackUsesTheSignatureSlot(t *testing.T) {
	// Trusted in period 1; finality signed at slot 64, the first slot of period 2.
	finality := finalityJSON(63, 50, 64, 900)
	beacon := &crossingBeacon{
		finality:  finality,
		bootstrap: true,
		updates: func(start, count string) (string, int) {
			// Neither the fetched range nor the one-period lookup has a committee.
			return "[" + updateJSON(40, 42, 700, false) + "]", http.StatusOK
		},
	}
	url := beacon.start(t)
	if _, err := buildCrossing(t, url, crossingTestState(40), 50, finality); err != nil {
		t.Fatalf("build: %v", err)
	}
	if asked := beacon.asked(); len(asked) != 1 || asked[0] != "64" {
		t.Fatalf("bootstrap block roots asked = %v, want [64], the signature slot", asked)
	}
}

// #575 review: the crossing header's fallback was moved to the signature slot,
// whose block is not a finalized checkpoint, so a node that serves bootstraps
// for checkpoints alone answered 404 where the finalized slot worked.
func TestBuildEthClientUpdateHeaders_CrossingFallbackAsksTheFinalizedCheckpoint(t *testing.T) {
	// Trusted in period 0; the period 1 update finalizes at 40, signed at 42.
	finality := finalityJSON(41, 40, 42, 700)
	beacon := &crossingBeacon{
		finality:    finality,
		checkpoints: map[string]bool{"40": true},
		updates: func(string, string) (string, int) {
			// No update carries next_sync_committee: the fallback has to answer.
			return "[" + updateJSON(40, 42, 700, false) + "]", http.StatusOK
		},
	}
	url := beacon.start(t)
	build, err := buildCrossing(t, url, crossingTestState(16), 40, finality)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(build.headers) != 1 {
		t.Fatalf("headers = %d, want the crossing header alone", len(build.headers))
	}
	if asked := beacon.asked(); len(asked) != 1 || asked[0] != "40" {
		t.Fatalf("bootstrap block roots asked = %v, want [40], the finalized checkpoint", asked)
	}
}

// When the one-period lookup fails, the error used to be dropped and only the
// bootstrap failure was reported.
func TestSyncCommitteeForPeriod_KeepsTheLookupError(t *testing.T) {
	beacon := &crossingBeacon{
		updates: func(string, string) (string, int) { return "upstream unavailable", http.StatusInternalServerError },
	}
	_, err := syncCommitteeForPeriod(context.Background(), beacon.start(t), nil, 2, 64)
	if err == nil {
		t.Fatal("no committee and no bootstrap, but no error")
	}
	for _, want := range []string{"fetch the period 1 update", "500", "bootstrap at slot 64"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("err = %v, want it to say %q", err, want)
		}
	}
}

// ethClientStateNode serves one 08-wasm Ethereum client state.
func ethClientStateNode(t *testing.T, state *relayerclient.EthereumClientState) *rpcmock.CosmosNode {
	t.Helper()
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	wasmState, err := codectypes.NewAnyWithValue(&ibcwasmtypes.ClientState{Data: data})
	if err != nil {
		t.Fatal(err)
	}
	response, err := (&clienttypes.QueryClientStateResponse{ClientState: wasmState}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	node := rpcmock.NewCosmosNode(t)
	node.Handle("abci_query", func(json.RawMessage) (any, error) {
		return &coretypes.ResultABCIQuery{Response: abci.ResponseQuery{Value: response}}, nil
	})
	return node
}

// #575: when light_client/updates is ahead of finality_update, the crossing
// header is the last one installed, but the proof state and the catch-up slot
// were taken from the older finality update, so proofs were built against a
// state the client was never given.
func TestBuildEthClientUpdateHeaders_ProofStateIsTheLastInstalledHeader(t *testing.T) {
	beacon := &crossingBeacon{
		// Finalized at 40 (block 500), signed at 42.
		finality: finalityJSON(41, 40, 42, 500),
		updates: func(string, string) (string, int) {
			// The period 1 update is newer: finalized at 50 (block 600), signed at 52.
			return "[" + updateJSON(16, 18, 456, true) + "," + updateJSON(50, 52, 600, true) + "]", http.StatusOK
		},
	}
	url := beacon.start(t)
	node := ethClientStateNode(t, crossingTestState(16))

	result, err := (&Worker{}).BuildEthClientUpdateHeaders(context.Background(),
		CosmosEndpoint{Client: node.Client(t)}, EVMEndpoint{BeaconAPIURL: url}, "08-wasm-7")
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(result.Headers) != 1 {
		t.Fatalf("headers = %d, want the crossing header alone (the finality update is older)", len(result.Headers))
	}
	if got := result.EthClientState; got.LatestSlot != 50 || got.LatestExecutionBlockNumber != 600 {
		t.Fatalf("proof state = slot %d block %d, want slot 50 block 600: the header installed last", got.LatestSlot, got.LatestExecutionBlockNumber)
	}
	if result.SigSlot != 52 {
		t.Fatalf("catch-up slot = %d, want 52, the installed header's signature slot", result.SigSlot)
	}
	if result.ProofTimestamp != 50 {
		t.Fatalf("proof timestamp = %d, want 50", result.ProofTimestamp)
	}
}

// The catch-up wait must cover every header in the update, not only the last:
// here the crossing header is signed after the finality update.
func TestBuildEthClientUpdateHeaders_CatchUpSlotIsTheHighestSignatureSlot(t *testing.T) {
	beacon := &crossingBeacon{
		finality: finalityJSON(46, 45, 47, 500),
		updates: func(string, string) (string, int) {
			return "[" + updateJSON(16, 18, 456, true) + "," + updateJSON(40, 62, 480, true) + "]", http.StatusOK
		},
	}
	url := beacon.start(t)
	node := ethClientStateNode(t, crossingTestState(16))

	result, err := (&Worker{}).BuildEthClientUpdateHeaders(context.Background(),
		CosmosEndpoint{Client: node.Client(t)}, EVMEndpoint{BeaconAPIURL: url}, "08-wasm-7")
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(result.Headers) != 2 || result.EthClientState.LatestSlot != 45 {
		t.Fatalf("result = %d headers at slot %d, want 2 at slot 45", len(result.Headers), result.EthClientState.LatestSlot)
	}
	if result.SigSlot != 62 {
		t.Fatalf("catch-up slot = %d, want 62: the crossing header is signed last", result.SigSlot)
	}
}

// The common case is unchanged: the finality update is the last header.
func TestBuildEthClientUpdateHeaders_ProofStateIsTheFinalityUpdateWhenItIsNewest(t *testing.T) {
	beacon := &crossingBeacon{
		finality: finalityJSON(41, 40, 42, 500),
		updates: func(string, string) (string, int) {
			return "[" + updateJSON(16, 18, 456, true) + "," + updateJSON(32, 34, 480, true) + "]", http.StatusOK
		},
	}
	url := beacon.start(t)
	node := ethClientStateNode(t, crossingTestState(16))

	result, err := (&Worker{}).BuildEthClientUpdateHeaders(context.Background(),
		CosmosEndpoint{Client: node.Client(t)}, EVMEndpoint{BeaconAPIURL: url}, "08-wasm-7")
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(result.Headers) != 2 || result.EthClientState.LatestSlot != 40 || result.EthClientState.LatestExecutionBlockNumber != 500 || result.SigSlot != 42 {
		t.Fatalf("result = %d headers, slot %d, block %d, sig %d; want 2, 40, 500, 42",
			len(result.Headers), result.EthClientState.LatestSlot, result.EthClientState.LatestExecutionBlockNumber, result.SigSlot)
	}
}
