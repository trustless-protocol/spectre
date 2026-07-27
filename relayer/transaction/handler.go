package transaction

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	routerContract "relayer/bindings/ICS26Router"
	spectreContract "relayer/bindings/SpectreClient"
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
	sdkbech32 "github.com/cosmos/cosmos-sdk/types/bech32"
	txservice "github.com/cosmos/cosmos-sdk/types/tx"
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
	"github.com/ethereum/go-ethereum/ethclient"
)

var validatorCacheRaceErrorSelectors = map[[4]byte]string{
	errorSelector("ValidatorSetCacheMiss(bytes32)"): "ValidatorSetCacheMiss",
}

func errorSelector(signature string) [4]byte {
	hash := crypto.Keccak256([]byte(signature))
	var selector [4]byte
	copy(selector[:], hash[:4])
	return selector
}

func validatorCacheRaceErrorName(callErr error) (string, bool) {
	type dataErr interface {
		ErrorData() interface{}
	}
	de, ok := callErr.(dataErr)
	if !ok {
		return "", false
	}

	var data []byte
	switch raw := de.ErrorData().(type) {
	case string:
		data = common.FromHex(raw)
	case fmt.Stringer:
		data = common.FromHex(raw.String())
	case []byte:
		data = raw
	default:
		return "", false
	}
	if len(data) < 4 {
		return "", false
	}

	var selector [4]byte
	copy(selector[:], data[:4])
	name, ok := validatorCacheRaceErrorSelectors[selector]
	return name, ok
}

type Handler struct {
	mu            sync.Mutex
	cosmosMu      sync.Mutex
	nonce         uint64
	nonceValid    bool
	lastNonce     uint64
	lastGasPrice  *big.Int
	lastGasFeeCap *big.Int
	lastGasTipCap *big.Int
}

var ethTxBroadcastTimeout = 30 * time.Second
var ethTxReceiptTimeout = 45 * time.Second
var ethTxReceiptPollInterval = 2 * time.Second

const ethDeployGasHeadroomPercent uint64 = 20

// Cosmos RPC deadlines (issue #119): bound every Cosmos broadcast/query so a
// stalled node can never wedge the relay goroutine. The relay path previously
// used BroadcastTxCommit, which blocks until the tx is committed — or forever.
const cosmosRPCTimeout = 30 * time.Second       // a single Cosmos RPC call (query / sync broadcast)
const cosmosInclusionTimeout = 90 * time.Second // total bounded wait for a broadcast tx to land in a block

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
func logInnerGasFromTrace(stdCtx context.Context, ctx services.Context, txHash common.Hash, labels []string) error {
	rpcClient := ctx.EthClient().Client()
	if rpcClient == nil {
		return fmt.Errorf("nil rpc client")
	}

	var root callTrace
	tracerCfg := map[string]any{"tracer": "callTracer"}
	if err := rpcClient.CallContext(stdCtx, &root, "debug_traceTransaction", txHash, tracerCfg); err != nil {
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

func toGroth16ValidatorSet(in relayerclient.ContractValidatorSet) spectreContract.IICS07TendermintMsgsValidatorSet {
	validators := make([]spectreContract.IICS07TendermintMsgsValidatorInfo, len(in.Validators))
	for i, val := range in.Validators {
		validators[i] = spectreContract.IICS07TendermintMsgsValidatorInfo{
			ValAddress:       val.ValAddress,
			PubKey:           val.PubKey,
			VotingPower:      val.VotingPower,
			ProposerPriority: val.ProposerPriority,
		}
	}
	return spectreContract.IICS07TendermintMsgsValidatorSet{
		Validators:  validators,
		HasProposer: in.HasProposer,
		Proposer: spectreContract.IICS07TendermintMsgsValidatorInfo{
			ValAddress:       in.Proposer.ValAddress,
			PubKey:           in.Proposer.PubKey,
			VotingPower:      in.Proposer.VotingPower,
			ProposerPriority: in.Proposer.ProposerPriority,
		},
		TotalVotingPower: in.TotalVotingPower,
	}
}

func estimateCosmosClientDeployGas(
	stdCtx context.Context,
	ctx services.Context,
	from common.Address,
	gasPrice *big.Int,
	clientState []byte,
	consensusHash []byte,
	initialPinnedValidatorSet spectreContract.IICS07TendermintMsgsValidatorSet,
) (uint64, uint64, error) {
	parsed, err := spectreContract.ContractSpectreClientMetaData.GetAbi()
	if err != nil {
		return 0, 0, fmt.Errorf("parse ICS07 ABI: %w", err)
	}
	constructorInput, err := parsed.Pack(
		"",
		*ctx.UpdateClientContract(),
		*ctx.MembershipContract(),
		*ctx.MisbehaviourContract(),
		clientState,
		utils.BytesToBytes32(consensusHash),
		initialPinnedValidatorSet,
		*ctx.RoleManagerAddress(),
	)
	if err != nil {
		return 0, 0, fmt.Errorf("pack ICS07 constructor args: %w", err)
	}
	deployData := append(common.FromHex(spectreContract.ContractSpectreClientBin), constructorInput...)
	estimate, err := ctx.EthClient().EstimateGas(stdCtx, ethereum.CallMsg{
		From:     from,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     deployData,
	})
	if err != nil {
		return 0, 0, fmt.Errorf("estimate ICS07 deploy gas: %w", err)
	}

	gasLimit := estimate + (estimate*ethDeployGasHeadroomPercent)/100
	if gasLimit < estimate {
		gasLimit = estimate
	}
	if header, err := ctx.EthClient().HeaderByNumber(stdCtx, nil); err == nil && header.GasLimit > 0 {
		if estimate >= header.GasLimit {
			return estimate, 0, fmt.Errorf(
				"estimated ICS07 deploy gas %d exceeds latest block gas limit %d",
				estimate,
				header.GasLimit,
			)
		}
		if gasLimit >= header.GasLimit {
			gasLimit = header.GasLimit - 1
		}
	}
	return estimate, gasLimit, nil
}

func (h *Handler) CreateCosmosClientContract(stdCtx context.Context, ctx services.Context, clientState, consensusHash []byte, initialPinnedValidatorSet relayerclient.ContractValidatorSet) (common.Address, error) {
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

	gasPrice, err := ctx.EthClient().SuggestGasPrice(stdCtx)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] failed to suggest gas price: %w", err)
	}

	pinnedForDeploy := toGroth16ValidatorSet(initialPinnedValidatorSet)
	estimatedDeployGas, deployGasLimit, err := estimateCosmosClientDeployGas(stdCtx, ctx, fromAddress, gasPrice, clientState, consensusHash, pinnedForDeploy)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] %w", err)
	}
	log.Printf("[CreateCosmosClient] ICS07 deploy gas estimate=%d limit=%d", estimatedDeployGas, deployGasLimit)

	var address common.Address
	deployFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		addr, tx, _, err := spectreContract.DeployContractSpectreClient(
			auth,
			ctx.EthClient(),
			*ctx.UpdateClientContract(),
			*ctx.MembershipContract(),
			*ctx.MisbehaviourContract(),
			clientState,
			utils.BytesToBytes32(consensusHash),
			pinnedForDeploy,
			*ctx.RoleManagerAddress(),
		)
		if err == nil {
			address = addr
		}
		return tx, err
	}

	receipt, _, _, err := h.executeWithRetryAndResubmission(stdCtx, ctx, privateKey, deployGasLimit, deployFn)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed waiting for deploy receipt: %w", err)
	}
	log.Printf("[CreateCosmosClient] ICS07 deployed at %s (block %d, gasUsed=%d)", address.String(), receipt.BlockNumber.Uint64(), receipt.GasUsed)
	ctx.SetClient(address)

	ics26Router, err := routerContract.NewContractICS26Router(*ctx.RouterContract(), ctx.EthClient())
	if err != nil {
		return common.Address{}, err
	}

	addClientFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return ics26Router.AddClient(
			auth,
			cosmosClientID,
			routerContract.IICS02ClientMsgsCounterpartyInfo{
				ClientId:     wasmClientID,
				MerklePrefix: [][]byte{[]byte("ibc"), []byte("")},
			},
			*ctx.SpectreClientContract(),
		)
	}

	addClientReceipt, _, _, err := h.executeWithRetryAndResubmission(stdCtx, ctx, privateKey, 16000000, addClientFn)
	if err != nil {
		if errors.Is(err, services.ErrPermanentRelayFailure) {
			// Fallback to MigrateClient only for confirmed on-chain reverts
			log.Printf("[CreateCosmosClient] AddClient failed permanently (%v) — falling back to MigrateClient to repoint %s to new ICS07 %s",
				err, cosmosClientID, address.Hex())

			migrateClientFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
				return ics26Router.MigrateClient(
					auth,
					cosmosClientID,
					routerContract.IICS02ClientMsgsCounterpartyInfo{
						ClientId:     wasmClientID,
						MerklePrefix: [][]byte{[]byte("ibc"), []byte("")},
					},
					*ctx.SpectreClientContract(),
				)
			}

			migrateReceipt, _, _, mErr := h.executeWithRetryAndResubmission(stdCtx, ctx, privateKey, 16000000, migrateClientFn)
			if mErr != nil {
				return common.Address{}, fmt.Errorf("MigrateClient call failed: %w", mErr)
			}
			log.Printf("[CreateCosmosClient] MigrateClient confirmed (block %d, gasUsed=%d)", migrateReceipt.BlockNumber.Uint64(), migrateReceipt.GasUsed)
			return address, nil
		}
		// For transient wait/RPC errors, return the error immediately so the caller can retry
		return common.Address{}, err
	}

	log.Printf("[CreateCosmosClient] AddClient confirmed (block %d, gasUsed=%d)", addClientReceipt.BlockNumber.Uint64(), addClientReceipt.GasUsed)
	return address, nil
}

