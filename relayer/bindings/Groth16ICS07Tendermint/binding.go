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

// ContractGroth16ICS07TendermintMetaData contains all meta data concerning the ContractGroth16ICS07Tendermint contract.
var ContractGroth16ICS07TendermintMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membership_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviour_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"updateClient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_clientState\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"initialPinnedValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPinnedValidatorSet\",\"inputs\":[],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reAnchorPinnedSet\",\"inputs\":[{\"name\":\"reAnchorMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"updateClientMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ClientFrozen\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClientUnfrozen\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClientUpdated\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PinnedSetReAnchored\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CannotHandleMisbehavior\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientNotFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ClientStateMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ClockDriftMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustedVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidMembershipProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NonMonotonicHeightUpdate\",\"inputs\":[{\"name\":\"latestHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"updateHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofHeightMismatch\",\"inputs\":[{\"name\":\"expectedRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"expectedRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofSignerCommitSigMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustThresholdMismatch\",\"inputs\":[{\"name\":\"expectedNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expectedDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnbondingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"UnknownZkAlgorithm\",\"inputs\":[{\"name\":\"algorithm\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorCountExceedsLimit\",\"inputs\":[{\"name\":\"validatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxValidatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"VerificationKeyMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x610160604052346101d35761758b803803809161001b8261021c565b61016039806101600161010082126101d3576100356102d1565b906100416101806102e8565b61004c6101a06102e8565b6100576101c06102e8565b6101e0516001600160401b0381116101d357846100779161016001610317565b916102005193610220519760018060401b0389116101d357608090899003126101d357604051956100a787610248565b6101608901516001600160401b0381116101d35789018161017f820112156101d3576101608101516100d88161035b565b916100e6604051938461027e565b81835260206101608185019360051b83010101918483116101d3576101808201905b8382106101d7575050505087526101226101808a016103f1565b60208801526101a0890151906001600160401b0382116101d357896101566101c092610160610161956101779e0101610386565b60408a015201610372565b60608701526101716102406102e8565b966109a0565b60405161588d9081611c3e823960805181615428015260a05181613b3e015260c05181610b15015260e05181818161207901526124a9015261010051818181610a0c0152610a54015261012051816142a2015261014051815050f35b5f80fd5b81516001600160401b0381116101d3576020916101fd8884610160819589010101610386565b815201910190610108565b634e487b7160e01b5f52604160045260245ffd5b610160601f91909101601f19168101906001600160401b0382119082101761024357604052565b610208565b608081019081106001600160401b0382111761024357604052565b604081019081106001600160401b0382111761024357604052565b601f909101601f19168101906001600160401b0382119082101761024357604052565b604051906102b16101008361027e565b565b604051906102b160408361027e565b604051906102b160808361027e565b61016051906001600160a01b03821682036101d357565b51906001600160a01b03821682036101d357565b6001600160401b03811161024357601f01601f191660200190565b81601f820112156101d357602081519101610331826102fc565b9261033f604051948561027e565b828452828201116101d357815f926020928386015e8301015290565b6001600160401b0381116102435760051b60200190565b51906001600160401b03821682036101d357565b91906080838203126101d3576040519061039f82610248565b8351919384926001600160401b0381116101d3576060926103c1918301610317565b8352602081015160208401526103d960408201610372565b60408401520151908160070b82036101d35760600152565b519081151582036101d357565b519060ff821682036101d357565b91908260409103126101d35760405161042481610263565b602061043d818395610435816103fe565b8552016103fe565b910152565b91908260409103126101d35760405161045a81610263565b602061043d81839561046b81610372565b855201610372565b519063ffffffff821682036101d357565b519060028210156101d357565b6020818303126101d3578051906001600160401b0382116101d35701610140818303126101d3576104c06102a1565b8151909290916001600160401b0383116101d357610507826104ea61012094610557968501610317565b86526104f9816020850161040c565b602087015260608301610442565b604085015261051860a08201610473565b606085015261052960c08201610473565b608085015261053a60e082016103f1565b60a085015261054c6101008201610484565b60c085015201610473565b60e082015290565b90600182811c9216801561058d575b602083101461057957565b634e487b7160e01b5f52602260045260245ffd5b91607f169161056e565b601f82116105a457505050565b5f5260205f20906020601f840160051c830193106105dc575b601f0160051c01905b8181106105d1575050565b5f81556001016105c6565b90915081906105bd565b600211156105f057565b634e487b7160e01b5f52602160045260245ffd5b9060028110156105f05769ff00000000000000000082549160481b169069ff0000000000000000001916179055565b80518051906001600160401b0382116102435761065c8261065560015461055f565b6001610597565b602090601f83116001146107fd5792610695836107c29460e0946102b1975f926107f2575b50508160011b915f199060031b1c19161790565b6001555b6020818101518051600280549284015161ff0060089190911b1660ff90921661ffff199093169290921717905560408083015180516003805492909401516fffffffffffffffff0000000000000000931b929092166001600160401b039092166001600160801b031990911617179055606081015161072f9063ffffffff1660049063ffffffff1663ffffffff19825416179055565b610767610743608083015163ffffffff1690565b60049067ffffffff0000000082549160201b169067ffffffff000000001916179055565b61079f61077760a0830151151590565b60049068ff0000000000000000825491151560401b169068ff00000000000000001916179055565b6107b760c08201516107b0816105e6565b6004610604565b015163ffffffff1690565b6004906dffffffff0000000000000000000082549160501b16906dffffffff000000000000000000001916179055565b015190505f80610681565b60015f52601f19831691905f51602061754b5f395f51905f52925f5b81811061085e57509360e0936102b19693600193836107c29810610846575b505050811b01600155610699565b01515f1960f88460031b161c191690555f8080610838565b92936020600181928786015181550195019301610819565b604051905f82600154916108898361055f565b80835292600181169081156108f957506001146108ad575b6102b19250038361027e565b5060015f90815290915f51602061754b5f395f51905f525b8183106108dd5750509060206102b1928201016108a1565b60209193508060019154838589010152019101909184926108c5565b602092506102b194915060ff191682840152151560051b8201016108a1565b15610921575050565b6355bace6f60e01b5f9081526001600160401b039182166004529116602452604490fd5b634e487b7160e01b5f52601160045260245ffd5b9063ffffffff8091169116019063ffffffff821161097357565b610945565b15610981575050565b9063ffffffff80926333fae18560e21b5f52166004521660245260445ffd5b90919293946109e96109e4610ae798977f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a612336101005260208082518301019101610491565b610633565b6109f1610876565b6020815191012061012052610a0c610a07610876565b610b47565b61014052610a7f610a66610a396020610a2b610a26610876565b610c18565b01516001600160401b031690565b60035490610a57906001600160401b03808416919081168214610918565b60401c6001600160401b031690565b6001600160401b03165f90815260056020526040902090565b556001600160a01b0390811660805290811660a05290811660c0521660e052600454610ae29063ffffffff8116610ad0610ac3605084901c63ffffffff1683610959565b9260201c63ffffffff1690565b9163ffffffff80841691161115610978565b610dfe565b600354610aff9060401c6001600160401b0316611093565b6001600160a01b038116610b265750610b166112ab565b50610b236101005161132d565b50565b80610b33610b23926111ad565b50610b3d81611223565b5061010051611386565b610b53610b5891611453565b6114ed565b90565b60405190610b6882610263565b5f602083606081520152565b8015610973575f190190565b5f1981019190821161097357565b9190820391821161097357565b634e487b7160e01b5f52603260045260245ffd5b908151811015610bc0570160200190565b610b9b565b906001820180921161097357565b603001908160301161097357565b906004820180921161097357565b90600c820180921161097357565b90602c820180921161097357565b9190820180921161097357565b610c20610b5b565b5080518015908115610df2575b50610de3575f19908051805b610d93575b505f198214610d7557600360fc1b6001600160f81b0319610c78610c6a610c6486610bc5565b85610baf565b516001600160f81b03191690565b161480610d7f575b610d75575f90610c8f83610bc5565b915b8151831015610d2657610cb0610caa610c6a8585610baf565b60f81c90565b60ff811660308110908115610d1b575b50610d0e57600a82026001600160401b03908116602f1990920160ff1691909101811691168110610cf657600190920191610c91565b50915050610d026102b3565b9081525f602082015290565b5050915050610d026102b3565b60399150115f610cc0565b9150916001811190811591610d69575b50610d5a57610b5890610d476102b3565b9283526001600160401b03166020830152565b63384c687760e21b5f5260045ffd5b602b915010155f610d36565b9050610d026102b3565b506002610d8d838351610b8e565b11610c80565b602d60f81b610dbd610db0610c6a610daa85610b80565b86610baf565b6001600160f81b03191690565b14610dd157610dcb90610b74565b80610c39565b610ddc919250610b80565b905f610c3e565b6329120bff60e21b5f5260045ffd5b6040915010155f610c2d565b610e0781611556565b9051805190811561100c5760b48211610ff3575f93610e48610e38610e33610e2e86611600565b610bd3565b611421565b93610e4285611ad4565b84611b1c565b610e5e610e57610b5386611c05565b84602e0152565b610e6783611b2c565b60305f955b8351871015610f1757610f0c849392610f07600193610ef0610e8f8c809a611542565b5191610ecf63ffffffff610ec66040860193610ec0610eb4865160018060401b031690565b6001600160401b031690565b90610c0b565b9b16868d611af8565b610ee9610edb86610be1565b91516001600160401b031690565b908b611b80565b6020610efb84610bef565b91015190890160200152565b610bfd565b960195909192610e6c565b9195506102b194610fd2946020945092610f8f9250610f5290610f436001600160401b03821115611616565b6001600160401b031684611b3a565b610f8a610f68610f628584611641565b9461177a565b600780546001600160a01b0319166001600160a01b0392909216919091179055565b600655565b8051610fc9906001600160401b031660078054600160a01b600160e01b03191660a09290921b600160a01b600160e01b0316919091179055565b015161ffff1690565b6007805461ffff60e01b191660e09290921b61ffff60e01b16919091179055565b63156f758160e31b5f52600482905260b460245260445ffd5b6305f8ded760e21b5f5260045ffd5b6009549068010000000000000000821015610243576001820180600955821015610bc05760095f52600282901c7f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af0180546001600160401b0360069490941b60c01684811b199091169390921690911b919091179055565b6001600160401b038181165f908152600860205260409020600101546006546007546001600160a01b03928316936111939361111b9261ffff60e082901c16926111109260a083901c9091169161110091166110ed6102c2565b9687526001600160a01b03166020870152565b6001600160401b03166040850152565b61ffff166060830152565b6001600160401b0384165f90815260086020526040902081518155602082015160019190910180546040840151606094909401516001600160f01b03199091166001600160a01b03939093169290921760a09390931b600160a01b600160e01b03169290921760e09190911b61ffff60e01b16179055565b6001600160a01b0316156111a45750565b6102b19061101b565b6001600160a01b0381165f9081525f51602061756b5f395f51905f52602052604090205460ff1661121e576001600160a01b03165f8181525f51602061756b5f395f51905f5260205260408120805460ff191660011790553391905f5160206174cb5f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f5160206174eb5f395f51905f52602052604090205460ff1661121e576001600160a01b0381165f9081525f5160206174eb5f395f51905f5260205260409020805460ff1916600117905533906001600160a01b03165f51602061752b5f395f51905f525f5160206174cb5f395f51905f525f80a4600190565b5f80525f5160206174eb5f395f51905f526020525f51602061750b5f395f51905f525460ff16611329575f8080525f5160206174eb5f395f51905f526020525f51602061750b5f395f51905f52805460ff1916600117905533905f51602061752b5f395f51905f525f5160206174cb5f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f1461121e575f818152602081815260408083208380529091528120805460ff1916600117905533915f5160206174cb5f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff166113f9575f818152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905533916001600160a01b0316905f5160206174cb5f395f51905f525f80a4600190565b50505f90565b60405160809190611410838261027e565b6041815291601f1901366020840137565b9061142b826102fc565b611438604051918261027e565b8281528092611449601f19916102fc565b0190602036910137565b8051156114a857611464815161182b565b80600101908160011161097357600190835101018091116109735761148b6114a491611421565b91600a602084015361149e8151846118b1565b8361194c565b5090565b506040516114b760208261027e565b5f8152601f196114c65f6102fc565b0136602083013790565b805191908290602001825e015f815290565b6040513d5f823e3d90fd5b8051600181018091116109735761150390611421565b805115610bc05761152d816115206020945f868196015382611922565b50604051918280926114d0565b039060025afa1561153d575f5190565b6114e2565b8051821015610bc05760209160051b010190565b8051519081156113f9576115698261035b565b91611577604051938461027e565b808352601f196115868261035b565b013660208501375f5b825180518210156115f257906115e1610b5360206115af84600196611542565b5101516115dc6115d460406115c5878b51611542565b5101516001600160401b031690565b610d476102b3565b611977565b6115eb8287611542565b520161158f565b505090505f610b5892611a3a565b90602c820291808304602c149015171561097357565b1561100c57565b6040519061162a82610248565b5f6060838281528260208201528260408201520152565b91909161164c61161d565b506030835110611701576356414c34602084015160e01c0361170157602483015160c01c602c84015160f01c936116d7602e82015192604e83015160f01c936116a56116966102c2565b6001600160401b039093168352565b6116b76020830198899061ffff169052565b60408201526116ce6060820194859061ffff169052565b955161ffff1690565b9061ffff8216928315938415611737575b508315611728575b508215611713575b50506117015750565b634724a0fd60e01b5f5260045260245ffd5b51915061171f90611bc8565b14155f806116f8565b5161ffff16151592505f6116f0565b60b41093505f6116e8565b606160f81b81526001600160f01b031990911660018201526680600a3d393df360c81b6003820152610b5891600a91909101906114d0565b604051905f60208301526117a38261179560218201846114d0565b03601f19810184528361027e565b61600082511161181857506117e56117f36117d36117c3845161ffff1690565b60f01b6001600160f01b03191690565b92604051928391602083019586611742565b03601f19810183528261027e565b51905ff0906001600160a01b0382161561180957565b63fbad885d60e01b5f5260045ffd5b516331734c2160e21b5f5260045260245ffd5b906001915b608081101561183c5750565b60019060071c920191611830565b6020600a910153600190565b602082602292010153600181018091116109735790565b602082600a92010153600181018091116109735790565b6020828192010153600181018091116109735790565b602082601092010153600181018091116109735790565b91906021600193015b60808210156118ce57906001929391530190565b600180916080607f85161781530193019060071c90926118ba565b9092919083016020015b608082101561190757906001929391530190565b600180916080607f85161781530193019060071c90926118f3565b908051918215611944576021602084930191015e600101806001116109735790565b505050600190565b90809291825192831561197057839260208092019201015e81018091116109735790565b5050505090565b6020810180516024906001600160401b031680611a12575b5061199c6119ca91611421565b926119c16119bb6119b56119af8761184a565b87611856565b8661186d565b85611884565b90519084611bed565b81519091906119e1906001600160401b0316610eb4565b6119ea57505090565b611a0b610eb46119fd6114a4948661189a565b92516001600160401b031690565b90836118e9565b611a1c915061182b565b6001018060011161097357602401806024116109735761199c61198f565b90929192808403938085116109735760018514611ac35760015b8060011b9086821015611a675750611a54565b919293949550508201918281116109735782611a839185611a3a565b91611a8e9293611a3a565b611a966113ff565b91825115610bc057825f9261152d92600160208097015360218301526041820152604051918280926114d0565b5090611ad0929350611542565b5190565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602d908260081c602c8201530153565b604f5f9182604e8201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b61ffff16602c810290808204602c149015171561097357603001806030116109735790565b81602091939293010152602081018091116109735790565b60405190606090611c16828461027e565b602283526114a491600a906020850190601f19013682375360206021840153600283611bed56fe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a7146101245780630bece3561461011f578063248a9ca31461011a5780632f2ff15d146101155780632fe76c3c1461011057806336568abe1461010b5780636a28f000146101065780636edfe3af146101015780638a8e4c5d146100fc57806391d14854146100f7578063974a74c4146100f2578063a217fddf146100ed578063a6f031bb146100e8578063d547741f146100e3578063db3e1fa4146100de578063ddba6537146100d95763ef913a4b146100d4575f80fd5b610c53565b610a2f565b6109f5565b6109c6565b6108b9565b61089f565b61073f565b6106fe565b6106c6565b6105ba565b610428565b6103cf565b61034e565b610318565b6102c0565b61023b565b346101c55760203660031901126101c5576004357fffffffff0000000000000000000000000000000000000000000000000000000081168091036101c557807f7965db0b000000000000000000000000000000000000000000000000000000006020921490811561019b575b506040519015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150145f610190565b5f80fd5b9060206003198301126101c5576004356001600160401b0381116101c557826023820112156101c5578060040135926001600160401b0384116101c557602484830101116101c5576024019190565b634e487b7160e01b5f52602160045260245ffd5b6003111561023657565b610218565b346101c55760206102a261024e366101c9565b9061026160ff60045460401c1615610d3e565b7fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5f525f845260405f205f8052845260ff60405f205416156102b357612044565b604051906102af8161022c565b8152f35b6102bb612d19565b612044565b346101c55760203660031901126101c55760206102ea6004355f525f602052600160405f20015490565b604051908152f35b60409060031901126101c557600435906024356001600160a01b03811681036101c55790565b346101c55761034c610329366102f2565b90610347610342825f525f602052600160405f20015490565b612d88565b613320565b005b346101c55761034c61035f366101c9565b9061037260ff60045460401c1615610d3e565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff16612472576103ca612d19565b612472565b346101c5576103dd366102f2565b336001600160a01b038216036103f65761034c91613872565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b5f9103126101c557565b346101c5575f3660031901126101c557335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff16156104da5760045460ff8160401c16156104b25768ff000000000000000019166004557f8e0da210a2ffecc744373b4860cd9b8dc471810f5fb61b7d2c3234ccd4aae30d5f80a1005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b90602080835192838152019201905f5b81811061050e5750505090565b8251845260209384019390920191600101610501565b919060608301916060845281518093526020608085019201925f5b81811061059e57505061055a925083820360208501526104f1565b906040818303910152602080835192838152019201905f5b81811061057f5750505090565b82516001600160401b0316845260209384019390920191600101610572565b845163ffffffff1684526020948501949093019260010161053f565b346101c5575f3660031901126101c5576105d26127a6565b5060206105dd613902565b91016105fd6105f86105f1835161ffff1690565b61ffff1690565b6127ca565b61060f6105f86105f1845161ffff1690565b926106226105f86105f1855161ffff1690565b60065490915f5b866106396105f1885161ffff1690565b61ffff8316908110156106b057916106a960019261069b858061069461068c828f8f6106868f928f9261ffff9f6106738161067e93612810565b9063ffffffff169052565b5161ffff1690565b916139c8565b929095612810565b5289612810565b906001600160401b03169052565b0116610629565b50856106c28660405193849384610524565b0390f35b346101c5576106d4366101c9565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b346101c557602060ff610733610713366102f2565b905f525f845260405f20906001600160a01b03165f5260205260405f2090565b54166040519015158152f35b346101c55760203660031901126101c5576004356001600160401b0381116101c557806004019061016060031982360301126101c55761078760ff60045460401c1615610d3e565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610892575b6101448101906107e98284612829565b90501561086a5761084b61085a926106c294610808604485018261285b565b610818606487949394018361285b565b906084880135926101048901359561082f87610fa3565b60a46108526108426101248d018961285b565b9b909a89612829565b3691610e8d565b9a0195613a84565b6040519081529081906020820190565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b61089a612d19565b6107d9565b346101c5575f3660031901126101c55760206040515f8152f35b346101c55760203660031901126101c5576004356001600160401b0381116101c5578060040161014060031983360301126101c5576106c29161085a9161090860ff60045460401c1615610d3e565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475416156109b9575b610967604483018261285b565b91610975606485018261285b565b61010486013592916084870135919061098d85610fa3565b61099b61012489018561285b565b97909660a46040519a6109af60208d610df2565b5f8c520195613a84565b6109c1612d19565b61095a565b346101c55761034c6109d7366102f2565b906109f0610342825f525f602052600160405f20015490565b613872565b346101c5575f3660031901126101c55760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346101c557610aa1610a40366101c9565b610a5260ff60045460401c1615610d3e565b7f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260ff610a9060405f205f805260205260405f2090565b541615610c0c575b50810190612a34565b805190602081018051906040830190815160806060860194855197610b0983890199610ad48b516001600160801b031690565b9060405196879586957f85fbdd9c00000000000000000000000000000000000000000000000000000000875260048701612b5e565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa8015610c0757610b8c9660c095604095610b6a945f94610bd6575b5088519051915192516001600160801b031693613c99565b610b7e60208251015160a086015190613d44565b505051015191015190613d44565b5050610bb06801000000000000000068ff0000000000000000196004541617600455565b7f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1005b610bf991945060803d608011610c00575b610bf18183610df2565b810190612b27565b925f610b52565b503d610be7565b611fc7565b610c1590612d88565b5f610a98565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b906020610c50928181520190610c1b565b90565b346101c5575f3660031901126101c5576106c26040516020808201526101406040820152610d3281610c886101808201612c7a565b610ca460608301602060ff600254818116845260081c16910152565b610cc660a0830160206001600160401b03600354818116845260401c16910152565b60045463ffffffff811660e0840152610d2490602081901c63ffffffff16610100850152610cff610120850160ff8360401c1615159052565b610d13610140850160ff8360481c16611a7e565b60501c63ffffffff16610160840152565b03601f198101835282610df2565b60405191829182610c3f565b15610d4557565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b03821117610d9c57604052565b610d6d565b606081019081106001600160401b03821117610d9c57604052565b608081019081106001600160401b03821117610d9c57604052565b60c081019081106001600160401b03821117610d9c57604052565b90601f801991011681019081106001600160401b03821117610d9c57604052565b60405190610e2361010083610df2565b565b60405190610e2361026083610df2565b60405190610e2361018083610df2565b60405190610e2360e083610df2565b60405190610e23608083610df2565b60405190610e23606083610df2565b6001600160401b038111610d9c57601f01601f191660200190565b929192610e9982610e72565b91610ea76040519384610df2565b8294818452818301116101c5578281602093845f960137010152565b9080601f830112156101c557816020610c5093359101610e8d565b60ff8116036101c557565b91908260409103126101c557604051610f0181610d81565b60208082948035610f1181610ede565b8452013591610f1f83610ede565b0152565b6001600160401b038116036101c557565b3590610e2382610f23565b91908260409103126101c557604051610f5781610d81565b60208082948035610f6781610f23565b8452013591610f1f83610f23565b63ffffffff8116036101c557565b3590610e2382610f75565b801515036101c557565b3590610e2382610f8e565b600211156101c557565b3590610e2382610fa3565b919091610140818403126101c557610fce610e13565b928135916001600160401b0383116101c55761101382610ff661012094611063968501610ec3565b87526110058160208501610ee9565b602088015260608301610f3f565b604086015261102460a08201610f83565b606086015261103560c08201610f83565b608086015261104660e08201610f98565b60a08601526110586101008201610fad565b60c086015201610f83565b60e0830152565b6001600160801b038116036101c557565b3590610e238261106a565b91908260609103126101c55760405161109e81610da1565b604080829480356110ae8161106a565b8452602081013560208501520135910152565b8092910391606083126101c5576040516110da81610d81565b6040819483358352601f1901126101c55760209060408051936110fc85610d81565b8381013561110981610f75565b85520135828401520152565b6001600160401b038111610d9c5760051b60200190565b919060c0838203126101c55760405161114481610dbc565b8093803561115181610f23565b8252602081013561116181610f75565b602083015261117383604083016110c1565b604083015260a0810135906001600160401b0382116101c557019180601f840112156101c5578235926111a584611115565b936111b36040519586610df2565b80855260208086019160051b830101918383116101c55760208101915b8383106111e257505050505060600152565b82356001600160401b0381116101c5578201906040828703601f1901126101c5576040519161121083610d81565b602081013560048110156101c557835260408101356001600160401b0381116101c5576020910101906080828803126101c5576040519261125084610dbc565b82356001600160401b0381116101c5578861126c918501610ec3565b8452602083013561127c8161106a565b6020850152604083013561128f81610f8e565b60408501526060830135936001600160401b0385116101c5576112b789602096879601610ec3565b6060820152838201528152019201916111d0565b91906060838203126101c557604051906112e482610d81565b819380356001600160401b0381116101c5578101916040838203126101c5576040519261131084610d81565b80356001600160401b0381116101c55781016102c0818403126101c557611335610e25565b906113408482610f3f565b825260408101356001600160401b0381116101c55784611361918301610ec3565b602083015261137260608201610f34565b60408301526113836080820161107b565b606083015261139460a08201610f98565b60808301526113a68460c083016110c1565b60a08301526113b86101208201610f98565b60c083015261014081013560e08301526113d56101608201610f98565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526114246102208201610f98565b6101c08301526102408101356101e08301526114436102608201610f98565b6102008301526102808101356102208301526102a0810135906001600160401b0382116101c55761147691859101610ec3565b61024082015284526020810135926001600160401b0384116101c5576020946114a5846114b19688950161112c565b83820152865201610f3f565b910152565b9080601f830112156101c557610100604051926114d38285610df2565b839181019283116101c557905b8282106114ed5750505090565b81358152602091820191016114e0565b9080601f830112156101c55760405191611518604084610df2565b8290604081019283116101c557905b8282106115345750505090565b8135815260209182019101611527565b359061ffff821682036101c557565b9080601f830112156101c557813561156a81611115565b926115786040519485610df2565b81845260208085019260051b8201019283116101c557602001905b8282106115a05750505090565b6020809183356115af81610f75565b815201910190611593565b9080601f830112156101c55781356115d181611115565b926115df6040519485610df2565b81845260208085019260051b8201019283116101c557602001905b8282106116075750505090565b81358152602091820191016115fa565b9080601f830112156101c557813561162e81611115565b9261163c6040519485610df2565b81845260208085019260051b8201019283116101c557602001905b8282106116645750505090565b60208091833561167381610f8e565b815201910190611657565b9190916102e0818403126101c557611694610e35565b9281356001600160401b0381116101c557816116b1918401610fb8565b84526116c08160208401611086565b602085015260808201356001600160401b0381116101c557816116e49184016112cb565b60408501526116f560a0830161107b565b60608501526117078160c084016114b6565b608085015261171a816101c084016114fd565b60a085015261172d8161020084016114fd565b60c085015261173f6102408301611544565b60e08501526102608201356001600160401b0381116101c55781611764918401611553565b6101008501526102808201356001600160401b0381116101c5578161178a918401611553565b6101208501526102a08201356001600160401b0381116101c557816117b09184016115ba565b6101408501526102c08201356001600160401b0381116101c5576117d49201611617565b610160830152565b906020828203126101c55781356001600160401b0381116101c557610c50920161167e565b81601f820112156101c55780519061181882610e72565b926118266040519485610df2565b828452602083830101116101c557815f9260208093018386015e8301015290565b91908260409103126101c55760405161185f81610d81565b6020808294805161186f81610ede565b8452015191610f1f83610ede565b91908260409103126101c55760405161189581610d81565b602080829480516118a581610f23565b8452015191610f1f83610f23565b5190610e2382610f75565b5190610e2382610f8e565b5190610e2382610fa3565b919091610140818403126101c5576118ea610e13565b928151916001600160401b0383116101c55761192f8261191261012094611063968501611801565b87526119218160208501611847565b60208801526060830161187d565b604086015261194060a082016118b3565b606086015261195160c082016118b3565b608086015261196260e082016118be565b60a086015261197461010082016118c9565b60c0860152016118b3565b5190610e238261106a565b91908260609103126101c5576040516119a281610da1565b604080829480516119b28161106a565b8452602081015160208501520151910152565b6020818303126101c5578051906001600160401b0382116101c55701610180818303126101c557604051916119f983610dd7565b81516001600160401b0381116101c55782611a1c8361014093611a6c96016118d4565b8552611a2b836020830161198a565b6020860152611a3d836080830161198a565b6040860152611a4e60e0820161197f565b6060860152611a6183610100830161187d565b60808601520161187d565b60a082015290565b6002111561023657565b90611a8882611a74565b52565b90610c509061012060e0611aaa85516101408552610140850190610c1b565b9460ff60208083015182815116828801520151166040850152611aea604082015160608601906001600160401b0360208092828151168552015116910152565b606081015163ffffffff1660a0850152608081015163ffffffff1660c085015260a0810151151584830152611b2860c0820151610100860190611a7e565b015163ffffffff16910152565b90606060c08201926001600160401b03815116835263ffffffff6020820151166020840152611b866040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519160c060a0830152825180915260e0820190602060e08260051b8501019401925f905b828210611bba57505050505090565b909192939460df198282030185528551908151600481101561023657611c388260206001958195948295520151906040848201526060611c0683516080604085015260c0840190610c1b565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610c1b565b9701950193920190611bab565b90610c5090602080611dbe85516060855282611daa825160406060890152611c8660a0890182516001600160401b0360208092828151168552015116910152565b610240611ca3848301516102c060e08c01526103608b0190610c1b565b60408301516001600160401b03166101008b01529160608101516001600160801b03166101208b0152608081015115156101408b015260a081015180516101608c0152602090810151805163ffffffff166101808d015201516101a08b015260c081015115156101c08b015260e08101516101e08b015261010081015115156102008b01526101208101516102208b0152610140810151828b01526101608101516102608b01526101808101516102808b01526101a08101516102a08b01526101c081015115156102c08b01526101e08101516102e08b015261020081015115156103008b01526102208101516103208b01520151888203609f19016103408a0152610c1b565b910151858203605f19016080870152611b35565b9401519101906001600160401b0360208092828151168552015116910152565b905f905b60088210611def57505050565b6020806001928551815201930191019091611de2565b905f905b60028210611e1657505050565b6020806001928551815201930191019091611e09565b90602080835192838152019201905f5b818110611e495750505090565b825163ffffffff16845260209384019390920191600101611e3c565b90602080835192838152019201905f5b818110611e825750505090565b82511515845260209384019390920191600101611e75565b90610c509160208152610160611fb1611f99611f81611f0f611eca87516102e06020890152610300880190611a8b565b611ef960208901516040890190604080916001600160801b038151168452602081015160208501520151910152565b6040880151878203601f190160a0890152611c45565b60608701516001600160801b031660c0870152611f34608088015160e0880190611dde565b611f4760a08801516101e0880190611e05565b611f5a60c0880151610220880190611e05565b60e087015161ffff16610260870152610100870151868203601f1901610280880152611e2c565b610120860151858203601f19016102a0870152611e2c565b610140850151848203601f19016102c08601526104f1565b920151906102e0601f1982850301910152611e65565b6040513d5f823e3d90fd5b15611fdb575050565b906001600160401b0380927f7d939bd8000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b610e23909291926060810193604080916001600160801b038151168452602081015160208501520151910152565b612050918101906117dc565b604051630913b29560e41b81525f818061206d8560048301611e9a565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115610c07575f91612283575b506120b381612e06565b6120c56120bf82612e87565b926130d1565b50506120d08261022c565b816122115761220b6121f1602083604060a07f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb3196019161214661211d85855101516001600160401b031690565b60035460401c6001600160401b03166001600160401b0381166001600160401b03831611611fd2565b61219e83516001600160401b038151166fffffffffffffffffffffffffffffffff196fffffffffffffffff0000000000000000602060035494846001600160401b0319871617600355015160401b1692161717600355565b01516040516121b481610d248682019485612016565b51902081518301516121e1906001600160401b03165b6001600160401b03165f52600560205260405f2090565b555101516001600160401b031690565b6040516001600160401b0390911681529081906020820190565b0390a190565b5061221b8161022c565b6001810361226c576122456801000000000000000068ff0000000000000000196004541617600455565b7f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a190565b6122758161022c565b60028103610c505750600290565b61229f91503d805f833e6122978183610df2565b8101906119c5565b5f6120a9565b91906080838203126101c557604051906122be82610dbc565b819380356001600160401b0381116101c5576060926122de918301610ec3565b83526020810135602084015260408101356122f881610f23565b60408401520135908160070b82036101c55760600152565b91906040838203126101c55782356001600160401b0381116101c5578161233891850161167e565b926020810135906001600160401b0382116101c557016080818303126101c5576040519161236583610dbc565b81356001600160401b0381116101c557820181601f820112156101c557803561238d81611115565b9161239b6040519384610df2565b81835260208084019260051b820101918483116101c55760208201905b83821061240e575050505083526123d160208301610f98565b60208401526040820135916001600160401b0383116101c5576123fb6060926124069483016122a5565b604085015201610f34565b606082015290565b81356001600160401b0381116101c557602091612430888480948801016122a5565b8152019101906123b8565b15612444575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b61247e91810190612310565b9060405191630913b29560e41b83525f838061249d8560048301611e9a565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa928315610c07575f9361278a575b506124e383612e06565b6124ec83612e87565b6124f5836130d1565b50506125008161022c565b6001811461273e5761252e60409361016061251a856133b8565b95869201515151018051821490519061243b565b6125378161022c565b80612693575060206126769160408561258b60a07f16f9bd20fd25a0828c55e0202f115427e910995861d180b9cf422f347cc0157598019361258661211d87875101516001600160401b031690565b613455565b6125e383516001600160401b038151166fffffffffffffffffffffffffffffffff196fffffffffffffffff0000000000000000602060035494846001600160401b0319871617600355015160401b1692161717600355565b01516040516125f981610d248682019485612016565b5190208151830151612613906001600160401b03166121ca565b55805182015161262b906001600160401b0316613723565b7f9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb316126646121f1848451016001600160401b0390511690565b0390a15101516001600160401b031690565b604080516001600160401b039290921682526020820192909252a1565b8061269f60029261022c565b146126a957505050565b60206126769161271860a07f16f9bd20fd25a0828c55e0202f115427e910995861d180b9cf422f347cc015759601916125866126ef85855101516001600160401b031690565b60035460401c6001600160401b03166001600160401b0381166001600160401b03831614611fd2565b805182015161272f906001600160401b0316613723565b5101516001600160401b031690565b505050506127646801000000000000000068ff0000000000000000196004541617600455565b7f59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c5f80a1565b61279f9193503d805f833e6122978183610df2565b915f6124d9565b604051906127b382610dbc565b5f6060838281528260208201528260408201520152565b906127d482611115565b6127e16040519182610df2565b82815280926127f2601f1991611115565b0190602036910137565b634e487b7160e01b5f52603260045260245ffd5b80518210156128245760209160051b010190565b6127fc565b903590601e19813603018212156101c557018035906001600160401b0382116101c5576020019181360383136101c557565b903590601e19813603018212156101c557018035906001600160401b0382116101c557602001918160051b360383136101c557565b91906060838203126101c557604051906128a982610da1565b819380356001600160401b0381116101c55781016040818403126101c557604051906128d482610d81565b8035906001600160401b0382116101c5576128f3856020938301610ec3565b8352013561290081610f23565b6020820152835260208101356001600160401b0381116101c557826129269183016112cb565b60208401526040810135916001600160401b0383116101c5576040926114b192016112cb565b919091610220818403126101c557612962610e13565b9261296d81836114b6565b845261297d8161010084016114fd565b60208501526129908161014084016114fd565b60408501526129a26101808301611544565b60608501526101a08201356001600160401b0381116101c557816129c7918401611553565b60808501526101c08201356001600160401b0381116101c557816129ec918401611553565b60a08501526101e08201356001600160401b0381116101c55781612a119184016115ba565b60c08501526102008201356001600160401b0381116101c5576110639201611617565b6020818303126101c5578035906001600160401b0382116101c55701610160818303126101c557612a63610e45565b9181356001600160401b0381116101c55781612a80918401610fb8565b835260208201356001600160401b0381116101c55781612aa1918401612890565b6020840152612ab38160408401611086565b6040840152612ac58160a08401611086565b6060840152612ad7610100830161107b565b60808401526101208201356001600160401b0381116101c55781612afc91840161294c565b60a08401526101408201356001600160401b0381116101c557612b1f920161294c565b60c082015290565b906080828203126101c557612b56906040805193612b4485610d81565b612b4e838261187d565b85520161187d565b602082015290565b90610e2394612c0e612be661010095612b88612c33959b9a989b6101208852610120880190611a8b565b86810360208801526040612bd58351606084526001600160401b036020612bba835186606089015260a0880190610c1b565b92015116608085015260208501518482036020860152611c45565b920151906040818403910152611c45565b986040850190604080916001600160801b038151168452602081015160208501520151910152565b80516001600160801b031660a0840152602081015160c08401526040015160e0830152565b01906001600160801b03169052565b90600182811c92168015612c70575b6020831014612c5c57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612c51565b6001545f9291612c8982612c42565b8082529160018116908115612cfd5750600114612ca4575050565b60015f9081529293509091907fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b838310612ce3575060209250010190565b600181602092949394548385870101520191019190612cd2565b9050602093945060ff929192191683830152151560051b010190565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff1615612d5157565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260ff612daf3360405f20906001600160a01b03165f5260205260405f2090565b541615612db95750565b63e2517d3f60e01b5f523360045260245260445ffd5b15612dd8575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b610e2390612e2381516001600160801b036060840151169061422f565b612e7f6001600160401b036020608081850151604051612e638482018093604080916001600160801b038151168452602081015160208501520151910152565b60608152612e718382610df2565b5190209401510151166143b0565b808214612dcf565b612ea26121ca602060a084015101516001600160401b031690565b5480612eae5750505f90565b60408201908151604051612eca81610d24602082019485612016565b5190201491821592612ee8575b505015612ee357600190565b600290565b6001600160801b03919250612f1e612f0f6020612f2a9301516001600160801b0390511690565b9351516001600160801b031690565b6001600160801b031690565b911610155f80612ed7565b15612f3c57565b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b15612f6c5750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15612fa65750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15612fe05750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b906001600160401b03809116911601906001600160401b03821161304657565b613012565b9060038202918083046003149015171561304657565b908160011b918083046002149015171561304657565b90602c820291808304602c149015171561304657565b15613096575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b905f6101008301928351519161311560e08301936130f46105f1865161ffff1690565b8091149081613310575b81613300575b816132f0575b509593949195612f35565b61311d6127a6565b50613126613902565b909561313e6007546001600160401b039060a01c1690565b945f905f985f975f96600654995b61315b6105f18d5161ffff1690565b89101561328b5761318161317d6131778b6101608e0151612810565b51151590565b1590565b61327b5761319d6131938a8751612810565b5163ffffffff1690565b809d61325d575b505060019b958989610120820151906131bc91612810565b5163ffffffff168a8d828c60208a019b828d516131da9061ffff1690565b61ffff1663ffffffff821610906131f091612f9e565b600163ffffffff84161b906132088482841615612f64565b179b516132169061ffff1690565b90613220936139c8565b91909361014001519061323291612810565b51149061323e91612fd8565b61324791613026565b976105f1600161315b925b01999791505061314c565b63ffffffff613274921663ffffffff821611612f64565b5f8c6131a4565b95976105f1600161315b92613252565b50985099505091965050610e239392506132eb91506132cf87876132b76001600160401b03821661304b565b6132c96001600160401b038416613061565b1061308d565b60606020604085015151015101519051610160840151916143f6565b6144a0565b905061016084015151145f61310a565b6101408501515181149150613104565b61012085015151811491506130fe565b805f525f60205260ff6133478360405f20906001600160a01b03165f5260205260405f2090565b54166133b257805f525f6020526133728260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916600117905533916001600160a01b0316907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b8051519081156133b2576133cb826127ca565b915f5b82518051821015613447579061343661343160206133ee84600196612810565b51015161342c6001600160401b036040613409878b51612810565b510151166040519261341a84610d81565b83526001600160401b03166020830152565b614520565b61460a565b6134408287612810565b52016133ce565b505090505f610c509261465a565b61345e816133b8565b90518051908115612f3c5760b4821161365d575f9361349f61348f61348a61348586613077565b613959565b6144f8565b9361349985615572565b846155ba565b6134b56134ae61343186615804565b84602e0152565b6134be836155ca565b60305f955b835187101561356f5761356484939261355f6001936135486134e68c809a612810565b519161352763ffffffff61351e604086019361351861350c86516001600160401b031690565b6001600160401b031690565b906139bb565b9b16868d615596565b61354161353386613967565b91516001600160401b031690565b908b61561e565b602061355384613975565b91015190890160200152565b613983565b9601959091926134c3565b61364094929650602093506135e69150946135a26001600160401b03610e239761359b82821115612f35565b16846155d8565b6135e16135b86135b285846146f4565b9461486d565b6001600160a01b031673ffffffffffffffffffffffffffffffffffffffff196007541617600755565b600655565b6136376135fa82516001600160401b031690565b7fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6007549260a01b16911617600755565b015161ffff1690565b61ffff60e01b1961ffff60e01b6007549260e01b16911617600755565b7fab7bac08000000000000000000000000000000000000000000000000000000005f52600482905260b460245260445b5ffd5b906009548210156128245760095f52600282901c7f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af019160031b60181690565b6009549068010000000000000000821015610d9c57600182016009556009548210156128245760095f5260205f208260021c01916001600160401b038060c085549360061b169316831b921b1916179055565b6001600160401b0381165f5260086020526001600160a01b0380600160405f200154166138606006546137af6007546137a4868216916137946137786001600160401b038360a01c169261ffff9060e01c1690565b93613781610e54565b9687526001600160a01b03166020870152565b6001600160401b03166040850152565b61ffff166060830152565b6137ca856001600160401b03165f52600860205260405f2090565b6060600161ffff928451815501926001600160a01b0360208201511673ffffffffffffffffffffffffffffffffffffffff1985541617845560408101517fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b8087549360a01b1616911617845501511661ffff60e01b1961ffff60e01b83549260e01b169116179055565b16156138695750565b610e23906136d0565b805f525f60205260ff6138998360405f20906001600160a01b03165f5260205260405f2090565b5416156133b257805f525f6020526138c58260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916905533916001600160a01b0316907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b61390a6127a6565b5061395560065461ffff6007546040519261392484610dbc565b83526001600160a01b03811660208401526001600160401b038160a01c16604084015260e01c16606082015261495a565b9091565b603001908160301161304657565b906004820180921161304657565b90600c820180921161304657565b90602c820180921161304657565b600101908160011161304657565b602401908160241161304657565b906001820180921161304657565b9190820180921161304657565b9091939263ffffffff61ffff91169416841015613a54576139eb61348585613077565b9363ffffffff6020868501015160e01c1603613a295750610c5090613a216020613a1486613967565b8301015160c01c94613975565b016020015190565b7f4724a0fd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b83907ff81f5140000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b999093969998929897919597613a998b611a74565b8a15613ada5761368d8b613aac81611a74565b7f112d89cc000000000000000000000000000000000000000000000000000000005f5260ff16600452602490565b869798999a506020613b0c9188613b1d94959697989915159081613c8c575b613b069198979698614a0b565b01614a49565b613b16368b611086565b908761568b565b5f5f935b878510613bef575b50613b349350614bf0565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001691823b156101c557613ba15f95604051988996879586957f87d3a9b100000000000000000000000000000000000000000000000000000000875260048701614fbd565b03915afa908115610c0757610c5092613bc192613bd5575b503690611086565b6001600160801b03633b9aca009151160490565b80613be35f613be993610df2565b8061041e565b5f613bb9565b90613c2561317d613c14613c0d613c07898d8c614a53565b8061285b565b3691614a75565b613c1f368888614a75565b9061575d565b613c805750613c6290613c4c61084b613c42613b34978b8a614a53565b6020810190612829565b90815181518082149182613c6a575b5050614ae0565b60015f613b29565b9091506020840120906020830120145f80613c5b565b60019094019390613b21565b61ffff8111159150613af9565b9260208091613d08612e7f95613cba610e23996001600160401b039761422f565b604051613ce78582018093604080916001600160801b038151168452602081015160208501520151910152565b60608152613cf6608082610df2565b519020612e7f86858a510151166143b0565b604051613d358382018093604080916001600160801b038151168452602081015160208501520151910152565b60608152612e71608082610df2565b9190915f92608081019384515193613d8c6060840195613d696105f1885161ffff1690565b8091149081613f79575b81613f6a575b81613f5b575b5090869391959295612f35565b6020828101510151613da6906001600160401b03166150d7565b90613daf6127a6565b50613db98261495a565b939096613dd060408501516001600160401b031690565b905f925f945f613de66105f15f9b5161ffff1690565b8a1015613f1257908b949392918e613e0961317d6131778f8f9060e00151612810565b613eff576131938c613e1b9251612810565b8098613ee1575b50508a600197968b8060a084015190613e3a91612810565b5163ffffffff16918260208a0151613e539061ffff1690565b61ffff1663ffffffff82161090613e6991612f9e565b600163ffffffff84161b90613e818482841615612f64565b1797828d8d519260200151613e979061ffff1690565b90613ea1936139c8565b91909360c0015190613eb291612810565b511490613ebe91612fd8565b613ec791613026565b986105f16001613de6925b019a929394959691908e61067e565b63ffffffff613ef8921663ffffffff821611612f64565b5f87613e22565b50959450986105f16001613de692613ed2565b509a509550975098915050610e23949350613f569150613f3f88886132b76001600160401b03821661304b565b60606020845101510151905160e0850151916143f6565b6151fc565b905060e085015151145f613d7f565b60c08601515181149150613d79565b60a08601515181149150613d73565b15613f91575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b5f1981019190821161304657565b9190820391821161304657565b15613fe3575050565b7f12e74add000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b156140195750565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f60015461405581612c42565b908160448501526001811690815f146140ec575060011461408c575b506140889192600319848303016024850152610c1b565b0390fd5b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b8183106140d25750919291508101606401614088614071565b8054606484880101528594506020909201916001016140b9565b60ff191660648086019190915291151560051b840190910191506140889050614071565b9392919093156141205750505050565b9160ff8092816084969581604051977fd382033a000000000000000000000000000000000000000000000000000000008952166004880152166024860152166044840152166064820152fd5b15614175575050565b9063ffffffff80927f73b1bde3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b156141b6575050565b9063ffffffff80927f1eb5c1eb000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b156141f7575050565b9063ffffffff80927f87f10409000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6142496001600160801b03610e239316633b9aca00900490565b614257814242821115613f88565b61437560e06142668342613fcd565b9361436a6004546142946142818263ffffffff9060501c1690565b9663ffffffff8816988942911115613fda565b6142c78351805160208201207f000000000000000000000000000000000000000000000000000000000000000014614011565b61430f60208401516142da815160ff1690565b6002549160ff831660ff811660ff8416149384614383575b6143099060209060081c60ff165b93015160ff1690565b93614110565b614335614323606085015163ffffffff1690565b63ffffffff838116908216811461416c565b614356614349608085015163ffffffff1690565b9160201c63ffffffff1690565b63ffffffff811663ffffffff8316146141ad565b015163ffffffff1690565b9163ffffffff8316146141ee565b9350614309602061430061439a8286015160ff1690565b60ff600889901c811691161496925050506142f2565b6001600160401b03165f52600560205260405f205480156143ce5790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b92916144058251825114612f35565b5f5b82518110156144995761441a8183612810565b51156144915763ffffffff61442f8285612810565b511661443e8187518110612f9e565b6144488187612810565b51516004811015610236576001190161446657506001905b01614407565b7f5b5c80eb000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b600190614460565b5050509050565b610e239060408101519061ffff60e082015116608082015160a08301519060c08401519261016061014086015195015195615320565b604051608091906144e78382610df2565b6041815291601f1901366020840137565b9061450282610e72565b61450f6040519182610df2565b82815280926127f2601f1991610e72565b602460208201906001600160401b03825116806145bd575b50614545614573916144f8565b9261456a61456461455e61455887615490565b8761549c565b866154b3565b856154ca565b905190846154f7565b9061458861350c82516001600160401b031690565b61459157505090565b6145b261350c6145a46145b994866154e0565b92516001600160401b031690565b908361550f565b5090565b600191505b60808110156145ea57506145456145e36145de61457393613991565b61399f565b9150614538565b60019060071c9101906145c2565b805191908290602001825e015f815290565b80516001810180911161304657614620906144f8565b8051156128245761464a8161463d6020945f868196015382615548565b50604051918280926145f8565b039060025afa15610c07575f5190565b909291928084039380851161304657600185146146e35760015b8060011b90868210156146875750614674565b9192939495505082019182811161304657826146a3918561465a565b916146ae929361465a565b6146b66144d6565b9182511561282457825f9261464a92600160208097015360218301526041820152604051918280926145f8565b50906146f0929350612810565b5190565b9190916146ff6127a6565b506030835110613a29576356414c34602084015160e01c03613a2957602483015160c01c602c84015160f01c9361478a602e82015192604e83015160f01c93614758614749610e54565b6001600160401b039093168352565b61476a6020830198899061ffff169052565b60408201526147816060820194859061ffff169052565b955161ffff1690565b9061ffff82169283159384156147e3575b5083156147c9575b5082156147b4575b5050613a295750565b5191506147c090615666565b14155f806147ab565b519092506147da9061ffff166105f1565b1515915f6147a3565b60b41093505f61479b565b600a907fffff000000000000000000000000000000000000000000000000000000000000610c5094937f610000000000000000000000000000000000000000000000000000000000000083521660018201527f80600a3d393df300000000000000000000000000000000000000000000000000600382015201906145f8565b604051905f60208301526148968261488860218201846145f8565b03601f198101845283610df2565b61600082511161492e5750610d246148f06148de6148b6845161ffff1690565b60f01b7fffff0000000000000000000000000000000000000000000000000000000000001690565b926040519283916020830195866147ee565b51905ff0906001600160a01b0382161561490657565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b6149626127a6565b50602081016001600160a01b03815116156149f7576001600160a01b0390511661499361ffff606084015116615666565b813b60018111156149e4575f1981019081116130465781116149d157908160016149bf6149ce946144f8565b9260208401903c8092516146f4565b91565b50633610565160e21b5f5260045260245ffd5b82633610565160e21b5f5260045260245ffd5b5051633a517eed60e21b5f5260045260245ffd5b15614a135750565b7fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b35610c5081610f23565b91908110156128245760051b81013590603e19813603018212156101c5570190565b929190614a8181611115565b93614a8f6040519586610df2565b602085838152019160051b8101918383116101c55781905b838210614ab5575050505050565b81356001600160401b0381116101c557602091614ad58784938701610ec3565b815201910190614aa7565b91909115614aec575050565b90614088614b2e926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610c1b565b83810360031901602485015290610c1b565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156101c55701602081359101916001600160401b0382116101c55781360383136101c557565b90602083828152019260208260051b82010193835f925b848410614bb85750505050505090565b909192939495602080614be0600193601f19868203018852614bda8b88614b60565b90614b40565b9801940194019294939190614ba8565b91909115614bfc575050565b6140886040519283927ffef760c7000000000000000000000000000000000000000000000000000000008452602060048501526024840191614b91565b9035601e19823603018112156101c55701602081359101916001600160401b0382116101c5578160051b360383136101c557565b9035607e19823603018112156101c5570190565b359060038210156101c557565b9035605e19823603018112156101c5570190565b90602083828152019260208260051b82010193835f925b848410614cc95750505050505090565b909192939495602080614d3b600193601f19868203018852614ceb8b88614c8e565b90614cf582614c81565b614cfe8161022c565b8152614d2d614d22614d1286850185614b60565b6060888601526060850191614b40565b926040810190614b60565b916040818503910152614b40565b9801940194019294939190614cb9565b610c5091614e15614e0a614d8e614d73614d658680614b60565b608087526080870191614b40565b614d806020870187614b60565b908683036020880152614b40565b6080614dfa614da06040880188614c6d565b8684036040880152614db181614c81565b614dba8161022c565b8452614dc860208201614c81565b614dd18161022c565b6020850152614de260408201614c81565b614deb8161022c565b60408501526060810190614b60565b9190928160608201520191614b40565b926060810190614c39565b916060818503910152614ca2565b90602083828152019260208260051b82010193835f925b848410614e4a5750505050505090565b909192939495601f198282030184528635601e19843603018112156101c557830190614e7a602082019280614c39565b8091936020845252604082019060408160051b8401019380935f915b838310614eb9575050505050506020806001929801940194019294939190614e3a565b909192939495603f19838203018652614ed28783614c8e565b8035614edd81610fa3565b614ee681611a74565b8252614f09614ef86020830183614c6d565b606060208501526060840190614d4b565b906040810135609e19823603018112156101c5576001936020938493614faf9301916040818303910152614fa1614f81614f54614f468580614b60565b60a0865260a0860191614b40565b86850135614f6181610f8e565b151587850152614f746040860186614c6d565b8482036040860152614d4b565b926060810135614f9081610f8e565b151560608401526080810190614c6d565b906080818403910152614d4b565b980196019493019190614e96565b9092809492969593606083019083526060602084015252608081019560808560051b83010194815f90603e1981360301995b83831061500e575050505050610c509495506040818503910152614e23565b9091929397607f1986820301825288358b8112156101c5576020615068600193868394019061505b6150516150438480614c39565b604085526040850191614b91565b9285810190614b60565b9185818503910152614b40565b9a019201930191909392614fef565b1561507e57565b633a517eed60e21b5f525f60045260245ffd5b9060405161509e81610dbc565b606061ffff600183958054855201546001600160a01b03811660208501526001600160401b038160a01c16604085015260e01c16910152565b6150df6127a6565b5060095480156151e6576150f35f91613fbf565b808210615194575061515b9161513f6001600160401b0361512c61511961515695613690565b90546001600160401b039160031b1c1690565b92166001600160401b0383161115615077565b6001600160401b03165f52600860205260405f2090565b615091565b906001600160a01b0361517860208401516001600160a01b031690565b161561518057565b8151633a517eed60e21b5f5260045260245ffd5b906151b06151aa6151a584846139bb565b6139ad565b60011c90565b906151bd61511983613690565b6001600160401b038581169116116151d65750906150f3565b91506151e190613fbf565b6150f3565b633a517eed60e21b5f5261368d6024905f600452565b610e239161ffff606082015116815160208301519060408401519260e060c086015195015195615320565b908160209103126101c55751610c5081610f8e565b9796959493929161ffff8992168252602082015f905b600882106152d7575050509361529e6040959461528a6101e09561527f6152ad966101208b9a0190611e05565b6101608c0190611e05565b6102406101a08b01526102408a01906104f1565b908882036101c08a0152611e65565b9501926001600160401b0381511684526001600160401b0360208201511660208501520151910152565b829350602080916001939451815201930191018992615252565b156152f857565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b9560209561541b93959294975f6040805161533a81610da1565b828152828b82015201525192615387888501519461537761350c604061536789516001600160401b031690565b935101516001600160401b031690565b6001600160401b038216146157c7565b83516001600160401b0316936153e260406153b46153ab8c85015163ffffffff1690565b63ffffffff1690565b92015151916153d36153c4610e63565b6001600160401b039098168852565b6001600160401b0316868b0152565b604085015260405198899788977f7d1a86950000000000000000000000000000000000000000000000000000000089526004890161523c565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af18015610c0757610e23915f91615461575b506152f1565b615483915060203d602011615489575b61547b8183610df2565b810190615227565b5f61545b565b503d615471565b6020600a910153600190565b602082602292010153600181018091116130465790565b602082600a92010153600181018091116130465790565b6020828192010153600181018091116130465790565b602082601092010153600181018091116130465790565b81602091939293010152602081018091116130465790565b9092919083016020015b608082101561552d57906001929391530190565b600180916080607f85161781530193019060071c9092615519565b90805191821561556a576021602084930191015e600101806001116130465790565b505050600190565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602d908260081c602c8201530153565b604f5f9182604e8201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b61ffff16602c810290808204602c149015171561304657603001806030116130465790565b906156da90612e7f60405160208101906156c28288604080916001600160801b038151168452602081015160208501520151910152565b606081526156d1608082610df2565b519020916143b0565b602082015180820361572f5750506001600160801b03633b9aca0091511604615707814242821115613f88565b4203428111613046576001600160801b03610e23911663ffffffff600454169081811061583c565b7f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b9081518151036133b2575f5b825181101561556a5761577c8184612810565b51516157888284612810565b5151036157c0576157998184612810565b51602081519101206157ab8284612810565b5160208151910120036157c057600101615769565b5050505f90565b156157cf5750565b6001600160401b03907f90f4dbed000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b604051906060906158158284610df2565b602283526145b991600a906020850190601f190136823753602060218401536002836154f7565b15615845575050565b906001600160801b0380927fdb020436000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffdfea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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

