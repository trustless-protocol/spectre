// This file is the beacon side of the Ethereum light client: sync-committee
// lookup, bootstrap checkpoint resolution, and turning a finality update into the
// proof state an update message needs.
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	relayerclient "relayer/client"
	"sort"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func cosmosCurrentSlotReady(currentSlot, sigSlot uint64) bool {
	return currentSlot >= sigSlot+cosmosCatchUpSafetySlots
}

// executionHeaderRLPHex returns the hex RLP execution header the wasm client needs for a finalized
// header in the Gloas format, or "" for a pre-Gloas-format header, which carries its own payload
// header. Fetching it here also checks the header against its proven block hash, so a node that
// serves the wrong bytes fails now rather than after a transaction has been paid for.
func executionHeaderRLPHex(ctx context.Context, source relayerclient.RawHeaderSource, finalized relayerclient.LightClientHeader) (string, error) {
	if !finalized.IsGloasFormat() {
		return "", nil
	}
	_, raw, err := relayerclient.ResolveExecution(ctx, source, finalized)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(raw), nil
}

// maxLightClientUpdatesPerRequest is the most sync-committee periods one client
// update may cross. It is the beacon API's own limit on one
// light_client/updates request (MAX_REQUEST_LIGHT_CLIENT_UPDATES).
const maxLightClientUpdatesPerRequest = 128

// lightClientUpdateCount is how many light-client updates take the client from
// trustedPeriod to targetPeriod, one per period touched.
//
// Both periods derive from slots the RPCs supplied, and the count sizes a slice.
// Unchecked, a finalized slot far in the future asked makeslice for more memory
// than exists and panicked (#575 L8); a target behind the client wrapped around
// to the same thing.
func lightClientUpdateCount(trustedPeriod, targetPeriod uint64) (uint64, error) {
	if targetPeriod < trustedPeriod {
		return 0, fmt.Errorf("finalized sync-committee period %d is behind the client's period %d", targetPeriod, trustedPeriod)
	}
	if gap := targetPeriod - trustedPeriod; gap >= maxLightClientUpdatesPerRequest {
		return 0, fmt.Errorf(
			"finalized sync-committee period %d is %d periods past the client's period %d; one update crosses at most %d",
			targetPeriod, gap, trustedPeriod, maxLightClientUpdatesPerRequest-1)
	}
	return targetPeriod - trustedPeriod + 1, nil
}

// ethUpdateBuild is one client update: the headers; the finalized header of the
// LAST one, the state the client holds once they land; and the highest signature
// slot, which the destination's clock must reach before every header verifies.
type ethUpdateBuild struct {
	headers              [][]byte
	lastFinalized        relayerclient.LightClientHeader
	lastFinalizedSlot    uint64
	highestSignatureSlot uint64
}

