package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	tendermintContract "relayer/bindings/Groth16ICS07Tendermint"
	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/client"
	"relayer/prover"
	"relayer/services"
	"relayer/subscriber"
	"relayer/transaction"
	"relayer/utils"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/joho/godotenv"
)

type CosmosToEthConfig struct {
	TmRpcUrl        string `json:"tm_rpc_url"`
	ICS26Address    string `json:"ics26_address"`
	EthRpcUrl       string `json:"eth_rpc_url"`
	WrapperVerifier string `json:"wrapper_verifier"`
	Membership      string `json:"membership"`
	Misbehaviour    string `json:"misbehaviour"`
	UpdateClient    string `json:"update_client"`
}

type EthToCosmosConfig struct {
	TmRpcUrl      string `json:"tm_rpc_url"`
	ICS26Address  string `json:"ics26_address"`
	EthRpcUrl     string `json:"eth_rpc_url"`
	BeaconUrl     string `json:"eth_beacon_api_url"`
	SignerAddress string `json:"signer_address"`
}

type Module struct {
	Name     string          `json:"name"`
	SrcChain string          `json:"src_chain"`
	DstChain string          `json:"dst_chain"`
	Config   json.RawMessage `json:"config"`
}

type ServerConfig struct {
	LogLevel string `json:"log_level"`
	Address  string `json:"address"`
	Port     uint64 `json:"port"`
}

type AppJsonConfig struct {
	SeverConfig ServerConfig `json:"server"`
	Modules     []Module     `json:"modules"`
}

type AppConfig struct {
	SeverConfig       ServerConfig
	EthToCosmosConfig EthToCosmosConfig
	CosmosToEthConfig CosmosToEthConfig
}

