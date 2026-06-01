package transaction

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	tendermintContract "relayer/bindings/Groth16ICS07Tendermint"
	contractICS26Router "relayer/bindings/ICS26Router"
	routerContract "relayer/bindings/ICS26Router"
	updateclient "relayer/bindings/UpdateClient"
	relayerclient "relayer/client"
	"relayer/keys"
	services "relayer/services"
	utils "relayer/utils"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdksigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	clienttypesv2 "github.com/cosmos/ibc-go/v10/modules/core/02-client/v2/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	exported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Handler struct {
}

const ethTxReceiptTimeout = 45 * time.Second

func routerManagesProofSubmission(ctx services.Context) bool {
	roleManager := ctx.RoleManagerAddress()
	router := ctx.RouterContract()
	if roleManager == nil || router == nil {
		return false
	}
	if *roleManager == (common.Address{}) || *router == (common.Address{}) {
		return false
	}
	return *roleManager == *router
}

// callTrace mirrors the geth callTracer output. We only need the gasUsed of
// the top-level call's direct children (one entry per multicall inner). Other
// fields are ignored.
type callTrace struct {
	GasUsed string      `json:"gasUsed"`
	Calls   []callTrace `json:"calls"`
}

// logInnerGasFromTrace runs debug_traceTransaction with the callTracer against
// txHash and logs the gasUsed of each multicall inner call. It returns nil on
// success and an error if the trace could not be obtained or parsed — callers
// log the error and continue (this is best-effort instrumentation).
//
// Requires the RPC endpoint to expose the `debug_` namespace. Kurtosis-Geth
// devnets enable it by default; production endpoints typically do not.
//
// UUPS proxies (e.g. ICS26Router) introduce an intermediate proxy→impl
// delegatecall frame between the tx entry and the multicall body. Descend
// single-child wrapper frames until we find the frame whose children count
// matches the expected inner-call count.
func logInnerGasFromTrace(ctx services.Context, txHash common.Hash, labels []string) error {
	rpcClient := ctx.EthClient().Client()
	if rpcClient == nil {
		return fmt.Errorf("nil rpc client")
	}

	var root callTrace
	tracerCfg := map[string]any{"tracer": "callTracer"}
	if err := rpcClient.CallContext(context.Background(), &root, "debug_traceTransaction", txHash, tracerCfg); err != nil {
		return fmt.Errorf("debug_traceTransaction: %w", err)
	}

	frame := findMulticallFrame(&root, len(labels))
	if frame == nil || len(frame.Calls) == 0 {
		return fmt.Errorf("trace returned no inner calls (expected %d)", len(labels))
	}

	for i, child := range frame.Calls {
		gasUsed, perr := hexToUint64(child.GasUsed)
		if perr != nil {
			log.Printf("[bench][eth] inner[%d] gas parse failed: %v", i, perr)
			continue
		}
		label := fmt.Sprintf("inner[%d]", i)
		if i < len(labels) {
			label = fmt.Sprintf("inner[%d] %s", i, labels[i])
		}
		log.Printf("[bench][eth] %s gas=%d (from trace)", label, gasUsed)
	}
	return nil
}

// findMulticallFrame walks down the trace, skipping single-child wrapper
// frames (UUPS proxy → impl delegatecall) until it finds the frame whose
// direct children match the expected inner-call count. If no exact match is
// found, returns the deepest single-chain frame — its children are still
// the best approximation of the inner calls.
func findMulticallFrame(node *callTrace, expectedInners int) *callTrace {
	if node == nil {
		return nil
	}
	cur := node
	for {
		if expectedInners > 0 && len(cur.Calls) == expectedInners {
			return cur
		}
		if len(cur.Calls) == 1 {
			cur = &cur.Calls[0]
			continue
		}
		return cur
	}
}

func hexToUint64(s string) (uint64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty hex string")
	}
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}
	if s == "" {
		return 0, nil
	}
	var v uint64
	if _, err := fmt.Sscanf(s, "%x", &v); err != nil {
		return 0, fmt.Errorf("parse %q: %w", s, err)
	}
	return v, nil
}

func cosmosRouterClientID(ctx services.Context) (string, error) {
	clientID := ctx.CosmosRouterClientID()
	if clientID == "" {
		return "", fmt.Errorf("cosmos router client id is not configured")
	}
	return clientID, nil
}

func ethLightClientIDOnCosmos(ctx services.Context) (string, error) {
	clientID := ctx.EthClientID()
	if clientID == "" {
		return "", fmt.Errorf("cosmos wasm client id is not configured")
	}
	return clientID, nil
}

