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

// ISpectreClientMsgsMsgUpdateApplicationState is an auto generated low-level Go binding around an user-defined struct.
type ISpectreClientMsgsMsgUpdateApplicationState struct {
	TrustedConsensusState ISpectreTendermintMsgsConsensusState
	ProposedHeader        ISpectreTendermintMsgsHeader
	Time                  *big.Int
	Proof                 ISpectreClientMsgsBatchProof
}

// ISpectreClientMsgsVerifyHeaderOutput is an auto generated low-level Go binding around an user-defined struct.
type ISpectreClientMsgsVerifyHeaderOutput struct {
	TrustedConsensusState ISpectreTendermintMsgsConsensusState
	NewConsensusState     ISpectreTendermintMsgsConsensusState
	TrustedHeight         IICS02ClientMsgsHeight
	NewHeight             IICS02ClientMsgsHeight
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

// ContractUpdateClientMetaData contains all meta data concerning the ContractUpdateClient contract.
var ContractUpdateClientMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"signatureVerifier\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyHeader\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.MsgUpdateApplicationState\",\"components\":[{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structISpectreTendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumISpectreTendermintMsgs.CommitSigFlag\"},{\"name\":\"validatorAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.BatchProof\",\"components\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[4]\",\"internalType\":\"uint256[4]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pinnedValidatorIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.VerifyHeaderOutput\",\"components\":[{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"newConsensusState\",\"type\":\"tuple\",\"internalType\":\"structISpectreTendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"newHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"ChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DirectCallNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeights\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]}]",
	Bin: "0x60c034607457601f611ff538819003918201601f19168301916001600160401b03831184841017607857808492602094604052833981010312607457516001600160a01b038116908190036074576080523060a052604051611f68908161008d823960805181610c53015260a05181607e0152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c63c67440cc14610024575f80fd5b34610ea3576020366003190112610ea35760043567ffffffffffffffff8111610ea35780600401908036039160031983019160c08312610ea3576100666116c7565b5073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016301461160b5767ffffffffffffffff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1702541693604051946100e0604087611687565b6040515f7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1700548060011c9160018216918215611601575b6020841083146115ed5783855284929081156115ce5750600114611550575b61014292500382611687565b8652602086019081527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17035463ffffffff8082169160481c16966040805161018a606082611687565b81516101958161164f565b60ff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170154818116835260081c1660208201528152602081019384520197885260648501976101e38988611722565b926101f060848801611737565b9261022361021e61021761020d6102078980611774565b80611789565b604081019061179f565b36916117ee565b611916565b6020810167ffffffffffffffff815116602088019067ffffffffffffffff61024a8361175f565b16810361150e57506102aa886102a561027d606061027761020767ffffffffffffffff8951169580611774565b0161175f565b936040519261028b8461164f565b835267ffffffffffffffff6020840195168552369061185d565b611bc6565b60038110156114fa57600281149081156114ef575b506114d1575067ffffffffffffffff808951169151160361147757506102e86102078680611774565b8036036102c08112610ea35760405191610260830183811067ffffffffffffffff8211176114635760405261031d368261185d565b8352604081013567ffffffffffffffff8111610ea357810136601f82011215610ea3576103519036906020813591016117ee565b916020840192835261036560608301611848565b60408501526103766080830161174b565b606085015261038760a08301611893565b6080850152606060bf19820112610ea35760408051916103a68361164f565b60c0840135835260df190112610ea3576040516103c28161164f565b60e083013563ffffffff81168103610ea35781526101008301356020820152602082015260a08401526103f86101208201611893565b60c084015261014081013560e08401526104156101608201611893565b6101008401526101808101356101208401526101a08101356101408401526101c08101356101608401526101e08101356101808401526102008101356101a08401526104646102208201611893565b6101c08401526102408101356101e08401526104836102608201611893565b6102008401526102808101356102208401526102a08101359067ffffffffffffffff8211610ea3570136601f82011215610ea3576104e0916104cf6104db9236906020813591016117ee565b61024085015251611c64565b611cc7565b906101e0604051926104f28285611687565b600e845260208401601f19830136823783515f9067ffffffffffffffff81511680611446575b50602081019067ffffffffffffffff8251168061141b575b5061053d61056a93611c32565b915f9167ffffffffffffffff8151166113f3575b5067ffffffffffffffff8151166113d3575b5050611cc7565b90855115611280575283516001101561128057604084015261059c6104db67ffffffffffffffff604085015116611d1a565b8351600210156112805760608401526001600160801b036060830151165f906001600160801b03633b9aca008204169081151590816113b7575b633b9aca006001600160801b039106168015158061138b575b6105fb61060f95611c32565b935f9361136f575b50611354575050611cc7565b835160031015611280576080840152608082015115155f1461132e5760a082015160208101519061064663ffffffff835116611e47565b80600101908160011161104457602301809111611044576106696106a291611c32565b9260086020850153602061069861069261068c63ffffffff855116600189611eaa565b87611e7d565b86611e94565b9101519084611f35565b5081516106ae81611e47565b602301928360231161104457610704826106fe6106f86106dc6106d761070a976107109a611909565b611c32565b96600a60208901536106ef600189611e94565b90519088611f35565b86611e7d565b85611eaa565b83611f0d565b50611cc7565b8351600410156112805760a084015260c082015115611308576107396104db60e0840151611d50565b8351600510156112805760c0840152610100820151156112e2576107646104db610120840151611d50565b8351600610156112805760e08401526107846104db610140840151611d50565b835160071015611280576101008401526107a56104db610160840151611d50565b835160081015611280576101208401526107c66104db610180840151611d50565b835160091015611280576101408401526107e76104db6101a0840151611d50565b8351600a1015611280576101608401526101c0820151156112bb576104db61081191830151611d50565b8251600b101561128057610180830152610200810151156112945761083d6104db610220830151611d50565b905b8251600c1015611280576102406104db91610861936101a08601520151611c64565b908051600d10156112805761087f916101c08201525f815191611d86565b604061089861088e8880611774565b60208101906118a0565b0135036112155763ffffffff6108b19151925116611b8e565b6001600160801b03806108c38b611737565b1694169384109081156111db575b50611170576108ee60806108e86102078780611774565b01611737565b6001600160801b03806109008b611737565b16911611156111065761091c61021761020d6102078780611774565b602081519101209060208151910120036110c25761094863ffffffff6001600160801b03925116611b8e565b16016001600160801b038111611044576001600160801b038061097360806108e86102078780611774565b92169116101561105857600167ffffffffffffffff6109946040840161175f565b16019067ffffffffffffffff82116110445767ffffffffffffffff91826109c360606102776102078680611774565b9116928391161115610fb3575b50506040519160208301936001600160801b036109ec8761174b565b16855260248101359485604086015260448201359485606082015260608152610a16608082611687565b51902067ffffffffffffffff610a3160406102778c8b611722565b165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f20548015610f8b57808203610f5d575050610a768887611722565b9160a4820135906102621901811215610ea35701600401610a9a61088e8380611774565b9167ffffffffffffffff80610ac06060610277610207610ab98961175f565b9680611774565b16911614610acd8361175f565b9015610f405750610add8261175f565b91602081013563ffffffff8116809103610ea35767ffffffffffffffff60405194610b0786611633565b16845260208401908152604080850192013582526101c08301359361ffff8516809503610ea357610180610b3f6102208601866118b5565b906040610b506102408901896118b5565b94909882519a7fb710e364000000000000000000000000000000000000000000000000000000008c5260048c01526101008160248d0137608061010082016101248d0137016101a48a01376102806101e4890152816102848901527f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8211610ea3576102c4908896979594939260051b80916102a48901378601826102a482016003196102a48a850301016102048a0152520193905f5b818110610f15575050509367ffffffffffffffff84939281602097511661022486015251166102448401525161026483015203815f73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af1908115610f0a575f91610ecf575b5015610ea7575167ffffffffffffffff1691610c9d6116c7565b50610ca88685611722565b604001610cb49061175f565b9060405191610cc28361164f565b84835267ffffffffffffffff166020830152610cde8786611722565b80610ce891611774565b80610cf291611789565b606001610cfe9061175f565b60405194610d0b8661164f565b855267ffffffffffffffff166020850152610d268786611722565b80610d3091611774565b80610d3a91611789565b608001610d4690611737565b96610d518187611722565b80610d5b91611774565b80610d6591611789565b610200013590610d759087611722565b80610d7f91611774565b80610d8991611789565b6101c001359060405198610d9c8a611633565b6001600160801b031689526020890152604088015260405195610dbe8761166b565b606013610ea35761014096610ea195610e7f94610e5e93610dea60405193610de585611633565b61174b565b8352602083015260408201528752602087019081526040870192835260608701948552610e38604051809851604080916001600160801b038151168452602081015160208501520151910152565b5180516001600160801b03166060880152602081015160808801526040015160a0870152565b5160c085019067ffffffffffffffff60208092828151168552015116910152565b5161010083019067ffffffffffffffff60208092828151168552015116910152565bf35b5f80fd5b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90506020813d602011610f02575b81610eea60209383611687565b81010312610ea357518015158103610ea3575f610c83565b3d9150610edd565b6040513d5f823e3d90fd5b91969550919293602080600192610f2b8a611893565b15158152019701910191879596949392610c07565b67ffffffffffffffff906390f4dbed60e01b5f521660045260245ffd5b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b610fce60606102776102078467ffffffffffffffff95611774565b1603610fdb575f806109d0565b608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b634e487b7160e01b5f52601160045260245ffd5b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b608460405163f492ef2b60e01b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152fd5b90506001600160801b036111ee8a611737565b168403906001600160801b038211611044576001600160801b03809116911610155f6108d1565b608460405163f492ef2b60e01b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f74636800000000000000000000000000000000000000000000000000000000006064820152fd5b634e487b7160e01b5f52603260045260245ffd5b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d9061083f565b507f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d610811565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d610764565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d610739565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d610710565b6113616113679284611e66565b83611eaa565b505f80610563565b61138491935060086020860153600185611eaa565b915f610603565b61139482611e47565b60010180600111611044576113af6105fb9161060f97611909565b9550506105ef565b92506113c282611e47565b6001018060011161104457926105d6565b67ffffffffffffffff6113e96113679385611e66565b9151169083611eaa565b61141491925067ffffffffffffffff90600860208601535116600184611eaa565b905f610551565b61142490611e47565b600101806001116110445761143f61053d9161056a95611909565b9350610530565b611451919250611e47565b6001018060011161104457905f610518565b634e487b7160e01b5f52604160045260245ffd5b6114cd6114bb925191516040519384937ff6b6676b000000000000000000000000000000000000000000000000000000008552604060048601526044850190611824565b83810360031901602485015290611824565b0390fd5b67ffffffffffffffff9051166390f4dbed60e01b5f5260045260245ffd5b60019150145f6102bf565b634e487b7160e01b5f52602160045260245ffd5b61152067ffffffffffffffff9261175f565b907f55bace6f000000000000000000000000000000000000000000000000000000005f526004521660245260445ffd5b507f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005f90815290917fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be5b8183106115b257505090602061014292820101610136565b602091935080600191548385880101520191019091839261159a565b6020925061014294915060ff191682840152151560051b820101610136565b634e487b7160e01b5f52602260045260245ffd5b92607f1692610117565b7f3921c703000000000000000000000000000000000000000000000000000000005f5260045ffd5b6060810190811067ffffffffffffffff82111761146357604052565b6040810190811067ffffffffffffffff82111761146357604052565b6080810190811067ffffffffffffffff82111761146357604052565b90601f8019910116810190811067ffffffffffffffff82111761146357604052565b604051906116b682611633565b5f6040838281528260208201520152565b604051906116d48261166b565b816116dd6116a9565b81526116e76116a9565b60208201526040516116f88161164f565b5f81525f602082015260408201526060604051916117158361164f565b5f83525f60208401520152565b903590605e1981360301821215610ea3570190565b356001600160801b0381168103610ea35790565b35906001600160801b0382168203610ea357565b3567ffffffffffffffff81168103610ea35790565b903590603e1981360301821215610ea3570190565b9035906102be1981360301821215610ea3570190565b903590601e1981360301821215610ea3570180359067ffffffffffffffff8211610ea357602001918136038313610ea357565b67ffffffffffffffff811161146357601f01601f191660200190565b9291926117fa826117d2565b916118086040519384611687565b829481845281830111610ea3578281602093845f960137010152565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b359067ffffffffffffffff82168203610ea357565b9190826040910312610ea3576040516118758161164f565b602061188e81839561188681611848565b855201611848565b910152565b35908115158203610ea357565b90359060be1981360301821215610ea3570190565b903590601e1981360301821215610ea3570180359067ffffffffffffffff8211610ea357602001918160051b36038313610ea357565b9190820391821161104457565b908151811015611280570160200190565b9190820180921161104457565b6040516119228161164f565b606081525f60208201525080518015908115611b82575b50611b5a575f19908051805b611ae2575b505f198214611ad4576001820190818311611044577f30000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000006119ab84846118f8565b51161480611ac0575b611ab0575f5b8151831015611a45576119cd83836118f8565b5160f81c603081108015611a3b575b611a295767ffffffffffffffff60ff8192602f19011681600a850216011691168110611a0d576001909201916119ba565b5091505060405190611a1e8261164f565b81525f602082015290565b505091505060405190611a1e8261164f565b50603981116119dc565b9150916001811190811591611aa4575b50611a7c5767ffffffffffffffff9060405192611a718461164f565b835216602082015290565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f611a55565b91505060405190611a1e8261164f565b506002611ace8483516118eb565b116119b4565b6040519150611a1e8261164f565b5f198101818111611044577f2d000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000611b3983866118f8565b511614611b5057508015611044575f190180611945565b92505f905061194a565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f611939565b6001600160801b03633b9aca00911602906001600160801b03821691820361104457565b80518210156112805760209160051b010190565b67ffffffffffffffff81511667ffffffffffffffff835116908181105f14611bf057505050505f90565b1115611bfd575050600290565b602067ffffffffffffffff81819301511692015116908181105f14611c225750505f90565b1115611c2d57600290565b600190565b90611c3c826117d2565b611c496040519182611687565b8281528092611c5a601f19916117d2565b0190602036910137565b805115611cab57611c758151611e47565b6001018060011161104457611c926106d7611ca792845190611909565b91600a60208401536107048151600185611eaa565b5090565b50604051611cba602082611687565b5f80825236602083013790565b80516001810180911161104457611cdd90611c32565b8051156112805760209181611cf9845f94019284845382611ee3565b50604051918291518091835e8101838152039060025afa15610f0a575f5190565b8015611cab57611d2981611e47565b6001018060011161104457611d40611ca791611c32565b9160086020840153600183611eaa565b611ca760405191611d62606084611687565b602283526040366020850137600a6020840153611d80600184611e94565b83611f35565b929192611d9382856118eb565b9360018514611e375760015b8060011b9086821015611db25750611d9f565b939495505090611dcc611dc58486611909565b8583611d86565b611de093611dda9195611909565b90611d86565b90604051611def608082611687565b6041815260208101906060368337805115611280576020935f936001845360218301526041820152604051918291518091835e8101838152039060025afa15610f0a575f5190565b50611e43929350611bb2565b5190565b906001915b6080811015611e585750565b60019060071c920191611e4c565b602082601092010153600181018091116110445790565b602082601292010153600181018091116110445790565b6020828192010153600181018091116110445790565b9092919083016020015b6080821015611ec857906001929391530190565b600180916080607f85161781530193019060071c9092611eb4565b908051918215611f05576021602084930191015e600101806001116110445790565b505050600190565b825191928215611f2f5783916020611f2c95818694019201015e611909565b90565b50505090565b8160209193929301015260208101809111611044579056fea26c6578706572696d656e74616cf564736f6c63430008250018",
}

// ContractUpdateClientABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractUpdateClientMetaData.ABI instead.
var ContractUpdateClientABI = ContractUpdateClientMetaData.ABI

// ContractUpdateClientBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractUpdateClientMetaData.Bin instead.
var ContractUpdateClientBin = ContractUpdateClientMetaData.Bin

// DeployContractUpdateClient deploys a new Ethereum contract, binding an instance of ContractUpdateClient to it.
func DeployContractUpdateClient(auth *bind.TransactOpts, backend bind.ContractBackend, signatureVerifier common.Address) (common.Address, *types.Transaction, *ContractUpdateClient, error) {
	parsed, err := ContractUpdateClientMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractUpdateClientBin), backend, signatureVerifier)
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

// VerifyHeader is a paid mutator transaction binding the contract method 0xc67440cc.
//
// Solidity: function verifyHeader(((uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,bytes20)[])),(uint64,uint64)),uint128,(uint256[8],uint256[4],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns(((uint128,bytes32,bytes32),(uint128,bytes32,bytes32),(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientTransactor) VerifyHeader(opts *bind.TransactOpts, msg_ ISpectreClientMsgsMsgUpdateApplicationState) (*types.Transaction, error) {
	return _ContractUpdateClient.contract.Transact(opts, "verifyHeader", msg_)
}

// VerifyHeader is a paid mutator transaction binding the contract method 0xc67440cc.
//
// Solidity: function verifyHeader(((uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,bytes20)[])),(uint64,uint64)),uint128,(uint256[8],uint256[4],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns(((uint128,bytes32,bytes32),(uint128,bytes32,bytes32),(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientSession) VerifyHeader(msg_ ISpectreClientMsgsMsgUpdateApplicationState) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.VerifyHeader(&_ContractUpdateClient.TransactOpts, msg_)
}

// VerifyHeader is a paid mutator transaction binding the contract method 0xc67440cc.
//
// Solidity: function verifyHeader(((uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,bytes20)[])),(uint64,uint64)),uint128,(uint256[8],uint256[4],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns(((uint128,bytes32,bytes32),(uint128,bytes32,bytes32),(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientTransactorSession) VerifyHeader(msg_ ISpectreClientMsgsMsgUpdateApplicationState) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.VerifyHeader(&_ContractUpdateClient.TransactOpts, msg_)
}
