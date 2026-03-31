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
	ABI: "[{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.MsgUpdateClient\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.UpdateClientOutput\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"newConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"newHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"HeaderChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ValSetHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x6080806040523460b9575f54600181811c9116801560b0575b6020821014609c57601f81116058575b507f30372d74656e6465726d696e742d30000000000000000000000000000000001e5f5561295590816100be8239f35b5f8052601f0160051c7f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563908101905b818110609257506028565b5f81556001016087565b634e487b7160e01b5f52602260045260245ffd5b90607f16906018565b5f80fdfe6080806040526004361015610012575f80fd5b5f3560e01c6358d248c014610025575f80fd5b346111525760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112611152576004359067ffffffffffffffff8211611152576102407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc83360301126111525761009d816119dc565b6040516100a9816119f8565b606081526040516100b981611a4c565b5f81525f602082015260208201526040516100d381611a4c565b5f81525f602082015260408201525f60608201525f60808201525f60a08201525f60c08201528152610103611a8b565b6020820152610110611a8b565b60408201525f606082015260405161012781611a4c565b5f81525f6020820152608082015260a06040519161014483611a4c565b5f8084526020840152015261015c6004820180611aa9565b8035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215611152570180359067ffffffffffffffff821161115257602001908036038213611152576101bf916101ba913691611af8565b61206b565b60206101ce6004840180611aa9565b019060a06101df6004850180611aa9565b013563ffffffff811680910361115257610206604051936101ff85611a14565b3690611b3c565b83526020830152600f604083015261022f60406102296084860186600401611b72565b01611ba5565b604067ffffffffffffffff61024f60606102296084890189600401611b72565b8183519461025c86611a14565b610264611bba565b8652166020850152169101526102806084840184600401611b72565b9061028d60a48501611cd0565b9360a08336031261115257604051926102a584611a30565b803567ffffffffffffffff811161115257810160408136031261115257604051906102cf82611a4c565b803567ffffffffffffffff81116111525781016102c0813603126111525760405190610260820182811067ffffffffffffffff82111761166b576040526103163682611d02565b8252604081013567ffffffffffffffff8111611152576103399036908301611d33565b602083015261034a60608201611ced565b604083015261035b60808201611d51565b606083015261036c60a08201611d6e565b608083015261037e3660c08301611d8c565b60a08301526103906101208201611d6e565b60c083015261014081013560e08301526103ad6101608201611d6e565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526103fc6102208201611d6e565b6101c08301526102408101356101e083015261041b6102608201611d6e565b6102008301526102808101356102208301526102a08101359067ffffffffffffffff82116111525761044f91369101611d33565b610240820152825260208101359067ffffffffffffffff8211611152570160c081360312611152576040519061048482611a30565b61048d81611ced565b825261049b60208201611d7b565b60208301526104ad3660408301611d8c565b604083015260a08101359067ffffffffffffffff8211611152570136601f820112156111525780356104de81611dde565b916104ec6040519384611a68565b81835260208084019260051b820101903682116111525760208101925b8284106118cb5750505050606082015260208201528452602081013567ffffffffffffffff8111611152576105419036908301611e60565b60208501526105533660408301611d02565b604085015260808101359067ffffffffffffffff82116111525761057991369101611e60565b6060840152610586611bba565b506105943660248301611f50565b9160206105a581865151015161206b565b0167ffffffffffffffff81511667ffffffffffffffff604087015151169081810361189d57505067ffffffffffffffff90511661060d67ffffffffffffffff604087515101511691604051906105fa82611a4c565b815260208101928352604087015161266c565b61061681612554565b60028114908115611889575b506118525750610664602080860151604051809381927fc87f1f69000000000000000000000000000000000000000000000000000000008352600483016123d7565b038173__$7440a880b7578767f72184d998805816e4$__5af49081156115c4575f91611820575b5061014085515101518082036117f25750506106ac6020855151015161206b565b67ffffffffffffffff602081818501511692015116036117745761070360206060860151604051809381927fc87f1f69000000000000000000000000000000000000000000000000000000008352600483016123d7565b038173__$7440a880b7578767f72184d998805816e4$__5af49081156115c4575f91611742575b5060408401518103611698576fffffffffffffffffffffffffffffffff825194511667ffffffffffffffff602060408801510151166060870151916040519660a0880188811067ffffffffffffffff82111761166b57604052875260208701526040860152606085015260808401526107eb6020808651960151604051966107b188611a4c565b87528082880152604051809381927fc87f1f69000000000000000000000000000000000000000000000000000000008352600483016123d7565b038173__$7440a880b7578767f72184d998805816e4$__5af49081156115c4575f91611639575b506101408551510151036115cf576109c46020855151604051809381927f95b25f7900000000000000000000000000000000000000000000000000000000835284600484015261087c60248401825167ffffffffffffffff60208092828151168552015116910152565b610240610899868301516102c060648701526102e48601906119b7565b9167ffffffffffffffff60408201511660848601526fffffffffffffffffffffffffffffffff60608201511660a48601526080810151151560c4860152868060a0830151805160e4890152015163ffffffff815116610104880152015161012486015260c0810151151561014486015260e081015161016486015261010081015115156101848601526101208101516101a48601526101408101516101c48601526101608101516101e48601526101808101516102048601526101a08101516102248601526101c081015115156102448601526101e081015161026486015261020081015115156102848601526102208101516102a486015201517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc848303016102c48501526119b7565b038173__$7440a880b7578767f72184d998805816e4$__5af49081156115c4575f91611592575b506040602086510151015151036115285783519260606020808701519501510151945f965f5b8751811015610a4e57610a2481896126d8565b5151610a2f81612554565b610a3881612554565b610a45575b600101610a11565b60019850610a3d565b508897969596156114be5785518751510361142e575f5b8651811015610ba157610a7881886126d8565b518051610a8481612554565b610a8d81612554565b610a9c57506001905b01610a65565b6020015151909890965f96929591949093875b8a518051821015610b915781610ac4916126d8565b5151604051610af26020828180820195805191829101875e81015f838201520301601f198101835282611a68565b5190208a604051610b226020828180820195805191829101875e81015f838201520301601f198101835282611a68565b51902014610b3257600101610aaf565b509397509398909491955060015b15610b4d57600190610a96565b606460405162461bcd60e51b815260206004820152601d60248201527f696e76616c696420636f6d6d69743a206661756c7479207369676e65720000006044820152fd5b5050939750939890949195610b40565b5087610bb663ffffffff602084015116612494565b906fffffffffffffffffffffffffffffffff806020870151169116918183109182156113f5575b505061138b576fffffffffffffffffffffffffffffffff60608451510151166fffffffffffffffffffffffffffffffff60208601511610156113215760208351510151604051610c4c6020828180820195805191829101875e81015f838201520301601f198101835282611a68565b5190208451604051610c7d6020828180820195805191829101875e81015f838201520301601f198101835282611a68565b519020036112dd57610c9c67ffffffffffffffff60408601511661255e565b8351516040015167ffffffffffffffff918216911681810361126f5750506101408351510151608085015103611205575b6fffffffffffffffffffffffffffffffff610cf163ffffffff604085015116612494565b16016fffffffffffffffffffffffffffffffff81116111d8576fffffffffffffffffffffffffffffffff8060608551510151169116111561116e57610d4367ffffffffffffffff60408501511661255e565b8251516040015167ffffffffffffffff91821691161461115657610d9992610d749160608451920151905191612711565b60405190610d8182611a4c565b60028252600360208301526020815191015190612711565b67ffffffffffffffff6020610dd16080610dcb610dc5610dbf6084890189600401611b72565b80611f89565b80611fbc565b01611cd0565b92610200610deb610dc5610dbf6084890189600401611b72565b01356101c0610e06610dc5610dbf60848a018a600401611b72565b0135906fffffffffffffffffffffffffffffffff60405196610e2788611a14565b1686528386015260408501520151169067ffffffffffffffff610e5b6060610229610dc5610dbf6084890189600401611b72565b60405193610e6885611a4c565b8452166020830152610e7d6004840180611aa9565b610e8960a48501611cd0565b926040610e9c6084870187600401611b72565b019160405195610eab876119dc565b610120823603126111525760405191610ec3836119f8565b80359067ffffffffffffffff821161115257610ee56101009236908301611d33565b8452610ef43660208301611b3c565b6020850152610f063660608301611d02565b6040850152610f1760a08201611d7b565b6060850152610f2860c08201611d7b565b6080850152610f3960e08201611d6e565b60a085015201359060028210156111525782610f639260c0610f9295015288526024369101611f50565b9260208701938452604087019485526fffffffffffffffffffffffffffffffff60608801961686523690611d02565b926080860193845260a0860191825260405195602087525193610180602088015260c0610fce86516101206101a08b01526102c08a01906119b7565b9560ff602080830151828151166101c08d01520151166101e08a015261101360408201516102008b019067ffffffffffffffff60208092828151168552015116910152565b63ffffffff6060820151166102408a015263ffffffff6080820151166102608a015260a081015115156102808a0152015160028110156111255787966110d56110ff946110a6611121986fffffffffffffffffffffffffffffffff956102a08d01525160408c0190604080916fffffffffffffffffffffffffffffffff8151168452602081015160208501520151910152565b5180516fffffffffffffffffffffffffffffffff1660a08b0152602081015160c08b01526040015160e08a0152565b51166101008701525161012086019067ffffffffffffffff60208092828151168552015116910152565b5161016084019067ffffffffffffffff60208092828151168552015116910152565b0390f35b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b5f80fd5b50611169915060405190610d8182611a4c565b610d99565b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b608460405162461bcd60e51b815260206004820152602f60248201527f696e76616c696420626c6f636b3a206e6578742076616c696461746f7220736560448201527f742068617368206d69736d6174636800000000000000000000000000000000006064820152fd5b11610ccd57608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152fd5b830391506fffffffffffffffffffffffffffffffff82116111d8576fffffffffffffffffffffffffffffffff8091169116118780610bdd565b60a460405162461bcd60e51b815260206004820152604860248201527f696e76616c696420636f6d6d69743a206e756d626572206f66207369676e617460448201527f7572657320646f6573206e6f74206d61746368206e756d626572206f6620766160648201527f6c696461746f72730000000000000000000000000000000000000000000000006084820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420636f6d6d69743a206e6f2070726573656e74207369676e6160448201527f74757265730000000000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152fd5b90506020813d6020116115bc575b816115ad60209383611a68565b8101031261115257515f6109eb565b3d91506115a0565b6040513d5f823e3d90fd5b608460405162461bcd60e51b815260206004820152602a60248201527f696e76616c696420626c6f636b3a2076616c696461746f72207365742068617360448201527f68206d69736d61746368000000000000000000000000000000000000000000006064820152fd5b90506020813d602011611663575b8161165460209383611a68565b8101031261115257515f610812565b3d9150611647565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60a46040517ff492ef2b00000000000000000000000000000000000000000000000000000000815260206004820152604360248201527f74727573746564206e6578742076616c696461746f722073657420686173682060448201527f646f6573206e6f74206d6174636820686173682073746f726564206f6e20636860648201527f61696e00000000000000000000000000000000000000000000000000000000006084820152fd5b90506020813d60201161176c575b8161175d60209383611a68565b8101031261115257515f61072a565b3d9150611750565b6117be906117ee6020865151015191516040519384937fc6913e9f0000000000000000000000000000000000000000000000000000000085526040600486015260448501906119b7565b907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8483030160248501526119b7565b0390fd5b7f4d6d8014000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b90506020813d60201161184a575b8161183b60209383611a68565b8101031261115257515f61068b565b3d915061182e565b67ffffffffffffffff9051167f90f4dbed000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b6001915061189681612554565b145f610622565b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b833567ffffffffffffffff8111611152578201906040601f19833603011261115257604051916118fa83611a4c565b60208101356003811015611152578352604081013567ffffffffffffffff811161115257602091010190608082360312611152576040519261193b84611a30565b823567ffffffffffffffff8111611152576119599036908501611d33565b845261196760208401611d51565b602085015261197860408401611d6e565b604085015260608301359367ffffffffffffffff8511611152576119a3602095948695369101611d33565b606082015283820152815201930192610509565b90601f19601f602080948051918291828752018686015e5f8582860101520116010190565b60c0810190811067ffffffffffffffff82111761166b57604052565b60e0810190811067ffffffffffffffff82111761166b57604052565b6060810190811067ffffffffffffffff82111761166b57604052565b6080810190811067ffffffffffffffff82111761166b57604052565b6040810190811067ffffffffffffffff82111761166b57604052565b90601f601f19910116810190811067ffffffffffffffff82111761166b57604052565b60405190611a9882611a14565b5f6040838281528260208201520152565b9035907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee181360301821215611152570190565b67ffffffffffffffff811161166b57601f01601f191660200190565b929192611b0482611adc565b91611b126040519384611a68565b829481845281830111611152578281602093845f960137010152565b359060ff8216820361115257565b919082604091031261115257604051611b5481611a4c565b6020611b6d818395611b6581611b2e565b855201611b2e565b910152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6181360301821215611152570190565b3567ffffffffffffffff811681036111525790565b604051905f5f548060011c9160018216918215611cc6575b602084108314611c99578386528592908115611c5c5750600114611bff575b611bfd92500383611a68565b565b505f80805290917f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e5635b818310611c40575050906020611bfd92820101611bf1565b6020919350806001915483858901015201910190918492611c28565b60209250611bfd9491507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001682840152151560051b820101611bf1565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b92607f1692611bd2565b356fffffffffffffffffffffffffffffffff811681036111525790565b359067ffffffffffffffff8216820361115257565b919082604091031261115257604051611d1a81611a4c565b6020611b6d818395611d2b81611ced565b855201611ced565b9080601f8301121561115257816020611d4e93359101611af8565b90565b35906fffffffffffffffffffffffffffffffff8216820361115257565b3590811515820361115257565b359063ffffffff8216820361115257565b80929103916060831261115257604051611da581611a4c565b6040601f198295843584520112611152576020906040805193611dc785611a4c565b611dd2848201611d7b565b85520135828401520152565b67ffffffffffffffff811161166b5760051b60200190565b91906080838203126111525760405190611e0f82611a30565b8193803567ffffffffffffffff811161115257606092611e30918301611d33565b835260208101356020840152611e4860408201611ced565b60408401520135908160070b82036111525760600152565b9190916080818403126111525760405190611e7a82611a30565b8193813567ffffffffffffffff811161115257820181601f82011215611152578035611ea581611dde565b91611eb36040519384611a68565b81835260208084019260051b820101918483116111525760208201905b838210611f2257505050508352611ee960208301611d6e565b602084015260408201359067ffffffffffffffff82116111525782611f1760609492611b6d94869401611df6565b604086015201611ced565b813567ffffffffffffffff811161115257602091611f4588848094880101611df6565b815201910190611ed0565b919082606091031261115257604051611f6881611a14565b6040808294611f7681611d51565b8452602081013560208501520135910152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc181360301821215611152570190565b9035907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd4181360301821215611152570190565b908151811015612000570160200190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b9061203782611adc565b6120446040519182611a68565b828152601f196120548294611adc565b0190602036910137565b919082018092116111d857565b60405161207781611a4c565b606081525f60209091015280515f1981019081116111d8575b7f2d000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000006120dc8385611fef565b5116146120f15780156111d8575f1901612090565b91905f1983146123735780518381039081116111d8575f1981019081116111d85761211b9061202d565b905f5b825181101561217e5760018501908186116111d8577fff000000000000000000000000000000000000000000000000000000000000006121696121638360019561205e565b85611fef565b51165f1a6121778286611fef565b530161211e565b509192908051158015612315575b61228657805161219b916124ca565b9190158015612304575b612286576121b28161202d565b905f5b818110612243575050516001811190811591612237575b506121f35767ffffffffffffffff90604051926121e884611a4c565b835216602082015290565b606460405162461bcd60e51b815260206004820152601b60248201527f496e76616c696420636861696e20707265666978206c656e67746800000000006044820152fd5b602b915010155f6121cc565b807fff0000000000000000000000000000000000000000000000000000000000000061227160019388611fef565b51165f1a61227f8286611fef565b53016121b5565b50508051806001111590816122f9575b50156122b457604051906122a982611a4c565b81525f602082015290565b60405162461bcd60e51b815260206004820152601760248201527f496e76616c696420636861696e204944206c656e6774680000000000000000006044820152606490fd5b60409150105f612296565b5067ffffffffffffffff82116121a5565b612000577f30000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000060208301511614801561218c5750600181511161218c565b8091925051806001111590816122f95750156122b457604051906122a982611a4c565b906060806123ad84516080855260808501906119b7565b936020810151602085015267ffffffffffffffff6040820151166040850152015160070b91015290565b6020815260a0810182519060806020840152815180915260c0830190602060c08260051b8601019301915f905b82821061244b575050505067ffffffffffffffff60606124416080936020870151151560408701526040870151601f198783030184880152612396565b9401511691015290565b90919293602080612486837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff408a600196030186528851612396565b960192019201909291612404565b6fffffffffffffffffffffffffffffffff633b9aca00911602906fffffffffffffffffffffffffffffffff82169182036111d857565b5f929183915b8183106124e05750505060019190565b90919360ff6125167fff0000000000000000000000000000000000000000000000000000000000000060208886010151166125a0565b16906009821161254857600a810290808204600a14901517156111d85760019161253f9161205e565b940191906124d0565b5050505090505f905f90565b6003111561112557565b67ffffffffffffffff60019116019067ffffffffffffffff82116111d857565b9067ffffffffffffffff8091169116019067ffffffffffffffff82116111d857565b60f81c602f811180612662575b156125da577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd00160ff1690565b6060811180612658575b15612611577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa90160ff1690565b604081118061264e575b15612648577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc90160ff1690565b5060ff90565b506047811061261b565b50606781106125e4565b50603a81106125ad565b67ffffffffffffffff81511667ffffffffffffffff835116908181105f1461269657505050505f90565b11156126a3575050600290565b602067ffffffffffffffff81819301511692015116908181105f146126c85750505f90565b11156126d357600290565b600190565b80518210156120005760209160051b010190565b9067ffffffffffffffff8091169116029067ffffffffffffffff82169182036111d857565b6020015160600151925f9291835b815180518610156127575760019167ffffffffffffffff60406127458961274f956126d8565b510151169061257e565b94019361271f565b509092919493505f5f935b85518510156128a35761277585876126d8565b515161278081612554565b61278981612554565b1561289957602061279a86886126d8565b510151519560208701905f5b8351805182101561288d57816127bb916126d8565b51516040516127e96020828180820195805191829101875e81015f838201520301601f198101835282611a68565b519020896040516128156020828181019451808a875e81015f838201520301601f198101835282611a68565b51902014612825576001016127a6565b83985067ffffffffffffffff919794925061274560409161284695516126d8565b905b61285960ff602086015116836126ec565b67ffffffffffffffff8061287160ff885116876126ec565b16911611612885576001905b019394612762565b505050505050565b50509591965050612848565b949360019061287d565b5067ffffffffffffffff9350839294509060ff6128ca6128d39382602089015116906126ec565b955116906126ec565b16911611156128de57565b608460405162461bcd60e51b815260206004820152602160248201527f696e73756666696369656e7420766f74696e6720706f776572206f7665726c6160448201527f70000000000000000000000000000000000000000000000000000000000000006064820152fdfea164736f6c634300081c000a",
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

// UpdateClient is a free data retrieval call binding the contract method 0x58d248c0.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2]) msg_) view returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCaller) UpdateClient(opts *bind.CallOpts, msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	var out []interface{}
	err := _ContractUpdateClient.contract.Call(opts, &out, "updateClient", msg_)

	if err != nil {
		return *new(IUpdateClientMsgsUpdateClientOutput), err
	}

	out0 := *abi.ConvertType(out[0], new(IUpdateClientMsgsUpdateClientOutput)).(*IUpdateClientMsgsUpdateClientOutput)

	return out0, err

}

// UpdateClient is a free data retrieval call binding the contract method 0x58d248c0.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2]) msg_) view returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientSession) UpdateClient(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClient(&_ContractUpdateClient.CallOpts, msg_)
}

// UpdateClient is a free data retrieval call binding the contract method 0x58d248c0.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64),(uint64,uint64),((bytes,bytes32,uint64,int64)[],bool,(bytes,bytes32,uint64,int64),uint64)),uint128,uint256[8],uint256[2],uint256[2]) msg_) view returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCallerSession) UpdateClient(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClient(&_ContractUpdateClient.CallOpts, msg_)
}