// ReAnchorPinnedSet is a paid mutator transaction binding the contract method 0x2fe76c3c.
//
// Solidity: function reAnchorPinnedSet(bytes reAnchorMsg) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactor) ReAnchorPinnedSet(opts *bind.TransactOpts, reAnchorMsg []byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.contract.Transact(opts, "reAnchorPinnedSet", reAnchorMsg)
}

// ReAnchorPinnedSet is a paid mutator transaction binding the contract method 0x2fe76c3c.
//
// Solidity: function reAnchorPinnedSet(bytes reAnchorMsg) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) ReAnchorPinnedSet(reAnchorMsg []byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.ReAnchorPinnedSet(&_ContractGroth16ICS07Tendermint.TransactOpts, reAnchorMsg)
}

// ReAnchorPinnedSet is a paid mutator transaction binding the contract method 0x2fe76c3c.
//
// Solidity: function reAnchorPinnedSet(bytes reAnchorMsg) returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorSession) ReAnchorPinnedSet(reAnchorMsg []byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.ReAnchorPinnedSet(&_ContractGroth16ICS07Tendermint.TransactOpts, reAnchorMsg)
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

// ContractGroth16ICS07TendermintClientFrozenIterator is returned from FilterClientFrozen and is used to iterate over the raw logs and unpacked data for ClientFrozen events raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintClientFrozenIterator struct {
	Event *ContractGroth16ICS07TendermintClientFrozen // Event containing the contract specifics and raw log

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
func (it *ContractGroth16ICS07TendermintClientFrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractGroth16ICS07TendermintClientFrozen)
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
		it.Event = new(ContractGroth16ICS07TendermintClientFrozen)
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
func (it *ContractGroth16ICS07TendermintClientFrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractGroth16ICS07TendermintClientFrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractGroth16ICS07TendermintClientFrozen represents a ClientFrozen event raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintClientFrozen struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterClientFrozen is a free log retrieval operation binding the contract event 0x59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c.
//
// Solidity: event ClientFrozen()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) FilterClientFrozen(opts *bind.FilterOpts) (*ContractGroth16ICS07TendermintClientFrozenIterator, error) {

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.FilterLogs(opts, "ClientFrozen")
	if err != nil {
		return nil, err
	}
	return &ContractGroth16ICS07TendermintClientFrozenIterator{contract: _ContractGroth16ICS07Tendermint.contract, event: "ClientFrozen", logs: logs, sub: sub}, nil
}

