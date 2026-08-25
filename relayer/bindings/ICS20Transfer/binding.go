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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activateEscrow\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"tokens\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createEscrow\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"escrow\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"enableEscrowLaunchGate\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getEscrow\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEscrowBeacon\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIBCERC20Beacon\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPermit2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ibcERC20Contract\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ibcERC20Denom\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ics26\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"ics26Router\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"escrowLogic\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"ibcERC20Logic\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permit2\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializeV2\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isConsumingScheduledOp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isEscrowActive\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onAcknowledgementPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIIBCAppCallbacks.OnAcknowledgementPacketCallback\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destinationClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"acknowledgement\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayer\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onRecvPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIIBCAppCallbacks.OnRecvPacketCallback\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destinationClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"relayer\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onTimeoutPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIIBCAppCallbacks.OnTimeoutPacketCallback\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destinationClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"relayer\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requiresPrecreatedEscrows\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sendTransfer\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS20TransferMsgs.SendTransferMsg\",\"components\":[{\"name\":\"denom\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"memo\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendTransferWithPermit2\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS20TransferMsgs.SendTransferMsg\",\"components\":[{\"name\":\"denom\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"memo\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"permit\",\"type\":\"tuple\",\"internalType\":\"structISignatureTransfer.PermitTransferFrom\",\"components\":[{\"name\":\"permitted\",\"type\":\"tuple\",\"internalType\":\"structISignatureTransfer.TokenPermissions\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendTransferWithSender\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS20TransferMsgs.SendTransferMsg\",\"components\":[{\"name\":\"denom\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"memo\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAuthority\",\"inputs\":[{\"name\":\"newAuthority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setCustomERC20\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setIBCERC20Metadata\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decimals_\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeEscrowTo\",\"inputs\":[{\"name\":\"newEscrowLogic\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeIBCERC20To\",\"inputs\":[{\"name\":\"newIBCERC20Logic\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCERC20ContractCreated\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fullDenomPath\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCSenderAckPacketCallbackError\",\"inputs\":[{\"name\":\"callbackAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCSenderAckPacketCallbackError\",\"inputs\":[{\"name\":\"callbackAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCSenderTimeoutPacketCallbackError\",\"inputs\":[{\"name\":\"callbackAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCSenderTimeoutPacketCallbackError\",\"inputs\":[{\"name\":\"callbackAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS20EscrowActivated\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"escrow\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS20EscrowCreated\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"escrow\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS20EscrowLaunchGateEnabled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessManagedInvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessManagedRequiredDelay\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"AccessManagedUnauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ICS20DenomAlreadyExists\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20DenomNotFound\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20EscrowNotActive\",\"inputs\":[{\"name\":\"clientID\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20EscrowNotFound\",\"inputs\":[{\"name\":\"clientID\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20EscrowNotProvisioned\",\"inputs\":[{\"name\":\"clientID\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20EscrowRateLimitNotSet\",\"inputs\":[{\"name\":\"clientID\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ICS20EscrowTokenListEmpty\",\"inputs\":[{\"name\":\"clientID\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20InvalidAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20InvalidAmount\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ICS20InvalidPort\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20Permit2TokenMismatch\",\"inputs\":[{\"name\":\"permitToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sentToken\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ICS20TokenAlreadyExists\",\"inputs\":[{\"name\":\"denom\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ICS20UnauthorizedPacketSender\",\"inputs\":[{\"name\":\"packetSender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ICS20UnexpectedERC20Balance\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ICS20UnexpectedEncoding\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ICS20UnexpectedVersion\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60a080604052346100c257306080525f5160206157975f395f51905f525460ff8160401c166100b3576002600160401b03196001600160401b03821601610060575b6040516156d090816100c782396080518181816118b701526119670152f35b6001600160401b0319166001600160401b039081175f5160206157975f395f51905f525581527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f80610041565b63f92ee8a960e01b5f5260045ffd5b5f80fdfe60806040526004361015610011575f80fd5b5f5f3560e01c806306ab20bc14612f15578063078c4a79146127af5780631459457a1461241d5780631bbf2e23146123d75780631cf0a67f146123925780631e5150e41461234c57806329b6eca9146121415780632ac3dc38146120a35780633f4ba83a14611fe2578063428e4e1714611bb85780634f1ef2861461191657806352d1902d1461189c57806353816a7c1461148657806357a85bea1461143b5780635c975abb146113f95780635e32b6b6146111c957806362bfa1721461108c5780637a9e5e4b14610ffb578063826cae7a14610fb55780638456cb5914610f1a5780638fb3603714610eb057806390f6d62014610e5657806393dcfee414610bd7578063969631d514610b5b578063a1d28f5714610845578063a50ee2b4146107e7578063aaa2c3431461073b578063ac9650d8146105be578063ad3cb1cc1461055d578063b29c715d146103be578063bf7e214f1461038b578063d413227d14610345578063dab8ed9c146102a25763e163b1af14610190575f80fd5b3461029f57602036600319011261029f576101e26101ac612fba565b6001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9060405191818154916101f4836136db565b80865292600181169081156102755750600114610234575b6102308561021c81870382613082565b60405191829160208352602083019061302e565b0390f35b815260208120939250905b80821061025b5750909150810160200161021c8261023061020c565b91926001816020925483858801015201910190929161023f565b8695506102309693506020925061021c94915060ff191682840152151560051b820101929361020c565b80fd5b503461029f578060031936011261029f576102bd36336137cd565b7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065460ff8160a01c16156102ef575080f35b60ff60a01b1916600160a01b177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f806557f602a536c4b89ff28b2a974a997fe044117170c12aa39af142bdbf06a097bb7fd8180a180f35b503461029f578060031936011261029f5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416604051908152f35b503461029f578060031936011261029f5760206001600160a01b035f5160206156845f395f51905f525416604051908152f35b503461029f57604036600319011261029f576004359067ffffffffffffffff821161029f578160040160e06003198436030112610559576103fd612fd0565b906104066139e8565b61040e613974565b61041836336137cd565b6024840135938415610546579060206001600160a01b0361045261044d61044660646104a8989701866131a7565b36916130c0565b613a3b565b16956104688161046185613501565b8933613c02565b61047183613501565b60405163b4f22eb760e01b81526001600160a01b039091166004820152336024820152604481019190915293849081906064820190565b038187895af192831561053b578493610500575b50906020946104cb9392613d9c565b907f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d67ffffffffffffffff60405191168152f35b9250906020833d602011610533575b8161051c60209383613082565b8101031261052f579151919060206104bc565b5f80fd5b3d915061050f565b6040513d86823e3d90fd5b6024846304f6df8d60e41b815280600452fd5b5080fd5b503461029f578060031936011261029f5750610230604051610580604082613082565b600581527f352e302e30000000000000000000000000000000000000000000000000000000602082015260405191829160208352602083019061302e565b503461029f57602036600319011261029f5760043567ffffffffffffffff8111610559576105f0903690600401613124565b9060206040516106008282613082565b84815281810191601f198101368437610618856137a1565b936106266040519586613082565b858552601f19610635876137a1565b0182885b82811061072b57505050865b868110156106cc576001906106b0898061069c886106688660051b8901896131a7565b8a8d6040959395519483869484860198893784019083820190898252519283915e010185815203601f198101835282613082565b5190305af46106a961439b565b9030614a97565b6106ba82896137b9565b526106c581886137b9565b5001610645565b82888760405191838301848452825180915260408401948060408360051b870101940192955b8287106106ff5785850386f35b90919293828061071b600193603f198a8203018652885161302e565b96019201960195929190926106f2565b606082828a010152018390610639565b503461029f57602036600319011261029f5780610756612fba565b61076036336137cd565b6001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805541690813b156107e3576001600160a01b0360248492836040519586948593631b2ce7f360e11b85521660048401525af180156107d8576107c75750f35b816107d191613082565b61029f5780f35b6040513d84823e3d90fd5b5050fd5b503461029f57602036600319011261029f576004359067ffffffffffffffff821161029f57602061083d61081e36600486016130f6565b6001600160a01b03610832828495946135ba565b54169283151561362a565b604051908152f35b503461029f57604036600319011261029f5760043567ffffffffffffffff8111610559576108779036906004016130f6565b610882929192612fd0565b61088c36336137cd565b6001600160a01b0361089e83866135ba565b5416610b1b576108e76108e1826001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b546136db565b15610922826001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9015610a5f57506109889061093783866135ba565b6001600160a01b03808316166001600160a01b03198254161790556001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9067ffffffffffffffff8111610a4b576109ac816109a684546136db565b84613713565b82601f82116001146109eb57819084956109db9495926109e0575b50508160011b915f199060031b1c19161790565b905580f35b013590505f806109c7565b601f198216948385526020852091855b878110610a33575083600195969710610a1a575b505050811b01905580f35b01355f19600384901b60f8161c191690555f8080610a0f565b909260206001819286860135815501940191016109fb565b602483634e487b7160e01b81526041600452fd5b83906040519182917f778769c40000000000000000000000000000000000000000000000000000000083526020600484015281815491610a9e836136db565b928360248701526001811690815f14610afb5750600114610ac1575b5050500390fd5b9080935052602082205b818310610ae15750508101604401838080610aba565b805460448487010152849350602090920191600101610acb565b925050506044925060ff191682840152151560051b820101838080610aba565b6040517f0c0ef5340000000000000000000000000000000000000000000000000000000081526020600482015280610b57602482018588613215565b0390fd5b503461029f57602036600319011261029f5760043567ffffffffffffffff81116105595760209182610b996001600160a01b039336906004016130f6565b925082604051938492833781017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8008152030190205416604051908152f35b503461029f57604036600319011261029f5760043567ffffffffffffffff811161055957610c099036906004016130f6565b60243567ffffffffffffffff8111610e5257610c29903690600401613124565b9290610c3536336137cd565b6001600160a01b03604051848482376020818681017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8008152030190205416938415610e2f578015610df357855b818110610cd957505050610c968282613582565b600160ff1982541617905581604051928392833781018481520390207f4091a6ce2a3c52563f54d5437475ea54363fc120a36395a84533922c2d5ad9918380a380f35b6001600160a01b03610cf4610cef83858761365a565b613501565b16151580610dd7575b610d17908686610d11610cef86888a61365a565b9261367e565b610d25610cef82848661365a565b906001600160a01b03604051927f0b0aee690000000000000000000000000000000000000000000000000000000084521660048301526020826024818a5afa918215610dcc578892610d97575b50610d916001928787610d89610cef86898b61365a565b92151561367e565b01610c82565b91506020823d8211610dc4575b81610db160209383613082565b8101031261052f57905190610d91610d72565b3d9150610da4565b6040513d8a823e3d90fd5b50610d17610de9610cef83858761365a565b3b15159050610cfd565b6040517fd5c055800000000000000000000000000000000000000000000000000000000081526020600482015280610b57602482018787613215565b6040516358d1c8a160e11b81526020600482015280610b57602482018787613215565b8380fd5b503461029f57602036600319011261029f576004359067ffffffffffffffff821161029f5760206001600160a01b03610ea7610ea2610e9836600488016130f6565b61044636336137cd565b6146d2565b16604051908152f35b503461029f578060031936011261029f575f5160206156845f395f51905f525460a01c60ff1615610f12575060207f8fb36037000000000000000000000000000000000000000000000000000000005b6001600160e01b031960405191168152f35b602090610f00565b503461029f578060031936011261029f57610f3536336137cd565b610f3d6139e8565b600160ff197fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005416177fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a180f35b503461029f578060031936011261029f5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8045416604051908152f35b503461029f57602036600319011261029f57611015612fba565b6001600160a01b035f5160206156845f395f51905f525416330361107a57803b156110465761104390614672565b80f35b7fc2f31e5e0000000000000000000000000000000000000000000000000000000082526001600160a01b0316600452602490fd5b60248262d1953b60e31b815233600452fd5b503461029f57608036600319011261029f578060043567ffffffffffffffff81116111c6576110bf9036906004016130f6565b60243567ffffffffffffffff81116111c3576110df9036906004016130f6565b91909260443567ffffffffffffffff81116111bf576111029036906004016130f6565b94906064359360ff85168095036111bb5761113e9061112136336137cd565b6001600160a01b0361113382876135ba565b54169485151561362a565b823b156111b757869461118e946111a08793604051998a98899788967f37d2c2f4000000000000000000000000000000000000000000000000000000008852606060048901526064880191613215565b85810360031901602487015291613215565b90604483015203925af180156107d8576107c75750f35b8680fd5b8780fd5b8580fd5b50505b50fd5b503461029f576111d836612ffa565b61120e336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314613155565b611216613974565b60608101908261123e61123661122c8585613192565b60808101906131a7565b810190613332565b9261127661125561124f8386613192565b806131a7565b959061126186806131a7565b9060408801986112708a613515565b936143ca565b905061128181614c10565b6112ad575b82807f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d80f35b6001600160a01b031693843b156113f557829160405194859283927f5e32b6b600000000000000000000000000000000000000000000000000000000845260048401602090526112fd83806142ab565b6024860160a0905260c486019061131392613215565b61132060208501856142ab565b8683036023190160448801526113369291613215565b91611340906142dd565b67ffffffffffffffff16606485015261135990836142f2565b83820360231901608485015261136f9190614306565b9060800161137c90612fe6565b6001600160a01b031660a483015203818387620186a0f191826113e0575b50506113da577f02c075b43cbb5fee7c14de7aea6a6ddde6e3646144712fd9c7159705152537e66113cc61021c61439b565b0390a25b5f80828180611286565b506113d0565b816113ea91613082565b6113f557825f61139a565b8280fd5b503461029f578060031936011261029f57602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b503461029f57602036600319011261029f576004359067ffffffffffffffff821161029f57602060ff61147a61147436600487016130f6565b90613582565b54166040519015158152f35b503461029f5760c036600319011261029f576004359067ffffffffffffffff821161029f57816004019160e060031982360301126105595760803660231901126105595760a43567ffffffffffffffff81116113f5576114ea9036906004016130f6565b6114f59291926139e8565b6114fd613974565b6024820135918215611889576115116134eb565b6001600160a01b038061152389613501565b1691161461152f6134eb565b61153888613501565b911561184f5750509061155661044d610446606460249501896131a7565b9360206001600160a01b038061156b8a613501565b16961695604051948580926370a0823160e01b82528960048301525afa928315611844578693611810575b506001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80654166040516040810181811067ffffffffffffffff8211176117fc576040528681526020810192868452823b156117f857888094939261167c6001600160a01b0393604051988997889687957f30f28b7a00000000000000000000000000000000000000000000000000000000875281611638612fd0565b166004880152604435602488015260643560448801526084356064880152511660848601525160a48501523360c485015261010060e4850152610104840191613215565b03925af180156117d8579085916117e3575b5050906024929160206001600160a01b036116a888613501565b16604051958680926370a0823160e01b82528760048301525afa9384156117d85785946117a0575b50906116f361173394836116e68460209661352a565b908211806117975761354b565b6116fc86613501565b60405163b4f22eb760e01b81526001600160a01b039091166004820152336024820152604481019190915292839081906064820190565b038186855af191821561178c578392611757575b50906020936104cb923391613d9c565b91506020823d602011611784575b8161177260209383613082565b8101031261052f579051906020611747565b3d9150611765565b6040513d85823e3d90fd5b5080821461354b565b9350906020843d6020116117d0575b816117bc60209383613082565b8101031261052f57925192906116f36116d0565b3d91506117af565b6040513d87823e3d90fd5b816117ed91613082565b610e5257835f61168e565b8880fd5b602489634e487b7160e01b81526041600452fd5b9092506020813d60201161183c575b8161182c60209383613082565b8101031261052f5751915f611596565b3d915061181f565b6040513d88823e3d90fd5b7fe36164780000000000000000000000000000000000000000000000000000000087526001600160a01b0390811660045216602452604485fd5b6024856304f6df8d60e41b815280600452fd5b503461029f578060031936011261029f576001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036119075760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b8063703e46dd60e11b60049252fd5b50604036600319011261029f5761192b612fba565b9060243567ffffffffffffffff811161055957366023820112156105595761195d9036906024816004013591016130c0565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016803014908115611b83575b50611b74576119a136336137cd565b6001600160a01b03831690604051937f52d1902d000000000000000000000000000000000000000000000000000000008552602085600481865afa80958596611b3c575b506119fd5760248484634c9c8ce360e01b8252600452fd5b9091847f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8103611b115750813b15611aff57806001600160a01b03197f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5416177f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b8480a28151839015611acc5780836020611ac895519101845af4611ac261439b565b91614a97565b5080f35b50505034611ad75780f35b807fb398979f0000000000000000000000000000000000000000000000000000000060049252fd5b634c9c8ce360e01b8452600452602483fd5b7faa1d49a4000000000000000000000000000000000000000000000000000000008552600452602484fd5b9095506020813d602011611b6c575b81611b5860209383613082565b81010312611b685751945f6119e5565b8480fd5b3d9150611b4b565b60048263703e46dd60e11b8152fd5b90506001600160a01b037f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc541614155f611992565b503461029f57602036600319011261029f5760043567ffffffffffffffff8111610559578060040160c060031983360301126113f557611c24336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314613155565b611c2c613974565b8260648301611c4161123661122c8386613192565b6084850191611c5361044684876131a7565b60208151910120611c626134b0565b805160209091012014611e3057611c9a611c7f61124f8388613192565b9390611c8b88806131a7565b9060448b019661127088613515565b9050611ca581614c10565b611cd8575b505050505050505b807f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d80f35b6001600160a01b031695863b15611b685784936040519687948594638dfcd9ad60e01b86528760048701526024860160409052611d1583806142ab565b6044880160c09052610104880190611d2c92613215565b611d3960248701856142ab565b8883036043190160648a0152611d4f9291613215565b91611d59906142dd565b67ffffffffffffffff166084870152611d7290836142f2565b8582036043190160a4870152611d889190614306565b91611d92916142ab565b8483036043190160c4860152611da89291613215565b9060a401611db590612fe6565b6001600160a01b031660e483015203818387620186a0f19182611e1b575b5050611e15577f898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e0611e0561021c61439b565b0390a25b5f808281808080611caa565b50611e09565b81611e2591613082565b6113f557825f611dd3565b6020611e5692611e4087806131a7565b949060448a0195611e5087613515565b9161426f565b500151611e6e611e678251836148c9565b92906133ff565b611e7781614c10565b611e88575b50505050505050611cb2565b6001600160a01b031695863b15611b685784936040519687948594638dfcd9ad60e01b865260048601600190526024860160409052611ec783806142ab565b6044880160c09052610104880190611ede92613215565b611eeb60248701856142ab565b8883036043190160648a0152611f019291613215565b91611f0b906142dd565b67ffffffffffffffff166084870152611f2490836142f2565b8582036043190160a4870152611f3a9190614306565b91611f44916142ab565b8483036043190160c4860152611f5a9291613215565b9060a401611f6790612fe6565b6001600160a01b031660e483015203818387620186a0f19182611fcd575b5050611fc7577f898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e0611fb761021c61439b565b0390a25b5f808281808080611e7c565b50611fbb565b81611fd791613082565b6113f557825f611f85565b503461029f578060031936011261029f57611ffd36336137cd565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff81161561207b5760ff19167fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a180f35b6004827f8dfc202b000000000000000000000000000000000000000000000000000000008152fd5b503461029f57602036600319011261029f576004359067ffffffffffffffff821161029f57816004019160e06003198236030112610559576120e36139e8565b6120eb613974565b602481013590811561212e579060206001600160a01b0361211861044d61044660646117339701896131a7565b16916116f38161212788613501565b8533613c02565b6024836304f6df8d60e41b815280600452fd5b503461029f57602036600319011261029f5761215b612fba565b5f5160206156a45f395f51905f525467ffffffffffffffff8116906001820361233d5760401c60ff16908115612331575b50612322576121c2600267ffffffffffffffff195f5160206156a45f395f51905f525416175f5160206156a45f395f51905f5255565b6801000000000000000068ff0000000000000000195f5160206156a45f395f51905f525416175f5160206156a45f395f51905f5255602460206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416604051928380927f24d7806c0000000000000000000000000000000000000000000000000000000082523360048301525afa90811561178c5783916122e3575b509061227661228b923390613155565b61227e614964565b612286614964565b614672565b68ff0000000000000000195f5160206156a45f395f51905f5254165f5160206156a45f395f51905f52557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a180f35b90506020813d60201161231a575b816122fe60209383613082565b810103126113f557519081151582036113f55790612276612266565b3d91506122f1565b60048263f92ee8a960e01b8152fd5b6002915010155f61218c565b60048463f92ee8a960e01b8152fd5b503461029f578060031936011261029f5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8055416604051908152f35b503461029f578060031936011261029f57602060ff7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065460a01c166040519015158152f35b503461029f578060031936011261029f5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065416604051908152f35b503461029f5760a036600319011261029f57612437612fba565b61243f612fd0565b6044356001600160a01b0381168103610e5257606435926001600160a01b038416809403611b68576084356001600160a01b03811681036111bf575f5160206156a45f395f51905f525467ffffffffffffffff811690816127a05760401c60ff16908115612794575b5061278557906125386001600160a01b03926124eb600267ffffffffffffffff195f5160206156a45f395f51905f525416175f5160206156a45f395f51905f5255565b6801000000000000000068ff0000000000000000195f5160206156a45f395f51905f525416175f5160206156a45f395f51905f5255612528614964565b612530614964565b612276614964565b166001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80355604051916104009283810181811067ffffffffffffffff821117612771576125ca8291614ddc95878785396001600160a01b0316815230602082015260400190565b039086f080156117d8576001600160a01b03166001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8045416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80455604051928084019284841067ffffffffffffffff85111761277157918493916126689385396001600160a01b0316815230602082015260400190565b039083f080156107d8576001600160a01b03166001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8055416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805556001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065568ff0000000000000000195f5160206156a45f395f51905f5254165f5160206156a45f395f51905f52557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a180f35b602487634e487b7160e01b81526041600452fd5b60048663f92ee8a960e01b8152fd5b6002915010155f6124a8565b60048863f92ee8a960e01b8152fd5b503461029f576127be36612ffa565b906127f5336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314613155565b6127fd613974565b6128056139e8565b606082016128236104466128198386613192565b60408101906131a7565b602081519101206128326131da565b60208151910120146128426131da565b61284f6128198487613192565b909215612edf575050506128a161286c61044661124f8487613192565b6020815191012061287b61325d565b602081519101201461288b61325d565b9061289961124f8588613192565b929091613298565b6128bb6104466128b18386613192565b60608101906131a7565b602081519101206128ca6132dc565b60208151910120146128da6132dc565b6128e76128b18487613192565b909215612ea95750505061293b61290e6104466129048487613192565b60208101906131a7565b6020815191012061291d61325d565b602081519101201461292d61325d565b906128996129048588613192565b61294b61123661122c8386613192565b92606084018051156105465760408501926129826001600160a01b038551612977611e678251836148c9565b1694518515156133ff565b602083019261299761044d61044686846131a7565b9651936129c66129aa61124f8585613192565b6129c16129b786806131a7565b93909236916130c0565b613b99565b92855184511080612e99575b15612b2c5750509051835195969495879590949092509060209085811890861102851880612a08612a038289613762565b61376f565b9603920101602085015e612a1d8351846148c9565b939015848115612b1a575b50612af1575b506001600160a01b03905b16905191813b15610e5257836064926001600160a01b039460405197889687957f0779afe6000000000000000000000000000000000000000000000000000000008752166004860152602485015260448401525af180156107d857612adc575b61023082612aa56134b0565b907f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d60405191829160208352602083019061302e565b612ae7828092613082565b61029f575f612a99565b92506001600160a01b0390612b1482612b0986613445565b541694851515613483565b90612a2e565b6001600160a01b03915016155f612a28565b612b969350612b618360209794936129c16129b7612b59612b508c98978998613192565b878101906131a7565b9390946131a7565b926040519684889551918291018487015e8401908282018a8152815193849201905e010186815203601f198101845283613082565b6001600160a01b038516946001600160a01b03612bb284613445565b5416928315612c3b575b5084956001600160a01b03849695969416908351823b156111b7576040516340c10f1960e01b81526001600160a01b0392909216600483015260248201529085908290604490829084905af19081156117d8578591612c26575b50506001600160a01b0390612a39565b81612c3091613082565b610e5257835f612c16565b92506001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80454166040517f4571e3a600000000000000000000000000000000000000000000000000000000602082015230602482015287604482015260606064820152612cc381612cb5608482018861302e565b03601f198101835282613082565b604051916104a88084019084821067ffffffffffffffff831117612e855791849391612cf3936151dc8639613be2565b039086f080156117d8576001600160a01b031695612d1084613445565b6001600160a01b0388166001600160a01b0319825416179055612d63876001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9684519767ffffffffffffffff8911612e7157612d8a89612d8483546136db565b83613713565b602098601f8111600114612e0e5780612dbb918a9b8b9a9b91612e03575b508160011b915f199060031b1c19161790565b90555b7f6031fab685dd6d86e4dbac9a69eae347145f332c95b3a0d728d3730fc5233d62612df8829660405191829160208352602083019061302e565b0390a2959493612bbc565b90508801515f612da8565b81895289892099601f1982168a5b818110612e5957509a82916001938c9d9c9b9c10612e41575b5050811b019055612dbe565b8901515f1960f88460031b161c191690555f80612e35565b898301518d556001909c019b60209283019201612e1c565b602488634e487b7160e01b81526041600452fd5b60248a634e487b7160e01b81526041600452fd5b50612ea48487614877565b6129d2565b610b57906040519384937fd1ca953a00000000000000000000000000000000000000000000000000000000855260048501613235565b610b57906040519384937f094af3b800000000000000000000000000000000000000000000000000000000855260048501613235565b503461052f57602036600319011261052f57612f2f612fba565b612f3936336137cd565b6001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f804541690813b1561052f576001600160a01b0360245f92836040519586948593631b2ce7f360e11b85521660048401525af18015612faf57612fa1575080f35b612fad91505f90613082565b005b6040513d5f823e3d90fd5b600435906001600160a01b038216820361052f57565b602435906001600160a01b038216820361052f57565b35906001600160a01b038216820361052f57565b602060031982011261052f576004359067ffffffffffffffff821161052f5760a090829003600319011261052f5760040190565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b60a0810190811067ffffffffffffffff82111761306e57604052565b634e487b7160e01b5f52604160045260245ffd5b90601f8019910116810190811067ffffffffffffffff82111761306e57604052565b67ffffffffffffffff811161306e57601f01601f191660200190565b9291926130cc826130a4565b916130da6040519384613082565b82948184528183011161052f578281602093845f960137010152565b9181601f8401121561052f5782359167ffffffffffffffff831161052f576020838186019501011161052f57565b9181601f8401121561052f5782359167ffffffffffffffff831161052f576020808501948460051b01011161052f57565b1561315d5750565b6001600160a01b03907f2ecb3242000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b903590609e198136030182121561052f570190565b903590601e198136030182121561052f570180359067ffffffffffffffff821161052f5760200191813603831361052f57565b604051906131e9604083613082565b600782527f69637332302d31000000000000000000000000000000000000000000000000006020830152565b908060209392818452848401375f828201840152601f01601f1916010190565b9161324c61325a949260408552604085019061302e565b926020818503910152613215565b90565b6040519061326c604083613082565b600882527f7472616e736665720000000000000000000000000000000000000000000000006020830152565b92909192156132a657505050565b610b57906040519384937f5d3a3cdd00000000000000000000000000000000000000000000000000000000855260048501613235565b604051906132eb604083613082565b601a82527f6170706c69636174696f6e2f782d736f6c69646974792d6162690000000000006020830152565b9080601f8301121561052f5781602061325a933591016130c0565b60208183031261052f5780359067ffffffffffffffff821161052f570160a08183031261052f576040519161336683613052565b813567ffffffffffffffff811161052f5781613383918401613317565b8352602082013567ffffffffffffffff811161052f57816133a5918401613317565b6020840152604082013567ffffffffffffffff811161052f57816133ca918401613317565b604084015260608201356060840152608082013567ffffffffffffffff811161052f576133f79201613317565b608082015290565b156134075750565b610b57906040519182917f3fed5d8700000000000000000000000000000000000000000000000000000000835260206004840152602483019061302e565b60208091604051928184925191829101835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80181520301902090565b1561348b5750565b610b579060405191829163e1275e2f60e01b835260206004840152602483019061302e565b604051906134bf604083613082565b601182527f7b22726573756c74223a2241513d3d227d0000000000000000000000000000006020830152565b6024356001600160a01b038116810361052f5790565b356001600160a01b038116810361052f5790565b3567ffffffffffffffff8116810361052f5790565b9190820180921161353757565b634e487b7160e01b5f52601160045260245ffd5b15613554575050565b7f2fb30cfc000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b60209082604051938492833781017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80781520301902090565b60209082604051938492833781017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80181520301902090565b60209082604051938492833781017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80881520301902090565b91909115613636575050565b610b5760405192839263e1275e2f60e01b8452602060048501526024840191613215565b919081101561366a5760051b0190565b634e487b7160e01b5f52603260045260245ffd5b9290921561368b57505050565b6001600160a01b036136d06040519485947f28203fe7000000000000000000000000000000000000000000000000000000008652604060048701526044860191613215565b911660248301520390fd5b90600182811c92168015613709575b60208310146136f557565b634e487b7160e01b5f52602260045260245ffd5b91607f16916136ea565b601f821161372057505050565b5f5260205f20906020601f840160051c83019310613758575b601f0160051c01905b81811061374d575050565b5f8155600101613742565b9091508190613739565b9190820391821161353757565b90613779826130a4565b6137866040519182613082565b8281528092613797601f19916130a4565b0190602036910137565b67ffffffffffffffff811161306e5760051b60200190565b805182101561366a5760209160051b010190565b5f5160206156845f395f51905f5254916001600160a01b038316928160041161052f575f5f9060405f8151966001600160a01b0360208901917fb700961300000000000000000000000000000000000000000000000000000000835216978860248201523060448201526001600160e01b0319833516606482015260648152613857608482613082565b828052826020525190895afa613961575b15613875575b5050505050565b63ffffffff161561394f5760ff60a01b1916600160a01b175f5160206156845f395f51905f5255823b1561052f576020925f92836040518096819582947f94c7d7ee00000000000000000000000000000000000000000000000000000000845260048401526040602484015260448301908082528085848401378181018301859052601f01601f1916010103925af18015612faf5761393f575b5060ff60a01b195f5160206156845f395f51905f5254165f5160206156845f395f51905f52555f8080808061386e565b5f61394991613082565b5f61390f565b8262d1953b60e31b5f5260045260245ffd5b50505f516020518060201c150290613868565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005c6139c05760017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d565b7f3ee5aeb5000000000000000000000000000000000000000000000000000000005f5260045ffd5b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005416613a1357565b7fd93c0665000000000000000000000000000000000000000000000000000000005f5260045ffd5b604051906001600160a01b03815192602081818501958087835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80081520301902054169182155f14613ae65750905060ff7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065460a01c16613ac15761325a906146d2565b610b57906040519182916358d1c8a160e11b835260206004840152602483019061302e565b60ff7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065460a01c16613b1757505090565b60ff906020604051809285518091835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f807815203019020541615613b5b575090565b610b57906040519182917f0d2e2ae300000000000000000000000000000000000000000000000000000000835260206004840152602483019061302e565b600180916020809561325a958160405198858a9651918291018688015e850191602f60f81b8584015260218301370101602f60f81b838201520301601e19810184520182613082565b6040906001600160a01b0361325a9493168152816020820152019061302e565b916001600160a01b031691604051906370a0823160e01b82526001600160a01b03831691826004820152602081602481885afa938415612faf5786915f95613d65575b506040517f23b872dd0000000000000000000000000000000000000000000000000000000060208281019182526001600160a01b03958616602484015294909216604482015260648101929092525f91613ca28160848101612cb5565b519082875af115612faf575f513d613d5c5750823b155b613d30576020906024604051809581936370a0823160e01b835260048301525afa918215612faf575f92613cf8575b506116e6613cf6938261352a565b565b9291506020833d602011613d28575b81613d1460209383613082565b8101031261052f57915190916116e6613ce8565b3d9150613d07565b827f5274afe7000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b60011415613cb9565b915093506020813d602011613d94575b81613d8260209383613082565b8101031261052f57519285905f613c45565b3d9150613d75565b92919092613dac6101ac82613501565b935f9060405195865f825492613dc1846136db565b808452936001811690811561424d5750600114614209575b50613de692500387613082565b858051155f14614124575050613f689450613e11613e0b613e0684613501565b6149a8565b936149a8565b613e1e60408401846131a7565b94613e6d613e52613e3260c08801886131a7565b94909860405194613e4286613052565b85526020850196875236916130c0565b966040830197885260608301936020880135855236916130c0565b9160808201928352613f116001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416946060880198613eb38a8a6131a7565b969094613f5e613ec560a08d01613515565b97613f50613ed161325d565b93613eda61325d565b95613f37613ee66131da565b98613f24613ef26132dc565b9b6040519d8e986020808b01525160a060408b015260e08a019061302e565b9051888203603f190160608a015261302e565b9051868203603f1901608088015261302e565b915160a085015251838203603f190160c085015261302e565b03601f198101875286613082565b6040519d8e613052565b8d5260208d015260408c015260608b015260808a0152604051926060840184811067ffffffffffffffff821117612771579361408d60209694613fbe614014958a9567ffffffffffffffff9960405236916130c0565b83528688840191168152604083019c8d526040519c8d97889687957f4d6e7ce30000000000000000000000000000000000000000000000000000000087528b60048801525160606024880152608487019061302e565b92511660448501525190602319848203016064850152608061407c61406a614058614048865160a0875260a087019061302e565b8d8701518682038f88015261302e565b6040860151858203604087015261302e565b6060850151848203606086015261302e565b92015190608081840391015261302e565b03925af19485156141175781956140cc575b50506140b4916140ae916131a7565b906135f2565b67ffffffffffffffff83165f5260205260405f205590565b909194506020813d60201161410f575b816140e960209383613082565b8101031261055957519067ffffffffffffffff8216820361029f575092816140ae61409f565b3d91506140dc565b50604051903d90823e3d90fd5b61414561413297959761325d565b61413f60608701876131a7565b91613b99565b815181511091826141f8575b5050614165575b50613e11613f68956149a8565b6001600160a01b0316946001600160a01b0361418084613501565b16863b1561052f575f9660448892604051998a9384927f9dc29fac0000000000000000000000000000000000000000000000000000000084526004840152602089013560248401525af1958615612faf57613f68966141e1575b5094614158565b6141ee9192505f90613082565b5f90613e116141da565b6142029250614877565b5f80614151565b90505f9291925260205f20905f915b818310614231575050906020613de6928201015f613dd9565b6020919350806001915483858d01015201910190918892614218565b905060209250613de694915060ff191682840152151560051b8201015f613dd9565b92919061429c67ffffffffffffffff9161428981876135f2565b8385165f5260205260405f2054956135f2565b91165f526020525f6040812055565b9035601e198236030181121561052f57016020813591019167ffffffffffffffff821161052f57813603831361052f57565b359067ffffffffffffffff8216820361052f57565b9035609e198236030181121561052f570190565b61325a9161438d61438261436761434c61433161432387806142ab565b60a0885260a0880191613215565b61433e60208801886142ab565b908783036020890152613215565b61435960408701876142ab565b908683036040880152613215565b61437460608601866142ab565b908583036060870152613215565b9260808101906142ab565b916080818503910152613215565b3d156143c5573d906143ac826130a4565b916143ba6040519384613082565b82523d5f602084013e565b606090565b95949392939190915f936001600160a01b03604051878582376020818981017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80081520301902054169384156146365786846129c16144389b61443f9461044660208901519e8f8051906148c9565b9f906133ff565b82519081518151109182614625575b5050156145a7576001600160a01b036144678351613445565b5416956144778351881515613483565b6060830151873b1561052f576040516340c10f1960e01b81526001600160a01b038716600482015260248101919091525f81604481838c5af18015612faf57614582575b509067ffffffffffffffff6144eb606094935b6144d881886135f2565b8385165f5260205260405f2054966135f2565b91165f526020528460405f20550151823b15610e525790608484928360405195869485937fc920a6350000000000000000000000000000000000000000000000000000000085526001600160a01b038b1660048601526001600160a01b038d166024860152604485015260648401525af180156107d85761456d575b50509190565b614578828092613082565b61029f5780614567565b606093929196505f61459391613082565b5f959192509067ffffffffffffffff6144bb565b6001600160a01b036145b98351613445565b54169586156145d9575b9067ffffffffffffffff6144eb606094936144ce565b95509060609167ffffffffffffffff6144eb8351986146036145fc8b518c6148c9565b9b906133ff565b61461a8a6001600160a01b03875191161515613483565b9293945050506145c3565b61462f9250614877565b5f8061444e565b6040517f5778f3780000000000000000000000000000000000000000000000000000000081526020600482015280610b57602482018a88613215565b60206001600160a01b037f2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad9216806001600160a01b03195f5160206156845f395f51905f525416175f5160206156845f395f51905f5255604051908152a1565b6040516001600160a01b03825191602081818601948086835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f800815203019020541691821561472257505090565b7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805545f5160206156845f395f51905f52546040517f485cc9550000000000000000000000000000000000000000000000000000000060208201523060248201526001600160a01b039182166044808301919091528152939450919291166147aa606483613082565b604051916104a8908184019284841067ffffffffffffffff85111761306e5784936147d9936151dc8639613be2565b03905ff08015612faf576001600160a01b031691829160405160208183518086835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8008152030190206001600160a01b0384166001600160a01b0319825416179055604051918291518091835e81015f81520390207fd71bcf9347e4c440706076986a93115883dc9a480f43ec5d00610402a437be655f80a390565b80519082518092106148c25780519182808210911802808318921415820282189060206148a7612a038486613762565b92808285019503920101835e51902090602081519101201490565b5050505f90565b80518211801561495d575b614921576001821180614929575b158015908160011b91820460021417156135375760280180602811613537578203614921576001600160a01b0392915f61491b92614b23565b90921690565b50505f905f90565b5061060f60f31b7fffff000000000000000000000000000000000000000000000000000000000000602083015116146148e2565b505f6148d4565b60ff5f5160206156a45f395f51905f525460401c161561498057565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f5260045ffd5b6001600160a01b0316806149bc602a6130a4565b916149ca6040519384613082565b602a83526149d8602a6130a4565b6020840190601f190136823783511561366a576030905382516001101561366a576078602184015360295b60018111614a445750614a14575090565b7fe22e27eb000000000000000000000000000000000000000000000000000000005f52600452601460245260445ffd5b90600f8116601081101561366a57845183101561366a577f3031323334353637383961626364656600000000000000000000000000000000901a8483016020015360041c908015613537575f1901614a03565b90614ad45750805115614aac57602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b81511580614b1a575b614ae5575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15614add565b9290926001840180851161353757831180614bda575b158015908160011b918204600214171561353757614b5c905f949293949561352a565b915b818310614b6e5750505060019190565b9092919360ff614ba57fff000000000000000000000000000000000000000000000000000000000000006020888601015116614d69565b16600f8111614bcf578160041b918083046010149015171561353757600191019401919290614b5e565b505f94508493505050565b5061060f60f31b7fffff000000000000000000000000000000000000000000000000000000000000602086840101511614614b39565b60205f604051828101906301ffc9a760e01b82526301ffc9a760e01b602482015260248152614c40604482613082565b519084617530fa903d5f519083614d5d575b5082614d53575b5081614ceb575b81614c69575090565b602091505f90604051838101906301ffc9a760e01b82526001600160e01b03197fd3ce6f1b0000000000000000000000000000000000000000000000000000000016602482015260248152614cbf604482613082565b5191617530fa5f513d82614cdf575b5081614cd8575090565b9050151590565b6020111591505f614cce565b905060205f604051828101906301ffc9a760e01b82526001600160e01b0319602482015260248152614d1e604482613082565b519084617530fa5f513d82614d47575b5081614d3d575b501590614c60565b905015155f614d35565b6020111591505f614d2e565b151591505f614c59565b6020111592505f614c52565b60f81c602f811180614dd1575b15614d8557602f190160ff1690565b6060811180614dc7575b15614d9e576056190160ff1690565b6040811180614dbd575b15614db7576036190160ff1690565b5060ff90565b5060478110614da8565b5060678110614d8f565b50603a8110614d7656fe60803461013457601f61040038819003918201601f19168301916001600160401b03831184841017610138578084926040948552833981010312610134576100468161014c565b906001600160a01b039061005c9060200161014c565b16908115610121575f80546001600160a01b031981168417825560405193916001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a3803b1561010157600180546001600160a01b0319166001600160a01b039290921691821790557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a261029f90816101618239f35b63211eb15960e21b5f9081526001600160a01b0391909116600452602490fd5b631e4fbdf760e01b5f525f60045260245ffd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b03821682036101345756fe60806040526004361015610011575f80fd5b5f3560e01c80633659cfe6146101af5780635c60da1b14610189578063715018a6146101255780638da5cb5b146101005763f2fde38b14610050575f80fd5b346100fc5760203660031901126100fc576004356001600160a01b0381168091036100fc5761007d610253565b80156100d0576001600160a01b035f548273ffffffffffffffffffffffffffffffffffffffff198216175f55167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3005b7f1e4fbdf7000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b5f80fd5b346100fc575f3660031901126100fc5760206001600160a01b035f5416604051908152f35b346100fc575f3660031901126100fc5761013d610253565b5f6001600160a01b03815473ffffffffffffffffffffffffffffffffffffffff1981168355167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b346100fc575f3660031901126100fc5760206001600160a01b0360015416604051908152f35b346100fc5760203660031901126100fc576004356001600160a01b038116908181036100fc576101dd610253565b3b15610228578073ffffffffffffffffffffffffffffffffffffffff1960015416176001557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a2005b7f847ac564000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b6001600160a01b035f5416330361026657565b7f118cdaa7000000000000000000000000000000000000000000000000000000005f523360045260245ffdfea164736f6c634300081c000a60a0806040526104a880380380916100178285610292565b833981016040828203126101eb5761002e826102c9565b602083015190926001600160401b0382116101eb57019080601f830112156101eb57815161005b816102dd565b926100696040519485610292565b8184526020840192602083830101116101eb57815f926020809301855e84010152823b15610274577fa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d5080546001600160a01b0319166001600160a01b038516908117909155604051635c60da1b60e01b8152909190602081600481865afa9081156101f7575f9161023a575b50803b1561021a5750817f1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e5f80a282511561020257602060049260405193848092635c60da1b60e01b82525afa9182156101f7575f926101ae575b505f809161018a945190845af43d156101a6573d9161016e836102dd565b9261017c6040519485610292565b83523d5f602085013e6102f8565b505b608052604051610151908161035782396080518160460152f35b6060916102f8565b9291506020833d6020116101ef575b816101ca60209383610292565b810103126101eb575f80916101e161018a956102c9565b9394509150610150565b5f80fd5b3d91506101bd565b6040513d5f823e3d90fd5b505050341561018c5763b398979f60e01b5f5260045ffd5b634c9c8ce360e01b5f9081526001600160a01b0391909116600452602490fd5b90506020813d60201161026c575b8161025560209383610292565b810103126101eb57610266906102c9565b5f6100f5565b3d9150610248565b631933b43b60e21b5f9081526001600160a01b038416600452602490fd5b601f909101601f19168101906001600160401b038211908210176102b557604052565b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b03821682036101eb57565b6001600160401b0381116102b557601f01601f191660200190565b9061031c575080511561030d57602081519101fd5b63d6bda27560e01b5f5260045ffd5b8151158061034d575b61032d575090565b639996b31560e01b5f9081526001600160a01b0391909116600452602490fd5b50803b1561032556fe60806040527f5c60da1b000000000000000000000000000000000000000000000000000000006080526020608060048173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa80156100e9575f9015610127575060203d6020116100e2575b601f19601f820116608001906080821067ffffffffffffffff8311176100b5576100b0916040526080016100f4565b610127565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b503d610081565b6040513d5f823e3d90fd5b602090607f1901126101235760805173ffffffffffffffffffffffffffffffffffffffff811681036101235790565b5f80fd5b5f8091368280378136915af43d5f803e15610140573d5ff35b3d5ffdfea164736f6c634300081c000af3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a00f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00a164736f6c634300081c000af0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
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

