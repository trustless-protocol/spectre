package services

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/big"
	"relayer/chain/wasmclient"
	"relayer/utils"
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
	CreateCosmosClientContract(stdCtx context.Context, endpoint EVMEndpoint, clientIDs ClientIDs, clientState []byte, consensusState spectreContract.IICS07TendermintMsgsConsensusState, initialPinnedValidatorSet client.ContractValidatorSet) (ethcommon.Address, error)
	// CreateWasmClient submits MsgCreateClient for any 08-wasm client (ETH beacon,
	// L2 rollup, ...) from an already-built wasm ClientState/ConsensusState and returns
	// the auto-assigned client id. counterpartyClientID is the client on the
	// counterparty chain that tracks Cosmos, registered inline (a naming binding); an
	// empty value skips registration (register it later once its id is known).
	CreateWasmClient(stdCtx context.Context, endpoint CosmosEndpoint, clientState ibcexported.ClientState, consensusState ibcexported.ConsensusState, counterpartyClientID string) (string, error)
	SendEthTx(stdCtx context.Context, endpoint EVMEndpoint, cosmosRouterClientID string, msg any) error
	SendEthTxBatch(stdCtx context.Context, endpoint EVMEndpoint, cosmosRouterClientID string, msgs []any) error
	SubmitMisbehaviour(stdCtx context.Context, endpoint EVMEndpoint, cosmosRouterClientID string, misbehaviourMsg []byte) error
	SendCosmosTxBatch(stdCtx context.Context, endpoint CosmosEndpoint, msgs []any) error
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
	GenerateMisbehaviourProof(header1Sigs, header2Sigs []prover.ValidatorSignature) (prover.MisbehaviourProof, error)
}

// Services holds the shared pieces the chain adapters reuse: the Worker (tx
// handler + prover), the resolved Cosmos-source config, and the BatchBuilder the
// subscriber pushes into and the adapters drain. (The legacy StartLoop's own
// event channels, listener and cross-chunk memoization were removed with it.)
type Services struct {
	worker        *Worker
	cosmosConfig  Config
	BatchBuilder  *BatchBuilder
	clientUpdates clientUpdateTimes
	recovery      *RecoveryStateStore
}

type evmTimeoutDeps struct {
	cosmos         CosmosEndpoint
	evm            EVMEndpoint
	routerClientID string
}

func New(txHandler TransactionHandler, prover Prover, cosmosConfig Config, recovery ...*RecoveryStateStore) *Services {
	var recoveryStore *RecoveryStateStore
	if len(recovery) > 0 {
		recoveryStore = recovery[0]
	}
	return &Services{
		cosmosConfig: cosmosConfig,
		worker: &Worker{
			txHandler,
			prover,
		},
		BatchBuilder: NewBatchBuilder(),
		recovery:     recoveryStore,
	}
}

// RecoveryState returns the durable recovery store wired for this source. It is
// nil for one-shot commands and source types that do not use gap recovery.
func (s *Services) RecoveryState() *RecoveryStateStore { return s.recovery }

// NewWithPendingState is the production relay constructor. It rehydrates
// pending packets and timeout tombstones before a Services value can be handed
// to source subscriptions or recovery loops.
func NewWithPendingState(txHandler TransactionHandler, prover Prover, cosmosConfig Config, stateDir string, recovery ...*RecoveryStateStore) (*Services, error) {
	batchBuilder, err := NewPersistentBatchBuilder(stateDir)
	if err != nil {
		return nil, fmt.Errorf("restore pending packet trackers: %w", err)
	}
	var recoveryStore *RecoveryStateStore
	if len(recovery) > 0 {
		recoveryStore = recovery[0]
	}
	return &Services{
		cosmosConfig: cosmosConfig,
		worker: &Worker{
			txHandler,
			prover,
		},
		BatchBuilder: batchBuilder,
		recovery:     recoveryStore,
	}, nil
}

