// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractSpectreClient

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

// ICS02ClientMsgsHeight is an auto generated low-level Go binding around an user-defined struct.
type ICS02ClientMsgsHeight struct {
	RevisionNumber uint64
	RevisionHeight uint64
}

// LightClientMsgsMsgVerifyMembership is an auto generated low-level Go binding around an user-defined struct.
type LightClientMsgsMsgVerifyMembership struct {
	Height                ICS02ClientMsgsHeight
	KvPairs               []MembershipMsgsKVPair
	MerkleProofs          []MembershipMsgsMerkleProof
	AppHash               [32]byte
	TrustedConsensusState SpectreMsgsConsensusState
	MembershipType        uint8
	Path                  [][]byte
	Value                 []byte
}

// LightClientMsgsMsgVerifyNonMembership is an auto generated low-level Go binding around an user-defined struct.
type LightClientMsgsMsgVerifyNonMembership struct {
	Height                ICS02ClientMsgsHeight
	KvPairs               []MembershipMsgsKVPair
	MerkleProofs          []MembershipMsgsMerkleProof
	AppHash               [32]byte
	TrustedConsensusState SpectreMsgsConsensusState
	MembershipType        uint8
	Path                  [][]byte
}

// MembershipMsgsCommitmentProof is an auto generated low-level Go binding around an user-defined struct.
type MembershipMsgsCommitmentProof struct {
	ProofType         uint8
	ExistenceProof    MembershipMsgsExistenceProof
	NonExistenceProof MembershipMsgsNonExistenceProof
}

// MembershipMsgsExistenceProof is an auto generated low-level Go binding around an user-defined struct.
type MembershipMsgsExistenceProof struct {
	Key   []byte
	Value []byte
	Leaf  MembershipMsgsLeafOp
	Path  []MembershipMsgsInnerOp
}

// MembershipMsgsInnerOp is an auto generated low-level Go binding around an user-defined struct.
type MembershipMsgsInnerOp struct {
	HashOp uint8
	Prefix []byte
	Suffix []byte
}

// MembershipMsgsKVPair is an auto generated low-level Go binding around an user-defined struct.
type MembershipMsgsKVPair struct {
	Path  [][]byte
	Value []byte
}

// MembershipMsgsLeafOp is an auto generated low-level Go binding around an user-defined struct.
type MembershipMsgsLeafOp struct {
	HashOp       uint8
	PrehashKey   uint8
	PrehashValue uint8
	Prefix       []byte
}

// MembershipMsgsMerkleProof is an auto generated low-level Go binding around an user-defined struct.
type MembershipMsgsMerkleProof struct {
	Proofs []MembershipMsgsCommitmentProof
}

// MembershipMsgsNonExistenceProof is an auto generated low-level Go binding around an user-defined struct.
type MembershipMsgsNonExistenceProof struct {
	Key      []byte
	HasLeft  bool
	Left     MembershipMsgsExistenceProof
	HasRight bool
	Right    MembershipMsgsExistenceProof
}

// SpectreMsgsConsensusState is an auto generated low-level Go binding around an user-defined struct.
type SpectreMsgsConsensusState struct {
	Timestamp          *big.Int
	Root               [32]byte
	NextValidatorsHash [32]byte
}

// SpectreMsgsValidatorInfo is an auto generated low-level Go binding around an user-defined struct.
type SpectreMsgsValidatorInfo struct {
	ValAddress       []byte
	PubKey           [32]byte
	VotingPower      uint64
	ProposerPriority int64
}

// SpectreMsgsValidatorSet is an auto generated low-level Go binding around an user-defined struct.
type SpectreMsgsValidatorSet struct {
	Validators       []SpectreMsgsValidatorInfo
	HasProposer      bool
	Proposer         SpectreMsgsValidatorInfo
	TotalVotingPower uint64
}

