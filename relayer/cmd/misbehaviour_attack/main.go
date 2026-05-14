package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	contractGroth16ICS07Tendermint "relayer/bindings/Groth16ICS07Tendermint"
	contractICS26Router "relayer/bindings/ICS26Router"
	contractMisbehaviour "relayer/bindings/Misbehaviour"
	contractUpdateClient "relayer/bindings/UpdateClient"
	"relayer/keys"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

type cosmosToEthConfig struct {
	TmRpcUrl           string `json:"tm_rpc_url"`
	ICS26Address       string `json:"ics26_address"`
	ICS26ClientID      string `json:"ics26_client_id"`
	EthRpcUrl          string `json:"eth_rpc_url"`
	ICS07Client        string `json:"ics07_client"`
	WrapperVerifier    string `json:"wrapper_verifier"`
	Membership         string `json:"membership"`
	Misbehaviour       string `json:"misbehaviour"`
	UpdateClient       string `json:"update_client"`
	CosmosWasmClientID string `json:"cosmos_wasm_client_id"`
}

type configModule struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

type rootConfig struct {
	Modules []configModule `json:"modules"`
}

type txData struct {
	TxHash    common.Hash
	ClientID  string
	UpdateMsg contractUpdateClient.IUpdateClientMsgsMsgUpdateClient
}

type preparedHeader struct {
	Data   txData
	Output contractUpdateClient.IUpdateClientMsgsUpdateClientOutput
}

type msgSubmitMisbehaviour struct {
	ClientState            contractMisbehaviour.IICS07TendermintMsgsClientState
	Misbehaviour           contractMisbehaviour.IMisbehaviourMsgsMisbehaviour
	TrustedConsensusState1 contractMisbehaviour.IICS07TendermintMsgsConsensusState
	TrustedConsensusState2 contractMisbehaviour.IICS07TendermintMsgsConsensusState
	Time                   *big.Int
}

type output struct {
	ClientID           string   `json:"client_id"`
	Router             string   `json:"router"`
	ICS07Client        string   `json:"ics07_client"`
	Header1Tx          string   `json:"header1_tx"`
	Header2Tx          string   `json:"header2_tx"`
	Header1Height      uint64   `json:"header1_height"`
	Header2Height      uint64   `json:"header2_height"`
	Header1Trusted     uint64   `json:"header1_trusted_height"`
	Header2Trusted     uint64   `json:"header2_trusted_height"`
	Header1AppHash     string   `json:"header1_app_hash"`
	Header2AppHash     string   `json:"header2_app_hash"`
	SameHeight         bool     `json:"same_height"`
	DifferentAppHash   bool     `json:"different_app_hash"`
	SubmissionTimeNs   string   `json:"submission_time_ns"`
	MisbehaviourMsg    string   `json:"misbehaviour_msg"`
	SubmitCalldata     string   `json:"submit_calldata"`
	SimulateOK         bool     `json:"simulate_ok"`
	SimulateError      string   `json:"simulate_error,omitempty"`
	SimulateErrorData  string   `json:"simulate_error_data,omitempty"`
	SimulateDecoded    string   `json:"simulate_decoded_error,omitempty"`
	ClientFrozenBefore bool     `json:"client_frozen_before"`
	ClientFrozenAfter  bool     `json:"client_frozen_after,omitempty"`
	Submitted          bool     `json:"submitted"`
	SubmitTx           string   `json:"submit_tx,omitempty"`
	ReceiptStatus      uint64   `json:"receipt_status,omitempty"`
	Notes              []string `json:"notes,omitempty"`
	CastCallCommand    string   `json:"cast_call_command,omitempty"`
	CastSendCommand    string   `json:"cast_send_command,omitempty"`
}

