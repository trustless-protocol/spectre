package gnark

import (
	"relayer/prover/witness"

	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"testing"

	"filippo.io/edwards25519"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/test"
)

// ZK-08: negative tests for the circuit. Every test here builds an assignment
// that an honest prover would never produce but a malicious one can, and asserts
// the circuit is UNSATISFIABLE. Each is paired with a positive control on the
// same vectors, because a negative test that fails for the wrong reason — a
// build error, a gadget panicking on the witness — is indistinguishable from one
// that passes.
//
// Read each case's comment for what it solved before the fix: an assignment that
// was never satisfiable is not a regression test.

// soundnessVote returns a no-round canonical vote over blockHash, long enough
// for a one-byte leading varint.
func soundnessVote(blockHash [32]byte) []byte {
	return makeCanonicalVote(false, 10, blockHash)
}

// witnessHashFor recomputes the witness commitment for a layout the assignment
// carries, rather than for what the signed bytes say. witness.ComputeWitnessHash derives
// roundPresent and blockHash from a real signed message, which is exactly what a
// malicious prover would not do — so the tests below need this to keep the public
// Hash input consistent with the tampered witness.
//
// Layout (witness/witness.go): PrefixHead(11) || Round(8) || BlockHash(32)
// || per slot: active(1) || A(32).
func witnessHashFor(
	prefixHead [witness.PrefixHeadLen]byte,
	round [witness.RoundLen]byte,
	blockHash [32]byte,
	pubkeys [][]byte,
	activeBytes []byte,
) [32]byte {
	buf := make([]byte, 0, witness.PrefixHeadLen+witness.RoundLen+32+len(pubkeys)*(1+32))
	buf = append(buf, prefixHead[:]...)
	buf = append(buf, round[:]...)
	buf = append(buf, blockHash[:]...)
	for i := range pubkeys {
		buf = append(buf, activeBytes[i])
		buf = append(buf, pubkeys[i]...)
	}
	return sha256.Sum256(buf)
}

// signedBatch returns a bucket-sized batch with one active signer over vote,
// signing signedPrefix (which may be a strict prefix of vote — see the MsgLens
// test). It returns the padded slots so callers can read back the dummy pubkeys.
func signedBatch(t *testing.T, n int, vote, signedPrefix []byte) []witness.ValidatorSignature {
	t.Helper()

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	padded, err := witness.PadWithDummies([]witness.ValidatorSignature{{
		Signature:   ed25519.Sign(priv, signedPrefix),
		PublicKey:   pub,
		SignedBytes: vote,
		Active:      true,
	}}, n)
	if err != nil {
		t.Fatalf("witness.PadWithDummies: %v", err)
	}
	return padded
}

func slotBytes(sigs []witness.ValidatorSignature) (pubkeys [][]byte, activeBytes []byte) {
	pubkeys = make([][]byte, len(sigs))
	activeBytes = make([]byte, len(sigs))
	for i, s := range sigs {
		pubkeys[i] = s.PublicKey
		if s.Active {
			activeBytes[i] = 1
		}
	}
	return pubkeys, activeBytes
}

// TestBatchCircuit_RejectsUnsignedBlockHashBytes is the ZK-05 regression for the
// message-length bound.
//
// The prefix and block-hash asserts read fixed offsets — up to 57 bytes into the
// message — while H_RAM covers only Msgs[i][:MsgLens[i]]. With no lower bound on
// MsgLens, a prover signs a short prefix and fills the rest of the buffer with
// whatever the committed block hash needs to be. Here the validator genuinely
// signs 20 bytes; bytes 20..47, which carry most of the block hash the circuit
// compares against, are never covered by the signature. Every other constraint
// is satisfied: the signature verifies over the real 20-byte prefix, the prefix
// head sits inside it, and the block-hash window matches by construction.
//
// This solved before MsgLens was bounded below. It is the whole finding: the
// committed block hash — which is what carries the AppHash the bridge trusts —
// could come from bytes no validator ever signed.
func TestBatchCircuit_RejectsUnsignedBlockHashBytes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}
	const n = 4
	const signedLen = 20

	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0x50 + i)
	}
	vote := soundnessVote(blockHash)

	// Positive control: the same signer over the whole vote must solve.
	honest := signedBatch(t, n, vote, vote)
	honestHash, err := witness.ComputeWitnessHash(honest)
	if err != nil {
		t.Fatalf("witness.ComputeWitnessHash: %v", err)
	}
	honestAssignment, err := buildBatchAssignment(honest, honestHash)
	if err != nil {
		t.Fatalf("buildBatchAssignment: %v", err)
	}
	if err := test.IsSolved(NewBatchCircuit(n), honestAssignment, ecc.BN254.ScalarField()); err != nil {
		t.Fatalf("honest assignment rejected: %v", err)
	}

	// Attack: sign only the first 20 bytes, present the full buffer, and declare
	// the short length. The bytes the block-hash assert reads past offset 20 are
	// then attacker-supplied.
	forged := signedBatch(t, n, vote, vote[:signedLen])
	forgedHash, err := witness.ComputeWitnessHash(forged)
	if err != nil {
		t.Fatalf("witness.ComputeWitnessHash: %v", err)
	}
	assignment, err := buildBatchAssignment(forged, forgedHash)
	if err != nil {
		t.Fatalf("buildBatchAssignment: %v", err)
	}
	assignment.MsgLens[0] = signedLen

	if err := test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField()); err == nil {
		t.Fatal("circuit accepted a block hash read from bytes outside the signed prefix (ZK-05)")
	}
}

