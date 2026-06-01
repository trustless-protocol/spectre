// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractUpdateClient

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

// IUpdateClientMsgsMsgUpdateClient is an auto generated low-level Go binding around an user-defined struct.
type IUpdateClientMsgsMsgUpdateClient struct {
	ClientState           IICS07TendermintMsgsClientState
	TrustedConsensusState IICS07TendermintMsgsConsensusState
	ProposedHeader        IICS07TendermintMsgsHeader
	Time                  *big.Int
	Proof                 [8]*big.Int
	Commitments           [2]*big.Int
	CommitmentPok         [2]*big.Int
	Bucket                uint16
	SignerIndices         []uint32
	SignerPubkeys         [][32]byte
	TimestampSeconds      []uint64
	TimestampNanos        []uint32
	Active                []bool
}

// IUpdateClientMsgsUpdateClientOutput is an auto generated low-level Go binding around an user-defined struct.
type IUpdateClientMsgsUpdateClientOutput struct {
	ClientState           IICS07TendermintMsgsClientState
	TrustedConsensusState IICS07TendermintMsgsConsensusState
	NewConsensusState     IICS07TendermintMsgsConsensusState
	Time                  *big.Int
	TrustedHeight         IICS02ClientMsgsHeight
	NewHeight             IICS02ClientMsgsHeight
}

