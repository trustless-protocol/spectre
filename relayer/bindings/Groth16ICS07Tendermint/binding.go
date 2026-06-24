// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractGroth16ICS07Tendermint

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
	SignedHeader  IICS07TendermintMsgsSignedHeader
	TrustedHeight IICS02ClientMsgsHeight
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

// ILightClientMsgsMsgVerifyMembership is an auto generated low-level Go binding around an user-defined struct.
type ILightClientMsgsMsgVerifyMembership struct {
	Height                IICS02ClientMsgsHeight
	KvPairs               []IMembershipMsgsKVPair
	MerkleProofs          []IMembershipMsgsMerkleProof
	AppHash               [32]byte
	TrustedConsensusState IICS07TendermintMsgsConsensusState
	MembershipType        uint8
	Path                  [][]byte
	Value                 []byte
}

// ILightClientMsgsMsgVerifyNonMembership is an auto generated low-level Go binding around an user-defined struct.
type ILightClientMsgsMsgVerifyNonMembership struct {
	Height                IICS02ClientMsgsHeight
	KvPairs               []IMembershipMsgsKVPair
	MerkleProofs          []IMembershipMsgsMerkleProof
	AppHash               [32]byte
	TrustedConsensusState IICS07TendermintMsgsConsensusState
	MembershipType        uint8
	Path                  [][]byte
}

// IMembershipMsgsCommitmentProof is an auto generated low-level Go binding around an user-defined struct.
type IMembershipMsgsCommitmentProof struct {
	ProofType         uint8
	ExistenceProof    IMembershipMsgsExistenceProof
	NonExistenceProof IMembershipMsgsNonExistenceProof
}

// IMembershipMsgsExistenceProof is an auto generated low-level Go binding around an user-defined struct.
type IMembershipMsgsExistenceProof struct {
	Key   []byte
	Value []byte
	Leaf  IMembershipMsgsLeafOp
	Path  []IMembershipMsgsInnerOp
}

// IMembershipMsgsInnerOp is an auto generated low-level Go binding around an user-defined struct.
type IMembershipMsgsInnerOp struct {
	HashOp uint8
	Prefix []byte
	Suffix []byte
}

// IMembershipMsgsKVPair is an auto generated low-level Go binding around an user-defined struct.
type IMembershipMsgsKVPair struct {
	Path  [][]byte
	Value []byte
}

// IMembershipMsgsLeafOp is an auto generated low-level Go binding around an user-defined struct.
type IMembershipMsgsLeafOp struct {
	HashOp       uint8
	PrehashKey   uint8
	PrehashValue uint8
	Prefix       []byte
}

// IMembershipMsgsMerkleProof is an auto generated low-level Go binding around an user-defined struct.
type IMembershipMsgsMerkleProof struct {
	Proofs []IMembershipMsgsCommitmentProof
}

// IMembershipMsgsNonExistenceProof is an auto generated low-level Go binding around an user-defined struct.
type IMembershipMsgsNonExistenceProof struct {
	Key      []byte
	HasLeft  bool
	Left     IMembershipMsgsExistenceProof
	HasRight bool
	Right    IMembershipMsgsExistenceProof
}

// IUpdateClientMsgsMsgUpdateClient is an auto generated low-level Go binding around an user-defined struct.
type IUpdateClientMsgsMsgUpdateClient struct {
	ClientState            IICS07TendermintMsgsClientState
	TrustedConsensusState  IICS07TendermintMsgsConsensusState
	ProposedHeader         IICS07TendermintMsgsHeader
	Time                   *big.Int
	Proof                  [8]*big.Int
	Commitments            [2]*big.Int
	CommitmentPok          [2]*big.Int
	Bucket                 uint16
	SignerIndices          []uint32
	PinnedValidatorIndices []uint32
	SignerPubkeys          [][32]byte
	TimestampSeconds       []uint64
	TimestampNanos         []uint32
	Active                 []bool
}

