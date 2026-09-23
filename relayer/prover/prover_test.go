package prover

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ---- fakes -----------------------------------------------------------------

type fakeEngine struct {
	proves atomic.Int64
	// prove is called under the worker; default returns instantly.
	prove func(ctx context.Context, req ProofRequest) (ProofResult, error)
}

func (f *fakeEngine) Prove(ctx context.Context, req ProofRequest) (ProofResult, error) {
	f.proves.Add(1)
	if f.prove != nil {
		return f.prove(ctx, req)
	}
	return ProofResult{Bucket: req.Bucket, PaddedSigs: req.Sigs}, nil
}
func (f *fakeEngine) Buckets() []int        { return []int{4} }
func (f *fakeEngine) VKDigest(int) [32]byte { return [32]byte{0x42} }

type fakeQueue struct {
	ch   chan ProofRequest
	acks atomic.Int64
}

func newFakeQueue(size int) *fakeQueue { return &fakeQueue{ch: make(chan ProofRequest, size)} }

func (q *fakeQueue) Enqueue(_ context.Context, job ProofRequest) error {
	select {
	case q.ch <- job:
		return nil
	default:
		return ErrQueueFull
	}
}

func (q *fakeQueue) Dequeue(ctx context.Context) (ProofRequest, func(), error) {
	select {
	case <-ctx.Done():
		return ProofRequest{}, nil, ctx.Err()
	case job := <-q.ch:
		return job, func() { q.acks.Add(1) }, nil
	}
}

type fakeStore struct {
	mu        sync.Mutex
	jobs      map[ProofID]*ProofJob
	failPut   error // injected PutPending failure
	failDone  int   // fail this many SetDone calls before succeeding
	closes    atomic.Int64
	nextToken ClaimToken
	tokens    map[ProofID]ClaimToken
}

func newFakeStore() *fakeStore {
	return &fakeStore{jobs: make(map[ProofID]*ProofJob), tokens: make(map[ProofID]ClaimToken)}
}

func (s *fakeStore) PutPending(_ context.Context, id ProofID) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.failPut != nil {
		return false, s.failPut
	}
	if job, ok := s.jobs[id]; ok {
		if job.Status != StatusFailed {
			return false, nil
		}
		s.jobs[id] = &ProofJob{ID: id, Status: StatusQueued}
		return true, nil
	}
	s.jobs[id] = &ProofJob{ID: id, Status: StatusQueued}
	return true, nil
}

func (s *fakeStore) Claim(_ context.Context, id ProofID, _ string, _ time.Duration) (ClaimToken, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[id]
	if !ok {
		return 0, false, ErrNotFound
	}
	if j.Status != StatusQueued {
		return 0, false, nil
	}
	j.Status = StatusClaimed
	s.nextToken++
	s.tokens[id] = s.nextToken
	return s.nextToken, true, nil
}

func (s *fakeStore) set(id ProofID, token ClaimToken, from Status, f func(*ProofJob)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[id]
	if !ok {
		return ErrNotFound
	}
	if token == 0 || token != s.tokens[id] || j.Status != from {
		return ErrStaleClaim
	}
	f(j)
	return nil
}

func (s *fakeStore) SetProving(_ context.Context, id ProofID, token ClaimToken) error {
	return s.set(id, token, StatusClaimed, func(j *ProofJob) { j.Status = StatusProving })
}

func (s *fakeStore) SetDone(_ context.Context, id ProofID, token ClaimToken, r ProofResult) error {
	s.mu.Lock()
	if s.failDone > 0 {
		s.failDone--
		s.mu.Unlock()
		return fmt.Errorf("injected SetDone failure")
	}
	s.mu.Unlock()
	return s.set(id, token, StatusProving, func(j *ProofJob) { j.Status = StatusDone; j.Result = &r })
}

func (s *fakeStore) SetFailed(_ context.Context, id ProofID, token ClaimToken, class FailureClass, cause string) error {
	return s.set(id, token, StatusProving, func(j *ProofJob) { j.Status = StatusFailed; j.Class = class; j.Cause = cause })
}

func (s *fakeStore) Delete(_ context.Context, id ProofID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, id)
	return nil
}

func (s *fakeStore) Get(_ context.Context, id ProofID) (ProofJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[id]
	if !ok {
		return ProofJob{}, ErrNotFound
	}
	return *j, nil
}

func (s *fakeStore) Close() { s.closes.Add(1) }

