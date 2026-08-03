package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap"

	"relayer/chain"
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
	L1RpcUrl     string `json:"l1_rpc_url"`
	L2RpcUrl     string `json:"l2_rpc_url"`
	OpNodeRpcUrl string `json:"op_node_rpc_url"` // op-node, for optimism_outputAtBlock (OP-Stack only)
	TmRpcUrl     string `json:"tm_rpc_url"`
	// EthBeaconAPIURL is the L1 beacon REST endpoint used to keep the pinned
	// Ethereum light client on Cosmos advancing. That client is this module's trust
	// anchor: the header builder proves the rollup contracts at whatever L1 block it
	// trusts, so a stalled client makes every proof unservable (#276). Modules that
	// pin the same ethereum_client must agree on this URL.
	EthBeaconAPIURL  string `json:"eth_beacon_api_url"`
	AttestorAddr     string `json:"attestor_addr"`
	AttestorSrcChain string `json:"attestor_src_chain"`
	L2WasmClientID   string `json:"l2_wasm_client_id"`
	L2ICS26ClientID  string `json:"l2_ics26_client_id"`
	HeadKind         string `json:"head_kind"`
	// IncludeProvisional accepts attestor verdicts that have not yet been re-derived
	// from finalized L1 data. It is deliberately its own field: it used to be inferred
	// from head_kind, which conflated two independent axes — head_kind selects which
	// L2 head is read, provisional is about whether the attestor's finality re-check
	// has completed. "safe" therefore silently implied "accept provisional", with no
	// way to ask for one without the other. Defaults to true for anything but
	// finalized, preserving the previous behaviour.
	IncludeProvisional *bool           `json:"include_provisional,omitempty"`
	RollupProfile      json.RawMessage `json:"rollup_profile"`

	// kind is the chain family (opstack/arbitrum) resolved from the module's
	// src_chain by loadConfig; it selects the per-L2 header builder.
	kind chain.ChainType
}

