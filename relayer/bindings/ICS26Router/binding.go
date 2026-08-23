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

// ICS02ClientMsgsCounterpartyInfo is an auto generated low-level Go binding around an user-defined struct.
type ICS02ClientMsgsCounterpartyInfo struct {
	ClientId     string
	MerklePrefix [][]byte
}

// ICS26RouterMsgsMsgAckPacket is an auto generated low-level Go binding around an user-defined struct.
type ICS26RouterMsgsMsgAckPacket struct {
	Packet          ICS26RouterMsgsPacket
	Acknowledgement []byte
	MembershipMsg   []byte
}

// ICS26RouterMsgsMsgRecvPacket is an auto generated low-level Go binding around an user-defined struct.
type ICS26RouterMsgsMsgRecvPacket struct {
	Packet        ICS26RouterMsgsPacket
	MembershipMsg []byte
}

// ICS26RouterMsgsMsgSendPacket is an auto generated low-level Go binding around an user-defined struct.
type ICS26RouterMsgsMsgSendPacket struct {
	SourceClient     string
	TimeoutTimestamp uint64
	Payload          ICS26RouterMsgsPayload
}

// ICS26RouterMsgsMsgTimeoutPacket is an auto generated low-level Go binding around an user-defined struct.
type ICS26RouterMsgsMsgTimeoutPacket struct {
	Packet           ICS26RouterMsgsPacket
	NonMembershipMsg []byte
}

// ICS26RouterMsgsPacket is an auto generated low-level Go binding around an user-defined struct.
type ICS26RouterMsgsPacket struct {
	Sequence         uint64
	SourceClient     string
	DestClient       string
	TimeoutTimestamp uint64
	Payloads         []ICS26RouterMsgsPayload
}

// ICS26RouterMsgsPayload is an auto generated low-level Go binding around an user-defined struct.
type ICS26RouterMsgsPayload struct {
	SourcePort string
	DestPort   string
	Version    string
	Encoding   string
	Value      []byte
}

