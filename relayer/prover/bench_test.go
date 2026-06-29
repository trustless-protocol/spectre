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

// makeSyntheticVote builds one #199-compliant canonical-vote-shaped message:
// a 1-byte length varint, the Type|Height prefix head, and a 32-byte block hash
// at the no-round offset (msg[16:48]). Mirrors makeSmokeVote in prover/cmd. The
// exact field values are irrelevant — the circuit binds only the prefix head and
// block hash, which every active signer must share.
func makeSyntheticVote() []byte {
	const voteLen = 120 // bodyLen 119 (<128) -> 1-byte leading varint
	msg := make([]byte, voteLen)
	msg[0] = byte(voteLen - 1)  // 1-byte length varint (high bit clear)
	msg[1], msg[2] = 0x08, 0x02 // Type = precommit (field 1)
	msg[3] = 0x11               // Height tag (field 2, sfixed64)
	for j := 0; j < 8; j++ {
		msg[4+j] = byte(j + 1)
	}
	// BlockID (field 4): tag 0x22, len 0x48, inner hash tag 0x0a, len 0x20.
	// msg[12] != round tag 0x19 -> roundPresent false -> hash at msg[16:48].
	msg[12], msg[13], msg[14], msg[15] = 0x22, 0x48, 0x0a, 0x20
	for j := 0; j < 32; j++ {
		msg[16+j] = byte(0xA0 + j) // block hash
	}
	for j := 48; j < voteLen; j++ {
		msg[j] = byte(j) // part-set header / timestamp / chain id (unconstrained)
	}
	return msg
}

// makeSyntheticSigs builds `count` valid Ed25519 signatures over a single
// #199-compliant canonical vote. The #199 circuit binds every ACTIVE slot's
// message to a common Type|Height prefix + block hash, so all signers must sign
// the SAME vote — distinct keypairs still yield distinct (R, A) points, which is
// what the ECIP aggregate requires. Proving cost is identical whether the
// per-slot messages are equal or not (the circuit shape is fixed per bucket), so
// this is a faithful proving benchmark and matches production, where validators
// sign votes sharing the same prefix + block hash.
func makeSyntheticSigs(tb testing.TB, count int) []ValidatorSignature {
	tb.Helper()
	vote := makeSyntheticVote()
	sigs := make([]ValidatorSignature, count)
	for i := range count {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			tb.Fatalf("ed25519 keygen: %v", err)
		}
		sigs[i] = ValidatorSignature{
			Signature:   ed25519.Sign(priv, vote),
			PublicKey:   pub,
			Index:       i,
			Power:       1,
			SignedBytes: vote,
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
