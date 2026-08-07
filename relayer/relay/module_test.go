package relay

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"relayer/chain"
)

// --- mock adapters (the whole point: Module is testable without any chain) ---

type mockSource struct {
	latest                uint64
	relayable             uint64 // highest provable height; if 0, defaults to latest
	membershipCalls       int
	nonMembershipCalls    int
	proofHeights          []uint64
	failMembershipOn      map[string]bool // Raw payload -> return a retryable proof error
	permanentMembershipOn map[string]bool // Raw payload -> return a permanent proof error
}

func (m *mockSource) Chain() chain.ChainType { return chain.Cosmos }
func (m *mockSource) Subscribe(context.Context, func(context.Context, []chain.Event) []int) error {
	return nil
}
func (m *mockSource) LatestHeight(context.Context) (uint64, error) { return m.latest, nil }
func (m *mockSource) RelayableHeight(context.Context) (uint64, error) {
	if m.relayable != 0 {
		return m.relayable, nil
	}
	return m.latest, nil
}
func (m *mockSource) QueryHeader(_ context.Context, h uint64) ([]byte, error) {
	return []byte{byte(h)}, nil
}
func (m *mockSource) MembershipProof(_ context.Context, packet []byte, height uint64, _ chain.EventType) ([]byte, error) {
	m.membershipCalls++
	m.proofHeights = append(m.proofHeights, height)
	if m.permanentMembershipOn[string(packet)] {
		return nil, chain.Permanent(errors.New("timed out"))
	}
	if m.failMembershipOn[string(packet)] {
		return nil, chain.Retryable(errors.New("proof unavailable"))
	}
	return []byte("membership"), nil
}
func (m *mockSource) NonMembershipProof(_ context.Context, _ []byte, height uint64) ([]byte, error) {
	m.nonMembershipCalls++
	m.proofHeights = append(m.proofHeights, height)
	return []byte("non-membership"), nil
}

type mockDest struct {
	updates    []chain.ClientUpdate // recorded UpdateClient calls
	relayed    []chain.RelayPacket  // recorded RelayPackets calls
	relayCalls int                  // number of RelayPackets invocations (multicall folding check)
	relayErr   error                // if set, RelayPackets returns it (transient)
	// poison: Raw payload -> this packet deterministically reverts any batch it is
	// in (chain.Permanent), like a timed-out/duplicate packet in a real multicall.
	poison map[string]bool
}

func (m *mockDest) Chain() chain.ChainType { return chain.Ethereum }
func (m *mockDest) UpdateClient(_ context.Context, _ string, u chain.ClientUpdate) error {
	m.updates = append(m.updates, u)
	return nil
}
func (m *mockDest) RelayPackets(_ context.Context, packets []chain.RelayPacket) error {
	m.relayCalls++
	if m.relayErr != nil {
		return m.relayErr
	}
	for _, p := range packets {
		if m.poison[string(p.Packet)] {
			return chain.Permanent(errors.New("poison packet reverted the batch"))
		}
	}
	m.relayed = append(m.relayed, packets...)
	return nil
}
func (m *mockDest) HasPacketReceipt(context.Context, []byte) (bool, error) { return false, nil }
func (m *mockDest) ClientExpiresAt(context.Context, string) (time.Time, error) {
	return time.Time{}, nil
}

type foldingMockDest struct {
	mockDest
	enabled       bool
	foldCalls     int
	foldedUpdates []chain.ClientUpdate
	foldedPackets [][]chain.RelayPacket
	foldErr       error
	foldFn        func(chain.ClientUpdate, []chain.RelayPacket) error
}

func (m *foldingMockDest) SupportsUpdatePacketFolding() bool { return m.enabled }
func (m *foldingMockDest) RelayWithUpdate(_ context.Context, _ string, update chain.ClientUpdate, packets []chain.RelayPacket) error {
	m.foldCalls++
	m.foldedUpdates = append(m.foldedUpdates, update)
	m.foldedPackets = append(m.foldedPackets, packets)
	if m.foldFn != nil {
		return m.foldFn(update, packets)
	}
	return m.foldErr
}

