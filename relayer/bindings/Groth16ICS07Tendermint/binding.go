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
	Bin: "0x610160806040523461056657616888803803809161001d8285610688565b8339810160e08282031261056657610034826106ab565b610040602084016106ab565b61004c604085016106ab565b91610059606086016106ab565b60808601519094906001600160401b0381116105665786019080601f8301121561056657815161008b926020016106bf565b9461009d60c060a083015192016106ab565b957f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a612336101005280518101906020820190602081840312610566576020810151906001600160401b038211610566570191829003601f1981019061012013610566576040519160e083016001600160401b038111848210176106595760405260208401516001600160401b038111610566576020908501019080601f8301121561056657815161014e926020016106bf565b8252604081126105665760408051916101668361066d565b610171828601610704565b835261017f60608601610704565b602084015260208401928352603f190112610566576040516101a08161066d565b6101ac60808501610712565b81526101ba60a08501610712565b6020820152604083019081526101d260c08501610726565b90606084019182526101e660e08601610726565b92608085019384526101008601519586151587036105665760a0860196875261012001519460028610156105665760c08101958652518051906001600160401b03821161065957610238600154610737565b601f8111610609575b50602090601f831160011461059c5763ffffffff95949392915f9183610591575b50508160011b915f199060031b1c1916176001555b5160ff81511661ff00602060025493015160081b169161ffff191617176002555160018060401b0381511660035491602068010000000000000000600160801b0391015160401b169160018060801b031916171760035551169267ffffffff00000000600454925160201b169051151560401b925191600283101561057d5769ff00000000000000000068ff00000000000000009360481b169469ff000000000000000000199160018060481b03191617161791161717600455604051610348816103418161076f565b0382610688565b602081519101206101205260405163685e272760e11b815260206004820152602081806103776024820161076f565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4908115610572575f9161053c575b50610140526040516103b3816103418161076f565b6001600160401b03906020906103c890610820565b0151600354916001600160401b0383169116818103610527575050604090811c6001600160401b03165f908152600560205220556001600160a01b0390811660805290811660a05290811660c0521660e05260045463ffffffff8082169061070882019081116105135763ffffffff809360201c1692839116116104fe57826001600160a01b0381166104d7575061045e610af7565b5061046b61010051610b79565b505b604051615b809081610c48823960805181615233015260a051816146ae015260c05181505060e0518181816122f701528181612cae0152612d4901526101005181818161026f015261030a0152610120518161236e0152610140518181816122c20152612c790152f35b806104e46104f892610a01565b506104ee81610a77565b5061010051610bd2565b5061046d565b6333fae18560e21b5f5260045260245260445ffd5b634e487b7160e01b5f52601160045260245ffd5b6355bace6f60e01b5f5260045260245260445ffd5b90506020813d60201161056a575b8161055760209383610688565b8101031261056657515f61039e565b5f80fd5b3d915061054a565b6040513d5f823e3d90fd5b634e487b7160e01b5f52602160045260245ffd5b015190505f80610262565b90601f1983169160015f52815f20925f5b8181106105f1575091600193918563ffffffff9998979694106105d9575b505050811b01600155610277565b01515f1960f88460031b161c191690555f80806105cb565b929360206001819287860151815501950193016105ad565b60015f525f5160206168485f395f51905f52601f840160051c8101916020851061064f575b601f0160051c01905b8181106106445750610241565b5f8155600101610637565b909150819061062e565b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b0382111761065957604052565b601f909101601f19168101906001600160401b0382119082101761065957604052565b51906001600160a01b038216820361056657565b9192916001600160401b03821161065957604051916106e8601f8201601f191660200184610688565b829481845281830111610566578281602093845f96015e010152565b519060ff8216820361056657565b51906001600160401b038216820361056657565b519063ffffffff8216820361056657565b90600182811c92168015610765575b602083101461075157565b634e487b7160e01b5f52602260045260245ffd5b91607f1691610746565b6001545f929161077e82610737565b80825291600181169081156107df5750600114610799575050565b60015f9081529293509091905f5160206168485f395f51905f525b8383106107c5575060209250010190565b6001816020929493945483858701015201910191906107b4565b9050602093945060ff929192191683830152151560051b010190565b90815181101561080c570160200190565b634e487b7160e01b5f52603260045260245ffd5b60405161082c8161066d565b606081525f602082015250805180159081156109f5575b506109e6575f19908051805b6109a2575b505f19821461099457600182019081831161051357600360fc1b6001600160f81b031961088184846107fb565b5116148061097f575b61096f575f5b815183101561091f576108a383836107fb565b5160f81c603081108015610915575b61090357600a82026001600160401b03908116602f1990920160ff16919091018116911681106108e757600190920191610890565b50915050604051906108f88261066d565b81525f602082015290565b5050915050604051906108f88261066d565b50603981116108b2565b9150916001811190811591610963575b5061095457604051916109418361066d565b82526001600160401b0316602082015290565b63384c687760e21b5f5260045ffd5b602b915010155f61092f565b915050604051906108f88261066d565b5080518381039081116105135760021061088a565b60405191506108f88261066d565b5f19810181811161051357602d60f81b6001600160f81b03196109c583866107fb565b5116146109dc57508015610513575f19018061084f565b92505f9050610854565b6329120bff60e21b5f5260045ffd5b6040915010155f610843565b6001600160a01b0381165f9081525f5160206168685f395f51905f52602052604090205460ff16610a72576001600160a01b03165f8181525f5160206168685f395f51905f5260205260408120805460ff191660011790553391905f5160206167c85f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f5160206167e85f395f51905f52602052604090205460ff16610a72576001600160a01b03165f8181525f5160206167e85f395f51905f5260205260408120805460ff191660011790553391905f5160206168285f395f51905f52905f5160206167c85f395f51905f529080a4600190565b5f80525f5160206167e85f395f51905f526020525f5160206168085f395f51905f525460ff16610b75575f8080525f5160206167e85f395f51905f526020525f5160206168085f395f51905f52805460ff1916600117905533905f5160206168285f395f51905f525f5160206167c85f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f14610a7257805f525f60205260405f205f805260205260405f20600160ff198254161790555f33915f5160206167c85f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16610c41575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291905f5160206167c85f395f51905f529080a4600190565b50505f9056fe610140806040526004361015610013575f80fd5b5f3560e01c90816301ffc9a714610a37575080630bece3561461099d578063248a9ca3146109735780632f2ff15d1461094457806336568abe146108f5578063536c2ad31461089d5780636a28f000146107f85780638a8e4c5d146107d957806391d148541461079d578063974a74c414610652578063a217fddf14610638578063a6f031bb14610528578063ac9650d814610363578063d547741f1461032d578063db3e1fa4146102f3578063ddba6537146102505763ef913a4b146100d8575f80fd5b3461024c575f36600319011261024c5760405160208082015261012060408201525f6001546101068161315b565b90816101608501526001811690815f1461022757506001146101c6575b6101c2836101ae818560ff600254818116606085015260081c1660808301526001600160401b0360035481811660a085015260401c1660c083015260ff60045463ffffffff811660e085015263ffffffff8160201c16610100850152818160401c16151561012085015260481c1661019a81611318565b61014083015203601f198101835282610ccf565b604051918291602083526020830190610c10565b0390f35b91905060015f527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6915f905b80821061020b57509091508101610180016101ae610123565b91926001816020925461018085880101520191019092916101f2565b60ff19166101808086019190915291151560051b840190910191506101ae9050610123565b5f80fd5b3461024c5761025e36610ad5565b505060ff60045460401c166102cb577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f205416156102bc575b636d40ebe160e11b5f5260045ffd5b6102c590613202565b806102ad565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b3461024c575f36600319011261024c5760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b3461024c5761036161033e36610b42565b9061035c610357825f525f602052600160405f20015490565b613202565b61425a565b005b3461024c57602036600319011261024c576004356001600160401b03811161024c573660238201121561024c5780600401356001600160401b03811161024c5760248201913660248360051b8301011161024c576020926040516103c78582610ccf565b5f815284810191601f1986013684376103df85610e74565b936103ed6040519586610ccf565b858552601f196103fc87610e74565b01875f5b828110610519575050505f5b868110156104bc576001906104985f808b8861046561043360248860051b8b01018b6130b1565b90938d6040519483869484860198893784019083820190898252519283915e010185815203601f198101835282610ccf565b5190305af43d156104b4573d9061047b82610cf0565b916104896040519384610ccf565b82523d5f8d84013e5b30615635565b6104a28289612f48565b526104ad8188612f48565b500161040c565b606090610492565b604080518981528751818b018190525f92600582901b83018101918a8d01918d9085015b8287106104ed5785850386f35b909192938280610509600193603f198a82030186528851610c10565b96019201960195929190926104e0565b60608882018301528101610400565b3461024c57602036600319011261024c576004356001600160401b03811161024c578060040190610140600319823603011261024c5760ff60045460401c166102cb575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff161561062b575b6105ca60448201836130e3565b916105d860648201856130e3565b93909161010481013590600282101561024c57602096610623966106006101248401836130e3565b969095604051986106118c8b610ccf565b5f8a52608460a48701960135946145b2565b604051908152f35b610633613193565b6105bd565b3461024c575f36600319011261024c5760206040515f8152f35b3461024c57602036600319011261024c576004356001600160401b03811161024c578060040190610160600319823603011261024c5760ff60045460401c166102cb575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff1615610790575b61014481016106f681846130b1565b9050156107685761070a60448301846130e3565b909161071960648501866130e3565b95909461010481013591600283101561024c5760209761062397610751966107586107486101248701866130e3565b999098866130b1565b3691610d0b565b98608460a48701960135946145b2565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b610798613193565b6106e7565b3461024c576107ab36610b42565b905f525f6020526001600160a01b0360405f2091165f52602052602060ff60405f2054166040519015158152f35b3461024c576107e736610ad5565b5050636d40ebe160e11b5f5260045ffd5b3461024c575f36600319011261024c57335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff16156108865760045460ff8160401c161561085e5768ff00000000000000001916600455005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b3461024c57602036600319011261024c576108d96108e76101c26108c2600435612f70565b919390604051958695606087526060870190610b68565b908582036020870152610ba1565b908382036040850152610bd4565b3461024c5761090336610b42565b336001600160a01b0382160361091c576103619161425a565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b3461024c5761036161095536610b42565b9061096e610357825f525f602052600160405f20015490565b6141cd565b3461024c57602036600319011261024c5760206106236004355f525f602052600160405f20015490565b3461024c576109ab36610ad5565b9060ff60045460401c166102cb575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060209081527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47549092610a1992909160ff1615610a2a57611a3e565b60405190610a2681610b24565b8152f35b610a32613193565b611a3e565b3461024c57602036600319011261024c57600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361024c57817f7965db0b0000000000000000000000000000000000000000000000000000000060209314908115610aab575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483610aa4565b90602060031983011261024c576004356001600160401b03811161024c578260238201121561024c578060040135926001600160401b03841161024c576024848301011161024c576024019190565b60031115610b2e57565b634e487b7160e01b5f52602160045260245ffd5b604090600319011261024c57600435906024356001600160a01b038116810361024c5790565b90602080835192838152019201905f5b818110610b855750505090565b825163ffffffff16845260209384019390920191600101610b78565b90602080835192838152019201905f5b818110610bbe5750505090565b8251845260209384019390920191600101610bb1565b90602080835192838152019201905f5b818110610bf15750505090565b82516001600160401b0316845260209384019390920191600101610be4565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b604081019081106001600160401b03821117610c4f57604052565b634e487b7160e01b5f52604160045260245ffd5b606081019081106001600160401b03821117610c4f57604052565b608081019081106001600160401b03821117610c4f57604052565b60c081019081106001600160401b03821117610c4f57604052565b60e081019081106001600160401b03821117610c4f57604052565b90601f801991011681019081106001600160401b03821117610c4f57604052565b6001600160401b038111610c4f57601f01601f191660200190565b929192610d1782610cf0565b91610d256040519384610ccf565b82948184528183011161024c578281602093845f960137010152565b9080601f8301121561024c57816020610d5c93359101610d0b565b90565b359060ff8216820361024c57565b35906001600160401b038216820361024c57565b919082604091031261024c57604051610d9981610c34565b6020610db2818395610daa81610d6d565b855201610d6d565b910152565b359063ffffffff8216820361024c57565b3590811515820361024c57565b35906001600160801b038216820361024c57565b919082606091031261024c57604051610e0181610c63565b6040808294610e0f81610dd5565b8452602081013560208501520135910152565b80929103916060831261024c57604051610e3b81610c34565b6040819483358352601f19011261024c576020906040805193610e5d85610c34565b610e68848201610db7565b85520135828401520152565b6001600160401b038111610c4f5760051b60200190565b919060808382031261024c5760405190610ea482610c7e565b819380356001600160401b03811161024c57606092610ec4918301610d41565b835260208101356020840152610edc60408201610d6d565b60408401520135908160070b820361024c5760600152565b91909160808184031261024c5760405190610f0e82610c7e565b819381356001600160401b03811161024c57820181601f8201121561024c578035610f3881610e74565b91610f466040519384610ccf565b81835260208084019260051b8201019184831161024c5760208201905b838210610fb457505050508352610f7c60208301610dc8565b60208401526040820135906001600160401b03821161024c5782610fa960609492610db294869401610e8b565b604086015201610d6d565b81356001600160401b03811161024c57602091610fd688848094880101610e8b565b815201910190610f63565b9080601f8301121561024c57604080519290610ffd9084610ccf565b82906040810192831161024c57905b8282106110195750505090565b813581526020918201910161100c565b9080601f8301121561024c57813561104081610e74565b9261104e6040519485610ccf565b81845260208085019260051b82010192831161024c57602001905b8282106110765750505090565b6020809161108384610db7565b815201910190611069565b6040519061109b82610c63565b5f6040838281528260208201520152565b519060ff8216820361024c57565b51906001600160401b038216820361024c57565b919082604091031261024c576040516110e681610c34565b6020610db28183956110f7816110ba565b8552016110ba565b519063ffffffff8216820361024c57565b5190811515820361024c57565b51906001600160801b038216820361024c57565b919082606091031261024c5760405161114981610c63565b60408082946111578161111d565b8452602081015160208501520151910152565b60208183031261024c578051906001600160401b03821161024c57016101808183031261024c576040519161119e83610c99565b81516001600160401b03811161024c5782019182820392610120841261024c57604051936111cb85610cb4565b81516001600160401b03811161024c5782019084601f8301121561024c5781516111f481610cf0565b906112026040519283610ccf565b808252866020828601011161024c576020815f9282604097018386015e830101528652601f19011261024c576101009060405161123e81610c34565b61124a602083016110ac565b8152611258604083016110ac565b6020820152602086015261126f84606083016110ce565b604086015261128060a082016110ff565b606086015261129160c082016110ff565b60808601526112a260e08201611110565b60a0860152015190600282101561024c57836101409260c061131096015285526112cf8360208301611131565b60208601526112e18360808301611131565b60408601526112f260e0820161111d565b60608601526113058361010083016110ce565b6080860152016110ce565b60a082015290565b60021115610b2e57565b906060806113398451608085526080850190610c10565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b8282106113cb57505050506001600160401b0360606113c1819360208701511515602087015260408701518682036040880152611322565b9401511691015290565b909192936020806113e8600193609f198a82030186528851611322565b960192019201909291611389565b905f905b6002821061140757505050565b60208060019285518152019301910190916113fa565b90602080835192838152019201905f5b81811061143a5750505090565b8251151584526020938401939092019160010161142d565b908151610960825260c06114758251610120610960860152610a80850190610c10565b9160ff602080830151828151166109808801520151166109a08501526114b960408201516109c08601906001600160401b0360208092828151168552015116910152565b63ffffffff606082015116610a0085015263ffffffff608082015116610a2085015260a08101511515610a4085015201516114f381611318565b610a6083015261152860208401516020840190604080916001600160801b038151168452602081015160208501520151910152565b6040830151908281036080840152815160a0825260206116928251604060a086015261156d60e0860182516001600160401b0360208092828151168552015116910152565b61024061158b848301516102c06101208901526103a0880190610c10565b60408301516001600160401b031661014088015260608301516001600160801b03166101608801526080830151151561018088015260a083015180516101a0890152602090810151805163ffffffff166101c08a015201516101e08801529160c0810151151561020088015260e08101516102208801526101008101511515828801526101208101516102608801526101408101516102808801526101608101516102a08801526101808101516102c08801526101a08101516102e08801526101c081015115156103008801526101e08101516103208801526102008101511515610340880152610220810151610360880152015160df1986830301610380870152610c10565b91015190609f198382030160c0840152606060c08201926001600160401b03815116835263ffffffff60208201511660208401526116f26040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519160c060a0830152825180915260e0820190602060e08260051b8501019401925f905b828210611986575050505050906060611740849360206117789601518482036020860152611362565b92611768604082015160408501906001600160401b0360208092828151168552015116910152565b0151906080818403910152611362565b60608301516001600160801b031660a083015260808301515f60c084015b60088210611970575050506117ee906117b860a08501516101c08501906113f6565b6117cb60c08501516102008501906113f6565b61ffff60e085015116610240840152610100840151838203610260850152610b68565b61012083015190828103610280840152602080835192838152019201905f5b81811061195a57505050610140830151908281036102a0840152602080835192838152019201905f5b81811061193b575050506118896118756118616101c0936101608701518682036102c0880152610b68565b6101808601518582036102e087015261141d565b6101a0850151848203610300860152610b68565b9201518051610320830152602081015160ff1661034083015260408101515f61036084015b6010821061191f5750505060608101515f61056084015b601082106119095750505060800151905f90610760015b601082106118ea5750505090565b6020806001926001600160401b038651168152019301910190916118dc565b60208060019285518152019301910190916118c5565b60208060019263ffffffff8651168152019301910190916118ae565b82516001600160401b0316845260209384019390920191600101611836565b825184526020938401939092019160010161180d565b6020806001928551815201930191019091611796565b909192939460df1982820301855285519081516004811015610b2e57611a0482602060019581959482955201519060408482015260606119d283516080604085015260c0840190610c10565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610c10565b9701950193920190611717565b6040513d5f823e3d90fd5b611a3460409295949395606083526060830190611452565b9460208201520152565b9081019060208183031261024c578035906001600160401b03821161024c5701916109608383031261024c576040516101e081018181106001600160401b03821117610c4f5760405283356001600160401b03811161024c57840180840390610120821261024c5760405191611ab383610cb4565b8135906001600160401b03821161024c57611ad2876040938501610d41565b8452601f19011261024c5761010090604051611aed81610c34565b611af960208301610d5f565b8152611b0760408301610d5f565b60208201526020840152611b1e8660608301610d81565b6040840152611b2f60a08201610db7565b6060840152611b4060c08201610db7565b6080840152611b5160e08201610dc8565b60a08401520135600281101561024c5760c08201528152611b758360208601610de9565b602082015260808401356001600160401b03811161024c5784019360a08585031261024c5760405194611ba786610c7e565b80356001600160401b03811161024c57810160408187031261024c5760405190611bd082610c34565b80356001600160401b03811161024c5781016102c08189031261024c576040519061026082018281106001600160401b03821117610c4f57604052611c158982610d81565b825260408101356001600160401b03811161024c5789611c36918301610d41565b6020830152611c4760608201610d6d565b6040830152611c5860808201610dd5565b6060830152611c6960a08201610dc8565b6080830152611c7b8960c08301610e22565b60a0830152611c8d6101208201610dc8565b60c083015261014081013560e0830152611caa6101608201610dc8565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a0830152611cf96102208201610dc8565b6101c08301526102408101356101e0830152611d186102608201610dc8565b6102008301526102808101356102208301526102a0810135906001600160401b03821161024c57611d4b918a9101610d41565b61024082015282526020810135906001600160401b03821161024c570160c08188031261024c5760405190611d7f82610c7e565b611d8881610d6d565b8252611d9660208201610db7565b6020830152611da88860408301610e22565b604083015260a0810135906001600160401b03821161024c570187601f8201121561024c57803590611dd982610e74565b91611de76040519384610ccf565b80835260208084019160051b830101918a831161024c5760208101915b838310612e2d575050505060608201526020820152865260208101356001600160401b03811161024c5785611e3a918301610ef4565b6020870152611e4c8560408301610d81565b60408701526080810135906001600160401b03821161024c57611e7191869101610ef4565b606086015260408201948552611e8960a08201610dd5565b60608301528360df8201121561024c57604051611ea861010082610ccf565b806101c083019186831161024c5790869160c08501905b848210612e1a5750506080850152611ed691610fe1565b60a0830152611ee9846102008301610fe1565b60c08301526102408101359361ffff8516850361024c5760e083019485526102608201356001600160401b03811161024c5781611f27918401611029565b6101008401526102808201356001600160401b03811161024c57820181601f8201121561024c57803590611f5a82610e74565b91611f686040519384610ccf565b80835260208084019160051b8301019184831161024c57602001905b828210612e0a575050506101208401526102a08201356001600160401b03811161024c57820181601f8201121561024c57803590611fc182610e74565b91611fcf6040519384610ccf565b80835260208084019160051b8301019184831161024c57602001905b828210612df25750505061014084019081526102c08301356001600160401b03811161024c578261201d918501611029565b9161016085019283526102e08401356001600160401b03811161024c57840181601f8201121561024c5780359061205382610e74565b916120616040519384610ccf565b80835260208084019160051b8301019184831161024c57602001905b828210612dda575050506101808601526103008401356001600160401b03811161024c57816120ad918601611029565b6101a086019081529361064081830361031f19011261024c576040519160a083018381106001600160401b03821117610c4f5760405261032082013583526120f86103408301610d5f565b60208401528061037f8301121561024c576102009160405161211a8482610ccf565b8061056083019184831161024c576103608401905b838210612dc257505060408601528261057f8301121561024c576040516121568582610ccf565b8061076084019285841161024c57905b838210612db257505060608601528261077f8301121561024c5761218d6040519485610ccf565b61096084920192831161024c57905b828210612d9a5750505060808201526101c08501526121ba84613242565b9892959194909460c05196996040516121d281610c99565b6040516121de81610cb4565b606081526040516121ee81610c34565b5f81525f6020820152602082015260405161220881610c34565b5f81525f602082015260408201525f60608201525f60808201525f60a08201525f60c0820152815261223861108e565b602082015261224561108e565b60408201525f606082015260405161225c81610c34565b5f81525f6020820152608082015260405161227681610c34565b5f81525f602082015260a0820152508a5f14612cfe5715612c48575f6122eb91604051809381927ffe0a61f30000000000000000000000000000000000000000000000000000000083527f00000000000000000000000000000000000000000000000000000000000000008d60048501611a1c565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115611a11575f91612c26575b505b995b8a633b9aca006001600160801b03606083519301511604612350814242821115614fe7565b61070861235d8242613126565b11612bf757508051805160208201207f000000000000000000000000000000000000000000000000000000000000000003612b005750602081015160ff815116906002549160ff8316928382149283612ae7575b6020015160ff169215612a9f575050505063ffffffff606082015116906004549163ffffffff8316808203612a7157505063ffffffff608081920151169160201c16808203612a435750506124648b61245c6001600160401b0360206080818501516040516124408482018093604080916001600160801b038151168452602081015160208501520151910152565b6060815261244e8382610ccf565b51902094015101511661501e565b808214613a09565b61246d8b613a40565b9915612830575061ffff610100880151519351168093149384612820575b84612814575b5083612808575b50826127f8575b826127ec575b5050156127c4576124b5826142dd565b946001600160a01b038616156127b15793906124e7600897939754966124e161ffff8960e01c16614313565b90614338565b966124f288826143a2565b6080525f925f965f9560095494600a54965f985b6101008b0151518a10156125d057908d979695949392918b61252d8c610180830151612f48565b51156125c4576125478c61010063ffffffff930151612f48565b5116809d6125ac575b50508a60019c8b88828d8c829e6080516020015161ffff1690612572956144ea565b91909361012001519061258491612f48565b51149061259091613b82565b61259991613aee565b986001905b019890919293949596612506565b63ffffffff6125bd92168111613b48565b5f8c612550565b5097509860019061259e565b50955096509750979195935097506001600160401b03821660038102908082046003149015171561279d576801fffffffffffffffe82609f1c16906001600160401b038360a01c1682046002146001600160401b038460a01c1615171561279d576001600160401b039361264b92858560a01c169211613bd2565b61265483615064565b60a01c16915b61278c575b50505061266b82610b24565b81612743576001600160401b036020604060a08401938451838101519085600354851c16868316116126ef575b505001516040516126c98382018093604080916001600160801b038151168452602081015160208501520151910152565b606081526126d8608082610ccf565b51902092510151165f52600560205260405f205590565b85905116851960035416176003557fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6fffffffffffffffff000000000000000060035492851b169116176003555f80612698565b5061274d81610b24565b60018103612775576801000000000000000068ff000000000000000019600454161760045590565b61277e81610b24565b60028103610d5c5750600290565b61279592613c16565b5f808061265f565b634e487b7160e01b5f52601160045260245ffd5b82633a517eed60e21b5f5260045260245ffd5b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b51511490505f806124a5565b610180860151518214925061249f565b5151821492505f612498565b5151831493505f612491565b610120880151518414945061248b565b9891939097925f9695965060205f9a510151519485519961ffff6101008b0151519351168093149384612a33575b84612a27575b5083612a1b575b5082612a0b575b826129ff575b5050156127c4575f9893985b8781106129d757505f935f995f965f9b5b6101008a0151518d1015612969576128b28d6101808c0151612f48565b511561295f578a8a888f80610100840151906128cd91612f48565b5163ffffffff169c8d80968180978110906128e791613b0e565b61291897612911956101209460209461290694612947575b5050612f48565b510151930151612f48565b5114613b82565b600161293d81986001600160401b0360406129338d8c612f48565b5101511690613aee565b9c5b019b96612895565b63ffffffff61295892168111613b48565b5f826128ff565b969b60019061293f565b509650969492995096509691506001600160401b03811660038102908082046003149015171561279d576001600160401b038316916801fffffffffffffffe8460011b16928084046002149015171561279d576129c892849211613bd2565b6129d182615064565b9161265a565b976129f56001916001600160401b0360406129338d999e9989612f48565b9801989398612884565b51511490505f80612878565b6101808901515182149250612872565b5151821492505f61286b565b5151831493505f612864565b6101208b0151518414945061285e565b7f1eb5c1eb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f73b1bde3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6084945060ff90604051947fd382033a000000000000000000000000000000000000000000000000000000008652600486015260081c16602484015260448301526064820152fd5b6020810151600883901c60ff90811691161493506123b1565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f600154612b3c8161315b565b908160448501526001811690815f14612bd35750600114612b73575b50612b6f9192600319848303016024850152610c10565b0390fd5b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b818310612bb95750919291508101606401612b6f612b58565b805460648488010152859450602090920191600101612ba0565b60ff191660648086019190915291151560051b84019091019150612b6f9050612b58565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b612c4291503d805f833e612c3a8183610ccf565b81019061116a565b5f612327565b5f612ca291604051809381927f893324000000000000000000000000000000000000000000000000000000000083527f00000000000000000000000000000000000000000000000000000000000000008d60048501611a1c565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115611a11575f91612ce4575b50612329565b612cf891503d805f833e612c3a8183610ccf565b5f612cde565b50506040517ffe284b13000000000000000000000000000000000000000000000000000000008152602060048201525f8180612d3d602482018c611452565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115611a11575f91612d80575b509961232b565b612d9491503d805f833e612c3a8183610ccf565b5f612d79565b60208091612da784610d6d565b81520191019061219c565b8135815260209182019101612166565b60208091612dcf84610db7565b81520191019061212f565b60208091612de784610dc8565b81520191019061207d565b60208091612dff84610d6d565b815201910190611feb565b8135815260209182019101611f84565b8135815288935060209182019101611ebf565b82356001600160401b03811161024c578201906040828e03601f19011261024c5760405190612e5b82610c34565b6020830135600481101561024c57825260408301356001600160401b03811161024c5760208f949160809201018094031261024c5760405191612e9d83610c7e565b83356001600160401b03811161024c578f612eb9918601610d41565b8352612ec760208501610dd5565b6020840152612ed860408501610dc8565b60408401526060840135926001600160401b03841161024c578f612f029060209695879601610d41565b606082015283820152815201920191611e04565b90612f2082610e74565b612f2d6040519182610ccf565b8281528092612f3e601f1991610e74565b0190602036910137565b8051821015612f5c5760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b90612f7a826142dd565b6001600160a01b0381161561306757612f9f906124e161ffff60085460e01c16614313565b6020612fab82856143a2565b0191612fbb61ffff845116612f16565b93612fca61ffff855116612f16565b938493612fdb61ffff835116612f16565b9260095491600a54935f5b61ffff8251168b61ffff83169182101561305b57996001600160401b0361304783806130408c9d9e9f828d9e8d9e8d9e8d9e61ffff9e8f9060019f9e61302f816130389a612f48565b525116916144ea565b929096612f48565b528c612f48565b911690520116908897969594939291612fe6565b50505050505050509150565b50905060206040519161307a8284610ccf565b5f83525f3681376040519261308f8385610ccf565b5f84525f368137604051926130a48185610ccf565b5f8452505f368137929190565b903590601e198136030182121561024c57018035906001600160401b03821161024c5760200191813603831361024c57565b903590601e198136030182121561024c57018035906001600160401b03821161024c57602001918160051b3603831361024c57565b5f1981019190821161279d57565b9190820391821161279d57565b9061313d82610cf0565b61314a6040519182610ccf565b8281528092612f3e601f1991610cf0565b90600182811c92168015613189575b602083101461317557565b634e487b7160e01b5f52602260045260245ffd5b91607f169161316a565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff16156131cb57565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260405f206001600160a01b0333165f5260205260ff60405f2054161561322c5750565b63e2517d3f60e01b5f523360045260245260445ffd5b5f60c0525f905f92604082015160016001600160401b03602060408451946101408651015160c052838280858b01510151975101511660a052015101511601926001600160401b03841161279d5761329b60c051614da9565b816139e7576101c08301805151909290156139bd575050518051906132bf826142dd565b916001600160a01b0383169081156139ab5750823b90600182111561399957505f19810190811161279d5760016132f582613133565b9360208501903c60206133098383516143a2565b0160e05261ffff60e05151169161332e6008549361ffff8560e01c1684519114614fce565b60ff6020830151168015158061398e575b835161334a91614fce565b6001600160401b0360095494600a546101005260a01c169261336b82612f16565b5f5f5b84811061380d57505061338c82516001600160401b03871115614fce565b81519061ffff60e0515116916040519061010082018281106001600160401b03821117610c4f5761341394613408945f93849360409b9a999b52855288602086015281604086015289606086015260808501528a60a08501526101005160c0850152610100518b1760e08501528160ff60208b01511694615862565b60c051818114614e3a565b5f925b8184106136b55750505050604051916338f49afb60e01b835260c05160048401526020836024817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4928315611a11575f93613681575b5060095561010051600a557fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b806008549360a01b16169116176008558060075560c051600655956001600160401b036001956134c8614dce565b6020604086015101525b1660a051146136635784158080613658575b15613505575050604091500151606060208201519101526001915b93929190565b156135245750604001516060015161351c91614e58565b6001916134ff565b939160c05114613535575b506134ff565b60609192959493506040015101519161354f60c0516142dd565b6001600160a01b0381161561364357613574906124e161ffff60085460e01c16614313565b926135818460c0516143a2565b9051946020865192019161ffff835116036135ef5760095492600a54945f5b8851811015613630576135c561ffff865116888863ffffffff8516928c60c0516144ea565b60206135d1848d612f48565b5101511490828b831593613608575b5050506135ef576001016135a0565b63769a20cb60e01b5f5260c0516004525f60245260445ffd5b6001600160401b039293506136208392604092612f48565b51015116911614155f828b6135e0565b509550955095925050506001915f61352f565b633a517eed60e21b5f5260c05160045260245ffd5b5060c05184146134e4565b50905060606040613675949394614dce565b92015101529291905f90565b9092506020813d6020116136ad575b8161369d60209383610ccf565b8101031261024c5751915f613462565b3d9150613690565b9091929461ffff6136ca876040850151614fd6565b5116906001821b906001600160401b036136e8896080870151614fd6565b5116926136f9896060870151614fd6565b5190855161ffff60e05151168210156137f85781602c810204602c148215171561279d57602c820260300160301161279d57816050602c82028b01015160e01c036137e657506034602c820201602c82026030011161279d57602c81028881016054015160c01c90603c810160309091011161279d57600195605c602c84028b0101519181145f146137bb575084196101005116610100525b82036137a8575050901916955b01929190613416565b5f52600b60205260405f2055179561379f565b825f52600c60205260405f20906001600160401b031982541617905584610100511761010052613792565b634724a0fd60e01b5f5260045260245ffd5b6303e07d4560e61b5f5260045260245260445ffd5b63ffffffff613820826040870151614fd6565b5116916138368361ffff60e05151168110613b0e565b81613975575b5081845190878a60e0515161ffff166101005192613859956144ea565b9190978160808701519061386c91614fd6565b516001600160401b0316988260608801519061388791614fd6565b51610120526001600160401b031692898085148015956020946138bd6138fa986138c7966138c295613967575b508c5190614fce565b613126565b614306565b604051639412e6b360e01b81526101205160048201526001600160401b03909a1660248b01529892839081906044820190565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af48015611a11575f90613935575b6001925061392e8286612f48565b520161336e565b506020823d821161395f575b8161394e60209383610ccf565b8101031261024c5760019151613920565b3d9150613941565b90506101205114155f6138b4565b84516139889163ffffffff168411614fce565b5f61383c565b50601081111561333f565b633610565160e21b5f5260045260245ffd5b633a517eed60e21b5f5260045260245ffd5b9596509690506139d760206040840151015160c051614e58565b6001600160401b036001966134d2565b9690946001600160401b03906139fb614dce565b6020604086015101526134d2565b15613a12575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160401b03602060a08301510151165f52600560205260405f205480155f14613a6c5750505f90565b60408201908151604051613aa1602082018093604080916001600160801b038151168452602081015160208501520151910152565b60608152613ab0608082610ccf565b5190201491821592613ace575b505015613ac957600190565b600290565b6001600160801b0391925060208291015151169151511611155f80613abd565b906001600160401b03809116911601906001600160401b03821161279d57565b15613b165750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15613b505750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15613b8a5750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b908160011b918083046002149015171561279d57565b15613bdb575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b909291926001600160a01b03613c2b836142dd565b161515806141bb575b6141b557604060209101510151519283519182156127c45760b4831161419d57613c6083949294612f16565b927317435cce3d1b4fa2e5f8a08ed921d57c6762a180925f5b8751811015613d3757806020896001600160401b036040613cab8585613ca2613cde9987612f48565b51015194612f48565b510151604051639412e6b360e01b81526004810193909352166001600160401b0316602482015292839081906044820190565b0381895af48015611a11575f90613d05575b60019250613cfe8289612f48565b5201613c79565b506020823d8211613d2f575b81613d1e60209383610ccf565b8101031261024c5760019151613cf0565b3d9150613d11565b5092935f9692959194613d6888613d55613d5089613bbc565b613118565b9388613d6086612f16565b9b8c926156c1565b50602c8602868104602c0361279d576030018060301161279d578260051b918383046020148415171561279d57613dab613da6602494602094614306565b613133565b97828901956056875361564160218b01536256414c60228b01536356414c3460238b01538060381c858b01538060301c60258b01538060281c60268b015380841c60278b01538060181c60288b01538060101c60298b01538060081c602a8b0153602b8a01538060081c602c8a0153602d890153604051928380926338f49afb60e01b82528b60048301525af4908115611a11575f9161416b575b50602e8601528060081c604e860153604f8501536030955f935b8351851015613f5d57613e738585612f48565b518887019060ff8760181c16602083015361ffff8760101c16602183015362ffffff8760081c16602283015363ffffffff8716602383015360048a018a1161279d57604081015166ffffffffffffff6001600160401b0382169160ff8160381c16602486015361ffff8160301c16602586015362ffffff8160281c16602686015363ffffffff8160201c16602786015364ffffffffff8160181c16602886015365ffffffffffff8160101c16602986015360081c16602a840153602b830153600c8a018a1161279d576020602c910151910152602c880180981161279d5760018895019450613e60565b92509250935f955b8351871015613f9657613f788785612f48565b51602082870101526020810180911161279d57600190960195613f65565b50939291509350613fa781846143a2565b9160405191613fd860218460208101945f8652845180918484015e81015f838201520301601f198101855284610ccf565b61600083511161413f5750906001600160a01b0391614090602a7fffff000000000000000000000000000000000000000000000000000000000000845160f01b169360405193849160208301967f6100000000000000000000000000000000000000000000000000000000000000885260218401527f80600a3d393df30000000000000000000000000000000000000000000000000060238401525180918484015e81015f838201520301601f198101835282610ccf565b51905ff0168015614117576008549260065560408201516007557fffff0000000000000000000000000000000000000000000000000000000000007dffff00000000000000000000000000000000000000000000000000000000602067ffffffffffffffff60a01b855160a01b1694015160e01b1693161717176008555f6009555f600a55565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b90506020813d602011614195575b8161418660209383610ccf565b8101031261024c57515f613e46565b3d9150614179565b8263156f758160e31b5f5260045260b460245260445ffd5b50509050565b5061ffff60085460e01c161515613c34565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f205416155f1461425457805f525f60205260405f206001600160a01b0383165f5260205260405f20600160ff198254161790556001600160a01b03339216907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f2054165f1461425457805f525f60205260405f206001600160a01b0383165f5260205260405f2060ff1981541690556001600160a01b03339216907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b600654036142f4576001600160a01b036008541690565b5f90565b906001820180921161279d57565b9190820180921161279d57565b61ffff16602c810290808204602c149015171561279d576030018060301161279d5790565b9190823b6001811115614386575f19810190811161279d57811161436a57600161436182613133565b9360208501903c565b6001600160a01b0383633610565160e21b5f521660045260245ffd5b6001600160a01b0384633610565160e21b5f521660045260245ffd5b9190915f60606040516143b481610c7e565b828152826020820152826040820152015260308351106137e6576356414c34602084015160e01c036137e657602483015160c01c92602c810151938460f01c91602e81015190604e81015160f01c6040519361440f85610c7e565b84526020840192858452604085015260608401938185529785159586156144df575b5085156144ad575b5050831561444c575b5050506137e65750565b519051915192509061ffff838116916144659116614313565b9081831493841561447e575b50505050155f8080614442565b621fffe0919293945060051b16908082046020149015171561279d576144a391614306565b145f808080614471565b9091945060ef1c6201fffe61fffe82169116810361279d575f190161ffff811161279d5761ffff161415925f80614439565b60b41095505f614431565b939092959461ffff63ffffffff8216931683101561459b5761ffff1690602c830292808404602c148115171561279d5783603001948560301161279d5784019581605088015160e01c036137e657506001901b9081166145805760348301841161279d57605485015160c01c5b961661456d5750603c011061279d57605c015190565b925050505f52600b60205260405f205490565b815f52600c6020526001600160401b0360405f205416614557565b82856303e07d4560e61b5f5260045260245260445ffd5b999493979198909695995f966145c781611318565b80614d68575089151580614d5c575b15614d2557602001946145e886615301565b6146406145f5368e610de9565b9161245c60405160208101906146288287604080916001600160801b038151168452602081015160208501520151910152565b60608152614637608082610ccf565b5190209161501e565b6020810151808703614cf55750614676633b9aca006001600160801b038093511604614670814242821115614fe7565b42613126565b1663ffffffff600454169081811015614cc75750505f905f5b8b8a818310614bda575b5050505015614b9d5750506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690813b1561024c579187926040519384927f87d3a9b100000000000000000000000000000000000000000000000000000000845260648401906004850152606060248501525260848201908960051b9860848a850101928b8a915f603e198d360301925b8110614b2f5750505050600319848403016044850152808352602083019260208260051b82010193835f92601e1982360301905b8585106149c0575050505050505091815f81819503925af18015611a11576149ab575b50600185116147b6575b50505050506001600160801b036147b0633b9aca00923690610de9565b51160490565b6147c7909691929496959395615301565b948335946001600160801b0386168096036149a7576147e588610e74565b976147f3604051998a610ccf565b885260208801918101903682116149a35780979597925b82841061490857505050506001600160401b03829316925b86518110156148eb576148358188612f48565b519085604051602081019087825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801918a5b8181106148b457505050508161489c60019760206148aa940151605f19848303016080850152610c10565b03601f198101835282610ccf565b5190205d01614822565b9193949650919496976020806148d660019360bf198b82030188528951610c10565b970194019101918c9694939298979598614871565b5094505050506001600160801b036147b0633b9aca005f80614793565b83989698356001600160401b03811161499f57820160408136031261499f576040519061493482610c34565b80356001600160401b03811161499b57810136601f8201121561499b57614962903690602081359101615337565b825260208101356001600160401b03811161499b5791614989602094928594369101610d41565b8382015281520193019297959761480a565b8880fd5b8680fd5b8480fd5b8380fd5b6149b89192505f90610ccf565b5f905f614789565b9193959750919395601f1982820301855287358381121561024c578401906149ec602082019280615452565b8091936020845252604082019060408160051b8401019380935f915b838310614a2f57505050505050602080600192990195019501929091899796949592614766565b909192939495603f19838203018652614a48878361549a565b8035600281101561024c57614a5c81611318565b8252614a7f614a6e6020830183615486565b6060602085015260608401906154ae565b906040810135609e198236030181121561024c576001936020938493614b219301916040818303910152614b13614af5614aca614abc85806153c2565b60a0865260a08601916153a2565b614ad5878601610dc8565b151587850152614ae86040860186615486565b84820360408601526154ae565b92614b0260608201610dc8565b151560608401526080810190615486565b9060808184039101526154ae565b980196019493019190614a08565b9193949650919460831989820301835285358481121561024c576020614b8b6001938f83940190614b7e614b74614b668480615452565b6040855260408501916153f3565b92858101906153c2565b91858185039101526153a2565b9701930191019088969493918e614732565b612b6f6040519283927ffef760c70000000000000000000000000000000000000000000000000000000084526020600485015260248401916153f3565b614bf8614bf284614c1094614bff9498969798615315565b806130e3565b3691615337565b614c0a368789615337565b906157e2565b15614cbd5750610751614c27614c31928d8c615315565b60208101906130b1565b805182518082149182614ca7575b505015614c5357505060015f808b8a614699565b90612b6f614c95926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610c10565b83810360031901602485015290610c10565b9091506020830120906020840120145f80614c3f565b919060010161468f565b7fdb020436000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b86907f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b897fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b5061ffff8a11156145d6565b8760ff602492614d7781611318565b614d8081611318565b7f112d89cc00000000000000000000000000000000000000000000000000000000835216600452fd5b614dba6001600160a01b03916142dd565b1615614dc857600754600191565b5f905f90565b60405190614ddb82610c7e565b606082525f6020830152604051614df181610c7e565b606081525f60208201525f60408201525f606082015260408301525f60608301528160405190614e22602083610ccf565b5f825252565b52565b9081602091031261024c575190565b15614e43575050565b63769a20cb60e01b5f5260045260245260445ffd5b8151518015614fbb5760b48111614fa45750604051917fc87f1f69000000000000000000000000000000000000000000000000000000008352602060048401528260a4810182519060806024840152815180915260c4830190602060c48260051b8601019301915f905b828210614f755750505050826001600160401b036060614f008594602080980151151560448701526040850151602319878303016064880152611322565b92015116608483015203817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4908115611a11575f91614f3f575b614f3d9250808214614e3a565b565b90506020823d602011614f6d575b81614f5a60209383610ccf565b8101031261024c57614f3d915190614f30565b3d9150614f4d565b91936001919395506020614f94819260c3198c82030186528851611322565b9601920192018794939192614ec2565b63156f758160e31b5f5260045260b460245260445ffd5b50633a517eed60e21b5f5260045260245ffd5b156137e65750565b906010811015612f5c5760051b0190565b15614ff0575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160401b03165f52600560205260405f2054801561503c5790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040810151519060046020830151926001600160401b0384511690604063ffffffff602087015116950151906020825192015160208063ffffffff8351169201519251015192604051946150b786610c99565b85526020850197885260408501908152606085019182526080850192835260a0850193845261ffff60e0880151169760808801519760a08101519060c081015161012082015161014083015191610180610160850151940151946040519e8f9d8e7f4cc22bb7000000000000000000000000000000000000000000000000000000008152015260248d019d5f5b600881106152cd57505060209d506151eb8d63ffffffff976151d8829f9d9a95986152269f9c986151c59060c09f9b6001600160401b039a61519e6151b2926151938e9c6101248c01906113f6565b6101648a01906113f6565b6102406101a4890152610244880190610ba1565b868103600319016101c488015290610bd4565b848103600319016101e486015290610b68565b916102046003198285030191015261141d565b8c8103600319016102248e0152995116895251168b880152516040870152511660608501525160808401525160a08301829052910190610c10565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af1908115611a11575f91615293575b501561526b57565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90506020813d6020116152c5575b816152ae60209383610ccf565b8101031261024c576152bf90611110565b5f615263565b3d91506152a1565b91939597999b9d5091939597999b9d602080600192855181520193019101908f9d9b99979593919e9c9a989694929e615144565b356001600160401b038116810361024c5790565b9190811015612f5c5760051b81013590603e198136030182121561024c570190565b92919061534381610e74565b936153516040519586610ccf565b602085838152019160051b81019183831161024c5781905b838210615377575050505050565b81356001600160401b03811161024c576020916153978784938701610d41565b815201910190615369565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e198236030181121561024c5701602081359101916001600160401b03821161024c57813603831361024c57565b90602083828152019260208260051b82010193835f925b84841061541a5750505050505090565b909192939495602080615442600193601f1986820301885261543c8b886153c2565b906153a2565b980194019401929493919061540a565b9035601e198236030181121561024c5701602081359101916001600160401b03821161024c578160051b3603831361024c57565b9035607e198236030181121561024c570190565b9035605e198236030181121561024c570190565b6154e76154cc6154be83806153c2565b6080865260808601916153a2565b6154d960208401846153c2565b9085830360208701526153a2565b6154f46040830183615486565b9083810360408501528135600381101561024c5761551181610b24565b81526020820135600381101561024c5761552a81610b24565b6020820152604082013590600382101561024c576080615566615581948461555761557696999899610b24565b604085015260608101906153c2565b91909281606082015201916153a2565b926060810190615452565b90916060818503910152808352602083019060208160051b85010193835f915b8383106155b15750505050505090565b909192939495601f198282030186526155ca878461549a565b803591600383101561024c576156276020928392856155ea600197610b24565b815261561961560e6155fe868501856153c2565b60608886015260608501916153a2565b9260408101906153c2565b9160408185039101526153a2565b9801960194930191906155a1565b90615672575080511561564a57602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b815115806156b8575b615683575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b1561567b565b9493919092946156d18483613126565b600181146157c2576156e38791615ac4565b92836157046156f28289614306565b846156fc896142f8565b918a886156c1565b615715615737966157319299614306565b9161572b613d506157258a6142f8565b92613bbc565b90614306565b936156c1565b60405163a5641f6f60e01b8152600481019390935260248301526020826044817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49182156157bd575f92615788575b50614e28908294612f48565b614e289192506157af9060203d6020116157b6575b6157a78183610ccf565b810190614e2b565b919061577c565b503d61579d565b611a11565b506157d59150926157de92949593612f48565b51928392612f48565b5290565b908151815103614254575f5b825181101561584c576158018184612f48565b515161580d8284612f48565b5151036158455761581e8184612f48565b51602081519101206158308284612f48565b516020815191012003615845576001016157ee565b5050505f90565b505050600190565b5f19811461279d5760010190565b9392909491828414808091615aa8575b615a805761ffff83116127c4576158898783613126565b906001821461599c575061589c90615ac4565b936158c36158aa8689614306565b9561572b613d506157256158bd886142f8565b976142f8565b92815b85811087828a8361596b575b505050156158e8576158e390615854565b6158c6565b94959697918591886158fa948b615862565b946159059596615862565b60405163a5641f6f60e01b81526004810192909252602482015260208180604481015b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49081156157bd575f91615952575090565b610d5c915060203d6020116157b6576157a78183610ccf565b63ffffffff92935061599291604060606159889301510151614fd6565b5163ffffffff1690565b161087828a6158d2565b925050509493929415615a155750506159e48360209261592894955191848101516159cc604083015161ffff1690565b9063ffffffff60c060a08501519401519416946144ea565b6040519384928392639412e6b360e01b8452600484019092916001600160401b036020916040840195845216910152565b615a208293926142f8565b1490811591615a55575b50615a41576080615a3d92930151612f48565b5190565b8251634724a0fd60e01b5f5260045260245ffd5b9050615a78615a6f61598884604060608901510151614fd6565b63ffffffff1690565b14155f615a2a565b505092915050610d5c9250805190615aa26040602083015192015161ffff1690565b91615b01565b50615abf615abb838960e08a0151615ae0565b1590565b615872565b9060015b8060011b9083821015615adb5750615ac8565b925050565b91615aed82600192613126565b1b905f19820191821161279d571b16151590565b9392909161ffff16602c810290808204602c149015171561279d57603001908160301161279d578060051b908082046020149015171561279d57615b4491614306565b906020820180831161279d57815110615b605701602001519150565b83634724a0fd60e01b5f5260045260245ffdfea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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
