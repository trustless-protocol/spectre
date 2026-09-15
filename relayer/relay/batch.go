// This file holds one batch's journey: partition by relayability, prepare the
// client update, prove each packet, submit, and isolate a poison packet on
// failure. It is separated from module.go because handleBatch was 256 lines --
// five times the repo's 50-line limit -- inside a file whose other concern is
// the process lifecycle. C1's rule (every retry branch must change something and
// must be bounded) cannot be expressed while the gas ladder, the floor, and the
// error classification are interleaved in one function.
package relay

import (
	"context"
	"fmt"
	"log"
	"time"

	"relayer/chain"
)

// handleBatch relays one batch of source packets. When the destination supports
// update/packet folding, it builds packet proofs against the update's target
// height and submits the update followed by the packets in one atomic tx. Other
// destinations retain the update-then-packets fallback. It returns the indices
// of events that failed TRANSIENTLY and must be re-queued; successfully-relayed
// and permanently-dropped events are omitted.
func (m *Module) handleBatch(ctx context.Context, events []chain.Event) []int {
	// Serialized: see Module.batchMu. Nothing below re-enters handleBatch --
	// relayIsolated and relayFoldedIsolated are the only retry paths and neither
	// calls back into it -- so this cannot deadlock on itself.
	m.batchMu.Lock()
	defer m.batchMu.Unlock()
	return m.handleBatchLocked(ctx, events)
}

