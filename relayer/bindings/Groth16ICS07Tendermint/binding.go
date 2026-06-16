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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membership_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviour_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"updateClient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_clientState\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCachedValidatorSet\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"updateClientMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CannotHandleMisbehavior\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientNotFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ClientStateMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ClockDriftMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidMembershipProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProofHeightMismatch\",\"inputs\":[{\"name\":\"expectedRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"expectedRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustThresholdMismatch\",\"inputs\":[{\"name\":\"expectedNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expectedDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnbondingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"UnknownZkAlgorithm\",\"inputs\":[{\"name\":\"algorithm\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorCountExceedsLimit\",\"inputs\":[{\"name\":\"validatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxValidatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"VerificationKeyMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x610160604052346100fe57616ef0803803809161001b82610116565b6101603960e0816101600191126100fe576100346101a1565b61003f6101806101b8565b61004a6101a06101b8565b6100556101c06101b8565b6101e0516001600160401b0381116100fe578561017f820112156100fe576100a2958161018061008b93610160015191016101e7565b91610200519361009c6102206101b8565b956107fd565b604051615e459081610feb823960805181614644015260a05181613aeb015260c05181505060e05181612a820152610100518181816109ac01526109f301526101205181612bcd015261014051818181612a4d01526140f60152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b610160601f91909101601f19168101906001600160401b0382119082101761013d57604052565b610102565b604081019081106001600160401b0382111761013d57604052565b601f909101601f19168101906001600160401b0382119082101761013d57604052565b604051906101906101008361015d565b565b6040519061019060408361015d565b61016051906001600160a01b03821682036100fe57565b51906001600160a01b03821682036100fe57565b6001600160401b03811161013d57601f01601f191660200190565b9291926101f3826101cc565b91610201604051938461015d565b8294818452818301116100fe578281602093845f96015e010152565b9080601f830112156100fe578151610237926020016101e7565b90565b519060ff821682036100fe57565b91908260409103126100fe5760405161026081610142565b60206102798183956102718161023a565b85520161023a565b910152565b51906001600160401b03821682036100fe57565b91908260409103126100fe576040516102aa81610142565b60206102798183956102bb8161027e565b85520161027e565b519063ffffffff821682036100fe57565b519081151582036100fe57565b519060028210156100fe57565b6020818303126100fe578051906001600160401b0382116100fe5701610140818303126100fe5761031d610180565b8151909290916001600160401b0383116100fe5761036482610347610120946103b496850161021d565b86526103568160208501610248565b602087015260608301610292565b604085015261037560a082016102c3565b606085015261038660c082016102c3565b608085015261039760e082016102d4565b60a08501526103a961010082016102e1565b60c0850152016102c3565b60e082015290565b90600182811c921680156103ea575b60208310146103d657565b634e487b7160e01b5f52602260045260245ffd5b91607f16916103cb565b601f821161040157505050565b5f5260205f20906020601f840160051c83019310610439575b601f0160051c01905b81811061042e575050565b5f8155600101610423565b909150819061041a565b6002111561044d57565b634e487b7160e01b5f52602160045260245ffd5b90600281101561044d5769ff00000000000000000082549160481b169069ff0000000000000000001916179055565b80518051906001600160401b03821161013d576104b9826104b26001546103bc565b60016103f4565b602090601f831160011461065a57926104f28361061f9460e094610190975f9261064f575b50508160011b915f199060031b1c19161790565b6001555b6020818101518051600280549284015161ff0060089190911b1660ff90921661ffff199093169290921717905560408083015180516003805492909401516fffffffffffffffff0000000000000000931b929092166001600160401b039092166001600160801b031990911617179055606081015161058c9063ffffffff1660049063ffffffff1663ffffffff19825416179055565b6105c46105a0608083015163ffffffff1690565b60049067ffffffff0000000082549160201b169067ffffffff000000001916179055565b6105fc6105d460a0830151151590565b60049068ff0000000000000000825491151560401b169068ff00000000000000001916179055565b61061460c082015161060d81610443565b6004610461565b015163ffffffff1690565b6004906dffffffff0000000000000000000082549160501b16906dffffffff000000000000000000001916179055565b015190505f806104de565b60015f52601f19831691905f516020616eb05f395f51905f52925f5b8181106106bb57509360e09361019096936001938361061f98106106a3575b505050811b016001556104f6565b01515f1960f88460031b161c191690555f8080610695565b92936020600181928786015181550195019301610676565b604051905f82600154916106e6836103bc565b8083529260018116908115610756575060011461070a575b6101909250038361015d565b5060015f90815290915f516020616eb05f395f51905f525b81831061073a575050906020610190928201016106fe565b6020919350806001915483858901015201910190918492610722565b6020925061019094915060ff191682840152151560051b8201016106fe565b1561077e575050565b6355bace6f60e01b5f9081526001600160401b039182166004529116602452604490fd5b634e487b7160e01b5f52601160045260245ffd5b9063ffffffff8091169116019063ffffffff82116107d057565b6107a2565b156107de575050565b9063ffffffff80926333fae18560e21b5f52166004521660245260445ffd5b919361083e61084391969294967f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a6123361010052602080825183010191016102ee565b610490565b61084b6106d3565b60208151910120610120526108666108616106d3565b610984565b610140526108d96108c061089360206108856108806106d3565b610a5e565b01516001600160401b031690565b600354906108b1906001600160401b03808416919081168214610775565b60401c6001600160401b031690565b6001600160401b03165f90815260056020526040902090565b556001600160a01b0390811660805290811660a05290811660c0521660e05260045461093c9063ffffffff811661092a61091d605084901c63ffffffff16836107b6565b9260201c63ffffffff1690565b9163ffffffff808416911611156107d5565b6001600160a01b0381166109635750610953610d42565b5061096061010051610dc4565b50565b8061097061096092610c44565b5061097a81610cba565b5061010051610e1d565b61098d90610ec8565b8051600181018091116107d0576109a390610e96565b8051156109eb57602091816109bf845f94019284845382610f95565b50604051918291518091835e8101838152039060025afa156109e0575f5190565b6040513d5f823e3d90fd5b634e487b7160e01b5f52603260045260245ffd5b60405190610a0c82610142565b5f602083606081520152565b80156107d0575f190190565b5f198101919082116107d057565b919082039182116107d057565b9081518110156109eb570160200190565b90600182018092116107d057565b610a666109ff565b5080518015908115610c38575b50610c29575f19908051805b610bd9575b505f198214610bbb57600360fc1b6001600160f81b0319610abe610ab0610aaa86610a50565b85610a3f565b516001600160f81b03191690565b161480610bc5575b610bbb575f90610ad583610a50565b915b8151831015610b6c57610af6610af0610ab08585610a3f565b60f81c90565b60ff811660308110908115610b61575b50610b5457600a82026001600160401b03908116602f1990920160ff1691909101811691168110610b3c57600190920191610ad7565b50915050610b48610192565b9081525f602082015290565b5050915050610b48610192565b60399150115f610b06565b9150916001811190811591610baf575b50610ba05761023790610b8d610192565b9283526001600160401b03166020830152565b63384c687760e21b5f5260045ffd5b602b915010155f610b7c565b9050610b48610192565b506002610bd3838351610a32565b11610ac6565b602d60f81b610c03610bf6610ab0610bf085610a24565b86610a3f565b6001600160f81b03191690565b14610c1757610c1190610a18565b80610a7f565b610c22919250610a24565b905f610a84565b6329120bff60e21b5f5260045ffd5b6040915010155f610a73565b6001600160a01b0381165f9081525f516020616ed05f395f51905f52602052604090205460ff16610cb5576001600160a01b03165f8181525f516020616ed05f395f51905f5260205260408120805460ff191660011790553391905f516020616e305f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f516020616e505f395f51905f52602052604090205460ff16610cb5576001600160a01b0381165f9081525f516020616e505f395f51905f5260205260409020805460ff1916600117905533906001600160a01b03165f516020616e905f395f51905f525f516020616e305f395f51905f525f80a4600190565b5f80525f516020616e505f395f51905f526020525f516020616e705f395f51905f525460ff16610dc0575f8080525f516020616e505f395f51905f526020525f516020616e705f395f51905f52805460ff1916600117905533905f516020616e905f395f51905f525f516020616e305f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f14610cb5575f818152602081815260408083208380529091528120805460ff1916600117905533915f516020616e305f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16610e90575f818152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905533916001600160a01b0316905f516020616e305f395f51905f525f80a4600190565b50505f90565b90610ea0826101cc565b610ead604051918261015d565b8281528092610ebe601f19916101cc565b0190602036910137565b805115610f345780516001905b6080811015610f2657508060010190816001116107d057600190835101018091116107d057610f06610f2291610e96565b91610f1c610f1384610f50565b82519085610f5c565b83610fbf565b5090565b60019060071c910190610ed5565b50604051610f4360208261015d565b5f81525f36602083013790565b6020600a910153600190565b9092919083016020015b6080821015610f7a57906001929391530190565b600180916080607f85161781530193019060071c9092610f66565b908051918215610fb7576021602084930191015e600101806001116107d05790565b505050600190565b908092918251928315610fe357839260208092019201015e81018091116107d05790565b505050509056fe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a7146101145780630bece3561461010f578063248a9ca31461010a5780632f2ff15d1461010557806336568abe14610100578063536c2ad3146100fb5780636a28f000146100f65780638a8e4c5d146100f157806391d14854146100ec578063974a74c4146100e7578063a217fddf146100e2578063a6f031bb146100dd578063d547741f146100d8578063db3e1fa4146100d3578063ddba6537146100ce5763ef913a4b146100c9575f80fd5b610a84565b6109cf565b610995565b610966565b610859565b61083f565b6106df565b61069e565b61067f565b6105da565b61046f565b61033e565b610308565b6102b0565b61022b565b346101b55760203660031901126101b5576004357fffffffff0000000000000000000000000000000000000000000000000000000081168091036101b557807f7965db0b000000000000000000000000000000000000000000000000000000006020921490811561018b575b506040519015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150145f610180565b5f80fd5b9060206003198301126101b5576004356001600160401b0381116101b557826023820112156101b5578060040135926001600160401b0384116101b557602484830101116101b5576024019190565b634e487b7160e01b5f52602160045260245ffd5b6003111561022657565b610208565b346101b557602061029261023e366101b9565b9061025160ff60045460401c1615610b6f565b7fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5f525f845260405f205f8052845260ff60405f205416156102a357611a32565b6040519061029f8161021c565b8152f35b6102ab611dd4565b611a32565b346101b55760203660031901126101b55760206102da6004355f525f602052600160405f20015490565b604051908152f35b60409060031901126101b557600435906024356001600160a01b03811681036101b55790565b346101b55761033c610319366102e2565b90610337610332825f525f602052600160405f20015490565b611e43565b613525565b005b346101b55761034c366102e2565b336001600160a01b038216036103655761033c916135bd565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b90602080835192838152019201905f5b8181106103aa5750505090565b825163ffffffff1684526020938401939092019160010161039d565b90602080835192838152019201905f5b8181106103e35750505090565b82518452602093840193909201916001016103d6565b90602080835192838152019201905f5b8181106104165750505090565b82516001600160401b0316845260209384019390920191600101610409565b9161045e9061045061046c959360608652606086019061038d565b9084820360208601526103c6565b9160408184039101526103f9565b90565b346101b55760203660031901126101b55760043561048c8161364d565b6001600160a01b038116156105bb576104be906104b86104b360085461ffff9060e01c1690565b6136e5565b9061374d565b9060206104cb8383613806565b01906104eb6104e66104df845161ffff1690565b61ffff1690565b611c20565b906104fe6104e66104df855161ffff1690565b936105116104e66104df865161ffff1690565b9160095490600a54925f5b8861052c6104df8a5161ffff1690565b61ffff8316908110156105a5579161059e60019261059085806105898f8f998e8e8e8e61057b61ffff9f9661057387610568816105819b611c66565b9063ffffffff169052565b5161ffff1690565b91613948565b929095611c66565b528b611c66565b906001600160401b03169052565b011661051c565b50876105b78860405193849384610435565b0390f35b633a517eed60e21b5f52600482905260245b5ffd5b5f9103126101b557565b346101b5575f3660031901126101b557335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff16156106685760045460ff8160401c16156106405768ff00000000000000001916600455005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b346101b55761068d366101b9565b5050636d40ebe160e11b5f5260045ffd5b346101b557602060ff6106d36106b3366102e2565b905f525f845260405f20906001600160a01b03165f5260205260405f2090565b54166040519015158152f35b346101b55760203660031901126101b5576004356001600160401b0381116101b557806004019061016060031982360301126101b55761072760ff60045460401c1615610b6f565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610832575b6101448101906107898284611c7f565b90501561080a576107eb6107fa926105b7946107a86044850182611cb1565b6107b86064879493940183611cb1565b90608488013592610104890135956107cf87610dfc565b60a46107f26107e26101248d0189611cb1565b9b909a89611c7f565b3691610ce6565b9a0195613a32565b6040519081529081906020820190565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b61083a611dd4565b610779565b346101b5575f3660031901126101b55760206040515f8152f35b346101b55760203660031901126101b5576004356001600160401b0381116101b5578060040161014060031983360301126101b5576105b7916107fa916108a860ff60045460401c1615610b6f565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610959575b6109076044830182611cb1565b916109156064850182611cb1565b61010486013592916084870135919061092d85610dfc565b61093b610124890185611cb1565b97909660a46040519a61094f60208d610c5a565b5f8c520195613a32565b610961611dd4565b6108fa565b346101b55761033c610977366102e2565b90610990610332825f525f602052600160405f20015490565b6135bd565b346101b5575f3660031901126101b55760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346101b5576109dd366101b9565b50506109f160ff60045460401c1615610b6f565b7f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610a40575b636d40ebe160e11b5f5260045ffd5b610a4990611e43565b5f610a31565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b90602061046c928181520190610a4f565b346101b5575f3660031901126101b5576105b76040516020808201526101406040820152610b6381610ab96101808201611d28565b610ad560608301602060ff600254818116845260081c16910152565b610af760a0830160206001600160401b03600354818116845260401c16910152565b60045463ffffffff811660e0840152610b5590602081901c63ffffffff16610100850152610b30610120850160ff8360401c1615159052565b610b44610140850160ff8360481c16611dc7565b60501c63ffffffff16610160840152565b03601f198101835282610c5a565b60405191829182610a73565b15610b7657565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b03821117610bcd57604052565b610b9e565b606081019081106001600160401b03821117610bcd57604052565b608081019081106001600160401b03821117610bcd57604052565b60a081019081106001600160401b03821117610bcd57604052565b60c081019081106001600160401b03821117610bcd57604052565b61010081019081106001600160401b03821117610bcd57604052565b90601f801991011681019081106001600160401b03821117610bcd57604052565b60405190610c8b61010083610c5a565b565b60405190610c8b61026083610c5a565b60405190610c8b6101e083610c5a565b60405190610c8b608083610c5a565b60405190610c8b60c083610c5a565b6001600160401b038111610bcd57601f01601f191660200190565b929192610cf282610ccb565b91610d006040519384610c5a565b8294818452818301116101b5578281602093845f960137010152565b9080601f830112156101b55781602061046c93359101610ce6565b60ff8116036101b557565b91908260409103126101b557604051610d5a81610bb2565b60208082948035610d6a81610d37565b8452013591610d7883610d37565b0152565b6001600160401b038116036101b557565b3590610c8b82610d7c565b91908260409103126101b557604051610db081610bb2565b60208082948035610dc081610d7c565b8452013591610d7883610d7c565b63ffffffff8116036101b557565b3590610c8b82610dce565b801515036101b557565b3590610c8b82610de7565b600211156101b557565b3590610c8b82610dfc565b919091610140818403126101b557610e27610c7b565b928135916001600160401b0383116101b557610e6c82610e4f61012094610ebc968501610d1c565b8752610e5e8160208501610d42565b602088015260608301610d98565b6040860152610e7d60a08201610ddc565b6060860152610e8e60c08201610ddc565b6080860152610e9f60e08201610df1565b60a0860152610eb16101008201610e06565b60c086015201610ddc565b60e0830152565b6001600160801b038116036101b557565b3590610c8b82610ec3565b91908260609103126101b557604051610ef781610bd2565b60408082948035610f0781610ec3565b8452602081013560208501520135910152565b8092910391606083126101b557604051610f3381610bb2565b6040819483358352601f1901126101b5576020906040805193610f5585610bb2565b83810135610f6281610dce565b85520135828401520152565b6001600160401b038111610bcd5760051b60200190565b919060c0838203126101b557604051610f9d81610bed565b80938035610faa81610d7c565b82526020810135610fba81610dce565b6020830152610fcc8360408301610f1a565b604083015260a0810135906001600160401b0382116101b557019180601f840112156101b557823592610ffe84610f6e565b9361100c6040519586610c5a565b80855260208086019160051b830101918383116101b55760208101915b83831061103b57505050505060600152565b82356001600160401b0381116101b5578201906040828703601f1901126101b5576040519161106983610bb2565b602081013560048110156101b557835260408101356001600160401b0381116101b5576020910101906080828803126101b557604051926110a984610bed565b82356001600160401b0381116101b557886110c5918501610d1c565b845260208301356110d581610ec3565b602085015260408301356110e881610de7565b60408501526060830135936001600160401b0385116101b55761111089602096879601610d1c565b606082015283820152815201920191611029565b91906040838203126101b5576040519061113d82610bb2565b819380356001600160401b0381116101b55781016102c0818403126101b557611164610c8d565b9061116f8482610d98565b825260408101356001600160401b0381116101b55784611190918301610d1c565b60208301526111a160608201610d8d565b60408301526111b260808201610ed4565b60608301526111c360a08201610df1565b60808301526111d58460c08301610f1a565b60a08301526111e76101208201610df1565b60c083015261014081013560e08301526112046101608201610df1565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526112536102208201610df1565b6101c08301526102408101356101e08301526112726102608201610df1565b6102008301526102808101356102208301526102a0810135906001600160401b0382116101b5576112a591859101610d1c565b61024082015283526020810135916001600160401b0383116101b5576020926112ce9201610f85565b910152565b91906080838203126101b557604051906112ec82610bed565b819380356001600160401b0381116101b55760609261130c918301610d1c565b835260208101356020840152604081013561132681610d7c565b60408401520135908160070b82036101b55760600152565b9190916080818403126101b5576040519061135882610bed565b819381356001600160401b0381116101b557820181601f820112156101b557803561138281610f6e565b916113906040519384610c5a565b81835260208084019260051b820101918483116101b55760208201905b8382106113fe575050505083526113c660208301610df1565b60208401526040820135906001600160401b0382116101b557826113f3606094926112ce948694016112d3565b604086015201610d8d565b81356001600160401b0381116101b557602091611420888480948801016112d3565b8152019101906113ad565b919060a0838203126101b5576040519061144482610bed565b819380356001600160401b0381116101b55782611462918301611124565b835260208101356001600160401b0381116101b5578261148391830161133e565b60208401526114958260408301610d98565b60408401526080810135916001600160401b0383116101b5576060926112ce920161133e565b9080601f830112156101b557610100604051926114d88285610c5a565b839181019283116101b557905b8282106114f25750505090565b81358152602091820191016114e5565b9080601f830112156101b5576040519161151d604084610c5a565b8290604081019283116101b557905b8282106115395750505090565b813581526020918201910161152c565b359061ffff821682036101b557565b9080601f830112156101b557813561156f81610f6e565b9261157d6040519485610c5a565b81845260208085019260051b8201019283116101b557602001905b8282106115a55750505090565b6020809183356115b481610dce565b815201910190611598565b9080601f830112156101b55781356115d681610f6e565b926115e46040519485610c5a565b81845260208085019260051b8201019283116101b557602001905b82821061160c5750505090565b81358152602091820191016115ff565b9080601f830112156101b557813561163381610f6e565b926116416040519485610c5a565b81845260208085019260051b8201019283116101b557602001905b8282106116695750505090565b60208091833561167881610d7c565b81520191019061165c565b9080601f830112156101b557813561169a81610f6e565b926116a86040519485610c5a565b81845260208085019260051b8201019283116101b557602001905b8282106116d05750505090565b6020809183356116df81610de7565b8152019101906116c3565b9080601f830112156101b557610200604051926117078285610c5a565b839181019283116101b557905b8282106117215750505090565b8135815260209182019101611714565b9080601f830112156101b5576102006040519261174e8285610c5a565b839181019283116101b557905b8282106117685750505090565b60208091833561177781610d7c565b81520191019061175b565b9190610640838203126101b5576040519061179c82610c08565b81938035835260208101356117b081610d37565b602084015281605f820112156101b5576040516117cf61020082610c5a565b806102408301918483116101b55760408401905b8382106118115750506112ce9261180685608096946104409460408a01526116ea565b606087015201611731565b60208091833561182081610dce565b8152019101906117e3565b6020818303126101b5578035906001600160401b0382116101b55701610960818303126101b55761185a610c9d565b9181356001600160401b0381116101b55781611877918401610e11565b83526118868160208401610edf565b602084015260808201356001600160401b0381116101b557816118aa91840161142b565b60408401526118bb60a08301610ed4565b60608401526118cd8160c084016114bb565b60808401526118e0816101c08401611502565b60a08401526118f3816102008401611502565b60c08401526119056102408301611549565b60e08401526102608201356001600160401b0381116101b5578161192a918401611558565b6101008401526102808201356001600160401b0381116101b557816119509184016115bf565b6101208401526102a08201356001600160401b0381116101b5578161197691840161161c565b6101408401526102c08201356001600160401b0381116101b5578161199c918401611558565b6101608401526102e08201356001600160401b0381116101b557816119c2918401611683565b6101808401526103008201356001600160401b0381116101b557826119ef83610320936119fb9601611558565b6101a086015201611782565b6101c082015290565b610c8b909291926060810193604080916001600160801b038151168452602081015160208501520151910152565b611a4290611a589281019061182b565b611a4b81611e8a565b84869793929495976129f3565b92611a6284612b50565b611a6b84612d28565b9415611c1157611a7b8184613186565b915b80611c00575b611bef575b505050611a948261021c565b81611ba457611b33611b1c6020604060a08501948551611abd848201516001600160401b031690565b6001600160401b03611aea611ade6003546001600160401b039060401c1690565b6001600160401b031690565b911611611b37575b500151604051611b0981610b558582019485611a04565b519020935101516001600160401b031690565b6001600160401b03165f52600560205260405f2090565b5590565b611b9e906020906001600160401b038151166001600160401b0319600354161760035501517fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6fffffffffffffffff00000000000000006003549260401b16911617600355565b5f611af2565b50611bae8161021c565b60018103611bd85761046c6801000000000000000068ff0000000000000000196004541617600455565b611be18161021c565b6002810361046c5750600290565b611bf89261330f565b5f8080611a88565b50611c0a8561021c565b8415611a83565b611c1a81612f88565b91611a7d565b90611c2a82610f6e565b611c376040519182610c5a565b8281528092611c48601f1991610f6e565b0190602036910137565b634e487b7160e01b5f52603260045260245ffd5b8051821015611c7a5760209160051b010190565b611c52565b903590601e19813603018212156101b557018035906001600160401b0382116101b5576020019181360383136101b557565b903590601e19813603018212156101b557018035906001600160401b0382116101b557602001918160051b360383136101b557565b6002111561022657565b90600182811c92168015611d1e575b6020831014611d0a57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691611cff565b6001545f9291611d3782611cf0565b8082529160018116908115611dab5750600114611d52575050565b60015f9081529293509091907fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b838310611d91575060209250010190565b600181602092949394548385870101520191019190611d80565b9050602093945060ff929192191683830152151560051b010190565b90611dd182611ce6565b52565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff1615611e0c57565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260ff611e6a3360405f20906001600160a01b03165f5260205260405f2090565b541615611e745750565b63e2517d3f60e01b5f523360045260245260445ffd5b5f5f9160408101938451610140815151015195611eaf60406020860151015192613c91565b93611eb988613cc9565b909182611fa6576101c00180515190929015611f89575050611edc905188613df1565b95600194611ee8613d12565b6020845101525b611f6d578785158080611f64575b15611f1d57505050905051606060208201519101526001915b9493929190565b15611f3a5750505160600151611f3291613d96565b600191611f16565b90949214611f49575b50611f16565b611f5b91935060609051015186613f6e565b6001915f611f43565b50818514611efd565b5090506060611f7d939293613d12565b91510152939291905f90565b959650969050611f9e60208351015189613d96565b600195611eef565b50969094611fb2613d12565b602084510152611eef565b60405190611fca82610bb2565b5f6020838281520152565b60405190611fe282610bd2565b5f6040838281528260208201520152565b6040519061200082610c23565b8160405161200d81610c3e565b60608152612019611fbd565b6020820152612026611fbd565b60408201525f60608201525f60808201525f60a08201525f60c08201525f60e08201528152612053611fd5565b6020820152612060611fd5565b60408201525f6060820152612073611fbd565b608082015260a06112ce611fbd565b81601f820112156101b55780519061209982610ccb565b926120a76040519485610c5a565b828452602083830101116101b557815f9260208093018386015e8301015290565b91908260409103126101b5576040516120e081610bb2565b602080829480516120f081610d37565b8452015191610d7883610d37565b91908260409103126101b55760405161211681610bb2565b6020808294805161212681610d7c565b8452015191610d7883610d7c565b5190610c8b82610dce565b5190610c8b82610de7565b5190610c8b82610dfc565b919091610140818403126101b55761216b610c7b565b928151916001600160401b0383116101b5576121b08261219361012094610ebc968501612082565b87526121a281602085016120c8565b6020880152606083016120fe565b60408601526121c160a08201612134565b60608601526121d260c08201612134565b60808601526121e360e0820161213f565b60a08601526121f5610100820161214a565b60c086015201612134565b5190610c8b82610ec3565b91908260609103126101b55760405161222381610bd2565b6040808294805161223381610ec3565b8452602081015160208501520151910152565b6020818303126101b5578051906001600160401b0382116101b55701610180818303126101b5576040519161227a83610c23565b81516001600160401b0381116101b5578261229d83610140936122ed9601612155565b85526122ac836020830161220b565b60208601526122be836080830161220b565b60408601526122cf60e08201612200565b60608601526122e28361010083016120fe565b6080860152016120fe565b60a082015290565b9061046c9061012060e061231485516101408552610140850190610a4f565b9460ff60208083015182815116828801520151166040850152612354604082015160608601906001600160401b0360208092828151168552015116910152565b606081015163ffffffff1660a0850152608081015163ffffffff1660c085015260a081015115158483015261239260c0820151610100860190611dc7565b015163ffffffff16910152565b90606060c08201926001600160401b03815116835263ffffffff60208201511660208401526123f06040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519160c060a0830152825180915260e0820190602060e08260051b8501019401925f905b82821061242457505050505090565b909192939460df1982820301855285519081516004811015610226576124a2826020600195819594829552015190604084820152606061247083516080604085015260c0840190610a4f565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610a4f565b9701950193920190612415565b906060806124c68451608085526080850190610a4f565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b82821061256257505050509060608061255161046c9461253f6020880151602087019015159052565b604087015185820360408701526124af565b9401516001600160401b0316910152565b9091929360208061257f600193609f198a820301865288516124af565b960192019201909291612516565b61046c91606061271a612708845160a0855260206126f48251604060a08901526125d060e0890182516001600160401b0360208092828151168552015116910152565b6102406125ee848301516102c06101208c01526103a08b0190610a4f565b60408301516001600160401b03166101408b015291808901516001600160801b03166101608b0152608081015115156101808b015260a081015180516101a08c0152602090810151805163ffffffff166101c08d015201516101e08b015260c081015115156102008b015260e08101516102208b01526101008101511515828b01526101208101516102608b01526101408101516102808b01526101608101516102a08b01526101808101516102c08b01526101a08101516102e08b01526101c081015115156103008b01526101e08101516103208b015261020081015115156103408b01526102208101516103608b0152015188820360df19016103808a0152610a4f565b910151858203609f190160c087015261239f565b602085015184820360208601526124ef565b92612742604082015160408501906001600160401b0360208092828151168552015116910152565b01519060808184039101526124ef565b905f905b6008821061276357505050565b6020806001928551815201930191019091612756565b905f905b6002821061278a57505050565b602080600192855181520193019101909161277d565b90602080835192838152019201905f5b8181106127bd5750505090565b825115158452602093840193909201916001016127b0565b905f905b601082106127e657505050565b60208060019285518152019301910190916127d9565b905f905b6010821061280d57505050565b6020806001926001600160401b03865116815201930191019091612800565b8051825260ff60208201511660208301526040810151604083015f905b6010821061287b57505050906104406080836128716060610c8b9601516102408601906127d5565b01519101906127fc565b60208060019263ffffffff865116815201930191019091612849565b9061046c906103206101c06129dd6129c96129b56129a161298d61297961290d8b6128fb60206128d18d61096085519181815201906122f5565b92015160208d0190604080916001600160801b038151168452602081015160208501520151910152565b60408d01518b820360808d015261258d565b60608c01516001600160801b031660a08b015261293260808d015160c08c0190612752565b61294360a08d0151898c0190612779565b61295660c08d01516102008c0190612779565b60e08c015161ffff166102408b01526101008c01518a82036102608c015261038d565b6101208b01518982036102808b01526103c6565b6101408a01518882036102a08a01526103f9565b6101608901518782036102c089015261038d565b6101808801518682036102e08801526127a0565b6101a087015185820361030087015261038d565b94015191019061282c565b6040513d5f823e3d90fd5b9291906129fe611ff3565b5015612ad757612a115761046c916140ba565b5f90612a7692612a1f611ff3565b5060405193849283927f9f3665420000000000000000000000000000000000000000000000000000000084527f00000000000000000000000000000000000000000000000000000000000000009060048501614098565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115612ad2575f91612ab6575b5090565b61046c91503d805f833e612aca8183610c5a565b810190612246565b6129e8565b50505f612a7691604051809381927fd27520d9000000000000000000000000000000000000000000000000000000008352602060048401526024830190612897565b15612b22575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b610c8b90612cae8151612b746001600160801b03606085015116633b9aca00900490565b612b8281424282111561411f565b612ca060e0612b918342613718565b93612c95600454612bbf612bac8263ffffffff9060501c1690565b9663ffffffff8816988942911115614156565b612bf28351805160208201207f00000000000000000000000000000000000000000000000000000000000000001461418d565b612c3a6020840151612c05815160ff1690565b6002549160ff831660ff811660ff8416149384612cf1575b612c349060209060081c60ff165b93015160ff1690565b9361428c565b612c60612c4e606085015163ffffffff1690565b63ffffffff83811690821681146142e8565b612c81612c74608085015163ffffffff1690565b9160201c63ffffffff1690565b63ffffffff811663ffffffff831614614329565b015163ffffffff1690565b9163ffffffff83161461436a565b612ce9612ce46020608081850151604051612cd081610b558682019485611a04565b51902094015101516001600160401b031690565b6143ab565b808214612b19565b9350612c346020612c2b612d088286015160ff1690565b60ff612d1b60088a901c82165b60ff1690565b9116149692505050612c1d565b612d43611b1c602060a084015101516001600160401b031690565b5480612d4f5750505f90565b60408201908151604051612d6b81610b55602082019485611a04565b5190201491821592612d89575b505015612d8457600190565b600290565b6001600160801b03919250612dbf612db06020612dcb9301516001600160801b0390511690565b9351516001600160801b031690565b6001600160801b031690565b911610155f80612d78565b15612ddd57565b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52601160045260245ffd5b906001600160401b03809116911601906001600160401b038211612e3957565b612e05565b15612e465750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15612e805750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15612eba5750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b90600382029180830460031490151715612e3957565b908160011b9180830460021490151715612e3957565b90602c820291808304602c1490151715612e3957565b908160051b9180830460201490151715612e3957565b15612f4d575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b905f9160206040820151015151805193610100830192612fde845151612fb66104df60e085015161ffff1690565b8091149081613176575b81613166575b81613156575b81613146575b81613136575b50612dd6565b5f5b86811061311957505f92839283805b8751518710156130d7579089929161301c6130186130128a6101808a0151611c66565b51151590565b1590565b6130cc576001916130a99161303f6130358b8d51611c66565b5163ffffffff1690565b80996130548263ffffffff8116998a10612e3e565b6130b4575b6130839150602061306a888a611c66565b51015161307c8c6101208c0151611c66565b5114612eb2565b6130a36040613094859a9789611c66565b5101516001600160401b031690565b90612e19565b965b01959091612fef565b63ffffffff6130c592168711612e78565b5f88613059565b9250956001906130ab565b5091509650610c8b94508691935061311492506130fc6001600160401b038216612eec565b61310e6001600160401b038416612f02565b10612f44565b614521565b9161312f6001916130a360406130948789611c66565b9201612fe0565b90506101a083015151145f612fd8565b6101808401515181149150612fd2565b6101608401515181149150612fcc565b6101408401515181149150612fc6565b6101208401515181149150612fc0565b906101008101906131a5825151612fb66104df60e085015161ffff1690565b6131ae8361364d565b926001600160a01b038416156132fd5793926131d9600854916104b86104b38461ffff9060e01c1690565b946131f76131e78783613806565b9260a01c6001600160401b031690565b955f945f5f9860095497600a549260205f98019b5b8551518910156132d65791898b94928e989796946132356130186130128e610180870151611c66565b6132c8576132476130358d8a51611c66565b80926132aa575b8c91508b90876001998c839e516132669061ffff1690565b9061327095613948565b91909361012001519061328291611c66565b51149061328e91612eb2565b61329791612e19565b976001905b01979193949592909261320c565b63ffffffff6132c1921663ffffffff821611612e78565b5f8161324e565b98509450509760019061329c565b5050935098505050610c8b945061311492508691506130fc6001600160401b038216612eec565b633a517eed60e21b5f5260045260245ffd5b91906001600160a01b036133228461364d565b16151580613513575b61350e576040015160200151518051938415612ddd5760b485116134f55761335285611c20565b925f5b83518110156133995780613388602061337060019488611c66565b5101516133826040613094858a611c66565b906153cb565b6133928288611c66565b5201613355565b5092613435906134145f979693966133cf896133bc6133b784612f02565b61370a565b94836133c787611c20565b9c8d926153f5565b5061340e6133fe6133f96133ea6133e585612f18565b613668565b6133f387612f2e565b906136d8565b613725565b9761340889615492565b886154da565b86615568565b61342f61342861342385615bf2565b615902565b86602e0152565b84615578565b5f9460305b83518710156134ad576134a56001916134a06134568a88611c66565b5161346863ffffffff8c16848b6154b6565b61348961347484613676565b60408301516001600160401b0316908b615520565b602061349484613684565b91015190890160200152565b613692565b96019561343a565b95509150925f945b82518610156134e8576134e06001916134db6134d18987611c66565b5187830160200152565b6136a0565b9501946134b5565b50935050610c8b916146ac565b63156f758160e31b5f52600485905260b460245260445ffd5b505050565b5061ffff60085460e01c16151561332b565b805f525f60205260ff61354c8360405f20906001600160a01b03165f5260205260405f2090565b54166135b757805f525f6020526135778260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916600117905533916001600160a01b0316907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260ff6135e48360405f20906001600160a01b03165f5260205260405f2090565b5416156135b757805f525f6020526136108260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916905533916001600160a01b0316907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b60065403613664576001600160a01b036008541690565b5f90565b6030019081603011612e3957565b9060048201809211612e3957565b90600c8201809211612e3957565b90602c8201809211612e3957565b9060208201809211612e3957565b6001019081600111612e3957565b6024019081602411612e3957565b9060018201809211612e3957565b91908201809211612e3957565b61ffff16602c810290808204602c1490151715612e395760300180603011612e395790565b5f19810191908211612e3957565b91908203918211612e3957565b9061372f82610ccb565b61373c6040519182610c5a565b8281528092611c48601f1991610ccb565b9190823b600181111561379b575f198101908111612e3957811161377f57600161377682613725565b9360208501903c565b6001600160a01b0383633610565160e21b5f521660045260245ffd5b6001600160a01b0384633610565160e21b5f521660045260245ffd5b604051906137c482610bed565b5f6060838281528260208201528260408201520152565b60011b906201fffe61fffe831692168203612e3957565b61ffff5f199116019061ffff8211612e3957565b9190916138116137b7565b5060308351106138ca576356414c34602084015160e01c036138ca57602483015160c01c92602c81015160f01c93602e82015190604e83015160f01c91613868613859610cad565b6001600160401b039093168352565b61387a6020830197889061ffff169052565b60408201526138916060820192839061ffff169052565b9461389e815161ffff1690565b9161ffff831692831593841561393d575b50831561390c575b5082156138dc575b505090506138ca5750565b634724a0fd60e01b5f5260045260245ffd5b61301892506138fe6138f56139049551935161ffff1690565b915161ffff1690565b916148b0565b805f806138bf565b90925061ffff6139326104df61392d613927875161ffff1690565b946137db565b6137f2565b91161415915f6138b7565b60b41093505f6138af565b939092959461ffff63ffffffff82169316831015613a155761ffff16916139716133e582612f18565b94816139896020888801015160e01c63ffffffff1690565b036138ca57506001901b9081166139ea576139b26139a685613676565b84016020015160c01c90565b95166139cd57506139c561046c92613684565b016020015190565b90506139e6915061ffff165f52600b60205260405f2090565b5490565b613a10613a038361ffff165f52600c60205260405f2090565b546001600160401b031690565b6139b2565b6303e07d4560e61b5f52600485905263ffffffff1660245260445ffd5b9995969998919790949298613a468b611ce6565b8a15613a87576105cd8b613a5981611ce6565b7f112d89cc000000000000000000000000000000000000000000000000000000005f5260ff16600452602490565b88999a50602090898095969798999a15159081613c84575b613aa8916148ff565b0196613ac7613ab68961493d565b613ac0368c610edf565b9087615619565b5f905f5b88868210613be3575b505050613ae19350614ae4565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016803b156101b557838792613b505f956040519b8c96879586957f87d3a9b100000000000000000000000000000000000000000000000000000000875260048701614eb1565b03915afa938415612ad25761046c95613b7c95613bc9575b5060018111613b90575b5050503690610edf565b6001600160801b03633b9aca009151160490565b6001600160801b03613bb9613ba7613bc19561493d565b93613bb187614f6b565b933691614f75565b91169161575d565b5f8080613b72565b80613bd75f613bdd93610c5a565b806105d0565b5f613b68565b613018613c0b613c04613bfe858b613c1c9699979899614947565b80611cb1565b3691614969565b613c16368989614969565b906156eb565b613c7a575090613c436107eb613c39613c5994613ae1988c614947565b6020810190611c7f565b90815181518082149182613c64575b50506149d4565b600189935f88613ad4565b9091506020840120906020830120145f80613c52565b9190600101613acb565b61ffff8111159150613a9f565b60016001600160401b03602060408281865151015116940151015116016001600160401b038111612e39576001600160401b03161490565b613cda6001600160a01b039161364d565b1615613ce857600754600191565b5f905f90565b60405190613cfb82610bed565b5f6060838181528260208201528260408201520152565b60405190613d1f82610bed565b606082525f6020830152613d31613cee565b60408301525f60608301528160209060405191613d4f602084610c5a565b5f808452805b818110613d625750505052565b8290613d6c613cee565b82828801015201613d55565b15613d81575050565b63769a20cb60e01b5f5260045260245260445ffd5b908051518015613dd65760b48111613dbf575090613db6610c8b92615040565b90808214613d78565b63156f758160e31b5f5260045260b460245260445ffd5b82633a517eed60e21b5f5260045260245ffd5b156138ca5750565b91908051613dfe8161364d565b906001600160a01b038216156132fd575092613f546001600160401b03613ed6613ef9613eeb610c8b96613e34613f169a6150d8565b613e3f818351613806565b906020820192613e76613e54855161ffff1690565b61ffff613e6b6104df60085461ffff9060e01c1690565b911614825190613de9565b613e9d613e8a612d15602084015160ff1690565b8015159081613f62575b50825190613de9565b613ef4600954958695600a54978185613ec58b84613ebe82975161ffff1690565b8b8561510b565b9d8f8f908495939551911115613de9565b88613ee58351925161ffff1690565b9161525b565b8b808214613d78565b6152b1565b613f11613f0b613423889b949b615bf2565b99600955565b600a55565b167fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6008549260a01b16911617600855565b613f5d84600755565b600655565b6010915011155f613e94565b90613f788261364d565b6001600160a01b0381161561408457613f9f906104b86104b360085461ffff9060e01c1690565b613fa98184613806565b91519160208351910190613fc26104df835161ffff1690565b0361406957600954600a54915f5b855181101561406057613ff9613fe8835161ffff1690565b858563ffffffff851692898c613948565b6020614005848a611c66565b510151149081159161403a575b5061401f57600101613fd0565b6105cd8763769a20cb60e01b5f52906044916004525f602452565b90506001600160401b03806140546040613094868c611c66565b9216911614155f614012565b50505050505050565b6105cd8463769a20cb60e01b5f52906044916004525f602452565b633a517eed60e21b5f52600483905260245ffd5b6140b060409295949395606083526060830190612897565b9460208201520152565b612a76915f916140c8611ff3565b5060405193849283927f22b541630000000000000000000000000000000000000000000000000000000084527f00000000000000000000000000000000000000000000000000000000000000009060048501614098565b15614128575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b1561415f575050565b7f12e74add000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b156141955750565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f6001546141d181611cf0565b908160448501526001811690815f146142685750600114614208575b506142049192600319848303016024850152610a4f565b0390fd5b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b81831061424e57509192915081016064016142046141ed565b805483870160640152859450602090920191600101614235565b60ff191660648086019190915291151560051b8401909101915061420490506141ed565b93929190931561429c5750505050565b9160ff8092816084969581604051977fd382033a000000000000000000000000000000000000000000000000000000008952166004880152166024860152166044840152166064820152fd5b156142f1575050565b9063ffffffff80927f73b1bde3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b15614332575050565b9063ffffffff80927f1eb5c1eb000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b15614373575050565b9063ffffffff80927f87f10409000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6001600160401b03165f52600560205260405f205480156143c95790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b908160209103126101b5575161046c81610de7565b979260a09661446761046c9b97614453614494986101608e61444c60c09f99986144416144859c61ffff6144769c1685526020850190612752565b610120830190612779565b0190612779565b6102406101a08d01526102408c01906103c6565b908a82036101c08c01526103f9565b908882036101e08a015261038d565b908682036102008801526127a0565b936102208186039101526001600160401b0381511684526001600160401b0360208201511660208501526040810151604085015263ffffffff6060820151166060850152608081015160808501520151918160a08201520190610a4f565b156144f957565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b60206040820151518181015161453e81516001600160401b031690565b916145c160406145606145578786015163ffffffff1690565b63ffffffff1690565b930151858151910151906145af878061457e8563ffffffff90511690565b94015195510151956145a0614591610cbc565b6001600160401b039099168952565b6001600160401b031687890152565b604086015263ffffffff166060850152565b608083015260a082015260e083015161ffff1661463760808501519260a08601519560c08101519061012081015161014082015190610180610160840151930151936040519a8b998a997f4cc22bb7000000000000000000000000000000000000000000000000000000008b5260048b01614406565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af18015612ad257610c8b915f9161467d575b506144f2565b61469f915060203d6020116146a5575b6146978183610c5a565b8101906143f1565b5f614677565b503d61468d565b6146b68282613806565b91604051905f60208301526146e0826146d26021820184615588565b03601f198101845283610c5a565b6160008251116148845750610b5561473a614728614700845161ffff1690565b60f01b7fffff0000000000000000000000000000000000000000000000000000000000001690565b9260405192839160208301958661559a565b51905ff0906001600160a01b0382161561485c5761484a92614793602092613f5d6147fa956001600160a01b03167fffffffffffffffffffffffff00000000000000000000000000000000000000006008541617600855565b6147a06040820151600755565b6147f16147b482516001600160401b031690565b7fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6008549260a01b16911617600855565b015161ffff1690565b7fffff0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff7dffff000000000000000000000000000000000000000000000000000000006008549260e01b16911617600855565b6148535f600955565b610c8b5f600a55565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b906148ba906136e5565b908181149283156148cc575b50505090565b90919250621fffe061ffff82169160051b169080820460201490151715612e39578201809211612e3957145f80806148c6565b156149075750565b7fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b3561046c81610d7c565b9190811015611c7a5760051b81013590603e19813603018212156101b5570190565b92919061497581610f6e565b936149836040519586610c5a565b602085838152019160051b8101918383116101b55781905b8382106149a9575050505050565b81356001600160401b0381116101b5576020916149c98784938701610d1c565b81520191019061499b565b919091156149e0575050565b90614204614a22926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610a4f565b83810360031901602485015290610a4f565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156101b55701602081359101916001600160401b0382116101b55781360383136101b557565b90602083828152019260208260051b82010193835f925b848410614aac5750505050505090565b909192939495602080614ad4600193601f19868203018852614ace8b88614a54565b90614a34565b9801940194019294939190614a9c565b91909115614af0575050565b6142046040519283927ffef760c7000000000000000000000000000000000000000000000000000000008452602060048501526024840191614a85565b9035601e19823603018112156101b55701602081359101916001600160401b0382116101b5578160051b360383136101b557565b9035607e19823603018112156101b5570190565b359060038210156101b557565b9035605e19823603018112156101b5570190565b90602083828152019260208260051b82010193835f925b848410614bbd5750505050505090565b909192939495602080614c2f600193601f19868203018852614bdf8b88614b82565b90614be982614b75565b614bf28161021c565b8152614c21614c16614c0686850185614a54565b6060888601526060850191614a34565b926040810190614a54565b916040818503910152614a34565b9801940194019294939190614bad565b61046c91614d09614cfe614c82614c67614c598680614a54565b608087526080870191614a34565b614c746020870187614a54565b908683036020880152614a34565b6080614cee614c946040880188614b61565b8684036040880152614ca581614b75565b614cae8161021c565b8452614cbc60208201614b75565b614cc58161021c565b6020850152614cd660408201614b75565b614cdf8161021c565b60408501526060810190614a54565b9190928160608201520191614a34565b926060810190614b2d565b916060818503910152614b96565b90602083828152019260208260051b82010193835f925b848410614d3e5750505050505090565b909192939495601f198282030184528635601e19843603018112156101b557830190614d6e602082019280614b2d565b8091936020845252604082019060408160051b8401019380935f915b838310614dad575050505050506020806001929801940194019294939190614d2e565b909192939495603f19838203018652614dc68783614b82565b8035614dd181610dfc565b614dda81611ce6565b8252614dfd614dec6020830183614b61565b606060208501526060840190614c3f565b906040810135609e19823603018112156101b5576001936020938493614ea39301916040818303910152614e95614e75614e48614e3a8580614a54565b60a0865260a0860191614a34565b86850135614e5581610de7565b151587850152614e686040860186614b61565b8482036040860152614c3f565b926060810135614e8481610de7565b151560608401526080810190614b61565b906080818403910152614c3f565b980196019493019190614d8a565b9092809492969593606083019083526060602084015252608081019560808560051b83010194815f90603e1981360301995b838310614f0257505050505061046c9495506040818503910152614d17565b9091929397607f1986820301825288358b8112156101b5576020614f5c6001938683940190614f4f614f45614f378480614b2d565b604085526040850191614a85565b9285810190614a54565b9185818503910152614a34565b9a019201930191909392614ee3565b3561046c81610ec3565b92919092614f8284610f6e565b93614f906040519586610c5a565b602085828152019060051b8201918383116101b55780915b838310614fb6575050505050565b82356001600160401b0381116101b55782016040818703126101b55760405191614fdf83610bb2565b81356001600160401b0381116101b557820187601f820112156101b5578781602061500c93359101614969565b83526020820135926001600160401b0384116101b55761503188602095869501610d1c565b83820152815201920191614fa8565b8051519081156135b75761505382611c20565b915f5b825180518210156150ca57906150b9613423602061507684600196611c66565b5101516150b46001600160401b036040615091878b51611c66565b51015116604051926150a284610bb2565b83526001600160401b03166020830152565b61582e565b6150c38287611c66565b5201615056565b505090505f61046c92615952565b90813b600181111561377f575f198101908111612e3957600161377682613725565b906010811015611c7a5760051b0190565b919490929360208301946151266104e6612d15885160ff1690565b96615140611ade6008546001600160401b039060a01c1690565b955f945f5b615153612d158b5160ff1690565b81101561524f5761516b6130358260408b01516150fa565b9663ffffffff8816906151848961ffff88168410612e3e565b82615235575b505082878787878c519461519d95613948565b908260808b0151906151ae916150fa565b516001600160401b03169a8360608c0151906151c9916150fa565b51926001600160401b038d16926001600160401b031690848483149182159261522a575b50508c516151fa91613de9565b61520391613718565b9061520d916136d8565b99615217916153cb565b615221828d611c66565b52600101615145565b14159050845f6151ed565b6152489163ffffffff8b51921610613de9565b5f8061518a565b50969750505050505050565b60205f95939461046c9861527d9861ffff8998899760ff976040519d8e610c3e565b8d52868d015216968760408c01528360608c015260808b01528060a08b01528160c08b01521760e0890152015116946159b4565b9294935f955b6152c8612d15602087015160ff1690565b8710156153c2576001600160401b03916001916152ef6104df6130358b60408b01516150fa565b91600161ffff84161b928361531861530b8d60808d01516150fa565b516001600160401b031690565b926153278d60608d01516150fa565b519361533b8c518c8c61ffff881692615b6e565b99166001600160401b038216145f146153865750901916955b8203615368575050901916965b01956152b7565b61537e9061ffff165f52600b60205260405f2090565b551796615361565b6153bb906153a08561ffff165f52600c60205260405f2090565b906001600160401b03166001600160401b0319825416179055565b1795615354565b95509392505050565b613423906001600160401b0361046c93604051926153e884610bb2565b835216602082015261582e565b9190949394808203828111612e3957600181146154735761541590615c2a565b9283820192838311612e39576001880192838911612e395783878661543a93866153f5565b6154496133b760019297612f02565b890101809311612e395761545e9386926153f5565b61546b90611dd192615d64565b938492611c66565b50615485915061548e92959495611c66565b51928392611c66565b5290565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b602d908260081c602c8201530153565b604f908260081c604e8201530153565b805191908290602001825e015f815290565b600a907fffff00000000000000000000000000000000000000000000000000000000000061046c94937f610000000000000000000000000000000000000000000000000000000000000083521660018201527f80600a3d393df30000000000000000000000000000000000000000000000000060038201520190615588565b9061566890612ce960405160208101906156508288604080916001600160801b038151168452602081015160208501520151910152565b6060815261565f608082610c5a565b519020916143ab565b60208201518082036156bd5750506001600160801b03633b9aca009151160461569581424282111561411f565b4203428111612e39576001600160801b03610c8b911663ffffffff6004541690818110615c46565b7f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b9081518151036135b7575f5b82518110156157555761570a8184611c66565b51516157168284611c66565b51510361574e576157278184611c66565b51602081519101206157398284611c66565b51602081519101200361574e576001016156f7565b5050505f90565b505050600190565b929190925f5b8451811015615827576157768186611c66565b51908360405160208101906001600160401b038616825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801915f5b8181106157f0575050505081610b5560019760206157e6940151605f19848303016080850152610a4f565b5190205d01615763565b91939496509194969760208061581260019360bf198b82030188528951610a4f565b970194019101918a96949392989795986157bb565b5050509050565b602460208201906001600160401b03825116806158c7575b5061585361588191613725565b9261587861587261586c61586687615c8a565b87615c96565b86615cad565b85615cc4565b90519084615cf1565b90615896611ade82516001600160401b031690565b61589f57505090565b6158c0611ade6158b2612ab29486615cda565b92516001600160401b031690565b9083615d09565b600191505b60808110156158f457506158536158ed6158e8615881936136ae565b6136bc565b9150615846565b60019060071c9101906158cc565b805160018101809111612e395761591890613725565b805115611c7a57615942816159356020945f868196015382615d42565b5060405191828092615588565b039060025afa15612ad2575f5190565b91818103818111612e3957600181146159975761596e90615c2a565b820191828111612e3957826159839185615952565b9161598e9293615952565b61046c91615d64565b50506159a291611c66565b5190565b5f198114612e395760010190565b9392909491828414808091615b56575b615b2e5761ffff8311612ddd576159db8783613718565b9060018214615a8457506159ee90615c2a565b93615a1b6159fc86896136d8565b956133f36133b7615a15615a0f886136ca565b976136ca565b92612f02565b92815b85811087828a83615a5d575b50505015615a4057615a3b906159a6565b615a1e565b9495969791859188615a52948b6159b4565b9461598e95966159b4565b63ffffffff929350615a7a916040606061303593015101516150fa565b161087828a615a2a565b925050509493929415615ad057505082615acb9161046c939451916020810151615ab3604083015161ffff1690565b9063ffffffff60c060a0850151940151941694613948565b6153cb565b615adb8293926136ca565b1490811591615b0c575b50615af85760806159a292930151611c66565b8251634724a0fd60e01b5f5260045260245ffd5b9050615b26614557613035846040606089015101516150fa565b14155f615ae5565b50509291505061046c9250805190615b506040602083015192015161ffff1690565b91615dd8565b50615b69613018838960e08a0151615db6565b6159c4565b9091602061ffff919594950151169363ffffffff811694851015615bd45750615b996133e585612f18565b936020858401015160e01c036138ca5750602090615bce615bc8615bbc86613676565b83016020015160c01c90565b94613684565b01015190565b6303e07d4560e61b5f5260049190915263ffffffff1660245260445ffd5b60405190606090615c038284610c5a565b60228352612ab291600a906020850190601f19013682375360206021840153600283615cf1565b9060015b8060011b9083821015615c415750615c2e565b925050565b15615c4f575050565b906001600160801b0380927fdb020436000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6020600a910153600190565b60208260229201015360018101809111612e395790565b602082600a9201015360018101809111612e395790565b602082819201015360018101809111612e395790565b60208260109201015360018101809111612e395790565b8160209193929301015260208101809111612e395790565b9092919083016020015b6080821015615d2757906001929391530190565b600180916080607f85161781530193019060071c9092615d13565b908051918215615755576021602084930191015e60010180600111612e395790565b604051615d72608082610c5a565b60418152615d806041610ccb565b602082019190601f1901368337805115611c7a576020935f93600161594294536021830152604182015260405191828092615588565b91818103908111612e39576001901b905f198201918211612e39571b16151590565b909161ffff1692602c840293808504602c1490151715612e3957836030019384603011612e39578160051b9180830460201490151715612e3957019260308401809111612e3957615e28906136a0565b8251106138ca575001605001519056fea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
}

