package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	var backendName string
	var gpuProve bool
	flag.StringVar(&backendName, "prover-backend", "", "proof backend to use for setup/proving (native or icicle)")
	flag.BoolVar(&gpuProve, "gpu-prove", false, "shorthand for -prover-backend=icicle")
	flag.Parse()

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
	proofBackend, err := prover.NewProofBackendFromEnv()
	if backendName != "" || gpuProve {
		proofBackend, err = prover.NewProofBackendFromSelection(backendName, gpuProve)
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("Using proof backend: %s\n", proofBackend.Name())

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
		pk, vk, err := proofBackend.Setup(r1csObj)
		if err != nil {
			panic(fmt.Errorf("setup n=%d: %w", n, err))
		}

		smokeTest(r1csObj, pk, vk, n, proofBackend)

		solPath := filepath.Join(solOutDir, fmt.Sprintf("Groth16Verifier_N%d.sol", n))
		if err := exportSolidityVerifier(vk, solPath, n); err != nil {
			panic(fmt.Errorf("export solidity n=%d: %w", n, err))
		}
		fmt.Printf("Saved %s\n", solPath)

		mustWriteArtifact(filepath.Join(bucketDir, "r1cs.bin"), r1csObj)
		mustWriteArtifact(filepath.Join(bucketDir, "pk.bin"), pk)
		mustWriteArtifact(filepath.Join(bucketDir, "vk.bin"), vk)
	}

	fmt.Println("\nAll buckets built.")
}

// smokeTest proves and verifies a fresh batch of n real Ed25519 signatures
// over random per-slot message bytes (truncated to a typical canonical-vote
// size). It exercises the full Sig/Pub decompression + msg-bytes hashing
// path so a broken circuit fails fast before the operator loads the
// artifacts into the relayer.
func smokeTest(cs constraint.ConstraintSystem, pk groth16.ProvingKey, vk groth16.VerifyingKey, n int, proofBackend prover.ProofBackend) {
	// Use a per-slot length comparable to a Tendermint canonical vote
	// (~110-175 bytes). Random bytes are fine for the circuit — semantics
	// don't matter, only that what's signed matches what's hashed.
	const smokeMsgLen = 113

	valSigs := make([]prover.ValidatorSignature, n)
	for i := 0; i < n; i++ {
		pub, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			panic(err)
		}
		msg := make([]byte, smokeMsgLen)
		if _, err := rand.Read(msg); err != nil {
			panic(err)
		}
		sig := ed25519.Sign(priv, msg)
		valSigs[i] = prover.ValidatorSignature{
			Signature:   sig,
			PublicKey:   pub,
			SignedBytes: msg,
			Active:      true,
		}
	}

	hash, err := prover.ComputeWitnessHash(valSigs)
	if err != nil {
		panic(fmt.Errorf("hash witness n=%d: %w", n, err))
	}

	assignment, err := buildSmokeAssignment(valSigs, hash)
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
	proof, err := proofBackend.Prove(cs, pk, witness)
	if err != nil {
		panic(fmt.Errorf("prove n=%d: %w", n, err))
	}
	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		panic(fmt.Errorf("verify n=%d: %w", n, err))
	}
	fmt.Printf("Smoke test passed for n=%d\n", n)
}

func buildSmokeAssignment(
	sigs []prover.ValidatorSignature,
	hash [32]byte,
) (*prover.BatchCircuit[prover.Fp25519, prover.Fr25519], error) {
	n := len(sigs)
	a := &prover.BatchCircuit[prover.Fp25519, prover.Fr25519]{
		Sig:     make([]eddsa.Signature[prover.Fp25519, prover.Fr25519], n),
		Pub:     make([]eddsa.PublicKey[prover.Fp25519, prover.Fr25519], n),
		Msgs:    make([][prover.MaxMsgLen]uints.U8, n),
		MsgLens: make([]frontend.Variable, n),
		Active:  make([]frontend.Variable, n),
	}
	publicInputs := prover.DigestPublicInputs(hash)
	for i := range publicInputs {
		a.Hash[i] = publicInputs[i]
	}
	for i := 0; i < n; i++ {
		v := sigs[i]
		R := v.Signature[:32]
		S, err := edwards25519.NewScalar().SetCanonicalBytes(v.Signature[32:])
		if err != nil {
			return nil, err
		}
		aX, aY, err := utils.DecompressPoint(v.PublicKey)
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
		for j := 0; j < prover.MaxMsgLen; j++ {
			if j < len(v.SignedBytes) {
				a.Msgs[i][j] = uints.NewU8(v.SignedBytes[j])
			} else {
				a.Msgs[i][j] = uints.NewU8(0)
			}
		}
		a.MsgLens[i] = len(v.SignedBytes)
		if v.Active {
			a.Active[i] = 1
		} else {
			a.Active[i] = 0
		}
	}
	return a, nil
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