func (h *Handler) CreateCosmosClientContract(ctx services.Context, clientState, consensusHash []byte) (common.Address, error) {
	cosmosClientID, err := cosmosRouterClientID(ctx)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] %w", err)
	}
	wasmClientID, err := ethLightClientIDOnCosmos(ctx)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] %w", err)
	}
	privKey := os.Getenv("ETH_PRIVATE_KEY")
	if privKey == "" {
		return common.Address{}, fmt.Errorf("ETH_PRIVATE_KEY environment variable is required in .env file")
	}
	privateKey, err := keys.RestoreKey(privKey)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to restore private key: %w", err)
	}

	publicKey, err := keys.PublicKey(privateKey)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] failed to derive public key: %w", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKey)
	nonce, err := ctx.EthClient().PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] failed to get nonce: %w", err)
	}
	gasPrice, err := ctx.EthClient().SuggestGasPrice(context.Background())
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] failed to suggest gas price: %w", err)
	}

	chainIdInt, err := ctx.EthClient().ChainID(context.Background())
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] invalid chain id: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIdInt)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] failed to create auth transactor: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)       // in wei
	auth.GasLimit = uint64(10000000) // in units
	auth.GasPrice = gasPrice

	address, tx, _, err := tendermintContract.DeployContractGroth16ICS07Tendermint(
		auth,
		ctx.EthClient(),
		*ctx.VerifierContract(),
		*ctx.MembershipContract(),
		*ctx.MisbehaviourContract(),
		*ctx.UpdateClientContract(),
		clientState,
		utils.BytesToBytes32(consensusHash),
		*ctx.RoleManagerAddress(),
	)

	if err != nil {
		return common.Address{}, fmt.Errorf("failed to deploy ics07 contract: %w", err)
	}
	log.Printf("[CreateCosmosClient] Deploy tx sent: %s. Waiting for receipt...", tx.Hash().Hex())
	receiptCtx, cancel := context.WithTimeout(context.Background(), ethTxReceiptTimeout)
	defer cancel()
	receipt, err := bind.WaitMined(receiptCtx, ctx.EthClient(), tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed waiting for deploy receipt: %w", err)
	}
	if receipt.Status == 0 {
		return common.Address{}, fmt.Errorf("deploy tx %s reverted (gasUsed=%d)", tx.Hash().Hex(), receipt.GasUsed)
	}
	log.Printf("[CreateCosmosClient] ICS07 deployed at %s (block %d, gasUsed=%d)", address.String(), receipt.BlockNumber.Uint64(), receipt.GasUsed)
	ctx.SetClient(address)

	// In Eureka mode, the router is typically both the admin and proof submitter
	// for the ICS07 client. Direct submission remains available only when the
	// role manager is not the router.

	ics26Router, err := routerContract.NewContractICS26Router(*ctx.RouterContract(), ctx.EthClient())
	if err != nil {
		return common.Address{}, err
	}

	nonce, err = ctx.EthClient().PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] failed to get nonce for AddClient: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err = ics26Router.AddClient(
		auth,
		cosmosClientID,
		routerContract.IICS02ClientMsgsCounterpartyInfo{
			ClientId:     wasmClientID,
			MerklePrefix: [][]byte{[]byte("ibc"), []byte("")},
		},
		*ctx.ClientContract(),
	)

	if err != nil {
		return common.Address{}, fmt.Errorf("failed to add client to ICS26Router: %w", err)
	}
	log.Printf("[CreateCosmosClient] AddClient tx sent: %s. Waiting for receipt...", tx.Hash().Hex())
	receiptCtx, cancel = context.WithTimeout(context.Background(), ethTxReceiptTimeout)
	defer cancel()
	receipt, err = bind.WaitMined(receiptCtx, ctx.EthClient(), tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed waiting for AddClient receipt: %w", err)
	}
	if receipt.Status == 0 {
		log.Printf("[CreateCosmosClient] AddClient reverted (gasUsed=%d) — falling back to MigrateClient to repoint %s to new ICS07 %s",
			receipt.GasUsed, cosmosClientID, address.Hex())

		nonce, err = ctx.EthClient().PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return common.Address{}, fmt.Errorf("[CreateCosmosClient] failed to get nonce for MigrateClient: %w", err)
		}
		auth.Nonce = big.NewInt(int64(nonce))

		mtx, err := ics26Router.MigrateClient(
			auth,
			cosmosClientID,
			routerContract.IICS02ClientMsgsCounterpartyInfo{
				ClientId:     wasmClientID,
				MerklePrefix: [][]byte{[]byte("ibc"), []byte("")},
			},
			*ctx.ClientContract(),
		)
		if err != nil {
			return common.Address{}, fmt.Errorf("[CreateCosmosClient] MigrateClient call failed: %w", err)
		}
		log.Printf("[CreateCosmosClient] MigrateClient tx sent: %s. Waiting for receipt...", mtx.Hash().Hex())
		mctx, mcancel := context.WithTimeout(context.Background(), ethTxReceiptTimeout)
		defer mcancel()
		mreceipt, err := bind.WaitMined(mctx, ctx.EthClient(), mtx)
		if err != nil {
			return common.Address{}, fmt.Errorf("failed waiting for MigrateClient receipt: %w", err)
		}
		if mreceipt.Status == 0 {
			return common.Address{}, fmt.Errorf("MigrateClient tx %s reverted (gasUsed=%d)", mtx.Hash().Hex(), mreceipt.GasUsed)
		}
		log.Printf("[CreateCosmosClient] MigrateClient confirmed (block %d, gasUsed=%d)", mreceipt.BlockNumber.Uint64(), mreceipt.GasUsed)
		return address, nil
	}
	log.Printf("[CreateCosmosClient] AddClient confirmed (block %d, gasUsed=%d)", receipt.BlockNumber.Uint64(), receipt.GasUsed)

	return address, nil
}

