// Package gnark implements the prover.ProofEngine port with the Groth16/BN254
// batch-Ed25519 circuit. The circuit definition (circuit.go), witness
// assignment (assignment.go), and proof backends (backend*.go) are carried
// verbatim from the audited implementation so the deployed verifying keys
// stay valid; only this engine shell is new. The circuit rewrite is a
// separate, later task.
package gnark

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	gnarkbackend "github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/solidity"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"

	"relayer/prover"
	"relayer/prover/witness"
)

type bucketArtifacts struct {
	r1cs     constraint.ConstraintSystem
	pk       groth16.ProvingKey
	vk       groth16.VerifyingKey
	vkDigest [32]byte
}

// Engine proves against per-bucket artifacts loaded from
// binDir/n{N}/{r1cs,pk,vk}.bin. Immutable after New, safe for concurrent
// Prove calls.
type Engine struct {
	byBucket map[int]*bucketArtifacts
	backend  ProofBackend
}

// New discovers and loads every bucket that has a complete artifact triple
// under binDir (n{N}/{r1cs,pk,vk}.bin). The disk IS the configuration: which
// buckets are usable is decided by what setup-circuits produced, never by a
// code or config constant. Discovering zero buckets is a startup error — the
// operator must run cmd/setup-circuits first. A directory whose N is outside
// witness.SupportedBuckets is an error too (a stray or mistyped artifact dir
// should not silently become a provable bucket). A nil backend means native
// gnark proving.
func New(binDir string, backend ProofBackend) (*Engine, error) {
	if backend == nil {
		backend = nativeProofBackend{}
	}
	buckets, err := discoverBuckets(binDir)
	if err != nil {
		return nil, err
	}
	if len(buckets) == 0 {
		return nil, fmt.Errorf("no complete bucket artifacts (n{N}/{r1cs,pk,vk}.bin) under %s — run cmd/setup-circuits first", binDir)
	}
	e := &Engine{byBucket: make(map[int]*bucketArtifacts, len(buckets)), backend: backend}
	for _, n := range buckets {
		art, err := loadBucketArtifacts(binDir, n, backend)
		if err != nil {
			return nil, fmt.Errorf("load bucket n=%d from %s: %w", n, binDir, err)
		}
		e.byBucket[n] = art
		log.Printf("[prover] bucket n=%d loaded, vk_digest=%x", n, art.vkDigest[:8])
	}
	return e, nil
}

// discoverBuckets scans binDir for n{N} directories carrying the full
// artifact triple and returns their sizes sorted ascending.
func discoverBuckets(binDir string) ([]int, error) {
	entries, err := os.ReadDir(binDir)
	if err != nil {
		return nil, fmt.Errorf("read artifact dir %s: %w", binDir, err)
	}
	var buckets []int
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "n") {
			continue
		}
		n, err := strconv.Atoi(entry.Name()[1:])
		if err != nil {
			continue // not an n{N} dir (e.g. "notes/")
		}
		complete := true
		for _, artifact := range []string{"r1cs.bin", "pk.bin", "vk.bin"} {
			if _, err := os.Stat(filepath.Join(binDir, entry.Name(), artifact)); err != nil {
				complete = false
				break
			}
		}
		if !complete {
			continue
		}
		if !witness.IsSupported(n) {
			return nil, fmt.Errorf("artifact dir %s/n%d is outside the supported bucket family %v", binDir, n, witness.SupportedBuckets)
		}
		buckets = append(buckets, n)
	}
	sort.Ints(buckets)
	return buckets, nil
}

// NewFromEnv builds the proof backend from the GPU_PROVE environment toggle
// and loads the engine. GPU on a build without -tags=icicle errors out —
// never a silent CPU fallback.
func NewFromEnv(binDir string) (*Engine, error) {
	b, err := NewProofBackend(GPUProveEnvEnabled())
	if err != nil {
		return nil, err
	}
	return New(binDir, b)
}

// Buckets declares which bucket sizes artifacts were loaded for, ascending.
func (e *Engine) Buckets() []int {
	out := make([]int, 0, len(e.byBucket))
	for n := range e.byBucket {
		out = append(out, n)
	}
	sort.Ints(out)
	return out
}