func main() {
	var (
		configPath = flag.String("config", "config.example.json", "path to relayer JSON config")
		rpcURL     = flag.String("rpc-url", "", "ethereum rpc url override")
		routerAddr = flag.String("router", "", "ICS26 router address override")
		clientID   = flag.String("client-id", "", "ICS26 client id override")
		tx1Hash    = flag.String("tx1", "", "newer updateClient tx hash override")
		tx2Hash    = flag.String("tx2", "", "older updateClient tx hash override")
		lookback   = flag.Uint64("lookback", 5000, "blocks to scan for historical ICS02ClientUpdated logs")
		fromBlock  = flag.Uint64("from-block", 0, "optional explicit start block")
		submit     = flag.Bool("submit", false, "send submitMisbehaviour transaction after simulation")
	)
	flag.Parse()

	_ = godotenv.Load()

	cfg, err := loadCosmosToEthConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	resolvedRPC := chooseFirstNonEmpty(*rpcURL, cfg.EthRpcUrl)
	if resolvedRPC == "" {
		log.Fatal("ethereum rpc url is required")
	}
	resolvedRouter := chooseFirstNonEmpty(*routerAddr, cfg.ICS26Address)
	if !common.IsHexAddress(resolvedRouter) {
		log.Fatalf("invalid router address: %s", resolvedRouter)
	}
	resolvedClientID := chooseFirstNonEmpty(*clientID, cfg.ICS26ClientID)
	if resolvedClientID == "" {
		log.Fatal("ics26 client id is required")
	}
	if !common.IsHexAddress(cfg.ICS07Client) {
		log.Fatalf("invalid ics07_client address in config: %s", cfg.ICS07Client)
	}
	if !common.IsHexAddress(cfg.UpdateClient) {
		log.Fatalf("invalid update_client address in config: %s", cfg.UpdateClient)
	}

	ctx := context.Background()
	ethClient, err := ethclient.DialContext(ctx, resolvedRPC)
	if err != nil {
		log.Fatalf("dial ethereum rpc: %v", err)
	}
	defer ethClient.Close()

	routerABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil || routerABI == nil {
		log.Fatalf("load router ABI: %v", err)
	}
	updateABI, err := contractUpdateClient.ContractUpdateClientMetaData.GetAbi()
	if err != nil || updateABI == nil {
		log.Fatalf("load update ABI: %v", err)
	}
	misbehaviourABI, err := contractMisbehaviour.ContractMisbehaviourMetaData.GetAbi()
	if err != nil || misbehaviourABI == nil {
		log.Fatalf("load misbehaviour ABI: %v", err)
	}

	header1, header2, err := loadHeaders(
		ctx,
		ethClient,
		*routerABI,
		*updateABI,
		cfg,
		resolvedRouter,
		resolvedClientID,
		*lookback,
		*fromBlock,
		*tx1Hash,
		*tx2Hash,
	)
	if err != nil {
		log.Fatal(err)
	}

	updateClientContract, err := contractUpdateClient.NewContractUpdateClient(common.HexToAddress(cfg.UpdateClient), ethClient)
	if err != nil {
		log.Fatalf("bind update_client contract: %v", err)
	}
	header1.Output, err = updateClientContract.UpdateClient(&bind.CallOpts{Context: ctx}, header1.Data.UpdateMsg)
	if err != nil {
		log.Fatalf("simulate updateClient for header1: %v", err)
	}
	header2.Output, err = updateClientContract.UpdateClient(&bind.CallOpts{Context: ctx}, header2.Data.UpdateMsg)
	if err != nil {
		log.Fatalf("simulate updateClient for header2: %v", err)
	}

	submissionTime := deriveSubmissionTime(header1, header2)
	misbehaviourMsg, notes, err := buildMisbehaviourMsg(misbehaviourABI, header1, header2, submissionTime)
	if err != nil {
		log.Fatal(err)
	}
	submitCalldata, err := routerABI.Pack("submitMisbehaviour", resolvedClientID, misbehaviourMsg)
	if err != nil {
		log.Fatalf("pack submitMisbehaviour calldata: %v", err)
	}

	ics07Client, err := contractGroth16ICS07Tendermint.NewContractGroth16ICS07Tendermint(common.HexToAddress(cfg.ICS07Client), ethClient)
	if err != nil {
		log.Fatalf("bind ics07 client contract: %v", err)
	}
	ics07ABI, err := contractGroth16ICS07Tendermint.ContractGroth16ICS07TendermintMetaData.GetAbi()
	if err != nil || ics07ABI == nil {
		log.Fatalf("load ICS07 ABI: %v", err)
	}
	clientStateBefore, err := ics07Client.ClientState(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Fatalf("read client state before submit: %v", err)
	}

	callErr, callErrData, callDecodedErr := simulateCall(ctx, resolvedRPC, common.HexToAddress(resolvedRouter), submitCalldata, []*abi.ABI{routerABI, ics07ABI, misbehaviourABI})

	out := output{
		ClientID:           resolvedClientID,
		Router:             resolvedRouter,
		ICS07Client:        cfg.ICS07Client,
		Header1Tx:          header1.Data.TxHash.Hex(),
		Header2Tx:          header2.Data.TxHash.Hex(),
		Header1Height:      header1.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.Height,
		Header2Height:      header2.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.Height,
		Header1Trusted:     header1.Data.UpdateMsg.ProposedHeader.TrustedHeight.RevisionHeight,
		Header2Trusted:     header2.Data.UpdateMsg.ProposedHeader.TrustedHeight.RevisionHeight,
		Header1AppHash:     hashHex(header1.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.AppHash),
		Header2AppHash:     hashHex(header2.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.AppHash),
		SameHeight:         header1.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.Height == header2.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.Height,
		DifferentAppHash:   header1.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.AppHash != header2.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.AppHash,
		SubmissionTimeNs:   submissionTime.String(),
		MisbehaviourMsg:    "0x" + hex.EncodeToString(misbehaviourMsg),
		SubmitCalldata:     "0x" + hex.EncodeToString(submitCalldata),
		SimulateOK:         callErr == nil,
		ClientFrozenBefore: clientStateBefore.IsFrozen,
		Submitted:          false,
		Notes:              notes,
		CastCallCommand: fmt.Sprintf(
			"cast call %s 'submitMisbehaviour(string,bytes)' %q %s --rpc-url %s",
			resolvedRouter,
			resolvedClientID,
			"0x"+hex.EncodeToString(misbehaviourMsg),
			resolvedRPC,
		),
		CastSendCommand: fmt.Sprintf(
			"cast send %s 'submitMisbehaviour(string,bytes)' %q %s --rpc-url %s --private-key \"$ETH_PRIVATE_KEY\"",
			resolvedRouter,
			resolvedClientID,
			"0x"+hex.EncodeToString(misbehaviourMsg),
			resolvedRPC,
		),
	}
	if callErr != nil {
		out.SimulateError = callErr.Error()
		out.SimulateErrorData = callErrData
		out.SimulateDecoded = callDecodedErr
	}

	if *submit {
		tx, receipt, frozenAfter, submitErr := submitMisbehaviour(ctx, ethClient, common.HexToAddress(resolvedRouter), common.HexToAddress(cfg.ICS07Client), resolvedClientID, misbehaviourMsg)
		if submitErr != nil {
			log.Fatalf("submit misbehaviour: %v", submitErr)
		}
		out.Submitted = true
		out.SubmitTx = tx.Hex()
		out.ReceiptStatus = receipt.Status
		out.ClientFrozenAfter = frozenAfter
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		log.Fatalf("encode output: %v", err)
	}
}

