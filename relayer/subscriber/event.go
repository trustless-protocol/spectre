package subscriber

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/services"
	"relayer/utils"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	commettypes "github.com/cometbft/cometbft/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogo/protobuf/proto"
)

// Match by IBC v2 event TYPE (send_packet / write_acknowledgement / timeout_packet)
// rather than by message.action. The Rust upstream relayer takes the same approach
// (packages/relayer/lib/src/events/eureka.rs::TryFrom<TmEvent>) — it iterates the tx's
// events and matches event.kind against cosmos_sdk::EVENT_TYPE_*. Using event TYPE
// here means we catch any tx that produces a send_packet (whether triggered by
// gaiad's MsgTransfer wrapper, a direct channeltypesv2.MsgSendPacket from a test
// harness, or any other module that ends up calling the v2 channel keeper), instead
// of being tied to one specific message.action.
const COMETBFT_SEND_PACKET_EVENT = "tm.event = 'Tx' AND send_packet.encoded_packet_hex EXISTS"
const COMETBFT_WRITE_ACK_PACKET_EVENT = "tm.event = 'Tx' AND write_acknowledgement.encoded_packet_hex EXISTS"
const COMETBFT_TIMEOUT_PACKET_EVENT = "tm.event = 'Tx' AND timeout_packet.encoded_packet_hex EXISTS"

const EVENT_SEND_PACKET_FIELD = "send_packet.encoded_packet_hex"
const EVENT_WRITE_ACK_PACKET_FIELD = "write_acknowledgement.encoded_packet_hex"
const EVENT_ACKNOWLEDGEMENT_FIELD = "write_acknowledgement.encoded_acknowledgement_hex"
const EVENT_TIMEOUT_PACKET_FIELD = "timeout_packet.encoded_packet_hex"
const EVENT_TX_HEIGHT_FIELD = "tx.height"

const ethStartupRecoveryLookbackEnv = "ETH_STARTUP_LOOKBACK_BLOCKS"
const defaultEthStartupRecoveryLookbackBlocks uint64 = 256
const ethSubscriptionReconnectDelay = 2 * time.Second
const cosmosStartupRecoveryLookbackEnv = "COSMOS_STARTUP_LOOKBACK_BLOCKS"
const defaultCosmosStartupRecoveryLookbackBlocks uint64 = 256
const cosmosRecoveryPollInterval = 5 * time.Second
const cosmosTxSearchPerPage = 100
const cosmosSeenRetentionBlocks uint64 = 1024

// normalizeTimeoutSeconds converts IBC v2 timeout timestamps from nanoseconds to seconds.
// ibc-go stores TimeoutTimestamp in nanoseconds, but the ETH side uses seconds.
// Values < 1e12 are assumed to already be in seconds (year ~33658 CE in seconds).
func normalizeTimeoutSeconds(ts uint64) uint64 {
	if ts >= 1e12 {
		return ts / 1e9
	}
	return ts
}

type Subscriber struct {
	cosmosSeenMtx sync.Mutex
	cosmosSeen    map[string]uint64
}

func NewSubscriber() *Subscriber {
	return &Subscriber{cosmosSeen: make(map[string]uint64)}
}

