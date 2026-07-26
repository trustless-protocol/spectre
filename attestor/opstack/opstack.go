package opstack

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"attestor"
)

// gameSource and replica are the two upstream dependencies of the attestation
// loop, defined at the consumer so tests can fake them.
type gameSource interface {
	GameCount(ctx context.Context) (uint64, error)
	GameAtIndex(ctx context.Context, index uint64) (ProposedRoot, error)
}

type replica interface {
	SyncStatus(ctx context.Context) (SyncStatus, error)
	OutputAtBlock(ctx context.Context, l2Block uint64) ([32]byte, error)
}

// OpStackAttestor replays the L2 state transition (via a verify-mode replica
// that derives solely from L1) to attest proposed output roots before anything
// relays against them. It implements attestor.Attestor.
//
// Goroutine ownership: the Run goroutine is the only writer of pending state,
// the ingest cursor, the backoff, and every store mutation; the wake-hint
// goroutine owns nothing shared (it only signals the buffered wake channel);
// external readers use the store's RWMutex-guarded accessors.
type OpStackAttestor struct {
	cfg     Config
	games   gameSource
	replica replica
	store   *AttestedRootStore
	hook    ChallengeHook
	metrics *Metrics
	logger  *zap.SugaredLogger
	backoff attestor.Backoff
	now     func() time.Time
	wake    chan struct{}
	// lastStatus is the most recent successful replica sync status, published
	// for read-only observers (the sidecar gRPC server); nil until the first
	// successful poll.
	lastStatus atomic.Pointer[SyncStatus]
}

// New assembles an attestor from already-validated config and dialed
// dependencies (see cmd/attest.go for the production wiring).
func New(cfg Config, games gameSource, rep replica, store *AttestedRootStore, hook ChallengeHook, metrics *Metrics, logger *zap.SugaredLogger) *OpStackAttestor {
	return &OpStackAttestor{
		cfg:     cfg,
		games:   games,
		replica: rep,
		store:   store,
		hook:    hook,
		metrics: metrics,
		logger:  logger,
		now:     time.Now,
		wake:    make(chan struct{}, 1),
	}
}

func (a *OpStackAttestor) Name() string {
	return "opstack:" + a.cfg.SrcChain
}

// AttestedUpTo reports the highest confirmed (non-provisional) attested root.
func (a *OpStackAttestor) AttestedUpTo() (attestor.Cursor, bool) {
	return a.store.AttestedUpTo(false)
}

// Store exposes the attested-root feed for consumers (the future op_to_cosmos
// module relays only against roots this store affirms).
func (a *OpStackAttestor) Store() *AttestedRootStore {
	return a.store
}

// SrcChain returns the configured source chain identifier.
func (a *OpStackAttestor) SrcChain() string {
	return a.cfg.SrcChain
}

// Head returns the configured gating head.
func (a *OpStackAttestor) Head() Head {
	return a.cfg.AttestationHead
}

// LastSyncStatus returns the most recent successful replica sync status; ok is
// false until the replica has responded once.
func (a *OpStackAttestor) LastSyncStatus() (SyncStatus, bool) {
	if s := a.lastStatus.Load(); s != nil {
		return *s, true
	}
	return SyncStatus{}, false
}

