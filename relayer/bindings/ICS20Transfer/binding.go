// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractICS20Transfer

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

// IIBCAppCallbacksOnAcknowledgementPacketCallback is an auto generated low-level Go binding around an user-defined struct.
type IIBCAppCallbacksOnAcknowledgementPacketCallback struct {
	SourceClient      string
	DestinationClient string
	Sequence          uint64
	Payload           IICS26RouterMsgsPayload
	Acknowledgement   []byte
	Relayer           common.Address
}

// IIBCAppCallbacksOnRecvPacketCallback is an auto generated low-level Go binding around an user-defined struct.
type IIBCAppCallbacksOnRecvPacketCallback struct {
	SourceClient      string
	DestinationClient string
	Sequence          uint64
	Payload           IICS26RouterMsgsPayload
	Relayer           common.Address
}

// IIBCAppCallbacksOnTimeoutPacketCallback is an auto generated low-level Go binding around an user-defined struct.
type IIBCAppCallbacksOnTimeoutPacketCallback struct {
	SourceClient      string
	DestinationClient string
	Sequence          uint64
	Payload           IICS26RouterMsgsPayload
	Relayer           common.Address
}

// IICS20TransferMsgsSendTransferMsg is an auto generated low-level Go binding around an user-defined struct.
type IICS20TransferMsgsSendTransferMsg struct {
	Denom            common.Address
	Amount           *big.Int
	Receiver         string
	SourceClient     string
	DestPort         string
	TimeoutTimestamp uint64
	Memo             string
}

// IICS26RouterMsgsPayload is an auto generated low-level Go binding around an user-defined struct.
type IICS26RouterMsgsPayload struct {
	SourcePort string
	DestPort   string
	Version    string
	Encoding   string
	Value      []byte
}

// ISignatureTransferPermitTransferFrom is an auto generated low-level Go binding around an user-defined struct.
type ISignatureTransferPermitTransferFrom struct {
	Permitted ISignatureTransferTokenPermissions
	Nonce     *big.Int
	Deadline  *big.Int
}

// ISignatureTransferTokenPermissions is an auto generated low-level Go binding around an user-defined struct.
type ISignatureTransferTokenPermissions struct {
	Token  common.Address
	Amount *big.Int
}

