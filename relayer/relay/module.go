// Package relay drives one relay path (source -> destination) using only the
// chain adapter contract. It replaces the hardcoded Cosmos<->ETH orchestration
// in services.StartLoop: chain-specific work (RPC, event decode, proof/attest,
// gap recovery) lives inside the adapters, so this loop is generic and depends
// only on the chain interfaces — no god-object, testable with mock adapters.
package relay

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"relayer/chain"
)

// refreshMargin is how far ahead of a client's expiry the refresh routine
// proactively advances it, so a slow proof/tx cannot let the client expire.
const refreshMargin = 30 * time.Minute

// refreshTick is how often the refresh routine re-checks client expiry.
// It is a var, not a const, only so tests can drive the loop without waiting a
// minute; nothing outside the package reassigns it (same arrangement as
// l2rollup's l2SubscribeInterval).
var refreshTick = time.Minute

// defaultScanInterval matches the legacy StartLoop timeout scan cadence (30s).
const defaultScanInterval = 30 * time.Second

// Periodic-update failure backoff bounds (exponential 1m→15m), mirroring the legacy
// routineBackoff: a persistently failing client-update — e.g. a rotation tx that
// keeps reverting — must not resubmit every interval and drain the signer's gas.
const (
	periodicUpdateBackoffMin = time.Minute
	periodicUpdateBackoffMax = 15 * time.Minute
)

// defaultPeriodicUpdateInterval is the fallback force-rotation cadence when the
// caller does not derive one from the on-chain trusting period (matches the
// legacy DEFAULT_REFRESH_INTERVAL).
const defaultPeriodicUpdateInterval = 24 * time.Hour

// ScanFunc runs one timeout-recovery sweep for this path: it detects packets that
// were recv-relayed but expired undelivered and submits their MsgTimeout back to
// the source chain. It is called periodically and must handle its own errors (the
// wrapped scanners log and recover internally), so a transient RPC failure never
// stops the loop.
type ScanFunc func(ctx context.Context)

// TrackFunc records a recv-relayed packet — proto-marshaled bytes plus the source
// height it was observed at — so the path's ScanFunc can later refund it if it is
// never delivered. It returns false when the tracker could not durably record
// the packet; the module then re-queues the event without attempting its relay.
// Add is idempotent by packet identity.
type TrackFunc func(packet []byte, height uint64) bool

// UntrackFunc removes a packet (by its proto-marshaled bytes) from the pending
// tracker once the destination confirms delivery. Called after a successful relay
// so the timeout scanner stops considering a packet that can no longer time out.
type UntrackFunc func(packet []byte)

// ClientUpdateObserver records the source-chain timestamp trusted by a client
// after an update is submitted or the builder confirms it is current.
type ClientUpdateObserver func(trustedAt time.Time)

// PeriodicUpdateFunc runs one periodic destination-client update on a fixed
// cadence, independent of packet flow and of the expiry-driven refresh. It
// exists for clients whose freshness need is not captured by ClientExpiresAt
// alone — notably the SpectreClient pinned validator set, which must be rotated
// before its overlap decays below quorum even during quiet periods. It returns an
// error so the loop can retry sooner on failure.
type PeriodicUpdateFunc func(ctx context.Context) error

// Option configures optional Module capabilities (timeout scanning, pending-packet
// tracking) without breaking the base NewModule contract.
type Option func(*Module)

// WithTimeoutScanner runs scan every interval until Run's context is cancelled.
// interval <= 0 falls back to defaultScanInterval.
func WithTimeoutScanner(interval time.Duration, scan ScanFunc) Option {
	return func(m *Module) {
		if interval <= 0 {
			interval = defaultScanInterval
		}
		m.scan = scan
		m.scanInterval = interval
	}
}

// WithPacketTracker records every recv-relayed (SendPacket) packet via track, so
// the path's timeout scanner can refund it if delivery never completes. Tracking
// happens regardless of relay outcome (mirrors the legacy "track every send"
// invariant), so a packet that fails to relay and then times out is still
// refunded. untrack removes a packet once the relay succeeds — after delivery it
// can no longer time out, so keeping it in the tracker only bloats it and makes
// the scanner query a receipt it will always find (mirrors the legacy handleCosmos
// PendingTracker.Remove-on-recv). untrack may be nil (no removal hook).
func WithPacketTracker(track TrackFunc, untrack UntrackFunc) Option {
	return func(m *Module) {
		m.track = track
		m.untrack = untrack
	}
}

