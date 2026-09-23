package memstore

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"relayer/prover"
	"relayer/prover/witness"
)

func newTestStore(t *testing.T, ttl time.Duration, capEntries int) (*Store, *time.Time) {
	t.Helper()
	s := New(ttl, capEntries)
	t.Cleanup(s.Close)
	now := time.Unix(1000, 0)
	var mu sync.Mutex
	s.now = func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		return now
	}
	// The returned pointer lets the test advance the clock; guarded writes go
	// through advance below in single-goroutine tests, which is safe here.
	return s, &now
}

func TestPutPending_Dedup(t *testing.T) {
	s, _ := newTestStore(t, time.Minute, 8)
	ctx := context.Background()
	id := prover.ProofID{1}

	created, err := s.PutPending(ctx, id)
	if err != nil || !created {
		t.Fatalf("first put: created=%v err=%v", created, err)
	}
	created, err = s.PutPending(ctx, id)
	if err != nil || created {
		t.Fatalf("second put must dedup: created=%v err=%v", created, err)
	}
}

func TestPutPending_ResetsFailedForExplicitRetry(t *testing.T) {
	s, _ := newTestStore(t, time.Minute, 8)
	ctx := context.Background()
	id := prover.ProofID{1}
	if created, err := s.PutPending(ctx, id); err != nil || !created {
		t.Fatalf("first put: created=%v err=%v", created, err)
	}
	token, ok, err := s.Claim(ctx, id, "w1", time.Minute)
	if err != nil || !ok {
		t.Fatalf("claim: ok=%v err=%v", ok, err)
	}
	if err := s.SetProving(ctx, id, token); err != nil {
		t.Fatal(err)
	}
	if err := s.SetFailed(ctx, id, token, prover.FailureTransient, "timeout"); err != nil {
		t.Fatal(err)
	}
	created, err := s.PutPending(ctx, id)
	if err != nil || !created {
		t.Fatalf("retry put: created=%v err=%v", created, err)
	}
	job, err := s.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if job.Status != prover.StatusQueued || job.Class != "" || job.Cause != "" || job.Result != nil {
		t.Fatalf("failed entry was not reset cleanly: %+v", job)
	}
}

func TestPutPending_CapRejects(t *testing.T) {
	s, _ := newTestStore(t, time.Minute, 2)
	ctx := context.Background()
	for i := byte(0); i < 2; i++ {
		if _, err := s.PutPending(ctx, prover.ProofID{i}); err != nil {
			t.Fatalf("put %d: %v", i, err)
		}
	}
	if _, err := s.PutPending(ctx, prover.ProofID{9}); !errors.Is(err, prover.ErrStoreUnavailable) {
		t.Fatalf("over-cap put must return ErrStoreUnavailable, got %v", err)
	}
}

func TestClaim_LifecycleAndLease(t *testing.T) {
	s, now := newTestStore(t, time.Hour, 8)
	ctx := context.Background()
	id := prover.ProofID{1}
	if _, err := s.PutPending(ctx, id); err != nil {
		t.Fatal(err)
	}

	token1, ok, err := s.Claim(ctx, id, "w1", time.Minute)
	if err != nil || !ok {
		t.Fatalf("first claim: ok=%v err=%v", ok, err)
	}
	// A live lease blocks a second worker.
	_, ok, err = s.Claim(ctx, id, "w2", time.Minute)
	if err != nil || ok {
		t.Fatalf("second claim under live lease must fail: ok=%v err=%v", ok, err)
	}
	// Lease expiry makes a dead worker's job reclaimable.
	*now = now.Add(2 * time.Minute)
	token2, ok, err := s.Claim(ctx, id, "w2", time.Minute)
	if err != nil || !ok {
		t.Fatalf("claim after lease expiry: ok=%v err=%v", ok, err)
	}
	// A finished job is never claimable, expired lease or not.
	if token1 == token2 {
		t.Fatal("reclaim reused claim token")
	}
	if err := s.SetProving(ctx, id, token2); err != nil {
		t.Fatal(err)
	}
	if err := s.SetDone(ctx, id, token2, prover.ProofResult{Bucket: 4}); err != nil {
		t.Fatal(err)
	}
	*now = now.Add(2 * time.Minute)
	_, ok, err = s.Claim(ctx, id, "w3", time.Minute)
	if err != nil || ok {
		t.Fatalf("done job must not be claimable: ok=%v err=%v", ok, err)
	}

	if _, _, err := s.Claim(ctx, prover.ProofID{9}, "w1", time.Minute); !errors.Is(err, prover.ErrNotFound) {
		t.Fatalf("claiming unknown id: %v", err)
	}
}

