package services

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/big"
	"relayer/utils"
	"strconv"
	"strings"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	spectreContract "relayer/bindings/SpectreClient"
	client "relayer/client"
	"relayer/prover"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ics23 "github.com/cosmos/ics23/go"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// read abi json file once in runtime
var tendermintAbiJson *abi.ABI
var initErr error

func init() {
	tendermintAbiJson, initErr = spectreContract.ContractSpectreClientMetaData.GetAbi()
	if initErr != nil {
		log.Fatal(initErr)
	}
}

type TransactionHandler interface {
	CreateCosmosClientContract(ctx Context, clientState, consensusHash []byte, initialPinnedValidatorSet client.ContractValidatorSet) (ethcommon.Address, error)
	CreateEthClient(ctx Context, clientState ibcexported.ClientState, consensusState ibcexported.ConsensusState) (string, error)
	SendEthTx(ctx Context, msg any) error
	SendEthTxBatch(ctx Context, msgs []any) error
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

	CosmosPackets chan CosmosBatch
	EthPackets    chan EthBatch
	BatchBuilder  *BatchBuilder

	// Cross-chunk memoization (issue #76). Each field is touched by exactly one
	// handler goroutine, so no locking is needed even after the #4 split:
	//   - lastCosmosAppHashHeight: written/read only by handleCosmos.
	//   - lastFinalizedExecBlock:  written/read only by handleEth.
	// They let a chunked flush skip setup work already done for an earlier chunk
	// of the same source block.
	//
	// lastCosmosAppHashHeight: highest Cosmos height confirmed committed into the
	// AppHash. handleCosmos skips the AppHash wait when the batch's packets are
	// all at or below this height.
	//
	// lastFinalizedExecBlock: highest beacon-finalized ETH execution block
	// observed. waitBeaconFinality returns immediately when the requested event
	// block is at or below this value.
	lastCosmosAppHashHeight uint64
	lastFinalizedExecBlock  uint64
}

type cosmosClientFreshness struct {
	trustedHeight  int64
	trustedTime    time.Time
	trustingPeriod time.Duration
	clockDrift     time.Duration
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
		CosmosPackets: make(chan CosmosBatch),
		EthPackets:    make(chan EthBatch),
		BatchBuilder:  NewBatchBuilder(),
	}
}

func (s *Services) StartLoop(ctx Context) error {
	freshness, err := seedCosmosClientFreshness(ctx)
	if err != nil {
		return fmt.Errorf("seed cosmos client freshness: %w", err)
	}
	// Derived once at startup from the current on-chain client; restart the
	// relayer after migrating to a client with a different trusting period.
	routineInterval, err := deriveCosmosRefreshInterval(s.cosmosConfig, freshness.trustingPeriod)
	if err != nil {
		return err
	}
	log.Printf("[Routine] Cosmos client refresh interval=%s (trustedHeight=%d trustedTime=%s trustingPeriod=%s clockDrift=%s)",
		routineInterval, freshness.trustedHeight, freshness.trustedTime.UTC().Format(time.RFC3339),
		freshness.trustingPeriod, freshness.clockDrift)
	seedEthClientFreshness(ctx)

	go s.listener.SubscribeCosmos(ctx, s.BatchBuilder)
	go s.listener.SubscribeEth(ctx, s.BatchBuilder)

	// routinely run update client
	go func() {
		// TODO: Revisit routine scheduling strategy (interval/backoff/event-driven mix) to ensure this is optimal for production.
		// One configured interval gates both freshness routines (cosmos client
		// refresh below and the ETH client update further down).
		var refreshBackoff routineBackoff
		var ethClientBackoff routineBackoff
		for {
			now := time.Now()

			// periodically advance the cosmos light client, rotating the
			// pinned validator set when it is stale. The refresh reads the
			// authoritative on-chain trusted height; if the client is already
			// caught up it skips the transaction and returns the current block.
			ethUpdateTime, _ := ctx.latestEthTimestamp.Snapshot()
			if ethUpdateTime.Add(routineInterval).Before(now) && refreshBackoff.Ready(now) {
				latestBlock, err := s.worker.RefreshCosmosClient(
					ctx,
					s.cosmosConfig.ProofType,
					s.cosmosConfig.TrustLevel,
				)
				if err != nil {
					delay := refreshBackoff.RecordFailure(time.Now())
					log.Printf("[Routine] Failed to refresh cosmos client: %v; retrying after %s", err, delay)
				} else if err := recordRefreshResult(ctx.latestEthTimestamp, latestBlock); err != nil {
					delay := refreshBackoff.RecordFailure(time.Now())
					log.Printf("[Routine] Failed to record cosmos client refresh result: %v; retrying after %s", err, delay)
				} else {
					refreshBackoff.RecordSuccess()
				}
			}

			// update client on Cosmos side routinely
			cosmosUpdateTime, _ := ctx.latestCosmosTimestamp.Snapshot()
			if cosmosUpdateTime.Add(routineInterval).Before(now) && ethClientBackoff.Ready(now) {
				if err := s.worker.UpdateEthClient(ctx); err != nil {
					delay := ethClientBackoff.RecordFailure(time.Now())
					log.Printf("[Routine] Failed to update ETH client on Cosmos: %v; retrying after %s", err, delay)
				} else {
					ethClientBackoff.RecordSuccess()
				}
			}

			time.Sleep(time.Second)
		}
	}()

	// check for batch builder
	go func() {
		for {
			// check for batch every 3 seconds
			time.Sleep(time.Second * 3)
			s.BatchBuilder.CheckCosmos(ctx.Config.BatchConfig, s.CosmosPackets)
			s.BatchBuilder.CheckEth(ctx.Config.BatchConfig, s.EthPackets)
		}
	}()

	// scan for cosmos-originated packets that have timed out on ETH
	go func() {
		for {
			time.Sleep(time.Second * 30)
			s.scanForCosmosTimeouts(ctx)
		}
	}()

	// scan for ETH-originated packets that have timed out on Cosmos
	go func() {
		for {
			time.Sleep(time.Second * 30)
			s.scanForEthTimeouts(ctx)
		}
	}()

	// Handle the two directions on independent goroutines so a slow Cosmos→ETH
	// chunk (proof gen, beacon-finality wait) doesn't stall the next ETH→Cosmos
	// chunk and vice versa (issue #76 #4). The two handlers touch disjoint
	// mutable state: handleCosmos owns latestEthTimestamp + lastCosmosAppHashHeight,
	// handleEth owns latestCosmosTimestamp + lastFinalizedExecBlock. Shared
	// dependencies (TxHandler nonce paths, BatchBuilder.PendingTracker, the
	// Timestamp accessors) are each independently goroutine-safe.
	go func() {
		for batch := range s.CosmosPackets {
			s.handleCosmos(ctx, batch)
		}
		log.Println("[StartLoop] Cosmos batch channel closed, cosmos handler exiting")
	}()

	for batch := range s.EthPackets {
		s.handleEth(ctx, batch)
	}
	log.Println("[StartLoop] Eth batch channel closed, exiting loop")
	return nil
}