// CosmosClientExpiry returns when the on-chain Cosmos SpectreClient expires: the
// trusted consensus timestamp plus the trusting period. Unlike
// seedCosmosClientFreshness this is side-effect-free (no timestamp mutation, no
// logging), so the relay module's refresh routine can poll it periodically.
func CosmosClientExpiry(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint) (time.Time, time.Duration, error) {
	readCtx, cancelRead := fetchCtx(stdCtx, defaultFetchTimeout)
	clientState, err := fetchOnChainClientState(readCtx, evm)
	cancelRead()
	if err != nil {
		return time.Time{}, 0, err
	}
	trustedHeight, err := clientStateRevisionHeightInt64(clientState)
	if err != nil {
		return time.Time{}, 0, err
	}
	lightCtx, cancelLight := fetchCtx(stdCtx, defaultFetchTimeout)
	defer cancelLight()
	lightBlock, err := client.GetLightBlock(lightCtx, cosmos.CosmosClient(), trustedHeight)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("cosmos client expiry: trusted light block %d: %w", trustedHeight, err)
	}
	period := time.Duration(clientState.TrustingPeriod) * time.Second
	return tendermintClientExpiry(lightBlock.SignedHeader.Header.Time, clientState.TrustingPeriod), period, nil
}

// tendermintClientExpiry derives when the Tendermint light client on the EVM side
// stops accepting updates: the trusting period measured from the TRUSTED HEADER's
// time, not from now.
//
// Measuring from now would report a client as healthy forever, since "now" always
// moves with the check. The clock that matters belongs to the block the client
// currently trusts.
//
// This is the mirror of ethClientExpiry in chain/cosmos, and the two are shaped
// differently on purpose: a Tendermint client has one explicit trusting period,
// while a beacon client has none and its expiry has to be inferred from sync
// committee geometry.
func tendermintClientExpiry(trustedHeaderTime time.Time, trustingPeriodSeconds uint32) time.Time {
	return trustedHeaderTime.Add(time.Duration(trustingPeriodSeconds) * time.Second)
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

// RefreshSafetyMargin is exported because the relay module needs the same rule
// the refresh-interval calculation already uses: a margin that is a FRACTION of
// the trusting period, not a fixed number of minutes. Two rules for one concept
// is how the module ended up with a hard-coded 30 minutes that has no relation
// to any client's real trusting period.
func RefreshSafetyMargin(trustingPeriod time.Duration) time.Duration {
	return refreshSafetyMargin(trustingPeriod)
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

func (s *Services) updateCosmosClientForEth(stdCtx context.Context, deps evmTimeoutDeps, tag string) (*client.LightBlock, bool) {
	latestLightBlock, err := s.worker.UpdateCosmosClient(stdCtx, deps.cosmos, deps.evm, deps.routerClientID, s.cosmosConfig.FetchTimeout, s.cosmosConfig.RotationThreshold, s.cosmosConfig.ProofType, 0, s.cosmosConfig.TrustLevel, false, 0)
	if err != nil {
		log.Printf("[%s] Failed to update cosmos light client: %v", tag, err)
		return nil, false
	}
	if latestLightBlock == nil {
		log.Printf("[%s] Failed to update cosmos light client: latestLightBlock is nil", tag)
		return nil, false
	}
	s.ObserveCosmosOnEVMUpdate(latestLightBlock.SignedHeader.Header.Time)

	return latestLightBlock, true
}

func (s *Services) timeoutEVMSend(stdCtx context.Context, deps evmTimeoutDeps, packet EthPacket, tag string, latestLightBlock *client.LightBlock) timeoutSendOutcome {
	log.Printf("[%sTimeout] seq=%d: packet expired, preparing timeout proof", tag, packet.Packet.Sequence)
	counterpartyTime := uint64(latestLightBlock.SignedHeader.Header.Time.Unix())
	if counterpartyTime < packet.Packet.TimeoutTimestamp {
		log.Printf("[%sTimeout] seq=%d: counterparty time %d < timeout %d, skipping", tag, packet.Packet.Sequence, counterpartyTime, packet.Packet.TimeoutTimestamp)
		return timeoutNotDue
	}

	calldata, err := CosmosNonMembership(stdCtx, deps.cosmos, *packet.Packet, packet.Packet.DestinationClient, []byte{2}, latestLightBlock)
	if err != nil {
		log.Printf("[%sTimeout] seq=%d: %v", tag, packet.Packet.Sequence, err)
		return timeoutOutcomeForError(err)
	}
	msgTimeoutPacket := contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket{
		Packet:           ToEthPacket(*packet.Packet),
		NonMembershipMsg: calldata,
	}
	if err := s.worker.TxHandler.SendEthTx(stdCtx, deps.evm, deps.routerClientID, msgTimeoutPacket); err != nil {
		log.Printf("[%sTimeout] seq=%d: SendEthTx failed: %v", tag, packet.Packet.Sequence, err)
		return timeoutOutcomeForError(err)
	}
	log.Printf("[%sTimeout] seq=%d: relay completed", tag, packet.Packet.Sequence)
	return timeoutSent
}

const pendingTrackerMaxAge = 1 * time.Hour

func (s *Services) scanForEthTimeouts(stdCtx context.Context, deps evmTimeoutDeps) {
	s.scanForEVMTimeouts(stdCtx, deps, evmTimeoutScanOptions{
		tag:     "Eth",
		tracker: s.BatchBuilder.EthPendingTracker,
		prepareTimeouts: func(c context.Context, scanCtx evmTimeoutDeps) (*client.LightBlock, bool) {
			return s.updateCosmosClientForEth(c, scanCtx, "EthTimeout")
		},
		hasPendingCommitment: func(c context.Context, scanCtx evmTimeoutDeps, packet channeltypesv2.Packet) (bool, error) {
			return HasPendingEthPacketCommitment(c, scanCtx.evm, packet)
		},
		timeoutSend: func(c context.Context, scanCtx evmTimeoutDeps, packet EthPacket, lightBlock *client.LightBlock) timeoutSendOutcome {
			return s.timeoutEVMSend(c, scanCtx, packet, "Eth", lightBlock)
		},
	})
}

func (s *Services) scanForL2Timeouts(stdCtx context.Context, deps evmTimeoutDeps) {
	s.scanForEVMTimeouts(stdCtx, deps, evmTimeoutScanOptions{
		tag:     "L2",
		tracker: s.BatchBuilder.L2PendingTracker,
		prepareTimeouts: func(c context.Context, scanCtx evmTimeoutDeps) (*client.LightBlock, bool) {
			return s.updateCosmosClientForEth(c, scanCtx, "L2Timeout")
		},
		hasPendingCommitment: func(c context.Context, scanCtx evmTimeoutDeps, packet channeltypesv2.Packet) (bool, error) {
			return HasPendingEthPacketCommitment(c, scanCtx.evm, packet)
		},
		timeoutSend: func(c context.Context, scanCtx evmTimeoutDeps, packet EthPacket, lightBlock *client.LightBlock) timeoutSendOutcome {
			return s.timeoutEVMSend(c, scanCtx, packet, "L2", lightBlock)
		},
	})
}

type evmTimeoutScanOptions struct {
	tag                  string
	tracker              *PendingPacketTracker
	prepareTimeouts      func(context.Context, evmTimeoutDeps) (*client.LightBlock, bool)
	hasPendingCommitment func(context.Context, evmTimeoutDeps, channeltypesv2.Packet) (bool, error)
	timeoutSend          func(context.Context, evmTimeoutDeps, EthPacket, *client.LightBlock) timeoutSendOutcome
}

func (s *Services) scanForEVMTimeouts(stdCtx context.Context, deps evmTimeoutDeps, opts evmTimeoutScanOptions) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[%sTimeoutScan] Panic recovered: %v", opts.tag, r)
		}
	}()
	if opts.tracker == nil {
		log.Printf("[%sTimeoutScan] pending tracker is nil", opts.tag)
		return
	}
	if allowed, err := opts.tracker.TimeoutRetriesAllowed(); !allowed {
		if err != nil {
			log.Printf("[%sTimeoutScan][ATTENTION] timeout retry circuit breaker remains open: %v", opts.tag, err)
		}
		return
	}
	if err := opts.tracker.PurgeStaleWithoutTimeout(pendingTrackerMaxAge); err != nil {
		log.Printf("[%sTimeoutScan][ATTENTION] failed to persist pending-state cleanup: %v", opts.tag, err)
		return
	}
	now := time.Now()
	pending := opts.tracker.GetDue(now)
	if len(pending) == 0 {
		return
	}
	expired := pendingPacketsTimedOutAtTimestamp(pending, uint64(now.Unix()))
	if len(expired) == 0 {
		return
	}
	log.Printf("[%sTimeoutScan] Found %d locally-expired EVM-origin packet(s)", opts.tag, len(expired))
	candidates := make([]pendingPacketInfo, 0, len(expired))
	for _, info := range expired {
		pendingCommitment, err := opts.hasPendingCommitment(stdCtx, deps, info.Packet)
		if err != nil {
			log.Printf("[%sTimeoutScan] seq=%d: failed to check EVM packet commitment: %v", opts.tag, info.Packet.Sequence, err)
			if !applyTimeoutOutcome(opts.tracker, info, timeoutDeferred, time.Now(), opts.tag+"Timeout") {
				return
			}
			continue
		}
		if !pendingCommitment {
			if err := opts.tracker.RemoveIfCurrent(info); err != nil {
				log.Printf("[%sTimeoutScan][ATTENTION] failed to persist cleared commitment: %v", opts.tag, err)
				return
			}
			log.Printf("[%sTimeoutScan] seq=%d: EVM commitment already cleared, removed from pending tracker", opts.tag, info.Packet.Sequence)
			continue
		}
		candidates = append(candidates, info)
	}
	if len(candidates) == 0 {
		return
	}

	// The Cosmos client update is a shared prerequisite for every timeout proof
	// in this scan. Resolve it once: a backlog of N expired packets must not
	// trigger N expensive update/proof attempts when that prerequisite is down.
	latestLightBlock, ok := opts.prepareTimeouts(stdCtx, deps)
	if !ok {
		deferTimeoutRetries(opts.tracker, candidates, time.Now(), opts.tag+"Timeout")
		return
	}
	if err := opts.tracker.ClearDeferralsIfCurrentBatch(candidates); err != nil {
		log.Printf("[%sTimeoutScan][ATTENTION] failed to persist recovered prerequisite state: %v", opts.tag, err)
		return
	}
	for _, info := range candidates {
		// The shared client update just succeeded. Any prior shared-outage
		// deferrals no longer describe this packet, so let a subsequent packet-
		// specific not-due or transient result start at the one-minute backoff.
		packet := info.Packet
		outcome := opts.timeoutSend(stdCtx, deps, EthPacket{Type: EthSend, Packet: &packet, BlockNumber: info.BlockNumber}, latestLightBlock)
		if !applyTimeoutOutcome(opts.tracker, info, outcome, time.Now(), opts.tag+"Timeout") {
			return
		}
	}
}

