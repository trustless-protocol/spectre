// This file is the EVM half of the transaction package: every send goes through
// the nonce block under h.mu, for BOTH the L1 and the L2s.
//
// It is a file rather than a section because the package's two halves are
// mirrors, and non-functional requirement #5 says every change must be checked
// against the opposite path. While both lived in one 2501-line handler.go the
// two halves sat about a thousand lines apart, so that check was a scrolling
// exercise instead of a diff between two files -- and every asymmetry bug found
// in this package so far appeared under exactly that condition.
//
// It stays over the 800-line limit on purpose. Cutting further (ethereum_gas.go,
// and so on) would buy a nicer number and break the thing the split is for: the
// mirror comparison would become four files instead of two. The line limit
// serves the symmetry rule, not the other way round.
//
// The shared declarations -- Handler, the signer accessors, the gas limits, the
// Solidity error selectors, the call-trace decoding -- live here because they are
// EVM-shaped and the target layout has no handler.go for them to sit in.
package transaction

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	routerContract "relayer/bindings/ICS26Router"
	spectreContract "relayer/bindings/SpectreClient"
	"relayer/chain"
	relayerclient "relayer/client"
	"relayer/keys"
	services "relayer/services"
	utils "relayer/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

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
	mu           sync.Mutex
	cosmosMu     sync.Mutex
	senderStates map[evmNonceKey]*ethTxSenderState

	signerOnce sync.Once
	signer     Signer
}

// SetSigner installs the ordinary signing-key source. It must be called before
// the first send; a nil value is ignored so it cannot consume signerOnce and
// turn the next send into a nil-pointer panic.
func (h *Handler) SetSigner(s Signer) {
	if s == nil {
		log.Printf("[Handler] SetSigner(nil) ignored; keeping the existing key source")
		return
	}
	h.signer = s
}

func (h *Handler) keySigner() Signer {
	h.signerOnce.Do(func() {
		if h.signer == nil {
			h.signer = NewEnvSigner()
		}
	})
	return h.signer
}

// ValidateKeys forces both ordinary signing keys through the same seam used by
// send paths and derives both signer addresses, so startup also validates the
// configured bech32 prefix.
func (h *Handler) ValidateKeys() error {
	if _, err := h.EthSignerAddress(); err != nil {
		return err
	}
	_, err := h.CosmosSignerAddress()
	return err
}

// EthSignerAddress returns the EVM address derived from ETH_PRIVATE_KEY.
//
// It belongs beside CosmosSignerAddress rather than at a submit call site: a
// caller that needs to identify the EVM nonce domain before sending must derive
// exactly the same address SendEthTx will use, through the same signer seam.
func (h *Handler) EthSignerAddress() (common.Address, error) {
	privateKey, err := h.keySigner().EthKey()
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to restore ETH private key: %w", err)
	}
	publicKey, err := keys.PublicKey(privateKey)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to derive ETH public key: %w", err)
	}
	return crypto.PubkeyToAddress(*publicKey), nil
}

// evmNonceKey identifies one independent EVM nonce domain. The relayer shares
// one Handler across chains and the incident command uses a separate signer.
type evmNonceKey struct {
	chainID string
	address common.Address
}

// ethTxSenderState stores nonce and replacement-fee state for one chain/account
// pair. Callers must hold Handler.mu while accessing it.
type ethTxSenderState struct {
	nonce         uint64
	nonceValid    bool
	lastNonce     uint64
	lastGasPrice  *big.Int
	lastGasFeeCap *big.Int
	lastGasTipCap *big.Int
}

func (h *Handler) senderState(chainID string, sender common.Address) *ethTxSenderState {
	if h.senderStates == nil {
		h.senderStates = make(map[evmNonceKey]*ethTxSenderState)
	}
	key := evmNonceKey{chainID: chainID, address: sender}
	state := h.senderStates[key]
	if state == nil {
		state = &ethTxSenderState{}
		h.senderStates[key] = state
	}
	return state
}

var ethTxBroadcastTimeout = 30 * time.Second
var ethTxReceiptTimeout = 45 * time.Second
var ethTxReceiptPollInterval = 2 * time.Second

const (
	defaultEthGasLimit          uint64 = 3_000_000
	defaultMisbehaviourGasLimit uint64 = 16_000_000
)

// EthGasLimit returns the regular EVM transaction gas limit. A zero override
// would make go-ethereum estimate ordinary sends but leave the OOG retry ladder
// with a zero base, so retain the safe default instead.
func EthGasLimit(gasText string) (uint64, error) {
	if gasText == "" {
		return defaultEthGasLimit, nil
	}
	gasLimit, err := strconv.ParseUint(gasText, 10, 64)
	if err != nil || gasLimit == 0 {
		return 0, fmt.Errorf("invalid ETH_GAS_LIMIT %q", gasText)
	}
	return gasLimit, nil
}

// envUint64 reads an optional unsigned environment override, returning def when
// the variable is unset.
//
// strconv, not fmt.Sscanf("%d"): Sscanf stops at the first non-digit and reports
// success on the prefix, so "150000x" configured a gas limit of 150000 and
// "0x1e8480" configured one of 0. An operator who sets a limit and gets a
// different one silently is worse off than one who gets an error, because the
// symptom appears later as an out-of-gas revert with no obvious cause.
func envUint64(name string, def uint64) (uint64, error) {
	raw := os.Getenv(name)
	if raw == "" {
		return def, nil
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", name, err)
	}
	return value, nil
}

