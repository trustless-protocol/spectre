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

// IUpdateClientMsgsMsgUpdateClient is an auto generated low-level Go binding around an user-defined struct.
type IUpdateClientMsgsMsgUpdateClient struct {
	ClientState            IICS07TendermintMsgsClientState
	TrustedConsensusState  IICS07TendermintMsgsConsensusState
	ProposedHeader         IICS07TendermintMsgsHeader
	Time                   *big.Int
	Proof                  [8]*big.Int
	Commitments            [2]*big.Int
	CommitmentPok          [2]*big.Int
	Bucket                 uint16
	SignerIndices          []uint32
	PinnedValidatorIndices []uint32
	SignerPubkeys          [][32]byte
	TimestampSeconds       []uint64
	TimestampNanos         []uint32
	Active                 []bool
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
	ABI: "[{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.MsgUpdateClient\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"},{\"name\":\"clockDrift\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pinnedValidatorIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"timestampSeconds\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"timestampNanos\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIUpdateClientMsgs.UpdateClientOutput\",\"components\":[{\"name\":\"clientState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ClientState\",\"components\":[{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"trustLevel\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.TrustThreshold\",\"components\":[{\"name\":\"numerator\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"denominator\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"latestHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unbondingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isFrozen\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"zkAlgorithm\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.SupportedZkAlgorithm\"},{\"name\":\"clockDrift\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"newConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"newHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"HeaderChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]",
	Bin: "0x60808060405234601557612159908161001a8239f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c63c8f5eea014610024575f80fd5b3461018d57602036600319011261018d576004356001600160401b03811161018d578060040190610320600319823603011261018d576101899161017761017d9261006d610489565b5061016a61008461007e8580610514565b8061052a565b6100b061009f60606100998980989698610514565b0161056d565b946100a8610375565b923691610592565b81526100c9602082019485906001600160401b03169052565b60206100d58780610514565b016101316100ee60a06100e88a80610514565b016105d0565b6101246101016101206100e88c80610514565b9161011561010d610399565b9536906105e8565b855263ffffffff166020850152565b63ffffffff166040830152565b61013e608484018861061a565b9161016461015961015160a48701610640565b9436906108b0565b946024369101610a8a565b93610ac1565b516001600160401b031690565b90610d64565b604051918291826101db565b0390f35b5f80fd5b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b634e487b7160e01b5f52602160045260245ffd5b9060028210156101d65752565b6101b5565b61035e906020815261016060a084519461018060208501526102a660e061021188516101406101a08901526102e0880190610191565b9760ff602080830151828151166101c08b01520151166101e088015261025560408201516102008901906001600160401b0360208092828151168552015116910152565b606081015163ffffffff16610240880152608081015163ffffffff1661026088015280850151151561028088015261029660c08201516102a08901906101c9565b015163ffffffff166102c0860152565b6102d560208201516040860190604080916001600160801b038151168452602081015160208501520151910152565b610303604082015183860190604080916001600160801b038151168452602081015160208501520151910152565b60608101516001600160801b031661010085015261033f60808201516101208601906001600160401b0360208092828151168552015116910152565b01519101906001600160401b0360208092828151168552015116910152565b90565b634e487b7160e01b5f52604160045260245ffd5b60405190604082018281106001600160401b0382111761039457604052565b610361565b60405190606082018281106001600160401b0382111761039457604052565b6040519060c082018281106001600160401b0382111761039457604052565b6040519061010082018281106001600160401b0382111761039457604052565b60405190608082018281106001600160401b0382111761039457604052565b6040519061026082018281106001600160401b0382111761039457604052565b6040519190601f01601f191682016001600160401b0381118382101761039457604052565b610463610375565b905f82525f6020830152565b610477610399565b905f82525f60208301525f6040830152565b6104916103b8565b9061049a6103d7565b606081526104a661045b565b60208201526104b361045b565b60408201525f60608201525f60808201525f60a08201525f60c08201525f60e082015282526104e061046f565b60208301526104ed61046f565b60408301525f606083015261050061045b565b608083015261050d61045b565b60a0830152565b90359061013e198136030182121561018d570190565b903590601e198136030182121561018d57018035906001600160401b03821161018d5760200191813603831361018d57565b6001600160401b0381160361018d57565b3561035e8161055c565b6001600160401b03811161039457601f01601f191660200190565b9291926105a66105a183610577565b610436565b938285528282011161018d57815f926020928387013784010152565b63ffffffff81160361018d57565b3561035e816105c2565b359060ff8216820361018d57565b919082604091031261018d576106136020610601610375565b9361060b816105da565b8552016105da565b6020830152565b903590605e198136030182121561018d570190565b6001600160801b0381160361018d57565b3561035e8161062f565b35906106558261055c565b565b919082604091031261018d57602061066d610375565b9280356106798161055c565b845201356106138161055c565b9080601f8301121561018d5781602061035e93359101610592565b35906106558261062f565b3590811515820361018d57565b3590610655826105c2565b80929103916060831261018d5760406106db610375565b8235815293601f19011261018d5760406106f3610375565b916020810135610702816105c2565b8352013560208201526020830152565b6001600160401b0381116103945760051b60200190565b919060c08382031261018d5761073d6103f7565b9280356107498161055c565b84526020810135610759816105c2565b602085015261076b82604083016106c4565b604085015260a0810135906001600160401b03821161018d57019080601f8301121561018d578135916107a06105a184610712565b9260208085838152019160051b8301019183831161018d5760208101915b8383106107d15750505050506060830152565b82356001600160401b03811161018d578201906040828703601f19011261018d576107fa610375565b916020810135600481101561018d57835260408101356001600160401b03811161018d5760209101019060808288031261018d576108366103f7565b9282356001600160401b03811161018d5788610853918501610686565b845260208301356108638161062f565b6020850152610874604084016106ac565b60408501526060830135936001600160401b03851161018d5761089c89602096879601610686565b6060820152838201528152019201916107be565b91909160608184031261018d576108c5610375565b9281356001600160401b03811161018d57820160408183031261018d576108ea610375565b9281356001600160401b03811161018d5782016102c08185031261018d57610910610416565b9061091b8582610657565b825260408101356001600160401b03811161018d578561093c918301610686565b602083015261094d6060820161064a565b604083015261095e608082016106a1565b606083015261096f60a082016106ac565b60808301526109818560c083016106c4565b60a083015261099361012082016106ac565b60c083015261014081013560e08301526109b061016082016106ac565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526109ff61022082016106ac565b6101c08301526102408101356101e0830152610a1e61026082016106ac565b6102008301526102808101356102208301526102a0810135906001600160401b03821161018d57610a5191869101610686565b61024082015284526020820135936001600160401b03851161018d57610a7e846106139660209501610729565b83820152865201610657565b919082606091031261018d576040610aa0610399565b928035610aac8161062f565b84526020810135602085015201356040830152565b9392919093610ad560208251510151610f7e565b602081018051906020840191610af38351516001600160401b031690565b906001600160401b0382166001600160401b03821603610c4057505051610b65906001600160401b0316845151604001516001600160401b031692610b48610b39610375565b6001600160401b039093168352565b610b5f602083019485906001600160401b03169052565b51611967565b610c0b5750610b749086611211565b610b7f81515161128e565b604060208351015101515103610b9a57610655945190611759565b60405163f492ef2b60e01b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f7463680000000000000000000000000000000000000000000000000000000000606482015280608481015b0390fd5b517f90f4dbed000000000000000000000000000000000000000000000000000000005f526001600160401b031660045260245ffd5b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f526001600160401b039081166004521660245260445ffd5b903590603e198136030182121561018d570190565b9035906102be198136030182121561018d570190565b3590600282101561018d57565b9190916101408184031261018d57610cc86103d7565b928135916001600160401b03831161018d57610d0d82610cf061012094610d5d968501610686565b8752610cff81602085016105e8565b602088015260608301610657565b6040860152610d1e60a082016106b9565b6060860152610d2f60c082016106b9565b6080860152610d4060e082016106ac565b60a0860152610d526101008201610ca5565b60c0860152016106b9565b60e0830152565b90610d6d610489565b506080820191610ebb610d856040610099868561061a565b91610dae610d91610375565b6001600160401b0386168152936001600160401b03166020850152565b610dfc610dd26060610099610dcc610dc68a8761061a565b80610c7a565b80610c8f565b610dec610ddd610375565b6001600160401b039097168752565b6001600160401b03166020860152565b610e176080610e11610dcc610dc6898661061a565b01610640565b946101c0610e40610dcc610dc6610200610e37610dcc610dc6888a61061a565b0135948661061a565b013590610e5d610e4e610399565b6001600160801b039098168852565b60208701526040860152610e718180610514565b94610ea1610e8160a08401610640565b92610e95610e8d6103b8565b983690610cb2565b88526020369101610a8a565b602087015260408601526001600160801b03166060850152565b608083015260a082015290565b610ed0610375565b90606082525f6020830152565b634e487b7160e01b5f52601160045260245ffd5b8015610efd575f190190565b610edd565b5f19810191908211610efd57565b91908203918211610efd57565b634e487b7160e01b5f52603260045260245ffd5b908151811015610f42570160200190565b610f1d565b9060018201809211610efd57565b6001019081600111610efd57565b6023019081602311610efd57565b91908201809211610efd57565b610f86610ec8565b5080518015908115611205575b506111dd575f19908051805b611159575b505f19821461113b577f30000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000061102a611004610ffe86610f47565b85610f31565b517fff000000000000000000000000000000000000000000000000000000000000001690565b161480611145575b61113b575f9061104183610f47565b915b81518310156110d35761106261105c6110048585610f31565b60f81c90565b60ff8116603081109081156110c8575b506110bb576001600160401b0360ff8192602f19011681600a8502160116911681106110a357600190920191611043565b509150506110af610375565b9081525f602082015290565b50509150506110af610375565b60399150115f611072565b915091600181119081159161112f575b506111075761035e906110f4610375565b9283526001600160401b03166020830152565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f6110e3565b90506110af610375565b506002611153838351610f10565b11611032565b7f2d000000000000000000000000000000000000000000000000000000000000006111b761119261100461118c85610f02565b86610f31565b7fff000000000000000000000000000000000000000000000000000000000000001690565b146111cb576111c590610ef1565b80610f9f565b6111d6919250610f02565b905f610fa4565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f610f93565b906001600160401b036020830151166001600160401b0360208301511603611237575050565b90610c0761127c925191516040519384937fc6913e9f000000000000000000000000000000000000000000000000000000008552604060048601526044850190610191565b83810360031901602485015290610191565b61035e906114696112ab6102406112b06112ab6020860151611b67565b611bf1565b936112b9611991565b946112c76112ab8351611c46565b6112d0876119af565b526112da866119bc565b526113096112ab6113046112f860408501516001600160401b031690565b6001600160401b031690565b611d5c565b611312866119cc565b526113326112ab61132d60608401516001600160801b031690565b611d90565b61133b866119dc565b52608081015115611514576113566112ab60a0830151611e53565b61135f866119ec565b5260c0810151156114ee5761137a6112ab60e0830151611f1f565b611383866119fc565b52610100810151156114c8576113a06112ab610120830151611f1f565b6113a986611a0c565b526113bb6112ab610140830151611f1f565b6113c486611a1c565b526113d66112ab610160830151611f1f565b6113df86611a2d565b526113f16112ab610180830151611f1f565b6113fa86611a3e565b5261140c6112ab6101a0830151611f1f565b61141586611a4f565b526101c0810151156114a2576114326112ab6101e0830151611f1f565b61143b86611a60565b526102008101511561147c576114586112ab610220830151611f1f565b61146186611a71565b520151611b67565b61147282611a82565b525f815191611f4c565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611458565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611432565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d6113a0565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d61137a565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611356565b6001600160801b03633b9aca00911602906001600160801b038216918203610efd57565b906001600160801b03809116911603906001600160801b038211610efd57565b1561158557565b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b156115f657565b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b906001600160801b03809116911601906001600160801b038211610efd57565b1561166157565b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b6001600160401b036001911601906001600160401b038211610efd57565b156116f057565b608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b9391929061178261177d611774602087015163ffffffff1690565b63ffffffff1690565b61153a565b83516001600160801b0316906001600160801b0382166001600160801b03851610918215611945575b50506118d95761177d61177460406118699661182e6118399561181261183f996001600160801b0361180a6117fe8f60606117f091515101516001600160801b031690565b93516001600160801b031690565b6001600160801b031690565b91161161157e565b60208b51510151602081519101209060208151910120146115ef565b015163ffffffff1690565b9061163a565b6001600160801b038061185f6060865151016001600160801b0390511690565b921691161061165a565b6001600160401b03604061189361188e602080860151016001600160401b0390511690565b6116cb565b9251510191816118ab846001600160401b0390511690565b91169182911611156118bb575050565b6118d36112f8610655936001600160401b0390511690565b146116e9565b60405163f492ef2b60e01b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152608490fd5b6001600160801b0391925061195b82918661155e565b92169116115f806117ab565b9061197191611aa7565b60038110156101d65760028114908115611989575090565b600191501490565b6101e09061199e82610436565b600e815291601f1901366020840137565b805115610f425760200190565b805160011015610f425760400190565b805160021015610f425760600190565b805160031015610f425760800190565b805160041015610f425760a00190565b805160051015610f425760c00190565b805160061015610f425760e00190565b805160071015610f42576101000190565b805160081015610f42576101200190565b805160091015610f42576101400190565b8051600a1015610f42576101600190565b8051600b1015610f42576101800190565b8051600c1015610f42576101a00190565b8051600d1015610f42576101c00190565b8051821015610f425760209160051b010190565b80516001600160401b03166001600160401b03611ace6112f885516001600160401b031690565b911681811015611ae057505050505f90565b1115611aed575050600290565b611b1f6112f86020611b10816001600160401b039501516001600160401b031690565b9401516001600160401b031690565b911681811015611b2f5750505f90565b1115611b3a57600290565b600190565b90611b4c6105a183610577565b8281528092611b5d601f1991610577565b0190602036910137565b805115611bbc57611b788151611ff3565b806001019081600111610efd5760019083510101809111610efd57611b9f611bb891611b3f565b91600a6020840153611bb281518461206e565b83612109565b5090565b50611bc76020610436565b5f80825236602083013790565b805191908290602001825e015f815290565b6040513d5f823e3d90fd5b805160018101809111610efd57611c0790611b3f565b805115610f4257611c3181611c246020945f8681960153826120df565b5060405191828092611bd4565b039060025afa15611c41575f5190565b611be6565b5f90611c5981516001600160401b031690565b6001600160401b038116611d37575b50611c916020820192611c856112f885516001600160401b031690565b80611d16575b50611b3f565b915f91611ca86112f882516001600160401b031690565b611cf3575b50611cc26112f882516001600160401b031690565b611ccb57505090565b611cec6112f8611cde611bb8948661201e565b92516001600160401b031690565b90836120a6565b611d0f919250611d086112f8611cde86612012565b90846120a6565b905f611cad565b90611d2b611d26611d3193611ff3565b610f55565b90610f71565b5f611c8b565b611d55919250611d50611d26916001600160401b031690565b611ff3565b905f611c68565b8015611bbc57611d6b81611ff3565b60010180600111610efd57611d82611bb891611b3f565b91600860208401538261206e565b6001600160801b0316633b9aca0081066001600160801b03633b9aca005f930416918215159182611e35575b6001600160801b0316908115159081611e1a575b611dd990611b3f565b935f93611dfe575b50611deb57505090565b611df8611bb8928461201e565b836120a6565b611e13919350600860208601536001856120a6565b915f611de1565b611e26611d2684611ff3565b810180911115611dd057610edd565b90506001600160801b03611e4b611d2685611ff3565b919050611dbc565b6020810151611e6863ffffffff825116611ff3565b90816001019182600111610efd57602301809211610efd57611ecf611e8f611bb893611b3f565b91600860208401536020611ec5611ebf611eb9611eb3611774865163ffffffff1690565b8761206e565b86612035565b8561204c565b9101519083612134565b50611bb28151611f19611eb9611efd611ef884611ef3611eee82611ff3565b610f63565b610f71565b611b3f565b96611f10611f0a89612062565b8961204c565b90519088612134565b856120a6565b611bb8611f2c6060610436565b6022815291600a6020840160403682375360206021840153600283612134565b9092919280840393808511610efd5760018514611fe25760015b8060011b9086821015611f795750611f66565b91929394955050820191828111610efd5782611f959185611f4c565b91611fa09293611f4c565b90611fab6080610436565b6041815260208101906060368337805115610f42576020935f936001611c3194536021830152604182015260405191828092611bd4565b5090611fef929350611a93565b5190565b906001915b60808110156120045750565b60019060071c920191611ff8565b60206008910153600190565b60208260109201015360018101809111610efd5790565b60208260129201015360018101809111610efd5790565b602082819201015360018101809111610efd5790565b6020600a910153600190565b91906021600193015b608082101561208b57906001929391530190565b600180916080607f85161781530193019060071c9092612077565b9092919083016020015b60808210156120c457906001929391530190565b600180916080607f85161781530193019060071c90926120b0565b908051918215612101576021602084930191015e60010180600111610efd5790565b505050600190565b90809291825192831561212d57839260208092019201015e8101809111610efd5790565b5050505090565b8160209193929301015260208101809111610efd579056fea164736f6c634300081c000a",
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

// UpdateClient is a free data retrieval call binding the contract method 0xc8f5eea0.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCaller) UpdateClient(opts *bind.CallOpts, msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	var out []interface{}
	err := _ContractUpdateClient.contract.Call(opts, &out, "updateClient", msg_)

	if err != nil {
		return *new(IUpdateClientMsgsUpdateClientOutput), err
	}

	out0 := *abi.ConvertType(out[0], new(IUpdateClientMsgsUpdateClientOutput)).(*IUpdateClientMsgsUpdateClientOutput)

	return out0, err

}

// UpdateClient is a free data retrieval call binding the contract method 0xc8f5eea0.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientSession) UpdateClient(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClient(&_ContractUpdateClient.CallOpts, msg_)
}

// UpdateClient is a free data retrieval call binding the contract method 0xc8f5eea0.
//
// Solidity: function updateClient(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32),(uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),uint128,uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],uint64[],uint32[],bool[]) msg_) pure returns(((string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8,uint32),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientCallerSession) UpdateClient(msg_ IUpdateClientMsgsMsgUpdateClient) (IUpdateClientMsgsUpdateClientOutput, error) {
	return _ContractUpdateClient.Contract.UpdateClient(&_ContractUpdateClient.CallOpts, msg_)
}