func (c l2ToCosmosConfig) validate() error {
	for name, v := range map[string]string{
		"l1_rpc_url": c.L1RpcUrl, "l2_rpc_url": c.L2RpcUrl, "tm_rpc_url": c.TmRpcUrl,
		"attestor_addr": c.AttestorAddr, "attestor_src_chain": c.AttestorSrcChain,
		"eth_beacon_api_url": c.EthBeaconAPIURL,
		"l2_wasm_client_id":  c.L2WasmClientID, "l2_ics26_client_id": c.L2ICS26ClientID,
	} {
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

// l2RouterFromProfile extracts rollup_profile.common.l2_router (the L2 ICS26Router).
// l1ClientIDFromProfile reads rollup_profile.common.ethereum_client.client_id — the
// Ethereum light client on Cosmos this L2 client is anchored to. Several L2 modules
// normally share one such client (create-clients-cosmos injects the same id into
// every --l2-config), which is what makes the shared-refresher grouping meaningful.
func l1ClientIDFromProfile(profile json.RawMessage) (string, error) {
	var pc struct {
		Common struct {
			EthereumClient struct {
				ClientID string `json:"client_id"`
			} `json:"ethereum_client"`
		} `json:"common"`
	}
	if err := json.Unmarshal(profile, &pc); err != nil {
		return "", fmt.Errorf("l2_to_cosmos config: parse rollup_profile.common: %w", err)
	}
	if pc.Common.EthereumClient.ClientID == "" {
		return "", fmt.Errorf("l2_to_cosmos config: rollup_profile.common.ethereum_client.client_id is required")
	}
	return pc.Common.EthereumClient.ClientID, nil
}

// ethClientBeacon names one module's view of a Cosmos-side Ethereum light client:
// which client it advances, from which beacon, and a label for the error message.
type ethClientBeacon struct {
	clientID  string
	beaconURL string
	module    string
}

// validateSharedEthClients rejects a config where two modules that advance the SAME
// Cosmos-side Ethereum client disagree on the beacon endpoint. One client driven from
// two sources of truth is divergent and racy.
//
// Both families must be considered, not just the L2 ones: an eth_to_cosmos module
// advances the client named by its paired cosmos_to_eth.cosmos_wasm_client_id, using
// eth_to_cosmos.eth_beacon_api_url, and an L2 module can pin that very same client
// through rollup_profile.common.ethereum_client.client_id. Checking only the L2 list
// let that pair diverge silently — and config.example.json already ships two different
// beacon URLs.
func validateSharedEthClients(l2s []l2ToCosmosConfig, others ...ethClientBeacon) error {
	seen := map[string]ethClientBeacon{}
	claims := make([]ethClientBeacon, 0, len(l2s)+len(others))
	for i := range l2s {
		clientID, err := l1ClientIDFromProfile(l2s[i].RollupProfile)
		if err != nil {
			return err
		}
		claims = append(claims, ethClientBeacon{clientID, l2s[i].EthBeaconAPIURL, "l2_to_cosmos " + l2s[i].AttestorSrcChain})
	}
	for _, o := range others {
		if o.clientID == "" || o.beaconURL == "" {
			continue // nothing to compare against
		}
		claims = append(claims, o)
	}

	for _, c := range claims {
		prev, ok := seen[c.clientID]
		if !ok {
			seen[c.clientID] = c
			continue
		}
		if prev.beaconURL != c.beaconURL {
			return fmt.Errorf(
				"modules advancing Ethereum client %q disagree on eth_beacon_api_url: %s uses %q, %s uses %q; "+
					"modules sharing a client must share its beacon endpoint",
				c.clientID, prev.module, prev.beaconURL, c.module, c.beaconURL)
		}
	}
	return nil
}

// validateL2SourceEthClient refuses to start when the config names a different
// Ethereum client than the L2 wasm client was created with.
//
// create-clients-cosmos writes the id into both places, but a config assembled by
// hand or carried over from an earlier devnet run can still drift — and drift here
// is silent and total: header builds pin proofs to the client the config names while
// the contract queries the client it was created with, so every update fails with
// "IBC host query failed: codespace: undefined, code: 1", which names neither. In
// production the pinned client is also the one nobody advances, so it expires.
// One query at startup turns all of that into a sentence.
func validateL2SourceEthClient(cosmosClient *rpchttp.HTTP, l2WasmClientID, configuredL1ClientID string) error {
	pinned, err := relayerclient.GetL2PinnedEthClientID(cosmosClient, l2WasmClientID)
	if err != nil {
		return fmt.Errorf("read the Ethereum client pinned by L2 client %s: %w", l2WasmClientID, err)
	}
	if pinned != configuredL1ClientID {
		return fmt.Errorf(
			"l2_to_cosmos config: rollup_profile.common.ethereum_client.client_id is %q but L2 wasm client %s was created against %q; "+
				"the client state is authoritative — set the config to %q, or re-create the clients",
			configuredL1ClientID, l2WasmClientID, pinned, pinned)
	}
	return nil
}

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
	includeProvisional := headKind != l2rollup.Finalized
	if cfg.IncludeProvisional != nil {
		includeProvisional = *cfg.IncludeProvisional
	}

	l1, err := ethclient.Dial(cfg.L1RpcUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("dial L1 rpc: %w", err)
	}
	l2, err := ethclient.Dial(cfg.L2RpcUrl)
	if err != nil {
		l1.Close()
		return nil, nil, fmt.Errorf("dial L2 rpc: %w", err)
	}
	cosmosClient, err := rpchttp.New(cfg.TmRpcUrl, "/websocket")
	if err != nil {
		l1.Close()
		l2.Close()
		return nil, nil, fmt.Errorf("create Cosmos rpc client: %w", err)
	}
	if err := cosmosClient.Start(); err != nil {
		l1.Close()
		l2.Close()
		return nil, nil, fmt.Errorf("start Cosmos rpc client: %w", err)
	}
	attestor, err := attestorgrpc.Dial(cfg.AttestorAddr)
	if err != nil {
		l1.Close()
		l2.Close()
		_ = cosmosClient.Stop()
		return nil, nil, err
	}

	// Cosmos-side context for the destination (it uses only CosmosClient + the
	// signer via the tx handler); the L2 exec client fills the eth slot unused here.
	svcCtx := services.NewCtx(cosmosClient, l2)
	worker := services.NewWorker(txHandler, nil)

	// Wire the on-demand ETH client updater: the header builders ask it for freshness
	// right before they prove against that client's L1 block (#276).
	l1ClientID, err := l1ClientIDFromProfile(cfg.RollupProfile)
	if err != nil {
		l1.Close()
		l2.Close()
		_ = cosmosClient.Stop()
		_ = attestor.Close()
		return nil, nil, err
	}
	if err := validateL2SourceEthClient(cosmosClient, cfg.L2WasmClientID, l1ClientID); err != nil {
		l1.Close()
		l2.Close()
		_ = cosmosClient.Stop()
		_ = attestor.Close()
		return nil, nil, err
	}
	updater, updaterCleanup, err := buildEthClientUpdater(cfg, l1ClientID, txHandler,
		func(c context.Context) (uint64, error) { return l1.BlockNumber(c) })
	if err != nil {
		l1.Close()
		l2.Close()
		_ = cosmosClient.Stop()
		_ = attestor.Close()
		return nil, nil, err
	}
	cosmosReader := &cosmosEthClientReader{cosmos: cosmosClient, updater: updater}

	var headerBuilder l2rollup.HeaderBuilder
	switch cfg.kind {
	case chain.OPStack:
		opProfile, perr := parseOPProfile(cfg.RollupProfile)
		if perr != nil {
			l1.Close()
			l2.Close()
			_ = cosmosClient.Stop()
			_ = attestor.Close()
			updaterCleanup()
			return nil, nil, perr
		}
		if cfg.OpNodeRpcUrl == "" {
			l1.Close()
			l2.Close()
			_ = cosmosClient.Stop()
			_ = attestor.Close()
			updaterCleanup()
			return nil, nil, fmt.Errorf("l2_to_cosmos config: op_node_rpc_url is required for opstack (optimism_outputAtBlock)")
		}
		opNode, derr := rpc.Dial(cfg.OpNodeRpcUrl)
		if derr != nil {
			l1.Close()
			l2.Close()
			_ = cosmosClient.Stop()
			_ = attestor.Close()
			updaterCleanup()
			return nil, nil, fmt.Errorf("dial op-node rpc: %w", derr)
		}
		headerBuilder = l2rollup.NewOPStackHeaderBuilder(
			l1, l2, opNode, cosmosReader, attestor, cfg.AttestorSrcChain, opProfile, includeProvisional)
	case chain.Arbitrum:
		closeAll := func() {
			l1.Close()
			l2.Close()
			_ = cosmosClient.Stop()
			_ = attestor.Close()
			updaterCleanup()
		}
		boldProfile, berr := parseArbBoldProfile(cfg.RollupProfile)
		if berr != nil {
			closeAll()
			return nil, nil, berr
		}
		headerBuilder = l2rollup.NewArbitrumBoldHeaderBuilder(l1, l2, cosmosReader, boldProfile)
	default:
		l1.Close()
		l2.Close()
		_ = cosmosClient.Stop()
		_ = attestor.Close()
		return nil, nil, fmt.Errorf("l2_to_cosmos config: unsupported chain kind %q", cfg.kind)
	}

	source := l2rollup.NewSource(cfg.kind, l2, headKind, cfg.L2ICS26ClientID, cfg.L2WasmClientID, cfg.AttestorSrcChain, router, attestor, includeProvisional)
	dest := l2rollup.NewDestination(worker, svcCtx, cfg.L2WasmClientID)
	builder := l2rollup.NewBuilder(headerBuilder)
	trackL2Pending, untrackL2Pending := l2PendingTrackerHooks(timeoutReturn.svc)

	module := relay.NewModule(
		fmt.Sprintf("%s->cosmos", cfg.kind),
		cfg.L2ICS26ClientID,
		source, dest, builder,
		relay.WithTimeoutScanner(0, func(c context.Context) { timeoutReturn.svc.ScanL2Timeouts(c, timeoutReturn.ctx) }),
		relay.WithPacketTracker(trackL2Pending, untrackL2Pending),
	)
	logger.Sugar().Infof("l2->cosmos source: %s (attestor=%s src_chain=%s wasm_client=%s head=%s)",
		cfg.kind, cfg.AttestorAddr, cfg.AttestorSrcChain, cfg.L2WasmClientID, cfg.HeadKind)

	// Keep the pinned Ethereum client advancing even when no packets are flowing —
	// the on-demand path only runs during a header build, so an idle relayer would
	// let the client drift out of the L1's state window and eventually expire.
	freshCtx, stopFreshness := context.WithCancel(context.Background())
	go runEthClientFreshnessLoop(freshCtx, cosmosReader, l1ClientID)

	cleanup := func() {
		stopFreshness()
		l1.Close()
		l2.Close()
		_ = cosmosClient.Stop()
		_ = attestor.Close()
		updaterCleanup()
	}
	return module, cleanup, nil
}

func l2PendingTrackerHooks(svc *services.Services) (relay.TrackFunc, relay.UntrackFunc) {
	trackL2Pending := func(raw []byte, height uint64) {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[adapter l2->cosmos] track pending: decode packet: %v", err)
			return
		}
		svc.TrackL2Pending(pkt, height)
	}
	untrackL2Pending := func(raw []byte) {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[adapter l2->cosmos] untrack pending: decode packet: %v", err)
			return
		}
		svc.UntrackL2Pending(pkt)
	}
	return trackL2Pending, untrackL2Pending
}

