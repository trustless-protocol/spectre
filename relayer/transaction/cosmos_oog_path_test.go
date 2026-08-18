package transaction

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"relayer/services"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtbytes "github.com/cometbft/cometbft/libs/bytes"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	txservice "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/gogoproto/proto"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
)

type cosmosOOGRPC struct {
	t *testing.T

	checkTxOOG   int
	deliverTxOOG int
	broadcasts   int
	txQueries    int
	gasLimits    []uint64
	feeAmounts   []int64
}

type cosmosRPCRequest struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type cosmosRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   any             `json:"error,omitempty"`
}

func (m *cosmosOOGRPC) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.t.Helper()
	var req cosmosRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		m.t.Errorf("decode RPC request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "consensus_params":
		m.writeError(w, req.ID, -32603, "block max unavailable")
	case "abci_query":
		m.writeResult(w, req.ID, &coretypes.ResultABCIQuery{
			Response: abci.ResponseQuery{Code: 1, Codespace: "sdk", Log: "simulation unavailable"},
		})
	case "broadcast_tx_sync":
		m.broadcasts++
		m.captureGas(req.Params)
		if m.broadcasts <= m.checkTxOOG {
			m.writeResult(w, req.ID, &coretypes.ResultBroadcastTx{
				Code:      errortypes.ErrOutOfGas.ABCICode(),
				Codespace: errortypes.ErrOutOfGas.Codespace(),
				Log:       "out of gas in CheckTx",
			})
			return
		}
		m.writeResult(w, req.ID, &coretypes.ResultBroadcastTx{
			Hash: cmtbytes.HexBytes(testCosmosTxHash(m.broadcasts)),
		})
	case "tx":
		m.txQueries++
		result := abci.ExecTxResult{}
		if m.txQueries <= m.deliverTxOOG {
			result.Code = errortypes.ErrOutOfGas.ABCICode()
			result.Codespace = errortypes.ErrOutOfGas.Codespace()
			result.Log = "out of gas in DeliverTx"
		}
		m.writeResult(w, req.ID, &coretypes.ResultTx{
			Hash:     cmtbytes.HexBytes(testCosmosTxHash(m.broadcasts)),
			Height:   int64(m.txQueries),
			TxResult: result,
		})
	default:
		m.writeError(w, req.ID, -32601, "method not found")
	}
}

func testCosmosTxHash(n int) []byte {
	hash := make([]byte, 32)
	hash[len(hash)-1] = byte(n)
	return hash
}

func (m *cosmosOOGRPC) captureGas(params json.RawMessage) {
	var named map[string]json.RawMessage
	if err := json.Unmarshal(params, &named); err != nil {
		m.t.Fatalf("decode broadcast params: %v", err)
	}
	var encoded string
	if err := json.Unmarshal(named["tx"], &encoded); err != nil {
		m.t.Fatalf("decode broadcast tx parameter: %v", err)
	}
	txBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		m.t.Fatalf("decode broadcast tx bytes: %v", err)
	}
	var raw txservice.TxRaw
	if err := proto.Unmarshal(txBytes, &raw); err != nil {
		m.t.Fatalf("decode TxRaw: %v", err)
	}
	var authInfo txservice.AuthInfo
	if err := proto.Unmarshal(raw.AuthInfoBytes, &authInfo); err != nil {
		m.t.Fatalf("decode AuthInfo: %v", err)
	}
	if authInfo.Fee == nil {
		m.t.Fatal("broadcast transaction has no fee")
	}
	m.gasLimits = append(m.gasLimits, authInfo.Fee.GasLimit)
	if len(authInfo.Fee.Amount) != 1 {
		m.t.Fatalf("broadcast fee coins = %v, want exactly one", authInfo.Fee.Amount)
	}
	m.feeAmounts = append(m.feeAmounts, authInfo.Fee.Amount[0].Amount.Int64())
}

func (m *cosmosOOGRPC) writeResult(w http.ResponseWriter, id json.RawMessage, result any) {
	m.t.Helper()
	raw, err := cmtjson.Marshal(result)
	if err != nil {
		m.t.Fatalf("marshal RPC result: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cosmosRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  raw,
	}); err != nil {
		m.t.Errorf("encode RPC response: %v", err)
	}
}

func (m *cosmosOOGRPC) writeError(w http.ResponseWriter, id json.RawMessage, code int, message string) {
	m.t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cosmosRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: map[string]any{
			"code":    code,
			"message": message,
		},
	}); err != nil {
		m.t.Errorf("encode RPC error: %v", err)
	}
}

func newCosmosOOGEndpoint(t *testing.T, mock *cosmosOOGRPC) services.CosmosEndpoint {
	t.Helper()
	server := httptest.NewServer(mock)
	t.Cleanup(server.Close)
	client, err := rpchttp.New(server.URL, "/websocket")
	if err != nil {
		t.Fatalf("create CometBFT client: %v", err)
	}
	return services.CosmosEndpoint{Client: client}
}

func cosmosOOGTestMessage() sdk.Msg {
	return &clienttypes.MsgUpdateClient{
		ClientId: "08-wasm-0",
		ClientMessage: &codectypes.Any{
			TypeUrl: "/test.ClientMessage",
			Value:   []byte{1},
		},
	}
}