func (s *Subscriber) SubscribeCosmos(ctx services.Context, batchBuilder *services.BatchBuilder) {
	lookback := cosmosStartupRecoveryLookbackBlocks()
	var nextSendRecoveryStartBlock uint64
	var nextAckRecoveryStartBlock uint64
	var nextTimeoutRecoveryStartBlock uint64
	for {
		latestBlock, err := latestCosmosBlockHeight(ctx)
		if err != nil {
			ctx.Logger.Printf("[SubscribeCosmos] failed to get latest Cosmos block before subscription: %v", err)
			time.Sleep(ethSubscriptionReconnectDelay)
			continue
		}
		if nextSendRecoveryStartBlock == 0 {
			nextSendRecoveryStartBlock = cosmosStartupRecoveryStartBlock(latestBlock, lookback)
		}
		if nextAckRecoveryStartBlock == 0 {
			nextAckRecoveryStartBlock = cosmosStartupRecoveryStartBlock(latestBlock, lookback)
		}
		if nextTimeoutRecoveryStartBlock == 0 {
			nextTimeoutRecoveryStartBlock = cosmosStartupRecoveryStartBlock(latestBlock, lookback)
		}

		s.recoverCosmosGaps(ctx, batchBuilder, latestBlock,
			&nextSendRecoveryStartBlock,
			&nextAckRecoveryStartBlock,
			&nextTimeoutRecoveryStartBlock,
		)

		err = s.subscribeCosmosOnce(ctx, batchBuilder,
			&nextSendRecoveryStartBlock,
			&nextAckRecoveryStartBlock,
			&nextTimeoutRecoveryStartBlock,
		)
		ctx.Logger.Printf("[SubscribeCosmos] Subscription loop ended: %v", err)
		time.Sleep(ethSubscriptionReconnectDelay)
	}
}

func (s *Subscriber) subscribeCosmosOnce(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	nextSendRecoveryStartBlock *uint64,
	nextAckRecoveryStartBlock *uint64,
	nextTimeoutRecoveryStartBlock *uint64,
) error {
	subCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer ctx.CosmosClient().UnsubscribeAll(context.Background(), "")

	sendPacketSub, err := ctx.CosmosClient().WSEvents.Subscribe(subCtx, "", COMETBFT_SEND_PACKET_EVENT)
	if err != nil {
		return fmt.Errorf("failed to subscribe to send_packet events: %w", err)
	}
	ackPacketSub, err := ctx.CosmosClient().WSEvents.Subscribe(subCtx, "", COMETBFT_WRITE_ACK_PACKET_EVENT)
	if err != nil {
		return fmt.Errorf("failed to subscribe to write_acknowledgement events: %w", err)
	}
	timeoutPacketSub, err := ctx.CosmosClient().WSEvents.Subscribe(subCtx, "", COMETBFT_TIMEOUT_PACKET_EVENT)
	if err != nil {
		return fmt.Errorf("failed to subscribe to timeout_packet events: %w", err)
	}
	ctx.Logger.Println("[SubscribeCosmos] Successfully subscribed to CometBFT events")

	recoveryTicker := time.NewTicker(cosmosRecoveryPollInterval)
	defer recoveryTicker.Stop()

	for {
		select {
		case e, ok := <-sendPacketSub:
			if !ok {
				return fmt.Errorf("send_packet subscription channel closed")
			}
			blockNumber := s.handleCosmosSendPacketEvent(ctx, batchBuilder, e.Data, e.Events, "live")
			advanceRecoveryStart(nextSendRecoveryStartBlock, blockNumber+1)

		case e, ok := <-ackPacketSub:
			if !ok {
				return fmt.Errorf("write_acknowledgement subscription channel closed")
			}
			blockNumber := s.handleCosmosWriteAckEvent(ctx, batchBuilder, e.Data, e.Events, "live")
			advanceRecoveryStart(nextAckRecoveryStartBlock, blockNumber+1)

		case e, ok := <-timeoutPacketSub:
			if !ok {
				return fmt.Errorf("timeout_packet subscription channel closed")
			}
			blockNumber := s.handleCosmosTimeoutEvent(ctx, batchBuilder, e.Data, e.Events, "live")
			advanceRecoveryStart(nextTimeoutRecoveryStartBlock, blockNumber+1)

		case <-recoveryTicker.C:
			latestBlock, err := latestCosmosBlockHeight(ctx)
			if err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] gap recovery: failed to get latest Cosmos block: %v", err)
				continue
			}
			s.recoverCosmosGaps(ctx, batchBuilder, latestBlock,
				nextSendRecoveryStartBlock,
				nextAckRecoveryStartBlock,
				nextTimeoutRecoveryStartBlock,
			)
		}
	}
}

