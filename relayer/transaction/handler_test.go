package transaction

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/services"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
)

type testDataError struct {
	data interface{}
}

func (e testDataError) Error() string {
	return "execution reverted"
}

func (e testDataError) ErrorData() interface{} {
	return e.data
}

// These tests cover the calldata packing layer that SendEthTxBatch relies on.
// We do not exercise the RPC submission path here — that requires a full
// endpoint dependency + ethclient mock and is validated end-to-end via the
// benchmark run described in /Users/ducnt/.claude/plans/...
//
// Selector reference (from `cast sig` against the ICS26Router ABI):
//   recvPacket(...)              → 0x596e00b9
//   ackPacket(...)               → 0xfdbd955d
//   timeoutPacket(...)           → 0x223e357a
//   updateApplicationState(string,bytes) → 0x9c11bece
//   updateConsensusState(string,bytes)   → 0xd395af7f
//   multicall(bytes[])                   → 0xac9650d8

func mustDecodeHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("decode %q: %v", s, err)
	}
	return b
}

func samplePacket() contractICS26Router.IICS26RouterMsgsPacket {
	return contractICS26Router.IICS26RouterMsgsPacket{
		Sequence:         42,
		SourceClient:     "cosmoshub-1",
		DestClient:       "08-wasm-0",
		TimeoutTimestamp: 1_700_000_000,
		Payloads: []contractICS26Router.IICS26RouterMsgsPayload{{
			SourcePort: "transfer",
			DestPort:   "transfer",
			Version:    "ics20-2",
			Encoding:   "abi",
			Value:      []byte{0x01, 0x02, 0x03},
		}},
	}
}

func TestSelectorsForBatchedMsgs(t *testing.T) {
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	pkt := samplePacket()

	cases := []struct {
		name     string
		msg      any
		method   string
		selector string
	}{
		{
			name: "recvPacket",
			msg: contractICS26Router.IICS26RouterMsgsMsgRecvPacket{
				Packet:        pkt,
				MembershipMsg: []byte{0xaa, 0xbb},
			},
			method:   "recvPacket",
			selector: "596e00b9",
		},
		{
			name: "ackPacket",
			msg: contractICS26Router.IICS26RouterMsgsMsgAckPacket{
				Packet:          pkt,
				Acknowledgement: []byte{0xde, 0xad},
				MembershipMsg:   []byte{0xbe, 0xef},
			},
			method:   "ackPacket",
			selector: "fdbd955d",
		},
		{
			name: "timeoutPacket",
			msg: contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket{
				Packet:           pkt,
				NonMembershipMsg: []byte{0xca, 0xfe},
			},
			method:   "timeoutPacket",
			selector: "223e357a",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := parsedABI.Pack(tc.method, tc.msg)
			if err != nil {
				t.Fatalf("pack: %v", err)
			}
			if len(data) < 4 {
				t.Fatalf("calldata too short: %d", len(data))
			}
			gotSelector := hex.EncodeToString(data[:4])
			if gotSelector != tc.selector {
				t.Fatalf("selector mismatch: got %s, want %s", gotSelector, tc.selector)
			}
		})
	}
}

func TestValidatorCacheRaceErrorName(t *testing.T) {
	selector := errorSelector("ValidatorSetCacheMiss(bytes32)")
	data := append(selector[:], make([]byte, 32)...)
	name, ok := validatorCacheRaceErrorName(testDataError{data: "0x" + hex.EncodeToString(data)})
	if !ok {
		t.Fatal("expected ValidatorSetCacheMiss to be classified as a cache race")
	}
	if name != "ValidatorSetCacheMiss" {
		t.Fatalf("error name: got %q want ValidatorSetCacheMiss", name)
	}

	unknown := errorSelector("ProofVerificationFailed()")
	if name, ok := validatorCacheRaceErrorName(testDataError{data: unknown[:]}); ok {
		t.Fatalf("did not expect unknown selector to be classified, got %q", name)
	}
}

func TestCloneCosmosSDKMsgWithSignerDoesNotMutateOriginal(t *testing.T) {
	original := &clienttypes.MsgUpdateClient{ClientId: "08-wasm-0"}

	cloned, err := cloneCosmosSDKMsgWithSigner(original, "cosmos1signer", -1)
	if err != nil {
		t.Fatalf("clone: %v", err)
	}
	got := cloned.(*clienttypes.MsgUpdateClient)
	if got.Signer != "cosmos1signer" {
		t.Fatalf("cloned signer = %q, want cosmos1signer", got.Signer)
	}
	if original.Signer != "" {
		t.Fatalf("original signer mutated to %q", original.Signer)
	}
}

func TestCloneCosmosSDKMsgWithSignerPreservesExistingSigner(t *testing.T) {
	original := &clienttypes.MsgUpdateClient{ClientId: "08-wasm-0", Signer: "cosmos1original"}

	cloned, err := cloneCosmosSDKMsgWithSigner(original, "cosmos1default", -1)
	if err != nil {
		t.Fatalf("clone: %v", err)
	}
	got := cloned.(*clienttypes.MsgUpdateClient)
	if got.Signer != "cosmos1original" {
		t.Fatalf("cloned signer = %q, want existing signer", got.Signer)
	}
	if original.Signer != "cosmos1original" {
		t.Fatalf("original signer mutated to %q", original.Signer)
	}
}

// TestMulticallWraps verifies that wrapping per-msg calldata into a multicall
// produces the expected outer selector + that each inner blob is recoverable
// after Unpack. This is the round-trip that on-chain MulticallUpgradeable
// would do.
func TestMulticallWraps(t *testing.T) {
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	pkt := samplePacket()

	recv, err := parsedABI.Pack("recvPacket", contractICS26Router.IICS26RouterMsgsMsgRecvPacket{
		Packet:        pkt,
		MembershipMsg: []byte{0x01},
	})
	if err != nil {
		t.Fatalf("pack recv: %v", err)
	}
	ack, err := parsedABI.Pack("ackPacket", contractICS26Router.IICS26RouterMsgsMsgAckPacket{
		Packet:          pkt,
		Acknowledgement: []byte{0x02},
		MembershipMsg:   []byte{0x03},
	})
	if err != nil {
		t.Fatalf("pack ack: %v", err)
	}

	innerCalldata := [][]byte{recv, ack}
	outer, err := parsedABI.Pack("multicall", innerCalldata)
	if err != nil {
		t.Fatalf("pack multicall: %v", err)
	}
	wantSelector := mustDecodeHex(t, "ac9650d8")
	if !bytes.Equal(outer[:4], wantSelector) {
		t.Fatalf("multicall selector mismatch: got %x, want %x", outer[:4], wantSelector)
	}

	// Decode the outer multicall payload and verify we recover the same two
	// inner calldata blobs in order.
	method, err := parsedABI.MethodById(outer[:4])
	if err != nil {
		t.Fatalf("MethodById: %v", err)
	}
	if method.Name != "multicall" {
		t.Fatalf("unexpected method: %s", method.Name)
	}
	unpacked, err := method.Inputs.Unpack(outer[4:])
	if err != nil {
		t.Fatalf("unpack inputs: %v", err)
	}
	if len(unpacked) != 1 {
		t.Fatalf("expected 1 top-level arg, got %d", len(unpacked))
	}
	gotInner, ok := unpacked[0].([][]byte)
	if !ok {
		t.Fatalf("unexpected inner type: %T", unpacked[0])
	}
	if len(gotInner) != 2 {
		t.Fatalf("expected 2 inner calls, got %d", len(gotInner))
	}
	if !bytes.Equal(gotInner[0], recv) {
		t.Fatalf("inner[0] mismatch")
	}
	if !bytes.Equal(gotInner[1], ack) {
		t.Fatalf("inner[1] mismatch")
	}
}

