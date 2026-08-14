package prover

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"testing"

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
// carries, rather than for what the signed bytes say. ComputeWitnessHash derives
// roundPresent and blockHash from a real signed message, which is exactly what a
// malicious prover would not do — so the tests below need this to keep the public
// Hash input consistent with the tampered witness.
//
// Layout (hash_witness.go): PrefixHead(11) || roundPresent(1) || BlockHash(32)
// || per slot: active(1) || A(32).
func witnessHashFor(
	prefixHead [prefixHeadLen]byte,
	roundPresent byte,
	blockHash [32]byte,
	pubkeys [][]byte,
	activeBytes []byte,
) [32]byte {
	buf := make([]byte, 0, prefixHeadLen+1+32+len(pubkeys)*(1+32))
	buf = append(buf, prefixHead[:]...)
	buf = append(buf, roundPresent)
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
func signedBatch(t *testing.T, n int, vote, signedPrefix []byte) []ValidatorSignature {
	t.Helper()

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	padded, err := padWithDummies([]ValidatorSignature{{
		Signature:   ed25519.Sign(priv, signedPrefix),
		PublicKey:   pub,
		SignedBytes: vote,
		Active:      true,
	}}, generateDummySlots(n))
	if err != nil {
		t.Fatalf("padWithDummies: %v", err)
	}
	return padded
}

func slotBytes(sigs []ValidatorSignature) (pubkeys [][]byte, activeBytes []byte) {
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
	honestHash, err := ComputeWitnessHash(honest)
	if err != nil {
		t.Fatalf("ComputeWitnessHash: %v", err)
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
	forgedHash, err := ComputeWitnessHash(forged)
	if err != nil {
		t.Fatalf("ComputeWitnessHash: %v", err)
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
// RoundPresent selects which 32-byte window of the message is compared against
// the committed block hash: offset 15 when the round field is absent, 24 when it
// is present. Left untied to the message, it is a free witness — so a prover
// takes a real vote with no round, declares RoundPresent = 1, and commits to
// whatever 32 bytes happen to sit at the round-present offset (here, inside the
// timestamp and chain id). The signature still verifies, the prefix head still
// matches, and the on-chain SharedBlock then claims a round the validators never
// voted in, over a block hash lifted from unrelated bytes.
//
// This solved before the body byte at offset 11 was checked against the round
// tag. Now it is 0x22 (BlockID) where RoundPresent = 1 demands 0x19 (round).
func TestBatchCircuit_RejectsMismatchedRoundPresent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in-circuit solve in -short mode")
	}
	const n = 4

	var blockHash [32]byte
	for i := range blockHash {
		blockHash[i] = byte(0x60 + i)
	}
	vote := soundnessVote(blockHash)
	padded := signedBatch(t, n, vote, vote)

	// The window a round-present layout would read, in a message that has no
	// round field. One-byte leading varint, so body offset O is at index 1+O.
	var shiftedWindow [32]byte
	copy(shiftedWindow[:], vote[1+blockHashBodyOffRound:1+blockHashBodyOffRound+32])

	prefixHead, _, roundPresent, err := commonPrefixFields(padded)
	if err != nil {
		t.Fatalf("commonPrefixFields: %v", err)
	}
	if roundPresent {
		t.Fatal("test vector was built without a round field")
	}

	pubkeys, activeBytes := slotBytes(padded)
	forgedHash := witnessHashFor(prefixHead, 1, shiftedWindow, pubkeys, activeBytes)

	assignment, err := buildBatchAssignment(padded, forgedHash)
	if err != nil {
		t.Fatalf("buildBatchAssignment: %v", err)
	}
	assignment.RoundPresent = 1
	for j := 0; j < 32; j++ {
		assignment.BlockHash[j] = uints.NewU8(shiftedWindow[j])
	}

	if err := test.IsSolved(NewBatchCircuit(n), assignment, ecc.BN254.ScalarField()); err == nil {
		t.Fatal("circuit accepted RoundPresent that contradicts the signed message (ZK-05)")
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

	prefixHead, committedHash, _, err := commonPrefixFields(padded)
	if err != nil {
		t.Fatalf("commonPrefixFields: %v", err)
	}

	pubkeys, activeBytes := slotBytes(padded)
	activeBytes[0] = 2
	forgedHash := witnessHashFor(prefixHead, 0, committedHash, pubkeys, activeBytes)

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

	hash, err := ComputeWitnessHash(padded)
	if err != nil {
		t.Fatalf("ComputeWitnessHash: %v", err)
	}
	if _, err := buildBatchAssignment(padded, hash); err != nil {
		t.Fatalf("buildBatchAssignment: %v", err)
	}

	// Flip a low bit of S. The scalar stays canonical, so the assignment builds;
	// only the group equation can reject it.
	tampered := make([]ValidatorSignature, len(padded))
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