func seedCosmosClientFreshness(ctx Context) (cosmosClientFreshness, error) {
	clientState, err := fetchOnChainClientState(ctx)
	if err != nil {
		return cosmosClientFreshness{}, err
	}
	trustedHeight, err := clientStateRevisionHeightInt64(clientState)
	if err != nil {
		return cosmosClientFreshness{}, err
	}
	lightBlock, err := client.GetLightBlock(ctx.CosmosClient(), trustedHeight)
	if err != nil {
		staleTime := seedStaleTimestamp(ctx.latestEthTimestamp, uint64(trustedHeight))
		log.Printf("[Routine] Failed to fetch trusted cosmos light block %d: %v; seeded stale timestamp so refresh routine fires immediately",
			trustedHeight, err)
		return cosmosClientFreshness{
			trustedHeight:  trustedHeight,
			trustedTime:    staleTime,
			trustingPeriod: time.Duration(clientState.TrustingPeriod) * time.Second,
			clockDrift:     time.Duration(clientState.ClockDrift) * time.Second,
		}, nil
	}
	if err := recordRefreshResult(ctx.latestEthTimestamp, lightBlock); err != nil {
		return cosmosClientFreshness{}, err
	}
	log.Printf("[Routine] Seeded cosmos client freshness from trusted height %d at %s",
		trustedHeight, lightBlock.SignedHeader.Header.Time.UTC().Format(time.RFC3339))
	return cosmosClientFreshness{
		trustedHeight:  trustedHeight,
		trustedTime:    lightBlock.SignedHeader.Header.Time,
		trustingPeriod: time.Duration(clientState.TrustingPeriod) * time.Second,
		clockDrift:     time.Duration(clientState.ClockDrift) * time.Second,
	}, nil
}

func seedEthClientFreshness(ctx Context) {
	stale := func(reason string) {
		staleTime := seedStaleTimestamp(ctx.latestCosmosTimestamp, 0)
		log.Printf("[Routine] %s; seeded ETH client freshness as stale at %s",
			reason, staleTime.UTC().Format(time.RFC3339))
	}

	ethClientID := ctx.EthClientID()
	if ethClientID == "" {
		stale("Ethereum client ID is not configured")
		return
	}
	if ctx.CosmosClient() == nil {
		stale("Cosmos RPC client is not configured")
		return
	}
	ethClientState, err := client.GetEthereumClientState(ctx.CosmosClient(), ethClientID)
	if err != nil {
		stale(fmt.Sprintf("Failed to fetch Ethereum client state %q from Cosmos: %v", ethClientID, err))
		return
	}
	proofTimestamp := ethClientState.ComputeTimestampAtSlot(ethClientState.LatestSlot)
	result := &EthClientUpdateResult{
		EthClientState: ethClientState,
		ProofTimestamp: proofTimestamp,
	}
	if err := recordEthClientUpdateResult(ctx.latestCosmosTimestamp, result); err != nil {
		stale(fmt.Sprintf("Failed to record ETH client freshness from state %q: %v", ethClientID, err))
		return
	}
	trustedTime := time.Unix(int64(proofTimestamp), 0)
	log.Printf("[Routine] Seeded ETH client freshness from trusted slot %d at %s",
		ethClientState.LatestSlot, trustedTime.UTC().Format(time.RFC3339))
}

func seedStaleTimestamp(timestamp *Timestamp, height uint64) time.Time {
	staleTime := time.Unix(0, 0)
	timestamp.Set(staleTime, height)
	return staleTime
}

func recordRefreshResult(timestamp *Timestamp, lightBlock *client.LightBlock) error {
	if lightBlock == nil {
		return fmt.Errorf("no light block returned")
	}
	if lightBlock.BlockHeight < 0 {
		return fmt.Errorf("negative light block height: %d", lightBlock.BlockHeight)
	}
	if lightBlock.SignedHeader.Header == nil {
		return fmt.Errorf("missing light block header at height %d", lightBlock.BlockHeight)
	}
	trustedTime := lightBlock.SignedHeader.Header.Time
	if trustedTime.IsZero() {
		return fmt.Errorf("zero light block timestamp at height %d", lightBlock.BlockHeight)
	}
	timestamp.Set(trustedTime, uint64(lightBlock.BlockHeight))
	return nil
}

func recordEthClientUpdateResult(timestamp *Timestamp, result *EthClientUpdateResult) error {
	if result == nil {
		return fmt.Errorf("no ethereum client update result")
	}
	if result.EthClientState == nil {
		return fmt.Errorf("missing ethereum client state")
	}
	if result.ProofTimestamp == 0 {
		return fmt.Errorf("zero ethereum proof timestamp at slot %d", result.EthClientState.LatestSlot)
	}
	if result.ProofTimestamp > uint64(^uint64(0)>>1) {
		return fmt.Errorf("ethereum proof timestamp overflows int64: %d", result.ProofTimestamp)
	}
	proofTime := time.Unix(int64(result.ProofTimestamp), 0)
	timestamp.Set(proofTime, result.EthClientState.LatestSlot)
	return nil
}

func deriveCosmosRefreshInterval(cfg Config, trustingPeriod time.Duration) (time.Duration, error) {
	if trustingPeriod <= 0 {
		return 0, fmt.Errorf("on-chain trusting period must be positive: %s", trustingPeriod)
	}
	configured := cfg.RefreshInterval
	if configured == 0 {
		configured = DEFAULT_REFRESH_INTERVAL
	}
	margin := refreshSafetyMargin(trustingPeriod)
	maxInterval := trustingPeriod - margin
	if maxInterval <= 0 {
		return 0, fmt.Errorf("on-chain trusting period %s is too short for refresh safety margin %s", trustingPeriod, margin)
	}
	if cfg.RefreshIntervalConfigured && configured >= maxInterval {
		return 0, fmt.Errorf(
			"configured refresh interval %s must be less than on-chain trusting period %s minus safety margin %s",
			configured, trustingPeriod, margin,
		)
	}
	if configured >= maxInterval {
		log.Printf("[Routine] Default refresh interval %s exceeds safe interval %s; deriving from on-chain trusting period %s",
			configured, maxInterval, trustingPeriod)
		return maxInterval, nil
	}
	return configured, nil
}

func refreshSafetyMargin(trustingPeriod time.Duration) time.Duration {
	margin := DEFAULT_REFRESH_SAFETY_MARGIN
	maxMargin := trustingPeriod / 4
	if maxMargin < MIN_REFRESH_SAFETY_MARGIN {
		return maxMargin
	}
	if margin > maxMargin {
		margin = maxMargin
	}
	if margin < MIN_REFRESH_SAFETY_MARGIN {
		margin = MIN_REFRESH_SAFETY_MARGIN
	}
	return margin
}

