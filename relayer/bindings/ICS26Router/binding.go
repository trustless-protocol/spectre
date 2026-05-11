// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractICS26Router

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

// IICS02ClientMsgsCounterpartyInfo is an auto generated low-level Go binding around an user-defined struct.
type IICS02ClientMsgsCounterpartyInfo struct {
	ClientId     string
	MerklePrefix [][]byte
}

// IICS26RouterMsgsMsgAckPacket is an auto generated low-level Go binding around an user-defined struct.
type IICS26RouterMsgsMsgAckPacket struct {
	Packet          IICS26RouterMsgsPacket
	Acknowledgement []byte
	MembershipMsg   []byte
}

// IICS26RouterMsgsMsgRecvPacket is an auto generated low-level Go binding around an user-defined struct.
type IICS26RouterMsgsMsgRecvPacket struct {
	Packet        IICS26RouterMsgsPacket
	MembershipMsg []byte
}

// IICS26RouterMsgsMsgSendPacket is an auto generated low-level Go binding around an user-defined struct.
type IICS26RouterMsgsMsgSendPacket struct {
	SourceClient     string
	TimeoutTimestamp uint64
	Payload          IICS26RouterMsgsPayload
}

// IICS26RouterMsgsMsgTimeoutPacket is an auto generated low-level Go binding around an user-defined struct.
type IICS26RouterMsgsMsgTimeoutPacket struct {
	Packet           IICS26RouterMsgsPacket
	NonMembershipMsg []byte
}

// IICS26RouterMsgsPacket is an auto generated low-level Go binding around an user-defined struct.
type IICS26RouterMsgsPacket struct {
	Sequence         uint64
	SourceClient     string
	DestClient       string
	TimeoutTimestamp uint64
	Payloads         []IICS26RouterMsgsPayload
}

// IICS26RouterMsgsPayload is an auto generated low-level Go binding around an user-defined struct.
type IICS26RouterMsgsPayload struct {
	SourcePort string
	DestPort   string
	Version    string
	Encoding   string
	Value      []byte
}

