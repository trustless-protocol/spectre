package main

import (
	"context"
	"encoding/json"
	"fmt"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
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
	L1RpcUrl         string          `json:"l1_rpc_url"`
	L2RpcUrl         string          `json:"l2_rpc_url"`
	OpNodeRpcUrl     string          `json:"op_node_rpc_url"` // op-node, for optimism_outputAtBlock (OP-Stack only)
	TmRpcUrl         string          `json:"tm_rpc_url"`
	AttestorAddr     string          `json:"attestor_addr"`
	AttestorSrcChain string          `json:"attestor_src_chain"`
	L2WasmClientID   string          `json:"l2_wasm_client_id"`
	L2ICS26ClientID  string          `json:"l2_ics26_client_id"`
	HeadKind         string          `json:"head_kind"`
	RollupProfile    json.RawMessage `json:"rollup_profile"`

	// kind is the chain family (opstack/arbitrum) resolved from the module's
	// src_chain by loadConfig; it selects the per-L2 header builder.
	kind chain.ChainType
}

func (c l2ToCosmosConfig) validate() error {
	for name, v := range map[string]string{
		"l1_rpc_url": c.L1RpcUrl, "l2_rpc_url": c.L2RpcUrl, "tm_rpc_url": c.TmRpcUrl,
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

// l2RouterFromProfile extracts rollup_profile.common.l2_router (the L2 ICS26Router).
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
// (gated on the attestor), a Cosmos Destination hosting the L2 wasm client, and the
// per-L2 header builder. It returns the module and a cleanup that closes the dialed
// clients. The caller runs module.Run(ctx) on its own goroutine.
func buildL2ToCosmosModule(logger *zap.Logger, cfg l2ToCosmosConfig, txHandler services.TransactionHandler) (*relay.Module, func(), error) {
	headKind, err := parseHeadKind(cfg.HeadKind)
	if err != nil {
		return nil, nil, err
	}
	router, err := l2RouterFromProfile(cfg.RollupProfile)
	if err != nil {
		return nil, nil, err
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

	cosmosReader := &cosmosEthClientReader{cosmos: cosmosClient}

	var headerBuilder l2rollup.HeaderBuilder
	switch cfg.kind {
	case chain.OPStack:
		opProfile, perr := parseOPProfile(cfg.RollupProfile)
		if perr != nil {
			l1.Close()
			l2.Close()
			_ = cosmosClient.Stop()
			_ = attestor.Close()
			return nil, nil, perr
		}
		if cfg.OpNodeRpcUrl == "" {
			l1.Close()
			l2.Close()
			_ = cosmosClient.Stop()
			_ = attestor.Close()
			return nil, nil, fmt.Errorf("l2_to_cosmos config: op_node_rpc_url is required for opstack (optimism_outputAtBlock)")
		}
		opNode, derr := rpc.Dial(cfg.OpNodeRpcUrl)
		if derr != nil {
			l1.Close()
			l2.Close()
			_ = cosmosClient.Stop()
			_ = attestor.Close()
			return nil, nil, fmt.Errorf("dial op-node rpc: %w", derr)
		}
		// includeProvisional must match the source's head policy so the builder proves
		// the same game the source gated on (Finalized = confirmed only).
		headerBuilder = l2rollup.NewOPStackHeaderBuilder(
			l1, l2, opNode, cosmosReader, attestor, cfg.AttestorSrcChain, headKind != l2rollup.Finalized, opProfile)
	case chain.Arbitrum:
		closeAll := func() {
			l1.Close()
			l2.Close()
			_ = cosmosClient.Stop()
			_ = attestor.Close()
		}
		protocol, perr := arbRollupProtocol(cfg.RollupProfile)
		if perr != nil {
			closeAll()
			return nil, nil, perr
		}
		switch protocol {
		case "bold_v2":
			boldProfile, berr := parseArbBoldProfile(cfg.RollupProfile)
			if berr != nil {
				closeAll()
				return nil, nil, berr
			}
			headerBuilder = l2rollup.NewArbitrumBoldHeaderBuilder(l1, l2, cosmosReader, boldProfile)
		case "legacy_nitro":
			legacyProfile, lerr := parseArbLegacyProfile(cfg.RollupProfile)
			if lerr != nil {
				closeAll()
				return nil, nil, lerr
			}
			// includeProvisional must match the source's head policy so the builder
			// resolves the same node the source gated on (Finalized = confirmed only).
			headerBuilder = l2rollup.NewArbitrumLegacyHeaderBuilder(
				l1, l2, cosmosReader, attestor, cfg.AttestorSrcChain, headKind != l2rollup.Finalized, legacyProfile)
		default:
			closeAll()
			return nil, nil, fmt.Errorf("l2_to_cosmos config: unsupported arbitrum rollup protocol %q (want bold_v2|legacy_nitro)", protocol)
		}
	default:
		l1.Close()
		l2.Close()
		_ = cosmosClient.Stop()
		_ = attestor.Close()
		return nil, nil, fmt.Errorf("l2_to_cosmos config: unsupported chain kind %q", cfg.kind)
	}

	source := l2rollup.NewSource(cfg.kind, l2, headKind, cfg.L2ICS26ClientID, cfg.AttestorSrcChain, router, attestor)
	dest := l2rollup.NewDestination(worker, svcCtx, cfg.L2WasmClientID)
	builder := l2rollup.NewBuilder(headerBuilder)

	module := relay.NewModule(
		fmt.Sprintf("%s->cosmos", cfg.kind),
		cfg.L2ICS26ClientID,
		source, dest, builder,
	)
	logger.Sugar().Infof("l2->cosmos source: %s (attestor=%s src_chain=%s wasm_client=%s head=%s)",
		cfg.kind, cfg.AttestorAddr, cfg.AttestorSrcChain, cfg.L2WasmClientID, cfg.HeadKind)

	cleanup := func() {
		l1.Close()
		l2.Close()
		_ = cosmosClient.Stop()
		_ = attestor.Close()
	}
	return module, cleanup, nil
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

// parseArbProfile extracts the Arbitrum header builder's fields from the rollup
// profile. It mirrors the arbitrum-verifier Profile (packages/arbitrum-verifier): the
// L1 RollupCore (`rollup`), the BoLD `_assertions` mapping slot, the L2 router, and the
// shared ETH client id. assertion_scan_range is an optional relayer-only override.
// arbRollupProtocol reads rollup_profile.protocol.type (bold_v2 | legacy_nitro) — the
// tagged RollupProtocol discriminant that selects the Arbitrum header builder.
func arbRollupProtocol(profile json.RawMessage) (string, error) {
	var p struct {
		Protocol struct {
			Type string `json:"type"`
		} `json:"protocol"`
	}
	if err := json.Unmarshal(profile, &p); err != nil {
		return "", fmt.Errorf("parse Arbitrum rollup_profile.protocol: %w", err)
	}
	if p.Protocol.Type == "" {
		return "", fmt.Errorf("rollup_profile.protocol.type is required (bold_v2|legacy_nitro)")
	}
	return p.Protocol.Type, nil
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

// parseArbBoldProfile extracts the BoLD v2 header builder's fields (protocol.value holds
// assertions_mapping_slot).
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

// parseArbLegacyProfile extracts the legacy Nitro header builder's fields
// (protocol.value holds the node-lifecycle slot, nodes mapping slot, and offsets).
func parseArbLegacyProfile(profile json.RawMessage) (l2rollup.ArbLegacyProfile, error) {
	var p struct {
		Common struct {
			L2Router       string `json:"l2_router"`
			EthereumClient struct {
				ClientID string `json:"client_id"`
			} `json:"ethereum_client"`
		} `json:"common"`
		Rollup   string `json:"rollup"`
		Protocol struct {
			Value struct {
				NodeLifecycleSlot     string `json:"node_lifecycle_slot"`
				NodesMappingSlot      string `json:"nodes_mapping_slot"`
				LatestConfirmedOffset uint8  `json:"latest_confirmed_offset"`
				FirstUnresolvedOffset uint8  `json:"first_unresolved_offset"`
				LatestCreatedOffset   uint8  `json:"latest_created_offset"`
				ConfirmDataOffset     uint8  `json:"confirm_data_offset"`
			} `json:"value"`
		} `json:"protocol"`
	}
	if err := json.Unmarshal(profile, &p); err != nil {
		return l2rollup.ArbLegacyProfile{}, fmt.Errorf("parse Arbitrum legacy_nitro rollup_profile: %w", err)
	}
	rollup, l2Router, clientID, err := arbCommonFields(p.Rollup, p.Common.L2Router, p.Common.EthereumClient.ClientID)
	if err != nil {
		return l2rollup.ArbLegacyProfile{}, err
	}
	v := p.Protocol.Value
	if v.NodeLifecycleSlot == "" || v.NodesMappingSlot == "" {
		return l2rollup.ArbLegacyProfile{}, fmt.Errorf("rollup_profile.protocol.value.{node_lifecycle_slot,nodes_mapping_slot} are required")
	}
	return l2rollup.ArbLegacyProfile{
		RollupCore:            rollup,
		L2Router:              l2Router,
		L1ClientID:            clientID,
		NodeLifecycleSlot:     common.HexToHash(v.NodeLifecycleSlot),
		NodesMappingSlot:      common.HexToHash(v.NodesMappingSlot),
		LatestConfirmedOffset: v.LatestConfirmedOffset,
		FirstUnresolvedOffset: v.FirstUnresolvedOffset,
		LatestCreatedOffset:   v.LatestCreatedOffset,
		ConfirmDataOffset:     v.ConfirmDataOffset,
	}, nil
}

// cosmosEthClientReader adapts the Cosmos RPC client to the header builder's
// cosmosClientStateReader: it reads the shared ETH (08-wasm) client's trusted L1
// slot + execution block from Cosmos.
type cosmosEthClientReader struct {
	cosmos *rpchttp.HTTP
}

func (r *cosmosEthClientReader) EthClientLatestSlotAndBlock(l1ClientID string) (uint64, uint64, error) {
	cs, err := relayerclient.GetEthereumClientState(r.cosmos, l1ClientID)
	if err != nil {
		return 0, 0, err
	}
	return cs.LatestSlot, cs.LatestExecutionBlockNumber, nil
}
