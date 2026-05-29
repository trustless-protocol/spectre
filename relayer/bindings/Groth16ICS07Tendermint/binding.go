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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"membership_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"misbehaviour_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"updateClient_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_clientState\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_consensusState\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"roleManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ALLOWED_CLOCK_DRIFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MEMBERSHIP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMembership\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MISBEHAVIOUR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMisbehaviour\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PROOF_SUBMITTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPDATE_CLIENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUpdateClient\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VERIFIER\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"clientState\",\"inputs\":[],\"outputs\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCachedValidatorSet\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"indices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"votingPowers\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConsensusStateHash\",\"inputs\":[{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasCachedValidatorSet\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"updateClientMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeClient\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyNonMembership\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structILightClientMsgs.MsgVerifyNonMembership\",\"components\":[{\"name\":\"height\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"kvPairs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.KVPair[]\",\"components\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"merkleProofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.MerkleProof[]\",\"components\":[{\"name\":\"proofs\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.CommitmentProof[]\",\"components\":[{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.ProofType\"},{\"name\":\"existenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonExistenceProof\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.NonExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hasLeft\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"left\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"hasRight\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"right\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.ExistenceProof\",\"components\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"leaf\",\"type\":\"tuple\",\"internalType\":\"structIMembershipMsgs.LeafOp\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashKey\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prehashValue\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structIMembershipMsgs.InnerOp[]\",\"components\":[{\"name\":\"hashOp\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.HashOp\"},{\"name\":\"prefix\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"suffix\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]}]},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"enumIMembershipMsgs.MembershipType\"},{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BatchLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CachedSignerNotFound\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"CachedValidatorSetCorrupted\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CannotHandleMisbehavior\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ClientStateMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateRootMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"EmptyValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"FeatureNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenClientState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMisbehaviourHeaderHeight\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[{\"name\":\"accumulated\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"total\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidMembershipProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"KeyValuePairNotInCache\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"LengthIsOutOfRange\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"min\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MembershipProofKeyNotFound\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"MembershipProofValueMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedValidatorHashes\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ProofHeightMismatch\",\"inputs\":[{\"name\":\"expectedRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"expectedRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualRevisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofIsInTheFuture\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofIsTooOld\",\"inputs\":[{\"name\":\"now\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubkeyMismatch\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"SSTORE2DataTooLarge\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SSTORE2InvalidPointer\",\"inputs\":[{\"name\":\"pointer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SSTORE2WriteFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignerIndexOutOfRange\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"TrustThresholdMismatch\",\"inputs\":[{\"name\":\"expectedNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expectedDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualNumerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualDenominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TrustingPeriodTooLong\",\"inputs\":[{\"name\":\"trustingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnbondingPeriodMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownMembershipType\",\"inputs\":[{\"name\":\"membershipType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"UnknownZkAlgorithm\",\"inputs\":[{\"name\":\"algorithm\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"ValidatorSetCacheMiss\",\"inputs\":[{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"VerificationKeyMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x61014080604052346105ec5761589b803803809161001d828561060b565b8339810160e0828203126105ec576100348261062e565b6100406020840161062e565b61004c6040850161062e565b916100596060860161062e565b60808601519094906001600160401b0381116105ec5786019080601f830112156105ec57815161008b92602001610642565b9461009d60c060a0830151920161062e565b957fbd893629a699470e4ec82a5715bb4981fdaacc5d0a728bf5f55b801d8f4ef10b61010052805181019060208201906020818403126105ec576020810151906001600160401b0382116105ec570191829003601f19810190610120136105ec576040519160e083016001600160401b038111848210176105d85760405260208401516001600160401b0381116105ec576020908501019080601f830112156105ec57815161014e92602001610642565b8252604081126105ec576040805191610166836105f0565b610171828601610687565b835261017f60608601610687565b602084015260208401928352603f1901126105ec576040516101a0816105f0565b6101ac60808501610695565b81526101ba60a08501610695565b6020820152604083019081526101d260c085016106a9565b90606084019182526101e660e086016106a9565b92608085019384526101008601519586151587036105ec5760a0860196875261012001519460028610156105ec5760c08101958652518051906001600160401b0382116105d8576102386001546106ba565b601f8111610588575b50602090601f831160011461051b5763ffffffff95949392915f9183610510575b50508160011b915f199060031b1c1916176001555b5160ff81511661ff00602060025493015160081b169161ffff191617176002555160018060401b0381511660035491602068010000000000000000600160801b0391015160401b169160018060801b031916171760035551169267ffffffff00000000600454925160201b169051151560401b92519160028310156104fc5769ff00000000000000000068ff00000000000000009360481b169469ff000000000000000000199160018060481b0319161716179116171760045560405161034881610341816106f2565b038261060b565b602081519101206101205260405161036381610341816106f2565b6001600160401b0390602090610378906107a3565b0151600354916001600160401b03831691168181036104e7575050604090811c6001600160401b03165f908152600560205220556001600160a01b0390811660805290811660a05290811660c0521660e05260045463ffffffff8082169061070882019081116104d35763ffffffff809360201c1692839116116104be57826001600160a01b0381166104a15750610412610100516109fa565b505b604051614d729081610ac98239608051818181610fd301526147af015260a051818181610d5d0152613b2f015260c0518181816104e1015261103e015260e051818181610da0015281816127eb0152612eaa01526101005181818161021601528181610a5c01528181610b7701528181610cc901528181610ddb0152610f520152610120518161440d0152f35b806104ae6104b892610984565b5061010051610a53565b50610414565b6333fae18560e21b5f5260045260245260445ffd5b634e487b7160e01b5f52601160045260245ffd5b6355bace6f60e01b5f5260045260245260445ffd5b634e487b7160e01b5f52602160045260245ffd5b015190505f80610262565b90601f1983169160015f52815f20925f5b818110610570575091600193918563ffffffff999897969410610558575b505050811b01600155610277565b01515f1960f88460031b161c191690555f808061054a565b9293602060018192878601518155019501930161052c565b60015f525f51602061585b5f395f51905f52601f840160051c810191602085106105ce575b601f0160051c01905b8181106105c35750610241565b5f81556001016105b6565b90915081906105ad565b634e487b7160e01b5f52604160045260245ffd5b5f80fd5b604081019081106001600160401b038211176105d857604052565b601f909101601f19168101906001600160401b038211908210176105d857604052565b51906001600160a01b03821682036105ec57565b9192916001600160401b0382116105d8576040519161066b601f8201601f19166020018461060b565b8294818452818301116105ec578281602093845f96015e010152565b519060ff821682036105ec57565b51906001600160401b03821682036105ec57565b519063ffffffff821682036105ec57565b90600182811c921680156106e8575b60208310146106d457565b634e487b7160e01b5f52602260045260245ffd5b91607f16916106c9565b6001545f9291610701826106ba565b8082529160018116908115610762575060011461071c575050565b60015f9081529293509091905f51602061585b5f395f51905f525b838310610748575060209250010190565b600181602092949394548385870101520191019190610737565b9050602093945060ff929192191683830152151560051b010190565b90815181101561078f570160200190565b634e487b7160e01b5f52603260045260245ffd5b6040516107af816105f0565b606081525f60208201525080518015908115610978575b50610969575f19908051805b610925575b505f1982146109175760018201908183116104d357600360fc1b6001600160f81b0319610804848461077e565b51161480610902575b6108f2575f5b81518310156108a257610826838361077e565b5160f81c603081108015610898575b61088657600a82026001600160401b03908116602f1990920160ff169190910181169116811061086a57600190920191610813565b509150506040519061087b826105f0565b81525f602082015290565b50509150506040519061087b826105f0565b5060398111610835565b91509160018111908115916108e6575b506108d757604051916108c4836105f0565b82526001600160401b0316602082015290565b63384c687760e21b5f5260045ffd5b602b915010155f6108b2565b9150506040519061087b826105f0565b5080518381039081116104d35760021061080d565b604051915061087b826105f0565b5f1981018181116104d357602d60f81b6001600160f81b0319610948838661077e565b51161461095f575080156104d3575f1901806107d2565b92505f90506107d7565b6329120bff60e21b5f5260045ffd5b6040915010155f6107c6565b6001600160a01b0381165f9081525f51602061587b5f395f51905f52602052604090205460ff166109f5576001600160a01b03165f8181525f51602061587b5f395f51905f5260205260408120805460ff191660011790553391905f51602061583b5f395f51905f528180a4600190565b505f90565b805f525f60205260405f205f805260205260ff60405f205416155f146109f557805f525f60205260405f205f805260205260405f20600160ff198254161790555f33915f51602061583b5f395f51905f528280a4600190565b5f818152602081815260408083206001600160a01b038616845290915290205460ff16610ac2575f818152602081815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291905f51602061583b5f395f51905f529080a4600190565b50505f9056fe6080806040526004361015610012575f80fd5b5f3560e01c90816301ffc9a7146110625750806302cf29521461101f578063083a1fee14610ff757806308c84e7014610fb45780630bece35614610f2d57806323842fb814610efe578063248a9ca314610ed45780632f2ff15d14610ea557806336568abe14610e56578063536c2ad314610dfe5780635972185a14610dc457806387d4332f14610d8157806389df51f114610d3e5780638a8e4c5d14610caa57806391d1485414610c6e578063974a74c414610b32578063a217fddf14610b18578063a6f031bb14610a17578063ac9650d814610852578063aef1f78a14610836578063bd3ce6b01461072e578063d547741f146106f8578063ddba6537146101f75763ef913a4b14610124575f80fd5b346101f3575f3660031901126101f3576101ef60405160208082015261012060408201526101db816101596101608201611273565b60ff600254818116606085015260081c1660808301526001600160401b0360035481811660a085015260401c1660c083015260ff60045463ffffffff811660e085015263ffffffff8160201c16610100850152818160401c16151561012085015260481c166101c7816113e5565b61014083015203601f1981018352826113c4565b60405191829160208352602083019061124f565b0390f35b5f80fd5b346101f35761020536611100565b6004549060ff8260401c166106d0577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f205416156106c1575b508201916020818403126101f3578035906001600160401b0382116101f35701610120818403126101f3576040519060a082018281106001600160401b038211176106ad5760405280356001600160401b0381116101f357846102b89183016114c0565b825260208101356001600160401b0381116101f3578101936060858203126101f357604051946102e786611358565b80356001600160401b0381116101f35781016040818403126101f357604051906103108261133d565b80356001600160401b0381116101f357816103328660209361033a9501611440565b84520161116d565b6020820152865260208101356001600160401b0381116101f3578261036091830161179d565b602087015260408101356001600160401b0381116101f3576103849183910161179d565b60408601526020830194855280604083019061039f916115a5565b60408401908152906103b49060a084016115a5565b926060810192848452610100016103ca90611591565b9560808201948786528251915197845191604051998a947fa6fe8f56000000000000000000000000000000000000000000000000000000008652600486016101209052610124860161041b91611e81565b906003198683030160248701528051606083528051606084016040905260a084016104459161124f565b90602001516001600160401b03166080840152602082015190838103602085015261046f91611fee565b906040015191808203906040015261048691611fee565b83516001600160801b0316604486015260208401516064860152604090930151608485015280516001600160801b031660a4850152602081015160c48501526040015160e48401526001600160801b031661010483015203867f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031691815a93608094fa9586156106a2575f9661060e575b506020806105fc956105b268010000000000000000999661055c6105aa976001600160801b036001600160401b03985191519351955116906143d7565b6040516105898582018093604080916001600160801b038151168452602081015160208501520151910152565b606081526105986080826113c4565b5190206105aa86858a51015116612f4e565b808214613339565b6040516105df8382018093604080916001600160801b038151168452602081015160208501520151910152565b606081526105ee6080826113c4565b519020940151015116612f4e565b68ff0000000000000000191617600455005b9195509160803d60801161069b575b61062781846113c4565b8201916080818403126101f35760206105fc956105b26001600160401b039461055c680100000000000000009b6001600160801b0386976106856105aa9b60408051936106738561133d565b61067d8382611c37565b855201611c37565b888201529d50509750509694505095505061051f565b503d61061d565b6040513d5f823e3d90fd5b634e487b7160e01b5f52604160045260245ffd5b6106ca906131ed565b83610254565b7f928b1233000000000000000000000000000000000000000000000000000000005f5260045ffd5b346101f35761072c61070936611181565b90610727610722825f525f602052600160405f20015490565b6131ed565b6138cd565b005b346101f3575f3660031901126101f3576107c46040516107588161075181611273565b03826113c4565b6040516107648161133d565b60ff600254818116835260081c1660208201526040516107838161133d565b6001600160401b03600354818116835260401c16602082015260ff6004546107fd828260481c16936107de6040519889986101208a526101208a019061124f565b96602089019060ff60208092828151168552015116910152565b60608701906001600160401b0360208092828151168552015116910152565b63ffffffff811660a086015263ffffffff8160201c1660c086015260401c16151560e084015261082c816113e5565b6101008301520390f35b346101f3575f3660031901126101f35760206040516107088152f35b346101f35760203660031901126101f3576004356001600160401b0381116101f357366023820112156101f35780600401356001600160401b0381116101f35760248201913660248360051b830101116101f3576020926040516108b685826113c4565b5f815284810191601f1986013684376108ce85611630565b936108dc60405195866113c4565b858552601f196108eb87611630565b01875f5b828110610a08575050505f5b868110156109ab576001906109875f808b8861095461092260248860051b8b01018b61312d565b90938d6040519483869484860198893784019083820190898252519283915e010185815203601f1981018352826113c4565b5190305af43d156109a3573d9061096a826113ef565b9161097860405193846113c4565b82523d5f8d84013e5b30614c67565b6109918289612fc6565b5261099c8188612fc6565b50016108fb565b606090610981565b604080518981528751818b018190525f92600582901b83018101918a8d01918d9085015b8287106109dc5785850386f35b9091929382806109f8600193603f198a8203018652885161124f565b96019201960195929190926109cf565b606088820183015281016108ef565b346101f35760203660031901126101f3576004356001600160401b0381116101f357806004019061014060031982360301126101f35760ff60045460401c166106d0577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610b09575b50610aa8604482018361315f565b91610ab6606482018561315f565b9390916101048101359060028210156101f357602096610b0196610ade61012484018361315f565b96909560405198610aef8c8b6113c4565b5f8a52608460a4870196013594613a71565b604051908152f35b610b12906131ed565b82610a9a565b346101f3575f3660031901126101f35760206040515f8152f35b346101f35760203660031901126101f3576004356001600160401b0381116101f357806004019061016060031982360301126101f35760ff60045460401c166106d0577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610c5f575b506101448101610bc5818461312d565b905015610c3757610bd9604483018461315f565b9091610be8606485018661315f565b9590946101048101359160028310156101f357602097610b0197610c2096610c27610c1761012487018661315f565b9990988661312d565b369161140a565b98608460a4870196013594613a71565b7f1208b21b000000000000000000000000000000000000000000000000000000005f5260045ffd5b610c68906131ed565b82610bb5565b346101f357610c7c36611181565b905f525f6020526001600160a01b0360405f2091165f52602052602060ff60405f2054166040519015158152f35b346101f357610cb836611100565b505060ff60045460401c166106d0577f0000000000000000000000000000000000000000000000000000000000000000805f525f60205260405f205f805260205260ff60405f20541615610d2f575b7fda81d7c2000000000000000000000000000000000000000000000000000000005f5260045ffd5b610d38906131ed565b80610d07565b346101f3575f3660031901126101f35760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101f3575f3660031901126101f35760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101f3575f3660031901126101f35760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b346101f35760203660031901126101f357610e3a610e486101ef610e23600435612fee565b9193906040519586956060875260608701906111a7565b9085820360208701526111e0565b908382036040850152611213565b346101f357610e6436611181565b336001600160a01b03821603610e7d5761072c916138cd565b7f6697b232000000000000000000000000000000000000000000000000000000005f5260045ffd5b346101f35761072c610eb636611181565b90610ecf610722825f525f602052600160405f20015490565b613840565b346101f35760203660031901126101f3576020610b016004355f525f602052600160405f20015490565b346101f35760203660031901126101f3576004356001600160401b03811681036101f357610b01602091612f4e565b346101f357610f3b36611100565b9060ff60045460401c166106d057602091610f94917f0000000000000000000000000000000000000000000000000000000000000000805f525f855260405f205f8052855260ff60405f20541615610fa5575b506124e5565b60405190610fa18161114f565b8152f35b610fae906131ed565b84610f8e565b346101f3575f3660031901126101f35760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101f35760203660031901126101f35760206110156004356131bc565b6040519015158152f35b346101f3575f3660031901126101f35760206040516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101f35760203660031901126101f357600435907fffffffff0000000000000000000000000000000000000000000000000000000082168092036101f357817f7965db0b00000000000000000000000000000000000000000000000000000000602093149081156110d6575b5015158152f35b7f01ffc9a700000000000000000000000000000000000000000000000000000000915014836110cf565b9060206003198301126101f3576004356001600160401b0381116101f357826023820112156101f3578060040135926001600160401b0384116101f357602484830101116101f3576024019190565b6003111561115957565b634e487b7160e01b5f52602160045260245ffd5b35906001600160401b03821682036101f357565b60409060031901126101f357600435906024356001600160a01b03811681036101f35790565b90602080835192838152019201905f5b8181106111c45750505090565b825163ffffffff168452602093840193909201916001016111b7565b90602080835192838152019201905f5b8181106111fd5750505090565b82518452602093840193909201916001016111f0565b90602080835192838152019201905f5b8181106112305750505090565b82516001600160401b0316845260209384019390920191600101611223565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b905f91600154908160011c91600181168015611333575b60208410811461131f5783835290811561130357506001146112aa575050565b60015f9081529293509091907fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b8383106112e9575060209250010190565b6001816020929493945483858701015201910191906112d8565b9050602093945060ff929192191683830152151560051b010190565b634e487b7160e01b5f52602260045260245ffd5b92607f169261128a565b604081019081106001600160401b038211176106ad57604052565b606081019081106001600160401b038211176106ad57604052565b60e081019081106001600160401b038211176106ad57604052565b608081019081106001600160401b038211176106ad57604052565b60c081019081106001600160401b038211176106ad57604052565b90601f801991011681019081106001600160401b038211176106ad57604052565b6002111561115957565b6001600160401b0381116106ad57601f01601f191660200190565b929192611416826113ef565b9161142460405193846113c4565b8294818452818301116101f3578281602093845f960137010152565b9080601f830112156101f35781602061145b9335910161140a565b90565b359060ff821682036101f357565b91908260409103126101f3576040516114848161133d565b602061149d8183956114958161116d565b85520161116d565b910152565b359063ffffffff821682036101f357565b359081151582036101f357565b808203929161012084126101f357604051916114db83611373565b82948135906001600160401b0382116101f3576114fc846040938501611440565b8552601f1901126101f3576115476101009260405161151a8161133d565b6115266020850161145e565b81526115346040850161145e565b602082015260208601526060830161146c565b604084015261155860a082016114a2565b606084015261156960c082016114a2565b608084015261157a60e082016114b3565b60a084015201359060028210156101f35760c00152565b35906001600160801b03821682036101f357565b91908260609103126101f3576040516115bd81611358565b60408082946115cb81611591565b8452602081013560208501520135910152565b8092910391606083126101f3576040516115f78161133d565b6040819483358352601f1901126101f35760209060408051936116198561133d565b6116248482016114a2565b85520135828401520152565b6001600160401b0381116106ad5760051b60200190565b91906080838203126101f357604051906116608261138e565b819380356001600160401b0381116101f357606092611680918301611440565b8352602081013560208401526116986040820161116d565b60408401520135908160070b82036101f35760600152565b9190916080818403126101f357604051906116ca8261138e565b819381356001600160401b0381116101f357820181601f820112156101f35780356116f481611630565b9161170260405193846113c4565b81835260208084019260051b820101918483116101f35760208201905b83821061177057505050508352611738602083016114b3565b60208401526040820135906001600160401b0382116101f357826117656060949261149d94869401611647565b60408601520161116d565b81356001600160401b0381116101f35760209161179288848094880101611647565b81520191019061171f565b919060a0838203126101f357604051906117b68261138e565b819380356001600160401b0381116101f35781016040818403126101f357604051906117e18261133d565b80356001600160401b0381116101f35781016102c0818603126101f3576040519061026082018281106001600160401b038211176106ad57604052611826868261146c565b825260408101356001600160401b0381116101f35786611847918301611440565b60208301526118586060820161116d565b604083015261186960808201611591565b606083015261187a60a082016114b3565b608083015261188c8660c083016115de565b60a083015261189e61012082016114b3565b60c083015261014081013560e08301526118bb61016082016114b3565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a083015261190a61022082016114b3565b6101c08301526102408101356101e083015261192961026082016114b3565b6102008301526102808101356102208301526102a0810135906001600160401b0382116101f35761195c91879101611440565b61024082015282526020810135906001600160401b0382116101f3570160c0818503126101f357604051906119908261138e565b6119998161116d565b82526119a7602082016114a2565b60208301526119b985604083016115de565b604083015260a0810135906001600160401b0382116101f3570184601f820112156101f3578035906119ea82611630565b916119f860405193846113c4565b80835260208084019160051b830101918783116101f35760208101915b838310611a83575050505060608201526020820152835260208101356001600160401b0381116101f35782611a4b9183016116b0565b6020840152611a5d826040830161146c565b60408401526080810135916001600160401b0383116101f35760609261149d92016116b0565b82356001600160401b0381116101f3578201906040828b03601f1901126101f35760405191611ab18361133d565b602081013560048110156101f357835260408101356001600160401b0381116101f3576020910101906080828c03126101f35760405192611af18461138e565b82356001600160401b0381116101f3578c611b0d918501611440565b8452611b1b60208401611591565b6020850152611b2c604084016114b3565b60408501526060830135936001600160401b0385116101f357611b548d602096879601611440565b606082015283820152815201920191611a15565b9080601f830112156101f357604080519290611b8490846113c4565b8290604081019283116101f357905b828210611ba05750505090565b8135815260209182019101611b93565b9080601f830112156101f3578135611bc781611630565b92611bd560405194856113c4565b81845260208085019260051b8201019283116101f357602001905b828210611bfd5750505090565b60208091611c0a846114a2565b815201910190611bf0565b519060ff821682036101f357565b51906001600160401b03821682036101f357565b91908260409103126101f357604051611c4f8161133d565b602061149d818395611c6081611c23565b855201611c23565b519063ffffffff821682036101f357565b519081151582036101f357565b51906001600160801b03821682036101f357565b91908260609103126101f357604051611cb281611358565b6040808294611cc081611c86565b8452602081015160208501520151910152565b6020818303126101f3578051906001600160401b0382116101f35701610180818303126101f35760405191611d07836113a9565b81516001600160401b0381116101f3578201918282039261012084126101f35760405193611d3485611373565b81516001600160401b0381116101f35782019084601f830112156101f3578151611d5d816113ef565b90611d6b60405192836113c4565b80825286602082860101116101f3576020815f9282604097018386015e830101528652601f1901126101f35761010090604051611da78161133d565b611db360208301611c15565b8152611dc160408301611c15565b60208201526020860152611dd88460608301611c37565b6040860152611de960a08201611c68565b6060860152611dfa60c08201611c68565b6080860152611e0b60e08201611c79565b60a086015201519060028210156101f357836101409260c0611e799601528552611e388360208301611c9a565b6020860152611e4a8360808301611c9a565b6040860152611e5b60e08201611c86565b6060860152611e6e836101008301611c37565b608086015201611c37565b60a082015290565b9061010060c0611e9c8451610120855261012085019061124f565b602080860151805160ff9081168784015291015116604085015293611ede604082015160608601906001600160401b0360208092828151168552015116910152565b63ffffffff60608201511660a085015263ffffffff6080820151168285015260a0810151151560e0850152015191611f15836113e5565b015290565b90606080611f31845160808552608085019061124f565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b906080810182519060808352815180915260a0830190602060a08260051b8601019301915f905b828210611fc357505050506001600160401b036060611fb9819360208701511515602087015260408701518682036040880152611f1a565b9401511691015290565b90919293602080611fe0600193609f198a82030186528851611f1a565b960192019201909291611f81565b91909180519260a08152602061214e8551604060a085015261202960e0850182516001600160401b0360208092828151168552015116910152565b610240612047848301516102c06101208801526103a087019061124f565b60408301516001600160401b031661014087015260608301516001600160801b03166101608701526080830151151561018087015260a083015180516101a0880152602090810151805163ffffffff166101c089015201516101e08701529160c0810151151561020087015260e08101516102208701526101008101511515828701526101208101516102608701526101408101516102808701526101608101516102a08701526101808101516102c08701526101a08101516102e08701526101c081015115156103008701526101e08101516103208701526102008101511515610340870152610220810151610360870152015160df198583030161038086015261124f565b94015193609f198282030160c0830152606060c08201956001600160401b03815116835263ffffffff60208201511660208401526121ae6040820151604085019060208060409280518552015163ffffffff815116828501520151910152565b01519460c060a0830152855180915260e0820190602060e08260051b8501019701925f905b82821061223357505050505060606121fb61145b949560208501518482036020860152611f5a565b92612223604082015160408501906001600160401b0360208092828151168552015116910152565b0151906080818403910152611f5a565b909192939760df1982820301855288519081516004811015611159576122b1826020600195819594829552015190604084820152606061227f83516080604085015260c084019061124f565b926001600160801b0386820151168284015260408101511515608084015201519060a0603f198285030191015261124f565b9a019501939201906121d3565b905f905b600882106122cf57505050565b60208060019285518152019301910190916122c2565b905f905b600282106122f657505050565b60208060019285518152019301910190916122e9565b90602080835192838152019201905f5b8181106123295750505090565b8251151584526020938401939092019160010161231c565b91906020835261241c6123aa61236583516103006020880152610320870190611e81565b61239460208501516040880190604080916001600160801b038151168452602081015160208501520151910152565b6040840151868203601f190160a0880152611fee565b6001600160801b0360608401511660c08601526123cf608084015160e08701906122be565b6123e260a08401516101e08701906122e5565b6123f560c08401516102208701906122e5565b60e083015161ffff16610260860152610100830151858203601f19016102808701526111a7565b9261012082015193601f19828203016102a0830152602080865192838152019501905f5b8181106124cf5750505061014082015193601f19828203016102c0830152602080865192838152019501905f5b8181106124b05750505061018061249a61145b9495610160850151601f19858303016102e08601526111a7565b92015190610300601f198285030191015261230c565b82516001600160401b031687526020968701969092019160010161246d565b8251875260209687019690920191600101612440565b908101906020818303126101f3578035906001600160401b0382116101f3570191610300838303126101f3576040516101a081018181106001600160401b038211176106ad5760405283356001600160401b0381116101f3578361254a9186016114c0565b815261255983602086016115a5565b602082015260808401356001600160401b0381116101f3578361257d91860161179d565b916040820192835261259160a08601611591565b60608301528360df860112156101f3576040516125b0610100826113c4565b806101c08701918683116101f35790869160c08901905b848210612f3b57505060808501526125de91611b68565b60a08301526125f1846102008701611b68565b60c083015261024085013561ffff811681036101f35760e083019081526102608601356001600160401b0381116101f3578561262e918801611bb0565b6101008401526102808601356001600160401b0381116101f357860185601f820112156101f35780359061266182611630565b9161266f60405193846113c4565b80835260208084019160051b830101918883116101f357602001905b828210612f2b575050506101208401526102a08601356001600160401b0381116101f35786019480601f870112156101f3578535956126c987611630565b966126d760405198896113c4565b80885260208089019160051b830101918383116101f357602001905b828210612f135750505061014084019586526102c08701356001600160401b0381116101f35781612725918901611bb0565b9661016085019788526102e0810135906001600160401b0382116101f357019080601f830112156101f357813561275b81611630565b9261276960405194856113c4565b81845260208085019260051b8201019283116101f357602001905b828210612efb5750505061018084015261279d83613246565b9691929097835f14612e68576040517f155bbe880000000000000000000000000000000000000000000000000000000081525f81806127df8a60048301612341565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156106a2575f91612e46575b50965b61283788516001600160801b0360608b015116906143d7565b612899602089015160405161286d602082018093604080916001600160801b038151168452602081015160208501520151910152565b6060815261287c6080826113c4565b5190206105aa6001600160401b03602060808d0151015116612f4e565b6128a288613370565b9415612c0f575061ffff610100870151519251168092149283612bff575b83612bf3575b5082612be7575b5081612bd7575b5015612baf57855f5260066020526001600160a01b0360405f2054168015612b9c57612904909594939295613950565b61290e81886139f7565b90965f96879687959186805b6101008a0151518910156129e3578d906129398a6101808d0151612fc6565b51156129d157988a9b9c899796959493928c9b63ffffffff612962859f61010087910151612fc6565b511680936129b9575b905060019f8b829d61297c94614858565b919390939c61012001519061299091612fc6565b51149061299c91613478565b6129a59161341e565b986001905b0197999890919293949661291a565b63ffffffff6129ca9216811161343e565b5f8261296b565b979594939291999a98600191506129aa565b50945096955097509792509794506001600160401b038116600381029080820460031490151715612b88576001600160401b038316916801fffffffffffffffe8460011b169280840460021490151715612b8857612a43928492116134b2565b612a4c826145f8565b915b612b77575b505050612a5f8261114f565b81612b2e576001600160401b036020604060a0840193845183810151600354918683861c1687831611612ae5575b5050500151604051612abf8382018093604080916001600160801b038151168452602081015160208501520151910152565b60608152612ace6080826113c4565b51902092510151165f52600560205260405f205590565b6fffffffffffffffff0000000000000000877fffffffffffffffffffffffffffffffff0000000000000000000000000000000092511692861b16921617176003555f8080612a8d565b50612b388161114f565b60018103612b60576801000000000000000068ff000000000000000019600454161760045590565b612b698161114f565b6002810361145b5750600290565b612b80926134f6565b5f8080612a53565b634e487b7160e01b5f52601160045260245ffd5b86633a517eed60e21b5f5260045260245ffd5b7f17e37b5c000000000000000000000000000000000000000000000000000000005f5260045ffd5b905061018084015151145f6128d4565b5151811491505f6128cd565b5151821492505f6128c6565b61012087015151831493506128c0565b98925f98949198969592965060205f9a510151519687519961ffff610100890151519251168092149283612e36575b83612e2a575b5082612e1e575b5081612e0e575b5015612baf575f989698959493955b888110612de357505f985f935f975f985b610100880151518a1015612d7757612c8f8a6101808a0151612fc6565b5115612d685763ffffffff612ca98b6101008b0151612fc6565b5116968c881015612d3c578a9b9c9d612cf5898b612cee8e9f8f9e9f6020612ce3866101209360019c82612d159d612d24575b5050612fc6565b510151930151612fc6565b5114613478565b6001600160401b036040612d0b859b809d612fc6565b510151169061341e565b9a5b019897969b9a999b612c72565b63ffffffff612d359216811161343e565b5f82612cdc565b877fc9d6366c000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b9697986001909c9a9b9c612d17565b50965096509692509650966001600160401b038116600381029080820460031490151715612b88576001600160401b038316916801fffffffffffffffe8460011b169280840460021490151715612b8857612dd4928492116134b2565b612ddd826145f8565b91612a4e565b9896612e026001916001600160401b036040612d0b8e8b9c9a9b612fc6565b97999496959401612c61565b905061018086015151145f612c52565b5151811491505f612c4b565b5151821492505f612c44565b6101208901515183149350612c3e565b612e6291503d805f833e612e5a81836113c4565b810190611cd3565b5f61281b565b6040517fb4e415790000000000000000000000000000000000000000000000000000000081525f8180612e9e8a60048301612341565b03816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9081156106a2575f91612ee1575b509661281e565b612ef591503d805f833e612e5a81836113c4565b5f612eda565b60208091612f08846114b3565b815201910190612784565b60208091612f208461116d565b8152019101906126f3565b813581526020918201910161268b565b81358152889350602091820191016125c7565b6001600160401b03165f52600560205260405f20548015612f6c5790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b90612f9e82611630565b612fab60405191826113c4565b8281528092612fbc601f1991611630565b0190602036910137565b8051821015612fda5760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b90815f5260066020526001600160a01b0360405f20541680156130e35761301490613950565b9161301e836139b0565b156130e3576130308361ffff926139f7565b9050169061303d82612f94565b9261304783612f94565b9061305184612f94565b5f5b61ffff811690868210156130da57602c820290828204602c1483151715612b885781600e0180600e11612b885782860192602e84015160e01c613096868d612fc6565b52601281018211612b8857601a90603285015160c01c6130b68789612fc6565b520110612b885761ffff926130d2603a60019401519188612fc6565b520116613053565b50509291935050565b5090506020604051916130f682846113c4565b5f83525f3681376040519261310b83856113c4565b5f84525f3681376040519261312081856113c4565b5f8452505f368137929190565b903590601e19813603018212156101f357018035906001600160401b0382116101f3576020019181360383136101f357565b903590601e19813603018212156101f357018035906001600160401b0382116101f357602001918160051b360383136101f357565b9061319e826113ef565b6131ab60405191826113c4565b8281528092612fbc601f19916113ef565b5f5260066020526001600160a01b0360405f20541680156131e8576131e361145b91613950565b6139b0565b505f90565b805f525f60205260405f206001600160a01b0333165f5260205260ff60405f205416156132175750565b7fe2517d3f000000000000000000000000000000000000000000000000000000005f523360045260245260445ffd5b905f906040830191825192835160016001600160401b036020604082818084610140895101519d01510151965101511698015101511601946001600160401b038611612b8857613295876131bc565b938461331c57506132ab6020845101518861425a565b6001600160401b036001965b161461330857821580806132ff575b156132dd5750505160606020820151910152929190565b6132e9575b5050929190565b60606132f8925101519061425a565b5f806132e2565b508682146132c6565b506060613313614200565b91510152929190565b956001600160401b039061332e614200565b6020865101526132b7565b15613342575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6001600160401b03602060a08301510151165f52600560205260405f205480155f1461339c5750505f90565b604082019081516040516133d1602082018093604080916001600160801b038151168452602081015160208501520151910152565b606081526133e06080826113c4565b51902014918215926133fe575b5050156133f957600190565b600290565b6001600160801b0391925060208291015151169151511611155f806133ed565b906001600160401b03809116911601906001600160401b038211612b8857565b156134465750565b63ffffffff907fea334314000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156134805750565b63ffffffff907fc3271bcd000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b156134bb575050565b906001600160401b0380927f6e3083c3000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b9290613501846131bc565b61383a5760400151602001515180519161ffff8311612baf57602c8302838104602c1484151715612b8857600e0180600e11612b885761354090613194565b9260208401946056865361564160218601536256414c60228601536356414c3160238601538060381c60248601538060301c60258601538060281c60268601538060201c60278601538060181c60288601538060101c60298601538060081c602a860153602b8501538060081c602c850153602d8401535f94600e5b83518710156136b7576135cf8785612fc6565b518186019060ff8960181c16602083015361ffff8960101c16602183015362ffffff8960081c16602283015363ffffffff89166023830153600483018311612b88576001600160401b03604082015160ff8160381c16602485015361ffff8160301c16602585015362ffffff8160281c16602685015363ffffffff8160201c16602785015364ffffffffff8160181c16602885015365ffffffffffff8160101c16602985015366ffffffffffffff8160081c16602a85015316602b830153600c83018311612b88576020602c910151910152602c8101809111612b88576001909601956135bc565b509391509350604051916136ed60218460208101945f8652845180918484015e81015f838201520301601f1981018552846113c4565b61600083511161380e5750906001600160a01b03916137a5602a7fffff000000000000000000000000000000000000000000000000000000000000845160f01b169360405193849160208301967f6100000000000000000000000000000000000000000000000000000000000000885260218401527f80600a3d393df30000000000000000000000000000000000000000000000000060238401525180918484015e81015f838201520301601f1981018352826113c4565b51905ff0169081156137e6575f52600660205260405f20907fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055565b7ffbad885d000000000000000000000000000000000000000000000000000000005f5260045ffd5b517fc5cd3084000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b50915050565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f205416155f146138c757805f525f60205260405f206001600160a01b0383165f5260205260405f20600160ff198254161790556001600160a01b03339216907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d5f80a4600190565b50505f90565b805f525f60205260405f206001600160a01b0383165f5260205260ff60405f2054165f146138c757805f525f60205260405f206001600160a01b0383165f5260205260405f2060ff1981541690556001600160a01b03339216907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b90813b600181111561397b575f198101908111612b8857600161397282613194565b9360208501903c565b6001600160a01b03837fd8415944000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b805190600e82106138c7576356414c31602082015160e01c036138c757602c015160f01c602c810290808204602c1490151715612b8857600e019081600e11612b88571490565b91600e825110613a5e576356414c31602083015160e01c03613a5e57602482015160c01c92602c83015160f01c9251602c8402848104602c1485151715612b8857600e019081600e11612b885703613a4c5750565b634724a0fd60e01b5f5260045260245ffd5b82634724a0fd60e01b5f5260045260245ffd5b999493979198909695995f96613a86816113e5565b806141bf5750891515806141b3575b1561417c576020019460208b613b02613ab8613ab08a614933565b9236906115a5565b916105aa60405185810190613aea8287604080916001600160801b038151168452602081015160208501520151910152565b60608152613af96080826113c4565b51902091612f4e565b015180860361414c57505f905f5b8b8a81831061405f575b505050501561401e5750506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690813b156101f3579187926040519384927f87d3a9b100000000000000000000000000000000000000000000000000000000845260648401906004850152606060248501525260848201908960051b9860848a850101928b8a915f603e198d360301925b8110613fb05750505050600319848403016044850152808352602083019260208260051b82010193835f92601e1982360301905b858510613e41575050505050505091815f81819503925af180156106a257613e2c575b5060018511613c37575b50505050506001600160801b03613c31633b9aca009236906115a5565b51160490565b613c48909691929496959395614933565b948335946001600160801b038616809603613e2857613c6688611630565b97613c74604051998a6113c4565b88526020880191810190368211613e245780979597925b828410613d8957505050506001600160401b03829316925b8651811015613d6c57613cb68188612fc6565b519085604051602081019087825260408082015260a081019480519560406060840152865180915260c08301602060c08360051b8601019801918a5b818110613d35575050505081613d1d6001976020613d2b940151605f1984830301608085015261124f565b03601f1981018352826113c4565b5190205d01613ca3565b919394965091949697602080613d5760019360bf198b8203018852895161124f565b970194019101918c9694939298979598613cf2565b5094505050506001600160801b03613c31633b9aca005f80613c14565b83989698356001600160401b038111613e20578201604081360312613e205760405190613db58261133d565b80356001600160401b038111613e1c57810136601f82011215613e1c57613de3903690602081359101614969565b825260208101356001600160401b038111613e1c5791613e0a602094928594369101611440565b83820152815201930192979597613c8b565b8880fd5b8680fd5b8480fd5b8380fd5b613e399192505f906113c4565b5f905f613c0a565b9193959750919395601f198282030185528735838112156101f357840190613e6d602082019280614a84565b8091936020845252604082019060408160051b8401019380935f915b838310613eb057505050505050602080600192990195019501929091899796949592613be7565b909192939495603f19838203018652613ec98783614acc565b803560028110156101f357613edd816113e5565b8252613f00613eef6020830183614ab8565b606060208501526060840190614ae0565b906040810135609e19823603018112156101f3576001936020938493613fa29301916040818303910152613f94613f76613f4b613f3d85806149f4565b60a0865260a08601916149d4565b613f568786016114b3565b151587850152613f696040860186614ab8565b8482036040860152614ae0565b92613f83606082016114b3565b151560608401526080810190614ab8565b906080818403910152614ae0565b980196019493019190613e89565b919394965091946083198982030183528535848112156101f357602061400c6001938f83940190613fff613ff5613fe78480614a84565b604085526040850191614a25565b92858101906149f4565b91858185039101526149d4565b9701930191019088969493918e613bb3565b61405b6040519283927ffef760c7000000000000000000000000000000000000000000000000000000008452602060048501526024840191614a25565b0390fd5b61407d61407784614095946140849498969798614947565b8061315f565b3691614969565b61408f368789614969565b90614cf3565b156141425750610c206140ac6140b6928d8c614947565b602081019061312d565b80518251808214918261412c575b5050156140d857505060015f808b8a613b1a565b9061405b61411a926040519384937f5f1ca38100000000000000000000000000000000000000000000000000000000855260406004860152604485019061124f565b8381036003190160248501529061124f565b9091506020830120906020840120145f806140c4565b9190600101613b10565b85907f4b2dfe98000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b897fb0369c31000000000000000000000000000000000000000000000000000000005f52600452600160245261ffff60445260645ffd5b5061ffff8a1115613a95565b8760ff6024926141ce816113e5565b6141d7816113e5565b7f112d89cc00000000000000000000000000000000000000000000000000000000835216600452fd5b6040519061420d8261138e565b606082525f60208301526040516142238161138e565b606081525f60208201525f60408201525f606082015260408301525f606083015281604051906142546020836113c4565b5f825252565b815151156143c557604051917fc87f1f69000000000000000000000000000000000000000000000000000000008352602060048401528260a4810182519060806024840152815180915260c4830190602060c48260051b8601019301915f905b8282106143965750505050826001600160401b0360606142f88594602080980151151560448701526040850151602319878303016064880152611f1a565b92015116608483015203817317435cce3d1b4fa2e5f8a08ed921d57c6762a1805af49182156106a2575f92614362575b50808203614334575050565b7f769a20cb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b9091506020813d60201161438e575b8161437e602093836113c4565b810103126101f35751905f614328565b3d9150614371565b919360019193955060206143b5819260c3198c82030186528851611f1a565b96019201920187949391926142ba565b633a517eed60e21b5f5260045260245ffd5b906001600160801b03633b9aca009116044281116145c957804203428111612b88576107081061459a57508051805160208201207f00000000000000000000000000000000000000000000000000000000000000000361455a5750602081015160ff815116906002549160ff8316928382149283614541575b6020015160ff1692156144f9575050505063ffffffff606082015116906004549163ffffffff83168082036144cb57505063ffffffff608081920151169160201c1680820361449d575050565b7f1eb5c1eb000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f73b1bde3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b6084945060ff90604051947fd382033a000000000000000000000000000000000000000000000000000000008652600486015260081c16602484015260448301526064820152fd5b6020810151600883901c60ff9081169116149350614450565b61405b906040519182917ff6b6676b0000000000000000000000000000000000000000000000000000000083526040600484015261411a60448401611273565b7f12e74add000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b7f20daedb6000000000000000000000000000000000000000000000000000000005f524260045260245260445ffd5b60408101515160208101519081516001600160401b031690602083015163ffffffff1692604001519081519160200151805163ffffffff1690602001519151602001519260405194614649866113a9565b85526020850195865260408501908152606085019182526080850192835260a0850193845260e087015161ffff169560808801519760a08101519060c081015161012082015161014083015191610160840151936101800151946040519d8e809e7f4cc22bb700000000000000000000000000000000000000000000000000000000825260048201526024016146de916122be565b6101248d016146ec916122e5565b6101648c016146fa916122e5565b6101a48b0161024090526102448b01614712916111e0565b8a8103600319016101c48c015261472891611213565b898103600319016101e48b015261473e916111a7565b888103600319016102048a01526147549161230c565b9560031988880301610224890152516001600160401b03168652516001600160401b031660208601525160408501525163ffffffff166060840152516080830152519060a0810160c0905260c0016147ab9161124f565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031691815a6020945f91f19081156106a2575f9161481e575b50156147f657565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90506020813d602011614850575b81614839602093836113c4565b810103126101f35761484a90611c79565b5f6147ee565b3d915061482c565b92949091945b61ffff861681106148a0575b63ffffffff85857ff81f5140000000000000000000000000000000000000000000000000000000005f526004521660245260445ffd5b602c8102818104602c1482151715612b885780600e019081600e11612b885780850190602e82015160e01c9163ffffffff8916938484146148f857505050116148f3575f198114612b885760010161485e565b61486a565b9597509550975050509350601285018111612b8857601a603283015160c01c950110612b8857603a01519160018101809111612b8857929190565b356001600160401b03811681036101f35790565b9190811015612fda5760051b81013590603e19813603018212156101f3570190565b92919061497581611630565b9361498360405195866113c4565b602085838152019160051b8101918383116101f35781905b8382106149a9575050505050565b81356001600160401b0381116101f3576020916149c98784938701611440565b81520191019061499b565b908060209392818452848401375f828201840152601f01601f1916010190565b9035601e19823603018112156101f35701602081359101916001600160401b0382116101f35781360383136101f357565b90602083828152019260208260051b82010193835f925b848410614a4c5750505050505090565b909192939495602080614a74600193601f19868203018852614a6e8b886149f4565b906149d4565b9801940194019294939190614a3c565b9035601e19823603018112156101f35701602081359101916001600160401b0382116101f3578160051b360383136101f357565b9035607e19823603018112156101f3570190565b9035605e19823603018112156101f3570190565b614b19614afe614af083806149f4565b6080865260808601916149d4565b614b0b60208401846149f4565b9085830360208701526149d4565b614b266040830183614ab8565b908381036040850152813560038110156101f357614b438161114f565b8152602082013560038110156101f357614b5c8161114f565b602082015260408201359060038210156101f3576080614b98614bb39484614b89614ba89699989961114f565b604085015260608101906149f4565b91909281606082015201916149d4565b926060810190614a84565b90916060818503910152808352602083019060208160051b85010193835f915b838310614be35750505050505090565b909192939495601f19828203018652614bfc8784614acc565b80359160038310156101f357614c59602092839285614c1c60019761114f565b8152614c4b614c40614c30868501856149f4565b60608886015260608501916149d4565b9260408101906149f4565b9160408185039101526149d4565b980196019493019190614bd3565b90614ca45750805115614c7c57602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b81511580614cea575b614cb5575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15614cad565b9081518151036138c7575f5b8251811015614d5d57614d128184612fc6565b5151614d1e8284612fc6565b515103614d5657614d2f8184612fc6565b5160208151910120614d418284612fc6565b516020815191012003614d5657600101614cff565b5050505f90565b50505060019056fea164736f6c634300081c000a2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0db10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
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