// IsEscrowActive is a free data retrieval call binding the contract method 0x57a85bea.
//
// Solidity: function isEscrowActive(string clientId) view returns(bool)
func (_ContractICS20Transfer *ContractICS20TransferCaller) IsEscrowActive(opts *bind.CallOpts, clientId string) (bool, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "isEscrowActive", clientId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEscrowActive is a free data retrieval call binding the contract method 0x57a85bea.
//
// Solidity: function isEscrowActive(string clientId) view returns(bool)
func (_ContractICS20Transfer *ContractICS20TransferSession) IsEscrowActive(clientId string) (bool, error) {
	return _ContractICS20Transfer.Contract.IsEscrowActive(&_ContractICS20Transfer.CallOpts, clientId)
}

// IsEscrowActive is a free data retrieval call binding the contract method 0x57a85bea.
//
// Solidity: function isEscrowActive(string clientId) view returns(bool)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) IsEscrowActive(clientId string) (bool, error) {
	return _ContractICS20Transfer.Contract.IsEscrowActive(&_ContractICS20Transfer.CallOpts, clientId)
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

// RequiresPrecreatedEscrows is a free data retrieval call binding the contract method 0x1cf0a67f.
//
// Solidity: function requiresPrecreatedEscrows() view returns(bool)
func (_ContractICS20Transfer *ContractICS20TransferCaller) RequiresPrecreatedEscrows(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ContractICS20Transfer.contract.Call(opts, &out, "requiresPrecreatedEscrows")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// RequiresPrecreatedEscrows is a free data retrieval call binding the contract method 0x1cf0a67f.
//
// Solidity: function requiresPrecreatedEscrows() view returns(bool)
func (_ContractICS20Transfer *ContractICS20TransferSession) RequiresPrecreatedEscrows() (bool, error) {
	return _ContractICS20Transfer.Contract.RequiresPrecreatedEscrows(&_ContractICS20Transfer.CallOpts)
}

// RequiresPrecreatedEscrows is a free data retrieval call binding the contract method 0x1cf0a67f.
//
// Solidity: function requiresPrecreatedEscrows() view returns(bool)
func (_ContractICS20Transfer *ContractICS20TransferCallerSession) RequiresPrecreatedEscrows() (bool, error) {
	return _ContractICS20Transfer.Contract.RequiresPrecreatedEscrows(&_ContractICS20Transfer.CallOpts)
}

// ActivateEscrow is a paid mutator transaction binding the contract method 0x93dcfee4.
//
// Solidity: function activateEscrow(string clientId, address[] tokens) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) ActivateEscrow(opts *bind.TransactOpts, clientId string, tokens []common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "activateEscrow", clientId, tokens)
}

// ActivateEscrow is a paid mutator transaction binding the contract method 0x93dcfee4.
//
// Solidity: function activateEscrow(string clientId, address[] tokens) returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) ActivateEscrow(clientId string, tokens []common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.ActivateEscrow(&_ContractICS20Transfer.TransactOpts, clientId, tokens)
}