func (h *Handler) SendEthTx(stdCtx context.Context, ctx services.Context, msg any) error {
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

	gasLimit := uint64(3000000) // in units
	if gasStr := os.Getenv("ETH_GAS_LIMIT"); gasStr != "" {
		var val uint64
		if _, err := fmt.Sscanf(gasStr, "%d", &val); err == nil {
			gasLimit = val
		}
	}

	ics07Tendermint, err := spectreContract.NewContractSpectreClient(
		*ctx.SpectreClientContract(),
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

	benchEnabled := utils.BenchEnabled()
	var benchStart time.Time
	if benchEnabled {
		benchStart = time.Now()
	}

	var txLabel string
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		switch msg := msg.(type) {
		case services.CosmosClientUpdateBuildResult:
			switch msg.Kind {
			case services.ConsensusUpdate:
				txLabel = "updateConsensusState"
				data, err := relayerclient.EncodeUpdateConsensusStateMsg(msg.AppMsg, msg.NewValSet)
				if err != nil {
					return nil, fmt.Errorf("failed to encode updateConsensusState msg: %w", err)
				}
				if routerManagesProofSubmission(ctx) {
					log.Printf("[SendEthTx] Sending ICS26Router.updateConsensusState tx for clientId=%s...", cosmosClientID)
					return ics26Router.UpdateConsensusState(auth, cosmosClientID, data)
				}
				log.Printf("[SendEthTx] Sending direct SpectreClient.updateConsensusState tx...")
				return ics07Tendermint.UpdateConsensusState(auth, data)
			default:
				txLabel = "updateApplicationState"
				data, err := relayerclient.EncodeUpdateApplicationStateMsg(msg.AppMsg)
				if err != nil {
					return nil, fmt.Errorf("failed to encode updateApplicationState msg: %w", err)
				}
				if routerManagesProofSubmission(ctx) {
					log.Printf("[SendEthTx] Sending ICS26Router.updateApplicationState tx for clientId=%s...", cosmosClientID)
					return ics26Router.UpdateApplicationState(auth, cosmosClientID, data)
				}
				log.Printf("[SendEthTx] Sending direct SpectreClient.updateApplicationState tx...")
				return ics07Tendermint.UpdateApplicationState(auth, data)
			}
		case spectreContract.ILightClientMsgsMsgVerifyMembership:
			txLabel = "verifyMembership"
			if routerManagesProofSubmission(ctx) {
				return nil, fmt.Errorf("direct verifyMembership is disabled when ROLE_MANAGER is the ICS26 router; use ICS26Router packet flows instead")
			}
			log.Printf("[SendEthTx] Sending verifyMembership tx...")
			return ics07Tendermint.VerifyMembership(auth, msg)
		case spectreContract.ILightClientMsgsMsgVerifyNonMembership:
			txLabel = "verifyNonMembership"
			if routerManagesProofSubmission(ctx) {
				return nil, fmt.Errorf("direct verifyNonMembership is disabled when ROLE_MANAGER is the ICS26 router; use ICS26Router packet flows instead")
			}
			log.Printf("[SendEthTx] Sending verifyNonMembership tx...")
			return ics07Tendermint.VerifyNonMembership(auth, msg)
		case contractICS26Router.IICS26RouterMsgsMsgRecvPacket:
			txLabel = fmt.Sprintf("recvPacket seq=%d", msg.Packet.Sequence)
			log.Printf("[SendEthTx] Sending recvPacket seq=%d...", msg.Packet.Sequence)
			return ics26Router.RecvPacket(auth, msg)
		case contractICS26Router.IICS26RouterMsgsMsgAckPacket:
			txLabel = fmt.Sprintf("ackPacket seq=%d", msg.Packet.Sequence)
			log.Printf("[SendEthTx] Sending ackPacket seq=%d...", msg.Packet.Sequence)
			return ics26Router.AckPacket(auth, msg)
		case contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket:
			txLabel = fmt.Sprintf("timeoutPacket seq=%d", msg.Packet.Sequence)
			log.Printf("[SendEthTx] Sending timeoutPacket seq=%d...", msg.Packet.Sequence)
			return ics26Router.TimeoutPacket(auth, msg)
		default:
			return nil, fmt.Errorf("unsupported message type: %T", msg)
		}
	}

	receipt, submitDur, waitDur, err := h.executeWithRetryAndResubmission(stdCtx, ctx, privateKey, gasLimit, senderFn)
	if err != nil {
		return err
	}

	log.Printf("[SendEthTx] Tx %s confirmed in block %d (gasUsed=%d)", receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), receipt.GasUsed)
	if benchEnabled {
		log.Printf("[bench][eth] %s gasUsed=%d submit=%s wait=%s total=%s tx=%s",
			txLabel, receipt.GasUsed, submitDur, waitDur, time.Since(benchStart), receipt.TxHash.Hex())
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
func (h *Handler) SendEthTxBatch(stdCtx context.Context, ctx services.Context, msgs []any) error {
	if len(msgs) == 0 {
		return nil
	}
	if len(msgs) == 1 {
		return h.SendEthTx(stdCtx, ctx, msgs[0])
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
		case services.CosmosClientUpdateBuildResult:
			// Folding a client update into a multicall only works when the
			// ICS26Router is the proof submitter — the multicall is dispatched
			// on the router, so every inner call must target a router method.
			// In direct-submission mode SendEthTx calls SpectreClient directly;
			// that path can't be expressed inside multicall, so caller must
			// submit the update as a standalone tx.
			if !routerManagesProofSubmission(ctx) {
				return fmt.Errorf("[SendEthTxBatch] client update cannot be batched when ICS26Router is not the proof submitter; submit it via SendEthTx instead")
			}
			if !clientIDResolved {
				cosmosClientID, perr = cosmosRouterClientID(ctx)
				if perr != nil {
					return fmt.Errorf("[SendEthTxBatch] resolve cosmos client id: %w", perr)
				}
				clientIDResolved = true
			}
			if m.Kind == services.ConsensusUpdate {
				encoded, encErr := relayerclient.EncodeUpdateConsensusStateMsg(m.AppMsg, m.NewValSet)
				if encErr != nil {
					return fmt.Errorf("[SendEthTxBatch] encode updateConsensusState msg %d: %w", i, encErr)
				}
				data, perr = parsedABI.Pack("updateConsensusState", cosmosClientID, encoded)
				lbl = "updateConsensusState"
			} else {
				encoded, encErr := relayerclient.EncodeUpdateApplicationStateMsg(m.AppMsg)
				if encErr != nil {
					return fmt.Errorf("[SendEthTxBatch] encode updateApplicationState msg %d: %w", i, encErr)
				}
				data, perr = parsedABI.Pack("updateApplicationState", cosmosClientID, encoded)
				lbl = "updateApplicationState"
			}
		default:
			return fmt.Errorf("[SendEthTxBatch] unsupported message type at index %d: %T (only client update/recvPacket/ackPacket/timeoutPacket allowed in multicall)", i, msg)
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
	multicallGasLimit := uint64(16000000)
	if gasStr := os.Getenv("ETH_MULTICALL_GAS_LIMIT"); gasStr != "" {
		var val uint64
		if _, err := fmt.Sscanf(gasStr, "%d", &val); err == nil {
			multicallGasLimit = val
		}
	}

	ics26Router, err := contractICS26Router.NewContractICS26Router(*ctx.RouterContract(), ctx.EthClient())
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed to bind ICS26Router: %w", err)
	}

	benchEnabled := utils.BenchEnabled()
	var benchStart time.Time
	if benchEnabled {
		benchStart = time.Now()
	}

	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		log.Printf("[SendEthTxBatch] Submitting multicall: %d inner calls (%s)", len(calldata), labelStr)
		return ics26Router.Multicall(auth, calldata)
	}

	receipt, submitDur, waitDur, err := h.executeWithRetryAndResubmission(stdCtx, ctx, privateKey, multicallGasLimit, senderFn)
	if err != nil {
		return fmt.Errorf("multicall labels=%s: %w", labelStr, err)
	}

	log.Printf("[SendEthTxBatch] Tx %s confirmed in block %d (gasUsed=%d, inner=%d)",
		receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), receipt.GasUsed, len(calldata))
	if benchEnabled {
		log.Printf("[bench][eth] multicall labels=%s gasUsed=%d submit=%s wait=%s total=%s tx=%s",
			labelStr, receipt.GasUsed, submitDur, waitDur, time.Since(benchStart), receipt.TxHash.Hex())
		if traceErr := logInnerGasFromTrace(stdCtx, ctx, receipt.TxHash, labels); traceErr != nil {
			log.Printf("[bench][eth] inner gas trace unavailable (RPC may lack debug_ namespace): %v", traceErr)
		}
	}

	return nil
}