// runL2Engine drives one L2->Cosmos module until ctx is cancelled. A cancelled
// context is a clean shutdown, not a relay failure.
func runL2Engine(ctx context.Context, module *relay.Module) error {
	if err := module.Run(ctx); err != nil && ctx.Err() == nil {
		return err
	}
	return nil
}

// parseOPProfile extracts the OP header builder's fields from the rollup profile.
func parseOPProfile(profile json.RawMessage) (l2rollup.OPProfile, error) {
	var p struct {
		Common struct {
			L2Router       string `json:"l2_router"`
			EthereumClient struct {
				ClientID string `json:"client_id"`
			} `json:"ethereum_client"`
		} `json:"common"`
		DisputeGameFactory string `json:"dispute_game_factory"`
		GameListSlot       string `json:"game_list_slot"`
	}
	if err := json.Unmarshal(profile, &p); err != nil {
		return l2rollup.OPProfile{}, fmt.Errorf("parse OP rollup_profile: %w", err)
	}
	if !common.IsHexAddress(p.DisputeGameFactory) {
		return l2rollup.OPProfile{}, fmt.Errorf("rollup_profile.dispute_game_factory %q is not a valid hex address", p.DisputeGameFactory)
	}
	if !common.IsHexAddress(p.Common.L2Router) {
		return l2rollup.OPProfile{}, fmt.Errorf("rollup_profile.common.l2_router %q is not a valid hex address", p.Common.L2Router)
	}
	if p.Common.EthereumClient.ClientID == "" {
		return l2rollup.OPProfile{}, fmt.Errorf("rollup_profile.common.ethereum_client.client_id is required")
	}
	return l2rollup.OPProfile{
		DisputeGameFactory: common.HexToAddress(p.DisputeGameFactory),
		GameListSlot:       common.HexToHash(p.GameListSlot),
		L2Router:           common.HexToAddress(p.Common.L2Router),
		L1ClientID:         p.Common.EthereumClient.ClientID,
	}, nil
}

