package client

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/cometbft/cometbft/crypto/merkle"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/gogoproto/proto"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
)

type EthereumClientState struct {
	ChainID                      uint64         `json:"chain_id"`
	EpochsPerSyncCommitteePeriod uint64         `json:"epochs_per_sync_committee_period"`
	ForkParameters               ForkParameters `json:"fork_parameters"`
	GenesisSlot                  uint64         `json:"genesis_slot"`
	GenesisTime                  uint64         `json:"genesis_time"`
	GenesisValidatorsRoot        string         `json:"genesis_validators_root"`
	IbcCommitmentSlot            string         `json:"ibc_commitment_slot"`
	IbcContractAddress           string         `json:"ibc_contract_address"`
	IsFrozen                     bool           `json:"is_frozen"`
	LatestExecutionBlockNumber   uint64         `json:"latest_execution_block_number"`
	LatestSlot                   uint64         `json:"latest_slot"`
	MinSyncCommitteeParticipants uint64         `json:"min_sync_committee_participants"`
	SecondsPerSlot               uint64         `json:"seconds_per_slot"`
	SlotsPerEpoch                uint64         `json:"slots_per_epoch"`
	SyncCommitteeSize            uint64         `json:"sync_committee_size"`
}

type EthereumConsensusState struct {
	Slot                 uint64                   `json:"slot"`
	StateRoot            string                   `json:"state_root"`
	Timestamp            uint64                   `json:"timestamp"`
	CurrentSyncCommittee SummarizedSyncCommittee  `json:"current_sync_committee"`
	NextSyncCommittee    *SummarizedSyncCommittee `json:"next_sync_committee"`
	StorageRoot          string                   `json:"storage_root"`
}

type ForkParameters struct {
	Altair             Fork   `json:"altair"`
	Bellatrix          Fork   `json:"bellatrix"`
	Capella            Fork   `json:"capella"`
	Deneb              Fork   `json:"deneb"`
	Electra            Fork   `json:"electra"`
	GenesisForkVersion string `json:"genesis_fork_version"`
	GenesisSlot        uint64 `json:"genesis_slot"`
}

type Fork struct {
	Epoch   uint64 `json:"epoch"`
	Version string `json:"version"`
}

func (cs *EthereumClientState) ComputeSyncCommitteePeriodAtSlot(slot uint64) uint64 {
	epoch := slot / cs.SlotsPerEpoch
	return epoch / cs.EpochsPerSyncCommitteePeriod
}

type SyncCommittee struct {
	Pubkeys         []string `json:"pubkeys"`
	AggregatePubkey string   `json:"aggregate_pubkey"`
}

func (sc *SyncCommittee) ToSummarizedSyncCommittee() (*SummarizedSyncCommittee, error) {
	pks := [][]byte{}
	for _, pk := range sc.Pubkeys {
		pkTrimmed := strings.TrimPrefix(pk, "0x")
		pkBytes, err := hex.DecodeString(pkTrimmed)
		if err != nil {
			return nil, err
		}
		pks = append(pks, pkBytes)
	}

	return &SummarizedSyncCommittee{
		PubkeysHash:     hex.EncodeToString(merkle.HashFromByteSlices(pks)),
		AggregatePubkey: sc.AggregatePubkey,
	}, nil
}

type SummarizedSyncCommittee struct {
	PubkeysHash     string `json:"pubkeys_hash"`
	AggregatePubkey string `json:"aggregate_pubkey"`
}

type LightClientHeader struct {
	Beacon          BeaconBlockHeader      `json:"beacon"`
	Execution       ExecutionPayloadHeader `json:"execution"`
	ExecutionBranch []string               `json:"execution_branch"`
}

type BeaconBlockHeader struct {
	Slot          string `json:"slot"`
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
}

