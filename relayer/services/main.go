package services

import (
	"context"
	"encoding/binary"
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
			// update client on Eth side routinely
			if ctx.latestEthTimestamp.LatestUpdateTime.Add(routineInterval).Before(now) {
				latestBlock, err := s.worker.UpdateCosmosClient(ctx, s.cosmosConfig.ProofType, int64(ctx.latestEthTimestamp.LatestUpdateHeight), s.cosmosConfig.TrustLevel)
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

	// handle packets
	for {
		select {
		case batch, ok := <-s.CosmosPackets:
			if !ok {
				log.Println("[StartLoop] Cosmos batch channel closed, exiting loop")
				return
			}
			s.handleCosmos(ctx, batch)
		case batch, ok := <-s.EthPackets:
			if !ok {
				log.Println("[StartLoop] Eth batch channel closed, exiting loop")
				return
			}
			s.handleEth(ctx, batch)
		}
	}
}

func (s *Services) handleCosmos(ctx Context, batch CosmosBatch) {
	log.Printf("[StartLoop] Received cosmos batch: %d packets", len(batch.Packets))

	log.Printf("[StartLoop] Waiting 2 blocks for packet commitment to be included in AppHash...")
	time.Sleep(6 * time.Second)

	// V2: build the Cosmos→ETH updateClient msg but do NOT submit it as its
	// own tx. It will be the first inner call of the multicall so we save
	// one round-trip + ~21k base intrinsic gas per flush.
	updateBuild, err := s.worker.BuildCosmosClientUpdateMsg(
		ctx, s.cosmosConfig.ProofType,
		int64(ctx.latestEthTimestamp.LatestUpdateHeight),
		s.cosmosConfig.TrustLevel,
	)
	if err != nil {
		// The whole chunk was already sliced off the queue by CheckCosmos; if we
		// bail here without re-queueing it is lost. Re-queue so the next flush
		// retries once the transient cause (proof gen, RPC) clears (issue #80).
		log.Printf("[StartLoop] Failed to build cosmos light client update: %v", err)
		s.BatchBuilder.RequeueCosmos(batch.Packets)
		return
	}
	if updateBuild == nil || updateBuild.LightBlock == nil {
		log.Printf("[StartLoop] BuildCosmosClientUpdateMsg returned nil light block")
		s.BatchBuilder.RequeueCosmos(batch.Packets)
		return
	}
	latestLightBlock := updateBuild.LightBlock

	ethHeader, err := ctx.EthClient().HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Printf("[StartLoop] Failed to get eth block header: %v", err)
	}
	ethBlockTime := uint64(0)
	if ethHeader != nil {
		ethBlockTime = ethHeader.Time
	}

	// Build per-packet msgs into a single slice and submit one multicall when
	// N >= 2 (issue #67 / benchmark V1). Failures during pre-msg construction
	// (membership proof, missing ack bytes, etc.) still cause a per-packet
	// `continue` — that packet is just excluded from the batch.
	type batchedMsg struct {
		msg          any
		label        string        // log label, e.g. "RecvPacket"
		sequence     uint64        // for log lines
		sourceClient string        // for PendingTracker.Remove
		isRecv       bool          // true → call PendingTracker.Remove after success
		origin       *CosmosPacket // source packet to re-queue on failure; nil for folded updateClient
	}
	// +1 capacity for the optional updateClient prepend.
	msgs := make([]batchedMsg, 0, len(batch.Packets)+1)
	if updateBuild.HasMsg {
		msgs = append(msgs, batchedMsg{
			msg:   updateBuild.Msg,
			label: "UpdateClient",
		})
	}

	for _, packet := range batch.Packets {
		switch packet.Type {
		case CosmosSend:
			s.BatchBuilder.PendingTracker.Add(*packet.Packet, packet.BlockNumber)

			if ethBlockTime > 0 && packet.Packet.TimeoutTimestamp > 0 && ethBlockTime >= packet.Packet.TimeoutTimestamp {
				log.Printf("[RecvPacket] Packet seq=%d timed out (timeout=%d <= eth_block_time=%d), deferring to async timeout scanner",
					packet.Packet.Sequence, packet.Packet.TimeoutTimestamp, ethBlockTime)
				continue
			}

			calldata, err := s.cosmosMembership(ctx, *packet.Packet, packet.Packet.SourceClient, []byte{1}, latestLightBlock)
			if err != nil {
				log.Printf("[RecvPacket] seq=%d: %v", packet.Packet.Sequence, err)
				continue
			}

			msgs = append(msgs, batchedMsg{
				msg: contractICS26Router.IICS26RouterMsgsMsgRecvPacket{
					Packet:        toEthPacket(*packet.Packet),
					MembershipMsg: calldata,
				},
				label:        "RecvPacket",
				sequence:     packet.Packet.Sequence,
				sourceClient: packet.Packet.SourceClient,
				isRecv:       true,
				origin:       &packet,
			})
		case CosmosAck:
			if len(packet.AckBytes) == 0 {
				log.Printf("[AckPacket] seq=%d: acknowledgement bytes missing, skipping", packet.Packet.Sequence)
				continue
			}

			calldata, err := s.cosmosMembership(ctx, *packet.Packet, packet.Packet.DestinationClient, []byte{3}, latestLightBlock)
			if err != nil {
				log.Printf("[AckPacket] seq=%d: %v", packet.Packet.Sequence, err)
				continue
			}

			msgs = append(msgs, batchedMsg{
				msg: contractICS26Router.IICS26RouterMsgsMsgAckPacket{
					Packet:          toEthPacket(*packet.Packet),
					Acknowledgement: packet.AckBytes[0],
					MembershipMsg:   calldata,
				},
				label:    "AckPacket",
				sequence: packet.Packet.Sequence,
				origin:   &packet,
			})
		case CosmosTimeout:
			if packet.Packet.SourceClient == ctx.CosmosRouterClientID() {
				log.Printf("[Timeout] seq=%d: Cosmos-originated packet timeout already handled locally, skipping ETH relay", packet.Packet.Sequence)
				continue
			}

			calldata, err := s.cosmosNonMembership(ctx, *packet.Packet, packet.Packet.DestinationClient, []byte{2}, latestLightBlock)
			if err != nil {
				log.Printf("[Timeout] seq=%d: %v", packet.Packet.Sequence, err)
				continue
			}

			msgs = append(msgs, batchedMsg{
				msg: contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket{
					Packet:           toEthPacket(*packet.Packet),
					NonMembershipMsg: calldata,
				},
				label:    "Timeout",
				sequence: packet.Packet.Sequence,
				origin:   &packet,
			})
		default:
			log.Printf("[StartLoop] Unknown cosmos packet type: %d (seq=%d)", packet.Type, packet.Packet.Sequence)
		}
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
		// only the packet msgs need re-queueing (issue #80).
		s.BatchBuilder.RequeueCosmos(failed)
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
		ctx.latestEthTimestamp.mtx.Lock()
		ctx.latestEthTimestamp.LatestUpdateTime = time.Now()
		ctx.latestEthTimestamp.LatestUpdateHeight = uint64(latestLightBlock.BlockHeight)
		ctx.latestEthTimestamp.mtx.Unlock()
	}
}