func (h *Handler) CreateWasmClient(stdCtx context.Context, svcCtx services.Context, clientState exported.ClientState, consensusState exported.ConsensusState, counterpartyClientID string) (string, error) {
	log.Printf("[CreateWasmClientTx] starting")
	// counterpartyClientID is the client on the counterparty chain that tracks Cosmos.
	// It is a config value known upfront (registration is only a naming binding in
	// ICS26Router — the counterparty light client need not exist yet), so each caller
	// passes its own: the ETH beacon client passes the ETH-side router client id, an
	// L2 bootstrap passes the L2-side client id. An empty value skips registration.

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
	signerAddr, err := cosmosSignerBech32(privKey)
	if err != nil {
		return "", err
	}
	log.Printf("[CreateWasmClient] signer: %s", signerAddr)

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
	log.Printf("[CreateWasmClientTx] gas config: gasLimit=%d fee=%d%s", gasLimit, feeAmount, feeDenom)

	// Query account info (account number and sequence) from the chain
	log.Printf("[CreateWasmClientTx] querying cosmos account info")
	accountNumber, sequence, err := h.queryAccountInfo(stdCtx, svcCtx, signerAddr)
	if err != nil {
		return "", fmt.Errorf("failed to query account info: %w", err)
	}
	log.Printf("[CreateWasmClientTx] account info: accountNumber=%d sequence=%d", accountNumber, sequence)

	// Setup encoding config
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	channeltypesv2.RegisterInterfaces(interfaceRegistry)
	clienttypes.RegisterInterfaces(interfaceRegistry)
	ibcwasmtypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)

	msg, err := clienttypes.NewMsgCreateClient(clientState, consensusState, signerAddr)
	if err != nil {
		return "", err
	}
	log.Printf("[CreateWasmClientTx] MsgCreateClient built")

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
		Address:       signerAddr,
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	// Get sign bytes using the adapter function
	signBytes, err := authsigning.GetSignBytesAdapter(
		stdCtx,
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
	log.Printf("[CreateWasmClientTx] broadcasting MsgCreateClient")
	bctx, bcancel := context.WithTimeout(stdCtx, cosmosRPCTimeout)
	result, err := svcCtx.CosmosClient().BroadcastTxSync(bctx, txBytes)
	bcancel()
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("transaction failed with code %d: %s", result.Code, result.Log)
	}

	log.Printf("[CreateWasmClient] MsgCreateClient broadcast successfully. Hash: %s", result.Hash.String())

	// Wait for MsgCreateClient tx and extract the new client ID from events
	log.Printf("[CreateWasmClientTx] waiting for MsgCreateClient tx result")
	txResult, err := h.waitForTxResult(stdCtx, svcCtx, result.Hash, 30*time.Second)
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
	log.Printf("[CreateWasmClient] new client ID: %s", newClientID)

	// No counterparty id given (e.g. an L2 bootstrap whose L2-side client id is not yet
	// configured): the client is created, its counterparty registered later.
	if counterpartyClientID == "" {
		return newClientID, nil
	}

	// Build and broadcast MsgRegisterCounterparty as a separate transaction
	registerMsg := clienttypesv2.NewMsgRegisterCounterparty(
		newClientID,
		[][]byte{[]byte("")},
		counterpartyClientID,
		signerAddr,
	)
	log.Printf("[CreateWasmClientTx] MsgRegisterCounterparty built for clientID=%s", newClientID)

	// Re-query account info (sequence incremented after first tx)
	log.Printf("[CreateWasmClientTx] querying cosmos account info for register counterparty")
	accountNumber, sequence, err = h.queryAccountInfo(stdCtx, svcCtx, signerAddr)
	if err != nil {
		return "", fmt.Errorf("failed to query account info for register counterparty: %w", err)
	}
	log.Printf("[CreateWasmClientTx] register counterparty account info: accountNumber=%d sequence=%d", accountNumber, sequence)

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
		Address:       signerAddr,
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	signBytes2, err := authsigning.GetSignBytesAdapter(
		stdCtx,
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

	log.Printf("[CreateWasmClientTx] broadcasting MsgRegisterCounterparty")
	bctx2, bcancel2 := context.WithTimeout(stdCtx, cosmosRPCTimeout)
	result2, err := svcCtx.CosmosClient().BroadcastTxSync(bctx2, txBytes2)
	bcancel2()
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

// cosmosBech32Prefix returns the bech32 account prefix of the target Cosmos
// chain: COSMOS_ADDRESS_PREFIX, defaulting to "cosmos" (issue #220). Addresses
// are rendered with bech32.ConvertAndEncode instead of sdk.AccAddress.String()
// so the SDK's process-global bech32 config is never consulted or mutated —
// a future per-source prefix only needs to thread a string here.
func cosmosBech32Prefix() string {
	if prefix := os.Getenv("COSMOS_ADDRESS_PREFIX"); prefix != "" {
		return prefix
	}
	return "cosmos"
}

// cosmosSignerBech32 encodes the signer account address derived from privKey
// with the configured bech32 prefix.
func cosmosSignerBech32(privKey secp256k1.PrivKey) (string, error) {
	prefix := cosmosBech32Prefix()
	addr, err := sdkbech32.ConvertAndEncode(prefix, privKey.PubKey().Address())
	if err != nil {
		return "", fmt.Errorf("failed to bech32-encode signer address with prefix %q: %w", prefix, err)
	}
	return addr, nil
}

// CosmosSignerAddress returns the bech32 address derived from COSMOS_PRIVATE_KEY.
// Used to populate the Signer field in Cosmos messages before batching them.
func (h *Handler) CosmosSignerAddress() (string, error) {
	privKeyHex := strings.TrimPrefix(os.Getenv("COSMOS_PRIVATE_KEY"), "0x")
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode COSMOS_PRIVATE_KEY: %w", err)
	}
	return cosmosSignerBech32(secp256k1.PrivKey{Key: privKeyBytes})
}

