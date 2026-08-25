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
const subscriptionCleanupTimeout = 5 * time.Second
const subscriberRPCTimeout = 15 * time.Second
const cosmosRecoveryChunkEnv = "COSMOS_RECOVERY_CHUNK_HEIGHTS"
const defaultCosmosRecoveryChunkHeights uint64 = 256

// cosmosLiveEventBuffer is the per-subscription channel depth requested from
// CometBFT. Passing it is not a tuning nicety -- the default is 1, and a full
// channel makes CometBFT DROP the event:
//
//	// cometbft rpc/client/http/http.go, eventListener
//	select {
//	case out <- *result:
//	default:
//	    w.Logger.Error("wanted to publish ResultEvent, but out channel is full", ...)
//	}
//
// The send is non-blocking, so there is no backpressure, and w.Logger is a nop
// logger unless SetLogger is called -- so a dropped packet leaves no trace
// anywhere. Meanwhile this subscriber's consuming goroutine also runs the 30s
// gap-recovery scan (several paginated TxSearch calls) in the same select, so it
// stops reading for seconds at a time, every 30 seconds.
//
// The Ethereum side does not need this: go-ethereum's generated watcher sends
// blocking and its RPC client queues 20,000 events before failing loudly with
// ErrSubscriptionQueueOverflow, which subscribeEthOnce already surfaces.
//
// A buffer is a mitigation, not a fix -- it still drops once full. The fix is to
// keep slow work off the consuming goroutine; see #376.
const cosmosLiveEventBuffer = 1024

// cosmosLivePathStaleAfter is how long the live subscription may deliver nothing while
// gap recovery is finding packets before that combination is reported.
//
// Recovery is the backstop; the live subscription is meant to be what delivers.
// When recovery starts finding packets the live path should have delivered, the
// live path is failing -- and it fails silently, because a subscription that
// delivers nothing looks exactly like a quiet chain. That is how a Cosmos->Base
// packet was lost unnoticed: every packet for 14 minutes arrived via recovery,
// and nothing said so.
const cosmosLivePathStaleAfter = 2 * time.Minute

// cosmosIndexerLagBlocks is how far behind the head gap recovery stops scanning.
//
// The scan reads two different subsystems and assumes they agree. The range end
// comes from Status.SyncInfo.LatestBlockHeight, which is set when the block
// commits; the results come from TxSearch, which reads the tx indexer -- and the
// indexer is a SEPARATE goroutine consuming the event bus after commit
// (cometbft state/txindex/indexer_service.go). So a tx in the newest block is
// routinely committed but not yet indexed.
//
// That race loses the packet permanently, because the scan then advances the
// cursor past the height it just failed to read:
//
//	latest := Status()                    // 202195, committed
//	TxSearch(..., "tx.height <= 202195")  // indexer has not written 202195 yet -> empty
//	cursor = 202195 + 1                   // 202195 is never scanned again
//
// and it is silent: an empty scan increments nothing, so nothing is logged.
// Observed on a live devnet -- one Cosmos->Base transfer at height 202195 was
// never relayed and never timed out, leaving its funds escrowed.
//
// Two blocks is far more than the indexer needs (its lag is milliseconds) and
// costs the backstop about ten seconds of extra detection delay. Rescanning is
// free: seenEvents dedups, and recovery re-checks on-chain state anyway.
const cosmosIndexerLagBlocks uint64 = 2

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

const ethRecoveryChunkEnv = "ETH_RECOVERY_CHUNK_BLOCKS"
const defaultEthRecoveryChunkBlocks uint64 = 256

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
	recovery *services.RecoveryStateStore
}

type cosmosDeps struct {
	Cosmos services.CosmosEndpoint
	EVM    services.EVMEndpoint
	IDs    services.ClientIDs
	Logger *log.Logger
}

type ethDeps = cosmosDeps

func NewSubscriber(recovery ...*services.RecoveryStateStore) *Subscriber {
	var store *services.RecoveryStateStore
	if len(recovery) > 0 {
		store = recovery[0]
	}
	return &Subscriber{recovery: store}
}

func recoverySourceID(ids services.ClientIDs) string {
	if ids.CosmosOnEVM != "" {
		return ids.CosmosOnEVM
	}
	return "default"
}

func clampRecoveryCursor(cursor, lowestUnsubmitted uint64) uint64 {
	if lowestUnsubmitted != 0 && lowestUnsubmitted < cursor {
		return lowestUnsubmitted
	}
	return cursor
}

func recoveryChunkSize(envName string, fallback uint64) uint64 {
	if raw := os.Getenv(envName); raw != "" {
		value, err := strconv.ParseUint(raw, 10, 64)
		if err == nil && value > 0 {
			return value
		}
		log.Printf("[recovery] ignoring invalid %s=%q; using %d", envName, raw, fallback)
	}
	return fallback
}