func (s *Services) handleEth(ctx Context, batch EthBatch) {
	log.Printf("[StartLoop] Received eth batch: %d packets", len(batch.Packets))

	// V2: build per-packet Cosmos msgs into a single slice, fold the wasm
	// MsgUpdateClient(s) at the head, and submit one SendCosmosTxBatch.
	// Cosmos tx is atomic ⇒ MsgUpdateClient applies first, packet msgs verify
	// against the freshly-advanced state in the same block.
	type batchedCosmosMsg struct {
		msg      any
		label    string     // log label, e.g. "EthSend", "EthWriteAck", "UpdateClient"
		sequence uint64
		origin   *EthPacket // source packet to re-queue on failure; nil for folded updateClient
	}
	cosmosMsgs := make([]batchedCosmosMsg, 0, len(batch.Packets)+1)

	// Pre-filter expired EthSend packets — they bypass the batch and go to the
	// async timeout scanner. Anything else is kept for proof generation.
	type relayablePacket struct {
		packet EthPacket
		signer string
	}
	var relayable []relayablePacket
	maxEventBlock := uint64(0)
	for _, p := range batch.Packets {
		switch p.Type {
		case EthSend:
			if ethPacketExpired(p) {
				s.timeoutEthSend(ctx, p)
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
			log.Printf("[EthAck] seq=%d: terminal event handled", p.Packet.Sequence)
		case EthTimeout:
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
	requeueRelayable := func() {
		pkts := make([]EthPacket, len(relayable))
		for i, r := range relayable {
			pkts[i] = r.packet
		}
		s.BatchBuilder.RequeueEth(pkts)
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
		cosmosMsgs = append(cosmosMsgs, batchedCosmosMsg{msg: m, label: "UpdateClient"})
	}

	for _, r := range relayable {
		packet := r.packet
		origin := packet
		switch packet.Type {
		case EthSend:
			if ethPacketExpired(packet) {
				s.timeoutEthSend(ctx, packet)
				continue
			}
			proofBytes, err := client.GetEthMembershipProof(
				ctx.EthClient(), *ctx.RouterContract(), ethPath(packet.Packet.SourceClient, packet.Packet.Sequence, 1),
				ethcommon.HexToHash(ICS26_IBC_STORAGE_SLOT), new(big.Int).SetUint64(proofBlockNumber))
			if err != nil {
				log.Printf("[EthSend] seq=%d: failed to get ETH membership proof: %v", packet.Packet.Sequence, err)
				continue
			}
			cosmosMsgs = append(cosmosMsgs, batchedCosmosMsg{
				msg: &channeltypesv2.MsgRecvPacket{
					Packet:          *packet.Packet,
					ProofCommitment: proofBytes,
					ProofHeight:     clienttypes.Height{RevisionNumber: 0, RevisionHeight: proofSlot},
					Signer:          signerAddr,
				},
				label:    "EthSend",
				sequence: packet.Packet.Sequence,
				origin:   &origin,
			})
		case EthWriteAck:
			proofBytes, err := client.GetEthMembershipProof(
				ctx.EthClient(), *ctx.RouterContract(), ethPath(packet.Packet.DestinationClient, packet.Packet.Sequence, 3),
				ethcommon.HexToHash(ICS26_IBC_STORAGE_SLOT), new(big.Int).SetUint64(proofBlockNumber))
			if err != nil {
				log.Printf("[EthWriteAck] seq=%d: failed to get ETH membership proof: %v", packet.Packet.Sequence, err)
				continue
			}
			cosmosMsgs = append(cosmosMsgs, batchedCosmosMsg{
				msg: &channeltypesv2.MsgAcknowledgement{
					Packet: *packet.Packet,
					Acknowledgement: channeltypesv2.Acknowledgement{
						AppAcknowledgements: packet.AckBytes,
					},
					ProofAcked:  proofBytes,
					ProofHeight: clienttypes.Height{RevisionNumber: 0, RevisionHeight: proofSlot},
					Signer:      signerAddr,
				},
				label:    "EthWriteAck",
				sequence: packet.Packet.Sequence,
				origin:   &origin,
			})
		}
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
		failed := make([]EthPacket, 0, len(cosmosMsgs))
		for _, m := range cosmosMsgs {
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
		// packet msgs are re-queued (issue #80).
		s.BatchBuilder.RequeueEth(failed)
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
		ctx.latestCosmosTimestamp.mtx.Lock()
		ctx.latestCosmosTimestamp.LatestUpdateTime = time.Now()
		ctx.latestCosmosTimestamp.mtx.Unlock()
	}
}

// waitBeaconFinality polls beacon finality until execution block ≥ target.
// Extracted from the per-packet ethProofHeight so handleEth can amortize the
// wait over a whole batch. Returns false on timeout (60 polls × 10s = 10min).
func (s *Services) waitBeaconFinality(ctx Context, eventBlock uint64, tag string) bool {
	log.Printf("[%s] waiting for beacon finality at block %d", tag, eventBlock)
	for attempt := 0; attempt < 60; attempt++ {
		if attempt > 0 {
			time.Sleep(10 * time.Second)
		}
		finalityUpdate, err := client.GetFinalityUpdate(ctx.BeaconAPIURL())
		if err != nil {
			log.Printf("[%s] failed to get finality update: %v", tag, err)
			continue
		}
		execBlock, _ := strconv.ParseUint(finalityUpdate.FinalizedHeader.Execution.BlockNumber, 10, 64)
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
	latestLightBlock, err := s.worker.UpdateCosmosClient(ctx, s.cosmosConfig.ProofType, int64(ctx.latestEthTimestamp.LatestUpdateHeight), s.cosmosConfig.TrustLevel)
	if err != nil {
		log.Printf("[%s] Failed to update cosmos light client: %v", tag, err)
		return nil, false
	}
	if latestLightBlock == nil {
		log.Printf("[%s] Failed to update cosmos light client: latestLightBlock is nil", tag)
		return nil, false
	}

	ctx.latestEthTimestamp.mtx.Lock()
	ctx.latestEthTimestamp.LatestUpdateTime = time.Now()
	ctx.latestEthTimestamp.LatestUpdateHeight = uint64(latestLightBlock.BlockHeight)
	ctx.latestEthTimestamp.mtx.Unlock()

	return latestLightBlock, true
}

func (s *Services) timeoutEthSend(ctx Context, packet EthPacket) {
	log.Printf("[EthTimeout] seq=%d: packet expired, preparing timeout proof", packet.Packet.Sequence)

	latestLightBlock, ok := s.updateCosmosClientForEth(ctx, "EthTimeout")
	if !ok {
		return
	}

	counterpartyTime := uint64(latestLightBlock.SignedHeader.Header.Time.Unix())
	if counterpartyTime < packet.Packet.TimeoutTimestamp {
		log.Printf("[EthTimeout] seq=%d: counterparty time %d < timeout %d, skipping",
			packet.Packet.Sequence, counterpartyTime, packet.Packet.TimeoutTimestamp)
		return
	}

	calldata, err := s.cosmosNonMembership(ctx, *packet.Packet, packet.Packet.DestinationClient, []byte{2}, latestLightBlock)
	if err != nil {
		log.Printf("[EthTimeout] seq=%d: %v", packet.Packet.Sequence, err)
		return
	}

	msgTimeoutPacket := contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket{
		Packet:           toEthPacket(*packet.Packet),
		NonMembershipMsg: calldata,
	}

	if err := s.worker.TxHandler.SendEthTx(ctx, msgTimeoutPacket); err != nil {
		log.Printf("[EthTimeout] seq=%d: SendEthTx failed: %v", packet.Packet.Sequence, err)
		return
	}
	log.Printf("[EthTimeout] seq=%d: relay completed", packet.Packet.Sequence)
}

const pendingTrackerMaxAge = 1 * time.Hour

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
		return
	}

	for _, info := range processed {
		s.BatchBuilder.PendingTracker.Remove(info.Packet.SourceClient, info.Packet.Sequence)
		log.Printf("[CosmosTimeout] seq=%d: timeout relay completed (bundled with %d update msgs)", info.Packet.Sequence, len(updateResult.Msgs))
	}
}

func (s *Services) buildCosmosTimeoutMsg(ctx Context, packet channeltypesv2.Packet, ethClientState *client.EthereumClientState) (*channeltypesv2.MsgTimeout, error) {
	receiptPath := ethPath(packet.DestinationClient, packet.Sequence, 2)
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
		finalityUpdate, err := client.GetFinalityUpdate(ctx.BeaconAPIURL())
		if err != nil {
			log.Printf("[%s] seq=%d: failed to get finality update: %v", tag, sequence, err)
			continue
		}
		execBlock, _ := strconv.ParseUint(finalityUpdate.FinalizedHeader.Execution.BlockNumber, 10, 64)
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

func ethPath(clientID string, sequence uint64, pathType byte) []byte {
	seqBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(seqBytes, sequence)
	path := append([]byte(clientID), pathType)
	return append(path, seqBytes...)
}
