package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"

	tendermintContract "relayer/bindings/Groth16ICS07Tendermint"
	contractICS26Router "relayer/bindings/ICS26Router"
	tendermintClient "relayer/client"
	"relayer/keys"
	"relayer/prover"
	"relayer/runner"
	"relayer/services"
	"relayer/subscriber"
	"relayer/transaction"
	"relayer/utils"
)

const (
	flagConfigPath     = "config"
	flagOnlyOnce       = "only-once"
	flagProofType      = "proof-type"
	flagOutput         = "output"
	flagOutputPath     = "output-path"
	flagTrustLevel     = "trust-level"
	flagTrustingPeriod = "trusting-period"
	flagTrustedBlock   = "trusted-block"
	flagMembership     = "membership"
	flagWasmChecksum   = "wasm-checksum"
)

// --- Config types for JSON config file ---

type cosmosToEthConfig struct {
	TmRpcUrl        string `json:"tm_rpc_url"`
	ICS26Address    string `json:"ics26_address"`
	EthRpcUrl       string `json:"eth_rpc_url"`
	ICS07Client     string `json:"ics07_client"`
	WrapperVerifier string `json:"wrapper_verifier"`
	Membership      string `json:"membership"`
	Misbehaviour    string `json:"misbehaviour"`
	UpdateClient    string `json:"update_client"`
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

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
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
		Fixtures(zLogger),
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

			// Connect to Ethereum
			ethClient, err := ethclient.Dial(cfg.CosmosToEthConfig.EthRpcUrl)
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum: %w", err)
			}

			// Connect to Cosmos
			cosmosClient, err := rpchttp.New(cfg.CosmosToEthConfig.TmRpcUrl, "/websocket")
			if err != nil {
				return fmt.Errorf("failed to create Cosmos RPC client: %w", err)
			}

			// Load prover
			r1csPath := envOrDefault("PROVER_R1CS_PATH", "./bin/r1cs.bin")
			pkPath := envOrDefault("PROVER_PK_PATH", "./bin/pk.bin")
			vkPath := envOrDefault("PROVER_VK_PATH", "./bin/vk.bin")

			p, err := prover.NewProver(r1csPath, pkPath, vkPath)
			if err != nil {
				return fmt.Errorf("failed to load prover: %w", err)
			}

			worker := services.NewWorker(&transaction.Handler{}, p)

			// Create context
			ctx := services.NewCtxWithBeacon(
				cosmosClient, ethClient,
				cfg.EthToCosmosConfig.BeaconUrl,
				"08-wasm-0",
			)

			roleManager := envOrDefault("ROLE_MANAGER", "0x0000000000000000000000000000000000000000")
			ctx.SetAddresses(
				cfg.CosmosToEthConfig.ICS26Address,
				cfg.CosmosToEthConfig.WrapperVerifier,
				cfg.CosmosToEthConfig.Membership,
				cfg.CosmosToEthConfig.Misbehaviour,
				cfg.CosmosToEthConfig.UpdateClient,
				roleManager,
			)

			// Start Cosmos WS (needed for queries)
			if err := cosmosClient.Start(); err != nil {
				return fmt.Errorf("failed to start Cosmos WS client: %w", err)
			}
			defer cosmosClient.Stop()

			// --- 1. Create Cosmos light client on Ethereum (deploy ICS07) ---
			trustLevel, err := cmd.Flags().GetString(flagTrustLevel)
			if err != nil {
				return fmt.Errorf("failed to get trust level: %w", err)
			}

			unbondingPeriod, err := tendermintClient.GetUnbondingTime(cosmosClient)
			if err != nil {
				return fmt.Errorf("failed to fetch unbonding time: %w", err)
			}
			trustingPeriod := 2 * uint32(unbondingPeriod) / 3

			logger.Sugar().Infof("Creating Cosmos light client on Ethereum (trustingPeriod=%d, trustLevel=%s)...", trustingPeriod, trustLevel)
			if err := worker.CreateCosmosClient(ctx, "groth16", trustingPeriod, 0, trustLevel); err != nil {
				return fmt.Errorf("failed to create Cosmos client on Ethereum: %w", err)
			}
			// ICS07 address is printed by handler: "[CreateCosmosClient] ICS07 deployed at <address>"
			// Copy that address into config.json cosmos_to_eth.ics07_client

			// --- 2. Create Ethereum light client on Cosmos (wasm) ---
			wasmChecksum, err := cmd.Flags().GetString(flagWasmChecksum)
			if err != nil {
				return fmt.Errorf("failed to get wasm checksum: %w", err)
			}
			if wasmChecksum == "" {
				wasmChecksum = os.Getenv("WASM_CHECKSUM")
			}

			if wasmChecksum != "" && cfg.EthToCosmosConfig.BeaconUrl != "" {
				logger.Sugar().Infof("Creating Ethereum light client on Cosmos (checksum=%s)...", wasmChecksum)
				if err := worker.CreateEthClient(ctx, wasmChecksum); err != nil {
					return fmt.Errorf("failed to create Ethereum client on Cosmos: %w", err)
				}
				logger.Sugar().Info("Ethereum light client created on Cosmos")
			} else {
				if wasmChecksum == "" {
					logger.Sugar().Warn("Skipping ETH client creation: --wasm-checksum not provided")
				}
				if cfg.EthToCosmosConfig.BeaconUrl == "" {
					logger.Sugar().Warn("Skipping ETH client creation: beacon URL not configured")
				}
			}

			logger.Sugar().Infof("=== Setup Complete ===")
			logger.Sugar().Infof("Copy ICS07 address from log above into config.json cosmos_to_eth.ics07_client before running 'start'")

			return nil
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().String(flagTrustLevel, "1/3", "trust level for Cosmos light client (e.g., 1/3, 2/3)")
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

			// Load JSON config
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Connect to Ethereum
			ethClient, err := ethclient.Dial(cfg.CosmosToEthConfig.EthRpcUrl)
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum: %w", err)
			}

			// Connect to Cosmos
			cosmosClient, err := rpchttp.New(cfg.CosmosToEthConfig.TmRpcUrl, "/websocket")
			if err != nil {
				return fmt.Errorf("failed to create Cosmos RPC client: %w", err)
			}

			// Load prover
			r1csPath := envOrDefault("PROVER_R1CS_PATH", "./bin/r1cs.bin")
			pkPath := envOrDefault("PROVER_PK_PATH", "./bin/pk.bin")
			vkPath := envOrDefault("PROVER_VK_PATH", "./bin/vk.bin")

			p, err := prover.NewProver(r1csPath, pkPath, vkPath)
			if err != nil {
				return fmt.Errorf("failed to load prover: %w", err)
			}

			// Create worker
			worker := services.NewWorker(&transaction.Handler{}, p)

			// Create context with beacon API
			ctx := services.NewCtxWithBeacon(
				cosmosClient, ethClient,
				cfg.EthToCosmosConfig.BeaconUrl,
				"08-wasm-0",
			)

			// Set contract addresses from config
			roleManager := envOrDefault("ROLE_MANAGER", "0x0000000000000000000000000000000000000000")
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

			// Cosmos → ETH relay (blocking)
			// ETH → Cosmos is not yet fully implemented (needs WS endpoint + storage proofs)
			subscribeCosmos(ctx, worker, logger)

			return nil
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	return cmd
}

