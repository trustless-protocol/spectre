package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	tendermintContract "relayer/bindings/Groth16ICS07Tendermint"
	updateclientContract "relayer/bindings/UpdateClient"
	relayerclient "relayer/client"
	"relayer/prover"
	"sort"
	"strconv"
	"strings"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const ICS26_IBC_STORAGE_SLOT = "0x1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600"

type Worker struct {
	TxHandler TransactionHandler
	Prover    Prover
}

const cosmosCatchUpSafetySlots uint64 = 3

func NewWorker(txHandler TransactionHandler, prover Prover) *Worker {
	return &Worker{
		txHandler,
		prover,
	}
}

func (w *Worker) CreateCosmosClient(ctx Context, proofType string, trustingPeriod uint32, trustedBlock int64, trustLevel string) (common.Address, error) {
	genesis, err := relayerclient.GetGenesis(ctx.CosmosClient(), trustedBlock, trustingPeriod, trustLevel, proofType)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get genesis: %w", err)
	}

	clientState := genesis.TrustedClientState
	consensusState := genesis.TrustedConsensusState

	log.Printf("[CreateCosmosClient] clientState: chainId=%s trustLevel=%d/%d height=%d/%d trustingPeriod=%d unbondingPeriod=%d isFrozen=%v zkAlgorithm=%d",
		clientState.ChainId, clientState.TrustLevel.Numerator, clientState.TrustLevel.Denominator,
		clientState.LatestHeight.RevisionNumber, clientState.LatestHeight.RevisionHeight,
		clientState.TrustingPeriod, clientState.UnbondingPeriod, clientState.IsFrozen, clientState.ZkAlgorithm)

	clientStateEncoded, err := relayerclient.EncodeClientState(clientState)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to encode client state: %w", err)
	}

	consensusStateEncoded, err := relayerclient.EncodeConsensusState(consensusState)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to encode consensus state: %w", err)
	}

	consensusHash := crypto.Keccak256(consensusStateEncoded)
	log.Printf("[CreateCosmosClient] consensusHash=%x", consensusHash)
	return w.TxHandler.CreateCosmosClientContract(ctx, clientStateEncoded, consensusHash, genesis.InitialPinnedValidatorSet)
}

// CosmosClientUpdateBuildResult is the output of BuildCosmosClientUpdateMsg.
// When HasMsg is false the on-chain client is already at the latest height
// and no updateClient tx is needed — only LightBlock is populated.
type CosmosClientUpdateBuildResult struct {
	Msg        updateclientContract.IUpdateClientMsgsMsgUpdateClient
	HasMsg     bool
	LightBlock *relayerclient.LightBlock
}

// UpdateCosmosClient builds the next MsgUpdateClient for the Tendermint light
// client on Ethereum and submits it via SendEthTx. Used by background routines
// that want to advance the client without packets attached. handleCosmos uses
// the split BuildCosmosClientUpdateMsg builder so it can fold updateClient
// into the same multicall as its packet calls (issue #67 V2).
func (w *Worker) UpdateCosmosClient(ctx Context, proofType string, trustedBlock int64, trustLevel string) (*relayerclient.LightBlock, error) {
	result, err := w.BuildCosmosClientUpdateMsg(ctx, proofType, trustedBlock, trustLevel)
	if err != nil {
		return nil, err
	}
	if !result.HasMsg {
		return result.LightBlock, nil
	}
	if err := w.TxHandler.SendEthTx(ctx, result.Msg); err != nil {
		log.Printf("[UpdateCosmosClient] SendEthTx failed: %v", err)
		return nil, err
	}
	log.Printf("[UpdateCosmosClient] SendEthTx succeeded")
	return result.LightBlock, nil
}

// ReAnchorCosmosPinnedSet builds a reAnchorPinnedSet transaction for the ICS07
// light client on Ethereum. It reads the authoritative on-chain trusted height,
// generates a Groth16 proof from that stored height to the latest CometBFT block,
// then calls reAnchorPinnedSet to advance both the light client and the pinned
// validator set to the block's nextValidatorsHash.
//
// When the on-chain client is already caught up (no newer block exists) the
// function returns the current light block without submitting any tx — the
// caller can then update its cached timestamp/height from the returned block.
func (w *Worker) ReAnchorCosmosPinnedSet(
	ctx Context,
	proofType string,
	trustLevel string,
) (*relayerclient.LightBlock, error) {
	onChainTrusted, err := fetchOnChainTrustedHeight(ctx)
	if err != nil {
		return nil, fmt.Errorf("[ReAnchorCosmosPinnedSet] fetch on-chain height: %w", err)
	}
	status, err := ctx.CosmosClient().Status(context.Background())
	if err != nil {
		return nil, fmt.Errorf("[ReAnchorCosmosPinnedSet] get status: %w", err)
	}
	needsReAnchor, err := shouldReAnchor(onChainTrusted, status.SyncInfo.LatestBlockHeight)
	if err != nil {
		return nil, fmt.Errorf("[ReAnchorCosmosPinnedSet] compare heights: %w", err)
	}
	if !needsReAnchor {
		log.Printf("[ReAnchorCosmosPinnedSet] client is up to date (trusted=%d, latest=%d), skipping",
			onChainTrusted, status.SyncInfo.LatestBlockHeight)
		lightBlock, err := relayerclient.GetLightBlock(ctx.CosmosClient(), onChainTrusted)
		if err != nil {
			return nil, fmt.Errorf("[ReAnchorCosmosPinnedSet] get current light block: %w", err)
		}
		return lightBlock, nil
	}

	result, err := w.BuildCosmosClientUpdateMsg(ctx, proofType, onChainTrusted, trustLevel)
	if err != nil {
		return nil, fmt.Errorf("[ReAnchorCosmosPinnedSet] build update msg: %w", err)
	}
	if !result.HasMsg || result.LightBlock == nil {
		return nil, fmt.Errorf("[ReAnchorCosmosPinnedSet] build returned no update light block")
	}

	newPinnedSet, err := relayerclient.ValidatorSetToContract(result.LightBlock.NextValSet, "re-anchor")
	if err != nil {
		return nil, fmt.Errorf("[ReAnchorCosmosPinnedSet] convert next validator set: %w", err)
	}

	if err := w.TxHandler.SendReAnchorPinnedSet(ctx, result.Msg, newPinnedSet); err != nil {
		return nil, fmt.Errorf("[ReAnchorCosmosPinnedSet] send tx: %w", err)
	}
	log.Printf("[ReAnchorCosmosPinnedSet] re-anchor succeeded")
	return result.LightBlock, nil
}