// ContractICS26RouterMetaData contains all meta data concerning the ContractICS26Router contract.
var ContractICS26RouterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"clientMigrationProposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"clientMigrationExecutor\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ackPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structICS26RouterMsgs.MsgAckPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgement\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"membershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addClient\",\"inputs\":[{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addIBCApp\",\"inputs\":[{\"name\":\"app\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addIBCApp\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"app\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelClientMigration\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeClientMigration\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractILightClient\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientMigration\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"executeAfter\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"expireAfter\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitment\",\"inputs\":[{\"name\":\"hashedPath\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCounterparty\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIBCApp\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIIBCApp\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLightClientMigratorRole\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getNextClientSeq\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializeV2\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isConsumingScheduledOp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposeClientMigration\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recvPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structICS26RouterMsgs.MsgRecvPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"membershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structICS26RouterMsgs.MsgSendPacket\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAuthority\",\"inputs\":[{\"name\":\"newAuthority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitMisbehaviour\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"timeoutPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structICS26RouterMsgs.MsgTimeoutPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonMembershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unfreezeClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateApplicationState\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumLightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateConsensusState\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumLightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AckPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgement\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCAppAdded\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"app\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCAppRecvPacketCallbackError\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientAdded\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientMigrated\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientMigrationCancelled\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientMigrationProposed\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"executeAfter\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"expireAfter\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"proposer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientUnfrozen\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientUpdated\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"result\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumLightClientMsgs.UpdateResult\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02MisbehaviourSubmitted\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Noop\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SendPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TimeoutPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WriteAcknowledgement\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgements\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessManagedInvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessManagedRequiredDelay\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"AccessManagedUnauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"DefaultAdminRoleCannotBeGranted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCAppNotFound\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCAsyncAcknowledgementNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCClientAlreadyExists\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationAlreadyProposed\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationDelayRequired\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationDirectCallNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCClientMigrationExpired\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"expireAfter\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationMismatch\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationModuleMismatch\",\"inputs\":[{\"name\":\"module\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expectedModuleId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationModuleMissing\",\"inputs\":[{\"name\":\"module\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationNotProposed\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationNotReady\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"executeAfter\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCClientNotFound\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCCounterpartyClientNotFound\",\"inputs\":[{\"name\":\"counterpartyClientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCCounterpartyMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCErrorUniversalAcknowledgement\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCFailedCallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCInvalidClientId\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidCounterparty\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidCounterpartyClientId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCInvalidLocalClientId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCInvalidPortIdentifier\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidTimeoutDuration\",\"inputs\":[{\"name\":\"maxTimeoutDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualTimeoutDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidTimeoutTimestamp\",\"inputs\":[{\"name\":\"timeoutTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"comparedTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCMultiPayloadPacketNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCPacketAcknowledgementAlreadyExists\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"IBCPacketCommitmentAlreadyExists\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"IBCPacketCommitmentMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IBCPacketReceiptMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IBCPortAlreadyExists\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCUnauthorizedMigrator\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"IBCUnauthorizedSender\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMerklePrefix\",\"inputs\":[{\"name\":\"prefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"NoAcknowledgements\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"Unreachable\",\"inputs\":[]}]",
	Bin: "0x60e080604052346102795760408161637f803803809161001f8285610356565b8339810103126102795761003e602061003783610379565b9201610379565b90803b15610336575f80604051602081019063a1308f2760e01b825260048152610069602482610356565b5190845afa3d1561032e573d906001600160401b03821161028a576040519161009c601f8201601f191660200184610356565b82523d5f602084013e5b15908115610321575b81156102f8575b506102c757813b156102a6575f80604051602081019063a1308f2760e01b8252600481526100e5602482610356565b5190855afa3d1561029e573d906001600160401b03821161028a5760405191610118601f8201601f191660200184610356565b82523d5f602084013e5b1590811561027d575b8115610250575b5061021e576001600160a01b039081166080521660a0523060c0525f51602061635f5f395f51905f5254604081901c60ff1661020f576002600160401b03196001600160401b038216016101b9575b604051615f91908161038e82396080518161035a015260a0518181816104010152611982015260c051818181610cdc0152610e7f0152f35b6001600160401b0319166001600160401b039081175f51602061635f5f395f51905f52556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f610181565b63f92ee8a960e01b5f5260045ffd5b506382f0c90760e01b5f9081526001600160a01b03919091166004525f51602061633f5f395f51905f52602452604490fd5b905060208180518101031261027957602001515f51602061633f5f395f51905f5214155f610132565b5f80fd5b805160201415915061012b565b634e487b7160e01b5f52604160045260245ffd5b606090610122565b50632b833ccb60e11b5f9081526001600160a01b0391909116600452602490fd5b6382f0c90760e01b5f9081526001600160a01b03919091166004525f51602061631f5f395f51905f52602452604490fd5b905060208180518101031261027957602001515f51602061631f5f395f51905f5214155f6100b6565b80516020141591506100af565b6060906100a6565b632b833ccb60e11b5f9081526001600160a01b0391909116600452602490fd5b601f909101601f19168101906001600160401b0382119082101761028a57604052565b51906001600160a01b03821682036102795756fe60806040526004361015610011575f80fd5b5f3560e01c80630cbdf9f61461023f57806319e843991461023a5780631ec43e2314610235578063223e357a146102305780632447af291461022b57806327f146f31461022657806329b6eca9146102215780633f4ba83a1461021c5780634b720d5b146102175780634d6e7ce3146102125780634f1ef2861461020d57806352d1902d14610208578063596e00b9146102035780635c975abb146101fe5780635f516889146101f957806362d42988146101f45780637795820c146101ef5780637a9e5e4b146101ea5780637eb78932146101e55780638456cb59146101e05780638fb36037146101db5780639c11bece146101d65780639e2e5c83146101d1578063ac9650d8146101cc578063ad3cb1cc146101c7578063b0777bfa146101c2578063b0830ab9146101bd578063bf7e214f146101b8578063c45f24ca146101b3578063c4d66de8146101ae578063cce0b2651461019f578063d395af7f146101a9578063e3cb36a0146101a4578063ee5c877c1461019f5763fdbd955d1461019a575f80fd5b611c47565b61193a565b611a1d565b6119b1565b611870565b611790565b61175e565b6116ee565b6116a5565b611646565b6115c6565b61145e565b611343565b611259565b6111c0565b6111a7565b611119565b6110d0565b610fdf565b610f40565b610eff565b610ecf565b610e65565b610c9e565b610932565b6107f6565b610737565b6105e5565b6105a9565b610575565b610528565b610468565b6103bb565b610312565b9181601f84011215610271578235916001600160401b038311610271576020838186019501011161027157565b5f80fd5b908160409103126102715790565b600435906001600160a01b038216820361027157565b602435906001600160a01b038216820361027157565b6060600319820112610271576004356001600160401b03811161027157816102d991600401610244565b92909291602435906001600160401b038211610271576102fb91600401610275565b906044356001600160a01b03811681036102715790565b3461027157610320366102af565b505050507f0cbdf9f6000000000000000000000000000000000000000000000000000000005f5236600480375f8036816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af43d5f803e15610389573d5ff35b3d5ffd5b602060031982011261027157600435906001600160401b038211610271576103b791600401610244565b9091565b34610271576103c93661038d565b50507f19e84399000000000000000000000000000000000000000000000000000000005f5236600480375f8036816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af43d5f803e15610389573d5ff35b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b906020610465928181520190610430565b90565b3461027157610476366102af565b9161048493919336336145fa565b610497610492368385610c4d565b61479b565b156104d6576104b76104be936104d2956104b2368587610c4d565b614ca6565b3691610c4d565b604051918291602083526020830190610430565b0390f35b7fc27d827e000000000000000000000000000000000000000000000000000000005f5260045ffd5b602060031982011261027157600435906001600160401b0382116102715761046591600401610275565b3461027157610550610539366104fe565b610541614dc0565b61054b36336145fa565b612ca6565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d005b3461027157602061058e6105883661038d565b90613057565b6001600160a01b0360405191168152f35b5f91031261027157565b34610271575f3660031901126102715760207f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960254604051908152f35b34610271576020366003190112610271576105fe610283565b5f516020615f655f395f51905f525460ff6001600160401b03821691610626600184146130dd565b60401c1690811561072b575b50610703576106a59061066b60026001600160401b03195f516020615f655f395f51905f525416175f516020615f655f395f51905f5255565b6801000000000000000068ff0000000000000000195f516020615f655f395f51905f525416175f516020615f655f395f51905f5255613121565b6106d468ff0000000000000000195f516020615f655f395f51905f5254165f516020615f655f395f51905f5255565b604051600281527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a1005b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f5260045ffd5b6002915010155f610632565b34610271575f3660031901126102715761075136336145fa565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff8116156107ce5760ff19167fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a1005b7f8dfc202b000000000000000000000000000000000000000000000000000000005f5260045ffd5b346102715760203660031901126102715761080f610283565b610817614dc0565b6001600160a01b03811690816108356108306028615a95565b613dd4565b90603061084183613310565b53607861084d836158f1565b536108586028615aa3565b600181116108ab57506108795761086f925061509b565b610877614e34565b005b7fe22e27eb000000000000000000000000000000000000000000000000000000005f526004839052601460245260445ffd5b90600f811660108110156108ff576108fa917f30313233343536373839616263646566000000000000000000000000000000006108f4921a6108ed8587615901565b5360041c90565b91615abe565b610858565b611ce8565b6020600319820112610271576004356001600160401b03811161027157606090600401809203126102715790565b34610271576104d261094336610904565b61094b6151a4565b610953614dc0565b604081019061099261098061097461058861096e8686611c6f565b80611d0c565b6001600160a01b031690565b33906001600160a01b031681146130e4565b610ac5610abf6109ab6109a58480611d0c565b9061402a565b5192610ab5610a52610aa660208401976109e46109c78a611dd5565b42906109d28c611dd5565b906001600160401b03429116116128ba565b610a2462015180610a0c42610a076109fb8e611dd5565b6001600160401b031690565b613256565b1115610a1e42610a076109fb8e611dd5565b90613263565b610a90610a3a610a348780611d0c565b906151f7565b998a99610a76610a4a8980611d0c565b979093611dd5565b92610a5b61329d565b976104b7610a67610be4565b6001600160401b03909f168f52565b60208c015260408b01526001600160401b031660608a0152565b826080890152610aa1369186611c6f565b6128fb565b610aaf82613310565b52613310565b5061096e84615265565b90612a9b565b7fab3a4458a269be61dfa43faa33aa7b1f5d570716f83ad078bc2ba5dab039abae60405180610afe6001600160401b0387169582613331565b0390a3610b09614e34565b6040516001600160401b0390911681529081906020820190565b634e487b7160e01b5f52604160045260245ffd5b608081019081106001600160401b03821117610b5257604052565b610b23565b604081019081106001600160401b03821117610b5257604052565b606081019081106001600160401b03821117610b5257604052565b602081019081106001600160401b03821117610b5257604052565b60a081019081106001600160401b03821117610b5257604052565b90601f801991011681019081106001600160401b03821117610b5257604052565b60405190610bf360a083610bc3565b565b60405190610bf3608083610bc3565b60405190610bf360e083610bc3565b60405190610bf361010083610bc3565b60405190610bf360c083610bc3565b6001600160401b038111610b5257601f01601f191660200190565b929192610c5982610c32565b91610c676040519384610bc3565b829481845281830111610271578281602093845f960137010152565b9080601f830112156102715781602061046593359101610c4d565b604036600319011261027157610cb2610283565b6024356001600160401b03811161027157610cd1903690600401610c83565b906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016803014908115610e30575b50610e0857610d1636336145fa565b604051917f52d1902d0000000000000000000000000000000000000000000000000000000083526020836004816001600160a01b0386165afa5f9381610dd7575b50610d7857634c9c8ce360e01b5f526001600160a01b03821660045260245ffd5b907f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8303610daa576108779250615aca565b7faa1d49a4000000000000000000000000000000000000000000000000000000005f52600483905260245ffd5b610dfa91945060203d602011610e01575b610df28183610bc3565b810190612500565b925f610d57565b503d610de8565b7fe07c8dba000000000000000000000000000000000000000000000000000000005f5260045ffd5b90506001600160a01b037f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc541614155f610d07565b34610271575f366003190112610271576001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000163003610e085760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b3461027157610550610ee0366104fe565b610ee86151a4565b610ef0614dc0565b610efa36336145fa565b613903565b34610271575f36600319011261027157602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b34610271576040366003190112610271576004356001600160401b03811161027157610fda610f76610550923690600401610244565b9190610f80610299565b92610f89614dc0565b610f9336336145fa565b610fa08183811515613cc1565b610fc28183610fbb610fb3368484610c4d565b805190615c9c565b5015613cc1565b6104b78183610fd5610492368484610c4d565b613cc1565b61509b565b34610271576020610fef3661038d565b919082604051938492833781017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496038152030190206040519061103082610b37565b6104d2606061109761108a61107e600186549687895201549665ffffffffffff88169081602082015265ffffffffffff8960301c1698896040830152861c958691015265ffffffffffff1690565b9565ffffffffffff1690565b916001600160a01b031690565b9060405194859485929365ffffffffffff6001600160a01b03929695816060956080880199885216602087015216604085015216910152565b34610271576020366003190112610271576004355f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600602052602060405f2054604051908152f35b3461027157602036600319011261027157611132610283565b6001600160a01b035f516020615f455f395f51905f525416330361119557803b15611160576108779061569b565b6001600160a01b03907fc2f31e5e000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b62d1953b60e31b5f523360045260245ffd5b3461027157602061058e6111ba3661038d565b90613d09565b34610271575f366003190112610271576111da36336145fa565b6111e26151a4565b600160ff197fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005416177fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a1005b34610271575f366003190112610271575f516020615f455f395f51905f525460a01c60ff16156112b95760207f8fb36037000000000000000000000000000000000000000000000000000000005b6001600160e01b031960405191168152f35b60205f6112a7565b6040600319820112610271576004356001600160401b03811161027157816112eb91600401610244565b92909291602435906001600160401b038211610271576103b791600401610244565b634e487b7160e01b5f52602160045260245ffd5b6003111561132b57565b61130d565b9190602083019261134082611321565b52565b34610271576113af6020611356366112c1565b906113659492939436336145fa565b6001600160a01b036113778587613d09565b16905f6040518098819582947f4b1872be00000000000000000000000000000000000000000000000000000000845260048401613046565b03925af1908115611459576104d2935f92611408575b5081926113f97f87bbef2779889a19f0435ddca81fda94132c06ffddb0ea73def256307a293aef9360405193849384613da4565b0390a160405191829182611330565b7f87bbef2779889a19f0435ddca81fda94132c06ffddb0ea73def256307a293aef925061144c9060203d602011611452575b6114448183610bc3565b810190613d8f565b916113c5565b503d61143a565b6128af565b346102715761146c366112c1565b9192905f9261147b36336145fa565b6001600160a01b0361148d8685613d09565b1691823b15610271576114d9925f92836040518096819582947fddba65370000000000000000000000000000000000000000000000000000000084526020600485018181520191611d3e565b03925af1801561145957611529575b507fa263f0a976b2937a51fd2e416491cf0ca724d5499fa870715929dfde4ee4a4309192611523604051928392602084526020840191611d3e565b0390a180f35b7fa263f0a976b2937a51fd2e416491cf0ca724d5499fa870715929dfde4ee4a43092505f61155691610bc3565b5f916114e8565b9080602083519182815201916020808360051b8301019401925f915b83831061158857505050505090565b90919293946020806115a6600193601f198682030187528951610430565b97019301930191939290611579565b90602061046592818152019061155d565b34610271576020366003190112610271576004356001600160401b0381116102715736602382011215610271578060040135906001600160401b038211610271573660248360051b83010111610271576104d29160246116269201613e48565b604051918291826115b5565b60405190611641602083610bc3565b5f8252565b34610271575f366003190112610271576104d2604051611667604082610bc3565b600581527f352e302e300000000000000000000000000000000000000000000000000000006020820152604051918291602083526020830190610430565b34610271576104d26116b96109a53661038d565b6040519182916020835260206116da82516040838701526060860190610430565b910151838203601f1901604085015261155d565b346102715760206001600160401b036117063661038d565b611752603b604051838194888301967f4c494748545f434c49454e545f4d49475241544f525f524f4c455f000000000088528484013781015f838201520301601f198101835282610bc3565b51902016604051908152f35b34610271575f3660031901126102715760206001600160a01b035f516020615f455f395f51905f525416604051908152f35b346102715761179e3661038d565b905f906117ab36336145fa565b6001600160a01b036117bd8483613d09565b16803b15610271575f80916004604051809481937f6a28f0000000000000000000000000000000000000000000000000000000000083525af180156114595761183c575b507f7aa8dd0b42a8f6969fbecb42ee2a8fdec2a1a3a576bfc4623b9e7f98ba76fa999192611523604051928392602084526020840191611d3e565b7f7aa8dd0b42a8f6969fbecb42ee2a8fdec2a1a3a576bfc4623b9e7f98ba76fa9992505f61186991610bc3565b5f91611801565b3461027157602036600319011261027157611889610283565b5f516020615f655f395f51905f525460ff6001600160401b038216916118af83156130dd565b60401c1690811561192e575b50610703576106a5906118f460026001600160401b03195f516020615f655f395f51905f525416175f516020615f655f395f51905f5255565b6801000000000000000068ff0000000000000000195f516020615f655f395f51905f525416175f516020615f655f395f51905f5255614112565b6002915010155f6118bb565b3461027157611948366102af565b505050507fee5c877c000000000000000000000000000000000000000000000000000000005f5236600480375f8036816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af43d5f803e15610389573d5ff35b34610271576113af60206119c4366112c1565b906119d39492939436336145fa565b6001600160a01b036119e58587613d09565b16905f6040518098819582947f92c19bdc00000000000000000000000000000000000000000000000000000000845260048401613046565b34610271576040366003190112610271576004356001600160401b03811161027157611a4d903690600401610275565b611a55610299565b90611a5e614760565b7f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496025490925f198214611c4257600182017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496025581825f937a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000811015611c17575b50806d04ee2d6d415b85acef8100000000600a921015611bfb575b662386f26fc10000811015611be6575b6305f5e100811015611bd4575b612710811015611bc4575b6064811015611bb5575b1015611baa575b611b756021611b3d60018601613dd4565b948501015b5f1901917f3031323334353637383961626364656600000000000000000000000000000000600a82061a8353600a900490565b908115611b8557611b7590611b42565b5050611b97611b9e926104d295615718565b9283614ca6565b60405191829182610454565b600190920191611b2c565b60029060649004940193611b25565b6004906127109004940193611b1b565b6008906305f5e1009004940193611b10565b601090662386f26fc100009004940193611b03565b6020906d04ee2d6d415b85acef81000000009004940193611af3565b604094507a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000090049050600a611ad8565b613234565b3461027157610550611c5836610904565b611c60614dc0565b611c6a36336145fa565b6141e5565b903590609e1981360301821215610271570190565b903590601e198136030182121561027157018035906001600160401b03821161027157602001918160051b3603831361027157565b15611cc057565b7f356f4dbd000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52603260045260245ffd5b90156108ff578061046591611c6f565b903590601e198136030182121561027157018035906001600160401b0382116102715760200191813603831361027157565b908060209392818452848401375f828201840152601f01601f1916010190565b92909215611d6b57505050565b611dae9291611dc0916040519485947f9fff831f000000000000000000000000000000000000000000000000000000008652604060048701526044860190610430565b84810360031901602486015291611d3e565b0390fd5b6001600160401b0381160361027157565b3561046581611dc4565b3590610bf382611dc4565b919082604091031261027157604051611e0281610b57565b60208082948035611e1281611dc4565b8452013591611e2083611dc4565b0152565b6001600160401b038111610b525760051b60200190565b9080601f83011215610271578135611e5281611e24565b92611e606040519485610bc3565b81845260208085019260051b820101918383116102715760208201905b838210611e8c57505050505090565b81356001600160401b03811161027157602091611eae87848094880101610c83565b815201910190611e7d565b9080601f8301121561027157813591611ed183611e24565b92611edf6040519485610bc3565b80845260208085019160051b830101918383116102715760208101915b838310611f0b57505050505090565b82356001600160401b038111610271578201906040828703601f1901126102715760405190611f3982610b57565b60208301356001600160401b03811161027157876020611f5b92860101611e3b565b82526040830135916001600160401b03831161027157611f8388602080969581960101610c83565b83820152815201920191611efc565b6003111561027157565b91906080838203126102715760405190611fb582610b37565b81938035611fc281611f92565b83526020810135611fd281611f92565b60208401526040810135611fe581611f92565b60408401526060810135916001600160401b0383116102715760609261200b9201610c83565b910152565b9080601f830112156102715781359161202883611e24565b926120366040519485610bc3565b80845260208085019160051b830101918383116102715760208101915b83831061206257505050505090565b82356001600160401b038111610271578201906060828703601f190112610271576040519061209082610b72565b602083013561209e81611f92565b825260408301356001600160401b038111610271578760206120c292860101610c83565b60208301526060830135916001600160401b038311610271576120ed88602080969581960101610c83565b6040820152815201920191612053565b91909160808184031261027157612112610bf5565b9281356001600160401b038111610271578161212f918401610c83565b845260208201356001600160401b0381116102715781612150918401610c83565b602085015260408201356001600160401b0381116102715781612174918401611f9c565b604085015260608201356001600160401b038111610271576121969201612010565b6060830152565b3590811515820361027157565b9080601f83011215610271578135916121c283611e24565b926121d06040519485610bc3565b80845260208085019160051b830101918383116102715760208101915b8383106121fc57505050505090565b82356001600160401b0381116102715782016020818703601f190112610271576040519061222982610b8d565b60208101356001600160401b03811161027157602091010186601f820112156102715780359061225882611e24565b916122666040519384610bc3565b80835260208084019160051b830101918983116102715760208101915b83831061229f57505050908252508152602092830192016121ed565b82356001600160401b038111610271578201906060828d03601f19011261027157604051906122cd82610b72565b6020830135600281101561027157825260408301356001600160401b038111610271578d60206122ff928601016120fd565b602083015260608301356001600160401b03811161027157602060a0918f9501018094031261027157612330610be4565b9183356001600160401b038111610271578e61234d918601610c83565b835261235b6020850161219d565b602084015260408401356001600160401b038111610271578e61237f9186016120fd565b60408401526123906060850161219d565b60608401526080840135926001600160401b038411610271576123b98f602096958796016120fd565b60808201526040820152815201920191612283565b9190826060910312610271576040516123e681610b72565b809280356fffffffffffffffffffffffffffffffff811681036102715760409182918452602081013560208501520135910152565b3590600182101561027157565b602081830312610271578035906001600160401b03821161027157016101408183031261027157612457610c04565b916124628183611dea565b835260408201356001600160401b0381116102715781612483918401611eb9565b602084015260608201356001600160401b03811161027157816124a79184016121aa565b6040840152608082013560608401526124c38160a084016123ce565b60808401526124d5610100830161241b565b60a08401526101208201356001600160401b038111610271576124f89201611e3b565b60c082015290565b90816020910312610271575190565b9080602083519182815201916020808360051b8301019401925f915b83831061253a57505050505090565b9091929394602080612578600193601f1986820301875289519083612568835160408452604084019061155d565b9201519084818403910152610430565b9701930193019193929061252b565b9080602083519182815201916020808360051b8301019401925f915b8383106125b257505050505090565b9091929394602080612603600193601f1986820301875289519081516125d781611321565b815260406125f2858401516060878501526060840190610430565b920151906040818403910152610430565b970193019301919392906125a3565b6104659160606126956126436126318551608086526080860190610430565b60208601518582036020870152610430565b60808360408701518684036040880152805161265e81611321565b8452602081015161266e81611321565b6020850152604081015161268181611321565b604085015201519181858201520190610430565b920151906060818403910152612587565b9080602083519283815201916020808260051b8401019401925f925b8284106126d157505050505090565b9091929394601f198282030183528551602082019051916020815282518092526040810190602060408460051b8301019401925f915b818310612728575050505050602080600192970193019401929193906126c2565b9091929394603f19828203018552855190815191600283101561132b576127c98260209392600195859452604061276c858301516060878601526060850190612612565b91015191604081830391015260806127ac612790845160a0855260a0850190610430565b8685015115158785015260408501518482036040860152612612565b926060810151151560608401520151906080818403910152612612565b9701950193019190612707565b90600182101561132b5752565b90610465916020815261280f6020820183516001600160401b0360208092828151168552015116910152565b60c061284561282f6020850151610140606086015261016085019061250f565b6040850151848203601f190160808601526126a6565b92606081015160a0840152612887608082015183850190604080916fffffffffffffffffffffffffffffffff8151168452602081015160208501520151910152565b61289a60a08201516101208501906127d6565b015190610140601f198285030191015261155d565b6040513d5f823e3d90fd5b156128c3575050565b6001600160401b03907f65d30129000000000000000000000000000000000000000000000000000000005f521660045260245260445ffd5b91909160a08184031261027157612910610be4565b9281356001600160401b038111610271578161292d918401610c83565b845260208201356001600160401b038111610271578161294e918401610c83565b602085015260408201356001600160401b0381116102715781612972918401610c83565b604085015260608201356001600160401b0381116102715781612996918401610c83565b606085015260808201356001600160401b038111610271576129b89201610c83565b6080830152565b610465916080612a17612a056129f36129e1865160a0875260a0870190610430565b60208701518682036020880152610430565b60408601518582036040870152610430565b60608501518482036060860152610430565b920151906080818403910152610430565b6020815260a06001600160a01b036080612a91612a69612a53875186602089015260c0880190610430565b6020880151878203601f19016040890152610430565b6001600160401b0360408801511660608701526060870151601f1987830301848801526129bf565b9401511691015290565b81604051928392833781015f815203902090565b9035601e19823603018112156102715701602081359101916001600160401b03821161027157813603831361027157565b9035601e19823603018112156102715701602081359101916001600160401b038211610271578160051b3603831361027157565b906001600160401b038235612b2881611dc4565b168152612b94612b6d612b52612b416020860186612aaf565b60a0602087015260a0860191611d3e565b612b5f6040860186612aaf565b908583036040870152611d3e565b926001600160401b036060820135612b8481611dc4565b1660608401526080810190612ae0565b9290916080818303910152828152602081019260208160051b83010193835f91609e1982360301945b848410612bce575050505050505090565b90919293949596601f19828203018352873587811215610271576020612c846001938783940190612c76612c6b612c50612c35612c1c612c0e8780612aaf565b60a0885260a0880191611d3e565b612c2889880188612aaf565b908783038b890152611d3e565b612c426040870187612aaf565b908683036040880152611d3e565b612c5d6060860186612aaf565b908583036060870152611d3e565b926080810190612aaf565b916080818503910152611d3e565b990193019401929195949390612bbd565b906020610465928181520190612b14565b612cca6001612cc2612cb88480611c6f565b6080810190611c84565b905014611cb9565b612ce0612cda612cb88380611c6f565b90611cfc565b905f6020612df2612d006109a5612cf78680611c6f565b84810190611d0c565b612d448151848151910120612d256104b7612d1b8980611c6f565b6040810190611d0c565b85815191012014825190612d3c612d1b8980611c6f565b929091611d5e565b612d9b612d7b612d76612d5a612d1b8980611c6f565b9190612d6e612d698b80611c6f565b611dd5565b923691610c4d565b614e59565b84612d93612d8b828a018a611d0c565b810190612428565b930151614e9c565b60c0820152612dbc6109746111ba612db38880611c6f565b86810190611d0c565b906040519485809481937fa6f031bb000000000000000000000000000000000000000000000000000000008352600483016127e3565b03925af1801561145957612e38915f91612fef575b50612e206109fb6060612e1a8680611c6f565b01611dd5565b811015612e326060612e1a8680611c6f565b906128ba565b612e51612e4d612e488380611c6f565b614f8c565b1590565b612fc75781612edc612e6c61097461058884612ed597611d0c565b91612e84612e7a8580611c6f565b6020810190611d0c565b9590612ec0612e96612d1b8880611c6f565b612eb7612ea6612d698b80611c6f565b94612eaf610be4565b9b3691610c4d565b8a523691610c4d565b60208801526001600160401b03166040870152565b36906128fb565b6060840152336080840152803b1561027157612f2b5f939184926040519586809481937f5e32b6b600000000000000000000000000000000000000000000000000000000835260048301612a28565b03925af191821561145957612f9692612fad575b507f01e5ed58494819ef3f6480dd08e433b7c08ed75c7abdf2c22c6f04b71340a1686001600160401b03612f76612e7a8480611c6f565b612f90612f89612d698780999599611c6f565b9580611c6f565b95612a9b565b92612fa8604051928392169582612c95565b0390a3565b80612fbb5f612fc193610bc3565b8061059f565b5f612f3f565b50507fd08bf58b0e4eec5bfc697a4fdbb6839057fbf4dd06f1b1ce07445c0e5a654caf5f80a1565b613008915060203d602011610e0157610df28183610bc3565b5f612e07565b60209082604051938492833781017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960181520301902090565b916020610465938181520191611d3e565b6001600160a01b03604051838382376020818581017fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf60081520301902054169182156130a157505090565b611dc06040519283927fa09dbf590000000000000000000000000000000000000000000000000000000084526020600485018181520191611d3e565b1561070357565b156130ec5750565b6001600160a01b03907fbe2f2b45000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b610bf3906001600160a01b03197fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b00546001600160a01b038116331480156131fc575b61316e9033906130e4565b167fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b00556001600160a01b03197fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b0154167fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b01556131e7615a51565b6131ef615a51565b6131f7615a51565b61569b565b5061316e6001600160a01b037fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b01541633149050613163565b634e487b7160e01b5f52601160045260245ffd5b5f19810191908211611c4257565b91908203918211611c4257565b1561326b5750565b7f715fed60000000000000000000000000000000000000000000000000000000005f526201518060045260245260445ffd5b6132a76001611e24565b906132b56040519283610bc3565b6001825281601f196132c76001611e24565b01905f5b8281106132d757505050565b6020906040516132e681610ba8565b606081526060838201526060604082015260608082015260606080820152828285010152016132cb565b8051156108ff5760200190565b80518210156108ff5760209160051b010190565b90602082526001600160401b038151166020830152608061337a613364602084015160a0604087015260c0860190610430565b6040840151858203601f19016060870152610430565b916001600160401b036060820151168285015201519160a0601f1982840301910152815180825260208201916020808360051b8301019401925f915b8383106133c557505050505090565b90919293946020806133e3600193601f1986820301875289516129bf565b970193019301919392906133b6565b9080601f8301121561027157813561340981611e24565b926134176040519485610bc3565b81845260208085019260051b820101918383116102715760208201905b83821061344357505050505090565b81356001600160401b03811161027157602091613465878480948801016128fb565b815201910190613434565b91909160a08184031261027157613485610be4565b9261348f82611ddf565b845260208201356001600160401b03811161027157816134b0918401610c83565b602085015260408201356001600160401b03811161027157816134d4918401610c83565b60408501526134e560608301611ddf565b606085015260808201356001600160401b038111610271576129b892016133f2565b610465903690613470565b602081830312610271578035906001600160401b03821161027157016101608183031261027157613541610c13565b9161354c8183611dea565b835260408201356001600160401b038111610271578161356d918401611eb9565b602084015260608201356001600160401b03811161027157816135919184016121aa565b6040840152608082013560608401526135ad8160a084016123ce565b60808401526135bf610100830161241b565b60a08401526101208201356001600160401b03811161027157816135e4918401611e3b565b60c08401526101408201356001600160401b038111610271576136079201610c83565b60e082015290565b90610465916020815261363b6020820183516001600160401b0360208092828151168552015116910152565b60e06136da61367461365e6020860151610160606087015261018086019061250f565b6040860151858203601f190160808701526126a6565b606085015160a0850152608085015180516fffffffffffffffffffffffffffffffff1660c0860152602081015160e0860152604001516101008501526136c360a08601516101208601906127d6565b60c0850151848203601f190161014086015261155d565b92015190610160601f1982850301910152610430565b604080519091906137018382610bc3565b6001815291601f1901825f5b82811061371957505050565b80606060208093850101520161370d565b9061373482611e24565b6137416040519182610bc3565b8281528092613752601f1991611e24565b01905f5b82811061376257505050565b806060602080938501015201613756565b602081830312610271578051906001600160401b038211610271570181601f82011215610271578051906137a682610c32565b926137b46040519485610bc3565b8284526020838301011161027157815f9260208093018386015e8301015290565b156137dc57565b7fecfef798000000000000000000000000000000000000000000000000000000005f5260045ffd5b60405190613813604083610bc3565b602082527f4774d4a575993f963b1c06573736617a457abef8589178db8d10c94b4ab511ab6020830152565b613847613804565b6020815191012090565b1561385857565b7f6b2675e3000000000000000000000000000000000000000000000000000000005f5260045ffd5b3d156138aa573d9061389182610c32565b9161389f6040519384610bc3565b82523d5f602084013e565b606090565b156138b657565b7fadef7fb8000000000000000000000000000000000000000000000000000000005f5260045ffd5b90916138f561046593604084526040840190612b14565b91602081840391015261155d565b6139156001612cc2612cb88480611c6f565b613925612cda612cb88380611c6f565b5f6020613a6761393b6109a5612d1b8780611c6f565b613976815184815191012061395f6104b76139568a80611c6f565b87810190611d0c565b85815191012014825190612d3c6139568a80611c6f565b6139996139886060612e1a8980611c6f565b42906109d26060612e1a8b80611c6f565b613a24613a326139cf6139ca6139bb6139b28b80611c6f565b88810190611d0c565b9190612d6e612d698d80611c6f565b615346565b613a0d6139ec6139e76139e28c80611c6f565b613507565b6153b9565b9187613a056139fd8d83018e611d0c565b810190613512565b960151614e9c565b60c085015260405192839187830160209181520190565b03601f198101835282610bc3565b60e0820152613a4a6109746111ba612d1b8980611c6f565b906040519485809481936325d29d3160e21b83526004830161360f565b03925af1801561145957613ca4575b50613a8c612e4d613a878480611c6f565b6154ef565b612fc7577f76765590e2b799b0506100f8a6610cfecab2c71e8e1f8aa981b099aff0dfdb746001600160401b03613c33935f80613b86613aca6136f0565b96612ed5613b46613ae46109746105886020860186611d0c565b92613af2612e7a8980611c6f565b9390613b318a612d69613b28613b17613b0e612d1b8580611c6f565b93909480611c6f565b94613b20610be4565b993691610c4d565b88523691610c4d565b60208601526001600160401b03166040850152565b60608201523360808201526040519485809481937f078c4a7900000000000000000000000000000000000000000000000000000000835260048301612a28565b03925af15f9181613c80575b50613c4257507fb9edb487876e8be10f54e377c1a815a54ad92a6db1c9561dfe8fad2f0d1da84f613bd0613bc4613880565b611b9e815115156138af565b0390a1613bdb613804565b613be485613310565b52613bee84613310565b505b613c0384613bfe8380611c6f565b6155d5565b612fa8613c13612d1b8380611c6f565b613c2d613c26612d6986809b959b611c6f565b9480611c6f565b97612a9b565b946040519384931696836138de565b613c4e815115156137d5565b613c6781516020830120613c6061383f565b1415613851565b613c7085613310565b52613c7a84613310565b50613bf0565b613c9d9192503d805f833e613c958183610bc3565b810190613773565b905f613b92565b613cbc9060203d602011610e0157610df28183610bc3565b613a76565b91909115613ccd575050565b611dc06040519283927f14d712470000000000000000000000000000000000000000000000000000000084526020600485018181520191611d3e565b6001600160a01b03604051838382376020818581017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496008152030190205416918215613d5357505090565b611dc06040519283927fa0db16fe0000000000000000000000000000000000000000000000000000000084526020600485018181520191611d3e565b90816020910312610271575161046581611f92565b939291602091613dbc91604087526040870191611d3e565b93611e2083611321565b906004116102715790600490565b90613dde82610c32565b613deb6040519182610bc3565b8281528092613dfc601f1991610c32565b0190602036910137565b805191908290602001825e015f815290565b613e3a610bf392949360208660405197889583870137840101905f8252613e06565b03601f198101845283610bc3565b90919060405190613e5a602083610bc3565b5f8252601f19613e695f610c32565b01366020840137613e798461372a565b925f5b85811015613ec557600190613ea9613ea386613e9d8460051b880188611d0c565b90613e18565b306156fb565b613eb3828861331d565b52613ebe818761331d565b5001613e7c565b509350505090565b90600182811c92168015613efb575b6020831014613ee757565b634e487b7160e01b5f52602260045260245ffd5b91607f1691613edc565b90815491613f1283611e24565b92613f206040519485610bc3565b80845260208401915f5260205f20915f905b828210613f3f5750505050565b6040515f8554613f4e81613ecd565b8084529060018116908115613fbf5750600114613f88575b5060019282613f7a85946020940382610bc3565b815201940191019092613f32565b5f878152602081209092505b818310613fa957505081016020016001613f66565b6001816020925483868801015201920191613f94565b60ff191660208581019190915291151560051b8401909101915060019050613f66565b91909115613fee575050565b611dc06040519283927fdf95155a0000000000000000000000000000000000000000000000000000000084526020600485018181520191611d3e565b9190916060602060405161403d81610b57565b828152015261404c838261300e565b926040519161405a83610b57565b604051945f815461406a81613ecd565b80895290600181169081156140ee57506001146140b2575b506140a3918761409b6104659798996001940382610bc3565b875201613f05565b60208501528351511515613fe2565b5f838152602081209092505b8183106140d457505086016020016140a3614082565b6001816020929493945483858d01015201910191906140be565b60ff19166020808b019190915291151560051b890190910191506140a39050614082565b610bf39061411e615a51565b614126615a51565b61412e615a51565b614136615a51565b6131e7615a51565b6020815260c06001600160a01b0360a0612a916141ab61418261416c88518760208a015260e0890190610430565b6020890151888203601f190160408a0152610430565b6001600160401b0360408901511660608801526060880151601f198883030160808901526129bf565b6080870151868203601f190184880152610430565b916141d76104659492604085526040850190612b14565b926020818503910152611d3e565b6141ef8180611c6f565b608081016141fc91611c84565b6142099150600114611cb9565b6142138180611c6f565b6080810161422091611c84565b61422991611cfc565b6142338280611c6f565b6020810161424091611d0c565b6142499161402a565b9081518051906020012061425d8480611c6f565b6040810161426a91611d0c565b369061427592610c4d565b805190602001201482516142898580611c6f565b6040810161429691611d0c565b916142a093611d5e565b6142aa8380611c6f565b604081016142b791611d0c565b906142c28580611c6f565b6142cb90611dd5565b9136906142d792610c4d565b906142e191615734565b916142ea6136f0565b9260208501936142fa8587611d0c565b369061430592610c4d565b61430e82613310565b5261431881613310565b50614322906157d7565b906143306040870187611d0c565b810161433b91613512565b92602001519061434a91614e9c565b60c0830152604080516020808201939093529182526143699082610bc3565b60e08201526143788480611c6f565b6020810161438591611d0c565b61438e91613d09565b6001600160a01b03166040518080936325d29d3160e21b825260048201906143b59161360f565b03815a6020945f91f1801561145957614587575b506143da612e4d612e488580611c6f565b61455e57806144716143f56109746105888461446596611d0c565b91614447614406612e7a8880611c6f565b959092612ed58961445061441d612d1b8380611c6f565b6144348d61442e612d698780611c6f565b95611d0c565b98909961443f610c23565b9d3691610c4d565b8c523691610c4d565b60208a01526001600160401b03166040890152565b60608601523691610c4d565b60808301523360a0830152803b15610271576144c05f929183926040519485809481937f428e4e170000000000000000000000000000000000000000000000000000000083526004830161413e565b03925af1801561145957612fa8927ff9bab74bcdb634f4d3dd064cc42a13df056598e1c0336905d2f5750fbfb08b7b926001600160401b039261454a575b5061450c612e7a8680611c6f565b61453b614533614522612d698a809a969a611c6f565b9461452d8a80611c6f565b99611d0c565b929097612a9b565b956040519485941697846141c0565b80612fbb5f61455893610bc3565b5f6144fe565b5050507fd08bf58b0e4eec5bfc697a4fdbb6839057fbf4dd06f1b1ce07445c0e5a654caf5f80a1565b61459f9060203d602011610e0157610df28183610bc3565b6143c9565b919091356001600160e01b0319811692600481106145c0575050565b6001600160e01b0319929350829060040360031b1b161690565b6040906001600160a01b0361046595931681528160208201520191611d3e565b61463761461b5f516020615f455f395f51905f52546001600160a01b031690565b61462e614628855f613dc6565b906145a4565b90833091615864565b901561464257505050565b63ffffffff16156147465761468e7401000000000000000000000000000000000000000060ff60a01b195f516020615f455f395f51905f525416175f516020615f455f395f51905f5255565b6146b26109746109745f516020615f455f395f51905f52546001600160a01b031690565b91823b15610271576146f8925f808094604051968795869485937f94c7d7ee000000000000000000000000000000000000000000000000000000008552600485016145da565b03925af1801561145957614732575b50610bf360ff60a01b195f516020615f455f395f51905f5254165f516020615f455f395f51905f5255565b80612fbb5f61474093610bc3565b5f614707565b62d1953b60e31b5f526001600160a01b031660045260245ffd5b6040519061476f604083610bc3565b600782527f636c69656e742d000000000000000000000000000000000000000000000000006020830152565b6147a481615912565b15614819576147ea6040516147ba604082610bc3565b600881527f6368616e6e656c2d0000000000000000000000000000000000000000000000006020820152826159ff565b908115614800575b506147fc57600190565b5f90565b614813915061480d614760565b906159ff565b5f6147f2565b505f90565b60206148309160405192838092613e06565b7f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960081520301902090565b602061486c9160405192838092613e06565b7f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960181520301902090565b1561489d57565b7fd7274bcf000000000000000000000000000000000000000000000000000000005f5260045ffd5b8181106148d0575050565b5f81556001016148c5565b9190601f81116148ea57505050565b610bf3925f5260205f20906020601f840160051c83019310614914575b601f0160051c01906148c5565b9091508190614907565b908160011b9180830460021490151715611c4257565b908160041b9180830460101490151715611c4257565b9092916001600160401b038111610b525761496f816149698454613ecd565b846148db565b5f601f82116001146149ad57819061499e9394955f926149a2575b50508160011b915f199060031b1c19161790565b9055565b013590505f8061498a565b601f198216946149c0845f5260205f2090565b915f5b8781106149fa5750836001959697106149e1575b505050811b019055565b01355f19600384901b60f8161c191690555f80806149d7565b909260206001819286860135815501940191016149c3565b680100000000000000008311610b52578054838255808410614a77575b5090614a3f81925f5260205f2090565b905f925b848410614a51575050505050565b6001602082614a6b614a64849587611d0c565b908861494a565b01930193019291614a43565b815f528360205f2091820191015b818110614a925750614a2f565b80614a9f60019254613ecd565b80614aac575b5001614a85565b601f81118314614ac157505f81555b5f614aa5565b614ae39083601f614ad5855f5260205f2090565b920160051c820191016148c5565b5f8181526020812081835555614abb565b90614aff8180611d0c565b906001600160401b038211610b5257614b2282614b1c8654613ecd565b866148db565b5f90601f8311600114614b705792614b5a83610bf39694614b67946001975f926149a25750508160011b915f199060031b1c19161790565b83555b6020810190611c84565b92909101614a12565b601f19831691614b83865f5260205f2090565b925f5b818110614bcb575093614b679360019693879383610bf39a9810614bb2575b505050811b018355614b5d565b01355f19600384901b60f8161c191690555f8080614ba5565b91936020600181928787013581550195019201614b86565b93929190614bf990606086526060860190610430565b8481036020860152614c2a614c1f614c118480612aaf565b604085526040850191611d3e565b926020810190612ae0565b90916020818503910152808352602083019260208260051b82010193835f925b848410614c6e57505050505050906040610bf3929401906001600160a01b03169052565b909192939495602080614c96600193601f19868203018852614c908b88612aaf565b90611d3e565b9801940194019294939190614c4a565b92916001600160a01b03604051602081614cc08189613e06565b7f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496008152030190205416614d855792614d807f0ecded31ecd211a73abf0fb3bc09150bbe321a05550fbe29ea0f16b6e25fbfa89394614d31614d2c614d276104b78780611d0c565b615912565b614896565b614d626001600160a01b038216614d478561481e565b906001600160a01b03166001600160a01b0319825416179055565b614d7484614d6f8561485a565b614af4565b60405193849384614be3565b0390a1565b6040517f87dfb2670000000000000000000000000000000000000000000000000000000081526020600482015280611dc06024820187610430565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005c614e0c5760017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d565b7f3ee5aeb5000000000000000000000000000000000000000000000000000000005f5260045ffd5b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d565b600961046591614e7c936001600160c01b03196040519586936020850190613e06565b91600160f91b835260c01b1660018201520301601f198101835282610bc3565b919091805115614f5657614eb0815161372a565b915f5b614ebd8351613248565b811015614eee5780614ed16001928561331d565b51614edc828761331d565b52614ee7818661331d565b5001614eb3565b5092614f40614f5292614f3a614f2693614f2c614f14614f0e8551613248565b8561331d565b51916040519687936020850190613e06565b90613e06565b03601f198101855284610bc3565b51613248565b90614f4b828561331d565b528261331d565b5090565b611dc0906040519182917fa7c34e4f000000000000000000000000000000000000000000000000000000008352600483016115b5565b614fad6139ca614f9f6020840184611d0c565b91908435926104b784611dc4565b6020815191012090815f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f2054801561506d576150046139e7614ffa6139e73686613470565b8314933690613470565b911561503f5750505f9081527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600602052604081205b55600190565b7f3f87a2ec000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b5050505f90565b906001600160a01b03615094602092959495604085526040850190610430565b9416910152565b91906001600160a01b036040516020816150b58188613e06565b7fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf6008152030190205416615169577fa6ec8e860960e638347460dc632fbe0175c51a5ca130e336138bbe26ff304499919261515a6040516020816151188186613e06565b7fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf60081520301902080546001600160a01b0319166001600160a01b038516179055565b614d8060405192839283615074565b6040517f837f46a60000000000000000000000000000000000000000000000000000000081526020600482015280611dc06024820186610430565b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166151cf57565b7fd93c0665000000000000000000000000000000000000000000000000000000005f5260045ffd5b60209082604051938492833781017f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e07476018152030190206001600160401b038154166001600160401b038114611c4257600101906001600160401b0382166001600160401b031982541617905590565b602081019061528082516001600160401b0383511690615346565b6020815191012091825f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e07476006020526152cb60405f20541591516001600160401b0384511690615346565b901561530957506152db906153b9565b905f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f2055565b611dc0906040519182917f91ffd9240000000000000000000000000000000000000000000000000000000083526020600484018181520190610430565b600961046591615369936001600160c01b03196040519586936020850190613e06565b917f0100000000000000000000000000000000000000000000000000000000000000835260c01b1660018201520301601f198101835282610bc3565b6153b29060209392613e06565b9081520190565b906153c2611632565b5f905b6080840151805183101561540757906153ff6153ec6153e68560019561331d565b51615bb0565b91613a24604051938492602084016153a5565b9101906153c5565b50905091909160205f615424604085015160405191828092613e06565b039060025afa156114595760205f61547c613a24615470615451606085519801516001600160401b031690565b6040519283918783016001600160c01b031960089260c01b1681520190565b60405191828092613e06565b039060025afa156114595760205f61549c81519360405191828092613e06565b039060025afa15611459576154df5f91615470602094613a2485516040519485938985019160619391600160f91b84526001840152602183015260418201520190565b039060025afa15611459575f5190565b615522615513612d766155056040850185611d0c565b91908535926104b784611dc4565b60208151910120913690613470565b60405161553781613a24602082019485613331565b51902090805f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f205482811461506d57806155a55750615039905f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f2090565b90507f657b94fe000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b90604082019161564e61560a6156056156456155f18786611d0c565b61560587959295359586926104b784611dc4565b615734565b6020815191012096875f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f20541595611d0c565b6104b784611dc4565b901561565e57506152db906157d7565b611dc0906040519182917f40470d740000000000000000000000000000000000000000000000000000000083526020600484018181520190610430565b60206001600160a01b037f2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad9216806001600160a01b03195f516020615f455f395f51905f525416175f516020615f455f395f51905f5255604051908152a1565b5f8061046593602081519101845af4615712613880565b91615d37565b610bf390613e3a614f2694936040519586936020850190613e06565b600961046591615757936001600160c01b03196040519586936020850190613e06565b917f0300000000000000000000000000000000000000000000000000000000000000835260c01b1660018201520301601f198101835282610bc3565b1561579a57565b7f760d6a9b000000000000000000000000000000000000000000000000000000005f5260045ffd5b90600161046592600160f91b81520190613e06565b906157e482511515615793565b6157ec611632565b915f925b81518410156158445760205f615816615809878661331d565b5160405191828092613e06565b039060025afa156114595760019061583c5f5191613a24604051938492602084016153a5565b9301926157f0565b60209293505f9150613a246154706154df926040519283918783016157c2565b5f9060409295939582966001600160e01b031984976001600160a01b038751938160208601967fb7009613000000000000000000000000000000000000000000000000000000008852166024860152166044840152166064820152606481526158ce608482610bc3565b8380528360205251915afa6158df57565b9150505f51906020518060201c150290565b8051600110156108ff5760210190565b9081518110156108ff570160200190565b8051600481109081156159f4575b50614819575f5b81518110156159ed5761597561596f6159696159438486615901565b517fff000000000000000000000000000000000000000000000000000000000000001690565b60f81c90565b60ff1690565b6020811790606182101591826159e1575b5081156159c3575b81156159a8575b506159a05750505f90565b600101615927565b600191501b6ba800000050006808000000001615155f615995565b90506030811015806159d6575b9061598e565b5060398111156159d0565b607a101591505f615986565b5050600190565b60809150115f615920565b80518251808210615a4957818082109118028082189114158102811890818103908111611c42576020615a3182613dd4565b92818401940101835e51902090602081519101201490565b505050505f90565b60ff5f516020615f655f395f51905f525460401c1615615a6d57565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f5260045ffd5b9060028201809211611c4257565b9060018201809211611c4257565b91908201809211611c4257565b8015611c42575f190190565b90813b15615b94576001600160a01b038216806001600160a01b03197f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5416177f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a2805115615b6357615b60916156fb565b50565b505034615b6c57565b7fb398979f000000000000000000000000000000000000000000000000000000005f5260045ffd5b6001600160a01b0382634c9c8ce360e01b5f521660045260245ffd5b60205f615bc4835160405191828092613e06565b039060025afa15611459575f5160205f615be78285015160405191828092613e06565b039060025afa15611459575f519160205f615c0c604084015160405191828092613e06565b039060025afa15611459575f519260205f615c31606085015160405191828092613e06565b039060025afa156114595760205f615c556080825195015160405191828092613e06565b039060025afa15611459576154705f93613a246020966154df958751916040519687958b8701939160a0959391855260208501526040840152606083015260808201520190565b805182118015615d30575b615cf4576001821180615cfc575b615cc090151561491e565b60280180602811611c4257615cd55f84613256565b03615cf4576001600160a01b0392915f615cee92615dc3565b90921690565b50505f905f90565b5060208101517fffff0000000000000000000000000000000000000000000000000000000000001661060f60f31b14615cb5565b505f615ca7565b90615d745750805115615d4c57602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b81511580615dba575b615d85575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15615d7d565b929092615dcf84615aa3565b831180615e6b575b615de9615df09194929394151561491e565b5f95615ab1565b915b818310615e025750505060019190565b9092919360ff615e3f615e3a602088860101517fff000000000000000000000000000000000000000000000000000000000000001690565b615ed2565b1690600f8211615e605790615e55600192614934565b019401919290615df2565b505f94508493505050565b50615df0615de961060f60f31b7fffff000000000000000000000000000000000000000000000000000000000000615ec8602089870101517fffff0000000000000000000000000000000000000000000000000000000000001690565b1614915050615dd7565b60f81c602f811180615f3a575b15615eee57602f190160ff1690565b6060811180615f30575b15615f07576056190160ff1690565b6040811180615f26575b15615f20576036190160ff1690565b5060ff90565b5060478110615f11565b5060678110615ef8565b50603a8110615edf56fef3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a00f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00a164736f6c634300081c000a1b9a2a4db4fdabb5b32192c8c2b90c1835653cf1aedfe8794a443b5c718f4c66d68750c3950334e058f8ccdd1a08540c778394eba7c6aa6b16155c1c6d05a9cef0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
}

