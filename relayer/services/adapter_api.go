package services

import (
	"context"
	"fmt"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"

	client "relayer/client"
)

// This file exposes the timeout-recovery surface (previously driven internally by
// the now-removed StartLoop) so the chain-adapter relay path (package relay) can
// reuse it verbatim instead of reimplementing the (subtle) scanners and
// pending-tracker bookkeeping. These are thin wrappers over the unexported
// originals — no logic is duplicated.

// ScanCosmosTimeouts detects Cosmos-origin packets that expired on ETH without a
// receipt and relays their MsgTimeout back to Cosmos. Safe to call periodically;
// it purges/prunes the pending tracker and recovers from panics internally. stdCtx
// aborts the in-scan Cosmos catch-up wait promptly on shutdown.
func (s *Services) ScanCosmosTimeouts(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint, ethClientID string) {
	s.scanForCosmosTimeouts(stdCtx, cosmos, evm, ethClientID)
}

// ScanEthTimeouts detects ETH-origin packets that expired on Cosmos without a
// receipt and relays their MsgTimeout back to ETH. Same periodic-call contract as
// ScanCosmosTimeouts.
func (s *Services) ScanEthTimeouts(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint, routerClientID string) {
	s.scanForEthTimeouts(stdCtx, evmTimeoutDeps{cosmos: cosmos, evm: evm, routerClientID: routerClientID})
}

// ScanL2Timeouts detects L2-origin packets that expired on Cosmos without a
// receipt and relays their timeoutPacket back to the L2 ICS26Router. It uses the
// same Cosmos non-membership proof path as ScanEthTimeouts, but drains the L2
// pending tracker so it never interferes with the legacy ETH source tracker.
func (s *Services) ScanL2Timeouts(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint, routerClientID string) {
	s.scanForL2Timeouts(stdCtx, evmTimeoutDeps{cosmos: cosmos, evm: evm, routerClientID: routerClientID})
}

// TrackCosmosPending records a Cosmos-origin packet just recv-relayed to ETH so
// ScanCosmosTimeouts can later refund it if it expires undelivered. Mirrors the
// PendingTracker.Add that handleCosmos performs in the StartLoop path. False
// means the durable state write failed and the source event must be retried.
func (s *Services) TrackCosmosPending(packet channeltypesv2.Packet, blockNumber uint64) bool {
	return s.BatchBuilder.PendingTracker.Add(packet, blockNumber)
}

// UntrackCosmosPending removes a Cosmos-origin packet from the pending tracker
// once it has been successfully received on ETH — it can no longer time out, so
// the timeout scanner need not keep querying its receipt. Mirrors the
// PendingTracker.Remove-on-recv that handleCosmos performs in the StartLoop path.
func (s *Services) UntrackCosmosPending(packet channeltypesv2.Packet) {
	s.BatchBuilder.PendingTracker.RemovePacketIfCurrent(packet)
}

// TrackL2Pending records an L2-origin packet observed on the L2 source path so
// ScanL2Timeouts can refund it if Cosmos never receives it before timeout. False
// means the durable state write failed and the source event must be retried.
func (s *Services) TrackL2Pending(packet channeltypesv2.Packet, blockNumber uint64) bool {
	return s.BatchBuilder.L2PendingTracker.Add(packet, blockNumber)
}

// UntrackL2Pending removes an L2-origin packet once the L2->Cosmos receive relay
// succeeds; after a Cosmos receipt exists, a timeout refund must not be attempted.
func (s *Services) UntrackL2Pending(packet channeltypesv2.Packet) {
	s.BatchBuilder.L2PendingTracker.RemovePacketIfCurrent(packet)
}

// Worker exposes the shared Worker (TxHandler + Prover) so the chain adapters —
// which wrap the existing Build/Send pipeline — can be constructed with the same
// worker the Services instance uses. Read-only handle; the adapters do not mutate
// it.
func (s *Services) Worker() *Worker { return s.worker }

// CosmosConfig exposes the resolved Cosmos-source config (ProofType, TrustLevel,
// BatchConfig, …) so the adapter wiring can build the groth16 client-update
// builder with the same parameters StartLoop uses.
func (s *Services) CosmosConfig() Config { return s.cosmosConfig }

