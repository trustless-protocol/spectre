package prover

import (
	"0x5ea000000/ecip-gnark/signature/canonvote"
	"0x5ea000000/ecip-gnark/signature/eddsa"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/uints"
)

// Type aliases for Ed25519 field parameters
type Fp25519 = emulated.Curve25519Fp
type Fr25519 = emulated.Curve25519Fr

// BatchCircuit verifies N Tendermint precommit Ed25519 signatures against a
// shared block + per-validator timestamps. The signed bytes are rebuilt
// inside the circuit so the on-chain verifier only consumes ~22 Fr per slot
// (R, S, A limbs + 2 timestamp fields) instead of the full ~256-byte vote.
//
// N (slice lengths) is fixed at compile time — one circuit per bucket.
type BatchCircuit[Base, Scalars emulated.FieldParams] struct {
	Sig []eddsa.Signature[Base, Scalars] `gnark:",public"`
	Pub []eddsa.PublicKey[Base, Scalars] `gnark:",public"`

	TsSeconds []frontend.Variable `gnark:",public"`
	TsNanos   []frontend.Variable `gnark:",public"`

	Height       frontend.Variable `gnark:",public"`
	Round        frontend.Variable `gnark:",public"`
	BlockIDHash  [32]uints.U8      `gnark:",public"`
	PartSetTotal frontend.Variable `gnark:",public"`
	PartSetHash  [32]uints.U8      `gnark:",public"`
	ChainID      [canonvote.MaxChainIDLen]uints.U8 `gnark:",public"`
	ChainIDLen   frontend.Variable `gnark:",public"`
}

func (c *BatchCircuit[Base, Scalars]) Define(api frontend.API) error {
	shared := canonvote.SharedBlockData{
		Height:       c.Height,
		Round:        c.Round,
		BlockIDHash:  c.BlockIDHash[:],
		PartSetTotal: c.PartSetTotal,
		PartSetHash:  c.PartSetHash[:],
		ChainID:      c.ChainID[:],
		ChainIDLen:   c.ChainIDLen,
	}
	ts := make([]canonvote.Timestamp, len(c.TsSeconds))
	for i := range ts {
		ts[i] = canonvote.Timestamp{Seconds: c.TsSeconds[i], Nanos: c.TsNanos[i]}
	}
	return eddsa.VerifyBatchWithCanonicalVote[Base, Scalars](
		api, c.Sig, c.Pub, shared, ts, eddsa.Config{FromWei: false},
	)
}

// NewBatchCircuit returns a circuit skeleton sized to bucket N. All dynamic
// slices are allocated at length N so frontend.Compile produces a circuit of
// the exact bucket shape.
func NewBatchCircuit(n int) *BatchCircuit[Fp25519, Fr25519] {
	return &BatchCircuit[Fp25519, Fr25519]{
		Sig:       make([]eddsa.Signature[Fp25519, Fr25519], n),
		Pub:       make([]eddsa.PublicKey[Fp25519, Fr25519], n),
		TsSeconds: make([]frontend.Variable, n),
		TsNanos:   make([]frontend.Variable, n),
	}
}