// arbCommonFields validates and returns the RollupCore, L2 router, and shared ETH
// client id every Arbitrum profile carries (common.* + top-level rollup).
func arbCommonFields(rollup, l2Router, clientID string) (common.Address, common.Address, string, error) {
	if !common.IsHexAddress(rollup) {
		return common.Address{}, common.Address{}, "", fmt.Errorf("rollup_profile.rollup %q is not a valid hex address", rollup)
	}
	if !common.IsHexAddress(l2Router) {
		return common.Address{}, common.Address{}, "", fmt.Errorf("rollup_profile.common.l2_router %q is not a valid hex address", l2Router)
	}
	if clientID == "" {
		return common.Address{}, common.Address{}, "", fmt.Errorf("rollup_profile.common.ethereum_client.client_id is required")
	}
	return common.HexToAddress(rollup), common.HexToAddress(l2Router), clientID, nil
}

// parseArbBoldProfile extracts the BoLD v2 header builder's fields from the rollup
// profile. It mirrors the arbitrum-verifier Profile (packages/arbitrum-verifier): the L1
// RollupCore (`rollup`), the BoLD `_assertions` mapping slot (protocol.value), the L2
// router, and the shared ETH client id. assertion_scan_range is an optional relayer-only
// override.
//
// protocol.type is required and must be "bold_v2": it is the tagged RollupProtocol
// discriminant, and the legacy Nitro variant it used to select no longer exists. A
// profile still carrying the old value is rejected here rather than parsed into a
// BoLD profile it does not describe.
func parseArbBoldProfile(profile json.RawMessage) (l2rollup.ArbBoldProfile, error) {
	var p struct {
		Common struct {
			L2Router       string `json:"l2_router"`
			EthereumClient struct {
				ClientID string `json:"client_id"`
			} `json:"ethereum_client"`
		} `json:"common"`
		Rollup   string `json:"rollup"`
		Protocol struct {
			Type  string `json:"type"`
			Value struct {
				AssertionsMappingSlot string `json:"assertions_mapping_slot"`
				// AssertionScanRange is an optional relayer-only override of how far
				// back BuildHeader scans AssertionCreated logs (0 → builder default).
				AssertionScanRange uint64 `json:"assertion_scan_range"`
			} `json:"value"`
		} `json:"protocol"`
	}
	if err := json.Unmarshal(profile, &p); err != nil {
		return l2rollup.ArbBoldProfile{}, fmt.Errorf("parse Arbitrum bold_v2 rollup_profile: %w", err)
	}
	if p.Protocol.Type != "bold_v2" {
		// Name the offending value. The common way to reach this is a config carried
		// over from the legacy Nitro path ("legacy_nitro"), and an error that does not
		// echo what it read sends the reader looking for a missing field instead.
		got := p.Protocol.Type
		if got == "" {
			got = "<absent>"
		}
		return l2rollup.ArbBoldProfile{}, fmt.Errorf(
			"rollup_profile.protocol.type is %s, must be %q (the legacy Nitro protocol was removed)",
			got, "bold_v2",
		)
	}
	rollup, l2Router, clientID, err := arbCommonFields(p.Rollup, p.Common.L2Router, p.Common.EthereumClient.ClientID)
	if err != nil {
		return l2rollup.ArbBoldProfile{}, err
	}
	if p.Protocol.Value.AssertionsMappingSlot == "" {
		return l2rollup.ArbBoldProfile{}, fmt.Errorf("rollup_profile.protocol.value.assertions_mapping_slot is required")
	}
	return l2rollup.ArbBoldProfile{
		RollupCore:            rollup,
		AssertionsMappingSlot: common.HexToHash(p.Protocol.Value.AssertionsMappingSlot),
		L2Router:              l2Router,
		L1ClientID:            clientID,
		AssertionScanRange:    p.Protocol.Value.AssertionScanRange,
	}, nil
}