// ContractICS26RouterABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractICS26RouterMetaData.ABI instead.
var ContractICS26RouterABI = ContractICS26RouterMetaData.ABI

// ContractICS26RouterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractICS26RouterMetaData.Bin instead.
var ContractICS26RouterBin = ContractICS26RouterMetaData.Bin

// DeployContractICS26Router deploys a new Ethereum contract, binding an instance of ContractICS26Router to it.
func DeployContractICS26Router(auth *bind.TransactOpts, backend bind.ContractBackend, clientMigrationProposer common.Address, clientMigrationExecutor common.Address) (common.Address, *types.Transaction, *ContractICS26Router, error) {
	parsed, err := ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractICS26RouterBin), backend, clientMigrationProposer, clientMigrationExecutor)
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

// GetClientMigration is a free data retrieval call binding the contract method 0x62d42988.
//
// Solidity: function getClientMigration(string clientId) view returns(bytes32 digest, uint48 executeAfter, uint48 expireAfter, address proposer)
func (_ContractICS26Router *ContractICS26RouterCaller) GetClientMigration(opts *bind.CallOpts, clientId string) (struct {
	Digest       [32]byte
	ExecuteAfter *big.Int
	ExpireAfter  *big.Int
	Proposer     common.Address
}, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "getClientMigration", clientId)

	outstruct := new(struct {
		Digest       [32]byte
		ExecuteAfter *big.Int
		ExpireAfter  *big.Int
		Proposer     common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Digest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.ExecuteAfter = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ExpireAfter = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Proposer = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// GetClientMigration is a free data retrieval call binding the contract method 0x62d42988.
//
// Solidity: function getClientMigration(string clientId) view returns(bytes32 digest, uint48 executeAfter, uint48 expireAfter, address proposer)
func (_ContractICS26Router *ContractICS26RouterSession) GetClientMigration(clientId string) (struct {
	Digest       [32]byte
	ExecuteAfter *big.Int
	ExpireAfter  *big.Int
	Proposer     common.Address
}, error) {
	return _ContractICS26Router.Contract.GetClientMigration(&_ContractICS26Router.CallOpts, clientId)
}

// GetClientMigration is a free data retrieval call binding the contract method 0x62d42988.
//
// Solidity: function getClientMigration(string clientId) view returns(bytes32 digest, uint48 executeAfter, uint48 expireAfter, address proposer)
func (_ContractICS26Router *ContractICS26RouterCallerSession) GetClientMigration(clientId string) (struct {
	Digest       [32]byte
	ExecuteAfter *big.Int
	ExpireAfter  *big.Int
	Proposer     common.Address
}, error) {
	return _ContractICS26Router.Contract.GetClientMigration(&_ContractICS26Router.CallOpts, clientId)
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
func (_ContractICS26Router *ContractICS26RouterCaller) GetCounterparty(opts *bind.CallOpts, clientId string) (ICS02ClientMsgsCounterpartyInfo, error) {
	var out []interface{}
	err := _ContractICS26Router.contract.Call(opts, &out, "getCounterparty", clientId)

	if err != nil {
		return *new(ICS02ClientMsgsCounterpartyInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(ICS02ClientMsgsCounterpartyInfo)).(*ICS02ClientMsgsCounterpartyInfo)

	return out0, err

}

// GetCounterparty is a free data retrieval call binding the contract method 0xb0777bfa.
//
// Solidity: function getCounterparty(string clientId) view returns((string,bytes[]))
func (_ContractICS26Router *ContractICS26RouterSession) GetCounterparty(clientId string) (ICS02ClientMsgsCounterpartyInfo, error) {
	return _ContractICS26Router.Contract.GetCounterparty(&_ContractICS26Router.CallOpts, clientId)
}

// GetCounterparty is a free data retrieval call binding the contract method 0xb0777bfa.
//
// Solidity: function getCounterparty(string clientId) view returns((string,bytes[]))
func (_ContractICS26Router *ContractICS26RouterCallerSession) GetCounterparty(clientId string) (ICS02ClientMsgsCounterpartyInfo, error) {
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
func (_ContractICS26Router *ContractICS26RouterTransactor) AckPacket(opts *bind.TransactOpts, msg_ ICS26RouterMsgsMsgAckPacket) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "ackPacket", msg_)
}

// AckPacket is a paid mutator transaction binding the contract method 0xfdbd955d.
//
// Solidity: function ackPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes,bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterSession) AckPacket(msg_ ICS26RouterMsgsMsgAckPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AckPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// AckPacket is a paid mutator transaction binding the contract method 0xfdbd955d.
//
// Solidity: function ackPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes,bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) AckPacket(msg_ ICS26RouterMsgsMsgAckPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AckPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// AddClient is a paid mutator transaction binding the contract method 0x1ec43e23.
//
// Solidity: function addClient(string clientId, (string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterTransactor) AddClient(opts *bind.TransactOpts, clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "addClient", clientId, counterpartyInfo, client)
}

// AddClient is a paid mutator transaction binding the contract method 0x1ec43e23.
//
// Solidity: function addClient(string clientId, (string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterSession) AddClient(clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AddClient(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// AddClient is a paid mutator transaction binding the contract method 0x1ec43e23.
//
// Solidity: function addClient(string clientId, (string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterTransactorSession) AddClient(clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AddClient(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// AddClient0 is a paid mutator transaction binding the contract method 0xe3cb36a0.
//
// Solidity: function addClient((string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterTransactor) AddClient0(opts *bind.TransactOpts, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "addClient0", counterpartyInfo, client)
}

// AddClient0 is a paid mutator transaction binding the contract method 0xe3cb36a0.
//
// Solidity: function addClient((string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterSession) AddClient0(counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.AddClient0(&_ContractICS26Router.TransactOpts, counterpartyInfo, client)
}