// mockBuilder returns an update whose Height is the height it was asked to prove,
// unless overridden (to simulate a "no update needed" zero-height result).
type mockBuilder struct {
	forceHeight uint64 // if non-zero, always return this height
	noPayload   bool   // if true, return no payloads ("client already current")
	built       int    // count of Build calls
}

func (m *mockBuilder) Name() string { return "mock" }
func (m *mockBuilder) Build(_ context.Context, header []byte) (chain.ClientUpdate, error) {
	m.built++
	h := uint64(header[0])
	if m.forceHeight != 0 {
		h = m.forceHeight
	}
	if m.noPayload {
		return chain.ClientUpdate{Height: h}, nil
	}
	return chain.ClientUpdate{Height: h, Payloads: [][]byte{header}}, nil
}

func TestHandleBatch_ProofByType(t *testing.T) {
	src := &mockSource{latest: 100} // chain tip well above the packet heights
	dst := &mockDest{}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	// recv/ack packets prove membership (commitment exists)
	if rq := m.handleBatch(context.Background(), []chain.Event{{Type: chain.SendPacket, Height: 5, Raw: []byte("pkt-a")}}); len(rq) != 0 {
		t.Fatalf("clean send must not re-queue, got %v", rq)
	}
	if src.membershipCalls != 1 || src.nonMembershipCalls != 0 {
		t.Fatalf("send packet must use membership proof: m=%d nm=%d", src.membershipCalls, src.nonMembershipCalls)
	}
	if len(dst.relayed) != 1 || string(dst.relayed[0].Proof) != "membership" {
		t.Fatalf("relayed packet must carry membership proof, got %+v", dst.relayed)
	}

	// timeouts prove non-membership (receipt absent)
	if rq := m.handleBatch(context.Background(), []chain.Event{{Type: chain.TimeoutPacket, Height: 6, Raw: []byte("pkt-b")}}); len(rq) != 0 {
		t.Fatalf("clean timeout must not re-queue, got %v", rq)
	}
	if src.nonMembershipCalls != 1 {
		t.Fatalf("timeout must use non-membership proof: nm=%d", src.nonMembershipCalls)
	}
	if len(dst.relayed) != 2 || string(dst.relayed[1].Proof) != "non-membership" {
		t.Fatalf("timeout relayed packet must carry non-membership proof, got %+v", dst.relayed)
	}
}

// TestHandleBatch_Multicall verifies a multi-event batch folds into ONE
// RelayPackets call (the multicall that amortizes per-tx intrinsic gas).
func TestHandleBatch_Multicall(t *testing.T) {
	src := &mockSource{latest: 100} // chain tip well above the packet heights
	dst := &mockDest{}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	batch := []chain.Event{
		{Type: chain.SendPacket, Height: 7, Raw: []byte("a")},
		{Type: chain.SendPacket, Height: 9, Raw: []byte("b")},
		{Type: chain.AckPacket, Height: 8, Raw: []byte("c")},
	}
	if rq := m.handleBatch(context.Background(), batch); len(rq) != 0 {
		t.Fatalf("clean batch must not re-queue, got %v", rq)
	}
	// One client update advancing to the source tip (100), one submit for all 3.
	if len(dst.updates) != 1 || dst.updates[0].Height != 100 {
		t.Fatalf("want single update at source tip 100, got %+v", dst.updates)
	}
	if dst.relayCalls != 1 {
		t.Fatalf("want packets folded into ONE RelayPackets call, got %d calls", dst.relayCalls)
	}
	if len(dst.relayed) != 3 {
		t.Fatalf("want 3 packets relayed, got %d", len(dst.relayed))
	}
}