func TestGet_StatesAndResultIsolation(t *testing.T) {
	s, _ := newTestStore(t, time.Hour, 8)
	ctx := context.Background()
	id := prover.ProofID{1}

	if _, err := s.Get(ctx, id); !errors.Is(err, prover.ErrNotFound) {
		t.Fatalf("unknown id must be ErrNotFound, got %v", err)
	}
	if _, err := s.PutPending(ctx, id); err != nil {
		t.Fatal(err)
	}
	token, ok, err := s.Claim(ctx, id, "w1", time.Minute)
	if err != nil || !ok {
		t.Fatalf("claim: ok=%v err=%v", ok, err)
	}
	if err := s.SetProving(ctx, id, token); err != nil {
		t.Fatal(err)
	}
	if err := s.SetFailed(ctx, id, token, prover.FailureTransient, "boom"); err != nil {
		t.Fatal(err)
	}
	job, err := s.Get(ctx, id)
	if err != nil || job.Status != prover.StatusFailed || job.Class != prover.FailureTransient || job.Cause != "boom" {
		t.Fatalf("failed job readback wrong: %+v err=%v", job, err)
	}

	done := prover.ProofResult{
		Bucket:     4,
		PaddedSigs: []witness.ValidatorSignature{{PublicKey: []byte{1, 2, 3}}},
	}
	done.Proof[0] = big.NewInt(7)
	if created, err := s.PutPending(ctx, id); err != nil || !created {
		t.Fatalf("retry: created=%v err=%v", created, err)
	}
	token, ok, err = s.Claim(ctx, id, "w2", time.Minute)
	if err != nil || !ok {
		t.Fatalf("reclaim: ok=%v err=%v", ok, err)
	}
	if err := s.SetProving(ctx, id, token); err != nil {
		t.Fatal(err)
	}
	if err := s.SetDone(ctx, id, token, done); err != nil {
		t.Fatal(err)
	}
	job, _ = s.Get(ctx, id)
	// Caller mutations — scalar, slice content, big.Int — must not reach the
	// store: results are shared state read by concurrent pollers.
	job.Result.Bucket = 99
	job.Result.PaddedSigs[0].PublicKey[0] = 0xFF
	job.Result.Proof[0].SetInt64(1234)
	again, _ := s.Get(ctx, id)
	if again.Result.Bucket != 4 || again.Result.PaddedSigs[0].PublicKey[0] != 1 || again.Result.Proof[0].Int64() != 7 {
		t.Fatalf("Get must return a deep copy of the result, got %+v", again.Result)
	}
}

func TestSweep_ExpiresOldEntries(t *testing.T) {
	s, now := newTestStore(t, time.Minute, 8)
	ctx := context.Background()
	old := prover.ProofID{1}
	fresh := prover.ProofID{2}
	if _, err := s.PutPending(ctx, old); err != nil {
		t.Fatal(err)
	}
	*now = now.Add(2 * time.Minute)
	if _, err := s.PutPending(ctx, fresh); err != nil {
		t.Fatal(err)
	}
	s.sweepOnce()
	if _, err := s.Get(ctx, old); !errors.Is(err, prover.ErrNotFound) {
		t.Fatalf("stale entry must be swept, got %v", err)
	}
	if _, err := s.Get(ctx, fresh); err != nil {
		t.Fatalf("fresh entry must survive: %v", err)
	}
	// A swept id can be resubmitted — the stuck-queued reconciliation path.
	if created, err := s.PutPending(ctx, old); err != nil || !created {
		t.Fatalf("resubmit after sweep: created=%v err=%v", created, err)
	}
}

// A proof can outlive the TTL while its claim remains live. TTL applies only
// after it reaches a terminal state (or to a stuck queued entry).
func TestSweep_KeepsActiveProofPastTTL(t *testing.T) {
	s, now := newTestStore(t, time.Minute, 8)
	ctx := context.Background()
	id := prover.ProofID{3}
	if _, err := s.PutPending(ctx, id); err != nil {
		t.Fatal(err)
	}
	token, ok, err := s.Claim(ctx, id, "w1", time.Minute)
	if err != nil || !ok {
		t.Fatalf("claim: ok=%v err=%v", ok, err)
	}
	if err := s.SetProving(ctx, id, token); err != nil {
		t.Fatal(err)
	}
	*now = now.Add(2 * time.Minute)
	s.sweepOnce()
	if job, err := s.Get(ctx, id); err != nil || job.Status != prover.StatusProving {
		t.Fatalf("active proof swept: job=%+v err=%v", job, err)
	}
	if err := s.SetDone(ctx, id, token, prover.ProofResult{Bucket: 4}); err != nil {
		t.Fatal(err)
	}
	s.sweepOnce()
	if _, err := s.Get(ctx, id); err != nil {
		t.Fatalf("fresh result swept: %v", err)
	}
	*now = now.Add(2 * time.Minute)
	s.sweepOnce()
	if _, err := s.Get(ctx, id); !errors.Is(err, prover.ErrNotFound) {
		t.Fatalf("old terminal job remains: %v", err)
	}
}