// handleBatchLocked is handleBatch's body, for callers that must do something
// else under the same lock. flushOnce is the one: its receipt check has to be
// serialized with submission, or a live batch can relay the packet in between
// and the flush submits a duplicate receive.
//
// Callers MUST hold batchMu.
func (m *Module) handleBatchLocked(ctx context.Context, events []chain.Event) []int {
	if len(events) == 0 {
		return nil
	}
	// Honor shutdown: if the context is already cancelled, re-queue everything
	// untouched rather than starting expensive proof/tx work that would run past
	// cancellation. The source re-queues these indices (nothing was relayed).
	if ctx.Err() != nil {
		return allIndices(len(events))
	}

	// Record every recv-relayed packet for timeout recovery BEFORE attempting the
	// relay: even if the update/proof/submit below fails, the scanner must still be
	// able to refund a packet that later expires undelivered. Only SendPacket opens
	// a receive that can time out; Add is idempotent.
	trackFailed := make([]bool, len(events))
	var requeue []int
	for i, e := range events {
		if e.Type == chain.SendPacket && m.track != nil {
			if !m.track(e.Raw, e.Height) {
				log.Printf("[relay %s] track pending packet seq=%d height=%d: durable write failed; re-queueing", m.name, e.Sequence, e.Height)
				trackFailed[i] = true
				requeue = append(requeue, i)
			}
		}
	}

	// Cheap precondition: the highest source height provable right now. A packet
	// above it (its source state not yet available — Cosmos AppHash H+2 lag, or ETH
	// beacon finality lag) is re-queued WITHOUT any client update or proof, so the
	// module does not burn an expensive Groth16 proof / beacon build every flush
	// while waiting. This replaces the legacy blocking waitCosmosAppHash /
	// waitBeaconFinality with a non-blocking value check.
	//
	// Held off while the source is in permanent failure. The events are already
	// TRACKED above, so a packet queued during the hold still gets its timeout
	// refund -- the hold skips the RPC and the proof, never the durable record.
	if m.sourceHeldOff(time.Now()) {
		m.lastWait = waitLogState{}
		return allIndices(len(events))
	}
	relayable, err := m.src.RelayableHeight(ctx)
	if err != nil {
		m.noteSourceFailure(err, "relayable height")
		// Forget the last wait: this flush learned nothing about the frontier, so a
		// wait that resumes identically afterwards is news again rather than a
		// repeat, and the error lines in between do not leave a silent gap.
		m.lastWait = waitLogState{}
		return allIndices(len(events))
	}
	m.noteSourceHealthy()

	var stalled []waiting // still-waiting packets, with how long each has waited
	provable := make([]chain.Event, 0, len(events))
	provableIdx := make([]int, 0, len(events)) // original index behind provable[j]
	for i, e := range events {
		if trackFailed[i] {
			continue
		}
		if e.Height > relayable {
			requeue = append(requeue, i) // not yet provable — wait, no expensive work
			stalled = append(stalled, waiting{event: e, age: m.waits.observe(e)})
			continue
		}
		// Deliberately NOT cleared here. Clearing at this gate would reset the age of
		// a packet that clears the source gate every flush but parks at the
		// destination-coverage gate below — the exact shape of a stuck L2 ack, and the
		// case the age exists to expose. An entry is cleared only when the packet
		// stops waiting for good: relayed, or dropped as permanently dead.
		provable = append(provable, e)
		provableIdx = append(provableIdx, i)
	}
	m.waits.purgeStale()
	// Surface the wait so the finality/AppHash lag is visible: relayable moving up
	// shows the source's head advancing toward the packet, and each packet is named
	// (type, sequence, height, age) because the count alone cannot answer the two
	// questions an operator has when a leg goes quiet — is my packet in this queue,
	// and has it been here long enough to be stuck rather than merely waiting?
	//
	// Logged only when that picture CHANGES. This runs once per flush — every batch
	// period, on no backoff — and a wait is normally long. Measured on OP Sepolia:
	// with the attestor's derived-root gap at 150 blocks the frontier sat still for
	// the whole wait, so 67 byte-identical lines collapse to 1; at a gap of 10 it
	// advanced, so 4 lines collapse to 3 — one per real advance.
	//
	// The description is part of the key, not just the counts: a packet whose age
	// crosses a reporting threshold is new information even when the frontier has
	// not moved, and suppressing that would hide the very case the age exists for.
	if len(stalled) > 0 {
		description := describeWaiting(stalled)
		now := waitLogState{packets: len(stalled), relayable: relayable, description: description, logged: true}
		if now != m.lastWait {
			log.Printf("[relay %s] waiting: %d packet(s) not yet relayable [%s] (source relayable height=%d)",
				m.name, len(stalled), description, relayable)
			m.lastWait = now
		}
	} else {
		m.lastWait = waitLogState{} // next wait reports itself, even if identical to the last one
	}
	if len(provable) == 0 {
		return requeue // nothing relayable this flush — skip the client update entirely
	}

	// Bail before the expensive client update / proof / submit if we are shutting
	// down; re-queue the provable packets so nothing is lost.
	if ctx.Err() != nil {
		return append(requeue, provableIdx...)
	}

	// Prepare the client state that packet proofs will target. A folding destination
	// leaves a needed update unsubmitted until the packet tx; the fallback submits
	// it immediately. The folded plan holds the client-update lock until its atomic
	// tx completes, preventing a refresh update from racing the planned height.
	foldPlan, proofHeight, err := m.prepareBatchUpdate(ctx, relayable)
	if err != nil {
		m.noteSourceFailure(err, fmt.Sprintf("update client to height %d", relayable))
		return append(requeue, provableIdx...)
	}
	if foldPlan != nil {
		defer foldPlan.release()
	}

	// Build proofs per packet. A proof-build failure is transient and per-packet:
	// re-queue just that event and keep the rest in the multicall (mirrors the
	// planCosmosPacketMsgs transient-failure routing).
	packets := make([]chain.RelayPacket, 0, len(provable))
	relayedIdx := make([]int, 0, len(provable)) // original event index behind packets[j]
	for j, e := range provable {
		i := provableIdx[j]
		// The RelayableHeight partition says the SOURCE can prove e.Height; it says
		// nothing about how far the DESTINATION client actually advanced. Those differ
		// whenever a client update commits a lower height than it was asked for — an
		// L2 update commits the block its attested game covers, which is usually below
		// the requested height. Proving against a client that does not cover the packet
		// is never valid, so re-queue and wait for an update that does.
		//
		// This is logged because it used to be silent: a packet parked here re-ran a
		// full client update every flush, forever, with nothing in the log. Rate-limit
		// via the same cadence as the waiting log if it ever gets noisy.
		if e.Height > proofHeight {
			log.Printf("[relay %s] waiting: %s seq=%d at height %d not yet covered by the destination client (trusts %d, waiting %s); re-queued",
				m.name, e.Type, e.Sequence, e.Height, proofHeight, m.waits.observe(e).Round(time.Second))
			requeue = append(requeue, i)
			continue
		}

		var (
			proof []byte
			err   error
		)
		if e.Type == chain.TimeoutPacket {
			proof, err = m.src.NonMembershipProof(ctx, e.Raw, proofHeight)
		} else {
			proof, err = m.src.MembershipProof(ctx, e.Raw, proofHeight, e.Type)
		}
		if err != nil {
			// A permanent proof failure means the packet is dead (timed out — its
			// commitment can never be proven for a receive). DROP it rather than
			// re-queueing: retrying re-runs a full client update each flush and
			// drains gas, and the timeout scanner (fed independently) refunds it.
			// Any other proof failure is transient (RPC blip, H+2 lag) → re-queue.
			if chain.IsPermanent(err) {
				// Dead for good — forget any recorded wait so the entry does not sit in
				// the tracker until its TTL.
				m.waits.clear(e)
				log.Printf("[relay %s] DROP packet (%s seq=%d height=%d): %v", m.name, e.Type, e.Sequence, e.Height, err)
				continue
			}
			// Every proof failure here re-queues, so logging one line per flush
			// says the same thing forever. What an operator actually needs is the
			// AGE: a proof failing for ten seconds is an RPC blip, one failing for
			// an hour is a packet nobody will ever deliver -- and the two used to
			// produce byte-identical lines.
			//
			// Whether the packet has another way out depends on its type, and the
			// STUCK line says which -- an operator's next action is different for
			// each, so a single wording would be wrong for one of them.
			//
			// A SendPacket does have one: once it is past its timeout every source
			// reports it permanent (evm/source.go, cosmos/source.go,
			// l2rollup/source.go all return chain.Permanent for an expired send),
			// the branch above drops it, and the timeout scanner refunds it. So
			// "no other exit" would be false for a send -- naming the real exit is
			// more useful than a warning to act.
			//
			// An AckPacket has none: its counterparty packet already has a receipt
			// on the destination, so it can never be timed out either, and its
			// escrow stays locked. A TimeoutPacket that cannot be proven is the
			// same -- nothing else refunds it. Both are worth saying out loud
			// rather than dropping, because dropping would lose the only record
			// that the funds are stranded.
			if age, report := m.waits.observeProofFailure(e); report {
				switch {
				case age < proofFailureStuckAfter:
					log.Printf("[relay %s] proof for packet (%s seq=%d height=%d, failing for %s): %v",
						m.name, e.Type, e.Sequence, e.Height, age.Round(time.Second), err)
				case e.Type == chain.SendPacket:
					log.Printf("[relay %s] STUCK: %s seq=%d height=%d has failed to prove for %s; it is refunded by the timeout scanner once past its timeout; last error: %v",
						m.name, e.Type, e.Sequence, e.Height, age.Round(time.Second), err)
				default:
					log.Printf("[relay %s] STUCK: %s seq=%d height=%d has failed to prove for %s and has no other exit, so its escrow stays locked; last error: %v",
						m.name, e.Type, e.Sequence, e.Height, age.Round(time.Second), err)
				}
			}
			requeue = append(requeue, i)
			continue
		}
		packets = append(packets, chain.RelayPacket{
			Type: e.Type, Sequence: e.Sequence, Packet: e.Raw,
			Proof: proof, Height: proofHeight, AckBytes: e.AckBytes,
		})
		relayedIdx = append(relayedIdx, i)
	}

	if len(packets) == 0 {
		// Nothing was provable this pass. Submit the folded client update on its
		// own before returning, or the light client never advances again.
		//
		// A folding destination leaves the update unsubmitted so it can ride in
		// the packet tx. When every proof fails, returning here drops it — and the
		// proofs are built at the height the client already trusts, so the client
		// staying put is exactly why they failed. That closes a loop the relayer
		// cannot leave on its own:
		//
		//   proofs target the trusted height  ->  that height falls out of the
		//   execution node's state window  ->  every proof fails  ->  no packets
		//   ->  the update is dropped  ->  the client stays put  ->  the gap only
		//   widens
		//
		// Observed on Sepolia: trustedSlot pinned at 10885728 for 20 minutes while
		// finalizedSlot advanced, every pass logging `eth_getProof failed:
		// historical state ... is not available`, with no update ever submitted.
		// Advancing the client is independent of whether any packet can be
		// relayed, so it must not be conditional on one.
		if foldPlan != nil {
			// Destination.UpdateClient, not the folding path: RelayWithUpdate is
			// specified to carry packets, and there are none.
			if err := m.dst.UpdateClient(ctx, m.clientID, foldPlan.update); err != nil {
				log.Printf("[relay %s] client update (no relayable packets) to height %d: %v",
					m.name, foldPlan.update.Height, err)
				return requeue
			}
			m.lastHeight = foldPlan.update.Height // foldPlan still holds m.mu
			m.recordClientUpdate(foldPlan.update)
			foldPlan.release()
			log.Printf("[relay %s] advanced client to height %d with no packets to relay",
				m.name, foldPlan.update.Height)
		}
		return requeue
	}

	// Last cancellation check before the tx: on shutdown, re-queue the proven
	// packets rather than broadcasting during teardown.
	if ctx.Err() != nil {
		return append(requeue, relayedIdx...)
	}

	// Submit the whole batch. A folded plan prepends the not-yet-on-chain client
	// update in the same atomic tx; otherwise the client is already current and the
	// destination submits packets only.
	if foldPlan != nil {
		err = foldPlan.destination.RelayWithUpdate(ctx, m.clientID, foldPlan.update, packets)
		if err == nil {
			m.lastHeight = foldPlan.update.Height // foldPlan still holds m.mu
			m.recordClientUpdate(foldPlan.update)
			foldPlan.release()
		}
	} else {
		err = m.dst.RelayPackets(ctx, packets)
	}
	if err == nil {
		m.settleDelivered(packets)
		return requeue
	}
	if foldPlan != nil {
		foldPlan.release()
	}
	if !chain.IsPermanent(err) {
		// Transient (RPC, nonce, broadcast) → re-queue the whole batch for retry.
		log.Printf("[relay %s] relay %d packet(s) at height %d: %v", m.name, len(packets), proofHeight, err)
		return append(requeue, relayedIdx...)
	}
	// Permanent (deterministic on-chain revert). A single packet is the poison —
	// DROP it (the timeout scanner refunds it if it later expires). But a multicall
	// is atomic: one bad packet (timed out / duplicate / already received) reverts
	// the whole tx, so dropping the batch would discard its VALID siblings too.
	// Isolate: retry each packet on its own, so the poison fails (and drops) alone
	// while its siblings still relay. This is the failure path only (rare), so the
	// extra per-packet txs are acceptable to avoid losing valid packets.
	if len(packets) == 1 {
		m.waits.clearPacket(packets[0]) // dead for good — stop reporting a wait for it
		log.Printf("[relay %s] DROP packet (%s seq=%d) at height %d (permanent): %v",
			m.name, packets[0].Type, packets[0].Sequence, proofHeight, err)
		return requeue
	}
	log.Printf("[relay %s] permanent batch failure at height %d (%v); isolating %d packet(s) individually", m.name, proofHeight, err, len(packets))
	if foldPlan != nil {
		return append(requeue, m.relayFoldedIsolated(ctx, foldPlan.destination, foldPlan.update, packets, relayedIdx)...)
	}
	return append(requeue, m.relayIsolated(ctx, packets, relayedIdx)...)
}