func TestHandleBatch_FoldsUpdateAndProofsAtTargetHeight(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &foldingMockDest{enabled: true}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{forceHeight: 80})
	m.lastHeight = 10

	batch := []chain.Event{
		{Type: chain.SendPacket, Height: 30, Raw: []byte("a")},
		{Type: chain.AckPacket, Height: 50, Raw: []byte("b")},
	}
	if rq := m.handleBatch(context.Background(), batch); len(rq) != 0 {
		t.Fatalf("clean folded batch must not re-queue, got %v", rq)
	}
	if dst.foldCalls != 1 || len(dst.foldedUpdates) != 1 || dst.foldedUpdates[0].Height != 80 {
		t.Fatalf("want one folded update at target height 80, calls=%d updates=%+v", dst.foldCalls, dst.foldedUpdates)
	}
	if len(dst.updates) != 0 || dst.relayCalls != 0 {
		t.Fatalf("folded path must not submit standalone txs: updates=%d relays=%d", len(dst.updates), dst.relayCalls)
	}
	if len(src.proofHeights) != 2 || src.proofHeights[0] != 80 || src.proofHeights[1] != 80 {
		t.Fatalf("proofs must target not-yet-on-chain update height 80, got %v", src.proofHeights)
	}
	for _, packet := range dst.foldedPackets[0] {
		if packet.Height != 80 {
			t.Fatalf("folded packet height = %d, want 80", packet.Height)
		}
	}
	if m.lastHeight != 80 {
		t.Fatalf("successful atomic tx must advance lastHeight to 80, got %d", m.lastHeight)
	}
}

func TestHandleBatch_FoldingCapabilityDisabledFallsBack(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &foldingMockDest{enabled: false}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	if rq := m.handleBatch(context.Background(), []chain.Event{{Type: chain.SendPacket, Height: 20, Raw: []byte("p")}}); len(rq) != 0 {
		t.Fatalf("fallback batch must not re-queue, got %v", rq)
	}
	if dst.foldCalls != 0 || len(dst.updates) != 1 || dst.relayCalls != 1 {
		t.Fatalf("disabled folding must use two-tx fallback: folds=%d updates=%d relays=%d", dst.foldCalls, len(dst.updates), dst.relayCalls)
	}
}

func TestHandleBatch_FoldFailureDoesNotAdvanceState(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &foldingMockDest{enabled: true, foldErr: errors.New("rpc unavailable")}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	rq := m.handleBatch(context.Background(), []chain.Event{{Type: chain.SendPacket, Height: 20, Raw: []byte("p")}})
	if len(rq) != 1 || rq[0] != 0 {
		t.Fatalf("transient folded failure must re-queue packet, got %v", rq)
	}
	if m.lastHeight != 0 {
		t.Fatalf("failed atomic tx must not advance lastHeight, got %d", m.lastHeight)
	}
	if len(dst.updates) != 0 || dst.relayCalls != 0 {
		t.Fatalf("failed folded path must not fall through to standalone txs")
	}
}

func TestHandleBatch_NoUpdateNeededRelaysPacketsOnly(t *testing.T) {
	src := &mockSource{latest: 200}
	dst := &foldingMockDest{enabled: true}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{noPayload: true, forceHeight: 200})

	if rq := m.handleBatch(context.Background(), []chain.Event{{Type: chain.SendPacket, Height: 180, Raw: []byte("p")}}); len(rq) != 0 {
		t.Fatalf("current-client batch must relay, got %v", rq)
	}
	if dst.foldCalls != 0 || len(dst.updates) != 0 || dst.relayCalls != 1 {
		t.Fatalf("no-update branch must submit packets only: folds=%d updates=%d relays=%d", dst.foldCalls, len(dst.updates), dst.relayCalls)
	}
	if len(src.proofHeights) != 1 || src.proofHeights[0] != 200 {
		t.Fatalf("no-update proof height = %v, want [200]", src.proofHeights)
	}
}