// envInt64 is envUint64 for the fee amounts, and it rejects a negative value
// here rather than letting it travel.
//
// An earlier version accepted one on the reasoning that the chain would reject it
// with a better message. That was wrong: every path feeds the amount into
// sdk.NewCoin, which PANICS on a negative amount ("negative coin amount: -1").
// Nothing reaches the chain -- the relayer process dies while building the
// transaction, on an operator typo.
//
// The second return reports whether the variable was SET, which is not the same
// question as whether the value is non-zero. COSMOS_FEE_AMOUNT=0 is a legitimate
// setting on a chain with no minimum gas price, and a caller that decides on the
// value alone silently ignores it -- see the batch path below.
//
// The value type stays int64 because sdkmath.NewInt takes one and the fee is
// compared against int64 overflow bounds downstream.
func envInt64(name string, def int64) (int64, bool, error) {
	raw := os.Getenv(name)
	if raw == "" {
		return def, false, nil
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("invalid %s: %w", name, err)
	}
	if value < 0 {
		return 0, false, fmt.Errorf("invalid %s: %d is negative; a fee amount cannot be", name, value)
	}
	return value, true, nil
}

// MisbehaviourGasLimit returns the configured EVM gas limit for an incident
// submission. Zero and malformed overrides fail before proof generation.
func MisbehaviourGasLimit() (uint64, error) {
	gasText := os.Getenv("ETH_MISBEHAVIOUR_GAS_LIMIT")
	if gasText == "" {
		return defaultMisbehaviourGasLimit, nil
	}
	gasLimit, err := strconv.ParseUint(gasText, 10, 64)
	if err != nil || gasLimit == 0 {
		return 0, fmt.Errorf("invalid ETH_MISBEHAVIOUR_GAS_LIMIT %q", gasText)
	}
	return gasLimit, nil
}

func misbehaviourPrivateKey() (*ecdsa.PrivateKey, error) {
	privateKeyText := os.Getenv("MISBEHAVIOUR_PRIVATE_KEY")
	if privateKeyText == "" {
		return nil, fmt.Errorf("MISBEHAVIOUR_PRIVATE_KEY environment variable is required")
	}
	privateKey, err := keys.RestoreKey(privateKeyText)
	if err != nil {
		return nil, fmt.Errorf("invalid MISBEHAVIOUR_PRIVATE_KEY: %w", err)
	}
	return privateKey, nil
}

// ValidateMisbehaviourPrivateKey validates the dedicated key before expensive
// proving begins.
func ValidateMisbehaviourPrivateKey() error {
	_, err := misbehaviourPrivateKey()
	return err
}

const ethDeployGasHeadroomPercent uint64 = 20

// Cosmos RPC deadlines (issue #119): bound every Cosmos broadcast/query so a
// stalled node can never wedge the relay goroutine. The relay path previously
// used BroadcastTxCommit, which blocks until the tx is committed — or forever.
const cosmosRPCTimeout = 30 * time.Second       // a single Cosmos RPC call (query / sync broadcast)
const cosmosInclusionTimeout = 90 * time.Second // total bounded wait for a broadcast tx to land in a block

func routerManagesProofSubmission(endpoint services.EVMEndpoint) bool {
	roleManager := endpoint.RoleManagerAddress()
	router := endpoint.RouterContract()
	if roleManager == nil || router == nil {
		return false
	}
	if *roleManager == (common.Address{}) || *router == (common.Address{}) {
		return false
	}
	return *roleManager == *router
}

// RouterManagesProofSubmission reports whether proof calls must route through
// ICS26Router rather than directly to SpectreClient.
func RouterManagesProofSubmission(endpoint services.EVMEndpoint) bool {
	return routerManagesProofSubmission(endpoint)
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
func logInnerGasFromTrace(stdCtx context.Context, endpoint services.EVMEndpoint, txHash common.Hash, labels []string) error {
	rpcClient := endpoint.EthClient().Client()
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
	// Base 16 explicitly, and the whole string: Sscanf("%x") accepted "1fzz" as
	// 0x1f, which turned a corrupted trace field into a plausible gas number.
	v, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %q: %w", s, err)
	}
	return v, nil
}

func cosmosRouterClientID(ids services.ClientIDs) (string, error) {
	clientID := ids.CosmosOnEVM
	if clientID == "" {
		return "", fmt.Errorf("cosmos router client id is not configured")
	}
	return clientID, nil
}

