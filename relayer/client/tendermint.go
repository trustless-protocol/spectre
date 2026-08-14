package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strconv"
	"strings"

	misbehaviourContract "relayer/bindings/Misbehaviour"
	spectreContract "relayer/bindings/SpectreClient"
	updateClientContract "relayer/bindings/UpdateClient"

	"github.com/cometbft/cometbft/p2p"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	commettypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v10/modules/core/23-commitment/types"
	ics23 "github.com/cosmos/ics23/go"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	MAX_TOTAL_VOTING_POWER = math.MaxInt64 / 8
	ProofType_EXIST        = 0
	ProofType_NON_EXIST    = 1
)

var clientStateType abi.Type
var consensusStateType abi.Type

// updateApplicationStateMsgType is the ABI tuple for ISpectreClientMsgs.MsgUpdateApplicationState.
// Note: no clientState field (the client reads it from its Store), and the proof
// fields are nested inside a `proof` (BatchProof) sub-tuple.
var updateApplicationStateMsgType abi.Type

// updateConsensusStateMsgType is the ABI tuple for ISpectreClientMsgs.MsgUpdateConsensusState:
// { MsgUpdateApplicationState update; ValidatorSet newValidatorSet }.
var updateConsensusStateMsgType abi.Type

// misbehaviourMsgType is the ABI tuple for
// ISpectreClientMsgs.MsgSubmitMisbehaviour.  It is hand-defined because the
// router accepts the message as opaque bytes rather than exposing a generated
// binding for this nested tuple.
var misbehaviourMsgType abi.Type

// ClientState mirrors IICS07TendermintMsgs.ClientState. It is defined locally
// because the on-chain client no longer exposes this struct in any ABI (client
// state is passed/returned as opaque `bytes`), so abigen emits no Go type for it.
// Field order/types must match clientStateComponents below.
type ClientState struct {
	ChainId         string
	TrustLevel      TrustThreshold
	LatestHeight    updateClientContract.IICS02ClientMsgsHeight
	TrustingPeriod  uint32
	UnbondingPeriod uint32
	IsFrozen        bool
	ClockDrift      uint32
}

// TrustThreshold mirrors IICS07TendermintMsgs.TrustThreshold (a Fraction).
type TrustThreshold struct {
	Numerator   uint8
	Denominator uint8
}

