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
	"sort"
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

func NewWorker(txHandler TransactionHandler, prover Prover) *Worker {
	return &Worker{
		txHandler,
		prover,
	}
}

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
		onChainTrusted, err := FetchOnChainTrustedHeightWithContext(readCtx, deps.evm)
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
func FetchOnChainTrustedHeight(ctx EVMEndpoint) (int64, error) {
	return FetchOnChainTrustedHeightWithContext(context.Background(), ctx)
}

func FetchOnChainTrustedHeightWithContext(stdCtx context.Context, ctx EVMEndpoint) (int64, error) {
	onChainClientState, err := fetchOnChainClientStateWithContext(stdCtx, ctx)
	if err != nil {
		return 0, err
	}
	log.Printf("[UpdateCosmosClient] On-chain client state: chainId=%s height=(%d,%d) frozen=%v",
		onChainClientState.ChainId, onChainClientState.LatestHeight.RevisionNumber,
		onChainClientState.LatestHeight.RevisionHeight, onChainClientState.IsFrozen)

	return clientStateRevisionHeightInt64(onChainClientState)
}

func fetchOnChainClientState(ctx EVMEndpoint) (relayerclient.ClientState, error) {
	return fetchOnChainClientStateWithContext(context.Background(), ctx)
}

func fetchOnChainClientStateWithContext(stdCtx context.Context, ctx EVMEndpoint) (relayerclient.ClientState, error) {
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

func getPinnedCosmosValidatorSet(stdCtx context.Context, ctx EVMEndpoint) (pinnedCosmosValidatorSet, error) {
	ics07, err := spectreContract.NewContractSpectreClient(*ctx.SpectreClientContract(), ctx.EthClient())
	if err != nil {
		return pinnedCosmosValidatorSet{}, fmt.Errorf("failed to create ICS07 instance: %w", err)
	}
	out, err := ics07.GetPinnedValidatorSet(&bind.CallOpts{Context: stdCtx})
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

// pinnedOverlapPower sums the pinned-set voting power held by the target
// block's commit signers — the power available to prove against the pinned
// set. Signers outside the pinned set contribute zero. The measurement is
// per-commit: a pinned validator that merely missed this one block counts as
// zero, which makes the rotation trigger conservative (fires early, not late).
func pinnedOverlapPower(candidates []prover.ValidatorSignature, pinnedSet pinnedCosmosValidatorSet) int64 {
	powerByPubkey := pinnedSet.powerByPubkey()
	var overlap int64
	for _, candidate := range candidates {
		overlap += powerByPubkey[bytesToBytes32(candidate.PublicKey)]
	}
	return overlap
}

// shouldRotatePinnedSet decides whether an update whose target
// nextValidatorsHash differs from the on-chain pinned hash rotates the pinned
// set (updateConsensusState) or defers and stays on updateApplicationState.
// forceRotation (the refresh-routine cadence path) always rotates; otherwise
// rotate once overlap/total ≤ threshold. With threshold "1/1" this rotates on
// any pinned-set change.
func shouldRotatePinnedSet(forceRotation bool, overlapPower, totalPower int64, threshold relayerclient.TrustThreshold) bool {
	if forceRotation {
		return true
	}
	overlap := new(big.Int).Mul(big.NewInt(overlapPower), big.NewInt(int64(threshold.Denominator)))
	total := new(big.Int).Mul(big.NewInt(totalPower), big.NewInt(int64(threshold.Numerator)))
	return overlap.Cmp(total) <= 0
}

func ethLatestHeaderTimestampNanos(stdCtx context.Context, ctx EVMEndpoint, fetchTimeout time.Duration) (*big.Int, error) {
	hctx, cancel := fetchCtx(stdCtx, fetchTimeout)
	defer cancel()

	header, err := ctx.EthClient().HeaderByNumber(hctx, nil)
	if err != nil {
		return nil, err
	}
	if header == nil || header.Time == 0 {
		return nil, fmt.Errorf("latest ethereum header is nil or has zero timestamp")
	}
	return new(big.Int).Mul(new(big.Int).SetUint64(header.Time), big.NewInt(1_000_000_000)), nil
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

// binarySearchHighestFeasible finds the highest height in (trusted, latest]
// for which feasible returns true, assuming feasibility is non-increasing as
// height moves away from trusted (validator churn accumulates over time; it
// does not spontaneously reverse). feasible must distinguish "infeasible at
// this height" (ok=false, err=nil) from a genuine fetch/RPC error (err!=nil);
// the latter is propagated immediately rather than treated as infeasibility.
// latest itself is assumed already-probed-and-infeasible by the caller and is
// never re-probed here. Returns ok=false if no intermediate height is feasible.
//
// Pure and side-effect free so it's unit-testable without RPC mocking.
func binarySearchHighestFeasible(trusted, latest int64, feasible func(height int64) (ok bool, err error)) (int64, bool, error) {
	if latest <= trusted+1 {
		return 0, false, nil
	}

	lo, hi := trusted, latest
	found := int64(0)
	haveFound := false
	for hi-lo > 1 {
		mid := lo + (hi-lo)/2
		ok, err := feasible(mid)
		if err != nil {
			return 0, false, fmt.Errorf("probe height %d: %w", mid, err)
		}
		if ok {
			lo, found, haveFound = mid, mid, true
		} else {
			hi = mid
		}
	}
	return found, haveFound, nil
}

// findHighestFeasibleHop searches for the highest Cosmos height between
// trusted and latest whose commit still carries enough pinned-set signing
// power to clear selectSignaturesForPinnedSet's >2/3 quorum gate. This is the
// RLY-01 multi-hop fallback: used when the direct trusted->latest update no
// longer clears quorum due to validator churn (docs/RELIABILITY.md).
func findHighestFeasibleHop(
	stdCtx context.Context,
	ctx cosmosClientDeps,
	trusted, latest int64,
	chainId string,
	pinnedSet pinnedCosmosValidatorSet,
	directQuorumErr error,
) (hopCandidate, error) {
	probed := make(map[int64]struct{})
	var lastQuorumHeight int64
	var lastQuorumErr error

	probe := func(height int64) (hopProbeResult, error) {
		if err := stdCtx.Err(); err != nil {
			return hopProbeResult{}, err
		}
		probed[height] = struct{}{}
		fetchCtx, cancel := fetchCtx(stdCtx, ctx.fetchTimeout)
		defer cancel()
		lb, err := relayerclient.GetLightBlockWithContext(fetchCtx, ctx.cosmos.CosmosClient(), height)
		if err != nil {
			return hopProbeResult{}, fmt.Errorf("fetch light block at height %d: %w", height, err)
		}
		extracted, err := prover.ExtractValidatorSignatures(lb, chainId, nil)
		if err != nil {
			return hopProbeResult{}, fmt.Errorf("extract signatures at height %d: %w", height, err)
		}
		selected, err := selectSignaturesForPinnedSet(extracted.Candidates, pinnedSet)
		if err != nil {
			return hopProbeResult{
				lightBlock: lb,
				candidates: extracted.Candidates,
				quorumErr:  err,
			}, nil // infeasible at this height, not a fetch error
		}
		return hopProbeResult{
			lightBlock: lb,
			candidates: extracted.Candidates,
			selected:   selected,
		}, nil
	}

	feasible := func(height int64) (bool, error) {
		result, err := probe(height)
		if err != nil {
			return false, err
		}
		if result.selected == nil {
			lastQuorumHeight = height
			lastQuorumErr = result.quorumErr
			return false, nil
		}
		return true, nil
	}

	height, ok, err := binarySearchHighestFeasible(trusted, latest, feasible)
	if err != nil {
		return hopCandidate{}, fmt.Errorf("hop search: %w", err)
	}
	if !ok {
		for h, sampled := trusted+1, int64(0); h < latest && sampled < hopExhaustionSampleLimit; h++ {
			if _, seen := probed[h]; seen {
				continue
			}
			sampled++
			result, err := probe(h)
			if err != nil {
				return hopCandidate{}, fmt.Errorf("hop search sample height %d: %w", h, err)
			}
			if result.selected != nil {
				log.Printf("RLY01_HOP_NON_MONOTONIC_SAMPLE trusted=%d target=%d hop=%d",
					trusted, latest, h)
				return hopCandidate{
					lightBlock: result.lightBlock,
					candidates: result.candidates,
					selected:   result.selected,
				}, nil
			}
			lastQuorumHeight = h
			lastQuorumErr = result.quorumErr
		}
		msg := fmt.Sprintf("no provable height in (%d, %d]: pinned quorum lost immediately past trusted height", trusted, latest)
		if directQuorumErr != nil {
			msg += fmt.Sprintf("; target quorum shortfall: %v", directQuorumErr)
		}
		if lastQuorumErr != nil {
			msg += fmt.Sprintf("; sampled height %d quorum shortfall: %v", lastQuorumHeight, lastQuorumErr)
		}
		return hopCandidate{}, fmt.Errorf("%s", msg)
	}

	result, err := probe(height)
	if err != nil {
		return hopCandidate{}, fmt.Errorf("re-probe chosen hop height %d: %w", height, err)
	}
	if result.selected == nil {
		return hopCandidate{}, fmt.Errorf("hop height %d unexpectedly infeasible on re-check: %w", height, result.quorumErr)
	}
	log.Printf("RLY01_HOP_FOUND trusted=%d target=%d hop=%d", trusted, latest, height)
	return hopCandidate{lightBlock: result.lightBlock, candidates: result.candidates, selected: result.selected}, nil
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
func (w *Worker) BuildCosmosClientUpdateMsg(cosmos CosmosEndpoint, evm EVMEndpoint, fetchTimeout time.Duration, rotationThreshold string, proofType string, trustedBlock int64, trustLevel string, forceRotation bool, targetHeight int64) (*CosmosClientUpdateBuildResult, error) {
	return w.BuildCosmosClientUpdateMsgWithContext(context.Background(), cosmos, evm, fetchTimeout, rotationThreshold, proofType, trustedBlock, trustLevel, forceRotation, targetHeight)
}

func (w *Worker) BuildCosmosClientUpdateMsgWithContext(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint, fetchTimeout time.Duration, rotationThreshold string, proofType string, trustedBlock int64, trustLevel string, forceRotation bool, targetHeight int64) (*CosmosClientUpdateBuildResult, error) {
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
	onChainTrusted, err := FetchOnChainTrustedHeightWithContext(onChainCtx, ctx.evm)
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
			lightBlock, err := relayerclient.GetLightBlockWithContext(lightCtx, ctx.cosmos.CosmosClient(), trustedBlock)
			if err != nil {
				return nil, fmt.Errorf("failed to get current light block while up-to-date: %w", err)
			}
			return &CosmosClientUpdateBuildResult{LightBlock: lightBlock}, nil
		}
		return nil, fmt.Errorf("trusted block is ahead of target height (trusted=%d, target=%d)", trustedBlock, target)
	}

	log.Printf("[UpdateCosmosClient] Fetching trustedLightBlock at height %d, targetLightBlock at height %d", trustedBlock, target)
	trustedCtx, cancelTrusted := fetchCtx(stdCtx, ctx.fetchTimeout)
	trustedLightBlock, err := relayerclient.GetLightBlockWithContext(trustedCtx, ctx.cosmos.CosmosClient(), trustedBlock)
	cancelTrusted()
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted light block: %w", err)
	}

	targetCtx, cancelTarget := fetchCtx(stdCtx, ctx.fetchTimeout)
	latestLightBlock, err := relayerclient.GetLightBlockWithContext(targetCtx, ctx.cosmos.CosmosClient(), target)
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
	// failing outright (see docs/RELIABILITY.md).
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
			log.Printf("RLY01_HOP_EXHAUSTED no provable height beyond trusted=%d; manual intervention required (see docs/RELIABILITY.md)", trustedBlock)
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
	ethClientState, err := relayerclient.GetEthereumClientStateWithContext(readCtx, cosmos.CosmosClient(), ethClientID)
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

func cosmosCurrentSlotReady(currentSlot, sigSlot uint64) bool {
	return currentSlot >= sigSlot+cosmosCatchUpSafetySlots
}

func (w *Worker) buildEthClientUpdateHeadersWithPeriodCrossing(stdCtx context.Context, beaconAPIURL string, ethClientState *relayerclient.EthereumClientState, trustedSlot, trustedPeriod, targetPeriod uint64, finalityUpdate *relayerclient.LightClientFinalityUpdate, finalizedSlot uint64) ([][]byte, error) {
	count := targetPeriod - trustedPeriod + 1
	bctx, bcancel := context.WithTimeout(stdCtx, 15*time.Second)
	lightClientUpdates, err := relayerclient.GetLightClientUpdates(bctx, beaconAPIURL, trustedPeriod, count)
	bcancel()
	if err != nil {
		return nil, fmt.Errorf("failed to get light client updates: %w", err)
	}

	if len(lightClientUpdates) == 0 {
		return nil, fmt.Errorf("no light client updates available for period range %d to %d", trustedPeriod, targetPeriod)
	}

	// Index the fetched updates by the period they belong to. Crossing into period P
	// needs the FULL sync committee of period P, and the update for period P-1 already
	// carries it as next_sync_committee — Merkle-proven against its attested header,
	// which is exactly what next_sync_committee_branch exists for. Reading it from
	// here avoids a light_client/bootstrap call that beacon nodes only answer for the
	// checkpoint roots they happen to retain (see the fallback below).
	updatesByPeriod := make(map[uint64]relayerclient.LightClientUpdate, len(lightClientUpdates))
	for _, u := range lightClientUpdates {
		slot, err := parseSlot(u.FinalizedHeader.Beacon.Slot)
		if err != nil {
			return nil, fmt.Errorf("failed to parse update finalized slot: %w", err)
		}
		updatesByPeriod[ethClientState.ComputeSyncCommitteePeriodAtSlot(slot)] = u
	}

	headers := make([][]byte, 0, count)
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

		syncCommittee, err := syncCommitteeForPeriod(stdCtx, beaconAPIURL, updatesByPeriod, updatePeriod, updateFinalizedSlot)
		if err != nil {
			return nil, err
		}

		header := relayerclient.EthereumHeader{
			ActiveSyncCommittee: relayerclient.ActiveSyncCommittee{
				Next: &syncCommittee,
			},
			ConsensusUpdate: update,
			TrustedSlot:     latestTrustedSlot,
		}

		headerBytes, err := json.Marshal(header)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal update header: %w", err)
		}

		headers = append(headers, headerBytes)
		latestPeriod = updatePeriod
		latestTrustedSlot = updateFinalizedSlot
	}

	// If the latest header is earlier than the finality update, add a header for the finality update.
	if finalizedSlot > latestTrustedSlot {
		attestedSlot := finalityUpdate.AttestedHeader.Beacon.Slot
		log.Printf("[updateEthClient] final update: attestedSlot=%s finalizedSlot=%d latestTrustedSlot=%d",
			attestedSlot, finalizedSlot, latestTrustedSlot)

		// The finality update is signed by the committee active at its attested slot,
		// which the client checks against the current_sync_committee it already trusts.
		// Same sourcing problem as the crossing loop above: after crossing into a new
		// period the beacon will not serve a bootstrap for that period, so prefer the
		// committee carried by the preceding period's update.
		attestedSlotNum, err := parseSlot(attestedSlot)
		if err != nil {
			return nil, fmt.Errorf("failed to parse attested slot: %w", err)
		}
		syncCommittee, err := syncCommitteeForPeriod(
			stdCtx,
			beaconAPIURL, updatesByPeriod,
			ethClientState.ComputeSyncCommitteePeriodAtSlot(attestedSlotNum), attestedSlotNum)
		if err != nil {
			return nil, err
		}

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

		headerBytes, err := json.Marshal(header)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal update header: %w", err)
		}

		headers = append(headers, headerBytes)
	}

	return headers, nil
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