func TestHandleBatch_FoldedIsolationPreservesValidSiblings(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &foldingMockDest{enabled: true}
	dst.poison = map[string]bool{"bad": true}
	dst.foldFn = func(_ chain.ClientUpdate, packets []chain.RelayPacket) error {
		if len(packets) > 1 {
			return chain.Permanent(errors.New("folded batch reverted"))
		}
		return nil
	}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	batch := []chain.Event{
		{Type: chain.SendPacket, Height: 7, Raw: []byte("good-1")},
		{Type: chain.SendPacket, Height: 8, Raw: []byte("bad")},
		{Type: chain.SendPacket, Height: 9, Raw: []byte("good-2")},
	}
	if rq := m.handleBatch(context.Background(), batch); len(rq) != 0 {
		t.Fatalf("folded isolation must drop only poison, got requeue %v", rq)
	}
	if dst.foldCalls != 2 { // failed batch, then first successful singleton installs update
		t.Fatalf("folded update must be retried only until installed, got %d fold calls", dst.foldCalls)
	}
	if m.lastHeight != 100 {
		t.Fatalf("successful singleton fold must advance lastHeight, got %d", m.lastHeight)
	}
	if len(dst.relayed) != 1 || string(dst.relayed[0].Packet) != "good-2" {
		t.Fatalf("packets after update should use packets-only isolation, got %+v", dst.relayed)
	}
}

// TestHandleBatch_UntracksDeliveredSends verifies that after a successful relay
// the module removes each delivered SendPacket from the pending tracker (via the
// untrack hook) — and only SendPacket, not acks — so the timeout scanner stops
// considering a packet that can no longer time out.
func TestHandleBatch_UntracksDeliveredSends(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &mockDest{}
	var untracked [][]byte
	m := NewModule("test", "client-0", src, dst, &mockBuilder{},
		WithPacketTracker(
			func([]byte, uint64) {}, // track: no-op for this test
			func(pkt []byte) { untracked = append(untracked, pkt) },
		),
	)

	batch := []chain.Event{
		{Type: chain.SendPacket, Height: 7, Raw: []byte("send-1")},
		{Type: chain.AckPacket, Height: 8, Raw: []byte("ack-1")},
		{Type: chain.SendPacket, Height: 9, Raw: []byte("send-2")},
	}
	if rq := m.handleBatch(context.Background(), batch); len(rq) != 0 {
		t.Fatalf("clean batch must not re-queue, got %v", rq)
	}
	if len(untracked) != 2 {
		t.Fatalf("want 2 sends untracked, got %d (%q)", len(untracked), untracked)
	}
	if string(untracked[0]) != "send-1" || string(untracked[1]) != "send-2" {
		t.Fatalf("untracked the wrong packets (acks must not be untracked): %q", untracked)
	}
}

// TestHandleBatch_NoUntrackOnRelayFailure verifies a failed relay does NOT untrack
// any packet — a packet that was not delivered must stay in the tracker so the
// timeout scanner can still refund it.
func TestHandleBatch_NoUntrackOnRelayFailure(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &mockDest{relayErr: errors.New("rpc down")} // transient relay failure
	var untracked [][]byte
	m := NewModule("test", "client-0", src, dst, &mockBuilder{},
		WithPacketTracker(
			func([]byte, uint64) {},
			func(pkt []byte) { untracked = append(untracked, pkt) },
		),
	)

	rq := m.handleBatch(context.Background(), []chain.Event{{Type: chain.SendPacket, Height: 7, Raw: []byte("send-1")}})
	if len(rq) != 1 {
		t.Fatalf("transient relay failure must re-queue the send, got %v", rq)
	}
	if len(untracked) != 0 {
		t.Fatalf("a failed relay must not untrack anything, got %q", untracked)
	}
}

// TestHandleBatch_IsolatesPoisonPacket verifies that when a multi-packet batch
// hits a permanent (deterministic revert) failure, the module isolates the poison
// packet by retrying each packet individually — so the poison is dropped alone and
// its VALID siblings still relay, rather than the whole batch being discarded.
func TestHandleBatch_IsolatesPoisonPacket(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &mockDest{poison: map[string]bool{"bad": true}}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	batch := []chain.Event{
		{Type: chain.SendPacket, Height: 7, Raw: []byte("good-1")},
		{Type: chain.SendPacket, Height: 8, Raw: []byte("bad")},
		{Type: chain.SendPacket, Height: 9, Raw: []byte("good-2")},
	}
	if rq := m.handleBatch(context.Background(), batch); len(rq) != 0 {
		t.Fatalf("isolation drops only the poison and re-queues nothing, got %v", rq)
	}
	var relayed []string
	for _, p := range dst.relayed {
		relayed = append(relayed, string(p.Packet))
	}
	if len(relayed) != 2 || relayed[0] != "good-1" || relayed[1] != "good-2" {
		t.Fatalf("valid siblings must still relay (poison dropped), got %v", relayed)
	}
}