func shouldReAnchor(onChainTrusted, latestBlockHeight int64) (bool, error) {
	if onChainTrusted > latestBlockHeight {
		return false, fmt.Errorf(
			"trusted block is ahead of latest chain height (trusted=%d, latest=%d)",
			onChainTrusted,
			latestBlockHeight,
		)
	}
	return onChainTrusted < latestBlockHeight, nil
}

// fetchOnChainTrustedHeight reads the ICS07 client state on ETH and returns its
// latest trusted revision height. This is a cheap eth_call relative to the
// Groth16 proof, so it's always worth doing before committing to proof gen.
func fetchOnChainTrustedHeight(ctx Context) (int64, error) {
	ics07, err := tendermintContract.NewContractGroth16ICS07Tendermint(*ctx.ClientContract(), ctx.EthClient())
	if err != nil {
		return 0, fmt.Errorf("failed to create ICS07 instance: %w", err)
	}
	clientStateBytes, err := ics07.GetClientState(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get on-chain client state: %w", err)
	}
	onChainClientState, err := relayerclient.DecodeClientState(clientStateBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to decode on-chain client state: %w", err)
	}
	log.Printf("[UpdateCosmosClient] On-chain client state: chainId=%s height=(%d,%d) frozen=%v",
		onChainClientState.ChainId, onChainClientState.LatestHeight.RevisionNumber,
		onChainClientState.LatestHeight.RevisionHeight, onChainClientState.IsFrozen)
	return int64(onChainClientState.LatestHeight.RevisionHeight), nil
}

type pinnedCosmosValidatorSet struct {
	indices      []uint32
	pubkeys      [][32]byte
	votingPowers []uint64
	totalPower   int64
}

func (s pinnedCosmosValidatorSet) indexByPubkey() map[[32]byte]uint32 {
	out := make(map[[32]byte]uint32, len(s.pubkeys))
	for i, pubkey := range s.pubkeys {
		out[pubkey] = s.indices[i]
	}
	return out
}

func (s pinnedCosmosValidatorSet) powerByPubkey() map[[32]byte]int64 {
	out := make(map[[32]byte]int64, len(s.pubkeys))
	for i, pubkey := range s.pubkeys {
		out[pubkey] = int64(s.votingPowers[i])
	}
	return out
}

func getPinnedCosmosValidatorSet(ctx Context) (pinnedCosmosValidatorSet, error) {
	ics07, err := tendermintContract.NewContractGroth16ICS07Tendermint(*ctx.ClientContract(), ctx.EthClient())
	if err != nil {
		return pinnedCosmosValidatorSet{}, fmt.Errorf("failed to create ICS07 instance: %w", err)
	}
	out, err := ics07.GetPinnedValidatorSet(nil)
	if err != nil {
		return pinnedCosmosValidatorSet{}, fmt.Errorf("getPinnedValidatorSet: %w", err)
	}
	if len(out.Indices) != len(out.Pubkeys) || len(out.Indices) != len(out.VotingPowers) {
		return pinnedCosmosValidatorSet{}, fmt.Errorf(
			"pinned validator set length mismatch: indices=%d pubkeys=%d powers=%d",
			len(out.Indices), len(out.Pubkeys), len(out.VotingPowers),
		)
	}
	var totalPower int64
	for i, idx := range out.Indices {
		if idx != uint32(i) {
			return pinnedCosmosValidatorSet{}, fmt.Errorf("pinned validator set index mismatch at position %d: got %d", i, idx)
		}
		if out.VotingPowers[i] > uint64(^uint64(0)>>1) {
			return pinnedCosmosValidatorSet{}, fmt.Errorf("pinned validator voting power exceeds int64: index=%d", i)
		}
		totalPower += int64(out.VotingPowers[i])
	}
	return pinnedCosmosValidatorSet{indices: out.Indices, pubkeys: out.Pubkeys, votingPowers: out.VotingPowers, totalPower: totalPower}, nil
}

