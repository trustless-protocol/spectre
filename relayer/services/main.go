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

	tendermintContract "relayer/bindings/Groth16ICS07Tendermint"
	contractICS26Router "relayer/bindings/ICS26Router"
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
	tendermintAbiJson, initErr = tendermintContract.ContractGroth16ICS07TendermintMetaData.GetAbi()
	if initErr != nil {
		log.Fatal(initErr)
	}
}

type TransactionHandler interface {
	CreateCosmosClientContract(ctx Context, clientState, consensusHash []byte, initialPinnedValidatorSet client.ContractValidatorSet) (ethcommon.Address, error)
	CreateEthClient(ctx Context, clientState ibcexported.ClientState, consensusState ibcexported.ConsensusState) (string, error)
	SendEthTx(ctx Context, msg any) error
	SendEthTxBatch(ctx Context, msgs []any) error
	SendReAnchorPinnedSet(ctx Context, updateMsg any, newPinnedValidatorSet any) error
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

func (s *Services) StartLoop(ctx Context) {
	go s.listener.SubscribeCosmos(ctx, s.BatchBuilder)
	go s.listener.SubscribeEth(ctx, s.BatchBuilder)

	// routinely run update client
	go func() {
		// TODO: Revisit routine scheduling strategy (interval/backoff/event-driven mix) to ensure this is optimal for production.
		routineInterval := 24 * time.Hour
		for {
			now := time.Now()

			// periodically advance the cosmos light client and re-anchor the
			// pinned validator set in one shot. The re-anchor reads the
			// authoritative on-chain trusted height; if the client is already
			// caught up it skips the transaction and returns the current block.
			ethUpdateTime, _ := ctx.latestEthTimestamp.Snapshot()
			if ethUpdateTime.Add(routineInterval).Before(now) {
				latestBlock, err := s.worker.ReAnchorCosmosPinnedSet(
					ctx,
					s.cosmosConfig.ProofType,
					s.cosmosConfig.TrustLevel,
				)
				if err != nil {
					log.Printf("[Routine] Failed to re-anchor cosmos pinned set: %v", err)
					time.Sleep(time.Second)
					continue
				}
				if err := recordReAnchorResult(ctx.latestEthTimestamp, latestBlock, time.Now()); err != nil {
					log.Printf("[Routine] Failed to re-anchor cosmos pinned set: %v", err)
					time.Sleep(time.Second)
					continue
				}
			}

			// update client on Cosmos side routinely
			cosmosUpdateTime, _ := ctx.latestCosmosTimestamp.Snapshot()
			if cosmosUpdateTime.Add(routineInterval).Before(now) {
				s.worker.UpdateEthClient(ctx)
				ctx.latestCosmosTimestamp.SetTime(now)
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
}

func recordReAnchorResult(timestamp *Timestamp, lightBlock *client.LightBlock, now time.Time) error {
	if lightBlock == nil {
		return fmt.Errorf("no light block returned")
	}
	if lightBlock.BlockHeight < 0 {
		return fmt.Errorf("negative light block height: %d", lightBlock.BlockHeight)
	}
	timestamp.Set(now, uint64(lightBlock.BlockHeight))
	return nil
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
	updateBuild, err := s.worker.BuildCosmosClientUpdateMsg(
		ctx, s.cosmosConfig.ProofType,
		int64(ethTrustedHeight),
		s.cosmosConfig.TrustLevel,
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

	// +1 capacity for the optional updateClient prepend.
	msgs := make([]cosmosBatchMsg, 0, len(planned)+1)
	if updateBuild.HasMsg {
		msgs = append(msgs, cosmosBatchMsg{
			msg:   updateBuild.Msg,
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

	// Advance the trusted ETH-side height only after the batch (including the
	// folded updateClient, if any) confirmed on-chain.
	if updateBuild.HasMsg {
		ctx.latestEthTimestamp.Set(time.Now(), uint64(latestLightBlock.BlockHeight))
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
	for _, p := range batch.Packets {
		switch p.Type {
		case EthSend:
			if ethPacketExpired(p) {
				if s.timeoutEthSend(ctx, p) {
					s.BatchBuilder.EthPendingTracker.Remove(p.Packet.SourceClient, p.Packet.Sequence)
				}
				continue
			}
			relayable = append(relayable, relayablePacket{packet: p})
		case EthWriteAck:
			if len(p.AckBytes) == 0 {
				log.Printf("[EthWriteAck] seq=%d: acknowledgement bytes missing, skipping", p.Packet.Sequence)
				continue
			}
			relayable = append(relayable, relayablePacket{packet: p})
		case EthAck:
			s.BatchBuilder.EthPendingTracker.Remove(p.Packet.SourceClient, p.Packet.Sequence)
			log.Printf("[EthAck] seq=%d: terminal event handled", p.Packet.Sequence)
		case EthTimeout:
			s.BatchBuilder.EthPendingTracker.Remove(p.Packet.SourceClient, p.Packet.Sequence)
			log.Printf("[EthTimeout] seq=%d: terminal event handled", p.Packet.Sequence)
		default:
			log.Printf("[StartLoop] Unknown eth packet type: %d (seq=%d)", p.Type, p.Packet.Sequence)
		}
		if p.BlockNumber > maxEventBlock {
			maxEventBlock = p.BlockNumber
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
			s.BatchBuilder.RequeueEthPermanent(failed)
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

	// Advance latestCosmosTimestamp only after the batch (including the folded
	// wasm updateClient, if any) confirmed on-chain.
	if len(buildResult.Msgs) > 0 {
		ctx.latestCosmosTimestamp.SetTime(time.Now())
	}
}

// waitCosmosAppHash blocks until the Cosmos chain height reaches targetHeight,
// at which point a packet commitment written ≤ targetHeight-2 is guaranteed to
// be reflected in the queried AppHash. It memoizes the highest confirmed height
// so later chunks of the same source block skip the RPC entirely (issue #76 #1).
//
// Timeout: 30 polls x 1s = 30s. On timeout it returns false so the caller can
// re-queue the chunk instead of generating proofs against stale AppHash state.
func (s *Services) waitCosmosAppHash(ctx Context, targetHeight uint64) bool {
	if targetHeight <= s.lastCosmosAppHashHeight {
		return true
	}
	for attempt := 0; attempt < 30; attempt++ {
		status, err := ctx.CosmosClient().Status(context.Background())
		if err != nil {
			log.Printf("[StartLoop] failed to query cosmos status: %v", err)
			time.Sleep(1 * time.Second)
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
		time.Sleep(1 * time.Second)
	}
	log.Printf("[StartLoop] cosmos AppHash wait for height %d timed out after 30s", targetHeight)
	return false
}

// waitBeaconFinality polls beacon finality until execution block ≥ target.
// Extracted from the per-packet ethProofHeight so handleEth can amortize the
// wait over a whole batch. Returns false on timeout (60 polls × 10s = 10min).
//
// It memoizes the highest finalized execution block seen, so later chunks whose
// event block is already covered return immediately without an RPC (issue #76 #3).
func (s *Services) waitBeaconFinality(ctx Context, eventBlock uint64, tag string) bool {
	if eventBlock <= s.lastFinalizedExecBlock {
		log.Printf("[%s] beacon finality already covers block %d (finalized %d)",
			tag, eventBlock, s.lastFinalizedExecBlock)
		return true
	}
	log.Printf("[%s] waiting for beacon finality at block %d", tag, eventBlock)
	for attempt := 0; attempt < 60; attempt++ {
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
		log.Printf("[%s] beacon finalized block %d < event block %d, waiting... (%d/60)",
			tag, execBlock, eventBlock, attempt+1)
	}
	log.Printf("[%s] beacon finality did not reach block %d after 60 retries", tag, eventBlock)
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
	msg := err.Error()
	return strings.Contains(msg, "timeout elapsed") ||
		strings.Contains(msg, "IBCInvalidTimeoutTimestamp") ||
		strings.Contains(msg, "timed out")
}

func (s *Services) updateCosmosClientForEth(ctx Context, tag string) (*client.LightBlock, bool) {
	_, ethTrustedHeight := ctx.latestEthTimestamp.Snapshot()
	latestLightBlock, err := s.worker.UpdateCosmosClient(ctx, s.cosmosConfig.ProofType, int64(ethTrustedHeight), s.cosmosConfig.TrustLevel)
	if err != nil {
		log.Printf("[%s] Failed to update cosmos light client: %v", tag, err)
		return nil, false
	}
	if latestLightBlock == nil {
		log.Printf("[%s] Failed to update cosmos light client: latestLightBlock is nil", tag)
		return nil, false
	}

	ctx.latestEthTimestamp.Set(time.Now(), uint64(latestLightBlock.BlockHeight))

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

	s.BatchBuilder.PendingTracker.PurgeStale(pendingTrackerMaxAge)

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

	membershipMsg := tendermintContract.ILightClientMsgsMsgVerifyMembership{
		Height: tendermintContract.IICS02ClientMsgsHeight{
			RevisionHeight: uint64(height),
			RevisionNumber: 0,
		},
		KvPairs: []tendermintContract.IMembershipMsgsKVPair{
			{
				Path:  ibcPath,
				Value: value,
			},
		},
		MerkleProofs: []tendermintContract.IMembershipMsgsMerkleProof{merkleProof},
		AppHash:      utils.BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
		TrustedConsensusState: tendermintContract.IICS07TendermintMsgsConsensusState{
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

	nonMembershipMsg := tendermintContract.ILightClientMsgsMsgVerifyNonMembership{
		Height: tendermintContract.IICS02ClientMsgsHeight{
			RevisionHeight: uint64(height),
			RevisionNumber: 0,
		},
		KvPairs: []tendermintContract.IMembershipMsgsKVPair{
			{
				Path:  ibcPath,
				Value: value,
			},
		},
		MerkleProofs: []tendermintContract.IMembershipMsgsMerkleProof{merkleProof},
		AppHash:      utils.BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
		TrustedConsensusState: tendermintContract.IICS07TendermintMsgsConsensusState{
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

func (s *Services) ethProofHeight(ctx Context, eventBlock uint64, sequence uint64, tag string) (uint64, uint64, bool) {
	log.Printf("[%s] seq=%d: waiting for beacon finality at block %d", tag, sequence, eventBlock)
	finalized := false
	for attempt := 0; attempt < 60; attempt++ {
		if attempt > 0 {
			time.Sleep(10 * time.Second)
		}
		bctx, bcancel := context.WithTimeout(context.Background(), 15*time.Second)
		finalityUpdate, err := client.GetFinalityUpdate(bctx, ctx.BeaconAPIURL())
		bcancel()
		if err != nil {
			log.Printf("[%s] seq=%d: failed to get finality update: %v", tag, sequence, err)
			continue
		}
		execBlock, err := strconv.ParseUint(finalityUpdate.FinalizedHeader.Execution.BlockNumber, 10, 64)
		if err != nil {
			log.Printf("[%s] seq=%d: failed to parse finalized execution block %q: %v",
				tag, sequence, finalityUpdate.FinalizedHeader.Execution.BlockNumber, err)
			continue
		}
		if execBlock >= eventBlock {
			log.Printf("[%s] seq=%d: beacon finalized block %d >= event block %d", tag, sequence, execBlock, eventBlock)
			finalized = true
			break
		}
		log.Printf("[%s] seq=%d: beacon finalized block %d < event block %d, waiting... (%d/60)",
			tag, sequence, execBlock, eventBlock, attempt+1)
	}
	if !finalized {
		log.Printf("[%s] seq=%d: beacon finality did not reach block %d after 60 retries", tag, sequence, eventBlock)
		return 0, 0, false
	}

	if err := s.worker.UpdateEthClient(ctx); err != nil {
		log.Printf("[%s] seq=%d: failed to update ETH client: %v", tag, sequence, err)
		return 0, 0, false
	}

	ethClientState, err := client.GetEthereumClientState(ctx.CosmosClient(), ctx.EthClientID())
	if err != nil {
		log.Printf("[%s] seq=%d: failed to get ETH client state: %v", tag, sequence, err)
		return 0, 0, false
	}

	if ethClientState.LatestExecutionBlockNumber < eventBlock {
		log.Printf("[%s] seq=%d: ETH client at block %d still < event block %d after update, skipping",
			tag, sequence, ethClientState.LatestExecutionBlockNumber, eventBlock)
		return 0, 0, false
	}

	return ethClientState.LatestExecutionBlockNumber, ethClientState.LatestSlot, true
}

func parseMerkleProof(proofs []*ics23.CommitmentProof, sequence uint64) (tendermintContract.IMembershipMsgsMerkleProof, error) {
	merkleProof := tendermintContract.IMembershipMsgsMerkleProof{
		Proofs: []tendermintContract.IMembershipMsgsCommitmentProof{},
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
