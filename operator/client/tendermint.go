package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strconv"
	"strings"

	tendermintContract "operator/bindings/SP1ICS07Tendermint"
	updateClientContract "operator/bindings/UpdateClient"

	"github.com/cometbft/cometbft/p2p"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
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
var updateClientMsgType abi.Type

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
		{Name: "zkAlgorithm", Type: "uint8"},
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
					{Name: "data", Type: "tuple", Components: []abi.ArgumentMarshaling{
						{Name: "validatorAddress", Type: "bytes"},
						{Name: "timestamp", Type: "uint128"},
						{Name: "hasSignature", Type: "bool"},
						{Name: "signature", Type: "bytes"},
					}},
				}},
			},
		},
	}

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

	headerComponents := []abi.ArgumentMarshaling{
		{Name: "signedHeader", Type: "tuple", Components: signedHeaderComponents},
		{Name: "validatorSet", Type: "tuple", Components: validatorSetComponents},
		{Name: "trustedHeight", Type: "tuple", Components: []abi.ArgumentMarshaling{
			{Name: "revisionNumber", Type: "uint64"},
			{Name: "revisionHeight", Type: "uint64"},
		}},
		{Name: "trustedNextValidatorSet", Type: "tuple", Components: validatorSetComponents},
	}

	updateClientMsgType, _ = abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "clientState", Type: "tuple", Components: clientStateComponents},
		{Name: "trustedConsensusState", Type: "tuple", Components: consensusStateComponents},
		{Name: "proposedHeader", Type: "tuple", Components: headerComponents},
		{Name: "time", Type: "uint128"},
		{Name: "proof", Type: "uint256[8]"},
		{Name: "commitments", Type: "uint256[2]"},
		{Name: "commitmentPok", Type: "uint256[2]"},
	})
}

type LightBlock struct {
	SignedHeader commettypes.SignedHeader
	ValSet       commettypes.ValidatorSet
	NextValSet   commettypes.ValidatorSet
	PeerId       p2p.ID
	BlockHeight  int64
}