func simulateCall(ctx context.Context, rpcURL string, to common.Address, data []byte, abis []*abi.ABI) (error, string, string) {
	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		return err, "", ""
	}
	defer rpcClient.Close()

	var result hexutil.Bytes
	callArg := map[string]any{
		"to":   to.Hex(),
		"gas":  hexutil.Uint64(30_000_000),
		"data": hexutil.Encode(data),
	}
	err = rpcClient.CallContext(ctx, &result, "eth_call", callArg, "latest")
	if err == nil {
		return nil, "", ""
	}
	if dataErr, ok := err.(rpc.DataError); ok {
		raw := dataErr.ErrorData()
		if s, ok := raw.(string); ok {
			return err, s, decodeErrorString(s, abis)
		}
	}
	return err, "", ""
}

func decodeErrorString(raw string, abis []*abi.ABI) string {
	if !strings.HasPrefix(raw, "0x") || len(raw) < 10 {
		return ""
	}
	data := common.FromHex(raw)
	if len(data) < 4 {
		return ""
	}
	selector := data[:4]
	for _, parsedABI := range abis {
		if parsedABI == nil {
			continue
		}
		for name, abiError := range parsedABI.Errors {
			if bytes.Equal(abiError.ID.Bytes()[:4], selector) {
				values, unpackErr := abiError.Unpack(data[4:])
				if unpackErr != nil {
					return name
				}
				return fmt.Sprintf("%s%v", name, values)
			}
		}
	}
	return ""
}

