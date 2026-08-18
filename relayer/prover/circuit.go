package prover

import (
	"0x5ea000000/ecip-gnark/signature/eddsa"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/sha2"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/uints"
)

// Type aliases for Ed25519 field parameters
type Fp25519 = emulated.Curve25519Fp
type Fr25519 = emulated.Curve25519Fr

// MaxMsgLen is the fixed-width buffer length the circuit allocates per slot
// for the validator-signed bytes (CometBFT VoteSignBytes: length prefix +
// CanonicalVote body). It is chosen from two hard bounds, both asserted by
// TestMaxMsgLenBounds — do not raise it casually.
//
// Upper bound (cost). The signed bytes are hashed by SHA-512 for H_RAM, never
// by the SHA-256 witness commitment, so what matters is the 128-byte SHA-512
// block — not a 64-byte SHA-256 one. FixedLengthSum hashes 64 bytes of R||A
// plus this buffer, then pads:
//
//	blocks = ceil((64 + MaxMsgLen + padding) / 128), padding >= 17
//
// 175 pads 239 -> 256 = 2 blocks; 176 pads 240 -> 384 = 3 blocks. Every slot
// pays that third permutation, so crossing 175 costs ~17% of the circuit
// (N=16: 1,569,995 -> 1,899,694 constraints, 6.03s -> 7.63s per proof on CPU).
//
// Lower bound (correctness). A vote must always fit, or GenerateProof rejects
// the batch (see prover.go). The largest VoteSignBytes CometBFT can produce is
// 167 bytes: a 50-char chain id (types.MaxChainIDLen, enforced at validation),
// max-width height/round encodings and PartSetHeader.Total, a full BlockID,
// and a pathological far-future timestamp — measured by TestMaxMsgLenBounds,
// which rebuilds that vote from the CometBFT types.
//
// 175 is the largest value that keeps SHA-512 at two blocks, leaving 8 bytes
// over the protocol maximum. If a future CometBFT raises MaxChainIDLen past
// 50, the guard test fails and this constant must be revisited — note that
// any value above 175 re-adds the third block.
const MaxMsgLen = 175

// #199 offsets within each signer's canonical-vote message, measured from the
// start of the BODY — i.e. after the leading length varint that
// MarshalDelimited prepends. That varint is 1 byte when bodyLen < 128 and 2
// bytes otherwise; its width is detected PER-SLOT (see Define) rather than
// assumed, because the per-validator timestamp shifts bodyLen across the 128
// boundary, so signers in one commit can mix 1- and 2-byte varints.
//
//	body = Type(0x08 0x02)            // [0,2)
//	     | Height(0x11 + 8B sfixed64) // [2,11)
//	     | Round(0x19 + 8B)           // [11,20) — present only when round > 0
//	     | BlockID(0x22 len 0x0a 0x20 hash[32] ...)
//	     | Timestamp | ChainID
//
// PrefixHead = Type|Height = first 11 body bytes (always present). The 32-byte
// block hash sits 4 bytes into BlockID, so at body offset 15 (no round) or 24
// (round shifts BlockID by 9). The in-message index adds the per-slot varint
// width (1 or 2) to these body offsets.
const (
	prefixHeadLen           = 11 // Type(2) + Height(9)
	blockHashBodyOffNoRound = 15
	blockHashBodyOffRound   = 24 // round shifts BlockID by 9

	// blockIDFieldTag is the proto wire tag for CanonicalVote.block_id (field 4,
	// wiretype 2 / length-delimited). It occupies body offset prefixHeadLen when
	// the round field is absent, where roundFieldTag (hash_witness.go) sits when
	// it is present — so this one byte distinguishes the two body layouts and is
	// what RoundPresent is checked against (ZK-05).
	blockIDFieldTag = 0x22

	// blockHashEnd is the body offset one past the last byte the fixed-offset
	// asserts below read, in the round-present layout — the widest of the two.
	// A signed prefix shorter than this leaves those bytes outside what H_RAM
	// covers, so it is the per-slot lower bound on MsgLens for active slots.
	blockHashEnd = blockHashBodyOffRound + 32
)