func (s *Services) handleCosmos(ctx Context, batch CosmosBatch) {
	log.Printf("[StartLoop] Received cosmos batch: %d packets", len(batch.Packets))

	// A packet commitment written at block H is only reflected in the AppHash
	// queried at H+2. Wait until the chain has advanced that far instead of
	// sleeping a fixed 6s on every chunk: later chunks of the same source block
	// find the chain already advanced and return immediately (issue #76 #1).
	var maxPacketHeight uint64
	for _, p := range batch.Packets {
		if p.BlockNumber > maxPacketHeight {
			maxPacketHeight = p.BlockNumber
		}
	}
	if !s.waitCosmosAppHash(ctx, maxPacketHeight+2) {
		log.Printf("[StartLoop] cosmos AppHash wait failed; re-queueing %d packet(s)", len(batch.Packets))
		s.BatchBuilder.RequeueCosmosTransient(batch.Packets)
		return
	}

	// V2: build the Cosmos→ETH updateClient msg but do NOT submit it as its
	// own tx. It will be the first inner call of the multicall so we save
	// one round-trip + ~21k base intrinsic gas per flush.
	//
	// Pass the cached height as a hint only; BuildCosmosClientUpdateMsg
	// re-checks the authoritative on-chain trusted height before deciding to
	// regenerate the (expensive) Groth16 proof (issue #76 #2).
	_, ethTrustedHeight := ctx.latestEthTimestamp.Snapshot()
	// forceRotation=false: per-packet flushes rotate the pinned set only when
	// its overlap has decayed to the configured rotation threshold.
	updateBuild, err := s.worker.BuildCosmosClientUpdateMsg(
		ctx, s.cosmosConfig.ProofType,
		int64(ethTrustedHeight),
		s.cosmosConfig.TrustLevel,
		false,
	)
	if err != nil {
		// The whole chunk was already sliced off the queue by CheckCosmos; if we
		// bail here without re-queueing it is lost. Re-queue so the next flush
		// retries once the transient cause (proof gen, RPC) clears (issue #80).
		log.Printf("[StartLoop] Failed to build cosmos light client update: %v", err)
		s.BatchBuilder.RequeueCosmosTransient(batch.Packets)
		return
	}
	if updateBuild == nil || updateBuild.LightBlock == nil {
		log.Printf("[StartLoop] BuildCosmosClientUpdateMsg returned nil light block")
		s.BatchBuilder.RequeueCosmosTransient(batch.Packets)
		return
	}
	latestLightBlock := updateBuild.LightBlock

	hctx, hcancel := context.WithTimeout(context.Background(), ctx.Config.FetchTimeout)
	ethHeader, err := ctx.EthClient().HeaderByNumber(hctx, nil)
	hcancel()
	if err != nil {
		log.Printf("[StartLoop] Failed to get eth block header: %v", err)
		s.BatchBuilder.RequeueCosmosTransient(batch.Packets)
		return
	}
	if ethHeader == nil || ethHeader.Time == 0 {
		log.Printf("[StartLoop] Fetched eth block header is nil or has zero time")
		s.BatchBuilder.RequeueCosmosTransient(batch.Packets)
		return
	}
	ethBlockTime := ethHeader.Time

	// Build per-packet msgs into a single slice and submit one multicall when
	// N >= 2 (issue #67 / benchmark V1). The per-packet planning (including the
	// proof-build failure routing) is extracted into planCosmosPacketMsgs so the
	// failure branches are unit-tested without a live Context/RPC (issue #106).
	membershipFn := func(packet channeltypesv2.Packet, clientID string, pathType []byte) ([]byte, error) {
		return s.cosmosMembership(ctx, packet, clientID, pathType, latestLightBlock)
	}
	nonMembershipFn := func(packet channeltypesv2.Packet, clientID string, pathType []byte) ([]byte, error) {
		return s.cosmosNonMembership(ctx, packet, clientID, pathType, latestLightBlock)
	}
	planned, trackerAdds, transientFailures, trackerRemoves := planCosmosPacketMsgs(
		batch.Packets, ethBlockTime, ctx.CosmosRouterClientID(), membershipFn, nonMembershipFn)

	for _, p := range trackerAdds {
		s.BatchBuilder.PendingTracker.Add(*p.Packet, p.BlockNumber)
	}

	// +1 capacity for the optional client-update prepend.
	msgs := make([]cosmosBatchMsg, 0, len(planned)+1)
	if updateBuild.HasMsg {
		msgs = append(msgs, cosmosBatchMsg{
			msg:   *updateBuild,
			label: "UpdateClient",
		})
	}
	msgs = append(msgs, planned...)

	// Re-queue packets that failed proof construction this round (transient) so a
	// brief RPC/AppHash hiccup does not permanently drop a valid packet (issue #106).
	// CheckCosmos already sliced them off the queue, so without this they are lost.
	if len(transientFailures) > 0 {
		s.BatchBuilder.RequeueCosmosTransient(transientFailures)
	}
	for _, p := range trackerRemoves {
		s.BatchBuilder.PendingTracker.Remove(p.Packet.SourceClient, p.Packet.Sequence)
	}

	if len(msgs) == 0 {
		return
	}

	// V1 threshold: only use multicall when there are at least 2 msgs. Single-
	// msg flushes go through SendEthTx to avoid the tiny multicall overhead.
	rawMsgs := make([]any, len(msgs))
	for i, m := range msgs {
		rawMsgs[i] = m.msg
	}

	var sendErr error
	if len(msgs) >= 2 {
		log.Printf("[StartLoop] Submitting %d-packet multicall to ETH", len(msgs))
		sendErr = s.worker.TxHandler.SendEthTxBatch(ctx, rawMsgs)
	} else {
		sendErr = s.worker.TxHandler.SendEthTx(ctx, rawMsgs[0])
	}

	if sendErr != nil {
		failed := make([]CosmosPacket, 0, len(msgs))
		for _, m := range msgs {
			log.Printf("[%s] seq=%d: batch submission failed: %v", m.label, m.sequence, sendErr)
			if m.origin != nil {
				failed = append(failed, *m.origin)
			}
		}
		// The folded updateClient (origin == nil) is rebuilt fresh each flush, so
		// only the packet msgs need re-queueing (issue #80). Classify the failure:
		// an on-chain revert is deterministic (consumes retry budget → eventual
		// dead-letter); anything else (RPC, timeout) is transient and must not
		// burn the budget for a valid packet (issue #80 review).
		if errors.Is(sendErr, ErrPermanentRelayFailure) {
			s.BatchBuilder.RequeueCosmosPermanent(failed)
		} else {
			s.BatchBuilder.RequeueCosmosTransient(failed)
		}
		return
	}

	// Multicall is all-or-nothing: on success, every inner call applied. Remove
	// CosmosSend entries from the pending tracker and emit per-packet completion
	// logs to preserve the existing log shape consumers expect.
	for _, m := range msgs {
		if m.isRecv {
			s.BatchBuilder.PendingTracker.Remove(m.sourceClient, m.sequence)
		}
		if m.label == "UpdateClient" {
			log.Printf("[UpdateClient] relay completed (folded into batch)")
		} else {
			log.Printf("[%s] seq=%d: relay completed", m.label, m.sequence)
		}
	}

	// Refresh scheduling is based on the trusted Cosmos consensus timestamp,
	// which is what the on-chain trusting-period check uses.
	if err := recordRefreshResult(ctx.latestEthTimestamp, latestLightBlock); err != nil {
		log.Printf("[UpdateClient] failed to record trusted cosmos timestamp: %v", err)
	}
}