// AddClient0 is a paid mutator transaction binding the contract method 0xe3cb36a0.
//
// Solidity: function addClient((string,bytes[]) counterpartyInfo, address client) returns(string)
func (_ContractICS26Router *ContractICS26RouterTransactorSession) AddClient0(counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
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

// CancelClientMigration is a paid mutator transaction binding the contract method 0x19e84399.
//
// Solidity: function cancelClientMigration(string clientId) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) CancelClientMigration(opts *bind.TransactOpts, clientId string) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "cancelClientMigration", clientId)
}

// CancelClientMigration is a paid mutator transaction binding the contract method 0x19e84399.
//
// Solidity: function cancelClientMigration(string clientId) returns()
func (_ContractICS26Router *ContractICS26RouterSession) CancelClientMigration(clientId string) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.CancelClientMigration(&_ContractICS26Router.TransactOpts, clientId)
}

// CancelClientMigration is a paid mutator transaction binding the contract method 0x19e84399.
//
// Solidity: function cancelClientMigration(string clientId) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) CancelClientMigration(clientId string) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.CancelClientMigration(&_ContractICS26Router.TransactOpts, clientId)
}

// ExecuteClientMigration is a paid mutator transaction binding the contract method 0xee5c877c.
//
// Solidity: function executeClientMigration(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) ExecuteClientMigration(opts *bind.TransactOpts, clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "executeClientMigration", clientId, counterpartyInfo, client)
}

