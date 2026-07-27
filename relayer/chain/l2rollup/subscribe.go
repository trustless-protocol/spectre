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
// pending buffer for re-queued (not-yet-relayable) events. Polling from a persisted
// cursor is inherently gap-recovery-safe: the cursor only advances past a range that
// scanned successfully, so a crash/RPC hiccup re-scans rather than skips (mistake #11).
//
// LIMITATION (unsafe head-kind): the cursor only moves FORWARD — it never rewinds on an
// L2 reorg. With head_kind=unsafe a scanned block can be reorged out and replaced; an
// event only present on the losing branch is then missed (the cursor already passed
// it). This is a non-issue for safe/finalized head-kinds (the relayable frontier is
// post-reorg), which is the reason unsafe is the lowest-trust setting.
const (
	l2SubscribeInterval = 4 * time.Second
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

	head, err := s.head(ctx)
	if err != nil {
		return fmt.Errorf("l2 source: initial head: %w", err)
	}
	from := uint64(0)
	if head > l2StartupLookback {
		from = head - l2StartupLookback
	}
	log.Printf("[SubscribeL2] polling ICS26Router %s from block %d (client_id=%s)", s.router.Hex(), from, s.l2ClientID)

	ticker := time.NewTicker(l2SubscribeInterval)
	defer ticker.Stop()

	var pending []chain.Event // events the handler re-queued (not yet relayable)
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
		if head < from {
			continue
		}

		fresh, err := s.scanPacketLogs(ctx, filterer, from, head)
		if err != nil {
			log.Printf("[SubscribeL2] scan [%d,%d]: %v", from, head, err)
			continue // do NOT advance the cursor on failure (re-scan next tick)
		}
		from = head + 1 // range consumed; fresh events are now carried in the batch

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
// client id in [from,to] and maps them to chain.Events.
func (s *Source) scanPacketLogs(ctx context.Context, filterer *contractICS26Router.ContractICS26RouterFilterer, from, to uint64) ([]chain.Event, error) {
	opts := &bind.FilterOpts{Start: from, End: &to, Context: ctx}
	clientFilter := []string{s.l2ClientID}

	var events []chain.Event

	sends, err := filterer.FilterSendPacket(opts, clientFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("filter SendPacket: %w", err)
	}
	defer sends.Close()
	for sends.Next() {
		if e, ok := l2SendToEvent(sends.Event); ok {
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
		if e, ok := l2AckToEvent(acks.Event); ok {
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
func l2SendToEvent(ev *contractICS26Router.ContractICS26RouterSendPacket) (chain.Event, bool) {
	pkt := subscriber.EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
	raw, err := pkt.Marshal()
	if err != nil {
		return chain.Event{}, false
	}
	return chain.Event{
		Type:     chain.SendPacket,
		Height:   ev.Raw.BlockNumber,
		ClientID: pkt.DestinationClient,
		Raw:      raw,
	}, true
}

// l2AckToEvent maps a WriteAcknowledgement log to an AckPacket chain.Event. An empty
// acknowledgement is skipped — it cannot be built into a MsgAcknowledgement (mirrors
// the ETH source's no-ack skip).
func l2AckToEvent(ev *contractICS26Router.ContractICS26RouterWriteAcknowledgement) (chain.Event, bool) {
	if len(ev.Acknowledgements) == 0 {
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