func ethLightClientIDOnCosmos(ids services.ClientIDs) (string, error) {
	clientID := ids.EVMOnCosmos
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

// baseFeeHeadroomPercent is how far above the current base fee a price is lifted
// before it is used. EIP-1559 lets the base fee rise 12.5% per block, so 100%
// covers roughly six consecutive full blocks — enough for the gap between reading
// a price and the transaction landing, without overpaying: the base fee is burned
// at its actual value, and the excess is refunded.
const baseFeeHeadroomPercent = 100

// gasPriceWithBaseFeeHeadroom lifts a suggested price clear of the current base
// fee.
//
// SuggestGasPrice reflects the chain at the moment it is called. On a chain whose
// base fee is climbing, the value is already stale by the time it reaches
// EstimateGas, and the node rejects the call outright:
//
//	max fee per gas less than block base fee: maxFeePerGas: 21486000, baseFee: 22414000
//
// which surfaces as an estimation failure rather than a pricing one. Local devnets
// never show this — their base fee sits at the floor — so it only appears against
// a public network, at the point where a deployment is being set up.
//
// A failure to read the header leaves the suggested price untouched: this is a
// safety margin, not a correctness requirement.
func gasPriceWithBaseFeeHeadroom(stdCtx context.Context, endpoint services.EVMEndpoint, suggested *big.Int) *big.Int {
	header, err := endpoint.EthClient().HeaderByNumber(stdCtx, nil)
	if err != nil || header == nil || header.BaseFee == nil {
		return suggested
	}
	lifted := liftAboveBaseFee(suggested, header.BaseFee)
	if lifted != suggested {
		log.Printf("[gas] lifting suggested price %v to %v to clear base fee %v",
			suggested, lifted, header.BaseFee)
	}
	return lifted
}

// liftAboveBaseFee raises a price to baseFee + headroom when it sits below that,
// and otherwise returns it untouched. Split out from the RPC call so the
// arithmetic is testable on its own.
func liftAboveBaseFee(suggested, baseFee *big.Int) *big.Int {
	if baseFee == nil {
		return suggested
	}
	floor := new(big.Int).Mul(baseFee, big.NewInt(100+baseFeeHeadroomPercent))
	floor.Div(floor, big.NewInt(100))
	if suggested == nil || suggested.Cmp(floor) < 0 {
		return floor
	}
	return suggested
}

func estimateCosmosClientDeployGas(
	stdCtx context.Context,
	endpoint services.EVMEndpoint,
	from common.Address,
	gasPrice *big.Int,
	clientState []byte,
	consensusState spectreContract.IICS07TendermintMsgsConsensusState,
	initialPinnedValidatorSet spectreContract.IICS07TendermintMsgsValidatorSet,
) (uint64, uint64, error) {
	parsed, err := spectreContract.ContractSpectreClientMetaData.GetAbi()
	if err != nil {
		return 0, 0, fmt.Errorf("parse ICS07 ABI: %w", err)
	}
	constructorInput, err := parsed.Pack(
		"",
		*endpoint.UpdateClientContract(),
		*endpoint.MembershipContract(),
		*endpoint.MisbehaviourContract(),
		clientState,
		consensusState,
		initialPinnedValidatorSet,
		*endpoint.RoleManagerAddress(),
	)
	if err != nil {
		return 0, 0, fmt.Errorf("pack ICS07 constructor args: %w", err)
	}
	deployData := append(common.FromHex(spectreContract.ContractSpectreClientBin), constructorInput...)
	estimate, err := endpoint.EthClient().EstimateGas(stdCtx, ethereum.CallMsg{
		From:     from,
		GasPrice: gasPriceWithBaseFeeHeadroom(stdCtx, endpoint, gasPrice),
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
	if header, err := endpoint.EthClient().HeaderByNumber(stdCtx, nil); err == nil && header.GasLimit > 0 {
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

func (h *Handler) CreateCosmosClientContract(stdCtx context.Context, endpoint services.EVMEndpoint, clientIDs services.ClientIDs, clientState []byte, consensusState spectreContract.IICS07TendermintMsgsConsensusState, initialPinnedValidatorSet relayerclient.ContractValidatorSet) (common.Address, error) {
	cosmosClientID, err := cosmosRouterClientID(clientIDs)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] %w", err)
	}
	wasmClientID, err := ethLightClientIDOnCosmos(clientIDs)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] %w", err)
	}
	privateKey, err := h.keySigner().EthKey()
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to restore private key: %w", err)
	}

	publicKey, err := keys.PublicKey(privateKey)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] failed to derive public key: %w", err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	gasPrice, err := endpoint.EthClient().SuggestGasPrice(stdCtx)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] failed to suggest gas price: %w", err)
	}

	pinnedForDeploy := toGroth16ValidatorSet(initialPinnedValidatorSet)
	estimatedDeployGas, deployGasLimit, err := estimateCosmosClientDeployGas(stdCtx, endpoint, fromAddress, gasPrice, clientState, consensusState, pinnedForDeploy)
	if err != nil {
		return common.Address{}, fmt.Errorf("[CreateCosmosClient] %w", err)
	}
	log.Printf("[CreateCosmosClient] ICS07 deploy gas estimate=%d limit=%d", estimatedDeployGas, deployGasLimit)

	var address common.Address
	deployFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		addr, tx, _, err := spectreContract.DeployContractSpectreClient(
			auth,
			endpoint.EthClient(),
			*endpoint.UpdateClientContract(),
			*endpoint.MembershipContract(),
			*endpoint.MisbehaviourContract(),
			clientState,
			consensusState,
			pinnedForDeploy,
			*endpoint.RoleManagerAddress(),
		)
		if err == nil {
			address = addr
		}
		return tx, err
	}

	receipt, _, _, err := h.executeWithRetryAndResubmission(stdCtx, endpoint, privateKey, deployGasLimit, deployFn, knobDeployEstimate)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed waiting for deploy receipt: %w", err)
	}
	log.Printf("[CreateCosmosClient] ICS07 deployed at %s (block %d, gasUsed=%d)", address.String(), receipt.BlockNumber.Uint64(), receipt.GasUsed)

	ics26Router, err := routerContract.NewContractICS26Router(*endpoint.RouterContract(), endpoint.EthClient())
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
			address,
		)
	}

	addClientReceipt, _, _, err := h.executeWithRetryAndResubmission(stdCtx, endpoint, privateKey, 16000000, addClientFn, knobFixedClientLimit)
	if err != nil {
		if errors.Is(err, services.ErrPermanentRelayFailure) {
			return common.Address{}, fmt.Errorf(
				"AddClient permanently failed for %s; governed client migration is required and is not submitted by the relayer hot key: %w",
				cosmosClientID,
				err,
			)
		}
		// For transient wait/RPC errors, return the error immediately so the caller can retry
		return common.Address{}, err
	}

	log.Printf("[CreateCosmosClient] AddClient confirmed (block %d, gasUsed=%d)", addClientReceipt.BlockNumber.Uint64(), addClientReceipt.GasUsed)
	return address, nil
}

