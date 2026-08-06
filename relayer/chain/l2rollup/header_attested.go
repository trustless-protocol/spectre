package l2rollup

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"

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
// Everything the client checks is derivable from these two fields, and everything
// it does NOT check — who attested this, at which head, from which L1 origin — is
// deliberately absent from the wire format rather than carried unverified.
//
// That makes the attestor the whole trust boundary, so it has to bound both axes.
// Source.RelayableHeight bounds HOW FAR we may relay; the VerifyStateRoot call below
// bounds WHAT is relayed at that height. Without the second, an L2 RPC that reorged
// past the attestor's frontier — or simply points at a different chain — hands back a
// replacement block at an approved height and the client accepts it, because the
// client only checks the header against itself and the router proof against the
// header. The attestor never certified that block.
//
// This binds against accidental divergence, not against a hostile relayer: nothing
// in the wire format lets the client re-check the answer, so a relayer that skips the
// call still gets its header accepted. Closing that needs the signed attestations the
// redesign defers to the next wire version.
type attestedHeaderBuilder struct {
	l2       *ethclient.Client  // L2 exec: l2_header + router eth_getProof
	router   ethcommon.Address  // the L2 ICS26Router (rollup_profile.common.l2_router)
	attestor AttestorClient     // nil in skip-finality mode, where nothing attests
	runMode  attestorpb.RunMode // replica head the attestor must answer against
	name     string             // registry builder name, for logs and errors

	// warnUnsupported keeps the "attestor cannot answer" warning to once per process
	// rather than once per header build.
	warnUnsupported sync.Once
}

// NewAttestedHeaderBuilder wires the builder to one L2 exec endpoint and router.
// attestor may be nil, which disables the binding check — that is skip-finality mode,
// where Source.RelayableHeight also degrades to the raw L2 head and there is no
// attestation to bind to in the first place.
func NewAttestedHeaderBuilder(l2 *ethclient.Client, router ethcommon.Address, attestor AttestorClient, runMode attestorpb.RunMode, name string) HeaderBuilder {
	return &attestedHeaderBuilder{l2: l2, router: router, attestor: attestor, runMode: runMode, name: name}
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
	if err := a.bindToAttestation(ctx, request.Height, l2Header); err != nil {
		return nil, 0, err
	}
	routerProof, err := relayerclient.EthGetProof(ctx, a.l2, a.router, nil, height)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: router proof at L2 height %d: %w", a.name, request.Height, err)
	}

	return &AttestedL2Header{
		L2Header:    toCanonicalHeader(l2Header),
		RouterProof: EvmAccountProof{Proof: routerProof.AccountProof},
	}, request.Height, nil
}

// bindToAttestation refuses a header the attestor's replica does not recognise as the
// canonical block at that height. state root AND block hash are both sent: the state
// root is what the client will trust, the block hash is what makes the answer specific
// to one block rather than to any block sharing a state root.
//
// A mismatch is transient by nature — the usual cause is the attestor replica lagging
// the L2 RPC by a block, or a reorg the attestor has not yet re-derived — so it is
// returned as a plain error and the relay loop retries. It never advances anything.
func (a *attestedHeaderBuilder) bindToAttestation(ctx context.Context, height uint64, l2Header *types.Header) error {
	if a.attestor == nil {
		return nil
	}
	stateRoot := l2Header.Root
	blockHash := l2Header.Hash()
	valid, err := a.attestor.VerifyStateRoot(ctx, height, stateRoot.Bytes(), blockHash.Bytes(), a.runMode)
	if errors.Is(err, ErrVerifyStateRootUnsupported) {
		// Only the Arbitrum attestor serves this RPC; the OP-Stack one does not, and
		// failing here would stop OP and Base relaying entirely. Degrade to the
		// pre-binding behaviour and say so loudly, once per process, so an operator
		// can see which chains are running unbound instead of assuming otherwise.
		a.warnUnsupported.Do(func() {
			log.Printf("[%s] WARNING: this attestor does not implement VerifyStateRoot, so built headers are NOT "+
				"bound to attested state — the relayer trusts its L2 RPC for block identity at approved heights. "+
				"Relaying continues; implement VerifyStateRoot in the attestor to close this.", a.name)
		})
		return nil
	}
	if err != nil {
		return fmt.Errorf("%s: verify L2 block %d against the attestor: %w", a.name, height, err)
	}
	if !valid {
		return fmt.Errorf(
			"%s: the attestor does not recognise L2 block %d (state_root=%s block_hash=%s) as canonical at run_mode=%s; "+
				"the L2 RPC and the attestor replica disagree, so this header is not attested state",
			a.name, height, stateRoot.Hex(), blockHash.Hex(), a.runMode)
	}
	return nil
}
