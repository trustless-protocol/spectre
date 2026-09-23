// Package memqueue is the in-memory JobQueue adapter: a bounded channel.
//
// Contract notes: delivery is at-most-once in practice — an un-acked job dies
// with the process, which the design accepts (deterministic ProofIDs make a
// re-submit reproduce the job). ack is a no-op; the at-least-once semantics
// the Service is written against become real when an external backend
// (Redis/Kafka) replaces this adapter.
package memqueue

import (
	"context"

	"relayer/prover"
)

type Queue struct {
	ch chan prover.ProofRequest
}

// New returns a queue holding at most size jobs; size <= 0 defaults to 64.
func New(size int) *Queue {
	if size <= 0 {
		size = 64
	}
	return &Queue{ch: make(chan prover.ProofRequest, size)}
}

// Enqueue never blocks: a full queue returns ErrQueueFull immediately so the
// caller can apply its own backpressure policy.
func (q *Queue) Enqueue(_ context.Context, job prover.ProofRequest) error {
	select {
	case q.ch <- job:
		return nil
	default:
		return prover.ErrQueueFull
	}
}

// Dequeue blocks until a job arrives or ctx is cancelled.
func (q *Queue) Dequeue(ctx context.Context) (prover.ProofRequest, func(), error) {
	select {
	case <-ctx.Done():
		return prover.ProofRequest{}, nil, ctx.Err()
	case job := <-q.ch:
		return job, func() {}, nil
	}
}

// Depth reports the number of queued jobs (for metrics/monitoring).
func (q *Queue) Depth() int { return len(q.ch) }
