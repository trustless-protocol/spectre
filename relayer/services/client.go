// This file is the light-client lifecycle for BOTH ends: create, update, refresh,
// and read the on-chain client state.
//
// The two ends are together on purpose. The symmetry rule is about DISTANCE --
// two mirrors must be comparable -- and at this size they sit within one
// scrolling window. They are also not true mirrors: CreateCosmosClient drives an
// 08-wasm client from beacon data and CreateEthClient drives a SpectreClient from
// a Groth16 proof. There is no diff to run between them, so separating them would
// buy a smaller number and no new ability.
package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	spectreContract "relayer/bindings/SpectreClient"
	updateclientContract "relayer/bindings/UpdateClient"
	relayerclient "relayer/client"
	"relayer/prover"
	"strconv"
	"strings"
	"time"

	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const ICS26_IBC_STORAGE_SLOT = "0x1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600"

type Worker struct {
	TxHandler TransactionHandler
	Prover    Prover
}

// cosmosClientDeps is the narrow capability set shared by the Cosmos light
// client update helpers. It is package-private so the command layer cannot
// accidentally turn it into another process-wide context object.
type cosmosClientDeps struct {
	cosmos            CosmosEndpoint
	evm               EVMEndpoint
	fetchTimeout      time.Duration
	rotationThreshold string
}

const cosmosCatchUpSafetySlots uint64 = 3

func (w *Worker) CreateCosmosClient(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint, ids ClientIDs, proofType string, trustingPeriod uint32, trustedBlock int64, trustLevel string, clockDrift uint32) (common.Address, error) {
	genesis, err := relayerclient.GetGenesis(cosmos.CosmosClient(), trustedBlock, trustingPeriod, trustLevel, proofType, clockDrift)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get genesis: %w", err)
	}

	clientState := genesis.TrustedClientState
	consensusState := genesis.TrustedConsensusState

	log.Printf("[CreateCosmosClient] clientState: chainId=%s trustLevel=%d/%d height=%d/%d trustingPeriod=%d unbondingPeriod=%d isFrozen=%v",
		clientState.ChainId, clientState.TrustLevel.Numerator, clientState.TrustLevel.Denominator,
		clientState.LatestHeight.RevisionNumber, clientState.LatestHeight.RevisionHeight,
		clientState.TrustingPeriod, clientState.UnbondingPeriod, clientState.IsFrozen)

	clientStateEncoded, err := relayerclient.EncodeClientState(clientState)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to encode client state: %w", err)
	}

	// LC-03: the constructor now takes the full consensus state (not a pre-hashed
	// bytes32) so it can assert the genesis pinned validator set actually matches
	// consensusState.nextValidatorsHash on-chain, rather than trusting the deployer's
	// hash and pin to agree. spectreContract.IICS07TendermintMsgsConsensusState is
	// structurally identical to updateClientContract's (same ABI type, different
	// generated Go package), so this is a direct field-for-field conversion.
	spectreConsensusState := spectreContract.IICS07TendermintMsgsConsensusState{
		Timestamp:          consensusState.Timestamp,
		Root:               consensusState.Root,
		NextValidatorsHash: consensusState.NextValidatorsHash,
	}
	return w.TxHandler.CreateCosmosClientContract(stdCtx, evm, ids, clientStateEncoded, spectreConsensusState, genesis.InitialPinnedValidatorSet)
}

// ClientUpdateKind discriminates the two split-API entry points the Spectre
// client now exposes: advance appHash against the pinned set (application), or
// rotate + re-pin the validator set (consensus).
type ClientUpdateKind int

const (
	// ApplicationUpdate advances the appHash against the already-pinned
	// validator set (updateApplicationState) — the frequent, cheap path.
	ApplicationUpdate ClientUpdateKind = iota
	// ConsensusUpdate rotates and re-pins the validator set
	// (updateConsensusState) — needed when the target header's
	// nextValidatorsHash differs from the on-chain pinned set.
	ConsensusUpdate
)