func newTestService(t *testing.T, engine *fakeEngine, queueSize int, cfg Config) (*Service, *fakeStore) {
	t.Helper()
	store := newFakeStore()
	s, err := New(engine, newFakeQueue(queueSize), store, cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(s.Stop)
	return s, store
}

func waitStatus(t *testing.T, s *Service, id ProofID, want Status) ProofJob {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		job, err := s.GetProof(id)
		if err == nil && job.Status == want {
			return job
		}
		time.Sleep(5 * time.Millisecond)
	}
	job, err := s.GetProof(id)
	t.Fatalf("job %s never reached %s (last: %+v, err=%v)", id, want, job, err)
	return ProofJob{}
}

// ---- scenarios --------------------------------------------------------------

func TestSubmitProof_IsIdempotent(t *testing.T) {
	release := make(chan struct{})
	engine := &fakeEngine{prove: func(ctx context.Context, req ProofRequest) (ProofResult, error) {
		<-release
		return ProofResult{Bucket: req.Bucket, PaddedSigs: req.Sigs}, nil
	}}
	s, _ := newTestService(t, engine, 8, Config{})

	sigs := MakeTestSigs(t, 2)
	id1, err := s.SubmitProof(sigs)
	if err != nil {
		t.Fatalf("submit 1: %v", err)
	}
	id2, err := s.SubmitProof(sigs)
	if err != nil {
		t.Fatalf("submit 2: %v", err)
	}
	if id1 != id2 {
		t.Fatalf("duplicate submit produced a different id: %s vs %s", id1, id2)
	}
	close(release)
	waitStatus(t, s, id1, StatusDone)
	if got := engine.proves.Load(); got != 1 {
		t.Fatalf("engine proved %d times, want exactly 1", got)
	}
}

func TestSubmitProof_QueueFullIsImmediate(t *testing.T) {
	block := make(chan struct{})
	engine := &fakeEngine{prove: func(ctx context.Context, req ProofRequest) (ProofResult, error) {
		<-block
		return ProofResult{}, nil
	}}
	defer close(block)
	s, store := newTestService(t, engine, 1, Config{})

	// First job occupies the worker, second fills the size-1 queue.
	if _, err := s.SubmitProof(MakeTestSigs(t, 1)); err != nil {
		t.Fatalf("submit 1: %v", err)
	}
	// Wait until the worker picked job 1 up so the queue is genuinely empty.
	time.Sleep(50 * time.Millisecond)
	if _, err := s.SubmitProof(MakeTestSigs(t, 2)); err != nil {
		t.Fatalf("submit 2: %v", err)
	}
	id3, err := s.SubmitProof(MakeTestSigs(t, 3))
	if !errors.Is(err, ErrQueueFull) {
		t.Fatalf("want ErrQueueFull, got id=%s err=%v", id3, err)
	}
	// The rolled-back entry must not linger as a ghost queued job.
	if _, err := store.Get(context.Background(), id3); !errors.Is(err, ErrNotFound) {
		t.Fatalf("rejected job left a store entry: %v", err)
	}
}

func TestWorker_SurvivesEnginePanic(t *testing.T) {
	var first atomic.Bool
	first.Store(true)
	engine := &fakeEngine{prove: func(ctx context.Context, req ProofRequest) (ProofResult, error) {
		if first.CompareAndSwap(true, false) {
			panic("poisoned job")
		}
		return ProofResult{Bucket: req.Bucket}, nil
	}}
	s, _ := newTestService(t, engine, 8, Config{})

	id1, err := s.SubmitProof(MakeTestSigs(t, 1))
	if err != nil {
		t.Fatalf("submit 1: %v", err)
	}
	job := waitStatus(t, s, id1, StatusFailed)
	if job.Class != FailureInternal {
		t.Fatalf("panic must classify internal, got %s", job.Class)
	}

	id2, err := s.SubmitProof(MakeTestSigs(t, 2))
	if err != nil {
		t.Fatalf("submit 2: %v", err)
	}
	waitStatus(t, s, id2, StatusDone) // worker is still alive
}

func TestSubmitProof_FastWhileProving(t *testing.T) {
	block := make(chan struct{})
	engine := &fakeEngine{prove: func(ctx context.Context, req ProofRequest) (ProofResult, error) {
		<-block
		return ProofResult{}, nil
	}}
	defer close(block)
	s, _ := newTestService(t, engine, 8, Config{})

	if _, err := s.SubmitProof(MakeTestSigs(t, 1)); err != nil {
		t.Fatalf("submit 1: %v", err)
	}
	time.Sleep(20 * time.Millisecond) // let the worker start blocking

	start := time.Now()
	if _, err := s.SubmitProof(MakeTestSigs(t, 2)); err != nil {
		t.Fatalf("submit 2: %v", err)
	}
	if d := time.Since(start); d > 50*time.Millisecond {
		t.Fatalf("SubmitProof blocked %s while a prove was running", d)
	}
}