func init() {
	clientStateComponents := []abi.ArgumentMarshaling{
		{Name: "chainId", Type: "string"},
		{Name: "trustLevel", Type: "tuple", Components: []abi.ArgumentMarshaling{
			{Name: "numerator", Type: "uint8"},
			{Name: "denominator", Type: "uint8"},
		}},
		{Name: "latestHeight", Type: "tuple", Components: []abi.ArgumentMarshaling{
			{Name: "revisionNumber", Type: "uint64"},
			{Name: "revisionHeight", Type: "uint64"},
		}},
		{Name: "trustingPeriod", Type: "uint32"},
		{Name: "unbondingPeriod", Type: "uint32"},
		{Name: "isFrozen", Type: "bool"},
		{Name: "clockDrift", Type: "uint32"},
	}
	clientStateType, _ = abi.NewType("tuple", "", clientStateComponents)

	consensusStateComponents := []abi.ArgumentMarshaling{
		{Name: "timestamp", Type: "uint128"},
		{Name: "root", Type: "bytes32"},
		{Name: "nextValidatorsHash", Type: "bytes32"},
	}
	consensusStateType, _ = abi.NewType("tuple", "", consensusStateComponents)

	signedHeaderComponents := []abi.ArgumentMarshaling{
		{
			Name: "header",
			Type: "tuple",
			Components: []abi.ArgumentMarshaling{
				{Name: "version", Type: "tuple", Components: []abi.ArgumentMarshaling{
					{Name: "blockVersion", Type: "uint64"},
					{Name: "appVersion", Type: "uint64"},
				}},
				{Name: "chainId", Type: "string"},
				{Name: "height", Type: "uint64"},
				{Name: "time", Type: "uint128"},
				{Name: "hasLastBlockId", Type: "bool"},
				{Name: "lastBlockId", Type: "tuple", Components: []abi.ArgumentMarshaling{
					{Name: "hashData", Type: "bytes32"},
					{Name: "partSetHeader", Type: "tuple", Components: []abi.ArgumentMarshaling{
						{Name: "total", Type: "uint32"},
						{Name: "hashData", Type: "bytes32"},
					}},
				}},
				{Name: "hasLastCommitHash", Type: "bool"},
				{Name: "lastCommitHash", Type: "bytes32"},
				{Name: "hasDataHash", Type: "bool"},
				{Name: "dataHash", Type: "bytes32"},
				{Name: "validatorsHash", Type: "bytes32"},
				{Name: "nextValidatorsHash", Type: "bytes32"},
				{Name: "consensusHash", Type: "bytes32"},
				{Name: "appHash", Type: "bytes32"},
				{Name: "hasLastResultsHash", Type: "bool"},
				{Name: "lastResultsHash", Type: "bytes32"},
				{Name: "hasEvidenceHash", Type: "bool"},
				{Name: "evidenceHash", Type: "bytes32"},
				{Name: "proposerAddress", Type: "bytes"},
			},
		},
		{
			Name: "commit",
			Type: "tuple",
			Components: []abi.ArgumentMarshaling{
				{Name: "height", Type: "uint64"},
				{Name: "round", Type: "uint32"},
				{Name: "blockId", Type: "tuple", Components: []abi.ArgumentMarshaling{
					{Name: "hashData", Type: "bytes32"},
					{Name: "partSetHeader", Type: "tuple", Components: []abi.ArgumentMarshaling{
						{Name: "total", Type: "uint32"},
						{Name: "hashData", Type: "bytes32"},
					}},
				}},
				{Name: "commitSigs", Type: "tuple[]", Components: []abi.ArgumentMarshaling{
					{Name: "flag", Type: "uint8"},
				}},
			},
		},
	}

	headerComponents := []abi.ArgumentMarshaling{
		{Name: "signedHeader", Type: "tuple", Components: signedHeaderComponents},
		{Name: "trustedHeight", Type: "tuple", Components: []abi.ArgumentMarshaling{
			{Name: "revisionNumber", Type: "uint64"},
			{Name: "revisionHeight", Type: "uint64"},
		}},
	}

	// BatchProof sub-tuple, shared by all flows (matches ISpectreClientMsgs.BatchProof).
	batchProofComponents := []abi.ArgumentMarshaling{
		{Name: "proof", Type: "uint256[8]"},
		{Name: "commitments", Type: "uint256[2]"},
		{Name: "commitmentPok", Type: "uint256[2]"},
		{Name: "bucket", Type: "uint16"},
		{Name: "signerIndices", Type: "uint32[]"},
		{Name: "pinnedValidatorIndices", Type: "uint32[]"},
		{Name: "signerPubkeys", Type: "bytes32[]"},
		{Name: "active", Type: "bool[]"},
	}

	// MsgUpdateApplicationState: no clientState field; proof nested as BatchProof.
	applicationStateComponents := []abi.ArgumentMarshaling{
		{Name: "trustedConsensusState", Type: "tuple", Components: consensusStateComponents},
		{Name: "proposedHeader", Type: "tuple", Components: headerComponents},
		{Name: "time", Type: "uint128"},
		{Name: "proof", Type: "tuple", Components: batchProofComponents},
	}
	updateApplicationStateMsgType, _ = abi.NewType("tuple", "", applicationStateComponents)

	// ValidatorSet sub-tuple (matches IICS07TendermintMsgs.ValidatorSet).
	validatorInfoComponents := []abi.ArgumentMarshaling{
		{Name: "valAddress", Type: "bytes"},
		{Name: "pubKey", Type: "bytes32"},
		{Name: "votingPower", Type: "uint64"},
		{Name: "proposerPriority", Type: "int64"},
	}
	validatorSetComponents := []abi.ArgumentMarshaling{
		{Name: "validators", Type: "tuple[]", Components: validatorInfoComponents},
		{Name: "hasProposer", Type: "bool"},
		{Name: "proposer", Type: "tuple", Components: validatorInfoComponents},
		{Name: "totalVotingPower", Type: "uint64"},
	}

	// MsgUpdateConsensusState: { update: MsgUpdateApplicationState, newValidatorSet: ValidatorSet }.
	updateConsensusStateMsgType, _ = abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "update", Type: "tuple", Components: applicationStateComponents},
		{Name: "newValidatorSet", Type: "tuple", Components: validatorSetComponents},
	})

	misbehaviourMsgType, _ = abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "misbehaviour", Type: "tuple", Components: []abi.ArgumentMarshaling{
			{Name: "header1", Type: "tuple", Components: headerComponents},
			{Name: "header2", Type: "tuple", Components: headerComponents},
		}},
		{Name: "trustedConsensusState1", Type: "tuple", Components: consensusStateComponents},
		{Name: "trustedConsensusState2", Type: "tuple", Components: consensusStateComponents},
		{Name: "time", Type: "uint128"},
		{Name: "proof1", Type: "tuple", Components: batchProofComponents},
		{Name: "proof2", Type: "tuple", Components: batchProofComponents},
	})
}

type LightBlock struct {
	SignedHeader commettypes.SignedHeader
	ValSet       commettypes.ValidatorSet
	NextValSet   commettypes.ValidatorSet
	PeerId       p2p.ID
	BlockHeight  int64
}