// WithClientUpdateObserver exposes successful client progress to observability
// without coupling the generic relay loop to a particular reporter.
func WithClientUpdateObserver(observer ClientUpdateObserver) Option {
	return func(m *Module) { m.observeClientUpdate = observer }
}

// WithPeriodicUpdate runs a forced client update every interval (independent of
// the expiry-driven refresh) to keep the destination client fresh in ways
// ClientExpiresAt does not capture — e.g. the SpectreClient pinned-set rotation.
// interval <= 0 falls back to defaultPeriodicUpdateInterval.
//
// initialDelay is how long to wait before the FIRST update, derived from the
// client's on-chain freshness (time until it is next due), not process uptime.
// A restart near the rotation deadline therefore fires promptly instead of
// waiting a whole fresh interval — the guarantee that keeps the pinned set above
// quorum. A negative initialDelay is treated as 0 (due now).
func WithPeriodicUpdate(interval, initialDelay time.Duration, periodicUpdate PeriodicUpdateFunc) Option {
	return func(m *Module) {
		if interval <= 0 {
			interval = defaultPeriodicUpdateInterval
		}
		if initialDelay < 0 {
			initialDelay = 0
		}
		m.periodicUpdate = periodicUpdate
		m.periodicUpdateInterval = interval
		m.periodicUpdateInitialDelay = initialDelay
	}
}

// Module drives a single source -> destination relay path.
type Module struct {
	name     string
	clientID string
	src      chain.Source
	dst      chain.Destination
	builder  chain.ClientUpdateBuilder

	// Optional timeout-recovery wiring (nil when the path has no timeout scanner,
	// e.g. a mock or a permissioned client that cannot time out).
	scan                ScanFunc
	scanInterval        time.Duration
	track               TrackFunc
	untrack             UntrackFunc
	observeClientUpdate ClientUpdateObserver

	// Optional fixed-cadence forced update (nil when the destination client's
	// freshness is fully captured by ClientExpiresAt, e.g. the beacon client).
	periodicUpdate             PeriodicUpdateFunc
	periodicUpdateInterval     time.Duration
	periodicUpdateInitialDelay time.Duration

	// waits records how long each parked packet has been waiting, so the waiting
	// logs report an age and not just a count. Purely observational — it never
	// gates a relay decision.
	waits *waitTracker

	// lastHeight is the highest source height whose ClientUpdate we have
	// submitted. It enforces the append-only / monotonic invariant across the
	// event and refresh goroutines, so a lower or stale update (e.g. a builder
	// returning "no update needed") is skipped rather than replayed.
	mu         sync.Mutex
	lastHeight uint64

	// lastWait is the last "not yet relayable" state logged, so a wait that is not
	// progressing prints once instead of once per flush. Touched only from the
	// Subscribe callback goroutine (handleBatch), which is the sole caller.
	lastWait waitLogState
}

// waitLogState is the identity of a wait: an identical value means nothing changed
// since the last line, so there is nothing new to tell the operator.
//
// description carries the per-packet detail (type, sequence, height, age), so a
// packet whose reported age advances re-logs even while the frontier is static —
// which is the case the age was added to expose.
type waitLogState struct {
	packets     int
	relayable   uint64
	description string
	logged      bool
}

