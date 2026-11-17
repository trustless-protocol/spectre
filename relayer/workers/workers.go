// Package workers gives a set of goroutines a name, so a shutdown that does not
// finish can say WHICH goroutine is still running rather than only that one is.
//
// It lives in its own package rather than being exported from relay/ because
// both relay/ and the chain adapters need it, and the dependency edge between
// those two runs relay -> chain: relay drives adapters through the chain
// interfaces, and adding a chain is meant to touch chain/ alone. An adapter
// importing relay/ would reverse that and make the engine a dependency of every
// adapter. A leaf package neither of them owns adds no edge in either direction.
package workers

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sort"
	"sync"
	"time"
)

// ShutdownDrainTimeout bounds how long one group waits for its workers after
// cancellation before reporting the ones that are stuck. Modules run in
// parallel, so this does not accumulate across them; the process-level budget
// that contains it lives in cmd (ShutdownBudget).
const ShutdownDrainTimeout = 30 * time.Second

// SourceDrainTimeout bounds the drain a source adapter runs over its OWN inner
// goroutines, nested inside the module drain above.
//
// It is strictly shorter, and that is the whole point. With equal budgets the two
// drains start together on the same cancellation and expire together, so which
// one reports is a race: if the inner one finishes first the module sees its
// worker return, drains cleanly, and reports a successful shutdown even though a
// named goroutine inside the source never stopped. The outer drain has to outlast
// the inner one to be able to receive what it found.
const SourceDrainTimeout = ShutdownDrainTimeout * 2 / 3

// ErrWorkersStillRunning marks a drain that gave up with a worker alive.
//
// It is a sentinel rather than a message so a caller can act on it, not only
// log it: closing a client underneath a goroutine that has not returned trades
// a slow shutdown for a crash, which is the reason stopRelaysAndCleanup already
// leaves clients open when its own wait overruns. A source whose inner drain
// overran has exactly the same problem one level down, and errors.Is is how the
// engine above it finds out.
var ErrWorkersStillRunning = errors.New("workers still running after the drain timeout")

// ErrWorkerPanicked marks an error that is a recovered panic.
//
// No relay goroutine used to recover, so one nil dereference -- typically on an
// RPC answer in a shape nobody expected -- ended the process with a bare stack
// trace on stderr, both directions at once (#575). A panic is now an error that
// names its worker and travels the same path as any other failure: the owner
// stops the rest in order and returns it.
//
// It is a reason to stop, not to carry on. A worker that panicked has stopped
// mid-step, and may have left a mutex locked or state half-written.
var ErrWorkerPanicked = errors.New("worker panicked")

// PanicError is a recovered panic. The stack is logged once where the panic was
// recovered and deliberately not carried in the error.
type PanicError struct {
	Worker string
	Value  any
}

func (e *PanicError) Error() string { return fmt.Sprintf("%s panicked: %v", e.Worker, e.Value) }

// Is makes errors.Is(err, ErrWorkerPanicked) true for a recovered panic.
func (e *PanicError) Is(target error) bool { return target == ErrWorkerPanicked }

func recovered(name string, value any) error {
	// The level token leads the line because logbridge reads it from the start.
	log.Printf("ERROR [workers][ATTENTION] %s panicked: %v\n%s", name, value, debug.Stack())
	return &PanicError{Worker: name, Value: value}
}

// Recover runs fn and returns its error, or a PanicError if fn panicked. It is
// for a goroutine that is not part of a Group: an engine, or one module's Run.
func Recover(name string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recovered(name, r)
		}
	}()
	return fn()
}

// Group tracks a set of named goroutines and can report, after a bounded wait,
// which of them have not returned.
type Group struct {
	wg sync.WaitGroup

	mu      sync.Mutex
	running map[string]struct{}

	// failed carries the first panic of any worker. It is buffered so a worker
	// never blocks delivering it; a later panic is logged and not queued.
	failed chan error
}

// New returns an empty group.
func New() *Group {
	return &Group{running: make(map[string]struct{}), failed: make(chan error, 1)}
}

// Go runs fn in a new goroutine recorded under name.
//
// A panic in fn is recovered, logged with its stack, and delivered on Failed.
func (g *Group) Go(name string, fn func()) {
	g.wg.Add(1)
	g.mu.Lock()
	g.running[name] = struct{}{}
	g.mu.Unlock()

	go func() {
		defer func() {
			g.mu.Lock()
			delete(g.running, name)
			g.mu.Unlock()
			g.wg.Done()
		}()
		// Declared last so it runs first: the panic is reported while the worker
		// is still listed as running.
		defer func() {
			if r := recover(); r != nil {
				err := recovered(name, r)
				select {
				case g.failed <- err:
				default:
				}
			}
		}()
		fn()
	}()
}

// Failed delivers the first panic of any worker in the group as an error.
//
// The owner of the group must select on it. A recovered worker has stopped, and
// a group that carries on without it does part of its job in silence -- which is
// worse than the crash the recovery replaced.
func (g *Group) Failed() <-chan error { return g.failed }

func (g *Group) stillRunning() []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	names := make([]string, 0, len(g.running))
	for name := range g.running {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Drain returns the names of workers still running when timeout expires. An
// empty result means every worker returned.
func (g *Group) Drain(timeout time.Duration) []string {
	done := make(chan struct{})
	go func() {
		g.wg.Wait()
		close(done)
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-done:
		return nil
	case <-timer.C:
		return g.stillRunning()
	}
}
