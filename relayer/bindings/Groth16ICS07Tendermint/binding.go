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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membership_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviour_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"updateClient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_clientState\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCachedValidatorSet\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"updateClientMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CannotHandleMisbehavior\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientNotFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ClientStateMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ClockDriftMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustedVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidMembershipProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NonMonotonicHeightUpdate\",\"inputs\":[{\"name\":\"latestHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"updateHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofHeightMismatch\",\"inputs\":[{\"name\":\"expectedRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"expectedRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofSignerCommitSigMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustThresholdMismatch\",\"inputs\":[{\"name\":\"expectedNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expectedDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnbondingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"UnknownZkAlgorithm\",\"inputs\":[{\"name\":\"algorithm\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorCountExceedsLimit\",\"inputs\":[{\"name\":\"validatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxValidatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"VerificationKeyMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x6101606040523461010157617c4f803803809161001b82610119565b6101603960e081610160019112610101576100346101a4565b61003f6101806101bb565b61004a6101a06101bb565b6100556101c06101bb565b6101e0516001600160401b038111610101578561017f82011215610101576100a2958161018061008b93610160015191016101ea565b91610200519361009c6102206101bb565b95610800565b604051616ba19081610fee82396080518161602a015260a051816140d7015260c05181610ad1015260e0518161306c0152610100518181816109c50152610a0e01526101205181614ca5015261014051818181613037015261497d0152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b610160601f91909101601f19168101906001600160401b0382119082101761014057604052565b610105565b604081019081106001600160401b0382111761014057604052565b601f909101601f19168101906001600160401b0382119082101761014057604052565b6040519061019361010083610160565b565b60405190610193604083610160565b61016051906001600160a01b038216820361010157565b51906001600160a01b038216820361010157565b6001600160401b03811161014057601f01601f191660200190565b9291926101f6826101cf565b916102046040519384610160565b829481845281830111610101578281602093845f96015e010152565b9080601f8301121561010157815161023a926020016101ea565b90565b519060ff8216820361010157565b91908260409103126101015760405161026381610145565b602061027c8183956102748161023d565b85520161023d565b910152565b51906001600160401b038216820361010157565b9190826040910312610101576040516102ad81610145565b602061027c8183956102be81610281565b855201610281565b519063ffffffff8216820361010157565b5190811515820361010157565b5190600282101561010157565b602081830312610101578051906001600160401b03821161010157016101408183031261010157610320610183565b8151909290916001600160401b038311610101576103678261034a610120946103b7968501610220565b8652610359816020850161024b565b602087015260608301610295565b604085015261037860a082016102c6565b606085015261038960c082016102c6565b608085015261039a60e082016102d7565b60a08501526103ac61010082016102e4565b60c0850152016102c6565b60e082015290565b90600182811c921680156103ed575b60208310146103d957565b634e487b7160e01b5f52602260045260245ffd5b91607f16916103ce565b601f821161040457505050565b5f5260205f20906020601f840160051c8301931061043c575b601f0160051c01905b818110610431575050565b5f8155600101610426565b909150819061041d565b6002111561045057565b634e487b7160e01b5f52602160045260245ffd5b9060028110156104505769ff00000000000000000082549160481b169069ff0000000000000000001916179055565b80518051906001600160401b038211610140576104bc826104b56001546103bf565b60016103f7565b602090601f831160011461065d57926104f5836106229460e094610193975f92610652575b50508160011b915f199060031b1c19161790565b6001555b6020818101518051600280549284015161ff0060089190911b1660ff90921661ffff199093169290921717905560408083015180516003805492909401516fffffffffffffffff0000000000000000931b929092166001600160401b039092166001600160801b031990911617179055606081015161058f9063ffffffff1660049063ffffffff1663ffffffff19825416179055565b6105c76105a3608083015163ffffffff1690565b60049067ffffffff0000000082549160201b169067ffffffff000000001916179055565b6105ff6105d760a0830151151590565b60049068ff0000000000000000825491151560401b169068ff00000000000000001916179055565b61061760c082015161061081610446565b6004610464565b015163ffffffff1690565b6004906dffffffff0000000000000000000082549160501b16906dffffffff000000000000000000001916179055565b015190505f806104e1565b60015f52601f19831691905f516020617c0f5f395f51905f52925f5b8181106106be57509360e09361019396936001938361062298106106a6575b505050811b016001556104f9565b01515f1960f88460031b161c191690555f8080610698565b92936020600181928786015181550195019301610679565b604051905f82600154916106e9836103bf565b8083529260018116908115610759575060011461070d575b61019392500383610160565b5060015f90815290915f516020617c0f5f395f51905f525b81831061073d57505090602061019392820101610701565b6020919350806001915483858901015201910190918492610725565b6020925061019394915060ff191682840152151560051b820101610701565b15610781575050565b6355bace6f60e01b5f9081526001600160401b039182166004529116602452604490fd5b634e487b7160e01b5f52601160045260245ffd5b9063ffffffff8091169116019063ffffffff82116107d357565b6107a5565b156107e1575050565b9063ffffffff80926333fae18560e21b5f52166004521660245260445ffd5b919361084161084691969294967f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a6123361010052602080825183010191016102f1565b610493565b61084e6106d6565b60208151910120610120526108696108646106d6565b610987565b610140526108dc6108c361089660206108886108836106d6565b610a61565b01516001600160401b031690565b600354906108b4906001600160401b03808416919081168214610778565b60401c6001600160401b031690565b6001600160401b03165f90815260056020526040902090565b556001600160a01b0390811660805290811660a05290811660c0521660e05260045461093f9063ffffffff811661092d610920605084901c63ffffffff16836107b9565b9260201c63ffffffff1690565b9163ffffffff808416911611156107d8565b6001600160a01b0381166109665750610956610d45565b5061096361010051610dc7565b50565b8061097361096392610c47565b5061097d81610cbd565b5061010051610e20565b61099090610ecb565b8051600181018091116107d3576109a690610e99565b8051156109ee57602091816109c2845f94019284845382610f98565b50604051918291518091835e8101838152039060025afa156109e3575f5190565b6040513d5f823e3d90fd5b634e487b7160e01b5f52603260045260245ffd5b60405190610a0f82610145565b5f602083606081520152565b80156107d3575f190190565b5f198101919082116107d357565b919082039182116107d357565b9081518110156109ee570160200190565b90600182018092116107d357565b610a69610a02565b5080518015908115610c3b575b50610c2c575f19908051805b610bdc575b505f198214610bbe57600360fc1b6001600160f81b0319610ac1610ab3610aad86610a53565b85610a42565b516001600160f81b03191690565b161480610bc8575b610bbe575f90610ad883610a53565b915b8151831015610b6f57610af9610af3610ab38585610a42565b60f81c90565b60ff811660308110908115610b64575b50610b5757600a82026001600160401b03908116602f1990920160ff1691909101811691168110610b3f57600190920191610ada565b50915050610b4b610195565b9081525f602082015290565b5050915050610b4b610195565b60399150115f610b09565b9150916001811190811591610bb2575b50610ba35761023a90610b90610195565b9283526001600160401b03166020830152565b63384c687760e21b5f5260045ffd5b602b915010155f610b7f565b9050610b4b610195565b506002610bd6838351610a35565b11610ac9565b602d60f81b610c06610bf9610ab3610bf385610a27565b86610a42565b6001600160f81b03191690565b14610c1a57610c1490610a1b565b80610a82565b610c25919250610a27565b905f610a87565b6329120bff60e21b5f5260045ffd5b6040915010155f610a76565b6001600160a01b0381165f9081525f516020617c2f5f395f51905f52602052604090205460ff16610cb8576001600160a01b03165f8181525f516020617c2f5f395f51905f5260205260408120805460ff191660011790553391905f516020617b8f5f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f516020617baf5f395f51905f52602052604090205460ff16610cb8576001600160a01b0381165f9081525f516020617baf5f395f51905f5260205260409020805460ff1916600117905533906001600160a01b03165f516020617bef5f395f51905f525f516020617b8f5f395f51905f525f80a4600190565b5f80525f516020617baf5f395f51905f526020525f516020617bcf5f395f51905f525460ff16610dc3575f8080525f516020617baf5f395f51905f526020525f516020617bcf5f395f51905f52805460ff1916600117905533905f516020617bef5f395f51905f525f516020617b8f5f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f14610cb8575f818152602081815260408083208380529091528120805460ff1916600117905533915f516020617b8f5f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16610e93575f818152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905533916001600160a01b0316905f516020617b8f5f395f51905f525f80a4600190565b50505f90565b90610ea3826101cf565b610eb06040519182610160565b8281528092610ec1601f19916101cf565b0190602036910137565b805115610f375780516001905b6080811015610f2957508060010190816001116107d357600190835101018091116107d357610f09610f2591610e99565b91610f1f610f1684610f53565b82519085610f5f565b83610fc2565b5090565b60019060071c910190610ed8565b50604051610f46602082610160565b5f81525f36602083013790565b6020600a910153600190565b9092919083016020015b6080821015610f7d57906001929391530190565b600180916080607f85161781530193019060071c9092610f69565b908051918215610fba576021602084930191015e600101806001116107d35790565b505050600190565b908092918251928315610fe657839260208092019201015e81018091116107d35790565b505050509056fe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a7146101145780630bece3561461010f578063248a9ca31461010a5780632f2ff15d1461010557806336568abe14610100578063536c2ad3146100fb5780636a28f000146100f65780638a8e4c5d146100f157806391d14854146100ec578063974a74c4146100e7578063a217fddf146100e2578063a6f031bb146100dd578063d547741f146100d8578063db3e1fa4146100d3578063ddba6537146100ce5763ef913a4b146100c9575f80fd5b610c2b565b6109e8565b6109ae565b61097f565b610872565b610858565b6106f8565b6106b7565b61067f565b6105da565b61046f565b61033e565b610308565b6102b0565b61022b565b346101b55760203660031901126101b5576004357fffffffff0000000000000000000000000000000000000000000000000000000081168091036101b557807f7965db0b000000000000000000000000000000000000000000000000000000006020921490811561018b575b506040519015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150145f610180565b5f80fd5b9060206003198301126101b5576004356001600160401b0381116101b557826023820112156101b5578060040135926001600160401b0384116101b557602484830101116101b5576024019190565b634e487b7160e01b5f52602160045260245ffd5b6003111561022657565b610208565b346101b557602061029261023e366101b9565b9061025160ff60045460401c1615610d16565b7fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5f525f845260405f205f8052845260ff60405f205416156102a357611c3c565b6040519061029f8161021c565b8152f35b6102ab61285c565b611c3c565b346101b55760203660031901126101b55760206102da6004355f525f602052600160405f20015490565b604051908152f35b60409060031901126101b557600435906024356001600160a01b03811681036101b55790565b346101b55761033c610319366102e2565b90610337610332825f525f602052600160405f20015490565b6128cb565b613b11565b005b346101b55761034c366102e2565b336001600160a01b038216036103655761033c91613ba9565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b90602080835192838152019201905f5b8181106103aa5750505090565b825163ffffffff1684526020938401939092019160010161039d565b90602080835192838152019201905f5b8181106103e35750505090565b82518452602093840193909201916001016103d6565b90602080835192838152019201905f5b8181106104165750505090565b82516001600160401b0316845260209384019390920191600101610409565b9161045e9061045061046c959360608652606086019061038d565b9084820360208601526103c6565b9160408184039101526103f9565b90565b346101b55760203660031901126101b55760043561048c81613c39565b6001600160a01b038116156105bb576104be906104b86104b360085461ffff9060e01c1690565b613cd1565b90613d39565b9060206104cb8383613df2565b01906104eb6104e66104df845161ffff1690565b61ffff1690565b611e2b565b906104fe6104e66104df855161ffff1690565b936105116104e66104df865161ffff1690565b9160095490600a54925f5b8861052c6104df8a5161ffff1690565b61ffff8316908110156105a5579161059e60019261059085806105898f8f998e8e8e8e61057b61ffff9f9661057387610568816105819b611e71565b9063ffffffff169052565b5161ffff1690565b91613f34565b929095611e71565b528b611e71565b906001600160401b03169052565b011661051c565b50876105b78860405193849384610435565b0390f35b633a517eed60e21b5f52600482905260245b5ffd5b5f9103126101b557565b346101b5575f3660031901126101b557335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff16156106685760045460ff8160401c16156106405768ff00000000000000001916600455005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b346101b55761068d366101b9565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b346101b557602060ff6106ec6106cc366102e2565b905f525f845260405f20906001600160a01b03165f5260205260405f2090565b54166040519015158152f35b346101b55760203660031901126101b5576004356001600160401b0381116101b557806004019061016060031982360301126101b55761074060ff60045460401c1615610d16565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac4754161561084b575b6101448101906107a28284611e8a565b90501561082357610804610813926105b7946107c16044850182611ebc565b6107d16064879493940183611ebc565b90608488013592610104890135956107e887610fc2565b60a461080b6107fb6101248d0189611ebc565b9b909a89611e8a565b3691610eac565b9a019561401e565b6040519081529081906020820190565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b61085361285c565b610792565b346101b5575f3660031901126101b55760206040515f8152f35b346101b55760203660031901126101b5576004356001600160401b0381116101b5578060040161014060031983360301126101b5576105b791610813916108c160ff60045460401c1615610d16565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610972575b6109206044830182611ebc565b9161092e6064850182611ebc565b61010486013592916084870135919061094685610fc2565b610954610124890185611ebc565b97909660a46040519a61096860208d610e01565b5f8c52019561401e565b61097a61285c565b610913565b346101b55761033c610990366102e2565b906109a9610332825f525f602052600160405f20015490565b613ba9565b346101b5575f3660031901126101b55760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346101b557610a5b6109f9366101b9565b600454610a0c9060401c60ff1615610d16565b7f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260ff610a4a60405f205f805260205260405f2090565b541615610be7575b508101906120cc565b8051906020810190815191604082019384519460806060850195865193610ac583880195610a9087516001600160801b031690565b906040519b8c9586957f265942ff00000000000000000000000000000000000000000000000000000000875260048701612696565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa958615610be257610b8f96610b5f93610b2c925f92610bb1575b508651845190610b268a5193516001600160801b031690565b9361427d565b602083510151610b4260a0860191825190614328565b5050604060208551015192510151905190602086510151926144d6565b6020604080808451015193610b7a60c0870195865190614328565b505051015194510151915192510151926144d6565b61033c6801000000000000000068ff0000000000000000196004541617600455565b610bd491925060803d608011610bdb575b610bcc8183610e01565b8101906121f5565b905f610b0d565b503d610bc2565b61277a565b610bf0906128cb565b5f610a52565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b90602061046c928181520190610bf6565b346101b5575f3660031901126101b5576105b76040516020808201526101406040820152610d0a81610c6061018082016127bd565b610c7c60608301602060ff600254818116845260081c16910152565b610c9e60a0830160206001600160401b03600354818116845260401c16910152565b60045463ffffffff811660e0840152610cfc90602081901c63ffffffff16610100850152610cd7610120850160ff8360401c1615159052565b610ceb610140850160ff8360481c1661222c565b60501c63ffffffff16610160840152565b03601f198101835282610e01565b60405191829182610c1a565b15610d1d57565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b03821117610d7457604052565b610d45565b606081019081106001600160401b03821117610d7457604052565b608081019081106001600160401b03821117610d7457604052565b60a081019081106001600160401b03821117610d7457604052565b60c081019081106001600160401b03821117610d7457604052565b61010081019081106001600160401b03821117610d7457604052565b90601f801991011681019081106001600160401b03821117610d7457604052565b60405190610e3261010083610e01565b565b60405190610e3261026083610e01565b60405190610e326101e083610e01565b60405190610e3261012083610e01565b60405190610e3260e083610e01565b60405190610e32608083610e01565b60405190610e3260c083610e01565b6001600160401b038111610d7457601f01601f191660200190565b929192610eb882610e91565b91610ec66040519384610e01565b8294818452818301116101b5578281602093845f960137010152565b9080601f830112156101b55781602061046c93359101610eac565b60ff8116036101b557565b91908260409103126101b557604051610f2081610d59565b60208082948035610f3081610efd565b8452013591610f3e83610efd565b0152565b6001600160401b038116036101b557565b3590610e3282610f42565b91908260409103126101b557604051610f7681610d59565b60208082948035610f8681610f42565b8452013591610f3e83610f42565b63ffffffff8116036101b557565b3590610e3282610f94565b801515036101b557565b3590610e3282610fad565b600211156101b557565b3590610e3282610fc2565b919091610140818403126101b557610fed610e22565b928135916001600160401b0383116101b5576110328261101561012094611082968501610ee2565b87526110248160208501610f08565b602088015260608301610f5e565b604086015261104360a08201610fa2565b606086015261105460c08201610fa2565b608086015261106560e08201610fb7565b60a08601526110776101008201610fcc565b60c086015201610fa2565b60e0830152565b6001600160801b038116036101b557565b3590610e3282611089565b91908260609103126101b5576040516110bd81610d79565b604080829480356110cd81611089565b8452602081013560208501520135910152565b8092910391606083126101b5576040516110f981610d59565b6040819483358352601f1901126101b557602090604080519361111b85610d59565b8381013561112881610f94565b85520135828401520152565b6001600160401b038111610d745760051b60200190565b919060c0838203126101b55760405161116381610d94565b8093803561117081610f42565b8252602081013561118081610f94565b602083015261119283604083016110e0565b604083015260a0810135906001600160401b0382116101b557019180601f840112156101b5578235926111c484611134565b936111d26040519586610e01565b80855260208086019160051b830101918383116101b55760208101915b83831061120157505050505060600152565b82356001600160401b0381116101b5578201906040828703601f1901126101b5576040519161122f83610d59565b602081013560048110156101b557835260408101356001600160401b0381116101b5576020910101906080828803126101b5576040519261126f84610d94565b82356001600160401b0381116101b5578861128b918501610ee2565b8452602083013561129b81611089565b602085015260408301356112ae81610fad565b60408501526060830135936001600160401b0385116101b5576112d689602096879601610ee2565b6060820152838201528152019201916111ef565b91906040838203126101b5576040519061130382610d59565b819380356001600160401b0381116101b55781016102c0818403126101b55761132a610e34565b906113358482610f5e565b825260408101356001600160401b0381116101b55784611356918301610ee2565b602083015261136760608201610f53565b60408301526113786080820161109a565b606083015261138960a08201610fb7565b608083015261139b8460c083016110e0565b60a08301526113ad6101208201610fb7565b60c083015261014081013560e08301526113ca6101608201610fb7565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526114196102208201610fb7565b6101c08301526102408101356101e08301526114386102608201610fb7565b6102008301526102808101356102208301526102a0810135906001600160401b0382116101b55761146b91859101610ee2565b61024082015283526020810135916001600160401b0383116101b557602092611494920161114b565b910152565b91906080838203126101b557604051906114b282610d94565b819380356001600160401b0381116101b5576060926114d2918301610ee2565b83526020810135602084015260408101356114ec81610f42565b60408401520135908160070b82036101b55760600152565b9190916080818403126101b5576040519061151e82610d94565b819381356001600160401b0381116101b557820181601f820112156101b557803561154881611134565b916115566040519384610e01565b81835260208084019260051b820101918483116101b55760208201905b8382106115c45750505050835261158c60208301610fb7565b60208401526040820135906001600160401b0382116101b557826115b96060949261149494869401611499565b604086015201610f53565b81356001600160401b0381116101b5576020916115e688848094880101611499565b815201910190611573565b919060a0838203126101b5576040519061160a82610d94565b819380356001600160401b0381116101b557826116289183016112ea565b835260208101356001600160401b0381116101b55782611649918301611504565b602084015261165b8260408301610f5e565b60408401526080810135916001600160401b0383116101b5576060926114949201611504565b9080601f830112156101b5576101006040519261169e8285610e01565b839181019283116101b557905b8282106116b85750505090565b81358152602091820191016116ab565b9080601f830112156101b557604051916116e3604084610e01565b8290604081019283116101b557905b8282106116ff5750505090565b81358152602091820191016116f2565b359061ffff821682036101b557565b9080601f830112156101b557813561173581611134565b926117436040519485610e01565b81845260208085019260051b8201019283116101b557602001905b82821061176b5750505090565b60208091833561177a81610f94565b81520191019061175e565b9080601f830112156101b557813561179c81611134565b926117aa6040519485610e01565b81845260208085019260051b8201019283116101b557602001905b8282106117d25750505090565b81358152602091820191016117c5565b9080601f830112156101b55781356117f981611134565b926118076040519485610e01565b81845260208085019260051b8201019283116101b557602001905b82821061182f5750505090565b60208091833561183e81610f42565b815201910190611822565b9080601f830112156101b557813561186081611134565b9261186e6040519485610e01565b81845260208085019260051b8201019283116101b557602001905b8282106118965750505090565b6020809183356118a581610fad565b815201910190611889565b9080601f830112156101b557610200604051926118cd8285610e01565b839181019283116101b557905b8282106118e75750505090565b81358152602091820191016118da565b9080601f830112156101b557610200604051926119148285610e01565b839181019283116101b557905b82821061192e5750505090565b60208091833561193d81610f42565b815201910190611921565b9190610640838203126101b5576040519061196282610daf565b819380358352602081013561197681610efd565b602084015281605f820112156101b55760405161199561020082610e01565b806102408301918483116101b55760408401905b8382106119d7575050611494926119cc85608096946104409460408a01526118b0565b6060870152016118f7565b6020809183356119e681610f94565b8152019101906119a9565b6020818303126101b5578035906001600160401b0382116101b55701610960818303126101b557611a20610e44565b9181356001600160401b0381116101b55781611a3d918401610fd7565b8352611a4c81602084016110a5565b602084015260808201356001600160401b0381116101b55781611a709184016115f1565b6040840152611a8160a0830161109a565b6060840152611a938160c08401611681565b6080840152611aa6816101c084016116c8565b60a0840152611ab98161020084016116c8565b60c0840152611acb610240830161170f565b60e08401526102608201356001600160401b0381116101b55781611af091840161171e565b6101008401526102808201356001600160401b0381116101b55781611b16918401611785565b6101208401526102a08201356001600160401b0381116101b55781611b3c9184016117e2565b6101408401526102c08201356001600160401b0381116101b55781611b6291840161171e565b6101608401526102e08201356001600160401b0381116101b55781611b88918401611849565b6101808401526103008201356001600160401b0381116101b55782611bb58361032093611bc1960161171e565b6101a086015201611948565b6101c082015290565b15611bd3575050565b906001600160401b0380927f7d939bd8000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b610e32909291926060810193604080916001600160801b038151168452602081015160208501520151910152565b611c4c90611c96928101906119f1565b611c67611c5882612912565b84879496939297989598612fdd565b93611c7185613135565b611c7a856131b6565b9515611e1a57611c8a8285613652565b9390935b848685613876565b80611e09575b611df8575b505050611cad8261021c565b81611dad57611da9611d926020604060a0850194611d01611cd884885101516001600160401b031690565b60035460401c6001600160401b03166001600160401b0381166001600160401b03831611611bca565b611d6986516020906001600160401b038151166001600160401b0319600354161760035501517fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6fffffffffffffffff00000000000000006003549260401b16911617600355565b0151604051611d7f81610cfc8582019485611c0e565b519020935101516001600160401b031690565b6001600160401b03165f52600560205260405f2090565b5590565b50611db78161021c565b60018103611de15761046c6801000000000000000068ff0000000000000000196004541617600455565b611dea8161021c565b6002810361046c5750600290565b611e01926138fb565b5f8080611ca1565b50611e138561021c565b8415611c9c565b611e2382613429565b939093611c8e565b90611e3582611134565b611e426040519182610e01565b8281528092611e53601f1991611134565b0190602036910137565b634e487b7160e01b5f52603260045260245ffd5b8051821015611e855760209160051b010190565b611e5d565b903590601e19813603018212156101b557018035906001600160401b0382116101b5576020019181360383136101b557565b903590601e19813603018212156101b557018035906001600160401b0382116101b557602001918160051b360383136101b557565b6002111561022657565b91906060838203126101b55760405190611f1482610d79565b819380356001600160401b0381116101b55781016040818403126101b55760405190611f3f82610d59565b8035906001600160401b0382116101b557611f5e856020938301610ee2565b83520135611f6b81610f42565b6020820152835260208101356001600160401b0381116101b55782611f919183016115f1565b60208401526040810135916001600160401b0383116101b55760409261149492016115f1565b919091610240818403126101b557611fcd610e54565b92611fd88183611681565b8452611fe88161010084016116c8565b6020850152611ffb8161014084016116c8565b604085015261200d610180830161170f565b60608501526101a08201356001600160401b0381116101b5578161203291840161171e565b60808501526101c08201356001600160401b0381116101b55781612057918401611785565b60a08501526101e08201356001600160401b0381116101b5578161207c9184016117e2565b60c08501526102008201356001600160401b0381116101b557816120a191840161171e565b60e08501526102208201356001600160401b0381116101b5576120c49201611849565b610100830152565b6020818303126101b5578035906001600160401b0382116101b55701610160818303126101b5576120fb610e64565b9181356001600160401b0381116101b55781612118918401610fd7565b835260208201356001600160401b0381116101b55781612139918401611efb565b602084015261214b81604084016110a5565b604084015261215d8160a084016110a5565b606084015261216f610100830161109a565b60808401526101208201356001600160401b0381116101b55781612194918401611fb7565b60a08401526101408201356001600160401b0381116101b5576121b79201611fb7565b60c082015290565b91908260409103126101b5576040516121d781610d59565b602080829480516121e781610f42565b8452015191610f3e83610f42565b906080828203126101b55761222490604080519361221285610d59565b61221c83826121bf565b8552016121bf565b602082015290565b9061223682611ef1565b52565b9061046c9061012060e061225885516101408552610140850190610bf6565b9460ff60208083015182815116828801520151166040850152612298604082015160608601906001600160401b0360208092828151168552015116910152565b606081015163ffffffff1660a0850152608081015163ffffffff1660c085015260a08101511515848301526122d660c082015161010086019061222c565b015163ffffffff16910152565b90606060c08201926001600160401b03815116835263ffffffff60208201511660208401526123346040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519160c060a0830152825180915260e0820190602060e08260051b8501019401925f905b82821061236857505050505090565b909192939460df1982820301855285519081516004811015610226576123e682602060019581959482955201519060408482015260606123b483516080604085015260c0840190610bf6565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610bf6565b9701950193920190612359565b9060608061240a8451608085526080850190610bf6565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b8282106124a657505050509060608061249561046c946124836020880151602087019015159052565b604087015185820360408701526123f3565b9401516001600160401b0316910152565b909192936020806124c3600193609f198a820301865288516123f3565b96019201920190929161245a565b61046c91606061265e61264c845160a0855260206126388251604060a089015261251460e0890182516001600160401b0360208092828151168552015116910152565b610240612532848301516102c06101208c01526103a08b0190610bf6565b60408301516001600160401b03166101408b015291808901516001600160801b03166101608b0152608081015115156101808b015260a081015180516101a08c0152602090810151805163ffffffff166101c08d015201516101e08b015260c081015115156102008b015260e08101516102208b01526101008101511515828b01526101208101516102608b01526101408101516102808b01526101608101516102a08b01526101808101516102c08b01526101a08101516102e08b01526101c081015115156103008b01526101e08101516103208b015261020081015115156103408b01526102208101516103608b0152015188820360df19016103808a0152610bf6565b910151858203609f190160c08701526122e3565b60208501518482036020860152612433565b92612686604082015160408501906001600160401b0360208092828151168552015116910152565b0151906080818403910152612433565b90610e329461274661271e610100956126c061276b959b9a989b6101208852610120880190612239565b8681036020880152604061270d8351606084526001600160401b0360206126f2835186606089015260a0880190610bf6565b920151166080850152602085015184820360208601526124d1565b9201519060408184039101526124d1565b986040850190604080916001600160801b038151168452602081015160208501520151910152565b80516001600160801b031660a0840152602081015160c08401526040015160e0830152565b01906001600160801b03169052565b6040513d5f823e3d90fd5b90600182811c921680156127b3575b602083101461279f57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612794565b6001545f92916127cc82612785565b808252916001811690811561284057506001146127e7575050565b60015f9081529293509091907fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b838310612826575060209250010190565b600181602092949394548385870101520191019190612815565b9050602093945060ff929192191683830152151560051b010190565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff161561289457565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260ff6128f23360405f20906001600160a01b03165f5260205260405f2090565b5416156128fc5750565b63e2517d3f60e01b5f523360045260245260445ffd5b5f5f916040810193845161014081515101519561293760406020860151015192614512565b936129418861454a565b909182612a2e576101c00180515190929015612a11575050612964905188614672565b95600194612970614593565b6020845101525b6129f55787851580806129ec575b156129a557505050905051606060208201519101526001915b9493929190565b156129c257505051606001516129ba91614617565b60019161299e565b909492146129d1575b5061299e565b6129e3919350606090510151866147f5565b6001915f6129cb565b50818514612985565b5090506060612a05939293614593565b91510152939291905f90565b959650969050612a2660208351015189614617565b600195612977565b50969094612a3a614593565b602084510152612977565b60405190612a5282610d59565b5f6020838281520152565b60405190612a6a82610d79565b5f6040838281528260208201520152565b60405190612a8882610dca565b81604051612a9581610de5565b60608152612aa1612a45565b6020820152612aae612a45565b60408201525f60608201525f60808201525f60a08201525f60c08201525f60e08201528152612adb612a5d565b6020820152612ae8612a5d565b60408201525f6060820152612afb612a45565b608082015260a0611494612a45565b81601f820112156101b557805190612b2182610e91565b92612b2f6040519485610e01565b828452602083830101116101b557815f9260208093018386015e8301015290565b91908260409103126101b557604051612b6881610d59565b60208082948051612b7881610efd565b8452015191610f3e83610efd565b5190610e3282610f94565b5190610e3282610fad565b5190610e3282610fc2565b919091610140818403126101b557612bbd610e22565b928151916001600160401b0383116101b557612c0282612be561012094611082968501612b0a565b8752612bf48160208501612b50565b6020880152606083016121bf565b6040860152612c1360a08201612b86565b6060860152612c2460c08201612b86565b6080860152612c3560e08201612b91565b60a0860152612c476101008201612b9c565b60c086015201612b86565b5190610e3282611089565b91908260609103126101b557604051612c7581610d79565b60408082948051612c8581611089565b8452602081015160208501520151910152565b6020818303126101b5578051906001600160401b0382116101b55701610180818303126101b55760405191612ccc83610dca565b81516001600160401b0381116101b55782612cef8361014093612d3f9601612ba7565b8552612cfe8360208301612c5d565b6020860152612d108360808301612c5d565b6040860152612d2160e08201612c52565b6060860152612d348361010083016121bf565b6080860152016121bf565b60a082015290565b905f905b60088210612d5857505050565b6020806001928551815201930191019091612d4b565b905f905b60028210612d7f57505050565b6020806001928551815201930191019091612d72565b90602080835192838152019201905f5b818110612db25750505090565b82511515845260209384019390920191600101612da5565b905f905b60108210612ddb57505050565b6020806001928551815201930191019091612dce565b905f905b60108210612e0257505050565b6020806001926001600160401b03865116815201930191019091612df5565b8051825260ff60208201511660208301526040810151604083015f905b60108210612e705750505090610440608083612e666060610e32960151610240860190612dca565b0151910190612df1565b60208060019263ffffffff865116815201930191019091612e3e565b9061046c906103206101c0612fd2612fbe612faa612f96612f82612f6e612f028b612ef06020612ec68d6109608551918181520190612239565b92015160208d0190604080916001600160801b038151168452602081015160208501520151910152565b60408d01518b820360808d01526124d1565b60608c01516001600160801b031660a08b0152612f2760808d015160c08c0190612d47565b612f3860a08d0151898c0190612d6e565b612f4b60c08d01516102008c0190612d6e565b60e08c015161ffff166102408b01526101008c01518a82036102608c015261038d565b6101208b01518982036102808b01526103c6565b6101408a01518882036102a08a01526103f9565b6101608901518782036102c089015261038d565b6101808801518682036102e0880152612d95565b6101a087015185820361030087015261038d565b940151910190612e21565b929190612fe8612a7b565b50156130bc57612ffb5761046c91614941565b5f9061306092613009612a7b565b5060405193849283927f9f3665420000000000000000000000000000000000000000000000000000000084527f0000000000000000000000000000000000000000000000000000000000000000906004850161491f565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115610be2575f916130a0575b5090565b61046c91503d805f833e6130b48183610e01565b810190612c98565b50505f61306091604051809381927fd27520d9000000000000000000000000000000000000000000000000000000008352602060048401526024830190612e8c565b15613107575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b610e329061315281516001600160801b0360608401511690614c32565b6131ae6001600160401b0360206080818501516040516131928482018093604080916001600160801b038151168452602081015160208501520151910152565b606081526131a08382610e01565b519020940151015116614dbb565b8082146130fe565b6131d1611d92602060a084015101516001600160401b031690565b54806131dd5750505f90565b604082019081516040516131f981610cfc602082019485611c0e565b5190201491821592613217575b50501561321257600190565b600290565b6001600160801b0391925061324d61323e60206132599301516001600160801b0390511690565b9351516001600160801b031690565b6001600160801b031690565b911610155f80613206565b1561326b57565b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52601160045260245ffd5b906001600160401b03809116911601906001600160401b0382116132c757565b613293565b156132d45750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b1561330e5750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156133485750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b906003820291808304600314901517156132c757565b908160011b91808304600214901517156132c757565b90602c820291808304602c14901517156132c757565b908160051b91808304602014901517156132c757565b818102929181159184041417156132c757565b156133ee575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b905f905f6040840192602084510151519485519561010082019461348686515161345b6104df60e087015161ffff1690565b8091149081613642575b81613632575b81613622575b81613612575b81613602575b50949194613264565b5f5b8881106135e357505f935f955f965b88515188101561358357908a939291896134c66134c26134bc8c6101808c0151611e71565b51151590565b1590565b6135775791613553916134e96134df8c60019651611e71565b5163ffffffff1690565b809a6134fe8263ffffffff81169a8b106132cc565b61355f575b61352d91506020613514898b611e71565b5101516135268d6101208d0151611e71565b5114613340565b61354d604061353e859b988a611e71565b5101516001600160401b031690565b906132a7565b975b0196909192613497565b63ffffffff61357092168811613306565b5f89613503565b50935096600190613555565b50976060929950610e3296506135de9495506020919793506135ca8a8a6135b26001600160401b03821661337a565b6135c46001600160401b038416613390565b106133e5565b515101510151905161018084015191614e01565b614eab565b93906135fa60019161354d604061353e8988611e71565b919401613488565b90506101a085015151145f61347d565b6101808601515181149150613477565b6101608601515181149150613471565b610140860151518114915061346b565b6101208601515181149150613465565b905f9061010081019161369b8351516136736104df60e086015161ffff1690565b8091149081613866575b81613856575b81613846575b81613836575b81613826575b50613264565b6136a484613c39565b936001600160a01b03851615613814576136d26008969592949654926104b86104b38561ffff9060e01c1690565b6136ef6136df8284613df2565b9360a01c6001600160401b031690565b965f5f9860095497600a549260205f98019b5b8551518910156137cc5791898b94928e9897969461372b6134c26134bc8e610180870151611e71565b6137be5761373d6134df8d8a51611e71565b80926137a0575b8c91508b90876001998c839e5161375c9061ffff1690565b9061376695613f34565b91909361012001519061377891611e71565b51149061378491613340565b61378d916132a7565b976001905b019791939495929092613702565b63ffffffff6137b7921663ffffffff821611613306565b5f81613744565b985094505097600190613792565b50509897509850509050610e329392506135de91506137f887876135b26001600160401b03821661337a565b6060602060408501515101510151905161018084015191614e01565b633a517eed60e21b5f5260045260245ffd5b90506101a084015151145f613695565b610180850151518114915061368f565b6101608501515181149150613689565b6101408501515181149150613683565b610120850151518114915061367d565b9092919260408201936138898551614512565b6138f45760208351015193604060208501510151928314806138e5575b6138d7575050836138c160609283610e329751015190614617565b5101519061018061012082015191015191614f3f565b91509150610e329350614eee565b506060865101515151156138a6565b5050505050565b91906001600160a01b0361390e84613c39565b16151580613aff575b613afa57604001516020015151805193841561326b5760b48511613ae15761393e85611e2b565b925f5b83518110156139855780613974602061395c60019488611e71565b51015161396e604061353e858a611e71565b906160c1565b61397e8288611e71565b5201613941565b5092613a2190613a005f979693966139bb896139a86139a384613390565b613cf6565b94836139b387611e2b565b9c8d926160eb565b506139fa6139ea6139e56139d66139d1856133a6565b613c54565b6139df876133bc565b90613cc4565b613d11565b976139f489616188565b886161d0565b8661625e565b613a1b613a14613a0f856168e1565b6165f1565b86602e0152565b8461626e565b5f9460305b8351871015613a9957613a91600191613a8c613a428a88611e71565b51613a5463ffffffff8c16848b6161ac565b613a75613a6084613c62565b60408301516001600160401b0316908b616216565b6020613a8084613c70565b91015190890160200152565b613c7e565b960195613a26565b95509150925f945b8251861015613ad457613acc600191613ac7613abd8987611e71565b5187830160200152565b613c8c565b950194613aa1565b50935050610e3291615064565b63156f758160e31b5f52600485905260b460245260445ffd5b505050565b5061ffff60085460e01c161515613917565b805f525f60205260ff613b388360405f20906001600160a01b03165f5260205260405f2090565b5416613ba357805f525f602052613b638260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916600117905533916001600160a01b0316907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260ff613bd08360405f20906001600160a01b03165f5260205260405f2090565b541615613ba357805f525f602052613bfc8260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916905533916001600160a01b0316907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b60065403613c50576001600160a01b036008541690565b5f90565b60300190816030116132c757565b90600482018092116132c757565b90600c82018092116132c757565b90602c82018092116132c757565b90602082018092116132c757565b60010190816001116132c757565b60240190816024116132c757565b90600182018092116132c757565b919082018092116132c757565b61ffff16602c810290808204602c14901517156132c757603001806030116132c75790565b5f198101919082116132c757565b919082039182116132c757565b90613d1b82610e91565b613d286040519182610e01565b8281528092611e53601f1991610e91565b9190823b6001811115613d87575f1981019081116132c7578111613d6b576001613d6282613d11565b9360208501903c565b6001600160a01b0383633610565160e21b5f521660045260245ffd5b6001600160a01b0384633610565160e21b5f521660045260245ffd5b60405190613db082610d94565b5f6060838281528260208201528260408201520152565b60011b906201fffe61fffe8316921682036132c757565b61ffff5f199116019061ffff82116132c757565b919091613dfd613da3565b506030835110613eb6576356414c34602084015160e01c03613eb657602483015160c01c92602c81015160f01c93602e82015190604e83015160f01c91613e54613e45610e73565b6001600160401b039093168352565b613e666020830197889061ffff169052565b6040820152613e7d6060820192839061ffff169052565b94613e8a815161ffff1690565b9161ffff8316928315938415613f29575b508315613ef8575b508215613ec8575b50509050613eb65750565b634724a0fd60e01b5f5260045260245ffd5b6134c29250613eea613ee1613ef09551935161ffff1690565b915161ffff1690565b91615268565b805f80613eab565b90925061ffff613f1e6104df613f19613f13875161ffff1690565b94613dc7565b613dde565b91161415915f613ea3565b60b41093505f613e9b565b939092959461ffff63ffffffff821693168310156140015761ffff1691613f5d6139d1826133a6565b9481613f756020888801015160e01c63ffffffff1690565b03613eb657506001901b908116613fd657613f9e613f9285613c62565b84016020015160c01c90565b9516613fb95750613fb161046c92613c70565b016020015190565b9050613fd2915061ffff165f52600b60205260405f2090565b5490565b613ffc613fef8361ffff165f52600c60205260405f2090565b546001600160401b031690565b613f9e565b6303e07d4560e61b5f52600485905263ffffffff1660245260445ffd5b99959699989197909492986140328b611ef1565b8a15614073576105cd8b61404581611ef1565b7f112d89cc000000000000000000000000000000000000000000000000000000005f5260ff16600452602490565b88999a50602090898095969798999a15159081614270575b614094916152b7565b01966140b36140a2896152f5565b6140ac368c6110a5565b908761630f565b5f905f5b888682106141cf575b5050506140cd935061549c565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016803b156101b55783879261413c5f956040519b8c96879586957f87d3a9b100000000000000000000000000000000000000000000000000000000875260048701615869565b03915afa938415610be25761046c95614168956141b5575b506001811161417c575b50505036906110a5565b6001600160801b03633b9aca009151160490565b6001600160801b036141a56141936141ad956152f5565b9361419d87615923565b93369161592d565b911691616453565b5f808061415e565b806141c35f6141c993610e01565b806105d0565b5f614154565b6134c26141f76141f06141ea858b61420896999798996152ff565b80611ebc565b3691615321565b614202368989615321565b906163e1565b61426657509061422f610804614225614245946140cd988c6152ff565b6020810190611e8a565b90815181518082149182614250575b505061538c565b600189935f886140c0565b9091506020840120906020830120145f8061423e565b91906001016140b7565b61ffff811115915061408b565b92602080916142ec6131ae9561429e610e32996001600160401b0397614c32565b6040516142cb8582018093604080916001600160801b038151168452602081015160208501520151910152565b606081526142da608082610e01565b5190206131ae86858a51015116614dbb565b6040516143198382018093604080916001600160801b038151168452602081015160208501520151910152565b606081526131a0608082610e01565b90915f5f906020840151519485519560808201946143788651516143546104df606087015161ffff1690565b80911490816144c7575b816144b8575b816144a9575b816144995750949194613264565b5f5b88811061447a57505f935f955f965b88515188101561443357908a939291896143ae6134c26134bc8c6101008c0151611e71565b6144275791614403916143c76134df8c60019651611e71565b809a6143dc8263ffffffff81169a8b106132cc565b61440f575b61352d915060206143f2898b611e71565b5101516135268d60a08d0151611e71565b975b0196909192614389565b63ffffffff61442092168811613306565b5f896143e1565b50935096600190614405565b509793945095909750610e329450614475915061445d88886135b26001600160401b03821661337a565b60606020845101510151905161010085015191614e01565b6159f8565b939061449160019161354d604061353e8988611e71565b91940161437a565b905061010085015151145f61347d565b60e0860151518114915061436a565b60c08601515181149150614364565b60a0860151518114915061435e565b9291906144e284614512565b61450c576144f96060610e32950191825190614617565b519061010060a082015191015191614f3f565b50505050565b60016001600160401b03602060408281865151015116940151015116016001600160401b0381116132c7576001600160401b03161490565b61455b6001600160a01b0391613c39565b161561456957600754600191565b5f905f90565b6040519061457c82610d94565b5f6060838181528260208201528260408201520152565b604051906145a082610d94565b606082525f60208301526145b261456f565b60408301525f606083015281602090604051916145d0602084610e01565b5f808452805b8181106145e35750505052565b82906145ed61456f565b828288010152016145d6565b15614602575050565b63769a20cb60e01b5f5260045260245260445ffd5b9080515180156146575760b48111614640575090614637610e3292615a2f565b908082146145f9565b63156f758160e31b5f5260045260b460245260445ffd5b82633a517eed60e21b5f5260045260245ffd5b15613eb65750565b9190805161467f81613c39565b906001600160a01b038216156138145750926147db6001600160401b0361475d614780614772610e32966146b561479d9a615ac7565b6146c0818351613df2565b9060208201926146f76146d5855161ffff1690565b61ffff6146ec6104df60085461ffff9060e01c1690565b91161482519061466a565b61472461471161470b602084015160ff1690565b60ff1690565b80151590816147e9575b5082519061466a565b61477b600954958695600a5497818561474c8b8461474582975161ffff1690565b8b85615afa565b9d8f8f90849593955191111561466a565b8861476c8351925161ffff1690565b91615c56565b8b8082146145f9565b615cac565b614798614792613a0f889b949b6168e1565b99600955565b600a55565b167fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6008549260a01b16911617600855565b6147e484600755565b600655565b6010915011155f61471b565b906147ff82613c39565b6001600160a01b0381161561490b57614826906104b86104b360085461ffff9060e01c1690565b6148308184613df2565b915191602083519101906148496104df835161ffff1690565b036148f057600954600a54915f5b85518110156148e75761488061486f835161ffff1690565b858563ffffffff851692898c613f34565b602061488c848a611e71565b51015114908115916148c1575b506148a657600101614857565b6105cd8763769a20cb60e01b5f52906044916004525f602452565b90506001600160401b03806148db604061353e868c611e71565b9216911614155f614899565b50505050505050565b6105cd8463769a20cb60e01b5f52906044916004525f602452565b633a517eed60e21b5f52600483905260245ffd5b61493760409295949395606083526060830190612e8c565b9460208201520152565b613060915f9161494f612a7b565b5060405193849283927f22b541630000000000000000000000000000000000000000000000000000000084527f0000000000000000000000000000000000000000000000000000000000000000906004850161491f565b156149af575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b156149e6575050565b7f12e74add000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b15614a1c5750565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f600154614a5881612785565b908160448501526001811690815f14614aef5750600114614a8f575b50614a8b9192600319848303016024850152610bf6565b0390fd5b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b818310614ad55750919291508101606401614a8b614a74565b805460648488010152859450602090920191600101614abc565b60ff191660648086019190915291151560051b84019091019150614a8b9050614a74565b939291909315614b235750505050565b9160ff8092816084969581604051977fd382033a000000000000000000000000000000000000000000000000000000008952166004880152166024860152166044840152166064820152fd5b15614b78575050565b9063ffffffff80927f73b1bde3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b15614bb9575050565b9063ffffffff80927f1eb5c1eb000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b15614bfa575050565b9063ffffffff80927f87f10409000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b614c4c6001600160801b03610e329316633b9aca00900490565b614c5a8142428211156149a6565b614d7860e0614c698342613d04565b93614d6d600454614c97614c848263ffffffff9060501c1690565b9663ffffffff88169889429111156149dd565b614cca8351805160208201207f000000000000000000000000000000000000000000000000000000000000000014614a14565b614d126020840151614cdd815160ff1690565b6002549160ff831660ff811660ff8416149384614d86575b614d0c9060209060081c60ff165b93015160ff1690565b93614b13565b614d38614d26606085015163ffffffff1690565b63ffffffff8381169082168114614b6f565b614d59614d4c608085015163ffffffff1690565b9160201c63ffffffff1690565b63ffffffff811663ffffffff831614614bb0565b015163ffffffff1690565b9163ffffffff831614614bf1565b9350614d0c6020614d03614d9d8286015160ff1690565b60ff614dae60088a901c821661470b565b9116149692505050614cf5565b6001600160401b03165f52600560205260405f20548015614dd95790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b9291614e108251825114613264565b5f5b8251811015614ea457614e258183611e71565b5115614e9c5763ffffffff614e3a8285611e71565b5116614e4981875181106132cc565b614e538187611e71565b515160048110156102265760011901614e7157506001905b01614e12565b7f5b5c80eb000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b600190614e6b565b5050509050565b610e329060408101519061ffff60e082015116608082015160a083015160c084015190610120850151926101408601519461018061016088015197015197615ef6565b9091614efb908383616092565b15614f04575050565b906001600160401b0380927f4c7d4802000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b93929093614f508551845114613264565b5180519060b4821161504b575f915f5b81811061502e57505f915f975f5b815181101561501b57614f876134c26134bc838b611e71565b61501357614f958183611e71565b515f905b858210614fc8575b5050614fae878787616092565b614fbc576001905b01614f6e565b50505050505050509050565b806020614fd58488611e71565b5101510361500957506001811b808c16614fa157604061353e61354d926150019599949e179d87611e71565b935f80614fa1565b9060010190614f99565b600190614fb6565b5050505090919250610e32939450614eee565b9261504460019161354d604061353e8888611e71565b9301614f60565b63156f758160e31b5f52600482905260b460245260445ffd5b61506e8282613df2565b91604051905f60208301526150988261508a602182018461627e565b03601f198101845283610e01565b61600082511161523c5750610cfc6150f26150e06150b8845161ffff1690565b60f01b7fffff0000000000000000000000000000000000000000000000000000000000001690565b92604051928391602083019586616290565b51905ff0906001600160a01b03821615615214576152029261514b6020926147e46151b2956001600160a01b03167fffffffffffffffffffffffff00000000000000000000000000000000000000006008541617600855565b6151586040820151600755565b6151a961516c82516001600160401b031690565b7fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6008549260a01b16911617600855565b015161ffff1690565b7fffff0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff7dffff000000000000000000000000000000000000000000000000000000006008549260e01b16911617600855565b61520b5f600955565b610e325f600a55565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b9061527290613cd1565b90818114928315615284575b50505090565b90919250621fffe061ffff82169160051b1690808204602014901517156132c75782018092116132c757145f808061527e565b156152bf5750565b7fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b3561046c81610f42565b9190811015611e855760051b81013590603e19813603018212156101b5570190565b92919061532d81611134565b9361533b6040519586610e01565b602085838152019160051b8101918383116101b55781905b838210615361575050505050565b81356001600160401b0381116101b5576020916153818784938701610ee2565b815201910190615353565b91909115615398575050565b90614a8b6153da926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610bf6565b83810360031901602485015290610bf6565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156101b55701602081359101916001600160401b0382116101b55781360383136101b557565b90602083828152019260208260051b82010193835f925b8484106154645750505050505090565b90919293949560208061548c600193601f198682030188526154868b8861540c565b906153ec565b9801940194019294939190615454565b919091156154a8575050565b614a8b6040519283927ffef760c700000000000000000000000000000000000000000000000000000000845260206004850152602484019161543d565b9035601e19823603018112156101b55701602081359101916001600160401b0382116101b5578160051b360383136101b557565b9035607e19823603018112156101b5570190565b359060038210156101b557565b9035605e19823603018112156101b5570190565b90602083828152019260208260051b82010193835f925b8484106155755750505050505090565b9091929394956020806155e7600193601f198682030188526155978b8861553a565b906155a18261552d565b6155aa8161021c565b81526155d96155ce6155be8685018561540c565b60608886015260608501916153ec565b92604081019061540c565b9160408185039101526153ec565b9801940194019294939190615565565b61046c916156c16156b661563a61561f615611868061540c565b6080875260808701916153ec565b61562c602087018761540c565b9086830360208801526153ec565b60806156a661564c6040880188615519565b868403604088015261565d8161552d565b6156668161021c565b84526156746020820161552d565b61567d8161021c565b602085015261568e6040820161552d565b6156978161021c565b6040850152606081019061540c565b91909281606082015201916153ec565b9260608101906154e5565b91606081850391015261554e565b90602083828152019260208260051b82010193835f925b8484106156f65750505050505090565b909192939495601f198282030184528635601e19843603018112156101b5578301906157266020820192806154e5565b8091936020845252604082019060408160051b8401019380935f915b8383106157655750505050505060208060019298019401940192949391906156e6565b909192939495603f1983820301865261577e878361553a565b803561578981610fc2565b61579281611ef1565b82526157b56157a46020830183615519565b6060602085015260608401906155f7565b906040810135609e19823603018112156101b557600193602093849361585b930191604081830391015261584d61582d6158006157f2858061540c565b60a0865260a08601916153ec565b8685013561580d81610fad565b1515878501526158206040860186615519565b84820360408601526155f7565b92606081013561583c81610fad565b151560608401526080810190615519565b9060808184039101526155f7565b980196019493019190615742565b9092809492969593606083019083526060602084015252608081019560808560051b83010194815f90603e1981360301995b8383106158ba57505050505061046c94955060408185039101526156cf565b9091929397607f1986820301825288358b8112156101b557602061591460019386839401906159076158fd6158ef84806154e5565b60408552604085019161543d565b928581019061540c565b91858185039101526153ec565b9a01920193019190939261589b565b3561046c81611089565b9291909261593a84611134565b936159486040519586610e01565b602085828152019060051b8201918383116101b55780915b83831061596e575050505050565b82356001600160401b0381116101b55782016040818703126101b5576040519161599783610d59565b81356001600160401b0381116101b557820187601f820112156101b557878160206159c493359101615321565b83526020820135926001600160401b0384116101b5576159e988602095869501610ee2565b83820152815201920191615960565b610e329161ffff6060820151168151602083015160408401519060a08501519260c08601519461010060e088015197015197615ef6565b805151908115613ba357615a4282611e2b565b915f5b82518051821015615ab95790615aa8613a0f6020615a6584600196611e71565b510151615aa36001600160401b036040615a80878b51611e71565b5101511660405192615a9184610d59565b83526001600160401b03166020830152565b61651d565b615ab28287611e71565b5201615a45565b505090505f61046c92616641565b90813b6001811115613d6b575f1981019081116132c7576001613d6282613d11565b906010811015611e855760051b0190565b91949092936020830194615b156104e661470b885160ff1690565b96615b3b615b2f6008546001600160401b039060a01c1690565b6001600160401b031690565b955f945f5b615b4e61470b8b5160ff1690565b811015615c4a57615b666134df8260408b0151615ae9565b9663ffffffff881690615b7f8961ffff881684106132cc565b82615c30575b505082878787878c5194615b9895613f34565b908260808b015190615ba991615ae9565b516001600160401b03169a8360608c015190615bc491615ae9565b51926001600160401b038d16926001600160401b0316908484831491821592615c25575b50508c51615bf59161466a565b615bfe91613d04565b90615c0891613cc4565b99615c12916160c1565b615c1c828d611e71565b52600101615b40565b14159050845f615be8565b615c439163ffffffff8b5192161061466a565b5f80615b85565b50969750505050505050565b60205f95939461046c98615c789861ffff8998899760ff976040519d8e610de5565b8d52868d015216968760408c01528360608c015260808b01528060a08b01528160c08b01521760e0890152015116946166a3565b9294935f955b615cc361470b602087015160ff1690565b871015615dbd576001600160401b0391600191615cea6104df6134df8b60408b0151615ae9565b91600161ffff84161b9283615d13615d068d60808d0151615ae9565b516001600160401b031690565b92615d228d60608d0151615ae9565b5193615d368c518c8c61ffff88169261685d565b99166001600160401b038216145f14615d815750901916955b8203615d63575050901916965b0195615cb2565b615d799061ffff165f52600b60205260405f2090565b551796615d5c565b615db690615d9b8561ffff165f52600c60205260405f2090565b906001600160401b03166001600160401b0319825416179055565b1795615d4f565b95509392505050565b908160209103126101b5575161046c81610fad565b979260a096615e3c61046c9b97615e28615e69986101608e615e2160c09f9998615e16615e5a9c61ffff615e4b9c1685526020850190612d47565b610120830190612d6e565b0190612d6e565b6102406101a08d01526102408c01906103c6565b908a82036101c08c01526103f9565b908882036101e08a015261038d565b90868203610200880152612d95565b936102208186039101526001600160401b0381511684526001600160401b0360208201511660208501526040810151604085015263ffffffff6060820151166060850152608081015160808501520151918160a08201520190610bf6565b15615ece57565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b92979195602097919561601d95615f0b616919565b508551615f4f8b82015191615f3f615b2f6040615f2f86516001600160401b031690565b935101516001600160401b031690565b6001600160401b03821614616949565b80516001600160401b031696615fdf6040615f7c615f738f86015163ffffffff1690565b63ffffffff1690565b9301518d815191015190615fcd8f80615f9a8563ffffffff90511690565b940151955151015195615fbd615fae610e82565b6001600160401b03909e168e52565b6001600160401b031660208d0152565b60408b015263ffffffff1660608a0152565b608088015260a08701526040519a8b998a997f4cc22bb7000000000000000000000000000000000000000000000000000000008b5260048b01615ddb565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af18015610be257610e32915f91616063575b50615ec7565b616085915060203d60201161608b575b61607d8183610e01565b810190615dc6565b5f61605d565b503d616073565b906001600160401b0360ff6160b36160bd94838360208901511691166133d2565b94511691166133d2565b1090565b613a0f906001600160401b0361046c93604051926160de84610d59565b835216602082015261651d565b91909493948082038281116132c757600181146161695761610b90616986565b92838201928383116132c75760018801928389116132c75783878661613093866160eb565b61613f6139a360019297613390565b8901018093116132c7576161549386926160eb565b6161619061223692616ac0565b938492611e71565b5061617b915061618492959495611e71565b51928392611e71565b5290565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b602d908260081c602c8201530153565b604f908260081c604e8201530153565b805191908290602001825e015f815290565b600a907fffff00000000000000000000000000000000000000000000000000000000000061046c94937f610000000000000000000000000000000000000000000000000000000000000083521660018201527f80600a3d393df3000000000000000000000000000000000000000000000000006003820152019061627e565b9061635e906131ae60405160208101906163468288604080916001600160801b038151168452602081015160208501520151910152565b60608152616355608082610e01565b51902091614dbb565b60208201518082036163b35750506001600160801b03633b9aca009151160461638b8142428211156149a6565b42034281116132c7576001600160801b03610e32911663ffffffff60045416908181106169a2565b7f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b908151815103613ba3575f5b825181101561644b576164008184611e71565b515161640c8284611e71565b5151036164445761641d8184611e71565b516020815191012061642f8284611e71565b516020815191012003616444576001016163ed565b5050505f90565b505050600190565b929190925f5b8451811015614ea45761646c8186611e71565b51908360405160208101906001600160401b038616825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801915f5b8181106164e6575050505081610cfc60019760206164dc940151605f19848303016080850152610bf6565b5190205d01616459565b91939496509194969760208061650860019360bf198b82030188528951610bf6565b970194019101918a96949392989795986164b1565b602460208201906001600160401b03825116806165b6575b5061654261657091613d11565b9261656761656161655b616555876169e6565b876169f2565b86616a09565b85616a20565b90519084616a4d565b90616585615b2f82516001600160401b031690565b61658e57505090565b6165af615b2f6165a161309c9486616a36565b92516001600160401b031690565b9083616a65565b600191505b60808110156165e357506165426165dc6165d761657093613c9a565b613ca8565b9150616535565b60019060071c9101906165bb565b8051600181018091116132c75761660790613d11565b805115611e8557616631816166246020945f868196015382616a9e565b506040519182809261627e565b039060025afa15610be2575f5190565b918181038181116132c757600181146166865761665d90616986565b8201918281116132c757826166729185616641565b9161667d9293616641565b61046c91616ac0565b505061669191611e71565b5190565b5f1981146132c75760010190565b9392909491828414808091616845575b61681d5761ffff831161326b576166ca8783613d04565b906001821461677357506166dd90616986565b9361670a6166eb8689613cc4565b956139df6139a36167046166fe88613cb6565b97613cb6565b92613390565b92815b85811087828a8361674c575b5050501561672f5761672a90616695565b61670d565b9495969791859188616741948b6166a3565b9461667d95966166a3565b63ffffffff92935061676991604060606134df9301510151615ae9565b161087828a616719565b9250505094939294156167bf575050826167ba9161046c9394519160208101516167a2604083015161ffff1690565b9063ffffffff60c060a0850151940151941694613f34565b6160c1565b6167ca829392613cb6565b14908115916167fb575b506167e757608061669192930151611e71565b8251634724a0fd60e01b5f5260045260245ffd5b9050616815615f736134df84604060608901510151615ae9565b14155f6167d4565b50509291505061046c925080519061683f6040602083015192015161ffff1690565b91616b34565b506168586134c2838960e08a0151616b12565b6166b3565b9091602061ffff919594950151169363ffffffff8116948510156168c357506168886139d1856133a6565b936020858401015160e01c03613eb657506020906168bd6168b76168ab86613c62565b83016020015160c01c90565b94613c70565b01015190565b6303e07d4560e61b5f5260049190915263ffffffff1660245260445ffd5b604051906060906168f28284610e01565b6022835261309c91600a906020850190601f19013682375360206021840153600283616a4d565b6040519061692682610dca565b606060a0835f81525f60208201525f60408201525f838201525f60808201520152565b156169515750565b6001600160401b03907f90f4dbed000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b9060015b8060011b908382101561699d575061698a565b925050565b156169ab575050565b906001600160801b0380927fdb020436000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6020600a910153600190565b602082602292010153600181018091116132c75790565b602082600a92010153600181018091116132c75790565b6020828192010153600181018091116132c75790565b602082601092010153600181018091116132c75790565b81602091939293010152602081018091116132c75790565b9092919083016020015b6080821015616a8357906001929391530190565b600180916080607f85161781530193019060071c9092616a6f565b90805191821561644b576021602084930191015e600101806001116132c75790565b604051616ace608082610e01565b60418152616adc6041610e91565b602082019190601f1901368337805115611e85576020935f9360016166319453602183015260418201526040519182809261627e565b918181039081116132c7576001901b905f1982019182116132c7571b16151590565b909161ffff1692602c840293808504602c14901517156132c7578360300193846030116132c7578160051b91808304602014901517156132c7570192603084018091116132c757616b8490613c8c565b825110613eb6575001605001519056fea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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