// WatchClientFrozen is a free log subscription operation binding the contract event 0x59869f4e64ba4779d80f0c4bd28741de9276b991a817653df9563ee66d59652c.
//
// Solidity: event ClientFrozen()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) WatchClientFrozen(opts *bind.WatchOpts, sink chan<- *ContractGroth16ICS07TendermintClientFrozen) (event.Subscription, error) {

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.WatchLogs(opts, "ClientFrozen")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractGroth16ICS07TendermintClientFrozen)
				if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "ClientFrozen", log); err != nil {
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
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) ParseClientFrozen(log types.Log) (*ContractGroth16ICS07TendermintClientFrozen, error) {
	event := new(ContractGroth16ICS07TendermintClientFrozen)
	if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "ClientFrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractGroth16ICS07TendermintClientUnfrozenIterator is returned from FilterClientUnfrozen and is used to iterate over the raw logs and unpacked data for ClientUnfrozen events raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintClientUnfrozenIterator struct {
	Event *ContractGroth16ICS07TendermintClientUnfrozen // Event containing the contract specifics and raw log

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
func (it *ContractGroth16ICS07TendermintClientUnfrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractGroth16ICS07TendermintClientUnfrozen)
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
		it.Event = new(ContractGroth16ICS07TendermintClientUnfrozen)
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
func (it *ContractGroth16ICS07TendermintClientUnfrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractGroth16ICS07TendermintClientUnfrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractGroth16ICS07TendermintClientUnfrozen represents a ClientUnfrozen event raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintClientUnfrozen struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterClientUnfrozen is a free log retrieval operation binding the contract event 0x8e0da210a2ffecc744373b4860cd9b8dc471810f5fb61b7d2c3234ccd4aae30d.
//
// Solidity: event ClientUnfrozen()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) FilterClientUnfrozen(opts *bind.FilterOpts) (*ContractGroth16ICS07TendermintClientUnfrozenIterator, error) {

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.FilterLogs(opts, "ClientUnfrozen")
	if err != nil {
		return nil, err
	}
	return &ContractGroth16ICS07TendermintClientUnfrozenIterator{contract: _ContractGroth16ICS07Tendermint.contract, event: "ClientUnfrozen", logs: logs, sub: sub}, nil
}

// WatchClientUnfrozen is a free log subscription operation binding the contract event 0x8e0da210a2ffecc744373b4860cd9b8dc471810f5fb61b7d2c3234ccd4aae30d.
//
// Solidity: event ClientUnfrozen()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) WatchClientUnfrozen(opts *bind.WatchOpts, sink chan<- *ContractGroth16ICS07TendermintClientUnfrozen) (event.Subscription, error) {

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.WatchLogs(opts, "ClientUnfrozen")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractGroth16ICS07TendermintClientUnfrozen)
				if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "ClientUnfrozen", log); err != nil {
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
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) ParseClientUnfrozen(log types.Log) (*ContractGroth16ICS07TendermintClientUnfrozen, error) {
	event := new(ContractGroth16ICS07TendermintClientUnfrozen)
	if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "ClientUnfrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractGroth16ICS07TendermintClientUpdatedIterator is returned from FilterClientUpdated and is used to iterate over the raw logs and unpacked data for ClientUpdated events raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintClientUpdatedIterator struct {
	Event *ContractGroth16ICS07TendermintClientUpdated // Event containing the contract specifics and raw log

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
func (it *ContractGroth16ICS07TendermintClientUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractGroth16ICS07TendermintClientUpdated)
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
		it.Event = new(ContractGroth16ICS07TendermintClientUpdated)
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
func (it *ContractGroth16ICS07TendermintClientUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractGroth16ICS07TendermintClientUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractGroth16ICS07TendermintClientUpdated represents a ClientUpdated event raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintClientUpdated struct {
	Height uint64
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterClientUpdated is a free log retrieval operation binding the contract event 0x9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb31.
//
// Solidity: event ClientUpdated(uint64 height)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) FilterClientUpdated(opts *bind.FilterOpts) (*ContractGroth16ICS07TendermintClientUpdatedIterator, error) {

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.FilterLogs(opts, "ClientUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractGroth16ICS07TendermintClientUpdatedIterator{contract: _ContractGroth16ICS07Tendermint.contract, event: "ClientUpdated", logs: logs, sub: sub}, nil
}

// WatchClientUpdated is a free log subscription operation binding the contract event 0x9d676bcc5dc254c2c765d65800ad98819df9c1e898a93456b968f7d206cfeb31.
//
// Solidity: event ClientUpdated(uint64 height)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) WatchClientUpdated(opts *bind.WatchOpts, sink chan<- *ContractGroth16ICS07TendermintClientUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.WatchLogs(opts, "ClientUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractGroth16ICS07TendermintClientUpdated)
				if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "ClientUpdated", log); err != nil {
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
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) ParseClientUpdated(log types.Log) (*ContractGroth16ICS07TendermintClientUpdated, error) {
	event := new(ContractGroth16ICS07TendermintClientUpdated)
	if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "ClientUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractGroth16ICS07TendermintPinnedSetReAnchoredIterator is returned from FilterPinnedSetReAnchored and is used to iterate over the raw logs and unpacked data for PinnedSetReAnchored events raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintPinnedSetReAnchoredIterator struct {
	Event *ContractGroth16ICS07TendermintPinnedSetReAnchored // Event containing the contract specifics and raw log

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
func (it *ContractGroth16ICS07TendermintPinnedSetReAnchoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractGroth16ICS07TendermintPinnedSetReAnchored)
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
		it.Event = new(ContractGroth16ICS07TendermintPinnedSetReAnchored)
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
func (it *ContractGroth16ICS07TendermintPinnedSetReAnchoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractGroth16ICS07TendermintPinnedSetReAnchoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractGroth16ICS07TendermintPinnedSetReAnchored represents a PinnedSetReAnchored event raised by the ContractGroth16ICS07Tendermint contract.
type ContractGroth16ICS07TendermintPinnedSetReAnchored struct {
	Height         uint64
	ValidatorsHash [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPinnedSetReAnchored is a free log retrieval operation binding the contract event 0x16f9bd20fd25a0828c55e0202f115427e910995861d180b9cf422f347cc01575.
//
// Solidity: event PinnedSetReAnchored(uint64 height, bytes32 validatorsHash)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) FilterPinnedSetReAnchored(opts *bind.FilterOpts) (*ContractGroth16ICS07TendermintPinnedSetReAnchoredIterator, error) {

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.FilterLogs(opts, "PinnedSetReAnchored")
	if err != nil {
		return nil, err
	}
	return &ContractGroth16ICS07TendermintPinnedSetReAnchoredIterator{contract: _ContractGroth16ICS07Tendermint.contract, event: "PinnedSetReAnchored", logs: logs, sub: sub}, nil
}

// WatchPinnedSetReAnchored is a free log subscription operation binding the contract event 0x16f9bd20fd25a0828c55e0202f115427e910995861d180b9cf422f347cc01575.
//
// Solidity: event PinnedSetReAnchored(uint64 height, bytes32 validatorsHash)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) WatchPinnedSetReAnchored(opts *bind.WatchOpts, sink chan<- *ContractGroth16ICS07TendermintPinnedSetReAnchored) (event.Subscription, error) {

	logs, sub, err := _ContractGroth16ICS07Tendermint.contract.WatchLogs(opts, "PinnedSetReAnchored")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractGroth16ICS07TendermintPinnedSetReAnchored)
				if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "PinnedSetReAnchored", log); err != nil {
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

// ParsePinnedSetReAnchored is a log parse operation binding the contract event 0x16f9bd20fd25a0828c55e0202f115427e910995861d180b9cf422f347cc01575.
//
// Solidity: event PinnedSetReAnchored(uint64 height, bytes32 validatorsHash)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintFilterer) ParsePinnedSetReAnchored(log types.Log) (*ContractGroth16ICS07TendermintPinnedSetReAnchored, error) {
	event := new(ContractGroth16ICS07TendermintPinnedSetReAnchored)
	if err := _ContractGroth16ICS07Tendermint.contract.UnpackLog(event, "PinnedSetReAnchored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
