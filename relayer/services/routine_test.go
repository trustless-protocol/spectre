package services

import (
	"testing"
	"time"

	commettypes "github.com/cometbft/cometbft/types"

	relayerclient "relayer/client"
	"relayer/prover"
)

func TestRecordRefreshResult(t *testing.T) {
	timestamp := &Timestamp{}
	now := time.Unix(1_700_000_000, 0)
	trustedTime := now.Add(-10 * time.Minute)

	if err := recordRefreshResult(timestamp, nil); err == nil {
		t.Fatal("expected nil light block error")
	}

	if err := recordRefreshResult(timestamp, &relayerclient.LightBlock{BlockHeight: 123}); err == nil {
		t.Fatal("expected missing header error")
	}

	block := &relayerclient.LightBlock{
		BlockHeight: 123,
		SignedHeader: commettypes.SignedHeader{
			Header: &commettypes.Header{Time: trustedTime},
		},
	}
	if err := recordRefreshResult(timestamp, block); err != nil {
		t.Fatalf("record refresh result: %v", err)
	}
	gotTime, gotHeight := timestamp.Snapshot()
	if !gotTime.Equal(trustedTime) || gotHeight != 123 {
		t.Fatalf("timestamp = (%v, %d), want (%v, 123)", gotTime, gotHeight, trustedTime)
	}
}

func TestRecordEthClientUpdateResult(t *testing.T) {
	timestamp := &Timestamp{}
	proofTime := uint64(1_700_000_123)
	result := &EthClientUpdateResult{
		EthClientState: &relayerclient.EthereumClientState{
			LatestSlot: 42,
		},
		ProofTimestamp: proofTime,
	}

	if err := recordEthClientUpdateResult(timestamp, nil); err == nil {
		t.Fatal("expected nil ethereum client update result error")
	}
	if err := recordEthClientUpdateResult(timestamp, &EthClientUpdateResult{}); err == nil {
		t.Fatal("expected missing ethereum client state error")
	}
	if err := recordEthClientUpdateResult(timestamp, &EthClientUpdateResult{
		EthClientState: &relayerclient.EthereumClientState{LatestSlot: 42},
	}); err == nil {
		t.Fatal("expected zero proof timestamp error")
	}

	if err := recordEthClientUpdateResult(timestamp, result); err != nil {
		t.Fatalf("record eth client update result: %v", err)
	}
	gotTime, gotSlot := timestamp.Snapshot()
	wantTime := time.Unix(int64(proofTime), 0)
	if !gotTime.Equal(wantTime) || gotSlot != 42 {
		t.Fatalf("timestamp = (%v, %d), want (%v, 42)", gotTime, gotSlot, wantTime)
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
