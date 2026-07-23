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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEscrow\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEscrowBeacon\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIBCERC20Beacon\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPermit2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ibcERC20Contract\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ibcERC20Denom\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ics26\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"ics26Router\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"escrowLogic\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"ibcERC20Logic\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permit2\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializeV2\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isConsumingScheduledOp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onAcknowledgementPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIIBCAppCallbacks.OnAcknowledgementPacketCallback\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destinationClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"acknowledgement\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayer\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onRecvPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIIBCAppCallbacks.OnRecvPacketCallback\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destinationClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"relayer\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onTimeoutPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIIBCAppCallbacks.OnTimeoutPacketCallback\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destinationClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"relayer\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sendTransfer\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS20TransferMsgs.SendTransferMsg\",\"components\":[{\"name\":\"denom\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"memo\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendTransferWithPermit2\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS20TransferMsgs.SendTransferMsg\",\"components\":[{\"name\":\"denom\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"memo\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"permit\",\"type\":\"tuple\",\"internalType\":\"structISignatureTransfer.PermitTransferFrom\",\"components\":[{\"name\":\"permitted\",\"type\":\"tuple\",\"internalType\":\"structISignatureTransfer.TokenPermissions\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendTransferWithSender\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS20TransferMsgs.SendTransferMsg\",\"components\":[{\"name\":\"denom\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"memo\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAuthority\",\"inputs\":[{\"name\":\"newAuthority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setCustomERC20\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setIBCERC20Metadata\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decimals_\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeEscrowTo\",\"inputs\":[{\"name\":\"newEscrowLogic\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeIBCERC20To\",\"inputs\":[{\"name\":\"newIBCERC20Logic\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCERC20ContractCreated\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fullDenomPath\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCSenderAckPacketCallbackError\",\"inputs\":[{\"name\":\"callbackAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCSenderAckPacketCallbackError\",\"inputs\":[{\"name\":\"callbackAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCSenderTimeoutPacketCallbackError\",\"inputs\":[{\"name\":\"callbackAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCSenderTimeoutPacketCallbackError\",\"inputs\":[{\"name\":\"callbackAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessManagedInvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessManagedRequiredDelay\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"AccessManagedUnauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ICS20DenomAlreadyExists\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20DenomNotFound\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20EscrowNotFound\",\"inputs\":[{\"name\":\"clientID\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20InvalidAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20InvalidAmount\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ICS20InvalidPort\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20Permit2TokenMismatch\",\"inputs\":[{\"name\":\"permitToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sentToken\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ICS20TokenAlreadyExists\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ICS20UnauthorizedPacketSender\",\"inputs\":[{\"name\":\"packetSender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ICS20UnexpectedERC20Balance\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ICS20UnexpectedEncoding\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20UnexpectedVersion\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60a080604052346100c257306080525f516020614d765f395f51905f525460ff8160401c166100b3576002600160401b03196001600160401b03821601610060575b604051614caf90816100c7823960805181818161123b01526112eb0152f35b6001600160401b0319166001600160401b039081175f516020614d765f395f51905f525581527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f80610041565b63f92ee8a960e01b5f5260045ffd5b5f80fdfe60806040526004361015610011575f80fd5b5f5f3560e01c806306ab20bc146128ae578063078c4a791461210f5780631459457a14611d7d5780631bbf2e2314611d375780631e5150e414611cf157806329b6eca914611adb5780632ac3dc3814611a3a5780633f4ba83a14611979578063428e4e171461153c5780634f1ef2861461129a57806352d1902d1461122057806353816a7c14610f645780635c975abb14610f225780635e32b6b614610d0a57806362bfa17214610bcd5780637a9e5e4b14610b3c578063826cae7a14610af65780638456cb5914610a5b5780638fb36037146109f1578063969631d514610975578063a1d28f57146106dd578063a50ee2b41461067f578063aaa2c343146105d3578063ac9650d81461042f578063ad3cb1cc146103ce578063b29c715d1461025b578063bf7e214f14610228578063d413227d146101e25763e163b1af14610159575f80fd5b346101df5760203660031901126101df576101db6101c06101c76101b461017e612957565b6001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b60405192838092612ee6565b0382612a1f565b6040519182916020835260208301906129cb565b0390f35b80fd5b50346101df57806003193601126101df5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416604051908152f35b50346101df57806003193601126101df5760206001600160a01b035f516020614c635f395f51905f525416604051908152f35b50346101df5760403660031901126101df5760043567ffffffffffffffff81116103b357806004019060e060031982360301126103ca5761029a61296d565b916102a3613275565b6102ab613201565b6102b53633613049565b60248201359182156103b7576102e86102e36102dc60646001600160a01b03940185612b13565b3691612a5d565b6132e8565b16916102fe816102f784612e32565b85336134ad565b8461030883612e32565b91843b156103b35760405163b4f22eb760e01b81526001600160a01b0393909316600484015233602484015260448301528160648183875af180156103a85761038f575b60208561035a868686613689565b907f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d67ffffffffffffffff60405191168152f35b61039a858092612a1f565b6103a4575f61034c565b8380fd5b6040513d87823e3d90fd5b5080fd5b6024856304f6df8d60e41b815280600452fd5b8280fd5b50346101df57806003193601126101df57506101db6040516103f1604082612a1f565b600581527f352e302e3000000000000000000000000000000000000000000000000000000060208201526040519182916020835260208301906129cb565b50346101df5760203660031901126101df5760043567ffffffffffffffff81116103b357366023820112156103b35780600401359067ffffffffffffffff82116103ca57602481013660248460051b840101116103a45760405160206104958183612a1f565b85825280820192601f1982013685376104ad86613009565b946104bb6040519687612a1f565b868652601f196104ca88613009565b0183895b8281106105c357505050875b87811015610564576001906105488a8089896105348a61050260248960051b8c01018c612b13565b9190946040519483869484860198893784019083820190898252519283915e010185815203601f198101835282612a1f565b5190305af4610541613b98565b9030614076565b610552828a613021565b5261055d8189613021565b50016104da565b83898860405191838301848452825180915260408401948060408360051b870101940192955b8287106105975785850386f35b9091929382806105b3600193603f198a820301865288516129cb565b960192019601959291909261058a565b606082828b0101520184906104ce565b50346101df5760203660031901126101df57806105ee612957565b6105f83633613049565b6001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805541690813b1561067b576001600160a01b0360248492836040519586948593631b2ce7f360e11b85521660048401525af180156106705761065f5750f35b8161066991612a1f565b6101df5780f35b6040513d84823e3d90fd5b5050fd5b50346101df5760203660031901126101df576004359067ffffffffffffffff82116101df5760206106d56106b63660048601612a93565b6001600160a01b036106ca82849594612e46565b541692831515612e7e565b604051908152f35b50346101df5760403660031901126101df5760043567ffffffffffffffff81116103b35761070f903690600401612a93565b61071a92919261296d565b6107243633613049565b6001600160a01b036107368386612e46565b54166109395761077f610779826001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b54612eae565b156107ba826001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b90156108f75750610820906107cf8386612e46565b6001600160a01b03808316166001600160a01b03198254161790556001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9067ffffffffffffffff81116108e3576108448161083e8454612eae565b84612f67565b82601f82116001146108835781908495610873949592610878575b50508160011b915f199060031b1c19161790565b905580f35b013590505f8061085f565b601f198216948385526020852091855b8781106108cb5750836001959697106108b2575b505050811b01905580f35b01355f19600384901b60f8161c191690555f80806108a7565b90926020600181928686013581550194019101610893565b602483634e487b7160e01b81526041600452fd5b610935906040519182917f778769c4000000000000000000000000000000000000000000000000000000008352602060048401526024830190612ee6565b0390fd5b6040517f0c0ef5340000000000000000000000000000000000000000000000000000000081526020600482015280610935602482018588612b81565b50346101df5760203660031901126101df5760043567ffffffffffffffff81116103b357602091826109b36001600160a01b03933690600401612a93565b925082604051938492833781017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8008152030190205416604051908152f35b50346101df57806003193601126101df575f516020614c635f395f51905f525460a01c60ff1615610a53575060207f8fb36037000000000000000000000000000000000000000000000000000000005b6001600160e01b031960405191168152f35b602090610a41565b50346101df57806003193601126101df57610a763633613049565b610a7e613275565b600160ff197fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005416177fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a180f35b50346101df57806003193601126101df5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8045416604051908152f35b50346101df5760203660031901126101df57610b56612957565b6001600160a01b035f516020614c635f395f51905f5254163303610bbb57803b15610b8757610b8490613df6565b80f35b7fc2f31e5e0000000000000000000000000000000000000000000000000000000082526001600160a01b0316600452602490fd5b60248262d1953b60e31b815233600452fd5b50346101df5760803660031901126101df578060043567ffffffffffffffff8111610d0757610c00903690600401612a93565b60243567ffffffffffffffff8111610d0457610c20903690600401612a93565b91909260443567ffffffffffffffff8111610d0057610c43903690600401612a93565b94906064359360ff8516809503610cfc57610c7f90610c623633613049565b6001600160a01b03610c748287612e46565b541694851515612e7e565b823b15610cf8578694610ccf94610ce18793604051998a98899788967f37d2c2f4000000000000000000000000000000000000000000000000000000008852606060048901526064880191612b81565b85810360031901602487015291612b81565b90604483015203925af180156106705761065f5750f35b8680fd5b8780fd5b8580fd5b50505b50fd5b50346101df57610d1936612997565b610d4f336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314612ac1565b610d57613201565b606081019082610da8610d82610d7a610d708686612afe565b6080810190612b13565b810190612c9e565b610d95610d8f8686612afe565b80612b13565b90610da08680612b13565b929091613bc7565b9050610db3816141ef565b610ddf575b50807f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d80f35b6001600160a01b031692833b156103b357610eab82916080946001600160a01b03610eb260405197889586957f5e32b6b600000000000000000000000000000000000000000000000000000000875260206004880152610e9a610e78610e59610e488980613aa8565b60a060248d015260c48c0191612b81565b610e6660208a018a613aa8565b8b83036023190160448d015290612b81565b9167ffffffffffffffff610e8e60408a01613ada565b1660648a015287613aef565b878203602319016084890152613b03565b9301612983565b1660a483015203818387620186a0f19182610f0d575b5050610f07577f02c075b43cbb5fee7c14de7aea6a6ddde6e3646144712fd9c7159705152537e6610efa6101c7613b98565b0390a25b5f808281610db8565b50610efe565b81610f1791612a1f565b6103ca57825f610ec8565b50346101df57806003193601126101df57602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b50346101df5760c03660031901126101df5760043567ffffffffffffffff81116103b357806004019060e060031982360301126103ca5760803660231901126103ca5760a43567ffffffffffffffff81116103a457610fc7903690600401612a93565b610fd2929192613275565b610fda613201565b602482013591821561120d57610fee612e1c565b6001600160a01b038061100088612e32565b1691161461100c612e1c565b61101587612e32565b91156111d35750506102e36102dc6064611030930187612b13565b926001600160a01b03807f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065416941693604051906040820182811067ffffffffffffffff8211176111bf579088949392916040528682526020820191868352813b15610d0057856001600160a01b039161111f8296604051988997889687957f30f28b7a000000000000000000000000000000000000000000000000000000008752816110db61296d565b166004880152604435602488015260643560448801526084356064880152511660848601525160a48501523360c485015261010060e4850152610104840191612b81565b03925af18015610670576111aa575b5061113884612e32565b91833b156103b35760405163b4f22eb760e01b81526001600160a01b0393909316600484015233602484015260448301528160648183865af1801561119f5761118a575b60208461035a338587613689565b611195848092612a1f565b6103ca575f61117c565b6040513d86823e3d90fd5b816111b491612a1f565b6103a457835f61112e565b602489634e487b7160e01b81526041600452fd5b7fe36164780000000000000000000000000000000000000000000000000000000088526001600160a01b0390811660045216602452604486fd5b6024866304f6df8d60e41b815280600452fd5b50346101df57806003193601126101df576001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016300361128b5760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b8063703e46dd60e11b60049252fd5b5060403660031901126101df576112af612957565b9060243567ffffffffffffffff81116103b357366023820112156103b3576112e1903690602481600401359101612a5d565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016803014908115611507575b506114f8576113253633613049565b6001600160a01b03831690604051937f52d1902d000000000000000000000000000000000000000000000000000000008552602085600481865afa809585966114c0575b506113815760248484634c9c8ce360e01b8252600452fd5b9091847f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc81036114955750813b1561148357806001600160a01b03197f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5416177f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b8480a28151839015611450578083602061144c95519101845af4611446613b98565b91614076565b5080f35b5050503461145b5780f35b807fb398979f0000000000000000000000000000000000000000000000000000000060049252fd5b634c9c8ce360e01b8452600452602483fd5b7faa1d49a4000000000000000000000000000000000000000000000000000000008552600452602484fd5b9095506020813d6020116114f0575b816114dc60209383612a1f565b810103126114ec5751945f611369565b8480fd5b3d91506114cf565b60048263703e46dd60e11b8152fd5b90506001600160a01b037f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc541614155f611316565b50346101df5760203660031901126101df5760043567ffffffffffffffff81116103b357806004019060c060031982360301126103ca576115a9336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314612ac1565b6115b1613201565b6064810191836115c7610d7a610d708685612afe565b9260848101926115da6102dc8583612b13565b602081519101209460409560208088516115f48a82612a1f565b818152017f4774d4a575993f963b1c06573736617a457abef8589178db8d10c94b4ab511ab815220145f146117ea5761163f90611634610d8f8985612afe565b90610da08580612b13565b905061164a816141ef565b61167d575b505050505050505b807f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d80f35b6001600160a01b031695863b156103a457839286519586938493638dfcd9ad60e01b85528660048601528960248601526116b78280613aa8565b6044870160c090526101048701906116ce92612b81565b6116db6024860184613aa8565b8783036043190160648901526116f19291612b81565b906116fe60448601613ada565b67ffffffffffffffff1660848701526117179083613aef565b8582036043190160a487015261172d9190613b03565b9161173791613aa8565b8483036043190160c486015261174d9291612b81565b9060a40161175a90612983565b6001600160a01b031660e483015203818388620186a0f191826117d5575b50506117ce576117be7f898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e0916117ab613b98565b90519182916020835260208301906129cb565b0390a25b5f80808381808061164f565b50506117c2565b816117df91612a1f565b6103a457835f611778565b602001516118036117fc825183613ea8565b9290612d6b565b61180c816141ef565b61181d575b50505050505050611657565b6001600160a01b031695863b156103a457839286519586938493638dfcd9ad60e01b855260048501600190528960248601526118598280613aa8565b6044870160c0905261010487019061187092612b81565b61187d6024860184613aa8565b8783036043190160648901526118939291612b81565b906118a060448601613ada565b67ffffffffffffffff1660848701526118b99083613aef565b8582036043190160a48701526118cf9190613b03565b916118d991613aa8565b8483036043190160c48601526118ef9291612b81565b9060a4016118fc90612983565b6001600160a01b031660e483015203818388620186a0f19182611964575b505061195d5761194d7f898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e0916117ab613b98565b0390a25b5f808083818080611811565b5050611951565b8161196e91612a1f565b6103a457835f61191a565b50346101df57806003193601126101df576119943633613049565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff811615611a125760ff19167fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a180f35b6004827f8dfc202b000000000000000000000000000000000000000000000000000000008152fd5b50346101df5760203660031901126101df5760043567ffffffffffffffff81116103b357806004019060e060031982360301126103ca57611a79613275565b611a81613201565b6024810135908115611ac857611aa86102e36102dc60646001600160a01b03940186612b13565b1690611abe81611ab785612e32565b84336134ad565b8361113884612e32565b6024846304f6df8d60e41b815280600452fd5b50346101df5760203660031901126101df57611af5612957565b5f516020614c835f395f51905f525467ffffffffffffffff81169060018203611ce25760401c60ff16908115611cd6575b50611cc757611b5c600267ffffffffffffffff195f516020614c835f395f51905f525416175f516020614c835f395f51905f5255565b6801000000000000000068ff0000000000000000195f516020614c835f395f51905f525416175f516020614c835f395f51905f5255602460206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416604051928380927f24d7806c0000000000000000000000000000000000000000000000000000000082523360048301525afa908115611cbc578391611c7d575b5090611c10611c25923390612ac1565b611c18613f43565b611c20613f43565b613df6565b68ff0000000000000000195f516020614c835f395f51905f5254165f516020614c835f395f51905f52557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a180f35b90506020813d602011611cb4575b81611c9860209383612a1f565b810103126103ca57519081151582036103ca5790611c10611c00565b3d9150611c8b565b6040513d85823e3d90fd5b60048263f92ee8a960e01b8152fd5b6002915010155f611b26565b60048463f92ee8a960e01b8152fd5b50346101df57806003193601126101df5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8055416604051908152f35b50346101df57806003193601126101df5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065416604051908152f35b50346101df5760a03660031901126101df57611d97612957565b611d9f61296d565b6044356001600160a01b03811681036103a457606435926001600160a01b0384168094036114ec576084356001600160a01b0381168103610d00575f516020614c835f395f51905f525467ffffffffffffffff811690816121005760401c60ff169081156120f4575b506120e55790611e986001600160a01b0392611e4b600267ffffffffffffffff195f516020614c835f395f51905f525416175f516020614c835f395f51905f5255565b6801000000000000000068ff0000000000000000195f516020614c835f395f51905f525416175f516020614c835f395f51905f5255611e88613f43565b611e90613f43565b611c10613f43565b166001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80355604051916104009283810181811067ffffffffffffffff8211176120d157611f2a82916143bb95878785396001600160a01b0316815230602082015260400190565b039086f080156103a8576001600160a01b03166001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8045416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80455604051928084019284841067ffffffffffffffff8511176120d15791849391611fc89385396001600160a01b0316815230602082015260400190565b039083f08015610670576001600160a01b03166001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8055416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805556001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065568ff0000000000000000195f516020614c835f395f51905f5254165f516020614c835f395f51905f52557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a180f35b602487634e487b7160e01b81526041600452fd5b60048663f92ee8a960e01b8152fd5b6002915010155f611e08565b60048863f92ee8a960e01b8152fd5b50346101df5761211e36612997565b90612155336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314612ac1565b61215d613201565b612165613275565b60608201916121846102dc61217a8584612afe565b6040810190612b13565b60208151910120612193612b46565b60208151910120146121a3612b46565b6121b061217a8685612afe565b909215612878575050506122026121cd6102dc610d8f8685612afe565b602081519101206121dc612bc9565b60208151910120146121ec612bc9565b906121fa610d8f8786612afe565b929091612c04565b61221c6102dc6122128584612afe565b6060810190612b13565b6020815191012061222b612c48565b602081519101201461223b612c48565b6122486122128685612afe565b9092156128425750505061229c61226f6102dc6122658685612afe565b6020810190612b13565b6020815191012061227e612bc9565b602081519101201461228e612bc9565b906121fa6122658786612afe565b6122ac610d7a610d708584612afe565b9260608401805115611ac85760408501916122ea8351936122d86122d1865187613ea8565b9690612d6b565b516001600160a01b0385161515612d6b565b60208401936122ff6102e36102dc8784612b13565b96519461232e612312610d8f8585612afe565b61232961231f8680612b13565b9390923691612a5d565b613457565b92865184511080612832575b156124c75750509051845195968795909250906020908781189088110287188061236c612367828b612fb6565b612fd7565b9803920101602087015e612381855186613ea8565b9590158681156124b5575b5061248c575b506001600160a01b03905b16905193813b156103a4576040517f0779afe60000000000000000000000000000000000000000000000000000000081526001600160a01b039182166004820152921660248301526044820193909352918290606490829084905af1801561067057612477575b6101db8260405190612417604083612a1f565b601182527f7b22726573756c74223a2241513d3d227d00000000000000000000000000000060208301527f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d6040519182916020835260208301906129cb565b612482828092612a1f565b6101df575f612404565b94506001600160a01b03906124af826124a488612db1565b541696871515612def565b90612392565b6001600160a01b03915016155f61238c565b61253193506124fc83602098949361232961231f6124f46124eb8d98978998612afe565b87810190612b13565b939094612b13565b926040519784899551918291018487015e8401908282018a8152815193849201905e010186815203601f198101855284612a1f565b6001600160a01b038516946001600160a01b0361254d85612db1565b54169384156125d5575b5084956001600160a01b038596959616908351823b15610cf8576040516340c10f1960e01b81526001600160a01b0392909216600483015260248201529085908290604490829084905af19081156103a85785916125c0575b50506001600160a01b039061239d565b816125ca91612a1f565b6103a457835f6125b0565b93506001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80454166040517f4571e3a60000000000000000000000000000000000000000000000000000000060208201523060248201528760448201526060606482015261265d8161264f60848201896129cb565b03601f198101835282612a1f565b604051916104a88084019084821067ffffffffffffffff83111761281e579184939161268d936147bb86396132c8565b039086f080156103a8576001600160a01b0316956126aa85612db1565b6001600160a01b0388166001600160a01b03198254161790556126fd876001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9685519767ffffffffffffffff891161280a576127248961271e8354612eae565b83612f67565b602098601f81116001146127a85780612755918a9b8b9a9b9161279d575b508160011b915f199060031b1c19161790565b90555b7f6031fab685dd6d86e4dbac9a69eae347145f332c95b3a0d728d3730fc5233d6261279282986040519182916020835260208301906129cb565b0390a2959493612557565b90508a01515f612742565b818952898920601f1982168a5b8181106127f25750908a9b83600194939c9b9c106127da575b5050811b019055612758565b8b01515f1960f88460031b161c191690555f806127ce565b8a8d0151835560209c8d019c600190930192016127b5565b602488634e487b7160e01b81526041600452fd5b60248a634e487b7160e01b81526041600452fd5b5061283d8488613e56565b61233a565b610935906040519384937fd1ca953a00000000000000000000000000000000000000000000000000000000855260048501612ba1565b610935906040519384937f094af3b800000000000000000000000000000000000000000000000000000000855260048501612ba1565b5034612953576020366003190112612953576128c8612957565b6128d23633613049565b6001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f804541690813b15612953576001600160a01b0360245f92836040519586948593631b2ce7f360e11b85521660048401525af180156129485761293a575080f35b61294691505f90612a1f565b005b6040513d5f823e3d90fd5b5f80fd5b600435906001600160a01b038216820361295357565b602435906001600160a01b038216820361295357565b35906001600160a01b038216820361295357565b6020600319820112612953576004359067ffffffffffffffff82116129535760a09082900360031901126129535760040190565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b60a0810190811067ffffffffffffffff821117612a0b57604052565b634e487b7160e01b5f52604160045260245ffd5b90601f8019910116810190811067ffffffffffffffff821117612a0b57604052565b67ffffffffffffffff8111612a0b57601f01601f191660200190565b929192612a6982612a41565b91612a776040519384612a1f565b829481845281830111612953578281602093845f960137010152565b9181601f840112156129535782359167ffffffffffffffff8311612953576020838186019501011161295357565b15612ac95750565b6001600160a01b03907f2ecb3242000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b903590609e1981360301821215612953570190565b903590601e1981360301821215612953570180359067ffffffffffffffff82116129535760200191813603831361295357565b60405190612b55604083612a1f565b600782527f69637332302d31000000000000000000000000000000000000000000000000006020830152565b908060209392818452848401375f828201840152601f01601f1916010190565b91612bb8612bc694926040855260408501906129cb565b926020818503910152612b81565b90565b60405190612bd8604083612a1f565b600882527f7472616e736665720000000000000000000000000000000000000000000000006020830152565b9290919215612c1257505050565b610935906040519384937f5d3a3cdd00000000000000000000000000000000000000000000000000000000855260048501612ba1565b60405190612c57604083612a1f565b601a82527f6170706c69636174696f6e2f782d736f6c69646974792d6162690000000000006020830152565b9080601f8301121561295357816020612bc693359101612a5d565b6020818303126129535780359067ffffffffffffffff8211612953570160a0818303126129535760405191612cd2836129ef565b813567ffffffffffffffff81116129535781612cef918401612c83565b8352602082013567ffffffffffffffff81116129535781612d11918401612c83565b6020840152604082013567ffffffffffffffff81116129535781612d36918401612c83565b604084015260608201356060840152608082013567ffffffffffffffff811161295357612d639201612c83565b608082015290565b15612d735750565b610935906040519182917f3fed5d870000000000000000000000000000000000000000000000000000000083526020600484015260248301906129cb565b60208091604051928184925191829101835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80181520301902090565b15612df75750565b6109359060405191829163e1275e2f60e01b83526020600484015260248301906129cb565b6024356001600160a01b03811681036129535790565b356001600160a01b03811681036129535790565b60209082604051938492833781017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80181520301902090565b91909115612e8a575050565b61093560405192839263e1275e2f60e01b8452602060048501526024840191612b81565b90600182811c92168015612edc575b6020831014612ec857565b634e487b7160e01b5f52602260045260245ffd5b91607f1691612ebd565b5f9291815491612ef583612eae565b8083529260018116908115612f4a5750600114612f1157505050565b5f9081526020812093945091925b838310612f30575060209250010190565b600181602092949394548385870101520191019190612f1f565b915050602093945060ff929192191683830152151560051b010190565b601f8211612f7457505050565b5f5260205f20906020601f840160051c83019310612fac575b601f0160051c01905b818110612fa1575050565b5f8155600101612f96565b9091508190612f8d565b91908203918211612fc357565b634e487b7160e01b5f52601160045260245ffd5b90612fe182612a41565b612fee6040519182612a1f565b8281528092612fff601f1991612a41565b0190602036910137565b67ffffffffffffffff8111612a0b5760051b60200190565b80518210156130355760209160051b010190565b634e487b7160e01b5f52603260045260245ffd5b5f516020614c635f395f51905f5254916001600160a01b0383169281600411612953575f5f9060405f8151966001600160a01b0360208901917fb700961300000000000000000000000000000000000000000000000000000000835216978860248201523060448201526001600160e01b03198335166064820152606481526130d3608482612a1f565b828052826020525190895afa6131ee575b156130f1575b5050505050565b63ffffffff16156131dc5760ff60a01b191674010000000000000000000000000000000000000000175f516020614c635f395f51905f5255823b15612953576020925f92836040518096819582947f94c7d7ee00000000000000000000000000000000000000000000000000000000845260048401526040602484015260448301908082528085848401378181018301859052601f01601f1916010103925af18015612948576131cc575b5060ff60a01b195f516020614c635f395f51905f5254165f516020614c635f395f51905f52555f808080806130ea565b5f6131d691612a1f565b5f61319c565b8262d1953b60e31b5f5260045260245ffd5b50505f516020518060201c1502906130e4565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005c61324d5760017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d565b7f3ee5aeb5000000000000000000000000000000000000000000000000000000005f5260045ffd5b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166132a057565b7fd93c0665000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040906001600160a01b03612bc6949316815281602082015201906129cb565b604051906001600160a01b03815192602081818501958087835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f800815203019020541691821561333957505090565b7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805545f516020614c635f395f51905f52546040517f485cc9550000000000000000000000000000000000000000000000000000000060208201523060248201526001600160a01b039182166044808301919091528152939450919291166133c1606483612a1f565b604051916104a8908184019284841067ffffffffffffffff851117612a0b5784936133f0936147bb86396132c8565b03905ff08015612948576001600160a01b036020911692604051928391518091835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8008152030190206001600160a01b0382166001600160a01b031982541617905590565b6001809160208095612bc6958160405198858a9651918291018688015e850191602f60f81b8584015260218301370101602f60f81b838201520301601e19810184520182612a1f565b91908201809211612fc357565b916001600160a01b03909391931692604051926370a0823160e01b84526001600160a01b03821691826004860152602085602481895afa948515612948575f95613654575b506040517f23b872dd0000000000000000000000000000000000000000000000000000000060208281019182526001600160a01b039485166024840152939092166044820152606481018590525f919061354f816084810161264f565b519082885af115612948575f513d61364b5750833b155b61361f576020906024604051809681936370a0823160e01b835260048301525afa928315612948575f936135e9575b506135a090826134a0565b908211806135e0575b156135b2575050565b7f2fb30cfc000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b508082146135a9565b9092506020813d602011613617575b8161360560209383612a1f565b810103126129535751916135a0613595565b3d91506135f8565b837f5274afe7000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b60011415613566565b9094506020813d602011613681575b8161367060209383612a1f565b8101031261295357519360206134f2565b3d9150613663565b6136aa61369861017e83612e32565b926136b15f9460405193848092612ee6565b0383612a1f565b818051155f146139ce5750505060a06136da6136d46136cf84612e32565b613f87565b94613f87565b936136e86040840184612b13565b93909561373961371e6136fe60c0850185612b13565b9990976040519661370e886129ef565b8752602087019485523691612a5d565b9560408501968752606085019860208501358a523691612a5d565b94608084019586526001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354169561377a6060850185612b13565b96909401359467ffffffffffffffff86168096036139ca57906138346137e79493926138266137a7612bc9565b9c6137b0612bc9565b9461380d6137bc612b46565b976137fa6137c8612c48565b9a6040519c8d986020808b01525160a060408b015260e08a01906129cb565b9051888203603f190160608a01526129cb565b9051868203603f190160808801526129cb565b915160a085015251838203603f190160c08501526129cb565b03601f198101865285612a1f565b604051996138418b6129ef565b8a5260208a0152604089015260608801526080870152604051926060840184811067ffffffffffffffff8211176120d15793613964602096946138976138eb958a9567ffffffffffffffff996040523691612a5d565b835287830190815260408301998a52604051998a97889687957f4d6e7ce30000000000000000000000000000000000000000000000000000000087528b6004880152516060602488015260848701906129cb565b92511660448501525190602319848203016064850152608061395361394161392f61391f865160a0875260a08701906129cb565b8d8701518682038f8801526129cb565b604086015185820360408701526129cb565b606085015184820360608601526129cb565b9201519060808184039101526129cb565b03925af19182156139bd57819261397a57505090565b9091506020813d6020116139b5575b8161399660209383612a1f565b810103126103b357519067ffffffffffffffff821682036101df575090565b3d9150613989565b50604051903d90823e3d90fd5b8880fd5b6139ef6139dc969396612bc9565b6139e96060870187612b13565b91613457565b81518151109182613a97575b5050613a0e575b506136da60a091613f87565b6001600160a01b03613a1f84612e32565b16803b15612953576040517f9dc29fac0000000000000000000000000000000000000000000000000000000081526001600160a01b03929092166004830152602084013560248301525f908290604490829084905af180156129485715613a0257613a8d9193505f90612a1f565b5f916136da613a02565b613aa19250613e56565b5f806139fb565b9035601e198236030181121561295357016020813591019167ffffffffffffffff821161295357813603831361295357565b359067ffffffffffffffff8216820361295357565b9035609e1982360301811215612953570190565b612bc691613b8a613b7f613b64613b49613b2e613b208780613aa8565b60a0885260a0880191612b81565b613b3b6020880188613aa8565b908783036020890152612b81565b613b566040870187613aa8565b908683036040880152612b81565b613b716060860186613aa8565b908583036060870152612b81565b926080810190613aa8565b916080818503910152612b81565b3d15613bc2573d90613ba982612a41565b91613bb76040519384612a1f565b82523d5f602084013e565b606090565b909493925f926001600160a01b03604051838382376020818581017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8008152030190205416928315613db75790612329613c3693926020880151996102dc613c2f8c518d613ea8565b9c90612d6b565b83519081518151109182613da6575b505015613d4f576001600160a01b03613c5e8451612db1565b541692613c6e8151851515612def565b6060810151843b15612953576040516340c10f1960e01b81526001600160a01b038416600482015260248101919091525f8160448183895af1801561294857613d39575b506060905b0151813b156103ca576040517fb6163ac40000000000000000000000000000000000000000000000000000000081526001600160a01b0385811660048301528716602482015260448101919091529082908290606490829084905af1801561067057613d24575b50509190565b613d2f828092612a1f565b6101df5780613d1e565b613d469193505f90612a1f565b5f916060613cb2565b6001600160a01b03613d618451612db1565b5416928315613d73575b606090613cb7565b92506060835193613d886122d1865187613ea8565b613d9f856001600160a01b03835191161515612def565b9050613d6b565b613db09250613e56565b5f80613c45565b50906109356040519283927f5778f378000000000000000000000000000000000000000000000000000000008452602060048501526024840191612b81565b60206001600160a01b037f2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad9216806001600160a01b03195f516020614c635f395f51905f525416175f516020614c635f395f51905f5255604051908152a1565b8051908251809210613ea1578051918280821091180280831892141582028218906020613e866123678486612fb6565b92808285019503920101835e51902090602081519101201490565b5050505f90565b805182118015613f3c575b613f00576001821180613f08575b158015908160011b9182046002141715612fc35760280180602811612fc3578203613f00576001600160a01b0392915f613efa92614102565b90921690565b50505f905f90565b5061060f60f31b7fffff00000000000000000000000000000000000000000000000000000000000060208301511614613ec1565b505f613eb3565b60ff5f516020614c835f395f51905f525460401c1615613f5f57565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f5260045ffd5b6001600160a01b031680613f9b602a612a41565b91613fa96040519384612a1f565b602a8352613fb7602a612a41565b6020840190601f19013682378351156130355760309053825160011015613035576078602184015360295b600181116140235750613ff3575090565b7fe22e27eb000000000000000000000000000000000000000000000000000000005f52600452601460245260445ffd5b90600f81166010811015613035578451831015613035577f3031323334353637383961626364656600000000000000000000000000000000901a8483016020015360041c908015612fc3575f1901613fe2565b906140b3575080511561408b57602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b815115806140f9575b6140c4575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b156140bc565b92909260018401808511612fc3578311806141b9575b158015908160011b9182046002141715612fc35761413b905f94929394956134a0565b915b81831061414d5750505060019190565b9092919360ff6141847fff000000000000000000000000000000000000000000000000000000000000006020888601015116614348565b16600f81116141ae578160041b9180830460101490151715612fc35760019101940191929061413d565b505f94508493505050565b5061060f60f31b7fffff000000000000000000000000000000000000000000000000000000000000602086840101511614614118565b60205f604051828101906301ffc9a760e01b82526301ffc9a760e01b60248201526024815261421f604482612a1f565b519084617530fa903d5f51908361433c575b5082614332575b50816142ca575b81614248575090565b602091505f90604051838101906301ffc9a760e01b82526001600160e01b03197fd3ce6f1b000000000000000000000000000000000000000000000000000000001660248201526024815261429e604482612a1f565b5191617530fa5f513d826142be575b50816142b7575090565b9050151590565b6020111591505f6142ad565b905060205f604051828101906301ffc9a760e01b82526001600160e01b03196024820152602481526142fd604482612a1f565b519084617530fa5f513d82614326575b508161431c575b50159061423f565b905015155f614314565b6020111591505f61430d565b151591505f614238565b6020111592505f614231565b60f81c602f8111806143b0575b1561436457602f190160ff1690565b60608111806143a6575b1561437d576056190160ff1690565b604081118061439c575b15614396576036190160ff1690565b5060ff90565b5060478110614387565b506067811061436e565b50603a811061435556fe60803461013457601f61040038819003918201601f19168301916001600160401b03831184841017610138578084926040948552833981010312610134576100468161014c565b906001600160a01b039061005c9060200161014c565b16908115610121575f80546001600160a01b031981168417825560405193916001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a3803b1561010157600180546001600160a01b0319166001600160a01b039290921691821790557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a261029f90816101618239f35b63211eb15960e21b5f9081526001600160a01b0391909116600452602490fd5b631e4fbdf760e01b5f525f60045260245ffd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b03821682036101345756fe60806040526004361015610011575f80fd5b5f3560e01c80633659cfe6146101af5780635c60da1b14610189578063715018a6146101255780638da5cb5b146101005763f2fde38b14610050575f80fd5b346100fc5760203660031901126100fc576004356001600160a01b0381168091036100fc5761007d610253565b80156100d0576001600160a01b035f548273ffffffffffffffffffffffffffffffffffffffff198216175f55167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3005b7f1e4fbdf7000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b5f80fd5b346100fc575f3660031901126100fc5760206001600160a01b035f5416604051908152f35b346100fc575f3660031901126100fc5761013d610253565b5f6001600160a01b03815473ffffffffffffffffffffffffffffffffffffffff1981168355167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b346100fc575f3660031901126100fc5760206001600160a01b0360015416604051908152f35b346100fc5760203660031901126100fc576004356001600160a01b038116908181036100fc576101dd610253565b3b15610228578073ffffffffffffffffffffffffffffffffffffffff1960015416176001557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a2005b7f847ac564000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b6001600160a01b035f5416330361026657565b7f118cdaa7000000000000000000000000000000000000000000000000000000005f523360045260245ffdfea164736f6c634300081c000a60a0806040526104a880380380916100178285610292565b833981016040828203126101eb5761002e826102c9565b602083015190926001600160401b0382116101eb57019080601f830112156101eb57815161005b816102dd565b926100696040519485610292565b8184526020840192602083830101116101eb57815f926020809301855e84010152823b15610274577fa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d5080546001600160a01b0319166001600160a01b038516908117909155604051635c60da1b60e01b8152909190602081600481865afa9081156101f7575f9161023a575b50803b1561021a5750817f1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e5f80a282511561020257602060049260405193848092635c60da1b60e01b82525afa9182156101f7575f926101ae575b505f809161018a945190845af43d156101a6573d9161016e836102dd565b9261017c6040519485610292565b83523d5f602085013e6102f8565b505b608052604051610151908161035782396080518160460152f35b6060916102f8565b9291506020833d6020116101ef575b816101ca60209383610292565b810103126101eb575f80916101e161018a956102c9565b9394509150610150565b5f80fd5b3d91506101bd565b6040513d5f823e3d90fd5b505050341561018c5763b398979f60e01b5f5260045ffd5b634c9c8ce360e01b5f9081526001600160a01b0391909116600452602490fd5b90506020813d60201161026c575b8161025560209383610292565b810103126101eb57610266906102c9565b5f6100f5565b3d9150610248565b631933b43b60e21b5f9081526001600160a01b038416600452602490fd5b601f909101601f19168101906001600160401b038211908210176102b557604052565b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b03821682036101eb57565b6001600160401b0381116102b557601f01601f191660200190565b9061031c575080511561030d57602081519101fd5b63d6bda27560e01b5f5260045ffd5b8151158061034d575b61032d575090565b639996b31560e01b5f9081526001600160a01b0391909116600452602490fd5b50803b1561032556fe60806040527f5c60da1b000000000000000000000000000000000000000000000000000000006080526020608060048173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa80156100e9575f9015610127575060203d6020116100e2575b601f19601f820116608001906080821067ffffffffffffffff8311176100b5576100b0916040526080016100f4565b610127565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b503d610081565b6040513d5f823e3d90fd5b602090607f1901126101235760805173ffffffffffffffffffffffffffffffffffffffff811681036101235790565b5f80fd5b5f8091368280378136915af43d5f803e15610140573d5ff35b3d5ffdfea164736f6c634300081c000af3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a00f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00a164736f6c634300081c000af0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
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

