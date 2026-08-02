package l2rollup

import (
	"context"

	attestorpb "attestor/types/attestor"
)

// AttestorClient is the subset of the L2 attestor sidecar's gRPC surface the relayer
// gates on (AttestorService, #240/#258). The 2-method interface is defined HERE, on the
// consumer side, so the source/builders stay unit-testable with a small fake instead of
// a full gRPC server. It returns the generated proto AttestedRoot directly — the
// attestor module is the single source of truth for that shape (attestor/types), so no
// parallel relayer-side struct is kept.
//
// The attestor independently RE-DERIVES L2 state from L1 (OP: a verify-mode op-node
// replica compares each DisputeGameFactory root claim against its own L1-derived
// output root; Arbitrum: a Nitro replay). So an attested root is a trust statement,
// not a relay of the on-chain proposal — which is exactly why RelayableHeight gates
// on it rather than on the raw L2 RPC head.
//
// The relayer reads only the height/provenance fields: L2BlockNumber, Source, and
// the provenance oneof (GameIndex for OP games or AssertionHash for Arbitrum).
// It does NOT consume the Root bytes — the builders re-derive the header from
// L1/L2, so the root itself is never packaged, only used as an attestation gate.
type AttestorClient interface {
	// AttestedUpTo returns the highest L2 block the attestor has independently
	// confirmed for srcChain. includeProvisional accepts Safe-but-not-yet-finalized
	// roots (false = finalized only). found is false before the first attestation.
	AttestedUpTo(ctx context.Context, srcChain string, includeProvisional bool) (*attestorpb.AttestedRoot, bool, error)
	// AttestedRootAtOrBelow returns the best attested root at or below l2BlockNumber —
	// the canonical commitment for a target height (e.g. the OP game_index the header
	// builder proves). found is false when nothing qualifies at or below the bound.
	AttestedRootAtOrBelow(ctx context.Context, srcChain string, l2BlockNumber uint64, includeProvisional bool) (*attestorpb.AttestedRoot, bool, error)
}