func (h *Handler) SendEthTx(ctx services.Context, msg any) error {
	cosmosClientID, err := cosmosRouterClientID(ctx)
	if err != nil {
		return fmt.Errorf("[SendEthTx] %w", err)
	}
	privKey := os.Getenv("ETH_PRIVATE_KEY")
	if privKey == "" {
		return fmt.Errorf("ETH_PRIVATE_KEY environment variable is required in .env file")
	}
	privateKey, err := keys.RestoreKey(privKey)
	if err != nil {
		return fmt.Errorf("failed to restore private key: %w", err)
	}

	publicKey, err := keys.PublicKey(privateKey)
	if err != nil {
		return fmt.Errorf("[SendEthTx] failed to derive public key: %w", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKey)
	nonce, err := ctx.EthClient().PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("[SendEthTx] failed to get nonce: %w", err)
	}
	gasPrice, err := ctx.EthClient().SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("[SendEthTx] failed to suggest gas price: %w", err)
	}

	chainIdInt, err := ctx.EthClient().ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("[SendEthTx] invalid chain id: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIdInt)
	if err != nil {
		return fmt.Errorf("[SendEthTx] failed to create auth transactor: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	ics07Tendermint, err := tendermintContract.NewContractGroth16ICS07Tendermint(
		*ctx.ClientContract(),
		ctx.EthClient(),
	)
	if err != nil {
		return fmt.Errorf("failed to create ICS07 Tendermint contract: %w", err)
	}

	ics26Router, err := contractICS26Router.NewContractICS26Router(
		*ctx.RouterContract(),
		ctx.EthClient(),
	)
	if err != nil {
		return fmt.Errorf("failed to create ICS07 Tendermint contract: %w", err)
	}

	var tx *types.Transaction
	var txLabel string
	benchEnabled := utils.BenchEnabled()
	var benchStart time.Time
	if benchEnabled {
		benchStart = time.Now()
	}
	switch msg := msg.(type) {
	case updateclient.IUpdateClientMsgsMsgUpdateClient:
		txLabel = "updateClient"
		data, err := relayerclient.EncodeUpdateClientMsg(msg)
		if err != nil {
			return fmt.Errorf("[SendEthTx] failed to encode updateClient msg: %w", err)
		}
		if routerManagesProofSubmission(ctx) {
			log.Printf("[SendEthTx] Sending ICS26Router.updateClient tx for clientId=%s...", cosmosClientID)
			tx, err = ics26Router.UpdateClient(auth, cosmosClientID, data)
			if err != nil {
				return fmt.Errorf("[SendEthTx] failed to send router updateClient tx: %w", err)
			}
		} else {
			log.Printf("[SendEthTx] Sending direct ICS07 updateClient tx...")
			tx, err = ics07Tendermint.UpdateClient(auth, data)
			if err != nil {
				return fmt.Errorf("[SendEthTx] failed to send direct updateClient tx: %w", err)
			}
		}
	case tendermintContract.ILightClientMsgsMsgVerifyMembership:
		txLabel = "verifyMembership"
		if routerManagesProofSubmission(ctx) {
			return fmt.Errorf(
				"[SendEthTx] direct verifyMembership is disabled when ROLE_MANAGER is the ICS26 router; use ICS26Router packet flows instead",
			)
		}
		log.Printf("[SendEthTx] Sending verifyMembership tx...")
		tx, err = ics07Tendermint.VerifyMembership(auth, msg)
		if err != nil {
			return fmt.Errorf("[SendEthTx] failed to verify membership: %w", err)
		}
	case tendermintContract.ILightClientMsgsMsgVerifyNonMembership:
		txLabel = "verifyNonMembership"
		if routerManagesProofSubmission(ctx) {
			return fmt.Errorf(
				"[SendEthTx] direct verifyNonMembership is disabled when ROLE_MANAGER is the ICS26 router; use ICS26Router packet flows instead",
			)
		}
		log.Printf("[SendEthTx] Sending verifyNonMembership tx...")
		tx, err = ics07Tendermint.VerifyNonMembership(auth, msg)
		if err != nil {
			return fmt.Errorf("[SendEthTx] failed to verify non-membership: %w", err)
		}
	case contractICS26Router.IICS26RouterMsgsMsgRecvPacket:
		txLabel = fmt.Sprintf("recvPacket seq=%d", msg.Packet.Sequence)
		log.Printf("[SendEthTx] Sending recvPacket seq=%d...", msg.Packet.Sequence)
		tx, err = ics26Router.RecvPacket(auth, msg)
		if err != nil {
			return fmt.Errorf("[SendEthTx] failed to recv packet: %w", err)
		}
	case contractICS26Router.IICS26RouterMsgsMsgAckPacket:
		txLabel = fmt.Sprintf("ackPacket seq=%d", msg.Packet.Sequence)
		log.Printf("[SendEthTx] Sending ackPacket seq=%d...", msg.Packet.Sequence)
		tx, err = ics26Router.AckPacket(auth, msg)
		if err != nil {
			return fmt.Errorf("[SendEthTx] failed to ack packet: %w", err)
		}
	case contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket:
		txLabel = fmt.Sprintf("timeoutPacket seq=%d", msg.Packet.Sequence)
		log.Printf("[SendEthTx] Sending timeoutPacket seq=%d...", msg.Packet.Sequence)
		tx, err = ics26Router.TimeoutPacket(auth, msg)
		if err != nil {
			return fmt.Errorf("[SendEthTx] failed to timeout packet: %w", err)
		}
	default:
		return fmt.Errorf("[SendEthTx] unsupported message type: %T", msg)
	}

	var submitDur time.Duration
	if benchEnabled {
		submitDur = time.Since(benchStart)
	}
	log.Printf("[SendEthTx] Tx sent: %s. Waiting for receipt...", tx.Hash().Hex())
	var waitStart time.Time
	if benchEnabled {
		waitStart = time.Now()
	}
	receiptCtx, cancel := context.WithTimeout(context.Background(), ethTxReceiptTimeout)
	defer cancel()
	receipt, err := bind.WaitMined(receiptCtx, ctx.EthClient(), tx)
	if err != nil {
		return fmt.Errorf("failed waiting for tx receipt: %w", err)
	}
	if receipt.Status == 0 {
		// Try to get revert reason by replaying the tx via eth_call
		callMsg := ethereum.CallMsg{
			From:     fromAddress,
			To:       tx.To(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice(),
			Value:    tx.Value(),
			Data:     tx.Data(),
		}
		_, callErr := ctx.EthClient().CallContract(context.Background(), callMsg, receipt.BlockNumber)
		if callErr != nil {
			log.Printf("[SendEthTx] Revert reason: %v", callErr)
			// Extract hex-encoded revert data for custom error decoding
			type dataErr interface {
				ErrorData() interface{}
			}
			if de, ok := callErr.(dataErr); ok {
				log.Printf("[SendEthTx] Revert data (hex): %v", de.ErrorData())
			}
		}
		return fmt.Errorf("tx %s reverted (status=0, gasUsed=%d): %w", tx.Hash().Hex(), receipt.GasUsed, services.ErrPermanentRelayFailure)
	}
	var waitDur time.Duration
	if benchEnabled {
		waitDur = time.Since(waitStart)
	}
	log.Printf("[SendEthTx] Tx %s confirmed in block %d (gasUsed=%d)", tx.Hash().Hex(), receipt.BlockNumber.Uint64(), receipt.GasUsed)
	if benchEnabled {
		log.Printf("[bench][eth] %s gasUsed=%d submit=%s wait=%s total=%s tx=%s",
			txLabel, receipt.GasUsed, submitDur, waitDur, time.Since(benchStart), tx.Hash().Hex())
	}

	return nil
}

// SendEthTxBatch packs N packet-level ICS26Router calls into a single
// `multicall(bytes[])` tx. V1 supports only the per-packet msg types that
// `ICS26Router` exposes directly (recvPacket / ackPacket / timeoutPacket).
// `updateClient`, `verifyMembership` and `verifyNonMembership` are NOT
// permitted in a batch — they have separate submission paths.
//
// Empty input is a no-op. A single-msg input is forwarded to SendEthTx to
// avoid the multicall wrapper's small overhead for N=1 (matches the
// threshold contract in the V1 plan).
//
// Semantics: MulticallUpgradeable runs each inner call via delegatecall and
// reverts the whole tx if any inner call reverts — so caller can treat
// success as "every packet in the batch was relayed".
func (h *Handler) SendEthTxBatch(ctx services.Context, msgs []any) error {
	if len(msgs) == 0 {
		return nil
	}
	if len(msgs) == 1 {
		return h.SendEthTx(ctx, msgs[0])
	}

	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed to load ICS26Router ABI: %w", err)
	}

	// cosmosClientID is only consulted when the batch contains an updateClient
	// inner call; resolved lazily to avoid coupling pure-recvPacket batches to
	// the router-managed config check.
	var cosmosClientID string
	var clientIDResolved bool

	// Pack each per-packet call into raw calldata bytes that MulticallUpgradeable
	// will delegatecall back into the same contract.
	calldata := make([][]byte, 0, len(msgs))
	labels := make([]string, 0, len(msgs))
	for i, msg := range msgs {
		var (
			data []byte
			lbl  string
			perr error
		)
		switch m := msg.(type) {
		case contractICS26Router.IICS26RouterMsgsMsgRecvPacket:
			data, perr = parsedABI.Pack("recvPacket", m)
			lbl = fmt.Sprintf("recvPacket:%d", m.Packet.Sequence)
		case contractICS26Router.IICS26RouterMsgsMsgAckPacket:
			data, perr = parsedABI.Pack("ackPacket", m)
			lbl = fmt.Sprintf("ackPacket:%d", m.Packet.Sequence)
		case contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket:
			data, perr = parsedABI.Pack("timeoutPacket", m)
			lbl = fmt.Sprintf("timeoutPacket:%d", m.Packet.Sequence)
		case updateclient.IUpdateClientMsgsMsgUpdateClient:
			// Folding updateClient into a multicall only works when the
			// ICS26Router is the proof submitter — the multicall is dispatched
			// on the router, so every inner call must target a router method.
			// In direct-submission mode SendEthTx calls ICS07 directly; that
			// path can't be expressed inside multicall, so caller must submit
			// updateClient as a standalone tx.
			if !routerManagesProofSubmission(ctx) {
				return fmt.Errorf("[SendEthTxBatch] updateClient cannot be batched when ICS26Router is not the proof submitter; submit it via SendEthTx instead")
			}
			if !clientIDResolved {
				cosmosClientID, perr = cosmosRouterClientID(ctx)
				if perr != nil {
					return fmt.Errorf("[SendEthTxBatch] resolve cosmos client id: %w", perr)
				}
				clientIDResolved = true
			}
			encoded, encErr := relayerclient.EncodeUpdateClientMsg(m)
			if encErr != nil {
				return fmt.Errorf("[SendEthTxBatch] encode updateClient msg %d: %w", i, encErr)
			}
			data, perr = parsedABI.Pack("updateClient", cosmosClientID, encoded)
			lbl = "updateClient"
		default:
			return fmt.Errorf("[SendEthTxBatch] unsupported message type at index %d: %T (only updateClient/recvPacket/ackPacket/timeoutPacket allowed in multicall)", i, msg)
		}
		if perr != nil {
			return fmt.Errorf("[SendEthTxBatch] pack msg %d (%s): %w", i, lbl, perr)
		}
		calldata = append(calldata, data)
		labels = append(labels, lbl)
	}

	labelStr := strings.Join(labels, ",")

	privKey := os.Getenv("ETH_PRIVATE_KEY")
	if privKey == "" {
		return fmt.Errorf("ETH_PRIVATE_KEY environment variable is required in .env file")
	}
	privateKey, err := keys.RestoreKey(privKey)
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed to restore private key: %w", err)
	}
	publicKey, err := keys.PublicKey(privateKey)
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed to derive public key: %w", err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	nonce, err := ctx.EthClient().PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed to get nonce: %w", err)
	}
	gasPrice, err := ctx.EthClient().SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed to suggest gas price: %w", err)
	}
	chainIdInt, err := ctx.EthClient().ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] invalid chain id: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIdInt)
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed to create auth transactor: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(16000000)
	auth.GasPrice = gasPrice

	ics26Router, err := contractICS26Router.NewContractICS26Router(*ctx.RouterContract(), ctx.EthClient())
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed to bind ICS26Router: %w", err)
	}

	benchEnabled := utils.BenchEnabled()
	var benchStart time.Time
	if benchEnabled {
		benchStart = time.Now()
	}
	log.Printf("[SendEthTxBatch] Submitting multicall: %d inner calls (%s)", len(calldata), labelStr)
	tx, err := ics26Router.Multicall(auth, calldata)
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed to submit multicall: %w", err)
	}
	var submitDur time.Duration
	if benchEnabled {
		submitDur = time.Since(benchStart)
	}
	log.Printf("[SendEthTxBatch] Tx sent: %s. Waiting for receipt...", tx.Hash().Hex())

	var waitStart time.Time
	if benchEnabled {
		waitStart = time.Now()
	}
	receiptCtx, cancel := context.WithTimeout(context.Background(), ethTxReceiptTimeout)
	defer cancel()
	receipt, err := bind.WaitMined(receiptCtx, ctx.EthClient(), tx)
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed waiting for tx receipt: %w", err)
	}
	if receipt.Status == 0 {
		// Replay against the same calldata so MulticallUpgradeable re-throws the
		// first inner revert and we can extract its selector + args.
		callMsg := ethereum.CallMsg{
			From:     fromAddress,
			To:       tx.To(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice(),
			Value:    tx.Value(),
			Data:     tx.Data(),
		}
		_, callErr := ctx.EthClient().CallContract(context.Background(), callMsg, receipt.BlockNumber)
		if callErr != nil {
			log.Printf("[SendEthTxBatch] Revert reason: %v", callErr)
			type dataErr interface {
				ErrorData() interface{}
			}
			if de, ok := callErr.(dataErr); ok {
				log.Printf("[SendEthTxBatch] Revert data (hex): %v", de.ErrorData())
			}
		}
		return fmt.Errorf("multicall tx %s reverted (status=0, gasUsed=%d, labels=%s): %w", tx.Hash().Hex(), receipt.GasUsed, labelStr, services.ErrPermanentRelayFailure)
	}
	var waitDur time.Duration
	if benchEnabled {
		waitDur = time.Since(waitStart)
	}
	log.Printf("[SendEthTxBatch] Tx %s confirmed in block %d (gasUsed=%d, inner=%d)",
		tx.Hash().Hex(), receipt.BlockNumber.Uint64(), receipt.GasUsed, len(calldata))
	if benchEnabled {
		log.Printf("[bench][eth] multicall labels=%s gasUsed=%d submit=%s wait=%s total=%s tx=%s",
			labelStr, receipt.GasUsed, submitDur, waitDur, time.Since(benchStart), tx.Hash().Hex())
	}
	if benchEnabled {
		if traceErr := logInnerGasFromTrace(ctx, tx.Hash(), labels); traceErr != nil {
			log.Printf("[bench][eth] inner gas trace unavailable (RPC may lack debug_ namespace): %v", traceErr)
		}
	}

	return nil
}