func TestClaim_ExpiredLeaseWithoutReclaimCanFinish(t *testing.T) {
	s, now := newTestStore(t, time.Minute, 8)
	ctx := context.Background()
	id := prover.ProofID{5}
	if _, err := s.PutPending(ctx, id); err != nil {
		t.Fatal(err)
	}
	token, ok, err := s.Claim(ctx, id, "w1", time.Minute)
	if err != nil || !ok {
		t.Fatalf("claim: ok=%v err=%v", ok, err)
	}
	*now = now.Add(2 * time.Minute)
	// A long prove may pass the lease deadline before the worker reports it.
	// The lease allows another claim; time alone does not revoke this one.
	if err := s.SetProving(ctx, id, token); err != nil {
		t.Fatalf("late SetProving: %v", err)
	}
	s.sweepOnce()
	if err := s.SetDone(ctx, id, token, prover.ProofResult{Bucket: 4}); err != nil {
		t.Fatalf("late SetDone: %v", err)
	}
	job, err := s.Get(ctx, id)
	if err != nil || job.Status != prover.StatusDone || job.Result.Bucket != 4 {
		t.Fatalf("result lost after lease expiry: job=%+v err=%v", job, err)
	}
	if _, ok, err := s.Claim(ctx, id, "w2", time.Minute); err != nil || ok {
		t.Fatalf("completed job reclaimed: ok=%v err=%v", ok, err)
	}
}

func TestClaim_FencesExpiredAndReclaimedWorkers(t *testing.T) {
	s, now := newTestStore(t, time.Minute, 8)
	ctx := context.Background()
	id := prover.ProofID{4}
	if _, err := s.PutPending(ctx, id); err != nil {
		t.Fatal(err)
	}
	first, ok, err := s.Claim(ctx, id, "same-worker", time.Minute)
	if err != nil || !ok {
		t.Fatalf("first claim: ok=%v err=%v", ok, err)
	}
	if err := s.SetProving(ctx, id, first); err != nil {
		t.Fatal(err)
	}
	*now = now.Add(time.Minute)
	s.sweepOnce()
	second, ok, err := s.Claim(ctx, id, "same-worker", time.Minute)
	if err != nil || !ok || second == first {
		t.Fatalf("reclaim: token=%d ok=%v err=%v", second, ok, err)
	}
	if err := s.SetProving(ctx, id, first); !errors.Is(err, prover.ErrStaleClaim) {
		t.Fatalf("stale SetProving: %v", err)
	}
	if err := s.SetProving(ctx, id, second); err != nil {
		t.Fatal(err)
	}
	if err := s.SetDone(ctx, id, first, prover.ProofResult{Bucket: 1}); !errors.Is(err, prover.ErrStaleClaim) {
		t.Fatalf("stale SetDone: %v", err)
	}
	if err := s.SetFailed(ctx, id, first, prover.FailureInternal, "late"); !errors.Is(err, prover.ErrStaleClaim) {
		t.Fatalf("stale SetFailed: %v", err)
	}
	if err := s.SetDone(ctx, id, second, prover.ProofResult{Bucket: 4}); err != nil {
		t.Fatal(err)
	}
	if err := s.SetFailed(ctx, id, second, prover.FailureInternal, "late"); !errors.Is(err, prover.ErrStaleClaim) {
		t.Fatalf("done overwritten: %v", err)
	}
	job, err := s.Get(ctx, id)
	if err != nil || job.Status != prover.StatusDone || job.Result.Bucket != 4 {
		t.Fatalf("winner result lost: job=%+v err=%v", job, err)
	}

	// Deleting and recreating the same deterministic id must not recycle the token.
	if err := s.Delete(ctx, id); err != nil {
		t.Fatal(err)
	}
	if _, err := s.PutPending(ctx, id); err != nil {
		t.Fatal(err)
	}
	third, ok, err := s.Claim(ctx, id, "same-worker", time.Minute)
	if err != nil || !ok || third == second {
		t.Fatalf("new job token: %d ok=%v err=%v", third, ok, err)
	}
	if err := s.SetProving(ctx, id, second); !errors.Is(err, prover.ErrStaleClaim) {
		t.Fatalf("prior job token accepted: %v", err)
	}
}