type foldedUpdatePlan struct {
	destination chain.FoldingDestination
	update      chain.ClientUpdate
	release     func()
}

// prepareBatchUpdate returns the proof height and, when supported and needed, a
// still-unsubmitted update to fold into the packet tx. A returned folded plan
// owns m.mu until release is called; release is idempotent.
func (m *Module) prepareBatchUpdate(ctx context.Context, height uint64) (*foldedUpdatePlan, uint64, error) {
	folding, ok := m.dst.(chain.FoldingDestination)
	if !ok || !folding.SupportsUpdatePacketFolding() {
		if err := m.updateClientTo(ctx, height); err != nil {
			return nil, 0, err
		}
		m.mu.Lock()
		proofHeight := m.lastHeight
		m.mu.Unlock()
		return nil, proofHeight, nil
	}

	m.mu.Lock()
	locked := true
	release := func() {
		if locked {
			locked = false
			m.mu.Unlock()
		}
	}
	fail := func(err error) (*foldedUpdatePlan, uint64, error) {
		release()
		return nil, 0, err
	}

	if height <= m.lastHeight {
		proofHeight := m.lastHeight
		release()
		return nil, proofHeight, nil
	}
	header, err := m.src.QueryHeader(ctx, height)
	if err != nil {
		return fail(fmt.Errorf("query header at %d: %w", height, err))
	}
	update, err := m.builder.Build(ctx, header)
	if err != nil {
		return fail(fmt.Errorf("build client update: %w", err))
	}
	if len(update.Payloads) == 0 {
		if update.Height > m.lastHeight {
			m.lastHeight = update.Height
		}
		m.recordClientUpdate(update)
		proofHeight := m.lastHeight
		release()
		return nil, proofHeight, nil
	}
	if update.Height <= m.lastHeight {
		proofHeight := m.lastHeight
		release()
		return nil, proofHeight, nil
	}

	return &foldedUpdatePlan{destination: folding, update: update, release: release}, update.Height, nil
}

