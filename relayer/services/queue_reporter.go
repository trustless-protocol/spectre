package services

import (
	"context"
	"log"
	"sync"
	"time"
)

// queueReportInterval is a variable so the reporter can be exercised without
// waiting a minute. Production keeps the intentionally low-noise cadence.
var queueReportInterval = time.Minute

type clientUpdateTimes struct {
	mu          sync.RWMutex
	cosmosOnEVM time.Time
	evmOnCosmos time.Time
}

func (s *Services) ObserveCosmosOnEVMUpdate(trustedAt time.Time) {
	s.clientUpdates.mu.Lock()
	defer s.clientUpdates.mu.Unlock()
	if trustedAt.After(s.clientUpdates.cosmosOnEVM) {
		s.clientUpdates.cosmosOnEVM = trustedAt
	}
}

func (s *Services) ObserveEVMOnCosmosUpdate(trustedAt time.Time) {
	s.clientUpdates.mu.Lock()
	defer s.clientUpdates.mu.Unlock()
	if trustedAt.After(s.clientUpdates.evmOnCosmos) {
		s.clientUpdates.evmOnCosmos = trustedAt
	}
}

// clientUpdateAges reports staleness using the source-chain timestamps trusted
// by each light client, not process uptime.
func (s *Services) clientUpdateAges(now time.Time) (cosmosOnEVM, evmOnCosmos time.Duration, ok bool) {
	s.clientUpdates.mu.RLock()
	defer s.clientUpdates.mu.RUnlock()
	if s.clientUpdates.cosmosOnEVM.IsZero() && s.clientUpdates.evmOnCosmos.IsZero() {
		return 0, 0, false
	}
	if !s.clientUpdates.cosmosOnEVM.IsZero() {
		cosmosOnEVM = now.Sub(s.clientUpdates.cosmosOnEVM)
	}
	if !s.clientUpdates.evmOnCosmos.IsZero() {
		evmOnCosmos = now.Sub(s.clientUpdates.evmOnCosmos)
	}
	return cosmosOnEVM, evmOnCosmos, true
}

type queueState struct {
	cosmosQueued, ethQueued                                   int
	cosmosPending, ethPending, l2Pending                      int
	cosmosTimeoutDead, ethTimeoutDead, l2TimeoutDead          int
	cosmosStuck, ethStuck, l2Stuck                            int
	cosmosWorstDeferrals, ethWorstDeferrals, l2WorstDeferrals int
	cosmosTimeoutDeadAge, ethTimeoutDeadAge, l2TimeoutDeadAge time.Duration
}

func (s *Services) queueStateAt(now time.Time) queueState {
	cosmosQueued, ethQueued := s.BatchBuilder.QueueDepths()
	cosmosStuck, cosmosWorst := s.BatchBuilder.PendingTracker.StuckTimeouts()
	ethStuck, ethWorst := s.BatchBuilder.EthPendingTracker.StuckTimeouts()
	l2Stuck, l2Worst := s.BatchBuilder.L2PendingTracker.StuckTimeouts()
	return queueState{
		cosmosQueued:         cosmosQueued,
		ethQueued:            ethQueued,
		cosmosPending:        s.BatchBuilder.PendingTracker.Len(),
		ethPending:           s.BatchBuilder.EthPendingTracker.Len(),
		l2Pending:            s.BatchBuilder.L2PendingTracker.Len(),
		cosmosTimeoutDead:    s.BatchBuilder.PendingTracker.DeadLetteredTimeouts(),
		ethTimeoutDead:       s.BatchBuilder.EthPendingTracker.DeadLetteredTimeouts(),
		l2TimeoutDead:        s.BatchBuilder.L2PendingTracker.DeadLetteredTimeouts(),
		cosmosStuck:          cosmosStuck,
		ethStuck:             ethStuck,
		l2Stuck:              l2Stuck,
		cosmosWorstDeferrals: cosmosWorst,
		ethWorstDeferrals:    ethWorst,
		l2WorstDeferrals:     l2Worst,
		cosmosTimeoutDeadAge: s.BatchBuilder.PendingTracker.OldestDeadLetteredTimeout(now),
		ethTimeoutDeadAge:    s.BatchBuilder.EthPendingTracker.OldestDeadLetteredTimeout(now),
		l2TimeoutDeadAge:     s.BatchBuilder.L2PendingTracker.OldestDeadLetteredTimeout(now),
	}
}

func (s *Services) reportQueueState(now time.Time) {
	state := s.queueStateAt(now)
	log.Printf("[QueueState] queued(cosmos=%d eth=%d) pending(cosmos=%d eth=%d l2=%d) timeout-dead-letter(cosmos=%d eth=%d l2=%d)",
		state.cosmosQueued, state.ethQueued,
		state.cosmosPending, state.ethPending, state.l2Pending,
		state.cosmosTimeoutDead, state.ethTimeoutDead, state.l2TimeoutDead)
	if state.cosmosTimeoutDead > 0 || state.ethTimeoutDead > 0 || state.l2TimeoutDead > 0 {
		log.Printf("[QueueState][ATTENTION] timeout submission gave up for cosmos=%d eth=%d l2=%d packet(s) (oldest: cosmos=%s eth=%s l2=%s); their escrowed funds require operator action.",
			state.cosmosTimeoutDead, state.ethTimeoutDead, state.l2TimeoutDead,
			state.cosmosTimeoutDeadAge.Truncate(time.Second), state.ethTimeoutDeadAge.Truncate(time.Second), state.l2TimeoutDeadAge.Truncate(time.Second))
	}
	if state.cosmosStuck > 0 || state.ethStuck > 0 || state.l2Stuck > 0 {
		log.Printf("[QueueState][ATTENTION] timeout deferrals are stuck for cosmos=%d eth=%d l2=%d packet(s) (most deferrals: cosmos=%d eth=%d l2=%d); check counterparty RPC and light-client health.",
			state.cosmosStuck, state.ethStuck, state.l2Stuck,
			state.cosmosWorstDeferrals, state.ethWorstDeferrals, state.l2WorstDeferrals)
	}
	if cosmosOnEVM, evmOnCosmos, ok := s.clientUpdateAges(now); ok {
		log.Printf("[QueueState] last client update: cosmos-on-eth=%s ago, eth-on-cosmos=%s ago",
			cosmosOnEVM.Truncate(time.Second), evmOnCosmos.Truncate(time.Second))
	}
}

// RunQueueReporter emits an immediate structured queue-state line and then
// reports periodically until the engine context is cancelled. The adapter owns
// the goroutine, so each Services instance has exactly one reporter.
func (s *Services) RunQueueReporter(ctx context.Context) {
	s.reportQueueState(time.Now())
	ticker := time.NewTicker(queueReportInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			s.reportQueueState(now)
		}
	}
}
