package services

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/big"
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
	CreateCosmosClientContract(stdCtx context.Context, ctx Context, clientState, consensusHash []byte, initialPinnedValidatorSet client.ContractValidatorSet) (ethcommon.Address, error)
	// CreateWasmClient submits MsgCreateClient for any 08-wasm client (ETH beacon,
	// L2 rollup, ...) from an already-built wasm ClientState/ConsensusState and
	// returns the auto-assigned client id. registerCounterparty additionally binds
	// the counterparty client id (the ETH-side router client): true for the ETH
	// beacon client, false for an L2 bootstrap whose counterparty does not exist yet.
	CreateWasmClient(stdCtx context.Context, ctx Context, clientState ibcexported.ClientState, consensusState ibcexported.ConsensusState, registerCounterparty bool) (string, error)
	SendEthTx(stdCtx context.Context, ctx Context, msg any) error
	SendEthTxBatch(stdCtx context.Context, ctx Context, msgs []any) error
	SendCosmosTxBatch(stdCtx context.Context, ctx Context, msgs []any) error
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

// Services holds the shared pieces the chain adapters reuse: the Worker (tx
// handler + prover), the resolved Cosmos-source config, and the BatchBuilder the
// subscriber pushes into and the adapters drain. (The legacy StartLoop's own
// event channels, listener and cross-chunk memoization were removed with it.)
type Services struct {
	worker       *Worker
	cosmosConfig Config
	BatchBuilder *BatchBuilder
}

func New(txHandler TransactionHandler, prover Prover, cosmosConfig Config) *Services {
	return &Services{
		cosmosConfig: cosmosConfig,
		worker: &Worker{
			txHandler,
			prover,
		},
		BatchBuilder: NewBatchBuilder(),
	}
}

// CosmosClientExpiry returns when the on-chain Cosmos SpectreClient expires: the
// trusted consensus timestamp plus the trusting period. Unlike
// seedCosmosClientFreshness this is side-effect-free (no timestamp mutation, no
// logging), so the relay module's refresh routine can poll it periodically.
func CosmosClientExpiry(ctx Context) (time.Time, error) {
	clientState, err := fetchOnChainClientState(ctx)
	if err != nil {
		return time.Time{}, err
	}
	trustedHeight, err := clientStateRevisionHeightInt64(clientState)
	if err != nil {
		return time.Time{}, err
	}
	lightBlock, err := client.GetLightBlock(ctx.CosmosClient(), trustedHeight)
	if err != nil {
		return time.Time{}, fmt.Errorf("cosmos client expiry: trusted light block %d: %w", trustedHeight, err)
	}
	trustingPeriod := time.Duration(clientState.TrustingPeriod) * time.Second
	return lightBlock.SignedHeader.Header.Time.Add(trustingPeriod), nil
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

func (s *Services) updateCosmosClientForEth(stdCtx context.Context, ctx Context, tag string) (*client.LightBlock, bool) {
	_, ethTrustedHeight := ctx.latestEthTimestamp.Snapshot()
	latestLightBlock, err := s.worker.UpdateCosmosClient(stdCtx, ctx, s.cosmosConfig.ProofType, int64(ethTrustedHeight), s.cosmosConfig.TrustLevel, false)
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

func (s *Services) timeoutEthSend(stdCtx context.Context, ctx Context, packet EthPacket) bool {
	log.Printf("[EthTimeout] seq=%d: packet expired, preparing timeout proof", packet.Packet.Sequence)

	latestLightBlock, ok := s.updateCosmosClientForEth(stdCtx, ctx, "EthTimeout")
	if !ok {
		return false
	}

	counterpartyTime := uint64(latestLightBlock.SignedHeader.Header.Time.Unix())
	if counterpartyTime < packet.Packet.TimeoutTimestamp {
		log.Printf("[EthTimeout] seq=%d: counterparty time %d < timeout %d, skipping",
			packet.Packet.Sequence, counterpartyTime, packet.Packet.TimeoutTimestamp)
		return false
	}

	calldata, err := CosmosNonMembership(ctx, *packet.Packet, packet.Packet.DestinationClient, []byte{2}, latestLightBlock)
	if err != nil {
		log.Printf("[EthTimeout] seq=%d: %v", packet.Packet.Sequence, err)
		return false
	}

	msgTimeoutPacket := contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket{
		Packet:           ToEthPacket(*packet.Packet),
		NonMembershipMsg: calldata,
	}

	if err := s.worker.TxHandler.SendEthTx(stdCtx, ctx, msgTimeoutPacket); err != nil {
		log.Printf("[EthTimeout] seq=%d: SendEthTx failed: %v", packet.Packet.Sequence, err)
		return false
	}
	log.Printf("[EthTimeout] seq=%d: relay completed", packet.Packet.Sequence)
	return true
}

const pendingTrackerMaxAge = 1 * time.Hour

func (s *Services) scanForEthTimeouts(stdCtx context.Context, ctx Context) {
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
		if s.timeoutEthSend(stdCtx, ctx, EthPacket{
			Type:        EthSend,
			Packet:      &packet,
			BlockNumber: info.BlockNumber,
		}) {
			s.BatchBuilder.EthPendingTracker.Remove(info.Packet.SourceClient, info.Packet.Sequence)
		}
	}
}

func (s *Services) scanForCosmosTimeouts(stdCtx context.Context, ctx Context) {
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

	log.Printf("[CosmosTimeoutScan] Found %d head-expired packets, checking ETH receipts", len(expired))
	unreceived := make([]pendingPacketInfo, 0, len(expired))
	for _, info := range expired {
		received, err := HasEthPacketReceipt(ctx, info.Packet)
		if err != nil {
			log.Printf("[CosmosTimeoutScan] seq=%d: failed to check ETH packet receipt: %v", info.Packet.Sequence, err)
			continue
		}
		if received {
			s.BatchBuilder.PendingTracker.Remove(info.Packet.SourceClient, info.Packet.Sequence)
			log.Printf("[CosmosTimeoutScan] seq=%d: ETH receipt already exists, removed from pending tracker", info.Packet.Sequence)
			continue
		}
		unreceived = append(unreceived, info)
	}
	expired = unreceived
	if len(expired) == 0 {
		return
	}

	log.Printf("[CosmosTimeoutScan] Building proof state for %d unreceived expired packets", len(expired))

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
		s.worker.WaitForCosmosCatchUp(stdCtx, ctx, updateResult.EthClientState, updateResult.SigSlot)
	}

	if err := s.worker.TxHandler.SendCosmosTxBatch(stdCtx, ctx, batchMsgs); err != nil {
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

func CosmosMembership(ctx Context, packet channeltypesv2.Packet, clientID string, pathType []byte, latestLightBlock *client.LightBlock) ([]byte, error) {
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

func CosmosNonMembership(ctx Context, packet channeltypesv2.Packet, clientID string, pathType []byte, latestLightBlock *client.LightBlock) ([]byte, error) {
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
