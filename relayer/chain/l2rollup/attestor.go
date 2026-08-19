package l2rollup

import (
	"context"
	"errors"

	attestorpb "attestor/types/attestor"
)

// ErrVerifyStateRootUnsupported reports that this attestor build does not serve
// VerifyStateRoot at all. Both in-tree attestors — Arbitrum and OP-Stack, which
// also backs Base — implement it today, so this fires only on version skew: a
// daemon built before VerifyStateRoot shipped never overrides the embedded
// UnimplementedAttestorServiceServer, and the call returns gRPC Unimplemented.
//
// It is a distinct error because the header builder must not read "cannot answer" as
// "answered no": treating it as a refusal would fail every header build against a
// version-skewed deployment, which is worse than the gap it was meant to close. See
// the builder for how it degrades, and note the degradation is only correct while it
// is this narrow — any other failure stays fatal.
var ErrVerifyStateRootUnsupported = errors.New("attestor does not implement VerifyStateRoot")

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
// AttestedUpTo/AttestedRootAtOrBelow are height gates: the relayer reads
// L2BlockNumber from them and nothing else. AttestedRoot.Root is deliberately not
// consumed there because it is chain-specific — an OP output root, but an Arbitrum
// L2 state root — and the chain-agnostic builder has no way to recompute the OP
// form without the settlement machinery the attestor-trusted design removed.
//
// VerifyStateRoot is how the built header is bound to the attestor instead. The
// relayer sends the block identity it is about to package and the attestor compares
// it against its own independently-synced replica, so the comparison stays on the
// attestor's side of the chain-specific detail.
type AttestorClient interface {
	// AttestedUpTo returns the highest L2 block the attestor has independently
	// confirmed for srcChain. includeProvisional accepts Safe-but-not-yet-finalized
	// roots (false = finalized only). found is false before the first attestation.
	AttestedUpTo(ctx context.Context, srcChain string, includeProvisional bool) (*attestorpb.AttestedRoot, bool, error)
	// AttestedRootAtOrBelow returns the best attested root at or below l2BlockNumber —
	// the canonical commitment for a target height. found is false when nothing
	// qualifies at or below the bound.
	AttestedRootAtOrBelow(ctx context.Context, srcChain string, l2BlockNumber uint64, includeProvisional bool) (*attestorpb.AttestedRoot, bool, error)
	// VerifyStateRoot asks the attestor whether stateRoot and blockHash are what its
	// replica has at l2BlockNumber under runMode. valid is false when the replica
	// disagrees; err is reserved for transport and request failures, so a false valid
	// is a real divergence and not a degraded answer.
	//
	// srcChain names the chain, exactly as for the oracle calls above: a daemon
	// serving several chains would otherwise verify against whichever replica it
	// guessed, and a wrong guess answers valid=false for a good header — which the
	// caller cannot tell apart from a real divergence.
	VerifyStateRoot(ctx context.Context, srcChain string, l2BlockNumber uint64, stateRoot, blockHash []byte, runMode attestorpb.RunMode) (bool, error)
}
