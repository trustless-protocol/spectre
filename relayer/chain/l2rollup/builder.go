package l2rollup

import (
	"context"
	"encoding/binary"
	"fmt"

	"relayer/chain"
)

// HeaderBuilder gathers the chain-specific trustless proof for one L2 block into a
// per-L2 ClientMessage header (OpStackHeader vs ArbitrumHeader): the beacon_slot +
// l1_state_root to verify against the shared ETH client, the L1 rollup-contract MPT
// witnesses (RollupCore assertion for Arbitrum, DisputeGameFactory game for
// OP-Stack), the canonical L2 header, and the L2 router account proof. It is the
// only per-L2 part of the client-update path.
type HeaderBuilder interface {
	// Name is the registry `builder` name (e.g. "l2-opstack", "l2-arbitrum").
	Name() string
	// BuildHeader builds the client message header for L2 block l2Height and returns
	// the L2 block the header actually COMMITS. That committed height can be lower
	// than the requested l2Height (e.g. the attested game/assertion at or below the
	// request), so the caller must advance the client state by the committed height,
	// not the requested one — otherwise proofs are built at an unproven height and
	// recvs fail. It reads the L1 + L2 chains; a transient RPC / not-yet-available
	// error is fine (the module re-queues via chain.Retryable in Build).
	BuildHeader(ctx context.Context, l2Height uint64) (msg ClientMessage, committedHeight uint64, err error)
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
// transient, so it is wrapped chain.Retryable and the module re-queues; a marshal
// failure of a fully-assembled header is a programming error, not transient.
func (b *Builder) Build(ctx context.Context, header []byte) (chain.ClientUpdate, error) {
	height, err := decodeHeight(header)
	if err != nil {
		return chain.ClientUpdate{}, chain.Retryable(err)
	}
	msg, committedHeight, err := b.headerBuilder.BuildHeader(ctx, height)
	if err != nil {
		return chain.ClientUpdate{}, chain.Retryable(fmt.Errorf("l2: assemble proof at height %d: %w", height, err))
	}
	payload, err := msg.EncodeClientMessage()
	if err != nil {
		return chain.ClientUpdate{}, err
	}
	// Advance by the height the header commits (may be < requested), so the module's
	// lastHeight never overstates what the client can actually prove.
	return chain.ClientUpdate{Height: committedHeight, Payload: payload}, nil
}

// decodeHeight reads the 8-byte big-endian L2 height the Source.QueryHeader emits.
func decodeHeight(header []byte) (uint64, error) {
	if len(header) != 8 {
		return 0, fmt.Errorf("l2: malformed header (want 8 bytes, got %d)", len(header))
	}
	return binary.BigEndian.Uint64(header), nil
}
