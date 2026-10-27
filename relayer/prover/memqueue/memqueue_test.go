package memqueue

import (
	"context"
	"errors"
	"testing"
	"time"

	"relayer/prover"
)

func TestEnqueue_FullRejectsImmediately(t *testing.T) {
	q := New(1)
	ctx := context.Background()
	if err := q.Enqueue(ctx, prover.ProofRequest{ID: prover.ProofID{1}}); err != nil {
		t.Fatalf("enqueue 1: %v", err)
	}
	start := time.Now()
	err := q.Enqueue(ctx, prover.ProofRequest{ID: prover.ProofID{2}})
	if !errors.Is(err, prover.ErrQueueFull) {
		t.Fatalf("want ErrQueueFull, got %v", err)
	}
	if time.Since(start) > 10*time.Millisecond {
		t.Fatal("full enqueue must not block")
	}
}

func TestDequeue_DeliversInOrderAndHonorsCtx(t *testing.T) {
	q := New(4)
	ctx := context.Background()
	for i := byte(1); i <= 3; i++ {
		if err := q.Enqueue(ctx, prover.ProofRequest{ID: prover.ProofID{i}}); err != nil {
			t.Fatalf("enqueue %d: %v", i, err)
		}
	}
	for i := byte(1); i <= 3; i++ {
		job, ack, err := q.Dequeue(ctx)
		if err != nil {
			t.Fatalf("dequeue %d: %v", i, err)
		}
		if job.ID != (prover.ProofID{i}) {
			t.Fatalf("dequeue %d returned %s", i, job.ID)
		}
		ack()
	}

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := q.Dequeue(cancelled); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled dequeue: %v", err)
	}
}
