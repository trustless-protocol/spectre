// Package relay drives one relay path (source -> destination) using only the
// chain adapter contract. It replaces the hardcoded Cosmos<->ETH orchestration
// in services.StartLoop: chain-specific work (RPC, event decode, proof/attest,
// gap recovery) lives inside the adapters, so this loop is generic and depends
// only on the chain interfaces — no god-object, testable with mock adapters.
package relay

import (
	"context"
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
const refreshTick = time.Minute

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
// never delivered. Add is idempotent by (source client, sequence).
type TrackFunc func(packet []byte, height uint64)

// UntrackFunc removes a packet (by its proto-marshaled bytes) from the pending
// tracker once the destination confirms delivery. Called after a successful relay
// so the timeout scanner stops considering a packet that can no longer time out.
type UntrackFunc func(packet []byte)

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
	scan         ScanFunc
	scanInterval time.Duration
	track        TrackFunc
	untrack      UntrackFunc

	// Optional fixed-cadence forced update (nil when the destination client's
	// freshness is fully captured by ClientExpiresAt, e.g. the beacon client).
	periodicUpdate             PeriodicUpdateFunc
	periodicUpdateInterval     time.Duration
	periodicUpdateInitialDelay time.Duration

	// lastHeight is the highest source height whose ClientUpdate we have
	// submitted. It enforces the append-only / monotonic invariant across the
	// event and refresh goroutines, so a lower or stale update (e.g. a builder
	// returning "no update needed") is skipped rather than replayed.
	mu         sync.Mutex
	lastHeight uint64
}

