package main

import (
	"bytes"
	"context"
	"errors"
	"math/big"
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

func TestAssertionSourceIdentityPinsBoLDStorageLayout(t *testing.T) {
	config := arbitrum.DaemonConfig{
		L1ChainID:             11_155_111,
		L2ChainID:             421_614,
		RollupCoreAddress:     "0x042B2E6C5E99d4c521bd49beeD5E99651D9B0Cf4",
		AssertionsMappingSlot: "0x75",
		AssertionStatusOffset: 25,
	}
	identity := assertionSourceIdentity(config)
	for _, expected := range []string{
		"l1:11155111",
		"l2:421614",
		"rollup:0x042B2E6C5E99d4c521bd49beeD5E99651D9B0Cf4",
		"assertions:0x0000000000000000000000000000000000000000000000000000000000000075",
		"status-offset:25",
	} {
		if !strings.Contains(identity, expected) {
			t.Fatalf("source identity %q does not contain %q", identity, expected)
		}
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

func TestServeAttestorStopsOnContextCancellation(t *testing.T) {
	listener := bufconn.Listen(1024)
	server := grpc.NewServer()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := serveAttestor(ctx, server, listener)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("context cancellation: got %v", err)
	}
}

func TestValidateNitroChainID(t *testing.T) {
	if err := validateNitroChainID(
		context.Background(),
		chainIDReaderFunc(func(context.Context) (*big.Int, error) {
			return big.NewInt(421614), nil
		}),
		421614,
		"RPC",
	); err != nil {
		t.Fatalf("validate matching chain ID: %v", err)
	}

	err := validateNitroChainID(
		context.Background(),
		chainIDReaderFunc(func(context.Context) (*big.Int, error) {
			return big.NewInt(42161), nil
		}),
		421614,
		"WebSocket",
	)
	if err == nil || !strings.Contains(err.Error(), "expected 421614") {
		t.Fatalf("validate mismatched chain ID: got %v", err)
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

type chainIDReaderFunc func(context.Context) (*big.Int, error)

func (f chainIDReaderFunc) ChainID(ctx context.Context) (*big.Int, error) {
	return f(ctx)
}

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