type ExecutionPayloadHeader struct {
	ParentHash       string `json:"parent_hash"`
	FeeRecipient     string `json:"fee_recipient"`
	StateRoot        string `json:"state_root"`
	ReceiptsRoot     string `json:"receipts_root"`
	LogsBloom        string `json:"logs_bloom"`
	PrevRandao       string `json:"prev_randao"`
	BlockNumber      string `json:"block_number"`
	GasLimit         string `json:"gas_limit"`
	GasUsed          string `json:"gas_used"`
	Timestamp        string `json:"timestamp"`
	ExtraData        string `json:"extra_data"`
	BaseFeePerGas    string `json:"base_fee_per_gas"`
	BlockHash        string `json:"block_hash"`
	TransactionsRoot string `json:"transactions_root"`
	WithdrawalsRoot  string `json:"withdrawals_root"`
	BlobGasUsed      string `json:"blob_gas_used"`
	ExcessBlobGas    string `json:"excess_blob_gas"`
}

type SyncAggregate struct {
	SyncCommitteeBits      string `json:"sync_committee_bits"`
	SyncCommitteeSignature string `json:"sync_committee_signature"`
}

type LightClientUpdate struct {
	AttestedHeader          LightClientHeader `json:"attested_header"`
	NextSyncCommittee       *SyncCommittee    `json:"next_sync_committee"`
	NextSyncCommitteeBranch []string          `json:"next_sync_committee_branch"`
	FinalizedHeader         LightClientHeader `json:"finalized_header"`
	FinalityBranch          []string          `json:"finality_branch"`
	SyncAggregate           SyncAggregate     `json:"sync_aggregate"`
	SignatureSlot           string            `json:"signature_slot"`
}

type LightClientFinalityUpdate struct {
	AttestedHeader  LightClientHeader `json:"attested_header"`
	FinalizedHeader LightClientHeader `json:"finalized_header"`
	FinalityBranch  []string          `json:"finality_branch"`
	SyncAggregate   SyncAggregate     `json:"sync_aggregate"`
	SignatureSlot   string            `json:"signature_slot"`
}

type ActiveSyncCommittee struct {
	Current *SyncCommittee `json:"Current,omitempty"`
	Next    *SyncCommittee `json:"Next,omitempty"`
}

type ExecutionPayload struct {
	ParentHash    string `json:"parent_hash"`
	FeeRecipient  string `json:"fee_recipient"`
	StateRoot     string `json:"state_root"`
	ReceiptsRoot  string `json:"receipts_root"`
	LogsBloom     string `json:"logs_bloom"`
	PrevRandao    string `json:"prev_randao"`
	BlockNumber   string `json:"block_number"`
	GasLimit      string `json:"gas_limit"`
	GasUsed       string `json:"gas_used"`
	Timestamp     string `json:"timestamp"`
	ExtraData     string `json:"extra_data"`
	BaseFeePerGas string `json:"base_fee_per_gas"`
	BlockHash     string `json:"block_hash"`
	BlobGasUsed   string `json:"blob_gas_used"`
	ExcessBlobGas string `json:"excess_blob_gas"`
}

type BeaconBlockBody struct {
	SyncAggregate    SyncAggregate    `json:"sync_aggregate"`
	ExecutionPayload ExecutionPayload `json:"execution_payload"`
}

