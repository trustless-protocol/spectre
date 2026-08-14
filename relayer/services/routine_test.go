package services

import (
	"context"
	"errors"
	"testing"
	"time"

	relayerclient "relayer/client"
	"relayer/prover"
)

// TestWaitForCosmosCatchUp_AbortsOnCancelledContext verifies the catch-up wait
// returns promptly when the context is already cancelled, instead of polling the
// Cosmos client (which is never reached — a zero endpoint would panic if it were).
func TestWaitForCosmosCatchUp_AbortsOnCancelledContext(t *testing.T) {
	w := &Worker{}
	stdCtx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	done := make(chan struct{})
	go func() {
		// A zero endpoint set is fine: the cancelled stdCtx short-circuits
		// before ctx.CosmosClient() is ever called.
		w.WaitForCosmosCatchUp(stdCtx, CosmosEndpoint{}, &relayerclient.EthereumClientState{}, 0)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("WaitForCosmosCatchUp did not abort on a cancelled context")
	}
}

func TestDeriveCosmosRefreshInterval(t *testing.T) {
	trustingPeriod := 16 * time.Hour

	got, err := deriveCosmosRefreshInterval(Config{
		RefreshInterval: DEFAULT_REFRESH_INTERVAL,
	}, trustingPeriod)
	if err != nil {
		t.Fatalf("derive default interval: %v", err)
	}
	if got != 15*time.Hour {
		t.Fatalf("default interval = %s, want 15h", got)
	}

	_, err = deriveCosmosRefreshInterval(Config{
		RefreshInterval:           15 * time.Hour,
		RefreshIntervalConfigured: true,
	}, trustingPeriod)
	if err == nil {
		t.Fatal("expected configured interval at safety boundary to be rejected")
	}

	got, err = deriveCosmosRefreshInterval(Config{
		RefreshInterval:           14 * time.Hour,
		RefreshIntervalConfigured: true,
	}, trustingPeriod)
	if err != nil {
		t.Fatalf("derive configured safe interval: %v", err)
	}
	if got != 14*time.Hour {
		t.Fatalf("configured interval = %s, want 14h", got)
	}
}

func TestCosmosCurrentSlotReady(t *testing.T) {
	if cosmosCurrentSlotReady(932, 933) {
		t.Fatal("did not expect current slot below safety margin to be ready")
	}

	if cosmosCurrentSlotReady(935, 933) {
		t.Fatal("did not expect current slot equal to the minimum safety slot to be treated as not ready")
	}

	if !cosmosCurrentSlotReady(936, 933) {
		t.Fatal("expected current slot above safety margin to be ready")
	}
}

func TestSelectSignaturesForPinnedSetChoosesPinnedQuorum(t *testing.T) {
	pubkey0 := [32]byte{0x01}
	pubkey1 := [32]byte{0x02}
	pubkey2 := [32]byte{0x03}
	pubkey3 := [32]byte{0x04}

	candidates := []prover.ValidatorSignature{
		{Index: 0, PublicKey: pubkey0[:], Power: 40, Active: true},
		{Index: 1, PublicKey: pubkey1[:], Power: 40, Active: true},
		{Index: 2, PublicKey: pubkey2[:], Power: 30, Active: true},
		{Index: 3, PublicKey: pubkey3[:], Power: 30, Active: true},
	}
	pinned := pinnedCosmosValidatorSet{
		indices:      []uint32{0, 1, 2},
		pubkeys:      [][32]byte{pubkey2, pubkey3, pubkey0},
		votingPowers: []uint64{40, 40, 20},
		totalPower:   100,
	}

	selected, err := selectSignaturesForPinnedSet(candidates, pinned)
	if err != nil {
		t.Fatalf("select signatures: %v", err)
	}

	got := make([]int, len(selected))
	for i, sig := range selected {
		got[i] = sig.Index
	}
	want := []int{2, 3}
	if len(got) != len(want) {
		t.Fatalf("selected indices: got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("selected indices: got %v want %v", got, want)
		}
	}
}

func TestPinnedOverlapPower(t *testing.T) {
	pubkey0 := [32]byte{0x01}
	pubkey1 := [32]byte{0x02}
	pubkey2 := [32]byte{0x03}

	candidates := []prover.ValidatorSignature{
		{Index: 0, PublicKey: pubkey0[:], Power: 40, Active: true},
		{Index: 1, PublicKey: pubkey1[:], Power: 40, Active: true}, // not pinned: contributes 0
	}
	pinned := pinnedCosmosValidatorSet{
		indices:      []uint32{0, 1},
		pubkeys:      [][32]byte{pubkey0, pubkey2},
		votingPowers: []uint64{30, 70},
		totalPower:   100,
	}

	if got := pinnedOverlapPower(candidates, pinned); got != 30 {
		t.Fatalf("overlap = %d, want 30 (pinned power of pubkey0 only)", got)
	}
	if got := pinnedOverlapPower(nil, pinned); got != 0 {
		t.Fatalf("overlap with no candidates = %d, want 0", got)
	}
}

