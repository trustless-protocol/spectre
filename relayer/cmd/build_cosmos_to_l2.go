package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"

	"relayer/chain/cosmos"
	"relayer/chain/evm"
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
// sharing the passed-in prover. It returns the built Services, its Context, and a
// cleanup that stops the chain clients. Unlike buildCosmosToEthSource it needs no
// beacon URL and no eth-client-on-Cosmos id — those belong to the ETH→Cosmos reverse
// direction, which this path does not run.
func buildCosmosToL2Dest(
	logger *zap.Logger,
	c2l cosmosToEthConfig,
	batchCfg services.BatchConfig,
	p *prover.EcipProver,
	txHandler services.TransactionHandler,
) (*services.Services, services.Context, func(), error) {
	var zero services.Context

	ethClient, err := ethclient.Dial(c2l.EthRpcUrl)
	if err != nil {
		return nil, zero, nil, fmt.Errorf("failed to connect to L2 exec rpc: %w", err)
	}

	var ethWsClient *ethclient.Client
	if c2l.EthWsUrl != "" {
		if !strings.HasPrefix(c2l.EthWsUrl, "ws://") && !strings.HasPrefix(c2l.EthWsUrl, "wss://") {
			return nil, zero, nil, fmt.Errorf("eth_ws_url must use ws:// or wss://, got: %s", c2l.EthWsUrl)
		}
		ethWsClient, err = ethclient.Dial(c2l.EthWsUrl)
		if err != nil {
			return nil, zero, nil, fmt.Errorf("failed to connect to L2 exec ws: %w", err)
		}
	}

	cosmosClient, err := rpchttp.New(c2l.TmRpcUrl, "/websocket")
	if err != nil {
		return nil, zero, nil, fmt.Errorf("failed to create Cosmos RPC client: %w", err)
	}

	// No beacon URL and no eth-client-on-Cosmos id: the Cosmos→L2 direction consumes
	// neither (both belong to the ETH→Cosmos beacon path).
	ctx := services.NewCtxWithBeacon(cosmosClient, ethClient, ethWsClient, c2l.EthWsUrl, "", "")

	if c2l.ICS26ClientID == "" {
		return nil, zero, nil, fmt.Errorf("ics26_client_id is required in cosmos_to_l2 config")
	}
	ctx.SetCosmosRouterClientID(c2l.ICS26ClientID)
	ctx.SetAddresses(
		c2l.ICS26Address, c2l.SignatureVerifier, c2l.Membership,
		c2l.Misbehaviour, c2l.UpdateClient, c2l.ICS26Address,
	)
	if c2l.SpectreClient == "" {
		return nil, zero, nil, fmt.Errorf("spectre_client address is required in cosmos_to_l2 config")
	}
	ctx.SetClient(common.HexToAddress(c2l.SpectreClient))

	if err := cosmosClient.Start(); err != nil {
		return nil, zero, nil, fmt.Errorf("failed to start Cosmos WS client: %w", err)
	}
	cleanup := func() { ctx.StopClient() }

	cosmosConfig := buildCosmosConfig(c2l, batchCfg)
	ctx.Config = cosmosConfig

	logger.Sugar().Infof("cosmos->l2 dest %q: spectre_client=%s l2_rpc=%s tm=%s",
		c2l.ICS26ClientID, c2l.SpectreClient, c2l.EthRpcUrl, c2l.TmRpcUrl)

	svc := services.New(txHandler, p, cosmosConfig)
	return svc, ctx, cleanup, nil
}

// runCosmosToL2Engine drives the Cosmos→L2 direction of one destination: the same
// pipeline as the Cosmos→ETH half of runAdapterEngine (Cosmos source → EVM
// SpectreClient destination → groth16 builder, with timeout scanning + pinned-set
// rotation), but WITHOUT the ETH→Cosmos beacon module. It returns nil on clean context
// cancellation, or the module's first fatal error.
func runCosmosToL2Engine(ctx context.Context, svc *services.Services, dstCtx services.Context) error {
	worker := svc.Worker()
	bb := svc.BatchBuilder
	cfg := svc.CosmosConfig()

	// Cosmos-origin sends aren't tracked by the subscriber; the module records each
	// one so ScanCosmosTimeouts can refund it, and removes it once received on the L2.
	trackCosmosPending := func(raw []byte, height uint64) {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[adapter cosmos->l2] track pending: decode packet: %v", err)
			return
		}
		svc.TrackCosmosPending(pkt, height)
	}
	untrackCosmosPending := func(raw []byte) {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[adapter cosmos->l2] untrack pending: decode packet: %v", err)
			return
		}
		svc.UntrackCosmosPending(pkt)
	}

	// Pinned-set rotation cadence from the on-chain trusting period (a derivation
	// failure is FATAL, as in the Cosmos→ETH path, so the client can't silently expire).
	periodicUpdateInterval, err := svc.PinnedSetRotationInterval(dstCtx)
	if err != nil {
		return fmt.Errorf("cosmos->l2: derive pinned-set rotation interval: %w", err)
	}
	initialRotationDelay, err := svc.PinnedSetRotationDueIn(dstCtx)
	if err != nil {
		log.Printf("[adapter cosmos->l2] derive initial rotation delay: %v; rotating on startup", err)
		initialRotationDelay = 0
	}

	module := relay.NewModule(
		"cosmos->l2",
		dstCtx.CosmosRouterClientID(),
		cosmos.NewSource(dstCtx, bb),
		evm.NewDestination(worker, dstCtx),
		cosmos.NewGroth16Builder(worker, dstCtx, cfg.ProofType, cfg.TrustLevel),
		relay.WithTimeoutScanner(0, func(c context.Context) { svc.ScanCosmosTimeouts(c, dstCtx) }),
		relay.WithPacketTracker(trackCosmosPending, untrackCosmosPending),
		relay.WithPeriodicUpdate(periodicUpdateInterval, initialRotationDelay, func(c context.Context) error {
			return svc.RotatePinnedSet(c, dstCtx)
		}),
	)

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := module.Run(runCtx); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
