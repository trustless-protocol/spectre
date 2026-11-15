package l2rollup

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strings"
	"time"

	"attestor/types/attestation"
	attestorpb "attestor/types/attestor"
	"relayer/chain"
	relayerclient "relayer/client"
	workergroup "relayer/workers"

	ethereum "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// stepDownFrontierTimeout bounds the frontier round of a step-down. The answers
// only order the attempts, so an endpoint slow to answer is left out rather than
// waited for. A var only so tests can shorten it.
var stepDownFrontierTimeout = 5 * time.Second

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
// The client verifies the indexed signatures before it touches the router proof.
// Source.RelayableHeight decides HOW FAR the relay may advance; this builder asks
// independent attestors to sign WHAT it packages at that exact height.
type attestedHeaderBuilder struct {
	l2        *ethclient.Client // L2 exec: l2_header + router eth_getProof
	router    ethcommon.Address // the L2 ICS26Router (rollup_profile.common.l2_router)
	chainID   uint64            // profile.common.l2_chain_id
	attestors []SigningAttestor // ordered by immutable wasm attestor index
	set       attestation.AttestorConfig
	setHash   [attestation.HashLength]byte
	srcChain  string             // attestor src_chain key, same as the source's
	runMode   attestorpb.RunMode // replica head the attestor must answer against
	// includeProvisional is the source's frontier policy, so a step-down asks for
	// the same frontiers the source gates on.
	includeProvisional bool
	name               string // registry builder name, for logs and errors
}

// SigningAttestor is one reachable attestor endpoint assigned to its immutable
// ClientState public-key index. It never carries a private key.
type SigningAttestor struct {
	Index  uint16
	Client AttestorClient
}

// NewAttestedHeaderBuilder wires the builder to all configured signing attestors.
// The endpoint list may contain more than the threshold to provide availability;
// each emitted update carries exactly threshold valid signatures in index order.
// includeProvisional must be the source's (NewSource): a step-down asks the
// endpoints for the frontiers the source gates on.
func NewAttestedHeaderBuilder(l2 *ethclient.Client, router ethcommon.Address, chainID uint64, set attestation.AttestorConfig, attestors []SigningAttestor, srcChain string, runMode attestorpb.RunMode, includeProvisional bool, name string) (HeaderBuilder, error) {
	if l2 == nil {
		return nil, fmt.Errorf("%s: L2 client must not be nil", name)
	}
	if chainID == 0 {
		return nil, fmt.Errorf("%s: L2 chain ID must not be zero", name)
	}
	if srcChain == "" {
		return nil, fmt.Errorf("%s: attestor src_chain must not be empty", name)
	}
	if err := set.Validate(); err != nil {
		return nil, fmt.Errorf("%s: invalid attestor set: %w", name, err)
	}
	if len(attestors) < int(set.Threshold) {
		return nil, fmt.Errorf("%s: %d attestor endpoints cannot satisfy threshold %d", name, len(attestors), set.Threshold)
	}
	ordered := append([]SigningAttestor(nil), attestors...)
	seen := make(map[uint16]struct{}, len(ordered))
	for _, endpoint := range ordered {
		if endpoint.Client == nil {
			return nil, fmt.Errorf("%s: attestor endpoint at index %d is nil", name, endpoint.Index)
		}
		if int(endpoint.Index) >= len(set.PublicKeys) {
			return nil, fmt.Errorf("%s: attestor endpoint index %d is outside the configured set", name, endpoint.Index)
		}
		if _, duplicate := seen[endpoint.Index]; duplicate {
			return nil, fmt.Errorf("%s: duplicate attestor endpoint index %d", name, endpoint.Index)
		}
		seen[endpoint.Index] = struct{}{}
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Index < ordered[j].Index })
	setHash, err := set.SetHash()
	if err != nil {
		return nil, fmt.Errorf("%s: derive attestor set hash: %w", name, err)
	}
	return &attestedHeaderBuilder{l2: l2, router: router, chainID: chainID, set: set, setHash: setHash, attestors: ordered, srcChain: srcChain, runMode: runMode, includeProvisional: includeProvisional, name: name}, nil
}

func (a *attestedHeaderBuilder) Name() string { return a.name }

