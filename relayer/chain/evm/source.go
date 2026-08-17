package evm

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"relayer/chain"
	relayerclient "relayer/client"
	"relayer/services"
	"relayer/subscriber"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Source is the Ethereum L1 implementation of chain.Source for the ETH->Cosmos
// direction. It holds the endpoints and IDs needed by beacon proofs and state reads,
// and the shared *services.BatchBuilder so Subscribe drains the same queue the
// Services instance owns (where the ETH subscriber records EthPendingTracker
// entries) and can re-queue handler failures (with a waiting backoff) — the
// re-queue behavior of the legacy handleEth, expressed generically.
// Constructed by the wiring, not the cfg-only registry.
type Source struct {
	cosmos      services.CosmosEndpoint
	evm         services.EVMEndpoint
	ids         services.ClientIDs
	batchConfig services.BatchConfig
	logger      *log.Logger
	bb          *services.BatchBuilder
	recovery    *services.RecoveryStateStore
}

// NewSource wires the Ethereum source to its endpoints and the shared batch builder.
func NewSource(cosmos services.CosmosEndpoint, evm services.EVMEndpoint, ids services.ClientIDs, batchConfig services.BatchConfig, logger *log.Logger, bb *services.BatchBuilder, recovery ...*services.RecoveryStateStore) *Source {
	var store *services.RecoveryStateStore
	if len(recovery) > 0 {
		store = recovery[0]
	}
	return &Source{cosmos: cosmos, evm: evm, ids: ids, batchConfig: batchConfig, logger: logger, bb: bb, recovery: store}
}

func (s *Source) Chain() chain.ChainType { return chain.Ethereum }

// LatestHeight returns the execution block number of the latest FINALIZED beacon
// header. It MUST be the same number space as RelayableHeight, event heights, and
// ClientUpdate.Height (all execution block numbers, see BeaconBuilder.Build) —
// the module's refreshLoop feeds this into updateClientTo, which compares it
// against m.lastHeight (an execution block number). Returning a beacon slot here
// would mix number spaces and make the anti-expiry refresh skip incorrectly. The
// beacon light client only advances to finalized headers, so there is no
// unsafe-head exposure — latest and relayable coincide on the ETH source side.
func (s *Source) LatestHeight(ctx context.Context) (uint64, error) {
	return s.finalizedExecBlock(ctx)
}

// RelayableHeight returns the execution block number of the latest FINALIZED
// beacon header — the highest ETH block the beacon light client can prove against
// right now. The module compares packet event heights (execution block numbers)
// against this, so a packet whose block is not yet finalized is re-queued without
// building a proof (the legacy waitBeaconFinality precondition, as a value).
func (s *Source) RelayableHeight(ctx context.Context) (uint64, error) {
	return s.finalizedExecBlock(ctx)
}