func (b *LightBlock) IntoHeader(trustedBlock LightBlock) (updateClientContract.IICS07TendermintMsgsHeader, error) {
	if b == nil {
		return updateClientContract.IICS07TendermintMsgsHeader{}, fmt.Errorf("light block is nil")
	}
	revisionNumer := clienttypes.ParseChainID(trustedBlock.SignedHeader.ChainID)
	trustedHeight, err := nonNegativeInt64ToUint64("trusted block height", trustedBlock.BlockHeight)
	if err != nil {
		return updateClientContract.IICS07TendermintMsgsHeader{}, err
	}
	headerHeight, err := nonNegativeInt64ToUint64("proposed block height", b.BlockHeight)
	if err != nil {
		return updateClientContract.IICS07TendermintMsgsHeader{}, err
	}
	commitHeight, err := nonNegativeInt64ToUint64("commit height", b.SignedHeader.Commit.Height)
	if err != nil {
		return updateClientContract.IICS07TendermintMsgsHeader{}, err
	}

	commitSigs := []updateClientContract.IICS07TendermintMsgsCommitSig{}
	for _, sig := range b.SignedHeader.Commit.Signatures {
		// The validator address ties a commit slot to the validator it belongs to, so the
		// on-chain quorum check can confirm that the slot a proof cites is the same validator
		// whose pinned-set voting power it claims (ZK-09). CometBFT derives it the same way
		// Solidity does — the first 20 bytes of sha256 over the Ed25519 pubkey. ABSENT slots
		// carry no address and copy as zero, which is fine: only active signers are checked.
		var valAddr [20]byte
		copy(valAddr[:], sig.ValidatorAddress)

		// CometBFT: 0=UNKNOWN, 1=ABSENT, 2=COMMIT, 3=NIL
		// Solidity:  0=UNKNOWN, 1=ABSENT, 2=COMMIT, 3=NIL
		commitSigs = append(commitSigs, updateClientContract.IICS07TendermintMsgsCommitSig{
			Flag:             uint8(sig.BlockIDFlag),
			ValidatorAddress: valAddr,
		})
	}

	header := updateClientContract.IICS07TendermintMsgsHeader{
		TrustedHeight: updateClientContract.IICS02ClientMsgsHeight{
			RevisionNumber: revisionNumer,
			RevisionHeight: trustedHeight,
		},
		SignedHeader: updateClientContract.IICS07TendermintMsgsSignedHeader{
			Header: updateClientContract.IICS07TendermintMsgsBlockHeader{
				Version: updateClientContract.IICS07TendermintMsgsVersion{
					BlockVersion: b.SignedHeader.Version.Block,
					AppVersion:   b.SignedHeader.Version.App,
				},
				ChainId:        b.SignedHeader.ChainID,
				Height:         headerHeight,
				Time:           big.NewInt(b.SignedHeader.Time.UnixNano()),
				HasLastBlockId: !b.SignedHeader.LastBlockID.IsZero(),
				LastBlockId: updateClientContract.IICS07TendermintMsgsBlockId{
					HashData: bytesToBytes32(b.SignedHeader.LastBlockID.Hash),
					PartSetHeader: updateClientContract.IICS07TendermintMsgsPartSetHeader{
						Total:    b.SignedHeader.LastBlockID.PartSetHeader.Total,
						HashData: bytesToBytes32(b.SignedHeader.LastBlockID.PartSetHeader.Hash),
					},
				},
				HasLastCommitHash:  b.SignedHeader.LastCommitHash != nil,
				LastCommitHash:     bytesToBytes32(b.SignedHeader.LastCommitHash),
				HasDataHash:        b.SignedHeader.DataHash != nil,
				DataHash:           bytesToBytes32(b.SignedHeader.DataHash),
				ValidatorsHash:     bytesToBytes32(b.SignedHeader.ValidatorsHash),
				NextValidatorsHash: bytesToBytes32(b.SignedHeader.NextValidatorsHash),
				ConsensusHash:      bytesToBytes32(b.SignedHeader.ConsensusHash),
				AppHash:            bytesToBytes32(b.SignedHeader.AppHash),
				HasLastResultsHash: b.SignedHeader.LastResultsHash != nil,
				LastResultsHash:    bytesToBytes32(b.SignedHeader.LastResultsHash),
				HasEvidenceHash:    b.SignedHeader.EvidenceHash != nil,
				EvidenceHash:       bytesToBytes32(b.SignedHeader.EvidenceHash),
				ProposerAddress:    b.SignedHeader.ProposerAddress,
			},
			Commit: updateClientContract.IICS07TendermintMsgsBlockCommit{
				Height: commitHeight,
				Round:  uint32(b.SignedHeader.Commit.Round),
				BlockId: updateClientContract.IICS07TendermintMsgsBlockId{
					HashData: bytesToBytes32(b.SignedHeader.Commit.BlockID.Hash),
					PartSetHeader: updateClientContract.IICS07TendermintMsgsPartSetHeader{
						Total:    b.SignedHeader.Commit.BlockID.PartSetHeader.Total,
						HashData: bytesToBytes32(b.SignedHeader.Commit.BlockID.PartSetHeader.Hash),
					},
				},
				CommitSigs: commitSigs,
			},
		},
	}

	return header, nil
}

func nonNegativeInt64ToUint64(name string, value int64) (uint64, error) {
	if value < 0 {
		return 0, fmt.Errorf("%s cannot be negative: %d", name, value)
	}
	return uint64(value), nil
}