func cloneCosmosSDKMsgWithSigner(msg any, signer string, index int) (sdk.Msg, error) {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		if index >= 0 {
			return nil, fmt.Errorf("message %d must be a proto.Message", index)
		}
		return nil, fmt.Errorf("message must be a proto.Message")
	}
	clonedProto := proto.Clone(protoMsg)
	if clonedProto == nil {
		if index >= 0 {
			return nil, fmt.Errorf("message %d cloned to nil", index)
		}
		return nil, fmt.Errorf("message cloned to nil")
	}
	sdkMsg, ok := clonedProto.(sdk.Msg)
	if !ok {
		if index >= 0 {
			return nil, fmt.Errorf("message %d does not implement sdk.Msg interface", index)
		}
		return nil, fmt.Errorf("message does not implement sdk.Msg interface")
	}
	fillEmptyCosmosSigner(sdkMsg, signer)
	return sdkMsg, nil
}

func fillEmptyCosmosSigner(msg sdk.Msg, signer string) {
	switch m := msg.(type) {
	case *clienttypes.MsgUpdateClient:
		if m.Signer == "" {
			m.Signer = signer
		}
	case *channeltypesv2.MsgRecvPacket:
		if m.Signer == "" {
			m.Signer = signer
		}
	case *channeltypesv2.MsgAcknowledgement:
		if m.Signer == "" {
			m.Signer = signer
		}
	case *channeltypesv2.MsgTimeout:
		if m.Signer == "" {
			m.Signer = signer
		}
	}
}

// isCosmosDuplicatePacketError returns true when the ABCI response indicates the
// packet(s) were already processed, so retrying the identical tx can never
// succeed and it is safe to drop. Matches two deterministic duplicate signals:
//
//	codespace "channelv2", code 11 = ErrAcknowledgementExists (ack already written)
//	codespace "channelv2", code 12 = ErrNoOpMsg (duplicate-delivery signal)
//	codespace "channel",   code 22 = ErrRedundantTx — the RedundantRelayDecorator
//	  ante handler rejects a tx whose packet messages are ALL already processed
//	  (raised at CheckTx). Without dropping this the relayer re-queues the tx and
//	  re-submits it forever, each retry re-running a full client update — draining
//	  gas — even though it can never stop being redundant.
//
// Substring matching on log messages is avoided because "commitment not found" can
// mean a packet was never sent — not just a duplicate — and would silently drop it.
func isCosmosDuplicatePacketError(codespace string, code uint32) bool {
	return (codespace == "channelv2" && (code == 11 || code == 12)) ||
		(codespace == "channel" && code == 22)
}

// simulateMsgs builds a transaction with the given messages, signs it with an empty signature, and simulates its gas consumption.
func (h *Handler) simulateMsgs(stdCtx context.Context, svcCtx services.Context, sdkMsgs []sdk.Msg, sequence uint64) (uint64, error) {
	// Setup encoding config
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	channeltypesv2.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)

	// Get key details to construct empty signature
	privKeyHex := os.Getenv("COSMOS_PRIVATE_KEY")
	if privKeyHex == "" {
		return 0, fmt.Errorf("COSMOS_PRIVATE_KEY environment variable is required")
	}
	privKeyBytes, err := hex.DecodeString(strings.TrimPrefix(privKeyHex, "0x"))
	if err != nil {
		return 0, fmt.Errorf("failed to decode private key: %w", err)
	}
	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	pubKey := privKey.PubKey()

	// Build the transaction
	txBuilder := txConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(sdkMsgs...); err != nil {
		return 0, fmt.Errorf("failed to set messages: %w", err)
	}

	// We set a reasonable gas limit and empty fee for simulation
	txBuilder.SetGasLimit(20000000) // generous gas limit for simulation
	feeDenom := os.Getenv("COSMOS_FEE_DENOM")
	if feeDenom == "" {
		feeDenom = "stake"
	}
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, sdkmath.NewInt(0))))

	// Set empty signature
	emptySig := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	if err := txBuilder.SetSignatures(emptySig); err != nil {
		return 0, fmt.Errorf("failed to set empty signature: %w", err)
	}

	// Encode the transaction
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return 0, fmt.Errorf("failed to encode transaction: %w", err)
	}

	// Create and marshal simulation request
	simReq := &txservice.SimulateRequest{
		TxBytes: txBytes,
	}
	simReqBytes, err := proto.Marshal(simReq)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal simulate request: %w", err)
	}

	// Query simulate endpoint
	qctx, qcancel := context.WithTimeout(stdCtx, cosmosRPCTimeout)
	result, err := svcCtx.CosmosClient().ABCIQuery(qctx, "/cosmos.tx.v1beta1.Service/Simulate", simReqBytes)
	qcancel()
	if err != nil {
		return 0, fmt.Errorf("simulate ABCI query failed: %w", err)
	}

	if result.Response.Code != 0 {
		return 0, fmt.Errorf("simulation failed with code %d: %s", result.Response.Code, result.Response.Log)
	}

	var simResp txservice.SimulateResponse
	if err := proto.Unmarshal(result.Response.Value, &simResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal simulate response: %w", err)
	}

	if simResp.GasInfo == nil {
		return 0, fmt.Errorf("gas info is nil in simulation response")
	}

	return simResp.GasInfo.GasUsed, nil
}

