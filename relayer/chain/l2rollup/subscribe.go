package l2rollup

import (
	"context"
	"fmt"
	"log"
	"time"

	"relayer/chain"
	"relayer/subscriber"

	contractICS26Router "relayer/bindings/ICS26Router"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// A standalone poll-based L2 event listener (deliberately NOT the ETH
// subscriber.SubscribeEth, which is bound to the L1 endpoint + BatchBuilder via
// services.Context). It polls eth_getLogs on the L2 RPC for the ICS26Router
// SendPacket / WriteAcknowledgement events of this source's client id, decodes them
// with the shared ICS26Router bindings, and drives the adapter handler with its own
// pending buffer for re-queued (not-yet-relayable) events. Within one run the cursor is
// gap-recovery-safe: it only advances past a range that scanned successfully, so an RPC
// hiccup re-scans rather than skips (mistake #11).
//
// The cursor is NOT persisted across runs — it lives only in Subscribe's stack frame.
// A restart resumes at `head - l2StartupLookback` (see below), so anything older than
// that window is never re-scanned, and a packet whose acknowledgement fell outside it
// stays pending forever with nothing in the log to say why.
//
// LIMITATION (unsafe head-kind): the cursor only moves FORWARD — it never rewinds on an
// L2 reorg. With head_kind=unsafe a scanned block can be reorged out and replaced; an
// event only present on the losing branch is then missed (the cursor already passed
// it). This is a non-issue for safe/finalized head-kinds (the relayable frontier is
// post-reorg), which is the reason unsafe is the lowest-trust setting.
// l2SubscribeInterval is a var, not a const, purely so tests can shorten the poll
// period; nothing outside the package reassigns it.
var l2SubscribeInterval = 4 * time.Second

const (
	// l2StartupLookback rescans a window below the head at startup so packets emitted
	// while the relayer was down are picked up. Re-emitting an already-relayed packet
	// is safe — Cosmos rejects the duplicate recv and the module drops it permanently.
	//
	// NOTE: this fixed 256-block (~8 min on OP) window bounds crash-recovery. For a
	// finalized frontier that can lag hours behind the L2 head, a send that finalizes
	// only after a longer downtime would fall outside the window; size it against the
	// worst-case attestor/finality lag for the configured head-kind, or add the
	// receipt-checked recovery the ETH mirror has, before relying on it in production.
	l2StartupLookback = uint64(256)
)

// Subscribe polls the L2 for ICS26Router packet events and drives handler in batches
// until ctx is cancelled.
func (s *Source) Subscribe(ctx context.Context, handler func(context.Context, []chain.Event) []int) error {
	filterer, err := contractICS26Router.NewContractICS26RouterFilterer(s.router, s.eth)
	if err != nil {
		return fmt.Errorf("l2 source: new ICS26Router filterer: %w", err)
	}

	ticker := time.NewTicker(l2SubscribeInterval)
	defer ticker.Stop()

	// The cursor is seeded on the FIRST SUCCESSFUL head read inside the loop, not
	// before it. Reading it up front made a transient RPC error at startup fatal
	// while the identical error one tick later was logged and retried — the process
	// died on a rate-limit blip during boot, which is exactly when an operator can
	// least tell a transient failure from a misconfiguration. The mirrors already
	// get this right: SubscribeCosmos and SubscribeEth both read their starting
	// height inside the retry loop and only return early for PERMANENT conditions
	// (missing websocket URL, filterer construction) — the same distinction the
	// filterer check above makes. Errors are classified by kind, never by position.
	var (
		from    uint64
		seeded  bool
		pending []chain.Event // events the handler re-queued (not yet relayable)
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}

		head, err := s.head(ctx)
		if err != nil {
			log.Printf("[SubscribeL2] head: %v", err)
			continue // do NOT advance the cursor on failure
		}
		if !seeded {
			if head > l2StartupLookback {
				from = head - l2StartupLookback
			}
			seeded = true
			log.Printf("[SubscribeL2] polling ICS26Router %s from block %d (client_id=%s)",
				s.router.Hex(), from, s.l2ClientID)
		}

		// Only the SCAN is gated on there being a new block range — the pending
		// buffer must be re-offered on every tick regardless. A demand-driven
		// rollup (Arbitrum Nitro seals a block per transaction) stops at its last
		// block when traffic stops, so `head == from-1` becomes the steady state.
		// Skipping the whole iteration there strands every re-queued packet until
		// the next transaction happens to arrive, and the acknowledgement is the
		// worst case: it is the last event its packet writes on the L2, so nothing
		// further will be written and no new block will ever come. The ack then
		// waits forever, its escrow stays locked, and the log says nothing after
		// the first "waiting" line.
		var fresh []chain.Event
		if head >= from {
			fresh, err = s.scanPacketLogs(ctx, filterer, from, head)
			if err != nil {
				log.Printf("[SubscribeL2] scan [%d,%d]: %v", from, head, err)
				continue // do NOT advance the cursor on failure (re-scan next tick)
			}
			from = head + 1 // range consumed; fresh events are now carried in the batch
		}

		batch := append(pending[:len(pending):len(pending)], fresh...)
		if len(batch) == 0 {
			pending = nil
			continue
		}
		requeue := handler(ctx, batch)
		pending = keepIndexed(batch, requeue)
	}
}

