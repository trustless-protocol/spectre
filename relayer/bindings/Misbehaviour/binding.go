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

// IICS07TendermintMsgsVersion is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsVersion struct {
	BlockVersion uint64
	AppVersion   uint64
}

// ISpectreClientMsgsBatchProof is an auto generated low-level Go binding around an user-defined struct.
type ISpectreClientMsgsBatchProof struct {
	Proof                  [8]*big.Int
	Commitments            [2]*big.Int
	CommitmentPok          [2]*big.Int
	Bucket                 uint16
	SignerIndices          []uint32
	PinnedValidatorIndices []uint32
	SignerPubkeys          [][32]byte
	Active                 []bool
}

// ISpectreClientMsgsMisbehaviour is an auto generated low-level Go binding around an user-defined struct.
type ISpectreClientMsgsMisbehaviour struct {
	ClientId IICS07TendermintMsgsChainId
	Header1  IICS07TendermintMsgsHeader
	Header2  IICS07TendermintMsgsHeader
}

// ISpectreClientMsgsMisbehaviourOutput is an auto generated low-level Go binding around an user-defined struct.
type ISpectreClientMsgsMisbehaviourOutput struct {
	TrustedHeight1 IICS02ClientMsgsHeight
	TrustedHeight2 IICS02ClientMsgsHeight
}

// ISpectreClientMsgsMsgSubmitMisbehaviour is an auto generated low-level Go binding around an user-defined struct.
type ISpectreClientMsgsMsgSubmitMisbehaviour struct {
	Misbehaviour           ISpectreClientMsgsMisbehaviour
	TrustedConsensusState1 IICS07TendermintMsgsConsensusState
	TrustedConsensusState2 IICS07TendermintMsgsConsensusState
	Time                   *big.Int
	Proof1                 ISpectreClientMsgsBatchProof
	Proof2                 ISpectreClientMsgsBatchProof
}