func (h *Handler) CreateEthClient(svcCtx services.Context, clientState exported.ClientState, consensusState exported.ConsensusState) (string, error) {
	log.Printf("[CreateEthClientTx] starting")
	cosmosClientID, err := cosmosRouterClientID(svcCtx)
	if err != nil {
		return "", fmt.Errorf("[CreateEthClientTx] %w", err)
	}

	// Get the private key from environment variable
	privKeyHex := os.Getenv("COSMOS_PRIVATE_KEY")
	if privKeyHex == "" {
		return "", fmt.Errorf("COSMOS_PRIVATE_KEY environment variable is required in .env file")
	}

	// Decode the private key
	privKeyBytes, err := hex.DecodeString(strings.TrimPrefix(privKeyHex, "0x"))
	if err != nil {
		return "", fmt.Errorf("failed to decode private key: %w", err)
	}

	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	signerAddr := sdk.AccAddress(privKey.PubKey().Address())
	log.Printf("[CreateEthClient] signer: %s", signerAddr.String())

	// Get chain configuration from environment
	chainID := os.Getenv("COSMOS_CHAIN_ID")
	if chainID == "" {
		return "", fmt.Errorf("COSMOS_CHAIN_ID environment variable is required in .env file")
	}

	// Get gas and fee configuration
	gasLimit := uint64(1500000) // Default gas limit for MsgCreateClient with wasm payload
	if gasStr := os.Getenv("COSMOS_GAS_LIMIT"); gasStr != "" {
		if _, err := fmt.Sscanf(gasStr, "%d", &gasLimit); err != nil {
			return "", fmt.Errorf("failed to parse COSMOS_GAS_LIMIT: %w", err)
		}
	}

	feeDenom := os.Getenv("COSMOS_FEE_DENOM")
	if feeDenom == "" {
		feeDenom = "stake" // Default fee denom
	}

	feeAmount := int64(10000000) // Default fee amount
	if feeStr := os.Getenv("COSMOS_FEE_AMOUNT"); feeStr != "" {
		if _, err := fmt.Sscanf(feeStr, "%d", &feeAmount); err != nil {
			return "", fmt.Errorf("failed to parse COSMOS_FEE_AMOUNT: %w", err)
		}
	}
	log.Printf("[CreateEthClientTx] gas config: gasLimit=%d fee=%d%s", gasLimit, feeAmount, feeDenom)

	// Query account info (account number and sequence) from the chain
	log.Printf("[CreateEthClientTx] querying cosmos account info")
	accountNumber, sequence, err := h.queryAccountInfo(svcCtx, signerAddr.String())
	if err != nil {
		return "", fmt.Errorf("failed to query account info: %w", err)
	}
	log.Printf("[CreateEthClientTx] account info: accountNumber=%d sequence=%d", accountNumber, sequence)

	// Setup encoding config
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	channeltypesv2.RegisterInterfaces(interfaceRegistry)
	clienttypes.RegisterInterfaces(interfaceRegistry)
	ibcwasmtypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)

	msg, err := clienttypes.NewMsgCreateClient(clientState, consensusState, signerAddr.String())
	if err != nil {
		return "", err
	}
	log.Printf("[CreateEthClientTx] MsgCreateClient built")

	// Build the transaction
	txBuilder := txConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		return "", fmt.Errorf("failed to set messages: %w", err)
	}

	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, sdkmath.NewInt(feeAmount))))

	// First, set an empty signature to populate signer info for sign bytes generation
	pubKey := privKey.PubKey()
	emptySig := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	if err := txBuilder.SetSignatures(emptySig); err != nil {
		return "", fmt.Errorf("failed to set empty signature: %w", err)
	}

	// Create signer data
	signerData := authsigning.SignerData{
		Address:       signerAddr.String(),
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	// Get sign bytes using the adapter function
	signBytes, err := authsigning.GetSignBytesAdapter(
		context.Background(),
		txConfig.SignModeHandler(),
		sdksigning.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get sign bytes: %w", err)
	}

	// Sign the bytes
	sigRaw, err := privKey.Sign(signBytes)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Set the actual signature
	sigV2 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: sigRaw,
		},
		Sequence: sequence,
	}

	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return "", fmt.Errorf("failed to set signatures: %w", err)
	}

	// Encode the transaction
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return "", fmt.Errorf("failed to encode transaction: %w", err)
	}

	// Broadcast the transaction
	log.Printf("[CreateEthClientTx] broadcasting MsgCreateClient")
	result, err := svcCtx.CosmosClient().BroadcastTxSync(context.Background(), txBytes)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("transaction failed with code %d: %s", result.Code, result.Log)
	}

	log.Printf("[CreateEthClient] MsgCreateClient broadcast successfully. Hash: %s", result.Hash.String())

	// Wait for MsgCreateClient tx and extract the new client ID from events
	log.Printf("[CreateEthClientTx] waiting for MsgCreateClient tx result")
	txResult, err := h.waitForTxResult(svcCtx, result.Hash, 30*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed waiting for MsgCreateClient tx: %w", err)
	}

	if txResult.TxResult.Code != 0 {
		return "", fmt.Errorf("MsgCreateClient tx failed at DeliverTx: code=%d log=%s", txResult.TxResult.Code, txResult.TxResult.Log)
	}

	newClientID := extractClientID(txResult)
	if newClientID == "" {
		return "", fmt.Errorf("MsgCreateClient tx confirmed but client_id not found in events")
	}
	log.Printf("[CreateEthClient] new ETH client ID: %s", newClientID)

	// Build and broadcast MsgRegisterCounterparty as a separate transaction
	registerMsg := clienttypesv2.NewMsgRegisterCounterparty(
		newClientID,
		[][]byte{[]byte("")},
		cosmosClientID,
		signerAddr.String(),
	)
	log.Printf("[CreateEthClientTx] MsgRegisterCounterparty built for clientID=%s", newClientID)

	// Re-query account info (sequence incremented after first tx)
	log.Printf("[CreateEthClientTx] querying cosmos account info for register counterparty")
	accountNumber, sequence, err = h.queryAccountInfo(svcCtx, signerAddr.String())
	if err != nil {
		return "", fmt.Errorf("failed to query account info for register counterparty: %w", err)
	}
	log.Printf("[CreateEthClientTx] register counterparty account info: accountNumber=%d sequence=%d", accountNumber, sequence)

	txBuilder2 := txConfig.NewTxBuilder()
	if err := txBuilder2.SetMsgs(registerMsg); err != nil {
		return "", fmt.Errorf("failed to set register counterparty message: %w", err)
	}
	txBuilder2.SetGasLimit(gasLimit)
	txBuilder2.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, sdkmath.NewInt(feeAmount))))

	emptySig2 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	if err := txBuilder2.SetSignatures(emptySig2); err != nil {
		return "", fmt.Errorf("failed to set empty signature: %w", err)
	}

	signerData2 := authsigning.SignerData{
		Address:       signerAddr.String(),
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	signBytes2, err := authsigning.GetSignBytesAdapter(
		context.Background(),
		txConfig.SignModeHandler(),
		sdksigning.SignMode_SIGN_MODE_DIRECT,
		signerData2,
		txBuilder2.GetTx(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get sign bytes for register counterparty: %w", err)
	}

	sigRaw2, err := privKey.Sign(signBytes2)
	if err != nil {
		return "", fmt.Errorf("failed to sign register counterparty tx: %w", err)
	}

	sigV22 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: sigRaw2,
		},
		Sequence: sequence,
	}
	if err := txBuilder2.SetSignatures(sigV22); err != nil {
		return "", fmt.Errorf("failed to set signatures: %w", err)
	}

	txBytes2, err := txConfig.TxEncoder()(txBuilder2.GetTx())
	if err != nil {
		return "", fmt.Errorf("failed to encode register counterparty tx: %w", err)
	}

	log.Printf("[CreateEthClientTx] broadcasting MsgRegisterCounterparty")
	result2, err := svcCtx.CosmosClient().BroadcastTxSync(context.Background(), txBytes2)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast register counterparty tx: %w", err)
	}

	if result2.Code != 0 {
		return "", fmt.Errorf("register counterparty tx failed with code %d: %s", result2.Code, result2.Log)
	}

	log.Printf("[CreateEthClient] MsgRegisterCounterparty broadcast successfully. Hash: %s", result2.Hash.String())

	return newClientID, nil
}

