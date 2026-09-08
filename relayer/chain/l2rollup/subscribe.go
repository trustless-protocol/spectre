package l2rollup

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"relayer/chain"
	"relayer/subscriber"

	contractICS26Router "relayer/bindings/ICS26Router"

	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// A standalone poll-based L2 event listener (deliberately NOT the ETH
// subscriber.SubscribeEth, which is bound to the L1 endpoint + BatchBuilder via
// shared context). It polls eth_getLogs on the L2 RPC for the ICS26Router
// SendPacket / WriteAcknowledgement events of this source's client id, decodes them
// with the shared ICS26Router bindings, and drives the adapter handler with its own
// pending buffer for re-queued (not-yet-relayable) events. Within one run the cursor is
// gap-recovery-safe: it only advances past a range that scanned successfully, so an RPC
// hiccup re-scans rather than skips (mistake #11).
//
// The cursor is NOT persisted across runs — it lives only in Subscribe's stack frame.
// A restart resumes at `head - L2_STARTUP_LOOKBACK_BLOCKS` (see below), so anything older than
// that window is never re-scanned, and a packet whose acknowledgement fell outside it
// stays pending forever with nothing in the log to say why. A value of zero starts at
// the head and disables startup recovery.
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
	// l2StartupLookbackEnv optionally extends the recovery window below the head at
	// startup.  It matters when an operator restarts after the default window: an
	// unpersisted cursor cannot otherwise discover an older, still-unrelayed packet.
	l2StartupLookbackEnv = "L2_STARTUP_LOOKBACK_BLOCKS"

	// defaultL2StartupLookback rescans a window below the head at startup so packets emitted
	// while the relayer was down are picked up. Re-emitting an already-relayed packet
	// is safe — Cosmos rejects the duplicate recv and the module drops it permanently.
	//
	// NOTE: this fixed 256-block (~8 min on OP) window bounds crash-recovery. For a
	// finalized frontier that can lag hours behind the L2 head, a send that finalizes
	// only after a longer downtime would fall outside the window; size it against the
	// worst-case attestor/finality lag for the configured head-kind, or add the
	// receipt-checked recovery the ETH mirror has, before relying on it in production.
	defaultL2StartupLookback = uint64(256)
)

func l2StartupLookbackBlocksFromEnv(raw string) uint64 {
	if raw == "" {
		return defaultL2StartupLookback
	}

	lookback, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		log.Printf("[recovery] ignoring invalid %s=%q; using %d", l2StartupLookbackEnv, raw, defaultL2StartupLookback)
		return defaultL2StartupLookback
	}
	return lookback
}

