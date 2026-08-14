package prover

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/uints"
)

const (
	// WitnessDigestBytes is the width of the SHA-256 witness commitment.
	WitnessDigestBytes = sha256.Size
	// WitnessDigestFieldElements is the number of BN254 field elements we use
	// to expose the 32-byte SHA-256 digest as public inputs. We intentionally
	// use two 128-bit limbs instead of one field element: a single BN254 scalar
	// cannot injectively encode a full 256-bit digest without truncation or
	// modular reduction.
	WitnessDigestFieldElements = 2
	WitnessDigestFieldBytes    = WitnessDigestBytes / WitnessDigestFieldElements
)

// Witness byte layout. Same byte sequence is hashed in three places — Go off-
// chain (ComputeWitnessHash), in-circuit (BatchCircuit.Define), and on-chain
// (SignatureVerifier._hashWitness). Any drift fails the in-circuit byte assert
// against the public Hash field.
//
//	PrefixHead[11]           // Type || Height, first 11 bytes of CanonicalVote body
//	roundPresent[1]          // 0x01 when round is encoded, else 0x00
//	BlockHash[32]            // BlockID hash from the canonical vote body
//	per slot i ∈ [0, bucket):
//	  active[1]              // 0x00 padding slot, 0x01 real signer
//	  A[32 LE]               // raw Ed25519 pubkey (compressed)
//
// A is hashed because Solidity uses pubkeys[i] to look up validator voting
// power; without binding it, an attacker could swap calldata pubkeys to
// misattribute the quorum. R/S are deliberately NOT in calldata or hash —
// the Groth16 proof itself binds them via the in-circuit Ed25519 verify, and
// no on-chain logic consumes them, so omitting them saves calldata + hash
// work without weakening security. The active byte marks a slot as a real
// signer; it does not gate the slot's contribution to the in-circuit ECIP
// aggregate (every slot is verified and aggregated), it is what the on-chain
// quorum check reads to decide which slots carry voting power.

// ComputeWitnessHash serializes the per-slot data into the canonical layout
// and returns the SHA-256 digest. The on-chain SignatureVerifier recomputes the
// same hash from calldata and feeds it as the proof's two packed public inputs.
func ComputeWitnessHash(sigs []ValidatorSignature) ([32]byte, error) {
	buf, err := encodeWitnessBytes(sigs)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(buf), nil
}

// DigestPublicInputs splits the 32-byte witness SHA-256 digest into two
// 128-bit big-endian limbs so the circuit can expose it as two BN254 public
// inputs instead of 32 separate byte-sized inputs.
func DigestPublicInputs(hash [WitnessDigestBytes]byte) [WitnessDigestFieldElements]*big.Int {
	var out [WitnessDigestFieldElements]*big.Int
	for i := 0; i < WitnessDigestFieldElements; i++ {
		start := i * WitnessDigestFieldBytes
		end := start + WitnessDigestFieldBytes
		out[i] = new(big.Int).SetBytes(hash[start:end])
	}
	return out
}

// roundFieldTag is the proto wire tag for CanonicalVote.round (field 3,
// wiretype 1 / fixed64). Its presence at body offset prefixHeadLen tells the
// reader whether the round field is encoded (round > 0) — which shifts the
// BlockID (and the block hash within it) 9 bytes later.
const roundFieldTag = 0x19

