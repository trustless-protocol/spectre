package rpcmock

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EVMNode is a JSON-RPC stub for an EVM execution client.
//
// Defaults are chosen so a caller that only wants "a working chain" writes no
// setup at all: chain id 1, nonce 0, 1 gwei gas price, a generous block gas
// limit, and eth_sendRawTransaction accepted. Override only what the test is
// actually about.
//
// Every field is guarded, because ethclient issues concurrent requests and the
// tests run under -race.
type EVMNode struct {
	server *httptest.Server

	mu             sync.Mutex
	chainID        uint64
	gasPrice       *big.Int
	blockGasLim    uint64
	nonces         map[common.Address]uint64
	receipts       map[common.Hash]*RawReceipt
	defaultRcpt    *RawReceipt
	sendErr        error
	callResult     string
	callBySelector map[string]string
	callErr        error
	sentRawTxs     []string
	nonceQueries   map[common.Address]int
	methodCalls    map[string]int
	code           map[common.Address]string
}

// RawReceipt is the subset of a transaction receipt the stub serves. It is a
// deliberate subset: anything the relayer does not read is not modelled, so a
// test cannot accidentally depend on a field the stub invents.
type RawReceipt struct {
	Status      uint64
	GasUsed     uint64
	BlockNumber uint64
}

// NewEVMNode starts a stub EVM node and stops it when the test finishes.
func NewEVMNode(t testing.TB) *EVMNode {
	t.Helper()
	n := &EVMNode{
		chainID:        1,
		gasPrice:       big.NewInt(1_000_000_000), // 1 gwei
		blockGasLim:    0xffffff,
		nonces:         make(map[common.Address]uint64),
		receipts:       make(map[common.Hash]*RawReceipt),
		callResult:     "0x",
		callBySelector: map[string]string{},
		nonceQueries:   make(map[common.Address]int),
		methodCalls:    make(map[string]int),
		code:           make(map[common.Address]string),
	}
	n.server = httptest.NewServer(n)
	t.Cleanup(n.server.Close)
	return n
}

// URL is the stub's address, for a caller that dials it itself.
func (n *EVMNode) URL() string { return n.server.URL }

// Client dials the stub and returns a real *ethclient.Client.
func (n *EVMNode) Client(t testing.TB) *ethclient.Client {
	t.Helper()
	client, err := ethclient.Dial(n.server.URL)
	if err != nil {
		t.Fatalf("dial stub EVM node: %v", err)
	}
	t.Cleanup(client.Close)
	return client
}

// --- knobs ---

// SetChainID changes the id reported by eth_chainId. Distinct ids are what keep
// nonce state separate per chain, so tests about that must set it.
func (n *EVMNode) SetChainID(id uint64) { n.withLock(func() { n.chainID = id }) }

// SetNonce sets the account nonce eth_getTransactionCount reports for addr.
func (n *EVMNode) SetNonce(addr common.Address, nonce uint64) {
	n.withLock(func() { n.nonces[addr] = nonce })
}

// SetReceipt makes eth_getTransactionReceipt answer for hash. Without one the
// stub reports "not mined yet" (a null result), which is what a caller polling
// for inclusion sees.
func (n *EVMNode) SetReceipt(hash common.Hash, r RawReceipt) {
	n.withLock(func() { n.receipts[hash] = &r })
}

// MineEverything makes eth_getTransactionReceipt answer for ANY hash, modelling a
// chain that includes whatever it is given. Use it when the test cannot know the
// hash in advance: a signed transaction's hash is derived from its contents, so a
// caller driving a real signer has no way to pre-register a receipt for it.
// SetReceipt still wins for hashes registered explicitly.
func (n *EVMNode) MineEverything(r RawReceipt) { n.withLock(func() { n.defaultRcpt = &r }) }

// FailSend makes eth_sendRawTransaction return an RPC error. Pass nil to accept
// again.
func (n *EVMNode) FailSend(err error) { n.withLock(func() { n.sendErr = err }) }

// SetCode makes eth_getCode report codeHex ("0x...") for addr, modelling a
// deployed contract. An address left unset answers "0x", which is what a real
// node returns for one nobody deployed to — so a test must opt IN to a contract
// existing, and can never accidentally assume one does.
func (n *EVMNode) SetCode(addr common.Address, codeHex string) {
	n.withLock(func() { n.code[addr] = codeHex })
}

// SetCallResult sets the ABI-encoded return data eth_call answers with.
func (n *EVMNode) SetCallResult(hexData string) { n.withLock(func() { n.callResult = hexData }) }

// SetCallResultFor answers one 4-byte selector differently from the rest.
//
// SetCallResult alone cannot express a contract whose methods disagree, and that
// is exactly what a wiring check reads: getClient returns an address while
// getCounterparty returns a struct, and a test that cannot tell them apart can
// only assert that some call happened. selector is the hex method id, with or
// without the 0x prefix.
func (n *EVMNode) SetCallResultFor(selector, hexData string) {
	n.withLock(func() { n.callBySelector[normalizeSelector(selector)] = hexData })
}

func normalizeSelector(selector string) string {
	return strings.ToLower(strings.TrimPrefix(selector, "0x"))
}

// FailCall makes eth_call return an RPC error, the shape a reverting contract
// read takes.
func (n *EVMNode) FailCall(err error) { n.withLock(func() { n.callErr = err }) }

// SetBlockGasLimit changes the limit reported by eth_getBlockByNumber, which is
// what batch-splitting decisions are measured against.
func (n *EVMNode) SetBlockGasLimit(limit uint64) { n.withLock(func() { n.blockGasLim = limit }) }

