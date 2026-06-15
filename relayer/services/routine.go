package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
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
const noTrustedOverlapIndex = ^uint32(0)

type Worker struct {
	TxHandler TransactionHandler
	Prover    Prover
}

const cosmosCatchUpSafetySlots uint64 = 3
const maxValidatorDeltaLeafCount = 16

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
	return w.TxHandler.CreateCosmosClientContract(ctx, clientStateEncoded, consensusHash)
}

// CosmosClientUpdateBuildResult is the output of BuildCosmosClientUpdateMsg.
// When HasMsg is false the on-chain client is already at the latest height
// and no updateClient tx is needed — only LightBlock is populated.
type CosmosClientUpdateBuildResult struct {
	Msg                         updateclientContract.IUpdateClientMsgsMsgUpdateClient
	HasMsg                      bool
	LightBlock                  *relayerclient.LightBlock
	UsedValidatorCache          bool
	FullValidatorSetFallbackMsg updateclientContract.IUpdateClientMsgsMsgUpdateClient
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
		if errors.Is(err, ErrValidatorCacheRace) && result.UsedValidatorCache {
			log.Printf("[UpdateCosmosClient] validator cache changed during submission; retrying updateClient with full validator set")
			if retryErr := w.TxHandler.SendEthTx(ctx, result.FullValidatorSetFallbackMsg); retryErr != nil {
				log.Printf("[UpdateCosmosClient] full validator-set retry failed: %v", retryErr)
				return nil, retryErr
			}
			log.Printf("[UpdateCosmosClient] full validator-set retry succeeded")
			return result.LightBlock, nil
		}
		return nil, err
	}
	log.Printf("[UpdateCosmosClient] SendEthTx succeeded")
	return result.LightBlock, nil
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

type cachedCosmosValidatorSet struct {
	indices      []uint32
	pubkeys      [][32]byte
	votingPowers []uint64
}

func (s cachedCosmosValidatorSet) isEmpty() bool {
	return len(s.indices) == 0
}

// getCachedCosmosValidatorSet fetches the on-chain validator cache snapshot for
// a validatorsHash. The contract returns an empty snapshot when the hash is not
// cached or the stored data is unusable.
func getCachedCosmosValidatorSet(ctx Context, validatorsHash [32]byte) (cachedCosmosValidatorSet, error) {
	ics07, err := tendermintContract.NewContractGroth16ICS07Tendermint(*ctx.ClientContract(), ctx.EthClient())
	if err != nil {
		return cachedCosmosValidatorSet{}, fmt.Errorf("failed to create ICS07 instance: %w", err)
	}
	out, err := ics07.GetCachedValidatorSet(nil, validatorsHash)
	if err != nil {
		return cachedCosmosValidatorSet{}, fmt.Errorf("getCachedValidatorSet(%x): %w", validatorsHash, err)
	}
	if len(out.Indices) != len(out.Pubkeys) || len(out.Indices) != len(out.VotingPowers) {
		return cachedCosmosValidatorSet{}, fmt.Errorf(
			"cached validator set length mismatch for hash=%x: indices=%d pubkeys=%d powers=%d",
			validatorsHash, len(out.Indices), len(out.Pubkeys), len(out.VotingPowers),
		)
	}
	for i, idx := range out.Indices {
		if idx != uint32(i) {
			return cachedCosmosValidatorSet{}, fmt.Errorf(
				"cached validator set index mismatch for hash=%x at position %d: got %d",
				validatorsHash, i, idx,
			)
		}
	}
	return cachedCosmosValidatorSet{
		indices:      out.Indices,
		pubkeys:      out.Pubkeys,
		votingPowers: out.VotingPowers,
	}, nil
}