func (s *Services) scanForCosmosTimeouts(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint, ethClientID string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[CosmosTimeoutScan] Panic recovered: %v", r)
		}
	}()

	tracker := s.BatchBuilder.PendingTracker
	if allowed, err := tracker.TimeoutRetriesAllowed(); !allowed {
		if err != nil {
			log.Printf("[CosmosTimeoutScan][ATTENTION] timeout retry circuit breaker remains open: %v", err)
		}
		return
	}
	if err := tracker.PurgeStaleWithoutTimeout(pendingTrackerMaxAge); err != nil {
		log.Printf("[CosmosTimeoutScan][ATTENTION] failed to persist pending-state cleanup: %v", err)
		return
	}
	now := time.Now()
	pending := tracker.GetDue(now)
	if len(pending) == 0 {
		return
	}

	headerCtx, cancelHeader := fetchCtx(stdCtx, defaultFetchTimeout)
	ethHeader, err := evm.EthClient().HeaderByNumber(headerCtx, nil)
	cancelHeader()
	if err != nil {
		log.Printf("[CosmosTimeoutScan] Failed to get ETH block header: %v", err)
		deferTimeoutRetries(tracker, pending, time.Now(), "CosmosTimeout")
		return
	}
	expired := pendingPacketsTimedOutAtTimestamp(pending, ethHeader.Time)
	if len(expired) == 0 {
		return
	}

	unreceived := make([]pendingPacketInfo, 0, len(expired))
	for _, info := range expired {
		received, err := HasEthPacketReceipt(stdCtx, evm, info.Packet)
		if err != nil {
			log.Printf("[CosmosTimeoutScan] seq=%d: failed to check ETH packet receipt: %v", info.Packet.Sequence, err)
			if !deferTimeoutRetries(tracker, []pendingPacketInfo{info}, time.Now(), "CosmosTimeout") {
				return
			}
			continue
		}
		if received {
			if err := tracker.RemoveIfCurrent(info); err != nil {
				log.Printf("[CosmosTimeoutScan][ATTENTION] failed to persist received packet removal: %v", err)
				return
			}
			log.Printf("[CosmosTimeoutScan] seq=%d: ETH receipt already exists, removed from pending tracker", info.Packet.Sequence)
			continue
		}
		unreceived = append(unreceived, info)
	}
	expired = unreceived
	if len(expired) == 0 {
		return
	}

	updateResult, ethClientState, ok := s.resolveCosmosTimeoutProofState(
		tracker,
		expired,
		func() (*EthClientUpdateResult, error) {
			return s.worker.BuildEthClientUpdateHeaders(stdCtx, cosmos, evm, ethClientID)
		},
		func() (*client.EthereumClientState, error) {
			readCtx, cancel := fetchCtx(stdCtx, defaultFetchTimeout)
			defer cancel()
			return client.GetEthereumClientState(readCtx, cosmos.CosmosClient(), ethClientID)
		},
	)
	if !ok {
		return
	}
	if len(updateResult.Headers) == 0 && updateResult.ProofTimestamp > 0 {
		s.ObserveEVMOnCosmosUpdate(time.Unix(int64(updateResult.ProofTimestamp), 0))
	}

	expired = pendingPacketsTimedOutAtTimestamp(expired, updateResult.ProofTimestamp)
	if len(expired) == 0 {
		// The shared proof/update path succeeded; a packet merely is not due at
		// this proof timestamp, so stale outage deferrals must not keep it stuck.
		if err := tracker.ClearDeferralsIfCurrentBatch(pending); err != nil {
			log.Printf("[CosmosTimeoutScan][ATTENTION] failed to persist not-due state: %v", err)
		}
		return
	}

	timeoutMsgs := make([]any, 0, len(expired))
	processed := make([]pendingPacketInfo, 0, len(expired))
	for _, info := range expired {
		msgTimeout, err := s.buildCosmosTimeoutMsg(stdCtx, evm, info.Packet, ethClientState)
		if err != nil {
			if errors.Is(err, client.ErrPacketAlreadyReceived) {
				if err := tracker.RemoveIfCurrent(info); err != nil {
					log.Printf("[CosmosTimeoutScan][ATTENTION] failed to persist received packet removal: %v", err)
					return
				}
				log.Printf("[CosmosTimeout] seq=%d: receipt exists at proof height, removed from pending tracker", info.Packet.Sequence)
				continue
			}
			log.Printf("[CosmosTimeout] seq=%d: %v", info.Packet.Sequence, err)
			if !deferTimeoutRetries(tracker, []pendingPacketInfo{info}, time.Now(), "CosmosTimeout") {
				return
			}
			continue
		}
		timeoutMsgs = append(timeoutMsgs, msgTimeout)
		processed = append(processed, info)
	}
	if len(timeoutMsgs) == 0 {
		return
	}
	if err := tracker.ClearDeferralsIfCurrentBatch(processed); err != nil {
		log.Printf("[CosmosTimeoutScan][ATTENTION] failed to persist recovered prerequisite state: %v", err)
		return
	}

	updateMsgs := make([]any, 0, len(updateResult.Headers))
	for i, header := range updateResult.Headers {
		msg, err := wasmclient.BuildUpdateClient("", ethClientID, header)
		if err != nil {
			log.Printf("[CosmosTimeoutScan] Failed to wrap ETH update header %d: %v", i, err)
			deferTimeoutRetries(tracker, processed, time.Now(), "CosmosTimeout")
			return
		}
		updateMsgs = append(updateMsgs, msg)
	}
	if len(updateMsgs) > 0 {
		if err := s.worker.WaitForCosmosCatchUp(stdCtx, cosmos, ethClientState, updateResult.SigSlot); err != nil {
			log.Printf("[CosmosTimeoutScan] target chain did not catch up: %v", err)
			deferTimeoutRetries(tracker, processed, time.Now(), "CosmosTimeout")
			return
		}
	}

	batchMsgs := append(append(make([]any, 0, len(updateMsgs)+len(timeoutMsgs)), updateMsgs...), timeoutMsgs...)
	if err := s.worker.TxHandler.SendCosmosTxBatch(stdCtx, cosmos, batchMsgs); err != nil {
		log.Printf("[CosmosTimeoutScan] SendCosmosTxBatch failed: %v", err)
		s.handleCosmosTimeoutBatchFailure(stdCtx, cosmos, tracker, updateMsgs, timeoutMsgs, processed, err)
		return
	}
	for _, info := range processed {
		if err := tracker.RemoveIfCurrent(info); err != nil {
			log.Printf("[CosmosTimeoutScan][ATTENTION] timeout succeeded but completion state was not durable: %v", err)
			return
		}
		log.Printf("[CosmosTimeout] seq=%d: timeout relay completed", info.Packet.Sequence)
	}
	if updateResult.ProofTimestamp > 0 {
		s.ObserveEVMOnCosmosUpdate(time.Unix(int64(updateResult.ProofTimestamp), 0))
	}
}

