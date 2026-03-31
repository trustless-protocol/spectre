// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractSP1ICS07Tendermint

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

// IICS07TendermintMsgsTrustThreshold is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsTrustThreshold struct {
	Numerator   uint8
	Denominator uint8
}

// ILightClientMsgsMsgVerifyMembership is an auto generated low-level Go binding around an user-defined struct.
type ILightClientMsgsMsgVerifyMembership struct {
	Height                IICS02ClientMsgsHeight
	KvPairs               []IMembershipMsgsKVPair
	MerkleProofs          []IMembershipMsgsMerkleProof
	AppHash               [32]byte
	TrustedConsensusState IICS07TendermintMsgsConsensusState
	MembershipType        uint8
}

// ILightClientMsgsMsgVerifyNonMembership is an auto generated low-level Go binding around an user-defined struct.
type ILightClientMsgsMsgVerifyNonMembership struct {
	Height                IICS02ClientMsgsHeight
	KvPairs               []IMembershipMsgsKVPair
	MerkleProofs          []IMembershipMsgsMerkleProof
	AppHash               [32]byte
	TrustedConsensusState IICS07TendermintMsgsConsensusState
	MembershipType        uint8
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
	Value [32]byte
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
	Value [32]byte
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

// ContractSP1ICS07TendermintMetaData contains all meta data concerning the ContractSP1ICS07Tendermint contract.
var ContractSP1ICS07TendermintMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membership_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviour_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"updateClient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_clientState\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ALLOWED_SP1_CLOCK_DRIFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MEMBERSHIP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMembership\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMisbehaviour\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PROOF_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPDATE_CLIENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUpdateClient\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VERIFIER\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"clientState\",\"inputs\":[],\"outputs\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConsensusStateHash\",\"inputs\":[{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"updateClientMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CannotHandleMisbehavior\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientStateMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidMembershipProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"KeyValuePairNotInCache\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProofHeightMismatch\",\"inputs\":[{\"name\":\"expectedRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"expectedRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TrustThresholdMismatch\",\"inputs\":[{\"name\":\"expectedNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expectedDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnbondingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"UnknownZkAlgorithm\",\"inputs\":[{\"name\":\"algorithm\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"VerificationKeyMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x6101006040523461055257614aaf8038038061001a81610575565b928339810160e082820312610552576100328261059a565b61003e6020840161059a565b61004a6040850161059a565b916100576060860161059a565b60808601519094906001600160401b0381116105525786019080601f83011215610552578151610089926020016105ae565b9461009b60c060a0830151920161059a565b9580518101906020820190602081840312610552576020810151906001600160401b038211610552570191829003601f1981019061012013610552576040519160e083016001600160401b0381118482101761053e5760405260208401516001600160401b038111610552576020908501019080601f83011215610552578151610127926020016105ae565b82526040811261055257604061013b610556565b916101478286016105ee565b8352610155606086016105ee565b602084015260208401928352603f19011261055257610172610556565b61017e608085016105fc565b815261018c60a085016105fc565b6020820152604083019081526101a460c08501610610565b90606084019182526101b860e08601610610565b92608085019384526101008601519586151587036105525760a0860196875261012001519460028610156105525760c08101958652518051906001600160401b03821161053e57600154600181811c91168015610534575b602082101461052057601f81116104bd575b50602090601f83116001146104505763ffffffff95949392915f9183610445575b50508160011b915f199060031b1c1916176001555b5160ff81511661ff00602060025493015160081b169161ffff191617176002555160018060401b0381511660035491602068010000000000000000600160801b0391015160401b169160018060801b031916171760035551169267ffffffff00000000600454925160201b169051151560401b92519160028310156104315769ff00000000000000000068ff00000000000000009360481b169469ff000000000000000000199160018060481b0319161716179116171760045560018060401b0360035460401c165f52600560205260405f205560018060a01b031660805260018060a01b031660a05260018060a01b031660c05260018060a01b031660e05260045463ffffffff8116610708810163ffffffff811161041d5763ffffffff809360201c16928391161161040857826001600160a01b0381166103ef575061039e610717565b505b604051614275908161079a823960805181818161039b015261285b015260a0518181816107990152612fe9015260c0518181816102450152610f0d015260e05181818161074901526126400152f35b806103fc61040292610621565b50610697565b506103a0565b6333fae18560e21b5f5260045260245260445ffd5b634e487b7160e01b5f52601160045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b015190505f80610243565b90601f1983169160015f52815f20925f5b8181106104a5575091600193918563ffffffff99989796941061048d575b505050811b01600155610258565b01515f1960f88460031b161c191690555f808061047f565b92936020600181928786015181550195019301610461565b60015f527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6601f840160051c81019160208510610516575b601f0160051c01905b81811061050b5750610222565b5f81556001016104fe565b90915081906104f5565b634e487b7160e01b5f52602260045260245ffd5b90607f1690610210565b634e487b7160e01b5f52604160045260245ffd5b5f80fd5b60408051919082016001600160401b0381118382101761053e57604052565b6040519190601f01601f191682016001600160401b0381118382101761053e57604052565b51906001600160a01b038216820361055257565b9192916001600160401b03821161053e576105d2601f8301601f1916602001610575565b938285528282011161055257815f926020928387015e84010152565b519060ff8216820361055257565b51906001600160401b038216820361055257565b519063ffffffff8216820361055257565b6001600160a01b0381165f9081525f516020614a8f5f395f51905f52602052604090205460ff16610692576001600160a01b03165f8181525f516020614a8f5f395f51905f5260205260408120805460ff191660011790553391905f516020614a0f5f395f51905f528180a4600190565b505f90565b6001600160a01b0381165f9081525f516020614a2f5f395f51905f52602052604090205460ff16610692576001600160a01b03165f8181525f516020614a2f5f395f51905f5260205260408120805460ff191660011790553391905f516020614a6f5f395f51905f52905f516020614a0f5f395f51905f529080a4600190565b5f80525f516020614a2f5f395f51905f526020525f516020614a4f5f395f51905f525460ff16610795575f8080525f516020614a2f5f395f51905f526020525f516020614a4f5f395f51905f52805460ff1916600117905533905f516020614a6f5f395f51905f525f516020614a0f5f395f51905f528280a4600190565b5f9056fe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a71461017457806302cf29521461016f578063038d5cb31461016a57806308c84e70146101655780630bece356146101605780630c6faf591461015b57806323842fb814610156578063248a9ca3146101515780632c3ee4741461014c5780632f2ff15d1461014757806336568abe146101425780635972185a1461013d57806387d4332f1461013857806389df51f1146101335780638a8e4c5d1461012e57806391d1485414610129578063a217fddf14610124578063ac9650d81461011f578063bd3ce6b01461011a578063d547741f14610115578063ddba6537146101105763ef913a4b1461010b575f80fd5b610ff5565b610e06565b610dd7565b610d29565b610970565b6108b4565b610866565b6107bd565b61076d565b61071d565b6106e3565b610687565b610651565b610602565b6105d8565b6105a9565b6104d1565b61044c565b61036f565b610299565b610219565b34610215576020600319360112610215576004357fffffffff00000000000000000000000000000000000000000000000000000000811680910361021557807f7965db0b00000000000000000000000000000000000000000000000000000000602092149081156101eb575b506040519015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150145f6101e0565b5f80fd5b34610215575f60031936011261021557602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b60206003198201126102155760043567ffffffffffffffff81116102155761012090600401809203126102155790565b346102155761035e61034e6102ad36610269565b6102bf60ff60045460401c16156110c7565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610362575b61031e60408201826110f6565b61032d606084939401836110f6565b90608084013592610100850135956103448761114a565b60a0860195612f6a565b6040519081529081906020820190565b0390f35b61036a612e75565b610311565b34610215575f60031936011261021557602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b9060206003198301126102155760043567ffffffffffffffff811161021557826023820112156102155780600401359267ffffffffffffffff84116102155760248483010111610215576024019190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b6003111561044757565b610410565b346102155760206104b361045f366103bf565b9061047260ff60045460401c16156110c7565b7fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b5f525f845260405f205f8052845260ff60405f205416156104c4576125e4565b604051906104c08161043d565b8152f35b6104cc612e75565b6125e4565b346102155761035e61034e6104e536610269565b6104f760ff60045460401c16156110c7565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac4754161561057d575b60808101359061055c60408201826110f6565b61056c60608495939501846110f6565b91610100850135956103448761114a565b610585612e75565b610549565b67ffffffffffffffff81160361021557565b35906105a78261058a565b565b346102155760206003193601126102155760206105d06004356105cb8161058a565b6129f7565b604051908152f35b346102155760206003193601126102155760206105d06004355f525f602052600160405f20015490565b34610215575f6003193601126102155760206040516107088152f35b6003196040910112610215576004359060243573ffffffffffffffffffffffffffffffffffffffff811681036102155790565b34610215576106856106623661061e565b9061068061067b825f525f602052600160405f20015490565b612efd565b6132b2565b005b34610215576106953661061e565b3373ffffffffffffffffffffffffffffffffffffffff8216036106bb576106859161338f565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b34610215575f6003193601126102155760206040517fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b8152f35b34610215575f60031936011261021557602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b34610215575f60031936011261021557602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b34610215576107cb366103bf565b50506107df60ff60045460401c16156110c7565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06020527f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac475460ff1615610859575b7fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b610861612e75565b610831565b3461021557602060ff6108a861087b3661061e565b905f525f845260405f209073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b54166040519015158152f35b34610215575f6003193601126102155760206040515f8152f35b90601f19601f602080948051918291828752018686015e5f8582860101520116010190565b602081016020825282518091526040820191602060408360051b8301019401925f915b83831061092557505050505090565b9091929394602080610961837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0866001960301875289516108ce565b97019301930191939290610916565b346102155760206003193601126102155760043567ffffffffffffffff811161021557366023820112156102155780600401359067ffffffffffffffff8211610215573660248360051b830101116102155761035e9160246109d29201612b10565b604051918291826108f3565b90600182811c92168015610a25575b60208310146109f857565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f16916109ed565b6001545f9291610a3e826109de565b8082529160018116908115610ab25750600114610a59575050565b60015f9081529293509091907fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b838310610a98575060209250010190565b600181602092949394548385870101520191019190610a87565b60209495507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091509291921683830152151560051b010190565b5f9291815491610afb836109de565b8083529260018116908115610b505750600114610b1757505050565b5f9081526020812093945091925b838310610b36575060209250010190565b600181602092949394548385870101520191019190610b25565b905060209495507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091509291921683830152151560051b010190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6040810190811067ffffffffffffffff821117610bd557604052565b610b8c565b6060810190811067ffffffffffffffff821117610bd557604052565b6080810190811067ffffffffffffffff821117610bd557604052565b60c0810190811067ffffffffffffffff821117610bd557604052565b60a0810190811067ffffffffffffffff821117610bd557604052565b90601f601f19910116810190811067ffffffffffffffff821117610bd557604052565b604051906105a760e083610c4a565b604051906105a761026083610c4a565b6002111561044757565b90610ca082610c8c565b52565b9363ffffffff929897969390610d098492610ce9610ccf6101009a966101208b526101208b01906108ce565b9c60208a019060ff60208092828151168552015116910152565b606088019067ffffffffffffffff60208092828151168552015116910152565b1660a08501521660c0830152151560e0820152610d2583610c8c565b0152565b34610215575f60031936011261021557604051610d5081610d4981610a2f565b0382610c4a565b60405190610d5d82610bb9565b60ff600254818116845260081c16602083015261035e604051610d7f81610bb9565b67ffffffffffffffff600354818116835260401c166020820152600454610da98163ffffffff1690565b602082901c63ffffffff1690610dca604084901c60ff169360481c60ff1690565b9360405197889788610ca3565b3461021557610685610de83661061e565b90610e0161067b825f525f602052600160405f20015490565b61338f565b3461021557610e83610e17366103bf565b610e2960ff60045460401c16156110c7565b5f80527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e060205260ff7f6e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47541615610fd4575b810190612c8c565b805190602081015191604082019283519360806060850192835194610ef483880196610ebf88516fffffffffffffffffffffffffffffffff1690565b906040519a8b9586957fa6fe8f5600000000000000000000000000000000000000000000000000000000875260048701612d75565b038173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa938415610fcf57610f62955f95610f9a575b50519051915192516fffffffffffffffffffffffffffffffff165b936134a8565b610685680100000000000000007fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff6004541617600455565b610f5c919550610fc19060803d608011610fc8575b610fb98183610c4a565b810190612d3e565b9490610f41565b503d610faf565b6122df565b610fdc612e75565b610e7b565b906020610ff29281815201906108ce565b90565b34610215575f6003193601126102155761035e60405160208082015261012060408201526110bb8161102a6101608201610a2f565b60ff600254818116606085015260081c16608083015261106360a08301602067ffffffffffffffff600354818116845260401c16910152565b60045463ffffffff811660e08401526110ad90602081901c63ffffffff1661010085015261109c610120850160ff8360401c1615159052565b60ff61014085019160481c16610c96565b03601f198101835282610c4a565b60405191829182610fe1565b156110ce57565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610215570180359067ffffffffffffffff821161021557602001918160051b3603831361021557565b6002111561021557565b67ffffffffffffffff8111610bd557601f01601f191660200190565b81601f820112156102155760208135910161118a82611154565b926111986040519485610c4a565b8284528282011161021557815f92602092838601378301015290565b60ff81160361021557565b9190826040910312610215576040516111d781610bb9565b602080829480356111e7816111b4565b8452013591610d25836111b4565b91908260409103126102155760405161120d81610bb9565b6020808294803561121d8161058a565b8452013591610d258361058a565b63ffffffff81160361021557565b35906105a78261122b565b8015150361021557565b35906105a782611244565b35906105a78261114a565b919091610120818403126102155761127a610c6d565b9281359167ffffffffffffffff8311610215576112c0826112a3610100946112fe968501611170565b87526112b281602085016111bf565b6020880152606083016111f5565b60408601526112d160a08201611239565b60608601526112e260c08201611239565b60808601526112f360e0820161124e565b60a086015201611259565b60c0830152565b6fffffffffffffffffffffffffffffffff81160361021557565b35906105a782611305565b91908260609103126102155760405161134281610bda565b6040808294803561135281611305565b8452602081013560208501520135910152565b8092910391606083126102155760405161137e81610bb9565b6040601f1982958435845201126102155760209060408051936113a085610bb9565b838101356113ad8161122b565b85520135828401520152565b67ffffffffffffffff8111610bd55760051b60200190565b6003111561021557565b919060c083820312610215576040516113f381610bf6565b809380356114008161058a565b825260208101356114108161122b565b60208301526114228360408301611365565b604083015260a08101359067ffffffffffffffff821161021557019180601f8401121561021557823592611455846113b9565b936114636040519586610c4a565b80855260208086019160051b830101918383116102155760208101915b83831061149257505050505060600152565b823567ffffffffffffffff8111610215578201906040601f19838803011261021557604051916114c183610bb9565b60208101356114cf816113d1565b8352604081013567ffffffffffffffff811161021557602091010190608082880312610215576040519261150284610bf6565b823567ffffffffffffffff8111610215578861151f918501611170565b8452602083013561152f81611305565b6020850152604083013561154281611244565b604085015260608301359367ffffffffffffffff85116102155761156b89602096879601611170565b606082015283820152815201920191611480565b9190604083820312610215576040519061159882610bb9565b8193803567ffffffffffffffff81116102155781016102c081840312610215576115c0610c7c565b906115cb84826111f5565b8252604081013567ffffffffffffffff811161021557846115ed918301611170565b60208301526115fe6060820161059c565b604083015261160f6080820161131f565b606083015261162060a0820161124e565b60808301526116328460c08301611365565b60a0830152611644610120820161124e565b60c083015261014081013560e0830152611661610160820161124e565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526116b0610220820161124e565b6101c08301526102408101356101e08301526116cf610260820161124e565b6102008301526102808101356102208301526102a08101359067ffffffffffffffff82116102155761170391859101611170565b610240820152835260208101359167ffffffffffffffff83116102155760209261172d92016113db565b910152565b9190608083820312610215576040519061174b82610bf6565b8193803567ffffffffffffffff81116102155760609261176c918301611170565b83526020810135602084015260408101356117868161058a565b60408401520135908160070b82036102155760600152565b91909160808184031261021557604051906117b882610bf6565b8193813567ffffffffffffffff811161021557820181601f820112156102155780356117e3816113b9565b916117f16040519384610c4a565b81835260208084019260051b820101918483116102155760208201905b838210611860575050505083526118276020830161124e565b602084015260408201359067ffffffffffffffff821161021557826118556060949261172d94869401611732565b60408601520161059c565b813567ffffffffffffffff81116102155760209161188388848094880101611732565b81520191019061180e565b919060a08382031261021557604051906118a782610bf6565b8193803567ffffffffffffffff811161021557826118c691830161157f565b8352602081013567ffffffffffffffff811161021557826118e891830161179e565b60208401526118fa82604083016111f5565b604084015260808101359167ffffffffffffffff83116102155760609261172d920161179e565b9080601f83011215610215576101006040519261193e8285610c4a565b8391810192831161021557905b8282106119585750505090565b813581526020918201910161194b565b9080601f830112156102155760405191611983604084610c4a565b82906040810192831161021557905b82821061199f5750505090565b8135815260209182019101611992565b6020818303126102155780359067ffffffffffffffff8211610215570161024081830312610215576119df610c6d565b91813567ffffffffffffffff811161021557816119fd918401611264565b8352611a0c816020840161132a565b6020840152608082013567ffffffffffffffff81116102155782611a388361020093611a79960161188e565b6040860152611a4960a0820161131f565b6060860152611a5b8360c08301611921565b6080860152611a6e836101c08301611968565b60a086015201611968565b60c082015290565b81601f8201121561021557602081519101611a9b82611154565b92611aa96040519485610c4a565b8284528282011161021557815f926020928386015e8301015290565b919082604091031261021557604051611add81610bb9565b60208082948051611aed816111b4565b8452015191610d25836111b4565b919082604091031261021557604051611b1381610bb9565b60208082948051611b238161058a565b8452015191610d258361058a565b51906105a78261122b565b51906105a782611244565b51906105a78261114a565b9190916101208184031261021557611b68610c6d565b9281519167ffffffffffffffff831161021557611bae82611b91610100946112fe968501611a81565b8752611ba08160208501611ac5565b602088015260608301611afb565b6040860152611bbf60a08201611b31565b6060860152611bd060c08201611b31565b6080860152611be160e08201611b3c565b60a086015201611b47565b51906105a782611305565b919082606091031261021557604051611c0f81610bda565b60408082948051611c1f81611305565b8452602081015160208501520151910152565b6020818303126102155780519067ffffffffffffffff82116102155701610180818303126102155760405191611c6783610c12565b815167ffffffffffffffff81116102155782611c8b8361014093611cdb9601611b52565b8552611c9a8360208301611bf7565b6020860152611cac8360808301611bf7565b6040860152611cbd60e08201611bec565b6060860152611cd0836101008301611afb565b608086015201611afb565b60a082015290565b90610ff29061010060c0611d02855161012085526101208501906108ce565b602080870151805160ff9081168784015291015116604085015294611d456040820151606086019067ffffffffffffffff60208092828151168552015116910152565b63ffffffff60608201511660a085015263ffffffff6080820151168285015260a0810151151560e08501520151910190610c96565b90606060c082019267ffffffffffffffff815116835263ffffffff6020820151166020840152611dcd6040820151604085019080518252602090810151805163ffffffff16828401520151604090910152565b01519160c060a0830152825180915260e0820191602060e08360051b8301019401925f915b838310611e0157505050505090565b9091929394602080611ec0837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff208660019603018752828a518051611e448161043d565b83520151906040848201526060611e6783516080604085015260c08401906108ce565b926fffffffffffffffffffffffffffffffff86820151168284015260408101511515608084015201519060a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0828503019101526108ce565b97019301930191939290611df2565b90606080611ee684516080855260808501906108ce565b936020810151602085015267ffffffffffffffff6040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b828210611f84575050505090606080611f72610ff294611f606020880151602087019015159052565b60408701518582036040870152611ecf565b94015167ffffffffffffffff16910152565b90919293602080611fbf837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff608a600196030186528851611ecf565b960192019201909291611f37565b610ff29160606121a161218f845160a08552602061215d8251604060a089015261201160e08901825167ffffffffffffffff60208092828151168552015116910152565b61024061202f848301516102c06101208c01526103a08b01906108ce565b604083015167ffffffffffffffff166101408b0152888301516fffffffffffffffffffffffffffffffff166101608b0152608083015115156101808b015260a083015180516101a08c0152602090810151805163ffffffff166101c08d015201516101e08b01529160c081015115156102008b015260e08101516102208b01526101008101511515828b01526101208101516102608b01526101408101516102808b01526101608101516102a08b01526101808101516102c08b01526101a08101516102e08b01526101c081015115156103008b01526101e08101516103208b015261020081015115156103408b01526102208101516103608b015201517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff20898303016103808a01526108ce565b9101517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff608683030160c0870152611d7a565b60208501518482036020860152611f10565b926121ca6040820151604085019067ffffffffffffffff60208092828151168552015116910152565b0151906080818403910152611f10565b905f905b600282106121eb57505050565b60208060019285518152019301910190916121de565b6020815261226e61222083516102406020850152610260840190611ce3565b61225860208501516040850190604080916fffffffffffffffffffffffffffffffff8151168452602081015160208501520151910152565b6040840151601f198483030160a0850152611fcd565b916fffffffffffffffffffffffffffffffff60608201511660c0830152608081015160e083015f905b600882106122c9575050509061022060c0836122bf60a0610ff29601516101e08601906121da565b01519101906121da565b6020806001928551815201930191019091612297565b6040513d5f823e3d90fd5b6105a7909291926060810193604080916fffffffffffffffffffffffffffffffff8151168452602081015160208501520151910152565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b80511561235b5760200190565b612321565b805182101561235b5760209160051b010190565b1561237b57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f696e76616c6964207369676e6174757265206c656e67746800000000000000006044820152fd5b6040908151916123e98184610c4a565b368337565b9060208282031261021557815167ffffffffffffffff811161021557610ff29201611a81565b929160608452606061012085019267ffffffffffffffff8151168287015263ffffffff602082015116608087015261246f60206040830151805160a08a0152015160c08801906020809163ffffffff81511684520151910152565b01519160c0610100860152825180915261014085019060206101408260051b8801019401915f905b8282106124c157505050506124b982604092866105a7950360208801526108ce565b5f9190940152565b90919294602080612503837ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec08c60019603018652828a518051611e448161043d565b970192019201909291612497565b908160209103126102155751610ff281611244565b939695949291905f855b6008821061259f575050509061254e612559926101008601906121da565b6101408401906121da565b5f61018083015b60028210612589575050610ff293945090610200916101c0820152816101e082015201906108ce565b6020806001928951815201970191019095612560565b6020806001928551815201930191019091612530565b156125bc57565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b6125f0918101906119af565b604051907f58d248c00000000000000000000000000000000000000000000000000000000082525f82806126278460048301612201565b038173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa918215610fcf575f926129d3575b5061267a82613160565b612683826131f4565b9161268d8361043d565b8261296c5761273261271a6020604060a085019485516126b78482015167ffffffffffffffff1690565b67ffffffffffffffff6126e76126da60035467ffffffffffffffff9060401c1690565b67ffffffffffffffff1690565b9116116128e7575b500151604051612706816110ad85820194856122ea565b5190209351015167ffffffffffffffff1690565b67ffffffffffffffff165f52600560205260405f2090565b555b6040810151908151906020612755818080865101519501519501515161234e565b510151916127de6060602061276c8288015161234e565b51015101519361277f6040865114612374565b60406127896123d9565b9560208101518752015160208601525f60808501519360c060a08701519601519760405194859283927fdf9a94a400000000000000000000000000000000000000000000000000000000845260048401612414565b038173b4b46bdaa835f8e4b4d8e208b6559cd2678510515af4908115610fcf57602095612841935f936128c3575b5060405197889687967fce5f643300000000000000000000000000000000000000000000000000000000885260048801612526565b03815f73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af18015610fcf57610ff2915f91612894575b506125b5565b6128b6915060203d6020116128bc575b6128ae8183610c4a565b810190612511565b5f61288e565b503d6128a4565b6128e09193503d805f833e6128d88183610c4a565b8101906123ee565b915f61280c565b6129669067ffffffffffffffff8151167fffffffffffffffffffffffffffffffff000000000000000000000000000000006fffffffffffffffff0000000000000000602060035494847fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000871617600355015160401b1692161717600355565b5f6126ef565b506129768261043d565b600182036129bb576129b6680100000000000000007fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff6004541617600455565b612734565b6129c48261043d565b60028203612734575050600290565b6129f09192503d805f833e6129e88183610c4a565b810190611c32565b905f612670565b67ffffffffffffffff165f52600560205260405f20548015612a165790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b91908203918211612a4b57565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b919081101561235b5760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18136030182121561021557019081359167ffffffffffffffff8311610215576020018236038113610215579190565b6020806105a793959486604051978895848701378401908282015f8152815193849201905e01015f815203601f198101845283610c4a565b612b195f611154565b612b266040519182610c4a565b5f8152601f19612b355f611154565b01366020830137612b45836113b9565b92612b536040519485610c4a565b808452601f19612b62826113b9565b015f5b818110612bbb5750505f5b818110612b7e575050505090565b80612b9f612b9985612b93600195878a612a78565b90612ad8565b30613464565b612ba98288612360565b52612bb48187612360565b5001612b70565b806060602080938901015201612b65565b91906060838203126102155760405190612be582610bda565b8193803567ffffffffffffffff81116102155781016040818403126102155760405190612c1182610bb9565b80359067ffffffffffffffff821161021557612c31856020938301611170565b83520135612c3e8161058a565b60208201528352602081013567ffffffffffffffff81116102155782612c6591830161188e565b602084015260408101359167ffffffffffffffff83116102155760409261172d920161188e565b6020818303126102155780359067ffffffffffffffff82116102155701610120818303126102155760405191612cc183610c2e565b813567ffffffffffffffff81116102155781612cde918401611264565b835260208201359167ffffffffffffffff831161021557612d2b82612d0b61010094612d36968501612bcc565b6020870152612d1d816040850161132a565b604087015260a0830161132a565b60608501520161131f565b608082015290565b9060808282031261021557612d6d906040805193612d5b85610bb9565b612d658382611afb565b855201611afb565b602082015290565b906105a794612e2f612dfe61010095612d9f612e5d959b9a989b6101208852610120880190611ce3565b86810360208801526040612ded83516060845267ffffffffffffffff6020612dd2835186606089015260a08801906108ce565b92015116608085015260208501518482036020860152611fcd565b920151906040818403910152611fcd565b986040850190604080916fffffffffffffffffffffffffffffffff8151168452602081015160208501520151910152565b80516fffffffffffffffffffffffffffffffff1660a0840152602081015160c08401526040015160e0830152565b01906fffffffffffffffffffffffffffffffff169052565b335f9081527f4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e0602052604090205460ff1615612ead57565b7fe2517d3f000000000000000000000000000000000000000000000000000000005f52336004527fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b60245260445ffd5b805f525f60205260ff612f313360405f209073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b541615612f3b5750565b7fe2517d3f000000000000000000000000000000000000000000000000000000005f523360045260245260445ffd5b95949290939196612f7a81610c8c565b806130e257505f939291612fcf9188151589816130d5575b612f9b91613566565b60405198899586957f13542b7a00000000000000000000000000000000000000000000000000000000875260048701613bb7565b03818373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af1918215610fcf57602061305b93610ff2955f916130b3575b5061304682825194019361303685613d02565b613040368861132a565b916140e5565b01516001815111613078575b5050369061132a565b6fffffffffffffffffffffffffffffffff633b9aca009151160490565b6130846130ac92613d02565b906130a661309185613d0c565b6fffffffffffffffffffffffffffffffff1690565b91614179565b5f80613052565b6130cf91503d805f833e6130c78183610c4a565b8101906135a4565b5f613023565b61ffff8111159150612f92565b806130ef61312692610c8c565b6130f881610c8c565b7f112d89cc000000000000000000000000000000000000000000000000000000005f5260ff16600452602490565b5ffd5b15613132575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6105a79061318681516fffffffffffffffffffffffffffffffff60608401511690613edc565b6131ec67ffffffffffffffff60206080818501516040516131d08482018093604080916fffffffffffffffffffffffffffffffff8151168452602081015160208501520151910152565b606081526131de8382610c4a565b5190209401510151166129f7565b808214613129565b61321061271a602060a0840151015167ffffffffffffffff1690565b548061321c5750505f90565b60408201908151604051613238816110ad6020820194856122ea565b5190201491821592613256575b50501561325157600190565b600290565b6fffffffffffffffffffffffffffffffff91925061309161328f60206132a79301516fffffffffffffffffffffffffffffffff90511690565b9351516fffffffffffffffffffffffffffffffff1690565b911610155f80613245565b805f525f60205260ff6132e68360405f209073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b541661338957805f525f60205261331e8260405f209073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082541617905573ffffffffffffffffffffffffffffffffffffffff339216907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260ff6133c38360405f209073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b54161561338957805f525f6020526133fc8260405f209073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00815416905573ffffffffffffffffffffffffffffffffffffffff339216907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b5f80610ff293602081519101845af43d156134a0573d9161348483611154565b926134926040519485610c4a565b83523d5f602085013e61404c565b60609161404c565b92602080916135216131ec956134ca6105a79967ffffffffffffffff97613edc565b6040516135008582018093604080916fffffffffffffffffffffffffffffffff8151168452602081015160208501520151910152565b6060815261350f608082610c4a565b5190206131ec86858a510151166129f7565b6040516135578382018093604080916fffffffffffffffffffffffffffffffff8151168452602081015160208501520151910152565b606081526131de608082610c4a565b1561356e5750565b7fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b6020818303126102155780519067ffffffffffffffff821161021557019060408282031261021557604051916135d983610bb9565b8051835260208101519067ffffffffffffffff8211610215570181601f820112156102155780519061360a826113b9565b926136186040519485610c4a565b82845260208085019360051b830101918183116102155760208101935b838510613649575050505050602082015290565b845167ffffffffffffffff81116102155782016040601f198286030112610215576040519061367782610bb9565b602081015167ffffffffffffffff81116102155760209082010185601f820112156102155780516136a7816113b9565b916136b56040519384610c4a565b81835260208084019260051b820101918883116102155760208201905b8382106136f957505050509160406020949285948352015183820152815201940193613635565b815167ffffffffffffffff81116102155760209161371c8c848094880101611a81565b8152019101906136d2565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18236030181121561021557016020813591019167ffffffffffffffff8211610215578160051b3603831361021557565b601f8260209493601f1993818652868601375f8582860101520116010190565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18236030181121561021557016020813591019167ffffffffffffffff821161021557813603831361021557565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8182360301811215610215570190565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa182360301811215610215570190565b90602083828152019260208260051b82010193835f925b8484106138755750505050505090565b9091929394956020806138e983601f1986600196030188526138978b8861381c565b9081356138a3816113d1565b6138ac8161043d565b81526138db6138d06138c08685018561379a565b606088860152606085019161377a565b92604081019061379a565b91604081850391015261377a565b9801940194019294939190613865565b610ff2916139b56139aa61391e613910858061379a565b60808652608086019161377a565b60208501356020850152608061399a61393a60408801886137ea565b8684036040880152803561394d816113d1565b6139568161043d565b84526020810135613966816113d1565b61396f8161043d565b60208501526040810135613982816113d1565b61398b8161043d565b6040850152606081019061379a565b919092816060820152019161377a565b926060810190613727565b91606081850391015261384e565b90602083828152019260208260051b82010193835f925b8484106139ea5750505050505090565b909192939495601f1982820301845286357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181121561021557830190613a38602082019280613727565b8091936020845252604082019060408160051b8401019380935f915b838310613a775750505050505060208060019298019401940192949391906139da565b9091929394957fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0838203018652613aae878361381c565b8035613ab98161114a565b613ac281610c8c565b8252613ae5613ad460208301836137ea565b6060602085015260608401906138f9565b9060408101357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6182360301811215610215576001936020938493613ba99301916040818303910152613b9b613b7b613b4e613b40858061379a565b60a0865260a086019161377a565b86850135613b5b81611244565b151587850152613b6e60408601866137ea565b84820360408601526138f9565b926060810135613b8a81611244565b1515606084015260808101906137ea565b9060808184039101526138f9565b980196019493019190613a54565b90928094929695936060830190835260606020840152526080810160808560051b8301019487915f907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc18a360301995b838310613c26575050505050610ff294955060408185039101526139c3565b90919293977fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8086820301835288358b811215610215578201906040810191613c6e8180613727565b809460408552526060830160608560051b85010194825f5b828110613cb05750505050506001926020928380809401359101529a019301930191939290613c07565b9091929396602080613cf5837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08b60019603018952613cef8c8861379a565b9061377a565b9901950193929101613c86565b35610ff28161058a565b35610ff281611305565b15613d1f575050565b7f20daedb6000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b15613d56575050565b7f12e74add000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b604051906105a782613d97816001610aec565b0383610c4a565b15613da65750565b613dfa906040519182917ff6b6676b00000000000000000000000000000000000000000000000000000000835260406004840152613de8604484016001610aec565b906003198483030160248501526108ce565b0390fd5b939291909315613e0e5750505050565b9160ff8092816084969581604051977fd382033a000000000000000000000000000000000000000000000000000000008952166004880152166024860152166044840152166064820152fd5b15613e63575050565b9063ffffffff80927f73b1bde3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b15613ea4575050565b9063ffffffff80927f1eb5c1eb000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b613f25613f026fffffffffffffffffffffffffffffffff6105a79416633b9aca00900490565b613f10814242821115613d16565b42610708613f1e8342612a3e565b1115613d4d565b805151613f336001546109de565b1480614028575b8151613f4591613d9e565b613f8d6020820151613f58815160ff1690565b6002549160ff831660ff811660ff8416149384613ffb575b613f879060209060081c60ff165b93015160ff1690565b93613dfe565b613fe7613fda6080613fa6606085015163ffffffff1690565b93613fcf60045495613fbb8763ffffffff1690565b63ffffffff811663ffffffff831614613e5a565b015163ffffffff1690565b9160201c63ffffffff1690565b63ffffffff811663ffffffff831614613e9b565b9350613f876020613f7e6140128286015160ff1690565b60ff600889901c81169116149692505050613f70565b50613f4581516020815191012061403d613d84565b60208151910120149050613f3a565b90614089575080511561406157602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b815115806140dc575b61409a575090565b73ffffffffffffffffffffffffffffffffffffffff907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15614092565b9161413e6020926131ec604051858101906141268287604080916fffffffffffffffffffffffffffffffff8151168452602081015160208501520151910152565b60608152614135608082610c4a565b519020916129f7565b015180820361414b575050565b7f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b929190925f5b8451811015614261576141928186612360565b519083604051602081019067ffffffffffffffff8616825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801915f5b81811061420c5750505050816001966020614202930151608083015203601f198101835282610c4a565b5190205d0161417f565b91939496509194969760208061424c837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff408b6001960301885289516108ce565b970194019101918a96949392989795986141d8565b505050905056fea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d4c48d1d2b0f4f7485fc28e3db22341d96a20aa29e6efa8149da9751603abd4e06e0c24a6e293ff9b755263dbaa15ba3796b0b8d3fe17cfb4ddf8143b268eac47bd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10bad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
}