// ContractICS26RouterMetaData contains all meta data concerning the ContractICS26Router contract.
var ContractICS26RouterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ackPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgAckPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgement\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"membershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addClient\",\"inputs\":[{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addIBCApp\",\"inputs\":[{\"name\":\"app\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addIBCApp\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"app\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractILightClient\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitment\",\"inputs\":[{\"name\":\"hashedPath\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCounterparty\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIBCApp\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIIBCApp\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNextClientSeq\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializeV2\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isConsumingScheduledOp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recvPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgRecvPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"membershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgSendPacket\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAuthority\",\"inputs\":[{\"name\":\"newAuthority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitMisbehaviour\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"timeoutPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgTimeoutPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonMembershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AckPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgement\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCAppAdded\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"app\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCAppRecvPacketCallbackError\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientAdded\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientMigrated\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientUpdated\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"result\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02MisbehaviourSubmitted\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Noop\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SendPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TimeoutPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WriteAcknowledgement\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgements\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessManagedInvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessManagedRequiredDelay\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"AccessManagedUnauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"DefaultAdminRoleCannotBeGranted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCAppNotFound\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCAsyncAcknowledgementNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCClientAlreadyExists\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientNotFound\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCCounterpartyClientNotFound\",\"inputs\":[{\"name\":\"counterpartyClientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCErrorUniversalAcknowledgement\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCFailedCallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCInvalidClientId\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidCounterparty\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidPortIdentifier\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidTimeoutDuration\",\"inputs\":[{\"name\":\"maxTimeoutDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualTimeoutDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidTimeoutTimestamp\",\"inputs\":[{\"name\":\"timeoutTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"comparedTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCMultiPayloadPacketNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCPacketAcknowledgementAlreadyExists\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"IBCPacketCommitmentAlreadyExists\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"IBCPacketCommitmentMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IBCPacketReceiptMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IBCPortAlreadyExists\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCUnauthorizedSender\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMerklePrefix\",\"inputs\":[{\"name\":\"prefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"NoAcknowledgements\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"Unreachable\",\"inputs\":[]}]",
	Bin: "0x60a080604052346100c257306080525f5160206163535f395f51905f525460ff8160401c166100b3576002600160401b03196001600160401b03821601610060575b60405161628c90816100c78239608051818181611258015261131b0152f35b6001600160401b0319166001600160401b039081175f5160206163535f395f51905f525581527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f80610041565b63f92ee8a960e01b5f5260045ffd5b5f80fdfe60806040526004361015610011575f80fd5b5f5f3560e01c80631ec43e2314611edf578063223e357a14611eb75780632447af2914611e8157806327f146f314611e4557806329b6eca914611af55780634b720d5b146119a75780634d6e7ce3146115b35780634f1ef286146112d057806352d1902d14611231578063596e00b9146112095780635f516889146111385780636fbf80791461100a5780637795820c14610fc15780637a9e5e4b14610eea5780637eb7893214610e905780638fb3603714610dfd5780639e2e5c8314610cec578063ac9650d814610b8f578063ad3cb1cc14610b2e578063b0777bfa14610abb578063bf7e214f14610a68578063c4d66de814610800578063cce0b2651461042b578063e3cb36a01461017e5763fdbd955d1461012d575f80fd5b3461017b5761015561013e3661207d565b61014661546e565b6101503633614beb565b61479e565b807f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d80f35b80fd5b503461017b57604060031936011261017b576004359067ffffffffffffffff821161017b576040600319833603011261017b576101b9611f8b565b916101c2614e70565b927f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960254915f1983146103fe57600183017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a449602558383807a184f03e93ff9f4daa797ed6e38ed64bf6a1f0100000000000000008110156103d3575b50806d04ee2d6d415b85acef8100000000600a9210156103b8575b662386f26fc100008110156103a4575b6305f5e100811015610393575b612710811015610384575b6064811015610376575b101561036e575b6001810193600a5f1960216102bc6102a689612170565b986102b46040519a8b61214d565b808a52612170565b94601f1960208a0196013687378801015b01917f30313233343536373839616263646566000000000000000000000000000000008282061a83530490811561030957600a905f19906102cd565b505061036a956020958661034d936103569760405199858b9651918291018588015e85019083820190858252519283915e010190815203601f19810186528561214d565b60040183615079565b604051918291602083526020830190612025565b0390f35b60010161028f565b606460029104920191610288565b6127106004910492019161027e565b6305f5e10060089104920191610273565b662386f26fc1000060109104920191610266565b6d04ee2d6d415b85acef810000000060209104920191610256565b604092507a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000090049050600a61023b565b6024847f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b503461017b5761043a36611fae565b9291906104473633614beb565b6104518284614398565b5061045c8284613794565b61046682806123b1565b9067ffffffffffffffff821161076d5761048a82610484855461442c565b856146c7565b8790601f831160011461079a5791806104bb92600195948b92610652575b50505f198260011b9260031b1c19161790565b81555b016104cc602083018361235d565b9168010000000000000000831161076d5780548382558084106106f3575b508752602087208791805b8484106105e05750505050506105d473ffffffffffffffffffffffffffffffffffffffff7f23c2e29d6ae84e79fa116b8afd6e28ddc1de7f473d3edb407fbd08093c3ed6bf951691604051848682376020818681017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960081520301902073ffffffffffffffffffffffffffffffffffffffff84167fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790556105c66040519586956060875260608701916122c1565b90848203602086015261470c565b9060408301520390a180f35b6105ea81836123b1565b9067ffffffffffffffff82116106c65761060e82610608875461442c565b876146c7565b8b908c601f841160011461065d57836001959294602094879661064394926106525750505f198260011b9260031b1c19161790565b86555b019301930192916104f5565b013590505f806104a8565b91601f19841687845260208420935b8181106106ae5750936020936001969387969383889510610695575b505050811b018655610646565b5f1960f88560031b161c199101351690555f8080610688565b9193602060018192878701358155019501920161066c565b60248c7f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b8189528360208a2091820191015b81811061070e57506104ea565b808a61071c6001935461442c565b8061072a575b505001610701565b601f811184146107415750508a81555b8a5f610722565b83601f6020848661075c965220920160051c820191016146b1565b808b528a602081208183555561073a565b6024887f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b8389526020892091601f1984168a5b8181106107e857509160019594929183879593106107cf575b505050811b0181556104be565b5f1960f88560031b161c199101351690555f80806107c2565b919360206001819287870135815501950192016107a9565b503461017b57602060031936011261017b5761081a611f68565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005467ffffffffffffffff81169081610a405760401c60ff16908115610a34575b50610a0c57610978906108d260027fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000007ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005416177ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0055565b680100000000000000007fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005416177ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005561094b615ef9565b610953615ef9565b61095b615ef9565b610963615ef9565b61096b615ef9565b610973615ef9565b615c66565b7fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054167ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a180f35b6004827ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b6002915010155f61085b565b6004847ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b503461017b578060031936011261017b57602073ffffffffffffffffffffffffffffffffffffffff7ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a005416604051908152f35b503461017b57602060031936011261017b576004359067ffffffffffffffff821161017b5761036a610af9610af33660048601611f3a565b9061447d565b604051918291602083526020610b1a82516040838701526060860190612025565b910151601f19848303016040850152612269565b503461017b578060031936011261017b575061036a604051610b5160408261214d565b600581527f352e302e300000000000000000000000000000000000000000000000000000006020820152604051918291602083526020830190612025565b503461017b57602060031936011261017b576004359067ffffffffffffffff821161017b573660238301121561017b5781600401359067ffffffffffffffff821161017b57602483013660248460051b86010111610ce8576040516020610bf6818361214d565b83825280820192601f198201368537610c0e866124c4565b96610c1c604051988961214d565b868852601f19610c2b886124c4565b0183875b828110610cd857505050855b87811015610cc557600190610ca988808989610c958a610c6360248960051b8c01018c6123b1565b9190946040519483869484860198893784019083820190898252519283915e010185815203601f19810183528261214d565b5190305af4610ca2613d25565b9030616013565b610cb3828c6138c4565b52610cbe818b6138c4565b5001610c3b565b6040518481528061036a8187018c612269565b606082828d010152018490610c2f565b5080fd5b5034610df957610cfb366121e0565b73ffffffffffffffffffffffffffffffffffffffff610d1d8486959796614398565b1691823b15610df957610d6a925f92836040518096819582947fddba65370000000000000000000000000000000000000000000000000000000084526020600485015260248401916122c1565b03925af18015610dee57610dba575b507fa263f0a976b2937a51fd2e416491cf0ca724d5499fa870715929dfde4ee4a4309192610db46040519283926020845260208401916122c1565b0390a180f35b7fa263f0a976b2937a51fd2e416491cf0ca724d5499fa870715929dfde4ee4a43092505f610de79161214d565b5f91610d79565b6040513d5f823e3d90fd5b5f80fd5b34610df9575f600319360112610df9577ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a005460a01c60ff1615610e885760207f8fb36037000000000000000000000000000000000000000000000000000000005b7fffffffff0000000000000000000000000000000000000000000000000000000060405191168152f35b60205f610e5e565b34610df9576020600319360112610df95760043567ffffffffffffffff8111610df957610ecc610ec66020923690600401611f3a565b90614398565b73ffffffffffffffffffffffffffffffffffffffff60405191168152f35b34610df9576020600319360112610df957610f03611f68565b73ffffffffffffffffffffffffffffffffffffffff7ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a0054163303610f9557803b15610f5357610f5190615c66565b005b73ffffffffffffffffffffffffffffffffffffffff907fc2f31e5e000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b7f068ca9d8000000000000000000000000000000000000000000000000000000005f523360045260245ffd5b34610df9576020600319360112610df9576004355f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600602052602060405f2054604051908152f35b34610df957602061108961101d366121e0565b61102b959293953633614beb565b73ffffffffffffffffffffffffffffffffffffffff61104a8588614398565b16905f6040518097819582947f0bece35600000000000000000000000000000000000000000000000000000000845288600485015260248401916122c1565b03925af1918215610dee575f926110fa575b506020927f87bbef2779889a19f0435ddca81fda94132c06ffddb0ea73def256307a293aef916110d86040519283926040845260408401916122c1565b6110e185612232565b84868301520390a1604051906110f681612232565b8152f35b9091506020813d602011611130575b816111166020938361214d565b81010312610df957516003811015610df95790602061109b565b3d9150611109565b34610df9576040600319360112610df95760043567ffffffffffffffff8111610df9576111df61116f6111e4923690600401611f3a565b9190611179611f8b565b9261118261546e565b61118c3633614beb565b611199818381151561434f565b6111bb81836111b46111ac36848461218c565b805190615f50565b501561434f565b6111d881836111d36111ce36848461218c565b614ebc565b61434f565b369161218c565b6156fd565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d005b34610df9576111e461121a3661204a565b61122261546e565b61122c3633614beb565b613d54565b34610df9575f600319360112610df95773ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001630036112a85760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b7fe07c8dba000000000000000000000000000000000000000000000000000000005f5260045ffd5b6040600319360112610df9576112e4611f68565b60243567ffffffffffffffff8111610df9576113049036906004016121c2565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016803014908115611571575b506112a8576113553633614beb565b73ffffffffffffffffffffffffffffffffffffffff8216916040517f52d1902d000000000000000000000000000000000000000000000000000000008152602081600481875afa5f918161153d575b506113d557837f4c9c8ce3000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8592036115125750813b156114e757807fffffffffffffffffffffffff00000000000000000000000000000000000000007f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5416177f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a28151156114b6575f80836020610f5195519101845af46114b0613d25565b91616013565b5050346114bf57005b7fb398979f000000000000000000000000000000000000000000000000000000005f5260045ffd5b7f4c9c8ce3000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b7faa1d49a4000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b9091506020813d602011611569575b816115596020938361214d565b81010312610df9575190856113a4565b3d915061154c565b905073ffffffffffffffffffffffffffffffffffffffff7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5416141583611346565b34610df9576115c13661207d565b6115c961546e565b604081019061160c73ffffffffffffffffffffffffffffffffffffffff6116026115fc6115f6868661232a565b806123b1565b906137cc565b1633903314613860565b611619610af382806123b1565b5190602081019061164a61162c83612464565b429061163785612464565b9067ffffffffffffffff42911611612da9565b6201518061166a4267ffffffffffffffff61166486612464565b166138aa565b116116814267ffffffffffffffff61166486612464565b906119755750602061169382806123b1565b919082604051938492833781017f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e07476018152030190209267ffffffffffffffff8454169467ffffffffffffffff86146119485767ffffffffffffffff600161173097011694857fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000082541617905561172883806123b1565b969094612464565b604094855197611740878a61214d565b60018952601f1987015f5b8181106119125750509267ffffffffffffffff8899936117816117a9946117b8978b519d8e611779816120b0565b52369161218c565b9660208c01978852898c01521660608a01528260808a01526117a436918761232a565b612deb565b6117b2826138b7565b526138b7565b506117d0815167ffffffffffffffff87511690615860565b6020815191012090815f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205261181b845f205415915167ffffffffffffffff88511690615860565b90156118d15750937fab3a4458a269be61dfa43faa33aa7b1f5d570716f83ad078bc2ba5dab039abae6118a56118888694602098611858866158f1565b905f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e07476008a52875f2055806123b1565b90818751928392833781015f8152039020928551918291826138d8565b0390a35f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d51908152f35b61190e9084519182917f91ffd924000000000000000000000000000000000000000000000000000000008352602060048401526024830190612025565b0390fd5b8089602080938e6060845194611927866120b0565b8186528185870152850152606080850152606060808501520101520161174b565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b7f715fed60000000000000000000000000000000000000000000000000000000005f526201518060045260245260445ffd5b34610df9576020600319360112610df9576119c0611f68565b6119c861546e565b73ffffffffffffffffffffffffffffffffffffffff811690816119eb602a612170565b906119f9604051928361214d565b602a8252611a07602a612170565b601f19602084019101368237825115611ac85760309053815160011015611ac8576078602183015360295b60018111611a7a5750611a49576111e492506156fd565b827fe22e27eb000000000000000000000000000000000000000000000000000000005f52600452601460245260445ffd5b90600f81166010811015611ac8577f3031323334353637383961626364656600000000000000000000000000000000901a611ab58385614eab565b5360041c908015611948575f1901611a32565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b34610df9576020600319360112610df957611b0e611f68565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005467ffffffffffffffff81169060018203611e115760401c60ff16908115611e39575b50611e1157611d3990611bc960027fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000007ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005416177ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0055565b680100000000000000007fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005416177ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00557fffffffffffffffffffffffff00000000000000000000000000000000000000007fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b005473ffffffffffffffffffffffffffffffffffffffff811633148015611dcc575b611ca8903390613860565b167fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b00557fffffffffffffffffffffffff00000000000000000000000000000000000000007fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b0154167fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b0155610963615ef9565b7fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054167ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a1005b50611ca873ffffffffffffffffffffffffffffffffffffffff7fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b01541633149050611c9d565b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f5260045ffd5b60029150101582611b52565b34610df9575f600319360112610df95760207f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960254604051908152f35b34610df9576020600319360112610df95760043567ffffffffffffffff8111610df957610ecc6115fc6020923690600401611f3a565b34610df9576111e4611ec83661204a565b611ed061546e565b611eda3633614beb565b6131d3565b34610df95761036a6103566111d8611ef636611fae565b90611f05949293943633614beb565b611f1284868115156122e1565b611f2a8486611f256111ce36848461218c565b6122e1565b611f3536858761218c565b615079565b9181601f84011215610df95782359167ffffffffffffffff8311610df95760208381860195010111610df957565b6004359073ffffffffffffffffffffffffffffffffffffffff82168203610df957565b6024359073ffffffffffffffffffffffffffffffffffffffff82168203610df957565b6060600319820112610df95760043567ffffffffffffffff8111610df95781611fd991600401611f3a565b929092916024359067ffffffffffffffff8211610df95760031982604092030112610df9576004019060443573ffffffffffffffffffffffffffffffffffffffff81168103610df95790565b90601f19601f602080948051918291828752018686015e5f8582860101520116010190565b6020600319820112610df9576004359067ffffffffffffffff8211610df95760031982604092030112610df95760040190565b6020600319820112610df9576004359067ffffffffffffffff8211610df95760031982606092030112610df95760040190565b60a0810190811067ffffffffffffffff8211176120cc57604052565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6040810190811067ffffffffffffffff8211176120cc57604052565b6080810190811067ffffffffffffffff8211176120cc57604052565b6060810190811067ffffffffffffffff8211176120cc57604052565b90601f601f19910116810190811067ffffffffffffffff8211176120cc57604052565b67ffffffffffffffff81116120cc57601f01601f191660200190565b92919261219882612170565b916121a6604051938461214d565b829481845281830111610df9578281602093845f960137010152565b9080601f83011215610df9578160206121dd9335910161218c565b90565b6040600319820112610df95760043567ffffffffffffffff8111610df9578161220b91600401611f3a565b929092916024359067ffffffffffffffff8211610df95761222e91600401611f3a565b9091565b6003111561223c57565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b9080602083519182815201916020808360051b8301019401925f915b83831061229457505050505090565b90919293946020806122b283601f1986600196030187528951612025565b97019301930191939290612285565b601f8260209493601f1993818652868601375f8582860101520116010190565b919091156122ed575050565b61190e6040519283927f4870bd740000000000000000000000000000000000000000000000000000000084526020600485015260248401916122c1565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6181360301821215610df9570190565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610df9570180359067ffffffffffffffff8211610df957602001918160051b36038313610df957565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610df9570180359067ffffffffffffffff8211610df957602001918136038313610df957565b9290921561240f57505050565b612452929161190e916040519485947f9fff831f000000000000000000000000000000000000000000000000000000008652604060048701526044860190612025565b916003198584030160248601526122c1565b3567ffffffffffffffff81168103610df95790565b359067ffffffffffffffff82168203610df957565b9190826040910312610df9576040516124a6816120f9565b60206124bf8183956124b781612479565b855201612479565b910152565b67ffffffffffffffff81116120cc5760051b60200190565b9080601f83011215610df95781356124f3816124c4565b92612501604051948561214d565b81845260208085019260051b82010191838311610df95760208201905b83821061252d57505050505090565b813567ffffffffffffffff8111610df957602091612550878480948801016121c2565b81520191019061251e565b9080601f83011215610df957813591612573836124c4565b92612581604051948561214d565b80845260208085019160051b83010191838311610df95760208101915b8383106125ad57505050505090565b823567ffffffffffffffff8111610df9578201906040601f198388030112610df957604051906125dc826120f9565b602083013567ffffffffffffffff8111610df9578760206125ff928601016124dc565b825260408301359167ffffffffffffffff8311610df957612628886020809695819601016121c2565b8382015281520192019161259e565b9190608083820312610df95760405161264f81612115565b8093803567ffffffffffffffff8111610df9578361266e9183016121c2565b8252602081013567ffffffffffffffff8111610df957836126909183016121c2565b6020830152604081013567ffffffffffffffff8111610df9578101608081850312610df957604051906126c282612115565b80356003811015610df957825260208101356003811015610df957602083015260408101356003811015610df957604083015260608101359067ffffffffffffffff8211610df957612716918691016121c2565b6060820152604083015260608101359067ffffffffffffffff8211610df957019180601f84011215610df95782359261274e846124c4565b9361275c604051958661214d565b80855260208086019160051b83010191838311610df95760208101915b83831061278b57505050505060600152565b823567ffffffffffffffff8111610df9578201906060601f198388030112610df957604051906127ba82612131565b60208301356003811015610df9578252604083013567ffffffffffffffff8111610df9578760206127ed928601016121c2565b602083015260608301359167ffffffffffffffff8311610df957612819886020809695819601016121c2565b6040820152815201920191612779565b35908115158203610df957565b9080601f83011215610df95781359161284e836124c4565b9261285c604051948561214d565b80845260208085019160051b83010191838311610df95760208101915b83831061288857505050505090565b823567ffffffffffffffff8111610df95782016020601f198288030112610df957604051906020820182811067ffffffffffffffff8211176120cc57604052602081013567ffffffffffffffff8111610df957602091010186601f82011215610df9578035906128f7826124c4565b91612905604051938461214d565b80835260208084019160051b83010191898311610df95760208101915b83831061293e5750505090825250815260209283019201612879565b82359067ffffffffffffffff8211610df9576060601f198d9385018094030112610df9576040519061296f82612131565b60208301356002811015610df9578252604083013567ffffffffffffffff8111610df9578d60206129a292860101612637565b6020830152606083013567ffffffffffffffff8111610df957602060a0918f95010180940312610df957604051916129d9836120b0565b833567ffffffffffffffff8111610df9578e6129f69186016121c2565b8352612a0460208501612829565b6020840152604084013567ffffffffffffffff8111610df9578e612a29918601612637565b6040840152612a3a60608501612829565b606084015260808401359267ffffffffffffffff8411610df957612a648f60209695879601612637565b60808201526040820152815201920191612922565b9190826060910312610df957604051612a9181612131565b809280356fffffffffffffffffffffffffffffffff81168103610df95760409182918452602081013560208501520135910152565b9080602083519182815201916020808360051b8301019401925f915b838310612af157505050505090565b9091929394602080612b2f83601f19866001960301875289519083612b1f8351604084526040840190612269565b9201519084818403910152612025565b97019301930191939290612ae2565b6002111561223c57565b6060612bc7612b75612b638451608087526080870190612025565b60208501518682036020880152612025565b608083604086015187840360408901528051612b9081612232565b84526020810151612ba081612232565b60208501526040810151612bb381612232565b604085015201519181858201520190612025565b910151916060818303910152815180825260208201916020808360051b8301019401925f915b838310612bfc57505050505090565b9091929394602080612c4d83601f1986600196030187528951908151612c2181612232565b81526040612c3c858401516060878501526060840190612025565b920151906040818403910152612025565b97019301930191939290612bed565b9080602083519182815201906020808260051b8501019401915f905b828210612c8757505050505090565b909192939594601f19878203018252845190602081019151916020825282518091526040820190602060408260051b8501019401925f5b828110612ce257505050505060208060019296019201920190929195939495612c78565b9091929394602080612d9c837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0876001960301895289518051612d2481612b3e565b82526040612d3f858301516060878601526060850190612b48565b9101519160408183039101526080612d7f612d63845160a0855260a0850190612025565b8685015115158785015260408501518482036040860152612b48565b926060810151151560608401520151906080818403910152612b48565b9701950193929101612cbe565b15612db2575050565b67ffffffffffffffff907f65d30129000000000000000000000000000000000000000000000000000000005f521660045260245260445ffd5b919060a083820312610df95760405190612e04826120b0565b8193803567ffffffffffffffff8111610df95782612e239183016121c2565b8352602081013567ffffffffffffffff8111610df95782612e459183016121c2565b6020840152604081013567ffffffffffffffff8111610df95782612e6a9183016121c2565b6040840152606081013567ffffffffffffffff8111610df95782612e8f9183016121c2565b606084015260808101359167ffffffffffffffff8311610df9576080926124bf92016121c2565b6121dd916080612f0e612efc612eea612ed8865160a0875260a0870190612025565b60208701518682036020880152612025565b60408601518582036040870152612025565b60608501518482036060860152612025565b920151906080818403910152612025565b90608073ffffffffffffffffffffffffffffffffffffffff81612f89612f63612f51875160a0885260a0880190612025565b60208801518782036020890152612025565b67ffffffffffffffff604088015116604087015260608701518682036060880152612eb6565b9401511691015290565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811215610df957016020813591019167ffffffffffffffff8211610df9578136038313610df957565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811215610df957016020813591019167ffffffffffffffff8211610df9578160051b36038313610df957565b9067ffffffffffffffff61304983612479565b1681526130b461308e6130736130626020860186612f93565b60a0602087015260a08601916122c1565b6130806040860186612f93565b9085830360408701526122c1565b9267ffffffffffffffff6130a460608301612479565b1660608401526080810190612fe3565b9290916080818303910152828152602081019260208160051b83010193835f917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6182360301945b84841061310c575050505050505090565b90919293949596601f19828203018352873587811215610df95760206131c260019387839401906131b46131a961318e61317361315a61314c8780612f93565b60a0885260a08801916122c1565b61316689880188612f93565b908783038b8901526122c1565b6131806040870187612f93565b9086830360408801526122c1565b61319b6060860186612f93565b9085830360608701526122c1565b926080810190612f93565b9160808185039101526122c1565b9901930194019291959493906130fb565b60016131ec6131e2838061232a565b608081019061235d565b90500361376c576132006131e2828061232a565b15611ac8578061320f9161232a565b9061322a610af3613220838061232a565b60208101906123b1565b6132708151602081519101206132506111d8613246868061232a565b60408101906123b1565b6020815191012014825190613268613246868061232a565b929091612402565b6132a461329f613283613246858061232a565b9190613297613292878061232a565b612464565b92369161218c565b6154e2565b906132b260208401846123b1565b81019190602081840312610df95780359067ffffffffffffffff8211610df957019261014084840312610df9576040519060e0820182811067ffffffffffffffff8211176120cc57604052613307848661248e565b8252604085013567ffffffffffffffff8111610df9578461332991870161255b565b9260208301938452606086013567ffffffffffffffff8111610df95785613351918801612836565b95604084019687526060840192608082013584526133728760a08401612a79565b9660808601978852610100830135926002841015610df95760a087019384526101208101359067ffffffffffffffff8211610df95701906133b2916124dc565b9260c0860193845260200151906133c89161555c565b82526133d4888061232a565b602081016133e1916123b1565b6133ea91614398565b73ffffffffffffffffffffffffffffffffffffffff169560405197889687967fa6f031bb00000000000000000000000000000000000000000000000000000000885260048801602090526024880190519061345a9167ffffffffffffffff60208092828151168552015116910152565b51606487016101409052610164870161347291612ac6565b9051908681037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc0160848801526134a891612c5c565b935160a48601525180516fffffffffffffffffffffffffffffffff1660c4860152602081015160e486015260400151610104850152516134e781612b3e565b61012484015251908281037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc0161014484015261352391612269565b03815a6020945f91f18015610dee575f90613738575b613575915067ffffffffffffffff61355c6060613556868061232a565b01612464565b1681101561356f6060613556868061232a565b90612da9565b613587613582828061232a565b61561a565b15613710578161361573ffffffffffffffffffffffffffffffffffffffff6135bb6115fc8467ffffffffffffffff976123b1565b16916135ca613220858061232a565b95906136036135dc613246888061232a565b6135fa6135ec6132928b8061232a565b946040519b6111d88d6120b0565b8a52369161218c565b60208801521660408601523690612deb565b6060840152336080840152803b15610df95761366c5f939184926040519586809481937f5e32b6b6000000000000000000000000000000000000000000000000000000008352602060048401526024830190612f1f565b03925af1918215610dee5767ffffffffffffffff92613700575b507f01e5ed58494819ef3f6480dd08e433b7c08ed75c7abdf2c22c6f04b71340a1686136b5613220838061232a565b6136cf6136c8613292868098959861232a565b948061232a565b9481604051928392833781015f8152039020926136fb6040519283926020845216956020830190613036565b0390a3565b5f61370a9161214d565b5f613686565b50507fd08bf58b0e4eec5bfc697a4fdbb6839057fbf4dd06f1b1ce07445c0e5a654caf5f80a1565b506020813d602011613764575b816137526020938361214d565b81010312610df9576135759051613539565b3d9150613745565b7f356f4dbd000000000000000000000000000000000000000000000000000000005f5260045ffd5b60209082604051938492833781017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960181520301902090565b73ffffffffffffffffffffffffffffffffffffffff604051838382376020818581017fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf600815203019020541691821561382357505090565b61190e6040519283927fa09dbf590000000000000000000000000000000000000000000000000000000084526020600485015260248401916122c1565b156138685750565b73ffffffffffffffffffffffffffffffffffffffff907fbe2f2b45000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b9190820391821161194857565b805115611ac85760200190565b8051821015611ac85760209160051b010190565b906020825267ffffffffffffffff8151166020830152608061392261390c602084015160a0604087015260c0860190612025565b6040840151601f19868303016060870152612025565b9167ffffffffffffffff6060820151168285015201519160a0601f1982840301910152815180825260208201916020808360051b8301019401925f915b83831061396e57505050505090565b909192939460208061398c83601f1986600196030187528951612eb6565b9701930193019193929061395f565b919060a083820312610df9576040516139b3816120b0565b80936139be81612479565b8252602081013567ffffffffffffffff8111610df957836139e09183016121c2565b6020830152604081013567ffffffffffffffff8111610df95783613a059183016121c2565b6040830152613a1660608201612479565b606083015260808101359067ffffffffffffffff8211610df957019180601f84011215610df9578235613a48816124c4565b93613a56604051958661214d565b81855260208086019260051b82010191838311610df95760208201905b838210613a8557505050505060800152565b813567ffffffffffffffff8111610df957602091613aa887848094880101612deb565b815201910190613a73565b602081830312610df95780359067ffffffffffffffff8211610df9570161016081830312610df95760405191610100830183811067ffffffffffffffff8211176120cc57604052613b04818361248e565b8352604082013567ffffffffffffffff8111610df95781613b2691840161255b565b6020840152606082013567ffffffffffffffff8111610df95781613b4b918401612836565b604084015260808201356060840152613b678160a08401612a79565b60808401526101008201356002811015610df95760a084015261012082013567ffffffffffffffff8111610df95781613ba19184016124dc565b60c084015261014082013567ffffffffffffffff8111610df957613bc592016121c2565b60e082015290565b906121dd9160208152613bfa60208201835167ffffffffffffffff60208092828151168552015116910152565b60e0613c9a613c33613c1d60208601516101606060870152610180860190612ac6565b6040860151601f19868303016080870152612c5c565b606085015160a0850152608085015180516fffffffffffffffffffffffffffffffff1660c0860152602081015160e08601526040015161010085015260a0850151613c7d81612b3e565b61012085015260c0850151601f1985830301610140860152612269565b92015190610160601f1982850301910152612025565b60408051909190613cc1838261214d565b6001815291601f1901825f5b828110613cd957505050565b806060602080938501015201613ccd565b60405190613cf960408361214d565b602082527f4774d4a575993f963b1c06573736617a457abef8589178db8d10c94b4ab511ab6020830152565b3d15613d4f573d90613d3682612170565b91613d44604051938461214d565b82523d5f602084013e565b606090565b6001613d636131e2838061232a565b90500361376c57613d776131e2828061232a565b15611ac85780613d869161232a565b905f6020613ee6613d9d610af3613246868061232a565b613dd88151848151910120613dc16111d8613db8898061232a565b878101906123b1565b85815191012014825190613268613db8898061232a565b613dfb613dea6060613556888061232a565b429061163760606135568a8061232a565b613e2b613e26613e17613e0e888061232a565b868101906123b1565b91906132976132928a8061232a565b615860565b613e6a613e49613e4436613e3f8a8061232a565b61399b565b6158f1565b9185613e62613e5a828b018b6123b1565b810190613ab3565b94015161555c565b60c08301526040519084820152838152613e8560408261214d565b60e082015273ffffffffffffffffffffffffffffffffffffffff613eaf610ec6613246888061232a565b16906040519485809481937f974a74c400000000000000000000000000000000000000000000000000000000835260048301613bcd565b03925af18015610dee57614320575b50613f08613f03828061232a565b615b84565b15613710575f80613ff1613f1a613cb0565b9467ffffffffffffffff613fa973ffffffffffffffffffffffffffffffffffffffff613f4c6115fc60208601866123b1565b1692613f5b613220898061232a565b9390613f978a613292613f8e613f80613f77613246858061232a565b9390948061232a565b94604051996111d88b6120b0565b8852369161218c565b60208601521660408401523690612deb565b60608201523360808201526040519485809481937f078c4a79000000000000000000000000000000000000000000000000000000008352602060048401526024830190612f1f565b03925af15f91816142a4575b50614219575061400b613d25565b8051156141f15761404b7fb9edb487876e8be10f54e377c1a815a54ad92a6db1c9561dfe8fad2f0d1da84f91604051918291602083526020830190612025565b0390a1614056613cea565b61405f836138b7565b52614069826138b7565b505b614075818061232a565b604081016140e961329761409d6140a261409d61409286886123b1565b919061329789612464565b615d11565b6020815191012094855f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e07476006020526140e160405f20541595826123b1565b939091612464565b90156141b357506140f983615d8b565b905f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f20557f76765590e2b799b0506100f8a6610cfecab2c71e8e1f8aa981b099aff0dfdb746141a36136fb614159613246858061232a565b61417361416c613292888096959661232a565b968061232a565b9281604051928392833781015f81520390209467ffffffffffffffff604051948594604086526040860190613036565b9184830360208601521696612269565b61190e906040519182917f40470d74000000000000000000000000000000000000000000000000000000008352602060048401526024830190612025565b7fadef7fb8000000000000000000000000000000000000000000000000000000005f5260045ffd5b80511561427c578051602082012061422f613cea565b602081519101201461425457614244836138b7565b5261424e826138b7565b5061406b565b7f6b2675e3000000000000000000000000000000000000000000000000000000005f5260045ffd5b7fecfef798000000000000000000000000000000000000000000000000000000005f5260045ffd5b9091503d805f833e6142b6818361214d565b810190602081830312610df95780519067ffffffffffffffff8211610df9570181601f82011215610df9578051906142ed82612170565b926142fb604051948561214d565b82845260208383010111610df957815f9260208093018386015e83010152905f613ffd565b6020813d602011614347575b816143396020938361214d565b81010312610df95751613ef5565b3d915061432c565b9190911561435b575050565b61190e6040519283927f14d712470000000000000000000000000000000000000000000000000000000084526020600485015260248401916122c1565b73ffffffffffffffffffffffffffffffffffffffff604051838382376020818581017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960081520301902054169182156143ef57505090565b61190e6040519283927fa0db16fe0000000000000000000000000000000000000000000000000000000084526020600485015260248401916122c1565b90600182811c92168015614473575b602083101461444657565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f169161443b565b6060602060405161448d816120f9565b828152015261449c8282613794565b91604051926144aa846120f9565b6040515f82546144b98161442c565b808452906001811690811561466f575060011461462c575b50906144e28160019493038261214d565b8552018054906144f1826124c4565b916144ff604051938461214d565b80835260208301915f5260205f20915f905b82821061456b575050505060208401528251511561452e57505090565b61190e6040519283927fdf95155a0000000000000000000000000000000000000000000000000000000084526020600485015260248401916122c1565b6040515f855461457a8161442c565b80845290600181169081156145eb57506001146145b4575b50600192826145a68594602094038261214d565b815201940191019092614511565b5f878152602081209092505b8183106145d557505081016020016001614592565b60018160209254838688010152019201916145c0565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660208581019190915291151560051b8401909101915060019050614592565b5f8481526020812094939250905b80821061465357509192509081016020016144e26144d1565b919293600181602092548385880101520191019093929161463a565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660208086019190915291151560051b840190910191506144e290506144d1565b8181106146bc575050565b5f81556001016146b1565b9190601f81116146d657505050565b614700925f5260205f20906020601f840160051c83019310614702575b601f0160051c01906146b1565b565b90915081906146f3565b9061473661472b61471d8480612f93565b6040855260408501916122c1565b926020810190612fe3565b90916020818503910152808352602083019260208260051b82010193835f925b8484106147665750505050505090565b90919293949560208061478e83601f1986600196030188526147888b88612f93565b906122c1565b9801940194019294939190614756565b60016147ad6131e2838061232a565b90500361376c576147c16131e2828061232a565b15611ac857806147d09161232a565b906147e1610af3613220838061232a565b916148168351602081519101206147fe6111d8613246868061232a565b6020815191012014845190613268613246868061232a565b5f60206148c061482f61409d613e17613246888061232a565b95614838613cb0565b9661487b614869858901996148506111d88c8c6123b1565b614859826138b7565b52614863816138b7565b50615d8b565b9185613e62613e5a60408c018c6123b1565b60c0830152604051908482015283815261489660408261214d565b60e082015273ffffffffffffffffffffffffffffffffffffffff613eaf610ec6613e0e898061232a565b03925af18015610dee57614bbc575b506148dd613582838061232a565b15614b935773ffffffffffffffffffffffffffffffffffffffff6149046115fc83806123b1565b16614912613220848061232a565b614922613246868096949661232a565b614932613292888095949561232a565b9161493d89896123b1565b9790916040519560c087019487861067ffffffffffffffff8711176120cc57613f8e61499594614975946149a498604052369161218c565b956020860196875267ffffffffffffffff60408701951685523690612deb565b9660608501978852369161218c565b6080830190815260a0830191338352853b15610df95760405196879586957f428e4e170000000000000000000000000000000000000000000000000000000087526004870160209052516024870160c0905260e48701614a0391612025565b9051908681037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc016044880152614a3991612025565b915167ffffffffffffffff16606486015251908481037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc016084860152614a7f91612eb6565b9051908381037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc0160a4850152614ab591612025565b905173ffffffffffffffffffffffffffffffffffffffff1660c483015203815a5f948591f18015610dee576136fb927ff9bab74bcdb634f4d3dd064cc42a13df056598e1c0336905d2f5750fbfb08b7b92614b7392614b83575b50614b1d613220828061232a565b949095614b41614b30613292858061232a565b91614b3b858061232a565b946123b1565b96909781604051928392833781015f81520390209567ffffffffffffffff604051958695604087526040870190613036565b92858403602087015216976122c1565b5f614b8d9161214d565b5f614b0f565b5050507fd08bf58b0e4eec5bfc697a4fdbb6839057fbf4dd06f1b1ce07445c0e5a654caf5f80a1565b6020813d602011614be3575b81614bd56020938361214d565b81010312610df957516148cf565b3d9150614bc8565b7ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a00549173ffffffffffffffffffffffffffffffffffffffff83169281600411610df9575f5f9060405f81519673ffffffffffffffffffffffffffffffffffffffff60208901917fb700961300000000000000000000000000000000000000000000000000000000835216978860248201523060448201527fffffffff00000000000000000000000000000000000000000000000000000000833516606482015260648152614cba60848261214d565b828052826020525190895afa614e5d575b15614cd8575b5050505050565b63ffffffff1615614e31577fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1674010000000000000000000000000000000000000000177ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a0055823b15610df9576020925f92836040518096819582947f94c7d7ee000000000000000000000000000000000000000000000000000000008452600484015260406024840152601f19601f6044850192808452808786860137868582860101520116010103925af18015610dee57614e21575b507fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff7ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a0054167ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a00555f80808080614cd1565b5f614e2b9161214d565b5f614db0565b827f068ca9d8000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b50505f516020518060201c150290614ccb565b60405190614e7f60408361214d565b600782527f636c69656e742d000000000000000000000000000000000000000000000000006020830152565b908151811015611ac8570160200190565b80516004811090811561506e575b5061505257614f10604051614ee060408261214d565b600881527f6368616e6e656c2d000000000000000000000000000000000000000000000000602082015282615e82565b8015615057575b615052575f5b815181101561504b57614f308183614eab565b5160f81c6061811015908161503f575b8115615021575b8115615003575b8115614fc4575b8115614f6f575b50614f675750505f90565b600101614f1d565b6023811491508115614fb9575b8115614fae575b8115614fa3575b8115614f98575b505f614f5c565b603e9150145f614f91565b603c81149150614f8a565b605d81149150614f83565b605b81149150614f7c565b9050602e81148015614ff9575b8015614fef575b8015614fe5575b90614f55565b50602d8114614fdf565b50602b8114614fd8565b50605f8114614fd1565b9050604181101580615016575b90614f4e565b50605a811115615010565b9050603081101580615034575b90614f47565b50603981111561502e565b607a8111159150614f40565b5050600190565b505f90565b50615069615063614e70565b82615e82565b614f17565b60809150115f614eca565b91906040519173ffffffffffffffffffffffffffffffffffffffff845193602081818801968088835e81017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960081520301902054166154335773ffffffffffffffffffffffffffffffffffffffff169160405160208186518085835e81017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960081520301902073ffffffffffffffffffffffffffffffffffffffff84167fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790556020604051809286518091835e81017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960181520301902061519b82806123b1565b9067ffffffffffffffff82116120cc576151b982610484855461442c565b5f90601f83116001146153cc5791806151e992600195945f926106525750505f198260011b9260031b1c19161790565b81555b016151fa602083018361235d565b906801000000000000000082116120cc578254828455808310615355575b505f928352602083209290805b83831061527a575050505050916105c69161526f7f0ecded31ecd211a73abf0fb3bc09150bbe321a05550fbe29ea0f16b6e25fbfa894604051948594606086526060860190612025565b9060408301520390a1565b61528481836123b1565b9067ffffffffffffffff82116120cc576152a8826152a2895461442c565b896146c7565b5f90601f83116001146152eb57926152dc836001959460209487965f926106525750505f198260011b9260031b1c19161790565b88555b01950192019193615225565b601f19831691885f5260205f20925f5b81811061533d5750936020936001969387969383889510615324575b505050811b0188556152df565b5f1960f88560031b161c199101351690555f8080615317565b919360206001819287870135815501950192016152fb565b835f528260205f2091820191015b8181106153705750615218565b8061537d6001925461442c565b8061538a575b5001615363565b601f8111831461539f57505f81555b5f615383565b6153bb90825f5283601f60205f20920160051c820191016146b1565b805f525f6020812081835555615399565b601f19831691845f5260205f20925f5b81811061541b5750916001959492918387959310615402575b505050811b0181556151ec565b5f1960f88560031b161c199101351690555f80806153f5565b919360206001819287870135815501950192016153dc565b6040517f87dfb267000000000000000000000000000000000000000000000000000000008152602060048201528061190e6024820187612025565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005c6154ba5760017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d565b7f3ee5aeb5000000000000000000000000000000000000000000000000000000005f5260045ffd5b60096121dd916020937fffffffffffffffff000000000000000000000000000000000000000000000000856040519687948051918291018387015e840101917f0200000000000000000000000000000000000000000000000000000000000000835260c01b1660018201520301601f19810183528261214d565b908151156155df5781515f1981019081116119485760209182806155836155b994876138c4565b51926040519584879551918291018487015e8401908282015f8152815193849201905e01015f815203601f19810183528261214d565b81515f198101908111611948576155db916155d482856138c4565b52826138c4565b5090565b6040517fa7c34e4f000000000000000000000000000000000000000000000000000000008152602060048201528061190e6024820185612269565b615638613e2661562d60208401846123b1565b919061329785612464565b6020815191012090815f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f205480156156f65761568f613e44615685613e44368661399b565b831493369061399b565b91156156c85750505f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e07476006020525f6040812055600190565b7f3f87a2ec000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b5050505f90565b906040519073ffffffffffffffffffffffffffffffffffffffff835192602081818701958087835e81017fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf6008152030190205416615825579161581a9173ffffffffffffffffffffffffffffffffffffffff7fa6ec8e860960e638347460dc632fbe0175c51a5ca130e336138bbe26ff3044999416906020604051809285518091835e81017fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf60081520301902073ffffffffffffffffffffffffffffffffffffffff82167fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055604051928392604084526040840190612025565b9060208301520390a1565b6040517f837f46a6000000000000000000000000000000000000000000000000000000008152602060048201528061190e6024820186612025565b60096121dd916020937fffffffffffffffff000000000000000000000000000000000000000000000000856040519687948051918291018387015e840101917f0100000000000000000000000000000000000000000000000000000000000000835260c01b1660018201520301601f19810183528261214d565b60209291908391805192839101825e019081520190565b90602091604051615902848261214d565b5f8152905f915b60808201518051841015615a665783615921916138c4565b51855f818351604051918183925191829101835e8101838152039060025afa15610dee575f5190865f8180840151604051918183925191829101835e8101838152039060025afa15610dee575f5191875f816040850151604051918183925191829101835e8101838152039060025afa15610dee575f5192885f816060860151604051918183925191829101835e8101838152039060025afa15610dee57885f8160808251960151604051918183925191829101835e8101838152039060025afa15610dee5788935f938451916040519387850195865260408501526060840152608083015260a082015260a08152615a1b60c08261214d565b604051918291518091835e8101838152039060025afa15610dee57600190615a5e5f5191615a506040519384928a84016158da565b03601f19810183528261214d565b920191615909565b509150929192825f816040840151604051918183925191829101835e8101838152039060025afa15610dee57825f606081519301516040517fffffffffffffffff0000000000000000000000000000000000000000000000008482019260c01b16825260088152615ad860288261214d565b604051918291518091835e8101838152039060025afa15610dee57825f81815194604051918183925191829101835e8101838152039060025afa15610dee575f91825160405191858301937f0200000000000000000000000000000000000000000000000000000000000000855260218401526041830152606182015260618152615b6460818261214d565b604051918291518091835e8101838152039060025afa15610dee575f5190565b615bb4615ba561329f615b9a60408501856123b1565b919061329786612464565b6020815191012091369061399b565b604051615bc981615a506020820194856138d8565b51902090805f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f20548281146156f65780615c3657505f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f2055600190565b90507f657b94fe000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b602073ffffffffffffffffffffffffffffffffffffffff7f2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad9216807fffffffffffffffffffffffff00000000000000000000000000000000000000007ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a005416177ff3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a0055604051908152a1565b60096121dd916020937fffffffffffffffff000000000000000000000000000000000000000000000000856040519687948051918291018387015e840101917f0300000000000000000000000000000000000000000000000000000000000000835260c01b1660018201520301601f19810183528261214d565b90815115615e5a57602091604051615da3848261214d565b5f8152905f915b8151831015615e0157845f81615dc086866138c4565b51604051918183925191829101835e8101838152039060025afa15610dee57600190615df95f5191615a506040519384928a84016158da565b920191615daa565b90505f9150929192604051615b6460218286808201957f020000000000000000000000000000000000000000000000000000000000000087528051918291018484015e810186838201520301601f19810183528261214d565b7f760d6a9b000000000000000000000000000000000000000000000000000000005f5260045ffd5b8051908251808310615ef1578280821091180280831892141582028218906020615eac83856138aa565b9280615ed0615eba86612170565b95615ec8604051978861214d565b808752612170565b95601f19848701970136883703920101835e51902090602081519101201490565b505050505f90565b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c1615615f2857565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f5260045ffd5b80518211801561600c575b615fb5576001821180615fbd575b158015908160011b91820460021417156119485760280180602811611948578203615fb55773ffffffffffffffffffffffffffffffffffffffff92915f615faf926160ac565b90921690565b50505f905f90565b507f30780000000000000000000000000000000000000000000000000000000000007fffff00000000000000000000000000000000000000000000000000000000000060208301511614615f69565b505f615f5b565b90616050575080511561602857602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b815115806160a3575b616061575090565b73ffffffffffffffffffffffffffffffffffffffff907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15616059565b9290926001840180851161194857831180616162575b15938415948560011b9586046002141715611948575f948101809111611948579192905b8183106160f65750505060019190565b9092919360ff61612d7fff0000000000000000000000000000000000000000000000000000000000000060208886010151166161b3565b16600f8111616157578160041b9180830460101490151715611948576001910194019192906160e6565b505f94508493505050565b507f30780000000000000000000000000000000000000000000000000000000000007fffff0000000000000000000000000000000000000000000000000000000000006020868401015116146160c2565b60f81c602f811180616275575b156161ed577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd00160ff1690565b606081118061626b575b15616224577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa90160ff1690565b6040811180616261575b1561625b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc90160ff1690565b5060ff90565b506047811061622e565b50606781106161f7565b50603a81106161c056fea164736f6c634300081c000af0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
}

