package types

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

// DefaultPassTimeout bounds one pass of a feed runner: the OP attestation pass
// and the Arbitrum assertion sync. It is the catch-all for upstream transports
// that carry no timeout of their own (WebSocket and IPC connections, where
// DialRPC applies none), and it sits above DefaultRPCTimeout so the HTTP
// transport timeout still fires first on an HTTP endpoint.
//
// A pass that hits the deadline fails like any other failed pass. Both runners
// commit progress item by item (OP per game and per verdict, Arbitrum per log
// chunk), so the next pass resumes where this one stopped and a backlog larger
// than one deadline still drains.
const DefaultPassTimeout = 5 * time.Minute

// RunPass runs one feed-runner pass under the per-pass deadline. A timeout of
// zero or less runs the pass with no deadline of its own.
//
// An error caused by the deadline is labelled as such, so the operator sees
// which bound fired rather than a bare "context deadline exceeded". The
// label is added only when the pass context itself expired: an HTTP transport
// timeout also satisfies errors.Is(err, context.DeadlineExceeded), and must
// not be reported as this deadline. A cancelled parent context is returned
// as is, for the caller to treat as shutdown.
func RunPass(ctx context.Context, timeout time.Duration, pass func(context.Context) error) error {
	if timeout <= 0 {
		return pass(ctx)
	}
	passCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	err := pass(passCtx)
	if err != nil && ctx.Err() == nil &&
		errors.Is(passCtx.Err(), context.DeadlineExceeded) && errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("pass exceeded its %s deadline: %w", timeout, err)
	}
	return err
}

// FeedStaleAfter is how long a feed runner's loop may go without coming back
// around before it is no longer considered alive: one pass bounded by the
// pass timeout, then one poll tick, plus one more tick of slack.
func FeedStaleAfter(passTimeout, pollInterval time.Duration) time.Duration {
	return passTimeout + 2*pollInterval
}

// Heartbeat records when a feed runner's loop last came back around. It tells
// a runner that is failing (it still beats, and the failure shows in its log,
// its failure counter and its backoff) from one that is parked on a call that
// never returns (it stops beating, and nothing else reports it).
//
// The zero value has never beaten and is not alive. It is safe for concurrent
// use: the runner goroutine beats, readiness probes read.
type Heartbeat struct {
	last atomic.Int64 // UnixNano of the last Beat; 0 until the first one
}

// Beat records that the runner's loop is advancing at now.
func (h *Heartbeat) Beat(now time.Time) {
	h.last.Store(now.UnixNano())
}

// Alive reports whether the last beat is no more than maxAge before now.
func (h *Heartbeat) Alive(now time.Time, maxAge time.Duration) bool {
	last := h.last.Load()
	return last != 0 && now.Sub(time.Unix(0, last)) <= maxAge
}
