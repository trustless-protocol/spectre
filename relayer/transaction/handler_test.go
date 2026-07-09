package transaction

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
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
// services.Context + ethclient mock and is validated end-to-end via the
// benchmark run described in /Users/ducnt/.claude/plans/...
//
// Selector reference (from `cast sig` against the ICS26Router ABI):
//   recvPacket(...)              → 0x596e00b9
//   ackPacket(...)               → 0xfdbd955d
//   timeoutPacket(...)           → 0x223e357a
//   updateClient(string,bytes)   → 0x6fbf8079
//   reAnchorPinnedSet(...)       → 0x08ffd30f
//   multicall(bytes[])           → 0xac9650d8

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

// TestSelectorForUpdateClient pins the selector for ICS26Router.updateClient,
// which V2 batches into the same multicall as recvPacket/ackPacket/
// timeoutPacket.
func TestSelectorForUpdateClient(t *testing.T) {
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	data, err := parsedABI.Pack("updateClient", "cosmoshub-1", []byte{0xde, 0xad, 0xbe, 0xef})
	if err != nil {
		t.Fatalf("pack: %v", err)
	}
	if len(data) < 4 {
		t.Fatalf("calldata too short: %d", len(data))
	}
	got := hex.EncodeToString(data[:4])
	if got != "6fbf8079" {
		t.Fatalf("selector mismatch: got %s, want 6fbf8079", got)
	}
}

func TestSelectorForReAnchorPinnedSet(t *testing.T) {
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	// The selector is the first 4 bytes of any Pack(...) of this method, i.e.
	// Method.ID. Reading it directly avoids packing sample args, which would
	// panic in the ABI encoder on the empty (nil-slice) nested tuples.
	method, ok := parsedABI.Methods["reAnchorPinnedSet"]
	if !ok {
		t.Fatal("reAnchorPinnedSet not found in ICS26Router ABI")
	}
	got := hex.EncodeToString(method.ID)
	if got != "08ffd30f" {
		t.Fatalf("selector mismatch: got %s, want 08ffd30f", got)
	}
}