// sendCosmosTxBatchWithSplitting handles gas simulation, clamping, fee scaling, and recursive batch splitting.
func (h *Handler) sendCosmosTxBatchWithSplitting(stdCtx context.Context, svcCtx services.Context, sdkMsgs []sdk.Msg, accountNumber, sequence uint64) (uint64, int, error) {
	if len(sdkMsgs) == 0 {
		return sequence, 0, nil
	}

	// Get block gas limit
	var maxBlockGas uint64 = 0
	if params, err := svcCtx.CosmosClient().ConsensusParams(stdCtx, nil); err == nil && params != nil {
		if params.ConsensusParams.Block.MaxGas > 0 {
			maxBlockGas = uint64(params.ConsensusParams.Block.MaxGas)
		}
	}

	shouldSplit := false
	var finalGasLimit uint64

	// Simulate gas consumption for the messages in the batch.
	if len(sdkMsgs) > 0 {
		simulatedGas, err := h.simulateMsgs(stdCtx, svcCtx, sdkMsgs, sequence)
		if err != nil {
			log.Printf("[SendCosmosTxBatch] Simulation failed for batch of size %d: %v", len(sdkMsgs), err)
			if len(sdkMsgs) > 1 {
				log.Printf("[SendCosmosTxBatch] Splitting batch...")
				shouldSplit = true
			}
		} else {
			// Apply a 1.3 gas adjustment factor
			adjustedGas := uint64(float64(simulatedGas) * 1.3)
			if maxBlockGas > 0 && adjustedGas >= maxBlockGas {
				log.Printf("[SendCosmosTxBatch] Adjusted gas %d exceeds max block gas %d for batch of size %d", adjustedGas, maxBlockGas, len(sdkMsgs))
				if len(sdkMsgs) > 1 {
					log.Printf("[SendCosmosTxBatch] Splitting batch...")
					shouldSplit = true
				} else {
					// Single message exceeds block limit; clamp it to max block gas.
					finalGasLimit = maxBlockGas
				}
			} else {
				finalGasLimit = adjustedGas
			}
		}
	}

	if shouldSplit {
		mid := len(sdkMsgs) / 2
		nextSeq, succ1, err := h.sendCosmosTxBatchWithSplitting(stdCtx, svcCtx, sdkMsgs[:mid], accountNumber, sequence)
		if err != nil {
			return sequence, succ1, err
		}
		nextSeq, succ2, err := h.sendCosmosTxBatchWithSplitting(stdCtx, svcCtx, sdkMsgs[mid:], accountNumber, nextSeq)
		if err != nil {
			return nextSeq, succ1 + succ2, err
		}
		return nextSeq, succ1 + succ2, nil
	}

	// If simulation wasn't run or failed, calculate the fallback gas limit
	if finalGasLimit == 0 {
		baseGas := uint64(200000)
		if gasStr := os.Getenv("COSMOS_GAS_LIMIT"); gasStr != "" {
			if _, err := fmt.Sscanf(gasStr, "%d", &baseGas); err != nil {
				return sequence, 0, fmt.Errorf("failed to parse COSMOS_GAS_LIMIT: %w", err)
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

		// Guard against gasLimit overflow: baseGas * len(sdkMsgs)
		if len(sdkMsgs) > 0 && baseGas > (1<<64-1)/uint64(len(sdkMsgs)) {
			finalGasLimit = 1<<64 - 1
		} else {
			finalGasLimit = baseGas * uint64(len(sdkMsgs))
		}
	}

	// Clamp to block gas limit and guard against zero limit
	if maxBlockGas > 0 && finalGasLimit > maxBlockGas {
		finalGasLimit = maxBlockGas
	}

	feeDenom := os.Getenv("COSMOS_FEE_DENOM")
	if feeDenom == "" {
		feeDenom = "stake"
	}

	// Set fee amount: default to matching gas limit
	feeAmount := int64(finalGasLimit)
	if feeStr := os.Getenv("COSMOS_FEE_AMOUNT"); feeStr != "" {
		var baseFee int64
		if _, err := fmt.Sscanf(feeStr, "%d", &baseFee); err != nil {
			return sequence, 0, fmt.Errorf("failed to parse COSMOS_FEE_AMOUNT: %w", err)
		}

		// Guard against feeAmount overflow: baseFee * len(sdkMsgs)
		if len(sdkMsgs) > 0 && baseFee > (1<<63-1)/int64(len(sdkMsgs)) {
			feeAmount = 1<<63 - 1
		} else {
			feeAmount = baseFee * int64(len(sdkMsgs))
		}
	}

	// Clamp feeAmount if gasLimit was clamped to block limit
	if maxBlockGas > 0 && feeAmount > int64(maxBlockGas) {
		feeAmount = int64(maxBlockGas)
	}

	// Now build, sign, and broadcast the transaction!
	privKeyHex := os.Getenv("COSMOS_PRIVATE_KEY")
	if privKeyHex == "" {
		return sequence, 0, fmt.Errorf("COSMOS_PRIVATE_KEY environment variable is required")
	}
	privKeyBytes, err := hex.DecodeString(strings.TrimPrefix(privKeyHex, "0x"))
	if err != nil {
		return sequence, 0, fmt.Errorf("failed to decode private key: %w", err)
	}
	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	pubKey := privKey.PubKey()
	signerAddr, err := cosmosSignerBech32(privKey)
	if err != nil {
		return sequence, 0, err
	}
	chainID := os.Getenv("COSMOS_CHAIN_ID")
	if chainID == "" {
		return sequence, 0, fmt.Errorf("COSMOS_CHAIN_ID environment variable is required")
	}

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	channeltypesv2.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)

	txBuilder := txConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(sdkMsgs...); err != nil {
		return sequence, 0, fmt.Errorf("failed to set messages: %w", err)
	}

	txBuilder.SetGasLimit(finalGasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, sdkmath.NewInt(feeAmount))))

	emptySig := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	if err := txBuilder.SetSignatures(emptySig); err != nil {
		return sequence, 0, fmt.Errorf("failed to set empty signature: %w", err)
	}

	signerData := authsigning.SignerData{
		Address:       signerAddr,
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	signBytes, err := authsigning.GetSignBytesAdapter(
		stdCtx,
		txConfig.SignModeHandler(),
		sdksigning.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		return sequence, 0, fmt.Errorf("failed to get sign bytes: %w", err)
	}

	sigRaw, err := privKey.Sign(signBytes)
	if err != nil {
		return sequence, 0, fmt.Errorf("failed to sign transaction: %w", err)
	}

	sigV2 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: sigRaw,
		},
		Sequence: sequence,
	}

	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return sequence, 0, fmt.Errorf("failed to set signatures: %w", err)
	}

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return sequence, 0, fmt.Errorf("failed to encode transaction: %w", err)
	}

	benchEnabled := utils.BenchEnabled()
	var broadcastStart time.Time
	if benchEnabled {
		broadcastStart = time.Now()
	}

	bctx, bcancel := context.WithTimeout(stdCtx, cosmosRPCTimeout)
	syncResult, err := svcCtx.CosmosClient().BroadcastTxSync(bctx, txBytes)
	bcancel()
	if err != nil {
		return sequence, 0, fmt.Errorf("failed to broadcast transaction: %w", err)
	}
	if syncResult.Code != 0 {
		log.Printf("[SendCosmosTxBatch] CheckTx FAILED: code=%d codespace=%s log=%s data=%x",
			syncResult.Code, syncResult.Codespace, syncResult.Log, syncResult.Data)
		if isCosmosDuplicatePacketError(syncResult.Codespace, syncResult.Code) {
			log.Printf("[SendCosmosTxBatch] duplicate packet (codespace=%s code=%d), dropping", syncResult.Codespace, syncResult.Code)
			return sequence, len(sdkMsgs), nil
		}
		return sequence, 0, &services.CosmosTxFailure{
			Stage:     "CheckTx",
			Code:      syncResult.Code,
			Codespace: syncResult.Codespace,
			Log:       syncResult.Log,
			Data:      syncResult.Data,
		}
	}

	txResult, err := h.waitForTxResult(stdCtx, svcCtx, syncResult.Hash, cosmosInclusionTimeout)
	if err != nil {
		return sequence, 0, fmt.Errorf("failed to confirm transaction inclusion: %w", err)
	}

	var broadcastDur time.Duration
	if benchEnabled {
		broadcastDur = time.Since(broadcastStart)
	}

	if txResult.TxResult.Code != 0 {
		log.Printf("[SendCosmosTxBatch] DeliverTx FAILED: code=%d codespace=%s log=%s data=%x",
			txResult.TxResult.Code, txResult.TxResult.Codespace, txResult.TxResult.Log, txResult.TxResult.Data)
		if isCosmosDuplicatePacketError(txResult.TxResult.Codespace, txResult.TxResult.Code) {
			log.Printf("[SendCosmosTxBatch] duplicate packet (codespace=%s code=%d), dropping", txResult.TxResult.Codespace, txResult.TxResult.Code)
			return sequence + 1, len(sdkMsgs), nil
		}
		return sequence, 0, &services.CosmosTxFailure{
			Stage:     "DeliverTx",
			Code:      txResult.TxResult.Code,
			Codespace: txResult.TxResult.Codespace,
			Log:       txResult.TxResult.Log,
			Data:      txResult.TxResult.Data,
			Err:       services.ErrPermanentRelayFailure,
		}
	}

	log.Printf("[SendCosmosTxBatch] Tx confirmed at height %d hash=%s (msgs=%d)", txResult.Height, txResult.Hash.String(), len(sdkMsgs))
	if benchEnabled {
		log.Printf("[bench][cosmos] batch msgs=%d gasWanted=%d gasUsed=%d broadcast=%s height=%d hash=%s",
			len(sdkMsgs), txResult.TxResult.GasWanted, txResult.TxResult.GasUsed,
			broadcastDur, txResult.Height, txResult.Hash.String())
	}

	return sequence + 1, len(sdkMsgs), nil
}