func sleepOrDone(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func unsubscribeCosmos(parent context.Context, cosmos services.CosmosEndpoint) {
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(parent), subscriptionCleanupTimeout)
	defer cancel()
	_ = cosmos.CosmosClient().UnsubscribeAll(cleanupCtx, "")
}

func (s *Subscriber) SubscribeCosmos(stdCtx context.Context, cosmos services.CosmosEndpoint, evm services.EVMEndpoint, ids services.ClientIDs, logger *log.Logger, batchBuilder *services.BatchBuilder) {
	ctx := cosmosDeps{Cosmos: cosmos, EVM: evm, IDs: ids, Logger: logger}
	lookback := cosmosStartupRecoveryLookbackBlocks()
	var nextRecoveryStartHeight uint64
	// quietScans counts consecutive recovery passes that found nothing. It lives
	// beside the cursor, outside the resubscribe loop, so a WS reconnect does not
	// restart the heartbeat — mirroring the Ethereum side.
	var quietScans uint64
	// Outside the resubscribe loop: a reconnect must not reset the picture of
	// whether the live path has been delivering.
	var liveHealth cosmosLiveHealth
	seenEvents := make(map[cosmosEventKey]struct{})

	for {
		if stdCtx.Err() != nil {
			return
		}
		if nextRecoveryStartHeight == 0 {
			latestHeight, err := latestCosmosHeight(stdCtx, ctx)
			if err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] Failed to get latest Cosmos height before subscription: %v", err)
				if !sleepOrDone(stdCtx, cosmosSubscriptionReconnectDelay) {
					return
				}
				continue
			}
			nextRecoveryStartHeight = s.resumeCosmosCursor(ctx, latestHeight, lookback)
		}

		err := s.subscribeCosmosOnce(stdCtx, ctx, batchBuilder, &nextRecoveryStartHeight, seenEvents, &quietScans, &liveHealth)
		ctx.Logger.Printf("[SubscribeCosmos] Subscription loop ended: %v", err)
		s.persistCosmosCursor(ctx, batchBuilder, nextRecoveryStartHeight)
		if !sleepOrDone(stdCtx, cosmosSubscriptionReconnectDelay) {
			return
		}
	}
}

func (s *Subscriber) resumeCosmosCursor(ctx cosmosDeps, latestHeight, lookback uint64) uint64 {
	lookbackStart := cosmosStartupRecoveryStartHeight(latestHeight, lookback)
	if s.recovery == nil {
		return lookbackStart
	}
	cursors, ok := s.recovery.Get(recoverySourceID(ctx.IDs))
	if !ok {
		return lookbackStart
	}
	start, resumed, behind := services.ResumeCursor(cursors.CosmosHeight, lookbackStart, latestHeight)
	if resumed && behind > 0 {
		ctx.Logger.Printf("[SubscribeCosmos] resuming from persisted height %d, %d block(s) behind head %d", start, behind, latestHeight)
	} else if resumed {
		ctx.Logger.Printf("[SubscribeCosmos] persisted cursor %d is inside the lookback window; rescanning from %d", cursors.CosmosHeight, start)
	}
	return start
}

func (s *Subscriber) persistCosmosCursor(ctx cosmosDeps, batchBuilder *services.BatchBuilder, height uint64) {
	if s.recovery == nil || height == 0 {
		return
	}
	if batchBuilder != nil {
		height = clampRecoveryCursor(height, batchBuilder.LowestUnsubmittedCosmosHeight())
	}
	if height == 0 {
		return
	}
	if err := s.recovery.SaveCheckpoint(recoverySourceID(ctx.IDs), services.RecoveryCursors{CosmosHeight: height}); err != nil {
		ctx.Logger.Printf("[SubscribeCosmos] failed to persist recovery cursor (height %d): %v", height, err)
	}
}

// cosmosLiveHealth tracks whether the live subscription is actually delivering.
//
// It exists because the failure it detects is invisible by construction: a
// subscription that stops delivering looks exactly like a chain with no traffic.
// Gap recovery quietly picks up the slack, the relayer keeps working, and
// nothing says the primary path is dead -- until recovery misses one too and a
// packet is lost.
//
// Owned by one goroutine (the subscribe loop), so it needs no lock.
type cosmosLiveHealth struct {
	events   uint64
	lastSeen time.Time
	warned   bool
}

func (h *cosmosLiveHealth) recordEvent() {
	h.events++
	h.lastSeen = time.Now()
	h.warned = false
}

// recoveryIsCoveringForLive reports whether gap recovery just relayed packets the
// live subscription should have delivered. It says so once per outage rather
// than on every scan, so a long outage does not flood the log.
func (h *cosmosLiveHealth) recoveryIsCoveringForLive(recovered uint64, now time.Time) bool {
	if recovered == 0 || h.warned {
		return false
	}
	// Before the first live event there is nothing to compare against; treat
	// process start as the reference point so a subscription that never
	// delivers is still reported.
	if h.lastSeen.IsZero() {
		h.lastSeen = now
		return false
	}
	if now.Sub(h.lastSeen) < cosmosLivePathStaleAfter {
		return false
	}
	h.warned = true
	return true
}

