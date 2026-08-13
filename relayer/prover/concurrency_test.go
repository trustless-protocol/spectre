package prover

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// TestGenerateProofIsSafeUnderConcurrency closes the last RLY-11 sub-item:
// "concurrent GenerateProof asserted safe but untested".
//
// One EcipProver is shared by every relay path — both directions' flushes, the
// refresh routine and both timeout scanners each run on their own goroutine and
// call GenerateProof against the same *bucketArtifacts (r1cs, proving key). The
// safety claim is that those artifacts are read-only after load, so no call
// mutates state another call can observe. That was an assertion, not a result.
//
// This runs real proofs concurrently under -race, which is the only thing that
// can actually falsify it: a write to shared prover state shows up as a data
// race, and a witness or key corrupted by interleaving fails the in-proof local
// verify inside GenerateProof.
//
// Requires bucket-4 artifacts. Run:
//
//	PROVER_BIN_DIR="$PWD/bin" go test -race -run TestGenerateProofIsSafeUnderConcurrency ./prover
func TestGenerateProofIsSafeUnderConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real proof generation in -short mode")
	}

	binDir := os.Getenv("PROVER_BIN_DIR")
	if binDir == "" {
		t.Skip("set PROVER_BIN_DIR to the prover artifacts directory (containing n{N}/{r1cs,pk,vk}.bin)")
	}
	const bucket = 4
	// All three, not just pk.bin: GenerateProof verifies locally before returning,
	// so a checkout carrying a partial bucket (r1cs+pk but no vk) must SKIP rather
	// than fail on a missing file and look like a concurrency defect.
	for _, name := range []string{"r1cs.bin", "pk.bin", "vk.bin"} {
		path := filepath.Join(binDir, fmt.Sprintf("n%d", bucket), name)
		if _, err := os.Stat(path); err != nil {
			t.Skipf("bucket n=%d is missing %s in %s (run `go run ./prover/cmd ./bin ../contracts/verifiers` to generate)", bucket, name, binDir)
		}
	}

	backend, err := NewProofBackend(GPUProveEnvEnabled())
	if err != nil {
		t.Fatalf("init proof backend: %v", err)
	}
	art, err := loadBucketArtifacts(binDir, bucket, backend)
	if err != nil {
		t.Fatalf("load bucket n=%d: %v", bucket, err)
	}
	p := &EcipProver{byBucket: map[int]*bucketArtifacts{bucket: art}, backend: backend}

	// Distinct signature sets per goroutine: identical inputs could mask a shared
	// witness buffer by making every writer store the same bytes.
	const goroutines = 4
	sigSets := make([][]ValidatorSignature, goroutines)
	for i := range sigSets {
		sigSets[i] = makeSyntheticSigs(t, bucket)
	}

	var wg sync.WaitGroup
	errs := make([]error, goroutines)
	start := make(chan struct{})
	for i := range goroutines {
		wg.Go(func() {
			<-start // widen the overlap window rather than letting them stagger
			gotBucket, padded, _, _, _, err := p.GenerateProof(sigSets[i])
			switch {
			case err != nil:
				errs[i] = err
			case gotBucket != bucket:
				errs[i] = fmt.Errorf("bucket = %d, want %d", gotBucket, bucket)
			case len(padded) != bucket:
				errs[i] = fmt.Errorf("padded signature count = %d, want %d", len(padded), bucket)
			}
		})
	}
	close(start)
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d: %v", i, err)
		}
	}
}