// ContractSpectreClientMetaData contains all meta data concerning the ContractSpectreClient contract.
var ContractSpectreClientMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"updateClientModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membershipModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviourModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"clientState_\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"consensusState_\",\"type\":\"tuple\",\"internalType\":\"structSpectreMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"initialPinnedValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structSpectreMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structSpectreMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structSpectreMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPinnedValidatorSet\",\"inputs\":[],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPinnedValidatorsHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateApplicationState\",\"inputs\":[{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumLightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateConsensusState\",\"inputs\":[{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumLightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structLightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structSpectreMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structLightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structSpectreMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ClientFrozen\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClientUnfrozen\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClientUpdated\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConsensusStateUpdated\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientNotFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DirectCallNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"DuplicateValidatorPubkey\",\"inputs\":[{\"name\":\"firstIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"secondIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GenesisPinnedValidatorSetMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NonMonotonicHeightUpdate\",\"inputs\":[{\"name\":\"latestHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"updateHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"PinnedIndexOverflowsBitmask\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"PinnedValidatorSetStale\",\"inputs\":[{\"name\":\"expectedNextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actualPinnedValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofSignerCommitSigMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ProofSignerValidatorMismatch\",\"inputs\":[{\"name\":\"signerIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorCountExceedsLimit\",\"inputs\":[{\"name\":\"validatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxValidatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroVotingPower\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	Bin: "0x6101008060405234610f2d57616729803803809161001d8285610f94565b833981018181036101208112610f2d5761003683610fb7565b9161004360208501610fb7565b9261005060408601610fb7565b60608601519093906001600160401b038111610f2d5761007484606092890161101c565b91607f190112610f2d5760405193606085016001600160401b03811186821017610cb55760405260808701516001600160801b0381168103610f2d57855260a0870151906020860191825260c0880151946040870195865260e089015160018060401b038111610f2d578901608081830312610f2d57604051996100f78b610f5e565b81516001600160401b038111610f2d57820183601f82011215610f2d57805161011f81611039565b9161012d6040519384610f94565b81835260208084019260051b82010191868311610f2d5760208201905b838210610f3157505050508b52610163602083016110cf565b60208c015260408201516001600160401b038111610f2d576060838d60406101966101aa9861019f966101009901611064565b91015201611050565b60608c015201610fb7565b967f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a6123360e05283518401946020860194602081880312610f2d576020810151906001600160401b038211610f2d570195869003601f198101949061012013610f2d576040519560e087016001600160401b03811188821017610cb55760405260208801516001600160401b038111610f2d576020908901019080601f83011215610f2d57815161025b92602001610fe6565b865260408512610f2d57604080519561027387610f79565b61027e828a016110dc565b875261028c60608a016110dc565b602088015260208801968752603f190112610f2d57604051976102ae89610f79565b6102ba60808901611050565b89526102c860a08901611050565b60208a0152604087019889526102e060c089016110ea565b96606081019788526102f460e08a016110ea565b9860808201998a5261031d61012061030f61010084016110cf565b9260a08501938452016110ea565b60c0830190815282518051919991906001600160401b038211610cb5575f5160206166a95f395f51905f5254600181811c91168015610f23575b6020821014610f0f57601f8111610ea0575b50602090601f8311600114610e1e576104b7949392915f9183610e13575b50508160011b915f199060031b1c1916175f5160206166a95f395f51905f52555b5160ff81511661ff0060205f5160206165c95f395f51905f525493015160081b169161ffff191617175f5160206165c95f395f51905f52558b5160018060401b038151165f5160206167095f395f51905f525491602068010000000000000000600160801b0391015160401b169160018060801b03191617175f5160206167095f395f51905f52558a63ffffffff8b51169168ff000000000000000067ffffffff000000005f5160206166095f395f51905f5254935160201b169151151560401b916cffffffff0000000000000000008d5160481b1693856cffffffff000000000000000000199160018060481b031916171617911617175f5160206166095f395f51905f52558015156110fb565b6104ca63ffffffff8851168015156110fb565b516001600160401b03906020906104e090611168565b01518a51516001600160401b03169116818103610dfe5750506105028c61135c565b8151808203610de957505060405190602082019260018060801b038c511684525160408301525160608201526060815261053d608082610f94565b51902060018060401b03602089510151165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f205560805260a05260c05263ffffffff80835116915116019063ffffffff82116107895763ffffffff80809451169384925116921611610dd45750506105bc8461135c565b9351908151938415610dad5760b48511610dbc575f93602c8602868104602c036107895760300180603011610789576105f4906117e7565b9260208401976056895361564160218601536256414c60228601536356414c346023860153602c8501978060081c8953602d86015361063a61063582611990565b611819565b97602e8601988952604e8601915f83535f604f8801535f9a60305b89518d10156107f1578c610669818c611348565b51604081018051909d91906001600160401b0316156107de57602001908d5f8e5b85821061079d575050516106a9916001600160401b039091169061115b565b9c63ffffffff848d019360ff8160181c16602086015361ffff8160101c16602186015362ffffff8160081c166022860153166023840153600484018411610789575160ff8160381c16602484015361ffff8160301c16602584015362ffffff8160281c16602684015363ffffffff8160201c16602784015364ffffffffff8160181c16602884015365ffffffffffff8160101c16602984015366ffffffffffffff8160081c16602a8401536001600160401b0316602b830153600c8301831161078957602c9051910152602c8101809111610789576001909c019b610655565b634e487b7160e01b5f52601160045260245ffd5b8160209293506107b09196949596611348565b5101518451146107c9576001018e908e9493929461068a565b63577d2b6160e11b5f5260045260245260445ffd5b8263afef9f1560e01b5f5260045260245ffd5b5087908b8b6001600160401b038111610dad57602484019060ff8160381c16825361ffff8160301c16602586015362ffffff8160281c16602686015363ffffffff8160201c16602786015364ffffffffff8160181c16602886015365ffffffffffff8160101c16602986015366ffffffffffffff8160081c16602a8601536001600160401b0316602b8501535f606060405161088c81610f5e565b82815282602082015282604082015201526030845110610d9a576356414c34835160e01c03610d9a575160c01c945160f01c9051955160f01c95604051956108d387610f5e565b865260208601968288526040870191825280606088015282159283918415610d8f575b8415610d85575b508315610d5a575b50508115610d43575b50610d30576040519161094360218460208101945f8652845180918484015e81015f838201520301601f198101855284610f94565b616000835111610d1d57506109a4602a61ffff60f01b845160f01b16936040519384916020830196606160f81b885260218401526680600a3d393df360c81b60238401525180918484015e81015f838201520301601f198101835282610f94565b51905ff06001600160a01b03168015610d0e575f5160206166495f395f51905f5280545f5160206166895f395f51905f52939093559251935161ffff60e01b60e09190911b16600160a01b600160e01b0360a09590951b949094166001600160f01b0319909216600160a01b600160f01b0319919091161717919091179055516020015190515f5160206165a95f395f51905f52546001600160401b03909216916001600160801b039091169080610cc9575b506001600160401b0382165f9081525f5160206166c95f395f51905f52602052604090206001600160a01b0390600101545f5160206166895f395f51905f52545f5160206166495f395f51905f52546040519390921615939260a08101916001600160401b03831182841017610cb55760029260405281526020810160018060a01b0384168152604082019360018060401b038160a01c16855261ffff606084019160e01c1681526080830194868652610b2d8960018060401b03165f525f5160206166c95f395f51905f5260205260405f2090565b9351845591516001840180549351925161ffff60e01b60e09190911b16600160a01b600160e01b0360a09490941b939093166001600160a01b039092166001600160f01b031990941693909317171790559151910180546001600160801b039092166001600160801b03199283161790557f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17098054909116919091179055610c55575b506001600160a01b038116610c2f5750610be7611697565b50610bf360e051611719565b505b604051614bbe90816119cb823960805181610e9c015260a05181613967015260c051816105ee015260e05181818161032801526106d00152f35b80610c3c610c4f926115a1565b50610c4681611617565b5060e051611772565b50610bf5565b5f5160206165a95f395f51905f525468010000000000000000811015610cb557806001610c9192015f5160206165a95f395f51905f5255611549565b81546001600160401b0360039290921b91821b191692901b91909117905581610bcf565b634e487b7160e01b5f52604160045260245ffd5b5f19810190811161078957610cdd90611549565b905460039190911b1c6001600160401b031680831015610a57579050630fb2737b60e31b5f5260045260245260445ffd5b63fbad885d60e01b5f5260045ffd5b516331734c2160e21b5f5260045260245ffd5b82634724a0fd60e01b5f5260045260245ffd5b905051610d5261063585611990565b14158961090e565b8551929350602c80820292918304141715610789576030019081603011610789571415908a80610905565b151593508c6108fd565b60b4821194506108f6565b84634724a0fd60e01b5f5260045260245ffd5b6305f8ded760e21b5f5260045ffd5b8463156f758160e31b5f5260045260b460245260445ffd5b6333fae18560e21b5f5260045260245260445ffd5b631d9feb7d60e11b5f5260045260245260445ffd5b6355bace6f60e01b5f5260045260245260445ffd5b015190505f80610387565b90601f198316915f5160206166a95f395f51905f525f52815f20925f5b818110610e8857509160019391856104b79897969410610e70575b505050811b015f5160206166a95f395f51905f52556103a8565b01515f1960f88460031b161c191690555f8080610e56565b92936020600181928786015181550195019301610e3b565b5f5160206166a95f395f51905f525f527fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be601f840160051c81019160208510610f05575b601f0160051c01905b818110610efa5750610369565b5f8155600101610eed565b9091508190610ee4565b634e487b7160e01b5f52602260045260245ffd5b90607f1690610357565b5f80fd5b81516001600160401b038111610f2d57602091610f538a848094880101611064565b81520191019061014a565b608081019081106001600160401b03821117610cb557604052565b604081019081106001600160401b03821117610cb557604052565b601f909101601f19168101906001600160401b03821190821017610cb557604052565b51906001600160a01b0382168203610f2d57565b6001600160401b038111610cb557601f01601f191660200190565b929192610ff282610fcb565b916110006040519384610f94565b829481845281830111610f2d578281602093845f96015e010152565b9080601f83011215610f2d57815161103692602001610fe6565b90565b6001600160401b038111610cb55760051b60200190565b51906001600160401b0382168203610f2d57565b9190608083820312610f2d576040519061107d82610f5e565b8351919384926001600160401b038111610f2d5760609261109f91830161101c565b8352602081015160208401526110b760408201611050565b60408401520151908160070b8203610f2d5760600152565b51908115158203610f2d57565b519060ff82168203610f2d57565b519063ffffffff82168203610f2d57565b156111035750565b63ffffffff9063b0369c3160e01b5f5216600452600160245263ffffffff60445260645ffd5b9190820391821161078957565b908151811015611147570160200190565b634e487b7160e01b5f52603260045260245ffd5b9190820180921161078957565b60405161117481610f79565b606081525f6020820152508051801590811561133c575b5061132d575f19908051805b6112e9575b505f1982146112db57600182019081831161078957600360fc1b6001600160f81b03196111c98484611136565b511614806112c7575b6112b7575f5b8151831015611267576111eb8383611136565b5160f81c60308110801561125d575b61124b57600a82026001600160401b03908116602f1990920160ff169190910181169116811061122f576001909201916111d8565b509150506040519061124082610f79565b81525f602082015290565b50509150506040519061124082610f79565b50603981116111fa565b91509160018111908115916112ab575b5061129c576040519161128983610f79565b82526001600160401b0316602082015290565b63384c687760e21b5f5260045ffd5b602b915010155f611277565b9150506040519061124082610f79565b5060026112d5848351611129565b116111d2565b604051915061124082610f79565b5f19810181811161078957602d60f81b6001600160f81b031961130c8386611136565b51161461132357508015610789575f190180611197565b92505f905061119c565b6329120bff60e21b5f5260045ffd5b6040915010155f61118b565b80518210156111475760209160051b010190565b9081515180156115435761136f81611039565b9061137d6040519283610f94565b808252601f1961138c82611039565b013660208401375f5b84518051821015611533576113ac82602092611348565b51015160018060401b0360406113c3848951611348565b51015116604051916113d483610f79565b82526020820190808252602490806114d1575b506113f1906117e7565b91600a6020840153600190818401906022602083015382800191828411610789576021600a910153600283018092116107895761143b916114329086611938565b9051908561194e565b82519092906001600160401b031661146e575b5050509061145d600192611819565b6114678286611348565b5201611395565b828401906010602083015382840180941161078957516001600160401b03169260219190910191905b60808410156114b45750509161145d91600194935391925f61144e565b91908084926080607f8397161781530192019060071c9290611497565b6001905b608081101561152557508060010160011161151157602501908181116114fe57506113f16113e7565b634e487b7160e01b5f9081526011600452fd5b50634e487b7160e01b5f9081526011600452fd5b60019060071c9101906114d5565b5050906110369293505f90611877565b505f9150565b905f5160206165a95f395f51905f5254821015611147575f5160206165a95f395f51905f525f52600282901c7f06ff75c3bdc3f03474569b3e952aab1bfb9d41b303b5638523d3e1200dcd1ae2019160031b60181690565b6001600160a01b0381165f9081525f5160206166e95f395f51905f52602052604090205460ff16611612576001600160a01b03165f8181525f5160206166e95f395f51905f5260205260408120805460ff191660011790553391905f5160206165895f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f5160206165e95f395f51905f52602052604090205460ff16611612576001600160a01b03165f8181525f5160206165e95f395f51905f5260205260408120805460ff191660011790553391905f5160206166695f395f51905f52905f5160206165895f395f51905f529080a4600190565b5f80525f5160206165e95f395f51905f526020525f5160206166295f395f51905f525460ff16611715575f8080525f5160206165e95f395f51905f526020525f5160206166295f395f51905f52805460ff1916600117905533905f5160206166695f395f51905f525f5160206165895f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f1461161257805f525f60205260405f205f805260205260405f20600160ff198254161790555f33915f5160206165895f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff166117e1575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291905f5160206165895f395f51905f529080a4600190565b50505f90565b906117f182610fcb565b6117fe6040519182610f94565b828152809261180f601f1991610fcb565b0190602036910137565b8051600181018091116107895761182f906117e7565b805115611147576020918161184b845f94019284845382611966565b50604051918291518091835e8101838152039060025afa1561186c575f5190565b6040513d5f823e3d90fd5b9291926118848285611129565b93600185146119285760015b8060011b90868210156118a35750611890565b9394955050906118bd6118b6848661115b565b8583611877565b6118d1936118cb919561115b565b90611877565b906040516118e0608082610f94565b6041815260208101906060368337805115611147576020935f936001845360218301526041820152604051918291518091835e8101838152039060025afa1561186c575f5190565b50611934929350611348565b5190565b6020828192010153600181018091116107895790565b81602091939293010152602081018091116107895790565b908051918215611988576021602084930191015e600101806001116107895790565b505050600190565b6119c6604051916119a2606084610f94565b602283526040366020850137600a60208401536119c0600184611938565b8361194e565b509056fe6080806040526004361015610012575f80fd5b5f3560e01c90816301ffc9a71461116957508063248a9ca31461113f5780632f2ff15d1461111057806336568abe146110c15780634b1872be14610d7e5780636a28f00014610c9d5780636edfe3af14610b0c5780638a8e4c5d14610ad45780638c80fbda14610a9857806391d1485414610a5c57806392c19bdc146109b6578063974a74c41461085f578063a217fddf14610845578063a6f031bb14610729578063d547741f146106f3578063db3e1fa4146106b9578063ddba6537146102ff5763ef913a4b146100e2575f80fd5b346102fb575f3660031901126102fb5760405160208082015261012060408201525f7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1700548060011c906001811680156102f1575b6020831081146102dd578261016086015290815f146102b85750600114610238575b61023483610220818560ff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170154818116606085015260081c1660808301526001600160401b037f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025481811660a085015260401c1660c083015263ffffffff5f516020614b925f395f51905f525481811660e0850152818160201c1661010085015260ff8160401c16151561012085015260481c1661014083015203601f19810183528261135d565b6040519182916020835260208301906112cd565b0390f35b9190507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005f527fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be915f905b80821061029c5750909150810161018001610220610158565b9192600181602092546101808588010152019101909291610283565b60ff19166101808086019190915291151560051b840190910191506102209050610158565b634e487b7160e01b5f52602260045260245ffd5b91607f1691610136565b5f80fd5b346102fb5761030d3661122d565b60ff5f516020614b925f395f51905f525460401c16610691577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610682575b508101906020818303126102fb578035906001600160401b0382116102fb5701610140818303126102fb5760405160c081018181106001600160401b0382111761066e5760405281356001600160401b0381116102fb5782016040818503126102fb57604051906103d6826112f1565b80356001600160401b0381116102fb57856103f291830161150b565b82526020810135906001600160401b0382116102fb576104149186910161150b565b602082015281526104288360208401611392565b916020820192835261043d8460808301611392565b6040830190815261045060e0830161137e565b92606081019384526101008301356001600160401b0381116102fb57866104789185016118ba565b9560808201968752610120840135956001600160401b0387116102fb57610612976001600160801b036105ec976105af6105de976105886105616104c66105ca9960209f61054b9e016118ba565b9960a081019a8b526104da878751166128a8565b6104e8885182515190613e4e565b6104f98b5160208351015190613e4e565b6040519e8f9d8e927fb9ffc5010000000000000000000000000000000000000000000000000000000083850152826024850152519261014060448201526101c484519160406101848201520190611b24565b9101518d820361018319016101a48f0152611b24565b965180516001600160801b031660648d0152602081015160848d01526040015160a48c0152565b5180516001600160801b031660c48b0152602081015160e48b0152604001516101048a0152565b51166101248701525185820360431901610144870152611dc1565b905183820360431901610164850152611dc1565b03601f19810183528261135d565b7f0000000000000000000000000000000000000000000000000000000000000000612933565b506801000000000000000068ff0000000000000000195f516020614b925f395f51905f525416175f516020614b925f395f51905f52557f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1005b634e487b7160e01b5f52604160045260245ffd5b61068b906126f3565b82610366565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102fb575f3660031901126102fb5760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346102fb5761072761070436611207565b9061072261071d825f525f602052600160405f20015490565b6126f3565b6127c0565b005b346102fb5760203660031901126102fb576004356001600160401b0381116102fb57806004019061014060031982360301126102fb5760ff5f516020614b925f395f51905f525460401c16610691575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff1615610838575b6107d7604482018361264f565b916107e5606482018561264f565b9390916101048101359060018210156102fb576020966108309661080d61012484018361264f565b9690956040519861081e8c8b61135d565b5f8a52608460a4870196013594613817565b604051908152f35b610840612684565b6107ca565b346102fb575f3660031901126102fb5760206040515f8152f35b346102fb5760203660031901126102fb576004356001600160401b0381116102fb57806004019061016060031982360301126102fb5760ff5f516020614b925f395f51905f525460401c16610691575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff16156109a9575b610144810161090f818461261d565b90501561098157610923604483018461264f565b9091610932606485018661264f565b9590946101048101359160018310156102fb576020976108309761096a9661097161096161012487018661264f565b9990988661261d565b3691611430565b98608460a4870196013594613817565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b6109b1612684565b610900565b346102fb576109c43661122d565b9060ff5f516020614b925f395f51905f525460401c16610691575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060209081527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47549092610a3e92909160ff1615610a4f5761210c565b60405190610a4b8161127c565b8152f35b610a57612684565b61210c565b346102fb57610a6a36611207565b905f525f6020526001600160a01b0360405f2091165f52602052602060ff60405f2054166040519015158152f35b346102fb575f3660031901126102fb5760207f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170554604051908152f35b346102fb57610ae23661122d565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102fb575f3660031901126102fb576020610b26612a75565b610b2f81612b56565b92019161ffff83511691610b5b610b45846114f4565b93610b53604051958661135d565b8085526114f4565b602084019190601f1901368337610b7661ffff865116612049565b9261ffff86511691610b8a610b45846114f4565b602084019690601f19013688375f5b61ffff89511661ffff821690811015610bfe5761ffff91886001600160401b03610bf384808f610be482610bec928e8c8f8f9d610bd960019f829061207b565b525192511691612c07565b92909661207b565b528a61207b565b911690520116610b99565b87838a888a604051948594606086019060608752518091526080860192905f5b818110610c7e5750505081610c3b9186602094038488015261129a565b91848303604086015251918281520191905f5b818110610c5c575050500390f35b82516001600160401b0316845285945060209384019390920191600101610c4e565b825163ffffffff16855288975060209485019490920191600101610c1e565b346102fb575f3660031901126102fb57335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff1615610d67575f516020614b925f395f51905f525460ff8160401c1615610d3f5768ff000000000000000019165f516020614b925f395f51905f52557f8e0da210a2ffecc744373b4860cd9b8dc471810f5fb61b7d2c3234ccd4aae30d5f80a1005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b346102fb57610d8c3661122d565b60ff5f516020614b925f395f51905f525460401c16610691575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff16156110b4575b8101906020818303126102fb578035906001600160401b0382116102fb576105de610e9a610e2f602095610ece95610ec09501611aa7565b610e456001600160801b036040830151166128a8565b610e7f60608201518783015190606089610e5d612a75565b93610e746001600160801b03608087015116614035565b510151015190614159565b6040519283916343c3b39360e11b8884015260248301611ea7565b7f0000000000000000000000000000000000000000000000000000000000000000612933565b828082518301019101611fa7565b610ed781612981565b90610ee18261127c565b8161104457826001600160401b0381838160607f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb3196019184828451015116610f537f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025491878360401c16808211612005565b83516fffffffffffffffffffffffffffffffff196fffffffffffffffff0000000000000000858984511693015160401b16921617177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1702550151604051610fd98482018093604080916001600160801b038151168452602081015160208501520151910152565b60608152610fe860808261135d565b51902061102b848484510151166001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b5551015116604051908152a160405190610a4b8161127c565b5061104e8161127c565b60018103610a3e576801000000000000000068ff0000000000000000195f516020614b925f395f51905f525416175f516020614b925f395f51905f52557f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1610a3e565b6110bc612684565b610df7565b346102fb576110cf36611207565b336001600160a01b038216036110e857610727916127c0565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102fb5761072761112136611207565b9061113a61071d825f525f602052600160405f20015490565b612733565b346102fb5760203660031901126102fb5760206108306004355f525f602052600160405f20015490565b346102fb5760203660031901126102fb57600435907fffffffff0000000000000000000000000000000000000000000000000000000082168092036102fb57817f7965db0b00000000000000000000000000000000000000000000000000000000602093149081156111dd575b5015158152f35b7f01ffc9a700000000000000000000000000000000000000000000000000000000915014836111d6565b60409060031901126102fb57600435906024356001600160a01b03811681036102fb5790565b9060206003198301126102fb576004356001600160401b0381116102fb57826023820112156102fb578060040135926001600160401b0384116102fb57602484830101116102fb576024019190565b6003111561128657565b634e487b7160e01b5f52602160045260245ffd5b90602080835192838152019201905f5b8181106112b75750505090565b82518452602093840193909201916001016112aa565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b604081019081106001600160401b0382111761066e57604052565b606081019081106001600160401b0382111761066e57604052565b608081019081106001600160401b0382111761066e57604052565b60a081019081106001600160401b0382111761066e57604052565b90601f801991011681019081106001600160401b0382111761066e57604052565b35906001600160801b03821682036102fb57565b91908260609103126102fb576040516113aa8161130c565b60408082946113b88161137e565b8452602081013560208501520135910152565b35906001600160401b03821682036102fb57565b91908260409103126102fb576040516113f7816112f1565b6020611410818395611408816113cb565b8552016113cb565b910152565b6001600160401b03811161066e57601f01601f191660200190565b92919261143c82611415565b9161144a604051938461135d565b8294818452818301116102fb578281602093845f960137010152565b359081151582036102fb57565b359063ffffffff821682036102fb57565b8092910391606083126102fb5760405161149d816112f1565b6040819483358352601f1901126102fb5760209060408051936114bf856112f1565b6114ca848201611473565b85520135828401520152565b9080601f830112156102fb578160206114f193359101611430565b90565b6001600160401b03811161066e5760051b60200190565b91906060838203126102fb5760405190611524826112f1565b819380356001600160401b0381116102fb5781016040818403126102fb576040519061154f826112f1565b80356001600160401b0381116102fb5781016102c0818603126102fb576040519061026082018281106001600160401b0382111761066e5760405261159486826113df565b825260408101356001600160401b0381116102fb57810186601f820112156102fb57868160206115c693359101611430565b60208301526115d7606082016113cb565b60408301526115e86080820161137e565b60608301526115f960a08201611466565b608083015261160b8660c08301611484565b60a083015261161d6101208201611466565b60c083015261014081013560e083015261163a6101608201611466565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526116896102208201611466565b6101c08301526102408101356101e08301526116a86102608201611466565b6102008301526102808101356102208301526102a0810135906001600160401b0382116102fb576116db918791016114d6565b61024082015282526020810135906001600160401b0382116102fb570160c0818503126102fb576040519061170f82611327565b611718816113cb565b825261172660208201611473565b60208301526117388560408301611484565b604083015260a0810135906001600160401b0382116102fb570184601f820112156102fb578035611768816114f4565b91611776604051938461135d565b81835260208084019260061b820101908782116102fb57602001915b8183106117b85750505092602094928261141095606088950152838201528652016113df565b6040838903126102fb57604051906117cf826112f1565b833560048110156102fb5782526020840135906bffffffffffffffffffffffff19821682036102fb5782602092836040950152815201920191611792565b9080601f830112156102fb57604080519290611829908461135d565b8290604081019283116102fb57905b8282106118455750505090565b8135815260209182019101611838565b9080601f830112156102fb57813561186c816114f4565b9261187a604051948561135d565b81845260208085019260051b8201019283116102fb57602001905b8282106118a25750505090565b602080916118af84611473565b815201910190611895565b9190610220838203126102fb5760405161010081018181106001600160401b0382111761066e57604052809382601f820112156102fb576040516119006101008261135d565b806101008301918583116102fb5790859184905b848210611a9457505084526119289161180d565b602083015261193b83610140830161180d565b604083015261018081013561ffff811681036102fb5760608301526101a08101356001600160401b0381116102fb5783611976918301611855565b60808301526101c08101356001600160401b0381116102fb578361199b918301611855565b60a08301526101e08101356001600160401b0381116102fb57810183601f820112156102fb578035906119cd826114f4565b916119db604051938461135d565b80835260208084019160051b830101918683116102fb57602001905b828210611a845750505060c0830152610200810135906001600160401b0382116102fb57019180601f840112156102fb578235611a33816114f4565b93611a41604051958661135d565b81855260208086019260051b8201019283116102fb57602001905b828210611a6c5750505060e00152565b60208091611a7984611466565b815201910190611a5c565b81358152602091820191016119f7565b8135815287935060209182019101611914565b919060c0838203126102fb5760405190611ac082611327565b8193611acc8282611392565b835260608101356001600160401b0381116102fb5782611aed91830161150b565b6020840152611afe6080820161137e565b604084015260a0810135916001600160401b0383116102fb5760609261141092016118ba565b90815191606082526020611c81845160406060860152611b5d60a0860182516001600160401b0360208092828151168552015116910152565b610240611b7a848301516102c060e08901526103608801906112cd565b60408301516001600160401b031661010088015260608301516001600160801b03166101208801526080830151151561014088015260a08301518051610160890152602090810151805163ffffffff166101808a015201516101a08801529160c081015115156101c088015260e08101516101e08801526101008101511515610200880152610120810151610220880152610140810151828801526101608101516102608801526101808101516102808801526101a08101516102a08801526101c081015115156102c08801526101e08101516102e088015261020081015115156103008801526102208101516103208801520151609f19868303016103408701526112cd565b930151605f19838503016080840152602060e0606060c08701936001600160401b03815116885263ffffffff848201511684890152611ce2604082015160408a019060208060409280518552015163ffffffff815116828501520151910152565b01519560c060a08201528651809452019401905f5b818110611d255750505060209081015180516001600160401b03908116838501529101511660409091015290565b9091948551805160048110156112865782526020908101516bffffffffffffffffffffffff191681830152604090910195019190600101611cf7565b905f905b60028210611d7257505050565b6020806001928551815201930191019091611d65565b90602080835192838152019201905f5b818110611da55750505090565b825163ffffffff16845260209384019390920191600101611d98565b80515f835b60088210611e9157505050611de46020820151610100840190611d61565b611df76040820151610140840190611d61565b61ffff60608201511661018083015260e0611e50611e3d611e2a60808501516102206101a0880152610220870190611d88565b60a08501518682036101c0880152611d88565b60c08401518582036101e087015261129a565b91015191610200818303910152602080835192838152019201905f5b818110611e795750505090565b82511515845260209384019390920191600101611e6c565b6020806001928551815201930191019091611dc6565b906114f19160208152611edb602082018351604080916001600160801b038151168452602081015160208501520151910152565b6060611ef6602084015160c0608085015260e0840190611b24565b926001600160801b0360408201511660a084015201519060c0601f1982850301910152611dc1565b91908260609103126102fb57604051611f368161130c565b809280516001600160801b03811681036102fb5760409182918452602081015160208501520151910152565b51906001600160401b03821682036102fb57565b91908260409103126102fb57604051611f8e816112f1565b6020611410818395611f9f81611f62565b855201611f62565b90610140828203126102fb57611ffd9061010060405193611fc785611327565b611fd18382611f1e565b8552611fe08360608301611f1e565b6020860152611ff28360c08301611f76565b604086015201611f76565b606082015290565b1561200e575050565b906001600160401b0380927f7d939bd8000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b90612053826114f4565b612060604051918261135d565b8281528092612071601f19916114f4565b0190602036910137565b805182101561208f5760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b91906080838203126102fb57604051906120bc82611327565b819380356001600160401b0381116102fb576060926120dc9183016114d6565b8352602081013560208401526120f4604082016113cb565b60408401520135908160070b82036102fb5760600152565b908101906020818303126102fb578035906001600160401b0382116102fb57016040818303126102fb5760405190612143826112f1565b80356001600160401b0381116102fb578361215f918301611aa7565b82526020810135906001600160401b0382116102fb57016080818403126102fb576040519261218d84611327565b81356001600160401b0381116102fb57820181601f820112156102fb5780356121b5816114f4565b916121c3604051938461135d565b81835260208084019260051b820101918483116102fb5760208201905b8382106125f0575050505084526121f960208301611466565b60208501526040820135916001600160401b0383116102fb5761222360609261222e9483016120a3565b6040860152016113cb565b606083015260208101918252516122516001600160801b036040830151166128a8565b6122a56122966105de610e9a60608501519461227a602082019687519060606020610e5d612a75565b6040519283916343c3b39360e11b602084015260248301611ea7565b60208082518301019101611fa7565b916122af83612981565b926122b98461127c565b60018414612591576101606122ce8351612cc9565b935151510151808403612561575060207f1411889dcdfcf9af1fcf4ff4521cab4580a1fd842401ead543815cb1148d643893926001600160401b03926123138761127c565b866124d8578281612364606061247994019461235e88858851015116897f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025460401c16808211612005565b51612e90565b8351868151166fffffffffffffffffffffffffffffffff196fffffffffffffffff0000000000000000857f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025494015160401b16921617177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025501805160405161240d8682018093604080916001600160801b038151168452602081015160208501520151910152565b6060815261241c60808261135d565b51902061245f868686510151166001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b556001600160801b0385858551015116915151169061350c565b7f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb318284818451015116604051908152a1510151166124d260405192839283602090939291936001600160401b0360408201951681520152565b0390a190565b8061251d606061253893019361235e87878751015116887f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025460401c16808214612005565b6001600160801b03848681865101511692015151169061350c565b510151166124d260405192839283602090939291936001600160401b0360408201951681520152565b83907f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b5050506801000000000000000068ff0000000000000000195f516020614b925f395f51905f525416175f516020614b925f395f51905f52557f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a190565b81356001600160401b0381116102fb57602091612612888480948801016120a3565b8152019101906121e0565b903590601e19813603018212156102fb57018035906001600160401b0382116102fb576020019181360383136102fb57565b903590601e19813603018212156102fb57018035906001600160401b0382116102fb57602001918160051b360383136102fb57565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff16156126bc57565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260405f206001600160a01b0333165f5260205260ff60405f2054161561271d5750565b63e2517d3f60e01b5f523360045260245260445ffd5b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f205416155f146127ba57805f525f60205260405f206001600160a01b0383165f5260205260405f20600160ff198254161790556001600160a01b03339216907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f2054165f146127ba57805f525f60205260405f206001600160a01b0383165f5260205260405f2060ff1981541690556001600160a01b03339216907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b9190820180921161285057565b634e487b7160e01b5f52601160045260245ffd5b1561286d575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b9190820391821161285057565b6001600160801b03633b9aca0091160463ffffffff5f516020614b925f395f51905f525460481c166128e782426128df8442612843565b821115612864565b4282106128f2575050565b6128fc824261289b565b116129045750565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b5f918291602082519201905af43d15612979573d9061295182611415565b9161295f604051938461135d565b82523d5f602084013e5b156129715790565b602081519101fd5b606090612969565b6001600160401b03602060608301510151165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f205480155f146129cc5750505f90565b60208201908151604051612a01602082018093604080916001600160801b038151168452602081015160208501520151910152565b60608152612a1060808261135d565b5190201491821592612a2e575b505015612a2957600190565b600290565b6001600160801b0391925081905151169151511611155f80612a1d565b60405190612a5882611342565b5f6080838281528260208201528260408201528260608201520152565b612a7d612a4b565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1705547f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17065461ffff6001600160801b037f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170954169160405193612afe85611342565b84526001600160a01b03811660208501526001600160401b038160a01c16604085015260e01c166060830152608082015290565b60405190612b3f82611327565b5f6060838281528260208201528260408201520152565b612b5e612b32565b50602081016001600160a01b0381511615612bf3576001600160a01b03905116612b8f61ffff60608401511661453e565b813b6001811115612be0575f198101908111612850578111612bcd5790816001612bbb612bca94614563565b9260208401903c80925161458b565b91565b50633610565160e21b5f5260045260245ffd5b82633610565160e21b5f5260045260245ffd5b5051633a517eed60e21b5f5260045260245ffd5b939263ffffffff169161ffff16821015612c9957602c820293828504602c14831517156128505784603001918260301161285057850192605084015160e01c03612c6e575060348401811161285057603c605483015160c01c94011061285057605c015190565b7f4724a0fd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b50827ff81f5140000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b908151518015612e8a57612cdc81612049565b905f5b84518051821015612e7a57612cf68260209261207b565b5101516001600160401b036040612d0e84895161207b565b5101511660405191612d1f836112f1565b8252602082019080825260249080612e1a575b50612d3c90614563565b91600a6020840153600190818401906022602083015382800191828411612850576021600a9101536002830180921161285057612d8691612d7d9086614a95565b90519085614aab565b916001600160401b03815116612db7575b50505090612da6600192614661565b612db0828661207b565b5201612cdf565b828401906010602083015382840180941161285057516001600160401b03169260219190910191905b6080841015612dfd57505091612da691600194935391925f612d97565b91908084926080607f8397161781530192019060071c9290612de0565b6001905b6080811015612e6c575080600101600111612e595760250190818111612e475750612d3c612d32565b634e487b7160e01b5f5260116004525ffd5b50634e487b7160e01b5f5260116004525ffd5b60019060071c910190612e1e565b5050906114f19293505f906146b4565b505f9150565b90612e9a82612cc9565b91519081519182156134145760b4831161343c575f91602c8402848104602c03612850576030018060301161285057612ed290614563565b9060208201946056865361564160218401536256414c60228401536356414c3460238401538060081c602c840153602d830153612f16612f1187614b57565b614661565b602e8301525f604e8301535f604f8301535f9560305b84518810156130d657612f3f888661207b565b519560408701906001600160401b03825116156130aa5760205f9801975b8a811061305c5750612f7a906001600160401b0383511690612843565b968286019160ff8b60181c16602084015361ffff8b60101c16602184015362ffffff8b60081c16602284015363ffffffff8b166023840153600484018411612850576001600160401b03905160ff8160381c16602485015361ffff8160301c16602585015362ffffff8160281c16602685015363ffffffff8160201c16602785015364ffffffffff8160181c16602885015365ffffffffffff8160101c16602985015366ffffffffffffff8160081c16602a85015316602b830153600c8301831161285057602c9051910152602c810180911161285057600190970196612f2c565b6020613068828a61207b565b51015189511461307a57600101612f5d565b8a907faefa56c2000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b897fafef9f15000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b50949391925094506001600160401b038111613414578060ff6001600160401b039260381c16602484015361ffff8160301c16602584015362ffffff8160281c16602684015363ffffffff8160201c16602784015364ffffffffff8160181c16602884015365ffffffffffff8160101c16602984015366ffffffffffffff8160081c16602a84015316602b82015361316e818461458b565b916040519161319f60218460208101945f8652845180918484015e81015f838201520301601f19810185528461135d565b6160008351116133e85750906001600160a01b0391613257602a7fffff000000000000000000000000000000000000000000000000000000000000845160f01b169360405193849160208301967f6100000000000000000000000000000000000000000000000000000000000000885260218401527f80600a3d393df30000000000000000000000000000000000000000000000000060238401525180918484015e81015f838201520301601f19810183528261135d565b51905ff0169182156133c05760209273ffffffffffffffffffffffffffffffffffffffff197f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17065416177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706557f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055580517fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b807f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706549360a01b16169116177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655015161ffff60e01b1961ffff60e01b807f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706549360e01b16169116177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b827fab7bac08000000000000000000000000000000000000000000000000000000005f5260045260b460245260445ffd5b907f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17085482101561208f577f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17085f52600282901c7f06ff75c3bdc3f03474569b3e952aab1bfb9d41b303b5638523d3e1200dcd1ae2019160031b60181690565b919091805483101561208f575f52601860205f208360021c019260031b1690565b907f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170854806137d7575b506001600160a01b03600161357a846001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f2090565b01541615907f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170554906001600160801b0360027f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17065493604051906135dc82611342565b8152602081016001600160a01b038616815260408201956001600160401b038160a01c16875261ffff606084019160e01c16815284608084019616968787526136558a6001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f2090565b935184556001600160a01b036001850193511673ffffffffffffffffffffffffffffffffffffffff19845416178355517fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b8085549360a01b161691161782555161ffff60e01b1961ffff60e01b8084549360e01b16169116179055019151166fffffffffffffffffffffffffffffffff198254161790556fffffffffffffffffffffffffffffffff197f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17095416177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1709556137575750565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708546801000000000000000081101561066e578060016137b992017f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17085561346d565b6001600160401b0380839493549260031b9316831b921b1916179055565b5f19810190811161285057826001600160401b036137f76138119361346d565b90549060031b1c16806001600160401b0383161015612005565b5f613535565b999a9890939294969a979195975f9a60018d1015611286578c156138645760248c60ff8f7f112d89cc00000000000000000000000000000000000000000000000000000000835216600452fd5b909192939495968098999a9c50151580613e42575b15613e0b57602001356001600160401b03811681036102fb5761389c368b611392565b9061391860405160208101906138cf8286604080916001600160801b038151168452602081015160208501520151910152565b606081526138de60808261135d565b519020916001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b548015613de357808203613db55750506020810151808a03613d8557506001600160801b03613948915116614035565b5f905f5b888110613c9c575b505015613c5b5750506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001693843b156102fb5794929085926040519687957f56455e790000000000000000000000000000000000000000000000000000000087526064870190600488015260606024880152526084850160848560051b87010194825f90603e1981360301935b838310613bea57505050505050600319848403016044850152808352602083019260208260051b82010193835f92601e1982360301905b858510613a8457505050505050509181805f9403915afa8015613a7957613a64575b5035906001600160801b038216809203613a615750633b9aca00900490565b80fd5b613a719192505f9061135d565b5f905f613a42565b6040513d5f823e3d90fd5b9193959750919395601f198282030185528735838112156102fb57840190613ab06020820192806148b2565b8091936020845252604082019060408160051b8401019380935f915b838310613af357505050505050602080600192990195019501929091899796949592613a20565b909192939495603f19838203018652613b0c87836148fa565b803560028110156102fb578252613b3a613b2960208301836148e6565b60606020850152606084019061490e565b906040810135609e19823603018112156102fb576001936020938493613bdc9301916040818303910152613bce613bb0613b85613b778580614822565b60a0865260a0860191614802565b613b90878601611466565b151587850152613ba360408601866148e6565b848203604086015261490e565b92613bbd60608201611466565b1515606084015260808101906148e6565b90608081840391015261490e565b980196019493019190613acc565b9193959698509193966083198b82030182528735868112156102fb576020613c486001938683940190613c3b613c31613c2384806148b2565b604085526040850191614853565b9285810190614822565b9185818503910152614802565b9901920193019093918a989695936139e9565b613c986040519283927ffef760c7000000000000000000000000000000000000000000000000000000008452602060048501526024840191614853565b0390fd5b9091613cd0613cbf613cb8613cb2858d8c614775565b8061264f565b3691614797565b613cca368789614797565b90614aed565b15613d7b575061096a613ce7613cf1928a89614775565b602081019061261d565b805182518082149182613d65575b505015613d1157505060015f80613954565b90613c98613d53926040519384937f5f1ca3810000000000000000000000000000000000000000000000000000000085526040600486015260448501906112cd565b838103600319016024850152906112cd565b9091506020830120906020840120145f80613cff565b919060010161394c565b89907f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b877fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b5061ffff881115613879565b9190916001600160401b03602080606081875101510151950151015116613e73612a4b565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708548015613f9d575f905f198101908111612850575b808210613fb05750613ee46001600160401b03917f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17086134eb565b90549060031b1c16908111613f9d575f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f209260405191613f2b83611342565b8454948584526001600160801b03600260018301549261ffff6001600160a01b038516948560208a01526001600160401b038160a01c1660408a015260e01c166060880152015416608085015215613f8a57613f88939450614159565b565b84633a517eed60e21b5f5260045260245ffd5b633a517eed60e21b5f525f60045260245ffd5b90613fbb8282612843565b600181018091116128505760011c90836001600160401b03613ffd847f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17086134eb565b905460039190911b1c1611614013575090613eaa565b91505f19810190811115613eaa57634e487b7160e01b5f52601160045260245ffd5b6001600160801b03633b9aca009116045f516020614b925f395f51905f525461406c82426128df63ffffffff8560481c1642612843565b5f914281106140c1575b5063ffffffff16906001600160801b031681811015614093575050565b7fdb020436000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160801b03809293506140dc63ffffffff924261289b565b16929150614076565b156140ed5750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156141275750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b90919260808201515161ffff606084015116809114908161452f575b81614520575b81614511575b5015613414579161419184612b56565b6040860151955f95869485919082805b61ffff60608b0151168b101561430e576141bf8b60e08c015161207b565b51156143025763ffffffff6141d88c60808d015161207b565b511680996142ea575b50506001979363ffffffff6141fa8c60a08d015161207b565b5116906142118261ffff60208c015116811061411f565b6101008210156142be578160c06142618e6142578f95968f8f908f908f9e9d9c9b9a60209061ffff92861b9061424a87838316156140e5565b179f519301511691612c07565b939095015161207b565b510361429357506001600160401b038091169116016001600160401b038111612850576001909a5b01999091926141a1565b7fc3271bcd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b507ff6190fea000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b63ffffffff6142fb921681116140e5565b5f886141e1565b92919099600190614289565b50945096989550965050506001600160401b039150166003810281810460031482151715612850576801fffffffffffffffe8360011b16906001600160401b03841682046002146001600160401b0385161517156128505711156144d8575050608083015160e060c08501519401518151815180911490816144cd575b5015613414575f5b82518110156144c5576143a6818361207b565b51156144bd5763ffffffff6143bb828561207b565b51166143c7828861207b565b516143d5828751811061411f565b6143df828761207b565b5151600481101561128657600119016144915760205f6bffffffffffffffffffffffff198261440e868b61207b565b510151169260405183810191825283815261442a60408261135d565b604051918291518091835e8101838152039060025afa15613a79575f516bffffffffffffffffffffffff19160361446657506001905b01614393565b7fdc59b237000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b507f5b5c80eb000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b600190614460565b505050509050565b90508551145f61438b565b6001600160401b0392507f6e3083c3000000000000000000000000000000000000000000000000000000005f526004521660245260445ffd5b905060e083015151145f614181565b60c0840151518114915061417b565b60a08401515181149150614175565b61ffff16602c810290808204602c149015171561285057603001806030116128505790565b9061456d82611415565b61457a604051918261135d565b8281528092612071601f1991611415565b919091614596612b32565b506030835110612c6e576356414c34602084015160e01c03612c6e57602483015160c01c602c84015160f01c602e85015190604e86015160f01c92604051906145de82611327565b815281602082015260408101928352836060820152958115938415614656575b841561464c575b508315614635575b5050811561461e575b50612c6e5750565b90505161462d612f1183614b57565b14155f614616565b519192506146429061453e565b1415905f8061460d565b151593505f614605565b60b4831194506145fe565b8051600181018091116128505761467790614563565b80511561208f5760209181614693845f94019284845382614ac3565b50604051918291518091835e8101838152039060025afa15613a79575f5190565b9291926146c1828561289b565b93600185146147655760015b8060011b90868210156146e057506146cd565b9394955050906146fa6146f38486612843565b85836146b4565b61470e936147089195612843565b906146b4565b9060405161471d60808261135d565b604181526020810190606036833780511561208f576020935f936001845360218301526041820152604051918291518091835e8101838152039060025afa15613a79575f5190565b5061477192935061207b565b5190565b919081101561208f5760051b81013590603e19813603018212156102fb570190565b9291906147a3816114f4565b936147b1604051958661135d565b602085838152019160051b8101918383116102fb5781905b8382106147d7575050505050565b81356001600160401b0381116102fb576020916147f787849387016114d6565b8152019101906147c9565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156102fb5701602081359101916001600160401b0382116102fb5781360383136102fb57565b90602083828152019260208260051b82010193835f925b84841061487a5750505050505090565b9091929394956020806148a2600193601f1986820301885261489c8b88614822565b90614802565b980194019401929493919061486a565b9035601e19823603018112156102fb5701602081359101916001600160401b0382116102fb578160051b360383136102fb57565b9035607e19823603018112156102fb570190565b9035605e19823603018112156102fb570190565b61494761492c61491e8380614822565b608086526080860191614802565b6149396020840184614822565b908583036020870152614802565b61495460408301836148e6565b908381036040850152813560038110156102fb576149718161127c565b8152602082013560038110156102fb5761498a8161127c565b602082015260408201359060038210156102fb5760806149c66149e194846149b76149d69699989961127c565b60408501526060810190614822565b9190928160608201520191614802565b9260608101906148b2565b90916060818503910152808352602083019060208160051b85010193835f915b838310614a115750505050505090565b909192939495601f19828203018652614a2a87846148fa565b80359160038310156102fb57614a87602092839285614a4a60019761127c565b8152614a79614a6e614a5e86850185614822565b6060888601526060850191614802565b926040810190614822565b916040818503910152614802565b980196019493019190614a01565b6020828192010153600181018091116128505790565b81602091939293010152602081018091116128505790565b908051918215614ae5576021602084930191015e600101806001116128505790565b505050600190565b9081518151036127ba575f5b8251811015614ae557614b0c818461207b565b5151614b18828461207b565b515103614b5057614b29818461207b565b5160208151910120614b3b828461207b565b516020815191012003614b5057600101614af9565b5050505f90565b614b8d60405191614b6960608461135d565b602283526040366020850137600a6020840153614b87600184614a95565b83614aab565b509056fe5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1703a164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17085e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17014c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e05e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17036e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1707ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb55e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1702",
}

// ContractSpectreClientABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractSpectreClientMetaData.ABI instead.
var ContractSpectreClientABI = ContractSpectreClientMetaData.ABI

// ContractSpectreClientBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractSpectreClientMetaData.Bin instead.
var ContractSpectreClientBin = ContractSpectreClientMetaData.Bin

// DeployContractSpectreClient deploys a new Ethereum contract, binding an instance of ContractSpectreClient to it.
func DeployContractSpectreClient(auth *bind.TransactOpts, backend bind.ContractBackend, updateClientModule common.Address, membershipModule common.Address, misbehaviourModule common.Address, clientState_ []byte, consensusState_ SpectreMsgsConsensusState, initialPinnedValidatorSet SpectreMsgsValidatorSet, roleManager common.Address) (common.Address, *types.Transaction, *ContractSpectreClient, error) {
	parsed, err := ContractSpectreClientMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractSpectreClientBin), backend, updateClientModule, membershipModule, misbehaviourModule, clientState_, consensusState_, initialPinnedValidatorSet, roleManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractSpectreClient{ContractSpectreClientCaller: ContractSpectreClientCaller{contract: contract}, ContractSpectreClientTransactor: ContractSpectreClientTransactor{contract: contract}, ContractSpectreClientFilterer: ContractSpectreClientFilterer{contract: contract}}, nil
}