// BatchCircuit verifies N Tendermint precommit Ed25519 signatures and commits
// to its entire witness via a single SHA-256 digest, exposed publicly as two
// 128-bit field elements. This keeps the generated Groth16 verifier compact
// without losing the full 256-bit binding of the witness commitment.
//
// Each slot carries the canonical vote bytes the validator actually signed,
// padded to MaxMsgLen and gated by a per-slot length variable. The on-chain
// SignatureVerifier rebuilds those bytes from `(blockHeader, timestamp_i)` via
// Solidity proto encoding, hashes the witness exactly the same way, and
// passes the digest in as two packed public inputs.
//
// N (slice lengths) is fixed per bucket; one circuit/(pk,vk) triple is built
// per bucket in prover.Buckets.
type BatchCircuit[Base, Scalars emulated.FieldParams] struct {
	Hash [WitnessDigestFieldElements]frontend.Variable `gnark:",public"`

	// #199: bind only the security-critical common fields in-circuit, at fixed
	// offsets in every signer's message, instead of hashing the full per-slot
	// canonical vote. PrefixHead = Type(precommit)+Height. BlockHash = the
	// 32-byte block hash inside BlockID (→ binds the header → AppHash, chainId,
	// height). RoundPresent selects the block-hash offset (round shifts BlockID
	// by 9 bytes). All three are committed (hashed) so calldata can't alter
	// them. Round/timestamp/chainID are left unbound here — the Ed25519 verify
	// binds them via H_RAM, and chainID is enforced on-chain.
	PrefixHead   [prefixHeadLen]uints.U8 `gnark:",secret"`
	BlockHash    [32]uints.U8            `gnark:",secret"`
	RoundPresent frontend.Variable       `gnark:",secret"`

	Sig []eddsa.Signature[Base, Scalars] `gnark:",secret"`
	Pub []eddsa.PublicKey[Base, Scalars] `gnark:",secret"`

	Msgs    [][MaxMsgLen]uints.U8 `gnark:",secret"` // per-slot signed bytes (zero-padded)
	MsgLens []frontend.Variable   `gnark:",secret"` // per-slot meaningful prefix length
	// Active marks a slot as a real signer rather than padding. It does NOT
	// gate the slot's contribution to the ECIP aggregate — every slot is
	// verified and aggregated alike, which is why padding slots carry
	// deterministic dummy signatures that genuinely verify. What Active does is
	// (a) mask the in-circuit asserts that only apply to real canonical votes
	// and (b) tell the on-chain quorum check which slots carry voting power.
	// It is bound into the witness hash so calldata cannot toggle a slot's
	// active bit after proving.
	Active []frontend.Variable `gnark:",secret"`
}

