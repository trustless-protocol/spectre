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
	ABI: "[{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"misbehaviour_\",\"type\":\"tuple\",\"internalType\":\"structIMisbehaviourMsgs.Misbehaviour\",\"components\":[{\"name\":\"client_id\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ChainId\",\"components\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"header1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"header2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}]},{\"name\":\"trustedConsensusState1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"trustedConsensusState2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIMisbehaviourMsgs.MisbehaviourOutput\",\"components\":[{\"name\":\"trustedHeight1\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight2\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckForMisbehaviourFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidChainId\",\"inputs\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidClientId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MisbehaviourNotDetected\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MisbehaviourVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ValSetHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60808060405234601557612fec908161001a8239f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c63a6fe8f5614610024575f80fd5b3461014457610120366003190112610144576004356001600160401b0381116101445761012060031982360301126101445761005e6101ee565b9080600401356001600160401b038111610144576100e8916100896101049260043691840101610248565b8452610098366024830161029a565b60208501526100aa36606483016102e4565b60408501526100bb60a48201610315565b60608501526100cc60c48201610315565b60808501526100dd60e48201610326565b60a085015201610333565b60c08201526024356001600160401b03811161014457610140916101136101349236906004016108fb565b61011c366109b5565b610125366109f6565b9161012e610340565b93610aef565b60405191829182610a37565b0390f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b0382111761017757604052565b610148565b608081019081106001600160401b0382111761017757604052565b606081019081106001600160401b0382111761017757604052565b60a081019081106001600160401b0382111761017757604052565b90601f801991011681019081106001600160401b0382111761017757604052565b604051906101fd60e0836101cd565b565b604051906101fd610260836101cd565b604051906101fd6060836101cd565b604051906101fd6040836101cd565b6001600160401b03811161017757601f01601f191660200190565b81601f82011215610144576020813591016102628261022d565b9261027060405194856101cd565b8284528282011161014457815f92602092838601378301015290565b359060ff8216820361014457565b9190826040910312610144576040516102b28161015c565b60206102cb8183956102c38161028c565b85520161028c565b910152565b35906001600160401b038216820361014457565b9190826040910312610144576040516102fc8161015c565b60206102cb81839561030d816102d0565b8552016102d0565b359063ffffffff8216820361014457565b3590811515820361014457565b3590600282101561014457565b61010435906001600160801b038216820361014457565b35906001600160801b038216820361014457565b809291039160608312610144576040516103848161015c565b6040819483358352601f1901126101445760209060408051936103a68561015c565b6103b1848201610315565b85520135828401520152565b6001600160401b0381116101775760051b60200190565b919060c083820312610144576040516103ec8161017c565b80936103f7816102d0565b825261040560208201610315565b6020830152610417836040830161036b565b604083015260a0810135906001600160401b03821161014457019180601f8401121561014457823592610449846103bd565b9361045760405195866101cd565b80855260208086019160051b830101918383116101445760208101915b83831061048657505050505060600152565b82356001600160401b038111610144578201906040828703601f19011261014457604051916104b48361015c565b6020810135600481101561014457835260408101356001600160401b0381116101445760209101019060808288031261014457604051926104f48461017c565b82356001600160401b0381116101445788610510918501610248565b845261051e60208401610357565b602085015261052f60408401610326565b60408501526060830135936001600160401b0385116101445761055789602096879601610248565b606082015283820152815201920191610474565b919060408382031261014457604051906105848261015c565b819380356001600160401b0381116101445781016102c081840312610144576105ab6101ff565b906105b684826102e4565b825260408101356001600160401b03811161014457846105d7918301610248565b60208301526105e8606082016102d0565b60408301526105f960808201610357565b606083015261060a60a08201610326565b608083015261061c8460c0830161036b565b60a083015261062e6101208201610326565b60c083015261014081013560e083015261064b6101608201610326565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a083015261069a6102208201610326565b6101c08301526102408101356101e08301526106b96102608201610326565b6102008301526102808101356102208301526102a0810135906001600160401b038211610144576106ec91859101610248565b61024082015283526020810135916001600160401b038311610144576020926102cb92016103d4565b9190608083820312610144576040519061072e8261017c565b819380356001600160401b0381116101445760609261074e918301610248565b835260208101356020840152610766604082016102d0565b60408401520135908160070b82036101445760600152565b91909160808184031261014457604051906107988261017c565b819381356001600160401b03811161014457820181601f820112156101445780356107c2816103bd565b916107d060405193846101cd565b81835260208084019260051b820101918483116101445760208201905b83821061083e5750505050835261080660208301610326565b60208401526040820135906001600160401b0382116101445782610833606094926102cb94869401610715565b6040860152016102d0565b81356001600160401b0381116101445760209161086088848094880101610715565b8152019101906107ed565b919060a08382031261014457604051906108848261017c565b819380356001600160401b03811161014457826108a291830161056b565b835260208101356001600160401b03811161014457826108c391830161077e565b60208401526108d582604083016102e4565b60408401526080810135916001600160401b038311610144576060926102cb920161077e565b9190606083820312610144576040519061091482610197565b819380356001600160401b038111610144578101604081840312610144576040519061093f8261015c565b80356001600160401b0381116101445781610961866020936109699501610248565b8452016102d0565b6020820152835260208101356001600160401b038111610144578261098f91830161086b565b60208401526040810135916001600160401b038311610144576040926102cb920161086b565b606090604319011261014457604051906109ce82610197565b816044356001600160801b038116810361014457815260643560208201526040608435910152565b60609060a31901126101445760405190610a0f82610197565b8160a4356001600160801b038116810361014457815260c4356020820152604060e435910152565b6101fd90929192604060206080830195610a678482516001600160401b0360208092828151168552015116910152565b01519101906001600160401b0360208092828151168552015116910152565b60405190610a938261015c565b5f6020838281520152565b60405190610aab8261015c565b81610ab4610a86565b815260206102cb610a86565b15610ac757565b7fa179f8c9000000000000000000000000000000000000000000000000000000005f5260045ffd5b93610af8610a9e565b5084518051906020012094602083019586515151602001518051906020012014610b2190610ac0565b610b2a83610cc4565b6020810151926060820151610b429063ffffffff1690565b610b4a61020f565b94855263ffffffff166020850152600f60408501528151604090920151516001600160401b031694610b7a61021e565b9283526020830195610b959087906001600160401b03169052565b8383895192878151610bad906001600160801b031690565b916040015192610bbc95610eff565b604001948551938151610bd5906001600160801b031690565b916040015192610be495610eff565b80516001600160401b0316925160400151602001516001600160401b0316610c0a61021e565b6001600160401b0390941684526001600160401b03166020840152516001600160401b0316905160400151602001516001600160401b0316610c4a61021e565b6001600160401b0390921682526001600160401b03166020820152610c6d61021e565b918252602082015290565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b9091610cb3610cc193604084526040840190610c78565b916020818403910152610c78565b90565b60406020820191610cd583516110f7565b0190610ce182516110f7565b6020815151510180516020815191012090602084515151019182516020815191012003610e89575050610dd16020610d1e81845151510151611300565b01610d86610d3382516001600160401b031690565b91610d79610d4f604087515151016001600160401b0390511690565b610d69610d5a61021e565b6001600160401b039096168652565b6001600160401b03166020850152565b516001600160401b031690565b90610dcc610da2604087515151016001600160401b0390511690565b610dbc610dad61021e565b6001600160401b039095168552565b6001600160401b03166020840152565b61156f565b610e1a576020604081819351510151015151925151015101515114610df257565b7f93c0a5f9000000000000000000000000000000000000000000000000000000005f5260045ffd5b90610e4c6040610e3a81610e8695515151016001600160401b0390511690565b9251515101516001600160401b031690565b7f4447469a000000000000000000000000000000000000000000000000000000005f526001600160401b0391821660045216602452604490565b5ffd5b51905190610ec26040519283927ff6b6676b00000000000000000000000000000000000000000000000000000000845260048401610c9c565b0390fd5b634e487b7160e01b5f52601160045260245ffd5b906001600160801b03809116911603906001600160801b038211610efa57565b610ec6565b90919395949580610f1360608401516120a7565b0361104c57633b9aca00610f31610f2a828a61158b565b918661158b565b906001600160801b0382166001600160801b038216106110175790610f5591610eda565b936020860194610f69865163ffffffff1690565b9063ffffffff82166001600160801b0382161015610fda575050610fc5610fb1610fd5946101fd999a94610fce94610fa189516115bd565b610faa83611890565b9851611967565b95610fbb866119ec565b5163ffffffff1690565b63ffffffff1690565b8484611ca1565b611ea8565b7fdb020436000000000000000000000000000000000000000000000000000000005f526001600160801b031660045263ffffffff1660245260445ffd5b7fe42a1980000000000000000000000000000000000000000000000000000000005f526001600160801b03861660045260245ffd5b6040517ff492ef2b00000000000000000000000000000000000000000000000000000000815260206004820152604360248201527f74727573746564206e6578742076616c696461746f722073657420686173682060448201527f646f6573206e6f74206d6174636820686173682073746f726564206f6e20636860648201527f61696e0000000000000000000000000000000000000000000000000000000000608482015260a490fd5b6020611107818351510151611300565b0180516001600160401b031690604083019161112b8351516001600160401b031690565b906001600160401b0382166001600160401b038216036112255750505161119d906001600160401b0316835151604001516001600160401b03169261118061117161021e565b6001600160401b039093168352565b611197602083019485906001600160401b03169052565b51611f30565b6111f057506101406111b260208301516120a7565b91515101518082036111c2575050565b7f4d6d8014000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b517f90f4dbed000000000000000000000000000000000000000000000000000000005f526001600160401b031660045260245ffd5b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f526001600160401b039081166004521660245260445ffd5b6040519061126c8261015c565b5f602083606081520152565b8015610efa575f190190565b5f19810191908211610efa57565b91908203918211610efa57565b634e487b7160e01b5f52603260045260245ffd5b9081518110156112c4570160200190565b61129f565b9060018201809211610efa57565b6001019081600111610efa57565b6023019081602311610efa57565b91908201809211610efa57565b61130861125f565b508051801590811561154f575b50611527575f19908051805b6114bf575b505f1982146114a157600360fc1b7fff0000000000000000000000000000000000000000000000000000000000000061139061136a611364866112c9565b856112b3565b517fff000000000000000000000000000000000000000000000000000000000000001690565b1614806114ab575b6114a1575f906113a7836112c9565b915b8151831015611439576113c86113c261136a85856112b3565b60f81c90565b60ff81166030811090811561142e575b50611421576001600160401b0360ff8192602f19011681600a850216011691168110611409576001909201916113a9565b5091505061141561021e565b9081525f602082015290565b505091505061141561021e565b60399150115f6113d8565b9150916001811190811591611495575b5061146d57610cc19061145a61021e565b9283526001600160401b03166020830152565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f611449565b905061141561021e565b5060026114b9838351611292565b11611398565b602d60f81b6115016114dc61136a6114d685611284565b866112b3565b7fff000000000000000000000000000000000000000000000000000000000000001690565b146115155761150f90611278565b80611321565b611520919250611284565b905f611326565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f611315565b634e487b7160e01b5f52602160045260245ffd5b906115799161214a565b6003811015611586571590565b61155b565b906001600160801b03169081156115a9576001600160801b03160490565b634e487b7160e01b5f52601260045260245ffd5b805180159081156117fc575b506117b7575f5b81518110156117b3576115e96114dc61136a83856112b3565b7f61000000000000000000000000000000000000000000000000000000000000008110159081611788575b811561172c575b81156116ec575b81156116de575b81156116b4575b811561168a575b5015611645576001016115d0565b60405162461bcd60e51b815260206004820152601860248201527f696e76616c696420636861696e206964206368617273657400000000000000006044820152606490fd5b7f2e000000000000000000000000000000000000000000000000000000000000009150145f611637565b7f5f0000000000000000000000000000000000000000000000000000000000000081149150611630565b602d60f81b81149150611629565b9050600360fc1b81101580611702575b90611622565b507f39000000000000000000000000000000000000000000000000000000000000008111156116fc565b90507f41000000000000000000000000000000000000000000000000000000000000008110158061175e575b9061161b565b507f5a00000000000000000000000000000000000000000000000000000000000000811115611758565b7f7a000000000000000000000000000000000000000000000000000000000000008111159150611614565b5050565b60405162461bcd60e51b815260206004820152601760248201527f496e76616c696420636861696e206964206c656e6774680000000000000000006044820152606490fd5b60329150115f6115c9565b604051906118148261015c565b815f815260206102cb610a86565b6040519061182f8261017c565b606080835f81525f6020820152611844611807565b60408201520152565b6040519061185a8261017c565b5f6060838181528260208201526040516118738161017c565b828152836020820152836040820152838382015260408201520152565b60405161189c8161015c565b6040516118a88161015c565b6118b06101ff565b6118b8610a86565b8152606060208201525f60408201525f60608201525f60808201526118db611807565b60a08201525f60c08201525f60e08201525f6101008201525f6101208201525f6101408201525f6101608201525f6101808201525f6101a08201525f6101c08201525f6101e08201525f6102008201525f61022082015260606102408201528152611944611822565b60208201528152602061195561184d565b91015260208151910151610c6d61021e565b926119df905f608060405161197b816101b2565b6060815282602082015282604082015261199361184d565b606082015201526001600160801b0360606001600160401b036020604085015101511692015193604051966119c7886101b2565b87521660208601526001600160401b03166040850152565b6060830152608082015290565b611a058151516119ff6020820151612820565b9061282c565b604060208351015101515103611a26578060206101fd9251910151906121f9565b608460405162461bcd60e51b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152fd5b6001600160801b03633b9aca00911602906001600160801b038216918203610efa57565b15611abb57565b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b805191908290602001825e015f815290565b15611b3e57565b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b6001600160401b036001911601906001600160401b038211610efa57565b906001600160401b03809116911601906001600160401b038211610efa57565b15611bc757565b608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b15611c3757565b608460405162461bcd60e51b815260206004820152602f60248201527f696e76616c696420626c6f636b3a206e6578742076616c696461746f7220736560448201527f742068617368206d69736d6174636800000000000000000000000000000000006064820152fd5b9290916001600160401b03611cb69116611a90565b906020830191611ccd83516001600160801b031690565b6001600160801b0381166001600160801b03841610928315611e85575b505050611e1a57611d3b906001600160801b03611d33611d27611d19606088515101516001600160801b031690565b93516001600160801b031690565b6001600160801b031690565b911611611ab4565b611d8b60208351510151604051611d6881611d5a602082018095611b25565b03601f1981018352826101cd565b5190208251604051611d8281611d5a602082018095611b25565b51902014611b37565b611da7611da260408301516001600160401b031690565b611b82565b916001600160401b036040825151019381611dca866001600160401b0390511690565b911691168103611ded57506101fd92506101406080915151015191015114611c30565b915050611e14611e086101fd936001600160401b0390511690565b6001600160401b031690565b11611bc0565b60405162461bcd60e51b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152608490fd5b6001600160801b039293508291611e9b91610eda565b92169116115f8080611cea565b9091611ec1611da260408501516001600160401b031690565b6001600160401b0380611ee16040865151016001600160401b0390511690565b9216911614611f22576101fd92611f019160608451920151905191612524565b611f0961021e565b6002815290600360208301526020815191015190612463565b506101fd9150611f0961021e565b90611f3a9161214a565b60038110156115865760028114908115611f52575090565b600191501490565b6040516101e09190611f6c83826101cd565b600e815291601f1901366020840137565b90611f87826103bd565b611f9460405191826101cd565b8281528092611fa5601f19916103bd565b0190602036910137565b8051156112c45760200190565b8051600110156112c45760400190565b8051600210156112c45760600190565b8051600310156112c45760800190565b8051600410156112c45760a00190565b8051600510156112c45760c00190565b8051600610156112c45760e00190565b8051600710156112c4576101000190565b8051600810156112c4576101200190565b8051600910156112c4576101400190565b8051600a10156112c4576101600190565b8051600b10156112c4576101800190565b8051600c10156112c4576101a00190565b8051600d10156112c4576101c00190565b80518210156112c45760209160051b010190565b805151908115612144576120ba82611f7d565b915f5b82518051821015612136579061212561212060206120dd84600196612093565b51015161211b6001600160401b0360406120f8878b51612093565b51015116604051926121098461015c565b83526001600160401b03166020830152565b612665565b612731565b61212f8287612093565b52016120bd565b505090505f610cc192612786565b50505f90565b80516001600160401b03166001600160401b03612171611e0885516001600160401b031690565b91168181101561218357505050505f90565b1115612190575050600290565b6121c2611e0860206121b3816001600160401b039501516001600160401b031690565b9401516001600160401b031690565b9116818110156121d25750505f90565b11156121dd57600290565b600190565b6004111561158657565b5160048110156115865790565b602001516060015180518251510361235a575f925f5b82518110156122e3576122228184612093565b519460018651612231816121e2565b61223a816121e2565b146122d957506001946020612250838751612093565b515160405161226681611d5a8582018095611b25565b5190209101515160405161228281611d5a602082018095611b25565b51902003612294576001905b0161220f565b60405162461bcd60e51b815260206004820152601d60248201527f696e76616c696420636f6d6d69743a206661756c7479207369676e65720000006044820152606490fd5b945060019061228e565b5092915050156122ef57565b60405162461bcd60e51b815260206004820152602560248201527f696e76616c696420636f6d6d69743a206e6f2070726573656e74207369676e6160448201527f74757265730000000000000000000000000000000000000000000000000000006064820152608490fd5b60405162461bcd60e51b815260206004820152604860248201527f696e76616c696420636f6d6d69743a206e756d626572206f66207369676e617460448201527f7572657320646f6573206e6f74206d61746368206e756d626572206f6620766160648201527f6c696461746f7273000000000000000000000000000000000000000000000000608482015260a490fd5b1561235a57565b156123f957565b608460405162461bcd60e51b815260206004820152602160248201527f696e73756666696369656e7420766f74696e6720706f776572206f7665726c6160448201527f70000000000000000000000000000000000000000000000000000000000000006064820152fd5b602060609193929301510151915161247e83518251146123eb565b61248781612ab7565b905f90815b855183101561250c5760016124aa6124a48589612093565b516121ec565b6124b3816121e2565b14612503576124e0906124da60406124cb8686612093565b5101516001600160401b031690565b90611ba0565b916124ec858585612b09565b6124fb576001905b019161248c565b505050505050565b916001906124f4565b90506101fd945061251f93929150612b09565b6123f2565b60200151606001519051909161253982612ab7565b5f5f935b85518510156125f0576125508587612093565b516001815161255e816121e2565b612567816121e2565b146125e557602001515160208151910120955f5b82518110156125db578761258f8285612093565b515160208151910120146125a55760010161257b565b82975060406124cb6124da926125be9599969499612093565b905b6125cb848484612b09565b6124fb576001905b01939461253d565b50949095506125c0565b5094936001906125d3565b506101fd945061251f9350612b09565b6040516080919061261183826101cd565b6041815291601f1901366020840137565b604051906126316020836101cd565b5f808352366020840137565b906126478261022d565b61265460405191826101cd565b8281528092611fa5601f199161022d565b6020810180516001600160401b0316602481612700575b612686915061263d565b91600a6020840153602260218401536126b6906126ad6126a7600286612b63565b85612b7a565b90519084612bca565b906126cb611e0882516001600160401b031690565b6126d457505090565b6126f5611e086126e76126fc9486612b90565b92516001600160401b031690565b9083612c1a565b5090565b5061270d61271291612b38565b6112d7565b60240180602411610efa576126869061267c565b6040513d5f823e3d90fd5b805160018101809111610efa576127479061263d565b8051156112c457612771816127646020945f868196015382612c53565b5060405191828092611b25565b039060025afa15612781575f5190565b612726565b9092919280840393808511610efa576001851461280f5760015b8060011b90868210156127b357506127a0565b91929394955050820191828111610efa57826127cf9185612786565b916127da9293612786565b6127e2612600565b918251156112c457825f926127719260016020809701536021830152604182015260405191828092611b25565b509061281c929350612093565b5190565b612120610cc191612ca8565b906129e6612120610240610cc194612842611f5a565b946128506121208351612d02565b61285987611faf565b5261286386611fbc565b52612886612120612881611e0860408501516001600160401b031690565b612de2565b61288f86611fcc565b526128af6121206128aa60608401516001600160801b031690565b612e16565b6128b886611fdc565b52608081015115612a91576128d361212060a0830151612ee1565b6128dc86611fec565b5260c081015115612a6b576128f761212060e0830151612fa7565b61290086611ffc565b5261010081015115612a455761291d612120610120830151612fa7565b6129268661200c565b52612938612120610140830151612fa7565b6129418661201c565b52612953612120610160830151612fa7565b61295c8661202d565b5261296e612120610180830151612fa7565b6129778661203e565b526129896121206101a0830151612fa7565b6129928661204f565b526101c081015115612a1f576129af6121206101e0830151612fa7565b6129b886612060565b52610200810151156129f9576129d5612120610220830151612fa7565b6129de86612071565b520151612ca8565b6129ef82612082565b525f815191612786565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d6129d5565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d6129af565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d61291d565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d6128f7565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d6128d3565b5f9190825b8151841015612af157612ae96001916001600160401b036040612adf8887612093565b5101511690611ba0565b930192612abc565b925050565b81810292918115918404141715610efa57565b906001600160401b0360ff612b2a612b349483836020890151169116612af6565b9451169116612af6565b1090565b906001915b6080811015612b495750565b60019060071c920191612b3d565b6020600a910153600190565b602082600a9201015360018101809111610efa5790565b602082819201015360018101809111610efa5790565b60208260109201015360018101809111610efa5790565b60206008910153600190565b60208260129201015360018101809111610efa5790565b8160209193929301015260208101809111610efa5790565b91906021600193015b6080821015612bff57906001929391530190565b600180916080607f85161781530193019060071c9092612beb565b9092919083016020015b6080821015612c3857906001929391530190565b600180916080607f85161781530193019060071c9092612c24565b908051918215612c75576021602084930191015e60010180600111610efa5790565b505050600190565b908092918251928315612ca157839260208092019201015e8101809111610efa5790565b5050505090565b805115612cf957612cb98151612b38565b806001019081600111610efa5760019083510101809111610efa57612ce06126fc9161263d565b91600a6020840153612cf3815184612be2565b83612c7d565b50610cc1612622565b5f90612d1581516001600160401b031690565b6001600160401b038116612dbd575b50612d4d6020820192612d41611e0885516001600160401b031690565b80612da1575b5061263d565b915f91612d64611e0882516001600160401b031690565b612d7e575b506126cb611e0882516001600160401b031690565b612d9a919250612d93611e086126e786612ba7565b9084612c1a565b905f612d69565b90612db161270d612db793612b38565b906112f3565b5f612d47565b612ddb919250612dd661270d916001600160401b031690565b612b38565b905f612d24565b8015612cf957612df181612b38565b60010180600111610efa57612e086126fc9161263d565b916008602084015382612be2565b6001600160801b0316633b9aca008104906001600160801b035f9216918215159182612ebd575b633b9aca006001600160801b03910616908115159081612ea2575b612e619061263d565b935f93612e86575b50612e7357505090565b612e806126fc9284612b90565b83612c1a565b612e9b91935060086020860153600185612c1a565b915f612e69565b612eae61270d84612b38565b810180911115612e5857610ec6565b90506001600160801b03633b9aca00612ed861270d86612b38565b92915050612e3d565b6020810151612ef663ffffffff825116612b38565b90816001019182600111610efa57602301809211610efa57612f57612f1d6126fc9361263d565b91600860208401536020612f4d6126a7612f47612f41610fc5865163ffffffff1690565b87612be2565b86612bb3565b9101519083612bca565b50612cf38151612fa1612f47612f85612f8084612f7b612f7682612b38565b6112e5565b6112f3565b61263d565b96612f98612f9289612b57565b89612b7a565b90519088612bca565b85612c1a565b60405190606090612fb882846101cd565b602283526126fc91600a906020850190601f19013682375360206021840153600283612bca56fea164736f6c634300081c000a",
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
