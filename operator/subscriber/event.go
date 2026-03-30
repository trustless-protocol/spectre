package subscriber

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	contractICS26Router "operator/bindings/ICS26Router"
	"operator/services"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogo/protobuf/proto"
)

const COMETBFT_SEND_PACKET_EVENT = "tm.event = 'Tx' AND message.action = '/ibc.applications.transfer.v1.MsgTransfer'"
const COMETBFT_ACK_PACKET_EVENT = "tm.event = 'Tx' AND message.action = '/ibc.applications.transfer.v1.MsgAcknowledgement'"
const COMETBFT_TIMEOUT_PACKET_EVENT = "tm.event = 'Tx' AND message.action = '/ibc.applications.transfer.v1.MsgTimeout'"

const EVENT_SEND_PACKET_FIELD = "send_packet.encoded_packet_hex"
const EVENT_ACK_PACKET_FIELD = "acknowledge_packet.encoded_packet_hex"
const EVENT_TIMEOUT_PACKET_FIELD = "timeout_packet.encoded_packet_hex"

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
		ctx.Logger.Println(err.Error())
	}
	ackPacketSub, err := ctx.CosmosClient().WSEvents.Subscribe(context.Background(), "", COMETBFT_ACK_PACKET_EVENT)
	if err != nil {
		ctx.Logger.Println(err.Error())
	}
	timeoutPacketSub, err := ctx.CosmosClient().WSEvents.Subscribe(context.Background(), "", COMETBFT_TIMEOUT_PACKET_EVENT)
	if err != nil {
		ctx.Logger.Println(err.Error())
	}
	defer ctx.CosmosClient().UnsubscribeAll(context.Background(), "")

	for {
		select {
		case e := <-sendPacketSub:
			// handle event
			sendPacketEvent := e.Events[EVENT_SEND_PACKET_FIELD]
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

			batchBuilder.InsertPacket(services.Packet{
				PacketType: services.Send,
				Packet:     &packet,
			})
		case e := <-ackPacketSub:
			// handle event
			ackPacketEvent := e.Events[EVENT_ACK_PACKET_FIELD]
			if ackPacketEvent == nil {
				continue
			}

			packetEncodedStr := ackPacketEvent[0]
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

			batchBuilder.InsertPacket(services.Packet{
				PacketType: services.Ack,
				Packet:     &packet,
			})
		case e := <-timeoutPacketSub:
			// handle event
			timeoutPacketEvent := e.Events[EVENT_TIMEOUT_PACKET_FIELD]
			if timeoutPacketEvent == nil {
				continue
			}

			packetEncodedStr := timeoutPacketEvent[0]
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

			batchBuilder.InsertPacket(services.Packet{
				PacketType: services.Timeout,
				Packet:     &packet,
			})
		case <-c.Done():
			return
		}
	}
}

// ethPacketToCosmosPacket converts an Ethereum ICS26Router packet to a Cosmos IBC v2 packet
func ethPacketToCosmosPacket(ethPacket contractICS26Router.IICS26RouterMsgsPacket, sequence *big.Int) channeltypesv2.Packet {
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
	// Get the ICS26Router contract address from environment variable
	ics26RouterAddr := os.Getenv("ICS26_ROUTER_ADDRESS")
	if ics26RouterAddr == "" {
		ctx.Logger.Println("ICS26_ROUTER_ADDRESS environment variable is required")
		return
	}

	// Parse the contract address
	contractAddr := common.HexToAddress(ics26RouterAddr)

	// Create a new ICS26Router filterer instance (only for event subscriptions)
	filterer, err := contractICS26Router.NewContractICS26RouterFilterer(contractAddr, ctx.EthClient())
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
			cosmosPacket := ethPacketToCosmosPacket(ev.Packet, ev.Sequence)
			batchBuilder.InsertPacket(services.Packet{
				PacketType: services.Send,
				Packet:     &cosmosPacket,
			})

		case ev := <-writeAckCh:
			// TODO: handle WriteAcknowledgement event
			// This event is emitted when an acknowledgement is written on Ethereum
			ctx.Logger.Printf("WriteAcknowledgement event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())

		case ev := <-ackPacketCh:
			// TODO: handle AckPacket event
			// This event is emitted when a packet acknowledgement is received on Ethereum
			ctx.Logger.Printf("AckPacket event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())

		case ev := <-timeoutPacketCh:
			// TODO: handle TimeoutPacket event
			// This event is emitted when a packet times out
			ctx.Logger.Printf("TimeoutPacket event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())

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
