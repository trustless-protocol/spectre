package prover

import (
	"crypto/sha512"
	"fmt"
	"math/big"
	"os"

	"0x5ea000000/ecip-gnark/signature/eddsa"
	"0x5ea000000/ecip-gnark/utils"

	"filippo.io/edwards25519"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/consensys/gnark/backend/solidity"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/emulated"
)

type EcipProver struct {
	r1cs constraint.ConstraintSystem
	pk   groth16.ProvingKey
	vk   groth16.VerifyingKey
}

func NewProver(r1csPath, pkPath, vkPath string) (*EcipProver, error) {
	r1csFile, err := os.Open(r1csPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open r1cs file: %w", err)
	}
	defer r1csFile.Close()

	r1cs := groth16.NewCS(ecc.BN254)
	if _, err := r1cs.ReadFrom(r1csFile); err != nil {
		return nil, fmt.Errorf("failed to read r1cs: %w", err)
	}

	pk := groth16.NewProvingKey(ecc.BN254)
	pkFile, err := os.Open(pkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open proving key file: %w", err)
	}
	defer pkFile.Close()

	if _, err := pk.ReadFrom(pkFile); err != nil {
		return nil, fmt.Errorf("failed to read proving key: %w", err)
	}

	vk := groth16.NewVerifyingKey(ecc.BN254)
	vkFile, err := os.Open(vkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open verifying key file: %w", err)
	}
	defer vkFile.Close()
	if _, err := vk.ReadFrom(vkFile); err != nil {
		return nil, fmt.Errorf("failed to read verifying key: %w", err)
	}

	return &EcipProver{
		r1cs: r1cs,
		pk:   pk,
		vk:   vk,
	}, nil
}

func (p *EcipProver) GenerateProof(sig, pub, msg []byte) (
	proof [8]*big.Int,
	commitments [2]*big.Int,
	commitmentPok [2]*big.Int,
	err error,
) {
	if len(sig) != 64 {
		err = fmt.Errorf("invalid signature length: %d, expected 64", len(sig))
		return
	}
	if len(pub) != 32 {
		err = fmt.Errorf("invalid public key length: %d, expected 32", len(pub))
		return
	}

	// Extract R (first 32 bytes) and S (last 32 bytes) from signature
	R := sig[:32]
	S, sErr := edwards25519.NewScalar().SetCanonicalBytes(sig[32:])
	if sErr != nil {
		err = fmt.Errorf("invalid signature scalar S: %w", sErr)
		return
	}

	// Decompress Ed25519 points to Weierstrass coordinates
	aX, aY, decompErr := utils.DecompressPoint(pub)
	if decompErr != nil {
		err = fmt.Errorf("failed to decompress public key: %w", decompErr)
		return
	}

	rX, rY, decompErr := utils.DecompressPoint(R)
	if decompErr != nil {
		err = fmt.Errorf("failed to decompress R point: %w", decompErr)
		return
	}

	// Convert S scalar to big.Int
	s := utils.ScalarToBigInt(S)

	// Compute H = SHA512(R || A || msg) off-chain
	hasher := sha512.New()
	hasher.Write(R)
	hasher.Write(pub)
	hasher.Write(msg)
	sum := hasher.Sum(nil)

	H, hErr := edwards25519.NewScalar().SetUniformBytes(sum)
	if hErr != nil {
		err = fmt.Errorf("failed to set hash scalar: %w", hErr)
		return
	}
	h := utils.ScalarToBigInt(H)

	// Build witness assignment
	assignment := PreHashCircuit[Fp25519, Fr25519]{
		Sig: eddsa.Signature[Fp25519, Fr25519]{
			R: sw_emulated.AffinePoint[Fp25519]{
				X: emulated.ValueOf[Fp25519](rX),
				Y: emulated.ValueOf[Fp25519](rY),
			},
			S: emulated.ValueOf[Fr25519](s),
		},
		Hash: emulated.ValueOf[Fr25519](h),
		Pub: eddsa.PublicKey[Fp25519, Fr25519]{
			A: sw_emulated.AffinePoint[Fp25519]{
				X: emulated.ValueOf[Fp25519](aX),
				Y: emulated.ValueOf[Fp25519](aY),
			},
		},
	}

	// Create witness
	witness, wErr := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if wErr != nil {
		err = fmt.Errorf("failed to create witness: %w", wErr)
		return
	}

	// Generate proof — use keccak256 for commitment hash to match Solidity verifier
	gnarkProof, pErr := groth16.Prove(p.r1cs, p.pk, witness, solidity.WithProverTargetSolidityVerifier(backend.GROTH16))
	if pErr != nil {
		err = fmt.Errorf("failed to generate proof: %w", pErr)
		return
	}

	// Local verification with VK (same keccak256 hash as Solidity)
	pubWitness, _ := witness.Public()
	if vErr := groth16.Verify(gnarkProof, p.vk, pubWitness, solidity.WithVerifierTargetSolidityVerifier(backend.GROTH16)); vErr != nil {
		err = fmt.Errorf("LOCAL VERIFICATION FAILED: %w", vErr)
		return
	}

	// Convert to Solidity-compatible format
	return ProofToBigInts(gnarkProof)
}

func ProofToBigInts(proof groth16.Proof) ([8]*big.Int, [2]*big.Int, [2]*big.Int, error) {
	var out [8]*big.Int
	for i := range out {
		out[i] = new(big.Int)
	}
	var commitmentPoks [2]*big.Int
	for i := range commitmentPoks {
		commitmentPoks[i] = new(big.Int)
	}
	var commitments [2]*big.Int
	for i := range commitments {
		commitments[i] = new(big.Int)
	}

	p, ok := proof.(*groth16_bn254.Proof)
	if !ok {
		return out, commitmentPoks, commitments, fmt.Errorf("expected BN254 proof")
	}

	// A (G1)
	p.Ar.X.BigInt(out[0])
	p.Ar.Y.BigInt(out[1])

	// B (G2) — EIP-197 and gnark's MarshalSolidity expect (A1, A0) order
	// A1 = imaginary part first, A0 = real part second
	p.Bs.X.A1.BigInt(out[2])
	p.Bs.X.A0.BigInt(out[3])
	p.Bs.Y.A1.BigInt(out[4])
	p.Bs.Y.A0.BigInt(out[5])

	// C (G1)
	p.Krs.X.BigInt(out[6])
	p.Krs.Y.BigInt(out[7])

	p.CommitmentPok.X.BigInt(commitmentPoks[0])
	p.CommitmentPok.Y.BigInt(commitmentPoks[1])

	p.Commitments[0].X.BigInt(commitments[0])
	p.Commitments[0].Y.BigInt(commitments[1])
	return out, commitments, commitmentPoks, nil
}