func loadConfig(configPath string) (*AppConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var appConfig AppJsonConfig
	if err := json.Unmarshal(data, &appConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	var c2eCfg CosmosToEthConfig
	var e2cfg EthToCosmosConfig
	for _, m := range appConfig.Modules {

		if m.Name == "cosmos_to_eth" {
			fmt.Println("m.Config: ", string(m.Config))
			if err := json.Unmarshal(m.Config, &c2eCfg); err != nil {
				return nil, fmt.Errorf("failed to parse cosmos_to_eth config: %w", err)
			}
		}
		if m.Name == "eth_to_cosmos" {
			if err := json.Unmarshal(m.Config, &e2cfg); err != nil {
				return nil, fmt.Errorf("failed to parse eth_to_cosmos config: %w", err)
			}
		}
	}

	return &AppConfig{
		SeverConfig:       appConfig.SeverConfig,
		CosmosToEthConfig: c2eCfg,
		EthToCosmosConfig: e2cfg,
	}, nil
}

// type EurekaEvent struct {
// 	eventType string
// 	packet    channeltypesv2.Packet
// 	ack       *channeltypesv2.Acknowledgement
// }

func init() {
	// tendermintAbiJson, initErr = os.ReadFile("../../abi/SP1ICS07Tendermint.json")
	// if initErr != nil {
	// 	log.Fatal(initErr)
	// }
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Listener struct{}

func (l *Listener) SubscribeCosmos(ctx services.Context, worker *services.Worker) {
	sub, err := ctx.CosmosClient().WSEvents.Subscribe(context.Background(), "", subscriber.COMETBFT_SEND_PACKET_EVENT)
	if err != nil {
		ctx.Logger.Println(err.Error())
	}
	for {
		select {
		case e := <-sub:
			// handle event
			sendPacketEvent := e.Events[subscriber.EVENT_SEND_PACKET_FIELD]
			if sendPacketEvent == nil {
				continue
			}

			packetEncodedStr := sendPacketEvent[0]
			packetBytes, err := hex.DecodeString(packetEncodedStr)
			if err != nil {
				// TODO handle log here
				ctx.Logger.Println(fmt.Errorf("Failed to decode packet hex: %s", err.Error()))
				continue
			}

			var packet channeltypesv2.Packet
			err = proto.Unmarshal(packetBytes, &packet)
			if err != nil {
				// TODO handle log here
				ctx.Logger.Println(fmt.Errorf("Failed to unmarshal packet: %s", err.Error()))
				continue
			}

			log.Printf("[Listener] Received send_packet event, packet seq=%d, src=%s, dst=%s timeout=%d", packet.Sequence, packet.SourceClient, packet.DestinationClient, packet.TimeoutTimestamp)

			// Check if packet has already timed out before doing expensive ZK work
			ethHeader, err := ctx.EthClient().HeaderByNumber(context.Background(), nil)
			if err != nil {
				ctx.Logger.Println(fmt.Errorf("Failed to get eth block header: %s", err.Error()))
				continue
			}
			if packet.TimeoutTimestamp > 0 && ethHeader.Time >= packet.TimeoutTimestamp {
				log.Printf("[Listener] Packet seq=%d already timed out (timeout=%d <= eth_block_time=%d), skipping",
					packet.Sequence, packet.TimeoutTimestamp, ethHeader.Time)
				continue
			}

			// Wait for next block so AppHash includes the packet commitment
			// AppHash at block N+1 contains the state after block N's txs
			log.Printf("[Listener] Waiting 2 blocks for packet commitment to be included in AppHash...")
			time.Sleep(6 * time.Second)

			latestEthTimestamp := ctx.LatestCosmosTimestamp()
			log.Printf("[Listener] Updating cosmos client from height %d...", latestEthTimestamp.LatestUpdateHeight)
			latestLightBlock, err := worker.UpdateCosmosClient(ctx, "groth16", int64(latestEthTimestamp.LatestUpdateHeight), "1/3")
			if err != nil {
				ctx.Logger.Println(fmt.Errorf("Failed to update cosmos light client: %s", err.Error()))
				continue
			}

			// Re-check timeout after updateClient (ZK proof takes time)
			ethHeader, err = ctx.EthClient().HeaderByNumber(context.Background(), nil)
			if err != nil {
				ctx.Logger.Println(fmt.Errorf("Failed to get eth block header: %s", err.Error()))
				continue
			}
			if packet.TimeoutTimestamp > 0 && ethHeader.Time >= packet.TimeoutTimestamp {
				log.Printf("[Listener] Packet seq=%d timed out during updateClient (timeout=%d <= eth_block_time=%d), skipping",
					packet.Sequence, packet.TimeoutTimestamp, ethHeader.Time)
				continue
			}
			latestEthTimestamp.LatestUpdateTime = time.Now()
			latestEthTimestamp.LatestUpdateHeight = uint64(latestLightBlock.BlockHeight)
			log.Printf("[Listener] Light client updated to height %d", latestLightBlock.BlockHeight)

			ibcPath := utils.IbcCommitmentPath(packet, []byte{1})
			log.Printf("[Listener] Proving membership at height %d...", latestLightBlock.BlockHeight)
			value, proof, err := client.ProvePath(ctx.CosmosClient(), latestLightBlock.BlockHeight, ibcPath)
			if err != nil {
				ctx.Logger.Println(fmt.Errorf("Failed to prove path: %s", err.Error()))
				continue
			}
			log.Printf("[Listener] Proof obtained, value len=%d, proofs count=%d", len(value), len(proof.Proofs))
			log.Printf("[Listener] AppHash: %x", latestLightBlock.SignedHeader.AppHash)
			log.Printf("[Listener] Value: %x", value)
			log.Printf("[Listener] ConsensusState: timestamp=%d, root(appHash)=%x, nextValHash=%x",
				latestLightBlock.SignedHeader.Header.Time.Unix(),
				latestLightBlock.SignedHeader.AppHash,
				latestLightBlock.SignedHeader.Header.NextValidatorsHash)

			if len(value) == 0 {
				ctx.Logger.Println("WARNING: packet commitment value is empty at this height, skipping")
				continue
			}

			merkleProof := tendermintContract.IMembershipMsgsMerkleProof{
				Proofs: []tendermintContract.IMembershipMsgsCommitmentProof{},
			}
			for _, p := range proof.Proofs {
				commitmentProof, err := client.ParseCommitmentProof(p)
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("failed to parse commitment proof: %w", err))
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
				// current appHash
				AppHash: utils.BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
				// trusted consensus from revision height block
				TrustedConsensusState: tendermintContract.IICS07TendermintMsgsConsensusState{
					Timestamp:          big.NewInt(latestLightBlock.SignedHeader.Header.Time.UnixNano()),
					Root:               utils.BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
					NextValidatorsHash: utils.BytesToBytes32(latestLightBlock.SignedHeader.Header.NextValidatorsHash),
				},
				MembershipType: 0, // Membership (not MembershipAndUpdateClient)
			}

			tendermintAbiJson, err := tendermintContract.ContractGroth16ICS07TendermintMetaData.GetAbi()
			if err != nil {
				ctx.Logger.Println(fmt.Errorf("Failed to abi encode verify msg: %s", err.Error()))
			}
			calldata, err := tendermintAbiJson.Pack("verifyMembership", membershipMsg)
			if err != nil {
				ctx.Logger.Println(fmt.Errorf("Failed to abi encode verify msg: %s", err.Error()))
			}
			// Strip 4-byte function selector — ICS26Router does abi.decode, not a function call
			calldata = calldata[4:]

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

			log.Printf("[Listener] Sending recvPacket tx to Ethereum...")
			err = worker.TxHandler.SendEthTx(ctx, msgRecvPacket)
			if err != nil {
				ctx.Logger.Println(fmt.Errorf("Failed to send recvPacket: %s", err.Error()))
				continue
			}
			log.Printf("[Listener] recvPacket tx succeeded!")
		}
	}
}
func main() {
	cfg, err := loadConfig("./config.example.json")
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err).Error())
	}
	fmt.Println("cfg: ", cfg)

	ethRpcEndpoint := cfg.EthToCosmosConfig.EthRpcUrl
	ethClient, err := ethclient.Dial(ethRpcEndpoint)
	if err != nil {
		panic(fmt.Errorf("failed to connect to client: %s: ", err.Error()))
	}

	cosmosRpcEndpoint := cfg.CosmosToEthConfig.TmRpcUrl
	cosmosClient, err := rpchttp.New(cosmosRpcEndpoint, "/websocket")
	if err != nil {
		panic(fmt.Errorf("failed to create RPC client: %w", err))
	}

	prover, err := prover.NewProver("./bin/r1cs.bin", "./bin/pk.bin", "./bin/vk.bin")
	if err != nil {
		panic(fmt.Errorf("failed reading prover key: %w", err))
	}
	worker := services.NewWorker(&transaction.Handler{}, prover)

	ctx := services.NewCtxWithBeacon(cosmosClient, ethClient, cfg.EthToCosmosConfig.BeaconUrl, "08-wasm-0")
	ctx.SetAddresses(cfg.CosmosToEthConfig.ICS26Address, cfg.CosmosToEthConfig.WrapperVerifier, cfg.CosmosToEthConfig.Membership, cfg.CosmosToEthConfig.Misbehaviour, cfg.CosmosToEthConfig.UpdateClient, "0x8943545177806ED17B9F23F0a21ee5948eCaa776")

	ics07 := common.HexToAddress("0xD1ea1592b7927a2f0EE5f8567928Df0cfA687C78")
	ctx.SetClient(ics07)
	err = ctx.CosmosClient().Start()
	if err != nil {
		panic(err)
	}
	defer ctx.CosmosClient().Stop()
	listener := Listener{}

	unbondingPeriod, err := client.GetUnbondingTime(cosmosClient)
	if err != nil {
		panic(fmt.Errorf("failed to fetch unbonding time client: %w", err))
	}
	trustingPeriod := 2 * uint32(unbondingPeriod) / 3
	err = worker.CreateCosmosClient(ctx, "groth16", trustingPeriod, 0, "1/3")
	if err != nil {
		panic(fmt.Errorf("create client err: %w", err))
	}
	err = worker.CreateEthClient(ctx, "0xc6d93045091f05f6c056ca8fa583126902967b4b829085042529d279c188391c")
	if err != nil {
		panic(fmt.Errorf("create client err: %w", err))
	}

	listener.SubscribeCosmos(ctx, worker)
}