// ContractICS20TransferMetaData contains all meta data concerning the ContractICS20Transfer contract.
var ContractICS20TransferMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEscrow\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEscrowBeacon\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIBCERC20Beacon\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPermit2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ibcERC20Contract\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ibcERC20Denom\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ics26\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"ics26Router\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"escrowLogic\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"ibcERC20Logic\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permit2\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializeV2\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isConsumingScheduledOp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onAcknowledgementPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIIBCAppCallbacks.OnAcknowledgementPacketCallback\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destinationClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"acknowledgement\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayer\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onRecvPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIIBCAppCallbacks.OnRecvPacketCallback\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destinationClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"relayer\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onTimeoutPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIIBCAppCallbacks.OnTimeoutPacketCallback\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destinationClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"relayer\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sendTransfer\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS20TransferMsgs.SendTransferMsg\",\"components\":[{\"name\":\"denom\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"memo\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendTransferWithPermit2\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS20TransferMsgs.SendTransferMsg\",\"components\":[{\"name\":\"denom\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"memo\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"permit\",\"type\":\"tuple\",\"internalType\":\"structISignatureTransfer.PermitTransferFrom\",\"components\":[{\"name\":\"permitted\",\"type\":\"tuple\",\"internalType\":\"structISignatureTransfer.TokenPermissions\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendTransferWithSender\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS20TransferMsgs.SendTransferMsg\",\"components\":[{\"name\":\"denom\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"memo\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAuthority\",\"inputs\":[{\"name\":\"newAuthority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setCustomERC20\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeEscrowTo\",\"inputs\":[{\"name\":\"newEscrowLogic\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeIBCERC20To\",\"inputs\":[{\"name\":\"newIBCERC20Logic\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCERC20ContractCreated\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fullDenomPath\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessManagedInvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessManagedRequiredDelay\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"AccessManagedUnauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ICS20DenomAlreadyExists\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20DenomNotFound\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20EscrowNotFound\",\"inputs\":[{\"name\":\"clientID\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20InvalidAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20InvalidAmount\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ICS20InvalidPort\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20Permit2TokenMismatch\",\"inputs\":[{\"name\":\"permitToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sentToken\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ICS20TokenAlreadyExists\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ICS20UnauthorizedPacketSender\",\"inputs\":[{\"name\":\"packetSender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ICS20UnexpectedERC20Balance\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ICS20UnexpectedEncoding\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20UnexpectedVersion\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60a080604052346100c257306080525f5160206156395f395f51905f525460ff8160401c166100b3576002600160401b03196001600160401b03821601610060575b60405161557290816100c782396080518181816112b5015261137f0152f35b6001600160401b0319166001600160401b039081175f5160206156395f395f51905f525581527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f80610041565b63f92ee8a960e01b5f5260045ffd5b5f80fdfe60806040526004361015610011575f80fd5b5f5f3560e01c806306ab20bc14612c11578063078c4a79146124455780631459457a14611f405780631bbf2e2314611efa5780631e5150e414611eb457806329b6eca914611ba45780632ac3dc3814611aea5780633f4ba83a14611a0b578063428e4e171461161d5780634f1ef2861461132d57806352d1902d1461129a57806353816a7c14610f715780635c975abb14610f2f5780635e32b6b614610d095780637a9e5e4b14610c4b578063826cae7a14610c055780638456cb5914610b4c5780638fb3603714610ab7578063969631d514610a3b578063a1d28f5714610777578063a50ee2b4146106dc578063aaa2c34314610617578063ac9650d814610469578063ad3cb1cc14610408578063b29c715d14610263578063bf7e214f1461021d578063d413227d146101d75763e163b1af1461014e575f80fd5b346101d45760206003193601126101d4576101d06101b56101bc6101a9610173612cd3565b6001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b60405192838092613274565b0382612db4565b604051918291602083526020830190612d46565b0390f35b80fd5b50346101d457806003193601126101d45760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416604051908152f35b50346101d457806003193601126101d45760206001600160a01b037ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a005416604051908152f35b50346101d45760406003193601126101d45760043567ffffffffffffffff81116103d457806004019060e06003198236030112610404576102a2612ce9565b916102ab613706565b6102b3613692565b6102bd3633613427565b60248201359182156103d8576102f06102eb6102e460646001600160a01b03940185612ec7565b3691612df3565b6137cd565b1691610306816102ff846131d7565b8533613a13565b84610310836131d7565b91843b156103d4576040517fb4f22eb70000000000000000000000000000000000000000000000000000000081526001600160a01b0393909316600484015233602484015260448301528160648183875af180156103c9576103b0575b60208561037b868686613c21565b907f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d67ffffffffffffffff60405191168152f35b6103bb858092612db4565b6103c5575f61036d565b8380fd5b6040513d87823e3d90fd5b5080fd5b6024857f4f6df8d000000000000000000000000000000000000000000000000000000000815280600452fd5b8280fd5b50346101d457806003193601126101d457506101d060405161042b604082612db4565b600581527f352e302e300000000000000000000000000000000000000000000000000000006020820152604051918291602083526020830190612d46565b50346101d45760206003193601126101d45760043567ffffffffffffffff81116103d457366023820112156103d45780600401359067ffffffffffffffff821161040457602481013660248460051b840101116103c55760405160206104cf8183612db4565b85825280820192601f1982013685376104e7866133ce565b946104f56040519687612db4565b868652601f19610504886133ce565b0183895b82811061060757505050875b8781101561058a5760019061056e61056861053760248460051b87010187612ec7565b91908d898c856040519687958487013784018281018481528e519283915e010190815203601f198101835282612db4565b30614490565b610578828a6133e6565b5261058381896133e6565b5001610514565b83898860405191838301848452825180915260408401948060408360051b870101940192955b8287106105bd5785850386f35b9091929382806105f7837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08a600196030186528851612d46565b96019201960195929190926105b0565b606082828b010152018490610508565b50346101d45760206003193601126101d45780610632612cd3565b61063c3633613427565b6001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805541690813b156106d8576001600160a01b03602484928360405195869485937f3659cfe60000000000000000000000000000000000000000000000000000000085521660048401525af180156106cd576106bc5750f35b816106c691612db4565b6101d45780f35b6040513d84823e3d90fd5b5050fd5b50346101d45760206003193601126101d45760043567ffffffffffffffff81116103d45761070e903690600401612e29565b6001600160a01b0361072082846131eb565b541690811561073457602082604051908152f35b90506107736040519283927fe1275e2f000000000000000000000000000000000000000000000000000000008452602060048501526024840191612f53565b0390fd5b50346101d45760406003193601126101d45760043567ffffffffffffffff81116103d4576107a9903690600401612e29565b6107b4929192612ce9565b6107be3633613427565b6001600160a01b036107d083866131eb565b54166109ff57610819610813826001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b54613223565b15610854826001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b90156109c157506108d29061086983866131eb565b6001600160a01b03808316167fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790556001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9067ffffffffffffffff8111610994576108f6816108f08454613223565b84613314565b82601f82116001146109345781908495610924949592610929575b50505f198260011b9260031b1c19161790565b905580f35b013590505f80610911565b601f198216948385526020852091855b87811061097c575083600195969710610963575b505050811b01905580f35b5f1960f88560031b161c199101351690555f8080610958565b90926020600181928686013581550194019101610944565b6024837f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b610773906040519182917f778769c4000000000000000000000000000000000000000000000000000000008352602060048401526024830190613274565b6040517f0c0ef5340000000000000000000000000000000000000000000000000000000081526020600482015280610773602482018588612f53565b50346101d45760206003193601126101d45760043567ffffffffffffffff81116103d45760209182610a796001600160a01b03933690600401612e29565b925082604051938492833781017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8008152030190205416604051908152f35b50346101d457806003193601126101d4577ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a005460a01c60ff1615610b44575060207f8fb36037000000000000000000000000000000000000000000000000000000005b7fffffffff0000000000000000000000000000000000000000000000000000000060405191168152f35b602090610b1a565b50346101d457806003193601126101d457610b673633613427565b610b6f613706565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005416177fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a180f35b50346101d457806003193601126101d45760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8045416604051908152f35b50346101d45760206003193601126101d457610c65612cd3565b6001600160a01b037ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a0054163303610cdd57803b15610ca957610ca6906143f2565b80f35b7fc2f31e5e0000000000000000000000000000000000000000000000000000000082526001600160a01b0316600452602490fd5b6024827f068ca9d800000000000000000000000000000000000000000000000000000000815233600452fd5b50346101d45780610d1936612d13565b610d4f336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314612e57565b610d57613692565b610d5f613706565b60608101610dae610d88610d80610d768486612e94565b6080810190612ec7565b810190613070565b610d9b610d958486612e94565b80612ec7565b90610da68680612ec7565b9290916141cf565b9050610db9816148b6565b610de6575b5050507f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d80f35b6001600160a01b031691823b15610f2a57610ef992610ef28580946080946001600160a01b03604051988997889687957f5e32b6b600000000000000000000000000000000000000000000000000000000875260206004880152610ec3610ea1610e64610e5389806140a3565b60a060248d015260c48c0191612f53565b610e7160208a018a6140a3565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc8c84030160448d0152612f53565b9167ffffffffffffffff610eb760408a016140f3565b1660648a015287614108565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc88830301608489015261413a565b9301612cff565b1660a483015203925af180156106cd57610f15575b8080610dbe565b81610f1f91612db4565b6101d457805f610f0e565b505050fd5b50346101d457806003193601126101d457602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b50346101d45760c06003193601126101d45760043567ffffffffffffffff81116103d457806004019060e060031982360301126104045760807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc3601126104045760a43567ffffffffffffffff81116103c557610ff2903690600401612e29565b610ffd929192613706565b611005613692565b602482013591821561126e576110196131c1565b6001600160a01b038061102b886131d7565b169116146110376131c1565b611040876131d7565b91156112345750506102eb6102e4606461105b930187612ec7565b926001600160a01b03807f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065416941693604051906040820182811067ffffffffffffffff821117611207579088949392916040528682526020820191868352813b1561120357856001600160a01b039161114a8296604051988997889687957f30f28b7a00000000000000000000000000000000000000000000000000000000875281611106612ce9565b166004880152604435602488015260643560448801526084356064880152511660848601525160a48501523360c485015261010060e4850152610104840191612f53565b03925af180156106cd576111ee575b50611163846131d7565b91833b156103d4576040517fb4f22eb70000000000000000000000000000000000000000000000000000000081526001600160a01b0393909316600484015233602484015260448301528160648183865af180156111e3576111ce575b60208461037b338587613c21565b6111d9848092612db4565b610404575f6111c0565b6040513d86823e3d90fd5b816111f891612db4565b6103c557835f611159565b8580fd5b6024897f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b7fe36164780000000000000000000000000000000000000000000000000000000088526001600160a01b0390811660045216602452604486fd5b6024867f4f6df8d000000000000000000000000000000000000000000000000000000000815280600452fd5b50346101d457806003193601126101d4576001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036113055760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b807fe07c8dba0000000000000000000000000000000000000000000000000000000060049252fd5b5060406003193601126101d457611342612cd3565b9060243567ffffffffffffffff81116103d457366023820112156103d457611374903690602481600401359101612df3565b916001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168030149081156115e8575b506115c0576113b93633613427565b6001600160a01b03811690604051937f52d1902d000000000000000000000000000000000000000000000000000000008552602085600481865afa80958596611588575b5061142e57602484847f4c9c8ce3000000000000000000000000000000000000000000000000000000008252600452fd5b9091847f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc810361155d5750823b1561153257807fffffffffffffffffffffffff00000000000000000000000000000000000000007f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5416177f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b8480a2805115611500576114fc91614490565b5080f35b50503461150a5780f35b807fb398979f0000000000000000000000000000000000000000000000000000000060049252fd5b7f4c9c8ce3000000000000000000000000000000000000000000000000000000008452600452602483fd5b7faa1d49a4000000000000000000000000000000000000000000000000000000008552600452602484fd5b9095506020813d6020116115b8575b816115a460209383612db4565b810103126115b45751945f6113fd565b8480fd5b3d9150611597565b6004827fe07c8dba000000000000000000000000000000000000000000000000000000008152fd5b90506001600160a01b037f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc541614155f6113aa565b50346101d45760206003193601126101d4578060043567ffffffffffffffff8111611a08578060040160c060031983360301126106d85761168a336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314612e57565b611692613692565b61169a613706565b606482016116ae610d80610d768385612e94565b9160848401906116c16102e48383612ec7565b602081519101209360409460208087516116db8982612db4565b818152017f4774d4a575993f963b1c06573736617a457abef8589178db8d10c94b4ab511ab815220145f146118a8576117269061171b610d958685612e94565b90610da68580612ec7565b9050611731816148b6565b611764575b505050505050505b807f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d80f35b6001600160a01b031690813b156118a4578680948651978895869485937f8dfcd9ad0000000000000000000000000000000000000000000000000000000085528560048601528a60248601526117ba82806140a3565b6044870160c090526101048701906117d192612f53565b6117de60248601846140a3565b6043198884030160648901526117f49291612f53565b90611801604486016140f3565b67ffffffffffffffff16608487015261181a9083614108565b908581036043190160a48701526118309161413a565b9161183a916140a3565b6043198584030160c48601526118509291612f53565b9060a40161185d90612cff565b6001600160a01b031660e483015203925af190811561189b5750611886575b8080808080611736565b8161189091612db4565b6101d457805f61187c565b513d84823e3d90fd5b8680fd5b60206118b5910151613759565b6118be816148b6565b6118cf575b5050505050505061173e565b6001600160a01b031690813b156118a4578680948651978895869485937f8dfcd9ad00000000000000000000000000000000000000000000000000000000855260048501600190528a602486015261192782806140a3565b6044870160c0905261010487019061193e92612f53565b61194b60248601846140a3565b6043198884030160648901526119619291612f53565b9061196e604486016140f3565b67ffffffffffffffff1660848701526119879083614108565b908581036043190160a487015261199d9161413a565b916119a7916140a3565b6043198584030160c48601526119bd9291612f53565b9060a4016119ca90612cff565b6001600160a01b031660e483015203925af190811561189b57506119f3575b80808080806118c3565b816119fd91612db4565b6101d457805f6119e9565b50fd5b50346101d457806003193601126101d457611a263633613427565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff811615611ac2577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00167fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a180f35b6004827f8dfc202b000000000000000000000000000000000000000000000000000000008152fd5b50346101d45760206003193601126101d45760043567ffffffffffffffff81116103d457806004019060e0600319823603011261040457611b29613706565b611b31613692565b6024810135908115611b7857611b586102eb6102e460646001600160a01b03940186612ec7565b1690611b6e81611b67856131d7565b8433613a13565b83611163846131d7565b6024847f4f6df8d000000000000000000000000000000000000000000000000000000000815280600452fd5b50346101d45760206003193601126101d457611bbe612cd3565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005467ffffffffffffffff81169060018203611e8c5760401c60ff16908115611e80575b50611e5857611c7560027fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000007ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005416177ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0055565b680100000000000000007fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005416177ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0055602460206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416604051928380927f24d7806c0000000000000000000000000000000000000000000000000000000082523360048301525afa908115611e4d578391611e0e575b5090611d65611d7a923390612e57565b611d6d6145dc565b611d756145dc565b6143f2565b7fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054167ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a180f35b90506020813d602011611e45575b81611e2960209383612db4565b8101031261040457519081151582036104045790611d65611d55565b3d9150611e1c565b6040513d85823e3d90fd5b6004827ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b6002915010155f611c02565b6004847ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b50346101d457806003193601126101d45760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8055416604051908152f35b50346101d457806003193601126101d45760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065416604051908152f35b50346101d45760a06003193601126101d457611f5a612cd3565b611f62612ce9565b6044356001600160a01b03811681036103c557606435926001600160a01b0384168094036115b4576084356001600160a01b0381168103611203577ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005467ffffffffffffffff8116908161241d5760401c60ff16908115612411575b506123e957906120e76001600160a01b039261205e60027fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000007ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005416177ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0055565b680100000000000000007fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005416177ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00556120d76145dc565b6120df6145dc565b611d656145dc565b167fffffffffffffffffffffffff00000000000000000000000000000000000000007f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80355604051916105129283810181811067ffffffffffffffff8211176123bc576121918291614b7095878785396001600160a01b0316815230602082015260400190565b039086f080156103c9576001600160a01b03167fffffffffffffffffffffffff00000000000000000000000000000000000000007f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8045416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80455604051928084019284841067ffffffffffffffff8511176123bc57918493916122479385396001600160a01b0316815230602082015260400190565b039083f080156106cd576001600160a01b03167fffffffffffffffffffffffff00000000000000000000000000000000000000007f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8055416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805557fffffffffffffffffffffffff00000000000000000000000000000000000000007f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f806557fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054167ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a180f35b6024877f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b6004867ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b6002915010155f611fde565b6004887ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b50346101d45761245436612d13565b9061248b336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314612e57565b612493613692565b61249b613706565b60608201916124ba6102e46124b08584612e94565b6040810190612ec7565b602081519101206124c9612f18565b60208151910120146124d9612f18565b6124e66124b08685612e94565b909215612bdb575050506125386125036102e4610d958685612e94565b60208151910120612512612f9b565b6020815191012014612522612f9b565b90612530610d958786612e94565b929091612fd6565b6125526102e46125488584612e94565b6060810190612ec7565b6020815191012061256161301a565b602081519101201461257161301a565b61257e6125488685612e94565b909215612ba5575050506125d26125a56102e461259b8685612e94565b6020810190612ec7565b602081519101206125b4612f9b565b60208151910120146125c4612f9b565b9061253061259b8786612e94565b6125e2610d80610d768584612e94565b9260608401805115611b78576125fb6040860151613759565b9160208401936126116102eb6102e48784612ec7565b965194612640612624610d958585612e94565b61263b6126318680612ec7565b9390923691612df3565b613967565b9261264b84886144d4565b156127d85750509051845195968795909250906020908781189088110287188061267d612678828b613363565b61339d565b9803920101602087015e612692855186614526565b9590158681156127c6575b5061279d575b506001600160a01b03905b16905193813b156103c5576040517f0779afe60000000000000000000000000000000000000000000000000000000081526001600160a01b039182166004820152921660248301526044820193909352918290606490829084905af180156106cd57612788575b6101d08260405190612728604083612db4565b601182527f7b22726573756c74223a2241513d3d227d00000000000000000000000000000060208301527f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d604051918291602083526020830190612d46565b612793828092612db4565b6101d4575f612715565b94506001600160a01b03906127c0826127b58861313d565b54169687151561317b565b906126a3565b6001600160a01b03915016155f61269d565b612842935061280d83602098949361263b6126316128056127fc8d98978998612e94565b87810190612ec7565b939094612ec7565b926040519784899551918291018487015e8401908282018a8152815193849201905e010186815203601f198101855284612db4565b6001600160a01b038516946001600160a01b0361285e8561313d565b54169384156128ff575b5084956001600160a01b038596959616908351823b156118a4576040517f40c10f190000000000000000000000000000000000000000000000000000000081526001600160a01b0392909216600483015260248201529085908290604490829084905af19081156103c95785916128ea575b50506001600160a01b03906126ae565b816128f491612db4565b6103c557835f6128da565b93506001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80454166040517f4571e3a600000000000000000000000000000000000000000000000000000000602082015230602482015287604482015260606064820152612987816129796084820189612d46565b03601f198101835282612db4565b604051916104e48084019084821067ffffffffffffffff831117612b7857918493916129b79361508286396137ad565b039086f080156103c9576001600160a01b0316956129d48561313d565b6001600160a01b0388167fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055612a3f876001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9685519767ffffffffffffffff8911612b4b57612a6689612a608354613223565b83613314565b602098601f8111600114612ae95780612a96918a9b8b9a9b91612ade575b505f198260011b9260031b1c19161790565b90555b7f6031fab685dd6d86e4dbac9a69eae347145f332c95b3a0d728d3730fc5233d62612ad38298604051918291602083526020830190612d46565b0390a2959493612868565b90508a01515f612a84565b818952898920601f1982168a5b818110612b335750908a9b83600194939c9b9c10612b1b575b5050811b019055612a99565b8b01515f1960f88460031b161c191690555f80612b0f565b8a8d0151835560209c8d019c60019093019201612af6565b6024887f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b60248a7f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b610773906040519384937fd1ca953a00000000000000000000000000000000000000000000000000000000855260048501612f73565b610773906040519384937f094af3b800000000000000000000000000000000000000000000000000000000855260048501612f73565b5034612ccf576020600319360112612ccf57612c2b612cd3565b612c353633613427565b6001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f804541690813b15612ccf576001600160a01b0360245f928360405195869485937f3659cfe60000000000000000000000000000000000000000000000000000000085521660048401525af18015612cc457612cb6575080f35b612cc291505f90612db4565b005b6040513d5f823e3d90fd5b5f80fd5b600435906001600160a01b0382168203612ccf57565b602435906001600160a01b0382168203612ccf57565b35906001600160a01b0382168203612ccf57565b6020600319820112612ccf576004359067ffffffffffffffff8211612ccf576003198260a092030112612ccf5760040190565b90601f19601f602080948051918291828752018686015e5f8582860101520116010190565b60a0810190811067ffffffffffffffff821117612d8757604052565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b90601f601f19910116810190811067ffffffffffffffff821117612d8757604052565b67ffffffffffffffff8111612d8757601f01601f191660200190565b929192612dff82612dd7565b91612e0d6040519384612db4565b829481845281830111612ccf578281602093845f960137010152565b9181601f84011215612ccf5782359167ffffffffffffffff8311612ccf5760208381860195010111612ccf57565b15612e5f5750565b6001600160a01b03907f2ecb3242000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6181360301821215612ccf570190565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215612ccf570180359067ffffffffffffffff8211612ccf57602001918136038313612ccf57565b60405190612f27604083612db4565b600782527f69637332302d31000000000000000000000000000000000000000000000000006020830152565b601f8260209493601f1993818652868601375f8582860101520116010190565b91612f8a612f989492604085526040850190612d46565b926020818503910152612f53565b90565b60405190612faa604083612db4565b600882527f7472616e736665720000000000000000000000000000000000000000000000006020830152565b9290919215612fe457505050565b610773906040519384937f5d3a3cdd00000000000000000000000000000000000000000000000000000000855260048501612f73565b60405190613029604083612db4565b601a82527f6170706c69636174696f6e2f782d736f6c69646974792d6162690000000000006020830152565b9080601f83011215612ccf57816020612f9893359101612df3565b602081830312612ccf5780359067ffffffffffffffff8211612ccf570160a081830312612ccf57604051916130a483612d6b565b813567ffffffffffffffff8111612ccf57816130c1918401613055565b8352602082013567ffffffffffffffff8111612ccf57816130e3918401613055565b6020840152604082013567ffffffffffffffff8111612ccf5781613108918401613055565b604084015260608201356060840152608082013567ffffffffffffffff8111612ccf576131359201613055565b608082015290565b60208091604051928184925191829101835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80181520301902090565b156131835750565b610773906040519182917fe1275e2f000000000000000000000000000000000000000000000000000000008352602060048401526024830190612d46565b6024356001600160a01b0381168103612ccf5790565b356001600160a01b0381168103612ccf5790565b60209082604051938492833781017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80181520301902090565b90600182811c9216801561326a575b602083101461323d57565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f1691613232565b5f929181549161328383613223565b80835292600181169081156132d8575060011461329f57505050565b5f9081526020812093945091925b8383106132be575060209250010190565b6001816020929493945483858701015201910191906132ad565b905060209495507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091509291921683830152151560051b010190565b601f821161332157505050565b5f5260205f20906020601f840160051c83019310613359575b601f0160051c01905b81811061334e575050565b5f8155600101613343565b909150819061333a565b9190820391821161337057565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b906133a782612dd7565b6133b46040519182612db4565b828152601f196133c48294612dd7565b0190602036910137565b67ffffffffffffffff8111612d875760051b60200190565b80518210156133fa5760209160051b010190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a0054916001600160a01b0383169281600411612ccf575f5f9060405f8151966001600160a01b0360208901917fb700961300000000000000000000000000000000000000000000000000000000835216978860248201523060448201527fffffffff000000000000000000000000000000000000000000000000000000008335166064820152606481526134dc608482612db4565b828052826020525190895afa61367f575b156134fa575b5050505050565b63ffffffff1615613653577fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1674010000000000000000000000000000000000000000177ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a0055823b15612ccf576020925f92836040518096819582947f94c7d7ee000000000000000000000000000000000000000000000000000000008452600484015260406024840152601f19601f6044850192808452808786860137868582860101520116010103925af18015612cc457613643575b507fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff7ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a0054167ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a00555f808080806134f3565b5f61364d91612db4565b5f6135d2565b827f068ca9d8000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b50505f516020518060201c1502906134ed565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005c6136de5760017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d565b7f3ee5aeb5000000000000000000000000000000000000000000000000000000005f5260045ffd5b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300541661373157565b7fd93c0665000000000000000000000000000000000000000000000000000000005f5260045ffd5b613764815182614526565b91901561376f575090565b610773906040519182917f3fed5d87000000000000000000000000000000000000000000000000000000008352602060048401526024830190612d46565b6040906001600160a01b03612f9894931681528160208201520190612d46565b604051906001600160a01b03815192602081818501958087835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f800815203019020541691821561381e57505090565b7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805547ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a00546040517f485cc9550000000000000000000000000000000000000000000000000000000060208201523060248201526001600160a01b039182166044808301919091528152939450919291166138b9606483612db4565b604051916104e4908184019284841067ffffffffffffffff851117612d875784936138e89361508286396137ad565b03905ff08015612cc4576001600160a01b036020911692604051928391518091835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8008152030190206001600160a01b0382167fffffffffffffffffffffffff000000000000000000000000000000000000000082541617905590565b6001809160208095612f98958160405198858a9651918291018688015e8501917f2f0000000000000000000000000000000000000000000000000000000000000085840152602183013701017f2f000000000000000000000000000000000000000000000000000000000000008382015203017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1810184520182612db4565b9190820180921161337057565b916001600160a01b03909391931692604051927f70a082310000000000000000000000000000000000000000000000000000000084526001600160a01b03821691826004860152602085602481895afa948515612cc4575f95613bec575b506040517f23b872dd0000000000000000000000000000000000000000000000000000000060208281019182526001600160a01b039485166024840152939092166044820152606481018590525f9190613ace8160848101612979565b519082885af115612cc4575f513d613be35750833b155b613bb7576020906024604051809681937f70a0823100000000000000000000000000000000000000000000000000000000835260048301525afa928315612cc4575f93613b81575b50613b389082613a06565b90821180613b78575b15613b4a575050565b7f2fb30cfc000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b50808214613b41565b9092506020813d602011613baf575b81613b9d60209383612db4565b81010312612ccf575191613b38613b2d565b3d9150613b90565b837f5274afe7000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b60011415613ae5565b9094506020813d602011613c19575b81613c0860209383612db4565b81010312612ccf5751936020613a71565b3d9150613bfb565b613c42613c30610173836131d7565b92613c495f9460405193848092613274565b0383612db4565b818051155f14613fde5750505060a0613c72613c6c613c67846131d7565b614633565b94614633565b93613c806040840184612ec7565b939095613cd1613cb6613c9660c0850185612ec7565b99909760405196613ca688612d6b565b8752602087019485523691612df3565b9560408501968752606085019860208501358a523691612df3565b94608084019586526001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f803541695613d126060850185612ec7565b96909401359467ffffffffffffffff8616809603613fda5790613e26613d7f949392613e18613d3f612f9b565b9c613d48612f9b565b94613de1613d54612f18565b97613db0613d6061301a565b9a6040519c8d986020808b01525160a060408b015260e08a0190612d46565b90517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08983030160608a0152612d46565b90517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0878303016080880152612d46565b915160a0850152517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08483030160c0850152612d46565b03601f198101865285612db4565b60405199613e338b612d6b565b8a5260208a0152604089015260608801526080870152604051926060840184811067ffffffffffffffff8211176123bc5793613f7460209694613e89613edd958a9567ffffffffffffffff996040523691612df3565b835287830190815260408301998a52604051998a97889687957f4d6e7ce30000000000000000000000000000000000000000000000000000000087528b600488015251606060248801526084870190612d46565b925116604485015251907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc8482030160648501526080613f63613f51613f3f613f2f865160a0875260a0870190612d46565b8d8701518682038f880152612d46565b60408601518582036040870152612d46565b60608501518482036060860152612d46565b920151906080818403910152612d46565b03925af1918215613fcd578192613f8a57505090565b9091506020813d602011613fc5575b81613fa660209383612db4565b810103126103d457519067ffffffffffffffff821682036101d4575090565b3d9150613f99565b50604051903d90823e3d90fd5b8880fd5b61400990614003613ff0979497612f9b565b613ffd6060880188612ec7565b91613967565b906144d4565b61401a575b50613c7260a091614633565b6001600160a01b0361402b846131d7565b16803b15612ccf576040517f9dc29fac0000000000000000000000000000000000000000000000000000000081526001600160a01b03929092166004830152602084013560248301525f908290604490829084905af18015612cc4571561400e576140999193505f90612db4565b5f91613c7261400e565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811215612ccf57016020813591019167ffffffffffffffff8211612ccf578136038313612ccf57565b359067ffffffffffffffff82168203612ccf57565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6182360301811215612ccf570190565b612f98916141c16141b661419b61418061416561415787806140a3565b60a0885260a0880191612f53565b61417260208801886140a3565b908783036020890152612f53565b61418d60408701876140a3565b908683036040880152612f53565b6141a860608601866140a3565b908583036060870152612f53565b9260808101906140a3565b916080818503910152612f53565b909493925f926001600160a01b03604051838382376020818581017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80081520301902054169283156143b3579161423b9161263b6142429461423360208a0151613759565b9a3691612df3565b84516144d4565b15614363576001600160a01b03614259845161313d565b541692614269815185151561317b565b6060810151843b15612ccf576040517f40c10f190000000000000000000000000000000000000000000000000000000081526001600160a01b038416600482015260248101919091525f8160448183895af18015612cc45761434d575b506060905b0151813b15610404576040517f0779afe60000000000000000000000000000000000000000000000000000000081526001600160a01b0385811660048301528716602482015260448101919091529082908290606490829084905af180156106cd57614338575b50509190565b614343828092612db4565b6101d45780614332565b61435a9193505f90612db4565b5f9160606142c6565b6001600160a01b03614375845161313d565b5416928315614387575b6060906142cb565b925060606143958451613759565b936143ac81516001600160a01b038716151561317b565b905061437f565b50906107736040519283927f5778f378000000000000000000000000000000000000000000000000000000008452602060048501526024840191612f53565b60206001600160a01b037f2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad9216807fffffffffffffffffffffffff00000000000000000000000000000000000000007ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a005416177ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a0055604051908152a1565b5f80612f9893602081519101845af43d156144cc573d916144b083612dd7565b926144be6040519485612db4565b83523d5f602085013e614722565b606091614722565b805190825180921061451f5780519182808210911802808318921415820282189060206145046126788486613363565b92808285019503920101835e51902090602081519101201490565b5050505f90565b8051821180156145d5575b61457e576001821180614586575b158015908160011b9182046002141715613370576028018060281161337057820361457e576001600160a01b0392915f614578926147ae565b90921690565b50505f905f90565b507f30780000000000000000000000000000000000000000000000000000000000007fffff0000000000000000000000000000000000000000000000000000000000006020830151161461453f565b505f614531565b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c161561460b57565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f5260045ffd5b6001600160a01b031680614647602a612dd7565b916146556040519384612db4565b602a8352614663602a612dd7565b601f196020850191013682378351156133fa57603090538251600110156133fa576078602184015360295b600181116146cf575061469f575090565b7fe22e27eb000000000000000000000000000000000000000000000000000000005f52600452601460245260445ffd5b90600f811660108110156133fa5784518310156133fa577f3031323334353637383961626364656600000000000000000000000000000000901a8483016020015360041c908015613370575f190161468e565b9061475f575080511561473757602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b815115806147a5575b614770575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15614768565b9290926001840180851161337057831180614865575b158015908160011b9182046002141715613370576147e7905f9492939495613a06565b915b8183106147f95750505060019190565b9092919360ff6148307fff000000000000000000000000000000000000000000000000000000000000006020888601015116614aa3565b16600f811161485a578160041b9180830460101490151715613370576001910194019192906147e9565b505f94508493505050565b507f30780000000000000000000000000000000000000000000000000000000000007fffff0000000000000000000000000000000000000000000000000000000000006020868401015116146147c4565b60205f604051828101907f01ffc9a70000000000000000000000000000000000000000000000000000000082527f01ffc9a700000000000000000000000000000000000000000000000000000000602482015260248152614918604482612db4565b519084617530fa903d5f519083614a97575b5082614a8d575b50816149f4575b81614941575090565b602091505f90604051838101907f01ffc9a70000000000000000000000000000000000000000000000000000000082527fffffffff000000000000000000000000000000000000000000000000000000007fd3ce6f1b00000000000000000000000000000000000000000000000000000000166024820152602481526149c8604482612db4565b5191617530fa5f513d826149e8575b50816149e1575090565b9050151590565b6020111591505f6149d7565b905060205f604051828101907f01ffc9a70000000000000000000000000000000000000000000000000000000082527fffffffff00000000000000000000000000000000000000000000000000000000602482015260248152614a58604482612db4565b519084617530fa5f513d82614a81575b5081614a77575b501590614938565b905015155f614a6f565b6020111591505f614a68565b151591505f614931565b6020111592505f61492a565b60f81c602f811180614b65575b15614add577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd00160ff1690565b6060811180614b5b575b15614b14577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa90160ff1690565b6040811180614b51575b15614b4b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc90160ff1690565b5060ff90565b5060478110614b1e565b5060678110614ae7565b50603a8110614ab056fe60803461013457601f61051238819003918201601f19168301916001600160401b03831184841017610138578084926040948552833981010312610134576100468161014c565b906001600160a01b039061005c9060200161014c565b16908115610121575f80546001600160a01b031981168417825560405193916001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a3803b1561010157600180546001600160a01b0319166001600160a01b039290921691821790557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a26103b190816101618239f35b63211eb15960e21b5f9081526001600160a01b0391909116600452602490fd5b631e4fbdf760e01b5f525f60045260245ffd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b03821682036101345756fe60806040526004361015610011575f80fd5b5f3560e01c80633659cfe61461027e5780635c60da1b1461022d578063715018a6146101935780638da5cb5b146101435763f2fde38b14610050575f80fd5b3461013f5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013f5760043573ffffffffffffffffffffffffffffffffffffffff811680910361013f576100a8610358565b80156101135773ffffffffffffffffffffffffffffffffffffffff5f54827fffffffffffffffffffffffff00000000000000000000000000000000000000008216175f55167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3005b7f1e4fbdf7000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b5f80fd5b3461013f575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013f57602073ffffffffffffffffffffffffffffffffffffffff5f5416604051908152f35b3461013f575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013f576101c9610358565b5f73ffffffffffffffffffffffffffffffffffffffff81547fffffffffffffffffffffffff000000000000000000000000000000000000000081168355167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b3461013f575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013f57602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b3461013f5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013f5760043573ffffffffffffffffffffffffffffffffffffffff81169081810361013f576102d7610358565b3b1561032d57807fffffffffffffffffffffffff000000000000000000000000000000000000000060015416176001557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a2005b7f847ac564000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b73ffffffffffffffffffffffffffffffffffffffff5f5416330361037857565b7f118cdaa7000000000000000000000000000000000000000000000000000000005f523360045260245ffdfea164736f6c634300081c000a60a0806040526104e480380380916100178285610292565b833981016040828203126101eb5761002e826102c9565b602083015190926001600160401b0382116101eb57019080601f830112156101eb57815161005b816102dd565b926100696040519485610292565b8184526020840192602083830101116101eb57815f926020809301855e84010152823b15610274577fa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d5080546001600160a01b0319166001600160a01b038516908117909155604051635c60da1b60e01b8152909190602081600481865afa9081156101f7575f9161023a575b50803b1561021a5750817f1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e5f80a282511561020257602060049260405193848092635c60da1b60e01b82525afa9182156101f7575f926101ae575b505f809161018a945190845af43d156101a6573d9161016e836102dd565b9261017c6040519485610292565b83523d5f602085013e6102f8565b505b60805260405161018d908161035782396080518160460152f35b6060916102f8565b9291506020833d6020116101ef575b816101ca60209383610292565b810103126101eb575f80916101e161018a956102c9565b9394509150610150565b5f80fd5b3d91506101bd565b6040513d5f823e3d90fd5b505050341561018c5763b398979f60e01b5f5260045ffd5b634c9c8ce360e01b5f9081526001600160a01b0391909116600452602490fd5b90506020813d60201161026c575b8161025560209383610292565b810103126101eb57610266906102c9565b5f6100f5565b3d9150610248565b631933b43b60e21b5f9081526001600160a01b038416600452602490fd5b601f909101601f19168101906001600160401b038211908210176102b557604052565b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b03821682036101eb57565b6001600160401b0381116102b557601f01601f191660200190565b9061031c575080511561030d57602081519101fd5b63d6bda27560e01b5f5260045ffd5b8151158061034d575b61032d575090565b639996b31560e01b5f9081526001600160a01b0391909116600452602490fd5b50803b1561032556fe60806040527f5c60da1b000000000000000000000000000000000000000000000000000000006080526020608060048173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa8015610107575f9015610163575060203d602011610100575b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f820116608001906080821067ffffffffffffffff8311176100d3576100ce91604052608001610112565b610163565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b503d610081565b6040513d5f823e3d90fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80602091011261015f5760805173ffffffffffffffffffffffffffffffffffffffff8116810361015f5790565b5f80fd5b5f8091368280378136915af43d5f803e1561017c573d5ff35b3d5ffdfea164736f6c634300081c000aa164736f6c634300081c000af0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
}

