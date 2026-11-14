package relay

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"relayer/chain"
	workergroup "relayer/workers"
)

// panickingSource panics in Subscribe, the way an unchecked RPC answer does
// inside a source's event loop.
type panickingSource struct{ drainingSource }

func (s *panickingSource) Subscribe(context.Context, func(context.Context, []chain.Event) []int) error {
	var m map[string]int
	m["x"]++
	return nil
}

// #575: no relay goroutine recovered a panic, so one nil dereference in any loop
// ended the process with a bare stack trace, both directions. A panic is now the
// path's error: the module stops its other loops in order and returns it.
func TestRun_APanickingLoopStopsThePathAndIsReported(t *testing.T) {
	src := &drainingSource{started: make(chan struct{}), stopped: make(chan struct{})}
	m := NewModule("cosmos->eth", "client-0", src, &mockDest{}, &mockBuilder{},
		WithTimeoutScanner(5*time.Millisecond, func(context.Context) { panic("scanner hit a nil field") }))

	done := make(chan error, 1)
	go func() { done <- m.Run(context.Background()) }()

	select {
	case err := <-done:
		if !errors.Is(err, workergroup.ErrWorkerPanicked) {
			t.Fatalf("Run err = %v, want the panic reported as ErrWorkerPanicked", err)
		}
		for _, want := range []string{"cosmos->eth", workerTimeoutScan, "scanner hit a nil field"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("Run err = %v, want it to name %q", err, want)
			}
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run kept going after a loop panicked: the path would look alive with a loop missing")
	}
	select {
	case <-src.stopped:
	default:
		t.Fatal("Run returned before stopping the subscriber")
	}
}

func TestRun_APanickingSubscriberIsReported(t *testing.T) {
	m := NewModule("l2->cosmos", "client-0", &panickingSource{}, &mockDest{}, &mockBuilder{})

	done := make(chan error, 1)
	go func() { done <- m.Run(context.Background()) }()

	select {
	case err := <-done:
		if !errors.Is(err, workergroup.ErrWorkerPanicked) || !strings.Contains(err.Error(), workerSubscriber) {
			t.Fatalf("Run err = %v, want the subscriber's panic", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after its subscriber panicked")
	}
}
