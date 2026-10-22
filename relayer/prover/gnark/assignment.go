package gnark

import (
	"fmt"
	"math/big"

	"github.com/decentrio/ecip-gnark/signature/eddsa"
	"github.com/decentrio/ecip-gnark/utils"

	"filippo.io/edwards25519"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/uints"

	"relayer/prover/witness"
)

// BuildBatchAssignment is the exported entry point to buildBatchAssignment, used
// by the setup-circuits smoke test so it proves against the exact assignment
// layout the relayer produces (including the #199 prefix/block-hash fields).
func BuildBatchAssignment(sigs []witness.ValidatorSignature, hash [32]byte) (*BatchCircuit[Fp25519, Fr25519], error) {
	return buildBatchAssignment(sigs, hash)
}

// buildBatchAssignment builds the BatchCircuit witness assignment: per-slot
// (sig, pub, signed canonical-vote bytes, msgLen). Signature R/S and pubkey A
// are decompressed off-circuit into gnark's emulated Ed25519 coordinates; the
// raw canonical-vote bytes flow in via Msgs/MsgLens (right-padded to witness.MaxMsgLen).
func buildBatchAssignment(sigs []witness.ValidatorSignature, hash [32]byte) (*BatchCircuit[Fp25519, Fr25519], error) {
	n := len(sigs)
	a := &BatchCircuit[Fp25519, Fr25519]{
		Sig:     make([]eddsa.Signature[Fp25519, Fr25519], n),
		Pub:     make([]eddsa.PublicKey[Fp25519, Fr25519], n),
		Msgs:    make([][witness.MaxMsgLen]uints.U8, n),
		MsgLens: make([]frontend.Variable, n),
		Active:  make([]frontend.Variable, n),
	}
	publicInputs := witness.DigestPublicInputs(hash)
	for i := range publicInputs {
		a.Hash[i] = publicInputs[i]
	}

	// #199 committed common fields (Type|Height prefix head, round, block hash)
	// — derived from a real active signer's bytes so they match the
	// witness commit and what the circuit binds each active slot against.
	prefixHead, round, blockHash, err := witness.CommonPrefixFields(sigs)
	if err != nil {
		return nil, err
	}
	for j := 0; j < witness.PrefixHeadLen; j++ {
		a.PrefixHead[j] = uints.NewU8(prefixHead[j])
	}
	for j := 0; j < witness.RoundLen; j++ {
		a.Round[j] = uints.NewU8(round[j])
	}
	for j := 0; j < 32; j++ {
		a.BlockHash[j] = uints.NewU8(blockHash[j])
	}

	for i := 0; i < n; i++ {
		v := sigs[i]
		if len(v.Signature) != 64 {
			return nil, fmt.Errorf("slot %d: invalid signature length %d", i, len(v.Signature))
		}
		if len(v.PublicKey) != 32 {
			return nil, fmt.Errorf("slot %d: invalid public key length %d", i, len(v.PublicKey))
		}
		if len(v.SignedBytes) > witness.MaxMsgLen {
			return nil, fmt.Errorf("slot %d: signed bytes length %d exceeds witness.MaxMsgLen=%d", i, len(v.SignedBytes), witness.MaxMsgLen)
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
		for j := 0; j < witness.MaxMsgLen; j++ {
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

func ProofToBigInts(gnarkProof groth16.Proof) ([8]*big.Int, [2]*big.Int, [2]*big.Int, error) {
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

	p, ok := gnarkProof.(*groth16_bn254.Proof)
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