// CosmosClientUpdateBuildResult is the output of BuildCosmosClientUpdateMsg.
// When HasMsg is false the on-chain client is already at the latest height and
// its pinned set is current — no update tx is needed; only LightBlock is
// populated. Otherwise Kind selects which split-API entry point to submit:
// AppMsg is always the built MsgUpdateApplicationState; NewValSet is populated
// only for ConsensusUpdate (the set to re-pin). IsHop is true only when RLY-01
// fallback selected an intermediate height; HopTarget is that intermediate
// height. A hop result is the only case callers should immediately build
// another update in the same refresh/invocation.
type CosmosClientUpdateBuildResult struct {
	Kind       ClientUpdateKind
	AppMsg     updateclientContract.ISpectreClientMsgsMsgUpdateApplicationState
	NewValSet  spectreContract.IICS07TendermintMsgsValidatorSet
	HasMsg     bool
	IsHop      bool
	HopTarget  int64
	LightBlock *relayerclient.LightBlock
}

// UpdateCosmosClient builds the next MsgUpdateClient for the Tendermint light
// client on Ethereum and submits it via SendEthTx. Used by background routines
// that want to advance the client without packets attached. handleCosmos uses
// the split BuildCosmosClientUpdateMsg builder so it can fold updateClient
// into the same multicall as its packet calls (issue #67 V2).

func (w *Worker) UpdateCosmosClient(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint, clientID string, fetchTimeout time.Duration, rotationThreshold, proofType string, trustedBlock int64, trustLevel string, forceRotation bool, targetHeight int64) (*relayerclient.LightBlock, error) {
	deps := cosmosClientDeps{cosmos: cosmos, evm: evm, fetchTimeout: fetchTimeout, rotationThreshold: rotationThreshold}
	result, err := w.buildCosmosClientUpdateMsg(stdCtx, deps, proofType, trustedBlock, trustLevel, forceRotation, targetHeight)
	if err != nil {
		return nil, err
	}
	if !result.HasMsg {
		return result.LightBlock, nil
	}
	if err := w.TxHandler.SendEthTx(stdCtx, deps.evm, clientID, *result); err != nil {
		log.Printf("[UpdateCosmosClient] SendEthTx failed: %v", err)
		return nil, err
	}
	log.Printf("[UpdateCosmosClient] SendEthTx succeeded")
	return result.LightBlock, nil
}

// maxHopsPerRefresh bounds how many RLY-01 multi-hop iterations one
// RefreshCosmosClient call will submit. If churn requires more hops than this
// to fully catch up, the function returns an error rather than looping
// unbounded, so the periodic-update backoff (relay/module.go) retries soon.
const maxHopsPerRefresh = 16
const hopExhaustionSampleLimit = 4

// RefreshCosmosClient advances the Spectre light client on Ethereum without any
// packets attached (background freshness routine). It reads the authoritative
// on-chain trusted height and pinned-set hash, builds the next update via the
// split-API builder, and submits whichever kind the builder decided
// (updateApplicationState vs updateConsensusState — the latter rotates the
// pinned validator set to the block's nextValidatorsHash).
//
// When the on-chain client is already caught up AND its pinned set is current,
// the builder returns HasMsg=false and this function submits nothing, returning
// the current light block so the caller can refresh its cached timestamp/height.
//
// RLY-01: validator churn can force BuildCosmosClientUpdateMsg onto a
// multi-hop path that only advances the client to one intermediate height per
// call. This loops (re-reading the on-chain trusted height each time, so it
// never advances a cursor speculatively) up to maxHopsPerRefresh times so one
// refresh tick catches up fully when possible.
func (w *Worker) RefreshCosmosClient(
	stdCtx context.Context,
	cosmos CosmosEndpoint,
	evm EVMEndpoint,
	clientID string,
	fetchTimeout time.Duration,
	rotationThreshold string,
	proofType string,
	trustLevel string,
) (*relayerclient.LightBlock, error) {
	deps := cosmosClientDeps{cosmos: cosmos, evm: evm, fetchTimeout: fetchTimeout, rotationThreshold: rotationThreshold}
	var lastBlock *relayerclient.LightBlock
	for hop := 0; hop < maxHopsPerRefresh; hop++ {
		if err := stdCtx.Err(); err != nil {
			return nil, err
		}
		readCtx, cancel := fetchCtx(stdCtx, fetchTimeout)
		onChainTrusted, err := FetchOnChainTrustedHeight(readCtx, deps.evm)
		cancel()
		if err != nil {
			return nil, fmt.Errorf("[RefreshCosmosClient] fetch on-chain height: %w", err)
		}

		// forceRotation: the refresh routine only fires after RefreshInterval
		// without updates, so it is the guaranteed rotation cadence — a stale
		// pinned set is rotated here regardless of how much overlap remains.
		result, err := w.buildCosmosClientUpdateMsg(stdCtx, deps, proofType, onChainTrusted, trustLevel, true, 0)
		if err != nil {
			return nil, fmt.Errorf("[RefreshCosmosClient] build update msg: %w", err)
		}
		if !result.HasMsg {
			log.Printf("[RefreshCosmosClient] client is up to date and pinned set current, skipping")
			return result.LightBlock, nil
		}

		if err := w.TxHandler.SendEthTx(stdCtx, deps.evm, clientID, *result); err != nil {
			return nil, fmt.Errorf("[RefreshCosmosClient] send tx (hop %d): %w", hop, err)
		}
		lastBlock = result.LightBlock
		log.Printf("[RefreshCosmosClient] refresh succeeded (hop=%d kind=%d height=%d isHop=%t hopTarget=%d)",
			hop, result.Kind, lastBlock.BlockHeight, result.IsHop, result.HopTarget)
		if !result.IsHop {
			return lastBlock, nil
		}
	}
	log.Printf("RLY01_HOP_CAP_REACHED refresh did not fully catch up after %d hops; will retry via periodic backoff", maxHopsPerRefresh)
	return lastBlock, fmt.Errorf("RLY01: pinned-set catch-up incomplete after %d hops", maxHopsPerRefresh)
}

