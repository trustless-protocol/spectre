package services

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"relayer/prover"
	"relayer/prover/memqueue"
	"relayer/prover/memstore"
	"relayer/prover/witness"
)

type proofClientFake struct {
	submit func(int) (prover.ProofID, error)
	get    func(int, prover.ProofID) (prover.ProofJob, error)
	wait   func(context.Context, prover.ProofID) (prover.ProofResult, error)

	submitCalls int
	getCalls    int
	waitCalls   int
}

func (f *proofClientFake) SubmitProof([]witness.ValidatorSignature) (prover.ProofID, error) {
	f.submitCalls++
	return f.submit(f.submitCalls)
}

func (*proofClientFake) MaxBucket() int { return 4 }

func (f *proofClientFake) GetProof(id prover.ProofID) (prover.ProofJob, error) {
	f.getCalls++
	return f.get(f.getCalls, id)
}

func (f *proofClientFake) WaitProof(ctx context.Context, id prover.ProofID) (prover.ProofResult, error) {
	f.waitCalls++
	if f.wait == nil {
		return prover.ProofResult{}, errors.New("unexpected WaitProof call")
	}
	return f.wait(ctx, id)
}

func TestSubmitAndAwaitProof_RetryTable(t *testing.T) {
	id := prover.ProofID{1}
	done := prover.ProofResult{Bucket: 4}
	tests := []struct {
		name        string
		fake        *proofClientFake
		wantErr     error
		wantSubmits int
	}{
		{
			name: "invalid input is permanent",
			fake: &proofClientFake{
				submit: func(int) (prover.ProofID, error) {
					return prover.ProofID{}, &prover.ErrInvalidInput{Reason: "bad", Slot: -1}
				},
				get: func(int, prover.ProofID) (prover.ProofJob, error) { return prover.ProofJob{}, nil },
			},
			wantErr:     ErrPermanentRelayFailure,
			wantSubmits: 1,
		},
		{
			name: "queue full backs off and resubmits",
			fake: &proofClientFake{
				submit: func(call int) (prover.ProofID, error) {
					if call == 1 {
						return prover.ProofID{}, prover.ErrQueueFull
					}
					return id, nil
				},
				get: func(int, prover.ProofID) (prover.ProofJob, error) {
					return prover.ProofJob{ID: id, Status: prover.StatusDone, Result: &done}, nil
				},
			},
			wantSubmits: 2,
		},
		{
			name: "store unavailable backs off and resubmits",
			fake: &proofClientFake{
				submit: func(call int) (prover.ProofID, error) {
					if call == 1 {
						return prover.ProofID{}, prover.ErrStoreUnavailable
					}
					return id, nil
				},
				get: func(int, prover.ProofID) (prover.ProofJob, error) {
					return prover.ProofJob{ID: id, Status: prover.StatusDone, Result: &done}, nil
				},
			},
			wantSubmits: 2,
		},
		{
			name: "not found resubmits",
			fake: &proofClientFake{
				submit: func(int) (prover.ProofID, error) { return id, nil },
				get: func(call int, _ prover.ProofID) (prover.ProofJob, error) {
					if call == 1 {
						return prover.ProofJob{}, prover.ErrNotFound
					}
					return prover.ProofJob{ID: id, Status: prover.StatusDone, Result: &done}, nil
				},
			},
			wantSubmits: 2,
		},
		{
			name: "transient failure resubmits",
			fake: &proofClientFake{
				submit: func(int) (prover.ProofID, error) { return id, nil },
				get: func(call int, _ prover.ProofID) (prover.ProofJob, error) {
					if call == 1 {
						return prover.ProofJob{ID: id, Status: prover.StatusFailed, Class: prover.FailureTransient, Cause: "timeout"}, nil
					}
					return prover.ProofJob{ID: id, Status: prover.StatusDone, Result: &done}, nil
				},
			},
			wantSubmits: 2,
		},
		{
			name: "internal failure is permanent",
			fake: &proofClientFake{
				submit: func(int) (prover.ProofID, error) { return id, nil },
				get: func(int, prover.ProofID) (prover.ProofJob, error) {
					return prover.ProofJob{ID: id, Status: prover.StatusFailed, Class: prover.FailureInternal, Cause: "verify"}, nil
				},
			},
			wantErr:     ErrPermanentRelayFailure,
			wantSubmits: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			result, err := submitAndAwaitProofWithTiming(ctx, tt.fake, nil, time.Millisecond, time.Millisecond)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
			} else if err != nil || result.Bucket != done.Bucket {
				t.Fatalf("result=%+v err=%v", result, err)
			}
			if tt.fake.submitCalls != tt.wantSubmits {
				t.Fatalf("SubmitProof calls = %d, want %d", tt.fake.submitCalls, tt.wantSubmits)
			}
		})
	}
}