// ContractICS26RouterABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractICS26RouterMetaData.ABI instead.
var ContractICS26RouterABI = ContractICS26RouterMetaData.ABI

// ContractICS26RouterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractICS26RouterMetaData.Bin instead.
var ContractICS26RouterBin = ContractICS26RouterMetaData.Bin

// DeployContractICS26Router deploys a new Ethereum contract, binding an instance of ContractICS26Router to it.
func DeployContractICS26Router(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractICS26Router, error) {
	parsed, err := ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractICS26RouterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractICS26Router{ContractICS26RouterCaller: ContractICS26RouterCaller{contract: contract}, ContractICS26RouterTransactor: ContractICS26RouterTransactor{contract: contract}, ContractICS26RouterFilterer: ContractICS26RouterFilterer{contract: contract}}, nil
}

// ContractICS26Router is an auto generated Go binding around an Ethereum contract.
type ContractICS26Router struct {
	ContractICS26RouterCaller     // Read-only binding to the contract
	ContractICS26RouterTransactor // Write-only binding to the contract
	ContractICS26RouterFilterer   // Log filterer for contract events
}

// ContractICS26RouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractICS26RouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractICS26RouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractICS26RouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractICS26RouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractICS26RouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractICS26RouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractICS26RouterSession struct {
	Contract     *ContractICS26Router // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ContractICS26RouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractICS26RouterCallerSession struct {
	Contract *ContractICS26RouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ContractICS26RouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractICS26RouterTransactorSession struct {
	Contract     *ContractICS26RouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ContractICS26RouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractICS26RouterRaw struct {
	Contract *ContractICS26Router // Generic contract binding to access the raw methods on
}

// ContractICS26RouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractICS26RouterCallerRaw struct {
	Contract *ContractICS26RouterCaller // Generic read-only contract binding to access the raw methods on
}

