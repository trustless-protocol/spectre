package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"attestor/types/attestation"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"relayer/chain"
	"relayer/chain/avalanche"
	"relayer/chain/l2rollup"
	"relayer/chain/l2rollup/attestorgrpc"
	relayerclient "relayer/client"
	"relayer/relay"
	"relayer/services"
)

// l2ToCosmosConfig is one L2->Cosmos source's `config` block. It is self-contained:
// the L1/L2 execution RPCs, the Cosmos RPC hosting the L2 wasm client, the attestor
// sidecar address + its src_chain key (#240), the two client ids (the wasm client on
// Cosmos and the ICS26Router client on the L2), the head policy, and the rollup
// profile (whose common.l2_router the membership proofs are taken against).
type l2ToCosmosConfig struct {
	// L1RpcUrl, OpNodeRpcUrl and EthBeaconAPIURL are gone with the settlement path:
	// the client no longer verifies anything against L1, so this module never dials
	// one. That also retires the pinned-Ethereum-client dependency (#276) for L2
	// sources — there is no client here that has to be kept advancing.
	L2RpcUrl string `json:"l2_rpc_url"`
	TmRpcUrl string `json:"tm_rpc_url"`

	AttestorEndpoints []l2AttestorEndpointConfig `json:"attestor_endpoints"`
	AttestorSrcChain  string                     `json:"attestor_src_chain"`
	// Warp configures the Avalanche C-Chain path (src_chain "avalanche"): the
	// wasm client there verifies the primary network's aggregate BLS signature
	// per update, so this module runs no attestors at all — signatures come
	// from validators through the aggregator sidecar.
	Warp *l2WarpConfig `json:"warp,omitempty"`
	// Attestors is the exact immutable key set stored in the paired wasm client.
	// Every endpoint below names one key by its canonical set index.
	Attestors       attestation.AttestorConfig `json:"attestors"`
	L2WasmClientID  string                     `json:"l2_wasm_client_id"`
	L2ICS26ClientID string                     `json:"l2_ics26_client_id"`
	HeadKind        string                     `json:"head_kind"`
	// IncludeProvisional accepts attestor verdicts that have not yet been re-derived
	// from finalized L1 data. It is deliberately its own field: it used to be inferred
	// from head_kind, which conflated two independent axes — head_kind selects which
	// L2 head is read, provisional is about whether the attestor's finality re-check
	// has completed. "safe" therefore silently implied "accept provisional", with no
	// way to ask for one without the other. Defaults to true for anything but
	// finalized, preserving the previous behaviour.
	IncludeProvisional *bool `json:"include_provisional,omitempty"`
	// LogScanChunk caps the block span of a single eth_getLogs on the L2 packet
	// scan. 0 (the default) issues one call per range, which is right for any
	// provider that does not cap the span. It covered the L1 assertion scan too
	// until #347 removed that path; only the L2 scan is left to bound.
	//
	// It lives here rather than in rollup_profile because the profile is forwarded
	// verbatim into MsgCreateClient and the Rust verifier parses it with
	// deny_unknown_fields — a relayer-only knob added there would break client
	// creation.
	//
	// Known caps: Alchemy free tier 10, drpc 10_000.
	LogScanChunk  uint64          `json:"log_scan_chunk,omitempty"`
	RollupProfile json.RawMessage `json:"rollup_profile"`

	// kind is the chain family (opstack/arbitrum) resolved from the module's
	// src_chain by loadConfig; it selects the per-L2 header builder.
	kind chain.ChainType
}

// l2WarpConfig is the Avalanche warp path's own block: where to aggregate
// primary-network signatures, and the message identity the wasm client pins
// (network id + C-Chain blockchain id). quorum_num defaults to Avalanche's 67.
type l2WarpConfig struct {
	AggregatorURL string `json:"aggregator_url"`
	NetworkID     uint32 `json:"network_id"`
	SourceChainID string `json:"source_chain_id"`
	QuorumNum     uint64 `json:"quorum_num,omitempty"`
}