// ContractICS20TransferABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractICS20TransferMetaData.ABI instead.
var ContractICS20TransferABI = ContractICS20TransferMetaData.ABI

// ContractICS20TransferBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractICS20TransferMetaData.Bin instead.
var ContractICS20TransferBin = ContractICS20TransferMetaData.Bin

// DeployContractICS20Transfer deploys a new Ethereum contract, binding an instance of ContractICS20Transfer to it.
func DeployContractICS20Transfer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractICS20Transfer, error) {
	parsed, err := ContractICS20TransferMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractICS20TransferBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractICS20Transfer{ContractICS20TransferCaller: ContractICS20TransferCaller{contract: contract}, ContractICS20TransferTransactor: ContractICS20TransferTransactor{contract: contract}, ContractICS20TransferFilterer: ContractICS20TransferFilterer{contract: contract}}, nil
}

// ContractICS20Transfer is an auto generated Go binding around an Ethereum contract.
type ContractICS20Transfer struct {
	ContractICS20TransferCaller     // Read-only binding to the contract
	ContractICS20TransferTransactor // Write-only binding to the contract
	ContractICS20TransferFilterer   // Log filterer for contract events
}

// ContractICS20TransferCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractICS20TransferCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractICS20TransferTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractICS20TransferTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractICS20TransferFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractICS20TransferFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractICS20TransferSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractICS20TransferSession struct {
	Contract     *ContractICS20Transfer // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ContractICS20TransferCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractICS20TransferCallerSession struct {
	Contract *ContractICS20TransferCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ContractICS20TransferTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractICS20TransferTransactorSession struct {
	Contract     *ContractICS20TransferTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ContractICS20TransferRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractICS20TransferRaw struct {
	Contract *ContractICS20Transfer // Generic contract binding to access the raw methods on
}

// ContractICS20TransferCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractICS20TransferCallerRaw struct {
	Contract *ContractICS20TransferCaller // Generic read-only contract binding to access the raw methods on
}

