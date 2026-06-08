package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	proto "github.com/cosmos/gogoproto/proto"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	utils "relayer/utils"
	tendermintClient "relayer/client"
	"relayer/keys"
	"relayer/prover"
	"relayer/services"
	"relayer/subscriber"
	"relayer/transaction"
)

const (
	flagConfigPath     = "config"
	flagGPUProve       = "gpu-prove"
	flagOnlyOnce       = "only-once"
	flagProofType      = "proof-type"
	flagOutput         = "output"
	flagOutputPath     = "output-path"
	flagTrustLevel     = "trust-level"
	flagTrustingPeriod = "trusting-period"
	flagTrustedBlock   = "trusted-block"
	flagWasmChecksum   = "wasm-checksum"
	flagBenchmark      = "benchmark"
)

// --- Config types for JSON config file ---

type cosmosToEthConfig struct {
	TmRpcUrl           string `json:"tm_rpc_url"`
	ICS26Address       string `json:"ics26_address"`
	ICS26ClientID      string `json:"ics26_client_id"`
	CosmosWasmClientID string `json:"cosmos_wasm_client_id"`
	EthRpcUrl          string `json:"eth_rpc_url"`
	EthWsUrl           string `json:"eth_ws_url"`
	ICS07Client        string `json:"ics07_client"`
	WrapperVerifier    string `json:"wrapper_verifier"`
	Membership         string `json:"membership"`
	Misbehaviour       string `json:"misbehaviour"`
	UpdateClient       string `json:"update_client"`
	TrustingPeriod     uint32 `json:"trusting_period"`
	TrustLevel         string `json:"trust_level"`
	ProofType          string `json:"proof_type"`
}

type ethToCosmosConfig struct {
	TmRpcUrl      string `json:"tm_rpc_url"`
	ICS26Address  string `json:"ics26_address"`
	EthRpcUrl     string `json:"eth_rpc_url"`
	BeaconUrl     string `json:"eth_beacon_api_url"`
	SignerAddress string `json:"signer_address"`
}

type configModule struct {
	Name     string          `json:"name"`
	SrcChain string          `json:"src_chain"`
	DstChain string          `json:"dst_chain"`
	Config   json.RawMessage `json:"config"`
}

type serverConfig struct {
	LogLevel string `json:"log_level"`
	Address  string `json:"address"`
	Port     uint64 `json:"port"`
}

type jsonConfig struct {
	Server  serverConfig   `json:"server"`
	Modules []configModule `json:"modules"`
}

type appConfig struct {
	Server            serverConfig
	CosmosToEthConfig cosmosToEthConfig
	EthToCosmosConfig ethToCosmosConfig
}

// writeICS07Address rewrites configPath in place, setting
// modules[name=="cosmos_to_eth"].config.ics07_client = addr. Other fields and
// JSON formatting are preserved as much as encoding/json indent allows.
func writeICS07Address(configPath, addr string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	modules, ok := raw["modules"].([]any)
	if !ok {
		return fmt.Errorf("config has no modules array")
	}
	updated := false
	for _, m := range modules {
		mod, ok := m.(map[string]any)
		if !ok || mod["name"] != "cosmos_to_eth" {
			continue
		}
		cfg, ok := mod["config"].(map[string]any)
		if !ok {
			return fmt.Errorf("cosmos_to_eth.config is not an object")
		}
		cfg["ics07_client"] = addr
		updated = true
		break
	}
	if !updated {
		return fmt.Errorf("module cosmos_to_eth not found in config")
	}
	out, err := json.MarshalIndent(raw, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, out, 0o644)
}

func loadConfig(configPath string) (*appConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var jc jsonConfig
	if err := json.Unmarshal(data, &jc); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	var c2e cosmosToEthConfig
	var e2c ethToCosmosConfig
	for _, m := range jc.Modules {
		switch m.Name {
		case "cosmos_to_eth":
			if err := json.Unmarshal(m.Config, &c2e); err != nil {
				return nil, fmt.Errorf("failed to parse cosmos_to_eth config: %w", err)
			}
			if c2e.ICS26ClientID == "" {
				c2e.ICS26ClientID = m.SrcChain
			}
		case "eth_to_cosmos":
			if err := json.Unmarshal(m.Config, &e2c); err != nil {
				return nil, fmt.Errorf("failed to parse eth_to_cosmos config: %w", err)
			}
		}
	}

	return &appConfig{
		Server:            jc.Server,
		CosmosToEthConfig: c2e,
		EthToCosmosConfig: e2c,
	}, nil
}

