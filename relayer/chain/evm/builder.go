package evm

import (
	"context"
	"fmt"

	"relayer/chain"
	"relayer/chain/codec"
	"relayer/services"
)

// BeaconBuilder is the "beacon" chain.ClientUpdateBuilder for an Ethereum L1
// source: it assembles the sync-committee + storage-proof update that the
// destination's 08-wasm beacon light client verifies (ETH L1 -> Cosmos). It does
// NO relayer-side proving or re-execution — the wasm client verifies BLS +
// Merkle itself. It WRAPS services.BuildEthClientUpdateMsgs unchanged.
//
// Holds *services.Worker + services.Context (the legacy builder reads the beacon
// API + the on-chain 08-wasm client state via the context), so it is constructed
// by the wiring, not the cfg-only registry.
type BeaconBuilder struct {
	worker *services.Worker
	svcCtx services.Context
}

// NewBeaconBuilder wires the builder to the shared worker + context.
func NewBeaconBuilder(worker *services.Worker, svcCtx services.Context) *BeaconBuilder {
	return &BeaconBuilder{worker: worker, svcCtx: svcCtx}
}

func (b *BeaconBuilder) Name() string { return "beacon" }

// Build runs the existing beacon update pipeline and packages the result as a
// chain.ClientUpdate.
//
// Height is the update's EXECUTION block number (not the beacon SigSlot): the
// RelayModule compares it against packet event heights, which are execution block
// numbers, so the append-only guard and the provability guard stay in one number
// space. The SigSlot the Cosmos-side proof needs is carried inside the payload
// (decoded by the destination), unaffected by this choice.
//
// A "no update needed" result still reports the current execution block with a nil
// payload, so the module learns the client's real height (seeding the provability
// guard) without submitting a tx. BuildEthClientUpdateMsgs returns the current
// EthClientState even on the no-op path, so this is always available.
func (b *BeaconBuilder) Build(_ context.Context, _ []byte) (chain.ClientUpdate, error) {
	result, err := b.worker.BuildEthClientUpdateMsgs(b.svcCtx)
	if err != nil {
		return chain.ClientUpdate{}, fmt.Errorf("beacon: build eth client update: %w", err)
	}
	if result.EthClientState == nil {
		return chain.ClientUpdate{}, fmt.Errorf("beacon: build returned nil client state")
	}
	execBlock := result.EthClientState.LatestExecutionBlockNumber
	if len(result.Msgs) == 0 {
		return chain.ClientUpdate{Height: execBlock}, nil // client already current — learn height, no tx
	}
	payload, err := codec.EncodeBeaconUpdate(result.Msgs, *result.EthClientState, result.ProofTimestamp, result.SigSlot)
	if err != nil {
		return chain.ClientUpdate{}, fmt.Errorf("beacon: encode update: %w", err)
	}
	return chain.ClientUpdate{Height: execBlock, Payload: payload}, nil
}