// ContractICS20TransferTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractICS20TransferTransactorRaw struct {
	Contract *ContractICS20TransferTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractICS20Transfer creates a new instance of ContractICS20Transfer, bound to a specific deployed contract.
func NewContractICS20Transfer(address common.Address, backend bind.ContractBackend) (*ContractICS20Transfer, error) {
	contract, err := bindContractICS20Transfer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractICS20Transfer{ContractICS20TransferCaller: ContractICS20TransferCaller{contract: contract}, ContractICS20TransferTransactor: ContractICS20TransferTransactor{contract: contract}, ContractICS20TransferFilterer: ContractICS20TransferFilterer{contract: contract}}, nil
}

// NewContractICS20TransferCaller creates a new read-only instance of ContractICS20Transfer, bound to a specific deployed contract.
func NewContractICS20TransferCaller(address common.Address, caller bind.ContractCaller) (*ContractICS20TransferCaller, error) {
	contract, err := bindContractICS20Transfer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferCaller{contract: contract}, nil
}

// NewContractICS20TransferTransactor creates a new write-only instance of ContractICS20Transfer, bound to a specific deployed contract.
func NewContractICS20TransferTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractICS20TransferTransactor, error) {
	contract, err := bindContractICS20Transfer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferTransactor{contract: contract}, nil
}