// NewModule wires a relay path. name is a human label for logs; clientID is the
// destination client this path advances. Optional timeout-recovery behavior is
// added via WithTimeoutScanner / WithPacketTracker.
func NewModule(name, clientID string, src chain.Source, dst chain.Destination, builder chain.ClientUpdateBuilder, opts ...Option) *Module {
	m := &Module{
		name: name, clientID: clientID, src: src, dst: dst, builder: builder,
		waits: newWaitTracker(),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Run drives the path until ctx is cancelled: it subscribes to source packet
// events (the adapter owns gap recovery) and runs a background refresh routine
// that keeps the destination client from expiring.
func (m *Module) Run(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	workers := newWorkerGroup()
	subscribeErr := make(chan error, 1)

	workers.Go("subscriber", func() {
		subscribeErr <- m.src.Subscribe(runCtx, m.handleBatch)
	})
	workers.Go("refresh", func() { m.refreshLoop(runCtx) })
	if m.scan != nil {
		workers.Go("timeout-scan", func() { m.scanLoop(runCtx) })
	}
	if m.periodicUpdate != nil {
		workers.Go("periodic-update", func() { m.periodicUpdateLoop(runCtx) })
	}

	var runErr error
	select {
	case <-ctx.Done():
		runErr = ctx.Err()
	case err := <-subscribeErr:
		runErr = err
	}
	cancel()

	if running := workers.drain(shutdownDrainTimeout); len(running) > 0 {
		return fmt.Errorf("relay %s: shutdown drain timed out after %s; workers still running: %v", m.name, shutdownDrainTimeout, running)
	}
	if runErr != nil {
		if ctx.Err() != nil && (errors.Is(runErr, context.Canceled) || errors.Is(runErr, context.DeadlineExceeded)) {
			return nil
		}
		return fmt.Errorf("relay %s: subscribe: %w", m.name, runErr)
	}
	return nil
}

// scanLoop periodically runs the timeout-recovery sweep so packets that were
// recv-relayed but never delivered are refunded on the source chain before the
// pending set grows unbounded. It is a no-op path when no scanner is configured
// (Run does not start it in that case).
func (m *Module) scanLoop(ctx context.Context) {
	ticker := time.NewTicker(m.scanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ctx.Err() != nil { // select may pick the tick even when ctx is done
				return
			}
			m.scan(ctx)
		}
	}
}

// nextPeriodicUpdateBackoff advances the exponential failure backoff: the first
// failure waits periodicUpdateBackoffMin, each subsequent failure doubles up to
// periodicUpdateBackoffMax. current == 0 means "no prior failure".
func nextPeriodicUpdateBackoff(current time.Duration) time.Duration {
	if current == 0 {
		return periodicUpdateBackoffMin
	}
	return min(current*2, periodicUpdateBackoffMax)
}

// periodicUpdateLoop runs the fixed-cadence forced client update (e.g. the
// pinned-set rotation). It fires on periodicUpdateInterval on success and on an
// exponential backoff (periodicUpdateBackoffMin→Max) after a failure, so a
// transient failure of the infrequent update is retried promptly rather than
// waiting a full interval.
func (m *Module) periodicUpdateLoop(ctx context.Context) {
	// The first fire uses the on-chain-freshness-derived initial delay (0 = due
	// now); every fire after that uses the full interval (a success resets the
	// on-chain timestamp to now) or the failure backoff.
	timer := time.NewTimer(m.periodicUpdateInitialDelay)
	defer timer.Stop()
	var backoff time.Duration // 0 = healthy; grows on consecutive failures
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			// select picks a ready case at random, so the timer branch can win even
			// when ctx is already done — re-check so shutdown never starts a rotation
			// tx (the wired RotatePinnedSet does not yet honor ctx itself).
			if ctx.Err() != nil {
				return
			}
			next := m.periodicUpdateInterval
			if err := m.periodicUpdate(ctx); err != nil {
				// Exponential backoff on repeated failure (1m→15m) so a persistently
				// failing client-update tx does not resubmit every interval and drain
				// the signer's gas (the legacy routineBackoff). Never slower than the
				// normal cadence (guards a configuration with a shorter interval).
				backoff = nextPeriodicUpdateBackoff(backoff)
				log.Printf("[relay %s] periodic update: %v; retrying after %s", m.name, err, backoff)
				next = min(backoff, m.periodicUpdateInterval)
			} else {
				backoff = 0 // success resets the backoff
			}
			timer.Reset(next)
		}
	}
}