func (b *LightBlock) IntoHeader(trustedBlock LightBlock) updateClientContract.IICS07TendermintMsgsHeader {
	revisionNumer := clienttypes.ParseChainID(trustedBlock.SignedHeader.ChainID)

	commitSigs := []updateClientContract.IICS07TendermintMsgsCommitSig{}
	for _, sig := range b.SignedHeader.Commit.Signatures {
		// CometBFT: 0=UNKNOWN, 1=ABSENT, 2=COMMIT, 3=NIL
		// Solidity:            0=ABSENT, 1=COMMIT, 2=NIL
		commitSigs = append(commitSigs, updateClientContract.IICS07TendermintMsgsCommitSig{
			Flag: uint8(sig.BlockIDFlag - 1),
			Data: updateClientContract.IICS07TendermintMsgsCommitSigData{
				ValidatorAddress: sig.ValidatorAddress,
				Timestamp:        big.NewInt(sig.Timestamp.UnixNano()),
				HasSignature:     sig.Signature != nil,
				Signature:        sig.Signature,
			},
		})
	}

	vals := []updateClientContract.IICS07TendermintMsgsValidatorInfo{}
	for _, val := range b.ValSet.Validators {
		vals = append(vals, updateClientContract.IICS07TendermintMsgsValidatorInfo{
			ValAddress:       val.Address,
			PubKey:           bytesToBytes32(val.PubKey.Bytes()),
			VotingPower:      uint64(val.VotingPower),
			ProposerPriority: val.ProposerPriority,
		})
	}

	nextVals := []updateClientContract.IICS07TendermintMsgsValidatorInfo{}
	for _, val := range trustedBlock.NextValSet.Validators {
		nextVals = append(nextVals, updateClientContract.IICS07TendermintMsgsValidatorInfo{
			ValAddress:       val.Address,
			PubKey:           bytesToBytes32(val.PubKey.Bytes()),
			VotingPower:      uint64(val.VotingPower),
			ProposerPriority: val.ProposerPriority,
		})
	}

	header := updateClientContract.IICS07TendermintMsgsHeader{
		TrustedHeight: updateClientContract.IICS02ClientMsgsHeight{
			RevisionNumber: revisionNumer,
			RevisionHeight: uint64(trustedBlock.BlockHeight),
		},
		SignedHeader: updateClientContract.IICS07TendermintMsgsSignedHeader{
			Header: updateClientContract.IICS07TendermintMsgsBlockHeader{
				Version: updateClientContract.IICS07TendermintMsgsVersion{
					BlockVersion: b.SignedHeader.Version.Block,
					AppVersion:   b.SignedHeader.Version.App,
				},
				ChainId:        b.SignedHeader.ChainID,
				Height:         uint64(b.BlockHeight),
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
				Height: uint64(b.SignedHeader.Commit.Height),
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
		ValidatorSet: updateClientContract.IICS07TendermintMsgsValidatorSet{
			Validators:  vals,
			HasProposer: b.ValSet.Proposer != nil,
			Proposer: updateClientContract.IICS07TendermintMsgsValidatorInfo{
				ValAddress:       b.ValSet.Proposer.Address,
				PubKey:           bytesToBytes32(b.ValSet.Proposer.PubKey.Bytes()),
				VotingPower:      uint64(b.ValSet.Proposer.VotingPower),
				ProposerPriority: b.ValSet.Proposer.ProposerPriority,
			},
			TotalVotingPower: uint64(b.ValSet.TotalVotingPower()),
		},
		TrustedNextValidatorSet: updateClientContract.IICS07TendermintMsgsValidatorSet{
			Validators:  nextVals,
			HasProposer: trustedBlock.ValSet.Proposer != nil,
			Proposer: updateClientContract.IICS07TendermintMsgsValidatorInfo{
				ValAddress:       trustedBlock.ValSet.Proposer.Address,
				PubKey:           bytesToBytes32(trustedBlock.ValSet.Proposer.PubKey.Bytes()),
				VotingPower:      uint64(trustedBlock.ValSet.Proposer.VotingPower),
				ProposerPriority: trustedBlock.ValSet.Proposer.ProposerPriority,
			},
			TotalVotingPower: uint64(trustedBlock.ValSet.TotalVotingPower()),
		},
	}

	return header
}

type SP1ICS07TendermintGenesis struct {
	TrustedClientState    updateClientContract.IICS07TendermintMsgsClientState
	TrustedConsensusState updateClientContract.IICS07TendermintMsgsConsensusState
}

type SupportedZkAlgorithm uint8

const (
	Groth16 SupportedZkAlgorithm = iota
	Plonk
)

// String returns the string representation of the algorithm
func (s SupportedZkAlgorithm) String() string {
	switch s {
	case Groth16:
		return "Groth16"
	case Plonk:
		return "Plonk"
	default:
		return "Unknown"
	}
}

func GetGenesis(client *rpchttp.HTTP, trustedBlock int64, trustingPeriod uint32, trustLevel string, proofType string) (*SP1ICS07TendermintGenesis, error) {
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
	_ = trustedLightBlock

	unbondingPeriod, err := GetUnbondingTime(client)
	if err != nil {
		return nil, fmt.Errorf("failed to get unbonding time: %w", err)
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

	var zkAlgorithm SupportedZkAlgorithm
	switch proofType {
	case "groth16":
		zkAlgorithm = Groth16
	case "plonk":
		zkAlgorithm = Plonk
	default:
		return nil, fmt.Errorf("unsupported proof type: %s, supported types are: groth16, plonk", proofType)
	}

	clientState := updateClientContract.IICS07TendermintMsgsClientState{
		ChainId:    chainId,
		TrustLevel: trustThreshold,
		LatestHeight: updateClientContract.IICS02ClientMsgsHeight{
			RevisionNumber: revision,
			RevisionHeight: uint64(trustedLightBlock.SignedHeader.Header.Height),
		},
		IsFrozen:        false,
		ZkAlgorithm:     uint8(zkAlgorithm),
		TrustingPeriod:  trustingPeriod,
		UnbondingPeriod: uint32(unbondingPeriod),
	}

	consensusState := updateClientContract.IICS07TendermintMsgsConsensusState{
		Timestamp:          big.NewInt(trustedLightBlock.SignedHeader.Header.Time.UnixNano()),
		Root:               bytesToBytes32(trustedLightBlock.SignedHeader.Header.AppHash),
		NextValidatorsHash: bytesToBytes32(trustedLightBlock.SignedHeader.NextValidatorsHash),
	}

	genesis := SP1ICS07TendermintGenesis{
		TrustedClientState:    clientState,
		TrustedConsensusState: consensusState,
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

func GetLightBlock(client *rpchttp.HTTP, height int64) (*LightBlock, error) {
	status, err := client.Status(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	peerId := status.NodeInfo.ID()
	commitResp, err := client.Commit(context.Background(), &height)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trusted commit: %w", err)
	}

	signedHeader := commitResp.SignedHeader
	proposerAddr := signedHeader.Header.ProposerAddress
	var proposer *commettypes.Validator
	validatorResp, err := client.Validators(context.Background(), &height, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch validators: %w", err)
	}

	for _, resp := range validatorResp.Validators {
		if resp != nil {
			if slices.Equal(resp.Address.Bytes(), proposerAddr.Bytes()) {
				proposer = resp
				break
			}
		}
	}

	valSet := commettypes.NewValidatorSet(validatorResp.Validators)
	valSet.Proposer = proposer

	nextHeight := height + 1
	nextValidatorResp, err := client.Validators(context.Background(), &nextHeight, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch next validators: %w", err)
	}

	for _, resp := range nextValidatorResp.Validators {
		if resp != nil {
			if slices.Equal(resp.Address.Bytes(), proposerAddr.Bytes()) {
				proposer = resp
				break
			}
		}
	}

	nextValSet := commettypes.NewValidatorSet(nextValidatorResp.Validators)
	nextValSet.Proposer = proposer

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
func ParseTrustThreshold(value string) (updateClientContract.IICS07TendermintMsgsTrustThreshold, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return updateClientContract.IICS07TendermintMsgsTrustThreshold{}, fmt.Errorf("invalid trust threshold format: %s (expected format: 'numerator/denominator')", value)
	}

	numerator, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return updateClientContract.IICS07TendermintMsgsTrustThreshold{}, fmt.Errorf("invalid numerator: %s", parts[0])
	}

	denominator, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return updateClientContract.IICS07TendermintMsgsTrustThreshold{}, fmt.Errorf("invalid denominator: %s", parts[1])
	}

	if denominator == 0 {
		return updateClientContract.IICS07TendermintMsgsTrustThreshold{}, fmt.Errorf("denominator cannot be zero")
	}

	return updateClientContract.IICS07TendermintMsgsTrustThreshold{
		Numerator:   uint8(numerator),
		Denominator: uint8(denominator),
	}, nil
}

func ParseCommitmentProof(proof *ics23.CommitmentProof) (*tendermintContract.IMembershipMsgsCommitmentProof, error) {
	if proof == nil {
		return nil, fmt.Errorf("proof is nil")
	}

	var parsedProof *tendermintContract.IMembershipMsgsCommitmentProof
	switch p := proof.Proof.(type) {
	case *ics23.CommitmentProof_Exist:
		parsedProof = &tendermintContract.IMembershipMsgsCommitmentProof{
			ProofType: ProofType_EXIST,
			ExistenceProof: tendermintContract.IMembershipMsgsExistenceProof{
				Key:   p.Exist.Key,
				Value: bytesToBytes32(p.Exist.Value),
				Leaf:  ParseLeafOp(p.Exist.Leaf),
				Path:  []tendermintContract.IMembershipMsgsInnerOp{},
			},
			NonExistenceProof: tendermintContract.IMembershipMsgsNonExistenceProof{},
		}

		for _, innerOp := range p.Exist.Path {
			parsedProof.ExistenceProof.Path = append(parsedProof.ExistenceProof.Path, ParseInnerOp(innerOp))
		}
	case *ics23.CommitmentProof_Nonexist:
		if p.Nonexist.Left == nil || p.Nonexist.Right == nil {
			return nil, fmt.Errorf("non-existence proof must at least left or right existence proofs")
		}
		parsedProof = &tendermintContract.IMembershipMsgsCommitmentProof{
			ProofType:      ProofType_NON_EXIST,
			ExistenceProof: tendermintContract.IMembershipMsgsExistenceProof{},
			NonExistenceProof: tendermintContract.IMembershipMsgsNonExistenceProof{
				Key:      p.Nonexist.Key,
				HasLeft:  false,
				Left:     tendermintContract.IMembershipMsgsExistenceProof{},
				HasRight: false,
				Right:    tendermintContract.IMembershipMsgsExistenceProof{},
			},
		}

		if p.Nonexist.Left != nil {
			parsedProof.NonExistenceProof.HasLeft = true
			parsedProof.NonExistenceProof.Left = tendermintContract.IMembershipMsgsExistenceProof{
				Key:   p.Nonexist.Left.Key,
				Value: bytesToBytes32(p.Nonexist.Left.Value),
				Leaf:  ParseLeafOp(p.Nonexist.Left.Leaf),
				Path:  []tendermintContract.IMembershipMsgsInnerOp{},
			}
			for _, innerOp := range p.Nonexist.Left.Path {
				parsedProof.NonExistenceProof.Left.Path = append(parsedProof.NonExistenceProof.Left.Path, ParseInnerOp(innerOp))
			}
		}

		if p.Nonexist.Right != nil {
			parsedProof.NonExistenceProof.HasRight = true
			parsedProof.NonExistenceProof.Right = tendermintContract.IMembershipMsgsExistenceProof{
				Key:   p.Nonexist.Right.Key,
				Value: bytesToBytes32(p.Nonexist.Right.Value),
				Leaf:  ParseLeafOp(p.Nonexist.Right.Leaf),
				Path:  []tendermintContract.IMembershipMsgsInnerOp{},
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
			parsedProof = &tendermintContract.IMembershipMsgsCommitmentProof{
				ProofType: ProofType_EXIST,
				ExistenceProof: tendermintContract.IMembershipMsgsExistenceProof{
					Key:   e.Key,
					Value: bytesToBytes32(e.Value),
					Leaf:  ParseLeafOp(e.Leaf),
					Path:  []tendermintContract.IMembershipMsgsInnerOp{},
				},
				NonExistenceProof: tendermintContract.IMembershipMsgsNonExistenceProof{},
			}

			for _, innerOp := range e.Path {
				parsedProof.ExistenceProof.Path = append(parsedProof.ExistenceProof.Path, ParseInnerOp(innerOp))
			}
		}

		if n := p.Batch.GetEntries()[0].GetNonexist(); n != nil {
			parsedProof = &tendermintContract.IMembershipMsgsCommitmentProof{
				ProofType:      ProofType_NON_EXIST,
				ExistenceProof: tendermintContract.IMembershipMsgsExistenceProof{},
				NonExistenceProof: tendermintContract.IMembershipMsgsNonExistenceProof{
					Key:      n.Key,
					HasLeft:  false,
					Left:     tendermintContract.IMembershipMsgsExistenceProof{},
					HasRight: false,
					Right:    tendermintContract.IMembershipMsgsExistenceProof{},
				},
			}

			if n.Left != nil {
				parsedProof.NonExistenceProof.HasLeft = true
				parsedProof.NonExistenceProof.Left = tendermintContract.IMembershipMsgsExistenceProof{
					Key:   n.Left.Key,
					Value: bytesToBytes32(n.Left.Value),
					Leaf:  ParseLeafOp(n.Left.Leaf),
					Path:  []tendermintContract.IMembershipMsgsInnerOp{},
				}
				for _, innerOp := range n.Left.Path {
					parsedProof.NonExistenceProof.Left.Path = append(parsedProof.NonExistenceProof.Left.Path, ParseInnerOp(innerOp))
				}
			}

			if n.Right != nil {
				parsedProof.NonExistenceProof.HasRight = true
				parsedProof.NonExistenceProof.Right = tendermintContract.IMembershipMsgsExistenceProof{
					Key:   n.Right.Key,
					Value: bytesToBytes32(n.Right.Value),
					Leaf:  ParseLeafOp(n.Right.Leaf),
					Path:  []tendermintContract.IMembershipMsgsInnerOp{},
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

func ParseLeafOp(leafOp *ics23.LeafOp) tendermintContract.IMembershipMsgsLeafOp {
	if leafOp == nil {
		return tendermintContract.IMembershipMsgsLeafOp{}
	}

	return tendermintContract.IMembershipMsgsLeafOp{
		HashOp:       uint8(leafOp.Hash),
		PrehashKey:   uint8(leafOp.PrehashKey),
		PrehashValue: uint8(leafOp.PrehashValue),
		Prefix:       leafOp.Prefix,
	}
}

func ParseInnerOp(innerOp *ics23.InnerOp) tendermintContract.IMembershipMsgsInnerOp {
	if innerOp == nil {
		return tendermintContract.IMembershipMsgsInnerOp{}
	}

	return tendermintContract.IMembershipMsgsInnerOp{
		HashOp: uint8(innerOp.Hash),
		Prefix: innerOp.Prefix,
		Suffix: innerOp.Suffix,
	}
}

func EncodeClientState(clientState updateClientContract.IICS07TendermintMsgsClientState) ([]byte, error) {
	args := abi.Arguments{
		{Type: clientStateType},
	}
	encoded, err := args.Pack(clientState)
	return encoded, err
}

func DecodeClientState(data []byte) (updateClientContract.IICS07TendermintMsgsClientState, error) {
	args := abi.Arguments{
		{Type: clientStateType},
	}
	unpacked, err := args.Unpack(data)
	if err != nil {
		return updateClientContract.IICS07TendermintMsgsClientState{}, fmt.Errorf("unpack: %w", err)
	}
	if len(unpacked) == 0 {
		return updateClientContract.IICS07TendermintMsgsClientState{}, fmt.Errorf("no data unpacked")
	}
	// unpacked[0] is an anonymous struct matching the tuple.
	// Use JSON roundtrip to convert to the named target type,
	// since args.Copy maps the tuple to the first field instead of the struct itself.
	jsonBytes, err := json.Marshal(unpacked[0])
	if err != nil {
		return updateClientContract.IICS07TendermintMsgsClientState{}, fmt.Errorf("marshal unpacked tuple: %w", err)
	}
	var clientState updateClientContract.IICS07TendermintMsgsClientState
	if err := json.Unmarshal(jsonBytes, &clientState); err != nil {
		return updateClientContract.IICS07TendermintMsgsClientState{}, fmt.Errorf("unmarshal to client state: %w", err)
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

func EncodeUpdateClientMsg(updateClientMsg updateClientContract.IUpdateClientMsgsMsgUpdateClient) ([]byte, error) {
	args := abi.Arguments{
		{Type: updateClientMsgType},
	}
	encoded, err := args.Pack(updateClientMsg)
	return encoded, err
}
func bytesToBytes32(data []byte) [32]byte {
	var result [32]byte
	copy(result[:], data)
	return result
}