func loadHeaders(
	ctx context.Context,
	client *ethclient.Client,
	routerABI abi.ABI,
	updateABI abi.ABI,
	cfg cosmosToEthConfig,
	routerAddress string,
	clientID string,
	lookback uint64,
	fromBlockHint uint64,
	tx1Raw string,
	tx2Raw string,
) (preparedHeader, preparedHeader, error) {
	if tx1Raw != "" || tx2Raw != "" {
		if tx1Raw == "" || tx2Raw == "" {
			return preparedHeader{}, preparedHeader{}, errors.New("--tx1 and --tx2 must be provided together")
		}
		tx1, err := parseTxHash(tx1Raw)
		if err != nil {
			return preparedHeader{}, preparedHeader{}, fmt.Errorf("invalid --tx1: %w", err)
		}
		tx2, err := parseTxHash(tx2Raw)
		if err != nil {
			return preparedHeader{}, preparedHeader{}, fmt.Errorf("invalid --tx2: %w", err)
		}
		header1, err := decodeUpdateTx(ctx, client, &routerABI, &updateABI, tx1)
		if err != nil {
			return preparedHeader{}, preparedHeader{}, fmt.Errorf("decode tx1: %w", err)
		}
		header2, err := decodeUpdateTx(ctx, client, &routerABI, &updateABI, tx2)
		if err != nil {
			return preparedHeader{}, preparedHeader{}, fmt.Errorf("decode tx2: %w", err)
		}
		if header1.ClientID != clientID || header2.ClientID != clientID {
			return preparedHeader{}, preparedHeader{}, fmt.Errorf("tx client ids must both be %q, got %q and %q", clientID, header1.ClientID, header2.ClientID)
		}
		return sortHeaders(header1, header2)
	}

	latest, err := client.BlockNumber(ctx)
	if err != nil {
		return preparedHeader{}, preparedHeader{}, fmt.Errorf("get latest block: %w", err)
	}
	from := uint64(0)
	if fromBlockHint > 0 {
		from = fromBlockHint
	} else if latest > lookback {
		from = latest - lookback
	}

	logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from),
		ToBlock:   nil,
		Addresses: []common.Address{common.HexToAddress(routerAddress)},
		Topics:    [][]common.Hash{{routerABI.Events["ICS02ClientUpdated"].ID}},
	})
	if err != nil {
		return preparedHeader{}, preparedHeader{}, fmt.Errorf("filter ICS02ClientUpdated logs: %w", err)
	}

	found := make([]txData, 0, 2)
	seen := make(map[common.Hash]struct{})
	for i := len(logs) - 1; i >= 0; i-- {
		txHash := logs[i].TxHash
		if _, ok := seen[txHash]; ok {
			continue
		}
		seen[txHash] = struct{}{}

		decoded, err := decodeUpdateTx(ctx, client, &routerABI, &updateABI, txHash)
		if err != nil {
			continue
		}
		if decoded.ClientID != clientID {
			continue
		}
		found = append(found, decoded)
		if len(found) == 2 {
			break
		}
	}
	if len(found) < 2 {
		return preparedHeader{}, preparedHeader{}, fmt.Errorf("found only %d updateClient txs for client %q in lookback window", len(found), clientID)
	}
	return sortHeaders(found[0], found[1])
}

func sortHeaders(a, b txData) (preparedHeader, preparedHeader, error) {
	ah := a.UpdateMsg.ProposedHeader.SignedHeader.Header.Height
	bh := b.UpdateMsg.ProposedHeader.SignedHeader.Header.Height
	if ah > bh || (ah == bh && a.TxHash.Hex() >= b.TxHash.Hex()) {
		return preparedHeader{Data: a}, preparedHeader{Data: b}, nil
	}
	return preparedHeader{Data: b}, preparedHeader{Data: a}, nil
}

