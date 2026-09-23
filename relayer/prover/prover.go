package prover

import (
	"relayer/prover/witness"

	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

// Config tunes the Service. Zero values get defaults; the defaults suit the
// in-process, single-machine deployment this phase ships.
type Config struct {
	// Workers is the number of concurrent prove workers. Default 1: a gnark
	// prove already saturates every core, more workers buy RAM, not
	// throughput.
	Workers int
	// JobTimeout cancels a single prove that runs away. The default (10m) is
	// deliberately a safe ceiling rather than the doc's ~10× bucket baseline:
	// baselines exist only for N=4/N=16, and a ceiling tighter than an
	// unmeasured N=64/128 prove would misclassify valid work as transient and
	// put the relayer in a resubmit loop. Tighten per bucket once task #13
	// measures the large buckets; operators can already override it here.
	JobTimeout time.Duration
	// Lease is how long a Claim holds a job before a dead worker's job
	// becomes reclaimable. Default 2×JobTimeout.
	Lease time.Duration
	// PollInterval paces WaitProof's store polling. Default 200ms.
	PollInterval time.Duration
}

func (c Config) withDefaults() Config {
	if c.Workers <= 0 {
		c.Workers = 1
	}
	if c.JobTimeout <= 0 {
		c.JobTimeout = 10 * time.Minute
	}
	if c.Lease <= 0 {
		c.Lease = 2 * c.JobTimeout
	}
	if c.PollInterval <= 0 {
		c.PollInterval = 200 * time.Millisecond
	}
	return c
}

// Service is the async proof API: SubmitProof returns a deterministic ProofID
// immediately, workers prove in the background, GetProof reads the outcome.
// The Service holds no policy — no retries, no prioritization; a failed job
// stays failed in the store and what to do about it is the caller's decision.
type Service struct {
	engine  ProofEngine
	queue   JobQueue
	store   ProofStore
	metrics Metrics
	cfg     Config

	buckets []int // engine's loaded buckets, ascending

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	stop   sync.Once
}

// New wires the ports together and starts the worker pool. It fails when the
// engine declares no loaded buckets — a service that can never prove should
// die at startup, not at the first submit. Call Stop to shut down; in-flight
// jobs are abandoned un-acked (not failed) so a redelivering backend hands
// them to the next start.
func New(engine ProofEngine, queue JobQueue, store ProofStore, cfg Config) (*Service, error) {
	return NewWithMetrics(engine, queue, store, cfg, NopMetrics{})
}

// NewWithMetrics is New with a Metrics sink.
func NewWithMetrics(engine ProofEngine, queue JobQueue, store ProofStore, cfg Config, m Metrics) (*Service, error) {
	if m == nil {
		m = NopMetrics{}
	}
	buckets := append([]int(nil), engine.Buckets()...)
	if len(buckets) == 0 {
		return nil, fmt.Errorf("proof engine declares no loaded buckets")
	}
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		engine:  engine,
		queue:   queue,
		store:   store,
		metrics: m,
		cfg:     cfg.withDefaults(),
		buckets: buckets,
		ctx:     ctx,
		cancel:  cancel,
	}
	for i := 0; i < s.cfg.Workers; i++ {
		s.wg.Add(1)
		go s.worker(fmt.Sprintf("worker-%d", i))
	}
	return s, nil
}

// MaxBucket returns the largest bucket the engine can prove — the ceiling
// selection policy in the relayer must respect.
func (s *Service) MaxBucket() int { return s.buckets[len(s.buckets)-1] }

// Stop cancels the workers, waits for them to exit, and closes an adapter that
// owns background work (the in-memory store's TTL sweeper). It is idempotent.
func (s *Service) Stop() {
	s.stop.Do(func() {
		s.cancel()
		s.wg.Wait()
		if closer, ok := s.store.(interface{ Close() }); ok {
			closer.Close()
		}
	})
}

// SubmitProof validates cheaply, derives the ProofID, and enqueues the job.
// It never blocks on proving: every return path is store/queue bookkeeping.
// A duplicate submit (same signatures, same circuit) returns the existing id
// without a second prove.
func (s *Service) SubmitProof(sigs []witness.ValidatorSignature) (ProofID, error) {
	if err := validateSignatures(sigs); err != nil {
		return ProofID{}, err
	}
	bucket, err := witness.SmallestBucketGEQ(len(sigs), s.buckets)
	if err != nil {
		return ProofID{}, &ErrInvalidInput{Reason: err.Error(), Slot: -1}
	}
	padded, err := witness.PadWithDummies(sigs, bucket)
	if err != nil {
		return ProofID{}, &ErrInvalidInput{Reason: err.Error(), Slot: -1}
	}
	witnessDigest, err := witness.ComputeWitnessHash(padded)
	if err != nil {
		return ProofID{}, &ErrInvalidInput{Reason: err.Error(), Slot: -1}
	}
	id := NewProofID(s.engine.VKDigest(bucket), witnessDigest)

	created, err := s.store.PutPending(s.ctx, id)
	if err != nil {
		return ProofID{}, fmt.Errorf("%w: put pending: %w", ErrStoreUnavailable, err)
	}
	if !created {
		// Dedup: the job is queued, in flight, or done; the caller polls the
		// same id. A failed entry is reset by PutPending and returns created=true.
		return id, nil
	}
	req := ProofRequest{ID: id, Bucket: bucket, Sigs: padded, WitnessDigest: witnessDigest, EnqueuedAt: time.Now()}
	if err := s.queue.Enqueue(s.ctx, req); err != nil {
		// Roll the pending entry back so no ghost job sits queued forever;
		// if the delete also fails, the store TTL collects it.
		if delErr := s.store.Delete(s.ctx, id); delErr != nil {
			log.Printf("[prover] proof_id=%s rollback delete failed: %v", id, delErr)
		}
		if errors.Is(err, ErrQueueFull) {
			s.metrics.QueueFull()
			return ProofID{}, err
		}
		return ProofID{}, fmt.Errorf("enqueue: %w", err)
	}
	s.metrics.JobQueued()
	log.Printf("[prover] proof_id=%s bucket=%d witness=%x queued", id, bucket, witnessDigest[:8])
	return id, nil
}