// ActivateEscrow is a paid mutator transaction binding the contract method 0x93dcfee4.
//
// Solidity: function activateEscrow(string clientId, address[] tokens) returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) ActivateEscrow(clientId string, tokens []common.Address) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.ActivateEscrow(&_ContractICS20Transfer.TransactOpts, clientId, tokens)
}

// CreateEscrow is a paid mutator transaction binding the contract method 0x90f6d620.
//
// Solidity: function createEscrow(string clientId) returns(address escrow)
func (_ContractICS20Transfer *ContractICS20TransferTransactor) CreateEscrow(opts *bind.TransactOpts, clientId string) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "createEscrow", clientId)
}

// CreateEscrow is a paid mutator transaction binding the contract method 0x90f6d620.
//
// Solidity: function createEscrow(string clientId) returns(address escrow)
func (_ContractICS20Transfer *ContractICS20TransferSession) CreateEscrow(clientId string) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.CreateEscrow(&_ContractICS20Transfer.TransactOpts, clientId)
}

// CreateEscrow is a paid mutator transaction binding the contract method 0x90f6d620.
//
// Solidity: function createEscrow(string clientId) returns(address escrow)
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) CreateEscrow(clientId string) (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.CreateEscrow(&_ContractICS20Transfer.TransactOpts, clientId)
}