// subscribeCosmos listens for send_packet events on Cosmos and relays them to Ethereum.
func subscribeCosmos(ctx services.Context, worker *services.Worker, logger *zap.Logger) {
	sub, err := ctx.CosmosClient().WSEvents.Subscribe(context.Background(), "", subscriber.COMETBFT_SEND_PACKET_EVENT)
	if err != nil {
		logger.Sugar().Fatalf("Failed to subscribe to Cosmos events: %v", err)
	}

	for e := range sub {
		sendPacketEvent := e.Events[subscriber.EVENT_SEND_PACKET_FIELD]
		if sendPacketEvent == nil {
			continue
		}

		packetBytes, err := hex.DecodeString(sendPacketEvent[0])
		if err != nil {
			logger.Sugar().Errorf("Failed to decode packet hex: %v", err)
			continue
		}

		var packet channeltypesv2.Packet
		if err := proto.Unmarshal(packetBytes, &packet); err != nil {
			logger.Sugar().Errorf("Failed to unmarshal packet: %v", err)
			continue
		}

		log.Printf("[Relay] Received send_packet seq=%d, src=%s, dst=%s, timeout=%d",
			packet.Sequence, packet.SourceClient, packet.DestinationClient, packet.TimeoutTimestamp)

		if err := relayPacket(ctx, worker, logger, packet); err != nil {
			logger.Sugar().Errorf("Failed to relay packet seq=%d: %v", packet.Sequence, err)
		}
	}
}