// BuildHeader returns the header for request.Height or, when threshold attestors
// cannot sign that block, for the highest attestor frontier below it that they do
// sign. The header is the canonical block at the height it returns.
//
// The request comes from Source.RelayableHeight: the threshold-th highest of the
// endpoints' frontier answers, which are not signed. One endpoint that reports
// more than its replica holds (faulty, or rewritten on the plaintext gRPC path)
// lifts that answer to a block only threshold-1 honest replicas have reached, and
// the signatures fall short there whenever the honest frontiers differ (#575). A
// step-down costs one signature round per such endpoint; starting from the
// threshold-th highest still keeps one low endpoint from pinning the client.
//
// The returned height is then below the request. The caller advances by it, and
// packets above it wait for the next update (see HeaderBuilder).
func (a *attestedHeaderBuilder) BuildHeader(ctx context.Context, request HeaderRequest) (ClientMessage, uint64, error) {
	height := request.Height
	l2Header, signatures, err := a.signedHeaderAt(ctx, height)
	if err != nil && stepDownMayHelp(err) {
		height, l2Header, signatures, err = a.stepDown(ctx, height, err)
	}
	if err != nil {
		return nil, 0, err
	}
	routerProof, err := relayerclient.EthGetProof(ctx, a.l2, a.router, nil, new(big.Int).SetUint64(height))
	if err != nil {
		return nil, 0, fmt.Errorf("%s: router proof at L2 height %d: %w", a.name, height, err)
	}

	return &SignedAttestedL2Header{
		L2Header:          toCanonicalHeader(l2Header),
		RouterProof:       EvmAccountProof{Proof: routerProof.AccountProof},
		AttestorSignature: signatures,
	}, height, nil
}

// signedHeaderAt reads the canonical L2 header at height and binds it to
// threshold attestor signatures.
func (a *attestedHeaderBuilder) signedHeaderAt(ctx context.Context, height uint64) (*types.Header, []IndexedAttestorSignature, error) {
	l2Header, err := a.l2.HeaderByNumber(ctx, new(big.Int).SetUint64(height))
	if err != nil {
		return nil, nil, fmt.Errorf("%s: L2 header at %d: %w", a.name, height, err)
	}
	signatures, err := a.bindToAttestation(ctx, height, l2Header)
	if err != nil {
		return nil, nil, err
	}
	return l2Header, signatures, nil
}

// stepDownMayHelp reports whether a lower block can succeed where this one
// failed: the L2 RPC has no block there yet, or an attestor replica has not
// reached it. Any other failure, and a quorum no retry can complete, is the same
// at every height.
func stepDownMayHelp(err error) bool {
	if chain.IsPermanent(err) {
		return false
	}
	return errors.Is(err, ethereum.NotFound) || errors.Is(err, ErrAttestorReplicaBehind)
}

// stepDown builds at the next lower attestor frontier, then the next, until
// threshold attestors sign one. Each endpoint names at most one candidate, so it
// costs at most one signature round per endpoint. It returns cause when no
// candidate is signed.
func (a *attestedHeaderBuilder) stepDown(ctx context.Context, height uint64, cause error) (uint64, *types.Header, []IndexedAttestorSignature, error) {
	frontiers := a.endpointFrontiers(ctx)
	for _, candidate := range stepDownCandidates(frontiers, height) {
		l2Header, signatures, err := a.signedHeaderAt(ctx, candidate)
		if err == nil {
			log.Printf("WARN [%s] L2 block %d could not be built (%v); built at %d, the next lower attestor frontier. "+
				"Frontiers: %s. If this repeats, an attestor reports a frontier its replica cannot sign, "+
				"or this relayer's L2 RPC lags the attestors",
				a.name, height, cause, candidate, describeFrontiers(frontiers))
			return candidate, l2Header, signatures, nil
		}
		if !stepDownMayHelp(err) {
			return 0, nil, nil, err
		}
	}
	return 0, nil, nil, cause
}

// endpointFrontier is one endpoint's AttestedUpTo answer.
type endpointFrontier struct {
	index  uint16
	height uint64
}

// endpointFrontiers asks every endpoint for its attested frontier. One that
// fails, has attested nothing, or does not answer within stepDownFrontierTimeout
// names no candidate: the answers only order the attempts, and signatures still
// decide each one.
func (a *attestedHeaderBuilder) endpointFrontiers(ctx context.Context) []endpointFrontier {
	queryCtx, cancel := context.WithTimeout(ctx, stepDownFrontierTimeout)
	defer cancel()
	type result struct {
		index uint16
		root  *attestorpb.AttestedRoot
	}
	results := make(chan result, len(a.attestors))
	for _, endpoint := range a.attestors {
		go func(endpoint SigningAttestor) {
			var root *attestorpb.AttestedRoot
			// A panic in one endpoint's client costs only its candidate.
			_ = workergroup.Recover(fmt.Sprintf("attestor %d frontier", endpoint.Index), func() error {
				answer, found, err := endpoint.Client.AttestedUpTo(queryCtx, a.srcChain, a.includeProvisional)
				if err == nil && found {
					root = answer
				}
				return err
			})
			results <- result{index: endpoint.Index, root: root}
		}(endpoint)
	}
	frontiers := make([]endpointFrontier, 0, len(a.attestors))
	for range a.attestors {
		select {
		case <-queryCtx.Done():
			return frontiers
		case one := <-results:
			if one.root != nil {
				frontiers = append(frontiers, endpointFrontier{index: one.index, height: one.root.GetL2BlockNumber()})
			}
		}
	}
	return frontiers
}

