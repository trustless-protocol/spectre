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
	"strings"

	contractICS26Router "relayer/bindings/ICS26Router"
	contractUpdateClient "relayer/bindings/UpdateClient"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type output struct {
	Mode             string `json:"mode"`
	BaseTx           string `json:"base_tx"`
	DonorTx          string `json:"donor_tx,omitempty"`
	ClientID         string `json:"client_id"`
	Mutation         string `json:"mutation"`
	BaseAppHash      string `json:"base_app_hash,omitempty"`
	DonorAppHash     string `json:"donor_app_hash,omitempty"`
	TamperedCalldata string `json:"tampered_calldata"`
}

type txData struct {
	TxHash    common.Hash
	ClientID  string
	UpdateMsg contractUpdateClient.IUpdateClientMsgsMsgUpdateClient
}

func main() {
	var (
		mode          = flag.String("mode", "4b", "attack mode: 4b or 4c")
		rpcURL        = flag.String("rpc-url", "", "ethereum rpc url")
		router        = flag.String("router", "", "ics26 router address")
		lookback      = flag.Uint64("lookback", 5000, "blocks to scan for ICS02ClientUpdated logs")
		baseTxHash    = flag.String("base-tx", "", "base updateClient tx hash")
		donorTxHash   = flag.String("donor-tx", "", "donor updateClient tx hash (required for 4b unless auto-detected)")
		baseBlockHint = flag.Uint64("from-block", 0, "optional explicit start block")
	)
	flag.Parse()

	if *rpcURL == "" {
		log.Fatal("--rpc-url is required")
	}
	if !common.IsHexAddress(*router) {
		log.Fatal("--router must be a valid hex address")
	}
	modeVal := strings.ToLower(*mode)
	if modeVal != "4b" && modeVal != "4c" {
		log.Fatal("--mode must be 4b or 4c")
	}

	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, *rpcURL)
	if err != nil {
		log.Fatalf("dial rpc: %v", err)
	}
	defer client.Close()

	routerAddr := common.HexToAddress(*router)
	routerABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil || routerABI == nil {
		log.Fatalf("load ICS26 router ABI: %v", err)
	}
	updateABI, err := contractUpdateClient.ContractUpdateClientMetaData.GetAbi()
	if err != nil || updateABI == nil {
		log.Fatalf("load UpdateClient ABI: %v", err)
	}

	baseTx, donorTx, err := resolveTxHashes(ctx, client, routerABI.Events["ICS02ClientUpdated"].ID, routerAddr, *lookback, *baseBlockHint, *baseTxHash, *donorTxHash, modeVal)
	if err != nil {
		log.Fatal(err)
	}

	baseData, err := decodeUpdateTx(ctx, client, routerABI, updateABI, baseTx)
	if err != nil {
		log.Fatalf("decode base tx: %v", err)
	}

	out := output{
		Mode:     modeVal,
		BaseTx:   baseData.TxHash.Hex(),
		ClientID: baseData.ClientID,
	}

	var tamperedMsg contractUpdateClient.IUpdateClientMsgsMsgUpdateClient

	switch modeVal {
	case "4b":
		donorData, err := decodeUpdateTx(ctx, client, routerABI, updateABI, donorTx)
		if err != nil {
			log.Fatalf("decode donor tx: %v", err)
		}
		if donorData.ClientID != baseData.ClientID {
			log.Fatalf("client id mismatch between donor(%s) and base(%s)", donorData.ClientID, baseData.ClientID)
		}

		tamperedMsg = donorData.UpdateMsg
		oldHash := tamperedMsg.ProposedHeader.SignedHeader.Header.AppHash
		newHash := baseData.UpdateMsg.ProposedHeader.SignedHeader.Header.AppHash

		if oldHash == newHash {
			log.Fatal("donor and base appHash are identical; choose txs from different heights")
		}

		tamperedMsg.ProposedHeader.SignedHeader.Header.AppHash = newHash
		out.DonorTx = donorData.TxHash.Hex()
		out.DonorAppHash = common.BytesToHash(oldHash[:]).Hex()
		out.BaseAppHash = common.BytesToHash(newHash[:]).Hex()
		out.Mutation = "4b: donor proof with base appHash (witness mismatch)"

	case "4c":
		tamperedMsg = baseData.UpdateMsg
		i, j, ok := findDistinctPubkeyPair(tamperedMsg.SignerPubkeys)
		if !ok {
			log.Fatal("unable to find two distinct signer pubkeys in base tx; try another tx via --base-tx")
		}
		tamperedMsg.SignerPubkeys[i], tamperedMsg.SignerPubkeys[j] = tamperedMsg.SignerPubkeys[j], tamperedMsg.SignerPubkeys[i]
		out.Mutation = fmt.Sprintf("4c: swapped signerPubkeys[%d] and signerPubkeys[%d]", i, j)
	}

	tamperedBytes, err := updateABI.Methods["updateClient"].Inputs.Pack(tamperedMsg)
	if err != nil {
		log.Fatalf("pack tampered updateMsg: %v", err)
	}

	callData, err := routerABI.Pack("updateClient", baseData.ClientID, tamperedBytes)
	if err != nil {
		log.Fatalf("pack ICS26 updateClient calldata: %v", err)
	}

	out.TamperedCalldata = "0x" + hex.EncodeToString(callData)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		log.Fatalf("encode output: %v", err)
	}
}

