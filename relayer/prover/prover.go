package prover

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"

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

	benchcfg "relayer/benchmark"
)

// bucketArtifacts holds the compiled circuit, Groth16 key pair, and the
// deterministic dummy padding slots for one validator-count bucket.
type bucketArtifacts struct {
	n      int
	r1cs   constraint.ConstraintSystem
	pk     groth16.ProvingKey
	vk     groth16.VerifyingKey
	dummys []dummySignature
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
	if benchcfg.Enabled() {
		log.Printf("[NewProver] loading artifacts from %s", binDir)
	}
	p := &EcipProver{byBucket: make(map[int]*bucketArtifacts, len(Buckets))}
	for _, n := range Buckets {
		if benchcfg.Enabled() {
			log.Printf("[NewProver] loading bucket n=%d", n)
		}
		art, err := loadBucketArtifacts(binDir, n)
		if err != nil {
			return nil, fmt.Errorf("load bucket n=%d: %w", n, err)
		}
		p.byBucket[n] = art
		if benchcfg.Enabled() {
			log.Printf("[NewProver] bucket n=%d loaded", n)
		}
	}
	if benchcfg.Enabled() {
		log.Printf("[NewProver] loaded %d bucket(s)", len(p.byBucket))
	}
	return p, nil
}

func loadBucketArtifacts(binDir string, n int) (*bucketArtifacts, error) {
	dir := filepath.Join(binDir, fmt.Sprintf("n%d", n))
	if benchcfg.Enabled() {
		log.Printf("[NewProver] bucket n=%d dir=%s", n, dir)
	}

	r1cs := groth16.NewCS(ecc.BN254)
	if err := readFromFile(filepath.Join(dir, "r1cs.bin"), "r1cs", n, r1cs); err != nil {
		return nil, fmt.Errorf("read r1cs: %w", err)
	}

	pk := groth16.NewProvingKey(ecc.BN254)
	if err := readFromFile(filepath.Join(dir, "pk.bin"), "pk", n, pk); err != nil {
		return nil, fmt.Errorf("read pk: %w", err)
	}

	vk := groth16.NewVerifyingKey(ecc.BN254)
	if err := readFromFile(filepath.Join(dir, "vk.bin"), "vk", n, vk); err != nil {
		return nil, fmt.Errorf("read vk: %w", err)
	}

	return &bucketArtifacts{
		n:      n,
		r1cs:   r1cs,
		pk:     pk,
		vk:     vk,
		dummys: generateDummySlots(n),
	}, nil
}

type readerFrom interface {
	ReadFrom(r io.Reader) (int64, error)
}

