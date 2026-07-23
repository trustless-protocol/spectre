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

// IICS07TendermintMsgsConsensusState is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsConsensusState struct {
	Timestamp          *big.Int
	Root               [32]byte
	NextValidatorsHash [32]byte
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

// ContractSpectreClientMetaData contains all meta data concerning the ContractSpectreClient contract.
var ContractSpectreClientMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"updateClientModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membershipModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviourModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"clientState_\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"initialPinnedValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPinnedValidatorSet\",\"inputs\":[],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPinnedValidatorsHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateApplicationState\",\"inputs\":[{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateConsensusState\",\"inputs\":[{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ClientFrozen\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClientUnfrozen\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClientUpdated\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConsensusStateUpdated\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientNotFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DirectCallNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NonMonotonicHeightUpdate\",\"inputs\":[{\"name\":\"latestHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"updateHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofSignerCommitSigMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorCountExceedsLimit\",\"inputs\":[{\"name\":\"validatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxValidatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x6101008060405234610d885761626b803803809161001d8285610def565b8339810160e082820312610d885761003482610e12565b61004060208401610e12565b61004c60408501610e12565b60608501519092906001600160401b038111610d88578461006e918701610e77565b9060808601519460a087015160018060401b038111610d88578701608081830312610d8857604051976100a089610db9565b81516001600160401b038111610d8857820183601f82011215610d885780516100c881610e94565b916100d66040519384610def565b81835260208084019260051b82010191868311610d885760208201905b838210610d8c5750505050895261010c60208301610f2a565b60208a015260408201516001600160401b038111610d885760608361013b610151966101469460c09701610ebf565b60408d015201610eab565b60608a015201610e12565b937f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a6123360e05282518301936020850193602081870312610d88576020810151906001600160401b038211610d88570194859003601f198101939061012013610d88576040519460e086016001600160401b03811187821017610b275760405260208701516001600160401b038111610d88576020908801019080601f83011215610d8857815161020292602001610e41565b855260408412610d8857604080519461021a86610dd4565b610225828901610f37565b865261023360608901610f37565b602087015260208701958652603f190112610d88576040519761025589610dd4565b61026160808801610eab565b895261026f60a08801610eab565b60208a01526040860198895261028760c08801610f45565b956060810196875261029b60e08901610f45565b97608082019889526102c46101206102b66101008401610f2a565b9260a0850193845201610f45565b60c0830190815282518051919891906001600160401b038211610b27575f5160206161eb5f395f51905f5254600181811c91168015610d7e575b6020821014610d6a57601f8111610cfb575b50602090601f8311600114610c795761045e949392915f9183610c6e575b50508160011b915f199060031b1c1916175f5160206161eb5f395f51905f52555b5160ff81511661ff0060205f51602061610b5f395f51905f525493015160081b169161ffff191617175f51602061610b5f395f51905f52558b5160018060401b038151165f51602061624b5f395f51905f525491602068010000000000000000600160801b0391015160401b169160018060801b03191617175f51602061624b5f395f51905f525563ffffffff895116905f51602061614b5f395f51905f52549068ff000000000000000067ffffffff000000008d5160201b169151151560401b916cffffffff0000000000000000008c5160481b1693856cffffffff000000000000000000199160018060481b031916171617911617175f51602061614b5f395f51905f5255801515610f56565b61047163ffffffff875116801515610f56565b516001600160401b039060209061048790610fc3565b01518a51516001600160401b03169116818103610c5957505060018060401b0360208a510151165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f205560805260a05260c05263ffffffff80835116915116019063ffffffff82116106c85763ffffffff80809451169384925116921611610c4457505061051c83611455565b9251918251928315610c1d5760b48411610c2c575f92602c8502858104602c036106c857603001806030116106c85761055490611642565b9260208401966056885361564160218601536256414c60228601536356414c346023860153602c8501968060081c8853602d86015361059a610595826117eb565b611674565b96602e8601978852604e8601915f83535f604f8801535f9960305b87518c10156106dc578b6105c9818a611441565b51604081018051909c6105e5916001600160401b031690610fb6565b9b63ffffffff848d019360ff8160181c16602086015361ffff8160101c16602186015362ffffff8160081c1660228601531660238401536004840184116106c8575160ff8160381c16602484015361ffff8160301c16602584015362ffffff8160281c16602684015363ffffffff8160201c16602784015364ffffffffff8160181c16602884015365ffffffffffff8160101c16602984015366ffffffffffffff8160081c16602a8401536001600160401b0316602b830153600c830183116106c8576020602c910151910152602c81018091116106c8576001909b019a6105b5565b634e487b7160e01b5f52601160045260245ffd5b5087908a8a6001600160401b038111610c1d57602484019060ff8160381c16825361ffff8160301c16602586015362ffffff8160281c16602686015363ffffffff8160201c16602786015364ffffffffff8160181c16602886015365ffffffffffff8160101c16602986015366ffffffffffffff8160081c16602a8601536001600160401b0316602b8501535f606060405161077781610db9565b82815282602082015282604082015201526030845110610c0a576356414c34835160e01c03610c0a575160c01c945160f01c9051955160f01c95604051956107be87610db9565b865260208601968288526040870191825280606088015282159283918415610bff575b8415610bf5575b508315610bca575b50508115610bb3575b50610ba0576040519161082e60218460208101945f8652845180918484015e81015f838201520301601f198101855284610def565b616000835111610b8d575061088f602a61ffff60f01b845160f01b16936040519384916020830196606160f81b885260218401526680600a3d393df360c81b60238401525180918484015e81015f838201520301601f198101835282610def565b51905ff06001600160a01b03168015610b7e575f51602061618b5f395f51905f5280545f5160206161cb5f395f51905f52939093559251935161ffff60e01b60e09190911b16600160a01b600160e01b0360a09590951b949094166001600160f01b0319909216600160a01b600160f01b031991909116171791909117905551602001515f5160206160eb5f395f51905f52546001600160401b039091169080610b3b575b506001600160401b0381165f9081525f51602061620b5f395f51905f52602052604090206001600160a01b03906001015416155f5160206161cb5f395f51905f52545f51602061618b5f395f51905f52546040519161099283610db9565b8252602082019160018060a01b038216835260016040820191818060401b038460a01c16835261ffff606082019460e01c1684526109ec8760018060401b03165f525f51602061620b5f395f51905f5260205260405f2090565b90518155935193018054915192516001600160f01b03199092166001600160a01b03949094169390931760a09290921b600160a01b600160e01b03169190911760e09190911b61ffff60e01b16179055610ac7575b506001600160a01b038116610aa15750610a596112f1565b50610a6560e051611373565b505b6040516148a59081611826823960805181610e84015260a051816137bc015260c051816105ee015260e05181818161032801526106d00152f35b80610aae610ac1926111fb565b50610ab881611271565b5060e0516113cc565b50610a67565b5f5160206160eb5f395f51905f525468010000000000000000811015610b2757806001610b0392015f5160206160eb5f395f51905f52556111a3565b81546001600160401b0360039290921b91821b191692901b91909117905581610a41565b634e487b7160e01b5f52604160045260245ffd5b5f1981019081116106c857610b4f906111a3565b905460039190911b1c6001600160401b03168082101561093457630fb2737b60e31b5f5260045260245260445ffd5b63fbad885d60e01b5f5260045ffd5b516331734c2160e21b5f5260045260245ffd5b82634724a0fd60e01b5f5260045260245ffd5b905051610bc2610595856117eb565b1415886107f9565b8551929350602c808202929183041417156106c85760300190816030116106c85714159089806107f0565b151593508b6107e8565b60b4821194506107e1565b84634724a0fd60e01b5f5260045260245ffd5b6305f8ded760e21b5f5260045ffd5b8363156f758160e31b5f5260045260b460245260445ffd5b6333fae18560e21b5f5260045260245260445ffd5b6355bace6f60e01b5f5260045260245260445ffd5b015190505f8061032e565b90601f198316915f5160206161eb5f395f51905f525f52815f20925f5b818110610ce3575091600193918561045e9897969410610ccb575b505050811b015f5160206161eb5f395f51905f525561034f565b01515f1960f88460031b161c191690555f8080610cb1565b92936020600181928786015181550195019301610c96565b5f5160206161eb5f395f51905f525f527fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be601f840160051c81019160208510610d60575b601f0160051c01905b818110610d555750610310565b5f8155600101610d48565b9091508190610d3f565b634e487b7160e01b5f52602260045260245ffd5b90607f16906102fe565b5f80fd5b81516001600160401b038111610d8857602091610dae8a848094880101610ebf565b8152019101906100f3565b608081019081106001600160401b03821117610b2757604052565b604081019081106001600160401b03821117610b2757604052565b601f909101601f19168101906001600160401b03821190821017610b2757604052565b51906001600160a01b0382168203610d8857565b6001600160401b038111610b2757601f01601f191660200190565b929192610e4d82610e26565b91610e5b6040519384610def565b829481845281830111610d88578281602093845f96015e010152565b9080601f83011215610d88578151610e9192602001610e41565b90565b6001600160401b038111610b275760051b60200190565b51906001600160401b0382168203610d8857565b9190608083820312610d885760405190610ed882610db9565b8351919384926001600160401b038111610d8857606092610efa918301610e77565b835260208101516020840152610f1260408201610eab565b60408401520151908160070b8203610d885760600152565b51908115158203610d8857565b519060ff82168203610d8857565b519063ffffffff82168203610d8857565b15610f5e5750565b63ffffffff9063b0369c3160e01b5f5216600452600160245263ffffffff60445260645ffd5b919082039182116106c857565b908151811015610fa2570160200190565b634e487b7160e01b5f52603260045260245ffd5b919082018092116106c857565b604051610fcf81610dd4565b606081525f60208201525080518015908115611197575b50611188575f19908051805b611144575b505f1982146111365760018201908183116106c857600360fc1b6001600160f81b03196110248484610f91565b51161480611122575b611112575f5b81518310156110c2576110468383610f91565b5160f81c6030811080156110b8575b6110a657600a82026001600160401b03908116602f1990920160ff169190910181169116811061108a57600190920191611033565b509150506040519061109b82610dd4565b81525f602082015290565b50509150506040519061109b82610dd4565b5060398111611055565b9150916001811190811591611106575b506110f757604051916110e483610dd4565b82526001600160401b0316602082015290565b63384c687760e21b5f5260045ffd5b602b915010155f6110d2565b9150506040519061109b82610dd4565b506002611130848351610f84565b1161102d565b604051915061109b82610dd4565b5f1981018181116106c857602d60f81b6001600160f81b03196111678386610f91565b51161461117e575080156106c8575f190180610ff2565b92505f9050610ff7565b6329120bff60e21b5f5260045ffd5b6040915010155f610fe6565b905f5160206160eb5f395f51905f5254821015610fa2575f5160206160eb5f395f51905f525f52600282901c7f06ff75c3bdc3f03474569b3e952aab1bfb9d41b303b5638523d3e1200dcd1ae2019160031b60181690565b6001600160a01b0381165f9081525f51602061622b5f395f51905f52602052604090205460ff1661126c576001600160a01b03165f8181525f51602061622b5f395f51905f5260205260408120805460ff191660011790553391905f5160206160cb5f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f51602061612b5f395f51905f52602052604090205460ff1661126c576001600160a01b03165f8181525f51602061612b5f395f51905f5260205260408120805460ff191660011790553391905f5160206161ab5f395f51905f52905f5160206160cb5f395f51905f529080a4600190565b5f80525f51602061612b5f395f51905f526020525f51602061616b5f395f51905f525460ff1661136f575f8080525f51602061612b5f395f51905f526020525f51602061616b5f395f51905f52805460ff1916600117905533905f5160206161ab5f395f51905f525f5160206160cb5f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f1461126c57805f525f60205260405f205f805260205260405f20600160ff198254161790555f33915f5160206160cb5f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff1661143b575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291905f5160206160cb5f395f51905f529080a4600190565b50505f90565b8051821015610fa25760209160051b010190565b90815151801561163c5761146881610e94565b906114766040519283610def565b808252601f1961148582610e94565b013660208401375f5b8451805182101561162c576114a582602092611441565b51015160018060401b0360406114bc848951611441565b51015116604051916114cd83610dd4565b82526020820190808252602490806115ca575b506114ea90611642565b91600a60208401536001908184019060226020830153828001918284116106c8576021600a910153600283018092116106c8576115349161152b9086611793565b905190856117a9565b82519092906001600160401b0316611567575b50505090611556600192611674565b6115608286611441565b520161148e565b82840190601060208301538284018094116106c857516001600160401b03169260219190910191905b60808410156115ad5750509161155691600194935391925f611547565b91908084926080607f8397161781530192019060071c9290611590565b6001905b608081101561161e57508060010160011161160a57602501908181116115f757506114ea6114e0565b634e487b7160e01b5f9081526011600452fd5b50634e487b7160e01b5f9081526011600452fd5b60019060071c9101906115ce565b505090610e919293505f906116d2565b505f9150565b9061164c82610e26565b6116596040519182610def565b828152809261166a601f1991610e26565b0190602036910137565b8051600181018091116106c85761168a90611642565b805115610fa257602091816116a6845f940192848453826117c1565b50604051918291518091835e8101838152039060025afa156116c7575f5190565b6040513d5f823e3d90fd5b9291926116df8285610f84565b93600185146117835760015b8060011b90868210156116fe57506116eb565b9394955050906117186117118486610fb6565b85836116d2565b61172c936117269195610fb6565b906116d2565b9060405161173b608082610def565b6041815260208101906060368337805115610fa2576020935f936001845360218301526041820152604051918291518091835e8101838152039060025afa156116c7575f5190565b5061178f929350611441565b5190565b6020828192010153600181018091116106c85790565b81602091939293010152602081018091116106c85790565b9080519182156117e3576021602084930191015e600101806001116106c85790565b505050600190565b611821604051916117fd606084610def565b602283526040366020850137600a602084015361181b600184611793565b836117a9565b509056fe6080806040526004361015610012575f80fd5b5f3560e01c90816301ffc9a71461115157508063248a9ca3146111275780632f2ff15d146110f857806336568abe146110a95780634b1872be14610d7e5780636a28f00014610c9d5780636edfe3af14610b0c5780638a8e4c5d14610ad45780638c80fbda14610a9857806391d1485414610a5c57806392c19bdc146109b6578063974a74c41461085f578063a217fddf14610845578063a6f031bb14610729578063d547741f146106f3578063db3e1fa4146106b9578063ddba6537146102ff5763ef913a4b146100e2575f80fd5b346102fb575f3660031901126102fb5760405160208082015261012060408201525f7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1700548060011c906001811680156102f1575b6020831081146102dd578261016086015290815f146102b85750600114610238575b61023483610220818560ff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170154818116606085015260081c1660808301526001600160401b037f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025481811660a085015260401c1660c083015263ffffffff5f5160206148795f395f51905f525481811660e0850152818160201c1661010085015260ff8160401c16151561012085015260481c1661014083015203601f19810183528261132a565b6040519182916020835260208301906112b5565b0390f35b9190507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005f527fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be915f905b80821061029c5750909150810161018001610220610158565b9192600181602092546101808588010152019101909291610283565b60ff19166101808086019190915291151560051b840190910191506102209050610158565b634e487b7160e01b5f52602260045260245ffd5b91607f1691610136565b5f80fd5b346102fb5761030d36611215565b60ff5f5160206148795f395f51905f525460401c16610691577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610682575b508101906020818303126102fb578035906001600160401b0382116102fb5701610140818303126102fb5760405160c081018181106001600160401b0382111761066e5760405281356001600160401b0381116102fb5782016040818503126102fb57604051906103d6826112d9565b80356001600160401b0381116102fb57856103f29183016114d8565b82526020810135906001600160401b0382116102fb57610414918691016114d8565b60208201528152610428836020840161135f565b916020820192835261043d846080830161135f565b6040830190815261045060e0830161134b565b92606081019384526101008301356001600160401b0381116102fb5786610478918501611877565b9560808201968752610120840135956001600160401b0387116102fb57610612976001600160801b036105ec976105af6105de976105886105616104c66105ca9960209f61054b9e01611877565b9960a081019a8b526104da87875116612834565b6104e8885182515190613cf5565b6104f98b5160208351015190613cf5565b6040519e8f9d8e927f60312a7c0000000000000000000000000000000000000000000000000000000083850152826024850152519261014060448201526101c484519160406101848201520190611ae1565b9101518d820361018319016101a48f0152611ae1565b965180516001600160801b031660648d0152602081015160848d01526040015160a48c0152565b5180516001600160801b031660c48b0152602081015160e48b0152604001516101048a0152565b51166101248701525185820360431901610144870152611d67565b905183820360431901610164850152611d67565b03601f19810183528261132a565b7f00000000000000000000000000000000000000000000000000000000000000006128bf565b506801000000000000000068ff0000000000000000195f5160206148795f395f51905f525416175f5160206148795f395f51905f52557f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1005b634e487b7160e01b5f52604160045260245ffd5b61068b9061267f565b82610366565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102fb575f3660031901126102fb5760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346102fb57610727610704366111ef565b9061072261071d825f525f602052600160405f20015490565b61267f565b61274c565b005b346102fb5760203660031901126102fb576004356001600160401b0381116102fb57806004019061014060031982360301126102fb5760ff5f5160206148795f395f51905f525460401c16610691575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff1615610838575b6107d760448201836125db565b916107e560648201856125db565b9390916101048101359060018210156102fb576020966108309661080d6101248401836125db565b9690956040519861081e8c8b61132a565b5f8a52608460a4870196013594613622565b604051908152f35b610840612610565b6107ca565b346102fb575f3660031901126102fb5760206040515f8152f35b346102fb5760203660031901126102fb576004356001600160401b0381116102fb57806004019061016060031982360301126102fb5760ff5f5160206148795f395f51905f525460401c16610691575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff16156109a9575b610144810161090f81846125a9565b9050156109815761092360448301846125db565b909161093260648501866125db565b9590946101048101359160018310156102fb576020976108309761096a966109716109616101248701866125db565b999098866125a9565b36916113fd565b98608460a4870196013594613622565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b6109b1612610565b610900565b346102fb576109c436611215565b9060ff5f5160206148795f395f51905f525460401c16610691575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060209081527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47549092610a3e92909160ff1615610a4f576120b2565b60405190610a4b81611264565b8152f35b610a57612610565b6120b2565b346102fb57610a6a366111ef565b905f525f6020526001600160a01b0360405f2091165f52602052602060ff60405f2054166040519015158152f35b346102fb575f3660031901126102fb5760207f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170554604051908152f35b346102fb57610ae236611215565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102fb575f3660031901126102fb576020610b266129fb565b610b2f81612a87565b92019161ffff83511691610b5b610b45846114c1565b93610b53604051958661132a565b8085526114c1565b602084019190601f1901368337610b7661ffff865116611fef565b9261ffff86511691610b8a610b45846114c1565b602084019690601f19013688375f5b61ffff89511661ffff821690811015610bfe5761ffff91886001600160401b03610bf384808f610be482610bec928e8c8f8f9d610bd960019f8290612021565b525192511691612b38565b929096612021565b528a612021565b911690520116610b99565b87838a888a604051948594606086019060608752518091526080860192905f5b818110610c7e5750505081610c3b91866020940384880152611282565b91848303604086015251918281520191905f5b818110610c5c575050500390f35b82516001600160401b0316845285945060209384019390920191600101610c4e565b825163ffffffff16855288975060209485019490920191600101610c1e565b346102fb575f3660031901126102fb57335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff1615610d67575f5160206148795f395f51905f525460ff8160401c1615610d3f5768ff000000000000000019165f5160206148795f395f51905f52557f8e0da210a2ffecc744373b4860cd9b8dc471810f5fb61b7d2c3234ccd4aae30d5f80a1005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b346102fb57610d8c36611215565b60ff5f5160206148795f395f51905f525460401c16610691575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff161561109c575b8101906020818303126102fb578035906001600160401b0382116102fb576105de610e82610e2f602095610eb695610ea89501611a64565b610e456001600160801b03604083015116612834565b610e676060820151606088808501515101510151610e616129fb565b91613f3c565b604051928391631f011b7b60e21b8884015260248301611e4d565b7f00000000000000000000000000000000000000000000000000000000000000006128bf565b828082518301019101611f4d565b610ebf8161290d565b90610ec982611264565b8161102c57826001600160401b0381838160607f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb3196019184828451015116610f3b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025491878360401c16808211611fab565b83516fffffffffffffffffffffffffffffffff196fffffffffffffffff0000000000000000858984511693015160401b16921617177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1702550151604051610fc18482018093604080916001600160801b038151168452602081015160208501520151910152565b60608152610fd060808261132a565b519020611013848484510151166001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b5551015116604051908152a160405190610a4b81611264565b5061103681611264565b60018103610a3e576801000000000000000068ff0000000000000000195f5160206148795f395f51905f525416175f5160206148795f395f51905f52557f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1610a3e565b6110a4612610565b610df7565b346102fb576110b7366111ef565b336001600160a01b038216036110d0576107279161274c565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102fb57610727611109366111ef565b9061112261071d825f525f602052600160405f20015490565b6126bf565b346102fb5760203660031901126102fb5760206108306004355f525f602052600160405f20015490565b346102fb5760203660031901126102fb57600435907fffffffff0000000000000000000000000000000000000000000000000000000082168092036102fb57817f7965db0b00000000000000000000000000000000000000000000000000000000602093149081156111c5575b5015158152f35b7f01ffc9a700000000000000000000000000000000000000000000000000000000915014836111be565b60409060031901126102fb57600435906024356001600160a01b03811681036102fb5790565b9060206003198301126102fb576004356001600160401b0381116102fb57826023820112156102fb578060040135926001600160401b0384116102fb57602484830101116102fb576024019190565b6003111561126e57565b634e487b7160e01b5f52602160045260245ffd5b90602080835192838152019201905f5b81811061129f5750505090565b8251845260209384019390920191600101611292565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b604081019081106001600160401b0382111761066e57604052565b606081019081106001600160401b0382111761066e57604052565b608081019081106001600160401b0382111761066e57604052565b90601f801991011681019081106001600160401b0382111761066e57604052565b35906001600160801b03821682036102fb57565b91908260609103126102fb57604051611377816112f4565b60408082946113858161134b565b8452602081013560208501520135910152565b35906001600160401b03821682036102fb57565b91908260409103126102fb576040516113c4816112d9565b60206113dd8183956113d581611398565b855201611398565b910152565b6001600160401b03811161066e57601f01601f191660200190565b929192611409826113e2565b91611417604051938461132a565b8294818452818301116102fb578281602093845f960137010152565b359081151582036102fb57565b359063ffffffff821682036102fb57565b8092910391606083126102fb5760405161146a816112d9565b6040819483358352601f1901126102fb57602090604080519361148c856112d9565b611497848201611440565b85520135828401520152565b9080601f830112156102fb578160206114be933591016113fd565b90565b6001600160401b03811161066e5760051b60200190565b91906060838203126102fb57604051906114f1826112d9565b819380356001600160401b0381116102fb5781016040818403126102fb576040519061151c826112d9565b80356001600160401b0381116102fb5781016102c0818603126102fb576040519061026082018281106001600160401b0382111761066e5760405261156186826113ac565b825260408101356001600160401b0381116102fb57810186601f820112156102fb5786816020611593933591016113fd565b60208301526115a460608201611398565b60408301526115b56080820161134b565b60608301526115c660a08201611433565b60808301526115d88660c08301611451565b60a08301526115ea6101208201611433565b60c083015261014081013560e08301526116076101608201611433565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526116566102208201611433565b6101c08301526102408101356101e08301526116756102608201611433565b6102008301526102808101356102208301526102a0810135906001600160401b0382116102fb576116a8918791016114a3565b61024082015282526020810135906001600160401b0382116102fb570160c0818503126102fb57604051906116dc8261130f565b6116e581611398565b82526116f360208201611440565b60208301526117058560408301611451565b604083015260a0810135906001600160401b0382116102fb570184601f820112156102fb578035611735816114c1565b91611743604051938461132a565b81835260208084019260051b820101908782116102fb57602001915b818310611785575050509260209492826113dd95606088950152838201528652016113ac565b6020838903126102fb5760405190602082018281106001600160401b0382111761066e5760405283359060048210156102fb579082529081526020928301920161175f565b9080601f830112156102fb576040805192906117e6908461132a565b8290604081019283116102fb57905b8282106118025750505090565b81358152602091820191016117f5565b9080601f830112156102fb578135611829816114c1565b92611837604051948561132a565b81845260208085019260051b8201019283116102fb57602001905b82821061185f5750505090565b6020809161186c84611440565b815201910190611852565b9190610220838203126102fb5760405161010081018181106001600160401b0382111761066e57604052809382601f820112156102fb576040516118bd6101008261132a565b806101008301918583116102fb5790859184905b848210611a5157505084526118e5916117ca565b60208301526118f88361014083016117ca565b604083015261018081013561ffff811681036102fb5760608301526101a08101356001600160401b0381116102fb5783611933918301611812565b60808301526101c08101356001600160401b0381116102fb5783611958918301611812565b60a08301526101e08101356001600160401b0381116102fb57810183601f820112156102fb5780359061198a826114c1565b91611998604051938461132a565b80835260208084019160051b830101918683116102fb57602001905b828210611a415750505060c0830152610200810135906001600160401b0382116102fb57019180601f840112156102fb5782356119f0816114c1565b936119fe604051958661132a565b81855260208086019260051b8201019283116102fb57602001905b828210611a295750505060e00152565b60208091611a3684611433565b815201910190611a19565b81358152602091820191016119b4565b81358152879350602091820191016118d1565b919060c0838203126102fb5760405190611a7d8261130f565b8193611a89828261135f565b835260608101356001600160401b0381116102fb5782611aaa9183016114d8565b6020840152611abb6080820161134b565b604084015260a0810135916001600160401b0383116102fb576060926113dd9201611877565b90815191606082526020611c3e845160406060860152611b1a60a0860182516001600160401b0360208092828151168552015116910152565b610240611b37848301516102c060e08901526103608801906112b5565b60408301516001600160401b031661010088015260608301516001600160801b03166101208801526080830151151561014088015260a08301518051610160890152602090810151805163ffffffff166101808a015201516101a08801529160c081015115156101c088015260e08101516101e08801526101008101511515610200880152610120810151610220880152610140810151828801526101608101516102608801526101808101516102808801526101a08101516102a08801526101c081015115156102c08801526101e08101516102e088015261020081015115156103008801526102208101516103208801520151609f19868303016103408701526112b5565b930151605f19838503016080840152602060e0606060c08701936001600160401b03815116885263ffffffff848201511684890152611c9f604082015160408a019060208060409280518552015163ffffffff815116828501520151910152565b01519560c060a08201528651809452019401905f905b808210611ce35750505060209081015180516001600160401b03908116838501529101511660409091015290565b90919485515190600482101561126e57602081600193829352019601920190611cb5565b905f905b60028210611d1857505050565b6020806001928551815201930191019091611d0b565b90602080835192838152019201905f5b818110611d4b5750505090565b825163ffffffff16845260209384019390920191600101611d3e565b80515f835b60088210611e3757505050611d8a6020820151610100840190611d07565b611d9d6040820151610140840190611d07565b61ffff60608201511661018083015260e0611df6611de3611dd060808501516102206101a0880152610220870190611d2e565b60a08501518682036101c0880152611d2e565b60c08401518582036101e0870152611282565b91015191610200818303910152602080835192838152019201905f5b818110611e1f5750505090565b82511515845260209384019390920191600101611e12565b6020806001928551815201930191019091611d6c565b906114be9160208152611e81602082018351604080916001600160801b038151168452602081015160208501520151910152565b6060611e9c602084015160c0608085015260e0840190611ae1565b926001600160801b0360408201511660a084015201519060c0601f1982850301910152611d67565b91908260609103126102fb57604051611edc816112f4565b809280516001600160801b03811681036102fb5760409182918452602081015160208501520151910152565b51906001600160401b03821682036102fb57565b91908260409103126102fb57604051611f34816112d9565b60206113dd818395611f4581611f08565b855201611f08565b90610140828203126102fb57611fa39061010060405193611f6d8561130f565b611f778382611ec4565b8552611f868360608301611ec4565b6020860152611f988360c08301611f1c565b604086015201611f1c565b606082015290565b15611fb4575050565b906001600160401b0380927f7d939bd8000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b90611ff9826114c1565b612006604051918261132a565b8281528092612017601f19916114c1565b0190602036910137565b80518210156120355760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b91906080838203126102fb57604051906120628261130f565b819380356001600160401b0381116102fb576060926120829183016114a3565b83526020810135602084015261209a60408201611398565b60408401520135908160070b82036102fb5760600152565b908101906020818303126102fb578035906001600160401b0382116102fb57016040818303126102fb57604051906120e9826112d9565b80356001600160401b0381116102fb5783612105918301611a64565b82526020810135906001600160401b0382116102fb57016080818403126102fb57604051926121338461130f565b81356001600160401b0381116102fb57820181601f820112156102fb57803561215b816114c1565b91612169604051938461132a565b81835260208084019260051b820101918483116102fb5760208201905b83821061257c5750505050845261219f60208301611433565b60208501526040820135916001600160401b0383116102fb576121c96060926121d4948301612049565b604086015201611398565b606083015260208101918252516121f76001600160801b03604083015116612834565b61224f6122406105de610e8260608501519461222460208201966060602089515101510151610e616129fb565b604051928391631f011b7b60e21b602084015260248301611e4d565b60208082518301019101611f4d565b916122598361290d565b9261226384611264565b6001841461251d576101606122788351612bfa565b9351515101518084036124ed575060207f1411889dcdfcf9af1fcf4ff4521cab4580a1fd842401ead543815cb1148d643893926001600160401b03926122bd87611264565b86612474578061230b60608593019361230587858751015116887f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025460401c16808211611fab565b51612dc1565b8251858151166fffffffffffffffffffffffffffffffff196fffffffffffffffff0000000000000000857f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025494015160401b16921617177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025501516040516123b38482018093604080916001600160801b038151168452602081015160208501520151910152565b606081526123c260808261132a565b519020612405848484510151166001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b55612415838383510151166133a4565b7f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb318284818451015116604051908152a15101511661246e60405192839283602090939291936001600160401b0360408201951681520152565b0390a190565b60606124b591019161230585858551015116867f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025460401c16808214611fab565b6124c4838383510151166133a4565b5101511661246e60405192839283602090939291936001600160401b0360408201951681520152565b83907f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b5050506801000000000000000068ff0000000000000000195f5160206148795f395f51905f525416175f5160206148795f395f51905f52557f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a190565b81356001600160401b0381116102fb5760209161259e88848094880101612049565b815201910190612186565b903590601e19813603018212156102fb57018035906001600160401b0382116102fb576020019181360383136102fb57565b903590601e19813603018212156102fb57018035906001600160401b0382116102fb57602001918160051b360383136102fb57565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff161561264857565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260405f206001600160a01b0333165f5260205260ff60405f205416156126a95750565b63e2517d3f60e01b5f523360045260245260445ffd5b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f205416155f1461274657805f525f60205260405f206001600160a01b0383165f5260205260405f20600160ff198254161790556001600160a01b03339216907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f2054165f1461274657805f525f60205260405f206001600160a01b0383165f5260205260405f2060ff1981541690556001600160a01b03339216907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b919082018092116127dc57565b634e487b7160e01b5f52601160045260245ffd5b156127f9575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b919082039182116127dc57565b6001600160801b03633b9aca0091160463ffffffff5f5160206148795f395f51905f525460481c16612873824261286b84426127cf565b8211156127f0565b42821061287e575050565b6128888242612827565b116128905750565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b5f918291602082519201905af43d15612905573d906128dd826113e2565b916128eb604051938461132a565b82523d5f602084013e5b156128fd5790565b602081519101fd5b6060906128f5565b6001600160401b03602060608301510151165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f205480155f146129585750505f90565b6020820190815160405161298d602082018093604080916001600160801b038151168452602081015160208501520151910152565b6060815261299c60808261132a565b51902014918215926129ba575b5050156129b557600190565b600290565b6001600160801b0391925081905151169151511611155f806129a9565b604051906129e48261130f565b5f6060838281528260208201528260408201520152565b612a036129d7565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055461ffff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17065460405192612a588461130f565b83526001600160a01b03811660208401526001600160401b038160a01c16604084015260e01c16606082015290565b612a8f6129d7565b50602081016001600160a01b0381511615612b24576001600160a01b03905116612ac061ffff606084015116614225565b813b6001811115612b11575f1981019081116127dc578111612afe5790816001612aec612afb9461424a565b9260208401903c809251614272565b91565b50633610565160e21b5f5260045260245ffd5b82633610565160e21b5f5260045260245ffd5b5051633a517eed60e21b5f5260045260245ffd5b939263ffffffff169161ffff16821015612bca57602c820293828504602c14831517156127dc578460300191826030116127dc57850192605084015160e01c03612b9f57506034840181116127dc57603c605483015160c01c9401106127dc57605c015190565b7f4724a0fd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b50827ff81f5140000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b908151518015612dbb57612c0d81611fef565b905f5b84518051821015612dab57612c2782602092612021565b5101516001600160401b036040612c3f848951612021565b5101511660405191612c50836112d9565b8252602082019080825260249080612d4b575b50612c6d9061424a565b91600a60208401536001908184019060226020830153828001918284116127dc576021600a910153600283018092116127dc57612cb791612cae908661477c565b90519085614792565b916001600160401b03815116612ce8575b50505090612cd7600192614348565b612ce18286612021565b5201612c10565b82840190601060208301538284018094116127dc57516001600160401b03169260219190910191905b6080841015612d2e57505091612cd791600194935391925f612cc8565b91908084926080607f8397161781530192019060071c9290612d11565b6001905b6080811015612d9d575080600101600111612d8a5760250190818111612d785750612c6d612c63565b634e487b7160e01b5f5260116004525ffd5b50634e487b7160e01b5f5260116004525ffd5b60019060071c910190612d4f565b5050906114be9293505f9061439b565b505f9150565b90612dcb82612bfa565b915180519182156132ac5760b483116132d4575f91602c8402848104602c036127dc57603001806030116127dc57612e029061424a565b9160208301946056865361564160218501536256414c60228501536356414c3460238501538060081c602c850153602d840153612e46612e418761483e565b614348565b602e8401525f604e8401535f604f8401535f9560305b8351881015612f6f57612e6f8885612021565b5195612e8a60408801916001600160401b03835116906127cf565b968287019160ff8b60181c16602084015361ffff8b60101c16602184015362ffffff8b60081c16602284015363ffffffff8b1660238401536004840184116127dc576001600160401b03905160ff8160381c16602485015361ffff8160301c16602585015362ffffff8160281c16602685015363ffffffff8160201c16602785015364ffffffffff8160181c16602885015365ffffffffffff8160101c16602985015366ffffffffffffff8160081c16602a85015316602b830153600c830183116127dc576020602c910151910152602c81018091116127dc57600190970196612e5c565b509493915094506001600160401b0381116132ac578060ff6001600160401b039260381c16602484015361ffff8160301c16602584015362ffffff8160281c16602684015363ffffffff8160201c16602784015364ffffffffff8160181c16602884015365ffffffffffff8160101c16602984015366ffffffffffffff8160081c16602a84015316602b8201536130068184614272565b916040519161303760218460208101945f8652845180918484015e81015f838201520301601f19810185528461132a565b6160008351116132805750906001600160a01b03916130ef602a7fffff000000000000000000000000000000000000000000000000000000000000845160f01b169360405193849160208301967f6100000000000000000000000000000000000000000000000000000000000000885260218401527f80600a3d393df30000000000000000000000000000000000000000000000000060238401525180918484015e81015f838201520301601f19810183528261132a565b51905ff0169182156132585760209273ffffffffffffffffffffffffffffffffffffffff197f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17065416177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706557f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055580517fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b807f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706549360a01b16169116177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655015161ffff60e01b1961ffff60e01b807f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706549360e01b16169116177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b827fab7bac08000000000000000000000000000000000000000000000000000000005f5260045260b460245260445ffd5b907f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170854821015612035577f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17085f52600282901c7f06ff75c3bdc3f03474569b3e952aab1bfb9d41b303b5638523d3e1200dcd1ae2019160031b60181690565b9190918054831015612035575f52601860205f208360021c019260031b1690565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170854806135e2575b506001600160a01b036001613411836001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f2090565b015416157f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1705547f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170654604051916134668361130f565b825260208201916001600160a01b03821683526001600160a01b03600160408301926001600160401b038560a01c16845261ffff606082019560e01c1685526134df886001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f2090565b905181550193511673ffffffffffffffffffffffffffffffffffffffff19845416178355517fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b8085549360a01b161691161782555161ffff60e01b1961ffff60e01b8084549360e01b161691161790556135625750565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708546801000000000000000081101561066e578060016135c492017f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170855613305565b6001600160401b0380839493549260031b9316831b921b1916179055565b5f1981019081116127dc57816001600160401b0361360261361c93613305565b90549060031b1c16806001600160401b0383161015611fab565b5f6133cc565b999a9890939294969a979195975f9a60018d101561126e578c1561366f5760248c60ff8f7f112d89cc00000000000000000000000000000000000000000000000000000000835216600452fd5b909192939495968098999a9c50151580613ce9575b15613cb257602001356001600160401b03811681036102fb576136a7368b61135f565b9061372360405160208101906136da8286604080916001600160801b038151168452602081015160208501520151910152565b606081526136e960808261132a565b519020916001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b548015613c8a57808203613c5c5750506020810151808a03613c2c57506001600160801b03633b9aca00915116045f5160206148795f395f51905f5254613778824261286b63ffffffff8560481c16426127cf565b5f91428110613c08575b5063ffffffff16906001600160801b031681811015613bda5750505f905f5b888110613af1575b505015613ab05750506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001693843b156102fb5794929085926040519687957f56455e790000000000000000000000000000000000000000000000000000000087526064870190600488015260606024880152526084850160848560051b87010194825f90603e1981360301935b838310613a3f57505050505050600319848403016044850152808352602083019260208260051b82010193835f92601e1982360301905b8585106138d957505050505050509181805f9403915afa80156138ce576138b9575b5035906001600160801b0382168092036138b65750633b9aca00900490565b80fd5b6138c69192505f9061132a565b5f905f613897565b6040513d5f823e3d90fd5b9193959750919395601f198282030185528735838112156102fb57840190613905602082019280614599565b8091936020845252604082019060408160051b8401019380935f915b83831061394857505050505050602080600192990195019501929091899796949592613875565b909192939495603f1983820301865261396187836145e1565b803560028110156102fb57825261398f61397e60208301836145cd565b6060602085015260608401906145f5565b906040810135609e19823603018112156102fb576001936020938493613a319301916040818303910152613a23613a056139da6139cc8580614509565b60a0865260a08601916144e9565b6139e5878601611433565b1515878501526139f860408601866145cd565b84820360408601526145f5565b92613a1260608201611433565b1515606084015260808101906145cd565b9060808184039101526145f5565b980196019493019190613921565b9193959698509193966083198b82030182528735868112156102fb576020613a9d6001938683940190613a90613a86613a788480614599565b60408552604085019161453a565b9285810190614509565b91858185039101526144e9565b9901920193019093918a9896959361383e565b613aed6040519283927ffef760c700000000000000000000000000000000000000000000000000000000845260206004850152602484019161453a565b0390fd5b9091613b25613b14613b0d613b07858d8c61445c565b806125db565b369161447e565b613b1f36878961447e565b906147d4565b15613bd0575061096a613b3c613b46928a8961445c565b60208101906125a9565b805182518082149182613bba575b505015613b6657505060015f806137a9565b90613aed613ba8926040519384937f5f1ca3810000000000000000000000000000000000000000000000000000000085526040600486015260448501906112b5565b838103600319016024850152906112b5565b9091506020830120906020840120145f80613b54565b91906001016137a1565b7fdb020436000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160801b0380929350613c2363ffffffff9242612827565b16929150613782565b89907f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b877fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b5061ffff881115613684565b9190916001600160401b03602080606081875101510151950151015116613d1a6129d7565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708548015613e30575f905f1981019081116127dc575b808210613e435750613d8b6001600160401b03917f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708613383565b90549060031b1c16908111613e30575f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f209260405191613dd28361130f565b6001855495868552015461ffff6001600160a01b038216918260208701526001600160401b038160a01c16604087015260e01c16606085015215613e1d57613e1b939450613f3c565b565b84633a517eed60e21b5f5260045260245ffd5b633a517eed60e21b5f525f60045260245ffd5b90613e4e82826127cf565b600181018091116127dc5760011c90836001600160401b03613e90847f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708613383565b905460039190911b1c1611613ea6575090613d51565b91505f19810190811115613d5157634e487b7160e01b5f52601160045260245ffd5b15613ed05750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15613f0a5750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b9092919260808201515161ffff6060840151168091149081614216575b81614207575b816141f8575b50156132ac57613f7484612a87565b6040860151955f959293869384918291825b61ffff60608b0151168b10156140b257613fa48b60e08c0151612021565b51156140a85763ffffffff613fbd8c60808d0151612021565b51168098614090575b505060019692898b60a082015190613fdd91612021565b5163ffffffff168c60208c019782895161ffff16811090613ffd91613f02565b8b831b9061400e8482841615613ec8565b1797828b8b51925161ffff169061402493612b38565b91909360c001519061403591612021565b510361406557506001600160401b038091169116016001600160401b0381116127dc576001909a5b019992613f86565b7fc3271bcd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b63ffffffff6140a192168111613ec8565b5f87613fc6565b929960019061405d565b509850989550925050506001600160401b0391501660038102818104600314821517156127dc576801fffffffffffffffe8360011b16906001600160401b03841682046002146001600160401b0385161517156127dc5711156141bf57505060e0608082015191015181518151036132ac575f5b82518110156141b8576141398183612021565b51156141b05763ffffffff61414e8285612021565b511661415d8187518110613f02565b6141678187612021565b5151600481101561126e576001190161418557506001905b01614126565b7f5b5c80eb000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b60019061417f565b5050509050565b6001600160401b0392507f6e3083c3000000000000000000000000000000000000000000000000000000005f526004521660245260445ffd5b905060e083015151145f613f65565b60c08401515181149150613f5f565b60a08401515181149150613f59565b61ffff16602c810290808204602c14901517156127dc57603001806030116127dc5790565b90614254826113e2565b614261604051918261132a565b8281528092612017601f19916113e2565b91909161427d6129d7565b506030835110612b9f576356414c34602084015160e01c03612b9f57602483015160c01c602c84015160f01c602e85015190604e86015160f01c92604051906142c58261130f565b81528160208201526040810192835283606082015295811593841561433d575b8415614333575b50831561431c575b50508115614305575b50612b9f5750565b905051614314612e418361483e565b14155f6142fd565b5191925061432990614225565b1415905f806142f4565b151593505f6142ec565b60b4831194506142e5565b8051600181018091116127dc5761435e9061424a565b805115612035576020918161437a845f940192848453826147aa565b50604051918291518091835e8101838152039060025afa156138ce575f5190565b9291926143a88285612827565b936001851461444c5760015b8060011b90868210156143c757506143b4565b9394955050906143e16143da84866127cf565b858361439b565b6143f5936143ef91956127cf565b9061439b565b9060405161440460808261132a565b6041815260208101906060368337805115612035576020935f936001845360218301526041820152604051918291518091835e8101838152039060025afa156138ce575f5190565b50614458929350612021565b5190565b91908110156120355760051b81013590603e19813603018212156102fb570190565b92919061448a816114c1565b93614498604051958661132a565b602085838152019160051b8101918383116102fb5781905b8382106144be575050505050565b81356001600160401b0381116102fb576020916144de87849387016114a3565b8152019101906144b0565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156102fb5701602081359101916001600160401b0382116102fb5781360383136102fb57565b90602083828152019260208260051b82010193835f925b8484106145615750505050505090565b909192939495602080614589600193601f198682030188526145838b88614509565b906144e9565b9801940194019294939190614551565b9035601e19823603018112156102fb5701602081359101916001600160401b0382116102fb578160051b360383136102fb57565b9035607e19823603018112156102fb570190565b9035605e19823603018112156102fb570190565b61462e6146136146058380614509565b6080865260808601916144e9565b6146206020840184614509565b9085830360208701526144e9565b61463b60408301836145cd565b908381036040850152813560038110156102fb5761465881611264565b8152602082013560038110156102fb5761467181611264565b602082015260408201359060038210156102fb5760806146ad6146c8948461469e6146bd96999899611264565b60408501526060810190614509565b91909281606082015201916144e9565b926060810190614599565b90916060818503910152808352602083019060208160051b85010193835f915b8383106146f85750505050505090565b909192939495601f1982820301865261471187846145e1565b80359160038310156102fb5761476e602092839285614731600197611264565b815261476061475561474586850185614509565b60608886015260608501916144e9565b926040810190614509565b9160408185039101526144e9565b9801960194930191906146e8565b6020828192010153600181018091116127dc5790565b81602091939293010152602081018091116127dc5790565b9080519182156147cc576021602084930191015e600101806001116127dc5790565b505050600190565b908151815103612746575f5b82518110156147cc576147f38184612021565b51516147ff8284612021565b515103614837576148108184612021565b51602081519101206148228284612021565b516020815191012003614837576001016147e0565b5050505f90565b6148746040519161485060608461132a565b602283526040366020850137600a602084015361486e60018461477c565b83614792565b509056fe5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1703a164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17085e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17014c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e05e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17036e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1707ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb55e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1702",
}

// ContractSpectreClientABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractSpectreClientMetaData.ABI instead.
var ContractSpectreClientABI = ContractSpectreClientMetaData.ABI

// ContractSpectreClientBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractSpectreClientMetaData.Bin instead.
var ContractSpectreClientBin = ContractSpectreClientMetaData.Bin

// DeployContractSpectreClient deploys a new Ethereum contract, binding an instance of ContractSpectreClient to it.
func DeployContractSpectreClient(auth *bind.TransactOpts, backend bind.ContractBackend, updateClientModule common.Address, membershipModule common.Address, misbehaviourModule common.Address, clientState_ []byte, consensusState [32]byte, initialPinnedValidatorSet IICS07TendermintMsgsValidatorSet, roleManager common.Address) (common.Address, *types.Transaction, *ContractSpectreClient, error) {
	parsed, err := ContractSpectreClientMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractSpectreClientBin), backend, updateClientModule, membershipModule, misbehaviourModule, clientState_, consensusState, initialPinnedValidatorSet, roleManager)
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