func votingPowerToUint64(name string, value int64) (uint64, error) {
	if value < 0 {
		return 0, fmt.Errorf("%s cannot be negative: %d", name, value)
	}
	if value > MAX_TOTAL_VOTING_POWER {
		return 0, fmt.Errorf("%s %d exceeds max total voting power %d", name, value, MAX_TOTAL_VOTING_POWER)
	}
	return uint64(value), nil
}

func validatorInfoToContract(name string, val *commettypes.Validator) (spectreContract.IICS07TendermintMsgsValidatorInfo, error) {
	if val == nil {
		return spectreContract.IICS07TendermintMsgsValidatorInfo{}, nil
	}
	votingPower, err := votingPowerToUint64(name+" voting power", val.VotingPower)
	if err != nil {
		return spectreContract.IICS07TendermintMsgsValidatorInfo{}, err
	}
	return spectreContract.IICS07TendermintMsgsValidatorInfo{
		ValAddress:       val.Address,
		PubKey:           bytesToBytes32(val.PubKey.Bytes()),
		VotingPower:      votingPower,
		ProposerPriority: val.ProposerPriority,
	}, nil
}

func ValidatorSetToContract(valSet commettypes.ValidatorSet, name string) (spectreContract.IICS07TendermintMsgsValidatorSet, error) {
	vals := []spectreContract.IICS07TendermintMsgsValidatorInfo{}
	for i, val := range valSet.Validators {
		info, err := validatorInfoToContract(fmt.Sprintf("%s validator[%d]", name, i), val)
		if err != nil {
			return spectreContract.IICS07TendermintMsgsValidatorSet{}, err
		}
		vals = append(vals, info)
	}
	proposer, err := validatorInfoToContract(name+" proposer", valSet.Proposer)
	if err != nil {
		return spectreContract.IICS07TendermintMsgsValidatorSet{}, err
	}
	totalVotingPower, err := votingPowerToUint64(name+" total voting power", valSet.TotalVotingPower())
	if err != nil {
		return spectreContract.IICS07TendermintMsgsValidatorSet{}, err
	}
	return spectreContract.IICS07TendermintMsgsValidatorSet{
		Validators:       vals,
		HasProposer:      valSet.Proposer != nil,
		Proposer:         proposer,
		TotalVotingPower: totalVotingPower,
	}, nil
}

type ContractValidatorSet = spectreContract.IICS07TendermintMsgsValidatorSet

type SpectreClientGenesis struct {
	TrustedClientState        ClientState
	TrustedConsensusState     updateClientContract.IICS07TendermintMsgsConsensusState
	InitialPinnedValidatorSet spectreContract.IICS07TendermintMsgsValidatorSet
}

// DefaultClockDrift is the allowed gap (in seconds) between the proven
// consensus-state timestamp and the verifying chain's block time. It must be
// generous enough to cover relay latency (proof gen + destination block time +
// queueing); too small a value makes the light client reject otherwise-valid
// updates with ProofIsTooOld. The same value is stored in the client state at
// creation and drives the on-chain freshness window in `_requireFreshness`
// (ProofIsInTheFuture / ProofIsTooOld).
const DefaultClockDrift uint32 = 30