// getOnChainPinnedValidatorsHash reads the CometBFT validatorsHash of the
// currently-pinned validator set from the Spectre client on ETH. Used to decide
// whether an update must rotate the pinned set (updateConsensusState) or can
// take the cheap application-state path.
func getOnChainPinnedValidatorsHash(stdCtx context.Context, ctx EVMEndpoint) ([32]byte, error) {
	spectre, err := spectreContract.NewContractSpectreClient(*ctx.SpectreClientContract(), ctx.EthClient())
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to create Spectre client instance: %w", err)
	}
	return spectre.GetPinnedValidatorsHash(&bind.CallOpts{Context: stdCtx})
}

// FetchOnChainTrustedHeight reads the ICS07 client state on ETH and returns its
// latest trusted revision height. This is a cheap eth_call relative to the
// Groth16 proof, so it's always worth doing before committing to proof gen.
func FetchOnChainTrustedHeight(stdCtx context.Context, ctx EVMEndpoint) (int64, error) {
	onChainClientState, err := fetchOnChainClientState(stdCtx, ctx)
	if err != nil {
		return 0, err
	}
	log.Printf("[UpdateCosmosClient] On-chain client state: chainId=%s height=(%d,%d) frozen=%v",
		onChainClientState.ChainId, onChainClientState.LatestHeight.RevisionNumber,
		onChainClientState.LatestHeight.RevisionHeight, onChainClientState.IsFrozen)

	return clientStateRevisionHeightInt64(onChainClientState)
}

func fetchOnChainClientState(stdCtx context.Context, ctx EVMEndpoint) (relayerclient.ClientState, error) {
	ics07, err := spectreContract.NewContractSpectreClient(*ctx.SpectreClientContract(), ctx.EthClient())
	if err != nil {
		return relayerclient.ClientState{}, fmt.Errorf("failed to create ICS07 instance: %w", err)
	}
	clientStateBytes, err := ics07.GetClientState(&bind.CallOpts{Context: stdCtx})
	if err != nil {
		return relayerclient.ClientState{}, fmt.Errorf("failed to get on-chain client state: %w", err)
	}
	onChainClientState, err := relayerclient.DecodeClientState(clientStateBytes)
	if err != nil {
		return relayerclient.ClientState{}, fmt.Errorf("failed to decode on-chain client state: %w", err)
	}
	return onChainClientState, nil
}

func clientStateRevisionHeightInt64(clientState relayerclient.ClientState) (int64, error) {
	height := clientState.LatestHeight.RevisionHeight
	if height > uint64(^uint64(0)>>1) {
		return 0, fmt.Errorf("on-chain trusted height overflows int64: %d", height)
	}
	return int64(height), nil
}

// hopCandidate holds the light block and extracted signatures for a height
// chosen by findHighestFeasibleHop as an RLY-01 multi-hop update target.
type hopCandidate struct {
	lightBlock *relayerclient.LightBlock
	candidates []prover.ValidatorSignature
	selected   []prover.ValidatorSignature
}

