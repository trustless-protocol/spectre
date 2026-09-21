package evm

import (
	"context"
	"fmt"
	"time"

	"relayer/chain"
	"relayer/chain/evmflush"
	relayerclient "relayer/client"
	"relayer/services"
	"relayer/subscriber"
)

// flushState is the enumeration backstop's per-source state: one rotating
// cursor per leg and the log-range policy the lookups narrow in place.
type flushState struct {
	// sendCursor is the lowest sequence the last send probe reached; 0 means
	// "start from the newest". See evmflush.Window.
	sendCursor uint64
	// ackCursor is the last Cosmos commitment sequence the ack leg took, so the
	// next pass continues after it. See services.CosmosCommitmentWindow.
	ackCursor uint64
	span      relayerclient.LogSpan
}

// flushQueryTimeout bounds one enumeration pass, for the same reason the Cosmos
// source bounds its own: a hung endpoint must not wedge the flush worker.
const flushQueryTimeout = 60 * time.Second

// UnrelayedPackets enumerates what this ETH source still owes the Cosmos side,
// without being told which packets to look for. It implements
// chain.PacketLister; see that interface for why enumeration exists alongside
// block scanning.
//
// Two legs, both bounded, both answering "is it still open on the SENDING chain
// and not yet settled on this one":
//
//   - sends: the router still holds the packet's commitment. Probed below
//     prevSequenceSend in a rotating window (evmflush), never 1..N.
//   - acknowledgements: Cosmos still holds a send's commitment AND this router
//     has written an acknowledgement for it. The Cosmos commitment store is the
//     enumerable side, so it is paged (services.CosmosCommitmentWindow) and the
//     router is probed only for the sequences it lists.
//
// A leg whose query fails is reported and skipped rather than failing the pass:
// the other leg's packets are exactly what the pass exists to recover.
func (s *Source) UnrelayedPackets(ctx context.Context) ([]chain.Event, error) {
	if s.evm.EthClient() == nil || s.evm.RouterContract() == nil {
		return nil, fmt.Errorf("eth source: enumeration needs an ETH client and a router address")
	}
	queryCtx, cancel := context.WithTimeout(ctx, flushQueryTimeout)
	defer cancel()

	router := evmflush.Router{
		Client: s.evm.EthClient(), Address: *s.evm.RouterContract(),
		Span: &s.flush.span, Label: "[eth->cosmos Flush]", Logger: s.logger,
	}
	sends, err := s.unrelayedSends(queryCtx, router)
	if err != nil {
		return nil, err
	}
	acks, err := s.unrelayedAcks(queryCtx, router)
	if err != nil {
		s.logger.Printf("[eth->cosmos Flush] acknowledgement leg: %v; the send leg's %d packet(s) still stand", err, len(sends))
	}
	return append(sends, acks...), nil
}

// unrelayedSends probes the router's commitment slots for this source's client
// and rebuilds a SendPacket event for each one still held.
func (s *Source) unrelayedSends(ctx context.Context, router evmflush.Router) ([]chain.Event, error) {
	clientID := s.ids.CosmosOnEVM
	prev, err := router.PrevSequenceSend(ctx, clientID)
	if err != nil {
		return nil, err
	}
	lo, hi, ok := evmflush.Window(prev, s.flush.sendCursor, evmflush.ProbeCap)
	if !ok {
		return nil, nil
	}
	slots, err := router.OutstandingSlots(ctx, evmflush.Sequences(lo, hi), func(seq uint64) []byte {
		return services.EthPath(clientID, seq, services.EthPathCommitment)
	})
	if err != nil {
		return nil, err
	}
	// The cursor moves only once the probe has answered. A failed probe leaves
	// the window to the next pass; advancing first would skip these sequences
	// for a whole rotation on one RPC blip.
	s.flush.sendCursor = lo

	var events []chain.Event
	for _, slot := range slots {
		// The slot's value says WHICH packet this sequence names; SendLog
		// resolves it past the older generations a re-created client leaves in
		// the log index.
		ev, found, err := router.SendLog(ctx, clientID, slot.Sequence, slot.Value)
		if err != nil {
			// One unresolvable sequence must not discard the rest.
			s.logger.Printf("[eth->cosmos Flush] resolve send seq=%d: %v", slot.Sequence, err)
			continue
		}
		if !found {
			continue // SendLog has already said why
		}
		pkt := subscriber.EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
		// Past its timeout it can never be received; the timeout scanner refunds
		// it. Same rule the live path applies.
		if sendPacketExpired(pkt, time.Now()) {
			continue
		}
		if e, ok := ethPacketToEvent(services.EthPacket{Type: services.EthSend, Packet: &pkt, BlockNumber: ev.Raw.BlockNumber}); ok {
			events = append(events, e)
		}
	}
	return events, nil
}

// unrelayedAcks lists the acknowledgements this router has written for Cosmos
// sends whose commitment Cosmos still holds -- acks the Cosmos side has not
// received back.
func (s *Source) unrelayedAcks(ctx context.Context, router evmflush.Router) ([]chain.Event, error) {
	if s.cosmos.CosmosClient() == nil {
		return nil, nil
	}
	// Cosmos->ETH packets are committed on Cosmos under ITS client for ETH, and
	// acknowledged here under THIS router's client for Cosmos.
	window, err := services.CosmosCommitmentWindow(ctx, s.cosmos, s.ids.EVMOnCosmos, s.flush.ackCursor, evmflush.ProbeCap)
	if err != nil {
		return nil, fmt.Errorf("cosmos commitments: %w", err)
	}
	if len(window) == 0 {
		return nil, nil
	}
	committed := make(map[uint64][]byte, len(window))
	sequences := make([]uint64, 0, len(window))
	for _, c := range window {
		committed[c.Sequence] = c.Data
		sequences = append(sequences, c.Sequence)
	}
	written, err := router.OutstandingSlots(ctx, sequences, func(seq uint64) []byte {
		return services.EthPath(s.ids.CosmosOnEVM, seq, services.EthPathAck)
	})
	if err != nil {
		return nil, err
	}
	// The cursor moves only once the probe has answered. A failed probe leaves
	// the window to the next pass; advancing first would skip these sequences
	// for a whole rotation on one RPC blip.
	s.flush.ackCursor = window[len(window)-1].Sequence

	var events []chain.Event
	for _, slot := range written {
		ev, found, err := router.WriteAckLog(ctx, s.ids.CosmosOnEVM, slot.Sequence, committed[slot.Sequence])
		if err != nil {
			s.logger.Printf("[eth->cosmos Flush] resolve ack seq=%d: %v", slot.Sequence, err)
			continue
		}
		if !found {
			continue
		}
		pkt := subscriber.EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
		if e, ok := ethPacketToEvent(services.EthPacket{
			Type: services.EthWriteAck, Packet: &pkt, AckBytes: ev.Acknowledgements, BlockNumber: ev.Raw.BlockNumber,
		}); ok {
			events = append(events, e)
		}
	}
	return events, nil
}

var _ chain.PacketLister = (*Source)(nil)