// syncCommitteeForPeriod returns the full sync committee that is active in period.
//
// The Ethereum light client checks the supplied committee against the summary it
// already trusts (ActiveSyncCommittee::Next vs ConsensusState.next_sync_committee),
// so this must be the committee of the period being crossed INTO — which the light
// client update for the preceding period carries as next_sync_committee, proven by
// next_sync_committee_branch.
//
// It used to come from a light_client/bootstrap at the update's block root instead.
// That works only while the beacon node still serves a bootstrap for that particular
// root; most nodes serve bootstraps for a small set of retained checkpoints, so the
// first sync-committee period boundary answered:
//
//	404 NOT_FOUND: Sync committee for period 1 not found
//
// and the client stopped advancing entirely. On mainnet periods roll about every 27
// hours, so that is a client-expiry bug, not just a devnet annoyance. The bootstrap
// path is kept as a fallback for the case where the preceding period's update was not
// returned in the requested range.
func syncCommitteeForPeriod(
	stdCtx context.Context,
	beaconAPIURL string,
	updatesByPeriod map[uint64]relayerclient.LightClientUpdate,
	period, updateFinalizedSlot uint64,
) (relayerclient.SyncCommittee, error) {
	if period > 0 {
		if prev, ok := updatesByPeriod[period-1]; ok && prev.NextSyncCommittee != nil {
			return *prev.NextSyncCommittee, nil
		}
		// The preceding period's update is outside the range the caller fetched — the
		// steady-state case, where trusted and target are the same period so only that
		// one update was requested. Fetch it on its own rather than falling through to
		// a bootstrap the beacon will not serve for this period.
		fctx, fcancel := context.WithTimeout(stdCtx, 15*time.Second)
		prevUpdates, err := relayerclient.GetLightClientUpdates(fctx, beaconAPIURL, period-1, 1)
		fcancel()
		if err == nil {
			for _, u := range prevUpdates {
				if u.NextSyncCommittee != nil {
					return *u.NextSyncCommittee, nil
				}
			}
		}
	}

	bctx, bcancel := context.WithTimeout(stdCtx, 15*time.Second)
	blockRoot, err := relayerclient.GetBeaconBlockRoot(bctx, beaconAPIURL, fmt.Sprintf("%d", updateFinalizedSlot))
	bcancel()
	if err != nil {
		return relayerclient.SyncCommittee{}, fmt.Errorf(
			"period %d: no preceding update carries next_sync_committee and beacon block root lookup failed: %w", period, err)
	}

	bctx, bcancel = context.WithTimeout(stdCtx, 15*time.Second)
	bootstrap, err := relayerclient.GetLightClientBootstrap(bctx, beaconAPIURL, blockRoot)
	bcancel()
	if err != nil {
		return relayerclient.SyncCommittee{}, fmt.Errorf(
			"period %d: no preceding update carries next_sync_committee and bootstrap at slot %d is unavailable: %w",
			period, updateFinalizedSlot, err)
	}
	return bootstrap.Data.CurrentSyncCommittee, nil
}

