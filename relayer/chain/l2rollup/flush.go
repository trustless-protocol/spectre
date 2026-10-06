package l2rollup

import (
	"context"
	"fmt"
	"log"
	"time"

	"relayer/chain"
	"relayer/chain/evmflush"
	relayerclient "relayer/client"
	"relayer/services"
)

// flushState is the enumeration backstop's per-source state: one rotating
// cursor per leg and the log-range policy the lookups narrow in place.
type flushState struct {
	sendCursor uint64 // lowest sequence the last send probe reached; 0 = newest first
	ackCursor  uint64 // last Cosmos commitment sequence the ack leg took
	span       relayerclient.LogSpan
	// seeded records that span took its starting chunk from logScanChunk. The
	// scan loop seeds its own LogSpan the same way; both narrow from there.
	seeded bool
}

// flushQueryTimeout bounds one enumeration pass.
const flushQueryTimeout = 60 * time.Second

// cosmosCommitmentLister is the optional half of CosmosPacketState the
// acknowledgement leg needs: page the Cosmos commitment store for the client
// this L2 is paired with. The reconciliation seam asks point questions; this is
// the only enumeration it needs, and a test double that does not implement it
// simply has no acknowledgements to list.
type cosmosCommitmentLister interface {
	CommitmentWindow(ctx context.Context, clientID string, cursor uint64, limit int) ([]services.PacketCommitment, error)
}

func (s cosmosEndpointState) CommitmentWindow(ctx context.Context, clientID string, cursor uint64, limit int) ([]services.PacketCommitment, error) {
	return services.CosmosCommitmentWindow(ctx, s.endpoint, clientID, cursor, limit)
}

// UnrelayedPackets enumerates what this L2 source still owes the Cosmos side.
// It implements chain.PacketLister; the two legs and their bounds are the ones
// the ETH source describes, against the L2's router:
//
//   - sends: the L2 router still holds a commitment under this source's client;
//   - acknowledgements: Cosmos still holds a send's commitment under the paired
//     wasm client AND the L2 router has written an acknowledgement for it.
//
// This is the leg that needs the backstop most: the L2 scan window is the
// shortest of the four paths (#308), and an acknowledgement written on the L2
// while the relayer was down was, before this, gone for good.
func (s *Source) UnrelayedPackets(ctx context.Context) ([]chain.Event, error) {
	if s.eth == nil {
		return nil, fmt.Errorf("l2 source: enumeration needs an L2 client")
	}
	if !s.flush.seeded {
		s.flush.span.Chunk = s.logScanChunk
		s.flush.seeded = true
	}
	queryCtx, cancel := context.WithTimeout(ctx, flushQueryTimeout)
	defer cancel()

	router := evmflush.Router{
		Client: s.eth, Address: s.router, Span: &s.flush.span, DeployBlock: s.routerDeployBlock, Label: "[l2->cosmos Flush]",
	}
	sends, err := s.unrelayedSends(queryCtx, router)
	if err != nil {
		return nil, err
	}
	acks, err := s.unrelayedAcks(queryCtx, router)
	if err != nil {
		log.Printf("[l2->cosmos Flush] acknowledgement leg: %v; the send leg's %d packet(s) still stand", err, len(sends))
	}
	return append(sends, acks...), nil
}

func (s *Source) unrelayedSends(ctx context.Context, router evmflush.Router) ([]chain.Event, error) {
	prev, err := router.PrevSequenceSend(ctx, s.l2ClientID)
	if err != nil {
		return nil, err
	}
	lo, hi, ok := evmflush.Window(prev, s.flush.sendCursor, evmflush.ProbeCap)
	if !ok {
		return nil, nil
	}
	slots, err := router.OutstandingSlots(ctx, evmflush.Sequences(lo, hi), func(seq uint64) []byte {
		return services.EthPath(s.l2ClientID, seq, services.EthPathCommitment)
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
		ev, found, err := router.SendLog(ctx, s.l2ClientID, slot.Sequence, slot.Value)
		if err != nil {
			log.Printf("[l2->cosmos Flush] resolve send seq=%d: %v", slot.Sequence, err)
			continue
		}
		if !found {
			continue
		}
		// Past its timeout it can never be received on Cosmos; the L2 timeout
		// scanner refunds it. The same rule MembershipProof applies, applied here
		// so an old expired send does not cost a proof build per pass -- and so
		// this leg matches the ETH source's.
		if ev.Packet.TimeoutTimestamp > 0 && uint64(time.Now().Unix()) >= ev.Packet.TimeoutTimestamp {
			continue
		}
		// l2SendToEvent applies the same destination-client filter the scan does.
		if e, ok := l2SendToEvent(ev, s.cosmosWasmClientID); ok {
			events = append(events, e)
		}
	}
	return events, nil
}

func (s *Source) unrelayedAcks(ctx context.Context, router evmflush.Router) ([]chain.Event, error) {
	lister, ok := s.cosmosState.(cosmosCommitmentLister)
	if !ok || s.cosmosWasmClientID == "" {
		return nil, nil
	}
	// Cosmos->L2 packets are committed on Cosmos under the paired wasm client,
	// and acknowledged on the L2 under this router's client.
	window, err := lister.CommitmentWindow(ctx, s.cosmosWasmClientID, s.flush.ackCursor, evmflush.ProbeCap)
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
		return services.EthPath(s.l2ClientID, seq, services.EthPathAck)
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
		ev, found, err := router.WriteAckLog(ctx, s.l2ClientID, slot.Sequence, committed[slot.Sequence])
		if err != nil {
			log.Printf("[l2->cosmos Flush] resolve ack seq=%d: %v", slot.Sequence, err)
			continue
		}
		if !found {
			continue
		}
		if e, ok := l2AckToEvent(ev, s.cosmosWasmClientID); ok {
			events = append(events, e)
		}
	}
	return events, nil
}

var _ chain.PacketLister = (*Source)(nil)
