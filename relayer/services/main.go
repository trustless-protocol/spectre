package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"relayer/utils"
	"os"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	tendermintContract "relayer/bindings/SP1ICS07Tendermint"
	"relayer/client"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
)

// read abi json file once in runtime
var tendermintAbiJson *abi.ABI
var initErr error

func init() {
	tendermintAbiJson, initErr = tendermintContract.ContractSP1ICS07TendermintMetaData.GetAbi()
	if initErr != nil {
		log.Fatal(initErr)
	}
}

type TransactionHandler interface {
	CreateCosmosClientContract(ctx Context, clientState, consensusHash []byte) error
	CreateEthClient(ctx Context, clientState ibcexported.ClientState, consensusState ibcexported.ConsensusState) error
	SendEthTx(ctx Context, msg any) error
	SendCosmosTx(ctx Context, msg any) error
	SendCosmosTxBatch(ctx Context, msgs []any) error
}

type Prover interface {
	GenerateProof(sig, pub, msg []byte) (
		proof [8]*big.Int,
		commitments [2]*big.Int,
		commitmentPok [2]*big.Int,
		err error,
	)
}
type EventListener interface {
	SubscribeCosmos(ctx Context, batchBuilder *BatchBuilder)
	SubscribeEth(ctx Context, batchBuilder *BatchBuilder)
}

type Services struct {
	// event listener
	listener EventListener
	worker   *Worker

	ethConfig    Config
	cosmosConfig Config

	BatchPackets chan BatchPackets
	BatchBuilder *BatchBuilder

	txHandler TransactionHandler
}

func New(rpcEndpoint string, eventListener EventListener, txHandler TransactionHandler, prover Prover, ethConfig, cosmosConfig Config) *Services {
	return &Services{
		listener:     eventListener,
		ethConfig:    ethConfig,
		cosmosConfig: cosmosConfig,
		worker: &Worker{
			txHandler,
			prover,
		},
		BatchPackets: make(chan BatchPackets),
		BatchBuilder: NewBatchBuilder(),
	}
}