// SendCosmosTxBatch sends multiple messages in a single Cosmos transaction.
//
// NOTE on Non-Atomicity:
// If a batch is large or simulation indicates it would exceed the block gas limit,
// the batch is recursively split into smaller sub-batches and sent as MULTIPLE separate transactions.
// These transactions are executed sequentially, and the account sequence is updated accordingly.
// If an earlier sub-batch succeeds but a later one fails, an error is returned.
//
// Retry Behavior and Idempotency:
// Callers should be aware that the overall batch operation is not atomic. On failure,
// the relayer's main loop will retry. However, because the loop is state-based, it queries
// the current chain state at the start of each iteration. Any messages/packets successfully
// committed by the earlier succeeded sub-batches will not be included in the retried batch.
// Downstream Cosmos modules/contracts are idempotent and tolerate already-processed packets safely.
func (h *Handler) SendCosmosTxBatch(stdCtx context.Context, svcCtx services.Context, msgs []any) error {
	if len(msgs) == 0 {
		return nil
	}

	benchEnabled := utils.BenchEnabled()
	var benchStart time.Time
	if benchEnabled {
		benchStart = time.Now()
	}

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
	signerAddr, err := cosmosSignerBech32(privKey)
	if err != nil {
		return err
	}

	// Convert all messages to sdk.Msg, filling empty Signer fields
	var sdkMsgs []sdk.Msg
	for i, msg := range msgs {
		sdkMsg, err := cloneCosmosSDKMsgWithSigner(msg, signerAddr, i)
		if err != nil {
			return err
		}
		sdkMsgs = append(sdkMsgs, sdkMsg)
	}

	// Serialize Cosmos sequence use across relay goroutines. A split batch may
	// submit multiple txs, so the lock covers the whole sequence chain.
	h.cosmosMu.Lock()
	defer h.cosmosMu.Unlock()

	// Query account info (account number and sequence) from the chain
	accountNumber, sequence, err := h.queryAccountInfo(stdCtx, svcCtx, signerAddr)
	if err != nil {
		return fmt.Errorf("failed to query account info: %w", err)
	}

	_, succCount, err := h.sendCosmosTxBatchWithSplitting(stdCtx, svcCtx, sdkMsgs, accountNumber, sequence)
	if benchEnabled {
		log.Printf("[bench][cosmos] batch msgs=%d total=%s", len(msgs), time.Since(benchStart))
	}
	if err != nil {
		if succCount > 0 {
			return &services.BatchPartialError{
				SucceededCount: succCount,
				Err:            err,
			}
		}
		return err
	}
	return nil
}

// queryAccountInfo queries the account number and sequence for the given address
func (h *Handler) queryAccountInfo(stdCtx context.Context, svcCtx services.Context, address string) (uint64, uint64, error) {
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
	qctx, qcancel := context.WithTimeout(stdCtx, cosmosRPCTimeout)
	result, err := svcCtx.CosmosClient().ABCIQuery(qctx, queryPath, reqBytes)
	qcancel()
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
func (h *Handler) waitForTxResult(stdCtx context.Context, svcCtx services.Context, txHash []byte, timeout time.Duration) (*coretypes.ResultTx, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		// Abort promptly on shutdown. Without this the poll loop runs until the local
		// timeout (up to cosmosInclusionTimeout) — and because SendCosmosTxBatch holds
		// cosmosMu across this call, that would stall shutdown for the whole window.
		if err := stdCtx.Err(); err != nil {
			return nil, fmt.Errorf("cancelled waiting for tx %X: %w", txHash, err)
		}
		qctx, qcancel := context.WithTimeout(stdCtx, cosmosRPCTimeout)
		result, err := svcCtx.CosmosClient().Tx(qctx, txHash, false)
		qcancel()
		if err == nil && result != nil && result.Height > 0 {
			log.Printf("[WaitForTx] Tx %X confirmed at height %d", txHash, result.Height)
			return result, nil
		}
		// Wait before the next poll, but cancel out immediately instead of sleeping a
		// full second past a shutdown signal.
		timer := time.NewTimer(1 * time.Second)
		select {
		case <-stdCtx.Done():
			timer.Stop()
			return nil, fmt.Errorf("cancelled waiting for tx %X: %w", txHash, stdCtx.Err())
		case <-timer.C:
		}
	}
	return nil, fmt.Errorf("timeout waiting for tx %X to be included in a block", txHash)
}