// ExecuteClientMigration is a paid mutator transaction binding the contract method 0xee5c877c.
//
// Solidity: function executeClientMigration(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterSession) ExecuteClientMigration(clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.ExecuteClientMigration(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// ExecuteClientMigration is a paid mutator transaction binding the contract method 0xee5c877c.
//
// Solidity: function executeClientMigration(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) ExecuteClientMigration(clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.ExecuteClientMigration(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
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
func (_ContractICS26Router *ContractICS26RouterTransactor) MigrateClient(opts *bind.TransactOpts, clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "migrateClient", clientId, counterpartyInfo, client)
}

// MigrateClient is a paid mutator transaction binding the contract method 0xcce0b265.
//
// Solidity: function migrateClient(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterSession) MigrateClient(clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.MigrateClient(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// MigrateClient is a paid mutator transaction binding the contract method 0xcce0b265.
//
// Solidity: function migrateClient(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) MigrateClient(clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
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

// ProposeClientMigration is a paid mutator transaction binding the contract method 0x0cbdf9f6.
//
// Solidity: function proposeClientMigration(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) ProposeClientMigration(opts *bind.TransactOpts, clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "proposeClientMigration", clientId, counterpartyInfo, client)
}

// ProposeClientMigration is a paid mutator transaction binding the contract method 0x0cbdf9f6.
//
// Solidity: function proposeClientMigration(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterSession) ProposeClientMigration(clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.ProposeClientMigration(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// ProposeClientMigration is a paid mutator transaction binding the contract method 0x0cbdf9f6.
//
// Solidity: function proposeClientMigration(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) ProposeClientMigration(clientId string, counterpartyInfo ICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.ProposeClientMigration(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// RecvPacket is a paid mutator transaction binding the contract method 0x596e00b9.
//
// Solidity: function recvPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) RecvPacket(opts *bind.TransactOpts, msg_ ICS26RouterMsgsMsgRecvPacket) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "recvPacket", msg_)
}

