package subscriber

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/services"
	"relayer/utils"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogo/protobuf/proto"
)

const COMETBFT_SEND_PACKET_EVENT = "tm.event = 'Tx' AND message.action = '/ibc.applications.transfer.v1.MsgTransfer'"
const COMETBFT_WRITE_ACK_PACKET_EVENT = "tm.event = 'Tx' AND message.action = '/ibc.core.channel.v2.MsgRecvPacket'"
const COMETBFT_TIMEOUT_PACKET_EVENT = "tm.event = 'Tx' AND message.action = '/ibc.applications.transfer.v1.MsgTimeout'"

const EVENT_SEND_PACKET_FIELD = "send_packet.encoded_packet_hex"
const EVENT_WRITE_ACK_PACKET_FIELD = "write_acknowledgement.encoded_packet_hex"
const EVENT_ACKNOWLEDGEMENT_FIELD = "write_acknowledgement.encoded_acknowledgement_hex"
const EVENT_TIMEOUT_PACKET_FIELD = "timeout_packet.encoded_packet_hex"

const ethStartupRecoveryLookbackEnv = "ETH_STARTUP_LOOKBACK_BLOCKS"
const defaultEthStartupRecoveryLookbackBlocks uint64 = 256
const ethSubscriptionReconnectDelay = 2 * time.Second

type Subscriber struct {
}

func NewSubscriber() *Subscriber {
	return &Subscriber{}
}

func (s *Subscriber) SubscribeCosmos(ctx services.Context, batchBuilder *services.BatchBuilder) {
	c, cancel := context.WithCancel(context.Background())
	defer cancel()

	sendPacketSub, err := ctx.CosmosClient().WSEvents.Subscribe(context.Background(), "", COMETBFT_SEND_PACKET_EVENT)
	if err != nil {
		ctx.Logger.Printf("[SubscribeCosmos] Failed to subscribe to send_packet events: %v", err)
	}
	ackPacketSub, err := ctx.CosmosClient().WSEvents.Subscribe(context.Background(), "", COMETBFT_WRITE_ACK_PACKET_EVENT)
	if err != nil {
		ctx.Logger.Printf("[SubscribeCosmos] Failed to subscribe to write_acknowledgement events: %v", err)
	}
	timeoutPacketSub, err := ctx.CosmosClient().WSEvents.Subscribe(context.Background(), "", COMETBFT_TIMEOUT_PACKET_EVENT)
	if err != nil {
		ctx.Logger.Printf("[SubscribeCosmos] Failed to subscribe to timeout_packet events: %v", err)
	}
	ctx.Logger.Println("[SubscribeCosmos] Successfully subscribed to CometBFT events")
	defer ctx.CosmosClient().UnsubscribeAll(context.Background(), "")

	for {
		select {
		case e := <-sendPacketSub:
			sendPacketEvent := e.Events[EVENT_SEND_PACKET_FIELD]
			if sendPacketEvent == nil {
				continue
			}

			packetEncodedStr := sendPacketEvent[0]
			packetBytes, err := hex.DecodeString(packetEncodedStr)
			if err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] send_packet: failed to decode hex: %v", err)
				continue
			}

			var packet channeltypesv2.Packet
			err = proto.Unmarshal(packetBytes, &packet)
			if err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] send_packet: failed to unmarshal: %v", err)
				continue
			}

			ctx.Logger.Printf("[SubscribeCosmos] send_packet received: seq=%d src=%s",
				packet.Sequence, packet.SourceClient)
			batchBuilder.AddCosmos(services.CosmosPacket{
				Type:   services.CosmosSend,
				Packet: &packet,
			})
		case e := <-ackPacketSub:
			ackPacketEvent := e.Events[EVENT_WRITE_ACK_PACKET_FIELD]
			ackEvent := e.Events[EVENT_ACKNOWLEDGEMENT_FIELD]
			if len(ackPacketEvent) == 0 || len(ackEvent) == 0 {
				continue
			}

			packetEncodedStr := ackPacketEvent[0]
			packetBytes, err := hex.DecodeString(packetEncodedStr)
			if err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] write_ack: failed to decode packet hex: %v", err)
				continue
			}

			var packet channeltypesv2.Packet
			err = proto.Unmarshal(packetBytes, &packet)
			if err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] write_ack: failed to unmarshal packet: %v", err)
				continue
			}

			ackBytes, err := hex.DecodeString(ackEvent[0])
			if err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] write_ack seq=%d: failed to decode ack hex: %v",
					packet.Sequence, err)
				continue
			}

			var acknowledgement channeltypesv2.Acknowledgement
			err = proto.Unmarshal(ackBytes, &acknowledgement)
			if err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] write_ack seq=%d: failed to unmarshal ack: %v",
					packet.Sequence, err)
				continue
			}
			if len(acknowledgement.AppAcknowledgements) == 0 {
				ctx.Logger.Printf("[SubscribeCosmos] write_ack seq=%d: missing app acknowledgements", packet.Sequence)
				continue
			}

			ctx.Logger.Printf("[SubscribeCosmos] write_ack received: seq=%d src=%s",
				packet.Sequence, packet.SourceClient)
			batchBuilder.AddCosmos(services.CosmosPacket{
				Type:     services.CosmosAck,
				Packet:   &packet,
				AckBytes: acknowledgement.AppAcknowledgements,
			})
		case e := <-timeoutPacketSub:
			timeoutPacketEvent := e.Events[EVENT_TIMEOUT_PACKET_FIELD]
			if timeoutPacketEvent == nil {
				continue
			}

			packetEncodedStr := timeoutPacketEvent[0]
			packetBytes, err := hex.DecodeString(packetEncodedStr)
			if err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] timeout: failed to decode hex: %v", err)
				continue
			}

			var packet channeltypesv2.Packet
			err = proto.Unmarshal(packetBytes, &packet)
			if err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] timeout: failed to unmarshal: %v", err)
				continue
			}

			ctx.Logger.Printf("[SubscribeCosmos] timeout received: seq=%d src=%s",
				packet.Sequence, packet.SourceClient)
			batchBuilder.AddCosmos(services.CosmosPacket{
				Type:   services.CosmosTimeout,
				Packet: &packet,
			})
		case <-c.Done():
			return
		}
	}
}

