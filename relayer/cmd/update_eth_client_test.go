package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtbytes "github.com/cometbft/cometbft/libs/bytes"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	relayerclient "relayer/client"
	"relayer/rpcmock"
	"relayer/services"
)

func TestUpdateEthClient(t *testing.T) {
	// Keep local .env files and operator overrides out of command execution.
	t.Chdir(t.TempDir())
	t.Setenv("COSMOS_WASM_CLIENT_ID", "")
	t.Setenv("COSMOS_PRIVATE_KEY", "")
	submitErr := errors.New("broadcast rejected")
	for _, tc := range []struct {
		name           string
		alreadyCurrent bool
		submitErr      error
	}{
		{name: "already current skips signing and submission", alreadyCurrent: true},
		{name: "submits ordered headers and reports resulting state"},
		{name: "returns submission error without reporting success", submitErr: submitErr},
	} {
		t.Run(tc.name, func(t *testing.T) {
			configPath, node := ethClientUpdateTestConfig(t, tc.alreadyCurrent)
			tx := &ethClientUpdateTestHandler{submitErr: tc.submitErr}
			cmd := newUpdateEthClientCommand(zap.NewNop(), tx)
			var output bytes.Buffer
			cmd.SetOut(&output)
			cmd.SetErr(&output)
			cmd.SetArgs([]string{"--config", configPath})
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := cmd.ExecuteContext(ctx)

			if tc.submitErr != nil {
				if !errors.Is(err, tc.submitErr) {
					t.Fatalf("ExecuteContext error = %v, want wrapped %v", err, tc.submitErr)
				}
				if !strings.Contains(err.Error(), "failed to update Ethereum client on Cosmos") ||
					!strings.Contains(err.Error(), "exec block 500") {
					t.Fatalf("submission error lacks command/height context: %v", err)
				}
				if strings.Contains(output.String(), "update-eth-client complete") || strings.Contains(output.String(), "already current") {
					t.Fatalf("failed submission reported success: %s", output.String())
				}
			} else if err != nil {
				t.Fatalf("ExecuteContext: %v", err)
			}

			if tc.alreadyCurrent {
				if tx.signerCalls != 0 || tx.submitCalls != 0 {
					t.Fatalf("already-current client attempted signing/submission: signer=%d submit=%d", tx.signerCalls, tx.submitCalls)
				}
				if node.Calls("abci_query") != 1 || node.Calls("status") != 0 {
					t.Fatal("already-current client entered the destination submission path")
				}
				const want = "update-eth-client: Ethereum light client on Cosmos is already current; no transaction submitted (latest_slot=40 latest_execution_block_number=500)\n"
				if output.String() != want {
					t.Fatalf("output = %q, want %q", output.String(), want)
				}
				return
			}

			if tx.signerCalls != 1 || tx.submitCalls != 1 {
				t.Fatalf("signer calls=%d submit calls=%d, want one each", tx.signerCalls, tx.submitCalls)
			}
			if tx.ctx != ctx {
				t.Fatal("submission did not receive the command context")
			}
			if tx.endpoint.CosmosClient().IsRunning() {
				t.Fatal("recovery command started the Cosmos WebSocket client")
			}
			if node.Calls("abci_query") != 2 || node.Calls("status") != 1 {
				t.Fatal("submission skipped the destination's client-state/timing checks")
			}
			if len(tx.msgs) != 2 {
				t.Fatalf("submitted %d messages, want crossing header then finality header", len(tx.msgs))
			}
			for i, want := range []struct {
				trustedSlot   uint64
				finalizedSlot string
			}{{16, "32"}, {32, "40"}} {
				msg, ok := tx.msgs[i].(*clienttypes.MsgUpdateClient)
				if !ok {
					t.Fatalf("message %d = %T, want *MsgUpdateClient", i, tx.msgs[i])
				}
				if msg.ClientId != "08-wasm-7" || msg.Signer != "cosmos1testsigner" {
					t.Fatalf("message %d has client=%q signer=%q", i, msg.ClientId, msg.Signer)
				}
				var wrapped ibcwasmtypes.ClientMessage
				if err := wrapped.Unmarshal(msg.ClientMessage.Value); err != nil {
					t.Fatalf("unwrap message %d: %v", i, err)
				}
				var header relayerclient.EthereumHeader
				if err := json.Unmarshal(wrapped.Data, &header); err != nil {
					t.Fatalf("decode header %d: %v", i, err)
				}
				if header.TrustedSlot != want.trustedSlot || header.ConsensusUpdate.FinalizedHeader.Beacon.Slot != want.finalizedSlot {
					t.Fatalf("header %d: trusted=%d finalized=%s, want trusted=%d finalized=%s", i,
						header.TrustedSlot, header.ConsensusUpdate.FinalizedHeader.Beacon.Slot, want.trustedSlot, want.finalizedSlot)
				}
			}
			if tc.submitErr == nil {
				const want = "update-eth-client complete: latest_slot=40 latest_execution_block_number=500\n"
				if output.String() != want {
					t.Fatalf("output = %q, want %q", output.String(), want)
				}
			}
		})
	}
}

// Only the signing/broadcast boundary is replaced. The command builds headers
// from local RPC responses and submits through the real cosmos.Destination.
type ethClientUpdateTestHandler struct {
	services.TransactionHandler
	signerCalls, submitCalls int
	submitErr                error
	ctx                      context.Context
	endpoint                 services.CosmosEndpoint
	msgs                     []any
}

func (h *ethClientUpdateTestHandler) CosmosSignerAddress() (string, error) {
	h.signerCalls++
	return "cosmos1testsigner", nil
}