func latestCosmosBlockHeight(ctx services.Context) (uint64, error) {
	status, err := ctx.CosmosClient().Status(context.Background())
	if err != nil {
		return 0, err
	}
	if status.SyncInfo.LatestBlockHeight < 0 {
		return 0, fmt.Errorf("negative latest block height %d", status.SyncInfo.LatestBlockHeight)
	}
	return uint64(status.SyncInfo.LatestBlockHeight), nil
}

func (s *Subscriber) recoverCosmosGaps(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	latestBlock uint64,
	nextSendRecoveryStartBlock *uint64,
	nextAckRecoveryStartBlock *uint64,
	nextTimeoutRecoveryStartBlock *uint64,
) {
	*nextSendRecoveryStartBlock = s.recoverCosmosEventQuery(
		ctx, batchBuilder, "send_packet", COMETBFT_SEND_PACKET_EVENT,
		*nextSendRecoveryStartBlock, latestBlock, s.handleCosmosSendPacketEvent,
	)
	*nextAckRecoveryStartBlock = s.recoverCosmosEventQuery(
		ctx, batchBuilder, "write_acknowledgement", COMETBFT_WRITE_ACK_PACKET_EVENT,
		*nextAckRecoveryStartBlock, latestBlock, s.handleCosmosWriteAckEvent,
	)
	*nextTimeoutRecoveryStartBlock = s.recoverCosmosEventQuery(
		ctx, batchBuilder, "timeout_packet", COMETBFT_TIMEOUT_PACKET_EVENT,
		*nextTimeoutRecoveryStartBlock, latestBlock, s.handleCosmosTimeoutEvent,
	)
}

func (s *Subscriber) recoverCosmosEventQuery(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	name string,
	baseQuery string,
	startBlock uint64,
	endBlock uint64,
	handle func(services.Context, *services.BatchBuilder, commettypes.TMEventData, map[string][]string, string) uint64,
) uint64 {
	if endBlock < startBlock {
		return startBlock
	}

	ctx.Logger.Printf("[SubscribeCosmos] recovery scanning %s txs in [%d,%d]", name, startBlock, endBlock)
	query := fmt.Sprintf("%s AND tx.height >= %d AND tx.height <= %d", baseQuery, startBlock, endBlock)
	perPage := cosmosTxSearchPerPage

	for page := 1; ; page++ {
		pageNumber := page
		result, err := ctx.CosmosClient().TxSearch(context.Background(), query, false, &pageNumber, &perPage, "asc")
		if err != nil {
			ctx.Logger.Printf("[SubscribeCosmos] recovery: failed to search %s txs in [%d,%d]: %v",
				name, startBlock, endBlock, err)
			return startBlock
		}
		for _, tx := range result.Txs {
			events := cosmosEventsFromTx(tx)
			handle(ctx, batchBuilder, nil, events, "recovery")
		}
		if len(result.Txs) == 0 || pageNumber*perPage >= result.TotalCount {
			break
		}
	}

	return endBlock + 1
}

func cosmosEventsFromTx(tx *coretypes.ResultTx) map[string][]string {
	events := make(map[string][]string)
	if tx == nil {
		return events
	}
	for _, event := range tx.TxResult.Events {
		for _, attr := range event.Attributes {
			key := fmt.Sprintf("%s.%s", event.Type, attr.Key)
			events[key] = append(events[key], attr.Value)
		}
	}
	if tx.Height > 0 {
		events[EVENT_TX_HEIGHT_FIELD] = append(events[EVENT_TX_HEIGHT_FIELD], strconv.FormatInt(tx.Height, 10))
	}
	return events
}

func cosmosPacketMatchesConfiguredClient(ctx services.Context, packet *channeltypesv2.Packet) bool {
	if packet == nil {
		return false
	}
	ethClientID := ctx.EthClientID()
	routerClientID := ctx.CosmosRouterClientID()
	if ethClientID == "" && routerClientID == "" {
		return true
	}
	return packet.SourceClient == ethClientID ||
		packet.DestinationClient == ethClientID ||
		packet.SourceClient == routerClientID ||
		packet.DestinationClient == routerClientID
}