func configureCosmosOOGTest(t *testing.T) {
	t.Helper()
	t.Setenv("COSMOS_PRIVATE_KEY", testCosmosPrivKeyHex)
	t.Setenv("COSMOS_CHAIN_ID", "test-chain")
	t.Setenv("COSMOS_GAS_LIMIT", "200000")
	t.Setenv("COSMOS_FEE_DENOM", "stake")
	t.Setenv("COSMOS_FEE_AMOUNT", "")
}

func assertGasStrictlyIncreases(t *testing.T, limits []uint64) {
	t.Helper()
	if len(limits) < 2 {
		t.Fatalf("broadcast gas limits = %v, want at least two attempts", limits)
	}
	for i := 1; i < len(limits); i++ {
		if limits[i] <= limits[i-1] {
			t.Fatalf("broadcast gas limits = %v; attempt %d did not increase", limits, i+1)
		}
	}
}

func TestFallbackGasIncreasesAfterDeliverTxOutOfGas(t *testing.T) {
	configureCosmosOOGTest(t)
	mock := &cosmosOOGRPC{t: t, deliverTxOOG: 1}
	endpoint := newCosmosOOGEndpoint(t, mock)

	h := &Handler{}
	nextSequence, succeeded, err := h.sendCosmosTxBatchWithSplitting(
		context.Background(), endpoint, []sdk.Msg{cosmosOOGTestMessage()}, 7, 10, true,
	)
	if err != nil {
		t.Fatalf("send after DeliverTx OOG: %v", err)
	}
	if nextSequence != 12 || succeeded != 1 {
		t.Fatalf("result = sequence:%d succeeded:%d, want 12:1", nextSequence, succeeded)
	}
	assertGasStrictlyIncreases(t, mock.gasLimits)
}

func TestCheckTxOutOfGasRetriesWithMoreGasAndRetainsSequence(t *testing.T) {
	configureCosmosOOGTest(t)
	mock := &cosmosOOGRPC{t: t, checkTxOOG: 1}
	endpoint := newCosmosOOGEndpoint(t, mock)

	h := &Handler{}
	nextSequence, succeeded, err := h.sendCosmosTxBatchWithSplitting(
		context.Background(), endpoint, []sdk.Msg{cosmosOOGTestMessage()}, 7, 10, true,
	)
	if err != nil {
		t.Fatalf("send after CheckTx OOG: %v", err)
	}
	if nextSequence != 11 || succeeded != 1 {
		t.Fatalf("result = sequence:%d succeeded:%d, want 11:1; CheckTx must not consume sequence 10", nextSequence, succeeded)
	}
	assertGasStrictlyIncreases(t, mock.gasLimits)
}

func TestCheckTxOutOfGasScalesConfiguredFeeWithGas(t *testing.T) {
	configureCosmosOOGTest(t)
	t.Setenv("COSMOS_FEE_AMOUNT", "1000")
	mock := &cosmosOOGRPC{t: t, checkTxOOG: 1}
	endpoint := newCosmosOOGEndpoint(t, mock)

	h := &Handler{}
	_, _, err := h.sendCosmosTxBatchWithSplitting(
		context.Background(), endpoint, []sdk.Msg{cosmosOOGTestMessage()}, 7, 10, true,
	)
	if err != nil {
		t.Fatalf("send after CheckTx OOG: %v", err)
	}
	if len(mock.feeAmounts) != 2 {
		t.Fatalf("broadcast fee amounts = %v, want two attempts", mock.feeAmounts)
	}
	if mock.feeAmounts[0] != 1000 {
		t.Fatalf("initial configured fee = %d, want historical value 1000", mock.feeAmounts[0])
	}
	wantRetryFee := applyCosmosFeeHeadroom(1000, 1)
	if mock.feeAmounts[1] != wantRetryFee || mock.feeAmounts[1] <= mock.feeAmounts[0] {
		t.Fatalf("broadcast fee amounts = %v, want retry fee to grow to %d with gas headroom", mock.feeAmounts, wantRetryFee)
	}
	if mock.feeAmounts[1]*int64(mock.gasLimits[0]) < mock.feeAmounts[0]*int64(mock.gasLimits[1]) {
		t.Fatalf("fee per gas decreased across retry: fees=%v gas=%v", mock.feeAmounts, mock.gasLimits)
	}
	assertGasStrictlyIncreases(t, mock.gasLimits)
}

func TestCheckTxOutOfGasExhaustsTheHeadroomLadder(t *testing.T) {
	configureCosmosOOGTest(t)
	mock := &cosmosOOGRPC{t: t, checkTxOOG: 100}
	endpoint := newCosmosOOGEndpoint(t, mock)

	h := &Handler{}
	nextSequence, succeeded, err := h.sendCosmosTxBatchWithSplitting(
		context.Background(), endpoint, []sdk.Msg{cosmosOOGTestMessage()}, 7, 10, true,
	)
	if !errors.Is(err, services.ErrPermanentRelayFailure) {
		t.Fatalf("exhausted CheckTx ladder error = %v, want permanent relay failure", err)
	}
	if nextSequence != 10 || succeeded != 0 {
		t.Fatalf("result = sequence:%d succeeded:%d, want 10:0", nextSequence, succeeded)
	}
	if len(mock.gasLimits) != len(cosmosGasHeadroom) {
		t.Fatalf("broadcast attempts = %d, want finite ladder length %d", len(mock.gasLimits), len(cosmosGasHeadroom))
	}
	assertGasStrictlyIncreases(t, mock.gasLimits)
}
