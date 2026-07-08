package services

import (
	"testing"
	"time"

	relayerclient "relayer/client"
	"relayer/prover"
)

func TestShouldReAnchor(t *testing.T) {
	tests := []struct {
		name      string
		trusted   int64
		latest    int64
		want      bool
		wantError bool
	}{
		{name: "caught up", trusted: 100, latest: 100, want: false},
		{name: "newer block", trusted: 100, latest: 101, want: true},
		{name: "trusted ahead", trusted: 101, latest: 100, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shouldReAnchor(tt.trusted, tt.latest)
			if (err != nil) != tt.wantError {
				t.Fatalf("shouldReAnchor error = %v, wantError %v", err, tt.wantError)
			}
			if got != tt.want {
				t.Fatalf("shouldReAnchor = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordReAnchorResult(t *testing.T) {
	timestamp := &Timestamp{}
	now := time.Unix(1_700_000_000, 0)

	if err := recordReAnchorResult(timestamp, nil, now); err == nil {
		t.Fatal("expected nil light block error")
	}

	block := &relayerclient.LightBlock{BlockHeight: 123}
	if err := recordReAnchorResult(timestamp, block, now); err != nil {
		t.Fatalf("record re-anchor result: %v", err)
	}
	gotTime, gotHeight := timestamp.Snapshot()
	if !gotTime.Equal(now) || gotHeight != 123 {
		t.Fatalf("timestamp = (%v, %d), want (%v, 123)", gotTime, gotHeight, now)
	}
}

func TestRoutineBackoff(t *testing.T) {
	var backoff routineBackoff
	now := time.Unix(1_700_000_000, 0)

	if !backoff.Ready(now) {
		t.Fatal("new backoff should be ready")
	}

	if got := backoff.RecordFailure(now); got != time.Minute {
		t.Fatalf("first failure delay = %s, want 1m", got)
	}
	if backoff.Ready(now.Add(59 * time.Second)) {
		t.Fatal("backoff should not be ready before next attempt")
	}
	if !backoff.Ready(now.Add(time.Minute)) {
		t.Fatal("backoff should be ready at next attempt")
	}

	delay := time.Duration(0)
	for i := 0; i < 10; i++ {
		delay = backoff.RecordFailure(now)
	}
	if delay != routineFailureBackoffMax {
		t.Fatalf("capped delay = %s, want %s", delay, routineFailureBackoffMax)
	}

	backoff.RecordSuccess()
	if !backoff.Ready(now) {
		t.Fatal("successful attempt should reset backoff")
	}
	if got := backoff.RecordFailure(now); got != time.Minute {
		t.Fatalf("delay after reset = %s, want 1m", got)
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