func (s *Subscriber) markCosmosPacketSeen(packetType services.CosmosPacketType, packet *channeltypesv2.Packet, blockNumber uint64) bool {
	if packet == nil {
		return false
	}
	key := fmt.Sprintf("%d|%s|%s|%d|%d", int(packetType), packet.SourceClient, packet.DestinationClient, packet.Sequence, blockNumber)

	s.cosmosSeenMtx.Lock()
	defer s.cosmosSeenMtx.Unlock()
	if s.cosmosSeen == nil {
		s.cosmosSeen = make(map[string]uint64)
	}
	if _, ok := s.cosmosSeen[key]; ok {
		return false
	}
	s.cosmosSeen[key] = blockNumber
	if blockNumber > cosmosSeenRetentionBlocks {
		minHeight := blockNumber - cosmosSeenRetentionBlocks
		for seenKey, seenHeight := range s.cosmosSeen {
			if seenHeight != 0 && seenHeight < minHeight {
				delete(s.cosmosSeen, seenKey)
			}
		}
	}
	return true
}

func decodeCosmosPacket(ctx services.Context, label string, encoded string) (*channeltypesv2.Packet, bool) {
	packetBytes, err := hex.DecodeString(encoded)
	if err != nil {
		ctx.Logger.Printf("[SubscribeCosmos] %s: failed to decode packet hex: %v", label, err)
		return nil, false
	}

	var packet channeltypesv2.Packet
	if err := proto.Unmarshal(packetBytes, &packet); err != nil {
		ctx.Logger.Printf("[SubscribeCosmos] %s: failed to unmarshal packet: %v", label, err)
		return nil, false
	}
	return &packet, true
}

func (s *Subscriber) handleCosmosSendPacketEvent(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	data commettypes.TMEventData,
	events map[string][]string,
	source string,
) uint64 {
	sendPacketEvent := events[EVENT_SEND_PACKET_FIELD]
	blockNumber := txHeightFromEvent(data, events)
	if len(sendPacketEvent) == 0 {
		return blockNumber
	}

	for _, packetEncodedStr := range sendPacketEvent {
		packet, ok := decodeCosmosPacket(ctx, "send_packet", packetEncodedStr)
		if !ok {
			continue
		}
		if !cosmosPacketMatchesConfiguredClient(ctx, packet) {
			ctx.Logger.Printf("[SubscribeCosmos] send_packet seq=%d src=%s dst=%s ignored: unrelated client",
				packet.Sequence, packet.SourceClient, packet.DestinationClient)
			continue
		}
		if !s.markCosmosPacketSeen(services.CosmosSend, packet, blockNumber) {
			ctx.Logger.Printf("[SubscribeCosmos] send_packet seq=%d at height %d ignored: duplicate %s event",
				packet.Sequence, blockNumber, source)
			continue
		}

		ctx.Logger.Printf("[SubscribeCosmos] send_packet received: seq=%d src=%s dst=%s height=%d source=%s",
			packet.Sequence, packet.SourceClient, packet.DestinationClient, blockNumber, source)
		packet.TimeoutTimestamp = normalizeTimeoutSeconds(packet.TimeoutTimestamp)
		batchBuilder.AddCosmos(services.CosmosPacket{
			Type:        services.CosmosSend,
			Packet:      packet,
			BlockNumber: blockNumber,
		})
	}
	return blockNumber
}

