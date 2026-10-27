package witness

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

// This file is the canonical definition of the witness byte contract. The
// exact same byte sequence is hashed in three places — here (off-chain),
// in-circuit (gnark/circuit.go BatchCircuit.Define), and on-chain
// (SignatureVerifier._hashWitness) — and a proof is only accepted when all
// three agree. Changing anything here is a protocol change: new golden
// vectors, circuit re-setup, verifier redeploy.
//
// Layout:
//
//	PrefixHead[11]           // Type || Height, first 11 bytes of CanonicalVote body
//	Round[8]                 // sfixed64 LE round as on the wire; 0x00*8 when the field is absent
//	BlockHash[32]            // BlockID hash from the canonical vote body
//	per slot i ∈ [0, bucket):
//	  active[1]              // 0x00 padding slot, 0x01 real signer
//	  A[32 LE]               // raw Ed25519 pubkey (compressed)
//
// A is hashed because Solidity uses pubkeys[i] to look up validator voting
// power; without binding it, an attacker could swap calldata pubkeys to
// misattribute the quorum. R/S are deliberately NOT in calldata or hash —
// the Groth16 proof itself binds them via the in-circuit Ed25519 verify, and
// no on-chain logic consumes them.
//
// Round is committed by VALUE (not just presence) and asserted equal in every
// active slot's signed bytes. Each Ed25519 verify binds only that slot's own
// round, so without this a single proof could aggregate precommits for one
// block id cast in different rounds — which CometBFT's light-client rule
// (one Commit.Round for every signature) does not allow. Presence is derived
// from the value: proto3 omits the field exactly when round == 0.

const (
	// WitnessDigestBytes is the width of the SHA-256 witness commitment.
	WitnessDigestBytes = sha256.Size
	// WitnessDigestFieldElements is the number of BN254 field elements used to
	// expose the 32-byte digest as public inputs: two 128-bit limbs, because a
	// single BN254 scalar cannot injectively encode 256 bits.
	WitnessDigestFieldElements = 2
	WitnessDigestFieldBytes    = WitnessDigestBytes / WitnessDigestFieldElements
)

// MaxMsgLen is the fixed-width buffer length the circuit allocates per slot
// for the validator-signed bytes (CometBFT VoteSignBytes: length prefix +
// CanonicalVote body). Two hard bounds pin it, both asserted by
// TestMaxMsgLenBounds — do not raise it casually.
//
// Upper bound (cost): FixedLengthSum hashes 64 bytes of R||A plus this buffer
// through SHA-512; 175 pads to 2 blocks, 176 forces a 3rd block on every slot
// (~17% more constraints). Lower bound (correctness): the largest vote
// CometBFT can produce is 167 bytes (50-char chain id per types.MaxChainIDLen,
// max-width integer encodings, pathological timestamp). 175 is the largest
// 2-block value, leaving 8 bytes of headroom.
const MaxMsgLen = 175

// #199 offsets within each signer's canonical-vote message, measured from the
// start of the BODY — i.e. after the leading length varint (1 byte when
// bodyLen < 128, else 2; detected per slot, since the per-validator timestamp
// shifts bodyLen across the 128 boundary).
//
//	body = Type(0x08 0x02)            // [0,2)
//	     | Height(0x11 + 8B sfixed64) // [2,11)
//	     | Round(0x19 + 8B)           // [11,20) — present only when round > 0
//	     | BlockID(0x22 len 0x0a 0x20 hash[32] ...)
//	     | Timestamp | ChainID
const (
	PrefixHeadLen           = 11 // Type(2) + Height(9)
	RoundLen                = 8  // sfixed64 payload; follows the 0x19 tag at body offset PrefixHeadLen
	RoundBodyOff            = PrefixHeadLen + 1
	BlockHashBodyOffNoRound = 15
	BlockHashBodyOffRound   = 24 // round shifts BlockID by 9

	// BlockIDFieldTag is the proto wire tag for CanonicalVote.block_id
	// (field 4, length-delimited). It sits at body offset PrefixHeadLen when
	// the round field is absent, where RoundFieldTag sits when present — one
	// byte distinguishing the two body layouts (ZK-05).
	BlockIDFieldTag = 0x22

	// BlockHashEnd is one past the last body offset the circuit's fixed-offset
	// asserts read (round-present layout, the wider one). It is the per-slot
	// lower bound on MsgLens for active slots.
	BlockHashEnd = BlockHashBodyOffRound + 32
)

