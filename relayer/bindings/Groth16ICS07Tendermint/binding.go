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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membership_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviour_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"updateClient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_clientState\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ALLOWED_CLOCK_DRIFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MEMBERSHIP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMembership\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMisbehaviour\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PROOF_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPDATE_CLIENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUpdateClient\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VERIFIER\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"clientState\",\"inputs\":[],\"outputs\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConsensusStateHash\",\"inputs\":[{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasCachedValidatorSet\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"updateClientMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CannotHandleMisbehavior\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientStateMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidMembershipProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"KeyValuePairNotInCache\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProofHeightMismatch\",\"inputs\":[{\"name\":\"expectedRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"expectedRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustThresholdMismatch\",\"inputs\":[{\"name\":\"expectedNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expectedDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnbondingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"UnknownZkAlgorithm\",\"inputs\":[{\"name\":\"algorithm\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"VerificationKeyMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x61014080604052346105ec576155f4803803809161001d828561060b565b8339810160e0828203126105ec576100348261062e565b6100406020840161062e565b61004c6040850161062e565b916100596060860161062e565b60808601519094906001600160401b0381116105ec5786019080601f830112156105ec57815161008b92602001610642565b9461009d60c060a0830151920161062e565b957fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b61010052805181019060208201906020818403126105ec576020810151906001600160401b0382116105ec570191829003601f19810190610120136105ec576040519160e083016001600160401b038111848210176105d85760405260208401516001600160401b0381116105ec576020908501019080601f830112156105ec57815161014e92602001610642565b8252604081126105ec576040805191610166836105f0565b610171828601610687565b835261017f60608601610687565b602084015260208401928352603f1901126105ec576040516101a0816105f0565b6101ac60808501610695565b81526101ba60a08501610695565b6020820152604083019081526101d260c085016106a9565b90606084019182526101e660e086016106a9565b92608085019384526101008601519586151587036105ec5760a0860196875261012001519460028610156105ec5760c08101958652518051906001600160401b0382116105d8576102386001546106ba565b601f8111610588575b50602090601f831160011461051b5763ffffffff95949392915f9183610510575b50508160011b915f199060031b1c1916176001555b5160ff81511661ff00602060025493015160081b169161ffff191617176002555160018060401b0381511660035491602068010000000000000000600160801b0391015160401b169160018060801b031916171760035551169267ffffffff00000000600454925160201b169051151560401b92519160028310156104fc5769ff00000000000000000068ff00000000000000009360481b169469ff000000000000000000199160018060481b0319161716179116171760045560405161034881610341816106f2565b038261060b565b602081519101206101205260405161036381610341816106f2565b6001600160401b0390602090610378906107a3565b0151600354916001600160401b03831691168181036104e7575050604090811c6001600160401b03165f908152600560205220556001600160a01b0390811660805290811660a05290811660c0521660e05260045463ffffffff8082169061070882019081116104d35763ffffffff809360201c1692839116116104be57826001600160a01b0381166104a15750610412610100516109fa565b505b604051614acb9081610ac98239608051818181610f700152612c2b015260a051818181610d52015261397e015260c0518181816104d60152610fe2015260e051818181610d9501528181612787015261307001526101005181818161020b01528181610a5101528181610b6c01528181610cbe01528181610dd00152610eef015261012051816144a10152f35b806104ae6104b892610984565b5061010051610a53565b50610414565b6333fae18560e21b5f5260045260245260445ffd5b634e487b7160e01b5f52601160045260245ffd5b6355bace6f60e01b5f5260045260245260445ffd5b634e487b7160e01b5f52602160045260245ffd5b015190505f80610262565b90601f1983169160015f52815f20925f5b818110610570575091600193918563ffffffff999897969410610558575b505050811b01600155610277565b01515f1960f88460031b161c191690555f808061054a565b9293602060018192878601518155019501930161052c565b60015f525f5160206155b45f395f51905f52601f840160051c810191602085106105ce575b601f0160051c01905b8181106105c35750610241565b5f81556001016105b6565b90915081906105ad565b634e487b7160e01b5f52604160045260245ffd5b5f80fd5b604081019081106001600160401b038211176105d857604052565b601f909101601f19168101906001600160401b038211908210176105d857604052565b51906001600160a01b03821682036105ec57565b9192916001600160401b0382116105d8576040519161066b601f8201601f19166020018461060b565b8294818452818301116105ec578281602093845f96015e010152565b519060ff821682036105ec57565b51906001600160401b03821682036105ec57565b519063ffffffff821682036105ec57565b90600182811c921680156106e8575b60208310146106d457565b634e487b7160e01b5f52602260045260245ffd5b91607f16916106c9565b6001545f9291610701826106ba565b8082529160018116908115610762575060011461071c575050565b60015f9081529293509091905f5160206155b45f395f51905f525b838310610748575060209250010190565b600181602092949394548385870101520191019190610737565b9050602093945060ff929192191683830152151560051b010190565b90815181101561078f570160200190565b634e487b7160e01b5f52603260045260245ffd5b6040516107af816105f0565b606081525f60208201525080518015908115610978575b50610969575f19908051805b610925575b505f1982146109175760018201908183116104d357600360fc1b6001600160f81b0319610804848461077e565b51161480610902575b6108f2575f5b81518310156108a257610826838361077e565b5160f81c603081108015610898575b61088657600a82026001600160401b03908116602f1990920160ff169190910181169116811061086a57600190920191610813565b509150506040519061087b826105f0565b81525f602082015290565b50509150506040519061087b826105f0565b5060398111610835565b91509160018111908115916108e6575b506108d757604051916108c4836105f0565b82526001600160401b0316602082015290565b63384c687760e21b5f5260045ffd5b602b915010155f6108b2565b9150506040519061087b826105f0565b5080518381039081116104d35760021061080d565b604051915061087b826105f0565b5f1981018181116104d357602d60f81b6001600160f81b0319610948838661077e565b51161461095f575080156104d3575f1901806107d2565b92505f90506107d7565b6329120bff60e21b5f5260045ffd5b6040915010155f6107c6565b6001600160a01b0381165f9081525f5160206155d45f395f51905f52602052604090205460ff166109f5576001600160a01b03165f8181525f5160206155d45f395f51905f5260205260408120805460ff191660011790553391905f5160206155945f395f51905f528180a4600190565b505f90565b805f525f60205260405f205f805260205260ff60405f205416155f146109f557805f525f60205260405f205f805260205260405f20600160ff198254161790555f33915f5160206155945f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16610ac2575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291905f5160206155945f395f51905f529080a4600190565b50505f9056fe6080806040526004361015610012575f80fd5b5f3560e01c90816301ffc9a7146110065750806302cf295214610fc3578063083a1fee14610f9457806308c84e7014610f515780630bece35614610eca57806323842fb814610e9b578063248a9ca314610e715780632f2ff15d14610e4257806336568abe14610df35780635972185a14610db957806387d4332f14610d7657806389df51f114610d335780638a8e4c5d14610c9f57806391d1485414610c63578063974a74c414610b27578063a217fddf14610b0d578063a6f031bb14610a0c578063ac9650d814610847578063aef1f78a1461082b578063bd3ce6b014610723578063d547741f146106ed578063ddba6537146101ec5763ef913a4b14610119575f80fd5b346101e8575f3660031901126101e8576101e460405160208082015261012060408201526101d08161014e61016082016111a7565b60ff600254818116606085015260081c1660808301526001600160401b0360035481811660a085015260401c1660c083015260ff60045463ffffffff811660e085015263ffffffff8160201c16610100850152818160401c16151561012085015260481c166101bc816112ee565b61014083015203601f1981018352826112cd565b60405191829160208352602083019061114b565b0390f35b5f80fd5b346101e8576101fa366110a4565b6004549060ff8260401c166106c5577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f205416156106b6575b508201916020818403126101e8578035906001600160401b0382116101e85701610120818403126101e8576040519060a082018281106001600160401b038211176106a25760405280356001600160401b0381116101e857846102ad9183016113c9565b825260208101356001600160401b0381116101e8578101936060858203126101e857604051946102dc86611261565b80356001600160401b0381116101e85781016040818403126101e8576040519061030582611246565b80356001600160401b0381116101e857816103278660209361032f9501611349565b845201611111565b6020820152865260208101356001600160401b0381116101e857826103559183016116a6565b602087015260408101356001600160401b0381116101e857610379918391016116a6565b604086015260208301948552806040830190610394916114ae565b60408401908152906103a99060a084016114ae565b926060810192848452610100016103bf9061149a565b9560808201948786528251915197845191604051998a947fa6fe8f56000000000000000000000000000000000000000000000000000000008652600486016101209052610124860161041091611d8a565b906003198683030160248701528051606083528051606084016040905260a0840161043a9161114b565b90602001516001600160401b03166080840152602082015190838103602085015261046491611ef7565b906040015191808203906040015261047b91611ef7565b83516001600160801b0316604486015260208401516064860152604090930151608485015280516001600160801b031660a4850152602081015160c48501526040015160e48401526001600160801b031661010483015203867f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031691815a93608094fa958615610697575f96610603575b506020806105f1956105a768010000000000000000999661055161059f976001600160801b036001600160401b039851915193519551169061446b565b60405161057e8582018093604080916001600160801b038151168452602081015160208501520151910152565b6060815261058d6080826112cd565b51902061059f86858a51015116613125565b808214613421565b6040516105d48382018093604080916001600160801b038151168452602081015160208501520151910152565b606081526105e36080826112cd565b519020940151015116613125565b68ff0000000000000000191617600455005b9195509160803d608011610690575b61061c81846112cd565b8201916080818403126101e85760206105f1956105a76001600160401b0394610551680100000000000000009b6001600160801b03869761067a61059f9b604080519361066885611246565b6106728382611b40565b855201611b40565b888201529d505097505096945050955050610514565b503d610612565b6040513d5f823e3d90fd5b634e487b7160e01b5f52604160045260245ffd5b6106bf906131fa565b83610249565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b346101e8576107216106fe36611125565b9061071c610717825f525f602052600160405f20015490565b6131fa565b61383d565b005b346101e8575f3660031901126101e8576107b960405161074d81610746816111a7565b03826112cd565b60405161075981611246565b60ff600254818116835260081c16602082015260405161077881611246565b6001600160401b03600354818116835260401c16602082015260ff6004546107f2828260481c16936107d36040519889986101208a526101208a019061114b565b96602089019060ff60208092828151168552015116910152565b60608701906001600160401b0360208092828151168552015116910152565b63ffffffff811660a086015263ffffffff8160201c1660c086015260401c16151560e0840152610821816112ee565b6101008301520390f35b346101e8575f3660031901126101e85760206040516107088152f35b346101e85760203660031901126101e8576004356001600160401b0381116101e857366023820112156101e85780600401356001600160401b0381116101e85760248201913660248360051b830101116101e8576020926040516108ab85826112cd565b5f815284810191601f1986013684376108c385611539565b936108d160405195866112cd565b858552601f196108e087611539565b01875f5b8281106109fd575050505f5b868110156109a05760019061097c5f808b8861094961091760248860051b8b01018b61316b565b90938d6040519483869484860198893784019083820190898252519283915e010185815203601f1981018352826112cd565b5190305af43d15610998573d9061095f826112f8565b9161096d60405193846112cd565b82523d5f8d84013e5b306149c0565b61098682896131d2565b5261099181886131d2565b50016108f0565b606090610976565b604080518981528751818b018190525f92600582901b83018101918a8d01918d9085015b8287106109d15785850386f35b9091929382806109ed600193603f198a8203018652885161114b565b96019201960195929190926109c4565b606088820183015281016108e4565b346101e85760203660031901126101e8576004356001600160401b0381116101e857806004019061014060031982360301126101e85760ff60045460401c166106c5577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610afe575b50610a9d604482018361319d565b91610aab606482018561319d565b9390916101048101359060028210156101e857602096610af696610ad361012484018361319d565b96909560405198610ae48c8b6112cd565b5f8a52608460a48701960135946138c0565b604051908152f35b610b07906131fa565b82610a8f565b346101e8575f3660031901126101e85760206040515f8152f35b346101e85760203660031901126101e8576004356001600160401b0381116101e857806004019061016060031982360301126101e85760ff60045460401c166106c5577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610c54575b506101448101610bba818461316b565b905015610c2c57610bce604483018461319d565b9091610bdd606485018661319d565b9590946101048101359160028310156101e857602097610af697610c1596610c1c610c0c61012487018661319d565b9990988661316b565b3691611313565b98608460a48701960135946138c0565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b610c5d906131fa565b82610baa565b346101e857610c7136611125565b905f525f6020526001600160a01b0360405f2091165f52602052602060ff60405f2054166040519015158152f35b346101e857610cad366110a4565b505060ff60045460401c166106c5577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610d24575b7fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b610d2d906131fa565b80610cfc565b346101e8575f3660031901126101e85760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101e8575f3660031901126101e85760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101e8575f3660031901126101e85760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346101e857610e0136611125565b336001600160a01b03821603610e1a576107219161383d565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b346101e857610721610e5336611125565b90610e6c610717825f525f602052600160405f20015490565b6137b0565b346101e85760203660031901126101e8576020610af66004355f525f602052600160405f20015490565b346101e85760203660031901126101e8576004356001600160401b03811681036101e857610af6602091613125565b346101e857610ed8366110a4565b9060ff60045460401c166106c557602091610f31917f0000000000000000000000000000000000000000000000000000000000000000805f525f855260405f205f8052855260ff60405f20541615610f42575b506122cb565b60405190610f3e816110f3565b8152f35b610f4b906131fa565b84610f2b565b346101e8575f3660031901126101e85760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101e85760203660031901126101e8576004355f526006602052602060ff60405f2054166040519015158152f35b346101e8575f3660031901126101e85760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101e85760203660031901126101e857600435907fffffffff0000000000000000000000000000000000000000000000000000000082168092036101e857817f7965db0b000000000000000000000000000000000000000000000000000000006020931490811561107a575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483611073565b9060206003198301126101e8576004356001600160401b0381116101e857826023820112156101e8578060040135926001600160401b0384116101e857602484830101116101e8576024019190565b600311156110fd57565b634e487b7160e01b5f52602160045260245ffd5b35906001600160401b03821682036101e857565b60409060031901126101e857600435906024356001600160a01b03811681036101e85790565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b90600182811c9216801561119d575b602083101461118957565b634e487b7160e01b5f52602260045260245ffd5b91607f169161117e565b6001545f92916111b68261116f565b808252916001811690811561122a57506001146111d1575050565b60015f9081529293509091907fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b838310611210575060209250010190565b6001816020929493945483858701015201910191906111ff565b9050602093945060ff929192191683830152151560051b010190565b604081019081106001600160401b038211176106a257604052565b606081019081106001600160401b038211176106a257604052565b60e081019081106001600160401b038211176106a257604052565b608081019081106001600160401b038211176106a257604052565b60c081019081106001600160401b038211176106a257604052565b90601f801991011681019081106001600160401b038211176106a257604052565b600211156110fd57565b6001600160401b0381116106a257601f01601f191660200190565b92919261131f826112f8565b9161132d60405193846112cd565b8294818452818301116101e8578281602093845f960137010152565b9080601f830112156101e85781602061136493359101611313565b90565b359060ff821682036101e857565b91908260409103126101e85760405161138d81611246565b60206113a681839561139e81611111565b855201611111565b910152565b359063ffffffff821682036101e857565b359081151582036101e857565b808203929161012084126101e857604051916113e48361127c565b82948135906001600160401b0382116101e857611405846040938501611349565b8552601f1901126101e8576114506101009260405161142381611246565b61142f60208501611367565b815261143d60408501611367565b6020820152602086015260608301611375565b604084015261146160a082016113ab565b606084015261147260c082016113ab565b608084015261148360e082016113bc565b60a084015201359060028210156101e85760c00152565b35906001600160801b03821682036101e857565b91908260609103126101e8576040516114c681611261565b60408082946114d48161149a565b8452602081013560208501520135910152565b8092910391606083126101e85760405161150081611246565b6040819483358352601f1901126101e857602090604080519361152285611246565b61152d8482016113ab565b85520135828401520152565b6001600160401b0381116106a25760051b60200190565b91906080838203126101e8576040519061156982611297565b819380356001600160401b0381116101e857606092611589918301611349565b8352602081013560208401526115a160408201611111565b60408401520135908160070b82036101e85760600152565b9190916080818403126101e857604051906115d382611297565b819381356001600160401b0381116101e857820181601f820112156101e85780356115fd81611539565b9161160b60405193846112cd565b81835260208084019260051b820101918483116101e85760208201905b83821061167957505050508352611641602083016113bc565b60208401526040820135906001600160401b0382116101e8578261166e606094926113a694869401611550565b604086015201611111565b81356001600160401b0381116101e85760209161169b88848094880101611550565b815201910190611628565b919060a0838203126101e857604051906116bf82611297565b819380356001600160401b0381116101e85781016040818403126101e857604051906116ea82611246565b80356001600160401b0381116101e85781016102c0818603126101e8576040519061026082018281106001600160401b038211176106a25760405261172f8682611375565b825260408101356001600160401b0381116101e85786611750918301611349565b602083015261176160608201611111565b60408301526117726080820161149a565b606083015261178360a082016113bc565b60808301526117958660c083016114e7565b60a08301526117a761012082016113bc565b60c083015261014081013560e08301526117c461016082016113bc565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a083015261181361022082016113bc565b6101c08301526102408101356101e083015261183261026082016113bc565b6102008301526102808101356102208301526102a0810135906001600160401b0382116101e85761186591879101611349565b61024082015282526020810135906001600160401b0382116101e8570160c0818503126101e8576040519061189982611297565b6118a281611111565b82526118b0602082016113ab565b60208301526118c285604083016114e7565b604083015260a0810135906001600160401b0382116101e8570184601f820112156101e8578035906118f382611539565b9161190160405193846112cd565b80835260208084019160051b830101918783116101e85760208101915b83831061198c575050505060608201526020820152835260208101356001600160401b0381116101e857826119549183016115b9565b60208401526119668260408301611375565b60408401526080810135916001600160401b0383116101e8576060926113a692016115b9565b82356001600160401b0381116101e8578201906040828b03601f1901126101e857604051916119ba83611246565b602081013560048110156101e857835260408101356001600160401b0381116101e8576020910101906080828c03126101e857604051926119fa84611297565b82356001600160401b0381116101e8578c611a16918501611349565b8452611a246020840161149a565b6020850152611a35604084016113bc565b60408501526060830135936001600160401b0385116101e857611a5d8d602096879601611349565b60608201528382015281520192019161191e565b9080601f830112156101e857604080519290611a8d90846112cd565b8290604081019283116101e857905b828210611aa95750505090565b8135815260209182019101611a9c565b9080601f830112156101e8578135611ad081611539565b92611ade60405194856112cd565b81845260208085019260051b8201019283116101e857602001905b828210611b065750505090565b60208091611b13846113ab565b815201910190611af9565b519060ff821682036101e857565b51906001600160401b03821682036101e857565b91908260409103126101e857604051611b5881611246565b60206113a6818395611b6981611b2c565b855201611b2c565b519063ffffffff821682036101e857565b519081151582036101e857565b51906001600160801b03821682036101e857565b91908260609103126101e857604051611bbb81611261565b6040808294611bc981611b8f565b8452602081015160208501520151910152565b6020818303126101e8578051906001600160401b0382116101e85701610180818303126101e85760405191611c10836112b2565b81516001600160401b0381116101e8578201918282039261012084126101e85760405193611c3d8561127c565b81516001600160401b0381116101e85782019084601f830112156101e8578151611c66816112f8565b90611c7460405192836112cd565b80825286602082860101116101e8576020815f9282604097018386015e830101528652601f1901126101e85761010090604051611cb081611246565b611cbc60208301611b1e565b8152611cca60408301611b1e565b60208201526020860152611ce18460608301611b40565b6040860152611cf260a08201611b71565b6060860152611d0360c08201611b71565b6080860152611d1460e08201611b82565b60a086015201519060028210156101e857836101409260c0611d829601528552611d418360208301611ba3565b6020860152611d538360808301611ba3565b6040860152611d6460e08201611b8f565b6060860152611d77836101008301611b40565b608086015201611b40565b60a082015290565b9061010060c0611da58451610120855261012085019061114b565b602080860151805160ff9081168784015291015116604085015293611de7604082015160608601906001600160401b0360208092828151168552015116910152565b63ffffffff60608201511660a085015263ffffffff6080820151168285015260a0810151151560e0850152015191611e1e836112ee565b015290565b90606080611e3a845160808552608085019061114b565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b828210611ecc57505050506001600160401b036060611ec2819360208701511515602087015260408701518682036040880152611e23565b9401511691015290565b90919293602080611ee9600193609f198a82030186528851611e23565b960192019201909291611e8a565b91909180519260a0815260206120578551604060a0850152611f3260e0850182516001600160401b0360208092828151168552015116910152565b610240611f50848301516102c06101208801526103a087019061114b565b60408301516001600160401b031661014087015260608301516001600160801b03166101608701526080830151151561018087015260a083015180516101a0880152602090810151805163ffffffff166101c089015201516101e08701529160c0810151151561020087015260e08101516102208701526101008101511515828701526101208101516102608701526101408101516102808701526101608101516102a08701526101808101516102c08701526101a08101516102e08701526101c081015115156103008701526101e08101516103208701526102008101511515610340870152610220810151610360870152015160df198583030161038086015261114b565b94015193609f198282030160c0830152606060c08201956001600160401b03815116835263ffffffff60208201511660208401526120b76040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519460c060a0830152855180915260e0820190602060e08260051b8501019701925f905b82821061213c5750505050506060612104611364949560208501518482036020860152611e63565b9261212c604082015160408501906001600160401b0360208092828151168552015116910152565b0151906080818403910152611e63565b909192939760df19828203018552885190815160048110156110fd576121ba826020600195819594829552015190604084820152606061218883516080604085015260c084019061114b565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f198285030191015261114b565b9a019501939201906120dc565b905f905b600282106121d857505050565b60208060019285518152019301910190916121cb565b90602080835192838152019201905f5b81811061220b5750505090565b825163ffffffff168452602093840193909201916001016121fe565b90602080835192838152019201905f5b8181106122445750505090565b8251845260209384019390920191600101612237565b90602080835192838152019201905f5b8181106122775750505090565b82516001600160401b031684526020938401939092019160010161226a565b90602080835192838152019201905f5b8181106122b35750505090565b825115158452602093840193909201916001016122a6565b9160208383810103126101e8578235906001600160401b0382116101e85761030082850184860103126101e857604051916101a083018381106001600160401b038211176106a257604052808501356001600160401b0381116101e8576123399085870190838801016113c9565b835261234c8486016020838801016114ae565b6020840152608081860101356001600160401b0381116101e8576123779085870190838801016116a6565b604084015261238a60a08287010161149a565b606084015283850160df8287010112156101e8576040516123ad610100826112cd565b806101c0838801019186880183116101e85760c084890101905b83821061311557505060808501526123e29086860190611a71565b60a08401526123f984860161020083880101611a71565b60c0840152610240818601013561ffff811681036101e85760e084015261026081860101356001600160401b0381116101e85761243d908587019083880101611ab9565b61010084015261028081860101356001600160401b0381116101e857848601601f82848901010112156101e85780828701013561247981611539565b9161248760405193846112cd565b818352602083019087890160208460051b83888d01010101116101e85790602082868b010101915b60208460051b82888d01010101831061310157505050506101208401526102a081860101356001600160401b0381116101e857848601601f82848901010112156101e85780828701013561250281611539565b9161251060405193846112cd565b818352602083019087890160208460051b83888d01010101116101e85790602082868b010101915b60208460051b82888d0101010183106130e657505050506101408401526102c081860101356001600160401b0381116101e85761257c908587019083880101611ab9565b6101608401526102e08186010135906001600160401b0382116101e857850101838501601f820112156101e8578035906125b582611539565b916125c360405193846112cd565b80835260208084019160051b8301019186880183116101e857602001905b8282106130ce575050506101808301526125fa82613253565b9097929690959391929115612f745750506040517fb4e415790000000000000000000000000000000000000000000000000000000081526020600482015261269161265384516103006024850152610324840190611d8a565b60208581015180516001600160801b0316604486015290810151606485015260400151608484015260408501518382036023190160a4850152611ef7565b6001600160801b0360608501511660c483015281608085015160e482015f905b60088210612f5a575050508061277b61276361274b61273361271b5f976126e160a08d01516101e48901906121c7565b6126f460c08d01516102248901906121c7565b60e08c015161ffff166102648801526101008c0151878203602319016102848901526121ee565b6101208b0151868203602319016102a4880152612227565b6101408a0151858203602319016102c487015261225a565b610160890151848203602319016102e48601526121ee565b61018088015183820360231901610304850152612296565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa908115610697575f91612f38575b50945b6127d386516001600160801b036060890151169061446b565b6128356020870151604051612809602082018093604080916001600160801b038151168452602081015160208501520151910152565b606081526128186080826112cd565b51902061059f6001600160401b03602060808b0151015116613125565b61283e86613458565b9360206040850151015151948551976101008601515161ffff60e0880151168091149081612f28575b81612f18575b81612f08575b81612ef8575b5015612ed057909192985f96959496975f5b8a8110612e9857505f9a8b9889805b6101008a0151518c1015612a01576128b78c6101808c01516131d2565b51156129ea57908d989796959493929163ffffffff6128db8e6101008e01516131d2565b5116998a10156129be5761297f575b5060206128f7898c6131d2565b5101516129098c6101208c01516131d2565b510361295357600161293c8b9c9d9e9f6001600160401b0360406129329d9e9d859d809f6131d2565b5101511690613506565b9c5b019a999890919293949596979d9c9b9d61289a565b877fc3271bcd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b63ffffffff16881115612992575f6128ea565b877fea334314000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b897fc9d6366c000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b979695949392919098999a6001909e9c9d9e61293e565b5050949a50949a909597506001600160401b03919650979197166003810281810460031482151715612e84576801fffffffffffffffe6001600160401b0384169360011b169083820460021484151715612e84571115612e56575050604085015151602081015160046001600160401b03825116604063ffffffff60208501511693015193602085519501519060208063ffffffff845116930151915101519160405193612aae856112b2565b84526020840195865260408401968752606084019081526080840191825260a0840192835261ffff60e08d0151169660808d0151968d60a08101519060c081015161012082015161014083015191610180610160850151940151946040519e8f9d8e7f4cc22bb7000000000000000000000000000000000000000000000000000000008152015260248d019d5f5b60088110612e2257505060209d50612be38d63ffffffff97612bd0829f9d9a9598612c1e9f9c98612bbd9060c09f9b6001600160401b039a612b96612baa92612b8b8e9c6101248c01906121c7565b6101648a01906121c7565b6102406101a4890152610244880190612227565b868103600319016101c48801529061225a565b848103600319016101e4860152906121ee565b9161020460031982850301910152612296565b8c8103600319016102248e0152995116895251168b880152516040870152511660608501525160808401525160a0830182905291019061114b565b03815f6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af1908115610697575f91612de8575b5015612dc057612da8575b50612d8f575b5050612c77826110f3565b81612d46576001600160401b036020604060a0840193845183810151600354918683861c1687831611612cfd575b5050500151604051612cd78382018093604080916001600160801b038151168452602081015160208501520151910152565b60608152612ce66080826112cd565b51902092510151165f52600560205260405f205590565b6fffffffffffffffff0000000000000000877fffffffffffffffffffffffffffffffff0000000000000000000000000000000092511692861b16921617176003555f8080612ca5565b50612d50816110f3565b60018103612d78576801000000000000000068ff000000000000000019600454161760045590565b612d81816110f3565b600281036113645750600290565b60606040612da193015101519061355c565b5f80612c6c565b604084015160200151612dba9161355c565b5f612c66565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90506020813d602011612e1a575b81612e03602093836112cd565b810103126101e857612e1490611b82565b5f612c5b565b3d9150612df6565b91939597999b9d5091939597999b9d602080600192855181520193019101908f9d9b99979593919e9c9a989694929e612b3c565b7f6e3083c3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b634e487b7160e01b5f52601160045260245ffd5b9a949392919098612ebd6001916001600160401b0360408f612932908d9e9c9d6131d2565b9b019a9890919293949a9796959761288b565b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b905061018087015151145f612879565b6101608801515181149150612873565b610140880151518114915061286d565b6101208801515181149150612867565b612f5491503d805f833e612f4c81836112cd565b810190611bdc565b5f6127b7565b8293506020809160019394518152019301910184926126b1565b5f8091604051612f83816112b2565b604051612f8f8161127c565b60608152604051612f9f81611246565b8481528460208201526020820152604051612fb981611246565b84815284602082015260408201528360608201528360808201528360a08201528360c08201528152612fe9613403565b6020820152612ff6613403565b604082015282606082015260405161300d81611246565b838152836020820152608082015260a06040519161302a83611246565b848352846020840152015280604051947f513a5c0e0000000000000000000000000000000000000000000000000000000086526004860137600401836001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa156130c6576130c09060203d9182815281810192805f853e8101603f01601f19166040528051010190611bdc565b946127ba565b3d5f823e3d90fd5b602080916130db846113bc565b8152019101906125e1565b60208080936130f486611111565b8152019301929150612538565b6020808093853581520193019291506124af565b81358152602091820191016123c7565b6001600160401b03165f52600560205260405f205480156131435790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b903590601e19813603018212156101e857018035906001600160401b0382116101e8576020019181360383136101e857565b903590601e19813603018212156101e857018035906001600160401b0382116101e857602001918160051b360383136101e857565b80518210156131e65760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b805f525f60205260405f206001600160a01b0333165f5260205260ff60405f205416156132245750565b7fe2517d3f000000000000000000000000000000000000000000000000000000005f523360045260245260445ffd5b905f905f604084019182519283519360016001600160401b0360206040828180846101408d5101519e015101519a51015116940151015116016001600160401b038111612e84576001600160401b031614865f52600660205260ff60405f2054169087861490801592836133de575b8081156133d7575b80156133d0575b80156133c9575b156133b6571561339d576132eb8961421b565b6020855101525b61334a57613332571561331857606061330a8561421b565b915101525b60019493929190565b61332a9192506060905101518361404f565b60019061330f565b50516020810151606090910152600194939291505f90565b50509091506133576141f0565b906040516133666020826112cd565b5f81525f805b818110613386575050825251606001526001939291905f90565b6020906133916141cc565b8282860101520161336c565b96506133ae6020845101518961404f565b6001966132f2565b50505050509250505f9291600191908290565b50826132d8565b50836132d1565b50816132ca565b925081156133ee57825b926132c2565b865f52600660205260ff60405f2054166133e8565b6040519061341082611261565b5f6040838281528260208201520152565b1561342a575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160401b03602060a08301510151165f52600560205260405f205480155f146134845750505f90565b604082019081516040516134b9602082018093604080916001600160801b038151168452602081015160208501520151910152565b606081526134c86080826112cd565b51902014918215926134e6575b5050156134e157600190565b600290565b6001600160801b0391925060208291015151169151511611155f806134d5565b906001600160401b03809116911601906001600160401b038211612e8457565b80548210156131e6575f5260205f2001905f90565b91909180548310156131e6575f52601860205f208360021c019260031b1690565b805f52600660205260ff60405f2054166137ac575f92919252600660205260405f2091600160ff198454161783555f92600181016003600283019201945b835180518210156137a357816135af916131d2565b51518254680100000000000000008110156106a2578060016135d49201855584613526565b919091613790578051906001600160401b0382116106a25781906135f8845461116f565b601f8111613740575b50602090601f83116001146136dd575f926136d2575b50508160011b915f199060031b1c19161790555b60206136388286516131d2565b5101518354680100000000000000008110156106a25780600161365e9201865585613526565b819291549060031b91821b915f19901b19161790556001600160401b0360406136888387516131d2565b5101511690865491680100000000000000008310156106a2576136b28360018095018a558961353b565b6001600160401b03829392549160031b92831b921b19161790550161359a565b015190505f80613617565b5f8581528281209350601f198516905b8181106137285750908460019594939210613710575b505050811b01905561362b565b01515f1960f88460031b161c191690555f8080613703565b929360206001819287860151815501950193016136ed565b909150835f5260205f20601f840160051c81019160208510613786575b90601f859493920160051c01905b8181106137785750613601565b5f815584935060010161376b565b909150819061375d565b634e487b7160e01b5f525f60045260245ffd5b50505050509050565b5050565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f205416155f1461383757805f525f60205260405f206001600160a01b0383165f5260205260405f20600160ff198254161790556001600160a01b03339216907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f2054165f1461383757805f525f60205260405f206001600160a01b0383165f5260205260405f2060ff1981541690556001600160a01b03339216907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b999493979198909695995f966138d5816112ee565b8061400e575089151580614002575b15613fcb576020019460208b6139516139076138ff8a61468c565b9236906114ae565b9161059f604051858101906139398287604080916001600160801b038151168452602081015160208501520151910152565b606081526139486080826112cd565b51902091613125565b0151808603613f9b57505f905f5b8b8a818310613eae575b5050505015613e6d5750506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690813b156101e8579187926040519384927f87d3a9b100000000000000000000000000000000000000000000000000000000845260648401906004850152606060248501525260848201908960051b9860848a850101928b8a915f603e198d360301925b8110613dff5750505050600319848403016044850152808352602083019260208260051b82010193835f92601e1982360301905b858510613c90575050505050505091815f81819503925af1801561069757613c7b575b5060018511613a86575b50505050506001600160801b03613a80633b9aca009236906114ae565b51160490565b613a9790969192949695939561468c565b948335946001600160801b038616809603613c7757613ab588611539565b97613ac3604051998a6112cd565b88526020880191810190368211613c735780979597925b828410613bd857505050506001600160401b03829316925b8651811015613bbb57613b0581886131d2565b519085604051602081019087825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801918a5b818110613b84575050505081613b6c6001976020613b7a940151605f1984830301608085015261114b565b03601f1981018352826112cd565b5190205d01613af2565b919394965091949697602080613ba660019360bf198b8203018852895161114b565b970194019101918c9694939298979598613b41565b5094505050506001600160801b03613a80633b9aca005f80613a63565b83989698356001600160401b038111613c6f578201604081360312613c6f5760405190613c0482611246565b80356001600160401b038111613c6b57810136601f82011215613c6b57613c329036906020813591016146c2565b825260208101356001600160401b038111613c6b5791613c59602094928594369101611349565b83820152815201930192979597613ada565b8880fd5b8680fd5b8480fd5b8380fd5b613c889192505f906112cd565b5f905f613a59565b9193959750919395601f198282030185528735838112156101e857840190613cbc6020820192806147dd565b8091936020845252604082019060408160051b8401019380935f915b838310613cff57505050505050602080600192990195019501929091899796949592613a36565b909192939495603f19838203018652613d188783614825565b803560028110156101e857613d2c816112ee565b8252613d4f613d3e6020830183614811565b606060208501526060840190614839565b906040810135609e19823603018112156101e8576001936020938493613df19301916040818303910152613de3613dc5613d9a613d8c858061474d565b60a0865260a086019161472d565b613da58786016113bc565b151587850152613db86040860186614811565b8482036040860152614839565b92613dd2606082016113bc565b151560608401526080810190614811565b906080818403910152614839565b980196019493019190613cd8565b919394965091946083198982030183528535848112156101e8576020613e5b6001938f83940190613e4e613e44613e3684806147dd565b60408552604085019161477e565b928581019061474d565b918581850391015261472d565b9701930191019088969493918e613a02565b613eaa6040519283927ffef760c700000000000000000000000000000000000000000000000000000000845260206004850152602484019161477e565b0390fd5b613ecc613ec684613ee494613ed394989697986146a0565b8061319d565b36916146c2565b613ede3687896146c2565b90614a4c565b15613f915750610c15613efb613f05928d8c6146a0565b602081019061316b565b805182518082149182613f7b575b505015613f2757505060015f808b8a613969565b90613eaa613f69926040519384937f5f1ca38100000000000000000000000000000000000000000000000000000000855260406004860152604485019061114b565b8381036003190160248501529061114b565b9091506020830120906020840120145f80613f13565b919060010161395f565b85907f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b897fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b5061ffff8a11156138e4565b8760ff60249261401d816112ee565b614026816112ee565b7f112d89cc00000000000000000000000000000000000000000000000000000000835216600452fd5b815151156141ba57604051917fc87f1f69000000000000000000000000000000000000000000000000000000008352602060048401528260a4810182519060806024840152815180915260c4830190602060c48260051b8601019301915f905b82821061418b5750505050826001600160401b0360606140ed8594602080980151151560448701526040850151602319878303016064880152611e23565b920151166084830152038173a3c616dd54F6BB35a736cD6968c8EF7176faCACc5af4918215610697575f92614157575b50808203614129575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b9091506020813d602011614183575b81614173602093836112cd565b810103126101e85751905f61411d565b3d9150614166565b919360019193955060206141aa819260c3198c82030186528851611e23565b96019201920187949391926140af565b633a517eed60e21b5f5260045260245ffd5b604051906141d982611297565b5f6060838181528260208201528260408201520152565b604051906141fd82611297565b5f6060838181528260208201526142126141cc565b60408201520152565b6142236141f0565b50805f52600660205260405f2060ff81541615614458576002810190815460018201938454821480159061444a575b61441f57509061426182611539565b9361426f60405195866112cd565b828552601f1961427e84611539565b015f5b81811061440257505060035f9201915b8381106142f15750505050506040516142a981611297565b6040516142b76020826112cd565b5f815281525f60208201525f60408201525f6060820152604051916142db83611297565b82525f602083015260408201525f606082015290565b6142fb8183613526565b506143068287613526565b90549060031b1c6001600160401b0361431f848761353b565b90549060031b1c16906040519261433584611297565b604051905f908054906143478261116f565b80855291600181169081156143db57506001146143a2575b505090614374816001979695949303826112cd565b8352602083015260408201525f606082015261439082896131d2565b5261439b81886131d2565b5001614291565b5f908152602081209092505b8183106143c557505081016020016143748261435f565b60018160209254838688010152019201916143ae565b60ff191660208087019190915292151560051b85019092019250614374915083905061435f565b6020906144109594956141cc565b82828a01015201939293614281565b7f4724a0fd000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b506003830154821415614252565b50633a517eed60e21b5f5260045260245ffd5b906001600160801b03633b9aca0091160442811161465d57804203428111612e84576107081061462e57508051805160208201207f0000000000000000000000000000000000000000000000000000000000000000036145ee5750602081015160ff815116906002549160ff83169283821492836145d5575b6020015160ff16921561458d575050505063ffffffff606082015116906004549163ffffffff831680820361455f57505063ffffffff608081920151169160201c16808203614531575050565b7f1eb5c1eb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f73b1bde3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6084945060ff90604051947fd382033a000000000000000000000000000000000000000000000000000000008652600486015260081c16602484015260448301526064820152fd5b6020810151600883901c60ff90811691161493506144e4565b613eaa906040519182917ff6b6676b00000000000000000000000000000000000000000000000000000000835260406004840152613f69604484016111a7565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b7f20daedb6000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b356001600160401b03811681036101e85790565b91908110156131e65760051b81013590603e19813603018212156101e8570190565b9291906146ce81611539565b936146dc60405195866112cd565b602085838152019160051b8101918383116101e85781905b838210614702575050505050565b81356001600160401b0381116101e8576020916147228784938701611349565b8152019101906146f4565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156101e85701602081359101916001600160401b0382116101e85781360383136101e857565b90602083828152019260208260051b82010193835f925b8484106147a55750505050505090565b9091929394956020806147cd600193601f198682030188526147c78b8861474d565b9061472d565b9801940194019294939190614795565b9035601e19823603018112156101e85701602081359101916001600160401b0382116101e8578160051b360383136101e857565b9035607e19823603018112156101e8570190565b9035605e19823603018112156101e8570190565b614872614857614849838061474d565b60808652608086019161472d565b614864602084018461474d565b90858303602087015261472d565b61487f6040830183614811565b908381036040850152813560038110156101e85761489c816110f3565b8152602082013560038110156101e8576148b5816110f3565b602082015260408201359060038210156101e85760806148f161490c94846148e2614901969998996110f3565b6040850152606081019061474d565b919092816060820152019161472d565b9260608101906147dd565b90916060818503910152808352602083019060208160051b85010193835f915b83831061493c5750505050505090565b909192939495601f198282030186526149558784614825565b80359160038310156101e8576149b26020928392856149756001976110f3565b81526149a46149996149898685018561474d565b606088860152606085019161472d565b92604081019061474d565b91604081850391015261472d565b98019601949301919061492c565b906149fd57508051156149d557602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b81511580614a43575b614a0e575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15614a06565b908151815103613837575f5b8251811015614ab657614a6b81846131d2565b5151614a7782846131d2565b515103614aaf57614a8881846131d2565b5160208151910120614a9a82846131d2565b516020815191012003614aaf57600101614a58565b5050505f90565b50505060019056fea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0db10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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

