package prover

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sync"
	"testing"
)

// makeSyntheticSigs builds `count` valid Ed25519 signatures over per-slot
// messages. The circuit verifies each (pubkey, message, signature) tuple
// independently, so synthetic keypairs are sufficient to benchmark proving —
// no Cosmos commit, validator set, or relayer is involved. Each slot uses a
// distinct message so the (R, A) points across the batch are distinct.
func makeSyntheticSigs(tb testing.TB, count int) []ValidatorSignature {
	tb.Helper()
	sigs := make([]ValidatorSignature, count)
	for i := range count {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			tb.Fatalf("ed25519 keygen: %v", err)
		}
		// Distinct, realistic-width (<= MaxMsgLen) message per slot.
		msg := make([]byte, 110)
		for j := range msg {
			msg[j] = byte(j + i)
		}
		sigs[i] = ValidatorSignature{
			Signature:   ed25519.Sign(priv, msg),
			PublicKey:   pub,
			Index:       i,
			Power:       1,
			SignedBytes: msg,
			Active:      true,
		}
	}
	return sigs
}

// BenchmarkGenerateProof measures end-to-end proof generation
// (witness build + groth16.Prove + local verify) for each bucket that has
// artifacts on disk, in isolation from the relayer / Cosmos / Ethereum.
//
// Each bucket is loaded individually (not via NewProver, which is all-or-nothing
// over every configured bucket); buckets without artifacts are skipped, so you
// can benchmark {4,8,16} without compiling {32,64,128}.
//
// Proving takes seconds, so use a fixed iteration count and (for mean ± stddev
// like a Criterion report) several runs aggregated with benchstat:
//
//	# artifacts live in $PROVER_BIN_DIR/n{N} (generate: go run ./prover/cmd ./bin ../contracts/verifiers)
//	# IMPORTANT: `go test` runs with CWD = the package dir, so PROVER_BIN_DIR must
//	# be ABSOLUTE (a relative ./bin would resolve against relayer/prover, not relayer).
//	# CPU (run from relayer/):
//	PROVER_BIN_DIR="$PWD/bin" go test -bench=GenerateProof -benchtime=10x -count=6 ./prover | tee cpu.txt
//	# GPU (on a CUDA host, e.g. AWS g7e.2xlarge): run ONE bucket per process for a
//	# clean device context — the ICICLE backend reclaims proving-key / MSM device
//	# memory via finalizers, so proving several buckets back-to-back in one process
//	# can pile up on the GPU and fault a later, larger bucket's MSM.
//	for n in 4 8 16 32 64 128; do
//	  PROVER_BIN_DIR="$PWD/bin" go test -tags=icicle -bench="GenerateProof/bucket=$n" -benchtime=10x -count=6 ./prover | tee -a gpu.txt
//	done
//	benchstat cpu.txt          # mean ± variation per bucket
//	# add -v to see which buckets were skipped for missing artifacts
//
// ns/op is dominated by groth16.Prove (witness build and local verify are
// sub-10ms); divide by 1e9 for seconds-per-proof.
func BenchmarkGenerateProof(b *testing.B) {
	binDir := os.Getenv("PROVER_BIN_DIR")
	if binDir == "" {
		b.Skip("set PROVER_BIN_DIR to the prover artifacts directory (containing n{N}/{r1cs,pk,vk}.bin) to run this benchmark")
	}

	backend, err := NewProofBackend(GPUProveEnvEnabled())
	if err != nil {
		b.Fatalf("init proof backend: %v", err)
	}

	for _, n := range Buckets {
		name := fmt.Sprintf("bucket=%d", n)

		// Skip buckets without artifacts (e.g. 32/64/128 if not compiled).
		if _, statErr := os.Stat(filepath.Join(binDir, fmt.Sprintf("n%d", n), "pk.bin")); statErr != nil {
			b.Run(name, func(b *testing.B) {
				b.Skipf("no artifacts for bucket n=%d in %s (run ./prover/cmd to generate)", n, binDir)
			})
			continue
		}

		// Load inside b.Run so the (multi-hundred-MB, up to ~3 GB) artifacts are read
		// ONLY for buckets whose sub-benchmark actually runs — a filtered run such as
		// -bench=GenerateProof/bucket=16 must not load n32/n64/n128. The framework
		// invokes the b.Run body more than once (a priming b.N=1 pass, then the
		// measured pass), so a sync.Once caches the load across those passes.
		var (
			art     *bucketArtifacts
			prover  *EcipProver
			sigs    []ValidatorSignature
			loadErr error
			once    sync.Once
		)

		b.Run(name, func(b *testing.B) {
			once.Do(func() {
				art, loadErr = loadBucketArtifacts(binDir, n, backend)
				if loadErr != nil {
					return
				}
				prover = &EcipProver{byBucket: map[int]*bucketArtifacts{n: art}, backend: backend}
				sigs = makeSyntheticSigs(b, n) // fill the bucket (all active) — worst case
			})
			if loadErr != nil {
				b.Fatalf("load bucket n=%d: %v", n, loadErr)
			}

			b.ResetTimer() // exclude the one-time artifact load from ns/op
			for range b.N {
				if _, _, _, _, _, genErr := prover.GenerateProof(sigs); genErr != nil {
					b.Fatalf("GenerateProof(bucket=%d): %v", n, genErr)
				}
			}
		})

		// Release this bucket's proving key + device buffers before the next, larger
		// bucket loads. On GPU (ICICLE) they are reclaimed via finalizers, so without
		// forcing a GC between buckets they accumulate on the device and a later
		// bucket's MSM faults with an illegal memory access. Runs after the
		// sub-benchmark has completed (b.Run is synchronous), outside any measured
		// region. No-op for buckets that were filtered out (nothing loaded). For the
		// strictest isolation, run one bucket per process (see the -bench example).
		art, prover, sigs = nil, nil, nil
		runtime.GC()
		debug.FreeOSMemory()
	}
}
