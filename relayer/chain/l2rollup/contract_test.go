package l2rollup

import (
	"testing"

	"relayer/chain"
	"relayer/chain/chaintest"
)

// The adapter contract suite. Everything asserted here is asserted identically
// on the other two adapters: adding a fourth chain is a package plus this file.
func TestAdapterContract(t *testing.T) {
	chaintest.RunSourceContract(t, chaintest.SourceFixture{Name: "l2rollup", Unwired: &Source{chainType: chain.OPStack}})
	chaintest.RunDestinationContract(t, chaintest.DestinationFixture{Name: "l2rollup", Unwired: &Destination{}})
}