// ALLOWEDCLOCKDRIFT is a free data retrieval call binding the contract method 0xaef1f78a.
//
// Solidity: function ALLOWED_CLOCK_DRIFT() view returns(uint16)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) ALLOWEDCLOCKDRIFT(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "ALLOWED_CLOCK_DRIFT")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// ALLOWEDCLOCKDRIFT is a free data retrieval call binding the contract method 0xaef1f78a.
//
// Solidity: function ALLOWED_CLOCK_DRIFT() view returns(uint16)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) ALLOWEDCLOCKDRIFT() (uint16, error) {
	return _ContractGroth16ICS07Tendermint.Contract.ALLOWEDCLOCKDRIFT(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// ALLOWEDCLOCKDRIFT is a free data retrieval call binding the contract method 0xaef1f78a.
//
// Solidity: function ALLOWED_CLOCK_DRIFT() view returns(uint16)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) ALLOWEDCLOCKDRIFT() (uint16, error) {
	return _ContractGroth16ICS07Tendermint.Contract.ALLOWEDCLOCKDRIFT(&_ContractGroth16ICS07Tendermint.CallOpts)
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

// MEMBERSHIP is a free data retrieval call binding the contract method 0x89df51f1.
//
// Solidity: function MEMBERSHIP() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) MEMBERSHIP(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "MEMBERSHIP")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MEMBERSHIP is a free data retrieval call binding the contract method 0x89df51f1.
//
// Solidity: function MEMBERSHIP() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) MEMBERSHIP() (common.Address, error) {
	return _ContractGroth16ICS07Tendermint.Contract.MEMBERSHIP(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// MEMBERSHIP is a free data retrieval call binding the contract method 0x89df51f1.
//
// Solidity: function MEMBERSHIP() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) MEMBERSHIP() (common.Address, error) {
	return _ContractGroth16ICS07Tendermint.Contract.MEMBERSHIP(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// MISBEHAVIOUR is a free data retrieval call binding the contract method 0x02cf2952.
//
// Solidity: function MISBEHAVIOUR() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) MISBEHAVIOUR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "MISBEHAVIOUR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MISBEHAVIOUR is a free data retrieval call binding the contract method 0x02cf2952.
//
// Solidity: function MISBEHAVIOUR() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) MISBEHAVIOUR() (common.Address, error) {
	return _ContractGroth16ICS07Tendermint.Contract.MISBEHAVIOUR(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// MISBEHAVIOUR is a free data retrieval call binding the contract method 0x02cf2952.
//
// Solidity: function MISBEHAVIOUR() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) MISBEHAVIOUR() (common.Address, error) {
	return _ContractGroth16ICS07Tendermint.Contract.MISBEHAVIOUR(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// PROOFSUBMITTERROLE is a free data retrieval call binding the contract method 0x5972185a.
//
// Solidity: function PROOF_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) PROOFSUBMITTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "PROOF_SUBMITTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PROOFSUBMITTERROLE is a free data retrieval call binding the contract method 0x5972185a.
//
// Solidity: function PROOF_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) PROOFSUBMITTERROLE() ([32]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.PROOFSUBMITTERROLE(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// PROOFSUBMITTERROLE is a free data retrieval call binding the contract method 0x5972185a.
//
// Solidity: function PROOF_SUBMITTER_ROLE() view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) PROOFSUBMITTERROLE() ([32]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.PROOFSUBMITTERROLE(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// UPDATECLIENT is a free data retrieval call binding the contract method 0x87d4332f.
//
// Solidity: function UPDATE_CLIENT() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) UPDATECLIENT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "UPDATE_CLIENT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UPDATECLIENT is a free data retrieval call binding the contract method 0x87d4332f.
//
// Solidity: function UPDATE_CLIENT() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) UPDATECLIENT() (common.Address, error) {
	return _ContractGroth16ICS07Tendermint.Contract.UPDATECLIENT(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// UPDATECLIENT is a free data retrieval call binding the contract method 0x87d4332f.
//
// Solidity: function UPDATE_CLIENT() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) UPDATECLIENT() (common.Address, error) {
	return _ContractGroth16ICS07Tendermint.Contract.UPDATECLIENT(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// VERIFIER is a free data retrieval call binding the contract method 0x08c84e70.
//
// Solidity: function VERIFIER() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) VERIFIER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "VERIFIER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VERIFIER is a free data retrieval call binding the contract method 0x08c84e70.
//
// Solidity: function VERIFIER() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) VERIFIER() (common.Address, error) {
	return _ContractGroth16ICS07Tendermint.Contract.VERIFIER(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// VERIFIER is a free data retrieval call binding the contract method 0x08c84e70.
//
// Solidity: function VERIFIER() view returns(address)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) VERIFIER() (common.Address, error) {
	return _ContractGroth16ICS07Tendermint.Contract.VERIFIER(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// ClientState is a free data retrieval call binding the contract method 0xbd3ce6b0.
//
// Solidity: function clientState() view returns(string chainId, (uint8,uint8) trustLevel, (uint64,uint64) latestHeight, uint32 trustingPeriod, uint32 unbondingPeriod, bool isFrozen, uint8 zkAlgorithm)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) ClientState(opts *bind.CallOpts) (struct {
	ChainId         string
	TrustLevel      IICS07TendermintMsgsTrustThreshold
	LatestHeight    IICS02ClientMsgsHeight
	TrustingPeriod  uint32
	UnbondingPeriod uint32
	IsFrozen        bool
	ZkAlgorithm     uint8
}, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "clientState")

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
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) ClientState() (struct {
	ChainId         string
	TrustLevel      IICS07TendermintMsgsTrustThreshold
	LatestHeight    IICS02ClientMsgsHeight
	TrustingPeriod  uint32
	UnbondingPeriod uint32
	IsFrozen        bool
	ZkAlgorithm     uint8
}, error) {
	return _ContractGroth16ICS07Tendermint.Contract.ClientState(&_ContractGroth16ICS07Tendermint.CallOpts)
}

// ClientState is a free data retrieval call binding the contract method 0xbd3ce6b0.
//
// Solidity: function clientState() view returns(string chainId, (uint8,uint8) trustLevel, (uint64,uint64) latestHeight, uint32 trustingPeriod, uint32 unbondingPeriod, bool isFrozen, uint8 zkAlgorithm)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) ClientState() (struct {
	ChainId         string
	TrustLevel      IICS07TendermintMsgsTrustThreshold
	LatestHeight    IICS02ClientMsgsHeight
	TrustingPeriod  uint32
	UnbondingPeriod uint32
	IsFrozen        bool
	ZkAlgorithm     uint8
}, error) {
	return _ContractGroth16ICS07Tendermint.Contract.ClientState(&_ContractGroth16ICS07Tendermint.CallOpts)
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

// GetConsensusStateHash is a free data retrieval call binding the contract method 0x23842fb8.
//
// Solidity: function getConsensusStateHash(uint64 revisionHeight) view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) GetConsensusStateHash(opts *bind.CallOpts, revisionHeight uint64) ([32]byte, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "getConsensusStateHash", revisionHeight)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetConsensusStateHash is a free data retrieval call binding the contract method 0x23842fb8.
//
// Solidity: function getConsensusStateHash(uint64 revisionHeight) view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) GetConsensusStateHash(revisionHeight uint64) ([32]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GetConsensusStateHash(&_ContractGroth16ICS07Tendermint.CallOpts, revisionHeight)
}

// GetConsensusStateHash is a free data retrieval call binding the contract method 0x23842fb8.
//
// Solidity: function getConsensusStateHash(uint64 revisionHeight) view returns(bytes32)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) GetConsensusStateHash(revisionHeight uint64) ([32]byte, error) {
	return _ContractGroth16ICS07Tendermint.Contract.GetConsensusStateHash(&_ContractGroth16ICS07Tendermint.CallOpts, revisionHeight)
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

// HasCachedValidatorSet is a free data retrieval call binding the contract method 0x083a1fee.
//
// Solidity: function hasCachedValidatorSet(bytes32 validatorsHash) view returns(bool)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCaller) HasCachedValidatorSet(opts *bind.CallOpts, validatorsHash [32]byte) (bool, error) {
	var out []interface{}
	err := _ContractGroth16ICS07Tendermint.contract.Call(opts, &out, "hasCachedValidatorSet", validatorsHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasCachedValidatorSet is a free data retrieval call binding the contract method 0x083a1fee.
//
// Solidity: function hasCachedValidatorSet(bytes32 validatorsHash) view returns(bool)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) HasCachedValidatorSet(validatorsHash [32]byte) (bool, error) {
	return _ContractGroth16ICS07Tendermint.Contract.HasCachedValidatorSet(&_ContractGroth16ICS07Tendermint.CallOpts, validatorsHash)
}

// HasCachedValidatorSet is a free data retrieval call binding the contract method 0x083a1fee.
//
// Solidity: function hasCachedValidatorSet(bytes32 validatorsHash) view returns(bool)
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintCallerSession) HasCachedValidatorSet(validatorsHash [32]byte) (bool, error) {
	return _ContractGroth16ICS07Tendermint.Contract.HasCachedValidatorSet(&_ContractGroth16ICS07Tendermint.CallOpts, validatorsHash)
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
// Solidity: function upgradeClient(bytes ) view returns()
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
// Solidity: function upgradeClient(bytes ) view returns()
func (_ContractGroth16ICS07Tendermint *ContractGroth16ICS07TendermintSession) UpgradeClient(arg0 []byte) error {
	return _ContractGroth16ICS07Tendermint.Contract.UpgradeClient(&_ContractGroth16ICS07Tendermint.CallOpts, arg0)
}

// UpgradeClient is a free data retrieval call binding the contract method 0x8a8e4c5d.
//
// Solidity: function upgradeClient(bytes ) view returns()
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