func preflightCreateClients(cfg *appConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ethClient, err := ethclient.DialContext(ctx, cfg.CosmosToEthConfig.EthRpcUrl)
	if err != nil {
		return fmt.Errorf("ethereum rpc unavailable at %s: %w", cfg.CosmosToEthConfig.EthRpcUrl, err)
	}
	defer ethClient.Close()

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("ethereum rpc not responding at %s: %w", cfg.CosmosToEthConfig.EthRpcUrl, err)
	}

	cosmosClient, err := rpchttp.New(cfg.CosmosToEthConfig.TmRpcUrl, "/websocket")
	if err != nil {
		return fmt.Errorf("failed to create cosmos rpc client for %s: %w", cfg.CosmosToEthConfig.TmRpcUrl, err)
	}
	status, err := cosmosClient.Status(ctx)
	if err != nil {
		return fmt.Errorf("cosmos rpc unavailable at %s: %w", cfg.CosmosToEthConfig.TmRpcUrl, err)
	}

	if cfg.EthToCosmosConfig.BeaconUrl != "" {
		if _, err := tendermintClient.GetBeaconGenesis(ctx, cfg.EthToCosmosConfig.BeaconUrl); err != nil {
			return fmt.Errorf("beacon api unavailable at %s: %w", cfg.EthToCosmosConfig.BeaconUrl, err)
		}
	}

	log.Printf("[create-clients] preflight OK: cosmos_height=%d eth_chain_id=%s", status.SyncInfo.LatestBlockHeight, chainID.String())
	return nil
}