// relayPacket handles a single send_packet event: update client, prove membership, send recvPacket.
func relayPacket(ctx services.Context, worker *services.Worker, logger *zap.Logger, packet channeltypesv2.Packet) error {
	// Check if packet has already timed out
	ethHeader, err := ctx.EthClient().HeaderByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to get eth block header: %w", err)
	}
	if packet.TimeoutTimestamp > 0 && ethHeader.Time >= packet.TimeoutTimestamp {
		log.Printf("[Relay] Packet seq=%d already timed out, skipping", packet.Sequence)
		return nil
	}

	// Wait for next block so AppHash includes the packet commitment
	log.Printf("[Relay] Waiting 2 blocks for packet commitment in AppHash...")
	time.Sleep(6 * time.Second)

	// Update Cosmos light client on Ethereum
	latestTimestamp := ctx.LatestCosmosTimestamp()
	log.Printf("[Relay] Updating cosmos client from height %d...", latestTimestamp.LatestUpdateHeight)

	latestLightBlock, err := worker.UpdateCosmosClient(ctx, "groth16", int64(latestTimestamp.LatestUpdateHeight), "1/3")
	if err != nil {
		return fmt.Errorf("failed to update cosmos light client: %w", err)
	}

	// Re-check timeout after updateClient (ZK proof takes time)
	ethHeader, err = ctx.EthClient().HeaderByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to get eth block header: %w", err)
	}
	if packet.TimeoutTimestamp > 0 && ethHeader.Time >= packet.TimeoutTimestamp {
		log.Printf("[Relay] Packet seq=%d timed out during updateClient, skipping", packet.Sequence)
		return nil
	}

	latestTimestamp.LatestUpdateTime = time.Now()
	latestTimestamp.LatestUpdateHeight = uint64(latestLightBlock.BlockHeight)
	log.Printf("[Relay] Light client updated to height %d", latestLightBlock.BlockHeight)

	// Prove membership of the packet commitment
	ibcPath := utils.IbcCommitmentPath(packet, []byte{1})
	log.Printf("[Relay] Proving membership at height %d...", latestLightBlock.BlockHeight)

	value, proof, err := tendermintClient.ProvePath(ctx.CosmosClient(), latestLightBlock.BlockHeight, ibcPath)
	if err != nil {
		return fmt.Errorf("failed to prove path: %w", err)
	}

	if len(value) == 0 {
		return fmt.Errorf("packet commitment value is empty at height %d", latestLightBlock.BlockHeight)
	}

	// Build membership proof
	merkleProof := tendermintContract.IMembershipMsgsMerkleProof{
		Proofs: []tendermintContract.IMembershipMsgsCommitmentProof{},
	}
	for _, p := range proof.Proofs {
		commitmentProof, err := tendermintClient.ParseCommitmentProof(p)
		if err != nil {
			return fmt.Errorf("failed to parse commitment proof: %w", err)
		}
		merkleProof.Proofs = append(merkleProof.Proofs, *commitmentProof)
	}

	membershipMsg := tendermintContract.ILightClientMsgsMsgVerifyMembership{
		Height: tendermintContract.IICS02ClientMsgsHeight{
			RevisionHeight: uint64(latestLightBlock.BlockHeight),
			RevisionNumber: 0,
		},
		KvPairs: []tendermintContract.IMembershipMsgsKVPair{
			{
				Path:  ibcPath,
				Value: utils.BytesToBytes32(value),
			},
		},
		MerkleProofs: []tendermintContract.IMembershipMsgsMerkleProof{
			merkleProof,
		},
		AppHash: utils.BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
		TrustedConsensusState: tendermintContract.IICS07TendermintMsgsConsensusState{
			Timestamp:          big.NewInt(latestLightBlock.SignedHeader.Header.Time.UnixNano()),
			Root:               utils.BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
			NextValidatorsHash: utils.BytesToBytes32(latestLightBlock.SignedHeader.Header.NextValidatorsHash),
		},
		MembershipType: 0,
	}

	tendermintAbiJson, err := tendermintContract.ContractGroth16ICS07TendermintMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get ABI: %w", err)
	}
	calldata, err := tendermintAbiJson.Pack("verifyMembership", membershipMsg)
	if err != nil {
		return fmt.Errorf("failed to ABI encode verify msg: %w", err)
	}
	// Strip 4-byte function selector — ICS26Router does abi.decode, not a function call
	calldata = calldata[4:]

	// Build MsgRecvPacket
	payloads := make([]contractICS26Router.IICS26RouterMsgsPayload, 0, len(packet.Payloads))
	for _, p := range packet.Payloads {
		payloads = append(payloads, contractICS26Router.IICS26RouterMsgsPayload{
			SourcePort: p.SourcePort,
			DestPort:   p.DestinationPort,
			Version:    p.Version,
			Encoding:   p.Encoding,
			Value:      p.Value,
		})
	}

	msgRecvPacket := contractICS26Router.IICS26RouterMsgsMsgRecvPacket{
		Packet: contractICS26Router.IICS26RouterMsgsPacket{
			Sequence:         packet.Sequence,
			SourceClient:     packet.SourceClient,
			DestClient:       packet.DestinationClient,
			TimeoutTimestamp: packet.TimeoutTimestamp,
			Payloads:         payloads,
		},
		MembershipMsg: calldata,
	}

	log.Printf("[Relay] Sending recvPacket tx to Ethereum...")
	if err := worker.TxHandler.SendEthTx(ctx, msgRecvPacket); err != nil {
		return fmt.Errorf("failed to send recvPacket: %w", err)
	}
	log.Printf("[Relay] recvPacket tx succeeded for seq=%d!", packet.Sequence)
	return nil
}