// TestMulticallWithUpdateClient verifies that an updateClient + packet combo
// produces a valid multicall payload — the layout we expect from the V2
// handleCosmos refactor.
func TestMulticallWithUpdateClient(t *testing.T) {
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	pkt := samplePacket()

	updateData, err := parsedABI.Pack("updateClient", "cosmoshub-1", []byte{0x01, 0x02})
	if err != nil {
		t.Fatalf("pack updateClient: %v", err)
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
	// First inner must be updateClient — atomicity demands this so packet
	// proofs verify against the just-applied client state.
	if !bytes.Equal(inner[0][:4], mustDecodeHex(t, "6fbf8079")) {
		t.Fatalf("inner[0] is not updateClient: %x", inner[0][:4])
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
	nonce        uint64
	gasPrice     *big.Int
}

func newMockJSONRPC() *mockJSONRPC {
	return &mockJSONRPC{
		receiptResps: make(map[common.Hash]*types.Receipt),
		gasPrice:     big.NewInt(1000000000), // 1 Gwei
	}
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
		result = "0x1"
	case "eth_gasPrice":
		result = fmt.Sprintf("0x%x", m.gasPrice)
	case "eth_maxPriorityFeePerGas":
		result = fmt.Sprintf("0x%x", m.gasPrice)
	case "eth_getTransactionCount":
		result = fmt.Sprintf("0x%x", m.nonce)
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
			"gasLimit":         "0xffffff",
			"gasUsed":          "0x0",
			"timestamp":        "0x0",
			"extraData":        "0x",
			"mixHash":          "0x0000000000000000000000000000000000000000000000000000000000000000",
			"nonce":            "0x0000000000000000",
			"baseFeePerGas":    fmt.Sprintf("0x%x", m.gasPrice),
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

	ctx := services.NewCtx(nil, client)

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

	receipt, _, _, err := h.executeWithRetryAndResubmission(ctx, privKey, 100000, senderFnWrapped)
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

func TestExecuteWithRetryAndResubmission_AlreadyKnown(t *testing.T) {
	mockRPC := newMockJSONRPC()
	srv := httptest.NewServer(mockRPC)
	defer srv.Close()

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	ctx := services.NewCtx(nil, client)

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

	receipt, _, _, err := h.executeWithRetryAndResubmission(ctx, privKey, 100000, senderFn)
	if err != nil {
		t.Fatalf("executeWithRetryAndResubmission failed: %v", err)
	}

	if receipt == nil || receipt.BlockNumber.Uint64() != 300 {
		t.Errorf("expected receipt at block 300, got %v", receipt)
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

	ctx := services.NewCtx(nil, client)

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

	bumpedTx, err := h.bumpGasAndResubmit(ctx, tx, auth, 1)
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

	bumpedDynamicTx, err := h.bumpGasAndResubmit(ctx, dynamicTx, auth, 1)
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

	bumpedAccessListTx, err := h.bumpGasAndResubmit(ctx, accessListTx, auth, 1)
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

	ctx := services.NewCtx(nil, client)

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{
		nonce:      10,
		nonceValid: true,
	}

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

	_, _, _, err = h.executeWithRetryAndResubmission(ctx, privKey, 100000, senderFn)
	if err == nil {
		t.Fatal("expected error from executeWithRetryAndResubmission due to receipt wait failure, got nil")
	}

	h.mu.Lock()
	valid := h.nonceValid
	h.mu.Unlock()

	if valid {
		t.Error("expected h.nonceValid to be false after receipt wait error, but it was true")
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

	ctx := services.NewCtx(nil, client)

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	chainIdInt, err := client.ChainID(context.Background())
	if err != nil {
		t.Fatalf("chain ID: %v", err)
	}

	h := &Handler{
		nonce:      10,
		nonceValid: true,
	}

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

	_, _, _, err = h.executeWithRetryAndResubmission(ctx, privKey, 100000, senderFnWrapped1)
	if err != nil {
		t.Fatalf("first execution failed: %v", err)
	}

	if tx1 == nil {
		t.Fatal("first transaction was not captured")
	}

	// Invalidate nonce cache manually (simulating timeout that occurred)
	h.mu.Lock()
	h.nonceValid = false
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

	_, _, _, err = h.executeWithRetryAndResubmission(ctx, privKey, 100000, senderFnWrapped2)
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

	ctx := services.NewCtx(nil, client)

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{
		nonce:      10,
		nonceValid: true,
	}

	// senderFn returns (nil, nil)
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return nil, nil
	}

	_, _, _, err = h.executeWithRetryAndResubmission(ctx, privKey, 100000, senderFn)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "senderFn returned nil transaction without error"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected error to contain %q, got %q", expectedErr, err.Error())
	}

	h.mu.Lock()
	valid := h.nonceValid
	h.mu.Unlock()

	if valid {
		t.Error("expected h.nonceValid to be false after nil transaction error")
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

	ctx := services.NewCtx(nil, client)

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{
		nonce:      10,
		nonceValid: true,
	}

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

			_, _, _, err := h.executeWithRetryAndResubmission(ctx, privKey, 100000, senderFn)
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

	ctx := services.NewCtx(nil, client)

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

	_, _, _, err = h.executeWithRetryAndResubmission(ctx, privKey, 100000, senderFn)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, services.ErrPermanentRelayFailure) {
		t.Errorf("expected ErrPermanentRelayFailure, got %v", err)
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

	ctx := services.NewCtx(nil, client)

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

	_, _, _, err = h.executeWithRetryAndResubmission(ctx, privKey, 100000, senderFn)
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

	ctx := services.NewCtx(nil, client)

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	h := &Handler{
		nonce:      10,
		nonceValid: true,
	}

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

	_, _, _, err = h.executeWithRetryAndResubmission(ctx, privKey, 100000, senderFn)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	h.mu.Lock()
	valid := h.nonceValid
	h.mu.Unlock()

	if valid {
		t.Error("expected h.nonceValid to be false after SuggestGasPrice failure, but it was true")
	}
}
