package gnark

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// SetupBucket compiles the batch circuit for bucket n, runs the Groth16 setup
// on the given backend (nil = native), and writes outDir/n{N}/{r1cs,pk,vk}.bin
// — the artifact layout New loads. It returns the in-memory objects so callers
// (cmd/setup-circuits, the E2E test) can smoke-test or export without
// re-reading the files.
//
// The setup is single-party (groth16.Setup): fine for dev/test artifacts,
// NOT for production keys — mainnet needs the MPC ceremony (task #15).
func SetupBucket(outDir string, n int, backend ProofBackend) (constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey, error) {
	if backend == nil {
		backend = nativeProofBackend{}
	}
	bucketDir := filepath.Join(outDir, fmt.Sprintf("n%d", n))
	if err := os.MkdirAll(bucketDir, 0o755); err != nil {
		return nil, nil, nil, fmt.Errorf("create %s: %w", bucketDir, err)
	}

	circuit := NewBatchCircuit(n)
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("compile n=%d: %w", n, err)
	}
	pk, vk, err := backend.Setup(cs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("setup n=%d: %w", n, err)
	}

	for name, artifact := range map[string]io.WriterTo{"r1cs.bin": cs, "pk.bin": pk, "vk.bin": vk} {
		if err := writeArtifact(filepath.Join(bucketDir, name), artifact); err != nil {
			return nil, nil, nil, err
		}
	}
	return cs, pk, vk, nil
}

func writeArtifact(path string, value io.WriterTo) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()
	if _, err := value.WriteTo(f); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
