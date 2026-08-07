package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
)

const testSlotsPerEpoch = 32

// fakeBeacon serves the two endpoints resolveBootstrapCheckpoint uses. A slot is
// "servable" only if it is listed in servable; every other slot 404s the same way a
// real node does for a skipped checkpoint. Requested slots are recorded so a test can
// assert exactly which boundaries were walked.
//
// noBootstrap models the other, more common failure: the block exists (so the root
// lookup succeeds) but the node has not retained a light-client bootstrap for it.
// The two endpoints fail independently on a real node, so the fake must too.
type fakeBeacon struct {
	servable    map[uint64]bool
	noBootstrap map[uint64]bool

	mu            sync.Mutex
	rootsAsked    []uint64
	bootstrapsFor []uint64
}

// rootFor is a deterministic stand-in for a beacon block root, and encodes the slot so
// the bootstrap handler can map a root back to the slot it belongs to.
func rootFor(slot uint64) string { return fmt.Sprintf("0x%064x", slot) }

func (f *fakeBeacon) start(t *testing.T) string {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/eth/v1/beacon/blocks/", func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/eth/v1/beacon/blocks/"), "/root")
		slot, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			http.Error(w, "bad slot", http.StatusBadRequest)
			return
		}
		f.mu.Lock()
		f.rootsAsked = append(f.rootsAsked, slot)
		f.mu.Unlock()
		// Skipped slots have no block at all, so even the root lookup fails.
		if !f.servable[slot] {
			http.Error(w, "NOT_FOUND: no block at slot", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, `{"data":{"root":%q}}`, rootFor(slot))
	})

	mux.HandleFunc("/eth/v1/beacon/light_client/bootstrap/", func(w http.ResponseWriter, r *http.Request) {
		root := strings.TrimPrefix(r.URL.Path, "/eth/v1/beacon/light_client/bootstrap/")
		slot, err := strconv.ParseUint(strings.TrimPrefix(root, "0x"), 16, 64)
		if err != nil {
			http.Error(w, "bad root", http.StatusBadRequest)
			return
		}
		f.mu.Lock()
		f.bootstrapsFor = append(f.bootstrapsFor, slot)
		f.mu.Unlock()
		// The block exists but its bootstrap was pruned — the walk must treat this the
		// same as a missing block and step back.
		if f.noBootstrap[slot] {
			http.Error(w, "NOT_FOUND: no light client bootstrap for this block root", http.StatusNotFound)
			return
		}
		// The aggregate pubkey carries the slot so the test can tell which checkpoint
		// the returned bootstrap actually came from.
		fmt.Fprintf(w, `{"data":{"current_sync_committee":{"pubkeys":[],"aggregate_pubkey":%q}}}`, rootFor(slot))
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv.URL
}

func (f *fakeBeacon) asked() []uint64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]uint64(nil), f.rootsAsked...)
}

// TestResolveBootstrapCheckpointUsesEpochBoundary: a finalized header reported
// mid-epoch must be rounded DOWN to its epoch boundary, because that is the slot the
// beacon indexes bootstraps by.
func TestResolveBootstrapCheckpointUsesEpochBoundary(t *testing.T) {
	beacon := &fakeBeacon{servable: map[uint64]bool{4096: true}}
	url := beacon.start(t)

	slot, root, bootstrap, err := resolveBootstrapCheckpoint(url, "4127", testSlotsPerEpoch)
	if err != nil {
		t.Fatalf("resolveBootstrapCheckpoint: %v", err)
	}
	if slot != "4096" {
		t.Errorf("checkpoint slot = %q, want 4096 (the epoch boundary below 4127)", slot)
	}
	if root != rootFor(4096) {
		t.Errorf("root = %q, want %q", root, rootFor(4096))
	}
	if got := bootstrap.Data.CurrentSyncCommittee.AggregatePubkey; got != rootFor(4096) {
		t.Errorf("bootstrap came from the wrong checkpoint: aggregate_pubkey = %q", got)
	}
	if asked := beacon.asked(); len(asked) != 1 || asked[0] != 4096 {
		t.Errorf("asked for roots %v, want exactly [4096] — no step-back was needed", asked)
	}
}

// TestResolveBootstrapCheckpointStepsBackOverSkippedBoundary is the bug this fix is
// for: the boundary slot was skipped, so it has no block and no bootstrap. The walk
// must fall back to the previous boundary rather than failing client creation.
func TestResolveBootstrapCheckpointStepsBackOverSkippedBoundary(t *testing.T) {
	// 4096 skipped, 4064 servable.
	beacon := &fakeBeacon{servable: map[uint64]bool{4064: true}}
	url := beacon.start(t)

	slot, root, bootstrap, err := resolveBootstrapCheckpoint(url, "4100", testSlotsPerEpoch)
	if err != nil {
		t.Fatalf("resolveBootstrapCheckpoint: %v", err)
	}
	if slot != "4064" {
		t.Errorf("checkpoint slot = %q, want 4064 (one epoch back from the skipped 4096)", slot)
	}
	if root != rootFor(4064) {
		t.Errorf("root = %q, want %q", root, rootFor(4064))
	}
	if got := bootstrap.Data.CurrentSyncCommittee.AggregatePubkey; got != rootFor(4064) {
		t.Errorf("bootstrap came from the wrong checkpoint: aggregate_pubkey = %q", got)
	}
	if asked := beacon.asked(); len(asked) != 2 || asked[0] != 4096 || asked[1] != 4064 {
		t.Errorf("asked for roots %v, want [4096 4064]", asked)
	}
}

