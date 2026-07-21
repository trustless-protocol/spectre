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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ackPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgAckPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgement\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"membershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addClient\",\"inputs\":[{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addIBCApp\",\"inputs\":[{\"name\":\"app\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addIBCApp\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"app\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractILightClient\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitment\",\"inputs\":[{\"name\":\"hashedPath\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCounterparty\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIBCApp\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIIBCApp\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLightClientMigratorRole\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getNextClientSeq\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializeV2\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isConsumingScheduledOp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recvPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgRecvPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"membershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgSendPacket\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAuthority\",\"inputs\":[{\"name\":\"newAuthority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitMisbehaviour\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"timeoutPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgTimeoutPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonMembershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unfreezeClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateApplicationState\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateConsensusState\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AckPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgement\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCAppAdded\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"app\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCAppRecvPacketCallbackError\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientAdded\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientMigrated\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientUnfrozen\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientUpdated\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"result\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02MisbehaviourSubmitted\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Noop\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SendPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TimeoutPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WriteAcknowledgement\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgements\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessManagedInvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessManagedRequiredDelay\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"AccessManagedUnauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"DefaultAdminRoleCannotBeGranted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCAppNotFound\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCAsyncAcknowledgementNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCClientAlreadyExists\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientNotFound\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCCounterpartyClientNotFound\",\"inputs\":[{\"name\":\"counterpartyClientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCErrorUniversalAcknowledgement\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCFailedCallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCInvalidClientId\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidCounterparty\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidPortIdentifier\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidTimeoutDuration\",\"inputs\":[{\"name\":\"maxTimeoutDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualTimeoutDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidTimeoutTimestamp\",\"inputs\":[{\"name\":\"timeoutTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"comparedTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCMultiPayloadPacketNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCPacketAcknowledgementAlreadyExists\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"IBCPacketCommitmentAlreadyExists\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"IBCPacketCommitmentMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IBCPacketReceiptMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IBCPortAlreadyExists\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCUnauthorizedMigrator\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"IBCUnauthorizedSender\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMerklePrefix\",\"inputs\":[{\"name\":\"prefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"NoAcknowledgements\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"Unreachable\",\"inputs\":[]}]",
	Bin: "0x60a080604052346100c257306080525f5160206160b25f395f51905f525460ff8160401c166100b3576002600160401b03196001600160401b03821601610060575b604051615feb90816100c782396080518181816114af01526115640152f35b6001600160401b0319166001600160401b039081175f5160206160b25f395f51905f525581527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f80610041565b63f92ee8a960e01b5f5260045ffd5b5f80fdfe60806040526004361015610011575f80fd5b5f5f3560e01c80631ec43e2314611ff0578063223e357a14611fc85780632447af2914611f9357806327f146f314611f5757806329b6eca914611d345780633f4ba83a14611c755780634b720d5b14611b4d5780634d6e7ce3146117985780634f1ef2861461152757806352d1902d14611495578063596e00b9146114655780635c975abb146114245780635f516889146113545780637795820c1461130b5780637a9e5e4b1461127b5780637eb789321461122f5780638456cb59146111965780638fb36037146111165780639c11bece1461105c5780639e2e5c8314610f53578063ac9650d814610e36578063ad3cb1cc14610dd5578063b0777bfa14610d63578063b0830ab914610d16578063bf7e214f14610ce3578063c45f24ca14610be7578063c4d66de814610a73578063cce0b26514610579578063d395af7f14610450578063e3cb36a0146101c05763fdbd955d1461016f575f80fd5b346101bd5761019761018036612164565b610188615289565b6101923633614afa565b614732565b807f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d80f35b80fd5b50346101bd5760403660031901126101bd57600435906001600160401b0382116101bd57604060031983360301126101bd576101fa61208e565b91610203614cca565b927f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960254915f19831461043c57600183017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a449602558383807a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000811015610411575b50806d04ee2d6d415b85acef8100000000600a9210156103f6575b662386f26fc100008110156103e2575b6305f5e1008110156103d1575b6127108110156103c2575b60648110156103b4575b10156103ac575b6001810193600a60216102fb6102e588612238565b976102f3604051998a612217565b808952612238565b602088019490601f19013686378701015b5f1901917f30313233343536373839616263646566000000000000000000000000000000008282061a83530490811561034757600a9061030c565b50506103a8956020958661038b936103949760405199858b9651918291018588015e85019083820190858252519283915e010190815203601f198101865285612217565b60040183614ed3565b60405191829160208352602083019061210d565b0390f35b6001016102d0565b6064600291049201916102c9565b612710600491049201916102bf565b6305f5e100600891049201916102b4565b662386f26fc10000601091049201916102a7565b6d04ee2d6d415b85acef810000000060209104920191610297565b604092507a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000090049050600a61027c565b602484634e487b7160e01b81526011600452fd5b50346101bd5760206104c3610464366122a7565b610472959293953633614afa565b6001600160a01b0361048485886142e6565b1690876040518097819582947f92c19bdc000000000000000000000000000000000000000000000000000000008452886004850152602484019161236d565b03925af191821561056e576020949261051f575b50819261050b7f87bbef2779889a19f0435ddca81fda94132c06ffddb0ea73def256307a293aef9360405193849384614385565b0390a16040519061051b816122f7565b8152f35b7f87bbef2779889a19f0435ddca81fda94132c06ffddb0ea73def256307a293aef925061056190853d8711610567575b6105598183612217565b81019061436d565b916104d7565b503d61054f565b6040513d86823e3d90fd5b50346101bd57610588366120a4565b92919061059582846142e6565b506001600160a01b035f516020615f9f5f395f51905f5254166040517fa166aa89000000000000000000000000000000000000000000000000000000008152306004820152602081602481855afa8015610a68578790610a2d575b6105ff915084863392156145e8565b60406001600160401b0360446106158688615ac4565b835194859384927fd1f856ee0000000000000000000000000000000000000000000000000000000084521660048301523360248301525afa8015610a2257869087906109ca575b6106719250816109bb575b50838533926145e8565b61067b82846136de565b6106858280612420565b906001600160401b038211610941576106a8826106a285546143ab565b8561465b565b8790601f83116001146109555791806106da92600195948b9261083f575b50508160011b915f199060031b1c19161790565b81555b016106eb60208301836123eb565b916801000000000000000083116109415780548382558084106108c7575b508752602087208791805b8484106107cd5750505050506107c16001600160a01b037f23c2e29d6ae84e79fa116b8afd6e28ddc1de7f473d3edb407fbd08093c3ed6bf951691604051848682376020818681017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496008152030190206001600160a01b0384166001600160a01b03198254161790556107b360405195869560608752606087019161236d565b9084820360208601526146a0565b9060408301520390a180f35b6107d78183612420565b906001600160401b0382116108b3576107fa826107f487546143ab565b8761465b565b8b908c601f841160011461084a578360019592946020948796610830949261083f5750508160011b915f199060031b1c19161790565b86555b01930193019291610714565b013590505f806106c6565b91601f19841687845260208420935b81811061089b5750936020936001969387969383889510610882575b505050811b018655610833565b01355f19600384901b60f8161c191690555f8080610875565b91936020600181928787013581550195019201610859565b60248c634e487b7160e01b81526041600452fd5b8189528360208a2091820191015b8181106108e25750610709565b808a6108f0600193546143ab565b806108fe575b5050016108d5565b601f811184146109155750508a81555b8a5f6108f6565b83601f60208486610930965220920160051c82019101614645565b808b528a602081208183555561090e565b602488634e487b7160e01b81526041600452fd5b8389526020892091601f1984168a5b8181106109a3575091600195949291838795931061098a575b505050811b0181556106dd565b01355f19600384901b60f8161c191690555f808061097d565b91936020600181928787013581550195019201610964565b63ffffffff915016155f610667565b50506040813d604011610a1a575b816109e560409383612217565b81010312610a165760206109f8826145db565b91015163ffffffff81168103610a1257610671919061065c565b8680fd5b8580fd5b3d91506109d8565b6040513d88823e3d90fd5b506020813d602011610a60575b81610a4760209383612217565b81010312610a1257610a5b6105ff916145db565b6105f0565b3d9150610a3a565b6040513d89823e3d90fd5b50346101bd5760203660031901126101bd57610a8d612078565b5f516020615fbf5f395f51905f52546001600160401b0381169081610bd85760401c60ff16908115610bcc575b50610bbd57610b6590610af360026001600160401b03195f516020615fbf5f395f51905f525416175f516020615fbf5f395f51905f5255565b6801000000000000000068ff0000000000000000195f516020615fbf5f395f51905f525416175f516020615fbf5f395f51905f5255610b30615cd5565b610b38615cd5565b610b40615cd5565b610b48615cd5565b610b50615cd5565b610b58615cd5565b610b60615cd5565b615a64565b68ff0000000000000000195f516020615fbf5f395f51905f5254165f516020615fbf5f395f51905f52557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a180f35b60048263f92ee8a960e01b8152fd5b6002915010155f610aba565b60048463f92ee8a960e01b8152fd5b50346101bd5760203660031901126101bd576004356001600160401b038111610cdf57610c1890369060040161204b565b90610c233633614afa565b826001600160a01b03610c3684846142e6565b16803b15610cdf578180916004604051809481937f6a28f0000000000000000000000000000000000000000000000000000000000083525af18015610cd457610cbb575b50507f7aa8dd0b42a8f6969fbecb42ee2a8fdec2a1a3a576bfc4623b9e7f98ba76fa9991610cb560405192839260208452602084019161236d565b0390a180f35b81610cc591612217565b610cd057825f610c7a565b8280fd5b6040513d84823e3d90fd5b5080fd5b50346101bd57806003193601126101bd5760206001600160a01b035f516020615f9f5f395f51905f525416604051908152f35b50346101bd5760203660031901126101bd57600435906001600160401b0382116101bd576020610d52610d4c366004860161204b565b90615ac4565b6001600160401b0360405191168152f35b50346101bd5760203660031901126101bd57600435906001600160401b0382116101bd576103a8610da0610d9a366004860161204b565b906143e3565b604051918291602083526020610dc18251604083870152606086019061210d565b910151838203601f19016040850152612315565b50346101bd57806003193601126101bd57506103a8604051610df8604082612217565b600581527f352e302e30000000000000000000000000000000000000000000000000000000602082015260405191829160208352602083019061210d565b50346101bd5760203660031901126101bd57600435906001600160401b0382116101bd57366023830112156101bd578160040135906001600160401b0382116101bd57602483013660248460051b86010111610cdf576040516020610e9b8183612217565b83825280820192601f198201368537610eb386613c0c565b96855b87811015610f4057600190610f2488808989610f108a610ede60248960051b8c01018c612420565b9190946040519483869484860198893784019083820190898252519283915e010185815203601f198101835282612217565b5190305af4610f1d613c90565b9030615db4565b610f2e828c6137f4565b52610f39818b6137f4565b5001610eb6565b604051848152806103a88187018c612315565b503461105857610f62366122a7565b610f70949392943633614afa565b6001600160a01b03610f8286856142e6565b1691823b1561105857610fcf925f92836040518096819582947fddba653700000000000000000000000000000000000000000000000000000000845260206004850152602484019161236d565b03925af1801561104d57611019575b507fa263f0a976b2937a51fd2e416491cf0ca724d5499fa870715929dfde4ee4a4309192610cb560405192839260208452602084019161236d565b7fa263f0a976b2937a51fd2e416491cf0ca724d5499fa870715929dfde4ee4a43092505f61104691612217565b5f91610fde565b6040513d5f823e3d90fd5b5f80fd5b34611058576110ce602061106f366122a7565b61107d949293943633614afa565b6001600160a01b0361108f85876142e6565b16905f6040518098819582947f4b1872be000000000000000000000000000000000000000000000000000000008452886004850152602484019161236d565b03925af190811561104d576020935f9261051f5750819261050b7f87bbef2779889a19f0435ddca81fda94132c06ffddb0ea73def256307a293aef9360405193849384614385565b34611058575f366003190112611058575f516020615f9f5f395f51905f525460a01c60ff161561118e5760207f8fb36037000000000000000000000000000000000000000000000000000000005b7fffffffff0000000000000000000000000000000000000000000000000000000060405191168152f35b60205f611164565b34611058575f366003190112611058576111b03633614afa565b6111b8615657565b600160ff197fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005416177fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a1005b34611058576020366003190112611058576004356001600160401b0381116110585761126a611264602092369060040161204b565b906142e6565b6001600160a01b0360405191168152f35b3461105857602036600319011261105857611294612078565b6001600160a01b035f516020615f9f5f395f51905f52541633036112f957803b156112c4576112c290615a64565b005b6001600160a01b03907fc2f31e5e000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b62d1953b60e31b5f523360045260245ffd5b34611058576020366003190112611058576004355f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600602052602060405f2054604051908152f35b34611058576040366003190112611058576004356001600160401b038111611058576113fa61138a6113ff92369060040161204b565b919061139461208e565b9261139d615289565b6113a73633614afa565b6113b4818381151561429d565b6113d681836113cf6113c7368484612253565b805190615d19565b501561429d565b6113f381836113ee6113e9368484612253565b614d16565b61429d565b3691612253565b615533565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d005b34611058575f36600319011261105857602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b34611058576113ff61147636612131565b61147e615657565b611486615289565b6114903633614afa565b613cbf565b34611058575f366003190112611058576001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036114ff5760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b7fe07c8dba000000000000000000000000000000000000000000000000000000005f5260045ffd5b60403660031901126110585761153b612078565b6024356001600160401b0381116110585761155a903690600401612289565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016803014908115611763575b506114ff5761159e3633614afa565b6001600160a01b038216916040517f52d1902d000000000000000000000000000000000000000000000000000000008152602081600481875afa5f918161172f575b506115f85783634c9c8ce360e01b5f5260045260245ffd5b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8592036117045750813b156116f257806001600160a01b03197f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5416177f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a28151156116c1575f808360206112c295519101845af46116bb613c90565b91615db4565b5050346116ca57005b7fb398979f000000000000000000000000000000000000000000000000000000005f5260045ffd5b634c9c8ce360e01b5f5260045260245ffd5b7faa1d49a4000000000000000000000000000000000000000000000000000000005f5260045260245ffd5b9091506020813d60201161175b575b8161174b60209383612217565b81010312611058575190856115e0565b3d915061173e565b90506001600160a01b037f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc541614158361158f565b34611058576117a636612164565b6117ae615657565b6117b6615289565b60408101906117ec6001600160a01b036117e26117dc6117d686866123d6565b80612420565b90613716565b163390331461379d565b6117f9610d9a8280612420565b5190602081019061182961180c836124b4565b4290611817856124b4565b906001600160401b0342911611612dc5565b62015180611848426001600160401b03611842866124b4565b166137da565b1161185e426001600160401b03611842866124b4565b90611b1b575060206118708280612420565b919082604051938492833781017f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747601815203019020926001600160401b03845416946001600160401b038614611b07576001600160401b0360016118f297011694856001600160401b03198254161790556118ea8380612420565b9690946124b4565b604094855197611902878a612217565b60018952601f1987015f5b818110611ad1575050926001600160401b0388999361194261196a94611979978b519d8e61193a81612197565b523691612253565b9660208c01978852898c01521660608a01528260808a01526119653691876123d6565b612e06565b611973826137e7565b526137e7565b5061199081516001600160401b03875116906156aa565b6020815191012090815f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e07476006020526119da845f20541591516001600160401b03885116906156aa565b9015611a905750937fab3a4458a269be61dfa43faa33aa7b1f5d570716f83ad078bc2ba5dab039abae611a64611a478694602098611a1786615723565b905f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e07476008a52875f205580612420565b90818751928392833781015f815203902092855191829182613808565b0390a35f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d51908152f35b611acd9084519182917f91ffd92400000000000000000000000000000000000000000000000000000000835260206004840152602483019061210d565b0390fd5b8089602080938e6060845194611ae686612197565b8186528185870152850152606080850152606060808501520101520161190d565b634e487b7160e01b5f52601160045260245ffd5b7f715fed60000000000000000000000000000000000000000000000000000000005f526201518060045260245260445ffd5b3461105857602036600319011261105857611b66612078565b611b6e615289565b6001600160a01b0381169081611b84602a612238565b90611b926040519283612217565b602a8252611ba0602a612238565b6020830190601f1901368237825115611c615760309053815160011015611c61576078602183015360295b60018111611c135750611be2576113ff9250615533565b827fe22e27eb000000000000000000000000000000000000000000000000000000005f52600452601460245260445ffd5b90600f81166010811015611c61577f3031323334353637383961626364656600000000000000000000000000000000901a611c4e8385614d05565b5360041c908015611b07575f1901611bcb565b634e487b7160e01b5f52603260045260245ffd5b34611058575f36600319011261105857611c8f3633614afa565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff811615611d0c5760ff19167fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a1005b7f8dfc202b000000000000000000000000000000000000000000000000000000005f5260045ffd5b3461105857602036600319011261105857611d4d612078565b5f516020615fbf5f395f51905f52546001600160401b0381169060018203611f3c5760401c60ff16908115611f4b575b50611f3c57611ead90611db660026001600160401b03195f516020615fbf5f395f51905f525416175f516020615fbf5f395f51905f5255565b6801000000000000000068ff0000000000000000195f516020615fbf5f395f51905f525416175f516020615fbf5f395f51905f52556001600160a01b03197fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b00546001600160a01b03811633148015611f04575b611e3490339061379d565b167fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b00556001600160a01b03197fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b0154167fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b0155610b50615cd5565b68ff0000000000000000195f516020615fbf5f395f51905f5254165f516020615fbf5f395f51905f52557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160028152a1005b50611e346001600160a01b037fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b01541633149050611e29565b63f92ee8a960e01b5f5260045ffd5b60029150101582611d7d565b34611058575f3660031901126110585760207f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960254604051908152f35b34611058576020366003190112611058576004356001600160401b0381116110585761126a6117dc602092369060040161204b565b34611058576113ff611fd936612131565b611fe1615289565b611feb3633614afa565b61317d565b34611058576103a86103946113f3612007366120a4565b90612016949293943633614afa565b612023848681151561238d565b61203b84866120366113e9368484612253565b61238d565b612046368587612253565b614ed3565b9181601f84011215611058578235916001600160401b038311611058576020838186019501011161105857565b600435906001600160a01b038216820361105857565b602435906001600160a01b038216820361105857565b6060600319820112611058576004356001600160401b03811161105857816120ce9160040161204b565b92909291602435906001600160401b03821161105857604090829003600319011261105857600401906044356001600160a01b03811681036110585790565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b602060031982011261105857600435906001600160401b0382116110585760409082900360031901126110585760040190565b602060031982011261105857600435906001600160401b0382116110585760609082900360031901126110585760040190565b60a081019081106001600160401b038211176121b257604052565b634e487b7160e01b5f52604160045260245ffd5b604081019081106001600160401b038211176121b257604052565b608081019081106001600160401b038211176121b257604052565b606081019081106001600160401b038211176121b257604052565b90601f801991011681019081106001600160401b038211176121b257604052565b6001600160401b0381116121b257601f01601f191660200190565b92919261225f82612238565b9161226d6040519384612217565b829481845281830111611058578281602093845f960137010152565b9080601f83011215611058578160206122a493359101612253565b90565b6040600319820112611058576004356001600160401b03811161105857816122d19160040161204b565b92909291602435906001600160401b038211611058576122f39160040161204b565b9091565b6003111561230157565b634e487b7160e01b5f52602160045260245ffd5b9080602083519182815201916020808360051b8301019401925f915b83831061234057505050505090565b909192939460208061235e600193601f19868203018752895161210d565b97019301930191939290612331565b908060209392818452848401375f828201840152601f01601f1916010190565b91909115612399575050565b611acd6040519283927f4870bd7400000000000000000000000000000000000000000000000000000000845260206004850152602484019161236d565b903590609e1981360301821215611058570190565b903590601e198136030182121561105857018035906001600160401b03821161105857602001918160051b3603831361105857565b903590601e198136030182121561105857018035906001600160401b0382116110585760200191813603831361105857565b9290921561245f57505050565b6124a29291611acd916040519485947f9fff831f00000000000000000000000000000000000000000000000000000000865260406004870152604486019061210d565b8481036003190160248601529161236d565b356001600160401b03811681036110585790565b35906001600160401b038216820361105857565b9190826040910312611058576040516124f4816121c6565b602061250d818395612505816124c8565b8552016124c8565b910152565b6001600160401b0381116121b25760051b60200190565b9080601f8301121561105857813561254081612512565b9261254e6040519485612217565b81845260208085019260051b820101918383116110585760208201905b83821061257a57505050505090565b81356001600160401b0381116110585760209161259c87848094880101612289565b81520191019061256b565b9080601f83011215611058578135916125bf83612512565b926125cd6040519485612217565b80845260208085019160051b830101918383116110585760208101915b8383106125f957505050505090565b82356001600160401b038111611058578201906040828703601f1901126110585760405190612627826121c6565b60208301356001600160401b0381116110585787602061264992860101612529565b82526040830135916001600160401b0383116110585761267188602080969581960101612289565b838201528152019201916125ea565b919060808382031261105857604051612698816121e1565b809380356001600160401b03811161105857836126b6918301612289565b825260208101356001600160401b03811161105857836126d7918301612289565b602083015260408101356001600160401b0381116110585781016080818503126110585760405190612708826121e1565b8035600381101561105857825260208101356003811015611058576020830152604081013560038110156110585760408301526060810135906001600160401b0382116110585761275b91869101612289565b606082015260408301526060810135906001600160401b03821161105857019180601f840112156110585782359261279284612512565b936127a06040519586612217565b80855260208086019160051b830101918383116110585760208101915b8383106127cf57505050505060600152565b82356001600160401b038111611058578201906060828703601f19011261105857604051906127fd826121fc565b6020830135600381101561105857825260408301356001600160401b0381116110585787602061282f92860101612289565b60208301526060830135916001600160401b0383116110585761285a88602080969581960101612289565b60408201528152019201916127bd565b3590811515820361105857565b9080601f830112156110585781359161288f83612512565b9261289d6040519485612217565b80845260208085019160051b830101918383116110585760208101915b8383106128c957505050505090565b82356001600160401b0381116110585782016020818703601f1901126110585760405190602082018281106001600160401b038211176121b25760405260208101356001600160401b03811161105857602091010186601f820112156110585780359061293582612512565b916129436040519384612217565b80835260208084019160051b830101918983116110585760208101915b83831061297c57505050908252508152602092830192016128ba565b82356001600160401b038111611058578201906060828d03601f19011261105857604051906129aa826121fc565b6020830135600281101561105857825260408301356001600160401b038111611058578d60206129dc92860101612680565b602083015260608301356001600160401b03811161105857602060a0918f950101809403126110585760405191612a1283612197565b83356001600160401b038111611058578e612a2e918601612289565b8352612a3c6020850161286a565b602084015260408401356001600160401b038111611058578e612a60918601612680565b6040840152612a716060850161286a565b60608401526080840135926001600160401b03841161105857612a9a8f60209695879601612680565b60808201526040820152815201920191612960565b919082606091031261105857604051612ac7816121fc565b809280356fffffffffffffffffffffffffffffffff811681036110585760409182918452602081013560208501520135910152565b9080602083519182815201916020808360051b8301019401925f915b838310612b2757505050505090565b9091929394602080612b65600193601f1986820301875289519083612b558351604084526040840190612315565b920151908481840391015261210d565b97019301930191939290612b18565b6060612bf3612ba1612b8f845160808752608087019061210d565b6020850151868203602088015261210d565b608083604086015187840360408901528051612bbc816122f7565b84526020810151612bcc816122f7565b60208501526040810151612bdf816122f7565b60408501520151918185820152019061210d565b910151916060818303910152815180825260208201916020808360051b8301019401925f915b838310612c2857505050505090565b9091929394602080612c79600193601f198682030187528951908151612c4d816122f7565b81526040612c6885840151606087850152606084019061210d565b92015190604081840391015261210d565b97019301930191939290612c19565b9080602083519283815201916020808260051b8401019401925f925b828410612cb357505050505090565b9091929394601f198282030183528551602082019051916020815282518092526040810190602060408460051b8301019401925f915b818310612d0a57505050505060208060019297019301940192919390612ca4565b9091929394603f19828203018552855190815191600283101561230157612dab82602093926001958594526040612d4e858301516060878601526060850190612b74565b9101519160408183039101526080612d8e612d72845160a0855260a085019061210d565b8685015115158785015260408501518482036040860152612b74565b926060810151151560608401520151906080818403910152612b74565b9701950193019190612ce9565b9060018210156123015752565b15612dce575050565b6001600160401b03907f65d30129000000000000000000000000000000000000000000000000000000005f521660045260245260445ffd5b919060a0838203126110585760405190612e1f82612197565b819380356001600160401b0381116110585782612e3d918301612289565b835260208101356001600160401b0381116110585782612e5e918301612289565b602084015260408101356001600160401b0381116110585782612e82918301612289565b604084015260608101356001600160401b0381116110585782612ea6918301612289565b60608401526080810135916001600160401b0383116110585760809261250d9201612289565b6122a4916080612f24612f12612f00612eee865160a0875260a087019061210d565b6020870151868203602088015261210d565b6040860151858203604087015261210d565b6060850151848203606086015261210d565b92015190608081840391015261210d565b9060806001600160a01b0381612f91612f6c612f5a875160a0885260a088019061210d565b6020880151878203602089015261210d565b6001600160401b03604088015116604087015260608701518682036060880152612ecc565b9401511691015290565b9035601e19823603018112156110585701602081359101916001600160401b03821161105857813603831361105857565b9035601e19823603018112156110585701602081359101916001600160401b038211611058578160051b3603831361105857565b906001600160401b03613012836124c8565b16815261307c61305761303c61302b6020860186612f9b565b60a0602087015260a086019161236d565b6130496040860186612f9b565b90858303604087015261236d565b926001600160401b0361306c606083016124c8565b1660608401526080810190612fcc565b9290916080818303910152828152602081019260208160051b83010193835f91609e1982360301945b8484106130b6575050505050505090565b90919293949596601f1982820301835287358781121561105857602061316c600193878394019061315e61315361313861311d6131046130f68780612f9b565b60a0885260a088019161236d565b61311089880188612f9b565b908783038b89015261236d565b61312a6040870187612f9b565b90868303604088015261236d565b6131456060860186612f9b565b90858303606087015261236d565b926080810190612f9b565b91608081850391015261236d565b9901930194019291959493906130a5565b600161319661318c83806123d6565b60808101906123eb565b9050036136b6576131aa61318c82806123d6565b15611c6157806131b9916123d6565b906131d4610d9a6131ca83806123d6565b6020810190612420565b61321a8151602081519101206131fa6113f36131f086806123d6565b6040810190612420565b60208151910120148251906132126131f086806123d6565b929091612452565b61324e61324961322d6131f085806123d6565b919061324161323c87806123d6565b6124b4565b923691612253565b6152fd565b9061325c6020840184612420565b81019190602081840312611058578035906001600160401b03821161105857019261014084840312611058576040519060e082018281106001600160401b038211176121b2576040526132af84866124dc565b825260408501356001600160401b03811161105857846132d09187016125a7565b926020830193845260608601356001600160401b03811161105857856132f7918801612877565b95604084019687526060840192608082013584526133188760a08401612aaf565b96608086019788526101008301359260018410156110585760a08701938452610120810135906001600160401b03821161105857019061335791612529565b9260c08601938452602001519061336d91615343565b825261337988806123d6565b6020810161338691612420565b61338f916142e6565b6001600160a01b03169560405197889687967fa6f031bb0000000000000000000000000000000000000000000000000000000088526004880160209052602488019051906133f1916001600160401b0360208092828151168552015116910152565b51606487016101409052610164870161340991612afc565b90518682036023190160848801526134219190612c88565b935160a48601525180516fffffffffffffffffffffffffffffffff1660c4860152602081015160e48601526040015161010485015251610124840161346591612db8565b518282036023190161014484015261347d9190612315565b03815a6020945f91f1801561104d575f90613682575b6134ce91506001600160401b036134b560606134af86806123d6565b016124b4565b168110156134c860606134af86806123d6565b90612dc5565b6134e06134db82806123d6565b615450565b1561365a57816135606001600160a01b036135066117dc846001600160401b0397612420565b16916135156131ca85806123d6565b959061354e6135276131f088806123d6565b61354561353761323c8b806123d6565b946040519b6113f38d612197565b8a523691612253565b60208801521660408601523690612e06565b6060840152336080840152803b15611058576135b75f939184926040519586809481937f5e32b6b6000000000000000000000000000000000000000000000000000000008352602060048401526024830190612f35565b03925af191821561104d576001600160401b039261364a575b507f01e5ed58494819ef3f6480dd08e433b7c08ed75c7abdf2c22c6f04b71340a1686135ff6131ca83806123d6565b61361961361261323c86809895986123d6565b94806123d6565b9481604051928392833781015f8152039020926136456040519283926020845216956020830190613000565b0390a3565b5f61365491612217565b5f6135d0565b50507fd08bf58b0e4eec5bfc697a4fdbb6839057fbf4dd06f1b1ce07445c0e5a654caf5f80a1565b506020813d6020116136ae575b8161369c60209383612217565b81010312611058576134ce9051613493565b3d915061368f565b7f356f4dbd000000000000000000000000000000000000000000000000000000005f5260045ffd5b60209082604051938492833781017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960181520301902090565b6001600160a01b03604051838382376020818581017fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf600815203019020541691821561376057505090565b611acd6040519283927fa09dbf5900000000000000000000000000000000000000000000000000000000845260206004850152602484019161236d565b156137a55750565b6001600160a01b03907fbe2f2b45000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b91908203918211611b0757565b805115611c615760200190565b8051821015611c615760209160051b010190565b90602082526001600160401b038151166020830152608061385161383b602084015160a0604087015260c086019061210d565b6040840151858203601f1901606087015261210d565b916001600160401b036060820151168285015201519160a0601f1982840301910152815180825260208201916020808360051b8301019401925f915b83831061389c57505050505090565b90919293946020806138ba600193601f198682030187528951612ecc565b9701930193019193929061388d565b919060a083820312611058576040516138e181612197565b80936138ec816124c8565b825260208101356001600160401b038111611058578361390d918301612289565b602083015260408101356001600160401b0381116110585783613931918301612289565b6040830152613942606082016124c8565b60608301526080810135906001600160401b03821161105857019180601f8401121561105857823561397381612512565b936139816040519586612217565b81855260208086019260051b820101918383116110585760208201905b8382106139b057505050505060800152565b81356001600160401b038111611058576020916139d287848094880101612e06565b81520191019061399e565b602081830312611058578035906001600160401b038211611058570161016081830312611058576040519161010083018381106001600160401b038211176121b257604052613a2c81836124dc565b835260408201356001600160401b0381116110585781613a4d9184016125a7565b602084015260608201356001600160401b0381116110585781613a71918401612877565b604084015260808201356060840152613a8d8160a08401612aaf565b608084015261010082013560018110156110585760a08401526101208201356001600160401b0381116110585781613ac6918401612529565b60c08401526101408201356001600160401b03811161105857613ae99201612289565b60e082015290565b906122a49160208152613b1d6020820183516001600160401b0360208092828151168552015116910152565b60e0613bbc613b56613b4060208601516101606060870152610180860190612afc565b6040860151858203601f19016080870152612c88565b606085015160a0850152608085015180516fffffffffffffffffffffffffffffffff1660c0860152602081015160e086015260400151610100850152613ba560a0860151610120860190612db8565b60c0850151848203601f1901610140860152612315565b92015190610160601f198285030191015261210d565b60408051909190613be38382612217565b6001815291601f1901825f5b828110613bfb57505050565b806060602080938501015201613bef565b90613c1682612512565b613c236040519182612217565b8281528092613c34601f1991612512565b01905f5b828110613c4457505050565b806060602080938501015201613c38565b60405190613c64604083612217565b602082527f4774d4a575993f963b1c06573736617a457abef8589178db8d10c94b4ab511ab6020830152565b3d15613cba573d90613ca182612238565b91613caf6040519384612217565b82523d5f602084013e565b606090565b6001613cce61318c83806123d6565b9050036136b657613ce261318c82806123d6565b15611c615780613cf1916123d6565b905f6020613e44613d08610d9a6131f086806123d6565b613d438151848151910120613d2c6113f3613d2389806123d6565b87810190612420565b85815191012014825190613212613d2389806123d6565b613d66613d5560606134af88806123d6565b429061181760606134af8a806123d6565b613d96613d91613d82613d7988806123d6565b86810190612420565b919061324161323c8a806123d6565b6156aa565b613dd5613db4613daf36613daa8a806123d6565b6138c9565b615723565b9185613dcd613dc5828b018b612420565b8101906139dd565b940151615343565b60c08301526040519084820152838152613df0604082612217565b60e08201526001600160a01b03613e0d6112646131f088806123d6565b16906040519485809481937f974a74c400000000000000000000000000000000000000000000000000000000835260048301613af1565b03925af1801561104d5761426e575b50613e66613e6182806123d6565b615982565b1561365a575f80613f41613e78613bd2565b946001600160401b03613ef96001600160a01b03613e9c6117dc6020860186612420565b1692613eab6131ca89806123d6565b9390613ee78a61323c613ede613ed0613ec76131f085806123d6565b939094806123d6565b94604051996113f38b612197565b88523691612253565b60208601521660408401523690612e06565b60608201523360808201526040519485809481937f078c4a79000000000000000000000000000000000000000000000000000000008352602060048401526024830190612f35565b03925af15f91816141f3575b506141685750613f5b613c90565b80511561414057613f9b7fb9edb487876e8be10f54e377c1a815a54ad92a6db1c9561dfe8fad2f0d1da84f9160405191829160208352602083019061210d565b0390a1613fa6613c55565b613faf836137e7565b52613fb9826137e7565b505b613fc581806123d6565b60408101614039613241613fed613ff2613fed613fe28688612420565b9190613241896124b4565b615b21565b6020815191012094855f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205261403160405f2054159582612420565b9390916124b4565b9015614102575061404983615b83565b905f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f20557f76765590e2b799b0506100f8a6610cfecab2c71e8e1f8aa981b099aff0dfdb746140f26136456140a96131f085806123d6565b6140c36140bc61323c88809695966123d6565b96806123d6565b9281604051928392833781015f8152039020946001600160401b03604051948594604086526040860190613000565b9184830360208601521696612315565b611acd906040519182917f40470d7400000000000000000000000000000000000000000000000000000000835260206004840152602483019061210d565b7fadef7fb8000000000000000000000000000000000000000000000000000000005f5260045ffd5b8051156141cb578051602082012061417e613c55565b60208151910120146141a357614193836137e7565b5261419d826137e7565b50613fbb565b7f6b2675e3000000000000000000000000000000000000000000000000000000005f5260045ffd5b7fecfef798000000000000000000000000000000000000000000000000000000005f5260045ffd5b9091503d805f833e6142058183612217565b810190602081830312611058578051906001600160401b038211611058570181601f820112156110585780519061423b82612238565b926142496040519485612217565b8284526020838301011161105857815f9260208093018386015e83010152905f613f4d565b6020813d602011614295575b8161428760209383612217565b810103126110585751613e53565b3d915061427a565b919091156142a9575050565b611acd6040519283927f14d7124700000000000000000000000000000000000000000000000000000000845260206004850152602484019161236d565b6001600160a01b03604051838382376020818581017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a449600815203019020541691821561433057505090565b611acd6040519283927fa0db16fe00000000000000000000000000000000000000000000000000000000845260206004850152602484019161236d565b90816020910312611058575160038110156110585790565b93929160209161439d9160408752604087019161236d565b936143a7836122f7565b0152565b90600182811c921680156143d9575b60208310146143c557565b634e487b7160e01b5f52602260045260245ffd5b91607f16916143ba565b606060206040516143f3816121c6565b828152015261440282826136de565b9160405192614410846121c6565b6040515f825461441f816143ab565b80845290600181169081156145b75750600114614574575b509061444881600194930382612217565b85520180549061445782612512565b916144656040519384612217565b80835260208301915f5260205f20915f905b8282106144d1575050505060208401528251511561449457505090565b611acd6040519283927fdf95155a00000000000000000000000000000000000000000000000000000000845260206004850152602484019161236d565b6040515f85546144e0816143ab565b8084529060018116908115614551575060011461451a575b506001928261450c85946020940382612217565b815201940191019092614477565b5f878152602081209092505b81831061453b575050810160200160016144f8565b6001816020925483868801015201920191614526565b60ff191660208581019190915291151560051b84019091019150600190506144f8565b5f8481526020812094939250905b80821061459b5750919250908101602001614448614437565b9192936001816020925483858801015201910190939291614582565b60ff191660208086019190915291151560051b840190910191506144489050614437565b5190811515820361105857565b929092156145f557505050565b6001600160a01b0361463a6040519485947f36de203600000000000000000000000000000000000000000000000000000000865260406004870152604486019161236d565b911660248301520390fd5b818110614650575050565b5f8155600101614645565b9190601f811161466a57505050565b614694925f5260205f20906020601f840160051c83019310614696575b601f0160051c0190614645565b565b9091508190614687565b906146ca6146bf6146b18480612f9b565b60408552604085019161236d565b926020810190612fcc565b90916020818503910152808352602083019260208260051b82010193835f925b8484106146fa5750505050505090565b909192939495602080614722600193601f1986820301885261471c8b88612f9b565b9061236d565b98019401940192949391906146ea565b600161474161318c83806123d6565b9050036136b65761475561318c82806123d6565b15611c615780614764916123d6565b90614775610d9a6131ca83806123d6565b916147aa8351602081519101206147926113f36131f086806123d6565b60208151910120148451906132126131f086806123d6565b5f60206148476147c3613fed613d826131f088806123d6565b956147cc613bd2565b9661480f6147fd858901996147e46113f38c8c612420565b6147ed826137e7565b526147f7816137e7565b50615b83565b9185613dcd613dc560408c018c612420565b60c0830152604051908482015283815261482a604082612217565b60e08201526001600160a01b03613e0d611264613d7989806123d6565b03925af1801561104d57614acb575b506148646134db83806123d6565b15614aa2576001600160a01b0361487e6117dc8380612420565b1661488c6131ca84806123d6565b61489c6131f086809694966123d6565b6148ac61323c88809594956123d6565b916148b78989612420565b9790916040519560c08701948786106001600160401b038711176121b257613ede61490d946148ee9461491c986040523691612253565b95602086019687526001600160401b0360408701951685523690612e06565b96606085019788523691612253565b6080830190815260a0830191338352853b156110585760405196879586957f428e4e170000000000000000000000000000000000000000000000000000000087526004870160209052516024870160c0905260e4870161497b9161210d565b9051868203602319016044880152614993919061210d565b91516001600160401b03166064860152518482036023190160848601526149ba9190612ecc565b90518382036023190160a48501526149d2919061210d565b90516001600160a01b031660c483015203815a5f948591f1801561104d57613645927ff9bab74bcdb634f4d3dd064cc42a13df056598e1c0336905d2f5750fbfb08b7b92614a8292614a92575b50614a2d6131ca82806123d6565b949095614a51614a4061323c85806123d6565b91614a4b85806123d6565b94612420565b96909781604051928392833781015f8152039020956001600160401b03604051958695604087526040870190613000565b928584036020870152169761236d565b5f614a9c91612217565b5f614a1f565b5050507fd08bf58b0e4eec5bfc697a4fdbb6839057fbf4dd06f1b1ce07445c0e5a654caf5f80a1565b6020813d602011614af2575b81614ae460209383612217565b810103126110585751614856565b3d9150614ad7565b5f516020615f9f5f395f51905f5254916001600160a01b0383169281600411611058575f5f9060405f8151966001600160a01b0360208901917fb700961300000000000000000000000000000000000000000000000000000000835216978860248201523060448201527fffffffff00000000000000000000000000000000000000000000000000000000833516606482015260648152614b9c608482612217565b828052826020525190895afa614cb7575b15614bba575b5050505050565b63ffffffff1615614ca55760ff60a01b191674010000000000000000000000000000000000000000175f516020615f9f5f395f51905f5255823b15611058576020925f92836040518096819582947f94c7d7ee00000000000000000000000000000000000000000000000000000000845260048401526040602484015260448301908082528085848401378181018301859052601f01601f1916010103925af1801561104d57614c95575b5060ff60a01b195f516020615f9f5f395f51905f5254165f516020615f9f5f395f51905f52555f80808080614bb3565b5f614c9f91612217565b5f614c65565b8262d1953b60e31b5f5260045260245ffd5b50505f516020518060201c150290614bad565b60405190614cd9604083612217565b600782527f636c69656e742d000000000000000000000000000000000000000000000000006020830152565b908151811015611c61570160200190565b805160048110908115614ec8575b50614eac57614d6a604051614d3a604082612217565b600881527f6368616e6e656c2d000000000000000000000000000000000000000000000000602082015282615c5e565b8015614eb1575b614eac575f5b8151811015614ea557614d8a8183614d05565b5160f81c60618110159081614e99575b8115614e7b575b8115614e5d575b8115614e1e575b8115614dc9575b50614dc15750505f90565b600101614d77565b6023811491508115614e13575b8115614e08575b8115614dfd575b8115614df2575b505f614db6565b603e9150145f614deb565b603c81149150614de4565b605d81149150614ddd565b605b81149150614dd6565b9050602e81148015614e53575b8015614e49575b8015614e3f575b90614daf565b50602d8114614e39565b50602b8114614e32565b50605f8114614e2b565b9050604181101580614e70575b90614da8565b50605a811115614e6a565b9050603081101580614e8e575b90614da1565b506039811115614e88565b607a8111159150614d9a565b5050600190565b505f90565b50614ec3614ebd614cca565b82615c5e565b614d71565b60809150115f614d24565b9190604051916001600160a01b03845193602081818801968088835e81017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a449600815203019020541661524e576001600160a01b03169160405160208186518085835e81017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496008152030190206001600160a01b0384166001600160a01b03198254161790556020604051809286518091835e81017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a449601815203019020614fb68280612420565b906001600160401b0382116121b257614fd3826106a285546143ab565b5f90601f83116001146151e757918061500492600195945f9261083f5750508160011b915f199060031b1c19161790565b81555b0161501560208301836123eb565b906801000000000000000082116121b2578254828455808310615170575b505f928352602083209290805b838310615095575050505050916107b39161508a7f0ecded31ecd211a73abf0fb3bc09150bbe321a05550fbe29ea0f16b6e25fbfa89460405194859460608652606086019061210d565b9060408301520390a1565b61509f8183612420565b906001600160401b0382116121b2576150c2826150bc89546143ab565b8961465b565b5f90601f831160011461510657926150f7836001959460209487965f9261083f5750508160011b915f199060031b1c19161790565b88555b01950192019193615040565b601f19831691885f5260205f20925f5b818110615158575093602093600196938796938388951061513f575b505050811b0188556150fa565b01355f19600384901b60f8161c191690555f8080615132565b91936020600181928787013581550195019201615116565b835f528260205f2091820191015b81811061518b5750615033565b80615198600192546143ab565b806151a5575b500161517e565b601f811183146151ba57505f81555b5f61519e565b6151d690825f5283601f60205f20920160051c82019101614645565b805f525f60208120818355556151b4565b601f19831691845f5260205f20925f5b818110615236575091600195949291838795931061521d575b505050811b018155615007565b01355f19600384901b60f8161c191690555f8080615210565b919360206001819287870135815501950192016151f7565b6040517f87dfb2670000000000000000000000000000000000000000000000000000000081526020600482015280611acd602482018761210d565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005c6152d55760017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d565b7f3ee5aeb5000000000000000000000000000000000000000000000000000000005f5260045ffd5b60096122a4916020936001600160c01b0319856040519687948051918291018387015e84010191600160f91b835260c01b1660018201520301601f198101835282612217565b91825115615415576153558351613c0c565b905f5b84515f198101908111611b07578110156153965780615379600192876137f4565b5161538482866137f4565b5261538f81856137f4565b5001615358565b509291909180515f198101908111611b075760209283806153ba6153f094866137f4565b51926040519684889551918291018487015e8401908282015f8152815193849201905e01015f815203601f198101845283612217565b515f198101908111611b07576154119161540a82856137f4565b52826137f4565b5090565b6040517fa7c34e4f0000000000000000000000000000000000000000000000000000000081526020600482015280611acd6024820186612315565b61546e613d916154636020840184612420565b9190613241856124b4565b6020815191012090815f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f2054801561552c576154c5613daf6154bb613daf36866138c9565b83149336906138c9565b91156154fe5750505f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e07476006020525f6040812055600190565b7f3f87a2ec000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b5050505f90565b90604051906001600160a01b03835192602081818701958087835e81017fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf600815203019020541661561c5791615611916001600160a01b037fa6ec8e860960e638347460dc632fbe0175c51a5ca130e336138bbe26ff3044999416906020604051809285518091835e81017fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf6008152030190206001600160a01b0382166001600160a01b031982541617905560405192839260408452604084019061210d565b9060208301520390a1565b6040517f837f46a60000000000000000000000000000000000000000000000000000000081526020600482015280611acd602482018661210d565b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300541661568257565b7fd93c0665000000000000000000000000000000000000000000000000000000005f5260045ffd5b60096122a4916020936001600160c01b0319856040519687948051918291018387015e840101917f0100000000000000000000000000000000000000000000000000000000000000835260c01b1660018201520301601f198101835282612217565b60209291908391805192839101825e019081520190565b906020916040516157348482612217565b5f8152905f915b608082015180518410156158985783615753916137f4565b51855f818351604051918183925191829101835e8101838152039060025afa1561104d575f5190865f8180840151604051918183925191829101835e8101838152039060025afa1561104d575f5191875f816040850151604051918183925191829101835e8101838152039060025afa1561104d575f5192885f816060860151604051918183925191829101835e8101838152039060025afa1561104d57885f8160808251960151604051918183925191829101835e8101838152039060025afa1561104d5788935f938451916040519387850195865260408501526060840152608083015260a082015260a0815261584d60c082612217565b604051918291518091835e8101838152039060025afa1561104d576001906158905f51916158826040519384928a840161570c565b03601f198101835282612217565b92019161573b565b509150929192825f816040840151604051918183925191829101835e8101838152039060025afa1561104d57825f606081519301516040516001600160c01b03198482019260c01b168252600881526158f2602882612217565b604051918291518091835e8101838152039060025afa1561104d57825f81815194604051918183925191829101835e8101838152039060025afa1561104d575f9182516040519185830193600160f91b855260218401526041830152606182015260618152615962608182612217565b604051918291518091835e8101838152039060025afa1561104d575f5190565b6159b26159a36132496159986040850185612420565b9190613241866124b4565b602081519101209136906138c9565b6040516159c781615882602082019485613808565b51902090805f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f205482811461552c5780615a3457505f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f2055600190565b90507f657b94fe000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b60206001600160a01b037f2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad9216806001600160a01b03195f516020615f9f5f395f51905f525416175f516020615f9f5f395f51905f5255604051908152a1565b6001600160401b0391615b1a603b60405183819460208301967f4c494748545f434c49454e545f4d49475241544f525f524f4c455f000000000088528484013781015f838201520301601f198101835282612217565b5190201690565b60096122a4916020936001600160c01b0319856040519687948051918291018387015e840101917f0300000000000000000000000000000000000000000000000000000000000000835260c01b1660018201520301601f198101835282612217565b90815115615c3657602091604051615b9b8482612217565b5f8152905f915b8151831015615bf957845f81615bb886866137f4565b51604051918183925191829101835e8101838152039060025afa1561104d57600190615bf15f51916158826040519384928a840161570c565b920191615ba2565b90505f91509291926040516159626021828680820195600160f91b87528051918291018484015e810186838201520301601f198101835282612217565b7f760d6a9b000000000000000000000000000000000000000000000000000000005f5260045ffd5b8051908251808310615ccd578280821091180280831892141582028218906020615c8883856137da565b9280615cac615c9686612238565b95615ca46040519788612217565b808752612238565b8584019690601f190136883703920101835e51902090602081519101201490565b505050505f90565b60ff5f516020615fbf5f395f51905f525460401c1615615cf157565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f5260045ffd5b805182118015615dad575b615d71576001821180615d79575b158015908160011b9182046002141715611b075760280180602811611b07578203615d71576001600160a01b0392915f615d6b92615e40565b90921690565b50505f905f90565b5061060f60f31b7fffff00000000000000000000000000000000000000000000000000000000000060208301511614615d32565b505f615d24565b90615df15750805115615dc957602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b81511580615e37575b615e02575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15615dfa565b92909260018401808511611b0757831180615ef6575b15938415948560011b9586046002141715611b07575f948101809111611b07579192905b818310615e8a5750505060019190565b9092919360ff615ec17fff000000000000000000000000000000000000000000000000000000000000006020888601015116615f2c565b16600f8111615eeb578160041b9180830460101490151715611b0757600191019401919290615e7a565b505f94508493505050565b5061060f60f31b7fffff000000000000000000000000000000000000000000000000000000000000602086840101511614615e56565b60f81c602f811180615f94575b15615f4857602f190160ff1690565b6060811180615f8a575b15615f61576056190160ff1690565b6040811180615f80575b15615f7a576036190160ff1690565b5060ff90565b5060478110615f6b565b5060678110615f52565b50603a8110615f3956fef3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a00f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00a164736f6c634300081c000af0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
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

