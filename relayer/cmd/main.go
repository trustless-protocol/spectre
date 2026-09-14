// main.go is the process entry point and nothing else: wire the commands, run
// one, and shut down cleanly. Each command lives in its own file, following
// submit_misbehaviour.go, which already had that shape.
package main

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"relayer/prover"
	workergroup "relayer/workers"
)

func isShutdownErr(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

// Shutdown budget for the whole process, split between the two phases it has.
//
// shutdownDrainBudget matches workergroup.ShutdownDrainTimeout: modules drain in
// parallel, so the per-module bound is also the bound for all of them together
// and the two must not drift apart. shutdownCleanupBudget covers the client
// close layer (cosmosClient.Stop, ethClient.Close). Together they are the 45
// seconds a SIGTERM is allowed to take.
//
// Widening these is almost always the wrong fix. A drain that overruns means an
// RPC is running without a deadline derived from the relay context, and the
// remedy is to give it one, not to wait longer for it.
// Variables, not constants, only so the shutdown tests can run in milliseconds
// instead of 45 seconds. Production never assigns them.
var (
	shutdownDrainBudget   = workergroup.ShutdownDrainTimeout
	shutdownCleanupBudget = 15 * time.Second
)

// shutdownBudget is the number the operator-facing contract is stated in: a
// SIGTERM finishes, one way or the other, inside this.
const shutdownBudget = workergroup.ShutdownDrainTimeout + 15*time.Second

// stopRelaysAndCleanup cancels the relay context, waits for the engines to
// drain, then closes the clients -- in that order, because closing a client out
// from under a running goroutine turns a clean shutdown into an error cascade.
//
// Each phase is bounded. An unbounded wait here is not a hang that someone
// notices: the process sits there holding its Cosmos account sequence and its
// ETH nonce, and an operator restarting it hits both. Exceeding either budget
// returns an error rather than waiting, so the caller can exit non-zero.
//
// A transaction that has been broadcast but not yet confirmed is NOT waited for
// and NOT cancelled. Waiting reopens an unbounded window; cancelling races a
// transaction that may already be in a block. The packet is in the durable
// pending tracker, and on restart the on-chain receipt query decides whether it
// was delivered -- which is the whole reason those two mechanisms exist. Do not
// "improve" this into a wait.
// workersStuck reports whether any engine gave up on its own inner drain with a
// worker still alive. wait below cannot see that: it waits on the engine
// goroutines, and an engine returns as soon as ITS drain times out, leaving the
// orphaned worker outside every WaitGroup. Without this the wait succeeds and
// the cleanups run underneath a live worker -- the crash the skip above exists
// to avoid, reached one level down. Found in review.
func stopRelaysAndCleanup(cancel func(), wait func(), cleanups []func(), workersStuck func() bool) error {
	cancel()

	if err := runBounded(wait, shutdownDrainBudget); err != nil {
		// The clients are deliberately left open: a goroutine that has not
		// returned may still be using them, and closing underneath it trades a
		// slow shutdown for a crash. Report and let the process exit.
		return fmt.Errorf("relay engines did not drain within %s", shutdownDrainBudget)
	}
	if workersStuck != nil && workersStuck() {
		// Same reason, different level: every engine goroutine has returned, but
		// one of them reported a worker it could not stop.
		return fmt.Errorf("a relay source could not stop its own workers within %s; "+
			"leaving clients open rather than closing them underneath one that is still running",
			workergroup.SourceDrainTimeout)
	}

	if err := runBounded(func() {
		for _, cleanup := range cleanups {
			cleanup()
		}
	}, shutdownCleanupBudget); err != nil {
		return fmt.Errorf("client cleanup did not finish within %s", shutdownCleanupBudget)
	}
	return nil
}

// awaitRelayExit blocks until the relay path ends and returns the process's exit
// error: nil for a clean stop, non-nil for anything an operator has to act on.
//
// Extracted from Start so all three endings can be tested. The one that used to
// be wrong is unreachable from a unit test while it lives inside the command
// closure, and it was wrong for exactly that reason.
//
// loopErrCh is BUFFERED, which is what makes the reads below load-bearing: an
// engine can record its failure and return without anyone ever selecting on the
// channel.
func awaitRelayExit(
	runCtx context.Context,
	done <-chan struct{},
	loopErrCh <-chan error,
	stopAndCleanup func() error,
	logger *zap.Logger,
) error {
	select {
	case err := <-loopErrCh:
		// An engine failed while running. That error is the actionable one, so a
		// shutdown that also overruns is logged rather than returned in its place.
		logShutdownFailure(logger, stopAndCleanup())
		return err

	case <-done:
		// A worker sends its error before its deferred wg.Done. If all workers
		// have returned, prefer that buffered error over treating the coincident
		// done signal as a clean shutdown.
		select {
		case err := <-loopErrCh:
			return err
		default:
			return nil
		}

	case <-runCtx.Done():
		// SIGTERM. A module that overran its drain reports which workers are still
		// running (relay.Module.Run) and that error lands in the buffered channel.
		// This branch did not read it, so a stuck worker produced exit code 0 with
		// "Relayer clients stopped; exiting" as the last line -- the failure was
		// recorded and then dropped, which is the worst of the three outcomes
		// because it looks exactly like success. The read is the same non-blocking
		// one the <-done branch above already does; the two must stay in step.
		logger.Sugar().Infof("Relayer shutdown requested: %v", runCtx.Err())
		shutdownErr := stopAndCleanup()
		select {
		case err := <-loopErrCh:
			if err != nil {
				return err
			}
		default:
		}
		if shutdownErr != nil {
			return shutdownErr
		}
		logger.Sugar().Info("Relayer clients stopped; exiting")
		return nil
	}
}

// logShutdownFailure reports a shutdown that overran its budget on a path whose
// return value is already carrying a more actionable error (a build failure, or
// the loop error that caused the shutdown). Swallowing it silently is what made
// a stuck worker invisible; returning it would hide the cause behind the symptom.
func logShutdownFailure(logger *zap.Logger, err error) {
	if err != nil {
		logger.Sugar().Errorf("Shutdown did not finish cleanly: %v", err)
	}
}

// runBounded runs fn and returns an error if it has not returned within budget.
// fn keeps running in its goroutine afterwards -- there is no way to interrupt an
// arbitrary func() -- so the caller must treat the timeout as "exit now", not as
// "retry".
func runBounded(fn func(), budget time.Duration) error {
	done := make(chan struct{})
	go func() {
		fn()
		close(done)
	}()
	timer := time.NewTimer(budget)
	defer timer.Stop()
	select {
	case <-done:
		return nil
	case <-timer.C:
		return context.DeadlineExceeded
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

// releaseOrHold closes one relay path's clients, unless that path just reported
// a worker it could not stop.
//
// The per-goroutine cleanup used to be an unconditional `defer cleanup()`, which
// defeated the whole point of the skip in stopRelaysAndCleanup: an engine
// returns the instant ITS drain times out, so the defer fired while the named
// worker was still using the clients it then closed. onceCleanup made that the
// winning call, so the central skip never got a say. Found in review.
//
// It also raises the flag stopRelaysAndCleanup reads, so the central cleanups --
// which hold the same closure -- stand down too.
func releaseOrHold(
	logger *zap.Logger,
	cleanup func(),
	workersStuck *atomic.Bool,
	err error,
	kind, id string,
) {
	if errors.Is(err, workergroup.ErrWorkersStillRunning) {
		workersStuck.Store(true)
		logger.Sugar().Errorf("%s %q: a worker did not stop, so its clients are left open "+
			"rather than closed underneath it: %v", kind, id, err)
		return
	}
	cleanup()
}