func buildMisbehaviourMsg(
	misbehaviourABI *abi.ABI,
	header1 preparedHeader,
	header2 preparedHeader,
	submissionTime *big.Int,
) ([]byte, []string, error) {
	clientState := toMisClientState(header1.Data.UpdateMsg.ClientState)
	chainID, revision := parseChainID(clientState.ChainId)
	misbehaviour := contractMisbehaviour.IMisbehaviourMsgsMisbehaviour{
		ClientId: contractMisbehaviour.IICS07TendermintMsgsChainId{
			Id:             chainID,
			RevisionNumber: revision,
		},
		Header1: toMisHeader(header1.Data.UpdateMsg.ProposedHeader),
		Header2: toMisHeader(header2.Data.UpdateMsg.ProposedHeader),
	}
	msg := msgSubmitMisbehaviour{
		ClientState:            clientState,
		Misbehaviour:           misbehaviour,
		TrustedConsensusState1: toMisConsensusState(header1.Output.TrustedConsensusState),
		TrustedConsensusState2: toMisConsensusState(header2.Output.TrustedConsensusState),
		Time:                   submissionTime,
	}
	payload, err := packMsgSubmitMisbehaviour(misbehaviourABI, msg)
	if err != nil {
		return nil, nil, fmt.Errorf("pack MsgSubmitMisbehaviour: %w", err)
	}

	notes := make([]string, 0, 3)
	if header1.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.Height != header2.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.Height {
		notes = append(notes, "headers are from different heights; this exercises the current on-chain implementation, not true same-height equivocation")
	}
	if header1.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.AppHash == header2.Data.UpdateMsg.ProposedHeader.SignedHeader.Header.AppHash {
		notes = append(notes, "headers have the same appHash; if submission still succeeds, freeze is caused by implementation behavior rather than conflicting state roots")
	}
	notes = append(notes, "trustedConsensusState1/2 are populated from each header's trusted consensus state at header.trustedHeight")
	return payload, notes, nil
}

func deriveSubmissionTime(header1 preparedHeader, header2 preparedHeader) *big.Int {
	candidates := []*big.Int{
		header1.Output.TrustedConsensusState.Timestamp,
		header2.Output.TrustedConsensusState.Timestamp,
	}
	maxTs := big.NewInt(0)
	for _, ts := range candidates {
		if ts != nil && ts.Cmp(maxTs) > 0 {
			maxTs = new(big.Int).Set(ts)
		}
	}
	if maxTs.Sign() == 0 {
		return big.NewInt(time.Now().UnixNano())
	}
	return new(big.Int).Add(maxTs, big.NewInt(1))
}

func packMsgSubmitMisbehaviour(misbehaviourABI *abi.ABI, msg msgSubmitMisbehaviour) ([]byte, error) {
	method := misbehaviourABI.Methods["misbehaviour"]
	components := make([]abi.ArgumentMarshaling, 0, len(method.Inputs))
	for _, input := range method.Inputs {
		components = append(components, typeToArgumentMarshaling(input.Name, input.Type))
	}
	msgType, err := abi.NewType("tuple", "IMisbehaviourMsgs.MsgSubmitMisbehaviour", components)
	if err != nil {
		return nil, err
	}
	args := abi.Arguments{{Name: "msg_", Type: msgType}}
	return args.Pack(msg)
}

func typeToArgumentMarshaling(name string, t abi.Type) abi.ArgumentMarshaling {
	m := abi.ArgumentMarshaling{Name: name}
	switch t.T {
	case abi.TupleTy:
		m.Type = "tuple"
		m.Components = make([]abi.ArgumentMarshaling, 0, len(t.TupleElems))
		for i, elem := range t.TupleElems {
			childName := ""
			if i < len(t.TupleRawNames) {
				childName = t.TupleRawNames[i]
			}
			m.Components = append(m.Components, typeToArgumentMarshaling(childName, *elem))
		}
	case abi.ArrayTy:
		if t.Elem != nil && t.Elem.T == abi.TupleTy {
			m.Type = fmt.Sprintf("tuple[%d]", t.Size)
			child := typeToArgumentMarshaling("", *t.Elem)
			m.Components = child.Components
		} else {
			m.Type = t.String()
		}
	case abi.SliceTy:
		if t.Elem != nil && t.Elem.T == abi.TupleTy {
			m.Type = "tuple[]"
			child := typeToArgumentMarshaling("", *t.Elem)
			m.Components = child.Components
		} else {
			m.Type = t.String()
		}
	default:
		m.Type = t.String()
	}
	return m
}

