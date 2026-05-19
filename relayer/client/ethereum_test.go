package client

import (
	"crypto/sha256"
	"encoding/hex"
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

func TestComputeTimestampAtSlot(t *testing.T) {
	cs := &EthereumClientState{
		GenesisSlot:    0,
		GenesisTime:    1_700_000_000,
		SecondsPerSlot: 12,
	}

	tests := []struct {
		name string
		slot uint64
		want uint64
	}{
		{name: "genesis slot", slot: 0, want: 1_700_000_000},
		{name: "later slot", slot: 160, want: 1_700_001_920},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := cs.ComputeTimestampAtSlot(tc.slot)
			if got != tc.want {
				t.Fatalf("slot %d: got timestamp %d, want %d", tc.slot, got, tc.want)
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

// TestSyncCommitteeTreeHash verifies sszTreeHashBLSPubkeys against the known Rust test vector
// from packages/ethereum/types/src/consensus/sync_committee.rs test_sync_committee_tree_hash_root.
//
// The Rust test uses 32 pubkeys and checks SyncCommittee::tree_hash_root().
// SyncCommittee is a container with 2 fields: [pubkeys_root, aggregate_pubkey_root]
// container_root = sha256(pubkeys_root || aggregate_pubkey_root)
// expected container_root = 0x5361eb179f7499edbf09e514d317002f1d365d72e14a56c931e9edaccca3ff29
func TestSyncCommitteeTreeHash(t *testing.T) {
	pubkeys := []string{
		"0x81ea9f74ef7d935b807474e38954ae3934856219a23e074954b2e860c5a3c400f9aedb42cd27cb4ceb697ca36d1e58cb",
		"0x84d08d58c31bcd3cddf93e13d6f50203897384afa34644bff1135efe8e01c81c6a91ca6c234bb1e51ca32e41b828aaf9",
		"0xa759f6bcca8f35fcaadc406cc4b828c016c0ed23882987a79f52f2933b5cedefe24e31df6fd0d38e8a802dbafd750d01",
		"0x8d028a021c5c31a1aa1e18eda74cfaf0fba1c454c17c2e0fc730dd07a19d0c77f7a905d54017292f3e800ca06b6977cd",
		"0xb27ad13afc8ff30e087797b344c8382bb0a84447549f1b0274059ddd652276e7b148ba8808a10cc45746762957d4efbe",
		"0xa804e4fa8d1391a9d078aa93985a12503b84ce4f6f1f9e70ab7fca421e1cf972538666299d4c1bfc39327b469b2db7a8",
		"0x996323af7e545fb6363ace53f1538c7ddc3eb0d985b2479da3ee4ace10cbc393b518bf02d1a2ddb2f5bdf09b473933ea",
		"0x96947de9e6068c22a7716656a2755a9551b0b66c2d1a741bf84a088fe1e840e992dc39861bf8ba3e8d5b6d21e8f57e64",
		"0xae5302796cfeca685eaf37ffd5baeb32121f2f07415bee26cc0051ee513ff3932d2c365e3d9f87b0949a5980445cb64c",
		"0x996d10c3026b9344532b06c70a596f972a1e779a1f6106d3da9f6ba376bbf7ec82d2f52629e5dbf3f7d03b00f6b862af",
		"0xa35c6004f387430c3797ab0157af7b824c8fe106241c7cdeb897d900c0f9e4bb945ff2a6b88cbd10e35ec48aaa554ecb",
		"0xabd12678c73463ecea5867a80caf256d5c5e6ba53ff188b143a4d5be83365ad257edf39eaa1ba8753c4cdf4c632ff99e",
		"0x81fa222737fe818b43f55f209f42adaee135b2801d02709617fc88c2871852358260ace97cf323e761b5cc18bc7325b3",
		"0xab64f900c770e2b99de6b86b4390bbd1579bd48dccec55800adbcf52e006f22128e9971bbf3a92cc0105b0974849935a",
		"0x930743bfc7e18d3bd7351eaa74f477505268c1e4e1fd1ca3ccccdefb2595517343bbb8f5589c435c3c39323a4c0080f8",
		"0xab72cbc6575c3179680a58c0ecd5de46d2678ccbafc016746348ee5688edcb21b4e15bd37c70c508e3ea73103c2d566b",
		"0x84dc37ca3cd621d3da0fbdd11ca84021e0cd81a73d772dd6fcf19775b72eb64af4e573213378ccee0915dde92ac83ba6",
		"0x8d46e9aa0c1986056e407efc7013b7f271027d3c98ce96667faa98074ab0588a61681faf78644c11819a459a95689dab",
		"0xb5e898a1fc06d51c695712928f44646d15451340d1b3e480a40f03250160bc07d3b6691ec94361dd524d59d9df7f76d3",
		"0xa4ee6d37dc259cbb5237e4265429a9fd8ab5643af81628cc101e0d8b4a333ef2618a37df89ea3f92b5ea4333d8cda393",
		"0x8aa5bbee21e98c7b9e7a4c8ea45aa99f89e22992fa4fc2d73869d77da4cc8a05b25b61931ff521986677dd7f7159e8e6",
		"0x91709ee06497b9ac049325853d64947290189a8c2322e3a500d91e23ea02dc158b6db63ae558b3b7670357a151cd6071",
		"0x8fda66b8607af873f4c2c8218dd3ffc7940d411047eb199b5cd010156af4845d21dd2e65b0e44cfffb5e78271e9bb29d",
		"0xb72cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360",
		"0x896a51e0b0de0f29029af38b796db1f1e6d0f9f9085ade40a313a60cb723fa3d58f6587175570086c4fbf0fe5331f1c8",
		"0xaaf6c1251e73fb600624937760fef218aace5b253bf068ed45398aeb29d821e4d2899343ddcbbe37cb3f6cf500dff26c",
		"0x9918433b8f0bc5e126da3fdef8d7b71456492dae6d2d07f2e10c7a7f852046f84ed0ce6d3bfec42200670db27dcf3037",
		"0xa03c2a82374e04b2e0594c4ce14fb3f225b46f13188f0d8002a523c7dcfb939ae4856053c2c9c695374d7c3685df1ca5",
		"0x8d8985e5dd341c9035b37bf7391c5944c28131b47c7d5359d18fca598010ba9a63e27c55e6b421a807038c320564db17",
		"0xb24391aa97bfff29adc935d06a2b6d583433caf82f92de1980e0192d3b270323bdbf24b86dc61520a40c419dde3df4b3",
		"0xaf61f263addfb41c46d66e60ecfb598a5942f648f58718b6b4e4c92019fdb12328efbff98703134bcf28e9c1fab4bb60",
		"0xb63f327df68581cdc02a66c1c65e906a06a1a3a8d7a6e38f7b6da944e8e6cc2db85fced5327d8c12945ceb33018272ca",
	}
	aggregatePubkey := "0xa7b9141877f397e9d2a36cd86407387bbcec6d557b30ccd9e62adca217e458d7495b581e048fa1084218cadf8f45b9ff"

	// Decode pubkeys
	pks := make([][]byte, len(pubkeys))
	for i, pk := range pubkeys {
		b, err := hex.DecodeString(pk[2:]) // strip "0x"
		if err != nil {
			t.Fatalf("decode pubkey %d: %v", i, err)
		}
		pks[i] = b
	}

	// Compute pubkeys_root using our Go implementation
	pubkeysRoot := sszTreeHashBLSPubkeys(pks)

	// Compute aggregate_pubkey_root: tree_hash of [u8; 48]
	// = sha256(bytes[0..32] || bytes[32..48] || zeros[16])
	aggBytes, _ := hex.DecodeString(aggregatePubkey[2:])
	var aggChunk0, aggChunk1 [32]byte
	copy(aggChunk0[:], aggBytes[:32])
	copy(aggChunk1[:], aggBytes[32:]) // 16 bytes + 16 zeros
	h := sha256.New()
	h.Write(aggChunk0[:])
	h.Write(aggChunk1[:])
	aggRoot := [32]byte{}
	copy(aggRoot[:], h.Sum(nil))

	// Container hash: sha256(pubkeys_root || aggregate_pubkey_root)
	h2 := sha256.New()
	h2.Write(pubkeysRoot[:])
	h2.Write(aggRoot[:])
	containerRoot := h2.Sum(nil)

	expected := "5361eb179f7499edbf09e514d317002f1d365d72e14a56c931e9edaccca3ff29"
	got := hex.EncodeToString(containerRoot)
	if got != expected {
		t.Errorf("SyncCommittee tree hash mismatch:\n  got:  %s\n  want: %s", got, expected)
		t.Errorf("pubkeys_root: %s", hex.EncodeToString(pubkeysRoot[:]))
		t.Errorf("aggRoot:      %s", hex.EncodeToString(aggRoot[:]))
	}
}