// stepDownCandidates are the distinct frontiers below height, highest first.
func stepDownCandidates(frontiers []endpointFrontier, height uint64) []uint64 {
	seen := make(map[uint64]bool, len(frontiers))
	candidates := make([]uint64, 0, len(frontiers))
	for _, frontier := range frontiers {
		if frontier.height == 0 || frontier.height >= height || seen[frontier.height] {
			continue
		}
		seen[frontier.height] = true
		candidates = append(candidates, frontier.height)
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i] > candidates[j] })
	return candidates
}

// describeFrontiers names each endpoint's answer, in index order, for the
// step-down log: the endpoint that over-reports is the one above the rest.
func describeFrontiers(frontiers []endpointFrontier) string {
	if len(frontiers) == 0 {
		return "none answered"
	}
	sorted := append([]endpointFrontier(nil), frontiers...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].index < sorted[j].index })
	parts := make([]string, 0, len(sorted))
	for _, frontier := range sorted {
		parts = append(parts, fmt.Sprintf("attestor %d at %d", frontier.index, frontier.height))
	}
	return strings.Join(parts, ", ")
}

// bindToAttestation refuses a header the attestor's replica does not recognise as the
// canonical block at that height. state root AND block hash are both sent: the state
// root is what the client will trust, the block hash is what makes the answer specific
// to one block rather than to any block sharing a state root.
//
// Endpoints are queried concurrently. Individual causes retain their typed
// classification, while the aggregate is permanent only when no recoverable
// endpoint can bring the valid count up to the configured threshold.
func (a *attestedHeaderBuilder) bindToAttestation(ctx context.Context, height uint64, l2Header *types.Header) ([]IndexedAttestorSignature, error) {
	stateRoot := l2Header.Root
	blockHash := l2Header.Hash()
	statement, err := attestation.SigningBytes(a.chainID, a.router.Bytes(), a.setHash[:], height, blockHash.Bytes(), stateRoot.Bytes())
	if err != nil {
		return nil, fmt.Errorf("%s: build attestation statement: %w", a.name, err)
	}
	request := VerificationRequest{
		SrcChain: a.srcChain, BlockNumber: height, StateRoot: stateRoot.Bytes(), BlockHash: blockHash.Bytes(), RunMode: a.runMode, AttestorSetHash: a.setHash,
	}
	copy(request.L2Router[:], a.router.Bytes())

	signatures := make([]IndexedAttestorSignature, 0, a.set.Threshold)
	failures := make([]attestorFailure, 0, len(a.attestors))
	type result struct {
		index   uint16
		verdict SignedVerdict
		err     error
	}
	queryCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	results := make(chan result, len(a.attestors))
	for _, endpoint := range a.attestors {
		go func(endpoint SigningAttestor) {
			var verdict SignedVerdict
			// A panic in one endpoint's client is that endpoint's failure: the
			// signature quorum below already tolerates failed endpoints.
			err := workergroup.Recover(fmt.Sprintf("attestor %d verdict", endpoint.Index), func() error {
				var err error
				verdict, err = endpoint.Client.VerifyStateRoot(queryCtx, request)
				return err
			})
			results <- result{index: endpoint.Index, verdict: verdict, err: err}
		}(endpoint)
	}
	for range a.attestors {
		var one result
		select {
		case <-ctx.Done():
			return nil, chain.Transient(fmt.Errorf("%s: attestor signature quorum: %w", a.name, ctx.Err()))
		case one = <-results:
		}
		if one.err != nil {
			failures = append(failures, newAttestorFailure("Attest", one.index, one.err))
			continue
		}
		if !one.verdict.Valid {
			failures = append(failures, newAttestorFailure("Attest", one.index, fmt.Errorf("%w: did not recognise L2 block %d", ErrAttestorDivergence, height)))
			continue
		}
		if one.verdict.BlockNumber != height || string(one.verdict.BlockHash) != string(blockHash.Bytes()) || string(one.verdict.StateRoot) != string(stateRoot.Bytes()) {
			failures = append(failures, newAttestorFailure("Attest", one.index, fmt.Errorf("%w: returned a different block identity", ErrAttestorDivergence)))
			continue
		}
		if len(one.verdict.Signature) != attestation.SignatureLength || !ed25519.Verify(ed25519.PublicKey(a.set.PublicKeys[one.index]), statement, one.verdict.Signature) {
			failures = append(failures, newAttestorFailure("Attest", one.index, fmt.Errorf("%w: returned an invalid signature", ErrAttestorSignature)))
			continue
		}
		signatures = append(signatures, IndexedAttestorSignature{AttestorIndex: one.index, Signature: append([]byte(nil), one.verdict.Signature...)})
		if len(signatures) == int(a.set.Threshold) {
			cancel()
			sort.Slice(signatures, func(i, j int) bool { return signatures[i].AttestorIndex < signatures[j].AttestorIndex })
			return signatures, nil
		}
	}
	return nil, quorumFailure(a.name, fmt.Sprintf("valid attestor signatures for L2 block %d", height), len(signatures), 0, int(a.set.Threshold), failures)
}