// GetLightClientMigratorRole is a free data retrieval call binding the contract method 0xb0830ab9.
//
// Solidity: function getLightClientMigratorRole(string clientId) pure returns(uint64)
func (_ContractICS26Router *ContractICS26RouterCaller) GetLightClientMigratorRole(opts *bind.CallOpts, clientId string) (uint64, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "getLightClientMigratorRole", clientId)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetLightClientMigratorRole is a free data retrieval call binding the contract method 0xb0830ab9.
//
// Solidity: function getLightClientMigratorRole(string clientId) pure returns(uint64)
func (_ContractICS26Router *ContractICS26RouterSession) GetLightClientMigratorRole(clientId string) (uint64, error) {
	return _ContractICS26Router.Contract.GetLightClientMigratorRole(&_ContractICS26Router.CallOpts, clientId)
}

// GetLightClientMigratorRole is a free data retrieval call binding the contract method 0xb0830ab9.
//
// Solidity: function getLightClientMigratorRole(string clientId) pure returns(uint64)
func (_ContractICS26Router *ContractICS26RouterCallerSession) GetLightClientMigratorRole(clientId string) (uint64, error) {
	return _ContractICS26Router.Contract.GetLightClientMigratorRole(&_ContractICS26Router.CallOpts, clientId)
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

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ContractICS26Router *ContractICS26RouterCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ContractICS26Router *ContractICS26RouterSession) Paused() (bool, error) {
	return _ContractICS26Router.Contract.Paused(&_ContractICS26Router.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ContractICS26Router *ContractICS26RouterCallerSession) Paused() (bool, error) {
	return _ContractICS26Router.Contract.Paused(&_ContractICS26Router.CallOpts)
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

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ContractICS26Router *ContractICS26RouterSession) Pause() (*types.Transaction, error) {
	return _ContractICS26Router.Contract.Pause(&_ContractICS26Router.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) Pause() (*types.Transaction, error) {
	return _ContractICS26Router.Contract.Pause(&_ContractICS26Router.TransactOpts)
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

// UnfreezeClient is a paid mutator transaction binding the contract method 0xc45f24ca.
//
// Solidity: function unfreezeClient(string clientId) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) UnfreezeClient(opts *bind.TransactOpts, clientId string) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "unfreezeClient", clientId)
}

// UnfreezeClient is a paid mutator transaction binding the contract method 0xc45f24ca.
//
// Solidity: function unfreezeClient(string clientId) returns()
func (_ContractICS26Router *ContractICS26RouterSession) UnfreezeClient(clientId string) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.UnfreezeClient(&_ContractICS26Router.TransactOpts, clientId)
}

// UnfreezeClient is a paid mutator transaction binding the contract method 0xc45f24ca.
//
// Solidity: function unfreezeClient(string clientId) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) UnfreezeClient(clientId string) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.UnfreezeClient(&_ContractICS26Router.TransactOpts, clientId)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ContractICS26Router *ContractICS26RouterSession) Unpause() (*types.Transaction, error) {
	return _ContractICS26Router.Contract.Unpause(&_ContractICS26Router.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) Unpause() (*types.Transaction, error) {
	return _ContractICS26Router.Contract.Unpause(&_ContractICS26Router.TransactOpts)
}

// UpdateApplicationState is a paid mutator transaction binding the contract method 0x9c11bece.
//
// Solidity: function updateApplicationState(string clientId, bytes updateMsg) returns(uint8)
func (_ContractICS26Router *ContractICS26RouterTransactor) UpdateApplicationState(opts *bind.TransactOpts, clientId string, updateMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "updateApplicationState", clientId, updateMsg)
}