// VKDigest returns the SHA-256 of the bucket's verifying key bytes — the
// circuit-version half of every ProofID. Unknown buckets return the zero
// digest; SubmitProof validates the bucket before it ever derives an id.
func (e *Engine) VKDigest(bucket int) [32]byte {
	if art, ok := e.byBucket[bucket]; ok {
		return art.vkDigest
	}
	return [32]byte{}
}

// Prove builds the witness assignment, runs the Groth16 prove, and — always —
// local-verifies against the loaded vk before returning. No unverified proof
// leaves the engine: a verify failure here means artifact/layout drift or a
// soundness bug, and surfaces as an error instead of an on-chain revert.
// ctx gates phases that have not started; gnark itself cannot be interrupted
// mid-prove. Once Prove returns a proof, verification and serialization finish
// even if ctx expired meanwhile, so completed work is not discarded.
func (e *Engine) Prove(ctx context.Context, req prover.ProofRequest) (prover.ProofResult, error) {
	art, ok := e.byBucket[req.Bucket]
	if !ok {
		return prover.ProofResult{}, fmt.Errorf("bucket n=%d not loaded", req.Bucket)
	}
	if len(req.Sigs) != req.Bucket {
		return prover.ProofResult{}, fmt.Errorf("request carries %d slots, bucket is %d", len(req.Sigs), req.Bucket)
	}
	if err := ctx.Err(); err != nil {
		return prover.ProofResult{}, err
	}

	assignment, err := buildBatchAssignment(req.Sigs, req.WitnessDigest)
	if err != nil {
		return prover.ProofResult{}, fmt.Errorf("build assignment: %w", err)
	}
	fullWitness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return prover.ProofResult{}, fmt.Errorf("create witness: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return prover.ProofResult{}, err
	}

	gnarkProof, err := e.backend.Prove(art.r1cs, art.pk, fullWitness,
		solidity.WithProverTargetSolidityVerifier(gnarkbackend.GROTH16))
	if err != nil {
		return prover.ProofResult{}, fmt.Errorf("generate proof: %w", err)
	}

	pubWitness, err := fullWitness.Public()
	if err != nil {
		return prover.ProofResult{}, fmt.Errorf("public witness: %w", err)
	}
	if err := groth16.Verify(gnarkProof, art.vk, pubWitness,
		solidity.WithVerifierTargetSolidityVerifier(gnarkbackend.GROTH16)); err != nil {
		return prover.ProofResult{}, fmt.Errorf("local verification failed: %w", err)
	}

	proofInts, commitments, commitmentPok, err := ProofToBigInts(gnarkProof)
	if err != nil {
		return prover.ProofResult{}, err
	}
	return prover.ProofResult{
		Bucket:        req.Bucket,
		PaddedSigs:    req.Sigs,
		Proof:         proofInts,
		Commitments:   commitments,
		CommitmentPok: commitmentPok,
	}, nil
}

func loadBucketArtifacts(binDir string, n int, b ProofBackend) (*bucketArtifacts, error) {
	dir := filepath.Join(binDir, fmt.Sprintf("n%d", n))

	r1cs := groth16.NewCS(ecc.BN254)
	if err := readFromFile(filepath.Join(dir, "r1cs.bin"), r1cs); err != nil {
		return nil, fmt.Errorf("read r1cs: %w", err)
	}
	pk := b.NewProvingKey(ecc.BN254)
	if err := readFromFile(filepath.Join(dir, "pk.bin"), pk); err != nil {
		return nil, fmt.Errorf("read pk: %w", err)
	}
	vk := groth16.NewVerifyingKey(ecc.BN254)
	if err := readFromFile(filepath.Join(dir, "vk.bin"), vk); err != nil {
		return nil, fmt.Errorf("read vk: %w", err)
	}

	h := sha256.New()
	if _, err := vk.WriteTo(h); err != nil {
		return nil, fmt.Errorf("hash vk: %w", err)
	}
	var digest [32]byte
	copy(digest[:], h.Sum(nil))

	return &bucketArtifacts{r1cs: r1cs, pk: pk, vk: vk, vkDigest: digest}, nil
}

func readFromFile(path string, dst io.ReaderFrom) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = dst.ReadFrom(f)
	return err
}