func TestJobTimeout_ClassifiesTransient(t *testing.T) {
	engine := &fakeEngine{prove: func(ctx context.Context, req ProofRequest) (ProofResult, error) {
		<-ctx.Done()
		return ProofResult{}, ctx.Err()
	}}
	s, _ := newTestService(t, engine, 8, Config{JobTimeout: 30 * time.Millisecond})

	id, err := s.SubmitProof(MakeTestSigs(t, 1))
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	job := waitStatus(t, s, id, StatusFailed)
	if job.Class != FailureTransient {
		t.Fatalf("timeout must classify transient, got %s (%s)", job.Class, job.Cause)
	}
}

func TestSubmitProof_RequeuesFailedJob(t *testing.T) {
	var first atomic.Bool
	first.Store(true)
	engine := &fakeEngine{prove: func(_ context.Context, req ProofRequest) (ProofResult, error) {
		if first.CompareAndSwap(true, false) {
			return ProofResult{}, context.DeadlineExceeded
		}
		return ProofResult{Bucket: req.Bucket, PaddedSigs: req.Sigs}, nil
	}}
	s, _ := newTestService(t, engine, 8, Config{})
	sigs := MakeTestSigs(t, 1)

	id, err := s.SubmitProof(sigs)
	if err != nil {
		t.Fatalf("first submit: %v", err)
	}
	job := waitStatus(t, s, id, StatusFailed)
	if job.Class != FailureTransient {
		t.Fatalf("first attempt class = %s, want transient", job.Class)
	}
	retryID, err := s.SubmitProof(sigs)
	if err != nil {
		t.Fatalf("retry submit: %v", err)
	}
	if retryID != id {
		t.Fatalf("retry id = %s, want deterministic id %s", retryID, id)
	}
	waitStatus(t, s, id, StatusDone)
	if got := engine.proves.Load(); got != 2 {
		t.Fatalf("engine prove calls = %d, want 2", got)
	}
}