func cosmosHasWasmChecksum(cosmosClient *rpchttp.HTTP, checksum string) (bool, error) {
	checksum = strings.ToLower(strings.TrimPrefix(checksum, "0x"))

	reqBytes, err := proto.Marshal(&ibcwasmtypes.QueryChecksumsRequest{})
	if err != nil {
		return false, fmt.Errorf("failed to marshal checksum query: %w", err)
	}

	result, err := cosmosClient.ABCIQuery(context.Background(), "/ibc.lightclients.wasm.v1.Query/Checksums", reqBytes)
	if err != nil {
		return false, fmt.Errorf("failed to query wasm checksums: %w", err)
	}
	if result.Response.Code != 0 {
		return false, fmt.Errorf("checksum query failed with code %d: %s", result.Response.Code, result.Response.Log)
	}

	var resp ibcwasmtypes.QueryChecksumsResponse
	if err := proto.Unmarshal(result.Response.Value, &resp); err != nil {
		return false, fmt.Errorf("failed to unmarshal checksum query response: %w", err)
	}

	for _, existing := range resp.Checksums {
		if strings.ToLower(existing) == checksum {
			return true, nil
		}
	}

	return false, nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func roleManagerOrDefault(cfg *appConfig) string {
	return envOrDefault("ROLE_MANAGER", cfg.CosmosToEthConfig.ICS26Address)
}

func cosmosRouterClientIDOrDefault(cfg *appConfig) string {
	return envOrDefault("ICS26_CLIENT_ID", cfg.CosmosToEthConfig.ICS26ClientID)
}

func validateStartupKeys() error {
	ethPrivKey := os.Getenv("ETH_PRIVATE_KEY")
	if ethPrivKey == "" {
		return fmt.Errorf("ETH_PRIVATE_KEY environment variable is required in .env file")
	}
	if _, err := keys.RestoreKey(ethPrivKey); err != nil {
		return fmt.Errorf("failed to restore ETH private key: %w", err)
	}

	if _, err := (&transaction.Handler{}).CosmosSignerAddress(); err != nil {
		if os.Getenv("COSMOS_PRIVATE_KEY") == "" {
			return fmt.Errorf("COSMOS_PRIVATE_KEY environment variable is required in .env file")
		}
		return fmt.Errorf("failed to decode COSMOS_PRIVATE_KEY: %w", err)
	}

	return nil
}

func cosmosWasmClientIDOrDefault(cfg *appConfig) string {
	return envOrDefault("COSMOS_WASM_CLIENT_ID", cfg.CosmosToEthConfig.CosmosWasmClientID)
}

// proofBackendFromFlags resolves the GPU/CPU backend from --gpu-prove or the
// GPU_PROVE env var. Returns (backend, true) when an explicit selection was
// made, otherwise (nil, false) so the caller falls back to env-only defaults.
func proofBackendFromFlags(cmd *cobra.Command) (prover.ProofBackend, bool, error) {
	flagSet := cmd.Flags().Changed(flagGPUProve)
	envSet := prover.GPUProveEnvEnabled()
	if !flagSet && !envSet {
		return nil, false, nil
	}

	useGPU := envSet
	if flagSet {
		v, err := cmd.Flags().GetBool(flagGPUProve)
		if err != nil {
			return nil, false, fmt.Errorf("failed to get gpu prove flag: %w", err)
		}
		useGPU = v
	}

	backend, err := prover.NewProofBackend(useGPU)
	if err != nil {
		return nil, false, err
	}
	return backend, true, nil
}

// --- Main ---

func main() {
	zLogger, _ := zap.NewProduction(zap.AddStacktrace(zap.DPanicLevel))
	defer zLogger.Sync()

	logger := zLogger.Sugar()

	rootCmd := &cobra.Command{
		Use:   "relayer [command]",
		Short: "fast-ibc operator — relay IBC packets between Cosmos and Ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(
		Start(zLogger),
		CreateClients(zLogger),
		Genesis(zLogger),
	)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

// CreateClients deploys ICS07 Tendermint light client on Ethereum and
// creates a wasm Ethereum light client on Cosmos, then registers counterparties.
// This is a one-time setup step before running `start`.
func CreateClients(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-clients",
		Short: "deploy light clients and register counterparties on both chains",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigPath)
			if err != nil {
				return fmt.Errorf("failed to get config path: %w", err)
			}

			_ = godotenv.Load()

			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			logger.Sugar().Infof("create-clients: config loaded from %s", configPath)
			logger.Sugar().Infof(
				"create-clients: endpoints cosmos_rpc=%s eth_rpc=%s beacon=%s",
				cfg.CosmosToEthConfig.TmRpcUrl,
				cfg.CosmosToEthConfig.EthRpcUrl,
				cfg.EthToCosmosConfig.BeaconUrl,
			)
			if err := preflightCreateClients(cfg); err != nil {
				return err
			}
			logger.Sugar().Info("create-clients: preflight passed")

			// Connect to Ethereum
			logger.Sugar().Infof("create-clients: dialing ethereum rpc %s", cfg.CosmosToEthConfig.EthRpcUrl)
			ethClient, err := ethclient.Dial(cfg.CosmosToEthConfig.EthRpcUrl)
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum: %w", err)
			}
			logger.Sugar().Info("create-clients: ethereum rpc connected")

			// Connect to Cosmos
			logger.Sugar().Infof("create-clients: creating cosmos rpc client %s", cfg.CosmosToEthConfig.TmRpcUrl)
			cosmosClient, err := rpchttp.New(cfg.CosmosToEthConfig.TmRpcUrl, "/websocket")
			if err != nil {
				return fmt.Errorf("failed to create Cosmos RPC client: %w", err)
			}
			logger.Sugar().Info("create-clients: cosmos rpc client created")

			worker := services.NewWorker(&transaction.Handler{}, nil)

			cosmosWasmClientID := cosmosWasmClientIDOrDefault(cfg)
			if cosmosWasmClientID == "" {
				return fmt.Errorf("cosmos_wasm_client_id is required in cosmos_to_eth config")
			}

			// Create context (no WS client needed for create-clients)
			ctx := services.NewCtxWithBeacon(
				cosmosClient, ethClient, nil,
				"",
				cfg.EthToCosmosConfig.BeaconUrl,
				cosmosWasmClientID,
			)
			ctx.SetCosmosRouterClientID(cosmosRouterClientIDOrDefault(cfg))

			roleManager := roleManagerOrDefault(cfg)
			ctx.SetAddresses(
				cfg.CosmosToEthConfig.ICS26Address,
				cfg.CosmosToEthConfig.WrapperVerifier,
				cfg.CosmosToEthConfig.Membership,
				cfg.CosmosToEthConfig.Misbehaviour,
				cfg.CosmosToEthConfig.UpdateClient,
				roleManager,
			)

			// Start Cosmos WS (needed for queries)
			logger.Sugar().Info("create-clients: starting cosmos websocket client")
			if err := cosmosClient.Start(); err != nil {
				return fmt.Errorf("failed to start Cosmos WS client: %w", err)
			}
			defer cosmosClient.Stop()
			logger.Sugar().Info("create-clients: cosmos websocket client started")

			wasmChecksum, err := cmd.Flags().GetString(flagWasmChecksum)
			if err != nil {
				return fmt.Errorf("failed to get wasm checksum: %w", err)
			}
			if wasmChecksum == "" {
				wasmChecksum = os.Getenv("WASM_CHECKSUM")
			}
			logger.Sugar().Infof("create-clients: wasm checksum present=%t", wasmChecksum != "")

			if wasmChecksum != "" && cfg.EthToCosmosConfig.BeaconUrl != "" {
				logger.Sugar().Infof("create-clients: validating wasm checksum on Cosmos before mutating ETH/config: %s", wasmChecksum)
				ok, err := cosmosHasWasmChecksum(cosmosClient, wasmChecksum)
				if err != nil {
					return fmt.Errorf("failed to validate wasm checksum on Cosmos: %w", err)
				}
				if !ok {
					return fmt.Errorf("wasm checksum %s has not been previously stored on Cosmos", wasmChecksum)
				}
				logger.Sugar().Info("create-clients: wasm checksum preflight passed")
			}

			// --- 1. Create Cosmos light client on Ethereum (deploy ICS07) ---
			trustLevel, err := cmd.Flags().GetString(flagTrustLevel)
			if err != nil {
				return fmt.Errorf("failed to get trust level: %w", err)
			}

			trustingPeriod, err := cmd.Flags().GetUint32(flagTrustingPeriod)
			if err != nil {
				return fmt.Errorf("failed to get trusting period: %w", err)
			}
			if trustingPeriod == 0 {
				unbondingPeriod, err := tendermintClient.GetUnbondingTime(cosmosClient)
				if err != nil {
					return fmt.Errorf("failed to fetch unbonding time: %w", err)
				}
				trustingPeriod = 2 * uint32(unbondingPeriod) / 3
			}

			logger.Sugar().Infof("Creating Cosmos light client on Ethereum (trustingPeriod=%d, trustLevel=%s)...", trustingPeriod, trustLevel)
			ics07Addr, err := worker.CreateCosmosClient(ctx, "groth16", trustingPeriod, 0, trustLevel)
			if err != nil {
				return fmt.Errorf("failed to create Cosmos client on Ethereum: %w", err)
			}
			if (ics07Addr == common.Address{}) {
				return fmt.Errorf("ics07 address missing after deploy")
			}
			ctx.SetClient(ics07Addr)

			// --- 2. Create Ethereum light client on Cosmos (wasm) ---
			if wasmChecksum != "" && cfg.EthToCosmosConfig.BeaconUrl != "" {
				logger.Sugar().Infof("Creating Ethereum light client on Cosmos (checksum=%s)...", wasmChecksum)
				ethClientID, err := worker.CreateEthClient(ctx, wasmChecksum)
				if err != nil {
					return fmt.Errorf("failed to create Ethereum client on Cosmos: %w", err)
				}
				if ethClientID != cosmosWasmClientID {
					return fmt.Errorf(
						"created Ethereum light client ID %s does not match configured cosmos_wasm_client_id %s",
						ethClientID, cosmosWasmClientID,
					)
				}
				logger.Sugar().Infof("Ethereum light client created on Cosmos: clientID=%s", ethClientID)
			} else {
				if wasmChecksum == "" {
					logger.Sugar().Warn("Skipping ETH client creation: --wasm-checksum not provided")
				}
				if cfg.EthToCosmosConfig.BeaconUrl == "" {
					logger.Sugar().Warn("Skipping ETH client creation: beacon URL not configured")
				}
			}

			if err := writeICS07Address(configPath, ics07Addr.Hex()); err != nil {
				return fmt.Errorf("persist ics07 address to %s: %w", configPath, err)
			}
			logger.Sugar().Infof("create-clients: wrote ics07_client=%s into %s", ics07Addr.Hex(), configPath)

			logger.Sugar().Infof("=== Setup Complete ===")
			logger.Sugar().Infof("ICS07 address has been persisted to %s; ready for 'start'", configPath)

			return nil
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().String(flagTrustLevel, "2/3", "trust level for Cosmos light client (e.g., 1/3, 2/3)")
	cmd.Flags().Uint32(flagTrustingPeriod, 0, "trusting period in seconds for Cosmos light client (default: 2/3 of chain unbonding period)")
	cmd.Flags().String(flagWasmChecksum, "", "wasm checksum for Ethereum light client (hex)")
	return cmd
}

