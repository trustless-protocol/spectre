// Package cosmos adapts a Cosmos (Tendermint) chain to the chain.Source and
// chain.ClientUpdateBuilder contracts. It wraps the existing relayer/client and
// relayer/services code — it does not reimplement signature extraction, proof
// generation, or gap-recovery.
package cosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"relayer/chain"
	relayerclient "relayer/client"
	"relayer/services"
	"relayer/subscriber"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// ICS-24 path type discriminators (mirrors the constants used inline by the
// legacy packet planner): 1 = packet commitment, 2 = receipt, 3 = ack.
var (
	pathCommitment = []byte{1}
	pathReceipt    = []byte{2}
	pathAck        = []byte{3}
)

// Source is the Cosmos implementation of chain.Source. It holds a
// services.Context because the packet-proof methods wrap services.CosmosMembership
// (which reads the Cosmos RPC via the context); it is constructed by the wiring,
// not the cfg-only registry.
//
// It also holds the shared *services.BatchBuilder so Subscribe drains the same
// queue the (single) Services instance owns and can re-queue handler failures
// back onto it (with a waiting backoff) — the re-queue behavior of the legacy
// handleCosmos, expressed generically.
type Source struct {
	svcCtx services.Context
	bb     *services.BatchBuilder
}

// NewSource wires the Cosmos source to the shared context + batch builder.
func NewSource(svcCtx services.Context, bb *services.BatchBuilder) *Source {
	return &Source{svcCtx: svcCtx, bb: bb}
}

func (s *Source) Chain() chain.ChainType { return chain.Cosmos }

// LatestHeight returns the latest committed Tendermint height. Tendermint has
// BFT instant finality, so the latest committed block is already final.
func (s *Source) LatestHeight(_ context.Context) (uint64, error) {
	lb, err := relayerclient.GetLatestLightBlock(s.svcCtx.CosmosClient())
	if err != nil {
		return 0, fmt.Errorf("cosmos source: latest light block: %w", err)
	}
	if lb.BlockHeight < 0 {
		return 0, fmt.Errorf("cosmos source: negative block height %d", lb.BlockHeight)
	}
	return uint64(lb.BlockHeight), nil
}

// cosmosAppHashLag is how far the queryable app state trails the committed
// height: a packet commitment written at block H is only reflected in the AppHash
// (and thus provable) at H+2. RelayableHeight subtracts this.
const cosmosAppHashLag = 2

// RelayableHeight returns latest committed height minus the AppHash lag, so the
// module only proves a Cosmos packet once its commitment is queryable (the legacy
// waitCosmosAppHash(maxPacketHeight+2) precondition, as a value).
func (s *Source) RelayableHeight(ctx context.Context) (uint64, error) {
	latest, err := s.LatestHeight(ctx)
	if err != nil {
		return 0, err
	}
	if latest < cosmosAppHashLag {
		return 0, nil
	}
	return latest - cosmosAppHashLag, nil
}

// QueryHeader returns the JSON-encoded light block at height.
func (s *Source) QueryHeader(_ context.Context, height uint64) ([]byte, error) {
	lb, err := relayerclient.GetLightBlock(s.svcCtx.CosmosClient(), int64(height))
	if err != nil {
		return nil, fmt.Errorf("cosmos source: light block at %d: %w", height, err)
	}
	raw, err := json.Marshal(lb)
	if err != nil {
		return nil, fmt.Errorf("cosmos source: encode light block: %w", err)
	}
	return raw, nil
}

// MembershipProof proves the packet's commitment (recv) or ack (ack) exists at
// height, wrapping services.CosmosMembership. eventType selects the ICS-24 path.
func (s *Source) MembershipProof(_ context.Context, packet []byte, height uint64, eventType chain.EventType) ([]byte, error) {
	pkt, lb, err := s.decodePacketAndBlock(packet, height)
	if err != nil {
		return nil, err
	}
	var clientID string
	var pathType []byte
	switch eventType {
	case chain.SendPacket:
		// A Cosmos->ETH send whose timeout has passed (measured against ETH block
		// time, the destination clock) can never be received on ETH — recvPacket
		// would revert. Report it permanent so the module DROPS it and the async
		// timeout scanner refunds it on Cosmos instead (mirrors the legacy
		// planCosmosPacketMsgs ethBlockTime>=TimeoutTimestamp skip). Without this the
		// module retries the dead packet forever, each retry re-running a full
		// client update and draining gas.
		if timedOut, err := s.sendTimedOutOnEth(pkt); err != nil {
			return nil, err // transient (RPC) — let the module retry
		} else if timedOut {
			return nil, chain.Permanent(fmt.Errorf("cosmos source: send seq=%d timed out; deferred to timeout scanner", pkt.Sequence))
		}
		clientID, pathType = pkt.SourceClient, pathCommitment
	case chain.AckPacket:
		clientID, pathType = pkt.DestinationClient, pathAck
	default:
		return nil, fmt.Errorf("cosmos source: MembershipProof: unsupported event type %d", eventType)
	}
	return services.CosmosMembership(s.svcCtx, pkt, clientID, pathType, lb)
}

// sendTimedOutOnEth reports whether the packet's timeout timestamp has passed on
// the ETH side, comparing against the latest ETH block time (the destination
// clock the on-chain timeout check uses). A zero timeout never expires.
func (s *Source) sendTimedOutOnEth(pkt channeltypesv2.Packet) (bool, error) {
	if pkt.TimeoutTimestamp == 0 {
		return false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.svcCtx.Config.FetchTimeout)
	defer cancel()
	header, err := s.svcCtx.EthClient().HeaderByNumber(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("cosmos source: eth block time: %w", err)
	}
	return header.Time >= pkt.TimeoutTimestamp, nil
}