// NewContractICS20TransferFilterer creates a new log filterer instance of ContractICS20Transfer, bound to a specific deployed contract.
func NewContractICS20TransferFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractICS20TransferFilterer, error) {
	contract, err := bindContractICS20Transfer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferFilterer{contract: contract}, nil
}

// bindContractICS20Transfer binds a generic wrapper to an already deployed contract.
func bindContractICS20Transfer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractICS20TransferMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractICS20Transfer *ContractICS20TransferRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractICS20Transfer.Contract.ContractICS20TransferCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractICS20Transfer *ContractICS20TransferRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.ContractICS20TransferTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractICS20Transfer *ContractICS20TransferRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.ContractICS20TransferTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractICS20Transfer *ContractICS20TransferCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractICS20Transfer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractICS20Transfer *ContractICS20TransferTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractICS20Transfer *ContractICS20TransferTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ContractICS20Transfer *ContractICS20TransferCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ContractICS20Transfer *ContractICS20TransferSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ContractICS20Transfer.Contract.UPGRADEINTERFACEVERSION(&_ContractICS20Transfer.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ContractICS20Transfer.Contract.UPGRADEINTERFACEVERSION(&_ContractICS20Transfer.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCaller) Authority(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "authority")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferSession) Authority() (common.Address, error) {
	return _ContractICS20Transfer.Contract.Authority(&_ContractICS20Transfer.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) Authority() (common.Address, error) {
	return _ContractICS20Transfer.Contract.Authority(&_ContractICS20Transfer.CallOpts)
}

// GetEscrow is a free data retrieval call binding the contract method 0x969631d5.
//
// Solidity: function getEscrow(string clientId) view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCaller) GetEscrow(opts *bind.CallOpts, clientId string) (common.Address, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "getEscrow", clientId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetEscrow is a free data retrieval call binding the contract method 0x969631d5.
//
// Solidity: function getEscrow(string clientId) view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferSession) GetEscrow(clientId string) (common.Address, error) {
	return _ContractICS20Transfer.Contract.GetEscrow(&_ContractICS20Transfer.CallOpts, clientId)
}

// GetEscrow is a free data retrieval call binding the contract method 0x969631d5.
//
// Solidity: function getEscrow(string clientId) view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) GetEscrow(clientId string) (common.Address, error) {
	return _ContractICS20Transfer.Contract.GetEscrow(&_ContractICS20Transfer.CallOpts, clientId)
}