type hopProbeResult struct {
	lightBlock *relayerclient.LightBlock
	candidates []prover.ValidatorSignature
	selected   []prover.ValidatorSignature
	quorumErr  error
}

// BuildCosmosClientUpdateMsg fetches the latest Tendermint light block,
// generates the Groth16 batch proof, and returns a CosmosClientUpdateBuildResult
// (discriminated by Kind into application-state vs consensus-state update)
// WITHOUT submitting it. Callers either pass the result to SendEthTx directly
// (UpdateCosmosClient/RefreshCosmosClient) or fold it into a multicall alongside
// packet calls (handleCosmos / V2).
//
// When the target header's nextValidatorsHash differs from the on-chain pinned
// hash, forceRotation=true (the refresh-routine cadence path) always upgrades
// to ConsensusUpdate; forceRotation=false upgrades only once the pinned set's
// overlap with the target block's signers has decayed to the configured
// rotation threshold (Config.RotationThreshold). Frequent callers pass false so
// routine voting-power churn — which moves nextValidatorsHash almost every
// block on live chains — doesn't push every update onto the heavy path.
//
// HasMsg=false signals the on-chain client is already at the latest block —
// no update needed; LightBlock still returned so callers can use it for
// membership proofs.
//
// targetHeight overrides the update destination: 0 uses the chain's current
// latest height; a nonzero value must not exceed it. This is the RLY-01
// --target-height operator stopgap (relayer/cmd/main.go's update-client).
func (w *Worker) BuildCosmosClientUpdateMsg(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint, fetchTimeout time.Duration, rotationThreshold string, proofType string, trustedBlock int64, trustLevel string, forceRotation bool, targetHeight int64) (*CosmosClientUpdateBuildResult, error) {
	return w.buildCosmosClientUpdateMsg(stdCtx, cosmosClientDeps{cosmos: cosmos, evm: evm, fetchTimeout: fetchTimeout, rotationThreshold: rotationThreshold}, proofType, trustedBlock, trustLevel, forceRotation, targetHeight)
}