// validate checks the block and decodes the 32-byte source chain id.
func (w *l2WarpConfig) validate() ([32]byte, error) {
	var chainID [32]byte
	if w == nil {
		return chainID, fmt.Errorf("l2_to_cosmos config: warp block is required for an avalanche source")
	}
	if w.AggregatorURL == "" {
		return chainID, fmt.Errorf("l2_to_cosmos config: warp.aggregator_url is required")
	}
	if w.NetworkID == 0 {
		return chainID, fmt.Errorf("l2_to_cosmos config: warp.network_id is required")
	}
	raw := strings.TrimPrefix(w.SourceChainID, "0x")
	decoded, err := hex.DecodeString(raw)
	if err != nil || len(decoded) != 32 {
		return chainID, fmt.Errorf("l2_to_cosmos config: warp.source_chain_id must be 32 hex bytes")
	}
	copy(chainID[:], decoded)
	if chainID == [32]byte{} {
		return chainID, fmt.Errorf("l2_to_cosmos config: warp.source_chain_id must be non-zero")
	}
	if w.QuorumNum > 100 {
		return chainID, fmt.Errorf("l2_to_cosmos config: warp.quorum_num must be in 0..100")
	}
	return chainID, nil
}

// quorum returns the configured quorum numerator, defaulting to Avalanche's 67.
func (w *l2WarpConfig) quorum() uint64 {
	if w == nil || w.QuorumNum == 0 {
		return 67
	}
	return w.QuorumNum
}

// l2AttestorEndpointConfig maps a reachable attestor daemon to its immutable
// ClientState public-key index. Keeping the key in one canonical set prevents a
// typo from making the relayer sign for a different client than it bootstrapped.
type l2AttestorEndpointConfig struct {
	Address       string `json:"address"`
	AttestorIndex uint16 `json:"attestor_index"`
}

// validate checks everything the relay path needs, including the id of an L2 wasm
// client that must already exist on Cosmos.
func (c l2ToCosmosConfig) validate() error {
	return c.validateWith(true)
}

// validateForClientCreation is validate minus the one field create-clients-cosmos
// is about to produce. That command creates the L2 wasm client and writes its id
// back into this module, so demanding the id up front made a fresh config
// unusable: the only way forward was to invent a plausible one ("08-wasm-1"),
// which the command then overwrote — pure ceremony, and a wrong guess was accepted
// just as silently as a right one (#309).
//
// The exemption is deliberately this narrow. Every other field is needed to reach
// a chain at all, and failing on those before spending an on-chain MsgCreateClient
// is the point of validating here.
func (c l2ToCosmosConfig) validateForClientCreation() error {
	return c.validateWith(false)
}

func (c l2ToCosmosConfig) validateWith(requireWasmClientID bool) error {
	required := map[string]string{
		"l2_rpc_url": c.L2RpcUrl, "tm_rpc_url": c.TmRpcUrl,
		"attestor_src_chain": c.AttestorSrcChain,
		"l2_ics26_client_id": c.L2ICS26ClientID,
	}
	if c.kind == chain.Avalanche {
		// The warp path runs no attestors: per-update trust is the primary
		// network's aggregate signature, collected through the aggregator.
		delete(required, "attestor_src_chain")
		if _, err := c.Warp.validate(); err != nil {
			return err
		}
		if len(c.AttestorEndpoints) != 0 || len(c.Attestors.PublicKeys) != 0 {
			return fmt.Errorf("l2_to_cosmos config: an avalanche source takes no attestor fields; remove them")
		}
	}
	if requireWasmClientID {
		required["l2_wasm_client_id"] = c.L2WasmClientID
	}
	for name, v := range required {
		if v == "" {
			return fmt.Errorf("l2_to_cosmos config: %s is required", name)
		}
	}
	if len(c.RollupProfile) == 0 {
		return fmt.Errorf("l2_to_cosmos config: rollup_profile is required")
	}
	if _, err := parseHeadKind(c.HeadKind); err != nil {
		return err
	}
	if _, err := l2RouterFromProfile(c.RollupProfile); err != nil {
		return err
	}
	if _, err := l2ChainIDFromProfile(c.RollupProfile); err != nil {
		return fmt.Errorf("l2_to_cosmos config: %w", err)
	}
	if c.kind == chain.Avalanche {
		return nil
	}
	if err := c.Attestors.Validate(); err != nil {
		return fmt.Errorf("l2_to_cosmos config: invalid attestors: %w", err)
	}
	if len(c.AttestorEndpoints) < int(c.Attestors.Threshold) {
		return fmt.Errorf("l2_to_cosmos config: %d attestor_endpoints cannot satisfy threshold %d", len(c.AttestorEndpoints), c.Attestors.Threshold)
	}
	seen := make(map[uint16]struct{}, len(c.AttestorEndpoints))
	for _, endpoint := range c.AttestorEndpoints {
		if endpoint.Address == "" {
			return fmt.Errorf("l2_to_cosmos config: attestor endpoint address is required")
		}
		if int(endpoint.AttestorIndex) >= len(c.Attestors.PublicKeys) {
			return fmt.Errorf("l2_to_cosmos config: attestor endpoint index %d is outside attestors.public_keys", endpoint.AttestorIndex)
		}
		if _, duplicate := seen[endpoint.AttestorIndex]; duplicate {
			return fmt.Errorf("l2_to_cosmos config: duplicate attestor endpoint index %d", endpoint.AttestorIndex)
		}
		seen[endpoint.AttestorIndex] = struct{}{}
	}
	return nil
}

