package witness

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// Dummy padding fills the unused slots of a bucket. Every slot in a batch is
// Ed25519-verified in-circuit — padding is not gated — so each dummy carries a
// real signature from a deterministic keypair, and each (R, A) point in the
// batch stays distinct (duplicated points would hit the doubling branch of the
// ECIP divisor). What marks a slot as padding is the committed Active byte
// plus the on-chain quorum check.
//
// Neither string below is checked on-chain: SignatureVerifier hashes only each
// slot's Active byte and pubkey, and the dummy pubkeys travel in calldata. The
// LENGTH of dummyMagic is a circuit constant, though — DummyMsgLen is the
// circuit's SHA-512 MinMsgLen — so changing it changes the R1CS and requires
// re-setup + redeploying every Groth16Verifier_N{N} + setBucket. Keep it 14 bytes.

// dummyMagic prefixes every dummy slot's signed-bytes buffer. Full layout:
// dummyMagic || bucket(2 BE) || slot(2 BE) = 18 bytes.
const dummyMagic = "spectre-dummy"

// dummySeedFmt derives each slot's deterministic keypair seed.
const dummySeedFmt = "spectre-dummy-seed|bucket=%d|slot=%d"

// DummyMsgLen is the length of every padding slot's signed-bytes buffer — the
// shortest message any slot carries, hence the floor the circuit hands SHA-512
// as minimal input length (ZK-05).
const DummyMsgLen = len(dummyMagic) + 4

// DummyMsgBytes is the on-chain-reproducible signed-bytes buffer for a padding
// slot: dummyMagic || bucket(2 BE) || slot(2 BE).
func DummyMsgBytes(bucket, slot uint16) []byte {
	out := make([]byte, 0, DummyMsgLen)
	out = append(out, dummyMagic...)
	out = binary.BigEndian.AppendUint16(out, bucket)
	out = binary.BigEndian.AppendUint16(out, slot)
	return out
}

// PadWithDummies extends sigs to exactly n slots with deterministic dummy
// signatures (Active=false). Deterministic in (n, slot), so calldata is
// reproducible across runs.
func PadWithDummies(sigs []ValidatorSignature, n int) ([]ValidatorSignature, error) {
	if len(sigs) > n {
		return nil, fmt.Errorf("bucket n=%d too small for %d signers", n, len(sigs))
	}
	if n > 0xFFFF {
		return nil, fmt.Errorf("bucket size %d exceeds 16-bit dummy slot encoding", n)
	}
	out := make([]ValidatorSignature, n)
	copy(out, sigs)
	for i := len(sigs); i < n; i++ {
		seed := sha256.Sum256(fmt.Appendf(nil, dummySeedFmt, n, i))
		priv := ed25519.NewKeyFromSeed(seed[:])
		msg := DummyMsgBytes(uint16(n), uint16(i))
		out[i] = ValidatorSignature{
			Signature:   ed25519.Sign(priv, msg),
			PublicKey:   priv.Public().(ed25519.PublicKey),
			SignedBytes: msg,
			Active:      false,
		}
	}
	return out, nil
}