// scanPacketLogs fetches the SendPacket + WriteAcknowledgement logs of this source's
// client id in [from,to] and maps them to chain.Events, splitting the request into
// log_scan_chunk-sized spans when one is configured.
//
// A failure in any span fails the whole scan, so the caller leaves its cursor
// untouched: a partially scanned range must never be mistaken for a complete one.
func (s *Source) scanPacketLogs(ctx context.Context, filterer *contractICS26Router.ContractICS26RouterFilterer, from, to uint64) ([]chain.Event, error) {
	if chunk := s.logScanChunk; chunk > 0 && to >= from && to-from >= chunk {
		var all []chain.Event
		for start := from; start <= to; start += chunk {
			end := start + chunk - 1
			if end > to {
				end = to
			}
			events, err := s.scanPacketLogRange(ctx, filterer, start, end)
			if err != nil {
				return nil, fmt.Errorf("span [%d,%d]: %w", start, end, err)
			}
			all = append(all, events...)
		}
		return all, nil
	}
	return s.scanPacketLogRange(ctx, filterer, from, to)
}

// scanPacketLogRange is one eth_getLogs pair over a span the provider will serve.
func (s *Source) scanPacketLogRange(ctx context.Context, filterer *contractICS26Router.ContractICS26RouterFilterer, from, to uint64) ([]chain.Event, error) {
	opts := &bind.FilterOpts{Start: from, End: &to, Context: ctx}
	clientFilter := []string{s.l2ClientID}

	var events []chain.Event

	sends, err := filterer.FilterSendPacket(opts, clientFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("filter SendPacket: %w", err)
	}
	defer sends.Close()
	for sends.Next() {
		if e, ok := l2SendToEvent(sends.Event, s.cosmosWasmClientID); ok {
			events = append(events, e)
		}
	}
	if err := sends.Error(); err != nil {
		return nil, fmt.Errorf("iterate SendPacket: %w", err)
	}

	acks, err := filterer.FilterWriteAcknowledgement(opts, clientFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("filter WriteAcknowledgement: %w", err)
	}
	defer acks.Close()
	for acks.Next() {
		if e, ok := l2AckToEvent(acks.Event, s.cosmosWasmClientID); ok {
			events = append(events, e)
		}
	}
	if err := acks.Error(); err != nil {
		return nil, fmt.Errorf("iterate WriteAcknowledgement: %w", err)
	}
	return events, nil
}

// l2SendToEvent maps a SendPacket log to a recv (SendPacket) chain.Event, reusing the
// ETH source's packet decoder. Returns false if the packet cannot be marshaled.
func l2SendToEvent(ev *contractICS26Router.ContractICS26RouterSendPacket, cosmosWasmClientID string) (chain.Event, bool) {
	if cosmosWasmClientID != "" && ev.Packet.DestClient != cosmosWasmClientID {
		return chain.Event{}, false
	}
	pkt := subscriber.EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
	raw, err := pkt.Marshal()
	if err != nil {
		return chain.Event{}, false
	}
	return chain.Event{
		Type:     chain.SendPacket,
		Height:   ev.Raw.BlockNumber,
		Sequence: pkt.Sequence,
		ClientID: pkt.DestinationClient,
		Raw:      raw,
	}, true
}

// l2AckToEvent maps a WriteAcknowledgement log to an AckPacket chain.Event. An empty
// acknowledgement is skipped — it cannot be built into a MsgAcknowledgement (mirrors
// the ETH source's no-ack skip).
func l2AckToEvent(ev *contractICS26Router.ContractICS26RouterWriteAcknowledgement, cosmosWasmClientID string) (chain.Event, bool) {
	if len(ev.Acknowledgements) == 0 {
		return chain.Event{}, false
	}
	if cosmosWasmClientID != "" && ev.Packet.SourceClient != cosmosWasmClientID {
		return chain.Event{}, false
	}
	pkt := subscriber.EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
	raw, err := pkt.Marshal()
	if err != nil {
		return chain.Event{}, false
	}
	return chain.Event{
		Type:     chain.AckPacket,
		Height:   ev.Raw.BlockNumber,
		Sequence: pkt.Sequence,
		ClientID: pkt.DestinationClient,
		Raw:      raw,
		AckBytes: ev.Acknowledgements,
	}, true
}

// keepIndexed returns the events at the given indices (the handler's re-queue set),
// preserving order and ignoring out-of-range indices.
func keepIndexed(events []chain.Event, indices []int) []chain.Event {
	if len(indices) == 0 {
		return nil
	}
	kept := make([]chain.Event, 0, len(indices))
	for _, i := range indices {
		if i >= 0 && i < len(events) {
			kept = append(kept, events[i])
		}
	}
	return kept
}
