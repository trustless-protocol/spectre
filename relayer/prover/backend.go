package prover

import (
	"os"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	gnarkbackend "github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
)

const (
	proverBackendNative = "native"
	proverBackendICICLE = "icicle"
)

// ProofBackend wraps the Groth16 operations that differ between native gnark
// and accelerated ICICLE proving. Circuit definitions and verification stay
// backend-agnostic.
type ProofBackend interface {
	Name() string
	NewProvingKey(curveID ecc.ID) groth16.ProvingKey
	Setup(r1cs constraint.ConstraintSystem) (groth16.ProvingKey, groth16.VerifyingKey, error)
	Prove(r1cs constraint.ConstraintSystem, pk groth16.ProvingKey, fullWitness witness.Witness, opts ...gnarkbackend.ProverOption) (groth16.Proof, error)
}

type nativeProofBackend struct{}

func (nativeProofBackend) Name() string {
	return proverBackendNative
}

func (nativeProofBackend) NewProvingKey(curveID ecc.ID) groth16.ProvingKey {
	return groth16.NewProvingKey(curveID)
}

func (nativeProofBackend) Setup(r1cs constraint.ConstraintSystem) (groth16.ProvingKey, groth16.VerifyingKey, error) {
	return groth16.Setup(r1cs)
}

func (nativeProofBackend) Prove(r1cs constraint.ConstraintSystem, pk groth16.ProvingKey, fullWitness witness.Witness, opts ...gnarkbackend.ProverOption) (groth16.Proof, error) {
	return groth16.Prove(r1cs, pk, fullWitness, opts...)
}

// GPUProveEnvEnabled reports whether GPU_PROVE is set to a truthy value.
// Accepts: 1, true, yes, on (case-insensitive).
func GPUProveEnvEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GPU_PROVE"))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// NewProofBackend returns the ICICLE backend when useGPU is true, otherwise
// the native gnark backend. Native is the default so hosts without CUDA/ICICLE
// keep working. Selecting ICICLE on a build without -tags=icicle returns a
// clear error.
func NewProofBackend(useGPU bool) (ProofBackend, error) {
	if useGPU {
		return newICICLEProofBackendFromEnv()
	}
	return nativeProofBackend{}, nil
}