// ContractGroth16ICS07TendermintMetaData contains all meta data concerning the ContractGroth16ICS07Tendermint contract.
var ContractGroth16ICS07TendermintMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membership_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviour_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"updateClient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_clientState\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"initialPinnedValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPinnedValidatorSet\",\"inputs\":[],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reAnchorPinnedSet\",\"inputs\":[{\"name\":\"updateMsg\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.MsgUpdateClient\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"},{\"name\":\"clockDrift\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pinnedValidatorIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"timestampSeconds\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"timestampNanos\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]},{\"name\":\"newPinnedValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"updateClientMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CannotHandleMisbehavior\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientNotFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ClientStateMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ClockDriftMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustedVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidMembershipProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NonMonotonicHeightUpdate\",\"inputs\":[{\"name\":\"latestHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"updateHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofHeightMismatch\",\"inputs\":[{\"name\":\"expectedRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"expectedRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofSignerCommitSigMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustThresholdMismatch\",\"inputs\":[{\"name\":\"expectedNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expectedDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnbondingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"UnknownZkAlgorithm\",\"inputs\":[{\"name\":\"algorithm\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorCountExceedsLimit\",\"inputs\":[{\"name\":\"validatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxValidatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"VerificationKeyMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x610160604052346101d35761785f803803809161001b8261021c565b61016039806101600161010082126101d3576100356102d1565b906100416101806102e8565b61004c6101a06102e8565b6100576101c06102e8565b6101e0516001600160401b0381116101d357846100779161016001610317565b916102005193610220519760018060401b0389116101d357608090899003126101d357604051956100a787610248565b6101608901516001600160401b0381116101d35789018161017f820112156101d3576101608101516100d88161035b565b916100e6604051938461027e565b81835260206101608185019360051b83010101918483116101d3576101808201905b8382106101d7575050505087526101226101808a016103f1565b60208801526101a0890151906001600160401b0382116101d357896101566101c092610160610161956101779e0101610386565b60408a015201610372565b60608701526101716102406102e8565b966109a0565b604051615b619081611c3e823960805181615602015260a0518161366a015260c05181610b09015260e05181818161210401526125ca015261010051818181610a000152610a48015261012051816142fe015261014051815050f35b5f80fd5b81516001600160401b0381116101d3576020916101fd8884610160819589010101610386565b815201910190610108565b634e487b7160e01b5f52604160045260245ffd5b610160601f91909101601f19168101906001600160401b0382119082101761024357604052565b610208565b608081019081106001600160401b0382111761024357604052565b604081019081106001600160401b0382111761024357604052565b601f909101601f19168101906001600160401b0382119082101761024357604052565b604051906102b16101008361027e565b565b604051906102b160408361027e565b604051906102b160808361027e565b61016051906001600160a01b03821682036101d357565b51906001600160a01b03821682036101d357565b6001600160401b03811161024357601f01601f191660200190565b81601f820112156101d357602081519101610331826102fc565b9261033f604051948561027e565b828452828201116101d357815f926020928386015e8301015290565b6001600160401b0381116102435760051b60200190565b51906001600160401b03821682036101d357565b91906080838203126101d3576040519061039f82610248565b8351919384926001600160401b0381116101d3576060926103c1918301610317565b8352602081015160208401526103d960408201610372565b60408401520151908160070b82036101d35760600152565b519081151582036101d357565b519060ff821682036101d357565b91908260409103126101d35760405161042481610263565b602061043d818395610435816103fe565b8552016103fe565b910152565b91908260409103126101d35760405161045a81610263565b602061043d81839561046b81610372565b855201610372565b519063ffffffff821682036101d357565b519060028210156101d357565b6020818303126101d3578051906001600160401b0382116101d35701610140818303126101d3576104c06102a1565b8151909290916001600160401b0383116101d357610507826104ea61012094610557968501610317565b86526104f9816020850161040c565b602087015260608301610442565b604085015261051860a08201610473565b606085015261052960c08201610473565b608085015261053a60e082016103f1565b60a085015261054c6101008201610484565b60c085015201610473565b60e082015290565b90600182811c9216801561058d575b602083101461057957565b634e487b7160e01b5f52602260045260245ffd5b91607f169161056e565b601f82116105a457505050565b5f5260205f20906020601f840160051c830193106105dc575b601f0160051c01905b8181106105d1575050565b5f81556001016105c6565b90915081906105bd565b600211156105f057565b634e487b7160e01b5f52602160045260245ffd5b9060028110156105f05769ff00000000000000000082549160481b169069ff0000000000000000001916179055565b80518051906001600160401b0382116102435761065c8261065560015461055f565b6001610597565b602090601f83116001146107fd5792610695836107c29460e0946102b1975f926107f2575b50508160011b915f199060031b1c19161790565b6001555b6020818101518051600280549284015161ff0060089190911b1660ff90921661ffff199093169290921717905560408083015180516003805492909401516fffffffffffffffff0000000000000000931b929092166001600160401b039092166001600160801b031990911617179055606081015161072f9063ffffffff1660049063ffffffff1663ffffffff19825416179055565b610767610743608083015163ffffffff1690565b60049067ffffffff0000000082549160201b169067ffffffff000000001916179055565b61079f61077760a0830151151590565b60049068ff0000000000000000825491151560401b169068ff00000000000000001916179055565b6107b760c08201516107b0816105e6565b6004610604565b015163ffffffff1690565b6004906dffffffff0000000000000000000082549160501b16906dffffffff000000000000000000001916179055565b015190505f80610681565b60015f52601f19831691905f51602061781f5f395f51905f52925f5b81811061085e57509360e0936102b19693600193836107c29810610846575b505050811b01600155610699565b01515f1960f88460031b161c191690555f8080610838565b92936020600181928786015181550195019301610819565b604051905f82600154916108898361055f565b80835292600181169081156108f957506001146108ad575b6102b19250038361027e565b5060015f90815290915f51602061781f5f395f51905f525b8183106108dd5750509060206102b1928201016108a1565b60209193508060019154838589010152019101909184926108c5565b602092506102b194915060ff191682840152151560051b8201016108a1565b15610921575050565b6355bace6f60e01b5f9081526001600160401b039182166004529116602452604490fd5b634e487b7160e01b5f52601160045260245ffd5b9063ffffffff8091169116019063ffffffff821161097357565b610945565b15610981575050565b9063ffffffff80926333fae18560e21b5f52166004521660245260445ffd5b90919293946109e96109e4610ae798977f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a612336101005260208082518301019101610491565b610633565b6109f1610876565b6020815191012061012052610a0c610a07610876565b610b47565b61014052610a7f610a66610a396020610a2b610a26610876565b610c18565b01516001600160401b031690565b60035490610a57906001600160401b03808416919081168214610918565b60401c6001600160401b031690565b6001600160401b03165f90815260056020526040902090565b556001600160a01b0390811660805290811660a05290811660c0521660e052600454610ae29063ffffffff8116610ad0610ac3605084901c63ffffffff1683610959565b9260201c63ffffffff1690565b9163ffffffff80841691161115610978565b610dfe565b600354610aff9060401c6001600160401b0316611093565b6001600160a01b038116610b265750610b166112ab565b50610b236101005161132d565b50565b80610b33610b23926111ad565b50610b3d81611223565b5061010051611386565b610b53610b5891611453565b6114ed565b90565b60405190610b6882610263565b5f602083606081520152565b8015610973575f190190565b5f1981019190821161097357565b9190820391821161097357565b634e487b7160e01b5f52603260045260245ffd5b908151811015610bc0570160200190565b610b9b565b906001820180921161097357565b603001908160301161097357565b906004820180921161097357565b90600c820180921161097357565b90602c820180921161097357565b9190820180921161097357565b610c20610b5b565b5080518015908115610df2575b50610de3575f19908051805b610d93575b505f198214610d7557600360fc1b6001600160f81b0319610c78610c6a610c6486610bc5565b85610baf565b516001600160f81b03191690565b161480610d7f575b610d75575f90610c8f83610bc5565b915b8151831015610d2657610cb0610caa610c6a8585610baf565b60f81c90565b60ff811660308110908115610d1b575b50610d0e57600a82026001600160401b03908116602f1990920160ff1691909101811691168110610cf657600190920191610c91565b50915050610d026102b3565b9081525f602082015290565b5050915050610d026102b3565b60399150115f610cc0565b9150916001811190811591610d69575b50610d5a57610b5890610d476102b3565b9283526001600160401b03166020830152565b63384c687760e21b5f5260045ffd5b602b915010155f610d36565b9050610d026102b3565b506002610d8d838351610b8e565b11610c80565b602d60f81b610dbd610db0610c6a610daa85610b80565b86610baf565b6001600160f81b03191690565b14610dd157610dcb90610b74565b80610c39565b610ddc919250610b80565b905f610c3e565b6329120bff60e21b5f5260045ffd5b6040915010155f610c2d565b610e0781611556565b9051805190811561100c5760b48211610ff3575f93610e48610e38610e33610e2e86611600565b610bd3565b611421565b93610e4285611ad4565b84611b1c565b610e5e610e57610b5386611c05565b84602e0152565b610e6783611b2c565b60305f955b8351871015610f1757610f0c849392610f07600193610ef0610e8f8c809a611542565b5191610ecf63ffffffff610ec66040860193610ec0610eb4865160018060401b031690565b6001600160401b031690565b90610c0b565b9b16868d611af8565b610ee9610edb86610be1565b91516001600160401b031690565b908b611b80565b6020610efb84610bef565b91015190890160200152565b610bfd565b960195909192610e6c565b9195506102b194610fd2946020945092610f8f9250610f5290610f436001600160401b03821115611616565b6001600160401b031684611b3a565b610f8a610f68610f628584611641565b9461177a565b600780546001600160a01b0319166001600160a01b0392909216919091179055565b600655565b8051610fc9906001600160401b031660078054600160a01b600160e01b03191660a09290921b600160a01b600160e01b0316919091179055565b015161ffff1690565b6007805461ffff60e01b191660e09290921b61ffff60e01b16919091179055565b63156f758160e31b5f52600482905260b460245260445ffd5b6305f8ded760e21b5f5260045ffd5b6009549068010000000000000000821015610243576001820180600955821015610bc05760095f52600282901c7f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af0180546001600160401b0360069490941b60c01684811b199091169390921690911b919091179055565b6001600160401b038181165f908152600860205260409020600101546006546007546001600160a01b03928316936111939361111b9261ffff60e082901c16926111109260a083901c9091169161110091166110ed6102c2565b9687526001600160a01b03166020870152565b6001600160401b03166040850152565b61ffff166060830152565b6001600160401b0384165f90815260086020526040902081518155602082015160019190910180546040840151606094909401516001600160f01b03199091166001600160a01b03939093169290921760a09390931b600160a01b600160e01b03169290921760e09190911b61ffff60e01b16179055565b6001600160a01b0316156111a45750565b6102b19061101b565b6001600160a01b0381165f9081525f51602061783f5f395f51905f52602052604090205460ff1661121e576001600160a01b03165f8181525f51602061783f5f395f51905f5260205260408120805460ff191660011790553391905f51602061779f5f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f5160206177bf5f395f51905f52602052604090205460ff1661121e576001600160a01b0381165f9081525f5160206177bf5f395f51905f5260205260409020805460ff1916600117905533906001600160a01b03165f5160206177ff5f395f51905f525f51602061779f5f395f51905f525f80a4600190565b5f80525f5160206177bf5f395f51905f526020525f5160206177df5f395f51905f525460ff16611329575f8080525f5160206177bf5f395f51905f526020525f5160206177df5f395f51905f52805460ff1916600117905533905f5160206177ff5f395f51905f525f51602061779f5f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f1461121e575f818152602081815260408083208380529091528120805460ff1916600117905533915f51602061779f5f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff166113f9575f818152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905533916001600160a01b0316905f51602061779f5f395f51905f525f80a4600190565b50505f90565b60405160809190611410838261027e565b6041815291601f1901366020840137565b9061142b826102fc565b611438604051918261027e565b8281528092611449601f19916102fc565b0190602036910137565b8051156114a857611464815161182b565b80600101908160011161097357600190835101018091116109735761148b6114a491611421565b91600a602084015361149e8151846118b1565b8361194c565b5090565b506040516114b760208261027e565b5f8152601f196114c65f6102fc565b0136602083013790565b805191908290602001825e015f815290565b6040513d5f823e3d90fd5b8051600181018091116109735761150390611421565b805115610bc05761152d816115206020945f868196015382611922565b50604051918280926114d0565b039060025afa1561153d575f5190565b6114e2565b8051821015610bc05760209160051b010190565b8051519081156113f9576115698261035b565b91611577604051938461027e565b808352601f196115868261035b565b013660208501375f5b825180518210156115f257906115e1610b5360206115af84600196611542565b5101516115dc6115d460406115c5878b51611542565b5101516001600160401b031690565b610d476102b3565b611977565b6115eb8287611542565b520161158f565b505090505f610b5892611a3a565b90602c820291808304602c149015171561097357565b1561100c57565b6040519061162a82610248565b5f6060838281528260208201528260408201520152565b91909161164c61161d565b506030835110611701576356414c34602084015160e01c0361170157602483015160c01c602c84015160f01c936116d7602e82015192604e83015160f01c936116a56116966102c2565b6001600160401b039093168352565b6116b76020830198899061ffff169052565b60408201526116ce6060820194859061ffff169052565b955161ffff1690565b9061ffff8216928315938415611737575b508315611728575b508215611713575b50506117015750565b634724a0fd60e01b5f5260045260245ffd5b51915061171f90611bc8565b14155f806116f8565b5161ffff16151592505f6116f0565b60b41093505f6116e8565b606160f81b81526001600160f01b031990911660018201526680600a3d393df360c81b6003820152610b5891600a91909101906114d0565b604051905f60208301526117a38261179560218201846114d0565b03601f19810184528361027e565b61600082511161181857506117e56117f36117d36117c3845161ffff1690565b60f01b6001600160f01b03191690565b92604051928391602083019586611742565b03601f19810183528261027e565b51905ff0906001600160a01b0382161561180957565b63fbad885d60e01b5f5260045ffd5b516331734c2160e21b5f5260045260245ffd5b906001915b608081101561183c5750565b60019060071c920191611830565b6020600a910153600190565b602082602292010153600181018091116109735790565b602082600a92010153600181018091116109735790565b6020828192010153600181018091116109735790565b602082601092010153600181018091116109735790565b91906021600193015b60808210156118ce57906001929391530190565b600180916080607f85161781530193019060071c90926118ba565b9092919083016020015b608082101561190757906001929391530190565b600180916080607f85161781530193019060071c90926118f3565b908051918215611944576021602084930191015e600101806001116109735790565b505050600190565b90809291825192831561197057839260208092019201015e81018091116109735790565b5050505090565b6020810180516024906001600160401b031680611a12575b5061199c6119ca91611421565b926119c16119bb6119b56119af8761184a565b87611856565b8661186d565b85611884565b90519084611bed565b81519091906119e1906001600160401b0316610eb4565b6119ea57505090565b611a0b610eb46119fd6114a4948661189a565b92516001600160401b031690565b90836118e9565b611a1c915061182b565b6001018060011161097357602401806024116109735761199c61198f565b90929192808403938085116109735760018514611ac35760015b8060011b9086821015611a675750611a54565b919293949550508201918281116109735782611a839185611a3a565b91611a8e9293611a3a565b611a966113ff565b91825115610bc057825f9261152d92600160208097015360218301526041820152604051918280926114d0565b5090611ad0929350611542565b5190565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602d908260081c602c8201530153565b604f5f9182604e8201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b61ffff16602c810290808204602c149015171561097357603001806030116109735790565b81602091939293010152602081018091116109735790565b60405190606090611c16828461027e565b602283526114a491600a906020850190601f19013682375360206021840153600283611bed56fe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a7146101245780630bece3561461011f578063248a9ca31461011a5780632f2ff15d1461011557806336568abe146101105780636a28f0001461010b5780636edfe3af146101065780638a8e4c5d1461010157806391d14854146100fc578063974a74c4146100f7578063a217fddf146100f2578063a6f031bb146100ed578063c3096a69146100e8578063d547741f146100e3578063db3e1fa4146100de578063ddba6537146100d95763ef913a4b146100d4575f80fd5b610c1e565b610a23565b6109e9565b6109ba565b61093a565b61082d565b610813565b6106b3565b610672565b61063a565b61052e565b6103a7565b61034e565b610318565b6102c0565b61023b565b346101c55760203660031901126101c5576004357fffffffff0000000000000000000000000000000000000000000000000000000081168091036101c557807f7965db0b000000000000000000000000000000000000000000000000000000006020921490811561019b575b506040519015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150145f610190565b5f80fd5b9060206003198301126101c5576004356001600160401b0381116101c557826023820112156101c5578060040135926001600160401b0384116101c557602484830101116101c5576024019190565b634e487b7160e01b5f52602160045260245ffd5b6003111561023657565b610218565b346101c55760206102a261024e366101c9565b9061026160ff60045460401c1615610d09565b7fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5f525f845260405f205f8052845260ff60405f205416156102b3576120cf565b604051906102af8161022c565b8152f35b6102bb612cd0565b6120cf565b346101c55760203660031901126101c55760206102ea6004355f525f602052600160405f20015490565b604051908152f35b60409060031901126101c557600435906024356001600160a01b03811681036101c55790565b346101c55761034c610329366102f2565b90610347610342825f525f602052600160405f20015490565b612d3f565b6132ff565b005b346101c55761035c366102f2565b336001600160a01b038216036103755761034c91613397565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b5f9103126101c557565b346101c5575f3660031901126101c557335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff16156104355760045460ff8160401c161561040d5768ff00000000000000001916600455005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b90602080835192838152019201905f5b8181106104695750505090565b825163ffffffff1684526020938401939092019160010161045c565b90602080835192838152019201905f5b8181106104a25750505090565b8251845260209384019390920191600101610495565b90602080835192838152019201905f5b8181106104d55750505090565b82516001600160401b03168452602093840193909201916001016104c8565b9161051d9061050f61052b959360608652606086019061044c565b908482036020860152610485565b9160408184039101526104b8565b90565b346101c5575f3660031901126101c5576105466122b8565b506020610551613427565b910161057161056c610565835161ffff1690565b61ffff1690565b6122dc565b61058361056c610565845161ffff1690565b9261059661056c610565855161ffff1690565b60065490915f5b866105ad610565885161ffff1690565b61ffff831690811015610624579161061d60019261060f8580610608610600828f8f6105fa8f928f9261ffff9f6105e7816105f293612322565b9063ffffffff169052565b5161ffff1690565b916134ed565b929095612322565b5289612322565b906001600160401b03169052565b011661059d565b508561063686604051938493846104f4565b0390f35b346101c557610648366101c9565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b346101c557602060ff6106a7610687366102f2565b905f525f845260405f20906001600160a01b03165f5260205260405f2090565b54166040519015158152f35b346101c55760203660031901126101c5576004356001600160401b0381116101c557806004019061016060031982360301126101c5576106fb60ff60045460401c1615610d09565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610806575b61014481019061075d828461233b565b9050156107de576107bf6107ce926106369461077c604485018261236d565b61078c606487949394018361236d565b90608488013592610104890135956107a387610f7e565b60a46107c66107b66101248d018961236d565b9b909a8961233b565b3691610e68565b9a01956135ae565b6040519081529081906020820190565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b61080e612cd0565b61074d565b346101c5575f3660031901126101c55760206040515f8152f35b346101c55760203660031901126101c5576004356001600160401b0381116101c5578060040161014060031983360301126101c557610636916107ce9161087c60ff60045460401c1615610d09565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac4754161561092d575b6108db604483018261236d565b916108e9606485018261236d565b61010486013592916084870135919061090185610f7e565b61090f61012489018561236d565b97909660a46040519a61092360208d610dbd565b5f8c5201956135ae565b610935612cd0565b6108ce565b346101c55760403660031901126101c5576004356001600160401b0381116101c55761032060031982360301126101c557602435906001600160401b0382116101c557608060031983360301126101c55761034c916109ae6109a96109a560045460ff9060401c1690565b1590565b610d09565b600401906004016123a2565b346101c55761034c6109cb366102f2565b906109e4610342825f525f602052600160405f20015490565b613397565b346101c5575f3660031901126101c55760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346101c557610a95610a34366101c9565b610a4660ff60045460401c1615610d09565b7f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260ff610a8460405f205f805260205260405f2090565b541615610bda575b508101906129eb565b805190602081018051906040830190815160806060860194855197610afd83890199610ac88b516001600160801b031690565b9060405196879586957f85fbdd9c00000000000000000000000000000000000000000000000000000000875260048701612b15565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa8015610bd557610b809660c095604095610b5e945f94610ba4575b5088519051915192516001600160801b031693613cc7565b610b7260208251015160a086015190613d72565b505051015191015190613d72565b505061034c6801000000000000000068ff0000000000000000196004541617600455565b610bc791945060803d608011610bce575b610bbf8183610dbd565b810190612ade565b925f610b46565b503d610bb5565b612052565b610be390612d3f565b5f610a8c565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b90602061052b928181520190610be9565b346101c5575f3660031901126101c5576106366040516020808201526101406040820152610cfd81610c536101808201612c31565b610c6f60608301602060ff600254818116845260081c16910152565b610c9160a0830160206001600160401b03600354818116845260401c16910152565b60045463ffffffff811660e0840152610cef90602081901c63ffffffff16610100850152610cca610120850160ff8360401c1615159052565b610cde610140850160ff8360481c16611b0c565b60501c63ffffffff16610160840152565b03601f198101835282610dbd565b60405191829182610c0d565b15610d1057565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b03821117610d6757604052565b610d38565b606081019081106001600160401b03821117610d6757604052565b608081019081106001600160401b03821117610d6757604052565b60c081019081106001600160401b03821117610d6757604052565b90601f801991011681019081106001600160401b03821117610d6757604052565b60405190610dee61010083610dbd565b565b60405190610dee61026083610dbd565b60405190610dee6101c083610dbd565b60405190610dee61014083610dbd565b60405190610dee60e083610dbd565b60405190610dee608083610dbd565b60405190610dee60c083610dbd565b6001600160401b038111610d6757601f01601f191660200190565b929192610e7482610e4d565b91610e826040519384610dbd565b8294818452818301116101c5578281602093845f960137010152565b9080601f830112156101c55781602061052b93359101610e68565b60ff8116036101c557565b91908260409103126101c557604051610edc81610d4c565b60208082948035610eec81610eb9565b8452013591610efa83610eb9565b0152565b6001600160401b038116036101c557565b3590610dee82610efe565b91908260409103126101c557604051610f3281610d4c565b60208082948035610f4281610efe565b8452013591610efa83610efe565b63ffffffff8116036101c557565b3590610dee82610f50565b801515036101c557565b3590610dee82610f69565b600211156101c557565b3590610dee82610f7e565b919091610140818403126101c557610fa9610dde565b928135916001600160401b0383116101c557610fee82610fd16101209461103e968501610e9e565b8752610fe08160208501610ec4565b602088015260608301610f1a565b6040860152610fff60a08201610f5e565b606086015261101060c08201610f5e565b608086015261102160e08201610f73565b60a08601526110336101008201610f88565b60c086015201610f5e565b60e0830152565b6001600160801b038116036101c557565b3590610dee82611045565b91908260609103126101c55760405161107981610d6c565b6040808294803561108981611045565b8452602081013560208501520135910152565b8092910391606083126101c5576040516110b581610d4c565b6040819483358352601f1901126101c55760209060408051936110d785610d4c565b838101356110e481610f50565b85520135828401520152565b6001600160401b038111610d675760051b60200190565b919060c0838203126101c55760405161111f81610d87565b8093803561112c81610efe565b8252602081013561113c81610f50565b602083015261114e836040830161109c565b604083015260a0810135906001600160401b0382116101c557019180601f840112156101c557823592611180846110f0565b9361118e6040519586610dbd565b80855260208086019160051b830101918383116101c55760208101915b8383106111bd57505050505060600152565b82356001600160401b0381116101c5578201906040828703601f1901126101c557604051916111eb83610d4c565b602081013560048110156101c557835260408101356001600160401b0381116101c5576020910101906080828803126101c5576040519261122b84610d87565b82356001600160401b0381116101c55788611247918501610e9e565b8452602083013561125781611045565b6020850152604083013561126a81610f69565b60408501526060830135936001600160401b0385116101c55761129289602096879601610e9e565b6060820152838201528152019201916111ab565b91906060838203126101c557604051906112bf82610d4c565b819380356001600160401b0381116101c5578101916040838203126101c557604051926112eb84610d4c565b80356001600160401b0381116101c55781016102c0818403126101c557611310610df0565b9061131b8482610f1a565b825260408101356001600160401b0381116101c5578461133c918301610e9e565b602083015261134d60608201610f0f565b604083015261135e60808201611056565b606083015261136f60a08201610f73565b60808301526113818460c0830161109c565b60a08301526113936101208201610f73565b60c083015261014081013560e08301526113b06101608201610f73565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526113ff6102208201610f73565b6101c08301526102408101356101e083015261141e6102608201610f73565b6102008301526102808101356102208301526102a0810135906001600160401b0382116101c55761145191859101610e9e565b61024082015284526020810135926001600160401b0384116101c5576020946114808461148c96889501611107565b83820152865201610f1a565b910152565b9080601f830112156101c557610100604051926114ae8285610dbd565b839181019283116101c557905b8282106114c85750505090565b81358152602091820191016114bb565b9080601f830112156101c557604051916114f3604084610dbd565b8290604081019283116101c557905b82821061150f5750505090565b8135815260209182019101611502565b359061ffff821682036101c557565b9080601f830112156101c5578135611545816110f0565b926115536040519485610dbd565b81845260208085019260051b8201019283116101c557602001905b82821061157b5750505090565b60208091833561158a81610f50565b81520191019061156e565b9080601f830112156101c55781356115ac816110f0565b926115ba6040519485610dbd565b81845260208085019260051b8201019283116101c557602001905b8282106115e25750505090565b81358152602091820191016115d5565b9080601f830112156101c5578135611609816110f0565b926116176040519485610dbd565b81845260208085019260051b8201019283116101c557602001905b82821061163f5750505090565b60208091833561164e81610efe565b815201910190611632565b9080601f830112156101c5578135611670816110f0565b9261167e6040519485610dbd565b81845260208085019260051b8201019283116101c557602001905b8282106116a65750505090565b6020809183356116b581610f69565b815201910190611699565b919091610320818403126101c5576116d6610e00565b9281356001600160401b0381116101c557816116f3918401610f93565b84526117028160208401611061565b602085015260808201356001600160401b0381116101c557816117269184016112a6565b604085015261173760a08301611056565b60608501526117498160c08401611491565b608085015261175c816101c084016114d8565b60a085015261176f8161020084016114d8565b60c0850152611781610240830161151f565b60e08501526102608201356001600160401b0381116101c557816117a691840161152e565b6101008501526102808201356001600160401b0381116101c557816117cc91840161152e565b6101208501526102a08201356001600160401b0381116101c557816117f2918401611595565b6101408501526102c08201356001600160401b0381116101c557816118189184016115f2565b6101608501526102e08201356001600160401b0381116101c5578161183e91840161152e565b6101808501526103008201356001600160401b0381116101c5576118629201611659565b6101a0830152565b906020828203126101c55781356001600160401b0381116101c55761052b92016116c0565b81601f820112156101c5578051906118a682610e4d565b926118b46040519485610dbd565b828452602083830101116101c557815f9260208093018386015e8301015290565b91908260409103126101c5576040516118ed81610d4c565b602080829480516118fd81610eb9565b8452015191610efa83610eb9565b91908260409103126101c55760405161192381610d4c565b6020808294805161193381610efe565b8452015191610efa83610efe565b5190610dee82610f50565b5190610dee82610f69565b5190610dee82610f7e565b919091610140818403126101c557611978610dde565b928151916001600160401b0383116101c5576119bd826119a06101209461103e96850161188f565b87526119af81602085016118d5565b60208801526060830161190b565b60408601526119ce60a08201611941565b60608601526119df60c08201611941565b60808601526119f060e0820161194c565b60a0860152611a026101008201611957565b60c086015201611941565b5190610dee82611045565b91908260609103126101c557604051611a3081610d6c565b60408082948051611a4081611045565b8452602081015160208501520151910152565b6020818303126101c5578051906001600160401b0382116101c55701610180818303126101c55760405191611a8783610da2565b81516001600160401b0381116101c55782611aaa8361014093611afa9601611962565b8552611ab98360208301611a18565b6020860152611acb8360808301611a18565b6040860152611adc60e08201611a0d565b6060860152611aef83610100830161190b565b60808601520161190b565b60a082015290565b6002111561023657565b90611b1682611b02565b52565b9061052b9061012060e0611b3885516101408552610140850190610be9565b9460ff60208083015182815116828801520151166040850152611b78604082015160608601906001600160401b0360208092828151168552015116910152565b606081015163ffffffff1660a0850152608081015163ffffffff1660c085015260a0810151151584830152611bb660c0820151610100860190611b0c565b015163ffffffff16910152565b90606060c08201926001600160401b03815116835263ffffffff6020820151166020840152611c146040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519160c060a0830152825180915260e0820190602060e08260051b8501019401925f905b828210611c4857505050505090565b909192939460df198282030185528551908151600481101561023657611cc68260206001958195948295520151906040848201526060611c9483516080604085015260c0840190610be9565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610be9565b9701950193920190611c39565b9061052b90602080611e4c85516060855282611e38825160406060890152611d1460a0890182516001600160401b0360208092828151168552015116910152565b610240611d31848301516102c060e08c01526103608b0190610be9565b60408301516001600160401b03166101008b01529160608101516001600160801b03166101208b0152608081015115156101408b015260a081015180516101608c0152602090810151805163ffffffff166101808d015201516101a08b015260c081015115156101c08b015260e08101516101e08b015261010081015115156102008b01526101208101516102208b0152610140810151828b01526101608101516102608b01526101808101516102808b01526101a08101516102a08b01526101c081015115156102c08b01526101e08101516102e08b015261020081015115156103008b01526102208101516103208b01520151888203609f19016103408a0152610be9565b910151858203605f19016080870152611bc3565b9401519101906001600160401b0360208092828151168552015116910152565b905f905b60088210611e7d57505050565b6020806001928551815201930191019091611e70565b905f905b60028210611ea457505050565b6020806001928551815201930191019091611e97565b90602080835192838152019201905f5b818110611ed75750505090565b82511515845260209384019390920191600101611eca565b9061052b91602081526101a061203c61202461200c611ff4611fdc611f6a611f25895161032060208b01526103408a0190611b19565b611f5460208b015160408b0190604080916001600160801b038151168452602081015160208501520151910152565b60408a0151898203601f190160a08b0152611cd3565b60608901516001600160801b031660c0890152611f8f60808a015160e08a0190611e6c565b611fa260a08a01516101e08a0190611e93565b611fb560c08a01516102208a0190611e93565b60e089015161ffff16610260890152610100890151888203601f19016102808a015261044c565b610120880151878203601f19016102a089015261044c565b610140870151868203601f19016102c0880152610485565b610160860151858203601f19016102e08701526104b8565b610180850151848203601f190161030086015261044c565b92015190610320601f1982850301910152611eba565b6040513d5f823e3d90fd5b15612066575050565b906001600160401b0380927f7d939bd8000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b610dee909291926060810193604080916001600160801b038151168452602081015160208501520151910152565b6120db9181019061186a565b604051630647af7560e51b81525f81806120f88560048301611eef565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115610bd5575f91612296575b5061213e81612dbd565b61215061214a82612e3e565b92613088565b505061215b8261022c565b8161224b576122476122306020604060a08501946121af61218684885101516001600160401b031690565b60035460401c6001600160401b03166001600160401b0381166001600160401b0383161161205d565b61220786516001600160401b038151166fffffffffffffffffffffffffffffffff196fffffffffffffffff0000000000000000602060035494846001600160401b0319871617600355015160401b1692161717600355565b015160405161221d81610cef85820194856120a1565b519020935101516001600160401b031690565b6001600160401b03165f52600560205260405f2090565b5590565b506122558161022c565b6001810361227f5761052b6801000000000000000068ff0000000000000000196004541617600455565b6122888161022c565b6002810361052b5750600290565b6122b291503d805f833e6122aa8183610dbd565b810190611a53565b5f612134565b604051906122c582610d87565b5f6060838281528260208201528260408201520152565b906122e6826110f0565b6122f36040519182610dbd565b8281528092612304601f19916110f0565b0190602036910137565b634e487b7160e01b5f52603260045260245ffd5b80518210156123365760209160051b010190565b61230e565b903590601e19813603018212156101c557018035906001600160401b0382116101c5576020019181360383136101c557565b903590601e19813603018212156101c557018035906001600160401b0382116101c557602001918160051b360383136101c557565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac4754610dee92919060ff1661259457612400612cd0565b612594565b91906080838203126101c5576040519061241e82610d87565b819380356001600160401b0381116101c55760609261243e918301610e9e565b835260208101356020840152604081013561245881610efe565b60408401520135908160070b82036101c55760600152565b9190916080818403126101c5576040519061248a82610d87565b819381356001600160401b0381116101c557820181601f820112156101c55780356124b4816110f0565b916124c26040519384610dbd565b81835260208084019260051b820101918483116101c55760208201905b838210612530575050505083526124f860208301610f73565b60208401526040820135906001600160401b0382116101c557826125256060949261148c94869401612405565b604086015201610f0f565b81356001600160401b0381116101c55760209161255288848094880101612405565b8152019101906124df565b15612566575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b61259f9036906116c0565b9060405191630647af7560e51b83525f83806125be8460048301611eef565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa928315610bd5575f936127d8575b5061260483612dbd565b61260d83612e3e565b9061261781613088565b50506126228261022c565b600182146127b2576126549061016060406126456126403688612470565b613810565b9201515151015180821461255d565b61265d8161022c565b8061273a575060206127359160408461269d61269860a0610dee98019461269161218688885101516001600160401b031690565b3690612470565b6138ad565b6126f583516001600160401b038151166fffffffffffffffffffffffffffffffff196fffffffffffffffff0000000000000000602060035494846001600160401b0319871617600355015160401b1692161717600355565b015160405161270b81610cef86820194856120a1565b5190208151830151612725906001600160401b0316612230565b555101516001600160401b031690565b613b78565b8061274660029261022c565b1461274f575050565b6020612735916127a361269860a0610dee96019261269161277a86865101516001600160401b031690565b60035460401c6001600160401b03166001600160401b0381166001600160401b0383161461205d565b5101516001600160401b031690565b50505050610dee6801000000000000000068ff0000000000000000196004541617600455565b6127ed9193503d805f833e6122aa8183610dbd565b915f6125fa565b91906060838203126101c5576040519061280d82610d6c565b819380356001600160401b0381116101c55781016040818403126101c5576040519061283882610d4c565b8035906001600160401b0382116101c557612857856020938301610e9e565b8352013561286481610efe565b6020820152835260208101356001600160401b0381116101c5578261288a9183016112a6565b60208401526040810135916001600160401b0383116101c55760409261148c92016112a6565b919091610260818403126101c5576128c6610e10565b926128d18183611491565b84526128e18161010084016114d8565b60208501526128f48161014084016114d8565b6040850152612906610180830161151f565b60608501526101a08201356001600160401b0381116101c5578161292b91840161152e565b60808501526101c08201356001600160401b0381116101c5578161295091840161152e565b60a08501526101e08201356001600160401b0381116101c55781612975918401611595565b60c08501526102008201356001600160401b0381116101c5578161299a9184016115f2565b60e08501526102208201356001600160401b0381116101c557816129bf91840161152e565b6101008501526102408201356001600160401b0381116101c5576129e39201611659565b610120830152565b6020818303126101c5578035906001600160401b0382116101c55701610160818303126101c557612a1a610e20565b9181356001600160401b0381116101c55781612a37918401610f93565b835260208201356001600160401b0381116101c55781612a589184016127f4565b6020840152612a6a8160408401611061565b6040840152612a7c8160a08401611061565b6060840152612a8e6101008301611056565b60808401526101208201356001600160401b0381116101c55781612ab39184016128b0565b60a08401526101408201356001600160401b0381116101c557612ad692016128b0565b60c082015290565b906080828203126101c557612b0d906040805193612afb85610d4c565b612b05838261190b565b85520161190b565b602082015290565b90610dee94612bc5612b9d61010095612b3f612bea959b9a989b6101208852610120880190611b19565b86810360208801526040612b8c8351606084526001600160401b036020612b71835186606089015260a0880190610be9565b92015116608085015260208501518482036020860152611cd3565b920151906040818403910152611cd3565b986040850190604080916001600160801b038151168452602081015160208501520151910152565b80516001600160801b031660a0840152602081015160c08401526040015160e0830152565b01906001600160801b03169052565b90600182811c92168015612c27575b6020831014612c1357565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612c08565b6001545f9291612c4082612bf9565b8082529160018116908115612cb45750600114612c5b575050565b60015f9081529293509091907fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b838310612c9a575060209250010190565b600181602092949394548385870101520191019190612c89565b9050602093945060ff929192191683830152151560051b010190565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff1615612d0857565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260ff612d663360405f20906001600160a01b03165f5260205260405f2090565b541615612d705750565b63e2517d3f60e01b5f523360045260245260445ffd5b15612d8f575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b610dee90612dda81516001600160801b036060840151169061428b565b612e366001600160401b036020608081850151604051612e1a8482018093604080916001600160801b038151168452602081015160208501520151910152565b60608152612e288382610dbd565b51902094015101511661440c565b808214612d86565b612e59612230602060a084015101516001600160401b031690565b5480612e655750505f90565b60408201908151604051612e8181610cef6020820194856120a1565b5190201491821592612e9f575b505015612e9a57600190565b600290565b6001600160801b03919250612ed5612ec66020612ee19301516001600160801b0390511690565b9351516001600160801b031690565b6001600160801b031690565b911610155f80612e8e565b15612ef357565b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b15612f235750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15612f5d5750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15612f975750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b906001600160401b03809116911601906001600160401b038211612ffd57565b612fc9565b90600382029180830460031490151715612ffd57565b908160011b9180830460021490151715612ffd57565b90602c820291808304602c1490151715612ffd57565b1561304d575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b905f610100830192835151916130d860e08301936130ab610565865161ffff1690565b80911490816132ef575b816132df575b816132cf575b816132bf575b816132af575b509593949195612eec565b6130e06122b8565b506130e9613427565b90956131016007546001600160401b039060a01c1690565b945f905f985f975f96600654995b61311e6105658d5161ffff1690565b89101561324a576131406109a561313a8b6101a08e0151612322565b51151590565b61323a5761315c6131528a8751612322565b5163ffffffff1690565b809d61321c575b505060019b9589896101208201519061317b91612322565b5163ffffffff168a8d828c60208a019b828d516131999061ffff1690565b61ffff1663ffffffff821610906131af91612f55565b600163ffffffff84161b906131c78482841615612f1b565b179b516131d59061ffff1690565b906131df936134ed565b9190936101400151906131f191612322565b5114906131fd91612f8f565b61320691612fdd565b97610565600161311e925b01999791505061310f565b63ffffffff613233921663ffffffff821611612f1b565b5f8c613163565b9597610565600161311e92613211565b50985099505091965050610dee9392506132aa915061328e87876132766001600160401b038216613002565b6132886001600160401b038416613018565b10613044565b606060206040850151510151015190516101a084015191614452565b6144fc565b90506101a084015151145f6130cd565b61018085015151811491506130c7565b61016085015151811491506130c1565b61014085015151811491506130bb565b61012085015151811491506130b5565b805f525f60205260ff6133268360405f20906001600160a01b03165f5260205260405f2090565b541661339157805f525f6020526133518260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916600117905533916001600160a01b0316907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260ff6133be8360405f20906001600160a01b03165f5260205260405f2090565b54161561339157805f525f6020526133ea8260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916905533916001600160a01b0316907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b61342f6122b8565b5061347a60065461ffff6007546040519261344984610d87565b83526001600160a01b03811660208401526001600160401b038160a01c16604084015260e01c16606082015261453f565b9091565b6030019081603011612ffd57565b9060048201809211612ffd57565b90600c8201809211612ffd57565b90602c8201809211612ffd57565b6001019081600111612ffd57565b6024019081602411612ffd57565b9060018201809211612ffd57565b91908201809211612ffd57565b9091939263ffffffff61ffff9116941684101561357e576135156135108561302e565b61347e565b9363ffffffff6020868501015160e01c1603613553575061052b9061354b602061353e8661348c565b8301015160c01c9461349a565b016020015190565b7f4724a0fd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b83907ff81f5140000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b99959699989197909492986135c28b611b02565b8a15613606576136038b6135d581611b02565b7f112d89cc000000000000000000000000000000000000000000000000000000005f5260ff16600452602490565b5ffd5b88999a50602090898095969798999a15159081613803575b613627916145f0565b01966136466136358961462e565b61363f368c611061565b908761568f565b5f905f5b88868210613762575b50505061366093506147d5565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016803b156101c5578387926136cf5f956040519b8c96879586957f87d3a9b100000000000000000000000000000000000000000000000000000000875260048701614ba2565b03915afa938415610bd55761052b956136fb95613748575b506001811161370f575b5050503690611061565b6001600160801b03633b9aca009151160490565b6001600160801b036137386137266137409561462e565b9361373087614c5c565b933691614c66565b9116916157d3565b5f80806136f1565b806137565f61375c93610dbd565b8061039d565b5f6136e7565b6109a561378a61378361377d858b61379b9699979899614638565b8061236d565b369161465a565b61379536898961465a565b90615761565b6137f95750906137c26107bf6137b86137d894613660988c614638565b602081019061233b565b908151815180821491826137e3575b50506146c5565b600189935f88613653565b9091506020840120906020830120145f806137d1565b919060010161364a565b61ffff811115915061361e565b80515190811561339157613823826122dc565b915f5b8251805182101561389f579061388e613889602061384684600196612322565b5101516138846001600160401b036040613861878b51612322565b510151166040519261387284610d4c565b83526001600160401b03166020830152565b614d7b565b614e65565b6138988287612322565b5201613826565b505090505f61052b92614eb5565b6138b681613810565b90518051908115612ef35760b48211613ab0575f936138f26138e26138dd6135108661302e565b614d53565b936138ec85615977565b846159bf565b61390861390161388986615b1c565b84602e0152565b613911836159cf565b60305f955b83518710156139c2576139b78493926139b260019361399b6139398c809a612322565b519161397a63ffffffff613971604086019361396b61395f86516001600160401b031690565b6001600160401b031690565b906134e0565b9b16868d61599b565b6139946139868661348c565b91516001600160401b031690565b908b615a23565b60206139a68461349a565b91015190890160200152565b6134a8565b960195909192613916565b613a939492965060209350613a399150946139f56001600160401b03610dee976139ee82821115612eec565b16846159dd565b613a34613a0b613a058584614f4f565b946150c8565b6001600160a01b031673ffffffffffffffffffffffffffffffffffffffff196007541617600755565b600655565b613a8a613a4d82516001600160401b031690565b7fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6007549260a01b16911617600755565b015161ffff1690565b61ffff60e01b1961ffff60e01b6007549260e01b16911617600755565b7fab7bac08000000000000000000000000000000000000000000000000000000005f52600482905260b460245260445ffd5b906009548210156123365760095f52600282901c7f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af019160031b60181690565b6009549068010000000000000000821015610d6757600182016009556009548210156123365760095f5260205f208260021c019160031b6001600160401b038060c085549360031b169316831b921b1916179055565b6001600160401b0381165f5260086020526001600160a01b0380600160405f20015416613cb5600654613c04600754613bf986821691613be9613bcd6001600160401b038360a01c169261ffff9060e01c1690565b93613bd6610e2f565b9687526001600160a01b03166020870152565b6001600160401b03166040850152565b61ffff166060830152565b613c1f856001600160401b03165f52600860205260405f2090565b6060600161ffff928451815501926001600160a01b0360208201511673ffffffffffffffffffffffffffffffffffffffff1985541617845560408101517fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b8087549360a01b1616911617845501511661ffff60e01b1961ffff60e01b83549260e01b169116179055565b1615613cbe5750565b610dee90613b22565b9260208091613d36612e3695613ce8610dee996001600160401b039761428b565b604051613d158582018093604080916001600160801b038151168452602081015160208501520151910152565b60608152613d24608082610dbd565b519020612e3686858a5101511661440c565b604051613d638382018093604080916001600160801b038151168452602081015160208501520151910152565b60608152612e28608082610dbd565b9190915f92608081019384515193613dc66060840195613d97610565885161ffff1690565b8091149081613fd5575b81613fc6575b81613fb7575b81613fa7575b81613f97575b5090869391959295612eec565b6020828101510151613de0906001600160401b0316615215565b90613de96122b8565b50613df38261453f565b939096613e0a60408501516001600160401b031690565b905f925f945f613e206105655f9b5161ffff1690565b8a1015613f4d57908b949392918e613e446109a561313a8f8f906101200151612322565b613f3a576131528c613e569251612322565b8098613f1c575b50508a600197968b8060a084015190613e7591612322565b5163ffffffff16918260208a0151613e8e9061ffff1690565b61ffff1663ffffffff82161090613ea491612f55565b600163ffffffff84161b90613ebc8482841615612f1b565b1797828d8d519260200151613ed29061ffff1690565b90613edc936134ed565b91909360c0015190613eed91612322565b511490613ef991612f8f565b613f0291612fdd565b986105656001613e20925b019a929394959691908e6105f2565b63ffffffff613f33921663ffffffff821611612f1b565b5f87613e5d565b50959450986105656001613e2092613f0d565b509a509550975098915050610dee949350613f929150613f7a88886132766001600160401b038216613002565b60606020845101510151905161012085015191614452565b61533a565b905061012085015151145f613db9565b6101008601515181149150613db3565b60e08601515181149150613dad565b60c08601515181149150613da7565b60a08601515181149150613da1565b15613fed575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b5f19810191908211612ffd57565b91908203918211612ffd57565b1561403f575050565b7f12e74add000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b156140755750565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f6001546140b181612bf9565b908160448501526001811690815f1461414857506001146140e8575b506140e49192600319848303016024850152610be9565b0390fd5b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b81831061412e57509192915081016064016140e46140cd565b805460648488010152859450602090920191600101614115565b60ff191660648086019190915291151560051b840190910191506140e490506140cd565b93929190931561417c5750505050565b9160ff8092816084969581604051977fd382033a000000000000000000000000000000000000000000000000000000008952166004880152166024860152166044840152166064820152fd5b156141d1575050565b9063ffffffff80927f73b1bde3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b15614212575050565b9063ffffffff80927f1eb5c1eb000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b15614253575050565b9063ffffffff80927f87f10409000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6142a56001600160801b03610dee9316633b9aca00900490565b6142b3814242821115613fe4565b6143d160e06142c28342614029565b936143c66004546142f06142dd8263ffffffff9060501c1690565b9663ffffffff8816988942911115614036565b6143238351805160208201207f00000000000000000000000000000000000000000000000000000000000000001461406d565b61436b6020840151614336815160ff1690565b6002549160ff831660ff811660ff84161493846143df575b6143659060209060081c60ff165b93015160ff1690565b9361416c565b61439161437f606085015163ffffffff1690565b63ffffffff83811690821681146141c8565b6143b26143a5608085015163ffffffff1690565b9160201c63ffffffff1690565b63ffffffff811663ffffffff831614614209565b015163ffffffff1690565b9163ffffffff83161461424a565b9350614365602061435c6143f68286015160ff1690565b60ff600889901c8116911614969250505061434e565b6001600160401b03165f52600560205260405f2054801561442a5790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b92916144618251825114612eec565b5f5b82518110156144f5576144768183612322565b51156144ed5763ffffffff61448b8285612322565b511661449a8187518110612f55565b6144a48187612322565b5151600481101561023657600119016144c257506001905b01614463565b7f5b5c80eb000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b6001906144bc565b5050509050565b610dee9060408101519061ffff60e082015116608082015160a083015160c08401519061014085015192610160860151946101a0610180880151970151976154ce565b6145476122b8565b50602081016001600160a01b03815116156145dc576001600160a01b0390511661457861ffff60608401511661566a565b813b60018111156145c9575f198101908111612ffd5781116145b657908160016145a46145b394614d53565b9260208401903c809251614f4f565b91565b50633610565160e21b5f5260045260245ffd5b82633610565160e21b5f5260045260245ffd5b5051633a517eed60e21b5f5260045260245ffd5b156145f85750565b7fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b3561052b81610efe565b91908110156123365760051b81013590603e19813603018212156101c5570190565b929190614666816110f0565b936146746040519586610dbd565b602085838152019160051b8101918383116101c55781905b83821061469a575050505050565b81356001600160401b0381116101c5576020916146ba8784938701610e9e565b81520191019061468c565b919091156146d1575050565b906140e4614713926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610be9565b83810360031901602485015290610be9565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156101c55701602081359101916001600160401b0382116101c55781360383136101c557565b90602083828152019260208260051b82010193835f925b84841061479d5750505050505090565b9091929394956020806147c5600193601f198682030188526147bf8b88614745565b90614725565b980194019401929493919061478d565b919091156147e1575050565b6140e46040519283927ffef760c7000000000000000000000000000000000000000000000000000000008452602060048501526024840191614776565b9035601e19823603018112156101c55701602081359101916001600160401b0382116101c5578160051b360383136101c557565b9035607e19823603018112156101c5570190565b359060038210156101c557565b9035605e19823603018112156101c5570190565b90602083828152019260208260051b82010193835f925b8484106148ae5750505050505090565b909192939495602080614920600193601f198682030188526148d08b88614873565b906148da82614866565b6148e38161022c565b81526149126149076148f786850185614745565b6060888601526060850191614725565b926040810190614745565b916040818503910152614725565b980194019401929493919061489e565b61052b916149fa6149ef61497361495861494a8680614745565b608087526080870191614725565b6149656020870187614745565b908683036020880152614725565b60806149df6149856040880188614852565b868403604088015261499681614866565b61499f8161022c565b84526149ad60208201614866565b6149b68161022c565b60208501526149c760408201614866565b6149d08161022c565b60408501526060810190614745565b9190928160608201520191614725565b92606081019061481e565b916060818503910152614887565b90602083828152019260208260051b82010193835f925b848410614a2f5750505050505090565b909192939495601f198282030184528635601e19843603018112156101c557830190614a5f60208201928061481e565b8091936020845252604082019060408160051b8401019380935f915b838310614a9e575050505050506020806001929801940194019294939190614a1f565b909192939495603f19838203018652614ab78783614873565b8035614ac281610f7e565b614acb81611b02565b8252614aee614add6020830183614852565b606060208501526060840190614930565b906040810135609e19823603018112156101c5576001936020938493614b949301916040818303910152614b86614b66614b39614b2b8580614745565b60a0865260a0860191614725565b86850135614b4681610f69565b151587850152614b596040860186614852565b8482036040860152614930565b926060810135614b7581610f69565b151560608401526080810190614852565b906080818403910152614930565b980196019493019190614a7b565b9092809492969593606083019083526060602084015252608081019560808560051b83010194815f90603e1981360301995b838310614bf357505050505061052b9495506040818503910152614a08565b9091929397607f1986820301825288358b8112156101c5576020614c4d6001938683940190614c40614c36614c28848061481e565b604085526040850191614776565b9285810190614745565b9185818503910152614725565b9a019201930191909392614bd4565b3561052b81611045565b92919092614c73846110f0565b93614c816040519586610dbd565b602085828152019060051b8201918383116101c55780915b838310614ca7575050505050565b82356001600160401b0381116101c55782016040818703126101c55760405191614cd083610d4c565b81356001600160401b0381116101c557820187601f820112156101c55787816020614cfd9335910161465a565b83526020820135926001600160401b0384116101c557614d2288602095869501610e9e565b83820152815201920191614c99565b60405160809190614d428382610dbd565b6041815291601f1901366020840137565b90614d5d82610e4d565b614d6a6040519182610dbd565b8281528092612304601f1991610e4d565b602460208201906001600160401b0382511680614e18575b50614da0614dce91614d53565b92614dc5614dbf614db9614db38761589d565b876158a9565b866158c0565b856158d7565b90519084615904565b90614de361395f82516001600160401b031690565b614dec57505090565b614e0d61395f614dff614e1494866158ed565b92516001600160401b031690565b908361591c565b5090565b600191505b6080811015614e455750614da0614e3e614e39614dce936134b6565b6134c4565b9150614d93565b60019060071c910190614e1d565b805191908290602001825e015f815290565b805160018101809111612ffd57614e7b90614d53565b80511561233657614ea581614e986020945f868196015382615955565b5060405191828092614e53565b039060025afa15610bd5575f5190565b9092919280840393808511612ffd5760018514614f3e5760015b8060011b9086821015614ee25750614ecf565b91929394955050820191828111612ffd5782614efe9185614eb5565b91614f099293614eb5565b614f11614d31565b9182511561233657825f92614ea59260016020809701536021830152604182015260405191828092614e53565b5090614f4b929350612322565b5190565b919091614f5a6122b8565b506030835110613553576356414c34602084015160e01c0361355357602483015160c01c602c84015160f01c93614fe5602e82015192604e83015160f01c93614fb3614fa4610e2f565b6001600160401b039093168352565b614fc56020830198899061ffff169052565b6040820152614fdc6060820194859061ffff169052565b955161ffff1690565b9061ffff821692831593841561503e575b508315615024575b50821561500f575b50506135535750565b51915061501b9061566a565b14155f80615006565b519092506150359061ffff16610565565b1515915f614ffe565b60b41093505f614ff6565b600a907fffff00000000000000000000000000000000000000000000000000000000000061052b94937f610000000000000000000000000000000000000000000000000000000000000083521660018201527f80600a3d393df30000000000000000000000000000000000000000000000000060038201520190614e53565b604051905f60208301526150f1826150e36021820184614e53565b03601f198101845283610dbd565b6160008251116151895750610cef61514b615139615111845161ffff1690565b60f01b7fffff0000000000000000000000000000000000000000000000000000000000001690565b92604051928391602083019586615049565b51905ff0906001600160a01b0382161561516157565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b156151bc57565b633a517eed60e21b5f525f60045260245ffd5b906040516151dc81610d87565b606061ffff600183958054855201546001600160a01b03811660208501526001600160401b038160a01c16604085015260e01c16910152565b61521d6122b8565b506009548015615324576152315f9161401b565b8082106152d257506152999161527d6001600160401b0361526a61525761529495613ae2565b90546001600160401b039160031b1c1690565b92166001600160401b03831611156151b5565b6001600160401b03165f52600860205260405f2090565b6151cf565b906001600160a01b036152b660208401516001600160a01b031690565b16156152be57565b8151633a517eed60e21b5f5260045260245ffd5b906152ee6152e86152e384846134e0565b6134d2565b60011c90565b906152fb61525783613ae2565b6001600160401b03858116911611615314575090615231565b915061531f9061401b565b615231565b633a517eed60e21b5f526136036024905f600452565b610dee9161ffff6060820151168151602083015160408401519060c08501519260e086015194610120610100880151970151976154ce565b908160209103126101c5575161052b81610f69565b9060c060a061052b936001600160401b0381511684526001600160401b0360208201511660208501526040810151604085015263ffffffff6060820151166060850152608081015160808501520151918160a08201520190610be9565b9795939199989694929061ffff168852602088015f905b6008821061548957505061052b9899509261544d61547a96959361543961545c9461542e61546b986101208e0190611e93565b6101608c0190611e93565b6102406101a08b01526102408a0190610485565b908882036101c08a01526104b8565b908682036101e088015261044c565b90848203610200860152611eba565b91610220818403910152615387565b6020806001928e518152019c019101909a6153fb565b156154a657565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b9297919560209791956155f5956154e3615a6b565b5085516155278b8201519161551761395f604061550786516001600160401b031690565b935101516001600160401b031690565b6001600160401b03821614615a9b565b80516001600160401b0316966155b7604061555461554b8f86015163ffffffff1690565b63ffffffff1690565b9301518d8151910151906155a58f806155728563ffffffff90511690565b940151955151015195615595615586610e3e565b6001600160401b03909e168e52565b6001600160401b031660208d0152565b60408b015263ffffffff1660608a0152565b608088015260a08701526040519a8b998a997f4cc22bb7000000000000000000000000000000000000000000000000000000008b5260048b016153e4565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af18015610bd557610dee915f9161563b575b5061549f565b61565d915060203d602011615663575b6156558183610dbd565b810190615372565b5f615635565b503d61564b565b61ffff16602c810290808204602c1490151715612ffd5760300180603011612ffd5790565b906156de90612e3660405160208101906156c68288604080916001600160801b038151168452602081015160208501520151910152565b606081526156d5608082610dbd565b5190209161440c565b60208201518082036157335750506001600160801b03633b9aca009151160461570b814242821115613fe4565b4203428111612ffd576001600160801b03610dee911663ffffffff6004541690818110615ad8565b7f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b908151815103613391575f5b82518110156157cb576157808184612322565b515161578c8284612322565b5151036157c45761579d8184612322565b51602081519101206157af8284612322565b5160208151910120036157c45760010161576d565b5050505f90565b505050600190565b929190925f5b84518110156144f5576157ec8186612322565b51908360405160208101906001600160401b038616825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801915f5b818110615866575050505081610cef600197602061585c940151605f19848303016080850152610be9565b5190205d016157d9565b91939496509194969760208061588860019360bf198b82030188528951610be9565b970194019101918a9694939298979598615831565b6020600a910153600190565b60208260229201015360018101809111612ffd5790565b602082600a9201015360018101809111612ffd5790565b602082819201015360018101809111612ffd5790565b60208260109201015360018101809111612ffd5790565b8160209193929301015260208101809111612ffd5790565b9092919083016020015b608082101561593a57906001929391530190565b600180916080607f85161781530193019060071c9092615926565b9080519182156157cb576021602084930191015e60010180600111612ffd5790565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602d908260081c602c8201530153565b604f5f9182604e8201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b60405190615a7882610da2565b606060a0835f81525f60208201525f60408201525f838201525f60808201520152565b15615aa35750565b6001600160401b03907f90f4dbed000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15615ae1575050565b906001600160801b0380927fdb020436000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b60405190606090615b2d8284610dbd565b60228352614e1491600a906020850190601f1901368237536020602184015360028361590456fea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
}

