package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	tendermintContract "operator/bindings/SP1ICS07Tendermint"
	updateclientContract "operator/bindings/UpdateClient"
	operatorclient "operator/client"
	"operator/prover"
	"strconv"
	"strings"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const ICS26_IBC_STORAGE_SLOT = "0x1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600"

type Worker struct {
	TxHandler TransactionHandler
	Prover    Prover
}

func NewWorker(txHandler TransactionHandler, prover Prover) *Worker {
	return &Worker{
		txHandler,
		prover,
	}
}

func (w *Worker) CreateCosmosClient(ctx Context, proofType string, trustingPeriod uint32, trustedBlock int64, trustLevel string) error {
	genesis, err := operatorclient.GetGenesis(ctx.CosmosClient(), trustedBlock, trustingPeriod, trustLevel, proofType)
	if err != nil {
		return fmt.Errorf("failed to get genesis: %w", err)
	}

	clientState := genesis.TrustedClientState
	consensusState := genesis.TrustedConsensusState

	log.Printf("[CreateCosmosClient] clientState: chainId=%s trustLevel=%d/%d height=%d/%d trustingPeriod=%d unbondingPeriod=%d isFrozen=%v zkAlgorithm=%d",
		clientState.ChainId, clientState.TrustLevel.Numerator, clientState.TrustLevel.Denominator,
		clientState.LatestHeight.RevisionNumber, clientState.LatestHeight.RevisionHeight,
		clientState.TrustingPeriod, clientState.UnbondingPeriod, clientState.IsFrozen, clientState.ZkAlgorithm)

	clientStateEncoded, err := operatorclient.EncodeClientState(clientState)
	if err != nil {
		return fmt.Errorf("failed to encode client state: %w", err)
	}
	log.Printf("[CreateCosmosClient] clientStateEncoded len=%d hex=%x", len(clientStateEncoded), clientStateEncoded[:min(64, len(clientStateEncoded))])

	consensusStateEncoded, err := operatorclient.EncodeConsensusState(consensusState)
	if err != nil {
		return fmt.Errorf("failed to encode consensus state: %w", err)
	}

	consensusHash := crypto.Keccak256(consensusStateEncoded)
	log.Printf("[CreateCosmosClient] consensusHash=%x", consensusHash)
	return w.TxHandler.CreateCosmosClientContract(ctx, clientStateEncoded, consensusHash)
}