func GetGenesis(client *rpchttp.HTTP, trustedBlock int64, trustingPeriod uint32, trustLevel string, _ string, clockDrift uint32) (*SpectreClientGenesis, error) {
	status, err := client.Status(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	if trustedBlock == 0 {
		// get latest block
		trustedBlock = status.SyncInfo.LatestBlockHeight
	}

	trustedLightBlock, err := GetLightBlock(client, trustedBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get light block: %w", err)
	}

	unbondingPeriod, err := GetUnbondingTime(client)
	if err != nil {
		return nil, fmt.Errorf("failed to get unbonding time: %w", err)
	}

	return spectreClientGenesisFromLightBlock(trustedLightBlock, unbondingPeriod, trustingPeriod, trustLevel, clockDrift)
}

func spectreClientGenesisFromLightBlock(trustedLightBlock *LightBlock, unbondingPeriod float64, trustingPeriod uint32, trustLevel string, clockDrift uint32) (*SpectreClientGenesis, error) {
	if trustedLightBlock == nil {
		return nil, fmt.Errorf("trusted light block is nil")
	}
	if clockDrift == 0 {
		clockDrift = DefaultClockDrift
	}
	if trustingPeriod == 0 {
		trustingPeriod = uint32(unbondingPeriod * 2 / 3)
	}

	if trustingPeriod > uint32(unbondingPeriod) {
		return nil, fmt.Errorf("trusting period %d cannot be greater than unbonding period %d", trustingPeriod, uint32(unbondingPeriod))
	}

	chainId := trustedLightBlock.SignedHeader.Header.ChainID
	revision := clienttypes.ParseChainID(chainId)

	trustThreshold, err := ParseTrustThreshold(trustLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trust level: %w", err)
	}

	latestRevisionHeight, err := nonNegativeInt64ToUint64("trusted light block header height", trustedLightBlock.SignedHeader.Header.Height)
	if err != nil {
		return nil, err
	}

	clientState := ClientState{
		ChainId:    chainId,
		TrustLevel: trustThreshold,
		LatestHeight: updateClientContract.IICS02ClientMsgsHeight{
			RevisionNumber: revision,
			RevisionHeight: latestRevisionHeight,
		},
		IsFrozen:        false,
		TrustingPeriod:  trustingPeriod,
		UnbondingPeriod: uint32(unbondingPeriod),
		ClockDrift:      clockDrift,
	}

	consensusState := updateClientContract.IICS07TendermintMsgsConsensusState{
		Timestamp:          big.NewInt(trustedLightBlock.SignedHeader.Header.Time.UnixNano()),
		Root:               bytesToBytes32(trustedLightBlock.SignedHeader.Header.AppHash),
		NextValidatorsHash: bytesToBytes32(trustedLightBlock.SignedHeader.NextValidatorsHash),
	}

	initialPinnedValidatorSet, err := ValidatorSetToContract(trustedLightBlock.NextValSet, "initial pinned")
	if err != nil {
		return nil, err
	}

	genesis := SpectreClientGenesis{
		TrustedClientState:        clientState,
		TrustedConsensusState:     consensusState,
		InitialPinnedValidatorSet: initialPinnedValidatorSet,
	}
	return &genesis, nil
}

func GetUnbondingTime(client *rpchttp.HTTP) (float64, error) {

	// Query path for staking parameters
	queryPath := "/cosmos.staking.v1beta1.Query/Params"

	// Make ABCI query
	result, err := client.ABCIQuery(context.Background(), queryPath, nil)
	if err != nil {
		return 0, fmt.Errorf("ABCI query failed: %w", err)
	}

	if result.Response.Code != 0 {
		return 0, fmt.Errorf("query failed with code %d: %s", result.Response.Code, result.Response.Log)
	}

	// Create codec for decoding
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	// Decode the response
	var params stakingtypes.QueryParamsResponse
	if err := cdc.Unmarshal(result.Response.Value, &params); err != nil {
		return 0, fmt.Errorf("failed to unmarshal params: %w", err)
	}
	return params.Params.UnbondingTime.Seconds(), nil
}

func GetLatestLightBlock(client *rpchttp.HTTP) (*LightBlock, error) {
	status, err := client.Status(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	return GetLightBlock(client, status.SyncInfo.LatestBlockHeight)
}

// cometBFTMaxPerPage is the largest page size the CometBFT /validators RPC
// honors (rpc/core/env.go maxPerPage). Requesting this many per page minimizes
// round-trips while still reading the full set.
const cometBFTMaxPerPage = 100

// validatorsPager is the subset of the CometBFT RPC client used to page the
// validator set. *rpchttp.HTTP satisfies it; tests supply a fake.
type validatorsPager interface {
	Validators(ctx context.Context, height *int64, page, perPage *int) (*coretypes.ResultValidators, error)
}

// fetchAllValidators returns the complete validator set at the given height.
//
// The CometBFT /validators RPC paginates and defaults to perPage=30 when perPage
// is nil (rpc/core/env.go validatePerPage). Calling it once with nil therefore
// silently truncates the set to the first 30 validators on any chain with more
// than 30 — the assembled set then has the wrong hash and signer extraction
// fails for indices >= 30. We page explicitly and loop until the reported Total
// is collected (issue #105).
func fetchAllValidators(client validatorsPager, height int64) ([]*commettypes.Validator, error) {
	return fetchAllValidatorsWithContext(context.Background(), client, height)
}

func fetchAllValidatorsWithContext(ctx context.Context, client validatorsPager, height int64) ([]*commettypes.Validator, error) {
	var collected []*commettypes.Validator
	perPage := cometBFTMaxPerPage
	for page := 1; ; page++ {
		p := page
		resp, err := client.Validators(ctx, &height, &p, &perPage)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch validators (height %d, page %d): %w", height, page, err)
		}
		collected = append(collected, resp.Validators...)
		// Stop once we have the full set, or defensively if a page is empty.
		if len(collected) >= resp.Total || len(resp.Validators) == 0 {
			break
		}
	}
	return collected, nil
}

func GetLightBlock(client *rpchttp.HTTP, height int64) (*LightBlock, error) {
	return GetLightBlockWithContext(context.Background(), client, height)
}

// GetLightBlockWithContext fetches a complete light block while propagating
// cancellation through every CometBFT request. GetLightBlock remains as the
// compatibility wrapper for relay paths that do not yet carry a context.
func GetLightBlockWithContext(ctx context.Context, client *rpchttp.HTTP, height int64) (*LightBlock, error) {
	status, err := client.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	peerId := status.NodeInfo.ID()
	commitResp, err := client.Commit(ctx, &height)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trusted commit: %w", err)
	}

	signedHeader := commitResp.SignedHeader
	proposerAddr := signedHeader.Header.ProposerAddress
	var proposer *commettypes.Validator
	validators, err := fetchAllValidatorsWithContext(ctx, client, height)
	if err != nil {
		return nil, err
	}

	for _, resp := range validators {
		if resp != nil {
			if slices.Equal(resp.Address.Bytes(), proposerAddr.Bytes()) {
				proposer = resp
				break
			}
		}
	}

	valSet := commettypes.NewValidatorSet(validators)
	valSet.Proposer = proposer
	// Guard against a truncated/incomplete fetch: the assembled set must hash to
	// the value committed in the header, otherwise the on-chain validator-set
	// hash check would fail and signer indices would be misaligned (issue #105).
	if !bytes.Equal(valSet.Hash(), signedHeader.Header.ValidatorsHash) {
		return nil, fmt.Errorf(
			"validator set hash mismatch at height %d: assembled %X != header ValidatorsHash %X (got %d validators)",
			height, valSet.Hash(), signedHeader.Header.ValidatorsHash.Bytes(), len(validators))
	}

	nextHeight := height + 1
	nextValidators, err := fetchAllValidatorsWithContext(ctx, client, nextHeight)
	if err != nil {
		return nil, err
	}

	for _, resp := range nextValidators {
		if resp != nil {
			if slices.Equal(resp.Address.Bytes(), proposerAddr.Bytes()) {
				proposer = resp
				break
			}
		}
	}

	nextValSet := commettypes.NewValidatorSet(nextValidators)
	nextValSet.Proposer = proposer
	if !bytes.Equal(nextValSet.Hash(), signedHeader.Header.NextValidatorsHash) {
		return nil, fmt.Errorf(
			"next validator set hash mismatch at height %d: assembled %X != header NextValidatorsHash %X (got %d validators)",
			nextHeight, nextValSet.Hash(), signedHeader.Header.NextValidatorsHash.Bytes(), len(nextValidators))
	}

	return &LightBlock{
		SignedHeader: signedHeader,
		ValSet:       *valSet,
		NextValSet:   *nextValSet,
		PeerId:       peerId,
		BlockHeight:  height,
	}, nil

}

