//go:build !icicle

package prover

import "fmt"

func newICICLEProofBackendFromEnv() (ProofBackend, error) {
	return nil, fmt.Errorf("GNARK_PROVER_BACKEND=icicle or GPU_PROVE=1 requires building with -tags=icicle and installing the ICICLE runtime")
}
