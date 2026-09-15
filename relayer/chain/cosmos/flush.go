package cosmos

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"relayer/chain"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/gogoproto/proto"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"

	abcitypes "github.com/cometbft/cometbft/abci/types"
)

// flushQueryTimeout bounds one enumeration pass. It is generous because the pass
// runs every few minutes and a paginated commitment query against a busy chain
// is legitimately slow; the point of the bound is that a hung endpoint cannot
// wedge the flush worker, not to enforce a latency target.
const flushQueryTimeout = 60 * time.Second

// flushCommitmentPageSize caps one page of the commitment query. The healthy
// answer is a handful of sequences, so the page size only matters on a chain
// with a real backlog -- where paging is exactly what stops one query from
// returning everything at once.
const flushCommitmentPageSize = 100

// flushSequenceCap bounds how many outstanding sequences one pass will resolve
// into packets. Each sequence costs a TxSearch, so an unbounded backlog would
// turn a backstop into a stampede. The remainder is not lost: the cursor makes
// the next pass start where this one stopped.
//
// The cap bounds the COMMITMENT QUERY too, which is the part that used to be
// unbounded: the pass asks for a window of this size starting after the cursor
// rather than paging the whole set and slicing it. See outstandingCommitments.
const flushSequenceCap = 200

// flushTxSearchResults is how many transactions one sequence lookup will scan.
//
// One would do if every node matched a multi-condition query per EVENT, which
// CometBFT's kv indexer does since v0.38 -- its index key carries the event
// sequence, so `packet_sequence=N AND packet_source_client=X` must be satisfied
// by one event. That is a property of the NODE's indexer, not of this query: a
// psql-indexed node, or kv entries written in the legacy format (which record
// event sequence 0 for every attribute), match the two conditions across
// DIFFERENT events of the same transaction. Such a transaction then comes back
// first and carries no matching event, and asking for a single result would let
// it hide the real one. Scanning a few costs nothing on the healthy path, where
// the first result is the answer.
const flushTxSearchResults = 10

// flushTxSearchMaxResults bounds how far one sequence lookup will page.
//
// Scanning one page was not enough. On exactly the nodes the comment above
// describes, the stale cross-matches can FILL a page and leave the live send at
// result 11 or later -- the lookup then reports "not found", the cursor moves
// on, and the next pass asks the identical question and gets the identical
// answer. The packet is outstanding forever while every pass says there is
// nothing to do. Found in review.
//
// Paging to exhaustion is the other wrong answer: TotalCount on a
// cross-matching index is whatever that index chose to conflate, and this
// lookup runs once per outstanding sequence, up to flushSequenceCap of them per
// pass. A cap keeps the pass bounded, which is the property the commitment
// enumeration was just fixed to have. Reaching it is reported rather than
// silent, because "we looked at 100 transactions for one sequence" is a broken
// index and the operator is the only one who can fix it.
const flushTxSearchMaxResults = 100

const packetCommitmentsQueryPath = "/ibc.core.channel.v2.Query/PacketCommitments"

// UnrelayedPackets enumerates the packets this Cosmos source has sent that are
// still outstanding, so the relay module can find work it never observed.
//
// It implements chain.PacketLister; see that interface for why enumeration
// exists alongside block scanning.
//
// An outstanding packet COMMITMENT is the right signal: Cosmos deletes it when
// the packet is acknowledged or timed out, so its presence means the packet's
// lifecycle has not closed. It does NOT mean the destination never received it
// -- a delivered packet whose ack is still in flight still has a commitment --
// which is why the caller filters on the destination's receipt rather than this
// query pretending to know.
//
// Only the send leg is enumerated. The ack leg needs the same treatment and is
// not implemented here: enumerating acknowledgements Cosmos has written says
// nothing on its own about whether the counterparty consumed them, and
// establishing that needs a counterparty-side commitment query this adapter does
// not have. Tracked as the remaining half of the enumeration work.
func (s *Source) UnrelayedPackets(ctx context.Context) ([]chain.Event, error) {
	queryCtx, cancel := context.WithTimeout(ctx, flushQueryTimeout)
	defer cancel()

	window, err := s.outstandingCommitments(queryCtx)
	if err != nil {
		return nil, err
	}
	if len(window) == 0 {
		return nil, nil
	}
	s.advanceFlushCursor(window)

	events := make([]chain.Event, 0, len(window))
	for _, commitment := range window {
		event, found, err := s.sendEventForCommitment(queryCtx, commitment)
		if err != nil {
			// One unresolvable sequence must not discard the rest: the others are
			// exactly the packets this pass exists to recover.
			s.logger.Printf("[FlushCosmos] resolve send packet seq=%d: %v", commitment.sequence, err)
			continue
		}
		if !found {
			// The commitment exists but no indexed send_packet tx carries the packet
			// it commits to. On a pruned node that is expected for an old packet, and
			// it is worth saying so once per pass rather than silently returning a
			// shorter list.
			s.logger.Printf("[FlushCosmos] commitment for seq=%d has no indexed send_packet tx matching it (pruned history?)",
				commitment.sequence)
			continue
		}
		events = append(events, event)
	}
	return events, nil
}