// GetEscrowBeacon is a free data retrieval call binding the contract method 0x1e5150e4.
//
// Solidity: function getEscrowBeacon() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCaller) GetEscrowBeacon(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "getEscrowBeacon")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetEscrowBeacon is a free data retrieval call binding the contract method 0x1e5150e4.
//
// Solidity: function getEscrowBeacon() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferSession) GetEscrowBeacon() (common.Address, error) {
	return _ContractICS20Transfer.Contract.GetEscrowBeacon(&_ContractICS20Transfer.CallOpts)
}

// GetEscrowBeacon is a free data retrieval call binding the contract method 0x1e5150e4.
//
// Solidity: function getEscrowBeacon() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) GetEscrowBeacon() (common.Address, error) {
	return _ContractICS20Transfer.Contract.GetEscrowBeacon(&_ContractICS20Transfer.CallOpts)
}

// GetIBCERC20Beacon is a free data retrieval call binding the contract method 0x826cae7a.
//
// Solidity: function getIBCERC20Beacon() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCaller) GetIBCERC20Beacon(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "getIBCERC20Beacon")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetIBCERC20Beacon is a free data retrieval call binding the contract method 0x826cae7a.
//
// Solidity: function getIBCERC20Beacon() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferSession) GetIBCERC20Beacon() (common.Address, error) {
	return _ContractICS20Transfer.Contract.GetIBCERC20Beacon(&_ContractICS20Transfer.CallOpts)
}

// GetIBCERC20Beacon is a free data retrieval call binding the contract method 0x826cae7a.
//
// Solidity: function getIBCERC20Beacon() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) GetIBCERC20Beacon() (common.Address, error) {
	return _ContractICS20Transfer.Contract.GetIBCERC20Beacon(&_ContractICS20Transfer.CallOpts)
}

// GetPermit2 is a free data retrieval call binding the contract method 0x1bbf2e23.
//
// Solidity: function getPermit2() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCaller) GetPermit2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "getPermit2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPermit2 is a free data retrieval call binding the contract method 0x1bbf2e23.
//
// Solidity: function getPermit2() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferSession) GetPermit2() (common.Address, error) {
	return _ContractICS20Transfer.Contract.GetPermit2(&_ContractICS20Transfer.CallOpts)
}

// GetPermit2 is a free data retrieval call binding the contract method 0x1bbf2e23.
//
// Solidity: function getPermit2() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) GetPermit2() (common.Address, error) {
	return _ContractICS20Transfer.Contract.GetPermit2(&_ContractICS20Transfer.CallOpts)
}

// IbcERC20Contract is a free data retrieval call binding the contract method 0xa50ee2b4.
//
// Solidity: function ibcERC20Contract(string denom) view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCaller) IbcERC20Contract(opts *bind.CallOpts, denom string) (common.Address, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "ibcERC20Contract", denom)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IbcERC20Contract is a free data retrieval call binding the contract method 0xa50ee2b4.
//
// Solidity: function ibcERC20Contract(string denom) view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferSession) IbcERC20Contract(denom string) (common.Address, error) {
	return _ContractICS20Transfer.Contract.IbcERC20Contract(&_ContractICS20Transfer.CallOpts, denom)
}

// IbcERC20Contract is a free data retrieval call binding the contract method 0xa50ee2b4.
//
// Solidity: function ibcERC20Contract(string denom) view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) IbcERC20Contract(denom string) (common.Address, error) {
	return _ContractICS20Transfer.Contract.IbcERC20Contract(&_ContractICS20Transfer.CallOpts, denom)
}

// IbcERC20Denom is a free data retrieval call binding the contract method 0xe163b1af.
//
// Solidity: function ibcERC20Denom(address token) view returns(string)
func (_ContractICS20Transfer *ContractICS20TransferCaller) IbcERC20Denom(opts *bind.CallOpts, token common.Address) (string, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "ibcERC20Denom", token)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// IbcERC20Denom is a free data retrieval call binding the contract method 0xe163b1af.
//
// Solidity: function ibcERC20Denom(address token) view returns(string)
func (_ContractICS20Transfer *ContractICS20TransferSession) IbcERC20Denom(token common.Address) (string, error) {
	return _ContractICS20Transfer.Contract.IbcERC20Denom(&_ContractICS20Transfer.CallOpts, token)
}

// IbcERC20Denom is a free data retrieval call binding the contract method 0xe163b1af.
//
// Solidity: function ibcERC20Denom(address token) view returns(string)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) IbcERC20Denom(token common.Address) (string, error) {
	return _ContractICS20Transfer.Contract.IbcERC20Denom(&_ContractICS20Transfer.CallOpts, token)
}

// Ics26 is a free data retrieval call binding the contract method 0xd413227d.
//
// Solidity: function ics26() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCaller) Ics26(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "ics26")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Ics26 is a free data retrieval call binding the contract method 0xd413227d.
//
// Solidity: function ics26() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferSession) Ics26() (common.Address, error) {
	return _ContractICS20Transfer.Contract.Ics26(&_ContractICS20Transfer.CallOpts)
}

// Ics26 is a free data retrieval call binding the contract method 0xd413227d.
//
// Solidity: function ics26() view returns(address)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) Ics26() (common.Address, error) {
	return _ContractICS20Transfer.Contract.Ics26(&_ContractICS20Transfer.CallOpts)
}

// IsConsumingScheduledOp is a free data retrieval call binding the contract method 0x8fb36037.
//
// Solidity: function isConsumingScheduledOp() view returns(bytes4)
func (_ContractICS20Transfer *ContractICS20TransferCaller) IsConsumingScheduledOp(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "isConsumingScheduledOp")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// IsConsumingScheduledOp is a free data retrieval call binding the contract method 0x8fb36037.
//
// Solidity: function isConsumingScheduledOp() view returns(bytes4)
func (_ContractICS20Transfer *ContractICS20TransferSession) IsConsumingScheduledOp() ([4]byte, error) {
	return _ContractICS20Transfer.Contract.IsConsumingScheduledOp(&_ContractICS20Transfer.CallOpts)
}

// IsConsumingScheduledOp is a free data retrieval call binding the contract method 0x8fb36037.
//
// Solidity: function isConsumingScheduledOp() view returns(bytes4)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) IsConsumingScheduledOp() ([4]byte, error) {
	return _ContractICS20Transfer.Contract.IsConsumingScheduledOp(&_ContractICS20Transfer.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ContractICS20Transfer *ContractICS20TransferCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ContractICS20Transfer *ContractICS20TransferSession) Paused() (bool, error) {
	return _ContractICS20Transfer.Contract.Paused(&_ContractICS20Transfer.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) Paused() (bool, error) {
	return _ContractICS20Transfer.Contract.Paused(&_ContractICS20Transfer.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ContractICS20Transfer *ContractICS20TransferCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ContractICS20Transfer *ContractICS20TransferSession) ProxiableUUID() ([32]byte, error) {
	return _ContractICS20Transfer.Contract.ProxiableUUID(&_ContractICS20Transfer.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) ProxiableUUID() ([32]byte, error) {
	return _ContractICS20Transfer.Contract.ProxiableUUID(&_ContractICS20Transfer.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x1459457a.
//
// Solidity: function initialize(address ics26Router, address escrowLogic, address ibcERC20Logic, address permit2, address authority) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) Initialize(opts *bind.TransactOpts, ics26Router common.Address, escrowLogic common.Address, ibcERC20Logic common.Address, permit2 common.Address, authority common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "initialize", ics26Router, escrowLogic, ibcERC20Logic, permit2, authority)
}

// Initialize is a paid mutator transaction binding the contract method 0x1459457a.
//
// Solidity: function initialize(address ics26Router, address escrowLogic, address ibcERC20Logic, address permit2, address authority) returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) Initialize(ics26Router common.Address, escrowLogic common.Address, ibcERC20Logic common.Address, permit2 common.Address, authority common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.Initialize(&_ContractICS20Transfer.TransactOpts, ics26Router, escrowLogic, ibcERC20Logic, permit2, authority)
}

// Initialize is a paid mutator transaction binding the contract method 0x1459457a.
//
// Solidity: function initialize(address ics26Router, address escrowLogic, address ibcERC20Logic, address permit2, address authority) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) Initialize(ics26Router common.Address, escrowLogic common.Address, ibcERC20Logic common.Address, permit2 common.Address, authority common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.Initialize(&_ContractICS20Transfer.TransactOpts, ics26Router, escrowLogic, ibcERC20Logic, permit2, authority)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0x29b6eca9.
//
// Solidity: function initializeV2(address authority) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) InitializeV2(opts *bind.TransactOpts, authority common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "initializeV2", authority)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0x29b6eca9.
//
// Solidity: function initializeV2(address authority) returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) InitializeV2(authority common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.InitializeV2(&_ContractICS20Transfer.TransactOpts, authority)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0x29b6eca9.
//
// Solidity: function initializeV2(address authority) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) InitializeV2(authority common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.InitializeV2(&_ContractICS20Transfer.TransactOpts, authority)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractICS20Transfer *ContractICS20TransferTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractICS20Transfer *ContractICS20TransferSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.Multicall(&_ContractICS20Transfer.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.Multicall(&_ContractICS20Transfer.TransactOpts, data)
}

// OnAcknowledgementPacket is a paid mutator transaction binding the contract method 0x428e4e17.
//
// Solidity: function onAcknowledgementPacket((string,string,uint64,(string,string,string,string,bytes),bytes,address) msg_) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) OnAcknowledgementPacket(opts *bind.TransactOpts, msg_ IIBCAppCallbacksOnAcknowledgementPacketCallback) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "onAcknowledgementPacket", msg_)
}

// OnAcknowledgementPacket is a paid mutator transaction binding the contract method 0x428e4e17.
//
// Solidity: function onAcknowledgementPacket((string,string,uint64,(string,string,string,string,bytes),bytes,address) msg_) returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) OnAcknowledgementPacket(msg_ IIBCAppCallbacksOnAcknowledgementPacketCallback) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.OnAcknowledgementPacket(&_ContractICS20Transfer.TransactOpts, msg_)
}

