package prover

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// dummyMagic is the prefix every dummy slot's signed-bytes buffer starts
// with. The full layout is `dummyMagic || bucket(2 BE) || slot(2 BE)` (= 18
// bytes), reproducible byte-for-byte on-chain by SignatureVerifier so the
// in-circuit hash matches without paying calldata for padding messages.
const dummyMagic = "fast-ibc-dummy"

// dummySignature is a deterministic per-slot Ed25519 signature used to fill
// padding slots in a bucketed batch. Padding slots have Active=false so their
// signature does not need to verify against the canonical-vote bytes — but the
// (R, A) coordinates must be distinct from every other slot in the same batch
// so the in-circuit ECIP divisor does not hit the doubling branch with
// duplicated points.
//
// Distinctness is guaranteed by deriving each slot's seed from
// SHA256("fast-ibc-dummy" || bucket || slotIndex) and signing a slot-specific
// message; the chance two seeds collide is 2^-256.
type dummySignature struct {
	signature   []byte // 64 bytes: R || S
	publicKey   []byte // 32 bytes: compressed Ed25519 public key
	signedBytes []byte // dummy message that was signed (length < MaxMsgLen)
}

// generateDummySlots returns one dummy signature per slot in [0, n). Each
// dummy is deterministic in (n, slotIndex) so calldata is reproducible across
// runs — useful for debugging and gas-snapshot stability.
func generateDummySlots(n int) []dummySignature {
	if n > 0xFFFF {
		panic(fmt.Sprintf("bucket size %d exceeds 16-bit dummy slot encoding", n))
	}
	out := make([]dummySignature, n)
	for i := 0; i < n; i++ {
		seedSrc := fmt.Sprintf("fast-ibc-dummy-seed|bucket=%d|slot=%d", n, i)
		seed := sha256.Sum256([]byte(seedSrc))
		priv := ed25519.NewKeyFromSeed(seed[:])
		pub := priv.Public().(ed25519.PublicKey)
		msg := DummyMsgBytes(uint16(n), uint16(i))
		sig := ed25519.Sign(priv, msg)
		out[i] = dummySignature{
			signature:   sig,
			publicKey:   pub,
			signedBytes: msg,
		}
	}
	return out
}

// DummyMsgBytes is the on-chain-reproducible signed-bytes buffer for a padding
// slot. Layout: dummyMagic || bucket(2 BE) || slot(2 BE). SignatureVerifier
// rebuilds the exact same bytes using its own copy of dummyMagic so the
// witness hash matches without requiring calldata for inactive slots.
func DummyMsgBytes(bucket, slot uint16) []byte {
	out := make([]byte, 0, len(dummyMagic)+4)
	out = append(out, []byte(dummyMagic)...)
	out = binary.BigEndian.AppendUint16(out, bucket)
	out = binary.BigEndian.AppendUint16(out, slot)
	return out
}

// asValidatorSignature converts a dummy into a ValidatorSignature suitable for
// padding a real-signer slice up to a bucket boundary. Active=false marks it
// as padding for the in-circuit gate and on-chain quorum check.
func (d dummySignature) asValidatorSignature() ValidatorSignature {
	return ValidatorSignature{
		Signature:   d.signature,
		PublicKey:   d.publicKey,
		SignedBytes: d.signedBytes,
		Active:      false,
	}
}