func (w *Worker) buildCosmosClientUpdateMsg(stdCtx context.Context, ctx cosmosClientDeps, proofType string, trustedBlock int64, trustLevel string, forceRotation bool, targetHeight int64) (*CosmosClientUpdateBuildResult, error) {
	if err := stdCtx.Err(); err != nil {
		return nil, err
	}
	statusCtx, cancelStatus := fetchCtx(stdCtx, ctx.fetchTimeout)
	status, err := ctx.cosmos.CosmosClient().Status(statusCtx)
	cancelStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	target := status.SyncInfo.LatestBlockHeight
	if targetHeight != 0 {
		if targetHeight > status.SyncInfo.LatestBlockHeight {
			return nil, fmt.Errorf("target height %d exceeds chain latest %d", targetHeight, status.SyncInfo.LatestBlockHeight)
		}
		target = targetHeight
	}

	log.Printf("[UpdateCosmosClient] called with trustedBlock=%d, latestBlockHeight=%d, target=%d", trustedBlock, status.SyncInfo.LatestBlockHeight, target)

	// Always read the authoritative on-chain trusted height before deciding
	// whether to (re)generate the expensive Groth16 proof. The caller passes a
	// cached `trustedBlock` hint, but across chunked flushes — or after another
	// relayer / a crash-replay advanced the client — that cache can lag the real
	// chain state, causing redundant proofs for a range already on-chain
	// (issue #76 #2). The on-chain height always wins because only it is
	// guaranteed to identify a stored consensus state.
	onChainCtx, cancelOnChain := fetchCtx(stdCtx, ctx.fetchTimeout)
	onChainTrusted, err := FetchOnChainTrustedHeight(onChainCtx, ctx.evm)
	cancelOnChain()
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
	if trustedBlock >= target {
		if trustedBlock == target {
			log.Printf("[UpdateCosmosClient] client is up to date (trusted=%d, target=%d), skipping tx",
				trustedBlock, target)
			lightCtx, cancel := fetchCtx(stdCtx, ctx.fetchTimeout)
			defer cancel()
			lightBlock, err := relayerclient.GetLightBlock(lightCtx, ctx.cosmos.CosmosClient(), trustedBlock)
			if err != nil {
				return nil, fmt.Errorf("failed to get current light block while up-to-date: %w", err)
			}
			return &CosmosClientUpdateBuildResult{LightBlock: lightBlock}, nil
		}
		return nil, fmt.Errorf("trusted block is ahead of target height (trusted=%d, target=%d)", trustedBlock, target)
	}

	log.Printf("[UpdateCosmosClient] Fetching trustedLightBlock at height %d, targetLightBlock at height %d", trustedBlock, target)
	trustedCtx, cancelTrusted := fetchCtx(stdCtx, ctx.fetchTimeout)
	trustedLightBlock, err := relayerclient.GetLightBlock(trustedCtx, ctx.cosmos.CosmosClient(), trustedBlock)
	cancelTrusted()
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted light block: %w", err)
	}

	targetCtx, cancelTarget := fetchCtx(stdCtx, ctx.fetchTimeout)
	latestLightBlock, err := relayerclient.GetLightBlock(targetCtx, ctx.cosmos.CosmosClient(), target)
	cancelTarget()
	if err != nil {
		return nil, fmt.Errorf("failed to get target light block: %w", err)
	}

	if trustedLightBlock.SignedHeader.Header.Height < 0 {
		return nil, fmt.Errorf("trusted light block header height cannot be negative: %d", trustedLightBlock.SignedHeader.Header.Height)
	}

	// The Spectre split-API msgs no longer carry client state — SpectreClient
	// reads it from its own Store. proofType/trustLevel are validated by the
	// client on-chain, so the builder only needs chainId (for signature
	// extraction) and the trusted consensus state to prove against.
	chainId := trustedLightBlock.SignedHeader.Header.ChainID

	consensusState := updateclientContract.IICS07TendermintMsgsConsensusState{
		Timestamp:          big.NewInt(trustedLightBlock.SignedHeader.Header.Time.UnixNano()),
		Root:               bytesToBytes32(trustedLightBlock.SignedHeader.Header.AppHash),
		NextValidatorsHash: bytesToBytes32(trustedLightBlock.SignedHeader.NextValidatorsHash),
	}

	readCtx, cancelRead := fetchCtx(stdCtx, ctx.fetchTimeout)
	pinnedValidatorSet, err := getPinnedCosmosValidatorSet(readCtx, ctx.evm)
	cancelRead()
	if err != nil {
		return nil, fmt.Errorf("failed to query pinned validator set: %w", err)
	}

	// Extract non-absent validator signatures, then select enough signers that
	// overlap the pinned validator set to exceed 2/3 of pinned voting power.
	//
	// RLY-01: selectSignaturesForPinnedSet requires >2/3 of the *pinned* set
	// among the target block's signers. After large validator churn the
	// pinned set may no longer sign the target block with 2/3 power; when
	// that happens, fall back to a multi-hop update through the highest
	// intermediate height that still clears pinned-set quorum, rather than
	// failing outright.
	extracted, err := prover.ExtractValidatorSignatures(latestLightBlock, chainId, nil)
	if err != nil {
		return nil, fmt.Errorf("extract validator signatures: %w", err)
	}
	selected, err := selectSignaturesForPinnedSet(extracted.Candidates, pinnedValidatorSet)
	isHop := false
	if err != nil {
		log.Printf("RLY01_HOP_FALLBACK direct update trusted=%d->target=%d failed pinned quorum (%v); searching for intermediate hop",
			trustedBlock, target, err)
		hop, hopErr := findHighestFeasibleHop(stdCtx, ctx, trustedBlock, target, chainId, pinnedValidatorSet, err)
		if hopErr != nil {
			log.Printf("RLY01_HOP_EXHAUSTED no provable height beyond trusted=%d; manual intervention required", trustedBlock)
			return nil, fmt.Errorf("RLY01_HOP_EXHAUSTED: no provable hop above trusted height %d: %w", trustedBlock, hopErr)
		}
		latestLightBlock = hop.lightBlock
		extracted.Candidates = hop.candidates
		selected = hop.selected
		isHop = true
	} else if overlap := pinnedOverlapPower(extracted.Candidates, pinnedValidatorSet); overlap*4 <= pinnedValidatorSet.totalPower*3 {
		// RLY-01 leading indicator: overlap is closing in on the 2/3 floor
		// (<=3/4 of total) even though this update still succeeds outright.
		// Greppable warning for ops until a real metrics stack exists (RLY-02).
		log.Printf("RLY01_QUORUM_WARN pinned-set overlap approaching 2/3 floor: overlap=%d total=%d height=%d chainId=%s",
			overlap, pinnedValidatorSet.totalPower, latestLightBlock.BlockHeight, chainId)
	}
	extracted.Signatures = selected

	proposedHeader, err := latestLightBlock.IntoHeader(*trustedLightBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to convert light block into header: %w", err)
	}
	log.Printf("[UpdateCosmosClient] proposedHeader.Height=%d trustedBlock=%d updateToBlock=%d isHop=%t",
		proposedHeader.SignedHeader.Header.Height, trustedLightBlock.BlockHeight, latestLightBlock.BlockHeight, isHop)
	log.Printf("[UpdateCosmosClient] Generating Groth16 batch proof for %d validator signatures...", len(extracted.Signatures))
	if err := stdCtx.Err(); err != nil {
		return nil, err
	}
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

	updateTime, err := ethLatestHeaderTimestampNanos(stdCtx, ctx.evm, ctx.fetchTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to get ethereum timestamp for update freshness: %w", err)
	}

	appMsg := updateclientContract.ISpectreClientMsgsMsgUpdateApplicationState{
		TrustedConsensusState: consensusState,
		ProposedHeader:        proposedHeader,
		Time:                  updateTime,
		Proof: updateclientContract.ISpectreClientMsgsBatchProof{
			Proof:                  proof,
			Commitments:            commitments,
			CommitmentPok:          commitmentPok,
			Bucket:                 uint16(bucket),
			SignerIndices:          signerIndices,
			PinnedValidatorIndices: pinnedValidatorIndices,
			SignerPubkeys:          signerPubkeys,
			Active:                 active,
		},
	}

	var hopTarget int64
	if isHop {
		hopTarget = latestLightBlock.BlockHeight
	}
	result := &CosmosClientUpdateBuildResult{
		Kind:       ApplicationUpdate,
		AppMsg:     appMsg,
		HasMsg:     true,
		IsHop:      isHop,
		HopTarget:  hopTarget,
		LightBlock: latestLightBlock,
	}

	// Decide application-state vs consensus-state update: if the block we're
	// advancing to commits to a validator set (nextValidatorsHash) that differs
	// from the currently-pinned set on-chain, rotating + re-pinning keeps the
	// pinned set current for future proofs — but the rotation is the heavy path
	// (full valset calldata + buildCache + SSTORE2 write + snapshot push), so
	// it is gated by shouldRotatePinnedSet rather than taken on every hash
	// difference. Otherwise the cheap application-state path suffices.
	pinnedHashCtx, cancelPinnedHash := fetchCtx(stdCtx, ctx.fetchTimeout)
	pinnedHash, err := getOnChainPinnedValidatorsHash(pinnedHashCtx, ctx.evm)
	cancelPinnedHash()
	if err != nil {
		return nil, fmt.Errorf("failed to query on-chain pinned validators hash: %w", err)
	}
	targetNextValHash := bytesToBytes32(latestLightBlock.SignedHeader.NextValidatorsHash)
	switch {
	case isHop:
		// A hop's entire purpose is to re-pin at an intermediate validator
		// set so a subsequent call can make further progress toward the
		// original target — always rotate, bypassing shouldRotatePinnedSet's
		// threshold gating (RLY-01).
		newValSet, err := relayerclient.ValidatorSetToContract(latestLightBlock.NextValSet, "updateConsensusState")
		if err != nil {
			return nil, fmt.Errorf("convert next validator set: %w", err)
		}
		result.Kind = ConsensusUpdate
		result.NewValSet = newValSet
		log.Printf("RLY01_HOP_ROTATE hop height=%d pinned=%x target nextValHash=%x; rotating via updateConsensusState",
			latestLightBlock.BlockHeight, pinnedHash[:4], targetNextValHash[:4])
	case targetNextValHash != pinnedHash:
		rotationThreshold, err := ParseRotationThreshold(ctx.rotationThreshold)
		if err != nil {
			return nil, fmt.Errorf("invalid rotation threshold: %w", err)
		}
		overlap := pinnedOverlapPower(extracted.Candidates, pinnedValidatorSet)
		if shouldRotatePinnedSet(forceRotation, overlap, pinnedValidatorSet.totalPower, rotationThreshold) {
			newValSet, err := relayerclient.ValidatorSetToContract(latestLightBlock.NextValSet, "updateConsensusState")
			if err != nil {
				return nil, fmt.Errorf("convert next validator set: %w", err)
			}
			result.Kind = ConsensusUpdate
			result.NewValSet = newValSet
			log.Printf("[UpdateCosmosClient] pinned set stale (pinned=%x target nextValHash=%x overlap=%d/%d force=%t); rotating via updateConsensusState",
				pinnedHash[:4], targetNextValHash[:4], overlap, pinnedValidatorSet.totalPower, forceRotation)
		} else {
			log.Printf("[UpdateCosmosClient] pinned set differs from target nextValidatorsHash (pinned=%x target=%x) but overlap %d/%d is above rotation threshold %d/%d; staying on updateApplicationState",
				pinnedHash[:4], targetNextValHash[:4], overlap, pinnedValidatorSet.totalPower,
				rotationThreshold.Numerator, rotationThreshold.Denominator)
		}
	}

	return result, nil
}