func isNonceTooLowError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "nonce too low") ||
		strings.Contains(msg, "old nonce") ||
		strings.Contains(msg, "nonce has already been used")
}

func isAlreadyKnownError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "already known") ||
		strings.Contains(msg, "transaction already imported") ||
		strings.Contains(msg, "already exists")
}

func waitForReceipts(ctx context.Context, client *ethclient.Client, hashes []common.Hash) (*types.Receipt, error) {
	ticker := time.NewTicker(ethTxReceiptPollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			for i := len(hashes) - 1; i >= 0; i-- {
				hash := hashes[i]
				receipt, err := client.TransactionReceipt(ctx, hash)
				if err == nil {
					if receipt != nil {
						return receipt, nil
					}
					// Treat nil receipt as not found
					continue
				}
				if err != nil && !errors.Is(err, ethereum.NotFound) {
					log.Printf("[EthTxSender] Transient error polling receipt for %s: %v. Retrying...", hash.Hex(), err)
					continue
				}
			}
		}
	}
}

func (h *Handler) bumpGasAndResubmit(stdCtx context.Context, ctx services.Context, tx *types.Transaction, auth *bind.TransactOpts, attempt int) (*types.Transaction, error) {
	currentGasPrice := tx.GasPrice()

	suggestedGasPrice, err := ctx.EthClient().SuggestGasPrice(stdCtx)
	if err != nil {
		suggestedGasPrice = currentGasPrice
	}

	bumpedGasPrice := new(big.Int).Mul(currentGasPrice, big.NewInt(115))
	bumpedGasPrice.Div(bumpedGasPrice, big.NewInt(100))

	if suggestedGasPrice.Cmp(bumpedGasPrice) > 0 {
		bumpedGasPrice = suggestedGasPrice
	}

	var newTx *types.Transaction
	if tx.Type() == types.DynamicFeeTxType {
		tipCap := tx.GasTipCap()
		feeCap := tx.GasFeeCap()

		bumpedTipCap := new(big.Int).Mul(tipCap, big.NewInt(115))
		bumpedTipCap.Div(bumpedTipCap, big.NewInt(100))

		bumpedFeeCap := new(big.Int).Mul(feeCap, big.NewInt(115))
		bumpedFeeCap.Div(bumpedFeeCap, big.NewInt(100))

		if suggestedTipCap, err := ctx.EthClient().SuggestGasTipCap(stdCtx); err == nil {
			if suggestedTipCap.Cmp(bumpedTipCap) > 0 {
				bumpedTipCap = suggestedTipCap
			}
		}
		if suggestedGasPrice.Cmp(bumpedFeeCap) > 0 {
			bumpedFeeCap = suggestedGasPrice
		}

		newTx = types.NewTx(&types.DynamicFeeTx{
			ChainID:    tx.ChainId(),
			Nonce:      tx.Nonce(),
			GasTipCap:  bumpedTipCap,
			GasFeeCap:  bumpedFeeCap,
			Gas:        tx.Gas(),
			To:         tx.To(),
			Value:      tx.Value(),
			Data:       tx.Data(),
			AccessList: tx.AccessList(),
		})
	} else if tx.Type() == types.AccessListTxType {
		newTx = types.NewTx(&types.AccessListTx{
			ChainID:    tx.ChainId(),
			Nonce:      tx.Nonce(),
			GasPrice:   bumpedGasPrice,
			Gas:        tx.Gas(),
			To:         tx.To(),
			Value:      tx.Value(),
			Data:       tx.Data(),
			AccessList: tx.AccessList(),
		})
	} else {
		newTx = types.NewTx(&types.LegacyTx{
			Nonce:    tx.Nonce(),
			GasPrice: bumpedGasPrice,
			Gas:      tx.Gas(),
			To:       tx.To(),
			Value:    tx.Value(),
			Data:     tx.Data(),
		})
	}

	signedTx, err := auth.Signer(auth.From, newTx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign bumped transaction: %w", err)
	}

	sendCtx, cancel := context.WithTimeout(stdCtx, ethTxBroadcastTimeout)
	defer cancel()
	sendErr := ctx.EthClient().SendTransaction(sendCtx, signedTx)
	if sendErr != nil {
		return nil, fmt.Errorf("failed to send bumped transaction: %w", sendErr)
	}

	return signedTx, nil
}

