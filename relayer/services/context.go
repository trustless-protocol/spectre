package services

import (
	"log"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// CosmosEndpoint is the transport to a Cosmos chain. It deliberately carries
// no counterparty state or transaction configuration.
type CosmosEndpoint struct{ Client *rpchttp.HTTP }

// EVMContracts contains the contracts used by the EVM side of a relay path.
// Zero addresses retain the former "not configured" behaviour; consumers that
// require an address validate it at their boundary.
type EVMContracts struct {
	Router, SignatureVerifier, Membership, Misbehaviour, UpdateClient, RoleManager, SpectreClient common.Address
}

// EVMEndpoint is the transport and contract set for one EVM chain.
type EVMEndpoint struct {
	Client       *ethclient.Client
	WSURL        string
	BeaconAPIURL string
	Contracts    EVMContracts
}

// ClientIDs name the two counterparty clients without tying either identifier
// to a shared process-wide context.
type ClientIDs struct {
	CosmosOnEVM string
	EVMOnCosmos string
}

func (e CosmosEndpoint) CosmosClient() *rpchttp.HTTP { return e.Client }

func (e EVMEndpoint) EthClient() *ethclient.Client { return e.Client }
func (e EVMEndpoint) EthWsURL() string             { return e.WSURL }
func (e EVMEndpoint) SignatureVerifierContract() *common.Address {
	return &e.Contracts.SignatureVerifier
}
func (e EVMEndpoint) MembershipContract() *common.Address {
	return &e.Contracts.Membership
}
func (e EVMEndpoint) MisbehaviourContract() *common.Address {
	return &e.Contracts.Misbehaviour
}
func (e EVMEndpoint) UpdateClientContract() *common.Address {
	return &e.Contracts.UpdateClient
}
func (e EVMEndpoint) RoleManagerAddress() *common.Address {
	return &e.Contracts.RoleManager
}
func (e EVMEndpoint) RouterContract() *common.Address { return &e.Contracts.Router }
func (e EVMEndpoint) SpectreClientContract() *common.Address {
	return &e.Contracts.SpectreClient
}

// RelayDeps contains the complete wiring assembled by the command layer for a
// relay process. It is intentionally a composition-root value: adapters,
// subscribers, workers, and transaction handlers receive only the endpoint,
// identifier, and configuration values they actually use.
type RelayDeps struct {
	Cosmos CosmosEndpoint
	EVM    EVMEndpoint
	IDs    ClientIDs
	Config Config
	Logger *log.Logger
}