// flushCommitment is one outstanding commitment: the sequence to resolve, and
// the commitment hash that says WHICH packet that sequence currently names.
//
// The hash is carried rather than dropped because a sequence alone does not
// identify a packet to anything outside the chain's own store. The commitment is
// what the chain will accept a proof against, so it is also the only thing that
// can confirm the transaction this pass finds is the packet still outstanding
// and not an older one that happens to share the sequence.
type flushCommitment struct {
	sequence uint64
	data     []byte
}

// outstandingCommitments returns at most flushSequenceCap commitments, starting
// after the sequence the previous pass stopped at and wrapping to the oldest.
//
// It ASKS FOR THE WINDOW rather than paging the whole set and slicing it. The
// difference is the whole cost of the pass. The commitment store is keyed by the
// big-endian sequence, and the SDK's paginator starts its iterator at the key it
// is given, so "the 200 sequences after N" is two page queries -- against a
// backlog of any size. Paging everything first made the query cost proportional
// to the backlog: a chain with 100k outstanding commitments spent a thousand
// serial ABCI round-trips, and could burn the whole pass timeout before relaying
// a single packet. The cap advertised on flushSequenceCap only ever bounded the
// TxSearch half; this is the half that was unbounded.
//
// What it gives up is the backlog total, which used to be logged. Counting is
// exactly the unbounded work being removed, and a number nobody can act on is
// not worth a thousand queries.
func (s *Source) outstandingCommitments(ctx context.Context) ([]flushCommitment, error) {
	// After the cursor: the packets this pass has not looked at most recently.
	window, err := s.pageCommitments(ctx, sequenceKey(s.flushCursor+1), flushSequenceCap, 0)
	if err != nil {
		return nil, err
	}

	// Short means the tail ran out, so wrap to the oldest. Bounded by the cursor
	// so the two segments cannot overlap: everything taken above is beyond it.
	if len(window) < flushSequenceCap && s.flushCursor > 0 {
		head, err := s.pageCommitments(ctx, nil, flushSequenceCap-len(window), s.flushCursor)
		if err != nil {
			return nil, err
		}
		window = append(window, head...)
	}

	if len(window) == flushSequenceCap {
		s.logger.Printf("[FlushCosmos] resolving %d outstanding commitments on %s this pass (seq %d..%d); "+
			"the window is full, so more may remain for the next pass",
			len(window), s.commitmentQueryClientID(), window[0].sequence, window[len(window)-1].sequence)
	}
	return window, nil
}

