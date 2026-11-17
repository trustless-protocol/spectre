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

// ISpectreClientMsgsBatchProof is an auto generated low-level Go binding around an user-defined struct.
type ISpectreClientMsgsBatchProof struct {
	Proof                  [8]*big.Int
	Commitments            [4]*big.Int
	CommitmentPok          [2]*big.Int
	Bucket                 uint16
	SignerIndices          []uint32
	PinnedValidatorIndices []uint32
	SignerPubkeys          [][32]byte
	Active                 []bool
}

// ISpectreClientMsgsMisbehaviour is an auto generated low-level Go binding around an user-defined struct.
type ISpectreClientMsgsMisbehaviour struct {
	Header1 ISpectreTendermintMsgsHeader
	Header2 ISpectreTendermintMsgsHeader
}

// ISpectreClientMsgsMsgSubmitMisbehaviour is an auto generated low-level Go binding around an user-defined struct.
type ISpectreClientMsgsMsgSubmitMisbehaviour struct {
	Misbehaviour           ISpectreClientMsgsMisbehaviour
	TrustedConsensusState1 ISpectreTendermintMsgsConsensusState
	TrustedConsensusState2 ISpectreTendermintMsgsConsensusState
	Time                   *big.Int
	Proof1                 ISpectreClientMsgsBatchProof
	Proof2                 ISpectreClientMsgsBatchProof
}

// ISpectreTendermintMsgsBlockCommit is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsBlockCommit struct {
	Height     uint64
	Round      uint32
	BlockId    ISpectreTendermintMsgsBlockId
	CommitSigs []ISpectreTendermintMsgsCommitSig
}

// ISpectreTendermintMsgsBlockHeader is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsBlockHeader struct {
	Version            ISpectreTendermintMsgsVersion
	ChainId            string
	Height             uint64
	Time               *big.Int
	HasLastBlockId     bool
	LastBlockId        ISpectreTendermintMsgsBlockId
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

// ISpectreTendermintMsgsBlockId is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsBlockId struct {
	HashData      [32]byte
	PartSetHeader ISpectreTendermintMsgsPartSetHeader
}

// ISpectreTendermintMsgsCommitSig is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsCommitSig struct {
	Flag             uint8
	ValidatorAddress [20]byte
}

// ISpectreTendermintMsgsConsensusState is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsConsensusState struct {
	Timestamp          *big.Int
	Root               [32]byte
	NextValidatorsHash [32]byte
}

// ISpectreTendermintMsgsHeader is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsHeader struct {
	SignedHeader  ISpectreTendermintMsgsSignedHeader
	TrustedHeight IICS02ClientMsgsHeight
}

// ISpectreTendermintMsgsPartSetHeader is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsPartSetHeader struct {
	Total    uint32
	HashData [32]byte
}

// ISpectreTendermintMsgsSignedHeader is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsSignedHeader struct {
	Header ISpectreTendermintMsgsBlockHeader
	Commit ISpectreTendermintMsgsBlockCommit
}

// ISpectreTendermintMsgsVersion is an auto generated low-level Go binding around an user-defined struct.
type ISpectreTendermintMsgsVersion struct {
	BlockVersion uint64
	AppVersion   uint64
}

