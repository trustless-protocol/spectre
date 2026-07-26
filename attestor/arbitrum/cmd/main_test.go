package main

import (
	"bytes"
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"attestor/arbitrum"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestStartCommandPassesConfigAndContext(t *testing.T) {
	type contextKey struct{}
	ctx := context.WithValue(context.Background(), contextKey{}, "command-context")

	var receivedConfig string
	var receivedContextValue any
	command := newRootCommand(func(ctx context.Context, configPath string) error {
		receivedConfig = configPath
		receivedContextValue = ctx.Value(contextKey{})
		return nil
	})
	command.SetArgs([]string{"start", "--config", "custom-arbitrum.json"})
	command.SetOut(&bytes.Buffer{})
	command.SetErr(&bytes.Buffer{})

	if err := command.ExecuteContext(ctx); err != nil {
		t.Fatalf("execute start command: %v", err)
	}
	if receivedConfig != "custom-arbitrum.json" {
		t.Fatalf("config path: got %q want %q", receivedConfig, "custom-arbitrum.json")
	}
	if receivedContextValue != "command-context" {
		t.Fatalf("command context value: got %v", receivedContextValue)
	}
}

func TestStartCommandUsesDefaultConfig(t *testing.T) {
	var receivedConfig string
	command := newRootCommand(func(_ context.Context, configPath string) error {
		receivedConfig = configPath
		return nil
	})
	command.SetArgs([]string{"start"})
	command.SetOut(&bytes.Buffer{})
	command.SetErr(&bytes.Buffer{})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("execute start command: %v", err)
	}
	if receivedConfig != defaultConfigPath {
		t.Fatalf("default config path: got %q want %q", receivedConfig, defaultConfigPath)
	}
}

func TestStartCommandRejectsArguments(t *testing.T) {
	command := newRootCommand(func(context.Context, string) error {
		t.Fatal("start runner called with positional arguments")
		return nil
	})
	command.SetArgs([]string{"start", "unexpected"})
	command.SetOut(&bytes.Buffer{})
	command.SetErr(&bytes.Buffer{})

	if err := command.ExecuteContext(context.Background()); err == nil {
		t.Fatal("expected positional argument error")
	}
}

func TestMonitorRuntimeStateRefreshesOnHeadsAndReconnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	refresher := &testRuntimeRefresher{calls: make(chan struct{}, 16)}
	attempts := make(chan testHeadSubscriptionAttempt, 2)
	subscriber := headSubscriberFunc(func(
		_ context.Context,
		headers chan<- *types.Header,
	) (ethereum.Subscription, error) {
		subscription := newTestSubscription()
		attempts <- testHeadSubscriptionAttempt{
			headers:      headers,
			subscription: subscription,
		}
		return subscription, nil
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		monitorRuntimeStateWithRetry(
			ctx,
			refresher,
			subscriber,
			time.Hour,
			time.Millisecond,
		)
	}()

	first := waitForSubscriptionAttempt(t, attempts)
	waitForRefreshes(t, refresher.calls, 2) // Startup and post-subscribe reconciliation.

	first.headers <- &types.Header{}
	waitForRefreshes(t, refresher.calls, 1)

	first.subscription.err <- errors.New("test disconnect")
	second := waitForSubscriptionAttempt(t, attempts)
	waitForRefreshes(t, refresher.calls, 1) // Post-reconnect reconciliation.

	cancel()
	select {
	case <-second.subscription.unsubscribed:
	case <-time.After(time.Second):
		t.Fatal("reconnected subscription was not unsubscribed")
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runtime monitor did not stop")
	}
}