// ContractICS26RouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractICS26RouterTransactorRaw struct {
	Contract *ContractICS26RouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractICS26Router creates a new instance of ContractICS26Router, bound to a specific deployed contract.
func NewContractICS26Router(address common.Address, backend bind.ContractBackend) (*ContractICS26Router, error) {
	contract, err := bindContractICS26Router(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractICS26Router{ContractICS26RouterCaller: ContractICS26RouterCaller{contract: contract}, ContractICS26RouterTransactor: ContractICS26RouterTransactor{contract: contract}, ContractICS26RouterFilterer: ContractICS26RouterFilterer{contract: contract}}, nil
}

// NewContractICS26RouterCaller creates a new read-only instance of ContractICS26Router, bound to a specific deployed contract.
func NewContractICS26RouterCaller(address common.Address, caller bind.ContractCaller) (*ContractICS26RouterCaller, error) {
	contract, err := bindContractICS26Router(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterCaller{contract: contract}, nil
}

// NewContractICS26RouterTransactor creates a new write-only instance of ContractICS26Router, bound to a specific deployed contract.
func NewContractICS26RouterTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractICS26RouterTransactor, error) {
	contract, err := bindContractICS26Router(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterTransactor{contract: contract}, nil
}

// NewContractICS26RouterFilterer creates a new log filterer instance of ContractICS26Router, bound to a specific deployed contract.
func NewContractICS26RouterFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractICS26RouterFilterer, error) {
	contract, err := bindContractICS26Router(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterFilterer{contract: contract}, nil
}

// bindContractICS26Router binds a generic wrapper to an already deployed contract.
func bindContractICS26Router(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractICS26Router *ContractICS26RouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractICS26Router.Contract.ContractICS26RouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractICS26Router *ContractICS26RouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.ContractICS26RouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractICS26Router *ContractICS26RouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.ContractICS26RouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractICS26Router *ContractICS26RouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractICS26Router.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractICS26Router *ContractICS26RouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractICS26Router *ContractICS26RouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ContractICS26Router *ContractICS26RouterCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ContractICS26Router *ContractICS26RouterSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ContractICS26Router.Contract.UPGRADEINTERFACEVERSION(&_ContractICS26Router.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ContractICS26Router *ContractICS26RouterCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ContractICS26Router.Contract.UPGRADEINTERFACEVERSION(&_ContractICS26Router.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_ContractICS26Router *ContractICS26RouterCaller) Authority(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "authority")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_ContractICS26Router *ContractICS26RouterSession) Authority() (common.Address, error) {
	return _ContractICS26Router.Contract.Authority(&_ContractICS26Router.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_ContractICS26Router *ContractICS26RouterCallerSession) Authority() (common.Address, error) {
	return _ContractICS26Router.Contract.Authority(&_ContractICS26Router.CallOpts)
}

// GetClient is a free data retrieval call binding the contract method 0x7eb78932.
//
// Solidity: function getClient(string clientId) view returns(address)
func (_ContractICS26Router *ContractICS26RouterCaller) GetClient(opts *bind.CallOpts, clientId string) (common.Address, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "getClient", clientId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetClient is a free data retrieval call binding the contract method 0x7eb78932.
//
// Solidity: function getClient(string clientId) view returns(address)
func (_ContractICS26Router *ContractICS26RouterSession) GetClient(clientId string) (common.Address, error) {
	return _ContractICS26Router.Contract.GetClient(&_ContractICS26Router.CallOpts, clientId)
}

// GetClient is a free data retrieval call binding the contract method 0x7eb78932.
//
// Solidity: function getClient(string clientId) view returns(address)
func (_ContractICS26Router *ContractICS26RouterCallerSession) GetClient(clientId string) (common.Address, error) {
	return _ContractICS26Router.Contract.GetClient(&_ContractICS26Router.CallOpts, clientId)
}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 hashedPath) view returns(bytes32)
func (_ContractICS26Router *ContractICS26RouterCaller) GetCommitment(opts *bind.CallOpts, hashedPath [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "getCommitment", hashedPath)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 hashedPath) view returns(bytes32)
func (_ContractICS26Router *ContractICS26RouterSession) GetCommitment(hashedPath [32]byte) ([32]byte, error) {
	return _ContractICS26Router.Contract.GetCommitment(&_ContractICS26Router.CallOpts, hashedPath)
}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 hashedPath) view returns(bytes32)
func (_ContractICS26Router *ContractICS26RouterCallerSession) GetCommitment(hashedPath [32]byte) ([32]byte, error) {
	return _ContractICS26Router.Contract.GetCommitment(&_ContractICS26Router.CallOpts, hashedPath)
}

// GetCounterparty is a free data retrieval call binding the contract method 0xb0777bfa.
//
// Solidity: function getCounterparty(string clientId) view returns((string,bytes[]))
func (_ContractICS26Router *ContractICS26RouterCaller) GetCounterparty(opts *bind.CallOpts, clientId string) (IICS02ClientMsgsCounterpartyInfo, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "getCounterparty", clientId)

	if err != nil {
		return *new(IICS02ClientMsgsCounterpartyInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IICS02ClientMsgsCounterpartyInfo)).(*IICS02ClientMsgsCounterpartyInfo)

	return out0, err

}

// GetCounterparty is a free data retrieval call binding the contract method 0xb0777bfa.
//
// Solidity: function getCounterparty(string clientId) view returns((string,bytes[]))
func (_ContractICS26Router *ContractICS26RouterSession) GetCounterparty(clientId string) (IICS02ClientMsgsCounterpartyInfo, error) {
	return _ContractICS26Router.Contract.GetCounterparty(&_ContractICS26Router.CallOpts, clientId)
}

// GetCounterparty is a free data retrieval call binding the contract method 0xb0777bfa.
//
// Solidity: function getCounterparty(string clientId) view returns((string,bytes[]))
func (_ContractICS26Router *ContractICS26RouterCallerSession) GetCounterparty(clientId string) (IICS02ClientMsgsCounterpartyInfo, error) {
	return _ContractICS26Router.Contract.GetCounterparty(&_ContractICS26Router.CallOpts, clientId)
}

// GetIBCApp is a free data retrieval call binding the contract method 0x2447af29.
//
// Solidity: function getIBCApp(string portId) view returns(address)
func (_ContractICS26Router *ContractICS26RouterCaller) GetIBCApp(opts *bind.CallOpts, portId string) (common.Address, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "getIBCApp", portId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetIBCApp is a free data retrieval call binding the contract method 0x2447af29.
//
// Solidity: function getIBCApp(string portId) view returns(address)
func (_ContractICS26Router *ContractICS26RouterSession) GetIBCApp(portId string) (common.Address, error) {
	return _ContractICS26Router.Contract.GetIBCApp(&_ContractICS26Router.CallOpts, portId)
}

// GetIBCApp is a free data retrieval call binding the contract method 0x2447af29.
//
// Solidity: function getIBCApp(string portId) view returns(address)
func (_ContractICS26Router *ContractICS26RouterCallerSession) GetIBCApp(portId string) (common.Address, error) {
	return _ContractICS26Router.Contract.GetIBCApp(&_ContractICS26Router.CallOpts, portId)
}

// GetNextClientSeq is a free data retrieval call binding the contract method 0x27f146f3.
//
// Solidity: function getNextClientSeq() view returns(uint256)
func (_ContractICS26Router *ContractICS26RouterCaller) GetNextClientSeq(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "getNextClientSeq")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNextClientSeq is a free data retrieval call binding the contract method 0x27f146f3.
//
// Solidity: function getNextClientSeq() view returns(uint256)
func (_ContractICS26Router *ContractICS26RouterSession) GetNextClientSeq() (*big.Int, error) {
	return _ContractICS26Router.Contract.GetNextClientSeq(&_ContractICS26Router.CallOpts)
}

// GetNextClientSeq is a free data retrieval call binding the contract method 0x27f146f3.
//
// Solidity: function getNextClientSeq() view returns(uint256)
func (_ContractICS26Router *ContractICS26RouterCallerSession) GetNextClientSeq() (*big.Int, error) {
	return _ContractICS26Router.Contract.GetNextClientSeq(&_ContractICS26Router.CallOpts)
}

// IsConsumingScheduledOp is a free data retrieval call binding the contract method 0x8fb36037.
//
// Solidity: function isConsumingScheduledOp() view returns(bytes4)
func (_ContractICS26Router *ContractICS26RouterCaller) IsConsumingScheduledOp(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "isConsumingScheduledOp")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// IsConsumingScheduledOp is a free data retrieval call binding the contract method 0x8fb36037.
//
// Solidity: function isConsumingScheduledOp() view returns(bytes4)
func (_ContractICS26Router *ContractICS26RouterSession) IsConsumingScheduledOp() ([4]byte, error) {
	return _ContractICS26Router.Contract.IsConsumingScheduledOp(&_ContractICS26Router.CallOpts)
}

// IsConsumingScheduledOp is a free data retrieval call binding the contract method 0x8fb36037.
//
// Solidity: function isConsumingScheduledOp() view returns(bytes4)
func (_ContractICS26Router *ContractICS26RouterCallerSession) IsConsumingScheduledOp() ([4]byte, error) {
	return _ContractICS26Router.Contract.IsConsumingScheduledOp(&_ContractICS26Router.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ContractICS26Router *ContractICS26RouterCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ContractICS26Router *ContractICS26RouterSession) ProxiableUUID() ([32]byte, error) {
	return _ContractICS26Router.Contract.ProxiableUUID(&_ContractICS26Router.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ContractICS26Router *ContractICS26RouterCallerSession) ProxiableUUID() ([32]byte, error) {
	return _ContractICS26Router.Contract.ProxiableUUID(&_ContractICS26Router.CallOpts)
}

// AckPacket is a paid mutator transaction binding the contract method 0xfdbd955d.
//
// Solidity: function ackPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes,bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) AckPacket(opts *bind.TransactOpts, msg_ IICS26RouterMsgsMsgAckPacket) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "ackPacket", msg_)
}

// AckPacket is a paid mutator transaction binding the contract method 0xfdbd955d.
//
// Solidity: function ackPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes,bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterSession) AckPacket(msg_ IICS26RouterMsgsMsgAckPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AckPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// AckPacket is a paid mutator transaction binding the contract method 0xfdbd955d.
//
// Solidity: function ackPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes,bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) AckPacket(msg_ IICS26RouterMsgsMsgAckPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AckPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// AddClient is a paid mutator transaction binding the contract method 0x1ec43e23.
//
// Solidity: function addClient(string clientId, (string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterTransactor) AddClient(opts *bind.TransactOpts, clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "addClient", clientId, counterpartyInfo, client)
}

// AddClient is a paid mutator transaction binding the contract method 0x1ec43e23.
//
// Solidity: function addClient(string clientId, (string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterSession) AddClient(clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AddClient(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// AddClient is a paid mutator transaction binding the contract method 0x1ec43e23.
//
// Solidity: function addClient(string clientId, (string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterTransactorSession) AddClient(clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AddClient(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// AddClient0 is a paid mutator transaction binding the contract method 0xe3cb36a0.
//
// Solidity: function addClient((string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterTransactor) AddClient0(opts *bind.TransactOpts, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "addClient0", counterpartyInfo, client)
}

// AddClient0 is a paid mutator transaction binding the contract method 0xe3cb36a0.
//
// Solidity: function addClient((string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterSession) AddClient0(counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AddClient0(&_ContractICS26Router.TransactOpts, counterpartyInfo, client)
}

// AddClient0 is a paid mutator transaction binding the contract method 0xe3cb36a0.
//
// Solidity: function addClient((string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterTransactorSession) AddClient0(counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AddClient0(&_ContractICS26Router.TransactOpts, counterpartyInfo, client)
}