// ContractGroth16ICS07TendermintABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractGroth16ICS07TendermintMetaData.ABI instead.
var ContractGroth16ICS07TendermintABI = ContractGroth16ICS07TendermintMetaData.ABI

// ContractGroth16ICS07TendermintBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractGroth16ICS07TendermintMetaData.Bin instead.
var ContractGroth16ICS07TendermintBin = ContractGroth16ICS07TendermintMetaData.Bin

// DeployContractGroth16ICS07Tendermint deploys a new Ethereum contract, binding an instance of ContractGroth16ICS07Tendermint to it.
func DeployContractGroth16ICS07Tendermint(auth *bind.TransactOpts, backend bind.ContractBackend, verifier common.Address, membership_ common.Address, misbehaviour_ common.Address, updateClient_ common.Address, _clientState []byte, _consensusState [32]byte, initialPinnedValidatorSet IICS07TendermintMsgsValidatorSet, roleManager common.Address) (common.Address, *types.Transaction, *ContractGroth16ICS07Tendermint, error) {
	parsed, err := ContractGroth16ICS07TendermintMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractGroth16ICS07TendermintBin), backend, verifier, membership_, misbehaviour_, updateClient_, _clientState, _consensusState, initialPinnedValidatorSet, roleManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractGroth16ICS07Tendermint{ContractGroth16ICS07TendermintCaller: ContractGroth16ICS07TendermintCaller{contract: contract}, ContractGroth16ICS07TendermintTransactor: ContractGroth16ICS07TendermintTransactor{contract: contract}, ContractGroth16ICS07TendermintFilterer: ContractGroth16ICS07TendermintFilterer{contract: contract}}, nil
}

