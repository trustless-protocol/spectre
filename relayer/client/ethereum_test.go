package client

import (
	"testing"
)

func TestComputeSyncCommitteePeriodAtSlot(t *testing.T) {
	cs := &EthereumClientState{
		SlotsPerEpoch:                32,
		EpochsPerSyncCommitteePeriod: 256,
	}

	tests := []struct {
		name string
		slot uint64
		want uint64
	}{
		{name: "slot 0", slot: 0, want: 0},
		{name: "last slot in first period", slot: 8191, want: 0},
		{name: "first slot in second period", slot: 8192, want: 1},
		{name: "slot 100000", slot: 100000, want: 100000 / 32 / 256},
		{name: "slot 16384", slot: 16384, want: 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := cs.ComputeSyncCommitteePeriodAtSlot(tc.slot)
			if got != tc.want {
				t.Errorf("slot %d: got period %d, want %d", tc.slot, got, tc.want)
			}
		})
	}
}

func TestToSummarizedSyncCommittee(t *testing.T) {
	t.Run("valid pubkeys with 0x prefix", func(t *testing.T) {
		sc := &SyncCommittee{
			Pubkeys:         []string{"0xaabbccdd", "0x11223344"},
			AggregatePubkey: "0xdeadbeef",
		}
		result, err := sc.ToSummarizedSyncCommittee()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.PubkeysHash == "" {
			t.Error("expected non-empty PubkeysHash")
		}
		if result.AggregatePubkey != "0xdeadbeef" {
			t.Errorf("AggregatePubkey: got %q, want %q", result.AggregatePubkey, "0xdeadbeef")
		}
	})

	t.Run("valid pubkeys without prefix", func(t *testing.T) {
		sc := &SyncCommittee{
			Pubkeys:         []string{"aabbccdd"},
			AggregatePubkey: "aggregate",
		}
		result, err := sc.ToSummarizedSyncCommittee()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.PubkeysHash == "" {
			t.Error("expected non-empty PubkeysHash")
		}
	})

	t.Run("invalid hex pubkey", func(t *testing.T) {
		sc := &SyncCommittee{
			Pubkeys: []string{"0xnothex"},
		}
		_, err := sc.ToSummarizedSyncCommittee()
		if err == nil {
			t.Fatal("expected error for invalid hex")
		}
	})

	t.Run("empty pubkeys", func(t *testing.T) {
		sc := &SyncCommittee{
			Pubkeys:         []string{},
			AggregatePubkey: "agg",
		}
		result, err := sc.ToSummarizedSyncCommittee()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.AggregatePubkey != "agg" {
			t.Errorf("AggregatePubkey: got %q", result.AggregatePubkey)
		}
	})
}

func TestToForkParameters(t *testing.T) {
	t.Run("valid spec", func(t *testing.T) {
		spec := &BeaconSpec{
			GenesisForkVersion:   "0x00000000",
			AltairForkVersion:    "0x01000000",
			AltairForkEpoch:      "74240",
			BellatrixForkVersion: "0x02000000",
			BellatrixForkEpoch:   "144896",
			CapellaForkVersion:   "0x03000000",
			CapellaForkEpoch:     "194048",
			DenebForkVersion:     "0x04000000",
			DenebForkEpoch:       "269568",
			ElectraForkVersion:   "0x05000000",
			ElectraForkEpoch:     "364544",
		}

		fp, err := spec.ToForkParameters()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fp.GenesisForkVersion != "0x00000000" {
			t.Errorf("GenesisForkVersion: got %q", fp.GenesisForkVersion)
		}
		if fp.Altair.Epoch != 74240 {
			t.Errorf("Altair.Epoch: got %d, want 74240", fp.Altair.Epoch)
		}
		if fp.Altair.Version != "0x01000000" {
			t.Errorf("Altair.Version: got %q", fp.Altair.Version)
		}
		if fp.Bellatrix.Epoch != 144896 {
			t.Errorf("Bellatrix.Epoch: got %d", fp.Bellatrix.Epoch)
		}
		if fp.Capella.Epoch != 194048 {
			t.Errorf("Capella.Epoch: got %d", fp.Capella.Epoch)
		}
		if fp.Deneb.Epoch != 269568 {
			t.Errorf("Deneb.Epoch: got %d", fp.Deneb.Epoch)
		}
		if fp.Electra.Epoch != 364544 {
			t.Errorf("Electra.Epoch: got %d", fp.Electra.Epoch)
		}
	})

	t.Run("invalid altair epoch", func(t *testing.T) {
		spec := &BeaconSpec{AltairForkEpoch: "notanumber"}
		_, err := spec.ToForkParameters()
		if err == nil {
			t.Fatal("expected error for invalid epoch")
		}
	})

	t.Run("invalid bellatrix epoch", func(t *testing.T) {
		spec := &BeaconSpec{
			AltairForkEpoch:    "0",
			BellatrixForkEpoch: "bad",
		}
		_, err := spec.ToForkParameters()
		if err == nil {
			t.Fatal("expected error for invalid epoch")
		}
	})
}
