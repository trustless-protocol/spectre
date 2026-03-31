package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Timestamp struct {
	mtx                sync.Mutex
	LatestUpdateTime   time.Time
	LatestUpdateHeight uint64
}

type Context struct {
	Logger *log.Logger
	Config Config

	latestEthTimestamp    *Timestamp
	latestCosmosTimestamp *Timestamp

	cosmosClient *rpchttp.HTTP
	ethClient    *ethclient.Client
	beaconAPIURL string

	// Ethereum light client configuration
	ethClientID  string
	verifier     *common.Address
	membership   *common.Address
	misbehaviour *common.Address
	updateClient *common.Address
	roleManager  *common.Address
	ics26Router  *common.Address
	ics07Client  *common.Address
}

func NewCtx(cosmosClient *rpchttp.HTTP, ethClient *ethclient.Client) Context {
	return Context{
		Logger:       log.Default(),
		cosmosClient: cosmosClient,
		ethClient:    ethClient,
		beaconAPIURL: "",
		ethClientID:  "",
		latestEthTimestamp: &Timestamp{
			LatestUpdateTime:   time.Now(),
			LatestUpdateHeight: 0,
		},
		latestCosmosTimestamp: &Timestamp{
			LatestUpdateTime:   time.Now(),
			LatestUpdateHeight: 0,
		},
	}
}

func NewCtxWithBeacon(cosmosClient *rpchttp.HTTP, ethClient *ethclient.Client, beaconAPIURL string, ethClientID string) Context {
	return Context{
		Logger:       log.Default(),
		cosmosClient: cosmosClient,
		ethClient:    ethClient,
		beaconAPIURL: beaconAPIURL,
		ethClientID:  ethClientID,
		latestEthTimestamp: &Timestamp{
			LatestUpdateTime:   time.Now(),
			LatestUpdateHeight: 0,
		},
		latestCosmosTimestamp: &Timestamp{
			LatestUpdateTime:   time.Now(),
			LatestUpdateHeight: 0,
		},
	}
}

func (c *Context) EthClient() *ethclient.Client {
	return c.ethClient
}

func (c *Context) CosmosClient() *rpchttp.HTTP {
	return c.cosmosClient
}

func (c *Context) BeaconAPIURL() string {
	return c.beaconAPIURL
}

func (c *Context) EthClientID() string {
	return c.ethClientID
}

func (c *Context) SetAddresses(ics26Router, verifier, membership, misbehaviour, updateClient, roleManager string) {
	ics26RouterAddr := common.HexToAddress(ics26Router)
	verifierAddr := common.HexToAddress(verifier)
	membershipAddr := common.HexToAddress(membership)
	misbehaviourAddr := common.HexToAddress(misbehaviour)
	updateClientAddr := common.HexToAddress(updateClient)
	roleManagerAddr := common.HexToAddress(roleManager)

	c.ics26Router = &ics26RouterAddr
	c.verifier = &verifierAddr
	c.membership = &membershipAddr
	c.misbehaviour = &misbehaviourAddr
	c.updateClient = &updateClientAddr
	c.roleManager = &roleManagerAddr
}

func (c *Context) SetClient(client common.Address) {
	c.ics07Client = &client
}
func (c *Context) VerifierContract() *common.Address {
	return c.verifier
}

func (c *Context) MembershipContract() *common.Address {
	return c.membership
}

func (c *Context) MisbehaviourContract() *common.Address {
	return c.misbehaviour
}

func (c *Context) UpdateClientContract() *common.Address {
	return c.updateClient
}

func (c *Context) RoleManagerAddress() *common.Address {
	return c.roleManager
}

func (c *Context) RouterContract() *common.Address {
	return c.ics26Router
}

func (c *Context) ClientContract() *common.Address {
	return c.ics07Client
}

func (c *Context) StopClient() {
	err := c.cosmosClient.Stop()
	if err != nil {
		panic(fmt.Errorf("failed to terminate cosmos client: %v", err))
	}
	c.ethClient.Close()
}

func (c *Context) LatestCosmosTimestamp() *Timestamp {
	return c.latestCosmosTimestamp
}

func (c *Context) LatestEthTimestamp() *Timestamp {
	return c.latestEthTimestamp
}