// ContractGroth16ICS07Tendermint is an auto generated Go binding around an Ethereum contract.
type ContractGroth16ICS07Tendermint struct {
	ContractGroth16ICS07TendermintCaller     // Read-only binding to the contract
	ContractGroth16ICS07TendermintTransactor // Write-only binding to the contract
	ContractGroth16ICS07TendermintFilterer   // Log filterer for contract events
}

// ContractGroth16ICS07TendermintCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractGroth16ICS07TendermintCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractGroth16ICS07TendermintTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractGroth16ICS07TendermintTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractGroth16ICS07TendermintFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractGroth16ICS07TendermintFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractGroth16ICS07TendermintSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractGroth16ICS07TendermintSession struct {
	Contract     *ContractGroth16ICS07Tendermint // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                   // Call options to use throughout this session
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractGroth16ICS07TendermintCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractGroth16ICS07TendermintCallerSession struct {
	Contract *ContractGroth16ICS07TendermintCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                         // Call options to use throughout this session
}

// ContractGroth16ICS07TendermintTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractGroth16ICS07TendermintTransactorSession struct {
	Contract     *ContractGroth16ICS07TendermintTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                         // Transaction auth options to use throughout this session
}

// ContractGroth16ICS07TendermintRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractGroth16ICS07TendermintRaw struct {
	Contract *ContractGroth16ICS07Tendermint // Generic contract binding to access the raw methods on
}

