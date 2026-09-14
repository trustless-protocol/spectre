package main

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	workergroup "relayer/workers"

	"go.uber.org/zap"
)

// P1c, found in review. stopRelaysAndCleanup already refuses to close clients
// when its own wait overruns, for a stated reason: "a goroutine that has not
// returned may still be using them, and closing underneath it trades a slow
// shutdown for a crash". Its wait cannot see a worker orphaned INSIDE a source,
// because the engine returns as soon as its own drain times out and the
// orphaned goroutine is in no WaitGroup. The wait then succeeds and the
// cleanups run over a live worker -- the same crash, one level down.
func TestCleanupIsHeldWhenASourceCouldNotStopItsWorkers(t *testing.T) {
	t.Run("held when a source reports stuck workers", func(t *testing.T) {
		cleanupRan := false
		err := stopRelaysAndCleanup(
			func() {}, func() {},
			[]func(){func() { cleanupRan = true }},
			func() bool { return true },
		)
		if cleanupRan {
			t.Fatal("closed the clients while a worker was still using them; " +
				"that is the crash the drain-overrun branch already refuses to cause")
		}
		if err == nil {
			t.Fatal("a shutdown that left clients open must be reported, or it looks clean")
		}
		if !strings.Contains(err.Error(), "leaving clients open") {
			t.Fatalf("error %q does not say what was skipped or why", err)
		}
	})

	t.Run("runs when every worker stopped", func(t *testing.T) {
		cleanupRan := false
		err := stopRelaysAndCleanup(
			func() {}, func() {},
			[]func(){func() { cleanupRan = true }},
			func() bool { return false },
		)
		if err != nil {
			t.Fatalf("clean shutdown reported an error: %v", err)
		}
		if !cleanupRan {
			t.Fatal("clients were never closed on a clean shutdown; the guard is too eager " +
				"and every ordinary stop now leaks its connections")
		}
	})
}

// The per-path half of the same contract: the goroutine's own cleanup must
// stand down too. It used to be an unconditional `defer cleanup()`, which fired
// the instant the engine returned -- and onceCleanup made it the winning call,
// so the central skip never got a say.
func TestReleaseOrHold(t *testing.T) {
	for _, tc := range []struct {
		name       string
		err        error
		wantClosed bool
		wantFlag   bool
	}{
		{name: "a clean stop closes the clients", err: nil, wantClosed: true},
		{name: "an ordinary failure still closes them", err: errors.New("rpc broke"), wantClosed: true},
		{
			name: "a stuck worker holds them open and raises the flag",
			err:  errWrap(workergroup.ErrWorkersStillRunning),
			// no close, and the central cleanups must stand down as well
			wantFlag: true,
		},
		{
			name:     "a stuck worker joined with a real failure still holds them",
			err:      errors.Join(errors.New("proof failed"), errWrap(workergroup.ErrWorkersStillRunning)),
			wantFlag: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			closed := false
			var flag atomic.Bool
			releaseOrHold(zap.NewNop(), func() { closed = true }, &flag, tc.err, "cosmos_to_eth", "client-0")

			if closed != tc.wantClosed {
				t.Fatalf("clients closed = %v, want %v", closed, tc.wantClosed)
			}
			if flag.Load() != tc.wantFlag {
				t.Fatalf("workersStuck = %v, want %v: the central cleanups read this, so a "+
					"missed flag closes the same clients a moment later", flag.Load(), tc.wantFlag)
			}
		})
	}
}

func errWrap(err error) error { return errors.Join(errors.New("cosmos source"), err) }

// P1b: awaitRelayExit existed, was documented as the fix, was tested -- and
// Start never called it, so the SIGTERM ending still dropped a recorded failure
// and returned nil. This pins the contract awaitRelayExit implements, so the
// wiring test below has something to be the wiring OF.
func TestAwaitRelayExitReportsAFailureRecordedDuringShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // SIGTERM already landed

	loopErrCh := make(chan error, 1)
	loopErrCh <- errors.New("cosmos source: workers still running")

	err := awaitRelayExit(ctx, make(chan struct{}), loopErrCh,
		func() error { return nil }, zap.NewNop())
	if err == nil {
		t.Fatal("exit code 0 with a failure sitting in the buffered channel: the relayer " +
			"reports success and the operator never learns a worker was stuck")
	}
}