func TestStop_ClosesStoreBackgroundWork(t *testing.T) {
	store := newFakeStore()
	s, err := New(&fakeEngine{}, newFakeQueue(1), store, Config{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	s.Stop()
	s.Stop()
	if got := store.closes.Load(); got != 1 {
		t.Fatalf("store Close calls = %d, want 1", got)
	}
}

func TestWaitProof_ReturnsResult(t *testing.T) {
	engine := &fakeEngine{}
	s, _ := newTestService(t, engine, 8, Config{PollInterval: 5 * time.Millisecond})

	id, err := s.SubmitProof(MakeTestSigs(t, 1))
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := s.WaitProof(ctx, id)
	if err != nil {
		t.Fatalf("wait: %v", err)
	}
	if res.Bucket == 0 {
		t.Fatal("empty result")
	}
}

func TestWaitProof_ReturnsTypedFailure(t *testing.T) {
	engine := &fakeEngine{prove: func(context.Context, ProofRequest) (ProofResult, error) {
		return ProofResult{}, errors.New("broken artifacts")
	}}
	s, _ := newTestService(t, engine, 8, Config{PollInterval: time.Millisecond})
	id, err := s.SubmitProof(MakeTestSigs(t, 1))
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = s.WaitProof(ctx, id)
	var failed *ErrProofFailed
	if !errors.As(err, &failed) {
		t.Fatalf("error = %v, want ErrProofFailed", err)
	}
	if failed.ID != id || failed.Class != FailureInternal || failed.Cause != "broken artifacts" {
		t.Fatalf("typed failure = %+v", failed)
	}
}

func TestSubmitProof_InvalidInputAndStoreDown(t *testing.T) {
	engine := &fakeEngine{}
	s, store := newTestService(t, engine, 8, Config{})

	var invalid *ErrInvalidInput
	if _, err := s.SubmitProof(nil); !errors.As(err, &invalid) {
		t.Fatalf("want ErrInvalidInput, got %v", err)
	}

	store.mu.Lock()
	store.failPut = fmt.Errorf("backend down")
	store.mu.Unlock()
	if _, err := s.SubmitProof(MakeTestSigs(t, 1)); !errors.Is(err, ErrStoreUnavailable) {
		t.Fatalf("want ErrStoreUnavailable, got %v", err)
	}
}

func TestService_ParallelSubmitAndQuery(t *testing.T) {
	engine := &fakeEngine{}
	s, _ := newTestService(t, engine, 64, Config{})

	var wg sync.WaitGroup
	ids := make([]ProofID, 8)
	for i := range ids {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id, err := s.SubmitProof(MakeTestSigs(t, i%3+1))
			if err != nil && !errors.Is(err, ErrQueueFull) {
				t.Errorf("submit %d: %v", i, err)
				return
			}
			ids[i] = id
			for j := 0; j < 20; j++ {
				if _, err := s.GetProof(id); err != nil && !errors.Is(err, ErrNotFound) {
					t.Errorf("get %d: %v", i, err)
					return
				}
			}
		}(i)
	}
	wg.Wait()
}

func TestWorker_RedeliveryAckSemantics(t *testing.T) {
	engine := &fakeEngine{}
	queue := newFakeQueue(8)
	store := newFakeStore()
	s, err := New(engine, queue, store, Config{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(s.Stop)

	// Terminal job redelivered → dropped with an ack, engine untouched.
	doneID := ProofID{0xD0}
	if _, err := store.PutPending(context.Background(), doneID); err != nil {
		t.Fatal(err)
	}
	token, claimed, err := store.Claim(context.Background(), doneID, "setup", time.Minute)
	if err != nil || !claimed {
		t.Fatalf("claim done fixture: claimed=%v err=%v", claimed, err)
	}
	if err := store.SetProving(context.Background(), doneID, token); err != nil {
		t.Fatal(err)
	}
	if err := store.SetDone(context.Background(), doneID, token, ProofResult{Bucket: 4}); err != nil {
		t.Fatal(err)
	}
	if err := queue.Enqueue(context.Background(), ProofRequest{ID: doneID, Bucket: 4}); err != nil {
		t.Fatal(err)
	}

	// Unknown id redelivered (rollback/TTL took the entry) → dropped, no error loop.
	if err := queue.Enqueue(context.Background(), ProofRequest{ID: ProofID{0xEE}, Bucket: 4}); err != nil {
		t.Fatal(err)
	}

	// A live job still proves normally after the two drops above.
	id, err := s.SubmitProof(MakeTestSigs(t, 1))
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	waitStatus(t, s, id, StatusDone)
	if got := engine.proves.Load(); got != 1 {
		t.Fatalf("engine proved %d times, want 1 (redeliveries must not prove)", got)
	}
}

// TestWorker_AckOnlyAfterStore pins the ordering contract redelivery rests
// on: the ack happens strictly AFTER the outcome is stored. A SetDone failure
// must leave the job un-acked so the backend redelivers and the result is
// re-proven rather than lost.
func TestWorker_AckOnlyAfterStore(t *testing.T) {
	engine := &fakeEngine{}
	queue := newFakeQueue(8)
	store := newFakeStore()
	store.failDone = 1 // first SetDone fails after the prove succeeded
	s, err := New(engine, queue, store, Config{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(s.Stop)

	id, err := s.SubmitProof(MakeTestSigs(t, 1))
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	// The worker proves, SetDone fails, and it must NOT ack.
	deadline := time.Now().Add(2 * time.Second)
	for engine.proves.Load() < 1 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	time.Sleep(50 * time.Millisecond) // let the failed SetDone path finish
	if queue.acks.Load() != 0 {
		t.Fatalf("job acked despite SetDone failure — result would be lost on a redelivering backend")
	}

	// Simulate the backend redelivery the un-acked job earns; the second pass
	// re-proves (the store CAS reset is the mem-store's job — here the fake
	// keeps it claimable) and this time stores + acks.
	store.mu.Lock()
	store.jobs[id].Status = StatusQueued // redelivery finds it reclaimable
	store.mu.Unlock()
	if err := queue.Enqueue(context.Background(), ProofRequest{ID: id, Bucket: 4, Sigs: nil}); err != nil {
		t.Fatalf("redeliver: %v", err)
	}
	waitStatus(t, s, id, StatusDone)
	if queue.acks.Load() == 0 {
		t.Fatal("successful SetDone must be followed by an ack")
	}
}

// TestShutdown_AbandonsWithoutFailing pins the shutdown row of the error
// catalog: a service stop mid-prove must NOT mark the job failed — it is
// unfinished, not broken — and must not ack it.
func TestShutdown_AbandonsWithoutFailing(t *testing.T) {
	proving := make(chan struct{})
	engine := &fakeEngine{prove: func(ctx context.Context, req ProofRequest) (ProofResult, error) {
		close(proving)
		<-ctx.Done() // blocks until the service ctx is cancelled by Stop
		return ProofResult{}, ctx.Err()
	}}
	queue := newFakeQueue(8)
	store := newFakeStore()
	s, err := New(engine, queue, store, Config{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	id, err := s.SubmitProof(MakeTestSigs(t, 1))
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	<-proving // the worker is inside the prove
	s.Stop()  // shutdown mid-prove

	job, err := store.Get(context.Background(), id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if job.Status == StatusFailed {
		t.Fatalf("shutdown marked the job failed (%s: %s) — it is unfinished, not broken", job.Class, job.Cause)
	}
	if queue.acks.Load() != 0 {
		t.Fatal("shutdown must not ack the in-flight job")
	}
}