// parseHeadKind maps the config string to the anti-reorg head policy (default Safe).
func parseHeadKind(s string) (l2rollup.HeadKind, error) {
	switch s {
	case "", "safe":
		return l2rollup.Safe, nil
	case "unsafe":
		return l2rollup.Unsafe, nil
	case "finalized":
		return l2rollup.Finalized, nil
	default:
		return 0, fmt.Errorf("l2_to_cosmos config: head_kind %q must be unsafe|safe|finalized", s)
	}
}

// ethClientBeacon names one module's view of a Cosmos-side Ethereum light client:

func l2RouterFromProfile(profile json.RawMessage) (common.Address, error) {
	var pc struct {
		Common struct {
			L2Router string `json:"l2_router"`
		} `json:"common"`
	}
	if err := json.Unmarshal(profile, &pc); err != nil {
		return common.Address{}, fmt.Errorf("l2_to_cosmos config: parse rollup_profile.common: %w", err)
	}
	if !common.IsHexAddress(pc.Common.L2Router) {
		return common.Address{}, fmt.Errorf("l2_to_cosmos config: rollup_profile.common.l2_router %q is not a valid hex address", pc.Common.L2Router)
	}
	return common.HexToAddress(pc.Common.L2Router), nil
}

// buildL2ToCosmosModule constructs one L2->Cosmos relay module: an l2rollup Source
// (gated on the attestor), a Cosmos Destination hosting the L2 wasm client, the
// per-L2 header builder, and timeout recovery through the matching Cosmos->L2
// SpectreClient/ICS26Router return path. It returns the module and a cleanup that
// closes the dialed clients. The caller runs module.Run(ctx) on its own goroutine.
func buildL2ToCosmosModule(logger *zap.Logger, cfg l2ToCosmosConfig, txHandler services.TransactionHandler, timeoutReturn l2TimeoutReturnPath) (*relay.Module, func(), error) {
	// Before anything dials: a retired env key that silently sizes nothing is the
	// failure NFR 10 forbids. Checked here rather than globally because the key
	// only ever affected this path, so a deployment without an l2_to_cosmos
	// module has nothing to correct.
	if err := l2rollup.RejectRetiredLookbackEnv(); err != nil {
		return nil, nil, err
	}
	headKind, err := parseHeadKind(cfg.HeadKind)
	if err != nil {
		return nil, nil, err
	}
	router, err := l2RouterFromProfile(cfg.RollupProfile)
	if err != nil {
		return nil, nil, err
	}
	if err := timeoutReturn.validate(); err != nil {
		return nil, nil, err
	}
	chainID, err := l2ChainIDFromProfile(cfg.RollupProfile)
	if err != nil {
		return nil, nil, fmt.Errorf("l2_to_cosmos config: %w", err)
	}
	includeProvisional := headKind != l2rollup.Finalized
	if cfg.IncludeProvisional != nil {
		includeProvisional = *cfg.IncludeProvisional
	}

	l2, err := relayerclient.DialEthRPC(context.Background(), cfg.L2RpcUrl, relayerclient.DefaultRPCTimeout)
	if err != nil {
		return nil, nil, fmt.Errorf("dial L2 rpc: %w", err)
	}
	cosmosClient, err := relayerclient.DialCosmosRPC(cfg.TmRpcUrl, "/websocket", relayerclient.DefaultRPCTimeout)
	if err != nil {
		l2.Close()
		return nil, nil, fmt.Errorf("create Cosmos rpc client: %w", err)
	}
	if err := cosmosClient.Start(); err != nil {
		l2.Close()
		return nil, nil, fmt.Errorf("start Cosmos rpc client: %w", err)
	}
	attestors := make([]*attestorgrpc.Client, 0, len(cfg.AttestorEndpoints))
	signingAttestors := make([]l2rollup.SigningAttestor, 0, len(cfg.AttestorEndpoints))
	for _, endpoint := range cfg.AttestorEndpoints {
		attestor, err := attestorgrpc.Dial(endpoint.Address)
		if err != nil {
			l2.Close()
			if stopErr := cosmosClient.Stop(); stopErr != nil {
				log.Printf("failed to terminate cosmos client: %v", stopErr)
			}
			for _, opened := range attestors {
				_ = opened.Close()
			}
			return nil, nil, fmt.Errorf("dial attestor index %d (%s): %w", endpoint.AttestorIndex, endpoint.Address, err)
		}
		attestors = append(attestors, attestor)
		signingAttestors = append(signingAttestors, l2rollup.SigningAttestor{Index: endpoint.AttestorIndex, Client: attestor})
	}

	// Cosmos-side context for the destination (it uses only CosmosClient + the
	// signer via the tx handler); the L2 exec client fills the eth slot unused here.
	svcCtx := services.CosmosEndpoint{Client: cosmosClient}
	worker := services.NewWorker(txHandler, nil)

	// The builder and the height gate differ by trust model. Rollups: the
	// attested-header builder binds each update to a threshold of independent
	// replicas, and RelayableHeight gates on their quorum frontier. Avalanche's
	// C-Chain (an L1; acceptance is finality): the warp builder collects the
	// primary network's aggregate BLS signature through the aggregator sidecar,
	// no attestors exist, and the frontier is the finalized (= accepted) head.
	var headerBuilder l2rollup.HeaderBuilder
	var frontier l2rollup.AttestationFrontier
	if cfg.kind == chain.Avalanche {
		sourceChainID, warpErr := cfg.Warp.validate()
		if warpErr == nil {
			headerBuilder, warpErr = avalanche.NewWarpHeaderBuilder(
				l2, router, cfg.Warp.NetworkID, sourceChainID, cfg.Warp.AggregatorURL, cfg.Warp.quorum(), "avalanche-warp")
		}
		if warpErr != nil {
			l2.Close()
			if stopErr := cosmosClient.Stop(); stopErr != nil {
				log.Printf("failed to terminate cosmos client: %v", stopErr)
			}
			return nil, nil, warpErr
		}
	} else {
		headerBuilder, err = l2rollup.NewAttestedHeaderBuilder(l2, router, chainID, cfg.Attestors, signingAttestors, cfg.AttestorSrcChain, headKind.RunMode(), fmt.Sprintf("l2-%s", cfg.kind))
		if err != nil {
			l2.Close()
			if stopErr := cosmosClient.Stop(); stopErr != nil {
				log.Printf("failed to terminate cosmos client: %v", stopErr)
			}
			for _, opened := range attestors {
				_ = opened.Close()
			}
			return nil, nil, err
		}
		frontier, err = l2rollup.NewQuorumAttestationFrontier(signingAttestors, cfg.Attestors.Threshold)
		if err != nil {
			l2.Close()
			if stopErr := cosmosClient.Stop(); stopErr != nil {
				log.Printf("failed to terminate cosmos client: %v", stopErr)
			}
			for _, opened := range attestors {
				_ = opened.Close()
			}
			return nil, nil, err
		}
	}
	trackL2Pending, untrackL2Pending := l2PendingTrackerHooks(timeoutReturn.svc)
	source := l2rollup.NewSource(cfg.kind, l2, headKind, cfg.L2ICS26ClientID, cfg.L2WasmClientID, cfg.AttestorSrcChain, router, frontier, includeProvisional).
		WithLogScanChunk(cfg.LogScanChunk).
		// AckPacket / TimeoutPacket on the L2 close a packet's lifecycle. Reading
		// them lets the tracker drop a packet another relayer settled without
		// waiting for a timeout scan to query for it. Same hook the module uses to
		// untrack after its own successful relay -- one settlement path, two ways
		// of learning about it.
		WithSettleHook(untrackL2Pending).
		// The acknowledged half also closes the owed-acknowledgement record. The
		// same Services instance backs both directions of this pair, which is what
		// lets a debt one side recorded be settled from the other.
		WithAckSettleHook(settleL2OwedAck(timeoutReturn.svc))
	if cfg.kind == chain.Avalanche {
		// The settled-height proof mapping is the Avalanche provider's rule
		// (asynchronous execution); the shared source stays chain-agnostic.
		source = source.WithProofHeightResolver(avalanche.NewSettledProofHeightResolver(l2.Client()))
	}
	dest := l2rollup.NewDestination(worker, svcCtx, cfg.L2WasmClientID)
	builder := l2rollup.NewBuilder(headerBuilder)

	ackDue, ackSettled := ackWatchHooks(timeoutReturn.svc)
	module := relay.NewModule(
		fmt.Sprintf("%s->cosmos", cfg.kind),
		cfg.L2ICS26ClientID,
		source, dest, builder,
		relay.WithTimeoutScanner(0, func(c context.Context) {
			timeoutReturn.svc.ScanL2Timeouts(c, timeoutReturn.deps.Cosmos, timeoutReturn.deps.EVM, timeoutReturn.deps.IDs.CosmosOnEVM)
		}),
		relay.WithPacketTracker(trackL2Pending, untrackL2Pending),
		// The mirror of the cosmos->l2 hooks, over the SAME Services: that leg
		// records what a delivered receive owes, this one settles it when the
		// acknowledgement returns. The overdue reporter lives here rather than on
		// the outbound leg because this module exists exactly when the pair is
		// complete -- a forward-only deployment has no return leg and must not
		// alarm about acknowledgements it was never going to relay.
		relay.WithAckWatch(0, 0, ackDue, ackSettled, timeoutReturn.svc.OverdueAcks),
	)
	if cfg.kind == chain.Avalanche {
		logger.Sugar().Infof("avalanche->cosmos source: warp aggregation via %s (quorum=%d wasm_client=%s head=%s)",
			cfg.Warp.AggregatorURL, cfg.Warp.quorum(), cfg.L2WasmClientID, cfg.HeadKind)
	} else {
		logger.Sugar().Infof("l2->cosmos source: %s (attestors=%d threshold=%d src_chain=%s wasm_client=%s head=%s)",
			cfg.kind, len(cfg.AttestorEndpoints), cfg.Attestors.Threshold, cfg.AttestorSrcChain, cfg.L2WasmClientID, cfg.HeadKind)
	}

	cleanup := func() {
		l2.Close()
		if err := cosmosClient.Stop(); err != nil {
			log.Printf("failed to terminate cosmos client: %v", err)
		}
		for _, attestor := range attestors {
			if err := attestor.Close(); err != nil {
				log.Printf("failed to close attestor client: %v", err)
			}
		}
	}
	return module, cleanup, nil
}

// settleL2OwedAck closes an owed-acknowledgement record from a terminal L2
// AckPacket, whoever relayed it.
func settleL2OwedAck(svc *services.Services) func(raw []byte) {
	return func(raw []byte) {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[AckWatch] settle owed ack from terminal L2 event: decode packet: %v", err)
			return
		}
		svc.ClearAckDue(pkt)
	}
}

func l2PendingTrackerHooks(svc *services.Services) (relay.TrackFunc, relay.UntrackFunc) {
	trackL2Pending := func(raw []byte, height uint64) bool {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[l2->cosmos Relay] track pending: decode packet: %v", err)
			return false
		}
		return svc.TrackL2Pending(pkt, height)
	}
	untrackL2Pending := func(raw []byte) {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[l2->cosmos Relay] untrack pending: decode packet: %v", err)
			return
		}
		svc.UntrackL2Pending(pkt)
	}
	return trackL2Pending, untrackL2Pending
}

// runL2Engine drives one L2->Cosmos module until ctx is cancelled. A cancelled
// context is a clean shutdown, not a relay failure.
func runL2Engine(ctx context.Context, module *relay.Module) error {
	return module.Run(ctx)
}