// TestBatchCircuit_RejectsMismatchedRoundPresent is the ZK-05 regression for the
// round-presence tie.
//
// Round-presence selects which 32-byte window of the message is compared against
// the committed block hash: offset 15 when the round field is absent, 24 when it
// is present. Left untied to the message, it is a free witness — so a prover
// takes a real vote with no round, commits a non-zero Round, and commits to
// whatever 32 bytes happen to sit at the round-present offset (here, inside the
// timestamp and chain id). The signature still verifies, the prefix head still
// matches, and the on-chain SharedBlock then claims a round the validators never
// voted in, over a block hash lifted from unrelated bytes.
//
// This solved before the body byte at offset 11 was checked against the round
// tag. Now it is 0x22 (BlockID) where a non-zero Round demands 0x19 (round).
// Presence is derived from the committed Round value, so the converse is also
// covered: a vote that does carry a round, committed with Round = 0.
func TestBatchCircuit_RejectsMismatchedRoundPresent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}
	const n = 4

	t.Run("round committed for a vote without one", func(t *testing.T) {
		var blockHash [32]byte
		for i := range blockHash {
			blockHash[i] = byte(0x60 + i)
		}
		vote := soundnessVote(blockHash)
		padded := signedBatch(t, n, vote, vote)

		// The window a round-present layout would read, in a message that has no
		// round field. One-byte leading varint, so body offset O is at index 1+O.
		var shiftedWindow [32]byte
		copy(shiftedWindow[:], vote[1+witness.BlockHashBodyOffRound:1+witness.BlockHashBodyOffRound+32])

		prefixHead, round, _, err := witness.CommonPrefixFields(padded)
		if err != nil {
			t.Fatalf("CommonPrefixFields: %v", err)
		}
		if round != ([witness.RoundLen]byte{}) {
			t.Fatal("test vector was built without a round field")
		}

		forgedRound := [witness.RoundLen]byte{1}
		pubkeys, activeBytes := slotBytes(padded)
		forgedHash := witnessHashFor(prefixHead, forgedRound, shiftedWindow, pubkeys, activeBytes)

		assignment, err := buildBatchAssignment(padded, forgedHash)
		if err != nil {
			t.Fatalf("buildBatchAssignment: %v", err)
		}
		for j := range forgedRound {
			assignment.Round[j] = uints.NewU8(forgedRound[j])
		}
		for j := 0; j < 32; j++ {
			assignment.BlockHash[j] = uints.NewU8(shiftedWindow[j])
		}

		if err := test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField()); err == nil {
			t.Fatal("circuit accepted a committed Round that contradicts the signed message (ZK-05)")
		}
	})

	t.Run("round omitted for a vote that has one", func(t *testing.T) {
		var blockHash [32]byte
		for i := range blockHash {
			blockHash[i] = byte(0x68 + i)
		}
		vote := makeCanonicalVoteAtRound(3, 10, blockHash)
		padded := signedBatch(t, n, vote, vote)

		prefixHead, _, _, err := witness.CommonPrefixFields(padded)
		if err != nil {
			t.Fatalf("CommonPrefixFields: %v", err)
		}

		// Commit Round = 0 and the no-round window: the committed block hash is
		// then read from inside the round field's shifted BlockID.
		var noRoundWindow [32]byte
		copy(noRoundWindow[:], vote[1+witness.BlockHashBodyOffNoRound:1+witness.BlockHashBodyOffNoRound+32])
		pubkeys, activeBytes := slotBytes(padded)
		forgedHash := witnessHashFor(prefixHead, [witness.RoundLen]byte{}, noRoundWindow, pubkeys, activeBytes)

		assignment, err := buildBatchAssignment(padded, forgedHash)
		if err != nil {
			t.Fatalf("buildBatchAssignment: %v", err)
		}
		for j := 0; j < witness.RoundLen; j++ {
			assignment.Round[j] = uints.NewU8(0)
		}
		for j := 0; j < 32; j++ {
			assignment.BlockHash[j] = uints.NewU8(noRoundWindow[j])
		}

		if err := test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField()); err == nil {
			t.Fatal("circuit accepted Round = 0 for a vote whose round field is present (ZK-05)")
		}
	})
}

