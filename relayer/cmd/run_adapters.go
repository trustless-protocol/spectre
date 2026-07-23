package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"relayer/chain/cosmos"
	"relayer/chain/evm"
	"relayer/relay"
	"relayer/services"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// runAdapterEngine drives one source's bidirectional relay using the chain-adapter
// RelayModule instead of services.StartLoop. It is the cutover target: same
// battle-tested pipeline (proof gen, gap recovery, timeout scanning) reused via
// adapters, but orchestrated by the generic module.
//
// Two modules run concurrently on a shared child context and the source's shared
// services.BatchBuilder (SubscribeCosmos/SubscribeEth push to disjoint queues, so
// the two drains never collide):
//
//   - Cosmos→ETH ("groth16"): Cosmos source → SpectreClient-on-ETH destination.
//     The subscriber does NOT track Cosmos-origin sends, so the module tracks them
//     (WithPacketTracker) to feed ScanCosmosTimeouts.
//   - ETH→Cosmos ("beacon"): ETH source → 08-wasm-on-Cosmos destination. The ETH
//     subscriber already records EthPendingTracker entries, so the module adds no
//     tracker (avoiding a redundant double-add); ScanEthTimeouts drains that
//     tracker.
//
// It returns nil on clean context cancellation, or the first module's fatal error
// (cancelling the other).
func runAdapterEngine(ctx context.Context, svc *services.Services, srcCtx services.Context) error {
	worker := svc.Worker()
	bb := svc.BatchBuilder
	cfg := svc.CosmosConfig()

	// Cosmos-origin sends aren't tracked by the subscriber; the module records
	// each one (proto bytes -> Packet) so ScanCosmosTimeouts can refund it, and
	// removes it once the send is received on ETH (it can no longer time out).
	trackCosmosPending := func(raw []byte, height uint64) {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[adapter cosmos->eth] track pending: decode packet: %v", err)
			return
		}
		svc.TrackCosmosPending(pkt, height)
	}
	untrackCosmosPending := func(raw []byte) {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[adapter cosmos->eth] untrack pending: decode packet: %v", err)
			return
		}
		svc.UntrackCosmosPending(pkt)
	}

	// Pinned-set rotation cadence, derived from the on-chain trusting period (as
	// the legacy routine did). A derivation failure is FATAL — the legacy StartLoop
	// rejected an unsafe/underivable interval at startup rather than silently
	// falling back to a fixed default that could exceed the trusting period and let
	// the client expire. Fail loud so the operator fixes the config.
	periodicUpdateInterval, err := svc.PinnedSetRotationInterval(srcCtx)
	if err != nil {
		return fmt.Errorf("cosmos->eth: derive pinned-set rotation interval: %w", err)
	}
	// Initial delay from the client's ON-CHAIN freshness (time until next due), so
	// a restart near the rotation deadline fires promptly instead of waiting a full
	// fresh interval. On error, default to 0 (rotate now) — the safe direction:
	// never let the set decay by delaying the first rotation.
	initialRotationDelay, err := svc.PinnedSetRotationDueIn(srcCtx)
	if err != nil {
		log.Printf("[adapter cosmos->eth] derive initial rotation delay: %v; rotating on startup", err)
		initialRotationDelay = 0
	}

	cosmosToEth := relay.NewModule(
		"cosmos->eth",
		srcCtx.CosmosRouterClientID(),
		cosmos.NewSource(srcCtx, bb),
		evm.NewDestination(worker, srcCtx),
		cosmos.NewGroth16Builder(worker, srcCtx, cfg.ProofType, cfg.TrustLevel),
		relay.WithTimeoutScanner(0, func(context.Context) { svc.ScanCosmosTimeouts(srcCtx) }),
		relay.WithPacketTracker(trackCosmosPending, untrackCosmosPending),
		// Force-rotate the pinned validator set on a fixed cadence so it never
		// decays below quorum during a quiet period (the guaranteed rotation the
		// legacy StartLoop routine provided). ETH->Cosmos needs no equivalent —
		// the beacon client has no pinned set.
		relay.WithPeriodicUpdate(periodicUpdateInterval, initialRotationDelay, func(context.Context) error {
			return svc.RotatePinnedSet(srcCtx)
		}),
	)

	ethToCosmos := relay.NewModule(
		"eth->cosmos",
		srcCtx.EthClientID(),
		evm.NewSource(srcCtx, bb),
		cosmos.NewDestination(worker, srcCtx),
		evm.NewBeaconBuilder(worker, srcCtx),
		relay.WithTimeoutScanner(0, func(context.Context) { svc.ScanEthTimeouts(srcCtx) }),
	)

	// Child context so the first fatal error stops both modules.
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 2)
	go func() { errCh <- cosmosToEth.Run(runCtx) }()
	go func() { errCh <- ethToCosmos.Run(runCtx) }()

	err = <-errCh // first module to exit
	cancel()      // stop the other
	<-errCh       // wait for it so no goroutine leaks

	// A cancelled context is a clean shutdown, not a relay failure.
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}
