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
	Bin: "0x61016080604052346105585761645c803803809161001d828561067a565b8339810160e082820312610558576100348261069d565b6100406020840161069d565b61004c6040850161069d565b916100596060860161069d565b60808601519094906001600160401b0381116105585786019080601f8301121561055857815161008b926020016106b1565b9461009d60c060a0830151920161069d565b957f1a863411baffac77322e211abff616fd73c59fe2bb867e61a77bb762a1a612336101005280518101906020820190602081840312610558576020810151906001600160401b038211610558570191829003601f1981019061012013610558576040519160e083016001600160401b0381118482101761064b5760405260208401516001600160401b038111610558576020908501019080601f8301121561055857815161014e926020016106b1565b8252604081126105585760408051916101668361065f565b6101718286016106f6565b835261017f606086016106f6565b602084015260208401928352603f190112610558576040516101a08161065f565b6101ac60808501610704565b81526101ba60a08501610704565b6020820152604083019081526101d260c08501610718565b90606084019182526101e660e08601610718565b92608085019384526101008601519586151587036105585760a0860196875261012001519460028610156105585760c08101958652518051906001600160401b03821161064b57610238600154610729565b601f81116105fb575b50602090601f831160011461058e5763ffffffff95949392915f9183610583575b50508160011b915f199060031b1c1916176001555b5160ff81511661ff00602060025493015160081b169161ffff191617176002555160018060401b0381511660035491602068010000000000000000600160801b0391015160401b169160018060801b031916171760035551169267ffffffff00000000600454925160201b169051151560401b925191600283101561056f5769ff00000000000000000068ff00000000000000009360481b169469ff000000000000000000199160018060481b031916171617911617176004556040516103488161034181610761565b038261067a565b602081519101206101205260405163685e272760e11b8152602060048201526020818061037760248201610761565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4908115610564575f9161052e575b50610140526040516103b38161034181610761565b6001600160401b03906020906103c890610812565b0151600354916001600160401b0383169116818103610519575050604090811c6001600160401b03165f908152600560205220556001600160a01b0390811660805290811660a05290811660c0521660e05260045463ffffffff8082169061070882019081116105055763ffffffff809360201c1692839116116104f057826001600160a01b0381166104c9575061045e610ae9565b5061046b61010051610b6b565b505b6040516157629081610c3a823960805181614e15015260a051816142dc015260c05181505060e0518181816121690152612b0301526101005181818161026f015261030a015261012051816121d7015261014051816121340152f35b806104d66104ea926109f3565b506104e081610a69565b5061010051610bc4565b5061046d565b6333fae18560e21b5f5260045260245260445ffd5b634e487b7160e01b5f52601160045260245ffd5b6355bace6f60e01b5f5260045260245260445ffd5b90506020813d60201161055c575b816105496020938361067a565b8101031261055857515f61039e565b5f80fd5b3d915061053c565b6040513d5f823e3d90fd5b634e487b7160e01b5f52602160045260245ffd5b015190505f80610262565b90601f1983169160015f52815f20925f5b8181106105e3575091600193918563ffffffff9998979694106105cb575b505050811b01600155610277565b01515f1960f88460031b161c191690555f80806105bd565b9293602060018192878601518155019501930161059f565b60015f525f51602061641c5f395f51905f52601f840160051c81019160208510610641575b601f0160051c01905b8181106106365750610241565b5f8155600101610629565b9091508190610620565b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b0382111761064b57604052565b601f909101601f19168101906001600160401b0382119082101761064b57604052565b51906001600160a01b038216820361055857565b9192916001600160401b03821161064b57604051916106da601f8201601f19166020018461067a565b829481845281830111610558578281602093845f96015e010152565b519060ff8216820361055857565b51906001600160401b038216820361055857565b519063ffffffff8216820361055857565b90600182811c92168015610757575b602083101461074357565b634e487b7160e01b5f52602260045260245ffd5b91607f1691610738565b6001545f929161077082610729565b80825291600181169081156107d1575060011461078b575050565b60015f9081529293509091905f51602061641c5f395f51905f525b8383106107b7575060209250010190565b6001816020929493945483858701015201910191906107a6565b9050602093945060ff929192191683830152151560051b010190565b9081518110156107fe570160200190565b634e487b7160e01b5f52603260045260245ffd5b60405161081e8161065f565b606081525f602082015250805180159081156109e7575b506109d8575f19908051805b610994575b505f19821461098657600182019081831161050557600360fc1b6001600160f81b031961087384846107ed565b51161480610971575b610961575f5b81518310156109115761089583836107ed565b5160f81c603081108015610907575b6108f557600a82026001600160401b03908116602f1990920160ff16919091018116911681106108d957600190920191610882565b50915050604051906108ea8261065f565b81525f602082015290565b5050915050604051906108ea8261065f565b50603981116108a4565b9150916001811190811591610955575b5061094657604051916109338361065f565b82526001600160401b0316602082015290565b63384c687760e21b5f5260045ffd5b602b915010155f610921565b915050604051906108ea8261065f565b5080518381039081116105055760021061087c565b60405191506108ea8261065f565b5f19810181811161050557602d60f81b6001600160f81b03196109b783866107ed565b5116146109ce57508015610505575f190180610841565b92505f9050610846565b6329120bff60e21b5f5260045ffd5b6040915010155f610835565b6001600160a01b0381165f9081525f51602061643c5f395f51905f52602052604090205460ff16610a64576001600160a01b03165f8181525f51602061643c5f395f51905f5260205260408120805460ff191660011790553391905f51602061639c5f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f5160206163bc5f395f51905f52602052604090205460ff16610a64576001600160a01b03165f8181525f5160206163bc5f395f51905f5260205260408120805460ff191660011790553391905f5160206163fc5f395f51905f52905f51602061639c5f395f51905f529080a4600190565b5f80525f5160206163bc5f395f51905f526020525f5160206163dc5f395f51905f525460ff16610b67575f8080525f5160206163bc5f395f51905f526020525f5160206163dc5f395f51905f52805460ff1916600117905533905f5160206163fc5f395f51905f525f51602061639c5f395f51905f528280a4600190565b5f90565b805f525f60205260405f205f805260205260ff60405f205416155f14610a6457805f525f60205260405f205f805260205260405f20600160ff198254161790555f33915f51602061639c5f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16610c33575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291905f51602061639c5f395f51905f529080a4600190565b50505f9056fe610120806040526004361015610013575f80fd5b5f3560e01c90816301ffc9a714610a37575080630bece3561461099d578063248a9ca3146109735780632f2ff15d1461094457806336568abe146108f5578063536c2ad31461089d5780636a28f000146107f85780638a8e4c5d146107d957806391d148541461079d578063974a74c414610652578063a217fddf14610638578063a6f031bb14610528578063ac9650d814610363578063d547741f1461032d578063db3e1fa4146102f3578063ddba6537146102505763ef913a4b146100d8575f80fd5b3461024c575f36600319011261024c5760405160208082015261012060408201525f60015461010681612f11565b90816101608501526001811690815f1461022757506001146101c6575b6101c2836101ae818560ff600254818116606085015260081c1660808301526001600160401b0360035481811660a085015260401c1660c083015260ff60045463ffffffff811660e085015263ffffffff8160201c16610100850152818160401c16151561012085015260481c1661019a816112fa565b61014083015203601f198101835282610ccf565b604051918291602083526020830190610c10565b0390f35b91905060015f527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6915f905b80821061020b57509091508101610180016101ae610123565b91926001816020925461018085880101520191019092916101f2565b60ff19166101808086019190915291151560051b840190910191506101ae9050610123565b5f80fd5b3461024c5761025e36610ad5565b505060ff60045460401c166102cb577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f205416156102bc575b636d40ebe160e11b5f5260045ffd5b6102c590612fb8565b806102ad565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b3461024c575f36600319011261024c5760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b3461024c5761036161033e36610b42565b9061035c610357825f525f602052600160405f20015490565b612fb8565b613ec6565b005b3461024c57602036600319011261024c576004356001600160401b03811161024c573660238201121561024c5780600401356001600160401b03811161024c5760248201913660248360051b8301011161024c576020926040516103c78582610ccf565b5f815284810191601f1986013684376103df85610e74565b936103ed6040519586610ccf565b858552601f196103fc87610e74565b01875f5b828110610519575050505f5b868110156104bc576001906104985f808b8861046561043360248860051b8b01018b612e67565b90938d6040519483869484860198893784019083820190898252519283915e010185815203601f198101835282610ccf565b5190305af43d156104b4573d9061047b82610cf0565b916104896040519384610ccf565b82523d5f8d84013e5b30615217565b6104a28289612cfe565b526104ad8188612cfe565b500161040c565b606090610492565b604080518981528751818b018190525f92600582901b83018101918a8d01918d9085015b8287106104ed5785850386f35b909192938280610509600193603f198a82030186528851610c10565b96019201960195929190926104e0565b60608882018301528101610400565b3461024c57602036600319011261024c576004356001600160401b03811161024c578060040190610140600319823603011261024c5760ff60045460401c166102cb575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff161561062b575b6105ca6044820183612e99565b916105d86064820185612e99565b93909161010481013590600282101561024c5760209661062396610600610124840183612e99565b969095604051986106118c8b610ccf565b5f8a52608460a487019601359461421e565b604051908152f35b610633612f49565b6105bd565b3461024c575f36600319011261024c5760206040515f8152f35b3461024c57602036600319011261024c576004356001600160401b03811161024c578060040190610160600319823603011261024c5760ff60045460401c166102cb575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff1615610790575b61014481016106f68184612e67565b9050156107685761070a6044830184612e99565b90916107196064850186612e99565b95909461010481013591600283101561024c576020976106239761075196610758610748610124870186612e99565b99909886612e67565b3691610d0b565b98608460a487019601359461421e565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b610798612f49565b6106e7565b3461024c576107ab36610b42565b905f525f6020526001600160a01b0360405f2091165f52602052602060ff60405f2054166040519015158152f35b3461024c576107e736610ad5565b5050636d40ebe160e11b5f5260045ffd5b3461024c575f36600319011261024c57335f9081527fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5602052604090205460ff16156108865760045460ff8160401c161561085e5768ff00000000000000001916600455005b7fa5e1ca77000000000000000000000000000000000000000000000000000000005f5260045ffd5b63e2517d3f60e01b5f52336004525f60245260445ffd5b3461024c57602036600319011261024c576108d96108e76101c26108c2600435612d26565b919390604051958695606087526060870190610b68565b908582036020870152610ba1565b908382036040850152610bd4565b3461024c5761090336610b42565b336001600160a01b0382160361091c5761036191613ec6565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b3461024c5761036161095536610b42565b9061096e610357825f525f602052600160405f20015490565b613e39565b3461024c57602036600319011261024c5760206106236004355f525f602052600160405f20015490565b3461024c576109ab36610ad5565b9060ff60045460401c166102cb575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060209081527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47549092610a1992909160ff1615610a2a5761198d565b60405190610a2681610b24565b8152f35b610a32612f49565b61198d565b3461024c57602036600319011261024c57600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361024c57817f7965db0b0000000000000000000000000000000000000000000000000000000060209314908115610aab575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483610aa4565b90602060031983011261024c576004356001600160401b03811161024c578260238201121561024c578060040135926001600160401b03841161024c576024848301011161024c576024019190565b60031115610b2e57565b634e487b7160e01b5f52602160045260245ffd5b604090600319011261024c57600435906024356001600160a01b038116810361024c5790565b90602080835192838152019201905f5b818110610b855750505090565b825163ffffffff16845260209384019390920191600101610b78565b90602080835192838152019201905f5b818110610bbe5750505090565b8251845260209384019390920191600101610bb1565b90602080835192838152019201905f5b818110610bf15750505090565b82516001600160401b0316845260209384019390920191600101610be4565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b604081019081106001600160401b03821117610c4f57604052565b634e487b7160e01b5f52604160045260245ffd5b606081019081106001600160401b03821117610c4f57604052565b608081019081106001600160401b03821117610c4f57604052565b60c081019081106001600160401b03821117610c4f57604052565b60e081019081106001600160401b03821117610c4f57604052565b90601f801991011681019081106001600160401b03821117610c4f57604052565b6001600160401b038111610c4f57601f01601f191660200190565b929192610d1782610cf0565b91610d256040519384610ccf565b82948184528183011161024c578281602093845f960137010152565b9080601f8301121561024c57816020610d5c93359101610d0b565b90565b359060ff8216820361024c57565b35906001600160401b038216820361024c57565b919082604091031261024c57604051610d9981610c34565b6020610db2818395610daa81610d6d565b855201610d6d565b910152565b359063ffffffff8216820361024c57565b3590811515820361024c57565b35906001600160801b038216820361024c57565b919082606091031261024c57604051610e0181610c63565b6040808294610e0f81610dd5565b8452602081013560208501520135910152565b80929103916060831261024c57604051610e3b81610c34565b6040819483358352601f19011261024c576020906040805193610e5d85610c34565b610e68848201610db7565b85520135828401520152565b6001600160401b038111610c4f5760051b60200190565b919060808382031261024c5760405190610ea482610c7e565b819380356001600160401b03811161024c57606092610ec4918301610d41565b835260208101356020840152610edc60408201610d6d565b60408401520135908160070b820361024c5760600152565b91909160808184031261024c5760405190610f0e82610c7e565b819381356001600160401b03811161024c57820181601f8201121561024c578035610f3881610e74565b91610f466040519384610ccf565b81835260208084019260051b8201019184831161024c5760208201905b838210610fb457505050508352610f7c60208301610dc8565b60208401526040820135906001600160401b03821161024c5782610fa960609492610db294869401610e8b565b604086015201610d6d565b81356001600160401b03811161024c57602091610fd688848094880101610e8b565b815201910190610f63565b9080601f8301121561024c57604080519290610ffd9084610ccf565b82906040810192831161024c57905b8282106110195750505090565b813581526020918201910161100c565b9080601f8301121561024c57813561104081610e74565b9261104e6040519485610ccf565b81845260208085019260051b82010192831161024c57602001905b8282106110765750505090565b6020809161108384610db7565b815201910190611069565b519060ff8216820361024c57565b51906001600160401b038216820361024c57565b919082604091031261024c576040516110c881610c34565b6020610db28183956110d98161109c565b85520161109c565b519063ffffffff8216820361024c57565b5190811515820361024c57565b51906001600160801b038216820361024c57565b919082606091031261024c5760405161112b81610c63565b6040808294611139816110ff565b8452602081015160208501520151910152565b60208183031261024c578051906001600160401b03821161024c57016101808183031261024c576040519161118083610c99565b81516001600160401b03811161024c5782019182820392610120841261024c57604051936111ad85610cb4565b81516001600160401b03811161024c5782019084601f8301121561024c5781516111d681610cf0565b906111e46040519283610ccf565b808252866020828601011161024c576020815f9282604097018386015e830101528652601f19011261024c576101009060405161122081610c34565b61122c6020830161108e565b815261123a6040830161108e565b6020820152602086015261125184606083016110b0565b604086015261126260a082016110e1565b606086015261127360c082016110e1565b608086015261128460e082016110f2565b60a0860152015190600282101561024c57836101409260c06112f296015285526112b18360208301611113565b60208601526112c38360808301611113565b60408601526112d460e082016110ff565b60608601526112e78361010083016110b0565b6080860152016110b0565b60a082015290565b60021115610b2e57565b9060608061131b8451608085526080850190610c10565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b8282106113ad57505050506001600160401b0360606113a3819360208701511515602087015260408701518682036040880152611304565b9401511691015290565b909192936020806113ca600193609f198a82030186528851611304565b96019201920190929161136b565b905f905b600282106113e957505050565b60208060019285518152019301910190916113dc565b90602080835192838152019201905f5b81811061141c5750505090565b8251151584526020938401939092019160010161140f565b908151610940825260c06114578251610120610940860152610a60850190610c10565b9160ff6020808301518281511661096088015201511661098085015261149b60408201516109a08601906001600160401b0360208092828151168552015116910152565b63ffffffff6060820151166109e085015263ffffffff608082015116610a0085015260a08101511515610a2085015201516114d5816112fa565b610a4083015261150a60208401516020840190604080916001600160801b038151168452602081015160208501520151910152565b6040830151908281036080840152815160a0825260206116748251604060a086015261154f60e0860182516001600160401b0360208092828151168552015116910152565b61024061156d848301516102c06101208901526103a0880190610c10565b60408301516001600160401b031661014088015260608301516001600160801b03166101608801526080830151151561018088015260a083015180516101a0890152602090810151805163ffffffff166101c08a015201516101e08801529160c0810151151561020088015260e08101516102208801526101008101511515828801526101208101516102608801526101408101516102808801526101608101516102a08801526101808101516102c08801526101a08101516102e08801526101c081015115156103008801526101e08101516103208801526102008101511515610340880152610220810151610360880152015160df1986830301610380870152610c10565b91015190609f198382030160c0840152606060c08201926001600160401b03815116835263ffffffff60208201511660208401526116d46040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519160c060a0830152825180915260e0820190602060e08260051b8501019401925f905b8282106118f75750505050509060606117228493602061175a9601518482036020860152611344565b9261174a604082015160408501906001600160401b0360208092828151168552015116910152565b0151906080818403910152611344565b60608301516001600160801b031660a083015260808301515f60c084015b600882106118e15750505061182f61181b6118076117f36117df6101a0956117a960a08a01516101c08a01906113d8565b6117bc60c08a01516102008a01906113d8565b61ffff60e08a0151166102408901526101008901518882036102608a0152610b68565b610120880151878203610280890152610ba1565b6101408701518682036102a0880152610bd4565b6101608601518582036102c0870152610b68565b6101808501518482036102e08601526113ff565b9201518051610300830152602081015160ff1661032083015260408101515f61034084015b601082106118c55750505060608101515f61054084015b601082106118af5750505060800151905f90610740015b601082106118905750505090565b6020806001926001600160401b03865116815201930191019091611882565b602080600192855181520193019101909161186b565b60208060019263ffffffff865116815201930191019091611854565b6020806001928551815201930191019091611778565b909192939460df1982820301855285519081516004811015610b2e57611975826020600195819594829552015190604084820152606061194383516080604085015260c0840190610c10565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610c10565b97019501939201906116f9565b6040513d5f823e3d90fd5b919082019160208184031261024c578035906001600160401b03821161024c5701906109408284031261024c576040516101c081018181106001600160401b03821117610c4f5760405282356001600160401b03811161024c57830180850390610120821261024c5760405191611a0383610cb4565b8135906001600160401b03821161024c57611a22886040938501610d41565b8452601f19011261024c5761010090604051611a3d81610c34565b611a4960208301610d5f565b8152611a5760408301610d5f565b60208201526020840152611a6e8760608301610d81565b6040840152611a7f60a08201610db7565b6060840152611a9060c08201610db7565b6080840152611aa160e08201610dc8565b60a08401520135600281101561024c5760c08201528152611ac58460208501610de9565b602082015260808301356001600160401b03811161024c5783019360a08582031261024c5760405194611af786610c7e565b80356001600160401b03811161024c57810160408184031261024c5760405190611b2082610c34565b80356001600160401b03811161024c5781016102c08186031261024c576040519061026082018281106001600160401b03821117610c4f57604052611b658682610d81565b825260408101356001600160401b03811161024c5786611b86918301610d41565b6020830152611b9760608201610d6d565b6040830152611ba860808201610dd5565b6060830152611bb960a08201610dc8565b6080830152611bcb8660c08301610e22565b60a0830152611bdd6101208201610dc8565b60c083015261014081013560e0830152611bfa6101608201610dc8565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a0830152611c496102208201610dc8565b6101c08301526102408101356101e0830152611c686102608201610dc8565b6102008301526102808101356102208301526102a0810135906001600160401b03821161024c57611c9b91879101610d41565b61024082015282526020810135906001600160401b03821161024c570160c08185031261024c5760405190611ccf82610c7e565b611cd881610d6d565b8252611ce660208201610db7565b6020830152611cf88560408301610e22565b604083015260a0810135906001600160401b03821161024c570184601f8201121561024c57803590611d2982610e74565b91611d376040519384610ccf565b80835260208084019160051b8301019187831161024c5760208101915b838310612be7575050505060608201526020820152865260208101356001600160401b03811161024c5782611d8a918301610ef4565b6020870152611d9c8260408301610d81565b60408701526080810135906001600160401b03821161024c57611dc191839101610ef4565b606086015260408201948552611dd960a08501610dd5565b60608301528060df8501121561024c57604051611df861010082610ccf565b806101c086019183831161024c5790839160c08801905b848210612bd45750506080850152611e2691610fe1565b60a0830152611e39816102008601610fe1565b60c083015261024084013561ffff8116810361024c5760e083019081526102608501356001600160401b03811161024c5782611e76918701611029565b6101008401526102808501356001600160401b03811161024c57850182601f8201121561024c57803590611ea982610e74565b91611eb76040519384610ccf565b80835260208084019160051b8301019185831161024c57602001905b828210612bc4575050506101208401526102a08501356001600160401b03811161024c5785019180601f8401121561024c57823592611f1184610e74565b93611f1f6040519586610ccf565b80855260208086019160051b8301019183831161024c57602001905b828210612bac5750505061014084019283526102c08601356001600160401b03811161024c5781611f6d918801611029565b9561016085019687526102e08101356001600160401b03811161024c57810182601f8201121561024c57803590611fa382610e74565b91611fb16040519384610ccf565b80835260208084019160051b8301019185831161024c57602001905b828210612b94575050506101808601526106408183036102ff19011261024c576040519160a083018381106001600160401b03821117610c4f57604052610300820135835261201f6103208301610d5f565b60208401528061035f8301121561024c57610200916040516120418482610ccf565b8061054083019184831161024c576103408401905b838210612b7c57505060408601528261055f8301121561024c5760405161207d8582610ccf565b8061074084019285841161024c57905b838210612b6c57505060608601528261075f8301121561024c576120b46040519485610ccf565b61094084920192831161024c57905b828210612b545750505060808201526101a08401526120e183612ff8565b93919392909660a0519497885f14612ab9575f6040518080937f127dd05200000000000000000000000000000000000000000000000000000000825260606004830152612131606483018c611434565b907f00000000000000000000000000000000000000000000000000000000000000006024840152604483015203816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115611982575f91612a97575b50985b8951633b9aca006001600160801b0360608d01511604428111612a68576107086121c68242612edc565b11612a3957508051805160208201207f0000000000000000000000000000000000000000000000000000000000000000036129425750602081015160ff815116906002549160ff8316928382149283612929575b6020015160ff1692156128e1575050505063ffffffff606082015116906004549163ffffffff83168082036128b357505063ffffffff608081920151169160201c168082036128855750506122cd8a6122c56001600160401b0360206080818501516040516122a98482018093604080916001600160801b038151168452602081015160208501520151910152565b606081526122b78382610ccf565b519020940151015116614c00565b808214613675565b6122d68a6136ac565b9815612686575061ffff610100870151519251168092149283612676575b8361266a575b508261265e575b508161264e575b50156126265761231782613f49565b946001600160a01b038616156126135793906123496008979397549661234361ffff8960e01c16613f7f565b90613fa4565b96612354888261400e565b6080525f925f965f9560095494600a54965f985b6101008b0151518a101561243257908d979695949392918b61238f8c610180830151612cfe565b5115612426576123a98c61010063ffffffff930151612cfe565b5116809d61240e575b50508a60019c8b88828d8c829e6080516020015161ffff16906123d495614156565b9190936101200151906123e691612cfe565b5114906123f2916137ee565b6123fb9161375a565b986001905b019890919293949596612368565b63ffffffff61241f921681116137b4565b5f8c6123b2565b50975098600190612400565b50955096509750979195935097506001600160401b0382166003810290808204600314901517156125ff576801fffffffffffffffe82609f1c16906001600160401b038360a01c1682046002146001600160401b038460a01c161517156125ff576001600160401b03936124ad92858560a01c16921161383e565b6124b683614c46565b60a01c16915b6125ee575b5050506124cd82610b24565b816125a5576001600160401b036020604060a08401938451838101519085600354851c1686831611612551575b5050015160405161252b8382018093604080916001600160801b038151168452602081015160208501520151910152565b6060815261253a608082610ccf565b51902092510151165f52600560205260405f205590565b85905116851960035416176003557fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6fffffffffffffffff000000000000000060035492851b169116176003555f806124fa565b506125af81610b24565b600181036125d7576801000000000000000068ff000000000000000019600454161760045590565b6125e081610b24565b60028103610d5c5750600290565b6125f792613882565b5f80806124c1565b634e487b7160e01b5f52601160045260245ffd5b82633a517eed60e21b5f5260045260245ffd5b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b905061018084015151145f612308565b5151811491505f612301565b5151821492505f6122fa565b61012087015151831493506122f4565b97909691925f9594955060205f99510151519384519861ffff6101008a0151519251168092149283612875575b83612869575b508261285d575b508161284d575b5015612626575f9893985b87811061282557505f935f995f965f9b5b6101008a0151518d10156127b7576127008d6101808c0151612cfe565b51156127ad578a8a888f806101008401519061271b91612cfe565b5163ffffffff169c8d80968180978110906127359161377a565b6127669761275f956101209460209461275494612795575b5050612cfe565b510151930151612cfe565b51146137ee565b600161278b81986001600160401b0360406127818d8c612cfe565b510151169061375a565b9c5b019b966126e3565b63ffffffff6127a6921681116137b4565b5f8261274d565b969b60019061278d565b509650969492995096509691506001600160401b0381166003810290808204600314901517156125ff576001600160401b038316916801fffffffffffffffe8460011b1692808404600214901517156125ff576128169284921161383e565b61281f82614c46565b916124bc565b976128436001916001600160401b0360406127818d999e9989612cfe565b98019893986126d2565b905061018087015151145f6126c7565b5151811491505f6126c0565b5151821492505f6126b9565b6101208a015151831493506126b3565b7f1eb5c1eb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f73b1bde3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6084945060ff90604051947fd382033a000000000000000000000000000000000000000000000000000000008652600486015260081c16602484015260448301526064820152fd5b6020810151600883901c60ff908116911614935061221a565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f60015461297e81612f11565b908160448501526001811690815f14612a1557506001146129b5575b506129b19192600319848303016024850152610c10565b0390fd5b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b8183106129fb57509192915081016064016129b161299a565b8054606484880101528594506020909201916001016129e2565b60ff191660648086019190915291151560051b840190910191506129b1905061299a565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b7f20daedb6000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b612ab391503d805f833e612aab8183610ccf565b81019061114c565b5f612199565b506040517fccd771d6000000000000000000000000000000000000000000000000000000008152602060048201525f8180612af7602482018b611434565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115611982575f91612b3a575b509861219c565b612b4e91503d805f833e612aab8183610ccf565b5f612b33565b60208091612b6184610d6d565b8152019101906120c3565b813581526020918201910161208d565b60208091612b8984610db7565b815201910190612056565b60208091612ba184610dc8565b815201910190611fcd565b60208091612bb984610d6d565b815201910190611f3b565b8135815260209182019101611ed3565b8135815285935060209182019101611e0f565b82356001600160401b03811161024c578201906040828b03601f19011261024c5760405191612c1583610c34565b6020810135600481101561024c57835260408101356001600160401b03811161024c576020910101906080828c031261024c5760405192612c5584610c7e565b82356001600160401b03811161024c578c612c71918501610d41565b8452612c7f60208401610dd5565b6020850152612c9060408401610dc8565b60408501526060830135936001600160401b03851161024c57612cb88d602096879601610d41565b606082015283820152815201920191611d54565b90612cd682610e74565b612ce36040519182610ccf565b8281528092612cf4601f1991610e74565b0190602036910137565b8051821015612d125760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b90612d3082613f49565b6001600160a01b03811615612e1d57612d559061234361ffff60085460e01c16613f7f565b6020612d61828561400e565b0191612d7161ffff845116612ccc565b93612d8061ffff855116612ccc565b938493612d9161ffff835116612ccc565b9260095491600a54935f5b61ffff8251168b61ffff831691821015612e1157996001600160401b03612dfd8380612df68c9d9e9f828d9e8d9e8d9e8d9e61ffff9e8f9060019f9e612de581612dee9a612cfe565b52511691614156565b929096612cfe565b528c612cfe565b911690520116908897969594939291612d9c565b50505050505050509150565b509050602060405191612e308284610ccf565b5f83525f36813760405192612e458385610ccf565b5f84525f36813760405192612e5a8185610ccf565b5f8452505f368137929190565b903590601e198136030182121561024c57018035906001600160401b03821161024c5760200191813603831361024c57565b903590601e198136030182121561024c57018035906001600160401b03821161024c57602001918160051b3603831361024c57565b5f198101919082116125ff57565b919082039182116125ff57565b90612ef382610cf0565b612f006040519182610ccf565b8281528092612cf4601f1991610cf0565b90600182811c92168015612f3f575b6020831014612f2b57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612f20565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff1615612f8157565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260405f206001600160a01b0333165f5260205260ff60405f20541615612fe25750565b63e2517d3f60e01b5f523360045260245260445ffd5b905f60a0525f906040830151928351906101408251015160a05260016001600160401b0360206040828180848801510151975101511698015101511601916001600160401b0383116125ff5761304f60a0516149a9565b81613653576101a08401805151909290156136295750505180519061307382613f49565b916001600160a01b0383169081156136175750823b90600182111561360557505f1981019081116125ff5760016130a982612ee9565b9360208501903c60206130bd83835161400e565b0160c05261ffff60c0515116916130e26008549361ffff8560e01c1684519114614be7565b60ff602083015116801515806135fa575b83516130fe91614be7565b6001600160401b0360095494600a5460e05260a01c169261311e82612ccc565b5f5f5b84811061347a57505061313f82516001600160401b03871115614be7565b81519061ffff60c0515116916040519061010082018281106001600160401b03821117610c4f576131c4946131b9945f93849360409b9a999b52855288602086015281604086015289606086015260808501528a60a085015260e05160c085015260e0518b1760e08501528160ff60208b01511694615444565b60a051818114614a3a565b5f925b8184106133265750505050604051916338f49afb60e01b835260a05160048401526020836024817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4928315611982575f936132f2575b5060095560e051600a557fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b806008549360a01b16169116176008558060075560a051600655956001600160401b036001946132786149ce565b6020604087015101525b16146132db57821580806132d0575b156132ab5750506040015160606020820151910152929190565b6132b7575b5050929190565b606060406132c9930151015190614a71565b5f806132b0565b5060a0518214613291565b50606060406132e86149ce565b9201510152929190565b9092506020813d60201161331e575b8161330e60209383610ccf565b8101031261024c5751915f613213565b3d9150613301565b9091929461ffff61333b876040850151614bef565b5116906001821b906001600160401b03613359896080870151614bef565b51169261336a896060870151614bef565b5190855161ffff60c05151168210156134655781602c810204602c14821517156125ff57602c82026030016030116125ff57816050602c82028b01015160e01c0361345357506034602c820201602c8202603001116125ff57602c81028881016054015160c01c90603c81016030909101116125ff57600195605c602c84028b0101519181145f1461342a5750841960e0511660e0525b8203613417575050901916955b019291906131c7565b5f52600b60205260405f2055179561340e565b825f52600c60205260405f20906001600160401b03198254161790558460e0511760e052613401565b634724a0fd60e01b5f5260045260245ffd5b6303e07d4560e61b5f5260045260245260445ffd5b63ffffffff61348d826040870151614bef565b5116916134a38361ffff60c0515116811061377a565b816135e1575b5081845190878a60c0515161ffff1660e051926134c595614156565b919097816080870151906134d891614bef565b516001600160401b031698826060880151906134f391614bef565b51610100526001600160401b03169289808514801595602094613529613566986135339661352e956135d3575b508c5190614be7565b612edc565b613f72565b604051639412e6b360e01b81526101005160048201526001600160401b03909a1660248b01529892839081906044820190565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af48015611982575f906135a1575b6001925061359a8286612cfe565b5201613121565b506020823d82116135cb575b816135ba60209383610ccf565b8101031261024c576001915161358c565b3d91506135ad565b90506101005114155f613520565b84516135f49163ffffffff168411614be7565b5f6134a9565b5060108111156130f3565b633610565160e21b5f5260045260245ffd5b633a517eed60e21b5f5260045260245ffd5b94965096905061364360206040850151015160a051614a71565b6001600160401b03600196613282565b9690936001600160401b03906136676149ce565b602060408701510152613282565b1561367e575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160401b03602060a08301510151165f52600560205260405f205480155f146136d85750505f90565b6040820190815160405161370d602082018093604080916001600160801b038151168452602081015160208501520151910152565b6060815261371c608082610ccf565b519020149182159261373a575b50501561373557600190565b600290565b6001600160801b0391925060208291015151169151511611155f80613729565b906001600160401b03809116911601906001600160401b0382116125ff57565b156137825750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156137bc5750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156137f65750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b908160011b91808304600214901517156125ff57565b15613847575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b909291926001600160a01b0361389783613f49565b16151580613e27575b613e2157604060209101510151519283519182156126265760b48311613e09576138cc83949294612ccc565b927317435cce3d1b4fa2e5f8a08ed921d57c6762a180925f5b87518110156139a357806020896001600160401b036040613917858561390e61394a9987612cfe565b51015194612cfe565b510151604051639412e6b360e01b81526004810193909352166001600160401b0316602482015292839081906044820190565b0381895af48015611982575f90613971575b6001925061396a8289612cfe565b52016138e5565b506020823d821161399b575b8161398a60209383610ccf565b8101031261024c576001915161395c565b3d915061397d565b5092935f96929591946139d4886139c16139bc89613828565b612ece565b93886139cc86612ccc565b9b8c926152a3565b50602c8602868104602c036125ff57603001806030116125ff578260051b91838304602014841517156125ff57613a17613a12602494602094613f72565b612ee9565b97828901956056875361564160218b01536256414c60228b01536356414c3460238b01538060381c858b01538060301c60258b01538060281c60268b015380841c60278b01538060181c60288b01538060101c60298b01538060081c602a8b0153602b8a01538060081c602c8a0153602d890153604051928380926338f49afb60e01b82528b60048301525af4908115611982575f91613dd7575b50602e8601528060081c604e860153604f8501536030955f935b8351851015613bc957613adf8585612cfe565b518887019060ff8760181c16602083015361ffff8760101c16602183015362ffffff8760081c16602283015363ffffffff8716602383015360048a018a116125ff57604081015166ffffffffffffff6001600160401b0382169160ff8160381c16602486015361ffff8160301c16602586015362ffffff8160281c16602686015363ffffffff8160201c16602786015364ffffffffff8160181c16602886015365ffffffffffff8160101c16602986015360081c16602a840153602b830153600c8a018a116125ff576020602c910151910152602c88018098116125ff5760018895019450613acc565b92509250935f955b8351871015613c0257613be48785612cfe565b5160208287010152602081018091116125ff57600190960195613bd1565b50939291509350613c13818461400e565b9160405191613c4460218460208101945f8652845180918484015e81015f838201520301601f198101855284610ccf565b616000835111613dab5750906001600160a01b0391613cfc602a7fffff000000000000000000000000000000000000000000000000000000000000845160f01b169360405193849160208301967f6100000000000000000000000000000000000000000000000000000000000000885260218401527f80600a3d393df30000000000000000000000000000000000000000000000000060238401525180918484015e81015f838201520301601f198101835282610ccf565b51905ff0168015613d83576008549260065560408201516007557fffff0000000000000000000000000000000000000000000000000000000000007dffff00000000000000000000000000000000000000000000000000000000602067ffffffffffffffff60a01b855160a01b1694015160e01b1693161717176008555f6009555f600a55565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b90506020813d602011613e01575b81613df260209383610ccf565b8101031261024c57515f613ab2565b3d9150613de5565b8263156f758160e31b5f5260045260b460245260445ffd5b50509050565b5061ffff60085460e01c1615156138a0565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f205416155f14613ec057805f525f60205260405f206001600160a01b0383165f5260205260405f20600160ff198254161790556001600160a01b03339216907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f2054165f14613ec057805f525f60205260405f206001600160a01b0383165f5260205260405f2060ff1981541690556001600160a01b03339216907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b60065403613f60576001600160a01b036008541690565b5f90565b90600182018092116125ff57565b919082018092116125ff57565b61ffff16602c810290808204602c14901517156125ff57603001806030116125ff5790565b9190823b6001811115613ff2575f1981019081116125ff578111613fd6576001613fcd82612ee9565b9360208501903c565b6001600160a01b0383633610565160e21b5f521660045260245ffd5b6001600160a01b0384633610565160e21b5f521660045260245ffd5b9190915f606060405161402081610c7e565b82815282602082015282604082015201526030835110613453576356414c34602084015160e01c0361345357602483015160c01c92602c810151938460f01c91602e81015190604e81015160f01c6040519361407b85610c7e565b845260208401928584526040850152606084019381855297851595861561414b575b508515614119575b505083156140b8575b5050506134535750565b519051915192509061ffff838116916140d19116613f7f565b908183149384156140ea575b50505050155f80806140ae565b621fffe0919293945060051b1690808204602014901517156125ff5761410f91613f72565b145f8080806140dd565b9091945060ef1c6201fffe61fffe8216911681036125ff575f190161ffff81116125ff5761ffff161415925f806140a5565b60b41095505f61409d565b939092959461ffff63ffffffff821693168310156142075761ffff1690602c830292808404602c14811517156125ff578360300194856030116125ff5784019581605088015160e01c0361345357506001901b9081166141ec576034830184116125ff57605485015160c01c5b96166141d95750603c01106125ff57605c015190565b925050505f52600b60205260405f205490565b815f52600c6020526001600160401b0360405f2054166141c3565b82856303e07d4560e61b5f5260045260245260445ffd5b999493979198909695995f96614233816112fa565b8061496857508915158061495c575b15614925576020019460208b6142af61426561425d8a614ee3565b923690610de9565b916122c5604051858101906142978287604080916001600160801b038151168452602081015160208501520151910152565b606081526142a6608082610ccf565b51902091614c00565b01518086036148f557505f905f5b8b8a818310614808575b50505050156147cb5750506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690813b1561024c579187926040519384927f87d3a9b100000000000000000000000000000000000000000000000000000000845260648401906004850152606060248501525260848201908960051b9860848a850101928b8a915f603e198d360301925b811061475d5750505050600319848403016044850152808352602083019260208260051b82010193835f92601e1982360301905b8585106145ee575050505050505091815f81819503925af18015611982576145d9575b50600185116143e4575b50505050506001600160801b036143de633b9aca00923690610de9565b51160490565b6143f5909691929496959395614ee3565b948335946001600160801b0386168096036145d55761441388610e74565b97614421604051998a610ccf565b885260208801918101903682116145d15780979597925b82841061453657505050506001600160401b03829316925b8651811015614519576144638188612cfe565b519085604051602081019087825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801918a5b8181106144e25750505050816144ca60019760206144d8940151605f19848303016080850152610c10565b03601f198101835282610ccf565b5190205d01614450565b91939496509194969760208061450460019360bf198b82030188528951610c10565b970194019101918c969493929897959861449f565b5094505050506001600160801b036143de633b9aca005f806143c1565b83989698356001600160401b0381116145cd5782016040813603126145cd576040519061456282610c34565b80356001600160401b0381116145c957810136601f820112156145c957614590903690602081359101614f19565b825260208101356001600160401b0381116145c957916145b7602094928594369101610d41565b83820152815201930192979597614438565b8880fd5b8680fd5b8480fd5b8380fd5b6145e69192505f90610ccf565b5f905f6143b7565b9193959750919395601f1982820301855287358381121561024c5784019061461a602082019280615034565b8091936020845252604082019060408160051b8401019380935f915b83831061465d57505050505050602080600192990195019501929091899796949592614394565b909192939495603f19838203018652614676878361507c565b8035600281101561024c5761468a816112fa565b82526146ad61469c6020830183615068565b606060208501526060840190615090565b906040810135609e198236030181121561024c57600193602093849361474f93019160408183039101526147416147236146f86146ea8580614fa4565b60a0865260a0860191614f84565b614703878601610dc8565b1515878501526147166040860186615068565b8482036040860152615090565b9261473060608201610dc8565b151560608401526080810190615068565b906080818403910152615090565b980196019493019190614636565b9193949650919460831989820301835285358481121561024c5760206147b96001938f839401906147ac6147a26147948480615034565b604085526040850191614fd5565b9285810190614fa4565b9185818503910152614f84565b9701930191019088969493918e614360565b6129b16040519283927ffef760c7000000000000000000000000000000000000000000000000000000008452602060048501526024840191614fd5565b6148266148208461483e9461482d9498969798614ef7565b80612e99565b3691614f19565b614838368789614f19565b906153c4565b156148eb575061075161485561485f928d8c614ef7565b6020810190612e67565b8051825180821491826148d5575b50501561488157505060015f808b8a6142c7565b906129b16148c3926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610c10565b83810360031901602485015290610c10565b9091506020830120906020840120145f8061486d565b91906001016142bd565b85907f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b897fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b5061ffff8a1115614242565b8760ff602492614977816112fa565b614980816112fa565b7f112d89cc00000000000000000000000000000000000000000000000000000000835216600452fd5b6149ba6001600160a01b0391613f49565b16156149c857600754600191565b5f905f90565b604051906149db82610c7e565b606082525f60208301526040516149f181610c7e565b606081525f60208201525f60408201525f606082015260408301525f60608301528160405190614a22602083610ccf565b5f825252565b52565b9081602091031261024c575190565b15614a43575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b8151518015614bd45760b48111614bbd5750604051917fc87f1f69000000000000000000000000000000000000000000000000000000008352602060048401528260a4810182519060806024840152815180915260c4830190602060c48260051b8601019301915f905b828210614b8e5750505050826001600160401b036060614b198594602080980151151560448701526040850151602319878303016064880152611304565b92015116608483015203817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4908115611982575f91614b58575b614b569250808214614a3a565b565b90506020823d602011614b86575b81614b7360209383610ccf565b8101031261024c57614b56915190614b49565b3d9150614b66565b91936001919395506020614bad819260c3198c82030186528851611304565b9601920192018794939192614adb565b63156f758160e31b5f5260045260b460245260445ffd5b50633a517eed60e21b5f5260045260245ffd5b156134535750565b906010811015612d125760051b0190565b6001600160401b03165f52600560205260405f20548015614c1e5790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040810151519060046020830151926001600160401b0384511690604063ffffffff602087015116950151906020825192015160208063ffffffff835116920151925101519260405194614c9986610c99565b85526020850197885260408501908152606085019182526080850192835260a0850193845261ffff60e0880151169760808801519760a08101519060c081015161012082015161014083015191610180610160850151940151946040519e8f9d8e7f4cc22bb7000000000000000000000000000000000000000000000000000000008152015260248d019d5f5b60088110614eaf57505060209d50614dcd8d63ffffffff97614dba829f9d9a9598614e089f9c98614da79060c09f9b6001600160401b039a614d80614d9492614d758e9c6101248c01906113d8565b6101648a01906113d8565b6102406101a4890152610244880190610ba1565b868103600319016101c488015290610bd4565b848103600319016101e486015290610b68565b91610204600319828503019101526113ff565b8c8103600319016102248e0152995116895251168b880152516040870152511660608501525160808401525160a08301829052910190610c10565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af1908115611982575f91614e75575b5015614e4d57565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90506020813d602011614ea7575b81614e9060209383610ccf565b8101031261024c57614ea1906110f2565b5f614e45565b3d9150614e83565b91939597999b9d5091939597999b9d602080600192855181520193019101908f9d9b99979593919e9c9a989694929e614d26565b356001600160401b038116810361024c5790565b9190811015612d125760051b81013590603e198136030182121561024c570190565b929190614f2581610e74565b93614f336040519586610ccf565b602085838152019160051b81019183831161024c5781905b838210614f59575050505050565b81356001600160401b03811161024c57602091614f798784938701610d41565b815201910190614f4b565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e198236030181121561024c5701602081359101916001600160401b03821161024c57813603831361024c57565b90602083828152019260208260051b82010193835f925b848410614ffc5750505050505090565b909192939495602080615024600193601f1986820301885261501e8b88614fa4565b90614f84565b9801940194019294939190614fec565b9035601e198236030181121561024c5701602081359101916001600160401b03821161024c578160051b3603831361024c57565b9035607e198236030181121561024c570190565b9035605e198236030181121561024c570190565b6150c96150ae6150a08380614fa4565b608086526080860191614f84565b6150bb6020840184614fa4565b908583036020870152614f84565b6150d66040830183615068565b9083810360408501528135600381101561024c576150f381610b24565b81526020820135600381101561024c5761510c81610b24565b6020820152604082013590600382101561024c576080615148615163948461513961515896999899610b24565b60408501526060810190614fa4565b9190928160608201520191614f84565b926060810190615034565b90916060818503910152808352602083019060208160051b85010193835f915b8383106151935750505050505090565b909192939495601f198282030186526151ac878461507c565b803591600383101561024c576152096020928392856151cc600197610b24565b81526151fb6151f06151e086850185614fa4565b6060888601526060850191614f84565b926040810190614fa4565b916040818503910152614f84565b980196019493019190615183565b90615254575080511561522c57602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b8151158061529a575b615265575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b1561525d565b9493919092946152b38483612edc565b600181146153a4576152c587916156a6565b92836152e66152d48289613f72565b846152de89613f64565b918a886152a3565b6152f7615319966153139299613f72565b9161530d6139bc6153078a613f64565b92613828565b90613f72565b936152a3565b60405163a5641f6f60e01b8152600481019390935260248301526020826044817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af491821561539f575f9261536a575b50614a28908294612cfe565b614a289192506153919060203d602011615398575b6153898183610ccf565b810190614a2b565b919061535e565b503d61537f565b611982565b506153b79150926153c092949593612cfe565b51928392612cfe565b5290565b908151815103613ec0575f5b825181101561542e576153e38184612cfe565b51516153ef8284612cfe565b515103615427576154008184612cfe565b51602081519101206154128284612cfe565b516020815191012003615427576001016153d0565b5050505f90565b505050600190565b5f1981146125ff5760010190565b939290949182841480809161568a575b6156625761ffff83116126265761546b8783612edc565b906001821461557e575061547e906156a6565b936154a561548c8689613f72565b9561530d6139bc61530761549f88613f64565b97613f64565b92815b85811087828a8361554d575b505050156154ca576154c590615436565b6154a8565b94959697918591886154dc948b615444565b946154e79596615444565b60405163a5641f6f60e01b81526004810192909252602482015260208180604481015b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af490811561539f575f91615534575090565b610d5c915060203d602011615398576153898183610ccf565b63ffffffff929350615574916040606061556a9301510151614bef565b5163ffffffff1690565b161087828a6154b4565b9250505094939294156155f75750506155c68360209261550a94955191848101516155ae604083015161ffff1690565b9063ffffffff60c060a0850151940151941694614156565b6040519384928392639412e6b360e01b8452600484019092916001600160401b036020916040840195845216910152565b615602829392613f64565b1490811591615637575b5061562357608061561f92930151612cfe565b5190565b8251634724a0fd60e01b5f5260045260245ffd5b905061565a61565161556a84604060608901510151614bef565b63ffffffff1690565b14155f61560c565b505092915050610d5c92508051906156846040602083015192015161ffff1690565b916156e3565b506156a161569d838960e08a01516156c2565b1590565b615454565b9060015b8060011b90838210156156bd57506156aa565b925050565b916156cf82600192612edc565b1b905f1982019182116125ff571b16151590565b9392909161ffff16602c810290808204602c14901517156125ff5760300190816030116125ff578060051b90808204602014901517156125ff5761572691613f72565b90602082018083116125ff578151106157425701602001519150565b83634724a0fd60e01b5f5260045260245ffdfea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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
