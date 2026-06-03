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
	Bin: "0x610140806040523461050a5761678b803803809161001d828561062c565b8339810160e08282031261050a576100348261064f565b6100406020840161064f565b61004c6040850161064f565b916100596060860161064f565b60808601519094906001600160401b03811161050a5786019080601f8301121561050a57815161008b92602001610663565b9461009d60c060a0830151920161064f565b958051810190602082019060208184031261050a576020810151906001600160401b03821161050a570191829003601f198101906101201361050a576040519160e083016001600160401b038111848210176105fd5760405260208401516001600160401b03811161050a576020908501019080601f8301121561050a57815161012992602001610663565b82526040811261050a57604080519161014183610611565b61014c8286016106a8565b835261015a606086016106a8565b602084015260208401928352603f19011261050a5760405161017b81610611565b610187608085016106b6565b815261019560a085016106b6565b6020820152604083019081526101ad60c085016106ca565b90606084019182526101c160e086016106ca565b926080850193845261010086015195861515870361050a5760a08601968752610120015194600286101561050a5760c08101958652518051906001600160401b0382116105fd576102136001546106db565b601f81116105ad575b50602090601f83116001146105405763ffffffff95949392915f9183610535575b50508160011b915f199060031b1c1916176001555b5160ff81511661ff00602060025493015160081b169161ffff191617176002555160018060401b0381511660035491602068010000000000000000600160801b0391015160401b169160018060801b031916171760035551169267ffffffff00000000600454925160201b169051151560401b92519160028310156105215769ff00000000000000000068ff00000000000000009360481b169469ff000000000000000000199160018060481b031916171617911617176004556040516103238161031c81610713565b038261062c565b602081519101206101005260405163685e272760e11b8152602060048201526020818061035260248201610713565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4908115610516575f916104e0575b506101205260405161038e8161031c81610713565b6001600160401b03906020906103a3906107c4565b0151600354916001600160401b03831691168181036104cb575050604090811c6001600160401b03165f908152600560205220556001600160a01b0390811660805290811660a05290811660c0521660e05260045463ffffffff8082169061070882019081116104b75763ffffffff809360201c1692839116116104a257826001600160a01b0381166104895750610439610a9b565b505b604051615bad9081610b1e8239608051816151fa015260a0518161441d015260c05181610525015260e0518181816125e80152612c9c01526101005181614d47015261012051816125b30152f35b8061049661049c926109a5565b50610a1b565b5061043b565b6333fae18560e21b5f5260045260245260445ffd5b634e487b7160e01b5f52601160045260245ffd5b6355bace6f60e01b5f5260045260245260445ffd5b90506020813d60201161050e575b816104fb6020938361062c565b8101031261050a57515f610379565b5f80fd5b3d91506104ee565b6040513d5f823e3d90fd5b634e487b7160e01b5f52602160045260245ffd5b015190505f8061023d565b90601f1983169160015f52815f20925f5b818110610595575091600193918563ffffffff99989796941061057d575b505050811b01600155610252565b01515f1960f88460031b161c191690555f808061056f565b92936020600181928786015181550195019301610551565b60015f525f51602061674b5f395f51905f52601f840160051c810191602085106105f3575b601f0160051c01905b8181106105e8575061021c565b5f81556001016105db565b90915081906105d2565b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b038211176105fd57604052565b601f909101601f19168101906001600160401b038211908210176105fd57604052565b51906001600160a01b038216820361050a57565b9192916001600160401b0382116105fd576040519161068c601f8201601f19166020018461062c565b82948184528183011161050a578281602093845f96015e010152565b519060ff8216820361050a57565b51906001600160401b038216820361050a57565b519063ffffffff8216820361050a57565b90600182811c92168015610709575b60208310146106f557565b634e487b7160e01b5f52602260045260245ffd5b91607f16916106ea565b6001545f9291610722826106db565b8082529160018116908115610783575060011461073d575050565b60015f9081529293509091905f51602061674b5f395f51905f525b838310610769575060209250010190565b600181602092949394548385870101520191019190610758565b9050602093945060ff929192191683830152151560051b010190565b9081518110156107b0570160200190565b634e487b7160e01b5f52603260045260245ffd5b6040516107d081610611565b606081525f60208201525080518015908115610999575b5061098a575f19908051805b610946575b505f1982146109385760018201908183116104b757600360fc1b6001600160f81b0319610825848461079f565b51161480610923575b610913575f5b81518310156108c357610847838361079f565b5160f81c6030811080156108b9575b6108a757600a82026001600160401b03908116602f1990920160ff169190910181169116811061088b57600190920191610834565b509150506040519061089c82610611565b81525f602082015290565b50509150506040519061089c82610611565b5060398111610856565b9150916001811190811591610907575b506108f857604051916108e583610611565b82526001600160401b0316602082015290565b63384c687760e21b5f5260045ffd5b602b915010155f6108d3565b9150506040519061089c82610611565b5080518381039081116104b75760021061082e565b604051915061089c82610611565b5f1981018181116104b757602d60f81b6001600160f81b0319610969838661079f565b511614610980575080156104b7575f1901806107f3565b92505f90506107f8565b6329120bff60e21b5f5260045ffd5b6040915010155f6107e7565b6001600160a01b0381165f9081525f51602061676b5f395f51905f52602052604090205460ff16610a16576001600160a01b03165f8181525f51602061676b5f395f51905f5260205260408120805460ff191660011790553391905f5160206166cb5f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f5160206166eb5f395f51905f52602052604090205460ff16610a16576001600160a01b03165f8181525f5160206166eb5f395f51905f5260205260408120805460ff191660011790553391905f51602061672b5f395f51905f52905f5160206166cb5f395f51905f529080a4600190565b5f80525f5160206166eb5f395f51905f526020525f51602061670b5f395f51905f525460ff16610b19575f8080525f5160206166eb5f395f51905f526020525f51602061670b5f395f51905f52805460ff1916600117905533905f51602061672b5f395f51905f525f5160206166cb5f395f51905f528280a4600190565b5f9056fe6101a0806040526004361015610013575f80fd5b5f3560e01c90816301ffc9a714610da4575080630bece35614610d0a578063248a9ca314610ce05780632f2ff15d14610cb157806336568abe14610c62578063536c2ad314610c0a5780638a8e4c5d14610bd257806391d1485414610b96578063974a74c414610a4b578063a217fddf14610a31578063a6f031bb14610921578063ac9650d81461075c578063d547741f14610726578063ddba65371461023a5763ef913a4b146100c2575f80fd5b34610236575f3660031901126102365760405160208082015261012060408201525f6001546100f081612fdb565b90816101608501526001811690815f1461021157506001146101b0575b6101ac83610198818560ff600254818116606085015260081c1660808301526001600160401b0360035481811660a085015260401c1660c083015260ff60045463ffffffff811660e085015263ffffffff8160201c16610100850152818160401c16151561012085015260481c1661018481611b1e565b61014083015203601f198101835282611057565b604051918291602083526020830190610f7d565b0390f35b91905060015f527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6915f905b8082106101f5575090915081016101800161019861010d565b91926001816020925461018085880101520191019092916101dc565b60ff19166101808086019190915291151560051b84019091019150610198905061010d565b5f80fd5b346102365761024836610e42565b6004549060ff8260401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff16156106f1575b820191602081840312610236578035906001600160401b03821161023657016101208184031261023657604051906102e082610fa1565b80356001600160401b03811161023657846102fc91830161115d565b825260208101356001600160401b03811161023657810193606085820312610236576040519461032b86610fd0565b80356001600160401b038111610236578101604081840312610236576040519061035482610feb565b80356001600160401b03811161023657816103768660209361037e95016110c9565b8452016110f5565b6020820152865260208101356001600160401b03811161023657826103a491830161143a565b602087015260408101356001600160401b038111610236576103c89183910161143a565b6040860152602083019485528060408301906103e391611242565b60408401908152906103f89060a08401611242565b9260608101928484526101000161040e9061122e565b9560808201948786528251915197845191604051998a947fa6fe8f56000000000000000000000000000000000000000000000000000000008652600486016101209052610124860161045f91611b28565b906003198683030160248701528051606083528051606084016040905260a0840161048991610f7d565b90602001516001600160401b0316608084015260208201519083810360208501526104b391611c93565b90604001519180820390604001526104ca91611c93565b83516001600160801b0316604486015260208401516064860152604090930151608485015280516001600160801b031660a4850152602081015160c48501526040015160e48401526001600160801b031661010483015203867f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031691815a93608094fa9586156106e6575f96610652575b50602080610640956105f66801000000000000000099966105a06105ee976001600160801b036001600160401b0398519151935195511690614d11565b6040516105cd8582018093604080916001600160801b038151168452602081015160208501520151910152565b606081526105dc608082611057565b5190206105ee86858a51015116614fe5565b8082146137f1565b6040516106238382018093604080916001600160801b038151168452602081015160208501520151910152565b60608152610632608082611057565b519020940151015116614fe5565b68ff0000000000000000191617600455005b9195509160803d6080116106df575b61066b8184611057565b820191608081840312610236576020610640956105f66001600160401b03946105a0680100000000000000009b6001600160801b0386976106c96105ee9b60408051936106b785610feb565b6106c183826118d4565b8552016118d4565b888201529d505097505096945050955050610563565b503d610661565b6040513d5f823e3d90fd5b6106f9613013565b6102a9565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102365761075a61073736610eaf565b90610755610750825f525f602052600160405f20015490565b613082565b61403d565b005b34610236576020366003190112610236576004356001600160401b03811161023657366023820112156102365780600401356001600160401b0381116102365760248201913660248360051b83010111610236576020926040516107c08582611057565b5f815284810191601f1986013684376107d8856112cd565b936107e66040519586611057565b858552601f196107f5876112cd565b01875f5b828110610912575050505f5b868110156108b5576001906108915f808b8861085e61082c60248860051b8b01018b612f31565b90938d6040519483869484860198893784019083820190898252519283915e010185815203601f198101835282611057565b5190305af43d156108ad573d9061087482611078565b916108826040519384611057565b82523d5f8d84013e5b306155fc565b61089b8289612dd5565b526108a68188612dd5565b5001610805565b60609061088b565b604080518981528751818b018190525f92600582901b83018101918a8d01918d9085015b8287106108e65785850386f35b909192938280610902600193603f198a82030186528851610f7d565b96019201960195929190926108d9565b606088820183015281016107f9565b34610236576020366003190112610236576004356001600160401b03811161023657806004019061014060031982360301126102365760ff60045460401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff1615610a24575b6109c36044820183612f63565b916109d16064820185612f63565b93909161010481013590600282101561023657602096610a1c966109f9610124840183612f63565b96909560405198610a0a8c8b611057565b5f8a52608460a487019601359461435f565b604051908152f35b610a2c613013565b6109b6565b34610236575f3660031901126102365760206040515f8152f35b34610236576020366003190112610236576004356001600160401b03811161023657806004019061016060031982360301126102365760ff60045460401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff1615610b89575b6101448101610aef8184612f31565b905015610b6157610b036044830184612f63565b9091610b126064850186612f63565b95909461010481013591600283101561023657602097610a1c97610b4a96610b51610b41610124870186612f63565b99909886612f31565b3691611093565b98608460a487019601359461435f565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b610b91613013565b610ae0565b3461023657610ba436610eaf565b905f525f6020526001600160a01b0360405f2091165f52602052602060ff60405f2054166040519015158152f35b3461023657610be036610e42565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b3461023657602036600319011261023657610c46610c546101ac610c2f600435612dfd565b919390604051958695606087526060870190610ed5565b908582036020870152610f0e565b908382036040850152610f41565b3461023657610c7036610eaf565b336001600160a01b03821603610c895761075a9161403d565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102365761075a610cc236610eaf565b90610cdb610750825f525f602052600160405f20015490565b613fb0565b34610236576020366003190112610236576020610a1c6004355f525f602052600160405f20015490565b3461023657610d1836610e42565b9060ff60045460401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060209081527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47549092610d8692909160ff1615610d97576121b6565b60405190610d9381610e91565b8152f35b610d9f613013565b6121b6565b3461023657602036600319011261023657600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361023657817f7965db0b0000000000000000000000000000000000000000000000000000000060209314908115610e18575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483610e11565b906020600319830112610236576004356001600160401b0381116102365782602382011215610236578060040135926001600160401b0384116102365760248483010111610236576024019190565b60031115610e9b57565b634e487b7160e01b5f52602160045260245ffd5b604090600319011261023657600435906024356001600160a01b03811681036102365790565b90602080835192838152019201905f5b818110610ef25750505090565b825163ffffffff16845260209384019390920191600101610ee5565b90602080835192838152019201905f5b818110610f2b5750505090565b8251845260209384019390920191600101610f1e565b90602080835192838152019201905f5b818110610f5e5750505090565b82516001600160401b0316845260209384019390920191600101610f51565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b60a081019081106001600160401b03821117610fbc57604052565b634e487b7160e01b5f52604160045260245ffd5b606081019081106001600160401b03821117610fbc57604052565b604081019081106001600160401b03821117610fbc57604052565b60e081019081106001600160401b03821117610fbc57604052565b608081019081106001600160401b03821117610fbc57604052565b60c081019081106001600160401b03821117610fbc57604052565b90601f801991011681019081106001600160401b03821117610fbc57604052565b6001600160401b038111610fbc57601f01601f191660200190565b92919261109f82611078565b916110ad6040519384611057565b829481845281830111610236578281602093845f960137010152565b9080601f83011215610236578160206110e493359101611093565b90565b359060ff8216820361023657565b35906001600160401b038216820361023657565b91908260409103126102365760405161112181610feb565b602061113a818395611132816110f5565b8552016110f5565b910152565b359063ffffffff8216820361023657565b3590811515820361023657565b80820392916101208412610236576040519161117883611006565b82948135906001600160401b038211610236576111998460409385016110c9565b8552601f190112610236576111e4610100926040516111b781610feb565b6111c3602085016110e7565b81526111d1604085016110e7565b6020820152602086015260608301611109565b60408401526111f560a0820161113f565b606084015261120660c0820161113f565b608084015261121760e08201611150565b60a084015201359060028210156102365760c00152565b35906001600160801b038216820361023657565b91908260609103126102365760405161125a81610fd0565b60408082946112688161122e565b8452602081013560208501520135910152565b8092910391606083126102365760405161129481610feb565b6040819483358352601f1901126102365760209060408051936112b685610feb565b6112c184820161113f565b85520135828401520152565b6001600160401b038111610fbc5760051b60200190565b919060808382031261023657604051906112fd82611021565b819380356001600160401b0381116102365760609261131d9183016110c9565b835260208101356020840152611335604082016110f5565b60408401520135908160070b82036102365760600152565b919091608081840312610236576040519061136782611021565b819381356001600160401b03811161023657820181601f82011215610236578035611391816112cd565b9161139f6040519384611057565b81835260208084019260051b820101918483116102365760208201905b83821061140d575050505083526113d560208301611150565b60208401526040820135906001600160401b03821161023657826114026060949261113a948694016112e4565b6040860152016110f5565b81356001600160401b0381116102365760209161142f888480948801016112e4565b8152019101906113bc565b919060a083820312610236576040519061145382611021565b819380356001600160401b038111610236578101604081840312610236576040519061147e82610feb565b80356001600160401b0381116102365781016102c081860312610236576040519061026082018281106001600160401b03821117610fbc576040526114c38682611109565b825260408101356001600160401b03811161023657866114e49183016110c9565b60208301526114f5606082016110f5565b60408301526115066080820161122e565b606083015261151760a08201611150565b60808301526115298660c0830161127b565b60a083015261153b6101208201611150565b60c083015261014081013560e08301526115586101608201611150565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526115a76102208201611150565b6101c08301526102408101356101e08301526115c66102608201611150565b6102008301526102808101356102208301526102a0810135906001600160401b038211610236576115f9918791016110c9565b61024082015282526020810135906001600160401b038211610236570160c081850312610236576040519061162d82611021565b611636816110f5565b82526116446020820161113f565b6020830152611656856040830161127b565b604083015260a0810135906001600160401b038211610236570184601f8201121561023657803590611687826112cd565b916116956040519384611057565b80835260208084019160051b830101918783116102365760208101915b838310611720575050505060608201526020820152835260208101356001600160401b03811161023657826116e891830161134d565b60208401526116fa8260408301611109565b60408401526080810135916001600160401b0383116102365760609261113a920161134d565b82356001600160401b038111610236578201906040828b03601f190112610236576040519161174e83610feb565b6020810135600481101561023657835260408101356001600160401b038111610236576020910101906080828c0312610236576040519261178e84611021565b82356001600160401b038111610236578c6117aa9185016110c9565b84526117b86020840161122e565b60208501526117c960408401611150565b60408501526060830135936001600160401b038511610236576117f18d6020968796016110c9565b6060820152838201528152019201916116b2565b9080601f83011215610236576040805192906118219084611057565b82906040810192831161023657905b82821061183d5750505090565b8135815260209182019101611830565b9080601f83011215610236578135611864816112cd565b926118726040519485611057565b81845260208085019260051b82010192831161023657602001905b82821061189a5750505090565b602080916118a78461113f565b81520191019061188d565b519060ff8216820361023657565b51906001600160401b038216820361023657565b9190826040910312610236576040516118ec81610feb565b602061113a8183956118fd816118c0565b8552016118c0565b519063ffffffff8216820361023657565b5190811515820361023657565b51906001600160801b038216820361023657565b91908260609103126102365760405161194f81610fd0565b604080829461195d81611923565b8452602081015160208501520151910152565b602081830312610236578051906001600160401b03821161023657016101808183031261023657604051916119a48361103c565b81516001600160401b0381116102365782019182820392610120841261023657604051936119d185611006565b81516001600160401b0381116102365782019084601f830112156102365781516119fa81611078565b90611a086040519283611057565b8082528660208286010111610236576020815f9282604097018386015e830101528652601f1901126102365761010090604051611a4481610feb565b611a50602083016118b2565b8152611a5e604083016118b2565b60208201526020860152611a7584606083016118d4565b6040860152611a8660a08201611905565b6060860152611a9760c08201611905565b6080860152611aa860e08201611916565b60a0860152015190600282101561023657836101409260c0611b169601528552611ad58360208301611937565b6020860152611ae78360808301611937565b6040860152611af860e08201611923565b6060860152611b0b8361010083016118d4565b6080860152016118d4565b60a082015290565b60021115610e9b57565b9061010060c0611b4384516101208552610120850190610f7d565b9360ff60208083015182815116828801520151166040850152611b83604082015160608601906001600160401b0360208092828151168552015116910152565b63ffffffff60608201511660a085015263ffffffff6080820151168285015260a0810151151560e0850152015191611bba83611b1e565b015290565b90606080611bd68451608085526080850190610f7d565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b828210611c6857505050506001600160401b036060611c5e819360208701511515602087015260408701518682036040880152611bbf565b9401511691015290565b90919293602080611c85600193609f198a82030186528851611bbf565b960192019201909291611c26565b91909180519260a081526020611df38551604060a0850152611cce60e0850182516001600160401b0360208092828151168552015116910152565b610240611cec848301516102c06101208801526103a0870190610f7d565b60408301516001600160401b031661014087015260608301516001600160801b03166101608701526080830151151561018087015260a083015180516101a0880152602090810151805163ffffffff166101c089015201516101e08701529160c0810151151561020087015260e08101516102208701526101008101511515828701526101208101516102608701526101408101516102808701526101608101516102a08701526101808101516102c08701526101a08101516102e08701526101c081015115156103008701526101e08101516103208701526102008101511515610340870152610220810151610360870152015160df1985830301610380860152610f7d565b94015193609f198282030160c0830152606060c08201956001600160401b03815116835263ffffffff6020820151166020840152611e536040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519460c060a0830152855180915260e0820190602060e08260051b8501019701925f905b828210611ed85750505050506060611ea06110e4949560208501518482036020860152611bff565b92611ec8604082015160408501906001600160401b0360208092828151168552015116910152565b0151906080818403910152611bff565b909192939760df1982820301855288519081516004811015610e9b57611f568260206001958195948295520151906040848201526060611f2483516080604085015260c0840190610f7d565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610f7d565b9a01950193920190611e78565b905f905b60028210611f7457505050565b6020806001928551815201930191019091611f67565b90602080835192838152019201905f5b818110611fa75750505090565b82511515845260209384019390920191600101611f9a565b90612019611fd883516109408452610940840190611b28565b61200760208501516020850190604080916001600160801b038151168452602081015160208501520151910152565b60408401518382036080850152611c93565b60608301516001600160801b031660a083015260808301515f60c084015b600882106121a0575050506120ee6120da6120c66120b261209e6101a09561206860a08a01516101c08a0190611f63565b61207b60c08a01516102008a0190611f63565b61ffff60e08a0151166102408901526101008901518882036102608a0152610ed5565b610120880151878203610280890152610f0e565b6101408701518682036102a0880152610f41565b6101608601518582036102c0870152610ed5565b6101808501518482036102e0860152611f8a565b9201518051610300830152602081015160ff1661032083015260408101515f61034084015b601082106121845750505060608101515f61054084015b6010821061216e5750505060800151905f90610740015b6010821061214f5750505090565b6020806001926001600160401b03865116815201930191019091612141565b602080600192855181520193019101909161212a565b60208060019263ffffffff865116815201930191019091612113565b6020806001928551815201930191019091612037565b9190820191602081840312610236578035906001600160401b03821161023657019081830392610940841261023657604051906101c082018281106001600160401b03821117610fbc5760405283356001600160401b038111610236578161221f91860161115d565b825261222e8160208601611242565b602083015260808401356001600160401b038111610236578161225291860161143a565b936040830194855261226660a0820161122e565b60608401528160df820112156102365760405161228561010082611057565b806101c08301918483116102365790849160c08501905b848210612d6d57505060808601526122b391611805565b60a08401526122c6826102008301611805565b60c08401526102408101359161ffff831683036102365760e084019283526102608201356001600160401b038111610236578161230491840161184d565b6101008501526102808201356001600160401b03811161023657820181601f8201121561023657803590612337826112cd565b916123456040519384611057565b80835260208084019160051b8301019184831161023657602001905b828210612d5d575050506101208501526102a08201356001600160401b0381116102365782019181601f840112156102365782359261239f846112cd565b936123ad6040519586611057565b80855260208086019160051b8301019184831161023657602001905b828210612d455750505061014085019283526102c08101356001600160401b03811161023657826123fb91830161184d565b9761016086019889526102e08201356001600160401b03811161023657820183601f8201121561023657803590612431826112cd565b9161243f6040519384611057565b80835260208084019160051b8301019186831161023657602001905b828210612d2d575050506101808701526106406102ff1990910112610236576040519161248783610fa1565b610300820135835261249c61032083016110e7565b60208401528061035f8301121561023657610200916040516124be8482611057565b80610540830191848311610236576103408401905b838210612d1557505060408601528261055f83011215610236576040516124fa8582611057565b8061074084019285841161023657905b838210612d0557505060608601528261075f83011215610236576125316040519485611057565b61094084920192831161023657905b828210612ced5750505060808201526101a084015261255e836130c2565b91909260a051929760e0519498895f14612c52575f6040518080937f127dd052000000000000000000000000000000000000000000000000000000008252606060048301526125b0606483018c611fbf565b907f00000000000000000000000000000000000000000000000000000000000000006024840152604483015203816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156106e6575f91612c30575b50975b61263489516001600160801b0360608c01511690614d11565b612686896105ee6001600160401b0360206080818501516040516126788482018093604080916001600160801b038151168452602081015160208501520151910152565b606081526106328382611057565b61268f89613828565b9915612a30575061ffff610100870151519251168092149283612a20575b83612a14575b5082612a08575b50816129f8575b50156129d0576126d0826140c0565b6001600160a01b038116156129bd576126ef90969491969592956140db565b956126fa8787614156565b608052600854955f925f965f9560095494600a54965f985b6101008b0151518a10156127dc57908d979695949392918b6127398c610180830151612dd5565b51156127d0576127538c61010063ffffffff930151612dd5565b5116809d6127b8575b50508a60019c8b88828d8c829e6080516020015161ffff169061277e95614297565b91909361012001519061279091612dd5565b51149061279c9161396a565b6127a5916138d6565b986001905b019890919293949596612712565b63ffffffff6127c992168111613930565b5f8c61275c565b509750986001906127aa565b50955096509750979195935097506001600160401b0382166003810290808204600314901517156129a9576801fffffffffffffffe82609f1c16906001600160401b038360a01c1682046002146001600160401b038460a01c161517156129a9576001600160401b039361285792858560a01c1692116139ba565b6128608361502b565b60a01c16915b612998575b50505061287782610e91565b8161294f576001600160401b036020604060a08401938451838101519085600354851c16868316116128fb575b505001516040516128d58382018093604080916001600160801b038151168452602081015160208501520151910152565b606081526128e4608082611057565b51902092510151165f52600560205260405f205590565b85905116851960035416176003557fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6fffffffffffffffff000000000000000060035492851b169116176003555f806128a4565b5061295981610e91565b60018103612981576801000000000000000068ff000000000000000019600454161760045590565b61298a81610e91565b600281036110e45750600290565b6129a1926139fe565b5f808061286b565b634e487b7160e01b5f52601160045260245ffd5b82633a517eed60e21b5f5260045260245ffd5b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b905061018084015151145f6126c1565b5151811491505f6126ba565b5151821492505f6126b3565b61012087015151831493506126ad565b9798909691925f9594955060205f99510151519384519861ffff6101008a0151519251168092149283612c20575b83612c14575b5082612c08575b5081612bf8575b50156129d0575f9893985b878110612bd057505f935f995f965f9b5b6101008a0151518d1015612b6257612aab8d6101808c0151612dd5565b5115612b58578a8a888f8061010084015190612ac691612dd5565b5163ffffffff169c8d8096818097811090612ae0916138f6565b612b1197612b0a9561012094602094612aff94612b40575b5050612dd5565b510151930151612dd5565b511461396a565b6001612b3681986001600160401b036040612b2c8d8c612dd5565b51015116906138d6565b9c5b019b96612a8e565b63ffffffff612b5192168111613930565b5f82612af8565b969b600190612b38565b509650969492995096509691506001600160401b0381166003810290808204600314901517156129a9576001600160401b038316916801fffffffffffffffe8460011b1692808404600214901517156129a957612bc1928492116139ba565b612bca8261502b565b91612866565b97612bee6001916001600160401b036040612b2c8d999e9989612dd5565b9801989398612a7d565b905061018087015151145f612a72565b5151811491505f612a6b565b5151821492505f612a64565b6101208a01515183149350612a5e565b612c4c91503d805f833e612c448183611057565b810190611970565b5f612618565b506040517fccd771d6000000000000000000000000000000000000000000000000000000008152602060048201525f8180612c90602482018b611fbf565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156106e6575f91612cd3575b509761261b565b612ce791503d805f833e612c448183611057565b5f612ccc565b60208091612cfa846110f5565b815201910190612540565b813581526020918201910161250a565b60208091612d228461113f565b8152019101906124d3565b60208091612d3a84611150565b81520191019061245b565b60208091612d52846110f5565b8152019101906123c9565b8135815260209182019101612361565b813581528693506020918201910161229c565b90612d8a826112cd565b612d976040519182611057565b8281528092612da8601f19916112cd565b0190602036910137565b6040516112209190612dc48382611057565b6090815291601f1901366020840137565b8051821015612de95760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b90612e07826140c0565b6001600160a01b03811615612ee757612e1f906140db565b6020612e2b8285614156565b0191612e3b61ffff845116612d80565b93612e4a61ffff855116612d80565b938493612e5b61ffff835116612d80565b9260095491600a54935f5b61ffff8251168b61ffff831691821015612edb57996001600160401b03612ec78380612ec08c9d9e9f828d9e8d9e8d9e8d9e61ffff9e8f9060019f9e612eaf81612eb89a612dd5565b52511691614297565b929096612dd5565b528c612dd5565b911690520116908897969594939291612e66565b50505050505050509150565b509050602060405191612efa8284611057565b5f83525f36813760405192612f0f8385611057565b5f84525f36813760405192612f248185611057565b5f8452505f368137929190565b903590601e198136030182121561023657018035906001600160401b0382116102365760200191813603831361023657565b903590601e198136030182121561023657018035906001600160401b03821161023657602001918160051b3603831361023657565b5f198101919082116129a957565b919082039182116129a957565b90612fbd82611078565b612fca6040519182611057565b8281528092612da8601f1991611078565b90600182811c92168015613009575b6020831014612ff557565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612fea565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff161561304b57565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260405f206001600160a01b0333165f5260205260ff60405f205416156130ac5750565b63e2517d3f60e01b5f523360045260245260445ffd5b905f60e0525f60a0525f60a05260408201519160016001600160401b03602060408651966101408851015160e052838280858901510151995101511660c052015101511601906001600160401b0382116129a95761312160e051614aee565b816137cf576101a08301805151909290156137a4575050518051613144816140c0565b906001600160a01b03821615613792575061315e906140db565b602061316b828451614156565b0161ffff8151169261318b6008549461ffff8660e01c1683519114614cf8565b60ff60208201511680151580613787575b82516131a791614cf8565b6001600160401b0360095495600a5461012052600b5461014052600c546101605260a01c16936131d682612d80565b610100525f5f5b8381106136075750506131fb83516001600160401b03871115614cf8565b613203612db2565b9161320c612db2565b9484519461ffff82511660405196876101208101106001600160401b036101208a011117610fbc57875f806132989481946101206132a79d9b9a99989b0160405284528860208501528060408501528960608501526101005160808501528a60a08501528c60c08501526101405160e0850152610160516101008501528160ff60208c0151169461583d565b969060e05160e0518214614b7f565b5f935b8285106134b657505050505061ffff5f9216915b82811061342c5750505050604051916338f49afb60e01b835260e05160048401526020836024817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49283156106e6575f936133f8575b5060095561012051600a5561014051600b5561016051600c557fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b806008549360a01b16169116176008558060075560e051600655936001600160401b0360019361337e614b13565b6020604086015101525b1660c051146133e257821580806133d7575b156133b357505060400151606060208201519101529190565b6133be575b50509190565b606060406133d0930151015190614bb6565b5f806133b8565b5060e051821461339a565b50606060406133ef614b13565b92015101529190565b9092506020813d602011613424575b8161341460209383611057565b810103126102365751915f61330a565b3d9150613407565b61ffff6134398284612dd5565b5116906134468186612dd5565b5161010083101561347f57826001939161346d859361ffff165f52600f60205260405f2090565b551b6101405117610140525b016132be565b9160ff1981018181116129a9576001936134a6859361ffff165f52600f60205260405f2090565b551b610160511761016052613479565b909192939861ffff6134cc8b6040850151614d00565b5116906001821b906001600160401b036134ea8d6080870151614d00565b5116926134fb8d6060870151614d00565b5190855161ffff8851168210156135f25781602c810204602c14821517156129a957602c820260300190816030116129a957826050602c82028d01015160e01c036135e057506034602c83020181116129a9576054602c83028b01015160c01c90603c602c840201106129a957600195605c602c84028c0101519181145f146135b5575084196101205116610120525b82036135a2575050901916995b01939291906132aa565b5f52600d60205260405f20551799613598565b825f52600e60205260405f20906001600160401b03198254161790558461012051176101205261358b565b634724a0fd60e01b5f5260045260245ffd5b6303e07d4560e61b5f5260045260245260445ffd5b63ffffffff61361a826040880151614d00565b51169161362e8361ffff89511681106138f6565b8161376e575b5081855190848a895161ffff16610120519261364f95614297565b9190978160808801519061366291614d00565b516001600160401b0316988260608901519061367d91614d00565b51610180526001600160401b031692898085148015956020946136b36136f0986136bd966136b895613760575b508d5190614cf8565b612fa6565b614149565b604051639412e6b360e01b81526101805160048201526001600160401b03909a1660248b01529892839081906044820190565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af480156106e6575f9061372e575b600192506137278261010051612dd5565b52016131dd565b506020823d8211613758575b8161374760209383611057565b810103126102365760019151613716565b3d915061373a565b90506101805114155f6136aa565b85516137819163ffffffff168411614cf8565b5f613634565b50601081111561319c565b633a517eed60e21b5f5260045260245ffd5b936001600160401b03919692506137c560206040860151015160e051614bb6565b600160a052613388565b9490926001600160401b03906137e3614b13565b602060408601510152613388565b156137fa575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160401b03602060a08301510151165f52600560205260405f205480155f146138545750505f90565b60408201908151604051613889602082018093604080916001600160801b038151168452602081015160208501520151910152565b60608152613898608082611057565b51902014918215926138b6575b5050156138b157600190565b600290565b6001600160801b0391925060208291015151169151511611155f806138a5565b906001600160401b03809116911601906001600160401b0382116129a957565b156138fe5750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156139385750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156139725750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b908160011b91808304600214901517156129a957565b156139c3575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b92906001600160a01b03613a11856140c0565b16151580613f9e575b613f9857604001516020015151805180159490858015613f8e575b6129d057613a4581969596612d80565b935f967317435cce3d1b4fa2e5f8a08ed921d57c6762a180975b8551811015613b1c5780602080613a79613ac3948a612dd5565b5101516001600160401b036040613a90858c612dd5565b510151604051639412e6b360e01b81526004810193909352166001600160401b0316602482015292839081906044820190565b03818d5af480156106e6575f90613aea575b60019250613ae3828a612dd5565b5201613a5f565b506020823d8211613b14575b81613b0360209383611057565b810103126102365760019151613ad5565b3d9150613af6565b50919490939295965f97613b4e89613b3b613b368a6139a4565b612f98565b9489613b4687612d80565b9c8d92615688565b50602c870290878204602c1417156129a957603001806030116129a9578260051b91838304602014841517156129a957613b94613b8f602494602094614149565b612fb3565b97828901956056875361564160218b01536256414c60228b01536356414c3460238b01538060381c858b01538060301c60258b01538060281c60268b015380841c60278b01538060181c60288b01538060101c60298b01538060081c602a8b0153602b8a01538060081c602c8a0153602d890153604051928380926338f49afb60e01b82528b60048301525af49081156106e6575f91613f5c575b50602e8601528060081c604e860153604f8501536030955f935b8351851015613d4657613c5c8585612dd5565b518887019060ff8760181c16602083015361ffff8760101c16602183015362ffffff8760081c16602283015363ffffffff8716602383015360048a018a116129a957604081015166ffffffffffffff6001600160401b0382169160ff8160381c16602486015361ffff8160301c16602586015362ffffff8160281c16602686015363ffffffff8160201c16602786015364ffffffffff8160181c16602886015365ffffffffffff8160101c16602986015360081c16602a840153602b830153600c8a018a116129a9576020602c910151910152602c88018098116129a95760018895019450613c49565b92509250935f955b8351871015613d7f57613d618785612dd5565b5160208287010152602081018091116129a957600190960195613d4e565b50939291509350613d908184614156565b9160405191613dc160218460208101945f8652845180918484015e81015f838201520301601f198101855284611057565b616000835111613f305750906001600160a01b0391613e79602a7fffff000000000000000000000000000000000000000000000000000000000000845160f01b169360405193849160208301967f6100000000000000000000000000000000000000000000000000000000000000885260218401527f80600a3d393df30000000000000000000000000000000000000000000000000060238401525180918484015e81015f838201520301601f198101835282611057565b51905ff0168015613f08576008549260065560408201516007557fffff0000000000000000000000000000000000000000000000000000000000007dffff00000000000000000000000000000000000000000000000000000000602067ffffffffffffffff60a01b855160a01b1694015160e01b1693161717176008555f6009555f600a555f600b555f600c55565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b90506020813d602011613f86575b81613f7760209383611057565b8101031261023657515f613c2f565b3d9150613f6a565b5060b48111613a35565b50915050565b5061ffff60085460e01c161515613a1a565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f205416155f1461403757805f525f60205260405f206001600160a01b0383165f5260205260405f20600160ff198254161790556001600160a01b03339216907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f2054165f1461403757805f525f60205260405f206001600160a01b0383165f5260205260405f2060ff1981541690556001600160a01b03339216907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b600654036140d7576001600160a01b036008541690565b5f90565b90813b6001811115614106575f1981019081116129a95760016140fd82612fb3565b9360208501903c565b6001600160a01b03837fd8415944000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b90600182018092116129a957565b919082018092116129a957565b9190915f606060405161416881611021565b828152826020820152826040820152015260308351106135e0576356414c34602084015160e01c036135e057602483015160c01c602c840151938460f01c91602e82015191604e81015160f01c604051926141c284611021565b835260208301938585526040840152606083019281845297851595861561428c575b50851561425a575b505083156141ff575b5050506135e05750565b61ffff9192935051925116602c810290808204602c14901517156129a95760300190816030116129a95751621fffe061ffff82169160051b1690808204602014901517156129a95761425091614149565b14155f80806141f5565b9091945060ef1c6201fffe61fffe8216911681036129a9575f190161ffff81116129a95761ffff161415925f806141ec565b60b41095505f6141e4565b939092959461ffff63ffffffff821693168310156143485761ffff1690602c830292808404602c14811517156129a9578360300194856030116129a95784019581605088015160e01c036135e057506001901b90811661432d576034830184116129a957605485015160c01c5b961661431a5750603c01106129a957605c015190565b925050505f52600d60205260405f205490565b815f52600e6020526001600160401b0360405f205416614304565b82856303e07d4560e61b5f5260045260245260445ffd5b999493979198909695995f9661437481611b1e565b80614aad575089151580614aa1575b15614a6a576020019460208b6143f06143a661439e8a6152c8565b923690611242565b916105ee604051858101906143d88287604080916001600160801b038151168452602081015160208501520151910152565b606081526143e7608082611057565b51902091614fe5565b0151808603614a3a57505f905f5b8b8a81831061494d575b505050501561490c5750506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690813b15610236579187926040519384927f87d3a9b100000000000000000000000000000000000000000000000000000000845260648401906004850152606060248501525260848201908960051b9860848a850101928b8a915f603e198d360301925b811061489e5750505050600319848403016044850152808352602083019260208260051b82010193835f92601e1982360301905b85851061472f575050505050505091815f81819503925af180156106e65761471a575b5060018511614525575b50505050506001600160801b0361451f633b9aca00923690611242565b51160490565b6145369096919294969593956152c8565b948335946001600160801b03861680960361471657614554886112cd565b97614562604051998a611057565b885260208801918101903682116147125780979597925b82841061467757505050506001600160401b03829316925b865181101561465a576145a48188612dd5565b519085604051602081019087825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801918a5b81811061462357505050508161460b6001976020614619940151605f19848303016080850152610f7d565b03601f198101835282611057565b5190205d01614591565b91939496509194969760208061464560019360bf198b82030188528951610f7d565b970194019101918c96949392989795986145e0565b5094505050506001600160801b0361451f633b9aca005f80614502565b83989698356001600160401b03811161470e57820160408136031261470e57604051906146a382610feb565b80356001600160401b03811161470a57810136601f8201121561470a576146d19036906020813591016152fe565b825260208101356001600160401b03811161470a57916146f86020949285943691016110c9565b83820152815201930192979597614579565b8880fd5b8680fd5b8480fd5b8380fd5b6147279192505f90611057565b5f905f6144f8565b9193959750919395601f198282030185528735838112156102365784019061475b602082019280615419565b8091936020845252604082019060408160051b8401019380935f915b83831061479e575050505050506020806001929901950195019290918997969495926144d5565b909192939495603f198382030186526147b78783615461565b80356002811015610236576147cb81611b1e565b82526147ee6147dd602083018361544d565b606060208501526060840190615475565b906040810135609e1982360301811215610236576001936020938493614890930191604081830391015261488261486461483961482b8580615389565b60a0865260a0860191615369565b614844878601611150565b151587850152614857604086018661544d565b8482036040860152615475565b9261487160608201611150565b15156060840152608081019061544d565b906080818403910152615475565b980196019493019190614777565b919394965091946083198982030183528535848112156102365760206148fa6001938f839401906148ed6148e36148d58480615419565b6040855260408501916153ba565b9285810190615389565b9185818503910152615369565b9701930191019088969493918e6144a1565b6149496040519283927ffef760c70000000000000000000000000000000000000000000000000000000084526020600485015260248401916153ba565b0390fd5b61496b614965846149839461497294989697986152dc565b80612f63565b36916152fe565b61497d3687896152fe565b906157a9565b15614a305750610b4a61499a6149a4928d8c6152dc565b6020810190612f31565b805182518082149182614a1a575b5050156149c657505060015f808b8a614408565b90614949614a08926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610f7d565b83810360031901602485015290610f7d565b9091506020830120906020840120145f806149b2565b91906001016143fe565b85907f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b897fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b5061ffff8a1115614383565b8760ff602492614abc81611b1e565b614ac581611b1e565b7f112d89cc00000000000000000000000000000000000000000000000000000000835216600452fd5b614aff6001600160a01b03916140c0565b1615614b0d57600754600191565b5f905f90565b60405190614b2082611021565b606082525f6020830152604051614b3681611021565b606081525f60208201525f60408201525f606082015260408301525f60608301528160405190614b67602083611057565b5f825252565b52565b90816020910312610236575190565b15614b88575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b8151511561379257604051917fc87f1f69000000000000000000000000000000000000000000000000000000008352602060048401528260a4810182519060806024840152815180915260c4830190602060c48260051b8601019301915f905b828210614cc95750505050826001600160401b036060614c548594602080980151151560448701526040850151602319878303016064880152611bbf565b92015116608483015203817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49081156106e6575f91614c93575b614c919250808214614b7f565b565b90506020823d602011614cc1575b81614cae60209383611057565b8101031261023657614c91915190614c84565b3d9150614ca1565b91936001919395506020614ce8819260c3198c82030186528851611bbf565b9601920192018794939192614c16565b156135e05750565b906010811015612de95760051b0190565b906001600160801b03633b9aca00911604428111614fb657610708614d368242612fa6565b11614f8757508051805160208201207f000000000000000000000000000000000000000000000000000000000000000003614e945750602081015160ff815116906002549160ff8316928382149283614e7b575b6020015160ff169215614e33575050505063ffffffff606082015116906004549163ffffffff8316808203614e0557505063ffffffff608081920151169160201c16808203614dd7575050565b7f1eb5c1eb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f73b1bde3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6084945060ff90604051947fd382033a000000000000000000000000000000000000000000000000000000008652600486015260081c16602484015260448301526064820152fd5b6020810151600883901c60ff9081169116149350614d8a565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f600154614ed081612fdb565b908160448501526001811690815f14614f635750600114614f03575b506149499192600319848303016024850152610f7d565b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b818310614f495750919291508101606401614949614eec565b805460648488010152859450602090920191600101614f30565b60ff191660648086019190915291151560051b840190910191506149499050614eec565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b7f20daedb6000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b6001600160401b03165f52600560205260405f205480156150035790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040810151519060046020830151926001600160401b0384511690604063ffffffff602087015116950151906020825192015160208063ffffffff83511692015192510151926040519461507e8661103c565b85526020850197885260408501908152606085019182526080850192835260a0850193845261ffff60e0880151169760808801519760a08101519060c081015161012082015161014083015191610180610160850151940151946040519e8f9d8e7f4cc22bb7000000000000000000000000000000000000000000000000000000008152015260248d019d5f5b6008811061529457505060209d506151b28d63ffffffff9761519f829f9d9a95986151ed9f9c9861518c9060c09f9b6001600160401b039a6151656151799261515a8e9c6101248c0190611f63565b6101648a0190611f63565b6102406101a4890152610244880190610f0e565b868103600319016101c488015290610f41565b848103600319016101e486015290610ed5565b9161020460031982850301910152611f8a565b8c8103600319016102248e0152995116895251168b880152516040870152511660608501525160808401525160a08301829052910190610f7d565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af19081156106e6575f9161525a575b501561523257565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90506020813d60201161528c575b8161527560209383611057565b810103126102365761528690611916565b5f61522a565b3d9150615268565b91939597999b9d5091939597999b9d602080600192855181520193019101908f9d9b99979593919e9c9a989694929e61510b565b356001600160401b03811681036102365790565b9190811015612de95760051b81013590603e1981360301821215610236570190565b92919061530a816112cd565b936153186040519586611057565b602085838152019160051b8101918383116102365781905b83821061533e575050505050565b81356001600160401b0381116102365760209161535e87849387016110c9565b815201910190615330565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156102365701602081359101916001600160401b03821161023657813603831361023657565b90602083828152019260208260051b82010193835f925b8484106153e15750505050505090565b909192939495602080615409600193601f198682030188526154038b88615389565b90615369565b98019401940192949391906153d1565b9035601e19823603018112156102365701602081359101916001600160401b038211610236578160051b3603831361023657565b9035607e1982360301811215610236570190565b9035605e1982360301811215610236570190565b6154ae6154936154858380615389565b608086526080860191615369565b6154a06020840184615389565b908583036020870152615369565b6154bb604083018361544d565b90838103604085015281356003811015610236576154d881610e91565b815260208201356003811015610236576154f181610e91565b6020820152604082013590600382101561023657608061552d615548948461551e61553d96999899610e91565b60408501526060810190615389565b9190928160608201520191615369565b926060810190615419565b90916060818503910152808352602083019060208160051b85010193835f915b8383106155785750505050505090565b909192939495601f198282030186526155918784615461565b8035916003831015610236576155ee6020928392856155b1600197610e91565b81526155e06155d56155c586850185615389565b6060888601526060850191615369565b926040810190615389565b916040818503910152615369565b980196019493019190615568565b90615639575080511561561157602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b8151158061567f575b61564a575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15615642565b9493919092946156988483612fa6565b60018114615789576156aa8791615a8b565b92836156cb6156b98289614149565b846156c38961413b565b918a88615688565b6156dc6156fe966156f89299614149565b916156f2613b366156ec8a61413b565b926139a4565b90614149565b93615688565b60405163a5641f6f60e01b8152600481019390935260248301526020826044817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4918215615784575f9261574f575b50614b6d908294612dd5565b614b6d9192506157769060203d60201161577d575b61576e8183611057565b810190614b70565b9190615743565b503d615764565b6106e6565b5061579c9150926157a592949593612dd5565b51928392612dd5565b5290565b908151815103614037575f5b8251811015615813576157c88184612dd5565b51516157d48284612dd5565b51510361580c576157e58184612dd5565b51602081519101206157f78284612dd5565b51602081519101200361580c576001016157b5565b5050505f90565b505050600190565b61ffff60019116019061ffff82116129a957565b5f1981146129a95760010190565b969495919095939293808414615a775761ffff831660908110801590615a6c575b6129d05761586c8884612fa6565b90600182146159da575061587f90615a8b565b9261588a8489614149565b946158a76158978861413b565b956156f2613b366156ec8b61413b565b94815b8b858210806159ad575b156158c857506158c39061582f565b6158aa565b918882969798999a9b9c6158de9695939461583d565b906158ed95919490968a61583d565b60405163a5641f6f60e01b815260048101939093526024830191909152906020816044817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4908115615784575f9161598e575b50809461ffff831660908110801590615983575b6129d05760c061597d9261597661ffff6110e4981661596d8560a0850151612dd5565b9061ffff169052565b0151612dd5565b5261581b565b5061ffff851161594a565b6159a7915060203d60201161577d5761576e8183611057565b5f615936565b508863ffffffff6159d36159c985604060608701510151614d00565b5163ffffffff1690565b16106158b4565b948094979893506159eb915061413b565b1490811591615a41575b50615a2d576110e493929160c087615976615a1a61597d95608061ffff9c0151612dd5565b51998a961661596d8560a0850151612dd5565b8551634724a0fd60e01b5f5260045260245ffd5b9050615a64615a5b6159c984604060608c01510151614d00565b63ffffffff1690565b14155f6159f5565b5061ffff861161585e565b5050949050615a87929350615aa7565b9190565b9060015b8060011b9083821015615aa25750615a8f565b925050565b91909161ffff8311615b8d57610100831015615b5c5760e08101516001841b16615b47575b80519261ffff6040602084015193015116602c810290808204602c14901517156129a95760300190816030116129a9578060051b90808204602014901517156129a957615b1891614149565b90602082018083116129a957815110615b345701602001519150565b83634724a0fd60e01b5f5260045260245ffd5b509061ffff165f52600f60205260405f205490565b61010081015160ff1984018481116129a9576001901b1615615acc57509061ffff165f52600f60205260405f205490565b51634724a0fd60e01b5f5260045260245ffdfea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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