// ContractMisbehaviourMetaData contains all meta data concerning the ContractMisbehaviour contract.
var ContractMisbehaviourMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"signatureVerifier\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyMisbehaviour\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.MsgSubmitMisbehaviour\",\"components\":[{\"name\":\"misbehaviour\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.Misbehaviour\",\"components\":[{\"name\":\"header1\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structISpectreTendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumISpectreTendermintMsgs.CommitSigFlag\"},{\"name\":\"validatorAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"header2\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structISpectreTendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumISpectreTendermintMsgs.CommitSigFlag\"},{\"name\":\"validatorAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}]},{\"name\":\"trustedConsensusState1\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"trustedConsensusState2\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof1\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.BatchProof\",\"components\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[4]\",\"internalType\":\"uint256[4]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pinnedValidatorIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]},{\"name\":\"proof2\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.BatchProof\",\"components\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[4]\",\"internalType\":\"uint256[4]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pinnedValidatorIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DirectCallNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InsufficientTrustingPeriod\",\"inputs\":[{\"name\":\"durationSinceConsensusState\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"trustingPeriod\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConsensusStateTimestamp\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MisbehaviourNotDetected\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedMisbehaviourHeaderHeights\",\"inputs\":[{\"name\":\"height1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"height2\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]}]",
	Bin: "0x60c034607457601f611f9d38819003918201601f19168301916001600160401b03831184841017607857808492602094604052833981010312607457516001600160a01b038116908190036074576080523060a052604051611f10908161008d823960805181610e17015260a0518160720152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c6335aff9ac14610024575f80fd5b346105a75760203660031901126105a75760043567ffffffffffffffff81116105a757806004019061014060031982360301126105a75773ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016301461057f576100a282806105ab565b916040516100ba816100b3816105f8565b0382610721565b602081519101206100f46100ed6100e36100dd6100d78880610743565b806105ab565b80610758565b604081019061076e565b36916107bd565b602081519101201461010f6100e36100dd6100d78780610743565b90911561055957505061012a6101258480610743565b610f50565b6020830161013b6101258286610743565b6101516100ed6100e36100dd6100d78880610743565b6020815191012061016e6100ed6100e36100dd6100d7868a610743565b60208151910120036104f357610195606061018f6100dd6100d78880610743565b0161083d565b67ffffffffffffffff806101b4606061018f6100dd6100d7888c610743565b1691160361048e5760406101d86101ce6100d78780610743565b6020810190610852565b013560406101ec6101ce6100d78589610743565b013514610466577f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1703549363ffffffff60405195610228876106d5565b60405161023481610705565b60ff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170154818116835260081c1660208201528752818116602088015260481c16604086015267ffffffffffffffff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1702541694604051906102b382610705565b604051965f7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1700546102e3816105c0565b808b52906001811690811561044257506001146103c8575b506103be9361038e6103b398979461012497948c6103216103c69e610398970382610721565b8452602084015261037f8b6103368780610743565b602482019361035e61034786610813565b9260e4850193838a61035887610813565b93610867565b608461036a8b8b610743565b93019661035861037989610813565b93610813565b6103898580610743565b610b25565b6103898484610743565b6103b96103a58280610743565b6103b36101048a0188610827565b90610c50565b610743565b930190610827565b005b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005f90815291507fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be5b81831061042757505088016020016103be6102fb565b6001818c602086819597969754920101520191019190610411565b60ff19166020808d019190915291151560051b8b0190910191506103be90506102fb565b7f93c0a5f9000000000000000000000000000000000000000000000000000000005f5260045ffd5b8367ffffffffffffffff6104c2606061018f6100dd6100d785976104bc8561018f6100dd6100d78c80610743565b97610743565b917f494c87bc000000000000000000000000000000000000000000000000000000005f52166004521660245260445ffd5b839061055561051f6100e36100dd6100d76105176100e36100dd6100d78a80610743565b969097610743565b61054360405195869563f6b6676b60e01b87526040600488015260448701916107f3565b848103600319016024860152916107f3565b0390fd5b61055560405192839263f6b6676b60e01b845260406004850152610543604485016105f8565b7f3921c703000000000000000000000000000000000000000000000000000000005f5260045ffd5b5f80fd5b903590603e19813603018212156105a7570190565b90600182811c921680156105ee575b60208310146105da57565b634e487b7160e01b5f52602260045260245ffd5b91607f16916105cf565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1700545f9291610626826105c0565b80825291600181169081156106b95750600114610641575050565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005f9081529293509091907fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be5b83831061069f575060209250010190565b60018160209294939454838587010152019101919061068e565b9050602093945060ff929192191683830152151560051b010190565b6060810190811067ffffffffffffffff8211176106f157604052565b634e487b7160e01b5f52604160045260245ffd5b6040810190811067ffffffffffffffff8211176106f157604052565b90601f8019910116810190811067ffffffffffffffff8211176106f157604052565b903590605e19813603018212156105a7570190565b9035906102be19813603018212156105a7570190565b903590601e19813603018212156105a7570180359067ffffffffffffffff82116105a7576020019181360383136105a757565b67ffffffffffffffff81116106f157601f01601f191660200190565b9291926107c9826107a1565b916107d76040519384610721565b8294818452818301116105a7578281602093845f960137010152565b908060209392818452848401375f828201840152601f01601f1916010190565b356001600160801b03811681036105a75790565b90359061025e19813603018212156105a7570190565b3567ffffffffffffffff811681036105a75790565b90359060be19813603018212156105a7570190565b936001600160801b039081169381169190633b9aca0080840482169190860416818110610ae557036001600160801b0381116109e9576001600160801b0363ffffffff602086015116911681811015610ab75750506108cf6100ed6100e36100dd88806105ab565b60208151910120905180519160208201928320146108f36100e36100dd89806105ab565b909115610a6757505050506001600160801b0361091e60806109186100dd88806105ab565b01610813565b1611156109fd5763ffffffff6040633b9aca009201511602906001600160801b0382169182036109e95701906001600160801b0382116109e9576001600160801b0361097360806109186100dd8585966105ab565b92169116101561097f57565b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b634e487b7160e01b5f52601160045260245ffd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b9061055591606460405195869563f6b6676b60e01b8752604060048801525180918160448901528388015e5f868201830152601f01601f19168501858103820160031901602487015201916107f3565b7fdb020436000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b837fe42a1980000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b35906001600160801b03821682036105a757565b90610b7a604067ffffffffffffffff9281518260208201926001600160801b03610b4e82610b11565b1684526020810135828401520135606082015260608152610b70608082610721565b519020930161083d565b165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f20548015610be557808203610bb7575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b903590601e19813603018212156105a7570180359067ffffffffffffffff82116105a757602001918160051b360383136105a757565b359081151582036105a757565b90610c5e6101ce83806105ab565b9167ffffffffffffffff80610c84606061018f6100dd610c7d8961083d565b96806105ab565b16911614610c918361083d565b9015610ee85750610ca18261083d565b91602081013563ffffffff81168091036105a75767ffffffffffffffff60405194610ccb866106d5565b16845260208401908152604080850192013582526101c08301359361ffff85168095036105a757610180610d03610220860186610c0d565b906040610d14610240890189610c0d565b94909882519a7fb710e364000000000000000000000000000000000000000000000000000000008c5260048c01526101008160248d0137608061010082016101248d0137016101a48a01376102806101e4890152816102848901527f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82116105a7576102c4908896979594939260051b80916102a48901378601826102a482016003196102a48a850301016102048a0152520193905f5b818110610ebd575050509367ffffffffffffffff84939281602097511661022486015251166102448401525161026483015203815f73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af1908115610eb2575f91610e77575b5015610e4f57565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90506020813d602011610eaa575b81610e9260209383610721565b810103126105a7575180151581036105a7575f610e47565b3d9150610e85565b6040513d5f823e3d90fd5b91969550919293602080600192610ed38a610c43565b15158152019701910191879596949392610dcb565b67ffffffffffffffff906390f4dbed60e01b5f521660045260245ffd5b359067ffffffffffffffff821682036105a757565b91908260409103126105a757604051610f3281610705565b6020610f4b818395610f4381610f05565b855201610f05565b910152565b6020610f6d610f686100ed6100e36100dd86806105ab565b6118e2565b0167ffffffffffffffff81511690602083019167ffffffffffffffff610f928461083d565b168103611874575067ffffffffffffffff610feb915116610fe6610fbe606061018f6100dd88806105ab565b9360405192610fcc84610705565b835267ffffffffffffffff60208401951685523690610f1a565b611b5a565b60038110156118605760028114908115611855575b5061183757506110136100dd82806105ab565b803603906102c082126105a75760405191610260830183811067ffffffffffffffff8211176106f1576040526110493683610f1a565b8352604082013567ffffffffffffffff81116105a757820136601f820112156105a75761107d9036906020813591016107bd565b602084015261108e60608301610f05565b604084015261109f60808301610b11565b60608401526110b060a08301610c43565b6080840152606060bf198201126105a75760408051916110cf83610705565b60c0840135835260df1901126105a7576040516110eb81610705565b60e083013563ffffffff811681036105a75781526101008301356020820152602082015260a08301526111216101208201610c43565b60c083015261014081013560e083015261113e6101608201610c43565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a083015261118d6102208201610c43565b6101c08301526102408101356101e08301526111ac6102608201610c43565b6102008301526102808101356102208301526102a08101359067ffffffffffffffff82116105a7570136601f820112156105a7576111f19036906020813591016107bd565b61024082015261120c6112076020830151611c0c565b611c6f565b60405191906101e061121e8185610721565b600e8452601f190136602085013781515f9067ffffffffffffffff8151168061181a575b50602081019067ffffffffffffffff825116806117ef575b5061126761129493611bda565b915f9167ffffffffffffffff8151166117c7575b5067ffffffffffffffff8151166117a7575b5050611c6f565b8351156116555760208401528251600110156116555760408301526112c961120767ffffffffffffffff604084015116611cc2565b8251600210156116555760608301526001600160801b03606082015116633b9aca006001600160801b03821604906001600160801b035f921690811515908161178b575b633b9aca006001600160801b039106168015158061175f575b61133261134695611bda565b935f93611743575b50611728575050611c6f565b825160031015611655576080830152608081015115155f146117025760a081015160208101519061137d63ffffffff835116611def565b80600101806001116109e95760238201106109e9576113a160236113da9201611bda565b926008602085015360206113d06113ca6113c463ffffffff855116600189611e52565b87611e25565b86611e3c565b9101519084611edd565b5081516113e681611def565b60230192836023116109e95761143c8261143661143061141461140f611442976114489a6118d5565b611bda565b96600a6020890153611427600189611e3c565b90519088611edd565b86611e25565b85611e52565b83611eb5565b50611c6f565b8251600410156116555760a083015260c0810151156116dc5761147161120760e0830151611cf8565b8251600510156116555760c0830152610100810151156116b65761149c611207610120830151611cf8565b8251600610156116555760e08301526114bc611207610140830151611cf8565b825160071015611655576101008301526114dd611207610160830151611cf8565b825160081015611655576101208301526114fe611207610180830151611cf8565b8251600910156116555761014083015261151f6112076101a0830151611cf8565b8251600a1015611655576101608301526101c0810151156116905761154b6112076101e0830151611cf8565b8251600b1015611655576101808301526102008101511561166957611577611207610220830151611cf8565b905b8251600c1015611655576102406112079161159b936101a08601520151611c0c565b8151600d1015611655576101ce6115c1836115c8936101c060409601525f815191611d2e565b93806105ab565b0135036115d157565b60846040517ff492ef2b00000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152fd5b634e487b7160e01b5f52603260045260245ffd5b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d90611579565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d61154b565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d61149c565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611471565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611448565b61173561173b9284611e0e565b83611e52565b505f8061128d565b61175891935060086020860153600185611e52565b915f61133a565b61176882611def565b600101806001116109e95761178361133291611346976118d5565b955050611326565b925061179682611def565b600101806001116109e9579261130d565b67ffffffffffffffff6117bd61173b9385611e0e565b9151169083611e52565b6117e891925067ffffffffffffffff90600860208601535116600184611e52565b905f61127b565b6117f890611def565b600101806001116109e95761181361126791611294956118d5565b935061125a565b611825919250611def565b600101806001116109e957905f611242565b67ffffffffffffffff9051166390f4dbed60e01b5f5260045260245ffd5b60019150145f611000565b634e487b7160e01b5f52602160045260245ffd5b67ffffffffffffffff906118878461083d565b907f55bace6f000000000000000000000000000000000000000000000000000000005f526004521660245260445ffd5b919082039182116109e957565b908151811015611655570160200190565b919082018092116109e957565b6040516118ee81610705565b606081525f60208201525080518015908115611b4e575b50611b26575f19908051805b611aae575b505f198214611aa05760018201908183116109e9577f30000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000061197784846118c4565b51161480611a8c575b611a7c575f5b8151831015611a115761199983836118c4565b5160f81c603081108015611a07575b6119f55767ffffffffffffffff60ff8192602f19011681600a8502160116911681106119d957600190920191611986565b50915050604051906119ea82610705565b81525f602082015290565b5050915050604051906119ea82610705565b50603981116119a8565b9150916001811190811591611a70575b50611a485767ffffffffffffffff9060405192611a3d84610705565b835216602082015290565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f611a21565b915050604051906119ea82610705565b506002611a9a8483516118b7565b11611980565b60405191506119ea82610705565b5f1981018181116109e9577f2d000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000611b0583866118c4565b511614611b1c575080156109e9575f190180611911565b92505f9050611916565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f611905565b67ffffffffffffffff81511667ffffffffffffffff835116908181105f14611b8457505050505f90565b1115611b91575050600290565b602067ffffffffffffffff81819301511692015116908181105f14611bb65750505f90565b1115611bc157600290565b600190565b80518210156116555760209160051b010190565b90611be4826107a1565b611bf16040519182610721565b8281528092611c02601f19916107a1565b0190602036910137565b805115611c5357611c1d8151611def565b600101806001116109e957611c3a61140f611c4f928451906118d5565b91600a602084015361143c8151600185611e52565b5090565b50604051611c62602082610721565b5f80825236602083013790565b8051600181018091116109e957611c8590611bda565b8051156116555760209181611ca1845f94019284845382611e8b565b50604051918291518091835e8101838152039060025afa15610eb2575f5190565b8015611c5357611cd181611def565b600101806001116109e957611ce8611c4f91611bda565b9160086020840153600183611e52565b611c4f60405191611d0a606084610721565b602283526040366020850137600a6020840153611d28600184611e3c565b83611edd565b929192611d3b82856118b7565b9360018514611ddf5760015b8060011b9086821015611d5a5750611d47565b939495505090611d74611d6d84866118d5565b8583611d2e565b611d8893611d8291956118d5565b90611d2e565b90604051611d97608082610721565b6041815260208101906060368337805115611655576020935f936001845360218301526041820152604051918291518091835e8101838152039060025afa15610eb2575f5190565b50611deb929350611bc6565b5190565b906001915b6080811015611e005750565b60019060071c920191611df4565b602082601092010153600181018091116109e95790565b602082601292010153600181018091116109e95790565b6020828192010153600181018091116109e95790565b9092919083016020015b6080821015611e7057906001929391530190565b600180916080607f85161781530193019060071c9092611e5c565b908051918215611ead576021602084930191015e600101806001116109e95790565b505050600190565b825191928215611ed75783916020611ed495818694019201015e6118d5565b90565b50505090565b81602091939293010152602081018091116109e9579056fea26c6578706572696d656e74616cf564736f6c63430008250018",
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

