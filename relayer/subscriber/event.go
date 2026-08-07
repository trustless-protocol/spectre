package subscriber

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/services"
	"relayer/utils"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	commettypes "github.com/cometbft/cometbft/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
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
const ethGapRecoveryInterval = 30 * time.Second
const ethSeenEventRetentionBlocks uint64 = 2_000

const cosmosStartupRecoveryLookbackEnv = "COSMOS_STARTUP_LOOKBACK_BLOCKS"
const defaultCosmosStartupRecoveryLookbackBlocks uint64 = 256
const cosmosSubscriptionReconnectDelay = 2 * time.Second
const cosmosGapRecoveryInterval = 30 * time.Second

// quietScanHeartbeat is how many consecutive find-nothing recovery scans pass
// before one is reported, so a healthy subscriber stays visible (~10 min at the
// 30s interval) without printing every tick.
const quietScanHeartbeat = 20

// advanceQuietScans counts one find-nothing pass and reports whether this is the
// one to log. Shared by both directions so the two heartbeats cannot drift.
func advanceQuietScans(quietScans *uint64) bool {
	*quietScans++
	return *quietScans%quietScanHeartbeat == 0
}

const cosmosTxSearchPerPage = 100
const cosmosSeenEventRetentionHeights uint64 = 2_000

const cometBFTSendPacketTxSearch = "send_packet.encoded_packet_hex EXISTS"
const cometBFTWriteAckPacketTxSearch = "write_acknowledgement.encoded_packet_hex EXISTS"
const cometBFTTimeoutPacketTxSearch = "timeout_packet.encoded_packet_hex EXISTS"

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
}

func NewSubscriber() *Subscriber {
	return &Subscriber{}
}

func (s *Subscriber) SubscribeCosmos(ctx services.Context, batchBuilder *services.BatchBuilder) {
	lookback := cosmosStartupRecoveryLookbackBlocks()
	var nextRecoveryStartHeight uint64
	// quietScans counts consecutive recovery passes that found nothing. It lives
	// beside the cursor, outside the resubscribe loop, so a WS reconnect does not
	// restart the heartbeat — mirroring the Ethereum side.
	var quietScans uint64
	seenEvents := make(map[cosmosEventKey]struct{})

	for {
		if nextRecoveryStartHeight == 0 {
			latestHeight, err := latestCosmosHeight(ctx)
			if err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] Failed to get latest Cosmos height before subscription: %v", err)
				time.Sleep(cosmosSubscriptionReconnectDelay)
				continue
			}
			nextRecoveryStartHeight = cosmosStartupRecoveryStartHeight(latestHeight, lookback)
		}

		err := s.subscribeCosmosOnce(ctx, batchBuilder, &nextRecoveryStartHeight, seenEvents, &quietScans)
		ctx.Logger.Printf("[SubscribeCosmos] Subscription loop ended: %v", err)
		time.Sleep(cosmosSubscriptionReconnectDelay)
	}
}

func (s *Subscriber) subscribeCosmosOnce(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	nextRecoveryStartHeight *uint64,
	seenEvents map[cosmosEventKey]struct{},
	quietScans *uint64,
) error {
	sendPacketSub, err := ctx.CosmosClient().WSEvents.Subscribe(context.Background(), "", COMETBFT_SEND_PACKET_EVENT)
	if err != nil {
		return fmt.Errorf("failed to subscribe to send_packet events: %w", err)
	}
	ackPacketSub, err := ctx.CosmosClient().WSEvents.Subscribe(context.Background(), "", COMETBFT_WRITE_ACK_PACKET_EVENT)
	if err != nil {
		_ = ctx.CosmosClient().UnsubscribeAll(context.Background(), "")
		return fmt.Errorf("failed to subscribe to write_acknowledgement events: %w", err)
	}
	timeoutPacketSub, err := ctx.CosmosClient().WSEvents.Subscribe(context.Background(), "", COMETBFT_TIMEOUT_PACKET_EVENT)
	if err != nil {
		_ = ctx.CosmosClient().UnsubscribeAll(context.Background(), "")
		return fmt.Errorf("failed to subscribe to timeout_packet events: %w", err)
	}
	ctx.Logger.Println("[SubscribeCosmos] Successfully subscribed to CometBFT events")
	defer ctx.CosmosClient().UnsubscribeAll(context.Background(), "")

	if err := recoverCosmosGapToLatest(ctx, batchBuilder, nextRecoveryStartHeight, seenEvents, quietScans); err != nil {
		ctx.Logger.Printf("[SubscribeCosmos] startup recovery failed: %v", err)
	}

	gapRecoveryTicker := time.NewTicker(cosmosGapRecoveryInterval)
	defer gapRecoveryTicker.Stop()

	for {
		select {
		case e, ok := <-sendPacketSub:
			if !ok {
				return fmt.Errorf("send_packet subscription channel closed")
			}
			s.processLiveCosmosEvent(ctx, batchBuilder, nextRecoveryStartHeight, seenEvents, e)
		case e, ok := <-ackPacketSub:
			if !ok {
				return fmt.Errorf("write_acknowledgement subscription channel closed")
			}
			s.processLiveCosmosEvent(ctx, batchBuilder, nextRecoveryStartHeight, seenEvents, e)
		case e, ok := <-timeoutPacketSub:
			if !ok {
				return fmt.Errorf("timeout_packet subscription channel closed")
			}
			s.processLiveCosmosEvent(ctx, batchBuilder, nextRecoveryStartHeight, seenEvents, e)
		case <-gapRecoveryTicker.C:
			if err := recoverCosmosGapToLatest(ctx, batchBuilder, nextRecoveryStartHeight, seenEvents, quietScans); err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] periodic recovery failed: %v", err)
			}
		}
	}
}

