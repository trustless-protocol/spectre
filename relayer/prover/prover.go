package prover

import (
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"

	"0x5ea000000/ecip-gnark/signature/canonvote"
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
	"github.com/consensys/gnark/std/math/uints"
)

// bucketArtifacts holds the compiled circuit and Groth16 key pair for one
// validator-count bucket.
type bucketArtifacts struct {
	n    int
	r1cs constraint.ConstraintSystem
	pk   groth16.ProvingKey
	vk   groth16.VerifyingKey
}

// EcipProver is a registry of compiled circuits keyed by bucket size. The
// prover selects the smallest bucket that fits the required signer count for
// the current block and uses that bucket's artifacts to generate a proof.
type EcipProver struct {
	byBucket map[int]*bucketArtifacts
}

// NewProver loads every bucket's r1cs, proving key, and verifying key from
// binDir/n{N}/{r1cs,pk,vk}.bin. It errors if any bucket's artifacts are
// missing — the operator must run `cmd/setup-circuits` first.
func NewProver(binDir string) (*EcipProver, error) {
	p := &EcipProver{byBucket: make(map[int]*bucketArtifacts, len(Buckets))}
	for _, n := range Buckets {
		art, err := loadBucketArtifacts(binDir, n)
		if err != nil {
			return nil, fmt.Errorf("load bucket n=%d: %w", n, err)
		}
		p.byBucket[n] = art
	}
	return p, nil
}

func loadBucketArtifacts(binDir string, n int) (*bucketArtifacts, error) {
	dir := filepath.Join(binDir, fmt.Sprintf("n%d", n))

	r1cs := groth16.NewCS(ecc.BN254)
	if err := readFromFile(filepath.Join(dir, "r1cs.bin"), r1cs); err != nil {
		return nil, fmt.Errorf("read r1cs: %w", err)
	}

	pk := groth16.NewProvingKey(ecc.BN254)
	if err := readFromFile(filepath.Join(dir, "pk.bin"), pk); err != nil {
		return nil, fmt.Errorf("read pk: %w", err)
	}

	vk := groth16.NewVerifyingKey(ecc.BN254)
	if err := readFromFile(filepath.Join(dir, "vk.bin"), vk); err != nil {
		return nil, fmt.Errorf("read vk: %w", err)
	}

	return &bucketArtifacts{n: n, r1cs: r1cs, pk: pk, vk: vk}, nil
}

type readerFrom interface {
	ReadFrom(r io.Reader) (int64, error)
}

func readFromFile(path string, dst readerFrom) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = dst.ReadFrom(f)
	return err
}

// GenerateProof produces a Groth16 proof that every (sig[i], pub[i]) validly
// signed the canonical-vote bytes reconstructed from (shared, timestamps[i]).
// The prover picks the smallest bucket that fits len(sigs) and pads the
// witness up to the bucket's N with copies of slot 0. The returned bucket
// value tells the caller which Solidity verifier to dispatch to.
func (p *EcipProver) GenerateProof(shared SharedBlockData, sigs []ValidatorSignature) (
	bucket int,
	proof [8]*big.Int,
	commitments [2]*big.Int,
	commitmentPok [2]*big.Int,
	err error,
) {
	if len(sigs) == 0 {
		err = fmt.Errorf("no signatures provided")
		return
	}
	if len(shared.BlockIDHash) != 32 {
		err = fmt.Errorf("BlockIDHash must be 32 bytes, got %d", len(shared.BlockIDHash))
		return
	}
	if len(shared.PartSetHash) != 32 {
		err = fmt.Errorf("PartSetHash must be 32 bytes, got %d", len(shared.PartSetHash))
		return
	}
	if len(shared.ChainID) > canonvote.MaxChainIDLen {
		err = fmt.Errorf("chainID %q exceeds MaxChainIDLen=%d", shared.ChainID, canonvote.MaxChainIDLen)
		return
	}

	bucket, err = SmallestBucketGEQ(len(sigs))
	if err != nil {
		return
	}
	art, ok := p.byBucket[bucket]
	if !ok {
		err = fmt.Errorf("bucket n=%d not loaded", bucket)
		return
	}

	paddedSigs, err := padSigsToBucket(sigs, bucket)
	if err != nil {
		return
	}

	assignment, err := buildBatchAssignment(shared, paddedSigs)
	if err != nil {
		return
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		err = fmt.Errorf("create witness: %w", err)
		return
	}

	gnarkProof, err := groth16.Prove(art.r1cs, art.pk, witness, solidity.WithProverTargetSolidityVerifier(backend.GROTH16))
	if err != nil {
		err = fmt.Errorf("generate proof: %w", err)
		return
	}

	pubWitness, _ := witness.Public()
	if vErr := groth16.Verify(gnarkProof, art.vk, pubWitness, solidity.WithVerifierTargetSolidityVerifier(backend.GROTH16)); vErr != nil {
		err = fmt.Errorf("local verification failed: %w", vErr)
		return
	}

	proof, commitments, commitmentPok, err = ProofToBigInts(gnarkProof)
	return
}