// VerifyMisbehaviour is a paid mutator transaction binding the contract method 0x35aff9ac.
//
// Solidity: function verifyMisbehaviour(((((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,bytes20)[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,bytes20)[])),(uint64,uint64))),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint256[8],uint256[4],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[]),(uint256[8],uint256[4],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns()
func (_ContractMisbehaviour *ContractMisbehaviourTransactor) VerifyMisbehaviour(opts *bind.TransactOpts, msg_ ISpectreClientMsgsMsgSubmitMisbehaviour) (*types.Transaction, error) {
	return _ContractMisbehaviour.contract.Transact(opts, "verifyMisbehaviour", msg_)
}

// VerifyMisbehaviour is a paid mutator transaction binding the contract method 0x35aff9ac.
//
// Solidity: function verifyMisbehaviour(((((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,bytes20)[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,bytes20)[])),(uint64,uint64))),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint256[8],uint256[4],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[]),(uint256[8],uint256[4],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns()
func (_ContractMisbehaviour *ContractMisbehaviourSession) VerifyMisbehaviour(msg_ ISpectreClientMsgsMsgSubmitMisbehaviour) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.VerifyMisbehaviour(&_ContractMisbehaviour.TransactOpts, msg_)
}

// VerifyMisbehaviour is a paid mutator transaction binding the contract method 0x35aff9ac.
//
// Solidity: function verifyMisbehaviour(((((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,bytes20)[])),(uint64,uint64)),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,bytes20)[])),(uint64,uint64))),(uint128,bytes32,bytes32),(uint128,bytes32,bytes32),uint128,(uint256[8],uint256[4],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[]),(uint256[8],uint256[4],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns()
func (_ContractMisbehaviour *ContractMisbehaviourTransactorSession) VerifyMisbehaviour(msg_ ISpectreClientMsgsMsgSubmitMisbehaviour) (*types.Transaction, error) {
	return _ContractMisbehaviour.Contract.VerifyMisbehaviour(&_ContractMisbehaviour.TransactOpts, msg_)
}
