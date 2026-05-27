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
	ABI: "[{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.MsgUpdateClient\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"validatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedNextValidatorSet\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorSet\",\"components\":[{\"name\":\"validators\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo[]\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"hasProposer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"proposer\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ValidatorInfo\",\"components\":[{\"name\":\"valAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"votingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPriority\",\"type\":\"int64\",\"internalType\":\"int64\"}]},{\"name\":\"totalVotingPower\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"timestampSeconds\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"timestampNanos\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.UpdateClientOutput\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"newConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"newHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"HeaderChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ValSetHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x6080806040523460155761209c908161001a8239f35b5f80fdfe6080806040526004361015610012575f80fd5b5f3560e01c63513a5c0e14610025575f80fd5b34610e6c576020366003190112610e6c57600435906001600160401b038211610e6c576103006003198336030112610e6c57610060816115cc565b60405161006c816115e7565b6060815260405161007c81611602565b5f81525f6020820152602082015260405161009681611602565b5f81525f602082015260408201525f60608201525f60808201525f60a08201525f60c082015281526100c6611674565b60208201526100d3611674565b60408201525f60608201526040516100ea81611602565b5f81525f6020820152608082015260a06040519161010783611602565b5f8084526020840152015261011f6004820180611692565b803590601e1981360301821215610e6c5701908135916001600160401b038311610e6c576020018236038113610e6c576001600160401b0390610188610173606061016d6004870180611692565b016116a8565b916040519561018187611602565b36916116bc565b845216602083015260206101a28260040183600401611692565b019060a06101b36004830180611692565b013563ffffffff8116809103610e6c576101da604051936101d38561161d565b369061170f565b83526020830152600f60408301526101f86084820182600401611745565b9061020560a4820161175a565b9360a083360312610e6c576040519261021d84611638565b80356001600160401b038111610e6c578101604081360312610e6c576040519061024682611602565b80356001600160401b038111610e6c5781016102c081360312610e6c576040519061026082018281106001600160401b038211176112a65760405261028b3682611782565b825260408101356001600160401b038111610e6c576102ad90369083016117b3565b60208301526102be6060820161176e565b60408301526102cf608082016117d1565b60608301526102e060a082016117e5565b60808301526102f23660c08301611803565b60a083015261030461012082016117e5565b60c083015261014081013560e083015261032161016082016117e5565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a083015261037061022082016117e5565b6101c08301526102408101356101e083015261038f61026082016117e5565b6102008301526102808101356102208301526102a0810135906001600160401b038211610e6c576103c2913691016117b3565b61024082015282526020810135906001600160401b038211610e6c570160c081360312610e6c57604051906103f682611638565b6103ff8161176e565b825261040d602082016117f2565b602083015261041f3660408301611803565b604083015260a0810135906001600160401b038211610e6c570136601f82011215610e6c57803561044f81611855565b9161045d6040519384611653565b81835260208084019260051b82010190368211610e6c5760208101925b8284106114c0575050505060608201526020820152845260208101356001600160401b038111610e6c576104b190369083016118d5565b60208501526104c33660408301611782565b60408501526080810135906001600160401b038211610e6c576104e8913691016118d5565b60608401526104fa36602484016119c2565b9261050a60208251510151611b29565b6001600160401b036020820151166001600160401b0360408401515116908181036114925750506001600160401b036020820151166105736001600160401b036040855151015116916040519061056082611602565b8152602081019283526040850151611e4b565b600381101561099e5760028114908115611487575b5061145157506105b16020808401516040518093819263c87f1f6960e01b835260048301611a66565b038173__$7440a880b7578767f72184d998805816e4$__5af490811561129b575f9161141f575b5061014083515101518082036113f15750506001600160401b036020840151166001600160401b03602083015116036113965750610630602060608301516040518093819263c87f1f6960e01b835260048301611a66565b038173__$7440a880b7578767f72184d998805816e4$__5af490811561129b575f91611364575b50604085015181036112ba576001600160801b0383519551166001600160401b03602060408501510151166060840151916040519760a089018981106001600160401b038211176112a6576040528852602088015260408701526060860152608085015261084f602080835193015192604051936106d485611602565b8185528285015251604051809381927f95b25f7900000000000000000000000000000000000000000000000000000000835284600484015261072f6024840182516001600160401b0360208092828151168552015116910152565b61024061074c868301516102c060648701526102e48601906115a8565b60408301516001600160401b0316608486015260608301516001600160801b031660a48601526080830151151560c486015260a0830151805160e4870152870151805163ffffffff1661010487015287015161012486015260c0830151151561014486015260e083015161016486015261010083015115156101848601526101208301516101a48601526101408301516101c48601526101608301516101e48601526101808301516102048601526101a08301516102248601526101c083015115156102448601526101e083015161026486015261020083015115156102848601526102208301516102a4860152910151838203602319016102c48501526115a8565b038173__$7440a880b7578767f72184d998805816e4$__5af490811561129b575f91611269575b506040602083510151015151036111ff57805194606060208084015197015101519485518751510361116f575f975f5b87518110156109b2576108b98189611eb4565b51998a51600481101561099e57600114610994575060019960206108de838c51611eb4565b515160405161090b83828180820195805191829101875e81015f838201520301601f198101835282611653565b5190209101515160405161093e6020828180820195805191829101875e81015f838201520301601f198101835282611653565b51902003610950576001905b016108a6565b606460405162461bcd60e51b815260206004820152601d60248201527f696e76616c696420636f6d6d69743a206661756c7479207369676e65720000006044820152fd5b995060019061094a565b634e487b7160e01b5f52602160045260245ffd5b508815611105576109cc63ffffffff602085015116611da0565b906001600160801b03806020850151169116918183109182156110de575b5050611074576001600160801b0360608551510151166001600160801b03602084015116101561100a5760208451510151604051610a476020828180820195805191829101875e81015f838201520301601f198101835282611653565b5190208251604051610a786020828180820195805191829101875e81015f838201520301601f198101835282611653565b51902003610fc657610a966001600160401b03604084015116611dc4565b845151604001516001600160401b039182169116818103610f585750506101408451510151608083015103610eee575b6001600160801b03610ae163ffffffff604086015116611da0565b16016001600160801b038111610eda576001600160801b0380606086515101511691161115610e7057610b1392611e02565b6001600160401b036020610b4a6080610b44610b3e610b386084890189600401611745565b806119fb565b80611a10565b0161175a565b92610200610b64610b3e610b386084890189600401611745565b01356101c0610b7f610b3e610b3860848a018a600401611745565b0135906001600160801b0360405196610b978861161d565b168652838601526040850152015116906001600160401b03610bca606061016d610b3e610b386084890189600401611745565b60405193610bd785611602565b8452166020830152610bec6004840180611692565b610bf860a4850161175a565b926040610c0b6084870187600401611745565b019160405195610c1a876115cc565b61012082360312610e6c5760405191610c32836115e7565b8035906001600160401b038211610e6c57610c5361010092369083016117b3565b8452610c62366020830161170f565b6020850152610c743660608301611782565b6040850152610c8560a082016117f2565b6060850152610c9660c082016117f2565b6080850152610ca760e082016117e5565b60a08501520135906002821015610e6c5782610cd19260c0610cf7950152885260243691016119c2565b9260208701938452604087019485526001600160801b0360608801961686523690611782565b926080860193845260a0860191825260405195602087525193610180602088015260c0610d3386516101206101a08b01526102c08a01906115a8565b9560ff602080830151828151166101c08d01520151166101e08a0152610d7760408201516102008b01906001600160401b0360208092828151168552015116910152565b63ffffffff6060820151166102408a015263ffffffff6080820151166102608a015260a081015115156102808a01520151600281101561099e578796610e1e610e4794610df8610e68986001600160801b03956102a08d01525160408c0190604080916001600160801b038151168452602081015160208501520151910152565b5180516001600160801b031660a08b0152602081015160c08b01526040015160e08a0152565b5116610100870152516101208601906001600160401b0360208092828151168552015116910152565b516101608401906001600160401b0360208092828151168552015116910152565b0390f35b5f80fd5b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b634e487b7160e01b5f52601160045260245ffd5b608460405162461bcd60e51b815260206004820152602f60248201527f696e76616c696420626c6f636b3a206e6578742076616c696461746f7220736560448201527f742068617368206d69736d6174636800000000000000000000000000000000006064820152fd5b11610ac657608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152fd5b830391506001600160801b038211610eda576001600160801b0380911691161187806109ea565b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420636f6d6d69743a206e6f2070726573656e74207369676e6160448201527f74757265730000000000000000000000000000000000000000000000000000006064820152fd5b60a460405162461bcd60e51b815260206004820152604860248201527f696e76616c696420636f6d6d69743a206e756d626572206f66207369676e617460448201527f7572657320646f6573206e6f74206d61746368206e756d626572206f6620766160648201527f6c696461746f72730000000000000000000000000000000000000000000000006084820152fd5b608460405162461bcd60e51b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152fd5b90506020813d602011611293575b8161128460209383611653565b81010312610e6c57515f610876565b3d9150611277565b6040513d5f823e3d90fd5b634e487b7160e01b5f52604160045260245ffd5b60a46040517ff492ef2b00000000000000000000000000000000000000000000000000000000815260206004820152604360248201527f74727573746564206e6578742076616c696461746f722073657420686173682060448201527f646f6573206e6f74206d6174636820686173682073746f726564206f6e20636860648201527f61696e00000000000000000000000000000000000000000000000000000000006084820152fd5b90506020813d60201161138e575b8161137f60209383611653565b81010312610e6c57515f610657565b3d9150611372565b826113ed6113db925191516040519384937fc6913e9f0000000000000000000000000000000000000000000000000000000085526040600486015260448501906115a8565b838103600319016024850152906115a8565b0390fd5b7f4d6d8014000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b90506020813d602011611449575b8161143a60209383611653565b81010312610e6c57515f6105d8565b3d915061142d565b6001600160401b039051167f90f4dbed000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b60019150145f610588565b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b83356001600160401b038111610e6c578201906040601f198336030112610e6c57604051916114ee83611602565b60208101356004811015610e6c57835260408101356001600160401b038111610e6c57602091010190608082360312610e6c576040519261152e84611638565b82356001600160401b038111610e6c5761154b90369085016117b3565b8452611559602084016117d1565b602085015261156a604084016117e5565b60408501526060830135936001600160401b038511610e6c576115946020959486953691016117b3565b60608201528382015281520193019261047a565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b60c081019081106001600160401b038211176112a657604052565b60e081019081106001600160401b038211176112a657604052565b604081019081106001600160401b038211176112a657604052565b606081019081106001600160401b038211176112a657604052565b608081019081106001600160401b038211176112a657604052565b90601f801991011681019081106001600160401b038211176112a657604052565b604051906116818261161d565b5f6040838281528260208201520152565b90359061011e1981360301821215610e6c570190565b356001600160401b0381168103610e6c5790565b9291926001600160401b0382116112a657604051916116e5601f8201601f191660200184611653565b829481845281830111610e6c578281602093845f960137010152565b359060ff82168203610e6c57565b9190826040910312610e6c5760405161172781611602565b602061174081839561173881611701565b855201611701565b910152565b903590609e1981360301821215610e6c570190565b356001600160801b0381168103610e6c5790565b35906001600160401b0382168203610e6c57565b9190826040910312610e6c5760405161179a81611602565b60206117408183956117ab8161176e565b85520161176e565b9080601f83011215610e6c578160206117ce933591016116bc565b90565b35906001600160801b0382168203610e6c57565b35908115158203610e6c57565b359063ffffffff82168203610e6c57565b809291039160608312610e6c5760405161181c81611602565b6040819483358352601f190112610e6c57602090604080519361183e85611602565b6118498482016117f2565b85520135828401520152565b6001600160401b0381116112a65760051b60200190565b9190608083820312610e6c576040519061188582611638565b819380356001600160401b038111610e6c576060926118a59183016117b3565b8352602081013560208401526118bd6040820161176e565b60408401520135908160070b8203610e6c5760600152565b919091608081840312610e6c57604051906118ef82611638565b819381356001600160401b038111610e6c57820181601f82011215610e6c57803561191981611855565b916119276040519384611653565b81835260208084019260051b82010191848311610e6c5760208201905b8382106119955750505050835261195d602083016117e5565b60208401526040820135906001600160401b038211610e6c578261198a606094926117409486940161186c565b60408601520161176e565b81356001600160401b038111610e6c576020916119b78884809488010161186c565b815201910190611944565b9190826060910312610e6c576040516119da8161161d565b60408082946119e8816117d1565b8452602081013560208501520135910152565b903590603e1981360301821215610e6c570190565b9035906102be1981360301821215610e6c570190565b90606080611a3d84516080855260808501906115a8565b93602081015160208501526001600160401b036040820151166040850152015160070b91015290565b6020815260a0810182519060806020840152815180915260c0830190602060c08260051b8601019301915f905b828210611ad957505050506001600160401b036060611acf6080936020870151151560408701526040870151601f198783030184880152611a26565b9401511691015290565b90919293602080611af660019360bf198a82030186528851611a26565b960192019201909291611a93565b908151811015611b15570160200190565b634e487b7160e01b5f52603260045260245ffd5b604051611b3581611602565b606081525f60208201525080518015908115611d94575b50611d6c575f19908051805b611cf4575b505f198214611ce6576001820190818311610eda577f30000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000611bbe8484611b04565b51161480611cd1575b611cc1575f5b8151831015611c5757611be08383611b04565b5160f81c603081108015611c4d575b611c3b576001600160401b0360ff8192602f19011681600a850216011691168110611c1f57600190920191611bcd565b5091505060405190611c3082611602565b81525f602082015290565b505091505060405190611c3082611602565b5060398111611bef565b9150916001811190811591611cb5575b50611c8d576001600160401b039060405192611c8284611602565b835216602082015290565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f611c67565b91505060405190611c3082611602565b508051838103908111610eda57600210611bc7565b6040519150611c3082611602565b5f198101818111610eda577f2d000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000611d4b8386611b04565b511614611d6257508015610eda575f190180611b58565b92505f9050611b5d565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f611b4c565b6001600160801b03633b9aca00911602906001600160801b038216918203610eda57565b6001600160401b036001911601906001600160401b038211610eda57565b906001600160401b03809116911601906001600160401b038211610eda57565b9091611e1a6001600160401b03604085015116611dc4565b9151916001600160401b038060408551015116911614611e46576060611e44930151905191611ec8565b565b505050565b6001600160401b038151166001600160401b03835116908181105f14611e7357505050505f90565b1115611e80575050600290565b60206001600160401b0381819301511692015116908181105f14611ea45750505f90565b1115611eaf57600290565b600190565b8051821015611b155760209160051b010190565b60200151606001519051925f91825b8551841015611f0c57611f046001916001600160401b036040611efa888b611eb4565b5101511690611de2565b930192611ed7565b9093919492505f5f935b8551851015611fcd57611f298587611eb4565b518051600481101561099e57600114611fc257602001515160208151910120955f5b8251811015611fb85787611f5f8285611eb4565b51516020815191012014611f7557600101611f4b565b8297506040611efa6001600160401b0392611f939599969499611eb4565b905b611fa0848484612060565b611fb0576001905b019394611f16565b505050505050565b5094909550611f95565b509493600190611fa8565b5090919250611fdc9350612060565b15611fe357565b608460405162461bcd60e51b815260206004820152602160248201527f696e73756666696369656e7420766f74696e6720706f776572206f7665726c6160448201527f70000000000000000000000000000000000000000000000000000000000000006064820152fd5b81810292918115918404141715610eda57565b906001600160401b0360ff61208161208b948383602089015116911661204d565b945116911661204d565b109056fea164736f6c634300081c000a",
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