func (s *Subscriber) handleCosmosWriteAckEvent(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	data commettypes.TMEventData,
	events map[string][]string,
	source string,
) uint64 {
	ackPacketEvent := events[EVENT_WRITE_ACK_PACKET_FIELD]
	ackEvent := events[EVENT_ACKNOWLEDGEMENT_FIELD]
	blockNumber := txHeightFromEvent(data, events)
	if len(ackPacketEvent) == 0 || len(ackEvent) == 0 {
		return blockNumber
	}
	if len(ackPacketEvent) != len(ackEvent) {
		ctx.Logger.Printf("[SubscribeCosmos] write_ack: packet/ack count mismatch (%d vs %d), skipping",
			len(ackPacketEvent), len(ackEvent))
		return blockNumber
	}

	for i := range ackPacketEvent {
		packet, ok := decodeCosmosPacket(ctx, "write_ack", ackPacketEvent[i])
		if !ok {
			continue
		}
		if !cosmosPacketMatchesConfiguredClient(ctx, packet) {
			ctx.Logger.Printf("[SubscribeCosmos] write_ack seq=%d src=%s dst=%s ignored: unrelated client",
				packet.Sequence, packet.SourceClient, packet.DestinationClient)
			continue
		}

		ackBytes, err := hex.DecodeString(ackEvent[i])
		if err != nil {
			ctx.Logger.Printf("[SubscribeCosmos] write_ack seq=%d: failed to decode ack hex: %v",
				packet.Sequence, err)
			continue
		}

		var acknowledgement channeltypesv2.Acknowledgement
		if err := proto.Unmarshal(ackBytes, &acknowledgement); err != nil {
			ctx.Logger.Printf("[SubscribeCosmos] write_ack seq=%d: failed to unmarshal ack: %v",
				packet.Sequence, err)
			continue
		}
		if len(acknowledgement.AppAcknowledgements) == 0 {
			ctx.Logger.Printf("[SubscribeCosmos] write_ack seq=%d: missing app acknowledgements", packet.Sequence)
			continue
		}
		if !s.markCosmosPacketSeen(services.CosmosAck, packet, blockNumber) {
			ctx.Logger.Printf("[SubscribeCosmos] write_ack seq=%d at height %d ignored: duplicate %s event",
				packet.Sequence, blockNumber, source)
			continue
		}

		ctx.Logger.Printf("[SubscribeCosmos] write_ack received: seq=%d src=%s dst=%s height=%d source=%s",
			packet.Sequence, packet.SourceClient, packet.DestinationClient, blockNumber, source)
		packet.TimeoutTimestamp = normalizeTimeoutSeconds(packet.TimeoutTimestamp)
		batchBuilder.AddCosmos(services.CosmosPacket{
			Type:        services.CosmosAck,
			Packet:      packet,
			AckBytes:    acknowledgement.AppAcknowledgements,
			BlockNumber: blockNumber,
		})
	}
	return blockNumber
}

func (s *Subscriber) handleCosmosTimeoutEvent(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	data commettypes.TMEventData,
	events map[string][]string,
	source string,
) uint64 {
	timeoutPacketEvent := events[EVENT_TIMEOUT_PACKET_FIELD]
	blockNumber := txHeightFromEvent(data, events)
	if len(timeoutPacketEvent) == 0 {
		return blockNumber
	}

	for _, packetEncodedStr := range timeoutPacketEvent {
		packet, ok := decodeCosmosPacket(ctx, "timeout", packetEncodedStr)
		if !ok {
			continue
		}
		if !cosmosPacketMatchesConfiguredClient(ctx, packet) {
			ctx.Logger.Printf("[SubscribeCosmos] timeout seq=%d src=%s dst=%s ignored: unrelated client",
				packet.Sequence, packet.SourceClient, packet.DestinationClient)
			continue
		}
		if !s.markCosmosPacketSeen(services.CosmosTimeout, packet, blockNumber) {
			ctx.Logger.Printf("[SubscribeCosmos] timeout seq=%d at height %d ignored: duplicate %s event",
				packet.Sequence, blockNumber, source)
			continue
		}

		ctx.Logger.Printf("[SubscribeCosmos] timeout received: seq=%d src=%s dst=%s height=%d source=%s",
			packet.Sequence, packet.SourceClient, packet.DestinationClient, blockNumber, source)
		packet.TimeoutTimestamp = normalizeTimeoutSeconds(packet.TimeoutTimestamp)
		batchBuilder.AddCosmos(services.CosmosPacket{
			Type:        services.CosmosTimeout,
			Packet:      packet,
			BlockNumber: blockNumber,
		})
	}
	return blockNumber
}

