package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"

	"relayer/chain/cosmos"
	"relayer/chain/evm"
	relayerclient "relayer/client"
	"relayer/prover"
	"relayer/relay"
	"relayer/services"
)

// Cosmos→L2 reuses the Cosmos→ETH machinery verbatim: an L2 rollup that hosts a
// SpectreClient (ZK Tendermint light client) + ICS26Router is, to the relayer, just
// another EVM destination. Per the modular-refactor decision, the `groth16` builder
// covers "Cosmos → any EVM" — the relayer proves Cosmos headers and submits them to
// the destination's SpectreClient, whether that EVM chain is L1 Ethereum or an L2. The
// only differences from Cosmos→ETH are the endpoint (the L2 exec RPC + the L2's
// contract addresses) and that there is NO reverse beacon direction on this path (the
// L2→Cosmos return path is the separate optimistic l2rollup adapter).
//
// This is distinct from the L2→Cosmos return path, which is permissioned re-execution
// attestation (the executor), not a ZK light client.

// buildCosmosToL2Dest builds one Cosmos→L2 destination: its Tendermint RPC (source
// side) and the L2's exec RPC + SpectreClient + router client id (destination side),
// sharing the passed-in prover. It returns the built Services, composition-root
// dependencies, and a cleanup that stops the chain clients. Unlike buildCosmosToEthSource it needs no
// beacon URL and no eth-client-on-Cosmos id — those belong to the ETH→Cosmos reverse
// direction, which this path does not run.
func buildCosmosToL2Dest(
	logger *zap.Logger,
	c2l cosmosToEthConfig,
	batchCfg services.BatchConfig,
	p *prover.EcipProver,
	txHandler services.TransactionHandler,
	pendingStateDir string,
) (*services.Services, services.RelayDeps, func(), error) {
	var zero services.RelayDeps

	// Config-only validation runs first, for the reason spelled out in
	// buildCosmosToEthSource: nothing below here has a resource to release.
	//
	// No beacon URL and no eth-client-on-Cosmos id: the Cosmos→L2 direction consumes
	// neither (both belong to the ETH→Cosmos beacon path).
	if c2l.ICS26ClientID == "" {
		return nil, zero, nil, fmt.Errorf("ics26_client_id is required in cosmos_to_l2 config")
	}
	if c2l.SpectreClient == "" {
		return nil, zero, nil, fmt.Errorf("spectre_client address is required in cosmos_to_l2 config")
	}
	if c2l.EthWsUrl != "" {
		if !strings.HasPrefix(c2l.EthWsUrl, "ws://") && !strings.HasPrefix(c2l.EthWsUrl, "wss://") {
			return nil, zero, nil, fmt.Errorf("eth_ws_url must use ws:// or wss://, got: %s", c2l.EthWsUrl)
		}
	}
	// Parses FETCH_TIMEOUT, so it can fail on a malformed override.
	cosmosConfig, err := buildCosmosConfig(c2l, batchCfg)
	if err != nil {
		return nil, zero, nil, err
	}

	ethClient, err := relayerclient.DialEthRPC(context.Background(), c2l.EthRpcUrl, relayerclient.DefaultRPCTimeout)
	if err != nil {
		return nil, zero, nil, fmt.Errorf("failed to connect to L2 exec rpc: %w", err)
	}

	var ethWsClient *ethclient.Client
	if c2l.EthWsUrl != "" {
		ethWsClient, err = relayerclient.DialEthRPC(context.Background(), c2l.EthWsUrl, relayerclient.DefaultRPCTimeout)
		if err != nil {
			ethClient.Close()
			return nil, zero, nil, fmt.Errorf("failed to connect to L2 exec ws: %w", err)
		}
	}

	cosmosClient, err := relayerclient.DialCosmosRPC(c2l.TmRpcUrl, "/websocket", relayerclient.DefaultRPCTimeout)
	if err != nil {
		ethClient.Close()
		if ethWsClient != nil {
			ethWsClient.Close()
		}
		return nil, zero, nil, fmt.Errorf("failed to create Cosmos RPC client: %w", err)
	}

	if err := cosmosClient.Start(); err != nil {
		ethClient.Close()
		if ethWsClient != nil {
			ethWsClient.Close()
		}
		return nil, zero, nil, fmt.Errorf("failed to start Cosmos WS client: %w", err)
	}

	deps := services.RelayDeps{
		Cosmos: services.CosmosEndpoint{Client: cosmosClient},
		EVM: services.EVMEndpoint{Client: ethClient, WSURL: c2l.EthWsUrl, Contracts: services.EVMContracts{
			Router: common.HexToAddress(c2l.ICS26Address), SignatureVerifier: common.HexToAddress(c2l.SignatureVerifier),
			Membership: common.HexToAddress(c2l.Membership), Misbehaviour: common.HexToAddress(c2l.Misbehaviour),
			UpdateClient: common.HexToAddress(c2l.UpdateClient), RoleManager: common.HexToAddress(c2l.ICS26Address), SpectreClient: common.HexToAddress(c2l.SpectreClient),
		}},
		IDs: services.ClientIDs{CosmosOnEVM: c2l.ICS26ClientID}, Config: cosmosConfig, Logger: log.Default(),
	}
	cleanup := func() {
		if err := cosmosClient.Stop(); err != nil {
			log.Printf("failed to terminate cosmos client: %v", err)
		}
		ethClient.Close()
		if ethWsClient != nil {
			ethWsClient.Close()
		}
	}

	logger.Sugar().Infof("cosmos->l2 dest %q: spectre_client=%s l2_rpc=%s tm=%s",
		c2l.ICS26ClientID, c2l.SpectreClient, c2l.EthRpcUrl, c2l.TmRpcUrl)

	recoveryState, err := services.LoadRecoveryState(services.RecoveryStatePath())
	if err != nil {
		cleanup()
		return nil, zero, nil, fmt.Errorf("load recovery state: %w", err)
	}
	var svc *services.Services
	if pendingStateDir != "" {
		svc, err = services.NewWithPendingState(txHandler, p, cosmosConfig, pendingStateDir, recoveryState)
		if err != nil {
			cleanup()
			return nil, zero, nil, err
		}
	} else {
		svc = services.New(txHandler, p, cosmosConfig, recoveryState)
	}
	return svc, deps, cleanup, nil
}

