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

// ISpectreClientMsgsMsgUpdateApplicationState is an auto generated low-level Go binding around an user-defined struct.
type ISpectreClientMsgsMsgUpdateApplicationState struct {
	TrustedConsensusState IICS07TendermintMsgsConsensusState
	ProposedHeader        IICS07TendermintMsgsHeader
	Time                  *big.Int
	Proof                 ISpectreClientMsgsBatchProof
}

// ISpectreClientMsgsVerifyHeaderOutput is an auto generated low-level Go binding around an user-defined struct.
type ISpectreClientMsgsVerifyHeaderOutput struct {
	TrustedConsensusState IICS07TendermintMsgsConsensusState
	NewConsensusState     IICS07TendermintMsgsConsensusState
	TrustedHeight         IICS02ClientMsgsHeight
	NewHeight             IICS02ClientMsgsHeight
}

// ContractUpdateClientMetaData contains all meta data concerning the ContractUpdateClient contract.
var ContractUpdateClientMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"signatureVerifier\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyHeader\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.MsgUpdateApplicationState\",\"components\":[{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"proposedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Header\",\"components\":[{\"name\":\"signedHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.SignedHeader\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockHeader\",\"components\":[{\"name\":\"version\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.Version\",\"components\":[{\"name\":\"blockVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"chainId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasLastBlockId\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastBlockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"hasLastCommitHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastCommitHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasDataHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"consensusHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"appHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasLastResultsHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"lastResultsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hasEvidenceHash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"evidenceHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"commit\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockCommit\",\"components\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockId\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.BlockId\",\"components\":[{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"partSetHeader\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.PartSetHeader\",\"components\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashData\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]},{\"name\":\"commitSigs\",\"type\":\"tuple[]\",\"internalType\":\"structIICS07TendermintMsgs.CommitSig[]\",\"components\":[{\"name\":\"flag\",\"type\":\"uint8\",\"internalType\":\"enumIICS07TendermintMsgs.CommitSigFlag\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.CommitSigData\",\"components\":[{\"name\":\"validatorAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hasSignature\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}]}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]},{\"name\":\"time\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.BatchProof\",\"components\":[{\"name\":\"proof\",\"type\":\"uint256[8]\",\"internalType\":\"uint256[8]\"},{\"name\":\"commitments\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitmentPok\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"bucket\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"signerIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"pinnedValidatorIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"signerPubkeys\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"active\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structISpectreClientMsgs.VerifyHeaderOutput\",\"components\":[{\"name\":\"trustedConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"newConsensusState\",\"type\":\"tuple\",\"internalType\":\"structIICS07TendermintMsgs.ConsensusState\",\"components\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nextValidatorsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"trustedHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"newHeight\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.Height\",\"components\":[{\"name\":\"revisionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revisionHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"ConsensusStateHashMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ConsensusStateNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DirectCallNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToVerifyHeader\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"HeaderChainIdMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidChainIdLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChainPrefixLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidHeaderHeight\",\"inputs\":[{\"name\":\"height\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MismatchedRevisionHeight\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actual\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ProofVerificationFailed\",\"inputs\":[]}]",
	Bin: "0x60c034607557601f61242138819003918201601f19168301916001600160401b03831184841017607957808492602094604052833981010312607557516001600160a01b038116908190036075576080523060a052604051612393908161008e823960805181610e73015260a0518161026e0152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c63b615f64414610024575f80fd5b346101185760203660031901126101185760043567ffffffffffffffff81116101185760c060031982360301126101185761006c61014091600401610067610204565b610256565b61011660606040519261009d848251604080916001600160801b038151168452602081015160208501520151910152565b6100cb602082015183860190604080916001600160801b038151168452602081015160208501520151910152565b6100f3604082015160c086019067ffffffffffffffff60208092828151168552015116910152565b015161010083019067ffffffffffffffff60208092828151168552015116910152565bf35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b6060810190811067ffffffffffffffff82111761014c57604052565b61011c565b6040810190811067ffffffffffffffff82111761014c57604052565b90601f8019910116810190811067ffffffffffffffff82111761014c57604052565b6040519061019e60408361016d565b565b6040519061019e60608361016d565b6040519061019e6102608361016d565b6040519061019e60808361016d565b604051906101db82610130565b5f6040838281528260208201520152565b604051906101f982610151565b5f6020838281520152565b604051906080820182811067ffffffffffffffff82111761014c576040528161022b6101ce565b81526102356101ce565b60208201526102426101ec565b604082015260606102516101ec565b910152565b5073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163014610415576104129061040c6102cd7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17025467ffffffffffffffff1690565b6103fe6103eb6102db61018f565b6102e361043d565b81526102fd6020820194859067ffffffffffffffff169052565b610396867f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170354926103746103446103378663ffffffff1690565b9560501c63ffffffff1690565b61036761034f6101a0565b96610358610550565b885263ffffffff166020880152565b63ffffffff166040860152565b60608201936103838584610592565b91610390608085016105b8565b926109ad565b6103e560405160208101906103bd816103af8b856105cd565b03601f19810183528261016d565b5190206103dd6103d860406103d2868c610592565b01610611565b610b25565b80821461061b565b85610592565b6103f860a0860186610652565b90610d38565b5167ffffffffffffffff1690565b90610f1b565b90565b7f3921c703000000000000000000000000000000000000000000000000000000005f5260045ffd5b604051905f7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df1700548060011c91600182168015610546575b602084108114610532578386528592602084019190811561051957506001146104a5575b5061019e9250038361016d565b7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df17005f90815291507fa4cc147281017aae47c0db3f306061a3149e6999ca67777149d3e595a1e6b4be5b848310610502575061019e9350015f610498565b8054828401528693506020909201916001016104ee565b60ff191682525061019e93151560051b0190505f610498565b634e487b7160e01b5f52602260045260245ffd5b92607f1692610474565b6040519061055d82610151565b81602060ff7f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170154818116845260081c16910152565b903590605e1981360301821215610118570190565b6001600160801b0381160361011857565b35610412816105a7565b359061019e826105a7565b91909160408060608301946001600160801b0381356105eb816105a7565b168452602081013560208501520135910152565b67ffffffffffffffff81160361011857565b35610412816105ff565b15610624575050565b7f6d23e8a3000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b90359061021e1981360301821215610118570190565b903590603e1981360301821215610118570190565b9035906102be1981360301821215610118570190565b903590601e1981360301821215610118570180359067ffffffffffffffff82116101185760200191813603831361011857565b67ffffffffffffffff811161014c57601f01601f191660200190565b9291926106ee826106c6565b916106fc604051938461016d565b829481845281830111610118578281602093845f960137010152565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b909161075361041293604084526040840190610718565b916020818403910152610718565b359061019e826105ff565b91908260409103126101185760405161078481610151565b60208082948035610794816105ff565b84520135916107a2836105ff565b0152565b9080601f8301121561011857816020610412933591016106e2565b8015150361011857565b359061019e826107c1565b63ffffffff81160361011857565b809291039160608312610118576040516107fd81610151565b6040819483358352601f19011261011857602090604080519361081f85610151565b8381013561082c816107d6565b85520135828401520152565b6102c0813603126101185761084b6101af565b90610856368261076c565b8252604081013567ffffffffffffffff81116101185761087990369083016107a6565b602083015261088a60608201610761565b604083015261089b608082016105c2565b60608301526108ac60a082016107cb565b60808301526108be3660c083016107e4565b60a08301526108d061012082016107cb565b60c083015261014081013560e08301526108ed61016082016107cb565b6101008301526101808101356101208301526101a08101356101408301526101c08101356101608301526101e08101356101808301526102008101356101a083015261093c61022082016107cb565b6101c08301526102408101356101e083015261095b61026082016107cb565b6102008301526102808101356102208301526102a08101359067ffffffffffffffff82116101185761098f913691016107a6565b61024082015290565b90359060be1981360301821215610118570190565b93929190936109e46109df6109d86109ce6109c88580610668565b8061067d565b6040810190610693565b36916106e2565b6110f9565b6109ee818361138e565b602086015167ffffffffffffffff1667ffffffffffffffff610a2b610a1e602085015167ffffffffffffffff1690565b67ffffffffffffffff1690565b911603610aec5750610a50610a4b610a466109c88480610668565b610838565b6114a8565b6040610a69610a5f8480610668565b6020810190610998565b013503610a7b5761019e94519061196a565b60405163f492ef2b60e01b815260206004820152602360248201527f696e76616c696420626c6f636b3a206865616465722068617368206d69736d6160448201527f7463680000000000000000000000000000000000000000000000000000000000606482015280608481015b0390fd5b5185516040517fc6913e9f000000000000000000000000000000000000000000000000000000008152918291610ae8916004840161073c565b67ffffffffffffffff165f527f5e9f695da54d5993169597b18a19a3159f4a62cbfd7f6669a996ca09c6df170460205260405f20548015610b635790565b7f5b48b457000000000000000000000000000000000000000000000000000000005f5260045ffd5b15610b935750565b67ffffffffffffffff906390f4dbed60e01b5f521660045260245ffd5b35610412816107d6565b3561ffff811681036101185790565b903590601e1981360301821215610118570180359067ffffffffffffffff821161011857602001918160051b3603831361011857565b908160209103126101185751610412816107c1565b9998979593919261010060409461ffff8694168d5260208d01376101208b01376101608901376102406101a0880152806102408801527f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8111610118576020916102609160051b8091838a0137870181888203016101c0890152018281520191905f5b818110610cdb57505050906101e08261019e93509401906040809167ffffffffffffffff815116845267ffffffffffffffff60208201511660208501520151910152565b9091926020806001928635610cef816107c1565b15158152019401929101610c97565b6040513d5f823e3d90fd5b15610d1057565b7fd611c318000000000000000000000000000000000000000000000000000000005f5260045ffd5b90602090610d8f610d55610d4c8580610668565b84810190610998565b9367ffffffffffffffff610d7d610a1e60606103d26109c8610d768b610611565b9680610668565b911614610d8985610611565b90610b8b565b610d9883610611565b92610de46040610db5610dac868501610bb0565b63ffffffff1690565b92013591610dd4610dc46101a0565b67ffffffffffffffff9097168752565b67ffffffffffffffff1685850152565b6040840152610df66101808201610bba565b610e59610e076101e0840184610bc9565b610e18610200869893980186610bc9565b9160405198899788977f7d1a869500000000000000000000000000000000000000000000000000000000895261014082019161010081019160048b01610c14565b03815f73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af18015610edb5761019e915f91610eac575b50610d09565b610ece915060203d602011610ed4575b610ec6818361016d565b810190610bff565b5f610ea6565b503d610ebc565b610cfe565b919082606091031261011857604051610ef881610130565b60408082948035610f08816105a7565b8452602081013560208501520135910152565b90610f24610204565b5060608201610f3860406103d28386610592565b90610f63610f4461018f565b67ffffffffffffffff851681529267ffffffffffffffff166020840152565b610fad610f8160606103d26109c8610f7b868a610592565b80610668565b610f9c610f8c61018f565b67ffffffffffffffff9096168652565b67ffffffffffffffff166020850152565b610fc86080610fc26109c8610f7b8589610592565b016105b8565b906101c0610ff16109c8610f7b610200610fe86109c8610f7b888d610592565b01359489610592565b01359061100e610fff6101a0565b6001600160801b039094168452565b6020830152604082015261102b6110236101bf565b943690610ee0565b845260208401526040830152606082015290565b6040519061104c82610151565b5f602083606081520152565b634e487b7160e01b5f52601160045260245ffd5b8015611078575f190190565b611058565b5f1981019190821161107857565b9190820391821161107857565b634e487b7160e01b5f52603260045260245ffd5b9081518110156110bd570160200190565b611098565b906001820180921161107857565b600101908160011161107857565b602301908160231161107857565b9190820180921161107857565b61110161103f565b5080518015908115611382575b5061135a575f19908051805b6112d6575b505f1982146112b8577f30000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000006111a561117f611179866110c2565b856110ac565b517fff000000000000000000000000000000000000000000000000000000000000001690565b1614806112c2575b6112b8575f906111bc836110c2565b915b815183101561124f576111dd6111d761117f85856110ac565b60f81c90565b60ff811660308110908115611244575b506112375767ffffffffffffffff60ff8192602f19011681600a85021601169116811061121f576001909201916111be565b5091505061122b61018f565b9081525f602082015290565b505091505061122b61018f565b60399150115f6111ed565b91509160018111908115916112ac575b50611284576104129061127061018f565b92835267ffffffffffffffff166020830152565b7fe131a1dc000000000000000000000000000000000000000000000000000000005f5260045ffd5b602b915010155f61125f565b905061122b61018f565b5060026112d083835161108b565b116111ad565b7f2d0000000000000000000000000000000000000000000000000000000000000061133461130f61117f6113098561107d565b866110ac565b7fff000000000000000000000000000000000000000000000000000000000000001690565b14611348576113429061106c565b8061111a565b61135391925061107d565b905f61111f565b7fa4482ffc000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040915010155f61110e565b602091820180519282019267ffffffffffffffff84356113ad816105ff565b1667ffffffffffffffff82160361144f5750906114256113ea60606103d26109c86113e361142a975167ffffffffffffffff1690565b9580610668565b936114066113f661018f565b67ffffffffffffffff9094168452565b61141e6020840195869067ffffffffffffffff169052565b369061076c565b611b55565b6114315750565b516390f4dbed60e01b5f5267ffffffffffffffff1660045260245b5ffd5b8361146d61146761144c9367ffffffffffffffff1690565b91610611565b7f05b4c7a3000000000000000000000000000000000000000000000000000000005f5267ffffffffffffffff91821660045216602452604490565b610412906116786114c56102406114ca6114c56020860151611dbb565b611e2b565b936114d3611b93565b946114e16114c58351611e7b565b6114ea87611bb6565b526114f486611bc3565b526115186114c5611513610a1e604085015167ffffffffffffffff1690565b611f98565b61152186611bd3565b526115416114c561153c60608401516001600160801b031690565b611fcc565b61154a86611be3565b52608081015115611723576115656114c560a083015161208f565b61156e86611bf3565b5260c0810151156116fd576115896114c560e083015161215b565b61159286611c03565b52610100810151156116d7576115af6114c561012083015161215b565b6115b886611c13565b526115ca6114c561014083015161215b565b6115d386611c23565b526115e56114c561016083015161215b565b6115ee86611c34565b526116006114c561018083015161215b565b61160986611c45565b5261161b6114c56101a083015161215b565b61162486611c56565b526101c0810151156116b1576116416114c56101e083015161215b565b61164a86611c67565b526102008101511561168b576116676114c561022083015161215b565b61167086611c78565b520151611dbb565b61168182611c89565b525f815191612193565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611667565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611641565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d6115af565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611589565b7f6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d611565565b6001600160801b03633b9aca00911602906001600160801b03821691820361107857565b906001600160801b03809116911603906001600160801b03821161107857565b1561179457565b608460405162461bcd60e51b815260206004820152602560248201527f696e76616c696420626c6f636b3a206e6f6e206d6f6e6f746f6e69632062667460448201527f2074696d650000000000000000000000000000000000000000000000000000006064820152fd5b1561180557565b606460405162461bcd60e51b815260206004820152602060248201527f696e76616c696420626c6f636b3a20636861696e2d6964206d69736d617463686044820152fd5b906001600160801b03809116911601906001600160801b03821161107857565b1561187057565b608460405162461bcd60e51b815260206004820152602860248201527f696e76616c696420626c6f636b3a206865616465722069732066726f6d20746860448201527f65206675747572650000000000000000000000000000000000000000000000006064820152fd5b67ffffffffffffffff60019116019067ffffffffffffffff821161107857565b1561190157565b608460405162461bcd60e51b8152602060048201526024808201527f696e76616c696420626c6f636b3a206e6f6e20696e6372656173696e6720686560448201527f69676874000000000000000000000000000000000000000000000000000000006064820152fd5b9391929061198a611985610dac602087015163ffffffff1690565b611749565b6119a2611996856105b8565b6001600160801b031690565b6001600160801b03841610908115611b2d575b50611ac157611985610dac6040611a6396611a2d611a3895611a186109d8611a0f6109c88e610f7b611a3e9e6001600160801b03611a07611996611a016080610fc26109c88980610668565b936105b8565b91161161178d565b86810190610693565b602081519101209060208151910120146117fe565b015163ffffffff1690565b90611849565b6001600160801b0380611a596080610fc26109c88880610668565b9216911610611869565b67ffffffffffffffff611a80611a7b60408401610611565b6118da565b81611a9360606103d26109c88780610668565b9116918291161115611aa3575050565b611abb610a1e60606103d26109c88661019e97610668565b146118fa565b60405163f492ef2b60e01b815260206004820152603c60248201527f696e76616c696420626c6f636b3a20756e74727573746564207374617465206960448201527f73206f757473696465206f66207472757374696e6720706572696f64000000006064820152608490fd5b90506001600160801b0380611b4a611b44876105b8565b8661176d565b92169116115f6119b5565b90611b5f91611cae565b6003811015611b7f5760028114908115611b77575090565b600191501490565b634e487b7160e01b5f52602160045260245ffd5b6040516101e09190611ba5838261016d565b600e815291601f1901366020840137565b8051156110bd5760200190565b8051600110156110bd5760400190565b8051600210156110bd5760600190565b8051600310156110bd5760800190565b8051600410156110bd5760a00190565b8051600510156110bd5760c00190565b8051600610156110bd5760e00190565b8051600710156110bd576101000190565b8051600810156110bd576101200190565b8051600910156110bd576101400190565b8051600a10156110bd576101600190565b8051600b10156110bd576101800190565b8051600c10156110bd576101a00190565b8051600d10156110bd576101c00190565b80518210156110bd5760209160051b010190565b805167ffffffffffffffff1667ffffffffffffffff611cd8610a1e855167ffffffffffffffff1690565b911681811015611cea57505050505f90565b1115611cf7575050600290565b611d2c610a1e6020611d1c8167ffffffffffffffff95015167ffffffffffffffff1690565b94015167ffffffffffffffff1690565b911681811015611d3c5750505f90565b1115611d4757600290565b600190565b60405190611d5b60208361016d565b5f808352366020840137565b60405160809190611d78838261016d565b6041815291601f1901366020840137565b90611d93826106c6565b611da0604051918261016d565b8281528092611db1601f19916106c6565b0190602036910137565b805115611e1057611dcc815161222d565b806001019081600111611078576001908351010180911161107857611df3611e0c91611d89565b91600a6020840153611e068151846122a8565b83612343565b5090565b50610412611d4c565b805191908290602001825e015f815290565b80516001810180911161107857611e4190611d89565b8051156110bd57611e6b81611e5e6020945f868196015382612319565b5060405191828092611e19565b039060025afa15610edb575f5190565b5f90611e8f815167ffffffffffffffff1690565b67ffffffffffffffff8116611f72575b50611ec96020820192611ebd610a1e855167ffffffffffffffff1690565b80611f51575b50611d89565b915f91611ee1610a1e825167ffffffffffffffff1690565b611f2e575b50611efc610a1e825167ffffffffffffffff1690565b611f0557505090565b611f27610a1e611f18611e0c9486612258565b925167ffffffffffffffff1690565b90836122e0565b611f4a919250611f43610a1e611f188661224c565b90846122e0565b905f611ee6565b90611f66611f61611f6c9361222d565b6110d0565b906110ec565b5f611ec3565b611f91919250611f8c611f619167ffffffffffffffff1690565b61222d565b905f611e9f565b8015611e1057611fa78161222d565b6001018060011161107857611fbe611e0c91611d89565b9160086020840153826122a8565b6001600160801b0316633b9aca0081066001600160801b03633b9aca005f930416918215159182612071575b6001600160801b0316908115159081612056575b61201590611d89565b935f9361203a575b5061202757505090565b612034611e0c9284612258565b836122e0565b61204f919350600860208601536001856122e0565b915f61201d565b612062611f618461222d565b81018091111561200c57611058565b90506001600160801b03612087611f618561222d565b919050611ff8565b60208101516120a463ffffffff82511661222d565b90816001019182600111611078576023018092116110785761210b6120cb611e0c93611d89565b916008602084015360206121016120fb6120f56120ef610dac865163ffffffff1690565b876122a8565b8661226f565b85612286565b910151908361236e565b50611e0681516121556120f56121396121348461212f61212a8261222d565b6110de565b6110ec565b611d89565b9661214c6121468961229c565b89612286565b9051908861236e565b856122e0565b6040519060609061216c828461016d565b60228352611e0c91600a906020850190601f1901368237536020602184015360028361236e565b9092919280840393808511611078576001851461221c5760015b8060011b90868210156121c057506121ad565b9192939495505082019182811161107857826121dc9185612193565b916121e79293612193565b6121ef611d67565b918251156110bd57825f92611e6b9260016020809701536021830152604182015260405191828092611e19565b5090612229929350611c9a565b5190565b906001915b608081101561223e5750565b60019060071c920191612232565b60206008910153600190565b602082601092010153600181018091116110785790565b602082601292010153600181018091116110785790565b6020828192010153600181018091116110785790565b6020600a910153600190565b91906021600193015b60808210156122c557906001929391530190565b600180916080607f85161781530193019060071c90926122b1565b9092919083016020015b60808210156122fe57906001929391530190565b600180916080607f85161781530193019060071c90926122ea565b90805191821561233b576021602084930191015e600101806001116110785790565b505050600190565b90809291825192831561236757839260208092019201015e81018091116110785790565b5050505090565b8160209193929301015260208101809111611078579056fea164736f6c634300081c000a",
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

