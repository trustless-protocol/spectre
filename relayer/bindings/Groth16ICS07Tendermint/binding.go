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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membership_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviour_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"updateClient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_clientState\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCachedValidatorSet\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"updateClientMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CannotHandleMisbehavior\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientNotFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ClientStateMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidMembershipProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"KeyValuePairNotInCache\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProofHeightMismatch\",\"inputs\":[{\"name\":\"expectedRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"expectedRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustThresholdMismatch\",\"inputs\":[{\"name\":\"expectedNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expectedDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnbondingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"UnknownZkAlgorithm\",\"inputs\":[{\"name\":\"algorithm\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorCountExceedsLimit\",\"inputs\":[{\"name\":\"validatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxValidatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"VerificationKeyMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x61016060405234610101576176eb803803809161001b82610119565b6101603960e081610160019112610101576100346101a3565b61003f6101806101ba565b61004a6101a06101ba565b6100556101c06101ba565b6101e0516001600160401b038111610101578561017f82011215610101576100a2958161018061008b93610160015191016101e9565b91610200519361009c6102206101ba565b956107a5565b6040516166a39081610f88823960805181615be6015260a05181613f4e015260c05181610a8e015260e05181818161257d01526127c501526101005181818161098501526109cd01526101205181614a7b015261014051816125480152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b610160601f91909101601f19168101906001600160401b0382119082101761014057604052565b610105565b604081019081106001600160401b0382111761014057604052565b601f909101601f19168101906001600160401b0382119082101761014057604052565b6040519061019260e083610160565b565b60405190610192604083610160565b61016051906001600160a01b038216820361010157565b51906001600160a01b038216820361010157565b6001600160401b03811161014057601f01601f191660200190565b9291926101f5826101ce565b916102036040519384610160565b829481845281830111610101578281602093845f96015e010152565b9080601f83011215610101578151610239926020016101e9565b90565b519060ff8216820361010157565b91908260409103126101015760405161026281610145565b602061027b8183956102738161023c565b85520161023c565b910152565b51906001600160401b038216820361010157565b9190826040910312610101576040516102ac81610145565b602061027b8183956102bd81610280565b855201610280565b519063ffffffff8216820361010157565b5190811515820361010157565b5190600282101561010157565b602081830312610101578051906001600160401b0382116101015701610120818303126101015761031f610183565b8151909290916001600160401b0383116101015761036682610349610100946103a496850161021f565b8652610358816020850161024a565b602087015260608301610294565b604085015261037760a082016102c5565b606085015261038860c082016102c5565b608085015261039960e082016102d6565b60a0850152016102e3565b60c082015290565b90600182811c921680156103da575b60208310146103c657565b634e487b7160e01b5f52602260045260245ffd5b91607f16916103bb565b601f82116103f157505050565b5f5260205f20906020601f840160051c83019310610429575b601f0160051c01905b81811061041e575050565b5f8155600101610413565b909150819061040a565b6002111561043d57565b634e487b7160e01b5f52602160045260245ffd5b90600281101561043d5769ff00000000000000000082549160481b169069ff0000000000000000001916179055565b80518051906001600160401b038211610140576104a9826104a26001546103ac565b60016103e4565b602090601f8311600114610606578260c09361019295936104df935f926105fb575b50508160011b915f199060031b1c19161790565b6001555b6020818101518051600280549284015161ff0060089190911b1660ff90921661ffff199093169290921717905560408083015180516003805492909401516fffffffffffffffff0000000000000000931b929092166001600160401b039092166001600160801b03199091161717905560608101516105799063ffffffff1660049063ffffffff1663ffffffff19825416179055565b6105b161058d608083015163ffffffff1690565b60049067ffffffff0000000082549160201b169067ffffffff000000001916179055565b6105e96105c160a0830151151590565b60049068ff0000000000000000825491151560401b169068ff00000000000000001916179055565b01516105f481610433565b6004610451565b015190505f806104cb565b60015f52601f19831691905f5160206176ab5f395f51905f52925f5b818110610664575092600192859260c09661019298961061064c575b505050811b016001556104e3565b01515f1960f88460031b161c191690555f808061063e565b92936020600181928786015181550195019301610622565b604051905f826001549161068f836103ac565b80835292600181169081156106ff57506001146106b3575b61019292500383610160565b5060015f90815290915f5160206176ab5f395f51905f525b8183106106e3575050906020610192928201016106a7565b60209193508060019154838589010152019101909184926106cb565b6020925061019294915060ff191682840152151560051b8201016106a7565b15610727575050565b6355bace6f60e01b5f9081526001600160401b039182166004529116602452604490fd5b634e487b7160e01b5f52601160045260245ffd5b63ffffffff6107089116019063ffffffff821161077857565b61074b565b15610786575050565b9063ffffffff80926333fae18560e21b5f52166004521660245260445ffd5b91936107e66107eb91969294967f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a6123361010052602080825183010191016102f0565b610480565b6107f361067c565b602081519101206101205261080e61080961067c565b610921565b6101405261088161086861083b602061082d61082861067c565b6109fb565b01516001600160401b031690565b60035490610859906001600160401b0380841691908116821461071e565b60401c6001600160401b031690565b6001600160401b03165f90815260056020526040902090565b556001600160a01b0390811660805290811660a05290811660c0521660e0526004546108d99063ffffffff81166108c76108ba8261075f565b9260201c63ffffffff1690565b9163ffffffff8084169116111561077d565b6001600160a01b03811661090057506108f0610cdf565b506108fd61010051610d61565b50565b8061090d6108fd92610be1565b5061091781610c57565b5061010051610dba565b61092a90610e65565b8051600181018091116107785761094090610e33565b805115610988576020918161095c845f94019284845382610f32565b50604051918291518091835e8101838152039060025afa1561097d575f5190565b6040513d5f823e3d90fd5b634e487b7160e01b5f52603260045260245ffd5b604051906109a982610145565b5f602083606081520152565b8015610778575f190190565b5f1981019190821161077857565b9190820391821161077857565b908151811015610988570160200190565b906001820180921161077857565b610a0361099c565b5080518015908115610bd5575b50610bc6575f19908051805b610b76575b505f198214610b5857600360fc1b6001600160f81b0319610a5b610a4d610a47866109ed565b856109dc565b516001600160f81b03191690565b161480610b62575b610b58575f90610a72836109ed565b915b8151831015610b0957610a93610a8d610a4d85856109dc565b60f81c90565b60ff811660308110908115610afe575b50610af157600a82026001600160401b03908116602f1990920160ff1691909101811691168110610ad957600190920191610a74565b50915050610ae5610194565b9081525f602082015290565b5050915050610ae5610194565b60399150115f610aa3565b9150916001811190811591610b4c575b50610b3d5761023990610b2a610194565b9283526001600160401b03166020830152565b63384c687760e21b5f5260045ffd5b602b915010155f610b19565b9050610ae5610194565b506002610b708383516109cf565b11610a63565b602d60f81b610ba0610b93610a4d610b8d856109c1565b866109dc565b6001600160f81b03191690565b14610bb457610bae906109b5565b80610a1c565b610bbf9192506109c1565b905f610a21565b6329120bff60e21b5f5260045ffd5b6040915010155f610a10565b6001600160a01b0381165f9081525f5160206176cb5f395f51905f52602052604090205460ff16610c52576001600160a01b03165f8181525f5160206176cb5f395f51905f5260205260408120805460ff191660011790553391905f51602061762b5f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f51602061764b5f395f51905f52602052604090205460ff16610c52576001600160a01b0381165f9081525f51602061764b5f395f51905f5260205260409020805460ff1916600117905533906001600160a01b03165f51602061768b5f395f51905f525f51602061762b5f395f51905f525f80a4600190565b5f80525f51602061764b5f395f51905f526020525f51602061766b5f395f51905f525460ff16610d5d575f8080525f51602061764b5f395f51905f526020525f51602061766b5f395f51905f52805460ff1916600117905533905f51602061768b5f395f51905f525f51602061762b5f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f14610c52575f818152602081815260408083208380529091528120805460ff1916600117905533915f51602061762b5f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16610e2d575f818152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905533916001600160a01b0316905f51602061762b5f395f51905f525f80a4600190565b50505f90565b90610e3d826101ce565b610e4a6040519182610160565b8281528092610e5b601f19916101ce565b0190602036910137565b805115610ed15780516001905b6080811015610ec35750806001019081600111610778576001908351010180911161077857610ea3610ebf91610e33565b91610eb9610eb084610eed565b82519085610ef9565b83610f5c565b5090565b60019060071c910190610e72565b50604051610ee0602082610160565b5f81525f36602083013790565b6020600a910153600190565b9092919083016020015b6080821015610f1757906001929391530190565b600180916080607f85161781530193019060071c9092610f03565b908051918215610f54576021602084930191015e600101806001116107785790565b505050600190565b908092918251928315610f8057839260208092019201015e81018091116107785790565b505050509056fe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a7146101245780630bece3561461011f578063248a9ca31461011a5780632f2ff15d1461011557806336568abe14610110578063536c2ad31461010b5780636a28f000146101065780638a8e4c5d1461010157806391d14854146100fc578063974a74c4146100f7578063a217fddf146100f2578063a6f031bb146100ed578063ac9650d8146100e8578063d547741f146100e3578063db3e1fa4146100de578063ddba6537146100d95763ef913a4b146100d4575f80fd5b610b80565b6109a8565b61096e565b61093f565b6108d3565b610743565b610729565b6105c9565b610588565b610550565b6104ab565b610445565b61034e565b610318565b6102c0565b61023b565b346101c55760203660031901126101c5576004357fffffffff0000000000000000000000000000000000000000000000000000000081168091036101c557807f7965db0b000000000000000000000000000000000000000000000000000000006020921490811561019b575b506040519015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150145f610190565b5f80fd5b9060206003198301126101c5576004356001600160401b0381116101c557826023820112156101c5578060040135926001600160401b0384116101c557602484830101116101c5576024019190565b634e487b7160e01b5f52602160045260245ffd5b6003111561023657565b610218565b346101c55760206102a261024e366101c9565b9061026160ff60045460401c1615610cfc565b7fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5f525f845260405f205f8052845260ff60405f205416156102b3576124f8565b604051906102af8161022c565b8152f35b6102bb612fd5565b6124f8565b346101c55760203660031901126101c55760206102ea6004355f525f602052600160405f20015490565b604051908152f35b60409060031901126101c557600435906024356001600160a01b03811681036101c55790565b346101c55761034c610329366102f2565b90610347610342825f525f602052600160405f20015490565b613044565b613a10565b005b346101c55761035c366102f2565b336001600160a01b038216036103755761034c91613aa8565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b90602080835192838152019201905f5b8181106103ba5750505090565b825163ffffffff168452602093840193909201916001016103ad565b90602080835192838152019201905f5b8181106103f35750505090565b82518452602093840193909201916001016103e6565b90602080835192838152019201905f5b8181106104265750505090565b82516001600160401b0316845260209384019390920191600101610419565b346101c55760203660031901126101c55761048161048f61049d61046a600435612890565b91939060405195869560608752606087019061039d565b9085820360208701526103d6565b908382036040850152610409565b0390f35b5f9103126101c557565b346101c5575f3660031901126101c557335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff16156105395760045460ff8160401c16156105115768ff00000000000000001916600455005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b346101c55761055e366101c9565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b346101c557602060ff6105bd61059d366102f2565b905f525f845260405f20906001600160a01b03165f5260205260405f2090565b54166040519015158152f35b346101c55760203660031901126101c5576004356001600160401b0381116101c557806004019061016060031982360301126101c55761061160ff60045460401c1615610cfc565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac4754161561071c575b61014481019061067382846129e2565b9050156106f4576106d56106e49261049d946106926044850182612a14565b6106a26064879493940183612a14565b90608488013592610104890135956106b987610f7c565b60a46106dc6106cc6101248d0189612a14565b9b909a896129e2565b3691610e66565b9a0195613eda565b6040519081529081906020820190565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b610724612fd5565b610663565b346101c5575f3660031901126101c55760206040515f8152f35b346101c55760203660031901126101c5576004356001600160401b0381116101c5578060040161014060031983360301126101c55761049d916106e49161079260ff60045460401c1615610cfc565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610843575b6107f16044830182612a14565b916107ff6064850182612a14565b61010486013592916084870135919061081785610f7c565b610825610124890185612a14565b97909660a46040519a61083960208d610dcb565b5f8c520195613eda565b61084b612fd5565b6107e4565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b602081016020825282518091526040820191602060408360051b8301019401925f915b8383106108a657505050505090565b90919293946020806108c4600193603f198682030187528951610850565b97019301930191939290610897565b346101c55760203660031901126101c5576004356001600160401b0381116101c557366023820112156101c5578060040135906001600160401b0382116101c5573660248360051b830101116101c55761049d9160246109339201612b02565b60405191829182610874565b346101c55761034c610950366102f2565b90610969610342825f525f602052600160405f20015490565b613aa8565b346101c5575f3660031901126101c55760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346101c557610a1a6109b9366101c9565b6109cb60ff60045460401c1615610cfc565b7f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260ff610a0960405f205f805260205260405f2090565b541615610b5d575b50810190612d8f565b805190602081018051906040830190815160806060860194855197610a8283890199610a4d8b516001600160801b031690565b9060405196879586957fa6fe8f5600000000000000000000000000000000000000000000000000000000875260048701612eb9565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa8015610b5857610b049660c095604095610ae3945f94610b27575b5088519051915192516001600160801b03169361417d565b610af760208251015160a086015190614228565b5051015191015190614228565b5061034c6801000000000000000068ff0000000000000000196004541617600455565b610b4a91945060803d608011610b51575b610b428183610dcb565b810190612e82565b925f610acb565b503d610b38565b61249d565b610b6690613044565b5f610a11565b906020610b7d928181520190610850565b90565b346101c5575f3660031901126101c55760405160208082015261012060408201525f600154610bae81612f9d565b90816101608501526001811690815f14610cd75750600114610c76575b61049d83610c6a8185610bf060608301602060ff600254818116845260081c16910152565b610c1260a0830160206001600160401b03600354818116845260401c16910152565b60045463ffffffff811660e0840152610c5c90602081901c63ffffffff16610100850152610c4b610120850160ff8360401c1615159052565b60ff61014085019160481c16611db7565b03601f198101835282610dcb565b60405191829182610b6c565b91905060015f527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6915f905b808210610cbb5750909150810161018001610c6a610bcb565b9192600181602092546101808588010152019101909291610ca2565b60ff19166101808086019190915291151560051b84019091019150610c6a9050610bcb565b15610d0357565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b03821117610d5a57604052565b610d2b565b606081019081106001600160401b03821117610d5a57604052565b608081019081106001600160401b03821117610d5a57604052565b60a081019081106001600160401b03821117610d5a57604052565b60c081019081106001600160401b03821117610d5a57604052565b90601f801991011681019081106001600160401b03821117610d5a57604052565b60405190610dfb60e083610dcb565b565b60405190610dfb61026083610dcb565b60405190610dfb6101c083610dcb565b60405190610dfb61012083610dcb565b60405190610dfb608083610dcb565b60405190610dfb60c083610dcb565b6001600160401b038111610d5a57601f01601f191660200190565b929192610e7282610e4b565b91610e806040519384610dcb565b8294818452818301116101c5578281602093845f960137010152565b9080601f830112156101c557816020610b7d93359101610e66565b60ff8116036101c557565b91908260409103126101c557604051610eda81610d3f565b60208082948035610eea81610eb7565b8452013591610ef883610eb7565b0152565b6001600160401b038116036101c557565b3590610dfb82610efc565b91908260409103126101c557604051610f3081610d3f565b60208082948035610f4081610efc565b8452013591610ef883610efc565b63ffffffff8116036101c557565b3590610dfb82610f4e565b801515036101c557565b3590610dfb82610f67565b600211156101c557565b3590610dfb82610f7c565b919091610120818403126101c557610fa7610dec565b928135916001600160401b0383116101c557610fec82610fcf6101009461102a968501610e9c565b8752610fde8160208501610ec2565b602088015260608301610f18565b6040860152610ffd60a08201610f5c565b606086015261100e60c08201610f5c565b608086015261101f60e08201610f71565b60a086015201610f86565b60c0830152565b6001600160801b038116036101c557565b3590610dfb82611031565b91908260609103126101c55760405161106581610d5f565b6040808294803561107581611031565b8452602081013560208501520135910152565b8092910391606083126101c5576040516110a181610d3f565b6040819483358352601f1901126101c55760209060408051936110c385610d3f565b838101356110d081610f4e565b85520135828401520152565b6001600160401b038111610d5a5760051b60200190565b919060c0838203126101c55760405161110b81610d7a565b8093803561111881610efc565b8252602081013561112881610f4e565b602083015261113a8360408301611088565b604083015260a0810135906001600160401b0382116101c557019180601f840112156101c55782359261116c846110dc565b9361117a6040519586610dcb565b80855260208086019160051b830101918383116101c55760208101915b8383106111a957505050505060600152565b82356001600160401b0381116101c5578201906040828703601f1901126101c557604051916111d783610d3f565b602081013560048110156101c557835260408101356001600160401b0381116101c5576020910101906080828803126101c5576040519261121784610d7a565b82356001600160401b0381116101c55788611233918501610e9c565b8452602083013561124381611031565b6020850152604083013561125681610f67565b60408501526060830135936001600160401b0385116101c55761127e89602096879601610e9c565b606082015283820152815201920191611197565b91906040838203126101c557604051906112ab82610d3f565b819380356001600160401b0381116101c55781016102c0818403126101c5576112d2610dfd565b906112dd8482610f18565b825260408101356001600160401b0381116101c557846112fe918301610e9c565b602083015261130f60608201610f0d565b604083015261132060808201611042565b606083015261133160a08201610f71565b60808301526113438460c08301611088565b60a08301526113556101208201610f71565b60c083015261014081013560e08301526113726101608201610f71565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526113c16102208201610f71565b6101c08301526102408101356101e08301526113e06102608201610f71565b6102008301526102808101356102208301526102a0810135906001600160401b0382116101c55761141391859101610e9c565b61024082015283526020810135916001600160401b0383116101c55760209261143c92016110f3565b910152565b91906080838203126101c5576040519061145a82610d7a565b819380356001600160401b0381116101c55760609261147a918301610e9c565b835260208101356020840152604081013561149481610efc565b60408401520135908160070b82036101c55760600152565b9190916080818403126101c557604051906114c682610d7a565b819381356001600160401b0381116101c557820181601f820112156101c55780356114f0816110dc565b916114fe6040519384610dcb565b81835260208084019260051b820101918483116101c55760208201905b83821061156c5750505050835261153460208301610f71565b60208401526040820135906001600160401b0382116101c557826115616060949261143c94869401611441565b604086015201610f0d565b81356001600160401b0381116101c55760209161158e88848094880101611441565b81520191019061151b565b919060a0838203126101c557604051906115b282610d7a565b819380356001600160401b0381116101c557826115d0918301611292565b835260208101356001600160401b0381116101c557826115f19183016114ac565b60208401526116038260408301610f18565b60408401526080810135916001600160401b0383116101c55760609261143c92016114ac565b9080601f830112156101c557610100604051926116468285610dcb565b839181019283116101c557905b8282106116605750505090565b8135815260209182019101611653565b9080601f830112156101c5576040519161168b604084610dcb565b8290604081019283116101c557905b8282106116a75750505090565b813581526020918201910161169a565b359061ffff821682036101c557565b9080601f830112156101c55781356116dd816110dc565b926116eb6040519485610dcb565b81845260208085019260051b8201019283116101c557602001905b8282106117135750505090565b60208091833561172281610f4e565b815201910190611706565b9080601f830112156101c5578135611744816110dc565b926117526040519485610dcb565b81845260208085019260051b8201019283116101c557602001905b82821061177a5750505090565b813581526020918201910161176d565b9080601f830112156101c55781356117a1816110dc565b926117af6040519485610dcb565b81845260208085019260051b8201019283116101c557602001905b8282106117d75750505090565b6020809183356117e681610efc565b8152019101906117ca565b9080601f830112156101c5578135611808816110dc565b926118166040519485610dcb565b81845260208085019260051b8201019283116101c557602001905b82821061183e5750505090565b60208091833561184d81610f67565b815201910190611831565b9080601f830112156101c557610200604051926118758285610dcb565b839181019283116101c557905b82821061188f5750505090565b8135815260209182019101611882565b9080601f830112156101c557610200604051926118bc8285610dcb565b839181019283116101c557905b8282106118d65750505090565b6020809183356118e581610efc565b8152019101906118c9565b9190610640838203126101c5576040519061190a82610d95565b819380358352602081013561191e81610eb7565b602084015281605f820112156101c55760405161193d61020082610dcb565b806102408301918483116101c55760408401905b83821061197f57505061143c9261197485608096946104409460408a0152611858565b60608701520161189f565b60208091833561198e81610f4e565b815201910190611951565b6020818303126101c5578035906001600160401b0382116101c55701610940818303126101c5576119c8610e0d565b9181356001600160401b0381116101c557816119e5918401610f91565b83526119f4816020840161104d565b602084015260808201356001600160401b0381116101c55781611a18918401611599565b6040840152611a2960a08301611042565b6060840152611a3b8160c08401611629565b6080840152611a4e816101c08401611670565b60a0840152611a61816102008401611670565b60c0840152611a7361024083016116b7565b60e08401526102608201356001600160401b0381116101c55781611a989184016116c6565b6101008401526102808201356001600160401b0381116101c55781611abe91840161172d565b6101208401526102a08201356001600160401b0381116101c55781611ae491840161178a565b6101408401526102c08201356001600160401b0381116101c55781611b0a9184016116c6565b6101608401526102e08201356001600160401b0381116101c55782611b378361030093611b4396016117f1565b610180860152016118f0565b6101a082015290565b81601f820112156101c557805190611b6382610e4b565b92611b716040519485610dcb565b828452602083830101116101c557815f9260208093018386015e8301015290565b91908260409103126101c557604051611baa81610d3f565b60208082948051611bba81610eb7565b8452015191610ef883610eb7565b91908260409103126101c557604051611be081610d3f565b60208082948051611bf081610efc565b8452015191610ef883610efc565b5190610dfb82610f4e565b5190610dfb82610f67565b5190610dfb82610f7c565b919091610120818403126101c557611c35610dec565b928151916001600160401b0383116101c557611c7a82611c5d6101009461102a968501611b4c565b8752611c6c8160208501611b92565b602088015260608301611bc8565b6040860152611c8b60a08201611bfe565b6060860152611c9c60c08201611bfe565b6080860152611cad60e08201611c09565b60a086015201611c14565b5190610dfb82611031565b91908260609103126101c557604051611cdb81610d5f565b60408082948051611ceb81611031565b8452602081015160208501520151910152565b6020818303126101c5578051906001600160401b0382116101c55701610180818303126101c55760405191611d3283610db0565b81516001600160401b0381116101c55782611d558361014093611da59601611c1f565b8552611d648360208301611cc3565b6020860152611d768360808301611cc3565b6040860152611d8760e08201611cb8565b6060860152611d9a836101008301611bc8565b608086015201611bc8565b60a082015290565b6002111561023657565b90611dc182611dad565b52565b90610b7d9061010060c0611de385516101208552610120850190610850565b9460ff60208083015182815116828801520151166040850152611e23604082015160608601906001600160401b0360208092828151168552015116910152565b606081015163ffffffff1660a0850152608081015163ffffffff168483015260a0810151151560e08501520151910190611db7565b90606060c08201926001600160401b03815116835263ffffffff6020820151166020840152611ea96040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519160c060a0830152825180915260e0820190602060e08260051b8501019401925f905b828210611edd57505050505090565b909192939460df198282030185528551908151600481101561023657611f5b8260206001958195948295520151906040848201526060611f2983516080604085015260c0840190610850565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610850565b9701950193920190611ece565b90606080611f7f8451608085526080850190610850565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b82821061201b57505050509060608061200a610b7d94611ff86020880151602087019015159052565b60408701518582036040870152611f68565b9401516001600160401b0316910152565b90919293602080612038600193609f198a82030186528851611f68565b960192019201909291611fcf565b610b7d9160606121d36121c1845160a0855260206121ad8251604060a089015261208960e0890182516001600160401b0360208092828151168552015116910152565b6102406120a7848301516102c06101208c01526103a08b0190610850565b60408301516001600160401b03166101408b015291808901516001600160801b03166101608b0152608081015115156101808b015260a081015180516101a08c0152602090810151805163ffffffff166101c08d015201516101e08b015260c081015115156102008b015260e08101516102208b01526101008101511515828b01526101208101516102608b01526101408101516102808b01526101608101516102a08b01526101808101516102c08b01526101a08101516102e08b01526101c081015115156103008b01526101e08101516103208b015261020081015115156103408b01526102208101516103608b0152015188820360df19016103808a0152610850565b910151858203609f190160c0870152611e58565b60208501518482036020860152611fa8565b926121fb604082015160408501906001600160401b0360208092828151168552015116910152565b0151906080818403910152611fa8565b905f905b6008821061221c57505050565b602080600192855181520193019101909161220f565b905f905b6002821061224357505050565b6020806001928551815201930191019091612236565b90602080835192838152019201905f5b8181106122765750505090565b82511515845260209384019390920191600101612269565b905f905b6010821061229f57505050565b6020806001928551815201930191019091612292565b905f905b601082106122c657505050565b6020806001926001600160401b038651168152019301910190916122b9565b8051825260ff60208201511660208301526040810151604083015f905b60108210612334575050509061044060808361232a6060610dfb96015161024086019061228e565b01519101906122b5565b60208060019263ffffffff865116815201930191019091612302565b90610b7d906103006101a061248161246d6124596124456124316123c36123828b516109408b526109408b0190611dc4565b6123b160208d015160208c0190604080916001600160801b038151168452602081015160208501520151910152565b60408c01518a820360808c0152612046565b60608b01516001600160801b031660a08a01526123e860808c015160c08b019061220b565b6123fb60a08c01516101c08b0190612232565b61240e60c08c01516102008b0190612232565b60e08b015161ffff166102408a01526101008b01518982036102608b015261039d565b6101208a01518882036102808a01526103d6565b6101408901518782036102a0890152610409565b6101608801518682036102c088015261039d565b6101808701518582036102e0870152612259565b9401519101906122e5565b906020610b7d928181520190612350565b6040513d5f823e3d90fd5b6124c060409295949395606083526060830190612350565b9460208201520152565b610dfb909291926060810193604080916001600160801b038151168452602081015160208501520151910152565b61250491810190611999565b61250d8161308b565b929392908415612782575f61257191604051809381927f127dd0520000000000000000000000000000000000000000000000000000000083527f000000000000000000000000000000000000000000000000000000000000000089600485016124a8565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115610b58575f91612760575b50925b6125b9846131c0565b6125c284613241565b9415612751576125d28184613670565b915b612740575b5050506125e58261022c565b816126f55761268461266d6020604060a0850194855161260e848201516001600160401b031690565b6001600160401b0361263b61262f6003546001600160401b039060401c1690565b6001600160401b031690565b911611612688575b50015160405161265a81610c5c85820194856124ca565b519020935101516001600160401b031690565b6001600160401b03165f52600560205260405f2090565b5590565b6126ef906020906001600160401b038151166001600160401b0319600354161760035501517fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6fffffffffffffffff00000000000000006003549260401b16911617600355565b5f612643565b506126ff8161022c565b6001810361272957610b7d6801000000000000000068ff0000000000000000196004541617600455565b6127328161022c565b60028103610b7d5750600290565b612749926137fa565b5f80806125d9565b61275a81613488565b916125d4565b61277c91503d805f833e6127748183610dcb565b810190611cfe565b5f6125ad565b506040517fccd771d60000000000000000000000000000000000000000000000000000000081525f81806127b9876004830161248c565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115610b58575f916127fc575b50926125b0565b61281091503d805f833e6127748183610dcb565b5f6127f5565b60405190612825602083610dcb565b5f808352366020840137565b9061283b826110dc565b6128486040519182610dcb565b8281528092612859601f19916110dc565b0190602036910137565b634e487b7160e01b5f52603260045260245ffd5b805182101561288b5760209160051b010190565b612863565b9061289a82613b38565b6001600160a01b038116156129c0576128cc906128c66128c160085461ffff9060e01c1690565b613bd0565b90613bf5565b9060206128d98385613cae565b01926128f96128f46128ed865161ffff1690565b61ffff1690565b612831565b9361290c6128f46128ed835161ffff1690565b9361291f6128f46128ed845161ffff1690565b9260095491600a54935f5b8861293a6128ed845161ffff1690565b8b61ffff8416918210156129b45760019261299f8380612998612990828f8f8f8f8f61ffff9f9d90612982826129ad9f61298a9461297791612877565b9063ffffffff169052565b5161ffff1690565b91613df0565b929095612877565b528c612877565b906001600160401b03169052565b011661292a565b50505050505050505090565b5090506129cb612816565b6129d3612816565b916129dc612816565b91929190565b903590601e19813603018212156101c557018035906001600160401b0382116101c5576020019181360383136101c557565b903590601e19813603018212156101c557018035906001600160401b0382116101c557602001918160051b360383136101c557565b634e487b7160e01b5f52601160045260245ffd5b5f19810191908211612a6b57565b612a49565b91908203918211612a6b57565b90612a8782610e4b565b612a946040519182610dcb565b8281528092612859601f1991610e4b565b9082101561288b57612abc9160051b8101906129e2565b9091565b805191908290602001825e015f815290565b612af4610dfb92949360208660405197889583870137840101905f8252612ac0565b03601f198101845283610dcb565b612b0b5f610e4b565b612b186040519182610dcb565b5f8152601f19612b275f610e4b565b01366020830137612b37836110dc565b92612b456040519485610dcb565b808452601f19612b54826110dc565b015f5b818110612bad5750505f5b818110612b70575050505090565b80612b91612b8b85612b85600195878a612aa5565b90612ad2565b30614139565b612b9b8288612877565b52612ba68187612877565b5001612b62565b806060602080938901015201612b57565b91906060838203126101c55760405190612bd782610d5f565b819380356001600160401b0381116101c55781016040818403126101c55760405190612c0282610d3f565b8035906001600160401b0382116101c557612c21856020938301610e9c565b83520135612c2e81610efc565b6020820152835260208101356001600160401b0381116101c55782612c54918301611599565b60208401526040810135916001600160401b0383116101c55760409261143c9201611599565b919091610240818403126101c557612c90610e1d565b92612c9b8183611629565b8452612cab816101008401611670565b6020850152612cbe816101408401611670565b6040850152612cd061018083016116b7565b60608501526101a08201356001600160401b0381116101c55781612cf59184016116c6565b60808501526101c08201356001600160401b0381116101c55781612d1a91840161172d565b60a08501526101e08201356001600160401b0381116101c55781612d3f91840161178a565b60c08501526102008201356001600160401b0381116101c55781612d649184016116c6565b60e08501526102208201356001600160401b0381116101c557612d8792016117f1565b610100830152565b6020818303126101c5578035906001600160401b0382116101c55701610160818303126101c557612dbe610dec565b9181356001600160401b0381116101c55781612ddb918401610f91565b835260208201356001600160401b0381116101c55781612dfc918401612bbe565b6020840152612e0e816040840161104d565b6040840152612e208160a0840161104d565b6060840152612e326101008301611042565b60808401526101208201356001600160401b0381116101c55781612e57918401612c7a565b60a08401526101408201356001600160401b0381116101c557612e7a9201612c7a565b60c082015290565b906080828203126101c557612eb1906040805193612e9f85610d3f565b612ea98382611bc8565b855201611bc8565b602082015290565b90610dfb94612f69612f4161010095612ee3612f8e959b9a989b6101208852610120880190611dc4565b86810360208801526040612f308351606084526001600160401b036020612f15835186606089015260a0880190610850565b92015116608085015260208501518482036020860152612046565b920151906040818403910152612046565b986040850190604080916001600160801b038151168452602081015160208501520151910152565b80516001600160801b031660a0840152602081015160c08401526040015160e0830152565b01906001600160801b03169052565b90600182811c92168015612fcb575b6020831014612fb757565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612fac565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff161561300d57565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260ff61306b3360405f20906001600160a01b03165f5260205260405f2090565b5416156130755750565b63e2517d3f60e01b5f523360045260245260445ffd5b5f60408201928351926101408451510151946130af604060208401510151956143cc565b916130b987614404565b909182613172576101a001805151909290156131555750506130dc905187614556565b946001926130e861444d565b6020845101525b6131405782158080613137575b1561311557505051606060208201519101525b93929190565b613121575b505061310f565b606061313092510151906144ea565b5f8061311a565b508782146130fc565b50606061314b61444d565b9151015293929190565b93955095905061316a602083510151886144ea565b6001946130ef565b5095909261317e61444d565b6020845101526130ef565b15613192575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b610dfb906131dd81516001600160801b0360608401511690614a2d565b6132396001600160401b03602060808185015160405161321d8482018093604080916001600160801b038151168452602081015160208501520151910152565b6060815261322b8382610dcb565b519020940151015116614b8b565b808214613189565b61325c61266d602060a084015101516001600160401b031690565b54806132685750505f90565b6040820190815160405161328481610c5c6020820194856124ca565b51902014918215926132a2575b50501561329d57600190565b600290565b6001600160801b039192506132d86132c960206132e49301516001600160801b0390511690565b9351516001600160801b031690565b6001600160801b031690565b911610155f80613291565b156132f657565b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b906001600160401b03809116911601906001600160401b038211612a6b57565b156133465750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156133805750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156133ba5750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b90600382029180830460031490151715612a6b57565b908160011b9180830460021490151715612a6b57565b90602c820291808304602c1490151715612a6b57565b908160051b9180830460201490151715612a6b57565b1561344d575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b905f91602060408201510151518051936101008301926134d88451516134b66128ed60e085015161ffff1690565b8091149081613660575b81613650575b81613640575b81613630575b506132ef565b5f5b86811061361357505f92839283805b8751518710156135d1579089929161351661351261350c8a6101808a0151612877565b51151590565b1590565b6135c6576001916135a39161353961352f8b8d51612877565b5163ffffffff1690565b809961354e8263ffffffff8116998a1061333e565b6135ae575b61357d91506020613564888a612877565b5101516135768c6101208c0151612877565b51146133b2565b61359d604061358e859a9789612877565b5101516001600160401b031690565b9061331e565b965b019590916134e9565b63ffffffff6135bf92168711613378565b5f88613553565b9250956001906135a5565b5091509650610dfb94508691935061360e92506135f66001600160401b0382166133ec565b6136086001600160401b038416613402565b10613444565b614bd1565b9161362960019161359d604061358e8789612877565b92016134da565b905061018083015151145f6134d2565b61016084015151811491506134cc565b61014084015151811491506134c6565b61012084015151811491506134c0565b9061010081019061368f8251516134b66128ed60e085015161ffff1690565b61369883613b38565b926001600160a01b038416156137e75793926136c3600854916128c66128c18461ffff9060e01c1690565b946136e16136d18783613cae565b9260a01c6001600160401b031690565b955f945f5f9860095497600a549260205f98019b5b8551518910156137c05791898b94928e9897969461371f61351261350c8e610180870151612877565b6137b25761373161352f8d8a51612877565b8092613794575b8c91508b90876001998c839e516137509061ffff1690565b9061375a95613df0565b91909361012001519061376c91612877565b511490613778916133b2565b6137819161331e565b976001905b0197919394959290926136f6565b63ffffffff6137ab921663ffffffff821611613378565b5f81613738565b985094505097600190613786565b5050935098505050610dfb945061360e92508691506135f66001600160401b0382166133ec565b633a517eed60e21b5f5260045260245b5ffd5b91906001600160a01b0361380d84613b38565b161515806139fe575b6139f95760400151602001515180519384156132f65760b485116139e05761383d85612831565b925f5b83518110156138845780613873602061385b60019488612877565b51015161386d604061358e858a612877565b90615c4e565b61387d8288612877565b5201613840565b5092613920906138ff5f979693966138ba896138a76138a284613402565b612a5d565b94836138b287612831565b9c8d92615c78565b506138f96138e96138e46138d56138d085613418565b613b53565b6138de8761342e565b90613bc3565b612a7d565b976138f389615d15565b88615d5d565b86615deb565b61391a61391361390e856163e3565b616177565b86602e0152565b84615dfb565b5f9460305b83518710156139985761399060019161398b6139418a88612877565b5161395363ffffffff8c16848b615d39565b61397461395f84613b61565b60408301516001600160401b0316908b615da3565b602061397f84613b6f565b91015190890160200152565b613b7d565b960195613925565b95509150925f945b82518610156139d3576139cb6001916139c66139bc8987612877565b5187830160200152565b613b8b565b9501946139a0565b50935050610dfb91614c14565b63156f758160e31b5f52600485905260b460245260445ffd5b505050565b5061ffff60085460e01c161515613816565b805f525f60205260ff613a378360405f20906001600160a01b03165f5260205260405f2090565b5416613aa257805f525f602052613a628260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916600117905533916001600160a01b0316907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260ff613acf8360405f20906001600160a01b03165f5260205260405f2090565b541615613aa257805f525f602052613afb8260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916905533916001600160a01b0316907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b60065403613b4f576001600160a01b036008541690565b5f90565b6030019081603011612a6b57565b9060048201809211612a6b57565b90600c8201809211612a6b57565b90602c8201809211612a6b57565b9060208201809211612a6b57565b6001019081600111612a6b57565b6024019081602411612a6b57565b9060018201809211612a6b57565b91908201809211612a6b57565b61ffff16602c810290808204602c1490151715612a6b5760300180603011612a6b5790565b9190823b6001811115613c43575f198101908111612a6b578111613c27576001613c1e82612a7d565b9360208501903c565b6001600160a01b0383633610565160e21b5f521660045260245ffd5b6001600160a01b0384633610565160e21b5f521660045260245ffd5b60405190613c6c82610d7a565b5f6060838281528260208201528260408201520152565b60011b906201fffe61fffe831692168203612a6b57565b61ffff5f199116019061ffff8211612a6b57565b919091613cb9613c5f565b506030835110613d72576356414c34602084015160e01c03613d7257602483015160c01c92602c81015160f01c93602e82015190604e83015160f01c91613d10613d01610e2d565b6001600160401b039093168352565b613d226020830197889061ffff169052565b6040820152613d396060820192839061ffff169052565b94613d46815161ffff1690565b9161ffff8316928315938415613de5575b508315613db4575b508215613d84575b50509050613d725750565b634724a0fd60e01b5f5260045260245ffd5b6135129250613da6613d9d613dac9551935161ffff1690565b915161ffff1690565b91614e0a565b805f80613d67565b90925061ffff613dda6128ed613dd5613dcf875161ffff1690565b94613c83565b613c9a565b91161415915f613d5f565b60b41093505f613d57565b939092959461ffff63ffffffff82169316831015613ebd5761ffff1691613e196138d082613418565b9481613e316020888801015160e01c63ffffffff1690565b03613d7257506001901b908116613e9257613e5a613e4e85613b61565b84016020015160c01c90565b9516613e755750613e6d610b7d92613b6f565b016020015190565b9050613e8e915061ffff165f52600b60205260405f2090565b5490565b613eb8613eab8361ffff165f52600c60205260405f2090565b546001600160401b031690565b613e5a565b6303e07d4560e61b5f52600485905263ffffffff1660245260445ffd5b9095979198929996613eeb81611dad565b806140f557509060208993928480151590816140e8575b613f0b91614e59565b0196613f2a613f1989614e97565b613f23368c61104d565b9088615e8a565b5f905f5b88868210614047575b505050613f44935061503e565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016803b156101c557835f8094613fb48a956040519c8d97889687957f87d3a9b10000000000000000000000000000000000000000000000000000000087526004870161540b565b03925af1938415610b5857610b7d95613fe09561402d575b5060018111613ff4575b505050369061104d565b6001600160801b03633b9aca009151160490565b6001600160801b0361401d61400b61402595614e97565b93614015876154c5565b9336916154cf565b911691615fce565b5f8080613fd6565b8061403b5f61404193610dcb565b806104a1565b5f613fcc565b61351261406f614068614062858b6140809699979899614ea1565b80612a14565b3691614ec3565b61407a368989614ec3565b90615f5c565b6140de5750906140a76106d561409d6140bd94613f44988c614ea1565b60208101906129e2565b908151815180821491826140c8575b5050614f2e565b600189935f88613f37565b9091506020840120906020830120145f806140b6565b9190600101613f2e565b61ffff8111159150613f02565b806141026137f792611dad565b61410b81611dad565b7f112d89cc000000000000000000000000000000000000000000000000000000005f5260ff16600452602490565b5f80610b7d93602081519101845af43d15614175573d9161415983610e4b565b926141676040519485610dcb565b83523d5f602085013e61559a565b60609161559a565b92602080916141ec6132399561419e610dfb996001600160401b0397614a2d565b6040516141cb8582018093604080916001600160801b038151168452602081015160208501520151910152565b606081526141da608082610dcb565b51902061323986858a51015116614b8b565b6040516142198382018093604080916001600160801b038151168452602081015160208501520151910152565b6060815261322b608082610dcb565b90915f906020830151519182519460808101936142788551516142536128ed606086015161ffff1690565b80911490816143bd575b816143ae575b8161439f575b8161438f575b509493946132ef565b5f5b87811061437057505f93849384805b88515188101561434457908a939291896142ae61351261350c8c6101008c0151612877565b6143385791614314916142c761352f8c60019651612877565b809a6142dc8263ffffffff81169a8b1061333e565b614320575b614303915060206142f2898b612877565b5101516135768d60a08d0151612877565b61359d604061358e859b988a612877565b975b0196909192614289565b63ffffffff61433192168811613378565b5f896142e1565b50935096600190614316565b50919850969350610dfb955061436b9294508791506135f66001600160401b0382166133ec565b615626565b929361438760019161359d604061358e8887612877565b94930161427a565b905061010084015151145f61426f565b60e08501515181149150614269565b60c08501515181149150614263565b60a0850151518114915061425d565b60016001600160401b03602060408281865151015116940151015116016001600160401b038111612a6b576001600160401b03161490565b6144156001600160a01b0391613b38565b161561442357600754600191565b5f905f90565b6040519061443682610d7a565b5f6060838181528260208201528260408201520152565b6040519061445a82610d7a565b606082525f602083015261446c614429565b60408301525f6060830152816020906040519161448a602084610dcb565b5f808452805b81811061449d5750505052565b82906144a7614429565b82828801015201614490565b156144bc575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b90805151801561452a5760b4811161451357509061450a610dfb9261565d565b908082146144b3565b63156f758160e31b5f5260045260b460245260445ffd5b82633a517eed60e21b5f5260045260245ffd5b15613d725750565b90601081101561288b5760051b0190565b8151929161456384613b38565b936001600160a01b038516156137e75750614580614649946156f5565b61458b818351613cae565b908260208301926145c36145a1855161ffff1690565b61ffff6145b86128ed60085461ffff9060e01c1690565b91161483519061453d565b6145da6145d4602084015160ff1690565b60ff1690565b92831515806147d7575b6145f4908493929594519061453d565b614652600954998a96600a54978161461b8a8361461482965161ffff1690565b8c8b615717565b9b90916146348d6001600160401b03845191111561453d565b8a6146438351925161ffff1690565b91615867565b888082146144b3565b5f935b8385106146db5750505050506001600160401b03839261468f6146cd9361468a61468461390e610dfb996163e3565b99600955565b600a55565b167fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6008549260a01b16911617600855565b6146d684600755565b600655565b909192939860806001918b6001600160401b0398876147046128ed61352f856040850151614545565b9161473961472e61472186600161ffff88161b998a960151614545565b516001600160401b031690565b9460608c0151614545565b519361474d8b518b8b61ffff8816926158d2565b9d166001600160401b038216145f1461479b5750901916995b820361477d575050901916995b0193929190614655565b6147939061ffff165f52600b60205260405f2090565b551799614773565b6147d0906147b58561ffff165f52600c60205260405f2090565b906001600160401b03166001600160401b0319825416179055565b1799614766565b5060108411156145e4565b156147eb575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b15614822575050565b7f12e74add000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b156148585750565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f60015461489481612f9d565b908160448501526001811690815f1461492b57506001146148cb575b506148c79192600319848303016024850152610850565b0390fd5b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b81831061491157509192915081016064016148c76148b0565b8054606484880101528594506020909201916001016148f8565b60ff191660648086019190915291151560051b840190910191506148c790506148b0565b93929190931561495f5750505050565b9160ff8092816084969581604051977fd382033a000000000000000000000000000000000000000000000000000000008952166004880152166024860152166044840152166064820152fd5b156149b4575050565b9063ffffffff80927f73b1bde3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b156149f5575050565b9063ffffffff80927f1eb5c1eb000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b614a6d614a4a6001600160801b03610dfb9416633b9aca00900490565b614a588142428211156147e2565b42610708614a668342612a70565b1115614819565b614aa08151805160208201207f000000000000000000000000000000000000000000000000000000000000000014614850565b614ae86020820151614ab3815160ff1690565b6002549160ff831660ff811660ff8416149384614b56575b614ae29060209060081c60ff165b93015160ff1690565b9361494f565b614b42614b356080614b01606085015163ffffffff1690565b93614b2a60045495614b168763ffffffff1690565b63ffffffff811663ffffffff8316146149ab565b015163ffffffff1690565b9160201c63ffffffff1690565b63ffffffff811663ffffffff8316146149ec565b9350614ae26020614ad9614b6d8286015160ff1690565b60ff614b7e60088a901c82166145d4565b9116149692505050614acb565b6001600160401b03165f52600560205260405f20548015614ba95790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b610dfb9060408101519061ffff60e082015116608082015160a083015160c084015190610120850151926101408601519461018061016088015197015197615ab2565b614c1e8282613cae565b91604051905f6020830152614c3a82612af46021820184612ac0565b616000825111614dde5750610c5c614c94614c82614c5a845161ffff1690565b60f01b7fffff0000000000000000000000000000000000000000000000000000000000001690565b92604051928391602083019586615e0b565b51905ff0906001600160a01b03821615614db657614da492614ced6020926146d6614d54956001600160a01b03167fffffffffffffffffffffffff00000000000000000000000000000000000000006008541617600855565b614cfa6040820151600755565b614d4b614d0e82516001600160401b031690565b7fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6008549260a01b16911617600855565b015161ffff1690565b7fffff0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff7dffff000000000000000000000000000000000000000000000000000000006008549260e01b16911617600855565b614dad5f600955565b610dfb5f600a55565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b90614e1490613bd0565b90818114928315614e26575b50505090565b90919250621fffe061ffff82169160051b169080820460201490151715612a6b578201809211612a6b57145f8080614e20565b15614e615750565b7fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b35610b7d81610efc565b919081101561288b5760051b81013590603e19813603018212156101c5570190565b929190614ecf816110dc565b93614edd6040519586610dcb565b602085838152019160051b8101918383116101c55781905b838210614f03575050505050565b81356001600160401b0381116101c557602091614f238784938701610e9c565b815201910190614ef5565b91909115614f3a575050565b906148c7614f7c926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610850565b83810360031901602485015290610850565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156101c55701602081359101916001600160401b0382116101c55781360383136101c557565b90602083828152019260208260051b82010193835f925b8484106150065750505050505090565b90919293949560208061502e600193601f198682030188526150288b88614fae565b90614f8e565b9801940194019294939190614ff6565b9190911561504a575050565b6148c76040519283927ffef760c7000000000000000000000000000000000000000000000000000000008452602060048501526024840191614fdf565b9035601e19823603018112156101c55701602081359101916001600160401b0382116101c5578160051b360383136101c557565b9035607e19823603018112156101c5570190565b359060038210156101c557565b9035605e19823603018112156101c5570190565b90602083828152019260208260051b82010193835f925b8484106151175750505050505090565b909192939495602080615189600193601f198682030188526151398b886150dc565b90615143826150cf565b61514c8161022c565b815261517b61517061516086850185614fae565b6060888601526060850191614f8e565b926040810190614fae565b916040818503910152614f8e565b9801940194019294939190615107565b610b7d916152636152586151dc6151c16151b38680614fae565b608087526080870191614f8e565b6151ce6020870187614fae565b908683036020880152614f8e565b60806152486151ee60408801886150bb565b86840360408801526151ff816150cf565b6152088161022c565b8452615216602082016150cf565b61521f8161022c565b6020850152615230604082016150cf565b6152398161022c565b60408501526060810190614fae565b9190928160608201520191614f8e565b926060810190615087565b9160608185039101526150f0565b90602083828152019260208260051b82010193835f925b8484106152985750505050505090565b909192939495601f198282030184528635601e19843603018112156101c5578301906152c8602082019280615087565b8091936020845252604082019060408160051b8401019380935f915b838310615307575050505050506020806001929801940194019294939190615288565b909192939495603f1983820301865261532087836150dc565b803561532b81610f7c565b61533481611dad565b825261535761534660208301836150bb565b606060208501526060840190615199565b906040810135609e19823603018112156101c55760019360209384936153fd93019160408183039101526153ef6153cf6153a26153948580614fae565b60a0865260a0860191614f8e565b868501356153af81610f67565b1515878501526153c260408601866150bb565b8482036040860152615199565b9260608101356153de81610f67565b1515606084015260808101906150bb565b906080818403910152615199565b9801960194930191906152e4565b9092809492969593606083019083526060602084015252608081019560808560051b83010194815f90603e1981360301995b83831061545c575050505050610b7d9495506040818503910152615271565b9091929397607f1986820301825288358b8112156101c55760206154b660019386839401906154a961549f6154918480615087565b604085526040850191614fdf565b9285810190614fae565b9185818503910152614f8e565b9a01920193019190939261543d565b35610b7d81611031565b929190926154dc846110dc565b936154ea6040519586610dcb565b602085828152019060051b8201918383116101c55780915b838310615510575050505050565b82356001600160401b0381116101c55782016040818703126101c5576040519161553983610d3f565b81356001600160401b0381116101c557820187601f820112156101c5578781602061556693359101614ec3565b83526020820135926001600160401b0384116101c55761558b88602095869501610e9c565b83820152815201920191615502565b906155d757508051156155af57602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b8151158061561d575b6155e8575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b156155e0565b610dfb9161ffff6060820151168151602083015160408401519060a08501519260c08601519461010060e088015197015197615ab2565b805151908115613aa25761567082612831565b915f5b825180518210156156e757906156d661390e602061569384600196612877565b5101516156d16001600160401b0360406156ae878b51612877565b51015116604051926156bf84610d3f565b83526001600160401b03166020830152565b61609f565b6156e08287612877565b5201615673565b505090505f610b7d926161c7565b90813b6001811115613c27575f198101908111612a6b576001613c1e82612a7d565b919490929360208301946157326128f46145d4885160ff1690565b9661574c61262f6008546001600160401b039060a01c1690565b955f945f5b61575f6145d48b5160ff1690565b81101561585b5761577761352f8260408b0151614545565b9663ffffffff8816906157908961ffff8816841061333e565b82615841575b505082878787878c51946157a995613df0565b908260808b0151906157ba91614545565b516001600160401b03169a8360608c0151906157d591614545565b51926001600160401b038d16926001600160401b0316908484831491821592615836575b50508c516158069161453d565b61580f91612a70565b9061581991613bc3565b9961582391615c4e565b61582d828d612877565b52600101615751565b14159050845f6157f9565b6158549163ffffffff8b5192161061453d565b5f80615796565b50969750505050505050565b94939190604051956101008701948786106001600160401b03871117610d5a57610b7d985f9761ffff89989660ff966020968b996040528d52868d015216968760408c01528360608c015260808b01528060a08b01528160c08b01521760e089015201511694616229565b9091602061ffff919594950151169363ffffffff81169485101561593857506158fd6138d085613418565b936020858401015160e01c03613d72575060209061593261592c61592086613b61565b83016020015160c01c90565b94613b6f565b01015190565b6303e07d4560e61b5f5260049190915263ffffffff1660245260445ffd5b908160209103126101c55751610b7d81610f67565b9060c060a0610b7d936001600160401b0381511684526001600160401b0360208201511660208501526040810151604085015263ffffffff6060820151166060850152608081015160808501520151918160a08201520190610850565b9795939199989694929061ffff168852602088015f905b60088210615a6d575050610b7d98995092615a31615a5e969593615a1d615a4094615a12615a4f986101208e0190612232565b6101608c0190612232565b6102406101a08b01526102408a01906103d6565b908882036101c08a0152610409565b908682036101e088015261039d565b90848203610200860152612259565b9161022081840391015261596b565b6020806001928e518152019c019101909a6159df565b15615a8a57565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b929791956020979195615bd995615ac761641b565b508551615b0b8b82015191615afb61262f6040615aeb86516001600160401b031690565b935101516001600160401b031690565b6001600160401b0382161461644b565b80516001600160401b031696615b9b6040615b38615b2f8f86015163ffffffff1690565b63ffffffff1690565b9301518d815191015190615b898f80615b568563ffffffff90511690565b940151955151015195615b79615b6a610e3c565b6001600160401b03909e168e52565b6001600160401b031660208d0152565b60408b015263ffffffff1660608a0152565b608088015260a08701526040519a8b998a997f4cc22bb7000000000000000000000000000000000000000000000000000000008b5260048b016159c8565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af18015610b5857610dfb915f91615c1f575b50615a83565b615c41915060203d602011615c47575b615c398183610dcb565b810190615956565b5f615c19565b503d615c2f565b61390e906001600160401b03610b7d9360405192615c6b84610d3f565b835216602082015261609f565b9190949394808203828111612a6b5760018114615cf657615c9890616488565b9283820192838311612a6b576001880192838911612a6b57838786615cbd9386615c78565b615ccc6138a260019297613402565b890101809311612a6b57615ce1938692615c78565b615cee90611dc1926165c2565b938492612877565b50615d089150615d1192959495612877565b51928392612877565b5290565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b602d908260081c602c8201530153565b604f908260081c604e8201530153565b600a907fffff000000000000000000000000000000000000000000000000000000000000610b7d94937f610000000000000000000000000000000000000000000000000000000000000083521660018201527f80600a3d393df30000000000000000000000000000000000000000000000000060038201520190612ac0565b90615ed9906132396040516020810190615ec18288604080916001600160801b038151168452602081015160208501520151910152565b60608152615ed0608082610dcb565b51902091614b8b565b6020820151808203615f2e5750506001600160801b03633b9aca0091511604615f068142428211156147e2565b4203428111612a6b576001600160801b03610dfb911663ffffffff60045416908181106164a4565b7f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b908151815103613aa2575f5b8251811015615fc657615f7b8184612877565b5151615f878284612877565b515103615fbf57615f988184612877565b5160208151910120615faa8284612877565b516020815191012003615fbf57600101615f68565b5050505f90565b505050600190565b929190925f5b845181101561609857615fe78186612877565b51908360405160208101906001600160401b038616825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801915f5b818110616061575050505081610c5c6001976020616057940151605f19848303016080850152610850565b5190205d01615fd4565b91939496509194969760208061608360019360bf198b82030188528951610850565b970194019101918a969493929897959861602c565b5050509050565b602460208201906001600160401b038251168061613c575b506160c46160f291612a7d565b926160e96160e36160dd6160d7876164e8565b876164f4565b8661650b565b85616522565b9051908461654f565b9061610761262f82516001600160401b031690565b61611057505090565b61613161262f6161236161389486616538565b92516001600160401b031690565b9083616567565b5090565b600191505b608081101561616957506160c461616261615d6160f293613b99565b613ba7565b91506160b7565b60019060071c910190616141565b805160018101809111612a6b5761618d90612a7d565b80511561288b576161b7816161aa6020945f8681960153826165a0565b5060405191828092612ac0565b039060025afa15610b58575f5190565b91818103818111612a6b576001811461620c576161e390616488565b820191828111612a6b57826161f891856161c7565b9161620392936161c7565b610b7d916165c2565b505061621791612877565b5190565b5f198114612a6b5760010190565b93929094918284148080916163cb575b6163a35761ffff83116132f6576162508783612a70565b90600182146162f9575061626390616488565b936162906162718689613bc3565b956138de6138a261628a61628488613bb5565b97613bb5565b92613402565b92815b85811087828a836162d2575b505050156162b5576162b09061621b565b616293565b94959697918591886162c7948b616229565b946162039596616229565b63ffffffff9293506162ef916040606061352f9301510151614545565b161087828a61629f565b9250505094939294156163455750508261634091610b7d939451916020810151616328604083015161ffff1690565b9063ffffffff60c060a0850151940151941694613df0565b615c4e565b616350829392613bb5565b1490811591616381575b5061636d57608061621792930151612877565b8251634724a0fd60e01b5f5260045260245ffd5b905061639b615b2f61352f84604060608901510151614545565b14155f61635a565b505092915050610b7d92508051906163c56040602083015192015161ffff1690565b91616636565b506163de613512838960e08a0151616614565b616239565b604051906060906163f48284610dcb565b6022835261613891600a906020850190601f1901368237536020602184015360028361654f565b6040519061642882610db0565b606060a0835f81525f60208201525f60408201525f838201525f60808201520152565b156164535750565b6001600160401b03907f90f4dbed000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b9060015b8060011b908382101561649f575061648c565b925050565b156164ad575050565b906001600160801b0380927fdb020436000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6020600a910153600190565b60208260229201015360018101809111612a6b5790565b602082600a9201015360018101809111612a6b5790565b602082819201015360018101809111612a6b5790565b60208260109201015360018101809111612a6b5790565b8160209193929301015260208101809111612a6b5790565b9092919083016020015b608082101561658557906001929391530190565b600180916080607f85161781530193019060071c9092616571565b908051918215615fc6576021602084930191015e60010180600111612a6b5790565b6040516165d0608082610dcb565b604181526165de6041610e4b565b602082019190601f190136833780511561288b576020935f9360016161b794536021830152604182015260405191828092612ac0565b91818103908111612a6b576001901b905f198201918211612a6b571b16151590565b909161ffff1692602c840293808504602c1490151715612a6b57836030019384603011612a6b578160051b9180830460201490151715612a6b57019260308401809111612a6b5761668690613b8b565b825110613d72575001605001519056fea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.Multicall(&_ContractGroth16ICS07Tendermint.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ContractGroth16ICS07Tendermint.Contract.Multicall(&_ContractGroth16ICS07Tendermint.TransactOpts, data)
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