func (w *Worker) UpdateCosmosClient(ctx Context, proofType string, trustedBlock int64, trustLevel string) (*operatorclient.LightBlock, error) {
	status, err := ctx.CosmosClient().Status(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	log.Printf("[UpdateCosmosClient] called with trustedBlock=%d, latestBlockHeight=%d", trustedBlock, status.SyncInfo.LatestBlockHeight)

	if trustedBlock == 0 {
		// First update: query on-chain client state for the initial trusted height
		ics07, err := tendermintContract.NewContractSP1ICS07Tendermint(*ctx.ClientContract(), ctx.EthClient())
		if err != nil {
			return nil, fmt.Errorf("failed to create ICS07 instance: %w", err)
		}
		log.Printf("[UpdateCosmosClient] Querying on-chain client state at ICS07=%s", ctx.ClientContract().Hex())
		clientStateBytes, err := ics07.GetClientState(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get on-chain client state: %w", err)
		}
		log.Printf("[UpdateCosmosClient] GetClientState returned %d bytes: %x", len(clientStateBytes), clientStateBytes[:min(64, len(clientStateBytes))])
		onChainClientState, err := operatorclient.DecodeClientState(clientStateBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to decode on-chain client state: %w", err)
		}
		log.Printf("[UpdateCosmosClient] On-chain client state: chainId=%s height=(%d,%d) frozen=%v",
			onChainClientState.ChainId, onChainClientState.LatestHeight.RevisionNumber, onChainClientState.LatestHeight.RevisionHeight, onChainClientState.IsFrozen)
		trustedBlock = int64(onChainClientState.LatestHeight.RevisionHeight)
		log.Printf("[UpdateCosmosClient] Using on-chain client height %d as trusted block", trustedBlock)
	}
	if trustedBlock >= status.SyncInfo.LatestBlockHeight {
		return nil, fmt.Errorf("client is up to date (trusted=%d, latest=%d)", trustedBlock, status.SyncInfo.LatestBlockHeight)
	}

	log.Printf("[UpdateCosmosClient] Fetching trustedLightBlock at height %d, latestLightBlock at height %d", trustedBlock, status.SyncInfo.LatestBlockHeight)
	trustedLightBlock, err := operatorclient.GetLightBlock(ctx.CosmosClient(), trustedBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted light block: %w", err)
	}

	latestLightBlock, err := operatorclient.GetLightBlock(ctx.CosmosClient(), status.SyncInfo.LatestBlockHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest light block: %w", err)
	}

	unbondingPeriod, err := operatorclient.GetUnbondingTime(ctx.CosmosClient())
	if err != nil {
		return nil, fmt.Errorf("failed to get unbonding time: %w", err)
	}

	trustingPeriod := uint32(unbondingPeriod * 2 / 3)

	if trustingPeriod > uint32(unbondingPeriod) {
		return nil, fmt.Errorf("trusting period %d cannot be greater than unbonding period %d", trustingPeriod, uint32(unbondingPeriod))
	}

	chainId := trustedLightBlock.SignedHeader.Header.ChainID
	revision := clienttypes.ParseChainID(chainId)

	trustThreshold, err := operatorclient.ParseTrustThreshold(trustLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trust level: %w", err)
	}

	var zkAlgorithm operatorclient.SupportedZkAlgorithm
	switch proofType {
	case "groth16":
		zkAlgorithm = operatorclient.Groth16
	case "plonk":
		zkAlgorithm = operatorclient.Plonk
	default:
		return nil, fmt.Errorf("unsupported proof type: %s, supported types are: groth16, plonk", proofType)
	}

	clientState := updateclientContract.IICS07TendermintMsgsClientState{
		ChainId:    chainId,
		TrustLevel: trustThreshold,
		LatestHeight: updateclientContract.IICS02ClientMsgsHeight{
			RevisionNumber: revision,
			RevisionHeight: uint64(trustedLightBlock.SignedHeader.Header.Height),
		},
		IsFrozen:        false,
		ZkAlgorithm:     uint8(zkAlgorithm),
		TrustingPeriod:  trustingPeriod,
		UnbondingPeriod: uint32(unbondingPeriod),
	}

	consensusState := updateclientContract.IICS07TendermintMsgsConsensusState{
		Timestamp:          big.NewInt(trustedLightBlock.SignedHeader.Header.Time.UnixNano()),
		Root:               bytesToBytes32(trustedLightBlock.SignedHeader.Header.AppHash),
		NextValidatorsHash: bytesToBytes32(trustedLightBlock.SignedHeader.NextValidatorsHash),
	}

	// Debug: log consensus state fields and hash for comparison with on-chain
	csEncoded, _ := operatorclient.EncodeConsensusState(consensusState)
	csHash := crypto.Keccak256(csEncoded)
	log.Printf("[UpdateCosmosClient] trustedConsensusState: timestamp=%s root=%x nextValHash=%x",
		consensusState.Timestamp.String(), consensusState.Root, consensusState.NextValidatorsHash)
	log.Printf("[UpdateCosmosClient] csEncoded len=%d hex=%x", len(csEncoded), csEncoded)
	log.Printf("[UpdateCosmosClient] csHash=%x", csHash)

	proposedHeader := latestLightBlock.IntoHeader(*trustedLightBlock)

	log.Printf("[UpdateCosmosClient] clientState.LatestHeight=(%d,%d) proposedHeader.Height=%d trustedBlock=%d latestBlock=%d",
		clientState.LatestHeight.RevisionNumber, clientState.LatestHeight.RevisionHeight,
		proposedHeader.SignedHeader.Header.Height, trustedLightBlock.BlockHeight, latestLightBlock.BlockHeight)

	// TODO: proof for multiple sigs — currently only proves 1st valid validator signature
	// Extract first non-absent validator signature from the latest block
	valSig, err := prover.ExtractValidatorSignature(latestLightBlock, chainId)
	if err != nil {
		return nil, fmt.Errorf("extract validator signature: %w", err)
	}
	log.Printf("[UpdateCosmosClient] Generating Groth16 proof for validator signature...")
	proof, commitments, commitmentPok, err := w.Prover.GenerateProof(valSig.Signature, valSig.PublicKey, valSig.SignBytes)
	if err != nil {
		return nil, fmt.Errorf("error generating proof: %w", err)
	}
	log.Printf("[UpdateCosmosClient] Proof generated successfully. Sending Eth tx...")
	msg := updateclientContract.IUpdateClientMsgsMsgUpdateClient{
		ClientState:           clientState,
		TrustedConsensusState: consensusState,
		Time:                  big.NewInt(time.Now().UnixNano()),
		ProposedHeader:        proposedHeader,
		Proof:                 proof,
		Commitments:           commitments,
		CommitmentPok:         commitmentPok,
	}

	err = w.TxHandler.SendEthTx(ctx, msg)
	if err != nil {
		log.Printf("[UpdateCosmosClient] SendEthTx failed: %v", err)
		return nil, err
	}
	log.Printf("[UpdateCosmosClient] SendEthTx succeeded")
	return latestLightBlock, nil
}

func (w *Worker) CreateEthClient(ctx Context, checksum string) error {
	beaconAPIURL := ctx.BeaconAPIURL()
	if beaconAPIURL == "" {
		return fmt.Errorf("beacon API URL is not configured")
	}

	genesis, err := operatorclient.GetBeaconGenesis(beaconAPIURL)
	if err != nil {
		return fmt.Errorf("failed to get light client genesis: %w", err)
	}

	spec, err := operatorclient.GetBeaconSpec(beaconAPIURL)
	if err != nil {
		return fmt.Errorf("failed to get light client spec: %w", err)
	}
	beaconBlock, err := operatorclient.GetBeaconBlock(beaconAPIURL, "finalized")
	if err != nil {
		return fmt.Errorf("failed to get beacon block: %w", err)
	}

	blockRoot, err := operatorclient.GetBeaconBlockRoot(beaconAPIURL, beaconBlock.Message.Slot)
	if err != nil {
		return fmt.Errorf("failed to get beacon block root: %w", err)
	}

	bootstrap, err := operatorclient.GetLightClientBootstrap(beaconAPIURL, blockRoot)
	if err != nil {
		return fmt.Errorf("failed to get light client bootstrap: %w", err)
	}

	if bootstrap.Data.Header.Execution.BlockNumber != beaconBlock.Message.Body.ExecutionPayload.BlockNumber {
		return fmt.Errorf("Light client bootstrap block number does not match execution block number")
	}

	chainId, err := ctx.EthClient().ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get eth chain id: %w", err)
	}

	forkParameters, err := spec.ToForkParameters()
	if err != nil {
		return fmt.Errorf("failed to get fork parameters: %w", err)
	}

	epochsPerSyncCommitteePeriod, err := strconv.ParseUint(spec.EpochsPerSyncCommitteePeriod, 10, 64)
	if err != nil {
		return err
	}

	genesisTime, err := strconv.ParseUint(genesis.GenesisTime, 10, 64)
	if err != nil {
		return err
	}

	blockNumber, err := strconv.ParseUint(bootstrap.Data.Header.Execution.BlockNumber, 10, 64)
	if err != nil {
		return err
	}
	slot, err := strconv.ParseUint(bootstrap.Data.Header.Beacon.Slot, 10, 64)
	if err != nil {
		return err
	}

	syncCommitteeSize, err := strconv.ParseUint(spec.SyncCommitteeSize, 10, 64)
	if err != nil {
		return err
	}

	secondsPerSlot, err := strconv.ParseUint(spec.SecondsPerSlot, 10, 64)
	if err != nil {
		return err
	}
	slotsPerEpoch, err := strconv.ParseUint(spec.SlotsPerEpoch, 10, 64)
	if err != nil {
		return err
	}
	clientState := operatorclient.EthereumClientState{
		ChainID:                      chainId.Uint64(),
		EpochsPerSyncCommitteePeriod: epochsPerSyncCommitteePeriod,
		ForkParameters:               *forkParameters,
		GenesisSlot:                  0,
		GenesisTime:                  genesisTime,
		GenesisValidatorsRoot:        genesis.GenesisValidatorsRoot,
		IbcCommitmentSlot:            ICS26_IBC_STORAGE_SLOT,
		IbcContractAddress:           ctx.RouterContract().String(),
		IsFrozen:                     false,
		LatestExecutionBlockNumber:   blockNumber,
		LatestSlot:                   slot,
		MinSyncCommitteeParticipants: (syncCommitteeSize + 2) / 3,
		SecondsPerSlot:               secondsPerSlot,
		SlotsPerEpoch:                slotsPerEpoch,
		SyncCommitteeSize:            syncCommitteeSize,
	}
	clientStateBz, err := json.Marshal(clientState)
	if err != nil {
		return fmt.Errorf("error serializing client state: %w", err)
	}
	checksumTrimmed := strings.TrimPrefix(checksum, "0x")
	checksumBz, err := hex.DecodeString(checksumTrimmed)
	if err != nil {
		return fmt.Errorf("error parsing checksum: %w", err)
	}
	wasmClientState := ibcwasmtypes.ClientState{
		Data:     clientStateBz,
		Checksum: checksumBz,
		LatestHeight: clienttypes.Height{
			RevisionNumber: 0,
			RevisionHeight: clientState.LatestSlot,
		},
	}

	timestamp, err := strconv.ParseUint(bootstrap.Data.Header.Execution.Timestamp, 10, 64)
	if err != nil {
		return err
	}

	currentSyncCommittee, err := bootstrap.Data.CurrentSyncCommittee.ToSummarizedSyncCommittee()
	if err != nil {
		return fmt.Errorf("failed to hash pubkeys for CurrentSyncCommittee: %w", err)
	}

	latestPeriod := clientState.ComputeSyncCommitteePeriodAtSlot(clientState.LatestSlot)
	lightClientUpdates, err := operatorclient.GetLightClientUpdates(ctx.BeaconAPIURL(), latestPeriod, 1)
	if err != nil {
		return fmt.Errorf("failed to get light client updates: %w", err)
	}

	nextSyncCommittee, err := lightClientUpdates[0].NextSyncCommittee.ToSummarizedSyncCommittee()
	if err != nil {
		return fmt.Errorf("failed to hash pubkeys for NextSyncCommittee: %w", err)
	}

	consensusState := operatorclient.EthereumConsensusState{
		Slot:                 clientState.LatestSlot,
		StateRoot:            bootstrap.Data.Header.Execution.StateRoot,
		Timestamp:            timestamp,
		CurrentSyncCommittee: *currentSyncCommittee,
		NextSyncCommittee:    nextSyncCommittee,
		StorageRoot:          "0x0000000000000000000000000000000000000000000000000000000000000000",
	}
	fmt.Println("consensusState: ", consensusState)
	consensusStateBz, err := json.Marshal(consensusState)
	if err != nil {
		return fmt.Errorf("error serializing consensus state: %w", err)
	}

	wasmConsensusState := ibcwasmtypes.ConsensusState{
		Data: consensusStateBz,
	}

	return w.TxHandler.CreateEthClient(ctx, &wasmClientState, &wasmConsensusState)
}
func (w *Worker) UpdateEthClient(ctx Context) error {
	beaconAPIURL := ctx.BeaconAPIURL()
	if beaconAPIURL == "" {
		return fmt.Errorf("beacon API URL is not configured")
	}

	ethClientID := ctx.EthClientID()
	if ethClientID == "" {
		return fmt.Errorf("ethereum client ID is not configured")
	}

	ethClientState, err := operatorclient.GetEthereumClientState(ctx.CosmosClient(), ethClientID)
	if err != nil {
		return fmt.Errorf("failed to get ethereum client state: %w", err)
	}
	trustedSlot := ethClientState.LatestSlot

	finalityUpdate, err := operatorclient.GetFinalityUpdate(beaconAPIURL)
	if err != nil {
		return fmt.Errorf("failed to get finality update: %w", err)
	}

	finalizedSlot, err := parseSlot(finalityUpdate.FinalizedHeader.Beacon.Slot)
	if err != nil {
		return fmt.Errorf("failed to parse finalized slot: %w", err)
	}

	if finalizedSlot <= trustedSlot {
		return nil
	}

	trustedPeriod := ethClientState.ComputeSyncCommitteePeriodAtSlot(trustedSlot)
	targetPeriod := ethClientState.ComputeSyncCommitteePeriodAtSlot(finalizedSlot)

	if targetPeriod > trustedPeriod {
		return w.updateEthClientWithPeriodCrossing(ctx, beaconAPIURL, ethClientID, ethClientState, trustedSlot, trustedPeriod, targetPeriod, finalityUpdate, finalizedSlot)
	}

	return w.updateEthClientSamePeriod(ctx, beaconAPIURL, ethClientID, trustedSlot, finalityUpdate, finalizedSlot)
}