// RecvPacket is a paid mutator transaction binding the contract method 0x596e00b9.
//
// Solidity: function recvPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterSession) RecvPacket(msg_ ICS26RouterMsgsMsgRecvPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.RecvPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// RecvPacket is a paid mutator transaction binding the contract method 0x596e00b9.
//
// Solidity: function recvPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) RecvPacket(msg_ ICS26RouterMsgsMsgRecvPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.RecvPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// SendPacket is a paid mutator transaction binding the contract method 0x4d6e7ce3.
//
// Solidity: function sendPacket((string,uint64,(string,string,string,string,bytes)) msg_) returns(uint64)
func (_ContractICS26Router *ContractICS26RouterTransactor) SendPacket(opts *bind.TransactOpts, msg_ ICS26RouterMsgsMsgSendPacket) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "sendPacket", msg_)
}

// SendPacket is a paid mutator transaction binding the contract method 0x4d6e7ce3.
//
// Solidity: function sendPacket((string,uint64,(string,string,string,string,bytes)) msg_) returns(uint64)
func (_ContractICS26Router *ContractICS26RouterSession) SendPacket(msg_ ICS26RouterMsgsMsgSendPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.SendPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// SendPacket is a paid mutator transaction binding the contract method 0x4d6e7ce3.
//
// Solidity: function sendPacket((string,uint64,(string,string,string,string,bytes)) msg_) returns(uint64)
func (_ContractICS26Router *ContractICS26RouterTransactorSession) SendPacket(msg_ ICS26RouterMsgsMsgSendPacket) (*types.Transaction, error) {
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
func (_ContractICS26Router *ContractICS26RouterTransactor) TimeoutPacket(opts *bind.TransactOpts, msg_ ICS26RouterMsgsMsgTimeoutPacket) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "timeoutPacket", msg_)
}