type BeaconBlockMessage struct {
	Slot          string          `json:"slot"`
	ProposerIndex string          `json:"proposer_index"`
	ParentRoot    string          `json:"parent_root"`
	StateRoot     string          `json:"state_root"`
	Body          BeaconBlockBody `json:"body"`
}
type BeaconGenesis struct {
	GenesisTime           string `json:"genesis_time"`
	GenesisValidatorsRoot string `json:"genesis_validators_root"`
	GenesisForkstring     string `json:"genesis_fork_string"`
}
type BeaconSpec struct {
	// Identity
	ConfigName string `json:"CONFIG_NAME"`
	PresetBase string `json:"PRESET_BASE"`

	// Genesis
	MinGenesisActiveValidatorCount string `json:"MIN_GENESIS_ACTIVE_VALIDATOR_COUNT"`
	MinGenesisTime                 string `json:"MIN_GENESIS_TIME"`
	GenesisDelay                   string `json:"GENESIS_DELAY"`

	// Fork Versions & Epochs
	GenesisForkVersion   string `json:"GENESIS_FORK_VERSION"`
	AltairForkVersion    string `json:"ALTAIR_FORK_VERSION"`
	AltairForkEpoch      string `json:"ALTAIR_FORK_EPOCH"`
	BellatrixForkVersion string `json:"BELLATRIX_FORK_VERSION"`
	BellatrixForkEpoch   string `json:"BELLATRIX_FORK_EPOCH"`
	CapellaForkVersion   string `json:"CAPELLA_FORK_VERSION"`
	CapellaForkEpoch     string `json:"CAPELLA_FORK_EPOCH"`
	DenebForkVersion     string `json:"DENEB_FORK_VERSION"`
	DenebForkEpoch       string `json:"DENEB_FORK_EPOCH"`
	ElectraForkVersion   string `json:"ELECTRA_FORK_VERSION"`
	ElectraForkEpoch     string `json:"ELECTRA_FORK_EPOCH"`
	FuluForkVersion      string `json:"FULU_FORK_VERSION"`
	FuluForkEpoch        string `json:"FULU_FORK_EPOCH"`

	// Time
	SecondsPerSlot      string `json:"SECONDS_PER_SLOT"`
	SecondsPerEth1Block string `json:"SECONDS_PER_ETH1_BLOCK"`

	// Deposit Contract
	DepositChainID         string `json:"DEPOSIT_CHAIN_ID"`
	DepositNetworkID       string `json:"DEPOSIT_NETWORK_ID"`
	DepositContractAddress string `json:"DEPOSIT_CONTRACT_ADDRESS"`

	// Validator / Committee
	SlotsPerEpoch                    string `json:"SLOTS_PER_EPOCH"`
	MaxCommitteesPerSlot             string `json:"MAX_COMMITTEES_PER_SLOT"`
	TargetCommitteeSize              string `json:"TARGET_COMMITTEE_SIZE"`
	MaxValidatorsPerCommittee        string `json:"MAX_VALIDATORS_PER_COMMITTEE"`
	ShuffleRoundCount                string `json:"SHUFFLE_ROUND_COUNT"`
	MinSeedLookahead                 string `json:"MIN_SEED_LOOKAHEAD"`
	MaxSeedLookahead                 string `json:"MAX_SEED_LOOKAHEAD"`
	EpochsPerEth1VotingPeriod        string `json:"EPOCHS_PER_ETH1_VOTING_PERIOD"`
	SlotsPerHistoricalRoot           string `json:"SLOTS_PER_HISTORICAL_ROOT"`
	MinEpochsToInactivityPenalty     string `json:"MIN_EPOCHS_TO_INACTIVITY_PENALTY"`
	EpochsPerHistoricalVector        string `json:"EPOCHS_PER_HISTORICAL_VECTOR"`
	EpochsPerSlashingsVector         string `json:"EPOCHS_PER_SLASHINGS_VECTOR"`
	HistoricalRootsLimit             string `json:"HISTORICAL_ROOTS_LIMIT"`
	ValidatorRegistryLimit           string `json:"VALIDATOR_REGISTRY_LIMIT"`
	MinValidatorWithdrawabilityDelay string `json:"MIN_VALIDATOR_WITHDRAWABILITY_DELAY"`
	ShardCommitteePeriod             string `json:"SHARD_COMMITTEE_PERIOD"`
	EjectionBalance                  string `json:"EJECTION_BALANCE"`
	MinPerEpochChurnLimit            string `json:"MIN_PER_EPOCH_CHURN_LIMIT"`
	MaxPerEpochActivationChurnLimit  string `json:"MAX_PER_EPOCH_ACTIVATION_CHURN_LIMIT"`
	ChurnLimitQuotient               string `json:"CHURN_LIMIT_QUOTIENT"`

	// Balances
	MinDepositAmount                 string `json:"MIN_DEPOSIT_AMOUNT"`
	MaxEffectiveBalance              string `json:"MAX_EFFECTIVE_BALANCE"`
	EffectiveBalanceIncrement        string `json:"EFFECTIVE_BALANCE_INCREMENT"`
	MinActivationBalance             string `json:"MIN_ACTIVATION_BALANCE"`
	MaxEffectiveBalanceElectra       string `json:"MAX_EFFECTIVE_BALANCE_ELECTRA"`
	BalancePerAdditionalCustodyGroup string `json:"BALANCE_PER_ADDITIONAL_CUSTODY_GROUP"`

	// Rewards & Penalties
	BaseRewardFactor                        string `json:"BASE_REWARD_FACTOR"`
	WhistleblowerRewardQuotient             string `json:"WHISTLEBLOWER_REWARD_QUOTIENT"`
	ProposerRewardQuotient                  string `json:"PROPOSER_REWARD_QUOTIENT"`
	InactivityPenaltyQuotient               string `json:"INACTIVITY_PENALTY_QUOTIENT"`
	MinSlashingPenaltyQuotient              string `json:"MIN_SLASHING_PENALTY_QUOTIENT"`
	ProportionalSlashingMultiplier          string `json:"PROPORTIONAL_SLASHING_MULTIPLIER"`
	InactivityPenaltyQuotientAltair         string `json:"INACTIVITY_PENALTY_QUOTIENT_ALTAIR"`
	MinSlashingPenaltyQuotientAltair        string `json:"MIN_SLASHING_PENALTY_QUOTIENT_ALTAIR"`
	ProportionalSlashingMultiplierAltair    string `json:"PROPORTIONAL_SLASHING_MULTIPLIER_ALTAIR"`
	InactivityPenaltyQuotientBellatrix      string `json:"INACTIVITY_PENALTY_QUOTIENT_BELLATRIX"`
	MinSlashingPenaltyQuotientBellatrix     string `json:"MIN_SLASHING_PENALTY_QUOTIENT_BELLATRIX"`
	ProportionalSlashingMultiplierBellatrix string `json:"PROPORTIONAL_SLASHING_MULTIPLIER_BELLATRIX"`
	MinSlashingPenaltyQuotientElectra       string `json:"MIN_SLASHING_PENALTY_QUOTIENT_ELECTRA"`
	WhistleblowerRewardQuotientElectra      string `json:"WHISTLEBLOWER_REWARD_QUOTIENT_ELECTRA"`
	InactivityScoreBias                     string `json:"INACTIVITY_SCORE_BIAS"`
	InactivityScoreRecoveryRate             string `json:"INACTIVITY_SCORE_RECOVERY_RATE"`

	// Sync Committee
	SyncCommitteeSize                    string `json:"SYNC_COMMITTEE_SIZE"`
	EpochsPerSyncCommitteePeriod         string `json:"EPOCHS_PER_SYNC_COMMITTEE_PERIOD"`
	MinSyncCommitteeParticipants         string `json:"MIN_SYNC_COMMITTEE_PARTICIPANTS"`
	SyncCommitteeSubnetCount             string `json:"SYNC_COMMITTEE_SUBNET_COUNT"`
	TargetAggregatorsPerSyncSubcommittee string `json:"TARGET_AGGREGATORS_PER_SYNC_SUBCOMMITTEE"`
	TargetAggregatorsPerCommittee        string `json:"TARGET_AGGREGATORS_PER_COMMITTEE"`

	// Blobs & EIP-4844 / Deneb
	MaxBlobsPerBlock                 string `json:"MAX_BLOBS_PER_BLOCK"`
	MaxBlobCommitmentsPerBlock       string `json:"MAX_BLOB_COMMITMENTS_PER_BLOCK"`
	FieldElementsPerBlob             string `json:"FIELD_ELEMENTS_PER_BLOB"`
	BlobSidecarSubnetCount           string `json:"BLOB_SIDECAR_SUBNET_COUNT"`
	MaxRequestBlobSidecars           string `json:"MAX_REQUEST_BLOB_SIDECARS"`
	MinEpochsForBlobSidecarsRequests string `json:"MIN_EPOCHS_FOR_BLOB_SIDECARS_REQUESTS"`

	// Electra Blobs
	MaxBlobsPerBlockElectra       string `json:"MAX_BLOBS_PER_BLOCK_ELECTRA"`
	BlobSidecarSubnetCountElectra string `json:"BLOB_SIDECAR_SUBNET_COUNT_ELECTRA"`
	MaxRequestBlobSidecarsElectra string `json:"MAX_REQUEST_BLOB_SIDECARS_ELECTRA"`

	// PeerDAS / Data Columns
	NumberOfColumns                        string `json:"NUMBER_OF_COLUMNS"`
	NumberOfCustodyGroups                  string `json:"NUMBER_OF_CUSTODY_GROUPS"`
	DataColumnSidecarSubnetCount           string `json:"DATA_COLUMN_SIDECAR_SUBNET_COUNT"`
	SamplesPerSlot                         string `json:"SAMPLES_PER_SLOT"`
	CustodyRequirement                     string `json:"CUSTODY_REQUIREMENT"`
	ValidatorCustodyRequirement            string `json:"VALIDATOR_CUSTODY_REQUIREMENT"`
	MaxRequestDataColumnSidecars           string `json:"MAX_REQUEST_DATA_COLUMN_SIDECARS"`
	MinEpochsForDataColumnSidecarsRequests string `json:"MIN_EPOCHS_FOR_DATA_COLUMN_SIDECARS_REQUESTS"`
	FieldElementsPerCell                   string `json:"FIELD_ELEMENTS_PER_CELL"`
	FieldElementsPerExtBlob                string `json:"FIELD_ELEMENTS_PER_EXT_BLOB"`
	KzgCommitmentsInclusionProofDepth      string `json:"KZG_COMMITMENTS_INCLUSION_PROOF_DEPTH"`

	// Electra Limits
	PendingDepositsLimit                  string `json:"PENDING_DEPOSITS_LIMIT"`
	PendingPartialWithdrawalsLimit        string `json:"PENDING_PARTIAL_WITHDRAWALS_LIMIT"`
	PendingConsolidationsLimit            string `json:"PENDING_CONSOLIDATIONS_LIMIT"`
	MaxAttesterSlashingsElectra           string `json:"MAX_ATTESTER_SLASHINGS_ELECTRA"`
	MaxAttestationsElectra                string `json:"MAX_ATTESTATIONS_ELECTRA"`
	MaxDepositRequestsPerPayload          string `json:"MAX_DEPOSIT_REQUESTS_PER_PAYLOAD"`
	MaxWithdrawalRequestsPerPayload       string `json:"MAX_WITHDRAWAL_REQUESTS_PER_PAYLOAD"`
	MaxConsolidationRequestsPerPayload    string `json:"MAX_CONSOLIDATION_REQUESTS_PER_PAYLOAD"`
	MaxPendingPartialsPerWithdrawalsSweep string `json:"MAX_PENDING_PARTIALS_PER_WITHDRAWALS_SWEEP"`
	MaxPendingDepositsPerEpoch            string `json:"MAX_PENDING_DEPOSITS_PER_EPOCH"`
	MinPerEpochChurnLimitElectra          string `json:"MIN_PER_EPOCH_CHURN_LIMIT_ELECTRA"`
	MaxPerEpochActivationExitChurnLimit   string `json:"MAX_PER_EPOCH_ACTIVATION_EXIT_CHURN_LIMIT"`
	UnsetDepositRequestsStartIndex        string `json:"UNSET_DEPOSIT_REQUESTS_START_INDEX"`
	FullExitRequestAmount                 string `json:"FULL_EXIT_REQUEST_AMOUNT"`

	// Execution / Payload
	MaxBytesPerTransaction           string `json:"MAX_BYTES_PER_TRANSACTION"`
	MaxTransactionsPerPayload        string `json:"MAX_TRANSACTIONS_PER_PAYLOAD"`
	BytesPerLogsBloom                string `json:"BYTES_PER_LOGS_BLOOM"`
	MaxExtraDataBytes                string `json:"MAX_EXTRA_DATA_BYTES"`
	MaxWithdrawalsPerPayload         string `json:"MAX_WITHDRAWALS_PER_PAYLOAD"`
	MaxValidatorsPerWithdrawalsSweep string `json:"MAX_VALIDATORS_PER_WITHDRAWALS_SWEEP"`
	MaxBlsToExecutionChanges         string `json:"MAX_BLS_TO_EXECUTION_CHANGES"`
	GasLimitAdjustmentFactor         string `json:"GAS_LIMIT_ADJUSTMENT_FACTOR"`
	MaxPayloadSize                   string `json:"MAX_PAYLOAD_SIZE"`

	// Network
	MaxRequestBlocks                  string `json:"MAX_REQUEST_BLOCKS"`
	MaxRequestBlocksDeneb             string `json:"MAX_REQUEST_BLOCKS_DENEB"`
	MinEpochsForBlockRequests         string `json:"MIN_EPOCHS_FOR_BLOCK_REQUESTS"`
	Eth1FollowDistance                string `json:"ETH1_FOLLOW_DISTANCE"`
	SubnetsPerNode                    string `json:"SUBNETS_PER_NODE"`
	AttestationSubnetPrefixBits       string `json:"ATTESTATION_SUBNET_PREFIX_BITS"`
	TtfbTimeout                       string `json:"TTFB_TIMEOUT"`
	RespTimeout                       string `json:"RESP_TIMEOUT"`
	AttestationPropagationSlotRange   string `json:"ATTESTATION_PROPAGATION_SLOT_RANGE"`
	MaximumGossipClockDisparityMillis string `json:"MAXIMUM_GOSSIP_CLOCK_DISPARITY_MILLIS"`
	ProposerScoreBoost                string `json:"PROPOSER_SCORE_BOOST"`

	// Operation Limits
	MaxProposerSlashings         string `json:"MAX_PROPOSER_SLASHINGS"`
	MaxAttesterSlashings         string `json:"MAX_ATTESTER_SLASHINGS"`
	MaxAttestations              string `json:"MAX_ATTESTATIONS"`
	MaxDeposits                  string `json:"MAX_DEPOSITS"`
	MaxVoluntaryExits            string `json:"MAX_VOLUNTARY_EXITS"`
	MinAttestationInclusionDelay string `json:"MIN_ATTESTATION_INCLUSION_DELAY"`
	HysteresisQuotient           string `json:"HYSTERESIS_QUOTIENT"`
	HysteresisDownwardMultiplier string `json:"HYSTERESIS_DOWNWARD_MULTIPLIER"`
	HysteresisUpwardMultiplier   string `json:"HYSTERESIS_UPWARD_MULTIPLIER"`

	// Domains
	DomainBeaconProposer              string `json:"DOMAIN_BEACON_PROPOSER"`
	DomainBeaconAttester              string `json:"DOMAIN_BEACON_ATTESTER"`
	DomainRandao                      string `json:"DOMAIN_RANDAO"`
	DomainDeposit                     string `json:"DOMAIN_DEPOSIT"`
	DomainVoluntaryExit               string `json:"DOMAIN_VOLUNTARY_EXIT"`
	DomainSelectionProof              string `json:"DOMAIN_SELECTION_PROOF"`
	DomainAggregateAndProof           string `json:"DOMAIN_AGGREGATE_AND_PROOF"`
	DomainSyncCommittee               string `json:"DOMAIN_SYNC_COMMITTEE"`
	DomainSyncCommitteeSelectionProof string `json:"DOMAIN_SYNC_COMMITTEE_SELECTION_PROOF"`
	DomainContributionAndProof        string `json:"DOMAIN_CONTRIBUTION_AND_PROOF"`
	DomainApplicationMask             string `json:"DOMAIN_APPLICATION_MASK"`

	// Prefixes & Versioning
	BLSWithdrawalPrefix         string `json:"BLS_WITHDRAWAL_PREFIX"`
	Eth1AddressWithdrawalPrefix string `json:"ETH1_ADDRESS_WITHDRAWAL_PREFIX"`
	CompoundingWithdrawalPrefix string `json:"COMPOUNDING_WITHDRAWAL_PREFIX"`
	VersionedHashVersionKZG     string `json:"VERSIONED_HASH_VERSION_KZG"`

	// Message Domains
	MessageDomainInvalidSnappy string `json:"MESSAGE_DOMAIN_INVALID_SNAPPY"`
	MessageDomainValidSnappy   string `json:"MESSAGE_DOMAIN_VALID_SNAPPY"`

	// Terminal / Merge
	TerminalTotalDifficulty          string `json:"TERMINAL_TOTAL_DIFFICULTY"`
	TerminalBlockHash                string `json:"TERMINAL_BLOCK_HASH"`
	TerminalBlockHashActivationEpoch string `json:"TERMINAL_BLOCK_HASH_ACTIVATION_EPOCH"`

	// Blob Schedule (Fulu)
	BlobSchedule []BlobScheduleEntry `json:"BLOB_SCHEDULE"`
}