// commonPrefixFields derives the #199 committed common fields from the first
// active signer's canonical-vote bytes: PrefixHead (Type|Height), the 32-byte
// block hash inside BlockID, and whether the round field is present. All active
// signers of one block share these contents, and the circuit binds each active
// slot's message against them — deriving from a real signed message guarantees
// the off-chain hash matches what the circuit commits.
//
// The leading length varint is 1 byte when bodyLen < 128 and 2 bytes otherwise.
// Its width is per-signer (the per-validator timestamp shifts bodyLen across the
// 128 boundary), so it is detected here from the continuation bit rather than
// assumed; the body starts right after it. Offsets mirror the circuit constants.
func commonPrefixFields(sigs []ValidatorSignature) (prefixHead [prefixHeadLen]byte, blockHash [32]byte, roundPresent bool, err error) {
	for _, s := range sigs {
		if !s.Active {
			continue
		}
		sb := s.SignedBytes
		if len(sb) == 0 {
			return prefixHead, blockHash, false, fmt.Errorf("signed bytes empty")
		}
		// 1-byte varint unless the first byte's continuation bit is set.
		bodyStart := 1
		if sb[0] >= 0x80 {
			bodyStart = 2
		}
		if len(sb) < bodyStart+prefixHeadLen+1 {
			return prefixHead, blockHash, false, fmt.Errorf("signed bytes too short (%d) for prefix head at offset %d", len(sb), bodyStart)
		}
		roundPresent = sb[bodyStart+prefixHeadLen] == roundFieldTag
		off := bodyStart + blockHashBodyOffNoRound
		if roundPresent {
			off = bodyStart + blockHashBodyOffRound
		}
		if len(sb) < off+32 {
			return prefixHead, blockHash, false, fmt.Errorf("signed bytes too short (%d) for block hash at offset %d", len(sb), off)
		}
		copy(prefixHead[:], sb[bodyStart:bodyStart+prefixHeadLen])
		copy(blockHash[:], sb[off:off+32])
		return prefixHead, blockHash, roundPresent, nil
	}
	return prefixHead, blockHash, false, fmt.Errorf("no active signer to derive common prefix")
}

func encodeWitnessBytes(sigs []ValidatorSignature) ([]byte, error) {
	prefixHead, blockHash, roundPresent, err := commonPrefixFields(sigs)
	if err != nil {
		return nil, err
	}
	// Layout must match circuit BatchCircuit.Define and SignatureVerifier._hashWitness:
	//   PrefixHead(11) || roundPresent(1) || BlockHash(32) || per slot: active(1) || A(32)
	buf := make([]byte, 0, prefixHeadLen+1+32+len(sigs)*(1+32))
	buf = append(buf, prefixHead[:]...)
	var roundByte byte
	if roundPresent {
		roundByte = 1
	}
	buf = append(buf, roundByte)
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

// compressEdwardsToLE encodes an Edwards point to Ed25519's 32-byte canonical
// compressed form: Y in bits 0..254 little-endian, sign of X in bit 255.
func compressEdwardsToLE[Base emulated.FieldParams](
	api frontend.API,
	baseApi *emulated.Field[Base],
	p *sw_emulated.AffinePoint[Base],
) []uints.U8 {
	yBits := baseApi.ToBits(&p.Y)
	xBits := baseApi.ToBits(&p.X)
	out := make([]uints.U8, 32)
	for j := 0; j < 31; j++ {
		var b frontend.Variable = 0
		for k := 0; k < 8; k++ {
			idx := j*8 + k
			if idx < len(yBits) {
				b = api.Add(b, api.Mul(yBits[idx], 1<<k))
			}
		}
		out[j] = uints.U8{Val: b}
	}
	var last frontend.Variable = 0
	for j := 0; j < 7; j++ {
		idx := 31*8 + j
		if idx < len(yBits) {
			last = api.Add(last, api.Mul(yBits[idx], 1<<j))
		}
	}
	if len(xBits) > 0 {
		last = api.Add(last, api.Mul(xBits[0], 1<<7))
	}
	out[31] = uints.U8{Val: last}
	return out
}

// varToBytesBE converts a Variable to nBytes big-endian bytes. Used in-circuit
// for integer fields so the hash layout matches Solidity's abi.encodePacked
// default ordering.
func varToBytesBE(api frontend.API, v frontend.Variable, nBytes int) []uints.U8 {
	nBits := nBytes * 8
	bits := api.ToBinary(v, nBits)
	out := make([]uints.U8, nBytes)
	for j := 0; j < nBytes; j++ {
		var b frontend.Variable = 0
		for k := 0; k < 8; k++ {
			bitIdx := nBits - 1 - 8*j - k
			b = api.Add(b, api.Mul(bits[bitIdx], 1<<(7-k)))
		}
		out[j] = uints.U8{Val: b}
	}
	return out
}