// ContractMisbehaviourMetaData contains all meta data concerning the ContractMisbehaviour contract.
var ContractMisbehaviourMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"signatureVerifier\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyMisbehaviour\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.MsgSubmitMisbehaviour\",\"components\":[{\"name\":\"misbehaviour\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.Misbehaviour\",\"components\":[{\"name\":\"clientId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ChainId\",\"components\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"header1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"header2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}]},{\"name\":\"trustedConsensusState1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"trustedConsensusState2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof1\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.BatchProof\",\"components\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pinnedValidatorIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]},{\"name\":\"proof2\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.BatchProof\",\"components\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pinnedValidatorIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.MisbehaviourOutput\",\"components\":[{\"name\":\"trustedHeight1\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"trustedHeight2\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DirectCallNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MisbehaviourNotDetected\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]}]",
	Bin: "0x60c034607557601f6125bf38819003918201601f19168301916001600160401b03831184841017607957808492602094604052833981010312607557516001600160a01b038116908190036075576080523060a052604051612531908161008e82396080518161107a015260a0518161019c0152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c63789d9f3914610024575f80fd5b346100ba5760203660031901126100ba5760043567ffffffffffffffff81116100ba5761014060031982360301126100ba5761006c60809160040161006761015d565b610184565b6100b860206040519261009684825167ffffffffffffffff60208092828151168552015116910152565b0151604083019067ffffffffffffffff60208092828151168552015116910152565bf35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b6040810190811067ffffffffffffffff8211176100ee57604052565b6100be565b90601f8019910116810190811067ffffffffffffffff8211176100ee57604052565b604051906101246060836100f3565b565b604051906101246040836100f3565b60405190610124610260836100f3565b60405190610152826100d2565b5f6020838281520152565b6040519061016a826100d2565b81610173610145565b8152602061017f610145565b910152565b5073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163014610459576101cc8180610481565b906101d56104ce565b805190602001209160208101926101ec8483610481565b806101f6916105c1565b80610200916105d6565b6040810161020d916105ec565b36906102189261063b565b805190602001201461022990610671565b6102328161078e565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1703549163ffffffff83169260501c63ffffffff1661026e610115565b936102776106a0565b855263ffffffff16602085015263ffffffff1660408401527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025467ffffffffffffffff16906102c4610126565b936102cd6104ce565b855260208501926102e890849067ffffffffffffffff169052565b6102f28685610481565b6020830190610300826106f3565b908760e086019285610311856106f3565b9261031b94610aee565b604086019661032a8888610481565b916080860194610339866106f3565b91610343906106f3565b9261034d94610aee565b6103578786610481565b9061036191610cbf565b61036b8585610481565b9061037591610cbf565b61037f8584610481565b61038d610100830183610700565b61039691610f4f565b6103a08484610481565b9061012081016103af91610700565b6103b891610f4f565b805167ffffffffffffffff16936103cf9083610481565b6040016103db90610728565b6103e3610126565b67ffffffffffffffff909516855267ffffffffffffffff1660208501525167ffffffffffffffff169161041591610481565b60400161042190610728565b610429610126565b67ffffffffffffffff909216825267ffffffffffffffff16602082015261044e610126565b918252602082015290565b7f3921c703000000000000000000000000000000000000000000000000000000005f5260045ffd5b903590605e19813603018212156100ba570190565b90600182811c921680156104c4575b60208310146104b057565b634e487b7160e01b5f52602260045260245ffd5b91607f16916104a5565b604051905f827f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1700549161050083610496565b80835292600181169081156105a25750600114610524575b610124925003836100f3565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005f90815290917fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be5b81831061058657505090602061012492820101610518565b602091935080600191548385890101520191019091849261056e565b6020925061012494915060ff191682840152151560051b820101610518565b903590603e19813603018212156100ba570190565b9035906102be19813603018212156100ba570190565b903590601e19813603018212156100ba570180359067ffffffffffffffff82116100ba576020019181360383136100ba57565b67ffffffffffffffff81116100ee57601f01601f191660200190565b9291926106478261061f565b9161065560405193846100f3565b8294818452818301116100ba578281602093845f960137010152565b1561067857565b7fa179f8c9000000000000000000000000000000000000000000000000000000005f5260045ffd5b604051906106ad826100d2565b81602060ff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170154818116845260081c16910152565b6001600160801b038116036100ba57565b356106fd816106e2565b90565b90359061021e19813603018212156100ba570190565b67ffffffffffffffff8116036100ba57565b356106fd81610716565b908060209392818452848401375f828201840152601f01601f1916010190565b929061076b906106fd9593604086526040860191610732565b926020818503910152610732565b90359060be19813603018212156100ba570190565b602081016107a461079f8284610481565b6112fb565b60408201906107b661079f8385610481565b6107e96107e26107d86107d26107cc8588610481565b806105c1565b806105d6565b60408101906105ec565b369161063b565b602081519101206108066107e26107d86107d26107cc8789610481565b602081519101200361092d5761082d60606108276107d26107cc8588610481565b01610728565b67ffffffffffffffff61085b61084e60606108276107d26107cc898b610481565b67ffffffffffffffff1690565b9116036108c25761087f6107cc6040938461088961087f6107cc610891978a610481565b6020810190610779565b013595610481565b01351461089a57565b7f93c0a5f9000000000000000000000000000000000000000000000000000000005f5260045ffd5b60606108276107d26107cc86956108e9856108276107d26107cc6108ef9a61092a9d610481565b96610481565b7f494c87bc000000000000000000000000000000000000000000000000000000005f5267ffffffffffffffff91821660045216602452604490565b5ffd5b61098b6109586107d86107d26107cc6109506107d86107d26107cc8b998a610481565b979096610481565b906040519485947ff6b6676b00000000000000000000000000000000000000000000000000000000865260048601610752565b0390fd5b634e487b7160e01b5f52601160045260245ffd5b906001600160801b03809116911603906001600160801b0382116109c357565b61098f565b156109cf57565b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b6001600160801b03633b9aca00911602906001600160801b0382169182036109c357565b906001600160801b03809116911601906001600160801b0382116109c357565b15610a8457565b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b929493919094633b9aca00610b0d610b0682866114d0565b91836114d0565b906001600160801b0382166001600160801b03821610610c485790610b31916109a3565b602083015163ffffffff16906001600160801b038116821115610c0b5750506080610bcb6107d2610bfa6001600160801b0396610bf4610bef610be660408b9a6101249e9f610c019b8f8e610bd18e610bcb6107d284956107cc89610b99610bdb9b5161152c565b610bb56107e2610bac6107d286806105c1565b8f8101906105ec565b6020815191012090516020815191012014610671565b016106f3565b92169116116109c8565b015163ffffffff1690565b63ffffffff1690565b610a39565b90610a5d565b96806105c1565b9216911610610a7d565b7fdb020436000000000000000000000000000000000000000000000000000000005f526001600160801b031660045263ffffffff1660245260445ffd5b7fe42a1980000000000000000000000000000000000000000000000000000000005f526001600160801b03831660045260245ffd5b3590610124826106e2565b15610c91575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b90604067ffffffffffffffff9181518260208201926001600160801b038135610ce7816106e2565b1684526020810135828401520135606082015260608152610d096080826100f3565b519020920135610d1881610716565b165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2054908115610d575781610124928214610c88565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b15610d875750565b67ffffffffffffffff906390f4dbed60e01b5f521660045260245ffd5b63ffffffff8116036100ba57565b356106fd81610da4565b3561ffff811681036100ba5790565b903590601e19813603018212156100ba570180359067ffffffffffffffff82116100ba57602001918160051b360383136100ba57565b801515036100ba57565b908160209103126100ba57516106fd81610e01565b359061012482610e01565b9998979593919261010060409461ffff8694168d5260208d01376101208b01376101608901376102406101a0880152806102408801527f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81116100ba576020916102609160051b8091838a0137870181888203016101c0890152018281520191905f5b818110610ef257505050906101e08261012493509401906040809167ffffffffffffffff815116845267ffffffffffffffff60208201511660208501520151910152565b9091926020806001928635610f0681610e01565b15158152019401929101610eae565b6040513d5f823e3d90fd5b15610f2757565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90602090610f9f610f6c610f6385806105c1565b84810190610779565b9367ffffffffffffffff610f8d61084e60606108276107d2610bfa8b610728565b911614610f9985610728565b90610d7f565b610fa883610728565b92610feb6040610fbc610be6868501610db2565b92013591610fdb610fcb610115565b67ffffffffffffffff9097168752565b67ffffffffffffffff1685850152565b6040840152610ffd6101808201610dbc565b61106061100e6101e0840184610dcb565b61101f610200869893980186610dcb565b9160405198899788977f7d1a869500000000000000000000000000000000000000000000000000000000895261014082019161010081019160048b01610e2b565b03815f73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af180156110e257610124915f916110b3575b50610f20565b6110d5915060203d6020116110db575b6110cd81836100f3565b810190610e0b565b5f6110ad565b503d6110c3565b610f15565b359061012482610716565b91908260409103126100ba5760405161110a816100d2565b6020808294803561111a81610716565b845201359161112883610716565b0152565b9080601f830112156100ba578160206106fd9335910161063b565b8092910391606083126100ba57604051611160816100d2565b6040819483358352601f1901126100ba576020906040805193611182856100d2565b8381013561118f81610da4565b85520135828401520152565b6102c0813603126100ba576111ae610135565b906111b936826110f2565b8252604081013567ffffffffffffffff81116100ba576111dc903690830161112c565b60208301526111ed606082016110e7565b60408301526111fe60808201610c7d565b606083015261120f60a08201610e20565b60808301526112213660c08301611147565b60a08301526112336101208201610e20565b60c083015261014081013560e08301526112506101608201610e20565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a083015261129f6102208201610e20565b6101c08301526102408101356101e08301526112be6102608201610e20565b6102008301526102808101356102208301526102a08101359067ffffffffffffffff82116100ba576112f29136910161112c565b61024082015290565b60206113186113136107e26107d86107d286806105c1565b611838565b01805167ffffffffffffffff1690602083019161133761084e84610728565b67ffffffffffffffff8216036114885750516113ac9067ffffffffffffffff166113a761136c60606108276107d288806105c1565b93611388611378610126565b67ffffffffffffffff9094168452565b6113a06020840195869067ffffffffffffffff169052565b36906110f2565b611a4a565b61146b575060406113dd61087f6113d66113d16113cc6107d287806105c1565b61119b565b611a88565b93806105c1565b0135036113e657565b6040517ff492ef2b00000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152608490fd5b516390f4dbed60e01b5f5267ffffffffffffffff1660045260245ffd5b61092a9061149584610728565b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f5267ffffffffffffffff91821660045216602452604490565b906001600160801b03169081156114ee576001600160801b03160490565b634e487b7160e01b5f52601260045260245ffd5b634e487b7160e01b5f52603260045260245ffd5b908151811015611527570160200190565b611502565b805180159081156117b6575b50611771575f5b815181101561176d576115a361157e6115588385611516565b517fff000000000000000000000000000000000000000000000000000000000000001690565b7fff000000000000000000000000000000000000000000000000000000000000001690565b7f61000000000000000000000000000000000000000000000000000000000000008110159081611742575b81156116e6575b81156116a6575b8115611698575b811561166e575b8115611644575b50156115ff5760010161153f565b60405162461bcd60e51b815260206004820152601860248201527f696e76616c696420636861696e206964206368617273657400000000000000006044820152606490fd5b7f2e000000000000000000000000000000000000000000000000000000000000009150145f6115f1565b7f5f00000000000000000000000000000000000000000000000000000000000000811491506115ea565b602d60f81b811491506115e3565b9050600360fc1b811015806116bc575b906115dc565b507f39000000000000000000000000000000000000000000000000000000000000008111156116b6565b90507f410000000000000000000000000000000000000000000000000000000000000081101580611718575b906115d5565b507f5a00000000000000000000000000000000000000000000000000000000000000811115611712565b7f7a0000000000000000000000000000000000000000000000000000000000000081111591506115ce565b5050565b60405162461bcd60e51b815260206004820152601760248201527f496e76616c696420636861696e206964206c656e6774680000000000000000006044820152606490fd5b60329150115f611538565b604051906117ce826100d2565b5f602083606081520152565b80156109c3575f190190565b5f198101919082116109c357565b919082039182116109c357565b90600182018092116109c357565b60010190816001116109c357565b60230190816023116109c357565b919082018092116109c357565b6118406117c1565b5080518015908115611a3e575b50611a16575f19908051805b6119d3575b505f1982146119b557600360fc1b7fff000000000000000000000000000000000000000000000000000000000000006118a261155861189c86611801565b85611516565b1614806119bf575b6119b5575f906118b983611801565b915b815183101561194c576118da6118d46115588585611516565b60f81c90565b60ff811660308110908115611941575b506119345767ffffffffffffffff60ff8192602f19011681600a85021601169116811061191c576001909201916118bb565b50915050611928610126565b9081525f602082015290565b5050915050611928610126565b60399150115f6118ea565b91509160018111908115916119a9575b50611981576106fd9061196d610126565b92835267ffffffffffffffff166020830152565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f61195c565b9050611928610126565b5060026119cd8383516117f4565b116118aa565b602d60f81b6119f061157e6115586119ea856117e6565b86611516565b14611a04576119fe906117da565b80611859565b611a0f9192506117e6565b905f61185e565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f61184d565b90611a5491611d29565b6003811015611a745760028114908115611a6c575090565b600191501490565b634e487b7160e01b5f52602160045260245ffd5b6106fd90611c58611aa5610240611aaa611aa56020860151611f51565b611fc1565b93611ab3611dc7565b94611ac1611aa58351612011565b611aca87611dea565b52611ad486611df7565b52611af8611aa5611af361084e604085015167ffffffffffffffff1690565b61212e565b611b0186611e07565b52611b21611aa5611b1c60608401516001600160801b031690565b612162565b611b2a86611e17565b52608081015115611d0357611b45611aa560a083015161222d565b611b4e86611e27565b5260c081015115611cdd57611b69611aa560e08301516122f9565b611b7286611e37565b5261010081015115611cb757611b8f611aa56101208301516122f9565b611b9886611e47565b52611baa611aa56101408301516122f9565b611bb386611e57565b52611bc5611aa56101608301516122f9565b611bce86611e68565b52611be0611aa56101808301516122f9565b611be986611e79565b52611bfb611aa56101a08301516122f9565b611c0486611e8a565b526101c081015115611c9157611c21611aa56101e08301516122f9565b611c2a86611e9b565b5261020081015115611c6b57611c47611aa56102208301516122f9565b611c5086611eac565b520151611f51565b611c6182611ebd565b525f815191612331565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611c47565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611c21565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611b8f565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611b69565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611b45565b805167ffffffffffffffff1667ffffffffffffffff611d5361084e855167ffffffffffffffff1690565b911681811015611d6557505050505f90565b1115611d72575050600290565b611da761084e6020611d978167ffffffffffffffff95015167ffffffffffffffff1690565b94015167ffffffffffffffff1690565b911681811015611db75750505f90565b1115611dc257600290565b600190565b6040516101e09190611dd983826100f3565b600e815291601f1901366020840137565b8051156115275760200190565b8051600110156115275760400190565b8051600210156115275760600190565b8051600310156115275760800190565b8051600410156115275760a00190565b8051600510156115275760c00190565b8051600610156115275760e00190565b805160071015611527576101000190565b805160081015611527576101200190565b805160091015611527576101400190565b8051600a1015611527576101600190565b8051600b1015611527576101800190565b8051600c1015611527576101a00190565b8051600d1015611527576101c00190565b80518210156115275760209160051b010190565b60405190611ef16020836100f3565b5f808352366020840137565b60405160809190611f0e83826100f3565b6041815291601f1901366020840137565b90611f298261061f565b611f3660405191826100f3565b8281528092611f47601f199161061f565b0190602036910137565b805115611fa657611f6281516123cb565b8060010190816001116109c357600190835101018091116109c357611f89611fa291611f1f565b91600a6020840153611f9c815184612446565b836124e1565b5090565b506106fd611ee2565b805191908290602001825e015f815290565b8051600181018091116109c357611fd790611f1f565b8051156115275761200181611ff46020945f8681960153826124b7565b5060405191828092611faf565b039060025afa156110e2575f5190565b5f90612025815167ffffffffffffffff1690565b67ffffffffffffffff8116612108575b5061205f602082019261205361084e855167ffffffffffffffff1690565b806120e7575b50611f1f565b915f9161207761084e825167ffffffffffffffff1690565b6120c4575b5061209261084e825167ffffffffffffffff1690565b61209b57505090565b6120bd61084e6120ae611fa294866123f6565b925167ffffffffffffffff1690565b908361247e565b6120e09192506120d961084e6120ae866123ea565b908461247e565b905f61207c565b906120fc6120f7612102936123cb565b61180f565b9061182b565b5f612059565b6121279192506121226120f79167ffffffffffffffff1690565b6123cb565b905f612035565b8015611fa65761213d816123cb565b600101806001116109c357612154611fa291611f1f565b916008602084015382612446565b6001600160801b0316633b9aca008104906001600160801b035f9216918215159182612209575b633b9aca006001600160801b039106169081151590816121ee575b6121ad90611f1f565b935f936121d2575b506121bf57505090565b6121cc611fa292846123f6565b8361247e565b6121e79193506008602086015360018561247e565b915f6121b5565b6121fa6120f7846123cb565b8101809111156121a45761098f565b90506001600160801b03633b9aca006122246120f7866123cb565b92915050612189565b602081015161224263ffffffff8251166123cb565b908160010191826001116109c3576023018092116109c3576122a9612269611fa293611f1f565b9160086020840153602061229f61229961229361228d610be6865163ffffffff1690565b87612446565b8661240d565b85612424565b910151908361250c565b50611f9c81516122f36122936122d76122d2846122cd6122c8826123cb565b61181d565b61182b565b611f1f565b966122ea6122e48961243a565b89612424565b9051908861250c565b8561247e565b6040519060609061230a82846100f3565b60228352611fa291600a906020850190601f1901368237536020602184015360028361250c565b90929192808403938085116109c357600185146123ba5760015b8060011b908682101561235e575061234b565b919293949550508201918281116109c3578261237a9185612331565b916123859293612331565b61238d611efd565b9182511561152757825f926120019260016020809701536021830152604182015260405191828092611faf565b50906123c7929350611ece565b5190565b906001915b60808110156123dc5750565b60019060071c9201916123d0565b60206008910153600190565b602082601092010153600181018091116109c35790565b602082601292010153600181018091116109c35790565b6020828192010153600181018091116109c35790565b6020600a910153600190565b91906021600193015b608082101561246357906001929391530190565b600180916080607f85161781530193019060071c909261244f565b9092919083016020015b608082101561249c57906001929391530190565b600180916080607f85161781530193019060071c9092612488565b9080519182156124d9576021602084930191015e600101806001116109c35790565b505050600190565b90809291825192831561250557839260208092019201015e81018091116109c35790565b5050505090565b81602091939293010152602081018091116109c3579056fea164736f6c634300081c000a",
}

// ContractMisbehaviourABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMisbehaviourMetaData.ABI instead.
var ContractMisbehaviourABI = ContractMisbehaviourMetaData.ABI

// ContractMisbehaviourBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMisbehaviourMetaData.Bin instead.
var ContractMisbehaviourBin = ContractMisbehaviourMetaData.Bin

// DeployContractMisbehaviour deploys a new Ethereum contract, binding an instance of ContractMisbehaviour to it.
func DeployContractMisbehaviour(auth *bind.TransactOpts, backend bind.ContractBackend, signatureVerifier common.Address) (common.Address, *types.Transaction, *ContractMisbehaviour, error) {
	parsed, err := ContractMisbehaviourMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractMisbehaviourBin), backend, signatureVerifier)
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

// VerifyMisbehaviour is a paid mutator transaction binding the contract method 0x789d9f39.
//
// Solidity: function verifyMisbehaviour((((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64))),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[]),(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourTransactor) VerifyMisbehaviour(opts *bind.TransactOpts, msg_ ISpectreClientMsgsMsgSubmitMisbehaviour) (*types.Transaction, error) {
	return _ContractMisbehaviour.contract.Transact(opts, "verifyMisbehaviour", msg_)
}

// VerifyMisbehaviour is a paid mutator transaction binding the contract method 0x789d9f39.
//
// Solidity: function verifyMisbehaviour((((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64))),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[]),(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourSession) VerifyMisbehaviour(msg_ ISpectreClientMsgsMsgSubmitMisbehaviour) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.VerifyMisbehaviour(&_ContractMisbehaviour.TransactOpts, msg_)
}

// VerifyMisbehaviour is a paid mutator transaction binding the contract method 0x789d9f39.
//
// Solidity: function verifyMisbehaviour((((string,uint64),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64))),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[]),(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns(((uint64,uint64),(uint64,uint64)))
func (_ContractMisbehaviour *ContractMisbehaviourTransactorSession) VerifyMisbehaviour(msg_ ISpectreClientMsgsMsgSubmitMisbehaviour) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.VerifyMisbehaviour(&_ContractMisbehaviour.TransactOpts, msg_)
}