func detectValidatorSetDelta(
	baseHash [32]byte,
	baseSet cachedCosmosValidatorSet,
	currentSet updateclientContract.IICS07TendermintMsgsValidatorSet,
) (updateclientContract.IUpdateClientMsgsValidatorSetDelta, bool, string) {
	if baseSet.isEmpty() {
		return updateclientContract.IUpdateClientMsgsValidatorSetDelta{}, false, "base validator set is not cached"
	}
	if len(baseSet.pubkeys) != len(currentSet.Validators) {
		return updateclientContract.IUpdateClientMsgsValidatorSetDelta{}, false, fmt.Sprintf(
			"validator count changed: base=%d current=%d",
			len(baseSet.pubkeys), len(currentSet.Validators),
		)
	}

	delta := updateclientContract.IUpdateClientMsgsValidatorSetDelta{
		BaseValidatorsHash: baseHash,
	}
	changeCount := 0
	for i, current := range currentSet.Validators {
		if baseSet.pubkeys[i] == current.PubKey && baseSet.votingPowers[i] == current.VotingPower {
			continue
		}
		if changeCount >= maxValidatorDeltaLeafCount {
			return updateclientContract.IUpdateClientMsgsValidatorSetDelta{}, false,
				fmt.Sprintf("more than %d validator leaves changed", maxValidatorDeltaLeafCount)
		}
		delta.Indices[changeCount] = uint32(i)
		delta.PubKeys[changeCount] = current.PubKey
		delta.VotingPowers[changeCount] = current.VotingPower
		changeCount++
	}

	if changeCount == 0 {
		return updateclientContract.IUpdateClientMsgsValidatorSetDelta{}, false, "no validator leaf change detected"
	}
	delta.LeafCount = uint8(changeCount)
	return delta, true, ""
}

func buildTrustedOverlapIndices(
	trustedNextValidatorSet updateclientContract.IICS07TendermintMsgsValidatorSet,
	paddedSigs []prover.ValidatorSignature,
	bucket int,
) []uint32 {
	trustedByPubkey := make(map[[32]byte]uint32, len(trustedNextValidatorSet.Validators))
	for i, validator := range trustedNextValidatorSet.Validators {
		trustedByPubkey[validator.PubKey] = uint32(i)
	}

	indices := make([]uint32, bucket)
	for i := range indices {
		indices[i] = noTrustedOverlapIndex
	}
	for i, sig := range paddedSigs {
		if i >= len(indices) || !sig.Active {
			continue
		}
		pubkey := bytesToBytes32(sig.PublicKey)
		if trustedIndex, ok := trustedByPubkey[pubkey]; ok {
			indices[i] = trustedIndex
		}
	}
	return indices
}

func selectSignaturesForTrustedOverlap(
	candidates []prover.ValidatorSignature,
	currentTotalPower int64,
	trustedNextValidatorSet updateclientContract.IICS07TendermintMsgsValidatorSet,
	trustLevel updateclientContract.IICS07TendermintMsgsTrustThreshold,
) ([]prover.ValidatorSignature, error) {
	trustedPowerByPubkey := make(map[[32]byte]int64, len(trustedNextValidatorSet.Validators))
	var trustedTotalPower int64
	for _, validator := range trustedNextValidatorSet.Validators {
		power := int64(validator.VotingPower)
		trustedTotalPower += power
		trustedPowerByPubkey[validator.PubKey] = power
	}

	selected := make([]prover.ValidatorSignature, 0, len(candidates))
	selectedByIndex := make(map[int]bool, len(candidates))
	var currentAccum int64
	var trustedAccum int64

	overlapCandidates := append([]prover.ValidatorSignature(nil), candidates...)
	sort.SliceStable(overlapCandidates, func(i, j int) bool {
		leftPower := trustedPowerByPubkey[bytesToBytes32(overlapCandidates[i].PublicKey)]
		rightPower := trustedPowerByPubkey[bytesToBytes32(overlapCandidates[j].PublicKey)]
		if leftPower != rightPower {
			return leftPower > rightPower
		}
		if overlapCandidates[i].Power != overlapCandidates[j].Power {
			return overlapCandidates[i].Power > overlapCandidates[j].Power
		}
		return overlapCandidates[i].Index < overlapCandidates[j].Index
	})

	for _, candidate := range overlapCandidates {
		trustedPower := trustedPowerByPubkey[bytesToBytes32(candidate.PublicKey)]
		if trustedPower == 0 {
			continue
		}
		selected = append(selected, candidate)
		selectedByIndex[candidate.Index] = true
		currentAccum += candidate.Power
		trustedAccum += trustedPower
		if meetsTrustThreshold(trustedAccum, trustedTotalPower, int64(trustLevel.Numerator), int64(trustLevel.Denominator)) {
			break
		}
	}
	if !meetsTrustThreshold(trustedAccum, trustedTotalPower, int64(trustLevel.Numerator), int64(trustLevel.Denominator)) {
		return nil, fmt.Errorf("insufficient trusted overlap in proof candidates: have %d of %d", trustedAccum, trustedTotalPower)
	}

	currentCandidates := append([]prover.ValidatorSignature(nil), candidates...)
	sort.SliceStable(currentCandidates, func(i, j int) bool {
		if currentCandidates[i].Power != currentCandidates[j].Power {
			return currentCandidates[i].Power > currentCandidates[j].Power
		}
		return currentCandidates[i].Index < currentCandidates[j].Index
	})
	for _, candidate := range currentCandidates {
		if meetsTrustThreshold(currentAccum, currentTotalPower, 2, 3) {
			break
		}
		if selectedByIndex[candidate.Index] {
			continue
		}
		selected = append(selected, candidate)
		selectedByIndex[candidate.Index] = true
		currentAccum += candidate.Power
	}
	if !meetsTrustThreshold(currentAccum, currentTotalPower, 2, 3) {
		return nil, fmt.Errorf("insufficient current quorum in proof candidates: have %d of %d", currentAccum, currentTotalPower)
	}
	if len(selected) > prover.MaxBucket() {
		return nil, fmt.Errorf("trusted-overlap proof requires %d signers but largest bucket is %d", len(selected), prover.MaxBucket())
	}

	sort.Slice(selected, func(i, j int) bool {
		return selected[i].Index < selected[j].Index
	})
	return selected, nil
}