type BlobScheduleEntry struct {
	Epoch            string `json:"EPOCH"`
	MaxBlobsPerBlock string `json:"MAX_BLOBS_PER_BLOCK"`
}
type BeaconBlock struct {
	Message   BeaconBlockMessage `json:"message"`
	Signature string             `json:"signature"`
}

type EthereumHeader struct {
	ActiveSyncCommittee ActiveSyncCommittee `json:"active_sync_committee"`
	ConsensusUpdate     LightClientUpdate   `json:"consensus_update"`
	TrustedSlot         uint64              `json:"trusted_slot"`
}

type FinalityUpdateResponse struct {
	Version string                    `json:"version"`
	Data    LightClientFinalityUpdate `json:"data"`
}

type LightClientUpdateResponse struct {
	Data LightClientUpdate `json:"data"`
}

type BeaconGenesisResponse struct {
	Data BeaconGenesis `json:"data"`
}

type BeaconSpecResponse struct {
	Data BeaconSpec `json:"data"`
}

type BeaconBlockResponse struct {
	Version             string      `json:"version"`
	ExecutionOptimistic bool        `json:"execution_optimistic"`
	Finalized           bool        `json:"finalized"`
	Data                BeaconBlock `json:"data"`
}

type BootstrapResponse struct {
	Data struct {
		Header                     LightClientHeader `json:"header"`
		CurrentSyncCommittee       SyncCommittee     `json:"current_sync_committee"`
		CurrentSyncCommitteeBranch [6]string         `json:"current_sync_committee_branch"`
	} `json:"data"`
}

