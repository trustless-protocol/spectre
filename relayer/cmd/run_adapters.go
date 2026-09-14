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

// runAdapterEngine drives one source's bidirectional relay through the generic
// relay.Module. Proof generation, gap recovery and timeout scanning stay in
// services and are reached through the chain adapters; the module owns only the
// orchestration.
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
// The first module to exit cancels the other, and BOTH results are joined into
// the return value. Returning only the first is what hid the failure this engine
// now surfaces: on SIGTERM the unaffected direction usually returns nil first,
// and returning that discarded the other direction's drain failure. Each module
// normalizes cancellation of the run context to nil while preserving real RPC
// deadlines and drain failures, so a clean stop joins two nils and stays nil.
func runAdapterEngine(ctx context.Context, svc *services.Services, deps services.RelayDeps) error {
	worker := svc.Worker()
	bb := svc.BatchBuilder
	cfg := svc.CosmosConfig()

	// Cosmos-origin sends aren't tracked by the subscriber; the module records
	// each one (proto bytes -> Packet) so ScanCosmosTimeouts can refund it, and
	// removes it once the send is received on ETH (it can no longer time out).
	trackCosmosPending := func(raw []byte, height uint64) bool {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[cosmos->eth] track pending: decode packet: %v", err)
			return false
		}
		return svc.TrackCosmosPending(pkt, height)
	}
	untrackCosmosPending := func(raw []byte) {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(raw); err != nil {
			log.Printf("[cosmos->eth] untrack pending: decode packet: %v", err)
			return
		}
		svc.UntrackCosmosPending(pkt)
	}

	// Pinned-set rotation cadence, derived from the on-chain trusting period. A
	// derivation failure is FATAL: an unsafe or underivable interval must be
	// rejected at startup rather than silently replaced by a fixed default that
	// could exceed the trusting period and let the client expire. Fail loud so the
	// operator fixes the config.
	periodicUpdateInterval, err := svc.PinnedSetRotationInterval(ctx, deps.EVM)
	if err != nil {
		if isShutdownErr(err) && ctx.Err() != nil {
			return nil
		}
		return fmt.Errorf("cosmos->eth: derive pinned-set rotation interval: %w", err)
	}
	// Initial delay from the client's ON-CHAIN freshness (time until next due), so
	// a restart near the rotation deadline fires promptly instead of waiting a full
	// fresh interval. On error, default to 0 (rotate now) — the safe direction:
	// never let the set decay by delaying the first rotation.
	initialRotationDelay, err := svc.PinnedSetRotationDueIn(ctx, deps.Cosmos, deps.EVM)
	if err != nil {
		if isShutdownErr(err) && ctx.Err() != nil {
			return nil
		}
		log.Printf("[cosmos->eth] derive initial rotation delay: %v; rotating on startup", err)
		initialRotationDelay = 0
	}
	if err := svc.SeedEVMOnCosmosUpdate(ctx, deps.Cosmos, deps.IDs.EVMOnCosmos); err != nil {
		log.Printf("[eth->cosmos] seed client-update age: %v", err)
	}

	cosmosToEth := relay.NewModule(
		"cosmos->eth",
		deps.IDs.CosmosOnEVM,
		cosmos.NewSource(deps.Cosmos, deps.EVM, deps.IDs, deps.Config.FetchTimeout, deps.Config.BatchConfig, deps.Logger, bb, svc.RecoveryState()),
		evm.NewDestination(worker, deps.Cosmos, deps.EVM, deps.IDs.CosmosOnEVM),
		cosmos.NewGroth16Builder(worker, deps.Cosmos, deps.EVM, deps.Config.FetchTimeout, deps.Config.RotationThreshold, cfg.ProofType, cfg.TrustLevel),
		relay.WithClientUpdateObserver(svc.ObserveCosmosOnEVMUpdate),
		relay.WithTimeoutScanner(0, func(c context.Context) {
			svc.ScanCosmosTimeouts(c, deps.Cosmos, deps.EVM, deps.IDs.EVMOnCosmos)
		}),
		relay.WithPacketTracker(trackCosmosPending, untrackCosmosPending),
		// Force-rotate the pinned validator set on a fixed cadence so it never
		// decays below quorum during a quiet period. ETH->Cosmos needs no
		// equivalent — the beacon client has no pinned set.
		relay.WithPeriodicUpdate(periodicUpdateInterval, initialRotationDelay, func(c context.Context) error {
			return svc.RotatePinnedSet(c, deps.Cosmos, deps.EVM, deps.IDs.CosmosOnEVM)
		}),
	)

	ethToCosmos := relay.NewModule(
		"eth->cosmos",
		deps.IDs.EVMOnCosmos,
		evm.NewSource(deps.Cosmos, deps.EVM, deps.IDs, deps.Config.BatchConfig, deps.Logger, bb, svc.RecoveryState()),
		cosmos.NewDestination(worker, deps.Cosmos, deps.IDs.EVMOnCosmos),
		evm.NewBeaconBuilder(worker, deps.Cosmos, deps.EVM, deps.IDs.EVMOnCosmos),
		relay.WithClientUpdateObserver(svc.ObserveEVMOnCosmosUpdate),
		relay.WithTimeoutScanner(0, func(c context.Context) { svc.ScanEthTimeouts(c, deps.Cosmos, deps.EVM, deps.IDs.CosmosOnEVM) }),
	)

	// Child context so the first fatal error stops both modules.
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go svc.RunQueueReporter(runCtx)
	errCh := make(chan error, 2)
	go func() { errCh <- cosmosToEth.Run(runCtx) }()
	go func() { errCh <- ethToCosmos.Run(runCtx) }()

	first := <-errCh // first module to exit
	cancel()         // stop the other
	second := <-errCh

	// BOTH results, not just the first -- see the contract on this function.
	// Reading the second only to avoid a goroutine leak is what discarded the
	// other direction's drain failure before it could reach loopErrCh.
	return errors.Join(first, second)
}
