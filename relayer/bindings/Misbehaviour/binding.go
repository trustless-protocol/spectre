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
	ClockDrift      uint32
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
	ABI: "[{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"},{\"name\":\"clockDrift\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"misbehaviour_\",\"type\":\"tuple\",\"internalType\":\"structIMisbehaviourMsgs.Misbehaviour\",\"components\":[{\"name\":\"client_id\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ChainId\",\"components\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"header1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"header2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}]},{\"name\":\"trustedConsensusState1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"trustedConsensusState2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIMisbehaviourMsgs.MisbehaviourOutput\",\"components\":[{\"name\":\"trustedHeight1\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight2\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckForMisbehaviourFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidChainId\",\"inputs\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidClientId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MisbehaviourNotDetected\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MisbehaviourVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ValSetHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x608080604052346015576130d8908161001a8239f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c63265942ff14610024575f80fd5b3461015657610120366003190112610156576004356001600160401b0381116101565761014060031982360301126101565761005e610200565b9080600401356001600160401b038111610156576100fa91610089610124926004369184010161025b565b845261009836602483016102ad565b60208501526100aa36606483016102f7565b60408501526100bb60a48201610328565b60608501526100cc60c48201610328565b60808501526100dd60e48201610339565b60a08501526100ef6101048201610346565b60c085015201610328565b60e08201526024356001600160401b038111610156576101529161012561014692369060040161090e565b61012e366109c8565b61013736610a09565b91610140610353565b93610b02565b60405191829182610a4a565b0390f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b0382111761018957604052565b61015a565b608081019081106001600160401b0382111761018957604052565b606081019081106001600160401b0382111761018957604052565b60a081019081106001600160401b0382111761018957604052565b90601f801991011681019081106001600160401b0382111761018957604052565b60405190610210610100836101df565b565b60405190610210610260836101df565b604051906102106060836101df565b604051906102106040836101df565b6001600160401b03811161018957601f01601f191660200190565b81601f820112156101565760208135910161027582610240565b9261028360405194856101df565b8284528282011161015657815f92602092838601378301015290565b359060ff8216820361015657565b9190826040910312610156576040516102c58161016e565b60206102de8183956102d68161029f565b85520161029f565b910152565b35906001600160401b038216820361015657565b91908260409103126101565760405161030f8161016e565b60206102de818395610320816102e3565b8552016102e3565b359063ffffffff8216820361015657565b3590811515820361015657565b3590600282101561015657565b61010435906001600160801b038216820361015657565b35906001600160801b038216820361015657565b809291039160608312610156576040516103978161016e565b6040819483358352601f1901126101565760209060408051936103b98561016e565b6103c4848201610328565b85520135828401520152565b6001600160401b0381116101895760051b60200190565b919060c083820312610156576040516103ff8161018e565b809361040a816102e3565b825261041860208201610328565b602083015261042a836040830161037e565b604083015260a0810135906001600160401b03821161015657019180601f840112156101565782359261045c846103d0565b9361046a60405195866101df565b80855260208086019160051b830101918383116101565760208101915b83831061049957505050505060600152565b82356001600160401b038111610156578201906040828703601f19011261015657604051916104c78361016e565b6020810135600481101561015657835260408101356001600160401b0381116101565760209101019060808288031261015657604051926105078461018e565b82356001600160401b038111610156578861052391850161025b565b84526105316020840161036a565b602085015261054260408401610339565b60408501526060830135936001600160401b0385116101565761056a8960209687960161025b565b606082015283820152815201920191610487565b919060408382031261015657604051906105978261016e565b819380356001600160401b0381116101565781016102c081840312610156576105be610212565b906105c984826102f7565b825260408101356001600160401b03811161015657846105ea91830161025b565b60208301526105fb606082016102e3565b604083015261060c6080820161036a565b606083015261061d60a08201610339565b608083015261062f8460c0830161037e565b60a08301526106416101208201610339565b60c083015261014081013560e083015261065e6101608201610339565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526106ad6102208201610339565b6101c08301526102408101356101e08301526106cc6102608201610339565b6102008301526102808101356102208301526102a0810135906001600160401b038211610156576106ff9185910161025b565b61024082015283526020810135916001600160401b038311610156576020926102de92016103e7565b919060808382031261015657604051906107418261018e565b819380356001600160401b0381116101565760609261076191830161025b565b835260208101356020840152610779604082016102e3565b60408401520135908160070b82036101565760600152565b91909160808184031261015657604051906107ab8261018e565b819381356001600160401b03811161015657820181601f820112156101565780356107d5816103d0565b916107e360405193846101df565b81835260208084019260051b820101918483116101565760208201905b8382106108515750505050835261081960208301610339565b60208401526040820135906001600160401b0382116101565782610846606094926102de94869401610728565b6040860152016102e3565b81356001600160401b0381116101565760209161087388848094880101610728565b815201910190610800565b919060a08382031261015657604051906108978261018e565b819380356001600160401b03811161015657826108b591830161057e565b835260208101356001600160401b03811161015657826108d6918301610791565b60208401526108e882604083016102f7565b60408401526080810135916001600160401b038311610156576060926102de9201610791565b91906060838203126101565760405190610927826101a9565b819380356001600160401b03811161015657810160408184031261015657604051906109528261016e565b80356001600160401b03811161015657816109748660209361097c950161025b565b8452016102e3565b6020820152835260208101356001600160401b03811161015657826109a291830161087e565b60208401526040810135916001600160401b038311610156576040926102de920161087e565b606090604319011261015657604051906109e1826101a9565b816044356001600160801b038116810361015657815260643560208201526040608435910152565b60609060a31901126101565760405190610a22826101a9565b8160a4356001600160801b038116810361015657815260c4356020820152604060e435910152565b61021090929192604060206080830195610a7a8482516001600160401b0360208092828151168552015116910152565b01519101906001600160401b0360208092828151168552015116910152565b60405190610aa68261016e565b5f6020838281520152565b60405190610abe8261016e565b81610ac7610a99565b815260206102de610a99565b15610ada57565b7fa179f8c9000000000000000000000000000000000000000000000000000000005f5260045ffd5b93610b0b610ab1565b5084518051906020012094602083019586515151602001518051906020012014610b3490610ad3565b610b3d83610ce7565b6020810151926060820151610b559063ffffffff1690565b60e083015163ffffffff1690610b69610222565b95865263ffffffff16602086015263ffffffff1660408501528151604090920151516001600160401b031694610b9d610231565b9283526020830195610bb89087906001600160401b03169052565b8383895192878151610bd0906001600160801b031690565b916040015192610bdf95610f41565b604001948551938151610bf8906001600160801b031690565b916040015192610c0795610f41565b80516001600160401b0316925160400151602001516001600160401b0316610c2d610231565b6001600160401b0390941684526001600160401b03166020840152516001600160401b0316905160400151602001516001600160401b0316610c6d610231565b6001600160401b0390921682526001600160401b03166020820152610c90610231565b918252602082015290565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b9091610cd6610ce493604084526040840190610c9b565b916020818403910152610c9b565b90565b60406020820191610cf88351611139565b0190610d048251611139565b6020815151510180516020815191012090602084515151019182516020815191012003610ecb5750506020610d3e81835151510151611342565b016001600160401b03610e10610e04610df6610dbb610d6486516001600160401b031690565b95610dae6020610d8260408b515151016001600160401b0390511690565b98610d9d610d8e610231565b6001600160401b039092168252565b019788906001600160401b03169052565b516001600160401b031690565b94610dae6020610dd960408b515151016001600160401b0390511690565b97610de5610d8e610231565b019687906001600160401b03169052565b93516001600160401b031690565b6001600160401b031690565b911603610e5c576020604081819351510151015151925151015101515114610e3457565b7f93c0a5f9000000000000000000000000000000000000000000000000000000005f5260045ffd5b90610e8e6040610e7c81610ec895515151016001600160401b0390511690565b9251515101516001600160401b031690565b7f494c87bc000000000000000000000000000000000000000000000000000000005f526001600160401b0391821660045216602452604490565b5ffd5b51905190610f046040519283927ff6b6676b00000000000000000000000000000000000000000000000000000000845260048401610cbf565b0390fd5b634e487b7160e01b5f52601160045260245ffd5b906001600160801b03809116911603906001600160801b038211610f3c57565b610f08565b90919395949580610f5560608401516120c6565b0361108e57633b9aca00610f73610f6c828a61159d565b918661159d565b906001600160801b0382166001600160801b038216106110595790610f9791610f1c565b936020860194610fab865163ffffffff1690565b9063ffffffff82166001600160801b038216101561101c575050611007610ff361101794610210999a9461101094610fe389516115cf565b610fec836118a2565b9851611979565b95610ffd866119fe565b5163ffffffff1690565b63ffffffff1690565b8484611cb3565b611eae565b7fdb020436000000000000000000000000000000000000000000000000000000005f526001600160801b031660045263ffffffff1660245260445ffd5b7fe42a1980000000000000000000000000000000000000000000000000000000005f526001600160801b03861660045260245ffd5b6040517ff492ef2b00000000000000000000000000000000000000000000000000000000815260206004820152604360248201527f74727573746564206e6578742076616c696461746f722073657420686173682060448201527f646f6573206e6f74206d6174636820686173682073746f726564206f6e20636860648201527f61696e0000000000000000000000000000000000000000000000000000000000608482015260a490fd5b6020611149818351510151611342565b0180516001600160401b031690604083019161116d8351516001600160401b031690565b906001600160401b0382166001600160401b03821603611267575050516111df906001600160401b0316835151604001516001600160401b0316926111c26111b3610231565b6001600160401b039093168352565b6111d9602083019485906001600160401b03169052565b51611f4a565b61123257506101406111f460208301516120c6565b9151510151808203611204575050565b7f4d6d8014000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b517f90f4dbed000000000000000000000000000000000000000000000000000000005f526001600160401b031660045260245ffd5b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f526001600160401b039081166004521660245260445ffd5b604051906112ae8261016e565b5f602083606081520152565b8015610f3c575f190190565b5f19810191908211610f3c57565b91908203918211610f3c57565b634e487b7160e01b5f52603260045260245ffd5b908151811015611306570160200190565b6112e1565b9060018201809211610f3c57565b6001019081600111610f3c57565b6023019081602311610f3c57565b91908201809211610f3c57565b61134a6112a1565b5080518015908115611591575b50611569575f19908051805b611501575b505f1982146114e357600360fc1b7fff000000000000000000000000000000000000000000000000000000000000006113d26113ac6113a68661130b565b856112f5565b517fff000000000000000000000000000000000000000000000000000000000000001690565b1614806114ed575b6114e3575f906113e98361130b565b915b815183101561147b5761140a6114046113ac85856112f5565b60f81c90565b60ff811660308110908115611470575b50611463576001600160401b0360ff8192602f19011681600a85021601169116811061144b576001909201916113eb565b50915050611457610231565b9081525f602082015290565b5050915050611457610231565b60399150115f61141a565b91509160018111908115916114d7575b506114af57610ce49061149c610231565b9283526001600160401b03166020830152565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f61148b565b9050611457610231565b5060026114fb8383516112d4565b116113da565b602d60f81b61154361151e6113ac611518856112c6565b866112f5565b7fff000000000000000000000000000000000000000000000000000000000000001690565b1461155757611551906112ba565b80611363565b6115629192506112c6565b905f611368565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f611357565b906001600160801b03169081156115bb576001600160801b03160490565b634e487b7160e01b5f52601260045260245ffd5b8051801590811561180e575b506117c9575f5b81518110156117c5576115fb61151e6113ac83856112f5565b7f6100000000000000000000000000000000000000000000000000000000000000811015908161179a575b811561173e575b81156116fe575b81156116f0575b81156116c6575b811561169c575b5015611657576001016115e2565b60405162461bcd60e51b815260206004820152601860248201527f696e76616c696420636861696e206964206368617273657400000000000000006044820152606490fd5b7f2e000000000000000000000000000000000000000000000000000000000000009150145f611649565b7f5f0000000000000000000000000000000000000000000000000000000000000081149150611642565b602d60f81b8114915061163b565b9050600360fc1b81101580611714575b90611634565b507f390000000000000000000000000000000000000000000000000000000000000081111561170e565b90507f410000000000000000000000000000000000000000000000000000000000000081101580611770575b9061162d565b507f5a0000000000000000000000000000000000000000000000000000000000000081111561176a565b7f7a000000000000000000000000000000000000000000000000000000000000008111159150611626565b5050565b60405162461bcd60e51b815260206004820152601760248201527f496e76616c696420636861696e206964206c656e6774680000000000000000006044820152606490fd5b60329150115f6115db565b604051906118268261016e565b815f815260206102de610a99565b604051906118418261018e565b606080835f81525f6020820152611856611819565b60408201520152565b6040519061186c8261018e565b5f6060838181528260208201526040516118858161018e565b828152836020820152836040820152838382015260408201520152565b6040516118ae8161016e565b6040516118ba8161016e565b6118c2610212565b6118ca610a99565b8152606060208201525f60408201525f60608201525f60808201526118ed611819565b60a08201525f60c08201525f60e08201525f6101008201525f6101208201525f6101408201525f6101608201525f6101808201525f6101a08201525f6101c08201525f6101e08201525f6102008201525f61022082015260606102408201528152611956611834565b60208201528152602061196761185f565b91015260208151910151610c90610231565b926119f1905f608060405161198d816101c4565b606081528260208201528260408201526119a561185f565b606082015201526001600160801b0360606001600160401b036020604085015101511692015193604051966119d9886101c4565b87521660208601526001600160401b03166040850152565b6060830152608082015290565b611a17815151611a11602082015161290c565b90612918565b604060208351015101515103611a3857806020610210925191015190612180565b608460405162461bcd60e51b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152fd5b6001600160801b03633b9aca00911602906001600160801b038216918203610f3c57565b15611acd57565b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b805191908290602001825e015f815290565b15611b5057565b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b6001600160401b036001911601906001600160401b038211610f3c57565b906001600160401b03809116911601906001600160401b038211610f3c57565b15611bd957565b608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b15611c4957565b608460405162461bcd60e51b815260206004820152602f60248201527f696e76616c696420626c6f636b3a206e6578742076616c696461746f7220736560448201527f742068617368206d69736d6174636800000000000000000000000000000000006064820152fd5b9290916001600160401b03611cc89116611aa2565b906020830191611cdf83516001600160801b031690565b6001600160801b0381166001600160801b03841610928315611e8b575b505050611e2057611d4d906001600160801b03611d45611d39611d2b606088515101516001600160801b031690565b93516001600160801b031690565b6001600160801b031690565b911611611ac6565b611d9d60208351510151604051611d7a81611d6c602082018095611b37565b03601f1981018352826101df565b5190208251604051611d9481611d6c602082018095611b37565b51902014611b49565b611db9611db460408301516001600160401b031690565b611b94565b916001600160401b036040825151019381611ddc866001600160401b0390511690565b911691168103611dff575061021092506101406080915151015191015114611c42565b915050611e1a610e04610210936001600160401b0390511690565b11611bd2565b60405162461bcd60e51b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152608490fd5b6001600160801b039293508291611ea191610f1c565b92169116115f8080611cfc565b9091611ec7611db460408501516001600160401b031690565b6001600160401b0380611ee76040865151016001600160401b0390511690565b9216911614611f285761021092611f07916060845192015190519161252b565b611f0f610231565b600281529060036020830152602081519101519061246a565b506102109150611f0f610231565b634e487b7160e01b5f52602160045260245ffd5b90611f5491612654565b6003811015611f745760028114908115611f6c575090565b600191501490565b611f36565b6040516101e09190611f8b83826101df565b600e815291601f1901366020840137565b90611fa6826103d0565b611fb360405191826101df565b8281528092611fc4601f19916103d0565b0190602036910137565b8051156113065760200190565b8051600110156113065760400190565b8051600210156113065760600190565b8051600310156113065760800190565b8051600410156113065760a00190565b8051600510156113065760c00190565b8051600610156113065760e00190565b805160071015611306576101000190565b805160081015611306576101200190565b805160091015611306576101400190565b8051600a1015611306576101600190565b8051600b1015611306576101800190565b8051600c1015611306576101a00190565b8051600d1015611306576101c00190565b80518210156113065760209160051b010190565b805151908115612163576120d982611f9c565b915f5b82518051821015612155579061214461213f60206120fc846001966120b2565b51015161213a6001600160401b036040612117878b516120b2565b51015116604051926121288461016e565b83526001600160401b03166020830152565b612751565b61281d565b61214e82876120b2565b52016120dc565b505090505f610ce492612872565b50505f90565b60041115611f7457565b516004811015611f745790565b60208101906001600160401b036121b7610e0460406121a78651516001600160401b031690565b945101516001600160401b031690565b9116036123ad57516060015180518251510361231c575f925f5b82518110156122a5576121e481846120b2565b5194600186516121f381612169565b6121fc81612169565b1461229b575060019460206122128387516120b2565b515160405161222881611d6c8582018095611b37565b5190209101515160405161224481611d6c602082018095611b37565b51902003612256576001905b016121d1565b60405162461bcd60e51b815260206004820152601d60248201527f696e76616c696420636f6d6d69743a206661756c7479207369676e65720000006044820152606490fd5b9450600190612250565b5092915050156122b157565b60405162461bcd60e51b815260206004820152602560248201527f696e76616c696420636f6d6d69743a206e6f2070726573656e74207369676e6160448201527f74757265730000000000000000000000000000000000000000000000000000006064820152608490fd5b60405162461bcd60e51b815260206004820152604860248201527f696e76616c696420636f6d6d69743a206e756d626572206f66207369676e617460448201527f7572657320646f6573206e6f74206d61746368206e756d626572206f6620766160648201527f6c696461746f7273000000000000000000000000000000000000000000000000608482015260a490fd5b60405162461bcd60e51b815260206004820152601f60248201527f696e76616c696420636f6d6d69743a20686569676874206d69736d61746368006044820152606490fd5b1561231c57565b1561240057565b608460405162461bcd60e51b815260206004820152602160248201527f696e73756666696369656e7420766f74696e6720706f776572206f7665726c6160448201527f70000000000000000000000000000000000000000000000000000000000000006064820152fd5b602060609193929301510151915161248583518251146123f2565b61248e81612ba3565b905f90815b85518310156125135760016124b16124ab85896120b2565b51612173565b6124ba81612169565b1461250a576124e7906124e160406124d286866120b2565b5101516001600160401b031690565b90611bb2565b916124f3858585612bf5565b612502576001905b0191612493565b505050505050565b916001906124fb565b9050610210945061252693929150612bf5565b6123f9565b60200151606001519051909161254082612ba3565b5f61254b8451611f9c565b905f945b86518610156126425761256286886120b2565b516001815161257081612169565b61257981612169565b1461263757602001515160208151910120965f5b825181101561262d576125b06125ac6125a683886120b2565b51151590565b1590565b80612612575b6125c25760010161258d565b928298506125e16125e8916124e160406124d2886125ee979d986120b2565b93856120b2565b60019052565b6125f9858584612bf5565b612609576001905b01949561254f565b50505050505050565b508861261e82856120b2565b515160208151910120146125b6565b50959096506125ee565b509594600190612601565b50905061021094506125269350612bf5565b80516001600160401b03166001600160401b0361267b610e0485516001600160401b031690565b91168181101561268d57505050505f90565b111561269a575050600290565b6126cc610e0460206126bd816001600160401b039501516001600160401b031690565b9401516001600160401b031690565b9116818110156126dc5750505f90565b11156126e757600290565b600190565b604051608091906126fd83826101df565b6041815291601f1901366020840137565b6040519061271d6020836101df565b5f808352366020840137565b9061273382610240565b61274060405191826101df565b8281528092611fc4601f1991610240565b6020810180516001600160401b03166024816127ec575b6127729150612729565b91600a6020840153602260218401536127a290612799612793600286612c4f565b85612c66565b90519084612cb6565b906127b7610e0482516001600160401b031690565b6127c057505090565b6127e1610e046127d36127e89486612c7c565b92516001600160401b031690565b9083612d06565b5090565b506127f96127fe91612c24565b611319565b60240180602411610f3c5761277290612768565b6040513d5f823e3d90fd5b805160018101809111610f3c5761283390612729565b8051156113065761285d816128506020945f868196015382612d3f565b5060405191828092611b37565b039060025afa1561286d575f5190565b612812565b9092919280840393808511610f3c57600185146128fb5760015b8060011b908682101561289f575061288c565b91929394955050820191828111610f3c57826128bb9185612872565b916128c69293612872565b6128ce6126ec565b9182511561130657825f9261285d9260016020809701536021830152604182015260405191828092611b37565b50906129089293506120b2565b5190565b61213f610ce491612d94565b90612ad261213f610240610ce49461292e611f79565b9461293c61213f8351612dee565b61294587611fce565b5261294f86611fdb565b5261297261213f61296d610e0460408501516001600160401b031690565b612ece565b61297b86611feb565b5261299b61213f61299660608401516001600160801b031690565b612f02565b6129a486611ffb565b52608081015115612b7d576129bf61213f60a0830151612fcd565b6129c88661200b565b5260c081015115612b57576129e361213f60e0830151613093565b6129ec8661201b565b5261010081015115612b3157612a0961213f610120830151613093565b612a128661202b565b52612a2461213f610140830151613093565b612a2d8661203b565b52612a3f61213f610160830151613093565b612a488661204c565b52612a5a61213f610180830151613093565b612a638661205d565b52612a7561213f6101a0830151613093565b612a7e8661206e565b526101c081015115612b0b57612a9b61213f6101e0830151613093565b612aa48661207f565b5261020081015115612ae557612ac161213f610220830151613093565b612aca86612090565b520151612d94565b612adb826120a1565b525f815191612872565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d612ac1565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d612a9b565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d612a09565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d6129e3565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d6129bf565b5f9190825b8151841015612bdd57612bd56001916001600160401b036040612bcb88876120b2565b5101511690611bb2565b930192612ba8565b925050565b81810292918115918404141715610f3c57565b906001600160401b0360ff612c16612c209483836020890151169116612be2565b9451169116612be2565b1090565b906001915b6080811015612c355750565b60019060071c920191612c29565b6020600a910153600190565b602082600a9201015360018101809111610f3c5790565b602082819201015360018101809111610f3c5790565b60208260109201015360018101809111610f3c5790565b60206008910153600190565b60208260129201015360018101809111610f3c5790565b8160209193929301015260208101809111610f3c5790565b91906021600193015b6080821015612ceb57906001929391530190565b600180916080607f85161781530193019060071c9092612cd7565b9092919083016020015b6080821015612d2457906001929391530190565b600180916080607f85161781530193019060071c9092612d10565b908051918215612d61576021602084930191015e60010180600111610f3c5790565b505050600190565b908092918251928315612d8d57839260208092019201015e8101809111610f3c5790565b5050505090565b805115612de557612da58151612c24565b806001019081600111610f3c5760019083510101809111610f3c57612dcc6127e891612729565b91600a6020840153612ddf815184612cce565b83612d69565b50610ce461270e565b5f90612e0181516001600160401b031690565b6001600160401b038116612ea9575b50612e396020820192612e2d610e0485516001600160401b031690565b80612e8d575b50612729565b915f91612e50610e0482516001600160401b031690565b612e6a575b506127b7610e0482516001600160401b031690565b612e86919250612e7f610e046127d386612c93565b9084612d06565b905f612e55565b90612e9d6127f9612ea393612c24565b90611335565b5f612e33565b612ec7919250612ec26127f9916001600160401b031690565b612c24565b905f612e10565b8015612de557612edd81612c24565b60010180600111610f3c57612ef46127e891612729565b916008602084015382612cce565b6001600160801b0316633b9aca008104906001600160801b035f9216918215159182612fa9575b633b9aca006001600160801b03910616908115159081612f8e575b612f4d90612729565b935f93612f72575b50612f5f57505090565b612f6c6127e89284612c7c565b83612d06565b612f8791935060086020860153600185612d06565b915f612f55565b612f9a6127f984612c24565b810180911115612f4457610f08565b90506001600160801b03633b9aca00612fc46127f986612c24565b92915050612f29565b6020810151612fe263ffffffff825116612c24565b90816001019182600111610f3c57602301809211610f3c576130436130096127e893612729565b9160086020840153602061303961279361303361302d611007865163ffffffff1690565b87612cce565b86612c9f565b9101519083612cb6565b50612ddf815161308d61303361307161306c8461306761306282612c24565b611327565b611335565b612729565b9661308461307e89612c43565b89612c66565b90519088612cb6565b85612d06565b604051906060906130a482846101df565b602283526127e891600a906020850190601f19013682375360206021840153600283612cb656fea164736f6c634300081c000a",
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

// Misbehaviour is a free data retrieval call binding the contract method 0x265942ff.
//
// Solidity: function misbehaviour((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32) clientState, ((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64))) misbehaviour_, (uint128,bytes32,bytes32) trustedConsensusState1, (uint128,bytes32,bytes32) trustedConsensusState2, uint128 time) pure returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourCaller) Misbehaviour(opts *bind.CallOpts, clientState IICS07TendermintMsgsClientState, misbehaviour_ IMisbehaviourMsgsMisbehaviour, trustedConsensusState1 IICS07TendermintMsgsConsensusState, trustedConsensusState2 IICS07TendermintMsgsConsensusState, time *big.Int) (IMisbehaviourMsgsMisbehaviourOutput, error) {
	var out []interface{}
	err := _ContractMisbehaviour.contract.Call(opts, &out, "misbehaviour", clientState, misbehaviour_, trustedConsensusState1, trustedConsensusState2, time)

	if err != nil {
		return *new(IMisbehaviourMsgsMisbehaviourOutput), err
	}

	out0 := *abi.ConvertType(out[0], new(IMisbehaviourMsgsMisbehaviourOutput)).(*IMisbehaviourMsgsMisbehaviourOutput)

	return out0, err

}

