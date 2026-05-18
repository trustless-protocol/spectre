package subscriber

import (
	"context"
	"encoding/hex"
	"math/big"

	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/services"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gogo/protobuf/proto"
)

const COMETBFT_SEND_PACKET_EVENT = "tm.event = 'Tx' AND message.action = '/ibc.core.channel.v2.MsgSendPacket'"
const COMETBFT_WRITE_ACK_PACKET_EVENT = "tm.event = 'Tx' AND message.action = '/ibc.core.channel.v2.MsgRecvPacket'"
const COMETBFT_TIMEOUT_PACKET_EVENT = "tm.event = 'Tx' AND message.action = '/ibc.core.channel.v2.MsgTimeout'"

const EVENT_SEND_PACKET_FIELD = "send_packet.encoded_packet_hex"
const EVENT_WRITE_ACK_PACKET_FIELD = "write_acknowledgement.encoded_packet_hex"
const EVENT_ACKNOWLEDGEMENT_FIELD = "write_acknowledgement.encoded_acknowledgement_hex"
const EVENT_TIMEOUT_PACKET_FIELD = "timeout_packet.encoded_packet_hex"

// normalizeTimeoutSeconds converts IBC v2 timeout timestamps from nanoseconds to seconds.
// ibc-go stores TimeoutTimestamp in nanoseconds, but the ETH side uses seconds.
// Values <= 1e12 are assumed to already be in seconds (year ~33658 CE in seconds).
func normalizeTimeoutSeconds(ts uint64) uint64 {
	if ts > 1e12 {
		return ts / 1e9
	}
	return ts
}

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
			packet.TimeoutTimestamp = normalizeTimeoutSeconds(packet.TimeoutTimestamp)
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
			packet.TimeoutTimestamp = normalizeTimeoutSeconds(packet.TimeoutTimestamp)
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
			packet.TimeoutTimestamp = normalizeTimeoutSeconds(packet.TimeoutTimestamp)
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

// SubscribeEth subscribes to Ethereum events from the ICS26Router contract
func (s *Subscriber) SubscribeEth(ctx services.Context, batchBuilder *services.BatchBuilder) {
	filterer, err := contractICS26Router.NewContractICS26RouterFilterer(*ctx.RouterContract(), ctx.EthWsClient())
	if err != nil {
		ctx.Logger.Printf("Failed to create ICS26Router filterer instance: %v", err)
		return
	}

	// Create event channels for each event type
	sendPacketCh := make(chan *contractICS26Router.ContractICS26RouterSendPacket)
	writeAckCh := make(chan *contractICS26Router.ContractICS26RouterWriteAcknowledgement)
	ackPacketCh := make(chan *contractICS26Router.ContractICS26RouterAckPacket)
	timeoutPacketCh := make(chan *contractICS26Router.ContractICS26RouterTimeoutPacket)

	// Set up watch options (nil for all events, no filtering by clientId or sequence)
	watchOpts := &bind.WatchOpts{Context: context.Background()}
	// Subscribe to SendPacket events
	sendPacketSub, err := filterer.WatchSendPacket(watchOpts, sendPacketCh, nil, nil)
	if err != nil {
		ctx.Logger.Printf("Failed to subscribe to SendPacket events: %v", err)
		return
	}
	defer sendPacketSub.Unsubscribe()

	// Subscribe to WriteAcknowledgement events
	writeAckSub, err := filterer.WatchWriteAcknowledgement(watchOpts, writeAckCh, nil, nil)
	if err != nil {
		ctx.Logger.Printf("Failed to subscribe to WriteAcknowledgement events: %v", err)
		return
	}
	defer writeAckSub.Unsubscribe()

	// Subscribe to AckPacket events
	ackPacketSub, err := filterer.WatchAckPacket(watchOpts, ackPacketCh, nil, nil)
	if err != nil {
		ctx.Logger.Printf("Failed to subscribe to AckPacket events: %v", err)
		return
	}
	defer ackPacketSub.Unsubscribe()

	// Subscribe to TimeoutPacket events
	timeoutPacketSub, err := filterer.WatchTimeoutPacket(watchOpts, timeoutPacketCh, nil, nil)
	if err != nil {
		ctx.Logger.Printf("Failed to subscribe to TimeoutPacket events: %v", err)
		return
	}
	defer timeoutPacketSub.Unsubscribe()

	ctx.Logger.Println("Successfully subscribed to ICS26Router events")

	// Event loop to handle incoming events
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
			cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
			batchBuilder.AddEth(services.EthPacket{
				Type:        services.EthWriteAck,
				Packet:      &cosmosPacket,
				AckBytes:    ev.Acknowledgements,
				BlockNumber: ev.Raw.BlockNumber,
			})
			batchBuilder.PendingTracker.Remove(cosmosPacket.SourceClient, cosmosPacket.Sequence)

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
			ctx.Logger.Printf("SendPacket subscription error: %v", err)
			return

		case err := <-writeAckSub.Err():
			ctx.Logger.Printf("WriteAcknowledgement subscription error: %v", err)
			return

		case err := <-ackPacketSub.Err():
			ctx.Logger.Printf("AckPacket subscription error: %v", err)
			return

		case err := <-timeoutPacketSub.Err():
			ctx.Logger.Printf("TimeoutPacket subscription error: %v", err)
			return
		}
	}
}
