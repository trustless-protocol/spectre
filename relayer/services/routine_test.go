package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	updateClientContract "relayer/bindings/UpdateClient"
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

// cloneEthereumClientState is a one-line shallow copy: `cloned := *state`. That is
// a correct deep copy today only because every field of EthereumClientState is a
// value type. Add one slice, map or pointer -- a list of forks, say -- and the
// "clone" starts sharing it, so a mutation through the copy reaches the original.
//
// The callers rely on independence: ethProofStateFromFinalityUpdate clones the
// on-chain client state and then rewrites slots on the copy to build a proof
// state, while the original stays the trusted one.
//
// Testing that by mutating fields would prove nothing until the day someone adds
// a reference field, and would pass right up to it. This checks the property the
// shallow copy actually depends on, so the failure arrives with the change that
// breaks it rather than with the bug it later causes.
func TestCloneEthereumClientState_StaysADeepCopy(t *testing.T) {
	assertNoReferenceFields(t, reflect.TypeOf(relayerclient.EthereumClientState{}), "EthereumClientState")
}

func assertNoReferenceFields(t *testing.T, typ reflect.Type, path string) {
	t.Helper()
	switch typ.Kind() {
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			assertNoReferenceFields(t, typ.Field(i).Type, path+"."+typ.Field(i).Name)
		}
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Chan, reflect.Func, reflect.Interface:
		t.Errorf("%s is a %s, so `cloned := *state` now shares it with the original;"+
			" cloneEthereumClientState must copy it explicitly", path, typ.Kind())
	}
}

func TestCloneEthereumClientState_NilInNilOut(t *testing.T) {
	if got := cloneEthereumClientState(nil); got != nil {
		t.Fatalf("cloning nil returned %+v; the caller checks for nil to report a missing state", got)
	}
}

// The clone must not be the same pointer, or every "proof state" edit would land
// on the trusted state the caller kept.
func TestCloneEthereumClientState_ReturnsANewPointer(t *testing.T) {
	original := &relayerclient.EthereumClientState{LatestSlot: 100}
	cloned := cloneEthereumClientState(original)
	if cloned == original {
		t.Fatal("clone returned the same pointer")
	}
	cloned.LatestSlot = 200
	if original.LatestSlot != 100 {
		t.Fatalf("writing to the clone changed the original: LatestSlot = %d", original.LatestSlot)
	}
}

// The on-chain height is uint64 and the Tendermint RPC takes int64. The guard is
// the conversion boundary; without it a height above 2^63-1 wraps to a negative
// number, and a negative height reads to the RPC as "latest" rather than as an
// error.
func TestClientStateRevisionHeightInt64(t *testing.T) {
	const maxInt64 = uint64(1)<<63 - 1

	tests := []struct {
		name    string
		height  uint64
		want    int64
		wantErr bool
	}{
		{"zero", 0, 0, false},
		{"ordinary height", 12345, 12345, false},
		{"largest representable", maxInt64, int64(maxInt64), false},
		{"one past int64", maxInt64 + 1, 0, true},
		{"max uint64", ^uint64(0), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := clientStateRevisionHeightInt64(relayerclient.ClientState{
				LatestHeight: updateClientContract.IICS02ClientMsgsHeight{RevisionHeight: tt.height},
			})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("height %d was accepted; it does not fit in an int64", tt.height)
				}
				return
			}
			if err != nil {
				t.Fatalf("height %d rejected: %v", tt.height, err)
			}
			if got != tt.want {
				t.Fatalf("height = %d, want %d", got, tt.want)
			}
		})
	}
}

// parseSlot reads beacon slots out of JSON the relayer did not produce -- all
// five call sites pass a field straight from a beacon API response.
func TestParseSlot(t *testing.T) {
	t.Run("parses decimal slots", func(t *testing.T) {
		for _, tt := range []struct {
			in   string
			want uint64
		}{
			{"0", 0},
			{"123", 123},
			{"18446744073709551615", ^uint64(0)}, // max uint64
		} {
			got, err := parseSlot(tt.in)
			if err != nil {
				t.Errorf("parseSlot(%q): %v", tt.in, err)
				continue
			}
			if got != tt.want {
				t.Errorf("parseSlot(%q) = %d, want %d", tt.in, got, tt.want)
			}
		}
	})

	// Everything else is rejected. The three middle cases used to be ACCEPTED:
	// fmt.Sscanf("%d") stopped at the first non-digit and reported success on the
	// prefix, so "123abc" became 123, "1e3" became 1, and "0x10" became 0 -- which
	// downstream reads as genesis. A slot decides which state a proof is built
	// against, so a wrong-but-plausible number is the worst outcome available.
	t.Run("rejects anything that is not a plain decimal", func(t *testing.T) {
		for name, in := range map[string]string{
			"empty":               "",
			"negative":            "-5",
			"overflows uint64":    "99999999999999999999999",
			"digits then letters": "123abc",
			"scientific notation": "1e3",
			"hex":                 "0x10",
			"leading whitespace":  "  7",
			"trailing whitespace": "7 ",
			"decimal point":       "7.0",
			"thousands separator": "1,000",
		} {
			t.Run(name, func(t *testing.T) {
				got, err := parseSlot(in)
				if err == nil {
					t.Fatalf("parseSlot(%q) returned %d with no error", in, got)
				}
				// The value must be unusable, not merely accompanied by an error:
				// a caller that forgets to check err must not silently get genesis
				// back as a plausible slot... which is exactly what it would get.
				// Zero is the only honest answer here, so the error is the contract.
				if got != 0 {
					t.Fatalf("parseSlot(%q) = %d alongside an error", in, got)
				}
				// The message must carry BOTH halves: strconv quotes the offending
				// input, and the wrapper says which kind of value it was. Without
				// the wrapper an operator sees a bare parse error and has to guess
				// which field of which beacon response produced it.
				if in != "" && !strings.Contains(err.Error(), in) {
					t.Errorf("error %q does not quote the offending input %q", err, in)
				}
				if !strings.Contains(err.Error(), "beacon slot") {
					t.Errorf("error %q does not say what kind of value failed to parse", err)
				}
			})
		}
	})
}