func txHeightFromEvent(data commettypes.TMEventData, events map[string][]string) uint64 {
	switch txData := data.(type) {
	case commettypes.EventDataTx:
		if txData.Height > 0 {
			return uint64(txData.Height)
		}
	case *commettypes.EventDataTx:
		if txData != nil && txData.Height > 0 {
			return uint64(txData.Height)
		}
	}

	values := events[EVENT_TX_HEIGHT_FIELD]
	if len(values) == 0 {
		return 0
	}
	height, err := strconv.ParseUint(values[0], 10, 64)
	if err != nil {
		return 0
	}
	return height
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

func cosmosStartupRecoveryLookbackBlocksFromEnv(raw string) uint64 {
	if raw == "" {
		return defaultCosmosStartupRecoveryLookbackBlocks
	}

	lookback, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return defaultCosmosStartupRecoveryLookbackBlocks
	}

	return lookback
}

func cosmosStartupRecoveryLookbackBlocks() uint64 {
	return cosmosStartupRecoveryLookbackBlocksFromEnv(os.Getenv(cosmosStartupRecoveryLookbackEnv))
}

func ethStartupRecoveryStartBlock(latestBlock, lookback uint64) uint64 {
	if lookback >= latestBlock {
		return 0
	}

	return latestBlock - lookback
}

func cosmosStartupRecoveryStartBlock(latestBlock, lookback uint64) uint64 {
	return ethStartupRecoveryStartBlock(latestBlock, lookback)
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

func enqueueEthSendPacket(batchBuilder *services.BatchBuilder, ev *contractICS26Router.ContractICS26RouterSendPacket) {
	cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
	batchBuilder.AddEth(services.EthPacket{
		Type:        services.EthSend,
		Packet:      &cosmosPacket,
		BlockNumber: ev.Raw.BlockNumber,
	})
}

func hasCosmosIBCPathValue(ctx services.Context, path [][]byte) (bool, error) {
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

func hasPendingCosmosPacketCommitment(ctx services.Context, packet channeltypesv2.Packet) (bool, error) {
	return hasCosmosIBCPathValue(ctx, utils.IbcCommitmentPath(packet, []byte{1}))
}

func hasCosmosPacketReceipt(ctx services.Context, packet channeltypesv2.Packet) (bool, error) {
	return hasCosmosIBCPathValue(ctx, utils.IbcPath(packet.DestinationClient, packet.Sequence, []byte{2}))
}

func recoverEthSendPackets(
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
	iter, err := filterer.FilterSendPacket(filterOpts, nil, nil)
	if err != nil {
		ctx.Logger.Printf("[SubscribeEth] startup recovery: failed to filter SendPacket logs in [%d,%d]: %v",
			startBlock, endBlock, err)
		return
	}
	defer iter.Close()

	var recoveredCount uint64
	var skippedCount uint64
	for iter.Next() {
		ev := iter.Event
		cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)

		received, err := hasCosmosPacketReceipt(ctx, cosmosPacket)
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] startup recovery: seq=%d failed to check Cosmos packet receipt: %v",
				cosmosPacket.Sequence, err)
			continue
		}
		if received {
			skippedCount++
			ctx.Logger.Printf("[SubscribeEth] startup recovery: seq=%d already received on Cosmos, skipping historical SendPacket from ETH block %d",
				cosmosPacket.Sequence, ev.Raw.BlockNumber)
			continue
		}

		ctx.Logger.Printf("[SubscribeEth] startup recovery: recovered SendPacket seq=%d from ETH block %d",
			cosmosPacket.Sequence, ev.Raw.BlockNumber)
		enqueueEthSendPacket(batchBuilder, ev)
		recoveredCount++
	}

	if err := iter.Error(); err != nil {
		ctx.Logger.Printf("[SubscribeEth] startup recovery: SendPacket iterator error in [%d,%d]: %v",
			startBlock, endBlock, err)
	}

	ctx.Logger.Printf("[SubscribeEth] startup recovery complete for SendPacket: scanned [%d,%d], recovered=%d skipped=%d",
		startBlock, endBlock, recoveredCount, skippedCount)
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
	// Live handlers pass block+1 here, so periodic recovery trusts the live
	// subscription to have covered every relevant event in the observed block.
	if candidate > *nextRecoveryStartBlock {
		*nextRecoveryStartBlock = candidate
	}
}

