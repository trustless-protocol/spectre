package prover_test

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"relayer/prover"
	"relayer/prover/gnark"
	"relayer/prover/memqueue"
	"relayer/prover/memstore"
	"relayer/prover/witness"
)

// TestProverE2E drives the whole prover in isolation — no relayer, no chain:
// real artifacts on disk, the real gnark engine, the real in-memory queue and
// store, real Ed25519 signatures, and a real Groth16 prove + local verify.
//
// Artifacts: set PROVER_BIN_DIR to reuse existing bin/n{N} artifacts (fast).
// Unset, the test generates dev artifacts for every configured bucket into a
// temp dir via gnark.SetupBucket — slower (compile + setup, minutes) but
// fully self-contained and it additionally covers the setup→persist→load
// round-trip. Skipped under -short either way.
func TestProverE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("e2e: real Groth16 prove; skipped in -short mode")
	}
	if raceDetectorEnabled {
		// The race detector slows Groth16 setup and proving ~10-100x (a
		// single N=4 prove exceeds five minutes), and the ZK math is
		// single-threaded per prove — the Service's concurrency is
		// race-covered by the fake-engine unit tests instead.
		t.Skip("e2e: skipped under -race; run without -race for the real-prove path")
	}

	binDir := os.Getenv("PROVER_BIN_DIR")
	if binDir == "" {
		binDir = t.TempDir()
		t.Logf("PROVER_BIN_DIR not set: generating dev artifacts for bucket n=4 into %s (slow, dev-only trusted setup)", binDir)
		start := time.Now()
		if _, _, _, err := gnark.SetupBucket(binDir, 4, nil); err != nil {
			t.Fatalf("setup bucket n=4: %v", err)
		}
		t.Logf("bucket n=4 artifacts ready in %s", time.Since(start))
	}

	// An artifact-less directory must refuse to become an engine.
	if _, err := gnark.New(t.TempDir(), nil); err == nil {
		t.Fatal("empty artifact dir must fail engine construction")
	}

	// The engine discovers the bucket set from disk — nothing declares it.
	engine, err := gnark.New(binDir, nil)
	if err != nil {
		t.Fatalf("load engine from %s: %v", binDir, err)
	}
	if got := engine.Buckets(); len(got) == 0 || got[0] != 4 {
		t.Fatalf("discovery found %v, want it to include 4", got)
	}
	counting := &countingE2EEngine{engine: engine}
	store := memstore.New(0, 0)
	t.Cleanup(store.Close)
	svc, err := prover.New(counting, memqueue.New(0), store, prover.Config{PollInterval: 50 * time.Millisecond})
	if err != nil {
		t.Fatalf("service: %v", err)
	}
	t.Cleanup(svc.Stop)

	// Three real signers in the smallest bucket (4) forces one dummy slot, so
	// padding is part of what gets proven.
	sigs := prover.MakeTestSigs(t, 3)
	var id prover.ProofID

	t.Run("async submit->poll->done", func(t *testing.T) {
		id, err = svc.SubmitProof(sigs)
		if err != nil {
			t.Fatalf("submit: %v", err)
		}
		res := awaitDone(t, svc, id, 5*time.Minute)

		if res.Bucket != 4 || len(res.PaddedSigs) != 4 {
			t.Fatalf("bucket=%d len(padded)=%d, want 4/4", res.Bucket, len(res.PaddedSigs))
		}
		if res.PaddedSigs[3].Active {
			t.Fatal("slot 3 must be dummy padding")
		}
		for i, v := range res.Proof {
			if v == nil || v.Sign() == 0 {
				t.Fatalf("proof[%d] is empty: %v", i, v)
			}
		}
		// The id must bind exactly what was proven under exactly this circuit:
		// prover.ProofID = SHA-256(vkDigest || witnessDigest(paddedSigs)).
		witnessDigest, err := witness.ComputeWitnessHash(res.PaddedSigs)
		if err != nil {
			t.Fatalf("recompute witness hash: %v", err)
		}
		if want := prover.NewProofID(engine.VKDigest(res.Bucket), witnessDigest); id != want {
			t.Fatalf("prover.ProofID does not bind vk+witness: got %s want %s", id, want)
		}
		if got := counting.proves.Load(); got != 1 {
			t.Fatalf("engine proved %d times, want 1", got)
		}
	})

	t.Run("duplicate submit is free", func(t *testing.T) {
		id2, err := svc.SubmitProof(sigs)
		if err != nil {
			t.Fatalf("resubmit: %v", err)
		}
		if id2 != id {
			t.Fatalf("resubmit changed the id: %s vs %s", id2, id)
		}
		job, err := svc.GetProof(id2)
		if err != nil || job.Status != prover.StatusDone {
			t.Fatalf("resubmitted job not immediately done: %+v err=%v", job, err)
		}
		if got := counting.proves.Load(); got != 1 {
			t.Fatalf("duplicate submit re-proved: engine ran %d times", got)
		}
	})

	t.Run("WaitProof returns the finished result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		res, err := svc.WaitProof(ctx, id)
		if err != nil || res.Bucket != 4 {
			t.Fatalf("WaitProof: res=%+v err=%v", res, err)
		}
	})

	t.Run("bad input fails fast without touching the engine", func(t *testing.T) {
		before := counting.proves.Load()
		bad := prover.MakeTestSigs(t, 1)
		bad[0].Signature[0] ^= 0xFF
		start := time.Now()
		_, err := svc.SubmitProof(bad)
		var invalid *prover.ErrInvalidInput
		if !errors.As(err, &invalid) {
			t.Fatalf("want prover.ErrInvalidInput, got %v", err)
		}
		if d := time.Since(start); d > time.Second {
			t.Fatalf("validation took %s — must fail fast, not prove", d)
		}
		if counting.proves.Load() != before {
			t.Fatal("invalid input reached the engine")
		}
	})

}

func awaitDone(t *testing.T, svc *prover.Service, id prover.ProofID, timeout time.Duration) prover.ProofResult {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		job, err := svc.GetProof(id)
		if err != nil {
			t.Fatalf("get proof: %v", err)
		}
		switch job.Status {
		case prover.StatusDone:
			return *job.Result
		case prover.StatusFailed:
			t.Fatalf("job failed (%s): %s", job.Class, job.Cause)
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("job %s did not finish within %s", id, timeout)
	return prover.ProofResult{}
}

// countingE2EEngine wraps the real engine to assert how many real proves ran.
type countingE2EEngine struct {
	engine *gnark.Engine
	proves atomic.Int64
}

func (c *countingE2EEngine) Prove(ctx context.Context, req prover.ProofRequest) (prover.ProofResult, error) {
	c.proves.Add(1)
	return c.engine.Prove(ctx, req)
}
func (c *countingE2EEngine) Buckets() []int               { return c.engine.Buckets() }
func (c *countingE2EEngine) VKDigest(bucket int) [32]byte { return c.engine.VKDigest(bucket) }