// ContractSP1ICS07TendermintABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractSP1ICS07TendermintMetaData.ABI instead.
var ContractSP1ICS07TendermintABI = ContractSP1ICS07TendermintMetaData.ABI

// ContractSP1ICS07TendermintBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractSP1ICS07TendermintMetaData.Bin instead.
var ContractSP1ICS07TendermintBin = ContractSP1ICS07TendermintMetaData.Bin

// DeployContractSP1ICS07Tendermint deploys a new Ethereum contract, binding an instance of ContractSP1ICS07Tendermint to it.
func DeployContractSP1ICS07Tendermint(auth *bind.TransactOpts, backend bind.ContractBackend, verifier common.Address, membership_ common.Address, misbehaviour_ common.Address, updateClient_ common.Address, _clientState []byte, _consensusState [32]byte, roleManager common.Address) (common.Address, *types.Transaction, *ContractSP1ICS07Tendermint, error) {
	parsed, err := ContractSP1ICS07TendermintMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractSP1ICS07TendermintBin), backend, verifier, membership_, misbehaviour_, updateClient_, _clientState, _consensusState, roleManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractSP1ICS07Tendermint{ContractSP1ICS07TendermintCaller: ContractSP1ICS07TendermintCaller{contract: contract}, ContractSP1ICS07TendermintTransactor: ContractSP1ICS07TendermintTransactor{contract: contract}, ContractSP1ICS07TendermintFilterer: ContractSP1ICS07TendermintFilterer{contract: contract}}, nil
}

