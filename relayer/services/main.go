package services

import (
	"context"
	"encoding/binary"
	"log"
	"math/big"
	"relayer/utils"
	"strconv"
	"time"

	tendermintContract "relayer/bindings/Groth16ICS07Tendermint"
	contractICS26Router "relayer/bindings/ICS26Router"
	client "relayer/client"
	"relayer/prover"

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
	CreateCosmosClientContract(ctx Context, clientState, consensusHash []byte) (ethcommon.Address, error)
	CreateEthClient(ctx Context, clientState ibcexported.ClientState, consensusState ibcexported.ConsensusState) (string, error)
	SendEthTx(ctx Context, msg any) error
	SendCosmosTx(ctx Context, msg any) error
	SendCosmosTxBatch(ctx Context, msgs []any) error
	CosmosSignerAddress() (string, error)
}

type Prover interface {
	GenerateProof(sigs []prover.ValidatorSignature) (
		bucket int,
		paddedSigs []prover.ValidatorSignature,
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
					log.Printf("[Routine] Failed to update cosmos light client: %v", err)
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
			log.Println("[StartLoop] Batch channel closed, exiting loop")
			break // Exit the loop when the channel is closed
		}

		log.Printf("[StartLoop] Received batch: %d packets", len(batch.Packets))

		//TODO: the way better than wait
		log.Printf("[StartLoop] Waiting 2 blocks for packet commitment to be included in AppHash...")
		time.Sleep(6 * time.Second)

		// update client
		latestLightBlock, err := s.worker.UpdateCosmosClient(ctx, "groth16", int64(ctx.latestEthTimestamp.LatestUpdateHeight), "1/3")
		if err != nil {
			log.Printf("[StartLoop] Failed to update cosmos light client: %v", err)
			continue
		}
		if latestLightBlock == nil {
			log.Printf("[StartLoop] Failed to update cosmos light client: latestLightBlock is nil")
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
			log.Printf("[StartLoop] Failed to get eth block header: %v", err)
		}
		ethBlockTime := uint64(0)
		if ethHeader != nil {
			ethBlockTime = ethHeader.Time
		}

		// handle packets in batch
		for _, packet := range batch.Packets {
			switch packet.PacketType {
			case Send:
				// If the packet has timed out on ETH, submit MsgTimeout back to Cosmos to release escrow.
				if ethBlockTime > 0 && packet.Packet.TimeoutTimestamp > 0 && ethBlockTime >= packet.Packet.TimeoutTimestamp {
					log.Printf("[CosmosTimeout] seq=%d: timed out (eth_block_time=%d >= timeout=%d), submitting MsgTimeout to Cosmos",
						packet.Packet.Sequence, ethBlockTime, packet.Packet.TimeoutTimestamp)

					signerAddr, err := s.worker.TxHandler.CosmosSignerAddress()
					if err != nil {
						log.Printf("[CosmosTimeout] seq=%d: failed to get cosmos signer: %v", packet.Packet.Sequence, err)
						continue
					}

					// Wait for beacon to finalize an ETH block whose timestamp >= timeout timestamp.
					// Cosmos IBC verifies: proof_block.timestamp >= timeout_timestamp.
					finalized := false
					for attempt := 0; attempt < 60; attempt++ {
						if attempt > 0 {
							time.Sleep(10 * time.Second)
						}
						finalityUpdate, err := client.GetFinalityUpdate(ctx.BeaconAPIURL())
						if err != nil {
							log.Printf("[CosmosTimeout] seq=%d: failed to get finality update: %v", packet.Packet.Sequence, err)
							continue
						}
						execTimestamp, _ := strconv.ParseUint(finalityUpdate.FinalizedHeader.Execution.Timestamp, 10, 64)
						if execTimestamp >= packet.Packet.TimeoutTimestamp {
							log.Printf("[CosmosTimeout] seq=%d: finalized ETH block timestamp %d >= timeout %d",
								packet.Packet.Sequence, execTimestamp, packet.Packet.TimeoutTimestamp)
							finalized = true
							break
						}
						log.Printf("[CosmosTimeout] seq=%d: finalized ETH timestamp %d < timeout %d, waiting... (%d/60)",
							packet.Packet.Sequence, execTimestamp, packet.Packet.TimeoutTimestamp, attempt+1)
					}
					if !finalized {
						log.Printf("[CosmosTimeout] seq=%d: ETH finality did not reach timeout timestamp after 60 retries", packet.Packet.Sequence)
						continue
					}

					if err := s.worker.UpdateEthClient(ctx); err != nil {
						log.Printf("[CosmosTimeout] seq=%d: failed to update ETH client: %v", packet.Packet.Sequence, err)
						continue
					}

					ethClientState, err := client.GetEthereumClientState(ctx.CosmosClient(), ctx.EthClientID())
					if err != nil {
						log.Printf("[CosmosTimeout] seq=%d: failed to get ETH client state: %v", packet.Packet.Sequence, err)
						continue
					}
					proofBlockNumber := ethClientState.LatestExecutionBlockNumber
					proofSlot := ethClientState.LatestSlot

					// receipt path = destClientID + [0x02] + sequence.to_be_bytes(8)
					seqBytes := make([]byte, 8)
					binary.BigEndian.PutUint64(seqBytes, packet.Packet.Sequence)
					receiptPath := append([]byte(packet.Packet.DestinationClient), 0x02)
					receiptPath = append(receiptPath, seqBytes...)

					slot := ethcommon.HexToHash(ICS26_IBC_STORAGE_SLOT)
					proofBytes, err := client.GetEthMembershipProof(
						ctx.EthClient(), *ctx.RouterContract(), receiptPath, slot, new(big.Int).SetUint64(proofBlockNumber))
					if err != nil {
						log.Printf("[CosmosTimeout] seq=%d: failed to get ETH non-membership proof: %v", packet.Packet.Sequence, err)
						continue
					}

					timeoutMsg := &channeltypesv2.MsgTimeout{
						Packet:          *packet.Packet,
						ProofUnreceived: proofBytes,
						ProofHeight:     clienttypes.Height{RevisionNumber: 0, RevisionHeight: proofSlot},
						Signer:          signerAddr,
					}
					if err := s.worker.TxHandler.SendCosmosTx(ctx, timeoutMsg); err != nil {
						log.Printf("[CosmosTimeout] seq=%d: failed to send MsgTimeout: %v", packet.Packet.Sequence, err)
						continue
					}
					log.Printf("[CosmosTimeout] seq=%d: relay completed", packet.Packet.Sequence)
					continue
				}

				ibcPath := utils.IbcCommitmentPath(*packet.Packet, []byte{1})

				value, proof, err := client.ProvePath(ctx.CosmosClient(), latestLightBlock.BlockHeight, ibcPath)
				if err != nil {
					log.Printf("[RecvPacket] seq=%d: failed to prove path at height=%d: %v",
						packet.Packet.Sequence, latestLightBlock.BlockHeight, err)
					continue
				}

				if len(value) == 0 {
					log.Printf("[RecvPacket] seq=%d: packet commitment empty at height=%d, skipping",
						packet.Packet.Sequence, latestLightBlock.BlockHeight)
					continue
				}

				merkleProof := tendermintContract.IMembershipMsgsMerkleProof{
					Proofs: []tendermintContract.IMembershipMsgsCommitmentProof{},
				}
				for _, p := range proof.Proofs {
					commitmentProof, err := client.ParseCommitmentProof(p)
					if err != nil {
						log.Printf("[RecvPacket] seq=%d: failed to parse commitment proof: %v",
							packet.Packet.Sequence, err)
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
							Value: value,
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
					log.Printf("[RecvPacket] seq=%d: failed to ABI encode verifyMembership: %v",
						packet.Packet.Sequence, err)
					continue
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

				if err := s.worker.TxHandler.SendEthTx(ctx, msgRecvPacket); err != nil {
					log.Printf("[RecvPacket] seq=%d: SendEthTx failed: %v", packet.Packet.Sequence, err)
					continue
				}
				log.Printf("[RecvPacket] seq=%d: relay completed", packet.Packet.Sequence)
			case Ack:
				if len(packet.AckBytes) == 0 {
					log.Printf("[AckPacket] seq=%d: acknowledgement bytes missing, skipping", packet.Packet.Sequence)
					continue
				}

				ibcPath := utils.IbcCommitmentPath(*packet.Packet, []byte{3})

				// target height are the latest block height
				value, proof, err := client.ProvePath(ctx.CosmosClient(), latestLightBlock.BlockHeight, ibcPath)
				if err != nil {
					log.Printf("[AckPacket] seq=%d: failed to prove path at height=%d: %v",
						packet.Packet.Sequence, latestLightBlock.BlockHeight, err)
					continue
				}

				merkleProof := tendermintContract.IMembershipMsgsMerkleProof{
					Proofs: []tendermintContract.IMembershipMsgsCommitmentProof{},
				}
				for _, p := range proof.Proofs {
					commitmentProof, err := client.ParseCommitmentProof(p)
					if err != nil {
						log.Printf("[AckPacket] seq=%d: failed to parse commitment proof: %v",
							packet.Packet.Sequence, err)
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
							Value: value,
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
					log.Printf("[AckPacket] seq=%d: failed to ABI encode verifyMembership: %v",
						packet.Packet.Sequence, err)
					continue
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

				if err := s.worker.TxHandler.SendEthTx(ctx, msgAckPacket); err != nil {
					log.Printf("[AckPacket] seq=%d: SendEthTx failed: %v", packet.Packet.Sequence, err)
					continue
				}
				log.Printf("[AckPacket] seq=%d: relay completed", packet.Packet.Sequence)
			case Timeout:
				ibcPath := utils.IbcCommitmentPath(*packet.Packet, []byte{2})

				// target height are the latest block height
				value, proof, err := client.ProvePath(ctx.CosmosClient(), latestLightBlock.BlockHeight, ibcPath)
				if err != nil {
					log.Printf("[Timeout] seq=%d: failed to prove path at height=%d: %v",
						packet.Packet.Sequence, latestLightBlock.BlockHeight, err)
					continue
				}

				merkleProof := tendermintContract.IMembershipMsgsMerkleProof{
					Proofs: []tendermintContract.IMembershipMsgsCommitmentProof{},
				}
				for _, p := range proof.Proofs {
					commitmentProof, err := client.ParseCommitmentProof(p)
					if err != nil {
						log.Printf("[Timeout] seq=%d: failed to parse commitment proof: %v",
							packet.Packet.Sequence, err)
					}
					merkleProof.Proofs = append(merkleProof.Proofs, *commitmentProof)
				}

				nonMembershipMsg := tendermintContract.ILightClientMsgsMsgVerifyNonMembership{
					Height: tendermintContract.IICS02ClientMsgsHeight{
						RevisionHeight: uint64(latestLightBlock.BlockHeight),
						RevisionNumber: 0,
					},
					KvPairs: []tendermintContract.IMembershipMsgsKVPair{
						{
							Path:  ibcPath,
							Value: value,
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

				calldata, err := tendermintAbiJson.Pack("verifyNonMembership", nonMembershipMsg)
				if err != nil {
					log.Printf("[Timeout] seq=%d: failed to ABI encode verifyNonMembership: %v",
						packet.Packet.Sequence, err)
					continue
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

				msgTimeoutPacket := contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket{
					Packet: contractICS26Router.IICS26RouterMsgsPacket{
						Sequence:         packet.Packet.Sequence,
						SourceClient:     packet.Packet.SourceClient,
						DestClient:       packet.Packet.DestinationClient,
						TimeoutTimestamp: packet.Packet.TimeoutTimestamp,
						Payloads:         payloads,
					},
					NonMembershipMsg: calldata,
				}

				if err := s.worker.TxHandler.SendEthTx(ctx, msgTimeoutPacket); err != nil {
					log.Printf("[Timeout] seq=%d: SendEthTx failed: %v", packet.Packet.Sequence, err)
					continue
				}
				log.Printf("[Timeout] seq=%d: relay completed", packet.Packet.Sequence)

			case WriteAck:
				signerAddr, err := s.worker.TxHandler.CosmosSignerAddress()
				if err != nil {
					log.Printf("[WriteAck] seq=%d: failed to get cosmos signer: %v", packet.Packet.Sequence, err)
					continue
				}

				log.Printf("[WriteAck] seq=%d: waiting for beacon finality at block %d",
					packet.Packet.Sequence, packet.BlockNumber)
				finalized := false
				for attempt := 0; attempt < 60; attempt++ {
					if attempt > 0 {
						time.Sleep(10 * time.Second)
					}
					finalityUpdate, err := client.GetFinalityUpdate(ctx.BeaconAPIURL())
					if err != nil {
						log.Printf("[WriteAck] seq=%d: failed to get finality update: %v", packet.Packet.Sequence, err)
						continue
					}
					execBlock, _ := strconv.ParseUint(finalityUpdate.FinalizedHeader.Execution.BlockNumber, 10, 64)
					if execBlock >= packet.BlockNumber {
						log.Printf("[WriteAck] seq=%d: beacon finalized block %d >= event block %d",
							packet.Packet.Sequence, execBlock, packet.BlockNumber)
						finalized = true
						break
					}
					log.Printf("[WriteAck] seq=%d: beacon finalized block %d < event block %d, waiting... (%d/60)",
						packet.Packet.Sequence, execBlock, packet.BlockNumber, attempt+1)
				}
				if !finalized {
					log.Printf("[WriteAck] seq=%d: beacon finality did not reach block %d after 60 retries",
						packet.Packet.Sequence, packet.BlockNumber)
					continue
				}

				if err := s.worker.UpdateEthClient(ctx); err != nil {
					log.Printf("[WriteAck] seq=%d: failed to update ETH client: %v", packet.Packet.Sequence, err)
					continue
				}

				ethClientState, err := client.GetEthereumClientState(ctx.CosmosClient(), ctx.EthClientID())
				if err != nil {
					log.Printf("[WriteAck] seq=%d: failed to get ETH client state: %v", packet.Packet.Sequence, err)
					continue
				}
				proofBlockNumber := ethClientState.LatestExecutionBlockNumber
				proofSlot := ethClientState.LatestSlot

				if proofBlockNumber < packet.BlockNumber {
					log.Printf("[WriteAck] seq=%d: ETH client at block %d still < event block %d after update, skipping",
						packet.Packet.Sequence, proofBlockNumber, packet.BlockNumber)
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
					log.Printf("[WriteAck] seq=%d: failed to get ETH membership proof: %v", packet.Packet.Sequence, err)
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
					log.Printf("[WriteAck] seq=%d: failed to send MsgAcknowledgement: %v", packet.Packet.Sequence, err)
					continue
				}
				log.Printf("[WriteAck] seq=%d: relay completed", packet.Packet.Sequence)

			default:
				log.Printf("[StartLoop] Unknown packet type: %d (seq=%d)", packet.PacketType, packet.Packet.Sequence)
			}
		}

		defer ctx.StopClient()
	}
}