func (s *Services) handleEth(ctx Context, batch EthBatch) {
	log.Printf("[StartLoop] Received eth batch: %d packets", len(batch.Packets))

	// V2: build per-packet Cosmos msgs into a single slice, fold the wasm
	// MsgUpdateClient(s) at the head, and submit one SendCosmosTxBatch.
	// Cosmos tx is atomic ⇒ MsgUpdateClient applies first, packet msgs verify
	// against the freshly-advanced state in the same block.
	cosmosMsgs := make([]ethBatchMsg, 0, len(batch.Packets)+1)

	// Pre-filter expired EthSend packets — they bypass the batch and go to the
	// async timeout scanner. Anything else is kept for proof generation.
	var relayable []relayablePacket
	maxEventBlock := uint64(0)
	addRelayable := func(p EthPacket) {
		relayable = append(relayable, relayablePacket{packet: p})
		if p.BlockNumber > maxEventBlock {
			maxEventBlock = p.BlockNumber
		}
	}
	for _, p := range batch.Packets {
		switch p.Type {
		case EthSend:
			if ethPacketExpired(p) {
				if s.timeoutEthSend(ctx, p) {
					s.BatchBuilder.EthPendingTracker.Remove(p.Packet.SourceClient, p.Packet.Sequence)
				}
				continue
			}
			addRelayable(p)
		case EthWriteAck:
			if len(p.AckBytes) == 0 {
				log.Printf("[EthWriteAck] seq=%d: acknowledgement bytes missing, skipping", p.Packet.Sequence)
				continue
			}
			addRelayable(p)
		case EthAck:
			s.BatchBuilder.EthPendingTracker.Remove(p.Packet.SourceClient, p.Packet.Sequence)
			log.Printf("[EthAck] seq=%d: terminal event handled", p.Packet.Sequence)
		case EthTimeout:
			s.BatchBuilder.EthPendingTracker.Remove(p.Packet.SourceClient, p.Packet.Sequence)
			log.Printf("[EthTimeout] seq=%d: terminal event handled", p.Packet.Sequence)
		default:
			log.Printf("[StartLoop] Unknown eth packet type: %d (seq=%d)", p.Type, p.Packet.Sequence)
		}
	}
	if len(relayable) == 0 {
		return
	}

	// requeueRelayable puts the chunk's relayable packets back on the queue when
	// we bail before a successful submit. Every early return below would
	// otherwise drop the packets, which CheckEth already sliced off the queue
	// (issue #80) — the insufficient-sync-committee path is the common trigger.
	// Every early-return below is an infrastructure (transient) failure, so the
	// packets re-queue WITHOUT consuming the retry budget — a valid packet must
	// not be dropped just because finality/RPC was briefly unavailable.
	requeueRelayable := func() {
		pkts := make([]EthPacket, len(relayable))
		for i, r := range relayable {
			pkts[i] = r.packet
		}
		s.BatchBuilder.RequeueEthTransient(pkts)
	}

	// Wait once for beacon finality to cover the highest event block in the
	// batch. Subsequent packets in the same batch are by definition at
	// smaller-or-equal block numbers, so a single wait suffices.
	if !s.waitBeaconFinality(ctx, maxEventBlock, "EthBatch") {
		requeueRelayable()
		return
	}

	// Build wasm MsgUpdateClient(s) without submitting; the same EthClientState
	// output gives us the proof slot we'd see on-chain after the update applies.
	buildResult, err := s.worker.BuildEthClientUpdateMsgs(ctx)
	if err != nil {
		log.Printf("[EthBatch] BuildEthClientUpdateMsgs failed: %v", err)
		requeueRelayable()
		return
	}
	if buildResult == nil || buildResult.EthClientState == nil {
		log.Printf("[EthBatch] BuildEthClientUpdateMsgs returned nil result")
		requeueRelayable()
		return
	}

	// If we have update msgs, make sure the Cosmos chain has caught up enough
	// for the signature slot of the wasm update before we broadcast.
	if len(buildResult.Msgs) > 0 {
		s.worker.waitForCosmosCatchUp(ctx, buildResult.EthClientState, buildResult.SigSlot)
	}

	signerAddr, err := s.worker.TxHandler.CosmosSignerAddress()
	if err != nil {
		log.Printf("[EthBatch] failed to get cosmos signer: %v", err)
		requeueRelayable()
		return
	}

	proofBlockNumber := buildResult.EthClientState.LatestExecutionBlockNumber
	proofSlot := buildResult.EthClientState.LatestSlot
	if proofBlockNumber < maxEventBlock {
		log.Printf("[EthBatch] post-update proof block %d still < max event block %d, skipping",
			proofBlockNumber, maxEventBlock)
		requeueRelayable()
		return
	}

	// Prepend wasm MsgUpdateClient(s) — atomicity of the Cosmos tx applies them
	// before any packet msg verifies against the updated client state.
	for _, m := range buildResult.Msgs {
		cosmosMsgs = append(cosmosMsgs, ethBatchMsg{msg: m, label: "UpdateClient"})
	}

	// Plan the per-packet msgs; the proof builder is injected so the failure
	// branches are unit-tested without RPC (issue #106).
	buildProof := func(path []byte) ([]byte, error) {
		return client.GetEthMembershipProof(
			ctx.EthClient(), *ctx.RouterContract(), path,
			ethcommon.HexToHash(ICS26_IBC_STORAGE_SLOT), new(big.Int).SetUint64(proofBlockNumber))
	}
	planned, expired, failedRelayable := planEthPacketMsgs(relayable, proofSlot, signerAddr, buildProof)
	for _, p := range expired {
		if s.timeoutEthSend(ctx, p) {
			s.BatchBuilder.EthPendingTracker.Remove(p.Packet.SourceClient, p.Packet.Sequence)
		}
	}
	cosmosMsgs = append(cosmosMsgs, planned...)

	// Re-queue packets that failed proof construction this round (transient) so a
	// brief RPC/finality hiccup does not permanently drop a valid packet (issue #106).
	// CheckEth already sliced them off the queue, so without this they are lost.
	if len(failedRelayable) > 0 {
		s.BatchBuilder.RequeueEthTransient(failedRelayable)
	}

	if len(cosmosMsgs) == 0 {
		return
	}

	rawMsgs := make([]any, len(cosmosMsgs))
	for i, m := range cosmosMsgs {
		rawMsgs[i] = m.msg
	}

	log.Printf("[EthBatch] Submitting %d-msg batch to Cosmos", len(rawMsgs))
	if err := s.worker.TxHandler.SendCosmosTxBatch(ctx, rawMsgs); err != nil {
		var partialErr *BatchPartialError
		var failedMsgs []ethBatchMsg
		if errors.As(err, &partialErr) && partialErr.SucceededCount > 0 && partialErr.SucceededCount <= len(cosmosMsgs) {
			failedMsgs = cosmosMsgs[partialErr.SucceededCount:]
			for _, m := range cosmosMsgs[:partialErr.SucceededCount] {
				if m.label == "UpdateClient" {
					log.Printf("[UpdateClient] relay completed (succeeded sub-batch)")
				} else {
					log.Printf("[%s] seq=%d: relay completed (succeeded sub-batch)", m.label, m.sequence)
				}
			}
			err = partialErr.Err
		} else {
			failedMsgs = cosmosMsgs
		}

		failed := make([]EthPacket, 0, len(failedMsgs))
		for _, m := range failedMsgs {
			if m.label == "UpdateClient" {
				log.Printf("[UpdateClient] batch submission failed: %v", err)
			} else {
				log.Printf("[%s] seq=%d: batch submission failed: %v", m.label, m.sequence, err)
			}
			if m.origin != nil {
				failed = append(failed, *m.origin)
			}
		}
		// Folded wasm updateClient (origin == nil) is rebuilt each flush; only
		// packet msgs are re-queued (issue #80). Classify: a Cosmos DeliverTx
		// revert is deterministic (consumes budget → dead-letter); CheckTx /
		// broadcast / RPC errors are transient and must not burn the budget.
		if errors.Is(err, ErrPermanentRelayFailure) {
			timeoutPackets, permanentPackets := splitEthPermanentFailures(failedMsgs, err)
			for _, p := range timeoutPackets {
				log.Printf("[EthSend] seq=%d: Cosmos failure indicates timeout; routing to EthTimeout", p.Packet.Sequence)
				if s.timeoutEthSend(ctx, p) {
					s.BatchBuilder.EthPendingTracker.Remove(p.Packet.SourceClient, p.Packet.Sequence)
				}
			}
			s.BatchBuilder.RequeueEthPermanent(permanentPackets)
		} else {
			s.BatchBuilder.RequeueEthTransient(failed)
		}
		return
	}

	for _, m := range cosmosMsgs {
		if m.label == "UpdateClient" {
			log.Printf("[UpdateClient] relay completed (folded into Cosmos batch)")
		} else {
			log.Printf("[%s] seq=%d: relay completed", m.label, m.sequence)
		}
	}

	// Refresh scheduling uses the ETH timestamp trusted by the 08-wasm client,
	// not the relayer wall clock. If no update msg was needed, this records the
	// current on-chain client state that BuildEthClientUpdateMsgs already read.
	if err := recordEthClientUpdateResult(ctx.latestCosmosTimestamp, buildResult); err != nil {
		log.Printf("[UpdateClient] failed to record trusted ETH timestamp: %v", err)
	}
}

