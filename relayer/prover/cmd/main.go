package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
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

	"relayer/prover/gnark"
	"relayer/prover/witness"
)

// setup-circuits: for each bucket in proof.Buckets, compile BatchCircuit,
// run Groth16 setup, write bin/n{N}/{r1cs,pk,vk}.bin, then export the
// corresponding Solidity verifier to contracts/verifiers/Groth16Verifier_N{N}.sol.
//
// Also runs a smoke prove+verify with real Ed25519 signatures so a broken
// circuit fails fast before the operator tries to use it. With -onchain-fixture-dir it
// also writes that proof, in the words the on-chain verifier call takes, for the Foundry
// check that SignatureVerifier.verifyBatchProof accepts a proof of this circuit.
func main() {
	var gpuProve bool
	var listBuckets bool
	var bucketsFlag, fixtureDir string
	flag.BoolVar(&gpuProve, "gpu-prove", false, "use the ICICLE GPU backend for setup/proving (or set GPU_PROVE=1); requires an icicle-enabled build")
	flag.BoolVar(&listBuckets, "list-buckets", false, "print the buckets this invocation would build and exit")
	flag.StringVar(&fixtureDir, "onchain-fixture-dir", "", "write the smoke proof of each bucket to <dir>/onchain_fixture_n{N}.json for the Foundry on-chain check")
	flag.StringVar(&bucketsFlag, "buckets", "4", "comma-separated bucket sizes to build; each must be in the supported family "+fmt.Sprint(witness.SupportedBuckets))
	flag.Parse()
	buckets, err := parseBuckets(bucketsFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if listBuckets {
		for _, bucket := range buckets {
			fmt.Println(bucket)
		}
		return
	}

	outDir := "bin"
	solOutDir := filepath.Join("..", "contracts", "verifiers")
	if args := flag.Args(); len(args) > 0 {
		outDir = args[0]
	}
	if args := flag.Args(); len(args) > 1 {
		solOutDir = args[1]
	}
	if err := os.MkdirAll(solOutDir, 0o755); err != nil {
		panic(fmt.Errorf("create solidity output dir: %w", err))
	}

	useGPU := gpuProve || gnark.GPUProveEnvEnabled()
	proofBackend, err := gnark.NewProofBackend(useGPU)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Using proof backend: %s\n", proofBackend.Name())

	for _, n := range buckets {
		fmt.Printf("\n=== bucket n=%d ===\n", n)
		fmt.Printf("Compiling BatchCircuit (n=%d) and running Groth16 setup...\n", n)
		r1csObj, pk, vk, err := gnark.SetupBucket(outDir, n, proofBackend)
		if err != nil {
			panic(err)
		}
		fmt.Printf("R1CS: %d constraints, %d public witness values\n",
			r1csObj.GetNbConstraints(), r1csObj.GetNbPublicVariables())

		fixture := smokeTest(r1csObj, pk, vk, n, proofBackend)

		solPath := filepath.Join(solOutDir, fmt.Sprintf("Groth16Verifier_N%d.sol", n))
		if err := exportSolidityVerifier(vk, solPath, n); err != nil {
			panic(fmt.Errorf("export solidity n=%d: %w", n, err))
		}
		fmt.Printf("Saved %s\n", solPath)
		fmt.Printf("Saved artifacts under %s/n%d\n", outDir, n)
		if fixtureDir != "" {
			if err := writeOnchainFixture(fixtureDir, solPath, n, fixture); err != nil {
				panic(fmt.Errorf("write on-chain fixture n=%d: %w", n, err))
			}
		}
	}

	fmt.Println("\nAll buckets built.")
}

// smokeTest proves and verifies a fresh batch of n real Ed25519 signatures over
// a synthetic canonical vote that satisfies the #199 layout: a 2-byte leading
// length varint, the Type|Height prefix head, and a 32-byte block hash at the
// no-round offset. Every signer signs the SAME vote bytes so they share the
// committed common prefix the circuit binds each active slot against. This
// exercises the full Sig/Pub decompression + prefix/block-hash binding + witness
// hashing path so a broken circuit fails fast before the operator loads the
// artifacts into the relayer.
func smokeTest(cs constraint.ConstraintSystem, pk groth16.ProvingKey, vk groth16.VerifyingKey, n int, proofBackend gnark.ProofBackend) onchainFixture {
	msg := makeSmokeVote()

	valSigs := make([]witness.ValidatorSignature, n)
	for i := 0; i < n; i++ {
		pub, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			panic(err)
		}
		sig := ed25519.Sign(priv, msg)
		valSigs[i] = witness.ValidatorSignature{
			Signature:   sig,
			PublicKey:   pub,
			SignedBytes: msg,
			Active:      true,
		}
	}

	hash, err := witness.ComputeWitnessHash(valSigs)
	if err != nil {
		panic(fmt.Errorf("hash witness n=%d: %w", n, err))
	}

	prefixHead, round, blockHash, err := witness.CommonPrefixFields(valSigs)
	if err != nil {
		panic(fmt.Errorf("common prefix n=%d: %w", n, err))
	}

	assignment, err := gnark.BuildBatchAssignment(valSigs, hash)
	if err != nil {
		panic(fmt.Errorf("build assignment n=%d: %w", n, err))
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		panic(fmt.Errorf("witness n=%d: %w", n, err))
	}
	publicWitness, err := witness.Public()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Proving smoke test (n=%d)...\n", n)
	// The same options the engine proves and verifies with (gnark/engine.go): they select the
	// hash the Solidity verifier recomputes for the commitments. A proof made without them can
	// pass groth16.Verify here and still fail on-chain.
	proof, err := proofBackend.Prove(cs, pk, witness, solidity.WithProverTargetSolidityVerifier(gnarkbackend.GROTH16))
	if err != nil {
		panic(fmt.Errorf("prove n=%d: %w", n, err))
	}
	if err := groth16.Verify(proof, vk, publicWitness, solidity.WithVerifierTargetSolidityVerifier(gnarkbackend.GROTH16)); err != nil {
		panic(fmt.Errorf("verify n=%d: %w", n, err))
	}
	fmt.Printf("Smoke test passed for n=%d\n", n)

	proofWords, commitments, commitmentPok, err := gnark.ProofToBigInts(proof)
	if err != nil {
		panic(fmt.Errorf("proof words n=%d: %w", n, err))
	}
	fixture := onchainFixture{
		Bucket:        n,
		Proof:         hexWords(proofWords[:]),
		Commitments:   hexWords(commitments[:]),
		CommitmentPok: hexWords(commitmentPok[:]),
		PrefixHead:    "0x" + hex.EncodeToString(prefixHead[:]),
		Round:         "0x" + hex.EncodeToString(round[:]),
		BlockHash:     "0x" + hex.EncodeToString(blockHash[:]),
	}
	for _, s := range valSigs {
		fixture.Pubkeys = append(fixture.Pubkeys, "0x"+hex.EncodeToString(s.PublicKey))
		fixture.Active = append(fixture.Active, s.Active)
	}
	return fixture
}