// OnAcknowledgementPacket is a paid mutator transaction binding the contract method 0x428e4e17.
//
// Solidity: function onAcknowledgementPacket((string,string,uint64,(string,string,string,string,bytes),bytes,address) msg_) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) OnAcknowledgementPacket(msg_ IIBCAppCallbacksOnAcknowledgementPacketCallback) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.OnAcknowledgementPacket(&_ContractICS20Transfer.TransactOpts, msg_)
}

// OnRecvPacket is a paid mutator transaction binding the contract method 0x078c4a79.
//
// Solidity: function onRecvPacket((string,string,uint64,(string,string,string,string,bytes),address) msg_) returns(bytes)
func (_ContractICS20Transfer *ContractICS20TransferTransactor) OnRecvPacket(opts *bind.TransactOpts, msg_ IIBCAppCallbacksOnRecvPacketCallback) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "onRecvPacket", msg_)
}

// OnRecvPacket is a paid mutator transaction binding the contract method 0x078c4a79.
//
// Solidity: function onRecvPacket((string,string,uint64,(string,string,string,string,bytes),address) msg_) returns(bytes)
func (_ContractICS20Transfer *ContractICS20TransferSession) OnRecvPacket(msg_ IIBCAppCallbacksOnRecvPacketCallback) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.OnRecvPacket(&_ContractICS20Transfer.TransactOpts, msg_)
}

// OnRecvPacket is a paid mutator transaction binding the contract method 0x078c4a79.
//
// Solidity: function onRecvPacket((string,string,uint64,(string,string,string,string,bytes),address) msg_) returns(bytes)
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) OnRecvPacket(msg_ IIBCAppCallbacksOnRecvPacketCallback) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.OnRecvPacket(&_ContractICS20Transfer.TransactOpts, msg_)
}

// OnTimeoutPacket is a paid mutator transaction binding the contract method 0x5e32b6b6.
//
// Solidity: function onTimeoutPacket((string,string,uint64,(string,string,string,string,bytes),address) msg_) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) OnTimeoutPacket(opts *bind.TransactOpts, msg_ IIBCAppCallbacksOnTimeoutPacketCallback) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "onTimeoutPacket", msg_)
}

// OnTimeoutPacket is a paid mutator transaction binding the contract method 0x5e32b6b6.
//
// Solidity: function onTimeoutPacket((string,string,uint64,(string,string,string,string,bytes),address) msg_) returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) OnTimeoutPacket(msg_ IIBCAppCallbacksOnTimeoutPacketCallback) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.OnTimeoutPacket(&_ContractICS20Transfer.TransactOpts, msg_)
}

// OnTimeoutPacket is a paid mutator transaction binding the contract method 0x5e32b6b6.
//
// Solidity: function onTimeoutPacket((string,string,uint64,(string,string,string,string,bytes),address) msg_) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) OnTimeoutPacket(msg_ IIBCAppCallbacksOnTimeoutPacketCallback) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.OnTimeoutPacket(&_ContractICS20Transfer.TransactOpts, msg_)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) Pause() (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.Pause(&_ContractICS20Transfer.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) Pause() (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.Pause(&_ContractICS20Transfer.TransactOpts)
}

// SendTransfer is a paid mutator transaction binding the contract method 0x2ac3dc38.
//
// Solidity: function sendTransfer((address,uint256,string,string,string,uint64,string) msg_) returns(uint64)
func (_ContractICS20Transfer *ContractICS20TransferTransactor) SendTransfer(opts *bind.TransactOpts, msg_ IICS20TransferMsgsSendTransferMsg) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "sendTransfer", msg_)
}

// SendTransfer is a paid mutator transaction binding the contract method 0x2ac3dc38.
//
// Solidity: function sendTransfer((address,uint256,string,string,string,uint64,string) msg_) returns(uint64)
func (_ContractICS20Transfer *ContractICS20TransferSession) SendTransfer(msg_ IICS20TransferMsgsSendTransferMsg) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SendTransfer(&_ContractICS20Transfer.TransactOpts, msg_)
}

// SendTransfer is a paid mutator transaction binding the contract method 0x2ac3dc38.
//
// Solidity: function sendTransfer((address,uint256,string,string,string,uint64,string) msg_) returns(uint64)
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) SendTransfer(msg_ IICS20TransferMsgsSendTransferMsg) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SendTransfer(&_ContractICS20Transfer.TransactOpts, msg_)
}

// SendTransferWithPermit2 is a paid mutator transaction binding the contract method 0x53816a7c.
//
// Solidity: function sendTransferWithPermit2((address,uint256,string,string,string,uint64,string) msg_, ((address,uint256),uint256,uint256) permit, bytes signature) returns(uint64)
func (_ContractICS20Transfer *ContractICS20TransferTransactor) SendTransferWithPermit2(opts *bind.TransactOpts, msg_ IICS20TransferMsgsSendTransferMsg, permit ISignatureTransferPermitTransferFrom, signature []byte) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "sendTransferWithPermit2", msg_, permit, signature)
}

// SendTransferWithPermit2 is a paid mutator transaction binding the contract method 0x53816a7c.
//
// Solidity: function sendTransferWithPermit2((address,uint256,string,string,string,uint64,string) msg_, ((address,uint256),uint256,uint256) permit, bytes signature) returns(uint64)
func (_ContractICS20Transfer *ContractICS20TransferSession) SendTransferWithPermit2(msg_ IICS20TransferMsgsSendTransferMsg, permit ISignatureTransferPermitTransferFrom, signature []byte) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SendTransferWithPermit2(&_ContractICS20Transfer.TransactOpts, msg_, permit, signature)
}

// SendTransferWithPermit2 is a paid mutator transaction binding the contract method 0x53816a7c.
//
// Solidity: function sendTransferWithPermit2((address,uint256,string,string,string,uint64,string) msg_, ((address,uint256),uint256,uint256) permit, bytes signature) returns(uint64)
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) SendTransferWithPermit2(msg_ IICS20TransferMsgsSendTransferMsg, permit ISignatureTransferPermitTransferFrom, signature []byte) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SendTransferWithPermit2(&_ContractICS20Transfer.TransactOpts, msg_, permit, signature)
}

// SendTransferWithSender is a paid mutator transaction binding the contract method 0xb29c715d.
//
// Solidity: function sendTransferWithSender((address,uint256,string,string,string,uint64,string) msg_, address sender) returns(uint64)
func (_ContractICS20Transfer *ContractICS20TransferTransactor) SendTransferWithSender(opts *bind.TransactOpts, msg_ IICS20TransferMsgsSendTransferMsg, sender common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "sendTransferWithSender", msg_, sender)
}

// SendTransferWithSender is a paid mutator transaction binding the contract method 0xb29c715d.
//
// Solidity: function sendTransferWithSender((address,uint256,string,string,string,uint64,string) msg_, address sender) returns(uint64)
func (_ContractICS20Transfer *ContractICS20TransferSession) SendTransferWithSender(msg_ IICS20TransferMsgsSendTransferMsg, sender common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SendTransferWithSender(&_ContractICS20Transfer.TransactOpts, msg_, sender)
}

// SendTransferWithSender is a paid mutator transaction binding the contract method 0xb29c715d.
//
// Solidity: function sendTransferWithSender((address,uint256,string,string,string,uint64,string) msg_, address sender) returns(uint64)
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) SendTransferWithSender(msg_ IICS20TransferMsgsSendTransferMsg, sender common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SendTransferWithSender(&_ContractICS20Transfer.TransactOpts, msg_, sender)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address newAuthority) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) SetAuthority(opts *bind.TransactOpts, newAuthority common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "setAuthority", newAuthority)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address newAuthority) returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) SetAuthority(newAuthority common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SetAuthority(&_ContractICS20Transfer.TransactOpts, newAuthority)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address newAuthority) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) SetAuthority(newAuthority common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SetAuthority(&_ContractICS20Transfer.TransactOpts, newAuthority)
}

// SetCustomERC20 is a paid mutator transaction binding the contract method 0xa1d28f57.
//
// Solidity: function setCustomERC20(string denom, address token) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) SetCustomERC20(opts *bind.TransactOpts, denom string, token common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "setCustomERC20", denom, token)
}

// SetCustomERC20 is a paid mutator transaction binding the contract method 0xa1d28f57.
//
// Solidity: function setCustomERC20(string denom, address token) returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) SetCustomERC20(denom string, token common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SetCustomERC20(&_ContractICS20Transfer.TransactOpts, denom, token)
}

