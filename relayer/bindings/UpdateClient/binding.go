// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractUpdateClient

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

// IICS07TendermintMsgsBlockCommit is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsBlockCommit struct {
	Height     uint64
	Round      uint32
	BlockId    IICS07TendermintMsgsBlockId
	CommitSigs []IICS07TendermintMsgsCommitSig
}

// IICS07TendermintMsgsBlockHeader is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsBlockHeader struct {
	Version            IICS07TendermintMsgsVersion
	ChainId            string
	Height             uint64
	Time               *big.Int
	HasLastBlockId     bool
	LastBlockId        IICS07TendermintMsgsBlockId
	HasLastCommitHash  bool
	LastCommitHash     [32]byte
	HasDataHash        bool
	DataHash           [32]byte
	ValidatorsHash     [32]byte
	NextValidatorsHash [32]byte
	ConsensusHash      [32]byte
	AppHash            [32]byte
	HasLastResultsHash bool
	LastResultsHash    [32]byte
	HasEvidenceHash    bool
	EvidenceHash       [32]byte
	ProposerAddress    []byte
}

// IICS07TendermintMsgsBlockId is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsBlockId struct {
	HashData      [32]byte
	PartSetHeader IICS07TendermintMsgsPartSetHeader
}

// IICS07TendermintMsgsClientState is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsClientState struct {
	ChainId         string
	TrustLevel      IICS07TendermintMsgsTrustThreshold
	LatestHeight    IICS02ClientMsgsHeight
	TrustingPeriod  uint32
	UnbondingPeriod uint32
	IsFrozen        bool
	ZkAlgorithm     uint8
}

// IICS07TendermintMsgsCommitSig is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsCommitSig struct {
	Flag uint8
	Data IICS07TendermintMsgsCommitSigData
}

// IICS07TendermintMsgsCommitSigData is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsCommitSigData struct {
	ValidatorAddress []byte
	Timestamp        *big.Int
	HasSignature     bool
	Signature        []byte
}

// IICS07TendermintMsgsConsensusState is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsConsensusState struct {
	Timestamp          *big.Int
	Root               [32]byte
	NextValidatorsHash [32]byte
}

// IICS07TendermintMsgsHeader is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsHeader struct {
	SignedHeader            IICS07TendermintMsgsSignedHeader
	ValidatorSet            IICS07TendermintMsgsValidatorSet
	TrustedHeight           IICS02ClientMsgsHeight
	TrustedNextValidatorSet IICS07TendermintMsgsValidatorSet
}

// IICS07TendermintMsgsPartSetHeader is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsPartSetHeader struct {
	Total    uint32
	HashData [32]byte
}

// IICS07TendermintMsgsSignedHeader is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsSignedHeader struct {
	Header IICS07TendermintMsgsBlockHeader
	Commit IICS07TendermintMsgsBlockCommit
}

// IICS07TendermintMsgsTrustThreshold is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsTrustThreshold struct {
	Numerator   uint8
	Denominator uint8
}

// IICS07TendermintMsgsValidatorInfo is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsValidatorInfo struct {
	ValAddress       []byte
	PubKey           [32]byte
	VotingPower      uint64
	ProposerPriority int64
}

// IICS07TendermintMsgsValidatorSet is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsValidatorSet struct {
	Validators       []IICS07TendermintMsgsValidatorInfo
	HasProposer      bool
	Proposer         IICS07TendermintMsgsValidatorInfo
	TotalVotingPower uint64
}

// IICS07TendermintMsgsVersion is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsVersion struct {
	BlockVersion uint64
	AppVersion   uint64
}

// IUpdateClientMsgsMsgUpdateClient is an auto generated low-level Go binding around an user-defined struct.
type IUpdateClientMsgsMsgUpdateClient struct {
	ClientState           IICS07TendermintMsgsClientState
	TrustedConsensusState IICS07TendermintMsgsConsensusState
	ProposedHeader        IICS07TendermintMsgsHeader
	Time                  *big.Int
	Proof                 [8]*big.Int
	Commitments           [2]*big.Int
	CommitmentPok         [2]*big.Int
	Bucket                uint16
	SignerIndices         []uint32
	SignerPubkeys         [][32]byte
	TimestampSeconds      []uint64
	TimestampNanos        []uint32
	Active                []bool
}