func ProvePath(client *rpchttp.HTTP, height int64, path [][]byte) ([]byte, *commitmenttypes.MerkleProof, error) {
	queryPath := fmt.Sprintf("store/%s/key", string(path[0]))
	request := slices.Concat(path[1:]...)

	// Make ABCI query
	result, err := client.ABCIQueryWithOptions(context.Background(), queryPath, request, rpcclient.ABCIQueryOptions{
		// Proof height should be the block before the target block.
		Height: height - 1,
		Prove:  true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("ABCI query failed: %w", err)
	}

	if result.Response.Code != 0 {
		return nil, nil, fmt.Errorf("query failed with code %d: %s", result.Response.Code, result.Response.Log)
	}

	if result.Response.Height != height-1 {
		return nil, nil, fmt.Errorf("proof height mismatch")
	}

	if !slices.Equal(result.Response.Key, path[1]) {
		return nil, nil, fmt.Errorf("key mismatch")
	}

	proof, err := commitmenttypes.ConvertProofs(result.Response.ProofOps)
	if err != nil {
		return nil, nil, fmt.Errorf("proof could not be retrieved: %w", err)
	}
	if len(proof.Proofs) == 0 {
		return nil, nil, fmt.Errorf("proof is empty")
	}
	return result.Response.Value, &proof, nil
}

// parseTrustThreshold parses a trust threshold fraction string like "2/3"
func ParseTrustThreshold(value string) (TrustThreshold, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return TrustThreshold{}, fmt.Errorf("invalid trust threshold format: %s (expected format: 'numerator/denominator')", value)
	}

	// bitSize 8 so values that don't fit the uint8 contract fields are a
	// parse error instead of silently truncating (e.g. "300/6" -> 44/6).
	numerator, err := strconv.ParseUint(parts[0], 10, 8)
	if err != nil {
		return TrustThreshold{}, fmt.Errorf("invalid numerator: %s", parts[0])
	}

	denominator, err := strconv.ParseUint(parts[1], 10, 8)
	if err != nil {
		return TrustThreshold{}, fmt.Errorf("invalid denominator: %s", parts[1])
	}

	if denominator == 0 {
		return TrustThreshold{}, fmt.Errorf("denominator cannot be zero")
	}

	return TrustThreshold{
		Numerator:   uint8(numerator),
		Denominator: uint8(denominator),
	}, nil
}

