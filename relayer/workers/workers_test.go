package workers

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestDrain_WorkerGroupDrainsAndNamesStuckWorkers(t *testing.T) {
	g := New()
	done := make(chan struct{})
	g.Go("subscriber", func() { <-done })

	if running := g.stillRunning(); !slices.Equal(running, []string{"subscriber"}) {
		t.Fatalf("running = %v, want [subscriber]", running)
	}
	if running := g.Drain(20 * time.Millisecond); !slices.Equal(running, []string{"subscriber"}) {
		t.Fatalf("timed-out workers = %v, want [subscriber]", running)
	}
	close(done)
	if running := g.Drain(time.Second); len(running) != 0 {
		t.Fatalf("workers after release = %v, want none", running)
	}
}

// The nested drains must not share a budget. A source drains its own goroutines
// inside the drain the relay module runs over the source; with equal timeouts
// both expire on the same cancellation at the same moment, and whether the module
// observes the stuck worker or sees its subscriber return first is a race.
//
// This is an invariant rather than a behaviour test on purpose: reproducing the
// race takes the real timeouts, and the property that prevents it is one
// comparison. Reported by @DongLieu on #434.
func TestTheInnerDrainBudgetIsStrictlyShorterThanTheOuterOne(t *testing.T) {
	if SourceDrainTimeout >= ShutdownDrainTimeout {
		t.Fatalf("SourceDrainTimeout = %s, ShutdownDrainTimeout = %s; the inner drain must finish "+
			"while the outer one is still listening, or a stuck source worker is reported as a clean shutdown",
			SourceDrainTimeout, ShutdownDrainTimeout)
	}
	if SourceDrainTimeout <= 0 {
		t.Fatalf("SourceDrainTimeout = %s; a non-positive budget drains nothing", SourceDrainTimeout)
	}
}

// #575: no relay goroutine recovered a panic, so one bad RPC answer ended the
// process, both directions. A panic is now an error the group's owner receives.
func TestGo_RecoversAPanicAndReportsIt(t *testing.T) {
	g := New()
	g.Go("packet-flush", func() { panic("nil map") })

	select {
	case err := <-g.Failed():
		if !errors.Is(err, ErrWorkerPanicked) {
			t.Fatalf("err = %v, want ErrWorkerPanicked", err)
		}
		var panicErr *PanicError
		if !errors.As(err, &panicErr) || panicErr.Worker != "packet-flush" || panicErr.Value != "nil map" {
			t.Fatalf("err = %#v, want the worker and the panic value", err)
		}
	case <-time.After(time.Second):
		t.Fatal("the panic was not reported")
	}
	if running := g.Drain(time.Second); len(running) != 0 {
		t.Fatalf("the panicked worker is still listed as running: %v", running)
	}
}

// Failed holds one error. A second panic must neither block its worker nor be
// mistaken for the first.
func TestGo_ASecondPanicDoesNotBlock(t *testing.T) {
	g := New()
	g.Go("first", func() { panic("one") })
	g.Go("second", func() { panic("two") })

	if running := g.Drain(time.Second); len(running) != 0 {
		t.Fatalf("a panicking worker blocked delivering its error: %v still running", running)
	}
	select {
	case err := <-g.Failed():
		if !errors.Is(err, ErrWorkerPanicked) {
			t.Fatalf("err = %v, want ErrWorkerPanicked", err)
		}
	default:
		t.Fatal("no panic was reported")
	}
}

func TestGo_AWorkerThatReturnsReportsNothing(t *testing.T) {
	g := New()
	g.Go("quiet", func() {})
	g.Drain(time.Second)
	select {
	case err := <-g.Failed():
		t.Fatalf("a worker that returned normally reported %v", err)
	default:
	}
}

func TestRecover(t *testing.T) {
	want := errors.New("rpc down")
	if err := Recover("engine", func() error { return want }); err != want {
		t.Fatalf("err = %v, want the function's own error unchanged", err)
	}
	if err := Recover("engine", func() error { return nil }); err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	err := Recover("cosmos<->eth engine", func() error {
		var m map[string]int
		m["x"]++ // a nil-map write, the kind of panic an unchecked RPC answer causes
		return nil
	})
	if !errors.Is(err, ErrWorkerPanicked) || !strings.Contains(err.Error(), "cosmos<->eth engine panicked") {
		t.Fatalf("err = %v, want a PanicError naming the engine", err)
	}
}