// twoSignerBatch returns a bucket-sized batch with two active signers (distinct
// keys, so a quorum of two real validators) over the given votes, padded with
// dummies. Each signer genuinely signs its own vote.
func twoSignerBatch(t *testing.T, n int, voteA, voteB []byte) []witness.ValidatorSignature {
	t.Helper()
	sigs := make([]witness.ValidatorSignature, 0, n)
	for _, vote := range [][]byte{voteA, voteB} {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		sigs = append(sigs, witness.ValidatorSignature{
			Signature:   ed25519.Sign(priv, vote),
			PublicKey:   pub,
			SignedBytes: vote,
			Active:      true,
		})
	}
	padded, err := witness.PadWithDummies(sigs, n)
	if err != nil {
		t.Fatalf("witness.PadWithDummies: %v", err)
	}
	return padded
}

// TestBatchCircuit_RejectsCrossRoundSigners is the regression for committing
// only round PRESENCE.
//
// Each slot's Ed25519 verify binds that slot's own signed bytes, including its
// round — but nothing tied the slots to one round. Two validators who precommit
// the SAME block id in different rounds (round 1 and round 2) therefore verified
// under one proof, with the witness commitment only saying "a round field is
// present". CometBFT's light-client rule rebuilds every signature from the
// commit's single Round, so such a mix is not a commit for any round.
//
// Both signatures are valid, both votes carry the same Type, Height and block
// hash, and both have the round field present, so the prefix-head, block-hash
// and ZK-05 tag asserts are all satisfied. Only the per-slot comparison against
// the committed Round bytes can reject it.
func TestBatchCircuit_RejectsCrossRoundSigners(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}
	const n = 4

	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0x90 + i)
	}

	// Positive control: two signers in the SAME round solve.
	sameRound := twoSignerBatch(t, n,
		makeCanonicalVoteAtRound(2, 10, blockHash),
		makeCanonicalVoteAtRound(2, 90, blockHash), // also mixes varint widths
	)
	sameHash, err := witness.ComputeWitnessHash(sameRound)
	if err != nil {
		t.Fatalf("witness.ComputeWitnessHash: %v", err)
	}
	sameAssignment, err := buildBatchAssignment(sameRound, sameHash)
	if err != nil {
		t.Fatalf("buildBatchAssignment: %v", err)
	}
	if err := test.IsSolved(NewBatchCircuit(n), sameAssignment, ecc.BN254.ScalarField()); err != nil {
		t.Fatalf("same-round assignment rejected: %v", err)
	}

	// Attack: same block id, rounds 1 and 2. The honest-looking commitment takes
	// the first active signer's round.
	mixed := twoSignerBatch(t, n,
		makeCanonicalVoteAtRound(1, 10, blockHash),
		makeCanonicalVoteAtRound(2, 10, blockHash),
	)
	mixedHash, err := witness.ComputeWitnessHash(mixed)
	if err != nil {
		t.Fatalf("witness.ComputeWitnessHash: %v", err)
	}
	mixedAssignment, err := buildBatchAssignment(mixed, mixedHash)
	if err != nil {
		t.Fatalf("buildBatchAssignment: %v", err)
	}
	if err := test.IsSolved(NewBatchCircuit(n), mixedAssignment, ecc.BN254.ScalarField()); err == nil {
		t.Fatal("circuit accepted precommits for one block id from different rounds in one proof")
	}
}