// ContractUpdateClientMetaData contains all meta data concerning the ContractUpdateClient contract.
var ContractUpdateClientMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.MsgUpdateClient\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"timestampSeconds\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"timestampNanos\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.UpdateClientOutput\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"newConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"newHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateClientCachedCurrent\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.MsgUpdateClient\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"timestampSeconds\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"timestampNanos\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.UpdateClientOutput\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"newConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"newHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateClientResolved\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.MsgUpdateClient\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"timestampSeconds\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"timestampNanos\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.UpdateClientOutput\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"newConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"newHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"HeaderChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ValSetHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x6080806040523460155761275e908161001a8239f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c8063155bbe881461071c578063513a5c0e146103125763b4e415791461003a575f80fd5b3461030e5761004836610a82565b610050610d3d565b50610059610d3d565b5061006d6100678280610e0f565b80610e25565b909161009c61008760606100818480610e0f565b01610e57565b936040519361009585610c63565b3691610e6b565b82526001600160401b03602083019316835260206100ba8280610e0f565b019163ffffffff6100d660a06100d08580610e0f565b01610eb0565b6100ed604051956100e687610c92565b3690610ecf565b8552166020840152600f60408401526101096080830183610f05565b9361012161011960a08501610f1a565b95369061117f565b61012e366020860161154a565b9261015161014160208451510151611abf565b61014b8185611d36565b82611d9e565b51926001600160801b038151169160408101926001600160401b036020855101511692606083019360408551920151926040519861018e8a610cad565b89526020890152604088015260608701526080860152805192602082018051604051956101ba87610c63565b86526020860152515151156102a3576001600160401b0360406101e48260208195510151166116d0565b935151015116911614159081610297575b5061022c576102289561021c956001600160401b03946102149361234f565b5116906118e0565b60405191829182610ada565b0390f35b608460405163f492ef2b60e01b815260206004820152602760248201527f74727573746564206e6578742076616c696461746f7220736574206e6f74207260448201527f65736f6c766564000000000000000000000000000000000000000000000000006064820152fd5b9050515151155f6101f5565b608460405163f492ef2b60e01b815260206004820152602360248201527f70726f706f7365642076616c696461746f7220736574206e6f74207265736f6c60448201527f76656400000000000000000000000000000000000000000000000000000000006064820152fd5b5f80fd5b3461030e5761032036610a82565b610328610d3d565b50610331610d3d565b5061033f6100678280610e0f565b9161036061035260606100818480610e0f565b926040519461009586610c63565b83526001600160401b036020840192168252602061037e8280610e0f565b019263ffffffff61039460a06100d08580610e0f565b6103a4604051966100e688610c92565b8652166020850152600f60408501526103c06080830183610f05565b926103d86103d060a08501610f1a565b94369061117f565b916103e6366020860161154a565b946103f660208551510151611abf565b9660208801966001600160401b038851169760408701986001600160401b038a515116908181036106ee5750506001600160401b0390511661045f6001600160401b036040895151015116916040519061044f82610c63565b8152602081019283528a5161218c565b6106b857507317435cce3d1b4fa2e5f8a08ed921d57c6762a1809860208701906104a2602083518d604051808095819463c87f1f6960e01b835260048301611762565b03915af490811561064d575f91610686575b5061014089515101518082036106585750506104d09086611d9e565b606087019660206104f989519c604051809e8193829163c87f1f6960e01b835260048301611762565b03915af49a8b1561064d575f9b610619575b5060408301518b03610582576102289a61021c9a6001600160401b03998a60206001600160801b036102149b5198511693510151169051916040519661055088610cad565b875260208701526040860152606085015260808401525190516040519161057683610c63565b8252602082015261234f565b60405163f492ef2b60e01b815260206004820152604360248201527f74727573746564206e6578742076616c696461746f722073657420686173682060448201527f646f6573206e6f74206d6174636820686173682073746f726564206f6e20636860648201527f61696e000000000000000000000000000000000000000000000000000000000060848201528060a481015b0390fd5b909a506020813d602011610645575b8161063560209383610cfe565b8101031261030e5751998b61050b565b3d9150610628565b6040513d5f823e3d90fd5b7f4d6d8014000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b90506020813d6020116106b0575b816106a160209383610cfe565b8101031261030e57518c6104b4565b3d9150610694565b6001600160401b039051167f90f4dbed000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b3461030e5761072a36610a82565b610732610d3d565b506107406100678280610e0f565b61076061075260606100818680610e0f565b926040519261009584610c63565b81526001600160401b036020820192168252602061077e8480610e0f565b0161078e60a06100d08680610e0f565b61079e604051926100e684610c92565b825263ffffffff602083019116815260408201600f81526107c26080870187610f05565b906107da6107d260a08901610f1a565b92369061117f565b946107e83660208a0161154a565b9361080b6107fb60208951510151611abf565b610805818a611d36565b83611d9e565b7317435cce3d1b4fa2e5f8a08ed921d57c6762a1806108426020895151604051809381926395b25f7960e01b835260048301611583565b0381855af490811561064d575f91610a50575b50604060208a510151015151036109fe57849260408901916108846001600160401b03602085510151166116d0565b6001600160401b038060408d515101511691160361097d575b50519580516001600160801b03169151602001516001600160401b031660608a0151916040015192604051986108d28a610cad565b8952602089015260408801526060870152608086015286519660200151604051976108fc89610c63565b885260208801525163ffffffff1690610916918588611e1b565b5163ffffffff1661092690611800565b61092f91611824565b835151606001516001600160801b0316906001600160801b03161161095390611844565b61095c92612143565b516001600160401b031661096f916118e0565b604051610228819282610ada565b6109a791945060209060608b015190604051808095819463c87f1f6960e01b835260048301611762565b03915af490811561064d575f916109cc575b506040870151036105825784928b61089d565b90506020813d6020116109f6575b816109e760209383610cfe565b8101031261030e57518b6109b9565b3d91506109da565b60405163f492ef2b60e01b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d616044820152620e8c6d60eb1b6064820152608490fd5b90506020813d602011610a7a575b81610a6b60209383610cfe565b8101031261030e57518b610855565b3d9150610a5e565b602060031982011261030e57600435906001600160401b03821161030e5761030090829003600319011261030e5760040190565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b9060208252805191610180602082015260c0610b0584516101206101a08501526102c0840190610ab6565b9360ff602080830151828151166101c08701520151166101e0840152610b4960408201516102008501906001600160401b0360208092828151168552015116910152565b63ffffffff60608201511661024084015263ffffffff60808201511661026084015260a081015115156102808401520151916002831015610c4f5760a061016091610c4c946102a0850152610bc360208201516040860190604080916001600160801b038151168452602081015160208501520151910152565b610bf1604082015183860190604080916001600160801b038151168452602081015160208501520151910152565b6001600160801b03606082015116610100850152610c2d60808201516101208601906001600160401b0360208092828151168552015116910152565b01519101906001600160401b0360208092828151168552015116910152565b90565b634e487b7160e01b5f52602160045260245ffd5b604081019081106001600160401b03821117610c7e57604052565b634e487b7160e01b5f52604160045260245ffd5b606081019081106001600160401b03821117610c7e57604052565b60a081019081106001600160401b03821117610c7e57604052565b60e081019081106001600160401b03821117610c7e57604052565b608081019081106001600160401b03821117610c7e57604052565b90601f801991011681019081106001600160401b03821117610c7e57604052565b60405190610d2c82610c92565b5f6040838281528260208201520152565b6040519060c082018281106001600160401b03821117610c7e5760405281604051610d6781610cc8565b60608152604051610d7781610c63565b5f81525f60208201526020820152604051610d9181610c63565b5f81525f602082015260408201525f60608201525f60808201525f60a08201525f60c08201528152610dc1610d1f565b6020820152610dce610d1f565b60408201525f6060820152604051610de581610c63565b5f81525f6020820152608082015260a060405191610e0283610c63565b5f83525f60208401520152565b90359061011e198136030182121561030e570190565b903590601e198136030182121561030e57018035906001600160401b03821161030e5760200191813603831361030e57565b356001600160401b038116810361030e5790565b9291926001600160401b038211610c7e5760405191610e94601f8201601f191660200184610cfe565b82948184528183011161030e578281602093845f960137010152565b3563ffffffff8116810361030e5790565b359060ff8216820361030e57565b919082604091031261030e57604051610ee781610c63565b6020610f00818395610ef881610ec1565b855201610ec1565b910152565b903590609e198136030182121561030e570190565b356001600160801b038116810361030e5790565b35906001600160401b038216820361030e57565b919082604091031261030e57604051610f5a81610c63565b6020610f00818395610f6b81610f2e565b855201610f2e565b9080601f8301121561030e57816020610c4c93359101610e6b565b35906001600160801b038216820361030e57565b3590811515820361030e57565b359063ffffffff8216820361030e57565b80929103916060831261030e57604051610fd981610c63565b6040819483358352601f19011261030e576020906040805193610ffb85610c63565b611006848201610faf565b85520135828401520152565b6001600160401b038111610c7e5760051b60200190565b919060808382031261030e576040519061104282610ce3565b819380356001600160401b03811161030e57606092611062918301610f73565b83526020810135602084015261107a60408201610f2e565b60408401520135908160070b820361030e5760600152565b91909160808184031261030e57604051906110ac82610ce3565b819381356001600160401b03811161030e57820181601f8201121561030e5780356110d681611012565b916110e46040519384610cfe565b81835260208084019260051b8201019184831161030e5760208201905b8382106111525750505050835261111a60208301610fa2565b60208401526040820135906001600160401b03821161030e578261114760609492610f0094869401611029565b604086015201610f2e565b81356001600160401b03811161030e5760209161117488848094880101611029565b815201910190611101565b919060a08382031261030e576040519061119882610ce3565b819380356001600160401b03811161030e57810160408184031261030e57604051906111c382610c63565b80356001600160401b03811161030e5781016102c08186031261030e576040519061026082018281106001600160401b03821117610c7e576040526112088682610f42565b825260408101356001600160401b03811161030e5786611229918301610f73565b602083015261123a60608201610f2e565b604083015261124b60808201610f8e565b606083015261125c60a08201610fa2565b608083015261126e8660c08301610fc0565b60a08301526112806101208201610fa2565b60c083015261014081013560e083015261129d6101608201610fa2565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526112ec6102208201610fa2565b6101c08301526102408101356101e083015261130b6102608201610fa2565b6102008301526102808101356102208301526102a0810135906001600160401b03821161030e5761133e91879101610f73565b61024082015282526020810135906001600160401b03821161030e570160c08185031261030e576040519061137282610ce3565b61137b81610f2e565b825261138960208201610faf565b602083015261139b8560408301610fc0565b604083015260a0810135906001600160401b03821161030e570184601f8201121561030e578035906113cc82611012565b916113da6040519384610cfe565b80835260208084019160051b8301019187831161030e5760208101915b838310611465575050505060608201526020820152835260208101356001600160401b03811161030e578261142d918301611092565b602084015261143f8260408301610f42565b60408401526080810135916001600160401b03831161030e57606092610f009201611092565b82356001600160401b03811161030e578201906040828b03601f19011261030e576040519161149383610c63565b6020810135600481101561030e57835260408101356001600160401b03811161030e576020910101906080828c031261030e57604051926114d384610ce3565b82356001600160401b03811161030e578c6114ef918501610f73565b84526114fd60208401610f8e565b602085015261150e60408401610fa2565b60408501526060830135936001600160401b03851161030e576115368d602096879601610f73565b6060820152838201528152019201916113f7565b919082606091031261030e5760405161156281610c92565b604080829461157081610f8e565b8452602081013560208501520135910152565b90610c4c91602081526115af6020820183516001600160401b0360208092828151168552015116910152565b6102406115cd60208401516102c060608501526102e0840190610ab6565b926001600160401b0360408201511660808401526001600160801b0360608201511660a08401526080810151151560c084015260208060a0830151805160e0870152015163ffffffff815116610100860152015161012084015260c0810151151561014084015260e081015161016084015261010081015115156101808401526101208101516101a08401526101408101516101c08401526101608101516101e08401526101808101516102008401526101a08101516102208401526101c08101511515828401526101e081015161026084015261020081015115156102808401526102208101516102a08401520151906102c0601f1982850301910152610ab6565b6001600160401b036001911601906001600160401b0382116116ee57565b634e487b7160e01b5f52601160045260245ffd5b906001600160401b03809116911601906001600160401b0382116116ee57565b906060806117398451608085526080850190610ab6565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b6020815260a0810182519060806020840152815180915260c0830190602060c08260051b8601019301915f905b8282106117d557505050506001600160401b0360606117cb6080936020870151151560408701526040870151601f198783030184880152611722565b9401511691015290565b909192936020806117f260019360bf198a82030186528851611722565b96019201920190929161178f565b6001600160801b03633b9aca00911602906001600160801b0382169182036116ee57565b906001600160801b03809116911601906001600160801b0382116116ee57565b1561184b57565b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b903590603e198136030182121561030e570190565b9035906102be198136030182121561030e570190565b91906118ea610d3d565b926118f58180610e0f565b6101208136031261030e576040519061190d82610cc8565b8035906001600160401b03821161030e5761192e6101009236908301610f73565b835261193d3660208301610ecf565b602084015261194f3660608301610f42565b604084015261196060a08201610faf565b606084015261197160c08201610faf565b608084015261198260e08201610fa2565b60a08401520135600281101561030e576100816119e46119de6001600160401b03958560609560c0611a7c9801528a526119bf366020830161154a565b60208b015260808101906119f060806119ea6119e46119de8686610f05565b806118b5565b806118ca565b01610f1a565b610200611a036119e46119de8686610f05565b01356101c0611a186119e46119de8787610f05565b0135906001600160801b0360405193611a3085610c92565b1683526020830152604082015260408c01526001600160801b03611a5660a08301610f1a565b16868c0152611a726040611a6a8484610f05565b013690610f42565b60808c0152610f05565b8160405193611a8a85610c63565b16835216602082015260a0830152565b908151811015611aab570160200190565b634e487b7160e01b5f52603260045260245ffd5b604051611acb81610c63565b606081525f60208201525080518015908115611d2a575b50611d02575f19908051805b611c8a575b505f198214611c7c5760018201908183116116ee577f30000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000611b548484611a9a565b51161480611c67575b611c57575f5b8151831015611bed57611b768383611a9a565b5160f81c603081108015611be3575b611bd1576001600160401b0360ff8192602f19011681600a850216011691168110611bb557600190920191611b63565b5091505060405190611bc682610c63565b81525f602082015290565b505091505060405190611bc682610c63565b5060398111611b85565b9150916001811190811591611c4b575b50611c23576001600160401b039060405192611c1884610c63565b835216602082015290565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f611bfd565b91505060405190611bc682610c63565b5080518381039081116116ee57600210611b5d565b6040519150611bc682610c63565b5f1981018181116116ee577f2d000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000611ce18386611a9a565b511614611cf8575080156116ee575f190180611aee565b92505f9050611af3565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f611ae2565b90602001906001600160401b038251169160408201926001600160401b0384515116908181036106ee575050906001600160401b03604081611d979451169251510151169260405191611d8883610c63565b8252602082019384525161218c565b6106b85750565b906001600160401b036020830151166001600160401b0360208301511603611dc4575050565b90610615611e09925191516040519384937fc6913e9f000000000000000000000000000000000000000000000000000000008552604060048601526044850190610ab6565b83810360031901602485015290610ab6565b916001600160401b03611e319194929416611800565b9060208401916001600160801b038084511692169180831092831561211c575b5050506120b2576001600160801b0380606084515101511691511610156120485760208151510151604051611ea56020828180820195805191829101875e81015f838201520301601f198101835282610cfe565b5190208251604051611ed66020828180820195805191829101875e81015f838201520301601f198101835282610cfe565b5190200361200457611ef46001600160401b036040840151166116d0565b815151604001516001600160401b039182169116818103611f8f5750506101406080915151015191015103611f2557565b608460405162461bcd60e51b815260206004820152602f60248201527f696e76616c696420626c6f636b3a206e6578742076616c696461746f7220736560448201527f742068617368206d69736d6174636800000000000000000000000000000000006064820152fd5b915091501115611f9b57565b608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152fd5b0391506001600160801b0382116116ee576001600160801b038091169116115f8080611e51565b909161215b6001600160401b036040850151166116d0565b9151916001600160401b0380604085510151169116146121875760606121859301519051916121ca565b565b505050565b90612196916126a6565b6003811015610c4f57600281149081156121ae575090565b600191501490565b8051821015611aab5760209160051b010190565b60200151606001519051925f91825b855184101561220e576122066001916001600160401b0360406121fc888b6121b6565b5101511690611702565b9301926121d9565b9093919492505f5f935b85518510156122cf5761222b85876121b6565b5180516004811015610c4f576001146122c457602001515160208151910120955f5b82518110156122ba578761226182856121b6565b515160208151910120146122775760010161224d565b82975060406121fc6001600160401b039261229595999694996121b6565b905b6122a2848484612722565b6122b2576001905b019394612218565b505050505050565b5094909550612297565b5094936001906122aa565b50909192506122de9350612722565b156122e557565b608460405162461bcd60e51b815260206004820152602160248201527f696e73756666696369656e7420766f74696e6720706f776572206f7665726c6160448201527f70000000000000000000000000000000000000000000000000000000000000006064820152fd5b939091926123756020865151604051809381926395b25f7960e01b835260048301611583565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af490811561064d575f91612674575b50604060208751015101515103612623578451936060602080880151960151015193845186515103612593575f965f5b86518110156124c4576123df81886121b6565b519889516004811015610c4f576001146124ba57506001986020612404838b516121b6565b515160405161243183828180820195805191829101875e81015f838201520301601f198101835282610cfe565b519020910151516040516124646020828180820195805191829101875e81015f838201520301601f198101835282610cfe565b51902003612476576001905b016123cc565b606460405162461bcd60e51b815260206004820152601d60248201527f696e76616c696420636f6d6d69743a206661756c7479207369676e65720000006044820152fd5b9850600190612470565b509295509295909350156125295761252461250b856124f26121859763ffffffff6020890151168787611e1b565b61250563ffffffff604088015116611800565b90611824565b6001600160801b03806060855151015116911611611844565b612143565b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420636f6d6d69743a206e6f2070726573656e74207369676e6160448201527f74757265730000000000000000000000000000000000000000000000000000006064820152fd5b60a460405162461bcd60e51b815260206004820152604860248201527f696e76616c696420636f6d6d69743a206e756d626572206f66207369676e617460448201527f7572657320646f6573206e6f74206d61746368206e756d626572206f6620766160648201527f6c696461746f72730000000000000000000000000000000000000000000000006084820152fd5b60405162461bcd60e51b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d616044820152620e8c6d60eb1b6064820152608490fd5b90506020813d60201161269e575b8161268f60209383610cfe565b8101031261030e57515f61239c565b3d9150612682565b6001600160401b038151166001600160401b03835116908181105f146126ce57505050505f90565b11156126db575050600290565b60206001600160401b0381819301511692015116908181105f146126ff5750505f90565b111561270a57600290565b600190565b818102929181159184041417156116ee57565b906001600160401b0360ff61274361274d948383602089015116911661270f565b945116911661270f565b109056fea164736f6c634300081c000a",
}

// ContractUpdateClientABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractUpdateClientMetaData.ABI instead.
var ContractUpdateClientABI = ContractUpdateClientMetaData.ABI

// ContractUpdateClientBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractUpdateClientMetaData.Bin instead.
var ContractUpdateClientBin = ContractUpdateClientMetaData.Bin

// DeployContractUpdateClient deploys a new Ethereum contract, binding an instance of ContractUpdateClient to it.
func DeployContractUpdateClient(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractUpdateClient, error) {
	parsed, err := ContractUpdateClientMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractUpdateClientBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractUpdateClient{ContractUpdateClientCaller: ContractUpdateClientCaller{contract: contract}, ContractUpdateClientTransactor: ContractUpdateClientTransactor{contract: contract}, ContractUpdateClientFilterer: ContractUpdateClientFilterer{contract: contract}}, nil
}

// ContractUpdateClient is an auto generated Go binding around an Ethereum contract.
type ContractUpdateClient struct {
	ContractUpdateClientCaller     // Read-only binding to the contract
	ContractUpdateClientTransactor // Write-only binding to the contract
	ContractUpdateClientFilterer   // Log filterer for contract events
}

// ContractUpdateClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractUpdateClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractUpdateClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractUpdateClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractUpdateClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractUpdateClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractUpdateClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractUpdateClientSession struct {
	Contract     *ContractUpdateClient // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ContractUpdateClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractUpdateClientCallerSession struct {
	Contract *ContractUpdateClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// ContractUpdateClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractUpdateClientTransactorSession struct {
	Contract     *ContractUpdateClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractUpdateClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractUpdateClientRaw struct {
	Contract *ContractUpdateClient // Generic contract binding to access the raw methods on
}

// ContractUpdateClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractUpdateClientCallerRaw struct {
	Contract *ContractUpdateClientCaller // Generic read-only contract binding to access the raw methods on
}

// ContractUpdateClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractUpdateClientTransactorRaw struct {
	Contract *ContractUpdateClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractUpdateClient creates a new instance of ContractUpdateClient, bound to a specific deployed contract.
func NewContractUpdateClient(address common.Address, backend bind.ContractBackend) (*ContractUpdateClient, error) {
	contract, err := bindContractUpdateClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractUpdateClient{ContractUpdateClientCaller: ContractUpdateClientCaller{contract: contract}, ContractUpdateClientTransactor: ContractUpdateClientTransactor{contract: contract}, ContractUpdateClientFilterer: ContractUpdateClientFilterer{contract: contract}}, nil
}

// NewContractUpdateClientCaller creates a new read-only instance of ContractUpdateClient, bound to a specific deployed contract.
func NewContractUpdateClientCaller(address common.Address, caller bind.ContractCaller) (*ContractUpdateClientCaller, error) {
	contract, err := bindContractUpdateClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractUpdateClientCaller{contract: contract}, nil
}

// NewContractUpdateClientTransactor creates a new write-only instance of ContractUpdateClient, bound to a specific deployed contract.
func NewContractUpdateClientTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractUpdateClientTransactor, error) {
	contract, err := bindContractUpdateClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractUpdateClientTransactor{contract: contract}, nil
}