func selectSignaturesForPinnedSet(
	candidates []prover.ValidatorSignature,
	pinnedSet pinnedCosmosValidatorSet,
) ([]prover.ValidatorSignature, error) {
	powerByPubkey := pinnedSet.powerByPubkey()
	selected := make([]prover.ValidatorSignature, 0, len(candidates))
	var pinnedAccum int64

	pinnedCandidates := make([]prover.ValidatorSignature, 0, len(candidates))
	for _, candidate := range candidates {
		if powerByPubkey[bytesToBytes32(candidate.PublicKey)] == 0 {
			continue
		}
		pinnedCandidates = append(pinnedCandidates, candidate)
	}
	sort.SliceStable(pinnedCandidates, func(i, j int) bool {
		leftPower := powerByPubkey[bytesToBytes32(pinnedCandidates[i].PublicKey)]
		rightPower := powerByPubkey[bytesToBytes32(pinnedCandidates[j].PublicKey)]
		if leftPower != rightPower {
			return leftPower > rightPower
		}
		return pinnedCandidates[i].Index < pinnedCandidates[j].Index
	})

	total := new(big.Int).SetInt64(pinnedSet.totalPower)
	total.Mul(total, big.NewInt(2))
	for _, candidate := range pinnedCandidates {
		selected = append(selected, candidate)
		pinnedAccum += powerByPubkey[bytesToBytes32(candidate.PublicKey)]
		accum := new(big.Int).SetInt64(pinnedAccum)
		accum.Mul(accum, big.NewInt(3))
		if accum.Cmp(total) > 0 {
			break
		}
	}
	accum := new(big.Int).SetInt64(pinnedAccum)
	accum.Mul(accum, big.NewInt(3))
	if accum.Cmp(total) <= 0 {
		return nil, fmt.Errorf("insufficient pinned-set quorum in proof candidates: have %d of %d", pinnedAccum, pinnedSet.totalPower)
	}
	if len(selected) > prover.MaxBucket() {
		return nil, fmt.Errorf("pinned-set proof requires %d signers but largest bucket is %d", len(selected), prover.MaxBucket())
	}

	sort.Slice(selected, func(i, j int) bool {
		return selected[i].Index < selected[j].Index
	})
	return selected, nil
}