// extractClientID parses the client_id from the MsgCreateClient tx result events.
func extractClientID(txResult *coretypes.ResultTx) string {
	if txResult == nil {
		return ""
	}
	for _, event := range txResult.TxResult.Events {
		// Try all events — different ibc-go versions use different event types
		// (e.g. "create_client", "ibc_client", or attributes on generic events)
		for _, attr := range event.Attributes {
			if attr.Key == "client_id" {
				return attr.Value
			}
		}
	}
	return ""
}

// CosmosSignerAddress returns the bech32 address derived from COSMOS_PRIVATE_KEY.
// Used to populate the Signer field in Cosmos messages before calling SendCosmosTx.
func (h *Handler) CosmosSignerAddress() (string, error) {
	privKeyHex := strings.TrimPrefix(os.Getenv("COSMOS_PRIVATE_KEY"), "0x")
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode COSMOS_PRIVATE_KEY: %w", err)
	}
	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	return sdk.AccAddress(privKey.PubKey().Address()).String(), nil
}

func (h *Handler) SendCosmosTx(svcCtx services.Context, msg any) error {
	benchEnabled := utils.BenchEnabled()
	var benchStart time.Time
	if benchEnabled {
		benchStart = time.Now()
	}
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("message must be a proto.Message")
	}
	msgLabel := fmt.Sprintf("%T", protoMsg)
	// Get the private key from environment variable
	privKeyHex := os.Getenv("COSMOS_PRIVATE_KEY")
	if privKeyHex == "" {
		return fmt.Errorf("COSMOS_PRIVATE_KEY environment variable is required in .env file")
	}

	// Decode the private key
	privKeyBytes, err := hex.DecodeString(strings.TrimPrefix(privKeyHex, "0x"))
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	signerAddr := sdk.AccAddress(privKey.PubKey().Address())

	// Get chain configuration from environment
	chainID := os.Getenv("COSMOS_CHAIN_ID")
	if chainID == "" {
		return fmt.Errorf("COSMOS_CHAIN_ID environment variable is required in .env file")
	}

	// Get gas and fee configuration
	gasLimit := uint64(200000) // Default gas limit
	if gasStr := os.Getenv("COSMOS_GAS_LIMIT"); gasStr != "" {
		if _, err := fmt.Sscanf(gasStr, "%d", &gasLimit); err != nil {
			return fmt.Errorf("failed to parse COSMOS_GAS_LIMIT: %w", err)
		}
	}

	feeDenom := os.Getenv("COSMOS_FEE_DENOM")
	if feeDenom == "" {
		feeDenom = "stake" // Default fee denom
	}

	// Query account info (account number and sequence) from the chain
	accountNumber, sequence, err := h.queryAccountInfo(svcCtx, signerAddr.String())
	if err != nil {
		return fmt.Errorf("failed to query account info: %w", err)
	}

	// Setup encoding config
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	channeltypesv2.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)

	// Build the transaction
	txBuilder := txConfig.NewTxBuilder()

	// Convert the proto.Message to sdk.Msg
	sdkMsg, ok := protoMsg.(sdk.Msg)
	if !ok {
		return fmt.Errorf("message does not implement sdk.Msg interface")
	}

	// Fill empty Signer fields and apply per-message gas overrides.
	// Wasm light client verification (MsgUpdateClient, MsgRecvPacket, MsgAcknowledgement,
	// MsgTimeout) runs the cw-ics08-wasm-eth contract and exceeds the 200k default.
	wasmHeavyGas := os.Getenv("COSMOS_GAS_LIMIT") == ""
	if msg, ok := sdkMsg.(*clienttypes.MsgUpdateClient); ok {
		if msg.Signer == "" {
			msg.Signer = signerAddr.String()
		}
		if wasmHeavyGas {
			gasLimit = uint64(2000000)
		}
	}
	if msg, ok := sdkMsg.(*channeltypesv2.MsgRecvPacket); ok {
		if msg.Signer == "" {
			msg.Signer = signerAddr.String()
		}
		if wasmHeavyGas {
			gasLimit = uint64(2000000)
		}
	}
	if msg, ok := sdkMsg.(*channeltypesv2.MsgAcknowledgement); ok {
		if msg.Signer == "" {
			msg.Signer = signerAddr.String()
		}
		if wasmHeavyGas {
			gasLimit = uint64(2000000)
		}
	}
	if msg, ok := sdkMsg.(*channeltypesv2.MsgTimeout); ok {
		if msg.Signer == "" {
			msg.Signer = signerAddr.String()
		}
		if wasmHeavyGas {
			gasLimit = uint64(2000000)
		}
	}

	if err := txBuilder.SetMsgs(sdkMsg); err != nil {
		return fmt.Errorf("failed to set messages: %w", err)
	}

	// Set fee to match gas limit (gasPrice = 1stake per gas unit)
	feeAmount := int64(gasLimit)
	if feeStr := os.Getenv("COSMOS_FEE_AMOUNT"); feeStr != "" {
		if _, err := fmt.Sscanf(feeStr, "%d", &feeAmount); err != nil {
			return fmt.Errorf("failed to parse COSMOS_FEE_AMOUNT: %w", err)
		}
	}

	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, sdkmath.NewInt(feeAmount))))

	// First, set an empty signature to populate signer info for sign bytes generation
	pubKey := privKey.PubKey()
	emptySig := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	if err := txBuilder.SetSignatures(emptySig); err != nil {
		return fmt.Errorf("failed to set empty signature: %w", err)
	}

	// Create signer data
	signerData := authsigning.SignerData{
		Address:       signerAddr.String(),
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	// Get sign bytes using the adapter function
	signBytes, err := authsigning.GetSignBytesAdapter(
		context.Background(),
		txConfig.SignModeHandler(),
		sdksigning.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		return fmt.Errorf("failed to get sign bytes: %w", err)
	}

	// Sign the bytes
	sigRaw, err := privKey.Sign(signBytes)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Set the actual signature
	sigV2 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: sigRaw,
		},
		Sequence: sequence,
	}

	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return fmt.Errorf("failed to set signatures: %w", err)
	}

	// Encode the transaction
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return fmt.Errorf("failed to encode transaction: %w", err)
	}

	// Broadcast the transaction using BroadcastTxCommit for detailed error info
	var broadcastStart time.Time
	if benchEnabled {
		broadcastStart = time.Now()
	}
	commitResult, err := svcCtx.CosmosClient().BroadcastTxCommit(context.Background(), txBytes)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}
	var broadcastDur time.Duration
	if benchEnabled {
		broadcastDur = time.Since(broadcastStart)
	}

	if commitResult.CheckTx.Code != 0 {
		log.Printf("[SendCosmosTx] CheckTx FAILED: code=%d codespace=%s log=%s info=%s data=%x",
			commitResult.CheckTx.Code, commitResult.CheckTx.Codespace, commitResult.CheckTx.Log,
			commitResult.CheckTx.Info, commitResult.CheckTx.Data)
		return fmt.Errorf("transaction failed at CheckTx with code %d: %s", commitResult.CheckTx.Code, commitResult.CheckTx.Log)
	}

	if commitResult.TxResult.Code != 0 {
		log.Printf("[SendCosmosTx] DeliverTx FAILED: code=%d codespace=%s log=%s data=%x",
			commitResult.TxResult.Code, commitResult.TxResult.Codespace, commitResult.TxResult.Log, commitResult.TxResult.Data)
		// DeliverTx execution failure is deterministic (msg/proof rejected) — mark
		// permanent so the relay loop counts it toward the retry cap.
		return fmt.Errorf("transaction failed at DeliverTx with code %d: %s: %w", commitResult.TxResult.Code, commitResult.TxResult.Log, services.ErrPermanentRelayFailure)
	}

	log.Printf("[SendCosmosTx] Tx confirmed at height %d hash=%s", commitResult.Height, commitResult.Hash.String())
	if benchEnabled {
		log.Printf("[bench][cosmos] %s gasWanted=%d gasUsed=%d broadcast=%s total=%s height=%d hash=%s",
			msgLabel, commitResult.TxResult.GasWanted, commitResult.TxResult.GasUsed,
			broadcastDur, time.Since(benchStart), commitResult.Height, commitResult.Hash.String())
	}
	return nil
}

