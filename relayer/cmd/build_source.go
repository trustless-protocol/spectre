package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
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
	pendingStateDir    string

	// relaysEVMToCosmos says this caller will drive the eth->cosmos direction,
	// which is the only consumer of cosmos_wasm_client_id -- the id of the
	// Ethereum light client ON Cosmos.
	//
	// It is a property of the CALLER, not of the config, and treating it as the
	// latter broke submit-misbehaviour for a cosmos_to_l2 source: selectSource
	// offers L2 sources, cosmos_to_l2 deliberately carries no
	// cosmos_wasm_client_id, and the builder rejected it for a field the command
	// never reads. The evidence was refused before it was looked at.
	//
	// Default false, so a caller opts IN to the requirement. That is the safe
	// direction here: a caller that forgets it fails at the point of use with the
	// id it wanted, whereas defaulting to true would reject valid L2
	// configurations again -- the failure this field exists to remove.
	relaysEVMToCosmos bool
}

// pendingStateDir keeps restart state beside the selected config while
// separating relay kinds and client ids. The digest avoids treating client ids
// as paths and also scopes state to the config file when several deployments
// run from the same directory.
func pendingStateDir(configPath, relayKind, clientID string) string {
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		absConfigPath = configPath
	}
	sum := sha256.Sum256([]byte(absConfigPath + "\x00" + relayKind + "\x00" + clientID))
	return filepath.Join(filepath.Dir(configPath), ".fast-ibc-state", fmt.Sprintf("%s-%x", relayKind, sum[:8]))
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

	// Everything that can be decided from config alone is decided BEFORE any
	// client is dialled or started. A validation failure below this block has no
	// resource to release; one above it would have to unwind a dialled ETH client
	// and a started Cosmos WS subscription, and the cleanup closure that knows how
	// to do that is not built until the end of this function.
	cosmosWasmClientID := c2e.CosmosWasmClientID
	if options.allowEnvOverride {
		cosmosWasmClientID = envOrDefault("COSMOS_WASM_CLIENT_ID", cosmosWasmClientID)
	}
	if options.relaysEVMToCosmos && cosmosWasmClientID == "" {
		return nil, zero, nil, fmt.Errorf(
			"cosmos_wasm_client_id is required in cosmos_to_eth config to relay eth->cosmos")
	}
	cosmosRouterClientID := c2e.ICS26ClientID
	if options.allowEnvOverride {
		cosmosRouterClientID = envOrDefault("ICS26_CLIENT_ID", cosmosRouterClientID)
	}
	if cosmosRouterClientID == "" {
		return nil, zero, nil, fmt.Errorf("cosmos router client ID (ICS26_CLIENT_ID or ics26_client_id) is required and cannot be empty")
	}
	roleManager := c2e.ICS26Address
	if options.allowEnvOverride {
		roleManager = envOrDefault("ROLE_MANAGER", roleManager)
	}
	// SpectreClient is already deployed; this is its address, not a request to deploy.
	if c2e.SpectreClient == "" {
		return nil, zero, nil, fmt.Errorf("spectre_client address is required in cosmos_to_eth config")
	}
	if options.startSubscriptions && c2e.EthWsUrl != "" {
		if !strings.HasPrefix(c2e.EthWsUrl, "ws://") && !strings.HasPrefix(c2e.EthWsUrl, "wss://") {
			return nil, zero, nil, fmt.Errorf("eth_ws_url must use ws:// or wss://, got: %s", c2e.EthWsUrl)
		}
	}
	// Parses FETCH_TIMEOUT, so it can fail on a malformed override.
	cosmosConfig, err := buildCosmosConfig(c2e, batchCfg)
	if err != nil {
		return nil, zero, nil, err
	}

	// Connect to Ethereum (HTTP for queries)
	ethClient, err := relayerclient.DialEthRPC(context.Background(), c2e.EthRpcUrl, relayerclient.DefaultRPCTimeout)
	if err != nil {
		return nil, zero, nil, fmt.Errorf("failed to connect to Ethereum: %w", err)
	}

	// Connect to Ethereum (WS for subscriptions)
	var ethWsClient *ethclient.Client
	if options.startSubscriptions && c2e.EthWsUrl != "" {
		ethWsClient, err = relayerclient.DialEthRPC(context.Background(), c2e.EthWsUrl, relayerclient.DefaultRPCTimeout)
		if err != nil {
			ethClient.Close()
			return nil, zero, nil, fmt.Errorf("failed to connect to Ethereum WS: %w", err)
		}
	}

	// Connect to Cosmos
	cosmosClient, err := relayerclient.DialCosmosRPC(c2e.TmRpcUrl, "/websocket", relayerclient.DefaultRPCTimeout)
	if err != nil {
		ethClient.Close()
		if ethWsClient != nil {
			ethWsClient.Close()
		}
		return nil, zero, nil, fmt.Errorf("failed to create Cosmos RPC client: %w", err)
	}

	if options.startSubscriptions {
		if err := cosmosClient.Start(); err != nil {
			ethClient.Close()
			if ethWsClient != nil {
				ethWsClient.Close()
			}
			return nil, zero, nil, fmt.Errorf("failed to start Cosmos WS client: %w", err)
		}
	}

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

	var recoveryState *services.RecoveryStateStore
	if options.startSubscriptions {
		recoveryState, err = services.LoadRecoveryState(services.RecoveryStatePath())
		if err != nil {
			cleanup()
			return nil, zero, nil, fmt.Errorf("load recovery state: %w", err)
		}
	}
	var svc *services.Services
	if options.pendingStateDir != "" {
		svc, err = services.NewWithPendingState(txHandler, p, cosmosConfig, options.pendingStateDir, recoveryState)
		if err != nil {
			cleanup()
			return nil, zero, nil, err
		}
	} else {
		svc = services.New(txHandler, p, cosmosConfig, recoveryState)
	}
	return svc, deps, cleanup, nil
}