func meetsTrustThreshold(accumulated, total, numerator, denominator int64) bool {
	return accumulated*denominator > total*numerator
}

func emptyContractValidatorSet() updateclientContract.IICS07TendermintMsgsValidatorSet {
	return updateclientContract.IICS07TendermintMsgsValidatorSet{
		Validators:       []updateclientContract.IICS07TendermintMsgsValidatorInfo{},
		HasProposer:      false,
		Proposer:         updateclientContract.IICS07TendermintMsgsValidatorInfo{},
		TotalVotingPower: 0,
	}
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
	// (issue #76 #2). The on-chain height wins when it is ahead.
	onChainTrusted, err := fetchOnChainTrustedHeight(ctx)
	if err != nil {
		return nil, err
	}
	if onChainTrusted > trustedBlock {
		log.Printf("[UpdateCosmosClient] on-chain trusted height %d ahead of cached hint %d; using on-chain",
			onChainTrusted, trustedBlock)
		trustedBlock = onChainTrusted
	} else if trustedBlock == 0 {
		trustedBlock = onChainTrusted
	}
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
	fullProposedHeader := proposedHeader
	currentValidatorsHash := proposedHeader.SignedHeader.Header.ValidatorsHash
	currentValidatorCache, err := getCachedCosmosValidatorSet(ctx, currentValidatorsHash)
	if err != nil {
		return nil, fmt.Errorf("failed to query validator-set cache: %w", err)
	}
	currentValidatorsCacheExists := !currentValidatorCache.isEmpty()
	usedValidatorCache := false
	currentValidatorSetDelta := updateclientContract.IUpdateClientMsgsValidatorSetDelta{}

	if currentValidatorsCacheExists {
		usedValidatorCache = true
		log.Printf("[UpdateCosmosClient] validator quorum cache hit: hash=%x; omitting current validator set", currentValidatorsHash)
		proposedHeader.ValidatorSet = emptyContractValidatorSet()
	} else {
		baseValidatorsHash := consensusState.NextValidatorsHash
		baseValidatorCache, err := getCachedCosmosValidatorSet(ctx, baseValidatorsHash)
		if err != nil {
			return nil, fmt.Errorf("failed to query validator-set delta base cache: %w", err)
		}
		if delta, ok, reason := detectValidatorSetDelta(
			baseValidatorsHash,
			baseValidatorCache,
			proposedHeader.ValidatorSet,
		); ok {
			currentValidatorSetDelta = delta
			currentValidatorsCacheExists = true
			usedValidatorCache = true
			log.Printf(
				"[UpdateCosmosClient] validator quorum delta-cache hit: currentHash=%x baseHash=%x changedLeaves=%d; omitting current validator set",
				currentValidatorsHash,
				baseValidatorsHash,
				delta.LeafCount,
			)
			proposedHeader.ValidatorSet = emptyContractValidatorSet()
		} else {
			log.Printf(
				"[UpdateCosmosClient] validator quorum cache miss: hash=%x; delta unavailable from baseHash=%x (%s); sending full validator set",
				currentValidatorsHash,
				baseValidatorsHash,
				reason,
			)
		}
	}

	adjacentUpdate := proposedHeader.SignedHeader.Header.Height == proposedHeader.TrustedHeight.RevisionHeight+1
	if adjacentUpdate || (!currentValidatorsCacheExists && proposedHeader.SignedHeader.Header.ValidatorsHash == consensusState.NextValidatorsHash) {
		proposedHeader.TrustedNextValidatorSet = emptyContractValidatorSet()
	}

	log.Printf("[UpdateCosmosClient] clientState.LatestHeight=(%d,%d) proposedHeader.Height=%d trustedBlock=%d latestBlock=%d",
		clientState.LatestHeight.RevisionNumber, clientState.LatestHeight.RevisionHeight,
		proposedHeader.SignedHeader.Header.Height, trustedLightBlock.BlockHeight, latestLightBlock.BlockHeight)

	// Extract enough non-absent validator signatures to hit 2/3 voting power,
	// then batch-prove them in a single Groth16 proof that reconstructs each
	// CanonicalVote in-circuit.
	extracted, err := prover.ExtractValidatorSignatures(latestLightBlock, chainId)
	if err != nil {
		return nil, fmt.Errorf("extract validator signatures: %w", err)
	}
	if !adjacentUpdate {
		selected, err := selectSignaturesForTrustedOverlap(
			extracted.Candidates,
			latestLightBlock.ValSet.TotalVotingPower(),
			fullProposedHeader.TrustedNextValidatorSet,
			clientState.TrustLevel,
		)
		if err != nil {
			return nil, fmt.Errorf("select trusted-overlap signatures: %w", err)
		}
		extracted.Signatures = selected
	}
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
	timestampSeconds := make([]uint64, bucket)
	timestampNanos := make([]uint32, bucket)
	active := make([]bool, bucket)
	trustedOverlapIndices := buildTrustedOverlapIndices(fullProposedHeader.TrustedNextValidatorSet, paddedSigs, bucket)
	for i, s := range paddedSigs {
		signerIndices[i] = uint32(s.Index)
		copy(signerPubkeys[i][:], s.PublicKey)
		timestampSeconds[i] = uint64(s.TimestampSeconds)
		timestampNanos[i] = uint32(s.TimestampNanos)
		active[i] = s.Active
	}

	msg := updateclientContract.IUpdateClientMsgsMsgUpdateClient{
		ClientState:              clientState,
		TrustedConsensusState:    consensusState,
		Time:                     big.NewInt(time.Now().UnixNano()),
		ProposedHeader:           proposedHeader,
		Proof:                    proof,
		Commitments:              commitments,
		CommitmentPok:            commitmentPok,
		Bucket:                   uint16(bucket),
		SignerIndices:            signerIndices,
		SignerPubkeys:            signerPubkeys,
		TimestampSeconds:         timestampSeconds,
		TimestampNanos:           timestampNanos,
		Active:                   active,
		TrustedOverlapIndices:    trustedOverlapIndices,
		CurrentValidatorSetDelta: currentValidatorSetDelta,
	}

	fullValidatorSetFallbackMsg := msg
	fullValidatorSetFallbackMsg.ProposedHeader = fullProposedHeader
	fullValidatorSetFallbackMsg.CurrentValidatorSetDelta = updateclientContract.IUpdateClientMsgsValidatorSetDelta{}

	return &CosmosClientUpdateBuildResult{
		Msg:                         msg,
		HasMsg:                      true,
		LightBlock:                  latestLightBlock,
		UsedValidatorCache:          usedValidatorCache,
		FullValidatorSetFallbackMsg: fullValidatorSetFallbackMsg,
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
