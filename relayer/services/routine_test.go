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

func TestDetectSingleVotingPowerDelta(t *testing.T) {
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

	delta, ok, reason := detectSingleVotingPowerDelta(baseHash, baseSet, currentSet)
	if !ok {
		t.Fatalf("expected delta, got reason=%q", reason)
	}
	if delta.BaseValidatorsHash != baseHash {
		t.Fatalf("base hash: got %x want %x", delta.BaseValidatorsHash, baseHash)
	}
	if delta.ChangedIndex != 1 {
		t.Fatalf("changed index: got %d want 1", delta.ChangedIndex)
	}
	if delta.NewVotingPower != 125 {
		t.Fatalf("new voting power: got %d want 125", delta.NewVotingPower)
	}
}

func TestDetectSingleVotingPowerDeltaRejectsUnsafeCases(t *testing.T) {
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
			name: "multiple power changes",
			set: updateclientContract.IICS07TendermintMsgsValidatorSet{
				Validators: []updateclientContract.IICS07TendermintMsgsValidatorInfo{
					{PubKey: pubkey0, VotingPower: 101},
					{PubKey: pubkey1, VotingPower: 125},
					{PubKey: pubkey2, VotingPower: 100},
				},
			},
		},
		{
			name: "validator reorder",
			set: updateclientContract.IICS07TendermintMsgsValidatorSet{
				Validators: []updateclientContract.IICS07TendermintMsgsValidatorInfo{
					{PubKey: pubkey1, VotingPower: 125},
					{PubKey: pubkey0, VotingPower: 100},
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
			if _, ok, _ := detectSingleVotingPowerDelta(baseHash, baseSet, tc.set); ok {
				t.Fatal("did not expect delta")
			}
		})
	}
}
