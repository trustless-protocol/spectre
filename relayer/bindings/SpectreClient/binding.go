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

// IICS02ClientMsgsHeight is an auto generated low-level Go binding around an user-defined struct.
type IICS02ClientMsgsHeight struct {
	RevisionNumber uint64
	RevisionHeight uint64
}

// ILightClientMsgsMsgVerifyMembership is an auto generated low-level Go binding around an user-defined struct.
type ILightClientMsgsMsgVerifyMembership struct {
	Height                IICS02ClientMsgsHeight
	KvPairs               []IMembershipMsgsKVPair
	MerkleProofs          []IMembershipMsgsMerkleProof
	AppHash               [32]byte
	TrustedConsensusState ISpectreTendermintMsgsConsensusState
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
	TrustedConsensusState ISpectreTendermintMsgsConsensusState
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

// ISpectreTendermintMsgsConsensusState is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsConsensusState struct {
	Timestamp          *big.Int
	Root               [32]byte
	NextValidatorsHash [32]byte
}

// ISpectreTendermintMsgsValidatorInfo is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsValidatorInfo struct {
	ValAddress       []byte
	PubKey           [32]byte
	VotingPower      uint64
	ProposerPriority int64
}

// ISpectreTendermintMsgsValidatorSet is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsValidatorSet struct {
	Validators       []ISpectreTendermintMsgsValidatorInfo
	HasProposer      bool
	Proposer         ISpectreTendermintMsgsValidatorInfo
	TotalVotingPower uint64
}

