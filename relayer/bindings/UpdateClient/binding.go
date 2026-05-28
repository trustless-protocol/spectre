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
	ABI: "[{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.MsgUpdateClient\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"timestampSeconds\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"timestampNanos\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.UpdateClientOutput\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"newConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"newHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateClientResolved\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.MsgUpdateClient\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"timestampSeconds\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"timestampNanos\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.UpdateClientOutput\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"newConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"newHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"HeaderChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ValSetHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x608080604052346015576123b4908161001a8239f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c8063513a5c0e146105915763b4e415791461002f575f80fd5b346104475761003d36610950565b610045610c0b565b5061004e610c0b565b506100598180610cdd565b803590601e198136030182121561044757018035906001600160401b03821161044757602001918136038313610447576100b66100a1606061009b8480610cdd565b01610cf3565b93604051936100af85610b31565b3691610d07565b82526001600160401b03602083019316835260206100d48280610cdd565b0160a06100e18380610cdd565b013563ffffffff8116809103610447576101086040519261010184610b60565b3690610d5a565b82526020820152600f604082015260808201906101258284610d90565b9460a084019461013486610da5565b6101436020870198369061100a565b9361014e368a6113d5565b61015d6020875151015161158e565b9360208501946001600160401b038651169560408901966001600160401b0388515116908181036105635750506001600160401b039051166101c66001600160401b0360408b515101511691604051906101b682610b31565b81526020810192835288516120f2565b61052d57506101d59082611805565b51956001600160801b038251166001600160401b036020875101511692606083019360408551920151926040519a61020c8c610bb1565b8b5260208b015260408a0152606089015260808801528051946020820180516040519761023889610b31565b88526020880152515151156104c2576001600160401b0360406102628260208195510151166114ec565b9351510151169116141590816104b6575b5061044b576001600160401b039461028a936118aa565b511691610295610c0b565b936102a08280610cdd565b906101208236031261044757604051916102b983610b7b565b8035906001600160401b038211610447576102da6101009236908301610dfe565b84526102e93660208301610d5a565b60208501526102fb3660608301610dcd565b604085015261030c60a08201610e3a565b606085015261031d60c08201610e3a565b608085015261032e60e08201610e2d565b60a08501520135906002821015610447576103996001600160401b039561041d956001600160801b036103f761039f958c60206104439f61037e908c61009b9c60c060609f0152845236906113d5565b9101528c60406103ab60806103a561039f6103998b8b610d90565b8061153e565b80611553565b01610da5565b916102006103bf61039f6103998b8b610d90565b01356101c06103d461039f6103998c8c610d90565b013590878451956103e487610b60565b1685526020850152828401520152610da5565b16868b0152610413604061040b8484610d90565b013690610dcd565b60808b0152610d90565b6040519261042a84610b31565b835216602082015260a0820152604051918291826109a8565b0390f35b5f80fd5b608460405163f492ef2b60e01b815260206004820152602760248201527f74727573746564206e6578742076616c696461746f7220736574206e6f74207260448201527f65736f6c766564000000000000000000000000000000000000000000000000006064820152fd5b9050515151155f610273565b608460405163f492ef2b60e01b815260206004820152602360248201527f70726f706f7365642076616c696461746f7220736574206e6f74207265736f6c60448201527f76656400000000000000000000000000000000000000000000000000000000006064820152fd5b6001600160401b039051167f90f4dbed000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b346104475761059f36610950565b6105a7610c0b565b506105b0610c0b565b506105bb8180610cdd565b803590601e198136030182121561044757018035906001600160401b03821161044757602001918136038313610447576105fd6100a1606061009b8480610cdd565b82526001600160401b036020830193168352602061061b8280610cdd565b0160a06106288380610cdd565b013563ffffffff8116809103610447576106486040519261010184610b60565b82526020820152600f604082015260808201906106658284610d90565b9460a084019461067486610da5565b6106836020870198369061100a565b9361068e368a6113d5565b61069d6020875151015161158e565b60208101966001600160401b038851169760408201986001600160401b038a515116908181036105635750506001600160401b039051166107056001600160401b03604084515101511691604051906106f582610b31565b8152602081019283528a516120f2565b61052d57507317435cce3d1b4fa2e5f8a08ed921d57c6762a180976020820192610748602085518c604051808095819463c87f1f6960e01b83526004830161144e565b03915af49081156108e5575f9161091e575b5061014084515101518082036108f05750509061077a61079f9289611805565b6020606084019a8b5190604051808096819463c87f1f6960e01b83526004830161144e565b03915af49182156108e5575f926108b1575b5060408501518203610820576001600160401b03998a60206001600160801b0361028a9b519851169351015116905191604051966107ee88610bb1565b875260208701526040860152606085015260808401525190516040519161081483610b31565b825260208201526118aa565b60a460405163f492ef2b60e01b815260206004820152604360248201527f74727573746564206e6578742076616c696461746f722073657420686173682060448201527f646f6573206e6f74206d6174636820686173682073746f726564206f6e20636860648201527f61696e00000000000000000000000000000000000000000000000000000000006084820152fd5b9091506020813d6020116108dd575b816108cd60209383610bcc565b810103126104475751908e6107b1565b3d91506108c0565b6040513d5f823e3d90fd5b7f4d6d8014000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b90506020813d602011610948575b8161093960209383610bcc565b8101031261044757518f61075a565b3d915061092c565b602060031982011261044757600435906001600160401b038211610447576103009082900360031901126104475760040190565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b9060208252805191610180602082015260c06109d384516101206101a08501526102c0840190610984565b9360ff602080830151828151166101c08701520151166101e0840152610a1760408201516102008501906001600160401b0360208092828151168552015116910152565b63ffffffff60608201511661024084015263ffffffff60808201511661026084015260a081015115156102808401520151916002831015610b1d5760a061016091610b1a946102a0850152610a9160208201516040860190604080916001600160801b038151168452602081015160208501520151910152565b610abf604082015183860190604080916001600160801b038151168452602081015160208501520151910152565b6001600160801b03606082015116610100850152610afb60808201516101208601906001600160401b0360208092828151168552015116910152565b01519101906001600160401b0360208092828151168552015116910152565b90565b634e487b7160e01b5f52602160045260245ffd5b604081019081106001600160401b03821117610b4c57604052565b634e487b7160e01b5f52604160045260245ffd5b606081019081106001600160401b03821117610b4c57604052565b60e081019081106001600160401b03821117610b4c57604052565b608081019081106001600160401b03821117610b4c57604052565b60a081019081106001600160401b03821117610b4c57604052565b90601f801991011681019081106001600160401b03821117610b4c57604052565b60405190610bfa82610b60565b5f6040838281528260208201520152565b6040519060c082018281106001600160401b03821117610b4c5760405281604051610c3581610b7b565b60608152604051610c4581610b31565b5f81525f60208201526020820152604051610c5f81610b31565b5f81525f602082015260408201525f60608201525f60808201525f60a08201525f60c08201528152610c8f610bed565b6020820152610c9c610bed565b60408201525f6060820152604051610cb381610b31565b5f81525f6020820152608082015260a060405191610cd083610b31565b5f83525f60208401520152565b90359061011e1981360301821215610447570190565b356001600160401b03811681036104475790565b9291926001600160401b038211610b4c5760405191610d30601f8201601f191660200184610bcc565b829481845281830111610447578281602093845f960137010152565b359060ff8216820361044757565b919082604091031261044757604051610d7281610b31565b6020610d8b818395610d8381610d4c565b855201610d4c565b910152565b903590609e1981360301821215610447570190565b356001600160801b03811681036104475790565b35906001600160401b038216820361044757565b919082604091031261044757604051610de581610b31565b6020610d8b818395610df681610db9565b855201610db9565b9080601f8301121561044757816020610b1a93359101610d07565b35906001600160801b038216820361044757565b3590811515820361044757565b359063ffffffff8216820361044757565b80929103916060831261044757604051610e6481610b31565b6040819483358352601f190112610447576020906040805193610e8685610b31565b610e91848201610e3a565b85520135828401520152565b6001600160401b038111610b4c5760051b60200190565b91906080838203126104475760405190610ecd82610b96565b819380356001600160401b03811161044757606092610eed918301610dfe565b835260208101356020840152610f0560408201610db9565b60408401520135908160070b82036104475760600152565b9190916080818403126104475760405190610f3782610b96565b819381356001600160401b03811161044757820181601f82011215610447578035610f6181610e9d565b91610f6f6040519384610bcc565b81835260208084019260051b820101918483116104475760208201905b838210610fdd57505050508352610fa560208301610e2d565b60208401526040820135906001600160401b0382116104475782610fd260609492610d8b94869401610eb4565b604086015201610db9565b81356001600160401b03811161044757602091610fff88848094880101610eb4565b815201910190610f8c565b919060a083820312610447576040519061102382610b96565b819380356001600160401b038111610447578101604081840312610447576040519061104e82610b31565b80356001600160401b0381116104475781016102c081860312610447576040519061026082018281106001600160401b03821117610b4c576040526110938682610dcd565b825260408101356001600160401b03811161044757866110b4918301610dfe565b60208301526110c560608201610db9565b60408301526110d660808201610e19565b60608301526110e760a08201610e2d565b60808301526110f98660c08301610e4b565b60a083015261110b6101208201610e2d565b60c083015261014081013560e08301526111286101608201610e2d565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526111776102208201610e2d565b6101c08301526102408101356101e08301526111966102608201610e2d565b6102008301526102808101356102208301526102a0810135906001600160401b038211610447576111c991879101610dfe565b61024082015282526020810135906001600160401b038211610447570160c08185031261044757604051906111fd82610b96565b61120681610db9565b825261121460208201610e3a565b60208301526112268560408301610e4b565b604083015260a0810135906001600160401b038211610447570184601f820112156104475780359061125782610e9d565b916112656040519384610bcc565b80835260208084019160051b830101918783116104475760208101915b8383106112f0575050505060608201526020820152835260208101356001600160401b03811161044757826112b8918301610f1d565b60208401526112ca8260408301610dcd565b60408401526080810135916001600160401b03831161044757606092610d8b9201610f1d565b82356001600160401b038111610447578201906040828b03601f190112610447576040519161131e83610b31565b6020810135600481101561044757835260408101356001600160401b038111610447576020910101906080828c0312610447576040519261135e84610b96565b82356001600160401b038111610447578c61137a918501610dfe565b845261138860208401610e19565b602085015261139960408401610e2d565b60408501526060830135936001600160401b038511610447576113c18d602096879601610dfe565b606082015283820152815201920191611282565b9190826060910312610447576040516113ed81610b60565b60408082946113fb81610e19565b8452602081013560208501520135910152565b906060806114258451608085526080850190610984565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b6020815260a0810182519060806020840152815180915260c0830190602060c08260051b8601019301915f905b8282106114c157505050506001600160401b0360606114b76080936020870151151560408701526040870151601f19878303018488015261140e565b9401511691015290565b909192936020806114de60019360bf198a8203018652885161140e565b96019201920190929161147b565b6001600160401b036001911601906001600160401b03821161150a57565b634e487b7160e01b5f52601160045260245ffd5b906001600160401b03809116911601906001600160401b03821161150a57565b903590603e1981360301821215610447570190565b9035906102be1981360301821215610447570190565b90815181101561157a570160200190565b634e487b7160e01b5f52603260045260245ffd5b60405161159a81610b31565b606081525f602082015250805180159081156117f9575b506117d1575f19908051805b611759575b505f19821461174b57600182019081831161150a577f30000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000006116238484611569565b51161480611736575b611726575f5b81518310156116bc576116458383611569565b5160f81c6030811080156116b2575b6116a0576001600160401b0360ff8192602f19011681600a85021601169116811061168457600190920191611632565b509150506040519061169582610b31565b81525f602082015290565b50509150506040519061169582610b31565b5060398111611654565b915091600181119081159161171a575b506116f2576001600160401b0390604051926116e784610b31565b835216602082015290565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f6116cc565b9150506040519061169582610b31565b50805183810390811161150a5760021061162c565b604051915061169582610b31565b5f19810181811161150a577f2d000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000006117b08386611569565b5116146117c75750801561150a575f1901806115bd565b92505f90506115c2565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f6115b1565b906001600160401b036020830151166001600160401b036020830151160361182b575050565b90611882611870925191516040519384937fc6913e9f000000000000000000000000000000000000000000000000000000008552604060048601526044850190610984565b83810360031901602485015290610984565b0390fd5b6001600160801b03633b9aca00911602906001600160801b03821691820361150a57565b9390919293611a2a6020825151604051809381927f95b25f7900000000000000000000000000000000000000000000000000000000835284600484015261190a6024840182516001600160401b0360208092828151168552015116910152565b610240611927868301516102c060648701526102e4860190610984565b60408301516001600160401b0316608486015260608301516001600160801b031660a48601526080830151151560c486015260a0830151805160e4870152870151805163ffffffff1661010487015287015161012486015260c0830151151561014486015260e083015161016486015261010083015115156101848601526101208301516101a48601526101408301516101c48601526101608301516101e48601526101808301516102048601526101a08301516102248601526101c083015115156102448601526101e083015161026486015261020083015115156102848601526102208301516102a4860152910151838203602319016102c4850152610984565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49081156108e5575f916120c0575b50604060208351015101515103612056578051936060602080840151960151015193845186515103611fc6575f965f5b8651811015611b7957611a9481886121cc565b519889516004811015610b1d57600114611b6f57506001986020611ab9838b516121cc565b5151604051611ae683828180820195805191829101875e81015f838201520301601f198101835282610bcc565b51902091015151604051611b196020828180820195805191829101875e81015f838201520301601f198101835282610bcc565b51902003611b2b576001905b01611a81565b606460405162461bcd60e51b815260206004820152601d60248201527f696e76616c696420636f6d6d69743a206661756c7479207369676e65720000006044820152fd5b9850600190611b25565b50929550929590935015611f5c57611b9a63ffffffff602085015116611886565b9060208301906001600160801b0380835116911692818410918215611f35575b5050611ecb576001600160801b038060608751510151169151161015611e615760208451510151604051611c0d6020828180820195805191829101875e81015f838201520301601f198101835282610bcc565b5190208251604051611c3e6020828180820195805191829101875e81015f838201520301601f198101835282610bcc565b51902003611e1d57611c5c6001600160401b036040840151166114ec565b845151604001516001600160401b039182169116818103611daf5750506101408451510151608083015103611d45575b6001600160801b03611ca763ffffffff604086015116611886565b16016001600160801b03811161150a576001600160801b0380606086515101511691161115611cdb57611cd99261211c565b565b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152602f60248201527f696e76616c696420626c6f636b3a206e6578742076616c696461746f7220736560448201527f742068617368206d69736d6174636800000000000000000000000000000000006064820152fd5b11611c8c57608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152fd5b840391506001600160801b03821161150a576001600160801b038091169116115f80611bba565b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420636f6d6d69743a206e6f2070726573656e74207369676e6160448201527f74757265730000000000000000000000000000000000000000000000000000006064820152fd5b60a460405162461bcd60e51b815260206004820152604860248201527f696e76616c696420636f6d6d69743a206e756d626572206f66207369676e617460448201527f7572657320646f6573206e6f74206d61746368206e756d626572206f6620766160648201527f6c696461746f72730000000000000000000000000000000000000000000000006084820152fd5b608460405162461bcd60e51b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152fd5b90506020813d6020116120ea575b816120db60209383610bcc565b8101031261044757515f611a51565b3d91506120ce565b906120fc91612163565b6003811015610b1d5760028114908115612114575090565b600191501490565b90916121346001600160401b036040850151166114ec565b9151916001600160401b03806040855101511691161461215e576060611cd99301519051916121e0565b505050565b6001600160401b038151166001600160401b03835116908181105f1461218b57505050505f90565b1115612198575050600290565b60206001600160401b0381819301511692015116908181105f146121bc5750505f90565b11156121c757600290565b600190565b805182101561157a5760209160051b010190565b60200151606001519051925f91825b85518410156122245761221c6001916001600160401b036040612212888b6121cc565b510151169061151e565b9301926121ef565b9093919492505f5f935b85518510156122e55761224185876121cc565b5180516004811015610b1d576001146122da57602001515160208151910120955f5b82518110156122d0578761227782856121cc565b5151602081519101201461228d57600101612263565b82975060406122126001600160401b03926122ab95999694996121cc565b905b6122b8848484612378565b6122c8576001905b01939461222e565b505050505050565b50949095506122ad565b5094936001906122c0565b50909192506122f49350612378565b156122fb57565b608460405162461bcd60e51b815260206004820152602160248201527f696e73756666696369656e7420766f74696e6720706f776572206f7665726c6160448201527f70000000000000000000000000000000000000000000000000000000000000006064820152fd5b8181029291811591840414171561150a57565b906001600160401b0360ff6123996123a39483836020890151169116612365565b9451169116612365565b109056fea164736f6c634300081c000a",
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
