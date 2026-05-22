package prover

import (
	"crypto/sha256"
	"encoding/binary"
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
// (WrapperVerifier._hashWitness). Any drift fails the in-circuit byte assert
// against the public Hash field.
//
//	per slot i ∈ [0, bucket):
//	  active[1]               // 0x00 padding slot, 0x01 real signer
//	  A[32 LE]                // raw Ed25519 pubkey (compressed)
//	  msgLen[2 BE]            // length of meaningful prefix in msg
//	  msg[MaxMsgLen padded]   // canonical vote bytes, zero-padded right
//
// A is hashed because Solidity uses pubkeys[i] to look up validator voting
// power; without binding it, an attacker could swap calldata pubkeys to
// misattribute the quorum. R/S are deliberately NOT in calldata or hash —
// the Groth16 proof itself binds them via the in-circuit Ed25519 verify, and
// no on-chain logic consumes them, so omitting them saves calldata + hash
// work without weakening security. The active byte gates each slot's
// contribution to the in-circuit ECIP aggregate.

// ComputeWitnessHash serializes the per-slot data into the canonical layout
// and returns the SHA-256 digest. The on-chain WrapperVerifier recomputes the
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

func encodeWitnessBytes(sigs []ValidatorSignature) ([]byte, error) {
	size := len(sigs) * (1 + 32 + 2 + MaxMsgLen)
	buf := make([]byte, 0, size)
	for i, v := range sigs {
		if len(v.PublicKey) != 32 {
			return nil, fmt.Errorf("slot %d: public key length %d, want 32", i, len(v.PublicKey))
		}
		if len(v.SignedBytes) > MaxMsgLen {
			return nil, fmt.Errorf("slot %d: signed bytes length %d exceeds MaxMsgLen=%d", i, len(v.SignedBytes), MaxMsgLen)
		}
		var activeByte byte
		if v.Active {
			activeByte = 1
		}
		buf = append(buf, activeByte)
		buf = append(buf, v.PublicKey...)
		buf = binary.BigEndian.AppendUint16(buf, uint16(len(v.SignedBytes)))
		padded := make([]byte, MaxMsgLen)
		copy(padded, v.SignedBytes)
		buf = append(buf, padded...)
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

// scalarToBytesLE serializes an emulated scalar to 32 little-endian bytes
// (byte 0 = bits 0..7). Ed25519's S is canonically < 2^253, so the top 3 bits
// of byte 31 will be zero for valid signatures.
func scalarToBytesLE[Scalars emulated.FieldParams](
	api frontend.API,
	scalarApi *emulated.Field[Scalars],
	s *emulated.Element[Scalars],
) []uints.U8 {
	bits := scalarApi.ToBits(s)
	out := make([]uints.U8, 32)
	for j := 0; j < 32; j++ {
		var b frontend.Variable = 0
		for k := 0; k < 8; k++ {
			idx := 8*j + k
			if idx < len(bits) {
				b = api.Add(b, api.Mul(bits[idx], 1<<k))
			}
		}
		out[j] = uints.U8{Val: b}
	}
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