// --- observations ---

// SentRawTxs returns the hex payloads passed to eth_sendRawTransaction, in order.
func (n *EVMNode) SentRawTxs() []string {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]string(nil), n.sentRawTxs...)
}

// NonceQueries counts eth_getTransactionCount calls for addr — the measurement
// behind "the nonce is cached, not re-read per transaction".
func (n *EVMNode) NonceQueries(addr common.Address) int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.nonceQueries[addr]
}

// MethodCalls counts calls to a JSON-RPC method by name.
func (n *EVMNode) MethodCalls(method string) int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.methodCalls[method]
}

func (n *EVMNode) withLock(f func()) {
	n.mu.Lock()
	defer n.mu.Unlock()
	f()
}

// --- transport ---

type rpcRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	ID      json.RawMessage   `json:"id"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (n *EVMNode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req rpcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	n.mu.Lock()
	n.methodCalls[req.Method]++
	result, rpcErr := n.dispatch(req)
	n.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rpcResponse{
		JSONRPC: "2.0", ID: req.ID, Result: result, Error: rpcErr,
	})
}

// dispatch answers one call. The caller holds n.mu.
func (n *EVMNode) dispatch(req rpcRequest) (any, *rpcError) {
	switch req.Method {
	case "eth_chainId":
		return hexUint(n.chainID), nil
	case "eth_gasPrice", "eth_maxPriorityFeePerGas":
		return "0x" + n.gasPrice.Text(16), nil
	case "eth_getTransactionCount":
		addr := addressParam(req.Params)
		n.nonceQueries[addr]++
		return hexUint(n.nonces[addr]), nil
	case "eth_sendRawTransaction":
		if n.sendErr != nil {
			return nil, &rpcError{Code: -32000, Message: n.sendErr.Error()}
		}
		n.sentRawTxs = append(n.sentRawTxs, stringParam(req.Params, 0))
		return zeroHash, nil
	case "eth_getTransactionReceipt":
		hash := common.HexToHash(stringParam(req.Params, 0))
		receipt, ok := n.receipts[hash]
		if !ok || receipt == nil {
			receipt = n.defaultRcpt
		}
		if receipt == nil {
			return nil, nil // not mined yet
		}
		return receiptJSON(hash, receipt), nil
	case "eth_getBlockByNumber":
		return n.blockJSON(), nil
	case "eth_getCode":
		// "0x" is the honest default: an address nobody deployed to holds no
		// code, and that is what a real node answers. A test about a deployed
		// contract says so with SetCode.
		addr := common.HexToAddress(stringParam(req.Params, 0))
		if code, ok := n.code[addr]; ok {
			return code, nil
		}
		return "0x", nil
	case "eth_call":
		if n.callErr != nil {
			return nil, &rpcError{Code: -32000, Message: n.callErr.Error()}
		}
		if result, ok := n.callResultForRequest(req.Params); ok {
			return result, nil
		}
		return n.callResult, nil
	case "eth_estimateGas":
		return hexUint(21_000), nil
	default:
		return nil, &rpcError{Code: -32601, Message: "method not found: " + req.Method}
	}
}

const zeroHash = "0x0000000000000000000000000000000000000000000000000000000000000000"

func hexUint(v uint64) string { return fmt.Sprintf("0x%x", v) }

func stringParam(params []json.RawMessage, i int) string {
	if len(params) <= i {
		return ""
	}
	var s string
	if err := json.Unmarshal(params[i], &s); err != nil {
		return ""
	}
	return s
}

func addressParam(params []json.RawMessage) common.Address {
	return common.HexToAddress(stringParam(params, 0))
}

func receiptJSON(hash common.Hash, r *RawReceipt) map[string]any {
	return map[string]any{
		"transactionHash":   hash.Hex(),
		"transactionIndex":  "0x0",
		"blockHash":         zeroHash,
		"blockNumber":       hexUint(r.BlockNumber),
		"gasUsed":           hexUint(r.GasUsed),
		"cumulativeGasUsed": hexUint(r.GasUsed),
		"status":            hexUint(r.Status),
		"logsBloom":         "0x" + strings.Repeat("0", 512),
		"logs":              []any{},
	}
}

func (n *EVMNode) blockJSON() map[string]any {
	return map[string]any{
		"parentHash":       zeroHash,
		"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"miner":            "0x0000000000000000000000000000000000000000",
		"stateRoot":        zeroHash,
		"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"receiptsRoot":     "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"logsBloom":        "0x" + strings.Repeat("0", 512),
		"difficulty":       "0x0",
		"number":           "0x1",
		"gasLimit":         hexUint(n.blockGasLim),
		"gasUsed":          "0x0",
		"timestamp":        "0x0",
		"extraData":        "0x",
		"mixHash":          zeroHash,
		"nonce":            "0x0000000000000000",
		"baseFeePerGas":    "0x" + n.gasPrice.Text(16),
	}
}

// callResultForRequest picks a per-selector answer when one was registered.
func (n *EVMNode) callResultForRequest(params []json.RawMessage) (string, bool) {
	if len(n.callBySelector) == 0 {
		return "", false
	}
	if len(params) == 0 {
		return "", false
	}
	var arg struct {
		Data  string `json:"data"`
		Input string `json:"input"`
	}
	if err := json.Unmarshal(params[0], &arg); err != nil {
		return "", false
	}
	data := arg.Data
	if data == "" {
		data = arg.Input
	}
	data = normalizeSelector(data)
	if len(data) < 8 {
		return "", false
	}
	result, ok := n.callBySelector[data[:8]]
	return result, ok
}