// pageCommitments reads up to want commitments starting at startKey, stopping
// early at stopAfter when that is non-zero.
//
// stopAfter is what keeps the wrap segment from re-reading what the segment
// after the cursor already took.
func (s *Source) pageCommitments(ctx context.Context, startKey []byte, want int, stopAfter uint64) ([]flushCommitment, error) {
	if want <= 0 {
		return nil, nil
	}
	var (
		out     []flushCommitment
		nextKey = startKey
	)
	for {
		limit := want - len(out)
		if limit > flushCommitmentPageSize {
			limit = flushCommitmentPageSize
		}
		request := &channeltypesv2.QueryPacketCommitmentsRequest{
			ClientId:   s.commitmentQueryClientID(),
			Pagination: &query.PageRequest{Key: nextKey, Limit: uint64(limit)},
		}
		raw, err := proto.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("cosmos source: marshal packet-commitments request: %w", err)
		}
		result, err := s.cosmos.CosmosClient().ABCIQuery(ctx, packetCommitmentsQueryPath, raw)
		if err != nil {
			return nil, fmt.Errorf("cosmos source: packet commitments: %w", err)
		}
		if result.Response.Code != 0 {
			return nil, fmt.Errorf("cosmos source: packet commitments: code %d: %s",
				result.Response.Code, result.Response.Log)
		}
		var response channeltypesv2.QueryPacketCommitmentsResponse
		if err := proto.Unmarshal(result.Response.Value, &response); err != nil {
			return nil, fmt.Errorf("cosmos source: decode packet-commitments response: %w", err)
		}
		for _, commitment := range response.Commitments {
			if stopAfter > 0 && commitment.Sequence > stopAfter {
				return out, nil
			}
			out = append(out, flushCommitment{sequence: commitment.Sequence, data: commitment.Data})
			if len(out) == want {
				return out, nil
			}
		}
		if response.Pagination == nil || len(response.Pagination.NextKey) == 0 {
			return out, nil
		}
		nextKey = response.Pagination.NextKey
	}
}

// sequenceKey renders a sequence as the paginator key for it: the commitment
// store is prefixed per client and keyed by the 8-byte big-endian sequence, so
// this is the "start here" the SDK iterator understands.
func sequenceKey(sequence uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, sequence)
	return key
}

// sendEventForCommitment finds the indexed send_packet transaction for one
// outstanding commitment and rebuilds the event the subscriber would have
// delivered.
//
// It searches by sequence rather than by height because the whole point of the
// enumeration path is that the height is unknown -- if it were known, the block
// scan would already have covered it. What comes back is then checked against
// the commitment, because the sequence is what was searched for and the
// commitment is what says which packet that sequence currently means.
func (s *Source) sendEventForCommitment(ctx context.Context, commitment flushCommitment) (chain.Event, bool, error) {
	q := fmt.Sprintf("send_packet.packet_sequence=%d AND send_packet.packet_source_client='%s'",
		commitment.sequence, s.commitmentQueryClientID())
	perPage := flushTxSearchResults
	scanned := 0
	for page := 1; ; page++ {
		result, err := s.cosmos.CosmosClient().TxSearch(ctx, q, false, &page, &perPage, "asc")
		if err != nil {
			return chain.Event{}, false, fmt.Errorf("tx search page %d: %w", page, err)
		}
		for _, tx := range result.Txs {
			packet, ok := matchSendEvent(tx.TxResult.Events, s.commitmentQueryClientID(), commitment)
			if !ok {
				continue
			}
			return chain.Event{
				Type:     chain.SendPacket,
				Height:   uint64(tx.Height),
				Sequence: packet.Sequence,
				ClientID: packet.DestinationClient,
				Raw:      mustMarshalPacket(packet),
			}, true, nil
		}
		scanned += len(result.Txs)
		// An empty page ends it regardless of what TotalCount claims: a node that
		// keeps reporting more than it will return must not spin this loop.
		if len(result.Txs) == 0 || scanned >= result.TotalCount {
			return chain.Event{}, false, nil
		}
		if scanned >= flushTxSearchMaxResults {
			log.Printf("[FlushCosmos] sequence %d on %s: scanned %d of %d indexed transactions without "+
				"finding the commitment the chain holds; giving up on this sequence for this pass. "+
				"That many matches for one (sequence, client) pair means the node's transaction "+
				"index is conflating events -- a psql index, or kv entries in the legacy format",
				commitment.sequence, s.commitmentQueryClientID(), scanned, result.TotalCount)
			return chain.Event{}, false, nil
		}
	}
}