func resolveTxHashes(
	ctx context.Context,
	client *ethclient.Client,
	topic common.Hash,
	router common.Address,
	lookback uint64,
	fromBlockHint uint64,
	baseTxRaw string,
	donorTxRaw string,
	mode string,
) (common.Hash, common.Hash, error) {
	var baseTx, donorTx common.Hash
	var err error

	if baseTxRaw != "" {
		baseTx, err = parseTxHash(baseTxRaw)
		if err != nil {
			return common.Hash{}, common.Hash{}, fmt.Errorf("invalid --base-tx: %w", err)
		}
	}
	if donorTxRaw != "" {
		donorTx, err = parseTxHash(donorTxRaw)
		if err != nil {
			return common.Hash{}, common.Hash{}, fmt.Errorf("invalid --donor-tx: %w", err)
		}
	}

	if baseTxRaw != "" && (mode == "4c" || donorTxRaw != "") {
		return baseTx, donorTx, nil
	}

	latest, err := client.BlockNumber(ctx)
	if err != nil {
		return common.Hash{}, common.Hash{}, fmt.Errorf("get latest block: %w", err)
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
		Addresses: []common.Address{router},
		Topics:    [][]common.Hash{{topic}},
	})
	if err != nil {
		return common.Hash{}, common.Hash{}, fmt.Errorf("filter ICS02ClientUpdated logs: %w", err)
	}
	if len(logs) == 0 {
		return common.Hash{}, common.Hash{}, fmt.Errorf("no ICS02ClientUpdated logs found in [%d, latest]", from)
	}

	latestLog := logs[len(logs)-1]
	prevLog := types.Log{}
	hasPrev := len(logs) >= 2
	if hasPrev {
		prevLog = logs[len(logs)-2]
	}

	if (baseTx == common.Hash{}) {
		baseTx = latestLog.TxHash
	}
	if mode == "4b" && (donorTx == common.Hash{}) {
		if !hasPrev {
			return common.Hash{}, common.Hash{}, errors.New("4b requires at least two ICS02ClientUpdated txs in lookback window; provide --donor-tx or increase --lookback")
		}
		donorTx = prevLog.TxHash
	}

	return baseTx, donorTx, nil
}

func parseTxHash(s string) (common.Hash, error) {
	if !common.IsHexHash(s) {
		return common.Hash{}, errors.New("not a valid 0x-prefixed 32-byte hash")
	}
	return common.HexToHash(s), nil
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
		return txData{}, fmt.Errorf("unpack ICS26 updateClient args: %w", err)
	}
	if len(args) != 2 {
		return txData{}, fmt.Errorf("unexpected ICS26.updateClient arg count: %d", len(args))
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
		return txData{}, fmt.Errorf("unpack inner MsgUpdateClient bytes: %w", err)
	}
	if len(msgArgs) != 1 {
		return txData{}, fmt.Errorf("unexpected MsgUpdateClient arg count: %d", len(msgArgs))
	}

	msg := *abi.ConvertType(msgArgs[0], new(contractUpdateClient.IUpdateClientMsgsMsgUpdateClient)).(*contractUpdateClient.IUpdateClientMsgsMsgUpdateClient)

	return txData{
		TxHash:    txHash,
		ClientID:  clientID,
		UpdateMsg: msg,
	}, nil
}

func findDistinctPubkeyPair(pubkeys [][32]byte) (int, int, bool) {
	for i := 0; i < len(pubkeys); i++ {
		for j := i + 1; j < len(pubkeys); j++ {
			if pubkeys[i] != pubkeys[j] {
				return i, j, true
			}
		}
	}
	return 0, 0, false
}