// TestBatchCircuit_RejectsMismatchedRound flips one byte of the committed Round
// (same width) over a real signed vote and expects the circuit to fail. It is
// the "committed round differs from the signed round" shape of the cross-round
// finding: the witness hash is recomputed for the tampered Round, so the public
// input stays consistent and only the per-slot Round assert can catch it.
//
// The positive control commits the true Round over the same vote and must solve;
// without it this test could pass for an unrelated reason.
func TestBatchCircuit_RejectsMismatchedRound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}
	const n = 4

	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0xB0 + i)
	}
	vote := makeCanonicalVoteAtRound(0x0102030405060708, 10, blockHash)
	padded := signedBatch(t, n, vote, vote)

	prefixHead, round, committedHash, err := witness.CommonPrefixFields(padded)
	if err != nil {
		t.Fatalf("CommonPrefixFields: %v", err)
	}
	pubkeys, activeBytes := slotBytes(padded)

	solve := func(round [witness.RoundLen]byte) error {
		hash := witnessHashFor(prefixHead, round, committedHash, pubkeys, activeBytes)
		assignment, err := buildBatchAssignment(padded, hash)
		if err != nil {
			t.Fatalf("buildBatchAssignment: %v", err)
		}
		for j := range round {
			assignment.Round[j] = uints.NewU8(round[j])
		}
		return test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField())
	}

	if err := solve(round); err != nil {
		t.Fatalf("honest Round rejected: %v", err)
	}

	for k := 0; k < witness.RoundLen; k++ {
		tampered := round
		tampered[k] ^= 0x01
		if err := solve(tampered); err == nil {
			t.Fatalf("circuit accepted a committed Round with byte %d flipped", k)
		}
	}
}

// TestBatchCircuit_RejectsNonBooleanActive is the ZK-09 regression.
//
// Active[i] masks every per-slot assert and is committed as one byte of the
// witness hash. varToBytesBE range-checks it to 8 bits, which does not make it a
// bit: at Active[i] = 2 the masked equalities still hold (both sides scale by
// the same factor) and the commitment is simply consistent with the byte 0x02.
// Nothing downstream is broken by that today — the on-chain side encodes active
// as a bool, so it would rebuild 0x00 or 0x01 and the hashes would not match —
// which is why the audit rated it informational. It is still a witness the
// circuit should not admit, and "the mask is a bit" is assumed throughout
// Define.
//
// The vote must be long enough that the Active-scaled ZK-05 length bound stays
// satisfiable: Active = 2 doubles `required` (2 × 48 for a no-round vote), so
// MsgLens must be >= 96 or the length check rejects the witness before the
// boolean assert is ever consulted — which is exactly the false-confidence
// trap this test fell into with the 58-byte soundnessVote (it kept passing
// with the assert removed). At 108 bytes only AssertIsBoolean can reject;
// verified by stripping the assert and watching this witness solve.
func TestBatchCircuit_RejectsNonBooleanActive(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}
	const n = 4

	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0x70 + i)
	}
	// 60 trailing bytes -> 108-byte vote: absorbs the doubled length bound
	// (108 - 96 = 12, inside the 8-bit slack) so the boolean assert is the
	// only constraint the Active[0] = 2 witness can fail.
	vote := makeCanonicalVote(false, 60, blockHash)
	padded := signedBatch(t, n, vote, vote)

	prefixHead, round, committedHash, err := witness.CommonPrefixFields(padded)
	if err != nil {
		t.Fatalf("CommonPrefixFields: %v", err)
	}

	pubkeys, activeBytes := slotBytes(padded)
	activeBytes[0] = 2
	forgedHash := witnessHashFor(prefixHead, round, committedHash, pubkeys, activeBytes)

	assignment, err := buildBatchAssignment(padded, forgedHash)
	if err != nil {
		t.Fatalf("buildBatchAssignment: %v", err)
	}
	assignment.Active[0] = 2

	if err := test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField()); err == nil {
		t.Fatal("circuit accepted a non-boolean Active flag (ZK-09)")
	}
}