// Run executes the attestation loop until ctx is done: poll ticker plus the
// optional DisputeGameCreated wake hint, with failure backoff.
func (a *OpStackAttestor) Run(ctx context.Context) error {
	if a.store.Fresh() {
		if err := a.bootstrap(ctx); err != nil {
			return fmt.Errorf("%s: bootstrap failed: %w", a.Name(), err)
		}
	}
	if a.cfg.L1WsUrl != "" {
		go wakeHintLoop(ctx, a.logger, a.cfg.L1WsUrl, a.cfg.DisputeGameFactory, a.wake)
	}
	a.logger.Infof("%s: attestation loop started (head=%s, poll=%s, factory=%s)",
		a.Name(), a.cfg.AttestationHead, a.cfg.PollInterval, a.cfg.DisputeGameFactory)

	ticker := time.NewTicker(a.cfg.PollInterval)
	defer ticker.Stop()
	for {
		if a.backoff.Ready(a.now()) {
			if err := a.runOnce(ctx); err != nil {
				if ctx.Err() != nil {
					return nil
				}
				delay := a.backoff.RecordFailure(a.now())
				a.logger.Warnf("%s: attestation pass failed (nothing affirmed, retrying in %s): %v", a.Name(), delay, err)
			} else {
				a.backoff.RecordSuccess()
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		case <-a.wake:
		}
	}
}

// bootstrap sets the initial ingest cursor on a fresh store: the first game
// created inside the configured L1 lookback window.
func (a *OpStackAttestor) bootstrap(ctx context.Context) error {
	count, err := a.games.GameCount(ctx)
	if err != nil {
		return err
	}
	cutoff := uint64(0)
	// L1 slots are 12s; the lookback is a bound on history replay, not a
	// precise window, so wall-clock arithmetic is sufficient.
	if lookbackSecs := a.cfg.BootstrapLookbackBlocks * 12; lookbackSecs > 0 {
		nowSecs := uint64(a.now().Unix())
		if nowSecs > lookbackSecs {
			cutoff = nowSecs - lookbackSecs
		}
	}
	start, err := bootstrapStartIndex(ctx, a.games, count, cutoff)
	if err != nil {
		return err
	}
	a.store.Bootstrap(start)
	if err := a.store.Save(); err != nil {
		return err
	}
	a.logger.Infof("%s: bootstrapped ingest cursor at game index %d of %d (lookback %d L1 blocks)",
		a.Name(), start, count, a.cfg.BootstrapLookbackBlocks)
	return nil
}

// gateHead returns the replica head that gates verdicts and derived
// attestations for the configured attestation level.
func (a *OpStackAttestor) gateHead(status SyncStatus) uint64 {
	switch a.cfg.AttestationHead {
	case HeadUnsafe:
		return status.UnsafeL2
	case HeadSafe:
		return status.SafeL2
	default:
		return status.FinalizedL2
	}
}

// runOnce is one attestation pass: ingest new games, read the replica heads,
// resolve provisional entries the finalized head now covers, decide pending
// games, and self-attest the current head. Any error returns with nothing
// affirmed — no cursor advances past a failure, no game leaves pending
// without an explicit verdict, and the feed is only ever appended on a
// verdict or a derived attestation.
func (a *OpStackAttestor) runOnce(ctx context.Context) error {
	if err := a.ingest(ctx); err != nil {
		return err
	}

	status, err := a.replica.SyncStatus(ctx)
	if err != nil {
		a.countReplicaError("rpc")
		return err
	}
	a.lastStatus.Store(&status)
	a.gauge(func(m *Metrics) {
		m.ReplicaFinalizedHead.WithLabelValues(a.cfg.SrcChain).Set(float64(status.FinalizedL2))
		m.ReplicaSafeHead.WithLabelValues(a.cfg.SrcChain).Set(float64(status.SafeL2))
		m.ReplicaUnsafeHead.WithLabelValues(a.cfg.SrcChain).Set(float64(status.UnsafeL2))
	})

	if err := a.resolveRechecks(ctx, status.FinalizedL2); err != nil {
		return err
	}
	if err := a.reverifyDerived(ctx, status.FinalizedL2); err != nil {
		return err
	}
	if err := a.decidePending(ctx, status); err != nil {
		return err
	}
	if err := a.attestDerived(ctx, status); err != nil {
		return err
	}
	a.updateLagMetrics(a.gateHead(status))
	return nil
}

// ingest advances the game-index cursor through the factory list. The cursor
// only ever moves past an index whose fetch succeeded.
func (a *OpStackAttestor) ingest(ctx context.Context) error {
	count, err := a.games.GameCount(ctx)
	if err != nil {
		return err
	}
	start := a.store.NextGameIndex()
	if start >= count {
		return nil
	}
	for idx := start; idx < count; idx++ {
		game, err := a.games.GameAtIndex(ctx, idx)
		if err != nil {
			// Persist what was ingested so far; the cursor stays at idx.
			if saveErr := a.store.Save(); saveErr != nil {
				a.logger.Warnf("%s: failed to persist state after partial ingest: %v", a.Name(), saveErr)
			}
			return fmt.Errorf("ingest stopped at game index %d: %w", idx, err)
		}
		a.count(func(m *Metrics) { m.GamesIngested.WithLabelValues(a.cfg.SrcChain).Inc() })
		if game.GameType != a.cfg.RespectedGameType {
			a.logger.Debugf("%s: skipping game %d: type %d != respected type %d", a.Name(), idx, game.GameType, a.cfg.RespectedGameType)
			a.store.AdvanceIngest(idx+1, nil)
			continue
		}
		a.logger.Infof("%s: ingested proposal: game %d, l2 block %d, root %x", a.Name(), idx, game.L2BlockNumber, game.RootClaim)
		g := game
		a.store.AdvanceIngest(idx+1, &g)
	}
	return a.store.Save()
}

// decidePending issues verdicts for pending games covered by the gating head.
// A game leaves pending only via an explicit verdict; replica errors keep it
// pending forever (there is no skip-after-N-retries for verification).
func (a *OpStackAttestor) decidePending(ctx context.Context, status SyncStatus) error {
	gate := a.gateHead(status)
	for _, game := range a.store.Pending() {
		if game.L2BlockNumber > gate {
			continue
		}
		local, err := a.replica.OutputAtBlock(ctx, game.L2BlockNumber)
		if err != nil {
			a.countReplicaError("rpc")
			return fmt.Errorf("deciding game %d (l2 block %d): %w", game.GameIndex, game.L2BlockNumber, err)
		}
		// Provisional iff the verdict covers a block the finalized head has
		// not yet re-derived from finalized L1 data.
		provisional := game.L2BlockNumber > status.FinalizedL2
		a.recordVerdict(ctx, game, local, provisional)
		if err := a.store.Save(); err != nil {
			return fmt.Errorf("persisting verdict for game %d: %w", game.GameIndex, err)
		}
	}
	return nil
}

// attestDerived self-attests the replica's output root at the gating head —
// the proposal-independent feed: a consumer waits only for the configured
// trust level plus this verification, never for a proposer.
func (a *OpStackAttestor) attestDerived(ctx context.Context, status SyncStatus) error {
	if a.cfg.DisableDerivedRoots {
		return nil
	}
	gate := a.gateHead(status)
	if gate == 0 {
		return nil
	}
	if last, ok := a.store.HighestDerivedBlock(); ok && gate < last+a.cfg.DerivedGapBlocks {
		return nil
	}
	root, err := a.replica.OutputAtBlock(ctx, gate)
	if err != nil {
		a.countReplicaError("rpc")
		return fmt.Errorf("deriving output root at head block %d: %w", gate, err)
	}
	provisional := gate > status.FinalizedL2
	a.store.AppendDerived(gate, root, a.now(), provisional)
	a.count(func(m *Metrics) { m.RootsDerived.WithLabelValues(a.cfg.SrcChain).Inc() })
	if provisional {
		a.count(func(m *Metrics) { m.RootsProvisional.WithLabelValues(a.cfg.SrcChain).Inc() })
	}
	a.logger.Infof("%s: derived root attested at l2 block %d (root %x, head=%s, provisional=%t)",
		a.Name(), gate, root, a.cfg.AttestationHead, provisional)
	if pruned := a.store.PruneDerived(int(a.cfg.MaxDerivedRoots)); pruned > 0 {
		a.logger.Debugf("%s: pruned %d confirmed derived roots beyond cap %d", a.Name(), pruned, a.cfg.MaxDerivedRoots)
	}
	return a.store.Save()
}

// reverifyDerived re-computes provisional derived roots once the finalized
// head covers them. The finalized result is authoritative: agreement confirms
// the entry, disagreement corrects it and raises the divergence alarm.
func (a *OpStackAttestor) reverifyDerived(ctx context.Context, finalized uint64) error {
	for _, entry := range a.store.DerivedProvisionalAtOrBelow(finalized) {
		root, err := a.replica.OutputAtBlock(ctx, entry.L2BlockNumber)
		if err != nil {
			a.countReplicaError("rpc")
			return fmt.Errorf("reverifying derived root at l2 block %d: %w", entry.L2BlockNumber, err)
		}
		if root == entry.Root {
			a.store.ConfirmDerived(entry.L2BlockNumber)
			a.logger.Infof("%s: finalized head confirmed derived root at l2 block %d", a.Name(), entry.L2BlockNumber)
		} else {
			a.store.CorrectDerived(entry.L2BlockNumber, root, a.now())
			a.count(func(m *Metrics) { m.HeadDivergence.WithLabelValues(a.cfg.SrcChain).Inc() })
			a.logger.Errorw("[opstack attestor] HEAD DIVERGENCE — finalized derivation contradicts provisional derived root",
				"chain", a.cfg.SrcChain,
				"l2_block", entry.L2BlockNumber,
				"provisional_root", fmt.Sprintf("%x", entry.Root),
				"finalized_root", fmt.Sprintf("%x", root),
			)
		}
		if err := a.store.Save(); err != nil {
			return fmt.Errorf("persisting derived reverification at l2 block %d: %w", entry.L2BlockNumber, err)
		}
	}
	return nil
}

func (a *OpStackAttestor) recordVerdict(ctx context.Context, game ProposedRoot, local [32]byte, provisional bool) {
	now := a.now()
	if local == game.RootClaim {
		a.store.RecordMatch(game, now, provisional)
		a.count(func(m *Metrics) { m.RootsMatched.WithLabelValues(a.cfg.SrcChain).Inc() })
		a.logger.Infof("%s: attested game %d (l2 block %d, root %x, provisional=%t)",
			a.Name(), game.GameIndex, game.L2BlockNumber, game.RootClaim, provisional)
	} else {
		a.store.RecordMismatch(game, local, now, provisional)
		a.count(func(m *Metrics) { m.RootsMismatched.WithLabelValues(a.cfg.SrcChain).Inc() })
		if err := a.hook.OnMismatch(ctx, game, local); err != nil {
			a.logger.Errorf("%s: challenge hook failed for game %d: %v", a.Name(), game.GameIndex, err)
		}
	}
	if provisional {
		a.count(func(m *Metrics) { m.RootsProvisional.WithLabelValues(a.cfg.SrcChain).Inc() })
	}
}

// resolveRechecks re-runs provisional game verdicts once the finalized head
// covers their block. The finalized result is always authoritative; a
// contradiction is a head-divergence alarm (sequencer divergence for verdicts
// made at unsafe, an L1 reorg for verdicts made at safe).
func (a *OpStackAttestor) resolveRechecks(ctx context.Context, finalized uint64) error {
	for _, entry := range a.store.RecheckList() {
		if entry.Game.L2BlockNumber > finalized {
			continue
		}
		local, err := a.replica.OutputAtBlock(ctx, entry.Game.L2BlockNumber)
		if err != nil {
			a.countReplicaError("rpc")
			return fmt.Errorf("rechecking game %d (l2 block %d): %w", entry.Game.GameIndex, entry.Game.L2BlockNumber, err)
		}
		finalMatch := local == entry.Game.RootClaim
		if finalMatch == entry.ProvisionalMatch {
			a.store.ConfirmProvisional(entry.Game.GameIndex)
			a.logger.Infof("%s: finalized head confirmed provisional verdict for game %d (match=%t)",
				a.Name(), entry.Game.GameIndex, finalMatch)
		} else {
			a.store.CorrectProvisional(entry, local, a.now())
			a.count(func(m *Metrics) { m.HeadDivergence.WithLabelValues(a.cfg.SrcChain).Inc() })
			a.logger.Errorw("[opstack attestor] HEAD DIVERGENCE — finalized recheck contradicts provisional verdict",
				"chain", a.cfg.SrcChain,
				"game_index", entry.Game.GameIndex,
				"l2_block", entry.Game.L2BlockNumber,
				"provisional_match", entry.ProvisionalMatch,
				"finalized_match", finalMatch,
			)
			if !finalMatch {
				if err := a.hook.OnMismatch(ctx, entry.Game, local); err != nil {
					a.logger.Errorf("%s: challenge hook failed for rechecked game %d: %v", a.Name(), entry.Game.GameIndex, err)
				}
			}
		}
		if err := a.store.Save(); err != nil {
			return fmt.Errorf("persisting recheck for game %d: %w", entry.Game.GameIndex, err)
		}
	}
	return nil
}

func (a *OpStackAttestor) updateLagMetrics(gate uint64) {
	pending := a.store.Pending()
	var lag uint64
	for _, g := range pending {
		if g.L2BlockNumber > gate {
			if d := g.L2BlockNumber - gate; d > lag {
				lag = d
			}
		}
	}
	a.gauge(func(m *Metrics) {
		m.AttestationLag.WithLabelValues(a.cfg.SrcChain).Set(float64(lag))
		m.PendingGames.WithLabelValues(a.cfg.SrcChain).Set(float64(len(pending)))
	})
}

func (a *OpStackAttestor) countReplicaError(kind string) {
	a.count(func(m *Metrics) { m.ReplicaErrors.WithLabelValues(a.cfg.SrcChain, kind).Inc() })
}

func (a *OpStackAttestor) count(f func(*Metrics)) {
	if a.metrics != nil {
		f(a.metrics)
	}
}

func (a *OpStackAttestor) gauge(f func(*Metrics)) {
	if a.metrics != nil {
		f(a.metrics)
	}
}