// ContractSpectreClient is an auto generated Go binding around an Ethereum contract.
type ContractSpectreClient struct {
	ContractSpectreClientCaller     // Read-only binding to the contract
	ContractSpectreClientTransactor // Write-only binding to the contract
	ContractSpectreClientFilterer   // Log filterer for contract events
}

// ContractSpectreClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractSpectreClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSpectreClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractSpectreClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSpectreClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractSpectreClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSpectreClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSpectreClientSession struct {
	Contract     *ContractSpectreClient // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ContractSpectreClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractSpectreClientCallerSession struct {
	Contract *ContractSpectreClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ContractSpectreClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractSpectreClientTransactorSession struct {
	Contract     *ContractSpectreClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ContractSpectreClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractSpectreClientRaw struct {
	Contract *ContractSpectreClient // Generic contract binding to access the raw methods on
}

// ContractSpectreClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractSpectreClientCallerRaw struct {
	Contract *ContractSpectreClientCaller // Generic read-only contract binding to access the raw methods on
}

// ContractSpectreClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractSpectreClientTransactorRaw struct {
	Contract *ContractSpectreClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractSpectreClient creates a new instance of ContractSpectreClient, bound to a specific deployed contract.
func NewContractSpectreClient(address common.Address, backend bind.ContractBackend) (*ContractSpectreClient, error) {
	contract, err := bindContractSpectreClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractSpectreClient{ContractSpectreClientCaller: ContractSpectreClientCaller{contract: contract}, ContractSpectreClientTransactor: ContractSpectreClientTransactor{contract: contract}, ContractSpectreClientFilterer: ContractSpectreClientFilterer{contract: contract}}, nil
}