// TestHandleBatch_IsolationRequeuesTransient verifies that if an isolated retry
// hits a TRANSIENT failure, that packet's index is re-queued (not dropped), while
// the poison packet is still dropped.
func TestHandleBatch_IsolationRequeuesTransient(t *testing.T) {
	src := &mockSource{latest: 100}
	// Whole batch reverts permanently (poison "bad"); when retried alone, "good-1"
	// fails transiently and "bad" fails permanently.
	dst := &transientOnIsolate{poison: "bad", transientFor: "good-1"}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})
	batch := []chain.Event{
		{Type: chain.SendPacket, Height: 7, Raw: []byte("good-1")},
		{Type: chain.SendPacket, Height: 8, Raw: []byte("bad")},
	}
	rq := m.handleBatch(context.Background(), batch)
	// good-1 (index 0) transient → re-queued; bad (index 1) permanent → dropped.
	if len(rq) != 1 || rq[0] != 0 {
		t.Fatalf("transient isolated retry must re-queue only good-1's index [0], got %v", rq)
	}
}

// transientOnIsolate is a destination where the whole batch reverts permanently
// (poison present) but the named packet fails TRANSIENTLY when retried alone.
type transientOnIsolate struct {
	mockDest
	poison       string
	transientFor string
}

func (d *transientOnIsolate) RelayPackets(_ context.Context, packets []chain.RelayPacket) error {
	for _, p := range packets {
		if string(p.Packet) == d.poison {
			return chain.Permanent(errors.New("poison"))
		}
	}
	if len(packets) == 1 && string(packets[0].Packet) == d.transientFor {
		return errors.New("rpc blip") // transient
	}
	return nil
}

// TestHandleBatch_CancelledContextRelaysNothing verifies that a cancelled context
// short-circuits handleBatch: everything is re-queued and no proof/relay work runs.
func TestHandleBatch_CancelledContextRelaysNothing(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &mockDest{}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	batch := []chain.Event{
		{Type: chain.SendPacket, Height: 7, Raw: []byte("a")},
		{Type: chain.SendPacket, Height: 8, Raw: []byte("b")},
	}
	rq := m.handleBatch(ctx, batch)
	if len(rq) != 2 {
		t.Fatalf("cancelled ctx must re-queue every event, got %v", rq)
	}
	if src.membershipCalls != 0 || dst.relayCalls != 0 || len(dst.updates) != 0 {
		t.Fatalf("cancelled ctx must not run proof/update/relay work: m=%d relay=%d upd=%d",
			src.membershipCalls, dst.relayCalls, len(dst.updates))
	}
}

// TestHandleBatch_RequeueTransient verifies a per-packet transient proof failure
// re-queues exactly that event's index while the rest still relay.
func TestHandleBatch_RequeueTransient(t *testing.T) {
	src := &mockSource{latest: 100, failMembershipOn: map[string]bool{"bad": true}}
	dst := &mockDest{}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	batch := []chain.Event{
		{Type: chain.SendPacket, Height: 3, Raw: []byte("ok1")},
		{Type: chain.SendPacket, Height: 4, Raw: []byte("bad")},
		{Type: chain.SendPacket, Height: 5, Raw: []byte("ok2")},
	}
	rq := m.handleBatch(context.Background(), batch)
	if len(rq) != 1 || rq[0] != 1 {
		t.Fatalf("want re-queue of index 1 only, got %v", rq)
	}
	if len(dst.relayed) != 2 {
		t.Fatalf("want the 2 good packets relayed, got %d", len(dst.relayed))
	}
}