// TestBatchCircuit_RejectsTamperedSignature is the liveness check on the batch
// equation itself.
//
// GF-09 moved the batch from the cofactorless S*G == R + H*A to the cofactored
// [8]Q == [8*sum(z_i*S_i)]G. Multiplying both sides by the cofactor is exactly
// the kind of change that can turn an equation vacuous if it lands wrong — a
// check that accepts everything passes every completeness test in the suite. So
// this asserts the equation still rejects a forged S, which is the cheapest
// witness that only the signature check can catch: it does not touch the
// message, the pubkey, the commitment or any of the offset asserts.
func TestBatchCircuit_RejectsTamperedSignature(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}
	const n = 4

	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0x80 + i)
	}
	vote := soundnessVote(blockHash)
	padded := signedBatch(t, n, vote, vote)

	hash, err := witness.ComputeWitnessHash(padded)
	if err != nil {
		t.Fatalf("witness.ComputeWitnessHash: %v", err)
	}
	if _, err := buildBatchAssignment(padded, hash); err != nil {
		t.Fatalf("buildBatchAssignment: %v", err)
	}

	// Flip a low bit of S. The scalar stays canonical, so the assignment builds;
	// only the group equation can reject it.
	tampered := make([]witness.ValidatorSignature, len(padded))
	copy(tampered, padded)
	sig := make([]byte, len(padded[0].Signature))
	copy(sig, padded[0].Signature)
	sig[32] ^= 0x01
	tampered[0].Signature = sig

	assignment, err := buildBatchAssignment(tampered, hash)
	if err != nil {
		t.Fatalf("buildBatchAssignment (tampered): %v", err)
	}

	if err := test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField()); err == nil {
		t.Fatal("cofactored batch equation accepted a tampered signature (GF-09)")
	}
}

// nonceSig signs msg with the nonce r = ±a, a being the signer's secret scalar, so R = rB
// is A itself (negate false) or -A (negate true). It is an ordinary valid Ed25519
// signature (S = r + k*a, k = H(R || A || msg)) with the one property that makes the batch
// verifier's per-row ECIP check degenerate: R and A are the same point up to sign.
func nonceSig(t *testing.T, seedByte byte, msg []byte, negate bool) (pub, sig []byte) {
	t.Helper()
	seed := make([]byte, ed25519.SeedSize)
	seed[0] = seedByte
	h := sha512.Sum512(seed)
	a, err := edwards25519.NewScalar().SetBytesWithClamping(h[:32])
	if err != nil {
		t.Fatalf("clamp secret scalar: %v", err)
	}
	A := new(edwards25519.Point).ScalarBaseMult(a)
	r := a
	if negate {
		r = new(edwards25519.Scalar).Negate(a)
	}
	R := new(edwards25519.Point).ScalarBaseMult(r)
	kh := sha512.New()
	kh.Write(R.Bytes())
	kh.Write(A.Bytes())
	kh.Write(msg)
	k, err := edwards25519.NewScalar().SetUniformBytes(kh.Sum(nil))
	if err != nil {
		t.Fatalf("challenge scalar: %v", err)
	}
	S := new(edwards25519.Scalar).MultiplyAdd(k, a, r)
	pub, sig = A.Bytes(), append(R.Bytes(), S.Bytes()...)
	if !ed25519.Verify(pub, msg, sig) {
		t.Fatalf("crafted signature (negate=%v) does not verify", negate)
	}
	if negate == (string(sig[:32]) == string(pub)) {
		t.Fatalf("crafted signature has the wrong R (negate=%v)", negate)
	}
	return pub, sig
}

// TestBatchCircuit_AcceptsRIsPlusMinusA is the decentrio/ecip-gnark#26 regression on the
// Spectre side: a batch whose signers have R = A or R = -A is valid Ed25519 and must still
// be provable. The ECIP soundness fix must not reject it, or a validator whose signature
// happens to take that shape would stall the light client.
func TestBatchCircuit_AcceptsRIsPlusMinusA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}
	const n = 4

	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0x70 + i)
	}
	vote := soundnessVote(blockHash)

	plainPub, plainPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	plusPub, plusSig := nonceSig(t, 1, vote, false)
	minusPub, minusSig := nonceSig(t, 2, vote, true)

	padded, err := witness.PadWithDummies([]witness.ValidatorSignature{
		{Signature: ed25519.Sign(plainPriv, vote), PublicKey: plainPub, SignedBytes: vote, Active: true},
		{Signature: plusSig, PublicKey: plusPub, SignedBytes: vote, Active: true},
		{Signature: minusSig, PublicKey: minusPub, SignedBytes: vote, Active: true},
	}, n)
	if err != nil {
		t.Fatalf("witness.PadWithDummies: %v", err)
	}
	hash, err := witness.ComputeWitnessHash(padded)
	if err != nil {
		t.Fatalf("witness.ComputeWitnessHash: %v", err)
	}
	assignment, err := buildBatchAssignment(padded, hash)
	if err != nil {
		t.Fatalf("buildBatchAssignment: %v", err)
	}
	if err := test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField()); err != nil {
		t.Fatalf("a batch with R = A and R = -A signers was rejected: %v", err)
	}
}