// EthPacketToCosmosPacket converts an Ethereum ICS26Router packet to a Cosmos IBC v2 packet
func EthPacketToCosmosPacket(ethPacket contractICS26Router.IICS26RouterMsgsPacket, sequence *big.Int) channeltypesv2.Packet {
	var payloads []channeltypesv2.Payload
	for _, p := range ethPacket.Payloads {
		payloads = append(payloads, channeltypesv2.Payload{
			SourcePort:      p.SourcePort,
			DestinationPort: p.DestPort,
			Version:         p.Version,
			Encoding:        p.Encoding,
			Value:           p.Value,
		})
	}

	return channeltypesv2.Packet{
		Sequence:          sequence.Uint64(),
		SourceClient:      ethPacket.SourceClient,
		DestinationClient: ethPacket.DestClient,
		TimeoutTimestamp:  ethPacket.TimeoutTimestamp,
		Payloads:          payloads,
	}
}

func ethStartupRecoveryLookbackBlocksFromEnv(raw string) uint64 {
	if raw == "" {
		return defaultEthStartupRecoveryLookbackBlocks
	}

	lookback, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return defaultEthStartupRecoveryLookbackBlocks
	}

	return lookback
}

func ethStartupRecoveryLookbackBlocks() uint64 {
	return ethStartupRecoveryLookbackBlocksFromEnv(os.Getenv(ethStartupRecoveryLookbackEnv))
}

func ethStartupRecoveryStartBlock(latestBlock, lookback uint64) uint64 {
	if lookback >= latestBlock {
		return 0
	}

	return latestBlock - lookback
}

func enqueueEthWriteAcknowledgement(batchBuilder *services.BatchBuilder, ev *contractICS26Router.ContractICS26RouterWriteAcknowledgement) {
	cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
	batchBuilder.AddEth(services.EthPacket{
		Type:        services.EthWriteAck,
		Packet:      &cosmosPacket,
		AckBytes:    ev.Acknowledgements,
		BlockNumber: ev.Raw.BlockNumber,
	})
}