// TestHandleBatch_DropPermanentProof verifies a permanent proof failure (a
// timed-out packet) is DROPPED, not re-queued — the fix for the gas-drain bug
// where a dead packet spun forever, each retry re-running a client update.
func TestHandleBatch_DropPermanentProof(t *testing.T) {
	src := &mockSource{latest: 100, permanentMembershipOn: map[string]bool{"dead": true}}
	dst := &mockDest{}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	batch := []chain.Event{
		{Type: chain.SendPacket, Height: 3, Raw: []byte("ok")},
		{Type: chain.SendPacket, Height: 4, Raw: []byte("dead")}, // timed out → permanent
	}
	rq := m.handleBatch(context.Background(), batch)
	if len(rq) != 0 {
		t.Fatalf("permanent packet must be DROPPED, not re-queued, got %v", rq)
	}
	if len(dst.relayed) != 1 || string(dst.relayed[0].Packet) != "ok" {
		t.Fatalf("only the live packet should relay, got %+v", dst.relayed)
	}
}

// TestHandleBatch_DropPermanentRelay verifies a permanent relay (on-chain revert)
// failure drops the batch instead of re-queueing it forever.
func TestHandleBatch_DropPermanentRelay(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &mockDest{relayErr: chain.Permanent(errors.New("reverted"))}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	rq := m.handleBatch(context.Background(), []chain.Event{{Type: chain.SendPacket, Height: 5, Raw: []byte("p")}})
	if len(rq) != 0 {
		t.Fatalf("permanent relay revert must DROP, not re-queue, got %v", rq)
	}
}

// TestHandleBatch_RequeueTransientRelay verifies a transient relay failure (RPC,
// nonce) re-queues the batch (unlike a permanent revert).
func TestHandleBatch_RequeueTransientRelay(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &mockDest{relayErr: errors.New("rpc unavailable")} // not permanent
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	rq := m.handleBatch(context.Background(), []chain.Event{{Type: chain.SendPacket, Height: 5, Raw: []byte("p")}})
	if len(rq) != 1 {
		t.Fatalf("transient relay failure must re-queue, got %v", rq)
	}
}

// TestHandleBatch_PreconditionSkip verifies the RelayableHeight precondition: a
// packet whose source height is above the relayable height (Cosmos AppHash lag /
// ETH finality lag) is re-queued WITHOUT any client update or proof — the fix for
// the re-prove/spam regression where the module burned a proof every flush while
// waiting.
func TestHandleBatch_PreconditionSkip(t *testing.T) {
	src := &mockSource{latest: 100, relayable: 50} // tip 100, but only <=50 provable
	dst := &mockDest{}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	rq := m.handleBatch(context.Background(), []chain.Event{{Type: chain.SendPacket, Height: 80, Raw: []byte("pending")}})
	if len(rq) != 1 || rq[0] != 0 {
		t.Fatalf("packet above relayable height must re-queue, got %v", rq)
	}
	if len(dst.updates) != 0 {
		t.Fatalf("must NOT advance the client for an unprovable batch, got %d updates", len(dst.updates))
	}
	if src.membershipCalls != 0 {
		t.Fatalf("must NOT build a proof for an unprovable batch, got %d", src.membershipCalls)
	}
}

// TestHandleBatch_ProvabilityGuard verifies that when the destination client can
// only advance to a height below the packet's (the ETH->Cosmos finality-lag case),
// the packet is re-queued instead of proven against a state that lacks it.
func TestHandleBatch_ProvabilityGuard(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &mockDest{}
	// Builder can only advance the client to height 50 (finality lag), regardless
	// of the requested height.
	m := NewModule("test", "client-0", src, dst, &mockBuilder{forceHeight: 50})

	// Packet observed at execution height 100, but the client only reaches 50.
	rq := m.handleBatch(context.Background(), []chain.Event{{Type: chain.SendPacket, Height: 100, Raw: []byte("future")}})
	if len(rq) != 1 || rq[0] != 0 {
		t.Fatalf("packet above provable height must re-queue, got %v", rq)
	}
	if src.membershipCalls != 0 {
		t.Fatalf("must not attempt a proof for an unprovable packet, got %d", src.membershipCalls)
	}
	if len(dst.relayed) != 0 {
		t.Fatalf("nothing should be relayed, got %d", len(dst.relayed))
	}
}

