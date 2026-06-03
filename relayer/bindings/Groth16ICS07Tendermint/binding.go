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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membership_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviour_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"updateClient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_clientState\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCachedValidatorSet\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"updateClientMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CannotHandleMisbehavior\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientStateMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidMembershipProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"KeyValuePairNotInCache\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProofHeightMismatch\",\"inputs\":[{\"name\":\"expectedRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"expectedRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustThresholdMismatch\",\"inputs\":[{\"name\":\"expectedNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expectedDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnbondingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"UnknownZkAlgorithm\",\"inputs\":[{\"name\":\"algorithm\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"VerificationKeyMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x610140806040523461050a57616ba5803803809161001d828561062c565b8339810160e08282031261050a576100348261064f565b6100406020840161064f565b61004c6040850161064f565b916100596060860161064f565b60808601519094906001600160401b03811161050a5786019080601f8301121561050a57815161008b92602001610663565b9461009d60c060a0830151920161064f565b958051810190602082019060208184031261050a576020810151906001600160401b03821161050a570191829003601f198101906101201361050a576040519160e083016001600160401b038111848210176105fd5760405260208401516001600160401b03811161050a576020908501019080601f8301121561050a57815161012992602001610663565b82526040811261050a57604080519161014183610611565b61014c8286016106a8565b835261015a606086016106a8565b602084015260208401928352603f19011261050a5760405161017b81610611565b610187608085016106b6565b815261019560a085016106b6565b6020820152604083019081526101ad60c085016106ca565b90606084019182526101c160e086016106ca565b926080850193845261010086015195861515870361050a5760a08601968752610120015194600286101561050a5760c08101958652518051906001600160401b0382116105fd576102136001546106db565b601f81116105ad575b50602090601f83116001146105405763ffffffff95949392915f9183610535575b50508160011b915f199060031b1c1916176001555b5160ff81511661ff00602060025493015160081b169161ffff191617176002555160018060401b0381511660035491602068010000000000000000600160801b0391015160401b169160018060801b031916171760035551169267ffffffff00000000600454925160201b169051151560401b92519160028310156105215769ff00000000000000000068ff00000000000000009360481b169469ff000000000000000000199160018060481b031916171617911617176004556040516103238161031c81610713565b038261062c565b602081519101206101005260405163685e272760e11b8152602060048201526020818061035260248201610713565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4908115610516575f916104e0575b506101205260405161038e8161031c81610713565b6001600160401b03906020906103a3906107c4565b0151600354916001600160401b03831691168181036104cb575050604090811c6001600160401b03165f908152600560205220556001600160a01b0390811660805290811660a05290811660c0521660e05260045463ffffffff8082169061070882019081116104b75763ffffffff809360201c1692839116116104a257826001600160a01b0381166104895750610439610a9b565b505b604051615fc79081610b1e823960805181615213015260a05181614402015260c05181610525015260e0518181816125fb0152612cb601526101005181614d60015261012051816125c60152f35b8061049661049c926109a5565b50610a1b565b5061043b565b6333fae18560e21b5f5260045260245260445ffd5b634e487b7160e01b5f52601160045260245ffd5b6355bace6f60e01b5f5260045260245260445ffd5b90506020813d60201161050e575b816104fb6020938361062c565b8101031261050a57515f610379565b5f80fd5b3d91506104ee565b6040513d5f823e3d90fd5b634e487b7160e01b5f52602160045260245ffd5b015190505f8061023d565b90601f1983169160015f52815f20925f5b818110610595575091600193918563ffffffff99989796941061057d575b505050811b01600155610252565b01515f1960f88460031b161c191690555f808061056f565b92936020600181928786015181550195019301610551565b60015f525f516020616b655f395f51905f52601f840160051c810191602085106105f3575b601f0160051c01905b8181106105e8575061021c565b5f81556001016105db565b90915081906105d2565b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b038211176105fd57604052565b601f909101601f19168101906001600160401b038211908210176105fd57604052565b51906001600160a01b038216820361050a57565b9192916001600160401b0382116105fd576040519161068c601f8201601f19166020018461062c565b82948184528183011161050a578281602093845f96015e010152565b519060ff8216820361050a57565b51906001600160401b038216820361050a57565b519063ffffffff8216820361050a57565b90600182811c92168015610709575b60208310146106f557565b634e487b7160e01b5f52602260045260245ffd5b91607f16916106ea565b6001545f9291610722826106db565b8082529160018116908115610783575060011461073d575050565b60015f9081529293509091905f516020616b655f395f51905f525b838310610769575060209250010190565b600181602092949394548385870101520191019190610758565b9050602093945060ff929192191683830152151560051b010190565b9081518110156107b0570160200190565b634e487b7160e01b5f52603260045260245ffd5b6040516107d081610611565b606081525f60208201525080518015908115610999575b5061098a575f19908051805b610946575b505f1982146109385760018201908183116104b757600360fc1b6001600160f81b0319610825848461079f565b51161480610923575b610913575f5b81518310156108c357610847838361079f565b5160f81c6030811080156108b9575b6108a757600a82026001600160401b03908116602f1990920160ff169190910181169116811061088b57600190920191610834565b509150506040519061089c82610611565b81525f602082015290565b50509150506040519061089c82610611565b5060398111610856565b9150916001811190811591610907575b506108f857604051916108e583610611565b82526001600160401b0316602082015290565b63384c687760e21b5f5260045ffd5b602b915010155f6108d3565b9150506040519061089c82610611565b5080518381039081116104b75760021061082e565b604051915061089c82610611565b5f1981018181116104b757602d60f81b6001600160f81b0319610969838661079f565b511614610980575080156104b7575f1901806107f3565b92505f90506107f8565b6329120bff60e21b5f5260045ffd5b6040915010155f6107e7565b6001600160a01b0381165f9081525f516020616b855f395f51905f52602052604090205460ff16610a16576001600160a01b03165f8181525f516020616b855f395f51905f5260205260408120805460ff191660011790553391905f516020616ae55f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f516020616b055f395f51905f52602052604090205460ff16610a16576001600160a01b03165f8181525f516020616b055f395f51905f5260205260408120805460ff191660011790553391905f516020616b455f395f51905f52905f516020616ae55f395f51905f529080a4600190565b5f80525f516020616b055f395f51905f526020525f516020616b255f395f51905f525460ff16610b19575f8080525f516020616b055f395f51905f526020525f516020616b255f395f51905f52805460ff1916600117905533905f516020616b455f395f51905f525f516020616ae55f395f51905f528280a4600190565b5f9056fe610140806040526004361015610013575f80fd5b5f3560e01c90816301ffc9a714610da4575080630bece35614610d0a578063248a9ca314610ce05780632f2ff15d14610cb157806336568abe14610c62578063536c2ad314610c0a5780638a8e4c5d14610bd257806391d1485414610b96578063974a74c414610a4b578063a217fddf14610a31578063a6f031bb14610921578063ac9650d81461075c578063d547741f14610726578063ddba65371461023a5763ef913a4b146100c2575f80fd5b34610236575f3660031901126102365760405160208082015261012060408201525f6001546100f08161312d565b90816101608501526001811690815f1461021157506001146101b0575b6101ac83610198818560ff600254818116606085015260081c1660808301526001600160401b0360035481811660a085015260401c1660c083015260ff60045463ffffffff811660e085015263ffffffff8160201c16610100850152818160401c16151561012085015260481c1661018481611b3a565b61014083015203601f198101835282611073565b604051918291602083526020830190610f7d565b0390f35b91905060015f527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6915f905b8082106101f5575090915081016101800161019861010d565b91926001816020925461018085880101520191019092916101dc565b60ff19166101808086019190915291151560051b84019091019150610198905061010d565b5f80fd5b346102365761024836610e42565b6004549060ff8260401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff16156106f1575b820191602081840312610236578035906001600160401b03821161023657016101208184031261023657604051906102e082610fa1565b80356001600160401b03811161023657846102fc918301611179565b825260208101356001600160401b03811161023657810193606085820312610236576040519461032b86610fd0565b80356001600160401b038111610236578101604081840312610236576040519061035482610feb565b80356001600160401b03811161023657816103768660209361037e95016110e5565b845201611111565b6020820152865260208101356001600160401b03811161023657826103a4918301611456565b602087015260408101356001600160401b038111610236576103c891839101611456565b6040860152602083019485528060408301906103e39161125e565b60408401908152906103f89060a0840161125e565b9260608101928484526101000161040e9061124a565b9560808201948786528251915197845191604051998a947fa6fe8f56000000000000000000000000000000000000000000000000000000008652600486016101209052610124860161045f91611b44565b906003198683030160248701528051606083528051606084016040905260a0840161048991610f7d565b90602001516001600160401b0316608084015260208201519083810360208501526104b391611caf565b90604001519180820390604001526104ca91611caf565b83516001600160801b0316604486015260208401516064860152604090930151608485015280516001600160801b031660a4850152602081015160c48501526040015160e48401526001600160801b031661010483015203867f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031691815a93608094fa9586156106e6575f96610652575b50602080610640956105f66801000000000000000099966105a06105ee976001600160801b036001600160401b0398519151935195511690614d2a565b6040516105cd8582018093604080916001600160801b038151168452602081015160208501520151910152565b606081526105dc608082611073565b5190206105ee86858a51015116614ffe565b80821461388b565b6040516106238382018093604080916001600160801b038151168452602081015160208501520151910152565b60608152610632608082611073565b519020940151015116614ffe565b68ff0000000000000000191617600455005b9195509160803d6080116106df575b61066b8184611073565b820191608081840312610236576020610640956105f66001600160401b03946105a0680100000000000000009b6001600160801b0386976106c96105ee9b60408051936106b785610feb565b6106c183826118f0565b8552016118f0565b888201529d505097505096945050955050610563565b503d610661565b6040513d5f823e3d90fd5b6106f9613165565b6102a9565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102365761075a61073736610eaf565b90610755610750825f525f602052600160405f20015490565b6131d4565b613e0b565b005b34610236576020366003190112610236576004356001600160401b03811161023657366023820112156102365780600401356001600160401b0381116102365760248201913660248360051b83010111610236576020926040516107c08582611073565b5f815284810191601f1986013684376107d8856112e9565b936107e66040519586611073565b858552601f196107f5876112e9565b01875f5b828110610912575050505f5b868110156108b5576001906108915f808b8861085e61082c60248860051b8b01018b613083565b90938d6040519483869484860198893784019083820190898252519283915e010185815203601f198101835282611073565b5190305af43d156108ad573d9061087482611094565b916108826040519384611073565b82523d5f8d84013e5b3061585f565b61089b8289612def565b526108a68188612def565b5001610805565b60609061088b565b604080518981528751818b018190525f92600582901b83018101918a8d01918d9085015b8287106108e65785850386f35b909192938280610902600193603f198a82030186528851610f7d565b96019201960195929190926108d9565b606088820183015281016107f9565b34610236576020366003190112610236576004356001600160401b03811161023657806004019061014060031982360301126102365760ff60045460401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff1615610a24575b6109c360448201836130b5565b916109d160648201856130b5565b93909161010481013590600282101561023657602096610a1c966109f96101248401836130b5565b96909560405198610a0a8c8b611073565b5f8a52608460a4870196013594614344565b604051908152f35b610a2c613165565b6109b6565b34610236575f3660031901126102365760206040515f8152f35b34610236576020366003190112610236576004356001600160401b03811161023657806004019061016060031982360301126102365760ff60045460401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff1615610b89575b6101448101610aef8184613083565b905015610b6157610b0360448301846130b5565b9091610b1260648501866130b5565b95909461010481013591600283101561023657602097610a1c97610b4a96610b51610b416101248701866130b5565b99909886613083565b36916110af565b98608460a4870196013594614344565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b610b91613165565b610ae0565b3461023657610ba436610eaf565b905f525f6020526001600160a01b0360405f2091165f52602052602060ff60405f2054166040519015158152f35b3461023657610be036610e42565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b3461023657602036600319011261023657610c46610c546101ac610c2f600435612ed6565b919390604051958695606087526060870190610ed5565b908582036020870152610f0e565b908382036040850152610f41565b3461023657610c7036610eaf565b336001600160a01b03821603610c895761075a91613e0b565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102365761075a610cc236610eaf565b90610cdb610750825f525f602052600160405f20015490565b613d7e565b34610236576020366003190112610236576020610a1c6004355f525f602052600160405f20015490565b3461023657610d1836610e42565b9060ff60045460401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060209081527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47549092610d8692909160ff1615610d97576121d2565b60405190610d9381610e91565b8152f35b610d9f613165565b6121d2565b3461023657602036600319011261023657600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361023657817f7965db0b0000000000000000000000000000000000000000000000000000000060209314908115610e18575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483610e11565b906020600319830112610236576004356001600160401b0381116102365782602382011215610236578060040135926001600160401b0384116102365760248483010111610236576024019190565b60031115610e9b57565b634e487b7160e01b5f52602160045260245ffd5b604090600319011261023657600435906024356001600160a01b03811681036102365790565b90602080835192838152019201905f5b818110610ef25750505090565b825163ffffffff16845260209384019390920191600101610ee5565b90602080835192838152019201905f5b818110610f2b5750505090565b8251845260209384019390920191600101610f1e565b90602080835192838152019201905f5b818110610f5e5750505090565b82516001600160401b0316845260209384019390920191600101610f51565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b60a081019081106001600160401b03821117610fbc57604052565b634e487b7160e01b5f52604160045260245ffd5b606081019081106001600160401b03821117610fbc57604052565b604081019081106001600160401b03821117610fbc57604052565b60e081019081106001600160401b03821117610fbc57604052565b608081019081106001600160401b03821117610fbc57604052565b60c081019081106001600160401b03821117610fbc57604052565b61010081019081106001600160401b03821117610fbc57604052565b90601f801991011681019081106001600160401b03821117610fbc57604052565b6001600160401b038111610fbc57601f01601f191660200190565b9291926110bb82611094565b916110c96040519384611073565b829481845281830111610236578281602093845f960137010152565b9080601f8301121561023657816020611100933591016110af565b90565b359060ff8216820361023657565b35906001600160401b038216820361023657565b91908260409103126102365760405161113d81610feb565b602061115681839561114e81611111565b855201611111565b910152565b359063ffffffff8216820361023657565b3590811515820361023657565b80820392916101208412610236576040519161119483611006565b82948135906001600160401b038211610236576111b58460409385016110e5565b8552601f19011261023657611200610100926040516111d381610feb565b6111df60208501611103565b81526111ed60408501611103565b6020820152602086015260608301611125565b604084015261121160a0820161115b565b606084015261122260c0820161115b565b608084015261123360e0820161116c565b60a084015201359060028210156102365760c00152565b35906001600160801b038216820361023657565b91908260609103126102365760405161127681610fd0565b60408082946112848161124a565b8452602081013560208501520135910152565b809291039160608312610236576040516112b081610feb565b6040819483358352601f1901126102365760209060408051936112d285610feb565b6112dd84820161115b565b85520135828401520152565b6001600160401b038111610fbc5760051b60200190565b9190608083820312610236576040519061131982611021565b819380356001600160401b038111610236576060926113399183016110e5565b83526020810135602084015261135160408201611111565b60408401520135908160070b82036102365760600152565b919091608081840312610236576040519061138382611021565b819381356001600160401b03811161023657820181601f820112156102365780356113ad816112e9565b916113bb6040519384611073565b81835260208084019260051b820101918483116102365760208201905b838210611429575050505083526113f16020830161116c565b60208401526040820135906001600160401b038211610236578261141e6060949261115694869401611300565b604086015201611111565b81356001600160401b0381116102365760209161144b88848094880101611300565b8152019101906113d8565b919060a083820312610236576040519061146f82611021565b819380356001600160401b038111610236578101604081840312610236576040519061149a82610feb565b80356001600160401b0381116102365781016102c081860312610236576040519061026082018281106001600160401b03821117610fbc576040526114df8682611125565b825260408101356001600160401b03811161023657866115009183016110e5565b602083015261151160608201611111565b60408301526115226080820161124a565b606083015261153360a0820161116c565b60808301526115458660c08301611297565b60a0830152611557610120820161116c565b60c083015261014081013560e0830152611574610160820161116c565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526115c3610220820161116c565b6101c08301526102408101356101e08301526115e2610260820161116c565b6102008301526102808101356102208301526102a0810135906001600160401b03821161023657611615918791016110e5565b61024082015282526020810135906001600160401b038211610236570160c081850312610236576040519061164982611021565b61165281611111565b82526116606020820161115b565b60208301526116728560408301611297565b604083015260a0810135906001600160401b038211610236570184601f82011215610236578035906116a3826112e9565b916116b16040519384611073565b80835260208084019160051b830101918783116102365760208101915b83831061173c575050505060608201526020820152835260208101356001600160401b0381116102365782611704918301611369565b60208401526117168260408301611125565b60408401526080810135916001600160401b038311610236576060926111569201611369565b82356001600160401b038111610236578201906040828b03601f190112610236576040519161176a83610feb565b6020810135600481101561023657835260408101356001600160401b038111610236576020910101906080828c031261023657604051926117aa84611021565b82356001600160401b038111610236578c6117c69185016110e5565b84526117d46020840161124a565b60208501526117e56040840161116c565b60408501526060830135936001600160401b0385116102365761180d8d6020968796016110e5565b6060820152838201528152019201916116ce565b9080601f830112156102365760408051929061183d9084611073565b82906040810192831161023657905b8282106118595750505090565b813581526020918201910161184c565b9080601f83011215610236578135611880816112e9565b9261188e6040519485611073565b81845260208085019260051b82010192831161023657602001905b8282106118b65750505090565b602080916118c38461115b565b8152019101906118a9565b519060ff8216820361023657565b51906001600160401b038216820361023657565b91908260409103126102365760405161190881610feb565b6020611156818395611919816118dc565b8552016118dc565b519063ffffffff8216820361023657565b5190811515820361023657565b51906001600160801b038216820361023657565b91908260609103126102365760405161196b81610fd0565b60408082946119798161193f565b8452602081015160208501520151910152565b602081830312610236578051906001600160401b03821161023657016101808183031261023657604051916119c08361103c565b81516001600160401b0381116102365782019182820392610120841261023657604051936119ed85611006565b81516001600160401b0381116102365782019084601f83011215610236578151611a1681611094565b90611a246040519283611073565b8082528660208286010111610236576020815f9282604097018386015e830101528652601f1901126102365761010090604051611a6081610feb565b611a6c602083016118ce565b8152611a7a604083016118ce565b60208201526020860152611a9184606083016118f0565b6040860152611aa260a08201611921565b6060860152611ab360c08201611921565b6080860152611ac460e08201611932565b60a0860152015190600282101561023657836101409260c0611b329601528552611af18360208301611953565b6020860152611b038360808301611953565b6040860152611b1460e0820161193f565b6060860152611b278361010083016118f0565b6080860152016118f0565b60a082015290565b60021115610e9b57565b9061010060c0611b5f84516101208552610120850190610f7d565b9360ff60208083015182815116828801520151166040850152611b9f604082015160608601906001600160401b0360208092828151168552015116910152565b63ffffffff60608201511660a085015263ffffffff6080820151168285015260a0810151151560e0850152015191611bd683611b3a565b015290565b90606080611bf28451608085526080850190610f7d565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b828210611c8457505050506001600160401b036060611c7a819360208701511515602087015260408701518682036040880152611bdb565b9401511691015290565b90919293602080611ca1600193609f198a82030186528851611bdb565b960192019201909291611c42565b91909180519260a081526020611e0f8551604060a0850152611cea60e0850182516001600160401b0360208092828151168552015116910152565b610240611d08848301516102c06101208801526103a0870190610f7d565b60408301516001600160401b031661014087015260608301516001600160801b03166101608701526080830151151561018087015260a083015180516101a0880152602090810151805163ffffffff166101c089015201516101e08701529160c0810151151561020087015260e08101516102208701526101008101511515828701526101208101516102608701526101408101516102808701526101608101516102a08701526101808101516102c08701526101a08101516102e08701526101c081015115156103008701526101e08101516103208701526102008101511515610340870152610220810151610360870152015160df1985830301610380860152610f7d565b94015193609f198282030160c0830152606060c08201956001600160401b03815116835263ffffffff6020820151166020840152611e6f6040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519460c060a0830152855180915260e0820190602060e08260051b8501019701925f905b828210611ef45750505050506060611ebc611100949560208501518482036020860152611c1b565b92611ee4604082015160408501906001600160401b0360208092828151168552015116910152565b0151906080818403910152611c1b565b909192939760df1982820301855288519081516004811015610e9b57611f728260206001958195948295520151906040848201526060611f4083516080604085015260c0840190610f7d565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610f7d565b9a01950193920190611e94565b905f905b60028210611f9057505050565b6020806001928551815201930191019091611f83565b90602080835192838152019201905f5b818110611fc35750505090565b82511515845260209384019390920191600101611fb6565b90612035611ff483516109408452610940840190611b44565b61202360208501516020850190604080916001600160801b038151168452602081015160208501520151910152565b60408401518382036080850152611caf565b60608301516001600160801b031660a083015260808301515f60c084015b600882106121bc5750505061210a6120f66120e26120ce6120ba6101a09561208460a08a01516101c08a0190611f7f565b61209760c08a01516102008a0190611f7f565b61ffff60e08a0151166102408901526101008901518882036102608a0152610ed5565b610120880151878203610280890152610f0e565b6101408701518682036102a0880152610f41565b6101608601518582036102c0870152610ed5565b6101808501518482036102e0860152611fa6565b9201518051610300830152602081015160ff1661032083015260408101515f61034084015b601082106121a05750505060608101515f61054084015b6010821061218a5750505060800151905f90610740015b6010821061216b5750505090565b6020806001926001600160401b0386511681520193019101909161215d565b6020806001928551815201930191019091612146565b60208060019263ffffffff86511681520193019101909161212f565b6020806001928551815201930191019091612053565b90810190602081830312610236578035906001600160401b038211610236570190818103610940811261023657604051906101c082018281106001600160401b03821117610fbc5760405283356001600160401b0381116102365783612239918601611179565b8252612248836020860161125e565b602083015260808401356001600160401b038111610236578361226c918601611456565b926040830193845261228060a0860161124a565b60608401528060df860112156102365760405161229f61010082611073565b806101c08701918383116102365790839160c08901905b848210612d8757505060808601526122cd91611821565b60a08401526122e0816102008701611821565b60c08401526102408501359061ffff821682036102365760e084019182526102608601356001600160401b038111610236578161231e918801611869565b6101008501526102808601356001600160401b03811161023657860181601f8201121561023657803590612351826112e9565b9161235f6040519384611073565b80835260208084019160051b8301019184831161023657602001905b828210612d77575050506101208501526102a08601356001600160401b0381116102365786019581601f88011215610236578635966123b9886112e9565b976123c7604051998a611073565b8089526020808a019160051b8301019184831161023657602001905b828210612d5f5750505061014085019687526102c08101356001600160401b0381116102365782612415918301611869565b9361016086019485526102e08201356001600160401b03811161023657820183601f820112156102365780359061244b826112e9565b916124596040519384611073565b80835260208084019160051b8301019186831161023657602001905b828210612d47575050506101808701526106406102ff199091011261023657604051916124a183610fa1565b61030082013583526124b66103208301611103565b60208401528061035f8301121561023657610200916040516124d88482611073565b80610540830191848311610236576103408401905b838210612d2f57505060408601528261055f83011215610236576040516125148582611073565b8061074084019285841161023657905b838210612d1f57505060608601528261075f830112156102365761254b6040519485611073565b61094084920192831161023657905b828210612d075750505060808201526101a084015261257883613214565b979293919490845f14612c6c575f6040518080937f127dd052000000000000000000000000000000000000000000000000000000008252606060048301526125c3606483018c611fdb565b907f00000000000000000000000000000000000000000000000000000000000000006024840152604483015203816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156106e6575f91612c4a575b50965b61264788516001600160801b0360608b01511690614d2a565b6126a9602089015160405161267d602082018093604080916001600160801b038151168452602081015160208501520151910152565b6060815261268c608082611073565b5190206105ee6001600160401b03602060808d0151015116614ffe565b6126b2886138c2565b9415612a54575061ffff610100870151519251168092149283612a44575b83612a38575b5082612a2c575b5081612a1c575b50156129f457815f5260066020526001600160a01b0360405f20541680156129e1576127139094929194613e8e565b61271d818361408c565b80516001600160401b03166080525f958694859492919085805b6101008a01515189101561284357612754896101808c0151612def565b51156128335790889493929163ffffffff612774876101008e0151612def565b5f60a052511660a052612819575b50600196898660a0519c8960ff606089015116155f146127f2576127e2936127c26127d694600197946127dd9461ffff60208e0151169060a051936152ef565b919b9095905b61012060a051940151612def565b5114613a04565b613970565b985b019796999091929399612737565b6127e2936127d66128115f600198959c968c6127dd9660a05192614184565b9190956127c8565b60a05161282d9163ffffffff1681116139ca565b5f612782565b96976001909a949392919a6127e4565b5050975097945050939190506001600160401b0381166003810290808204600314901517156129cd57612885916080519161287f608051612e17565b10613a3e565b61288e81615044565b608051915b6129bc575b5050506128a482610e91565b81612973576001600160401b036020604060a0840193845183810151600354918683861c168783161161292a575b50505001516040516129048382018093604080916001600160801b038151168452602081015160208501520151910152565b60608152612913608082611073565b51902092510151165f52600560205260405f205590565b6fffffffffffffffff0000000000000000877fffffffffffffffffffffffffffffffff0000000000000000000000000000000092511692861b16921617176003555f80806128d2565b5061297d81610e91565b600181036129a5576801000000000000000068ff000000000000000019600454161760045590565b6129ae81610e91565b600281036111005750600290565b6129c592613a82565b5f8080612898565b634e487b7160e01b5f52601160045260245ffd5b82633a517eed60e21b5f5260045260245ffd5b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b905061018084015151145f6126e4565b5151811491505f6126dd565b5151821492505f6126d6565b61012087015151831493506126d0565b979493909698919260205f99510151519384519861ffff6101008a0151519251168092149283612c3a575b83612c2e575b5082612c22575b5081612c12575b50156129f4575f9893985b878110612bea57505f935f995f965f9b5b6101008a0151518d1015612b7c57612acc8d6101808c0151612def565b5115612b72578a8a888f8061010084015190612ae791612def565b5163ffffffff169c8d8096818097811090612b0191613990565b612b2b976127d69561012094602094612b2094612b5a575b5050612def565b510151930151612def565b6001612b5081986001600160401b036040612b468d8c612def565b5101511690613970565b9c5b019b96612aaf565b63ffffffff612b6b921681116139ca565b5f82612b19565b969b600190612b52565b509650969492995096509691506001600160401b0381166003810290808204600314901517156129cd576001600160401b038316916801fffffffffffffffe8460011b1692808404600214901517156129cd57612bdb92849211613a3e565b612be482615044565b91612893565b97612c086001916001600160401b036040612b468d999e9989612def565b9801989398612a9e565b905061018087015151145f612a93565b5151811491505f612a8c565b5151821492505f612a85565b6101208a01515183149350612a7f565b612c6691503d805f833e612c5e8183611073565b81019061198c565b5f61262b565b506040517fccd771d6000000000000000000000000000000000000000000000000000000008152602060048201525f8180612caa602482018b611fdb565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156106e6575f91612ced575b509661262e565b612d0191503d805f833e612c5e8183611073565b5f612ce6565b60208091612d1484611111565b81520191019061255a565b8135815260209182019101612524565b60208091612d3c8461115b565b8152019101906124ed565b60208091612d548461116c565b815201910190612475565b60208091612d6c84611111565b8152019101906123e3565b813581526020918201910161237b565b81358152859350602091820191016122b6565b90612da4826112e9565b612db16040519182611073565b8281528092612dc2601f19916112e9565b0190602036910137565b6040516112209190612dde8382611073565b6090815291601f1901366020840137565b8051821015612e035760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b908160011b91808304600214901517156129cd57565b90602c820291808304602c14901517156129cd57565b908160051b91808304602014901517156129cd57565b607e019081607e116129cd57565b90600482018092116129cd57565b90600c82018092116129cd57565b90602c82018092116129cd57565b90600282018092116129cd57565b90602282018092116129cd57565b90602082018092116129cd57565b90600182018092116129cd57565b919082018092116129cd57565b90815f5260066020526001600160a01b0360405f205416801561303957612efc90613e8e565b91612f0683613eee565b156130395791612f16818461408c565b926020840191612f2a61ffff845116612d9a565b94612f3961ffff855116612d9a565b92612f4861ffff865116612d9a565b9660608301905f5b61ffff88511661ffff82169081101561302d5780612f6e8185612def565b52878b87898860ff895116155f14612ff0575050505050602c8102818104602c14821517156129cd5780607e0180607e116129cd576082820181116129cd57608a828e938b0193612fc78660a287015160c01c92612def565b5201106129cd5761ffff928392612fe560aa6001940151918c612def565b525b01169050612f50565b8561ffff97948161301e6130165f8c9b60019b996001600160401b039961302499614184565b929097612def565b52612def565b91169052612fe7565b50509795505050505091565b50905060206040519161304c8284611073565b5f83525f368137604051926130618385611073565b5f84525f368137604051926130768185611073565b5f8452505f368137929190565b903590601e198136030182121561023657018035906001600160401b0382116102365760200191813603831361023657565b903590601e198136030182121561023657018035906001600160401b03821161023657602001918160051b3603831361023657565b5f198101919082116129cd57565b919082039182116129cd57565b9061310f82611094565b61311c6040519182611073565b8281528092612dc2601f1991611094565b90600182811c9216801561315b575b602083101461314757565b634e487b7160e01b5f52602260045260245ffd5b91607f169161313c565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff161561319d57565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260405f206001600160a01b0333165f5260205260ff60405f205416156131fe5750565b63e2517d3f60e01b5f523360045260245260445ffd5b5f6040820192835192835194610140865101519560016001600160401b0360206040828180848a01510151965101511699015101511601926001600160401b0384116129cd5761326388614ad3565b90918261386b576101a00180515190929015613846575050516101005261010051515f5260066020526001600160a01b0360405f205416801561382f576132a990613e8e565b6132b781610100515161408c565b6132d060ff608083015116600861010051519110614158565b610100516020015160c081905260ff1660e0819052151580613822575b61010051516132fb91614158565b61330660e051612d9a565b906001600160401b03815116915f935f5b60e051811061368c5750506133b3935061333f61010051516001600160401b03851115614158565b613347612dcc565b906133a85f613354612dcc565b92818061010051519261ffff60208a01511690604051946133748661103c565b855281602086015261010051604086015260608501528760808501528660a08501528160ff60206101005101511694615b7b565b95908d808214614b8e565b604051946338f49afb60e01b86528c60048701526020866024817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49586156106e6575f96613658575b5061ffff61340961340460e051612e2d565b612e59565b91169360228502858104602214861517156129cd57608061347d926134566001600160401b0361344561344060ff98968997612ec9565b613105565b9a61344f8c6158eb565b168a615933565b61346861ffff6020830151168a6159c1565b89602e8a01526001604e8a0153015116614172565b16604f850153610100515160508501525f607085015360e05160718501536134a583856159d1565b607e915f5b60e05181106135d357505f925b8484106135815750505050506134cc906153e0565b885f5260066020526001600160a01b0360405f20911673ffffffffffffffffffffffffffffffffffffffff19825416179055956001600160401b03600194613512614b22565b6020865101525b161461356c5782158080613563575b1561354157505051606060208201519101525b93929190565b61354d575b505061353b565b606061355c9251015190614bc5565b5f80613546565b50878214613528565b506060613577614b22565b9151015293929190565b6135cb6001916135918685612def565b51602161ffff82169160ff848c019160081c16602082015301536135b481612e91565b60206135c08888612def565b51918a010152612e9f565b9301926134b7565b926136516001916135fc63ffffffff6135f3886040610100510151614d19565b5116828a61590f565b61362a61360882612e67565b6001600160401b03613621896080610100510151614d19565b5116908a615979565b61363381612e75565b6020613646886060610100510151614d19565b51918a010152612e83565b93016134aa565b9095506020813d602011613684575b8161367460209383611073565b810103126102365751945f6133f2565b3d9150613667565b9384829394956101005160400151906136a491614d19565b5163ffffffff1686826101005160800151906136bf91614d19565b516001600160401b03168099846101005160600151906136de91614d19565b516101205283602084015161ffff168110906136f991613990565b8415156137859661374861372d5f6020986137529861374d976001600160401b0397613806575b50819d6101005151614184565b931692858414908115916137f8575b50610100515190614158565b6130f8565b612ec9565b604051639412e6b360e01b81526101205160048201526001600160401b0390991660248a01529792839081906044820190565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af480156106e6575f906137c6575b600192506137b98287612def565b5201949392919094613317565b506020823d82116137f0575b816137df60209383611073565b8101031261023657600191516137ab565b3d91506137d2565b90506101205114155f61373c565b610100515161381c9163ffffffff168411614158565b5f613720565b5060e051601010156132ed565b6101005151633a517eed60e21b5f5260045260245ffd5b94965096905061385b60208451015189614bc5565b6001600160401b03600196613519565b509690936001600160401b0390613880614b22565b602086510152613519565b15613894575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160401b03602060a08301510151165f52600560205260405f205480155f146138ee5750505f90565b60408201908151604051613923602082018093604080916001600160801b038151168452602081015160208501520151910152565b60608152613932608082611073565b5190201491821592613950575b50501561394b57600190565b600290565b6001600160801b0391925060208291015151169151511611155f8061393f565b906001600160401b03809116911601906001600160401b0382116129cd57565b156139985750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156139d25750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15613a0c5750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b15613a47575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b9190613a8d836153b4565b613d795760400151602001515180519384158015613d6f575b6129f457613ab5859495612d9a565b925f957317435cce3d1b4fa2e5f8a08ed921d57c6762a180965b8451811015613b8c5780602080613ae9613b339489612def565b5101516001600160401b036040613b00858b612def565b510151604051639412e6b360e01b81526004810193909352166001600160401b0316602482015292839081906044820190565b03818c5af480156106e6575f90613b5a575b60019250613b538289612def565b5201613acf565b506020823d8211613b84575b81613b7360209383611073565b810103126102365760019151613b45565b3d9150613b66565b5091949360205f97613bfe60249497613bc38b613bb0613bab84612e17565b6130ea565b9683613bbb89612d9a565b9e8f926159e1565b50613bf8613be8613440613bd961340485612e2d565b613be289612e43565b90612ec9565b99613bf28b6158eb565b8a615933565b886159c1565b604051938480926338f49afb60e01b82528760048301525af480156106e6575f90613d3b575b613c529250602e8601525f604e8601535f604f8601535f60508601525f60708601535f6071860153846159d1565b5f94607e5b8351871015613cc457613cbc600191613c708987612def565b51613c8263ffffffff8b16838a61590f565b613ca3613c8e83612e67565b6001600160401b03604084015116908a615979565b602080613caf84612e75565b9201519189010152612e83565b960195613c57565b95509150925f945b8251861015613cfa57613cf2600191613ce58886612def565b5160208288010152612ead565b950194613ccc565b50935050613d07906153e0565b905f5260066020526001600160a01b0360405f20911673ffffffffffffffffffffffffffffffffffffffff19825416179055565b506020823d602011613d67575b81613d5560209383611073565b8101031261023657613c529151613c24565b3d9150613d48565b5060b48511613aa6565b505050565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f205416155f14613e0557805f525f60205260405f206001600160a01b0383165f5260205260405f20600160ff198254161790556001600160a01b03339216907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f2054165f14613e0557805f525f60205260405f206001600160a01b0383165f5260205260405f2060ff1981541690556001600160a01b03339216907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b90813b6001811115613eb9575f1981019081116129cd576001613eb082613105565b9360208501903c565b6001600160a01b03837fd8415944000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b607e815110614087576356414c34602082015160e01c0361408757602c8101518060f01c91604e8101515f1a92607082015160f01c93609c830151938460f01c948215613fe8575050600114613f4657505050505f90565b8015908115613fdd575b508015613fd5575b8015613fcb575b8015613fc3575b8015613fb9575b613fb2575191602c810290808204602c14901517156129cd57607e019081607e116129cd576022810290808204602214901517156129cd57613fae91612ec9565b1490565b5050505f90565b5060908211613f6d565b508115613f66565b5060108311613f5f565b508215613f58565b60b49150115f613f50565b94939150945081158095811561407c575b811561404d575b50614044575193602c8202918204602c1417156129cd57607e019182607e116129cd57621fffe09060eb1c1690808204602014901517156129cd57613fae91612ec9565b50505050505f90565b905060ef1c6201fffe61fffe8216911681036129cd575f190161ffff81116129cd5761ffff168314155f614000565b60b484119150613ff9565b505f90565b9190915f60e060405161409e81611057565b8281528260208201528260408201528260608201528260808201528260a08201528260c08201520152607e835110614160576356414c34602084015160e01c0361416057614158602484015160c01c602c85015160f01c602e860151604e8701515f1a604f8801515f1a60508901519160708a015160f01c93609c8b015160f01c956040519761412d89611057565b8852602088015260408701526060860152608085015260a084015260c083015260e082015293613eee565b156141605750565b634724a0fd60e01b5f5260045260245ffd5b60ff60019116019060ff82116129cd57565b94919390936141a261419b602083015161ffff1690565b61ffff1690565b9563ffffffff8516968710801590614337575b61431b5760ff6141c9606084015160ff1690565b16156142e357505f95607e5b6141e761419b60c085015161ffff1690565b8810156142d3578681016020015160e01c82811461429a57821061421857614210600191612e83565b9701966141d5565b50509092945060a09193505b019261424961423c85515f52600660205260405f2090565b546001600160a01b031690565b6001600160a01b0381161561428657936142666142829495613e8e565b61427c61427482845161408c565b925194614172565b93614184565b9091565b8451633a517eed60e21b5f5260045260245ffd5b5096509294611100945092506142cb91506142c590506142b986612e67565b83016020015160c01c90565b94612e75565b016020015190565b50509092945060a0919350614224565b9294935050506142f561340485612e2d565b8281016020015190949060e01c036141605750611100906142cb6142c56142b986612e67565b6303e07d4560e61b5f5260045263ffffffff841660245260445ffd5b50600860ff8416116141b5565b999493979198909695995f9661435981611b3a565b80614a92575089151580614a86575b15614a4f576020019460208b6143d561438b6143838a61552b565b92369061125e565b916105ee604051858101906143bd8287604080916001600160801b038151168452602081015160208501520151910152565b606081526143cc608082611073565b51902091614ffe565b0151808603614a1f57505f905f5b8b8a818310614932575b50505050156148f15750506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690813b15610236579187926040519384927f87d3a9b100000000000000000000000000000000000000000000000000000000845260648401906004850152606060248501525260848201908960051b9860848a850101928b8a915f603e198d360301925b81106148835750505050600319848403016044850152808352602083019260208260051b82010193835f92601e1982360301905b858510614714575050505050505091815f81819503925af180156106e6576146ff575b506001851161450a575b50505050506001600160801b03614504633b9aca0092369061125e565b51160490565b61451b90969192949695939561552b565b948335946001600160801b0386168096036146fb57614539886112e9565b97614547604051998a611073565b885260208801918101903682116146f75780979597925b82841061465c57505050506001600160401b03829316925b865181101561463f576145898188612def565b519085604051602081019087825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801918a5b8181106146085750505050816145f060019760206145fe940151605f19848303016080850152610f7d565b03601f198101835282611073565b5190205d01614576565b91939496509194969760208061462a60019360bf198b82030188528951610f7d565b970194019101918c96949392989795986145c5565b5094505050506001600160801b03614504633b9aca005f806144e7565b83989698356001600160401b0381116146f35782016040813603126146f3576040519061468882610feb565b80356001600160401b0381116146ef57810136601f820112156146ef576146b6903690602081359101615561565b825260208101356001600160401b0381116146ef57916146dd6020949285943691016110e5565b8382015281520193019297959761455e565b8880fd5b8680fd5b8480fd5b8380fd5b61470c9192505f90611073565b5f905f6144dd565b9193959750919395601f198282030185528735838112156102365784019061474060208201928061567c565b8091936020845252604082019060408160051b8401019380935f915b838310614783575050505050506020806001929901950195019290918997969495926144ba565b909192939495603f1983820301865261479c87836156c4565b80356002811015610236576147b081611b3a565b82526147d36147c260208301836156b0565b6060602085015260608401906156d8565b906040810135609e1982360301811215610236576001936020938493614875930191604081830391015261486761484961481e61481085806155ec565b60a0865260a08601916155cc565b61482987860161116c565b15158785015261483c60408601866156b0565b84820360408601526156d8565b926148566060820161116c565b1515606084015260808101906156b0565b9060808184039101526156d8565b98019601949301919061475c565b919394965091946083198982030183528535848112156102365760206148df6001938f839401906148d26148c86148ba848061567c565b60408552604085019161561d565b92858101906155ec565b91858185039101526155cc565b9701930191019088969493918e614486565b61492e6040519283927ffef760c700000000000000000000000000000000000000000000000000000000845260206004850152602484019161561d565b0390fd5b61495061494a8461496894614957949896979861553f565b806130b5565b3691615561565b614962368789615561565b90615afc565b15614a155750610b4a61497f614989928d8c61553f565b6020810190613083565b8051825180821491826149ff575b5050156149ab57505060015f808b8a6143ed565b9061492e6149ed926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610f7d565b83810360031901602485015290610f7d565b9091506020830120906020840120145f80614997565b91906001016143e3565b85907f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b897fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b5061ffff8a1115614368565b8760ff602492614aa181611b3a565b614aaa81611b3a565b7f112d89cc00000000000000000000000000000000000000000000000000000000835216600452fd5b805f5260066020526001600160a01b0360405f2054168015614b1a57614af890613e8e565b90614b0282613eee565b15614b1a57604091614b139161408c565b0151600191565b50505f905f90565b60405190614b2f82611021565b606082525f6020830152604051614b4581611021565b606081525f60208201525f60408201525f606082015260408301525f60608301528160405190614b76602083611073565b5f825252565b52565b90816020910312610236575190565b15614b97575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b81515115614d0757604051917fc87f1f69000000000000000000000000000000000000000000000000000000008352602060048401528260a4810182519060806024840152815180915260c4830190602060c48260051b8601019301915f905b828210614cd85750505050826001600160401b036060614c638594602080980151151560448701526040850151602319878303016064880152611bdb565b92015116608483015203817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49081156106e6575f91614ca2575b614ca09250808214614b8e565b565b90506020823d602011614cd0575b81614cbd60209383611073565b8101031261023657614ca0915190614c93565b3d9150614cb0565b91936001919395506020614cf7819260c3198c82030186528851611bdb565b9601920192018794939192614c25565b633a517eed60e21b5f5260045260245ffd5b906010811015612e035760051b0190565b906001600160801b03633b9aca00911604428111614fcf57610708614d4f82426130f8565b11614fa057508051805160208201207f000000000000000000000000000000000000000000000000000000000000000003614ead5750602081015160ff815116906002549160ff8316928382149283614e94575b6020015160ff169215614e4c575050505063ffffffff606082015116906004549163ffffffff8316808203614e1e57505063ffffffff608081920151169160201c16808203614df0575050565b7f1eb5c1eb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f73b1bde3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6084945060ff90604051947fd382033a000000000000000000000000000000000000000000000000000000008652600486015260081c16602484015260448301526064820152fd5b6020810151600883901c60ff9081169116149350614da3565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f600154614ee98161312d565b908160448501526001811690815f14614f7c5750600114614f1c575b5061492e9192600319848303016024850152610f7d565b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b818310614f62575091929150810160640161492e614f05565b805460648488010152859450602090920191600101614f49565b60ff191660648086019190915291151560051b8401909101915061492e9050614f05565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b7f20daedb6000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b6001600160401b03165f52600560205260405f2054801561501c5790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040810151519060046020830151926001600160401b0384511690604063ffffffff602087015116950151906020825192015160208063ffffffff8351169201519251015192604051946150978661103c565b85526020850197885260408501908152606085019182526080850192835260a0850193845261ffff60e0880151169760808801519760a08101519060c081015161012082015161014083015191610180610160850151940151946040519e8f9d8e7f4cc22bb7000000000000000000000000000000000000000000000000000000008152015260248d019d5f5b600881106152ad57505060209d506151cb8d63ffffffff976151b8829f9d9a95986152069f9c986151a59060c09f9b6001600160401b039a61517e615192926151738e9c6101248c0190611f7f565b6101648a0190611f7f565b6102406101a4890152610244880190610f0e565b868103600319016101c488015290610f41565b848103600319016101e486015290610ed5565b9161020460031982850301910152611fa6565b8c8103600319016102248e0152995116895251168b880152516040870152511660608501525160808401525160a08301829052910190610f7d565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af19081156106e6575f91615273575b501561524b57565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90506020813d6020116152a5575b8161528e60209383611073565b810103126102365761529f90611932565b5f615243565b3d9150615281565b91939597999b9d5091939597999b9d602080600192855181520193019101908f9d9b99979593919e9c9a989694929e615124565b5f1981146129cd5760010190565b92949091945b61ffff8616811061531e575b63ffffffff85856303e07d4560e61b5f526004521660245260445ffd5b602c8102818104602c14821517156129cd5780607e019081607e116129cd5780850190609e82015160e01c9163ffffffff89169384841461537957505050116153745761536d61ffff916152e1565b90506152f5565b615301565b95975095509750505093506082850181116129cd57608a60a283015160c01c9501106129cd5760aa015191600181018091116129cd57929190565b5f5260066020526001600160a01b0360405f2054168015614087576153db61110091613e8e565b613eee565b6040519060208201905f8252615413602184835180602086018484015e81015f838201520301601f198101855284611073565b6160008351116154ff57506154c1602a7fffff000000000000000000000000000000000000000000000000000000000000845160f01b169360405193849160208301967f6100000000000000000000000000000000000000000000000000000000000000885260218401527f80600a3d393df30000000000000000000000000000000000000000000000000060238401525180918484015e81015f838201520301601f198101835282611073565b51905ff0906001600160a01b038216156154d757565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b356001600160401b03811681036102365790565b9190811015612e035760051b81013590603e1981360301821215610236570190565b92919061556d816112e9565b9361557b6040519586611073565b602085838152019160051b8101918383116102365781905b8382106155a1575050505050565b81356001600160401b038111610236576020916155c187849387016110e5565b815201910190615593565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156102365701602081359101916001600160401b03821161023657813603831361023657565b90602083828152019260208260051b82010193835f925b8484106156445750505050505090565b90919293949560208061566c600193601f198682030188526156668b886155ec565b906155cc565b9801940194019294939190615634565b9035601e19823603018112156102365701602081359101916001600160401b038211610236578160051b3603831361023657565b9035607e1982360301811215610236570190565b9035605e1982360301811215610236570190565b6157116156f66156e883806155ec565b6080865260808601916155cc565b61570360208401846155ec565b9085830360208701526155cc565b61571e60408301836156b0565b908381036040850152813560038110156102365761573b81610e91565b8152602082013560038110156102365761575481610e91565b602082015260408201359060038210156102365760806157906157ab94846157816157a096999899610e91565b604085015260608101906155ec565b91909281606082015201916155cc565b92606081019061567c565b90916060818503910152808352602083019060208160051b85010193835f915b8383106157db5750505050505090565b909192939495601f198282030186526157f487846156c4565b803591600383101561023657615851602092839285615814600197610e91565b8152615843615838615828868501856155ec565b60608886015260608501916155cc565b9260408101906155ec565b9160408185039101526155cc565b9801960194930191906157cb565b9061589c575080511561587457602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b815115806158e2575b6158ad575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b156158a5565b60236356414c34916056602082015361564160218201536256414c60228201530153565b90602391018260181c60208201538260101c60218201538260081c60228201530153565b602b908260381c60248201538260301c60258201538260281c60268201538260201c60278201538260181c60288201538260101c60298201538260081c602a8201530153565b90602791018260381c60208201538260301c60218201538260281c60228201538260201c60238201538260181c60248201538260101c60258201538260081c60268201530153565b602d908260081c602c8201530153565b609d908260081c609c8201530153565b9493919092946159f184836130f8565b60018114615adc57615a038791615ddb565b9283615a24615a128289612ec9565b84615a1c89612ebb565b918a886159e1565b615a35615a5196615a4b9299612ec9565b91613be2613bab615a458a612ebb565b92612e17565b936159e1565b60405163a5641f6f60e01b8152600481019390935260248301526020826044817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4918215615ad7575f92615aa2575b50614b7c908294612def565b614b7c919250615ac99060203d602011615ad0575b615ac18183611073565b810190614b7f565b9190615a96565b503d615ab7565b6106e6565b50615aef915092615af892949593612def565b51928392612def565b5290565b908151815103613e05575f5b8251811015615b5f57615b1b8184612def565b5151615b278284612def565b515103613fb257615b388184612def565b5160208151910120615b4a8284612def565b516020815191012003613fb257600101615b08565b505050600190565b61ffff60019116019061ffff82116129cd57565b969495919095939293808414615db35761ffff831660908110801590615da8575b6129f457615baa88846130f8565b9060018214615d175750615bbd90615ddb565b92615bc88489612ec9565b94615be5615bd588612ebb565b95613be2613bab615a458b612ebb565b94815b8b85821080615ceb575b15615c065750615c01906152e1565b615be8565b918882969798999a9b9c615c1c96959394615b7b565b90615c2b95919490968a615b7b565b60405163a5641f6f60e01b815260048101939093526024830191909152906020816044817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4908115615ad7575f91615ccc575b50809461ffff831660908110801590615cc1575b6129f45760a0615cbb92615cb461ffff6111009816615cab856080850151612def565b9061ffff169052565b0151612def565b52615b67565b5061ffff8511615c88565b615ce5915060203d602011615ad057615ac18183611073565b5f615c74565b508863ffffffff615d10615d06856040808701510151614d19565b5163ffffffff1690565b1610615bf2565b94809497989350615d289150612ebb565b1490811591615d7e575b50615d6a5761110093929160a087615cb4615d57615cbb95606061ffff9c0151612def565b51998a9616615cab856080850151612def565b8551634724a0fd60e01b5f5260045260245ffd5b9050615da0615d97615d06846040808c01510151614d19565b63ffffffff1690565b14155f615d32565b5061ffff8611615b9c565b505094615dd79394505f929150615dd16020825192015161ffff1690565b90615df7565b9190565b9060015b8060011b9083821015615df25750615ddf565b925050565b939192909261ffff81118015615fad575b615f9957615e2161423c865f52600660205260405f2090565b6001600160a01b03811615615f8557615e3990613e8e565b91615e44838761408c565b9560208701615e55815161ffff1690565b61ffff808916911603615f715760ff615e7260608a015160ff1690565b1615615f19575050615e99613404615e9461419b60c08a999a015161ffff1690565b612e2d565b955f9560e08101975b615eb161419b8a5161ffff1690565b881015615efa576020818701015160f01c8514615ee45761419b6001615ed9615eb193612e9f565b990198915050615ea2565b93505050506142cb919294506111009350612e91565b5060a001519496506111009550919250615f1390614172565b92615df7565b615f3d919750615f4394959650613be29250615e9461419b613404925161ffff1690565b91612e43565b90615f4d82612ead565b815110615f5d5701602001519150565b634724a0fd60e01b5f52600484905260245ffd5b634724a0fd60e01b5f52600482905260245ffd5b633a517eed60e21b5f52600486905260245ffd5b634724a0fd60e01b5f52600485905260245ffd5b50600860ff831611615e0856fea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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