// padSigsToBucket pads the signer slice up to `n` with copies of slot 0.
// Duplicate slots are cryptographically sound — each slot still holds a real
// signature; the on-chain voting-power counter dedupes by validator index.
func padSigsToBucket(sigs []ValidatorSignature, n int) ([]ValidatorSignature, error) {
	if len(sigs) > n {
		return nil, fmt.Errorf("bucket n=%d too small for %d signers", n, len(sigs))
	}
	out := make([]ValidatorSignature, n)
	for i := 0; i < n; i++ {
		src := i
		if src >= len(sigs) {
			src = 0
		}
		out[i] = sigs[src]
	}
	return out, nil
}

// buildBatchAssignment builds the BatchCircuit witness assignment: shared
// block data + per-slot (sig, pub, timestamp). Signature values are
// decompressed off-circuit into gnark's emulated Ed25519 coordinates; the
// canonical-vote bytes themselves are reconstructed inside the circuit so
// they do not appear in this witness.
func buildBatchAssignment(shared SharedBlockData, sigs []ValidatorSignature) (*BatchCircuit[Fp25519, Fr25519], error) {
	n := len(sigs)
	a := &BatchCircuit[Fp25519, Fr25519]{
		Sig:       make([]eddsa.Signature[Fp25519, Fr25519], n),
		Pub:       make([]eddsa.PublicKey[Fp25519, Fr25519], n),
		TsSeconds: make([]frontend.Variable, n),
		TsNanos:   make([]frontend.Variable, n),
	}

	for i := 0; i < n; i++ {
		v := sigs[i]
		if len(v.Signature) != 64 {
			return nil, fmt.Errorf("slot %d: invalid signature length %d", i, len(v.Signature))
		}
		if len(v.PublicKey) != 32 {
			return nil, fmt.Errorf("slot %d: invalid public key length %d", i, len(v.PublicKey))
		}

		R := v.Signature[:32]
		S, err := edwards25519.NewScalar().SetCanonicalBytes(v.Signature[32:])
		if err != nil {
			return nil, fmt.Errorf("slot %d: scalar S: %w", i, err)
		}
		aX, aY, err := utils.DecompressPoint(v.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("slot %d: decompress pubkey: %w", i, err)
		}
		rX, rY, err := utils.DecompressPoint(R)
		if err != nil {
			return nil, fmt.Errorf("slot %d: decompress R: %w", i, err)
		}

		a.Sig[i] = eddsa.Signature[Fp25519, Fr25519]{
			R: sw_emulated.AffinePoint[Fp25519]{
				X: emulated.ValueOf[Fp25519](rX),
				Y: emulated.ValueOf[Fp25519](rY),
			},
			S: emulated.ValueOf[Fr25519](utils.ScalarToBigInt(S)),
		}
		a.Pub[i] = eddsa.PublicKey[Fp25519, Fr25519]{
			A: sw_emulated.AffinePoint[Fp25519]{
				X: emulated.ValueOf[Fp25519](aX),
				Y: emulated.ValueOf[Fp25519](aY),
			},
		}
		a.TsSeconds[i] = v.TimestampSeconds
		a.TsNanos[i] = v.TimestampNanos
	}

	// Shared block data — fixed-width byte arrays padded with zeros.
	for i := 0; i < 32; i++ {
		a.BlockIDHash[i] = uints.NewU8(shared.BlockIDHash[i])
		a.PartSetHash[i] = uints.NewU8(shared.PartSetHash[i])
	}
	for i := 0; i < canonvote.MaxChainIDLen; i++ {
		if i < len(shared.ChainID) {
			a.ChainID[i] = uints.NewU8(shared.ChainID[i])
		} else {
			a.ChainID[i] = uints.NewU8(0)
		}
	}
	a.Height = shared.Height
	a.Round = shared.Round
	a.PartSetTotal = shared.PartSetTotal
	a.ChainIDLen = len(shared.ChainID)

	return a, nil
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

	p.Ar.X.BigInt(out[0])
	p.Ar.Y.BigInt(out[1])

	p.Bs.X.A1.BigInt(out[2])
	p.Bs.X.A0.BigInt(out[3])
	p.Bs.Y.A1.BigInt(out[4])
	p.Bs.Y.A0.BigInt(out[5])

	p.Krs.X.BigInt(out[6])
	p.Krs.Y.BigInt(out[7])

	p.CommitmentPok.X.BigInt(commitmentPoks[0])
	p.CommitmentPok.Y.BigInt(commitmentPoks[1])

	p.Commitments[0].X.BigInt(commitments[0])
	p.Commitments[0].Y.BigInt(commitments[1])
	return out, commitments, commitmentPoks, nil
}
