package l2rollup

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"relayer/chain"
)

// headerBuildTimeout bounds one header assembly (an eth_getBlockByNumber and an
// eth_getProof). Generous enough that a slow-but-live node still finishes, short
// enough that a dead one is retried rather than waited on forever.
const headerBuildTimeout = 90 * time.Second

// HeaderBuilder assembles one L2 block into a ClientMessage header. There is a single
// implementation now (attestedHeaderBuilder): the settlement proofs that used to make
// this per-chain are no longer verified, so nothing distinguishes OP-Stack from
// Arbitrum on this path. The interface stays because the signature slice will add a
// second implementation — one that also carries an attestor signature.
type HeaderBuilder interface {
	// Name is the registry `builder` name (e.g. "l2-opstack", "l2-arbitrum").
	Name() string
	// BuildHeader builds the client message header for request and returns the L2
	// block the header COMMITS. The caller advances client state by that height, not
	// the requested one. Today they are always equal — the header is the canonical
	// block at the requested height — but the contract is kept because a settlement
	// builder committed whatever block its newest proven game covered, which could be
	// far lower, and re-introducing one must not silently overstate the client height.
	// A transient RPC error is fine: the module re-queues via chain.Retryable in Build.
	BuildHeader(ctx context.Context, request HeaderRequest) (msg ClientMessage, committedHeight uint64, err error)
}

// HeaderRequest is the internal Source-to-Builder contract.
//
// It used to carry the head kind alongside the height, so a Safe or Finalized source
// selection could not be silently rebuilt from provisional evidence. The header no
// longer records a head kind at all, and the Source already applies the selection when
// it decides which attestor frontier to gate on, so passing it here would only let the
// builder re-derive a decision that has already been made.
type HeaderRequest struct {
	Height uint64
}

// Builder is the chain.ClientUpdateBuilder for the L2->Cosmos path. It decodes the
// target L2 height (from the L2 Source.QueryHeader), delegates header assembly to a
// HeaderBuilder, and packages the JSON as the wasm ClientMessage data the Cosmos
// Destination submits.
type Builder struct {
	headerBuilder HeaderBuilder
}

// NewBuilder wires the generic builder to a per-L2 HeaderBuilder.
func NewBuilder(headerBuilder HeaderBuilder) *Builder { return &Builder{headerBuilder: headerBuilder} }

func (b *Builder) Name() string { return b.headerBuilder.Name() }

// Build assembles the trustless client message for the L2 height in header. A
// proof-assembly failure (L1/L2 RPC down, height not yet available on L1) is
// transient and the module re-queues; a marshal failure of a fully-assembled
// header is a programming error, not transient.
//
// The header-build failure uses chain.Keep, not chain.Transient. Transient here
// re-labelled everything the header builder had already classified — an attestor
// answering "that source chain does not exist" arrived at the engine looking like
// an RPC blip, and no amount of reading the engine would show why.
func (b *Builder) Build(ctx context.Context, header []byte) (chain.ClientUpdate, error) {
	request, err := decodeHeaderRequest(header)
	if err != nil {
		return chain.ClientUpdate{}, chain.Transient(err)
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
		return chain.ClientUpdate{}, chain.Keep(fmt.Errorf("l2: assemble header at height %d: %w", request.Height, err))
	}
	payload, err := msg.EncodeClientMessage()
	if err != nil {
		return chain.ClientUpdate{}, err
	}
	// Advance by the height the header commits (may be < requested), so the module's
	// lastHeight never overstates what the client can actually prove.
	return chain.ClientUpdate{Height: committedHeight, Payloads: [][]byte{payload}}, nil
}

// decodeHeaderRequest reads the big-endian L2 height emitted by Source.QueryHeader.
func decodeHeaderRequest(header []byte) (HeaderRequest, error) {
	if len(header) != 8 {
		return HeaderRequest{}, fmt.Errorf("l2: malformed header request (want 8 bytes, got %d)", len(header))
	}
	return HeaderRequest{Height: binary.BigEndian.Uint64(header)}, nil
}