// relayIsolated re-submits each packet of a permanently-failed batch on its own
// to isolate the poison packet from its valid siblings. Each singleton is
// classified independently: a permanent failure drops that packet alone (the
// timeout scanner refunds it), a transient failure re-queues its index, and a
// success untracks a delivered send. It returns the indices to re-queue.
func (m *Module) relayIsolated(ctx context.Context, packets []chain.RelayPacket, relayedIdx []int) []int {
	var requeue []int
	for j, p := range packets {
		if ctx.Err() != nil { // shutting down — re-queue the rest untouched
			requeue = append(requeue, relayedIdx[j:]...)
			break
		}
		single := []chain.RelayPacket{p}
		if err := m.dst.RelayPackets(ctx, single); err != nil {
			if chain.IsPermanent(err) {
				m.waits.clearPacket(p) // dead for good — stop reporting a wait for it
				log.Printf("[relay %s] DROP packet (%s seq=%d height=%d, permanent): %v", m.name, p.Type, p.Sequence, p.Height, err)
				continue
			}
			log.Printf("[relay %s] isolated relay (%s seq=%d height=%d): %v", m.name, p.Type, p.Sequence, p.Height, err)
			requeue = append(requeue, relayedIdx[j])
			continue
		}
		m.settleDelivered(single)
	}
	return requeue
}

