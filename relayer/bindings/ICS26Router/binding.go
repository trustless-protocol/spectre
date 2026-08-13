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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"clientMigrationProposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"clientMigrationExecutor\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ackPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgAckPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgement\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"membershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addClient\",\"inputs\":[{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addIBCApp\",\"inputs\":[{\"name\":\"app\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addIBCApp\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"app\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelClientMigration\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeClientMigration\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractILightClient\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getClientMigration\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"executeAfter\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"expireAfter\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitment\",\"inputs\":[{\"name\":\"hashedPath\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCounterparty\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIBCApp\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIIBCApp\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLightClientMigratorRole\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getNextClientSeq\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializeV2\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isConsumingScheduledOp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposeClientMigration\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recvPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgRecvPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"membershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgSendPacket\",\"components\":[{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payload\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Payload\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAuthority\",\"inputs\":[{\"name\":\"newAuthority\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitMisbehaviour\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"misbehaviourMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"timeoutPacket\",\"inputs\":[{\"name\":\"msg_\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.MsgTimeoutPacket\",\"components\":[{\"name\":\"packet\",\"type\":\"tuple\",\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"nonMembershipMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unfreezeClient\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateApplicationState\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateConsensusState\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"updateMsg\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AckPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgement\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCAppAdded\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"app\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IBCAppRecvPacketCallbackError\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientAdded\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientMigrated\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"counterpartyInfo\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS02ClientMsgs.CounterpartyInfo\",\"components\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"merklePrefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"name\":\"client\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientMigrationCancelled\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientMigrationProposed\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"digest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"executeAfter\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"expireAfter\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"proposer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientUnfrozen\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02ClientUpdated\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"result\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumILightClientMsgs.UpdateResult\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ICS02MisbehaviourSubmitted\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Noop\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SendPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TimeoutPacket\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WriteAcknowledgement\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"sequence\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packet\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIICS26RouterMsgs.Packet\",\"components\":[{\"name\":\"sequence\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destClient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timeoutTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"payloads\",\"type\":\"tuple[]\",\"internalType\":\"structIICS26RouterMsgs.Payload[]\",\"components\":[{\"name\":\"sourcePort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"destPort\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"encoding\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]},{\"name\":\"acknowledgements\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessManagedInvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessManagedRequiredDelay\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"AccessManagedUnauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"DefaultAdminRoleCannotBeGranted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCAppNotFound\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCAsyncAcknowledgementNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCClientAlreadyExists\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationAlreadyProposed\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationDelayRequired\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationDirectCallNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCClientMigrationExpired\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"expireAfter\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationMismatch\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationModuleMismatch\",\"inputs\":[{\"name\":\"module\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expectedModuleId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationModuleMissing\",\"inputs\":[{\"name\":\"module\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationNotProposed\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCClientMigrationNotReady\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"executeAfter\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCClientNotFound\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCCounterpartyClientNotFound\",\"inputs\":[{\"name\":\"counterpartyClientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCErrorUniversalAcknowledgement\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCFailedCallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCInvalidClientId\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidCounterparty\",\"inputs\":[{\"name\":\"expected\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"actual\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidPortIdentifier\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidTimeoutDuration\",\"inputs\":[{\"name\":\"maxTimeoutDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualTimeoutDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCInvalidTimeoutTimestamp\",\"inputs\":[{\"name\":\"timeoutTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"comparedTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IBCMultiPayloadPacketNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IBCPacketAcknowledgementAlreadyExists\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"IBCPacketCommitmentAlreadyExists\",\"inputs\":[{\"name\":\"path\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"IBCPacketCommitmentMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IBCPacketReceiptMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IBCPortAlreadyExists\",\"inputs\":[{\"name\":\"portId\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"IBCUnauthorizedMigrator\",\"inputs\":[{\"name\":\"clientId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"IBCUnauthorizedSender\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMerklePrefix\",\"inputs\":[{\"name\":\"prefix\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]},{\"type\":\"error\",\"name\":\"NoAcknowledgements\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"Unreachable\",\"inputs\":[]}]",
	Bin: "0x60e08060405234610279576040816163d5803803809161001f8285610356565b8339810103126102795761003e602061003783610379565b9201610379565b90803b15610336575f80604051602081019063a1308f2760e01b825260048152610069602482610356565b5190845afa3d1561032e573d906001600160401b03821161028a576040519161009c601f8201601f191660200184610356565b82523d5f602084013e5b15908115610321575b81156102f8575b506102c757813b156102a6575f80604051602081019063a1308f2760e01b8252600481526100e5602482610356565b5190855afa3d1561029e573d906001600160401b03821161028a5760405191610118601f8201601f191660200184610356565b82523d5f602084013e5b1590811561027d575b8115610250575b5061021e576001600160a01b039081166080521660a0523060c0525f5160206163b55f395f51905f5254604081901c60ff1661020f576002600160401b03196001600160401b038216016101b9575b604051615fe7908161038e82396080518161035a015260a051818181610401015261196b015260c051818181610cc50152610e680152f35b6001600160401b0319166001600160401b039081175f5160206163b55f395f51905f52556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f610181565b63f92ee8a960e01b5f5260045ffd5b506382f0c90760e01b5f9081526001600160a01b03919091166004525f5160206163955f395f51905f52602452604490fd5b905060208180518101031261027957602001515f5160206163955f395f51905f5214155f610132565b5f80fd5b805160201415915061012b565b634e487b7160e01b5f52604160045260245ffd5b606090610122565b50632b833ccb60e11b5f9081526001600160a01b0391909116600452602490fd5b6382f0c90760e01b5f9081526001600160a01b03919091166004525f5160206163755f395f51905f52602452604490fd5b905060208180518101031261027957602001515f5160206163755f395f51905f5214155f6100b6565b80516020141591506100af565b6060906100a6565b632b833ccb60e11b5f9081526001600160a01b0391909116600452602490fd5b601f909101601f19168101906001600160401b0382119082101761028a57604052565b51906001600160a01b03821682036102795756fe60806040526004361015610011575f80fd5b5f3560e01c80630cbdf9f61461023f57806319e843991461023a5780631ec43e2314610235578063223e357a146102305780632447af291461022b57806327f146f31461022657806329b6eca9146102215780633f4ba83a1461021c5780634b720d5b146102175780634d6e7ce3146102125780634f1ef2861461020d57806352d1902d14610208578063596e00b9146102035780635c975abb146101fe5780635f516889146101f957806362d42988146101f45780637795820c146101ef5780637a9e5e4b146101ea5780637eb78932146101e55780638456cb59146101e05780638fb36037146101db5780639c11bece146101d65780639e2e5c83146101d1578063ac9650d8146101cc578063ad3cb1cc146101c7578063b0777bfa146101c2578063b0830ab9146101bd578063bf7e214f146101b8578063c45f24ca146101b3578063c4d66de8146101ae578063cce0b2651461019f578063d395af7f146101a9578063e3cb36a0146101a4578063ee5c877c1461019f5763fdbd955d1461019a575f80fd5b611c30565b611923565b611a06565b61199a565b611859565b611779565b611747565b6116d7565b61168e565b61162f565b6115af565b611447565b61132c565b611242565b6111a9565b611190565b611102565b6110b9565b610fc8565b610f29565b610ee8565b610eb8565b610e4e565b610c87565b61091b565b6107df565b610720565b6105ce565b610592565b61055e565b610511565b610468565b6103bb565b610312565b9181601f84011215610271578235916001600160401b038311610271576020838186019501011161027157565b5f80fd5b908160409103126102715790565b600435906001600160a01b038216820361027157565b602435906001600160a01b038216820361027157565b6060600319820112610271576004356001600160401b03811161027157816102d991600401610244565b92909291602435906001600160401b038211610271576102fb91600401610275565b906044356001600160a01b03811681036102715790565b3461027157610320366102af565b505050507f0cbdf9f6000000000000000000000000000000000000000000000000000000005f5236600480375f8036816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af43d5f803e15610389573d5ff35b3d5ffd5b602060031982011261027157600435906001600160401b038211610271576103b791600401610244565b9091565b34610271576103c93661038d565b50507f19e84399000000000000000000000000000000000000000000000000000000005f5236600480375f8036816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af43d5f803e15610389573d5ff35b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b906020610465928181520190610430565b90565b34610271576104e36104cf6104c861047f366102af565b9061048e94929394363361462b565b61049b8486811515611c89565b6104b884866104b36104ae368484610c36565b614828565b611c89565b6104c3368587610c36565b614e03565b3691610c36565b604051918291602083526020830190610430565b0390f35b602060031982011261027157600435906001600160401b0382116102715761046591600401610275565b3461027157610539610522366104e7565b61052a614f24565b610534363361462b565b612ce8565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d005b346102715760206105776105713661038d565b90613088565b6001600160a01b0360405191168152f35b5f91031261027157565b34610271575f3660031901126102715760207f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960254604051908152f35b34610271576020366003190112610271576105e7610283565b5f516020615fbb5f395f51905f525460ff6001600160401b0382169161060f6001841461310e565b60401c16908115610714575b506106ec5761068e9061065460026001600160401b03195f516020615fbb5f395f51905f525416175f516020615fbb5f395f51905f5255565b6801000000000000000068ff0000000000000000195f516020615fbb5f395f51905f525416175f516020615fbb5f395f51905f5255613152565b6106bd68ff0000000000000000195f516020615fbb5f395f51905f5254165f516020615fbb5f395f51905f5255565b604051600281527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a1005b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f5260045ffd5b6002915010155f61061b565b34610271575f3660031901126102715761073a363361462b565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff8116156107b75760ff19167fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a1005b7f8dfc202b000000000000000000000000000000000000000000000000000000005f5260045ffd5b34610271576020366003190112610271576107f8610283565b610800614f24565b6001600160a01b038116908161081e6108196028615aeb565b613e05565b90603061082a83613341565b53607861083683614807565b536108416028615af9565b6001811161089457506108625761085892506151ff565b610860614f98565b005b7fe22e27eb000000000000000000000000000000000000000000000000000000005f526004839052601460245260445ffd5b90600f811660108110156108e8576108e3917f30313233343536373839616263646566000000000000000000000000000000006108dd921a6108d68587614817565b5360041c90565b91615b14565b610841565b611d4e565b6020600319820112610271576004356001600160401b03811161027157606090600401809203126102715790565b34610271576104e361092c366108ed565b610934615308565b61093c614f24565b604081019061097b61096961095d6105716109578686611cd5565b80611d72565b6001600160a01b031690565b33906001600160a01b03168114613115565b610aae610aa861099461098e8480611d72565b9061405b565b5192610a9e610a3b610a8f60208401976109cd6109b08a611e17565b42906109bb8c611e17565b906001600160401b03429116116128fc565b610a0d620151806109f5426109f06109e48e611e17565b6001600160401b031690565b613287565b1115610a07426109f06109e48e611e17565b90613294565b610a79610a23610a1d8780611d72565b9061535b565b998a99610a5f610a338980611d72565b979093611e17565b92610a446132ce565b976104c8610a50610bcd565b6001600160401b03909f168f52565b60208c015260408b01526001600160401b031660608a0152565b826080890152610a8a369186611cd5565b61293d565b610a9882613341565b52613341565b50610957846153c9565b90612add565b7fab3a4458a269be61dfa43faa33aa7b1f5d570716f83ad078bc2ba5dab039abae60405180610ae76001600160401b0387169582613362565b0390a3610af2614f98565b6040516001600160401b0390911681529081906020820190565b634e487b7160e01b5f52604160045260245ffd5b608081019081106001600160401b03821117610b3b57604052565b610b0c565b604081019081106001600160401b03821117610b3b57604052565b606081019081106001600160401b03821117610b3b57604052565b602081019081106001600160401b03821117610b3b57604052565b60a081019081106001600160401b03821117610b3b57604052565b90601f801991011681019081106001600160401b03821117610b3b57604052565b60405190610bdc60a083610bac565b565b60405190610bdc608083610bac565b60405190610bdc60e083610bac565b60405190610bdc61010083610bac565b60405190610bdc60c083610bac565b6001600160401b038111610b3b57601f01601f191660200190565b929192610c4282610c1b565b91610c506040519384610bac565b829481845281830111610271578281602093845f960137010152565b9080601f830112156102715781602061046593359101610c36565b604036600319011261027157610c9b610283565b6024356001600160401b03811161027157610cba903690600401610c6c565b906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016803014908115610e19575b50610df157610cff363361462b565b604051917f52d1902d0000000000000000000000000000000000000000000000000000000083526020836004816001600160a01b0386165afa5f9381610dc0575b50610d6157634c9c8ce360e01b5f526001600160a01b03821660045260245ffd5b907f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8303610d93576108609250615b20565b7faa1d49a4000000000000000000000000000000000000000000000000000000005f52600483905260245ffd5b610de391945060203d602011610dea575b610ddb8183610bac565b810190612542565b925f610d40565b503d610dd1565b7fe07c8dba000000000000000000000000000000000000000000000000000000005f5260045ffd5b90506001600160a01b037f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc541614155f610cf0565b34610271575f366003190112610271576001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000163003610df15760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b3461027157610539610ec9366104e7565b610ed1615308565b610ed9614f24565b610ee3363361462b565b613934565b34610271575f36600319011261027157602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b34610271576040366003190112610271576004356001600160401b03811161027157610fc3610f5f610539923690600401610244565b9190610f69610299565b92610f72614f24565b610f7c363361462b565b610f898183811515613cf2565b610fab8183610fa4610f9c368484610c36565b805190615cf2565b5015613cf2565b6104c88183610fbe6104ae368484610c36565b613cf2565b6151ff565b34610271576020610fd83661038d565b919082604051938492833781017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496038152030190206040519061101982610b20565b6104e36060611080611073611067600186549687895201549665ffffffffffff88169081602082015265ffffffffffff8960301c1698896040830152861c958691015265ffffffffffff1690565b9565ffffffffffff1690565b916001600160a01b031690565b9060405194859485929365ffffffffffff6001600160a01b03929695816060956080880199885216602087015216604085015216910152565b34610271576020366003190112610271576004355f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600602052602060405f2054604051908152f35b346102715760203660031901126102715761111b610283565b6001600160a01b035f516020615f9b5f395f51905f525416330361117e57803b1561114957610860906157ff565b6001600160a01b03907fc2f31e5e000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b62d1953b60e31b5f523360045260245ffd5b346102715760206105776111a33661038d565b90613d3a565b34610271575f366003190112610271576111c3363361462b565b6111cb615308565b600160ff197fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005416177fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a1005b34610271575f366003190112610271575f516020615f9b5f395f51905f525460a01c60ff16156112a25760207f8fb36037000000000000000000000000000000000000000000000000000000005b6001600160e01b031960405191168152f35b60205f611290565b6040600319820112610271576004356001600160401b03811161027157816112d491600401610244565b92909291602435906001600160401b038211610271576103b791600401610244565b634e487b7160e01b5f52602160045260245ffd5b6003111561131457565b6112f6565b919060208301926113298261130a565b52565b3461027157611398602061133f366112aa565b9061134e94929394363361462b565b6001600160a01b036113608587613d3a565b16905f6040518098819582947f4b1872be00000000000000000000000000000000000000000000000000000000845260048401611c78565b03925af1908115611442576104e3935f926113f1575b5081926113e27f87bbef2779889a19f0435ddca81fda94132c06ffddb0ea73def256307a293aef9360405193849384613dd5565b0390a160405191829182611319565b7f87bbef2779889a19f0435ddca81fda94132c06ffddb0ea73def256307a293aef92506114359060203d60201161143b575b61142d8183610bac565b810190613dc0565b916113ae565b503d611423565b6128f1565b3461027157611455366112aa565b9192905f92611464363361462b565b6001600160a01b036114768685613d3a565b1691823b15610271576114c2925f92836040518096819582947fddba65370000000000000000000000000000000000000000000000000000000084526020600485018181520191611c58565b03925af1801561144257611512575b507fa263f0a976b2937a51fd2e416491cf0ca724d5499fa870715929dfde4ee4a430919261150c604051928392602084526020840191611c58565b0390a180f35b7fa263f0a976b2937a51fd2e416491cf0ca724d5499fa870715929dfde4ee4a43092505f61153f91610bac565b5f916114d1565b9080602083519182815201916020808360051b8301019401925f915b83831061157157505050505090565b909192939460208061158f600193601f198682030187528951610430565b97019301930191939290611562565b906020610465928181520190611546565b34610271576020366003190112610271576004356001600160401b0381116102715736602382011215610271578060040135906001600160401b038211610271573660248360051b83010111610271576104e391602461160f9201613e79565b6040519182918261159e565b6040519061162a602083610bac565b5f8252565b34610271575f366003190112610271576104e3604051611650604082610bac565b600581527f352e302e300000000000000000000000000000000000000000000000000000006020820152604051918291602083526020830190610430565b34610271576104e36116a261098e3661038d565b6040519182916020835260206116c382516040838701526060860190610430565b910151838203601f19016040850152611546565b346102715760206001600160401b036116ef3661038d565b61173b603b604051838194888301967f4c494748545f434c49454e545f4d49475241544f525f524f4c455f000000000088528484013781015f838201520301601f198101835282610bac565b51902016604051908152f35b34610271575f3660031901126102715760206001600160a01b035f516020615f9b5f395f51905f525416604051908152f35b34610271576117873661038d565b905f90611794363361462b565b6001600160a01b036117a68483613d3a565b16803b15610271575f80916004604051809481937f6a28f0000000000000000000000000000000000000000000000000000000000083525af1801561144257611825575b507f7aa8dd0b42a8f6969fbecb42ee2a8fdec2a1a3a576bfc4623b9e7f98ba76fa99919261150c604051928392602084526020840191611c58565b7f7aa8dd0b42a8f6969fbecb42ee2a8fdec2a1a3a576bfc4623b9e7f98ba76fa9992505f61185291610bac565b5f916117ea565b3461027157602036600319011261027157611872610283565b5f516020615fbb5f395f51905f525460ff6001600160401b03821691611898831561310e565b60401c16908115611917575b506106ec5761068e906118dd60026001600160401b03195f516020615fbb5f395f51905f525416175f516020615fbb5f395f51905f5255565b6801000000000000000068ff0000000000000000195f516020615fbb5f395f51905f525416175f516020615fbb5f395f51905f5255614143565b6002915010155f6118a4565b3461027157611931366102af565b505050507fee5c877c000000000000000000000000000000000000000000000000000000005f5236600480375f8036816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af43d5f803e15610389573d5ff35b346102715761139860206119ad366112aa565b906119bc94929394363361462b565b6001600160a01b036119ce8587613d3a565b16905f6040518098819582947f92c19bdc00000000000000000000000000000000000000000000000000000000845260048401611c78565b34610271576040366003190112610271576004356001600160401b03811161027157611a36903690600401610275565b611a3e610299565b90611a476147cc565b7f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496025490925f198214611c2b57600182017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496025581825f937a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000811015611c00575b50806d04ee2d6d415b85acef8100000000600a921015611be4575b662386f26fc10000811015611bcf575b6305f5e100811015611bbd575b612710811015611bad575b6064811015611b9e575b1015611b93575b611b5e6021611b2660018601613e05565b948501015b5f1901917f3031323334353637383961626364656600000000000000000000000000000000600a82061a8353600a900490565b908115611b6e57611b5e90611b2b565b5050611b80611b87926104e39561587c565b9283614e03565b60405191829182610454565b600190920191611b15565b60029060649004940193611b0e565b6004906127109004940193611b04565b6008906305f5e1009004940193611af9565b601090662386f26fc100009004940193611aec565b6020906d04ee2d6d415b85acef81000000009004940193611adc565b604094507a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000090049050600a611ac1565b613265565b3461027157610539611c41366108ed565b611c49614f24565b611c53363361462b565b614216565b908060209392818452848401375f828201840152601f01601f1916010190565b916020610465938181520191611c58565b91909115611c95575050565b611cd16040519283927f4870bd740000000000000000000000000000000000000000000000000000000084526020600485018181520191611c58565b0390fd5b903590609e1981360301821215610271570190565b903590601e198136030182121561027157018035906001600160401b03821161027157602001918160051b3603831361027157565b15611d2657565b7f356f4dbd000000000000000000000000000000000000000000000000000000005f5260045ffd5b634e487b7160e01b5f52603260045260245ffd5b90156108e8578061046591611cd5565b903590601e198136030182121561027157018035906001600160401b0382116102715760200191813603831361027157565b92909215611db157505050565b611df49291611cd1916040519485947f9fff831f000000000000000000000000000000000000000000000000000000008652604060048701526044860190610430565b84810360031901602486015291611c58565b6001600160401b0381160361027157565b3561046581611e06565b3590610bdc82611e06565b919082604091031261027157604051611e4481610b40565b60208082948035611e5481611e06565b8452013591611e6283611e06565b0152565b6001600160401b038111610b3b5760051b60200190565b9080601f83011215610271578135611e9481611e66565b92611ea26040519485610bac565b81845260208085019260051b820101918383116102715760208201905b838210611ece57505050505090565b81356001600160401b03811161027157602091611ef087848094880101610c6c565b815201910190611ebf565b9080601f8301121561027157813591611f1383611e66565b92611f216040519485610bac565b80845260208085019160051b830101918383116102715760208101915b838310611f4d57505050505090565b82356001600160401b038111610271578201906040828703601f1901126102715760405190611f7b82610b40565b60208301356001600160401b03811161027157876020611f9d92860101611e7d565b82526040830135916001600160401b03831161027157611fc588602080969581960101610c6c565b83820152815201920191611f3e565b6003111561027157565b91906080838203126102715760405190611ff782610b20565b8193803561200481611fd4565b8352602081013561201481611fd4565b6020840152604081013561202781611fd4565b60408401526060810135916001600160401b0383116102715760609261204d9201610c6c565b910152565b9080601f830112156102715781359161206a83611e66565b926120786040519485610bac565b80845260208085019160051b830101918383116102715760208101915b8383106120a457505050505090565b82356001600160401b038111610271578201906060828703601f19011261027157604051906120d282610b5b565b60208301356120e081611fd4565b825260408301356001600160401b0381116102715787602061210492860101610c6c565b60208301526060830135916001600160401b0383116102715761212f88602080969581960101610c6c565b6040820152815201920191612095565b91909160808184031261027157612154610bde565b9281356001600160401b0381116102715781612171918401610c6c565b845260208201356001600160401b0381116102715781612192918401610c6c565b602085015260408201356001600160401b03811161027157816121b6918401611fde565b604085015260608201356001600160401b038111610271576121d89201612052565b6060830152565b3590811515820361027157565b9080601f830112156102715781359161220483611e66565b926122126040519485610bac565b80845260208085019160051b830101918383116102715760208101915b83831061223e57505050505090565b82356001600160401b0381116102715782016020818703601f190112610271576040519061226b82610b76565b60208101356001600160401b03811161027157602091010186601f820112156102715780359061229a82611e66565b916122a86040519384610bac565b80835260208084019160051b830101918983116102715760208101915b8383106122e1575050509082525081526020928301920161222f565b82356001600160401b038111610271578201906060828d03601f190112610271576040519061230f82610b5b565b6020830135600281101561027157825260408301356001600160401b038111610271578d60206123419286010161213f565b602083015260608301356001600160401b03811161027157602060a0918f9501018094031261027157612372610bcd565b9183356001600160401b038111610271578e61238f918601610c6c565b835261239d602085016121df565b602084015260408401356001600160401b038111610271578e6123c191860161213f565b60408401526123d2606085016121df565b60608401526080840135926001600160401b038411610271576123fb8f6020969587960161213f565b608082015260408201528152019201916122c5565b91908260609103126102715760405161242881610b5b565b809280356fffffffffffffffffffffffffffffffff811681036102715760409182918452602081013560208501520135910152565b3590600182101561027157565b602081830312610271578035906001600160401b03821161027157016101408183031261027157612499610bed565b916124a48183611e2c565b835260408201356001600160401b03811161027157816124c5918401611efb565b602084015260608201356001600160401b03811161027157816124e99184016121ec565b6040840152608082013560608401526125058160a08401612410565b6080840152612517610100830161245d565b60a08401526101208201356001600160401b0381116102715761253a9201611e7d565b60c082015290565b90816020910312610271575190565b9080602083519182815201916020808360051b8301019401925f915b83831061257c57505050505090565b90919293946020806125ba600193601f19868203018752895190836125aa8351604084526040840190611546565b9201519084818403910152610430565b9701930193019193929061256d565b9080602083519182815201916020808360051b8301019401925f915b8383106125f457505050505090565b9091929394602080612645600193601f1986820301875289519081516126198161130a565b81526040612634858401516060878501526060840190610430565b920151906040818403910152610430565b970193019301919392906125e5565b6104659160606126d76126856126738551608086526080860190610430565b60208601518582036020870152610430565b6080836040870151868403604088015280516126a08161130a565b845260208101516126b08161130a565b602085015260408101516126c38161130a565b604085015201519181858201520190610430565b9201519060608184039101526125c9565b9080602083519283815201916020808260051b8401019401925f925b82841061271357505050505090565b9091929394601f198282030183528551602082019051916020815282518092526040810190602060408460051b8301019401925f915b81831061276a57505050505060208060019297019301940192919390612704565b9091929394603f1982820301855285519081519160028310156113145761280b826020939260019585945260406127ae858301516060878601526060850190612654565b91015191604081830391015260806127ee6127d2845160a0855260a0850190610430565b8685015115158785015260408501518482036040860152612654565b926060810151151560608401520151906080818403910152612654565b9701950193019190612749565b9060018210156113145752565b9061046591602081526128516020820183516001600160401b0360208092828151168552015116910152565b60c061288761287160208501516101406060860152610160850190612551565b6040850151848203601f190160808601526126e8565b92606081015160a08401526128c9608082015183850190604080916fffffffffffffffffffffffffffffffff8151168452602081015160208501520151910152565b6128dc60a0820151610120850190612818565b015190610140601f1982850301910152611546565b6040513d5f823e3d90fd5b15612905575050565b6001600160401b03907f65d30129000000000000000000000000000000000000000000000000000000005f521660045260245260445ffd5b91909160a08184031261027157612952610bcd565b9281356001600160401b038111610271578161296f918401610c6c565b845260208201356001600160401b0381116102715781612990918401610c6c565b602085015260408201356001600160401b03811161027157816129b4918401610c6c565b604085015260608201356001600160401b03811161027157816129d8918401610c6c565b606085015260808201356001600160401b038111610271576129fa9201610c6c565b6080830152565b610465916080612a59612a47612a35612a23865160a0875260a0870190610430565b60208701518682036020880152610430565b60408601518582036040870152610430565b60608501518482036060860152610430565b920151906080818403910152610430565b6020815260a06001600160a01b036080612ad3612aab612a95875186602089015260c0880190610430565b6020880151878203601f19016040890152610430565b6001600160401b0360408801511660608701526060870151601f198783030184880152612a01565b9401511691015290565b81604051928392833781015f815203902090565b9035601e19823603018112156102715701602081359101916001600160401b03821161027157813603831361027157565b9035601e19823603018112156102715701602081359101916001600160401b038211610271578160051b3603831361027157565b906001600160401b038235612b6a81611e06565b168152612bd6612baf612b94612b836020860186612af1565b60a0602087015260a0860191611c58565b612ba16040860186612af1565b908583036040870152611c58565b926001600160401b036060820135612bc681611e06565b1660608401526080810190612b22565b9290916080818303910152828152602081019260208160051b83010193835f91609e1982360301945b848410612c10575050505050505090565b90919293949596601f19828203018352873587811215610271576020612cc66001938783940190612cb8612cad612c92612c77612c5e612c508780612af1565b60a0885260a0880191611c58565b612c6a89880188612af1565b908783038b890152611c58565b612c846040870187612af1565b908683036040880152611c58565b612c9f6060860186612af1565b908583036060870152611c58565b926080810190612af1565b916080818503910152611c58565b990193019401929195949390612bff565b906020610465928181520190612b56565b612d0c6001612d04612cfa8480611cd5565b6080810190611cea565b905014611d1f565b612d22612d1c612cfa8380611cd5565b90611d62565b905f6020612e34612d4261098e612d398680611cd5565b84810190611d72565b612d868151848151910120612d676104c8612d5d8980611cd5565b6040810190611d72565b85815191012014825190612d7e612d5d8980611cd5565b929091611da4565b612ddd612dbd612db8612d9c612d5d8980611cd5565b9190612db0612dab8b80611cd5565b611e17565b923691610c36565b614fbd565b84612dd5612dcd828a018a611d72565b81019061246a565b930151615000565b60c0820152612dfe61095d6111a3612df58880611cd5565b86810190611d72565b906040519485809481937fa6f031bb00000000000000000000000000000000000000000000000000000000835260048301612825565b03925af1801561144257612e7a915f91613031575b50612e626109e46060612e5c8680611cd5565b01611e17565b811015612e746060612e5c8680611cd5565b906128fc565b612e93612e8f612e8a8380611cd5565b6150f0565b1590565b6130095781612f1e612eae61095d61057184612f1797611d72565b91612ec6612ebc8580611cd5565b6020810190611d72565b9590612f02612ed8612d5d8880611cd5565b612ef9612ee8612dab8b80611cd5565b94612ef1610bcd565b9b3691610c36565b8a523691610c36565b60208801526001600160401b03166040870152565b369061293d565b6060840152336080840152803b1561027157612f6d5f939184926040519586809481937f5e32b6b600000000000000000000000000000000000000000000000000000000835260048301612a6a565b03925af191821561144257612fd892612fef575b507f01e5ed58494819ef3f6480dd08e433b7c08ed75c7abdf2c22c6f04b71340a1686001600160401b03612fb8612ebc8480611cd5565b612fd2612fcb612dab8780999599611cd5565b9580611cd5565b95612add565b92612fea604051928392169582612cd7565b0390a3565b80612ffd5f61300393610bac565b80610588565b5f612f81565b50507fd08bf58b0e4eec5bfc697a4fdbb6839057fbf4dd06f1b1ce07445c0e5a654caf5f80a1565b61304a915060203d602011610dea57610ddb8183610bac565b5f612e49565b60209082604051938492833781017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960181520301902090565b6001600160a01b03604051838382376020818581017fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf60081520301902054169182156130d257505090565b611cd16040519283927fa09dbf590000000000000000000000000000000000000000000000000000000084526020600485018181520191611c58565b156106ec57565b1561311d5750565b6001600160a01b03907fbe2f2b45000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b610bdc906001600160a01b03197fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b00546001600160a01b0381163314801561322d575b61319f903390613115565b167fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b00556001600160a01b03197fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b0154167fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b0155613218615aa7565b613220615aa7565b613228615aa7565b6157ff565b5061319f6001600160a01b037fba83ed17c16070da0debaa680185af188d82c999a75962a12a40699ca48a2b01541633149050613194565b634e487b7160e01b5f52601160045260245ffd5b5f19810191908211611c2b57565b91908203918211611c2b57565b1561329c5750565b7f715fed60000000000000000000000000000000000000000000000000000000005f526201518060045260245260445ffd5b6132d86001611e66565b906132e66040519283610bac565b6001825281601f196132f86001611e66565b01905f5b82811061330857505050565b60209060405161331781610b91565b606081526060838201526060604082015260608082015260606080820152828285010152016132fc565b8051156108e85760200190565b80518210156108e85760209160051b010190565b90602082526001600160401b03815116602083015260806133ab613395602084015160a0604087015260c0860190610430565b6040840151858203601f19016060870152610430565b916001600160401b036060820151168285015201519160a0601f1982840301910152815180825260208201916020808360051b8301019401925f915b8383106133f657505050505090565b9091929394602080613414600193601f198682030187528951612a01565b970193019301919392906133e7565b9080601f8301121561027157813561343a81611e66565b926134486040519485610bac565b81845260208085019260051b820101918383116102715760208201905b83821061347457505050505090565b81356001600160401b038111610271576020916134968784809488010161293d565b815201910190613465565b91909160a081840312610271576134b6610bcd565b926134c082611e21565b845260208201356001600160401b03811161027157816134e1918401610c6c565b602085015260408201356001600160401b0381116102715781613505918401610c6c565b604085015261351660608301611e21565b606085015260808201356001600160401b038111610271576129fa9201613423565b6104659036906134a1565b602081830312610271578035906001600160401b03821161027157016101608183031261027157613572610bfc565b9161357d8183611e2c565b835260408201356001600160401b038111610271578161359e918401611efb565b602084015260608201356001600160401b03811161027157816135c29184016121ec565b6040840152608082013560608401526135de8160a08401612410565b60808401526135f0610100830161245d565b60a08401526101208201356001600160401b0381116102715781613615918401611e7d565b60c08401526101408201356001600160401b038111610271576136389201610c6c565b60e082015290565b90610465916020815261366c6020820183516001600160401b0360208092828151168552015116910152565b60e061370b6136a561368f60208601516101606060870152610180860190612551565b6040860151858203601f190160808701526126e8565b606085015160a0850152608085015180516fffffffffffffffffffffffffffffffff1660c0860152602081015160e0860152604001516101008501526136f460a0860151610120860190612818565b60c0850151848203601f1901610140860152611546565b92015190610160601f1982850301910152610430565b604080519091906137328382610bac565b6001815291601f1901825f5b82811061374a57505050565b80606060208093850101520161373e565b9061376582611e66565b6137726040519182610bac565b8281528092613783601f1991611e66565b01905f5b82811061379357505050565b806060602080938501015201613787565b602081830312610271578051906001600160401b038211610271570181601f82011215610271578051906137d782610c1b565b926137e56040519485610bac565b8284526020838301011161027157815f9260208093018386015e8301015290565b1561380d57565b7fecfef798000000000000000000000000000000000000000000000000000000005f5260045ffd5b60405190613844604083610bac565b602082527f4774d4a575993f963b1c06573736617a457abef8589178db8d10c94b4ab511ab6020830152565b613878613835565b6020815191012090565b1561388957565b7f6b2675e3000000000000000000000000000000000000000000000000000000005f5260045ffd5b3d156138db573d906138c282610c1b565b916138d06040519384610bac565b82523d5f602084013e565b606090565b156138e757565b7fadef7fb8000000000000000000000000000000000000000000000000000000005f5260045ffd5b909161392661046593604084526040840190612b56565b916020818403910152611546565b6139466001612d04612cfa8480611cd5565b613956612d1c612cfa8380611cd5565b5f6020613a9861396c61098e612d5d8780611cd5565b6139a781518481519101206139906104c86139878a80611cd5565b87810190611d72565b85815191012014825190612d7e6139878a80611cd5565b6139ca6139b96060612e5c8980611cd5565b42906109bb6060612e5c8b80611cd5565b613a55613a63613a006139fb6139ec6139e38b80611cd5565b88810190611d72565b9190612db0612dab8d80611cd5565b6154aa565b613a3e613a1d613a18613a138c80611cd5565b613538565b61551d565b9187613a36613a2e8d83018e611d72565b810190613543565b960151615000565b60c085015260405192839187830160209181520190565b03601f198101835282610bac565b60e0820152613a7b61095d6111a3612d5d8980611cd5565b906040519485809481936325d29d3160e21b835260048301613640565b03925af1801561144257613cd5575b50613abd612e8f613ab88480611cd5565b615653565b613009577f76765590e2b799b0506100f8a6610cfecab2c71e8e1f8aa981b099aff0dfdb746001600160401b03613c64935f80613bb7613afb613721565b96612f17613b77613b1561095d6105716020860186611d72565b92613b23612ebc8980611cd5565b9390613b628a612dab613b59613b48613b3f612d5d8580611cd5565b93909480611cd5565b94613b51610bcd565b993691610c36565b88523691610c36565b60208601526001600160401b03166040850152565b60608201523360808201526040519485809481937f078c4a7900000000000000000000000000000000000000000000000000000000835260048301612a6a565b03925af15f9181613cb1575b50613c7357507fb9edb487876e8be10f54e377c1a815a54ad92a6db1c9561dfe8fad2f0d1da84f613c01613bf56138b1565b611b87815115156138e0565b0390a1613c0c613835565b613c1585613341565b52613c1f84613341565b505b613c3484613c2f8380611cd5565b615739565b612fea613c44612d5d8380611cd5565b613c5e613c57612dab86809b959b611cd5565b9480611cd5565b97612add565b9460405193849316968361390f565b613c7f81511515613806565b613c9881516020830120613c91613870565b1415613882565b613ca185613341565b52613cab84613341565b50613c21565b613cce9192503d805f833e613cc68183610bac565b8101906137a4565b905f613bc3565b613ced9060203d602011610dea57610ddb8183610bac565b613aa7565b91909115613cfe575050565b611cd16040519283927f14d712470000000000000000000000000000000000000000000000000000000084526020600485018181520191611c58565b6001600160a01b03604051838382376020818581017f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496008152030190205416918215613d8457505090565b611cd16040519283927fa0db16fe0000000000000000000000000000000000000000000000000000000084526020600485018181520191611c58565b90816020910312610271575161046581611fd4565b939291602091613ded91604087526040870191611c58565b93611e628361130a565b906004116102715790600490565b90613e0f82610c1b565b613e1c6040519182610bac565b8281528092613e2d601f1991610c1b565b0190602036910137565b805191908290602001825e015f815290565b613e6b610bdc92949360208660405197889583870137840101905f8252613e37565b03601f198101845283610bac565b90919060405190613e8b602083610bac565b5f8252601f19613e9a5f610c1b565b01366020840137613eaa8461375b565b925f5b85811015613ef657600190613eda613ed486613ece8460051b880188611d72565b90613e49565b3061585f565b613ee4828861334e565b52613eef818761334e565b5001613ead565b509350505090565b90600182811c92168015613f2c575b6020831014613f1857565b634e487b7160e01b5f52602260045260245ffd5b91607f1691613f0d565b90815491613f4383611e66565b92613f516040519485610bac565b80845260208401915f5260205f20915f905b828210613f705750505050565b6040515f8554613f7f81613efe565b8084529060018116908115613ff05750600114613fb9575b5060019282613fab85946020940382610bac565b815201940191019092613f63565b5f878152602081209092505b818310613fda57505081016020016001613f97565b6001816020925483868801015201920191613fc5565b60ff191660208581019190915291151560051b8401909101915060019050613f97565b9190911561401f575050565b611cd16040519283927fdf95155a0000000000000000000000000000000000000000000000000000000084526020600485018181520191611c58565b9190916060602060405161406e81610b40565b828152015261407d8382613050565b926040519161408b83610b40565b604051945f815461409b81613efe565b808952906001811690811561411f57506001146140e3575b506140d491876140cc6104659798996001940382610bac565b875201613f36565b60208501528351511515614013565b5f838152602081209092505b81831061410557505086016020016140d46140b3565b6001816020929493945483858d01015201910191906140ef565b60ff19166020808b019190915291151560051b890190910191506140d490506140b3565b610bdc9061414f615aa7565b614157615aa7565b61415f615aa7565b614167615aa7565b613218615aa7565b6020815260c06001600160a01b0360a0612ad36141dc6141b361419d88518760208a015260e0890190610430565b6020890151888203601f190160408a0152610430565b6001600160401b0360408901511660608801526060880151601f19888303016080890152612a01565b6080870151868203601f190184880152610430565b916142086104659492604085526040850190612b56565b926020818503910152611c58565b6142208180611cd5565b6080810161422d91611cea565b61423a9150600114611d1f565b6142448180611cd5565b6080810161425191611cea565b61425a91611d62565b6142648280611cd5565b6020810161427191611d72565b61427a9161405b565b9081518051906020012061428e8480611cd5565b6040810161429b91611d72565b36906142a692610c36565b805190602001201482516142ba8580611cd5565b604081016142c791611d72565b916142d193611da4565b6142db8380611cd5565b604081016142e891611d72565b906142f38580611cd5565b6142fc90611e17565b91369061430892610c36565b9061431291615898565b9161431b613721565b92602085019361432b8587611d72565b369061433692610c36565b61433f82613341565b5261434981613341565b506143539061593b565b906143616040870187611d72565b810161436c91613543565b92602001519061437b91615000565b60c08301526040805160208082019390935291825261439a9082610bac565b60e08201526143a98480611cd5565b602081016143b691611d72565b6143bf91613d3a565b6001600160a01b03166040518080936325d29d3160e21b825260048201906143e691613640565b03815a6020945f91f18015611442576145b8575b5061440b612e8f612e8a8580611cd5565b61458f57806144a261442661095d6105718461449696611d72565b91614478614437612ebc8880611cd5565b959092612f178961448161444e612d5d8380611cd5565b6144658d61445f612dab8780611cd5565b95611d72565b989099614470610c0c565b9d3691610c36565b8c523691610c36565b60208a01526001600160401b03166040890152565b60608601523691610c36565b60808301523360a0830152803b15610271576144f15f929183926040519485809481937f428e4e170000000000000000000000000000000000000000000000000000000083526004830161416f565b03925af1801561144257612fea927ff9bab74bcdb634f4d3dd064cc42a13df056598e1c0336905d2f5750fbfb08b7b926001600160401b039261457b575b5061453d612ebc8680611cd5565b61456c614564614553612dab8a809a969a611cd5565b9461455e8a80611cd5565b99611d72565b929097612add565b956040519485941697846141f1565b80612ffd5f61458993610bac565b5f61452f565b5050507fd08bf58b0e4eec5bfc697a4fdbb6839057fbf4dd06f1b1ce07445c0e5a654caf5f80a1565b6145d09060203d602011610dea57610ddb8183610bac565b6143fa565b919091356001600160e01b0319811692600481106145f1575050565b6001600160e01b0319929350829060040360031b1b161690565b6040906001600160a01b0361046595931681528160208201520191611c58565b61466861464c5f516020615f9b5f395f51905f52546001600160a01b031690565b61465f614659855f613df7565b906145d5565b908330916159c8565b901561467357505050565b63ffffffff1615614777576146bf7401000000000000000000000000000000000000000060ff60a01b195f516020615f9b5f395f51905f525416175f516020615f9b5f395f51905f5255565b6146e361095d61095d5f516020615f9b5f395f51905f52546001600160a01b031690565b91823b1561027157614729925f808094604051968795869485937f94c7d7ee0000000000000000000000000000000000000000000000000000000085526004850161460b565b03925af1801561144257614763575b50610bdc60ff60a01b195f516020615f9b5f395f51905f5254165f516020615f9b5f395f51905f5255565b80612ffd5f61477193610bac565b5f614738565b62d1953b60e31b5f526001600160a01b031660045260245ffd5b604051906147a0604083610bac565b600882527f6368616e6e656c2d0000000000000000000000000000000000000000000000006020830152565b604051906147db604083610bac565b600782527f636c69656e742d000000000000000000000000000000000000000000000000006020830152565b8051600110156108e85760210190565b9081518110156108e8570160200190565b8051600481109081156149db575b506149c55761484c614846614791565b82615a55565b80156149ca575b6149c5575f5b81518110156149be576148a76148a161489b6148758486614817565b517fff000000000000000000000000000000000000000000000000000000000000001690565b60f81c90565b60ff1690565b606181101590816149b2575b8115614994575b8115614976575b8115614937575b81156148e2575b506148da5750505f90565b600101614859565b602381149150811561492c575b8115614921575b8115614916575b811561490b575b505f6148cf565b603e9150145f614904565b603c811491506148fd565b605d811491506148f6565b605b811491506148ef565b9050602e8114801561496c575b8015614962575b8015614958575b906148c8565b50602d8114614952565b50602b811461494b565b50605f8114614944565b9050604181101580614989575b906148c1565b50605a811115614983565b90506030811015806149a7575b906148ba565b5060398111156149a1565b607a81111591506148b3565b5050600190565b505f90565b506149d66148466147cc565b614853565b60809150115f614836565b60206149f89160405192838092613e37565b7f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960181520301902090565b818110614a2d575050565b5f8155600101614a22565b9190601f8111614a4757505050565b610bdc925f5260205f20906020601f840160051c83019310614a71575b601f0160051c0190614a22565b9091508190614a64565b908160011b9180830460021490151715611c2b57565b908160041b9180830460101490151715611c2b57565b9092916001600160401b038111610b3b57614acc81614ac68454613efe565b84614a38565b5f601f8211600114614b0a578190614afb9394955f92614aff575b50508160011b915f199060031b1c19161790565b9055565b013590505f80614ae7565b601f19821694614b1d845f5260205f2090565b915f5b878110614b57575083600195969710614b3e575b505050811b019055565b01355f19600384901b60f8161c191690555f8080614b34565b90926020600181928686013581550194019101614b20565b680100000000000000008311610b3b578054838255808410614bd4575b5090614b9c81925f5260205f2090565b905f925b848410614bae575050505050565b6001602082614bc8614bc1849587611d72565b9088614aa7565b01930193019291614ba0565b815f528360205f2091820191015b818110614bef5750614b8c565b80614bfc60019254613efe565b80614c09575b5001614be2565b601f81118314614c1e57505f81555b5f614c02565b614c409083601f614c32855f5260205f2090565b920160051c82019101614a22565b5f8181526020812081835555614c18565b90614c5c8180611d72565b906001600160401b038211610b3b57614c7f82614c798654613efe565b86614a38565b5f90601f8311600114614ccd5792614cb783610bdc9694614cc4946001975f92614aff5750508160011b915f199060031b1c19161790565b83555b6020810190611cea565b92909101614b6f565b601f19831691614ce0865f5260205f2090565b925f5b818110614d28575093614cc49360019693879383610bdc9a9810614d0f575b505050811b018355614cba565b01355f19600384901b60f8161c191690555f8080614d02565b91936020600181928787013581550195019201614ce3565b93929190614d5690606086526060860190610430565b8481036020860152614d87614d7c614d6e8480612af1565b604085526040850191611c58565b926020810190612b22565b90916020818503910152808352602083019260208260051b82010193835f925b848410614dcb57505050505050906040610bdc929401906001600160a01b03169052565b909192939495602080614df3600193601f19868203018852614ded8b88612af1565b90611c58565b9801940194019294939190614da7565b92916001600160a01b03604051602081614e1d8189613e37565b7f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a4496008152030190205416614ee95792614ee47f0ecded31ecd211a73abf0fb3bc09150bbe321a05550fbe29ea0f16b6e25fbfa89394614ec6604051602081614e848188613e37565b7f515a8336edcaab4ae6524d41223c1782132890f89189ba6632107a7b5a44960081520301902080546001600160a01b0319166001600160a01b038416179055565b614ed884614ed3856149e6565b614c51565b60405193849384614d40565b0390a1565b6040517f87dfb2670000000000000000000000000000000000000000000000000000000081526020600482015280611cd16024820187610430565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005c614f705760017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d565b7f3ee5aeb5000000000000000000000000000000000000000000000000000000005f5260045ffd5b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005d565b600961046591614fe0936001600160c01b03196040519586936020850190613e37565b91600160f91b835260c01b1660018201520301601f198101835282610bac565b9190918051156150ba57615014815161375b565b915f5b6150218351613279565b81101561505257806150356001928561334e565b51615040828761334e565b5261504b818661334e565b5001615017565b50926150a46150b69261509e61508a936150906150786150728551613279565b8561334e565b51916040519687936020850190613e37565b90613e37565b03601f198101855284610bac565b51613279565b906150af828561334e565b528261334e565b5090565b611cd1906040519182917fa7c34e4f0000000000000000000000000000000000000000000000000000000083526004830161159e565b6151116139fb6151036020840184611d72565b91908435926104c884611e06565b6020815191012090815f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f205480156151d157615168613a1861515e613a1836866134a1565b83149336906134a1565b91156151a35750505f9081527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600602052604081205b55600190565b7f3f87a2ec000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b5050505f90565b906001600160a01b036151f8602092959495604085526040850190610430565b9416910152565b91906001600160a01b036040516020816152198188613e37565b7fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf60081520301902054166152cd577fa6ec8e860960e638347460dc632fbe0175c51a5ca130e336138bbe26ff30449991926152be60405160208161527c8186613e37565b7fc5779f3c2c21083eefa6d04f6a698bc0d8c10db124ad5e0df6ef394b6d7bf60081520301902080546001600160a01b0319166001600160a01b038516179055565b614ee4604051928392836151d8565b6040517f837f46a60000000000000000000000000000000000000000000000000000000081526020600482015280611cd16024820186610430565b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300541661533357565b7fd93c0665000000000000000000000000000000000000000000000000000000005f5260045ffd5b60209082604051938492833781017f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e07476018152030190206001600160401b038154166001600160401b038114611c2b57600101906001600160401b0382166001600160401b031982541617905590565b60208101906153e482516001600160401b03835116906154aa565b6020815191012091825f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205261542f60405f20541591516001600160401b03845116906154aa565b901561546d575061543f9061551d565b905f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f2055565b611cd1906040519182917f91ffd9240000000000000000000000000000000000000000000000000000000083526020600484018181520190610430565b6009610465916154cd936001600160c01b03196040519586936020850190613e37565b917f0100000000000000000000000000000000000000000000000000000000000000835260c01b1660018201520301601f198101835282610bac565b6155169060209392613e37565b9081520190565b9061552661161b565b5f905b6080840151805183101561556b579061556361555061554a8560019561334e565b51615c06565b91613a5560405193849260208401615509565b910190615529565b50905091909160205f615588604085015160405191828092613e37565b039060025afa156114425760205f6155e0613a556155d46155b5606085519801516001600160401b031690565b6040519283918783016001600160c01b031960089260c01b1681520190565b60405191828092613e37565b039060025afa156114425760205f61560081519360405191828092613e37565b039060025afa15611442576156435f916155d4602094613a5585516040519485938985019160619391600160f91b84526001840152602183015260418201520190565b039060025afa15611442575f5190565b615686615677612db86156696040850185611d72565b91908535926104c884611e06565b602081519101209136906134a1565b60405161569b81613a55602082019485613362565b51902090805f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f20548281146151d15780615709575061519d905f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f2090565b90507f657b94fe000000000000000000000000000000000000000000000000000000005f5260045260245260445ffd5b9060408201916157b261576e6157696157a96157558786611d72565b61576987959295359586926104c884611e06565b615898565b6020815191012096875f527f1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e074760060205260405f20541595611d72565b6104c884611e06565b90156157c2575061543f9061593b565b611cd1906040519182917f40470d740000000000000000000000000000000000000000000000000000000083526020600484018181520190610430565b60206001600160a01b037f2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad9216806001600160a01b03195f516020615f9b5f395f51905f525416175f516020615f9b5f395f51905f5255604051908152a1565b5f8061046593602081519101845af46158766138b1565b91615d8d565b610bdc90613e6b61508a94936040519586936020850190613e37565b6009610465916158bb936001600160c01b03196040519586936020850190613e37565b917f0300000000000000000000000000000000000000000000000000000000000000835260c01b1660018201520301601f198101835282610bac565b156158fe57565b7f760d6a9b000000000000000000000000000000000000000000000000000000005f5260045ffd5b90600161046592600160f91b81520190613e37565b90615948825115156158f7565b61595061161b565b915f925b81518410156159a85760205f61597a61596d878661334e565b5160405191828092613e37565b039060025afa15611442576001906159a05f5191613a5560405193849260208401615509565b930192615954565b60209293505f9150613a556155d461564392604051928391878301615926565b5f9060409295939582966001600160e01b031984976001600160a01b038751938160208601967fb700961300000000000000000000000000000000000000000000000000000000885216602486015216604484015216606482015260648152615a32608482610bac565b8380528360205251915afa615a4357565b9150505f51906020518060201c150290565b80518251808210615a9f57818082109118028082189114158102811890818103908111611c2b576020615a8782613e05565b92818401940101835e51902090602081519101201490565b505050505f90565b60ff5f516020615fbb5f395f51905f525460401c1615615ac357565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f5260045ffd5b9060028201809211611c2b57565b9060018201809211611c2b57565b91908201809211611c2b57565b8015611c2b575f190190565b90813b15615bea576001600160a01b038216806001600160a01b03197f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5416177f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a2805115615bb957615bb69161585f565b50565b505034615bc257565b7fb398979f000000000000000000000000000000000000000000000000000000005f5260045ffd5b6001600160a01b0382634c9c8ce360e01b5f521660045260245ffd5b60205f615c1a835160405191828092613e37565b039060025afa15611442575f5160205f615c3d8285015160405191828092613e37565b039060025afa15611442575f519160205f615c62604084015160405191828092613e37565b039060025afa15611442575f519260205f615c87606085015160405191828092613e37565b039060025afa156114425760205f615cab6080825195015160405191828092613e37565b039060025afa15611442576155d45f93613a55602096615643958751916040519687958b8701939160a0959391855260208501526040840152606083015260808201520190565b805182118015615d86575b615d4a576001821180615d52575b615d16901515614a7b565b60280180602811611c2b57615d2b5f84613287565b03615d4a576001600160a01b0392915f615d4492615e19565b90921690565b50505f905f90565b5060208101517fffff0000000000000000000000000000000000000000000000000000000000001661060f60f31b14615d0b565b505f615cfd565b90615dca5750805115615da257602081519101fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f5260045ffd5b81511580615e10575b615ddb575090565b6001600160a01b03907f9996b315000000000000000000000000000000000000000000000000000000005f521660045260245ffd5b50803b15615dd3565b929092615e2584615af9565b831180615ec1575b615e3f615e4691949293941515614a7b565b5f95615b07565b915b818310615e585750505060019190565b9092919360ff615e95615e90602088860101517fff000000000000000000000000000000000000000000000000000000000000001690565b615f28565b1690600f8211615eb65790615eab600192614a91565b019401919290615e48565b505f94508493505050565b50615e46615e3f61060f60f31b7fffff000000000000000000000000000000000000000000000000000000000000615f1e602089870101517fffff0000000000000000000000000000000000000000000000000000000000001690565b1614915050615e2d565b60f81c602f811180615f90575b15615f4457602f190160ff1690565b6060811180615f86575b15615f5d576056190160ff1690565b6040811180615f7c575b15615f76576036190160ff1690565b5060ff90565b5060478110615f67565b5060678110615f4e565b50603a8110615f3556fef3177357ab46d8af007ab3fdb9af81da189e1068fefdc0073dca88a2cab40a00f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00a164736f6c634300081c000a1b9a2a4db4fdabb5b32192c8c2b90c1835653cf1aedfe8794a443b5c718f4c66d68750c3950334e058f8ccdd1a08540c778394eba7c6aa6b16155c1c6d05a9cef0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
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
func (_ContractICS26Router *ContractICS26RouterTransactor) ExecuteClientMigration(opts *bind.TransactOpts, clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "executeClientMigration", clientId, counterpartyInfo, client)
}