func (s *Subscriber) processLiveCosmosEvent(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	nextRecoveryStartHeight *uint64,
	seenEvents map[cosmosEventKey]struct{},
	e coretypes.ResultEvent,
) {
	height := txHeightFromEvent(e.Data, e.Events)
	if height > 0 && *nextRecoveryStartHeight > 0 && height > *nextRecoveryStartHeight {
		stats, err := recoverCosmosEvents(ctx, batchBuilder, *nextRecoveryStartHeight, height-1, seenEvents)
		if err != nil {
			ctx.Logger.Printf("[SubscribeCosmos] gap recovery before live height %d failed: %v", height, err)
		} else {
			ctx.Logger.Printf("[SubscribeCosmos] gap recovery complete before live height %d: recovered=%d skipped=%d",
				height, stats.recovered, stats.skipped)
			*nextRecoveryStartHeight = height
		}
	}

	packets := decodeCosmosPacketsFromEvents(ctx.Logger, e.Data, e.Events, "SubscribeCosmos")
	stats, err := enqueueCosmosPackets(ctx, batchBuilder, packets, seenEvents, false)
	if err != nil {
		ctx.Logger.Printf("[SubscribeCosmos] live enqueue failed: %v", err)
	}
	if stats.skipped > 0 {
		ctx.Logger.Printf("[SubscribeCosmos] skipped %d live Cosmos event(s)", stats.skipped)
	}
	if height > 0 {
		pruneCosmosSeenEvents(seenEvents, height)
	}
}

type cosmosEventKey struct {
	PacketType        services.CosmosPacketType
	SourceClient      string
	DestinationClient string
	Sequence          uint64
	Height            uint64
}

type cosmosRecoveryStats struct {
	recovered uint64
	skipped   uint64
}

type ethEventKey struct {
	EventType   string
	TxHash      string
	LogIndex    uint
	BlockNumber uint64
}

type ethRecoveryStats struct {
	recovered uint64
	skipped   uint64
}