// TestResolveBootstrapCheckpointStepsBackWhenBootstrapPruned covers the failure that
// motivated the walk in the first place, and the one a real node produces most often:
// the checkpoint block is there — the root lookup succeeds — but the node no longer
// retains a light-client bootstrap for it. The step-back must trigger on the bootstrap
// error, not only on a missing block.
func TestResolveBootstrapCheckpointStepsBackWhenBootstrapPruned(t *testing.T) {
	// Both boundaries have blocks; only 4064 still has a bootstrap.
	beacon := &fakeBeacon{
		servable:    map[uint64]bool{4096: true, 4064: true},
		noBootstrap: map[uint64]bool{4096: true},
	}
	url := beacon.start(t)

	slot, root, bootstrap, err := resolveBootstrapCheckpoint(url, "4096", testSlotsPerEpoch)
	if err != nil {
		t.Fatalf("resolveBootstrapCheckpoint: %v", err)
	}
	if slot != "4064" {
		t.Errorf("checkpoint slot = %q, want 4064 (4096 has a block but no bootstrap)", slot)
	}
	if root != rootFor(4064) {
		t.Errorf("root = %q, want %q", root, rootFor(4064))
	}
	// The returned triple must be self-consistent: the caller checks this bootstrap
	// against the beacon block at the returned slot, so a mismatched pair would only
	// surface later as a confusing block-number error.
	if got := bootstrap.Data.CurrentSyncCommittee.AggregatePubkey; got != rootFor(4064) {
		t.Errorf("bootstrap came from the wrong checkpoint: aggregate_pubkey = %q", got)
	}
	// 4096's root was fetched and its bootstrap attempted before stepping back.
	if asked := beacon.asked(); len(asked) != 2 || asked[0] != 4096 || asked[1] != 4064 {
		t.Errorf("asked for roots %v, want [4096 4064]", asked)
	}
}

// TestResolveBootstrapCheckpointGivesUpAfterCap: a beacon that serves no bootstrap at
// all must produce one bounded, explanatory failure — not an unbounded walk.
func TestResolveBootstrapCheckpointGivesUpAfterCap(t *testing.T) {
	beacon := &fakeBeacon{servable: map[uint64]bool{}}
	url := beacon.start(t)

	_, _, _, err := resolveBootstrapCheckpoint(url, "4100", testSlotsPerEpoch)
	if err == nil {
		t.Fatal("expected an error when no checkpoint is servable")
	}
	if !strings.Contains(err.Error(), "no servable light-client bootstrap") {
		t.Errorf("error does not explain the failure: %v", err)
	}
	if asked := beacon.asked(); len(asked) != maxBootstrapCheckpointStepBack {
		t.Errorf("asked for %d roots, want the %d-attempt cap: %v",
			len(asked), maxBootstrapCheckpointStepBack, asked)
	}
}

// TestResolveBootstrapCheckpointStopsAtGenesisEpoch guards the unsigned-slot subtraction:
// inside the genesis epoch there is no earlier boundary, so the walk must stop at slot 0
// instead of wrapping around to an absurd slot number.
func TestResolveBootstrapCheckpointStopsAtGenesisEpoch(t *testing.T) {
	beacon := &fakeBeacon{servable: map[uint64]bool{}}
	url := beacon.start(t)

	_, _, _, err := resolveBootstrapCheckpoint(url, "17", testSlotsPerEpoch)
	if err == nil {
		t.Fatal("expected an error when no checkpoint is servable")
	}
	asked := beacon.asked()
	if len(asked) != 1 || asked[0] != 0 {
		t.Fatalf("asked for roots %v, want exactly [0] — the walk must not go below genesis", asked)
	}
}

// TestResolveBootstrapCheckpointRejectsBadInput covers the two argument errors, both of
// which would otherwise divide by zero or walk from a garbage slot.
func TestResolveBootstrapCheckpointRejectsBadInput(t *testing.T) {
	beacon := &fakeBeacon{servable: map[uint64]bool{}}
	url := beacon.start(t)

	if _, _, _, err := resolveBootstrapCheckpoint(url, "not-a-slot", testSlotsPerEpoch); err == nil {
		t.Error("expected an error for an unparseable finalized slot")
	}
	if _, _, _, err := resolveBootstrapCheckpoint(url, "4100", 0); err == nil {
		t.Error("expected an error for slots_per_epoch = 0")
	}
	if asked := beacon.asked(); len(asked) != 0 {
		t.Errorf("bad input reached the beacon: %v", asked)
	}
}