type BeaconBlockRootResponse struct {
	Data struct {
		Root string `json:"root"`
	} `json:"data"`
}

func httpGet[T any](url string) (T, error) {
	var result T

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	if resp.StatusCode != 200 {
		return result, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, body)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}

	return result, nil
}

func GetFinalityUpdate(beaconAPIURL string) (*LightClientFinalityUpdate, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/light_client/finality_update", beaconAPIURL)
	response, err := httpGet[FinalityUpdateResponse](url)
	if err != nil {
		return nil, fmt.Errorf("failed to get finality update: %w", err)
	}
	return &response.Data, nil
}

func GetLightClientUpdates(beaconAPIURL string, startPeriod, count uint64) ([]LightClientUpdate, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/light_client/updates?start_period=%d&count=%d", beaconAPIURL, startPeriod, count)
	responses, err := httpGet[[]LightClientUpdateResponse](url)
	if err != nil {
		return nil, fmt.Errorf("failed to get light client updates: %w", err)
	}

	updates := make([]LightClientUpdate, len(responses))
	for i, r := range responses {
		updates[i] = r.Data
	}
	return updates, nil
}

func GetBeaconBlockRoot(beaconAPIURL, blockID string) (string, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/blocks/%s/root", beaconAPIURL, blockID)
	response, err := httpGet[BeaconBlockRootResponse](url)
	if err != nil {
		return "", fmt.Errorf("failed to get beacon block root: %w", err)
	}
	return response.Data.Root, nil
}

