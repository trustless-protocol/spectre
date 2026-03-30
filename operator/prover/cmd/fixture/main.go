package main

import (
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"0x5ea000000/ecip-gnark/signature/eddsa"
	"0x5ea000000/ecip-gnark/utils"

	"filippo.io/edwards25519"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/emulated"

	"operator/prover"
)

// Fixture is the JSON structure for Solidity test fixtures.
type Fixture struct {
	Proof         [8]string `json:"proof"`
	Commitments   [2]string `json:"commitments"`
	CommitmentPok [2]string `json:"commitmentPok"`
	SignatureR    string    `json:"signatureR"`
	SignatureS    string    `json:"signatureS"`
	Pubkey        string    `json:"pubkey"`
	Message       string    `json:"message"`
}

func bigIntToHex(b *big.Int) string {
	return "0x" + b.Text(16)
}

func bytesToHex(b []byte) string {
	return fmt.Sprintf("0x%x", b)
}

func main() {
	binDir := "prover/bin"
	if len(os.Args) > 1 {
		binDir = os.Args[1]
	}
	outPath := "test/sp1-ics07/fixtures/groth16_fixture.json"
	if len(os.Args) > 2 {
		outPath = os.Args[2]
	}

	// Load R1CS and PK
	r1csFile, err := os.Open(binDir + "/r1cs.bin")
	if err != nil {
		panic(fmt.Errorf("open r1cs: %w", err))
	}
	defer r1csFile.Close()
	r1cs := groth16.NewCS(ecc.BN254)
	if _, err := r1cs.ReadFrom(r1csFile); err != nil {
		panic(fmt.Errorf("read r1cs: %w", err))
	}

	pkFile, err := os.Open(binDir + "/pk.bin")
	if err != nil {
		panic(fmt.Errorf("open pk: %w", err))
	}
	defer pkFile.Close()
	pk := groth16.NewProvingKey(ecc.BN254)
	if _, err := pk.ReadFrom(pkFile); err != nil {
		panic(fmt.Errorf("read pk: %w", err))
	}

	// Generate Ed25519 keypair and sign
	msg := []byte("fixture-test-message")
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	sig := ed25519.Sign(priv, msg)

	R := sig[:32]
	S, err := edwards25519.NewScalar().SetCanonicalBytes(sig[32:])
	if err != nil {
		panic("invalid S scalar")
	}

	// SHA512(R || A || msg)
	hasher := sha512.New()
	hasher.Write(R)
	hasher.Write([]byte(pub))
	hasher.Write(msg)
	sum := hasher.Sum(nil)
	H, err := edwards25519.NewScalar().SetUniformBytes(sum)
	if err != nil {
		panic("setting hash scalar failed")
	}

	aX, aY, _ := utils.DecompressPoint([]byte(pub))
	rX, rY, _ := utils.DecompressPoint(R)
	h := utils.ScalarToBigInt(H)
	s := utils.ScalarToBigInt(S)

	// Build witness
	assignment := prover.PreHashCircuit[prover.Fp25519, prover.Fr25519]{
		Sig: eddsa.Signature[prover.Fp25519, prover.Fr25519]{
			R: sw_emulated.AffinePoint[prover.Fp25519]{
				X: emulated.ValueOf[prover.Fp25519](rX),
				Y: emulated.ValueOf[prover.Fp25519](rY),
			},
			S: emulated.ValueOf[prover.Fr25519](s),
		},
		Hash: emulated.ValueOf[prover.Fr25519](h),
		Pub: eddsa.PublicKey[prover.Fp25519, prover.Fr25519]{
			A: sw_emulated.AffinePoint[prover.Fp25519]{
				X: emulated.ValueOf[prover.Fp25519](aX),
				Y: emulated.ValueOf[prover.Fp25519](aY),
			},
		},
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		panic(fmt.Errorf("witness: %w", err))
	}

	fmt.Println("Generating proof...")
	gnarkProof, err := groth16.Prove(r1cs.(constraint.ConstraintSystem), pk, witness)
	if err != nil {
		panic(fmt.Errorf("prove: %w", err))
	}

	// Convert proof to big ints
	proofBigInt, commitments, commitmentPok, err := prover.ProofToBigInts(gnarkProof)
	if err != nil {
		panic(fmt.Errorf("proof to big ints: %w", err))
	}

	// Build fixture
	var proofHex [8]string
	for i, v := range proofBigInt {
		proofHex[i] = bigIntToHex(v)
	}
	var commitmentsHex [2]string
	for i, v := range commitments {
		commitmentsHex[i] = bigIntToHex(v)
	}
	var commitmentPokHex [2]string
	for i, v := range commitmentPok {
		commitmentPokHex[i] = bigIntToHex(v)
	}

	// Message as keccak256 hash (bytes32) — matches what operator sends on-chain
	msgHash := make([]byte, 32)
	copy(msgHash, msg) // For test: just pad msg to 32 bytes

	fixture := Fixture{
		Proof:         proofHex,
		Commitments:   commitmentsHex,
		CommitmentPok: commitmentPokHex,
		SignatureR:    bytesToHex(R),
		SignatureS:    bytesToHex(sig[32:]),
		Pubkey:        bytesToHex([]byte(pub)),
		Message:       bytesToHex(msgHash),
	}

	data, err := json.MarshalIndent(fixture, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		panic(fmt.Errorf("write fixture: %w", err))
	}
	fmt.Printf("Saved fixture to %s\n", outPath)
}