// foundSomething reports whether a recovery pass has anything to say. Every
// branch that runs a scan must fold its result in through this: a pass that
// recovered events while another found none is NOT a quiet scan, and reporting
// it as one prints "found nothing" beside the line that found something.
func (s ethRecoveryStats) foundSomething() bool {
	return s.recovered > 0 || s.skipped > 0
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

func cosmosStartupRecoveryStartHeight(latestHeight, lookback uint64) uint64 {
	if latestHeight == 0 {
		return 0
	}
	if lookback == 0 || lookback >= latestHeight {
		return 1
	}
	return latestHeight - lookback + 1
}

func latestCosmosHeight(ctx services.Context) (uint64, error) {
	status, err := ctx.CosmosClient().Status(context.Background())
	if err != nil {
		return 0, err
	}
	if status.SyncInfo.LatestBlockHeight <= 0 {
		return 0, nil
	}
	return uint64(status.SyncInfo.LatestBlockHeight), nil
}

func recoverCosmosGapToLatest(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	nextRecoveryStartHeight *uint64,
	seenEvents map[cosmosEventKey]struct{},
	quietScans *uint64,
) error {
	if *nextRecoveryStartHeight == 0 {
		return nil
	}

	latestHeight, err := latestCosmosHeight(ctx)
	if err != nil {
		return err
	}
	if latestHeight < *nextRecoveryStartHeight {
		return nil
	}

	stats, err := recoverCosmosEvents(ctx, batchBuilder, *nextRecoveryStartHeight, latestHeight, seenEvents)
	if err != nil {
		return err
	}
	// Report a scan that FOUND something, and otherwise only a periodic heartbeat.
	// This ticks every 30s for the life of the process, and the overwhelmingly
	// common outcome is recovered=0 skipped=0 — two lines a tick, ~5.8k lines a day
	// per direction, burying the events worth reading. The heartbeat keeps "gap
	// recovery is alive and current" observable without the repetition.
	if stats.recovered > 0 || stats.skipped > 0 {
		ctx.Logger.Printf("[SubscribeCosmos] recovery scanned [%d,%d]: recovered=%d skipped=%d",
			*nextRecoveryStartHeight, latestHeight, stats.recovered, stats.skipped)
		*quietScans = 0
	} else if beat := advanceQuietScans(quietScans); beat {
		ctx.Logger.Printf("[SubscribeCosmos] gap recovery healthy: %d consecutive scans found nothing, now current at height %d",
			*quietScans, latestHeight)
	}

	*nextRecoveryStartHeight = latestHeight + 1
	pruneCosmosSeenEvents(seenEvents, latestHeight)
	return nil
}

func recoverCosmosEvents(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	startHeight uint64,
	endHeight uint64,
	seenEvents map[cosmosEventKey]struct{},
) (cosmosRecoveryStats, error) {
	var combined cosmosRecoveryStats
	if endHeight < startHeight || endHeight == 0 {
		return combined, nil
	}
	if startHeight == 0 {
		startHeight = 1
	}

	queries := []string{
		cometBFTSendPacketTxSearch,
		cometBFTWriteAckPacketTxSearch,
		cometBFTTimeoutPacketTxSearch,
	}

	var firstErr error
	for _, query := range queries {
		stats, err := recoverCosmosEventsForQuery(ctx, batchBuilder, query, startHeight, endHeight, seenEvents)
		combined.recovered += stats.recovered
		combined.skipped += stats.skipped
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return combined, firstErr
}

func recoverCosmosEventsForQuery(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	baseQuery string,
	startHeight uint64,
	endHeight uint64,
	seenEvents map[cosmosEventKey]struct{},
) (cosmosRecoveryStats, error) {
	var stats cosmosRecoveryStats
	page := 1
	perPage := cosmosTxSearchPerPage
	query := cosmosTxSearchQuery(baseQuery, startHeight, endHeight)

	for {
		result, err := ctx.CosmosClient().TxSearch(context.Background(), query, false, &page, &perPage, "asc")
		if err != nil {
			return stats, fmt.Errorf("TxSearch query %q page %d failed: %w", query, page, err)
		}
		if result == nil || len(result.Txs) == 0 {
			return stats, nil
		}

		for _, tx := range result.Txs {
			if tx.TxResult.Code != 0 {
				continue
			}
			data, events := cosmosEventsFromTxResult(tx)
			packets := decodeCosmosPacketsFromEvents(ctx.Logger, data, events, "SubscribeCosmos][recovery")
			enqueued, err := enqueueCosmosPackets(ctx, batchBuilder, packets, seenEvents, true)
			stats.recovered += enqueued.recovered
			stats.skipped += enqueued.skipped
			if err != nil {
				return stats, err
			}
		}

		if page*perPage >= result.TotalCount {
			return stats, nil
		}
		page++
	}
}

func cosmosTxSearchQuery(baseQuery string, startHeight uint64, endHeight uint64) string {
	if startHeight == 0 {
		startHeight = 1
	}
	return fmt.Sprintf("%s AND tx.height >= %d AND tx.height <= %d", baseQuery, startHeight, endHeight)
}

func cosmosEventsFromTxResult(tx *coretypes.ResultTx) (commettypes.TMEventData, map[string][]string) {
	events := map[string][]string{}
	if tx == nil {
		return commettypes.EventDataTx{}, events
	}

	events[EVENT_TX_HEIGHT_FIELD] = []string{strconv.FormatInt(tx.Height, 10)}
	for _, ev := range tx.TxResult.Events {
		for _, attr := range ev.Attributes {
			key := ev.Type + "." + attr.Key
			events[key] = append(events[key], attr.Value)
		}
	}

	return commettypes.EventDataTx{
		TxResult: abcitypes.TxResult{
			Height: tx.Height,
			Index:  tx.Index,
			Tx:     []byte(tx.Tx),
			Result: tx.TxResult,
		},
	}, events
}

func decodeCosmosPacketsFromEvents(
	logger *log.Logger,
	data commettypes.TMEventData,
	events map[string][]string,
	logPrefix string,
) []services.CosmosPacket {
	if logger == nil {
		logger = log.Default()
	}

	blockNumber := txHeightFromEvent(data, events)
	packets := make([]services.CosmosPacket, 0)

	sendPacketEvent := events[EVENT_SEND_PACKET_FIELD]
	for _, packetEncodedStr := range sendPacketEvent {
		packetBytes, err := hex.DecodeString(packetEncodedStr)
		if err != nil {
			logger.Printf("[%s] send_packet: failed to decode hex: %v", logPrefix, err)
			continue
		}

		var packet channeltypesv2.Packet
		err = proto.Unmarshal(packetBytes, &packet)
		if err != nil {
			logger.Printf("[%s] send_packet: failed to unmarshal: %v", logPrefix, err)
			continue
		}

		logger.Printf("[%s] send_packet received: seq=%d src=%s",
			logPrefix, packet.Sequence, packet.SourceClient)
		packet.TimeoutTimestamp = normalizeTimeoutSeconds(packet.TimeoutTimestamp)
		packets = append(packets, services.CosmosPacket{
			Type:        services.CosmosSend,
			Packet:      &packet,
			BlockNumber: blockNumber,
		})
	}

	ackPacketEvent := events[EVENT_WRITE_ACK_PACKET_FIELD]
	ackEvent := events[EVENT_ACKNOWLEDGEMENT_FIELD]
	if len(ackPacketEvent) != 0 || len(ackEvent) != 0 {
		if len(ackPacketEvent) != len(ackEvent) {
			logger.Printf("[%s] write_ack: packet/ack count mismatch (%d vs %d), skipping",
				logPrefix, len(ackPacketEvent), len(ackEvent))
		} else {
			for i := range ackPacketEvent {
				packetBytes, err := hex.DecodeString(ackPacketEvent[i])
				if err != nil {
					logger.Printf("[%s] write_ack: failed to decode packet hex: %v", logPrefix, err)
					continue
				}

				var packet channeltypesv2.Packet
				err = proto.Unmarshal(packetBytes, &packet)
				if err != nil {
					logger.Printf("[%s] write_ack: failed to unmarshal packet: %v", logPrefix, err)
					continue
				}

				ackBytes, err := hex.DecodeString(ackEvent[i])
				if err != nil {
					logger.Printf("[%s] write_ack seq=%d: failed to decode ack hex: %v",
						logPrefix, packet.Sequence, err)
					continue
				}

				var acknowledgement channeltypesv2.Acknowledgement
				err = proto.Unmarshal(ackBytes, &acknowledgement)
				if err != nil {
					logger.Printf("[%s] write_ack seq=%d: failed to unmarshal ack: %v",
						logPrefix, packet.Sequence, err)
					continue
				}
				if len(acknowledgement.AppAcknowledgements) == 0 {
					logger.Printf("[%s] write_ack seq=%d: missing app acknowledgements", logPrefix, packet.Sequence)
					continue
				}

				logger.Printf("[%s] write_ack received: seq=%d src=%s",
					logPrefix, packet.Sequence, packet.SourceClient)
				packet.TimeoutTimestamp = normalizeTimeoutSeconds(packet.TimeoutTimestamp)
				packets = append(packets, services.CosmosPacket{
					Type:        services.CosmosAck,
					Packet:      &packet,
					AckBytes:    acknowledgement.AppAcknowledgements,
					BlockNumber: blockNumber,
				})
			}
		}
	}

	timeoutPacketEvent := events[EVENT_TIMEOUT_PACKET_FIELD]
	for _, packetEncodedStr := range timeoutPacketEvent {
		packetBytes, err := hex.DecodeString(packetEncodedStr)
		if err != nil {
			logger.Printf("[%s] timeout: failed to decode hex: %v", logPrefix, err)
			continue
		}

		var packet channeltypesv2.Packet
		err = proto.Unmarshal(packetBytes, &packet)
		if err != nil {
			logger.Printf("[%s] timeout: failed to unmarshal: %v", logPrefix, err)
			continue
		}

		logger.Printf("[%s] timeout received: seq=%d src=%s",
			logPrefix, packet.Sequence, packet.SourceClient)
		packet.TimeoutTimestamp = normalizeTimeoutSeconds(packet.TimeoutTimestamp)
		packets = append(packets, services.CosmosPacket{
			Type:        services.CosmosTimeout,
			Packet:      &packet,
			BlockNumber: blockNumber,
		})
	}

	return packets
}

func enqueueCosmosPackets(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	packets []services.CosmosPacket,
	seenEvents map[cosmosEventKey]struct{},
	recovery bool,
) (cosmosRecoveryStats, error) {
	var stats cosmosRecoveryStats
	for _, packet := range packets {
		if packet.Packet == nil {
			stats.skipped++
			continue
		}
		if !cosmosPacketMatchesConfiguredClient(ctx, packet.Packet) {
			ctx.Logger.Printf("[SubscribeCosmos] %s seq=%d src=%s dst=%s ignored: unrelated client",
				packet.Type, packet.Packet.Sequence, packet.Packet.SourceClient, packet.Packet.DestinationClient)
			stats.skipped++
			continue
		}

		key := cosmosEventKeyForPacket(packet)
		if _, ok := seenEvents[key]; ok {
			stats.skipped++
			continue
		}

		if recovery {
			shouldEnqueue, err := shouldEnqueueRecoveredCosmosPacket(ctx, packet)
			if err != nil {
				return stats, err
			}
			if !shouldEnqueue {
				seenEvents[key] = struct{}{}
				stats.skipped++
				continue
			}
		}

		seenEvents[key] = struct{}{}
		batchBuilder.AddCosmos(packet)
		stats.recovered++
	}
	return stats, nil
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

func shouldEnqueueRecoveredCosmosPacket(ctx services.Context, packet services.CosmosPacket) (bool, error) {
	if packet.Packet == nil {
		return false, nil
	}

	switch packet.Type {
	case services.CosmosSend:
		received, err := services.HasEthPacketReceipt(ctx, *packet.Packet)
		if err != nil {
			return false, fmt.Errorf("seq=%d failed to check ETH packet receipt: %w", packet.Packet.Sequence, err)
		}
		if received {
			ctx.Logger.Printf("[SubscribeCosmos] recovery: seq=%d already received on ETH, skipping historical SendPacket from Cosmos height %d",
				packet.Packet.Sequence, packet.BlockNumber)
			return false, nil
		}
		return true, nil
	case services.CosmosAck:
		pending, err := services.HasPendingEthPacketCommitment(ctx, *packet.Packet)
		if err != nil {
			return false, fmt.Errorf("seq=%d failed to check ETH packet commitment: %w", packet.Packet.Sequence, err)
		}
		if !pending {
			ctx.Logger.Printf("[SubscribeCosmos] recovery: seq=%d already cleared on ETH, skipping historical WriteAcknowledgement from Cosmos height %d",
				packet.Packet.Sequence, packet.BlockNumber)
			return false, nil
		}
		return true, nil
	case services.CosmosTimeout:
		if !services.ShouldRelayCosmosTimeoutToEth(packet.Packet, ctx.CosmosRouterClientID()) {
			ctx.Logger.Printf("[SubscribeCosmos] recovery: seq=%d Cosmos-originated timeout already handled locally, skipping historical TimeoutPacket from Cosmos height %d",
				packet.Packet.Sequence, packet.BlockNumber)
			return false, nil
		}
		pending, err := services.HasPendingEthPacketCommitment(ctx, *packet.Packet)
		if err != nil {
			return false, fmt.Errorf("seq=%d failed to check ETH packet commitment before timeout recovery: %w", packet.Packet.Sequence, err)
		}
		if !pending {
			ctx.Logger.Printf("[SubscribeCosmos] recovery: seq=%d commitment already cleared on ETH, skipping historical TimeoutPacket from Cosmos height %d",
				packet.Packet.Sequence, packet.BlockNumber)
			return false, nil
		}
		return true, nil
	default:
		return false, nil
	}
}

func cosmosEventKeyForPacket(packet services.CosmosPacket) cosmosEventKey {
	key := cosmosEventKey{
		PacketType: packet.Type,
		Height:     packet.BlockNumber,
	}
	if packet.Packet != nil {
		key.SourceClient = packet.Packet.SourceClient
		key.DestinationClient = packet.Packet.DestinationClient
		key.Sequence = packet.Packet.Sequence
	}
	return key
}

func pruneCosmosSeenEvents(seenEvents map[cosmosEventKey]struct{}, currentHeight uint64) {
	if currentHeight <= cosmosSeenEventRetentionHeights {
		return
	}
	minHeight := currentHeight - cosmosSeenEventRetentionHeights
	for key := range seenEvents {
		if key.Height > 0 && key.Height < minHeight {
			delete(seenEvents, key)
		}
	}
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

func ethStartupRecoveryStartBlock(latestBlock, lookback uint64) uint64 {
	if lookback >= latestBlock {
		return 0
	}

	return latestBlock - lookback
}

func ethEventClientIDFilter(ctx services.Context) []string {
	clientID := ctx.CosmosRouterClientID()
	if clientID == "" {
		return nil
	}
	return []string{clientID}
}

func ethEventKeyForLog(eventType string, raw gethtypes.Log) ethEventKey {
	return ethEventKey{
		EventType:   eventType,
		TxHash:      raw.TxHash.Hex(),
		LogIndex:    raw.Index,
		BlockNumber: raw.BlockNumber,
	}
}

func markEthEventSeen(seenEvents map[ethEventKey]struct{}, key ethEventKey) bool {
	if seenEvents == nil {
		return true
	}
	if _, ok := seenEvents[key]; ok {
		return false
	}
	seenEvents[key] = struct{}{}
	return true
}

func pruneEthSeenEvents(seenEvents map[ethEventKey]struct{}, currentBlock uint64) {
	if currentBlock <= ethSeenEventRetentionBlocks {
		return
	}
	minBlock := currentBlock - ethSeenEventRetentionBlocks
	for key := range seenEvents {
		if key.BlockNumber > 0 && key.BlockNumber < minBlock {
			delete(seenEvents, key)
		}
	}
}

func advanceRecoveryStartFromLive(nextRecoveryStartBlock *uint64, eventBlock uint64) {
	if *nextRecoveryStartBlock != 0 && eventBlock <= *nextRecoveryStartBlock {
		advanceRecoveryStart(nextRecoveryStartBlock, eventBlock+1)
	}
}

func enqueueEthWriteAcknowledgement(
	batchBuilder *services.BatchBuilder,
	ev *contractICS26Router.ContractICS26RouterWriteAcknowledgement,
	seenEvents map[ethEventKey]struct{},
) bool {
	if !markEthEventSeen(seenEvents, ethEventKeyForLog("WriteAcknowledgement", ev.Raw)) {
		return false
	}
	cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
	batchBuilder.AddEth(services.EthPacket{
		Type:        services.EthWriteAck,
		Packet:      &cosmosPacket,
		AckBytes:    ev.Acknowledgements,
		BlockNumber: ev.Raw.BlockNumber,
	})
	return true
}

func enqueueEthSendPacket(
	batchBuilder *services.BatchBuilder,
	ev *contractICS26Router.ContractICS26RouterSendPacket,
	seenEvents map[ethEventKey]struct{},
) bool {
	if !markEthEventSeen(seenEvents, ethEventKeyForLog("SendPacket", ev.Raw)) {
		return false
	}
	cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
	batchBuilder.AddEth(services.EthPacket{
		Type:        services.EthSend,
		Packet:      &cosmosPacket,
		BlockNumber: ev.Raw.BlockNumber,
	})
	batchBuilder.EthPendingTracker.Add(cosmosPacket, ev.Raw.BlockNumber)
	return true
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

func HasCosmosPacketReceipt(ctx services.Context, packet channeltypesv2.Packet) (bool, error) {
	return hasCosmosIBCPathValue(ctx, utils.IbcPath(packet.DestinationClient, packet.Sequence, []byte{2}))
}

func recoverEthSendPackets(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	filterer *contractICS26Router.ContractICS26RouterFilterer,
	startBlock uint64,
	endBlock uint64,
	seenEvents map[ethEventKey]struct{},
) (ethRecoveryStats, error) {
	var stats ethRecoveryStats
	if endBlock < startBlock {
		return stats, nil
	}

	filterOpts := &bind.FilterOpts{
		Start:   startBlock,
		End:     &endBlock,
		Context: context.Background(),
	}
	iter, err := filterer.FilterSendPacket(filterOpts, ethEventClientIDFilter(ctx), nil)
	if err != nil {
		return stats, fmt.Errorf("failed to filter SendPacket logs in [%d,%d]: %w", startBlock, endBlock, err)
	}
	defer iter.Close()

	var firstErr error
	for iter.Next() {
		ev := iter.Event
		key := ethEventKeyForLog("SendPacket", ev.Raw)
		if seenEvents != nil {
			if _, ok := seenEvents[key]; ok {
				stats.skipped++
				continue
			}
		}

		cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)

		received, err := HasCosmosPacketReceipt(ctx, cosmosPacket)
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] recovery: seq=%d failed to check Cosmos packet receipt: %v",
				cosmosPacket.Sequence, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if received {
			markEthEventSeen(seenEvents, key)
			stats.skipped++
			ctx.Logger.Printf("[SubscribeEth] recovery: seq=%d already received on Cosmos, skipping historical SendPacket from ETH block %d",
				cosmosPacket.Sequence, ev.Raw.BlockNumber)
			continue
		}

		pending, err := services.HasPendingEthPacketCommitment(ctx, cosmosPacket)
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] recovery: seq=%d failed to check ETH packet commitment: %v",
				cosmosPacket.Sequence, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if !pending {
			markEthEventSeen(seenEvents, key)
			batchBuilder.EthPendingTracker.Remove(cosmosPacket.SourceClient, cosmosPacket.Sequence)
			stats.skipped++
			ctx.Logger.Printf("[SubscribeEth] recovery: seq=%d already cleared on ETH, skipping historical SendPacket from ETH block %d",
				cosmosPacket.Sequence, ev.Raw.BlockNumber)
			continue
		}

		if enqueueEthSendPacket(batchBuilder, ev, seenEvents) {
			ctx.Logger.Printf("[SubscribeEth] recovery: recovered SendPacket seq=%d from ETH block %d",
				cosmosPacket.Sequence, ev.Raw.BlockNumber)
			stats.recovered++
		} else {
			stats.skipped++
		}
	}

	if err := iter.Error(); err != nil {
		ctx.Logger.Printf("[SubscribeEth] recovery: SendPacket iterator error in [%d,%d]: %v", startBlock, endBlock, err)
		if firstErr == nil {
			firstErr = err
		}
	}

	if stats.recovered > 0 || stats.skipped > 0 {
		ctx.Logger.Printf("[SubscribeEth] recovery scanned [%d,%d] for SendPacket: recovered=%d skipped=%d",
			startBlock, endBlock, stats.recovered, stats.skipped)
	}
	return stats, firstErr
}

func recoverEthWriteAcknowledgements(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	filterer *contractICS26Router.ContractICS26RouterFilterer,
	startBlock uint64,
	endBlock uint64,
	seenEvents map[ethEventKey]struct{},
) (ethRecoveryStats, error) {
	var stats ethRecoveryStats
	if endBlock < startBlock {
		return stats, nil
	}

	filterOpts := &bind.FilterOpts{
		Start:   startBlock,
		End:     &endBlock,
		Context: context.Background(),
	}
	iter, err := filterer.FilterWriteAcknowledgement(filterOpts, ethEventClientIDFilter(ctx), nil)
	if err != nil {
		return stats, fmt.Errorf("failed to filter WriteAcknowledgement logs in [%d,%d]: %w", startBlock, endBlock, err)
	}
	defer iter.Close()

	var firstErr error
	for iter.Next() {
		ev := iter.Event
		key := ethEventKeyForLog("WriteAcknowledgement", ev.Raw)
		if seenEvents != nil {
			if _, ok := seenEvents[key]; ok {
				stats.skipped++
				continue
			}
		}

		cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)

		pending, err := hasPendingCosmosPacketCommitment(ctx, cosmosPacket)
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] recovery: seq=%d failed to check Cosmos packet commitment: %v",
				cosmosPacket.Sequence, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if !pending {
			markEthEventSeen(seenEvents, key)
			batchBuilder.PendingTracker.Remove(cosmosPacket.SourceClient, cosmosPacket.Sequence)
			stats.skipped++
			ctx.Logger.Printf("[SubscribeEth] recovery: seq=%d already cleared on Cosmos, skipping historical WriteAcknowledgement from ETH block %d",
				cosmosPacket.Sequence, ev.Raw.BlockNumber)
			continue
		}

		if enqueueEthWriteAcknowledgement(batchBuilder, ev, seenEvents) {
			batchBuilder.PendingTracker.Remove(cosmosPacket.SourceClient, cosmosPacket.Sequence)
			ctx.Logger.Printf("[SubscribeEth] recovery: recovered WriteAcknowledgement seq=%d from ETH block %d",
				cosmosPacket.Sequence, ev.Raw.BlockNumber)
			stats.recovered++
		} else {
			stats.skipped++
		}
	}

	if err := iter.Error(); err != nil {
		ctx.Logger.Printf("[SubscribeEth] recovery: WriteAcknowledgement iterator error in [%d,%d]: %v", startBlock, endBlock, err)
		if firstErr == nil {
			firstErr = err
		}
	}

	if stats.recovered > 0 || stats.skipped > 0 {
		ctx.Logger.Printf("[SubscribeEth] recovery scanned [%d,%d] for WriteAcknowledgement: recovered=%d skipped=%d",
			startBlock, endBlock, stats.recovered, stats.skipped)
	}
	return stats, firstErr
}

