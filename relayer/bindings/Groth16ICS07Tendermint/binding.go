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
	Bin: "0x610160604052346100fe57616e22803803809161001b82610116565b6101603960e0816101600191126100fe576100346101a0565b61003f6101806101b7565b61004a6101a06101b7565b6100556101c06101b7565b6101e0516001600160401b0381116100fe578561017f820112156100fe576100a2958161018061008b93610160015191016101e6565b91610200519361009c6102206101b7565b956107a2565b604051615ddd9081610f85823960805181614682015260a05181613b33015260c05181505060e051818181612410015261265d01526101005181818161096c01526109b301526101205181612cd1015261014051816123db0152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b610160601f91909101601f19168101906001600160401b0382119082101761013d57604052565b610102565b604081019081106001600160401b0382111761013d57604052565b601f909101601f19168101906001600160401b0382119082101761013d57604052565b6040519061018f60e08361015d565b565b6040519061018f60408361015d565b61016051906001600160a01b03821682036100fe57565b51906001600160a01b03821682036100fe57565b6001600160401b03811161013d57601f01601f191660200190565b9291926101f2826101cb565b91610200604051938461015d565b8294818452818301116100fe578281602093845f96015e010152565b9080601f830112156100fe578151610236926020016101e6565b90565b519060ff821682036100fe57565b91908260409103126100fe5760405161025f81610142565b602061027881839561027081610239565b855201610239565b910152565b51906001600160401b03821682036100fe57565b91908260409103126100fe576040516102a981610142565b60206102788183956102ba8161027d565b85520161027d565b519063ffffffff821682036100fe57565b519081151582036100fe57565b519060028210156100fe57565b6020818303126100fe578051906001600160401b0382116100fe5701610120818303126100fe5761031c610180565b8151909290916001600160401b0383116100fe5761036382610346610100946103a196850161021c565b86526103558160208501610247565b602087015260608301610291565b604085015261037460a082016102c2565b606085015261038560c082016102c2565b608085015261039660e082016102d3565b60a0850152016102e0565b60c082015290565b90600182811c921680156103d7575b60208310146103c357565b634e487b7160e01b5f52602260045260245ffd5b91607f16916103b8565b601f82116103ee57505050565b5f5260205f20906020601f840160051c83019310610426575b601f0160051c01905b81811061041b575050565b5f8155600101610410565b9091508190610407565b6002111561043a57565b634e487b7160e01b5f52602160045260245ffd5b90600281101561043a5769ff00000000000000000082549160481b169069ff0000000000000000001916179055565b80518051906001600160401b03821161013d576104a68261049f6001546103a9565b60016103e1565b602090601f8311600114610603578260c09361018f95936104dc935f926105f8575b50508160011b915f199060031b1c19161790565b6001555b6020818101518051600280549284015161ff0060089190911b1660ff90921661ffff199093169290921717905560408083015180516003805492909401516fffffffffffffffff0000000000000000931b929092166001600160401b039092166001600160801b03199091161717905560608101516105769063ffffffff1660049063ffffffff1663ffffffff19825416179055565b6105ae61058a608083015163ffffffff1690565b60049067ffffffff0000000082549160201b169067ffffffff000000001916179055565b6105e66105be60a0830151151590565b60049068ff0000000000000000825491151560401b169068ff00000000000000001916179055565b01516105f181610430565b600461044e565b015190505f806104c8565b60015f52601f19831691905f516020616de25f395f51905f52925f5b818110610661575092600192859260c09661018f989610610649575b505050811b016001556104e0565b01515f1960f88460031b161c191690555f808061063b565b9293602060018192878601518155019501930161061f565b604051905f826001549161068c836103a9565b80835292600181169081156106fc57506001146106b0575b61018f9250038361015d565b5060015f90815290915f516020616de25f395f51905f525b8183106106e057505090602061018f928201016106a4565b60209193508060019154838589010152019101909184926106c8565b6020925061018f94915060ff191682840152151560051b8201016106a4565b15610724575050565b6355bace6f60e01b5f9081526001600160401b039182166004529116602452604490fd5b634e487b7160e01b5f52601160045260245ffd5b63ffffffff6107089116019063ffffffff821161077557565b610748565b15610783575050565b9063ffffffff80926333fae18560e21b5f52166004521660245260445ffd5b91936107e36107e891969294967f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a6123361010052602080825183010191016102ed565b61047d565b6107f0610679565b602081519101206101205261080b610806610679565b61091e565b6101405261087e610865610838602061082a610825610679565b6109f8565b01516001600160401b031690565b60035490610856906001600160401b0380841691908116821461071b565b60401c6001600160401b031690565b6001600160401b03165f90815260056020526040902090565b556001600160a01b0390811660805290811660a05290811660c0521660e0526004546108d69063ffffffff81166108c46108b78261075c565b9260201c63ffffffff1690565b9163ffffffff8084169116111561077a565b6001600160a01b0381166108fd57506108ed610cdc565b506108fa61010051610d5e565b50565b8061090a6108fa92610bde565b5061091481610c54565b5061010051610db7565b61092790610e62565b8051600181018091116107755761093d90610e30565b8051156109855760209181610959845f94019284845382610f2f565b50604051918291518091835e8101838152039060025afa1561097a575f5190565b6040513d5f823e3d90fd5b634e487b7160e01b5f52603260045260245ffd5b604051906109a682610142565b5f602083606081520152565b8015610775575f190190565b5f1981019190821161077557565b9190820391821161077557565b908151811015610985570160200190565b906001820180921161077557565b610a00610999565b5080518015908115610bd2575b50610bc3575f19908051805b610b73575b505f198214610b5557600360fc1b6001600160f81b0319610a58610a4a610a44866109ea565b856109d9565b516001600160f81b03191690565b161480610b5f575b610b55575f90610a6f836109ea565b915b8151831015610b0657610a90610a8a610a4a85856109d9565b60f81c90565b60ff811660308110908115610afb575b50610aee57600a82026001600160401b03908116602f1990920160ff1691909101811691168110610ad657600190920191610a71565b50915050610ae2610191565b9081525f602082015290565b5050915050610ae2610191565b60399150115f610aa0565b9150916001811190811591610b49575b50610b3a5761023690610b27610191565b9283526001600160401b03166020830152565b63384c687760e21b5f5260045ffd5b602b915010155f610b16565b9050610ae2610191565b506002610b6d8383516109cc565b11610a60565b602d60f81b610b9d610b90610a4a610b8a856109be565b866109d9565b6001600160f81b03191690565b14610bb157610bab906109b2565b80610a19565b610bbc9192506109be565b905f610a1e565b6329120bff60e21b5f5260045ffd5b6040915010155f610a0d565b6001600160a01b0381165f9081525f516020616e025f395f51905f52602052604090205460ff16610c4f576001600160a01b03165f8181525f516020616e025f395f51905f5260205260408120805460ff191660011790553391905f516020616d625f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f516020616d825f395f51905f52602052604090205460ff16610c4f576001600160a01b0381165f9081525f516020616d825f395f51905f5260205260409020805460ff1916600117905533906001600160a01b03165f516020616dc25f395f51905f525f516020616d625f395f51905f525f80a4600190565b5f80525f516020616d825f395f51905f526020525f516020616da25f395f51905f525460ff16610d5a575f8080525f516020616d825f395f51905f526020525f516020616da25f395f51905f52805460ff1916600117905533905f516020616dc25f395f51905f525f516020616d625f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f14610c4f575f818152602081815260408083208380529091528120805460ff1916600117905533915f516020616d625f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16610e2a575f818152602081815260408083206001600160a01b03861684529091529020805460ff1916600117905533916001600160a01b0316905f516020616d625f395f51905f525f80a4600190565b50505f90565b90610e3a826101cb565b610e47604051918261015d565b8281528092610e58601f19916101cb565b0190602036910137565b805115610ece5780516001905b6080811015610ec05750806001019081600111610775576001908351010180911161077557610ea0610ebc91610e30565b91610eb6610ead84610eea565b82519085610ef6565b83610f59565b5090565b60019060071c910190610e6f565b50604051610edd60208261015d565b5f81525f36602083013790565b6020600a910153600190565b9092919083016020015b6080821015610f1457906001929391530190565b600180916080607f85161781530193019060071c9092610f00565b908051918215610f51576021602084930191015e600101806001116107755790565b505050600190565b908092918251928315610f7d57839260208092019201015e81018091116107755790565b505050509056fe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a7146101245780630bece3561461011f578063248a9ca31461011a5780632f2ff15d1461011557806336568abe14610110578063536c2ad31461010b5780636a28f000146101065780638a8e4c5d1461010157806391d14854146100fc578063974a74c4146100f7578063a217fddf146100f2578063a6f031bb146100ed578063ac9650d8146100e8578063d547741f146100e3578063db3e1fa4146100de578063ddba6537146100d95763ef913a4b146100d4575f80fd5b610a23565b61098f565b610955565b610926565b6108ba565b61072a565b610710565b6105b0565b61056f565b610550565b6104ab565b610445565b61034e565b610318565b6102c0565b61023b565b346101c55760203660031901126101c5576004357fffffffff0000000000000000000000000000000000000000000000000000000081168091036101c557807f7965db0b000000000000000000000000000000000000000000000000000000006020921490811561019b575b506040519015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150145f610190565b5f80fd5b9060206003198301126101c5576004356001600160401b0381116101c557826023820112156101c5578060040135926001600160401b0384116101c557602484830101116101c5576024019190565b634e487b7160e01b5f52602160045260245ffd5b6003111561023657565b610218565b346101c55760206102a261024e366101c9565b9061026160ff60045460401c1615610b9f565b7fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5f525f845260405f205f8052845260ff60405f205416156102b35761238b565b604051906102af8161022c565b8152f35b6102bb612a8e565b61238b565b346101c55760203660031901126101c55760206102ea6004355f525f602052600160405f20015490565b604051908152f35b60409060031901126101c557600435906024356001600160a01b03811681036101c55790565b346101c55761034c610329366102f2565b90610347610342825f525f602052600160405f20015490565b612afd565b6135f5565b005b346101c55761035c366102f2565b336001600160a01b038216036103755761034c9161368d565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b90602080835192838152019201905f5b8181106103ba5750505090565b825163ffffffff168452602093840193909201916001016103ad565b90602080835192838152019201905f5b8181106103f35750505090565b82518452602093840193909201916001016103e6565b90602080835192838152019201905f5b8181106104265750505090565b82516001600160401b0316845260209384019390920191600101610419565b346101c55760203660031901126101c55761048161048f61049d61046a600435612728565b91939060405195869560608752606087019061039d565b9085820360208701526103d6565b908382036040850152610409565b0390f35b5f9103126101c557565b346101c5575f3660031901126101c557335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff16156105395760045460ff8160401c16156105115768ff00000000000000001916600455005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b346101c55761055e366101c9565b5050636d40ebe160e11b5f5260045ffd5b346101c557602060ff6105a4610584366102f2565b905f525f845260405f20906001600160a01b03165f5260205260405f2090565b54166040519015158152f35b346101c55760203660031901126101c5576004356001600160401b0381116101c557806004019061016060031982360301126101c5576105f860ff60045460401c1615610b9f565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610703575b61014481019061065a828461287a565b9050156106db576106bc6106cb9261049d9461067960448501826128ac565b61068960648794939401836128ac565b90608488013592610104890135956106a087610e0f565b60a46106c36106b36101248d01896128ac565b9b909a8961287a565b3691610cf9565b9a0195613abf565b6040519081529081906020820190565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b61070b612a8e565b61064a565b346101c5575f3660031901126101c55760206040515f8152f35b346101c55760203660031901126101c5576004356001600160401b0381116101c5578060040161014060031983360301126101c55761049d916106cb9161077960ff60045460401c1615610b9f565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac4754161561082a575b6107d860448301826128ac565b916107e660648501826128ac565b6101048601359291608487013591906107fe85610e0f565b61080c6101248901856128ac565b97909660a46040519a61082060208d610c6e565b5f8c520195613abf565b610832612a8e565b6107cb565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b602081016020825282518091526040820191602060408360051b8301019401925f915b83831061088d57505050505090565b90919293946020806108ab600193603f198682030187528951610837565b9701930193019193929061087e565b346101c55760203660031901126101c5576004356001600160401b0381116101c557366023820112156101c5578060040135906001600160401b0382116101c5573660248360051b830101116101c55761049d91602461091a920161299a565b6040519182918261085b565b346101c55761034c610937366102f2565b90610950610342825f525f602052600160405f20015490565b61368d565b346101c5575f3660031901126101c55760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346101c55761099d366101c9565b50506109b160ff60045460401c1615610b9f565b7f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610a00575b636d40ebe160e11b5f5260045ffd5b610a0990612afd565b5f6109f1565b906020610a20928181520190610837565b90565b346101c5575f3660031901126101c55760405160208082015261012060408201525f600154610a5181612a56565b90816101608501526001811690815f14610b7a5750600114610b19575b61049d83610b0d8185610a9360608301602060ff600254818116845260081c16910152565b610ab560a0830160206001600160401b03600354818116845260401c16910152565b60045463ffffffff811660e0840152610aff90602081901c63ffffffff16610100850152610aee610120850160ff8360401c1615159052565b60ff61014085019160481c16611c4a565b03601f198101835282610c6e565b60405191829182610a0f565b91905060015f527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6915f905b808210610b5e5750909150810161018001610b0d610a6e565b9192600181602092546101808588010152019101909291610b45565b60ff19166101808086019190915291151560051b84019091019150610b0d9050610a6e565b15610ba657565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b03821117610bfd57604052565b610bce565b606081019081106001600160401b03821117610bfd57604052565b608081019081106001600160401b03821117610bfd57604052565b60a081019081106001600160401b03821117610bfd57604052565b60c081019081106001600160401b03821117610bfd57604052565b90601f801991011681019081106001600160401b03821117610bfd57604052565b60405190610c9e60e083610c6e565b565b60405190610c9e61026083610c6e565b60405190610c9e6101c083610c6e565b60405190610c9e608083610c6e565b60405190610c9e60c083610c6e565b6001600160401b038111610bfd57601f01601f191660200190565b929192610d0582610cde565b91610d136040519384610c6e565b8294818452818301116101c5578281602093845f960137010152565b9080601f830112156101c557816020610a2093359101610cf9565b60ff8116036101c557565b91908260409103126101c557604051610d6d81610be2565b60208082948035610d7d81610d4a565b8452013591610d8b83610d4a565b0152565b6001600160401b038116036101c557565b3590610c9e82610d8f565b91908260409103126101c557604051610dc381610be2565b60208082948035610dd381610d8f565b8452013591610d8b83610d8f565b63ffffffff8116036101c557565b3590610c9e82610de1565b801515036101c557565b3590610c9e82610dfa565b600211156101c557565b3590610c9e82610e0f565b919091610120818403126101c557610e3a610c8f565b928135916001600160401b0383116101c557610e7f82610e6261010094610ebd968501610d2f565b8752610e718160208501610d55565b602088015260608301610dab565b6040860152610e9060a08201610def565b6060860152610ea160c08201610def565b6080860152610eb260e08201610e04565b60a086015201610e19565b60c0830152565b6001600160801b038116036101c557565b3590610c9e82610ec4565b91908260609103126101c557604051610ef881610c02565b60408082948035610f0881610ec4565b8452602081013560208501520135910152565b8092910391606083126101c557604051610f3481610be2565b6040819483358352601f1901126101c5576020906040805193610f5685610be2565b83810135610f6381610de1565b85520135828401520152565b6001600160401b038111610bfd5760051b60200190565b919060c0838203126101c557604051610f9e81610c1d565b80938035610fab81610d8f565b82526020810135610fbb81610de1565b6020830152610fcd8360408301610f1b565b604083015260a0810135906001600160401b0382116101c557019180601f840112156101c557823592610fff84610f6f565b9361100d6040519586610c6e565b80855260208086019160051b830101918383116101c55760208101915b83831061103c57505050505060600152565b82356001600160401b0381116101c5578201906040828703601f1901126101c5576040519161106a83610be2565b602081013560048110156101c557835260408101356001600160401b0381116101c5576020910101906080828803126101c557604051926110aa84610c1d565b82356001600160401b0381116101c557886110c6918501610d2f565b845260208301356110d681610ec4565b602085015260408301356110e981610dfa565b60408501526060830135936001600160401b0385116101c55761111189602096879601610d2f565b60608201528382015281520192019161102a565b91906040838203126101c5576040519061113e82610be2565b819380356001600160401b0381116101c55781016102c0818403126101c557611165610ca0565b906111708482610dab565b825260408101356001600160401b0381116101c55784611191918301610d2f565b60208301526111a260608201610da0565b60408301526111b360808201610ed5565b60608301526111c460a08201610e04565b60808301526111d68460c08301610f1b565b60a08301526111e86101208201610e04565b60c083015261014081013560e08301526112056101608201610e04565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526112546102208201610e04565b6101c08301526102408101356101e08301526112736102608201610e04565b6102008301526102808101356102208301526102a0810135906001600160401b0382116101c5576112a691859101610d2f565b61024082015283526020810135916001600160401b0383116101c5576020926112cf9201610f86565b910152565b91906080838203126101c557604051906112ed82610c1d565b819380356001600160401b0381116101c55760609261130d918301610d2f565b835260208101356020840152604081013561132781610d8f565b60408401520135908160070b82036101c55760600152565b9190916080818403126101c5576040519061135982610c1d565b819381356001600160401b0381116101c557820181601f820112156101c557803561138381610f6f565b916113916040519384610c6e565b81835260208084019260051b820101918483116101c55760208201905b8382106113ff575050505083526113c760208301610e04565b60208401526040820135906001600160401b0382116101c557826113f4606094926112cf948694016112d4565b604086015201610da0565b81356001600160401b0381116101c557602091611421888480948801016112d4565b8152019101906113ae565b919060a0838203126101c5576040519061144582610c1d565b819380356001600160401b0381116101c55782611463918301611125565b835260208101356001600160401b0381116101c5578261148491830161133f565b60208401526114968260408301610dab565b60408401526080810135916001600160401b0383116101c5576060926112cf920161133f565b9080601f830112156101c557610100604051926114d98285610c6e565b839181019283116101c557905b8282106114f35750505090565b81358152602091820191016114e6565b9080601f830112156101c5576040519161151e604084610c6e565b8290604081019283116101c557905b82821061153a5750505090565b813581526020918201910161152d565b359061ffff821682036101c557565b9080601f830112156101c557813561157081610f6f565b9261157e6040519485610c6e565b81845260208085019260051b8201019283116101c557602001905b8282106115a65750505090565b6020809183356115b581610de1565b815201910190611599565b9080601f830112156101c55781356115d781610f6f565b926115e56040519485610c6e565b81845260208085019260051b8201019283116101c557602001905b82821061160d5750505090565b8135815260209182019101611600565b9080601f830112156101c557813561163481610f6f565b926116426040519485610c6e565b81845260208085019260051b8201019283116101c557602001905b82821061166a5750505090565b60208091833561167981610d8f565b81520191019061165d565b9080601f830112156101c557813561169b81610f6f565b926116a96040519485610c6e565b81845260208085019260051b8201019283116101c557602001905b8282106116d15750505090565b6020809183356116e081610dfa565b8152019101906116c4565b9080601f830112156101c557610200604051926117088285610c6e565b839181019283116101c557905b8282106117225750505090565b8135815260209182019101611715565b9080601f830112156101c5576102006040519261174f8285610c6e565b839181019283116101c557905b8282106117695750505090565b60208091833561177881610d8f565b81520191019061175c565b9190610640838203126101c5576040519061179d82610c38565b81938035835260208101356117b181610d4a565b602084015281605f820112156101c5576040516117d061020082610c6e565b806102408301918483116101c55760408401905b8382106118125750506112cf9261180785608096946104409460408a01526116eb565b606087015201611732565b60208091833561182181610de1565b8152019101906117e4565b6020818303126101c5578035906001600160401b0382116101c55701610940818303126101c55761185b610cb0565b9181356001600160401b0381116101c55781611878918401610e24565b83526118878160208401610ee0565b602084015260808201356001600160401b0381116101c557816118ab91840161142c565b60408401526118bc60a08301610ed5565b60608401526118ce8160c084016114bc565b60808401526118e1816101c08401611503565b60a08401526118f4816102008401611503565b60c0840152611906610240830161154a565b60e08401526102608201356001600160401b0381116101c5578161192b918401611559565b6101008401526102808201356001600160401b0381116101c557816119519184016115c0565b6101208401526102a08201356001600160401b0381116101c5578161197791840161161d565b6101408401526102c08201356001600160401b0381116101c5578161199d918401611559565b6101608401526102e08201356001600160401b0381116101c557826119ca83610300936119d69601611684565b61018086015201611783565b6101a082015290565b81601f820112156101c5578051906119f682610cde565b92611a046040519485610c6e565b828452602083830101116101c557815f9260208093018386015e8301015290565b91908260409103126101c557604051611a3d81610be2565b60208082948051611a4d81610d4a565b8452015191610d8b83610d4a565b91908260409103126101c557604051611a7381610be2565b60208082948051611a8381610d8f565b8452015191610d8b83610d8f565b5190610c9e82610de1565b5190610c9e82610dfa565b5190610c9e82610e0f565b919091610120818403126101c557611ac8610c8f565b928151916001600160401b0383116101c557611b0d82611af061010094610ebd9685016119df565b8752611aff8160208501611a25565b602088015260608301611a5b565b6040860152611b1e60a08201611a91565b6060860152611b2f60c08201611a91565b6080860152611b4060e08201611a9c565b60a086015201611aa7565b5190610c9e82610ec4565b91908260609103126101c557604051611b6e81610c02565b60408082948051611b7e81610ec4565b8452602081015160208501520151910152565b6020818303126101c5578051906001600160401b0382116101c55701610180818303126101c55760405191611bc583610c53565b81516001600160401b0381116101c55782611be88361014093611c389601611ab2565b8552611bf78360208301611b56565b6020860152611c098360808301611b56565b6040860152611c1a60e08201611b4b565b6060860152611c2d836101008301611a5b565b608086015201611a5b565b60a082015290565b6002111561023657565b90611c5482611c40565b52565b90610a209061010060c0611c7685516101208552610120850190610837565b9460ff60208083015182815116828801520151166040850152611cb6604082015160608601906001600160401b0360208092828151168552015116910152565b606081015163ffffffff1660a0850152608081015163ffffffff168483015260a0810151151560e08501520151910190611c4a565b90606060c08201926001600160401b03815116835263ffffffff6020820151166020840152611d3c6040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519160c060a0830152825180915260e0820190602060e08260051b8501019401925f905b828210611d7057505050505090565b909192939460df198282030185528551908151600481101561023657611dee8260206001958195948295520151906040848201526060611dbc83516080604085015260c0840190610837565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610837565b9701950193920190611d61565b90606080611e128451608085526080850190610837565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b828210611eae575050505090606080611e9d610a2094611e8b6020880151602087019015159052565b60408701518582036040870152611dfb565b9401516001600160401b0316910152565b90919293602080611ecb600193609f198a82030186528851611dfb565b960192019201909291611e62565b610a20916060612066612054845160a0855260206120408251604060a0890152611f1c60e0890182516001600160401b0360208092828151168552015116910152565b610240611f3a848301516102c06101208c01526103a08b0190610837565b60408301516001600160401b03166101408b015291808901516001600160801b03166101608b0152608081015115156101808b015260a081015180516101a08c0152602090810151805163ffffffff166101c08d015201516101e08b015260c081015115156102008b015260e08101516102208b01526101008101511515828b01526101208101516102608b01526101408101516102808b01526101608101516102a08b01526101808101516102c08b01526101a08101516102e08b01526101c081015115156103008b01526101e08101516103208b015261020081015115156103408b01526102208101516103608b0152015188820360df19016103808a0152610837565b910151858203609f190160c0870152611ceb565b60208501518482036020860152611e3b565b9261208e604082015160408501906001600160401b0360208092828151168552015116910152565b0151906080818403910152611e3b565b905f905b600882106120af57505050565b60208060019285518152019301910190916120a2565b905f905b600282106120d657505050565b60208060019285518152019301910190916120c9565b90602080835192838152019201905f5b8181106121095750505090565b825115158452602093840193909201916001016120fc565b905f905b6010821061213257505050565b6020806001928551815201930191019091612125565b905f905b6010821061215957505050565b6020806001926001600160401b0386511681520193019101909161214c565b8051825260ff60208201511660208301526040810151604083015f905b601082106121c757505050906104406080836121bd6060610c9e960151610240860190612121565b0151910190612148565b60208060019263ffffffff865116815201930191019091612195565b90610a20906103006101a06123146123006122ec6122d86122c46122566122158b516109408b526109408b0190611c57565b61224460208d015160208c0190604080916001600160801b038151168452602081015160208501520151910152565b60408c01518a820360808c0152611ed9565b60608b01516001600160801b031660a08a015261227b60808c015160c08b019061209e565b61228e60a08c01516101c08b01906120c5565b6122a160c08c01516102008b01906120c5565b60e08b015161ffff166102408a01526101008b01518982036102608b015261039d565b6101208a01518882036102808a01526103d6565b6101408901518782036102a0890152610409565b6101608801518682036102c088015261039d565b6101808701518582036102e08701526120ec565b940151910190612178565b906020610a209281815201906121e3565b6040513d5f823e3d90fd5b612353604092959493956060835260608301906121e3565b9460208201520152565b610c9e909291926060810193604080916001600160801b038151168452602081015160208501520151910152565b6123979181019061182c565b6123a081612b44565b92939290841561261a575f61240491604051809381927f127dd0520000000000000000000000000000000000000000000000000000000083527f0000000000000000000000000000000000000000000000000000000000000000896004850161233b565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115612615575f916125f3575b50925b61244c84612c79565b61245584612e26565b94156125e4576124658184613255565b915b6125d3575b5050506124788261022c565b81612588576125176125006020604060a085019485516124a1848201516001600160401b031690565b6001600160401b036124ce6124c26003546001600160401b039060401c1690565b6001600160401b031690565b91161161251b575b5001516040516124ed81610aff858201948561235d565b519020935101516001600160401b031690565b6001600160401b03165f52600560205260405f2090565b5590565b612582906020906001600160401b038151166001600160401b0319600354161760035501517fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6fffffffffffffffff00000000000000006003549260401b16911617600355565b5f6124d6565b506125928161022c565b600181036125bc57610a206801000000000000000068ff0000000000000000196004541617600455565b6125c58161022c565b60028103610a205750600290565b6125dc926133df565b5f808061246c565b6125ed8161306d565b91612467565b61260f91503d805f833e6126078183610c6e565b810190611b91565b5f612440565b612330565b506040517fccd771d60000000000000000000000000000000000000000000000000000000081525f8180612651876004830161231f565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115612615575f91612694575b5092612443565b6126a891503d805f833e6126078183610c6e565b5f61268d565b604051906126bd602083610c6e565b5f808352366020840137565b906126d382610f6f565b6126e06040519182610c6e565b82815280926126f1601f1991610f6f565b0190602036910137565b634e487b7160e01b5f52603260045260245ffd5b80518210156127235760209160051b010190565b6126fb565b906127328261371d565b6001600160a01b03811615612858576127649061275e61275960085461ffff9060e01c1690565b6137b5565b906137da565b9060206127718385613893565b019261279161278c612785865161ffff1690565b61ffff1690565b6126c9565b936127a461278c612785835161ffff1690565b936127b761278c612785845161ffff1690565b9260095491600a54935f5b886127d2612785845161ffff1690565b8b61ffff84169182101561284c576001926128378380612830612828828f8f8f8f8f61ffff9f9d9061281a826128459f6128229461280f9161270f565b9063ffffffff169052565b5161ffff1690565b916139d5565b92909561270f565b528c61270f565b906001600160401b03169052565b01166127c2565b50505050505050505090565b5090506128636126ae565b61286b6126ae565b916128746126ae565b91929190565b903590601e19813603018212156101c557018035906001600160401b0382116101c5576020019181360383136101c557565b903590601e19813603018212156101c557018035906001600160401b0382116101c557602001918160051b360383136101c557565b634e487b7160e01b5f52601160045260245ffd5b5f1981019190821161290357565b6128e1565b9190820391821161290357565b9061291f82610cde565b61292c6040519182610c6e565b82815280926126f1601f1991610cde565b90821015612723576129549160051b81019061287a565b9091565b805191908290602001825e015f815290565b61298c610c9e92949360208660405197889583870137840101905f8252612958565b03601f198101845283610c6e565b6129a35f610cde565b6129b06040519182610c6e565b5f8152601f196129bf5f610cde565b013660208301376129cf83610f6f565b926129dd6040519485610c6e565b808452601f196129ec82610f6f565b015f5b818110612a455750505f5b818110612a08575050505090565b80612a29612a2385612a1d600195878a61293d565b9061296a565b30613d1e565b612a33828861270f565b52612a3e818761270f565b50016129fa565b8060606020809389010152016129ef565b90600182811c92168015612a84575b6020831014612a7057565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612a65565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff1615612ac657565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260ff612b243360405f20906001600160a01b03165f5260205260405f2090565b541615612b2e5750565b63e2517d3f60e01b5f523360045260245260445ffd5b5f6040820192835192610140845151015194612b6860406020840151015195613d62565b91612b7287613d9a565b909182612c2b576101a00180515190929015612c0e575050612b95905187613eec565b94600192612ba1613de3565b6020845101525b612bf95782158080612bf0575b15612bce57505051606060208201519101525b93929190565b612bda575b5050612bc8565b6060612be99251015190613e80565b5f80612bd3565b50878214612bb5565b506060612c04613de3565b9151015293929190565b939550959050612c2360208351015188613e80565b600194612ba8565b50959092612c37613de3565b602084510152612ba8565b15612c4b575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b610c9e90612dac8151612cc3612ca06001600160801b03606086015116633b9aca00900490565b612cae814242821115614172565b42610708612cbc8342612908565b11156141a9565b612cf68151805160208201207f0000000000000000000000000000000000000000000000000000000000000000146141e0565b612d3e6020820151612d09815160ff1690565b6002549160ff831660ff811660ff8416149384612def575b612d389060209060081c60ff165b93015160ff1690565b936142df565b612d98612d8b6080612d57606085015163ffffffff1690565b93612d8060045495612d6c8763ffffffff1690565b63ffffffff811663ffffffff83161461433b565b015163ffffffff1690565b9160201c63ffffffff1690565b63ffffffff811663ffffffff83161461437c565b612de7612de26020608081850151604051612dce81610aff868201948561235d565b51902094015101516001600160401b031690565b6143bd565b808214612c42565b9350612d386020612d2f612e068286015160ff1690565b60ff612e1960088a901c82165b60ff1690565b9116149692505050612d21565b612e41612500602060a084015101516001600160401b031690565b5480612e4d5750505f90565b60408201908151604051612e6981610aff60208201948561235d565b5190201491821592612e87575b505015612e8257600190565b600290565b6001600160801b03919250612ebd612eae6020612ec99301516001600160801b0390511690565b9351516001600160801b031690565b6001600160801b031690565b911610155f80612e76565b15612edb57565b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b906001600160401b03809116911601906001600160401b03821161290357565b15612f2b5750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15612f655750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15612f9f5750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b9060038202918083046003149015171561290357565b908160011b918083046002149015171561290357565b90602c820291808304602c149015171561290357565b908160051b918083046020149015171561290357565b15613032575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b905f91602060408201510151518051936101008301926130bd84515161309b61278560e085015161ffff1690565b8091149081613245575b81613235575b81613225575b81613215575b50612ed4565b5f5b8681106131f857505f92839283805b8751518710156131b657908992916130fb6130f76130f18a6101808a015161270f565b51151590565b1590565b6131ab576001916131889161311e6131148b8d5161270f565b5163ffffffff1690565b80996131338263ffffffff8116998a10612f23565b613193575b61316291506020613149888a61270f565b51015161315b8c6101208c015161270f565b5114612f97565b6131826040613173859a978961270f565b5101516001600160401b031690565b90612f03565b965b019590916130ce565b63ffffffff6131a492168711612f5d565b5f88613138565b92509560019061318a565b5091509650610c9e9450869193506131f392506131db6001600160401b038216612fd1565b6131ed6001600160401b038416612fe7565b10613029565b61455f565b9161320e6001916131826040613173878961270f565b92016130bf565b905061018083015151145f6130b7565b61016084015151811491506130b1565b61014084015151811491506130ab565b61012084015151811491506130a5565b9061010081019061327482515161309b61278560e085015161ffff1690565b61327d8361371d565b926001600160a01b038416156133cc5793926132a86008549161275e6127598461ffff9060e01c1690565b946132c66132b68783613893565b9260a01c6001600160401b031690565b955f945f5f9860095497600a549260205f98019b5b8551518910156133a55791898b94928e989796946133046130f76130f18e61018087015161270f565b613397576133166131148d8a5161270f565b8092613379575b8c91508b90876001998c839e516133359061ffff1690565b9061333f956139d5565b9190936101200151906133519161270f565b51149061335d91612f97565b61336691612f03565b976001905b0197919394959290926132db565b63ffffffff613390921663ffffffff821611612f5d565b5f8161331d565b98509450509760019061336b565b5050935098505050610c9e94506131f392508691506131db6001600160401b038216612fd1565b633a517eed60e21b5f5260045260245b5ffd5b91906001600160a01b036133f28461371d565b161515806135e3575b6135de576040015160200151518051938415612edb5760b485116135c557613422856126c9565b925f5b8351811015613469578061345860206134406001948861270f565b5101516134526040613173858a61270f565b906153f5565b613462828861270f565b5201613425565b5092613505906134e45f9796939661349f8961348c61348784612fe7565b6128f5565b9483613497876126c9565b9c8d9261541f565b506134de6134ce6134c96134ba6134b585612ffd565b613738565b6134c387613013565b906137a8565b612915565b976134d8896154bc565b88615504565b86615592565b6134ff6134f86134f385615b8a565b61591e565b86602e0152565b846155a2565b5f9460305b835187101561357d576135756001916135706135268a8861270f565b5161353863ffffffff8c16848b6154e0565b61355961354484613746565b60408301516001600160401b0316908b61554a565b602061356484613754565b91015190890160200152565b613762565b96019561350a565b95509150925f945b82518610156135b8576135b06001916135ab6135a1898761270f565b5187830160200152565b613770565b950194613585565b50935050610c9e916146ea565b63156f758160e31b5f52600485905260b460245260445ffd5b505050565b5061ffff60085460e01c1615156133fb565b805f525f60205260ff61361c8360405f20906001600160a01b03165f5260205260405f2090565b541661368757805f525f6020526136478260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916600117905533916001600160a01b0316907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260ff6136b48360405f20906001600160a01b03165f5260205260405f2090565b54161561368757805f525f6020526136e08260405f20906001600160a01b03165f5260205260405f2090565b805460ff1916905533916001600160a01b0316907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b60065403613734576001600160a01b036008541690565b5f90565b603001908160301161290357565b906004820180921161290357565b90600c820180921161290357565b90602c820180921161290357565b906020820180921161290357565b600101908160011161290357565b602401908160241161290357565b906001820180921161290357565b9190820180921161290357565b61ffff16602c810290808204602c149015171561290357603001806030116129035790565b9190823b6001811115613828575f19810190811161290357811161380c57600161380382612915565b9360208501903c565b6001600160a01b0383633610565160e21b5f521660045260245ffd5b6001600160a01b0384633610565160e21b5f521660045260245ffd5b6040519061385182610c1d565b5f6060838281528260208201528260408201520152565b60011b906201fffe61fffe83169216820361290357565b61ffff5f199116019061ffff821161290357565b91909161389e613844565b506030835110613957576356414c34602084015160e01c0361395757602483015160c01c92602c81015160f01c93602e82015190604e83015160f01c916138f56138e6610cc0565b6001600160401b039093168352565b6139076020830197889061ffff169052565b604082015261391e6060820192839061ffff169052565b9461392b815161ffff1690565b9161ffff83169283159384156139ca575b508315613999575b508215613969575b505090506139575750565b634724a0fd60e01b5f5260045260245ffd5b6130f7925061398b6139826139919551935161ffff1690565b915161ffff1690565b916148e0565b805f8061394c565b90925061ffff6139bf6127856139ba6139b4875161ffff1690565b94613868565b61387f565b91161415915f613944565b60b41093505f61393c565b939092959461ffff63ffffffff82169316831015613aa25761ffff16916139fe6134b582612ffd565b9481613a166020888801015160e01c63ffffffff1690565b0361395757506001901b908116613a7757613a3f613a3385613746565b84016020015160c01c90565b9516613a5a5750613a52610a2092613754565b016020015190565b9050613a73915061ffff165f52600b60205260405f2090565b5490565b613a9d613a908361ffff165f52600c60205260405f2090565b546001600160401b031690565b613a3f565b6303e07d4560e61b5f52600485905263ffffffff1660245260445ffd5b9095979198929996613ad081611c40565b80613cda5750906020899392848015159081613ccd575b613af09161492f565b0196613b0f613afe8961496d565b613b08368c610ee0565b9088615631565b5f905f5b88868210613c2c575b505050613b299350614b14565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016803b156101c557835f8094613b998a956040519c8d97889687957f87d3a9b100000000000000000000000000000000000000000000000000000000875260048701614ee1565b03925af193841561261557610a2095613bc595613c12575b5060018111613bd9575b5050503690610ee0565b6001600160801b03633b9aca009151160490565b6001600160801b03613c02613bf0613c0a9561496d565b93613bfa87614f9b565b933691614fa5565b911691615775565b5f8080613bbb565b80613c205f613c2693610c6e565b806104a1565b5f613bb1565b6130f7613c54613c4d613c47858b613c659699979899614977565b806128ac565b3691614999565b613c5f368989614999565b90615703565b613cc3575090613c8c6106bc613c82613ca294613b29988c614977565b602081019061287a565b90815181518082149182613cad575b5050614a04565b600189935f88613b1c565b9091506020840120906020830120145f80613c9b565b9190600101613b13565b61ffff8111159150613ae7565b80613ce76133dc92611c40565b613cf081611c40565b7f112d89cc000000000000000000000000000000000000000000000000000000005f5260ff16600452602490565b5f80610a2093602081519101845af43d15613d5a573d91613d3e83610cde565b92613d4c6040519485610c6e565b83523d5f602085013e615070565b606091615070565b60016001600160401b03602060408281865151015116940151015116016001600160401b038111612903576001600160401b03161490565b613dab6001600160a01b039161371d565b1615613db957600754600191565b5f905f90565b60405190613dcc82610c1d565b5f6060838181528260208201528260408201520152565b60405190613df082610c1d565b606082525f6020830152613e02613dbf565b60408301525f60608301528160209060405191613e20602084610c6e565b5f808452805b818110613e335750505052565b8290613e3d613dbf565b82828801015201613e26565b15613e52575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b908051518015613ec05760b48111613ea9575090613ea0610c9e926150fc565b90808214613e49565b63156f758160e31b5f5260045260b460245260445ffd5b82633a517eed60e21b5f5260045260245ffd5b156139575750565b9060108110156127235760051b0190565b81519291613ef98461371d565b936001600160a01b038516156133cc5750613f16613fd994615194565b613f21818351613893565b90826020830192613f59613f37855161ffff1690565b61ffff613f4e61278560085461ffff9060e01c1690565b911614835190613ed3565b613f6a612e13602084015160ff1690565b9283151580614167575b613f849084939295945190613ed3565b613fe2600954998a96600a549781613fab8a83613fa482965161ffff1690565b8c8b6151b6565b9b9091613fc48d6001600160401b038451911115613ed3565b8a613fd38351925161ffff1690565b91615306565b88808214613e49565b5f935b83851061406b5750505050506001600160401b03839261401f61405d9361401a6140146134f3610c9e99615b8a565b99600955565b600a55565b167fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6008549260a01b16911617600855565b61406684600755565b600655565b909192939860806001918b6001600160401b039887614094612785613114856040850151613edb565b916140c96140be6140b186600161ffff88161b998a960151613edb565b516001600160401b031690565b9460608c0151613edb565b51936140dd8b518b8b61ffff881692615371565b9d166001600160401b038216145f1461412b5750901916995b820361410d575050901916995b0193929190613fe5565b6141239061ffff165f52600b60205260405f2090565b551799614103565b614160906141458561ffff165f52600c60205260405f2090565b906001600160401b03166001600160401b0319825416179055565b17996140f6565b506010841115613f74565b1561417b575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b156141b2575050565b7f12e74add000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b156141e85750565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f60015461422481612a56565b908160448501526001811690815f146142bb575060011461425b575b506142579192600319848303016024850152610837565b0390fd5b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b8183106142a15750919291508101606401614257614240565b805460648488010152859450602090920191600101614288565b60ff191660648086019190915291151560051b840190910191506142579050614240565b9392919093156142ef5750505050565b9160ff8092816084969581604051977fd382033a000000000000000000000000000000000000000000000000000000008952166004880152166024860152166044840152166064820152fd5b15614344575050565b9063ffffffff80927f73b1bde3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b15614385575050565b9063ffffffff80927f1eb5c1eb000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6001600160401b03165f52600560205260405f205480156143db5790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b908160209103126101c55751610a2081610dfa565b9060c060a0610a20936001600160401b0381511684526001600160401b0360208201511660208501526040810151604085015263ffffffff6060820151166060850152608081015160808501520151918160a08201520190610837565b9795939199989694929061ffff168852602088015f905b6008821061451a575050610a20989950926144de61450b9695936144ca6144ed946144bf6144fc986101208e01906120c5565b6101608c01906120c5565b6102406101a08b01526102408a01906103d6565b908882036101c08a0152610409565b908682036101e088015261039d565b908482036102008601526120ec565b91610220818403910152614418565b6020806001928e518152019c019101909a61448c565b1561453757565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b60206040820151518181015161457c81516001600160401b031690565b916145ff604061459e6145958786015163ffffffff1690565b63ffffffff1690565b930151858151910151906145ed87806145bc8563ffffffff90511690565b94015195510151956145de6145cf610ccf565b6001600160401b039099168952565b6001600160401b031687890152565b604086015263ffffffff166060850152565b608083015260a082015260e083015161ffff1661467560808501519260a08601519560c08101519061012081015161014082015190610180610160840151930151936040519a8b998a997f4cc22bb7000000000000000000000000000000000000000000000000000000008b5260048b01614475565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af1801561261557610c9e915f916146bb575b50614530565b6146dd915060203d6020116146e3575b6146d58183610c6e565b810190614403565b5f6146b5565b503d6146cb565b6146f48282613893565b91604051905f60208301526147108261298c6021820184612958565b6160008251116148b45750610aff61476a614758614730845161ffff1690565b60f01b7fffff0000000000000000000000000000000000000000000000000000000000001690565b926040519283916020830195866155b2565b51905ff0906001600160a01b0382161561488c5761487a926147c360209261406661482a956001600160a01b03167fffffffffffffffffffffffff00000000000000000000000000000000000000006008541617600855565b6147d06040820151600755565b6148216147e482516001600160401b031690565b7fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b6008549260a01b16911617600855565b015161ffff1690565b7fffff0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff7dffff000000000000000000000000000000000000000000000000000000006008549260e01b16911617600855565b6148835f600955565b610c9e5f600a55565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b906148ea906137b5565b908181149283156148fc575b50505090565b90919250621fffe061ffff82169160051b16908082046020149015171561290357820180921161290357145f80806148f6565b156149375750565b7fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b35610a2081610d8f565b91908110156127235760051b81013590603e19813603018212156101c5570190565b9291906149a581610f6f565b936149b36040519586610c6e565b602085838152019160051b8101918383116101c55781905b8382106149d9575050505050565b81356001600160401b0381116101c5576020916149f98784938701610d2f565b8152019101906149cb565b91909115614a10575050565b90614257614a52926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610837565b83810360031901602485015290610837565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156101c55701602081359101916001600160401b0382116101c55781360383136101c557565b90602083828152019260208260051b82010193835f925b848410614adc5750505050505090565b909192939495602080614b04600193601f19868203018852614afe8b88614a84565b90614a64565b9801940194019294939190614acc565b91909115614b20575050565b6142576040519283927ffef760c7000000000000000000000000000000000000000000000000000000008452602060048501526024840191614ab5565b9035601e19823603018112156101c55701602081359101916001600160401b0382116101c5578160051b360383136101c557565b9035607e19823603018112156101c5570190565b359060038210156101c557565b9035605e19823603018112156101c5570190565b90602083828152019260208260051b82010193835f925b848410614bed5750505050505090565b909192939495602080614c5f600193601f19868203018852614c0f8b88614bb2565b90614c1982614ba5565b614c228161022c565b8152614c51614c46614c3686850185614a84565b6060888601526060850191614a64565b926040810190614a84565b916040818503910152614a64565b9801940194019294939190614bdd565b610a2091614d39614d2e614cb2614c97614c898680614a84565b608087526080870191614a64565b614ca46020870187614a84565b908683036020880152614a64565b6080614d1e614cc46040880188614b91565b8684036040880152614cd581614ba5565b614cde8161022c565b8452614cec60208201614ba5565b614cf58161022c565b6020850152614d0660408201614ba5565b614d0f8161022c565b60408501526060810190614a84565b9190928160608201520191614a64565b926060810190614b5d565b916060818503910152614bc6565b90602083828152019260208260051b82010193835f925b848410614d6e5750505050505090565b909192939495601f198282030184528635601e19843603018112156101c557830190614d9e602082019280614b5d565b8091936020845252604082019060408160051b8401019380935f915b838310614ddd575050505050506020806001929801940194019294939190614d5e565b909192939495603f19838203018652614df68783614bb2565b8035614e0181610e0f565b614e0a81611c40565b8252614e2d614e1c6020830183614b91565b606060208501526060840190614c6f565b906040810135609e19823603018112156101c5576001936020938493614ed39301916040818303910152614ec5614ea5614e78614e6a8580614a84565b60a0865260a0860191614a64565b86850135614e8581610dfa565b151587850152614e986040860186614b91565b8482036040860152614c6f565b926060810135614eb481610dfa565b151560608401526080810190614b91565b906080818403910152614c6f565b980196019493019190614dba565b9092809492969593606083019083526060602084015252608081019560808560051b83010194815f90603e1981360301995b838310614f32575050505050610a209495506040818503910152614d47565b9091929397607f1986820301825288358b8112156101c5576020614f8c6001938683940190614f7f614f75614f678480614b5d565b604085526040850191614ab5565b9285810190614a84565b9185818503910152614a64565b9a019201930191909392614f13565b35610a2081610ec4565b92919092614fb284610f6f565b93614fc06040519586610c6e565b602085828152019060051b8201918383116101c55780915b838310614fe6575050505050565b82356001600160401b0381116101c55782016040818703126101c5576040519161500f83610be2565b81356001600160401b0381116101c557820187601f820112156101c5578781602061503c93359101614999565b83526020820135926001600160401b0384116101c55761506188602095869501610d2f565b83820152815201920191614fd8565b906150ad575080511561508557602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b815115806150f3575b6150be575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b156150b6565b8051519081156136875761510f826126c9565b915f5b8251805182101561518657906151756134f360206151328460019661270f565b5101516151706001600160401b03604061514d878b5161270f565b510151166040519261515e84610be2565b83526001600160401b03166020830152565b615846565b61517f828761270f565b5201615112565b505090505f610a209261596e565b90813b600181111561380c575f19810190811161290357600161380382612915565b919490929360208301946151d161278c612e13885160ff1690565b966151eb6124c26008546001600160401b039060a01c1690565b955f945f5b6151fe612e138b5160ff1690565b8110156152fa576152166131148260408b0151613edb565b9663ffffffff88169061522f8961ffff88168410612f23565b826152e0575b505082878787878c5194615248956139d5565b908260808b01519061525991613edb565b516001600160401b03169a8360608c01519061527491613edb565b51926001600160401b038d16926001600160401b03169084848314918215926152d5575b50508c516152a591613ed3565b6152ae91612908565b906152b8916137a8565b996152c2916153f5565b6152cc828d61270f565b526001016151f0565b14159050845f615298565b6152f39163ffffffff8b51921610613ed3565b5f80615235565b50969750505050505050565b94939190604051956101008701948786106001600160401b03871117610bfd57610a20985f9761ffff89989660ff966020968b996040528d52868d015216968760408c01528360608c015260808b01528060a08b01528160c08b01521760e0890152015116946159d0565b9091602061ffff919594950151169363ffffffff8116948510156153d7575061539c6134b585612ffd565b936020858401015160e01c0361395757506020906153d16153cb6153bf86613746565b83016020015160c01c90565b94613754565b01015190565b6303e07d4560e61b5f5260049190915263ffffffff1660245260445ffd5b6134f3906001600160401b03610a20936040519261541284610be2565b8352166020820152615846565b9190949394808203828111612903576001811461549d5761543f90615bc2565b928382019283831161290357600188019283891161290357838786615464938661541f565b61547361348760019297612fe7565b8901018093116129035761548893869261541f565b61549590611c5492615cfc565b93849261270f565b506154af91506154b89295949561270f565b5192839261270f565b5290565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b602d908260081c602c8201530153565b604f908260081c604e8201530153565b600a907fffff000000000000000000000000000000000000000000000000000000000000610a2094937f610000000000000000000000000000000000000000000000000000000000000083521660018201527f80600a3d393df30000000000000000000000000000000000000000000000000060038201520190612958565b9061568090612de760405160208101906156688288604080916001600160801b038151168452602081015160208501520151910152565b60608152615677608082610c6e565b519020916143bd565b60208201518082036156d55750506001600160801b03633b9aca00915116046156ad814242821115614172565b4203428111612903576001600160801b03610c9e911663ffffffff6004541690818110615bde565b7f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b908151815103613687575f5b825181101561576d57615722818461270f565b515161572e828461270f565b5151036157665761573f818461270f565b5160208151910120615751828461270f565b5160208151910120036157665760010161570f565b5050505f90565b505050600190565b929190925f5b845181101561583f5761578e818661270f565b51908360405160208101906001600160401b038616825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801915f5b818110615808575050505081610aff60019760206157fe940151605f19848303016080850152610837565b5190205d0161577b565b91939496509194969760208061582a60019360bf198b82030188528951610837565b970194019101918a96949392989795986157d3565b5050509050565b602460208201906001600160401b03825116806158e3575b5061586b61589991612915565b9261589061588a61588461587e87615c22565b87615c2e565b86615c45565b85615c5c565b90519084615c89565b906158ae6124c282516001600160401b031690565b6158b757505090565b6158d86124c26158ca6158df9486615c72565b92516001600160401b031690565b9083615ca1565b5090565b600191505b6080811015615910575061586b6159096159046158999361377e565b61378c565b915061585e565b60019060071c9101906158e8565b8051600181018091116129035761593490612915565b8051156127235761595e816159516020945f868196015382615cda565b5060405191828092612958565b039060025afa15612615575f5190565b9181810381811161290357600181146159b35761598a90615bc2565b820191828111612903578261599f918561596e565b916159aa929361596e565b610a2091615cfc565b50506159be9161270f565b5190565b5f1981146129035760010190565b9392909491828414808091615b72575b615b4a5761ffff8311612edb576159f78783612908565b9060018214615aa05750615a0a90615bc2565b93615a37615a1886896137a8565b956134c3613487615a31615a2b8861379a565b9761379a565b92612fe7565b92815b85811087828a83615a79575b50505015615a5c57615a57906159c2565b615a3a565b9495969791859188615a6e948b6159d0565b946159aa95966159d0565b63ffffffff929350615a9691604060606131149301510151613edb565b161087828a615a46565b925050509493929415615aec57505082615ae791610a20939451916020810151615acf604083015161ffff1690565b9063ffffffff60c060a08501519401519416946139d5565b6153f5565b615af782939261379a565b1490811591615b28575b50615b145760806159be9293015161270f565b8251634724a0fd60e01b5f5260045260245ffd5b9050615b4261459561311484604060608901510151613edb565b14155f615b01565b505092915050610a209250805190615b6c6040602083015192015161ffff1690565b91615d70565b50615b856130f7838960e08a0151615d4e565b6159e0565b60405190606090615b9b8284610c6e565b602283526158df91600a906020850190601f19013682375360206021840153600283615c89565b9060015b8060011b9083821015615bd95750615bc6565b925050565b15615be7575050565b906001600160801b0380927fdb020436000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b6020600a910153600190565b602082602292010153600181018091116129035790565b602082600a92010153600181018091116129035790565b6020828192010153600181018091116129035790565b602082601092010153600181018091116129035790565b81602091939293010152602081018091116129035790565b9092919083016020015b6080821015615cbf57906001929391530190565b600180916080607f85161781530193019060071c9092615cab565b90805191821561576d576021602084930191015e600101806001116129035790565b604051615d0a608082610c6e565b60418152615d186041610cde565b602082019190601f1901368337805115612723576020935f93600161595e94536021830152604182015260405191828092612958565b91818103908111612903576001901b905f198201918211612903571b16151590565b909161ffff1692602c840293808504602c149015171561290357836030019384603011612903578160051b91808304602014901517156129035701926030840180911161290357615dc090613770565b825110613957575001605001519056fea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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