// IUpdateClientMsgsUpdateClientOutput is an auto generated low-level Go binding around an user-defined struct.
type IUpdateClientMsgsUpdateClientOutput struct {
	ClientState           IICS07TendermintMsgsClientState
	TrustedConsensusState IICS07TendermintMsgsConsensusState
	NewConsensusState     IICS07TendermintMsgsConsensusState
	Time                  *big.Int
	TrustedHeight         IICS02ClientMsgsHeight
	NewHeight             IICS02ClientMsgsHeight
}

// ContractUpdateClientMetaData contains all meta data concerning the ContractUpdateClient contract.
var ContractUpdateClientMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.MsgUpdateClient\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"timestampSeconds\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"timestampNanos\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.UpdateClientOutput\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"newConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"newHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"HeaderChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ValSetHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x6080806040523460b9575f54600181811c9116801560b0575b6020821014609c57601f81116058575b507f30372d74656e6465726d696e742d30000000000000000000000000000000001e5f5561294490816100be8239f35b5f8052601f0160051c7f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563908101905b818110609257506028565b5f81556001016087565b634e487b7160e01b5f52602260045260245ffd5b90607f16906018565b5f80fdfe6080806040526004361015610012575f80fd5b5f3560e01c63513a5c0e14610025575f80fd5b3461114d5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261114d576004359067ffffffffffffffff821161114d576103007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc833603011261114d5761009d816119ce565b6040516100a9816119ea565b606081526040516100b981611a3e565b5f81525f602082015260208201526040516100d381611a3e565b5f81525f602082015260408201525f60608201525f60808201525f60a08201525f60c08201528152610103611a7d565b6020820152610110611a7d565b60408201525f606082015260405161012781611a3e565b5f81525f6020820152608082015260a06040519161014483611a3e565b5f8084526020840152015261015c6004820180611a9b565b8035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18136030182121561114d570180359067ffffffffffffffff821161114d5760200190803603821361114d576101bf916101ba913691611aea565b61205d565b60206101ce6004840180611a9b565b019160a06101df6004830180611a9b565b013563ffffffff811680910361114d57610206604051946101ff86611a06565b3690611b2e565b84526020840152600f604084015261022f60406102296084840184600401611b64565b01611b97565b604067ffffffffffffffff61024f60606102296084870187600401611b64565b8183519461025c86611a06565b610264611bac565b8652166020850152169101526102806084820182600401611b64565b9161028d60a48301611cc2565b9160a08436031261114d57604051936102a585611a22565b803567ffffffffffffffff811161114d57810160408136031261114d57604051906102cf82611a3e565b803567ffffffffffffffff811161114d5781016102c08136031261114d5760405190610260820182811067ffffffffffffffff821117611666576040526103163682611cf4565b8252604081013567ffffffffffffffff811161114d576103399036908301611d25565b602083015261034a60608201611cdf565b604083015261035b60808201611d43565b606083015261036c60a08201611d60565b608083015261037e3660c08301611d7e565b60a08301526103906101208201611d60565b60c083015261014081013560e08301526103ad6101608201611d60565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526103fc6102208201611d60565b6101c08301526102408101356101e083015261041b6102608201611d60565b6102008301526102808101356102208301526102a08101359067ffffffffffffffff821161114d5761044f91369101611d25565b610240820152825260208101359067ffffffffffffffff821161114d570160c08136031261114d576040519061048482611a22565b61048d81611cdf565b825261049b60208201611d6d565b60208301526104ad3660408301611d7e565b604083015260a08101359067ffffffffffffffff821161114d570136601f8201121561114d5780356104de81611dd0565b916104ec6040519384611a5a565b81835260208084019260051b8201019036821161114d5760208101925b8284106118bd5750505050606082015260208201528552602081013567ffffffffffffffff811161114d576105419036908301611e52565b60208601526105533660408301611cf4565b604086015260808101359067ffffffffffffffff821161114d5761057991369101611e52565b6060850152610586611bac565b506105943660248301611f42565b9160206105a581875151015161205d565b0167ffffffffffffffff81511667ffffffffffffffff604088015151169081810361188f57505067ffffffffffffffff90511661060d67ffffffffffffffff604088515101511691604051906105fa82611a3e565b8152602081019283526040880151612654565b6003811015610a4f5760028114908115611884575b5061184d5750610664602080870151604051809381927fc87f1f69000000000000000000000000000000000000000000000000000000008352600483016123c9565b038173__$7440a880b7578767f72184d998805816e4$__5af49081156115bf575f9161181b575b5061014086515101518082036117ed5750506106ac6020865151015161205d565b67ffffffffffffffff6020818185015116920151160361176f5761070360206060870151604051809381927fc87f1f69000000000000000000000000000000000000000000000000000000008352600483016123c9565b038173__$7440a880b7578767f72184d998805816e4$__5af49081156115bf575f9161173d575b5060408401518103611693576fffffffffffffffffffffffffffffffff825194511667ffffffffffffffff602060408901510151166060880151916040519660a0880188811067ffffffffffffffff82111761166657604052875260208701526040860152606085015260808401526107eb6020808751970151604051976107b189611a3e565b88528082890152604051809381927fc87f1f69000000000000000000000000000000000000000000000000000000008352600483016123c9565b038173__$7440a880b7578767f72184d998805816e4$__5af49081156115bf575f91611634575b506101408651510151036115ca576109c46020865151604051809381927f95b25f7900000000000000000000000000000000000000000000000000000000835284600484015261087c60248401825167ffffffffffffffff60208092828151168552015116910152565b610240610899868301516102c060648701526102e48601906119a9565b9167ffffffffffffffff60408201511660848601526fffffffffffffffffffffffffffffffff60608201511660a48601526080810151151560c4860152868060a0830151805160e4890152015163ffffffff815116610104880152015161012486015260c0810151151561014486015260e081015161016486015261010081015115156101848601526101208101516101a48601526101408101516101c48601526101608101516101e48601526101808101516102048601526101a08101516102248601526101c081015115156102448601526101e081015161026486015261020081015115156102848601526102208101516102a486015201517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc848303016102c48501526119a9565b038173__$7440a880b7578767f72184d998805816e4$__5af49081156115bf575f9161158d575b506040602087510151015151036115235784519260606020808801519501510151955f975f5b8851811015610a7c57610a2d610a27828b6126c0565b516126d4565b6004811015610a4f57600103610a46575b600101610a11565b60019950610a3e565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b508897959697156114b957855187515103611429575f5b8651811015610bc957610aa681886126c0565b5180516004811015610a4f57600103610ac457506001905b01610a93565b6020015151909890965f96929591949093875b8a518051821015610bb95781610aec916126c0565b5151604051610b1a6020828180820195805191829101875e81015f838201520301601f198101835282611a5a565b5190208a604051610b4a6020828180820195805191829101875e81015f838201520301601f198101835282611a5a565b51902014610b5a57600101610ad7565b509397509398909491955060015b15610b7557600190610abe565b606460405162461bcd60e51b815260206004820152601d60248201527f696e76616c696420636f6d6d69743a206661756c7479207369676e65720000006044820152fd5b5050939750939890949195610b68565b5087610bde63ffffffff602084015116612486565b906fffffffffffffffffffffffffffffffff806020870151169116918183109182156113f0575b5050611386576fffffffffffffffffffffffffffffffff60608451510151166fffffffffffffffffffffffffffffffff602086015116101561131c5760208351510151604051610c746020828180820195805191829101875e81015f838201520301601f198101835282611a5a565b5190208451604051610ca56020828180820195805191829101875e81015f838201520301601f198101835282611a5a565b519020036112d857610cc467ffffffffffffffff604086015116612546565b8351516040015167ffffffffffffffff918216911681810361126a5750506101408351510151608085015103611200575b6fffffffffffffffffffffffffffffffff610d1963ffffffff604085015116612486565b16016fffffffffffffffffffffffffffffffff81116111d3576fffffffffffffffffffffffffffffffff8060608551510151169116111561116957610d6b67ffffffffffffffff604085015116612546565b8251516040015167ffffffffffffffff91821691161461115157610dc192610d9c9160608451920151905191612706565b60405190610da982611a3e565b60028252600360208301526020815191015190612706565b67ffffffffffffffff6020610df96080610df3610ded610de76084890189600401611b64565b80611f7b565b80611fae565b01611cc2565b92610200610e13610ded610de76084890189600401611b64565b01356101c0610e2e610ded610de760848a018a600401611b64565b0135906fffffffffffffffffffffffffffffffff60405196610e4f88611a06565b1686528386015260408501520151169067ffffffffffffffff610e836060610229610ded610de76084890189600401611b64565b60405193610e9085611a3e565b8452166020830152610ea56004840180611a9b565b610eb160a48501611cc2565b926040610ec46084870187600401611b64565b019160405195610ed3876119ce565b6101208236031261114d5760405191610eeb836119ea565b80359067ffffffffffffffff821161114d57610f0d6101009236908301611d25565b8452610f1c3660208301611b2e565b6020850152610f2e3660608301611cf4565b6040850152610f3f60a08201611d6d565b6060850152610f5060c08201611d6d565b6080850152610f6160e08201611d60565b60a0850152013590600282101561114d5782610f8b9260c0610fba95015288526024369101611f42565b9260208701938452604087019485526fffffffffffffffffffffffffffffffff60608801961686523690611cf4565b926080860193845260a0860191825260405195602087525193610180602088015260c0610ff686516101206101a08b01526102c08a01906119a9565b9560ff602080830151828151166101c08d01520151166101e08a015261103b60408201516102008b019067ffffffffffffffff60208092828151168552015116910152565b63ffffffff6060820151166102408a015263ffffffff6080820151166102608a015260a081015115156102808a015201516002811015610a4f5787966110fd611127946110ce611149986fffffffffffffffffffffffffffffffff956102a08d01525160408c0190604080916fffffffffffffffffffffffffffffffff8151168452602081015160208501520151910152565b5180516fffffffffffffffffffffffffffffffff1660a08b0152602081015160c08b01526040015160e08a0152565b51166101008701525161012086019067ffffffffffffffff60208092828151168552015116910152565b5161016084019067ffffffffffffffff60208092828151168552015116910152565b0390f35b5f80fd5b50611164915060405190610da982611a3e565b610dc1565b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b608460405162461bcd60e51b815260206004820152602f60248201527f696e76616c696420626c6f636b3a206e6578742076616c696461746f7220736560448201527f742068617368206d69736d6174636800000000000000000000000000000000006064820152fd5b11610cf557608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152fd5b830391506fffffffffffffffffffffffffffffffff82116111d3576fffffffffffffffffffffffffffffffff8091169116118780610c05565b60a460405162461bcd60e51b815260206004820152604860248201527f696e76616c696420636f6d6d69743a206e756d626572206f66207369676e617460448201527f7572657320646f6573206e6f74206d61746368206e756d626572206f6620766160648201527f6c696461746f72730000000000000000000000000000000000000000000000006084820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420636f6d6d69743a206e6f2070726573656e74207369676e6160448201527f74757265730000000000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152fd5b90506020813d6020116115b7575b816115a860209383611a5a565b8101031261114d57515f6109eb565b3d915061159b565b6040513d5f823e3d90fd5b608460405162461bcd60e51b815260206004820152602a60248201527f696e76616c696420626c6f636b3a2076616c696461746f72207365742068617360448201527f68206d69736d61746368000000000000000000000000000000000000000000006064820152fd5b90506020813d60201161165e575b8161164f60209383611a5a565b8101031261114d57515f610812565b3d9150611642565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60a46040517ff492ef2b00000000000000000000000000000000000000000000000000000000815260206004820152604360248201527f74727573746564206e6578742076616c696461746f722073657420686173682060448201527f646f6573206e6f74206d6174636820686173682073746f726564206f6e20636860648201527f61696e00000000000000000000000000000000000000000000000000000000006084820152fd5b90506020813d602011611767575b8161175860209383611a5a565b8101031261114d57515f61072a565b3d915061174b565b6117b9906117e96020875151015191516040519384937fc6913e9f0000000000000000000000000000000000000000000000000000000085526040600486015260448501906119a9565b907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8483030160248501526119a9565b0390fd5b7f4d6d8014000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b90506020813d602011611845575b8161183660209383611a5a565b8101031261114d57515f61068b565b3d9150611829565b67ffffffffffffffff9051167f90f4dbed000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b60019150145f610622565b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b833567ffffffffffffffff811161114d578201906040601f19833603011261114d57604051916118ec83611a3e565b6020810135600481101561114d578352604081013567ffffffffffffffff811161114d5760209101019060808236031261114d576040519261192d84611a22565b823567ffffffffffffffff811161114d5761194b9036908501611d25565b845261195960208401611d43565b602085015261196a60408401611d60565b604085015260608301359367ffffffffffffffff851161114d57611995602095948695369101611d25565b606082015283820152815201930192610509565b90601f19601f602080948051918291828752018686015e5f8582860101520116010190565b60c0810190811067ffffffffffffffff82111761166657604052565b60e0810190811067ffffffffffffffff82111761166657604052565b6060810190811067ffffffffffffffff82111761166657604052565b6080810190811067ffffffffffffffff82111761166657604052565b6040810190811067ffffffffffffffff82111761166657604052565b90601f601f19910116810190811067ffffffffffffffff82111761166657604052565b60405190611a8a82611a06565b5f6040838281528260208201520152565b9035907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee18136030182121561114d570190565b67ffffffffffffffff811161166657601f01601f191660200190565b929192611af682611ace565b91611b046040519384611a5a565b82948184528183011161114d578281602093845f960137010152565b359060ff8216820361114d57565b919082604091031261114d57604051611b4681611a3e565b6020611b5f818395611b5781611b20565b855201611b20565b910152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff618136030182121561114d570190565b3567ffffffffffffffff8116810361114d5790565b604051905f5f548060011c9160018216918215611cb8575b602084108314611c8b578386528592908115611c4e5750600114611bf1575b611bef92500383611a5a565b565b505f80805290917f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e5635b818310611c32575050906020611bef92820101611be3565b6020919350806001915483858901015201910190918492611c1a565b60209250611bef9491507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001682840152151560051b820101611be3565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b92607f1692611bc4565b356fffffffffffffffffffffffffffffffff8116810361114d5790565b359067ffffffffffffffff8216820361114d57565b919082604091031261114d57604051611d0c81611a3e565b6020611b5f818395611d1d81611cdf565b855201611cdf565b9080601f8301121561114d57816020611d4093359101611aea565b90565b35906fffffffffffffffffffffffffffffffff8216820361114d57565b3590811515820361114d57565b359063ffffffff8216820361114d57565b80929103916060831261114d57604051611d9781611a3e565b6040601f19829584358452011261114d576020906040805193611db985611a3e565b611dc4848201611d6d565b85520135828401520152565b67ffffffffffffffff81116116665760051b60200190565b919060808382031261114d5760405190611e0182611a22565b8193803567ffffffffffffffff811161114d57606092611e22918301611d25565b835260208101356020840152611e3a60408201611cdf565b60408401520135908160070b820361114d5760600152565b91909160808184031261114d5760405190611e6c82611a22565b8193813567ffffffffffffffff811161114d57820181601f8201121561114d578035611e9781611dd0565b91611ea56040519384611a5a565b81835260208084019260051b8201019184831161114d5760208201905b838210611f1457505050508352611edb60208301611d60565b602084015260408201359067ffffffffffffffff821161114d5782611f0960609492611b5f94869401611de8565b604086015201611cdf565b813567ffffffffffffffff811161114d57602091611f3788848094880101611de8565b815201910190611ec2565b919082606091031261114d57604051611f5a81611a06565b6040808294611f6881611d43565b8452602081013560208501520135910152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc18136030182121561114d570190565b9035907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd418136030182121561114d570190565b908151811015611ff2570160200190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b9061202982611ace565b6120366040519182611a5a565b828152601f196120468294611ace565b0190602036910137565b919082018092116111d357565b60405161206981611a3e565b606081525f60209091015280515f1981019081116111d3575b7f2d000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000006120ce8385611fe1565b5116146120e35780156111d3575f1901612082565b91905f1983146123655780518381039081116111d3575f1981019081116111d35761210d9061201f565b905f5b82518110156121705760018501908186116111d3577fff0000000000000000000000000000000000000000000000000000000000000061215b61215583600195612050565b85611fe1565b51165f1a6121698286611fe1565b5301612110565b509192908051158015612307575b61227857805161218d916124bc565b91901580156122f6575b612278576121a48161201f565b905f5b818110612235575050516001811190811591612229575b506121e55767ffffffffffffffff90604051926121da84611a3e565b835216602082015290565b606460405162461bcd60e51b815260206004820152601b60248201527f496e76616c696420636861696e20707265666978206c656e67746800000000006044820152fd5b602b915010155f6121be565b807fff0000000000000000000000000000000000000000000000000000000000000061226360019388611fe1565b51165f1a6122718286611fe1565b53016121a7565b50508051806001111590816122eb575b50156122a6576040519061229b82611a3e565b81525f602082015290565b60405162461bcd60e51b815260206004820152601760248201527f496e76616c696420636861696e204944206c656e6774680000000000000000006044820152606490fd5b60409150105f612288565b5067ffffffffffffffff8211612197565b611ff2577f30000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000060208301511614801561217e5750600181511161217e565b8091925051806001111590816122eb5750156122a6576040519061229b82611a3e565b9060608061239f84516080855260808501906119a9565b936020810151602085015267ffffffffffffffff6040820151166040850152015160070b91015290565b6020815260a0810182519060806020840152815180915260c0830190602060c08260051b8601019301915f905b82821061243d575050505067ffffffffffffffff60606124336080936020870151151560408701526040870151601f198783030184880152612388565b9401511691015290565b90919293602080612478837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff408a600196030186528851612388565b9601920192019092916123f6565b6fffffffffffffffffffffffffffffffff633b9aca00911602906fffffffffffffffffffffffffffffffff82169182036111d357565b5f929183915b8183106124d25750505060019190565b90919360ff6125087fff000000000000000000000000000000000000000000000000000000000000006020888601015116612588565b16906009821161253a57600a810290808204600a14901517156111d35760019161253191612050565b940191906124c2565b5050505090505f905f90565b67ffffffffffffffff60019116019067ffffffffffffffff82116111d357565b9067ffffffffffffffff8091169116019067ffffffffffffffff82116111d357565b60f81c602f81118061264a575b156125c2577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd00160ff1690565b6060811180612640575b156125f9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa90160ff1690565b6040811180612636575b15612630577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc90160ff1690565b5060ff90565b5060478110612603565b50606781106125cc565b50603a8110612595565b67ffffffffffffffff81511667ffffffffffffffff835116908181105f1461267e57505050505f90565b111561268b575050600290565b602067ffffffffffffffff81819301511692015116908181105f146126b05750505f90565b11156126bb57600290565b600190565b8051821015611ff25760209160051b010190565b516004811015610a4f5790565b9067ffffffffffffffff8091169116029067ffffffffffffffff82169182036111d357565b6020015160600151925f9291835b8151805186101561274c5760019167ffffffffffffffff604061273a89612744956126c0565b5101511690612566565b940193612714565b509092919493505f5f935b85518510156128925761276d610a2786886126c0565b6004811015610a4f5760011461288857602061278986886126c0565b510151519560208701905f5b8351805182101561287c57816127aa916126c0565b51516040516127d86020828180820195805191829101875e81015f838201520301601f198101835282611a5a565b519020896040516128046020828181019451808a875e81015f838201520301601f198101835282611a5a565b5190201461281457600101612795565b83985067ffffffffffffffff919794925061273a60409161283595516126c0565b905b61284860ff602086015116836126e1565b67ffffffffffffffff8061286060ff885116876126e1565b16911611612874576001905b019394612757565b505050505050565b50509591965050612837565b949360019061286c565b5067ffffffffffffffff9350839294509060ff6128b96128c29382602089015116906126e1565b955116906126e1565b16911611156128cd57565b608460405162461bcd60e51b815260206004820152602160248201527f696e73756666696369656e7420766f74696e6720706f776572206f7665726c6160448201527f70000000000000000000000000000000000000000000000000000000000000006064820152fdfea164736f6c634300081c000a",
}

// ContractUpdateClientABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractUpdateClientMetaData.ABI instead.
var ContractUpdateClientABI = ContractUpdateClientMetaData.ABI

// ContractUpdateClientBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractUpdateClientMetaData.Bin instead.
var ContractUpdateClientBin = ContractUpdateClientMetaData.Bin

// DeployContractUpdateClient deploys a new Ethereum contract, binding an instance of ContractUpdateClient to it.
func DeployContractUpdateClient(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractUpdateClient, error) {
	parsed, err := ContractUpdateClientMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractUpdateClientBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractUpdateClient{ContractUpdateClientCaller: ContractUpdateClientCaller{contract: contract}, ContractUpdateClientTransactor: ContractUpdateClientTransactor{contract: contract}, ContractUpdateClientFilterer: ContractUpdateClientFilterer{contract: contract}}, nil
}

// ContractUpdateClient is an auto generated Go binding around an Ethereum contract.
type ContractUpdateClient struct {
	ContractUpdateClientCaller     // Read-only binding to the contract
	ContractUpdateClientTransactor // Write-only binding to the contract
	ContractUpdateClientFilterer   // Log filterer for contract events
}

// ContractUpdateClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractUpdateClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractUpdateClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractUpdateClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractUpdateClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractUpdateClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractUpdateClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractUpdateClientSession struct {
	Contract     *ContractUpdateClient // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ContractUpdateClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractUpdateClientCallerSession struct {
	Contract *ContractUpdateClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// ContractUpdateClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractUpdateClientTransactorSession struct {
	Contract     *ContractUpdateClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractUpdateClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractUpdateClientRaw struct {
	Contract *ContractUpdateClient // Generic contract binding to access the raw methods on
}

// ContractUpdateClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractUpdateClientCallerRaw struct {
	Contract *ContractUpdateClientCaller // Generic read-only contract binding to access the raw methods on
}

// ContractUpdateClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractUpdateClientTransactorRaw struct {
	Contract *ContractUpdateClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractUpdateClient creates a new instance of ContractUpdateClient, bound to a specific deployed contract.
func NewContractUpdateClient(address common.Address, backend bind.ContractBackend) (*ContractUpdateClient, error) {
	contract, err := bindContractUpdateClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractUpdateClient{ContractUpdateClientCaller: ContractUpdateClientCaller{contract: contract}, ContractUpdateClientTransactor: ContractUpdateClientTransactor{contract: contract}, ContractUpdateClientFilterer: ContractUpdateClientFilterer{contract: contract}}, nil
}

// NewContractUpdateClientCaller creates a new read-only instance of ContractUpdateClient, bound to a specific deployed contract.
func NewContractUpdateClientCaller(address common.Address, caller bind.ContractCaller) (*ContractUpdateClientCaller, error) {
	contract, err := bindContractUpdateClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractUpdateClientCaller{contract: contract}, nil
}

// NewContractUpdateClientTransactor creates a new write-only instance of ContractUpdateClient, bound to a specific deployed contract.
func NewContractUpdateClientTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractUpdateClientTransactor, error) {
	contract, err := bindContractUpdateClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractUpdateClientTransactor{contract: contract}, nil
}