// AddIBCApp is a paid mutator transaction binding the contract method 0x4b720d5b.
//
// Solidity: function addIBCApp(address app) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) AddIBCApp(opts *bind.TransactOpts, app common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "addIBCApp", app)
}

// AddIBCApp is a paid mutator transaction binding the contract method 0x4b720d5b.
//
// Solidity: function addIBCApp(address app) returns()
func (_ContractICS26Router *ContractICS26RouterSession) AddIBCApp(app common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AddIBCApp(&_ContractICS26Router.TransactOpts, app)
}

// AddIBCApp is a paid mutator transaction binding the contract method 0x4b720d5b.
//
// Solidity: function addIBCApp(address app) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) AddIBCApp(app common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AddIBCApp(&_ContractICS26Router.TransactOpts, app)
}

// AddIBCApp0 is a paid mutator transaction binding the contract method 0x5f516889.
//
// Solidity: function addIBCApp(string portId, address app) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) AddIBCApp0(opts *bind.TransactOpts, portId string, app common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "addIBCApp0", portId, app)
}

// AddIBCApp0 is a paid mutator transaction binding the contract method 0x5f516889.
//
// Solidity: function addIBCApp(string portId, address app) returns()
func (_ContractICS26Router *ContractICS26RouterSession) AddIBCApp0(portId string, app common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AddIBCApp0(&_ContractICS26Router.TransactOpts, portId, app)
}

