package status

import (
	"context"
	"strings"
	"testing"
	"time"

	"relayer/relay"
	"relayer/services"
)

// Status is off by default, so an engine registers with a registry that may not
// exist. Making that the engine's problem would put a nil check at every call
// site and eventually miss one.
func TestNilRegistryAnswersInsteadOfPanicking(t *testing.T) {
	var reg *Registry
	reg.AddPath(func(context.Context) relay.Status { return relay.Status{} })
	reg.AddQueues("cosmos<->eth", func(time.Time) []services.QueueStatus { return nil })

	snap := reg.Snapshot(context.Background())
	if len(snap.Paths) != 0 || len(snap.Queues) != 0 {
		t.Fatalf("a nil registry produced %+v, want an empty snapshot", snap)
	}
}

func TestRegistryMergesEverySourceInOrder(t *testing.T) {
	reg := &Registry{}
	reg.AddPath(func(context.Context) relay.Status { return relay.Status{Path: "cosmos->eth"} })
	reg.AddPath(func(context.Context) relay.Status { return relay.Status{Path: "eth->cosmos"} })
	reg.AddQueues("cosmos<->eth", func(time.Time) []services.QueueStatus {
		return []services.QueueStatus{{Direction: "cosmos"}, {Direction: "eth"}}
	})

	snap := reg.Snapshot(context.Background())
	if len(snap.Paths) != 2 || snap.Paths[0].Path != "cosmos->eth" || snap.Paths[1].Path != "eth->cosmos" {
		t.Fatalf("paths = %+v, want both in registration order", snap.Paths)
	}
	if len(snap.Queues) != 2 {
		t.Fatalf("queues = %+v, want both directions flattened into one list", snap.Queues)
	}
	if snap.Now.IsZero() {
		t.Fatal("the snapshot carries no timestamp, so every age in it is unanchored")
	}
}

// The registry must not hold its lock across the functions it calls. One of them
// makes a bounded RPC, and holding the lock across that would let a slow endpoint
// block a second concurrent request -- during exactly the incident the endpoint
// exists for.
func TestRegistryDoesNotHoldItsLockWhileAsking(t *testing.T) {
	reg := &Registry{}
	reentered := false
	reg.AddPath(func(ctx context.Context) relay.Status {
		// Reenter exactly ONCE, guarded by a plain flag. Unconditional reentry is
		// infinite recursion rather than a lock test, and sync.Once cannot express
		// this either: calling Do from inside f deadlocks by design, which would
		// fail the test for a reason that has nothing to do with the registry.
		// One goroutine drives this, so the flag needs no synchronisation.
		if !reentered {
			reentered = true
			reg.Snapshot(ctx)
		}
		return relay.Status{Path: "reentrant"}
	})

	done := make(chan struct{})
	go func() { defer close(done); reg.Snapshot(context.Background()) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Snapshot deadlocked: the registry is holding its lock while calling out")
	}
}

// The regression @neitdung asked for: two engines, one holding the same
// (source_client, sequence) as the other, and the requested one alone changes.
//
// Client ids are chain-local, so a pair that identifies a packet inside one
// tracker identifies nothing across a process. The old Reactivate offered the
// pair to each registered engine in order and returned on the first success, so
// an operator acting on the SECOND row `relayer status` printed gave a fresh
// retry budget to the first and left the packet they meant dead-lettered.
func TestReactivateActsOnTheNamedEngineAlone(t *testing.T) {
	type call struct {
		direction    string
		sourceClient string
		sequence     uint64
	}
	var first, second []call

	reg := &Registry{}
	reg.AddReactivator("cosmos<->eth", func(direction, sourceClient string, sequence uint64) error {
		first = append(first, call{direction, sourceClient, sequence})
		return nil
	})
	reg.AddReactivator("cosmos->l2", func(direction, sourceClient string, sequence uint64) error {
		second = append(second, call{direction, sourceClient, sequence})
		return nil
	})

	// The SECOND engine, named explicitly. Under the old first-success scan this
	// reached the first one instead, and reported success for it.
	sel := DeadLetterSelector{Path: "cosmos->l2", Direction: "l2", SourceClient: "client-0", Sequence: 9}
	if err := reg.Reactivate(sel); err != nil {
		t.Fatalf("Reactivate(%s): %v", sel, err)
	}
	if len(first) != 0 {
		t.Fatalf("the engine that was NOT named was asked: %+v; registration order must not decide a write", first)
	}
	if len(second) != 1 || second[0] != (call{"l2", "client-0", 9}) {
		t.Fatalf("named engine received %+v, want one call {l2 client-0 9}", second)
	}
}

// A name nobody registered must be refused, and the refusal must list the names
// that exist -- the likely mistake is a typo in a string the operator just read.
func TestReactivateRefusesAPathThatIsNotRegistered(t *testing.T) {
	reg := &Registry{}
	reg.AddReactivator("cosmos<->eth", func(string, string, uint64) error {
		t.Fatal("an engine was asked for a path that does not name it")
		return nil
	})

	err := reg.Reactivate(DeadLetterSelector{Path: "typo", Direction: "cosmos", SourceClient: "c", Sequence: 1})
	if err == nil {
		t.Fatal("reactivating an unregistered path reported success")
	}
	if !strings.Contains(err.Error(), "cosmos<->eth") {
		t.Errorf("error %q does not name the paths this process serves", err)
	}
}

// Validation happens before anything is resolved, so a request that names no
// packet never reaches an engine.
func TestReactivateValidatesBeforeItResolves(t *testing.T) {
	reg := &Registry{}
	reg.AddReactivator("cosmos<->eth", func(string, string, uint64) error {
		t.Fatal("an incomplete selector reached an engine")
		return nil
	})

	for _, sel := range []DeadLetterSelector{
		{},
		{Direction: "cosmos", SourceClient: "c", Sequence: 1},
		{Path: "cosmos<->eth", SourceClient: "c", Sequence: 1},
		{Path: "cosmos<->eth", Direction: "cosmos", Sequence: 1},
		{Path: "cosmos<->eth", Direction: "cosmos", SourceClient: "c"},
	} {
		if err := reg.Reactivate(sel); err == nil {
			t.Errorf("Reactivate(%+v) reported success; every field is required", sel)
		}
	}
}

// The engine name a queue is reported under is the one a reactivation must
// quote back. If the snapshot printed one string and the write route accepted
// another, the operator would have nothing to copy.
func TestSnapshotStampsTheEngineThatReportedTheQueue(t *testing.T) {
	reg := &Registry{}
	reg.AddQueues("cosmos<->eth", func(time.Time) []services.QueueStatus {
		return []services.QueueStatus{{Direction: "cosmos"}, {Direction: "eth"}}
	})
	reg.AddQueues("cosmos->l2", func(time.Time) []services.QueueStatus {
		return []services.QueueStatus{{Direction: "l2"}}
	})

	snap := reg.Snapshot(context.Background())
	got := map[string]string{}
	for _, q := range snap.Queues {
		got[q.Direction] = q.Path
	}
	want := map[string]string{"cosmos": "cosmos<->eth", "eth": "cosmos<->eth", "l2": "cosmos->l2"}
	for direction, path := range want {
		if got[direction] != path {
			t.Errorf("queue %q reported path %q, want %q", direction, got[direction], path)
		}
	}
}