// handleBatch relays one batch of source packets. When the destination supports
// update/packet folding, it builds packet proofs against the update's target
// height and submits the update followed by the packets in one atomic tx. Other
// destinations retain the update-then-packets fallback. It returns the indices
// of events that failed TRANSIENTLY and must be re-queued; successfully-relayed
// and permanently-dropped events are omitted.
func (m *Module) handleBatch(ctx context.Context, events []chain.Event) []int {
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
	relayable, err := m.src.RelayableHeight(ctx)
	if err != nil {
		log.Printf("[relay %s] relayable height: %v", m.name, err)
		// Forget the last wait: this flush learned nothing about the frontier, so a
		// wait that resumes identically afterwards is news again rather than a
		// repeat, and the error lines in between do not leave a silent gap.
		m.lastWait = waitLogState{}
		return allIndices(len(events))
	}

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
		log.Printf("[relay %s] update client to height %d: %v", m.name, relayable, err)
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

// settleDelivered closes out packets the destination accepted.
//
// It clears each one's wait entry — the packet landed, so any age recorded for
// it must not be reported again (nor inflate a later packet that reuses the same
// sequence after a client re-creation).
//
// It also removes each delivered SendPacket from the pending tracker: once
// received on the destination it can no longer time out, so leaving it in the
// tracker only bloats it and makes the timeout scanner keep querying a receipt
// it will always find. Mirrors the legacy handleCosmos
// PendingTracker.Remove-on-recv. Only paths that supply an untracker
// (cosmos->eth) do this; eth->cosmos removes via its source's terminal
// EthAck/EthTimeout instead.
func (m *Module) settleDelivered(packets []chain.RelayPacket) {
	for _, p := range packets {
		m.waits.clearPacket(p)
		if m.untrack != nil && p.Type == chain.SendPacket {
			m.untrack(p.Packet)
		}
	}
}

// allIndices returns [0, n) — every position in a batch (whole-batch re-queue).
func allIndices(n int) []int {
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	return idx
}

// updateClientTo advances the destination client so it covers source height,
// honoring the append-only invariant. It is safe to call from multiple
// goroutines.
func (m *Module) updateClientTo(ctx context.Context, height uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if height <= m.lastHeight {
		return nil // already covered — nothing to prove
	}

	header, err := m.src.QueryHeader(ctx, height)
	if err != nil {
		return fmt.Errorf("query header at %d: %w", height, err)
	}
	update, err := m.builder.Build(ctx, header)
	if err != nil {
		return fmt.Errorf("build client update: %w", err)
	}
	// An empty payload list means "client already current"; the builder still reports
	// its current height, so learn it (seeding lastHeight from on-chain reality)
	// without submitting a tx. This keeps the provability guard correct across a
	// restart, where the client may already cover incoming packets and every
	// Build would otherwise report a no-op the module could not learn from.
	if len(update.Payloads) == 0 {
		if update.Height > m.lastHeight {
			m.lastHeight = update.Height
		}
		m.recordClientUpdate(update)
		return nil
	}
	// A stale update (at or below what we already trust) is skipped without a tx.
	if update.Height <= m.lastHeight {
		return nil
	}

	if err := m.dst.UpdateClient(ctx, m.clientID, update); err != nil {
		return fmt.Errorf("submit client update at %d: %w", update.Height, err)
	}
	m.lastHeight = update.Height
	m.recordClientUpdate(update)
	return nil
}

func (m *Module) recordClientUpdate(update chain.ClientUpdate) {
	if m.observeClientUpdate != nil && !update.TrustedAt.IsZero() {
		m.observeClientUpdate(update.TrustedAt)
	}
}

// needsRefresh reports whether the destination client is close enough to expiry
// that the anti-expiry routine must advance it now.
//
// A zero expiry means "never expires" (e.g. a permissioned client) and needs no
// refresh. Otherwise the client is refreshed once it is within refreshMargin of
// expiring — note the direction: MORE than a margin of headroom means there is
// nothing to do yet, and that is the comparison worth pinning, because inverting
// it produces a routine that refreshes only while there is plenty of time and
// goes quiet exactly when the client is about to expire.
func needsRefresh(expiresAt, now time.Time) bool {
	if expiresAt.IsZero() {
		return false
	}
	return expiresAt.Sub(now) <= refreshMargin
}

// refreshLoop proactively advances the destination client before it expires,
// covering quiet periods with no packet traffic (the anti-expiry routine).
func (m *Module) refreshLoop(ctx context.Context) {
	ticker := time.NewTicker(refreshTick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ctx.Err() != nil { // select may pick the tick even when ctx is done
				return
			}
			expiresAt, err := m.dst.ClientExpiresAt(ctx, m.clientID)
			if err != nil {
				log.Printf("[relay %s] query client expiry: %v", m.name, err)
				continue
			}
			if !needsRefresh(expiresAt, time.Now()) {
				continue
			}
			latest, err := m.src.LatestHeight(ctx)
			if err != nil {
				log.Printf("[relay %s] latest source height for refresh: %v", m.name, err)
				continue
			}
			if err := m.updateClientTo(ctx, latest); err != nil {
				log.Printf("[relay %s] refresh update to %d: %v", m.name, latest, err)
			}
		}
	}
}
