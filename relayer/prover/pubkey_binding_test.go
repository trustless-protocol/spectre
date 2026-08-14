package prover

import (
	"crypto/ed25519"
	"crypto/rand"
	"math/big"
	"testing"

	"0x5ea000000/ecip-gnark/utils"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/test"
)

// compressBindingCircuit isolates compressEdwardsToLE — the gadget this PR
// fixes — from the rest of the batch circuit. A point is committed by its
// compressed encoding (public Enc), and nothing else constrains it. This is the
// ZK-06 threat model in miniature: if the encoding does not determine the point,
// two different witnesses satisfy the same Enc.
//
// The isolation is deliberate. In the full BatchCircuit the ECIP batch verify
// multiplies each coordinate, and those emulated field mul-checks reject a
// non-canonical limb packing on their own — so a full-circuit witness cannot
// tell the fixed gadget from the broken one (it is rejected either way, for the
// wrong reason). Only a circuit that exercises the compression alone can prove
// that compressEdwardsToLE's own asserts are what bind the representation.
type compressBindingCircuit struct {
	P   sw_emulated.AffinePoint[Fp25519]
	Enc [32]frontend.Variable `gnark:",public"`
}

func (c *compressBindingCircuit) Define(api frontend.API) error {
	baseApi, err := emulated.NewField[Fp25519](api)
	if err != nil {
		return err
	}
	enc := compressEdwardsToLE(api, baseApi, &c.P)
	for i := range enc {
		api.AssertIsEqual(enc[i].Val, c.Enc[i])
	}
	return nil
}

// rawElement packs v into Curve25519Fp's four 64-bit limbs without reducing it
// modulo p, so a non-canonical representation (a coordinate >= p) can be
// submitted. emulated.ValueOf reduces during witness parsing and cannot express
// one; an element with Limbs set directly is read as given.
func rawElement(v *big.Int) emulated.Element[Fp25519] {
	var fp Fp25519
	limbs := make([]frontend.Variable, fp.NbLimbs())
	mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), fp.BitsPerLimb()), big.NewInt(1))
	rest := new(big.Int).Set(v)
	for i := range limbs {
		limbs[i] = new(big.Int).And(rest, mask)
		rest = new(big.Int).Rsh(rest, uint(fp.BitsPerLimb()))
	}
	return emulated.Element[Fp25519]{Limbs: limbs}
}

// encodingOf returns the 32-byte encoding compressEdwardsToLE computes for the
// given limb representation: the low 255 bits of y little-endian, with the
// parity of x in bit 255. It reads the representation as given, exactly like the
// in-circuit ToBits — so for a non-canonical y it differs from y's canonical
// encoding.
func encodingOf(x, y *big.Int) [32]byte {
	var out [32]byte
	yLow := new(big.Int).Mod(y, new(big.Int).Lsh(big.NewInt(1), 255))
	yLow.FillBytes(out[:]) // big-endian
	for i, j := 0, 31; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i] // to little-endian
	}
	out[31] |= byte(x.Bit(0)) << 7
	return out
}

func encWitness(enc [32]byte) [32]frontend.Variable {
	var out [32]frontend.Variable
	for i := range enc {
		out[i] = enc[i]
	}
	return out
}

func solveBinding(t *testing.T, x, y *big.Int, enc [32]byte) error {
	t.Helper()
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &compressBindingCircuit{})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	w, err := frontend.NewWitness(&compressBindingCircuit{
		P:   sw_emulated.AffinePoint[Fp25519]{X: rawElement(x), Y: rawElement(y)},
		Enc: encWitness(enc),
	}, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatalf("witness: %v", err)
	}
	_, err = cs.Solve(w)
	return err
}