// resolveCosmosTimeoutProofState keeps the shared-prerequisite deferral beside
// both failure exits. Every packet in the scan depends on the same update and
// client state, so neither failure is attributable to an individual packet;
// all expired packets back off without consuming their permanent-failure budget.
func (s *Services) resolveCosmosTimeoutProofState(
	tracker *PendingPacketTracker,
	expired []pendingPacketInfo,
	buildUpdate func() (*EthClientUpdateResult, error),
	fetchClientState func() (*client.EthereumClientState, error),
) (*EthClientUpdateResult, *client.EthereumClientState, bool) {
	updateResult, err := buildUpdate()
	if err != nil {
		log.Printf("[CosmosTimeoutScan] Failed to build ETH client update headers: %v; deferring %d timeout(s) without charging", err, len(expired))
		deferTimeoutRetries(tracker, expired, time.Now(), "CosmosTimeout")
		return nil, nil, false
	}

	ethClientState := updateResult.EthClientState
	if ethClientState == nil {
		ethClientState, err = fetchClientState()
		if err != nil {
			log.Printf("[CosmosTimeoutScan] Failed to get ETH client state: %v; deferring %d timeout(s) without charging", err, len(expired))
			deferTimeoutRetries(tracker, expired, time.Now(), "CosmosTimeout")
			return nil, nil, false
		}
	}
	return updateResult, ethClientState, true
}