func (w *Worker) buildEthClientUpdateHeadersWithPeriodCrossing(stdCtx context.Context, beaconAPIURL string, source relayerclient.RawHeaderSource, ethClientState *relayerclient.EthereumClientState, trustedSlot, trustedPeriod, targetPeriod uint64, finalityUpdate *relayerclient.LightClientFinalityUpdate, finalizedSlot uint64) (ethUpdateBuild, error) {
	count, err := lightClientUpdateCount(trustedPeriod, targetPeriod)
	if err != nil {
		return ethUpdateBuild{}, err
	}
	bctx, bcancel := context.WithTimeout(stdCtx, 15*time.Second)
	lightClientUpdates, err := relayerclient.GetLightClientUpdates(bctx, beaconAPIURL, trustedPeriod, count)
	bcancel()
	if err != nil {
		return ethUpdateBuild{}, fmt.Errorf("failed to get light client updates: %w", err)
	}

	if len(lightClientUpdates) == 0 {
		return ethUpdateBuild{}, fmt.Errorf("no light client updates available for period range %d to %d", trustedPeriod, targetPeriod)
	}

	// Index the fetched updates by the period they belong to. Crossing into period P
	// needs the FULL sync committee of period P, and the update for period P-1 already
	// carries it as next_sync_committee — Merkle-proven against its attested header,
	// which is exactly what next_sync_committee_branch exists for. Reading it from
	// here avoids a light_client/bootstrap call that beacon nodes only answer for the
	// checkpoint roots they happen to retain (see the fallback below).
	type slottedUpdate struct {
		update    relayerclient.LightClientUpdate
		finalized uint64
	}
	updatesByPeriod := make(map[uint64]relayerclient.LightClientUpdate, len(lightClientUpdates))
	sorted := make([]slottedUpdate, 0, len(lightClientUpdates))
	for _, u := range lightClientUpdates {
		slot, err := parseSlot(u.FinalizedHeader.Beacon.Slot)
		if err != nil {
			return ethUpdateBuild{}, fmt.Errorf("failed to parse update finalized slot: %w", err)
		}
		updatesByPeriod[ethClientState.ComputeSyncCommitteePeriodAtSlot(slot)] = u
		sorted = append(sorted, slottedUpdate{update: u, finalized: slot})
	}
	// In slot order, whatever order the node answered in: each header is checked
	// against the state the one before it leaves (#575).
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].finalized < sorted[j].finalized })

	build := ethUpdateBuild{headers: make([][]byte, 0, count)}
	latestTrustedSlot := trustedSlot
	latestPeriod := trustedPeriod

	for _, entry := range sorted {
		update, updateFinalizedSlot := entry.update, entry.finalized

		if updateFinalizedSlot <= latestTrustedSlot {
			continue
		}

		updatePeriod := ethClientState.ComputeSyncCommitteePeriodAtSlot(updateFinalizedSlot)
		if updatePeriod == latestPeriod {
			continue
		}
		// The client verifies one period at a time: a crossing header must finalize
		// in the next period and be signed in it, by the committee the client stores
		// as next. It rejects anything else; name the cause here instead.
		if updatePeriod != latestPeriod+1 {
			return ethUpdateBuild{}, fmt.Errorf("beacon light-client updates skip from period %d to %d (finalized slot %d); the client crosses one period at a time",
				latestPeriod, updatePeriod, updateFinalizedSlot)
		}
		signatureSlot, err := parseSlot(update.SignatureSlot)
		if err != nil {
			return ethUpdateBuild{}, fmt.Errorf("update finalized at slot %d: signature slot: %w", updateFinalizedSlot, err)
		}
		if signaturePeriod := ethClientState.ComputeSyncCommitteePeriodAtSlot(signatureSlot); signaturePeriod != updatePeriod {
			return ethUpdateBuild{}, fmt.Errorf("update finalized at slot %d (period %d) is signed at slot %d (period %d); the client can verify only a signature in period %d",
				updateFinalizedSlot, updatePeriod, signatureSlot, signaturePeriod, updatePeriod)
		}

		// The bootstrap fallback is asked at the finalized slot: a finalized
		// checkpoint, the only kind of block nodes serve a bootstrap for (see
		// resolveBootstrapCheckpoint). The signature slot is in the same period,
		// checked above, so the committee is the same; its block is not a
		// checkpoint and is answered 404.
		syncCommittee, err := syncCommitteeForPeriod(stdCtx, beaconAPIURL, updatesByPeriod, updatePeriod, updateFinalizedSlot)
		if err != nil {
			return ethUpdateBuild{}, err
		}

		executionRLP, err := executionHeaderRLPHex(stdCtx, source, update.FinalizedHeader)
		if err != nil {
			return ethUpdateBuild{}, fmt.Errorf("execution header for finalized slot %d: %w", updateFinalizedSlot, err)
		}
		header := relayerclient.EthereumHeader{
			ActiveSyncCommittee: relayerclient.ActiveSyncCommittee{
				Next: &syncCommittee,
			},
			ConsensusUpdate:    update,
			TrustedSlot:        latestTrustedSlot,
			ExecutionHeaderRLP: executionRLP,
		}

		headerBytes, err := json.Marshal(header)
		if err != nil {
			return ethUpdateBuild{}, fmt.Errorf("failed to marshal update header: %w", err)
		}

		build.headers = append(build.headers, headerBytes)
		latestPeriod = updatePeriod
		latestTrustedSlot = updateFinalizedSlot
		build.lastFinalized, build.lastFinalizedSlot = update.FinalizedHeader, updateFinalizedSlot
		build.highestSignatureSlot = max(build.highestSignatureSlot, signatureSlot)
	}

	// If the latest header is earlier than the finality update, add a header for the finality update.
	if finalizedSlot > latestTrustedSlot {
		attestedSlot := finalityUpdate.AttestedHeader.Beacon.Slot
		log.Printf("[UpdateEthClient] final update: attestedSlot=%s finalizedSlot=%d latestTrustedSlot=%d",
			attestedSlot, finalizedSlot, latestTrustedSlot)

		// The finality update is signed by the committee of its signature slot's
		// period, which the client checks against a committee it already trusts.
		// Same sourcing problem as the crossing loop above: after crossing into a new
		// period the beacon will not serve a bootstrap for that period, so prefer the
		// committee carried by the preceding period's update.
		//
		// The client selects the committee by the SIGNATURE slot's period, not the
		// attested slot's (ethereum/light-client verify.rs: signature_period =
		// compute_sync_committee_period_at_slot(update.signature_slot)). They differ
		// for exactly one slot -- signature_slot is attested_slot + 1 -- and that one
		// slot is a period boundary once every 8192.
		signatureSlotNum, err := parseSlot(finalityUpdate.SignatureSlot)
		if err != nil {
			return ethUpdateBuild{}, fmt.Errorf("failed to parse signature slot: %w", err)
		}
		signaturePeriod := ethClientState.ComputeSyncCommitteePeriodAtSlot(signatureSlotNum)
		// The bootstrap fallback reads the committee of the block it is asked for, so
		// it is asked for the signature slot: the attested slot is one period early
		// at a boundary (#575).
		syncCommittee, err := syncCommitteeForPeriod(
			stdCtx, beaconAPIURL, updatesByPeriod, signaturePeriod, signatureSlotNum)
		if err != nil {
			return ethUpdateBuild{}, err
		}
		activeSyncCommittee, err := activeCommitteeFor(signaturePeriod, latestPeriod, &syncCommittee)
		if err != nil {
			return ethUpdateBuild{}, err
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

		executionRLP, err := executionHeaderRLPHex(stdCtx, source, finalityUpdate.FinalizedHeader)
		if err != nil {
			return ethUpdateBuild{}, fmt.Errorf("execution header for finalized slot %d: %w", finalizedSlot, err)
		}
		header := relayerclient.EthereumHeader{
			ActiveSyncCommittee: activeSyncCommittee,
			ConsensusUpdate:     consensusUpdate,
			TrustedSlot:         latestTrustedSlot,
			ExecutionHeaderRLP:  executionRLP,
		}

		headerBytes, err := json.Marshal(header)
		if err != nil {
			return ethUpdateBuild{}, fmt.Errorf("failed to marshal update header: %w", err)
		}

		build.headers = append(build.headers, headerBytes)
		build.lastFinalized, build.lastFinalizedSlot = finalityUpdate.FinalizedHeader, finalizedSlot
		build.highestSignatureSlot = max(build.highestSignatureSlot, signatureSlotNum)
	}

	return build, nil
}

