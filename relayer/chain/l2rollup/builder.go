package l2rollup

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"relayer/chain"
)

// headerBuildTimeout bounds one full proof assembly (a few eth_getProof /
// eth_getBlockByNumber / optimism_outputAtBlock round trips, plus the on-demand
// Ethereum client update the OP builder performs first). Generous enough that a
// slow-but-live node still finishes, short enough that a dead one is retried rather
// than waited on forever.
const headerBuildTimeout = 90 * time.Second

// HeaderBuilder gathers the chain-specific trustless proof for one L2 block into a
// per-L2 ClientMessage header (OpStackHeader vs ArbitrumHeader): the beacon_slot +
// l1_state_root to verify against the shared ETH client, the L1 rollup-contract MPT
// witnesses (RollupCore assertion for Arbitrum, DisputeGameFactory game for
// OP-Stack), the canonical L2 header, and the L2 router account proof. It is the
// only per-L2 part of the client-update path.
type HeaderBuilder interface {
	// Name is the registry `builder` name (e.g. "l2-opstack", "l2-arbitrum").
	Name() string
	// BuildHeader builds the client message header for request and returns
	// the L2 block the header actually COMMITS. That committed height can be lower
	// than request.Height (e.g. the attested game/assertion at or below the
	// request), so the caller must advance the client state by the committed height,
	// not the requested one — otherwise proofs are built at an unproven height and
	// recvs fail. It reads the L1 + L2 chains; a transient RPC / not-yet-available
	// error is fine (the module re-queues via chain.Retryable in Build).
	BuildHeader(ctx context.Context, request HeaderRequest) (msg ClientMessage, committedHeight uint64, err error)
}

// HeaderRequest is the internal Source-to-Builder contract. Carrying Finality
// with Height prevents a Safe or Finalized source selection from being silently
// rebuilt using provisional evidence.
type HeaderRequest struct {
	Height   uint64
	Finality HeadKind
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
	request, err := decodeHeaderRequest(header)
	if err != nil {
		return chain.ClientUpdate{}, chain.Retryable(err)
	}
	// Bound the whole proof assembly. A header build is a handful of L1/L2 JSON-RPC
	// calls, and go-ethereum's HTTP client has no timeout of its own — a node that
	// accepts the connection and then never answers blocks forever. The relay module
	// drives one direction on a single goroutine, so that hangs the direction with no
	// error, no retry, and nothing in the log. The relay ctx alone does not save us:
	// it is only cancelled at shutdown. Timing out turns a wedged direction into a
	// retryable error, which is what the caller already knows how to handle.
	ctx, cancel := context.WithTimeout(ctx, headerBuildTimeout)
	defer cancel()

	msg, committedHeight, err := b.headerBuilder.BuildHeader(ctx, request)
	if err != nil {
		return chain.ClientUpdate{}, chain.Retryable(fmt.Errorf("l2: assemble %s proof at height %d: %w", request.Finality, request.Height, err))
	}
	payload, err := msg.EncodeClientMessage()
	if err != nil {
		return chain.ClientUpdate{}, err
	}
	// Advance by the height the header commits (may be < requested), so the module's
	// lastHeight never overstates what the client can actually prove.
	return chain.ClientUpdate{Height: committedHeight, Payloads: [][]byte{payload}}, nil
}

// decodeHeaderRequest reads the finality byte and big-endian L2 height emitted by
// Source.QueryHeader.
func decodeHeaderRequest(header []byte) (HeaderRequest, error) {
	if len(header) != 9 {
		return HeaderRequest{}, fmt.Errorf("l2: malformed header request (want 9 bytes, got %d)", len(header))
	}
	request := HeaderRequest{
		Finality: HeadKind(header[0]),
		Height:   binary.BigEndian.Uint64(header[1:]),
	}
	if err := request.Finality.validate(); err != nil {
		return HeaderRequest{}, err
	}
	return request, nil
}