func submitMisbehaviour(
	ctx context.Context,
	ethClient *ethclient.Client,
	routerAddress common.Address,
	clientAddress common.Address,
	clientID string,
	misbehaviourMsg []byte,
) (common.Hash, *types.Receipt, bool, error) {
	const gasCap uint64 = 16_700_000

	privKey := os.Getenv("ETH_PRIVATE_KEY")
	if privKey == "" {
		return common.Hash{}, nil, false, errors.New("ETH_PRIVATE_KEY environment variable is required")
	}
	privateKey, err := keys.RestoreKey(privKey)
	if err != nil {
		return common.Hash{}, nil, false, fmt.Errorf("restore private key: %w", err)
	}
	publicKey, err := keys.PublicKey(privateKey)
	if err != nil {
		return common.Hash{}, nil, false, fmt.Errorf("derive public key: %w", err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKey)
	nonce, err := ethClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return common.Hash{}, nil, false, fmt.Errorf("get nonce: %w", err)
	}
	gasPrice, err := ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, nil, false, fmt.Errorf("suggest gas price: %w", err)
	}
	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return common.Hash{}, nil, false, fmt.Errorf("read chain id: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return common.Hash{}, nil, false, fmt.Errorf("create transactor: %w", err)
	}

	routerABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil || routerABI == nil {
		return common.Hash{}, nil, false, fmt.Errorf("load router ABI: %w", err)
	}
	submitCalldata, err := routerABI.Pack("submitMisbehaviour", clientID, misbehaviourMsg)
	if err != nil {
		return common.Hash{}, nil, false, fmt.Errorf("pack submitMisbehaviour calldata: %w", err)
	}
	estimatedGas, err := ethClient.EstimateGas(ctx, ethereum.CallMsg{
		From:     fromAddress,
		To:       ptr(routerAddress),
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     submitCalldata,
	})
	if err != nil {
		estimatedGas = gasCap
	} else {
		estimatedGas = estimatedGas + estimatedGas/5 + 50_000
		if estimatedGas > gasCap {
			estimatedGas = gasCap
		}
	}

	auth.Nonce = new(big.Int).SetUint64(nonce)
	auth.Value = big.NewInt(0)
	auth.GasLimit = estimatedGas
	auth.GasPrice = gasPrice

	router, err := contractICS26Router.NewContractICS26Router(routerAddress, ethClient)
	if err != nil {
		return common.Hash{}, nil, false, fmt.Errorf("bind router: %w", err)
	}
	tx, err := router.SubmitMisbehaviour(auth, clientID, misbehaviourMsg)
	if err != nil {
		return common.Hash{}, nil, false, fmt.Errorf("send tx: %w", err)
	}
	receipt, err := bind.WaitMined(ctx, ethClient, tx)
	if err != nil {
		return tx.Hash(), nil, false, fmt.Errorf("wait mined: %w", err)
	}

	ics07Client, err := contractGroth16ICS07Tendermint.NewContractGroth16ICS07Tendermint(clientAddress, ethClient)
	if err != nil {
		return tx.Hash(), receipt, false, fmt.Errorf("bind ics07 client: %w", err)
	}
	stateAfter, err := ics07Client.ClientState(&bind.CallOpts{Context: ctx})
	if err != nil {
		return tx.Hash(), receipt, false, fmt.Errorf("read client state after submit: %w", err)
	}
	return tx.Hash(), receipt, stateAfter.IsFrozen, nil
}

func decodeUpdateTx(
	ctx context.Context,
	client *ethclient.Client,
	routerABI *abi.ABI,
	updateABI *abi.ABI,
	txHash common.Hash,
) (txData, error) {
	tx, _, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return txData{}, fmt.Errorf("fetch tx %s: %w", txHash.Hex(), err)
	}
	data := tx.Data()
	if len(data) < 4 {
		return txData{}, fmt.Errorf("tx %s has short calldata", txHash.Hex())
	}

	updateMethod := routerABI.Methods["updateClient"]
	if !bytes.Equal(data[:4], updateMethod.ID) {
		return txData{}, fmt.Errorf("tx %s is not ICS26.updateClient", txHash.Hex())
	}

	args, err := updateMethod.Inputs.Unpack(data[4:])
	if err != nil {
		return txData{}, fmt.Errorf("unpack router updateClient args: %w", err)
	}
	clientID, ok := args[0].(string)
	if !ok {
		return txData{}, errors.New("clientId arg is not string")
	}
	updateMsgBytes, ok := args[1].([]byte)
	if !ok {
		return txData{}, errors.New("updateMsg arg is not bytes")
	}

	msgMethod := updateABI.Methods["updateClient"]
	msgArgs, err := msgMethod.Inputs.Unpack(updateMsgBytes)
	if err != nil {
		return txData{}, fmt.Errorf("unpack inner updateClient message: %w", err)
	}
	msg := *abi.ConvertType(msgArgs[0], new(contractUpdateClient.IUpdateClientMsgsMsgUpdateClient)).(*contractUpdateClient.IUpdateClientMsgsMsgUpdateClient)
	return txData{
		TxHash:    txHash,
		ClientID:  clientID,
		UpdateMsg: msg,
	}, nil
}