// activeCommitteeFor labels the committee the way the light client will look it
// up, given the period its stored consensus state is in.
//
// The client does this (ethereum/light-client verify.rs):
//
//	sync_committee = if signature_period == stored_period { current } else { next }
//
// so the label is not cosmetic -- it selects which stored summary the supplied
// committee is checked against, and a wrong label fails with
// "current sync committee (X) does not match with the one in the current state (Y)".
//
// The two are not always equal. signature_slot runs roughly 65 slots ahead of the
// finalized slot the client stores, so for the last ~65 slots of every
// sync-committee period the signature is already in the NEXT period while the
// client is still finalized in this one. Labelling that Current is what stalled
// the ETH->Cosmos client for ~13 minutes once per period.
//
// Next is safe there and does not rotate anything: the client keeps
// next_sync_committee across same-period updates (update.rs only rotates when the
// update's FINALIZED period advances), so the committee this names is the one it
// already stores, and the rotation still happens later on the finalized crossing.
func activeCommitteeFor(signaturePeriod, storedPeriod uint64, committee *relayerclient.SyncCommittee) (relayerclient.ActiveSyncCommittee, error) {
	switch signaturePeriod {
	case storedPeriod:
		return relayerclient.ActiveSyncCommittee{Current: committee}, nil
	case storedPeriod + 1:
		return relayerclient.ActiveSyncCommittee{Next: committee}, nil
	default:
		// The client rejects this outright (InvalidSignaturePeriodWhenNextSyncCommitteeExists),
		// so building the header would only spend gas to be told so.
		return relayerclient.ActiveSyncCommittee{}, fmt.Errorf(
			"sync committee period %d is neither the client's period %d nor the one after it; the client cannot verify this update",
			signaturePeriod, storedPeriod)
	}
}