// runCosmosToL2Engine drives the Cosmos→L2 direction of one destination: the same
// pipeline as the Cosmos→ETH half of runAdapterEngine (Cosmos source → EVM
// SpectreClient destination → groth16 builder, with timeout scanning + pinned-set
// rotation), but WITHOUT the ETH→Cosmos beacon module. It returns nil on clean context
// cancellation, or the module's first fatal error.
func runCosmosToL2Engine(ctx context.Context, svc *services.Services, deps services.RelayDeps) error {
	worker := svc.Worker()
	bb := svc.BatchBuilder
	cfg := svc.CosmosConfig()

	// Cosmos-origin sends aren't tracked by the subscriber; the module records each
	// one so ScanCosmosTimeouts can refund it, and removes it once received on the L2.
	trackCosmosPending := func(raw []byte, height uint64) bool {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[cosmos->l2] track pending: decode packet: %v", err)
			return false
		}
		return svc.TrackCosmosPending(pkt, height)
	}
	untrackCosmosPending := func(raw []byte) {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[cosmos->l2] untrack pending: decode packet: %v", err)
			return
		}
		svc.UntrackCosmosPending(pkt)
	}

	// Pinned-set rotation cadence from the on-chain trusting period (a derivation
	// failure is FATAL, as in the Cosmos→ETH path, so the client can't silently expire).
	periodicUpdateInterval, err := svc.PinnedSetRotationInterval(ctx, deps.EVM)
	if err != nil {
		if isShutdownErr(err) && ctx.Err() != nil {
			return nil
		}
		return fmt.Errorf("cosmos->l2: derive pinned-set rotation interval: %w", err)
	}
	initialRotationDelay, err := svc.PinnedSetRotationDueIn(ctx, deps.Cosmos, deps.EVM)
	if err != nil {
		if isShutdownErr(err) && ctx.Err() != nil {
			return nil
		}
		log.Printf("[cosmos->l2] derive initial rotation delay: %v; rotating on startup", err)
		initialRotationDelay = 0
	}

	module := relay.NewModule(
		"cosmos->l2",
		deps.IDs.CosmosOnEVM,
		cosmos.NewSource(deps.Cosmos, deps.EVM, deps.IDs, deps.Config.FetchTimeout, deps.Config.BatchConfig, deps.Logger, bb, svc.RecoveryState()),
		evm.NewDestination(worker, deps.Cosmos, deps.EVM, deps.IDs.CosmosOnEVM),
		cosmos.NewGroth16Builder(worker, deps.Cosmos, deps.EVM, deps.Config.FetchTimeout, deps.Config.RotationThreshold, cfg.ProofType, cfg.TrustLevel),
		relay.WithTimeoutScanner(0, func(c context.Context) {
			svc.ScanCosmosTimeouts(c, deps.Cosmos, deps.EVM, deps.IDs.EVMOnCosmos)
		}),
		relay.WithPacketTracker(trackCosmosPending, untrackCosmosPending),
		relay.WithPeriodicUpdate(periodicUpdateInterval, initialRotationDelay, func(c context.Context) error {
			return svc.RotatePinnedSet(c, deps.Cosmos, deps.EVM, deps.IDs.CosmosOnEVM)
		}),
	)

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go svc.RunQueueReporter(runCtx)
	if err := module.Run(runCtx); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
