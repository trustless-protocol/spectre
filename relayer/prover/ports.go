package prover

import (
	"context"
	"time"
)

// The Service owns these ports (Go style: interfaces declared where they are
// consumed). Adapters implement them in subpackages; the Service never learns
// which one it got. Semantics every implementation must honor — the contract
// test suite in porttest is the gate:
//
//   - JobQueue delivery is at-least-once, ordering best-effort FIFO with no
//     guarantee. Ack only after the result is in the store; an unacked job is
//     redelivered, and the store's Claim/dedup absorbs the duplicate.
//   - ProofStore.PutPending is an atomic check-and-create (the dedup barrier),
//     except that an existing failed entry is reset to queued: an explicit
//     caller resubmit must be able to prove the deterministic id again. Claim
//     is an atomic queued→claimed CAS with an expiring lease; TTL/GC is the
//     store's own job.
//   - PutPending → Enqueue is deliberately not atomic: an enqueue failure
//     rolls the pending entry back best-effort, and an entry stuck queued past
//     its TTL is dead — the TTL cleans it and a re-submit re-enqueues the same
//     id.
//
// The sentinel errors these contracts speak in (ErrQueueFull, ErrNotFound,
// ErrStoreUnavailable) live in the proof package next to the job model, so
// adapters depend only on proof.

// ProofEngine produces verified proofs. Implementations are stateless after
// construction, safe for concurrent calls, and MUST local-verify against the
// vk whose digest they declare before returning — no unverified proof leaves
// an engine. Cancelling ctx abandons the attempt.
type ProofEngine interface {
	Prove(ctx context.Context, req ProofRequest) (ProofResult, error)
	// Buckets declares which bucket sizes this engine has artifacts for.
	Buckets() []int
	// VKDigest returns the SHA-256 of the verifying key for a bucket — the
	// circuit-version half of every ProofID.
	VKDigest(bucket int) [32]byte
}

// JobQueue transports ProofRequests from submit to a worker.
type JobQueue interface {
	// Enqueue rejects immediately with ErrQueueFull when full — never
	// blocks, never drops silently.
	Enqueue(ctx context.Context, job ProofRequest) error
	// Dequeue blocks until a job is available or ctx is cancelled. ack must be
	// called only after the job's outcome is stored; not acking hands the job
	// back to the backend for redelivery.
	Dequeue(ctx context.Context) (job ProofRequest, ack func(), err error)
}

// ProofStore keeps job state and results. It is a cache with a TTL, not a
// ledger — losing it costs a re-prove, never correctness.
type ProofStore interface {
	// PutPending atomically creates the entry. It returns created=false when a
	// queued, in-flight, or done id already exists (the dedup barrier), and
	// resets failed→queued with created=true so an explicit resubmit enqueues a
	// fresh attempt.
	PutPending(ctx context.Context, id ProofID) (created bool, err error)
	// Claim atomically moves queued→claimed and takes an expiring lease.
	// ok=false when the job is not claimable: already done/failed, or
	// claimed/proving under a live lease. An expired lease makes the job
	// claimable again — a dead worker never strands a job. Claim returns a
	// fresh token on success. Lease expiry permits a new Claim; it does not
	// revoke the owner until a replacement claim succeeds. The Set methods
	// reject superseded tokens with ErrStaleClaim.
	Claim(ctx context.Context, id ProofID, workerID string, lease time.Duration) (token ClaimToken, ok bool, err error)
	SetProving(ctx context.Context, id ProofID, token ClaimToken) error
	SetDone(ctx context.Context, id ProofID, token ClaimToken, r ProofResult) error
	SetFailed(ctx context.Context, id ProofID, token ClaimToken, class FailureClass, cause string) error
	// Delete removes an entry; used for best-effort rollback when Enqueue
	// fails after PutPending. Deleting a missing id is not an error.
	Delete(ctx context.Context, id ProofID) error
	// Get returns the job in any state, or ErrNotFound.
	Get(ctx context.Context, id ProofID) (ProofJob, error)
}

// Metrics receives lifecycle signals. The zero implementation is a no-op;
// wiring real metrics is refactor task #11.
type Metrics interface {
	JobQueued()
	JobStarted(queueWait time.Duration)
	JobDone(proveDuration time.Duration)
	JobFailed(class FailureClass)
	QueueFull()
}

// NopMetrics is the default Metrics implementation.
type NopMetrics struct{}

func (NopMetrics) JobQueued()               {}
func (NopMetrics) JobStarted(time.Duration) {}
func (NopMetrics) JobDone(time.Duration)    {}
func (NopMetrics) JobFailed(FailureClass)   {}
func (NopMetrics) QueueFull()               {}