func (w *Worker) CreateEthClient(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint, counterpartyClientID string, checksum string) (string, error) {
	beaconAPIURL := evm.BeaconAPIURL
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

	slotsPerEpoch, err := strconv.ParseUint(spec.SlotsPerEpoch, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parse slots_per_epoch %q: %w", spec.SlotsPerEpoch, err)
	}
	checkpointSlot, blockRoot, bootstrap, err := resolveBootstrapCheckpoint(beaconAPIURL, checkpointSlot, slotsPerEpoch)
	if err != nil {
		return "", err
	}
	log.Printf("[CreateEthClient] bootstrap checkpoint slot=%s root=%s", checkpointSlot, blockRoot)
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
	chainId, err := evm.EthClient().ChainID(stdCtx)
	if err != nil {
		return "", fmt.Errorf("failed to get eth chain id: %w", err)
	}
	log.Printf("[CreateEthClient] ethereum chain id=%s", chainId.String())

	epochsPerSyncCommitteePeriod, err := strconv.ParseUint(spec.EpochsPerSyncCommitteePeriod, 10, 64)
	if err != nil {
		return "", err
	}

	// The full fork schedule (through Fulu) is baked into the client state; the
	// light client selects the active version per header — see ToForkParameters.
	forkParameters, err := spec.ToForkParameters()
	if err != nil {
		return "", fmt.Errorf("failed to get fork parameters: %w", err)
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
	clientState := relayerclient.EthereumClientState{
		ChainID:                      chainId.Uint64(),
		EpochsPerSyncCommitteePeriod: epochsPerSyncCommitteePeriod,
		ForkParameters:               *forkParameters,
		GenesisSlot:                  0,
		GenesisTime:                  genesisTime,
		GenesisValidatorsRoot:        genesis.GenesisValidatorsRoot,
		IbcCommitmentSlot:            ICS26_IBC_STORAGE_SLOT,
		IbcContractAddress:           evm.RouterContract().String(),
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
	lightClientUpdates, err := relayerclient.GetLightClientUpdates(bctx, evm.BeaconAPIURL, latestPeriod, 1)
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

	// The ETH beacon client's counterparty is the ETH-side Cosmos router client id
	// (a config value); it must be present.
	if counterpartyClientID == "" {
		return "", fmt.Errorf("[CreateEthClient] cosmos router client id is not configured")
	}
	return w.TxHandler.CreateWasmClient(stdCtx, cosmos, &wasmClientState, &wasmConsensusState, counterpartyClientID)
}

type EthClientUpdateResult struct {
	Headers        [][]byte
	EthClientState *relayerclient.EthereumClientState
	ProofTimestamp uint64
	SigSlot        uint64
}

func (w *Worker) BuildEthClientUpdateHeaders(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint, ethClientID string) (*EthClientUpdateResult, error) {
	beaconAPIURL := evm.BeaconAPIURL
	if beaconAPIURL == "" {
		return nil, fmt.Errorf("beacon API URL is not configured")
	}

	if ethClientID == "" {
		return nil, fmt.Errorf("ethereum client ID is not configured")
	}

	readCtx, cancelRead := fetchCtx(stdCtx, defaultFetchTimeout)
	ethClientState, err := relayerclient.GetEthereumClientState(readCtx, cosmos.CosmosClient(), ethClientID)
	cancelRead()
	if err != nil {
		return nil, fmt.Errorf("failed to get ethereum client state: %w", err)
	}
	trustedSlot := ethClientState.LatestSlot

	bctx, bcancel := context.WithTimeout(stdCtx, 15*time.Second)
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
			Headers:        nil,
			EthClientState: cloneEthereumClientState(ethClientState),
			ProofTimestamp: ethClientState.ComputeTimestampAtSlot(trustedSlot),
		}, nil
	}

	trustedPeriod := ethClientState.ComputeSyncCommitteePeriodAtSlot(trustedSlot)
	targetPeriod := ethClientState.ComputeSyncCommitteePeriodAtSlot(finalizedSlot)

	log.Printf("[UpdateEthClient] trustedPeriod=%d targetPeriod=%d", trustedPeriod, targetPeriod)

	headers, err := w.buildEthClientUpdateHeadersWithPeriodCrossing(stdCtx, beaconAPIURL, ethClientState, trustedSlot, trustedPeriod, targetPeriod, finalityUpdate, finalizedSlot)
	if err != nil {
		return nil, err
	}

	proofState, proofTimestamp, err := ethProofStateFromFinalityUpdate(ethClientState, finalityUpdate, finalizedSlot)
	if err != nil {
		return nil, err
	}

	sigSlot, err := parseSlot(finalityUpdate.SignatureSlot)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signature slot: %w", err)
	}
	return &EthClientUpdateResult{
		Headers:        headers,
		EthClientState: proofState,
		ProofTimestamp: proofTimestamp,
		SigSlot:        sigSlot,
	}, nil
}