// SendCosmosTxBatch sends multiple messages in a single Cosmos transaction
func (h *Handler) SendCosmosTxBatch(svcCtx services.Context, msgs []any) error {
	if len(msgs) == 0 {
		return nil
	}
	benchEnabled := utils.BenchEnabled()
	var benchStart time.Time
	if benchEnabled {
		benchStart = time.Now()
	}

	// Convert all messages to sdk.Msg
	// Get the private key from environment variable
	privKeyHex := os.Getenv("COSMOS_PRIVATE_KEY")
	if privKeyHex == "" {
		return fmt.Errorf("COSMOS_PRIVATE_KEY environment variable is required in .env file")
	}

	// Decode the private key
	privKeyBytes, err := hex.DecodeString(strings.TrimPrefix(privKeyHex, "0x"))
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	signerAddr := sdk.AccAddress(privKey.PubKey().Address())

	// Convert all messages to sdk.Msg, filling empty Signer fields
	var sdkMsgs []sdk.Msg
	for i, msg := range msgs {
		protoMsg, ok := msg.(proto.Message)
		if !ok {
			return fmt.Errorf("message %d must be a proto.Message", i)
		}
		sdkMsg, ok := protoMsg.(sdk.Msg)
		if !ok {
			return fmt.Errorf("message %d does not implement sdk.Msg interface", i)
		}
		if m, ok := sdkMsg.(*clienttypes.MsgUpdateClient); ok && m.Signer == "" {
			m.Signer = signerAddr.String()
		}
		if m, ok := sdkMsg.(*channeltypesv2.MsgTimeout); ok && m.Signer == "" {
			m.Signer = signerAddr.String()
		}
		sdkMsgs = append(sdkMsgs, sdkMsg)
	}

	// Get chain configuration from environment
	chainID := os.Getenv("COSMOS_CHAIN_ID")
	if chainID == "" {
		return fmt.Errorf("COSMOS_CHAIN_ID environment variable is required in .env file")
	}

	// Get gas and fee configuration - use higher gas for batch transactions
	baseGas := uint64(200000)
	if gasStr := os.Getenv("COSMOS_GAS_LIMIT"); gasStr != "" {
		if _, err := fmt.Sscanf(gasStr, "%d", &baseGas); err != nil {
			return fmt.Errorf("failed to parse COSMOS_GAS_LIMIT: %w", err)
		}
	}
	// MsgUpdateClient requires significantly more gas due to wasm verification
	for _, msg := range sdkMsgs {
		if _, ok := msg.(*clienttypes.MsgUpdateClient); ok {
			if baseGas < 2000000 {
				baseGas = 2000000
			}
			break
		}
	}
	gasLimit := baseGas * uint64(len(sdkMsgs))

	feeDenom := os.Getenv("COSMOS_FEE_DENOM")
	if feeDenom == "" {
		feeDenom = "stake" // Default fee denom
	}

	// Set fee to match gas limit (gasPrice = 1stake per gas unit)
	feeAmount := int64(gasLimit)
	if feeStr := os.Getenv("COSMOS_FEE_AMOUNT"); feeStr != "" {
		var baseFee int64
		if _, err := fmt.Sscanf(feeStr, "%d", &baseFee); err != nil {
			return fmt.Errorf("failed to parse COSMOS_FEE_AMOUNT: %w", err)
		}
		feeAmount = baseFee * int64(len(sdkMsgs))
	}

	// Query account info (account number and sequence) from the chain
	accountNumber, sequence, err := h.queryAccountInfo(svcCtx, signerAddr.String())
	if err != nil {
		return fmt.Errorf("failed to query account info: %w", err)
	}

	// Setup encoding config
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	channeltypesv2.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)

	// Build the transaction
	txBuilder := txConfig.NewTxBuilder()

	if err := txBuilder.SetMsgs(sdkMsgs...); err != nil {
		return fmt.Errorf("failed to set messages: %w", err)
	}

	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, sdkmath.NewInt(feeAmount))))

	// First, set an empty signature to populate signer info for sign bytes generation
	pubKey := privKey.PubKey()
	emptySig := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	if err := txBuilder.SetSignatures(emptySig); err != nil {
		return fmt.Errorf("failed to set empty signature: %w", err)
	}

	// Create signer data
	signerData := authsigning.SignerData{
		Address:       signerAddr.String(),
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	// Get sign bytes using the adapter function
	signBytes, err := authsigning.GetSignBytesAdapter(
		context.Background(),
		txConfig.SignModeHandler(),
		sdksigning.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		return fmt.Errorf("failed to get sign bytes: %w", err)
	}

	// Sign the bytes
	sigRaw, err := privKey.Sign(signBytes)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Set the actual signature
	sigV2 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: sigRaw,
		},
		Sequence: sequence,
	}

	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return fmt.Errorf("failed to set signatures: %w", err)
	}

	// Encode the transaction
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return fmt.Errorf("failed to encode transaction: %w", err)
	}

	// Broadcast the transaction using BroadcastTxCommit for detailed error info
	var broadcastStart time.Time
	if benchEnabled {
		broadcastStart = time.Now()
	}
	commitResult, err := svcCtx.CosmosClient().BroadcastTxCommit(context.Background(), txBytes)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}
	var broadcastDur time.Duration
	if benchEnabled {
		broadcastDur = time.Since(broadcastStart)
	}

	if commitResult.CheckTx.Code != 0 {
		log.Printf("[SendCosmosTxBatch] CheckTx FAILED: code=%d codespace=%s log=%s info=%s data=%x",
			commitResult.CheckTx.Code, commitResult.CheckTx.Codespace, commitResult.CheckTx.Log,
			commitResult.CheckTx.Info, commitResult.CheckTx.Data)
		return fmt.Errorf("transaction failed at CheckTx with code %d: %s", commitResult.CheckTx.Code, commitResult.CheckTx.Log)
	}

	if commitResult.TxResult.Code != 0 {
		log.Printf("[SendCosmosTxBatch] DeliverTx FAILED: code=%d codespace=%s log=%s data=%x",
			commitResult.TxResult.Code, commitResult.TxResult.Codespace, commitResult.TxResult.Log, commitResult.TxResult.Data)
		// DeliverTx execution failure is deterministic — mark permanent.
		return fmt.Errorf("transaction failed at DeliverTx with code %d: %s: %w", commitResult.TxResult.Code, commitResult.TxResult.Log, services.ErrPermanentRelayFailure)
	}

	log.Printf("[SendCosmosTxBatch] Tx confirmed at height %d hash=%s (msgs=%d)", commitResult.Height, commitResult.Hash.String(), len(sdkMsgs))
	if benchEnabled {
		log.Printf("[bench][cosmos] batch msgs=%d gasWanted=%d gasUsed=%d broadcast=%s total=%s height=%d hash=%s",
			len(sdkMsgs), commitResult.TxResult.GasWanted, commitResult.TxResult.GasUsed,
			broadcastDur, time.Since(benchStart), commitResult.Height, commitResult.Hash.String())
	}

	return nil
}