func splitEthPermanentFailures(failedMsgs []ethBatchMsg, err error) (timeoutPackets, permanentPackets []EthPacket) {
	for _, m := range failedMsgs {
		if m.origin == nil {
			continue
		}
		packet := *m.origin
		if packet.Type == EthSend && shouldTimeoutEthSend(packet, err) {
			timeoutPackets = append(timeoutPackets, packet)
			continue
		}
		permanentPackets = append(permanentPackets, packet)
	}
	return timeoutPackets, permanentPackets
}

// waitCosmosAppHash blocks until the Cosmos chain height reaches targetHeight,
// at which point a packet commitment written ≤ targetHeight-2 is guaranteed to
// be reflected in the queried AppHash. It memoizes the highest confirmed height
// so later chunks of the same source block skip the RPC entirely (issue #76 #1).
//
// On timeout it returns false so the caller can re-queue the chunk instead of
// generating proofs against stale AppHash state.
func (s *Services) waitCosmosAppHash(ctx Context, targetHeight uint64) bool {
	if targetHeight <= s.lastCosmosAppHashHeight {
		return true
	}
	maxRetries := ctx.Config.AppHashWaitRetries
	if maxRetries == 0 {
		maxRetries = DEFAULT_COSMOS_APP_HASH_WAIT_RETRIES
	}
	interval := ctx.Config.AppHashWaitInterval
	if interval == 0 {
		interval = DEFAULT_COSMOS_APP_HASH_WAIT_INTERVAL
	}
	for attempt := uint32(0); attempt < maxRetries; attempt++ {
		status, err := ctx.CosmosClient().Status(context.Background())
		if err != nil {
			log.Printf("[StartLoop] failed to query cosmos status: %v", err)
			time.Sleep(interval)
			continue
		}
		if status.SyncInfo.LatestBlockHeight < 0 {
			log.Printf("[StartLoop] cosmos status returned negative latest block height: %d", status.SyncInfo.LatestBlockHeight)
			return false
		}
		current := uint64(status.SyncInfo.LatestBlockHeight)
		if current >= targetHeight {
			s.lastCosmosAppHashHeight = current
			return true
		}
		if attempt == 0 {
			log.Printf("[StartLoop] waiting for cosmos AppHash to cover height %d (current %d)...",
				targetHeight, current)
		}
		time.Sleep(interval)
	}
	log.Printf("[StartLoop] cosmos AppHash wait for height %d timed out after %s", targetHeight, time.Duration(maxRetries)*interval)
	return false
}

// waitBeaconFinality polls beacon finality until execution block ≥ target.
// handleEth uses it once per batch so every relayable packet in the batch can
// share the same finality wait.
//
// It memoizes the highest finalized execution block seen, so later chunks whose
// event block is already covered return immediately without an RPC (issue #76 #3).
func (s *Services) waitBeaconFinality(ctx Context, eventBlock uint64, tag string) bool {
	if eventBlock <= s.lastFinalizedExecBlock {
		log.Printf("[%s] beacon finality already covers block %d (finalized %d)",
			tag, eventBlock, s.lastFinalizedExecBlock)
		return true
	}
	maxRetries := ctx.Config.BeaconFinalityRetries
	if maxRetries == 0 {
		maxRetries = DEFAULT_BEACON_FINALITY_RETRIES
	}
	log.Printf("[%s] waiting for beacon finality at block %d", tag, eventBlock)
	for attempt := uint32(0); attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(10 * time.Second)
		}
		bctx, bcancel := context.WithTimeout(context.Background(), 15*time.Second)
		finalityUpdate, err := client.GetFinalityUpdate(bctx, ctx.BeaconAPIURL())
		bcancel()
		if err != nil {
			log.Printf("[%s] failed to get finality update: %v", tag, err)
			continue
		}
		execBlock, err := strconv.ParseUint(finalityUpdate.FinalizedHeader.Execution.BlockNumber, 10, 64)
		if err != nil {
			log.Printf("[%s] failed to parse finalized execution block %q: %v",
				tag, finalityUpdate.FinalizedHeader.Execution.BlockNumber, err)
			continue
		}
		if execBlock > s.lastFinalizedExecBlock {
			s.lastFinalizedExecBlock = execBlock
		}
		if execBlock >= eventBlock {
			log.Printf("[%s] beacon finalized block %d >= event block %d", tag, execBlock, eventBlock)
			return true
		}
		log.Printf("[%s] beacon finalized block %d < event block %d, waiting... (%d/%d)",
			tag, execBlock, eventBlock, attempt+1, maxRetries)
	}
	log.Printf("[%s] beacon finality did not reach block %d after %d retries", tag, eventBlock, maxRetries)
	return false
}