// Start starts the relay service loop.
// It loads config from a JSON file, connects to Cosmos and Ethereum,
// and subscribes to send_packet events to relay packets.
func Start(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start the relay loop",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigPath)
			if err != nil {
				return fmt.Errorf("failed to get config path: %w", err)
			}

			// Load .env for prover paths, private keys, etc.
			_ = godotenv.Load()

			benchmarkFlag, err := cmd.Flags().GetBool(flagBenchmark)
			if err != nil {
				return fmt.Errorf("failed to get benchmark flag: %w", err)
			}
			utils.SetBenchEnabled(benchmarkFlag || utils.BenchEnabled())
			if utils.BenchEnabled() {
				log.Printf("[benchmark] enabled: detailed gas/timing logs are active")
			}

			if err := validateStartupKeys(); err != nil {
				return err
			}

			// Load JSON config
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Connect to Ethereum (HTTP for queries)
			ethClient, err := ethclient.Dial(cfg.CosmosToEthConfig.EthRpcUrl)
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum: %w", err)
			}

			// Connect to Ethereum (WS for subscriptions)
			var ethWsClient *ethclient.Client
			if cfg.CosmosToEthConfig.EthWsUrl != "" {
				if !strings.HasPrefix(cfg.CosmosToEthConfig.EthWsUrl, "ws://") &&
					!strings.HasPrefix(cfg.CosmosToEthConfig.EthWsUrl, "wss://") {
					return fmt.Errorf("eth_ws_url must use ws:// or wss://, got: %s", cfg.CosmosToEthConfig.EthWsUrl)
				}
				ethWsClient, err = ethclient.Dial(cfg.CosmosToEthConfig.EthWsUrl)
				if err != nil {
					return fmt.Errorf("failed to connect to Ethereum WS: %w", err)
				}
			}

			// Connect to Cosmos
			cosmosClient, err := rpchttp.New(cfg.CosmosToEthConfig.TmRpcUrl, "/websocket")
			if err != nil {
				return fmt.Errorf("failed to create Cosmos RPC client: %w", err)
			}

			// Load prover (one bucket per supported validator count)
			binDir := envOrDefault("PROVER_BIN_DIR", "./bin")
			selectedBackend, hasBackendOverride, err := proofBackendFromFlags(cmd)
			if err != nil {
				return fmt.Errorf("failed to resolve proof backend: %w", err)
			}

			var p *prover.EcipProver
			if hasBackendOverride {
				logger.Sugar().Infof("start: overriding proof backend via flags: %s", selectedBackend.Name())
				p, err = prover.NewProverWithBackend(binDir, selectedBackend)
			} else {
				p, err = prover.NewProver(binDir)
			}
			if err != nil {
				return fmt.Errorf("failed to load prover: %w", err)
			}

			cosmosWasmClientID := cosmosWasmClientIDOrDefault(cfg)
			if cosmosWasmClientID == "" {
				return fmt.Errorf("cosmos_wasm_client_id is required in cosmos_to_eth config")
			}

			// Create context with beacon API
			ctx := services.NewCtxWithBeacon(
				cosmosClient, ethClient, ethWsClient,
				cfg.CosmosToEthConfig.EthWsUrl,
				cfg.EthToCosmosConfig.BeaconUrl,
				cosmosWasmClientID,
			)
			ctx.SetCosmosRouterClientID(cosmosRouterClientIDOrDefault(cfg))

			// Set contract addresses from config
			roleManager := roleManagerOrDefault(cfg)
			ctx.SetAddresses(
				cfg.CosmosToEthConfig.ICS26Address,
				cfg.CosmosToEthConfig.WrapperVerifier,
				cfg.CosmosToEthConfig.Membership,
				cfg.CosmosToEthConfig.Misbehaviour,
				cfg.CosmosToEthConfig.UpdateClient,
				roleManager,
			)

			// Set ICS07 client address (already deployed)
			if cfg.CosmosToEthConfig.ICS07Client == "" {
				return fmt.Errorf("ics07_client address is required in cosmos_to_eth config")
			}
			ctx.SetClient(common.HexToAddress(cfg.CosmosToEthConfig.ICS07Client))

			// Start Cosmos WebSocket client
			if err := cosmosClient.Start(); err != nil {
				return fmt.Errorf("failed to start Cosmos WS client: %w", err)
			}
			defer cosmosClient.Stop()

			logger.Sugar().Info("Relayer started, subscribing to events...")

			cosmosConfig := services.DefaultConfig()
			if cfg.CosmosToEthConfig.TrustingPeriod != 0 {
				cosmosConfig.TrustingPeriod = cfg.CosmosToEthConfig.TrustingPeriod
			}
			if cfg.CosmosToEthConfig.TrustLevel != "" {
				cosmosConfig.TrustLevel = cfg.CosmosToEthConfig.TrustLevel
			}
			if cfg.CosmosToEthConfig.ProofType != "" {
				cosmosConfig.ProofType = cfg.CosmosToEthConfig.ProofType
			}
			ctx.Config = cosmosConfig
			svc := services.New(
				subscriber.NewSubscriber(),
				&transaction.Handler{},
				p,
				services.DefaultConfig(),
				cosmosConfig,
			)
			svc.StartLoop(ctx)

			return nil
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().Bool(flagGPUProve, false, "use the ICICLE GPU backend for proving (or set GPU_PROVE=1); requires an icicle-enabled build")
	cmd.Flags().Bool(flagBenchmark, false, "enable detailed benchmark gas/timing logs (or set RELAYER_BENCHMARK=1)")
	return cmd
}

