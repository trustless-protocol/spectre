package l2rollup

import (
	"context"
	"fmt"
	"math/big"

	attestorpb "attestor/types/attestor"
	relayerclient "relayer/client"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// attestedHeaderBuilder assembles the only header shape the attestor-trusted L2
// clients accept: the canonical L2 execution header at the target height, plus an
// account proof of that chain's ICS26Router against the state root that header
// declares.
//
// It is chain-agnostic on purpose. The per-chain builders it replaces existed to
// gather settlement evidence — DisputeGameFactory games for OP-Stack, RollupCore
// assertions for Arbitrum, each pinned to the L1 block the shared Ethereum client
// trusted. None of that is verified any more: the three verifier crates are now
// identical adapters differing only in `profile_version`, so there is nothing left
// to branch on.
//
// The attestor is the L2 trust boundary. Source.RelayableHeight bounds HOW FAR
// the relayer may go, and VerifyStateRoot binds WHAT it packages. On a match,
// the attestor returns an Ed25519 signature; this builder carries it in the
// header and the wasm client verifies it against its immutable profile key.
// A relayer cannot bypass that verification by skipping this RPC.
type attestedHeaderBuilder struct {
	l2       *ethclient.Client   // L2 exec: l2_header + router eth_getProof
	router   ethcommon.Address   // the L2 ICS26Router (rollup_profile.common.l2_router)
	attestor AttestorClient      // nil in skip-finality mode, where nothing attests
	verifier AttestationVerifier // configured attestor key, pinned again by the wasm client
	srcChain string              // attestor src_chain key, same as the source's
	runMode  attestorpb.RunMode  // replica head the attestor must answer against
	name     string              // registry builder name, for logs and errors

}

// NewAttestedHeaderBuilder wires the builder to one L2 exec endpoint and router.
// attestor may be nil, which disables the binding check — that is skip-finality mode,
// where Source.RelayableHeight also degrades to the raw L2 head and there is no
// attestation to bind to in the first place.
func NewAttestedHeaderBuilder(l2 *ethclient.Client, router ethcommon.Address, attestor AttestorClient, verifier AttestationVerifier, srcChain string, runMode attestorpb.RunMode, name string) HeaderBuilder {
	return &attestedHeaderBuilder{l2: l2, router: router, attestor: attestor, verifier: verifier, srcChain: srcChain, runMode: runMode, name: name}
}

func (a *attestedHeaderBuilder) Name() string { return a.name }

// BuildHeader returns the header for exactly request.Height. Unlike the settlement
// builders, there is no step-down: those had to commit whatever L2 block the newest
// proven game or assertion covered, which could be far below the request. Here the
// header IS the canonical block at the requested height, so the committed height
// always equals the requested one and the caller advances by it directly.
func (a *attestedHeaderBuilder) BuildHeader(ctx context.Context, request HeaderRequest) (ClientMessage, uint64, error) {
	height := new(big.Int).SetUint64(request.Height)

	l2Header, err := a.l2.HeaderByNumber(ctx, height)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: L2 header at %d: %w", a.name, request.Height, err)
	}
	signature, err := a.bindToAttestation(ctx, request.Height, l2Header)
	if err != nil {
		return nil, 0, err
	}
	routerProof, err := relayerclient.EthGetProof(ctx, a.l2, a.router, nil, height)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: router proof at L2 height %d: %w", a.name, request.Height, err)
	}

	return &AttestedL2Header{
		L2Header:          toCanonicalHeader(l2Header),
		RouterProof:       EvmAccountProof{Proof: routerProof.AccountProof},
		AttestorSignature: byteList(signature),
	}, request.Height, nil
}

// bindToAttestation refuses a header the attestor's replica does not recognise as the
// canonical block at that height and returns its signature. State root AND block hash
// are both signed: the root is what the client will trust, the block hash binds the
// complete execution header.
//
// Every failure here goes through classifyAttestorFailure, which is the single
// place that decides the class AND raises the alarm. A mismatch is transient by
// nature — the usual cause is the attestor replica lagging the L2 RPC by a block,
// or a reorg it has not yet re-derived — but it is alarmed, because a real
// divergence looks identical and retrying one in silence hides the one condition
// the attestor exists to detect. Nothing is ever advanced on any of these paths.
func (a *attestedHeaderBuilder) bindToAttestation(ctx context.Context, height uint64, l2Header *types.Header) ([]byte, error) {
	if a.attestor == nil {
		return nil, fmt.Errorf("%s: attestor is required to build an authenticated L2 header", a.name)
	}
	stateRoot := l2Header.Root
	blockHash := l2Header.Hash()
	attestation, err := a.attestor.VerifyStateRoot(ctx, a.srcChain, height, stateRoot.Bytes(), blockHash.Bytes(), a.runMode)
	if err != nil {
		return nil, classifyAttestorFailure(a.name,
			fmt.Errorf("%s: verify L2 block %d against the attestor: %w", a.name, height, err))
	}
	if !attestation.Valid {
		return nil, classifyAttestorFailure(a.name, fmt.Errorf(
			"%w: %s: the attestor does not recognise L2 block %d (state_root=%s block_hash=%s) as canonical at run_mode=%s; "+
				"the L2 RPC and the attestor replica disagree, so this header is not attested state",
			ErrAttestorDivergence, a.name, height, stateRoot.Hex(), blockHash.Hex(), a.runMode))
	}
	if err := a.verifier.Verify(height, stateRoot.Bytes(), blockHash.Bytes(), attestation.Signature); err != nil {
		// Through the table, not raw: an unclassified error defaults to transient,
		// and a signature that does not verify is the one attestor failure where
		// retrying unchanged is certainly useless.
		return nil, classifyAttestorFailure(a.name, fmt.Errorf("%w: L2 block %d: %v",
			ErrAttestorSignature, height, err))
	}
	return attestation.Signature, nil
}