// ContractGroth16ICS07TendermintABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractGroth16ICS07TendermintMetaData.ABI instead.
var ContractGroth16ICS07TendermintABI = ContractGroth16ICS07TendermintMetaData.ABI

// ContractGroth16ICS07TendermintBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractGroth16ICS07TendermintMetaData.Bin instead.
var ContractGroth16ICS07TendermintBin = ContractGroth16ICS07TendermintMetaData.Bin

// DeployContractGroth16ICS07Tendermint deploys a new Ethereum contract, binding an instance of ContractGroth16ICS07Tendermint to it.
func DeployContractGroth16ICS07Tendermint(auth *bind.TransactOpts, backend bind.ContractBackend, verifier common.Address, membership_ common.Address, misbehaviour_ common.Address, updateClient_ common.Address, _clientState []byte, _consensusState [32]byte, roleManager common.Address) (common.Address, *types.Transaction, *ContractGroth16ICS07Tendermint, error) {
	parsed, err := ContractGroth16ICS07TendermintMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractGroth16ICS07TendermintBin), backend, verifier, membership_, misbehaviour_, updateClient_, _clientState, _consensusState, roleManager)
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

// GetCachedValidatorSet is a free data retrieval call binding the contract method 0x536c2ad3.
//
// Solidity: function getCachedValidatorSet(bytes32 validatorsHash) view returns(uint32[] indices, bytes32[] pubkeys, uint64[] votingPowers)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) GetCachedValidatorSet(opts *bind.CallOpts, validatorsHash [32]byte) (struct {
	Indices      []uint32
	Pubkeys      [][32]byte
	VotingPowers []uint64
}, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "getCachedValidatorSet", validatorsHash)

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