// ContractSpectreClientMetaData contains all meta data concerning the ContractSpectreClient contract.
var ContractSpectreClientMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"updateClientModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membershipModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviourModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"clientState_\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"consensusState_\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"initialPinnedValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structISpectreTendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConsensusStateHash\",\"inputs\":[{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPinnedSnapshot\",\"inputs\":[],\"outputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPinnedValidatorSet\",\"inputs\":[],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPinnedValidatorsHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateApplicationState\",\"inputs\":[{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateConsensusState\",\"inputs\":[{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ClientFrozen\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClientUnfrozen\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClientUpdated\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConsensusStateUpdated\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientNotFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DirectCallNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"DuplicateValidatorPubkey\",\"inputs\":[{\"name\":\"firstIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"secondIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GenesisPinnedValidatorSetMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ModuleHasNoCode\",\"inputs\":[{\"name\":\"module\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NonMonotonicHeightUpdate\",\"inputs\":[{\"name\":\"latestHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"updateHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"PinnedIndexOverflowsBitmask\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"PinnedValidatorSetStale\",\"inputs\":[{\"name\":\"expectedNextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actualPinnedValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofSignerCommitSigMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ProofSignerValidatorMismatch\",\"inputs\":[{\"name\":\"signerIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorCountExceedsLimit\",\"inputs\":[{\"name\":\"validatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxValidatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroVotingPower\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	Bin: "0x6101008060405234610f5357616918803803809161001d8285610fba565b833981018181036101208112610f535761003683610fdd565b9161004360208501610fdd565b9261005060408601610fdd565b60608601519093906001600160401b038111610f5357610074846060928901611042565b91607f190112610f535760405193606085016001600160401b03811186821017610cd05760405260808701516001600160801b0381168103610f5357855260a0870151906020860191825260c0880151946040870195865260e089015160018060401b038111610f53578901608081830312610f5357604051996100f78b610f84565b81516001600160401b038111610f5357820183601f82011215610f5357805161011f8161105f565b9161012d6040519384610fba565b81835260208084019260051b82010191868311610f535760208201905b838210610f5757505050508b52610163602083016110f5565b60208c015260408201516001600160401b038111610f53576060838d60406101966101aa9861019f96610100990161108a565b91015201611076565b60608c015201610fdd565b967f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a6123360e0526101d88561114f565b6101e18161114f565b6101ea8261114f565b83518401946020860194602081880312610f53576020810151906001600160401b038211610f53570195869003601f198101949061012013610f53576040519560e087016001600160401b03811188821017610cd05760405260208801516001600160401b038111610f53576020908901019080601f83011215610f535781516102769260200161100c565b865260408512610f5357604080519561028e87610f9f565b610299828a01611102565b87526102a760608a01611102565b602088015260208801968752603f190112610f5357604051976102c989610f9f565b6102d560808901611076565b89526102e360a08901611076565b60208a0152604087019889526102fb60c08901611110565b966060810197885261030f60e08a01611110565b9860808201998a5261033861012061032a61010084016110f5565b9260a0850193845201611110565b60c0830190815282518051919991906001600160401b038211610cd0575f5160206168985f395f51905f5254600181811c91168015610f49575b6020821014610f3557601f8111610ebb575b50602090601f8311600114610e39576104d2949392915f9183610e2e575b50508160011b915f199060031b1c1916175f5160206168985f395f51905f52555b5160ff81511661ff0060205f5160206167b85f395f51905f525493015160081b169161ffff191617175f5160206167b85f395f51905f52558b5160018060401b038151165f5160206168f85f395f51905f525491602068010000000000000000600160801b0391015160401b169160018060801b03191617175f5160206168f85f395f51905f52558a63ffffffff8b51169168ff000000000000000067ffffffff000000005f5160206167f85f395f51905f5254935160201b169151151560401b916cffffffff0000000000000000008d5160481b1693856cffffffff000000000000000000199160018060481b031916171617911617175f5160206167f85f395f51905f5255801515611121565b6104e563ffffffff885116801515611121565b516001600160401b03906020906104fb906111b8565b01518a51516001600160401b03169116818103610e1957505061051d8c6113ac565b8151808203610e0457505060405190602082019260018060801b038c5116845251604083015251606082015260608152610558608082610fba565b51902060018060401b03602089510151165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f205560805260a05260c05263ffffffff80835116915116019063ffffffff82116107a45763ffffffff80809451169384925116921611610def5750506105d7846113ac565b9351908151938415610dc85760b48511610dd7575f93602c8602868104602c036107a457603001806030116107a45761060f90611837565b9260208401976056895361564160218601536256414c60228601536356414c346023860153602c8501978060081c8953602d860153610655610650826119e0565b611869565b97602e8601988952604e8601915f83535f604f8801535f9a60305b89518d101561080c578c610684818c611398565b51604081018051909d91906001600160401b0316156107f957602001908d5f8e5b8582106107b8575050516106c4916001600160401b03909116906111ab565b9c63ffffffff848d019360ff8160181c16602086015361ffff8160101c16602186015362ffffff8160081c1660228601531660238401536004840184116107a4575160ff8160381c16602484015361ffff8160301c16602584015362ffffff8160281c16602684015363ffffffff8160201c16602784015364ffffffffff8160181c16602884015365ffffffffffff8160101c16602984015366ffffffffffffff8160081c16602a8401536001600160401b0316602b830153600c830183116107a457602c9051910152602c81018091116107a4576001909c019b610670565b634e487b7160e01b5f52601160045260245ffd5b8160209293506107cb9196949596611398565b5101518451146107e4576001018e908e949392946106a5565b63577d2b6160e11b5f5260045260245260445ffd5b8263afef9f1560e01b5f5260045260245ffd5b5087908b8b6001600160401b038111610dc857602484019060ff8160381c16825361ffff8160301c16602586015362ffffff8160281c16602686015363ffffffff8160201c16602786015364ffffffffff8160181c16602886015365ffffffffffff8160101c16602986015366ffffffffffffff8160081c16602a8601536001600160401b0316602b8501535f60606040516108a781610f84565b82815282602082015282604082015201526030845110610db5576356414c34835160e01c03610db5575160c01c945160f01c9051955160f01c95604051956108ee87610f84565b865260208601968288526040870191825280606088015282159283918415610daa575b8415610da0575b508315610d75575b50508115610d5e575b50610d4b576040519161095e60218460208101945f8652845180918484015e81015f838201520301601f198101855284610fba565b616000835111610d3857506109bf602a61ffff60f01b845160f01b16936040519384916020830196606160f81b885260218401526680600a3d393df360c81b60238401525180918484015e81015f838201520301601f198101835282610fba565b51905ff06001600160a01b03168015610d29575f5160206168385f395f51905f5280545f5160206168785f395f51905f52939093559251935161ffff60e01b60e09190911b16600160a01b600160e01b0360a09590951b949094166001600160f01b0319909216600160a01b600160f01b0319919091161717919091179055516020015190515f5160206167985f395f51905f52546001600160401b03909216916001600160801b039091169080610ce4575b506001600160401b0382165f9081525f5160206168b85f395f51905f52602052604090206001600160a01b0390600101545f5160206168785f395f51905f52545f5160206168385f395f51905f52546040519390921615939260a08101916001600160401b03831182841017610cd05760029260405281526020810160018060a01b0384168152604082019360018060401b038160a01c16855261ffff606084019160e01c1681526080830194868652610b488960018060401b03165f525f5160206168b85f395f51905f5260205260405f2090565b9351845591516001840180549351925161ffff60e01b60e09190911b16600160a01b600160e01b0360a09490941b939093166001600160a01b039092166001600160f01b031990941693909317171790559151910180546001600160801b039092166001600160801b03199283161790557f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17098054909116919091179055610c70575b506001600160a01b038116610c4a5750610c026116e7565b50610c0e60e051611769565b505b604051614d5d9081611a1b823960805181610eb2015260a05181613b0b015260c05181610604015260e05181818161033e01526106e60152f35b80610c57610c6a926115f1565b50610c6181611667565b5060e0516117c2565b50610c10565b5f5160206167985f395f51905f525468010000000000000000811015610cd057806001610cac92015f5160206167985f395f51905f5255611599565b81546001600160401b0360039290921b91821b191692901b91909117905581610bea565b634e487b7160e01b5f52604160045260245ffd5b5f1981019081116107a457610cf890611599565b905460039190911b1c6001600160401b031680831015610a72579050630fb2737b60e31b5f5260045260245260445ffd5b63fbad885d60e01b5f5260045ffd5b516331734c2160e21b5f5260045260245ffd5b82634724a0fd60e01b5f5260045260245ffd5b905051610d6d610650856119e0565b141589610929565b8551929350602c808202929183041417156107a45760300190816030116107a4571415908a80610920565b151593508c610918565b60b482119450610911565b84634724a0fd60e01b5f5260045260245ffd5b6305f8ded760e21b5f5260045ffd5b8463156f758160e31b5f5260045260b460245260445ffd5b6333fae18560e21b5f5260045260245260445ffd5b631d9feb7d60e11b5f5260045260245260445ffd5b6355bace6f60e01b5f5260045260245260445ffd5b015190505f806103a2565b90601f198316915f5160206168985f395f51905f525f52815f20925f5b818110610ea357509160019391856104d29897969410610e8b575b505050811b015f5160206168985f395f51905f52556103c3565b01515f1960f88460031b161c191690555f8080610e71565b92936020600181928786015181550195019301610e56565b82811115610384575f5160206168985f395f51905f525f52601f830160051c7fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be60208510610f2d575b81601f9101920160051c03905f5b828110610f20575050610384565b5f82820155600101610f12565b5f9150610f04565b634e487b7160e01b5f52602260045260245ffd5b90607f1690610372565b5f80fd5b81516001600160401b038111610f5357602091610f798a84809488010161108a565b81520191019061014a565b608081019081106001600160401b03821117610cd057604052565b604081019081106001600160401b03821117610cd057604052565b601f909101601f19168101906001600160401b03821190821017610cd057604052565b51906001600160a01b0382168203610f5357565b6001600160401b038111610cd057601f01601f191660200190565b92919261101882610ff1565b916110266040519384610fba565b829481845281830111610f53578281602093845f96015e010152565b9080601f83011215610f5357815161105c9260200161100c565b90565b6001600160401b038111610cd05760051b60200190565b51906001600160401b0382168203610f5357565b9190608083820312610f5357604051906110a382610f84565b8351919384926001600160401b038111610f53576060926110c5918301611042565b8352602081015160208401526110dd60408201611076565b60408401520151908160070b8203610f535760600152565b51908115158203610f5357565b519060ff82168203610f5357565b519063ffffffff82168203610f5357565b156111295750565b63ffffffff9063b0369c3160e01b5f5216600452600160245263ffffffff60445260645ffd5b803b156111595750565b630c7f75dd60e01b5f9081526001600160a01b0391909116600452602490fd5b919082039182116107a457565b908151811015611197570160200190565b634e487b7160e01b5f52603260045260245ffd5b919082018092116107a457565b6040516111c481610f9f565b606081525f6020820152508051801590811561138c575b5061137d575f19908051805b611339575b505f19821461132b5760018201908183116107a457600360fc1b6001600160f81b03196112198484611186565b51161480611317575b611307575f5b81518310156112b75761123b8383611186565b5160f81c6030811080156112ad575b61129b57600a82026001600160401b03908116602f1990920160ff169190910181169116811061127f57600190920191611228565b509150506040519061129082610f9f565b81525f602082015290565b50509150506040519061129082610f9f565b506039811161124a565b91509160018111908115916112fb575b506112ec57604051916112d983610f9f565b82526001600160401b0316602082015290565b63384c687760e21b5f5260045ffd5b602b915010155f6112c7565b9150506040519061129082610f9f565b506002611325848351611179565b11611222565b604051915061129082610f9f565b5f1981018181116107a457602d60f81b6001600160f81b031961135c8386611186565b511614611373575080156107a4575f1901806111e7565b92505f90506111ec565b6329120bff60e21b5f5260045ffd5b6040915010155f6111db565b80518210156111975760209160051b010190565b908151518015611593576113bf8161105f565b906113cd6040519283610fba565b808252601f196113dc8261105f565b013660208401375f5b84518051821015611583576113fc82602092611398565b51015160018060401b036040611413848951611398565b510151166040519161142483610f9f565b8252602082019080825260249080611521575b5061144190611837565b91600a60208401536001908184019060226020830153828001918284116107a4576021600a910153600283018092116107a45761148b916114829086611988565b9051908561199e565b82519092906001600160401b03166114be575b505050906114ad600192611869565b6114b78286611398565b52016113e5565b82840190601060208301538284018094116107a457516001600160401b03169260219190910191905b6080841015611504575050916114ad91600194935391925f61149e565b91908084926080607f8397161781530192019060071c92906114e7565b6001905b6080811015611575575080600101600111611561576025019081811161154e5750611441611437565b634e487b7160e01b5f9081526011600452fd5b50634e487b7160e01b5f9081526011600452fd5b60019060071c910190611525565b50509061105c9293505f906118c7565b505f9150565b905f5160206167985f395f51905f5254821015611197575f5160206167985f395f51905f525f52600282901c7f06ff75c3bdc3f03474569b3e952aab1bfb9d41b303b5638523d3e1200dcd1ae2019160031b60181690565b6001600160a01b0381165f9081525f5160206168d85f395f51905f52602052604090205460ff16611662576001600160a01b03165f8181525f5160206168d85f395f51905f5260205260408120805460ff191660011790553391905f5160206167785f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f5160206167d85f395f51905f52602052604090205460ff16611662576001600160a01b03165f8181525f5160206167d85f395f51905f5260205260408120805460ff191660011790553391905f5160206168585f395f51905f52905f5160206167785f395f51905f529080a4600190565b5f80525f5160206167d85f395f51905f526020525f5160206168185f395f51905f525460ff16611765575f8080525f5160206167d85f395f51905f526020525f5160206168185f395f51905f52805460ff1916600117905533905f5160206168585f395f51905f525f5160206167785f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f1461166257805f525f60205260405f205f805260205260405f20600160ff198254161790555f33915f5160206167785f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16611831575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291905f5160206167785f395f51905f529080a4600190565b50505f90565b9061184182610ff1565b61184e6040519182610fba565b828152809261185f601f1991610ff1565b0190602036910137565b8051600181018091116107a45761187f90611837565b805115611197576020918161189b845f940192848453826119b6565b50604051918291518091835e8101838152039060025afa156118bc575f5190565b6040513d5f823e3d90fd5b9291926118d48285611179565b93600185146119785760015b8060011b90868210156118f357506118e0565b93949550509061190d61190684866111ab565b85836118c7565b6119219361191b91956111ab565b906118c7565b90604051611930608082610fba565b6041815260208101906060368337805115611197576020935f936001845360218301526041820152604051918291518091835e8101838152039060025afa156118bc575f5190565b50611984929350611398565b5190565b6020828192010153600181018091116107a45790565b81602091939293010152602081018091116107a45790565b9080519182156119d8576021602084930191015e600101806001116107a45790565b505050600190565b611a16604051916119f2606084610fba565b602283526040366020850137600a6020840153611a10600184611988565b8361199e565b509056fe6080806040526004361015610012575f80fd5b5f3560e01c90816301ffc9a7146112ba5750806310d417e31461121757806323842fb81461117f578063248a9ca3146111555780632f2ff15d1461112657806336568abe146110d75780634b1872be14610d945780636a28f00014610cb35780636edfe3af14610b225780638a8e4c5d14610aea5780638c80fbda14610aae57806391d1485414610a7257806392c19bdc146109cc578063974a74c414610875578063a217fddf1461085b578063a6f031bb1461073f578063d547741f14610709578063db3e1fa4146106cf578063ddba6537146103155763ef913a4b146100f8575f80fd5b34610311575f3660031901126103115760405160208082015261012060408201525f7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1700548060011c90600181168015610307575b6020831081146102f3578261016086015290815f146102ce575060011461024e575b61024a83610236818560ff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170154818116606085015260081c1660808301526001600160401b037f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025481811660a085015260401c1660c083015263ffffffff5f516020614d235f395f51905f525481811660e0850152818160201c1661010085015260ff8160401c16151561012085015260481c1661014083015203601f198101835282611582565b604051918291602083526020830190611432565b0390f35b9190507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005f527fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be915f905b8082106102b2575090915081016101800161023661016e565b9192600181602092546101808588010152019101909291610299565b60ff19166101808086019190915291151560051b84019091019150610236905061016e565b634e487b7160e01b5f52602260045260245ffd5b91607f169161014c565b5f80fd5b346103115761032336611392565b60ff5f516020614d235f395f51905f525460401c166106a7577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610698575b50810190602081830312610311578035906001600160401b0382116103115701610140818303126103115760405160c081018181106001600160401b038211176106845760405281356001600160401b03811161031157820160408185031261031157604051906103ec82611516565b80356001600160401b038111610311578561040891830161171c565b82526020810135906001600160401b0382116103115761042a9186910161171c565b6020820152815261043e83602084016115b7565b916020820192835261045384608083016115b7565b6040830190815261046660e083016115a3565b92606081019384526101008301356001600160401b038111610311578661048e918501611a83565b9560808201968752610120840135956001600160401b03871161031157610628976001600160801b03610602976105c56105f49761059e6105776104dc6105e09960209f6105619e01611a83565b9960a081019a8b526104f087875116612aae565b6104fe885182515190613fca565b61050f8b5160208351015190613fca565b6040519e8f9d8e927f35aff9ac0000000000000000000000000000000000000000000000000000000083850152826024850152519261014060448201526101c484519160406101848201520190611d52565b9101518d820361018319016101a48f0152611d52565b965180516001600160801b031660648d0152602081015160848d01526040015160a48c0152565b5180516001600160801b031660c48b0152602081015160e48b0152604001516101048a0152565b51166101248701525185820360431901610144870152611fc8565b905183820360431901610164850152611fc8565b03601f198101835282611582565b7f0000000000000000000000000000000000000000000000000000000000000000612b39565b506801000000000000000068ff0000000000000000195f516020614d235f395f51905f525416175f516020614d235f395f51905f52557f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1005b634e487b7160e01b5f52604160045260245ffd5b6106a19061291a565b8261037c565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b34610311575f3660031901126103115760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346103115761073d61071a3661136c565b90610738610733825f525f602052600160405f20015490565b61291a565b6129e7565b005b34610311576020366003190112610311576004356001600160401b03811161031157806004019061014060031982360301126103115760ff5f516020614d235f395f51905f525460401c166106a7575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff161561084e575b6107ed6044820183612876565b916107fb6064820185612876565b9390916101048101359060018210156103115760209661084696610823610124840183612876565b969095604051986108348c8b611582565b5f8a52608460a48701960135946139bb565b604051908152f35b6108566128ab565b6107e0565b34610311575f3660031901126103115760206040515f8152f35b34610311576020366003190112610311576004356001600160401b03811161031157806004019061016060031982360301126103115760ff5f516020614d235f395f51905f525460401c166106a7575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff16156109bf575b61014481016109258184612844565b905015610997576109396044830184612876565b90916109486064850186612876565b959094610104810135916001831015610311576020976108469761098096610987610977610124870186612876565b99909886612844565b3691611641565b98608460a48701960135946139bb565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b6109c76128ab565b610916565b34610311576109da36611392565b9060ff5f516020614d235f395f51905f525460401c166106a7575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060209081527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47549092610a5492909160ff1615610a6557612333565b60405190610a61816113e1565b8152f35b610a6d6128ab565b612333565b3461031157610a803661136c565b905f525f6020526001600160a01b0360405f2091165f52602052602060ff60405f2054166040519015158152f35b34610311575f3660031901126103115760207f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170554604051908152f35b3461031157610af836611392565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b34610311575f366003190112610311576020610b3c612cb8565b610b4581612d99565b92019161ffff83511691610b71610b5b84611705565b93610b696040519586611582565b808552611705565b602084019190601f1901368337610b8c61ffff865116612284565b9261ffff86511691610ba0610b5b84611705565b602084019690601f19013688375f5b61ffff89511661ffff821690811015610c145761ffff91886001600160401b03610c0984808f610bfa82610c02928e8c8f8f9d610bef60019f82906122b6565b525192511691612e4a565b9290966122b6565b528a6122b6565b911690520116610baf565b87838a888a604051948594606086019060608752518091526080860192905f5b818110610c945750505081610c51918660209403848801526113ff565b91848303604086015251918281520191905f5b818110610c72575050500390f35b82516001600160401b0316845285945060209384019390920191600101610c64565b825163ffffffff16855288975060209485019490920191600101610c34565b34610311575f36600319011261031157335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff1615610d7d575f516020614d235f395f51905f525460ff8160401c1615610d555768ff000000000000000019165f516020614d235f395f51905f52557f8e0da210a2ffecc744373b4860cd9b8dc471810f5fb61b7d2c3234ccd4aae30d5f80a1005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b3461031157610da236611392565b60ff5f516020614d235f395f51905f525460401c166106a7575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff16156110ca575b810190602081830312610311578035906001600160401b038211610311576105f4610eb0610e45602095610ee495610ed69501611cd5565b610e5b6001600160801b03604083015116612aae565b610e9560608201518783015190606089610e73612cb8565b93610e8a6001600160801b036080870151166141bb565b5101510151906142df565b60405192839163319d103360e21b88840152602483016120e2565b7f0000000000000000000000000000000000000000000000000000000000000000612b39565b8280825183010191016121e2565b610eed81612bc4565b90610ef7826113e1565b8161105a57826001600160401b0381838160607f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb3196019184828451015116610f697f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025491878360401c16808211612240565b83516fffffffffffffffffffffffffffffffff196fffffffffffffffff0000000000000000858984511693015160401b16921617177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1702550151604051610fef8482018093604080916001600160801b038151168452602081015160208501520151910152565b60608152610ffe608082611582565b519020611041848484510151166001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b5551015116604051908152a160405190610a61816113e1565b50611064816113e1565b60018103610a54576801000000000000000068ff0000000000000000195f516020614d235f395f51905f525416175f516020614d235f395f51905f52557f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1610a54565b6110d26128ab565b610e0d565b34610311576110e53661136c565b336001600160a01b038216036110fe5761073d916129e7565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b346103115761073d6111373661136c565b90611150610733825f525f602052600160405f20015490565b61295a565b346103115760203660031901126103115760206108466004355f525f602052600160405f20015490565b34610311576020366003190112610311576004356001600160401b0381168103610311576111dd906001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b5480156111ef57602090604051908152f35b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b34610311575f366003190112610311577f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708545f1981019081116112a6576001600160401b03611267604092611463565b90549060031b1c166001600160801b037f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1709541682519182526020820152f35b634e487b7160e01b5f52601160045260245ffd5b3461031157602036600319011261031157600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361031157817f7965db0b000000000000000000000000000000000000000000000000000000006020931490811561132e575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483611327565b35906001600160401b038216820361031157565b604090600319011261031157600435906024356001600160a01b03811681036103115790565b906020600319830112610311576004356001600160401b0381116103115782602382011215610311578060040135926001600160401b0384116103115760248483010111610311576024019190565b600311156113eb57565b634e487b7160e01b5f52602160045260245ffd5b90602080835192838152019201905f5b81811061141c5750505090565b825184526020938401939092019160010161140f565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b919082039182116112a657565b907f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708548210156114e1577f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17085f52600282901c7f06ff75c3bdc3f03474569b3e952aab1bfb9d41b303b5638523d3e1200dcd1ae2019160031b60181690565b634e487b7160e01b5f52603260045260245ffd5b91909180548310156114e1575f52601860205f208360021c019260031b1690565b604081019081106001600160401b0382111761068457604052565b606081019081106001600160401b0382111761068457604052565b608081019081106001600160401b0382111761068457604052565b60a081019081106001600160401b0382111761068457604052565b90601f801991011681019081106001600160401b0382111761068457604052565b35906001600160801b038216820361031157565b9190826060910312610311576040516115cf81611531565b60408082946115dd816115a3565b8452602081013560208501520135910152565b91908260409103126103115760405161160881611516565b602061162181839561161981611358565b855201611358565b910152565b6001600160401b03811161068457601f01601f191660200190565b92919261164d82611626565b9161165b6040519384611582565b829481845281830111610311578281602093845f960137010152565b3590811515820361031157565b359063ffffffff8216820361031157565b809291039160608312610311576040516116ae81611516565b6040819483358352601f1901126103115760209060408051936116d085611516565b6116db848201611684565b85520135828401520152565b9080601f830112156103115781602061170293359101611641565b90565b6001600160401b0381116106845760051b60200190565b9190606083820312610311576040519061173582611516565b819380356001600160401b038111610311578101604081840312610311576040519061176082611516565b80356001600160401b0381116103115781016102c081860312610311576040519061026082018281106001600160401b03821117610684576040526117a586826115f0565b825260408101356001600160401b03811161031157810186601f8201121561031157868160206117d793359101611641565b60208301526117e860608201611358565b60408301526117f9608082016115a3565b606083015261180a60a08201611677565b608083015261181c8660c08301611695565b60a083015261182e6101208201611677565b60c083015261014081013560e083015261184b6101608201611677565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a083015261189a6102208201611677565b6101c08301526102408101356101e08301526118b96102608201611677565b6102008301526102808101356102208301526102a0810135906001600160401b038211610311576118ec918791016116e7565b61024082015282526020810135906001600160401b038211610311570160c08185031261031157604051906119208261154c565b61192981611358565b825261193760208201611684565b60208301526119498560408301611695565b604083015260a0810135906001600160401b038211610311570184601f8201121561031157803561197981611705565b916119876040519384611582565b81835260208084019260061b8201019087821161031157602001915b8183106119c95750505092602094928261162195606088950152838201528652016115f0565b60408389031261031157604051906119e082611516565b833560048110156103115782526020840135906bffffffffffffffffffffffff198216820361031157826020928360409501528152019201916119a3565b9080601f83011215610311578135611a3581611705565b92611a436040519485611582565b81845260208085019260051b82010192831161031157602001905b828210611a6b5750505090565b60208091611a7884611684565b815201910190611a5e565b9190610260838203126103115760405161010081018181106001600160401b0382111761068457604052809382601f8201121561031157604051611ac961010082611582565b806101008301918583116103115783905b838210611cc557505083528361011f8301121561031157604051611aff608082611582565b8061018084019286841161031157905b838210611cb557505060208401528361019f83011215610311576040805190611b389082611582565b806101c084019286841161031157905b838210611ca557505060408401523561ffff811681036103115760608301526101e08101356001600160401b0381116103115783611b87918301611a1e565b60808301526102008101356001600160401b0381116103115783611bac918301611a1e565b60a08301526102208101356001600160401b03811161031157810183601f8201121561031157803590611bde82611705565b91611bec6040519384611582565b80835260208084019160051b8301019186831161031157602001905b828210611c955750505060c0830152610240810135906001600160401b03821161031157019180601f84011215610311578235611c4481611705565b93611c526040519586611582565b81855260208086019260051b82010192831161031157602001905b828210611c7d5750505060e00152565b60208091611c8a84611677565b815201910190611c6d565b8135815260209182019101611c08565b8135815260209182019101611b48565b8135815260209182019101611b0f565b8135815260209182019101611ada565b919060c0838203126103115760405190611cee8261154c565b8193611cfa82826115b7565b835260608101356001600160401b0381116103115782611d1b91830161171c565b6020840152611d2c608082016115a3565b604084015260a0810135916001600160401b038311610311576060926116219201611a83565b90815191606082526020611eaf845160406060860152611d8b60a0860182516001600160401b0360208092828151168552015116910152565b610240611da8848301516102c060e0890152610360880190611432565b60408301516001600160401b031661010088015260608301516001600160801b03166101208801526080830151151561014088015260a08301518051610160890152602090810151805163ffffffff166101808a015201516101a08801529160c081015115156101c088015260e08101516101e08801526101008101511515610200880152610120810151610220880152610140810151828801526101608101516102608801526101808101516102808801526101a08101516102a08801526101c081015115156102c08801526101e08101516102e088015261020081015115156103008801526102208101516103208801520151609f1986830301610340870152611432565b930151605f19838503016080840152602060e0606060c08701936001600160401b03815116885263ffffffff848201511684890152611f10604082015160408a019060208060409280518552015163ffffffff815116828501520151910152565b01519560c060a08201528651809452019401905f5b818110611f535750505060209081015180516001600160401b03908116838501529101511660409091015290565b9091948551805160048110156113eb5782526020908101516bffffffffffffffffffffffff191681830152604090910195019190600101611f25565b90602080835192838152019201905f5b818110611fac5750505090565b825163ffffffff16845260209384019390920191600101611f9f565b80515f835b600882106120cc5750505060208101515f61010084015b600482106120b65750505060408101515f61018084015b600282106120a05750505061ffff6060820151166101c083015260e061205f61204c61203960808501516102606101e0880152610260870190611f8f565b60a0850151868203610200880152611f8f565b60c08401518582036102208701526113ff565b91015191610240818303910152602080835192838152019201905f5b8181106120885750505090565b8251151584526020938401939092019160010161207b565b6020806001928551815201930191019091611ffb565b6020806001928551815201930191019091611fe4565b6020806001928551815201930191019091611fcd565b906117029160208152612116602082018351604080916001600160801b038151168452602081015160208501520151910152565b6060612131602084015160c0608085015260e0840190611d52565b926001600160801b0360408201511660a084015201519060c0601f1982850301910152611fc8565b91908260609103126103115760405161217181611531565b809280516001600160801b03811681036103115760409182918452602081015160208501520151910152565b51906001600160401b038216820361031157565b9190826040910312610311576040516121c981611516565b60206116218183956121da8161219d565b85520161219d565b90610140828203126103115761223890610100604051936122028561154c565b61220c8382612159565b855261221b8360608301612159565b602086015261222d8360c083016121b1565b6040860152016121b1565b606082015290565b15612249575050565b906001600160401b0380927f7d939bd8000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b9061228e82611705565b61229b6040519182611582565b82815280926122ac601f1991611705565b0190602036910137565b80518210156114e15760209160051b010190565b919060808382031261031157604051906122e38261154c565b819380356001600160401b038111610311576060926123039183016116e7565b83526020810135602084015261231b60408201611358565b60408401520135908160070b82036103115760600152565b90810190602081830312610311578035906001600160401b0382116103115701604081830312610311576040519061236a82611516565b80356001600160401b0381116103115783612386918301611cd5565b82526020810135906001600160401b038211610311570160808184031261031157604051926123b48461154c565b81356001600160401b03811161031157820181601f820112156103115780356123dc81611705565b916123ea6040519384611582565b81835260208084019260051b820101918483116103115760208201905b8382106128175750505050845261242060208301611677565b60208501526040820135916001600160401b0383116103115761244a6060926124559483016122ca565b604086015201611358565b606083015260208101918252516124786001600160801b03604083015116612aae565b6124cc6124bd6105f4610eb06060850151946124a1602082019687519060606020610e73612cb8565b60405192839163319d103360e21b6020840152602483016120e2565b602080825183010191016121e2565b916124d683612bc4565b926124e0846113e1565b600184146127b8576101606124f58351612f0c565b935151510151808403612788575060207f1411889dcdfcf9af1fcf4ff4521cab4580a1fd842401ead543815cb1148d643893926001600160401b039261253a876113e1565b866126ff57828161258b60606126a094019461258588858851015116897f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025460401c16808211612240565b516130d3565b8351868151166fffffffffffffffffffffffffffffffff196fffffffffffffffff0000000000000000857f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025494015160401b16921617177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1702550180516040516126348682018093604080916001600160801b038151168452602081015160208501520151910152565b60608152612643608082611582565b519020612686868686510151166001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b556001600160801b038585855101511691515116906136b0565b7f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb318284818451015116604051908152a1510151166126f960405192839283602090939291936001600160401b0360408201951681520152565b0390a190565b80612744606061275f93019361258587878751015116887f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025460401c16808214612240565b6001600160801b0384868186510151169201515116906136b0565b510151166126f960405192839283602090939291936001600160401b0360408201951681520152565b83907f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b5050506801000000000000000068ff0000000000000000195f516020614d235f395f51905f525416175f516020614d235f395f51905f52557f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a190565b81356001600160401b03811161031157602091612839888480948801016122ca565b815201910190612407565b903590601e198136030182121561031157018035906001600160401b0382116103115760200191813603831361031157565b903590601e198136030182121561031157018035906001600160401b03821161031157602001918160051b3603831361031157565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff16156128e357565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260405f206001600160a01b0333165f5260205260ff60405f205416156129445750565b63e2517d3f60e01b5f523360045260245260445ffd5b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f205416155f146129e157805f525f60205260405f206001600160a01b0383165f5260205260405f20600160ff198254161790556001600160a01b03339216907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f2054165f146129e157805f525f60205260405f206001600160a01b0383165f5260205260405f2060ff1981541690556001600160a01b03339216907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b919082018092116112a657565b15612a80575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160801b03633b9aca0091160463ffffffff5f516020614d235f395f51905f525460481c16612aed8242612ae58442612a6a565b821115612a77565b428210612af8575050565b612b028242611456565b11612b0a5750565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b803b15612b8f57815f929160208493519201905af43d15612b87573d90612b5f82611626565b91612b6d6040519384611582565b82523d5f602084013e5b15612b7f5790565b602081519101fd5b606090612b77565b6001600160a01b03907f0c7f75dd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b6001600160401b03602060608301510151165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f205480155f14612c0f5750505f90565b60208201908151604051612c44602082018093604080916001600160801b038151168452602081015160208501520151910152565b60608152612c53608082611582565b5190201491821592612c71575b505015612c6c57600190565b600290565b6001600160801b0391925081905151169151511611155f80612c60565b60405190612c9b82611567565b5f6080838281528260208201528260408201528260608201520152565b612cc0612c8e565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1705547f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17065461ffff6001600160801b037f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170954169160405193612d4185611567565b84526001600160a01b03811660208501526001600160401b038160a01c16604085015260e01c166060830152608082015290565b60405190612d828261154c565b5f6060838281528260208201528260408201520152565b612da1612d75565b50602081016001600160a01b0381511615612e36576001600160a01b03905116612dd261ffff6060840151166146cf565b813b6001811115612e23575f1981019081116112a6578111612e105790816001612dfe612e0d946146f4565b9260208401903c80925161471c565b91565b50633610565160e21b5f5260045260245ffd5b82633610565160e21b5f5260045260245ffd5b5051633a517eed60e21b5f5260045260245ffd5b939263ffffffff169161ffff16821015612edc57602c820293828504602c14831517156112a6578460300191826030116112a657850192605084015160e01c03612eb157506034840181116112a657603c605483015160c01c9401106112a657605c015190565b7f4724a0fd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b50827ff81f5140000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b9081515180156130cd57612f1f81612284565b905f5b845180518210156130bd57612f39826020926122b6565b5101516001600160401b036040612f518489516122b6565b5101511660405191612f6283611516565b825260208201908082526024908061305d575b50612f7f906146f4565b91600a60208401536001908184019060226020830153828001918284116112a6576021600a910153600283018092116112a657612fc991612fc09086614c26565b90519085614c3c565b916001600160401b03815116612ffa575b50505090612fe96001926147f2565b612ff382866122b6565b5201612f22565b82840190601060208301538284018094116112a657516001600160401b03169260219190910191905b608084101561304057505091612fe991600194935391925f612fda565b91908084926080607f8397161781530192019060071c9290613023565b6001905b60808110156130af57508060010160011161309c576025019081811161308a5750612f7f612f75565b634e487b7160e01b5f5260116004525ffd5b50634e487b7160e01b5f5260116004525ffd5b60019060071c910190613061565b5050906117029293505f90614845565b505f9150565b906130dd82612f0c565b91519081519182156136575760b4831161367f575f91602c8402848104602c036112a657603001806030116112a657613115906146f4565b9060208201946056865361564160218401536256414c60228401536356414c3460238401538060081c602c840153602d83015361315961315487614ce8565b6147f2565b602e8301525f604e8301535f604f8301535f9560305b84518810156133195761318288866122b6565b519560408701906001600160401b03825116156132ed5760205f9801975b8a811061329f57506131bd906001600160401b0383511690612a6a565b968286019160ff8b60181c16602084015361ffff8b60101c16602184015362ffffff8b60081c16602284015363ffffffff8b1660238401536004840184116112a6576001600160401b03905160ff8160381c16602485015361ffff8160301c16602585015362ffffff8160281c16602685015363ffffffff8160201c16602785015364ffffffffff8160181c16602885015365ffffffffffff8160101c16602985015366ffffffffffffff8160081c16602a85015316602b830153600c830183116112a657602c9051910152602c81018091116112a65760019097019661316f565b60206132ab828a6122b6565b5101518951146132bd576001016131a0565b8a907faefa56c2000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b897fafef9f15000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b50949391925094506001600160401b038111613657578060ff6001600160401b039260381c16602484015361ffff8160301c16602584015362ffffff8160281c16602684015363ffffffff8160201c16602784015364ffffffffff8160181c16602884015365ffffffffffff8160101c16602984015366ffffffffffffff8160081c16602a84015316602b8201536133b1818461471c565b91604051916133e260218460208101945f8652845180918484015e81015f838201520301601f198101855284611582565b61600083511161362b5750906001600160a01b039161349a602a7fffff000000000000000000000000000000000000000000000000000000000000845160f01b169360405193849160208301967f6100000000000000000000000000000000000000000000000000000000000000885260218401527f80600a3d393df30000000000000000000000000000000000000000000000000060238401525180918484015e81015f838201520301601f198101835282611582565b51905ff0169182156136035760209273ffffffffffffffffffffffffffffffffffffffff197f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17065416177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706557f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055580517fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b807f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706549360a01b16169116177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655015161ffff60e01b1961ffff60e01b807f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706549360e01b16169116177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b827fab7bac08000000000000000000000000000000000000000000000000000000005f5260045260b460245260445ffd5b907f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708548061397b575b506001600160a01b03600161371e846001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f2090565b01541615907f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170554906001600160801b0360027f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170654936040519061378082611567565b8152602081016001600160a01b038616815260408201956001600160401b038160a01c16875261ffff606084019160e01c16815284608084019616968787526137f98a6001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f2090565b935184556001600160a01b036001850193511673ffffffffffffffffffffffffffffffffffffffff19845416178355517fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b8085549360a01b161691161782555161ffff60e01b1961ffff60e01b8084549360e01b16169116179055019151166fffffffffffffffffffffffffffffffff198254161790556fffffffffffffffffffffffffffffffff197f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17095416177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1709556138fb5750565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170854680100000000000000008110156106845780600161395d92017f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170855611463565b6001600160401b0380839493549260031b9316831b921b1916179055565b5f1981019081116112a657826001600160401b0361399b6139b593611463565b90549060031b1c16806001600160401b0383161015612240565b5f6136d9565b999a9890939294969a979195975f9a60018d10156113eb578c15613a085760248c60ff8f7f112d89cc00000000000000000000000000000000000000000000000000000000835216600452fd5b909192939495968098999a9c50151580613fbe575b15613f8757602001356001600160401b038116810361031157613a40368b6115b7565b90613abc6040516020810190613a738286604080916001600160801b038151168452602081015160208501520151910152565b60608152613a82608082611582565b519020916001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b5480156111ef57808203613f595750506020810151808a03613f2957506001600160801b03613aec9151166141bb565b5f905f5b888110613e40575b505015613dff5750506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001693843b156103115794929085926040519687957f56455e790000000000000000000000000000000000000000000000000000000087526064870190600488015260606024880152526084850160848560051b87010194825f90603e1981360301935b838310613d8e57505050505050600319848403016044850152808352602083019260208260051b82010193835f92601e1982360301905b858510613c2857505050505050509181805f9403915afa8015613c1d57613c08575b5035906001600160801b038216809203613c055750633b9aca00900490565b80fd5b613c159192505f90611582565b5f905f613be6565b6040513d5f823e3d90fd5b9193959750919395601f1982820301855287358381121561031157840190613c54602082019280614a43565b8091936020845252604082019060408160051b8401019380935f915b838310613c9757505050505050602080600192990195019501929091899796949592613bc4565b909192939495603f19838203018652613cb08783614a8b565b80356002811015610311578252613cde613ccd6020830183614a77565b606060208501526060840190614a9f565b906040810135609e1982360301811215610311576001936020938493613d809301916040818303910152613d72613d54613d29613d1b85806149b3565b60a0865260a0860191614993565b613d34878601611677565b151587850152613d476040860186614a77565b8482036040860152614a9f565b92613d6160608201611677565b151560608401526080810190614a77565b906080818403910152614a9f565b980196019493019190613c70565b9193959698509193966083198b8203018252873586811215610311576020613dec6001938683940190613ddf613dd5613dc78480614a43565b6040855260408501916149e4565b92858101906149b3565b9185818503910152614993565b9901920193019093918a98969593613b8d565b613e3c6040519283927ffef760c70000000000000000000000000000000000000000000000000000000084526020600485015260248401916149e4565b0390fd5b9091613e74613e63613e5c613e56858d8c614906565b80612876565b3691614928565b613e6e368789614928565b90614c7e565b15613f1f5750610980613e8b613e95928a89614906565b6020810190612844565b805182518082149182613f09575b505015613eb557505060015f80613af8565b90613e3c613ef7926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190611432565b83810360031901602485015290611432565b9091506020830120906020840120145f80613ea3565b9190600101613af0565b89907f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b877fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b5061ffff881115613a1d565b9190916001600160401b03602080850151015116613fe6612c8e565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708548015614123575f905f1981019081116112a6575b80821061413657506140576001600160401b03917f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17086114f5565b90549060031b1c16908111614123575f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f20906040519161409e83611567565b80548084526001600160801b03600260018401549361ffff6001600160a01b038616958660208a01526001600160401b038160a01c1660408a015260e01c1660608801520154169160808501928352156141115750602061410f9495610e8a6001600160801b0360609451166141bb565b565b633a517eed60e21b5f5260045260245ffd5b633a517eed60e21b5f525f60045260245ffd5b906141418282612a6a565b600181018091116112a65760011c90836001600160401b03614183847f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17086114f5565b905460039190911b1c161161419957509061401d565b91505f1981019081111561401d57634e487b7160e01b5f52601160045260245ffd5b6001600160801b03633b9aca009116045f516020614d235f395f51905f52546141f28242612ae563ffffffff8560481c1642612a6a565b5f91428110614247575b5063ffffffff16906001600160801b031681811015614219575050565b7fdb020436000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160801b038092935061426263ffffffff9242611456565b169291506141fc565b156142735750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156142ad5750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b9092919260808201515193606083019461ffff86511680911490816146c0575b816146b1575b816146a2575b50156136575761431a81612d99565b60408301515f96879593949293869283805b61ffff8d51168c101561449e576143478c60e08d01516122b6565b511561448e57908a9493929163ffffffff6143668e60808901516122b6565b5116809b614476575b505063ffffffff6143878d60a060019d9801516122b6565b51169561439e8761ffff60208d01511681106142a5565b61010087101561444a578b908b881b906143bb898383161561426b565b17968d8b828c8c51926020015161ffff16906143d693612e4a565b91909360c00151906143e7916122b6565b510361441f57506001600160401b038091169116016001600160401b0381116112a657600161ffff919c5b019b91929394905061432c565b7fc3271bcd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b867ff6190fea000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b63ffffffff6144879216811161426b565b5f8a61436f565b939291909a600161ffff91614412565b5095509650969850965050506001600160401b0391501660038102818104600314821517156112a6576801fffffffffffffffe8360011b16906001600160401b03841682046002146001600160401b0385161517156112a6571115614669575050608083015160e060c085015194015181518151809114908161465e575b5015613657575f5b82518110156146565761453781836122b6565b511561464e5763ffffffff61454c82856122b6565b511661455882886122b6565b5161456682875181106142a5565b61457082876122b6565b515160048110156113eb57600119016146225760205f6bffffffffffffffffffffffff198261459f868b6122b6565b51015116926040518381019182528381526145bb604082611582565b604051918291518091835e8101838152039060025afa15613c1d575f516bffffffffffffffffffffffff1916036145f757506001905b01614524565b7fdc59b237000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b507f5b5c80eb000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b6001906145f1565b505050509050565b90508551145f61451c565b6001600160401b0392507f6e3083c3000000000000000000000000000000000000000000000000000000005f526004521660245260445ffd5b905060e084015151145f61430b565b60c08501515181149150614305565b60a085015151811491506142ff565b61ffff16602c810290808204602c14901517156112a657603001806030116112a65790565b906146fe82611626565b61470b6040519182611582565b82815280926122ac601f1991611626565b919091614727612d75565b506030835110612eb1576356414c34602084015160e01c03612eb157602483015160c01c602c84015160f01c602e85015190604e86015160f01c926040519061476f8261154c565b8152816020820152604081019283528360608201529581159384156147e7575b84156147dd575b5083156147c6575b505081156147af575b50612eb15750565b9050516147be61315483614ce8565b14155f6147a7565b519192506147d3906146cf565b1415905f8061479e565b151593505f614796565b60b48311945061478f565b8051600181018091116112a657614808906146f4565b8051156114e15760209181614824845f94019284845382614c54565b50604051918291518091835e8101838152039060025afa15613c1d575f5190565b9291926148528285611456565b93600185146148f65760015b8060011b9086821015614871575061485e565b93949550509061488b6148848486612a6a565b8583614845565b61489f936148999195612a6a565b90614845565b906040516148ae608082611582565b60418152602081019060603683378051156114e1576020935f936001845360218301526041820152604051918291518091835e8101838152039060025afa15613c1d575f5190565b506149029293506122b6565b5190565b91908110156114e15760051b81013590603e1981360301821215610311570190565b92919061493481611705565b936149426040519586611582565b602085838152019160051b8101918383116103115781905b838210614968575050505050565b81356001600160401b0381116103115760209161498887849387016116e7565b81520191019061495a565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156103115701602081359101916001600160401b03821161031157813603831361031157565b90602083828152019260208260051b82010193835f925b848410614a0b5750505050505090565b909192939495602080614a33600193601f19868203018852614a2d8b886149b3565b90614993565b98019401940192949391906149fb565b9035601e19823603018112156103115701602081359101916001600160401b038211610311578160051b3603831361031157565b9035607e1982360301811215610311570190565b9035605e1982360301811215610311570190565b614ad8614abd614aaf83806149b3565b608086526080860191614993565b614aca60208401846149b3565b908583036020870152614993565b614ae56040830183614a77565b9083810360408501528135600381101561031157614b02816113e1565b81526020820135600381101561031157614b1b816113e1565b60208201526040820135906003821015610311576080614b57614b729484614b48614b67969998996113e1565b604085015260608101906149b3565b9190928160608201520191614993565b926060810190614a43565b90916060818503910152808352602083019060208160051b85010193835f915b838310614ba25750505050505090565b909192939495601f19828203018652614bbb8784614a8b565b803591600383101561031157614c18602092839285614bdb6001976113e1565b8152614c0a614bff614bef868501856149b3565b6060888601526060850191614993565b9260408101906149b3565b916040818503910152614993565b980196019493019190614b92565b6020828192010153600181018091116112a65790565b81602091939293010152602081018091116112a65790565b908051918215614c76576021602084930191015e600101806001116112a65790565b505050600190565b9081518151036129e1575f5b8251811015614c7657614c9d81846122b6565b5151614ca982846122b6565b515103614ce157614cba81846122b6565b5160208151910120614ccc82846122b6565b516020815191012003614ce157600101614c8a565b5050505f90565b614d1e60405191614cfa606084611582565b602283526040366020850137600a6020840153614d18600184614c26565b83614c3c565b509056fe5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1703a26c6578706572696d656e74616cf564736f6c634300082500182f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17085e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17014c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e05e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17036e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1707ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb55e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1702",
}

// ContractSpectreClientABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractSpectreClientMetaData.ABI instead.
var ContractSpectreClientABI = ContractSpectreClientMetaData.ABI

// ContractSpectreClientBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractSpectreClientMetaData.Bin instead.
var ContractSpectreClientBin = ContractSpectreClientMetaData.Bin

// DeployContractSpectreClient deploys a new Ethereum contract, binding an instance of ContractSpectreClient to it.
func DeployContractSpectreClient(auth *bind.TransactOpts, backend bind.ContractBackend, updateClientModule common.Address, membershipModule common.Address, misbehaviourModule common.Address, clientState_ []byte, consensusState_ ISpectreTendermintMsgsConsensusState, initialPinnedValidatorSet ISpectreTendermintMsgsValidatorSet, roleManager common.Address) (common.Address, *types.Transaction, *ContractSpectreClient, error) {
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

// GetConsensusStateHash is a free data retrieval call binding the contract method 0x23842fb8.
//
// Solidity: function getConsensusStateHash(uint64 revisionHeight) view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientCaller) GetConsensusStateHash(opts *bind.CallOpts, revisionHeight uint64) ([32]byte, error) {
	var out []interface{}
	err := _ContractSpectreClient.contract.Call(opts, &out, "getConsensusStateHash", revisionHeight)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetConsensusStateHash is a free data retrieval call binding the contract method 0x23842fb8.
//
// Solidity: function getConsensusStateHash(uint64 revisionHeight) view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientSession) GetConsensusStateHash(revisionHeight uint64) ([32]byte, error) {
	return _ContractSpectreClient.Contract.GetConsensusStateHash(&_ContractSpectreClient.CallOpts, revisionHeight)
}