// handleCosmosTimeoutBatchFailure preserves any successful prefix, never charges
// packets behind an update-prefix failure, and isolates permanent timeout
// failures into singleton submissions so only the actual poison packet is charged.
func (s *Services) handleCosmosTimeoutBatchFailure(stdCtx context.Context, cosmos CosmosEndpoint, tracker *PendingPacketTracker, updateMsgs, timeoutMsgs []any, processed []pendingPacketInfo, sendErr error) {
	remainingInfos := processed
	remainingMsgs := timeoutMsgs
	if updateMsgsFailed(sendErr, len(updateMsgs)) {
		deferTimeoutRetries(tracker, remainingInfos, time.Now(), "CosmosTimeout")
		return
	}
	var partialErr *BatchPartialError
	if errors.As(sendErr, &partialErr) {
		succeeded := partialErr.SucceededCount - len(updateMsgs)
		if succeeded > len(remainingInfos) {
			succeeded = len(remainingInfos)
		}
		for _, info := range remainingInfos[:succeeded] {
			if err := tracker.RemoveIfCurrent(info); err != nil {
				log.Printf("[CosmosTimeout][ATTENTION] successful timeout completion was not durable: %v", err)
				return
			}
		}
		remainingInfos = remainingInfos[succeeded:]
		remainingMsgs = remainingMsgs[succeeded:]
	}
	if len(remainingInfos) == 0 {
		return
	}
	if !errors.Is(sendErr, ErrPermanentRelayFailure) {
		deferTimeoutRetries(tracker, remainingInfos, time.Now(), "CosmosTimeout")
		return
	}
	for i, info := range remainingInfos {
		err := s.worker.TxHandler.SendCosmosTxBatch(stdCtx, cosmos, []any{remainingMsgs[i]})
		if err == nil {
			if err := tracker.RemoveIfCurrent(info); err != nil {
				log.Printf("[CosmosTimeout][ATTENTION] successful timeout completion was not durable: %v", err)
				return
			}
			continue
		}
		if errors.Is(err, ErrPermanentRelayFailure) {
			if !chargeTimeoutFailure(tracker, info, time.Now(), "CosmosTimeout") {
				return
			}
		} else {
			if !deferTimeoutRetries(tracker, []pendingPacketInfo{info}, time.Now(), "CosmosTimeout") {
				return
			}
		}
	}
}