// SetIBCERC20Metadata is a paid mutator transaction binding the contract method 0x62bfa172.
//
// Solidity: function setIBCERC20Metadata(string denom, string name_, string symbol_, uint8 decimals_) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) SetIBCERC20Metadata(opts *bind.TransactOpts, denom string, name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "setIBCERC20Metadata", denom, name_, symbol_, decimals_)
}

// SetIBCERC20Metadata is a paid mutator transaction binding the contract method 0x62bfa172.
//
// Solidity: function setIBCERC20Metadata(string denom, string name_, string symbol_, uint8 decimals_) returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) SetIBCERC20Metadata(denom string, name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SetIBCERC20Metadata(&_ContractICS20Transfer.TransactOpts, denom, name_, symbol_, decimals_)
}

// SetIBCERC20Metadata is a paid mutator transaction binding the contract method 0x62bfa172.
//
// Solidity: function setIBCERC20Metadata(string denom, string name_, string symbol_, uint8 decimals_) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) SetIBCERC20Metadata(denom string, name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.SetIBCERC20Metadata(&_ContractICS20Transfer.TransactOpts, denom, name_, symbol_, decimals_)
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

// ContractICS20TransferIBCSenderAckPacketCallbackErrorIterator is returned from FilterIBCSenderAckPacketCallbackError and is used to iterate over the raw logs and unpacked data for IBCSenderAckPacketCallbackError events raised by the ContractICS20Transfer contract.
type ContractICS20TransferIBCSenderAckPacketCallbackErrorIterator struct {
	Event *ContractICS20TransferIBCSenderAckPacketCallbackError // Event containing the contract specifics and raw log

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
func (it *ContractICS20TransferIBCSenderAckPacketCallbackErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferIBCSenderAckPacketCallbackError)
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
		it.Event = new(ContractICS20TransferIBCSenderAckPacketCallbackError)
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
func (it *ContractICS20TransferIBCSenderAckPacketCallbackErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferIBCSenderAckPacketCallbackErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferIBCSenderAckPacketCallbackError represents a IBCSenderAckPacketCallbackError event raised by the ContractICS20Transfer contract.
type ContractICS20TransferIBCSenderAckPacketCallbackError struct {
	CallbackAddress common.Address
	Reason          []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterIBCSenderAckPacketCallbackError is a free log retrieval operation binding the contract event 0x898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e0.
//
// Solidity: event IBCSenderAckPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterIBCSenderAckPacketCallbackError(opts *bind.FilterOpts, callbackAddress []common.Address) (*ContractICS20TransferIBCSenderAckPacketCallbackErrorIterator, error) {

	var callbackAddressRule []interface{}
	for _, callbackAddressItem := range callbackAddress {
		callbackAddressRule = append(callbackAddressRule, callbackAddressItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "IBCSenderAckPacketCallbackError", callbackAddressRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferIBCSenderAckPacketCallbackErrorIterator{contract: _ContractICS20Transfer.contract, event: "IBCSenderAckPacketCallbackError", logs: logs, sub: sub}, nil
}

// WatchIBCSenderAckPacketCallbackError is a free log subscription operation binding the contract event 0x898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e0.
//
// Solidity: event IBCSenderAckPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchIBCSenderAckPacketCallbackError(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferIBCSenderAckPacketCallbackError, callbackAddress []common.Address) (event.Subscription, error) {

	var callbackAddressRule []interface{}
	for _, callbackAddressItem := range callbackAddress {
		callbackAddressRule = append(callbackAddressRule, callbackAddressItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "IBCSenderAckPacketCallbackError", callbackAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferIBCSenderAckPacketCallbackError)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "IBCSenderAckPacketCallbackError", log); err != nil {
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

// ParseIBCSenderAckPacketCallbackError is a log parse operation binding the contract event 0x898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e0.
//
// Solidity: event IBCSenderAckPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseIBCSenderAckPacketCallbackError(log types.Log) (*ContractICS20TransferIBCSenderAckPacketCallbackError, error) {
	event := new(ContractICS20TransferIBCSenderAckPacketCallbackError)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "IBCSenderAckPacketCallbackError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS20TransferIBCSenderAckPacketCallbackError0Iterator is returned from FilterIBCSenderAckPacketCallbackError0 and is used to iterate over the raw logs and unpacked data for IBCSenderAckPacketCallbackError0 events raised by the ContractICS20Transfer contract.
type ContractICS20TransferIBCSenderAckPacketCallbackError0Iterator struct {
	Event *ContractICS20TransferIBCSenderAckPacketCallbackError0 // Event containing the contract specifics and raw log

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
func (it *ContractICS20TransferIBCSenderAckPacketCallbackError0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferIBCSenderAckPacketCallbackError0)
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
		it.Event = new(ContractICS20TransferIBCSenderAckPacketCallbackError0)
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
func (it *ContractICS20TransferIBCSenderAckPacketCallbackError0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferIBCSenderAckPacketCallbackError0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferIBCSenderAckPacketCallbackError0 represents a IBCSenderAckPacketCallbackError0 event raised by the ContractICS20Transfer contract.
type ContractICS20TransferIBCSenderAckPacketCallbackError0 struct {
	CallbackAddress common.Address
	Reason          []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterIBCSenderAckPacketCallbackError0 is a free log retrieval operation binding the contract event 0x898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e0.
//
// Solidity: event IBCSenderAckPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterIBCSenderAckPacketCallbackError0(opts *bind.FilterOpts, callbackAddress []common.Address) (*ContractICS20TransferIBCSenderAckPacketCallbackError0Iterator, error) {

	var callbackAddressRule []interface{}
	for _, callbackAddressItem := range callbackAddress {
		callbackAddressRule = append(callbackAddressRule, callbackAddressItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "IBCSenderAckPacketCallbackError0", callbackAddressRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferIBCSenderAckPacketCallbackError0Iterator{contract: _ContractICS20Transfer.contract, event: "IBCSenderAckPacketCallbackError0", logs: logs, sub: sub}, nil
}

// WatchIBCSenderAckPacketCallbackError0 is a free log subscription operation binding the contract event 0x898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e0.
//
// Solidity: event IBCSenderAckPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchIBCSenderAckPacketCallbackError0(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferIBCSenderAckPacketCallbackError0, callbackAddress []common.Address) (event.Subscription, error) {

	var callbackAddressRule []interface{}
	for _, callbackAddressItem := range callbackAddress {
		callbackAddressRule = append(callbackAddressRule, callbackAddressItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "IBCSenderAckPacketCallbackError0", callbackAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferIBCSenderAckPacketCallbackError0)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "IBCSenderAckPacketCallbackError0", log); err != nil {
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

// ParseIBCSenderAckPacketCallbackError0 is a log parse operation binding the contract event 0x898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e0.
//
// Solidity: event IBCSenderAckPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseIBCSenderAckPacketCallbackError0(log types.Log) (*ContractICS20TransferIBCSenderAckPacketCallbackError0, error) {
	event := new(ContractICS20TransferIBCSenderAckPacketCallbackError0)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "IBCSenderAckPacketCallbackError0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS20TransferIBCSenderTimeoutPacketCallbackErrorIterator is returned from FilterIBCSenderTimeoutPacketCallbackError and is used to iterate over the raw logs and unpacked data for IBCSenderTimeoutPacketCallbackError events raised by the ContractICS20Transfer contract.
type ContractICS20TransferIBCSenderTimeoutPacketCallbackErrorIterator struct {
	Event *ContractICS20TransferIBCSenderTimeoutPacketCallbackError // Event containing the contract specifics and raw log

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
func (it *ContractICS20TransferIBCSenderTimeoutPacketCallbackErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferIBCSenderTimeoutPacketCallbackError)
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
		it.Event = new(ContractICS20TransferIBCSenderTimeoutPacketCallbackError)
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
func (it *ContractICS20TransferIBCSenderTimeoutPacketCallbackErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferIBCSenderTimeoutPacketCallbackErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferIBCSenderTimeoutPacketCallbackError represents a IBCSenderTimeoutPacketCallbackError event raised by the ContractICS20Transfer contract.
type ContractICS20TransferIBCSenderTimeoutPacketCallbackError struct {
	CallbackAddress common.Address
	Reason          []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterIBCSenderTimeoutPacketCallbackError is a free log retrieval operation binding the contract event 0x02c075b43cbb5fee7c14de7aea6a6ddde6e3646144712fd9c7159705152537e6.
//
// Solidity: event IBCSenderTimeoutPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterIBCSenderTimeoutPacketCallbackError(opts *bind.FilterOpts, callbackAddress []common.Address) (*ContractICS20TransferIBCSenderTimeoutPacketCallbackErrorIterator, error) {

	var callbackAddressRule []interface{}
	for _, callbackAddressItem := range callbackAddress {
		callbackAddressRule = append(callbackAddressRule, callbackAddressItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "IBCSenderTimeoutPacketCallbackError", callbackAddressRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferIBCSenderTimeoutPacketCallbackErrorIterator{contract: _ContractICS20Transfer.contract, event: "IBCSenderTimeoutPacketCallbackError", logs: logs, sub: sub}, nil
}

// WatchIBCSenderTimeoutPacketCallbackError is a free log subscription operation binding the contract event 0x02c075b43cbb5fee7c14de7aea6a6ddde6e3646144712fd9c7159705152537e6.
//
// Solidity: event IBCSenderTimeoutPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchIBCSenderTimeoutPacketCallbackError(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferIBCSenderTimeoutPacketCallbackError, callbackAddress []common.Address) (event.Subscription, error) {

	var callbackAddressRule []interface{}
	for _, callbackAddressItem := range callbackAddress {
		callbackAddressRule = append(callbackAddressRule, callbackAddressItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "IBCSenderTimeoutPacketCallbackError", callbackAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferIBCSenderTimeoutPacketCallbackError)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "IBCSenderTimeoutPacketCallbackError", log); err != nil {
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

// ParseIBCSenderTimeoutPacketCallbackError is a log parse operation binding the contract event 0x02c075b43cbb5fee7c14de7aea6a6ddde6e3646144712fd9c7159705152537e6.
//
// Solidity: event IBCSenderTimeoutPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseIBCSenderTimeoutPacketCallbackError(log types.Log) (*ContractICS20TransferIBCSenderTimeoutPacketCallbackError, error) {
	event := new(ContractICS20TransferIBCSenderTimeoutPacketCallbackError)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "IBCSenderTimeoutPacketCallbackError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS20TransferIBCSenderTimeoutPacketCallbackError0Iterator is returned from FilterIBCSenderTimeoutPacketCallbackError0 and is used to iterate over the raw logs and unpacked data for IBCSenderTimeoutPacketCallbackError0 events raised by the ContractICS20Transfer contract.
type ContractICS20TransferIBCSenderTimeoutPacketCallbackError0Iterator struct {
	Event *ContractICS20TransferIBCSenderTimeoutPacketCallbackError0 // Event containing the contract specifics and raw log

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
func (it *ContractICS20TransferIBCSenderTimeoutPacketCallbackError0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferIBCSenderTimeoutPacketCallbackError0)
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
		it.Event = new(ContractICS20TransferIBCSenderTimeoutPacketCallbackError0)
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
func (it *ContractICS20TransferIBCSenderTimeoutPacketCallbackError0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferIBCSenderTimeoutPacketCallbackError0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferIBCSenderTimeoutPacketCallbackError0 represents a IBCSenderTimeoutPacketCallbackError0 event raised by the ContractICS20Transfer contract.
type ContractICS20TransferIBCSenderTimeoutPacketCallbackError0 struct {
	CallbackAddress common.Address
	Reason          []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterIBCSenderTimeoutPacketCallbackError0 is a free log retrieval operation binding the contract event 0x02c075b43cbb5fee7c14de7aea6a6ddde6e3646144712fd9c7159705152537e6.
//
// Solidity: event IBCSenderTimeoutPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterIBCSenderTimeoutPacketCallbackError0(opts *bind.FilterOpts, callbackAddress []common.Address) (*ContractICS20TransferIBCSenderTimeoutPacketCallbackError0Iterator, error) {

	var callbackAddressRule []interface{}
	for _, callbackAddressItem := range callbackAddress {
		callbackAddressRule = append(callbackAddressRule, callbackAddressItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "IBCSenderTimeoutPacketCallbackError0", callbackAddressRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferIBCSenderTimeoutPacketCallbackError0Iterator{contract: _ContractICS20Transfer.contract, event: "IBCSenderTimeoutPacketCallbackError0", logs: logs, sub: sub}, nil
}

// WatchIBCSenderTimeoutPacketCallbackError0 is a free log subscription operation binding the contract event 0x02c075b43cbb5fee7c14de7aea6a6ddde6e3646144712fd9c7159705152537e6.
//
// Solidity: event IBCSenderTimeoutPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchIBCSenderTimeoutPacketCallbackError0(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferIBCSenderTimeoutPacketCallbackError0, callbackAddress []common.Address) (event.Subscription, error) {

	var callbackAddressRule []interface{}
	for _, callbackAddressItem := range callbackAddress {
		callbackAddressRule = append(callbackAddressRule, callbackAddressItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "IBCSenderTimeoutPacketCallbackError0", callbackAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferIBCSenderTimeoutPacketCallbackError0)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "IBCSenderTimeoutPacketCallbackError0", log); err != nil {
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

// ParseIBCSenderTimeoutPacketCallbackError0 is a log parse operation binding the contract event 0x02c075b43cbb5fee7c14de7aea6a6ddde6e3646144712fd9c7159705152537e6.
//
// Solidity: event IBCSenderTimeoutPacketCallbackError(address indexed callbackAddress, bytes reason)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseIBCSenderTimeoutPacketCallbackError0(log types.Log) (*ContractICS20TransferIBCSenderTimeoutPacketCallbackError0, error) {
	event := new(ContractICS20TransferIBCSenderTimeoutPacketCallbackError0)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "IBCSenderTimeoutPacketCallbackError0", log); err != nil {
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
