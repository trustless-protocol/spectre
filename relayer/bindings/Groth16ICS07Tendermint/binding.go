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
	Bin: "0x610140806040523461050a57616697803803809161001d828561062c565b8339810160e08282031261050a576100348261064f565b6100406020840161064f565b61004c6040850161064f565b916100596060860161064f565b60808601519094906001600160401b03811161050a5786019080601f8301121561050a57815161008b92602001610663565b9461009d60c060a0830151920161064f565b958051810190602082019060208184031261050a576020810151906001600160401b03821161050a570191829003601f198101906101201361050a576040519160e083016001600160401b038111848210176105fd5760405260208401516001600160401b03811161050a576020908501019080601f8301121561050a57815161012992602001610663565b82526040811261050a57604080519161014183610611565b61014c8286016106a8565b835261015a606086016106a8565b602084015260208401928352603f19011261050a5760405161017b81610611565b610187608085016106b6565b815261019560a085016106b6565b6020820152604083019081526101ad60c085016106ca565b90606084019182526101c160e086016106ca565b926080850193845261010086015195861515870361050a5760a08601968752610120015194600286101561050a5760c08101958652518051906001600160401b0382116105fd576102136001546106db565b601f81116105ad575b50602090601f83116001146105405763ffffffff95949392915f9183610535575b50508160011b915f199060031b1c1916176001555b5160ff81511661ff00602060025493015160081b169161ffff191617176002555160018060401b0381511660035491602068010000000000000000600160801b0391015160401b169160018060801b031916171760035551169267ffffffff00000000600454925160201b169051151560401b92519160028310156105215769ff00000000000000000068ff00000000000000009360481b169469ff000000000000000000199160018060481b031916171617911617176004556040516103238161031c81610713565b038261062c565b602081519101206101005260405163685e272760e11b8152602060048201526020818061035260248201610713565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af4908115610516575f916104e0575b506101205260405161038e8161031c81610713565b6001600160401b03906020906103a3906107c4565b0151600354916001600160401b03831691168181036104cb575050604090811c6001600160401b03165f908152600560205220556001600160a01b0390811660805290811660a05290811660c0521660e05260045463ffffffff8082169061070882019081116104b75763ffffffff809360201c1692839116116104a257826001600160a01b0381166104895750610439610a9b565b505b604051615ab99081610b1e82396080518161516c015260a0518161438f015260c05181610525015260e0518181816125e60152612ca801526101005181614cb9015261012051816125b10152f35b8061049661049c926109a5565b50610a1b565b5061043b565b6333fae18560e21b5f5260045260245260445ffd5b634e487b7160e01b5f52601160045260245ffd5b6355bace6f60e01b5f5260045260245260445ffd5b90506020813d60201161050e575b816104fb6020938361062c565b8101031261050a57515f610379565b5f80fd5b3d91506104ee565b6040513d5f823e3d90fd5b634e487b7160e01b5f52602160045260245ffd5b015190505f8061023d565b90601f1983169160015f52815f20925f5b818110610595575091600193918563ffffffff99989796941061057d575b505050811b01600155610252565b01515f1960f88460031b161c191690555f808061056f565b92936020600181928786015181550195019301610551565b60015f525f5160206166575f395f51905f52601f840160051c810191602085106105f3575b601f0160051c01905b8181106105e8575061021c565b5f81556001016105db565b90915081906105d2565b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b038211176105fd57604052565b601f909101601f19168101906001600160401b038211908210176105fd57604052565b51906001600160a01b038216820361050a57565b9192916001600160401b0382116105fd576040519161068c601f8201601f19166020018461062c565b82948184528183011161050a578281602093845f96015e010152565b519060ff8216820361050a57565b51906001600160401b038216820361050a57565b519063ffffffff8216820361050a57565b90600182811c92168015610709575b60208310146106f557565b634e487b7160e01b5f52602260045260245ffd5b91607f16916106ea565b6001545f9291610722826106db565b8082529160018116908115610783575060011461073d575050565b60015f9081529293509091905f5160206166575f395f51905f525b838310610769575060209250010190565b600181602092949394548385870101520191019190610758565b9050602093945060ff929192191683830152151560051b010190565b9081518110156107b0570160200190565b634e487b7160e01b5f52603260045260245ffd5b6040516107d081610611565b606081525f60208201525080518015908115610999575b5061098a575f19908051805b610946575b505f1982146109385760018201908183116104b757600360fc1b6001600160f81b0319610825848461079f565b51161480610923575b610913575f5b81518310156108c357610847838361079f565b5160f81c6030811080156108b9575b6108a757600a82026001600160401b03908116602f1990920160ff169190910181169116811061088b57600190920191610834565b509150506040519061089c82610611565b81525f602082015290565b50509150506040519061089c82610611565b5060398111610856565b9150916001811190811591610907575b506108f857604051916108e583610611565b82526001600160401b0316602082015290565b63384c687760e21b5f5260045ffd5b602b915010155f6108d3565b9150506040519061089c82610611565b5080518381039081116104b75760021061082e565b604051915061089c82610611565b5f1981018181116104b757602d60f81b6001600160f81b0319610969838661079f565b511614610980575080156104b7575f1901806107f3565b92505f90506107f8565b6329120bff60e21b5f5260045ffd5b6040915010155f6107e7565b6001600160a01b0381165f9081525f5160206166775f395f51905f52602052604090205460ff16610a16576001600160a01b03165f8181525f5160206166775f395f51905f5260205260408120805460ff191660011790553391905f5160206165d75f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f5160206165f75f395f51905f52602052604090205460ff16610a16576001600160a01b03165f8181525f5160206165f75f395f51905f5260205260408120805460ff191660011790553391905f5160206166375f395f51905f52905f5160206165d75f395f51905f529080a4600190565b5f80525f5160206165f75f395f51905f526020525f5160206166175f395f51905f525460ff16610b19575f8080525f5160206165f75f395f51905f526020525f5160206166175f395f51905f52805460ff1916600117905533905f5160206166375f395f51905f525f5160206165d75f395f51905f528280a4600190565b5f9056fe610120806040526004361015610013575f80fd5b5f3560e01c90816301ffc9a714610da4575080630bece35614610d0a578063248a9ca314610ce05780632f2ff15d14610cb157806336568abe14610c62578063536c2ad314610c0a5780638a8e4c5d14610bd257806391d1485414610b96578063974a74c414610a4b578063a217fddf14610a31578063a6f031bb14610921578063ac9650d81461075c578063d547741f14610726578063ddba65371461023a5763ef913a4b146100c2575f80fd5b34610236575f3660031901126102365760405160208082015261012060408201525f6001546100f081612fd1565b90816101608501526001811690815f1461021157506001146101b0575b6101ac83610198818560ff600254818116606085015260081c1660808301526001600160401b0360035481811660a085015260401c1660c083015260ff60045463ffffffff811660e085015263ffffffff8160201c16610100850152818160401c16151561012085015260481c1661018481611b1e565b61014083015203601f198101835282611057565b604051918291602083526020830190610f7d565b0390f35b91905060015f527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6915f905b8082106101f5575090915081016101800161019861010d565b91926001816020925461018085880101520191019092916101dc565b60ff19166101808086019190915291151560051b84019091019150610198905061010d565b5f80fd5b346102365761024836610e42565b6004549060ff8260401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff16156106f1575b820191602081840312610236578035906001600160401b03821161023657016101208184031261023657604051906102e082610fa1565b80356001600160401b03811161023657846102fc91830161115d565b825260208101356001600160401b03811161023657810193606085820312610236576040519461032b86610fd0565b80356001600160401b038111610236578101604081840312610236576040519061035482610feb565b80356001600160401b03811161023657816103768660209361037e95016110c9565b8452016110f5565b6020820152865260208101356001600160401b03811161023657826103a491830161143a565b602087015260408101356001600160401b038111610236576103c89183910161143a565b6040860152602083019485528060408301906103e391611242565b60408401908152906103f89060a08401611242565b9260608101928484526101000161040e9061122e565b9560808201948786528251915197845191604051998a947fa6fe8f56000000000000000000000000000000000000000000000000000000008652600486016101209052610124860161045f91611b28565b906003198683030160248701528051606083528051606084016040905260a0840161048991610f7d565b90602001516001600160401b0316608084015260208201519083810360208501526104b391611c93565b90604001519180820390604001526104ca91611c93565b83516001600160801b0316604486015260208401516064860152604090930151608485015280516001600160801b031660a4850152602081015160c48501526040015160e48401526001600160801b031661010483015203867f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031691815a93608094fa9586156106e6575f96610652575b50602080610640956105f66801000000000000000099966105a06105ee976001600160801b036001600160401b0398519151935195511690614c83565b6040516105cd8582018093604080916001600160801b038151168452602081015160208501520151910152565b606081526105dc608082611057565b5190206105ee86858a51015116614f57565b808214613735565b6040516106238382018093604080916001600160801b038151168452602081015160208501520151910152565b60608152610632608082611057565b519020940151015116614f57565b68ff0000000000000000191617600455005b9195509160803d6080116106df575b61066b8184611057565b820191608081840312610236576020610640956105f66001600160401b03946105a0680100000000000000009b6001600160801b0386976106c96105ee9b60408051936106b785610feb565b6106c183826118d4565b8552016118d4565b888201529d505097505096945050955050610563565b503d610661565b6040513d5f823e3d90fd5b6106f9613009565b6102a9565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102365761075a61073736610eaf565b90610755610750825f525f602052600160405f20015490565b613078565b613f79565b005b34610236576020366003190112610236576004356001600160401b03811161023657366023820112156102365780600401356001600160401b0381116102365760248201913660248360051b83010111610236576020926040516107c08582611057565b5f815284810191601f1986013684376107d8856112cd565b936107e66040519586611057565b858552601f196107f5876112cd565b01875f5b828110610912575050505f5b868110156108b5576001906108915f808b8861085e61082c60248860051b8b01018b612f27565b90938d6040519483869484860198893784019083820190898252519283915e010185815203601f198101835282611057565b5190305af43d156108ad573d9061087482611078565b916108826040519384611057565b82523d5f8d84013e5b3061556e565b61089b8289612dbe565b526108a68188612dbe565b5001610805565b60609061088b565b604080518981528751818b018190525f92600582901b83018101918a8d01918d9085015b8287106108e65785850386f35b909192938280610902600193603f198a82030186528851610f7d565b96019201960195929190926108d9565b606088820183015281016107f9565b34610236576020366003190112610236576004356001600160401b03811161023657806004019061014060031982360301126102365760ff60045460401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff1615610a24575b6109c36044820183612f59565b916109d16064820185612f59565b93909161010481013590600282101561023657602096610a1c966109f9610124840183612f59565b96909560405198610a0a8c8b611057565b5f8a52608460a48701960135946142d1565b604051908152f35b610a2c613009565b6109b6565b34610236575f3660031901126102365760206040515f8152f35b34610236576020366003190112610236576004356001600160401b03811161023657806004019061016060031982360301126102365760ff60045460401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff1615610b89575b6101448101610aef8184612f27565b905015610b6157610b036044830184612f59565b9091610b126064850186612f59565b95909461010481013591600283101561023657602097610a1c97610b4a96610b51610b41610124870186612f59565b99909886612f27565b3691611093565b98608460a48701960135946142d1565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b610b91613009565b610ae0565b3461023657610ba436610eaf565b905f525f6020526001600160a01b0360405f2091165f52602052602060ff60405f2054166040519015158152f35b3461023657610be036610e42565b50507fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b3461023657602036600319011261023657610c46610c546101ac610c2f600435612de6565b919390604051958695606087526060870190610ed5565b908582036020870152610f0e565b908382036040850152610f41565b3461023657610c7036610eaf565b336001600160a01b03821603610c895761075a91613f79565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102365761075a610cc236610eaf565b90610cdb610750825f525f602052600160405f20015490565b613eec565b34610236576020366003190112610236576020610a1c6004355f525f602052600160405f20015490565b3461023657610d1836610e42565b9060ff60045460401c166106fe575f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060209081527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47549092610d8692909160ff1615610d97576121b6565b60405190610d9381610e91565b8152f35b610d9f613009565b6121b6565b3461023657602036600319011261023657600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361023657817f7965db0b0000000000000000000000000000000000000000000000000000000060209314908115610e18575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483610e11565b906020600319830112610236576004356001600160401b0381116102365782602382011215610236578060040135926001600160401b0384116102365760248483010111610236576024019190565b60031115610e9b57565b634e487b7160e01b5f52602160045260245ffd5b604090600319011261023657600435906024356001600160a01b03811681036102365790565b90602080835192838152019201905f5b818110610ef25750505090565b825163ffffffff16845260209384019390920191600101610ee5565b90602080835192838152019201905f5b818110610f2b5750505090565b8251845260209384019390920191600101610f1e565b90602080835192838152019201905f5b818110610f5e5750505090565b82516001600160401b0316845260209384019390920191600101610f51565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b60a081019081106001600160401b03821117610fbc57604052565b634e487b7160e01b5f52604160045260245ffd5b606081019081106001600160401b03821117610fbc57604052565b604081019081106001600160401b03821117610fbc57604052565b60e081019081106001600160401b03821117610fbc57604052565b608081019081106001600160401b03821117610fbc57604052565b60c081019081106001600160401b03821117610fbc57604052565b90601f801991011681019081106001600160401b03821117610fbc57604052565b6001600160401b038111610fbc57601f01601f191660200190565b92919261109f82611078565b916110ad6040519384611057565b829481845281830111610236578281602093845f960137010152565b9080601f83011215610236578160206110e493359101611093565b90565b359060ff8216820361023657565b35906001600160401b038216820361023657565b91908260409103126102365760405161112181610feb565b602061113a818395611132816110f5565b8552016110f5565b910152565b359063ffffffff8216820361023657565b3590811515820361023657565b80820392916101208412610236576040519161117883611006565b82948135906001600160401b038211610236576111998460409385016110c9565b8552601f190112610236576111e4610100926040516111b781610feb565b6111c3602085016110e7565b81526111d1604085016110e7565b6020820152602086015260608301611109565b60408401526111f560a0820161113f565b606084015261120660c0820161113f565b608084015261121760e08201611150565b60a084015201359060028210156102365760c00152565b35906001600160801b038216820361023657565b91908260609103126102365760405161125a81610fd0565b60408082946112688161122e565b8452602081013560208501520135910152565b8092910391606083126102365760405161129481610feb565b6040819483358352601f1901126102365760209060408051936112b685610feb565b6112c184820161113f565b85520135828401520152565b6001600160401b038111610fbc5760051b60200190565b919060808382031261023657604051906112fd82611021565b819380356001600160401b0381116102365760609261131d9183016110c9565b835260208101356020840152611335604082016110f5565b60408401520135908160070b82036102365760600152565b919091608081840312610236576040519061136782611021565b819381356001600160401b03811161023657820181601f82011215610236578035611391816112cd565b9161139f6040519384611057565b81835260208084019260051b820101918483116102365760208201905b83821061140d575050505083526113d560208301611150565b60208401526040820135906001600160401b03821161023657826114026060949261113a948694016112e4565b6040860152016110f5565b81356001600160401b0381116102365760209161142f888480948801016112e4565b8152019101906113bc565b919060a083820312610236576040519061145382611021565b819380356001600160401b038111610236578101604081840312610236576040519061147e82610feb565b80356001600160401b0381116102365781016102c081860312610236576040519061026082018281106001600160401b03821117610fbc576040526114c38682611109565b825260408101356001600160401b03811161023657866114e49183016110c9565b60208301526114f5606082016110f5565b60408301526115066080820161122e565b606083015261151760a08201611150565b60808301526115298660c0830161127b565b60a083015261153b6101208201611150565b60c083015261014081013560e08301526115586101608201611150565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526115a76102208201611150565b6101c08301526102408101356101e08301526115c66102608201611150565b6102008301526102808101356102208301526102a0810135906001600160401b038211610236576115f9918791016110c9565b61024082015282526020810135906001600160401b038211610236570160c081850312610236576040519061162d82611021565b611636816110f5565b82526116446020820161113f565b6020830152611656856040830161127b565b604083015260a0810135906001600160401b038211610236570184601f8201121561023657803590611687826112cd565b916116956040519384611057565b80835260208084019160051b830101918783116102365760208101915b838310611720575050505060608201526020820152835260208101356001600160401b03811161023657826116e891830161134d565b60208401526116fa8260408301611109565b60408401526080810135916001600160401b0383116102365760609261113a920161134d565b82356001600160401b038111610236578201906040828b03601f190112610236576040519161174e83610feb565b6020810135600481101561023657835260408101356001600160401b038111610236576020910101906080828c0312610236576040519261178e84611021565b82356001600160401b038111610236578c6117aa9185016110c9565b84526117b86020840161122e565b60208501526117c960408401611150565b60408501526060830135936001600160401b038511610236576117f18d6020968796016110c9565b6060820152838201528152019201916116b2565b9080601f83011215610236576040805192906118219084611057565b82906040810192831161023657905b82821061183d5750505090565b8135815260209182019101611830565b9080601f83011215610236578135611864816112cd565b926118726040519485611057565b81845260208085019260051b82010192831161023657602001905b82821061189a5750505090565b602080916118a78461113f565b81520191019061188d565b519060ff8216820361023657565b51906001600160401b038216820361023657565b9190826040910312610236576040516118ec81610feb565b602061113a8183956118fd816118c0565b8552016118c0565b519063ffffffff8216820361023657565b5190811515820361023657565b51906001600160801b038216820361023657565b91908260609103126102365760405161194f81610fd0565b604080829461195d81611923565b8452602081015160208501520151910152565b602081830312610236578051906001600160401b03821161023657016101808183031261023657604051916119a48361103c565b81516001600160401b0381116102365782019182820392610120841261023657604051936119d185611006565b81516001600160401b0381116102365782019084601f830112156102365781516119fa81611078565b90611a086040519283611057565b8082528660208286010111610236576020815f9282604097018386015e830101528652601f1901126102365761010090604051611a4481610feb565b611a50602083016118b2565b8152611a5e604083016118b2565b60208201526020860152611a7584606083016118d4565b6040860152611a8660a08201611905565b6060860152611a9760c08201611905565b6080860152611aa860e08201611916565b60a0860152015190600282101561023657836101409260c0611b169601528552611ad58360208301611937565b6020860152611ae78360808301611937565b6040860152611af860e08201611923565b6060860152611b0b8361010083016118d4565b6080860152016118d4565b60a082015290565b60021115610e9b57565b9061010060c0611b4384516101208552610120850190610f7d565b9360ff60208083015182815116828801520151166040850152611b83604082015160608601906001600160401b0360208092828151168552015116910152565b63ffffffff60608201511660a085015263ffffffff6080820151168285015260a0810151151560e0850152015191611bba83611b1e565b015290565b90606080611bd68451608085526080850190610f7d565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b828210611c6857505050506001600160401b036060611c5e819360208701511515602087015260408701518682036040880152611bbf565b9401511691015290565b90919293602080611c85600193609f198a82030186528851611bbf565b960192019201909291611c26565b91909180519260a081526020611df38551604060a0850152611cce60e0850182516001600160401b0360208092828151168552015116910152565b610240611cec848301516102c06101208801526103a0870190610f7d565b60408301516001600160401b031661014087015260608301516001600160801b03166101608701526080830151151561018087015260a083015180516101a0880152602090810151805163ffffffff166101c089015201516101e08701529160c0810151151561020087015260e08101516102208701526101008101511515828701526101208101516102608701526101408101516102808701526101608101516102a08701526101808101516102c08701526101a08101516102e08701526101c081015115156103008701526101e08101516103208701526102008101511515610340870152610220810151610360870152015160df1985830301610380860152610f7d565b94015193609f198282030160c0830152606060c08201956001600160401b03815116835263ffffffff6020820151166020840152611e536040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519460c060a0830152855180915260e0820190602060e08260051b8501019701925f905b828210611ed85750505050506060611ea06110e4949560208501518482036020860152611bff565b92611ec8604082015160408501906001600160401b0360208092828151168552015116910152565b0151906080818403910152611bff565b909192939760df1982820301855288519081516004811015610e9b57611f568260206001958195948295520151906040848201526060611f2483516080604085015260c0840190610f7d565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f1982850301910152610f7d565b9a01950193920190611e78565b905f905b60028210611f7457505050565b6020806001928551815201930191019091611f67565b90602080835192838152019201905f5b818110611fa75750505090565b82511515845260209384019390920191600101611f9a565b90612019611fd883516109408452610940840190611b28565b61200760208501516020850190604080916001600160801b038151168452602081015160208501520151910152565b60408401518382036080850152611c93565b60608301516001600160801b031660a083015260808301515f60c084015b600882106121a0575050506120ee6120da6120c66120b261209e6101a09561206860a08a01516101c08a0190611f63565b61207b60c08a01516102008a0190611f63565b61ffff60e08a0151166102408901526101008901518882036102608a0152610ed5565b610120880151878203610280890152610f0e565b6101408701518682036102a0880152610f41565b6101608601518582036102c0870152610ed5565b6101808501518482036102e0860152611f8a565b9201518051610300830152602081015160ff1661032083015260408101515f61034084015b601082106121845750505060608101515f61054084015b6010821061216e5750505060800151905f90610740015b6010821061214f5750505090565b6020806001926001600160401b03865116815201930191019091612141565b602080600192855181520193019101909161212a565b60208060019263ffffffff865116815201930191019091612113565b6020806001928551815201930191019091612037565b9190820191602081840312610236578035906001600160401b03821161023657019081830391610940831261023657604051906101c082018281106001600160401b03821117610fbc5760405280356001600160401b038111610236578561221f91830161115d565b825261222e8560208301611242565b602083015260808101356001600160401b038111610236578561225291830161143a565b946040830195865261226660a0830161122e565b60608401528060df830112156102365760405161228561010082611057565b806101c08401918383116102365790839160c08601905b848210612d7957505060808601526122b391611805565b60a08401526122c6816102008401611805565b60c08401526102408201359061ffff821682036102365760e084019182526102608301356001600160401b038111610236578161230491850161184d565b6101008501526102808301356001600160401b03811161023657830181601f8201121561023657803590612337826112cd565b916123456040519384611057565b80835260208084019160051b8301019184831161023657602001905b828210612d69575050506101208501526102a08301356001600160401b0381116102365783019281601f850112156102365783359361239f856112cd565b946123ad6040519687611057565b80865260208087019160051b8301019184831161023657602001905b828210612d515750505061014085019384526102c08101356001600160401b03811161023657826123fb91830161184d565b9661016086019788526102e08201356001600160401b03811161023657820183601f8201121561023657803590612431826112cd565b9161243f6040519384611057565b80835260208084019160051b8301019186831161023657602001905b828210612d39575050506101808701526106406102ff1990910112610236576040519161248783610fa1565b610300820135835261249c61032083016110e7565b60208401528061035f8301121561023657610200916040516124be8482611057565b80610540830191848311610236576103408401905b838210612d2157505060408601528261055f83011215610236576040516124fa8582611057565b8061074084019285841161023657905b838210612d1157505060608601528261075f83011215610236576125316040519485611057565b61094084920192831161023657905b828210612cf95750505060808201526101a084015261255e836130b8565b93919392909660a0519497885f14612c5e575f6040518080937f127dd052000000000000000000000000000000000000000000000000000000008252606060048301526125ae606483018c611fbf565b907f00000000000000000000000000000000000000000000000000000000000000006024840152604483015203816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156106e6575f91612c3c575b50985b6126328a516001600160801b0360608d01511690614c83565b6126848a6105ee6001600160401b0360206080818501516040516126768482018093604080916001600160801b038151168452602081015160208501520151910152565b606081526106328382611057565b61268d8a61376c565b9815612a3d575061ffff610100870151519251168092149283612a2d575b83612a21575b5082612a15575b5081612a05575b50156129dd576126ce82613ffc565b946001600160a01b038616156129ca579390612700600897939754966126fa61ffff8960e01c16614032565b90614057565b9661270b88826140c1565b6080525f925f965f9560095494600a54965f985b6101008b0151518a10156127e957908d979695949392918b6127468c610180830151612dbe565b51156127dd576127608c61010063ffffffff930151612dbe565b5116809d6127c5575b50508a60019c8b88828d8c829e6080516020015161ffff169061278b95614209565b91909361012001519061279d91612dbe565b5114906127a9916138ae565b6127b29161381a565b986001905b01989091929394959661271f565b63ffffffff6127d692168111613874565b5f8c612769565b509750986001906127b7565b50955096509750979195935097506001600160401b0382166003810290808204600314901517156129b6576801fffffffffffffffe82609f1c16906001600160401b038360a01c1682046002146001600160401b038460a01c161517156129b6576001600160401b039361286492858560a01c1692116138fe565b61286d83614f9d565b60a01c16915b6129a5575b50505061288482610e91565b8161295c576001600160401b036020604060a08401938451838101519085600354851c1686831611612908575b505001516040516128e28382018093604080916001600160801b038151168452602081015160208501520151910152565b606081526128f1608082611057565b51902092510151165f52600560205260405f205590565b85905116851960035416176003557fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff6fffffffffffffffff000000000000000060035492851b169116176003555f806128b1565b5061296681610e91565b6001810361298e576801000000000000000068ff000000000000000019600454161760045590565b61299781610e91565b600281036110e45750600290565b6129ae92613942565b5f8080612878565b634e487b7160e01b5f52601160045260245ffd5b82633a517eed60e21b5f5260045260245ffd5b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b905061018084015151145f6126bf565b5151811491505f6126b8565b5151821492505f6126b1565b61012087015151831493506126ab565b97909691925f9594955060205f99510151519384519861ffff6101008a0151519251168092149283612c2c575b83612c20575b5082612c14575b5081612c04575b50156129dd575f9893985b878110612bdc57505f935f995f965f9b5b6101008a0151518d1015612b6e57612ab78d6101808c0151612dbe565b5115612b64578a8a888f8061010084015190612ad291612dbe565b5163ffffffff169c8d8096818097811090612aec9161383a565b612b1d97612b169561012094602094612b0b94612b4c575b5050612dbe565b510151930151612dbe565b51146138ae565b6001612b4281986001600160401b036040612b388d8c612dbe565b510151169061381a565b9c5b019b96612a9a565b63ffffffff612b5d92168111613874565b5f82612b04565b969b600190612b44565b509650969492995096509691506001600160401b0381166003810290808204600314901517156129b6576001600160401b038316916801fffffffffffffffe8460011b1692808404600214901517156129b657612bcd928492116138fe565b612bd682614f9d565b91612873565b97612bfa6001916001600160401b036040612b388d999e9989612dbe565b9801989398612a89565b905061018087015151145f612a7e565b5151811491505f612a77565b5151821492505f612a70565b6101208a01515183149350612a6a565b612c5891503d805f833e612c508183611057565b810190611970565b5f612616565b506040517fccd771d6000000000000000000000000000000000000000000000000000000008152602060048201525f8180612c9c602482018b611fbf565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156106e6575f91612cdf575b5098612619565b612cf391503d805f833e612c508183611057565b5f612cd8565b60208091612d06846110f5565b815201910190612540565b813581526020918201910161250a565b60208091612d2e8461113f565b8152019101906124d3565b60208091612d4684611150565b81520191019061245b565b60208091612d5e846110f5565b8152019101906123c9565b8135815260209182019101612361565b813581528593506020918201910161229c565b90612d96826112cd565b612da36040519182611057565b8281528092612db4601f19916112cd565b0190602036910137565b8051821015612dd25760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b90612df082613ffc565b6001600160a01b03811615612edd57612e15906126fa61ffff60085460e01c16614032565b6020612e2182856140c1565b0191612e3161ffff845116612d8c565b93612e4061ffff855116612d8c565b938493612e5161ffff835116612d8c565b9260095491600a54935f5b61ffff8251168b61ffff831691821015612ed157996001600160401b03612ebd8380612eb68c9d9e9f828d9e8d9e8d9e8d9e61ffff9e8f9060019f9e612ea581612eae9a612dbe565b52511691614209565b929096612dbe565b528c612dbe565b911690520116908897969594939291612e5c565b50505050505050509150565b509050602060405191612ef08284611057565b5f83525f36813760405192612f058385611057565b5f84525f36813760405192612f1a8185611057565b5f8452505f368137929190565b903590601e198136030182121561023657018035906001600160401b0382116102365760200191813603831361023657565b903590601e198136030182121561023657018035906001600160401b03821161023657602001918160051b3603831361023657565b5f198101919082116129b657565b919082039182116129b657565b90612fb382611078565b612fc06040519182611057565b8281528092612db4601f1991611078565b90600182811c92168015612fff575b6020831014612feb57565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612fe0565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff161561304157565b63e2517d3f60e01b5f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260405f206001600160a01b0333165f5260205260ff60405f205416156130a25750565b63e2517d3f60e01b5f523360045260245260445ffd5b905f60a0525f906040830151928351906101408251015160a05260016001600160401b0360206040828180848801510151975101511698015101511601916001600160401b0383116129b65761310f60a051614a60565b81613713576101a08401805151909290156136e95750505180519061313382613ffc565b916001600160a01b0383169081156136d75750823b9060018211156136c557505f1981019081116129b657600161316982612fa9565b9360208501903c602061317d8383516140c1565b0160c05261ffff60c0515116916131a26008549361ffff8560e01c1684519114614c6a565b60ff602083015116801515806136ba575b83516131be91614c6a565b6001600160401b0360095494600a5460e05260a01c16926131de82612d8c565b5f5f5b84811061353a5750506131ff82516001600160401b03871115614c6a565b81519061ffff60c0515116916040519061010082018281106001600160401b03821117610fbc5761328494613279945f93849360409b9a999b52855288602086015281604086015289606086015260808501528a60a085015260e05160c085015260e0518b1760e08501528160ff60208b0151169461579b565b60a051818114614af1565b5f925b8184106133e65750505050604051916338f49afb60e01b835260a05160048401526020836024817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49283156106e6575f936133b2575b5060095560e051600a557fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff67ffffffffffffffff60a01b806008549360a01b16169116176008558060075560a051600655956001600160401b03600194613338614a85565b6020604087015101525b161461339b5782158080613390575b1561336b5750506040015160606020820151910152929190565b613377575b5050929190565b60606040613389930151015190614b28565b5f80613370565b5060a0518214613351565b50606060406133a8614a85565b9201510152929190565b9092506020813d6020116133de575b816133ce60209383611057565b810103126102365751915f6132d3565b3d91506133c1565b9091929461ffff6133fb876040850151614c72565b5116906001821b906001600160401b03613419896080870151614c72565b51169261342a896060870151614c72565b5190855161ffff60c05151168210156135255781602c810204602c14821517156129b657602c82026030016030116129b657816050602c82028b01015160e01c0361351357506034602c820201602c8202603001116129b657602c81028881016054015160c01c90603c81016030909101116129b657600195605c602c84028b0101519181145f146134ea5750841960e0511660e0525b82036134d7575050901916955b01929190613287565b5f52600b60205260405f205517956134ce565b825f52600c60205260405f20906001600160401b03198254161790558460e0511760e0526134c1565b634724a0fd60e01b5f5260045260245ffd5b6303e07d4560e61b5f5260045260245260445ffd5b63ffffffff61354d826040870151614c72565b5116916135638361ffff60c0515116811061383a565b816136a1575b5081845190878a60c0515161ffff1660e0519261358595614209565b9190978160808701519061359891614c72565b516001600160401b031698826060880151906135b391614c72565b51610100526001600160401b031692898085148015956020946135e9613626986135f3966135ee95613693575b508c5190614c6a565b612f9c565b614025565b604051639412e6b360e01b81526101005160048201526001600160401b03909a1660248b01529892839081906044820190565b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af480156106e6575f90613661575b6001925061365a8286612dbe565b52016131e1565b506020823d821161368b575b8161367a60209383611057565b81010312610236576001915161364c565b3d915061366d565b90506101005114155f6135e0565b84516136b49163ffffffff168411614c6a565b5f613569565b5060108111156131b3565b633610565160e21b5f5260045260245ffd5b633a517eed60e21b5f5260045260245ffd5b94965096905061370360206040850151015160a051614b28565b6001600160401b03600196613342565b9690936001600160401b0390613727614a85565b602060408701510152613342565b1561373e575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160401b03602060a08301510151165f52600560205260405f205480155f146137985750505f90565b604082019081516040516137cd602082018093604080916001600160801b038151168452602081015160208501520151910152565b606081526137dc608082611057565b51902014918215926137fa575b5050156137f557600190565b600290565b6001600160801b0391925060208291015151169151511611155f806137e9565b906001600160401b03809116911601906001600160401b0382116129b657565b156138425750565b63ffffffff907fc9d6366c000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b1561387c5750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156138b65750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b908160011b91808304600214901517156129b657565b15613907575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b92906001600160a01b0361395585613ffc565b16151580613eda575b613ed457604001516020015151805180159490858015613eca575b6129dd5761398981969596612d8c565b935f967317435cce3d1b4fa2e5f8a08ed921d57c6762a180975b8551811015613a6057806020806139bd613a07948a612dbe565b5101516001600160401b0360406139d4858c612dbe565b510151604051639412e6b360e01b81526004810193909352166001600160401b0316602482015292839081906044820190565b03818d5af480156106e6575f90613a2e575b60019250613a27828a612dbe565b52016139a3565b506020823d8211613a58575b81613a4760209383611057565b810103126102365760019151613a19565b3d9150613a3a565b50919490939295965f97613a9289613a7f613a7a8a6138e8565b612f8e565b9489613a8a87612d8c565b9c8d926155fa565b50602c870290878204602c1417156129b657603001806030116129b6578260051b91838304602014841517156129b657613ad8613ad3602494602094614025565b612fa9565b97828901956056875361564160218b01536256414c60228b01536356414c3460238b01538060381c858b01538060301c60258b01538060281c60268b015380841c60278b01538060181c60288b01538060101c60298b01538060081c602a8b0153602b8a01538060081c602c8a0153602d890153604051928380926338f49afb60e01b82528b60048301525af49081156106e6575f91613e98575b50602e8601528060081c604e860153604f8501536030955f935b8351851015613c8a57613ba08585612dbe565b518887019060ff8760181c16602083015361ffff8760101c16602183015362ffffff8760081c16602283015363ffffffff8716602383015360048a018a116129b657604081015166ffffffffffffff6001600160401b0382169160ff8160381c16602486015361ffff8160301c16602586015362ffffff8160281c16602686015363ffffffff8160201c16602786015364ffffffffff8160181c16602886015365ffffffffffff8160101c16602986015360081c16602a840153602b830153600c8a018a116129b6576020602c910151910152602c88018098116129b65760018895019450613b8d565b92509250935f955b8351871015613cc357613ca58785612dbe565b5160208287010152602081018091116129b657600190960195613c92565b50939291509350613cd481846140c1565b9160405191613d0560218460208101945f8652845180918484015e81015f838201520301601f198101855284611057565b616000835111613e6c5750906001600160a01b0391613dbd602a7fffff000000000000000000000000000000000000000000000000000000000000845160f01b169360405193849160208301967f6100000000000000000000000000000000000000000000000000000000000000885260218401527f80600a3d393df30000000000000000000000000000000000000000000000000060238401525180918484015e81015f838201520301601f198101835282611057565b51905ff0168015613e44576008549260065560408201516007557fffff0000000000000000000000000000000000000000000000000000000000007dffff00000000000000000000000000000000000000000000000000000000602067ffffffffffffffff60a01b855160a01b1694015160e01b1693161717176008555f6009555f600a55565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b90506020813d602011613ec2575b81613eb360209383611057565b8101031261023657515f613b73565b3d9150613ea6565b5060b48111613979565b50915050565b5061ffff60085460e01c16151561395e565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f205416155f14613f7357805f525f60205260405f206001600160a01b0383165f5260205260405f20600160ff198254161790556001600160a01b03339216907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f2054165f14613f7357805f525f60205260405f206001600160a01b0383165f5260205260405f2060ff1981541690556001600160a01b03339216907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b60065403614013576001600160a01b036008541690565b5f90565b90600182018092116129b657565b919082018092116129b657565b61ffff16602c810290808204602c14901517156129b657603001806030116129b65790565b9190823b60018111156140a5575f1981019081116129b657811161408957600161408082612fa9565b9360208501903c565b6001600160a01b0383633610565160e21b5f521660045260245ffd5b6001600160a01b0384633610565160e21b5f521660045260245ffd5b9190915f60606040516140d381611021565b82815282602082015282604082015201526030835110613513576356414c34602084015160e01c0361351357602483015160c01c92602c810151938460f01c91602e81015190604e81015160f01c6040519361412e85611021565b84526020840192858452604085015260608401938185529785159586156141fe575b5085156141cc575b5050831561416b575b5050506135135750565b519051915192509061ffff838116916141849116614032565b9081831493841561419d575b50505050155f8080614161565b621fffe0919293945060051b1690808204602014901517156129b6576141c291614025565b145f808080614190565b9091945060ef1c6201fffe61fffe8216911681036129b6575f190161ffff81116129b65761ffff161415925f80614158565b60b41095505f614150565b939092959461ffff63ffffffff821693168310156142ba5761ffff1690602c830292808404602c14811517156129b6578360300194856030116129b65784019581605088015160e01c0361351357506001901b90811661429f576034830184116129b657605485015160c01c5b961661428c5750603c01106129b657605c015190565b925050505f52600b60205260405f205490565b815f52600c6020526001600160401b0360405f205416614276565b82856303e07d4560e61b5f5260045260245260445ffd5b999493979198909695995f966142e681611b1e565b80614a1f575089151580614a13575b156149dc576020019460208b6143626143186143108a61523a565b923690611242565b916105ee6040518581019061434a8287604080916001600160801b038151168452602081015160208501520151910152565b60608152614359608082611057565b51902091614f57565b01518086036149ac57505f905f5b8b8a8183106148bf575b505050501561487e5750506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690813b15610236579187926040519384927f87d3a9b100000000000000000000000000000000000000000000000000000000845260648401906004850152606060248501525260848201908960051b9860848a850101928b8a915f603e198d360301925b81106148105750505050600319848403016044850152808352602083019260208260051b82010193835f92601e1982360301905b8585106146a1575050505050505091815f81819503925af180156106e65761468c575b5060018511614497575b50505050506001600160801b03614491633b9aca00923690611242565b51160490565b6144a890969192949695939561523a565b948335946001600160801b038616809603614688576144c6886112cd565b976144d4604051998a611057565b885260208801918101903682116146845780979597925b8284106145e957505050506001600160401b03829316925b86518110156145cc576145168188612dbe565b519085604051602081019087825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801918a5b81811061459557505050508161457d600197602061458b940151605f19848303016080850152610f7d565b03601f198101835282611057565b5190205d01614503565b9193949650919496976020806145b760019360bf198b82030188528951610f7d565b970194019101918c9694939298979598614552565b5094505050506001600160801b03614491633b9aca005f80614474565b83989698356001600160401b038111614680578201604081360312614680576040519061461582610feb565b80356001600160401b03811161467c57810136601f8201121561467c57614643903690602081359101615270565b825260208101356001600160401b03811161467c579161466a6020949285943691016110c9565b838201528152019301929795976144eb565b8880fd5b8680fd5b8480fd5b8380fd5b6146999192505f90611057565b5f905f61446a565b9193959750919395601f19828203018552873583811215610236578401906146cd60208201928061538b565b8091936020845252604082019060408160051b8401019380935f915b83831061471057505050505050602080600192990195019501929091899796949592614447565b909192939495603f1983820301865261472987836153d3565b803560028110156102365761473d81611b1e565b825261476061474f60208301836153bf565b6060602085015260608401906153e7565b906040810135609e198236030181121561023657600193602093849361480293019160408183039101526147f46147d66147ab61479d85806152fb565b60a0865260a08601916152db565b6147b6878601611150565b1515878501526147c960408601866153bf565b84820360408601526153e7565b926147e360608201611150565b1515606084015260808101906153bf565b9060808184039101526153e7565b9801960194930191906146e9565b9193949650919460831989820301835285358481121561023657602061486c6001938f8394019061485f614855614847848061538b565b60408552604085019161532c565b92858101906152fb565b91858185039101526152db565b9701930191019088969493918e614413565b6148bb6040519283927ffef760c700000000000000000000000000000000000000000000000000000000845260206004850152602484019161532c565b0390fd5b6148dd6148d7846148f5946148e4949896979861524e565b80612f59565b3691615270565b6148ef368789615270565b9061571b565b156149a25750610b4a61490c614916928d8c61524e565b6020810190612f27565b80518251808214918261498c575b50501561493857505060015f808b8a61437a565b906148bb61497a926040519384937f5f1ca381000000000000000000000000000000000000000000000000000000008552604060048601526044850190610f7d565b83810360031901602485015290610f7d565b9091506020830120906020840120145f80614924565b9190600101614370565b85907f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b897fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b5061ffff8a11156142f5565b8760ff602492614a2e81611b1e565b614a3781611b1e565b7f112d89cc00000000000000000000000000000000000000000000000000000000835216600452fd5b614a716001600160a01b0391613ffc565b1615614a7f57600754600191565b5f905f90565b60405190614a9282611021565b606082525f6020830152604051614aa881611021565b606081525f60208201525f60408201525f606082015260408301525f60608301528160405190614ad9602083611057565b5f825252565b52565b90816020910312610236575190565b15614afa575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b815151156136d757604051917fc87f1f69000000000000000000000000000000000000000000000000000000008352602060048401528260a4810182519060806024840152815180915260c4830190602060c48260051b8601019301915f905b828210614c3b5750505050826001600160401b036060614bc68594602080980151151560448701526040850151602319878303016064880152611bbf565b92015116608483015203817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49081156106e6575f91614c05575b614c039250808214614af1565b565b90506020823d602011614c33575b81614c2060209383611057565b8101031261023657614c03915190614bf6565b3d9150614c13565b91936001919395506020614c5a819260c3198c82030186528851611bbf565b9601920192018794939192614b88565b156135135750565b906010811015612dd25760051b0190565b906001600160801b03633b9aca00911604428111614f2857610708614ca88242612f9c565b11614ef957508051805160208201207f000000000000000000000000000000000000000000000000000000000000000003614e065750602081015160ff815116906002549160ff8316928382149283614ded575b6020015160ff169215614da5575050505063ffffffff606082015116906004549163ffffffff8316808203614d7757505063ffffffff608081920151169160201c16808203614d49575050565b7f1eb5c1eb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f73b1bde3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6084945060ff90604051947fd382033a000000000000000000000000000000000000000000000000000000008652600486015260081c16602484015260448301526064820152fd5b6020810151600883901c60ff9081169116149350614cfc565b604051907ff6b6676b00000000000000000000000000000000000000000000000000000000825260406004830152815f600154614e4281612fd1565b908160448501526001811690815f14614ed55750600114614e75575b506148bb9192600319848303016024850152610f7d565b60015f90815292939291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b818310614ebb57509192915081016064016148bb614e5e565b805460648488010152859450602090920191600101614ea2565b60ff191660648086019190915291151560051b840190910191506148bb9050614e5e565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b7f20daedb6000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b6001600160401b03165f52600560205260405f20548015614f755790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040810151519060046020830151926001600160401b0384511690604063ffffffff602087015116950151906020825192015160208063ffffffff835116920151925101519260405194614ff08661103c565b85526020850197885260408501908152606085019182526080850192835260a0850193845261ffff60e0880151169760808801519760a08101519060c081015161012082015161014083015191610180610160850151940151946040519e8f9d8e7f4cc22bb7000000000000000000000000000000000000000000000000000000008152015260248d019d5f5b6008811061520657505060209d506151248d63ffffffff97615111829f9d9a959861515f9f9c986150fe9060c09f9b6001600160401b039a6150d76150eb926150cc8e9c6101248c0190611f63565b6101648a0190611f63565b6102406101a4890152610244880190610f0e565b868103600319016101c488015290610f41565b848103600319016101e486015290610ed5565b9161020460031982850301910152611f8a565b8c8103600319016102248e0152995116895251168b880152516040870152511660608501525160808401525160a08301829052910190610f7d565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af19081156106e6575f916151cc575b50156151a457565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90506020813d6020116151fe575b816151e760209383611057565b81010312610236576151f890611916565b5f61519c565b3d91506151da565b91939597999b9d5091939597999b9d602080600192855181520193019101908f9d9b99979593919e9c9a989694929e61507d565b356001600160401b03811681036102365790565b9190811015612dd25760051b81013590603e1981360301821215610236570190565b92919061527c816112cd565b9361528a6040519586611057565b602085838152019160051b8101918383116102365781905b8382106152b0575050505050565b81356001600160401b038111610236576020916152d087849387016110c9565b8152019101906152a2565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156102365701602081359101916001600160401b03821161023657813603831361023657565b90602083828152019260208260051b82010193835f925b8484106153535750505050505090565b90919293949560208061537b600193601f198682030188526153758b886152fb565b906152db565b9801940194019294939190615343565b9035601e19823603018112156102365701602081359101916001600160401b038211610236578160051b3603831361023657565b9035607e1982360301811215610236570190565b9035605e1982360301811215610236570190565b6154206154056153f783806152fb565b6080865260808601916152db565b61541260208401846152fb565b9085830360208701526152db565b61542d60408301836153bf565b908381036040850152813560038110156102365761544a81610e91565b8152602082013560038110156102365761546381610e91565b6020820152604082013590600382101561023657608061549f6154ba94846154906154af96999899610e91565b604085015260608101906152fb565b91909281606082015201916152db565b92606081019061538b565b90916060818503910152808352602083019060208160051b85010193835f915b8383106154ea5750505050505090565b909192939495601f1982820301865261550387846153d3565b803591600383101561023657615560602092839285615523600197610e91565b8152615552615547615537868501856152fb565b60608886015260608501916152db565b9260408101906152fb565b9160408185039101526152db565b9801960194930191906154da565b906155ab575080511561558357602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b815115806155f1575b6155bc575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b156155b4565b94939190929461560a8483612f9c565b600181146156fb5761561c87916159fd565b928361563d61562b8289614025565b8461563589614017565b918a886155fa565b61564e6156709661566a9299614025565b91615664613a7a61565e8a614017565b926138e8565b90614025565b936155fa565b60405163a5641f6f60e01b8152600481019390935260248301526020826044817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49182156156f6575f926156c1575b50614adf908294612dbe565b614adf9192506156e89060203d6020116156ef575b6156e08183611057565b810190614ae2565b91906156b5565b503d6156d6565b6106e6565b5061570e91509261571792949593612dbe565b51928392612dbe565b5290565b908151815103613f73575f5b82518110156157855761573a8184612dbe565b51516157468284612dbe565b51510361577e576157578184612dbe565b51602081519101206157698284612dbe565b51602081519101200361577e57600101615727565b5050505f90565b505050600190565b5f1981146129b65760010190565b93929094918284148080916159e1575b6159b95761ffff83116129dd576157c28783612f9c565b90600182146158d557506157d5906159fd565b936157fc6157e38689614025565b95615664613a7a61565e6157f688614017565b97614017565b92815b85811087828a836158a4575b505050156158215761581c9061578d565b6157ff565b9495969791859188615833948b61579b565b9461583e959661579b565b60405163a5641f6f60e01b81526004810192909252602482015260208180604481015b03817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49081156156f6575f9161588b575090565b6110e4915060203d6020116156ef576156e08183611057565b63ffffffff9293506158cb91604060606158c19301510151614c72565b5163ffffffff1690565b161087828a61580b565b92505050949392941561594e57505061591d836020926158619495519184810151615905604083015161ffff1690565b9063ffffffff60c060a0850151940151941694614209565b6040519384928392639412e6b360e01b8452600484019092916001600160401b036020916040840195845216910152565b615959829392614017565b149081159161598e575b5061597a57608061597692930151612dbe565b5190565b8251634724a0fd60e01b5f5260045260245ffd5b90506159b16159a86158c184604060608901510151614c72565b63ffffffff1690565b14155f615963565b5050929150506110e492508051906159db6040602083015192015161ffff1690565b91615a3a565b506159f86159f4838960e08a0151615a19565b1590565b6157ab565b9060015b8060011b9083821015615a145750615a01565b925050565b91615a2682600192612f9c565b1b905f1982019182116129b6571b16151590565b9392909161ffff16602c810290808204602c14901517156129b65760300190816030116129b6578060051b90808204602014901517156129b657615a7d91614025565b90602082018083116129b657815110615a995701602001519150565b83634724a0fd60e01b5f5260045260245ffdfea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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
