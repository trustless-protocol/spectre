package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	relayerclient "relayer/client"
	"relayer/services"
	"relayer/transaction"
)

const flagL2Config = "l2-config"

// l2ClientConfig is one L2 rollup's client-creation config, passed to
// `create-clients-cosmos --l2-config <file>` (repeatable, one per L2 source). It is
// kept separate from the modules[] schema, which the config-schema follow-up will
// restructure. Cosmos connection details come from the main --config (cosmos_to_eth);
// this file only carries the L2-specific deployment identity. The L1 client id inside
// rollup_profile.common.ethereum_client is filled in by create-clients-cosmos (from
// the L1 client it just created), so it may be left empty here.
type l2ClientConfig struct {
	// WasmChecksum is the hex checksum of the L2 client's wasm code, already
	// governance-stored on Cosmos.
	WasmChecksum string `json:"wasm_checksum"`
	// L2RPCURL is the L2 execution RPC used to read the bootstrap roots.
	L2RPCURL string `json:"l2_rpc_url"`
	// RollupProfile is the full ICS-08 verifier Profile JSON (the `common` block +
	// the rollup-specific finality fields), embedded verbatim as ClientState.profile.
	// The L2 router address, commitment slot, L1 client id/checksum, and chain ids
	// all live inside it (profile.common.*), so they are not repeated here.
	RollupProfile json.RawMessage `json:"rollup_profile"`
	// BootstrapBlock is the L2 block to bootstrap from; 0 (or absent) = latest.
	BootstrapBlock uint64 `json:"bootstrap_block"`
}

// profileCommon is the subset of the verifier Profile the command reads to locate
// the L2 router for the bootstrap eth_getProof (single source of truth: the router
// address is not duplicated in l2-config).
type profileCommon struct {
	Common struct {
		L2Router string `json:"l2_router"`
	} `json:"common"`
}

// routerAddress extracts profile.common.l2_router from the rollup profile.
func (c *l2ClientConfig) routerAddress() (string, error) {
	var pc profileCommon
	if err := json.Unmarshal(c.RollupProfile, &pc); err != nil {
		return "", fmt.Errorf("l2-config: parse rollup_profile.common: %w", err)
	}
	if pc.Common.L2Router == "" {
		return "", fmt.Errorf("l2-config: rollup_profile.common.l2_router is required")
	}
	// Fail fast on malformed input: common.HexToAddress silently truncates/zero-pads a
	// bad string, which would later hit a wrong (or zero) address. GetL2BootstrapState
	// still rejects a non-contract (empty storage root), but a clear error here beats a
	// confusing downstream one.
	if !common.IsHexAddress(pc.Common.L2Router) {
		return "", fmt.Errorf("l2-config: rollup_profile.common.l2_router %q is not a valid hex address", pc.Common.L2Router)
	}
	return pc.Common.L2Router, nil
}

func (c *l2ClientConfig) validate() error {
	if c.WasmChecksum == "" {
		return fmt.Errorf("l2-config: wasm_checksum is required")
	}
	if c.L2RPCURL == "" {
		return fmt.Errorf("l2-config: l2_rpc_url is required")
	}
	if len(c.RollupProfile) == 0 {
		return fmt.Errorf("l2-config: rollup_profile is required")
	}
	if _, err := c.routerAddress(); err != nil {
		return err
	}
	return nil
}

func loadL2ClientConfig(path string) (*l2ClientConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read l2-config %s: %w", path, err)
	}
	var cfg l2ClientConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse l2-config %s: %w", path, err)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// injectL1ClientID sets rollup_profile.common.ethereum_client.client_id to the L1
// (Ethereum) wasm client id that create-clients-cosmos just created on Cosmos. The
// L2 client is anchored to that L1 client (validated by the wasm client at update
// time), and its id is only known after MsgCreateClient lands — so we inject it
// rather than making the operator hand-copy it. Other profile fields are preserved
// verbatim (the profile is otherwise opaque, rollup-specific config).
func injectL1ClientID(profile json.RawMessage, l1ClientID string) (json.RawMessage, error) {
	var root map[string]json.RawMessage
	if err := json.Unmarshal(profile, &root); err != nil {
		return nil, fmt.Errorf("parse rollup_profile: %w", err)
	}
	commonRaw, ok := root["common"]
	if !ok {
		return nil, fmt.Errorf("rollup_profile.common is missing")
	}
	var common map[string]json.RawMessage
	if err := json.Unmarshal(commonRaw, &common); err != nil {
		return nil, fmt.Errorf("parse rollup_profile.common: %w", err)
	}
	ethRaw, ok := common["ethereum_client"]
	if !ok {
		return nil, fmt.Errorf("rollup_profile.common.ethereum_client is missing")
	}
	var eth map[string]json.RawMessage
	if err := json.Unmarshal(ethRaw, &eth); err != nil {
		return nil, fmt.Errorf("parse rollup_profile.common.ethereum_client: %w", err)
	}
	idBz, err := json.Marshal(l1ClientID)
	if err != nil {
		return nil, err
	}
	eth["client_id"] = idBz
	if common["ethereum_client"], err = json.Marshal(eth); err != nil {
		return nil, err
	}
	if root["common"], err = json.Marshal(common); err != nil {
		return nil, err
	}
	return json.Marshal(root)
}