func (w *Worker) updateEthClientSamePeriod(ctx Context, beaconAPIURL, ethClientID string, trustedSlot uint64, finalityUpdate *operatorclient.LightClientFinalityUpdate, _ uint64) error {
	attestedSlot := finalityUpdate.AttestedHeader.Beacon.Slot
	blockRoot, err := operatorclient.GetBeaconBlockRoot(beaconAPIURL, attestedSlot)
	if err != nil {
		return fmt.Errorf("failed to get beacon block root: %w", err)
	}

	bootstrap, err := operatorclient.GetLightClientBootstrap(beaconAPIURL, blockRoot)
	if err != nil {
		return fmt.Errorf("failed to get light client bootstrap: %w", err)
	}

	syncCommittee := bootstrap.Data.CurrentSyncCommittee

	consensusUpdate := operatorclient.LightClientUpdate{
		AttestedHeader:          finalityUpdate.AttestedHeader,
		NextSyncCommittee:       nil,
		NextSyncCommitteeBranch: nil,
		FinalizedHeader:         finalityUpdate.FinalizedHeader,
		FinalityBranch:          finalityUpdate.FinalityBranch,
		SyncAggregate:           finalityUpdate.SyncAggregate,
		SignatureSlot:           finalityUpdate.SignatureSlot,
	}

	header := operatorclient.EthereumHeader{
		ActiveSyncCommittee: operatorclient.ActiveSyncCommittee{
			Current: &syncCommittee,
		},
		ConsensusUpdate: consensusUpdate,
		TrustedSlot:     trustedSlot,
	}

	msg, err := buildMsgUpdateClient(ethClientID, header)
	if err != nil {
		return fmt.Errorf("failed to build update client message: %w", err)
	}

	return w.TxHandler.SendCosmosTx(ctx, msg)
}