// NonMembershipProof proves the packet receipt is absent at height (timeout),
// wrapping services.CosmosNonMembership.
func (s *Source) NonMembershipProof(_ context.Context, packet []byte, height uint64) ([]byte, error) {
	pkt, lb, err := s.decodePacketAndBlock(packet, height)
	if err != nil {
		return nil, err
	}
	return services.CosmosNonMembership(s.svcCtx, pkt, pkt.DestinationClient, pathReceipt, lb)
}

// decodePacketAndBlock decodes the proto packet and fetches the light block the
// proof is built against.
func (s *Source) decodePacketAndBlock(packet []byte, height uint64) (channeltypesv2.Packet, *relayerclient.LightBlock, error) {
	var pkt channeltypesv2.Packet
	if err := pkt.Unmarshal(packet); err != nil {
		return channeltypesv2.Packet{}, nil, fmt.Errorf("cosmos source: decode packet: %w", err)
	}
	lb, err := relayerclient.GetLightBlock(s.svcCtx.CosmosClient(), int64(height))
	if err != nil {
		return channeltypesv2.Packet{}, nil, fmt.Errorf("cosmos source: light block at %d: %w", height, err)
	}
	return pkt, lb, nil
}

// drainInterval is how often the Subscribe bridge flushes the batch builder.
const drainInterval = 500 * time.Millisecond

// Subscribe bridges the existing gap-recovery subscriber to the handler-based
// contract: it runs services' SubscribeCosmos (which keeps all the startup
// lookback / reconnect-gap recovery) pushing into a BatchBuilder, then drains
// the builder on the configured batch window and emits each queued CosmosPacket
// as a chain.Event. It uses the batch config (not BatchSize=1) so CheckCosmos
// returns multi-packet batches the module folds into one RelayPackets multicall.
func (s *Source) Subscribe(ctx context.Context, handler func(context.Context, []chain.Event) []int) error {
	sub := subscriber.NewSubscriber()
	go sub.SubscribeCosmos(s.svcCtx, s.bb)

	// Use the configured batch window so CheckCosmos returns multi-packet batches
	// the handler can fold into one multicall (BatchSize=1 would defeat that).
	cfg := s.svcCtx.Config.BatchConfig
	ch := make(chan services.CosmosBatch, 16)

	go func() {
		ticker := time.NewTicker(drainInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.bb.CheckCosmos(cfg, ch)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case batch := <-ch:
			// Build the event slice aligned 1:1 with orig[] so a re-queue index
			// from the handler maps back to the exact source packet.
			events := make([]chain.Event, 0, len(batch.Packets))
			orig := make([]services.CosmosPacket, 0, len(batch.Packets))
			for _, p := range batch.Packets {
				e, ok := cosmosPacketToEvent(p, s.svcCtx.CosmosRouterClientID())
				if !ok {
					continue
				}
				events = append(events, e)
				orig = append(orig, p)
			}
			// Re-queue un-relayed packets with a waiting backoff so a packet not yet
			// relayable (AppHash H+2 lag) or hit by a brief RPC hiccup is retried with
			// a growing delay instead of every batch period — quiet, and no per-flush
			// proof/RPC churn. CheckCosmos already sliced them off the queue.
			for _, idx := range handler(ctx, events) {
				s.bb.RequeueCosmosWaiting([]services.CosmosPacket{orig[idx]})
			}
		}
	}
}

// cosmosPacketToEvent maps a queued CosmosPacket to a chain.Event. It returns
// false when the packet cannot be encoded (skip rather than relay a bad event) or
// when the event must not be relayed to ETH.
//
// routerClientID is this Cosmos chain's ETH-side router client id, used to drop
// Cosmos-originated timeouts: an IBC v2 timeout is submitted to the SOURCE chain,
// so a Cosmos→ETH packet's timeout fires on Cosmos and ETH never held a
// commitment for it — relaying an ETH timeoutPacket would waste a Groth16 proof +
// tx (and historically blocked the relay loop). The live subscribe path enqueues
// these (destination_client == routerClientID matches the configured-client
// filter), so the check the legacy planCosmosPacketMsgs applied at plan time is
// applied here instead. ETH-originated timeouts (destination_client != routerClientID)
// still relay so ETH deletes its commitment + refunds.
func cosmosPacketToEvent(p services.CosmosPacket, routerClientID string) (chain.Event, bool) {
	if p.Packet == nil {
		return chain.Event{}, false
	}
	raw, err := p.Packet.Marshal()
	if err != nil {
		return chain.Event{}, false
	}
	e := chain.Event{
		Height:   p.BlockNumber,
		Sequence: p.Packet.Sequence,
		ClientID: p.Packet.DestinationClient,
		Raw:      raw,
	}
	switch p.Type {
	case services.CosmosSend:
		e.Type = chain.SendPacket
	case services.CosmosAck:
		// An ack with no acknowledgement bytes cannot be built into a MsgAckPacket
		// (the ETH side needs the ack) — skip it like the legacy planCosmosPacketMsgs
		// did, rather than emitting an event that would error and re-queue forever.
		if len(p.AckBytes) == 0 {
			return chain.Event{}, false
		}
		e.Type = chain.AckPacket
		e.AckBytes = p.AckBytes
	case services.CosmosTimeout:
		if !services.ShouldRelayCosmosTimeoutToEth(p.Packet, routerClientID) {
			return chain.Event{}, false // Cosmos-originated timeout handled locally on Cosmos
		}
		e.Type = chain.TimeoutPacket
	default:
		return chain.Event{}, false
	}
	return e, true
}
