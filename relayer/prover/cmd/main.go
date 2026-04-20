package main

import (
	"crypto/ed25519"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"0x5ea000000/ecip-gnark/signature/canonvote"
	"0x5ea000000/ecip-gnark/signature/eddsa"
	"0x5ea000000/ecip-gnark/utils"

	"filippo.io/edwards25519"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/uints"

	"relayer/prover"
)

// setup-circuits: for each bucket in prover.Buckets, compile BatchCircuit,
// run Groth16 setup, write bin/n{N}/{r1cs,pk,vk}.bin, then export the
// corresponding Solidity verifier to contracts/verifiers/Groth16Verifier_N{N}.sol.
//
// Also runs a smoke prove+verify with real Ed25519 signatures so a broken
// circuit fails fast before the operator tries to use it.
func main() {
	outDir := "bin"
	solOutDir := filepath.Join("..", "contracts", "verifiers")
	if len(os.Args) > 1 {
		outDir = os.Args[1]
	}
	if len(os.Args) > 2 {
		solOutDir = os.Args[2]
	}
	if err := os.MkdirAll(solOutDir, 0o755); err != nil {
		panic(fmt.Errorf("create solidity output dir: %w", err))
	}

	for _, n := range prover.Buckets {
		fmt.Printf("\n=== bucket n=%d ===\n", n)
		bucketDir := filepath.Join(outDir, fmt.Sprintf("n%d", n))
		if err := os.MkdirAll(bucketDir, 0o755); err != nil {
			panic(fmt.Errorf("create %s: %w", bucketDir, err))
		}

		circuit := prover.NewBatchCircuit(n)
		fmt.Printf("Compiling BatchCircuit (n=%d)...\n", n)
		r1csObj, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
		if err != nil {
			panic(fmt.Errorf("compile n=%d: %w", n, err))
		}
		fmt.Printf("R1CS: %d constraints, %d public witness values\n",
			r1csObj.GetNbConstraints(), r1csObj.GetNbPublicVariables())

		fmt.Println("Running Groth16 setup...")
		pk, vk, err := groth16.Setup(r1csObj)
		if err != nil {
			panic(fmt.Errorf("setup n=%d: %w", n, err))
		}

		smokeTest(r1csObj, pk, vk, n)

		solPath := filepath.Join(solOutDir, fmt.Sprintf("Groth16Verifier_N%d.sol", n))
		solFile, err := os.Create(solPath)
		if err != nil {
			panic(err)
		}
		if err := vk.ExportSolidity(solFile); err != nil {
			solFile.Close()
			panic(fmt.Errorf("export solidity n=%d: %w", n, err))
		}
		solFile.Close()
		fmt.Printf("Saved %s\n", solPath)

		mustWriteArtifact(filepath.Join(bucketDir, "r1cs.bin"), r1csObj)
		mustWriteArtifact(filepath.Join(bucketDir, "pk.bin"), pk)
		mustWriteArtifact(filepath.Join(bucketDir, "vk.bin"), vk)
	}

	fmt.Println("\nAll buckets built.")
}

// smokeTest proves and verifies a fresh batch of n real Ed25519 signatures
// over a synthetic CanonicalVote. It catches circuit-build regressions before
// the operator ever loads the artifacts into the relayer.
func smokeTest(cs constraint.ConstraintSystem, pk groth16.ProvingKey, vk groth16.VerifyingKey, n int) {
	chainID := "smoke-test-chain"
	height := int64(1)
	round := int64(0)
	blockIDHash := bytesWithSeed(0x11, 32)
	partSetTotal := uint32(1)
	partSetHash := bytesWithSeed(0x22, 32)

	sigs := make([][]byte, n)
	pubs := make([][]byte, n)
	tsSec := make([]int64, n)
	tsNanos := make([]int32, n)
	for i := 0; i < n; i++ {
		pub, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			panic(err)
		}
		tsSec[i] = int64(1700000000 + i)
		tsNanos[i] = int32(i)
		voteBytes := canonvote.ReferenceEncodeVote(
			height, round, blockIDHash, partSetTotal, partSetHash,
			tsSec[i], tsNanos[i], chainID,
		)
		sigs[i] = ed25519.Sign(priv, voteBytes)
		pubs[i] = pub
	}

	assignment, err := buildSmokeAssignment(
		height, round, blockIDHash, partSetTotal, partSetHash, chainID,
		sigs, pubs, tsSec, tsNanos,
	)
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
	proof, err := groth16.Prove(cs, pk, witness)
	if err != nil {
		panic(fmt.Errorf("prove n=%d: %w", n, err))
	}
	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		panic(fmt.Errorf("verify n=%d: %w", n, err))
	}
	fmt.Printf("Smoke test passed for n=%d\n", n)
}

