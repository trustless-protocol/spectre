package prover

import (
	"0x5ea000000/ecip-gnark/signature/eddsa"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
)

// Type aliases for Ed25519 field parameters
type Fp25519 = emulated.Curve25519Fp
type Fr25519 = emulated.Curve25519Fr

// PreHashCircuit verifies a single Ed25519 signature with pre-computed hash.
// H = SHA512(R || A || msg) is computed OUTSIDE the circuit and passed as public input.
// The circuit only verifies: [S]*(-G) + [H]*A + R == 0
// Public inputs (24 uint256): R.X(4), R.Y(4), S(4), Hash(4), A.X(4), A.Y(4)
type PreHashCircuit[Base, Scalars emulated.FieldParams] struct {
	Sig  eddsa.Signature[Base, Scalars] `gnark:",public"`
	Hash emulated.Element[Scalars]      `gnark:",public"`
	Pub  eddsa.PublicKey[Base, Scalars]  `gnark:",public"`
}

func (c *PreHashCircuit[Base, Scalars]) Define(api frontend.API) error {
	config := eddsa.Config{
		Hasher:  nil,
		FromWei: false,
	}
	return eddsa.Verify[Base, Scalars](api, c.Sig, c.Hash, c.Pub, config)
}