// ContractSP1ICS07Tendermint is an auto generated Go binding around an Ethereum contract.
type ContractSP1ICS07Tendermint struct {
	ContractSP1ICS07TendermintCaller     // Read-only binding to the contract
	ContractSP1ICS07TendermintTransactor // Write-only binding to the contract
	ContractSP1ICS07TendermintFilterer   // Log filterer for contract events
}

// ContractSP1ICS07TendermintCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractSP1ICS07TendermintCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSP1ICS07TendermintTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractSP1ICS07TendermintTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSP1ICS07TendermintFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractSP1ICS07TendermintFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSP1ICS07TendermintSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSP1ICS07TendermintSession struct {
	Contract     *ContractSP1ICS07Tendermint // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ContractSP1ICS07TendermintCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractSP1ICS07TendermintCallerSession struct {
	Contract *ContractSP1ICS07TendermintCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// ContractSP1ICS07TendermintTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractSP1ICS07TendermintTransactorSession struct {
	Contract     *ContractSP1ICS07TendermintTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// ContractSP1ICS07TendermintRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractSP1ICS07TendermintRaw struct {
	Contract *ContractSP1ICS07Tendermint // Generic contract binding to access the raw methods on
}

// ContractSP1ICS07TendermintCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractSP1ICS07TendermintCallerRaw struct {
	Contract *ContractSP1ICS07TendermintCaller // Generic read-only contract binding to access the raw methods on
}