// TestCompressEdwardsToLE_BindsCanonicalRepresentation is the ZK-06 regression
// for the half of the fix that lives in this repo: the two AssertIsInRange calls
// in compressEdwardsToLE (hash_witness.go).
//
// A witness coordinate is only width-constrained, so v and v+p are both legal
// limb packings of the same field value, and the compression reads the limbs as
// given (ToBits, not ToBitsCanonical). So a prover can hand the compression the
// bits of v+p while the field value is v, and the emitted encoding then differs
// from v's canonical encoding. Without the asserts the encoding is not injective
// in the representation, and the bytes the on-chain quorum check matches no
// longer pin the point.
//
// The negative case uses the point (sqrt(-1), 0) with Y supplied as 0+p: the low
// 255 bits of p are all set, so compression emits p's bytes for a coordinate
// whose field value is 0. It discriminates the fix on two axes. With the two
// AssertIsInRange lines removed it SOLVES (verified by stripping them and re-
// running, at the branch's own submodule pins). And the on-curve asserts landing
// in ecip-gnark#8 do not catch it either — (sqrt(-1), 0) genuinely satisfies
// -x^2 + y^2 = 1 — so it is the canonical-limb assert specifically, the half in
// this repo, that rejects it.
//
// This gadget is exercised in isolation on purpose: see compressBindingCircuit.
func TestCompressEdwardsToLE_BindsCanonicalRepresentation(t *testing.T) {
	p := new(Fp25519).Modulus()

	// Positive control: an honest pubkey and its own encoding must still solve,
	// so the negative case cannot pass on a broken build or panicking gadget.
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	hx, hy, err := utils.DecompressPoint(pub)
	if err != nil {
		t.Fatalf("decompress pubkey: %v", err)
	}
	var honestEnc [32]byte
	copy(honestEnc[:], pub)
	if got := encodingOf(hx, hy); got != honestEnc {
		t.Fatalf("encodingOf disagrees with canonical pubkey bytes: %x vs %x", got, honestEnc)
	}
	if err := solveBinding(t, hx, hy, honestEnc); err != nil {
		t.Fatalf("honest point rejected: %v", err)
	}

	// Negative: the point (sqrt(-1), 0) with Y supplied as 0 + p — a
	// non-canonical representation of Y = 0 — together with the encoding that
	// representation produces. Self-consistent and on the curve (-x^2 + y^2 = 1),
	// yet a second witness whose canonical encoding is different. Only
	// AssertIsInRange rejects it.
	exp := new(big.Int).Rsh(new(big.Int).Sub(p, big.NewInt(1)), 2)
	x := new(big.Int).Exp(big.NewInt(2), exp, p) // sqrt(-1) = 2^((p-1)/4) mod p
	if check := new(big.Int).Exp(x, big.NewInt(2), p); check.Cmp(new(big.Int).Sub(p, big.NewInt(1))) != 0 {
		t.Fatalf("sqrt(-1) is wrong: x^2 = %s", check)
	}
	y := new(big.Int).Set(p) // 0 + p
	if err := solveBinding(t, x, y, encodingOf(x, y)); err == nil {
		t.Fatal("non-canonical coordinate representation was accepted")
	}
}

// TestBatchCircuit_RejectsOffCurvePubkey is an end-to-end sanity check that a
// perturbed pubkey coordinate does not solve under the honest witness
// commitment. It is NOT the ZK-06 binding proof.
//
// The x+2 perturbation reuses the honest signature, which already fails the
// Ed25519 verification equation against the moved point, so this test passes
// with every ZK-06 assert removed — a green run here proves nothing about the
// fix. The discriminating cases are elsewhere: the canonical-limb half in
// TestCompressEdwardsToLE_BindsCanonicalRepresentation above, and the on-curve
// half (off-curve, 2-torsion) in the ecip-gnark#8 negative suite, each verified
// to solve with the asserts removed.
func TestBatchCircuit_RejectsOffCurvePubkey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}
	const n = 4
	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0x30 + i)
	}
	vote := makeCanonicalVote(false, 10, blockHash)
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	padded, err := padWithDummies([]ValidatorSignature{{
		Signature:   ed25519.Sign(priv, vote),
		PublicKey:   pub,
		SignedBytes: vote,
		Active:      true,
	}}, generateDummySlots(n))
	if err != nil {
		t.Fatalf("padWithDummies: %v", err)
	}
	hash, err := ComputeWitnessHash(padded)
	if err != nil {
		t.Fatalf("ComputeWitnessHash: %v", err)
	}
	assignment, err := buildBatchAssignment(padded, hash)
	if err != nil {
		t.Fatalf("buildBatchAssignment: %v", err)
	}
	if err := test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField()); err != nil {
		t.Fatalf("honest assignment rejected: %v", err)
	}

	// Same public inputs, same committed pubkey bytes, different point.
	aX, _, err := utils.DecompressPoint(pub)
	if err != nil {
		t.Fatalf("decompress pubkey: %v", err)
	}
	forged := new(big.Int).Add(aX, big.NewInt(2))
	if forged.Bit(0) != aX.Bit(0) {
		t.Fatal("forged X changed the sign bit; the encoding would differ")
	}
	assignment.Pub[0].A.X = emulated.ValueOf[Fp25519](forged)
	if err := test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField()); err == nil {
		t.Fatal("off-curve pubkey accepted under the honest witness commitment")
	}
}