func (h *ethClientUpdateTestHandler) SendCosmosTxBatch(ctx context.Context, endpoint services.CosmosEndpoint, msgs []any) error {
	h.submitCalls++
	h.ctx, h.endpoint, h.msgs = ctx, endpoint, msgs
	return h.submitErr
}

func ethClientUpdateTestConfig(t *testing.T, alreadyCurrent bool) (string, *rpcmock.CosmosNode) {
	t.Helper()
	state := relayerclient.EthereumClientState{
		LatestSlot: 16, LatestExecutionBlockNumber: 456,
		SlotsPerEpoch: 8, EpochsPerSyncCommitteePeriod: 4,
		SecondsPerSlot: 1, SyncCommitteeSize: 8,
	}
	if alreadyCurrent {
		state.LatestSlot, state.LatestExecutionBlockNumber = 40, 500
	}
	stateJSON, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	wasmState, err := codectypes.NewAnyWithValue(&ibcwasmtypes.ClientState{Data: stateJSON})
	if err != nil {
		t.Fatal(err)
	}
	response := clienttypes.QueryClientStateResponse{ClientState: wasmState}
	responseBytes, err := response.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	node := rpcmock.NewCosmosNode(t)
	node.Handle("abci_query", func(params json.RawMessage) (any, error) {
		var query struct {
			Path string            `json:"path"`
			Data cmtbytes.HexBytes `json:"data"`
		}
		if err := cmtjson.Unmarshal(params, &query); err != nil {
			return nil, err
		}
		var req clienttypes.QueryClientStateRequest
		if err := req.Unmarshal(query.Data); err != nil {
			return nil, err
		}
		if query.Path != "/ibc.core.client.v1.Query/ClientState" || req.ClientId != "08-wasm-7" {
			return nil, fmt.Errorf("unexpected client query: path=%q client=%q", query.Path, req.ClientId)
		}
		return &coretypes.ResultABCIQuery{Response: abci.ResponseQuery{Value: responseBytes}}, nil
	})
	node.Handle("status", func(json.RawMessage) (any, error) {
		return &coretypes.ResultStatus{SyncInfo: coretypes.SyncInfo{LatestBlockTime: time.Unix(1000, 0)}}, nil
	})
	beacon := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/eth/v1/beacon/light_client/finality_update":
			_, _ = fmt.Fprint(w, `{"data":{
				"attested_header":{"beacon":{"slot":"41"}},
				"finalized_header":{"beacon":{"slot":"40"},"execution":{"block_number":"500","timestamp":"40"}},
				"signature_slot":"42","sync_aggregate":{"sync_committee_bits":"0xff"}
			}}`)
		case "/eth/v1/beacon/light_client/updates":
			if alreadyCurrent || r.URL.Query().Get("start_period") != "0" || r.URL.Query().Get("count") != "2" {
				t.Errorf("unexpected period-update request: %s", r.URL)
				http.Error(w, "unexpected period range", http.StatusBadRequest)
				return
			}
			_, _ = fmt.Fprint(w, `[
				{"data":{"finalized_header":{"beacon":{"slot":"16"}},"next_sync_committee":{"pubkeys":[],"aggregate_pubkey":"0x00"}}},
				{"data":{"finalized_header":{"beacon":{"slot":"32"}},"signature_slot":"34"}}
			]`)
		default:
			t.Errorf("unexpected beacon request: %s", r.URL)
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(beacon.Close)
	// Execution RPC and WebSocket services are deliberately absent.
	configPath := writeTempConfig(t, fmt.Sprintf(`{"modules":[
		{"name":"cosmos_to_eth","config":{
			"tm_rpc_url":%q,"eth_rpc_url":"http://127.0.0.1:1",
			"ics26_address":"0x1111111111111111111111111111111111111111",
			"cosmos_wasm_client_id":"08-wasm-7"}},
		{"name":"eth_to_cosmos","config":{"eth_beacon_api_url":%q}}
	]}`, node.URL(), beacon.URL))
	return configPath, node
}

func TestUpdateClientDirectionHelp(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		cmd  func() *cobra.Command
		want string
	}{
		{
			name: "cosmos client on Ethereum",
			cmd:  func() *cobra.Command { return UpdateClient(zap.NewNop()) },
			want: "Cosmos → Ethereum",
		},
		{
			name: "Ethereum client on Cosmos",
			cmd:  func() *cobra.Command { return UpdateEthClient(zap.NewNop()) },
			want: "Ethereum → Cosmos",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cmd := tc.cmd()
			var output bytes.Buffer
			cmd.SetOut(&output)
			cmd.SetErr(&output)
			cmd.SetArgs([]string{"--help"})
			if err := cmd.Execute(); err != nil {
				t.Fatalf("help: %v", err)
			}
			if got := output.String(); !strings.Contains(got, tc.want) {
				t.Fatalf("help does not distinguish direction %q:\n%s", tc.want, got)
			}
		})
	}
}

func TestEthClientUpdateSummary(t *testing.T) {
	t.Parallel()
	state := &relayerclient.EthereumClientState{
		LatestSlot:                 123,
		LatestExecutionBlockNumber: 456,
	}

	if got := ethClientUpdateSummary(state, true); got != "update-eth-client complete: latest_slot=123 latest_execution_block_number=456" {
		t.Errorf("submitted summary = %q", got)
	}
	if got := ethClientUpdateSummary(state, false); got != "update-eth-client: Ethereum light client on Cosmos is already current; no transaction submitted (latest_slot=123 latest_execution_block_number=456)" {
		t.Errorf("no-op summary = %q", got)
	}
}