func buildSmokeAssignment(
	height, round int64,
	blockIDHash []byte, partSetTotal uint32, partSetHash []byte,
	chainID string,
	sigs, pubs [][]byte, tsSec []int64, tsNanos []int32,
) (*prover.BatchCircuit[prover.Fp25519, prover.Fr25519], error) {
	n := len(sigs)
	a := &prover.BatchCircuit[prover.Fp25519, prover.Fr25519]{
		Sig:       make([]eddsa.Signature[prover.Fp25519, prover.Fr25519], n),
		Pub:       make([]eddsa.PublicKey[prover.Fp25519, prover.Fr25519], n),
		TsSeconds: make([]frontend.Variable, n),
		TsNanos:   make([]frontend.Variable, n),
	}
	for i := 0; i < n; i++ {
		sig, pub := sigs[i], pubs[i]
		R := sig[:32]
		S, err := edwards25519.NewScalar().SetCanonicalBytes(sig[32:])
		if err != nil {
			return nil, err
		}
		aX, aY, err := utils.DecompressPoint(pub)
		if err != nil {
			return nil, err
		}
		rX, rY, err := utils.DecompressPoint(R)
		if err != nil {
			return nil, err
		}
		a.Sig[i] = eddsa.Signature[prover.Fp25519, prover.Fr25519]{
			R: sw_emulated.AffinePoint[prover.Fp25519]{
				X: emulated.ValueOf[prover.Fp25519](rX),
				Y: emulated.ValueOf[prover.Fp25519](rY),
			},
			S: emulated.ValueOf[prover.Fr25519](utils.ScalarToBigInt(S)),
		}
		a.Pub[i] = eddsa.PublicKey[prover.Fp25519, prover.Fr25519]{
			A: sw_emulated.AffinePoint[prover.Fp25519]{
				X: emulated.ValueOf[prover.Fp25519](aX),
				Y: emulated.ValueOf[prover.Fp25519](aY),
			},
		}
		a.TsSeconds[i] = tsSec[i]
		a.TsNanos[i] = tsNanos[i]
	}
	for i := 0; i < 32; i++ {
		a.BlockIDHash[i] = uints.NewU8(blockIDHash[i])
		a.PartSetHash[i] = uints.NewU8(partSetHash[i])
	}
	for i := 0; i < canonvote.MaxChainIDLen; i++ {
		if i < len(chainID) {
			a.ChainID[i] = uints.NewU8(chainID[i])
		} else {
			a.ChainID[i] = uints.NewU8(0)
		}
	}
	a.Height = height
	a.Round = round
	a.PartSetTotal = partSetTotal
	a.ChainIDLen = len(chainID)
	return a, nil
}

func bytesWithSeed(seed byte, n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = seed + byte(i)
	}
	return b
}

func mustWriteArtifact(path string, value interface {
	WriteTo(io.Writer) (int64, error)
}) {
	f, err := os.Create(path)
	if err != nil {
		panic(fmt.Errorf("create %s: %w", path, err))
	}
	defer f.Close()
	if _, err := value.WriteTo(f); err != nil {
		panic(fmt.Errorf("write %s: %w", path, err))
	}
	fmt.Printf("Saved %s\n", path)
}
