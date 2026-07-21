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
	Bin: "0x610100604052346101a557616cd4803803809161001b826101ee565b61010039806101000160e082126101a5576100346102a2565b906100406101206102b9565b61004b6101406102b9565b610160516001600160401b0381116101a5578361006b91610100016102e8565b9061018051926101a0519660018060401b0388116101a557608090889003126101a5576040519461009b8661021a565b6101008801516001600160401b0381116101a55788018161011f820112156101a5576101008101516100cc8161032c565b916100da6040519384610250565b81835260206101008185019360051b83010101918483116101a5576101208201905b8382106101a95750505050865261011661012089016103c2565b6020870152610140880151906001600160401b0382116101a5578861014a610160926101006101559561016b9d0101610357565b604089015201610343565b60608601526101656101c06102b9565b95610934565b604051614ed69081611c7e823960805181610447015260a051816137a0015260c05181610ef9015260e051818181610ddc0152610e390152f35b5f80fd5b81516001600160401b0381116101a5576020916101cf8884610100819589010101610357565b8152019101906100fc565b634e487b7160e01b5f52604160045260245ffd5b610100601f91909101601f19168101906001600160401b0382119082101761021557604052565b6101da565b608081019081106001600160401b0382111761021557604052565b604081019081106001600160401b0382111761021557604052565b601f909101601f19168101906001600160401b0382119082101761021557604052565b6040519061028260e083610250565b565b60405190610282604083610250565b60405190610282608083610250565b61010051906001600160a01b03821682036101a557565b51906001600160a01b03821682036101a557565b6001600160401b03811161021557601f01601f191660200190565b81601f820112156101a557602081519101610302826102cd565b926103106040519485610250565b828452828201116101a557815f926020928386015e8301015290565b6001600160401b0381116102155760051b60200190565b51906001600160401b03821682036101a557565b91906080838203126101a557604051906103708261021a565b8351919384926001600160401b0381116101a5576060926103929183016102e8565b8352602081015160208401526103aa60408201610343565b60408401520151908160070b82036101a55760600152565b519081151582036101a557565b519060ff821682036101a557565b91908260409103126101a5576040516103f581610235565b602061040e818395610406816103cf565b8552016103cf565b910152565b91908260409103126101a55760405161042b81610235565b602061040e81839561043c81610343565b855201610343565b519063ffffffff821682036101a557565b6020818303126101a5578051906001600160401b0382116101a55701610120818303126101a557610484610273565b8151909290916001600160401b0383116101a5576104cb826104ae610100946105099685016102e8565b86526104bd81602085016103dd565b602087015260608301610413565b60408501526104dc60a08201610444565b60608501526104ed60c08201610444565b60808501526104fe60e082016103c2565b60a085015201610444565b60c082015290565b90600182811c9216801561053f575b602083101461052b57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691610520565b601f821161055657505050565b5f5260205f20906020601f840160051c8301931061058e575b601f0160051c01905b818110610583575050565b5f8155600101610578565b909150819061056f565b80518051906001600160401b038211610215576105d9826105c65f516020616c745f395f51905f5254610511565b5f516020616c745f395f51905f52610549565b602090601f83116001146107da5792610612836107959460c094610282975f926107cf575b50508160011b915f199060031b1c19161790565b5f516020616c745f395f51905f52555b60208181015180517f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170180549284015161ff0060089190911b1660ff90921661ffff199093169290921717905560408083015180517f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1702805492909401516fffffffffffffffff0000000000000000931b929092166001600160401b039092166001600160801b03199091161717905560608101516107029063ffffffff165f516020616bb45f395f51905f529063ffffffff1663ffffffff19825416179055565b610746610716608083015163ffffffff1690565b5f516020616bb45f395f51905f529067ffffffff0000000082549160201b169067ffffffff000000001916179055565b61078a61075660a0830151151590565b5f516020616bb45f395f51905f529068ff0000000000000000825491151560401b169068ff00000000000000001916179055565b015163ffffffff1690565b5f516020616bb45f395f51905f52906cffffffff00000000000000000082549160481b16906cffffffff0000000000000000001916179055565b015190505f806105fe565b5f516020616c745f395f51905f525f52601f19831691907fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be925f5b81811061086657509360c093610282969360019383610795981061084e575b505050811b015f516020616c745f395f51905f5255610622565b01515f1960f88460031b161c191690555f8080610834565b92936020600181928786015181550195019301610815565b156108865750565b63ffffffff9063b0369c3160e01b5f5216600452600160245263ffffffff60445260645ffd5b156108b5575050565b6355bace6f60e01b5f9081526001600160401b039182166004529116602452604490fd5b634e487b7160e01b5f52601160045260245ffd5b9063ffffffff8091169116019063ffffffff821161090757565b6108d9565b15610915575050565b9063ffffffff80926333fae18560e21b5f52166004521660245260445ffd5b90929493947f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a6123360e05280518101602001906020019061097291610455565b9161097c83610598565b606083019384516109909063ffffffff1690565b6109a29063ffffffff8116151561087e565b60c084019283516109b69063ffffffff1690565b6109c89063ffffffff8116151561087e565b84516109d390610bc2565b60200151604086018051519099916109fa916001600160401b0390811691168181146108ac565b8851602001516001600160401b03166001600160401b03165f9081527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1704602052604090205560805260a05260c052825163ffffffff16905163ffffffff16610a61916108ed565b60809190910151915163ffffffff9283169183168210159216610a839261090c565b610a8c90610dab565b51602001516001600160401b0316610aa390611106565b6001600160a01b038116610ac95750610aba6113b0565b50610ac660e051611432565b50565b80610ad6610ac6926112b2565b50610ae081611328565b5060e05161148b565b60405190610af682610235565b5f602083606081520152565b8015610907575f190190565b5f1981019190821161090757565b9190820391821161090757565b634e487b7160e01b5f52603260045260245ffd5b908151811015610b4e570160200190565b610b29565b906001820180921161090757565b603001908160301161090757565b906004820180921161090757565b90600c820180921161090757565b90602c820180921161090757565b600101908160011161090757565b602401908160241161090757565b9190820180921161090757565b610bca610ae9565b5080518015908115610d9f575b50610d90575f19908051805b610d40575b505f198214610d2257600360fc1b6001600160f81b0319610c22610c14610c0e86610b53565b85610b3d565b516001600160f81b03191690565b161480610d2c575b610d22575f90610c3983610b53565b915b8151831015610cd057610c5a610c54610c148585610b3d565b60f81c90565b60ff811660308110908115610cc5575b50610cb857600a82026001600160401b03908116602f1990920160ff1691909101811691168110610ca057600190920191610c3b565b50915050610cac610284565b9081525f602082015290565b5050915050610cac610284565b60399150115f610c6a565b9150916001811190811591610d16575b50610d0757610d0490610cf1610284565b9283526001600160401b03166020830152565b90565b63384c687760e21b5f5260045ffd5b602b915010155f610ce0565b9050610cac610284565b506002610d3a838351610b1c565b11610c2a565b602d60f81b610d6a610d5d610c14610d5785610b0e565b86610b3d565b6001600160f81b03191690565b14610d7e57610d7890610b02565b80610be3565b610d89919250610b0e565b905f610be8565b6329120bff60e21b5f5260045ffd5b6040915010155f610bd7565b610db481611518565b90518051908115610fee5760b48211610fd5575f93610df5610de5610de0610ddb866115c2565b610b61565b6115fa565b93610def85611a4a565b84611a92565b610e10610e09610e0486611c45565b61195b565b84602e0152565b610e1983611aa2565b60305f955b8351871015610ec957610ebe849392610eb9600193610ea2610e418c809a611504565b5191610e8163ffffffff610e786040860193610e72610e66865160018060401b031690565b6001600160401b031690565b90610bb5565b9b16868d611a6e565b610e9b610e8d86610b6f565b91516001600160401b031690565b908b611af6565b6020610ead84610b7d565b91015190890160200152565b610b8b565b960195909192610e1e565b91955061028294610fa8946020945092610f599250610f0490610ef56001600160401b0382111561162c565b6001600160401b031684611ab0565b610f48610f1a610f148584611657565b946117c5565b5f516020616bf45f395f51905f5280546001600160a01b0319166001600160a01b0392909216919091179055565b5f516020616c545f395f51905f5255565b8051610f9f906001600160401b03165f516020616bf45f395f51905f528054600160a01b600160e01b03191660a09290921b600160a01b600160e01b0316919091179055565b015161ffff1690565b5f516020616bf45f395f51905f52805461ffff60e01b191660e09290921b61ffff60e01b16919091179055565b63156f758160e31b5f52600482905260b460245260445ffd5b6305f8ded760e21b5f5260045ffd5b905f516020616b745f395f51905f5254821015610b4e575f516020616b745f395f51905f525f52600282901c5f516020616c145f395f51905f52019160031b60181690565b1561104b575050565b630fb2737b60e31b5f9081526001600160401b039182166004529116602452604490fd5b5f516020616b745f395f51905f5254906801000000000000000082101561021557600182015f516020616b745f395f51905f52555f516020616b745f395f51905f5254821015610b4e575f516020616b745f395f51905f525f52600282901c5f516020616c145f395f51905f520180546001600160401b0360069490941b60c01684811b199091169390921690911b919091179055565b5f516020616b745f395f51905f5254818161126a575b50506001600160401b0381165f9081525f516020616c945f395f51905f5260205260409020600101545f516020616c545f395f51905f52545f516020616bf45f395f51905f52546001600160a01b03928316159261125a92916111d6916111cb908216916111bb61119f60a083901c6001600160401b03169260e01c61ffff1690565b936111a8610293565b9687526001600160a01b03166020870152565b6001600160401b03166040850152565b61ffff166060830152565b6001600160401b0384165f9081525f516020616c945f395f51905f526020526040902081518155602082015160019190910180546040840151606094909401516001600160f01b03199091166001600160a01b03939093169290921760a09390931b600160a01b600160e01b03169290921760e09190911b61ffff60e01b16179055565b6112615750565b6102829061106f565b61129661128161127c6112ab94610b0e565b610ffd565b905460039190911b1c6001600160401b031690565b6001600160401b038181169083161015611042565b5f8161111c565b6001600160a01b0381165f9081525f516020616cb45f395f51905f52602052604090205460ff16611323576001600160a01b03165f8181525f516020616cb45f395f51905f5260205260408120805460ff191660011790553391905f516020616b545f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f516020616b945f395f51905f52602052604090205460ff16611323576001600160a01b0381165f9081525f516020616b945f395f51905f5260205260409020805460ff1916600117905533906001600160a01b03165f516020616c345f395f51905f525f516020616b545f395f51905f525f80a4600190565b5f80525f516020616b945f395f51905f526020525f516020616bd45f395f51905f525460ff1661142e575f8080525f516020616b945f395f51905f526020525f516020616bd45f395f51905f52805460ff1916600117905533905f516020616c345f395f51905f525f516020616b545f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f14611323575f818152602081815260408083208380529091528120805460ff1916600117905533915f516020616b545f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff166114fe575f818152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905533916001600160a01b0316905f516020616b545f395f51905f525f80a4600190565b50505f90565b8051821015610b4e5760209160051b010190565b8051519081156114fe5761152b8261032c565b916115396040519384610250565b808352601f196115488261032c565b013660208501375f5b825180518210156115b457906115a3610e04602061157184600196611504565b51015161159e6115966040611587878b51611504565b5101516001600160401b031690565b610cf1610284565b611876565b6115ad8287611504565b5201611551565b505090505f610d04926119b0565b90602c820291808304602c149015171561090757565b604051608091906115e98382610250565b6041815291601f1901366020840137565b90611604826102cd565b6116116040519182610250565b8281528092611622601f19916102cd565b0190602036910137565b15610fee57565b604051906116408261021a565b5f6060838281528260208201528260408201520152565b919091611662611633565b506030835110611721576356414c34602084015160e01c0361172157602483015160c01c602c84015160f01c93602e810151906116ef604e82015160f01c936116bb6116ac610293565b6001600160401b039092168252565b6116cd6020820198899061ffff169052565b604081019384526116e66060820195869061ffff169052565b965161ffff1690565b9061ffff8216938415948515611770575b508415611761575b50831561174a575b50508115611733575b506117215750565b634724a0fd60e01b5f5260045260245ffd5b905051611742610e0483611c45565b14155f611719565b5191925061175790611b3e565b1415905f80611710565b5161ffff16151593505f611708565b60b41094505f611700565b805191908290602001825e015f815290565b606160f81b81526001600160f01b031990911660018201526680600a3d393df360c81b6003820152610d0491600a919091019061177b565b604051905f60208301526117ee826117e0602182018461177b565b03601f198101845283610250565b616000825111611863575061183061183e61181e61180e845161ffff1690565b60f01b6001600160f01b03191690565b9260405192839160208301958661178d565b03601f198101835282610250565b51905ff0906001600160a01b0382161561185457565b63fbad885d60e01b5f5260045ffd5b516331734c2160e21b5f5260045260245ffd5b6020810180516024906001600160401b031680611915575b5061189b6118c9916115fa565b926118c06118ba6118b46118ae87611b63565b87611b6f565b86611b86565b85611b9d565b90519084611bca565b81519091906118e0906001600160401b0316610e66565b6118e957505090565b61190a610e666118fc6119119486611bb3565b92516001600160401b031690565b9083611be2565b5090565b600191505b6080811015611942575061189b61193b6119366118c993610b99565b610ba7565b915061188e565b60019060071c91019061191a565b6040513d5f823e3d90fd5b80516001810180911161090757611971906115fa565b805115610b4e5761199b8161198e6020945f868196015382611c1b565b506040519182809261177b565b039060025afa156119ab575f5190565b611950565b90929192808403938085116109075760018514611a395760015b8060011b90868210156119dd57506119ca565b9192939495505082019182811161090757826119f991856119b0565b91611a0492936119b0565b611a0c6115d8565b91825115610b4e57825f9261199b926001602080970153602183015260418201526040519182809261177b565b5090611a46929350611504565b5190565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602d908260081c602c8201530153565b604f5f9182604e8201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b61ffff16602c810290808204602c149015171561090757603001806030116109075790565b6020600a910153600190565b602082602292010153600181018091116109075790565b602082600a92010153600181018091116109075790565b6020828192010153600181018091116109075790565b602082601092010153600181018091116109075790565b81602091939293010152602081018091116109075790565b9092919083016020015b6080821015611c0057906001929391530190565b600180916080607f85161781530193019060071c9092611bec565b908051918215611c3d576021602084930191015e600101806001116109075790565b505050600190565b60405190606090611c568284610250565b6022835261191191600a906020850190601f19013682375360206021840153600283611bca56fe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a714610134578063248a9ca31461012f5780632f2ff15d1461012a57806336568abe146101255780634b1872be146101205780636a28f0001461011b5780636edfe3af146101165780638a8e4c5d146101115780638c80fbda1461010c57806391d148541461010757806392c19bdc14610102578063974a74c4146100fd578063a217fddf146100f8578063a6f031bb146100f3578063d547741f146100ee578063db3e1fa4146100e9578063ddba6537146100e45763ef913a4b146100df575f80fd5b610fc5565b610dff565b610dc5565b610d96565b610c7e565b610c64565b610afa565b610a69565b610a28565b6109ec565b6109b4565b6108ac565b610702565b61033b565b610267565b610231565b6101d9565b346101d55760203660031901126101d5576004357fffffffff0000000000000000000000000000000000000000000000000000000081168091036101d557807f7965db0b00000000000000000000000000000000000000000000000000000000602092149081156101ab575b506040519015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150145f6101a0565b5f80fd5b346101d55760203660031901126101d55760206102036004355f525f602052600160405f20015490565b604051908152f35b60409060031901126101d557600435906024356001600160a01b03811681036101d55790565b346101d5576102656102423661020b565b9061026061025b825f525f602052600160405f20015490565b612891565b6128d8565b005b346101d5576102753661020b565b336001600160a01b0382160361028e5761026591612970565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b9060206003198301126101d5576004356001600160401b0381116101d557826023820112156101d5578060040135926001600160401b0384116101d557602484830101116101d5576024019190565b634e487b7160e01b5f52602160045260245ffd5b6003111561032357565b610305565b9190602083019261033882610319565b52565b346101d55761067261047a61046b6104376104456103d361035b366102b6565b61037960ff5f516020614eaa5f395f51905f525460401c16156111d1565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475416156106eb575b810190611a89565b6103ef6103ea60408301516001600160801b031690565b612acd565b6104026060820151602083015190612b5a565b6040519283917f7c046dec00000000000000000000000000000000000000000000000000000000602084015260248301611e3c565b03601f198101835282611285565b7f0000000000000000000000000000000000000000000000000000000000000000612bdd565b60208082518301019101611f24565b61048381612c2b565b9061048d82610319565b81610676576106626106486020838160607f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb319601916105216104d983855101516001600160401b031690565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025460401c6001600160401b03166001600160401b0381166001600160401b03831611611f82565b6105d683516001600160401b038151166fffffffffffffffffffffffffffffffff196fffffffffffffffff000000000000000060207f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025494846001600160401b03198716177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170255015160401b16921617177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170255565b01516040516105ec816104378682019485611fc6565b5190208151830151610638906001600160401b03165b6001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b555101516001600160401b031690565b6040516001600160401b0390911681529081906020820190565b0390a15b60405191829182610328565b0390f35b5061068081610319565b60018103610666576106c26801000000000000000068ff0000000000000000195f516020614eaa5f395f51905f525416175f516020614eaa5f395f51905f5255565b7f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1610666565b6106f3612822565b6103cb565b5f9103126101d557565b346101d5575f3660031901126101d557335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff16156107cc575f516020614eaa5f395f51905f525460ff8160401c16156107a45768ff000000000000000019165f516020614eaa5f395f51905f52557f8e0da210a2ffecc744373b4860cd9b8dc471810f5fb61b7d2c3234ccd4aae30d5f80a1005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b90602080835192838152019201905f5b8181106108005750505090565b82518452602093840193909201916001016107f3565b919060608301916060845281518093526020608085019201925f5b81811061089057505061084c925083820360208501526107e3565b906040818303910152602080835192838152019201905f5b8181106108715750505090565b82516001600160401b0316845260209384019390920191600101610864565b845163ffffffff16845260209485019490930192600101610831565b346101d5575f3660031901126101d55760206108c6612d31565b6108cf81612dbd565b9201916108f06108eb6108e4855161ffff1690565b61ffff1690565b611ff4565b906109036108eb6108e4865161ffff1690565b936109166108eb6108e4835161ffff1690565b915f5b866109296108e4855161ffff1690565b61ffff8316908110156109a2579161099b60019261098d858061098661097e828f8f6109788f928f9261ffff9f6109638161096e9361203a565b9063ffffffff169052565b51925161ffff1690565b91612eb0565b92909561203a565b528961203a565b906001600160401b03169052565b0116610919565b50856106728660405193849384610816565b346101d5576109c2366102b6565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b346101d5575f3660031901126101d55760207f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170554604051908152f35b346101d557602060ff610a5d610a3d3661020b565b905f525f845260405f20906001600160a01b03165f5260205260405f2090565b54166040519015158152f35b346101d5576020610adc610a7c366102b6565b90610a9b60ff5f516020614eaa5f395f51905f525460401c16156111d1565b7fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5f525f845260405f205f8052845260ff60405f20541615610aed5761224e565b60405190610ae981610319565b8152f35b610af5612822565b61224e565b346101d55760203660031901126101d5576004356001600160401b0381116101d557806004019061016060031982360301126101d557610b4e60ff5f516020614eaa5f395f51905f525460401c16156111d1565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610c57575b610144810190610bb08284612593565b905015610c2f57610c10610c1f9261067294610bcf60448501826125c5565b610bdf60648794939401836125c5565b90608488013592610bf36101048a01612604565b9560a4610c17610c076101248d01896125c5565b9b909a89612593565b36916113ae565b9a01956136e7565b6040519081529081906020820190565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b610c5f612822565b610ba0565b346101d5575f3660031901126101d55760206040515f8152f35b346101d55760203660031901126101d5576004356001600160401b0381116101d5578060040161014060031983360301126101d55761067291610c1f91610cd960ff5f516020614eaa5f395f51905f525460401c16156111d1565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610d89575b610d3860448301826125c5565b91610d4660648501826125c5565b909190608486013590610d5c6101048801612604565b93610d6b6101248901856125c5565b97909660a46040519a610d7f60208d611285565b5f8c5201956136e7565b610d91612822565b610d2b565b346101d557610265610da73661020b565b90610dc061025b825f525f602052600160405f20015490565b612970565b346101d5575f3660031901126101d55760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346101d557610f1d610437610ef7610e86610e19366102b6565b610e3760ff5f516020614eaa5f395f51905f525460401c16156111d1565b7f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260ff610e7560405f205f805260205260405f2090565b541615610f7e575b5081019061266b565b610e9d6103ea60608301516001600160801b031690565b610eae6080820151825151906138fb565b610ec260a0820151602083510151906138fb565b6040519283917f60312a7c00000000000000000000000000000000000000000000000000000000602084015260248301612739565b7f0000000000000000000000000000000000000000000000000000000000000000612bdd565b50610f586801000000000000000068ff0000000000000000195f516020614eaa5f395f51905f525416175f516020614eaa5f395f51905f5255565b7f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1005b610f8790612891565b5f610e7d565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b906020610fc2928181520190610f8d565b90565b346101d5575f3660031901126101d55760405160208082015261012060408201525f7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1700548060011c90600181169081156111c7575b6020831082146111b357610160850183905261018085019190811561119a5750600114611129575b6106728461111d818661108660608301602060ff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170154818116845260081c16910152565b6110c760a0830160206001600160401b037f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170254818116845260401c16910152565b5f516020614eaa5f395f51905f525463ffffffff811660e084015261043790602081901c63ffffffff1661010085015261110c610120850160ff8360401c1615159052565b60481c63ffffffff16610140840152565b60405191829182610fb1565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005f9081529250907fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be5b8184106111865750500161111d82611042565b805484840152602090930192600101611173565b60ff191682525090151560051b01905061111d82611042565b634e487b7160e01b5f52602260045260245ffd5b91607f169161101a565b156111d857565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52604160045260245ffd5b606081019081106001600160401b0382111761122f57604052565b611200565b604081019081106001600160401b0382111761122f57604052565b608081019081106001600160401b0382111761122f57604052565b602081019081106001600160401b0382111761122f57604052565b90601f801991011681019081106001600160401b0382111761122f57604052565b604051906112b661026083611285565b565b604051906112b661010083611285565b604051906112b660c083611285565b604051906112b6608083611285565b6001600160801b038116036101d557565b35906112b6826112e6565b91908260609103126101d55760405161131a81611214565b6040808294803561132a816112e6565b8452602081013560208501520135910152565b6001600160401b038116036101d557565b35906112b68261133d565b91908260409103126101d55760405161137181611234565b602080829480356113818161133d565b845201359161138f8361133d565b0152565b6001600160401b03811161122f57601f01601f191660200190565b9291926113ba82611393565b916113c86040519384611285565b8294818452818301116101d5578281602093845f960137010152565b9080601f830112156101d557816020610fc2933591016113ae565b359081151582036101d557565b359063ffffffff821682036101d557565b8092910391606083126101d55760405161143681611234565b6040819483358352601f1901126101d557602090604080519361145885611234565b61146384820161140c565b85520135828401520152565b6001600160401b03811161122f5760051b60200190565b919060c0838203126101d55760405161149e8161124f565b809380356114ab8161133d565b82526114b96020820161140c565b60208301526114cb836040830161141d565b604083015260a0810135906001600160401b0382116101d5570182601f820112156101d5578035906114fc8261146f565b9361150a6040519586611285565b82855260208086019360051b830101918183116101d557602001925b828410611537575050505060600152565b6020848303126101d5576040519061154e8261126a565b84359060048210156101d55790825290815260209384019301611526565b91906060838203126101d5576040519061158582611234565b819380356001600160401b0381116101d5578101916040838203126101d557604051926115b184611234565b80356001600160401b0381116101d55781016102c0818403126101d5576115d66112a6565b906115e18482611359565b825260408101356001600160401b0381116101d557846116029183016113e4565b60208301526116136060820161134e565b6040830152611624608082016112f7565b606083015261163560a082016113ff565b60808301526116478460c0830161141d565b60a083015261165961012082016113ff565b60c083015261014081013560e083015261167661016082016113ff565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526116c561022082016113ff565b6101c08301526102408101356101e08301526116e461026082016113ff565b6102008301526102808101356102208301526102a0810135906001600160401b0382116101d557611717918591016113e4565b61024082015284526020810135926001600160401b0384116101d5576020946117468461175296889501611486565b83820152865201611359565b910152565b9080601f830112156101d557610100604051926117748285611285565b839181019283116101d557905b82821061178e5750505090565b8135815260209182019101611781565b9080601f830112156101d557604051916117b9604084611285565b8290604081019283116101d557905b8282106117d55750505090565b81358152602091820191016117c8565b359061ffff821682036101d557565b9080601f830112156101d557813561180b8161146f565b926118196040519485611285565b81845260208085019260051b8201019283116101d557602001905b8282106118415750505090565b6020809161184e8461140c565b815201910190611834565b9080601f830112156101d55781356118708161146f565b9261187e6040519485611285565b81845260208085019260051b8201019283116101d557602001905b8282106118a65750505090565b8135815260209182019101611899565b9080601f830112156101d55781356118cd8161146f565b926118db6040519485611285565b81845260208085019260051b8201019283116101d557602001905b8282106119035750505090565b60208091611910846113ff565b8152019101906118f6565b919091610220818403126101d5576119316112b8565b9261193c8183611757565b845261194c81610100840161179e565b602085015261195f81610140840161179e565b604085015261197161018083016117e5565b60608501526101a08201356001600160401b0381116101d557816119969184016117f4565b60808501526101c08201356001600160401b0381116101d557816119bb9184016117f4565b60a08501526101e08201356001600160401b0381116101d557816119e0918401611859565b60c08501526102008201356001600160401b0381116101d557611a0392016118b6565b60e0830152565b919060c0838203126101d55760405190611a238261124f565b8193611a2f8282611302565b835260608101356001600160401b0381116101d55782611a5091830161156c565b60208401526080810135611a63816112e6565b604084015260a0810135916001600160401b0383116101d557606092611752920161191b565b906020828203126101d55781356001600160401b0381116101d557610fc29201611a0a565b602060e0606060c08501936001600160401b03815116865263ffffffff848201511684870152611b006040820151604088019060208060409280518552015163ffffffff815116828501520151910152565b01519360c060a08201528451809452019201905f905b808210611b235750505090565b90919283515190600482101561032357602081600193829352019401920190611b16565b90610fc290602080611cc085516060855282611cac825160406060890152611b8860a0890182516001600160401b0360208092828151168552015116910152565b610240611ba5848301516102c060e08c01526103608b0190610f8d565b60408301516001600160401b03166101008b01529160608101516001600160801b03166101208b0152608081015115156101408b015260a081015180516101608c0152602090810151805163ffffffff166101808d015201516101a08b015260c081015115156101c08b015260e08101516101e08b015261010081015115156102008b01526101208101516102208b0152610140810151828b01526101608101516102608b01526101808101516102808b01526101a08101516102a08b01526101c081015115156102c08b01526101e08101516102e08b015261020081015115156103008b01526102208101516103208b01520151888203609f19016103408a0152610f8d565b910151858203605f19016080870152611aae565b9401519101906001600160401b0360208092828151168552015116910152565b905f905b60088210611cf157505050565b6020806001928551815201930191019091611ce4565b905f905b60028210611d1857505050565b6020806001928551815201930191019091611d0b565b90602080835192838152019201905f5b818110611d4b5750505090565b825163ffffffff16845260209384019390920191600101611d3e565b90602080835192838152019201905f5b818110611d845750505090565b82511515845260209384019390920191600101611d77565b610fc291611dab818351611ce0565b611dbe6020830151610100830190611d07565b611dd16040830151610140830190611d07565b606082015161ffff1661018082015260e0611e2a611e17611e0460808601516102206101a0870152610220860190611d2e565b60a08601518582036101c0870152611d2e565b60c08501518482036101e08601526107e3565b92015190610200818403910152611d67565b90610fc29160208152611e70602082018351604080916001600160801b038151168452602081015160208501520151910152565b6060611e8b602084015160c0608085015260e0840190611b47565b926001600160801b0360408201511660a084015201519060c0601f1982850301910152611d9c565b91908260609103126101d557604051611ecb81611214565b60408082948051611edb816112e6565b8452602081015160208501520151910152565b91908260409103126101d557604051611f0681611234565b60208082948051611f168161133d565b845201519161138f8361133d565b90610140828203126101d557611f7a9061010060405193611f448561124f565b611f4e8382611eb3565b8552611f5d8360608301611eb3565b6020860152611f6f8360c08301611eee565b604086015201611eee565b606082015290565b15611f8b575050565b906001600160401b0380927f7d939bd8000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6112b6909291926060810193604080916001600160801b038151168452602081015160208501520151910152565b90611ffe8261146f565b61200b6040519182611285565b828152809261201c601f199161146f565b0190602036910137565b634e487b7160e01b5f52603260045260245ffd5b805182101561204e5760209160051b010190565b612026565b91906080838203126101d5576040519061206c8261124f565b819380356001600160401b0381116101d55760609261208c9183016113e4565b83526020810135602084015260408101356120a68161133d565b60408401520135908160070b82036101d55760600152565b6020818303126101d5578035906001600160401b0382116101d55701906040828203126101d557604051916120f283611234565b80356001600160401b0381116101d5578261210e918301611a0a565b83526020810135906001600160401b0382116101d557016080818303126101d5576040519161213c8361124f565b81356001600160401b0381116101d557820181601f820112156101d55780356121648161146f565b916121726040519384611285565b81835260208084019260051b820101918483116101d55760208201905b8382106121ea575050505083526121a8602083016113ff565b60208401526040820135916001600160401b0383116101d5576121d26060926121dd948301612053565b60408501520161134e565b6060820152602082015290565b81356001600160401b0381116101d55760209161220c88848094880101612053565b81520191019061218f565b15612220575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b61225a918101906120be565b8051906122746103ea60408401516001600160801b031690565b61229661046b6104376104456060860151956104026020820197885190612b5a565b906122a082612c2b565b926122aa84610319565b6001841461252f577f1411889dcdfcf9af1fcf4ff4521cab4580a1fd842401ead543815cb1148d6438926122fb60206101609401926122e98451612f71565b90515151909401518490808214612217565b61230485610319565b8461249c5761247e916020826123366060839501936123306104d985875101516001600160401b031690565b5161300e565b6123eb83516001600160401b038151166fffffffffffffffffffffffffffffffff196fffffffffffffffff000000000000000060207f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025494846001600160401b03198716177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170255015160401b16921617177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170255565b0151604051612401816104378682019485611fc6565b519020815183015161241b906001600160401b0316610602565b558051820151612433906001600160401b0316613485565b7f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb3161246c610648848451016001600160401b0390511690565b0390a15101516001600160401b031690565b604080516001600160401b039290921682526020820192909252a190565b61247e91612509606060209301916123306124c185855101516001600160401b031690565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025460401c6001600160401b03166001600160401b0381166001600160401b03831614611f82565b8051820151612520906001600160401b0316613485565b5101516001600160401b031690565b50505061256c6801000000000000000068ff0000000000000000195f516020614eaa5f395f51905f525416175f516020614eaa5f395f51905f5255565b7f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a190565b903590601e19813603018212156101d557018035906001600160401b0382116101d5576020019181360383136101d557565b903590601e19813603018212156101d557018035906001600160401b0382116101d557602001918160051b360383136101d557565b6001111561032357565b3560018110156101d55790565b91906040838203126101d5576040519061262a82611234565b819380356001600160401b0381116101d5578261264891830161156c565b83526020810135916001600160401b0383116101d557602092611752920161156c565b6020818303126101d5578035906001600160401b0382116101d55701610140818303126101d55761269a6112c8565b9181356001600160401b0381116101d557816126b7918401612611565b83526126c68160208401611302565b60208401526126d88160808401611302565b60408401526126e960e083016112f7565b60608401526101008201356001600160401b0381116101d5578161270e91840161191b565b60808401526101208201356001600160401b0381116101d557612731920161191b565b60a082015290565b90610fc2916020815260a061280c61278484516101406020860152602061276e825160406101608901526101a0880190611b47565b91015185820361015f1901610180870152611b47565b6127b360208601516040860190604080916001600160801b038151168452602081015160208501520151910152565b6127e1604086015184860190604080916001600160801b038151168452602081015160208501520151910152565b60608501516001600160801b03166101008501526080850151848203601f1901610120860152611d9c565b92015190610140601f1982850301910152611d9c565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff161561285a57565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260ff6128b83360405f20906001600160a01b03165f5260205260405f2090565b5416156128c25750565b63e2517d3f60e01b5f523360045260245260445ffd5b805f525f60205260ff6128ff8360405f20906001600160a01b03165f5260205260405f2090565b541661296a57805f525f60205261292a8260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916600117905533916001600160a01b0316907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260ff6129978360405f20906001600160a01b03165f5260205260405f2090565b54161561296a57805f525f6020526129c38260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916905533916001600160a01b0316907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b634e487b7160e01b5f52601160045260245ffd5b6030019081603011612a2257565b612a00565b9060048201809211612a2257565b90600c8201809211612a2257565b90602c8201809211612a2257565b9060018201809211612a2257565b6001019081600111612a2257565b6024019081602411612a2257565b91908201809211612a2257565b15612a91575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b5f19810191908211612a2257565b6001600160801b03633b9aca0091160463ffffffff5f516020614eaa5f395f51905f525460481c16804201804211612a225782612b0e914290821115612a88565b428210612b19575050565b814203428111612a225711612b2b5750565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b606060206112b69351015101517f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170091612b91612d0d565b5061ffff6006600585015494015460405194612bac8661124f565b85526001600160a01b03811660208601526001600160401b038160a01c16604086015260e01c166060840152613bbf565b5f918291602082519201905af43d15612c23573d90612bfb82611393565b91612c096040519384611285565b82523d5f602084013e5b15612c1b5790565b602081519101fd5b606090612c13565b612c7d612c466020606084015101516001600160401b031690565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1704906001600160401b03165f5260205260405f2090565b5480612c895750505f90565b60208201908151604051612ca581610437602082019485611fc6565b5190201491821592612cc3575b505015612cbe57600190565b600290565b6001600160801b03919250612cf6612ce7612d0292516001600160801b0390511690565b9351516001600160801b031690565b6001600160801b031690565b911610155f80612cb2565b60405190612d1a8261124f565b5f6060838281528260208201528260408201520152565b612d39612d0d565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055461ffff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17065460405192612d8e8461124f565b83526001600160a01b03811660208401526001600160401b038160a01c16604084015260e01c16606082015290565b612dc5612d0d565b50602081016001600160a01b0381511615612e5a576001600160a01b03905116612df661ffff606084015116613def565b813b6001811115612e47575f198101908111612a22578111612e345790816001612e22612e3194613e36565b9260208401903c809251613e5e565b91565b50633610565160e21b5f5260045260245ffd5b82633610565160e21b5f5260045260245ffd5b5051633a517eed60e21b5f5260045260245ffd5b90602c820291808304602c1490151715612a2257565b90600382029180830460031490151715612a2257565b908160011b9180830460021490151715612a2257565b9091939263ffffffff61ffff91169416841015612f4157612ed8612ed385612e6e565b612a14565b9363ffffffff6020868501015160e01c1603612f165750610fc290612f0e6020612f0186612a27565b8301015160c01c94612a35565b016020015190565b7f4724a0fd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b83907ff81f5140000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b80515190811561296a57612f8482611ff4565b915f5b825180518210156130005790612fef612fea6020612fa78460019661203a565b510151612fe56001600160401b036040612fc2878b5161203a565b5101511660405192612fd384611234565b83526001600160401b03166020830152565b613f7b565b614070565b612ff9828761203a565b5201612f87565b505090505f610fc2926140c0565b61301781612f71565b9051805190811561331d5760b482116132ea575f9361305361304361303e612ed386612e6e565b613e36565b9361304d85614b1f565b84614b67565b613069613062612fea86614dae565b84602e0152565b61307283614b77565b60305f955b8351871015613123576131188493926131136001936130fc61309a8c809a61203a565b51916130db63ffffffff6130d260408601936130cc6130c086516001600160401b031690565b6001600160401b031690565b90612a7b565b9b16868d614b43565b6130f56130e786612a27565b91516001600160401b031690565b908b614bcb565b602061310784612a35565b91015190890160200152565b612a43565b960195909192613077565b61328f94929650602093506131f79150946131566001600160401b036112b69761314f82821115613aa6565b1684614b85565b6131d361316c6131668584613e5e565b946141d9565b6001600160a01b031673ffffffffffffffffffffffffffffffffffffffff197f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17065416177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170555565b61328661320b82516001600160401b031690565b7fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706549260a01b169116177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655565b015161ffff1690565b61ffff60e01b1961ffff60e01b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706549260e01b169116177f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170655565b7fab7bac08000000000000000000000000000000000000000000000000000000005f52600482905260b460245260445b5ffd5b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b907f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17085482101561204e577f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17085f52600282901c7f06ff75c3bdc3f03474569b3e952aab1bfb9d41b303b5638523d3e1200dcd1ae2019160031b60181690565b919091805483101561204e575f52601860205f208360021c019260031b1690565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708546801000000000000000081101561122f5780600161346792017f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708557f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17086133c3565b6001600160401b0380839493549260031b9316831b921b1916179055565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170854818161369b575b6001600160a01b03915060016134f7613505926001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f2090565b01546001600160a01b031690565b161561368b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1705546135bb7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1706546135b0613564826001600160a01b031690565b916135a061358460a083901c6001600160401b03169260e01c61ffff1690565b9361358d6112d7565b9687526001600160a01b03166020870152565b6001600160401b03166040850152565b61ffff166060830152565b6135f5846001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170760205260405f2090565b6060600161ffff928451815501926001600160a01b0360208201511673ffffffffffffffffffffffffffffffffffffffff1985541617845560408101517fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b8087549360a01b1616911617845501511661ffff60e01b1961ffff60e01b83549260e01b169116179055565b6136925750565b6112b6906133e4565b6136c56136b26136ad6136e094612abf565b613345565b90546001600160401b039160031b1c1690565b6001600160401b0381166001600160401b0383161015611f82565b5f816134ae565b999093969998929791959860018b1015610323578a1561373c5761331a8b61370e816125fa565b7f112d89cc000000000000000000000000000000000000000000000000000000005f5260ff16600452602490565b869798999a50602061376e918861377f949596979899151590816138ee575b61376891989796986142c6565b01614304565b613778368c611302565b9087614c13565b5f5f935b87851061384d575b5061379693506144af565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001691823b156101d5576138035f95604051978896879586957f56455e790000000000000000000000000000000000000000000000000000000087526004870161486f565b03915afa801561384857610fc29261382492612cf69261382e575b50614929565b633b9aca00900490565b8061383c5f61384293611285565b806106f8565b5f61381e565b614065565b9061388761388361387261386b613865898d8c61430e565b806125c5565b3691614330565b61387d368888614330565b90614d44565b1590565b6138e257506138c4906138ae610c106138a4613796978b8a61430e565b6020810190612593565b908151815180821491826138cc575b505061439b565b60015f61378b565b9091506020840120906020830120145f806138bd565b60019094019390613783565b61ffff811115915061375b565b906001600160401b0360208060608185510151015193015101511661391e612d0d565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1708548015613a90576139515f91612abf565b808210613a2057506139de916139a26139906136b26139d9947f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17086133c3565b916001600160401b0383161115614933565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1707906001600160401b03165f5260205260405f2090565b61494d565b916001600160a01b036139fb60208501516001600160a01b031690565b1615613a0c57906112b69291613bbf565b8251633a517eed60e21b5f5260045260245ffd5b90613a3c613a36613a318484612a7b565b612a51565b60011c90565b90836001600160401b03613a736136b2857f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17086133c3565b1611613a80575090613951565b9150613a8b90612abf565b613951565b633a517eed60e21b5f5261331a6024905f600452565b1561331d57565b15613ab55750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15613aef5750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15613b295750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b906001600160401b03809116911601906001600160401b038211612a2257565b15613b84575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b608081019384515193613c016060840195613bdf6108e4885161ffff1690565b8091149081613de0575b81613dd1575b81613dc2575b50908693959291613aa6565b613c0a81612dbd565b929095613c2160408401516001600160401b031690565b905f925f925f945f613c3a6108e45f9b5b5161ffff1690565b8a1015613d7657908b949392918e613c63613883613c5d8f8f9060e0015161203a565b51151590565b613d6357613c758c613c7f925161203a565b5163ffffffff1690565b8098613d45575b50508a600197968b8060a084015190613c9e9161203a565b5163ffffffff16918260208a0151613cb79061ffff1690565b61ffff1663ffffffff82161090613ccd91613ae7565b600163ffffffff84161b90613ce58482841615613aad565b1797828d8d519260200151613cfb9061ffff1690565b90613d0593612eb0565b91909360c0015190613d169161203a565b511490613d2291613b21565b613d2b91613b5b565b986108e46001613c3a925b019a929394959691908e613c32565b63ffffffff613d5c921663ffffffff821611613aad565b5f87613c86565b50959450986108e46001613c3a92613d36565b509350935097509893506112b6975060e09450613db89250613da06001600160401b038216612e84565b613db26001600160401b038416612e9a565b10613b7b565b5191015191614993565b905060e085015151145f613bf5565b60c08601515181149150613bef565b60a08601515181149150613be9565b61ffff16602c810290808204602c1490151715612a225760300180603011612a225790565b60405160809190613e258382611285565b6041815291601f1901366020840137565b90613e4082611393565b613e4d6040519182611285565b828152809261201c601f1991611393565b919091613e69612d0d565b506030835110612f16576356414c34602084015160e01c03612f1657602483015160c01c602c84015160f01c93602e81015190613ef6604e82015160f01c93613ec2613eb36112d7565b6001600160401b039092168252565b613ed46020820198899061ffff169052565b60408101938452613eed6060820195869061ffff169052565b965161ffff1690565b9061ffff8216938415948515613f70575b508415613f56575b508315613f3f575b50508115613f28575b50612f165750565b905051613f37612fea83614dae565b14155f613f20565b51919250613f4c90613def565b1415905f80613f17565b51909350613f679061ffff166108e4565b1515925f613f0f565b60b41094505f613f07565b602460208201906001600160401b0382511680614018575b50613fa0613fce91613e36565b92613fc5613fbf613fb9613fb387614a3d565b87614a49565b86614a60565b85614a77565b90519084614aa4565b90613fe36130c082516001600160401b031690565b613fec57505090565b61400d6130c0613fff6140149486614a8d565b92516001600160401b031690565b9083614abc565b5090565b600191505b60808110156140455750613fa061403e614039613fce93612a5f565b612a6d565b9150613f93565b60019060071c91019061401d565b805191908290602001825e015f815290565b6040513d5f823e3d90fd5b805160018101809111612a225761408690613e36565b80511561204e576140b0816140a36020945f868196015382614af5565b5060405191828092614053565b039060025afa15613848575f5190565b9092919280840393808511612a2257600185146141495760015b8060011b90868210156140ed57506140da565b91929394955050820191828111612a22578261410991856140c0565b9161411492936140c0565b61411c613e14565b9182511561204e57825f926140b09260016020809701536021830152604182015260405191828092614053565b509061415692935061203a565b5190565b600a907fffff000000000000000000000000000000000000000000000000000000000000610fc294937f610000000000000000000000000000000000000000000000000000000000000083521660018201527f80600a3d393df30000000000000000000000000000000000000000000000000060038201520190614053565b604051905f6020830152614202826141f46021820184614053565b03601f198101845283611285565b61600082511161429a575061043761425c61424a614222845161ffff1690565b60f01b7fffff0000000000000000000000000000000000000000000000000000000000001690565b9260405192839160208301958661415a565b51905ff0906001600160a01b0382161561427257565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b156142ce5750565b7fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b35610fc28161133d565b919081101561204e5760051b81013590603e19813603018212156101d5570190565b92919061433c8161146f565b9361434a6040519586611285565b602085838152019160051b8101918383116101d55781905b838210614370575050505050565b81356001600160401b0381116101d55760209161439087849387016113e4565b815201910190614362565b919091156143a7575050565b906143fb6143e9926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610f8d565b83810360031901602485015290610f8d565b0390fd5b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156101d55701602081359101916001600160401b0382116101d55781360383136101d557565b90602083828152019260208260051b82010193835f925b8484106144775750505050505090565b90919293949560208061449f600193601f198682030188526144998b8861441f565b906143ff565b9801940194019294939190614467565b919091156144bb575050565b6143fb6040519283927ffef760c7000000000000000000000000000000000000000000000000000000008452602060048501526024840191614450565b9035601e19823603018112156101d55701602081359101916001600160401b0382116101d5578160051b360383136101d557565b9035607e19823603018112156101d5570190565b359060038210156101d557565b9035605e19823603018112156101d5570190565b90602083828152019260208260051b82010193835f925b8484106145885750505050505090565b9091929394956020806145fa600193601f198682030188526145aa8b8861454d565b906145b482614540565b6145bd81610319565b81526145ec6145e16145d18685018561441f565b60608886015260608501916143ff565b92604081019061441f565b9160408185039101526143ff565b9801940194019294939190614578565b610fc2916146d46146c961464d614632614624868061441f565b6080875260808701916143ff565b61463f602087018761441f565b9086830360208801526143ff565b60806146b961465f604088018861452c565b868403604088015261467081614540565b61467981610319565b845261468760208201614540565b61469081610319565b60208501526146a160408201614540565b6146aa81610319565b6040850152606081019061441f565b91909281606082015201916143ff565b9260608101906144f8565b916060818503910152614561565b90602083828152019260208260051b82010193835f925b8484106147095750505050505090565b909192939495601f198282030184528635601e19843603018112156101d5578301906147396020820192806144f8565b8091936020845252604082019060408160051b8401019380935f915b8383106147785750505050505060208060019298019401940192949391906146f9565b909192939495603f19838203018652614791878361454d565b803560028110156101d55782526147bf6147ae602083018361452c565b60606020850152606084019061460a565b906040810135609e19823603018112156101d5576001936020938493614861930191604081830391015261485361483561480a6147fc858061441f565b60a0865260a08601916143ff565b6148158786016113ff565b151587850152614828604086018661452c565b848203604086015261460a565b92614842606082016113ff565b15156060840152608081019061452c565b90608081840391015261460a565b980196019493019190614755565b9092809492969593606083019083526060602084015252608081019560808560051b83010194815f90603e1981360301995b8383106148c0575050505050610fc294955060408185039101526146e2565b9091929397607f1986820301825288358b8112156101d557602061491a600193868394019061490d6149036148f584806144f8565b604085526040850191614450565b928581019061441f565b91858185039101526143ff565b9a0192019301919093926148a1565b35610fc2816112e6565b1561493a57565b633a517eed60e21b5f525f60045260245ffd5b9060405161495a8161124f565b606061ffff600183958054855201546001600160a01b03811660208501526001600160401b038160a01c16604085015260e01c16910152565b92916149a28251825114613aa6565b5f5b8251811015614a36576149b7818361203a565b5115614a2e5763ffffffff6149cc828561203a565b51166149db8187518110613ae7565b6149e5818761203a565b515160048110156103235760011901614a0357506001905b016149a4565b7f5b5c80eb000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b6001906149fd565b5050509050565b6020600a910153600190565b60208260229201015360018101809111612a225790565b602082600a9201015360018101809111612a225790565b602082819201015360018101809111612a225790565b60208260109201015360018101809111612a225790565b8160209193929301015260208101809111612a225790565b9092919083016020015b6080821015614ada57906001929391530190565b600180916080607f85161781530193019060071c9092614ac6565b908051918215614b17576021602084930191015e60010180600111612a225790565b505050600190565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602d908260081c602c8201530153565b604f5f9182604e8201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b90614c8f6040516020810190614c468287604080916001600160801b038151168452602081015160208501520151910152565b60608152614c55608082611285565b519020916001600160401b03165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2090565b548015614d1c57808203614cee5750506020820151808203614cc0575050516112b6906001600160801b0316614de6565b7f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b90815181510361296a575f5b8251811015614b1757614d63818461203a565b5151614d6f828461203a565b515103614da757614d80818461203a565b5160208151910120614d92828461203a565b516020815191012003614da757600101614d50565b5050505f90565b60405190606090614dbf8284611285565b6022835261401491600a906020850190601f19013682375360206021840153600283614aa4565b6001600160801b03633b9aca009116045f516020614eaa5f395f51905f525463ffffffff8160481c164201804211612a225782614e27914290821115612a88565b5f91428110614e89575b5063ffffffff16906001600160801b038116821115614e4e575050565b906001600160801b0380927fdb020436000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b9091504203428111612a22576001600160801b03169063ffffffff614e3156fe5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1703a164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17084c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e05e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17036e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170606ff75c3bdc3f03474569b3e952aab1bfb9d41b303b5638523d3e1200dcd1ae2bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17055e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1707ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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