// maxBootstrapCheckpointStepBack bounds how far back resolveBootstrapCheckpoint walks
// looking for a servable checkpoint. One step is normally enough; the cap only stops a
// pathological walk against a node that serves no bootstraps at all.
const maxBootstrapCheckpointStepBack = 8

// resolveBootstrapCheckpoint finds a finalized checkpoint the beacon will actually serve
// a light-client bootstrap for, returning its slot, block root and bootstrap together so
// the caller's later consistency check against the beacon block still holds.
//
// The finality update's finalized header is NOT always on an epoch boundary, despite
// what the surrounding code used to assume. When the boundary slot is skipped — no block
// proposed — the checkpoint root points back to the last block before it, and beacon
// nodes index bootstraps by the block AT the boundary, so there is nothing to serve:
//
//	404 NOT_FOUND: Sync committee branch for block root 0x… not found. This typically
//	occurs when the block is not a finalized checkpoint.
//
// Observed on Sepolia with a finalized slot at offset 31 within its epoch; every earlier
// boundary answered 200. Client creation failed outright on that, and would keep failing
// for as long as the condition held, so walk back a boundary at a time until one is
// servable. An older checkpoint is a perfectly good trust anchor — it is still finalized,
// only slightly further back.
func resolveBootstrapCheckpoint(
	beaconAPIURL, finalizedSlot string,
	slotsPerEpoch uint64,
) (string, string, *relayerclient.BootstrapResponse, error) {
	slot, err := strconv.ParseUint(finalizedSlot, 10, 64)
	if err != nil {
		return "", "", nil, fmt.Errorf("parse finalized slot %q: %w", finalizedSlot, err)
	}
	if slotsPerEpoch == 0 {
		return "", "", nil, fmt.Errorf("slots_per_epoch is zero")
	}

	// stepBack moves to the previous epoch boundary, reporting false at the genesis
	// epoch where there is no earlier boundary to try. Subtracting unguarded would wrap
	// the unsigned slot around and send the next attempt at an absurd slot number.
	boundary := slot - slot%slotsPerEpoch
	stepBack := func() bool {
		if boundary < slotsPerEpoch {
			return false
		}
		boundary -= slotsPerEpoch
		return true
	}

	var lastErr error
	for attempt := 0; attempt < maxBootstrapCheckpointStepBack; attempt++ {
		candidate := strconv.FormatUint(boundary, 10)

		bctx, bcancel := context.WithTimeout(context.Background(), 15*time.Second)
		root, rootErr := relayerclient.GetBeaconBlockRoot(bctx, beaconAPIURL, candidate)
		bcancel()
		if rootErr != nil {
			// A skipped boundary slot has no block, so no root and no bootstrap.
			lastErr = fmt.Errorf("block root at slot %s: %w", candidate, rootErr)
			if !stepBack() {
				break
			}
			continue
		}

		bctx, bcancel = context.WithTimeout(context.Background(), 15*time.Second)
		bootstrap, bootErr := relayerclient.GetLightClientBootstrap(bctx, beaconAPIURL, root)
		bcancel()
		if bootErr != nil {
			lastErr = fmt.Errorf("bootstrap at slot %s (root %s): %w", candidate, root, bootErr)
			if !stepBack() {
				break
			}
			continue
		}
		if attempt > 0 {
			log.Printf("[CreateEthClient] finalized slot %s is not on a servable checkpoint; using slot %s (%d epoch(s) back)",
				finalizedSlot, candidate, attempt)
		}
		return candidate, root, bootstrap, nil
	}
	return "", "", nil, fmt.Errorf(
		"no servable light-client bootstrap within %d epochs below finalized slot %s; "+
			"the beacon may not serve light_client/bootstrap at all: %w",
		maxBootstrapCheckpointStepBack, finalizedSlot, lastErr)
}