func (h *Handler) executeWithRetryAndResubmission(
	stdCtx context.Context,
	ctx services.Context,
	privateKey *ecdsa.PrivateKey,
	gasLimit uint64,
	senderFn func(auth *bind.TransactOpts) (*types.Transaction, error),
) (*types.Receipt, time.Duration, time.Duration, error) {
	publicKey, err := keys.PublicKey(privateKey)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to derive public key: %w", err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	chainIdInt, err := ctx.EthClient().ChainID(stdCtx)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("invalid chain id: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIdInt)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to create auth transactor: %w", err)
	}
	auth.Value = big.NewInt(0)
	auth.GasLimit = gasLimit

	maxNonceRetries := 3
	var tx *types.Transaction
	var signedTx *types.Transaction

	// Capture the signed transaction
	originalSigner := auth.Signer
	auth.Signer = func(address common.Address, txToSign *types.Transaction) (*types.Transaction, error) {
		sTx, err := originalSigner(address, txToSign)
		if err == nil {
			signedTx = sTx
		}
		return sTx, err
	}
	defer func() {
		auth.Signer = originalSigner
	}()

	submitStart := time.Now()
	var submitDur time.Duration

	// Detect EIP-1559 support and fetch base fee once before the retry loop
	var isEIP1559 bool
	var suggestedTip *big.Int
	var baseFee *big.Int
	if tip, err := ctx.EthClient().SuggestGasTipCap(stdCtx); err == nil {
		isEIP1559 = true
		suggestedTip = tip
		if header, err := ctx.EthClient().HeaderByNumber(stdCtx, nil); err == nil && header.BaseFee != nil {
			baseFee = header.BaseFee
		} else {
			baseFee = big.NewInt(1000000000) // fallback 1 Gwei
		}
	}

	for nonceAttempt := 1; nonceAttempt <= maxNonceRetries; nonceAttempt++ {
		h.mu.Lock()
		if !h.nonceValid {
			n, err := ctx.EthClient().PendingNonceAt(stdCtx, fromAddress)
			if err != nil {
				h.mu.Unlock()
				return nil, 0, 0, fmt.Errorf("failed to get pending nonce: %w", err)
			}
			h.nonce = n
			h.nonceValid = true
		}
		currentNonce := h.nonce

		if isEIP1559 {
			gasTipCap := suggestedTip
			gasFeeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), gasTipCap)

			// Enforce minimum floor if this nonce matches the last attempted nonce (e.g. replacing a stuck tx)
			if currentNonce == h.lastNonce {
				if h.lastGasTipCap != nil {
					minTip := new(big.Int).Mul(h.lastGasTipCap, big.NewInt(115))
					minTip.Div(minTip, big.NewInt(100))
					if gasTipCap.Cmp(minTip) < 0 {
						gasTipCap = minTip
					}
				}
				if h.lastGasFeeCap != nil {
					minFee := new(big.Int).Mul(h.lastGasFeeCap, big.NewInt(115))
					minFee.Div(minFee, big.NewInt(100))
					if gasFeeCap.Cmp(minFee) < 0 {
						gasFeeCap = minFee
					}
				}
			}

			auth.GasTipCap = gasTipCap
			auth.GasFeeCap = gasFeeCap
			auth.GasPrice = nil
		} else {
			gasPrice, err := ctx.EthClient().SuggestGasPrice(stdCtx)
			if err != nil {
				h.nonceValid = false
				h.mu.Unlock()
				return nil, 0, 0, fmt.Errorf("failed to suggest gas price: %w", err)
			}

			// Enforce minimum floor if this nonce matches the last attempted nonce (e.g. replacing a stuck tx)
			if currentNonce == h.lastNonce {
				if h.lastGasPrice != nil {
					minPrice := new(big.Int).Mul(h.lastGasPrice, big.NewInt(115))
					minPrice.Div(minPrice, big.NewInt(100))
					if gasPrice.Cmp(minPrice) < 0 {
						gasPrice = minPrice
					}
				}
			}

			auth.GasPrice = gasPrice
			auth.GasTipCap = nil
			auth.GasFeeCap = nil
		}

		auth.Nonce = big.NewInt(int64(currentNonce))
		signedTx = nil

		var callErr error
		sendCtx, sendCancel := context.WithTimeout(stdCtx, ethTxBroadcastTimeout)
		auth.Context = sendCtx
		tx, callErr = senderFn(auth)
		sendCancel()
		auth.Context = nil

		if callErr == nil {
			if signedTx == nil {
				// senderFn returned (nil, nil) without calling auth.Signer — treat as a bug
				h.nonceValid = false
				h.mu.Unlock()
				return nil, 0, 0, fmt.Errorf("senderFn returned nil transaction without error")
			}
			tx = signedTx
			submitDur = time.Since(submitStart)

			h.nonce++
			h.lastNonce = tx.Nonce()
			if tx.Type() == types.DynamicFeeTxType {
				h.lastGasFeeCap = tx.GasFeeCap()
				h.lastGasTipCap = tx.GasTipCap()
				h.lastGasPrice = nil
			} else {
				h.lastGasPrice = tx.GasPrice()
				h.lastGasFeeCap = nil
				h.lastGasTipCap = nil
			}
			h.mu.Unlock()

			break
		}

		if isNonceTooLowError(callErr) {
			log.Printf("[EthTxSender] Nonce %d too low (attempt %d/%d). Resetting nonce cache.", currentNonce, nonceAttempt, maxNonceRetries)
			h.nonceValid = false
			h.mu.Unlock()
			continue
		}

		if isAlreadyKnownError(callErr) && signedTx != nil {
			log.Printf("[EthTxSender] Transaction already known in mempool: %s. Proceeding to wait.", signedTx.Hash().Hex())
			tx = signedTx
			submitDur = time.Since(submitStart)

			h.nonce++
			h.lastNonce = tx.Nonce()
			if tx.Type() == types.DynamicFeeTxType {
				h.lastGasFeeCap = tx.GasFeeCap()
				h.lastGasTipCap = tx.GasTipCap()
				h.lastGasPrice = nil
			} else {
				h.lastGasPrice = tx.GasPrice()
				h.lastGasFeeCap = nil
				h.lastGasTipCap = nil
			}
			h.mu.Unlock()

			break
		}

		// Other error: invalidate nonce just in case and return
		h.nonceValid = false
		h.mu.Unlock()
		return nil, 0, 0, fmt.Errorf("contract call failed: %w", callErr)
	}

	if tx == nil {
		return nil, 0, 0, fmt.Errorf("failed to build or send transaction after nonce retries")
	}

	sentHashes := []common.Hash{tx.Hash()}
	maxAttempts := 5
	baseTimeout := ethTxReceiptTimeout

	waitStart := time.Now()
	attempt := 1
	for {
		attemptTimeout := baseTimeout * (1 << (attempt - 1))
		if attemptTimeout > 5*time.Minute {
			attemptTimeout = 5 * time.Minute
		}

		receiptCtx, cancel := context.WithTimeout(stdCtx, attemptTimeout)
		receipt, waitErr := waitForReceipts(receiptCtx, ctx.EthClient(), sentHashes)
		cancel()

		if waitErr == nil {
			waitDur := time.Since(waitStart)
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
				_, callErr := ctx.EthClient().CallContract(stdCtx, callMsg, receipt.BlockNumber)
				if callErr != nil {
					log.Printf("[EthTxSender] Revert reason: %v", callErr)
					type dataErr interface {
						ErrorData() interface{}
					}
					if de, ok := callErr.(dataErr); ok {
						log.Printf("[EthTxSender] Revert data (hex): %v", de.ErrorData())
					}
					if name, ok := validatorCacheRaceErrorName(callErr); ok {
						return receipt, submitDur, waitDur, fmt.Errorf("tx %s reverted with %s (status=0, gasUsed=%d): %w",
							receipt.TxHash.Hex(), name, receipt.GasUsed, services.ErrValidatorCacheRace)
					}
				}
				return receipt, submitDur, waitDur, fmt.Errorf("tx %s reverted (status=0, gasUsed=%d): %w", receipt.TxHash.Hex(), receipt.GasUsed, services.ErrPermanentRelayFailure)
			}
			return receipt, submitDur, waitDur, nil
		}

		if errors.Is(waitErr, context.DeadlineExceeded) {
			if attempt >= maxAttempts {
				h.mu.Lock()
				h.nonceValid = false
				h.mu.Unlock()
				return nil, submitDur, time.Since(waitStart), fmt.Errorf("transaction wait mined timed out after %d attempts (last hash: %s): %w", attempt, tx.Hash().Hex(), waitErr)
			}

			log.Printf("[EthTxSender] Tx %s not mined in %s, bumping gas price...", tx.Hash().Hex(), attemptTimeout)
			bumpedTx, bumpErr := h.bumpGasAndResubmit(stdCtx, ctx, tx, auth, attempt)
			if bumpErr != nil {
				log.Printf("[EthTxSender] Gas bump attempt %d failed: %v. Will continue waiting.", attempt, bumpErr)
			} else {
				tx = bumpedTx
				sentHashes = append(sentHashes, tx.Hash())
				log.Printf("[EthTxSender] Gas bumped tx submitted: %s (attempt %d)", tx.Hash().Hex(), attempt+1)

				h.mu.Lock()
				h.lastNonce = tx.Nonce()
				if tx.Type() == types.DynamicFeeTxType {
					h.lastGasFeeCap = tx.GasFeeCap()
					h.lastGasTipCap = tx.GasTipCap()
					h.lastGasPrice = nil
				} else {
					h.lastGasPrice = tx.GasPrice()
					h.lastGasFeeCap = nil
					h.lastGasTipCap = nil
				}
				h.mu.Unlock()
			}
			attempt++
			continue
		}

		h.mu.Lock()
		h.nonceValid = false
		h.mu.Unlock()
		return nil, submitDur, time.Since(waitStart), fmt.Errorf("failed waiting for tx receipt: %w", waitErr)
	}
}