func GetBeaconGenesis(beaconAPIURL string) (*BeaconGenesis, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/genesis", beaconAPIURL)
	response, err := httpGet[BeaconGenesisResponse](url)
	if err != nil {
		return nil, fmt.Errorf("failed to get beacon block root: %w", err)
	}
	return &response.Data, nil
}

func GetBeaconSpec(beaconAPIURL string) (*BeaconSpec, error) {
	url := fmt.Sprintf("%s/eth/v1/config/spec", beaconAPIURL)
	response, err := httpGet[BeaconSpecResponse](url)
	if err != nil {
		return nil, fmt.Errorf("failed to get beacon block root: %w", err)
	}
	return &response.Data, nil
}

func GetBeaconBlock(beaconAPIURL string, blockId string) (*BeaconBlock, error) {
	url := fmt.Sprintf("%s/eth/v2/beacon/blocks/%s", beaconAPIURL, blockId)
	response, err := httpGet[BeaconBlockResponse](url)
	if err != nil {
		return nil, fmt.Errorf("failed to get beacon block root: %w", err)
	}
	return &response.Data, nil
}

func GetLightClientBootstrap(beaconAPIURL, blockRoot string) (*BootstrapResponse, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/light_client/bootstrap/%s", beaconAPIURL, blockRoot)
	response, err := httpGet[BootstrapResponse](url)
	if err != nil {
		return nil, fmt.Errorf("failed to get light client bootstrap: %w", err)
	}
	return &response, nil
}