// relayFoldedIsolated preserves poison-packet isolation after an atomic folded
// batch reverts. Until one singleton succeeds, each attempt carries the update;
// after that success installed the target height, remaining singletons use the
// packets-only path.
func (m *Module) relayFoldedIsolated(ctx context.Context, folding chain.FoldingDestination, update chain.ClientUpdate, packets []chain.RelayPacket, relayedIdx []int) []int {
	var requeue []int
	for j, p := range packets {
		if ctx.Err() != nil {
			requeue = append(requeue, relayedIdx[j:]...)
			break
		}

		single := []chain.RelayPacket{p}
		m.mu.Lock()
		needsUpdate := update.Height > m.lastHeight
		var err error
		if needsUpdate {
			err = folding.RelayWithUpdate(ctx, m.clientID, update, single)
			if err == nil {
				m.lastHeight = update.Height
				m.recordClientUpdate(update)
			}
			m.mu.Unlock()
		} else {
			m.mu.Unlock()
			err = m.dst.RelayPackets(ctx, single)
		}

		if err != nil {
			if chain.IsPermanent(err) {
				m.waits.clearPacket(p) // dead for good — stop reporting a wait for it
				log.Printf("[relay %s] DROP folded packet (%s seq=%d height=%d, permanent): %v", m.name, p.Type, p.Sequence, p.Height, err)
				continue
			}
			log.Printf("[relay %s] isolated folded relay (%s seq=%d height=%d): %v", m.name, p.Type, p.Sequence, p.Height, err)
			requeue = append(requeue, relayedIdx[j])
			continue
		}
		m.settleDelivered(single)
	}
	return requeue
}

// sourceHeldOff reports whether the source is inside a permanent-failure hold,
// and logs the first tick that ends one so the recovery is visible.
//
// Callers must hold batchMu.
func (m *Module) sourceHeldOff(now time.Time) bool {
	if m.sourceProbeAt.IsZero() {
		return false
	}
	if now.Before(m.sourceProbeAt) {
		return true
	}
	m.sourceProbeAt = time.Time{} // probe this pass; noteSourceFailure re-arms it
	return false
}

// noteSourceFailure logs a source failure and, when it is PERMANENT, holds the
// source off for a widening interval.
//
// The log line says which it was. Without that an operator sees the same
// "relayable height: ..." line for a route the attestor does not serve and for a
// replica that is a second behind, and the first looks like the second for as
// long as it takes someone to read the classifier.
//
// Callers must hold batchMu.
func (m *Module) noteSourceFailure(err error, what string) {
	if !chain.IsPermanent(err) {
		m.noteSourceHealthy() // a transient answer means the source is reachable
		log.Printf("[relay %s] %s: %v", m.name, what, err)
		return
	}
	// Same 1m->15m ladder the periodic client update uses; one cadence for
	// "this keeps failing, stop asking so often" rather than two.
	m.sourceBackoff = nextPeriodicUpdateBackoff(m.sourceBackoff)
	m.sourceProbeAt = time.Now().Add(m.sourceBackoff)
	log.Printf("[relay %s][ATTENTION] %s: %v; this is PERMANENT — no retry resolves it, "+
		"so the source is held off for %s rather than asked every pass. Queued packets are kept and "+
		"still time out normally; fix the configuration to clear it.",
		m.name, what, err, m.sourceBackoff)
}

// noteSourceHealthy clears the hold after the source answers.
//
// Callers must hold batchMu.
func (m *Module) noteSourceHealthy() {
	if m.sourceBackoff == 0 && m.sourceProbeAt.IsZero() {
		return
	}
	log.Printf("[relay %s] source answered again; permanent-failure hold cleared", m.name)
	m.sourceBackoff = 0
	m.sourceProbeAt = time.Time{}
}
