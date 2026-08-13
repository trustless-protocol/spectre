package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	relayerclient "relayer/client"
	"relayer/prover"
	"relayer/services"
)

// buildCosmosToEthSource builds one Cosmos→ETH source: its own Tendermint RPC,
// SpectreClient and router client id, sharing the passed-in prover and the ETH
// beacon endpoint from the single eth_to_cosmos module. It returns the built
// Services, the composition-root dependencies, and a cleanup func that stops
// the chain clients.
//
// The caller runs runAdapterEngine(ctx, svc, deps) — which blocks until an error
// or shutdown — on its own goroutine, so N sources relay independently. Their ETH
// event streams don't cross-feed: SubscribeEth filters ICS26Router logs by the
// per-source router client id.
//
// buildCosmosToEthSourceOptions separates the long-running subscription setup
// from one-shot operator commands. Incident commands use only HTTP RPCs and
// must not depend on websocket infrastructure they never consume.
type buildCosmosToEthSourceOptions struct {
	allowEnvOverride   bool
	startSubscriptions bool
}

func buildCosmosToEthSource(
	logger *zap.Logger,
	c2e cosmosToEthConfig,
	e2c ethToCosmosConfig,
	batchCfg services.BatchConfig,
	p *prover.EcipProver,
	txHandler services.TransactionHandler,
	options buildCosmosToEthSourceOptions,
) (*services.Services, services.RelayDeps, func(), error) {
	var zero services.RelayDeps

	// Connect to Ethereum (HTTP for queries)
	ethClient, err := relayerclient.DialEthRPC(context.Background(), c2e.EthRpcUrl, relayerclient.DefaultRPCTimeout)
	if err != nil {
		return nil, zero, nil, fmt.Errorf("failed to connect to Ethereum: %w", err)
	}

	// Connect to Ethereum (WS for subscriptions)
	var ethWsClient *ethclient.Client
	if options.startSubscriptions && c2e.EthWsUrl != "" {
		if !strings.HasPrefix(c2e.EthWsUrl, "ws://") && !strings.HasPrefix(c2e.EthWsUrl, "wss://") {
			return nil, zero, nil, fmt.Errorf("eth_ws_url must use ws:// or wss://, got: %s", c2e.EthWsUrl)
		}
		ethWsClient, err = relayerclient.DialEthRPC(context.Background(), c2e.EthWsUrl, relayerclient.DefaultRPCTimeout)
		if err != nil {
			return nil, zero, nil, fmt.Errorf("failed to connect to Ethereum WS: %w", err)
		}
	}

	// Connect to Cosmos
	cosmosClient, err := relayerclient.DialCosmosRPC(c2e.TmRpcUrl, "/websocket", relayerclient.DefaultRPCTimeout)
	if err != nil {
		return nil, zero, nil, fmt.Errorf("failed to create Cosmos RPC client: %w", err)
	}

	cosmosWasmClientID := c2e.CosmosWasmClientID
	if options.allowEnvOverride {
		cosmosWasmClientID = envOrDefault("COSMOS_WASM_CLIENT_ID", cosmosWasmClientID)
	}
	if cosmosWasmClientID == "" {
		return nil, zero, nil, fmt.Errorf("cosmos_wasm_client_id is required in cosmos_to_eth config")
	}

	// Assemble the scoped chain dependencies, including the beacon API.
	cosmosRouterClientID := c2e.ICS26ClientID
	if options.allowEnvOverride {
		cosmosRouterClientID = envOrDefault("ICS26_CLIENT_ID", cosmosRouterClientID)
	}
	if cosmosRouterClientID == "" {
		return nil, zero, nil, fmt.Errorf("cosmos router client ID (ICS26_CLIENT_ID or ics26_client_id) is required and cannot be empty")
	}
	// Set contract addresses from config
	roleManager := c2e.ICS26Address
	if options.allowEnvOverride {
		roleManager = envOrDefault("ROLE_MANAGER", roleManager)
	}
	// Set SpectreClient address (already deployed)
	if c2e.SpectreClient == "" {
		return nil, zero, nil, fmt.Errorf("spectre_client address is required in cosmos_to_eth config")
	}
	if options.startSubscriptions {
		if err := cosmosClient.Start(); err != nil {
			return nil, zero, nil, fmt.Errorf("failed to start Cosmos WS client: %w", err)
		}
	}
	cosmosConfig := buildCosmosConfig(c2e, batchCfg)
	deps := services.RelayDeps{
		Cosmos: services.CosmosEndpoint{Client: cosmosClient},
		EVM: services.EVMEndpoint{
			Client: ethClient, WSURL: c2e.EthWsUrl, BeaconAPIURL: e2c.BeaconUrl,
			Contracts: services.EVMContracts{
				Router: common.HexToAddress(c2e.ICS26Address), SignatureVerifier: common.HexToAddress(c2e.SignatureVerifier),
				Membership: common.HexToAddress(c2e.Membership), Misbehaviour: common.HexToAddress(c2e.Misbehaviour),
				UpdateClient: common.HexToAddress(c2e.UpdateClient), RoleManager: common.HexToAddress(roleManager),
				SpectreClient: common.HexToAddress(c2e.SpectreClient),
			},
		},
		IDs:    services.ClientIDs{CosmosOnEVM: cosmosRouterClientID, EVMOnCosmos: cosmosWasmClientID},
		Config: cosmosConfig, Logger: log.Default(),
	}
	cleanup := func() {
		if options.startSubscriptions {
			if err := cosmosClient.Stop(); err != nil {
				log.Printf("failed to terminate cosmos client: %v", err)
			}
		}
		ethClient.Close()
		if ethWsClient != nil {
			ethWsClient.Close()
		}
	}

	if options.startSubscriptions {
		logger.Sugar().Infof("source %q: subscribing to events (spectre_client=%s tm=%s)",
			cosmosRouterClientID, c2e.SpectreClient, c2e.TmRpcUrl)
	} else {
		logger.Sugar().Infof("source %q: configured one-shot RPC context (spectre_client=%s tm=%s)",
			cosmosRouterClientID, c2e.SpectreClient, c2e.TmRpcUrl)
	}

	svc := services.New(txHandler, p, cosmosConfig)
	return svc, deps, cleanup, nil
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