// NewContractSpectreClientCaller creates a new read-only instance of ContractSpectreClient, bound to a specific deployed contract.
func NewContractSpectreClientCaller(address common.Address, caller bind.ContractCaller) (*ContractSpectreClientCaller, error) {
	contract, err := bindContractSpectreClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractSpectreClientCaller{contract: contract}, nil
}

// NewContractSpectreClientTransactor creates a new write-only instance of ContractSpectreClient, bound to a specific deployed contract.
func NewContractSpectreClientTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractSpectreClientTransactor, error) {
	contract, err := bindContractSpectreClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractSpectreClientTransactor{contract: contract}, nil
}

// NewContractSpectreClientFilterer creates a new log filterer instance of ContractSpectreClient, bound to a specific deployed contract.
func NewContractSpectreClientFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractSpectreClientFilterer, error) {
	contract, err := bindContractSpectreClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractSpectreClientFilterer{contract: contract}, nil
}

// bindContractSpectreClient binds a generic wrapper to an already deployed contract.
func bindContractSpectreClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractSpectreClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractSpectreClient *ContractSpectreClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractSpectreClient.Contract.ContractSpectreClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractSpectreClient *ContractSpectreClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.ContractSpectreClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractSpectreClient *ContractSpectreClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.ContractSpectreClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractSpectreClient *ContractSpectreClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractSpectreClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractSpectreClient *ContractSpectreClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractSpectreClient *ContractSpectreClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractSpectreClient.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ContractSpectreClient.Contract.DEFAULTADMINROLE(&_ContractSpectreClient.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ContractSpectreClient.Contract.DEFAULTADMINROLE(&_ContractSpectreClient.CallOpts)
}