// NewContractUpdateClientFilterer creates a new log filterer instance of ContractUpdateClient, bound to a specific deployed contract.
func NewContractUpdateClientFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractUpdateClientFilterer, error) {
	contract, err := bindContractUpdateClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractUpdateClientFilterer{contract: contract}, nil
}

// bindContractUpdateClient binds a generic wrapper to an already deployed contract.
func bindContractUpdateClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractUpdateClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractUpdateClient *ContractUpdateClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractUpdateClient.Contract.ContractUpdateClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractUpdateClient *ContractUpdateClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.ContractUpdateClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractUpdateClient *ContractUpdateClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.ContractUpdateClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractUpdateClient *ContractUpdateClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractUpdateClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractUpdateClient *ContractUpdateClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractUpdateClient *ContractUpdateClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.contract.Transact(opts, method, params...)
}

// UpdateClient is a free data retrieval call binding the contract method 0x513a5c0e.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCaller) UpdateClient(opts *bind.CallOpts, msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	var out []interface{}
	err := _ContractUpdateClient.contract.Call(opts, &out, "updateClient", msg_)

	if err != nil {
		return *new(IUpdateClientMsgsUpdateClientOutput), err
	}

	out0 := *abi.ConvertType(out[0], new(IUpdateClientMsgsUpdateClientOutput)).(*IUpdateClientMsgsUpdateClientOutput)

	return out0, err

}

