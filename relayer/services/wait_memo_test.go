package services

import "testing"

func TestDefaultConfigAppHashWait(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.AppHashWaitRetries != DEFAULT_COSMOS_APP_HASH_WAIT_RETRIES {
		t.Fatalf("AppHashWaitRetries = %d, want %d", cfg.AppHashWaitRetries, DEFAULT_COSMOS_APP_HASH_WAIT_RETRIES)
	}
	if cfg.AppHashWaitInterval != DEFAULT_COSMOS_APP_HASH_WAIT_INTERVAL {
		t.Fatalf("AppHashWaitInterval = %s, want %s", cfg.AppHashWaitInterval, DEFAULT_COSMOS_APP_HASH_WAIT_INTERVAL)
	}
}

// TestWaitBeaconFinalityMemo verifies the cross-chunk memoization fast-path:
// once a finalized execution block has been observed, a later chunk whose event
// block is already covered returns immediately without hitting the beacon RPC
// (issue #76 #3). The early return happens before `ctx` is dereferenced, so a
// zero Context is safe here.
func TestWaitBeaconFinalityMemo(t *testing.T) {
	s := &Services{lastFinalizedExecBlock: 100}

	if !s.waitBeaconFinality(Context{}, 100, "test") {
		t.Fatal("expected immediate success when event block == cached finalized block")
	}
	if !s.waitBeaconFinality(Context{}, 50, "test") {
		t.Fatal("expected immediate success when event block < cached finalized block")
	}
}

// TestWaitCosmosAppHashMemo verifies the AppHash wait fast-path: a target height
// already confirmed by an earlier chunk returns without querying Cosmos status
// (issue #76 #1). The early return precedes any `ctx` use, so a zero Context is
// safe. Reaching the RPC path with a zero Context would panic, so a clean return
// is itself the assertion that the fast-path was taken.
func TestWaitCosmosAppHashMemo(t *testing.T) {
	s := &Services{lastCosmosAppHashHeight: 200}

	if !s.waitCosmosAppHash(Context{}, 200) { // target == cached
		t.Fatal("expected immediate success when target == cached AppHash height")
	}
	if !s.waitCosmosAppHash(Context{}, 150) { // target < cached
		t.Fatal("expected immediate success when target < cached AppHash height")
	}
}