func cloneEthereumClientState(state *relayerclient.EthereumClientState) *relayerclient.EthereumClientState {
	if state == nil {
		return nil
	}
	cloned := *state
	return &cloned
}

// ethProofStateFromHeader is the client state once the header finalized at
// finalizedSlot lands: the state packet proofs are built against.
func ethProofStateFromHeader(ctx context.Context, source relayerclient.RawHeaderSource, base *relayerclient.EthereumClientState, finalized relayerclient.LightClientHeader, finalizedSlot uint64) (*relayerclient.EthereumClientState, uint64, error) {
	proofState := cloneEthereumClientState(base)
	if proofState == nil {
		return nil, 0, fmt.Errorf("ethereum client state is nil")
	}

	// The same values the wasm client stores when it applies this update, so proofs are built
	// against the execution block the client will trust.
	execution, _, err := relayerclient.ResolveExecution(ctx, source, finalized)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to resolve finalized execution header: %w", err)
	}

	proofState.LatestSlot = finalizedSlot
	proofState.LatestExecutionBlockNumber = execution.BlockNumber
	return proofState, execution.Timestamp, nil
}

// parseSlot parses a slot string to uint64
func parseSlot(slotStr string) (uint64, error) {
	slot, err := strconv.ParseUint(slotStr, 10, 64)
	if err != nil {
		// No %q on slotStr: strconv already quotes the offending input, and
		// repeating it prints the same string twice in one line.
		return 0, fmt.Errorf("parse beacon slot: %w", err)
	}
	return slot, nil
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
	period, bootstrapSlot uint64,
) (relayerclient.SyncCommittee, error) {
	prevErr := fmt.Errorf("period 0 has no preceding update")
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
			prevErr = fmt.Errorf("the period %d update carries no next_sync_committee", period-1)
		} else {
			// Kept for the error below: a failed fetch is a different cause from an
			// update that simply has no committee, and it used to be dropped (#575).
			prevErr = fmt.Errorf("fetch the period %d update: %w", period-1, err)
		}
	}

	bctx, bcancel := context.WithTimeout(stdCtx, 15*time.Second)
	blockRoot, err := relayerclient.GetBeaconBlockRoot(bctx, beaconAPIURL, fmt.Sprintf("%d", bootstrapSlot))
	bcancel()
	if err != nil {
		return relayerclient.SyncCommittee{}, fmt.Errorf(
			"period %d: no preceding update gave next_sync_committee (%v) and beacon block root lookup failed: %w", period, prevErr, err)
	}

	bctx, bcancel = context.WithTimeout(stdCtx, 15*time.Second)
	bootstrap, err := relayerclient.GetLightClientBootstrap(bctx, beaconAPIURL, blockRoot)
	bcancel()
	if err != nil {
		return relayerclient.SyncCommittee{}, fmt.Errorf(
			"period %d: no preceding update gave next_sync_committee (%v) and bootstrap at slot %d is unavailable: %w",
			period, prevErr, bootstrapSlot, err)
	}
	return bootstrap.Data.CurrentSyncCommittee, nil
}

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