// ContractGroth16ICS07TendermintCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractGroth16ICS07TendermintCallerRaw struct {
	Contract *ContractGroth16ICS07TendermintCaller // Generic read-only contract binding to access the raw methods on
}

// ContractGroth16ICS07TendermintTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractGroth16ICS07TendermintTransactorRaw struct {
	Contract *ContractGroth16ICS07TendermintTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractGroth16ICS07Tendermint creates a new instance of ContractGroth16ICS07Tendermint, bound to a specific deployed contract.
func NewContractGroth16ICS07Tendermint(address common.Address, backend bind.ContractBackend) (*ContractGroth16ICS07Tendermint, error) {
	contract, err := bindContractGroth16ICS07Tendermint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractGroth16ICS07Tendermint{ContractGroth16ICS07TendermintCaller: ContractGroth16ICS07TendermintCaller{contract: contract}, ContractGroth16ICS07TendermintTransactor: ContractGroth16ICS07TendermintTransactor{contract: contract}, ContractGroth16ICS07TendermintFilterer: ContractGroth16ICS07TendermintFilterer{contract: contract}}, nil
}

// NewContractGroth16ICS07TendermintCaller creates a new read-only instance of ContractGroth16ICS07Tendermint, bound to a specific deployed contract.
func NewContractGroth16ICS07TendermintCaller(address common.Address, caller bind.ContractCaller) (*ContractGroth16ICS07TendermintCaller, error) {
	contract, err := bindContractGroth16ICS07Tendermint(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractGroth16ICS07TendermintCaller{contract: contract}, nil
}

// NewContractGroth16ICS07TendermintTransactor creates a new write-only instance of ContractGroth16ICS07Tendermint, bound to a specific deployed contract.
func NewContractGroth16ICS07TendermintTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractGroth16ICS07TendermintTransactor, error) {
	contract, err := bindContractGroth16ICS07Tendermint(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractGroth16ICS07TendermintTransactor{contract: contract}, nil
}

// NewContractGroth16ICS07TendermintFilterer creates a new log filterer instance of ContractGroth16ICS07Tendermint, bound to a specific deployed contract.
func NewContractGroth16ICS07TendermintFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractGroth16ICS07TendermintFilterer, error) {
	contract, err := bindContractGroth16ICS07Tendermint(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractGroth16ICS07TendermintFilterer{contract: contract}, nil
}

// bindContractGroth16ICS07Tendermint binds a generic wrapper to an already deployed contract.
func bindContractGroth16ICS07Tendermint(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractGroth16ICS07TendermintMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractGroth16ICS07Tendermint.Contract.ContractGroth16ICS07TendermintCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.ContractGroth16ICS07TendermintTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.ContractGroth16ICS07TendermintTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractGroth16ICS07Tendermint.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.DEFAULTADMINROLE(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.DEFAULTADMINROLE(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// MISBEHAVIOURSUBMITTERROLE is a free data retrieval call binding the contract method 0xdb3e1fa4.
//
// Solidity: function MISBEHAVIOUR_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) MISBEHAVIOURSUBMITTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "MISBEHAVIOUR_SUBMITTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MISBEHAVIOURSUBMITTERROLE is a free data retrieval call binding the contract method 0xdb3e1fa4.
//
// Solidity: function MISBEHAVIOUR_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) MISBEHAVIOURSUBMITTERROLE() ([32]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.MISBEHAVIOURSUBMITTERROLE(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// MISBEHAVIOURSUBMITTERROLE is a free data retrieval call binding the contract method 0xdb3e1fa4.
//
// Solidity: function MISBEHAVIOUR_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) MISBEHAVIOURSUBMITTERROLE() ([32]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.MISBEHAVIOURSUBMITTERROLE(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// GetClientState is a free data retrieval call binding the contract method 0xef913a4b.
//
// Solidity: function getClientState() view returns(bytes)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) GetClientState(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "getClientState")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetClientState is a free data retrieval call binding the contract method 0xef913a4b.
//
// Solidity: function getClientState() view returns(bytes)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) GetClientState() ([]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GetClientState(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// GetClientState is a free data retrieval call binding the contract method 0xef913a4b.
//
// Solidity: function getClientState() view returns(bytes)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) GetClientState() ([]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GetClientState(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// GetPinnedValidatorSet is a free data retrieval call binding the contract method 0x6edfe3af.
//
// Solidity: function getPinnedValidatorSet() view returns(uint32[] indices, bytes32[] pubkeys, uint64[] votingPowers)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) GetPinnedValidatorSet(opts *bind.CallOpts) (struct {
	Indices      []uint32
	Pubkeys      [][32]byte
	VotingPowers []uint64
}, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "getPinnedValidatorSet")

	outstruct := new(struct {
		Indices      []uint32
		Pubkeys      [][32]byte
		VotingPowers []uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Indices = *abi.ConvertType(out[0], new([]uint32)).(*[]uint32)
	outstruct.Pubkeys = *abi.ConvertType(out[1], new([][32]byte)).(*[][32]byte)
	outstruct.VotingPowers = *abi.ConvertType(out[2], new([]uint64)).(*[]uint64)

	return *outstruct, err

}

// GetPinnedValidatorSet is a free data retrieval call binding the contract method 0x6edfe3af.
//
// Solidity: function getPinnedValidatorSet() view returns(uint32[] indices, bytes32[] pubkeys, uint64[] votingPowers)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) GetPinnedValidatorSet() (struct {
	Indices      []uint32
	Pubkeys      [][32]byte
	VotingPowers []uint64
}, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GetPinnedValidatorSet(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// GetPinnedValidatorSet is a free data retrieval call binding the contract method 0x6edfe3af.
//
// Solidity: function getPinnedValidatorSet() view returns(uint32[] indices, bytes32[] pubkeys, uint64[] votingPowers)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) GetPinnedValidatorSet() (struct {
	Indices      []uint32
	Pubkeys      [][32]byte
	VotingPowers []uint64
}, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GetPinnedValidatorSet(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GetRoleAdmin(&_ContractGroth16ICS07Tendermint.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GetRoleAdmin(&_ContractGroth16ICS07Tendermint.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ContractGroth16ICS07Tendermint.Contract.HasRole(&_ContractGroth16ICS07Tendermint.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ContractGroth16ICS07Tendermint.Contract.HasRole(&_ContractGroth16ICS07Tendermint.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ContractGroth16ICS07Tendermint.Contract.SupportsInterface(&_ContractGroth16ICS07Tendermint.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ContractGroth16ICS07Tendermint.Contract.SupportsInterface(&_ContractGroth16ICS07Tendermint.CallOpts, interfaceId)
}

// UpgradeClient is a free data retrieval call binding the contract method 0x8a8e4c5d.
//
// Solidity: function upgradeClient(bytes ) pure returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) UpgradeClient(opts *bind.CallOpts, arg0 []byte) error {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "upgradeClient", arg0)

	if err != nil {
		return err
	}

	return err

}

// UpgradeClient is a free data retrieval call binding the contract method 0x8a8e4c5d.
//
// Solidity: function upgradeClient(bytes ) pure returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) UpgradeClient(arg0 []byte) error {
	return _ContractGroth16ICS07Tendermint.Contract.UpgradeClient(&_ContractGroth16ICS07Tendermint.CallOpts, arg0)
}

// UpgradeClient is a free data retrieval call binding the contract method 0x8a8e4c5d.
//
// Solidity: function upgradeClient(bytes ) pure returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) UpgradeClient(arg0 []byte) error {
	return _ContractGroth16ICS07Tendermint.Contract.UpgradeClient(&_ContractGroth16ICS07Tendermint.CallOpts, arg0)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GrantRole(&_ContractGroth16ICS07Tendermint.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GrantRole(&_ContractGroth16ICS07Tendermint.TransactOpts, role, account)
}

// Misbehaviour is a paid mutator transaction binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes misbehaviourMsg) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactor) Misbehaviour(opts *bind.TransactOpts, misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.contract.Transact(opts, "misbehaviour", misbehaviourMsg)
}

// Misbehaviour is a paid mutator transaction binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes misbehaviourMsg) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) Misbehaviour(misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.Misbehaviour(&_ContractGroth16ICS07Tendermint.TransactOpts, misbehaviourMsg)
}

// Misbehaviour is a paid mutator transaction binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes misbehaviourMsg) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorSession) Misbehaviour(misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.Misbehaviour(&_ContractGroth16ICS07Tendermint.TransactOpts, misbehaviourMsg)
}

// ReAnchorPinnedSet is a paid mutator transaction binding the contract method 0xc3096a69.
//
// Solidity: function reAnchorPinnedSet(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],uint64[],uint32[],bool[]) updateMsg, ((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64) newPinnedValidatorSet) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactor) ReAnchorPinnedSet(opts *bind.TransactOpts, updateMsg IUpdateClientMsgsMsgUpdateClient, newPinnedValidatorSet IICS07TendermintMsgsValidatorSet) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.contract.Transact(opts, "reAnchorPinnedSet", updateMsg, newPinnedValidatorSet)
}

// ReAnchorPinnedSet is a paid mutator transaction binding the contract method 0xc3096a69.
//
// Solidity: function reAnchorPinnedSet(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],uint64[],uint32[],bool[]) updateMsg, ((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64) newPinnedValidatorSet) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) ReAnchorPinnedSet(updateMsg IUpdateClientMsgsMsgUpdateClient, newPinnedValidatorSet IICS07TendermintMsgsValidatorSet) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.ReAnchorPinnedSet(&_ContractGroth16ICS07Tendermint.TransactOpts, updateMsg, newPinnedValidatorSet)
}