func (s *Subscriber) subscribeEthOnce(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	watchClient *ethclient.Client,
	watchStartBlock uint64,
	nextSendRecoveryStartBlock *uint64,
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
			advanceRecoveryStart(nextSendRecoveryStartBlock, ev.Raw.BlockNumber+1)
			enqueueEthSendPacket(batchBuilder, ev)

		case ev := <-writeAckCh:
			ctx.Logger.Printf("WriteAcknowledgement event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			advanceRecoveryStart(nextWriteAckRecoveryStartBlock, ev.Raw.BlockNumber+1)
			enqueueEthWriteAcknowledgement(batchBuilder, ev)
			batchBuilder.PendingTracker.Remove(ev.Packet.SourceClient, ev.Sequence.Uint64())

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
	var nextSendRecoveryStartBlock uint64
	var nextWriteAckRecoveryStartBlock uint64

	for {
		latestBlock, err := ctx.EthClient().BlockNumber(context.Background())
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] Failed to get latest Ethereum block before subscription: %v", err)
			time.Sleep(ethSubscriptionReconnectDelay)
			continue
		}

		if nextSendRecoveryStartBlock == 0 {
			nextSendRecoveryStartBlock = ethStartupRecoveryStartBlock(latestBlock, lookback)
		}
		if nextWriteAckRecoveryStartBlock == 0 {
			nextWriteAckRecoveryStartBlock = ethStartupRecoveryStartBlock(latestBlock, lookback)
		}

		if latestBlock >= nextSendRecoveryStartBlock {
			ctx.Logger.Printf("[SubscribeEth] recovery scanning SendPacket logs in [%d,%d]",
				nextSendRecoveryStartBlock, latestBlock)
			recoverEthSendPackets(ctx, batchBuilder, recoveryFilterer, nextSendRecoveryStartBlock, latestBlock)
		}

		if latestBlock >= nextWriteAckRecoveryStartBlock {
			ctx.Logger.Printf("[SubscribeEth] recovery scanning WriteAcknowledgement logs in [%d,%d]",
				nextWriteAckRecoveryStartBlock, latestBlock)
			recoverEthWriteAcknowledgements(ctx, batchBuilder, recoveryFilterer, nextWriteAckRecoveryStartBlock, latestBlock)
		}

		watchStartBlock := latestBlock + 1
		nextSendRecoveryStartBlock = watchStartBlock
		nextWriteAckRecoveryStartBlock = watchStartBlock

		watchClient, err := ethclient.DialContext(context.Background(), ctx.EthWsURL())
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] Failed to connect to Ethereum WS at %s: %v", ctx.EthWsURL(), err)
			time.Sleep(ethSubscriptionReconnectDelay)
			continue
		}

		err = s.subscribeEthOnce(ctx, batchBuilder, watchClient, watchStartBlock, &nextSendRecoveryStartBlock, &nextWriteAckRecoveryStartBlock)
		watchClient.Close()
		ctx.Logger.Printf("[SubscribeEth] Subscription loop ended: %v", err)
		time.Sleep(ethSubscriptionReconnectDelay)
	}
}
