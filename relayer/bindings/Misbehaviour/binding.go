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

// IICS07TendermintMsgsCommitSig is an auto generated low-level Go binding around an user-defined struct.
type IICS07TendermintMsgsCommitSig struct {
	Flag uint8
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
	Header1 IICS07TendermintMsgsHeader
	Header2 IICS07TendermintMsgsHeader
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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"signatureVerifier\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyMisbehaviour\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.MsgSubmitMisbehaviour\",\"components\":[{\"name\":\"misbehaviour\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.Misbehaviour\",\"components\":[{\"name\":\"header1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"header2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}]},{\"name\":\"trustedConsensusState1\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"trustedConsensusState2\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof1\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.BatchProof\",\"components\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pinnedValidatorIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]},{\"name\":\"proof2\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.BatchProof\",\"components\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pinnedValidatorIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DirectCallNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MisbehaviourNotDetected\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]}]",
	Bin: "0x60c034607457601f6123a438819003918201601f19168301916001600160401b03831184841017607857808492602094604052833981010312607457516001600160a01b038116908190036074576080523060a052604051612317908161008d823960805181611072015260a0518160740152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c6360312a7c14610024575f80fd5b346100a65760203660031901126100a65760043567ffffffffffffffff81116100a65761014060031982360301126100a6576100a49061009c73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163014156100aa565b600401610509565b005b5f80fd5b156100b157565b7f3921c703000000000000000000000000000000000000000000000000000000005f5260045ffd5b903590603e19813603018212156100a6570190565b90600182811c9216801561011c575b602083101461010857565b634e487b7160e01b5f52602260045260245ffd5b91607f16916100fd565b5f9291815491610135836100ee565b808352926001811690811561018a575060011461015157505050565b5f9081526020812093945091925b838310610170575060209250010190565b60018160209294939454838587010152019101919061015f565b915050602093945060ff929192191683830152151560051b010190565b634e487b7160e01b5f52604160045260245ffd5b6040810190811067ffffffffffffffff8211176101d757604052565b6101a7565b90601f8019910116810190811067ffffffffffffffff8211176101d757604052565b6040519061023782610230817f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1700610126565b03836101dc565b565b903590605e19813603018212156100a6570190565b9035906102be19813603018212156100a6570190565b903590601e19813603018212156100a6570180359067ffffffffffffffff82116100a6576020019181360383136100a657565b604051906102376060836101dc565b604051906102376040836101dc565b60405190610237610260836101dc565b67ffffffffffffffff81116101d757601f01601f191660200190565b9291926102ed826102c5565b916102fb60405193846101dc565b8294818452818301116100a6578281602093845f960137010152565b908060209392818452848401375f828201840152601f01601f1916010190565b91909115610343575050565b61039c60405192839263f6b6676b60e01b84526040600485015261038a604485017f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1700610126565b84810360031901602486015291610317565b0390fd5b604051906103ad826101bb565b81602060ff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170154818116845260081c16910152565b604051905f827f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005491610414836100ee565b80835292600181169081156104b65750600114610438575b610237925003836101dc565b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005f90815290917fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be5b81831061049a5750509060206102379282010161042c565b6020919350806001915483858901015201910190918492610482565b6020925061023794915060ff191682840152151560051b82010161042c565b6001600160801b038116036100a657565b356104f0816104d5565b90565b90359061021e19813603018212156100a6570190565b610237906106e96106f461051d83806100d9565b6105836105286101fe565b6020815191012061056261055b61055161054b6105458780610239565b806100d9565b8061024e565b6040810190610264565b36916102e1565b602081519101201461057d61055161054b6105458680610239565b91610337565b61058c81610758565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170354906105f96105c963ffffffff84169360481c63ffffffff1690565b6105ec6105d4610297565b946105dd6103a0565b865263ffffffff166020860152565b63ffffffff166040840152565b6106ce61062e7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025467ffffffffffffffff1690565b6106556106396102a6565b916106426103e2565b835267ffffffffffffffff166020830152565b6106c46106628480610239565b946106b560208a019161068e610677846104e6565b9860e08d019983886106888d6104e6565b93610adc565b602087019761069d8989610239565b9160808d01966106886106af896104e6565b936104e6565b6106bf8580610239565b610cb7565b6106bf8484610239565b6106ef6106db8280610239565b6106e96101008801886104f3565b90610f47565b610239565b916101208101906104f3565b9290610719906104f09593604086526040860191610317565b926020818503910152610317565b67ffffffffffffffff8116036100a657565b356104f081610727565b90359060be19813603018212156100a6570190565b61076a6107658280610239565b6112f3565b6020810161077b6107658284610239565b61079161055b61055161054b6105458680610239565b602081519101206107ae61055b61055161054b6105458688610239565b60208151910120036108d5576107d560606107cf61054b6105458680610239565b01610739565b67ffffffffffffffff6108036107f660606107cf61054b610545888a610239565b67ffffffffffffffff1690565b91160361086a57610839610827610545604093846108316108276105458980610239565b6020810190610743565b013595610239565b01351461084257565b7f93c0a5f9000000000000000000000000000000000000000000000000000000005f5260045ffd5b9061089760606107cf61054b6105456108d296610891856107cf61054b6105458b80610239565b96610239565b7f494c87bc000000000000000000000000000000000000000000000000000000005f5267ffffffffffffffff91821660045216602452604490565b5ffd5b9061039c61090061055161054b6105456108f861055161054b6105458980610239565b979096610239565b9060405194859463f6b6676b60e01b865260048601610700565b634e487b7160e01b5f52601160045260245ffd5b906001600160801b03809116911603906001600160801b03821161094e57565b61091a565b9290921561096057505050565b6020929161039c91606460405195869563f6b6676b60e01b87526040600488015280519182918260448a0152018388015e5f868201830152601f01601f1916850185810382016003190160248701520191610317565b156109bd57565b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b6001600160801b03633b9aca00911602906001600160801b03821691820361094e57565b906001600160801b03809116911601906001600160801b03821161094e57565b15610a7257565b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b929493919094633b9aca00610afb610af482866114c8565b91836114c8565b906001600160801b0382166001600160801b03821610610c405790610b1f9161092e565b602083015163ffffffff16906001600160801b038116821115610c035750506080610bc361054b610bf26001600160801b0396610bec610be7610bde60408b9a6102379e9f610bf99b8f8e610bc98e610bc361054b85610545610bd399610b9961055b610b9061054b868c9d6100d9565b8f810190610264565b602081519101209051908151602083012014610bbb610b9061054b86806100d9565b929091610953565b016104e6565b92169116116109b6565b015163ffffffff1690565b63ffffffff1690565b610a27565b90610a4b565b96806100d9565b9216911610610a6b565b7fdb020436000000000000000000000000000000000000000000000000000000005f526001600160801b031660045263ffffffff1660245260445ffd5b7fe42a1980000000000000000000000000000000000000000000000000000000005f526001600160801b03831660045260245ffd5b3590610237826104d5565b15610c89575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b90604067ffffffffffffffff9181518260208201926001600160801b038135610cdf816104d5565b1684526020810135828401520135606082015260608152610d016080826101dc565b519020920135610d1081610727565b165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f2054908115610d4f5781610237928214610c80565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b15610d7f5750565b67ffffffffffffffff906390f4dbed60e01b5f521660045260245ffd5b63ffffffff8116036100a657565b356104f081610d9c565b3561ffff811681036100a65790565b903590601e19813603018212156100a6570180359067ffffffffffffffff82116100a657602001918160051b360383136100a657565b801515036100a657565b908160209103126100a657516104f081610df9565b359061023782610df9565b9998979593919261010060409461ffff8694168d5260208d01376101208b01376101608901376102406101a0880152806102408801527f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81116100a6576020916102609160051b8091838a0137870181888203016101c0890152018281520191905f5b818110610eea57505050906101e08261023793509401906040809167ffffffffffffffff815116845267ffffffffffffffff60208201511660208501520151910152565b9091926020806001928635610efe81610df9565b15158152019401929101610ea6565b6040513d5f823e3d90fd5b15610f1f57565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90602090610f97610f64610f5b85806100d9565b84810190610743565b9367ffffffffffffffff610f856107f660606107cf61054b610bf28b610739565b911614610f9185610739565b90610d77565b610fa083610739565b92610fe36040610fb4610bde868501610daa565b92013591610fd3610fc3610297565b67ffffffffffffffff9097168752565b67ffffffffffffffff1685850152565b6040840152610ff56101808201610db4565b6110586110066101e0840184610dc3565b611017610200869893980186610dc3565b9160405198899788977f7d1a869500000000000000000000000000000000000000000000000000000000895261014082019161010081019160048b01610e23565b03815f73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af180156110da57610237915f916110ab575b50610f18565b6110cd915060203d6020116110d3575b6110c581836101dc565b810190610e03565b5f6110a5565b503d6110bb565b610f0d565b359061023782610727565b91908260409103126100a657604051611102816101bb565b6020808294803561111281610727565b845201359161112083610727565b0152565b9080601f830112156100a6578160206104f0933591016102e1565b8092910391606083126100a657604051611158816101bb565b6040819483358352601f1901126100a657602090604080519361117a856101bb565b8381013561118781610d9c565b85520135828401520152565b6102c0813603126100a6576111a66102b5565b906111b136826110ea565b8252604081013567ffffffffffffffff81116100a6576111d49036908301611124565b60208301526111e5606082016110df565b60408301526111f660808201610c75565b606083015261120760a08201610e18565b60808301526112193660c0830161113f565b60a083015261122b6101208201610e18565b60c083015261014081013560e08301526112486101608201610e18565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a08301526112976102208201610e18565b6101c08301526102408101356101e08301526112b66102608201610e18565b6102008301526102808101356102208301526102a08101359067ffffffffffffffff82116100a6576112ea91369101611124565b61024082015290565b602061131061130b61055b61055161054b86806100d9565b61159b565b01805167ffffffffffffffff1690602083019161132f6107f684610739565b67ffffffffffffffff8216036114805750516113a49067ffffffffffffffff1661139f61136460606107cf61054b88806100d9565b936113806113706102a6565b67ffffffffffffffff9094168452565b6113986020840195869067ffffffffffffffff169052565b36906110ea565b611830565b611463575060406113d56108276113ce6113c96113c461054b87806100d9565b611193565b61186e565b93806100d9565b0135036113de57565b6040517ff492ef2b00000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152608490fd5b516390f4dbed60e01b5f5267ffffffffffffffff1660045260245ffd5b6108d29061148d84610739565b7f55bace6f000000000000000000000000000000000000000000000000000000005f5267ffffffffffffffff91821660045216602452604490565b906001600160801b03169081156114e6576001600160801b03160490565b634e487b7160e01b5f52601260045260245ffd5b60405190611507826101bb565b5f602083606081520152565b801561094e575f190190565b5f1981019190821161094e57565b9190820391821161094e57565b634e487b7160e01b5f52603260045260245ffd5b90815181101561155f570160200190565b61153a565b906001820180921161094e57565b600101908160011161094e57565b602301908160231161094e57565b9190820180921161094e57565b6115a36114fa565b5080518015908115611824575b506117fc575f19908051805b611778575b505f19821461175a577f30000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000061164761162161161b86611564565b8561154e565b517fff000000000000000000000000000000000000000000000000000000000000001690565b161480611764575b61175a575f9061165e83611564565b915b81518310156116f15761167f611679611621858561154e565b60f81c90565b60ff8116603081109081156116e6575b506116d95767ffffffffffffffff60ff8192602f19011681600a8502160116911681106116c157600190920191611660565b509150506116cd6102a6565b9081525f602082015290565b50509150506116cd6102a6565b60399150115f61168f565b915091600181119081159161174e575b50611726576104f0906117126102a6565b92835267ffffffffffffffff166020830152565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f611701565b90506116cd6102a6565b50600261177283835161152d565b1161164f565b7f2d000000000000000000000000000000000000000000000000000000000000006117d66117b16116216117ab8561151f565b8661154e565b7fff000000000000000000000000000000000000000000000000000000000000001690565b146117ea576117e490611513565b806115bc565b6117f591925061151f565b905f6115c1565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f6115b0565b9061183a91611b0f565b600381101561185a5760028114908115611852575090565b600191501490565b634e487b7160e01b5f52602160045260245ffd5b6104f090611a3e61188b61024061189061188b6020860151611d37565b611da7565b93611899611bad565b946118a761188b8351611df7565b6118b087611bd0565b526118ba86611bdd565b526118de61188b6118d96107f6604085015167ffffffffffffffff1690565b611f14565b6118e786611bed565b5261190761188b61190260608401516001600160801b031690565b611f48565b61191086611bfd565b52608081015115611ae95761192b61188b60a0830151612013565b61193486611c0d565b5260c081015115611ac35761194f61188b60e08301516120df565b61195886611c1d565b5261010081015115611a9d5761197561188b6101208301516120df565b61197e86611c2d565b5261199061188b6101408301516120df565b61199986611c3d565b526119ab61188b6101608301516120df565b6119b486611c4e565b526119c661188b6101808301516120df565b6119cf86611c5f565b526119e161188b6101a08301516120df565b6119ea86611c70565b526101c081015115611a7757611a0761188b6101e08301516120df565b611a1086611c81565b5261020081015115611a5157611a2d61188b6102208301516120df565b611a3686611c92565b520151611d37565b611a4782611ca3565b525f815191612117565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611a2d565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611a07565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611975565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d61194f565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d61192b565b805167ffffffffffffffff1667ffffffffffffffff611b396107f6855167ffffffffffffffff1690565b911681811015611b4b57505050505f90565b1115611b58575050600290565b611b8d6107f66020611b7d8167ffffffffffffffff95015167ffffffffffffffff1690565b94015167ffffffffffffffff1690565b911681811015611b9d5750505f90565b1115611ba857600290565b600190565b6040516101e09190611bbf83826101dc565b600e815291601f1901366020840137565b80511561155f5760200190565b80516001101561155f5760400190565b80516002101561155f5760600190565b80516003101561155f5760800190565b80516004101561155f5760a00190565b80516005101561155f5760c00190565b80516006101561155f5760e00190565b80516007101561155f576101000190565b80516008101561155f576101200190565b80516009101561155f576101400190565b8051600a101561155f576101600190565b8051600b101561155f576101800190565b8051600c101561155f576101a00190565b8051600d101561155f576101c00190565b805182101561155f5760209160051b010190565b60405190611cd76020836101dc565b5f808352366020840137565b60405160809190611cf483826101dc565b6041815291601f1901366020840137565b90611d0f826102c5565b611d1c60405191826101dc565b8281528092611d2d601f19916102c5565b0190602036910137565b805115611d8c57611d4881516121b1565b80600101908160011161094e576001908351010180911161094e57611d6f611d8891611d05565b91600a6020840153611d8281518461222c565b836122c7565b5090565b506104f0611cc8565b805191908290602001825e015f815290565b80516001810180911161094e57611dbd90611d05565b80511561155f57611de781611dda6020945f86819601538261229d565b5060405191828092611d95565b039060025afa156110da575f5190565b5f90611e0b815167ffffffffffffffff1690565b67ffffffffffffffff8116611eee575b50611e456020820192611e396107f6855167ffffffffffffffff1690565b80611ecd575b50611d05565b915f91611e5d6107f6825167ffffffffffffffff1690565b611eaa575b50611e786107f6825167ffffffffffffffff1690565b611e8157505090565b611ea36107f6611e94611d8894866121dc565b925167ffffffffffffffff1690565b9083612264565b611ec6919250611ebf6107f6611e94866121d0565b9084612264565b905f611e62565b90611ee2611edd611ee8936121b1565b611572565b9061158e565b5f611e3f565b611f0d919250611f08611edd9167ffffffffffffffff1690565b6121b1565b905f611e1b565b8015611d8c57611f23816121b1565b6001018060011161094e57611f3a611d8891611d05565b91600860208401538261222c565b6001600160801b0316633b9aca008104906001600160801b035f9216918215159182611fef575b633b9aca006001600160801b03910616908115159081611fd4575b611f9390611d05565b935f93611fb8575b50611fa557505090565b611fb2611d8892846121dc565b83612264565b611fcd91935060086020860153600185612264565b915f611f9b565b611fe0611edd846121b1565b810180911115611f8a5761091a565b90506001600160801b03633b9aca0061200a611edd866121b1565b92915050611f6f565b602081015161202863ffffffff8251166121b1565b9081600101918260011161094e5760230180921161094e5761208f61204f611d8893611d05565b9160086020840153602061208561207f612079612073610bde865163ffffffff1690565b8761222c565b866121f3565b8561220a565b91015190836122f2565b50611d8281516120d96120796120bd6120b8846120b36120ae826121b1565b611580565b61158e565b611d05565b966120d06120ca89612220565b8961220a565b905190886122f2565b85612264565b604051906060906120f082846101dc565b60228352611d8891600a906020850190601f190136823753602060218401536002836122f2565b909291928084039380851161094e57600185146121a05760015b8060011b90868210156121445750612131565b9192939495505082019182811161094e57826121609185612117565b9161216b9293612117565b612173611ce3565b9182511561155f57825f92611de79260016020809701536021830152604182015260405191828092611d95565b50906121ad929350611cb4565b5190565b906001915b60808110156121c25750565b60019060071c9201916121b6565b60206008910153600190565b6020826010920101536001810180911161094e5790565b6020826012920101536001810180911161094e5790565b60208281920101536001810180911161094e5790565b6020600a910153600190565b91906021600193015b608082101561224957906001929391530190565b600180916080607f85161781530193019060071c9092612235565b9092919083016020015b608082101561228257906001929391530190565b600180916080607f85161781530193019060071c909261226e565b9080519182156122bf576021602084930191015e6001018060011161094e5790565b505050600190565b9080929182519283156122eb57839260208092019201015e810180911161094e5790565b5050505090565b816020919392930101526020810180911161094e579056fea164736f6c634300081c000a",
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

// VerifyMisbehaviour is a paid mutator transaction binding the contract method 0x60312a7c.
//
// Solidity: function verifyMisbehaviour(((((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8)[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8)[])),(uint64,uint64))),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[]),(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns()
func (_ContractMisbehaviour *ContractMisbehaviourTransactor) VerifyMisbehaviour(opts *bind.TransactOpts, msg_ ISpectreClientMsgsMsgSubmitMisbehaviour) (*types.Transaction, error) {
	return _ContractMisbehaviour.contract.Transact(opts, "verifyMisbehaviour", msg_)
}

// VerifyMisbehaviour is a paid mutator transaction binding the contract method 0x60312a7c.
//
// Solidity: function verifyMisbehaviour(((((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8)[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8)[])),(uint64,uint64))),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[]),(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns()
func (_ContractMisbehaviour *ContractMisbehaviourSession) VerifyMisbehaviour(msg_ ISpectreClientMsgsMsgSubmitMisbehaviour) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.VerifyMisbehaviour(&_ContractMisbehaviour.TransactOpts, msg_)
}

// VerifyMisbehaviour is a paid mutator transaction binding the contract method 0x60312a7c.
//
// Solidity: function verifyMisbehaviour(((((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8)[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8)[])),(uint64,uint64))),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[]),(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns()
func (_ContractMisbehaviour *ContractMisbehaviourTransactorSession) VerifyMisbehaviour(msg_ ISpectreClientMsgsMsgSubmitMisbehaviour) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.VerifyMisbehaviour(&_ContractMisbehaviour.TransactOpts, msg_)
}