func ethPacketExpired(packet EthPacket) bool {
	return packet.Packet.TimeoutTimestamp > 0 && uint64(time.Now().Unix()) >= packet.Packet.TimeoutTimestamp
}

// ShouldRelayCosmosTimeoutToEth decides whether a CosmosTimeout event needs a
// follow-up ICS26Router.timeoutPacket on the ETH side.
//
// IBC v2 timeout is submitted to the *source* chain: it deletes the source-side
// packet commitment and refunds the sender. For a Cosmos→ETH packet, that
// MsgTimeout fires on Cosmos and the destination (ETH) never held a commitment
// for the packet — calling ETH's timeoutPacket would be a NoOp at best and a
// wasted Groth16 proof + tx in any case.
//
// The Cosmos-originated case is identified by `destination_client` matching the
// router client ID we know ETH uses for this Cosmos chain (`ICS26ClientID` in
// the config). When the destination is *not* our Cosmos chain's ETH-side ID,
// the packet came from ETH and ETH still owns a commitment that needs the
// timeoutPacket call to delete + refund.
//
// NOTE: this predicate previously compared `source_client` against the same
// router ID, which never matched for Cosmos-originated packets (their
// source_client is the Cosmos-side client name e.g. "08-wasm-0"). That fired
// the no-op ETH relay for every Cosmos→ETH timeout — wasted gas at best, and
// blocked the relay loop entirely whenever the ETH-side updateClient was
// failing for unrelated reasons.
func ShouldRelayCosmosTimeoutToEth(packet *channeltypesv2.Packet, cosmosRouterClientID string) bool {
	if packet == nil || cosmosRouterClientID == "" {
		return false
	}
	return packet.DestinationClient != cosmosRouterClientID
}

func cosmosPacketExpiredOnEth(packet CosmosPacket, ethBlockTime uint64) bool {
	return ethBlockTime > 0 && packet.Packet.TimeoutTimestamp > 0 && ethBlockTime >= packet.Packet.TimeoutTimestamp
}

func pendingPacketsTimedOutAtTimestamp(pending []pendingPacketInfo, timestamp uint64) []pendingPacketInfo {
	if timestamp == 0 {
		return nil
	}

	expired := make([]pendingPacketInfo, 0, len(pending))
	for _, info := range pending {
		if info.Packet.TimeoutTimestamp > 0 && timestamp >= info.Packet.TimeoutTimestamp {
			expired = append(expired, info)
		}
	}
	return expired
}

func shouldTimeoutEthSend(packet EthPacket, err error) bool {
	if !ethPacketExpired(packet) {
		return false
	}
	if err == nil {
		return false
	}

	var cosmosErr *CosmosTxFailure
	if errors.As(err, &cosmosErr) && cosmosTxFailureIndicatesTimeout(cosmosErr) {
		return true
	}
	if _, ok := timeoutEVMErrorName(err); ok {
		return true
	}
	return timeoutErrorMessage(err.Error())
}

func cosmosTxFailureIndicatesTimeout(err *CosmosTxFailure) bool {
	if err == nil {
		return false
	}
	return timeoutErrorMessage(err.Log) ||
		timeoutErrorMessage(string(err.Data)) ||
		timeoutEVMErrorDataMatches(err.Data)
}

func timeoutErrorMessage(msg string) bool {
	msg = strings.ToLower(msg)
	return strings.Contains(msg, "timeout elapsed") ||
		strings.Contains(msg, "ibcinvalidtimeouttimestamp") ||
		strings.Contains(msg, "timed out")
}

var timeoutErrorSelectors = map[[4]byte]string{
	evmErrorSelector("IBCInvalidTimeoutTimestamp(uint256,uint256)"): "IBCInvalidTimeoutTimestamp",
}

func evmErrorSelector(signature string) [4]byte {
	hash := crypto.Keccak256([]byte(signature))
	var selector [4]byte
	copy(selector[:], hash[:4])
	return selector
}

func timeoutEVMErrorName(callErr error) (string, bool) {
	type dataErr interface {
		ErrorData() interface{}
	}
	var de dataErr
	if !errors.As(callErr, &de) {
		return "", false
	}

	data := evmErrorDataBytes(de.ErrorData())
	if len(data) < 4 {
		return "", false
	}
	var selector [4]byte
	copy(selector[:], data[:4])
	name, ok := timeoutErrorSelectors[selector]
	return name, ok
}

func timeoutEVMErrorDataMatches(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	var selector [4]byte
	copy(selector[:], data[:4])
	_, ok := timeoutErrorSelectors[selector]
	return ok
}

func evmErrorDataBytes(raw interface{}) []byte {
	switch v := raw.(type) {
	case string:
		return ethcommon.FromHex(v)
	case fmt.Stringer:
		return ethcommon.FromHex(v.String())
	case []byte:
		return v
	default:
		return nil
	}
}

func (s *Services) updateCosmosClientForEth(ctx Context, tag string) (*client.LightBlock, bool) {
	_, ethTrustedHeight := ctx.latestEthTimestamp.Snapshot()
	latestLightBlock, err := s.worker.UpdateCosmosClient(ctx, s.cosmosConfig.ProofType, int64(ethTrustedHeight), s.cosmosConfig.TrustLevel, false)
	if err != nil {
		log.Printf("[%s] Failed to update cosmos light client: %v", tag, err)
		return nil, false
	}
	if latestLightBlock == nil {
		log.Printf("[%s] Failed to update cosmos light client: latestLightBlock is nil", tag)
		return nil, false
	}

	if err := recordRefreshResult(ctx.latestEthTimestamp, latestLightBlock); err != nil {
		log.Printf("[%s] failed to record trusted cosmos timestamp: %v", tag, err)
	}

	return latestLightBlock, true
}

func (s *Services) timeoutEthSend(ctx Context, packet EthPacket) bool {
	log.Printf("[EthTimeout] seq=%d: packet expired, preparing timeout proof", packet.Packet.Sequence)

	latestLightBlock, ok := s.updateCosmosClientForEth(ctx, "EthTimeout")
	if !ok {
		return false
	}

	counterpartyTime := uint64(latestLightBlock.SignedHeader.Header.Time.Unix())
	if counterpartyTime < packet.Packet.TimeoutTimestamp {
		log.Printf("[EthTimeout] seq=%d: counterparty time %d < timeout %d, skipping",
			packet.Packet.Sequence, counterpartyTime, packet.Packet.TimeoutTimestamp)
		return false
	}

	calldata, err := s.cosmosNonMembership(ctx, *packet.Packet, packet.Packet.DestinationClient, []byte{2}, latestLightBlock)
	if err != nil {
		log.Printf("[EthTimeout] seq=%d: %v", packet.Packet.Sequence, err)
		return false
	}

	msgTimeoutPacket := contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket{
		Packet:           toEthPacket(*packet.Packet),
		NonMembershipMsg: calldata,
	}

	if err := s.worker.TxHandler.SendEthTx(ctx, msgTimeoutPacket); err != nil {
		log.Printf("[EthTimeout] seq=%d: SendEthTx failed: %v", packet.Packet.Sequence, err)
		return false
	}
	log.Printf("[EthTimeout] seq=%d: relay completed", packet.Packet.Sequence)
	return true
}

