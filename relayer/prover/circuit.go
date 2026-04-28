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
// for the validator-signed bytes. Tendermint canonical vote bytes (length
// prefix + inner message) sit at ~110-175 bytes for typical chains; 192
// gives headroom for chain ids near MaxChainIDLen and rounds to a clean
// SHA-256 block boundary count.
const MaxMsgLen = 192

// BatchCircuit verifies N Tendermint precommit Ed25519 signatures and commits
// to its entire witness via a single SHA-256 public input. Collapsing the
// public-input vector to 32 field elements keeps the generated Groth16
// verifier under Ethereum's EIP-170 contract size limit (24576 bytes).
//
// Each slot carries the canonical vote bytes the validator actually signed,
// padded to MaxMsgLen and gated by a per-slot length variable. The on-chain
// WrapperVerifier rebuilds those bytes from `(blockHeader, timestamp_i)` via
// Solidity proto encoding, hashes the witness exactly the same way, and
// passes the digest in as the proof's single public input.
//
// N (slice lengths) is fixed per bucket; one circuit/(pk,vk) triple is built
// per bucket in prover.Buckets.
type BatchCircuit[Base, Scalars emulated.FieldParams] struct {
	Hash [32]uints.U8 `gnark:",public"`

	Sig []eddsa.Signature[Base, Scalars] `gnark:",secret"`
	Pub []eddsa.PublicKey[Base, Scalars] `gnark:",secret"`

	Msgs    [][MaxMsgLen]uints.U8 `gnark:",secret"` // per-slot signed bytes (zero-padded)
	MsgLens []frontend.Variable   `gnark:",secret"` // per-slot meaningful prefix length
}

func (c *BatchCircuit[Base, Scalars]) Define(api frontend.API) error {
	baseApi, err := emulated.NewField[Base](api)
	if err != nil {
		return err
	}
	scalarApi, err := emulated.NewField[Scalars](api)
	if err != nil {
		return err
	}
	uapi, err := uints.New[uints.U32](api)
	if err != nil {
		return err
	}

	// 1. Hash all witness bytes and bind to the public commitment.
	//    Layout (must match hash_witness.go and WrapperVerifier._hashWitness):
	//      per slot: R(32) || S(32) || A(32) || msgLen(2 BE) || msg(MaxMsgLen padded)
	//    Including R/S/A bytes in the hash binds the calldata Sig/Pub to the
	//    proof — without it an attacker could keep the same proof but swap
	//    pubkeys in calldata to misattribute voting power.
	var buf []uints.U8
	for i := range c.Sig {
		buf = append(buf, compressEdwardsToLE(api, baseApi, &c.Sig[i].R)...)
		buf = append(buf, scalarToBytesLE(api, scalarApi, &c.Sig[i].S)...)
		buf = append(buf, compressEdwardsToLE(api, baseApi, &c.Pub[i].A)...)
		buf = append(buf, varToBytesBE(api, c.MsgLens[i], 2)...)
		buf = append(buf, c.Msgs[i][:]...)
	}

	h, err := sha2.New(api)
	if err != nil {
		return err
	}
	h.Write(buf)
	digest := h.Sum()
	for i := 0; i < 32; i++ {
		uapi.ByteAssertEq(c.Hash[i], digest[i])
	}

	// 2. Run Ed25519 batch verify over the per-slot signed bytes. SHA-512
	//    truncates to MsgLens[i] via FixedLengthSum so the canonical-vote
	//    prefix is what gets hashed for H_RAM = SHA-512(R || A || msg).
	msgs := make([][]uints.U8, len(c.Sig))
	for i := range c.Sig {
		msgs[i] = c.Msgs[i][:]
	}
	return eddsa.VerifyBatchWithMsgBytes[Base, Scalars](
		api, c.Sig, c.Pub, msgs, c.MsgLens, eddsa.Config{FromWei: false},
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
	}
}