func TestMonitorRuntimeStatePeriodicallyReconcilesWithoutHeads(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	refresher := &testRuntimeRefresher{calls: make(chan struct{}, 16)}
	attempts := make(chan testHeadSubscriptionAttempt, 1)
	subscriber := headSubscriberFunc(func(
		_ context.Context,
		headers chan<- *types.Header,
	) (ethereum.Subscription, error) {
		subscription := newTestSubscription()
		attempts <- testHeadSubscriptionAttempt{
			headers:      headers,
			subscription: subscription,
		}
		return subscription, nil
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		monitorRuntimeStateWithRetry(
			ctx,
			refresher,
			subscriber,
			10*time.Millisecond,
			time.Hour,
		)
	}()

	attempt := waitForSubscriptionAttempt(t, attempts)
	waitForRefreshes(t, refresher.calls, 3) // Startup, subscribe, then fallback tick.
	cancel()

	select {
	case <-attempt.subscription.unsubscribed:
	case <-time.After(time.Second):
		t.Fatal("subscription was not unsubscribed")
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runtime monitor did not stop")
	}
}

func TestServeWithManagedNitroStopsWhenNitroExits(t *testing.T) {
	listener := bufconn.Listen(1024)
	server := grpc.NewServer()
	done := make(chan struct{})
	close(done)
	nitro := testNitroMonitor{done: done, err: errors.New("exit status 17")}

	err := serveWithManagedNitro(context.Background(), server, listener, nitro)
	if err == nil || !strings.Contains(err.Error(), "exit status 17") {
		t.Fatalf("managed Nitro exit error: got %v", err)
	}
}

func TestServeWithManagedNitroStopsOnContextCancellation(t *testing.T) {
	listener := bufconn.Listen(1024)
	server := grpc.NewServer()
	nitroDone := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := serveWithManagedNitro(ctx, server, listener, testNitroMonitor{done: nitroDone})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("context cancellation: got %v", err)
	}
}

func TestGracefulStopReturns(t *testing.T) {
	listener := bufconn.Listen(1024)
	server := grpc.NewServer()
	serveDone := make(chan error, 1)
	go func() { serveDone <- server.Serve(listener) }()

	gracefulStop(server, time.Second)
	if err := <-serveDone; err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		t.Fatalf("serve after graceful stop: got %v", err)
	}
}

type testNitroMonitor struct {
	done <-chan struct{}
	err  error
}

func (m testNitroMonitor) Done() <-chan struct{} { return m.done }
func (m testNitroMonitor) Err() error            { return m.err }

type testRuntimeRefresher struct {
	calls chan struct{}
}

func (r *testRuntimeRefresher) Refresh(
	context.Context,
) ([]arbitrum.FinalizedConsistency, error) {
	r.calls <- struct{}{}
	return nil, nil
}

type headSubscriberFunc func(
	context.Context,
	chan<- *types.Header,
) (ethereum.Subscription, error)

func (f headSubscriberFunc) SubscribeNewHead(
	ctx context.Context,
	headers chan<- *types.Header,
) (ethereum.Subscription, error) {
	return f(ctx, headers)
}

type testHeadSubscriptionAttempt struct {
	headers      chan<- *types.Header
	subscription *testSubscription
}

type testSubscription struct {
	err          chan error
	unsubscribed chan struct{}
	once         sync.Once
}

func newTestSubscription() *testSubscription {
	return &testSubscription{
		err:          make(chan error, 1),
		unsubscribed: make(chan struct{}),
	}
}

func (s *testSubscription) Err() <-chan error {
	return s.err
}

func (s *testSubscription) Unsubscribe() {
	s.once.Do(func() {
		close(s.unsubscribed)
	})
}

func waitForSubscriptionAttempt(
	t *testing.T,
	attempts <-chan testHeadSubscriptionAttempt,
) testHeadSubscriptionAttempt {
	t.Helper()
	select {
	case attempt := <-attempts:
		return attempt
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Nitro subscription attempt")
		return testHeadSubscriptionAttempt{}
	}
}

func waitForRefreshes(t *testing.T, calls <-chan struct{}, count int) {
	t.Helper()
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for index := range count {
		select {
		case <-calls:
		case <-timer.C:
			t.Fatalf("timed out after %d of %d runtime refreshes", index, count)
		}
	}
}

var _ net.Listener = (*bufconn.Listener)(nil)