func l2StartupLookbackBlocks() uint64 {
	return l2StartupLookbackBlocksFromEnv(os.Getenv(l2StartupLookbackEnv))
}

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
	lookback := l2StartupLookbackBlocks()
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
			if head > lookback {
				from = head - lookback
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
		var settled map[settledKey]struct{}
		if head >= from {
			fresh, settled, err = s.scanPacketLogs(ctx, filterer, from, head)
			if err != nil {
				log.Printf("[SubscribeL2] scan [%d,%d]: %v", from, head, err)
				continue // do NOT advance the cursor on failure (re-scan next tick)
			}
			from = head + 1 // range consumed; fresh events are now carried in the batch
		}

		// Filter BEFORE the handler, and filter the whole batch. A packet settled
		// by this very scan must not reach handleBatch: it would be relayed for
		// nothing and, worse, tracked again -- after settle removed it and after
		// the cursor moved past the terminal log that would clear it. See
		// dropSettled.
		batch := dropSettled(append(pending[:len(pending):len(pending)], fresh...), settled)
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
func (s *Source) scanPacketLogs(ctx context.Context, filterer *contractICS26Router.ContractICS26RouterFilterer, from, to uint64) ([]chain.Event, map[settledKey]struct{}, error) {
	if chunk := s.logScanChunk; chunk > 0 && to >= from && to-from >= chunk {
		var all []chain.Event
		// One set across every span: a send in span 1 can be settled by a terminal
		// log in span 3, and chunking must not hide that from the filter.
		settled := map[settledKey]struct{}{}
		for start := from; start <= to; start += chunk {
			end := start + chunk - 1
			if end > to {
				end = to
			}
			events, spanSettled, err := s.scanPacketLogRange(ctx, filterer, start, end)
			if err != nil {
				return nil, nil, fmt.Errorf("span [%d,%d]: %w", start, end, err)
			}
			all = append(all, events...)
			for key := range spanSettled {
				settled[key] = struct{}{}
			}
		}
		return all, settled, nil
	}
	return s.scanPacketLogRange(ctx, filterer, from, to)
}

// scanPacketLogRange is one eth_getLogs pair over a span the provider will serve.
func (s *Source) scanPacketLogRange(ctx context.Context, filterer *contractICS26Router.ContractICS26RouterFilterer, from, to uint64) ([]chain.Event, map[settledKey]struct{}, error) {
	opts := &bind.FilterOpts{Start: from, End: &to, Context: ctx}
	clientFilter := []string{s.l2ClientID}

	var events []chain.Event

	sends, err := filterer.FilterSendPacket(opts, clientFilter, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("filter SendPacket: %w", err)
	}
	defer sends.Close()
	for sends.Next() {
		if e, ok := l2SendToEvent(sends.Event, s.cosmosWasmClientID); ok {
			events = append(events, e)
		}
	}
	if err := sends.Error(); err != nil {
		return nil, nil, fmt.Errorf("iterate SendPacket: %w", err)
	}

	acks, err := filterer.FilterWriteAcknowledgement(opts, clientFilter, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("filter WriteAcknowledgement: %w", err)
	}
	defer acks.Close()
	for acks.Next() {
		if e, ok := l2AckToEvent(acks.Event, s.cosmosWasmClientID); ok {
			events = append(events, e)
		}
	}
	if err := acks.Error(); err != nil {
		return nil, nil, fmt.Errorf("iterate WriteAcknowledgement: %w", err)
	}

	settled, err := s.scanTerminalLogs(opts, filterer, clientFilter)
	if err != nil {
		return nil, nil, err
	}
	return events, settled, nil
}

// scanTerminalLogs reads the two events that CLOSE a packet's lifecycle on the
// L2 -- AckPacket and TimeoutPacket -- and settles them.
//
// They return no chain.Event on purpose. A terminal event creates no relay work,
// and feeding one to the relay loop would produce an empty "relay" of a packet
// with nothing left to do. Their entire value is that they are free: without
// them a packet another relayer acknowledged stays in our pending tracker until
// a timeout scan happens to query for it.
//
// A nil settle hook makes this a no-op rather than an error: reading the logs
// still costs two eth_getLogs calls, but a source that does not own a tracker
// has nothing to do with them.
func (s *Source) scanTerminalLogs(opts *bind.FilterOpts, filterer *contractICS26Router.ContractICS26RouterFilterer, clientFilter []string) (map[settledKey]struct{}, error) {
	if s.settle == nil {
		return nil, nil
	}
	settled := map[settledKey]struct{}{}

	acked, err := filterer.FilterAckPacket(opts, clientFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("filter AckPacket: %w", err)
	}
	defer acked.Close()
	for acked.Next() {
		s.settleTerminal("AckPacket", acked.Event.Packet, acked.Event.Sequence, settled)
	}
	if err := acked.Error(); err != nil {
		return nil, fmt.Errorf("iterate AckPacket: %w", err)
	}

	timedOut, err := filterer.FilterTimeoutPacket(opts, clientFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("filter TimeoutPacket: %w", err)
	}
	defer timedOut.Close()
	for timedOut.Next() {
		s.settleTerminal("TimeoutPacket", timedOut.Event.Packet, timedOut.Event.Sequence, settled)
	}
	if err := timedOut.Error(); err != nil {
		return nil, fmt.Errorf("iterate TimeoutPacket: %w", err)
	}
	return settled, nil
}

// settledKey identifies a packet across the two shapes this file handles: the
// chain.Event a send produced, and the router log a terminal event carries.
//
// It is the MARSHALED PACKET, not the client and sequence. Both shapes reach it
// through the same EthPacketToCosmosPacket + proto.Marshal, so the bytes agree
// whenever the packets do -- and disagree when they do not. That matters because
// a client migration keeps the client id and restarts sequences from 1, so
// seq=N names two different packets over the life of one client; a key of
// client+sequence would let a stale terminal event drop the NEW seq=N from the
// batch, and nothing downstream would notice it was never relayed.
type settledKey string

// dropSettled removes packets whose lifecycle closed inside the very range this
// scan just read.
//
// Settling only the tracker is not enough, and the gap is not a small one. A
// scan window that contains BOTH a send and its later AckPacket -- ordinary
// after a restart, since the window is sized to cover the downtime -- settles
// the packet and then hands the send to the relay loop anyway. handleBatch
// relays it and TRACKS IT AGAIN, so the tracker entry comes back after settle
// removed it, and the cursor has already moved past the terminal log that would
// have cleared it. Nothing settles it a second time: the entry stays pending for
// the life of the process and the timeout scanner keeps querying it.
//
// The same applies to a send already sitting in the retry queue that another
// relayer acknowledged, which is why both halves of the batch are filtered and
// not just the fresh ones.
func dropSettled(batch []chain.Event, settled map[settledKey]struct{}) []chain.Event {
	if len(settled) == 0 || len(batch) == 0 {
		return batch
	}
	kept := batch[:0:0]
	for _, e := range batch {
		if _, done := settled[settledKey(e.Raw)]; done {
			log.Printf("[SubscribeL2] dropping %s seq=%d from this batch: it settled inside the same scan range",
				e.Type, e.Sequence)
			continue
		}
		kept = append(kept, e)
	}
	return kept
}

// settleTerminal converts one terminal log to the packet identity the tracker
// keys on and hands it to the hook.
func (s *Source) settleTerminal(kind string, packet contractICS26Router.IICS26RouterMsgsPacket, sequence *big.Int, settled map[settledKey]struct{}) {
	cosmosPacket := subscriber.EthPacketToCosmosPacket(packet, sequence)
	raw, err := proto.Marshal(&cosmosPacket)
	if err != nil {
		// The packet came from the chain's own log, so this cannot happen without
		// the binding and the proto type having diverged. Say so rather than
		// settling a packet identity nobody can reproduce.
		log.Printf("[SubscribeL2] %s seq=%d: marshal for settlement: %v", kind, cosmosPacket.Sequence, err)
		return
	}
	log.Printf("[SubscribeL2] %s seq=%d settled; dropped from the pending tracker", kind, cosmosPacket.Sequence)
	settled[settledKey(raw)] = struct{}{}
	s.settle(raw)
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