// SetCustomERC20 is a paid mutator transaction binding the contract method 0xa1d28f57.
//
// Solidity: function setCustomERC20(string denom, address token) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) SetCustomERC20(denom string, token common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SetCustomERC20(&_ContractICS20Transfer.TransactOpts, denom, token)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) Unpause() (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.Unpause(&_ContractICS20Transfer.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) Unpause() (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.Unpause(&_ContractICS20Transfer.TransactOpts)
}

// UpgradeEscrowTo is a paid mutator transaction binding the contract method 0xaaa2c343.
//
// Solidity: function upgradeEscrowTo(address newEscrowLogic) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) UpgradeEscrowTo(opts *bind.TransactOpts, newEscrowLogic common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "upgradeEscrowTo", newEscrowLogic)
}

// UpgradeEscrowTo is a paid mutator transaction binding the contract method 0xaaa2c343.
//
// Solidity: function upgradeEscrowTo(address newEscrowLogic) returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) UpgradeEscrowTo(newEscrowLogic common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.UpgradeEscrowTo(&_ContractICS20Transfer.TransactOpts, newEscrowLogic)
}

// UpgradeEscrowTo is a paid mutator transaction binding the contract method 0xaaa2c343.
//
// Solidity: function upgradeEscrowTo(address newEscrowLogic) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) UpgradeEscrowTo(newEscrowLogic common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.UpgradeEscrowTo(&_ContractICS20Transfer.TransactOpts, newEscrowLogic)
}

// UpgradeIBCERC20To is a paid mutator transaction binding the contract method 0x06ab20bc.
//
// Solidity: function upgradeIBCERC20To(address newIBCERC20Logic) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) UpgradeIBCERC20To(opts *bind.TransactOpts, newIBCERC20Logic common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "upgradeIBCERC20To", newIBCERC20Logic)
}

// UpgradeIBCERC20To is a paid mutator transaction binding the contract method 0x06ab20bc.
//
// Solidity: function upgradeIBCERC20To(address newIBCERC20Logic) returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) UpgradeIBCERC20To(newIBCERC20Logic common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.UpgradeIBCERC20To(&_ContractICS20Transfer.TransactOpts, newIBCERC20Logic)
}

// UpgradeIBCERC20To is a paid mutator transaction binding the contract method 0x06ab20bc.
//
// Solidity: function upgradeIBCERC20To(address newIBCERC20Logic) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) UpgradeIBCERC20To(newIBCERC20Logic common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.UpgradeIBCERC20To(&_ContractICS20Transfer.TransactOpts, newIBCERC20Logic)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.UpgradeToAndCall(&_ContractICS20Transfer.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.UpgradeToAndCall(&_ContractICS20Transfer.TransactOpts, newImplementation, data)
}

// ContractICS20TransferAuthorityUpdatedIterator is returned from FilterAuthorityUpdated and is used to iterate over the raw logs and unpacked data for AuthorityUpdated events raised by the ContractICS20Transfer contract.
type ContractICS20TransferAuthorityUpdatedIterator struct {
	Event *ContractICS20TransferAuthorityUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractICS20TransferAuthorityUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferAuthorityUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractICS20TransferAuthorityUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractICS20TransferAuthorityUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferAuthorityUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferAuthorityUpdated represents a AuthorityUpdated event raised by the ContractICS20Transfer contract.
type ContractICS20TransferAuthorityUpdated struct {
	Authority common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAuthorityUpdated is a free log retrieval operation binding the contract event 0x2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad.
//
// Solidity: event AuthorityUpdated(address authority)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterAuthorityUpdated(opts *bind.FilterOpts) (*ContractICS20TransferAuthorityUpdatedIterator, error) {

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "AuthorityUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferAuthorityUpdatedIterator{contract: _ContractICS20Transfer.contract, event: "AuthorityUpdated", logs: logs, sub: sub}, nil
}

// WatchAuthorityUpdated is a free log subscription operation binding the contract event 0x2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad.
//
// Solidity: event AuthorityUpdated(address authority)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchAuthorityUpdated(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferAuthorityUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "AuthorityUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferAuthorityUpdated)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "AuthorityUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAuthorityUpdated is a log parse operation binding the contract event 0x2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad.
//
// Solidity: event AuthorityUpdated(address authority)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseAuthorityUpdated(log types.Log) (*ContractICS20TransferAuthorityUpdated, error) {
	event := new(ContractICS20TransferAuthorityUpdated)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "AuthorityUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS20TransferIBCERC20ContractCreatedIterator is returned from FilterIBCERC20ContractCreated and is used to iterate over the raw logs and unpacked data for IBCERC20ContractCreated events raised by the ContractICS20Transfer contract.
type ContractICS20TransferIBCERC20ContractCreatedIterator struct {
	Event *ContractICS20TransferIBCERC20ContractCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractICS20TransferIBCERC20ContractCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferIBCERC20ContractCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractICS20TransferIBCERC20ContractCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractICS20TransferIBCERC20ContractCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferIBCERC20ContractCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferIBCERC20ContractCreated represents a IBCERC20ContractCreated event raised by the ContractICS20Transfer contract.
type ContractICS20TransferIBCERC20ContractCreated struct {
	ContractAddress common.Address
	FullDenomPath   string
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterIBCERC20ContractCreated is a free log retrieval operation binding the contract event 0x6031fab685dd6d86e4dbac9a69eae347145f332c95b3a0d728d3730fc5233d62.
//
// Solidity: event IBCERC20ContractCreated(address indexed contractAddress, string fullDenomPath)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterIBCERC20ContractCreated(opts *bind.FilterOpts, contractAddress []common.Address) (*ContractICS20TransferIBCERC20ContractCreatedIterator, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "IBCERC20ContractCreated", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferIBCERC20ContractCreatedIterator{contract: _ContractICS20Transfer.contract, event: "IBCERC20ContractCreated", logs: logs, sub: sub}, nil
}

// WatchIBCERC20ContractCreated is a free log subscription operation binding the contract event 0x6031fab685dd6d86e4dbac9a69eae347145f332c95b3a0d728d3730fc5233d62.
//
// Solidity: event IBCERC20ContractCreated(address indexed contractAddress, string fullDenomPath)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchIBCERC20ContractCreated(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferIBCERC20ContractCreated, contractAddress []common.Address) (event.Subscription, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "IBCERC20ContractCreated", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferIBCERC20ContractCreated)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "IBCERC20ContractCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIBCERC20ContractCreated is a log parse operation binding the contract event 0x6031fab685dd6d86e4dbac9a69eae347145f332c95b3a0d728d3730fc5233d62.
//
// Solidity: event IBCERC20ContractCreated(address indexed contractAddress, string fullDenomPath)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseIBCERC20ContractCreated(log types.Log) (*ContractICS20TransferIBCERC20ContractCreated, error) {
	event := new(ContractICS20TransferIBCERC20ContractCreated)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "IBCERC20ContractCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS20TransferInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractICS20Transfer contract.
type ContractICS20TransferInitializedIterator struct {
	Event *ContractICS20TransferInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractICS20TransferInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractICS20TransferInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractICS20TransferInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferInitialized represents a Initialized event raised by the ContractICS20Transfer contract.
type ContractICS20TransferInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractICS20TransferInitializedIterator, error) {

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferInitializedIterator{contract: _ContractICS20Transfer.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferInitialized)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseInitialized(log types.Log) (*ContractICS20TransferInitialized, error) {
	event := new(ContractICS20TransferInitialized)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS20TransferPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ContractICS20Transfer contract.
type ContractICS20TransferPausedIterator struct {
	Event *ContractICS20TransferPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractICS20TransferPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractICS20TransferPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractICS20TransferPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferPaused represents a Paused event raised by the ContractICS20Transfer contract.
type ContractICS20TransferPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterPaused(opts *bind.FilterOpts) (*ContractICS20TransferPausedIterator, error) {

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferPausedIterator{contract: _ContractICS20Transfer.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferPaused) (event.Subscription, error) {

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferPaused)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParsePaused(log types.Log) (*ContractICS20TransferPaused, error) {
	event := new(ContractICS20TransferPaused)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS20TransferUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ContractICS20Transfer contract.
type ContractICS20TransferUnpausedIterator struct {
	Event *ContractICS20TransferUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractICS20TransferUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractICS20TransferUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractICS20TransferUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferUnpaused represents a Unpaused event raised by the ContractICS20Transfer contract.
type ContractICS20TransferUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterUnpaused(opts *bind.FilterOpts) (*ContractICS20TransferUnpausedIterator, error) {

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferUnpausedIterator{contract: _ContractICS20Transfer.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferUnpaused) (event.Subscription, error) {

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferUnpaused)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseUnpaused(log types.Log) (*ContractICS20TransferUnpaused, error) {
	event := new(ContractICS20TransferUnpaused)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS20TransferUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ContractICS20Transfer contract.
type ContractICS20TransferUpgradedIterator struct {
	Event *ContractICS20TransferUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractICS20TransferUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractICS20TransferUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractICS20TransferUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferUpgraded represents a Upgraded event raised by the ContractICS20Transfer contract.
type ContractICS20TransferUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ContractICS20TransferUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferUpgradedIterator{contract: _ContractICS20Transfer.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferUpgraded)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseUpgraded(log types.Log) (*ContractICS20TransferUpgraded, error) {
	event := new(ContractICS20TransferUpgraded)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
