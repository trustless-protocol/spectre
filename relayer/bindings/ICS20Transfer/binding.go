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
	Bin: "0x60a080604052346100c257306080525f5160206154f05f395f51905f525460ff8160401c166100b3576002600160401b03196001600160401b03821601610060575b60405161542990816100c7823960805181818161171901526117c90152f35b6001600160401b0319166001600160401b039081175f5160206154f05f395f51905f525581527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f80610041565b63f92ee8a960e01b5f5260045ffd5b5f80fdfe60806040526004361015610011575f80fd5b5f5f3560e01c806306ab20bc14612dd1578063078c4a79146126325780631459457a146122a05780631bbf2e231461225a5780631cf0a67f146122155780631e5150e4146121cf57806329b6eca914611fb95780632ac3dc3814611f185780633f4ba83a14611e57578063428e4e1714611a1a5780634f1ef2861461177857806352d1902d146116fe57806353816a7c1461144257806357a85bea146113f75780635c975abb146113b55780635e32b6b61461119d57806362bfa172146110605780637a9e5e4b14610fcf578063826cae7a14610f895780638456cb5914610eee5780638fb3603714610e8457806390f6d62014610e2a57806393dcfee414610bab578063969631d514610b2f578063a1d28f5714610819578063a50ee2b4146107bb578063aaa2c3431461070f578063ac9650d814610592578063ad3cb1cc14610531578063b29c715d146103be578063bf7e214f1461038b578063d413227d14610345578063dab8ed9c146102a25763e163b1af14610190575f80fd5b3461029f57602036600319011261029f576101e26101ac612e76565b6001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9060405191818154916101f4836134b7565b80865292600181169081156102755750600114610234575b6102308561021c81870382612f3e565b604051918291602083526020830190612eea565b0390f35b815260208120939250905b80821061025b5750909150810160200161021c8261023061020c565b91926001816020925483858801015201910190929161023f565b8695506102309693506020925061021c94915060ff191682840152151560051b820101929361020c565b80fd5b503461029f578060031936011261029f576102bd36336135bd565b7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065460ff8160a01c16156102ef575080f35b60ff60a01b1916600160a01b177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f806557f602a536c4b89ff28b2a974a997fe044117170c12aa39af142bdbf06a097bb7fd8180a180f35b503461029f578060031936011261029f5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416604051908152f35b503461029f578060031936011261029f5760206001600160a01b035f5160206153dd5f395f51905f525416604051908152f35b503461029f57604036600319011261029f5760043567ffffffffffffffff811161051657806004019060e0600319823603011261052d576103fd612e8c565b916104066137d8565b61040e613764565b61041836336135bd565b602482013591821561051a5761044b61044661043f60646001600160a01b03940185613063565b3691612f7c565b61382b565b16916104618161045a84613382565b85336139ff565b8461046b83613382565b91843b156105165760405163b4f22eb760e01b81526001600160a01b0393909316600484015233602484015260448301528160648183875af1801561050b576104f2575b6020856104bd868686613bdb565b907f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d67ffffffffffffffff60405191168152f35b6104fd858092612f3e565b610507575f6104af565b8380fd5b6040513d87823e3d90fd5b5080fd5b6024856304f6df8d60e41b815280600452fd5b8280fd5b503461029f578060031936011261029f5750610230604051610554604082612f3e565b600581527f352e302e300000000000000000000000000000000000000000000000000000006020820152604051918291602083526020830190612eea565b503461029f57602036600319011261029f5760043567ffffffffffffffff8111610516576105c4903690600401612fe0565b9060206040516105d48282612f3e565b84815281810191601f1981013684376105ec85613591565b936105fa6040519586612f3e565b858552601f1961060987613591565b0182885b8281106106ff57505050865b868110156106a05760019061068489806106708861063c8660051b890189613063565b8a8d6040959395519483869484860198893784019083820190898252519283915e010185815203601f198101835282612f3e565b5190305af461067d61416d565b90306147f0565b61068e82896135a9565b5261069981886135a9565b5001610619565b82888760405191838301848452825180915260408401948060408360051b870101940192955b8287106106d35785850386f35b9091929382806106ef600193603f198a82030186528851612eea565b96019201960195929190926106c6565b606082828a01015201839061060d565b503461029f57602036600319011261029f578061072a612e76565b61073436336135bd565b6001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805541690813b156107b7576001600160a01b0360248492836040519586948593631b2ce7f360e11b85521660048401525af180156107ac5761079b5750f35b816107a591612f3e565b61029f5780f35b6040513d84823e3d90fd5b5050fd5b503461029f57602036600319011261029f576004359067ffffffffffffffff821161029f5760206108116107f23660048601612fb2565b6001600160a01b03610806828495946133ce565b541692831515613406565b604051908152f35b503461029f57604036600319011261029f5760043567ffffffffffffffff81116105165761084b903690600401612fb2565b610856929192612e8c565b61086036336135bd565b6001600160a01b0361087283866133ce565b5416610aef576108bb6108b5826001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b546134b7565b156108f6826001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9015610a33575061095c9061090b83866133ce565b6001600160a01b03808316166001600160a01b03198254161790556001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9067ffffffffffffffff8111610a1f576109808161097a84546134b7565b846134ef565b82601f82116001146109bf57819084956109af9495926109b4575b50508160011b915f199060031b1c19161790565b905580f35b013590505f8061099b565b601f198216948385526020852091855b878110610a075750836001959697106109ee575b505050811b01905580f35b01355f19600384901b60f8161c191690555f80806109e3565b909260206001819286860135815501940191016109cf565b602483634e487b7160e01b81526041600452fd5b83906040519182917f778769c40000000000000000000000000000000000000000000000000000000083526020600484015281815491610a72836134b7565b928360248701526001811690815f14610acf5750600114610a95575b5050500390fd5b9080935052602082205b818310610ab55750508101604401838080610a8e565b805460448487010152849350602090920191600101610a9f565b925050506044925060ff191682840152151560051b820101838080610a8e565b6040517f0c0ef5340000000000000000000000000000000000000000000000000000000081526020600482015280610b2b6024820185886130d1565b0390fd5b503461029f57602036600319011261029f5760043567ffffffffffffffff81116105165760209182610b6d6001600160a01b03933690600401612fb2565b925082604051938492833781017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8008152030190205416604051908152f35b503461029f57604036600319011261029f5760043567ffffffffffffffff811161051657610bdd903690600401612fb2565b60243567ffffffffffffffff811161050757610bfd903690600401612fe0565b9290610c0936336135bd565b6001600160a01b03604051848482376020818681017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8008152030190205416938415610e07578015610dcb57855b818110610cad57505050610c6a8282613396565b600160ff1982541617905581604051928392833781018481520390207f4091a6ce2a3c52563f54d5437475ea54363fc120a36395a84533922c2d5ad9918380a380f35b6001600160a01b03610cc8610cc3838587613436565b613382565b16151580610daf575b610ceb908686610ce5610cc386888a613436565b9261345a565b610cf9610cc3828486613436565b906001600160a01b03604051927f0b0aee690000000000000000000000000000000000000000000000000000000084521660048301526020826024818a5afa918215610da4578892610d6b575b50610d656001928787610d5d610cc386898b613436565b92151561345a565b01610c56565b91506020823d8211610d9c575b81610d8560209383612f3e565b81010312610d9857905190610d65610d46565b5f80fd5b3d9150610d78565b6040513d8a823e3d90fd5b50610ceb610dc1610cc3838587613436565b3b15159050610cd1565b6040517fd5c055800000000000000000000000000000000000000000000000000000000081526020600482015280610b2b6024820187876130d1565b6040516358d1c8a160e11b81526020600482015280610b2b6024820187876130d1565b503461029f57602036600319011261029f576004359067ffffffffffffffff821161029f5760206001600160a01b03610e7b610e76610e6c3660048801612fb2565b61043f36336135bd565b61442b565b16604051908152f35b503461029f578060031936011261029f575f5160206153dd5f395f51905f525460a01c60ff1615610ee6575060207f8fb36037000000000000000000000000000000000000000000000000000000005b6001600160e01b031960405191168152f35b602090610ed4565b503461029f578060031936011261029f57610f0936336135bd565b610f116137d8565b600160ff197fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005416177fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a180f35b503461029f578060031936011261029f5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8045416604051908152f35b503461029f57602036600319011261029f57610fe9612e76565b6001600160a01b035f5160206153dd5f395f51905f525416330361104e57803b1561101a57611017906143cb565b80f35b7fc2f31e5e0000000000000000000000000000000000000000000000000000000082526001600160a01b0316600452602490fd5b60248262d1953b60e31b815233600452fd5b503461029f57608036600319011261029f578060043567ffffffffffffffff811161119a57611093903690600401612fb2565b60243567ffffffffffffffff8111611197576110b3903690600401612fb2565b91909260443567ffffffffffffffff8111611193576110d6903690600401612fb2565b94906064359360ff851680950361118f57611112906110f536336135bd565b6001600160a01b0361110782876133ce565b541694851515613406565b823b1561118b578694611162946111748793604051998a98899788967f37d2c2f40000000000000000000000000000000000000000000000000000000088526060600489015260648801916130d1565b858103600319016024870152916130d1565b90604483015203925af180156107ac5761079b5750f35b8680fd5b8780fd5b8580fd5b50505b50fd5b503461029f576111ac36612eb6565b6111e2336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314613011565b6111ea613764565b60608101908261123b61121561120d611203868661304e565b6080810190613063565b8101906131ee565b611228611222868661304e565b80613063565b906112338680613063565b92909161419c565b905061124681614969565b611272575b50807f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d80f35b6001600160a01b031692833b156105165761133e82916080946001600160a01b0361134560405197889586957f5e32b6b60000000000000000000000000000000000000000000000000000000087526020600488015261132d61130b6112ec6112db898061407d565b60a060248d015260c48c01916130d1565b6112f960208a018a61407d565b8b83036023190160448d0152906130d1565b9167ffffffffffffffff61132160408a016140af565b1660648a0152876140c4565b8782036023190160848901526140d8565b9301612ea2565b1660a483015203818387620186a0f191826113a0575b505061139a577f02c075b43cbb5fee7c14de7aea6a6ddde6e3646144712fd9c7159705152537e661138d61021c61416d565b0390a25b5f80828161124b565b50611391565b816113aa91612f3e565b61052d57825f61135b565b503461029f578060031936011261029f57602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b503461029f57602036600319011261029f576004359067ffffffffffffffff821161029f57602060ff6114366114303660048701612fb2565b90613396565b54166040519015158152f35b503461029f5760c036600319011261029f5760043567ffffffffffffffff811161051657806004019060e0600319823603011261052d57608036602319011261052d5760a43567ffffffffffffffff8111610507576114a5903690600401612fb2565b6114b09291926137d8565b6114b8613764565b60248201359182156116eb576114cc61336c565b6001600160a01b03806114de88613382565b169116146114ea61336c565b6114f387613382565b91156116b157505061044661043f606461150e930187613063565b926001600160a01b03807f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065416941693604051906040820182811067ffffffffffffffff82111761169d579088949392916040528682526020820191868352813b1561119357856001600160a01b03916115fd8296604051988997889687957f30f28b7a000000000000000000000000000000000000000000000000000000008752816115b9612e8c565b166004880152604435602488015260643560448801526084356064880152511660848601525160a48501523360c485015261010060e48501526101048401916130d1565b03925af180156107ac57611688575b5061161684613382565b91833b156105165760405163b4f22eb760e01b81526001600160a01b0393909316600484015233602484015260448301528160648183865af1801561167d57611668575b6020846104bd338587613bdb565b611673848092612f3e565b61052d575f61165a565b6040513d86823e3d90fd5b8161169291612f3e565b61050757835f61160c565b602489634e487b7160e01b81526041600452fd5b7fe36164780000000000000000000000000000000000000000000000000000000088526001600160a01b0390811660045216602452604486fd5b6024866304f6df8d60e41b815280600452fd5b503461029f578060031936011261029f576001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036117695760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b8063703e46dd60e11b60049252fd5b50604036600319011261029f5761178d612e76565b9060243567ffffffffffffffff81116105165736602382011215610516576117bf903690602481600401359101612f7c565b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000168030149081156119e5575b506119d65761180336336135bd565b6001600160a01b03831690604051937f52d1902d000000000000000000000000000000000000000000000000000000008552602085600481865afa8095859661199e575b5061185f5760248484634c9c8ce360e01b8252600452fd5b9091847f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc81036119735750813b1561196157806001600160a01b03197f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5416177f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b8480a2815183901561192e578083602061192a95519101845af461192461416d565b916147f0565b5080f35b505050346119395780f35b807fb398979f0000000000000000000000000000000000000000000000000000000060049252fd5b634c9c8ce360e01b8452600452602483fd5b7faa1d49a4000000000000000000000000000000000000000000000000000000008552600452602484fd5b9095506020813d6020116119ce575b816119ba60209383612f3e565b810103126119ca5751945f611847565b8480fd5b3d91506119ad565b60048263703e46dd60e11b8152fd5b90506001600160a01b037f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc541614155f6117f4565b503461029f57602036600319011261029f5760043567ffffffffffffffff811161051657806004019060c0600319823603011261052d57611a87336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314613011565b611a8f613764565b606481019183611aa561120d611203868561304e565b926084810192611ab861043f8583613063565b60208151910120946040956020808851611ad28a82612f3e565b818152017f4774d4a575993f963b1c06573736617a457abef8589178db8d10c94b4ab511ab815220145f14611cc857611b1d90611b12611222898561304e565b906112338580613063565b9050611b2881614969565b611b5b575b505050505050505b807f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d80f35b6001600160a01b031695863b1561050757839286519586938493638dfcd9ad60e01b8552866004860152896024860152611b95828061407d565b6044870160c09052610104870190611bac926130d1565b611bb9602486018461407d565b878303604319016064890152611bcf92916130d1565b90611bdc604486016140af565b67ffffffffffffffff166084870152611bf590836140c4565b8582036043190160a4870152611c0b91906140d8565b91611c159161407d565b8483036043190160c4860152611c2b92916130d1565b9060a401611c3890612ea2565b6001600160a01b031660e483015203818388620186a0f19182611cb3575b5050611cac57611c9c7f898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e091611c8961416d565b9051918291602083526020830190612eea565b0390a25b5f808083818080611b2d565b5050611ca0565b81611cbd91612f3e565b61050757835f611c56565b60200151611ce1611cda825183614622565b92906132bb565b611cea81614969565b611cfb575b50505050505050611b35565b6001600160a01b031695863b1561050757839286519586938493638dfcd9ad60e01b85526004850160019052896024860152611d37828061407d565b6044870160c09052610104870190611d4e926130d1565b611d5b602486018461407d565b878303604319016064890152611d7192916130d1565b90611d7e604486016140af565b67ffffffffffffffff166084870152611d9790836140c4565b8582036043190160a4870152611dad91906140d8565b91611db79161407d565b8483036043190160c4860152611dcd92916130d1565b9060a401611dda90612ea2565b6001600160a01b031660e483015203818388620186a0f19182611e42575b5050611e3b57611e2b7f898f81cc34db3cf328a50a3f019f67f2f82077674357844e51ac4445baa466e091611c8961416d565b0390a25b5f808083818080611cef565b5050611e2f565b81611e4c91612f3e565b61050757835f611df8565b503461029f578060031936011261029f57611e7236336135bd565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff811615611ef05760ff19167fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a180f35b6004827f8dfc202b000000000000000000000000000000000000000000000000000000008152fd5b503461029f57602036600319011261029f5760043567ffffffffffffffff811161051657806004019060e0600319823603011261052d57611f576137d8565b611f5f613764565b6024810135908115611fa657611f8661044661043f60646001600160a01b03940186613063565b1690611f9c81611f9585613382565b84336139ff565b8361161684613382565b6024846304f6df8d60e41b815280600452fd5b503461029f57602036600319011261029f57611fd3612e76565b5f5160206153fd5f395f51905f525467ffffffffffffffff811690600182036121c05760401c60ff169081156121b4575b506121a55761203a600267ffffffffffffffff195f5160206153fd5f395f51905f525416175f5160206153fd5f395f51905f5255565b6801000000000000000068ff0000000000000000195f5160206153fd5f395f51905f525416175f5160206153fd5f395f51905f5255602460206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416604051928380927f24d7806c0000000000000000000000000000000000000000000000000000000082523360048301525afa90811561219a57839161215b575b50906120ee612103923390613011565b6120f66146bd565b6120fe6146bd565b6143cb565b68ff0000000000000000195f5160206153fd5f395f51905f5254165f5160206153fd5f395f51905f52557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a180f35b90506020813d602011612192575b8161217660209383612f3e565b8101031261052d575190811515820361052d57906120ee6120de565b3d9150612169565b6040513d85823e3d90fd5b60048263f92ee8a960e01b8152fd5b6002915010155f612004565b60048463f92ee8a960e01b8152fd5b503461029f578060031936011261029f5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8055416604051908152f35b503461029f578060031936011261029f57602060ff7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065460a01c166040519015158152f35b503461029f578060031936011261029f5760206001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065416604051908152f35b503461029f5760a036600319011261029f576122ba612e76565b6122c2612e8c565b6044356001600160a01b038116810361050757606435926001600160a01b0384168094036119ca576084356001600160a01b0381168103611193575f5160206153fd5f395f51905f525467ffffffffffffffff811690816126235760401c60ff16908115612617575b5061260857906123bb6001600160a01b039261236e600267ffffffffffffffff195f5160206153fd5f395f51905f525416175f5160206153fd5f395f51905f5255565b6801000000000000000068ff0000000000000000195f5160206153fd5f395f51905f525416175f5160206153fd5f395f51905f52556123ab6146bd565b6123b36146bd565b6120ee6146bd565b166001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8035416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80355604051916104009283810181811067ffffffffffffffff8211176125f45761244d8291614b3595878785396001600160a01b0316815230602082015260400190565b039086f0801561050b576001600160a01b03166001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8045416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80455604051928084019284841067ffffffffffffffff8511176125f457918493916124eb9385396001600160a01b0316815230602082015260400190565b039083f080156107ac576001600160a01b03166001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8055416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805556001600160a01b03197f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065416177f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065568ff0000000000000000195f5160206153fd5f395f51905f5254165f5160206153fd5f395f51905f52557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a180f35b602487634e487b7160e01b81526041600452fd5b60048663f92ee8a960e01b8152fd5b6002915010155f61232b565b60048863f92ee8a960e01b8152fd5b503461029f5761264136612eb6565b90612678336001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80354163314613011565b612680613764565b6126886137d8565b60608201916126a761043f61269d858461304e565b6040810190613063565b602081519101206126b6613096565b60208151910120146126c6613096565b6126d361269d868561304e565b909215612d9b575050506127256126f061043f611222868561304e565b602081519101206126ff613119565b602081519101201461270f613119565b9061271d611222878661304e565b929091613154565b61273f61043f612735858461304e565b6060810190613063565b6020815191012061274e613198565b602081519101201461275e613198565b61276b612735868561304e565b909215612d65575050506127bf61279261043f612788868561304e565b6020810190613063565b602081519101206127a1613119565b60208151910120146127b1613119565b9061271d612788878661304e565b6127cf61120d611203858461304e565b9260608401805115611fa657604085019161280d8351936127fb6127f4865187614622565b96906132bb565b516001600160a01b03851615156132bb565b602084019361282261044661043f8784613063565b965194612851612835611222858561304e565b61284c6128428680613063565b9390923691612f7c565b613989565b92865184511080612d55575b156129ea5750509051845195968795909250906020908781189088110287188061288f61288a828b61353e565b61355f565b9803920101602087015e6128a4855186614622565b9590158681156129d8575b506129af575b506001600160a01b03905b16905193813b15610507576040517f0779afe60000000000000000000000000000000000000000000000000000000081526001600160a01b039182166004820152921660248301526044820193909352918290606490829084905af180156107ac5761299a575b610230826040519061293a604083612f3e565b601182527f7b22726573756c74223a2241513d3d227d00000000000000000000000000000060208301527f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d604051918291602083526020830190612eea565b6129a5828092612f3e565b61029f575f612927565b94506001600160a01b03906129d2826129c788613301565b54169687151561333f565b906128b5565b6001600160a01b03915016155f6128af565b612a549350612a1f83602098949361284c612842612a17612a0e8d9897899861304e565b87810190613063565b939094613063565b926040519784899551918291018487015e8401908282018a8152815193849201905e010186815203601f198101855284612f3e565b6001600160a01b038516946001600160a01b03612a7085613301565b5416938415612af8575b5084956001600160a01b038596959616908351823b1561118b576040516340c10f1960e01b81526001600160a01b0392909216600483015260248201529085908290604490829084905af190811561050b578591612ae3575b50506001600160a01b03906128c0565b81612aed91612f3e565b61050757835f612ad3565b93506001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80454166040517f4571e3a600000000000000000000000000000000000000000000000000000000602082015230602482015287604482015260606064820152612b8081612b726084820189612eea565b03601f198101835282612f3e565b604051916104a88084019084821067ffffffffffffffff831117612d415791849391612bb093614f3586396139d2565b039086f0801561050b576001600160a01b031695612bcd85613301565b6001600160a01b0388166001600160a01b0319825416179055612c20876001600160a01b03165f527f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80260205260405f2090565b9685519767ffffffffffffffff8911612d2d57612c4789612c4183546134b7565b836134ef565b602098601f8111600114612ccb5780612c78918a9b8b9a9b91612cc0575b508160011b915f199060031b1c19161790565b90555b7f6031fab685dd6d86e4dbac9a69eae347145f332c95b3a0d728d3730fc5233d62612cb58298604051918291602083526020830190612eea565b0390a2959493612a7a565b90508a01515f612c65565b818952898920601f1982168a5b818110612d155750908a9b83600194939c9b9c10612cfd575b5050811b019055612c7b565b8b01515f1960f88460031b161c191690555f80612cf1565b8a8d0151835560209c8d019c60019093019201612cd8565b602488634e487b7160e01b81526041600452fd5b60248a634e487b7160e01b81526041600452fd5b50612d6084886145d0565b61285d565b610b2b906040519384937fd1ca953a000000000000000000000000000000000000000000000000000000008552600485016130f1565b610b2b906040519384937f094af3b8000000000000000000000000000000000000000000000000000000008552600485016130f1565b5034610d98576020366003190112610d9857612deb612e76565b612df536336135bd565b6001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f804541690813b15610d98576001600160a01b0360245f92836040519586948593631b2ce7f360e11b85521660048401525af18015612e6b57612e5d575080f35b612e6991505f90612f3e565b005b6040513d5f823e3d90fd5b600435906001600160a01b0382168203610d9857565b602435906001600160a01b0382168203610d9857565b35906001600160a01b0382168203610d9857565b6020600319820112610d98576004359067ffffffffffffffff8211610d985760a0908290036003190112610d985760040190565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b60a0810190811067ffffffffffffffff821117612f2a57604052565b634e487b7160e01b5f52604160045260245ffd5b90601f8019910116810190811067ffffffffffffffff821117612f2a57604052565b67ffffffffffffffff8111612f2a57601f01601f191660200190565b929192612f8882612f60565b91612f966040519384612f3e565b829481845281830111610d98578281602093845f960137010152565b9181601f84011215610d985782359167ffffffffffffffff8311610d985760208381860195010111610d9857565b9181601f84011215610d985782359167ffffffffffffffff8311610d98576020808501948460051b010111610d9857565b156130195750565b6001600160a01b03907f2ecb3242000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b903590609e1981360301821215610d98570190565b903590601e1981360301821215610d98570180359067ffffffffffffffff8211610d9857602001918136038313610d9857565b604051906130a5604083612f3e565b600782527f69637332302d31000000000000000000000000000000000000000000000000006020830152565b908060209392818452848401375f828201840152601f01601f1916010190565b916131086131169492604085526040850190612eea565b9260208185039101526130d1565b90565b60405190613128604083612f3e565b600882527f7472616e736665720000000000000000000000000000000000000000000000006020830152565b929091921561316257505050565b610b2b906040519384937f5d3a3cdd000000000000000000000000000000000000000000000000000000008552600485016130f1565b604051906131a7604083612f3e565b601a82527f6170706c69636174696f6e2f782d736f6c69646974792d6162690000000000006020830152565b9080601f83011215610d985781602061311693359101612f7c565b602081830312610d985780359067ffffffffffffffff8211610d98570160a081830312610d98576040519161322283612f0e565b813567ffffffffffffffff8111610d98578161323f9184016131d3565b8352602082013567ffffffffffffffff8111610d9857816132619184016131d3565b6020840152604082013567ffffffffffffffff8111610d9857816132869184016131d3565b604084015260608201356060840152608082013567ffffffffffffffff8111610d98576132b392016131d3565b608082015290565b156132c35750565b610b2b906040519182917f3fed5d87000000000000000000000000000000000000000000000000000000008352602060048401526024830190612eea565b60208091604051928184925191829101835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80181520301902090565b156133475750565b610b2b9060405191829163e1275e2f60e01b8352602060048401526024830190612eea565b6024356001600160a01b0381168103610d985790565b356001600160a01b0381168103610d985790565b60209082604051938492833781017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80781520301902090565b60209082604051938492833781017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80181520301902090565b91909115613412575050565b610b2b60405192839263e1275e2f60e01b84526020600485015260248401916130d1565b91908110156134465760051b0190565b634e487b7160e01b5f52603260045260245ffd5b9290921561346757505050565b6001600160a01b036134ac6040519485947f28203fe70000000000000000000000000000000000000000000000000000000086526040600487015260448601916130d1565b911660248301520390fd5b90600182811c921680156134e5575b60208310146134d157565b634e487b7160e01b5f52602260045260245ffd5b91607f16916134c6565b601f82116134fc57505050565b5f5260205f20906020601f840160051c83019310613534575b601f0160051c01905b818110613529575050565b5f815560010161351e565b9091508190613515565b9190820391821161354b57565b634e487b7160e01b5f52601160045260245ffd5b9061356982612f60565b6135766040519182612f3e565b8281528092613587601f1991612f60565b0190602036910137565b67ffffffffffffffff8111612f2a5760051b60200190565b80518210156134465760209160051b010190565b5f5160206153dd5f395f51905f5254916001600160a01b0383169281600411610d98575f5f9060405f8151966001600160a01b0360208901917fb700961300000000000000000000000000000000000000000000000000000000835216978860248201523060448201526001600160e01b0319833516606482015260648152613647608482612f3e565b828052826020525190895afa613751575b15613665575b5050505050565b63ffffffff161561373f5760ff60a01b1916600160a01b175f5160206153dd5f395f51905f5255823b15610d98576020925f92836040518096819582947f94c7d7ee00000000000000000000000000000000000000000000000000000000845260048401526040602484015260448301908082528085848401378181018301859052601f01601f1916010103925af18015612e6b5761372f575b5060ff60a01b195f5160206153dd5f395f51905f5254165f5160206153dd5f395f51905f52555f8080808061365e565b5f61373991612f3e565b5f6136ff565b8262d1953b60e31b5f5260045260245ffd5b50505f516020518060201c150290613658565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005c6137b05760017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d565b7f3ee5aeb5000000000000000000000000000000000000000000000000000000005f5260045ffd5b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300541661380357565b7fd93c0665000000000000000000000000000000000000000000000000000000005f5260045ffd5b604051906001600160a01b03815192602081818501958087835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80081520301902054169182155f146138d65750905060ff7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065460a01c166138b1576131169061442b565b610b2b906040519182916358d1c8a160e11b8352602060048401526024830190612eea565b60ff7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8065460a01c1661390757505090565b60ff906020604051809285518091835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f80781520301902054161561394b575090565b610b2b906040519182917f0d2e2ae3000000000000000000000000000000000000000000000000000000008352602060048401526024830190612eea565b6001809160208095613116958160405198858a9651918291018688015e850191602f60f81b8584015260218301370101602f60f81b838201520301601e19810184520182612f3e565b6040906001600160a01b0361311694931681528160208201520190612eea565b9190820180921161354b57565b916001600160a01b03909391931692604051926370a0823160e01b84526001600160a01b03821691826004860152602085602481895afa948515612e6b575f95613ba6575b506040517f23b872dd0000000000000000000000000000000000000000000000000000000060208281019182526001600160a01b039485166024840152939092166044820152606481018590525f9190613aa18160848101612b72565b519082885af115612e6b575f513d613b9d5750833b155b613b71576020906024604051809681936370a0823160e01b835260048301525afa928315612e6b575f93613b3b575b50613af290826139f2565b90821180613b32575b15613b04575050565b7f2fb30cfc000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b50808214613afb565b9092506020813d602011613b69575b81613b5760209383612f3e565b81010312610d98575191613af2613ae7565b3d9150613b4a565b837f5274afe7000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b60011415613ab8565b9094506020813d602011613bd3575b81613bc260209383612f3e565b81010312610d985751936020613a44565b3d9150613bb5565b613be76101ac82613382565b5f9260405191825f825492613bfb846134b7565b808452936001811690811561405b5750600114614017575b50613c2092500383612f3e565b818051155f14613f3d5750505060a0613c49613c43613c3e84613382565b614701565b94614701565b93613c576040840184613063565b939095613ca8613c8d613c6d60c0850185613063565b99909760405196613c7d88612f0e565b8752602087019485523691612f7c565b9560408501968752606085019860208501358a523691612f7c565b94608084019586526001600160a01b037f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f803541695613ce96060850185613063565b96909401359467ffffffffffffffff8616809603613f395790613da3613d56949392613d95613d16613119565b9c613d1f613119565b94613d7c613d2b613096565b97613d69613d37613198565b9a6040519c8d986020808b01525160a060408b015260e08a0190612eea565b9051888203603f190160608a0152612eea565b9051868203603f19016080880152612eea565b915160a085015251838203603f190160c0850152612eea565b03601f198101865285612f3e565b60405199613db08b612f0e565b8a5260208a0152604089015260608801526080870152604051926060840184811067ffffffffffffffff8211176125f45793613ed360209694613e06613e5a958a9567ffffffffffffffff996040523691612f7c565b835287830190815260408301998a52604051998a97889687957f4d6e7ce30000000000000000000000000000000000000000000000000000000087528b600488015251606060248801526084870190612eea565b925116604485015251906023198482030160648501526080613ec2613eb0613e9e613e8e865160a0875260a0870190612eea565b8d8701518682038f880152612eea565b60408601518582036040870152612eea565b60608501518482036060860152612eea565b920151906080818403910152612eea565b03925af1918215613f2c578192613ee957505090565b9091506020813d602011613f24575b81613f0560209383612f3e565b8101031261051657519067ffffffffffffffff8216820361029f575090565b3d9150613ef8565b50604051903d90823e3d90fd5b8880fd5b613f5e613f4b969396613119565b613f586060870187613063565b91613989565b81518151109182614006575b5050613f7d575b50613c4960a091614701565b6001600160a01b03613f8e84613382565b16803b15610d98576040517f9dc29fac0000000000000000000000000000000000000000000000000000000081526001600160a01b03929092166004830152602084013560248301525f908290604490829084905af18015612e6b5715613f7157613ffc9193505f90612f3e565b5f91613c49613f71565b61401092506145d0565b5f80613f6a565b90505f9291925260205f20905f915b81831061403f575050906020613c20928201015f613c13565b6020919350806001915483858901015201910190918492614026565b905060209250613c2094915060ff191682840152151560051b8201015f613c13565b9035601e1982360301811215610d9857016020813591019167ffffffffffffffff8211610d98578136038313610d9857565b359067ffffffffffffffff82168203610d9857565b9035609e1982360301811215610d98570190565b6131169161415f61415461413961411e6141036140f5878061407d565b60a0885260a08801916130d1565b614110602088018861407d565b9087830360208901526130d1565b61412b604087018761407d565b9086830360408801526130d1565b614146606086018661407d565b9085830360608701526130d1565b92608081019061407d565b9160808185039101526130d1565b3d15614197573d9061417e82612f60565b9161418c6040519384612f3e565b82523d5f602084013e565b606090565b909493925f926001600160a01b03604051838382376020818581017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f800815203019020541692831561438c579061284c61420b939260208801519961043f6142048c518d614622565b9c906132bb565b8351908151815110918261437b575b505015614324576001600160a01b036142338451613301565b541692614243815185151561333f565b6060810151843b15610d98576040516340c10f1960e01b81526001600160a01b038416600482015260248101919091525f8160448183895af18015612e6b5761430e575b506060905b0151813b1561052d576040517fb6163ac40000000000000000000000000000000000000000000000000000000081526001600160a01b0385811660048301528716602482015260448101919091529082908290606490829084905af180156107ac576142f9575b50509190565b614304828092612f3e565b61029f57806142f3565b61431b9193505f90612f3e565b5f916060614287565b6001600160a01b036143368451613301565b5416928315614348575b60609061428c565b9250606083519361435d6127f4865187614622565b614374856001600160a01b0383519116151561333f565b9050614340565b61438592506145d0565b5f8061421a565b5090610b2b6040519283927f5778f3780000000000000000000000000000000000000000000000000000000084526020600485015260248401916130d1565b60206001600160a01b037f2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad9216806001600160a01b03195f5160206153dd5f395f51905f525416175f5160206153dd5f395f51905f5255604051908152a1565b6040516001600160a01b03825191602081818601948086835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f800815203019020541691821561447b57505090565b7f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f805545f5160206153dd5f395f51905f52546040517f485cc9550000000000000000000000000000000000000000000000000000000060208201523060248201526001600160a01b03918216604480830191909152815293945091929116614503606483612f3e565b604051916104a8908184019284841067ffffffffffffffff851117612f2a57849361453293614f3586396139d2565b03905ff08015612e6b576001600160a01b031691829160405160208183518086835e81017f823f7a8ea9ae6df0eb03ec5e1682d7a2839417ad8a91774118e6acf2e8d2f8008152030190206001600160a01b0384166001600160a01b0319825416179055604051918291518091835e81015f81520390207fd71bcf9347e4c440706076986a93115883dc9a480f43ec5d00610402a437be655f80a390565b805190825180921061461b57805191828082109118028083189214158202821890602061460061288a848661353e565b92808285019503920101835e51902090602081519101201490565b5050505f90565b8051821180156146b6575b61467a576001821180614682575b158015908160011b918204600214171561354b576028018060281161354b57820361467a576001600160a01b0392915f6146749261487c565b90921690565b50505f905f90565b5061060f60f31b7fffff0000000000000000000000000000000000000000000000000000000000006020830151161461463b565b505f61462d565b60ff5f5160206153fd5f395f51905f525460401c16156146d957565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f5260045ffd5b6001600160a01b031680614715602a612f60565b916147236040519384612f3e565b602a8352614731602a612f60565b6020840190601f19013682378351156134465760309053825160011015613446576078602184015360295b6001811161479d575061476d575090565b7fe22e27eb000000000000000000000000000000000000000000000000000000005f52600452601460245260445ffd5b90600f81166010811015613446578451831015613446577f3031323334353637383961626364656600000000000000000000000000000000901a8483016020015360041c90801561354b575f190161475c565b9061482d575080511561480557602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b81511580614873575b61483e575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15614836565b9290926001840180851161354b57831180614933575b158015908160011b918204600214171561354b576148b5905f94929394956139f2565b915b8183106148c75750505060019190565b9092919360ff6148fe7fff000000000000000000000000000000000000000000000000000000000000006020888601015116614ac2565b16600f8111614928578160041b918083046010149015171561354b576001910194019192906148b7565b505f94508493505050565b5061060f60f31b7fffff000000000000000000000000000000000000000000000000000000000000602086840101511614614892565b60205f604051828101906301ffc9a760e01b82526301ffc9a760e01b602482015260248152614999604482612f3e565b519084617530fa903d5f519083614ab6575b5082614aac575b5081614a44575b816149c2575090565b602091505f90604051838101906301ffc9a760e01b82526001600160e01b03197fd3ce6f1b0000000000000000000000000000000000000000000000000000000016602482015260248152614a18604482612f3e565b5191617530fa5f513d82614a38575b5081614a31575090565b9050151590565b6020111591505f614a27565b905060205f604051828101906301ffc9a760e01b82526001600160e01b0319602482015260248152614a77604482612f3e565b519084617530fa5f513d82614aa0575b5081614a96575b5015906149b9565b905015155f614a8e565b6020111591505f614a87565b151591505f6149b2565b6020111592505f6149ab565b60f81c602f811180614b2a575b15614ade57602f190160ff1690565b6060811180614b20575b15614af7576056190160ff1690565b6040811180614b16575b15614b10576036190160ff1690565b5060ff90565b5060478110614b01565b5060678110614ae8565b50603a8110614acf56fe60803461013457601f61040038819003918201601f19168301916001600160401b03831184841017610138578084926040948552833981010312610134576100468161014c565b906001600160a01b039061005c9060200161014c565b16908115610121575f80546001600160a01b031981168417825560405193916001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a3803b1561010157600180546001600160a01b0319166001600160a01b039290921691821790557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a261029f90816101618239f35b63211eb15960e21b5f9081526001600160a01b0391909116600452602490fd5b631e4fbdf760e01b5f525f60045260245ffd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b03821682036101345756fe60806040526004361015610011575f80fd5b5f3560e01c80633659cfe6146101af5780635c60da1b14610189578063715018a6146101255780638da5cb5b146101005763f2fde38b14610050575f80fd5b346100fc5760203660031901126100fc576004356001600160a01b0381168091036100fc5761007d610253565b80156100d0576001600160a01b035f548273ffffffffffffffffffffffffffffffffffffffff198216175f55167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3005b7f1e4fbdf7000000000000000000000000000000000000000000000000000000005f525f60045260245ffd5b5f80fd5b346100fc575f3660031901126100fc5760206001600160a01b035f5416604051908152f35b346100fc575f3660031901126100fc5761013d610253565b5f6001600160a01b03815473ffffffffffffffffffffffffffffffffffffffff1981168355167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b346100fc575f3660031901126100fc5760206001600160a01b0360015416604051908152f35b346100fc5760203660031901126100fc576004356001600160a01b038116908181036100fc576101dd610253565b3b15610228578073ffffffffffffffffffffffffffffffffffffffff1960015416176001557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a2005b7f847ac564000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b6001600160a01b035f5416330361026657565b7f118cdaa7000000000000000000000000000000000000000000000000000000005f523360045260245ffdfea164736f6c634300081c000a60a0806040526104a880380380916100178285610292565b833981016040828203126101eb5761002e826102c9565b602083015190926001600160401b0382116101eb57019080601f830112156101eb57815161005b816102dd565b926100696040519485610292565b8184526020840192602083830101116101eb57815f926020809301855e84010152823b15610274577fa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d5080546001600160a01b0319166001600160a01b038516908117909155604051635c60da1b60e01b8152909190602081600481865afa9081156101f7575f9161023a575b50803b1561021a5750817f1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e5f80a282511561020257602060049260405193848092635c60da1b60e01b82525afa9182156101f7575f926101ae575b505f809161018a945190845af43d156101a6573d9161016e836102dd565b9261017c6040519485610292565b83523d5f602085013e6102f8565b505b608052604051610151908161035782396080518160460152f35b6060916102f8565b9291506020833d6020116101ef575b816101ca60209383610292565b810103126101eb575f80916101e161018a956102c9565b9394509150610150565b5f80fd5b3d91506101bd565b6040513d5f823e3d90fd5b505050341561018c5763b398979f60e01b5f5260045ffd5b634c9c8ce360e01b5f9081526001600160a01b0391909116600452602490fd5b90506020813d60201161026c575b8161025560209383610292565b810103126101eb57610266906102c9565b5f6100f5565b3d9150610248565b631933b43b60e21b5f9081526001600160a01b038416600452602490fd5b601f909101601f19168101906001600160401b038211908210176102b557604052565b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b03821682036101eb57565b6001600160401b0381116102b557601f01601f191660200190565b9061031c575080511561030d57602081519101fd5b63d6bda27560e01b5f5260045ffd5b8151158061034d575b61032d575090565b639996b31560e01b5f9081526001600160a01b0391909116600452602490fd5b50803b1561032556fe60806040527f5c60da1b000000000000000000000000000000000000000000000000000000006080526020608060048173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa80156100e9575f9015610127575060203d6020116100e2575b601f19601f820116608001906080821067ffffffffffffffff8311176100b5576100b0916040526080016100f4565b610127565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b503d610081565b6040513d5f823e3d90fd5b602090607f1901126101235760805173ffffffffffffffffffffffffffffffffffffffff811681036101235790565b5f80fd5b5f8091368280378136915af43d5f803e15610140573d5ff35b3d5ffdfea164736f6c634300081c000af3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a00f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00a164736f6c634300081c000af0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
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