func (w *Worker) updateEthClientWithPeriodCrossing(ctx Context, beaconAPIURL, ethClientID string, ethClientState *operatorclient.EthereumClientState, trustedSlot, trustedPeriod, targetPeriod uint64, finalityUpdate *operatorclient.LightClientFinalityUpdate, finalizedSlot uint64) error {
	count := targetPeriod - trustedPeriod + 1
	lightClientUpdates, err := operatorclient.GetLightClientUpdates(beaconAPIURL, trustedPeriod, count)
	if err != nil {
		return fmt.Errorf("failed to get light client updates: %w", err)
	}

	if len(lightClientUpdates) == 0 {
		return fmt.Errorf("no light client updates available for period range %d to %d", trustedPeriod, targetPeriod)
	}

	var msgs []any
	latestTrustedSlot := trustedSlot
	latestPeriod := trustedPeriod

	for _, update := range lightClientUpdates {
		updateFinalizedSlot, err := parseSlot(update.FinalizedHeader.Beacon.Slot)
		if err != nil {
			return fmt.Errorf("failed to parse update finalized slot: %w", err)
		}

		if updateFinalizedSlot <= latestTrustedSlot {
			continue
		}

		updatePeriod := ethClientState.ComputeSyncCommitteePeriodAtSlot(updateFinalizedSlot)
		if updatePeriod == latestPeriod {
			continue
		}

		blockRoot, err := operatorclient.GetBeaconBlockRoot(beaconAPIURL, fmt.Sprintf("%d", updateFinalizedSlot))
		if err != nil {
			return fmt.Errorf("failed to get beacon block root: %w", err)
		}

		bootstrap, err := operatorclient.GetLightClientBootstrap(beaconAPIURL, blockRoot)
		if err != nil {
			return fmt.Errorf("failed to get light client bootstrap: %w", err)
		}

		syncCommittee := bootstrap.Data.CurrentSyncCommittee

		header := operatorclient.EthereumHeader{
			ActiveSyncCommittee: operatorclient.ActiveSyncCommittee{
				Next: &syncCommittee,
			},
			ConsensusUpdate: update,
			TrustedSlot:     latestTrustedSlot,
		}

		msg, err := buildMsgUpdateClient(ethClientID, header)
		if err != nil {
			return fmt.Errorf("failed to build update client message: %w", err)
		}

		msgs = append(msgs, msg)
		latestPeriod = updatePeriod
		latestTrustedSlot = updateFinalizedSlot
	}

	if finalizedSlot > latestTrustedSlot {
		attestedSlot := finalityUpdate.AttestedHeader.Beacon.Slot
		blockRoot, err := operatorclient.GetBeaconBlockRoot(beaconAPIURL, attestedSlot)
		if err != nil {
			return fmt.Errorf("failed to get beacon block root: %w", err)
		}

		bootstrap, err := operatorclient.GetLightClientBootstrap(beaconAPIURL, blockRoot)
		if err != nil {
			return fmt.Errorf("failed to get light client bootstrap: %w", err)
		}

		syncCommittee := bootstrap.Data.CurrentSyncCommittee

		consensusUpdate := operatorclient.LightClientUpdate{
			AttestedHeader:          finalityUpdate.AttestedHeader,
			NextSyncCommittee:       nil,
			NextSyncCommitteeBranch: nil,
			FinalizedHeader:         finalityUpdate.FinalizedHeader,
			FinalityBranch:          finalityUpdate.FinalityBranch,
			SyncAggregate:           finalityUpdate.SyncAggregate,
			SignatureSlot:           finalityUpdate.SignatureSlot,
		}

		header := operatorclient.EthereumHeader{
			ActiveSyncCommittee: operatorclient.ActiveSyncCommittee{
				Current: &syncCommittee,
			},
			ConsensusUpdate: consensusUpdate,
			TrustedSlot:     latestTrustedSlot,
		}

		msg, err := buildMsgUpdateClient(ethClientID, header)
		if err != nil {
			return fmt.Errorf("failed to build update client message: %w", err)
		}

		msgs = append(msgs, msg)
	}

	if len(msgs) == 0 {
		return nil
	}

	return w.TxHandler.SendCosmosTxBatch(ctx, msgs)
}

// parseSlot parses a slot string to uint64
func parseSlot(slotStr string) (uint64, error) {
	var slot uint64
	_, err := fmt.Sscanf(slotStr, "%d", &slot)
	return slot, err
}

func bytesToBytes32(data []byte) [32]byte {
	var result [32]byte
	copy(result[:], data)
	return result
}

// buildMsgUpdateClient builds a MsgUpdateClient for the Ethereum light client
func buildMsgUpdateClient(clientID string, header operatorclient.EthereumHeader) (*clienttypes.MsgUpdateClient, error) {
	// Serialize the header to JSON
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal header: %w", err)
	}

	// Create the wasm ClientMessage
	clientMessage := &ibcwasmtypes.ClientMessage{
		Data: headerBytes,
	}

	// Pack into Any type
	clientMessageAny, err := codectypes.NewAnyWithValue(clientMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to create Any for client message: %w", err)
	}

	// Build the MsgUpdateClient
	msg := &clienttypes.MsgUpdateClient{
		ClientId:      clientID,
		ClientMessage: clientMessageAny,
		Signer:        "", // Will be set by the transaction handler
	}

	return msg, nil
}
