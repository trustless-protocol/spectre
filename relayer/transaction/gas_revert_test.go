package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"relayer/services"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

func revertProbeClient(t *testing.T, blockGas uint64, callError string, onCall func(uint64)) *ethclient.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     any              `json:"id"`
			Method string           `json:"method"`
			Params []map[string]any `json:"params"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp := map[string]any{"jsonrpc": "2.0", "id": req.ID}
		switch req.Method {
		case "eth_getBlockByNumber":
			resp["result"] = map[string]any{
				"parentHash":       "0x0000000000000000000000000000000000000000000000000000000000000000",
				"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"miner":            "0x0000000000000000000000000000000000000000",
				"stateRoot":        "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"receiptsRoot":     "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"logsBloom":        "0x" + strings.Repeat("0", 512),
				"difficulty":       "0x0", "number": "0x1", "gasLimit": fmt.Sprintf("0x%x", blockGas),
				"gasUsed": "0x0", "timestamp": "0x0", "extraData": "0x",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce":   "0x0000000000000000",
			}
		case "eth_call":
			if onCall != nil && len(req.Params) > 0 {
				if raw, ok := req.Params[0]["gas"].(string); ok {
					gas, err := hexutil.DecodeUint64(raw)
					if err == nil {
						onCall(gas)
					}
				}
			}
			if callError == "" {
				resp["result"] = "0x"
			} else {
				resp["error"] = map[string]any{"code": -32000, "message": callError}
			}
		default:
			resp["error"] = map[string]any{"code": -32601, "message": "method not found"}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)
	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("dial mock RPC: %v", err)
	}
	t.Cleanup(client.Close)
	return client
}

func TestClassifyEthRevertDoesNotTrustRejectedProbe(t *testing.T) {
	const gasLimit = uint64(16_000_000)
	client := revertProbeClient(t, 30_000_000, "gas required exceeds allowance", nil)
	endpoint := services.EVMEndpoint{Client: client}
	got := classifyEthRevert(context.Background(), endpoint, ethereum.CallMsg{}, big.NewInt(1), nil, gasLimit/4, gasLimit)
	if got != revertLogical {
		t.Fatalf("rejected diagnostic probe classified as %v, want logical", got)
	}
}

func TestClassifyEthRevertClampsProbeToBlockGas(t *testing.T) {
	const (
		gasLimit = uint64(16_000_000)
		blockGas = uint64(30_000_000)
	)
	var observed uint64
	client := revertProbeClient(t, blockGas, "", func(gas uint64) { observed = gas })
	endpoint := services.EVMEndpoint{Client: client}
	if got := classifyEthRevert(context.Background(), endpoint, ethereum.CallMsg{}, big.NewInt(1), nil, gasLimit/4, gasLimit); got != revertOutOfGas {
		t.Fatalf("successful higher-gas replay = %v, want out of gas", got)
	}
	if observed != blockGas {
		t.Fatalf("probe gas = %d, want clamped block gas %d", observed, blockGas)
	}
}
