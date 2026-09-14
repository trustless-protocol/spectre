package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"

	"relayer/workers"
)

func TestClientCleanupRunsAfterRelayDrain(t *testing.T) {
	var order []string
	stopRelaysAndCleanup(
		func() { order = append(order, "cancel") },
		func() { order = append(order, "drain") },
		[]func(){func() { order = append(order, "cleanup") }},
		nil,
	)
	want := []string{"cancel", "drain", "cleanup"}
	if fmt.Sprint(order) != fmt.Sprint(want) {
		t.Fatalf("order = %v, want %v", order, want)
	}
}

func TestIsShutdownErr(t *testing.T) {
	for _, tc := range []struct {
		err  error
		want bool
	}{
		{context.Canceled, true},
		{context.DeadlineExceeded, true},
		{fmt.Errorf("wrapped: %w", context.Canceled), true},
		{errors.New("request was cancelled by server"), false},
		{nil, false},
	} {
		if got := isShutdownErr(tc.err); got != tc.want {
			t.Fatalf("isShutdownErr(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}

// withShortShutdownBudgets runs the same code path in milliseconds. The budgets
// are the only thing shortened -- every branch under test is the production one.
func withShortShutdownBudgets(t *testing.T, drain, cleanup time.Duration) {
	t.Helper()
	origDrain, origCleanup := shutdownDrainBudget, shutdownCleanupBudget
	shutdownDrainBudget, shutdownCleanupBudget = drain, cleanup
	t.Cleanup(func() { shutdownDrainBudget, shutdownCleanupBudget = origDrain, origCleanup })
}

// The budget an operator is promised. Nothing reads this at runtime, so a
// change to either half would otherwise be invisible until a SIGTERM in
// production took longer than the runbook says it can.
func TestShutdownBudgetIsFortyFiveSeconds(t *testing.T) {
	if shutdownBudget != 45*time.Second {
		t.Fatalf("shutdown budget = %s, want 45s (docs and the runbook state this number)", shutdownBudget)
	}
	if shutdownDrainBudget != workers.ShutdownDrainTimeout {
		t.Fatalf("drain budget = %s, want the per-module timeout %s: modules drain in parallel, so the "+
			"two are the same bound and must not drift", shutdownDrainBudget, workers.ShutdownDrainTimeout)
	}
}

// A worker that will not return must not hold the process. Before this, wait()
// was called bare: the relayer sat holding its Cosmos account sequence and its
// ETH nonce, and an operator restarting it collided with both.
func TestStopRelaysAndCleanupGivesUpOnAStuckDrain(t *testing.T) {
	withShortShutdownBudgets(t, 40*time.Millisecond, time.Second)

	release := make(chan struct{})
	t.Cleanup(func() { close(release) })
	cleanupRan := false

	start := time.Now()
	err := stopRelaysAndCleanup(
		func() {},
		func() { <-release }, // the worker that refuses to exit
		[]func(){func() { cleanupRan = true }},
		nil,
	)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("a drain that never returned was reported as a clean shutdown")
	}
	if !strings.Contains(err.Error(), "drain") {
		t.Fatalf("error %q does not say which phase overran", err)
	}
	if elapsed > 2*shutdownDrainBudget {
		t.Fatalf("waited %s for a %s budget", elapsed, shutdownDrainBudget)
	}
	// Deliberate: a goroutine that has not returned may still be using these
	// clients, and closing underneath it trades a slow shutdown for a crash.
	if cleanupRan {
		t.Fatal("clients were closed while a worker was still running")
	}
}

// The cleanup phase is bounded too -- cosmosClient.Stop() and ethClient.Close()
// talk to sockets, so neither is guaranteed to return.
func TestStopRelaysAndCleanupGivesUpOnAStuckCleanup(t *testing.T) {
	withShortShutdownBudgets(t, time.Second, 40*time.Millisecond)

	release := make(chan struct{})
	t.Cleanup(func() { close(release) })

	err := stopRelaysAndCleanup(func() {}, func() {}, []func(){func() { <-release }}, nil)
	if err == nil {
		t.Fatal("a cleanup that never returned was reported as a clean shutdown")
	}
	if !strings.Contains(err.Error(), "cleanup") {
		t.Fatalf("error %q does not say which phase overran", err)
	}
}

// The control case: nothing is stuck, so nothing is reported. Without this the
// tests above are satisfied by a function that always errors.
func TestStopRelaysAndCleanupReportsNothingOnACleanStop(t *testing.T) {
	withShortShutdownBudgets(t, time.Second, time.Second)

	if err := stopRelaysAndCleanup(func() {}, func() {}, []func(){func() {}}, nil); err != nil {
		t.Fatalf("a clean shutdown reported %v", err)
	}
}

// awaitRelayExit decides the process exit code. The SIGTERM branch is the one
// that was wrong: a module reports the workers it could not drain, that error
// sits in the BUFFERED channel, and the branch returned nil anyway -- a stuck
// worker looked exactly like a clean stop.
func TestAwaitRelayExit(t *testing.T) {
	logger := zap.NewNop()
	noopStop := func() error { return nil }

	t.Run("clean stop returns nil", func(t *testing.T) {
		done := make(chan struct{})
		close(done)
		if err := awaitRelayExit(context.Background(), done, make(chan error, 1), noopStop, logger); err != nil {
			t.Fatalf("clean stop returned %v", err)
		}
	})

	t.Run("an engine failure is returned", func(t *testing.T) {
		loopErrCh := make(chan error, 1)
		loopErrCh <- errors.New("cosmos_to_eth source \"cosmos-0\": boom")
		err := awaitRelayExit(context.Background(), make(chan struct{}), loopErrCh, noopStop, logger)
		if err == nil || !strings.Contains(err.Error(), "boom") {
			t.Fatalf("engine failure = %v, want the engine's own error", err)
		}
	})

	t.Run("SIGTERM surfaces a drain failure left in the buffer", func(t *testing.T) {
		// Exactly the shape relay.Module.Run produces on an overrun.
		const stuck = `relay cosmos->eth: shutdown drain timed out after 30s; workers still running: [cosmos-subscribe]`
		loopErrCh := make(chan error, 1)
		loopErrCh <- errors.New(stuck)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := awaitRelayExit(ctx, make(chan struct{}), loopErrCh, noopStop, logger)
		if err == nil {
			t.Fatal("SIGTERM with a stuck worker exited cleanly; that is the bug this branch had")
		}
		if !strings.Contains(err.Error(), "cosmos-subscribe") {
			t.Fatalf("error %q does not name the worker that is stuck", err)
		}
	})

	t.Run("SIGTERM surfaces a shutdown that overran its budget", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		stopErr := errors.New("relay engines did not drain within 30s")

		err := awaitRelayExit(ctx, make(chan struct{}), make(chan error, 1),
			func() error { return stopErr }, logger)
		if !errors.Is(err, stopErr) {
			t.Fatalf("SIGTERM returned %v, want the shutdown failure", err)
		}
	})

	// The control for the two above: an ordinary SIGTERM with nothing wrong exits
	// 0. Without it, "return the error" is satisfied by always returning one.
	t.Run("an ordinary SIGTERM exits cleanly", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := awaitRelayExit(ctx, make(chan struct{}), make(chan error, 1), noopStop, logger); err != nil {
			t.Fatalf("ordinary SIGTERM returned %v, want nil", err)
		}
	})
}