// BuildCosmosClientUpdateMsg fetches the latest Tendermint light block,
// generates the Groth16 batch proof, and returns the resulting
// IUpdateClientMsgsMsgUpdateClient WITHOUT submitting it. Callers either
// pass the msg to SendEthTx directly (UpdateCosmosClient) or fold it into
// a multicall alongside packet calls (handleCosmos / V2).
//
// HasMsg=false signals the on-chain client is already at the latest block
// — no update needed; LightBlock still returned so callers can use it for
// membership proofs.
func (w *Worker) BuildCosmosClientUpdateMsg(ctx Context, proofType string, trustedBlock int64, trustLevel string) (*CosmosClientUpdateBuildResult, error) {
	status, err := ctx.CosmosClient().Status(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	log.Printf("[UpdateCosmosClient] called with trustedBlock=%d, latestBlockHeight=%d", trustedBlock, status.SyncInfo.LatestBlockHeight)

	// Always read the authoritative on-chain trusted height before deciding
	// whether to (re)generate the expensive Groth16 proof. The caller passes a
	// cached `trustedBlock` hint, but across chunked flushes — or after another
	// relayer / a crash-replay advanced the client — that cache can lag the real
	// chain state, causing redundant proofs for a range already on-chain
	// (issue #76 #2). The on-chain height always wins because only it is
	// guaranteed to identify a stored consensus state.
	onChainTrusted, err := fetchOnChainTrustedHeight(ctx)
	if err != nil {
		return nil, err
	}
	// Only the authoritative on-chain latest height is guaranteed to have a
	// stored consensus state. A cached hint can point at an unstored sparse
	// height and make _getConsensusStateHash revert with ConsensusStateNotFound.
	if trustedBlock != 0 && trustedBlock != onChainTrusted {
		log.Printf("[UpdateCosmosClient] cached trusted height %d differs from on-chain %d; using on-chain",
			trustedBlock, onChainTrusted)
	}
	trustedBlock = onChainTrusted
	if trustedBlock >= status.SyncInfo.LatestBlockHeight {
		if trustedBlock == status.SyncInfo.LatestBlockHeight {
			log.Printf("[UpdateCosmosClient] client is up to date (trusted=%d, latest=%d), skipping tx",
				trustedBlock, status.SyncInfo.LatestBlockHeight)
			lightBlock, err := relayerclient.GetLightBlock(ctx.CosmosClient(), trustedBlock)
			if err != nil {
				return nil, fmt.Errorf("failed to get current light block while up-to-date: %w", err)
			}
			return &CosmosClientUpdateBuildResult{LightBlock: lightBlock}, nil
		}
		return nil, fmt.Errorf("trusted block is ahead of latest chain height (trusted=%d, latest=%d)", trustedBlock, status.SyncInfo.LatestBlockHeight)
	}

	log.Printf("[UpdateCosmosClient] Fetching trustedLightBlock at height %d, latestLightBlock at height %d", trustedBlock, status.SyncInfo.LatestBlockHeight)
	trustedLightBlock, err := relayerclient.GetLightBlock(ctx.CosmosClient(), trustedBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted light block: %w", err)
	}

	latestLightBlock, err := relayerclient.GetLightBlock(ctx.CosmosClient(), status.SyncInfo.LatestBlockHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest light block: %w", err)
	}

	unbondingPeriod, err := relayerclient.GetUnbondingTime(ctx.CosmosClient())
	if err != nil {
		return nil, fmt.Errorf("failed to get unbonding time: %w", err)
	}

	trustingPeriod := uint32(unbondingPeriod * 2 / 3)
	if ctx.Config.TrustingPeriod != 0 {
		trustingPeriod = ctx.Config.TrustingPeriod
	}

	if trustingPeriod > uint32(unbondingPeriod) {
		return nil, fmt.Errorf("trusting period %d cannot be greater than unbonding period %d", trustingPeriod, uint32(unbondingPeriod))
	}

	chainId := trustedLightBlock.SignedHeader.Header.ChainID
	revision := clienttypes.ParseChainID(chainId)

	trustThreshold, err := relayerclient.ParseTrustThreshold(trustLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trust level: %w", err)
	}

	var zkAlgorithm relayerclient.SupportedZkAlgorithm
	switch proofType {
	case "groth16":
		zkAlgorithm = relayerclient.Groth16
	case "plonk":
		zkAlgorithm = relayerclient.Plonk
	default:
		return nil, fmt.Errorf("unsupported proof type: %s, supported types are: groth16, plonk", proofType)
	}

	if trustedLightBlock.SignedHeader.Header.Height < 0 {
		return nil, fmt.Errorf("trusted light block header height cannot be negative: %d", trustedLightBlock.SignedHeader.Header.Height)
	}

	clientState := updateclientContract.IICS07TendermintMsgsClientState{
		ChainId:    chainId,
		TrustLevel: trustThreshold,
		LatestHeight: updateclientContract.IICS02ClientMsgsHeight{
			RevisionNumber: revision,
			RevisionHeight: uint64(trustedLightBlock.SignedHeader.Header.Height),
		},
		IsFrozen:        false,
		ZkAlgorithm:     uint8(zkAlgorithm),
		TrustingPeriod:  trustingPeriod,
		UnbondingPeriod: uint32(unbondingPeriod),
		ClockDrift:      15,
	}

	consensusState := updateclientContract.IICS07TendermintMsgsConsensusState{
		Timestamp:          big.NewInt(trustedLightBlock.SignedHeader.Header.Time.UnixNano()),
		Root:               bytesToBytes32(trustedLightBlock.SignedHeader.Header.AppHash),
		NextValidatorsHash: bytesToBytes32(trustedLightBlock.SignedHeader.NextValidatorsHash),
	}

	proposedHeader, err := latestLightBlock.IntoHeader(*trustedLightBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to convert light block into header: %w", err)
	}

	pinnedValidatorSet, err := getPinnedCosmosValidatorSet(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query pinned validator set: %w", err)
	}

	log.Printf("[UpdateCosmosClient] clientState.LatestHeight=(%d,%d) proposedHeader.Height=%d trustedBlock=%d latestBlock=%d",
		clientState.LatestHeight.RevisionNumber, clientState.LatestHeight.RevisionHeight,
		proposedHeader.SignedHeader.Header.Height, trustedLightBlock.BlockHeight, latestLightBlock.BlockHeight)

	// Extract non-absent validator signatures, then select enough signers that
	// overlap the pinned validator set to exceed 2/3 of pinned voting power.
	extracted, err := prover.ExtractValidatorSignatures(latestLightBlock, chainId, nil)
	if err != nil {
		return nil, fmt.Errorf("extract validator signatures: %w", err)
	}
	selected, err := selectSignaturesForPinnedSet(extracted.Candidates, pinnedValidatorSet)
	if err != nil {
		return nil, fmt.Errorf("select pinned-set signatures: %w", err)
	}
	extracted.Signatures = selected
	log.Printf("[UpdateCosmosClient] Generating Groth16 batch proof for %d validator signatures...", len(extracted.Signatures))
	bucket, paddedSigs, proof, commitments, commitmentPok, err := w.Prover.GenerateProof(extracted.Signatures)
	if err != nil {
		return nil, fmt.Errorf("error generating proof: %w", err)
	}
	log.Printf("[UpdateCosmosClient] Proof generated (bucket=%d). Sending Eth tx...", bucket)

	// paddedSigs is the prover's bucket-sized slice (real signers + dummy
	// padding); the on-chain quorum check sees the same per-slot layout the
	// circuit committed to. Padding slots have Active=false and zero voting
	// power; Solidity skips them via the active flag.
	signerIndices := make([]uint32, bucket)
	signerPubkeys := make([][32]byte, bucket)
	active := make([]bool, bucket)
	pinnedValidatorIndices := make([]uint32, bucket)
	pinnedIndexByPubkey := pinnedValidatorSet.indexByPubkey()
	for i, s := range paddedSigs {
		signerIndices[i] = uint32(s.Index)
		if s.Active {
			pinnedIdx, ok := pinnedIndexByPubkey[bytesToBytes32(s.PublicKey)]
			if !ok {
				return nil, fmt.Errorf("active signer at slot %d is not in pinned validator set", i)
			}
			pinnedValidatorIndices[i] = pinnedIdx
		}
		copy(signerPubkeys[i][:], s.PublicKey)
		active[i] = s.Active
	}

	msg := updateclientContract.IUpdateClientMsgsMsgUpdateClient{
		ClientState:            clientState,
		TrustedConsensusState:  consensusState,
		Time:                   big.NewInt(time.Now().UnixNano()),
		ProposedHeader:         proposedHeader,
		Proof:                  proof,
		Commitments:            commitments,
		CommitmentPok:          commitmentPok,
		Bucket:                 uint16(bucket),
		SignerIndices:          signerIndices,
		PinnedValidatorIndices: pinnedValidatorIndices,
		SignerPubkeys:          signerPubkeys,
		Active:                 active,
	}

	return &CosmosClientUpdateBuildResult{
		Msg:        msg,
		HasMsg:     true,
		LightBlock: latestLightBlock,
	}, nil
}

func (w *Worker) CreateEthClient(ctx Context, checksum string) (string, error) {
	beaconAPIURL := ctx.BeaconAPIURL()
	if beaconAPIURL == "" {
		return "", fmt.Errorf("beacon API URL is not configured")
	}
	log.Printf("[CreateEthClient] starting: beacon=%s checksum=%s", beaconAPIURL, checksum)

	log.Printf("[CreateEthClient] fetching beacon genesis")
	bctx, bcancel := context.WithTimeout(context.Background(), 15*time.Second)
	genesis, err := relayerclient.GetBeaconGenesis(bctx, beaconAPIURL)
	bcancel()
	if err != nil {
		return "", fmt.Errorf("failed to get light client genesis: %w", err)
	}
	log.Printf("[CreateEthClient] beacon genesis fetched: genesisTime=%s genesisValidatorsRoot=%s", genesis.GenesisTime, genesis.GenesisValidatorsRoot)

	log.Printf("[CreateEthClient] fetching beacon spec")
	bctx, bcancel = context.WithTimeout(context.Background(), 15*time.Second)
	spec, err := relayerclient.GetBeaconSpec(bctx, beaconAPIURL)
	bcancel()
	if err != nil {
		return "", fmt.Errorf("failed to get light client spec: %w", err)
	}
	log.Printf("[CreateEthClient] beacon spec fetched: secondsPerSlot=%s slotsPerEpoch=%s syncCommitteeSize=%s",
		spec.SecondsPerSlot, spec.SlotsPerEpoch, spec.SyncCommitteeSize)
	// Use the finalized header from the finality update — this is always a checkpoint slot
	// (epoch boundary), unlike GetBeaconBlock("finalized") which may return a non-checkpoint slot.
	log.Printf("[CreateEthClient] fetching finality update")
	bctx, bcancel = context.WithTimeout(context.Background(), 15*time.Second)
	finalityUpdate, err := relayerclient.GetFinalityUpdate(bctx, beaconAPIURL)
	bcancel()
	if err != nil {
		return "", fmt.Errorf("failed to get finality update: %w", err)
	}
	checkpointSlot := finalityUpdate.FinalizedHeader.Beacon.Slot
	log.Printf("[CreateEthClient] finality update fetched: attestedSlot=%s finalizedSlot=%s signatureSlot=%s",
		finalityUpdate.AttestedHeader.Beacon.Slot, checkpointSlot, finalityUpdate.SignatureSlot)

	log.Printf("[CreateEthClient] fetching beacon block root for slot=%s", checkpointSlot)
	bctx, bcancel = context.WithTimeout(context.Background(), 15*time.Second)
	blockRoot, err := relayerclient.GetBeaconBlockRoot(bctx, beaconAPIURL, checkpointSlot)
	bcancel()
	if err != nil {
		return "", fmt.Errorf("failed to get beacon block root: %w", err)
	}
	log.Printf("[CreateEthClient] beacon block root=%s", blockRoot)

	log.Printf("[CreateEthClient] fetching light client bootstrap")
	bctx, bcancel = context.WithTimeout(context.Background(), 15*time.Second)
	bootstrap, err := relayerclient.GetLightClientBootstrap(bctx, beaconAPIURL, blockRoot)
	bcancel()
	if err != nil {
		return "", fmt.Errorf("failed to get light client bootstrap: %w", err)
	}
	log.Printf("[CreateEthClient] checkpointSlot=%s syncCommittee.AggregatePubkey=%s",
		checkpointSlot, bootstrap.Data.CurrentSyncCommittee.AggregatePubkey)

	log.Printf("[CreateEthClient] fetching beacon block for slot=%s", checkpointSlot)
	bctx, bcancel = context.WithTimeout(context.Background(), 15*time.Second)
	beaconBlock, err := relayerclient.GetBeaconBlock(bctx, beaconAPIURL, checkpointSlot)
	bcancel()
	if err != nil {
		return "", fmt.Errorf("failed to get beacon block: %w", err)
	}
	log.Printf("[CreateEthClient] beacon block fetched: executionBlock=%s", beaconBlock.Message.Body.ExecutionPayload.BlockNumber)

	if bootstrap.Data.Header.Execution.BlockNumber != beaconBlock.Message.Body.ExecutionPayload.BlockNumber {
		return "", fmt.Errorf("light client bootstrap block number does not match execution block number")
	}

	log.Printf("[CreateEthClient] querying ethereum chain id")
	chainId, err := ctx.EthClient().ChainID(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get eth chain id: %w", err)
	}
	log.Printf("[CreateEthClient] ethereum chain id=%s", chainId.String())

	forkParameters, err := spec.ToForkParameters()
	if err != nil {
		return "", fmt.Errorf("failed to get fork parameters: %w", err)
	}

	epochsPerSyncCommitteePeriod, err := strconv.ParseUint(spec.EpochsPerSyncCommitteePeriod, 10, 64)
	if err != nil {
		return "", err
	}

	genesisTime, err := strconv.ParseUint(genesis.GenesisTime, 10, 64)
	if err != nil {
		return "", err
	}

	blockNumber, err := strconv.ParseUint(bootstrap.Data.Header.Execution.BlockNumber, 10, 64)
	if err != nil {
		return "", err
	}
	slot, err := strconv.ParseUint(bootstrap.Data.Header.Beacon.Slot, 10, 64)
	if err != nil {
		return "", err
	}

	syncCommitteeSize, err := strconv.ParseUint(spec.SyncCommitteeSize, 10, 64)
	if err != nil {
		return "", err
	}

	secondsPerSlot, err := strconv.ParseUint(spec.SecondsPerSlot, 10, 64)
	if err != nil {
		return "", err
	}
	slotsPerEpoch, err := strconv.ParseUint(spec.SlotsPerEpoch, 10, 64)
	if err != nil {
		return "", err
	}
	clientState := relayerclient.EthereumClientState{
		ChainID:                      chainId.Uint64(),
		EpochsPerSyncCommitteePeriod: epochsPerSyncCommitteePeriod,
		ForkParameters:               *forkParameters,
		GenesisSlot:                  0,
		GenesisTime:                  genesisTime,
		GenesisValidatorsRoot:        genesis.GenesisValidatorsRoot,
		IbcCommitmentSlot:            ICS26_IBC_STORAGE_SLOT,
		IbcContractAddress:           ctx.RouterContract().String(),
		IsFrozen:                     false,
		LatestExecutionBlockNumber:   blockNumber,
		LatestSlot:                   slot,
		MinSyncCommitteeParticipants: (syncCommitteeSize + 2) / 3,
		SecondsPerSlot:               secondsPerSlot,
		SlotsPerEpoch:                slotsPerEpoch,
		SyncCommitteeSize:            syncCommitteeSize,
	}
	log.Printf("[CreateEthClient] clientState prepared: latestSlot=%d latestExecutionBlock=%d minSyncCommitteeParticipants=%d",
		clientState.LatestSlot, clientState.LatestExecutionBlockNumber, clientState.MinSyncCommitteeParticipants)
	clientStateBz, err := json.Marshal(clientState)
	if err != nil {
		return "", fmt.Errorf("error serializing client state: %w", err)
	}
	checksumTrimmed := strings.TrimPrefix(checksum, "0x")
	checksumBz, err := hex.DecodeString(checksumTrimmed)
	if err != nil {
		return "", fmt.Errorf("error parsing checksum: %w", err)
	}
	wasmClientState := ibcwasmtypes.ClientState{
		Data:     clientStateBz,
		Checksum: checksumBz,
		LatestHeight: clienttypes.Height{
			RevisionNumber: 0,
			RevisionHeight: clientState.LatestSlot,
		},
	}

	timestamp, err := strconv.ParseUint(bootstrap.Data.Header.Execution.Timestamp, 10, 64)
	if err != nil {
		return "", err
	}

	currentSyncCommittee, err := bootstrap.Data.CurrentSyncCommittee.ToSummarizedSyncCommittee()
	if err != nil {
		return "", fmt.Errorf("failed to hash pubkeys for CurrentSyncCommittee: %w", err)
	}

	latestPeriod := clientState.ComputeSyncCommitteePeriodAtSlot(clientState.LatestSlot)
	log.Printf("[CreateEthClient] fetching light client updates for latestPeriod=%d", latestPeriod)
	bctx, bcancel = context.WithTimeout(context.Background(), 15*time.Second)
	lightClientUpdates, err := relayerclient.GetLightClientUpdates(bctx, ctx.BeaconAPIURL(), latestPeriod, 1)
	bcancel()
	if err != nil {
		return "", fmt.Errorf("failed to get light client updates: %w", err)
	}
	log.Printf("[CreateEthClient] fetched %d light client update(s)", len(lightClientUpdates))
	if len(lightClientUpdates) == 0 {
		return "", fmt.Errorf("no light client updates returned for period %d", latestPeriod)
	}

	nextSyncCommittee, err := lightClientUpdates[0].NextSyncCommittee.ToSummarizedSyncCommittee()
	if err != nil {
		return "", fmt.Errorf("failed to hash pubkeys for NextSyncCommittee: %w", err)
	}

	consensusState := relayerclient.EthereumConsensusState{
		Slot:                 clientState.LatestSlot,
		StateRoot:            bootstrap.Data.Header.Execution.StateRoot,
		Timestamp:            timestamp,
		CurrentSyncCommittee: *currentSyncCommittee,
		NextSyncCommittee:    nextSyncCommittee,
	}
	consensusStateBz, err := json.Marshal(consensusState)
	if err != nil {
		return "", fmt.Errorf("error serializing consensus state: %w", err)
	}

	wasmConsensusState := ibcwasmtypes.ConsensusState{
		Data: consensusStateBz,
	}
	log.Printf("[CreateEthClient] wasm client/consensus state prepared, broadcasting MsgCreateClient")

	return w.TxHandler.CreateEthClient(ctx, &wasmClientState, &wasmConsensusState)
}

type EthClientUpdateResult struct {
	Msgs           []any
	EthClientState *relayerclient.EthereumClientState
	ProofTimestamp uint64
	SigSlot        uint64
}

func (w *Worker) UpdateEthClient(ctx Context) error {
	result, err := w.BuildEthClientUpdateMsgs(ctx)
	if err != nil {
		return err
	}
	if len(result.Msgs) == 0 {
		return nil
	}
	w.waitForCosmosCatchUp(ctx, result.EthClientState, result.SigSlot)
	return w.TxHandler.SendCosmosTxBatch(ctx, result.Msgs)
}

func (w *Worker) BuildEthClientUpdateMsgs(ctx Context) (*EthClientUpdateResult, error) {
	beaconAPIURL := ctx.BeaconAPIURL()
	if beaconAPIURL == "" {
		return nil, fmt.Errorf("beacon API URL is not configured")
	}

	ethClientID := ctx.EthClientID()
	if ethClientID == "" {
		return nil, fmt.Errorf("ethereum client ID is not configured")
	}

	ethClientState, err := relayerclient.GetEthereumClientState(ctx.CosmosClient(), ethClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ethereum client state: %w", err)
	}
	trustedSlot := ethClientState.LatestSlot

	bctx, bcancel := context.WithTimeout(context.Background(), 15*time.Second)
	finalityUpdate, err := relayerclient.GetFinalityUpdate(bctx, beaconAPIURL)
	bcancel()
	if err != nil {
		return nil, fmt.Errorf("failed to get finality update: %w", err)
	}

	participation := relayerclient.CountSyncCommitteeParticipants(finalityUpdate.SyncAggregate.SyncCommitteeBits)
	syncCommitteeSize := ethClientState.SyncCommitteeSize
	log.Printf("[UpdateEthClient] sync committee participation: %d/%d (%.1f%%)",
		participation, syncCommitteeSize, float64(participation)*100/float64(syncCommitteeSize))
	if participation*3 < syncCommitteeSize*2 {
		return nil, fmt.Errorf("insufficient sync committee participation: %d/%d (need 2/3 = %d)", participation, syncCommitteeSize, syncCommitteeSize*2/3)
	}

	finalizedSlot, err := parseSlot(finalityUpdate.FinalizedHeader.Beacon.Slot)
	if err != nil {
		return nil, fmt.Errorf("failed to parse finalized slot: %w", err)
	}

	log.Printf("[UpdateEthClient] trustedSlot=%d finalizedSlot=%d latestExecBlock=%d",
		trustedSlot, finalizedSlot, ethClientState.LatestExecutionBlockNumber)

	if finalizedSlot <= trustedSlot {
		log.Printf("[UpdateEthClient] already up to date, skipping")
		return &EthClientUpdateResult{
			Msgs:           nil,
			EthClientState: cloneEthereumClientState(ethClientState),
			ProofTimestamp: ethClientState.ComputeTimestampAtSlot(trustedSlot),
		}, nil
	}

	trustedPeriod := ethClientState.ComputeSyncCommitteePeriodAtSlot(trustedSlot)
	targetPeriod := ethClientState.ComputeSyncCommitteePeriodAtSlot(finalizedSlot)

	log.Printf("[UpdateEthClient] trustedPeriod=%d targetPeriod=%d", trustedPeriod, targetPeriod)

	msgs, err := w.buildEthClientUpdateMsgsWithPeriodCrossing(ctx, beaconAPIURL, ethClientID, ethClientState, trustedSlot, trustedPeriod, targetPeriod, finalityUpdate, finalizedSlot)
	if err != nil {
		return nil, err
	}

	proofState, proofTimestamp, err := ethProofStateFromFinalityUpdate(ethClientState, finalityUpdate, finalizedSlot)
	if err != nil {
		return nil, err
	}

	sigSlot, _ := parseSlot(finalityUpdate.SignatureSlot)
	return &EthClientUpdateResult{
		Msgs:           msgs,
		EthClientState: proofState,
		ProofTimestamp: proofTimestamp,
		SigSlot:        sigSlot,
	}, nil
}

func (w *Worker) waitForCosmosCatchUp(ctx Context, ethClientState *relayerclient.EthereumClientState, sigSlot uint64) {
	requiredSlot := sigSlot + cosmosCatchUpSafetySlots
	for range 60 {
		status, err := ctx.CosmosClient().Status(context.Background())
		if err != nil {
			break
		}
		cosmosTime := uint64(status.SyncInfo.LatestBlockTime.Unix())
		currentSlot := ethClientState.ComputeSlotAtTimestamp(cosmosTime)
		if cosmosCurrentSlotReady(currentSlot, sigSlot) {
			log.Printf("[updateEthClient] timing OK: currentSlot=%d >= requiredSlot=%d (signatureSlot=%d safety=%d)",
				currentSlot, requiredSlot, sigSlot, cosmosCatchUpSafetySlots)
			break
		}
		log.Printf("[updateEthClient] waiting for target chain to catch up to required slot %d (signatureSlot=%d current=%d safety=%d)",
			requiredSlot, sigSlot, currentSlot, cosmosCatchUpSafetySlots)
		time.Sleep(5 * time.Second)
	}
}

func cosmosCurrentSlotReady(currentSlot, sigSlot uint64) bool {
	return currentSlot >= sigSlot+cosmosCatchUpSafetySlots
}

func (w *Worker) buildEthClientUpdateMsgsWithPeriodCrossing(ctx Context, beaconAPIURL, ethClientID string, ethClientState *relayerclient.EthereumClientState, trustedSlot, trustedPeriod, targetPeriod uint64, finalityUpdate *relayerclient.LightClientFinalityUpdate, finalizedSlot uint64) ([]any, error) {
	count := targetPeriod - trustedPeriod + 1
	bctx, bcancel := context.WithTimeout(context.Background(), 15*time.Second)
	lightClientUpdates, err := relayerclient.GetLightClientUpdates(bctx, beaconAPIURL, trustedPeriod, count)
	bcancel()
	if err != nil {
		return nil, fmt.Errorf("failed to get light client updates: %w", err)
	}

	if len(lightClientUpdates) == 0 {
		return nil, fmt.Errorf("no light client updates available for period range %d to %d", trustedPeriod, targetPeriod)
	}

	var msgs []any
	latestTrustedSlot := trustedSlot
	latestPeriod := trustedPeriod

	for _, update := range lightClientUpdates {
		updateFinalizedSlot, err := parseSlot(update.FinalizedHeader.Beacon.Slot)
		if err != nil {
			return nil, fmt.Errorf("failed to parse update finalized slot: %w", err)
		}

		if updateFinalizedSlot <= latestTrustedSlot {
			continue
		}

		updatePeriod := ethClientState.ComputeSyncCommitteePeriodAtSlot(updateFinalizedSlot)
		if updatePeriod == latestPeriod {
			continue
		}

		bctx, bcancel = context.WithTimeout(context.Background(), 15*time.Second)
		blockRoot, err := relayerclient.GetBeaconBlockRoot(bctx, beaconAPIURL, fmt.Sprintf("%d", updateFinalizedSlot))
		bcancel()
		if err != nil {
			return nil, fmt.Errorf("failed to get beacon block root: %w", err)
		}

		bctx, bcancel = context.WithTimeout(context.Background(), 15*time.Second)
		bootstrap, err := relayerclient.GetLightClientBootstrap(bctx, beaconAPIURL, blockRoot)
		bcancel()
		if err != nil {
			return nil, fmt.Errorf("failed to get light client bootstrap: %w", err)
		}

		syncCommittee := bootstrap.Data.CurrentSyncCommittee

		header := relayerclient.EthereumHeader{
			ActiveSyncCommittee: relayerclient.ActiveSyncCommittee{
				Next: &syncCommittee,
			},
			ConsensusUpdate: update,
			TrustedSlot:     latestTrustedSlot,
		}

		msg, err := buildMsgUpdateClient("", ethClientID, header)
		if err != nil {
			return nil, fmt.Errorf("failed to build update client message: %w", err)
		}

		msgs = append(msgs, msg)
		latestPeriod = updatePeriod
		latestTrustedSlot = updateFinalizedSlot
	}

	// If the latest header is earlier than the finality update, add a header for the finality update.
	if finalizedSlot > latestTrustedSlot {
		attestedSlot := finalityUpdate.AttestedHeader.Beacon.Slot
		log.Printf("[updateEthClient] final update: attestedSlot=%s finalizedSlot=%d latestTrustedSlot=%d",
			attestedSlot, finalizedSlot, latestTrustedSlot)

		// Get sync committee from attested slot's bootstrap (matches eureka relayer behavior)
		bctx, bcancel = context.WithTimeout(context.Background(), 15*time.Second)
		blockRoot, err := relayerclient.GetBeaconBlockRoot(bctx, beaconAPIURL, attestedSlot)
		bcancel()
		if err != nil {
			return nil, fmt.Errorf("failed to get beacon block root: %w", err)
		}

		bctx, bcancel = context.WithTimeout(context.Background(), 15*time.Second)
		bootstrap, err := relayerclient.GetLightClientBootstrap(bctx, beaconAPIURL, blockRoot)
		bcancel()
		if err != nil {
			return nil, fmt.Errorf("failed to get light client bootstrap: %w", err)
		}

		syncCommittee := bootstrap.Data.CurrentSyncCommittee

		consensusUpdate := relayerclient.LightClientUpdate{
			AttestedHeader:          finalityUpdate.AttestedHeader,
			NextSyncCommittee:       nil,
			NextSyncCommitteeBranch: nil,
			FinalizedHeader:         finalityUpdate.FinalizedHeader,
			FinalityBranch:          finalityUpdate.FinalityBranch,
			SyncAggregate:           finalityUpdate.SyncAggregate,
			SignatureSlot:           finalityUpdate.SignatureSlot,
		}

		header := relayerclient.EthereumHeader{
			ActiveSyncCommittee: relayerclient.ActiveSyncCommittee{
				Current: &syncCommittee,
			},
			ConsensusUpdate: consensusUpdate,
			TrustedSlot:     latestTrustedSlot,
		}

		msg, err := buildMsgUpdateClient("", ethClientID, header)
		if err != nil {
			return nil, fmt.Errorf("failed to build update client message: %w", err)
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func cloneEthereumClientState(state *relayerclient.EthereumClientState) *relayerclient.EthereumClientState {
	if state == nil {
		return nil
	}
	cloned := *state
	return &cloned
}

func ethProofStateFromFinalityUpdate(base *relayerclient.EthereumClientState, finalityUpdate *relayerclient.LightClientFinalityUpdate, finalizedSlot uint64) (*relayerclient.EthereumClientState, uint64, error) {
	proofState := cloneEthereumClientState(base)
	if proofState == nil {
		return nil, 0, fmt.Errorf("ethereum client state is nil")
	}

	finalizedExecutionBlock, err := strconv.ParseUint(finalityUpdate.FinalizedHeader.Execution.BlockNumber, 10, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse finalized execution block number: %w", err)
	}
	finalizedTimestamp, err := strconv.ParseUint(finalityUpdate.FinalizedHeader.Execution.Timestamp, 10, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse finalized execution timestamp: %w", err)
	}

	proofState.LatestSlot = finalizedSlot
	proofState.LatestExecutionBlockNumber = finalizedExecutionBlock
	return proofState, finalizedTimestamp, nil
}

// parseSlot parses a slot string to uint64
func parseSlot(slotStr string) (uint64, error) {
	var slot uint64
	_, err := fmt.Sscanf(slotStr, "%d", &slot)
	return slot, err
}

func bytesToBytes32(data []byte) [32]byte {
	var result [32]byte
	copy(result[:], data)
	return result
}

// buildMsgUpdateClient builds a MsgUpdateClient for the Ethereum light client
func buildMsgUpdateClient(signerAddr string, clientID string, header relayerclient.EthereumHeader) (*clienttypes.MsgUpdateClient, error) {
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal header: %w", err)
	}

	clientMessage := &ibcwasmtypes.ClientMessage{
		Data: headerBytes,
	}

	clientMessageAny, err := codectypes.NewAnyWithValue(clientMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to create Any for client message: %w", err)
	}

	return &clienttypes.MsgUpdateClient{
		ClientId:      clientID,
		ClientMessage: clientMessageAny,
		Signer:        signerAddr,
	}, nil
}