// GetProof reads the job's current state. Unknown ids return ErrNotFound —
// never submitted, TTL-collected, or lost to a restart; re-submitting
// reproduces the same id.
func (s *Service) GetProof(id ProofID) (ProofJob, error) {
	return s.store.Get(s.ctx, id)
}

// WaitProof polls until the job finishes or ctx expires. It is a convenience
// for one-shot CLI flows; the relay loop uses SubmitProof + GetProof on its
// own cadence. The job keeps running when WaitProof gives up.
func (s *Service) WaitProof(ctx context.Context, id ProofID) (ProofResult, error) {
	ticker := time.NewTicker(s.cfg.PollInterval)
	defer ticker.Stop()
	for {
		job, err := s.store.Get(ctx, id)
		if err != nil {
			return ProofResult{}, err
		}
		switch job.Status {
		case StatusDone:
			return *job.Result, nil
		case StatusFailed:
			return ProofResult{}, &ErrProofFailed{ID: id, Class: job.Class, Cause: job.Cause}
		}
		select {
		case <-ctx.Done():
			return ProofResult{}, ctx.Err()
		case <-ticker.C:
		}
	}
}

// worker pulls jobs until the service stops. A failing queue backend puts the
// worker into backoff-reconnect, never a busy loop, never an exit.
func (s *Service) worker(workerID string) {
	defer s.wg.Done()
	for {
		job, ack, err := s.queue.Dequeue(s.ctx)
		if err != nil {
			if s.ctx.Err() != nil {
				return
			}
			log.Printf("[prover] %s dequeue error, backing off: %v", workerID, err)
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(time.Second):
			}
			continue
		}
		s.process(workerID, job, ack)
	}
}

// process drives one job through claim → prove → store → ack. Invariants:
// every error either reaches the store as failed or leaves the job un-acked
// for redelivery — no third path; ack strictly after the outcome is stored;
// a shutdown mid-prove neither fails nor acks the job.
func (s *Service) process(workerID string, job ProofRequest, ack func()) {
	token, claimed, err := s.store.Claim(s.ctx, job.ID, workerID, s.cfg.Lease)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			// The entry is gone — enqueue rollback or TTL collected it. The
			// job is unrecoverable from here; a re-submit owns recovery.
			ack()
			return
		}
		log.Printf("[prover] proof_id=%s claim error: %v", job.ID, err)
		return // no ack → redelivery
	}
	if !claimed {
		// Ack only when the job is terminal. A job claimed/proving under a
		// live lease must keep its redelivery un-acked: if that worker dies,
		// the redelivery is what hands the job to the next worker after the
		// lease expires. (mem queue: ack is a no-op either way.)
		// ponytail: an un-acked redelivery can retry hot on some backends;
		// add a nack-with-delay to the JobQueue port if that ever matters.
		if cur, err := s.store.Get(s.ctx, job.ID); err == nil &&
			(cur.Status == StatusDone || cur.Status == StatusFailed) {
			ack()
		}
		return
	}
	if err := s.store.SetProving(s.ctx, job.ID, token); err != nil {
		log.Printf("[prover] proof_id=%s set proving error: %v", job.ID, err)
		return // no ack → redelivery; dedup absorbs the duplicate
	}
	s.metrics.JobStarted(time.Since(job.EnqueuedAt))

	start := time.Now()
	log.Printf("[prover] proof_id=%s bucket=%d proving", job.ID, job.Bucket)
	result, err := s.prove(job)
	if err != nil {
		if s.ctx.Err() != nil {
			// Shutdown, not failure: the job is merely unfinished. No
			// SetFailed, no ack — redelivery (or a re-submit) finishes it.
			log.Printf("[prover] proof_id=%s abandoned by shutdown", job.ID)
			return
		}
		class := FailureInternal
		if errors.Is(err, context.DeadlineExceeded) {
			class = FailureTransient
		}
		s.metrics.JobFailed(class)
		log.Printf("[prover] proof_id=%s failed class=%s cause=%v", job.ID, class, err)
		if serr := s.store.SetFailed(s.ctx, job.ID, token, class, err.Error()); serr != nil {
			log.Printf("[prover] proof_id=%s set failed error: %v", job.ID, serr)
			return // no ack → redelivery re-proves; wasteful, never wrong
		}
		ack()
		return
	}
	if serr := s.store.SetDone(s.ctx, job.ID, token, result); serr != nil {
		log.Printf("[prover] proof_id=%s set done error: %v", job.ID, serr)
		return // no ack → redelivery re-proves rather than losing the result
	}
	s.metrics.JobDone(time.Since(start))
	log.Printf("[prover] proof_id=%s done in %s", job.ID, time.Since(start))
	ack()
}

// prove calls the engine with a per-job timeout and converts a panic into an
// error — a poisoned job must not take the worker down with it.
func (s *Service) prove(job ProofRequest) (result ProofResult, err error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.cfg.JobTimeout)
	defer cancel()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("engine panic: %v\n%s", r, debug.Stack())
		}
	}()
	return s.engine.Prove(ctx, job)
}