// ExecuteClientMigration is a paid mutator transaction binding the contract method 0xee5c877c.
//
// Solidity: function executeClientMigration(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterSession) ExecuteClientMigration(clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.ExecuteClientMigration(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// ExecuteClientMigration is a paid mutator transaction binding the contract method 0xee5c877c.
//
// Solidity: function executeClientMigration(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) ExecuteClientMigration(clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
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

// ProposeClientMigration is a paid mutator transaction binding the contract method 0x0cbdf9f6.
//
// Solidity: function proposeClientMigration(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterTransactor) ProposeClientMigration(opts *bind.TransactOpts, clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.contract.Transact(opts, "proposeClientMigration", clientId, counterpartyInfo, client)
}

// ProposeClientMigration is a paid mutator transaction binding the contract method 0x0cbdf9f6.
//
// Solidity: function proposeClientMigration(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterSession) ProposeClientMigration(clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.ProposeClientMigration(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
}

// ProposeClientMigration is a paid mutator transaction binding the contract method 0x0cbdf9f6.
//
// Solidity: function proposeClientMigration(string clientId, (string,bytes[]) counterpartyInfo, address client) returns()
func (_ContractICS26Router *ContractICS26RouterTransactorSession) ProposeClientMigration(clientId string, counterpartyInfo IICS02ClientMsgsCounterpartyInfo, client common.Address) (*types.Transaction, error) {
	return _ContractICS26Router.Contract.ProposeClientMigration(&_ContractICS26Router.TransactOpts, clientId, counterpartyInfo, client)
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