// Genesis generates the genesis state for a new client.
func Genesis(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genesis",
		Short: "genesis",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := godotenv.Load()
			if err != nil {
				return fmt.Errorf("error loading .env file: %v", err)
			}
			tendermintRpcEndpoint := os.Getenv("TENDERMINT_RPC_URL")
			if tendermintRpcEndpoint == "" {
				return fmt.Errorf("TENDERMINT_RPC_URL environment variable is required in .env file")
			}
			tendermintRpcClient, err := rpchttp.New(tendermintRpcEndpoint, "/websocket")
			if err != nil {
				return fmt.Errorf("failed to create RPC client: %w", err)
			}
			trustedBlock, err := cmd.Flags().GetInt64(flagTrustedBlock)
			if err != nil {
				return fmt.Errorf("failed to get trusted block: %w", err)
			}
			trustingPeriod, err := cmd.Flags().GetUint32(flagTrustingPeriod)
			if err != nil {
				return fmt.Errorf("failed to get trusting period: %w", err)
			}
			trustLevel, err := cmd.Flags().GetString(flagTrustLevel)
			if err != nil {
				return fmt.Errorf("failed to get trust level from flag: %w", err)
			}
			proofType, err := cmd.Flags().GetString(flagProofType)
			if err != nil {
				return fmt.Errorf("failed to get proof type from flag: %w", err)
			}

			genesis, err := tendermintClient.GetGenesis(tendermintRpcClient, trustedBlock, trustingPeriod, trustLevel, proofType)
			if err != nil {
				return fmt.Errorf("failed to get genesis: %w", err)
			}

			outputType, err := cmd.Flags().GetString(flagOutput)
			if err != nil {
				return fmt.Errorf("failed to get output type from flag: %w", err)
			}

			data, err := json.Marshal(genesis)
			if err != nil {
				return fmt.Errorf("failed to marshal genesis state: %w", err)
			}

			switch outputType {
			case "json":
				fmt.Println(string(data))
			case "file":
				outputDir, err := cmd.Flags().GetString(flagOutputPath)
				if err != nil {
					return fmt.Errorf("failed to get output path from flag: %w", err)
				}
				if err := os.WriteFile(outputDir, data, 0644); err != nil {
					return fmt.Errorf("failed to write genesis state to file: %w", err)
				}
			default:
				return fmt.Errorf("unsupported output type: %s, supported types are: json, file", outputType)
			}

			return nil
		},
	}
	cmd.Flags().String(flagProofType, "groth16", "the type of proof to use (groth16, plonk)")
	cmd.Flags().Int64(flagTrustedBlock, 0, "the trusted block height, if <height> is 0 then catch latest block")
	cmd.Flags().String(flagOutput, "json", "the output structure for the genesis state (json, file)")
	cmd.Flags().String(flagOutputPath, "./data/genesis.json", "the path to the output file for the genesis state")
	cmd.Flags().String(flagTrustLevel, "2/3", "the trust level for the genesis state (e.g., 2/3)")
	cmd.Flags().Uint32(flagTrustingPeriod, 0, "the trusting period for the genesis state")
	return cmd
}
