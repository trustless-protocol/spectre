package l2rollup

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

// The two per-L2 HeaderBuilders. Both read L1 (for the rollup contract proof +
// the l1_beacon_slot to verify against the shared ETH client) and L2 (for the RLP
// header + IBC-handler account proof); they differ only in the rollup-specific
// L1 witnesses, per Dũng's schema. Skeletons: the eth_getProof / RLP / slot-select
// mechanics are TODO, pending the L2 config (contract addresses, IBC slot) and the
// shared-ETH-client handle for l1_beacon_slot.

// opStackHeaderBuilder builds the OP-Stack (Optimism, Base) client message from the
// AnchorStateRegistry output-root + L2-height storage witnesses + output-root
// preimage + L2 header + L2 IBC account proof.
type opStackHeaderBuilder struct {
	l1 *ethclient.Client // AnchorStateRegistry storage proof + l1_beacon_slot
	l2 *ethclient.Client // L2 header + IBC-handler account proof
	// TODO(Đức): AnchorStateRegistry address, L2 IBC-handler address + storage slot,
	// and the shared cw-ics08-wasm-eth client handle (to pick l1_beacon_slot).
}

// NewOPStackHeaderBuilder wires the OP-Stack headerBuilder to the L1 + L2 RPC endpoints.
func NewOPStackHeaderBuilder(l1, l2 *ethclient.Client) *opStackHeaderBuilder {
	return &opStackHeaderBuilder{l1: l1, l2: l2}
}

func (a *opStackHeaderBuilder) Name() string { return "l2-opstack" }

func (a *opStackHeaderBuilder) BuildHeader(_ context.Context, _ uint64) (*ClientMessageData, error) {
	// TODO(Đức):
	//  1. l1_beacon_slot = the finalized ETH slot the shared client already covers.
	//  2. eth_getProof AnchorStateRegistry on L1 (account + storage) at that L1 block.
	//  3. fetch the output-root preimage (OptimismInputs.OutputRootPreimage).
	//  4. RLP-encode the L2 header at l2Height.
	//  5. eth_getProof the L2 IBC-handler account.
	return nil, fmt.Errorf("l2-opstack: BuildHeader not yet wired (needs AnchorStateRegistry + IBC-handler config)")
}

// arbitrumHeaderBuilder builds the Arbitrum (Nitro) client message from the assertion
// inputs (parent_assertion_hash, assertion state, inbox_acc) + RollupCore storage
// witness + L2 header + L2 IBC account proof.
type arbitrumHeaderBuilder struct {
	l1 *ethclient.Client // RollupCore storage proof + l1_beacon_slot
	l2 *ethclient.Client // L2 header + IBC-handler account proof
	// TODO(Đức): RollupCore address, L2 IBC-handler address + storage slot, and the
	// shared cw-ics08-wasm-eth client handle (to pick l1_beacon_slot).
}

// NewArbitrumHeaderBuilder wires the Arbitrum headerBuilder to the L1 + L2 RPC endpoints.
func NewArbitrumHeaderBuilder(l1, l2 *ethclient.Client) *arbitrumHeaderBuilder {
	return &arbitrumHeaderBuilder{l1: l1, l2: l2}
}

func (a *arbitrumHeaderBuilder) Name() string { return "l2-arbitrum" }

func (a *arbitrumHeaderBuilder) BuildHeader(_ context.Context, _ uint64) (*ClientMessageData, error) {
	// TODO(Đức):
	//  1. l1_beacon_slot = the finalized ETH slot the shared client already covers.
	//  2. gather the assertion inputs (parent_assertion_hash, assertion state, inbox_acc).
	//  3. eth_getProof RollupCore on L1 (account + storage) at that L1 block.
	//  4. RLP-encode the L2 header at l2Height.
	//  5. eth_getProof the L2 IBC-handler account.
	return nil, fmt.Errorf("l2-arbitrum: BuildHeader not yet wired (needs RollupCore + IBC-handler config)")
}

// Both headerBuilders satisfy HeaderBuilder.
var (
	_ HeaderBuilder = (*opStackHeaderBuilder)(nil)
	_ HeaderBuilder = (*arbitrumHeaderBuilder)(nil)
)