// ContractSP1ICS07TendermintTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractSP1ICS07TendermintTransactorRaw struct {
	Contract *ContractSP1ICS07TendermintTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractSP1ICS07Tendermint creates a new instance of ContractSP1ICS07Tendermint, bound to a specific deployed contract.
func NewContractSP1ICS07Tendermint(address common.Address, backend bind.ContractBackend) (*ContractSP1ICS07Tendermint, error) {
	contract, err := bindContractSP1ICS07Tendermint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractSP1ICS07Tendermint{ContractSP1ICS07TendermintCaller: ContractSP1ICS07TendermintCaller{contract: contract}, ContractSP1ICS07TendermintTransactor: ContractSP1ICS07TendermintTransactor{contract: contract}, ContractSP1ICS07TendermintFilterer: ContractSP1ICS07TendermintFilterer{contract: contract}}, nil
}

// NewContractSP1ICS07TendermintCaller creates a new read-only instance of ContractSP1ICS07Tendermint, bound to a specific deployed contract.
func NewContractSP1ICS07TendermintCaller(address common.Address, caller bind.ContractCaller) (*ContractSP1ICS07TendermintCaller, error) {
	contract, err := bindContractSP1ICS07Tendermint(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractSP1ICS07TendermintCaller{contract: contract}, nil
}

// NewContractSP1ICS07TendermintTransactor creates a new write-only instance of ContractSP1ICS07Tendermint, bound to a specific deployed contract.
func NewContractSP1ICS07TendermintTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractSP1ICS07TendermintTransactor, error) {
	contract, err := bindContractSP1ICS07Tendermint(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractSP1ICS07TendermintTransactor{contract: contract}, nil
}

// NewContractSP1ICS07TendermintFilterer creates a new log filterer instance of ContractSP1ICS07Tendermint, bound to a specific deployed contract.
func NewContractSP1ICS07TendermintFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractSP1ICS07TendermintFilterer, error) {
	contract, err := bindContractSP1ICS07Tendermint(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractSP1ICS07TendermintFilterer{contract: contract}, nil
}

// bindContractSP1ICS07Tendermint binds a generic wrapper to an already deployed contract.
func bindContractSP1ICS07Tendermint(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractSP1ICS07TendermintMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractSP1ICS07Tendermint.Contract.ContractSP1ICS07TendermintCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.ContractSP1ICS07TendermintTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.ContractSP1ICS07TendermintTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractSP1ICS07Tendermint.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.contract.Transact(opts, method, params...)
}

// ALLOWEDSP1CLOCKDRIFT is a free data retrieval call binding the contract method 0x2c3ee474.
//
// Solidity: function ALLOWED_SP1_CLOCK_DRIFT() view returns(uint16)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) ALLOWEDSP1CLOCKDRIFT(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "ALLOWED_SP1_CLOCK_DRIFT")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// ALLOWEDSP1CLOCKDRIFT is a free data retrieval call binding the contract method 0x2c3ee474.
//
// Solidity: function ALLOWED_SP1_CLOCK_DRIFT() view returns(uint16)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) ALLOWEDSP1CLOCKDRIFT() (uint16, error) {
	return _ContractSP1ICS07Tendermint.Contract.ALLOWEDSP1CLOCKDRIFT(&_ContractSP1ICS07Tendermint.CallOpts)
}

