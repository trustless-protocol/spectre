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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membership_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviour_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"updateClient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_clientState\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCachedValidatorSet\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"updateClientMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CannotHandleMisbehavior\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientNotFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ClientStateMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidMembershipProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"KeyValuePairNotInCache\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProofHeightMismatch\",\"inputs\":[{\"name\":\"expectedRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"expectedRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustThresholdMismatch\",\"inputs\":[{\"name\":\"expectedNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expectedDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnbondingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"UnknownZkAlgorithm\",\"inputs\":[{\"name\":\"algorithm\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorCountExceedsLimit\",\"inputs\":[{\"name\":\"validatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxValidatorCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"VerificationKeyMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x610160604052346100fe57617096803803809161001b82610116565b6101603960e0816101600191126100fe576100346101a0565b61003f6101806101b7565b61004a6101a06101b7565b6100556101c06101b7565b6101e0516001600160401b0381116100fe578561017f820112156100fe576100a2958161018061008b93610160015191016101e6565b91610200519361009c6102206101b7565b956107a2565b6040516160519081610f858239608051816147cf015260a05181613c2e015260c05181505060e05181612c9001526101005181818161096c01526109b301526101205181612db6015261014051818181612c5b01526142c20152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b610160601f91909101601f19168101906001600160401b0382119082101761013d57604052565b610102565b604081019081106001600160401b0382111761013d57604052565b601f909101601f19168101906001600160401b0382119082101761013d57604052565b6040519061018f60e08361015d565b565b6040519061018f60408361015d565b61016051906001600160a01b03821682036100fe57565b51906001600160a01b03821682036100fe57565b6001600160401b03811161013d57601f01601f191660200190565b9291926101f2826101cb565b91610200604051938461015d565b8294818452818301116100fe578281602093845f96015e010152565b9080601f830112156100fe578151610236926020016101e6565b90565b519060ff821682036100fe57565b91908260409103126100fe5760405161025f81610142565b602061027881839561027081610239565b855201610239565b910152565b51906001600160401b03821682036100fe57565b91908260409103126100fe576040516102a981610142565b60206102788183956102ba8161027d565b85520161027d565b519063ffffffff821682036100fe57565b519081151582036100fe57565b519060028210156100fe57565b6020818303126100fe578051906001600160401b0382116100fe5701610120818303126100fe5761031c610180565b8151909290916001600160401b0383116100fe5761036382610346610100946103a196850161021c565b86526103558160208501610247565b602087015260608301610291565b604085015261037460a082016102c2565b606085015261038560c082016102c2565b608085015261039660e082016102d3565b60a0850152016102e0565b60c082015290565b90600182811c921680156103d7575b60208310146103c357565b634e487b7160e01b5f52602260045260245ffd5b91607f16916103b8565b601f82116103ee57505050565b5f5260205f20906020601f840160051c83019310610426575b601f0160051c01905b81811061041b575050565b5f8155600101610410565b9091508190610407565b6002111561043a57565b634e487b7160e01b5f52602160045260245ffd5b90600281101561043a5769ff00000000000000000082549160481b169069ff0000000000000000001916179055565b80518051906001600160401b03821161013d576104a68261049f6001546103a9565b60016103e1565b602090601f8311600114610603578260c09361018f95936104dc935f926105f8575b50508160011b915f199060031b1c19161790565b6001555b6020818101518051600280549284015161ff0060089190911b1660ff90921661ffff199093169290921717905560408083015180516003805492909401516fffffffffffffffff0000000000000000931b929092166001600160401b039092166001600160801b03199091161717905560608101516105769063ffffffff1660049063ffffffff1663ffffffff19825416179055565b6105ae61058a608083015163ffffffff1690565b60049067ffffffff0000000082549160201b169067ffffffff000000001916179055565b6105e66105be60a0830151151590565b60049068ff0000000000000000825491151560401b169068ff00000000000000001916179055565b01516105f181610430565b600461044e565b015190505f806104c8565b60015f52601f19831691905f5160206170565f395f51905f52925f5b818110610661575092600192859260c09661018f989610610649575b505050811b016001556104e0565b01515f1960f88460031b161c191690555f808061063b565b9293602060018192878601518155019501930161061f565b604051905f826001549161068c836103a9565b80835292600181169081156106fc57506001146106b0575b61018f9250038361015d565b5060015f90815290915f5160206170565f395f51905f525b8183106106e057505090602061018f928201016106a4565b60209193508060019154838589010152019101909184926106c8565b6020925061018f94915060ff191682840152151560051b8201016106a4565b15610724575050565b6355bace6f60e01b5f9081526001600160401b039182166004529116602452604490fd5b634e487b7160e01b5f52601160045260245ffd5b63ffffffff6107089116019063ffffffff821161077557565b610748565b15610783575050565b9063ffffffff80926333fae18560e21b5f52166004521660245260445ffd5b91936107e36107e891969294967f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a6123361010052602080825183010191016102ed565b61047d565b6107f0610679565b602081519101206101205261080b610806610679565b61091e565b6101405261087e610865610838602061082a610825610679565b6109f8565b01516001600160401b031690565b60035490610856906001600160401b0380841691908116821461071b565b60401c6001600160401b031690565b6001600160401b03165f90815260056020526040902090565b556001600160a01b0390811660805290811660a05290811660c0521660e0526004546108d69063ffffffff81166108c46108b78261075c565b9260201c63ffffffff1690565b9163ffffffff8084169116111561077a565b6001600160a01b0381166108fd57506108ed610cdc565b506108fa61010051610d5e565b50565b8061090a6108fa92610bde565b5061091481610c54565b5061010051610db7565b61092790610e62565b8051600181018091116107755761093d90610e30565b8051156109855760209181610959845f94019284845382610f2f565b50604051918291518091835e8101838152039060025afa1561097a575f5190565b6040513d5f823e3d90fd5b634e487b7160e01b5f52603260045260245ffd5b604051906109a682610142565b5f602083606081520152565b8015610775575f190190565b5f1981019190821161077557565b9190820391821161077557565b908151811015610985570160200190565b906001820180921161077557565b610a00610999565b5080518015908115610bd2575b50610bc3575f19908051805b610b73575b505f198214610b5557600360fc1b6001600160f81b0319610a58610a4a610a44866109ea565b856109d9565b516001600160f81b03191690565b161480610b5f575b610b55575f90610a6f836109ea565b915b8151831015610b0657610a90610a8a610a4a85856109d9565b60f81c90565b60ff811660308110908115610afb575b50610aee57600a82026001600160401b03908116602f1990920160ff1691909101811691168110610ad657600190920191610a71565b50915050610ae2610191565b9081525f602082015290565b5050915050610ae2610191565b60399150115f610aa0565b9150916001811190811591610b49575b50610b3a5761023690610b27610191565b9283526001600160401b03166020830152565b63384c687760e21b5f5260045ffd5b602b915010155f610b16565b9050610ae2610191565b506002610b6d8383516109cc565b11610a60565b602d60f81b610b9d610b90610a4a610b8a856109be565b866109d9565b6001600160f81b03191690565b14610bb157610bab906109b2565b80610a19565b610bbc9192506109be565b905f610a1e565b6329120bff60e21b5f5260045ffd5b6040915010155f610a0d565b6001600160a01b0381165f9081525f5160206170765f395f51905f52602052604090205460ff16610c4f576001600160a01b03165f8181525f5160206170765f395f51905f5260205260408120805460ff191660011790553391905f516020616fd65f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f516020616ff65f395f51905f52602052604090205460ff16610c4f576001600160a01b0381165f9081525f516020616ff65f395f51905f5260205260409020805460ff1916600117905533906001600160a01b03165f5160206170365f395f51905f525f516020616fd65f395f51905f525f80a4600190565b5f80525f516020616ff65f395f51905f526020525f5160206170165f395f51905f525460ff16610d5a575f8080525f516020616ff65f395f51905f526020525f5160206170165f395f51905f52805460ff1916600117905533905f5160206170365f395f51905f525f516020616fd65f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f14610c4f575f818152602081815260408083208380529091528120805460ff1916600117905533915f516020616fd65f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16610e2a575f818152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905533916001600160a01b0316905f516020616fd65f395f51905f525f80a4600190565b50505f90565b90610e3a826101cb565b610e47604051918261015d565b8281528092610e58601f19916101cb565b0190602036910137565b805115610ece5780516001905b6080811015610ec05750806001019081600111610775576001908351010180911161077557610ea0610ebc91610e30565b91610eb6610ead84610eea565b82519085610ef6565b83610f59565b5090565b60019060071c910190610e6f565b50604051610edd60208261015d565b5f81525f36602083013790565b6020600a910153600190565b9092919083016020015b6080821015610f1457906001929391530190565b600180916080607f85161781530193019060071c9092610f00565b908051918215610f51576021602084930191015e600101806001116107755790565b505050600190565b908092918251928315610f7d57839260208092019201015e81018091116107755790565b505050509056fe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a7146101245780630bece3561461011f578063248a9ca31461011a5780632f2ff15d1461011557806336568abe14610110578063536c2ad31461010b5780636a28f000146101065780638a8e4c5d1461010157806391d14854146100fc578063974a74c4146100f7578063a217fddf146100f2578063a6f031bb146100ed578063ac9650d8146100e8578063d547741f146100e3578063db3e1fa4146100de578063ddba6537146100d95763ef913a4b146100d4575f80fd5b610a23565b61098f565b610955565b610926565b6108ba565b61072a565b610710565b6105b0565b61056f565b610550565b6104ab565b610445565b61034e565b610318565b6102c0565b61023b565b346101c55760203660031901126101c5576004357fffffffff0000000000000000000000000000000000000000000000000000000081168091036101c557807f7965db0b000000000000000000000000000000000000000000000000000000006020921490811561019b575b506040519015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150145f610190565b5f80fd5b9060206003198301126101c5576004356001600160401b0381116101c557826023820112156101c5578060040135926001600160401b0384116101c557602484830101116101c5576024019190565b634e487b7160e01b5f52602160045260245ffd5b6003111561023657565b610218565b346101c55760206102a261024e366101c9565b9061026160ff60045460401c1615610af4565b7fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5f525f845260405f205f8052845260ff60405f205416156102b3576119a3565b604051906102af8161022c565b8152f35b6102bb612010565b6119a3565b346101c55760203660031901126101c55760206102ea6004355f525f602052600160405f20015490565b604051908152f35b60409060031901126101c557600435906024356001600160a01b03811681036101c55790565b346101c55761034c610329366102f2565b90610347610342825f525f602052600160405f20015490565b61207f565b6136f0565b005b346101c55761035c366102f2565b336001600160a01b038216036103755761034c91613788565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b90602080835192838152019201905f5b8181106103ba5750505090565b825163ffffffff168452602093840193909201916001016103ad565b90602080835192838152019201905f5b8181106103f35750505090565b82518452602093840193909201916001016103e6565b90602080835192838152019201905f5b8181106104265750505090565b82516001600160401b0316845260209384019390920191600101610419565b346101c55760203660031901126101c55761048161048f61049d61046a600435611bf4565b91939060405195869560608752606087019061039d565b9085820360208701526103d6565b908382036040850152610409565b0390f35b5f9103126101c557565b346101c5575f3660031901126101c557335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff16156105395760045460ff8160401c16156105115768ff00000000000000001916600455005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b346101c55761055e366101c9565b5050636d40ebe160e11b5f5260045ffd5b346101c557602060ff6105a4610584366102f2565b905f525f845260405f20906001600160a01b03165f5260205260405f2090565b54166040519015158152f35b346101c55760203660031901126101c5576004356001600160401b0381116101c557806004019061016060031982360301126101c5576105f860ff60045460401c1615610af4565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610703575b61014481019061065a8284611d46565b9050156106db576106bc6106cb9261049d946106796044850182611d78565b6106896064879493940183611d78565b90608488013592610104890135956106a087610d7f565b60a46106c36106b36101248d0189611d78565b9b909a89611d46565b3691610c69565b9a0195613bba565b6040519081529081906020820190565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b61070b612010565b61064a565b346101c5575f3660031901126101c55760206040515f8152f35b346101c55760203660031901126101c5576004356001600160401b0381116101c5578060040161014060031983360301126101c55761049d916106cb9161077960ff60045460401c1615610af4565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac4754161561082a575b6107d86044830182611d78565b916107e66064850182611d78565b6101048601359291608487013591906107fe85610d7f565b61080c610124890185611d78565b97909660a46040519a61082060208d610bde565b5f8c520195613bba565b610832612010565b6107cb565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b602081016020825282518091526040820191602060408360051b8301019401925f915b83831061088d57505050505090565b90919293946020806108ab600193603f198682030187528951610837565b9701930193019193929061087e565b346101c55760203660031901126101c5576004356001600160401b0381116101c557366023820112156101c5578060040135906001600160401b0382116101c5573660248360051b830101116101c55761049d91602461091a9201611e70565b6040519182918261085b565b346101c55761034c610937366102f2565b90610950610342825f525f602052600160405f20015490565b613788565b346101c5575f3660031901126101c55760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346101c55761099d366101c9565b50506109b160ff60045460401c1615610af4565b7f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610a00575b636d40ebe160e11b5f5260045ffd5b610a099061207f565b5f6109f1565b906020610a20928181520190610837565b90565b346101c5575f3660031901126101c55761049d6040516020808201526101206040820152610ae881610a586101608201611f64565b60ff600254818116606085015260081c166080830152610a9060a0830160206001600160401b03600354818116845260401c16910152565b60045463ffffffff811660e0840152610ada90602081901c63ffffffff16610100850152610ac9610120850160ff8360401c1615159052565b60ff61014085019160481c16612003565b03601f198101835282610bde565b60405191829182610a0f565b15610afb57565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b03821117610b5257604052565b610b23565b606081019081106001600160401b03821117610b5257604052565b608081019081106001600160401b03821117610b5257604052565b60a081019081106001600160401b03821117610b5257604052565b60c081019081106001600160401b03821117610b5257604052565b60e081019081106001600160401b03821117610b5257604052565b90601f801991011681019081106001600160401b03821117610b5257604052565b60405190610c0e60e083610bde565b565b60405190610c0e61026083610bde565b60405190610c0e6101e083610bde565b60405190610c0e608083610bde565b60405190610c0e60c083610bde565b6001600160401b038111610b5257601f01601f191660200190565b929192610c7582610c4e565b91610c836040519384610bde565b8294818452818301116101c5578281602093845f960137010152565b9080601f830112156101c557816020610a2093359101610c69565b60ff8116036101c557565b91908260409103126101c557604051610cdd81610b37565b60208082948035610ced81610cba565b8452013591610cfb83610cba565b0152565b6001600160401b038116036101c557565b3590610c0e82610cff565b91908260409103126101c557604051610d3381610b37565b60208082948035610d4381610cff565b8452013591610cfb83610cff565b63ffffffff8116036101c557565b3590610c0e82610d51565b801515036101c557565b3590610c0e82610d6a565b600211156101c557565b3590610c0e82610d7f565b919091610120818403126101c557610daa610bff565b928135916001600160401b0383116101c557610def82610dd261010094610e2d968501610c9f565b8752610de18160208501610cc5565b602088015260608301610d1b565b6040860152610e0060a08201610d5f565b6060860152610e1160c08201610d5f565b6080860152610e2260e08201610d74565b60a086015201610d89565b60c0830152565b6001600160801b038116036101c557565b3590610c0e82610e34565b91908260609103126101c557604051610e6881610b57565b60408082948035610e7881610e34565b8452602081013560208501520135910152565b8092910391606083126101c557604051610ea481610b37565b6040819483358352601f1901126101c5576020906040805193610ec685610b37565b83810135610ed381610d51565b85520135828401520152565b6001600160401b038111610b525760051b60200190565b919060c0838203126101c557604051610f0e81610b72565b80938035610f1b81610cff565b82526020810135610f2b81610d51565b6020830152610f3d8360408301610e8b565b604083015260a0810135906001600160401b0382116101c557019180601f840112156101c557823592610f6f84610edf565b93610f7d6040519586610bde565b80855260208086019160051b830101918383116101c55760208101915b838310610fac57505050505060600152565b82356001600160401b0381116101c5578201906040828703601f1901126101c55760405191610fda83610b37565b602081013560048110156101c557835260408101356001600160401b0381116101c5576020910101906080828803126101c5576040519261101a84610b72565b82356001600160401b0381116101c55788611036918501610c9f565b8452602083013561104681610e34565b6020850152604083013561105981610d6a565b60408501526060830135936001600160401b0385116101c55761108189602096879601610c9f565b606082015283820152815201920191610f9a565b91906040838203126101c557604051906110ae82610b37565b819380356001600160401b0381116101c55781016102c0818403126101c5576110d5610c10565b906110e08482610d1b565b825260408101356001600160401b0381116101c55784611101918301610c9f565b602083015261111260608201610d10565b604083015261112360808201610e45565b606083015261113460a08201610d74565b60808301526111468460c08301610e8b565b60a08301526111586101208201610d74565b60c083015261014081013560e08301526111756101608201610d74565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526111c46102208201610d74565b6101c08301526102408101356101e08301526111e36102608201610d74565b6102008301526102808101356102208301526102a0810135906001600160401b0382116101c55761121691859101610c9f565b61024082015283526020810135916001600160401b0383116101c55760209261123f9201610ef6565b910152565b91906080838203126101c5576040519061125d82610b72565b819380356001600160401b0381116101c55760609261127d918301610c9f565b835260208101356020840152604081013561129781610cff565b60408401520135908160070b82036101c55760600152565b9190916080818403126101c557604051906112c982610b72565b819381356001600160401b0381116101c557820181601f820112156101c55780356112f381610edf565b916113016040519384610bde565b81835260208084019260051b820101918483116101c55760208201905b83821061136f5750505050835261133760208301610d74565b60208401526040820135906001600160401b0382116101c557826113646060949261123f94869401611244565b604086015201610d10565b81356001600160401b0381116101c55760209161139188848094880101611244565b81520191019061131e565b919060a0838203126101c557604051906113b582610b72565b819380356001600160401b0381116101c557826113d3918301611095565b835260208101356001600160401b0381116101c557826113f49183016112af565b60208401526114068260408301610d1b565b60408401526080810135916001600160401b0383116101c55760609261123f92016112af565b9080601f830112156101c557610100604051926114498285610bde565b839181019283116101c557905b8282106114635750505090565b8135815260209182019101611456565b9080601f830112156101c5576040519161148e604084610bde565b8290604081019283116101c557905b8282106114aa5750505090565b813581526020918201910161149d565b359061ffff821682036101c557565b9080601f830112156101c55781356114e081610edf565b926114ee6040519485610bde565b81845260208085019260051b8201019283116101c557602001905b8282106115165750505090565b60208091833561152581610d51565b815201910190611509565b9080601f830112156101c557813561154781610edf565b926115556040519485610bde565b81845260208085019260051b8201019283116101c557602001905b82821061157d5750505090565b8135815260209182019101611570565b9080601f830112156101c55781356115a481610edf565b926115b26040519485610bde565b81845260208085019260051b8201019283116101c557602001905b8282106115da5750505090565b6020809183356115e981610cff565b8152019101906115cd565b9080601f830112156101c557813561160b81610edf565b926116196040519485610bde565b81845260208085019260051b8201019283116101c557602001905b8282106116415750505090565b60208091833561165081610d6a565b815201910190611634565b9080601f830112156101c557610200604051926116788285610bde565b839181019283116101c557905b8282106116925750505090565b8135815260209182019101611685565b9080601f830112156101c557610200604051926116bf8285610bde565b839181019283116101c557905b8282106116d95750505090565b6020809183356116e881610cff565b8152019101906116cc565b9190610640838203126101c5576040519061170d82610b8d565b819380358352602081013561172181610cba565b602084015281605f820112156101c55760405161174061020082610bde565b806102408301918483116101c55760408401905b83821061178257505061123f9261177785608096946104409460408a015261165b565b6060870152016116a2565b60208091833561179181610d51565b815201910190611754565b6020818303126101c5578035906001600160401b0382116101c55701610960818303126101c5576117cb610c20565b9181356001600160401b0381116101c557816117e8918401610d94565b83526117f78160208401610e50565b602084015260808201356001600160401b0381116101c5578161181b91840161139c565b604084015261182c60a08301610e45565b606084015261183e8160c0840161142c565b6080840152611851816101c08401611473565b60a0840152611864816102008401611473565b60c084015261187661024083016114ba565b60e08401526102608201356001600160401b0381116101c5578161189b9184016114c9565b6101008401526102808201356001600160401b0381116101c557816118c1918401611530565b6101208401526102a08201356001600160401b0381116101c557816118e791840161158d565b6101408401526102c08201356001600160401b0381116101c5578161190d9184016114c9565b6101608401526102e08201356001600160401b0381116101c557816119339184016115f4565b6101808401526103008201356001600160401b0381116101c55782611960836103209361196c96016114c9565b6101a0860152016116f3565b6101c082015290565b610c0e909291926060810193604080916001600160801b038151168452602081015160208501520151910152565b6119b3906119c99281019061179c565b6119bc816120c6565b8486979392949597612c01565b926119d384612d5e565b6119dc84612f0b565b9415611b6b576119ec8184613350565b915b611b5a575b5050506119ff8261022c565b81611b0f57611a9e611a876020604060a08501948551611a28848201516001600160401b031690565b6001600160401b03611a55611a496003546001600160401b039060401c1690565b6001600160401b031690565b911611611aa2575b500151604051611a7481610ada8582019485611975565b519020935101516001600160401b031690565b6001600160401b03165f52600560205260405f2090565b5590565b611b09906020906001600160401b038151166001600160401b0319600354161760035501517fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6fffffffffffffffff00000000000000006003549260401b16911617600355565b5f611a5d565b50611b198161022c565b60018103611b4357610a206801000000000000000068ff0000000000000000196004541617600455565b611b4c8161022c565b60028103610a205750600290565b611b63926134da565b5f80806119f3565b611b7481613152565b916119ee565b60405190611b89602083610bde565b5f808352366020840137565b90611b9f82610edf565b611bac6040519182610bde565b8281528092611bbd601f1991610edf565b0190602036910137565b634e487b7160e01b5f52603260045260245ffd5b8051821015611bef5760209160051b010190565b611bc7565b90611bfe82613818565b6001600160a01b03811615611d2457611c3090611c2a611c2560085461ffff9060e01c1690565b6138b0565b906138d5565b906020611c3d838561398e565b0192611c5d611c58611c51865161ffff1690565b61ffff1690565b611b95565b93611c70611c58611c51835161ffff1690565b93611c83611c58611c51845161ffff1690565b9260095491600a54935f5b88611c9e611c51845161ffff1690565b8b61ffff841691821015611d1857600192611d038380611cfc611cf4828f8f8f8f8f61ffff9f9d90611ce682611d119f611cee94611cdb91611bdb565b9063ffffffff169052565b5161ffff1690565b91613ad0565b929095611bdb565b528c611bdb565b906001600160401b03169052565b0116611c8e565b50505050505050505090565b509050611d2f611b7a565b611d37611b7a565b91611d40611b7a565b91929190565b903590601e19813603018212156101c557018035906001600160401b0382116101c5576020019181360383136101c557565b903590601e19813603018212156101c557018035906001600160401b0382116101c557602001918160051b360383136101c557565b6002111561023657565b634e487b7160e01b5f52601160045260245ffd5b5f19810191908211611dd957565b611db7565b91908203918211611dd957565b90611df582610c4e565b611e026040519182610bde565b8281528092611bbd601f1991610c4e565b90821015611bef57611e2a9160051b810190611d46565b9091565b805191908290602001825e015f815290565b611e62610c0e92949360208660405197889583870137840101905f8252611e2e565b03601f198101845283610bde565b611e795f610c4e565b611e866040519182610bde565b5f8152601f19611e955f610c4e565b01366020830137611ea583610edf565b92611eb36040519485610bde565b808452601f19611ec282610edf565b015f5b818110611f1b5750505f5b818110611ede575050505090565b80611eff611ef985611ef3600195878a611e13565b90611e40565b30613e19565b611f098288611bdb565b52611f148187611bdb565b5001611ed0565b806060602080938901015201611ec5565b90600182811c92168015611f5a575b6020831014611f4657565b634e487b7160e01b5f52602260045260245ffd5b91607f1691611f3b565b6001545f9291611f7382611f2c565b8082529160018116908115611fe75750600114611f8e575050565b60015f9081529293509091907fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b838310611fcd575060209250010190565b600181602092949394548385870101520191019190611fbc565b9050602093945060ff929192191683830152151560051b010190565b9061200d82611dad565b52565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff161561204857565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260ff6120a63360405f20906001600160a01b03165f5260205260405f2090565b5416156120b05750565b63e2517d3f60e01b5f523360045260245260445ffd5b5f5f91604081019384516101408151510151956120eb60406020860151015192613e5d565b936120f588613e95565b9091826121e2576101c001805151909290156121c5575050612118905188613fbd565b95600194612124613ede565b6020845101525b6121a95787851580806121a0575b1561215957505050905051606060208201519101526001915b9493929190565b15612176575050516060015161216e91613f62565b600191612152565b90949214612185575b50612152565b6121979193506060905101518661413a565b6001915f61217f565b50818514612139565b50905060606121b9939293613ede565b91510152939291905f90565b9596509690506121da60208351015189613f62565b60019561212b565b509690946121ee613ede565b60208451015261212b565b6040519061220682610b37565b5f6020838281520152565b6040519061221e82610b57565b5f6040838281528260208201520152565b6040519061223c82610ba8565b8160405161224981610bc3565b606081526122556121f9565b60208201526122626121f9565b60408201525f60608201525f60808201525f60a08201525f60c08201528152612289612211565b6020820152612296612211565b60408201525f60608201526122a96121f9565b608082015260a061123f6121f9565b81601f820112156101c5578051906122cf82610c4e565b926122dd6040519485610bde565b828452602083830101116101c557815f9260208093018386015e8301015290565b91908260409103126101c55760405161231681610b37565b6020808294805161232681610cba565b8452015191610cfb83610cba565b91908260409103126101c55760405161234c81610b37565b6020808294805161235c81610cff565b8452015191610cfb83610cff565b5190610c0e82610d51565b5190610c0e82610d6a565b5190610c0e82610d7f565b919091610120818403126101c5576123a1610bff565b928151916001600160401b0383116101c5576123e6826123c961010094610e2d9685016122b8565b87526123d881602085016122fe565b602088015260608301612334565b60408601526123f760a0820161236a565b606086015261240860c0820161236a565b608086015261241960e08201612375565b60a086015201612380565b5190610c0e82610e34565b91908260609103126101c55760405161244781610b57565b6040808294805161245781610e34565b8452602081015160208501520151910152565b6020818303126101c5578051906001600160401b0382116101c55701610180818303126101c5576040519161249e83610ba8565b81516001600160401b0381116101c557826124c18361014093612511960161238b565b85526124d0836020830161242f565b60208601526124e2836080830161242f565b60408601526124f360e08201612424565b6060860152612506836101008301612334565b608086015201612334565b60a082015290565b90610a209061010060c061253885516101208552610120850190610837565b9460ff60208083015182815116828801520151166040850152612578604082015160608601906001600160401b0360208092828151168552015116910152565b606081015163ffffffff1660a0850152608081015163ffffffff168483015260a0810151151560e08501520151910190612003565b90606060c08201926001600160401b03815116835263ffffffff60208201511660208401526125fe6040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519160c060a0830152825180915260e0820190602060e08260051b8501019401925f905b82821061263257505050505090565b909192939460df1982820301855285519081516004811015610236576126b0826020600195819594829552015190604084820152606061267e83516080604085015260c0840190610837565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610837565b9701950193920190612623565b906060806126d48451608085526080850190610837565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b82821061277057505050509060608061275f610a209461274d6020880151602087019015159052565b604087015185820360408701526126bd565b9401516001600160401b0316910152565b9091929360208061278d600193609f198a820301865288516126bd565b960192019201909291612724565b610a20916060612928612916845160a0855260206129028251604060a08901526127de60e0890182516001600160401b0360208092828151168552015116910152565b6102406127fc848301516102c06101208c01526103a08b0190610837565b60408301516001600160401b03166101408b015291808901516001600160801b03166101608b0152608081015115156101808b015260a081015180516101a08c0152602090810151805163ffffffff166101c08d015201516101e08b015260c081015115156102008b015260e08101516102208b01526101008101511515828b01526101208101516102608b01526101408101516102808b01526101608101516102a08b01526101808101516102c08b01526101a08101516102e08b01526101c081015115156103008b01526101e08101516103208b015261020081015115156103408b01526102208101516103608b0152015188820360df19016103808a0152610837565b910151858203609f190160c08701526125ad565b602085015184820360208601526126fd565b92612950604082015160408501906001600160401b0360208092828151168552015116910152565b01519060808184039101526126fd565b905f905b6008821061297157505050565b6020806001928551815201930191019091612964565b905f905b6002821061299857505050565b602080600192855181520193019101909161298b565b90602080835192838152019201905f5b8181106129cb5750505090565b825115158452602093840193909201916001016129be565b905f905b601082106129f457505050565b60208060019285518152019301910190916129e7565b905f905b60108210612a1b57505050565b6020806001926001600160401b03865116815201930191019091612a0e565b8051825260ff60208201511660208301526040810151604083015f905b60108210612a895750505090610440608083612a7f6060610c0e9601516102408601906129e3565b0151910190612a0a565b60208060019263ffffffff865116815201930191019091612a57565b90610a20906103206101c0612beb612bd7612bc3612baf612b9b612b87612b1b8b612b096020612adf8d6109608551918181520190612519565b92015160208d0190604080916001600160801b038151168452602081015160208501520151910152565b60408d01518b820360808d015261279b565b60608c01516001600160801b031660a08b0152612b4060808d015160c08c0190612960565b612b5160a08d0151898c0190612987565b612b6460c08d01516102008c0190612987565b60e08c015161ffff166102408b01526101008c01518a82036102608c015261039d565b6101208b01518982036102808b01526103d6565b6101408a01518882036102a08a0152610409565b6101608901518782036102c089015261039d565b6101808801518682036102e08801526129ae565b6101a087015185820361030087015261039d565b940151910190612a3a565b6040513d5f823e3d90fd5b929190612c0c61222f565b5015612ce557612c1f57610a2091614286565b5f90612c8492612c2d61222f565b5060405193849283927ffe0a61f30000000000000000000000000000000000000000000000000000000084527f00000000000000000000000000000000000000000000000000000000000000009060048501614264565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115612ce0575f91612cc4575b5090565b610a2091503d805f833e612cd88183610bde565b81019061246a565b612bf6565b50505f612c8491604051809381927ffe284b13000000000000000000000000000000000000000000000000000000008352602060048401526024830190612aa5565b15612d30575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b610c0e90612e918151612da8612d856001600160801b03606086015116633b9aca00900490565b612d938142428211156142eb565b42610708612da18342611dde565b1115614322565b612ddb8151805160208201207f000000000000000000000000000000000000000000000000000000000000000014614359565b612e236020820151612dee815160ff1690565b6002549160ff831660ff811660ff8416149384612ed4575b612e1d9060209060081c60ff165b93015160ff1690565b93614458565b612e7d612e706080612e3c606085015163ffffffff1690565b93612e6560045495612e518763ffffffff1690565b63ffffffff811663ffffffff8316146144b4565b015163ffffffff1690565b9160201c63ffffffff1690565b63ffffffff811663ffffffff8316146144f5565b612ecc612ec76020608081850151604051612eb381610ada8682019485611975565b51902094015101516001600160401b031690565b614536565b808214612d27565b9350612e1d6020612e14612eeb8286015160ff1690565b60ff612efe60088a901c82165b60ff1690565b9116149692505050612e06565b612f26611a87602060a084015101516001600160401b031690565b5480612f325750505f90565b60408201908151604051612f4e81610ada602082019485611975565b5190201491821592612f6c575b505015612f6757600190565b600290565b6001600160801b03919250612fa2612f936020612fae9301516001600160801b0390511690565b9351516001600160801b031690565b6001600160801b031690565b911610155f80612f5b565b15612fc057565b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b906001600160401b03809116911601906001600160401b038211611dd957565b156130105750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b1561304a5750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156130845750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b90600382029180830460031490151715611dd957565b908160011b9180830460021490151715611dd957565b90602c820291808304602c1490151715611dd957565b908160051b9180830460201490151715611dd957565b15613117575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b905f91602060408201510151518051936101008301926131a8845151613180611c5160e085015161ffff1690565b8091149081613340575b81613330575b81613320575b81613310575b81613300575b50612fb9565b5f5b8681106132e357505f92839283805b8751518710156132a157908992916131e66131e26131dc8a6101808a0151611bdb565b51151590565b1590565b61329657600191613273916132096131ff8b8d51611bdb565b5163ffffffff1690565b809961321e8263ffffffff8116998a10613008565b61327e575b61324d91506020613234888a611bdb565b5101516132468c6101208c0151611bdb565b511461307c565b61326d604061325e859a9789611bdb565b5101516001600160401b031690565b90612fe8565b965b019590916131b9565b63ffffffff61328f92168711613042565b5f88613223565b925095600190613275565b5091509650610c0e9450869193506132de92506132c66001600160401b0382166130b6565b6132d86001600160401b0384166130cc565b1061310e565b6146ac565b916132f960019161326d604061325e8789611bdb565b92016131aa565b90506101a083015151145f6131a2565b610180840151518114915061319c565b6101608401515181149150613196565b6101408401515181149150613190565b610120840151518114915061318a565b9061010081019061336f825151613180611c5160e085015161ffff1690565b61337883613818565b926001600160a01b038416156134c75793926133a360085491611c2a611c258461ffff9060e01c1690565b946133c16133b1878361398e565b9260a01c6001600160401b031690565b955f945f5f9860095497600a549260205f98019b5b8551518910156134a05791898b94928e989796946133ff6131e26131dc8e610180870151611bdb565b613492576134116131ff8d8a51611bdb565b8092613474575b8c91508b90876001998c839e516134309061ffff1690565b9061343a95613ad0565b91909361012001519061344c91611bdb565b5114906134589161307c565b61346191612fe8565b976001905b0197919394959290926133d6565b63ffffffff61348b921663ffffffff821611613042565b5f81613418565b985094505097600190613466565b5050935098505050610c0e94506132de92508691506132c66001600160401b0382166130b6565b633a517eed60e21b5f5260045260245b5ffd5b91906001600160a01b036134ed84613818565b161515806136de575b6136d9576040015160200151518051938415612fc05760b485116136c05761351d85611b95565b925f5b83518110156135645780613553602061353b60019488611bdb565b51015161354d604061325e858a611bdb565b906155e9565b61355d8288611bdb565b5201613520565b5092613600906135df5f9796939661359a89613587613582846130cc565b611dcb565b948361359287611b95565b9c8d92615613565b506135d96135c96135c46135b56135b0856130e2565b613833565b6135be876130f8565b906138a3565b611deb565b976135d3896156b0565b886156f8565b86615786565b6135fa6135f36135ee85615dfe565b615b0e565b86602e0152565b84615796565b5f9460305b83518710156136785761367060019161366b6136218a88611bdb565b5161363363ffffffff8c16848b6156d4565b61365461363f84613841565b60408301516001600160401b0316908b61573e565b602061365f8461384f565b91015190890160200152565b61385d565b960195613605565b95509150925f945b82518610156136b3576136ab6001916136a661369c8987611bdb565b5187830160200152565b61386b565b950194613680565b50935050610c0e91614837565b63156f758160e31b5f52600485905260b460245260445ffd5b505050565b5061ffff60085460e01c1615156134f6565b805f525f60205260ff6137178360405f20906001600160a01b03165f5260205260405f2090565b541661378257805f525f6020526137428260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916600117905533916001600160a01b0316907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260ff6137af8360405f20906001600160a01b03165f5260205260405f2090565b54161561378257805f525f6020526137db8260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916905533916001600160a01b0316907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b6006540361382f576001600160a01b036008541690565b5f90565b6030019081603011611dd957565b9060048201809211611dd957565b90600c8201809211611dd957565b90602c8201809211611dd957565b9060208201809211611dd957565b6001019081600111611dd957565b6024019081602411611dd957565b9060018201809211611dd957565b91908201809211611dd957565b61ffff16602c810290808204602c1490151715611dd95760300180603011611dd95790565b9190823b6001811115613923575f198101908111611dd95781116139075760016138fe82611deb565b9360208501903c565b6001600160a01b0383633610565160e21b5f521660045260245ffd5b6001600160a01b0384633610565160e21b5f521660045260245ffd5b6040519061394c82610b72565b5f6060838281528260208201528260408201520152565b60011b906201fffe61fffe831692168203611dd957565b61ffff5f199116019061ffff8211611dd957565b91909161399961393f565b506030835110613a52576356414c34602084015160e01c03613a5257602483015160c01c92602c81015160f01c93602e82015190604e83015160f01c916139f06139e1610c30565b6001600160401b039093168352565b613a026020830197889061ffff169052565b6040820152613a196060820192839061ffff169052565b94613a26815161ffff1690565b9161ffff8316928315938415613ac5575b508315613a94575b508215613a64575b50509050613a525750565b634724a0fd60e01b5f5260045260245ffd5b6131e29250613a86613a7d613a8c9551935161ffff1690565b915161ffff1690565b91614a2d565b805f80613a47565b90925061ffff613aba611c51613ab5613aaf875161ffff1690565b94613963565b61397a565b91161415915f613a3f565b60b41093505f613a37565b939092959461ffff63ffffffff82169316831015613b9d5761ffff1691613af96135b0826130e2565b9481613b116020888801015160e01c63ffffffff1690565b03613a5257506001901b908116613b7257613b3a613b2e85613841565b84016020015160c01c90565b9516613b555750613b4d610a209261384f565b016020015190565b9050613b6e915061ffff165f52600b60205260405f2090565b5490565b613b98613b8b8361ffff165f52600c60205260405f2090565b546001600160401b031690565b613b3a565b6303e07d4560e61b5f52600485905263ffffffff1660245260445ffd5b9095979198929996613bcb81611dad565b80613dd55750906020899392848015159081613dc8575b613beb91614a7c565b0196613c0a613bf989614aba565b613c03368c610e50565b9088615825565b5f905f5b88868210613d27575b505050613c249350614c61565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016803b156101c557835f8094613c948a956040519c8d97889687957f87d3a9b10000000000000000000000000000000000000000000000000000000087526004870161502e565b03925af1938415612ce057610a2095613cc095613d0d575b5060018111613cd4575b5050503690610e50565b6001600160801b03633b9aca009151160490565b6001600160801b03613cfd613ceb613d0595614aba565b93613cf5876150e8565b9336916150f2565b911691615969565b5f8080613cb6565b80613d1b5f613d2193610bde565b806104a1565b5f613cac565b6131e2613d4f613d48613d42858b613d609699979899614ac4565b80611d78565b3691614ae6565b613d5a368989614ae6565b906158f7565b613dbe575090613d876106bc613d7d613d9d94613c24988c614ac4565b6020810190611d46565b90815181518082149182613da8575b5050614b51565b600189935f88613c17565b9091506020840120906020830120145f80613d96565b9190600101613c0e565b61ffff8111159150613be2565b80613de26134d792611dad565b613deb81611dad565b7f112d89cc000000000000000000000000000000000000000000000000000000005f5260ff16600452602490565b5f80610a2093602081519101845af43d15613e55573d91613e3983610c4e565b92613e476040519485610bde565b83523d5f602085013e6151bd565b6060916151bd565b60016001600160401b03602060408281865151015116940151015116016001600160401b038111611dd9576001600160401b03161490565b613ea66001600160a01b0391613818565b1615613eb457600754600191565b5f905f90565b60405190613ec782610b72565b5f6060838181528260208201528260408201520152565b60405190613eeb82610b72565b606082525f6020830152613efd613eba565b60408301525f60608301528160209060405191613f1b602084610bde565b5f808452805b818110613f2e5750505052565b8290613f38613eba565b82828801015201613f21565b15613f4d575050565b63769a20cb60e01b5f5260045260245260445ffd5b908051518015613fa25760b48111613f8b575090613f82610c0e92615249565b90808214613f44565b63156f758160e31b5f5260045260b460245260445ffd5b82633a517eed60e21b5f5260045260245ffd5b15613a525750565b91908051613fca81613818565b906001600160a01b038216156134c75750926141206001600160401b036140a26140c56140b7610c0e966140006140e29a6152e1565b61400b81835161398e565b906020820192614042614020855161ffff1690565b61ffff614037611c5160085461ffff9060e01c1690565b911614825190613fb5565b614069614056612ef8602084015160ff1690565b801515908161412e575b50825190613fb5565b6140c0600954958695600a549781856140918b8461408a82975161ffff1690565b8b85615314565b9d8f8f908495939551911115613fb5565b886140b18351925161ffff1690565b91615464565b8b808214613f44565b6154cf565b6140dd6140d76135ee889b949b615dfe565b99600955565b600a55565b167fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6008549260a01b16911617600855565b61412984600755565b600655565b6010915011155f614060565b9061414482613818565b6001600160a01b038116156142505761416b90611c2a611c2560085461ffff9060e01c1690565b614175818461398e565b9151916020835191019061418e611c51835161ffff1690565b0361423557600954600a54915f5b855181101561422c576141c56141b4835161ffff1690565b858563ffffffff851692898c613ad0565b60206141d1848a611bdb565b5101511490811591614206575b506141eb5760010161419c565b6134d78763769a20cb60e01b5f52906044916004525f602452565b90506001600160401b0380614220604061325e868c611bdb565b9216911614155f6141de565b50505050505050565b6134d78463769a20cb60e01b5f52906044916004525f602452565b633a517eed60e21b5f52600483905260245ffd5b61427c60409295949395606083526060830190612aa5565b9460208201520152565b612c84915f9161429461222f565b5060405193849283927f893324000000000000000000000000000000000000000000000000000000000084527f00000000000000000000000000000000000000000000000000000000000000009060048501614264565b156142f4575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b1561432b575050565b7f12e74add000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b156143615750565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f60015461439d81611f2c565b908160448501526001811690815f1461443457506001146143d4575b506143d09192600319848303016024850152610837565b0390fd5b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b81831061441a57509192915081016064016143d06143b9565b805483870160640152859450602090920191600101614401565b60ff191660648086019190915291151560051b840190910191506143d090506143b9565b9392919093156144685750505050565b9160ff8092816084969581604051977fd382033a000000000000000000000000000000000000000000000000000000008952166004880152166024860152166044840152166064820152fd5b156144bd575050565b9063ffffffff80927f73b1bde3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b156144fe575050565b9063ffffffff80927f1eb5c1eb000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6001600160401b03165f52600560205260405f205480156145545790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b908160209103126101c55751610a2081610d6a565b979260a0966145f2610a209b976145de61461f986101608e6145d760c09f99986145cc6146109c61ffff6146019c1685526020850190612960565b610120830190612987565b0190612987565b6102406101a08d01526102408c01906103d6565b908a82036101c08c0152610409565b908882036101e08a015261039d565b908682036102008801526129ae565b936102208186039101526001600160401b0381511684526001600160401b0360208201511660208501526040810151604085015263ffffffff6060820151166060850152608081015160808501520151918160a08201520190610837565b1561468457565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b6020604082015151818101516146c981516001600160401b031690565b9161474c60406146eb6146e28786015163ffffffff1690565b63ffffffff1690565b9301518581519101519061473a87806147098563ffffffff90511690565b940151955101519561472b61471c610c3f565b6001600160401b039099168952565b6001600160401b031687890152565b604086015263ffffffff166060850152565b608083015260a082015260e083015161ffff166147c260808501519260a08601519560c08101519061012081015161014082015190610180610160840151930151936040519a8b998a997f4cc22bb7000000000000000000000000000000000000000000000000000000008b5260048b01614591565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af18015612ce057610c0e915f91614808575b5061467d565b61482a915060203d602011614830575b6148228183610bde565b81019061457c565b5f614802565b503d614818565b614841828261398e565b91604051905f602083015261485d82611e626021820184611e2e565b616000825111614a015750610ada6148b76148a561487d845161ffff1690565b60f01b7fffff0000000000000000000000000000000000000000000000000000000000001690565b926040519283916020830195866157a6565b51905ff0906001600160a01b038216156149d9576149c792614910602092614129614977956001600160a01b03167fffffffffffffffffffffffff00000000000000000000000000000000000000006008541617600855565b61491d6040820151600755565b61496e61493182516001600160401b031690565b7fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6008549260a01b16911617600855565b015161ffff1690565b7fffff0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff7dffff000000000000000000000000000000000000000000000000000000006008549260e01b16911617600855565b6149d05f600955565b610c0e5f600a55565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b90614a37906138b0565b90818114928315614a49575b50505090565b90919250621fffe061ffff82169160051b169080820460201490151715611dd9578201809211611dd957145f8080614a43565b15614a845750565b7fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b35610a2081610cff565b9190811015611bef5760051b81013590603e19813603018212156101c5570190565b929190614af281610edf565b93614b006040519586610bde565b602085838152019160051b8101918383116101c55781905b838210614b26575050505050565b81356001600160401b0381116101c557602091614b468784938701610c9f565b815201910190614b18565b91909115614b5d575050565b906143d0614b9f926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610837565b83810360031901602485015290610837565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156101c55701602081359101916001600160401b0382116101c55781360383136101c557565b90602083828152019260208260051b82010193835f925b848410614c295750505050505090565b909192939495602080614c51600193601f19868203018852614c4b8b88614bd1565b90614bb1565b9801940194019294939190614c19565b91909115614c6d575050565b6143d06040519283927ffef760c7000000000000000000000000000000000000000000000000000000008452602060048501526024840191614c02565b9035601e19823603018112156101c55701602081359101916001600160401b0382116101c5578160051b360383136101c557565b9035607e19823603018112156101c5570190565b359060038210156101c557565b9035605e19823603018112156101c5570190565b90602083828152019260208260051b82010193835f925b848410614d3a5750505050505090565b909192939495602080614dac600193601f19868203018852614d5c8b88614cff565b90614d6682614cf2565b614d6f8161022c565b8152614d9e614d93614d8386850185614bd1565b6060888601526060850191614bb1565b926040810190614bd1565b916040818503910152614bb1565b9801940194019294939190614d2a565b610a2091614e86614e7b614dff614de4614dd68680614bd1565b608087526080870191614bb1565b614df16020870187614bd1565b908683036020880152614bb1565b6080614e6b614e116040880188614cde565b8684036040880152614e2281614cf2565b614e2b8161022c565b8452614e3960208201614cf2565b614e428161022c565b6020850152614e5360408201614cf2565b614e5c8161022c565b60408501526060810190614bd1565b9190928160608201520191614bb1565b926060810190614caa565b916060818503910152614d13565b90602083828152019260208260051b82010193835f925b848410614ebb5750505050505090565b909192939495601f198282030184528635601e19843603018112156101c557830190614eeb602082019280614caa565b8091936020845252604082019060408160051b8401019380935f915b838310614f2a575050505050506020806001929801940194019294939190614eab565b909192939495603f19838203018652614f438783614cff565b8035614f4e81610d7f565b614f5781611dad565b8252614f7a614f696020830183614cde565b606060208501526060840190614dbc565b906040810135609e19823603018112156101c55760019360209384936150209301916040818303910152615012614ff2614fc5614fb78580614bd1565b60a0865260a0860191614bb1565b86850135614fd281610d6a565b151587850152614fe56040860186614cde565b8482036040860152614dbc565b92606081013561500181610d6a565b151560608401526080810190614cde565b906080818403910152614dbc565b980196019493019190614f07565b9092809492969593606083019083526060602084015252608081019560808560051b83010194815f90603e1981360301995b83831061507f575050505050610a209495506040818503910152614e94565b9091929397607f1986820301825288358b8112156101c55760206150d960019386839401906150cc6150c26150b48480614caa565b604085526040850191614c02565b9285810190614bd1565b9185818503910152614bb1565b9a019201930191909392615060565b35610a2081610e34565b929190926150ff84610edf565b9361510d6040519586610bde565b602085828152019060051b8201918383116101c55780915b838310615133575050505050565b82356001600160401b0381116101c55782016040818703126101c5576040519161515c83610b37565b81356001600160401b0381116101c557820187601f820112156101c5578781602061518993359101614ae6565b83526020820135926001600160401b0384116101c5576151ae88602095869501610c9f565b83820152815201920191615125565b906151fa57508051156151d257602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b81511580615240575b61520b575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15615203565b8051519081156137825761525c82611b95565b915f5b825180518210156152d357906152c26135ee602061527f84600196611bdb565b5101516152bd6001600160401b03604061529a878b51611bdb565b51015116604051926152ab84610b37565b83526001600160401b03166020830152565b615a3a565b6152cc8287611bdb565b520161525f565b505090505f610a2092615b5e565b90813b6001811115613907575f198101908111611dd95760016138fe82611deb565b906010811015611bef5760051b0190565b9194909293602083019461532f611c58612ef8885160ff1690565b96615349611a496008546001600160401b039060a01c1690565b955f945f5b61535c612ef88b5160ff1690565b811015615458576153746131ff8260408b0151615303565b9663ffffffff88169061538d8961ffff88168410613008565b8261543e575b505082878787878c51946153a695613ad0565b908260808b0151906153b791615303565b516001600160401b03169a8360608c0151906153d291615303565b51926001600160401b038d16926001600160401b0316908484831491821592615433575b50508c5161540391613fb5565b61540c91611dde565b90615416916138a3565b99615420916155e9565b61542a828d611bdb565b5260010161534e565b14159050845f6153f6565b6154519163ffffffff8b51921610613fb5565b5f80615393565b50969750505050505050565b94939190604051956101008701948786106001600160401b03871117610b5257610a20985f9761ffff89989660ff966020968b996040528d52868d015216968760408c01528360608c015260808b01528060a08b01528160c08b01521760e089015201511694615bc0565b9294935f955b6154e6612ef8602087015160ff1690565b8710156155e0576001600160401b039160019161550d611c516131ff8b60408b0151615303565b91600161ffff84161b92836155366155298d60808d0151615303565b516001600160401b031690565b926155458d60608d0151615303565b51936155598c518c8c61ffff881692615d7a565b99166001600160401b038216145f146155a45750901916955b8203615586575050901916965b01956154d5565b61559c9061ffff165f52600b60205260405f2090565b55179661557f565b6155d9906155be8561ffff165f52600c60205260405f2090565b906001600160401b03166001600160401b0319825416179055565b1795615572565b95509392505050565b6135ee906001600160401b03610a20936040519261560684610b37565b8352166020820152615a3a565b9190949394808203828111611dd957600181146156915761563390615e36565b9283820192838311611dd9576001880192838911611dd9578387866156589386615613565b615667613582600192976130cc565b890101809311611dd95761567c938692615613565b6156899061200d92615f70565b938492611bdb565b506156a391506156ac92959495611bdb565b51928392611bdb565b5290565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b602d908260081c602c8201530153565b604f908260081c604e8201530153565b600a907fffff000000000000000000000000000000000000000000000000000000000000610a2094937f610000000000000000000000000000000000000000000000000000000000000083521660018201527f80600a3d393df30000000000000000000000000000000000000000000000000060038201520190611e2e565b9061587490612ecc604051602081019061585c8288604080916001600160801b038151168452602081015160208501520151910152565b6060815261586b608082610bde565b51902091614536565b60208201518082036158c95750506001600160801b03633b9aca00915116046158a18142428211156142eb565b4203428111611dd9576001600160801b03610c0e911663ffffffff6004541690818110615e52565b7f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b908151815103613782575f5b8251811015615961576159168184611bdb565b51516159228284611bdb565b51510361595a576159338184611bdb565b51602081519101206159458284611bdb565b51602081519101200361595a57600101615903565b5050505f90565b505050600190565b929190925f5b8451811015615a33576159828186611bdb565b51908360405160208101906001600160401b038616825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801915f5b8181106159fc575050505081610ada60019760206159f2940151605f19848303016080850152610837565b5190205d0161596f565b919394965091949697602080615a1e60019360bf198b82030188528951610837565b970194019101918a96949392989795986159c7565b5050509050565b602460208201906001600160401b0382511680615ad3575b50615a5f615a8d91611deb565b92615a84615a7e615a78615a7287615e96565b87615ea2565b86615eb9565b85615ed0565b90519084615efd565b90615aa2611a4982516001600160401b031690565b615aab57505090565b615acc611a49615abe612cc09486615ee6565b92516001600160401b031690565b9083615f15565b600191505b6080811015615b005750615a5f615af9615af4615a8d93613879565b613887565b9150615a52565b60019060071c910190615ad8565b805160018101809111611dd957615b2490611deb565b805115611bef57615b4e81615b416020945f868196015382615f4e565b5060405191828092611e2e565b039060025afa15612ce0575f5190565b91818103818111611dd95760018114615ba357615b7a90615e36565b820191828111611dd95782615b8f9185615b5e565b91615b9a9293615b5e565b610a2091615f70565b5050615bae91611bdb565b5190565b5f198114611dd95760010190565b9392909491828414808091615d62575b615d3a5761ffff8311612fc057615be78783611dde565b9060018214615c905750615bfa90615e36565b93615c27615c0886896138a3565b956135be613582615c21615c1b88613895565b97613895565b926130cc565b92815b85811087828a83615c69575b50505015615c4c57615c4790615bb2565b615c2a565b9495969791859188615c5e948b615bc0565b94615b9a9596615bc0565b63ffffffff929350615c8691604060606131ff9301510151615303565b161087828a615c36565b925050509493929415615cdc57505082615cd791610a20939451916020810151615cbf604083015161ffff1690565b9063ffffffff60c060a0850151940151941694613ad0565b6155e9565b615ce7829392613895565b1490811591615d18575b50615d04576080615bae92930151611bdb565b8251634724a0fd60e01b5f5260045260245ffd5b9050615d326146e26131ff84604060608901510151615303565b14155f615cf1565b505092915050610a209250805190615d5c6040602083015192015161ffff1690565b91615fe4565b50615d756131e2838960e08a0151615fc2565b615bd0565b9091602061ffff919594950151169363ffffffff811694851015615de05750615da56135b0856130e2565b936020858401015160e01c03613a525750602090615dda615dd4615dc886613841565b83016020015160c01c90565b9461384f565b01015190565b6303e07d4560e61b5f5260049190915263ffffffff1660245260445ffd5b60405190606090615e0f8284610bde565b60228352612cc091600a906020850190601f19013682375360206021840153600283615efd565b9060015b8060011b9083821015615e4d5750615e3a565b925050565b15615e5b575050565b906001600160801b0380927fdb020436000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6020600a910153600190565b60208260229201015360018101809111611dd95790565b602082600a9201015360018101809111611dd95790565b602082819201015360018101809111611dd95790565b60208260109201015360018101809111611dd95790565b8160209193929301015260208101809111611dd95790565b9092919083016020015b6080821015615f3357906001929391530190565b600180916080607f85161781530193019060071c9092615f1f565b908051918215615961576021602084930191015e60010180600111611dd95790565b604051615f7e608082610bde565b60418152615f8c6041610c4e565b602082019190601f1901368337805115611bef576020935f936001615b4e94536021830152604182015260405191828092611e2e565b91818103908111611dd9576001901b905f198201918211611dd9571b16151590565b909161ffff1692602c840293808504602c1490151715611dd957836030019384603011611dd9578160051b9180830460201490151715611dd957019260308401809111611dd9576160349061386b565b825110613a52575001605001519056fea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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
