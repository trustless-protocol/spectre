package gnark

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/uints"
)

// compressEdwardsToLE encodes an Edwards point to Ed25519's 32-byte canonical
// compressed form: Y in bits 0..254 little-endian, sign of X in bit 255.
//
// These bytes are what the on-chain quorum check matches against the pinned
// validator set, so the encoding has to determine the point (ZK-06). Two
// properties make it injective, and neither is free:
//
//   - on-curve: the encoding drops X except for its parity, so an unconstrained
//     X could be moved (x+2 keeps the parity) without changing a byte. The
//     Ed25519 gadget constrains both witness points to the twisted Edwards
//     curve, which leaves y exactly the two solutions +-x that the parity bit
//     separates.
//   - canonical limbs: an emulated element may legally be given as v or v+p,
//     and ToBits reads the limbs as given. AssertIsInRange below pins the
//     representation. It is memoized per element and shares its bit
//     decomposition with the compression, so the Ed25519 gadget's own call on
//     the same coordinates does not pay for it twice.
func compressEdwardsToLE[Base emulated.FieldParams](
	api frontend.API,
	baseApi *emulated.Field[Base],
	p *sw_emulated.AffinePoint[Base],
) []uints.U8 {
	baseApi.AssertIsInRange(&p.X)
	baseApi.AssertIsInRange(&p.Y)

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
