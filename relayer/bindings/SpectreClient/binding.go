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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"updateClientModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membershipModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviourModule\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"clientState_\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"initialPinnedValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPinnedValidatorSet\",\"inputs\":[],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPinnedValidatorsHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateApplicationState\",\"inputs\":[{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateConsensusState\",\"inputs\":[{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ClientFrozen\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClientUnfrozen\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClientUpdated\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConsensusStateUpdated\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientNotFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DirectCallNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustedVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidMembershipProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NonMonotonicHeightUpdate\",\"inputs\":[{\"name\":\"latestHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"updateHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofSignerCommitSigMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorCountExceedsLimit\",\"inputs\":[{\"name\":\"validatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxValidatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x610100604052346101a557616c44803803809161001b826101ee565b61010039806101000160e082126101a5576100346102a3565b906100406101206102ba565b61004b6101406102ba565b610160516001600160401b0381116101a5578361006b91610100016102e9565b9061018051926101a0519660018060401b0388116101a557608090889003126101a5576040519461009b8661021a565b6101008801516001600160401b0381116101a55788018161011f820112156101a5576101008101516100cc8161032d565b916100da6040519384610250565b81835260206101008185019360051b83010101918483116101a5576101208201905b8382106101a95750505050865261011661012089016103c3565b6020870152610140880151906001600160401b0382116101a5578861014a610160926101006101559561016b9d0101610358565b604089015201610344565b60608601526101656101c06102ba565b95610999565b604051614f249081611bc0823960805181610447015260a05181613838015260c05181610eff015260e051818181610ddf0152610e3c0152f35b5f80fd5b81516001600160401b0381116101a5576020916101cf8884610100819589010101610358565b8152019101906100fc565b634e487b7160e01b5f52604160045260245ffd5b610100601f91909101601f19168101906001600160401b0382119082101761021557604052565b6101da565b608081019081106001600160401b0382111761021557604052565b604081019081106001600160401b0382111761021557604052565b601f909101601f19168101906001600160401b0382119082101761021557604052565b6040519061028361010083610250565b565b60405190610283604083610250565b60405190610283608083610250565b61010051906001600160a01b03821682036101a557565b51906001600160a01b03821682036101a557565b6001600160401b03811161021557601f01601f191660200190565b81601f820112156101a557602081519101610303826102ce565b926103116040519485610250565b828452828201116101a557815f926020928386015e8301015290565b6001600160401b0381116102155760051b60200190565b51906001600160401b03821682036101a557565b91906080838203126101a557604051906103718261021a565b8351919384926001600160401b0381116101a5576060926103939183016102e9565b8352602081015160208401526103ab60408201610344565b60408401520151908160070b82036101a55760600152565b519081151582036101a557565b519060ff821682036101a557565b91908260409103126101a5576040516103f681610235565b602061040f818395610407816103d0565b8552016103d0565b910152565b91908260409103126101a55760405161042c81610235565b602061040f81839561043d81610344565b855201610344565b519063ffffffff821682036101a557565b519060028210156101a557565b6020818303126101a5578051906001600160401b0382116101a55701610140818303126101a557610492610273565b8151909290916001600160401b0383116101a5576104d9826104bc610120946105299685016102e9565b86526104cb81602085016103de565b602087015260608301610414565b60408501526104ea60a08201610445565b60608501526104fb60c08201610445565b608085015261050c60e082016103c3565b60a085015261051e6101008201610456565b60c085015201610445565b60e082015290565b90600182811c9216801561055f575b602083101461054b57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691610540565b601f821161057657505050565b5f5260205f20906020601f840160051c830193106105ae575b601f0160051c01905b8181106105a3575050565b5f8155600101610598565b909150819061058f565b600211156105c257565b634e487b7160e01b5f52602160045260245ffd5b9060028110156105c25769ff00000000000000000082549160481b169069ff0000000000000000001916179055565b80518051906001600160401b03821161021557610646826106335f516020616be45f395f51905f5254610531565b5f516020616be45f395f51905f52610569565b602090601f831160011461086d579261067f836108269460e094610283975f92610862575b50508160011b915f199060031b1c19161790565b5f516020616be45f395f51905f52555b60208181015180517f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170180549284015161ff0060089190911b1660ff90921661ffff199093169290921717905560408083015180517f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1702805492909401516fffffffffffffffff0000000000000000931b929092166001600160401b039092166001600160801b031990911617179055606081015161076f9063ffffffff165f516020616b445f395f51905f529063ffffffff1663ffffffff19825416179055565b6107b3610783608083015163ffffffff1690565b5f516020616b445f395f51905f529067ffffffff0000000082549160201b169067ffffffff000000001916179055565b6107f76107c360a0830151151590565b5f516020616b445f395f51905f529068ff0000000000000000825491151560401b169068ff00000000000000001916179055565b61081b60c0820151610808816105b8565b5f516020616b445f395f51905f526105d6565b015163ffffffff1690565b5f516020616b445f395f51905f52906dffffffff0000000000000000000082549160501b16906dffffffff000000000000000000001916179055565b015190505f8061066b565b5f516020616be45f395f51905f525f52601f19831691907fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be925f5b8181106108f957509360e09361028396936001938361082698106108e1575b505050811b015f516020616be45f395f51905f525561068f565b01515f1960f88460031b161c191690555f80806108c7565b929360206001819287860151815501950193016108a8565b1561091a575050565b6355bace6f60e01b5f9081526001600160401b039182166004529116602452604490fd5b634e487b7160e01b5f52601160045260245ffd5b9063ffffffff8091169116019063ffffffff821161096c57565b61093e565b1561097a575050565b9063ffffffff80926333fae18560e21b5f52166004521660245260445ffd5b94610ad6610aef96610aea9694610adb946109e26020987f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a6123360e052898082518301019101610463565b926109ec84610605565b610a08896109fa8651610c0e565b01516001600160401b031690565b60408501805151909991610a3a916001600160401b0316906001600160401b0382166001600160401b03821614610911565b88518a01516001600160401b03165f9081527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1704602052604090205560805260a05260c05260608101610ad0610ab86080610aac610a9b855163ffffffff1690565b60e087015163ffffffff1690610952565b94015163ffffffff1690565b9263ffffffff80851691161115915163ffffffff1690565b90610971565b610df7565b5101516001600160401b031690565b6110e5565b6001600160a01b038116610b155750610b06611324565b50610b1260e0516113a6565b50565b80610b22610b1292611226565b50610b2c8161129c565b5060e0516113ff565b60405190610b4282610235565b5f602083606081520152565b801561096c575f190190565b5f1981019190821161096c57565b9190820391821161096c57565b634e487b7160e01b5f52603260045260245ffd5b908151811015610b9a570160200190565b610b75565b906001820180921161096c57565b603001908160301161096c57565b906004820180921161096c57565b90600c820180921161096c57565b90602c820180921161096c57565b600101908160011161096c57565b602401908160241161096c57565b9190820180921161096c57565b610c16610b35565b5080518015908115610deb575b50610ddc575f19908051805b610d8c575b505f198214610d6e57600360fc1b6001600160f81b0319610c6e610c60610c5a86610b9f565b85610b89565b516001600160f81b03191690565b161480610d78575b610d6e575f90610c8583610b9f565b915b8151831015610d1c57610ca6610ca0610c608585610b89565b60f81c90565b60ff811660308110908115610d11575b50610d0457600a82026001600160401b03908116602f1990920160ff1691909101811691168110610cec57600190920191610c87565b50915050610cf8610285565b9081525f602082015290565b5050915050610cf8610285565b60399150115f610cb6565b9150916001811190811591610d62575b50610d5357610d5090610d3d610285565b9283526001600160401b03166020830152565b90565b63384c687760e21b5f5260045ffd5b602b915010155f610d2c565b9050610cf8610285565b506002610d86838351610b68565b11610c76565b602d60f81b610db6610da9610c60610da385610b5a565b86610b89565b6001600160f81b03191690565b14610dca57610dc490610b4e565b80610c2f565b610dd5919250610b5a565b905f610c34565b6329120bff60e21b5f5260045ffd5b6040915010155f610c23565b610e008161148c565b9051805190811561103a5760b48211611021575f93610e41610e31610e2c610e2786611527565b610bad565b61155f565b93610e3b8561198c565b846119d4565b610e5c610e55610e5086611b87565b61189d565b84602e0152565b610e65836119e4565b60305f955b8351871015610f1557610f0a849392610f05600193610eee610e8d8c809a611478565b5191610ecd63ffffffff610ec46040860193610ebe610eb2865160018060401b031690565b6001600160401b031690565b90610c01565b9b16868d6119b0565b610ee7610ed986610bbb565b91516001600160401b031690565b908b611a38565b6020610ef984610bc9565b91015190890160200152565b610bd7565b960195909192610e6a565b91955061028394610ff4946020945092610fa59250610f5090610f416001600160401b03821115611591565b6001600160401b0316846119f2565b610f94610f66610f6085846115bc565b94611707565b5f516020616b845f395f51905f5280546001600160a01b0319166001600160a01b0392909216919091179055565b5f516020616bc45f395f51905f5255565b8051610feb906001600160401b03165f516020616b845f395f51905f528054600160a01b600160e01b03191660a09290921b600160a01b600160e01b0316919091179055565b015161ffff1690565b5f516020616b845f395f51905f52805461ffff60e01b191660e09290921b61ffff60e01b16919091179055565b63156f758160e31b5f52600482905260b460245260445ffd5b6305f8ded760e21b5f5260045ffd5b5f516020616b045f395f51905f525490680100000000000000008210156102155760018201805f516020616b045f395f51905f5255821015610b9a575f516020616b045f395f51905f525f52600282901c7f06ff75c3bdc3f03474569b3e952aab1bfb9d41b303b5638523d3e1200dcd1ae20180546001600160401b0360069490941b60c01684811b199091169390921690911b919091179055565b6001600160401b0381165f9081525f516020616c045f395f51905f52602052604090206001600160a01b03906001015416156112165f516020616bc45f395f51905f52546111925f516020616b845f395f51905f525460018060a01b0381169061ffff60018060401b038260a01c169160e01c169160405194611169608087610250565b85526001600160a01b031660208501526001600160401b0316604084015261ffff166060830152565b6001600160401b0384165f9081525f516020616c045f395f51905f526020526040902081518155602082015160019190910180546040840151606094909401516001600160f01b03199091166001600160a01b03939093169290921760a09390931b600160a01b600160e01b03169290921760e09190911b61ffff60e01b16179055565b61121d5750565b61028390611049565b6001600160a01b0381165f9081525f516020616c245f395f51905f52602052604090205460ff16611297576001600160a01b03165f8181525f516020616c245f395f51905f5260205260408120805460ff191660011790553391905f516020616ae45f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f516020616b245f395f51905f52602052604090205460ff16611297576001600160a01b0381165f9081525f516020616b245f395f51905f5260205260409020805460ff1916600117905533906001600160a01b03165f516020616ba45f395f51905f525f516020616ae45f395f51905f525f80a4600190565b5f80525f516020616b245f395f51905f526020525f516020616b645f395f51905f525460ff166113a2575f8080525f516020616b245f395f51905f526020525f516020616b645f395f51905f52805460ff1916600117905533905f516020616ba45f395f51905f525f516020616ae45f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f14611297575f818152602081815260408083208380529091528120805460ff1916600117905533915f516020616ae45f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16611472575f818152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905533916001600160a01b0316905f516020616ae45f395f51905f525f80a4600190565b50505f90565b8051821015610b9a5760209160051b010190565b8051519081156114725761149f8261032d565b916114ad6040519384610250565b808352601f196114bc8261032d565b013660208501375f5b825180518210156115195790611508610e5060206114e584600196611478565b5101516115036114fb6040610adb878b51611478565b610d3d610285565b6117b8565b6115128287611478565b52016114c5565b505090505f610d50926118f2565b90602c820291808304602c149015171561096c57565b6040516080919061154e8382610250565b6041815291601f1901366020840137565b90611569826102ce565b6115766040519182610250565b8281528092611587601f19916102ce565b0190602036910137565b1561103a57565b604051906115a58261021a565b5f6060838281528260208201528260408201520152565b9190916115c7611598565b50603083511061167c576356414c34602084015160e01c0361167c57602483015160c01c602c84015160f01c93611652602e82015192604e83015160f01c93611620611611610294565b6001600160401b039093168352565b6116326020830198899061ffff169052565b60408201526116496060820194859061ffff169052565b955161ffff1690565b9061ffff82169283159384156116b2575b5083156116a3575b50821561168e575b505061167c5750565b634724a0fd60e01b5f5260045260245ffd5b51915061169a90611a80565b14155f80611673565b5161ffff16151592505f61166b565b60b41093505f611663565b805191908290602001825e015f815290565b606160f81b81526001600160f01b031990911660018201526680600a3d393df360c81b6003820152610d5091600a91909101906116bd565b604051905f60208301526117308261172260218201846116bd565b03601f198101845283610250565b6160008251116117a55750611772611780611760611750845161ffff1690565b60f01b6001600160f01b03191690565b926040519283916020830195866116cf565b03601f198101835282610250565b51905ff0906001600160a01b0382161561179657565b63fbad885d60e01b5f5260045ffd5b516331734c2160e21b5f5260045260245ffd5b6020810180516024906001600160401b031680611857575b506117dd61180b9161155f565b926118026117fc6117f66117f087611aa5565b87611ab1565b86611ac8565b85611adf565b90519084611b0c565b8151909190611822906001600160401b0316610eb2565b61182b57505090565b61184c610eb261183e6118539486611af5565b92516001600160401b031690565b9083611b24565b5090565b600191505b608081101561188457506117dd61187d61187861180b93610be5565b610bf3565b91506117d0565b60019060071c91019061185c565b6040513d5f823e3d90fd5b80516001810180911161096c576118b39061155f565b805115610b9a576118dd816118d06020945f868196015382611b5d565b50604051918280926116bd565b039060025afa156118ed575f5190565b611892565b909291928084039380851161096c576001851461197b5760015b8060011b908682101561191f575061190c565b9192939495505082019182811161096c578261193b91856118f2565b9161194692936118f2565b61194e61153d565b91825115610b9a57825f926118dd92600160208097015360218301526041820152604051918280926116bd565b5090611988929350611478565b5190565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602d908260081c602c8201530153565b604f5f9182604e8201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b61ffff16602c810290808204602c149015171561096c576030018060301161096c5790565b6020600a910153600190565b6020826022920101536001810180911161096c5790565b602082600a920101536001810180911161096c5790565b60208281920101536001810180911161096c5790565b6020826010920101536001810180911161096c5790565b816020919392930101526020810180911161096c5790565b9092919083016020015b6080821015611b4257906001929391530190565b600180916080607f85161781530193019060071c9092611b2e565b908051918215611b7f576021602084930191015e6001018060011161096c5790565b505050600190565b60405190606090611b988284610250565b6022835261185391600a906020850190601f19013682375360206021840153600283611b0c56fe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a714610134578063248a9ca31461012f5780632f2ff15d1461012a57806336568abe146101255780634b1872be146101205780636a28f0001461011b5780636edfe3af146101165780638a8e4c5d146101115780638c80fbda1461010c57806391d148541461010757806392c19bdc14610102578063974a74c4146100fd578063a217fddf146100f8578063a6f031bb146100f3578063d547741f146100ee578063db3e1fa4146100e9578063ddba6537146100e45763ef913a4b146100df575f80fd5b610fcb565b610e02565b610dc8565b610d99565b610c80565b610c66565b610afa565b610a69565b610a28565b6109ec565b6109b4565b6108ac565b610702565b61033b565b610267565b610231565b6101d9565b346101d55760203660031901126101d5576004357fffffffff0000000000000000000000000000000000000000000000000000000081168091036101d557807f7965db0b00000000000000000000000000000000000000000000000000000000602092149081156101ab575b506040519015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150145f6101a0565b5f80fd5b346101d55760203660031901126101d55760206102036004355f525f602052600160405f20015490565b604051908152f35b60409060031901126101d557600435906024356001600160a01b03811681036101d55790565b346101d5576102656102423661020b565b9061026061025b825f525f602052600160405f20015490565b612a5b565b612aa2565b005b346101d5576102753661020b565b336001600160a01b0382160361028e5761026591612b3a565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b9060206003198301126101d5576004356001600160401b0381116101d557826023820112156101d5578060040135926001600160401b0384116101d557602484830101116101d5576024019190565b634e487b7160e01b5f52602160045260245ffd5b6003111561032357565b610305565b9190602083019261033882610319565b52565b346101d55761067261047a61046b6104376104456103d361035b366102b6565b61037960ff5f516020614ef85f395f51905f525460401c16156110f2565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475416156106eb575b810190611a44565b6103ef6103ea60408301516001600160801b031690565b612c28565b6104026060820151602083015190612c9f565b6040519283917fb615f64400000000000000000000000000000000000000000000000000000000602084015260248301611e6e565b03601f19810183528261118b565b7f0000000000000000000000000000000000000000000000000000000000000000612d22565b60208082518301019101611f56565b61048381612d70565b9061048d82610319565b81610676576106626106486020838160607f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb319601916105216104d983855101516001600160401b031690565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025460401c6001600160401b03166001600160401b0381166001600160401b03831611611fb4565b6105d683516001600160401b038151166fffffffffffffffffffffffffffffffff196fffffffffffffffff000000000000000060207f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025494846001600160401b03198716177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170255015160401b16921617177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170255565b01516040516105ec816104378682019485611ff8565b5190208151830151610638906001600160401b03165b6001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b555101516001600160401b031690565b6040516001600160401b0390911681529081906020820190565b0390a15b60405191829182610328565b0390f35b5061068081610319565b60018103610666576106c26801000000000000000068ff0000000000000000195f516020614ef85f395f51905f525416175f516020614ef85f395f51905f5255565b7f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1610666565b6106f36129ec565b6103cb565b5f9103126101d557565b346101d5575f3660031901126101d557335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff16156107cc575f516020614ef85f395f51905f525460ff8160401c16156107a45768ff000000000000000019165f516020614ef85f395f51905f52557f8e0da210a2ffecc744373b4860cd9b8dc471810f5fb61b7d2c3234ccd4aae30d5f80a1005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b90602080835192838152019201905f5b8181106108005750505090565b82518452602093840193909201916001016107f3565b919060608301916060845281518093526020608085019201925f5b81811061089057505061084c925083820360208501526107e3565b906040818303910152602080835192838152019201905f5b8181106108715750505090565b82516001600160401b0316845260209384019390920191600101610864565b845163ffffffff16845260209485019490930192600101610831565b346101d5575f3660031901126101d55760206108c6612e76565b6108cf81612f02565b9201916108f06108eb6108e4855161ffff1690565b61ffff1690565b612026565b906109036108eb6108e4865161ffff1690565b936109166108eb6108e4835161ffff1690565b915f5b866109296108e4855161ffff1690565b61ffff8316908110156109a2579161099b60019261098d858061098661097e828f8f6109788f928f9261ffff9f6109638161096e9361206c565b9063ffffffff169052565b51925161ffff1690565b91613064565b92909561206c565b528961206c565b906001600160401b03169052565b0116610919565b50856106728660405193849384610816565b346101d5576109c2366102b6565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b346101d5575f3660031901126101d55760207f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170554604051908152f35b346101d557602060ff610a5d610a3d3661020b565b905f525f845260405f20906001600160a01b03165f5260205260405f2090565b54166040519015158152f35b346101d5576020610adc610a7c366102b6565b90610a9b60ff5f516020614ef85f395f51905f525460401c16156110f2565b7fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5f525f845260405f205f8052845260ff60405f20541615610aed57612280565b60405190610ae981610319565b8152f35b610af56129ec565b612280565b346101d55760203660031901126101d5576004356001600160401b0381116101d557806004019061016060031982360301126101d557610b4e60ff5f516020614ef85f395f51905f525460401c16156110f2565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610c59575b610144810190610bb082846125c5565b905015610c3157610c12610c219261067294610bcf60448501826125f7565b610bdf60648794939401836125f7565b9060848801359261010489013595610bf68761262c565b60a4610c19610c096101248d01896125f7565b9b909a896125c5565b36916112b4565b9a019561377f565b6040519081529081906020820190565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b610c616129ec565b610ba0565b346101d5575f3660031901126101d55760206040515f8152f35b346101d55760203660031901126101d5576004356001600160401b0381116101d5578060040161014060031983360301126101d55761067291610c2191610cdb60ff5f516020614ef85f395f51905f525460401c16156110f2565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610d8c575b610d3a60448301826125f7565b91610d4860648501826125f7565b610104860135929160848701359190610d608561262c565b610d6e6101248901856125f7565b97909660a46040519a610d8260208d61118b565b5f8c52019561377f565b610d946129ec565b610d2d565b346101d557610265610daa3661020b565b90610dc361025b825f525f602052600160405f20015490565b612b3a565b346101d5575f3660031901126101d55760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346101d557610f23610437610efd610e89610e1c366102b6565b610e3a60ff5f516020614ef85f395f51905f525460401c16156110f2565b7f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260ff610e7860405f205f805260205260405f2090565b541615610f84575b508101906126fc565b610ea06103ea60608301516001600160801b031690565b610eb4608082015160208351015190613993565b610ec860a082015160408351015190613993565b6040519283917f789d9f39000000000000000000000000000000000000000000000000000000006020840152602483016127ca565b7f0000000000000000000000000000000000000000000000000000000000000000612d22565b50610f5e6801000000000000000068ff0000000000000000195f516020614ef85f395f51905f525416175f516020614ef85f395f51905f5255565b7f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1005b610f8d90612a5b565b5f610e80565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b906020610fc8928181520190610f93565b90565b346101d5575f3660031901126101d55761067260405160208082015261014060408201526110e68161100061018082016128eb565b61103b60608301602060ff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170154818116845260081c16910152565b61107c60a0830160206001600160401b037f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170254818116845260401c16910152565b5f516020614ef85f395f51905f525463ffffffff811660e084015261043790602081901c63ffffffff166101008501526110c1610120850160ff8360401c1615159052565b6110d5610140850160ff8360481c166129e2565b60501c63ffffffff16610160840152565b60405191829182610fb7565b156110f957565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52604160045260245ffd5b606081019081106001600160401b0382111761115057604052565b611121565b604081019081106001600160401b0382111761115057604052565b608081019081106001600160401b0382111761115057604052565b90601f801991011681019081106001600160401b0382111761115057604052565b604051906111bc6102608361118b565b565b604051906111bc6101008361118b565b604051906111bc60c08361118b565b604051906111bc60808361118b565b6001600160801b038116036101d557565b35906111bc826111ec565b91908260609103126101d55760405161122081611135565b60408082948035611230816111ec565b8452602081013560208501520135910152565b6001600160401b038116036101d557565b35906111bc82611243565b91908260409103126101d55760405161127781611155565b6020808294803561128781611243565b845201359161129583611243565b0152565b6001600160401b03811161115057601f01601f191660200190565b9291926112c082611299565b916112ce604051938461118b565b8294818452818301116101d5578281602093845f960137010152565b9080601f830112156101d557816020610fc8933591016112b4565b359081151582036101d557565b359063ffffffff821682036101d557565b8092910391606083126101d55760405161133c81611155565b6040819483358352601f1901126101d557602090604080519361135e85611155565b611369848201611312565b85520135828401520152565b6001600160401b0381116111505760051b60200190565b919060c0838203126101d5576040516113a481611170565b809380356113b181611243565b82526113bf60208201611312565b60208301526113d18360408301611323565b604083015260a0810135906001600160401b0382116101d557019180601f840112156101d55782359261140384611375565b93611411604051958661118b565b80855260208086019160051b830101918383116101d55760208101915b83831061144057505050505060600152565b82356001600160401b0381116101d5578201906040828703601f1901126101d5576040519161146e83611155565b602081013560048110156101d557835260408101356001600160401b0381116101d5576020910101906080828803126101d557604051926114ae84611170565b82356001600160401b0381116101d557886114ca9185016112ea565b845260208301356114da816111ec565b60208501526114eb60408401611305565b60408501526060830135936001600160401b0385116101d557611513896020968796016112ea565b60608201528382015281520192019161142e565b91906060838203126101d5576040519061154082611155565b819380356001600160401b0381116101d5578101916040838203126101d5576040519261156c84611155565b80356001600160401b0381116101d55781016102c0818403126101d5576115916111ac565b9061159c848261125f565b825260408101356001600160401b0381116101d557846115bd9183016112ea565b60208301526115ce60608201611254565b60408301526115df608082016111fd565b60608301526115f060a08201611305565b60808301526116028460c08301611323565b60a08301526116146101208201611305565b60c083015261014081013560e08301526116316101608201611305565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526116806102208201611305565b6101c08301526102408101356101e083015261169f6102608201611305565b6102008301526102808101356102208301526102a0810135906001600160401b0382116101d5576116d2918591016112ea565b61024082015284526020810135926001600160401b0384116101d5576020946117018461170d9688950161138c565b8382015286520161125f565b910152565b9080601f830112156101d5576101006040519261172f828561118b565b839181019283116101d557905b8282106117495750505090565b813581526020918201910161173c565b9080601f830112156101d5576040519161177460408461118b565b8290604081019283116101d557905b8282106117905750505090565b8135815260209182019101611783565b359061ffff821682036101d557565b9080601f830112156101d55781356117c681611375565b926117d4604051948561118b565b81845260208085019260051b8201019283116101d557602001905b8282106117fc5750505090565b6020809161180984611312565b8152019101906117ef565b9080601f830112156101d557813561182b81611375565b92611839604051948561118b565b81845260208085019260051b8201019283116101d557602001905b8282106118615750505090565b8135815260209182019101611854565b9080601f830112156101d557813561188881611375565b92611896604051948561118b565b81845260208085019260051b8201019283116101d557602001905b8282106118be5750505090565b602080916118cb84611305565b8152019101906118b1565b919091610220818403126101d5576118ec6111be565b926118f78183611712565b8452611907816101008401611759565b602085015261191a816101408401611759565b604085015261192c61018083016117a0565b60608501526101a08201356001600160401b0381116101d557816119519184016117af565b60808501526101c08201356001600160401b0381116101d557816119769184016117af565b60a08501526101e08201356001600160401b0381116101d5578161199b918401611814565b60c08501526102008201356001600160401b0381116101d5576119be9201611871565b60e0830152565b919060c0838203126101d557604051906119de82611170565b81936119ea8282611208565b835260608101356001600160401b0381116101d55782611a0b918301611527565b60208401526080810135611a1e816111ec565b604084015260a0810135916001600160401b0383116101d55760609261170d92016118d6565b906020828203126101d55781356001600160401b0381116101d557610fc892016119c5565b90606060c08201926001600160401b03815116835263ffffffff6020820151166020840152611aba6040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519160c060a0830152825180915260e0820190602060e08260051b8501019401925f905b828210611aee57505050505090565b909192939460df198282030185528551908151600481101561032357611b6c8260206001958195948295520151906040848201526060611b3a83516080604085015260c0840190610f93565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610f93565b9701950193920190611adf565b90610fc890602080611cf285516060855282611cde825160406060890152611bba60a0890182516001600160401b0360208092828151168552015116910152565b610240611bd7848301516102c060e08c01526103608b0190610f93565b60408301516001600160401b03166101008b01529160608101516001600160801b03166101208b0152608081015115156101408b015260a081015180516101608c0152602090810151805163ffffffff166101808d015201516101a08b015260c081015115156101c08b015260e08101516101e08b015261010081015115156102008b01526101208101516102208b0152610140810151828b01526101608101516102608b01526101808101516102808b01526101a08101516102a08b01526101c081015115156102c08b01526101e08101516102e08b015261020081015115156103008b01526102208101516103208b01520151888203609f19016103408a0152610f93565b910151858203605f19016080870152611a69565b9401519101906001600160401b0360208092828151168552015116910152565b905f905b60088210611d2357505050565b6020806001928551815201930191019091611d16565b905f905b60028210611d4a57505050565b6020806001928551815201930191019091611d3d565b90602080835192838152019201905f5b818110611d7d5750505090565b825163ffffffff16845260209384019390920191600101611d70565b90602080835192838152019201905f5b818110611db65750505090565b82511515845260209384019390920191600101611da9565b610fc891611ddd818351611d12565b611df06020830151610100830190611d39565b611e036040830151610140830190611d39565b606082015161ffff1661018082015260e0611e5c611e49611e3660808601516102206101a0870152610220860190611d60565b60a08601518582036101c0870152611d60565b60c08501518482036101e08601526107e3565b92015190610200818403910152611d99565b90610fc89160208152611ea2602082018351604080916001600160801b038151168452602081015160208501520151910152565b6060611ebd602084015160c0608085015260e0840190611b79565b926001600160801b0360408201511660a084015201519060c0601f1982850301910152611dce565b91908260609103126101d557604051611efd81611135565b60408082948051611f0d816111ec565b8452602081015160208501520151910152565b91908260409103126101d557604051611f3881611155565b60208082948051611f4881611243565b845201519161129583611243565b90610140828203126101d557611fac9061010060405193611f7685611170565b611f808382611ee5565b8552611f8f8360608301611ee5565b6020860152611fa18360c08301611f20565b604086015201611f20565b606082015290565b15611fbd575050565b906001600160401b0380927f7d939bd8000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6111bc909291926060810193604080916001600160801b038151168452602081015160208501520151910152565b9061203082611375565b61203d604051918261118b565b828152809261204e601f1991611375565b0190602036910137565b634e487b7160e01b5f52603260045260245ffd5b80518210156120805760209160051b010190565b612058565b91906080838203126101d5576040519061209e82611170565b819380356001600160401b0381116101d5576060926120be9183016112ea565b83526020810135602084015260408101356120d881611243565b60408401520135908160070b82036101d55760600152565b6020818303126101d5578035906001600160401b0382116101d55701906040828203126101d5576040519161212483611155565b80356001600160401b0381116101d557826121409183016119c5565b83526020810135906001600160401b0382116101d557016080818303126101d5576040519161216e83611170565b81356001600160401b0381116101d557820181601f820112156101d557803561219681611375565b916121a4604051938461118b565b81835260208084019260051b820101918483116101d55760208201905b83821061221c575050505083526121da60208301611305565b60208401526040820135916001600160401b0383116101d55761220460609261220f948301612085565b604085015201611254565b6060820152602082015290565b81356001600160401b0381116101d55760209161223e88848094880101612085565b8152019101906121c1565b15612252575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b61228c918101906120f0565b8051906122a66103ea60408401516001600160801b031690565b6122c861046b6104376104456060860151956104026020820197885190612c9f565b906122d282612d70565b926122dc84610319565b60018414612561577f1411889dcdfcf9af1fcf4ff4521cab4580a1fd842401ead543815cb1148d64389261232d602061016094019261231b8451613125565b90515151909401518490808214612249565b61233685610319565b846124ce576124b0916020826123686060839501936123626104d985875101516001600160401b031690565b516131c2565b61241d83516001600160401b038151166fffffffffffffffffffffffffffffffff196fffffffffffffffff000000000000000060207f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025494846001600160401b03198716177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170255015160401b16921617177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170255565b0151604051612433816104378682019485611ff8565b519020815183015161244d906001600160401b0316610602565b558051820151612465906001600160401b03166135bb565b7f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb3161249e610648848451016001600160401b0390511690565b0390a15101516001600160401b031690565b604080516001600160401b039290921682526020820192909252a190565b6124b09161253b606060209301916123626124f385855101516001600160401b031690565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025460401c6001600160401b03166001600160401b0381166001600160401b03831614611fb4565b8051820151612552906001600160401b03166135bb565b5101516001600160401b031690565b50505061259e6801000000000000000068ff0000000000000000195f516020614ef85f395f51905f525416175f516020614ef85f395f51905f5255565b7f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a190565b903590601e19813603018212156101d557018035906001600160401b0382116101d5576020019181360383136101d557565b903590601e19813603018212156101d557018035906001600160401b0382116101d557602001918160051b360383136101d557565b600211156101d557565b6002111561032357565b91906060838203126101d5576040519061265982611135565b819380356001600160401b0381116101d55781016040818403126101d5576040519061268482611155565b8035906001600160401b0382116101d5576126a38560209383016112ea565b835201356126b081611243565b6020820152835260208101356001600160401b0381116101d557826126d6918301611527565b60208401526040810135916001600160401b0383116101d55760409261170d9201611527565b6020818303126101d5578035906001600160401b0382116101d55701610140818303126101d55761272b6111ce565b9181356001600160401b0381116101d55781612748918401612640565b83526127578160208401611208565b60208401526127698160808401611208565b604084015261277a60e083016111fd565b60608401526101008201356001600160401b0381116101d5578161279f9184016118d6565b60808401526101208201356001600160401b0381116101d5576127c292016118d6565b60a082015290565b90610fc8916020815260a06128d561284d845161014060208601526040612837825160606101608901526001600160401b0360206128158351866101c08d01526102008c0190610f93565b920151166101e0890152602084015188820361015f19016101808a0152611b79565b91015185820361015f19016101a0870152611b79565b61287c60208601516040860190604080916001600160801b038151168452602081015160208501520151910152565b6128aa604086015184860190604080916001600160801b038151168452602081015160208501520151910152565b60608501516001600160801b03166101008501526080850151848203601f1901610120860152611dce565b92015190610140601f1982850301910152611dce565b905f917f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170054908160011c91600181169182156129d8575b6020841083146129c457838152602001919081156129b25750600114612946575050565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005f908152929350907fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be5b81841061299e5750500190565b805484840152602090930192600101612991565b60ff191682525090151560051b019150565b634e487b7160e01b5f52602260045260245ffd5b92607f1692612922565b9061033882612636565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff1615612a2457565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260ff612a823360405f20906001600160a01b03165f5260205260405f2090565b541615612a8c5750565b63e2517d3f60e01b5f523360045260245260445ffd5b805f525f60205260ff612ac98360405f20906001600160a01b03165f5260205260405f2090565b5416612b3457805f525f602052612af48260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916600117905533916001600160a01b0316907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260ff612b618360405f20906001600160a01b03165f5260205260405f2090565b541615612b3457805f525f602052612b8d8260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916905533916001600160a01b0316907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b15612bd3575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b634e487b7160e01b5f52601160045260245ffd5b5f19810191908211612c2357565b612c01565b6001600160801b03633b9aca00911604612c46814242821115612bca565b804203428111612c235763ffffffff5f516020614ef85f395f51905f525460501c1610612c705750565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b606060206111bc9351015101517f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170091612cd6612e52565b5061ffff6006600585015494015460405194612cf186611170565b85526001600160a01b03811660208601526001600160401b038160a01c16604086015260e01c166060840152613c6a565b5f918291602082519201905af43d15612d68573d90612d4082611299565b91612d4e604051938461118b565b82523d5f602084013e5b15612d605790565b602081519101fd5b606090612d58565b612dc2612d8b6020606084015101516001600160401b031690565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1704906001600160401b03165f5260205260405f2090565b5480612dce5750505f90565b60208201908151604051612dea81610437602082019485611ff8565b5190201491821592612e08575b505015612e0357600190565b600290565b6001600160801b03919250612e3b612e2c612e4792516001600160801b0390511690565b9351516001600160801b031690565b6001600160801b031690565b911610155f80612df7565b60405190612e5f82611170565b5f6060838281528260208201528260408201520152565b612e7e612e52565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055461ffff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17065460405192612ed384611170565b83526001600160a01b03811660208401526001600160401b038160a01c16604084015260e01c16606082015290565b612f0a612e52565b50602081016001600160a01b0381511615612f9f576001600160a01b03905116612f3b61ffff606084015116613e9a565b813b6001811115612f8c575f198101908111612c23578111612f795790816001612f67612f7694613ee1565b9260208401903c809251613f09565b91565b50633610565160e21b5f5260045260245ffd5b82633610565160e21b5f5260045260245ffd5b5051633a517eed60e21b5f5260045260245ffd5b90602c820291808304602c1490151715612c2357565b90600382029180830460031490151715612c2357565b908160011b9180830460021490151715612c2357565b6030019081603011612c2357565b9060048201809211612c2357565b90600c8201809211612c2357565b90602c8201809211612c2357565b9060018201809211612c2357565b6001019081600111612c2357565b6024019081602411612c2357565b91908201809211612c2357565b9091939263ffffffff61ffff911694168410156130f55761308c61308785612fb3565b612ff5565b9363ffffffff6020868501015160e01c16036130ca5750610fc8906130c260206130b586613003565b8301015160c01c94613011565b016020015190565b7f4724a0fd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b83907ff81f5140000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b805151908115612b345761313882612026565b915f5b825180518210156131b457906131a361319e602061315b8460019661206c565b5101516131996001600160401b036040613176878b5161206c565b510151166040519261318784611155565b83526001600160401b03166020830152565b614003565b6140f8565b6131ad828761206c565b520161313b565b505090505f610fc892614148565b6131cb81613125565b905180519081156134d15760b4821161349e575f936132076131f76131f261308786612fb3565b613ee1565b9361320185614bb0565b84614bf8565b61321d61321661319e86614e3f565b84602e0152565b61322683614c08565b60305f955b83518710156132d7576132cc8493926132c76001936132b061324e8c809a61206c565b519161328f63ffffffff613286604086019361328061327486516001600160401b031690565b6001600160401b031690565b90613057565b9b16868d614bd4565b6132a961329b86613003565b91516001600160401b031690565b908b614c5c565b60206132bb84613011565b91015190890160200152565b61301f565b96019590919261322b565b61344394929650602093506133ab91509461330a6001600160401b036111bc9761330382821115613b51565b1684614c16565b61338761332061331a8584613f09565b94614261565b6001600160a01b031673ffffffffffffffffffffffffffffffffffffffff197f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17065416177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170555565b61343a6133bf82516001600160401b031690565b7fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706549260a01b169116177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655565b015161ffff1690565b61ffff60e01b1961ffff60e01b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706549260e01b169116177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655565b7fab7bac08000000000000000000000000000000000000000000000000000000005f52600482905260b460245260445b5ffd5b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b9190918054831015612080575f52601860205f208360021c019260031b1690565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170854680100000000000000008110156111505780600161359d92017f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708557f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17086134f9565b6001600160401b0380839493549260031b9316831b921b1916179055565b6001600160a01b0360016135ff836001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f2090565b0154161561376f7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055461369f7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706546001600160401b036001600160a01b038216916001600160a01b0361ffff838360a01c169260e01c16936040519661368660808961118b565b875216602086015216604084015261ffff166060830152565b6136d9846001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f2090565b6060600161ffff928451815501926001600160a01b0360208201511673ffffffffffffffffffffffffffffffffffffffff1985541617845560408101517fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b8087549360a01b1616911617845501511661ffff60e01b1961ffff60e01b83549260e01b169116179055565b6137765750565b6111bc9061351a565b99909396999892979195986137938b612636565b8a156137d4576134ce8b6137a681612636565b7f112d89cc000000000000000000000000000000000000000000000000000000005f5260ff16600452602490565b869798999a506020613806918861381794959697989915159081613986575b613800919897969861434e565b0161438c565b613810368c611208565b9087614ca4565b5f5f935b8785106138e5575b5061382e9350614537565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001691823b156101d55761389b5f95604051978896879586957f56455e7900000000000000000000000000000000000000000000000000000000875260048701614900565b03915afa80156138e057610fc8926138bc92612e3b926138c6575b506149ba565b633b9aca00900490565b806138d45f6138da9361118b565b806106f8565b5f6138b6565b6140ed565b9061391f61391b61390a6139036138fd898d8c614396565b806125f7565b36916143b8565b6139153688886143b8565b90614dd5565b1590565b61397a575061395c90613946610c1261393c61382e978b8a614396565b60208101906125c5565b90815181518082149182613964575b5050614423565b60015f613823565b9091506020840120906020830120145f80613955565b6001909401939061381b565b61ffff81111591506137f3565b906001600160401b036020806060818551015101519301510151166139b6612e52565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708548015613b3b576139e95f91612c15565b808210613acb5750613a8991613a4d613a3b613a28613a84947f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17086134f9565b90546001600160401b039160031b1c1690565b916001600160401b03831611156149c4565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1707906001600160401b03165f5260205260405f2090565b6149de565b916001600160a01b03613aa660208501516001600160a01b031690565b1615613ab757906111bc9291613c6a565b8251633a517eed60e21b5f5260045260245ffd5b90613ae7613ae1613adc8484613057565b61302d565b60011c90565b90836001600160401b03613b1e613a28857f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17086134f9565b1611613b2b5750906139e9565b9150613b3690612c15565b6139e9565b633a517eed60e21b5f526134ce6024905f600452565b156134d157565b15613b605750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15613b9a5750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15613bd45750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b906001600160401b03809116911601906001600160401b038211612c2357565b15613c2f575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b608081019384515193613cac6060840195613c8a6108e4885161ffff1690565b8091149081613e8b575b81613e7c575b81613e6d575b50908693959291613b51565b613cb581612f02565b929095613ccc60408401516001600160401b031690565b905f925f925f945f613ce56108e45f9b5b5161ffff1690565b8a1015613e2157908b949392918e613d0e61391b613d088f8f9060e0015161206c565b51151590565b613e0e57613d208c613d2a925161206c565b5163ffffffff1690565b8098613df0575b50508a600197968b8060a084015190613d499161206c565b5163ffffffff16918260208a0151613d629061ffff1690565b61ffff1663ffffffff82161090613d7891613b92565b600163ffffffff84161b90613d908482841615613b58565b1797828d8d519260200151613da69061ffff1690565b90613db093613064565b91909360c0015190613dc19161206c565b511490613dcd91613bcc565b613dd691613c06565b986108e46001613ce5925b019a929394959691908e613cdd565b63ffffffff613e07921663ffffffff821611613b58565b5f87613d31565b50959450986108e46001613ce592613de1565b509350935097509893506111bc975060e09450613e639250613e4b6001600160401b038216612fc9565b613e5d6001600160401b038416612fdf565b10613c26565b5191015191614a24565b905060e085015151145f613ca0565b60c08601515181149150613c9a565b60a08601515181149150613c94565b61ffff16602c810290808204602c1490151715612c235760300180603011612c235790565b60405160809190613ed0838261118b565b6041815291601f1901366020840137565b90613eeb82611299565b613ef8604051918261118b565b828152809261204e601f1991611299565b919091613f14612e52565b5060308351106130ca576356414c34602084015160e01c036130ca57602483015160c01c602c84015160f01c93613f9f602e82015192604e83015160f01c93613f6d613f5e6111dd565b6001600160401b039093168352565b613f7f6020830198899061ffff169052565b6040820152613f966060820194859061ffff169052565b955161ffff1690565b9061ffff8216928315938415613ff8575b508315613fde575b508215613fc9575b50506130ca5750565b519150613fd590613e9a565b14155f80613fc0565b51909250613fef9061ffff166108e4565b1515915f613fb8565b60b41093505f613fb0565b602460208201906001600160401b03825116806140a0575b5061402861405691613ee1565b9261404d61404761404161403b87614ace565b87614ada565b86614af1565b85614b08565b90519084614b35565b9061406b61327482516001600160401b031690565b61407457505090565b61409561327461408761409c9486614b1e565b92516001600160401b031690565b9083614b4d565b5090565b600191505b60808110156140cd57506140286140c66140c16140569361303b565b613049565b915061401b565b60019060071c9101906140a5565b805191908290602001825e015f815290565b6040513d5f823e3d90fd5b805160018101809111612c235761410e90613ee1565b805115612080576141388161412b6020945f868196015382614b86565b50604051918280926140db565b039060025afa156138e0575f5190565b9092919280840393808511612c2357600185146141d15760015b8060011b90868210156141755750614162565b91929394955050820191828111612c2357826141919185614148565b9161419c9293614148565b6141a4613ebf565b9182511561208057825f9261413892600160208097015360218301526041820152604051918280926140db565b50906141de92935061206c565b5190565b600a907fffff000000000000000000000000000000000000000000000000000000000000610fc894937f610000000000000000000000000000000000000000000000000000000000000083521660018201527f80600a3d393df300000000000000000000000000000000000000000000000000600382015201906140db565b604051905f602083015261428a8261427c60218201846140db565b03601f19810184528361118b565b61600082511161432257506104376142e46142d26142aa845161ffff1690565b60f01b7fffff0000000000000000000000000000000000000000000000000000000000001690565b926040519283916020830195866141e2565b51905ff0906001600160a01b038216156142fa57565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b156143565750565b7fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b35610fc881611243565b91908110156120805760051b81013590603e19813603018212156101d5570190565b9291906143c481611375565b936143d2604051958661118b565b602085838152019160051b8101918383116101d55781905b8382106143f8575050505050565b81356001600160401b0381116101d55760209161441887849387016112ea565b8152019101906143ea565b9190911561442f575050565b90614483614471926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610f93565b83810360031901602485015290610f93565b0390fd5b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156101d55701602081359101916001600160401b0382116101d55781360383136101d557565b90602083828152019260208260051b82010193835f925b8484106144ff5750505050505090565b909192939495602080614527600193601f198682030188526145218b886144a7565b90614487565b98019401940192949391906144ef565b91909115614543575050565b6144836040519283927ffef760c70000000000000000000000000000000000000000000000000000000084526020600485015260248401916144d8565b9035601e19823603018112156101d55701602081359101916001600160401b0382116101d5578160051b360383136101d557565b9035607e19823603018112156101d5570190565b359060038210156101d557565b9035605e19823603018112156101d5570190565b90602083828152019260208260051b82010193835f925b8484106146105750505050505090565b909192939495602080614682600193601f198682030188526146328b886145d5565b9061463c826145c8565b61464581610319565b8152614674614669614659868501856144a7565b6060888601526060850191614487565b9260408101906144a7565b916040818503910152614487565b9801940194019294939190614600565b610fc89161475c6147516146d56146ba6146ac86806144a7565b608087526080870191614487565b6146c760208701876144a7565b908683036020880152614487565b60806147416146e760408801886145b4565b86840360408801526146f8816145c8565b61470181610319565b845261470f602082016145c8565b61471881610319565b6020850152614729604082016145c8565b61473281610319565b604085015260608101906144a7565b9190928160608201520191614487565b926060810190614580565b9160608185039101526145e9565b90602083828152019260208260051b82010193835f925b8484106147915750505050505090565b909192939495601f198282030184528635601e19843603018112156101d5578301906147c1602082019280614580565b8091936020845252604082019060408160051b8401019380935f915b838310614800575050505050506020806001929801940194019294939190614781565b909192939495603f1983820301865261481987836145d5565b80356148248161262c565b61482d81612636565b825261485061483f60208301836145b4565b606060208501526060840190614692565b906040810135609e19823603018112156101d55760019360209384936148f293019160408183039101526148e46148c661489b61488d85806144a7565b60a0865260a0860191614487565b6148a6878601611305565b1515878501526148b960408601866145b4565b8482036040860152614692565b926148d360608201611305565b1515606084015260808101906145b4565b906080818403910152614692565b9801960194930191906147dd565b9092809492969593606083019083526060602084015252608081019560808560051b83010194815f90603e1981360301995b838310614951575050505050610fc8949550604081850391015261476a565b9091929397607f1986820301825288358b8112156101d55760206149ab600193868394019061499e6149946149868480614580565b6040855260408501916144d8565b92858101906144a7565b9185818503910152614487565b9a019201930191909392614932565b35610fc8816111ec565b156149cb57565b633a517eed60e21b5f525f60045260245ffd5b906040516149eb81611170565b606061ffff600183958054855201546001600160a01b03811660208501526001600160401b038160a01c16604085015260e01c16910152565b9291614a338251825114613b51565b5f5b8251811015614ac757614a48818361206c565b5115614abf5763ffffffff614a5d828561206c565b5116614a6c8187518110613b92565b614a76818761206c565b515160048110156103235760011901614a9457506001905b01614a35565b7f5b5c80eb000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b600190614a8e565b5050509050565b6020600a910153600190565b60208260229201015360018101809111612c235790565b602082600a9201015360018101809111612c235790565b602082819201015360018101809111612c235790565b60208260109201015360018101809111612c235790565b8160209193929301015260208101809111612c235790565b9092919083016020015b6080821015614b6b57906001929391530190565b600180916080607f85161781530193019060071c9092614b57565b908051918215614ba8576021602084930191015e60010180600111612c235790565b505050600190565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602d908260081c602c8201530153565b604f5f9182604e8201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b90614d206040516020810190614cd78287604080916001600160801b038151168452602081015160208501520151910152565b60608152614ce660808261118b565b519020916001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b548015614dad57808203614d7f5750506020820151808203614d51575050516111bc906001600160801b0316614e77565b7f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b908151815103612b34575f5b8251811015614ba857614df4818461206c565b5151614e00828461206c565b515103614e3857614e11818461206c565b5160208151910120614e23828461206c565b516020815191012003614e3857600101614de1565b5050505f90565b60405190606090614e50828461118b565b6022835261409c91600a906020850190601f19013682375360206021840153600283614b35565b6001600160801b03633b9aca00911604614e95814242821115612bca565b4203428111612c23576001600160801b031663ffffffff5f516020614ef85f395f51905f5254169081811015614ec9575050565b7fdb020436000000000000000000000000000000000000000000000000000000005f5260045260245260445ffdfe5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1703a164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17084c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e05e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17036e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1707ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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