func hasPendingCosmosPacketCommitment(ctx services.Context, packet channeltypesv2.Packet) (bool, error) {
	path := utils.IbcCommitmentPath(packet, []byte{1})
	queryPath := fmt.Sprintf("store/%s/key", string(path[0]))
	request := path[1]

	result, err := ctx.CosmosClient().ABCIQuery(context.Background(), queryPath, request)
	if err != nil {
		return false, fmt.Errorf("ABCI query failed: %w", err)
	}
	if result.Response.Code != 0 {
		return false, fmt.Errorf("query failed with code %d: %s", result.Response.Code, result.Response.Log)
	}

	return len(result.Response.Value) > 0, nil
}

func recoverEthWriteAcknowledgements(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	filterer *contractICS26Router.ContractICS26RouterFilterer,
	startBlock uint64,
	endBlock uint64,
) {
	if endBlock < startBlock {
		return
	}

	filterOpts := &bind.FilterOpts{
		Start:   startBlock,
		End:     &endBlock,
		Context: context.Background(),
	}
	iter, err := filterer.FilterWriteAcknowledgement(filterOpts, nil, nil)
	if err != nil {
		ctx.Logger.Printf("[SubscribeEth] startup recovery: failed to filter WriteAcknowledgement logs in [%d,%d]: %v",
			startBlock, endBlock, err)
		return
	}
	defer iter.Close()

	var recoveredCount uint64
	var skippedCount uint64
	for iter.Next() {
		ev := iter.Event
		cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)

		pending, err := hasPendingCosmosPacketCommitment(ctx, cosmosPacket)
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] startup recovery: seq=%d failed to check Cosmos packet commitment: %v",
				cosmosPacket.Sequence, err)
			continue
		}
		if !pending {
			skippedCount++
			ctx.Logger.Printf("[SubscribeEth] startup recovery: seq=%d already cleared on Cosmos, skipping historical WriteAcknowledgement from ETH block %d",
				cosmosPacket.Sequence, ev.Raw.BlockNumber)
			continue
		}

		ctx.Logger.Printf("[SubscribeEth] startup recovery: recovered WriteAcknowledgement seq=%d from ETH block %d",
			cosmosPacket.Sequence, ev.Raw.BlockNumber)
		enqueueEthWriteAcknowledgement(batchBuilder, ev)
		recoveredCount++
	}

	if err := iter.Error(); err != nil {
		ctx.Logger.Printf("[SubscribeEth] startup recovery: iterator error in [%d,%d]: %v",
			startBlock, endBlock, err)
	}

	ctx.Logger.Printf("[SubscribeEth] startup recovery complete: scanned [%d,%d], recovered=%d skipped=%d",
		startBlock, endBlock, recoveredCount, skippedCount)
}

func advanceRecoveryStart(nextRecoveryStartBlock *uint64, candidate uint64) {
	if candidate > *nextRecoveryStartBlock {
		*nextRecoveryStartBlock = candidate
	}
}