func (h *Handler) SendEthTx(stdCtx context.Context, endpoint services.EVMEndpoint, cosmosClientID string, msg any) error {
	if cosmosClientID == "" {
		return fmt.Errorf("[SendEthTx] cosmos router client id is not configured")
	}
	privateKey, err := h.keySigner().EthKey()
	if err != nil {
		return fmt.Errorf("failed to restore private key: %w", err)
	}

	gasLimit := defaultEthGasLimit
	if gasStr := os.Getenv("ETH_GAS_LIMIT"); gasStr != "" {
		configuredGasLimit, err := EthGasLimit(gasStr)
		if err != nil {
			log.Printf("[SendEthTx] %v; using default %d", err, defaultEthGasLimit)
		} else {
			gasLimit = configuredGasLimit
		}
	}

	ics07Tendermint, err := spectreContract.NewContractSpectreClient(
		*endpoint.SpectreClientContract(),
		endpoint.EthClient(),
	)
	if err != nil {
		return fmt.Errorf("failed to create ICS07 Tendermint contract: %w", err)
	}

	ics26Router, err := contractICS26Router.NewContractICS26Router(
		*endpoint.RouterContract(),
		endpoint.EthClient(),
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
				if routerManagesProofSubmission(endpoint) {
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
				if routerManagesProofSubmission(endpoint) {
					log.Printf("[SendEthTx] Sending ICS26Router.updateApplicationState tx for clientId=%s...", cosmosClientID)
					return ics26Router.UpdateApplicationState(auth, cosmosClientID, data)
				}
				log.Printf("[SendEthTx] Sending direct SpectreClient.updateApplicationState tx...")
				return ics07Tendermint.UpdateApplicationState(auth, data)
			}
		case spectreContract.ILightClientMsgsMsgVerifyMembership:
			txLabel = "verifyMembership"
			if routerManagesProofSubmission(endpoint) {
				return nil, fmt.Errorf("direct verifyMembership is disabled when ROLE_MANAGER is the ICS26 router; use ICS26Router packet flows instead")
			}
			log.Printf("[SendEthTx] Sending verifyMembership tx...")
			return ics07Tendermint.VerifyMembership(auth, msg)
		case spectreContract.ILightClientMsgsMsgVerifyNonMembership:
			txLabel = "verifyNonMembership"
			if routerManagesProofSubmission(endpoint) {
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

	receipt, submitDur, waitDur, err := h.executeWithRetryAndResubmission(stdCtx, endpoint, privateKey, gasLimit, senderFn, knobEthGasLimit)
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

// SubmitMisbehaviour submits proof-backed same-height equivocation evidence
// with the dedicated incident-response signer. It intentionally stays outside
// SendEthTx's generic message switch so it cannot fall back to ETH_PRIVATE_KEY.
func (h *Handler) SubmitMisbehaviour(
	stdCtx context.Context,
	endpoint services.EVMEndpoint,
	cosmosRouterClientID string,
	misbehaviourMsg []byte,
) error {
	if len(misbehaviourMsg) == 0 {
		return fmt.Errorf("empty misbehaviour message")
	}
	if endpoint.EthClient() == nil || endpoint.SpectreClientContract() == nil ||
		*endpoint.SpectreClientContract() == (common.Address{}) {
		return fmt.Errorf("SpectreClient contract and EVM RPC are required")
	}
	privateKey, err := misbehaviourPrivateKey()
	if err != nil {
		return err
	}
	publicKey, err := keys.PublicKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to derive misbehaviour key address: %w", err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	spectre, err := spectreContract.NewContractSpectreClient(*endpoint.SpectreClientContract(), endpoint.EthClient())
	if err != nil {
		return fmt.Errorf("failed to bind Spectre client: %w", err)
	}
	role, err := spectre.MISBEHAVIOURSUBMITTERROLE(&bind.CallOpts{Context: stdCtx})
	if err != nil {
		return fmt.Errorf("query MISBEHAVIOUR_SUBMITTER_ROLE: %w", err)
	}
	permissionless, err := spectre.HasRole(&bind.CallOpts{Context: stdCtx}, role, common.Address{})
	if err != nil {
		return fmt.Errorf("query permissionless misbehaviour role: %w", err)
	}

	useRouter := routerManagesProofSubmission(endpoint)
	var router *contractICS26Router.ContractICS26Router
	if useRouter {
		if endpoint.RouterContract() == nil || *endpoint.RouterContract() == (common.Address{}) {
			return fmt.Errorf("router contract address is required")
		}
		if cosmosRouterClientID == "" {
			return fmt.Errorf("cosmos router client id is required")
		}
		routerHasRole, roleErr := spectre.HasRole(
			&bind.CallOpts{Context: stdCtx}, role, *endpoint.RouterContract(),
		)
		if roleErr != nil {
			return fmt.Errorf("query router misbehaviour role: %w", roleErr)
		}
		if !permissionless && !routerHasRole {
			return fmt.Errorf(
				"ICS26Router %s does not hold SpectreClient.MISBEHAVIOUR_SUBMITTER_ROLE",
				endpoint.RouterContract().Hex(),
			)
		}
		router, err = contractICS26Router.NewContractICS26Router(*endpoint.RouterContract(), endpoint.EthClient())
		if err != nil {
			return fmt.Errorf("failed to bind ICS26 router: %w", err)
		}
	} else if !permissionless {
		signerHasRole, roleErr := spectre.HasRole(&bind.CallOpts{Context: stdCtx}, role, fromAddress)
		if roleErr != nil {
			return fmt.Errorf("query signer misbehaviour role: %w", roleErr)
		}
		if !signerHasRole {
			return fmt.Errorf(
				"misbehaviour signer %s does not hold SpectreClient.MISBEHAVIOUR_SUBMITTER_ROLE",
				fromAddress.Hex(),
			)
		}
	}

	gasLimit, err := MisbehaviourGasLimit()
	if err != nil {
		return err
	}
	sender := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		if useRouter {
			return router.SubmitMisbehaviour(auth, cosmosRouterClientID, misbehaviourMsg)
		}
		return spectre.Misbehaviour(auth, misbehaviourMsg)
	}
	receipt, _, _, err := h.executeWithRetryAndResubmission(stdCtx, endpoint, privateKey, gasLimit, sender, knobMisbehaviour)
	if err != nil {
		return fmt.Errorf("submit misbehaviour: %w", err)
	}
	log.Printf("[SubmitMisbehaviour] report signed by %s confirmed in block %d", fromAddress.Hex(), receipt.BlockNumber.Uint64())
	return nil
}

// SendEthTxBatch packs an optional Cosmos-client update followed by N
// packet-level ICS26Router calls into a single `multicall(bytes[])` tx. Folding
// the update is supported only when the ICS26Router is the proof submitter;
// direct SpectreClient deployments use the caller's standalone-update fallback.
//
// Empty input is a no-op. A single-msg input is forwarded to SendEthTx to
// avoid the multicall wrapper's small overhead for N=1 (matches the
// threshold contract in the V1 plan).
//
// Semantics: MulticallUpgradeable runs each inner call via delegatecall and
// reverts the whole tx if any inner call reverts — so caller can treat
// success as "every packet in the batch was relayed".
func (h *Handler) SendEthTxBatch(stdCtx context.Context, endpoint services.EVMEndpoint, cosmosRouterID string, msgs []any) error {
	if len(msgs) == 0 {
		return nil
	}
	if len(msgs) == 1 {
		return h.SendEthTx(stdCtx, endpoint, cosmosRouterID, msgs[0])
	}

	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed to load ICS26Router ABI: %w", err)
	}

	// cosmosClientID is only consulted when the batch contains an updateClient
	// inner call; resolved lazily to avoid coupling pure-recvPacket batches to
	// the router-managed config check.
	cosmosClientID := cosmosRouterID
	clientIDResolved := cosmosRouterID != ""

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
			if !routerManagesProofSubmission(endpoint) {
				return fmt.Errorf("[SendEthTxBatch] client update cannot be batched when ICS26Router is not the proof submitter; submit it via SendEthTx instead")
			}
			if !clientIDResolved {
				if cosmosClientID == "" {
					return fmt.Errorf("[SendEthTxBatch] resolve cosmos client id: cosmos router client id is not configured")
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

	privateKey, err := h.keySigner().EthKey()
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] failed to restore private key: %w", err)
	}
	// A malformed override used to be discarded in silence -- the operator set a
	// limit, the default was used, and nothing said so. Fail instead, matching
	// EthGasLimit above.
	multicallGasLimit, err := envUint64("ETH_MULTICALL_GAS_LIMIT", 16000000)
	if err != nil {
		return fmt.Errorf("[SendEthTxBatch] %w", err)
	}

	ics26Router, err := contractICS26Router.NewContractICS26Router(*endpoint.RouterContract(), endpoint.EthClient())
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

	receipt, submitDur, waitDur, err := h.executeWithRetryAndResubmission(stdCtx, endpoint, privateKey, multicallGasLimit, senderFn, knobMulticallGasLimit)
	if err != nil {
		return fmt.Errorf("multicall labels=%s: %w", labelStr, err)
	}

	log.Printf("[SendEthTxBatch] Tx %s confirmed in block %d (gasUsed=%d, inner=%d)",
		receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), receipt.GasUsed, len(calldata))
	if benchEnabled {
		log.Printf("[bench][eth] multicall labels=%s gasUsed=%d submit=%s wait=%s total=%s tx=%s",
			labelStr, receipt.GasUsed, submitDur, waitDur, time.Since(benchStart), receipt.TxHash.Hex())
		if traceErr := logInnerGasFromTrace(stdCtx, endpoint, receipt.TxHash, labels); traceErr != nil {
			log.Printf("[bench][eth] inner gas trace unavailable (RPC may lack debug_ namespace): %v", traceErr)
		}
	}

	return nil
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

func (h *Handler) bumpGasAndResubmit(stdCtx context.Context, endpoint services.EVMEndpoint, tx *types.Transaction, auth *bind.TransactOpts, attempt int) (*types.Transaction, error) {
	currentGasPrice := tx.GasPrice()

	suggestedGasPrice, err := endpoint.EthClient().SuggestGasPrice(stdCtx)
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

		if suggestedTipCap, err := endpoint.EthClient().SuggestGasTipCap(stdCtx); err == nil {
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
	sendErr := endpoint.EthClient().SendTransaction(sendCtx, signedTx)
	if sendErr != nil {
		return nil, fmt.Errorf("failed to send bumped transaction: %w", sendErr)
	}

	return signedTx, nil
}

func (h *Handler) executeWithRetryAndResubmission(
	stdCtx context.Context,
	endpoint services.EVMEndpoint,
	privateKey *ecdsa.PrivateKey,
	gasLimit uint64,
	senderFn func(auth *bind.TransactOpts) (*types.Transaction, error),
	knobs ...gasLimitKnob,
) (*types.Receipt, time.Duration, time.Duration, error) {
	knob := gasLimitKnob("the gas limit for this operation")
	if len(knobs) > 0 {
		knob = knobs[0]
	}
	if gasLimit == 0 {
		return nil, 0, 0, fmt.Errorf("cannot submit %s with a zero gas limit", knob)
	}
	return h.executeWithRetryAndResubmissionAtGasStep(
		stdCtx, endpoint, privateKey, gasLimit, 0, 0, senderFn, knob,
	)
}

func (h *Handler) executeWithRetryAndResubmissionAtGasStep(
	stdCtx context.Context,
	endpoint services.EVMEndpoint,
	privateKey *ecdsa.PrivateKey,
	baseGasLimit uint64,
	gasStep int,
	gasCeiling uint64,
	senderFn func(auth *bind.TransactOpts) (*types.Transaction, error),
	knob gasLimitKnob,
) (*types.Receipt, time.Duration, time.Duration, error) {
	gasLimit := applyEVMGasHeadroom(baseGasLimit, gasStep)
	if gasCeiling > 0 && gasLimit > gasCeiling {
		gasLimit = gasCeiling
	}
	publicKey, err := keys.PublicKey(privateKey)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to derive public key: %w", err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKey)
	chainIdInt, err := endpoint.EthClient().ChainID(stdCtx)
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
	if tip, err := endpoint.EthClient().SuggestGasTipCap(stdCtx); err == nil {
		isEIP1559 = true
		suggestedTip = tip
		if header, err := endpoint.EthClient().HeaderByNumber(stdCtx, nil); err == nil && header.BaseFee != nil {
			baseFee = header.BaseFee
		} else {
			baseFee = big.NewInt(1000000000) // fallback 1 Gwei
		}
	}

	for nonceAttempt := 1; nonceAttempt <= maxNonceRetries; nonceAttempt++ {
		h.mu.Lock()
		state := h.senderState(chainIdInt.String(), fromAddress)
		if !state.nonceValid {
			n, err := endpoint.EthClient().PendingNonceAt(stdCtx, fromAddress)
			if err != nil {
				h.mu.Unlock()
				return nil, 0, 0, fmt.Errorf("failed to get pending nonce: %w", err)
			}
			state.nonce = n
			state.nonceValid = true
		}
		currentNonce := state.nonce

		if isEIP1559 {
			gasTipCap := suggestedTip
			gasFeeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), gasTipCap)

			// Enforce minimum floor if this nonce matches the last attempted nonce (e.g. replacing a stuck tx)
			if currentNonce == state.lastNonce {
				if state.lastGasTipCap != nil {
					minTip := new(big.Int).Mul(state.lastGasTipCap, big.NewInt(115))
					minTip.Div(minTip, big.NewInt(100))
					if gasTipCap.Cmp(minTip) < 0 {
						gasTipCap = minTip
					}
				}
				if state.lastGasFeeCap != nil {
					minFee := new(big.Int).Mul(state.lastGasFeeCap, big.NewInt(115))
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
			gasPrice, err := endpoint.EthClient().SuggestGasPrice(stdCtx)
			if err != nil {
				state.nonceValid = false
				h.mu.Unlock()
				return nil, 0, 0, fmt.Errorf("failed to suggest gas price: %w", err)
			}

			// Enforce minimum floor if this nonce matches the last attempted nonce (e.g. replacing a stuck tx)
			if currentNonce == state.lastNonce {
				if state.lastGasPrice != nil {
					minPrice := new(big.Int).Mul(state.lastGasPrice, big.NewInt(115))
					minPrice.Div(minPrice, big.NewInt(100))
					if gasPrice.Cmp(minPrice) < 0 {
						gasPrice = minPrice
					}
				}
			}

			// Clear the current base fee: the suggestion was read a moment ago and a
			// rising base fee makes it stale, which the node rejects outright rather
			// than queueing.
			gasPrice = gasPriceWithBaseFeeHeadroom(stdCtx, endpoint, gasPrice)

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
				state.nonceValid = false
				h.mu.Unlock()
				return nil, 0, 0, fmt.Errorf("senderFn returned nil transaction without error")
			}
			tx = signedTx
			submitDur = time.Since(submitStart)

			state.nonce++
			state.lastNonce = tx.Nonce()
			if tx.Type() == types.DynamicFeeTxType {
				state.lastGasFeeCap = tx.GasFeeCap()
				state.lastGasTipCap = tx.GasTipCap()
				state.lastGasPrice = nil
			} else {
				state.lastGasPrice = tx.GasPrice()
				state.lastGasFeeCap = nil
				state.lastGasTipCap = nil
			}
			h.mu.Unlock()

			break
		}

		if isNonceTooLowError(callErr) {
			log.Printf("[EthTxSender] Nonce %d too low (attempt %d/%d). Resetting nonce cache.", currentNonce, nonceAttempt, maxNonceRetries)
			state.nonceValid = false
			h.mu.Unlock()
			continue
		}

		if isAlreadyKnownError(callErr) && signedTx != nil {
			log.Printf("[EthTxSender] Transaction already known in mempool: %s. Proceeding to wait.", signedTx.Hash().Hex())
			tx = signedTx
			submitDur = time.Since(submitStart)

			state.nonce++
			state.lastNonce = tx.Nonce()
			if tx.Type() == types.DynamicFeeTxType {
				state.lastGasFeeCap = tx.GasFeeCap()
				state.lastGasTipCap = tx.GasTipCap()
				state.lastGasPrice = nil
			} else {
				state.lastGasPrice = tx.GasPrice()
				state.lastGasFeeCap = nil
				state.lastGasTipCap = nil
			}
			h.mu.Unlock()

			break
		}

		// Other error: invalidate nonce just in case and return
		state.nonceValid = false
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
		receipt, waitErr := waitForReceipts(receiptCtx, endpoint.EthClient(), sentHashes)
		cancel()

		if waitErr == nil {
			waitDur := time.Since(waitStart)
			if receipt.Status == 0 {
				callMsg := ethereum.CallMsg{
					From:     fromAddress,
					To:       tx.To(),
					Gas:      tx.Gas(),
					GasPrice: tx.GasPrice(),
					Value:    tx.Value(),
					Data:     tx.Data(),
				}
				callCtx, callCancel := context.WithTimeout(stdCtx, ethRevertProbeTimeout)
				_, callErr := endpoint.EthClient().CallContract(callCtx, callMsg, receipt.BlockNumber)
				callCancel()
				if callErr != nil {
					log.Printf("[EthTxSender] Revert reason: %v", callErr)
					if name, ok := validatorCacheRaceErrorName(callErr); ok {
						return receipt, submitDur, waitDur, fmt.Errorf(
							"tx %s reverted with %s (status=0, gasUsed=%d): %w",
							receipt.TxHash.Hex(), name, receipt.GasUsed, services.ErrValidatorCacheRace,
						)
					}
				}
				// A successful original-gas replay is still diagnostic input. The
				// classifier intentionally accepts a nil callErr and uses the bounded
				// higher-gas probe to decide whether the receipt ran out of gas.
				if classifyEthRevert(stdCtx, endpoint, callMsg, receipt.BlockNumber, callErr, receipt.GasUsed, tx.Gas()) == revertOutOfGas {
					nextCeiling := gasCeiling
					headerCtx, headerCancel := context.WithTimeout(stdCtx, ethRevertProbeTimeout)
					if header, headerErr := endpoint.EthClient().HeaderByNumber(headerCtx, receipt.BlockNumber); headerErr == nil && header.GasLimit > 0 {
						nextCeiling = header.GasLimit
					}
					headerCancel()
					// Diagnose the chain ceiling before ladder exhaustion. Both conditions
					// can become true on the final step; checking the ladder first masked the
					// actionable block-limit reason behind the generic exhaustion error.
					if nextCeiling > 0 && tx.Gas() >= nextCeiling {
						return receipt, submitDur, waitDur, fmt.Errorf(
							"tx %s ran out of gas at the block gas ceiling %d for %s: %w",
							receipt.TxHash.Hex(), tx.Gas(), knob, services.ErrPermanentRelayFailure,
						)
					}

					// Exhaustion is not decided here. The ladder is described, and
					// chain.Climb is the single place that turns its bottom into a
					// permanent failure -- the same place the two Cosmos out-of-gas
					// sites use, so the three cannot drift apart.
					cause := fmt.Errorf(
						"tx %s ran out of gas (gasUsed=%d of limit %d); exhausted finite gas ladder for %s: %w",
						receipt.TxHash.Hex(), receipt.GasUsed, tx.Gas(), knob, services.ErrPermanentRelayFailure,
					)
					next, bottom, ok := chain.Climb(chain.NeedsChange(cause, EVMGasLadder(baseGasLimit)), gasStep)
					if !ok {
						return receipt, submitDur, waitDur, bottom
					}

					nextStep := gasStep + 1
					nextGasLimit := next.To
					if nextCeiling > 0 && nextGasLimit > nextCeiling {
						nextGasLimit = nextCeiling
					}
					log.Printf("[EthTxSender] Tx %s ran out of gas at %d; retrying with fresh nonce and gas limit %d (step %d/%d)",
						receipt.TxHash.Hex(), tx.Gas(), nextGasLimit, nextStep+1, len(evmGasHeadroomBasisPoints))
					retryReceipt, retrySubmitDur, retryWaitDur, retryErr := h.executeWithRetryAndResubmissionAtGasStep(
						stdCtx, endpoint, privateKey, baseGasLimit, nextStep, nextCeiling, senderFn, knob,
					)
					return retryReceipt, submitDur + retrySubmitDur, waitDur + retryWaitDur, retryErr
				}
				return receipt, submitDur, waitDur, fmt.Errorf(
					"tx %s reverted (status=0, gasUsed=%d): %w",
					receipt.TxHash.Hex(), receipt.GasUsed, services.ErrPermanentRelayFailure,
				)
			}
			return receipt, submitDur, waitDur, nil
		}

		if errors.Is(waitErr, context.DeadlineExceeded) {
			// A transaction that cannot mine because nonces below it are missing
			// will not mine with more gas either, so bumping only pays to resubmit
			// the same unminable transaction until the budget runs out.
			//
			// PendingNonceAt reports the next nonce the account needs. A queued
			// (non-executable) transaction does not advance it, so a pending nonce
			// BELOW ours means the gap is real: nonces [pending, ours) are missing.
			// Equal means our transaction is simply not in the pool — a drop, which
			// a resubmit can still fix — so only the strict inequality short-circuits.
			if pending, nonceErr := endpoint.EthClient().PendingNonceAt(stdCtx, fromAddress); nonceErr != nil {
				log.Printf("[EthTxSender] Tx %s not mined in %s; pending-nonce probe failed: %v",
					tx.Hash().Hex(), attemptTimeout, nonceErr)
			} else if pending < tx.Nonce() {
				h.mu.Lock()
				h.senderState(chainIdInt.String(), fromAddress).nonceValid = false
				h.mu.Unlock()
				return nil, submitDur, time.Since(waitStart), fmt.Errorf(
					"tx %s cannot mine: sent with nonce %d but account %s still needs nonce %d, so %d-%d are missing; "+
						"cached nonce invalidated, retry will re-query",
					tx.Hash().Hex(), tx.Nonce(), fromAddress.Hex(), pending, pending, tx.Nonce()-1)
			}

			if attempt >= maxAttempts {
				h.mu.Lock()
				h.senderState(chainIdInt.String(), fromAddress).nonceValid = false
				h.mu.Unlock()
				return nil, submitDur, time.Since(waitStart), fmt.Errorf("transaction wait mined timed out after %d attempts (last hash: %s): %w", attempt, tx.Hash().Hex(), waitErr)
			}

			log.Printf("[EthTxSender] Tx %s not mined in %s, bumping gas price...", tx.Hash().Hex(), attemptTimeout)
			bumpedTx, bumpErr := h.bumpGasAndResubmit(stdCtx, endpoint, tx, auth, attempt)
			if bumpErr != nil {
				log.Printf("[EthTxSender] Gas bump attempt %d failed: %v. Will continue waiting.", attempt, bumpErr)
			} else {
				tx = bumpedTx
				sentHashes = append(sentHashes, tx.Hash())
				log.Printf("[EthTxSender] Gas bumped tx submitted: %s (attempt %d)", tx.Hash().Hex(), attempt+1)

				h.mu.Lock()
				state := h.senderState(chainIdInt.String(), fromAddress)
				state.lastNonce = tx.Nonce()
				if tx.Type() == types.DynamicFeeTxType {
					state.lastGasFeeCap = tx.GasFeeCap()
					state.lastGasTipCap = tx.GasTipCap()
					state.lastGasPrice = nil
				} else {
					state.lastGasPrice = tx.GasPrice()
					state.lastGasFeeCap = nil
					state.lastGasTipCap = nil
				}
				h.mu.Unlock()
			}
			attempt++
			continue
		}

		h.mu.Lock()
		h.senderState(chainIdInt.String(), fromAddress).nonceValid = false
		h.mu.Unlock()
		return nil, submitDur, time.Since(waitStart), fmt.Errorf("failed waiting for tx receipt: %w", waitErr)
	}
}