// ReAnchorPinnedSet is a paid mutator transaction binding the contract method 0xc3096a69.
//
// Solidity: function reAnchorPinnedSet(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],uint64[],uint32[],bool[]) updateMsg, ((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64) newPinnedValidatorSet) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorSession) ReAnchorPinnedSet(updateMsg IUpdateClientMsgsMsgUpdateClient, newPinnedValidatorSet IICS07TendermintMsgsValidatorSet) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.ReAnchorPinnedSet(&_ContractGroth16ICS07Tendermint.TransactOpts, updateMsg, newPinnedValidatorSet)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.RenounceRole(&_ContractGroth16ICS07Tendermint.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.RenounceRole(&_ContractGroth16ICS07Tendermint.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.RevokeRole(&_ContractGroth16ICS07Tendermint.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.RevokeRole(&_ContractGroth16ICS07Tendermint.TransactOpts, role, account)
}

// Unfreeze is a paid mutator transaction binding the contract method 0x6a28f000.
//
// Solidity: function unfreeze() returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactor) Unfreeze(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.contract.Transact(opts, "unfreeze")
}

// Unfreeze is a paid mutator transaction binding the contract method 0x6a28f000.
//
// Solidity: function unfreeze() returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) Unfreeze() (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.Unfreeze(&_ContractGroth16ICS07Tendermint.TransactOpts)
}