// EnableEscrowLaunchGate is a paid mutator transaction binding the contract method 0xdab8ed9c.
//
// Solidity: function enableEscrowLaunchGate() returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactor) EnableEscrowLaunchGate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractICS20Transfer.contract.Transact(opts, "enableEscrowLaunchGate")
}

// EnableEscrowLaunchGate is a paid mutator transaction binding the contract method 0xdab8ed9c.
//
// Solidity: function enableEscrowLaunchGate() returns()
func (_ContractICS20Transfer *ContractICS20TransferSession) EnableEscrowLaunchGate() (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.EnableEscrowLaunchGate(&_ContractICS20Transfer.TransactOpts)
}

// EnableEscrowLaunchGate is a paid mutator transaction binding the contract method 0xdab8ed9c.
//
// Solidity: function enableEscrowLaunchGate() returns()
func (_ContractICS20Transfer *ContractICS20TransferTransactorSession) EnableEscrowLaunchGate() (*types.Transaction, error) {
	return _ContractICS20Transfer.Contract.EnableEscrowLaunchGate(&_ContractICS20Transfer.TransactOpts)
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

// ContractICS20TransferICS20EscrowActivatedIterator is returned from FilterICS20EscrowActivated and is used to iterate over the raw logs and unpacked data for ICS20EscrowActivated events raised by the ContractICS20Transfer contract.
type ContractICS20TransferICS20EscrowActivatedIterator struct {
	Event *ContractICS20TransferICS20EscrowActivated // Event containing the contract specifics and raw log

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
func (it *ContractICS20TransferICS20EscrowActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferICS20EscrowActivated)
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
		it.Event = new(ContractICS20TransferICS20EscrowActivated)
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
func (it *ContractICS20TransferICS20EscrowActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferICS20EscrowActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferICS20EscrowActivated represents a ICS20EscrowActivated event raised by the ContractICS20Transfer contract.
type ContractICS20TransferICS20EscrowActivated struct {
	ClientId common.Hash
	Escrow   common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterICS20EscrowActivated is a free log retrieval operation binding the contract event 0x4091a6ce2a3c52563f54d5437475ea54363fc120a36395a84533922c2d5ad991.
//
// Solidity: event ICS20EscrowActivated(string indexed clientId, address indexed escrow)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterICS20EscrowActivated(opts *bind.FilterOpts, clientId []string, escrow []common.Address) (*ContractICS20TransferICS20EscrowActivatedIterator, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var escrowRule []interface{}
	for _, escrowItem := range escrow {
		escrowRule = append(escrowRule, escrowItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "ICS20EscrowActivated", clientIdRule, escrowRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferICS20EscrowActivatedIterator{contract: _ContractICS20Transfer.contract, event: "ICS20EscrowActivated", logs: logs, sub: sub}, nil
}

// WatchICS20EscrowActivated is a free log subscription operation binding the contract event 0x4091a6ce2a3c52563f54d5437475ea54363fc120a36395a84533922c2d5ad991.
//
// Solidity: event ICS20EscrowActivated(string indexed clientId, address indexed escrow)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchICS20EscrowActivated(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferICS20EscrowActivated, clientId []string, escrow []common.Address) (event.Subscription, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var escrowRule []interface{}
	for _, escrowItem := range escrow {
		escrowRule = append(escrowRule, escrowItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "ICS20EscrowActivated", clientIdRule, escrowRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferICS20EscrowActivated)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "ICS20EscrowActivated", log); err != nil {
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

// ParseICS20EscrowActivated is a log parse operation binding the contract event 0x4091a6ce2a3c52563f54d5437475ea54363fc120a36395a84533922c2d5ad991.
//
// Solidity: event ICS20EscrowActivated(string indexed clientId, address indexed escrow)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseICS20EscrowActivated(log types.Log) (*ContractICS20TransferICS20EscrowActivated, error) {
	event := new(ContractICS20TransferICS20EscrowActivated)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "ICS20EscrowActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS20TransferICS20EscrowCreatedIterator is returned from FilterICS20EscrowCreated and is used to iterate over the raw logs and unpacked data for ICS20EscrowCreated events raised by the ContractICS20Transfer contract.
type ContractICS20TransferICS20EscrowCreatedIterator struct {
	Event *ContractICS20TransferICS20EscrowCreated // Event containing the contract specifics and raw log

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
func (it *ContractICS20TransferICS20EscrowCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferICS20EscrowCreated)
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
		it.Event = new(ContractICS20TransferICS20EscrowCreated)
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
func (it *ContractICS20TransferICS20EscrowCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferICS20EscrowCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferICS20EscrowCreated represents a ICS20EscrowCreated event raised by the ContractICS20Transfer contract.
type ContractICS20TransferICS20EscrowCreated struct {
	ClientId common.Hash
	Escrow   common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterICS20EscrowCreated is a free log retrieval operation binding the contract event 0xd71bcf9347e4c440706076986a93115883dc9a480f43ec5d00610402a437be65.
//
// Solidity: event ICS20EscrowCreated(string indexed clientId, address indexed escrow)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterICS20EscrowCreated(opts *bind.FilterOpts, clientId []string, escrow []common.Address) (*ContractICS20TransferICS20EscrowCreatedIterator, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var escrowRule []interface{}
	for _, escrowItem := range escrow {
		escrowRule = append(escrowRule, escrowItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "ICS20EscrowCreated", clientIdRule, escrowRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferICS20EscrowCreatedIterator{contract: _ContractICS20Transfer.contract, event: "ICS20EscrowCreated", logs: logs, sub: sub}, nil
}

// WatchICS20EscrowCreated is a free log subscription operation binding the contract event 0xd71bcf9347e4c440706076986a93115883dc9a480f43ec5d00610402a437be65.
//
// Solidity: event ICS20EscrowCreated(string indexed clientId, address indexed escrow)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchICS20EscrowCreated(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferICS20EscrowCreated, clientId []string, escrow []common.Address) (event.Subscription, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var escrowRule []interface{}
	for _, escrowItem := range escrow {
		escrowRule = append(escrowRule, escrowItem)
	}

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "ICS20EscrowCreated", clientIdRule, escrowRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferICS20EscrowCreated)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "ICS20EscrowCreated", log); err != nil {
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

// ParseICS20EscrowCreated is a log parse operation binding the contract event 0xd71bcf9347e4c440706076986a93115883dc9a480f43ec5d00610402a437be65.
//
// Solidity: event ICS20EscrowCreated(string indexed clientId, address indexed escrow)
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseICS20EscrowCreated(log types.Log) (*ContractICS20TransferICS20EscrowCreated, error) {
	event := new(ContractICS20TransferICS20EscrowCreated)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "ICS20EscrowCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS20TransferICS20EscrowLaunchGateEnabledIterator is returned from FilterICS20EscrowLaunchGateEnabled and is used to iterate over the raw logs and unpacked data for ICS20EscrowLaunchGateEnabled events raised by the ContractICS20Transfer contract.
type ContractICS20TransferICS20EscrowLaunchGateEnabledIterator struct {
	Event *ContractICS20TransferICS20EscrowLaunchGateEnabled // Event containing the contract specifics and raw log

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
func (it *ContractICS20TransferICS20EscrowLaunchGateEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS20TransferICS20EscrowLaunchGateEnabled)
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
		it.Event = new(ContractICS20TransferICS20EscrowLaunchGateEnabled)
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
func (it *ContractICS20TransferICS20EscrowLaunchGateEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS20TransferICS20EscrowLaunchGateEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS20TransferICS20EscrowLaunchGateEnabled represents a ICS20EscrowLaunchGateEnabled event raised by the ContractICS20Transfer contract.
type ContractICS20TransferICS20EscrowLaunchGateEnabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterICS20EscrowLaunchGateEnabled is a free log retrieval operation binding the contract event 0x602a536c4b89ff28b2a974a997fe044117170c12aa39af142bdbf06a097bb7fd.
//
// Solidity: event ICS20EscrowLaunchGateEnabled()
func (_ContractICS20Transfer *ContractICS20TransferFilterer) FilterICS20EscrowLaunchGateEnabled(opts *bind.FilterOpts) (*ContractICS20TransferICS20EscrowLaunchGateEnabledIterator, error) {

	logs, sub, err := _ContractICS20Transfer.contract.FilterLogs(opts, "ICS20EscrowLaunchGateEnabled")
	if err != nil {
		return nil, err
	}
	return &ContractICS20TransferICS20EscrowLaunchGateEnabledIterator{contract: _ContractICS20Transfer.contract, event: "ICS20EscrowLaunchGateEnabled", logs: logs, sub: sub}, nil
}

// WatchICS20EscrowLaunchGateEnabled is a free log subscription operation binding the contract event 0x602a536c4b89ff28b2a974a997fe044117170c12aa39af142bdbf06a097bb7fd.
//
// Solidity: event ICS20EscrowLaunchGateEnabled()
func (_ContractICS20Transfer *ContractICS20TransferFilterer) WatchICS20EscrowLaunchGateEnabled(opts *bind.WatchOpts, sink chan<- *ContractICS20TransferICS20EscrowLaunchGateEnabled) (event.Subscription, error) {

	logs, sub, err := _ContractICS20Transfer.contract.WatchLogs(opts, "ICS20EscrowLaunchGateEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS20TransferICS20EscrowLaunchGateEnabled)
				if err := _ContractICS20Transfer.contract.UnpackLog(event, "ICS20EscrowLaunchGateEnabled", log); err != nil {
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

// ParseICS20EscrowLaunchGateEnabled is a log parse operation binding the contract event 0x602a536c4b89ff28b2a974a997fe044117170c12aa39af142bdbf06a097bb7fd.
//
// Solidity: event ICS20EscrowLaunchGateEnabled()
func (_ContractICS20Transfer *ContractICS20TransferFilterer) ParseICS20EscrowLaunchGateEnabled(log types.Log) (*ContractICS20TransferICS20EscrowLaunchGateEnabled, error) {
	event := new(ContractICS20TransferICS20EscrowLaunchGateEnabled)
	if err := _ContractICS20Transfer.contract.UnpackLog(event, "ICS20EscrowLaunchGateEnabled", log); err != nil {
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