func advanceRecoveryStart(nextRecoveryStartBlock *uint64, candidate uint64) {
	if candidate > *nextRecoveryStartBlock {
		*nextRecoveryStartBlock = candidate
	}
}

func recoverEthGapToBlock(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	filterer *contractICS26Router.ContractICS26RouterFilterer,
	nextSendRecoveryStartBlock *uint64,
	nextWriteAckRecoveryStartBlock *uint64,
	endBlock uint64,
	seenEvents map[ethEventKey]struct{},
	quietScans *uint64,
) error {
	var firstErr error
	found := false

	if endBlock >= *nextSendRecoveryStartBlock {
		if stats, err := recoverEthSendPackets(ctx, batchBuilder, filterer, *nextSendRecoveryStartBlock, endBlock, seenEvents); err != nil {
			ctx.Logger.Printf("[SubscribeEth] SendPacket recovery failed: %v", err)
			firstErr = err
		} else {
			found = found || stats.foundSomething()
			*nextSendRecoveryStartBlock = endBlock + 1
		}
	}

	if endBlock >= *nextWriteAckRecoveryStartBlock {
		if stats, err := recoverEthWriteAcknowledgements(ctx, batchBuilder, filterer, *nextWriteAckRecoveryStartBlock, endBlock, seenEvents); err != nil {
			ctx.Logger.Printf("[SubscribeEth] WriteAcknowledgement recovery failed: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		} else {
			found = found || stats.foundSomething()
			*nextWriteAckRecoveryStartBlock = endBlock + 1
		}
	}

	// Mirror of the Cosmos side: a scan that found nothing says nothing, except a
	// periodic heartbeat so "gap recovery is alive" stays observable.
	if firstErr == nil {
		if found {
			*quietScans = 0
		} else if beat := advanceQuietScans(quietScans); beat {
			ctx.Logger.Printf("[SubscribeEth] gap recovery healthy: %d consecutive scans found nothing, now current at block %d",
				*quietScans, endBlock)
		}
		pruneEthSeenEvents(seenEvents, endBlock)
	}
	return firstErr
}

func recoverEthGapToLatest(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	filterer *contractICS26Router.ContractICS26RouterFilterer,
	nextSendRecoveryStartBlock *uint64,
	nextWriteAckRecoveryStartBlock *uint64,
	seenEvents map[ethEventKey]struct{},
	quietScans *uint64,
) error {
	latestBlock, err := ctx.EthClient().BlockNumber(context.Background())
	if err != nil {
		return err
	}
	return recoverEthGapToBlock(
		ctx,
		batchBuilder,
		filterer,
		nextSendRecoveryStartBlock,
		nextWriteAckRecoveryStartBlock,
		latestBlock,
		seenEvents,
		quietScans,
	)
}

func (s *Subscriber) subscribeEthOnce(
	ctx services.Context,
	batchBuilder *services.BatchBuilder,
	watchClient *ethclient.Client,
	recoveryFilterer *contractICS26Router.ContractICS26RouterFilterer,
	watchStartBlock uint64,
	nextSendRecoveryStartBlock *uint64,
	nextWriteAckRecoveryStartBlock *uint64,
	seenEvents map[ethEventKey]struct{},
	quietScans *uint64,
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
	clientIDFilter := ethEventClientIDFilter(ctx)

	sendPacketSub, err := watchFilterer.WatchSendPacket(watchOpts, sendPacketCh, clientIDFilter, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to SendPacket events: %w", err)
	}
	defer sendPacketSub.Unsubscribe()

	writeAckSub, err := watchFilterer.WatchWriteAcknowledgement(watchOpts, writeAckCh, clientIDFilter, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to WriteAcknowledgement events: %w", err)
	}
	defer writeAckSub.Unsubscribe()

	ackPacketSub, err := watchFilterer.WatchAckPacket(watchOpts, ackPacketCh, clientIDFilter, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AckPacket events: %w", err)
	}
	defer ackPacketSub.Unsubscribe()

	timeoutPacketSub, err := watchFilterer.WatchTimeoutPacket(watchOpts, timeoutPacketCh, clientIDFilter, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to TimeoutPacket events: %w", err)
	}
	defer timeoutPacketSub.Unsubscribe()

	if len(clientIDFilter) > 0 {
		ctx.Logger.Printf("[SubscribeEth] Successfully subscribed to ICS26Router events from block %d for client_id=%s",
			watchStartBlock, clientIDFilter[0])
	} else {
		ctx.Logger.Printf("[SubscribeEth] Successfully subscribed to ICS26Router events from block %d", watchStartBlock)
	}

	gapRecoveryTicker := time.NewTicker(ethGapRecoveryInterval)
	defer gapRecoveryTicker.Stop()

	for {
		select {
		case ev := <-sendPacketCh:
			ctx.Logger.Printf("SendPacket event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			if enqueueEthSendPacket(batchBuilder, ev, seenEvents) {
				advanceRecoveryStartFromLive(nextSendRecoveryStartBlock, ev.Raw.BlockNumber)
			}

		case ev := <-writeAckCh:
			ctx.Logger.Printf("WriteAcknowledgement event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			if enqueueEthWriteAcknowledgement(batchBuilder, ev, seenEvents) {
				advanceRecoveryStartFromLive(nextWriteAckRecoveryStartBlock, ev.Raw.BlockNumber)
			}
			batchBuilder.PendingTracker.Remove(ev.Packet.SourceClient, ev.Sequence.Uint64())

		case ev := <-ackPacketCh:
			ctx.Logger.Printf("AckPacket event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
			batchBuilder.EthPendingTracker.Remove(cosmosPacket.SourceClient, cosmosPacket.Sequence)
			batchBuilder.AddEth(services.EthPacket{
				Type:     services.EthAck,
				Packet:   &cosmosPacket,
				AckBytes: [][]byte{ev.Acknowledgement},
			})

		case ev := <-timeoutPacketCh:
			ctx.Logger.Printf("TimeoutPacket event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
			batchBuilder.EthPendingTracker.Remove(cosmosPacket.SourceClient, cosmosPacket.Sequence)
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

		case <-gapRecoveryTicker.C:
			if err := recoverEthGapToLatest(
				ctx,
				batchBuilder,
				recoveryFilterer,
				nextSendRecoveryStartBlock,
				nextWriteAckRecoveryStartBlock,
				seenEvents,
				quietScans,
			); err != nil {
				ctx.Logger.Printf("[SubscribeEth] periodic recovery failed: %v", err)
			}
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
	// quietScans counts consecutive recovery passes that found nothing. It lives
	// out here, beside the cursors, so a resubscribe does not reset the heartbeat.
	var quietScans uint64
	seenEvents := make(map[ethEventKey]struct{})

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

		if err := recoverEthGapToBlock(
			ctx,
			batchBuilder,
			recoveryFilterer,
			&nextSendRecoveryStartBlock,
			&nextWriteAckRecoveryStartBlock,
			latestBlock,
			seenEvents,
			&quietScans,
		); err != nil {
			ctx.Logger.Printf("[SubscribeEth] startup recovery failed: %v", err)
		}

		watchStartBlock := latestBlock + 1

		watchClient, err := ethclient.DialContext(context.Background(), ctx.EthWsURL())
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] Failed to connect to Ethereum WS at %s: %v", ctx.EthWsURL(), err)
			time.Sleep(ethSubscriptionReconnectDelay)
			continue
		}

		err = s.subscribeEthOnce(
			ctx,
			batchBuilder,
			watchClient,
			recoveryFilterer,
			watchStartBlock,
			&nextSendRecoveryStartBlock,
			&nextWriteAckRecoveryStartBlock,
			seenEvents,
			&quietScans,
		)
		watchClient.Close()
		ctx.Logger.Printf("[SubscribeEth] Subscription loop ended: %v", err)
		time.Sleep(ethSubscriptionReconnectDelay)
	}
}
