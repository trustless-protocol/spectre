package relay

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"relayer/chain"
)

func TestWorkerGroupDrainsAndNamesStuckWorkers(t *testing.T) {
	g := newWorkerGroup()
	done := make(chan struct{})
	g.Go("subscriber", func() { <-done })

	if running := g.stillRunning(); !slices.Equal(running, []string{"subscriber"}) {
		t.Fatalf("running = %v, want [subscriber]", running)
	}
	if running := g.drain(20 * time.Millisecond); !slices.Equal(running, []string{"subscriber"}) {
		t.Fatalf("timed-out workers = %v, want [subscriber]", running)
	}
	close(done)
	if running := g.drain(time.Second); len(running) != 0 {
		t.Fatalf("workers after release = %v, want none", running)
	}
}

type drainingSource struct {
	started chan struct{}
	stopped chan struct{}
}

func (s *drainingSource) Chain() chain.ChainType { return chain.Cosmos }
func (s *drainingSource) Subscribe(ctx context.Context, _ func(context.Context, []chain.Event) []int) error {
	close(s.started)
	<-ctx.Done()
	close(s.stopped)
	return ctx.Err()
}
func (s *drainingSource) LatestHeight(context.Context) (uint64, error)        { return 0, nil }
func (s *drainingSource) RelayableHeight(context.Context) (uint64, error)     { return 0, nil }
func (s *drainingSource) QueryHeader(context.Context, uint64) ([]byte, error) { return nil, nil }
func (s *drainingSource) MembershipProof(context.Context, []byte, uint64, chain.EventType) ([]byte, error) {
	return nil, nil
}
func (s *drainingSource) NonMembershipProof(context.Context, []byte, uint64) ([]byte, error) {
	return nil, nil
}

func TestModuleCleanCancellationDrainsSubscriber(t *testing.T) {
	src := &drainingSource{started: make(chan struct{}), stopped: make(chan struct{})}
	m := NewModule("test", "client", src, &mockDest{}, &mockBuilder{})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- m.Run(ctx) }()
	<-src.started
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run cancellation error = %v, want nil", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not drain after cancellation")
	}
	select {
	case <-src.stopped:
	default:
		t.Fatal("Run returned before subscriber stopped")
	}
}

type internallyCancelledSource struct{}

func (internallyCancelledSource) Chain() chain.ChainType { return chain.Cosmos }
func (internallyCancelledSource) Subscribe(context.Context, func(context.Context, []chain.Event) []int) error {
	return context.Canceled
}
func (internallyCancelledSource) LatestHeight(context.Context) (uint64, error) { return 0, nil }
func (internallyCancelledSource) RelayableHeight(context.Context) (uint64, error) {
	return 0, nil
}
func (internallyCancelledSource) QueryHeader(context.Context, uint64) ([]byte, error) {
	return nil, nil
}
func (internallyCancelledSource) MembershipProof(context.Context, []byte, uint64, chain.EventType) ([]byte, error) {
	return nil, nil
}
func (internallyCancelledSource) NonMembershipProof(context.Context, []byte, uint64) ([]byte, error) {
	return nil, nil
}

func TestModuleReportsInternalSubscriberCancellation(t *testing.T) {
	m := NewModule("test", "client", internallyCancelledSource{}, &mockDest{}, &mockBuilder{})
	err := m.Run(context.Background())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run error = %v, want wrapped context.Canceled", err)
	}
}