// Misbehaviour is a free data retrieval call binding the contract method 0x265942ff.
//
// Solidity: function misbehaviour((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32) clientState, ((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64))) misbehaviour_, (uint128,bytes32,bytes32) trustedConsensusState1, (uint128,bytes32,bytes32) trustedConsensusState2, uint128 time) pure returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourSession) Misbehaviour(clientState IICS07TendermintMsgsClientState, misbehaviour_ IMisbehaviourMsgsMisbehaviour, trustedConsensusState1 IICS07TendermintMsgsConsensusState, trustedConsensusState2 IICS07TendermintMsgsConsensusState, time *big.Int) (IMisbehaviourMsgsMisbehaviourOutput, error) {
	return _ContractMisbehaviour.Contract.Misbehaviour(&_ContractMisbehaviour.CallOpts, clientState, misbehaviour_, trustedConsensusState1, trustedConsensusState2, time)
}

// Misbehaviour is a free data retrieval call binding the contract method 0x265942ff.
//
// Solidity: function misbehaviour((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32) clientState, ((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64))) misbehaviour_, (uint128,bytes32,bytes32) trustedConsensusState1, (uint128,bytes32,bytes32) trustedConsensusState2, uint128 time) pure returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourCallerSession) Misbehaviour(clientState IICS07TendermintMsgsClientState, misbehaviour_ IMisbehaviourMsgsMisbehaviour, trustedConsensusState1 IICS07TendermintMsgsConsensusState, trustedConsensusState2 IICS07TendermintMsgsConsensusState, time *big.Int) (IMisbehaviourMsgsMisbehaviourOutput, error) {
	return _ContractMisbehaviour.Contract.Misbehaviour(&_ContractMisbehaviour.CallOpts, clientState, misbehaviour_, trustedConsensusState1, trustedConsensusState2, time)
}