// ALLOWEDSP1CLOCKDRIFT is a free data retrieval call binding the contract method 0x2c3ee474.
//
// Solidity: function ALLOWED_SP1_CLOCK_DRIFT() view returns(uint16)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) ALLOWEDSP1CLOCKDRIFT() (uint16, error) {
	return _ContractSP1ICS07Tendermint.Contract.ALLOWEDSP1CLOCKDRIFT(&_ContractSP1ICS07Tendermint.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ContractSP1ICS07Tendermint.Contract.DEFAULTADMINROLE(&_ContractSP1ICS07Tendermint.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ContractSP1ICS07Tendermint.Contract.DEFAULTADMINROLE(&_ContractSP1ICS07Tendermint.CallOpts)
}

// MEMBERSHIP is a free data retrieval call binding the contract method 0x89df51f1.
//
// Solidity: function MEMBERSHIP() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) MEMBERSHIP(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "MEMBERSHIP")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MEMBERSHIP is a free data retrieval call binding the contract method 0x89df51f1.
//
// Solidity: function MEMBERSHIP() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) MEMBERSHIP() (common.Address, error) {
	return _ContractSP1ICS07Tendermint.Contract.MEMBERSHIP(&_ContractSP1ICS07Tendermint.CallOpts)
}

// MEMBERSHIP is a free data retrieval call binding the contract method 0x89df51f1.
//
// Solidity: function MEMBERSHIP() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) MEMBERSHIP() (common.Address, error) {
	return _ContractSP1ICS07Tendermint.Contract.MEMBERSHIP(&_ContractSP1ICS07Tendermint.CallOpts)
}