// MISBEHAVIOURSUBMITTERROLE is a free data retrieval call binding the contract method 0xdb3e1fa4.
//
// Solidity: function MISBEHAVIOUR_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientCaller) MISBEHAVIOURSUBMITTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractSpectreClient.contract.Call(opts, &out, "MISBEHAVIOUR_SUBMITTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MISBEHAVIOURSUBMITTERROLE is a free data retrieval call binding the contract method 0xdb3e1fa4.
//
// Solidity: function MISBEHAVIOUR_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientSession) MISBEHAVIOURSUBMITTERROLE() ([32]byte, error) {
	return _ContractSpectreClient.Contract.MISBEHAVIOURSUBMITTERROLE(&_ContractSpectreClient.CallOpts)
}

// MISBEHAVIOURSUBMITTERROLE is a free data retrieval call binding the contract method 0xdb3e1fa4.
//
// Solidity: function MISBEHAVIOUR_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientCallerSession) MISBEHAVIOURSUBMITTERROLE() ([32]byte, error) {
	return _ContractSpectreClient.Contract.MISBEHAVIOURSUBMITTERROLE(&_ContractSpectreClient.CallOpts)
}

// GetClientState is a free data retrieval call binding the contract method 0xef913a4b.
//
// Solidity: function getClientState() view returns(bytes)
func (_ContractSpectreClient *ContractSpectreClientCaller) GetClientState(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractSpectreClient.contract.Call(opts, &out, "getClientState")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetClientState is a free data retrieval call binding the contract method 0xef913a4b.
//
// Solidity: function getClientState() view returns(bytes)
func (_ContractSpectreClient *ContractSpectreClientSession) GetClientState() ([]byte, error) {
	return _ContractSpectreClient.Contract.GetClientState(&_ContractSpectreClient.CallOpts)
}

// GetClientState is a free data retrieval call binding the contract method 0xef913a4b.
//
// Solidity: function getClientState() view returns(bytes)
func (_ContractSpectreClient *ContractSpectreClientCallerSession) GetClientState() ([]byte, error) {
	return _ContractSpectreClient.Contract.GetClientState(&_ContractSpectreClient.CallOpts)
}

// GetPinnedValidatorSet is a free data retrieval call binding the contract method 0x6edfe3af.
//
// Solidity: function getPinnedValidatorSet() view returns(uint32[] indices, bytes32[] pubkeys, uint64[] votingPowers)
func (_ContractSpectreClient *ContractSpectreClientCaller) GetPinnedValidatorSet(opts *bind.CallOpts) (struct {
	Indices      []uint32
	Pubkeys      [][32]byte
	VotingPowers []uint64
}, error) {
	var out []interface{}
	err := _ContractSpectreClient.contract.Call(opts, &out, "getPinnedValidatorSet")

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
func (_ContractSpectreClient *ContractSpectreClientSession) GetPinnedValidatorSet() (struct {
	Indices      []uint32
	Pubkeys      [][32]byte
	VotingPowers []uint64
}, error) {
	return _ContractSpectreClient.Contract.GetPinnedValidatorSet(&_ContractSpectreClient.CallOpts)
}

// GetPinnedValidatorSet is a free data retrieval call binding the contract method 0x6edfe3af.
//
// Solidity: function getPinnedValidatorSet() view returns(uint32[] indices, bytes32[] pubkeys, uint64[] votingPowers)
func (_ContractSpectreClient *ContractSpectreClientCallerSession) GetPinnedValidatorSet() (struct {
	Indices      []uint32
	Pubkeys      [][32]byte
	VotingPowers []uint64
}, error) {
	return _ContractSpectreClient.Contract.GetPinnedValidatorSet(&_ContractSpectreClient.CallOpts)
}

// GetPinnedValidatorsHash is a free data retrieval call binding the contract method 0x8c80fbda.
//
// Solidity: function getPinnedValidatorsHash() view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientCaller) GetPinnedValidatorsHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractSpectreClient.contract.Call(opts, &out, "getPinnedValidatorsHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPinnedValidatorsHash is a free data retrieval call binding the contract method 0x8c80fbda.
//
// Solidity: function getPinnedValidatorsHash() view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientSession) GetPinnedValidatorsHash() ([32]byte, error) {
	return _ContractSpectreClient.Contract.GetPinnedValidatorsHash(&_ContractSpectreClient.CallOpts)
}

// GetPinnedValidatorsHash is a free data retrieval call binding the contract method 0x8c80fbda.
//
// Solidity: function getPinnedValidatorsHash() view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientCallerSession) GetPinnedValidatorsHash() ([32]byte, error) {
	return _ContractSpectreClient.Contract.GetPinnedValidatorsHash(&_ContractSpectreClient.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ContractSpectreClient.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ContractSpectreClient.Contract.GetRoleAdmin(&_ContractSpectreClient.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ContractSpectreClient.Contract.GetRoleAdmin(&_ContractSpectreClient.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ContractSpectreClient *ContractSpectreClientCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ContractSpectreClient.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ContractSpectreClient *ContractSpectreClientSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ContractSpectreClient.Contract.HasRole(&_ContractSpectreClient.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ContractSpectreClient *ContractSpectreClientCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ContractSpectreClient.Contract.HasRole(&_ContractSpectreClient.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ContractSpectreClient *ContractSpectreClientCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ContractSpectreClient.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ContractSpectreClient *ContractSpectreClientSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ContractSpectreClient.Contract.SupportsInterface(&_ContractSpectreClient.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ContractSpectreClient *ContractSpectreClientCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ContractSpectreClient.Contract.SupportsInterface(&_ContractSpectreClient.CallOpts, interfaceId)
}

// UpgradeClient is a free data retrieval call binding the contract method 0x8a8e4c5d.
//
// Solidity: function upgradeClient(bytes ) pure returns()
func (_ContractSpectreClient *ContractSpectreClientCaller) UpgradeClient(opts *bind.CallOpts, arg0 []byte) error {
	var out []interface{}
	err := _ContractSpectreClient.contract.Call(opts, &out, "upgradeClient", arg0)

	if err != nil {
		return err
	}

	return err

}

// UpgradeClient is a free data retrieval call binding the contract method 0x8a8e4c5d.
//
// Solidity: function upgradeClient(bytes ) pure returns()
func (_ContractSpectreClient *ContractSpectreClientSession) UpgradeClient(arg0 []byte) error {
	return _ContractSpectreClient.Contract.UpgradeClient(&_ContractSpectreClient.CallOpts, arg0)
}

// UpgradeClient is a free data retrieval call binding the contract method 0x8a8e4c5d.
//
// Solidity: function upgradeClient(bytes ) pure returns()
func (_ContractSpectreClient *ContractSpectreClientCallerSession) UpgradeClient(arg0 []byte) error {
	return _ContractSpectreClient.Contract.UpgradeClient(&_ContractSpectreClient.CallOpts, arg0)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ContractSpectreClient *ContractSpectreClientTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSpectreClient.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ContractSpectreClient *ContractSpectreClientSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.GrantRole(&_ContractSpectreClient.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ContractSpectreClient *ContractSpectreClientTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.GrantRole(&_ContractSpectreClient.TransactOpts, role, account)
}

// Misbehaviour is a paid mutator transaction binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes misbehaviourMsg) returns()
func (_ContractSpectreClient *ContractSpectreClientTransactor) Misbehaviour(opts *bind.TransactOpts, misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractSpectreClient.contract.Transact(opts, "misbehaviour", misbehaviourMsg)
}

// Misbehaviour is a paid mutator transaction binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes misbehaviourMsg) returns()
func (_ContractSpectreClient *ContractSpectreClientSession) Misbehaviour(misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.Misbehaviour(&_ContractSpectreClient.TransactOpts, misbehaviourMsg)
}

// Misbehaviour is a paid mutator transaction binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes misbehaviourMsg) returns()
func (_ContractSpectreClient *ContractSpectreClientTransactorSession) Misbehaviour(misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.Misbehaviour(&_ContractSpectreClient.TransactOpts, misbehaviourMsg)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ContractSpectreClient *ContractSpectreClientTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ContractSpectreClient.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ContractSpectreClient *ContractSpectreClientSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.RenounceRole(&_ContractSpectreClient.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ContractSpectreClient *ContractSpectreClientTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.RenounceRole(&_ContractSpectreClient.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ContractSpectreClient *ContractSpectreClientTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSpectreClient.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ContractSpectreClient *ContractSpectreClientSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.RevokeRole(&_ContractSpectreClient.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ContractSpectreClient *ContractSpectreClientTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.RevokeRole(&_ContractSpectreClient.TransactOpts, role, account)
}

// Unfreeze is a paid mutator transaction binding the contract method 0x6a28f000.
//
// Solidity: function unfreeze() returns()
func (_ContractSpectreClient *ContractSpectreClientTransactor) Unfreeze(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractSpectreClient.contract.Transact(opts, "unfreeze")
}

// Unfreeze is a paid mutator transaction binding the contract method 0x6a28f000.
//
// Solidity: function unfreeze() returns()
func (_ContractSpectreClient *ContractSpectreClientSession) Unfreeze() (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.Unfreeze(&_ContractSpectreClient.TransactOpts)
}

// Unfreeze is a paid mutator transaction binding the contract method 0x6a28f000.
//
// Solidity: function unfreeze() returns()
func (_ContractSpectreClient *ContractSpectreClientTransactorSession) Unfreeze() (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.Unfreeze(&_ContractSpectreClient.TransactOpts)
}

// UpdateApplicationState is a paid mutator transaction binding the contract method 0x4b1872be.
//
// Solidity: function updateApplicationState(bytes updateMsg) returns(uint8)
func (_ContractSpectreClient *ContractSpectreClientTransactor) UpdateApplicationState(opts *bind.TransactOpts, updateMsg []byte) (*types.Transaction, error) {
	return _ContractSpectreClient.contract.Transact(opts, "updateApplicationState", updateMsg)
}

// UpdateApplicationState is a paid mutator transaction binding the contract method 0x4b1872be.
//
// Solidity: function updateApplicationState(bytes updateMsg) returns(uint8)
func (_ContractSpectreClient *ContractSpectreClientSession) UpdateApplicationState(updateMsg []byte) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.UpdateApplicationState(&_ContractSpectreClient.TransactOpts, updateMsg)
}

// UpdateApplicationState is a paid mutator transaction binding the contract method 0x4b1872be.
//
// Solidity: function updateApplicationState(bytes updateMsg) returns(uint8)
func (_ContractSpectreClient *ContractSpectreClientTransactorSession) UpdateApplicationState(updateMsg []byte) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.UpdateApplicationState(&_ContractSpectreClient.TransactOpts, updateMsg)
}

// UpdateConsensusState is a paid mutator transaction binding the contract method 0x92c19bdc.
//
// Solidity: function updateConsensusState(bytes updateMsg) returns(uint8)
func (_ContractSpectreClient *ContractSpectreClientTransactor) UpdateConsensusState(opts *bind.TransactOpts, updateMsg []byte) (*types.Transaction, error) {
	return _ContractSpectreClient.contract.Transact(opts, "updateConsensusState", updateMsg)
}

// UpdateConsensusState is a paid mutator transaction binding the contract method 0x92c19bdc.
//
// Solidity: function updateConsensusState(bytes updateMsg) returns(uint8)
func (_ContractSpectreClient *ContractSpectreClientSession) UpdateConsensusState(updateMsg []byte) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.UpdateConsensusState(&_ContractSpectreClient.TransactOpts, updateMsg)
}

// UpdateConsensusState is a paid mutator transaction binding the contract method 0x92c19bdc.
//
// Solidity: function updateConsensusState(bytes updateMsg) returns(uint8)
func (_ContractSpectreClient *ContractSpectreClientTransactorSession) UpdateConsensusState(updateMsg []byte) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.UpdateConsensusState(&_ContractSpectreClient.TransactOpts, updateMsg)
}

// VerifyMembership is a paid mutator transaction binding the contract method 0x974a74c4.
//
// Solidity: function verifyMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[],bytes) msg_) returns(uint256)
func (_ContractSpectreClient *ContractSpectreClientTransactor) VerifyMembership(opts *bind.TransactOpts, msg_ LightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractSpectreClient.contract.Transact(opts, "verifyMembership", msg_)
}

// VerifyMembership is a paid mutator transaction binding the contract method 0x974a74c4.
//
// Solidity: function verifyMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[],bytes) msg_) returns(uint256)
func (_ContractSpectreClient *ContractSpectreClientSession) VerifyMembership(msg_ LightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.VerifyMembership(&_ContractSpectreClient.TransactOpts, msg_)
}

// VerifyMembership is a paid mutator transaction binding the contract method 0x974a74c4.
//
// Solidity: function verifyMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[],bytes) msg_) returns(uint256)
func (_ContractSpectreClient *ContractSpectreClientTransactorSession) VerifyMembership(msg_ LightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.VerifyMembership(&_ContractSpectreClient.TransactOpts, msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0xa6f031bb.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[]) msg_) returns(uint256)
func (_ContractSpectreClient *ContractSpectreClientTransactor) VerifyNonMembership(opts *bind.TransactOpts, msg_ LightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
	return _ContractSpectreClient.contract.Transact(opts, "verifyNonMembership", msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0xa6f031bb.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[]) msg_) returns(uint256)
func (_ContractSpectreClient *ContractSpectreClientSession) VerifyNonMembership(msg_ LightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.VerifyNonMembership(&_ContractSpectreClient.TransactOpts, msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0xa6f031bb.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[]) msg_) returns(uint256)
func (_ContractSpectreClient *ContractSpectreClientTransactorSession) VerifyNonMembership(msg_ LightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.VerifyNonMembership(&_ContractSpectreClient.TransactOpts, msg_)
}

// ContractSpectreClientClientFrozenIterator is returned from FilterClientFrozen and is used to iterate over the raw logs and unpacked data for ClientFrozen events raised by the ContractSpectreClient contract.
type ContractSpectreClientClientFrozenIterator struct {
	Event *ContractSpectreClientClientFrozen // Event containing the contract specifics and raw log

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
func (it *ContractSpectreClientClientFrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSpectreClientClientFrozen)
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
		it.Event = new(ContractSpectreClientClientFrozen)
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
func (it *ContractSpectreClientClientFrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSpectreClientClientFrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSpectreClientClientFrozen represents a ClientFrozen event raised by the ContractSpectreClient contract.
type ContractSpectreClientClientFrozen struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterClientFrozen is a free log retrieval operation binding the contract event 0x59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c.
//
// Solidity: event ClientFrozen()
func (_ContractSpectreClient *ContractSpectreClientFilterer) FilterClientFrozen(opts *bind.FilterOpts) (*ContractSpectreClientClientFrozenIterator, error) {

	logs, sub, err := _ContractSpectreClient.contract.FilterLogs(opts, "ClientFrozen")
	if err != nil {
		return nil, err
	}
	return &ContractSpectreClientClientFrozenIterator{contract: _ContractSpectreClient.contract, event: "ClientFrozen", logs: logs, sub: sub}, nil
}

// WatchClientFrozen is a free log subscription operation binding the contract event 0x59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c.
//
// Solidity: event ClientFrozen()
func (_ContractSpectreClient *ContractSpectreClientFilterer) WatchClientFrozen(opts *bind.WatchOpts, sink chan<- *ContractSpectreClientClientFrozen) (event.Subscription, error) {

	logs, sub, err := _ContractSpectreClient.contract.WatchLogs(opts, "ClientFrozen")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSpectreClientClientFrozen)
				if err := _ContractSpectreClient.contract.UnpackLog(event, "ClientFrozen", log); err != nil {
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

// ParseClientFrozen is a log parse operation binding the contract event 0x59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c.
//
// Solidity: event ClientFrozen()
func (_ContractSpectreClient *ContractSpectreClientFilterer) ParseClientFrozen(log types.Log) (*ContractSpectreClientClientFrozen, error) {
	event := new(ContractSpectreClientClientFrozen)
	if err := _ContractSpectreClient.contract.UnpackLog(event, "ClientFrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractSpectreClientClientUnfrozenIterator is returned from FilterClientUnfrozen and is used to iterate over the raw logs and unpacked data for ClientUnfrozen events raised by the ContractSpectreClient contract.
type ContractSpectreClientClientUnfrozenIterator struct {
	Event *ContractSpectreClientClientUnfrozen // Event containing the contract specifics and raw log

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
func (it *ContractSpectreClientClientUnfrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSpectreClientClientUnfrozen)
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
		it.Event = new(ContractSpectreClientClientUnfrozen)
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
func (it *ContractSpectreClientClientUnfrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSpectreClientClientUnfrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSpectreClientClientUnfrozen represents a ClientUnfrozen event raised by the ContractSpectreClient contract.
type ContractSpectreClientClientUnfrozen struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterClientUnfrozen is a free log retrieval operation binding the contract event 0x8e0da210a2ffecc744373b4860cd9b8dc471810f5fb61b7d2c3234ccd4aae30d.
//
// Solidity: event ClientUnfrozen()
func (_ContractSpectreClient *ContractSpectreClientFilterer) FilterClientUnfrozen(opts *bind.FilterOpts) (*ContractSpectreClientClientUnfrozenIterator, error) {

	logs, sub, err := _ContractSpectreClient.contract.FilterLogs(opts, "ClientUnfrozen")
	if err != nil {
		return nil, err
	}
	return &ContractSpectreClientClientUnfrozenIterator{contract: _ContractSpectreClient.contract, event: "ClientUnfrozen", logs: logs, sub: sub}, nil
}

// WatchClientUnfrozen is a free log subscription operation binding the contract event 0x8e0da210a2ffecc744373b4860cd9b8dc471810f5fb61b7d2c3234ccd4aae30d.
//
// Solidity: event ClientUnfrozen()
func (_ContractSpectreClient *ContractSpectreClientFilterer) WatchClientUnfrozen(opts *bind.WatchOpts, sink chan<- *ContractSpectreClientClientUnfrozen) (event.Subscription, error) {

	logs, sub, err := _ContractSpectreClient.contract.WatchLogs(opts, "ClientUnfrozen")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSpectreClientClientUnfrozen)
				if err := _ContractSpectreClient.contract.UnpackLog(event, "ClientUnfrozen", log); err != nil {
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

// ParseClientUnfrozen is a log parse operation binding the contract event 0x8e0da210a2ffecc744373b4860cd9b8dc471810f5fb61b7d2c3234ccd4aae30d.
//
// Solidity: event ClientUnfrozen()
func (_ContractSpectreClient *ContractSpectreClientFilterer) ParseClientUnfrozen(log types.Log) (*ContractSpectreClientClientUnfrozen, error) {
	event := new(ContractSpectreClientClientUnfrozen)
	if err := _ContractSpectreClient.contract.UnpackLog(event, "ClientUnfrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractSpectreClientClientUpdatedIterator is returned from FilterClientUpdated and is used to iterate over the raw logs and unpacked data for ClientUpdated events raised by the ContractSpectreClient contract.
type ContractSpectreClientClientUpdatedIterator struct {
	Event *ContractSpectreClientClientUpdated // Event containing the contract specifics and raw log

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
func (it *ContractSpectreClientClientUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSpectreClientClientUpdated)
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
		it.Event = new(ContractSpectreClientClientUpdated)
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
func (it *ContractSpectreClientClientUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSpectreClientClientUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSpectreClientClientUpdated represents a ClientUpdated event raised by the ContractSpectreClient contract.
type ContractSpectreClientClientUpdated struct {
	Height uint64
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterClientUpdated is a free log retrieval operation binding the contract event 0x9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb31.
//
// Solidity: event ClientUpdated(uint64 height)
func (_ContractSpectreClient *ContractSpectreClientFilterer) FilterClientUpdated(opts *bind.FilterOpts) (*ContractSpectreClientClientUpdatedIterator, error) {

	logs, sub, err := _ContractSpectreClient.contract.FilterLogs(opts, "ClientUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractSpectreClientClientUpdatedIterator{contract: _ContractSpectreClient.contract, event: "ClientUpdated", logs: logs, sub: sub}, nil
}

// WatchClientUpdated is a free log subscription operation binding the contract event 0x9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb31.
//
// Solidity: event ClientUpdated(uint64 height)
func (_ContractSpectreClient *ContractSpectreClientFilterer) WatchClientUpdated(opts *bind.WatchOpts, sink chan<- *ContractSpectreClientClientUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractSpectreClient.contract.WatchLogs(opts, "ClientUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSpectreClientClientUpdated)
				if err := _ContractSpectreClient.contract.UnpackLog(event, "ClientUpdated", log); err != nil {
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

// ParseClientUpdated is a log parse operation binding the contract event 0x9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb31.
//
// Solidity: event ClientUpdated(uint64 height)
func (_ContractSpectreClient *ContractSpectreClientFilterer) ParseClientUpdated(log types.Log) (*ContractSpectreClientClientUpdated, error) {
	event := new(ContractSpectreClientClientUpdated)
	if err := _ContractSpectreClient.contract.UnpackLog(event, "ClientUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractSpectreClientConsensusStateUpdatedIterator is returned from FilterConsensusStateUpdated and is used to iterate over the raw logs and unpacked data for ConsensusStateUpdated events raised by the ContractSpectreClient contract.
type ContractSpectreClientConsensusStateUpdatedIterator struct {
	Event *ContractSpectreClientConsensusStateUpdated // Event containing the contract specifics and raw log

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
func (it *ContractSpectreClientConsensusStateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSpectreClientConsensusStateUpdated)
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
		it.Event = new(ContractSpectreClientConsensusStateUpdated)
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
func (it *ContractSpectreClientConsensusStateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSpectreClientConsensusStateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSpectreClientConsensusStateUpdated represents a ConsensusStateUpdated event raised by the ContractSpectreClient contract.
type ContractSpectreClientConsensusStateUpdated struct {
	Height         uint64
	ValidatorsHash [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterConsensusStateUpdated is a free log retrieval operation binding the contract event 0x1411889dcdfcf9af1fcf4ff4521cab4580a1fd842401ead543815cb1148d6438.
//
// Solidity: event ConsensusStateUpdated(uint64 height, bytes32 validatorsHash)
func (_ContractSpectreClient *ContractSpectreClientFilterer) FilterConsensusStateUpdated(opts *bind.FilterOpts) (*ContractSpectreClientConsensusStateUpdatedIterator, error) {

	logs, sub, err := _ContractSpectreClient.contract.FilterLogs(opts, "ConsensusStateUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractSpectreClientConsensusStateUpdatedIterator{contract: _ContractSpectreClient.contract, event: "ConsensusStateUpdated", logs: logs, sub: sub}, nil
}

// WatchConsensusStateUpdated is a free log subscription operation binding the contract event 0x1411889dcdfcf9af1fcf4ff4521cab4580a1fd842401ead543815cb1148d6438.
//
// Solidity: event ConsensusStateUpdated(uint64 height, bytes32 validatorsHash)
func (_ContractSpectreClient *ContractSpectreClientFilterer) WatchConsensusStateUpdated(opts *bind.WatchOpts, sink chan<- *ContractSpectreClientConsensusStateUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractSpectreClient.contract.WatchLogs(opts, "ConsensusStateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSpectreClientConsensusStateUpdated)
				if err := _ContractSpectreClient.contract.UnpackLog(event, "ConsensusStateUpdated", log); err != nil {
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

// ParseConsensusStateUpdated is a log parse operation binding the contract event 0x1411889dcdfcf9af1fcf4ff4521cab4580a1fd842401ead543815cb1148d6438.
//
// Solidity: event ConsensusStateUpdated(uint64 height, bytes32 validatorsHash)
func (_ContractSpectreClient *ContractSpectreClientFilterer) ParseConsensusStateUpdated(log types.Log) (*ContractSpectreClientConsensusStateUpdated, error) {
	event := new(ContractSpectreClientConsensusStateUpdated)
	if err := _ContractSpectreClient.contract.UnpackLog(event, "ConsensusStateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractSpectreClientRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ContractSpectreClient contract.
type ContractSpectreClientRoleAdminChangedIterator struct {
	Event *ContractSpectreClientRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ContractSpectreClientRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSpectreClientRoleAdminChanged)
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
		it.Event = new(ContractSpectreClientRoleAdminChanged)
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
func (it *ContractSpectreClientRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSpectreClientRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSpectreClientRoleAdminChanged represents a RoleAdminChanged event raised by the ContractSpectreClient contract.
type ContractSpectreClientRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ContractSpectreClient *ContractSpectreClientFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ContractSpectreClientRoleAdminChangedIterator, error) {

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

	logs, sub, err := _ContractSpectreClient.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ContractSpectreClientRoleAdminChangedIterator{contract: _ContractSpectreClient.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ContractSpectreClient *ContractSpectreClientFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ContractSpectreClientRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _ContractSpectreClient.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSpectreClientRoleAdminChanged)
				if err := _ContractSpectreClient.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_ContractSpectreClient *ContractSpectreClientFilterer) ParseRoleAdminChanged(log types.Log) (*ContractSpectreClientRoleAdminChanged, error) {
	event := new(ContractSpectreClientRoleAdminChanged)
	if err := _ContractSpectreClient.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractSpectreClientRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ContractSpectreClient contract.
type ContractSpectreClientRoleGrantedIterator struct {
	Event *ContractSpectreClientRoleGranted // Event containing the contract specifics and raw log

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
func (it *ContractSpectreClientRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSpectreClientRoleGranted)
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
		it.Event = new(ContractSpectreClientRoleGranted)
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
func (it *ContractSpectreClientRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSpectreClientRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSpectreClientRoleGranted represents a RoleGranted event raised by the ContractSpectreClient contract.
type ContractSpectreClientRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractSpectreClient *ContractSpectreClientFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ContractSpectreClientRoleGrantedIterator, error) {

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

	logs, sub, err := _ContractSpectreClient.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ContractSpectreClientRoleGrantedIterator{contract: _ContractSpectreClient.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractSpectreClient *ContractSpectreClientFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ContractSpectreClientRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ContractSpectreClient.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSpectreClientRoleGranted)
				if err := _ContractSpectreClient.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_ContractSpectreClient *ContractSpectreClientFilterer) ParseRoleGranted(log types.Log) (*ContractSpectreClientRoleGranted, error) {
	event := new(ContractSpectreClientRoleGranted)
	if err := _ContractSpectreClient.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractSpectreClientRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ContractSpectreClient contract.
type ContractSpectreClientRoleRevokedIterator struct {
	Event *ContractSpectreClientRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ContractSpectreClientRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSpectreClientRoleRevoked)
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
		it.Event = new(ContractSpectreClientRoleRevoked)
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
func (it *ContractSpectreClientRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSpectreClientRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSpectreClientRoleRevoked represents a RoleRevoked event raised by the ContractSpectreClient contract.
type ContractSpectreClientRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractSpectreClient *ContractSpectreClientFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ContractSpectreClientRoleRevokedIterator, error) {

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

	logs, sub, err := _ContractSpectreClient.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ContractSpectreClientRoleRevokedIterator{contract: _ContractSpectreClient.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractSpectreClient *ContractSpectreClientFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ContractSpectreClientRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ContractSpectreClient.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSpectreClientRoleRevoked)
				if err := _ContractSpectreClient.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_ContractSpectreClient *ContractSpectreClientFilterer) ParseRoleRevoked(log types.Log) (*ContractSpectreClientRoleRevoked, error) {
	event := new(ContractSpectreClientRoleRevoked)
	if err := _ContractSpectreClient.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
