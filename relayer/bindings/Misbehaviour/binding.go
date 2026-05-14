// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractMisbehaviour

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IICS02ClientMsgsHeight is an auto generated low-level Go binding around an user-defined struct.
type IICS02ClientMsgsHeight struct {
	RevisionNumber uint64
	RevisionHeight uint64
}

// IICS07TendermintMsgsBlockCommit is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsBlockCommit struct {
	Height     uint64
	Round      uint32
	BlockId    IICS07TendermintMsgsBlockId
	CommitSigs []IICS07TendermintMsgsCommitSig
}

// IICS07TendermintMsgsBlockHeader is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsBlockHeader struct {
	Version            IICS07TendermintMsgsVersion
	ChainId            string
	Height             uint64
	Time               *big.Int
	HasLastBlockId     bool
	LastBlockId        IICS07TendermintMsgsBlockId
	HasLastCommitHash  bool
	LastCommitHash     [32]byte
	HasDataHash        bool
	DataHash           [32]byte
	ValidatorsHash     [32]byte
	NextValidatorsHash [32]byte
	ConsensusHash      [32]byte
	AppHash            [32]byte
	HasLastResultsHash bool
	LastResultsHash    [32]byte
	HasEvidenceHash    bool
	EvidenceHash       [32]byte
	ProposerAddress    []byte
}

// IICS07TendermintMsgsBlockId is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsBlockId struct {
	HashData      [32]byte
	PartSetHeader IICS07TendermintMsgsPartSetHeader
}

// IICS07TendermintMsgsChainId is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsChainId struct {
	Id             string
	RevisionNumber uint64
}

// IICS07TendermintMsgsClientState is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsClientState struct {
	ChainId         string
	TrustLevel      IICS07TendermintMsgsTrustThreshold
	LatestHeight    IICS02ClientMsgsHeight
	TrustingPeriod  uint32
	UnbondingPeriod uint32
	IsFrozen        bool
	ZkAlgorithm     uint8
}

// IICS07TendermintMsgsCommitSig is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsCommitSig struct {
	Flag uint8
	Data IICS07TendermintMsgsCommitSigData
}

// IICS07TendermintMsgsCommitSigData is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsCommitSigData struct {
	ValidatorAddress []byte
	Timestamp        *big.Int
	HasSignature     bool
	Signature        []byte
}

// IICS07TendermintMsgsConsensusState is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsConsensusState struct {
	Timestamp          *big.Int
	Root               [32]byte
	NextValidatorsHash [32]byte
}

// IICS07TendermintMsgsHeader is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsHeader struct {
	SignedHeader            IICS07TendermintMsgsSignedHeader
	ValidatorSet            IICS07TendermintMsgsValidatorSet
	TrustedHeight           IICS02ClientMsgsHeight
	TrustedNextValidatorSet IICS07TendermintMsgsValidatorSet
}

// IICS07TendermintMsgsPartSetHeader is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsPartSetHeader struct {
	Total    uint32
	HashData [32]byte
}

// IICS07TendermintMsgsSignedHeader is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsSignedHeader struct {
	Header IICS07TendermintMsgsBlockHeader
	Commit IICS07TendermintMsgsBlockCommit
}

// IICS07TendermintMsgsTrustThreshold is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsTrustThreshold struct {
	Numerator   uint8
	Denominator uint8
}

// IICS07TendermintMsgsValidatorInfo is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsValidatorInfo struct {
	ValAddress       []byte
	PubKey           [32]byte
	VotingPower      uint64
	ProposerPriority int64
}

// IICS07TendermintMsgsValidatorSet is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsValidatorSet struct {
	Validators       []IICS07TendermintMsgsValidatorInfo
	HasProposer      bool
	Proposer         IICS07TendermintMsgsValidatorInfo
	TotalVotingPower uint64
}

// IICS07TendermintMsgsVersion is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsVersion struct {
	BlockVersion uint64
	AppVersion   uint64
}

