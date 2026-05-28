package prover

import (
	"fmt"
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

func gpuProveEnvEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GPU_PROVE"))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// NewProofBackendFromSelection selects the proving backend from explicit
// inputs. Native proving is the default so hosts without CUDA/ICICLE keep
// working. useGPU is a shorthand for selecting the ICICLE backend.
func NewProofBackendFromSelection(name string, useGPU bool) (ProofBackend, error) {
	if useGPU {
		return newICICLEProofBackendFromEnv()
	}
	name = strings.ToLower(strings.TrimSpace(name))
	switch name {
	case "", proverBackendNative, "cpu":
		return nativeProofBackend{}, nil
	case proverBackendICICLE, "gpu":
		return newICICLEProofBackendFromEnv()
	default:
		return nil, fmt.Errorf("unsupported GNARK_PROVER_BACKEND=%q (supported: native, icicle)", name)
	}
}

// NewProofBackendFromEnv selects the proving backend from environment
// variables. GPU_PROVE=1 is a shorthand for GNARK_PROVER_BACKEND=icicle.
func NewProofBackendFromEnv() (ProofBackend, error) {
	return NewProofBackendFromSelection(os.Getenv("GNARK_PROVER_BACKEND"), gpuProveEnvEnabled())
}