func loadCosmosToEthConfig(path string) (cosmosToEthConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return cosmosToEthConfig{}, fmt.Errorf("read config %s: %w", path, err)
	}
	var root rootConfig
	if err := json.Unmarshal(raw, &root); err != nil {
		return cosmosToEthConfig{}, fmt.Errorf("parse config %s: %w", path, err)
	}
	for _, m := range root.Modules {
		if m.Name != "cosmos_to_eth" {
			continue
		}
		var cfg cosmosToEthConfig
		if err := json.Unmarshal(m.Config, &cfg); err != nil {
			return cosmosToEthConfig{}, fmt.Errorf("parse cosmos_to_eth config: %w", err)
		}
		return cfg, nil
	}
	return cosmosToEthConfig{}, errors.New("module cosmos_to_eth not found in config")
}

func chooseFirstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func parseTxHash(s string) (common.Hash, error) {
	if !common.IsHexHash(s) {
		return common.Hash{}, errors.New("not a valid 0x-prefixed 32-byte hash")
	}
	return common.HexToHash(s), nil
}

func parseChainID(id string) (string, uint64) {
	lastDash := strings.LastIndexByte(id, '-')
	if lastDash == -1 || lastDash == len(id)-1 {
		return id, 0
	}
	revision, err := strconv.ParseUint(id[lastDash+1:], 10, 64)
	if err != nil {
		return id, 0
	}
	return id, revision
}

func hashHex(v [32]byte) string {
	return common.BytesToHash(v[:]).Hex()
}

func ptr[T any](v T) *T {
	return &v
}

func toMisClientState(v contractUpdateClient.IICS07TendermintMsgsClientState) contractMisbehaviour.IICS07TendermintMsgsClientState {
	return contractMisbehaviour.IICS07TendermintMsgsClientState{
		ChainId:         v.ChainId,
		TrustLevel:      contractMisbehaviour.IICS07TendermintMsgsTrustThreshold{Numerator: v.TrustLevel.Numerator, Denominator: v.TrustLevel.Denominator},
		LatestHeight:    contractMisbehaviour.IICS02ClientMsgsHeight{RevisionNumber: v.LatestHeight.RevisionNumber, RevisionHeight: v.LatestHeight.RevisionHeight},
		TrustingPeriod:  v.TrustingPeriod,
		UnbondingPeriod: v.UnbondingPeriod,
		IsFrozen:        v.IsFrozen,
		ZkAlgorithm:     v.ZkAlgorithm,
	}
}

func toMisConsensusState(v contractUpdateClient.IICS07TendermintMsgsConsensusState) contractMisbehaviour.IICS07TendermintMsgsConsensusState {
	return contractMisbehaviour.IICS07TendermintMsgsConsensusState{
		Timestamp:          v.Timestamp,
		Root:               v.Root,
		NextValidatorsHash: v.NextValidatorsHash,
	}
}

func toMisHeader(v contractUpdateClient.IICS07TendermintMsgsHeader) contractMisbehaviour.IICS07TendermintMsgsHeader {
	return contractMisbehaviour.IICS07TendermintMsgsHeader{
		SignedHeader:            toMisSignedHeader(v.SignedHeader),
		ValidatorSet:            toMisValidatorSet(v.ValidatorSet),
		TrustedHeight:           contractMisbehaviour.IICS02ClientMsgsHeight{RevisionNumber: v.TrustedHeight.RevisionNumber, RevisionHeight: v.TrustedHeight.RevisionHeight},
		TrustedNextValidatorSet: toMisValidatorSet(v.TrustedNextValidatorSet),
	}
}

func toMisSignedHeader(v contractUpdateClient.IICS07TendermintMsgsSignedHeader) contractMisbehaviour.IICS07TendermintMsgsSignedHeader {
	return contractMisbehaviour.IICS07TendermintMsgsSignedHeader{
		Header: toMisBlockHeader(v.Header),
		Commit: toMisBlockCommit(v.Commit),
	}
}