// queryAccountInfo queries the account number and sequence for the given address
func (h *Handler) queryAccountInfo(svcCtx services.Context, address string) (uint64, uint64, error) {
	// Build the query request
	queryReq := &authtypes.QueryAccountRequest{
		Address: address,
	}

	reqBytes, err := proto.Marshal(queryReq)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to marshal query request: %w", err)
	}

	// Query path for auth account
	queryPath := "/cosmos.auth.v1beta1.Query/Account"

	// Make ABCI query
	result, err := svcCtx.CosmosClient().ABCIQuery(context.Background(), queryPath, reqBytes)
	if err != nil {
		return 0, 0, fmt.Errorf("ABCI query failed: %w", err)
	}

	if result.Response.Code != 0 {
		return 0, 0, fmt.Errorf("query failed with code %d: %s", result.Response.Code, result.Response.Log)
	}

	// Setup interface registry to decode the account
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)

	// Decode the response
	var queryResp authtypes.QueryAccountResponse
	if err := proto.Unmarshal(result.Response.Value, &queryResp); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal query response: %w", err)
	}

	// Unpack the account from Any type
	var account sdk.AccountI
	if err := interfaceRegistry.UnpackAny(queryResp.Account, &account); err != nil {
		return 0, 0, fmt.Errorf("failed to unpack account: %w", err)
	}

	return account.GetAccountNumber(), account.GetSequence(), nil
}

// waitForTxResult polls the chain until the transaction with the given hash is included in a block or the timeout expires.
func (h *Handler) waitForTxResult(svcCtx services.Context, txHash []byte, timeout time.Duration) (*coretypes.ResultTx, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		result, err := svcCtx.CosmosClient().Tx(context.Background(), txHash, false)
		if err == nil && result != nil && result.Height > 0 {
			log.Printf("[WaitForTx] Tx %X confirmed at height %d", txHash, result.Height)
			return result, nil
		}
		time.Sleep(1 * time.Second)
	}
	return nil, fmt.Errorf("timeout waiting for tx %X to be included in a block", txHash)
}
