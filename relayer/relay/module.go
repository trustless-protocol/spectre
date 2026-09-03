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