// MISBEHAVIOUR is a free data retrieval call binding the contract method 0x02cf2952.
//
// Solidity: function MISBEHAVIOUR() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) MISBEHAVIOUR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "MISBEHAVIOUR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MISBEHAVIOUR is a free data retrieval call binding the contract method 0x02cf2952.
//
// Solidity: function MISBEHAVIOUR() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) MISBEHAVIOUR() (common.Address, error) {
	return _ContractSP1ICS07Tendermint.Contract.MISBEHAVIOUR(&_ContractSP1ICS07Tendermint.CallOpts)
}

// MISBEHAVIOUR is a free data retrieval call binding the contract method 0x02cf2952.
//
// Solidity: function MISBEHAVIOUR() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) MISBEHAVIOUR() (common.Address, error) {
	return _ContractSP1ICS07Tendermint.Contract.MISBEHAVIOUR(&_ContractSP1ICS07Tendermint.CallOpts)
}

// PROOFSUBMITTERROLE is a free data retrieval call binding the contract method 0x5972185a.
//
// Solidity: function PROOF_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) PROOFSUBMITTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "PROOF_SUBMITTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PROOFSUBMITTERROLE is a free data retrieval call binding the contract method 0x5972185a.
//
// Solidity: function PROOF_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) PROOFSUBMITTERROLE() ([32]byte, error) {
	return _ContractSP1ICS07Tendermint.Contract.PROOFSUBMITTERROLE(&_ContractSP1ICS07Tendermint.CallOpts)
}

// PROOFSUBMITTERROLE is a free data retrieval call binding the contract method 0x5972185a.
//
// Solidity: function PROOF_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) PROOFSUBMITTERROLE() ([32]byte, error) {
	return _ContractSP1ICS07Tendermint.Contract.PROOFSUBMITTERROLE(&_ContractSP1ICS07Tendermint.CallOpts)
}

// UPDATECLIENT is a free data retrieval call binding the contract method 0x87d4332f.
//
// Solidity: function UPDATE_CLIENT() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) UPDATECLIENT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "UPDATE_CLIENT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UPDATECLIENT is a free data retrieval call binding the contract method 0x87d4332f.
//
// Solidity: function UPDATE_CLIENT() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) UPDATECLIENT() (common.Address, error) {
	return _ContractSP1ICS07Tendermint.Contract.UPDATECLIENT(&_ContractSP1ICS07Tendermint.CallOpts)
}

// UPDATECLIENT is a free data retrieval call binding the contract method 0x87d4332f.
//
// Solidity: function UPDATE_CLIENT() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) UPDATECLIENT() (common.Address, error) {
	return _ContractSP1ICS07Tendermint.Contract.UPDATECLIENT(&_ContractSP1ICS07Tendermint.CallOpts)
}

// VERIFIER is a free data retrieval call binding the contract method 0x08c84e70.
//
// Solidity: function VERIFIER() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) VERIFIER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "VERIFIER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VERIFIER is a free data retrieval call binding the contract method 0x08c84e70.
//
// Solidity: function VERIFIER() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) VERIFIER() (common.Address, error) {
	return _ContractSP1ICS07Tendermint.Contract.VERIFIER(&_ContractSP1ICS07Tendermint.CallOpts)
}

// VERIFIER is a free data retrieval call binding the contract method 0x08c84e70.
//
// Solidity: function VERIFIER() view returns(address)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) VERIFIER() (common.Address, error) {
	return _ContractSP1ICS07Tendermint.Contract.VERIFIER(&_ContractSP1ICS07Tendermint.CallOpts)
}