func ParseCommitmentProof(proof *ics23.CommitmentProof) (*spectreContract.IMembershipMsgsCommitmentProof, error) {
	if proof == nil {
		return nil, fmt.Errorf("proof is nil")
	}

	var parsedProof *spectreContract.IMembershipMsgsCommitmentProof
	switch p := proof.Proof.(type) {
	case *ics23.CommitmentProof_Exist:
		parsedProof = &spectreContract.IMembershipMsgsCommitmentProof{
			ProofType: ProofType_EXIST,
			ExistenceProof: spectreContract.IMembershipMsgsExistenceProof{
				Key:   p.Exist.Key,
				Value: p.Exist.Value,
				Leaf:  ParseLeafOp(p.Exist.Leaf),
				Path:  []spectreContract.IMembershipMsgsInnerOp{},
			},
			NonExistenceProof: spectreContract.IMembershipMsgsNonExistenceProof{},
		}

		for _, innerOp := range p.Exist.Path {
			parsedProof.ExistenceProof.Path = append(parsedProof.ExistenceProof.Path, ParseInnerOp(innerOp))
		}
	case *ics23.CommitmentProof_Nonexist:
		if p.Nonexist.Left == nil && p.Nonexist.Right == nil {
			return nil, fmt.Errorf("non-existence proof must at least left or right existence proofs")
		}
		parsedProof = &spectreContract.IMembershipMsgsCommitmentProof{
			ProofType:      ProofType_NON_EXIST,
			ExistenceProof: spectreContract.IMembershipMsgsExistenceProof{},
			NonExistenceProof: spectreContract.IMembershipMsgsNonExistenceProof{
				Key:      p.Nonexist.Key,
				HasLeft:  false,
				Left:     spectreContract.IMembershipMsgsExistenceProof{},
				HasRight: false,
				Right:    spectreContract.IMembershipMsgsExistenceProof{},
			},
		}

		if p.Nonexist.Left != nil {
			parsedProof.NonExistenceProof.HasLeft = true
			parsedProof.NonExistenceProof.Left = spectreContract.IMembershipMsgsExistenceProof{
				Key:   p.Nonexist.Left.Key,
				Value: p.Nonexist.Left.Value,
				Leaf:  ParseLeafOp(p.Nonexist.Left.Leaf),
				Path:  []spectreContract.IMembershipMsgsInnerOp{},
			}
			for _, innerOp := range p.Nonexist.Left.Path {
				parsedProof.NonExistenceProof.Left.Path = append(parsedProof.NonExistenceProof.Left.Path, ParseInnerOp(innerOp))
			}
		}

		if p.Nonexist.Right != nil {
			parsedProof.NonExistenceProof.HasRight = true
			parsedProof.NonExistenceProof.Right = spectreContract.IMembershipMsgsExistenceProof{
				Key:   p.Nonexist.Right.Key,
				Value: p.Nonexist.Right.Value,
				Leaf:  ParseLeafOp(p.Nonexist.Right.Leaf),
				Path:  []spectreContract.IMembershipMsgsInnerOp{},
			}
			for _, innerOp := range p.Nonexist.Right.Path {
				parsedProof.NonExistenceProof.Right.Path = append(parsedProof.NonExistenceProof.Right.Path, ParseInnerOp(innerOp))
			}
		}

	case *ics23.CommitmentProof_Batch:
		if len(p.Batch.GetEntries()) == 0 || p.Batch.GetEntries()[0] == nil {
			return nil, fmt.Errorf("batch proof has empty entry")
		}

		if e := p.Batch.GetEntries()[0].GetExist(); e != nil {
			parsedProof = &spectreContract.IMembershipMsgsCommitmentProof{
				ProofType: ProofType_EXIST,
				ExistenceProof: spectreContract.IMembershipMsgsExistenceProof{
					Key:   e.Key,
					Value: e.Value,
					Leaf:  ParseLeafOp(e.Leaf),
					Path:  []spectreContract.IMembershipMsgsInnerOp{},
				},
				NonExistenceProof: spectreContract.IMembershipMsgsNonExistenceProof{},
			}

			for _, innerOp := range e.Path {
				parsedProof.ExistenceProof.Path = append(parsedProof.ExistenceProof.Path, ParseInnerOp(innerOp))
			}
		}

		if n := p.Batch.GetEntries()[0].GetNonexist(); n != nil {
			parsedProof = &spectreContract.IMembershipMsgsCommitmentProof{
				ProofType:      ProofType_NON_EXIST,
				ExistenceProof: spectreContract.IMembershipMsgsExistenceProof{},
				NonExistenceProof: spectreContract.IMembershipMsgsNonExistenceProof{
					Key:      n.Key,
					HasLeft:  false,
					Left:     spectreContract.IMembershipMsgsExistenceProof{},
					HasRight: false,
					Right:    spectreContract.IMembershipMsgsExistenceProof{},
				},
			}

			if n.Left != nil {
				parsedProof.NonExistenceProof.HasLeft = true
				parsedProof.NonExistenceProof.Left = spectreContract.IMembershipMsgsExistenceProof{
					Key:   n.Left.Key,
					Value: n.Left.Value,
					Leaf:  ParseLeafOp(n.Left.Leaf),
					Path:  []spectreContract.IMembershipMsgsInnerOp{},
				}
				for _, innerOp := range n.Left.Path {
					parsedProof.NonExistenceProof.Left.Path = append(parsedProof.NonExistenceProof.Left.Path, ParseInnerOp(innerOp))
				}
			}

			if n.Right != nil {
				parsedProof.NonExistenceProof.HasRight = true
				parsedProof.NonExistenceProof.Right = spectreContract.IMembershipMsgsExistenceProof{
					Key:   n.Right.Key,
					Value: n.Right.Value,
					Leaf:  ParseLeafOp(n.Right.Leaf),
					Path:  []spectreContract.IMembershipMsgsInnerOp{},
				}
				for _, innerOp := range n.Right.Path {
					parsedProof.NonExistenceProof.Right.Path = append(parsedProof.NonExistenceProof.Right.Path, ParseInnerOp(innerOp))
				}
			}
		}
	case *ics23.CommitmentProof_Compressed:
		decompressedProof := ics23.Decompress(proof)
		return ParseCommitmentProof(decompressedProof)
	default:
		return nil, fmt.Errorf("unrecognized proof type")
	}

	return parsedProof, nil
}