// RoundFieldTag is the proto wire tag for CanonicalVote.round (field 3,
// fixed64). Its presence at body offset PrefixHeadLen means round > 0, which
// shifts the BlockID — and the block hash inside it — 9 bytes later. The 8
// round bytes themselves sit at [RoundBodyOff, RoundBodyOff+RoundLen).
const RoundFieldTag = 0x19

// ComputeWitnessHash serializes the batch into the canonical layout above and
// returns its SHA-256 digest — the value the proof's two public inputs commit
// to, and the "what is being proven" half of a ID.
func ComputeWitnessHash(sigs []ValidatorSignature) ([32]byte, error) {
	buf, err := encodeWitnessBytes(sigs)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(buf), nil
}

// DigestPublicInputs splits the witness digest into two 128-bit big-endian
// limbs, matching how the circuit exposes it as BN254 public inputs.
func DigestPublicInputs(hash [WitnessDigestBytes]byte) [WitnessDigestFieldElements]*big.Int {
	var out [WitnessDigestFieldElements]*big.Int
	for i := 0; i < WitnessDigestFieldElements; i++ {
		out[i] = new(big.Int).SetBytes(hash[i*WitnessDigestFieldBytes : (i+1)*WitnessDigestFieldBytes])
	}
	return out
}

// CommonPrefixFields derives the #199 committed common fields from the first
// active signer's canonical-vote bytes: PrefixHead (Type|Height), the 8 round
// bytes exactly as they appear on the wire (all zero when the round field is
// absent), and the 32-byte block hash inside BlockID. All active signers of one
// block share these contents — the circuit asserts it per slot; deriving them
// from a real signed message guarantees the off-chain hash matches what the
// circuit commits.
func CommonPrefixFields(sigs []ValidatorSignature) (prefixHead [PrefixHeadLen]byte, round [RoundLen]byte, blockHash [32]byte, err error) {
	for _, s := range sigs {
		if !s.Active {
			continue
		}
		sb := s.SignedBytes
		if len(sb) == 0 {
			return prefixHead, round, blockHash, fmt.Errorf("signed bytes empty")
		}
		// Leading length varint: 1 byte unless the continuation bit is set.
		bodyStart := 1
		if sb[0] >= 0x80 {
			bodyStart = 2
		}
		if len(sb) < bodyStart+PrefixHeadLen+1 {
			return prefixHead, round, blockHash, fmt.Errorf("signed bytes too short (%d) for prefix head at offset %d", len(sb), bodyStart)
		}
		roundPresent := sb[bodyStart+PrefixHeadLen] == RoundFieldTag
		off := bodyStart + BlockHashBodyOffNoRound
		if roundPresent {
			off = bodyStart + BlockHashBodyOffRound
		}
		if len(sb) < off+32 {
			return prefixHead, round, blockHash, fmt.Errorf("signed bytes too short (%d) for block hash at offset %d", len(sb), off)
		}
		copy(prefixHead[:], sb[bodyStart:bodyStart+PrefixHeadLen])
		if roundPresent {
			copy(round[:], sb[bodyStart+RoundBodyOff:bodyStart+RoundBodyOff+RoundLen])
			// The circuit derives presence from the value (round != 0), so a
			// round field carrying zero would be rejected in-circuit. proto3
			// never emits it; fail here rather than at prove time.
			if round == ([RoundLen]byte{}) {
				return prefixHead, round, blockHash, fmt.Errorf("non-canonical vote: round field present with value 0")
			}
		}
		copy(blockHash[:], sb[off:off+32])
		return prefixHead, round, blockHash, nil
	}
	return prefixHead, round, blockHash, fmt.Errorf("no active signer to derive common prefix")
}

func encodeWitnessBytes(sigs []ValidatorSignature) ([]byte, error) {
	prefixHead, round, blockHash, err := CommonPrefixFields(sigs)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 0, PrefixHeadLen+RoundLen+32+len(sigs)*(1+32))
	buf = append(buf, prefixHead[:]...)
	buf = append(buf, round[:]...)
	buf = append(buf, blockHash[:]...)
	for i, v := range sigs {
		if len(v.PublicKey) != 32 {
			return nil, fmt.Errorf("slot %d: public key length %d, want 32", i, len(v.PublicKey))
		}
		var activeByte byte
		if v.Active {
			activeByte = 1
		}
		buf = append(buf, activeByte)
		buf = append(buf, v.PublicKey...)
	}
	return buf, nil
}
