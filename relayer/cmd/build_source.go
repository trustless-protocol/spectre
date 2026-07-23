package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	"relayer/prover"
	"relayer/services"
)

// startCosmosToEthSource wires up one Cosmos→ETH source: its own Tendermint RPC,
// SpectreClient and router client id, sharing the passed-in prover and the ETH
// beacon endpoint from the single eth_to_cosmos module. It returns the built
// Services, its Context, and a cleanup func that stops the chain clients.
//
// The caller runs runAdapterEngine(ctx, svc, ctx) — which blocks until an error
// or shutdown — on its own goroutine, so N sources relay independently. Their ETH
// event streams don't cross-feed: SubscribeEth filters ICS26Router logs by the
// per-source router client id.
//
// allowEnvOverride honors the single-source env overrides (ICS26_CLIENT_ID,
// COSMOS_WASM_CLIENT_ID, ROLE_MANAGER); it must be false when more than one
// source is configured, or the overrides would apply to every source.
func startCosmosToEthSource(
	logger *zap.Logger,
	c2e cosmosToEthConfig,
	e2c ethToCosmosConfig,
	batchCfg services.BatchConfig,
	p *prover.EcipProver,
	txHandler services.TransactionHandler,
	allowEnvOverride bool,
) (*services.Services, services.Context, func(), error) {
	var zero services.Context

	// Connect to Ethereum (HTTP for queries)
	ethClient, err := ethclient.Dial(c2e.EthRpcUrl)
	if err != nil {
		return nil, zero, nil, fmt.Errorf("failed to connect to Ethereum: %w", err)
	}

	// Connect to Ethereum (WS for subscriptions)
	var ethWsClient *ethclient.Client
	if c2e.EthWsUrl != "" {
		if !strings.HasPrefix(c2e.EthWsUrl, "ws://") && !strings.HasPrefix(c2e.EthWsUrl, "wss://") {
			return nil, zero, nil, fmt.Errorf("eth_ws_url must use ws:// or wss://, got: %s", c2e.EthWsUrl)
		}
		ethWsClient, err = ethclient.Dial(c2e.EthWsUrl)
		if err != nil {
			return nil, zero, nil, fmt.Errorf("failed to connect to Ethereum WS: %w", err)
		}
	}

	// Connect to Cosmos
	cosmosClient, err := rpchttp.New(c2e.TmRpcUrl, "/websocket")
	if err != nil {
		return nil, zero, nil, fmt.Errorf("failed to create Cosmos RPC client: %w", err)
	}

	cosmosWasmClientID := c2e.CosmosWasmClientID
	if allowEnvOverride {
		cosmosWasmClientID = envOrDefault("COSMOS_WASM_CLIENT_ID", cosmosWasmClientID)
	}
	if cosmosWasmClientID == "" {
		return nil, zero, nil, fmt.Errorf("cosmos_wasm_client_id is required in cosmos_to_eth config")
	}

	// Create context with beacon API
	ctx := services.NewCtxWithBeacon(
		cosmosClient, ethClient, ethWsClient,
		c2e.EthWsUrl, e2c.BeaconUrl, cosmosWasmClientID,
	)

	cosmosRouterClientID := c2e.ICS26ClientID
	if allowEnvOverride {
		cosmosRouterClientID = envOrDefault("ICS26_CLIENT_ID", cosmosRouterClientID)
	}
	if cosmosRouterClientID == "" {
		return nil, zero, nil, fmt.Errorf("cosmos router client ID (ICS26_CLIENT_ID or ics26_client_id) is required and cannot be empty")
	}
	ctx.SetCosmosRouterClientID(cosmosRouterClientID)

	// Set contract addresses from config
	roleManager := c2e.ICS26Address
	if allowEnvOverride {
		roleManager = envOrDefault("ROLE_MANAGER", roleManager)
	}
	ctx.SetAddresses(
		c2e.ICS26Address, c2e.SignatureVerifier, c2e.Membership,
		c2e.Misbehaviour, c2e.UpdateClient, roleManager,
	)

	// Set SpectreClient address (already deployed)
	if c2e.SpectreClient == "" {
		return nil, zero, nil, fmt.Errorf("spectre_client address is required in cosmos_to_eth config")
	}
	ctx.SetClient(common.HexToAddress(c2e.SpectreClient))

	// Start Cosmos WebSocket client
	if err := cosmosClient.Start(); err != nil {
		return nil, zero, nil, fmt.Errorf("failed to start Cosmos WS client: %w", err)
	}
	cleanup := func() { ctx.StopClient() }

	cosmosConfig := buildCosmosConfig(c2e, batchCfg)
	ctx.Config = cosmosConfig

	logger.Sugar().Infof("source %q: subscribing to events (spectre_client=%s tm=%s)",
		cosmosRouterClientID, c2e.SpectreClient, c2e.TmRpcUrl)

	svc := services.New(txHandler, p, cosmosConfig)
	return svc, ctx, cleanup, nil
}

// buildCosmosConfig layers the per-source overrides from c2e (and the global
// FETCH_TIMEOUT env) onto the service defaults.
func buildCosmosConfig(c2e cosmosToEthConfig, batchCfg services.BatchConfig) services.Config {
	cfg := services.DefaultConfig()
	if c2e.TrustingPeriod != 0 {
		cfg.TrustingPeriod = c2e.TrustingPeriod
	}
	if c2e.TrustLevel != "" {
		cfg.TrustLevel = c2e.TrustLevel
	}
	if c2e.ProofType != "" {
		cfg.ProofType = c2e.ProofType
	}
	if c2e.ClockDrift != 0 {
		cfg.ClockDrift = c2e.ClockDrift
	}
	if c2e.BeaconFinalityRetries != 0 {
		cfg.BeaconFinalityRetries = c2e.BeaconFinalityRetries
	}
	if c2e.AppHashWaitRetries != 0 {
		cfg.AppHashWaitRetries = c2e.AppHashWaitRetries
	}
	if c2e.AppHashWaitInterval != 0 {
		cfg.AppHashWaitInterval = time.Duration(c2e.AppHashWaitInterval) * time.Second
	}
	if c2e.FetchTimeout != 0 {
		cfg.FetchTimeout = time.Duration(c2e.FetchTimeout) * time.Second
	}
	if c2e.RotationThreshold != "" {
		cfg.RotationThreshold = c2e.RotationThreshold
	}
	if c2e.RefreshInterval != 0 {
		cfg.RefreshInterval = time.Duration(c2e.RefreshInterval) * time.Second
		cfg.RefreshIntervalConfigured = true
	}
	if envVal := os.Getenv("FETCH_TIMEOUT"); envVal != "" {
		if d, err := strconv.Atoi(envVal); err == nil && d > 0 {
			cfg.FetchTimeout = time.Duration(d) * time.Second
		}
	}
	cfg.BatchConfig = batchCfg
	return cfg
}
