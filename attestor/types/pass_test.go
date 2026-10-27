package types

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRunPassWithoutTimeoutLeavesTheContextAlone(t *testing.T) {
	for _, timeout := range []time.Duration{0, -time.Second} {
		parent := context.Background()
		err := RunPass(parent, timeout, func(ctx context.Context) error {
			if ctx != parent {
				t.Fatalf("timeout %s: pass got a derived context, want the caller's", timeout)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("timeout %s: %v", timeout, err)
		}
	}
}

func TestRunPassCutsAPassThatNeverReturns(t *testing.T) {
	started := time.Now()
	err := RunPass(context.Background(), 50*time.Millisecond, func(ctx context.Context) error {
		<-ctx.Done() // a call parked on a connection that never answers
		return ctx.Err()
	})
	if err == nil {
		t.Fatal("RunPass returned nil for a pass cut by its deadline")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error %q does not wrap context.DeadlineExceeded", err)
	}
	if !strings.Contains(err.Error(), "pass exceeded its 50ms deadline") {
		t.Fatalf("error %q does not name the pass deadline", err)
	}
	if elapsed := time.Since(started); elapsed > 5*time.Second {
		t.Fatalf("pass was cut after %s, want about 50ms", elapsed)
	}
}

func TestRunPassReturnsParentCancellationUnlabelled(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	cancel()
	err := RunPass(parent, time.Hour, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if strings.Contains(err.Error(), "deadline") {
		t.Fatalf("shutdown reported as a pass deadline: %q", err)
	}
}

func TestRunPassDoesNotBlameTheDeadlineForATransportTimeout(t *testing.T) {
	// An HTTP client timeout also satisfies errors.Is(err, DeadlineExceeded),
	// while the pass context is still live.
	transport := errors.New("Post: Client.Timeout exceeded")
	err := RunPass(context.Background(), time.Hour, func(context.Context) error {
		return errors.Join(transport, context.DeadlineExceeded)
	})
	if err == nil || strings.Contains(err.Error(), "pass exceeded") {
		t.Fatalf("transport timeout reported as %v, want it returned as is", err)
	}
}

func TestRunPassPassesThroughOrdinaryResults(t *testing.T) {
	if err := RunPass(context.Background(), time.Hour, func(context.Context) error { return nil }); err != nil {
		t.Fatalf("success became %v", err)
	}
	failure := errors.New("upstream said no")
	if err := RunPass(context.Background(), time.Hour, func(context.Context) error { return failure }); err != failure {
		t.Fatalf("failure became %v, want it returned as is", err)
	}
}

func TestFeedStaleAfterCoversOnePassAndTwoTicks(t *testing.T) {
	if got := FeedStaleAfter(5*time.Minute, 30*time.Second); got != 6*time.Minute {
		t.Fatalf("FeedStaleAfter = %s, want 6m0s", got)
	}
}

func TestHeartbeatAliveWindow(t *testing.T) {
	const maxAge = time.Minute
	start := time.Unix(1_000_000, 0)
	var beat Heartbeat

	if beat.Alive(start, maxAge) {
		t.Fatal("a heartbeat that never beat is alive")
	}
	beat.Beat(start)
	for _, tc := range []struct {
		name string
		now  time.Time
		want bool
	}{
		{"at the beat", start, true},
		{"just inside the window", start.Add(maxAge), true},
		{"just outside the window", start.Add(maxAge + time.Nanosecond), false},
	} {
		if got := beat.Alive(tc.now, maxAge); got != tc.want {
			t.Errorf("%s: Alive = %t, want %t", tc.name, got, tc.want)
		}
	}
	beat.Beat(start.Add(10 * time.Minute))
	if !beat.Alive(start.Add(10*time.Minute+time.Second), maxAge) {
		t.Fatal("a fresh beat did not revive the heartbeat")
	}
}

func TestDialRPCBoundsHangingHTTPCall(t *testing.T) {
	server := blackHoleServer(t)
	client, err := DialRPC(context.Background(), server.URL, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("DialRPC: %v", err)
	}
	t.Cleanup(client.Close)

	done := make(chan error, 1)
	go func() {
		var out string
		done <- client.CallContext(context.Background(), &out, "eth_blockNumber")
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("hanging HTTP call succeeded")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("hanging HTTP call was not bounded")
	}
}

func TestDialRPCRejectsAnUnparsableURL(t *testing.T) {
	if _, err := DialRPC(context.Background(), "http://[::1", time.Second); err == nil {
		t.Fatal("DialRPC accepted a malformed URL")
	}
}