func ParseLeafOp(leafOp *ics23.LeafOp) spectreContract.IMembershipMsgsLeafOp {
	if leafOp == nil {
		return spectreContract.IMembershipMsgsLeafOp{}
	}

	return spectreContract.IMembershipMsgsLeafOp{
		HashOp:       uint8(leafOp.Hash),
		PrehashKey:   uint8(leafOp.PrehashKey),
		PrehashValue: uint8(leafOp.PrehashValue),
		Prefix:       leafOp.Prefix,
	}
}

func ParseInnerOp(innerOp *ics23.InnerOp) spectreContract.IMembershipMsgsInnerOp {
	if innerOp == nil {
		return spectreContract.IMembershipMsgsInnerOp{}
	}

	return spectreContract.IMembershipMsgsInnerOp{
		HashOp: uint8(innerOp.Hash),
		Prefix: innerOp.Prefix,
		Suffix: innerOp.Suffix,
	}
}

func EncodeClientState(clientState ClientState) ([]byte, error) {
	args := abi.Arguments{
		{Type: clientStateType},
	}
	encoded, err := args.Pack(clientState)
	return encoded, err
}

func DecodeClientState(data []byte) (ClientState, error) {
	args := abi.Arguments{
		{Type: clientStateType},
	}
	unpacked, err := args.Unpack(data)
	if err != nil {
		return ClientState{}, fmt.Errorf("unpack: %w", err)
	}
	if len(unpacked) == 0 {
		return ClientState{}, fmt.Errorf("no data unpacked")
	}
	// unpacked[0] is an anonymous struct matching the tuple.
	// Use JSON roundtrip to convert to the named target type,
	// since args.Copy maps the tuple to the first field instead of the struct itself.
	jsonBytes, err := json.Marshal(unpacked[0])
	if err != nil {
		return ClientState{}, fmt.Errorf("marshal unpacked tuple: %w", err)
	}
	var clientState ClientState
	if err := json.Unmarshal(jsonBytes, &clientState); err != nil {
		return ClientState{}, fmt.Errorf("unmarshal to client state: %w", err)
	}
	return clientState, nil
}

func EncodeConsensusState(consensusState updateClientContract.IICS07TendermintMsgsConsensusState) ([]byte, error) {
	args := abi.Arguments{
		{Type: consensusStateType},
	}
	encoded, err := args.Pack(consensusState)
	return encoded, err

}

// EncodeUpdateApplicationStateMsg abi-encodes a MsgUpdateApplicationState for
// SpectreClient.updateApplicationState(bytes) / ICS26Router.updateApplicationState.
func EncodeUpdateApplicationStateMsg(msg updateClientContract.ISpectreClientMsgsMsgUpdateApplicationState) ([]byte, error) {
	args := abi.Arguments{
		{Type: updateApplicationStateMsgType},
	}
	return args.Pack(msg)
}

// MsgUpdateConsensusState mirrors ISpectreClientMsgs.MsgUpdateConsensusState.
// There is no generated Go type for it (nothing takes it as a typed on-chain
// param), so it is hand-defined and packed against updateConsensusStateMsgType.
type MsgUpdateConsensusState struct {
	Update          updateClientContract.ISpectreClientMsgsMsgUpdateApplicationState
	NewValidatorSet spectreContract.IICS07TendermintMsgsValidatorSet
}

// EncodeUpdateConsensusStateMsg abi-encodes a MsgUpdateConsensusState for
// SpectreClient.updateConsensusState(bytes) / ICS26Router.updateConsensusState.
func EncodeUpdateConsensusStateMsg(update updateClientContract.ISpectreClientMsgsMsgUpdateApplicationState, newValidatorSet spectreContract.IICS07TendermintMsgsValidatorSet) ([]byte, error) {
	args := abi.Arguments{
		{Type: updateConsensusStateMsgType},
	}
	return args.Pack(MsgUpdateConsensusState{Update: update, NewValidatorSet: newValidatorSet})
}

// EncodeMisbehaviourContractMsg encodes the generated Misbehaviour binding
// type. The generated binding uses package-local struct types, so this helper
// shares the same ABI tuple while allowing the services package to preserve
// the exact proof metadata types produced by that binding.
func EncodeMisbehaviourContractMsg(msg misbehaviourContract.ISpectreClientMsgsMsgSubmitMisbehaviour) ([]byte, error) {
	args := abi.Arguments{{Type: misbehaviourMsgType}}
	return args.Pack(msg)
}

func bytesToBytes32(data []byte) [32]byte {
	var result [32]byte
	copy(result[:], data)
	return result
}