// onchainFixture is a smoke proof laid out as SignatureVerifier.verifyBatchProof takes
// it: the proof words, the two commitments as (x, y) pairs, the folded proof of
// knowledge, the per-slot public keys and active flags, and the common vote fields the
// witness hash commits to. VerifierSHA256 ties it to the exact Groth16Verifier_N{N}.sol
// exported next to it, so the Foundry test can tell a fixture of another setup from a
// broken proof.
type onchainFixture struct {
	Bucket         int      `json:"bucket"`
	VerifierSHA256 string   `json:"verifierSha256"`
	Proof          []string `json:"proof"`
	Commitments    []string `json:"commitments"`
	CommitmentPok  []string `json:"commitmentPok"`
	Pubkeys        []string `json:"pubkeys"`
	Active         []bool   `json:"active"`
	PrefixHead     string   `json:"prefixHead"`
	Round          string   `json:"round"`
	BlockHash      string   `json:"blockHash"`
}

func hexWords(words []*big.Int) []string {
	out := make([]string, len(words))
	for i, w := range words {
		out[i] = fmt.Sprintf("0x%064x", w)
	}
	return out
}

// writeOnchainFixture stores fixture as <dir>/onchain_fixture_n{N}.json, stamped with
// the SHA-256 of the exported verifier at solPath.
func writeOnchainFixture(dir, solPath string, n int, fixture onchainFixture) error {
	sol, err := os.ReadFile(solPath)
	if err != nil {
		return fmt.Errorf("read exported verifier: %w", err)
	}
	sum := sha256.Sum256(sol)
	fixture.VerifierSHA256 = "0x" + hex.EncodeToString(sum[:])
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(fixture, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(dir, fmt.Sprintf("onchain_fixture_n%d.json", n))
	if err := os.WriteFile(path, append(raw, '\n'), 0o644); err != nil {
		return err
	}
	fmt.Printf("Saved %s\n", path)
	return nil
}

// makeSmokeVote builds a synthetic canonical-vote byte slice satisfying the
// #199 witness layout with a 1-BYTE leading length varint (bodyLen < 128) — the
// case real devnets hit with short chain ids that the original 2-byte-only
// assumption rejected. Layout: 1-byte varint (sb[0] < 0x80), the Type|Height
// prefix head at body offset 0, no round field (so the 32-byte block hash sits
// at no-round body offset 15 → in-message index 16), and arbitrary trailing
// bytes. The exact field values do not matter — the circuit binds only the
// prefix head and block hash against what every active slot signed, and here
// every slot signs this same slice.
func makeSmokeVote() []byte {
	const voteLen = 120 // bodyLen = 119 (< 128) → 1-byte leading varint
	msg := make([]byte, voteLen)
	// 1-byte leading length varint (high bit clear). The encoded value is
	// irrelevant to the circuit; only the varint width drives the body offset.
	msg[0] = byte(voteLen - 1)
	// Type = precommit: field 1 (0x08), value 2.
	msg[1] = 0x08
	msg[2] = 0x02
	// Height: field 2 sfixed64 (0x11) + 8 bytes.
	msg[3] = 0x11
	for j := 0; j < 8; j++ {
		msg[4+j] = byte(j + 1)
	}
	// BlockID: field 4 (0x22), length 0x48, inner hash tag 0x0a length 0x20.
	// msg[12] must not be the round tag 0x19 so the round is absent (committed Round = 0), which
	// places the 32-byte block hash at no-round in-message index 16 (msg[16:48]).
	msg[12] = 0x22
	msg[13] = 0x48
	msg[14] = 0x0a
	msg[15] = 0x20
	for j := 0; j < 32; j++ {
		msg[16+j] = byte(0xA0 + j)
	}
	// Remaining bytes (part-set header, timestamp, chain id) are unconstrained.
	for j := 48; j < voteLen; j++ {
		msg[j] = byte(j)
	}
	return msg
}

// exportSolidityVerifier writes gnark's generated verifier and renames the
// default `contract Verifier` to `contract Groth16Verifier_N{N}` so every
// bucket can live in one Solidity import graph.
func exportSolidityVerifier(vk groth16.VerifyingKey, path string, n int) error {
	var buf bytes.Buffer
	if err := vk.ExportSolidity(&buf); err != nil {
		return err
	}
	renamed := bytes.Replace(
		buf.Bytes(),
		[]byte("contract Verifier {"),
		[]byte(fmt.Sprintf("contract Groth16Verifier_N%d {", n)),
		1,
	)
	return os.WriteFile(path, renamed, 0o644)
}

// parseBuckets parses the -buckets flag and rejects sizes outside the
// protocol's supported family (each member needs its own circuit and on-chain
// verifier — an arbitrary N is not buildable).
func parseBuckets(raw string) ([]int, error) {
	var out []int
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid bucket %q in -buckets", part)
		}
		if !witness.IsSupported(n) {
			return nil, fmt.Errorf("bucket %d is outside the supported family %v", n, witness.SupportedBuckets)
		}
		out = append(out, n)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("-buckets is empty")
	}
	sort.Ints(out)
	return out, nil
}