// ClientState is a free data retrieval call binding the contract method 0xbd3ce6b0.
//
// Solidity: function clientState() view returns(string chainId, (uint8,uint8) trustLevel, (uint64,uint64) latestHeight, uint32 trustingPeriod, uint32 unbondingPeriod, bool isFrozen, uint8 zkAlgorithm)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) ClientState(opts *bind.CallOpts) (struct {
	ChainId         string
	TrustLevel      IICS07TendermintMsgsTrustThreshold
	LatestHeight    IICS02ClientMsgsHeight
	TrustingPeriod  uint32
	UnbondingPeriod uint32
	IsFrozen        bool
	ZkAlgorithm     uint8
}, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "clientState")

	outstruct := new(struct {
		ChainId         string
		TrustLevel      IICS07TendermintMsgsTrustThreshold
		LatestHeight    IICS02ClientMsgsHeight
		TrustingPeriod  uint32
		UnbondingPeriod uint32
		IsFrozen        bool
		ZkAlgorithm     uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainId = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.TrustLevel = *abi.ConvertType(out[1], new(IICS07TendermintMsgsTrustThreshold)).(*IICS07TendermintMsgsTrustThreshold)
	outstruct.LatestHeight = *abi.ConvertType(out[2], new(IICS02ClientMsgsHeight)).(*IICS02ClientMsgsHeight)
	outstruct.TrustingPeriod = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.UnbondingPeriod = *abi.ConvertType(out[4], new(uint32)).(*uint32)
	outstruct.IsFrozen = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.ZkAlgorithm = *abi.ConvertType(out[6], new(uint8)).(*uint8)

	return *outstruct, err

}

// ClientState is a free data retrieval call binding the contract method 0xbd3ce6b0.
//
// Solidity: function clientState() view returns(string chainId, (uint8,uint8) trustLevel, (uint64,uint64) latestHeight, uint32 trustingPeriod, uint32 unbondingPeriod, bool isFrozen, uint8 zkAlgorithm)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) ClientState() (struct {
	ChainId         string
	TrustLevel      IICS07TendermintMsgsTrustThreshold
	LatestHeight    IICS02ClientMsgsHeight
	TrustingPeriod  uint32
	UnbondingPeriod uint32
	IsFrozen        bool
	ZkAlgorithm     uint8
}, error) {
	return _ContractSP1ICS07Tendermint.Contract.ClientState(&_ContractSP1ICS07Tendermint.CallOpts)
}

// ClientState is a free data retrieval call binding the contract method 0xbd3ce6b0.
//
// Solidity: function clientState() view returns(string chainId, (uint8,uint8) trustLevel, (uint64,uint64) latestHeight, uint32 trustingPeriod, uint32 unbondingPeriod, bool isFrozen, uint8 zkAlgorithm)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) ClientState() (struct {
	ChainId         string
	TrustLevel      IICS07TendermintMsgsTrustThreshold
	LatestHeight    IICS02ClientMsgsHeight
	TrustingPeriod  uint32
	UnbondingPeriod uint32
	IsFrozen        bool
	ZkAlgorithm     uint8
}, error) {
	return _ContractSP1ICS07Tendermint.Contract.ClientState(&_ContractSP1ICS07Tendermint.CallOpts)
}

// GetClientState is a free data retrieval call binding the contract method 0xef913a4b.
//
// Solidity: function getClientState() view returns(bytes)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) GetClientState(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "getClientState")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetClientState is a free data retrieval call binding the contract method 0xef913a4b.
//
// Solidity: function getClientState() view returns(bytes)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) GetClientState() ([]byte, error) {
	return _ContractSP1ICS07Tendermint.Contract.GetClientState(&_ContractSP1ICS07Tendermint.CallOpts)
}

// GetClientState is a free data retrieval call binding the contract method 0xef913a4b.
//
// Solidity: function getClientState() view returns(bytes)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) GetClientState() ([]byte, error) {
	return _ContractSP1ICS07Tendermint.Contract.GetClientState(&_ContractSP1ICS07Tendermint.CallOpts)
}

// GetConsensusStateHash is a free data retrieval call binding the contract method 0x23842fb8.
//
// Solidity: function getConsensusStateHash(uint64 revisionHeight) view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) GetConsensusStateHash(opts *bind.CallOpts, revisionHeight uint64) ([32]byte, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "getConsensusStateHash", revisionHeight)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetConsensusStateHash is a free data retrieval call binding the contract method 0x23842fb8.
//
// Solidity: function getConsensusStateHash(uint64 revisionHeight) view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) GetConsensusStateHash(revisionHeight uint64) ([32]byte, error) {
	return _ContractSP1ICS07Tendermint.Contract.GetConsensusStateHash(&_ContractSP1ICS07Tendermint.CallOpts, revisionHeight)
}

// GetConsensusStateHash is a free data retrieval call binding the contract method 0x23842fb8.
//
// Solidity: function getConsensusStateHash(uint64 revisionHeight) view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) GetConsensusStateHash(revisionHeight uint64) ([32]byte, error) {
	return _ContractSP1ICS07Tendermint.Contract.GetConsensusStateHash(&_ContractSP1ICS07Tendermint.CallOpts, revisionHeight)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ContractSP1ICS07Tendermint.Contract.GetRoleAdmin(&_ContractSP1ICS07Tendermint.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ContractSP1ICS07Tendermint.Contract.GetRoleAdmin(&_ContractSP1ICS07Tendermint.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ContractSP1ICS07Tendermint.Contract.HasRole(&_ContractSP1ICS07Tendermint.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ContractSP1ICS07Tendermint.Contract.HasRole(&_ContractSP1ICS07Tendermint.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ContractSP1ICS07Tendermint.Contract.SupportsInterface(&_ContractSP1ICS07Tendermint.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ContractSP1ICS07Tendermint.Contract.SupportsInterface(&_ContractSP1ICS07Tendermint.CallOpts, interfaceId)
}

// UpgradeClient is a free data retrieval call binding the contract method 0x8a8e4c5d.
//
// Solidity: function upgradeClient(bytes ) view returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCaller) UpgradeClient(opts *bind.CallOpts, arg0 []byte) error {
	var out []interface{}
	err := _ContractSP1ICS07Tendermint.contract.Call(opts, &out, "upgradeClient", arg0)

	if err != nil {
		return err
	}

	return err

}

// UpgradeClient is a free data retrieval call binding the contract method 0x8a8e4c5d.
//
// Solidity: function upgradeClient(bytes ) view returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) UpgradeClient(arg0 []byte) error {
	return _ContractSP1ICS07Tendermint.Contract.UpgradeClient(&_ContractSP1ICS07Tendermint.CallOpts, arg0)
}

// UpgradeClient is a free data retrieval call binding the contract method 0x8a8e4c5d.
//
// Solidity: function upgradeClient(bytes ) view returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintCallerSession) UpgradeClient(arg0 []byte) error {
	return _ContractSP1ICS07Tendermint.Contract.UpgradeClient(&_ContractSP1ICS07Tendermint.CallOpts, arg0)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.GrantRole(&_ContractSP1ICS07Tendermint.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.GrantRole(&_ContractSP1ICS07Tendermint.TransactOpts, role, account)
}

// Misbehaviour is a paid mutator transaction binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes misbehaviourMsg) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactor) Misbehaviour(opts *bind.TransactOpts, misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.contract.Transact(opts, "misbehaviour", misbehaviourMsg)
}

// Misbehaviour is a paid mutator transaction binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes misbehaviourMsg) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) Misbehaviour(misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.Misbehaviour(&_ContractSP1ICS07Tendermint.TransactOpts, misbehaviourMsg)
}