// UpdateApplicationState is a paid mutator transaction binding the contract method 0x9c11bece.
//
// Solidity: function updateApplicationState(string clientId, bytes updateMsg) returns(uint8)
func (_ContractICS26Router *ContractICS26RouterSession) UpdateApplicationState(clientId string, updateMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.UpdateApplicationState(&_ContractICS26Router.TransactOpts, clientId, updateMsg)
}

// UpdateApplicationState is a paid mutator transaction binding the contract method 0x9c11bece.
//
// Solidity: function updateApplicationState(string clientId, bytes updateMsg) returns(uint8)
func (_ContractICS26Router *ContractICS26RouterTransactorSession) UpdateApplicationState(clientId string, updateMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.UpdateApplicationState(&_ContractICS26Router.TransactOpts, clientId, updateMsg)
}

// UpdateConsensusState is a paid mutator transaction binding the contract method 0xd395af7f.
//
// Solidity: function updateConsensusState(string clientId, bytes updateMsg) returns(uint8)
func (_ContractICS26Router *ContractICS26RouterTransactor) UpdateConsensusState(opts *bind.TransactOpts, clientId string, updateMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "updateConsensusState", clientId, updateMsg)
}

// UpdateConsensusState is a paid mutator transaction binding the contract method 0xd395af7f.
//
// Solidity: function updateConsensusState(string clientId, bytes updateMsg) returns(uint8)
func (_ContractICS26Router *ContractICS26RouterSession) UpdateConsensusState(clientId string, updateMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.UpdateConsensusState(&_ContractICS26Router.TransactOpts, clientId, updateMsg)
}

