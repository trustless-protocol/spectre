// Package memstore is the in-memory ProofStore adapter: a mutex-guarded map
// with an expiring lease per claim and a TTL sweeper. It is a cache, not a
// ledger — a restart loses it, and a re-submit reproduces every id.
package memstore

import (
	"context"
	"math/big"
	"sync"
	"time"

	"relayer/prover"
	"relayer/prover/witness"
)

type entry struct {
	job        prover.ProofJob
	worker     string
	token      prover.ClaimToken
	leaseUntil time.Time
	updatedAt  time.Time
}

type Store struct {
	mu        sync.Mutex
	nextToken prover.ClaimToken
	entries   map[prover.ProofID]*entry
	ttl       time.Duration
	cap       int
	stop      chan struct{}
	stopped   sync.Once
	now       func() time.Time // injectable for lease/TTL tests
}

// New returns a store whose finished entries expire after ttl (default 15m)
// and which holds at most capEntries entries (default 1024). Close stops the
// TTL sweeper.
func New(ttl time.Duration, capEntries int) *Store {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	if capEntries <= 0 {
		capEntries = 1024
	}
	s := &Store{
		entries: make(map[prover.ProofID]*entry),
		ttl:     ttl,
		cap:     capEntries,
		stop:    make(chan struct{}),
		now:     time.Now,
	}
	go s.sweep()
	return s
}

func (s *Store) Close() { s.stopped.Do(func() { close(s.stop) }) }

func (s *Store) PutPending(_ context.Context, id prover.ProofID) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.entries[id]; ok {
		if e.job.Status != prover.StatusFailed {
			return false, nil
		}
		now := s.now()
		e.job = prover.ProofJob{ID: id, Status: prover.StatusQueued}
		e.worker = ""
		e.leaseUntil = time.Time{}
		e.updatedAt = now
		return true, nil
	}
	if len(s.entries) >= s.cap {
		return false, prover.ErrStoreUnavailable
	}
	s.entries[id] = &entry{
		job:       prover.ProofJob{ID: id, Status: prover.StatusQueued},
		updatedAt: s.now(),
	}
	return true, nil
}

func (s *Store) Claim(_ context.Context, id prover.ProofID, workerID string, lease time.Duration) (prover.ClaimToken, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[id]
	if !ok {
		return 0, false, prover.ErrNotFound
	}
	now := s.now()
	switch e.job.Status {
	case prover.StatusQueued:
		// claimable
	case prover.StatusClaimed, prover.StatusProving:
		if now.Before(e.leaseUntil) {
			return 0, false, nil // a live worker owns it
		}
		// lease expired: the owning worker is dead, reclaim
	default: // done, failed
		return 0, false, nil
	}
	e.job.Status = prover.StatusClaimed
	e.worker = workerID
	s.nextToken++
	e.token = s.nextToken
	e.leaseUntil = now.Add(lease)
	e.updatedAt = now
	return e.token, true, nil
}

func (s *Store) SetProving(_ context.Context, id prover.ProofID, token prover.ClaimToken) error {
	return s.update(id, token, prover.StatusClaimed, func(e *entry) {
		e.job.Status = prover.StatusProving
	})
}

func (s *Store) SetDone(_ context.Context, id prover.ProofID, token prover.ClaimToken, r prover.ProofResult) error {
	return s.update(id, token, prover.StatusProving, func(e *entry) {
		e.job.Status = prover.StatusDone
		e.job.Result = &r
	})
}

func (s *Store) SetFailed(_ context.Context, id prover.ProofID, token prover.ClaimToken, class prover.FailureClass, cause string) error {
	return s.update(id, token, prover.StatusProving, func(e *entry) {
		e.job.Status = prover.StatusFailed
		e.job.Class = class
		e.job.Cause = cause
	})
}

func (s *Store) Delete(_ context.Context, id prover.ProofID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, id)
	return nil
}

func (s *Store) Get(_ context.Context, id prover.ProofID) (prover.ProofJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[id]
	if !ok {
		return prover.ProofJob{}, prover.ErrNotFound
	}
	job := e.job
	if job.Result != nil {
		r := deepCopyResult(*job.Result)
		job.Result = &r
	}
	return job, nil
}

// deepCopyResult clones a result down to its byte slices and big.Ints, so a
// caller mutating what Get returned can never corrupt the store's entry —
// results are shared state read by concurrent pollers.
func deepCopyResult(r prover.ProofResult) prover.ProofResult {
	out := r
	out.PaddedSigs = make([]witness.ValidatorSignature, len(r.PaddedSigs))
	for i, v := range r.PaddedSigs {
		v.Signature = append([]byte(nil), v.Signature...)
		v.PublicKey = append([]byte(nil), v.PublicKey...)
		v.SignedBytes = append([]byte(nil), v.SignedBytes...)
		out.PaddedSigs[i] = v
	}
	for i, v := range r.Proof {
		out.Proof[i] = cloneBig(v)
	}
	for i, v := range r.Commitments {
		out.Commitments[i] = cloneBig(v)
	}
	for i, v := range r.CommitmentPok {
		out.CommitmentPok[i] = cloneBig(v)
	}
	return out
}

func cloneBig(v *big.Int) *big.Int {
	if v == nil {
		return nil
	}
	return new(big.Int).Set(v)
}

// Len reports the entry count (for metrics/tests).
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries)
}

func (s *Store) update(id prover.ProofID, token prover.ClaimToken, from prover.Status, f func(*entry)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[id]
	if !ok {
		return prover.ErrNotFound
	}
	// Expiry only opens the door to a replacement Claim. The current owner
	// may still finish until another claim actually supersedes its token.
	if token == 0 || e.token != token || e.job.Status != from {
		return prover.ErrStaleClaim
	}
	f(e)
	e.updatedAt = s.now()
	return nil
}

// sweep drops terminal jobs and stale queued entries after the TTL. Active
// claims remain until a worker finishes or the lease expires and another
// worker reclaims them; a long proof must not lose its result to the sweeper.
func (s *Store) sweep() {
	interval := max(s.ttl/4, time.Second)
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-s.stop:
			return
		case <-t.C:
			s.sweepOnce()
		}
	}
}

func (s *Store) sweepOnce() {
	cutoff := s.now().Add(-s.ttl)
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, e := range s.entries {
		if e.updatedAt.Before(cutoff) && (e.job.Status == prover.StatusQueued || e.job.Status == prover.StatusDone || e.job.Status == prover.StatusFailed) {
			delete(s.entries, id)
		}
	}
}