// Misbehaviour is a paid mutator transaction binding the contract method 0xddba6537.
//
// Solidity: function misbehaviour(bytes misbehaviourMsg) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactorSession) Misbehaviour(misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.Misbehaviour(&_ContractSP1ICS07Tendermint.TransactOpts, misbehaviourMsg)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.Multicall(&_ContractSP1ICS07Tendermint.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.Multicall(&_ContractSP1ICS07Tendermint.TransactOpts, data)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.RenounceRole(&_ContractSP1ICS07Tendermint.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.RenounceRole(&_ContractSP1ICS07Tendermint.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.RevokeRole(&_ContractSP1ICS07Tendermint.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.RevokeRole(&_ContractSP1ICS07Tendermint.TransactOpts, role, account)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x0bece356.
//
// Solidity: function updateClient(bytes updateClientMsg) returns(uint8)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactor) UpdateClient(opts *bind.TransactOpts, updateClientMsg []byte) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.contract.Transact(opts, "updateClient", updateClientMsg)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x0bece356.
//
// Solidity: function updateClient(bytes updateClientMsg) returns(uint8)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) UpdateClient(updateClientMsg []byte) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.UpdateClient(&_ContractSP1ICS07Tendermint.TransactOpts, updateClientMsg)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x0bece356.
//
// Solidity: function updateClient(bytes updateClientMsg) returns(uint8)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactorSession) UpdateClient(updateClientMsg []byte) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.UpdateClient(&_ContractSP1ICS07Tendermint.TransactOpts, updateClientMsg)
}

// VerifyMembership is a paid mutator transaction binding the contract method 0x0c6faf59.
//
// Solidity: function verifyMembership(((uint64,uint64),(bytes[],bytes32)[],((uint8,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8) msg_) returns(uint256)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactor) VerifyMembership(opts *bind.TransactOpts, msg_ ILightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.contract.Transact(opts, "verifyMembership", msg_)
}

// VerifyMembership is a paid mutator transaction binding the contract method 0x0c6faf59.
//
// Solidity: function verifyMembership(((uint64,uint64),(bytes[],bytes32)[],((uint8,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8) msg_) returns(uint256)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) VerifyMembership(msg_ ILightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.VerifyMembership(&_ContractSP1ICS07Tendermint.TransactOpts, msg_)
}

// VerifyMembership is a paid mutator transaction binding the contract method 0x0c6faf59.
//
// Solidity: function verifyMembership(((uint64,uint64),(bytes[],bytes32)[],((uint8,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8) msg_) returns(uint256)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactorSession) VerifyMembership(msg_ ILightClientMsgsMsgVerifyMembership) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.VerifyMembership(&_ContractSP1ICS07Tendermint.TransactOpts, msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0x038d5cb3.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes32)[],((uint8,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8) msg_) returns(uint256)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactor) VerifyNonMembership(opts *bind.TransactOpts, msg_ ILightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.contract.Transact(opts, "verifyNonMembership", msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0x038d5cb3.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes32)[],((uint8,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8) msg_) returns(uint256)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintSession) VerifyNonMembership(msg_ ILightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.VerifyNonMembership(&_ContractSP1ICS07Tendermint.TransactOpts, msg_)
}

// VerifyNonMembership is a paid mutator transaction binding the contract method 0x038d5cb3.
//
// Solidity: function verifyNonMembership(((uint64,uint64),(bytes[],bytes32)[],((uint8,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),(bytes,bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[]),bool,(bytes,bytes32,(uint8,uint8,uint8,bytes),(uint8,bytes,bytes)[])))[])[],bytes32,(uint128,bytes32,bytes32),uint8) msg_) returns(uint256)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintTransactorSession) VerifyNonMembership(msg_ ILightClientMsgsMsgVerifyNonMembership) (*types.Transaction, error) {
	return _ContractSP1ICS07Tendermint.Contract.VerifyNonMembership(&_ContractSP1ICS07Tendermint.TransactOpts, msg_)
}

// ContractSP1ICS07TendermintRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ContractSP1ICS07Tendermint contract.
type ContractSP1ICS07TendermintRoleAdminChangedIterator struct {
	Event *ContractSP1ICS07TendermintRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ContractSP1ICS07TendermintRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSP1ICS07TendermintRoleAdminChanged)
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
		it.Event = new(ContractSP1ICS07TendermintRoleAdminChanged)
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
func (it *ContractSP1ICS07TendermintRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSP1ICS07TendermintRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSP1ICS07TendermintRoleAdminChanged represents a RoleAdminChanged event raised by the ContractSP1ICS07Tendermint contract.
type ContractSP1ICS07TendermintRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ContractSP1ICS07TendermintRoleAdminChangedIterator, error) {

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

	logs, sub, err := _ContractSP1ICS07Tendermint.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ContractSP1ICS07TendermintRoleAdminChangedIterator{contract: _ContractSP1ICS07Tendermint.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ContractSP1ICS07TendermintRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _ContractSP1ICS07Tendermint.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSP1ICS07TendermintRoleAdminChanged)
				if err := _ContractSP1ICS07Tendermint.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintFilterer) ParseRoleAdminChanged(log types.Log) (*ContractSP1ICS07TendermintRoleAdminChanged, error) {
	event := new(ContractSP1ICS07TendermintRoleAdminChanged)
	if err := _ContractSP1ICS07Tendermint.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractSP1ICS07TendermintRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ContractSP1ICS07Tendermint contract.
type ContractSP1ICS07TendermintRoleGrantedIterator struct {
	Event *ContractSP1ICS07TendermintRoleGranted // Event containing the contract specifics and raw log

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
func (it *ContractSP1ICS07TendermintRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSP1ICS07TendermintRoleGranted)
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
		it.Event = new(ContractSP1ICS07TendermintRoleGranted)
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
func (it *ContractSP1ICS07TendermintRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSP1ICS07TendermintRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSP1ICS07TendermintRoleGranted represents a RoleGranted event raised by the ContractSP1ICS07Tendermint contract.
type ContractSP1ICS07TendermintRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ContractSP1ICS07TendermintRoleGrantedIterator, error) {

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

	logs, sub, err := _ContractSP1ICS07Tendermint.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ContractSP1ICS07TendermintRoleGrantedIterator{contract: _ContractSP1ICS07Tendermint.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ContractSP1ICS07TendermintRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ContractSP1ICS07Tendermint.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSP1ICS07TendermintRoleGranted)
				if err := _ContractSP1ICS07Tendermint.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintFilterer) ParseRoleGranted(log types.Log) (*ContractSP1ICS07TendermintRoleGranted, error) {
	event := new(ContractSP1ICS07TendermintRoleGranted)
	if err := _ContractSP1ICS07Tendermint.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractSP1ICS07TendermintRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ContractSP1ICS07Tendermint contract.
type ContractSP1ICS07TendermintRoleRevokedIterator struct {
	Event *ContractSP1ICS07TendermintRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ContractSP1ICS07TendermintRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSP1ICS07TendermintRoleRevoked)
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
		it.Event = new(ContractSP1ICS07TendermintRoleRevoked)
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
func (it *ContractSP1ICS07TendermintRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSP1ICS07TendermintRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSP1ICS07TendermintRoleRevoked represents a RoleRevoked event raised by the ContractSP1ICS07Tendermint contract.
type ContractSP1ICS07TendermintRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ContractSP1ICS07TendermintRoleRevokedIterator, error) {

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

	logs, sub, err := _ContractSP1ICS07Tendermint.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ContractSP1ICS07TendermintRoleRevokedIterator{contract: _ContractSP1ICS07Tendermint.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ContractSP1ICS07TendermintRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ContractSP1ICS07Tendermint.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSP1ICS07TendermintRoleRevoked)
				if err := _ContractSP1ICS07Tendermint.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_ContractSP1ICS07Tendermint *ContractSP1ICS07TendermintFilterer) ParseRoleRevoked(log types.Log) (*ContractSP1ICS07TendermintRoleRevoked, error) {
	event := new(ContractSP1ICS07TendermintRoleRevoked)
	if err := _ContractSP1ICS07Tendermint.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