// A slot string the relayer cannot parse must not be reported as slot 0, which
// downstream reads as "genesis" and would rebuild the client from scratch.
func TestParseSlot_ErrorDoesNotLookLikeGenesis(t *testing.T) {
	got, err := parseSlot("")
	if err == nil {
		t.Fatal("an empty slot string was accepted")
	}
	if got != 0 {
		t.Fatalf("slot = %d alongside an error; callers must not use it", got)
	}
	if !strings.Contains(err.Error(), "EOF") && !strings.Contains(err.Error(), "unexpected") {
		t.Logf("error text is %q -- fine, but callers must check err, not the value", err)
	}
}

// Which of the two stored committees the client checks against is decided by the
// label on the header, and the client picks it by the SIGNATURE slot's period:
//
//	sync_committee = if signature_period == stored_period { current } else { next }
//
// (ethereum/light-client verify.rs). Getting the label wrong fails on chain with
// "current sync committee (X) does not match with the one in the current state (Y)"
// -- a message that names two aggregate pubkeys and neither period.
func TestActiveCommitteeFor(t *testing.T) {
	committee := &relayerclient.SyncCommittee{AggregatePubkey: "0xaggregate"}

	cases := []struct {
		name            string
		signaturePeriod uint64
		storedPeriod    uint64
		wantCurrent     bool
		wantNext        bool
		wantErr         bool
		why             string
	}{
		{
			name:            "same period is the current committee",
			signaturePeriod: 1343,
			storedPeriod:    1343,
			wantCurrent:     true,
			why:             "the steady state: the signature is in the period the client is finalized in",
		},
		{
			name:            "one period ahead is the next committee",
			signaturePeriod: 1344,
			storedPeriod:    1343,
			wantNext:        true,
			why: "the last ~65 slots of every period: signature_slot has crossed while the client's " +
				"finalized slot has not. Labelling this Current is the stall",
		},
		{
			name:            "two periods ahead is refused",
			signaturePeriod: 1345,
			storedPeriod:    1343,
			wantErr:         true,
			why:             "the client rejects it outright, so building the header only spends gas to be told so",
		},
		{
			name:            "a period behind is refused",
			signaturePeriod: 1342,
			storedPeriod:    1343,
			wantErr:         true,
			why:             "the client has no committee older than its current one to check against",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			active, err := activeCommitteeFor(tc.signaturePeriod, tc.storedPeriod, committee)
			if (err != nil) != tc.wantErr {
				t.Fatalf("activeCommitteeFor(%d, %d) error = %v, want error %v: %s",
					tc.signaturePeriod, tc.storedPeriod, err, tc.wantErr, tc.why)
			}
			if tc.wantErr {
				return
			}
			if (active.Current != nil) != tc.wantCurrent {
				t.Fatalf("Current set = %v, want %v: %s", active.Current != nil, tc.wantCurrent, tc.why)
			}
			if (active.Next != nil) != tc.wantNext {
				t.Fatalf("Next set = %v, want %v: %s", active.Next != nil, tc.wantNext, tc.why)
			}
			// Exactly one, never both: the client reads one field and the other
			// would be silently ignored.
			if (active.Current != nil) == (active.Next != nil) {
				t.Fatal("exactly one of Current/Next must be set")
			}
		})
	}

	// The boundary is a real slot pair, not a rounded one: signature_slot is
	// attested_slot + 1, so a single slot decides the label once every 8192.
	t.Run("the boundary falls on one real slot", func(t *testing.T) {
		const slotsPerPeriod = 8192 // EpochsPerSyncCommitteePeriod(256) * SlotsPerEpoch(32)
		cs := &relayerclient.EthereumClientState{
			EpochsPerSyncCommitteePeriod: 256,
			SlotsPerEpoch:                32,
		}

		// Period 1344 starts at 8192 * 1344 = 11010048. The attested slot one below it
		// is still 1343 while its signature slot is already 1344 -- the one slot where
		// choosing by attested rather than signature gives the wrong answer.
		const lastSlotOf1343 = slotsPerPeriod*1344 - 1

		if got := cs.ComputeSyncCommitteePeriodAtSlot(lastSlotOf1343); got != 1343 {
			t.Fatalf("attested slot %d is in period %d, want 1343", lastSlotOf1343, got)
		}
		if got := cs.ComputeSyncCommitteePeriodAtSlot(lastSlotOf1343 + 1); got != 1344 {
			t.Fatalf("signature slot %d is in period %d, want 1344", lastSlotOf1343+1, got)
		}

		// Reading the period from the attested slot would label this Current against a
		// client stored in 1343 -- which is what the on-chain rejection looked like.
		byAttested, err := activeCommitteeFor(cs.ComputeSyncCommitteePeriodAtSlot(lastSlotOf1343), 1343, &relayerclient.SyncCommittee{})
		if err != nil || byAttested.Current == nil {
			t.Fatalf("attested-slot period must resolve to Current, got %+v err=%v", byAttested, err)
		}
		bySignature, err := activeCommitteeFor(cs.ComputeSyncCommitteePeriodAtSlot(lastSlotOf1343+1), 1343, &relayerclient.SyncCommittee{})
		if err != nil || bySignature.Next == nil {
			t.Fatalf("signature-slot period must resolve to Next, got %+v err=%v", bySignature, err)
		}
	})
}
