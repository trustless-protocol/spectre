package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	"relayer/chain"
	"relayer/chain/l2rollup"
	"relayer/chain/l2rollup/attestorgrpc"
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
		"l2_rpc_url": c.L2RpcUrl, "tm_rpc_url": c.TmRpcUrl,
		"attestor_addr": c.AttestorAddr, "attestor_src_chain": c.AttestorSrcChain,
		"l2_wasm_client_id": c.L2WasmClientID, "l2_ics26_client_id": c.L2ICS26ClientID,
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

	l2, err := ethclient.Dial(cfg.L2RpcUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("dial L2 rpc: %w", err)
	}
	cosmosClient, err := rpchttp.New(cfg.TmRpcUrl, "/websocket")
	if err != nil {
		l2.Close()
		return nil, nil, fmt.Errorf("create Cosmos rpc client: %w", err)
	}
	if err := cosmosClient.Start(); err != nil {
		l2.Close()
		return nil, nil, fmt.Errorf("start Cosmos rpc client: %w", err)
	}
	attestor, err := attestorgrpc.Dial(cfg.AttestorAddr)
	if err != nil {
		l2.Close()
		_ = cosmosClient.Stop()
		return nil, nil, fmt.Errorf("dial attestor: %w", err)
	}

	// Cosmos-side context for the destination (it uses only CosmosClient + the
	// signer via the tx handler); the L2 exec client fills the eth slot unused here.
	svcCtx := services.NewCtx(cosmosClient, l2)
	worker := services.NewWorker(txHandler, nil)

	// One builder for every chain: the header is the canonical L2 block plus a router
	// account proof, and nothing about that is chain-specific any more. It takes the
	// same attestor the source gates heights on, so the block it packages is bound to
	// what that attestor's replica actually has at that height.
	headerBuilder := l2rollup.NewAttestedHeaderBuilder(l2, router, attestor, headKind.RunMode(), fmt.Sprintf("l2-%s", cfg.kind))

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

	cleanup := func() {
		l2.Close()
		_ = cosmosClient.Stop()
		_ = attestor.Close()
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
