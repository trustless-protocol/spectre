package evm

import (
	"testing"

	"relayer/chain/chaintest"
)

// The adapter contract suite. Everything asserted here is asserted identically
// on the other two adapters: adding a fourth chain is a package plus this file.
func TestAdapterContract(t *testing.T) {
	chaintest.RunSourceContract(t, chaintest.SourceFixture{Name: "evm", Unwired: &Source{}})
	chaintest.RunDestinationContract(t, chaintest.DestinationFixture{Name: "evm", Unwired: &Destination{}})
}