func (s *Subscriber) subscribeEthOnce(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	watchClient *ethclient.Client,
	watchStartBlock uint64,
	nextWriteAckRecoveryStartBlock *uint64,
) error {
	watchFilterer, err := contractICS26Router.NewContractICS26RouterFilterer(*ctx.RouterContract(), watchClient)
	if err != nil {
		return fmt.Errorf("failed to create ICS26Router watch filterer instance: %w", err)
	}

	sendPacketCh := make(chan *contractICS26Router.ContractICS26RouterSendPacket)
	writeAckCh := make(chan *contractICS26Router.ContractICS26RouterWriteAcknowledgement)
	ackPacketCh := make(chan *contractICS26Router.ContractICS26RouterAckPacket)
	timeoutPacketCh := make(chan *contractICS26Router.ContractICS26RouterTimeoutPacket)

	watchOpts := &bind.WatchOpts{Start: &watchStartBlock, Context: context.Background()}

	sendPacketSub, err := watchFilterer.WatchSendPacket(watchOpts, sendPacketCh, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to SendPacket events: %w", err)
	}
	defer sendPacketSub.Unsubscribe()

	writeAckSub, err := watchFilterer.WatchWriteAcknowledgement(watchOpts, writeAckCh, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to WriteAcknowledgement events: %w", err)
	}
	defer writeAckSub.Unsubscribe()

	ackPacketSub, err := watchFilterer.WatchAckPacket(watchOpts, ackPacketCh, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AckPacket events: %w", err)
	}
	defer ackPacketSub.Unsubscribe()

	timeoutPacketSub, err := watchFilterer.WatchTimeoutPacket(watchOpts, timeoutPacketCh, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to TimeoutPacket events: %w", err)
	}
	defer timeoutPacketSub.Unsubscribe()

	ctx.Logger.Printf("[SubscribeEth] Successfully subscribed to ICS26Router events from block %d", watchStartBlock)

	for {
		select {
		case ev := <-sendPacketCh:
			ctx.Logger.Printf("SendPacket event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
			batchBuilder.AddEth(services.EthPacket{
				Type:        services.EthSend,
				Packet:      &cosmosPacket,
				BlockNumber: ev.Raw.BlockNumber,
			})

		case ev := <-writeAckCh:
			ctx.Logger.Printf("WriteAcknowledgement event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			advanceRecoveryStart(nextWriteAckRecoveryStartBlock, ev.Raw.BlockNumber+1)
			enqueueEthWriteAcknowledgement(batchBuilder, ev)

		case ev := <-ackPacketCh:
			ctx.Logger.Printf("AckPacket event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
			batchBuilder.AddEth(services.EthPacket{
				Type:     services.EthAck,
				Packet:   &cosmosPacket,
				AckBytes: [][]byte{ev.Acknowledgement},
			})

		case ev := <-timeoutPacketCh:
			ctx.Logger.Printf("TimeoutPacket event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
			batchBuilder.AddEth(services.EthPacket{
				Type:   services.EthTimeout,
				Packet: &cosmosPacket,
			})

		case err := <-sendPacketSub.Err():
			return fmt.Errorf("SendPacket subscription error: %w", err)

		case err := <-writeAckSub.Err():
			return fmt.Errorf("WriteAcknowledgement subscription error: %w", err)

		case err := <-ackPacketSub.Err():
			return fmt.Errorf("AckPacket subscription error: %w", err)

		case err := <-timeoutPacketSub.Err():
			return fmt.Errorf("TimeoutPacket subscription error: %w", err)
		}
	}
}

// SubscribeEth subscribes to Ethereum events from the ICS26Router contract
func (s *Subscriber) SubscribeEth(ctx services.Context, batchBuilder *services.BatchBuilder) {
	if ctx.EthWsURL() == "" {
		ctx.Logger.Printf("Failed to subscribe to Ethereum events: eth websocket URL is not configured")
		return
	}

	recoveryFilterer, err := contractICS26Router.NewContractICS26RouterFilterer(*ctx.RouterContract(), ctx.EthClient())
	if err != nil {
		ctx.Logger.Printf("Failed to create ICS26Router recovery filterer instance: %v", err)
		return
	}

	lookback := ethStartupRecoveryLookbackBlocks()
	var nextWriteAckRecoveryStartBlock uint64

	for {
		latestBlock, err := ctx.EthClient().BlockNumber(context.Background())
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] Failed to get latest Ethereum block before subscription: %v", err)
			time.Sleep(ethSubscriptionReconnectDelay)
			continue
		}

		if nextWriteAckRecoveryStartBlock == 0 {
			nextWriteAckRecoveryStartBlock = ethStartupRecoveryStartBlock(latestBlock, lookback)
		}

		if latestBlock >= nextWriteAckRecoveryStartBlock {
			ctx.Logger.Printf("[SubscribeEth] recovery scanning WriteAcknowledgement logs in [%d,%d]",
				nextWriteAckRecoveryStartBlock, latestBlock)
			recoverEthWriteAcknowledgements(ctx, batchBuilder, recoveryFilterer, nextWriteAckRecoveryStartBlock, latestBlock)
		}

		watchStartBlock := latestBlock + 1
		nextWriteAckRecoveryStartBlock = watchStartBlock

		watchClient, err := ethclient.DialContext(context.Background(), ctx.EthWsURL())
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] Failed to connect to Ethereum WS at %s: %v", ctx.EthWsURL(), err)
			time.Sleep(ethSubscriptionReconnectDelay)
			continue
		}

		err = s.subscribeEthOnce(ctx, batchBuilder, watchClient, watchStartBlock, &nextWriteAckRecoveryStartBlock)
		watchClient.Close()
		ctx.Logger.Printf("[SubscribeEth] Subscription loop ended: %v", err)
		time.Sleep(ethSubscriptionReconnectDelay)
	}
}