// GetCachedValidatorSet is a free data retrieval call binding the contract method 0x536c2ad3.
//
// Solidity: function getCachedValidatorSet(bytes32 validatorsHash) view returns(uint32[] indices, bytes32[] pubkeys, uint64[] votingPowers)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) GetCachedValidatorSet(validatorsHash [32]byte) (struct {
	Indices      []uint32
	Pubkeys      [][32]byte
	VotingPowers []uint64
}, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GetCachedValidatorSet(&_ContractGroth16ICS07Tendermint.CallOpts, validatorsHash)
}

// GetCachedValidatorSet is a free data retrieval call binding the contract method 0x536c2ad3.
//
// Solidity: function getCachedValidatorSet(bytes32 validatorsHash) view returns(uint32[] indices, bytes32[] pubkeys, uint64[] votingPowers)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) GetCachedValidatorSet(validatorsHash [32]byte) (struct {
	Indices      []uint32
	Pubkeys      [][32]byte
	VotingPowers []uint64
}, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GetCachedValidatorSet(&_ContractGroth16ICS07Tendermint.CallOpts, validatorsHash)
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

// Misbehaviour is a free data retrieval call binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes ) view returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) Misbehaviour(opts *bind.CallOpts, arg0 []byte) error {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "misbehaviour", arg0)

	if err != nil {
		return err
	}

	return err

}

// Misbehaviour is a free data retrieval call binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes ) view returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) Misbehaviour(arg0 []byte) error {
	return _ContractGroth16ICS07Tendermint.Contract.Misbehaviour(&_ContractGroth16ICS07Tendermint.CallOpts, arg0)
}

// Misbehaviour is a free data retrieval call binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes ) view returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) Misbehaviour(arg0 []byte) error {
	return _ContractGroth16ICS07Tendermint.Contract.Misbehaviour(&_ContractGroth16ICS07Tendermint.CallOpts, arg0)
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
