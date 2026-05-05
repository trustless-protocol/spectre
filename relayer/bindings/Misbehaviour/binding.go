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
	Bin: "0x60808060405234601557612aee908161001a8239f35b5f80fdfe6080806040526004361015610012575f80fd5b5f3560e01c63a6fe8f5614610025575f80fd5b34610779576101207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610779576004359067ffffffffffffffff8211610779578136036101207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8201126107795760e0820182811067ffffffffffffffff82111761077d57604052826004013567ffffffffffffffff8111610779576040916100f87fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc9260043691880101610876565b845201126107795760405161010c816107aa565b610118602484016108ba565b8152610126604484016108ba565b60208201526020820190815261013f36606485016108dd565b604083015261010461015360a48501610913565b936060840194855261016760c48201610913565b608085015261017860e48201610924565b60a0850152013560028110156107795760c08301526024359167ffffffffffffffff83116107795760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc843603011261077957604051936101d9856107c6565b836004013567ffffffffffffffff811161077957840160407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc82360301126107795760405190610228826107aa565b600481013567ffffffffffffffff81116107795761025a9161025260249260043691840101610876565b8452016108c8565b60208201528552602484013567ffffffffffffffff8111610779576102859060043691870101610b12565b9360208601948552604481013567ffffffffffffffff81116107795760409160046102b39236920101610b12565b950194855260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffbc36011261077957604051916102ef836107c6565b6044356fffffffffffffffffffffffffffffffff8116810361077957835260643560208401526040830192608435845260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff5c3601126107795760405193610356856107c6565b60a4356fffffffffffffffffffffffffffffffff8116810361077957855260c4356020860152604085019360e435855261010435966fffffffffffffffffffffffffffffffff88168803610779576040516103b0816107aa565b6040516103bc816107aa565b5f81525f602082015281526020604051916103d6836107aa565b5f83525f828401520152845160405161040d816103ff6020820194602086526040830190610ed8565b03601f198101835282610837565b51902060208a5151510151604051610435816103ff6020820194602086526040830190610ed8565b5190200361075157610447895161241b565b6104518a5161241b565b6020895151510151604051610476816103ff6020820194602086526040830190610ed8565b51902060208b515151015160405161049e816103ff6020820194602086526040830190610ed8565b519020036106cd5761051b8a8a67ffffffffffffffff60408160206104c881865151510151610fa6565b828481846104db818c5151510151610fa6565b94015116975151510151168451966104f2886107aa565b875282870152015116935151510151166040519261050f846107aa565b835260208301526126a8565b60038110156106a05715610655578760809a9767ffffffffffffffff9760409761059d6106539c8b998f996fffffffffffffffffffffffffffffffff9960209b61058e63ffffffff6105bc9c519351169760405198899461057b866107c6565b85528f850152600f604085015251610fa6565b9b8c9151945116925193611303565b856fffffffffffffffffffffffffffffffff8c51945116925193611303565b01818381835116985151510151168351976105d6896107aa565b88526020880152511692515151015116604051916105f3836107aa565b825260208201526020604051610608816107aa565b84815201908152610632604051809467ffffffffffffffff60208092828151168552015116910152565b51604083019067ffffffffffffffff60208092828151168552015116910152565bf35b8967ffffffffffffffff604081818d51515101511692515151015116907f4447469a000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b61071d8a61074d6020808d51515101519251515101516040519384937ff6b6676b000000000000000000000000000000000000000000000000000000008552604060048601526044850190610ed8565b907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc848303016024850152610ed8565b0390fd5b7fa179f8c9000000000000000000000000000000000000000000000000000000005f5260045ffd5b5f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6040810190811067ffffffffffffffff82111761077d57604052565b6060810190811067ffffffffffffffff82111761077d57604052565b6080810190811067ffffffffffffffff82111761077d57604052565b610260810190811067ffffffffffffffff82111761077d57604052565b60a0810190811067ffffffffffffffff82111761077d57604052565b90601f601f19910116810190811067ffffffffffffffff82111761077d57604052565b67ffffffffffffffff811161077d57601f01601f191660200190565b81601f82011215610779576020813591016108908261085a565b9261089e6040519485610837565b8284528282011161077957815f92602092838601378301015290565b359060ff8216820361077957565b359067ffffffffffffffff8216820361077957565b9190826040910312610779576040516108f5816107aa565b602061090e818395610906816108c8565b8552016108c8565b910152565b359063ffffffff8216820361077957565b3590811515820361077957565b35906fffffffffffffffffffffffffffffffff8216820361077957565b80929103916060831261077957604051610967816107aa565b6040601f198295843584520112610779576020906040805193610989856107aa565b610994848201610913565b85520135828401520152565b67ffffffffffffffff811161077d5760051b60200190565b919060808382031261077957604051906109d1826107e2565b8193803567ffffffffffffffff8111610779576060926109f2918301610876565b835260208101356020840152610a0a604082016108c8565b60408401520135908160070b82036107795760600152565b9190916080818403126107795760405190610a3c826107e2565b8193813567ffffffffffffffff811161077957820181601f82011215610779578035610a67816109a0565b91610a756040519384610837565b81835260208084019260051b820101918483116107795760208201905b838210610ae457505050508352610aab60208301610924565b602084015260408201359067ffffffffffffffff82116107795782610ad96060949261090e948694016109b8565b6040860152016108c8565b813567ffffffffffffffff811161077957602091610b07888480948801016109b8565b815201910190610a92565b919060a0838203126107795760405190610b2b826107e2565b8193803567ffffffffffffffff81116107795781016040818403126107795760405190610b57826107aa565b803567ffffffffffffffff81116107795781016102c0818603126107795760405190610b82826107fe565b610b8c86826108dd565b8252604081013567ffffffffffffffff81116107795786610bae918301610876565b6020830152610bbf606082016108c8565b6040830152610bd060808201610931565b6060830152610be160a08201610924565b6080830152610bf38660c0830161094e565b60a0830152610c056101208201610924565b60c083015261014081013560e0830152610c226101608201610924565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a0830152610c716102208201610924565b6101c08301526102408101356101e0830152610c906102608201610924565b6102008301526102808101356102208301526102a08101359067ffffffffffffffff821161077957610cc491879101610876565b610240820152825260208101359067ffffffffffffffff8211610779570160c0818503126107795760405190610cf9826107e2565b610d02816108c8565b8252610d1060208201610913565b6020830152610d22856040830161094e565b604083015260a08101359067ffffffffffffffff8211610779570184601f8201121561077957803590610d54826109a0565b91610d626040519384610837565b80835260208084019160051b830101918783116107795760208101915b838310610def5750505050606082015260208201528352602081013567ffffffffffffffff81116107795782610db6918301610a22565b6020840152610dc882604083016108dd565b604084015260808101359167ffffffffffffffff83116107795760609261090e9201610a22565b823567ffffffffffffffff8111610779578201906040601f19838c0301126107795760405191610e1e836107aa565b60208101356004811015610779578352604081013567ffffffffffffffff8111610779576020910101906080828c03126107795760405192610e5f846107e2565b823567ffffffffffffffff8111610779578c610e7c918501610876565b8452610e8a60208401610931565b6020850152610e9b60408401610924565b604085015260608301359367ffffffffffffffff851161077957610ec48d602096879601610876565b606082015283820152815201920191610d7f565b90601f19601f602080948051918291828752018686015e5f8582860101520116010190565b908151811015610f0e570160200190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b90610f458261085a565b610f526040519182610837565b828152601f19610f62829461085a565b0190602036910137565b91908201809211610f7957565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b604051610fb2816107aa565b606081525f60209091015280515f198101908111610f79575b7f2d000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000006110178385610efd565b51161461102c578015610f79575f1901610fcb565b91905f1983146112ae578051838103908111610f79575f198101908111610f795761105690610f3b565b905f5b82518110156110b9576001850190818611610f79577fff000000000000000000000000000000000000000000000000000000000000006110a461109e83600195610f6c565b85610efd565b51165f1a6110b28286610efd565b5301611059565b509192908051158015611250575b6111c15780516110d691612714565b919015801561123f575b6111c1576110ed81610f3b565b905f5b81811061117e575050516001811190811591611172575b5061112e5767ffffffffffffffff9060405192611123846107aa565b835216602082015290565b606460405162461bcd60e51b815260206004820152601b60248201527f496e76616c696420636861696e20707265666978206c656e67746800000000006044820152fd5b602b915010155f611107565b807fff000000000000000000000000000000000000000000000000000000000000006111ac60019388610efd565b51165f1a6111ba8286610efd565b53016110f0565b5050805180600111159081611234575b50156111ef57604051906111e4826107aa565b81525f602082015290565b60405162461bcd60e51b815260206004820152601760248201527f496e76616c696420636861696e204944206c656e6774680000000000000000006044820152606490fd5b60409150105f6111d1565b5067ffffffffffffffff82116110e0565b610f0e577f30000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000006020830151161480156110c7575060018151116110c7565b8091925051806001111590816112345750156111ef57604051906111e4826107aa565b906fffffffffffffffffffffffffffffffff809116911603906fffffffffffffffffffffffffffffffff8211610f7957565b919395909594929473__$7440a880b7578767f72184d998805816e4$__90606084019661136060208951604051809381927fc87f1f690000000000000000000000000000000000000000000000000000000083526004830161235e565b0381875af480156120ee5783915f916122e8575b500361223e576fffffffffffffffffffffffffffffffff811690816fffffffffffffffffffffffffffffffff881610612212576113b190876112d1565b6fffffffffffffffffffffffffffffffff63ffffffff60208a0151169116818110156121e4575050885197885180159081156121d9575b50612195575f5b8951811015611626577fff00000000000000000000000000000000000000000000000000000000000000611423828c610efd565b51167f610000000000000000000000000000000000000000000000000000000000000081101590816115fb575b811561159f575b8115611543575b8115611519575b81156114ef575b81156114c5575b5015611481576001016113ef565b606460405162461bcd60e51b815260206004820152601860248201527f696e76616c696420636861696e206964206368617273657400000000000000006044820152fd5b7f2e000000000000000000000000000000000000000000000000000000000000009150145f611473565b7f5f000000000000000000000000000000000000000000000000000000000000008114915061146c565b7f2d0000000000000000000000000000000000000000000000000000000000000081149150611465565b90507f300000000000000000000000000000000000000000000000000000000000000081101580611575575b9061145e565b507f390000000000000000000000000000000000000000000000000000000000000081111561156f565b90507f4100000000000000000000000000000000000000000000000000000000000000811015806115d1575b90611457565b507f5a000000000000000000000000000000000000000000000000000000000000008111156115cb565b7f7a000000000000000000000000000000000000000000000000000000000000008111159150611450565b50909192939495969861173f98506020604067ffffffffffffffff92815161164d816107aa565b8251611658816107aa565b8351611663816107fe565b845161166e816107aa565b5f81525f8782015281526060868201525f858201525f60608201525f60808201526116976125f8565b60a08201525f60c08201525f60e08201525f6101008201525f6101208201525f6101408201525f6101608201525f6101808201525f6101a08201525f6101c08201525f6101e08201525f6102008201525f610220820152606061024082015281528351611703816107e2565b5f81525f868201526117136125f8565b858201526060808201528582015281528361172c612623565b9101528951838b01519083519d8e6107aa565b8d52838d015251985f608083516117558161081b565b606081528286820152828582015261176b612623565b606082015201520151015116905191604051966117878861081b565b875260208701526040860152606085015260808401526117d9602080870151604051809381927fc87f1f690000000000000000000000000000000000000000000000000000000083526004830161235e565b0381855af49081156120ee575f91612163575b506101408651510151036120f95760206119a1918651519060405180809581947f95b25f7900000000000000000000000000000000000000000000000000000000835286600484015261185960248401825167ffffffffffffffff60208092828151168552015116910152565b610240611876888301516102c060648701526102e4860190610ed8565b9167ffffffffffffffff60408201511660848601526fffffffffffffffffffffffffffffffff60608201511660a48601526080810151151560c4860152888060a0830151805160e4890152015163ffffffff815116610104880152015161012486015260c0810151151561014486015260e081015161016486015261010081015115156101848601526101208101516101a48601526101408101516101c48601526101608101516101e48601526101808101516102048601526101a08101516102248601526101c081015115156102448601526101e081015161026486015261020081015115156102848601526102208101516102a486015201517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc848303016102c4850152610ed8565b03915af49081156120ee575f916120bc575b506040602086510151015151036120525783519160606020808701519401510151945f965f5b8751811015611a17576119f56119ef828a61279e565b516127b2565b60048110156106a057600103611a0e575b6001016119d9565b60019850611a06565b509091929395949615611fe857845186515103611f58575f5b8551811015611b6657611a43818761279e565b51805160048110156106a057600103611a6157506001905b01611a30565b602090989395919896929496015151955f955f5b89518051821015611b565781611a8a9161279e565b5151604051611ab86020828180820195805191829101875e81015f838201520301601f198101835282610837565b51902089604051611ae86020828180820195805191829101875e81015f838201520301601f198101835282610837565b51902014611af857600101611a75565b5093989195509391955060015b15611b1257600190611a5b565b606460405162461bcd60e51b815260206004820152601d60248201527f696e76616c696420636f6d6d69743a206661756c7479207369676e65720000006044820152fd5b5050939891959094929650611b05565b5092959194509250633b9aca0063ffffffff602084015116026fffffffffffffffffffffffffffffffff8116908103610f79576fffffffffffffffffffffffffffffffff602086015116806fffffffffffffffffffffffffffffffff841610928315611f30575b505050611ec6576fffffffffffffffffffffffffffffffff60608351510151166fffffffffffffffffffffffffffffffff6020850151161015611e5c5760208251510151604051611c3d6020828180820195805191829101875e81015f838201520301601f198101835282610837565b5190208351604051611c6e6020828180820195805191829101875e81015f838201520301601f198101835282610837565b51902003611e1857611c8d67ffffffffffffffff604085015116612666565b8251516040015167ffffffffffffffff9182169116818103611daa5750506101408251510151608084015103611d40575b611cd567ffffffffffffffff604085015116612666565b8251516040015167ffffffffffffffff918216911614611d2d57611d2b92611d0691606084519201519051916127e4565b60405190611d13826107aa565b600282526003602083015260208151910151906127e4565b565b50611d2b915060405190611d13826107aa565b608460405162461bcd60e51b815260206004820152602f60248201527f696e76616c696420626c6f636b3a206e6578742076616c696461746f7220736560448201527f742068617368206d69736d6174636800000000000000000000000000000000006064820152fd5b11611cbe57608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152fd5b6fffffffffffffffffffffffffffffffff92935090611f4e916112d1565b16115f8080611bcd565b60a460405162461bcd60e51b815260206004820152604860248201527f696e76616c696420636f6d6d69743a206e756d626572206f66207369676e617460448201527f7572657320646f6573206e6f74206d61746368206e756d626572206f6620766160648201527f6c696461746f72730000000000000000000000000000000000000000000000006084820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420636f6d6d69743a206e6f2070726573656e74207369676e6160448201527f74757265730000000000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152fd5b90506020813d6020116120e6575b816120d760209383610837565b8101031261077957515f6119b3565b3d91506120ca565b6040513d5f823e3d90fd5b608460405162461bcd60e51b815260206004820152602a60248201527f696e76616c696420626c6f636b3a2076616c696461746f72207365742068617360448201527f68206d69736d61746368000000000000000000000000000000000000000000006064820152fd5b90506020813d60201161218d575b8161217e60209383610837565b8101031261077957515f6117ec565b3d9150612171565b606460405162461bcd60e51b815260206004820152601760248201527f496e76616c696420636861696e206964206c656e6774680000000000000000006044820152fd5b60329150115f6113e8565b7fdb020436000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b507fe42a1980000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b60a46040517ff492ef2b00000000000000000000000000000000000000000000000000000000815260206004820152604360248201527f74727573746564206e6578742076616c696461746f722073657420686173682060448201527f646f6573206e6f74206d6174636820686173682073746f726564206f6e20636860648201527f61696e00000000000000000000000000000000000000000000000000000000006084820152fd5b9150506020813d602011612315575b8161230460209383610837565b81010312610779578290515f611374565b3d91506122f7565b906060806123348451608085526080850190610ed8565b936020810151602085015267ffffffffffffffff6040820151166040850152015160070b91015290565b6020815260a0810182519060806020840152815180915260c0830190602060c08260051b8601019301915f905b8282106123d2575050505067ffffffffffffffff60606123c86080936020870151151560408701526040870151601f19878303018488015261231d565b9401511691015290565b9091929360208061240d837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff408a60019603018652885161231d565b96019201920190929161238b565b602061242b818351510151610fa6565b0167ffffffffffffffff81511690604083019167ffffffffffffffff83515116908181036125ca57505067ffffffffffffffff61249291511667ffffffffffffffff60408551510151169260405191612483836107aa565b825260208201938452516126a8565b60038110156106a057600281149081156125bf575b506125885750806020806124ea930151604051809481927fc87f1f690000000000000000000000000000000000000000000000000000000083526004830161235e565b038173__$7440a880b7578767f72184d998805816e4$__5af49182156120ee575f92612552575b5051516101400151808203612524575050565b7f4d6d8014000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b9091506020813d602011612580575b8161256e60209383610837565b81010312610779575190610140612511565b3d9150612561565b67ffffffffffffffff9051167f90f4dbed000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b60019150145f6124a7565b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b60405190612605826107aa565b5f8252604051602083612617836107aa565b5f83525f828401520152565b60405190612630826107e2565b5f606083818152826020820152604051612649816107e2565b828152836020820152836040820152838382015260408201520152565b67ffffffffffffffff60019116019067ffffffffffffffff8211610f7957565b9067ffffffffffffffff8091169116019067ffffffffffffffff8211610f7957565b67ffffffffffffffff81511667ffffffffffffffff835116908181105f146126d257505050505f90565b11156126df575050600290565b602067ffffffffffffffff81819301511692015116908181105f146127045750505f90565b111561270f57600290565b600190565b5f929183915b81831061272a5750505060019190565b90919360ff6127607fff000000000000000000000000000000000000000000000000000000000000006020888601015116612a15565b16906009821161279257600a810290808204600a1490151715610f795760019161278991610f6c565b9401919061271a565b5050505090505f905f90565b8051821015610f0e5760209160051b010190565b5160048110156106a05790565b9067ffffffffffffffff8091169116029067ffffffffffffffff8216918203610f7957565b6020015160600151925f9291835b8151805186101561282a5760019167ffffffffffffffff6040612818896128229561279e565b5101511690612686565b9401936127f2565b509092919493505f5f935b85518510156129705761284b6119ef868861279e565b60048110156106a057600114612966576020612867868861279e565b510151519560208701905f5b8351805182101561295a57816128889161279e565b51516040516128b66020828180820195805191829101875e81015f838201520301601f198101835282610837565b519020896040516128e26020828181019451808a875e81015f838201520301601f198101835282610837565b519020146128f257600101612873565b83985067ffffffffffffffff9197949250612818604091612913955161279e565b905b61292660ff602086015116836127bf565b67ffffffffffffffff8061293e60ff885116876127bf565b16911611612952576001905b019394612835565b505050505050565b50509591965050612915565b949360019061294a565b5067ffffffffffffffff9350839294509060ff6129976129a09382602089015116906127bf565b955116906127bf565b16911611156129ab57565b608460405162461bcd60e51b815260206004820152602160248201527f696e73756666696369656e7420766f74696e6720706f776572206f7665726c6160448201527f70000000000000000000000000000000000000000000000000000000000000006064820152fd5b60f81c602f811180612ad7575b15612a4f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd00160ff1690565b6060811180612acd575b15612a86577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa90160ff1690565b6040811180612ac3575b15612abd577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc90160ff1690565b5060ff90565b5060478110612a90565b5060678110612a59565b50603a8110612a2256fea164736f6c634300081c000a",
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
