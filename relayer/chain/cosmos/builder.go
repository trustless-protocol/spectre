package cosmos

import (
	"context"
	"fmt"

	"relayer/chain"
	"relayer/chain/codec"
	"relayer/services"
)

// Groth16Builder is the "groth16" chain.ClientUpdateBuilder for a Cosmos source:
// it produces a ZK Tendermint update that the destination's SpectreClient
// verifies. It WRAPS the existing services.BuildCosmosClientUpdateMsg pipeline
// (extract signatures -> select pinned quorum -> GenerateProof -> assemble) with
// zero changes to the prover I/O.
//
// It holds a *services.Worker and services.Context because the legacy builder
// reads the pinned validator set + trusted height from the on-chain SpectreClient
// via the Context. It is therefore constructed by the wiring (which owns those),
// not by the cfg-only registry — hence NewGroth16Builder, no RegisterClientUpdateBuilder.
//
// Note: BuildCosmosClientUpdateMsg fetches its own light block, so Build
// ignores the header argument for the Cosmos path. Splitting fetch from build to
// honor the header is deferred to slice 4 (Context dissolution).
type Groth16Builder struct {
	worker     *services.Worker
	svcCtx     services.Context
	proofType  string
	trustLevel string
}

// NewGroth16Builder wires the builder to the shared worker + context.
func NewGroth16Builder(worker *services.Worker, svcCtx services.Context, proofType, trustLevel string) *Groth16Builder {
	return &Groth16Builder{worker: worker, svcCtx: svcCtx, proofType: proofType, trustLevel: trustLevel}
}

func (b *Groth16Builder) Name() string { return "groth16" }

// Build runs the existing Cosmos update pipeline and packages the result as a
// chain.ClientUpdate. A "no update needed" result still reports the current
// on-chain trusted Cosmos height with a nil payload, so the module learns the
// client's real height (seeding the provability guard) without submitting a tx.
func (b *Groth16Builder) Build(_ context.Context, _ []byte) (chain.ClientUpdate, error) {
	trustedBlock, err := services.FetchOnChainTrustedHeight(b.svcCtx)
	if err != nil {
		return chain.ClientUpdate{}, fmt.Errorf("groth16: on-chain trusted height: %w", err)
	}
	if trustedBlock < 0 {
		return chain.ClientUpdate{}, fmt.Errorf("groth16: negative trusted height %d", trustedBlock)
	}
	result, err := b.worker.BuildCosmosClientUpdateMsg(b.svcCtx, b.proofType, trustedBlock, b.trustLevel, false)
	if err != nil {
		return chain.ClientUpdate{}, fmt.Errorf("groth16: build cosmos update: %w", err)
	}
	if !result.HasMsg || result.LightBlock == nil {
		// Client already current — report the trusted height so the module can seed
		// its provability guard; nil payload means "no tx".
		return chain.ClientUpdate{Height: uint64(trustedBlock)}, nil
	}
	if result.LightBlock.BlockHeight < 0 {
		return chain.ClientUpdate{}, fmt.Errorf("groth16: negative block height %d", result.LightBlock.BlockHeight)
	}

	payload, err := codec.EncodeCosmosUpdate(int(result.Kind), result.AppMsg, result.NewValSet)
	if err != nil {
		return chain.ClientUpdate{}, fmt.Errorf("groth16: encode update: %w", err)
	}
	return chain.ClientUpdate{Height: uint64(result.LightBlock.BlockHeight), Payloads: [][]byte{payload}}, nil
}