const pendingTrackerMaxAge = 1 * time.Hour

func (s *Services) scanForEthTimeouts(ctx Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[EthTimeoutScan] Panic recovered: %v", r)
		}
	}()

	s.BatchBuilder.EthPendingTracker.PurgeStaleWithoutTimeout(pendingTrackerMaxAge)

	pending := s.BatchBuilder.EthPendingTracker.GetAll()
	if len(pending) == 0 {
		return
	}

	now := uint64(time.Now().Unix())
	expired := pendingPacketsTimedOutAtTimestamp(pending, now)
	if len(expired) == 0 {
		return
	}

	log.Printf("[EthTimeoutScan] Found %d locally-expired ETH-origin packet(s) at time %d", len(expired), now)
	for _, info := range expired {
		pendingCommitment, err := HasPendingEthPacketCommitment(ctx, info.Packet)
		if err != nil {
			log.Printf("[EthTimeoutScan] seq=%d: failed to check ETH packet commitment: %v", info.Packet.Sequence, err)
			continue
		}
		if !pendingCommitment {
			s.BatchBuilder.EthPendingTracker.Remove(info.Packet.SourceClient, info.Packet.Sequence)
			log.Printf("[EthTimeoutScan] seq=%d: ETH commitment already cleared, removed from pending tracker", info.Packet.Sequence)
			continue
		}

		packet := info.Packet
		if s.timeoutEthSend(ctx, EthPacket{
			Type:        EthSend,
			Packet:      &packet,
			BlockNumber: info.BlockNumber,
		}) {
			s.BatchBuilder.EthPendingTracker.Remove(info.Packet.SourceClient, info.Packet.Sequence)
		}
	}
}

func (s *Services) scanForCosmosTimeouts(ctx Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[CosmosTimeoutScan] Panic recovered: %v", r)
		}
	}()

	s.BatchBuilder.PendingTracker.PurgeStaleWithoutTimeout(pendingTrackerMaxAge)

	pending := s.BatchBuilder.PendingTracker.GetAll()
	if len(pending) == 0 {
		return
	}

	ethHeader, err := ctx.EthClient().HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Printf("[CosmosTimeoutScan] Failed to get eth block header: %v", err)
		return
	}
	ethBlockTime := ethHeader.Time

	log.Printf("[CosmosTimeoutScan] Checking %d pending packets against eth block time %d",
		len(pending), ethBlockTime)

	expired := pendingPacketsTimedOutAtTimestamp(pending, ethBlockTime)
	if len(expired) == 0 {
		return
	}

	log.Printf("[CosmosTimeoutScan] Found %d head-expired packets, building proof state", len(expired))

	updateResult, err := s.worker.BuildEthClientUpdateMsgs(ctx)
	if err != nil {
		log.Printf("[CosmosTimeoutScan] Failed to build ETH client update messages: %v", err)
		return
	}

	ethClientState := updateResult.EthClientState
	if ethClientState == nil {
		ethClientState, err = client.GetEthereumClientState(ctx.CosmosClient(), ctx.EthClientID())
		if err != nil {
			log.Printf("[CosmosTimeoutScan] Failed to get ETH client state: %v", err)
			return
		}
	}

	proofTimestamp := updateResult.ProofTimestamp
	expired = pendingPacketsTimedOutAtTimestamp(expired, proofTimestamp)
	if len(expired) == 0 {
		log.Printf("[CosmosTimeoutScan] Proof state timestamp %d has not reached any candidate timeout yet; retrying later", proofTimestamp)
		return
	}

	log.Printf("[CosmosTimeoutScan] Found %d proof-expired packets at proof timestamp %d (proof slot=%d exec_block=%d)",
		len(expired), proofTimestamp, ethClientState.LatestSlot, ethClientState.LatestExecutionBlockNumber)

	var timeoutMsgs []any
	var processed []pendingPacketInfo
	for _, info := range expired {
		msgTimeout, err := s.buildCosmosTimeoutMsg(ctx, info.Packet, ethClientState)
		if err != nil {
			log.Printf("[CosmosTimeout] seq=%d: %v", info.Packet.Sequence, err)
			continue
		}
		timeoutMsgs = append(timeoutMsgs, msgTimeout)
		processed = append(processed, info)
	}

	if len(timeoutMsgs) == 0 {
		return
	}

	var batchMsgs []any
	if len(updateResult.Msgs) > 0 {
		batchMsgs = append(batchMsgs, updateResult.Msgs...)
	}
	batchMsgs = append(batchMsgs, timeoutMsgs...)

	if len(updateResult.Msgs) > 0 {
		s.worker.waitForCosmosCatchUp(ctx, updateResult.EthClientState, updateResult.SigSlot)
	}

	if err := s.worker.TxHandler.SendCosmosTxBatch(ctx, batchMsgs); err != nil {
		log.Printf("[CosmosTimeoutScan] SendCosmosTxBatch failed: %v", err)
		var partialErr *BatchPartialError
		if errors.As(err, &partialErr) && partialErr.SucceededCount >= len(updateResult.Msgs) {
			if len(updateResult.Msgs) > 0 {
				if err := recordEthClientUpdateResult(ctx.latestCosmosTimestamp, updateResult); err != nil {
					log.Printf("[CosmosTimeoutScan] failed to record trusted ETH timestamp after partial batch: %v", err)
				}
			}
			succeededTimeoutsCount := partialErr.SucceededCount - len(updateResult.Msgs)
			for i := 0; i < succeededTimeoutsCount; i++ {
				info := processed[i]
				s.BatchBuilder.PendingTracker.Remove(info.Packet.SourceClient, info.Packet.Sequence)
				log.Printf("[CosmosTimeout] seq=%d: timeout relay completed (bundled with %d update msgs) in partial batch", info.Packet.Sequence, len(updateResult.Msgs))
			}
		}
		return
	}

	for _, info := range processed {
		s.BatchBuilder.PendingTracker.Remove(info.Packet.SourceClient, info.Packet.Sequence)
		log.Printf("[CosmosTimeout] seq=%d: timeout relay completed (bundled with %d update msgs)", info.Packet.Sequence, len(updateResult.Msgs))
	}
	if err := recordEthClientUpdateResult(ctx.latestCosmosTimestamp, updateResult); err != nil {
		log.Printf("[CosmosTimeoutScan] failed to record trusted ETH timestamp: %v", err)
	}
}