func TestNewOneShotWorker_UsesWaitProof(t *testing.T) {
	id := prover.ProofID{1}
	fake := &proofClientFake{
		submit: func(int) (prover.ProofID, error) { return id, nil },
		get: func(int, prover.ProofID) (prover.ProofJob, error) {
			return prover.ProofJob{}, errors.New("unexpected GetProof call")
		},
		wait: func(context.Context, prover.ProofID) (prover.ProofResult, error) {
			return prover.ProofResult{Bucket: 4}, nil
		},
	}
	worker := NewOneShotWorker(nil, fake)
	result, err := worker.generateProof(context.Background(), nil)
	if err != nil || result.Bucket != 4 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if fake.waitCalls != 1 || fake.getCalls != 0 {
		t.Fatalf("WaitProof/GetProof calls = %d/%d, want 1/0", fake.waitCalls, fake.getCalls)
	}
}

func TestSubmitAndWaitProof_TransientFailureResubmits(t *testing.T) {
	id := prover.ProofID{1}
	fake := &proofClientFake{
		submit: func(int) (prover.ProofID, error) { return id, nil },
		get: func(int, prover.ProofID) (prover.ProofJob, error) {
			return prover.ProofJob{}, errors.New("unexpected GetProof call")
		},
	}
	var waits int
	fake.wait = func(context.Context, prover.ProofID) (prover.ProofResult, error) {
		waits++
		if waits == 1 {
			return prover.ProofResult{}, &prover.ErrProofFailed{ID: id, Class: prover.FailureTransient, Cause: "timeout"}
		}
		return prover.ProofResult{Bucket: 4}, nil
	}
	result, err := submitAndWaitProof(context.Background(), fake, nil)
	if err != nil || result.Bucket != 4 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if fake.submitCalls != 2 || fake.waitCalls != 2 {
		t.Fatalf("SubmitProof/WaitProof calls = %d/%d, want 2/2", fake.submitCalls, fake.waitCalls)
	}
}

type failOnceProofEngine struct{ calls atomic.Int64 }

func (e *failOnceProofEngine) Prove(_ context.Context, req prover.ProofRequest) (prover.ProofResult, error) {
	if e.calls.Add(1) == 1 {
		return prover.ProofResult{}, context.DeadlineExceeded
	}
	return prover.ProofResult{Bucket: req.Bucket, PaddedSigs: req.Sigs}, nil
}

func (*failOnceProofEngine) Buckets() []int        { return []int{4} }
func (*failOnceProofEngine) VKDigest(int) [32]byte { return [32]byte{0x42} }

func testProofSignatures(t *testing.T) []witness.ValidatorSignature {
	t.Helper()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	message := make([]byte, 120)
	message[0] = 119
	message[1] = 0x08
	message[2] = 0x02
	message[3] = 0x11
	for i := 0; i < 8; i++ {
		message[4+i] = byte(i + 1)
	}
	message[12] = 0x22
	message[13] = 0x48
	message[14] = 0x0a
	message[15] = 0x20
	for i := 0; i < 32; i++ {
		message[16+i] = byte(0xa0 + i)
	}
	for i := 48; i < len(message); i++ {
		message[i] = byte(i)
	}
	return []witness.ValidatorSignature{{
		Index:       0,
		Signature:   ed25519.Sign(privateKey, message),
		PublicKey:   publicKey,
		SignedBytes: message,
		Power:       1,
		Active:      true,
	}}
}

func TestSubmitAndAwaitProof_TransientFailureIsProvedAgain(t *testing.T) {
	engine := &failOnceProofEngine{}
	store := memstore.New(time.Minute, 8)
	service, err := prover.New(engine, memqueue.New(8), store, prover.Config{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(service.Stop)

	sigs := testProofSignatures(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := submitAndAwaitProofWithTiming(ctx, service, sigs, time.Millisecond, time.Millisecond)
	if err != nil {
		t.Fatalf("submitAndAwaitProof: %v", err)
	}
	if result.Bucket != 4 || engine.calls.Load() != 2 {
		t.Fatalf("result bucket=%d prove calls=%d, want 4/2", result.Bucket, engine.calls.Load())
	}
}

type slowProofEngine struct{ calls atomic.Int64 }

func (e *slowProofEngine) Prove(_ context.Context, req prover.ProofRequest) (prover.ProofResult, error) {
	e.calls.Add(1)
	time.Sleep(60 * time.Millisecond) // longer than the lease; Gnark prove is non-interruptible
	return prover.ProofResult{Bucket: req.Bucket, PaddedSigs: req.Sigs}, nil
}
func (*slowProofEngine) Buckets() []int        { return []int{4} }
func (*slowProofEngine) VKDigest(int) [32]byte { return [32]byte{0x42} }

func TestSubmitAndAwaitProof_CompletesAfterLeaseWithMemqueue(t *testing.T) {
	engine := &slowProofEngine{}
	store := memstore.New(time.Minute, 8)
	service, err := prover.New(engine, memqueue.New(8), store, prover.Config{Lease: 10 * time.Millisecond, JobTimeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(service.Stop)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := submitAndAwaitProofWithTiming(ctx, service, testProofSignatures(t), time.Millisecond, time.Millisecond)
	if err != nil {
		t.Fatalf("proof stuck after lease expiry: %v", err)
	}
	if result.Bucket != 4 || engine.calls.Load() != 1 {
		t.Fatalf("result bucket=%d prove calls=%d, want 4/1", result.Bucket, engine.calls.Load())
	}
}