// UpdateClient is a free data retrieval call binding the contract method 0x513a5c0e.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientSession) UpdateClient(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClient(&_ContractUpdateClient.CallOpts, msg_)
}

// UpdateClient is a free data retrieval call binding the contract method 0x513a5c0e.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCallerSession) UpdateClient(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClient(&_ContractUpdateClient.CallOpts, msg_)
}

// UpdateClientCachedCurrent is a free data retrieval call binding the contract method 0x155bbe88.
//
// Solidity: function updateClientCachedCurrent(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCaller) UpdateClientCachedCurrent(opts *bind.CallOpts, msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	var out []interface{}
	err := _ContractUpdateClient.contract.Call(opts, &out, "updateClientCachedCurrent", msg_)

	if err != nil {
		return *new(IUpdateClientMsgsUpdateClientOutput), err
	}

	out0 := *abi.ConvertType(out[0], new(IUpdateClientMsgsUpdateClientOutput)).(*IUpdateClientMsgsUpdateClientOutput)

	return out0, err

}

// UpdateClientCachedCurrent is a free data retrieval call binding the contract method 0x155bbe88.
//
// Solidity: function updateClientCachedCurrent(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientSession) UpdateClientCachedCurrent(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClientCachedCurrent(&_ContractUpdateClient.CallOpts, msg_)
}