// updateMsgsFailed reports whether a failed [update..., timeout...] batch
// stopped before any timeout message could execute. Such a failure belongs to
// the shared client-update prefix and must defer every packet without charging
// its packet-attributable retry budget.
func updateMsgsFailed(err error, updateMsgCount int) bool {
	if updateMsgCount == 0 {
		return false
	}
	var partialErr *BatchPartialError
	if errors.As(err, &partialErr) {
		return partialErr.SucceededCount < updateMsgCount
	}
	return true
}

// buildCosmosTimeoutMsg takes stdCtx so the eth_getProof it issues is bound to
// the scan's lifetime. Without it the scanner could block forever on an
// unresponsive node, and because relay.Module runs the scan synchronously that
// stops timeout recovery entirely rather than just delaying one packet.
func (s *Services) buildCosmosTimeoutMsg(stdCtx context.Context, evm EVMEndpoint, packet channeltypesv2.Packet, ethClientState *client.EthereumClientState) (*channeltypesv2.MsgTimeout, error) {
	receiptPath := EthPath(packet.DestinationClient, packet.Sequence, 2)
	proofCtx, cancel := fetchCtx(stdCtx, defaultFetchTimeout)
	defer cancel()
	proofBytes, err := client.GetEthNonMembershipProof(
		proofCtx, evm.EthClient(), evm.Contracts.Router, receiptPath, ethcommon.HexToHash(ICS26_IBC_STORAGE_SLOT), new(big.Int).SetUint64(ethClientState.LatestExecutionBlockNumber))
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

func CosmosMembership(stdCtx context.Context, ctx CosmosEndpoint, packet channeltypesv2.Packet, clientID string, pathType []byte, latestLightBlock *client.LightBlock) ([]byte, error) {
	height := latestLightBlock.BlockHeight
	ibcPath := utils.IbcPath(clientID, packet.Sequence, pathType)
	value, proof, err := client.ProvePath(stdCtx, ctx.CosmosClient(), height, ibcPath)
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

func CosmosNonMembership(stdCtx context.Context, ctx CosmosEndpoint, packet channeltypesv2.Packet, clientID string, pathType []byte, latestLightBlock *client.LightBlock) ([]byte, error) {
	height := latestLightBlock.BlockHeight
	ibcPath := utils.IbcPath(clientID, packet.Sequence, pathType)
	value, proof, err := client.ProvePath(stdCtx, ctx.CosmosClient(), height, ibcPath)
	if err != nil {
		return nil, err
	}
	if len(value) != 0 {
		return nil, cosmosReceiptPresentError(height, len(value))
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

// cosmosReceiptPresentError preserves the already-received sentinel on the
// non-membership path. Without it the timeout scanner mistakes a delivered
// packet for a transient proof failure and defers it forever.
func cosmosReceiptPresentError(height int64, valueLen int) error {
	return fmt.Errorf("non-membership expected empty value at height=%d, got %d bytes: %w",
		height, valueLen, client.ErrPacketAlreadyReceived)
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

func ToEthPacket(packet channeltypesv2.Packet) contractICS26Router.IICS26RouterMsgsPacket {
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

func HasEthIBCPathValue(stdCtx context.Context, ctx EVMEndpoint, path []byte) (bool, error) {
	if ctx.EthClient() == nil {
		return false, fmt.Errorf("eth client is nil")
	}
	if ctx.RouterContract() == nil {
		return false, fmt.Errorf("router contract address is nil")
	}

	callCtx, cancel := fetchCtx(stdCtx, defaultFetchTimeout)
	defer cancel()

	value, err := ctx.EthClient().StorageAt(callCtx, *ctx.RouterContract(), EthIBCStorageKey(path), nil)
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

func HasEthPacketReceipt(stdCtx context.Context, ctx EVMEndpoint, packet channeltypesv2.Packet) (bool, error) {
	return HasEthIBCPathValue(stdCtx, ctx, EthPath(packet.DestinationClient, packet.Sequence, 2))
}

func HasPendingEthPacketCommitment(stdCtx context.Context, ctx EVMEndpoint, packet channeltypesv2.Packet) (bool, error) {
	return HasEthIBCPathValue(stdCtx, ctx, EthPath(packet.SourceClient, packet.Sequence, 1))
}

func NewWorker(txHandler TransactionHandler, prover Prover) *Worker {
	return &Worker{
		txHandler,
		prover,
	}
}