// TestSelectorForUpdateApplicationState pins the selectors for the split
// ICS26Router update entry points, which handleCosmos batches into the same
// multicall as recvPacket/ackPacket/timeoutPacket.
func TestSelectorForUpdateApplicationState(t *testing.T) {
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	for _, tc := range []struct{ method, want string }{
		{"updateApplicationState", "9c11bece"},
		{"updateConsensusState", "d395af7f"},
	} {
		data, err := parsedABI.Pack(tc.method, "cosmoshub-1", []byte{0xde, 0xad, 0xbe, 0xef})
		if err != nil {
			t.Fatalf("pack %s: %v", tc.method, err)
		}
		if len(data) < 4 {
			t.Fatalf("%s calldata too short: %d", tc.method, len(data))
		}
		if got := hex.EncodeToString(data[:4]); got != tc.want {
			t.Fatalf("%s selector mismatch: got %s, want %s", tc.method, got, tc.want)
		}
	}
}

// TestMulticallWithUpdateApplicationState verifies that an updateApplicationState
// + packet combo produces a valid multicall payload — the layout handleCosmos
// folds so the packet proof verifies against the just-applied client state.
func TestMulticallWithUpdateApplicationState(t *testing.T) {
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	pkt := samplePacket()

	updateData, err := parsedABI.Pack("updateApplicationState", "cosmoshub-1", []byte{0x01, 0x02})
	if err != nil {
		t.Fatalf("pack updateApplicationState: %v", err)
	}
	recvData, err := parsedABI.Pack("recvPacket", contractICS26Router.IICS26RouterMsgsMsgRecvPacket{
		Packet:        pkt,
		MembershipMsg: []byte{0x03},
	})
	if err != nil {
		t.Fatalf("pack recvPacket: %v", err)
	}

	outer, err := parsedABI.Pack("multicall", [][]byte{updateData, recvData})
	if err != nil {
		t.Fatalf("pack multicall: %v", err)
	}
	if !bytes.Equal(outer[:4], mustDecodeHex(t, "ac9650d8")) {
		t.Fatalf("multicall selector mismatch: got %x", outer[:4])
	}

	method, err := parsedABI.MethodById(outer[:4])
	if err != nil {
		t.Fatalf("MethodById: %v", err)
	}
	unpacked, err := method.Inputs.Unpack(outer[4:])
	if err != nil {
		t.Fatalf("unpack: %v", err)
	}
	inner := unpacked[0].([][]byte)
	if len(inner) != 2 {
		t.Fatalf("expected 2 inner calls, got %d", len(inner))
	}
	// First inner must be updateApplicationState — atomicity demands this so
	// packet proofs verify against the just-applied client state.
	if !bytes.Equal(inner[0][:4], mustDecodeHex(t, "9c11bece")) {
		t.Fatalf("inner[0] is not updateApplicationState: %x", inner[0][:4])
	}
	if !bytes.Equal(inner[1][:4], mustDecodeHex(t, "596e00b9")) {
		t.Fatalf("inner[1] is not recvPacket: %x", inner[1][:4])
	}
}

func TestErrorClassification(t *testing.T) {
	cases := []struct {
		err     error
		isNonce bool
		isKnown bool
	}{
		{errors.New("nonce too low"), true, false},
		{errors.New("old nonce"), true, false},
		{errors.New("nonce has already been used"), true, false},
		{errors.New("already known"), false, true},
		{errors.New("transaction already imported"), false, true},
		{errors.New("already exists"), false, true},
		{errors.New("execution reverted"), false, false},
		{nil, false, false},
	}
	for _, tc := range cases {
		if got := isNonceTooLowError(tc.err); got != tc.isNonce {
			t.Errorf("isNonceTooLowError(%v) = %v; want %v", tc.err, got, tc.isNonce)
		}
		if got := isAlreadyKnownError(tc.err); got != tc.isKnown {
			t.Errorf("isAlreadyKnownError(%v) = %v; want %v", tc.err, got, tc.isKnown)
		}
	}
}