// Unfreeze is a paid mutator transaction binding the contract method 0x6a28f000.
//
// Solidity: function unfreeze() returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorSession) Unfreeze() (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.Unfreeze(&_ContractGroth16ICS07Tendermint.TransactOpts)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x0bece356.
//
// Solidity: function updateClient(bytes updateClientMsg) returns(uint8)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactor) UpdateClient(opts *bind.TransactOpts, updateClientMsg []byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.contract.Transact(opts, "updateClient", updateClientMsg)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x0bece356.
//
// Solidity: function updateClient(bytes updateClientMsg) returns(uint8)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) UpdateClient(updateClientMsg []byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.UpdateClient(&_ContractGroth16ICS07Tendermint.TransactOpts, updateClientMsg)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x0bece356.
//
// Solidity: function updateClient(bytes updateClientMsg) returns(uint8)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorSession) UpdateClient(updateClientMsg []byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.UpdateClient(&_ContractGroth16ICS07Tendermint.TransactOpts, updateClientMsg)
}

// VerifyMembership is a paid mutator transaction binding the contract method 0x974a74c4.
//
// Solidity: function verifyMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[],bytes) msg_) returns(uint256)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactor) VerifyMembership(opts *bind.TransactOpts, msg_ ILightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.contract.Transact(opts, "verifyMembership", msg_)
}

