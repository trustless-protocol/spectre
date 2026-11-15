package relay

import (
	"context"
	"fmt"
	"log"
	"time"

	"attestor/types/redact"

	"relayer/chain"
)

// The module's side of the drop ledger (services.DropLedger, #575 M2).
//
// A DROP is not the end of a packet: its commitment is still outstanding on the
// source, so the flush backstop lists it again five minutes later. These helpers
// are what make a drop last: the three DROP sites record it, the batch and the
// flush skip a held packet before any client update, and a packet that lands or
// is found settled is forgotten.
//
// Every helper is a no-op when no ledger is wired (WithDropLedger), which keeps
// a mock or a test path on the old behaviour.

func (m *Module) now() time.Time {
	if m.clock != nil {
		return m.clock()
	}
	return time.Now()
}

// heldAfterDrop reports whether the packet is inside the backoff of an earlier
// permanent failure.
func (m *Module) heldAfterDrop(typ chain.EventType, raw []byte) bool {
	if m.dropped == nil {
		return false
	}
	_, held := m.dropped.Held(typ, raw, m.now())
	return held
}

// recordDrop notes a permanent failure and returns the text a DROP line appends
// to say when the packet is retried.
//
// The reason is stored with its URLs cut to scheme and host: /status serves it
// and the ledger file keeps it, and an RPC error quotes the endpoint URL, API
// key included (#575 L7).
func (m *Module) recordDrop(typ chain.EventType, sequence uint64, raw []byte, cause error) string {
	if m.dropped == nil {
		return ""
	}
	entry, err := m.dropped.Record(typ, sequence, raw, redact.Error(cause).Error(), m.now())
	if err != nil {
		log.Printf("ERROR [%s Relay][ATTENTION] dropped %s seq=%d is held in memory only: %v", m.name, typ, sequence, err)
	}
	return fmt.Sprintf("; held until %s (drop %d)", entry.NextRetryAt.UTC().Format(time.RFC3339), entry.Attempts)
}

// clearDrop forgets a packet that was delivered or is settled.
func (m *Module) clearDrop(typ chain.EventType, raw []byte) {
	if m.dropped == nil {
		return
	}
	if _, err := m.dropped.Clear(typ, raw); err != nil {
		log.Printf("WARN [%s Relay] clear dropped %s record: %v", m.name, typ, err)
	}
}

// dropHeld removes the events held by the ledger and reports how many it held.
func (m *Module) dropHeld(events []chain.Event) ([]chain.Event, int) {
	if m.dropped == nil {
		return events, 0
	}
	free := make([]chain.Event, 0, len(events))
	for _, e := range events {
		if m.heldAfterDrop(e.Type, e.Raw) {
			continue
		}
		free = append(free, e)
	}
	return free, len(events) - len(free)
}

// reconcileDropped forgets dropped packets that closed without this process.
//
// The source stops listing a packet once it is settled, so nothing would ever
// offer it to the batch again and its entry would stay in status as stranded.
// A send can also close on the source alone: an expired one is refunded by the
// timeout scanner and never gets a receipt (see dropClosed).
//
// Only entries whose backoff has expired are asked about. One the flush still
// lists is retried in the same pass; one nothing offers again (an expired ETH
// or L2 send waiting for its refund) is asked on every pass until it closes.
func (m *Module) reconcileDropped(ctx context.Context) {
	if m.dropped == nil {
		return
	}
	for _, entry := range m.dropped.Due(m.now()) {
		if ctx.Err() != nil {
			return
		}
		event := chain.Event{Type: entry.Type, Sequence: entry.Sequence, Raw: entry.Packet}
		closed, err := m.dropClosed(ctx, event)
		if err != nil || !closed {
			continue // unknown or still outstanding: the flush offers it again
		}
		m.clearDrop(entry.Type, entry.Packet)
		log.Printf("[%s Relay] dropped %s seq=%d has closed; forgotten", m.name, entry.Type, entry.Sequence)
	}
}

// dropClosed reports whether a dropped packet needs nothing more from this
// relay: the destination settled it or, for a send, the source no longer holds
// its commitment (acknowledged, or timed out and refunded).
//
// The commitment is read at the source's latest state. A wrong "gone" (a reorg
// restores the commitment) costs one early retry and nothing else: the entry
// only holds the packet out of the relay path.
func (m *Module) dropClosed(ctx context.Context, e chain.Event) (bool, error) {
	settled, err := m.packetSettled(ctx, e)
	if err != nil || settled || e.Type != chain.SendPacket {
		return settled, err
	}
	checker, ok := m.src.(chain.CommitmentChecker)
	if !ok {
		return false, nil
	}
	committed, err := checker.HasPacketCommitment(ctx, e.Raw)
	if err != nil {
		return false, err
	}
	return !committed, nil
}

// DroppedPacket is one packet held after a permanent failure, as status shows it.
type DroppedPacket struct {
	Type      string    `json:"type"`
	Sequence  uint64    `json:"sequence"`
	Since     time.Time `json:"since"`
	Drops     int       `json:"drops"`
	NextRetry time.Time `json:"next_retry"`
	Reason    string    `json:"reason"`
}

func (m *Module) droppedSnapshot() []DroppedPacket {
	if m.dropped == nil {
		return nil
	}
	entries := m.dropped.Snapshot()
	out := make([]DroppedPacket, 0, len(entries))
	for _, entry := range entries {
		out = append(out, DroppedPacket{
			Type:      entry.Type.String(),
			Sequence:  entry.Sequence,
			Since:     entry.FirstDroppedAt,
			Drops:     entry.Attempts,
			NextRetry: entry.NextRetryAt,
			Reason:    entry.Reason,
		})
	}
	return out
}