// NewModule wires a relay path. name is a human label for logs; clientID is the
// destination client this path advances. Optional timeout-recovery behavior is
// added via WithTimeoutScanner / WithPacketTracker.
func NewModule(name, clientID string, src chain.Source, dst chain.Destination, builder chain.ClientUpdateBuilder, opts ...Option) *Module {
	m := &Module{name: name, clientID: clientID, src: src, dst: dst, builder: builder}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Run drives the path until ctx is cancelled: it subscribes to source packet
// events (the adapter owns gap recovery) and runs a background refresh routine
// that keeps the destination client from expiring.
func (m *Module) Run(ctx context.Context) error {
	go m.refreshLoop(ctx)
	if m.scan != nil {
		go m.scanLoop(ctx)
	}
	if m.periodicUpdate != nil {
		go m.periodicUpdateLoop(ctx)
	}

	// Subscribe carries the adapter's gap-recovery guarantees; it returns only on
	// ctx cancellation or a fatal subscription error.
	if err := m.src.Subscribe(ctx, m.handleBatch); err != nil {
		return fmt.Errorf("relay %s: subscribe: %w", m.name, err)
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

// handleBatch relays one batch of source packets: advance the destination client
// once to cover the whole batch, prove each packet against that trusted height,
// then submit them in a single multicall (folding the per-tx intrinsic gas over
// the batch). It returns the indices of events that failed TRANSIENTLY and must
// be re-queued; successfully-relayed and permanently-dropped events are omitted.
//
// The client update and the packet submit are currently SEPARATE txs (updateClientTo
// then RelayPackets), not one folded multicall like the legacy path — so the client
// state can advance even when the subsequent packet submit fails (the packets are
// re-queued and retried against the now-current client; folding them back into one
// tx is a tracked follow-up). A packet's tracker entry / cursor is never advanced on
// failure.
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
	for _, e := range events {
		if e.Type == chain.SendPacket && m.track != nil {
			m.track(e.Raw, e.Height)
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
		return allIndices(len(events))
	}

	var requeue []int
	var pendingMax uint64 // highest source height still waiting (for the log below)
	provable := make([]chain.Event, 0, len(events))
	provableIdx := make([]int, 0, len(events)) // original index behind provable[j]
	for i, e := range events {
		if e.Height > relayable {
			requeue = append(requeue, i) // not yet provable — wait, no expensive work
			if e.Height > pendingMax {
				pendingMax = e.Height
			}
			continue
		}
		provable = append(provable, e)
		provableIdx = append(provableIdx, i)
	}
	// Surface the wait so the finality/AppHash lag is visible. This runs at the
	// waiting-backoff cadence (the source re-queues these with a growing delay), so
	// it is informative, not spammy — and relayable moving up shows the source's
	// finalized/committed head advancing toward the packet.
	if len(requeue) > 0 {
		log.Printf("[relay %s] waiting: %d packet(s) not yet relayable (highest pending height=%d, source relayable height=%d)",
			m.name, len(requeue), pendingMax, relayable)
	}
	if len(provable) == 0 {
		return requeue // nothing relayable this flush — skip the client update entirely
	}

	// Bail before the expensive client update / proof / submit if we are shutting
	// down; re-queue the provable packets so nothing is lost.
	if ctx.Err() != nil {
		return append(requeue, provableIdx...)
	}

	// Advance the destination client to cover the relayable height. updateClientTo
	// skips the (expensive) build when the client already covers it, so a retry
	// where nothing new is provable costs no proof. On failure the provable packets
	// are transient — re-queue them alongside the pending ones.
	if err := m.updateClientTo(ctx, relayable); err != nil {
		log.Printf("[relay %s] update client to height %d: %v", m.name, relayable, err)
		return append(requeue, provableIdx...)
	}

	// Prove every provable packet against the height the destination client trusts.
	m.mu.Lock()
	proofHeight := m.lastHeight
	m.mu.Unlock()

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
			log.Printf("[relay %s] packet at height %d not yet covered by the destination client (trusts %d); re-queued",
				m.name, e.Height, proofHeight)
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
				log.Printf("[relay %s] DROP packet (type=%d height=%d): %v", m.name, e.Type, e.Height, err)
				continue
			}
			log.Printf("[relay %s] proof for packet (type=%d height=%d): %v", m.name, e.Type, e.Height, err)
			requeue = append(requeue, i)
			continue
		}
		packets = append(packets, chain.RelayPacket{
			Type: e.Type, Packet: e.Raw, Proof: proof, Height: proofHeight, AckBytes: e.AckBytes,
		})
		relayedIdx = append(relayedIdx, i)
	}

	if len(packets) == 0 {
		return requeue
	}

	// Last cancellation check before the tx: on shutdown, re-queue the proven
	// packets rather than broadcasting during teardown.
	if ctx.Err() != nil {
		return append(requeue, relayedIdx...)
	}

	// Submit the whole batch as one multicall.
	err = m.dst.RelayPackets(ctx, packets)
	if err == nil {
		m.untrackDelivered(packets)
		return requeue
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
		log.Printf("[relay %s] DROP packet at height %d (permanent): %v", m.name, proofHeight, err)
		return requeue
	}
	log.Printf("[relay %s] permanent batch failure at height %d (%v); isolating %d packet(s) individually", m.name, proofHeight, err, len(packets))
	return append(requeue, m.relayIsolated(ctx, packets, relayedIdx)...)
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
				log.Printf("[relay %s] DROP packet (type=%d height=%d, permanent): %v", m.name, p.Type, p.Height, err)
				continue
			}
			log.Printf("[relay %s] isolated relay (type=%d height=%d): %v", m.name, p.Type, p.Height, err)
			requeue = append(requeue, relayedIdx[j])
			continue
		}
		m.untrackDelivered(single)
	}
	return requeue
}

// untrackDelivered removes each delivered SendPacket from the pending tracker:
// once received on the destination it can no longer time out, so leaving it in
// the tracker only bloats it and makes the timeout scanner keep querying a
// receipt it will always find. Mirrors the legacy handleCosmos
// PendingTracker.Remove-on-recv. Only paths that supply an untracker
// (cosmos->eth) do this; eth->cosmos removes via its source's terminal
// EthAck/EthTimeout instead.
func (m *Module) untrackDelivered(packets []chain.RelayPacket) {
	if m.untrack == nil {
		return
	}
	for _, p := range packets {
		if p.Type == chain.SendPacket {
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
	return nil
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
			// Zero means "no expiry" (e.g. a permissioned client); skip.
			if expiresAt.IsZero() || time.Until(expiresAt) > refreshMargin {
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