// GetConsensusStateHash is a free data retrieval call binding the contract method 0x23842fb8.
//
// Solidity: function getConsensusStateHash(uint64 revisionHeight) view returns(bytes32)
func (_ContractSpectreClient *ContractSpectreClientCallerSession) GetConsensusStateHash(revisionHeight uint64) ([32]byte, error) {
	return _ContractSpectreClient.Contract.GetConsensusStateHash(&_ContractSpectreClient.CallOpts, revisionHeight)
}

// GetPinnedSnapshot is a free data retrieval call binding the contract method 0x10d417e3.
//
// Solidity: function getPinnedSnapshot() view returns(uint64 height, uint128 timestamp)
func (_ContractSpectreClient *ContractSpectreClientCaller) GetPinnedSnapshot(opts *bind.CallOpts) (struct {
	Height    uint64
	Timestamp *big.Int
}, error) {
	var out []interface{}
	err := _ContractSpectreClient.contract.Call(opts, &out, "getPinnedSnapshot")

	outstruct := new(struct {
		Height    uint64
		Timestamp *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Height = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPinnedSnapshot is a free data retrieval call binding the contract method 0x10d417e3.
//
// Solidity: function getPinnedSnapshot() view returns(uint64 height, uint128 timestamp)
func (_ContractSpectreClient *ContractSpectreClientSession) GetPinnedSnapshot() (struct {
	Height    uint64
	Timestamp *big.Int
}, error) {
	return _ContractSpectreClient.Contract.GetPinnedSnapshot(&_ContractSpectreClient.CallOpts)
}

// GetPinnedSnapshot is a free data retrieval call binding the contract method 0x10d417e3.
//
// Solidity: function getPinnedSnapshot() view returns(uint64 height, uint128 timestamp)
func (_ContractSpectreClient *ContractSpectreClientCallerSession) GetPinnedSnapshot() (struct {
	Height    uint64
	Timestamp *big.Int
}, error) {
	return _ContractSpectreClient.Contract.GetPinnedSnapshot(&_ContractSpectreClient.CallOpts)
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
func (_ContractSpectreClient *ContractSpectreClientTransactor) VerifyMembership(opts *bind.TransactOpts, msg_ ILightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractSpectreClient.contract.Transact(opts, "verifyMembership", msg_)
}

// VerifyMembership is a paid mutator transaction binding the contract method 0x974a74c4.
//
// Solidity: function verifyMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[],bytes) msg_) returns(uint256)
func (_ContractSpectreClient *ContractSpectreClientSession) VerifyMembership(msg_ ILightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.VerifyMembership(&_ContractSpectreClient.TransactOpts, msg_)
}

// VerifyMembership is a paid mutator transaction binding the contract method 0x974a74c4.
//
// Solidity: function verifyMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[],bytes) msg_) returns(uint256)
func (_ContractSpectreClient *ContractSpectreClientTransactorSession) VerifyMembership(msg_ ILightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.VerifyMembership(&_ContractSpectreClient.TransactOpts, msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0xa6f031bb.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[]) msg_) returns(uint256)
func (_ContractSpectreClient *ContractSpectreClientTransactor) VerifyNonMembership(opts *bind.TransactOpts, msg_ ILightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
	return _ContractSpectreClient.contract.Transact(opts, "verifyNonMembership", msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0xa6f031bb.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[]) msg_) returns(uint256)
func (_ContractSpectreClient *ContractSpectreClientSession) VerifyNonMembership(msg_ ILightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
	return _ContractSpectreClient.Contract.VerifyNonMembership(&_ContractSpectreClient.TransactOpts, msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0xa6f031bb.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes)[],((uint8,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8,bytes[]) msg_) returns(uint256)
func (_ContractSpectreClient *ContractSpectreClientTransactorSession) VerifyNonMembership(msg_ ILightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
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