func GetEthereumClientState(cosmosClient *rpchttp.HTTP, clientID string) (*EthereumClientState, error) {
	queryReq := &clienttypes.QueryClientStateRequest{
		ClientId: clientID,
	}

	reqBytes, err := proto.Marshal(queryReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query request: %w", err)
	}

	result, err := cosmosClient.ABCIQuery(context.Background(), "/ibc.core.client.v1.Query/ClientState", reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to query client state: %w", err)
	}

	if result.Response.Code != 0 {
		return nil, fmt.Errorf("query failed with code %d: %s", result.Response.Code, result.Response.Log)
	}

	var queryResp clienttypes.QueryClientStateResponse
	if err := proto.Unmarshal(result.Response.Value, &queryResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query response: %w", err)
	}

	var wasmClientState ibcwasmtypes.ClientState
	if err := proto.Unmarshal(queryResp.ClientState.Value, &wasmClientState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wasm client state: %w", err)
	}

	var ethClientState EthereumClientState
	if err := json.Unmarshal(wasmClientState.Data, &ethClientState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ethereum client state: %w", err)
	}

	return &ethClientState, nil
}

func (s *BeaconSpec) ToForkParameters() (*ForkParameters, error) {
	altairForkEpoch, err := strconv.ParseUint(s.AltairForkEpoch, 10, 64)
	if err != nil {
		return nil, err
	}
	bellatrixForkEpoch, err := strconv.ParseUint(s.BellatrixForkEpoch, 10, 64)
	if err != nil {
		return nil, err
	}
	capellaForkEpoch, err := strconv.ParseUint(s.CapellaForkEpoch, 10, 64)
	if err != nil {
		return nil, err
	}
	denebForkEpoch, err := strconv.ParseUint(s.DenebForkEpoch, 10, 64)
	if err != nil {
		return nil, err
	}
	electraForkEpoch, err := strconv.ParseUint(s.ElectraForkEpoch, 10, 64)
	if err != nil {
		return nil, err
	}
	return &ForkParameters{
		GenesisForkVersion: s.GenesisForkVersion,
		GenesisSlot:        0,
		Altair: Fork{
			Version: s.AltairForkVersion,
			Epoch:   altairForkEpoch,
		},
		Bellatrix: Fork{
			Version: s.BellatrixForkVersion,
			Epoch:   bellatrixForkEpoch,
		},
		Capella: Fork{
			Version: s.CapellaForkVersion,
			Epoch:   capellaForkEpoch,
		},
		Deneb: Fork{
			Version: s.DenebForkVersion,
			Epoch:   denebForkEpoch,
		},
		Electra: Fork{
			Version: s.ElectraForkVersion,
			Epoch:   electraForkEpoch,
		},
	}, nil
}