// runCreateClientsL2 bootstraps one L2 rollup wasm light client on Cosmos: it reads
// the initial state roots from the L2 chain, wraps them in the ICS-08 client state
// (ClientState<Profile>) + consensus state, and submits MsgCreateClient. The
// initial consensus state (state roots read here) is the client's trust anchor at
// creation; ongoing updates are verified trustlessly by the wasm client against the
// pinned L1 (Ethereum) client + rollup proofs. l1ClientID (when non-empty) is
// injected into the rollup profile's ethereum_client.client_id.
func runCreateClientsL2(logger *zap.Logger, cfg *appConfig, l2cfg *l2ClientConfig, l1ClientID string) error {
	// Cosmos context (MsgCreateClient is submitted to Cosmos).
	ctx, cosmosClient, err := buildCreateClientsContext(logger, cfg, "")
	if err != nil {
		return err
	}
	defer cosmosClient.Stop()

	// The L2 client's wasm code must already be governance-stored on Cosmos.
	logger.Sugar().Infof("create-clients-cosmos[l2]: validating wasm checksum on Cosmos: %s", l2cfg.WasmChecksum)
	ok, err := cosmosHasWasmChecksum(cosmosClient, l2cfg.WasmChecksum)
	if err != nil {
		return fmt.Errorf("validate wasm checksum on Cosmos: %w", err)
	}
	if !ok {
		return fmt.Errorf("wasm checksum %s has not been stored on Cosmos", l2cfg.WasmChecksum)
	}

	// Read the trusted bootstrap roots from the L2 chain.
	logger.Sugar().Infof("create-clients-cosmos[l2]: dialing L2 rpc %s", l2cfg.L2RPCURL)
	l2Client, err := ethclient.Dial(l2cfg.L2RPCURL)
	if err != nil {
		return fmt.Errorf("dial L2 rpc: %w", err)
	}
	defer l2Client.Close()

	router, err := l2cfg.routerAddress() // profile.common.l2_router — single source of truth
	if err != nil {
		return err
	}
	var block *big.Int
	if l2cfg.BootstrapBlock != 0 {
		block = new(big.Int).SetUint64(l2cfg.BootstrapBlock)
	}
	bootstrap, err := relayerclient.GetL2BootstrapState(l2Client, common.HexToAddress(router), block)
	if err != nil {
		return err
	}
	logger.Sugar().Infof("create-clients-cosmos[l2]: bootstrap at L2 block %d (state_root=%s router_storage_root=%s ts=%d)",
		bootstrap.Height, bootstrap.StateRoot.Hex(), bootstrap.RouterStorageRoot.Hex(), bootstrap.TimestampSeconds)

	profile := l2cfg.RollupProfile
	if l1ClientID != "" {
		profile, err = injectL1ClientID(profile, l1ClientID)
		if err != nil {
			return fmt.Errorf("inject L1 client id: %w", err)
		}
		logger.Sugar().Infof("create-clients-cosmos[l2]: anchored to L1 client id %s (injected into rollup_profile.common.ethereum_client)", l1ClientID)
	}

	worker := services.NewWorker(&transaction.Handler{}, nil)
	clientID, err := worker.CreateL2Client(context.Background(), ctx, services.L2ClientParams{
		WasmChecksum:  l2cfg.WasmChecksum,
		RollupProfile: profile,
		Bootstrap:     bootstrap,
	})
	if err != nil {
		return fmt.Errorf("create L2 client on Cosmos: %w", err)
	}
	logger.Sugar().Infof("create-clients-cosmos[l2]: created L2 wasm client on Cosmos: clientID=%s", clientID)
	fmt.Println(clientID)
	return nil
}
