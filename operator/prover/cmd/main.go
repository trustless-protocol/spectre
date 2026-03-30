package main

import (
	"crypto/ed25519"
	"crypto/sha512"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"0x5ea000000/ecip-gnark/signature/eddsa"
	"0x5ea000000/ecip-gnark/utils"

	"filippo.io/edwards25519"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/emulated"

	"operator/prover"
)

func main() {
	outDir := "bin"
	if len(os.Args) > 1 {
		outDir = os.Args[1]
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		panic(fmt.Errorf("create output dir: %w", err))
	}

	msg := []byte("setup-test")

	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	sig := ed25519.Sign(priv, msg)
	R := sig[:32]
	S, err := edwards25519.NewScalar().SetCanonicalBytes(sig[32:])
	if err != nil {
		panic("invalid signature scalar")
	}

	hasher := sha512.New()
	hasher.Write(R)
	hasher.Write([]byte(pub))
	hasher.Write(msg)
	sum := hasher.Sum(nil)
	H, err := edwards25519.NewScalar().SetUniformBytes(sum)
	if err != nil {
		panic("setting scalar failed")
	}

	if !ed25519.Verify(pub, msg, sig) {
		panic("failed to verify signature outside circuit")
	}

	aX, aY, _ := utils.DecompressPoint([]byte(pub))
	rX, rY, _ := utils.DecompressPoint(R)
	h := utils.ScalarToBigInt(H)
	s := utils.ScalarToBigInt(S)

	// Compile circuit
	fmt.Println("Compiling PreHashCircuit...")
	var circuit prover.PreHashCircuit[prover.Fp25519, prover.Fr25519]
	r1csObj, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(fmt.Errorf("compile: %w", err))
	}
	fmt.Printf("R1CS: %d constraints\n", r1csObj.GetNbConstraints())

	// Groth16 Setup
	fmt.Println("Running Groth16 setup...")
	pk, vk, err := groth16.Setup(r1csObj)
	if err != nil {
		panic(fmt.Errorf("setup: %w", err))
	}

	// Build assignment and prove
	var assignment prover.PreHashCircuit[prover.Fp25519, prover.Fr25519]
	assignment.Sig = eddsa.Signature[prover.Fp25519, prover.Fr25519]{
		R: sw_emulated.AffinePoint[prover.Fp25519]{
			X: emulated.ValueOf[prover.Fp25519](rX),
			Y: emulated.ValueOf[prover.Fp25519](rY),
		},
		S: emulated.ValueOf[prover.Fr25519](s),
	}
	assignment.Hash = emulated.ValueOf[prover.Fr25519](h)
	assignment.Pub = eddsa.PublicKey[prover.Fp25519, prover.Fr25519]{
		A: sw_emulated.AffinePoint[prover.Fp25519]{
			X: emulated.ValueOf[prover.Fp25519](aX),
			Y: emulated.ValueOf[prover.Fp25519](aY),
		},
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		panic(fmt.Errorf("witness: %w", err))
	}
	publicWitness, err := witness.Public()
	if err != nil {
		panic(fmt.Errorf("public witness: %w", err))
	}

	fmt.Println("Generating test proof...")
	proof, err := groth16.Prove(r1csObj, pk, witness)
	if err != nil {
		panic(fmt.Errorf("prove: %w", err))
	}

	fmt.Println("Verifying test proof...")
	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		panic(fmt.Errorf("verify: %w", err))
	}
	fmt.Println("Proof verified successfully!")

	// Export Solidity verifier
	solPath := filepath.Join(outDir, "Groth16Verifier.sol")
	solFile, err := os.Create(solPath)
	if err != nil {
		panic(err)
	}
	defer solFile.Close()
	if err := vk.ExportSolidity(solFile); err != nil {
		panic(fmt.Errorf("export solidity: %w", err))
	}
	fmt.Printf("Saved %s\n", solPath)

	// Save artifacts
	mustWriteArtifact(filepath.Join(outDir, "r1cs.bin"), r1csObj)
	mustWriteArtifact(filepath.Join(outDir, "pk.bin"), pk)
	mustWriteArtifact(filepath.Join(outDir, "vk.bin"), vk)

	fmt.Println("Done!")
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