// WaitForCosmosCatchUp polls until the Cosmos chain time covers the update's
// signature slot (so the wasm beacon update verifies), or until stdCtx is
// cancelled. stdCtx makes the up-to-60 × 5s poll abort promptly on shutdown
// instead of stranding a goroutine for minutes during teardown.
func (w *Worker) WaitForCosmosCatchUp(stdCtx context.Context, cosmos CosmosEndpoint, ethClientState *relayerclient.EthereumClientState, sigSlot uint64) error {
	var lastSlot uint64
	requiredSlot := sigSlot + cosmosCatchUpSafetySlots
	for range 60 {
		if stdCtx.Err() != nil {
			return stdCtx.Err()
		}
		statusCtx, cancel := fetchCtx(stdCtx, defaultFetchTimeout)
		status, err := cosmos.CosmosClient().Status(statusCtx)
		cancel()
		if err != nil {
			return fmt.Errorf("query Cosmos status while waiting for catch-up: %w", err)
		}
		cosmosTime := uint64(status.SyncInfo.LatestBlockTime.Unix())
		currentSlot := ethClientState.ComputeSlotAtTimestamp(cosmosTime)
		lastSlot = currentSlot
		if cosmosCurrentSlotReady(currentSlot, sigSlot) {
			log.Printf("[updateEthClient] timing OK: currentSlot=%d >= requiredSlot=%d (signatureSlot=%d safety=%d)",
				currentSlot, requiredSlot, sigSlot, cosmosCatchUpSafetySlots)
			return nil
		}
		log.Printf("[updateEthClient] waiting for target chain to catch up to required slot %d (signatureSlot=%d current=%d safety=%d)",
			requiredSlot, sigSlot, currentSlot, cosmosCatchUpSafetySlots)
		select {
		case <-stdCtx.Done():
			return stdCtx.Err()
		case <-time.After(5 * time.Second):
		}
	}
	return fmt.Errorf(
		"cosmos catch-up: chain still at slot %d after 60 polls, need %d (signatureSlot=%d safety=%d)",
		lastSlot, requiredSlot, sigSlot, cosmosCatchUpSafetySlots,
	)
}

func bytesToBytes32(data []byte) [32]byte {
	var result [32]byte
	copy(result[:], data)
	return result
}

// maxBootstrapCheckpointStepBack bounds how far back resolveBootstrapCheckpoint walks
// looking for a servable checkpoint. One step is normally enough; the cap only stops a
// pathological walk against a node that serves no bootstraps at all.
const maxBootstrapCheckpointStepBack = 8