// UpdateConsensusState is a paid mutator transaction binding the contract method 0xd395af7f.
//
// Solidity: function updateConsensusState(string clientId, bytes updateMsg) returns(uint8)
func (_ContractICS26Router *ContractICS26RouterTransactorSession) UpdateConsensusState(clientId string, updateMsg []byte) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.UpdateConsensusState(&_ContractICS26Router.TransactOpts, clientId, updateMsg)
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

// ContractICS26RouterICS02ClientUnfrozenIterator is returned from FilterICS02ClientUnfrozen and is used to iterate over the raw logs and unpacked data for ICS02ClientUnfrozen events raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientUnfrozenIterator struct {
	Event *ContractICS26RouterICS02ClientUnfrozen // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterICS02ClientUnfrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterICS02ClientUnfrozen)
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
		it.Event = new(ContractICS26RouterICS02ClientUnfrozen)
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
func (it *ContractICS26RouterICS02ClientUnfrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterICS02ClientUnfrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterICS02ClientUnfrozen represents a ICS02ClientUnfrozen event raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientUnfrozen struct {
	ClientId string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterICS02ClientUnfrozen is a free log retrieval operation binding the contract event 0x7aa8dd0b42a8f6969fbecb42ee2a8fdec2a1a3a576bfc4623b9e7f98ba76fa99.
//
// Solidity: event ICS02ClientUnfrozen(string clientId)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterICS02ClientUnfrozen(opts *bind.FilterOpts) (*ContractICS26RouterICS02ClientUnfrozenIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "ICS02ClientUnfrozen")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterICS02ClientUnfrozenIterator{contract: _ContractICS26Router.contract, event: "ICS02ClientUnfrozen", logs: logs, sub: sub}, nil
}

// WatchICS02ClientUnfrozen is a free log subscription operation binding the contract event 0x7aa8dd0b42a8f6969fbecb42ee2a8fdec2a1a3a576bfc4623b9e7f98ba76fa99.
//
// Solidity: event ICS02ClientUnfrozen(string clientId)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchICS02ClientUnfrozen(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterICS02ClientUnfrozen) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "ICS02ClientUnfrozen")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterICS02ClientUnfrozen)
				if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientUnfrozen", log); err != nil {
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

// ParseICS02ClientUnfrozen is a log parse operation binding the contract event 0x7aa8dd0b42a8f6969fbecb42ee2a8fdec2a1a3a576bfc4623b9e7f98ba76fa99.
//
// Solidity: event ICS02ClientUnfrozen(string clientId)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseICS02ClientUnfrozen(log types.Log) (*ContractICS26RouterICS02ClientUnfrozen, error) {
	event := new(ContractICS26RouterICS02ClientUnfrozen)
	if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientUnfrozen", log); err != nil {
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

// ContractICS26RouterPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ContractICS26Router contract.
type ContractICS26RouterPausedIterator struct {
	Event *ContractICS26RouterPaused // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterPaused)
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
		it.Event = new(ContractICS26RouterPaused)
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
func (it *ContractICS26RouterPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterPaused represents a Paused event raised by the ContractICS26Router contract.
type ContractICS26RouterPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterPaused(opts *bind.FilterOpts) (*ContractICS26RouterPausedIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterPausedIterator{contract: _ContractICS26Router.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterPaused) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterPaused)
				if err := _ContractICS26Router.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_ContractICS26Router *ContractICS26RouterFilterer) ParsePaused(log types.Log) (*ContractICS26RouterPaused, error) {
	event := new(ContractICS26RouterPaused)
	if err := _ContractICS26Router.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ContractICS26RouterUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ContractICS26Router contract.
type ContractICS26RouterUnpausedIterator struct {
	Event *ContractICS26RouterUnpaused // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterUnpaused)
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
		it.Event = new(ContractICS26RouterUnpaused)
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
func (it *ContractICS26RouterUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterUnpaused represents a Unpaused event raised by the ContractICS26Router contract.
type ContractICS26RouterUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterUnpaused(opts *bind.FilterOpts) (*ContractICS26RouterUnpausedIterator, error) {

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterUnpausedIterator{contract: _ContractICS26Router.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterUnpaused) (event.Subscription, error) {

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterUnpaused)
				if err := _ContractICS26Router.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseUnpaused(log types.Log) (*ContractICS26RouterUnpaused, error) {
	event := new(ContractICS26RouterUnpaused)
	if err := _ContractICS26Router.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