// VerifyHeader is a paid mutator transaction binding the contract method 0xb615f644.
//
// Solidity: function verifyHeader(((uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),uint128,(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns(((uint128,bytes32,bytes32),(uint128,bytes32,bytes32),(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientTransactor) VerifyHeader(opts *bind.TransactOpts, msg_ ISpectreClientMsgsMsgUpdateApplicationState) (*types.Transaction, error) {
	return _ContractUpdateClient.contract.Transact(opts, "verifyHeader", msg_)
}

// VerifyHeader is a paid mutator transaction binding the contract method 0xb615f644.
//
// Solidity: function verifyHeader(((uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),uint128,(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns(((uint128,bytes32,bytes32),(uint128,bytes32,bytes32),(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientSession) VerifyHeader(msg_ ISpectreClientMsgsMsgUpdateApplicationState) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.VerifyHeader(&_ContractUpdateClient.TransactOpts, msg_)
}

// VerifyHeader is a paid mutator transaction binding the contract method 0xb615f644.
//
// Solidity: function verifyHeader(((uint128,bytes32,bytes32),((((uint64,uint64),string,uint64,uint128,bool,(bytes32,(uint32,bytes32)),bool,bytes32,bool,bytes32,bytes32,bytes32,bytes32,bytes32,bool,bytes32,bool,bytes32,bytes),(uint64,uint32,(bytes32,(uint32,bytes32)),(uint8,(bytes,uint128,bool,bytes))[])),(uint64,uint64)),uint128,(uint256[8],uint256[2],uint256[2],uint16,uint32[],uint32[],bytes32[],bool[])) msg_) returns(((uint128,bytes32,bytes32),(uint128,bytes32,bytes32),(uint64,uint64),(uint64,uint64)))
func (_ContractUpdateClient *ContractUpdateClientTransactorSession) VerifyHeader(msg_ ISpectreClientMsgsMsgUpdateApplicationState) (*types.Transaction, error) {
	return _ContractUpdateClient.Contract.VerifyHeader(&_ContractUpdateClient.TransactOpts, msg_)
}