// subscribeEth listens for SendPacket events on Ethereum (ICS26Router) and relays them to Cosmos.
func subscribeEth(ctx services.Context, worker *services.Worker, logger *zap.Logger) {
	ics26RouterAddr := ctx.RouterContract()
	if ics26RouterAddr == nil {
		logger.Sugar().Fatal("ICS26Router address not set in context")
	}

	filterer, err := contractICS26Router.NewContractICS26RouterFilterer(*ics26RouterAddr, ctx.EthClient())
	if err != nil {
		logger.Sugar().Fatalf("Failed to create ICS26Router filterer: %v", err)
	}

	sendPacketCh := make(chan *contractICS26Router.ContractICS26RouterSendPacket)
	watchOpts := &bind.WatchOpts{Context: context.Background()}

	sendPacketSub, err := filterer.WatchSendPacket(watchOpts, sendPacketCh, nil, nil)
	if err != nil {
		logger.Sugar().Fatalf("Failed to subscribe to ETH SendPacket events: %v", err)
	}
	defer sendPacketSub.Unsubscribe()

	logger.Sugar().Info("Subscribed to ETH SendPacket events")

	for {
		select {
		case ev := <-sendPacketCh:
			packet := subscriber.EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
			log.Printf("[ETH→Cosmos] Received SendPacket seq=%d, src=%s, dst=%s",
				packet.Sequence, packet.SourceClient, packet.DestinationClient)

			// Update ETH light client on Cosmos
			if err := worker.UpdateEthClient(ctx); err != nil {
				logger.Sugar().Errorf("Failed to update ETH light client: %v", err)
				continue
			}

			// TODO: Prove ETH storage commitment and build MsgRecvPacket for Cosmos
			// This requires:
			// 1. Get ETH storage proof for the packet commitment slot
			// 2. Build MsgRecvPacket with the proof
			// 3. Send via worker.TxHandler.SendCosmosTx(ctx, msgRecvPacket)
			logger.Sugar().Warnf("ETH→Cosmos relay not yet implemented for seq=%d", packet.Sequence)

		case err := <-sendPacketSub.Err():
			logger.Sugar().Errorf("ETH SendPacket subscription error: %v", err)
			return
		}
	}
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

// Fixtures generates fixture files for testing.
func Fixtures(logger *zap.Logger) *cobra.Command {
	fixturesCmd := &cobra.Command{
		Use:   "fixtures",
		Short: "fixtures",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	fixturesCmd.AddCommand(MembershipCmd(logger))

	return fixturesCmd
}

// MembershipCmd verifies a membership/non-membership proof on-chain.
func MembershipCmd(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "membership",
		Short: "membership <key_path> <is_base64> <membership_type>",
		Args:  cobra.ExactArgs(3),
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
			ethRpcEndpoint := os.Getenv("ETH_RPC_URL")
			if ethRpcEndpoint == "" {
				return fmt.Errorf("ETH_RPC_URL environment variable is required in .env file")
			}

			hexAddress := os.Getenv("CONTRACT_ADDRESS")
			if hexAddress == "" {
				return fmt.Errorf("CONTRACT_ADDRESS environment variable is required in .env file")
			}

			privKey := os.Getenv("PRIVATE_KEY")
			if privKey == "" {
				return fmt.Errorf("PRIVATE_KEY environment variable is required in .env file")
			}
			privateKey, err := keys.RestoreKey(privKey)
			if err != nil {
				return fmt.Errorf("failed to restore private key: %w", err)
			}

			chainIdEth := os.Getenv("CHAIN_ID")
			if chainIdEth == "" {
				return fmt.Errorf("CHAIN_ID environment variable is required in .env file")
			}

			chainIdInt := big.NewInt(0)
			chainIdInt, ok := chainIdInt.SetString(chainIdEth, 10)
			if !ok {
				return fmt.Errorf("invalid chain id: %v", chainIdEth)
			}

			ethClient, err := ethclient.Dial(ethRpcEndpoint)
			if err != nil {
				return fmt.Errorf("failed to create Ethereum client: %w", err)
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

			tendermintAddr := common.HexToAddress(hexAddress)
			ics07Tendermint, err := tendermintContract.NewContractGroth16ICS07Tendermint(
				tendermintAddr,
				ethClient,
			)
			if err != nil {
				return fmt.Errorf("failed to create ICS07 Tendermint contract: %w", err)
			}

			isMembership, err := cmd.Flags().GetBool(flagMembership)
			if err != nil {
				return fmt.Errorf("failed to get membership flag: %w", err)
			}

			publicKey, err := keys.PublicKey(privateKey)
			if err != nil {
				log.Fatal(err)
			}

			fromAddress := crypto.PubkeyToAddress(*publicKey)
			nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
			if err != nil {
				log.Fatal(err)
			}
			gasPrice, err := ethClient.SuggestGasPrice(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIdInt)
			if err != nil {
				return fmt.Errorf("failed to create auth transactor: %w", err)
			}
			auth.Nonce = big.NewInt(int64(nonce))
			auth.Value = big.NewInt(0)
			auth.GasLimit = uint64(300000)
			auth.GasPrice = gasPrice

			kvPairs, proofs, err := runner.RunMembership(tendermintRpcClient, args[0], trustedBlock, args[1] == "true")
			if err != nil {
				return err
			}

			membershipType, err := cmd.Flags().GetInt(flagMembership)
			if err != nil {
				membershipType = 0
			}

			if isMembership {
				msg := tendermintContract.ILightClientMsgsMsgVerifyMembership{
					Height:                tendermintContract.IICS02ClientMsgsHeight(genesis.TrustedClientState.LatestHeight),
					KvPairs:               kvPairs,
					MerkleProofs:          proofs,
					AppHash:               genesis.TrustedConsensusState.Root,
					TrustedConsensusState: tendermintContract.IICS07TendermintMsgsConsensusState(genesis.TrustedConsensusState),
					MembershipType:        uint8(membershipType),
				}

				tx, err := ics07Tendermint.VerifyMembership(auth, msg)
				if err != nil {
					return fmt.Errorf("failed to verify membership: %w", err)
				}
				fmt.Println("tx hash:", tx.Hash().Hex())
			} else {
				msg := tendermintContract.ILightClientMsgsMsgVerifyNonMembership{
					Height:                tendermintContract.IICS02ClientMsgsHeight(genesis.TrustedClientState.LatestHeight),
					KvPairs:               kvPairs,
					MerkleProofs:          proofs,
					AppHash:               genesis.TrustedConsensusState.Root,
					TrustedConsensusState: tendermintContract.IICS07TendermintMsgsConsensusState(genesis.TrustedConsensusState),
					MembershipType:        uint8(membershipType),
				}

				tx, err := ics07Tendermint.VerifyNonMembership(auth, msg)
				if err != nil {
					return fmt.Errorf("failed to verify non-membership: %w", err)
				}
				fmt.Println("tx hash:", tx.Hash().Hex())
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
	cmd.Flags().Bool(flagMembership, true, "verify membership/non-membership proof")
	return cmd
}