// NewContractUpdateClientFilterer creates a new log filterer instance of ContractUpdateClient, bound to a specific deployed contract.
func NewContractUpdateClientFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractUpdateClientFilterer, error) {
	contract, err := bindContractUpdateClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractUpdateClientFilterer{contract: contract}, nil
}

// bindContractUpdateClient binds a generic wrapper to an already deployed contract.
func bindContractUpdateClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractUpdateClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractUpdateClient *ContractUpdateClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractUpdateClient.Contract.ContractUpdateClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractUpdateClient *ContractUpdateClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.ContractUpdateClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractUpdateClient *ContractUpdateClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.ContractUpdateClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractUpdateClient *ContractUpdateClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractUpdateClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractUpdateClient *ContractUpdateClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractUpdateClient *ContractUpdateClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.contract.Transact(opts, method, params...)
}

// UpdateClient is a free data retrieval call binding the contract method 0x513a5c0e.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) view returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCaller) UpdateClient(opts *bind.CallOpts, msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	var out []interface{}
	err := _ContractUpdateClient.contract.Call(opts, &out, "updateClient", msg_)

	if err != nil {
		return *new(IUpdateClientMsgsUpdateClientOutput), err
	}

	out0 := *abi.ConvertType(out[0], new(IUpdateClientMsgsUpdateClientOutput)).(*IUpdateClientMsgsUpdateClientOutput)

	return out0, err

}

// UpdateClient is a free data retrieval call binding the contract method 0x513a5c0e.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) view returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientSession) UpdateClient(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClient(&_ContractUpdateClient.CallOpts, msg_)
}

// UpdateClient is a free data retrieval call binding the contract method 0x513a5c0e.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) view returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCallerSession) UpdateClient(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClient(&_ContractUpdateClient.CallOpts, msg_)
}
