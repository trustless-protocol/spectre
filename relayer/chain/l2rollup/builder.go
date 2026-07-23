package l2rollup

import (
	"context"
	"encoding/binary"
	"fmt"

	"relayer/chain"
)

// HeaderBuilder gathers the chain-specific trustless proof for one L2 block into
// a ClientMessageData: the l1_beacon_slot to verify against the shared ETH client,
// the L1 rollup MPT witnesses (RollupCore for Arbitrum, AnchorStateRegistry for
// OP-Stack), the RLP-encoded L2 header, and the L2 IBC-handler account proof. It
// is the only per-L2 (OP-Stack vs Arbitrum) part of the client-update path.
type HeaderBuilder interface {
	// Name is the registry `builder` name (e.g. "l2-opstack", "l2-arbitrum").
	Name() string
	// BuildHeader builds the client message for L2 block l2Height. It reads the L1 +
	// L2 chains; a transient RPC / not-yet-available error is fine (the module
	// re-queues via chain.Retryable in Build).
	BuildHeader(ctx context.Context, l2Height uint64) (*ClientMessageData, error)
}

// Builder is the chain.ClientUpdateBuilder for the L2->Cosmos path. It decodes the
// target L2 height (from the L2 Source.QueryHeader), delegates the trustless-proof
// assembly to a per-L2 HeaderBuilder, and packages the JSON as the wasm
// ClientMessage data the Cosmos Destination submits — no relayer signature.
type Builder struct {
	headerBuilder HeaderBuilder
}

// NewBuilder wires the generic builder to a per-L2 HeaderBuilder.
func NewBuilder(headerBuilder HeaderBuilder) *Builder { return &Builder{headerBuilder: headerBuilder} }

func (b *Builder) Name() string { return b.headerBuilder.Name() }

// Build assembles the trustless client message for the L2 height in header. A
// proof-assembly failure (L1/L2 RPC down, height not yet available on L1) is
// transient, so it is wrapped chain.Retryable and the module re-queues.
func (b *Builder) Build(ctx context.Context, header []byte) (chain.ClientUpdate, error) {
	height, err := decodeHeight(header)
	if err != nil {
		return chain.ClientUpdate{}, chain.Retryable(err)
	}
	data, err := b.headerBuilder.BuildHeader(ctx, height)
	if err != nil {
		return chain.ClientUpdate{}, chain.Retryable(fmt.Errorf("l2: assemble proof at height %d: %w", height, err))
	}
	payload, err := data.Encode()
	if err != nil {
		return chain.ClientUpdate{}, err
	}
	return chain.ClientUpdate{Height: height, Payload: payload}, nil
}

// decodeHeight reads the 8-byte big-endian L2 height the Source.QueryHeader emits.
func decodeHeight(header []byte) (uint64, error) {
	if len(header) != 8 {
		return 0, fmt.Errorf("l2: malformed header (want 8 bytes, got %d)", len(header))
	}
	return binary.BigEndian.Uint64(header), nil
}