func toMisBlockHeader(v contractUpdateClient.IICS07TendermintMsgsBlockHeader) contractMisbehaviour.IICS07TendermintMsgsBlockHeader {
	return contractMisbehaviour.IICS07TendermintMsgsBlockHeader{
		Version:            contractMisbehaviour.IICS07TendermintMsgsVersion{BlockVersion: v.Version.BlockVersion, AppVersion: v.Version.AppVersion},
		ChainId:            v.ChainId,
		Height:             v.Height,
		Time:               v.Time,
		HasLastBlockId:     v.HasLastBlockId,
		LastBlockId:        toMisBlockID(v.LastBlockId),
		HasLastCommitHash:  v.HasLastCommitHash,
		LastCommitHash:     v.LastCommitHash,
		HasDataHash:        v.HasDataHash,
		DataHash:           v.DataHash,
		ValidatorsHash:     v.ValidatorsHash,
		NextValidatorsHash: v.NextValidatorsHash,
		ConsensusHash:      v.ConsensusHash,
		AppHash:            v.AppHash,
		HasLastResultsHash: v.HasLastResultsHash,
		LastResultsHash:    v.LastResultsHash,
		HasEvidenceHash:    v.HasEvidenceHash,
		EvidenceHash:       v.EvidenceHash,
		ProposerAddress:    v.ProposerAddress,
	}
}

func toMisBlockCommit(v contractUpdateClient.IICS07TendermintMsgsBlockCommit) contractMisbehaviour.IICS07TendermintMsgsBlockCommit {
	commitSigs := make([]contractMisbehaviour.IICS07TendermintMsgsCommitSig, 0, len(v.CommitSigs))
	for _, sig := range v.CommitSigs {
		commitSigs = append(commitSigs, contractMisbehaviour.IICS07TendermintMsgsCommitSig{
			Flag: sig.Flag,
			Data: contractMisbehaviour.IICS07TendermintMsgsCommitSigData{
				ValidatorAddress: sig.Data.ValidatorAddress,
				Timestamp:        sig.Data.Timestamp,
				HasSignature:     sig.Data.HasSignature,
				Signature:        sig.Data.Signature,
			},
		})
	}
	return contractMisbehaviour.IICS07TendermintMsgsBlockCommit{
		Height:     v.Height,
		Round:      v.Round,
		BlockId:    toMisBlockID(v.BlockId),
		CommitSigs: commitSigs,
	}
}

func toMisBlockID(v contractUpdateClient.IICS07TendermintMsgsBlockId) contractMisbehaviour.IICS07TendermintMsgsBlockId {
	return contractMisbehaviour.IICS07TendermintMsgsBlockId{
		HashData: v.HashData,
		PartSetHeader: contractMisbehaviour.IICS07TendermintMsgsPartSetHeader{
			Total:    v.PartSetHeader.Total,
			HashData: v.PartSetHeader.HashData,
		},
	}
}

func toMisValidatorSet(v contractUpdateClient.IICS07TendermintMsgsValidatorSet) contractMisbehaviour.IICS07TendermintMsgsValidatorSet {
	validators := make([]contractMisbehaviour.IICS07TendermintMsgsValidatorInfo, 0, len(v.Validators))
	for _, validator := range v.Validators {
		validators = append(validators, contractMisbehaviour.IICS07TendermintMsgsValidatorInfo{
			ValAddress:       validator.ValAddress,
			PubKey:           validator.PubKey,
			VotingPower:      validator.VotingPower,
			ProposerPriority: validator.ProposerPriority,
		})
	}
	return contractMisbehaviour.IICS07TendermintMsgsValidatorSet{
		Validators:       validators,
		HasProposer:      v.HasProposer,
		Proposer:         toMisValidatorInfo(v.Proposer),
		TotalVotingPower: v.TotalVotingPower,
	}
}

func toMisValidatorInfo(v contractUpdateClient.IICS07TendermintMsgsValidatorInfo) contractMisbehaviour.IICS07TendermintMsgsValidatorInfo {
	return contractMisbehaviour.IICS07TendermintMsgsValidatorInfo{
		ValAddress:       v.ValAddress,
		PubKey:           v.PubKey,
		VotingPower:      v.VotingPower,
		ProposerPriority: v.ProposerPriority,
	}
}