type jsonrpcReq struct {
	Jsonrpc string            `json:"jsonrpc"`
	Id      json.RawMessage   `json:"id"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
}

type jsonrpcResp struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *jsonrpcError   `json:"error,omitempty"`
}

type jsonrpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type mockJSONRPC struct {
	mu           sync.Mutex
	receiptCalls []common.Hash
	receiptResps map[common.Hash]*types.Receipt
	sendErr      error
	chainID      uint64
	nonce        uint64
	nonces       map[common.Address]uint64
	nonceCalls   map[common.Address]int
	gasPrice     *big.Int
	blockGas     uint64
	blockGases   map[string]uint64
	callSucceeds bool
}

func newMockJSONRPC() *mockJSONRPC {
	return &mockJSONRPC{
		receiptResps: make(map[common.Hash]*types.Receipt),
		chainID:      1,
		nonces:       make(map[common.Address]uint64),
		nonceCalls:   make(map[common.Address]int),
		gasPrice:     big.NewInt(1000000000), // 1 Gwei
		blockGases:   make(map[string]uint64),
	}
}

func nonceCacheIsValid(h *Handler, chainID string, privateKey *ecdsa.PrivateKey) bool {
	from := crypto.PubkeyToAddress(privateKey.PublicKey)
	h.mu.Lock()
	defer h.mu.Unlock()
	state := h.senderStates[evmNonceKey{chainID: chainID, address: from}]
	return state != nil && state.nonceValid
}

func (m *mockJSONRPC) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req jsonrpcReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var result interface{}
	var rpcErr *jsonrpcError

	switch req.Method {
	case "eth_chainId":
		result = fmt.Sprintf("0x%x", m.chainID)
	case "eth_gasPrice":
		result = fmt.Sprintf("0x%x", m.gasPrice)
	case "eth_maxPriorityFeePerGas":
		result = fmt.Sprintf("0x%x", m.gasPrice)
	case "eth_getTransactionCount":
		var address common.Address
		if len(req.Params) > 0 {
			var rawAddress string
			if err := json.Unmarshal(req.Params[0], &rawAddress); err == nil {
				address = common.HexToAddress(rawAddress)
			}
		}
		m.nonceCalls[address]++
		nonce := m.nonce
		if configured, ok := m.nonces[address]; ok {
			nonce = configured
		}
		result = fmt.Sprintf("0x%x", nonce)
	case "eth_sendRawTransaction":
		if m.sendErr != nil {
			rpcErr = &jsonrpcError{
				Code:    -32000,
				Message: m.sendErr.Error(),
			}
		} else {
			result = "0x0000000000000000000000000000000000000000000000000000000000000000"
		}
	case "eth_getTransactionReceipt":
		if len(req.Params) > 0 {
			var hashStr string
			if err := json.Unmarshal(req.Params[0], &hashStr); err == nil {
				hash := common.HexToHash(hashStr)
				m.receiptCalls = append(m.receiptCalls, hash)
				if receipt, ok := m.receiptResps[hash]; ok && receipt != nil {
					result = map[string]interface{}{
						"transactionHash":   hash.Hex(),
						"transactionIndex":  "0x1",
						"blockHash":         "0x0000000000000000000000000000000000000000000000000000000000000001",
						"blockNumber":       fmt.Sprintf("0x%x", receipt.BlockNumber.Uint64()),
						"gasUsed":           fmt.Sprintf("0x%x", receipt.GasUsed),
						"cumulativeGasUsed": fmt.Sprintf("0x%x", receipt.GasUsed),
						"status":            fmt.Sprintf("0x%x", receipt.Status),
						"logsBloom":         "0x" + strings.Repeat("0", 512),
						"logs":              []interface{}{},
					}
				} else {
					result = nil
				}
			}
		}
	case "eth_getBlockByNumber":
		blockGas := m.blockGas
		if len(req.Params) > 0 {
			var blockTag string
			if err := json.Unmarshal(req.Params[0], &blockTag); err == nil {
				if configured, ok := m.blockGases[blockTag]; ok {
					blockGas = configured
				}
			}
		}
		if blockGas == 0 {
			blockGas = 0xffffff
		}
		result = map[string]interface{}{
			"parentHash":       "0x0000000000000000000000000000000000000000000000000000000000000000",
			"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			"miner":            "0x0000000000000000000000000000000000000000",
			"stateRoot":        "0x0000000000000000000000000000000000000000000000000000000000000000",
			"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"receiptsRoot":     "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"logsBloom":        "0x" + strings.Repeat("0", 512),
			"difficulty":       "0x0",
			"number":           "0x1",
			"gasLimit":         fmt.Sprintf("0x%x", blockGas),
			"gasUsed":          "0x0",
			"timestamp":        "0x0",
			"extraData":        "0x",
			"mixHash":          "0x0000000000000000000000000000000000000000000000000000000000000000",
			"nonce":            "0x0000000000000000",
			"baseFeePerGas":    fmt.Sprintf("0x%x", m.gasPrice),
		}
	case "eth_call":
		if m.callSucceeds {
			result = "0x"
		} else {
			rpcErr = &jsonrpcError{Code: -32601, Message: "method not found"}
		}
	default:
		rpcErr = &jsonrpcError{
			Code:    -32601,
			Message: "method not found",
		}
	}

	resp := jsonrpcResp{
		Jsonrpc: "2.0",
		Id:      req.Id,
		Result:  result,
		Error:   rpcErr,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func TestWaitForReceipts_ReverseOrder(t *testing.T) {
	mockRPC := newMockJSONRPC()
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	hash1 := common.HexToHash("0x1")
	hash2 := common.HexToHash("0x2")
	hash3 := common.HexToHash("0x3")

	mockRPC.receiptResps[hash3] = &types.Receipt{
		Status:      1,
		GasUsed:     21000,
		BlockNumber: big.NewInt(100),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	receipt, err := waitForReceipts(ctx, client, []common.Hash{hash1, hash2, hash3})
	if err != nil {
		t.Fatalf("waitForReceipts failed: %v", err)
	}

	if receipt == nil || receipt.BlockNumber.Uint64() != 100 {
		t.Errorf("expected receipt at block 100, got %v", receipt)
	}

	mockRPC.mu.Lock()
	defer mockRPC.mu.Unlock()
	if len(mockRPC.receiptCalls) != 1 {
		t.Errorf("expected exactly 1 call to TransactionReceipt, got %d: %v", len(mockRPC.receiptCalls), mockRPC.receiptCalls)
	} else if mockRPC.receiptCalls[0] != hash3 {
		t.Errorf("expected first call to be for %v, got %v", hash3, mockRPC.receiptCalls[0])
	}
}

func TestWaitForReceipts_ErrorPropagation(t *testing.T) {
	oldTimeout := ethTxReceiptTimeout
	oldPollInterval := ethTxReceiptPollInterval
	ethTxReceiptTimeout = 100 * time.Millisecond
	ethTxReceiptPollInterval = 10 * time.Millisecond
	defer func() {
		ethTxReceiptTimeout = oldTimeout
		ethTxReceiptPollInterval = oldPollInterval
	}()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_, err = waitForReceipts(ctx, client, []common.Hash{common.HexToHash("0x1")})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded error (since transient errors are ignored), got %v", err)
	}
}

func TestExecuteWithRetryAndResubmission_NonceRetry(t *testing.T) {
	mockRPC := newMockJSONRPC()
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.EVMEndpoint{Client: client}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{}

	var attempts int
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		attempts++
		if attempts == 1 {
			return nil, errors.New("nonce too low")
		}

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    auth.Nonce.Uint64(),
			GasPrice: auth.GasPrice,
			Gas:      auth.GasLimit,
			To:       &common.Address{0x1},
			Value:    big.NewInt(0),
			Data:     []byte{},
		})
		return auth.Signer(auth.From, tx)
	}

	var txHash common.Hash
	senderFnWrapped := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx, err := senderFn(auth)
		if err == nil && tx != nil {
			txHash = tx.Hash()
			mockRPC.mu.Lock()
			mockRPC.receiptResps[txHash] = &types.Receipt{
				Status:      1,
				GasUsed:     21000,
				BlockNumber: big.NewInt(200),
			}
			mockRPC.mu.Unlock()
		}
		return tx, err
	}

	receipt, _, _, err := h.executeWithRetryAndResubmission(context.Background(), ctx, privKey, 100000, senderFnWrapped)
	if err != nil {
		t.Fatalf("executeWithRetryAndResubmission failed: %v", err)
	}

	if receipt == nil || receipt.BlockNumber.Uint64() != 200 {
		t.Errorf("expected receipt at block 200, got %v", receipt)
	}

	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestMisbehaviourGasLimit(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    uint64
		wantErr bool
	}{
		{name: "default", want: defaultMisbehaviourGasLimit},
		{name: "override", value: "17000000", want: 17_000_000},
		{name: "zero", value: "0", wantErr: true},
		{name: "negative", value: "-1", wantErr: true},
		{name: "trailing characters", value: "17000000oops", wantErr: true},
		{name: "overflow", value: "18446744073709551616", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ETH_MISBEHAVIOUR_GAS_LIMIT", tt.value)
			got, err := MisbehaviourGasLimit()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("MisbehaviourGasLimit() = %d, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("MisbehaviourGasLimit(): %v", err)
			}
			if got != tt.want {
				t.Fatalf("MisbehaviourGasLimit() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestValidateMisbehaviourPrivateKey(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "missing", wantErr: true},
		{name: "invalid", value: "invalid", wantErr: true},
		{name: "valid", value: strings.Repeat("1", 64)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MISBEHAVIOUR_PRIVATE_KEY", tt.value)
			err := ValidateMisbehaviourPrivateKey()
			if tt.wantErr && err == nil {
				t.Fatal("ValidateMisbehaviourPrivateKey() succeeded, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("ValidateMisbehaviourPrivateKey(): %v", err)
			}
		})
	}
}

func TestExecuteWithRetryAndResubmission_ScopesNonceCachesByChainAndSender(t *testing.T) {
	relayKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate relay key: %v", err)
	}
	incidentKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate incident key: %v", err)
	}
	relayAddress := crypto.PubkeyToAddress(relayKey.PublicKey)
	incidentAddress := crypto.PubkeyToAddress(incidentKey.PublicKey)

	chainOneRPC := newMockJSONRPC()
	chainOneRPC.chainID = 1
	chainOneRPC.nonces[relayAddress] = 100
	chainOneRPC.nonces[incidentAddress] = 0
	chainOneServer := httptest.NewServer(chainOneRPC)
	defer chainOneServer.Close()
	chainOneClient, err := ethclient.Dial(chainOneServer.URL)
	if err != nil {
		t.Fatalf("dial chain one: %v", err)
	}

	chainTwoRPC := newMockJSONRPC()
	chainTwoRPC.chainID = 2
	chainTwoRPC.nonces[relayAddress] = 7
	chainTwoServer := httptest.NewServer(chainTwoRPC)
	defer chainTwoServer.Close()
	chainTwoClient, err := ethclient.Dial(chainTwoServer.URL)
	if err != nil {
		t.Fatalf("dial chain two: %v", err)
	}

	h := &Handler{}
	execute := func(name string, client *ethclient.Client, rpc *mockJSONRPC, privateKey *ecdsa.PrivateKey) uint64 {
		t.Helper()
		endpoint := services.EVMEndpoint{Client: client}
		chainID, err := client.ChainID(context.Background())
		if err != nil {
			t.Fatalf("%s chain ID: %v", name, err)
		}
		var usedNonce uint64
		sender := func(auth *bind.TransactOpts) (*types.Transaction, error) {
			usedNonce = auth.Nonce.Uint64()
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID: chainID, Nonce: usedNonce, GasTipCap: auth.GasTipCap,
				GasFeeCap: auth.GasFeeCap, Gas: auth.GasLimit, To: &common.Address{0x1},
			})
			signed, signErr := auth.Signer(auth.From, tx)
			if signErr == nil {
				rpc.mu.Lock()
				rpc.receiptResps[signed.Hash()] = &types.Receipt{Status: 1, GasUsed: 21_000, BlockNumber: big.NewInt(1)}
				rpc.mu.Unlock()
			}
			return signed, signErr
		}
		if _, _, _, err := h.executeWithRetryAndResubmission(context.Background(), endpoint, privateKey, 100_000, sender); err != nil {
			t.Fatalf("%s execute: %v", name, err)
		}
		return usedNonce
	}

	if got := execute("relay chain one", chainOneClient, chainOneRPC, relayKey); got != 100 {
		t.Fatalf("relay chain-one nonce = %d, want 100", got)
	}
	if got := execute("incident chain one", chainOneClient, chainOneRPC, incidentKey); got != 0 {
		t.Fatalf("incident chain-one nonce = %d, want 0", got)
	}
	if got := execute("relay chain two", chainTwoClient, chainTwoRPC, relayKey); got != 7 {
		t.Fatalf("relay chain-two nonce = %d, want 7", got)
	}
	if chainOneRPC.nonceCalls[relayAddress] != 1 || chainOneRPC.nonceCalls[incidentAddress] != 1 {
		t.Fatalf("chain-one pending nonce calls relay/incident = %d/%d, want 1/1",
			chainOneRPC.nonceCalls[relayAddress], chainOneRPC.nonceCalls[incidentAddress])
	}
}

func TestExecuteWithRetryAndResubmission_AlreadyKnown(t *testing.T) {
	mockRPC := newMockJSONRPC()
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.EVMEndpoint{Client: client}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{}

	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    auth.Nonce.Uint64(),
			GasPrice: auth.GasPrice,
			Gas:      auth.GasLimit,
			To:       &common.Address{0x1},
			Value:    big.NewInt(0),
			Data:     []byte{},
		})
		signedTx, err := auth.Signer(auth.From, tx)
		if err != nil {
			return nil, err
		}

		mockRPC.mu.Lock()
		mockRPC.receiptResps[signedTx.Hash()] = &types.Receipt{
			Status:      1,
			GasUsed:     21000,
			BlockNumber: big.NewInt(300),
		}
		mockRPC.mu.Unlock()

		return signedTx, errors.New("already known")
	}

	receipt, _, _, err := h.executeWithRetryAndResubmission(context.Background(), ctx, privKey, 100000, senderFn)
	if err != nil {
		t.Fatalf("executeWithRetryAndResubmission failed: %v", err)
	}

	if receipt == nil || receipt.BlockNumber.Uint64() != 300 {
		t.Errorf("expected receipt at block 300, got %v", receipt)
	}
}

func TestExecuteWithRetryAndResubmission_BroadcastContextDeadline(t *testing.T) {
	oldBroadcastTimeout := ethTxBroadcastTimeout
	ethTxBroadcastTimeout = 30 * time.Millisecond
	defer func() {
		ethTxBroadcastTimeout = oldBroadcastTimeout
	}()

	mockRPC := newMockJSONRPC()
	mockRPC.nonce = 10
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.EVMEndpoint{Client: client}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{}

	var sawDeadline bool
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		if auth.Context == nil {
			return nil, errors.New("auth context is nil")
		}
		if _, ok := auth.Context.Deadline(); !ok {
			return nil, errors.New("auth context has no deadline")
		}
		sawDeadline = true

		select {
		case <-auth.Context.Done():
			return nil, auth.Context.Err()
		case <-time.After(500 * time.Millisecond):
			return nil, errors.New("auth context did not expire")
		}
	}

	start := time.Now()
	_, _, _, err = h.executeWithRetryAndResubmission(context.Background(), ctx, privKey, 100000, senderFn)
	if err == nil {
		t.Fatal("expected broadcast context deadline error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}
	if !sawDeadline {
		t.Fatal("sender did not observe a deadline on auth.Context")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("broadcast deadline took too long: %s", elapsed)
	}

	if nonceCacheIsValid(h, "1", privKey) {
		t.Error("expected nonce cache invalidated after broadcast deadline")
	}
}

func TestBumpGasAndResubmit(t *testing.T) {
	mockRPC := newMockJSONRPC()
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.EVMEndpoint{Client: client}

	h := &Handler{}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	chainIdInt, err := client.ChainID(context.Background())
	if err != nil {
		t.Fatalf("chain ID: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainIdInt)
	if err != nil {
		t.Fatalf("transactor: %v", err)
	}

	// 1. Legacy Transaction
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		GasPrice: big.NewInt(1000000000), // 1 Gwei
		Gas:      21000,
		To:       &common.Address{0x1},
		Value:    big.NewInt(0),
		Data:     []byte{},
	})

	bumpedTx, err := h.bumpGasAndResubmit(context.Background(), ctx, tx, auth, 1)
	if err != nil {
		t.Fatalf("bumpGasAndResubmit failed: %v", err)
	}

	expectedGasPrice := big.NewInt(1150000000)
	if bumpedTx.GasPrice().Cmp(expectedGasPrice) != 0 {
		t.Errorf("expected gas price %v, got %v", expectedGasPrice, bumpedTx.GasPrice())
	}

	// 2. Dynamic Fee Transaction
	dynamicTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainIdInt,
		Nonce:     0,
		GasTipCap: big.NewInt(1000000000), // 1 Gwei
		GasFeeCap: big.NewInt(2000000000), // 2 Gwei
		Gas:       21000,
		To:        &common.Address{0x1},
		Value:     big.NewInt(0),
		Data:      []byte{},
	})

	bumpedDynamicTx, err := h.bumpGasAndResubmit(context.Background(), ctx, dynamicTx, auth, 1)
	if err != nil {
		t.Fatalf("bumpGasAndResubmit (dynamic) failed: %v", err)
	}

	expectedTipCap := big.NewInt(1150000000)
	expectedFeeCap := big.NewInt(2300000000)
	if bumpedDynamicTx.GasTipCap().Cmp(expectedTipCap) != 0 {
		t.Errorf("expected tip cap %v, got %v", expectedTipCap, bumpedDynamicTx.GasTipCap())
	}
	if bumpedDynamicTx.GasFeeCap().Cmp(expectedFeeCap) != 0 {
		t.Errorf("expected fee cap %v, got %v", expectedFeeCap, bumpedDynamicTx.GasFeeCap())
	}

	// 3. Access List Transaction
	accessListTx := types.NewTx(&types.AccessListTx{
		ChainID:    chainIdInt,
		Nonce:      0,
		GasPrice:   big.NewInt(1000000000), // 1 Gwei
		Gas:        21000,
		To:         &common.Address{0x1},
		Value:      big.NewInt(0),
		Data:       []byte{},
		AccessList: types.AccessList{},
	})

	bumpedAccessListTx, err := h.bumpGasAndResubmit(context.Background(), ctx, accessListTx, auth, 1)
	if err != nil {
		t.Fatalf("bumpGasAndResubmit (access list) failed: %v", err)
	}

	if bumpedAccessListTx.Type() != types.AccessListTxType {
		t.Errorf("expected AccessListTxType, got %d", bumpedAccessListTx.Type())
	}
	expectedALPrice := big.NewInt(1150000000)
	if bumpedAccessListTx.GasPrice().Cmp(expectedALPrice) != 0 {
		t.Errorf("expected gas price %v, got %v", expectedALPrice, bumpedAccessListTx.GasPrice())
	}
}

func TestExecuteWithRetryAndResubmission_WaitErrorNonceInvalidation(t *testing.T) {
	oldTimeout := ethTxReceiptTimeout
	oldPollInterval := ethTxReceiptPollInterval
	ethTxReceiptTimeout = 100 * time.Millisecond
	ethTxReceiptPollInterval = 10 * time.Millisecond
	defer func() {
		ethTxReceiptTimeout = oldTimeout
		ethTxReceiptPollInterval = oldPollInterval
	}()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req jsonrpcReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var result interface{}
		var rpcErr *jsonrpcError

		switch req.Method {
		case "eth_chainId":
			result = "0x1"
		case "eth_gasPrice":
			result = "0x3b9aca00" // 1 Gwei
		case "eth_getTransactionCount":
			result = "0xa" // 10
		case "eth_sendRawTransaction":
			result = "0x0000000000000000000000000000000000000000000000000000000000000000"
		case "eth_getTransactionReceipt":
			// Fail the receipt lookup with an RPC error
			rpcErr = &jsonrpcError{
				Code:    -32000,
				Message: "node syncing / receipt unavailable due to internal error",
			}
		default:
			rpcErr = &jsonrpcError{
				Code:    -32601,
				Message: "method not found",
			}
		}

		resp := jsonrpcResp{
			Jsonrpc: "2.0",
			Id:      req.Id,
			Result:  result,
			Error:   rpcErr,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.EVMEndpoint{Client: client}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{}

	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    auth.Nonce.Uint64(),
			GasPrice: auth.GasPrice,
			Gas:      auth.GasLimit,
			To:       &common.Address{0x1},
			Value:    big.NewInt(0),
			Data:     []byte{},
		})
		return auth.Signer(auth.From, tx)
	}

	_, _, _, err = h.executeWithRetryAndResubmission(context.Background(), ctx, privKey, 100000, senderFn)
	if err == nil {
		t.Fatal("expected error from executeWithRetryAndResubmission due to receipt wait failure, got nil")
	}

	if nonceCacheIsValid(h, "1", privKey) {
		t.Error("expected nonce cache to be invalid after receipt wait error")
	}
}

func TestExecuteWithRetryAndResubmission_GasFloorsOnStuckNonce(t *testing.T) {
	mockRPC := newMockJSONRPC()
	mockRPC.nonce = 10
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.EVMEndpoint{Client: client}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	chainIdInt, err := client.ChainID(context.Background())
	if err != nil {
		t.Fatalf("chain ID: %v", err)
	}

	h := &Handler{}

	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainIdInt,
			Nonce:     auth.Nonce.Uint64(),
			GasTipCap: auth.GasTipCap,
			GasFeeCap: auth.GasFeeCap,
			Gas:       auth.GasLimit,
			To:        &common.Address{0x1},
			Value:     big.NewInt(0),
			Data:      []byte{},
		})
		return auth.Signer(auth.From, tx)
	}

	var tx1 *types.Transaction
	senderFnWrapped1 := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx, err := senderFn(auth)
		if err == nil {
			tx1 = tx
			mockRPC.mu.Lock()
			mockRPC.receiptResps[tx.Hash()] = &types.Receipt{
				Status:      1,
				GasUsed:     21000,
				BlockNumber: big.NewInt(100),
			}
			mockRPC.mu.Unlock()
		}
		return tx, err
	}

	_, _, _, err = h.executeWithRetryAndResubmission(context.Background(), ctx, privKey, 100000, senderFnWrapped1)
	if err != nil {
		t.Fatalf("first execution failed: %v", err)
	}

	if tx1 == nil {
		t.Fatal("first transaction was not captured")
	}

	// Invalidate nonce cache manually (simulating timeout that occurred)
	h.mu.Lock()
	h.senderState("1", crypto.PubkeyToAddress(privKey.PublicKey)).nonceValid = false
	h.mu.Unlock()

	// 2. Second execution: re-submitting at same nonce
	var tx2 *types.Transaction
	senderFnWrapped2 := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx, err := senderFn(auth)
		if err == nil {
			tx2 = tx
			mockRPC.mu.Lock()
			mockRPC.receiptResps[tx.Hash()] = &types.Receipt{
				Status:      1,
				GasUsed:     21000,
				BlockNumber: big.NewInt(101),
			}
			mockRPC.mu.Unlock()
		}
		return tx, err
	}

	_, _, _, err = h.executeWithRetryAndResubmission(context.Background(), ctx, privKey, 100000, senderFnWrapped2)
	if err != nil {
		t.Fatalf("second execution failed: %v", err)
	}

	if tx2 == nil {
		t.Fatal("second transaction was not captured")
	}

	// Verify gas floors were applied (15% bump)
	expectedTipCap := new(big.Int).Mul(tx1.GasTipCap(), big.NewInt(115))
	expectedTipCap.Div(expectedTipCap, big.NewInt(100))

	expectedFeeCap := new(big.Int).Mul(tx1.GasFeeCap(), big.NewInt(115))
	expectedFeeCap.Div(expectedFeeCap, big.NewInt(100))

	if tx2.GasTipCap().Cmp(expectedTipCap) != 0 {
		t.Errorf("expected bumped tip cap %v, got %v", expectedTipCap, tx2.GasTipCap())
	}

	if tx2.GasFeeCap().Cmp(expectedFeeCap) != 0 {
		t.Errorf("expected bumped fee cap %v, got %v", expectedFeeCap, tx2.GasFeeCap())
	}
}

func TestExecuteWithRetryAndResubmission_SenderFnReturnsNil(t *testing.T) {
	mockRPC := newMockJSONRPC()
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.EVMEndpoint{Client: client}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{}

	// senderFn returns (nil, nil)
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return nil, nil
	}

	_, _, _, err = h.executeWithRetryAndResubmission(context.Background(), ctx, privKey, 100000, senderFn)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "senderFn returned nil transaction without error"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected error to contain %q, got %q", expectedErr, err.Error())
	}

	if nonceCacheIsValid(h, "1", privKey) {
		t.Error("expected nonce cache to be invalid after nil transaction error")
	}
}

func TestExecuteWithRetryAndResubmission_Concurrency(t *testing.T) {
	mockRPC := newMockJSONRPC()
	mockRPC.nonce = 10
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.EVMEndpoint{Client: client}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{}

	numRequests := 10
	var mu sync.Mutex
	noncesUsed := make([]uint64, 0, numRequests)
	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
				// Simulate random latency to trigger concurrency races if they exist
				time.Sleep(time.Duration(10+time.Now().UnixNano()%30) * time.Millisecond)

				tx := types.NewTx(&types.LegacyTx{
					Nonce:    auth.Nonce.Uint64(),
					GasPrice: auth.GasPrice,
					Gas:      auth.GasLimit,
					To:       &common.Address{0x1},
					Value:    big.NewInt(0),
					Data:     []byte{},
				})
				signedTx, signErr := auth.Signer(auth.From, tx)
				if signErr != nil {
					return nil, signErr
				}
				mockRPC.mu.Lock()
				mockRPC.receiptResps[signedTx.Hash()] = &types.Receipt{
					Status:      1,
					GasUsed:     21000,
					BlockNumber: big.NewInt(100),
				}
				mockRPC.mu.Unlock()

				mu.Lock()
				noncesUsed = append(noncesUsed, auth.Nonce.Uint64())
				mu.Unlock()

				return signedTx, nil
			}

			_, _, _, err := h.executeWithRetryAndResubmission(context.Background(), ctx, privKey, 100000, senderFn)
			if err != nil {
				t.Errorf("executeWithRetryAndResubmission failed: %v", err)
				return
			}
		}()
	}

	wg.Wait()

	if len(noncesUsed) != numRequests {
		t.Fatalf("expected %d nonces, got %d", numRequests, len(noncesUsed))
	}

	seen := make(map[uint64]bool)
	for _, nonce := range noncesUsed {
		if nonce < 10 || nonce >= 20 {
			t.Errorf("nonce %d out of bounds [10, 19]", nonce)
		}
		if seen[nonce] {
			t.Errorf("duplicate nonce used: %d", nonce)
		}
		seen[nonce] = true
	}
}

func TestExecuteWithRetryAndResubmission_RevertPermanentFailure(t *testing.T) {
	mockRPC := newMockJSONRPC()
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.EVMEndpoint{Client: client}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{}

	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    auth.Nonce.Uint64(),
			GasPrice: auth.GasPrice,
			Gas:      auth.GasLimit,
			To:       &common.Address{0x1},
			Value:    big.NewInt(0),
			Data:     []byte{},
		})
		signedTx, signErr := auth.Signer(auth.From, tx)
		if signErr != nil {
			return nil, signErr
		}
		mockRPC.mu.Lock()
		mockRPC.receiptResps[signedTx.Hash()] = &types.Receipt{
			Status:      0, // Reverted!
			GasUsed:     21000,
			BlockNumber: big.NewInt(100),
		}
		mockRPC.mu.Unlock()
		return signedTx, nil
	}

	_, _, _, err = h.executeWithRetryAndResubmission(context.Background(), ctx, privKey, 100000, senderFn)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, services.ErrPermanentRelayFailure) {
		t.Errorf("expected ErrPermanentRelayFailure, got %v", err)
	}
}

func TestExecuteWithRetryAndResubmission_OutOfGasEscalatesInProcess(t *testing.T) {
	oldPollInterval := ethTxReceiptPollInterval
	ethTxReceiptPollInterval = time.Millisecond
	defer func() { ethTxReceiptPollInterval = oldPollInterval }()

	mockRPC := newMockJSONRPC()
	mockRPC.callSucceeds = true
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()
	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("dial mock RPC: %v", err)
	}
	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	var gasLimits, nonces []uint64
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		gasLimits = append(gasLimits, auth.GasLimit)
		nonces = append(nonces, auth.Nonce.Uint64())
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID: big.NewInt(1), Nonce: auth.Nonce.Uint64(), GasTipCap: auth.GasTipCap,
			GasFeeCap: auth.GasFeeCap, Gas: auth.GasLimit, To: &common.Address{0x1},
		})
		signed, signErr := auth.Signer(auth.From, tx)
		if signErr == nil {
			status := uint64(1)
			if len(gasLimits) == 1 {
				status = 0
			}
			mockRPC.mu.Lock()
			mockRPC.receiptResps[signed.Hash()] = &types.Receipt{Status: status, GasUsed: auth.GasLimit, BlockNumber: big.NewInt(1)}
			mockRPC.mu.Unlock()
		}
		return signed, signErr
	}

	receipt, _, _, err := (&Handler{}).executeWithRetryAndResubmission(
		context.Background(), services.EVMEndpoint{Client: client}, privKey, 100_000, senderFn, knobEthGasLimit,
	)
	if err != nil {
		t.Fatalf("execute after OOG: %v", err)
	}
	if receipt == nil || receipt.Status != 1 {
		t.Fatalf("final receipt = %v, want success", receipt)
	}
	if len(gasLimits) != 2 || gasLimits[1] <= gasLimits[0] {
		t.Fatalf("gas limits = %v, want one larger in-process retry", gasLimits)
	}
	if len(nonces) != 2 || nonces[1] != nonces[0]+1 {
		t.Fatalf("nonces = %v, want included OOG to consume one nonce", nonces)
	}
}

func TestExecuteWithRetryAndResubmission_OutOfGasLadderExhaustionIsPermanent(t *testing.T) {
	oldPollInterval := ethTxReceiptPollInterval
	ethTxReceiptPollInterval = time.Millisecond
	defer func() { ethTxReceiptPollInterval = oldPollInterval }()

	mockRPC := newMockJSONRPC()
	mockRPC.callSucceeds = true
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()
	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("dial mock RPC: %v", err)
	}
	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	var gasLimits []uint64
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		gasLimits = append(gasLimits, auth.GasLimit)
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID: big.NewInt(1), Nonce: auth.Nonce.Uint64(), GasTipCap: auth.GasTipCap,
			GasFeeCap: auth.GasFeeCap, Gas: auth.GasLimit, To: &common.Address{0x1},
		})
		signed, signErr := auth.Signer(auth.From, tx)
		if signErr == nil {
			mockRPC.mu.Lock()
			mockRPC.receiptResps[signed.Hash()] = &types.Receipt{Status: 0, GasUsed: auth.GasLimit, BlockNumber: big.NewInt(1)}
			mockRPC.mu.Unlock()
		}
		return signed, signErr
	}

	_, _, _, err = (&Handler{}).executeWithRetryAndResubmission(
		context.Background(), services.EVMEndpoint{Client: client}, privKey, 100_000, senderFn, knobEthGasLimit,
	)
	if !errors.Is(err, services.ErrPermanentRelayFailure) {
		t.Fatalf("exhausted OOG ladder error = %v, want permanent relay failure", err)
	}
	if len(gasLimits) != len(evmGasHeadroomBasisPoints) {
		t.Fatalf("gas attempts = %v, want finite ladder length %d", gasLimits, len(evmGasHeadroomBasisPoints))
	}
	for i := 1; i < len(gasLimits); i++ {
		if gasLimits[i] <= gasLimits[i-1] {
			t.Fatalf("gas attempts = %v, want strictly increasing limits", gasLimits)
		}
	}
}

func TestExecuteWithRetryAndResubmission_OutOfGasBlockCeilingIsPermanent(t *testing.T) {
	oldPollInterval := ethTxReceiptPollInterval
	ethTxReceiptPollInterval = time.Millisecond
	defer func() { ethTxReceiptPollInterval = oldPollInterval }()

	mockRPC := newMockJSONRPC()
	mockRPC.callSucceeds = true
	mockRPC.blockGas = 120_000
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()
	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("dial mock RPC: %v", err)
	}
	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	var gasLimits []uint64
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		gasLimits = append(gasLimits, auth.GasLimit)
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID: big.NewInt(1), Nonce: auth.Nonce.Uint64(), GasTipCap: auth.GasTipCap,
			GasFeeCap: auth.GasFeeCap, Gas: auth.GasLimit, To: &common.Address{0x1},
		})
		signed, signErr := auth.Signer(auth.From, tx)
		if signErr == nil {
			mockRPC.mu.Lock()
			mockRPC.receiptResps[signed.Hash()] = &types.Receipt{Status: 0, GasUsed: auth.GasLimit, BlockNumber: big.NewInt(1)}
			mockRPC.mu.Unlock()
		}
		return signed, signErr
	}

	_, _, _, err = (&Handler{}).executeWithRetryAndResubmission(
		context.Background(), services.EVMEndpoint{Client: client}, privKey, 100_000, senderFn, knobEthGasLimit,
	)
	if !errors.Is(err, services.ErrPermanentRelayFailure) {
		t.Fatalf("block-ceiling OOG error = %v, want permanent relay failure", err)
	}
	if !strings.Contains(err.Error(), "block gas ceiling") {
		t.Fatalf("block-ceiling OOG error = %v, want the block-ceiling termination reason", err)
	}
	if want := []uint64{100_000, 120_000}; !slices.Equal(gasLimits, want) {
		t.Fatalf("gas attempts = %v, want %v without re-broadcasting the clamped limit", gasLimits, want)
	}
}

func TestExecuteWithRetryAndResubmission_RefreshesBlockGasCeiling(t *testing.T) {
	oldPollInterval := ethTxReceiptPollInterval
	ethTxReceiptPollInterval = time.Millisecond
	defer func() { ethTxReceiptPollInterval = oldPollInterval }()

	mockRPC := newMockJSONRPC()
	mockRPC.callSucceeds = true
	mockRPC.blockGases["0x1"] = 120_000
	mockRPC.blockGases["0x2"] = 150_000
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()
	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("dial mock RPC: %v", err)
	}
	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	var gasLimits []uint64
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		gasLimits = append(gasLimits, auth.GasLimit)
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID: big.NewInt(1), Nonce: auth.Nonce.Uint64(), GasTipCap: auth.GasTipCap,
			GasFeeCap: auth.GasFeeCap, Gas: auth.GasLimit, To: &common.Address{0x1},
		})
		signed, signErr := auth.Signer(auth.From, tx)
		if signErr == nil {
			status := uint64(0)
			if len(gasLimits) == 3 {
				status = 1
			}
			mockRPC.mu.Lock()
			mockRPC.receiptResps[signed.Hash()] = &types.Receipt{
				Status: status, GasUsed: auth.GasLimit, BlockNumber: big.NewInt(int64(len(gasLimits))),
			}
			mockRPC.mu.Unlock()
		}
		return signed, signErr
	}

	receipt, _, _, err := (&Handler{}).executeWithRetryAndResubmission(
		context.Background(), services.EVMEndpoint{Client: client}, privKey, 100_000, senderFn, knobEthGasLimit,
	)
	if err != nil {
		t.Fatalf("execute after block gas ceiling increased: %v", err)
	}
	if receipt == nil || receipt.Status != 1 {
		t.Fatalf("final receipt = %v, want success", receipt)
	}
	if want := []uint64{100_000, 120_000, 150_000}; !slices.Equal(gasLimits, want) {
		t.Fatalf("gas attempts = %v, want refreshed ceilings %v", gasLimits, want)
	}
}

func TestExecuteWithRetryAndResubmission_BlockCeilingWinsAtFinalLadderStep(t *testing.T) {
	oldPollInterval := ethTxReceiptPollInterval
	ethTxReceiptPollInterval = time.Millisecond
	defer func() { ethTxReceiptPollInterval = oldPollInterval }()

	mockRPC := newMockJSONRPC()
	mockRPC.callSucceeds = true
	mockRPC.blockGas = 300_000
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()
	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("dial mock RPC: %v", err)
	}
	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	var gasLimits []uint64
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		gasLimits = append(gasLimits, auth.GasLimit)
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID: big.NewInt(1), Nonce: auth.Nonce.Uint64(), GasTipCap: auth.GasTipCap,
			GasFeeCap: auth.GasFeeCap, Gas: auth.GasLimit, To: &common.Address{0x1},
		})
		signed, signErr := auth.Signer(auth.From, tx)
		if signErr == nil {
			mockRPC.mu.Lock()
			mockRPC.receiptResps[signed.Hash()] = &types.Receipt{Status: 0, GasUsed: auth.GasLimit, BlockNumber: big.NewInt(1)}
			mockRPC.mu.Unlock()
		}
		return signed, signErr
	}

	finalStep := len(evmGasHeadroomBasisPoints) - 1
	_, _, _, err = (&Handler{}).executeWithRetryAndResubmissionAtGasStep(
		context.Background(), services.EVMEndpoint{Client: client}, privKey,
		100_000, finalStep, 300_000, senderFn, knobEthGasLimit,
	)
	if !errors.Is(err, services.ErrPermanentRelayFailure) {
		t.Fatalf("final-step block-ceiling OOG error = %v, want permanent relay failure", err)
	}
	if !strings.Contains(err.Error(), "block gas ceiling") {
		t.Fatalf("final-step block-ceiling OOG error = %v, want ceiling to win over ladder exhaustion", err)
	}
	if want := []uint64{300_000}; !slices.Equal(gasLimits, want) {
		t.Fatalf("gas attempts = %v, want %v without another final-step broadcast", gasLimits, want)
	}
}

func TestExecuteWithRetryAndResubmission_TransientWaitFailure(t *testing.T) {
	oldTimeout := ethTxReceiptTimeout
	oldPollInterval := ethTxReceiptPollInterval
	ethTxReceiptTimeout = 100 * time.Millisecond
	ethTxReceiptPollInterval = 10 * time.Millisecond
	defer func() {
		ethTxReceiptTimeout = oldTimeout
		ethTxReceiptPollInterval = oldPollInterval
	}()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req jsonrpcReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var result interface{}
		var rpcErr *jsonrpcError

		switch req.Method {
		case "eth_chainId":
			result = "0x1"
		case "eth_gasPrice":
			result = "0x3b9aca00" // 1 Gwei
		case "eth_getTransactionCount":
			result = "0xa" // 10
		case "eth_sendRawTransaction":
			result = "0x0000000000000000000000000000000000000000000000000000000000000000"
		case "eth_getTransactionReceipt":
			// Fail the receipt lookup with an RPC error (transient error)
			rpcErr = &jsonrpcError{
				Code:    -32000,
				Message: "node syncing / receipt unavailable due to internal error",
			}
		default:
			rpcErr = &jsonrpcError{
				Code:    -32601,
				Message: "method not found",
			}
		}

		resp := jsonrpcResp{
			Jsonrpc: "2.0",
			Id:      req.Id,
			Result:  result,
			Error:   rpcErr,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.EVMEndpoint{Client: client}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{}

	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    auth.Nonce.Uint64(),
			GasPrice: auth.GasPrice,
			Gas:      auth.GasLimit,
			To:       &common.Address{0x1},
			Value:    big.NewInt(0),
			Data:     []byte{},
		})
		return auth.Signer(auth.From, tx)
	}

	_, _, _, err = h.executeWithRetryAndResubmission(context.Background(), ctx, privKey, 100000, senderFn)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if errors.Is(err, services.ErrPermanentRelayFailure) {
		t.Errorf("expected transient error, but errors.Is(err, ErrPermanentRelayFailure) was true: %v", err)
	}
}
func TestExecuteWithRetryAndResubmission_SuggestGasPriceErrorNonceInvalidation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req jsonrpcReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var result interface{}
		var rpcErr *jsonrpcError

		switch req.Method {
		case "eth_chainId":
			result = "0x1"
		case "eth_getTransactionCount":
			result = "0xa" // 10
		case "eth_gasPrice":
			// Fail gas price suggestion
			rpcErr = &jsonrpcError{
				Code:    -32000,
				Message: "cannot suggest gas price",
			}
		default:
			rpcErr = &jsonrpcError{
				Code:    -32601,
				Message: "method not found",
			}
		}

		resp := jsonrpcResp{
			Jsonrpc: "2.0",
			Id:      req.Id,
			Result:  result,
			Error:   rpcErr,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.EVMEndpoint{Client: client}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{}

	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    auth.Nonce.Uint64(),
			GasPrice: auth.GasPrice,
			Gas:      auth.GasLimit,
			To:       &common.Address{0x1},
			Value:    big.NewInt(0),
			Data:     []byte{},
		})
		return auth.Signer(auth.From, tx)
	}

	_, _, _, err = h.executeWithRetryAndResubmission(context.Background(), ctx, privKey, 100000, senderFn)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if nonceCacheIsValid(h, "1", privKey) {
		t.Error("expected nonce cache to be invalid after SuggestGasPrice failure")
	}
}

// TestWaitForTxResult_CancelledCtxAbortsPromptly is the regression for the #252 P2
// finding: the poll loop must abort on shutdown instead of sleeping past
// cancellation until the local timeout — which, held under cosmosMu in
// SendCosmosTxBatch, would stall shutdown for the whole inclusion window. A
// pre-cancelled ctx returns on the first iteration before any CosmosClient poll, so
// no live client is needed.
func TestWaitForTxResult_CancelledCtxAbortsPromptly(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	h := &Handler{}
	start := time.Now()
	_, err := h.waitForTxResult(ctx, services.CosmosEndpoint{}, []byte{0x01}, cosmosInclusionTimeout)
	if err == nil {
		t.Fatal("expected a cancellation error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want wrapped context.Canceled", err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("waitForTxResult took %s; must abort promptly on cancel, not run to the inclusion timeout", elapsed)
	}
}
