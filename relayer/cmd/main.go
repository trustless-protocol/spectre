// main.go is the process entry point and nothing else: wire the commands, run
// one, and shut down cleanly. Each command lives in its own file, following
// submit_misbehaviour.go, which already had that shape.
package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"relayer/prover"
)

func isShutdownErr(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func stopRelaysAndCleanup(cancel func(), wait func(), cleanups []func()) {
	cancel()
	wait()
	for _, cleanup := range cleanups {
		cleanup()
	}
}

// proofBackendFromFlags resolves the GPU/CPU backend from --gpu-prove or the
// GPU_PROVE env var. Returns (backend, true) when an explicit selection was
// made, otherwise (nil, false) so the caller falls back to env-only defaults.
func proofBackendFromFlags(cmd *cobra.Command) (prover.ProofBackend, bool, error) {
	flagSet := cmd.Flags().Changed(flagGPUProve)
	envSet := prover.GPUProveEnvEnabled()
	if !flagSet && !envSet {
		return nil, false, nil
	}

	useGPU := envSet
	if flagSet {
		v, err := cmd.Flags().GetBool(flagGPUProve)
		if err != nil {
			return nil, false, fmt.Errorf("failed to get gpu prove flag: %w", err)
		}
		useGPU = v
	}

	backend, err := prover.NewProofBackend(useGPU)
	if err != nil {
		return nil, false, err
	}
	return backend, true, nil
}

// --- Main ---

func main() {
	zLogger, _ := zap.NewProduction(zap.AddStacktrace(zap.DPanicLevel))
	defer zLogger.Sync()

	logger := zLogger.Sugar()

	rootCmd := &cobra.Command{
		Use:   "relayer [command]",
		Short: "fast-ibc operator — relay IBC packets between Cosmos and Ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(
		Start(zLogger),
		CreateClientsCosmos(zLogger),
		CreateClientsEth(zLogger),
		UpdateClient(zLogger),
		SubmitMisbehaviour(zLogger),
		Genesis(zLogger),
	)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