func readFromFile(path string, kind string, bucket int, dst readerFrom) error {
	if benchcfg.Enabled() {
		if fi, err := os.Stat(path); err == nil {
			log.Printf("[NewProver] bucket n=%d reading %s from %s (%d bytes)", bucket, kind, path, fi.Size())
		} else {
			log.Printf("[NewProver] bucket n=%d stat failed for %s %s: %v", bucket, kind, path, err)
		}
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if benchcfg.Enabled() {
		log.Printf("[NewProver] bucket n=%d opened %s", bucket, kind)
	}
	_, err = dst.ReadFrom(f)
	if err == nil && benchcfg.Enabled() {
		log.Printf("[NewProver] bucket n=%d loaded %s", bucket, kind)
	}
	return err
}

// GenerateProof produces a Groth16 proof that every (sig[i], pub[i]) validly
// signed the canonical-vote bytes carried in sigs[i].SignedBytes. The prover
// picks the smallest bucket that fits len(sigs) and pads the witness up to
// the bucket's N with copies of slot 0. The returned bucket value tells the
// caller which Solidity verifier to dispatch to.
func (p *EcipProver) GenerateProof(sigs []ValidatorSignature) (
	bucket int,
	paddedSigs []ValidatorSignature,
	proof [8]*big.Int,
	commitments [2]*big.Int,
	commitmentPok [2]*big.Int,
	err error,
) {
	if len(sigs) == 0 {
		err = fmt.Errorf("no signatures provided")
		return
	}

	benchEnabled := benchcfg.Enabled()
	var totalStart time.Time
	if benchEnabled {
		totalStart = time.Now()
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
	if benchEnabled {
		log.Printf("[bench][prover] start sigs=%d bucket=%d", len(sigs), bucket)
	}

	paddedSigs, err = padWithDummies(sigs, art.dummys)
	if err != nil {
		return
	}

	var witnessStart time.Time
	if benchEnabled {
		witnessStart = time.Now()
	}
	hash, err := ComputeWitnessHash(paddedSigs)
	if err != nil {
		err = fmt.Errorf("compute witness hash: %w", err)
		return
	}

	assignment, err := buildBatchAssignment(paddedSigs, hash)
	if err != nil {
		return
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		err = fmt.Errorf("create witness: %w", err)
		return
	}
	var witnessDur time.Duration
	if benchEnabled {
		witnessDur = time.Since(witnessStart)
	}

	var proveStart time.Time
	if benchEnabled {
		proveStart = time.Now()
	}
	gnarkProof, err := groth16.Prove(art.r1cs, art.pk, witness, solidity.WithProverTargetSolidityVerifier(backend.GROTH16))
	if err != nil {
		err = fmt.Errorf("generate proof: %w", err)
		return
	}
	var proveDur time.Duration
	if benchEnabled {
		proveDur = time.Since(proveStart)
	}

	var verifyStart time.Time
	if benchEnabled {
		verifyStart = time.Now()
	}
	pubWitness, _ := witness.Public()
	if vErr := groth16.Verify(gnarkProof, art.vk, pubWitness, solidity.WithVerifierTargetSolidityVerifier(backend.GROTH16)); vErr != nil {
		err = fmt.Errorf("local verification failed: %w", vErr)
		return
	}
	var verifyDur time.Duration
	if benchEnabled {
		verifyDur = time.Since(verifyStart)
	}

	proof, commitments, commitmentPok, err = ProofToBigInts(gnarkProof)
	if benchEnabled {
		log.Printf("[bench][prover] done sigs=%d bucket=%d witness=%s prove=%s verify=%s total=%s",
			len(sigs), bucket, witnessDur, proveDur, verifyDur, time.Since(totalStart))
	}
	return
}

// padWithDummies fills the trailing slots of a bucket with deterministic dummy
// signatures (active=false) so every (R, A) point in the batch is distinct.
// Duplicating a real slot would re-introduce the ECIP doubling-branch issue;
// distinct dummies keep the divisor well-formed while the active gate zeroes
// their contribution to the aggregate.
func padWithDummies(sigs []ValidatorSignature, dummies []dummySignature) ([]ValidatorSignature, error) {
	n := len(dummies)
	if len(sigs) > n {
		return nil, fmt.Errorf("bucket n=%d too small for %d signers", n, len(sigs))
	}
	out := make([]ValidatorSignature, n)
	copy(out, sigs)
	for i := len(sigs); i < n; i++ {
		out[i] = dummies[i].asValidatorSignature()
	}
	return out, nil
}

// buildBatchAssignment builds the BatchCircuit witness assignment: per-slot
// (sig, pub, signed canonical-vote bytes, msgLen). Signature R/S and pubkey A
// are decompressed off-circuit into gnark's emulated Ed25519 coordinates; the
// raw canonical-vote bytes flow in via Msgs/MsgLens (right-padded to MaxMsgLen).
func buildBatchAssignment(sigs []ValidatorSignature, hash [32]byte) (*BatchCircuit[Fp25519, Fr25519], error) {
	n := len(sigs)
	a := &BatchCircuit[Fp25519, Fr25519]{
		Sig:     make([]eddsa.Signature[Fp25519, Fr25519], n),
		Pub:     make([]eddsa.PublicKey[Fp25519, Fr25519], n),
		Msgs:    make([][MaxMsgLen]uints.U8, n),
		MsgLens: make([]frontend.Variable, n),
		Active:  make([]frontend.Variable, n),
	}
	publicInputs := DigestPublicInputs(hash)
	for i := range publicInputs {
		a.Hash[i] = publicInputs[i]
	}

	for i := 0; i < n; i++ {
		v := sigs[i]
		if len(v.Signature) != 64 {
			return nil, fmt.Errorf("slot %d: invalid signature length %d", i, len(v.Signature))
		}
		if len(v.PublicKey) != 32 {
			return nil, fmt.Errorf("slot %d: invalid public key length %d", i, len(v.PublicKey))
		}
		if len(v.SignedBytes) > MaxMsgLen {
			return nil, fmt.Errorf("slot %d: signed bytes length %d exceeds MaxMsgLen=%d", i, len(v.SignedBytes), MaxMsgLen)
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
		for j := 0; j < MaxMsgLen; j++ {
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
