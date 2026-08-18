package relay

import (
	"sort"
	"sync"
	"time"
)

// shutdownDrainTimeout bounds how long a module waits for source, refresh and
// scan workers after cancellation before reporting the workers that are stuck.
const shutdownDrainTimeout = 30 * time.Second

type workerGroup struct {
	wg sync.WaitGroup

	mu      sync.Mutex
	running map[string]struct{}
}

func newWorkerGroup() *workerGroup {
	return &workerGroup{running: make(map[string]struct{})}
}

func (g *workerGroup) Go(name string, fn func()) {
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
		fn()
	}()
}

func (g *workerGroup) stillRunning() []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	names := make([]string, 0, len(g.running))
	for name := range g.running {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// drain returns the names of workers still running when timeout expires. An
// empty result means every worker returned.
func (g *workerGroup) drain(timeout time.Duration) []string {
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
