package cosmos

import (
	"context"
	"fmt"
	"time"

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
// It holds a *services.Worker and the endpoints/configuration needed by the legacy
// builder to read the pinned validator set and trusted height from the on-chain
// SpectreClient. It is therefore constructed by the wiring (which owns those),
// not by the cfg-only registry — hence NewGroth16Builder, no RegisterClientUpdateBuilder.
//
// Note: BuildCosmosClientUpdateMsg fetches its own light block, so Build
// ignores the header argument for the Cosmos path.
type Groth16Builder struct {
	worker            *services.Worker
	cosmos            services.CosmosEndpoint
	evm               services.EVMEndpoint
	fetchTimeout      time.Duration
	rotationThreshold string
	proofType         string
	trustLevel        string
}

// NewGroth16Builder wires the builder to the worker and the endpoints it reads.
func NewGroth16Builder(worker *services.Worker, cosmos services.CosmosEndpoint, evm services.EVMEndpoint, fetchTimeout time.Duration, rotationThreshold, proofType, trustLevel string) *Groth16Builder {
	return &Groth16Builder{worker: worker, cosmos: cosmos, evm: evm, fetchTimeout: fetchTimeout, rotationThreshold: rotationThreshold, proofType: proofType, trustLevel: trustLevel}
}

func (b *Groth16Builder) Name() string { return "groth16" }

// Build runs the existing Cosmos update pipeline and packages the result as a
// chain.ClientUpdate. A "no update needed" result still reports the current
// on-chain trusted Cosmos height with a nil payload, so the module learns the
// client's real height (seeding the provability guard) without submitting a tx.
func (b *Groth16Builder) Build(ctx context.Context, _ []byte) (chain.ClientUpdate, error) {
	readCtx, cancel := context.WithTimeout(ctx, b.fetchTimeout)
	defer cancel()
	trustedBlock, err := services.FetchOnChainTrustedHeightWithContext(readCtx, b.evm)
	if err != nil {
		return chain.ClientUpdate{}, fmt.Errorf("groth16: on-chain trusted height: %w", err)
	}
	if trustedBlock < 0 {
		return chain.ClientUpdate{}, fmt.Errorf("groth16: negative trusted height %d", trustedBlock)
	}
	result, err := b.worker.BuildCosmosClientUpdateMsgWithContext(ctx, b.cosmos, b.evm, b.fetchTimeout, b.rotationThreshold, b.proofType, trustedBlock, b.trustLevel, false, 0)
	if err != nil {
		return chain.ClientUpdate{}, fmt.Errorf("groth16: build cosmos update: %w", err)
	}
	if !result.HasMsg || result.LightBlock == nil {
		// Client already current — report the trusted height so the module can seed
		// its provability guard; nil payload means "no tx".
		update := chain.ClientUpdate{Height: uint64(trustedBlock)}
		if result.LightBlock != nil {
			update.TrustedAt = result.LightBlock.SignedHeader.Header.Time
		}
		return update, nil
	}
	if result.LightBlock.BlockHeight < 0 {
		return chain.ClientUpdate{}, fmt.Errorf("groth16: negative block height %d", result.LightBlock.BlockHeight)
	}

	payload, err := codec.EncodeCosmosUpdate(int(result.Kind), result.AppMsg, result.NewValSet)
	if err != nil {
		return chain.ClientUpdate{}, fmt.Errorf("groth16: encode update: %w", err)
	}
	return chain.ClientUpdate{
		Height: uint64(result.LightBlock.BlockHeight), Payloads: [][]byte{payload},
		TrustedAt: result.LightBlock.SignedHeader.Header.Time,
	}, nil
}