// cosmosEthClientReader adapts the Cosmos RPC client to the header builder's
// cosmosClientStateReader: it reads the shared ETH (08-wasm) client's trusted L1
// slot + execution block from Cosmos.
type cosmosEthClientReader struct {
	cosmos *rpchttp.HTTP
	// updater advances the pinned ETH client when a header build finds it too far
	// behind the L1 head. Nil when the deployment has no way to advance it (no beacon
	// endpoint), in which case builds proceed with whatever height the client has.
	updater *ethClientUpdater
}

// UpdateEthClientIfStale advances the pinned client if it has fallen behind. It never
// fails the build: the client may still be recent enough, or another relay direction
// may be advancing it, so the error is only surfaced for logging.
func (r *cosmosEthClientReader) UpdateEthClientIfStale(ctx context.Context, l1ClientID string) error {
	if r.updater == nil {
		return nil
	}
	cs, err := relayerclient.GetEthereumClientState(r.cosmos, l1ClientID)
	if err != nil {
		return fmt.Errorf("eth client %s: read client state: %w", l1ClientID, err)
	}
	if err := r.updater.updateIfStale(ctx, l1ClientID, cs.LatestExecutionBlockNumber); err != nil {
		log.Printf("[EthClientUpdate] %v", err)
		return err
	}
	return nil
}

func (r *cosmosEthClientReader) EthClientLatestSlotAndBlock(l1ClientID string) (uint64, uint64, error) {
	cs, err := relayerclient.GetEthereumClientState(r.cosmos, l1ClientID)
	if err != nil {
		return 0, 0, err
	}
	return cs.LatestSlot, cs.LatestExecutionBlockNumber, nil
}

// runEthClientFreshnessLoop keeps the pinned Ethereum client advancing while no
// packets are in flight.
//
// UpdateEthClientIfStale is only reached from a header build, so an idle relayer
// never advances the client at all. That is not just a latency problem: once the
// pinned block falls out of the execution node's state window, eth_getProof at that
// block fails, so the first packet to arrive after a quiet stretch cannot be relayed
// either — and left long enough the client expires. It mirrors the anti-expiry
// refresh the Cosmos→ETH direction already runs for its own destination client.
//
// The tick only CHECKS; updateIfStale no-ops unless the client is more than
// ethClientUpdateLag behind, and shares this updater's cooldown and mutex, so this
// never races or duplicates the on-demand path.
func runEthClientFreshnessLoop(ctx context.Context, reader *cosmosEthClientReader, l1ClientID string) {
	if reader == nil || reader.updater == nil {
		return // no beacon endpoint configured — nothing can advance the client here
	}
	ticker := time.NewTicker(ethClientPollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		// Errors are already logged by UpdateEthClientIfStale; a failed poll must not
		// stop the loop, since the next tick is exactly the retry.
		_ = reader.UpdateEthClientIfStale(ctx, l1ClientID)
	}
}