// VerifyMembership is a paid mutator transaction binding the contract method 0x974a74c4.
//
// Solidity: function verifyMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[],bytes) msg_) returns(uint256)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) VerifyMembership(msg_ ILightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.VerifyMembership(&_ContractGroth16ICS07Tendermint.TransactOpts, msg_)
}

// VerifyMembership is a paid mutator transaction binding the contract method 0x974a74c4.
//
// Solidity: function verifyMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[],bytes) msg_) returns(uint256)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorSession) VerifyMembership(msg_ ILightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.VerifyMembership(&_ContractGroth16ICS07Tendermint.TransactOpts, msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0xa6f031bb.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[]) msg_) returns(uint256)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactor) VerifyNonMembership(opts *bind.TransactOpts, msg_ ILightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.contract.Transact(opts, "verifyNonMembership", msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0xa6f031bb.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[]) msg_) returns(uint256)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) VerifyNonMembership(msg_ ILightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.VerifyNonMembership(&_ContractGroth16ICS07Tendermint.TransactOpts, msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0xa6f031bb.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[]) msg_) returns(uint256)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorSession) VerifyNonMembership(msg_ ILightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.VerifyNonMembership(&_ContractGroth16ICS07Tendermint.TransactOpts, msg_)
}

// ContractGroth16ICS07TendermintRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintRoleAdminChangedIterator struct {
	Event *ContractGroth16ICS07TendermintRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractGroth16ICS07TendermintRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractGroth16ICS07TendermintRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractGroth16ICS07TendermintRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractGroth16ICS07TendermintRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractGroth16ICS07TendermintRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractGroth16ICS07TendermintRoleAdminChanged represents a RoleAdminChanged event raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ContractGroth16ICS07TendermintRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ContractGroth16ICS07TendermintRoleAdminChangedIterator{contract: _ContractGroth16ICS07Tendermint.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ContractGroth16ICS07TendermintRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractGroth16ICS07TendermintRoleAdminChanged)
				if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) ParseRoleAdminChanged(log types.Log) (*ContractGroth16ICS07TendermintRoleAdminChanged, error) {
	event := new(ContractGroth16ICS07TendermintRoleAdminChanged)
	if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractGroth16ICS07TendermintRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintRoleGrantedIterator struct {
	Event *ContractGroth16ICS07TendermintRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractGroth16ICS07TendermintRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractGroth16ICS07TendermintRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractGroth16ICS07TendermintRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractGroth16ICS07TendermintRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractGroth16ICS07TendermintRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractGroth16ICS07TendermintRoleGranted represents a RoleGranted event raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ContractGroth16ICS07TendermintRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ContractGroth16ICS07TendermintRoleGrantedIterator{contract: _ContractGroth16ICS07Tendermint.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ContractGroth16ICS07TendermintRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractGroth16ICS07TendermintRoleGranted)
				if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) ParseRoleGranted(log types.Log) (*ContractGroth16ICS07TendermintRoleGranted, error) {
	event := new(ContractGroth16ICS07TendermintRoleGranted)
	if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractGroth16ICS07TendermintRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintRoleRevokedIterator struct {
	Event *ContractGroth16ICS07TendermintRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractGroth16ICS07TendermintRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractGroth16ICS07TendermintRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractGroth16ICS07TendermintRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractGroth16ICS07TendermintRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractGroth16ICS07TendermintRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractGroth16ICS07TendermintRoleRevoked represents a RoleRevoked event raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ContractGroth16ICS07TendermintRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ContractGroth16ICS07TendermintRoleRevokedIterator{contract: _ContractGroth16ICS07Tendermint.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ContractGroth16ICS07TendermintRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractGroth16ICS07TendermintRoleRevoked)
				if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) ParseRoleRevoked(log types.Log) (*ContractGroth16ICS07TendermintRoleRevoked, error) {
	event := new(ContractGroth16ICS07TendermintRoleRevoked)
	if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