func (c *BatchCircuit[Base, Scalars]) Define(api frontend.API) error {
	baseApi, err := emulated.NewField[Base](api)
	if err != nil {
		return err
	}

	// 1. Hash all witness bytes and bind to the public commitment.
	//    Layout (must match hash_witness.go and SignatureVerifier._hashWitness):
	//      PrefixHead(11) || roundPresent(1) || BlockHash(32) || per slot: active(1) || A(32)
	//    Binding A protects the on-chain quorum lookup — Solidity matches
	//    pubkeys[i] against the validator set to attribute voting power, so
	//    without this hash binding an attacker could swap calldata pubkeys.
	//    Hashing the encoding is only half of that guarantee: the encoding must
	//    also determine the point that enters the EC-MSM, which needs A to be
	//    canonically represented and on the curve (ZK-06). compressEdwardsToLE
	//    asserts the representation here; the curve assert lives in the ecip-gnark
	//    Ed25519 gadget (signature/eddsa/verify_msg_bytes.go asserts both R and A
	//    are on the curve, and the batch is cofactored per GF-09). That fork is
	//    pinned since the third_party submodule bump in #398, so both halves of
	//    ZK-06 are in force — see #282.
	//    R/S are NOT hashed and NOT in calldata: the Groth16 proof itself
	//    binds them via the in-circuit Ed25519 verify, and no on-chain logic
	//    consumes them. The active byte is placed first so the on-chain
	//    rebuild can short-circuit cheaply for padding slots if needed.
	//    The per-slot canonical-vote bytes are no longer hashed here (#199) — they
	//    are bound instead by the fixed-offset prefix/blockHash asserts below
	//    + the in-circuit Ed25519 verify. This shrinks the witness commit.
	var buf []uints.U8
	buf = append(buf, c.PrefixHead[:]...)
	buf = append(buf, varToBytesBE(api, c.RoundPresent, 1)...)
	buf = append(buf, c.BlockHash[:]...)
	for i := range c.Sig {
		buf = append(buf, varToBytesBE(api, c.Active[i], 1)...)
		buf = append(buf, compressEdwardsToLE(api, baseApi, &c.Pub[i].A)...)
	}

	h, err := sha2.New(api)
	if err != nil {
		return err
	}
	h.Write(buf)
	digest := h.Sum()
	for i := 0; i < WitnessDigestFieldElements; i++ {
		start := i * WitnessDigestFieldBytes
		end := start + WitnessDigestFieldBytes
		api.AssertIsEqual(c.Hash[i], packBytesToVariableBE(api, digest[start:end]))
	}

	// #199 prefix binding (cheap: fixed offsets + a single RoundPresent select,
	// no per-byte comparator). For each signer's message:
	//   - PrefixHead (Type=precommit | Height) is always at [prefixBodyOffset, +11).
	//   - the 32-byte block hash is at blockHashOffNoRound or blockHashOffRound,
	//     selected by RoundPresent (round, when present, shifts BlockID by 9 B).
	// Binds Type, Height and the block hash (→ header → AppHash/chainId/height).
	// A committed-vs-actual mismatch (wrong RoundPresent, wrong bytes) fails the
	// equality; everything else in the message is covered by the Ed25519 verify.
	// Only ACTIVE slots carry the real canonical vote — padding/dummy slots sign
	// deterministic dummy bytes, so their prefix must NOT be checked against the
	// committed common prefix. Mask each assert by Active[i] (a single bit
	// already in the witness): for active slots it enforces equality, for padding
	// it is vacuous (0 == 0).
	api.AssertIsBoolean(c.RoundPresent)
	for i := range c.Sig {
		// Active[i] is the mask on every assert below and one byte of the witness
		// commitment. varToBytesBE range-checks it to 8 bits, which is not the
		// same thing: any Active[i] in 2..255 scales both sides of every masked
		// equality by the same factor and still passes, with the commitment
		// simply consistent with that byte (0x02, ...). Not currently
		// exploitable — the on-chain side rebuilds active as 0x00/0x01, so the
		// hashes would not match — but "the mask is a bit" is assumed
		// everywhere below (ZK-09).
		api.AssertIsBoolean(c.Active[i])

		// Per-slot leading-varint width: the continuation bit (MSB) of the first
		// length byte is set iff a second varint byte follows (bodyLen >= 128).
		// So the body starts at index 1 (v2=0, 1-byte varint) or 2 (v2=1, 2-byte
		// varint), and every body offset O below maps to in-message index
		// (1+v2)+O. v2 is derived from the message itself; for padding slots it is
		// meaningless but harmless since the asserts are masked by Active[i].
		byte0Bits := api.ToBinary(c.Msgs[i][0].Val, 8)
		v2 := byte0Bits[7]

		for j := 0; j < prefixHeadLen; j++ {
			b := api.Select(v2, c.Msgs[i][2+j].Val, c.Msgs[i][1+j].Val)
			api.AssertIsEqual(api.Mul(c.Active[i], b), api.Mul(c.Active[i], c.PrefixHead[j].Val))
		}
		for j := 0; j < 32; j++ {
			noRound := api.Select(v2, c.Msgs[i][2+blockHashBodyOffNoRound+j].Val, c.Msgs[i][1+blockHashBodyOffNoRound+j].Val)
			round := api.Select(v2, c.Msgs[i][2+blockHashBodyOffRound+j].Val, c.Msgs[i][1+blockHashBodyOffRound+j].Val)
			b := api.Select(c.RoundPresent, round, noRound)
			api.AssertIsEqual(api.Mul(c.Active[i], b), api.Mul(c.Active[i], c.BlockHash[j].Val))
		}

		// ZK-05, part 1: tie RoundPresent to the message it selects an offset in.
		// RoundPresent is a free witness that chooses which 32-byte window is
		// compared against the committed block hash, so left untied the prover
		// picks the window rather than reading it. The body byte at
		// prefixHeadLen settles it from the message itself: proto emits the round
		// field (tag 0x19) there when round > 0 and BlockID (tag 0x22) otherwise,
		// so the tag and the offset must agree.
		tag := api.Select(v2, c.Msgs[i][2+prefixHeadLen].Val, c.Msgs[i][1+prefixHeadLen].Val)
		wantTag := api.Select(c.RoundPresent, roundFieldTag, blockIDFieldTag)
		api.AssertIsEqual(api.Mul(c.Active[i], tag), api.Mul(c.Active[i], wantTag))

		// ZK-05, part 2: keep every byte the asserts above read inside the signed
		// prefix. H_RAM covers only Msgs[i][:MsgLens[i]], while the offsets are
		// fixed and reach as far as blockHashEnd — so without a lower bound on
		// MsgLens the prover supplies attacker-chosen bytes past the prefix and
		// has them compared against the committed block hash, with the signature
		// covering none of it.
		//
		// required = varint width + block-hash end, masked by Active[i]: padding
		// slots carry an 18-byte dummy message and must stay exempt. The slack is
		// range-checked to 8 bits, which forces MsgLens[i] >= required — MsgLens
		// is already bounded above by MaxMsgLen inside FixedLengthSum, so the
		// difference cannot wrap.
		blockHashOff := api.Select(c.RoundPresent, blockHashEnd, blockHashEnd-(blockHashBodyOffRound-blockHashBodyOffNoRound))
		required := api.Add(api.Add(1, v2), blockHashOff)
		api.ToBinary(api.Sub(c.MsgLens[i], api.Mul(c.Active[i], required)), 8)
	}

	// 2. Run Ed25519 batch verify over the per-slot signed bytes. SHA-512
	//    truncates to MsgLens[i] via FixedLengthSum so the canonical-vote
	//    prefix is what gets hashed for H_RAM = SHA-512(R || A || msg).
	msgs := make([][]uints.U8, len(c.Sig))
	for i := range c.Sig {
		msgs[i] = c.Msgs[i][:]
	}
	// All N slots carry valid Ed25519 signatures — real signers sign canonical
	// vote bytes, padding slots sign deterministic dummy bytes via a generated
	// dummy keypair. Both verify under the same primitive, and the batch does
	// NOT gate inactive slots: every slot contributes to the ECIP aggregate and
	// every slot's signature must verify. Nothing needs gating, because each
	// (R, A) is distinct and each signature stands on its own; what distinguishes
	// a real signer from padding is the committed Active byte plus the on-chain
	// quorum check, which is where padding is excluded.
	//
	// MinMsgLen is the floor over ALL slots, so it is the dummy length. The
	// stricter per-slot bound that real votes need is asserted above, masked by
	// Active[i].
	return eddsa.VerifyBatchWithMsgBytes[Base, Scalars](
		api, c.Sig, c.Pub, msgs, c.MsgLens, eddsa.Config{FromWei: false, MinMsgLen: DummyMsgLen},
	)
}

// NewBatchCircuit returns a circuit skeleton sized to bucket N. All dynamic
// slices are allocated at length N so frontend.Compile produces a circuit of
// the exact bucket shape.
func NewBatchCircuit(n int) *BatchCircuit[Fp25519, Fr25519] {
	return &BatchCircuit[Fp25519, Fr25519]{
		Sig:     make([]eddsa.Signature[Fp25519, Fr25519], n),
		Pub:     make([]eddsa.PublicKey[Fp25519, Fr25519], n),
		Msgs:    make([][MaxMsgLen]uints.U8, n),
		MsgLens: make([]frontend.Variable, n),
		Active:  make([]frontend.Variable, n),
	}
}

func packBytesToVariableBE(api frontend.API, bs []uints.U8) frontend.Variable {
	var acc frontend.Variable = 0
	for _, b := range bs {
		acc = api.Add(api.Mul(acc, 256), b.Val)
	}
	return acc
}