func (s *Services) buildCosmosTimeoutMsg(ctx Context, packet channeltypesv2.Packet, ethClientState *client.EthereumClientState) (*channeltypesv2.MsgTimeout, error) {
	receiptPath := EthPath(packet.DestinationClient, packet.Sequence, 2)
	proofBytes, err := client.GetEthNonMembershipProof(
		ctx.EthClient(), *ctx.RouterContract(), receiptPath, ethcommon.HexToHash(ICS26_IBC_STORAGE_SLOT), new(big.Int).SetUint64(ethClientState.LatestExecutionBlockNumber))
	if err != nil {
		return nil, fmt.Errorf("failed to get ETH non-membership proof: %w", err)
	}

	return channeltypesv2.NewMsgTimeout(
		packet,
		proofBytes,
		clienttypes.Height{RevisionNumber: 0, RevisionHeight: ethClientState.LatestSlot},
		"",
	), nil
}

func (s *Services) cosmosMembership(ctx Context, packet channeltypesv2.Packet, clientID string, pathType []byte, latestLightBlock *client.LightBlock) ([]byte, error) {
	height := latestLightBlock.BlockHeight
	ibcPath := utils.IbcPath(clientID, packet.Sequence, pathType)
	value, proof, err := client.ProvePath(ctx.CosmosClient(), height, ibcPath)
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return nil, fmt.Errorf("commitment empty at height=%d", height)
	}

	merkleProof, err := parseMerkleProof(proof.Proofs, packet.Sequence)
	if err != nil {
		return nil, err
	}

	membershipMsg := spectreContract.ILightClientMsgsMsgVerifyMembership{
		Height: spectreContract.IICS02ClientMsgsHeight{
			RevisionHeight: uint64(height),
			RevisionNumber: 0,
		},
		KvPairs: []spectreContract.IMembershipMsgsKVPair{
			{
				Path:  ibcPath,
				Value: value,
			},
		},
		MerkleProofs: []spectreContract.IMembershipMsgsMerkleProof{merkleProof},
		AppHash:      utils.BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
		TrustedConsensusState: spectreContract.IICS07TendermintMsgsConsensusState{
			Timestamp:          big.NewInt(latestLightBlock.SignedHeader.Header.Time.UnixNano()),
			Root:               utils.BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
			NextValidatorsHash: utils.BytesToBytes32(latestLightBlock.SignedHeader.Header.NextValidatorsHash),
		},
		MembershipType: 0,
	}

	calldata, err := tendermintAbiJson.Pack("verifyMembership", membershipMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to ABI encode verifyMembership: %w", err)
	}
	return calldata[4:], nil
}

func (s *Services) cosmosNonMembership(ctx Context, packet channeltypesv2.Packet, clientID string, pathType []byte, latestLightBlock *client.LightBlock) ([]byte, error) {
	height := latestLightBlock.BlockHeight
	ibcPath := utils.IbcPath(clientID, packet.Sequence, pathType)
	value, proof, err := client.ProvePath(ctx.CosmosClient(), height, ibcPath)
	if err != nil {
		return nil, err
	}
	if len(value) != 0 {
		return nil, fmt.Errorf("non-membership expected empty value at height=%d, got %d bytes", height, len(value))
	}

	merkleProof, err := parseMerkleProof(proof.Proofs, packet.Sequence)
	if err != nil {
		return nil, err
	}

	nonMembershipMsg := spectreContract.ILightClientMsgsMsgVerifyNonMembership{
		Height: spectreContract.IICS02ClientMsgsHeight{
			RevisionHeight: uint64(height),
			RevisionNumber: 0,
		},
		KvPairs: []spectreContract.IMembershipMsgsKVPair{
			{
				Path:  ibcPath,
				Value: value,
			},
		},
		MerkleProofs: []spectreContract.IMembershipMsgsMerkleProof{merkleProof},
		AppHash:      utils.BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
		TrustedConsensusState: spectreContract.IICS07TendermintMsgsConsensusState{
			Timestamp:          big.NewInt(latestLightBlock.SignedHeader.Header.Time.UnixNano()),
			Root:               utils.BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
			NextValidatorsHash: utils.BytesToBytes32(latestLightBlock.SignedHeader.Header.NextValidatorsHash),
		},
		MembershipType: 0,
	}

	calldata, err := tendermintAbiJson.Pack("verifyNonMembership", nonMembershipMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to ABI encode verifyNonMembership: %w", err)
	}
	return calldata[4:], nil
}

func parseMerkleProof(proofs []*ics23.CommitmentProof, sequence uint64) (spectreContract.IMembershipMsgsMerkleProof, error) {
	merkleProof := spectreContract.IMembershipMsgsMerkleProof{
		Proofs: []spectreContract.IMembershipMsgsCommitmentProof{},
	}
	for _, p := range proofs {
		commitmentProof, err := client.ParseCommitmentProof(p)
		if err != nil {
			return merkleProof, fmt.Errorf("failed to parse commitment proof for seq=%d: %w", sequence, err)
		}
		merkleProof.Proofs = append(merkleProof.Proofs, *commitmentProof)
	}
	return merkleProof, nil
}

func toEthPacket(packet channeltypesv2.Packet) contractICS26Router.IICS26RouterMsgsPacket {
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

	return contractICS26Router.IICS26RouterMsgsPacket{
		Sequence:         packet.Sequence,
		SourceClient:     packet.SourceClient,
		DestClient:       packet.DestinationClient,
		TimeoutTimestamp: packet.TimeoutTimestamp,
		Payloads:         payloads,
	}
}

func EthPath(clientID string, sequence uint64, pathType byte) []byte {
	seqBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(seqBytes, sequence)
	path := append([]byte(clientID), pathType)
	return append(path, seqBytes...)
}

func EthIBCStorageKey(path []byte) ethcommon.Hash {
	pathHash := crypto.Keccak256(path)
	return crypto.Keccak256Hash(pathHash, ethcommon.HexToHash(ICS26_IBC_STORAGE_SLOT).Bytes())
}

func HasEthIBCPathValue(ctx Context, path []byte) (bool, error) {
	if ctx.EthClient() == nil {
		return false, fmt.Errorf("eth client is nil")
	}
	if ctx.RouterContract() == nil {
		return false, fmt.Errorf("router contract address is nil")
	}

	value, err := ctx.EthClient().StorageAt(context.Background(), *ctx.RouterContract(), EthIBCStorageKey(path), nil)
	if err != nil {
		return false, fmt.Errorf("eth storage query failed: %w", err)
	}

	for _, b := range value {
		if b != 0 {
			return true, nil
		}
	}
	return false, nil
}

func HasEthPacketReceipt(ctx Context, packet channeltypesv2.Packet) (bool, error) {
	return HasEthIBCPathValue(ctx, EthPath(packet.DestinationClient, packet.Sequence, 2))
}

func HasPendingEthPacketCommitment(ctx Context, packet channeltypesv2.Packet) (bool, error) {
	return HasEthIBCPathValue(ctx, EthPath(packet.SourceClient, packet.Sequence, 1))
}