// AddIBCApp0 is a paid mutator transaction binding the contract method 0x5f516889.
//
// Solidity: function addIBCApp(string portId, address app) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) AddIBCApp0(portId string, app common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AddIBCApp0(&_ContractICS26Router.TransactOpts, portId, app)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address authority) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) Initialize(opts *bind.TransactOpts, authority common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "initialize", authority)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address authority) returns()
func (_ContractICS26Router *ContractICS26RouterSession) Initialize(authority common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.Initialize(&_ContractICS26Router.TransactOpts, authority)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address authority) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) Initialize(authority common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.Initialize(&_ContractICS26Router.TransactOpts, authority)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0x29b6eca9.
//
// Solidity: function initializeV2(address authority) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) InitializeV2(opts *bind.TransactOpts, authority common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "initializeV2", authority)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0x29b6eca9.
//
// Solidity: function initializeV2(address authority) returns()
func (_ContractICS26Router *ContractICS26RouterSession) InitializeV2(authority common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.InitializeV2(&_ContractICS26Router.TransactOpts, authority)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0x29b6eca9.
//
// Solidity: function initializeV2(address authority) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) InitializeV2(authority common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.InitializeV2(&_ContractICS26Router.TransactOpts, authority)
}

// MigrateClient is a paid mutator transaction binding the contract method 0xcce0b265.
//
// Solidity: function migrateClient(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) MigrateClient(opts *bind.TransactOpts, clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "migrateClient", clientId, counterpartyInfo, client)
}

// MigrateClient is a paid mutator transaction binding the contract method 0xcce0b265.
//
// Solidity: function migrateClient(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterSession) MigrateClient(clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.MigrateClient(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// MigrateClient is a paid mutator transaction binding the contract method 0xcce0b265.
//
// Solidity: function migrateClient(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) MigrateClient(clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.MigrateClient(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractICS26Router *ContractICS26RouterTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractICS26Router *ContractICS26RouterSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.Multicall(&_ContractICS26Router.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_ContractICS26Router *ContractICS26RouterTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.Multicall(&_ContractICS26Router.TransactOpts, data)
}

// RecvPacket is a paid mutator transaction binding the contract method 0x596e00b9.
//
// Solidity: function recvPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) RecvPacket(opts *bind.TransactOpts, msg_ IICS26RouterMsgsMsgRecvPacket) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "recvPacket", msg_)
}

// RecvPacket is a paid mutator transaction binding the contract method 0x596e00b9.
//
// Solidity: function recvPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterSession) RecvPacket(msg_ IICS26RouterMsgsMsgRecvPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.RecvPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// RecvPacket is a paid mutator transaction binding the contract method 0x596e00b9.
//
// Solidity: function recvPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) RecvPacket(msg_ IICS26RouterMsgsMsgRecvPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.RecvPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// SendPacket is a paid mutator transaction binding the contract method 0x4d6e7ce3.
//
// Solidity: function sendPacket((string,uint64,(string,string,string,string,bytes)) msg_) returns(uint64)
func (_ContractICS26Router *ContractICS26RouterTransactor) SendPacket(opts *bind.TransactOpts, msg_ IICS26RouterMsgsMsgSendPacket) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "sendPacket", msg_)
}

// SendPacket is a paid mutator transaction binding the contract method 0x4d6e7ce3.
//
// Solidity: function sendPacket((string,uint64,(string,string,string,string,bytes)) msg_) returns(uint64)
func (_ContractICS26Router *ContractICS26RouterSession) SendPacket(msg_ IICS26RouterMsgsMsgSendPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.SendPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// SendPacket is a paid mutator transaction binding the contract method 0x4d6e7ce3.
//
// Solidity: function sendPacket((string,uint64,(string,string,string,string,bytes)) msg_) returns(uint64)
func (_ContractICS26Router *ContractICS26RouterTransactorSession) SendPacket(msg_ IICS26RouterMsgsMsgSendPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.SendPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address newAuthority) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) SetAuthority(opts *bind.TransactOpts, newAuthority common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "setAuthority", newAuthority)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address newAuthority) returns()
func (_ContractICS26Router *ContractICS26RouterSession) SetAuthority(newAuthority common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.SetAuthority(&_ContractICS26Router.TransactOpts, newAuthority)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address newAuthority) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) SetAuthority(newAuthority common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.SetAuthority(&_ContractICS26Router.TransactOpts, newAuthority)
}

// SubmitMisbehaviour is a paid mutator transaction binding the contract method 0x9e2e5c83.
//
// Solidity: function submitMisbehaviour(string clientId, bytes misbehaviourMsg) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) SubmitMisbehaviour(opts *bind.TransactOpts, clientId string, misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "submitMisbehaviour", clientId, misbehaviourMsg)
}

// SubmitMisbehaviour is a paid mutator transaction binding the contract method 0x9e2e5c83.
//
// Solidity: function submitMisbehaviour(string clientId, bytes misbehaviourMsg) returns()
func (_ContractICS26Router *ContractICS26RouterSession) SubmitMisbehaviour(clientId string, misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.SubmitMisbehaviour(&_ContractICS26Router.TransactOpts, clientId, misbehaviourMsg)
}

// SubmitMisbehaviour is a paid mutator transaction binding the contract method 0x9e2e5c83.
//
// Solidity: function submitMisbehaviour(string clientId, bytes misbehaviourMsg) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) SubmitMisbehaviour(clientId string, misbehaviourMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.SubmitMisbehaviour(&_ContractICS26Router.TransactOpts, clientId, misbehaviourMsg)
}

// TimeoutPacket is a paid mutator transaction binding the contract method 0x223e357a.
//
// Solidity: function timeoutPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) TimeoutPacket(opts *bind.TransactOpts, msg_ IICS26RouterMsgsMsgTimeoutPacket) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "timeoutPacket", msg_)
}

// TimeoutPacket is a paid mutator transaction binding the contract method 0x223e357a.
//
// Solidity: function timeoutPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterSession) TimeoutPacket(msg_ IICS26RouterMsgsMsgTimeoutPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.TimeoutPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// TimeoutPacket is a paid mutator transaction binding the contract method 0x223e357a.
//
// Solidity: function timeoutPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) TimeoutPacket(msg_ IICS26RouterMsgsMsgTimeoutPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.TimeoutPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x6fbf8079.
//
// Solidity: function updateClient(string clientId, bytes updateMsg) returns(uint8)
func (_ContractICS26Router *ContractICS26RouterTransactor) UpdateClient(opts *bind.TransactOpts, clientId string, updateMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "updateClient", clientId, updateMsg)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x6fbf8079.
//
// Solidity: function updateClient(string clientId, bytes updateMsg) returns(uint8)
func (_ContractICS26Router *ContractICS26RouterSession) UpdateClient(clientId string, updateMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.UpdateClient(&_ContractICS26Router.TransactOpts, clientId, updateMsg)
}

