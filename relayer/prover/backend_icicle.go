//go:build icicle

package prover

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	gnarkbackend "github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/accelerated/icicle"
	iciclegroth16 "github.com/consensys/gnark/backend/accelerated/icicle/groth16"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
)

type icicleProofBackend struct {
	opts []icicle.Option
}

func newICICLEProofBackendFromEnv() (ProofBackend, error) {
	var opts []icicle.Option
	if raw := strings.TrimSpace(os.Getenv("GNARK_ICICLE_DEVICE_ID")); raw != "" {
		deviceID, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("parse GNARK_ICICLE_DEVICE_ID: %w", err)
		}
		opts = append(opts, icicle.WithDeviceID(deviceID))
	}
	if libs := strings.TrimSpace(os.Getenv("GNARK_ICICLE_BACKEND_LIBS")); libs != "" {
		opts = append(opts, icicle.WithBackendLibrary(libs))
	}
	if raw := strings.TrimSpace(os.Getenv("GNARK_ICICLE_PIN_KEYS")); raw != "" {
		pin, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("parse GNARK_ICICLE_PIN_KEYS: %w", err)
		}
		opts = append(opts, icicle.WithPinKeysToGPU(pin))
	}
	return icicleProofBackend{opts: opts}, nil
}

func (b icicleProofBackend) Name() string {
	return proverBackendICICLE
}

func (b icicleProofBackend) NewProvingKey(curveID ecc.ID) groth16.ProvingKey {
	return iciclegroth16.NewProvingKey(curveID)
}

func (b icicleProofBackend) Setup(r1cs constraint.ConstraintSystem) (groth16.ProvingKey, groth16.VerifyingKey, error) {
	return iciclegroth16.Setup(r1cs)
}

func (b icicleProofBackend) Prove(r1cs constraint.ConstraintSystem, pk groth16.ProvingKey, fullWitness witness.Witness, opts ...gnarkbackend.ProverOption) (groth16.Proof, error) {
	icicleOpts := append([]icicle.Option{}, b.opts...)
	if len(opts) > 0 {
		icicleOpts = append(icicleOpts, icicle.WithProverOptions(opts...))
	}
	return iciclegroth16.Prove(r1cs, pk, fullWitness, icicleOpts...)
}