// finalizedExecBlock reads the beacon finality update and returns the execution
// block number of the finalized header — the single finality-gated height both
// LatestHeight and RelayableHeight report.
func (s *Source) finalizedExecBlock(_ context.Context) (uint64, error) {
	beaconURL := s.evm.BeaconAPIURL
	if beaconURL == "" {
		return 0, fmt.Errorf("eth source: beacon API URL is not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	fu, err := relayerclient.GetFinalityUpdate(ctx, beaconURL)
	if err != nil {
		return 0, fmt.Errorf("eth source: finality update: %w", err)
	}
	block, err := strconv.ParseUint(fu.FinalizedHeader.Execution.BlockNumber, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("eth source: parse finalized exec block %q: %w", fu.FinalizedHeader.Execution.BlockNumber, err)
	}
	return block, nil
}

// QueryHeader is a no-op for the beacon path: the BeaconBuilder self-fetches the
// finality/sync-committee data via the beacon API (like the Cosmos groth16
// builder), so it ignores the header argument.
func (s *Source) QueryHeader(_ context.Context, _ uint64) ([]byte, error) {
	return nil, nil
}

// ethDrainInterval is how often the Subscribe bridge flushes the batch builder.
const ethDrainInterval = 500 * time.Millisecond

// Subscribe bridges the existing SubscribeEth gap-recovery subscriber to the
// handler-based contract (mirror of the Cosmos source bridge): it runs
// SubscribeEth pushing into a BatchBuilder, drains it, and emits each EthSend as
// a recv event and each EthWriteAck as an ack event. EthTimeout is skipped —
// ETH-origin timeouts are handled by the async scanner, not this path.
func (s *Source) Subscribe(ctx context.Context, handler func(context.Context, []chain.Event) []int) error {
	sub := subscriber.NewSubscriber(s.recovery)
	go sub.SubscribeEth(s.cosmos, s.evm, s.ids, s.logger, s.bb)

	// Use the configured batch window so CheckEth returns multi-packet batches the
	// handler can fold into one multicall (BatchSize=1 would defeat that).
	cfg := s.batchConfig
	ch := make(chan services.EthBatch, services.BatchHandoffCapacity)

	go func() {
		ticker := time.NewTicker(ethDrainInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.bb.CheckEth(ctx, cfg, ch)
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
			orig := make([]services.EthPacket, 0, len(batch.Packets))
			for _, p := range batch.Packets {
				// Terminal events (the packet was acked or timed out on ETH) settle a
				// pending ETH-origin send: remove it from the tracker so the timeout
				// scanner stops considering it (mirrors the legacy handleEth EthAck/
				// EthTimeout tracker removal), then skip — nothing to relay.
				if p.Packet != nil && (p.Type == services.EthAck || p.Type == services.EthTimeout) {
					s.bb.EthPendingTracker.Remove(p.Packet.SourceClient, p.Packet.Sequence)
					continue
				}
				e, ok := ethPacketToEvent(p)
				if !ok {
					continue
				}
				events = append(events, e)
				orig = append(orig, p)
			}
			// Re-queue un-relayed packets with a waiting backoff so a packet not yet
			// relayable (beacon finality lag) or hit by a brief RPC hiccup is retried
			// with a growing delay instead of every batch period — quiet, and no
			// per-flush finality RPC. CheckEth already sliced them off the queue.
			for _, idx := range handler(ctx, events) {
				s.bb.RequeueEthWaiting([]services.EthPacket{orig[idx]})
			}
			s.bb.ReleaseEthInFlight(batch)
		}
	}
}

// ethPacketToEvent maps a queued EthPacket to a chain.Event, or false to skip
// (unencodable, or an EthTimeout the scanner owns).
func ethPacketToEvent(p services.EthPacket) (chain.Event, bool) {
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
	case services.EthSend:
		e.Type = chain.SendPacket
	case services.EthWriteAck:
		// Skip an ack with no acknowledgement bytes — it cannot be built into a
		// MsgAcknowledgement (mirrors the legacy handleEth EthWriteAck no-ack skip).
		if len(p.AckBytes) == 0 {
			return chain.Event{}, false
		}
		e.Type = chain.AckPacket
		e.AckBytes = p.AckBytes
	default:
		return chain.Event{}, false // EthTimeout is scanner-handled
	}
	return e, true
}

// MembershipProof builds an Ethereum storage proof that the packet's commitment
// (recv) or ack exists, wrapping client.GetEthMembershipProof. The proof is taken
// at the execution block of the on-chain 08-wasm client's latest slot (read from
// Cosmos), so it verifies against the client state the destination trusts. The
// height argument is not the proof block (the proof block is read from the
// on-chain client above), so it is ignored.
func (s *Source) MembershipProof(ctx context.Context, packet []byte, _ uint64, eventType chain.EventType) ([]byte, error) {
	pkt, ethClientState, err := s.decodePacketAndClientState(packet)
	if err != nil {
		return nil, err
	}
	var clientID string
	var pathType byte
	switch eventType {
	case chain.SendPacket:
		// An ETH->Cosmos send past its timeout can never be received on Cosmos —
		// MsgRecvPacket would fail. Report it permanent so the module DROPS it and
		// the async scanner (scanForEthTimeouts) refunds it on ETH instead (mirrors
		// the legacy ethPacketExpired pre-filter). Cosmos has ~wall-clock BFT time,
		// so compare against time.Now like the legacy path does. Without this a dead
		// send is retried forever, draining gas on repeated client updates.
		if pkt.TimeoutTimestamp > 0 && uint64(time.Now().Unix()) >= pkt.TimeoutTimestamp {
			return nil, chain.Permanent(fmt.Errorf("eth source: send seq=%d timed out; deferred to timeout scanner", pkt.Sequence))
		}
		clientID, pathType = pkt.SourceClient, 1 // packet commitment
	case chain.AckPacket:
		clientID, pathType = pkt.DestinationClient, 3 // ack
	default:
		return nil, fmt.Errorf("eth source: MembershipProof: unsupported event type %d", eventType)
	}
	path := services.EthPath(clientID, pkt.Sequence, pathType)
	return relayerclient.GetEthMembershipProof(
		ctx, s.evm.EthClient(), s.evm.Contracts.Router, path,
		ethcommon.HexToHash(services.ICS26_IBC_STORAGE_SLOT),
		new(big.Int).SetUint64(ethClientState.LatestExecutionBlockNumber),
	)
}

// NonMembershipProof is unused on the ETH source relay path: ETH-origin packet
// timeouts are handled by the async timeout scanner (scanForEthTimeouts), not the
// Subscribe->relay flow, so no TimeoutPacket event reaches here.
func (s *Source) NonMembershipProof(_ context.Context, _ []byte, _ uint64) ([]byte, error) {
	return nil, fmt.Errorf("eth source: NonMembershipProof unused (timeouts are scanner-handled)")
}

// decodePacketAndClientState decodes the proto packet and reads the on-chain
// 08-wasm ETH client state (for the proof's execution block).
func (s *Source) decodePacketAndClientState(packet []byte) (channeltypesv2.Packet, *relayerclient.EthereumClientState, error) {
	var pkt channeltypesv2.Packet
	if err := pkt.Unmarshal(packet); err != nil {
		return channeltypesv2.Packet{}, nil, fmt.Errorf("eth source: decode packet: %w", err)
	}
	cs, err := relayerclient.GetEthereumClientState(s.cosmos.CosmosClient(), s.ids.EVMOnCosmos)
	if err != nil {
		return channeltypesv2.Packet{}, nil, fmt.Errorf("eth source: eth client state: %w", err)
	}
	return pkt, cs, nil
}