func (s *Subscriber) subscribeCosmosOnce(
	stdCtx context.Context,
	ctx cosmosDeps,
	batchBuilder *services.BatchBuilder,
	nextRecoveryStartHeight *uint64,
	seenEvents map[cosmosEventKey]struct{},
	quietScans *uint64,
	liveHealth *cosmosLiveHealth,
) error {
	if err := stdCtx.Err(); err != nil {
		return err
	}
	sendPacketSub, err := ctx.Cosmos.CosmosClient().WSEvents.Subscribe(stdCtx, "", COMETBFT_SEND_PACKET_EVENT, cosmosLiveEventBuffer)
	if err != nil {
		return fmt.Errorf("failed to subscribe to send_packet events: %w", err)
	}
	ackPacketSub, err := ctx.Cosmos.CosmosClient().WSEvents.Subscribe(stdCtx, "", COMETBFT_WRITE_ACK_PACKET_EVENT, cosmosLiveEventBuffer)
	if err != nil {
		unsubscribeCosmos(stdCtx, ctx.Cosmos)
		return fmt.Errorf("failed to subscribe to write_acknowledgement events: %w", err)
	}
	timeoutPacketSub, err := ctx.Cosmos.CosmosClient().WSEvents.Subscribe(stdCtx, "", COMETBFT_TIMEOUT_PACKET_EVENT, cosmosLiveEventBuffer)
	if err != nil {
		unsubscribeCosmos(stdCtx, ctx.Cosmos)
		return fmt.Errorf("failed to subscribe to timeout_packet events: %w", err)
	}
	ctx.Logger.Println("[SubscribeCosmos] Successfully subscribed to CometBFT events")
	defer unsubscribeCosmos(stdCtx, ctx.Cosmos)

	if err := s.recoverCosmosGapToLatest(stdCtx, ctx, batchBuilder, nextRecoveryStartHeight, seenEvents, quietScans, liveHealth); err != nil {
		ctx.Logger.Printf("[SubscribeCosmos] startup recovery failed: %v", err)
	}

	gapRecoveryTicker := time.NewTicker(cosmosGapRecoveryInterval)
	defer gapRecoveryTicker.Stop()

	for {
		select {
		case <-stdCtx.Done():
			return stdCtx.Err()
		case e, ok := <-sendPacketSub:
			if !ok {
				return fmt.Errorf("send_packet subscription channel closed")
			}
			s.processLiveCosmosEvent(stdCtx, ctx, batchBuilder, nextRecoveryStartHeight, seenEvents, liveHealth, e)
		case e, ok := <-ackPacketSub:
			if !ok {
				return fmt.Errorf("write_acknowledgement subscription channel closed")
			}
			s.processLiveCosmosEvent(stdCtx, ctx, batchBuilder, nextRecoveryStartHeight, seenEvents, liveHealth, e)
		case e, ok := <-timeoutPacketSub:
			if !ok {
				return fmt.Errorf("timeout_packet subscription channel closed")
			}
			s.processLiveCosmosEvent(stdCtx, ctx, batchBuilder, nextRecoveryStartHeight, seenEvents, liveHealth, e)
		case <-gapRecoveryTicker.C:
			if err := s.recoverCosmosGapToLatest(stdCtx, ctx, batchBuilder, nextRecoveryStartHeight, seenEvents, quietScans, liveHealth); err != nil {
				ctx.Logger.Printf("[SubscribeCosmos] periodic recovery failed: %v", err)
			}
		}
	}
}