// matchSendEvent returns the send_packet event whose packet is the one this
// commitment commits to.
//
// Three checks, and each closes a hole the others do not. TxSearch matches
// TRANSACTIONS, and the event list it hands back is the whole transaction's,
// unfiltered -- the query decides which transactions come back, never which
// events within one.
//
//   - SEQUENCE, because a transaction can carry several sends.
//   - SOURCE CLIENT, because sequences are per-client: seq=N exists independently
//     on every client, so a sequence-only match would answer with whichever
//     client's send appeared first in the transaction.
//   - COMMITMENT, because the first two identify a POSITION, not a packet. The
//     chain reuses a position: a client migration keeps a client id and restarts
//     sequences, and the transaction index keeps the old send. Ordering by "asc"
//     then hands back the OLDEST transaction at that position, so the flush would
//     rebuild a packet that closed long ago -- relaying something already settled
//     while the packet actually outstanding stays undiscovered, every pass,
//     identically. Comparing CommitPacket against the commitment the chain holds
//     right now is what distinguishes the two, and it is the same hash the chain
//     would check a proof against.
func matchSendEvent(events []abcitypes.Event, sourceClient string, commitment flushCommitment) (channeltypesv2.Packet, bool) {
	for _, event := range events {
		if event.Type != "send_packet" {
			continue
		}
		packet, ok := packetFromSendEvent(event.Attributes)
		if !ok || packet.Sequence != commitment.sequence || packet.SourceClient != sourceClient {
			continue
		}
		if !bytes.Equal(channeltypesv2.CommitPacket(packet), commitment.data) {
			continue
		}
		return packet, true
	}
	return channeltypesv2.Packet{}, false
}

// packetFromSendEvent decodes the hex-encoded packet an ibc-go v2 send_packet
// event carries. The encoded form is the contract the subscriber already
// consumes; decoding it here keeps one wire format rather than two.
func packetFromSendEvent(attributes []abcitypes.EventAttribute) (channeltypesv2.Packet, bool) {
	for _, attribute := range attributes {
		if attribute.Key != "encoded_packet_hex" {
			continue
		}
		raw, err := hex.DecodeString(attribute.Value)
		if err != nil {
			return channeltypesv2.Packet{}, false
		}
		var packet channeltypesv2.Packet
		if err := proto.Unmarshal(raw, &packet); err != nil {
			return channeltypesv2.Packet{}, false
		}
		return packet, true
	}
	return channeltypesv2.Packet{}, false
}

// mustMarshalPacket re-encodes a packet that was just decoded from the chain's
// own event, so a failure here would mean the packet round-tripped differently
// than it was produced. Returning nil rather than panicking keeps a malformed
// event from taking the process down; handleBatch rejects an empty payload.
func mustMarshalPacket(packet channeltypesv2.Packet) []byte {
	raw, err := proto.Marshal(&packet)
	if err != nil {
		return nil
	}
	return raw
}

var _ chain.PacketLister = (*Source)(nil)

// commitmentQueryClientID names the client a Cosmos send commitment is stored
// under, which is the packet's LOCAL source client.
//
// For a Cosmos->EVM packet that is EVMOnCosmos -- the EVM light client living on
// Cosmos -- and NOT CosmosOnEVM, which is how the EVM router knows Cosmos and is
// the packet's DESTINATION client. Getting it backwards is not a partial
// failure: the query returns an empty set, so the enumeration backstop reports
// nothing outstanding on precisely the path it was built for, and looks healthy
// doing it. It is a named method so the choice is testable rather than repeated
// as a field access at each call site.
func (s *Source) commitmentQueryClientID() string { return s.ids.EVMOnCosmos }

// advanceFlushCursor records where this pass stopped, so the next one starts
// after it.
//
// The rotation is the point. Taking the oldest flushSequenceCap every pass looks
// fair and is not: a prefix that cannot resolve -- pruned history, or packets
// already delivered and waiting on an acknowledgement, which keep their
// commitment -- occupies the whole window forever, and every later packet is
// never even looked at. Rotating means an unresolvable prefix costs one pass,
// not all of them.
//
// A set that fits inside one window ends every pass at the same last sequence,
// so the cursor parks there and the wrap hands the whole set back next time.
// That is the intended behaviour for a healthy chain, not a stalled cursor.
func (s *Source) advanceFlushCursor(window []flushCommitment) {
	if len(window) == 0 {
		return
	}
	s.flushCursor = window[len(window)-1].sequence
}