func TestShouldRotatePinnedSet(t *testing.T) {
	threshold := relayerclient.TrustThreshold{Numerator: 5, Denominator: 6}

	if !shouldRotatePinnedSet(true, 100, 100, threshold) {
		t.Fatal("forceRotation should always rotate")
	}
	if shouldRotatePinnedSet(false, 90, 100, threshold) {
		t.Fatal("overlap 90/100 above 5/6 should not rotate")
	}
	// boundary: overlap/total exactly 5/6 (100*6 == 120*5)
	if !shouldRotatePinnedSet(false, 100, 120, threshold) {
		t.Fatal("overlap exactly at 5/6 should rotate")
	}
	if !shouldRotatePinnedSet(false, 70, 100, threshold) {
		t.Fatal("overlap 70/100 below 5/6 should rotate")
	}

	// threshold 1/1 rotates on any overlap (overlap <= total always holds)
	full := relayerclient.TrustThreshold{Numerator: 1, Denominator: 1}
	if !shouldRotatePinnedSet(false, 100, 100, full) {
		t.Fatal("threshold 1/1 should rotate even at full overlap")
	}
}

func TestParseRotationThreshold(t *testing.T) {
	def, err := ParseRotationThreshold("")
	if err != nil {
		t.Fatalf("empty value should fall back to default: %v", err)
	}
	if def.Numerator != 5 || def.Denominator != 6 {
		t.Fatalf("default = %d/%d, want 5/6", def.Numerator, def.Denominator)
	}

	for _, ok := range []string{"3/4", "1/1", "5/6"} {
		if _, err := ParseRotationThreshold(ok); err != nil {
			t.Fatalf("expected %q to be valid: %v", ok, err)
		}
	}
	for _, bad := range []string{"2/3", "1/2", "7/6", "0/1", "1/0", "abc", "5"} {
		if _, err := ParseRotationThreshold(bad); err == nil {
			t.Fatalf("expected %q to be rejected", bad)
		}
	}
}

func TestSelectSignaturesForPinnedSetRejectsInsufficientOverlap(t *testing.T) {
	pubkey0 := [32]byte{0x01}
	pubkey1 := [32]byte{0x02}

	candidates := []prover.ValidatorSignature{
		{Index: 0, PublicKey: pubkey0[:], Power: 40, Active: true},
	}
	pinned := pinnedCosmosValidatorSet{
		indices:      []uint32{0, 1},
		pubkeys:      [][32]byte{pubkey0, pubkey1},
		votingPowers: []uint64{40, 60},
		totalPower:   100,
	}

	if _, err := selectSignaturesForPinnedSet(candidates, pinned); err == nil {
		t.Fatal("expected insufficient pinned quorum error")
	}
}

// TestBinarySearchHighestFeasibleConverges is the RLY-01 multi-hop core: given
// a synthetic monotonic feasibility predicate (feasible(h) == h <= threshold),
// the search must find exactly the highest height in (trusted, latest] that
// clears it, including the boundary cases (threshold at trusted+1, at
// latest-1, and "no height clears it at all").
func TestBinarySearchHighestFeasibleConverges(t *testing.T) {
	cases := []struct {
		name            string
		trusted, latest int64
		threshold       int64
		wantHeight      int64
		wantOk          bool
	}{
		{"mid-range threshold", 100, 200, 150, 150, true},
		{"threshold at trusted+1", 100, 200, 101, 101, true},
		{"threshold at latest-1", 100, 200, 199, 199, true},
		{"no height clears it", 100, 200, 100, 0, false},
		{"single-block gap, target already infeasible", 100, 101, 101, 0, false},
		{"large gap", 1_000, 1_000_000, 654_321, 654_321, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			feasible := func(h int64) (bool, error) { return h <= tc.threshold, nil }
			height, ok, err := binarySearchHighestFeasible(tc.trusted, tc.latest, feasible)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok != tc.wantOk {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOk)
			}
			if ok && height != tc.wantHeight {
				t.Fatalf("height = %d, want %d", height, tc.wantHeight)
			}
		})
	}
}

func TestBinarySearchHighestFeasibleDoesNotReprobeKnownInfeasibleEdges(t *testing.T) {
	calls := map[int64]int{}
	feasible := func(h int64) (bool, error) {
		calls[h]++
		return false, nil
	}
	height, ok, err := binarySearchHighestFeasible(100, 102, feasible)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok || height != 0 {
		t.Fatalf("height=%d ok=%v, want no feasible intermediate", height, ok)
	}
	if calls[102] != 0 {
		t.Fatalf("latest height was re-probed %d times", calls[102])
	}
	if calls[101] != 1 {
		t.Fatalf("trusted+1 was probed %d times, want exactly once", calls[101])
	}

	calls = map[int64]int{}
	height, ok, err = binarySearchHighestFeasible(100, 101, feasible)
	if err != nil {
		t.Fatalf("unexpected single-block error: %v", err)
	}
	if ok || height != 0 {
		t.Fatalf("single-block height=%d ok=%v, want no feasible intermediate", height, ok)
	}
	if len(calls) != 0 {
		t.Fatalf("single-block gap re-probed known failed target: calls=%v", calls)
	}
}

// TestBinarySearchHighestFeasiblePropagatesError verifies a genuine probe
// error (e.g. an RPC failure fetching a candidate light block) is returned to
// the caller rather than being misread as "infeasible at this height" — the
// two must never be conflated, since one means "churn beat quorum here" and
// the other means "we don't actually know."
func TestBinarySearchHighestFeasiblePropagatesError(t *testing.T) {
	wantErr := errors.New("boom: rpc unavailable")
	feasible := func(h int64) (bool, error) {
		if h == 150 { // the first midpoint probed for (100, 200]
			return false, wantErr
		}
		return h <= 190, nil
	}
	_, _, err := binarySearchHighestFeasible(100, 200, feasible)
	if err == nil {
		t.Fatal("expected error to propagate")
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected wrapped %v, got %v", wantErr, err)
	}
}
