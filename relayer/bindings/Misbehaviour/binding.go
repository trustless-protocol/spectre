// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractMisbehaviour

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

// IICS07TendermintMsgsChainId is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsChainId struct {
	Id             string
	RevisionNumber uint64
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
	ClockDrift      uint32
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
	SignedHeader  IICS07TendermintMsgsSignedHeader
	TrustedHeight IICS02ClientMsgsHeight
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

// IICS07TendermintMsgsVersion is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsVersion struct {
	BlockVersion uint64
	AppVersion   uint64
}

// IMisbehaviourMsgsMisbehaviour is an auto generated low-level Go binding around an user-defined struct.
type IMisbehaviourMsgsMisbehaviour struct {
	ClientId IICS07TendermintMsgsChainId
	Header1  IICS07TendermintMsgsHeader
	Header2  IICS07TendermintMsgsHeader
}

// IMisbehaviourMsgsMisbehaviourOutput is an auto generated low-level Go binding around an user-defined struct.
type IMisbehaviourMsgsMisbehaviourOutput struct {
	TrustedHeight1 IICS02ClientMsgsHeight
	TrustedHeight2 IICS02ClientMsgsHeight
}

// ContractMisbehaviourMetaData contains all meta data concerning the ContractMisbehaviour contract.
var ContractMisbehaviourMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"misbehaviour\",\"inputs\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"},{\"name\":\"clockDrift\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"misbehaviour_\",\"type\":\"tuple\",\"internalType\":\"structIMisbehaviourMsgs.Misbehaviour\",\"components\":[{\"name\":\"client_id\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ChainId\",\"components\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"header1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"header2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}]},{\"name\":\"trustedConsensusState1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"trustedConsensusState2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIMisbehaviourMsgs.MisbehaviourOutput\",\"components\":[{\"name\":\"trustedHeight1\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight2\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MisbehaviourNotDetected\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]",
	Bin: "0x608080604052346015576121cb908161001a8239f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c6385fbdd9c14610024575f80fd5b3461015657610120366003190112610156576004356001600160401b0381116101565761014060031982360301126101565761005e61016e565b9080600401356001600160401b038111610156576100fa916100896101249260043691840101610250565b8452610098366024830161029d565b60208501526100aa36606483016102e3565b60408501526100bb60a4820161030e565b60608501526100cc60c4820161030e565b60808501526100dd60e4820161031f565b60a08501526100ef610104820161032c565b60c08501520161030e565b60e08201526024356001600160401b0381116101565761015291610125610146923690600401610722565b61012e366107d6565b61013736610812565b91610140610339565b93610902565b6040519182918261084e565b0390f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b6040519061010082018281106001600160401b0382111761018e57604052565b61015a565b60405190604082018281106001600160401b0382111761018e57604052565b60405190608082018281106001600160401b0382111761018e57604052565b6040519061026082018281106001600160401b0382111761018e57604052565b60405190606082018281106001600160401b0382111761018e57604052565b6040519190601f01601f191682016001600160401b0381118382101761018e57604052565b6001600160401b03811161018e57601f01601f191660200190565b81601f820112156101565760208135910161027261026d83610235565b610210565b928284528282011161015657815f92602092838601378301015290565b359060ff8216820361015657565b9190826040910312610156576102c860206102b6610193565b936102c08161028f565b85520161028f565b6020830152565b35906001600160401b038216820361015657565b9190826040910312610156576102c860206102fc610193565b93610306816102cf565b8552016102cf565b359063ffffffff8216820361015657565b3590811515820361015657565b3590600282101561015657565b61010435906001600160801b038216820361015657565b35906001600160801b038216820361015657565b80929103916060831261015657604061037b610193565b8235815293601f190112610156576040610393610193565b916103a06020820161030e565b8352013560208201526020830152565b6001600160401b03811161018e5760051b60200190565b919060c083820312610156576103db6101b2565b926103e5816102cf565b84526103f36020820161030e565b60208501526104058260408301610364565b604085015260a0810135906001600160401b03821161015657019080601f830112156101565781359161043a61026d846103b0565b9260208085838152019160051b830101918383116101565760208101915b83831061046b5750505050506060830152565b82356001600160401b038111610156578201906040828703601f19011261015657610494610193565b916020810135600481101561015657835260408101356001600160401b03811161015657602091010190608082880312610156576104d06101b2565b9282356001600160401b03811161015657886104ed918501610250565b84526104fb60208401610350565b602085015261050c6040840161031f565b60408501526060830135936001600160401b0385116101565761053489602096879601610250565b606082015283820152815201920191610458565b9190916060818403126101565761055d610193565b9281356001600160401b03811161015657820160408183031261015657610582610193565b9281356001600160401b0381116101565782016102c081850312610156576105a86101d1565b906105b385826102e3565b825260408101356001600160401b03811161015657856105d4918301610250565b60208301526105e5606082016102cf565b60408301526105f660808201610350565b606083015261060760a0820161031f565b60808301526106198560c08301610364565b60a083015261062b610120820161031f565b60c083015261014081013560e0830152610648610160820161031f565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a0830152610697610220820161031f565b6101c08301526102408101356101e08301526106b6610260820161031f565b6102008301526102808101356102208301526102a0810135906001600160401b038211610156576106e991869101610250565b61024082015284526020820135936001600160401b03851161015657610716846102c896602095016103c7565b838201528652016102e3565b919091606081840312610156576107376101f1565b9281356001600160401b0381116101565782016040818303126101565761075c610193565b9080356001600160401b038111610156578161077f856020936107879501610250565b8452016102cf565b6020820152845260208201356001600160401b03811161015657816107ad918401610548565b602085015260408201356001600160401b038111610156576107cf9201610548565b6040830152565b6060906043190112610156576107ea6101f1565b906044356001600160801b038116810361015657825260643560208301526084356040830152565b60609060a3190112610156576108266101f1565b9060a4356001600160801b038116810361015657825260c435602083015260e4356040830152565b61089d9092919260406020608083019561087e8482516001600160401b0360208092828151168552015116910152565b01519101906001600160401b0360208092828151168552015116910152565b565b6108a7610193565b906108b0610193565b5f81525f602082015282526108c3610193565b5f81525f60208201526020830152565b156108da57565b7fa179f8c9000000000000000000000000000000000000000000000000000000005f5260045ffd5b9392909261090e61089f565b5084518051906020012094602085019586515151602001518051906020012014610937906108d3565b61094085610ade565b60208101519460608201516109589063ffffffff1690565b60e083015163ffffffff169061096c6101f1565b97885263ffffffff16602088015263ffffffff1660408701528151604090920151516001600160401b0316946109a0610193565b92835260208301956109bb9087906001600160401b03169052565b8383888a5193516109d2906001600160801b031690565b916109dc94610e5e565b60400194855193516109f4906001600160801b031690565b916109fe94610e5e565b80516001600160401b0316925160209081015101516001600160401b0316610a24610193565b6001600160401b0390941684526001600160401b03166020840152516001600160401b0316905160209081015101516001600160401b0316610a64610193565b6001600160401b0390921682526001600160401b03166020820152610a87610193565b918252602082015290565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b9091610acd610adb93604084526040840190610a92565b916020818403910152610a92565b90565b60406020820191610aef8351610fed565b0190610afb8251610fed565b6020815151510180516020815191012090602084515151019182516020815191012003610cc25750506020610b3581835151510151611246565b016001600160401b03610c07610bfb610bed610bb2610b5b86516001600160401b031690565b95610ba56020610b7960408b515151016001600160401b0390511690565b98610b94610b85610193565b6001600160401b039092168252565b019788906001600160401b03169052565b516001600160401b031690565b94610ba56020610bd060408b515151016001600160401b0390511690565b97610bdc610b85610193565b019687906001600160401b03169052565b93516001600160401b031690565b6001600160401b031690565b911603610c53576020604081819351510151015151925151015101515114610c2b57565b7f93c0a5f9000000000000000000000000000000000000000000000000000000005f5260045ffd5b90610c856040610c7381610cbf95515151016001600160401b0390511690565b9251515101516001600160401b031690565b7f494c87bc000000000000000000000000000000000000000000000000000000005f526001600160401b0391821660045216602452604490565b5ffd5b51905190610cfb6040519283927ff6b6676b00000000000000000000000000000000000000000000000000000000845260048401610ab6565b0390fd5b634e487b7160e01b5f52601160045260245ffd5b906001600160801b03809116911603906001600160801b038211610d3357565b610cff565b15610d3f57565b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b6001600160801b03633b9aca00911602906001600160801b038216918203610d3357565b906001600160801b03809116911601906001600160801b038211610d3357565b15610df457565b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b929493919094633b9aca00610e7d610e7682866114a1565b91836114a1565b906001600160801b0382166001600160801b03821610610fb85790610ea191610d13565b95610eb3602084015163ffffffff1690565b63ffffffff81166001600160801b0389161015610f7a575061089d959650610f5f6001600160801b0394610f59610f54610f4b60408998610f40610f70998b808f6060610f3691610f26829f610f0981516114d3565b6020835151015160208151910120905160208151910120146108d3565b515101516001600160801b031690565b9216911611610d38565b015163ffffffff1690565b63ffffffff1690565b610da9565b90610dcd565b94515101516001600160801b031690565b9216911610610ded565b7fdb020436000000000000000000000000000000000000000000000000000000005f526001600160801b03881660045263ffffffff1660245260445ffd5b7fe42a1980000000000000000000000000000000000000000000000000000000005f526001600160801b03831660045260245ffd5b6020610ffd818351510151611246565b0180516001600160401b03169060208301916110218351516001600160401b031690565b906001600160401b0382166001600160401b0382160361116f57505051611093906001600160401b0316835151604001516001600160401b031692611076611067610193565b6001600160401b039093168352565b61108d602083019485906001600160401b03169052565b5161171d565b61113a5750604060206110a783515161175b565b92510151015151036110b557565b6040517ff492ef2b00000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152608490fd5b517f90f4dbed000000000000000000000000000000000000000000000000000000005f526001600160401b031660045260245ffd5b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f526001600160401b039081166004521660245260445ffd5b6111b1610193565b90606082525f6020830152565b8015610d33575f190190565b5f19810191908211610d3357565b91908203918211610d3357565b634e487b7160e01b5f52603260045260245ffd5b90815181101561120a570160200190565b6111e5565b9060018201809211610d3357565b6001019081600111610d3357565b6023019081602311610d3357565b91908201809211610d3357565b61124e6111a9565b5080518015908115611495575b5061146d575f19908051805b611405575b505f1982146113e757600360fc1b7fff000000000000000000000000000000000000000000000000000000000000006112d66112b06112aa8661120f565b856111f9565b517fff000000000000000000000000000000000000000000000000000000000000001690565b1614806113f1575b6113e7575f906112ed8361120f565b915b815183101561137f5761130e6113086112b085856111f9565b60f81c90565b60ff811660308110908115611374575b50611367576001600160401b0360ff8192602f19011681600a85021601169116811061134f576001909201916112ef565b5091505061135b610193565b9081525f602082015290565b505091505061135b610193565b60399150115f61131e565b91509160018111908115916113db575b506113b357610adb906113a0610193565b9283526001600160401b03166020830152565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f61138f565b905061135b610193565b5060026113ff8383516111d8565b116112de565b602d60f81b6114476114226112b061141c856111ca565b866111f9565b7fff000000000000000000000000000000000000000000000000000000000000001690565b1461145b57611455906111be565b80611267565b6114669192506111ca565b905f61126c565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f61125b565b906001600160801b03169081156114bf576001600160801b03160490565b634e487b7160e01b5f52601260045260245ffd5b80518015908115611712575b506116cd575f5b81518110156116c9576114ff6114226112b083856111f9565b7f6100000000000000000000000000000000000000000000000000000000000000811015908161169e575b8115611642575b8115611602575b81156115f4575b81156115ca575b81156115a0575b501561155b576001016114e6565b60405162461bcd60e51b815260206004820152601860248201527f696e76616c696420636861696e206964206368617273657400000000000000006044820152606490fd5b7f2e000000000000000000000000000000000000000000000000000000000000009150145f61154d565b7f5f0000000000000000000000000000000000000000000000000000000000000081149150611546565b602d60f81b8114915061153f565b9050600360fc1b81101580611618575b90611538565b507f3900000000000000000000000000000000000000000000000000000000000000811115611612565b90507f410000000000000000000000000000000000000000000000000000000000000081101580611674575b90611531565b507f5a0000000000000000000000000000000000000000000000000000000000000081111561166e565b7f7a00000000000000000000000000000000000000000000000000000000000000811115915061152a565b5050565b60405162461bcd60e51b815260206004820152601760248201527f496e76616c696420636861696e206964206c656e6774680000000000000000006044820152606490fd5b60329150115f6114df565b90611727916119fb565b6003811015611747576002811490811561173f575090565b600191501490565b634e487b7160e01b5f52602160045260245ffd5b610adb9061192a61177861024061177d6117786020860151611bd1565b611c5b565b93611786611a93565b946117946117788351611cb0565b61179d87611ab1565b526117a786611abe565b526117ca6117786117c5610bfb60408501516001600160401b031690565b611dc6565b6117d386611ace565b526117f36117786117ee60608401516001600160801b031690565b611dfa565b6117fc86611ade565b526080810151156119d55761181761177860a0830151611ec5565b61182086611aee565b5260c0810151156119af5761183b61177860e0830151611f91565b61184486611afe565b526101008101511561198957611861611778610120830151611f91565b61186a86611b0e565b5261187c611778610140830151611f91565b61188586611b1e565b52611897611778610160830151611f91565b6118a086611b2f565b526118b2611778610180830151611f91565b6118bb86611b40565b526118cd6117786101a0830151611f91565b6118d686611b51565b526101c081015115611963576118f36117786101e0830151611f91565b6118fc86611b62565b526102008101511561193d57611919611778610220830151611f91565b61192286611b73565b520151611bd1565b61193382611b84565b525f815191611fbe565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611919565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d6118f3565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611861565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d61183b565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611817565b80516001600160401b03166001600160401b03611a22610bfb85516001600160401b031690565b911681811015611a3457505050505f90565b1115611a41575050600290565b611a73610bfb6020611a64816001600160401b039501516001600160401b031690565b9401516001600160401b031690565b911681811015611a835750505f90565b1115611a8e57600290565b600190565b6101e090611aa082610210565b600e815291601f1901366020840137565b80511561120a5760200190565b80516001101561120a5760400190565b80516002101561120a5760600190565b80516003101561120a5760800190565b80516004101561120a5760a00190565b80516005101561120a5760c00190565b80516006101561120a5760e00190565b80516007101561120a576101000190565b80516008101561120a576101200190565b80516009101561120a576101400190565b8051600a101561120a576101600190565b8051600b101561120a576101800190565b8051600c101561120a576101a00190565b8051600d101561120a576101c00190565b805182101561120a5760209160051b010190565b90611bb661026d83610235565b8281528092611bc7601f1991610235565b0190602036910137565b805115611c2657611be28151612065565b806001019081600111610d335760019083510101809111610d3357611c09611c2291611ba9565b91600a6020840153611c1c8151846120e0565b8361217b565b5090565b50611c316020610210565b5f80825236602083013790565b805191908290602001825e015f815290565b6040513d5f823e3d90fd5b805160018101809111610d3357611c7190611ba9565b80511561120a57611c9b81611c8e6020945f868196015382612151565b5060405191828092611c3e565b039060025afa15611cab575f5190565b611c50565b5f90611cc381516001600160401b031690565b6001600160401b038116611da1575b50611cfb6020820192611cef610bfb85516001600160401b031690565b80611d80575b50611ba9565b915f91611d12610bfb82516001600160401b031690565b611d5d575b50611d2c610bfb82516001600160401b031690565b611d3557505090565b611d56610bfb611d48611c229486612090565b92516001600160401b031690565b9083612118565b611d79919250611d72610bfb611d4886612084565b9084612118565b905f611d17565b90611d95611d90611d9b93612065565b61121d565b90611239565b5f611cf5565b611dbf919250611dba611d90916001600160401b031690565b612065565b905f611cd2565b8015611c2657611dd581612065565b60010180600111610d3357611dec611c2291611ba9565b9160086020840153826120e0565b6001600160801b0316633b9aca008104906001600160801b035f9216918215159182611ea1575b633b9aca006001600160801b03910616908115159081611e86575b611e4590611ba9565b935f93611e6a575b50611e5757505090565b611e64611c229284612090565b83612118565b611e7f91935060086020860153600185612118565b915f611e4d565b611e92611d9084612065565b810180911115611e3c57610cff565b90506001600160801b03633b9aca00611ebc611d9086612065565b92915050611e21565b6020810151611eda63ffffffff825116612065565b90816001019182600111610d3357602301809211610d3357611f41611f01611c2293611ba9565b91600860208401536020611f37611f31611f2b611f25610f4b865163ffffffff1690565b876120e0565b866120a7565b856120be565b91015190836121a6565b50611c1c8151611f8b611f2b611f6f611f6a84611f65611f6082612065565b61122b565b611239565b611ba9565b96611f82611f7c896120d4565b896120be565b905190886121a6565b85612118565b611c22611f9e6060610210565b6022815291600a60208401604036823753602060218401536002836121a6565b9092919280840393808511610d3357600185146120545760015b8060011b9086821015611feb5750611fd8565b91929394955050820191828111610d3357826120079185611fbe565b916120129293611fbe565b9061201d6080610210565b604181526020810190606036833780511561120a576020935f936001611c9b94536021830152604182015260405191828092611c3e565b5090612061929350611b95565b5190565b906001915b60808110156120765750565b60019060071c92019161206a565b60206008910153600190565b60208260109201015360018101809111610d335790565b60208260129201015360018101809111610d335790565b602082819201015360018101809111610d335790565b6020600a910153600190565b91906021600193015b60808210156120fd57906001929391530190565b600180916080607f85161781530193019060071c90926120e9565b9092919083016020015b608082101561213657906001929391530190565b600180916080607f85161781530193019060071c9092612122565b908051918215612173576021602084930191015e60010180600111610d335790565b505050600190565b90809291825192831561219f57839260208092019201015e8101809111610d335790565b5050505090565b8160209193929301015260208101809111610d33579056fea164736f6c634300081c000a",
}

// ContractMisbehaviourABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMisbehaviourMetaData.ABI instead.
var ContractMisbehaviourABI = ContractMisbehaviourMetaData.ABI

// ContractMisbehaviourBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMisbehaviourMetaData.Bin instead.
var ContractMisbehaviourBin = ContractMisbehaviourMetaData.Bin

// DeployContractMisbehaviour deploys a new Ethereum contract, binding an instance of ContractMisbehaviour to it.
func DeployContractMisbehaviour(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractMisbehaviour, error) {
	parsed, err := ContractMisbehaviourMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractMisbehaviourBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractMisbehaviour{ContractMisbehaviourCaller: ContractMisbehaviourCaller{contract: contract}, ContractMisbehaviourTransactor: ContractMisbehaviourTransactor{contract: contract}, ContractMisbehaviourFilterer: ContractMisbehaviourFilterer{contract: contract}}, nil
}

// ContractMisbehaviour is an auto generated Go binding around an Ethereum contract.
type ContractMisbehaviour struct {
	ContractMisbehaviourCaller     // Read-only binding to the contract
	ContractMisbehaviourTransactor // Write-only binding to the contract
	ContractMisbehaviourFilterer   // Log filterer for contract events
}

// ContractMisbehaviourCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractMisbehaviourCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMisbehaviourTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractMisbehaviourTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMisbehaviourFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractMisbehaviourFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractMisbehaviourSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractMisbehaviourSession struct {
	Contract     *ContractMisbehaviour // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ContractMisbehaviourCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractMisbehaviourCallerSession struct {
	Contract *ContractMisbehaviourCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// ContractMisbehaviourTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractMisbehaviourTransactorSession struct {
	Contract     *ContractMisbehaviourTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractMisbehaviourRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractMisbehaviourRaw struct {
	Contract *ContractMisbehaviour // Generic contract binding to access the raw methods on
}

// ContractMisbehaviourCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractMisbehaviourCallerRaw struct {
	Contract *ContractMisbehaviourCaller // Generic read-only contract binding to access the raw methods on
}

// ContractMisbehaviourTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractMisbehaviourTransactorRaw struct {
	Contract *ContractMisbehaviourTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractMisbehaviour creates a new instance of ContractMisbehaviour, bound to a specific deployed contract.
func NewContractMisbehaviour(address common.Address, backend bind.ContractBackend) (*ContractMisbehaviour, error) {
	contract, err := bindContractMisbehaviour(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractMisbehaviour{ContractMisbehaviourCaller: ContractMisbehaviourCaller{contract: contract}, ContractMisbehaviourTransactor: ContractMisbehaviourTransactor{contract: contract}, ContractMisbehaviourFilterer: ContractMisbehaviourFilterer{contract: contract}}, nil
}

// NewContractMisbehaviourCaller creates a new read-only instance of ContractMisbehaviour, bound to a specific deployed contract.
func NewContractMisbehaviourCaller(address common.Address, caller bind.ContractCaller) (*ContractMisbehaviourCaller, error) {
	contract, err := bindContractMisbehaviour(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractMisbehaviourCaller{contract: contract}, nil
}

// NewContractMisbehaviourTransactor creates a new write-only instance of ContractMisbehaviour, bound to a specific deployed contract.
func NewContractMisbehaviourTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractMisbehaviourTransactor, error) {
	contract, err := bindContractMisbehaviour(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractMisbehaviourTransactor{contract: contract}, nil
}

// NewContractMisbehaviourFilterer creates a new log filterer instance of ContractMisbehaviour, bound to a specific deployed contract.
func NewContractMisbehaviourFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractMisbehaviourFilterer, error) {
	contract, err := bindContractMisbehaviour(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractMisbehaviourFilterer{contract: contract}, nil
}

// bindContractMisbehaviour binds a generic wrapper to an already deployed contract.
func bindContractMisbehaviour(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMisbehaviourMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractMisbehaviour *ContractMisbehaviourRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractMisbehaviour.Contract.ContractMisbehaviourCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractMisbehaviour *ContractMisbehaviourRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.ContractMisbehaviourTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractMisbehaviour *ContractMisbehaviourRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.ContractMisbehaviourTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractMisbehaviour *ContractMisbehaviourCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractMisbehaviour.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractMisbehaviour *ContractMisbehaviourTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractMisbehaviour *ContractMisbehaviourTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.contract.Transact(opts, method, params...)
}

// Misbehaviour is a free data retrieval call binding the contract method 0x85fbdd9c.
//
// Solidity: function misbehaviour((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32) clientState, ((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64))) misbehaviour_, (uint128,bytes32,bytes32) trustedConsensusState1, (uint128,bytes32,bytes32) trustedConsensusState2, uint128 time) pure returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourCaller) Misbehaviour(opts *bind.CallOpts, clientState IICS07TendermintMsgsClientState, misbehaviour_ IMisbehaviourMsgsMisbehaviour, trustedConsensusState1 IICS07TendermintMsgsConsensusState, trustedConsensusState2 IICS07TendermintMsgsConsensusState, time *big.Int) (IMisbehaviourMsgsMisbehaviourOutput, error) {
	var out []interface{}
	err := _ContractMisbehaviour.contract.Call(opts, &out, "misbehaviour", clientState, misbehaviour_, trustedConsensusState1, trustedConsensusState2, time)

	if err != nil {
		return *new(IMisbehaviourMsgsMisbehaviourOutput), err
	}

	out0 := *abi.ConvertType(out[0], new(IMisbehaviourMsgsMisbehaviourOutput)).(*IMisbehaviourMsgsMisbehaviourOutput)

	return out0, err

}

// Misbehaviour is a free data retrieval call binding the contract method 0x85fbdd9c.
//
// Solidity: function misbehaviour((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32) clientState, ((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64))) misbehaviour_, (uint128,bytes32,bytes32) trustedConsensusState1, (uint128,bytes32,bytes32) trustedConsensusState2, uint128 time) pure returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourSession) Misbehaviour(clientState IICS07TendermintMsgsClientState, misbehaviour_ IMisbehaviourMsgsMisbehaviour, trustedConsensusState1 IICS07TendermintMsgsConsensusState, trustedConsensusState2 IICS07TendermintMsgsConsensusState, time *big.Int) (IMisbehaviourMsgsMisbehaviourOutput, error) {
	return _ContractMisbehaviour.Contract.Misbehaviour(&_ContractMisbehaviour.CallOpts, clientState, misbehaviour_, trustedConsensusState1, trustedConsensusState2, time)
}

// Misbehaviour is a free data retrieval call binding the contract method 0x85fbdd9c.
//
// Solidity: function misbehaviour((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32) clientState, ((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64))) misbehaviour_, (uint128,bytes32,bytes32) trustedConsensusState1, (uint128,bytes32,bytes32) trustedConsensusState2, uint128 time) pure returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourCallerSession) Misbehaviour(clientState IICS07TendermintMsgsClientState, misbehaviour_ IMisbehaviourMsgsMisbehaviour, trustedConsensusState1 IICS07TendermintMsgsConsensusState, trustedConsensusState2 IICS07TendermintMsgsConsensusState, time *big.Int) (IMisbehaviourMsgsMisbehaviourOutput, error) {
	return _ContractMisbehaviour.Contract.Misbehaviour(&_ContractMisbehaviour.CallOpts, clientState, misbehaviour_, trustedConsensusState1, trustedConsensusState2, time)
}