// TimeoutPacket is a paid mutator transaction binding the contract method 0x223e357a.
//
// Solidity: function timeoutPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterSession) TimeoutPacket(msg_ ICS26RouterMsgsMsgTimeoutPacket) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.TimeoutPacket(&_ContractICS26Router.TransactOpts, msg_)
}

// TimeoutPacket is a paid mutator transaction binding the contract method 0x223e357a.
//
// Solidity: function timeoutPacket(((uint64,string,string,uint64,(string,string,string,string,bytes)[]),bytes) msg_) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) TimeoutPacket(msg_ ICS26RouterMsgsMsgTimeoutPacket) (*types.Transaction, error) {
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
	Packet          ICS26RouterMsgsPacket
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
	CounterpartyInfo ICS02ClientMsgsCounterpartyInfo
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
	CounterpartyInfo ICS02ClientMsgsCounterpartyInfo
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

// ContractICS26RouterICS02ClientMigrationCancelledIterator is returned from FilterICS02ClientMigrationCancelled and is used to iterate over the raw logs and unpacked data for ICS02ClientMigrationCancelled events raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientMigrationCancelledIterator struct {
	Event *ContractICS26RouterICS02ClientMigrationCancelled // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterICS02ClientMigrationCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterICS02ClientMigrationCancelled)
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
		it.Event = new(ContractICS26RouterICS02ClientMigrationCancelled)
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
func (it *ContractICS26RouterICS02ClientMigrationCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterICS02ClientMigrationCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterICS02ClientMigrationCancelled represents a ICS02ClientMigrationCancelled event raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientMigrationCancelled struct {
	ClientId common.Hash
	Digest   [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterICS02ClientMigrationCancelled is a free log retrieval operation binding the contract event 0x676356e432cc05e98d6a429f0c0f7518220b62035e48089c8212d703a6a633b9.
//
// Solidity: event ICS02ClientMigrationCancelled(string indexed clientId, bytes32 indexed digest)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterICS02ClientMigrationCancelled(opts *bind.FilterOpts, clientId []string, digest [][32]byte) (*ContractICS26RouterICS02ClientMigrationCancelledIterator, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var digestRule []interface{}
	for _, digestItem := range digest {
		digestRule = append(digestRule, digestItem)
	}

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "ICS02ClientMigrationCancelled", clientIdRule, digestRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterICS02ClientMigrationCancelledIterator{contract: _ContractICS26Router.contract, event: "ICS02ClientMigrationCancelled", logs: logs, sub: sub}, nil
}

// WatchICS02ClientMigrationCancelled is a free log subscription operation binding the contract event 0x676356e432cc05e98d6a429f0c0f7518220b62035e48089c8212d703a6a633b9.
//
// Solidity: event ICS02ClientMigrationCancelled(string indexed clientId, bytes32 indexed digest)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchICS02ClientMigrationCancelled(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterICS02ClientMigrationCancelled, clientId []string, digest [][32]byte) (event.Subscription, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var digestRule []interface{}
	for _, digestItem := range digest {
		digestRule = append(digestRule, digestItem)
	}

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "ICS02ClientMigrationCancelled", clientIdRule, digestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterICS02ClientMigrationCancelled)
				if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientMigrationCancelled", log); err != nil {
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

// ParseICS02ClientMigrationCancelled is a log parse operation binding the contract event 0x676356e432cc05e98d6a429f0c0f7518220b62035e48089c8212d703a6a633b9.
//
// Solidity: event ICS02ClientMigrationCancelled(string indexed clientId, bytes32 indexed digest)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseICS02ClientMigrationCancelled(log types.Log) (*ContractICS26RouterICS02ClientMigrationCancelled, error) {
	event := new(ContractICS26RouterICS02ClientMigrationCancelled)
	if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientMigrationCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractICS26RouterICS02ClientMigrationProposedIterator is returned from FilterICS02ClientMigrationProposed and is used to iterate over the raw logs and unpacked data for ICS02ClientMigrationProposed events raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientMigrationProposedIterator struct {
	Event *ContractICS26RouterICS02ClientMigrationProposed // Event containing the contract specifics and raw log

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
func (it *ContractICS26RouterICS02ClientMigrationProposedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractICS26RouterICS02ClientMigrationProposed)
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
		it.Event = new(ContractICS26RouterICS02ClientMigrationProposed)
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
func (it *ContractICS26RouterICS02ClientMigrationProposedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractICS26RouterICS02ClientMigrationProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractICS26RouterICS02ClientMigrationProposed represents a ICS02ClientMigrationProposed event raised by the ContractICS26Router contract.
type ContractICS26RouterICS02ClientMigrationProposed struct {
	ClientId     common.Hash
	Digest       [32]byte
	ExecuteAfter *big.Int
	ExpireAfter  *big.Int
	Proposer     common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterICS02ClientMigrationProposed is a free log retrieval operation binding the contract event 0x9bb264c2dd85dcb1e978b2be15709752a375d155b99d385f0b7f81020433e2d2.
//
// Solidity: event ICS02ClientMigrationProposed(string indexed clientId, bytes32 indexed digest, uint48 executeAfter, uint48 expireAfter, address proposer)
func (_ContractICS26Router *ContractICS26RouterFilterer) FilterICS02ClientMigrationProposed(opts *bind.FilterOpts, clientId []string, digest [][32]byte) (*ContractICS26RouterICS02ClientMigrationProposedIterator, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var digestRule []interface{}
	for _, digestItem := range digest {
		digestRule = append(digestRule, digestItem)
	}

	logs, sub, err := _ContractICS26Router.contract.FilterLogs(opts, "ICS02ClientMigrationProposed", clientIdRule, digestRule)
	if err != nil {
		return nil, err
	}
	return &ContractICS26RouterICS02ClientMigrationProposedIterator{contract: _ContractICS26Router.contract, event: "ICS02ClientMigrationProposed", logs: logs, sub: sub}, nil
}

// WatchICS02ClientMigrationProposed is a free log subscription operation binding the contract event 0x9bb264c2dd85dcb1e978b2be15709752a375d155b99d385f0b7f81020433e2d2.
//
// Solidity: event ICS02ClientMigrationProposed(string indexed clientId, bytes32 indexed digest, uint48 executeAfter, uint48 expireAfter, address proposer)
func (_ContractICS26Router *ContractICS26RouterFilterer) WatchICS02ClientMigrationProposed(opts *bind.WatchOpts, sink chan<- *ContractICS26RouterICS02ClientMigrationProposed, clientId []string, digest [][32]byte) (event.Subscription, error) {

	var clientIdRule []interface{}
	for _, clientIdItem := range clientId {
		clientIdRule = append(clientIdRule, clientIdItem)
	}
	var digestRule []interface{}
	for _, digestItem := range digest {
		digestRule = append(digestRule, digestItem)
	}

	logs, sub, err := _ContractICS26Router.contract.WatchLogs(opts, "ICS02ClientMigrationProposed", clientIdRule, digestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractICS26RouterICS02ClientMigrationProposed)
				if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientMigrationProposed", log); err != nil {
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

// ParseICS02ClientMigrationProposed is a log parse operation binding the contract event 0x9bb264c2dd85dcb1e978b2be15709752a375d155b99d385f0b7f81020433e2d2.
//
// Solidity: event ICS02ClientMigrationProposed(string indexed clientId, bytes32 indexed digest, uint48 executeAfter, uint48 expireAfter, address proposer)
func (_ContractICS26Router *ContractICS26RouterFilterer) ParseICS02ClientMigrationProposed(log types.Log) (*ContractICS26RouterICS02ClientMigrationProposed, error) {
	event := new(ContractICS26RouterICS02ClientMigrationProposed)
	if err := _ContractICS26Router.contract.UnpackLog(event, "ICS02ClientMigrationProposed", log); err != nil {
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
	Packet   ICS26RouterMsgsPacket
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
	Packet   ICS26RouterMsgsPacket
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
	Packet           ICS26RouterMsgsPacket
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
