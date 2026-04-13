package services

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"relayer/utils"
	"strconv"
	"time"

	tendermintContract "relayer/bindings/Groth16ICS07Tendermint"
	contractICS26Router "relayer/bindings/ICS26Router"
	client "relayer/client"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// read abi json file once in runtime
var tendermintAbiJson *abi.ABI
var initErr error

func init() {
	tendermintAbiJson, initErr = tendermintContract.ContractGroth16ICS07TendermintMetaData.GetAbi()
	if initErr != nil {
		log.Fatal(initErr)
	}
}

type TransactionHandler interface {
	CreateCosmosClientContract(ctx Context, clientState, consensusHash []byte) error
	CreateEthClient(ctx Context, clientState ibcexported.ClientState, consensusState ibcexported.ConsensusState) (string, error)
	SendEthTx(ctx Context, msg any) error
	SendCosmosTx(ctx Context, msg any) error
	SendCosmosTxBatch(ctx Context, msgs []any) error
	CosmosSignerAddress() (string, error)
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
	listener EventListener
	worker   *Worker

	ethConfig    Config
	cosmosConfig Config

	BatchPackets chan BatchPackets
	BatchBuilder *BatchBuilder
}

func New(eventListener EventListener, txHandler TransactionHandler, prover Prover, ethConfig, cosmosConfig Config) *Services {
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

func (s *Services) StartLoop(ctx Context) {
	go s.listener.SubscribeCosmos(ctx, s.BatchBuilder)
	go s.listener.SubscribeEth(ctx, s.BatchBuilder)

	// routinely run update client
	go func() {
		// TODO: Revisit routine scheduling strategy (interval/backoff/event-driven mix) to ensure this is optimal for production.
		routineInterval := 24 * time.Hour
		for {
			now := time.Now()
			// update client on Eth side routinely
			if ctx.latestEthTimestamp.LatestUpdateTime.Add(routineInterval).Before(now) {
				latestBlock, err := s.worker.UpdateCosmosClient(ctx, "groth16", int64(ctx.latestEthTimestamp.LatestUpdateHeight), "1/3")
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("Failed to update cosmos light client: %s", err.Error()))
					continue
				}

				// update latest update time
				ctx.latestEthTimestamp.mtx.Lock()
				ctx.latestEthTimestamp.LatestUpdateTime = time.Now()
				ctx.latestEthTimestamp.LatestUpdateHeight = uint64(latestBlock.BlockHeight)
				ctx.latestEthTimestamp.mtx.Unlock()
			}

			// update client on Cosmos side routinely
			if ctx.latestCosmosTimestamp.LatestUpdateTime.Add(routineInterval).Before(now) {
				s.worker.UpdateEthClient(ctx)

				// update latest update time
				ctx.latestCosmosTimestamp.mtx.Lock()
				ctx.latestCosmosTimestamp.LatestUpdateTime = now
				ctx.latestCosmosTimestamp.mtx.Unlock()

			}

			time.Sleep(time.Second)
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

		//TODO: the way better than wait
		log.Printf("[Listener] Waiting 2 blocks for packet commitment to be included in AppHash...")
		time.Sleep(6 * time.Second)

		// update client
		latestLightBlock, err := s.worker.UpdateCosmosClient(ctx, "groth16", int64(ctx.latestEthTimestamp.LatestUpdateHeight), "1/3")
		if err != nil {
			ctx.Logger.Println(fmt.Errorf("Failed to update cosmos light client: %s", err.Error()))
			continue
		}
		if latestLightBlock == nil {
			ctx.Logger.Println("Failed to update cosmos light client: latestLightBlock is nil")
			continue
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

				value, proof, err := client.ProvePath(ctx.CosmosClient(), latestLightBlock.BlockHeight, ibcPath)
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("failed to prove path: %w", err))
					continue
				}

				if len(value) == 0 {
					ctx.Logger.Println(fmt.Errorf("[RecvPacket] packet commitment empty at height %d, skipping seq=%d",
						latestLightBlock.BlockHeight, packet.Packet.Sequence))
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
					AppHash: utils.BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
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

				s.worker.TxHandler.SendEthTx(ctx, msgRecvPacket)
			case Ack:
				if len(packet.AckBytes) == 0 {
					ctx.Logger.Println(fmt.Errorf("acknowledgement bytes missing for packet seq=%d", packet.Packet.Sequence))
					continue
				}

				ibcPath := utils.IbcCommitmentPath(*packet.Packet, []byte{3})

				// target height are the latest block height
				value, proof, err := client.ProvePath(ctx.CosmosClient(), latestLightBlock.BlockHeight, ibcPath)
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("failed to prove path: %w", err))
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
					Acknowledgement: packet.AckBytes[0],
					MembershipMsg:   calldata,
				}

				s.worker.TxHandler.SendEthTx(ctx, msgAckPacket)
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

				s.worker.TxHandler.SendEthTx(ctx, msgRecvPacket)

			case WriteAck:
				signerAddr, err := s.worker.TxHandler.CosmosSignerAddress()
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("[WriteAck] failed to get cosmos signer: %w", err))
					continue
				}

				finalized := false
				for attempt := 0; attempt < 60; attempt++ {
					if attempt > 0 {
						time.Sleep(10 * time.Second)
					}
					finalityUpdate, err := client.GetFinalityUpdate(ctx.BeaconAPIURL())
					if err != nil {
						log.Printf("[WriteAck] failed to get finality update: %v", err)
						continue
					}
					execBlock, _ := strconv.ParseUint(finalityUpdate.FinalizedHeader.Execution.BlockNumber, 10, 64)
					if execBlock >= packet.BlockNumber {
						log.Printf("[WriteAck] Beacon finalized block %d >= event block %d", execBlock, packet.BlockNumber)
						finalized = true
						break
					}
					log.Printf("[WriteAck] Beacon finalized block %d < event block %d, waiting... (%d/60)",
						execBlock, packet.BlockNumber, attempt+1)
				}
				if !finalized {
					ctx.Logger.Println(fmt.Errorf("[WriteAck] beacon finality did not reach block %d after retries", packet.BlockNumber))
					continue
				}

				if err := s.worker.UpdateEthClient(ctx); err != nil {
					ctx.Logger.Println(fmt.Errorf("[WriteAck] failed to update ETH client: %w", err))
					continue
				}

				ethClientState, err := client.GetEthereumClientState(ctx.CosmosClient(), ctx.EthClientID())
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("[WriteAck] failed to get ETH client state: %w", err))
					continue
				}
				proofBlockNumber := ethClientState.LatestExecutionBlockNumber
				proofSlot := ethClientState.LatestSlot

				if proofBlockNumber < packet.BlockNumber {
					ctx.Logger.Println(fmt.Errorf("[WriteAck] ETH client at block %d still < event block %d after update, skipping",
						proofBlockNumber, packet.BlockNumber))
					continue
				}

				// ack_path = destClientID + [0x03] + sequence.to_be_bytes(8)
				seqBytes := make([]byte, 8)
				binary.BigEndian.PutUint64(seqBytes, packet.Packet.Sequence)
				ackPath := append([]byte(packet.Packet.DestinationClient), 0x03)
				ackPath = append(ackPath, seqBytes...)

				slot := ethcommon.HexToHash(ICS26_IBC_STORAGE_SLOT)
				proofBytes, err := client.GetEthMembershipProof(
					ctx.EthClient(), *ctx.RouterContract(), ackPath, slot, new(big.Int).SetUint64(proofBlockNumber))
				if err != nil {
					ctx.Logger.Println(fmt.Errorf("[WriteAck] failed to get ETH membership proof: %w", err))
					continue
				}

				ackMsg := &channeltypesv2.MsgAcknowledgement{
					Packet: *packet.Packet,
					Acknowledgement: channeltypesv2.Acknowledgement{
						AppAcknowledgements: packet.AckBytes,
					},
					ProofAcked:  proofBytes,
					ProofHeight: clienttypes.Height{RevisionNumber: 0, RevisionHeight: proofSlot},
					Signer:      signerAddr,
				}
				if err := s.worker.TxHandler.SendCosmosTx(ctx, ackMsg); err != nil {
					ctx.Logger.Println(fmt.Errorf("[WriteAck] failed to send MsgAcknowledgement: %w", err))
				}

			default:
				ctx.Logger.Println(fmt.Errorf("Invalid packet type"))
			}
		}

		defer ctx.StopClient()
	}
}
