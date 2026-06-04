package services

import (
	"testing"

	updateclientContract "relayer/bindings/UpdateClient"
)

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

func TestDetectValidatorSetDelta(t *testing.T) {
	baseHash := [32]byte{0x88, 0xbe}
	pubkey0 := [32]byte{0x01}
	pubkey1 := [32]byte{0x02}
	pubkey2 := [32]byte{0x03}
	baseSet := cachedCosmosValidatorSet{
		indices:      []uint32{0, 1, 2},
		pubkeys:      [][32]byte{pubkey0, pubkey1, pubkey2},
		votingPowers: []uint64{100, 100, 100},
	}

	currentSet := updateclientContract.IICS07TendermintMsgsValidatorSet{
		Validators: []updateclientContract.IICS07TendermintMsgsValidatorInfo{
			{PubKey: pubkey0, VotingPower: 100},
			{PubKey: pubkey1, VotingPower: 125},
			{PubKey: pubkey2, VotingPower: 100},
		},
	}

	delta, ok, reason := detectValidatorSetDelta(baseHash, baseSet, currentSet)
	if !ok {
		t.Fatalf("expected delta, got reason=%q", reason)
	}
	if delta.BaseValidatorsHash != baseHash {
		t.Fatalf("base hash: got %x want %x", delta.BaseValidatorsHash, baseHash)
	}
	if delta.LeafCount != 1 {
		t.Fatalf("leaf count: got %d want 1", delta.LeafCount)
	}
	if delta.Indices[0] != 1 {
		t.Fatalf("changed index: got %d want 1", delta.Indices[0])
	}
	if delta.PubKeys[0] != pubkey1 {
		t.Fatalf("new pubkey: got %x want %x", delta.PubKeys[0], pubkey1)
	}
	if delta.VotingPowers[0] != 125 {
		t.Fatalf("new voting power: got %d want 125", delta.VotingPowers[0])
	}
}

func TestDetectValidatorSetDeltaMultiLeafReorder(t *testing.T) {
	baseHash := [32]byte{0x88, 0xbe}
	pubkey0 := [32]byte{0x01}
	pubkey1 := [32]byte{0x02}
	pubkey2 := [32]byte{0x03}
	baseSet := cachedCosmosValidatorSet{
		indices:      []uint32{0, 1, 2},
		pubkeys:      [][32]byte{pubkey0, pubkey1, pubkey2},
		votingPowers: []uint64{100, 100, 100},
	}

	currentSet := updateclientContract.IICS07TendermintMsgsValidatorSet{
		Validators: []updateclientContract.IICS07TendermintMsgsValidatorInfo{
			{PubKey: pubkey0, VotingPower: 100},
			{PubKey: pubkey2, VotingPower: 125},
			{PubKey: pubkey1, VotingPower: 100},
		},
	}

	delta, ok, reason := detectValidatorSetDelta(baseHash, baseSet, currentSet)
	if !ok {
		t.Fatalf("expected multi-leaf delta, got reason=%q", reason)
	}
	if delta.LeafCount != 2 {
		t.Fatalf("leaf count: got %d want 2", delta.LeafCount)
	}
	if delta.Indices[0] != 1 || delta.PubKeys[0] != pubkey2 || delta.VotingPowers[0] != 125 {
		t.Fatalf("first changed leaf decoded incorrectly: %+v", delta)
	}
	if delta.Indices[1] != 2 || delta.PubKeys[1] != pubkey1 || delta.VotingPowers[1] != 100 {
		t.Fatalf("second changed leaf decoded incorrectly: %+v", delta)
	}
}

func TestDetectValidatorSetDeltaRejectsUnsafeCases(t *testing.T) {
	baseHash := [32]byte{0x88, 0xbe}
	pubkey0 := [32]byte{0x01}
	pubkey1 := [32]byte{0x02}
	pubkey2 := [32]byte{0x03}
	baseSet := cachedCosmosValidatorSet{
		indices:      []uint32{0, 1, 2},
		pubkeys:      [][32]byte{pubkey0, pubkey1, pubkey2},
		votingPowers: []uint64{100, 100, 100},
	}

	tests := []struct {
		name string
		set  updateclientContract.IICS07TendermintMsgsValidatorSet
	}{
		{
			name: "no change",
			set: updateclientContract.IICS07TendermintMsgsValidatorSet{
				Validators: []updateclientContract.IICS07TendermintMsgsValidatorInfo{
					{PubKey: pubkey0, VotingPower: 100},
					{PubKey: pubkey1, VotingPower: 100},
					{PubKey: pubkey2, VotingPower: 100},
				},
			},
		},
		{
			name: "validator count changed",
			set: updateclientContract.IICS07TendermintMsgsValidatorSet{
				Validators: []updateclientContract.IICS07TendermintMsgsValidatorInfo{
					{PubKey: pubkey0, VotingPower: 100},
					{PubKey: pubkey1, VotingPower: 125},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok, _ := detectValidatorSetDelta(baseHash, baseSet, tc.set); ok {
				t.Fatal("did not expect delta")
			}
		})
	}
}

func TestDetectValidatorSetDeltaRejectsTooManyLeaves(t *testing.T) {
	baseHash := [32]byte{0x88, 0xbe}
	baseSet := cachedCosmosValidatorSet{
		indices:      make([]uint32, 17),
		pubkeys:      make([][32]byte, 17),
		votingPowers: make([]uint64, 17),
	}
	currentSet := updateclientContract.IICS07TendermintMsgsValidatorSet{
		Validators: make([]updateclientContract.IICS07TendermintMsgsValidatorInfo, 17),
	}
	for i := 0; i < 17; i++ {
		pubkey := [32]byte{byte(i + 1)}
		baseSet.indices[i] = uint32(i)
		baseSet.pubkeys[i] = pubkey
		baseSet.votingPowers[i] = 100
		currentSet.Validators[i] = updateclientContract.IICS07TendermintMsgsValidatorInfo{
			PubKey:      pubkey,
			VotingPower: uint64(101 + i),
		}
	}

	if _, ok, _ := detectValidatorSetDelta(baseHash, baseSet, currentSet); ok {
		t.Fatal("did not expect delta over max leaf count")
	}
}