// UpdateClientCachedCurrent is a free data retrieval call binding the contract method 0x155bbe88.
//
// Solidity: function updateClientCachedCurrent(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCallerSession) UpdateClientCachedCurrent(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClientCachedCurrent(&_ContractUpdateClient.CallOpts, msg_)
}

// UpdateClientResolved is a free data retrieval call binding the contract method 0xb4e41579.
//
// Solidity: function updateClientResolved(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCaller) UpdateClientResolved(opts *bind.CallOpts, msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	var out []interface{}
	err := _ContractUpdateClient.contract.Call(opts, &out, "updateClientResolved", msg_)

	if err != nil {
		return *new(IUpdateClientMsgsUpdateClientOutput), err
	}

	out0 := *abi.ConvertType(out[0], new(IUpdateClientMsgsUpdateClientOutput)).(*IUpdateClientMsgsUpdateClientOutput)

	return out0, err

}

// UpdateClientResolved is a free data retrieval call binding the contract method 0xb4e41579.
//
// Solidity: function updateClientResolved(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientSession) UpdateClientResolved(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClientResolved(&_ContractUpdateClient.CallOpts, msg_)
}

// UpdateClientResolved is a free data retrieval call binding the contract method 0xb4e41579.
//
// Solidity: function updateClientResolved(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCallerSession) UpdateClientResolved(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClientResolved(&_ContractUpdateClient.CallOpts, msg_)
}