// IMisbehaviourMsgsMisbehaviour is an auto generated low-level Go binding around an user-defined struct.
type IMisbehaviourMsgsMisbehaviour struct {
	ClientId IICS07TendermintMsgsChainId
	Header1  IICS07TendermintMsgsHeader
	Header2  IICS07TendermintMsgsHeader
}

// IMisbehaviourMsgsMisbehaviourOutput is an auto generated low-level Go binding around an user-defined struct.
type IMisbehaviourMsgsMisbehaviourOutput struct {
	TrustedHeight1 IICS02ClientMsgsHeight
	TrustedHeight2 IICS02ClientMsgsHeight
}

// ContractMisbehaviourMetaData contains all meta data concerning the ContractMisbehaviour contract.
var ContractMisbehaviourMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"misbehaviour_\",\"type\":\"tuple\",\"internalType\":\"structIMisbehaviourMsgs.Misbehaviour\",\"components\":[{\"name\":\"client_id\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ChainId\",\"components\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"header1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"header2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}]},{\"name\":\"trustedConsensusState1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"trustedConsensusState2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIMisbehaviourMsgs.MisbehaviourOutput\",\"components\":[{\"name\":\"trustedHeight1\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight2\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckForMisbehaviourFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidChainId\",\"inputs\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidClientId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MisbehaviourNotDetected\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MisbehaviourVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ValSetHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x608080604052346015576126d3908161001a8239f35b5f80fdfe6080806040526004361015610012575f80fd5b5f3560e01c63a6fe8f5614610025575f80fd5b346106365761012036600319011261063657600435906001600160401b038211610636578136036101206003198201126106365760e082018281106001600160401b0382111761063a576040528260040135906001600160401b0382116106365761009860409260043691870101610712565b83526023190112610636576040516100af8161064e565b6100bb60248401610756565b81526100c960448401610756565b6020820152602082019081526100e23660648501610778565b60408301526101046100f660a485016107ae565b936060840194855261010a60c482016107ae565b608085015261011b60e482016107bf565b60a0850152013560028110156106365760c0830152602435916001600160401b0383116106365760606003198436030112610636576040519361015d85610669565b83600401356001600160401b03811161063657840160406003198236030112610636576040519061018d8261064e565b60048101356001600160401b038111610636576101be916101b660249260043691840101610712565b845201610764565b6020820152855260248401356001600160401b038111610636576101e8906004369187010161099f565b936020860194855260448101356001600160401b038111610636576040916004610215923692010161099f565b95019485526060366043190112610636576040519161023383610669565b6044356001600160801b0381168103610636578352606435602084015260408301926084358452606060a319360112610636576040519361027385610669565b60a4356001600160801b038116810361063657855260c4356020860152604085019360e435855261010435966001600160801b0388168803610636576040516102bb8161064e565b6040516102c78161064e565b5f81525f602082015281526020604051916102e18361064e565b5f83525f82840152015284516040516103188161030a6020820194602086526040830190610d59565b03601f1981018352826106d6565b51902060208a51515101516040516103408161030a6020820194602086526040830190610d59565b5190200361060e57610352895161209d565b61035c8a5161209d565b60208951515101516040516103818161030a6020820194602086526040830190610d59565b51902060208b51515101516040516103a98161030a6020820194602086526040830190610d59565b519020036105a8576104258a8a6001600160401b0360408160206103d281865151510151610df5565b828481846103e5818c5151510151610df5565b94015116975151510151168451966103fc8861064e565b87528287015201511693515151015116604051926104198461064e565b83526020830152612308565b6003811015610594571561054a578760809a976001600160401b039760409761049d6105489c8b998f996001600160801b039960209b61048e63ffffffff6104b39c519351169760405198899461047b86610669565b85528f850152600f604085015251610df5565b9b8c91519451169251936110a6565b856001600160801b038c519451169251936110a6565b01818381835116985151510151168351976104cd8961064e565b88526020880152511692515151015116604051916104ea8361064e565b8252602082015260206040516104ff8161064e565b8481520190815261052860405180946001600160401b0360208092828151168552015116910152565b5160408301906001600160401b0360208092828151168552015116910152565bf35b896001600160401b03604081818d51515101511692515151015116907f4447469a000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b634e487b7160e01b5f52602160045260245ffd5b6105f88a61060a6020808d51515101519251515101516040519384937ff6b6676b000000000000000000000000000000000000000000000000000000008552604060048601526044850190610d59565b83810360031901602485015290610d59565b0390fd5b7fa179f8c9000000000000000000000000000000000000000000000000000000005f5260045ffd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b0382111761063a57604052565b606081019081106001600160401b0382111761063a57604052565b608081019081106001600160401b0382111761063a57604052565b61026081019081106001600160401b0382111761063a57604052565b60a081019081106001600160401b0382111761063a57604052565b90601f801991011681019081106001600160401b0382111761063a57604052565b6001600160401b03811161063a57601f01601f191660200190565b81601f820112156106365760208135910161072c826106f7565b9261073a60405194856106d6565b8284528282011161063657815f92602092838601378301015290565b359060ff8216820361063657565b35906001600160401b038216820361063657565b9190826040910312610636576040516107908161064e565b60206107a98183956107a181610764565b855201610764565b910152565b359063ffffffff8216820361063657565b3590811515820361063657565b35906001600160801b038216820361063657565b809291039160608312610636576040516107f98161064e565b6040819483358352601f19011261063657602090604080519361081b8561064e565b6108268482016107ae565b85520135828401520152565b6001600160401b03811161063a5760051b60200190565b9190608083820312610636576040519061086282610684565b819380356001600160401b03811161063657606092610882918301610712565b83526020810135602084015261089a60408201610764565b60408401520135908160070b82036106365760600152565b91909160808184031261063657604051906108cc82610684565b819381356001600160401b03811161063657820181601f820112156106365780356108f681610832565b9161090460405193846106d6565b81835260208084019260051b820101918483116106365760208201905b8382106109725750505050835261093a602083016107bf565b60208401526040820135906001600160401b0382116106365782610967606094926107a994869401610849565b604086015201610764565b81356001600160401b0381116106365760209161099488848094880101610849565b815201910190610921565b919060a08382031261063657604051906109b882610684565b819380356001600160401b03811161063657810160408184031261063657604051906109e38261064e565b80356001600160401b0381116106365781016102c0818603126106365760405190610a0d8261069f565b610a178682610778565b825260408101356001600160401b0381116106365786610a38918301610712565b6020830152610a4960608201610764565b6040830152610a5a608082016107cc565b6060830152610a6b60a082016107bf565b6080830152610a7d8660c083016107e0565b60a0830152610a8f61012082016107bf565b60c083015261014081013560e0830152610aac61016082016107bf565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a0830152610afb61022082016107bf565b6101c08301526102408101356101e0830152610b1a61026082016107bf565b6102008301526102808101356102208301526102a0810135906001600160401b03821161063657610b4d91879101610712565b61024082015282526020810135906001600160401b038211610636570160c0818503126106365760405190610b8182610684565b610b8a81610764565b8252610b98602082016107ae565b6020830152610baa85604083016107e0565b604083015260a0810135906001600160401b038211610636570184601f8201121561063657803590610bdb82610832565b91610be960405193846106d6565b80835260208084019160051b830101918783116106365760208101915b838310610c74575050505060608201526020820152835260208101356001600160401b0381116106365782610c3c9183016108b2565b6020840152610c4e8260408301610778565b60408401526080810135916001600160401b038311610636576060926107a992016108b2565b82356001600160401b038111610636578201906040828b03601f1901126106365760405191610ca28361064e565b6020810135600481101561063657835260408101356001600160401b038111610636576020910101906080828c03126106365760405192610ce284610684565b82356001600160401b038111610636578c610cfe918501610712565b8452610d0c602084016107cc565b6020850152610d1d604084016107bf565b60408501526060830135936001600160401b03851161063657610d458d602096879601610712565b606082015283820152815201920191610c06565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b908151811015610d8e570160200190565b634e487b7160e01b5f52603260045260245ffd5b90610dac826106f7565b610db960405191826106d6565b8281528092610dca601f19916106f7565b0190602036910137565b91908201809211610de157565b634e487b7160e01b5f52601160045260245ffd5b604051610e018161064e565b606081525f60209091015280515f198101908111610de1575b602d60f81b6001600160f81b0319610e328385610d7d565b511614610e47578015610de1575f1901610e1a565b91905f198314611063578051838103908111610de1575f198101908111610de157610e7190610da2565b905f5b8251811015610ebc576001850190818611610de1576001600160f81b0319610ea7610ea183600195610dd4565b85610d7d565b51165f1a610eb58286610d7d565b5301610e74565b509192908051158015611039575b610fab578051610ed991612371565b9190158015611029575b610fab57610ef081610da2565b905f5b818110610f80575050516001811190811591610f74575b50610f30576001600160401b039060405192610f258461064e565b835216602082015290565b606460405162461bcd60e51b815260206004820152601b60248201527f496e76616c696420636861696e20707265666978206c656e67746800000000006044820152fd5b602b915010155f610f0a565b806001600160f81b0319610f9660019388610d7d565b51165f1a610fa48286610d7d565b5301610ef3565b505080518060011115908161101e575b5015610fd95760405190610fce8261064e565b81525f602082015290565b60405162461bcd60e51b815260206004820152601760248201527f496e76616c696420636861696e204944206c656e6774680000000000000000006044820152606490fd5b60409150105f610fbb565b506001600160401b038211610ee3565b610d8e57600360fc1b6001600160f81b0319602083015116148015610eca57506001815111610eca565b80919250518060011115908161101e575015610fd95760405190610fce8261064e565b906001600160801b03809116911603906001600160801b038211610de157565b919395909594929473__$7440a880b7578767f72184d998805816e4$__9060608401966110ea602089516040518093819263c87f1f6960e01b835260048301611fff565b0381875af48015611d905783915f91611f8a575b5003611ee0576001600160801b03811690816001600160801b03881610611eb4576111299087611086565b6001600160801b0363ffffffff60208a015116911681811015611e8657505088519788518015908115611e7b575b50611e37575f5b8951811015611345576001600160f81b031961117a828c610d7d565b51167f6100000000000000000000000000000000000000000000000000000000000000811015908161131a575b81156112be575b811561127e575b8115611270575b8115611246575b811561121c575b50156111d85760010161115e565b606460405162461bcd60e51b815260206004820152601860248201527f696e76616c696420636861696e206964206368617273657400000000000000006044820152fd5b7f2e000000000000000000000000000000000000000000000000000000000000009150145f6111ca565b7f5f00000000000000000000000000000000000000000000000000000000000000811491506111c3565b602d60f81b811491506111bc565b9050600360fc1b81101580611294575b906111b5565b507f390000000000000000000000000000000000000000000000000000000000000081111561128e565b90507f4100000000000000000000000000000000000000000000000000000000000000811015806112f0575b906111ae565b507f5a000000000000000000000000000000000000000000000000000000000000008111156112ea565b7f7a0000000000000000000000000000000000000000000000000000000000000081111591506111a7565b50909192939495969861145d9850602060406001600160401b0392815161136b8161064e565b82516113768161064e565b83516113818161069f565b845161138c8161064e565b5f81525f8782015281526060868201525f858201525f60608201525f60808201526113b561225c565b60a08201525f60c08201525f60e08201525f6101008201525f6101208201525f6101408201525f6101608201525f6101808201525f6101a08201525f6101c08201525f6101e08201525f6102008201525f61022082015260606102408201528152835161142181610684565b5f81525f8682015261143161225c565b858201526060808201528582015281528361144a612287565b9101528951838b01519083519d8e61064e565b8d52838d015251985f60808351611473816106bb565b6060815282868201528285820152611489612287565b606082015201520151015116905191604051966114a5886106bb565b875260208701526040860152606085015260808401526114de6020808701516040518093819263c87f1f6960e01b835260048301611fff565b0381855af4908115611d90575f91611e05575b50610140865151015103611d9b57602061167d918651519060405180809581947f95b25f7900000000000000000000000000000000000000000000000000000000835286600484015261155d6024840182516001600160401b0360208092828151168552015116910152565b61024061157a888301516102c060648701526102e4860190610d59565b60408301516001600160401b0316608486015260608301516001600160801b031660a48601526080830151151560c486015260a0830151805160e4870152890151805163ffffffff1661010487015289015161012486015260c0830151151561014486015260e083015161016486015261010083015115156101848601526101208301516101a48601526101408301516101c48601526101608301516101e48601526101808301516102048601526101a08301516102248601526101c083015115156102448601526101e083015161026486015261020083015115156102848601526102208301516102a4860152910151838203602319016102c4850152610d59565b03915af4908115611d90575f91611d5e575b50604060208651015101515103611cf45783519160606020808701519401510151945f965f5b87518110156116f3576116d16116cb828a6123e3565b516123f7565b6004811015610594576001036116ea575b6001016116b5565b600198506116e2565b509091929395949615611c8a57845186515103611bfa575f5b85518110156118425761171f81876123e3565b51805160048110156105945760010361173d57506001905b0161170c565b602090989395919896929496015151955f955f5b895180518210156118325781611766916123e3565b51516040516117946020828180820195805191829101875e81015f838201520301601f1981018352826106d6565b519020896040516117c46020828180820195805191829101875e81015f838201520301601f1981018352826106d6565b519020146117d457600101611751565b5093989195509391955060015b156117ee57600190611737565b606460405162461bcd60e51b815260206004820152601d60248201527f696e76616c696420636f6d6d69743a206661756c7479207369676e65720000006044820152fd5b50509398919590949296506117e1565b5092959194509250633b9aca0063ffffffff602084015116026001600160801b038116908103610de1576001600160801b03602086015116806001600160801b03841610928315611bdb575b505050611b71576001600160801b0360608351510151166001600160801b036020850151161015611b0757602082515101516040516118ec6020828180820195805191829101875e81015f838201520301601f1981018352826106d6565b519020835160405161191d6020828180820195805191829101875e81015f838201520301601f1981018352826106d6565b51902003611ac35761193b6001600160401b036040850151166122ca565b825151604001516001600160401b039182169116818103611a5557505061014082515101516080840151036119eb575b6119816001600160401b036040850151166122ca565b825151604001516001600160401b039182169116146119d8576119d6926119b19160608451920151905191612427565b604051906119be8261064e565b60028252600360208301526020815191015190612427565b565b506119d69150604051906119be8261064e565b608460405162461bcd60e51b815260206004820152602f60248201527f696e76616c696420626c6f636b3a206e6578742076616c696461746f7220736560448201527f742068617368206d69736d6174636800000000000000000000000000000000006064820152fd5b1161196b57608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152fd5b6001600160801b0392935090611bf091611086565b16115f808061188e565b60a460405162461bcd60e51b815260206004820152604860248201527f696e76616c696420636f6d6d69743a206e756d626572206f66207369676e617460448201527f7572657320646f6573206e6f74206d61746368206e756d626572206f6620766160648201527f6c696461746f72730000000000000000000000000000000000000000000000006084820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420636f6d6d69743a206e6f2070726573656e74207369676e6160448201527f74757265730000000000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152fd5b90506020813d602011611d88575b81611d79602093836106d6565b8101031261063657515f61168f565b3d9150611d6c565b6040513d5f823e3d90fd5b608460405162461bcd60e51b815260206004820152602a60248201527f696e76616c696420626c6f636b3a2076616c696461746f72207365742068617360448201527f68206d69736d61746368000000000000000000000000000000000000000000006064820152fd5b90506020813d602011611e2f575b81611e20602093836106d6565b8101031261063657515f6114f1565b3d9150611e13565b606460405162461bcd60e51b815260206004820152601760248201527f496e76616c696420636861696e206964206c656e6774680000000000000000006044820152fd5b60329150115f611157565b7fdb020436000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b507fe42a1980000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b60a46040517ff492ef2b00000000000000000000000000000000000000000000000000000000815260206004820152604360248201527f74727573746564206e6578742076616c696461746f722073657420686173682060448201527f646f6573206e6f74206d6174636820686173682073746f726564206f6e20636860648201527f61696e00000000000000000000000000000000000000000000000000000000006084820152fd5b9150506020813d602011611fb7575b81611fa6602093836106d6565b81010312610636578290515f6110fe565b3d9150611f99565b90606080611fd68451608085526080850190610d59565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b6020815260a0810182519060806020840152815180915260c0830190602060c08260051b8601019301915f905b82821061207257505050506001600160401b0360606120686080936020870151151560408701526040870151601f198783030184880152611fbf565b9401511691015290565b9091929360208061208f60019360bf198a82030186528851611fbf565b96019201920190929161202c565b60206120ad818351510151610df5565b016001600160401b038151169060408301916001600160401b03835151169081810361222e5750506001600160401b036121109151166001600160401b03604085515101511692604051916121018361064e565b82526020820193845251612308565b60038110156105945760028114908115612223575b506121ed57508060208061214f9301516040518094819263c87f1f6960e01b835260048301611fff565b038173__$7440a880b7578767f72184d998805816e4$__5af4918215611d90575f926121b7575b5051516101400151808203612189575050565b7f4d6d8014000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b9091506020813d6020116121e5575b816121d3602093836106d6565b81010312610636575190610140612176565b3d91506121c6565b6001600160401b039051167f90f4dbed000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b60019150145f612125565b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b604051906122698261064e565b5f825260405160208361227b8361064e565b5f83525f828401520152565b6040519061229482610684565b5f6060838181528260208201526040516122ad81610684565b828152836020820152836040820152838382015260408201520152565b6001600160401b036001911601906001600160401b038211610de157565b906001600160401b03809116911601906001600160401b038211610de157565b6001600160401b038151166001600160401b03835116908181105f1461233057505050505f90565b111561233d575050600290565b60206001600160401b0381819301511692015116908181105f146123615750505f90565b111561236c57600290565b600190565b5f929183915b8183106123875750505060019190565b90919360ff6123a56001600160f81b03196020888601015116612654565b1690600982116123d757600a810290808204600a1490151715610de1576001916123ce91610dd4565b94019190612377565b5050505090505f905f90565b8051821015610d8e5760209160051b010190565b5160048110156105945790565b906001600160401b03809116911602906001600160401b038216918203610de157565b6020015160600151925f9291835b8151805186101561246c576001916001600160401b03604061245a89612464956123e3565b51015116906122e8565b940193612435565b509092919493505f5f935b85518510156125b05761248d6116cb86886123e3565b6004811015610594576001146125a65760206124a986886123e3565b510151519560208701905f5b8351805182101561259a57816124ca916123e3565b51516040516124f86020828180820195805191829101875e81015f838201520301601f1981018352826106d6565b519020896040516125246020828181019451808a875e81015f838201520301601f1981018352826106d6565b51902014612534576001016124b5565b8398506001600160401b03919794925061245a60409161255495516123e3565b905b61256760ff60208601511683612404565b6001600160401b038061257e60ff88511687612404565b16911611612592576001905b019394612477565b505050505050565b50509591965050612556565b949360019061258a565b506001600160401b039350839294509060ff6125d66125df938260208901511690612404565b95511690612404565b16911611156125ea57565b608460405162461bcd60e51b815260206004820152602160248201527f696e73756666696369656e7420766f74696e6720706f776572206f7665726c6160448201527f70000000000000000000000000000000000000000000000000000000000000006064820152fd5b60f81c602f8111806126bc575b1561267057602f190160ff1690565b60608111806126b2575b15612689576056190160ff1690565b60408111806126a8575b156126a2576036190160ff1690565b5060ff90565b5060478110612693565b506067811061267a565b50603a811061266156fea164736f6c634300081c000a",
}

// ContractMisbehaviourABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMisbehaviourMetaData.ABI instead.
var ContractMisbehaviourABI = ContractMisbehaviourMetaData.ABI

// ContractMisbehaviourBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMisbehaviourMetaData.Bin instead.
var ContractMisbehaviourBin = ContractMisbehaviourMetaData.Bin

// DeployContractMisbehaviour deploys a new Ethereum contract, binding an instance of ContractMisbehaviour to it.
func DeployContractMisbehaviour(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractMisbehaviour, error) {
	parsed, err := ContractMisbehaviourMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractMisbehaviourBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractMisbehaviour{ContractMisbehaviourCaller: ContractMisbehaviourCaller{contract: contract}, ContractMisbehaviourTransactor: ContractMisbehaviourTransactor{contract: contract}, ContractMisbehaviourFilterer: ContractMisbehaviourFilterer{contract: contract}}, nil
}

// ContractMisbehaviour is an auto generated Go binding around an Ethereum contract.
type ContractMisbehaviour struct {
	ContractMisbehaviourCaller     // Read-only binding to the contract
	ContractMisbehaviourTransactor // Write-only binding to the contract
	ContractMisbehaviourFilterer   // Log filterer for contract events
}

// ContractMisbehaviourCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractMisbehaviourCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMisbehaviourTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractMisbehaviourTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMisbehaviourFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractMisbehaviourFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMisbehaviourSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractMisbehaviourSession struct {
	Contract     *ContractMisbehaviour // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ContractMisbehaviourCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractMisbehaviourCallerSession struct {
	Contract *ContractMisbehaviourCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// ContractMisbehaviourTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractMisbehaviourTransactorSession struct {
	Contract     *ContractMisbehaviourTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractMisbehaviourRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractMisbehaviourRaw struct {
	Contract *ContractMisbehaviour // Generic contract binding to access the raw methods on
}

// ContractMisbehaviourCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractMisbehaviourCallerRaw struct {
	Contract *ContractMisbehaviourCaller // Generic read-only contract binding to access the raw methods on
}

// ContractMisbehaviourTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractMisbehaviourTransactorRaw struct {
	Contract *ContractMisbehaviourTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractMisbehaviour creates a new instance of ContractMisbehaviour, bound to a specific deployed contract.
func NewContractMisbehaviour(address common.Address, backend bind.ContractBackend) (*ContractMisbehaviour, error) {
	contract, err := bindContractMisbehaviour(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractMisbehaviour{ContractMisbehaviourCaller: ContractMisbehaviourCaller{contract: contract}, ContractMisbehaviourTransactor: ContractMisbehaviourTransactor{contract: contract}, ContractMisbehaviourFilterer: ContractMisbehaviourFilterer{contract: contract}}, nil
}

// NewContractMisbehaviourCaller creates a new read-only instance of ContractMisbehaviour, bound to a specific deployed contract.
func NewContractMisbehaviourCaller(address common.Address, caller bind.ContractCaller) (*ContractMisbehaviourCaller, error) {
	contract, err := bindContractMisbehaviour(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractMisbehaviourCaller{contract: contract}, nil
}

// NewContractMisbehaviourTransactor creates a new write-only instance of ContractMisbehaviour, bound to a specific deployed contract.
func NewContractMisbehaviourTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractMisbehaviourTransactor, error) {
	contract, err := bindContractMisbehaviour(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractMisbehaviourTransactor{contract: contract}, nil
}

// NewContractMisbehaviourFilterer creates a new log filterer instance of ContractMisbehaviour, bound to a specific deployed contract.
func NewContractMisbehaviourFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractMisbehaviourFilterer, error) {
	contract, err := bindContractMisbehaviour(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractMisbehaviourFilterer{contract: contract}, nil
}

// bindContractMisbehaviour binds a generic wrapper to an already deployed contract.
func bindContractMisbehaviour(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMisbehaviourMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractMisbehaviour *ContractMisbehaviourRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractMisbehaviour.Contract.ContractMisbehaviourCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractMisbehaviour *ContractMisbehaviourRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.ContractMisbehaviourTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractMisbehaviour *ContractMisbehaviourRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.ContractMisbehaviourTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractMisbehaviour *ContractMisbehaviourCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractMisbehaviour.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractMisbehaviour *ContractMisbehaviourTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractMisbehaviour *ContractMisbehaviourTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.contract.Transact(opts, method, params...)
}

// Misbehaviour is a free data retrieval call binding the contract method 0xa6fe8f56.
//
// Solidity: function misbehaviour((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8) clientState, ((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64))) misbehaviour_, (uint128,bytes32,bytes32) trustedConsensusState1, (uint128,bytes32,bytes32) trustedConsensusState2, uint128 time) pure returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourCaller) Misbehaviour(opts *bind.CallOpts, clientState IICS07TendermintMsgsClientState, misbehaviour_ IMisbehaviourMsgsMisbehaviour, trustedConsensusState1 IICS07TendermintMsgsConsensusState, trustedConsensusState2 IICS07TendermintMsgsConsensusState, time *big.Int) (IMisbehaviourMsgsMisbehaviourOutput, error) {
	var out []interface{}
	err := _ContractMisbehaviour.contract.Call(opts, &out, "misbehaviour", clientState, misbehaviour_, trustedConsensusState1, trustedConsensusState2, time)

	if err != nil {
		return *new(IMisbehaviourMsgsMisbehaviourOutput), err
	}

	out0 := *abi.ConvertType(out[0], new(IMisbehaviourMsgsMisbehaviourOutput)).(*IMisbehaviourMsgsMisbehaviourOutput)

	return out0, err

}

// Misbehaviour is a free data retrieval call binding the contract method 0xa6fe8f56.
//
// Solidity: function misbehaviour((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8) clientState, ((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64))) misbehaviour_, (uint128,bytes32,bytes32) trustedConsensusState1, (uint128,bytes32,bytes32) trustedConsensusState2, uint128 time) pure returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourSession) Misbehaviour(clientState IICS07TendermintMsgsClientState, misbehaviour_ IMisbehaviourMsgsMisbehaviour, trustedConsensusState1 IICS07TendermintMsgsConsensusState, trustedConsensusState2 IICS07TendermintMsgsConsensusState, time *big.Int) (IMisbehaviourMsgsMisbehaviourOutput, error) {
	return _ContractMisbehaviour.Contract.Misbehaviour(&_ContractMisbehaviour.CallOpts, clientState, misbehaviour_, trustedConsensusState1, trustedConsensusState2, time)
}

// Misbehaviour is a free data retrieval call binding the contract method 0xa6fe8f56.
//
// Solidity: function misbehaviour((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8) clientState, ((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64))) misbehaviour_, (uint128,bytes32,bytes32) trustedConsensusState1, (uint128,bytes32,bytes32) trustedConsensusState2, uint128 time) pure returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourCallerSession) Misbehaviour(clientState IICS07TendermintMsgsClientState, misbehaviour_ IMisbehaviourMsgsMisbehaviour, trustedConsensusState1 IICS07TendermintMsgsConsensusState, trustedConsensusState2 IICS07TendermintMsgsConsensusState, time *big.Int) (IMisbehaviourMsgsMisbehaviourOutput, error) {
	return _ContractMisbehaviour.Contract.Misbehaviour(&_ContractMisbehaviour.CallOpts, clientState, misbehaviour_, trustedConsensusState1, trustedConsensusState2, time)
}
