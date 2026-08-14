package client

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	cmted25519 "github.com/cometbft/cometbft/crypto/ed25519"
	commettypes "github.com/cometbft/cometbft/types"
)

func TestSpectreClientGenesisPinsNextValidatorSet(t *testing.T) {
	currentValSet := genesisTestValSet("current", 10, 20, 30)
	nextValSet := genesisTestValSet("next", 10, 20, 30)
	if bytes.Equal(currentValSet.Hash(), nextValSet.Hash()) {
		t.Fatal("test setup requires different current and next validator-set hashes")
	}

	trustedLightBlock := &LightBlock{
		SignedHeader: commettypes.SignedHeader{
			Header: &commettypes.Header{
				ChainID:            "test-chain-1",
				Height:             123,
				Time:               time.Unix(1_700_000_000, 0),
				ValidatorsHash:     currentValSet.Hash(),
				NextValidatorsHash: nextValSet.Hash(),
				AppHash:            []byte("app-hash"),
			},
		},
		ValSet:      currentValSet,
		NextValSet:  nextValSet,
		BlockHeight: 123,
	}

	genesis, err := spectreClientGenesisFromLightBlock(trustedLightBlock, 1_814_400, 0, "2/3", 0)
	if err != nil {
		t.Fatalf("unexpected genesis error: %v", err)
	}

	currentPinned, err := ValidatorSetToContract(currentValSet, "current")
	if err != nil {
		t.Fatalf("current ValidatorSetToContract: %v", err)
	}
	nextPinned, err := ValidatorSetToContract(nextValSet, "next")
	if err != nil {
		t.Fatalf("next ValidatorSetToContract: %v", err)
	}

	if got, want := genesis.TrustedConsensusState.NextValidatorsHash, bytesToBytes32(nextValSet.Hash()); got != want {
		t.Fatalf("next validators hash: got %X, want %X", got, want)
	}
	if !reflect.DeepEqual(genesis.InitialPinnedValidatorSet, nextPinned) {
		t.Fatal("initial pinned validator set should come from the trusted light block's next validator set")
	}
	if reflect.DeepEqual(genesis.InitialPinnedValidatorSet, currentPinned) {
		t.Fatal("initial pinned validator set must not come from the current validator set")
	}
}

func genesisTestValSet(name string, powers ...int64) commettypes.ValidatorSet {
	validators := make([]*commettypes.Validator, len(powers))
	for i, power := range powers {
		secret := []byte(fmt.Sprintf("%s-%d", name, i))
		validators[i] = commettypes.NewValidator(cmted25519.GenPrivKeyFromSecret(secret).PubKey(), power)
	}
	return *commettypes.NewValidatorSet(validators)
}