func (s *Services) StartLoop() {
	rpcEndpoint := os.Getenv("TENDERMINT_RPC_URL")
	if rpcEndpoint == "" {
		panic(fmt.Errorf("TENDERMINT_RPC_URL environment variable is required in .env file"))
	}
	cosmosClient, err := rpchttp.New(rpcEndpoint, "/websocket")
	if err != nil {
		panic(fmt.Errorf("failed to create RPC client: %w", err))
	}

	ethRpcEndpoint := os.Getenv("ETH_RPC_URL")
	if ethRpcEndpoint == "" {
		panic(fmt.Errorf("ETH_RPC_URL environment variable is required in .env file"))
	}
	ethClient, err := ethclient.Dial(ethRpcEndpoint)
	if err != nil {
		panic(fmt.Errorf("failed to connect to client: %s: ", err.Error()))
	}

	ics26Router := os.Getenv("ICS26_ROUTER")
	if ics26Router == "" {
		panic(fmt.Errorf("ICS26_ROUTER environment variable is required in .env file"))
	}
	wrapVerifier := os.Getenv("WRAP_VERIFIER")
	if wrapVerifier == "" {
		panic(fmt.Errorf("WRAP_VERIFIER environment variable is required in .env file"))
	}

	membership := os.Getenv("MEMBERSHIP")
	if membership == "" {
		panic(fmt.Errorf("MEMBERSHIP environment variable is required in .env file"))
	}
	misbehaviour := os.Getenv("MISBEHAVIOUR")
	if misbehaviour == "" {
		panic(fmt.Errorf("MISBEHAVIOUR environment variable is required in .env file"))
	}
	updateClient := os.Getenv("UPDATE_CLIENT")
	if updateClient == "" {
		panic(fmt.Errorf("UPDATE_CLIENT environment variable is required in .env file"))
	}
	roleManager := os.Getenv("ROLE_MANAGER")
	if roleManager == "" {
		panic(fmt.Errorf("ROLE_MANAGER environment variable is required in .env file"))
	}

	ctx := NewCtx(cosmosClient, ethClient)
	ctx.SetAddresses(ics26Router, wrapVerifier, membership, misbehaviour, updateClient, roleManager)

	// listen to new tx events on Eth
	// add it to handler queue
	go func() {
		s.listener.SubscribeCosmos(ctx, s.BatchBuilder)
	}()

	// listen to new tx events on Cosmos
	// add it to handler queue
	go func() {
		s.listener.SubscribeEth(ctx, s.BatchBuilder)
	}()

	// routinely run update client
	go func() {
		for {
			// update client on Eth side routinely
			if ctx.latestEthTimestamp.LatestUpdateTime.Add(s.ethConfig.IntervalParams.blockTime).After(time.Now()) {
				latestBlock, err := s.worker.UpdateCosmosClient(ctx, "groth16", int64(ctx.latestEthTimestamp.LatestUpdateHeight), "2/3")
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("Failed to update cosmos light client: %s", err.Error()))
				}

				// update latest update time
				ctx.latestEthTimestamp.mtx.Lock()
				ctx.latestEthTimestamp.LatestUpdateTime = time.Now()
				ctx.latestEthTimestamp.LatestUpdateHeight = uint64(latestBlock.BlockHeight)
				ctx.latestEthTimestamp.mtx.Unlock()
			}

			// update client on Cosmos side routinely
			if ctx.latestCosmosTimestamp.LatestUpdateTime.Add(s.cosmosConfig.IntervalParams.blockTime).After(time.Now()) {
				s.worker.UpdateEthClient(ctx)

				// update latest update time
				ctx.latestCosmosTimestamp.mtx.Lock()
				ctx.latestCosmosTimestamp.LatestUpdateTime = time.Now()
				ctx.latestCosmosTimestamp.mtx.Unlock()

			}
		}
	}()

	// check for batch builder
	go func() {
		for {
			// check for batch every 3 seconds
			time.Sleep(time.Second * 3)
			s.BatchBuilder.CheckBatch(ctx.Config.BatchConfig, s.BatchPackets)
		}
	}()

	// handle packets
	for {
		batch, ok := <-s.BatchPackets
		if !ok {
			fmt.Println("Channel closed, exiting loop")
			break // Exit the loop when the channel is closed
		}

		// update client
		latestLightBlock, err := s.worker.UpdateCosmosClient(ctx, "groth16", int64(ctx.latestEthTimestamp.LatestUpdateHeight), "2/3")
		if err != nil {
			ctx.Logger.Println(fmt.Errorf("Failed to update cosmos light client: %s", err.Error()))
		}

		ctx.latestEthTimestamp.mtx.Lock()
		// update latest update time
		ctx.latestEthTimestamp.LatestUpdateTime = time.Now()
		// update latest trusted block height
		ctx.latestEthTimestamp.LatestUpdateHeight = uint64(latestLightBlock.BlockHeight)
		ctx.latestEthTimestamp.mtx.Unlock()

		// Check current Eth block timestamp for timeout comparisons
		ethHeader, err := ctx.EthClient().HeaderByNumber(context.Background(), nil)
		if err != nil {
			ctx.Logger.Println(fmt.Errorf("Failed to get eth block header: %s", err.Error()))
		}
		ethBlockTime := uint64(0)
		if ethHeader != nil {
			ethBlockTime = ethHeader.Time
		}

		// handle packets in batch
		for _, packet := range batch.Packets {
			switch packet.PacketType {
			case Send:
				// Skip packets that have already timed out
				if ethBlockTime > 0 && packet.Packet.TimeoutTimestamp > 0 && ethBlockTime >= packet.Packet.TimeoutTimestamp {
					log.Printf("[RecvPacket] Packet seq=%d timed out (timeout=%d <= eth_block_time=%d), skipping",
						packet.Packet.Sequence, packet.Packet.TimeoutTimestamp, ethBlockTime)
					continue
				}

				ibcPath := utils.IbcCommitmentPath(*packet.Packet, []byte{1})

				// target height are the latest block height
				value, proof, err := client.ProvePath(ctx.CosmosClient(), latestLightBlock.BlockHeight, ibcPath)
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("failed to prove path: %w", err))
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
					MembershipType: 0,
				}

				calldata, err := tendermintAbiJson.Pack("verifyMembership", membershipMsg)
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("Failed to abi encode verify msg: %s", err.Error()))
				}
				// Strip 4-byte function selector — ICS26Router does abi.decode, not a function call
				calldata = calldata[4:]

				payloads := make([]contractICS26Router.IICS26RouterMsgsPayload, 0, len(packet.Packet.Payloads))
				for _, p := range packet.Packet.Payloads {
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
						Sequence:         packet.Packet.Sequence,
						SourceClient:     packet.Packet.SourceClient,
						DestClient:       packet.Packet.DestinationClient,
						TimeoutTimestamp: packet.Packet.TimeoutTimestamp,
						Payloads:         payloads,
					},
					MembershipMsg: calldata,
				}

				s.txHandler.SendEthTx(ctx, msgRecvPacket)
			case Ack:
				ibcPath := utils.IbcCommitmentPath(*packet.Packet, []byte{2})

				// target height are the latest block height
				value, proof, err := client.ProvePath(ctx.CosmosClient(), latestLightBlock.BlockHeight, ibcPath)
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
					MembershipType: 0,
				}

				calldata, err := tendermintAbiJson.Pack("verifyMembership", membershipMsg)
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("Failed to abi encode verify msg: %s", err.Error()))
				}
				// Strip 4-byte function selector — ICS26Router does abi.decode, not a function call
				calldata = calldata[4:]

				payloads := make([]contractICS26Router.IICS26RouterMsgsPayload, 0, len(packet.Packet.Payloads))
				for _, p := range packet.Packet.Payloads {
					payloads = append(payloads, contractICS26Router.IICS26RouterMsgsPayload{
						SourcePort: p.SourcePort,
						DestPort:   p.DestinationPort,
						Version:    p.Version,
						Encoding:   p.Encoding,
						Value:      p.Value,
					})
				}

				msgAckPacket := contractICS26Router.IICS26RouterMsgsMsgAckPacket{
					Packet: contractICS26Router.IICS26RouterMsgsPacket{
						Sequence:         packet.Packet.Sequence,
						SourceClient:     packet.Packet.SourceClient,
						DestClient:       packet.Packet.DestinationClient,
						TimeoutTimestamp: packet.Packet.TimeoutTimestamp,
						Payloads:         payloads,
					},

					MembershipMsg: calldata,
				}

				s.txHandler.SendEthTx(ctx, msgAckPacket)
			case Timeout:
				ibcPath := utils.IbcCommitmentPath(*packet.Packet, []byte{3})

				// target height are the latest block height
				value, proof, err := client.ProvePath(ctx.CosmosClient(), latestLightBlock.BlockHeight, ibcPath)
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("failed to prove path: %w", err))
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

				membershipMsg := tendermintContract.ILightClientMsgsMsgVerifyNonMembership{
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
					MembershipType: 0,
				}

				calldata, err := tendermintAbiJson.Pack("verifyMembership", membershipMsg)
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("Failed to abi encode verify msg: %s", err.Error()))
				}
				// Strip 4-byte function selector — ICS26Router does abi.decode, not a function call
				calldata = calldata[4:]

				payloads := make([]contractICS26Router.IICS26RouterMsgsPayload, 0, len(packet.Packet.Payloads))
				for _, p := range packet.Packet.Payloads {
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
						Sequence:         packet.Packet.Sequence,
						SourceClient:     packet.Packet.SourceClient,
						DestClient:       packet.Packet.DestinationClient,
						TimeoutTimestamp: packet.Packet.TimeoutTimestamp,
						Payloads:         payloads,
					},
					MembershipMsg: calldata,
				}

				s.txHandler.SendEthTx(ctx, msgRecvPacket)

			default:
				ctx.Logger.Println(fmt.Errorf("Invalid packet type"))
			}
		}

		defer ctx.StopClient()
	}
}
