package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/common"
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
		"attestor_addr": c.AttestorAddr, "attestor_src_chain": c.AttestorSrcChain,
		"l2_ics26_client_id": c.L2ICS26ClientID,
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
	headKind, err := parseHeadKind(c.HeadKind)
	if err != nil {
		return err
	}
	if _, err := l2RouterFromProfile(c.RollupProfile); err != nil {
		return err
	}
	verifier, err := l2AttestationVerifierFromProfile(c.RollupProfile)
	if err != nil {
		return err
	}
	if !verifier.MatchesRunMode(headKind.RunMode()) {
		return fmt.Errorf("l2_to_cosmos config: head_kind %q must match rollup_profile.common.attestation_head", c.HeadKind)
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

// l2AttestationVerifierFromProfile parses the immutable key the wasm client
// will use. Verifying the attestor response here makes an operator key mismatch
// visible before submitting a transaction; it is not a substitute for the
// independent on-chain check.
func l2AttestationVerifierFromProfile(profile json.RawMessage) (l2rollup.AttestationVerifier, error) {
	var pc struct {
		Common struct {
			L2ChainID         uint64 `json:"l2_chain_id"`
			AttestorPublicKey string `json:"attestor_public_key"`
			AttestationHead   string `json:"attestation_head"`
		} `json:"common"`
	}
	if err := json.Unmarshal(profile, &pc); err != nil {
		return l2rollup.AttestationVerifier{}, fmt.Errorf("l2_to_cosmos config: parse rollup_profile.common: %w", err)
	}
	keyHex := pc.Common.AttestorPublicKey
	if len(keyHex) >= 2 && keyHex[:2] == "0x" {
		keyHex = keyHex[2:]
	}
	publicKey, err := hex.DecodeString(keyHex)
	if err != nil {
		return l2rollup.AttestationVerifier{}, fmt.Errorf("l2_to_cosmos config: decode rollup_profile.common.attestor_public_key: %w", err)
	}
	verifier, err := l2rollup.NewAttestationVerifier(pc.Common.L2ChainID, pc.Common.AttestationHead, publicKey)
	if err != nil {
		return l2rollup.AttestationVerifier{}, fmt.Errorf("l2_to_cosmos config: invalid attestor profile: %w", err)
	}
	return verifier, nil
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
	verifier, err := l2AttestationVerifierFromProfile(cfg.RollupProfile)
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
	attestor, err := attestorgrpc.Dial(cfg.AttestorAddr)
	if err != nil {
		l2.Close()
		if stopErr := cosmosClient.Stop(); stopErr != nil {
			log.Printf("failed to terminate cosmos client: %v", stopErr)
		}
		return nil, nil, fmt.Errorf("dial attestor: %w", err)
	}

	// Cosmos-side context for the destination (it uses only CosmosClient + the
	// signer via the tx handler); the L2 exec client fills the eth slot unused here.
	svcCtx := services.CosmosEndpoint{Client: cosmosClient}
	worker := services.NewWorker(txHandler, nil)

	// One builder for every chain: the header is the canonical L2 block plus a router
	// account proof, and nothing about that is chain-specific any more. It takes the
	// same attestor the source gates heights on, so the block it packages is bound to
	// what that attestor's replica actually has at that height.
	headerBuilder := l2rollup.NewAttestedHeaderBuilder(l2, router, attestor, verifier, cfg.AttestorSrcChain, headKind.RunMode(), fmt.Sprintf("l2-%s", cfg.kind))

	trackL2Pending, untrackL2Pending := l2PendingTrackerHooks(timeoutReturn.svc)
	source := l2rollup.NewSource(cfg.kind, l2, headKind, cfg.L2ICS26ClientID, cfg.L2WasmClientID, cfg.AttestorSrcChain, router, attestor, includeProvisional).
		WithLogScanChunk(cfg.LogScanChunk).
		// AckPacket / TimeoutPacket on the L2 close a packet's lifecycle. Reading
		// them lets the tracker drop a packet another relayer settled without
		// waiting for a timeout scan to query for it. Same hook the module uses to
		// untrack after its own successful relay -- one settlement path, two ways
		// of learning about it.
		WithSettleHook(untrackL2Pending)
	dest := l2rollup.NewDestination(worker, svcCtx, cfg.L2WasmClientID)
	builder := l2rollup.NewBuilder(headerBuilder)

	module := relay.NewModule(
		fmt.Sprintf("%s->cosmos", cfg.kind),
		cfg.L2ICS26ClientID,
		source, dest, builder,
		relay.WithTimeoutScanner(0, func(c context.Context) {
			timeoutReturn.svc.ScanL2Timeouts(c, timeoutReturn.deps.Cosmos, timeoutReturn.deps.EVM, timeoutReturn.deps.IDs.CosmosOnEVM)
		}),
		relay.WithPacketTracker(trackL2Pending, untrackL2Pending),
	)
	logger.Sugar().Infof("l2->cosmos source: %s (attestor=%s src_chain=%s wasm_client=%s head=%s)",
		cfg.kind, cfg.AttestorAddr, cfg.AttestorSrcChain, cfg.L2WasmClientID, cfg.HeadKind)

	cleanup := func() {
		l2.Close()
		if err := cosmosClient.Stop(); err != nil {
			log.Printf("failed to terminate cosmos client: %v", err)
		}
		if err := attestor.Close(); err != nil {
			log.Printf("failed to close attestor client: %v", err)
		}
	}
	return module, cleanup, nil
}

func l2PendingTrackerHooks(svc *services.Services) (relay.TrackFunc, relay.UntrackFunc) {
	trackL2Pending := func(raw []byte, height uint64) bool {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[adapter l2->cosmos] track pending: decode packet: %v", err)
			return false
		}
		return svc.TrackL2Pending(pkt, height)
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
	return module.Run(ctx)
}