// buildCosmosConfig layers the per-source overrides from c2e (and the global
// FETCH_TIMEOUT env) onto the service defaults.
// envSecondsDuration reads an optional whole-second duration override.
//
// It returns an error rather than falling back, because falling back is what the
// old code did: `if d, err := strconv.Atoi(envVal); err == nil && d > 0` treated
// BOTH a parse failure and a non-positive value as "not set". FETCH_TIMEOUT=30s
// -- the way a Go duration is normally written, and the natural guess for a
// field that is a time.Duration -- was silently discarded, and the operator got
// the default while believing they had set a deadline.
//
// That matters more here than for a gas knob: FetchTimeout is what bounds the
// Cosmos RPC calls, and an RPC without the deadline the operator asked for is
// the failure mode that wedged a whole relay direction before.
func envSecondsDuration(name string) (time.Duration, bool, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return 0, false, nil
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false, fmt.Errorf("invalid %s: %w (whole seconds, e.g. 30)", name, err)
	}
	if seconds <= 0 {
		return 0, false, fmt.Errorf("invalid %s: %d seconds is not a usable timeout", name, seconds)
	}
	d, err := secondsToDuration(name, uint64(seconds))
	if err != nil {
		return 0, false, err
	}
	return d, true, nil
}

// maxDurationSeconds is the largest whole-second value a time.Duration can hold.
// time.Duration is an int64 of NANOseconds, so the usable range is roughly 292
// years and anything past it wraps.
const maxDurationSeconds = uint64(math.MaxInt64 / int64(time.Second))

// secondsToDuration converts operator-supplied whole seconds, refusing the ones
// that do not fit.
//
// The multiplication is the trap: strconv accepts 9223372036854775807 happily
// and `time.Duration(seconds) * time.Second` then wraps to -1s. A negative
// timeout is not merely wrong, it is wrong in two different ways depending on
// who reads it -- context.WithTimeout cancels immediately, while
// services.fetchCtx substitutes its 15s default. One accepted setting, two
// behaviours, neither the one the operator asked for.
func secondsToDuration(name string, seconds uint64) (time.Duration, error) {
	if seconds > maxDurationSeconds {
		return 0, fmt.Errorf("invalid %s: %d seconds exceeds the maximum representable duration (%d seconds)",
			name, seconds, maxDurationSeconds)
	}
	// #nosec G115 -- bounded by maxDurationSeconds immediately above
	return time.Duration(seconds) * time.Second, nil
}

func buildCosmosConfig(c2e cosmosToEthConfig, batchCfg services.BatchConfig) (services.Config, error) {
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
	// Same overflow trap as the environment override below, same origin: an
	// operator typing a number. Config JSON is no safer than an env var.
	if c2e.AppHashWaitInterval != 0 {
		d, err := secondsToDuration("app_hash_wait_interval_seconds", c2e.AppHashWaitInterval)
		if err != nil {
			return services.Config{}, err
		}
		cfg.AppHashWaitInterval = d
	}
	if c2e.FetchTimeout != 0 {
		d, err := secondsToDuration("fetch_timeout", c2e.FetchTimeout)
		if err != nil {
			return services.Config{}, err
		}
		cfg.FetchTimeout = d
	}
	if c2e.RotationThreshold != "" {
		cfg.RotationThreshold = c2e.RotationThreshold
	}
	if c2e.RefreshInterval != 0 {
		d, err := secondsToDuration("refresh_interval_seconds", c2e.RefreshInterval)
		if err != nil {
			return services.Config{}, err
		}
		cfg.RefreshInterval = d
		cfg.RefreshIntervalConfigured = true
	}
	timeout, set, err := envSecondsDuration("FETCH_TIMEOUT")
	if err != nil {
		return services.Config{}, err
	}
	if set {
		cfg.FetchTimeout = timeout
	}
	cfg.BatchConfig = batchCfg
	return cfg, nil
}