// UpdateClient is a paid mutator transaction binding the contract method 0x6fbf8079.
//
// Solidity: function updateClient(string clientId, bytes updateMsg) returns(uint8)
func (_ContractICS26Router *ContractICS26RouterTransactorSession) UpdateClient(clientId string, updateMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.UpdateClient(&_ContractICS26Router.TransactOpts, clientId, updateMsg)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ContractICS26Router *ContractICS26RouterSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.UpgradeToAndCall(&_ContractICS26Router.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.UpgradeToAndCall(&_ContractICS26Router.TransactOpts, newImplementation, data)
}

// ContractICS26RouterAckPacketIterator is returned from FilterAckPacket and is used to iterate over the raw logs and unpacked data for AckPacket events raised by the ContractICS26Router contract.
type ContractICS26RouterAckPacketIterator struct {
	Event *ContractICS26RouterAckPacket // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterAckPacketIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterAckPacket)
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
		it.Event = new(ContractICS26RouterAckPacket)
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
func (it *ContractICS26RouterAckPacketIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterAckPacketIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterAckPacket represents a AckPacket event raised by the ContractICS26Router contract.
type ContractICS26RouterAckPacket struct {
	ClientId        common.Hash
	Sequence        *big.Int
	Packet          IICS26RouterMsgsPacket
	Acknowledgement []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterAckPacket is a free log retrieval operation binding the contract event 0xf9bab74bcdb634f4d3dd064cc42a13df056598e1c0336905d2f5750fbfb08b7b.
//
// Solidity: event AckPacket(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet, bytes acknowledgement)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterAckPacket(opts *bind.FilterOpts, clientId []string, sequence []*big.Int) (*ContractICS26RouterAckPacketIterator, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "AckPacket", clientIdRule, sequenceRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterAckPacketIterator{contract: _ContractICS26Router.contract, event: "AckPacket", logs: logs, sub: sub}, nil
}

// WatchAckPacket is a free log subscription operation binding the contract event 0xf9bab74bcdb634f4d3dd064cc42a13df056598e1c0336905d2f5750fbfb08b7b.
//
// Solidity: event AckPacket(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet, bytes acknowledgement)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchAckPacket(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterAckPacket, clientId []string, sequence []*big.Int) (event.Subscription, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "AckPacket", clientIdRule, sequenceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterAckPacket)
				if err := _ContractICS26Router.contract.UnpackLog(event, "AckPacket", log); err != nil {
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

// ParseAckPacket is a log parse operation binding the contract event 0xf9bab74bcdb634f4d3dd064cc42a13df056598e1c0336905d2f5750fbfb08b7b.
//
// Solidity: event AckPacket(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet, bytes acknowledgement)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseAckPacket(log types.Log) (*ContractICS26RouterAckPacket, error) {
	event := new(ContractICS26RouterAckPacket)
	if err := _ContractICS26Router.contract.UnpackLog(event, "AckPacket", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterAuthorityUpdatedIterator is returned from FilterAuthorityUpdated and is used to iterate over the raw logs and unpacked data for AuthorityUpdated events raised by the ContractICS26Router contract.
type ContractICS26RouterAuthorityUpdatedIterator struct {
	Event *ContractICS26RouterAuthorityUpdated // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterAuthorityUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterAuthorityUpdated)
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
		it.Event = new(ContractICS26RouterAuthorityUpdated)
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
func (it *ContractICS26RouterAuthorityUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterAuthorityUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterAuthorityUpdated represents a AuthorityUpdated event raised by the ContractICS26Router contract.
type ContractICS26RouterAuthorityUpdated struct {
	Authority common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAuthorityUpdated is a free log retrieval operation binding the contract event 0x2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad.
//
// Solidity: event AuthorityUpdated(address authority)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterAuthorityUpdated(opts *bind.FilterOpts) (*ContractICS26RouterAuthorityUpdatedIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "AuthorityUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterAuthorityUpdatedIterator{contract: _ContractICS26Router.contract, event: "AuthorityUpdated", logs: logs, sub: sub}, nil
}

// WatchAuthorityUpdated is a free log subscription operation binding the contract event 0x2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad.
//
// Solidity: event AuthorityUpdated(address authority)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchAuthorityUpdated(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterAuthorityUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "AuthorityUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterAuthorityUpdated)
				if err := _ContractICS26Router.contract.UnpackLog(event, "AuthorityUpdated", log); err != nil {
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
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseAuthorityUpdated(log types.Log) (*ContractICS26RouterAuthorityUpdated, error) {
	event := new(ContractICS26RouterAuthorityUpdated)
	if err := _ContractICS26Router.contract.UnpackLog(event, "AuthorityUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterIBCAppAddedIterator is returned from FilterIBCAppAdded and is used to iterate over the raw logs and unpacked data for IBCAppAdded events raised by the ContractICS26Router contract.
type ContractICS26RouterIBCAppAddedIterator struct {
	Event *ContractICS26RouterIBCAppAdded // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterIBCAppAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterIBCAppAdded)
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
		it.Event = new(ContractICS26RouterIBCAppAdded)
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
func (it *ContractICS26RouterIBCAppAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterIBCAppAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterIBCAppAdded represents a IBCAppAdded event raised by the ContractICS26Router contract.
type ContractICS26RouterIBCAppAdded struct {
	PortId string
	App    common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterIBCAppAdded is a free log retrieval operation binding the contract event 0xa6ec8e860960e638347460dc632fbe0175c51a5ca130e336138bbe26ff304499.
//
// Solidity: event IBCAppAdded(string portId, address app)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterIBCAppAdded(opts *bind.FilterOpts) (*ContractICS26RouterIBCAppAddedIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "IBCAppAdded")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterIBCAppAddedIterator{contract: _ContractICS26Router.contract, event: "IBCAppAdded", logs: logs, sub: sub}, nil
}

// WatchIBCAppAdded is a free log subscription operation binding the contract event 0xa6ec8e860960e638347460dc632fbe0175c51a5ca130e336138bbe26ff304499.
//
// Solidity: event IBCAppAdded(string portId, address app)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchIBCAppAdded(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterIBCAppAdded) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "IBCAppAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterIBCAppAdded)
				if err := _ContractICS26Router.contract.UnpackLog(event, "IBCAppAdded", log); err != nil {
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

// ParseIBCAppAdded is a log parse operation binding the contract event 0xa6ec8e860960e638347460dc632fbe0175c51a5ca130e336138bbe26ff304499.
//
// Solidity: event IBCAppAdded(string portId, address app)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseIBCAppAdded(log types.Log) (*ContractICS26RouterIBCAppAdded, error) {
	event := new(ContractICS26RouterIBCAppAdded)
	if err := _ContractICS26Router.contract.UnpackLog(event, "IBCAppAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterIBCAppRecvPacketCallbackErrorIterator is returned from FilterIBCAppRecvPacketCallbackError and is used to iterate over the raw logs and unpacked data for IBCAppRecvPacketCallbackError events raised by the ContractICS26Router contract.
type ContractICS26RouterIBCAppRecvPacketCallbackErrorIterator struct {
	Event *ContractICS26RouterIBCAppRecvPacketCallbackError // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterIBCAppRecvPacketCallbackErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterIBCAppRecvPacketCallbackError)
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
		it.Event = new(ContractICS26RouterIBCAppRecvPacketCallbackError)
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
func (it *ContractICS26RouterIBCAppRecvPacketCallbackErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterIBCAppRecvPacketCallbackErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterIBCAppRecvPacketCallbackError represents a IBCAppRecvPacketCallbackError event raised by the ContractICS26Router contract.
type ContractICS26RouterIBCAppRecvPacketCallbackError struct {
	Reason []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterIBCAppRecvPacketCallbackError is a free log retrieval operation binding the contract event 0xb9edb487876e8be10f54e377c1a815a54ad92a6db1c9561dfe8fad2f0d1da84f.
//
// Solidity: event IBCAppRecvPacketCallbackError(bytes reason)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterIBCAppRecvPacketCallbackError(opts *bind.FilterOpts) (*ContractICS26RouterIBCAppRecvPacketCallbackErrorIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "IBCAppRecvPacketCallbackError")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterIBCAppRecvPacketCallbackErrorIterator{contract: _ContractICS26Router.contract, event: "IBCAppRecvPacketCallbackError", logs: logs, sub: sub}, nil
}

// WatchIBCAppRecvPacketCallbackError is a free log subscription operation binding the contract event 0xb9edb487876e8be10f54e377c1a815a54ad92a6db1c9561dfe8fad2f0d1da84f.
//
// Solidity: event IBCAppRecvPacketCallbackError(bytes reason)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchIBCAppRecvPacketCallbackError(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterIBCAppRecvPacketCallbackError) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "IBCAppRecvPacketCallbackError")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterIBCAppRecvPacketCallbackError)
				if err := _ContractICS26Router.contract.UnpackLog(event, "IBCAppRecvPacketCallbackError", log); err != nil {
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

// ParseIBCAppRecvPacketCallbackError is a log parse operation binding the contract event 0xb9edb487876e8be10f54e377c1a815a54ad92a6db1c9561dfe8fad2f0d1da84f.
//
// Solidity: event IBCAppRecvPacketCallbackError(bytes reason)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseIBCAppRecvPacketCallbackError(log types.Log) (*ContractICS26RouterIBCAppRecvPacketCallbackError, error) {
	event := new(ContractICS26RouterIBCAppRecvPacketCallbackError)
	if err := _ContractICS26Router.contract.UnpackLog(event, "IBCAppRecvPacketCallbackError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterICS02ClientAddedIterator is returned from FilterICS02ClientAdded and is used to iterate over the raw logs and unpacked data for ICS02ClientAdded events raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientAddedIterator struct {
	Event *ContractICS26RouterICS02ClientAdded // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterICS02ClientAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterICS02ClientAdded)
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
		it.Event = new(ContractICS26RouterICS02ClientAdded)
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
func (it *ContractICS26RouterICS02ClientAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterICS02ClientAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterICS02ClientAdded represents a ICS02ClientAdded event raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientAdded struct {
	ClientId         string
	CounterpartyInfo IICS02ClientMsgsCounterpartyInfo
	Client           common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterICS02ClientAdded is a free log retrieval operation binding the contract event 0x0ecded31ecd211a73abf0fb3bc09150bbe321a05550fbe29ea0f16b6e25fbfa8.
//
// Solidity: event ICS02ClientAdded(string clientId, (string,bytes[]) counterpartyInfo, address client)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterICS02ClientAdded(opts *bind.FilterOpts) (*ContractICS26RouterICS02ClientAddedIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "ICS02ClientAdded")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterICS02ClientAddedIterator{contract: _ContractICS26Router.contract, event: "ICS02ClientAdded", logs: logs, sub: sub}, nil
}

// WatchICS02ClientAdded is a free log subscription operation binding the contract event 0x0ecded31ecd211a73abf0fb3bc09150bbe321a05550fbe29ea0f16b6e25fbfa8.
//
// Solidity: event ICS02ClientAdded(string clientId, (string,bytes[]) counterpartyInfo, address client)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchICS02ClientAdded(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterICS02ClientAdded) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "ICS02ClientAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterICS02ClientAdded)
				if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientAdded", log); err != nil {
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

// ParseICS02ClientAdded is a log parse operation binding the contract event 0x0ecded31ecd211a73abf0fb3bc09150bbe321a05550fbe29ea0f16b6e25fbfa8.
//
// Solidity: event ICS02ClientAdded(string clientId, (string,bytes[]) counterpartyInfo, address client)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseICS02ClientAdded(log types.Log) (*ContractICS26RouterICS02ClientAdded, error) {
	event := new(ContractICS26RouterICS02ClientAdded)
	if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterICS02ClientMigratedIterator is returned from FilterICS02ClientMigrated and is used to iterate over the raw logs and unpacked data for ICS02ClientMigrated events raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientMigratedIterator struct {
	Event *ContractICS26RouterICS02ClientMigrated // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterICS02ClientMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterICS02ClientMigrated)
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
		it.Event = new(ContractICS26RouterICS02ClientMigrated)
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
func (it *ContractICS26RouterICS02ClientMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterICS02ClientMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterICS02ClientMigrated represents a ICS02ClientMigrated event raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientMigrated struct {
	ClientId         string
	CounterpartyInfo IICS02ClientMsgsCounterpartyInfo
	Client           common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterICS02ClientMigrated is a free log retrieval operation binding the contract event 0x23c2e29d6ae84e79fa116b8afd6e28ddc1de7f473d3edb407fbd08093c3ed6bf.
//
// Solidity: event ICS02ClientMigrated(string clientId, (string,bytes[]) counterpartyInfo, address client)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterICS02ClientMigrated(opts *bind.FilterOpts) (*ContractICS26RouterICS02ClientMigratedIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "ICS02ClientMigrated")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterICS02ClientMigratedIterator{contract: _ContractICS26Router.contract, event: "ICS02ClientMigrated", logs: logs, sub: sub}, nil
}

// WatchICS02ClientMigrated is a free log subscription operation binding the contract event 0x23c2e29d6ae84e79fa116b8afd6e28ddc1de7f473d3edb407fbd08093c3ed6bf.
//
// Solidity: event ICS02ClientMigrated(string clientId, (string,bytes[]) counterpartyInfo, address client)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchICS02ClientMigrated(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterICS02ClientMigrated) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "ICS02ClientMigrated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterICS02ClientMigrated)
				if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientMigrated", log); err != nil {
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

// ParseICS02ClientMigrated is a log parse operation binding the contract event 0x23c2e29d6ae84e79fa116b8afd6e28ddc1de7f473d3edb407fbd08093c3ed6bf.
//
// Solidity: event ICS02ClientMigrated(string clientId, (string,bytes[]) counterpartyInfo, address client)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseICS02ClientMigrated(log types.Log) (*ContractICS26RouterICS02ClientMigrated, error) {
	event := new(ContractICS26RouterICS02ClientMigrated)
	if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterICS02ClientUpdatedIterator is returned from FilterICS02ClientUpdated and is used to iterate over the raw logs and unpacked data for ICS02ClientUpdated events raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientUpdatedIterator struct {
	Event *ContractICS26RouterICS02ClientUpdated // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterICS02ClientUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterICS02ClientUpdated)
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
		it.Event = new(ContractICS26RouterICS02ClientUpdated)
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
func (it *ContractICS26RouterICS02ClientUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterICS02ClientUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterICS02ClientUpdated represents a ICS02ClientUpdated event raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientUpdated struct {
	ClientId string
	Result   uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterICS02ClientUpdated is a free log retrieval operation binding the contract event 0x87bbef2779889a19f0435ddca81fda94132c06ffddb0ea73def256307a293aef.
//
// Solidity: event ICS02ClientUpdated(string clientId, uint8 result)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterICS02ClientUpdated(opts *bind.FilterOpts) (*ContractICS26RouterICS02ClientUpdatedIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "ICS02ClientUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterICS02ClientUpdatedIterator{contract: _ContractICS26Router.contract, event: "ICS02ClientUpdated", logs: logs, sub: sub}, nil
}

// WatchICS02ClientUpdated is a free log subscription operation binding the contract event 0x87bbef2779889a19f0435ddca81fda94132c06ffddb0ea73def256307a293aef.
//
// Solidity: event ICS02ClientUpdated(string clientId, uint8 result)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchICS02ClientUpdated(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterICS02ClientUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "ICS02ClientUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterICS02ClientUpdated)
				if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientUpdated", log); err != nil {
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

// ParseICS02ClientUpdated is a log parse operation binding the contract event 0x87bbef2779889a19f0435ddca81fda94132c06ffddb0ea73def256307a293aef.
//
// Solidity: event ICS02ClientUpdated(string clientId, uint8 result)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseICS02ClientUpdated(log types.Log) (*ContractICS26RouterICS02ClientUpdated, error) {
	event := new(ContractICS26RouterICS02ClientUpdated)
	if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterICS02MisbehaviourSubmittedIterator is returned from FilterICS02MisbehaviourSubmitted and is used to iterate over the raw logs and unpacked data for ICS02MisbehaviourSubmitted events raised by the ContractICS26Router contract.
type ContractICS26RouterICS02MisbehaviourSubmittedIterator struct {
	Event *ContractICS26RouterICS02MisbehaviourSubmitted // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterICS02MisbehaviourSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterICS02MisbehaviourSubmitted)
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
		it.Event = new(ContractICS26RouterICS02MisbehaviourSubmitted)
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
func (it *ContractICS26RouterICS02MisbehaviourSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterICS02MisbehaviourSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterICS02MisbehaviourSubmitted represents a ICS02MisbehaviourSubmitted event raised by the ContractICS26Router contract.
type ContractICS26RouterICS02MisbehaviourSubmitted struct {
	ClientId string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterICS02MisbehaviourSubmitted is a free log retrieval operation binding the contract event 0xa263f0a976b2937a51fd2e416491cf0ca724d5499fa870715929dfde4ee4a430.
//
// Solidity: event ICS02MisbehaviourSubmitted(string clientId)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterICS02MisbehaviourSubmitted(opts *bind.FilterOpts) (*ContractICS26RouterICS02MisbehaviourSubmittedIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "ICS02MisbehaviourSubmitted")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterICS02MisbehaviourSubmittedIterator{contract: _ContractICS26Router.contract, event: "ICS02MisbehaviourSubmitted", logs: logs, sub: sub}, nil
}

// WatchICS02MisbehaviourSubmitted is a free log subscription operation binding the contract event 0xa263f0a976b2937a51fd2e416491cf0ca724d5499fa870715929dfde4ee4a430.
//
// Solidity: event ICS02MisbehaviourSubmitted(string clientId)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchICS02MisbehaviourSubmitted(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterICS02MisbehaviourSubmitted) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "ICS02MisbehaviourSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterICS02MisbehaviourSubmitted)
				if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02MisbehaviourSubmitted", log); err != nil {
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

// ParseICS02MisbehaviourSubmitted is a log parse operation binding the contract event 0xa263f0a976b2937a51fd2e416491cf0ca724d5499fa870715929dfde4ee4a430.
//
// Solidity: event ICS02MisbehaviourSubmitted(string clientId)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseICS02MisbehaviourSubmitted(log types.Log) (*ContractICS26RouterICS02MisbehaviourSubmitted, error) {
	event := new(ContractICS26RouterICS02MisbehaviourSubmitted)
	if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02MisbehaviourSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractICS26Router contract.
type ContractICS26RouterInitializedIterator struct {
	Event *ContractICS26RouterInitialized // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterInitialized)
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
		it.Event = new(ContractICS26RouterInitialized)
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
func (it *ContractICS26RouterInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterInitialized represents a Initialized event raised by the ContractICS26Router contract.
type ContractICS26RouterInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractICS26RouterInitializedIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterInitializedIterator{contract: _ContractICS26Router.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterInitialized)
				if err := _ContractICS26Router.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseInitialized(log types.Log) (*ContractICS26RouterInitialized, error) {
	event := new(ContractICS26RouterInitialized)
	if err := _ContractICS26Router.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterNoopIterator is returned from FilterNoop and is used to iterate over the raw logs and unpacked data for Noop events raised by the ContractICS26Router contract.
type ContractICS26RouterNoopIterator struct {
	Event *ContractICS26RouterNoop // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterNoopIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterNoop)
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
		it.Event = new(ContractICS26RouterNoop)
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
func (it *ContractICS26RouterNoopIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterNoopIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterNoop represents a Noop event raised by the ContractICS26Router contract.
type ContractICS26RouterNoop struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNoop is a free log retrieval operation binding the contract event 0xd08bf58b0e4eec5bfc697a4fdbb6839057fbf4dd06f1b1ce07445c0e5a654caf.
//
// Solidity: event Noop()
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterNoop(opts *bind.FilterOpts) (*ContractICS26RouterNoopIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "Noop")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterNoopIterator{contract: _ContractICS26Router.contract, event: "Noop", logs: logs, sub: sub}, nil
}

// WatchNoop is a free log subscription operation binding the contract event 0xd08bf58b0e4eec5bfc697a4fdbb6839057fbf4dd06f1b1ce07445c0e5a654caf.
//
// Solidity: event Noop()
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchNoop(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterNoop) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "Noop")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterNoop)
				if err := _ContractICS26Router.contract.UnpackLog(event, "Noop", log); err != nil {
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

// ParseNoop is a log parse operation binding the contract event 0xd08bf58b0e4eec5bfc697a4fdbb6839057fbf4dd06f1b1ce07445c0e5a654caf.
//
// Solidity: event Noop()
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseNoop(log types.Log) (*ContractICS26RouterNoop, error) {
	event := new(ContractICS26RouterNoop)
	if err := _ContractICS26Router.contract.UnpackLog(event, "Noop", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterSendPacketIterator is returned from FilterSendPacket and is used to iterate over the raw logs and unpacked data for SendPacket events raised by the ContractICS26Router contract.
type ContractICS26RouterSendPacketIterator struct {
	Event *ContractICS26RouterSendPacket // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterSendPacketIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterSendPacket)
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
		it.Event = new(ContractICS26RouterSendPacket)
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
func (it *ContractICS26RouterSendPacketIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterSendPacketIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterSendPacket represents a SendPacket event raised by the ContractICS26Router contract.
type ContractICS26RouterSendPacket struct {
	ClientId common.Hash
	Sequence *big.Int
	Packet   IICS26RouterMsgsPacket
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSendPacket is a free log retrieval operation binding the contract event 0xab3a4458a269be61dfa43faa33aa7b1f5d570716f83ad078bc2ba5dab039abae.
//
// Solidity: event SendPacket(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterSendPacket(opts *bind.FilterOpts, clientId []string, sequence []*big.Int) (*ContractICS26RouterSendPacketIterator, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "SendPacket", clientIdRule, sequenceRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterSendPacketIterator{contract: _ContractICS26Router.contract, event: "SendPacket", logs: logs, sub: sub}, nil
}

// WatchSendPacket is a free log subscription operation binding the contract event 0xab3a4458a269be61dfa43faa33aa7b1f5d570716f83ad078bc2ba5dab039abae.
//
// Solidity: event SendPacket(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchSendPacket(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterSendPacket, clientId []string, sequence []*big.Int) (event.Subscription, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "SendPacket", clientIdRule, sequenceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterSendPacket)
				if err := _ContractICS26Router.contract.UnpackLog(event, "SendPacket", log); err != nil {
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

// ParseSendPacket is a log parse operation binding the contract event 0xab3a4458a269be61dfa43faa33aa7b1f5d570716f83ad078bc2ba5dab039abae.
//
// Solidity: event SendPacket(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseSendPacket(log types.Log) (*ContractICS26RouterSendPacket, error) {
	event := new(ContractICS26RouterSendPacket)
	if err := _ContractICS26Router.contract.UnpackLog(event, "SendPacket", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterTimeoutPacketIterator is returned from FilterTimeoutPacket and is used to iterate over the raw logs and unpacked data for TimeoutPacket events raised by the ContractICS26Router contract.
type ContractICS26RouterTimeoutPacketIterator struct {
	Event *ContractICS26RouterTimeoutPacket // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterTimeoutPacketIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterTimeoutPacket)
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
		it.Event = new(ContractICS26RouterTimeoutPacket)
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
func (it *ContractICS26RouterTimeoutPacketIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterTimeoutPacketIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterTimeoutPacket represents a TimeoutPacket event raised by the ContractICS26Router contract.
type ContractICS26RouterTimeoutPacket struct {
	ClientId common.Hash
	Sequence *big.Int
	Packet   IICS26RouterMsgsPacket
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTimeoutPacket is a free log retrieval operation binding the contract event 0x01e5ed58494819ef3f6480dd08e433b7c08ed75c7abdf2c22c6f04b71340a168.
//
// Solidity: event TimeoutPacket(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterTimeoutPacket(opts *bind.FilterOpts, clientId []string, sequence []*big.Int) (*ContractICS26RouterTimeoutPacketIterator, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "TimeoutPacket", clientIdRule, sequenceRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterTimeoutPacketIterator{contract: _ContractICS26Router.contract, event: "TimeoutPacket", logs: logs, sub: sub}, nil
}

// WatchTimeoutPacket is a free log subscription operation binding the contract event 0x01e5ed58494819ef3f6480dd08e433b7c08ed75c7abdf2c22c6f04b71340a168.
//
// Solidity: event TimeoutPacket(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchTimeoutPacket(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterTimeoutPacket, clientId []string, sequence []*big.Int) (event.Subscription, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "TimeoutPacket", clientIdRule, sequenceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterTimeoutPacket)
				if err := _ContractICS26Router.contract.UnpackLog(event, "TimeoutPacket", log); err != nil {
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

// ParseTimeoutPacket is a log parse operation binding the contract event 0x01e5ed58494819ef3f6480dd08e433b7c08ed75c7abdf2c22c6f04b71340a168.
//
// Solidity: event TimeoutPacket(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseTimeoutPacket(log types.Log) (*ContractICS26RouterTimeoutPacket, error) {
	event := new(ContractICS26RouterTimeoutPacket)
	if err := _ContractICS26Router.contract.UnpackLog(event, "TimeoutPacket", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ContractICS26Router contract.
type ContractICS26RouterUpgradedIterator struct {
	Event *ContractICS26RouterUpgraded // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterUpgraded)
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
		it.Event = new(ContractICS26RouterUpgraded)
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
func (it *ContractICS26RouterUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterUpgraded represents a Upgraded event raised by the ContractICS26Router contract.
type ContractICS26RouterUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ContractICS26RouterUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterUpgradedIterator{contract: _ContractICS26Router.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterUpgraded)
				if err := _ContractICS26Router.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseUpgraded(log types.Log) (*ContractICS26RouterUpgraded, error) {
	event := new(ContractICS26RouterUpgraded)
	if err := _ContractICS26Router.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterWriteAcknowledgementIterator is returned from FilterWriteAcknowledgement and is used to iterate over the raw logs and unpacked data for WriteAcknowledgement events raised by the ContractICS26Router contract.
type ContractICS26RouterWriteAcknowledgementIterator struct {
	Event *ContractICS26RouterWriteAcknowledgement // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterWriteAcknowledgementIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterWriteAcknowledgement)
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
		it.Event = new(ContractICS26RouterWriteAcknowledgement)
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
func (it *ContractICS26RouterWriteAcknowledgementIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterWriteAcknowledgementIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterWriteAcknowledgement represents a WriteAcknowledgement event raised by the ContractICS26Router contract.
type ContractICS26RouterWriteAcknowledgement struct {
	ClientId         common.Hash
	Sequence         *big.Int
	Packet           IICS26RouterMsgsPacket
	Acknowledgements [][]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterWriteAcknowledgement is a free log retrieval operation binding the contract event 0x76765590e2b799b0506100f8a6610cfecab2c71e8e1f8aa981b099aff0dfdb74.
//
// Solidity: event WriteAcknowledgement(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet, bytes[] acknowledgements)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterWriteAcknowledgement(opts *bind.FilterOpts, clientId []string, sequence []*big.Int) (*ContractICS26RouterWriteAcknowledgementIterator, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "WriteAcknowledgement", clientIdRule, sequenceRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterWriteAcknowledgementIterator{contract: _ContractICS26Router.contract, event: "WriteAcknowledgement", logs: logs, sub: sub}, nil
}

// WatchWriteAcknowledgement is a free log subscription operation binding the contract event 0x76765590e2b799b0506100f8a6610cfecab2c71e8e1f8aa981b099aff0dfdb74.
//
// Solidity: event WriteAcknowledgement(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet, bytes[] acknowledgements)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchWriteAcknowledgement(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterWriteAcknowledgement, clientId []string, sequence []*big.Int) (event.Subscription, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "WriteAcknowledgement", clientIdRule, sequenceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterWriteAcknowledgement)
				if err := _ContractICS26Router.contract.UnpackLog(event, "WriteAcknowledgement", log); err != nil {
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

// ParseWriteAcknowledgement is a log parse operation binding the contract event 0x76765590e2b799b0506100f8a6610cfecab2c71e8e1f8aa981b099aff0dfdb74.
//
// Solidity: event WriteAcknowledgement(string indexed clientId, uint256 indexed sequence, (uint64,string,string,uint64,(string,string,string,string,bytes)[]) packet, bytes[] acknowledgements)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseWriteAcknowledgement(log types.Log) (*ContractICS26RouterWriteAcknowledgement, error) {
	event := new(ContractICS26RouterWriteAcknowledgement)
	if err := _ContractICS26Router.contract.UnpackLog(event, "WriteAcknowledgement", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