// TestUpdateClientTo_SeedOnNoop verifies an empty-payload-list ("already current")
// build advances lastHeight without submitting a tx, so a restart where the
// client already covers incoming packets does not stall the provability guard.
func TestUpdateClientTo_SeedOnNoop(t *testing.T) {
	src := &mockSource{latest: 200}
	dst := &mockDest{}
	b := &mockBuilder{noPayload: true, forceHeight: 200} // "client current at 200"
	m := NewModule("test", "client-0", src, dst, b)

	if err := m.updateClientTo(context.Background(), 150); err != nil {
		t.Fatalf("seed update: %v", err)
	}
	if len(dst.updates) != 0 {
		t.Fatalf("no-op build must not submit a tx, got %d", len(dst.updates))
	}
	if m.lastHeight != 200 {
		t.Fatalf("lastHeight must be seeded to reported current height 200, got %d", m.lastHeight)
	}
	// A packet at 180 is now provable (<= seeded 200) even though no tx was sent.
	if rq := m.handleBatch(context.Background(), []chain.Event{{Type: chain.SendPacket, Height: 180, Raw: []byte("p")}}); len(rq) != 0 {
		t.Fatalf("packet under seeded height must relay, got re-queue %v", rq)
	}
}

// TestPeriodicUpdateLoop_Runs verifies the fixed-cadence forced update fires
// (the pinned-set rotation guarantee) and keeps firing until the context ends.
func TestPeriodicUpdateLoop_Runs(t *testing.T) {
	var calls atomic.Int32
	m := NewModule("test", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
		WithPeriodicUpdate(10*time.Millisecond, 0, func(context.Context) error {
			calls.Add(1)
			return nil
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	go m.periodicUpdateLoop(ctx)

	deadline := time.After(2 * time.Second)
	for calls.Load() < 2 {
		select {
		case <-deadline:
			cancel()
			t.Fatalf("periodicUpdate fired %d times, want >= 2", calls.Load())
		case <-time.After(5 * time.Millisecond):
		}
	}
	cancel()
}

// TestPeriodicUpdateLoop_RetriesOnFailure verifies a failing run does not stop the
// loop — it keeps retrying (on the shorter retry cadence).
func TestPeriodicUpdateLoop_RetriesOnFailure(t *testing.T) {
	var calls atomic.Int32
	m := NewModule("test", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
		WithPeriodicUpdate(10*time.Millisecond, 0, func(context.Context) error {
			calls.Add(1)
			return errors.New("rotate failed")
		}),
	)
	// Shrink the retry cadence for the test so failures re-fire quickly.
	m.periodicUpdateInterval = 10 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	m.periodicUpdateLoop(ctx) // returns when ctx times out

	if calls.Load() < 2 {
		t.Fatalf("failing periodicUpdate must keep retrying, fired %d times", calls.Load())
	}
}

// TestNextPeriodicUpdateBackoff verifies the exponential failure backoff sequence
// (1m→2m→4m→8m→15m cap) that prevents a persistently-failing client-update tx
// from draining the signer's gas.
func TestNextPeriodicUpdateBackoff(t *testing.T) {
	want := []time.Duration{
		1 * time.Minute, // first failure
		2 * time.Minute, // doubling...
		4 * time.Minute,
		8 * time.Minute,
		15 * time.Minute, // capped (would be 16m)
		15 * time.Minute, // stays capped
	}
	cur := time.Duration(0) // no prior failure
	for i, w := range want {
		cur = nextPeriodicUpdateBackoff(cur)
		if cur != w {
			t.Fatalf("backoff step %d: got %s, want %s", i, cur, w)
		}
	}
	// A success (reset to 0) restarts the sequence at the minimum.
	if got := nextPeriodicUpdateBackoff(0); got != periodicUpdateBackoffMin {
		t.Fatalf("reset backoff: got %s, want %s", got, periodicUpdateBackoffMin)
	}
}

func TestUpdateClientTo_AppendOnly(t *testing.T) {
	src := &mockSource{}
	dst := &mockDest{}
	b := &mockBuilder{}
	m := NewModule("test", "client-0", src, dst, b)
	ctx := context.Background()

	// 1. first update at height 10 -> submitted, lastHeight=10
	if err := m.updateClientTo(ctx, 10); err != nil {
		t.Fatalf("update to 10: %v", err)
	}
	if len(dst.updates) != 1 || dst.updates[0].Height != 10 {
		t.Fatalf("want one update at height 10, got %+v", dst.updates)
	}

	// 2. re-request height 10 (<= lastHeight) -> skipped before Build
	builtBefore := b.built
	if err := m.updateClientTo(ctx, 10); err != nil {
		t.Fatalf("re-update to 10: %v", err)
	}
	if len(dst.updates) != 1 {
		t.Fatalf("append-only violated: expected no new submit, got %d", len(dst.updates))
	}
	if b.built != builtBefore {
		t.Fatalf("expected Build skipped for covered height, but it ran")
	}

	// 3. height 20 but builder reports the client already current (Height 0)
	b.forceHeight = 0 // header[0]=20 -> but simulate no-op: override to 0? no: use a builder that returns stale
	b.forceHeight = 5 // stale (<= lastHeight 10): must be skipped without submit
	if err := m.updateClientTo(ctx, 20); err != nil {
		t.Fatalf("update to 20 (stale build): %v", err)
	}
	if len(dst.updates) != 1 {
		t.Fatalf("stale/no-op update must not submit, got %d submits", len(dst.updates))
	}

	// 4. height 20, builder returns a real advance -> submitted, lastHeight=20
	b.forceHeight = 0 // back to echoing the requested height (20)
	if err := m.updateClientTo(ctx, 20); err != nil {
		t.Fatalf("update to 20: %v", err)
	}
	if len(dst.updates) != 2 || dst.updates[1].Height != 20 {
		t.Fatalf("want second update at height 20, got %+v", dst.updates)
	}
}

// A packet waiting on the source frontier is re-queued every flush — every batch
// period, on no backoff. A measured OP Sepolia wait produced 67 identical lines in
// 4m12s, and a default 150-block attestor gap makes ~100 the normal case. The line
// must therefore report transitions, not repetitions: the same wait says nothing
// new, a frontier that moved does.
func TestHandleBatch_WaitingLogsOnlyOnChange(t *testing.T) {
	src := &mockSource{latest: 100, relayable: 5}
	m := NewModule("test", "client-0", src, &mockDest{}, &mockBuilder{})
	waiting := []chain.Event{{Type: chain.SendPacket, Height: 50, Raw: []byte("pkt-far")}}

	if rq := m.handleBatch(context.Background(), waiting); len(rq) != 1 {
		t.Fatalf("packet above the relayable height must re-queue, got %v", rq)
	}
	first := m.lastWait
	if !first.logged || first.packets != 1 || first.pendingMax != 50 || first.relayable != 5 {
		t.Fatalf("first wait must be recorded as logged: %+v", first)
	}

	// Same picture again: nothing to say.
	m.handleBatch(context.Background(), waiting)
	if m.lastWait != first {
		t.Fatalf("an unchanged wait must not update the log state: %+v -> %+v", first, m.lastWait)
	}

	// Frontier advances — that is the signal worth printing.
	src.relayable = 20
	m.handleBatch(context.Background(), waiting)
	if m.lastWait == first {
		t.Fatal("a moving frontier must be reported")
	}
	if m.lastWait.relayable != 20 {
		t.Fatalf("log state must track the new frontier, got %+v", m.lastWait)
	}

	// Once everything is relayable the state resets, so the next wait reports
	// itself even if it looks identical to the one before.
	src.relayable = 100
	if rq := m.handleBatch(context.Background(), waiting); len(rq) != 0 {
		t.Fatalf("packet below the relayable height must not re-queue, got %v", rq)
	}
	if m.lastWait.logged {
		t.Fatalf("a cleared wait must reset the log state, got %+v", m.lastWait)
	}
}