// RotatePinnedSet force-rotates the on-chain SpectreClient's pinned
// validator set (RefreshCosmosClient with forceRotation=true), regardless of how
// much overlap remains. This is the guaranteed rotation the legacy
// StartLoop routine ran on a fixed cadence: the per-packet and expiry-driven
// refresh paths only rotate when overlap has already decayed to the threshold, so
// during a quiet period (no packet traffic) the pinned set could otherwise decay
// below the 2/3 quorum needed to update OR rotate the client — bricking it. This
// call keeps the pinned set fresh independent of packet flow. Side-effect-free on
// the adapter's cursors (the expiry is driven by ClientExpiresAt, not a seeded
// timestamp), so it only submits the rotation tx.
func (s *Services) RotatePinnedSet(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint, routerClientID string) error {
	lightBlock, err := s.worker.RefreshCosmosClient(stdCtx, cosmos, evm, routerClientID, s.cosmosConfig.FetchTimeout, s.cosmosConfig.RotationThreshold, s.cosmosConfig.ProofType, s.cosmosConfig.TrustLevel)
	if err == nil && lightBlock != nil {
		s.ObserveCosmosOnEVMUpdate(lightBlock.SignedHeader.Header.Time)
	}
	return err
}

// PinnedSetRotationInterval derives the force-rotation cadence from the on-chain
// trusting period exactly as the legacy routine did (min of the configured
// refresh interval and trustingPeriod minus a safety margin), so the pinned set
// is refreshed well within the window where it stays above quorum.
func (s *Services) PinnedSetRotationInterval(stdCtx context.Context, evm EVMEndpoint) (time.Duration, error) {
	readCtx, cancel := fetchCtx(stdCtx, s.cosmosConfig.FetchTimeout)
	defer cancel()
	clientState, err := fetchOnChainClientStateWithContext(readCtx, evm)
	if err != nil {
		return 0, err
	}
	trustingPeriod := time.Duration(clientState.TrustingPeriod) * time.Second
	return deriveCosmosRefreshInterval(s.cosmosConfig, trustingPeriod)
}

// PinnedSetRotationDueIn reports how long until the next pinned-set rotation is
// due, based on the ON-CHAIN trusted timestamp (not process uptime): it is the
// remaining time until trustedTime + interval, clamped to >= 0. A client that is
// already at or past its rotation interval at startup returns 0 (rotate now),
// mirroring the legacy seedCosmosClientFreshness behavior — so a relayer restart
// near the rotation deadline does not wait a full fresh interval before the first
// rotation, which could let the pinned set decay below quorum.
func (s *Services) PinnedSetRotationDueIn(stdCtx context.Context, cosmos CosmosEndpoint, evm EVMEndpoint) (time.Duration, error) {
	interval, err := s.PinnedSetRotationInterval(stdCtx, evm)
	if err != nil {
		return 0, err
	}
	readCtx, cancel := fetchCtx(stdCtx, s.cosmosConfig.FetchTimeout)
	trustedHeight, err := FetchOnChainTrustedHeightWithContext(readCtx, evm)
	cancel()
	if err != nil {
		return 0, err
	}
	lightCtx, cancelLight := fetchCtx(stdCtx, s.cosmosConfig.FetchTimeout)
	defer cancelLight()
	lightBlock, err := client.GetLightBlockWithContext(lightCtx, cosmos.CosmosClient(), trustedHeight)
	if err != nil {
		return 0, err
	}
	trustedTime := lightBlock.SignedHeader.Header.Time
	if trustedTime.IsZero() {
		return 0, fmt.Errorf("zero trusted light block timestamp at height %d", trustedHeight)
	}
	s.ObserveCosmosOnEVMUpdate(trustedTime)
	due := time.Until(trustedTime.Add(interval))
	if due < 0 {
		return 0, nil
	}
	return due, nil
}

// SeedEVMOnCosmosUpdate records the timestamp already trusted by the on-chain
// beacon client so the first queue report is useful even before packet traffic
// or a refresh advances it in this process.
func (s *Services) SeedEVMOnCosmosUpdate(cosmos CosmosEndpoint, ethClientID string) error {
	state, err := client.GetEthereumClientState(cosmos.CosmosClient(), ethClientID)
	if err != nil {
		return err
	}
	trustedTimestamp := state.ComputeTimestampAtSlot(state.LatestSlot)
	if trustedTimestamp == 0 {
		return fmt.Errorf("zero trusted EVM timestamp at slot %d", state.LatestSlot)
	}
	s.ObserveEVMOnCosmosUpdate(time.Unix(int64(trustedTimestamp), 0))
	return nil
}