func (s *Subscriber) processLiveCosmosEvent(
	stdCtx context.Context,
	ctx cosmosDeps,
	batchBuilder *services.BatchBuilder,
	nextRecoveryStartHeight *uint64,
	seenEvents map[cosmosEventKey]struct{},
	liveHealth *cosmosLiveHealth,
	e coretypes.ResultEvent,
) {
	height := txHeightFromEvent(e.Data, e.Events)
	if height > 0 && *nextRecoveryStartHeight > 0 && height > *nextRecoveryStartHeight {
		stats, err := recoverCosmosEvents(stdCtx, ctx, batchBuilder, *nextRecoveryStartHeight, height-1, seenEvents)
		if err != nil {
			ctx.Logger.Printf("[SubscribeCosmos] gap recovery before live height %d failed: %v", height, err)
		} else {
			ctx.Logger.Printf("[SubscribeCosmos] gap recovery complete before live height %d: recovered=%d skipped=%d",
				height, stats.recovered, stats.skipped)
			*nextRecoveryStartHeight = height
			s.persistCosmosCursor(ctx, batchBuilder, *nextRecoveryStartHeight)
		}
	}

	liveHealth.recordEvent()

	packets := decodeCosmosPacketsFromEvents(ctx.Logger, e.Data, e.Events, "SubscribeCosmos")
	stats, err := enqueueCosmosPackets(stdCtx, ctx, batchBuilder, packets, seenEvents, false)
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

func latestCosmosHeight(stdCtx context.Context, ctx cosmosDeps) (uint64, error) {
	rpcCtx, cancel := context.WithTimeout(stdCtx, subscriberRPCTimeout)
	defer cancel()
	status, err := ctx.Cosmos.CosmosClient().Status(rpcCtx)
	if err != nil {
		return 0, err
	}
	if status.SyncInfo.LatestBlockHeight <= 0 {
		return 0, nil
	}
	return uint64(status.SyncInfo.LatestBlockHeight), nil
}

func (s *Subscriber) recoverCosmosGapToLatest(
	stdCtx context.Context,
	ctx cosmosDeps,
	batchBuilder *services.BatchBuilder,
	nextRecoveryStartHeight *uint64,
	seenEvents map[cosmosEventKey]struct{},
	quietScans *uint64,
	liveHealth *cosmosLiveHealth,
) error {
	if *nextRecoveryStartHeight == 0 {
		return nil
	}

	latestHeight, err := latestCosmosHeight(stdCtx, ctx)
	if err != nil {
		return err
	}
	scanFrom, scanTo, ok := cosmosScanRange(*nextRecoveryStartHeight, latestHeight)
	if !ok {
		return nil
	}
	passStart := scanFrom

	var stats cosmosRecoveryStats
	chunkSize := recoveryChunkSize(cosmosRecoveryChunkEnv, defaultCosmosRecoveryChunkHeights)
	for scanFrom <= scanTo {
		chunkEnd := scanFrom + chunkSize - 1
		if chunkEnd < scanFrom || chunkEnd > scanTo {
			chunkEnd = scanTo
		}
		chunkStats, err := recoverCosmosEvents(stdCtx, ctx, batchBuilder, scanFrom, chunkEnd, seenEvents)
		if err != nil {
			return err
		}
		stats.recovered += chunkStats.recovered
		stats.skipped += chunkStats.skipped
		*nextRecoveryStartHeight = chunkEnd + 1
		s.persistCosmosCursor(ctx, batchBuilder, *nextRecoveryStartHeight)
		pruneCosmosSeenEvents(seenEvents, chunkEnd)
		scanFrom = *nextRecoveryStartHeight
	}
	// Report a scan that FOUND something, and otherwise only a periodic heartbeat.
	// This ticks every 30s for the life of the process, and the overwhelmingly
	// common outcome is recovered=0 skipped=0 — two lines a tick, ~5.8k lines a day
	// per direction, burying the events worth reading. The heartbeat keeps "gap
	// recovery is alive and current" observable without the repetition.
	if liveHealth.recoveryIsCoveringForLive(stats.recovered, time.Now()) {
		ctx.Logger.Printf("[SubscribeCosmos][ATTENTION] gap recovery relayed %d packet(s) but the live "+
			"subscription has delivered nothing for %s (%d live event(s) this run). Recovery is the "+
			"BACKSTOP, not the delivery path -- while it is doing this job a packet in the newest blocks "+
			"can still be missed. Check the CometBFT websocket.",
			stats.recovered, cosmosLivePathStaleAfter, liveHealth.events)
	}

	if stats.recovered > 0 || stats.skipped > 0 {
		ctx.Logger.Printf("[SubscribeCosmos] recovery scanned [%d,%d]: recovered=%d skipped=%d",
			passStart, scanTo, stats.recovered, stats.skipped)
		*quietScans = 0
	} else if beat := advanceQuietScans(quietScans); beat {
		ctx.Logger.Printf("[SubscribeCosmos] gap recovery healthy: %d consecutive scans found nothing, now current at height %d",
			*quietScans, scanTo)
	}

	return nil
}

// cosmosScanRange decides the height range one recovery pass covers, given the
// cursor and the chain head. ok is false when nothing can be scanned safely yet.
//
// The whole decision lives here, in one pure function, because the property that
// matters is not "the arithmetic is right" but "the scan never reads up to the
// head" -- and a helper the caller could stop calling would let that property be
// removed without a test noticing.
func cosmosScanRange(cursor, latestHeight uint64) (from, to uint64, ok bool) {
	if cursor == 0 {
		return 0, 0, false
	}
	to = cosmosIndexedHeight(latestHeight)
	if to < cursor {
		return 0, 0, false
	}
	return cursor, to, true
}

// cosmosIndexedHeight is the newest height gap recovery may scan: far enough
// behind the head that the tx indexer has certainly caught up. Returns 0 on a
// chain too short to have one, which makes the caller wait.
func cosmosIndexedHeight(latestHeight uint64) uint64 {
	if latestHeight <= cosmosIndexerLagBlocks {
		return 0
	}
	return latestHeight - cosmosIndexerLagBlocks
}

func recoverCosmosEvents(
	stdCtx context.Context,
	ctx cosmosDeps,
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
		stats, err := recoverCosmosEventsForQuery(stdCtx, ctx, batchBuilder, query, startHeight, endHeight, seenEvents)
		combined.recovered += stats.recovered
		combined.skipped += stats.skipped
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return combined, firstErr
}

func recoverCosmosEventsForQuery(
	stdCtx context.Context,
	ctx cosmosDeps,
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
		rpcCtx, cancel := context.WithTimeout(stdCtx, subscriberRPCTimeout)
		result, err := ctx.Cosmos.CosmosClient().TxSearch(rpcCtx, query, false, &page, &perPage, "asc")
		cancel()
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
			enqueued, err := enqueueCosmosPackets(stdCtx, ctx, batchBuilder, packets, seenEvents, true)
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
	stdCtx context.Context,
	ctx cosmosDeps,
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
			shouldEnqueue, err := shouldEnqueueRecoveredCosmosPacket(stdCtx, ctx, packet)
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

func cosmosPacketMatchesConfiguredClient(ctx cosmosDeps, packet *channeltypesv2.Packet) bool {
	if packet == nil {
		return false
	}

	ethClientID := ctx.IDs.EVMOnCosmos
	routerClientID := ctx.IDs.CosmosOnEVM
	if ethClientID == "" && routerClientID == "" {
		return true
	}

	return packet.SourceClient == ethClientID ||
		packet.DestinationClient == ethClientID ||
		packet.SourceClient == routerClientID ||
		packet.DestinationClient == routerClientID
}

func shouldEnqueueRecoveredCosmosPacket(stdCtx context.Context, ctx cosmosDeps, packet services.CosmosPacket) (bool, error) {
	if packet.Packet == nil {
		return false, nil
	}

	switch packet.Type {
	case services.CosmosSend:
		received, err := services.HasEthPacketReceipt(stdCtx, ctx.EVM, *packet.Packet)
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
		pending, err := services.HasPendingEthPacketCommitment(stdCtx, ctx.EVM, *packet.Packet)
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
		if !services.ShouldRelayCosmosTimeoutToEth(packet.Packet, ctx.IDs.CosmosOnEVM) {
			ctx.Logger.Printf("[SubscribeCosmos] recovery: seq=%d Cosmos-originated timeout already handled locally, skipping historical TimeoutPacket from Cosmos height %d",
				packet.Packet.Sequence, packet.BlockNumber)
			return false, nil
		}
		pending, err := services.HasPendingEthPacketCommitment(stdCtx, ctx.EVM, *packet.Packet)
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

func ethEventClientIDFilter(ctx ethDeps) []string {
	clientID := ctx.IDs.CosmosOnEVM
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

// enqueueEthTerminal records an AckPacket or TimeoutPacket — the two events that
// SETTLE an ETH-origin send rather than starting new work. It removes the packet
// from the pending tracker (it can no longer time out) and queues it so the
// source bridge can observe the settlement.
//
// Extracted from the subscribe select so BlockNumber is set in one place. The
// two enqueues used to be written inline and both omitted it, which is harmless
// only for as long as nothing downstream reads Height for these types — a trap
// rather than a bug, and one worth closing where it can be tested.
func enqueueEthTerminal(
	batchBuilder *services.BatchBuilder,
	packetType services.EthPacketType,
	packet contractICS26Router.IICS26RouterMsgsPacket,
	sequence *big.Int,
	ackBytes [][]byte,
	blockNumber uint64,
) {
	cosmosPacket := EthPacketToCosmosPacket(packet, sequence)
	batchBuilder.EthPendingTracker.RemovePacketIfCurrent(cosmosPacket)
	batchBuilder.AddEth(services.EthPacket{
		Type:        packetType,
		Packet:      &cosmosPacket,
		AckBytes:    ackBytes,
		BlockNumber: blockNumber,
	})
}

func enqueueEthSendPacket(
	batchBuilder *services.BatchBuilder,
	ev *contractICS26Router.ContractICS26RouterSendPacket,
	seenEvents map[ethEventKey]struct{},
) bool {
	cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
	if !batchBuilder.EthPendingTracker.Add(cosmosPacket, ev.Raw.BlockNumber) {
		return false
	}
	if !markEthEventSeen(seenEvents, ethEventKeyForLog("SendPacket", ev.Raw)) {
		return false
	}
	batchBuilder.AddEth(services.EthPacket{
		Type:        services.EthSend,
		Packet:      &cosmosPacket,
		BlockNumber: ev.Raw.BlockNumber,
	})
	return true
}

func hasCosmosIBCPathValue(stdCtx context.Context, endpoint services.CosmosEndpoint, path [][]byte) (bool, error) {
	queryPath := fmt.Sprintf("store/%s/key", string(path[0]))
	request := path[1]

	rpcCtx, cancel := context.WithTimeout(stdCtx, subscriberRPCTimeout)
	defer cancel()
	result, err := endpoint.CosmosClient().ABCIQuery(rpcCtx, queryPath, request)
	if err != nil {
		return false, fmt.Errorf("ABCI query failed: %w", err)
	}
	if result.Response.Code != 0 {
		return false, fmt.Errorf("query failed with code %d: %s", result.Response.Code, result.Response.Log)
	}

	return len(result.Response.Value) > 0, nil
}

func hasPendingCosmosPacketCommitment(stdCtx context.Context, ctx ethDeps, packet channeltypesv2.Packet) (bool, error) {
	return hasCosmosIBCPathValue(stdCtx, ctx.Cosmos, utils.IbcCommitmentPath(packet, []byte{1}))
}

// HasCosmosPacketReceipt is retained for one-shot callers. Long-running relay
// paths must call HasCosmosPacketReceiptWithContext.
func HasCosmosPacketReceipt(endpoint services.CosmosEndpoint, packet channeltypesv2.Packet) (bool, error) {
	return HasCosmosPacketReceiptWithContext(context.Background(), endpoint, packet)
}

func HasCosmosPacketReceiptWithContext(stdCtx context.Context, endpoint services.CosmosEndpoint, packet channeltypesv2.Packet) (bool, error) {
	return hasCosmosIBCPathValue(stdCtx, endpoint, utils.IbcPath(packet.DestinationClient, packet.Sequence, []byte{2}))
}

func recoverEthSendPackets(
	stdCtx context.Context,
	ctx ethDeps,
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
	filterCtx, cancel := context.WithTimeout(stdCtx, subscriberRPCTimeout)
	defer cancel()

	filterOpts := &bind.FilterOpts{
		Start:   startBlock,
		End:     &endBlock,
		Context: filterCtx,
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

		received, err := HasCosmosPacketReceiptWithContext(stdCtx, ctx.Cosmos, cosmosPacket)
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

		pending, err := services.HasPendingEthPacketCommitment(stdCtx, ctx.EVM, cosmosPacket)
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
			batchBuilder.EthPendingTracker.RemovePacketIfCurrent(cosmosPacket)
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
	stdCtx context.Context,
	ctx ethDeps,
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
	filterCtx, cancel := context.WithTimeout(stdCtx, subscriberRPCTimeout)
	defer cancel()

	filterOpts := &bind.FilterOpts{
		Start:   startBlock,
		End:     &endBlock,
		Context: filterCtx,
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

		pending, err := hasPendingCosmosPacketCommitment(stdCtx, ctx, cosmosPacket)
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
			batchBuilder.PendingTracker.RemovePacketIfCurrent(cosmosPacket)
			stats.skipped++
			ctx.Logger.Printf("[SubscribeEth] recovery: seq=%d already cleared on Cosmos, skipping historical WriteAcknowledgement from ETH block %d",
				cosmosPacket.Sequence, ev.Raw.BlockNumber)
			continue
		}

		if enqueueEthWriteAcknowledgement(batchBuilder, ev, seenEvents) {
			batchBuilder.PendingTracker.RemovePacketIfCurrent(cosmosPacket)
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
	// Only historical scans move this cursor. Live subscriptions are not a
	// substitute for gap recovery: advancing past a reconnect gap would make the
	// missed logs permanently invisible.
	if candidate > *nextRecoveryStartBlock {
		*nextRecoveryStartBlock = candidate
	}
}

func (s *Subscriber) scanEthRangeInChunks(
	stdCtx context.Context,
	ctx ethDeps,
	label string,
	cursor *uint64,
	endBlock uint64,
	chunkSize uint64,
	scan func(from, to uint64) (ethRecoveryStats, error),
	persist func(),
) (ethRecoveryStats, error) {
	var combined ethRecoveryStats
	for *cursor <= endBlock {
		if err := stdCtx.Err(); err != nil {
			return combined, err
		}
		to := *cursor + chunkSize - 1
		if to < *cursor || to > endBlock {
			to = endBlock
		}
		stats, err := scan(*cursor, to)
		combined.recovered += stats.recovered
		combined.skipped += stats.skipped
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] %s recovery failed at [%d,%d]: %v", label, *cursor, to, err)
			return combined, err
		}
		*cursor = to + 1
		persist()
	}
	return combined, nil
}

func (s *Subscriber) recoverEthGapToBlock(
	stdCtx context.Context,
	ctx ethDeps,
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
	chunkSize := recoveryChunkSize(ethRecoveryChunkEnv, defaultEthRecoveryChunkBlocks)
	persist := func() {
		s.persistEthCursors(ctx, batchBuilder, *nextSendRecoveryStartBlock, *nextWriteAckRecoveryStartBlock)
	}

	if endBlock >= *nextSendRecoveryStartBlock {
		stats, err := s.scanEthRangeInChunks(stdCtx, ctx, "SendPacket", nextSendRecoveryStartBlock, endBlock, chunkSize,
			func(from, to uint64) (ethRecoveryStats, error) {
				return recoverEthSendPackets(stdCtx, ctx, batchBuilder, filterer, from, to, seenEvents)
			}, persist)
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] SendPacket recovery failed: %v", err)
			firstErr = err
		} else {
			found = found || stats.foundSomething()
		}
	}

	if endBlock >= *nextWriteAckRecoveryStartBlock {
		stats, err := s.scanEthRangeInChunks(stdCtx, ctx, "WriteAcknowledgement", nextWriteAckRecoveryStartBlock, endBlock, chunkSize,
			func(from, to uint64) (ethRecoveryStats, error) {
				return recoverEthWriteAcknowledgements(stdCtx, ctx, batchBuilder, filterer, from, to, seenEvents)
			}, persist)
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] WriteAcknowledgement recovery failed: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		} else {
			found = found || stats.foundSomething()
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

func (s *Subscriber) recoverEthGapToLatest(
	stdCtx context.Context,
	ctx ethDeps,
	batchBuilder *services.BatchBuilder,
	filterer *contractICS26Router.ContractICS26RouterFilterer,
	nextSendRecoveryStartBlock *uint64,
	nextWriteAckRecoveryStartBlock *uint64,
	seenEvents map[ethEventKey]struct{},
	quietScans *uint64,
) error {
	rpcCtx, cancel := context.WithTimeout(stdCtx, subscriberRPCTimeout)
	defer cancel()
	latestBlock, err := ctx.EVM.EthClient().BlockNumber(rpcCtx)
	if err != nil {
		return err
	}
	return s.recoverEthGapToBlock(
		stdCtx,
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
	stdCtx context.Context,
	ctx ethDeps,
	batchBuilder *services.BatchBuilder,
	watchClient *ethclient.Client,
	recoveryFilterer *contractICS26Router.ContractICS26RouterFilterer,
	watchStartBlock uint64,
	nextSendRecoveryStartBlock *uint64,
	nextWriteAckRecoveryStartBlock *uint64,
	seenEvents map[ethEventKey]struct{},
	quietScans *uint64,
) error {
	watchFilterer, err := contractICS26Router.NewContractICS26RouterFilterer(*ctx.EVM.RouterContract(), watchClient)
	if err != nil {
		return fmt.Errorf("failed to create ICS26Router watch filterer instance: %w", err)
	}

	sendPacketCh := make(chan *contractICS26Router.ContractICS26RouterSendPacket)
	writeAckCh := make(chan *contractICS26Router.ContractICS26RouterWriteAcknowledgement)
	ackPacketCh := make(chan *contractICS26Router.ContractICS26RouterAckPacket)
	timeoutPacketCh := make(chan *contractICS26Router.ContractICS26RouterTimeoutPacket)

	watchOpts := &bind.WatchOpts{Start: &watchStartBlock, Context: stdCtx}
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
		case <-stdCtx.Done():
			return stdCtx.Err()
		case ev, ok := <-sendPacketCh:
			if !ok || ev == nil {
				return fmt.Errorf("SendPacket event channel closed")
			}
			ctx.Logger.Printf("SendPacket event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			enqueueEthSendPacket(batchBuilder, ev, seenEvents)

		case ev, ok := <-writeAckCh:
			if !ok || ev == nil {
				return fmt.Errorf("WriteAcknowledgement event channel closed")
			}
			ctx.Logger.Printf("WriteAcknowledgement event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			enqueueEthWriteAcknowledgement(batchBuilder, ev, seenEvents)
			cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
			batchBuilder.PendingTracker.RemovePacketIfCurrent(cosmosPacket)

		case ev, ok := <-ackPacketCh:
			if !ok || ev == nil {
				return fmt.Errorf("AckPacket event channel closed")
			}
			ctx.Logger.Printf("AckPacket event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			enqueueEthTerminal(batchBuilder, services.EthAck, ev.Packet, ev.Sequence,
				[][]byte{ev.Acknowledgement}, ev.Raw.BlockNumber)

		case ev, ok := <-timeoutPacketCh:
			if !ok || ev == nil {
				return fmt.Errorf("TimeoutPacket event channel closed")
			}
			ctx.Logger.Printf("TimeoutPacket event received: clientId=%x, sequence=%s", ev.ClientId, ev.Sequence.String())
			enqueueEthTerminal(batchBuilder, services.EthTimeout, ev.Packet, ev.Sequence,
				nil, ev.Raw.BlockNumber)

		case err := <-sendPacketSub.Err():
			return fmt.Errorf("SendPacket subscription error: %w", err)

		case err := <-writeAckSub.Err():
			return fmt.Errorf("WriteAcknowledgement subscription error: %w", err)

		case err := <-ackPacketSub.Err():
			return fmt.Errorf("AckPacket subscription error: %w", err)

		case err := <-timeoutPacketSub.Err():
			return fmt.Errorf("TimeoutPacket subscription error: %w", err)

		case <-gapRecoveryTicker.C:
			if err := s.recoverEthGapToLatest(
				stdCtx,
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
func (s *Subscriber) SubscribeEth(stdCtx context.Context, cosmos services.CosmosEndpoint, evm services.EVMEndpoint, ids services.ClientIDs, logger *log.Logger, batchBuilder *services.BatchBuilder) {
	ctx := ethDeps{Cosmos: cosmos, EVM: evm, IDs: ids, Logger: logger}
	if stdCtx.Err() != nil {
		return
	}
	if ctx.EVM.EthWsURL() == "" {
		ctx.Logger.Printf("Failed to subscribe to Ethereum events: eth websocket URL is not configured")
		return
	}

	recoveryFilterer, err := contractICS26Router.NewContractICS26RouterFilterer(*ctx.EVM.RouterContract(), ctx.EVM.EthClient())
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
		if stdCtx.Err() != nil {
			return
		}
		rpcCtx, cancel := context.WithTimeout(stdCtx, subscriberRPCTimeout)
		latestBlock, err := ctx.EVM.EthClient().BlockNumber(rpcCtx)
		cancel()
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] Failed to get latest Ethereum block before subscription: %v", err)
			if !sleepOrDone(stdCtx, ethSubscriptionReconnectDelay) {
				return
			}
			continue
		}

		if nextSendRecoveryStartBlock == 0 || nextWriteAckRecoveryStartBlock == 0 {
			send, writeAck := s.resumeEthCursors(ctx, latestBlock, lookback)
			if nextSendRecoveryStartBlock == 0 {
				nextSendRecoveryStartBlock = send
			}
			if nextWriteAckRecoveryStartBlock == 0 {
				nextWriteAckRecoveryStartBlock = writeAck
			}
		}

		if err := s.recoverEthGapToBlock(
			stdCtx,
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

		if stdCtx.Err() != nil {
			return
		}
		dialCtx, cancelDial := context.WithTimeout(stdCtx, subscriberRPCTimeout)
		watchClient, err := ethclient.DialContext(dialCtx, ctx.EVM.EthWsURL())
		cancelDial()
		if err != nil {
			ctx.Logger.Printf("[SubscribeEth] Failed to connect to Ethereum WS at %s: %v", ctx.EVM.EthWsURL(), err)
			if !sleepOrDone(stdCtx, ethSubscriptionReconnectDelay) {
				return
			}
			continue
		}

		err = s.subscribeEthOnce(
			stdCtx,
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
		s.persistEthCursors(ctx, batchBuilder, nextSendRecoveryStartBlock, nextWriteAckRecoveryStartBlock)
		if !sleepOrDone(stdCtx, ethSubscriptionReconnectDelay) {
			return
		}
	}
}

func (s *Subscriber) resumeEthCursors(ctx ethDeps, latestBlock, lookback uint64) (send, writeAck uint64) {
	lookbackStart := ethStartupRecoveryStartBlock(latestBlock, lookback)
	if s.recovery == nil {
		return lookbackStart, lookbackStart
	}
	cursors, ok := s.recovery.Get(recoverySourceID(ctx.IDs))
	if !ok {
		return lookbackStart, lookbackStart
	}
	send, sendResumed, sendBehind := services.ResumeCursor(cursors.EthSendBlock, lookbackStart, latestBlock)
	writeAck, ackResumed, ackBehind := services.ResumeCursor(cursors.EthWriteAckBlk, lookbackStart, latestBlock)
	if sendResumed || ackResumed {
		ctx.Logger.Printf("[SubscribeEth] scanning from send=%d write_ack=%d (persisted send=%d write_ack=%d, head=%d); behind send=%d write_ack=%d",
			send, writeAck, cursors.EthSendBlock, cursors.EthWriteAckBlk, latestBlock, sendBehind, ackBehind)
	}
	return send, writeAck
}

func (s *Subscriber) persistEthCursors(ctx ethDeps, batchBuilder *services.BatchBuilder, send, writeAck uint64) {
	if s.recovery == nil || (send == 0 && writeAck == 0) {
		return
	}
	if batchBuilder != nil {
		floor := batchBuilder.LowestUnsubmittedEthHeight()
		send = clampRecoveryCursor(send, floor)
		writeAck = clampRecoveryCursor(writeAck, floor)
	}
	if err := s.recovery.SaveCheckpoint(recoverySourceID(ctx.IDs), services.RecoveryCursors{EthSendBlock: send, EthWriteAckBlk: writeAck}); err != nil {
		ctx.Logger.Printf("[SubscribeEth] failed to persist recovery cursors (send=%d write_ack=%d): %v", send, writeAck, err)
	}
}
