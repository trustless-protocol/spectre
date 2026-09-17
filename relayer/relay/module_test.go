package relay

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
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
	// relayableErr, when set, is returned by RelayableHeight instead of a height.
	// relayableCalls counts how often the source was actually asked, which is the
	// thing the permanent-failure hold is about.
	relayableErr   error
	relayableCalls int
}

func (m *mockSource) Chain() chain.ChainType { return chain.Cosmos }
func (m *mockSource) Subscribe(context.Context, func(context.Context, []chain.Event) []int) error {
	return nil
}
func (m *mockSource) LatestHeight(context.Context) (uint64, error) { return m.latest, nil }
func (m *mockSource) RelayableHeight(context.Context) (uint64, error) {
	m.relayableCalls++
	if m.relayableErr != nil {
		return 0, m.relayableErr
	}
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
		return nil, chain.Transient(errors.New("proof unavailable"))
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
	updateErr  error                // if set, UpdateClient returns it
	relayed    []chain.RelayPacket  // recorded RelayPackets calls
	relayCalls int                  // number of RelayPackets invocations (multicall folding check)
	relayErr   error                // if set, RelayPackets returns it (transient)
	// relayCommit is how many packets a FAILING RelayPackets reports as already on
	// chain, so a partial commit can be driven. Read only when relayErr is set.
	relayCommit int
	// hasReceipt / receiptErr drive HasPacketReceipt for the flush tests: an
	// outstanding source commitment does not say whether the destination already
	// delivered the packet, so the flush asks.
	hasReceipt bool
	receiptErr error
	// receiptAnswers, when non-empty, is consumed one per HasPacketReceipt call
	// and falls back to hasReceipt once exhausted. It is how a test makes a packet
	// become delivered BETWEEN the flush's two receipt checks.
	receiptAnswers []bool
	// poison: Raw payload -> this packet deterministically reverts any batch it is
	// in (chain.Permanent), like a timed-out/duplicate packet in a real multicall.
	poison map[string]bool
	// expiresAt / trustingPeriod / expiresErr drive ClientExpiresAt for the
	// anti-expiry refresh tests. A zero trustingPeriod means the destination could
	// not report one, which is what makes the module fall back to its default margin.
	trustingPeriod time.Duration
	expiresAt      time.Time
	expiresErr     error
}

func (m *mockDest) Chain() chain.ChainType { return chain.Ethereum }
func (m *mockDest) UpdateClient(_ context.Context, _ string, u chain.ClientUpdate) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.updates = append(m.updates, u)
	return nil
}
func (m *mockDest) RelayPackets(_ context.Context, packets []chain.RelayPacket) (int, error) {
	m.relayCalls++
	if m.relayErr != nil {
		// A reported prefix has to be RECORDED as relayed too. A stub that claims
		// a commit without delivering anything would let the engine mis-settle and
		// still pass.
		//
		// The COUNT is returned unclamped on purpose. Clamping it here would make
		// the stub incapable of the contract violation one subtest is about, and
		// that subtest would then pass with the engine's own clamp deleted. Only
		// the recording is bounded: a stub cannot deliver what it was not given.
		m.relayed = append(m.relayed, packets[:min(m.relayCommit, len(packets))]...)
		return m.relayCommit, m.relayErr
	}
	for _, p := range packets {
		if m.poison[string(p.Packet)] {
			// This mock is an EVM destination: one multicall, so a revert commits
			// nothing.
			return 0, chain.Permanent(errors.New("poison packet reverted the batch"))
		}
	}
	m.relayed = append(m.relayed, packets...)
	return len(packets), nil
}
func (m *mockDest) HasPacketReceipt(context.Context, []byte) (bool, error) {
	if len(m.receiptAnswers) > 0 {
		answer := m.receiptAnswers[0]
		m.receiptAnswers = m.receiptAnswers[1:]
		return answer, m.receiptErr
	}
	return m.hasReceipt, m.receiptErr
}
func (m *mockDest) ClientExpiresAt(context.Context, string) (time.Time, time.Duration, error) {
	return m.expiresAt, m.trustingPeriod, m.expiresErr
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
	// hops is the ladder a churn-blocked builder walks: each Build returns the
	// next entry and marks it IsHop, until the list runs out and the last update
	// is the real one. Empty means "never hops", which is every adapter but the
	// Groth16 one.
	hops []uint64
	// trustedAt is the timestamp the built update claims the client now trusts.
	// It is what the client-update observer reports, so a test that checks the
	// observer must be able to set it to something distinguishable.
	trustedAt time.Time
}

func (m *mockBuilder) Name() string { return "mock" }
func (m *mockBuilder) Build(_ context.Context, header []byte) (chain.ClientUpdate, error) {
	m.built++
	if len(m.hops) > 0 {
		h := m.hops[0]
		m.hops = m.hops[1:]
		return chain.ClientUpdate{
			Height: h, Payloads: [][]byte{header}, TrustedAt: m.trustedAt,
			IsHop: len(m.hops) > 0, // the last rung is the target, not a hop
		}, nil
	}
	h := uint64(header[0])
	if m.forceHeight != 0 {
		h = m.forceHeight
	}
	if m.noPayload {
		return chain.ClientUpdate{Height: h, TrustedAt: m.trustedAt}, nil
	}
	return chain.ClientUpdate{Height: h, Payloads: [][]byte{header}, TrustedAt: m.trustedAt}, nil
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
		t.Fatalf("isolation must drop only poison, got requeue %v", rq)
	}
	// One fold: the batch. Isolation no longer re-folds, because the update is
	// installed by the attribution step before any packet is judged -- see
	// TestHandleBatch_FoldedUpdateRejectionDoesNotDropInnocentPackets for why
	// that step exists.
	if dst.foldCalls != 1 {
		t.Fatalf("folded batch should be attempted once, got %d fold calls", dst.foldCalls)
	}
	if len(dst.updates) != 1 || dst.updates[0].Height != 100 {
		t.Fatalf("the update must be installed on its own before any packet is judged, got %+v", dst.updates)
	}
	if m.lastHeight != 100 {
		t.Fatalf("the standalone update must advance lastHeight, got %d", m.lastHeight)
	}
	if len(dst.relayed) != 2 ||
		string(dst.relayed[0].Packet) != "good-1" || string(dst.relayed[1].Packet) != "good-2" {
		t.Fatalf("both valid siblings must still relay, got %+v", dst.relayed)
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
			func([]byte, uint64) bool { return true }, // track: no-op for this test
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
			func([]byte, uint64) bool { return true },
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

func TestHandleBatch_RequeuesSendWhenPendingTrackingFails(t *testing.T) {
	src := &mockSource{latest: 100}
	dst := &mockDest{}
	trackCalls := 0
	m := NewModule("test", "client-0", src, dst, &mockBuilder{},
		WithPacketTracker(
			func([]byte, uint64) bool {
				trackCalls++
				return trackCalls > 1
			},
			nil,
		),
	)
	event := chain.Event{Type: chain.SendPacket, Height: 7, Sequence: 9, Raw: []byte("send-1")}

	if got := m.handleBatch(context.Background(), []chain.Event{event}); len(got) != 1 || got[0] != 0 {
		t.Fatalf("failed durable tracking requeue = %v, want [0]", got)
	}
	if dst.relayCalls != 0 {
		t.Fatalf("relay calls after failed durable tracking = %d, want 0", dst.relayCalls)
	}

	if got := m.handleBatch(context.Background(), []chain.Event{event}); len(got) != 0 {
		t.Fatalf("requeue after tracking recovered = %v, want none", got)
	}
	if dst.relayCalls != 1 {
		t.Fatalf("relay calls after tracking recovered = %d, want 1", dst.relayCalls)
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

func (d *transientOnIsolate) RelayPackets(_ context.Context, packets []chain.RelayPacket) (int, error) {
	for _, p := range packets {
		if string(p.Packet) == d.poison {
			return 0, chain.Permanent(errors.New("poison")) // one multicall: nothing lands
		}
	}
	if len(packets) == 1 && string(packets[0].Packet) == d.transientFor {
		return 0, errors.New("rpc blip") // transient
	}
	return len(packets), nil
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

// A permanently-failing packet must leave no wait entry behind, on EVERY terminal
// path. The folded isolated path was the one that did not clear: a batch fold
// reverts, isolation retries each packet carrying the update, and a packet that
// fails permanently there stayed in the tracker until its TTL — so `status` and
// the waiting log would keep reporting an age for a packet that was already dropped.
//
// Asserting on the tracker rather than on a log line is deliberate: the tracker is
// what a future status command reads, and it is what a fifth terminal path would
// also have to clear.
func TestFoldedIsolatedPermanentDropClearsTheWait(t *testing.T) {
	src := &mockSource{relayable: 100}
	dst := &foldingMockDest{enabled: true}
	// The folded batch reverts, the update lands on its own -- so the packets are
	// the poison -- and each isolated singleton then reverts too. Both halves are
	// needed: without the second the packets would simply relay, and this test
	// would assert an empty wait tracker for the wrong reason.
	dst.foldFn = func(_ chain.ClientUpdate, _ []chain.RelayPacket) error {
		return chain.Permanent(errors.New("reverted"))
	}
	dst.relayErr = chain.Permanent(errors.New("reverted"))
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	events := []chain.Event{
		{Type: chain.SendPacket, Height: 10, Sequence: 1, Raw: []byte("a")},
		{Type: chain.SendPacket, Height: 10, Sequence: 2, Raw: []byte("b")},
	}
	// First pass: nothing is relayable yet, so both packets are observed as waiting.
	src.relayable = 5
	if rq := m.handleBatch(context.Background(), events); len(rq) != 2 {
		t.Fatalf("both packets should be re-queued while unrelayable, got %d", len(rq))
	}
	if got := len(m.waits.entries); got != 2 {
		t.Fatalf("expected 2 tracked waits after the first pass, got %d", got)
	}

	// Second pass: relayable, folded batch reverts, isolation drops both permanently.
	src.relayable = 100
	if rq := m.handleBatch(context.Background(), events); len(rq) != 0 {
		t.Fatalf("permanently dropped packets must not be re-queued, got %d", len(rq))
	}
	if got := len(m.waits.entries); got != 0 {
		t.Fatalf("folded isolated drop left %d wait entry/entries behind; every terminal path must clear", got)
	}
	if dst.foldCalls != 1 {
		t.Fatalf("expected one folded batch attempt, got %d", dst.foldCalls)
	}
	if dst.relayCalls != len(events) {
		t.Fatalf("expected one packets-only attempt per packet, got %d for %d packet(s)", dst.relayCalls, len(events))
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
	if !first.logged || first.packets != 1 || first.relayable != 5 || first.description == "" {
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

// TestFoldedUpdateSubmittedWhenNoPacketIsProvable pins the escape from a loop
// the relayer could not leave on its own.
//
// A folding destination leaves the client update unsubmitted so it can ride in
// the packet transaction. When every proof fails there are no packets, and the
// update used to be dropped with them — but the proofs are built at the height
// the client already trusts, so the client standing still is exactly why they
// failed:
//
//	proofs target the trusted height -> that height falls out of the execution
//	node's state window -> every proof fails -> no packets -> the update is
//	dropped -> the client stays put -> the gap only widens
//
// Observed on Sepolia: trustedSlot pinned at 10885728 for twenty minutes while
// finalizedSlot advanced, every pass logging `eth_getProof failed: historical
// state ... is not available`, and no update ever submitted.
func TestFoldedUpdateSubmittedWhenNoPacketIsProvable(t *testing.T) {
	src := &mockSource{latest: 200, failMembershipOn: map[string]bool{"p": true}}
	dst := &foldingMockDest{enabled: true}
	b := &mockBuilder{}
	m := NewModule("test", "client-0", src, dst, b)

	rq := m.handleBatch(context.Background(), []chain.Event{
		{Type: chain.SendPacket, Height: 180, Raw: []byte("p")},
	})

	if len(rq) != 1 {
		t.Fatalf("the unprovable packet must be re-queued, got %v", rq)
	}
	if len(dst.updates) != 1 {
		t.Fatalf("the client update must still be submitted with no packets to fold it into, got %d", len(dst.updates))
	}
	if dst.foldCalls != 0 {
		t.Fatalf("RelayWithUpdate is specified to carry packets; with none it must not be called, got %d calls", dst.foldCalls)
	}
	if m.lastHeight != 200 {
		t.Fatalf("lastHeight must advance to the submitted update height, got %d", m.lastHeight)
	}
}

// The lock the folded plan holds must be released on this path too, or the next
// flush deadlocks — a stall that would look exactly like the bug being fixed.
func TestFoldedUpdateWithNoPacketsReleasesTheLock(t *testing.T) {
	src := &mockSource{latest: 200, failMembershipOn: map[string]bool{"p": true}}
	dst := &foldingMockDest{enabled: true}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	evts := []chain.Event{{Type: chain.SendPacket, Height: 180, Raw: []byte("p")}}
	m.handleBatch(context.Background(), evts)

	done := make(chan struct{})
	go func() { defer close(done); m.handleBatch(context.Background(), evts) }()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("second flush blocked: the folded plan's lock was not released after the packet-less update")
	}
}

// A failed update must not advance lastHeight — otherwise the module would
// believe the client covers heights it does not, and build proofs the
// destination cannot verify.
func TestFoldedUpdateFailureDoesNotAdvanceHeight(t *testing.T) {
	src := &mockSource{latest: 200, failMembershipOn: map[string]bool{"p": true}}
	dst := &foldingMockDest{enabled: true}
	dst.updateErr = errors.New("cosmos rpc down")
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	m.handleBatch(context.Background(), []chain.Event{
		{Type: chain.SendPacket, Height: 180, Raw: []byte("p")},
	})

	if m.lastHeight != 0 {
		t.Fatalf("a failed update must leave lastHeight at 0, got %d", m.lastHeight)
	}
}

// needsRefresh decides whether the anti-expiry routine acts. The direction of the
// comparison is the whole point: plenty of headroom means do nothing, little
// headroom means refresh now. Inverted, the routine refreshes constantly while the
// client is healthy and falls silent exactly as it approaches expiry — and an
// expired light client stops every relay in that direction.
//
// The boundary cases are listed explicitly because a > / >= slip is invisible in
// the middle of the range.
func TestNeedsRefresh(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)

	// With no reported trusting period the module keeps its fixed default. This
	// is the L2 case: no self-expiry, so nothing to take a fraction of.
	t.Run("falls back to the default margin when no period is reported", func(t *testing.T) {
		margin := defaultRefreshMargin
		for _, tt := range []struct {
			name      string
			expiresAt time.Time
			want      bool
		}{
			{"no expiry at all", time.Time{}, false},
			{"expiry far beyond the margin", now.Add(margin + time.Hour), false},
			{"one second more headroom than the margin", now.Add(margin + time.Second), false},
			{"exactly the margin", now.Add(margin), true},
			{"one second inside the margin", now.Add(margin - time.Second), true},
			{"about to expire", now.Add(time.Minute), true},
			{"already expired", now.Add(-time.Hour), true},
		} {
			t.Run(tt.name, func(t *testing.T) {
				if got := needsRefresh(tt.expiresAt, now, 0); got != tt.want {
					t.Fatalf("needsRefresh = %v, want %v (headroom %s, margin %s)",
						got, tt.want, tt.expiresAt.Sub(now), margin)
				}
			})
		}
	})

	// The defect this replaces: a margin fixed at 30 minutes against a trusting
	// period SHORTER than 30 minutes makes the condition always true, so the
	// routine refreshes on every tick forever instead of when the client needs it.
	t.Run("a short trusting period does not make every tick a refresh", func(t *testing.T) {
		const period = 10 * time.Minute
		// Eight minutes of headroom out of a ten-minute period is plenty; under the
		// old fixed 30-minute margin this reported "refresh now".
		if needsRefresh(now.Add(8*time.Minute), now, period) {
			t.Fatalf("refreshed with 8m of headroom on a %s period; the margin must scale with the period", period)
		}
		// One minute of headroom out of ten is not.
		if !needsRefresh(now.Add(time.Minute), now, period) {
			t.Fatalf("did not refresh with 1m of headroom on a %s period", period)
		}
	})

	// And the other direction: a period measured in days must not be refreshed
	// only in its last half hour.
	t.Run("a long trusting period gets a proportionally larger margin", func(t *testing.T) {
		const period = 14 * 24 * time.Hour
		if got := refreshMarginFor(period); got <= defaultRefreshMargin {
			t.Fatalf("margin for a %s period = %s, want more than the %s default", period, got, defaultRefreshMargin)
		}
	})
}

// TestRefreshMarginFor pins the rule itself: the margin is a fraction of the
// client's own trusting period, and it is the SAME rule the refresh-interval
// calculation uses. Two rules for one concept is how the fixed 30 minutes
// survived.
func TestRefreshMarginFor(t *testing.T) {
	for _, tt := range []struct {
		name   string
		period time.Duration
		want   time.Duration
	}{
		{"unknown period falls back", 0, defaultRefreshMargin},
		{"negative period falls back", -time.Hour, defaultRefreshMargin},
		{"a short period gets a quarter of itself", 10 * time.Minute, 150 * time.Second},
		{"a long period is capped at an hour", 14 * 24 * time.Hour, time.Hour},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := refreshMarginFor(tt.period); got != tt.want {
				t.Fatalf("refreshMarginFor(%s) = %s, want %s", tt.period, got, tt.want)
			}
		})
	}
}

// withFastRefreshTick shortens the refresh cadence so a loop test finishes in
// milliseconds instead of a minute, and restores it afterwards.
func withFastRefreshTick(t *testing.T) {
	t.Helper()
	prev := refreshTick
	refreshTick = time.Millisecond
	t.Cleanup(func() { refreshTick = prev })
}

// waitFor polls cond until it holds or ctx expires.
func waitFor(t *testing.T, ctx context.Context, cond func() bool, msg string) {
	t.Helper()
	for {
		if cond() {
			return
		}
		select {
		case <-ctx.Done():
			t.Fatal(msg)
		case <-time.After(time.Millisecond):
		}
	}
}

// updateClientTo is the single place the module advances its view of the
// destination client. Every scenario it must get right lives here.
func TestUpdateClientTo(t *testing.T) {
	// an empty-payload-list ("already current")
	// build advances lastHeight without submitting a tx, so a restart where the
	// client already covers incoming packets does not stall the provability guard.
	t.Run("seeds lastHeight from a no-op build", func(t *testing.T) {
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
	})

	t.Run("is append-only", func(t *testing.T) {
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
	})

	// A failed submission must leave lastHeight where it was. lastHeight is what the
	// module believes the destination client covers; advancing it on a failure makes
	// the module skip the very update that did not land, and the client then sits at
	// an older height than every later decision assumes.
	//
	// This is a named failure mode for this repo — "never advance a
	// timestamp/cursor/tracker on a failed operation" — and the check that catches it
	// is the retry: if lastHeight moved, the second attempt returns early instead of
	// submitting.
	t.Run("does not advance lastHeight on a failed submit", func(t *testing.T) {
		src := &mockSource{}
		dst := &mockDest{updateErr: errors.New("submit reverted")}
		m := NewModule("test", "client-0", src, dst, &mockBuilder{})
		ctx := context.Background()

		if err := m.updateClientTo(ctx, 10); err == nil {
			t.Fatal("a failing UpdateClient must surface as an error")
		}
		if m.lastHeight != 0 {
			t.Fatalf("lastHeight advanced to %d on a failed submit; the client is still at 0", m.lastHeight)
		}

		// The destination recovers: the same height must now actually be submitted.
		dst.updateErr = nil
		if err := m.updateClientTo(ctx, 10); err != nil {
			t.Fatalf("retry after recovery: %v", err)
		}
		if len(dst.updates) != 1 || dst.updates[0].Height != 10 {
			t.Fatalf("retry did not submit the update that previously failed, got %+v", dst.updates)
		}
	})

	// The stale-update guard is inclusive: an update AT the height already trusted is
	// as pointless as one below it, and submitting it spends a proof and a tx to move
	// the client nowhere. The equal case is the one worth pinning — a strictly-less
	// comparison still passes every test that only exercises heights below lastHeight.
	t.Run("skips an update at a height already trusted", func(t *testing.T) {
		src := &mockSource{}
		dst := &mockDest{}
		b := &mockBuilder{}
		m := NewModule("test", "client-0", src, dst, b)
		ctx := context.Background()

		if err := m.updateClientTo(ctx, 10); err != nil {
			t.Fatalf("first update: %v", err)
		}
		if len(dst.updates) != 1 {
			t.Fatalf("want the first update submitted, got %+v", dst.updates)
		}

		// Ask for a higher height (so the pre-Build guard lets it through) but have the
		// builder report exactly the height already trusted.
		b.forceHeight = 10
		if err := m.updateClientTo(ctx, 20); err != nil {
			t.Fatalf("update at trusted height: %v", err)
		}
		if len(dst.updates) != 1 {
			t.Fatalf("submitted an update at the height already trusted: %+v", dst.updates)
		}
		if m.lastHeight != 10 {
			t.Fatalf("lastHeight = %d, want it unchanged at 10", m.lastHeight)
		}
	})
}

// refreshLoop is the anti-expiry routine: it acts only on the decision
// needsRefresh makes, and must keep running through a failed query.
//
// No t.Parallel here, deliberately. withFastRefreshTick reassigns a
// package-level var and restores it with t.Cleanup; running these subtests in
// parallel would turn that sequence into a race that -race catches only
// sometimes.
func TestRefreshLoop(t *testing.T) {
	// The decision above is only useful if the loop acts on it. This covers the wiring:
	// a client inside the margin gets an update submitted without any packet traffic.
	t.Run("advances a client near expiry", func(t *testing.T) {
		withFastRefreshTick(t)

		src := &mockSource{latest: 42}
		dst := &mockDest{expiresAt: time.Now().Add(time.Minute)} // well inside the default margin
		m := NewModule("test", "client-0", src, dst, &mockBuilder{})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { defer wg.Done(); m.refreshLoop(ctx) }()

		waitFor(t, ctx, func() bool {
			m.mu.Lock()
			defer m.mu.Unlock()
			return m.lastHeight == 42
		}, "refresh loop never advanced the client to the source's latest height")
		cancel()
		wg.Wait()
	})

	// A client with plenty of headroom must be left alone: refreshing it early spends a
	// proof and a transaction for nothing, on a cadence of once a minute.
	t.Run("leaves a healthy client alone", func(t *testing.T) {
		withFastRefreshTick(t)

		src := &mockSource{latest: 42}
		dst := &mockDest{expiresAt: time.Now().Add(defaultRefreshMargin + time.Hour)}
		m := NewModule("test", "client-0", src, dst, &mockBuilder{})

		ctx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { defer wg.Done(); m.refreshLoop(ctx) }()
		time.Sleep(50 * time.Millisecond) // many ticks at 1ms
		cancel()
		wg.Wait()

		if len(dst.updates) != 0 {
			t.Fatalf("refreshed a client that was nowhere near expiry: %+v", dst.updates)
		}
	})

	// A failed expiry query must not be read as "no expiry". The loop logs and retries
	// on the next tick; it must not advance any state, and it must not exit — a refresh
	// loop that dies on one bad RPC leaves the client to expire in silence.
	t.Run("survives a failed expiry query", func(t *testing.T) {
		withFastRefreshTick(t)

		src := &mockSource{latest: 42}
		dst := &mockDest{expiresErr: errors.New("rpc down")}
		m := NewModule("test", "client-0", src, dst, &mockBuilder{})

		ctx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup
		wg.Add(1)
		done := make(chan struct{})
		go func() { defer wg.Done(); m.refreshLoop(ctx); close(done) }()

		time.Sleep(50 * time.Millisecond)
		select {
		case <-done:
			t.Fatal("refresh loop exited on a failed expiry query; the client is now unattended")
		default:
		}
		if len(dst.updates) != 0 {
			t.Fatalf("submitted an update off a failed expiry query: %+v", dst.updates)
		}
		m.mu.Lock()
		last := m.lastHeight
		m.mu.Unlock()
		if last != 0 {
			t.Fatalf("lastHeight advanced to %d off a failed expiry query", last)
		}
		cancel()
		wg.Wait()
	})
}

// --- the optional capabilities ---
//
// The With* options are how a path declares what it needs beyond the base relay
// loop. They were all at 0%: nothing pinned that an option actually installs what
// it names, and an option that silently does nothing is invisible -- the module
// runs, relays packets, and quietly never scans for timeouts.

func TestWithTimeoutScanner(t *testing.T) {
	t.Run("installs the scan function and its interval", func(t *testing.T) {
		called := 0
		m := NewModule("test", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
			WithTimeoutScanner(5*time.Second, func(context.Context) { called++ }))

		if m.scan == nil {
			t.Fatal("no scan function installed; Run would skip the timeout loop entirely")
		}
		if m.scanInterval != 5*time.Second {
			t.Fatalf("scanInterval = %s, want 5s", m.scanInterval)
		}
		m.scan(context.Background())
		if called != 1 {
			t.Fatalf("the installed function ran %d times, want 1", called)
		}
	})

	t.Run("a non-positive interval falls back to the default", func(t *testing.T) {
		// Zero is what a caller passes to mean "use the standard cadence" --
		// run_adapters.go does exactly that. Left as zero it becomes
		// time.NewTicker(0), which panics and takes the module down at startup.
		for _, interval := range []time.Duration{0, -time.Second} {
			m := NewModule("test", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
				WithTimeoutScanner(interval, func(context.Context) {}))
			if m.scanInterval != defaultScanInterval {
				t.Errorf("interval %s produced scanInterval %s, want the default %s",
					interval, m.scanInterval, defaultScanInterval)
			}
		}
	})
}

func TestWithClientUpdateObserver(t *testing.T) {
	var seen []time.Time
	trustedAt := time.Unix(1_700_000_000, 0)
	src := &mockSource{}
	dst := &mockDest{}
	b := &mockBuilder{trustedAt: trustedAt}
	m := NewModule("test", "client-0", src, dst, b,
		WithClientUpdateObserver(func(at time.Time) { seen = append(seen, at) }))

	if err := m.updateClientTo(context.Background(), 10); err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(seen) != 1 {
		t.Fatalf("observer called %d times for one client update, want 1", len(seen))
	}
	// The observed value is the height the client now trusts, not the moment the
	// relayer happened to submit it: the freshness metric is about the CLIENT.
	if !seen[0].Equal(trustedAt) {
		t.Fatalf("observer got %s, want the update's TrustedAt %s", seen[0], trustedAt)
	}
}

// A module with no observer must still update. The option is optional, and the
// nil check is what stops every path that does not set one from panicking on its
// first client update.
//
// The builder MUST report a non-zero TrustedAt here. recordClientUpdate is
// guarded by `observer != nil && !TrustedAt.IsZero()`, so a zero timestamp short-
// circuits before the nil check is reached -- an earlier version of this test
// used the default builder and passed with the nil guard deleted.
func TestUpdateClientTo_WorksWithoutAnObserver(t *testing.T) {
	b := &mockBuilder{trustedAt: time.Unix(1_700_000_000, 0)}
	m := NewModule("test", "client-0", &mockSource{}, &mockDest{}, b)
	if err := m.updateClientTo(context.Background(), 10); err != nil {
		t.Fatalf("update without an observer: %v", err)
	}
}

// withFastScanTick shortens the scan cadence so a loop test finishes in
// milliseconds rather than waiting out the 30s production interval.
func withFastScanTick(m *Module) { m.scanInterval = time.Millisecond }

func TestScanLoop(t *testing.T) {
	t.Run("keeps scanning until the context is cancelled", func(t *testing.T) {
		// The timeout sweep is what refunds a packet that was relayed but never
		// delivered. A loop that runs once and stops leaves the pending set
		// growing and the escrow locked, with nothing in the log to say so.
		var mu sync.Mutex
		scans := 0
		m := NewModule("test", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
			WithTimeoutScanner(time.Second, func(context.Context) {
				mu.Lock()
				scans++
				mu.Unlock()
			}))
		withFastScanTick(m)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { defer wg.Done(); m.scanLoop(ctx) }()

		waitFor(t, ctx, func() bool {
			mu.Lock()
			defer mu.Unlock()
			return scans >= 3
		}, "the scan loop did not keep running")
		cancel()
		wg.Wait()
	})

	t.Run("returns immediately on cancellation, not at the next tick", func(t *testing.T) {
		// The interval here is deliberately LONG and the deadline short. Waiting
		// for the next tick would be correct-looking and still wrong: in
		// production the cadence is 30s, and a shutdown that waits it out spends
		// most of the 45s process budget on one worker.
		//
		// A short interval hides this -- an earlier version used a 1ms tick and a
		// 2s deadline, which passed with the ctx.Done() arm deleted entirely.
		const scanInterval = 30 * time.Second
		m := NewModule("test", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
			WithTimeoutScanner(scanInterval, func(context.Context) {}))

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})
		go func() { m.scanLoop(ctx); close(done) }()

		cancel()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatalf("scan loop did not return within 2s of cancellation on a %s cadence; "+
				"shutdown waits for the next tick", scanInterval)
		}
	})

	t.Run("never scans on an already-cancelled context", func(t *testing.T) {
		var mu sync.Mutex
		scans := 0
		m := NewModule("test", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
			WithTimeoutScanner(time.Second, func(context.Context) {
				mu.Lock()
				scans++
				mu.Unlock()
			}))
		withFastScanTick(m)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		m.scanLoop(ctx)

		mu.Lock()
		defer mu.Unlock()
		if scans != 0 {
			t.Fatalf("ran %d scans on an already-cancelled context", scans)
		}
	})

	t.Run("does not start a sweep once cancellation and a tick are both ready", func(t *testing.T) {
		// The case the re-check inside the tick branch exists for. `select` picks
		// at random among ready cases, so when the ticker has already fired and
		// the context is done, half the time it takes the tick -- and without the
		// re-check a sweep begins during shutdown, issuing RPCs nobody is left to
		// read.
		//
		// Forcing that state deterministically is not possible; forcing it REPEATEDLY
		// is. Each round blocks inside the scan long enough for the ticker to fire,
		// cancels while blocked, and then releases. With the re-check the loop
		// returns every time; without it, each round is an independent coin flip,
		// so 20 rounds leave a 1-in-a-million chance of not catching it. The test
		// never fails spuriously -- correct code scans exactly once per round.
		for round := 0; round < 20; round++ {
			var mu sync.Mutex
			scans := 0
			blocked := make(chan struct{})
			release := make(chan struct{})

			ctx, cancel := context.WithCancel(context.Background())
			m := NewModule("test", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
				WithTimeoutScanner(time.Second, func(context.Context) {
					mu.Lock()
					first := scans == 0
					scans++
					mu.Unlock()
					if first {
						close(blocked)
						<-release
					}
				}))
			withFastScanTick(m)

			done := make(chan struct{})
			go func() { m.scanLoop(ctx); close(done) }()

			<-blocked                        // inside the first scan
			time.Sleep(5 * time.Millisecond) // let the 1ms ticker fire while blocked
			cancel()                         // now both cases are ready
			close(release)

			select {
			case <-done:
			case <-time.After(2 * time.Second):
				cancel()
				t.Fatalf("round %d: scan loop did not return", round)
			}

			mu.Lock()
			got := scans
			mu.Unlock()
			if got != 1 {
				t.Fatalf("round %d: ran %d sweeps, want 1; a sweep started after cancellation", round, got)
			}
		}
	})
}

// The ladder is dead code unless handleBatch actually goes through it. This is
// the production shape: an ack whose proof keeps failing, flushed repeatedly.
func TestHandleBatch_EscalatesARepeatedProofFailure(t *testing.T) {
	src := &mockSource{
		latest:           100,
		relayable:        50, // packet at 80 is not provable yet: it waits first
		failMembershipOn: map[string]bool{"ack": true},
	}
	m := NewModule("eth->cosmos", "client-0", src, &mockDest{}, &mockBuilder{})
	clk := &fakeClock{t: time.Unix(1_700_000_000, 0)}
	m.waits.now = clk.now

	batch := []chain.Event{{Type: chain.AckPacket, Sequence: 21, Height: 80, Raw: []byte("ack"), AckBytes: [][]byte{{1}}}}

	// Sixteen minutes parked at the source gate before a proof is ever attempted.
	// This is the normal case on an L2 return leg, and none of it is a proof
	// failure -- so it must not be charged to the ladder below.
	m.handleBatch(context.Background(), batch)
	clk.t = clk.t.Add(16 * time.Minute)
	src.relayable = 100

	first := captureLog(func() { m.handleBatch(context.Background(), batch) })
	if !strings.Contains(first, "ack seq=21") {
		t.Fatalf("the first proof failure must reach the log:\n%s", first)
	}
	if strings.Contains(first, "STUCK") {
		t.Fatalf("the FIRST proof failure was reported as STUCK because it inherited the pre-proof wait:\n%s", first)
	}

	// Ten more flushes within the same threshold: silent. This is the case that
	// produced 296 byte-identical lines.
	repeats := captureLog(func() {
		for i := 0; i < 10; i++ {
			clk.t = clk.t.Add(2 * time.Second)
			m.handleBatch(context.Background(), batch)
		}
	})
	if strings.Contains(repeats, "ack seq=21") {
		t.Fatalf("repeats inside one threshold must be suppressed:\n%s", repeats)
	}

	// Past the stuck threshold the wording has to change: an operator reading
	// "proof for packet" cannot tell a retry from a packet with no way out.
	clk.t = clk.t.Add(proofFailureStuckAfter + time.Minute)
	stuck := captureLog(func() { m.handleBatch(context.Background(), batch) })
	for _, want := range []string{"STUCK", "ack seq=21", "no other exit"} {
		if !strings.Contains(stuck, want) {
			t.Fatalf("the escalated line must contain %q:\n%s", want, stuck)
		}
	}
}

// TestHandleBatch_StuckSendNamesTheTimeoutExit: "no other exit" is true of an ack
// and false of a send. Every source reports an expired send as permanent, the
// module drops it, and the timeout scanner refunds it -- so telling an operator a
// send is unrecoverable sends them looking for a problem that resolves itself.
// The two packet types need different next actions, so they need different lines.
func TestHandleBatch_StuckSendNamesTheTimeoutExit(t *testing.T) {
	src := &mockSource{
		latest:           100,
		relayable:        100,
		failMembershipOn: map[string]bool{"snd": true},
	}
	m := NewModule("eth->cosmos", "client-0", src, &mockDest{}, &mockBuilder{})
	clk := &fakeClock{t: time.Unix(1_700_000_000, 0)}
	m.waits.now = clk.now

	batch := []chain.Event{{Type: chain.SendPacket, Sequence: 21, Height: 80, Raw: []byte("snd")}}

	m.handleBatch(context.Background(), batch)
	clk.t = clk.t.Add(proofFailureStuckAfter + time.Minute)
	stuck := captureLog(func() { m.handleBatch(context.Background(), batch) })

	// A SendPacket prints as "recv": the line names the message the destination
	// will be asked for, not the event's origin (see TestEventTypeNames).
	for _, want := range []string{"STUCK", "recv seq=21", "timeout scanner"} {
		if !strings.Contains(stuck, want) {
			t.Fatalf("a stuck send must name the exit it actually has (%q):\n%s", want, stuck)
		}
	}
	if strings.Contains(stuck, "no other exit") {
		t.Fatalf("a send is refunded once past its timeout, so it does have another exit:\n%s", stuck)
	}
}

type drainingSource struct {
	started chan struct{}
	stopped chan struct{}
}

func (s *drainingSource) Chain() chain.ChainType { return chain.Cosmos }
func (s *drainingSource) Subscribe(ctx context.Context, _ func(context.Context, []chain.Event) []int) error {
	close(s.started)
	<-ctx.Done()
	close(s.stopped)
	return ctx.Err()
}
func (s *drainingSource) LatestHeight(context.Context) (uint64, error)        { return 0, nil }
func (s *drainingSource) RelayableHeight(context.Context) (uint64, error)     { return 0, nil }
func (s *drainingSource) QueryHeader(context.Context, uint64) ([]byte, error) { return nil, nil }
func (s *drainingSource) MembershipProof(context.Context, []byte, uint64, chain.EventType) ([]byte, error) {
	return nil, nil
}
func (s *drainingSource) NonMembershipProof(context.Context, []byte, uint64) ([]byte, error) {
	return nil, nil
}

func TestModuleCleanCancellationDrainsSubscriber(t *testing.T) {
	src := &drainingSource{started: make(chan struct{}), stopped: make(chan struct{})}
	m := NewModule("test", "client", src, &mockDest{}, &mockBuilder{})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- m.Run(ctx) }()
	<-src.started
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run cancellation error = %v, want nil", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not drain after cancellation")
	}
	select {
	case <-src.stopped:
	default:
		t.Fatal("Run returned before subscriber stopped")
	}
}

type internallyCancelledSource struct{}

func (internallyCancelledSource) Chain() chain.ChainType { return chain.Cosmos }
func (internallyCancelledSource) Subscribe(context.Context, func(context.Context, []chain.Event) []int) error {
	return context.Canceled
}
func (internallyCancelledSource) LatestHeight(context.Context) (uint64, error) { return 0, nil }
func (internallyCancelledSource) RelayableHeight(context.Context) (uint64, error) {
	return 0, nil
}
func (internallyCancelledSource) QueryHeader(context.Context, uint64) ([]byte, error) {
	return nil, nil
}
func (internallyCancelledSource) MembershipProof(context.Context, []byte, uint64, chain.EventType) ([]byte, error) {
	return nil, nil
}
func (internallyCancelledSource) NonMembershipProof(context.Context, []byte, uint64) ([]byte, error) {
	return nil, nil
}

func TestModuleReportsInternalSubscriberCancellation(t *testing.T) {
	m := NewModule("test", "client", internallyCancelledSource{}, &mockDest{}, &mockBuilder{})
	err := m.Run(context.Background())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run error = %v, want wrapped context.Canceled", err)
	}
}

// stuckWorkerSource is what the real Cosmos and EVM adapters became: on
// cancellation they drain their own goroutines, and if one of them never
// returns they report THAT rather than the cancellation that triggered it.
type stuckWorkerSource struct {
	started chan struct{}
	err     error
}

func (s *stuckWorkerSource) Chain() chain.ChainType { return chain.Cosmos }
func (s *stuckWorkerSource) Subscribe(ctx context.Context, _ func(context.Context, []chain.Event) []int) error {
	close(s.started)
	<-ctx.Done()
	return s.err
}
func (s *stuckWorkerSource) LatestHeight(context.Context) (uint64, error)        { return 0, nil }
func (s *stuckWorkerSource) RelayableHeight(context.Context) (uint64, error)     { return 0, nil }
func (s *stuckWorkerSource) QueryHeader(context.Context, uint64) ([]byte, error) { return nil, nil }
func (s *stuckWorkerSource) MembershipProof(context.Context, []byte, uint64, chain.EventType) ([]byte, error) {
	return nil, nil
}
func (s *stuckWorkerSource) NonMembershipProof(context.Context, []byte, uint64) ([]byte, error) {
	return nil, nil
}

// Reported by @DongLieu on #434. Two things had to line up for a stuck worker to
// be reported as a clean shutdown, and both were true:
//
//  1. Run's select takes ctx.Done() on SIGTERM, sets runErr = context.Canceled,
//     and NEVER reads subscribeErr -- so whatever Subscribe returned was dropped.
//  2. A bare context.Canceled during shutdown is normalised to nil.
//
// So a source that knew one of its own goroutines had not stopped had no way to
// say so: the value it returned was discarded, and the value that replaced it
// meant "stopped cleanly". The process exited 0 with a goroutine still running.
func TestRunReportsASourceThatCouldNotStopItsOwnWorkers(t *testing.T) {
	stuck := errors.New("cosmos source: shutdown drain timed out after 20s; workers still running: [cosmos-subscribe]")
	src := &stuckWorkerSource{started: make(chan struct{}), err: stuck}
	m := NewModule("test", "client", src, &mockDest{}, &mockBuilder{})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- m.Run(ctx) }()
	<-src.started
	cancel()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Run reported a clean shutdown while a named source worker was still running; " +
				"SIGTERM would exit 0 with the goroutine alive")
		}
		if !errors.Is(err, stuck) {
			t.Fatalf("Run error = %v, want it to carry the source's drain failure", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not return after cancellation")
	}
}

// The mirror of TestModuleCleanCancellationDrainsSubscriber: preferring the
// source's error must not turn an ordinary stop into a failure. A source that
// returns the cancellation it was given is still a clean shutdown.
func TestRunStillExitsCleanWhenTheSourceOnlyReportsCancellation(t *testing.T) {
	src := &stuckWorkerSource{started: make(chan struct{}), err: context.Canceled}
	m := NewModule("test", "client", src, &mockDest{}, &mockBuilder{})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- m.Run(ctx) }()
	<-src.started
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run error = %v, want nil: an ordinary stop is not a failure", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not return after cancellation")
	}
}

// A permanently-failing source -- a misconfigured attestor route, a malformed
// request, an RPC the daemon does not implement -- used to be asked again on
// every pass. The L2 loop re-offers its queued packets every 4s and deliberately
// does so even when the chain produced no block, so one bad config line became
// roughly 21,600 attestor calls and log lines a day, each looking like an
// ordinary transient RPC failure.
//
// The hold changes only how often the source is asked. It must not drop a packet
// and must not skip the durable tracking that pays for the timeout refund.
func TestPermanentSourceFailureIsHeldOff(t *testing.T) {
	events := []chain.Event{{Type: chain.SendPacket, Sequence: 7, Height: 40, Raw: []byte("pkt-7")}}

	t.Run("the source is not asked again during the hold", func(t *testing.T) {
		src := &mockSource{latest: 100, relayableErr: chain.Permanent(errors.New("attestor does not serve this chain"))}
		m := NewModule("l2->cosmos", "client-0", src, &mockDest{}, &mockBuilder{})

		out := captureLog(func() { m.handleBatch(context.Background(), events) })
		if src.relayableCalls != 1 {
			t.Fatalf("first pass asked the source %d times, want 1", src.relayableCalls)
		}
		if !strings.Contains(out, "PERMANENT") {
			t.Fatalf("a permanent source failure must say so, not read as a transient RPC error:\n%s", out)
		}

		for i := 0; i < 5; i++ {
			m.handleBatch(context.Background(), events)
		}
		if src.relayableCalls != 1 {
			t.Fatalf("source asked %d times across 6 passes; the hold is not holding", src.relayableCalls)
		}
	})

	// The hold skips the RPC and the proof, never the durable record: a packet
	// queued while the source is held off must still be tracked, or it can never
	// be refunded when it expires.
	t.Run("packets are still tracked while held off", func(t *testing.T) {
		src := &mockSource{latest: 100, relayableErr: chain.Permanent(errors.New("attestor does not serve this chain"))}
		var tracked [][]byte
		m := NewModule("l2->cosmos", "client-0", src, &mockDest{}, &mockBuilder{},
			WithPacketTracker(func(raw []byte, _ uint64) bool {
				tracked = append(tracked, raw)
				return true
			}, nil))

		m.handleBatch(context.Background(), events) // arms the hold
		m.handleBatch(context.Background(), events) // held off

		if len(tracked) != 2 {
			t.Fatalf("tracked %d packets across two passes, want 2: a packet queued during the hold "+
				"that is never tracked can never be refunded", len(tracked))
		}
	})

	t.Run("the hold clears once the source answers", func(t *testing.T) {
		src := &mockSource{latest: 100, relayableErr: chain.Permanent(errors.New("attestor does not serve this chain"))}
		m := NewModule("l2->cosmos", "client-0", src, &mockDest{}, &mockBuilder{})

		m.handleBatch(context.Background(), events)
		m.sourceProbeAt = time.Now().Add(-time.Second) // the operator fixed the config
		src.relayableErr = nil

		out := captureLog(func() { m.handleBatch(context.Background(), events) })
		if src.relayableCalls != 2 {
			t.Fatalf("source asked %d times; the probe after the hold expired did not happen", src.relayableCalls)
		}
		if !strings.Contains(out, "permanent-failure hold cleared") {
			t.Fatalf("a recovered source must say so:\n%s", out)
		}
		if !m.sourceProbeAt.IsZero() || m.sourceBackoff != 0 {
			t.Fatalf("hold state survived recovery: probeAt=%v backoff=%v", m.sourceProbeAt, m.sourceBackoff)
		}
	})

	// A transient failure must not arm the hold, or a source that is merely
	// unreachable for a moment goes quiet for up to fifteen minutes.
	t.Run("a transient failure does not arm the hold", func(t *testing.T) {
		src := &mockSource{latest: 100, relayableErr: chain.Transient(errors.New("attestor unavailable"))}
		m := NewModule("l2->cosmos", "client-0", src, &mockDest{}, &mockBuilder{})

		m.handleBatch(context.Background(), events)
		m.handleBatch(context.Background(), events)

		if src.relayableCalls != 2 {
			t.Fatalf("source asked %d times across two passes, want 2: a transient failure was treated as permanent", src.relayableCalls)
		}
	})
}

// Validator churn can outrun the trust-overlap requirement: an update from the
// trusted height straight to the head then cannot be verified at all, and the
// only way across is a chain of updates that each overlap the last.
//
// That mechanism lived only in the manual `update-client` command, so the
// AUTOMATIC path could not cross a churn gap: the client stayed stuck until a
// person typed a command, while refreshLoop retried and failed underneath.
func TestUpdateClientTo_ClimbsThroughIntermediateHops(t *testing.T) {
	t.Run("crosses a churn gap without a manual command", func(t *testing.T) {
		src := &mockSource{latest: 100}
		dst := &mockDest{}
		b := &mockBuilder{hops: []uint64{30, 60, 90}} // two hops, then the target
		m := NewModule("test", "client-0", src, dst, b)

		if err := m.updateClientTo(context.Background(), 90); err != nil {
			t.Fatalf("climb: %v", err)
		}
		if len(dst.updates) != 3 {
			t.Fatalf("submitted %d updates, want the two hops and the target", len(dst.updates))
		}
		for i, want := range []uint64{30, 60, 90} {
			if dst.updates[i].Height != want {
				t.Fatalf("update %d at height %d, want %d: each hop must overlap the last",
					i, dst.updates[i].Height, want)
			}
		}
		if m.lastHeight != 90 {
			t.Fatalf("lastHeight = %d, want 90", m.lastHeight)
		}
	})

	// An adapter that never hops must behave exactly as before: one Build, one
	// update, no loop.
	t.Run("an adapter that does not hop is unchanged", func(t *testing.T) {
		src := &mockSource{latest: 100}
		dst := &mockDest{}
		b := &mockBuilder{}
		m := NewModule("test", "client-0", src, dst, b)

		if err := m.updateClientTo(context.Background(), 42); err != nil {
			t.Fatalf("update: %v", err)
		}
		if b.built != 1 || len(dst.updates) != 1 {
			t.Fatalf("built %d, submitted %d; want exactly one of each", b.built, len(dst.updates))
		}
	})

	// Each hop costs a proof and a transaction, so an unbounded climb is an
	// expensive silent loop -- the third failure class dressed as progress.
	t.Run("stops at the ceiling and says how far it got", func(t *testing.T) {
		endless := make([]uint64, MaxClientUpdateHops+5)
		for i := range endless {
			endless[i] = uint64(i + 1)
		}
		src := &mockSource{latest: 1000}
		dst := &mockDest{}
		b := &mockBuilder{hops: endless}
		m := NewModule("test", "client-0", src, dst, b)

		err := m.updateClientTo(context.Background(), 1000)
		if err == nil {
			t.Fatal("an endless climb reported success")
		}
		if !strings.Contains(err.Error(), "hop") {
			t.Errorf("error = %v; it must say the climb ran out, not look like a build failure", err)
		}
		if len(dst.updates) > MaxClientUpdateHops {
			t.Fatalf("submitted %d updates past a ceiling of %d", len(dst.updates), MaxClientUpdateHops)
		}
	})

	// A hop that installs nothing must end the climb rather than repeat it. This
	// is the anti-spin guard, and without it a builder stuck on one height would
	// burn the whole ceiling in proofs.
	t.Run("a hop that does not advance ends the climb", func(t *testing.T) {
		src := &mockSource{latest: 100}
		dst := &mockDest{}
		b := &mockBuilder{hops: []uint64{50, 50, 50, 50}}
		m := NewModule("test", "client-0", src, dst, b)

		if err := m.updateClientTo(context.Background(), 90); err != nil {
			t.Fatalf("climb: %v", err)
		}
		if len(dst.updates) != 1 {
			t.Fatalf("submitted %d updates for a builder stuck at one height, want 1", len(dst.updates))
		}
	})
}

func TestAckWatch(t *testing.T) {
	send := chain.RelayPacket{Type: chain.SendPacket, Sequence: 7, Packet: []byte("pkt-7")}
	ack := chain.RelayPacket{Type: chain.AckPacket, Sequence: 7, Packet: []byte("pkt-7")}

	newModule := func(due AckDueFunc, settled AckSettledFunc, overdue OverdueAcksFunc) *Module {
		return NewModule("cosmos->eth", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
			WithAckWatch(0, 0, due, settled, overdue))
	}

	t.Run("a delivered receive makes the acknowledgement owed", func(t *testing.T) {
		var recorded [][]byte
		m := newModule(func(raw []byte) bool { recorded = append(recorded, raw); return true }, nil, func(time.Duration) []string { return nil })

		m.settleDelivered([]chain.RelayPacket{send})

		if len(recorded) != 1 || string(recorded[0]) != "pkt-7" {
			t.Fatalf("recorded %v, want the delivered packet: the debt starts when WE deliver, not when we see the send", recorded)
		}
	})

	t.Run("a relayed acknowledgement closes the debt", func(t *testing.T) {
		var cleared [][]byte
		m := newModule(nil, func(raw []byte) { cleared = append(cleared, raw) }, func(time.Duration) []string { return nil })

		m.settleDelivered([]chain.RelayPacket{ack})

		if len(cleared) != 1 || string(cleared[0]) != "pkt-7" {
			t.Fatalf("cleared %v, want the acknowledged packet", cleared)
		}
	})

	// A durable write that fails is the case where the record would exist only in
	// memory -- and a restart is the most likely reason the ack was missed in the
	// first place, so failing quietly here loses exactly the packet this watches.
	t.Run("a failed record is reported", func(t *testing.T) {
		m := newModule(func([]byte) bool { return false }, nil, func(time.Duration) []string { return nil })

		out := captureLog(func() { m.settleDelivered([]chain.RelayPacket{send}) })

		if !strings.Contains(out, "could not durably record the owed acknowledgement for seq=7") {
			t.Fatalf("a failed durable write must reach the log:\n%s", out)
		}
	})

	// The retry itself is NOT tested here any more: it moved into
	// services.ackDueLedger, because the module that records a debt is not the one
	// that settles it, so a per-module queue could never be cancelled by the
	// settlement. See services/ackdue_test.go -- "a settlement cancels a record
	// still queued for writing".

	// The second module gets the hooks with a nil overdue reporter: the tracker
	// is shared, so a second watch loop would report every overdue ack twice.
	// nil must leave the loop unstarted while recording and settlement still work
	// -- that combination is what lets the debt be recorded on one relay
	// direction and cleared on the other.
	t.Run("a nil overdue reporter still wires recording and settlement", func(t *testing.T) {
		var recorded, cleared [][]byte
		m := newModule(
			func(raw []byte) bool { recorded = append(recorded, raw); return true },
			func(raw []byte) { cleared = append(cleared, raw) },
			nil,
		)
		if m.overdueAcks != nil {
			t.Fatal("a nil reporter must stay nil; Module.Run starts the watch loop on it")
		}

		m.settleDelivered([]chain.RelayPacket{send})
		m.settleDelivered([]chain.RelayPacket{ack})

		if len(recorded) != 1 || len(cleared) != 1 {
			t.Fatalf("recorded %d and cleared %d, want 1 each without a watch loop", len(recorded), len(cleared))
		}
	})

	t.Run("the watch names overdue acknowledgements and stays quiet otherwise", func(t *testing.T) {
		var threshold time.Duration
		overdue := []string{"seq=7 src=08-wasm-0"}
		m := newModule(nil, nil, func(d time.Duration) []string { threshold = d; return overdue })
		m.ackWatchInterval = time.Millisecond

		ctx, cancel := context.WithCancel(context.Background())
		out := captureLog(func() {
			done := make(chan struct{})
			go func() { defer close(done); m.ackWatchLoop(ctx) }()
			time.Sleep(20 * time.Millisecond)
			cancel()
			<-done
		})

		if !strings.Contains(out, "seq=7 src=08-wasm-0") {
			t.Fatalf("an overdue acknowledgement must be named:\n%s", out)
		}
		if threshold != defaultAckOverdueAfter {
			t.Fatalf("watch asked for packets older than %s, want the default %s", threshold, defaultAckOverdueAfter)
		}

		overdue = nil
		quiet := captureLog(func() {
			ctx2, cancel2 := context.WithCancel(context.Background())
			done := make(chan struct{})
			go func() { defer close(done); m.ackWatchLoop(ctx2) }()
			time.Sleep(20 * time.Millisecond)
			cancel2()
			<-done
		})
		if strings.Contains(quiet, "acknowledgement(s) owed") {
			t.Fatalf("nothing overdue must produce no line:\n%s", quiet)
		}
	})
}

// A destination can commit part of a batch and then fail: a Cosmos batch that
// splits submits its leading sub-batches, and the node can go away between two
// of them. Those packets are on chain.
//
// Before the committed count reached the engine, both endings were wrong for
// them. The transient ending re-queued the whole batch, so the delivered packets
// were submitted a second time — and then reverted as duplicates, which the
// engine reads as permanent and reports as DROP for packets that were in fact
// delivered. The permanent ending isolated all of them, one wasted transaction
// each, with the same misleading DROP.
func TestHandleBatch_PartialCommit(t *testing.T) {
	batch := []chain.Event{
		{Type: chain.SendPacket, Height: 7, Raw: []byte("send-1")},
		{Type: chain.SendPacket, Height: 8, Raw: []byte("send-2")},
		{Type: chain.SendPacket, Height: 9, Raw: []byte("send-3")},
	}

	run := func(t *testing.T, dst *mockDest) (requeue []int, untracked [][]byte) {
		t.Helper()
		m := NewModule("test", "client-0", &mockSource{latest: 100}, dst, &mockBuilder{},
			WithPacketTracker(
				func([]byte, uint64) bool { return true },
				func(pkt []byte) { untracked = append(untracked, pkt) },
			),
		)
		return m.handleBatch(context.Background(), batch), untracked
	}

	t.Run("a transient failure re-queues only the uncommitted tail", func(t *testing.T) {
		dst := &mockDest{relayErr: errors.New("node went away"), relayCommit: 2}
		requeue, untracked := run(t, dst)

		if len(requeue) != 1 || requeue[0] != 2 {
			t.Fatalf("requeue = %v, want only index 2; the first two are already on chain", requeue)
		}
		if len(untracked) != 2 {
			t.Fatalf("untracked %d packet(s) (%q), want the 2 that were committed", len(untracked), untracked)
		}
		if dst.relayCalls != 1 {
			t.Fatalf("relay calls = %d, want 1; a transient failure must not start isolating", dst.relayCalls)
		}
	})

	t.Run("a permanent failure isolates only the uncommitted tail", func(t *testing.T) {
		dst := &mockDest{relayErr: chain.Permanent(errors.New("reverted")), relayCommit: 2}
		requeue, untracked := run(t, dst)

		// The tail is a single packet, so it is dropped outright rather than
		// isolated — and the committed two must not be dropped with it.
		if len(requeue) != 0 {
			t.Fatalf("requeue = %v, want none: two landed and the third is permanently dead", requeue)
		}
		if len(untracked) != 2 {
			t.Fatalf("untracked %d packet(s) (%q), want the 2 that were committed", len(untracked), untracked)
		}
		if dst.relayCalls != 1 {
			t.Fatalf("relay calls = %d, want 1; the committed prefix must not be re-submitted for isolation", dst.relayCalls)
		}
	})

	t.Run("a batch that fully committed before the error re-queues nothing", func(t *testing.T) {
		dst := &mockDest{relayErr: errors.New("receipt wait timed out"), relayCommit: 3}
		requeue, untracked := run(t, dst)

		if len(requeue) != 0 {
			t.Fatalf("requeue = %v, want none: every packet is on chain", requeue)
		}
		if len(untracked) != 3 {
			t.Fatalf("untracked %d packet(s), want all 3", len(untracked))
		}
	})

	t.Run("a count larger than the batch is clamped, not indexed with", func(t *testing.T) {
		// An adapter breaking the contract must not take the relay loop down with
		// a slice bounds panic.
		dst := &mockDest{relayErr: errors.New("node went away"), relayCommit: 99}
		requeue, untracked := run(t, dst)

		if len(requeue) != 0 {
			t.Fatalf("requeue = %v, want none after clamping", requeue)
		}
		if len(untracked) != 3 {
			t.Fatalf("untracked %d packet(s), want at most the 3 in the batch", len(untracked))
		}
	})
}

func TestFlushOnce(t *testing.T) {
	send := chain.Event{Type: chain.SendPacket, Sequence: 7, Height: 40, Raw: []byte("old-packet")}

	t.Run("relays a packet the scan never saw", func(t *testing.T) {
		src := &listerSource{mockSource: mockSource{latest: 100, relayable: 100}, candidate: []chain.Event{send}}
		dst := &mockDest{}
		m := NewModule("cosmos->eth", "client-0", src, dst, &mockBuilder{})

		out := captureLog(func() { m.flushOnce(context.Background(), src) })

		if len(dst.relayed) != 1 {
			t.Fatalf("relayed %d packets, want 1: the flush exists to relay what the scan cannot reach", len(dst.relayed))
		}
		if !strings.Contains(out, "packet flush found 1 packet(s)") {
			t.Fatalf("a flush that finds work must say so:\n%s", out)
		}
	})

	// An outstanding source commitment only says the packet was never acked BACK.
	// A delivered packet whose ack is in flight still has one, and re-proving it
	// every pass for as long as the ack takes is pure waste.
	t.Run("skips packets the destination already holds a receipt for", func(t *testing.T) {
		src := &listerSource{mockSource: mockSource{latest: 100, relayable: 100}, candidate: []chain.Event{send}}
		dst := &mockDest{hasReceipt: true}
		m := NewModule("cosmos->eth", "client-0", src, dst, &mockBuilder{})

		out := captureLog(func() { m.flushOnce(context.Background(), src) })

		if len(dst.relayed) != 0 {
			t.Fatalf("relayed %d packets, want 0: the destination already has a receipt", len(dst.relayed))
		}
		if strings.Contains(out, "packet flush found") {
			t.Fatalf("a pass that finds nothing unsettled must stay quiet:\n%s", out)
		}
	})

	// The receipt check has to be serialized with submission. It used to run
	// entirely outside batchMu, so a live subscription batch could relay the
	// packet between the check and the lock -- and the flush then submitted a
	// duplicate receive: an on-chain revert, wasted gas, and a permanent failure
	// that takes any other packet folded into the same batch down with it.
	t.Run("re-checks receipts under the batch lock", func(t *testing.T) {
		src := &listerSource{mockSource: mockSource{latest: 100, relayable: 100}, candidate: []chain.Event{send}}
		// Outstanding at the first check; delivered by the time the lock is held.
		dst := &mockDest{receiptAnswers: []bool{false, true}}
		m := NewModule("cosmos->eth", "client-0", src, dst, &mockBuilder{})

		out := captureLog(func() { m.flushOnce(context.Background(), src) })

		if len(dst.relayed) != 0 {
			t.Fatalf("relayed %d packet(s), want 0: it was delivered between the two checks, so this is a duplicate receive", len(dst.relayed))
		}
		if strings.Contains(out, "packet flush found") {
			t.Fatalf("a pass whose packets were relayed underneath it must stay quiet:\n%s", out)
		}
	})

	// Dropping on a failed receipt query would silently skip exactly the packet
	// this loop exists to find, and the next pass would ask the same broken
	// endpoint again.
	t.Run("a failed receipt check relays rather than skips", func(t *testing.T) {
		src := &listerSource{mockSource: mockSource{latest: 100, relayable: 100}, candidate: []chain.Event{send}}
		dst := &mockDest{receiptErr: errors.New("rpc down")}
		m := NewModule("cosmos->eth", "client-0", src, dst, &mockBuilder{})

		captureLog(func() { m.flushOnce(context.Background(), src) })

		if len(dst.relayed) != 1 {
			t.Fatalf("relayed %d packets, want 1: an unknown receipt must not be read as delivered", len(dst.relayed))
		}
	})

	t.Run("an enumeration failure is logged and relays nothing", func(t *testing.T) {
		src := &listerSource{mockSource: mockSource{latest: 100, relayable: 100}, listErr: errors.New("query failed")}
		dst := &mockDest{}
		m := NewModule("cosmos->eth", "client-0", src, dst, &mockBuilder{})

		out := captureLog(func() { m.flushOnce(context.Background(), src) })

		if len(dst.relayed) != 0 {
			t.Fatalf("relayed %d packets on a failed query, want 0", len(dst.relayed))
		}
		if !strings.Contains(out, "packet flush:") {
			t.Fatalf("a failed enumeration must reach the log:\n%s", out)
		}
	})
}

// The restart boundary the ledger exists to close, and the one it did not until
// a review caught the order.
//
// A delivered receive moves from one durable obligation to another: the pending
// record that would time it out, and the owed-ack record that says an
// acknowledgement is expected. Between those two there must never be an instant
// where the packet is in neither -- the receive is committed on the destination,
// so it can no longer time out, and with no owed-ack record nothing durable is
// left from which the missing acknowledgement could be noticed.
func TestSettleDelivered_NeverLeavesAPacketInNeitherLedger(t *testing.T) {
	send := chain.RelayPacket{Type: chain.SendPacket, Sequence: 7, Packet: []byte("pkt-7")}

	// order records which durable obligation was touched, and when. The sequence
	// is the property: "recorded" must precede "untracked", because a crash
	// between them must leave the packet in the FIRST ledger, not in none.
	newModule := func(t *testing.T, durable bool, order *[]string) *Module {
		t.Helper()
		return NewModule("cosmos->eth", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
			WithPacketTracker(
				func([]byte, uint64) bool { return true },
				func([]byte) { *order = append(*order, "untracked") },
			),
			WithAckWatch(0, 0,
				func([]byte) bool { *order = append(*order, "recorded"); return durable },
				nil, nil),
		)
	}

	t.Run("records the owed acknowledgement before releasing the pending record", func(t *testing.T) {
		var order []string
		m := newModule(t, true, &order)

		m.settleDelivered([]chain.RelayPacket{send})

		want := []string{"recorded", "untracked"}
		if len(order) != 2 || order[0] != want[0] || order[1] != want[1] {
			t.Fatalf("order = %v, want %v: a crash between the two must land in the owed-ack ledger, not in neither",
				order, want)
		}
	})

	// The write failing is the case the reviewer named, and it is the easier half:
	// the ORDER covers the crash, this covers the error return.
	t.Run("keeps the pending record when the owed-ack write does not reach disk", func(t *testing.T) {
		var order []string
		m := newModule(t, false, &order)

		m.settleDelivered([]chain.RelayPacket{send})

		if len(order) != 1 || order[0] != "recorded" {
			t.Fatalf("order = %v; the pending record was released even though the owed-ack write failed, "+
				"leaving the packet in neither durable ledger", order)
		}
	})

	// A path with no return leg has no second obligation to hand the packet to,
	// so holding the pending record forever would strand it instead of protecting
	// it. Per-destination gating leaves ackDue nil there.
	t.Run("releases the pending record when there is no acknowledgement ledger", func(t *testing.T) {
		var order []string
		m := NewModule("cosmos->l2", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
			WithPacketTracker(
				func([]byte, uint64) bool { return true },
				func([]byte) { order = append(order, "untracked") },
			),
		)

		m.settleDelivered([]chain.RelayPacket{send})

		if len(order) != 1 || order[0] != "untracked" {
			t.Fatalf("order = %v, want the pending record released: with no ledger to hand it to, keeping it strands the packet",
				order)
		}
	})
}

// The churn-gap climb has to exist on BOTH sides of prepareBatchUpdate. The
// unfolded branch delegates to updateClientTo, which loops; the folded branch
// built exactly one update and returned it, so a client behind a churn gap
// advanced one hop per flush. With a 30s scan interval and a ceiling of 16 that
// is eight minutes of re-queueing to cross a gap the unfolded path crosses in
// one call -- on the packet-driven path, which is the one with packets waiting.
func TestHandleBatch_FoldedPathClimbsAChurnGapInOneFlush(t *testing.T) {
	t.Run("crosses the gap in one flush and folds only the last rung", func(t *testing.T) {
		src := &mockSource{latest: 100}
		dst := &foldingMockDest{enabled: true}
		m := NewModule("test", "client-0", src, dst, &mockBuilder{hops: []uint64{30, 60, 90}})
		m.lastHeight = 10

		batch := []chain.Event{{Type: chain.SendPacket, Height: 90, Raw: []byte("a")}}
		if rq := m.handleBatch(context.Background(), batch); len(rq) != 0 {
			t.Fatalf("a climbable gap must not re-queue the batch, got %v", rq)
		}

		// The two intermediate hops cannot carry the packets -- each lands below
		// the packet height -- so they go out on their own. Only the rung that
		// reaches the target is folded.
		if len(dst.updates) != 2 || dst.updates[0].Height != 30 || dst.updates[1].Height != 60 {
			t.Fatalf("standalone updates = %+v, want the two intermediate hops at 30 and 60", dst.updates)
		}
		if dst.foldCalls != 1 || len(dst.foldedUpdates) != 1 || dst.foldedUpdates[0].Height != 90 {
			t.Fatalf("want one folded update at the target height 90, calls=%d updates=%+v", dst.foldCalls, dst.foldedUpdates)
		}
		if len(dst.foldedPackets) != 1 || len(dst.foldedPackets[0]) != 1 {
			t.Fatalf("the packet must ride the final tx, got %+v", dst.foldedPackets)
		}
		if m.lastHeight != 90 {
			t.Fatalf("lastHeight = %d, want 90", m.lastHeight)
		}
	})

	// The ceiling is the whole reason the climb is safe to run on the packet path:
	// each hop is a proof and a transaction, so an unbounded one is a gas drain
	// that looks like progress.
	t.Run("stops at the ceiling instead of climbing forever", func(t *testing.T) {
		endless := make([]uint64, chain.MaxClientUpdateHops+5)
		for i := range endless {
			endless[i] = uint64(i + 1)
		}
		src := &mockSource{latest: 1000, relayable: 1000}
		dst := &foldingMockDest{enabled: true}
		m := NewModule("test", "client-0", src, dst, &mockBuilder{hops: endless})

		batch := []chain.Event{{Type: chain.SendPacket, Height: 900, Raw: []byte("a")}}
		if rq := m.handleBatch(context.Background(), batch); len(rq) != 1 {
			t.Fatalf("a batch behind an uncrossable gap must be re-queued, got %v", rq)
		}
		if len(dst.updates) > chain.MaxClientUpdateHops {
			t.Fatalf("submitted %d standalone updates past a ceiling of %d", len(dst.updates), chain.MaxClientUpdateHops)
		}
		if dst.foldCalls != 0 {
			t.Fatalf("nothing reached the target, so nothing may be folded; foldCalls=%d", dst.foldCalls)
		}
	})

	// A hop that installs nothing must end the climb rather than repeat it -- the
	// mirror of the same guard on updateClientTo. Without it a builder stuck on one
	// height burns the whole ceiling in proofs and transactions, per flush, on the
	// packet path.
	t.Run("a hop that does not advance ends the climb", func(t *testing.T) {
		src := &mockSource{latest: 100}
		dst := &foldingMockDest{enabled: true}
		m := NewModule("test", "client-0", src, dst, &mockBuilder{hops: []uint64{50, 50, 50, 50}})

		batch := []chain.Event{{Type: chain.SendPacket, Height: 90, Raw: []byte("a")}}
		if rq := m.handleBatch(context.Background(), batch); len(rq) != 1 {
			t.Fatalf("a packet above a stalled climb must be re-queued, got %v", rq)
		}
		if len(dst.updates) != 1 {
			t.Fatalf("submitted %d updates for a builder stuck at one height, want 1", len(dst.updates))
		}
		if dst.foldCalls != 0 {
			t.Fatalf("the climb never reached the target, so nothing may be folded; foldCalls=%d", dst.foldCalls)
		}
	})

	// An adapter that never hops must behave exactly as it did before the climb
	// existed: one Build, nothing submitted on its own, everything folded.
	t.Run("a builder that never hops is unchanged", func(t *testing.T) {
		src := &mockSource{latest: 100}
		dst := &foldingMockDest{enabled: true}
		b := &mockBuilder{forceHeight: 80}
		m := NewModule("test", "client-0", src, dst, b)
		m.lastHeight = 10

		batch := []chain.Event{{Type: chain.SendPacket, Height: 30, Raw: []byte("a")}}
		if rq := m.handleBatch(context.Background(), batch); len(rq) != 0 {
			t.Fatalf("re-queued %v", rq)
		}
		if b.built != 1 || len(dst.updates) != 0 || dst.foldCalls != 1 {
			t.Fatalf("built=%d standalone=%d folds=%d; want one build, no standalone tx, one fold",
				b.built, len(dst.updates), dst.foldCalls)
		}
	})
}

// C4 asks for a decision, not a default: the retry spacing after a failed
// refresh must be stated and bounded. The bound is what makes backing off safe
// here -- the loop exists to act inside the refresh margin, so a backoff that can
// grow past it trades a gas leak for an expired client.
func TestNextRefreshBackoff(t *testing.T) {
	// A long trusting period, so the margin is the default hour and the ceiling
	// is what bounds the backoff. The short-period case is the test below.
	const longPeriod = 24 * time.Hour
	ceiling := refreshBackoffMaxFor(longPeriod)

	t.Run("the first failure waits one tick, then doubles", func(t *testing.T) {
		first := nextRefreshBackoff(0, longPeriod)
		if first != refreshTick {
			t.Fatalf("first backoff = %s, want one tick (%s)", first, refreshTick)
		}
		if second := nextRefreshBackoff(first, longPeriod); second != 2*first {
			t.Fatalf("second backoff = %s, want %s: a retry that waits the same as the last one is not spacing out",
				second, 2*first)
		}
	})

	t.Run("it stops at the cap", func(t *testing.T) {
		if got := nextRefreshBackoff(ceiling, longPeriod); got != ceiling {
			t.Fatalf("backoff grew past the cap to %s, want %s", got, ceiling)
		}
		if got := nextRefreshBackoff(ceiling*4, longPeriod); got != ceiling {
			t.Fatalf("backoff = %s from above the cap, want it clamped to %s", got, ceiling)
		}
	})
}

// The constraint C4 states outright, checked at EVERY trusting period rather than
// at the one the fixed margin used to be.
//
// Found in review by @neitdung: the cap was a fixed five minutes, written when
// the margin was a fixed thirty. main derives the margin from the destination's
// own trusting period, so at twenty minutes or less the margin is five minutes or
// less and one backoff step could span the whole window -- the anti-expiry
// routine sitting out exactly the period it exists for. A fixed cap cannot be
// made safe by choosing a smaller number, because the margin has no lower bound
// that a constant can stay under; the cap has to be derived from the same input.
func TestRefreshBackoffStaysInsideTheWindowAtEveryTrustingPeriod(t *testing.T) {
	for _, trustingPeriod := range []time.Duration{
		0, // the destination reports none: the default margin stands in
		30 * time.Second, 2 * time.Minute, 4 * time.Minute, 10 * time.Minute,
		20 * time.Minute, time.Hour, 4 * time.Hour, 14 * 24 * time.Hour,
	} {
		t.Run(trustingPeriod.String(), func(t *testing.T) {
			margin := refreshMarginFor(trustingPeriod)
			cap := refreshBackoffMaxFor(trustingPeriod)
			if cap >= margin {
				t.Fatalf("backoff cap %s is not below the margin %s; one gap could span the window the loop exists to work inside",
					cap, margin)
			}

			attempts := 0
			for elapsed, backoff := time.Duration(0), time.Duration(0); elapsed < margin; attempts++ {
				backoff = nextRefreshBackoff(backoff, trustingPeriod)
				elapsed += backoff
			}
			if attempts < 4 {
				t.Fatalf("only %d attempts fit in the %s window (cap %s); the cap is too high to be an emergency routine",
					attempts, margin, cap)
			}
		})
	}
}

// The wiring half: the loop must actually space its attempts out. The pure
// function above can be perfect while the loop calls it and ignores the answer.
func TestRefreshLoop_BacksOffWhenRefreshKeepsFailing(t *testing.T) {
	withFastRefreshTick(t) // 1ms, so the backoff doubles in milliseconds

	src := &mockSource{latest: 42}
	dst := &countingExpiryDest{expiresErr: errors.New("rpc down")}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); m.refreshLoop(ctx) }()

	const window = 300 * time.Millisecond
	time.Sleep(window)
	cancel()
	wg.Wait()

	attempts := dst.calls()
	// Unspaced, a 1ms tick would try ~300 times in this window. Backing off 1, 2,
	// 4, 8 ... reaches 300ms in nine attempts. The bound is loose because the
	// scheduler is not, but it cannot be met by a loop that ignores the backoff.
	if attempts > 30 {
		t.Fatalf("%d attempts in %s: the loop is retrying on every tick, not backing off", attempts, window)
	}
	if attempts < 2 {
		t.Fatalf("%d attempts in %s: the loop stopped retrying altogether", attempts, window)
	}
}

// Backing off must not wedge the loop: once the cause clears, the client still
// gets advanced. A backoff that outlived its failure would be the expired-client
// outcome the cap exists to prevent, arriving by another route.
func TestRefreshLoop_RecoversAfterTheFailureClears(t *testing.T) {
	withFastRefreshTick(t)

	src := &mockSource{latest: 42}
	dst := &countingExpiryDest{expiresErr: errors.New("rpc down")}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); m.refreshLoop(ctx) }()

	waitFor(t, ctx, func() bool { return dst.calls() >= 3 }, "the loop never made its first few attempts")
	dst.heal(time.Now().Add(time.Minute)) // inside refreshMargin, so a refresh is due

	waitFor(t, ctx, func() bool {
		m.mu.Lock()
		defer m.mu.Unlock()
		return m.lastHeight == 42
	}, "the loop never advanced the client after the failure cleared; the backoff outlived its cause")
	cancel()
	wg.Wait()
}

// A backoff that survives its own cause is a slow leak: a client that failed a
// dozen times last week and recovered would wait the full cap before its next
// isolated blip is retried, inside a window that exists to keep trying.
func TestRefreshLoop_BackoffDoesNotOutliveItsCause(t *testing.T) {
	withFastRefreshTick(t)

	src := &mockSource{latest: 42}
	dst := &countingExpiryDest{expiresErr: errors.New("rpc down")}
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); m.refreshLoop(ctx) }()

	// Ten failures at a 1ms tick take the backoff to roughly half a second, which
	// is what makes the difference below measurable rather than theoretical.
	waitFor(t, ctx, func() bool { return dst.calls() >= 10 }, "the loop never accumulated a backoff to clear")
	dst.heal(time.Now().Add(time.Minute))
	waitFor(t, ctx, func() bool {
		m.mu.Lock()
		defer m.mu.Unlock()
		return m.lastHeight == 42
	}, "the loop never recovered")

	// Break it again. With the backoff cleared by the success, the next attempts
	// are one tick apart; carried over, the first would be half a second away.
	before := dst.calls()
	dst.fail(errors.New("rpc down again"))
	deadline, ok := context.WithTimeout(ctx, 300*time.Millisecond)
	defer ok()
	promptly := true
	waitForSoft(deadline, func() bool { return dst.calls() >= before+3 }, &promptly)
	if !promptly {
		t.Fatalf("only %d attempts in 300ms after a fresh failure; the backoff was carried over from the outage that already ended",
			dst.calls()-before)
	}
	cancel()
	wg.Wait()
}

// waitForSoft is waitFor without the t.Fatal: the caller decides what a timeout
// means, because here it is the assertion rather than a broken test.
func waitForSoft(ctx context.Context, cond func() bool, ok *bool) {
	for {
		if cond() {
			return
		}
		select {
		case <-ctx.Done():
			*ok = false
			return
		case <-time.After(time.Millisecond):
		}
	}
}

// D5: three states, and two of them used to be silence. An operator reading the
// log could not tell "the client has hours left" from "it has minutes left and
// the last refresh worked".
func TestRefreshLoop_SaysWhatItIsDoing(t *testing.T) {
	withFastRefreshTick(t)

	var logs bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(oldOutput) })

	src := &mockSource{latest: 42}
	dst := &mockDest{expiresAt: time.Now().Add(time.Minute)} // inside refreshMargin
	m := NewModule("test", "client-0", src, dst, &mockBuilder{})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); m.refreshLoop(ctx) }()
	waitFor(t, ctx, func() bool {
		m.mu.Lock()
		defer m.mu.Unlock()
		return m.lastHeight == 42
	}, "the refresh loop never advanced the client")
	cancel()
	wg.Wait()

	out := logs.String()
	if !strings.Contains(out, "client expires in") {
		t.Errorf("no line saying a refresh was due, so a due refresh is indistinguishable from a healthy client:\n%s", out)
	}
	if !strings.Contains(out, "advanced the client to height 42") {
		t.Errorf("no line saying the refresh succeeded, so success is indistinguishable from failure:\n%s", out)
	}
	if !strings.Contains(out, "of trust left") {
		t.Errorf("the lines carry no remaining lifetime, which is the number an operator is reading for:\n%s", out)
	}
}

// countingExpiryDest counts expiry queries -- one per refresh ATTEMPT, since it
// is the first call the loop makes -- and can be healed mid-test.
type countingExpiryDest struct {
	mockDest
	mu             sync.Mutex
	n              int
	expiresErr     error
	expiresAt      time.Time
	trustingPeriod time.Duration
}

func (d *countingExpiryDest) ClientExpiresAt(context.Context, string) (time.Time, time.Duration, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.n++
	return d.expiresAt, d.trustingPeriod, d.expiresErr
}

func (d *countingExpiryDest) calls() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.n
}

func (d *countingExpiryDest) heal(expiresAt time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.expiresErr, d.expiresAt = nil, expiresAt
}

func (d *countingExpiryDest) fail(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.expiresErr = err
}

func TestFlushLoop_RunsImmediately(t *testing.T) {
	src := &listerSource{mockSource: mockSource{latest: 100, relayable: 100}}
	m := NewModule("cosmos->eth", "client-0", src, &mockDest{}, &mockBuilder{})
	m.flushInterval = time.Hour // so only the immediate pass can run

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { defer close(done); m.flushLoop(ctx, src) }()

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && src.callCount() == 0 {
		time.Sleep(time.Millisecond)
	}
	cancel()
	<-done

	if got := src.callCount(); got != 1 {
		t.Fatalf("enumeration ran %d times before the first tick, want 1", got)
	}
}

type listerSource struct {
	mockSource
	mu        sync.Mutex
	calls     int
	candidate []chain.Event
	listErr   error
}

func (l *listerSource) UnrelayedPackets(context.Context) ([]chain.Event, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.calls++
	return append([]chain.Event(nil), l.candidate...), l.listErr
}

func (l *listerSource) callCount() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.calls
}

// handleBatch has two callers now -- the Subscribe callback and the flush loop,
// on separate goroutines -- and it writes m.lastWait, which no lock covered.
// Under -race that write is a data race without Module.batchMu.
//
// The race is the smaller half. Two batches in flight at once can each build and
// submit a client update, and can relay the same packet twice: the flush
// enumerates outstanding commitments, which includes the packets the scan is
// delivering at that moment.
//
// Overlap is measured from INSIDE the batch, through a source call handleBatch
// makes while holding the lock. Counting around the call instead would count
// goroutines queued on the lock, not batches running in it -- which is what a
// first version of this test did, and it failed against correct code.
type overlapProbeSource struct {
	mockSource
	inFlight atomic.Int32
	overlap  atomic.Bool
}

func (s *overlapProbeSource) RelayableHeight(ctx context.Context) (uint64, error) {
	if s.inFlight.Add(1) > 1 {
		s.overlap.Store(true)
	}
	time.Sleep(time.Millisecond) // widen the window a real batch would occupy
	s.inFlight.Add(-1)
	return s.mockSource.RelayableHeight(ctx)
}

func TestHandleBatchIsSerializedAcrossCallers(t *testing.T) {
	source := &overlapProbeSource{}
	m := NewModule("test", "client", source, &mockDest{}, &mockBuilder{})
	events := []chain.Event{{Type: chain.SendPacket, Sequence: 1, Height: 500}}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.handleBatch(context.Background(), events)
		}()
	}
	wg.Wait()

	if source.overlap.Load() {
		t.Fatal("two handleBatch calls were inside the batch at once; batchMu is not held for the whole of it")
	}
}

// Reported by @neitdung. A folded transaction carries the client update AND the
// packets, so a deterministic rejection of it says nothing about WHICH of the
// two the chain refused -- but the code read every permanent failure as "the
// packet is poison". A single packet was therefore dropped outright, and a batch
// was retried as folded singletons carrying the SAME update, so each failed for
// the same reason and every packet was dropped in turn.
//
// An ack has no independent escape: its counterparty already holds a receipt, so
// nothing times it out and nothing refunds it. Dropping it on a failure never
// shown to be its own strands the escrow.
//
// The fixture has to reject the update on BOTH paths, folded and standalone --
// that is what "the update is the problem" means. A fold that fails while the
// standalone update succeeds is the opposite case, and the last subtest pins
// that the two are still told apart.
func TestHandleBatch_FoldedUpdateRejectionDoesNotDropInnocentPackets(t *testing.T) {
	rejected := func() *foldingMockDest {
		dst := &foldingMockDest{enabled: true, foldErr: chain.Permanent(errors.New("client update rejected"))}
		dst.updateErr = chain.Permanent(errors.New("client update rejected"))
		return dst
	}

	t.Run("a single ack is held, not dropped", func(t *testing.T) {
		src := &mockSource{latest: 100}
		dst := rejected()
		m := NewModule("test", "client-0", src, dst, &mockBuilder{forceHeight: 80})
		m.lastHeight = 10

		batch := []chain.Event{{Type: chain.AckPacket, Sequence: 7, Height: 30, Raw: []byte("a")}}
		requeue := m.handleBatch(context.Background(), batch)

		if len(requeue) != 1 {
			t.Fatalf("re-queued %v, want the ack held: the update was refused on its own, "+
				"so nothing was ever shown about the packet", requeue)
		}
		if m.lastHeight != 10 {
			t.Fatalf("lastHeight = %d, want 10: an update the chain refused must not be recorded", m.lastHeight)
		}
	})

	// The batch shape is the worse one: isolation would have retried every packet
	// with the same rejected update and dropped all of them.
	t.Run("a whole batch is held, not dropped one by one", func(t *testing.T) {
		src := &mockSource{latest: 100}
		dst := rejected()
		m := NewModule("test", "client-0", src, dst, &mockBuilder{forceHeight: 80})
		m.lastHeight = 10

		batch := []chain.Event{
			{Type: chain.AckPacket, Sequence: 1, Height: 30, Raw: []byte("a")},
			{Type: chain.AckPacket, Sequence: 2, Height: 31, Raw: []byte("b")},
			{Type: chain.TimeoutPacket, Sequence: 3, Height: 32, Raw: []byte("c")},
		}
		requeue := m.handleBatch(context.Background(), batch)

		if len(requeue) != len(batch) {
			t.Fatalf("re-queued %v, want all %d held", requeue, len(batch))
		}
	})

	// Second round, also @neitdung. The update landing alone clears the UPDATE,
	// and nothing more: the folded transaction was one atomic unit, and a unit can
	// be refused for what the two are together (its gas envelope, its route)
	// while each half passes on its own. The first fix drew the stronger
	// conclusion -- "the update landed, so the packet is poison" -- and dropped a
	// single packet on the folded error without ever sending it by itself. The
	// batch shape went through packet-only isolation; the singleton skipped it.
	t.Run("a packet that fails only in the folded composition is delivered, not dropped", func(t *testing.T) {
		src := &mockSource{latest: 100}
		dst := &foldingMockDest{enabled: true, foldErr: chain.Permanent(errors.New("folded tx reverted"))}
		m := NewModule("test", "client-0", src, dst, &mockBuilder{forceHeight: 80})
		m.lastHeight = 10

		batch := []chain.Event{{Type: chain.AckPacket, Sequence: 7, Height: 30, Raw: []byte("a")}}
		if rq := m.handleBatch(context.Background(), batch); len(rq) != 0 {
			t.Fatalf("re-queued %v, want nothing outstanding: the packet relays on its own", rq)
		}
		if dst.relayCalls != 1 {
			t.Fatalf("packet-only relay calls = %d, want 1: the folded error says nothing about the packet alone", dst.relayCalls)
		}
		if len(dst.relayed) != 1 || dst.relayed[0].Sequence != 7 {
			t.Fatalf("delivered %+v, want the ack: it was dropped on a failure never shown to be its own", dst.relayed)
		}
		if m.lastHeight != 80 {
			t.Fatalf("lastHeight = %d, want 80: the standalone update installed", m.lastHeight)
		}
		if len(dst.updates) != 1 || dst.updates[0].Height != 80 {
			t.Fatalf("standalone update calls = %+v, want exactly the folded update", dst.updates)
		}
	})

	// The control: a packet that fails on its own too is still dropped -- after
	// that attempt, on that attempt's error, not on the folded one.
	t.Run("a packet that really is poison is dropped once it has failed alone", func(t *testing.T) {
		src := &mockSource{latest: 100}
		dst := &foldingMockDest{enabled: true, foldErr: chain.Permanent(errors.New("folded tx reverted"))}
		dst.poison = map[string]bool{"a": true}
		m := NewModule("test", "client-0", src, dst, &mockBuilder{forceHeight: 80})
		m.lastHeight = 10

		batch := []chain.Event{{Type: chain.AckPacket, Sequence: 7, Height: 30, Raw: []byte("a")}}
		if rq := m.handleBatch(context.Background(), batch); len(rq) != 0 {
			t.Fatalf("re-queued %v; the packet failed on its own, so it is dropped", rq)
		}
		if dst.relayCalls != 1 {
			t.Fatalf("packet-only relay calls = %d, want 1: dropping needs the packet's own failure first", dst.relayCalls)
		}
		if len(dst.relayed) != 0 {
			t.Fatalf("delivered %+v, want nothing: the packet is poison", dst.relayed)
		}
		if m.lastHeight != 80 {
			t.Fatalf("lastHeight = %d, want 80: the standalone update installed", m.lastHeight)
		}
	})
}

// A failing expiry query must keep sizing its retry against THIS client's margin.
//
// Reported by @DongLieu on #487, against the commit that fixed the same class of
// bug one line away. The failure path passed a zero trusting period, which falls
// back to the default margin, whose cap is the five-minute ceiling. On a
// ten-minute trusting period the real margin is 2m30s, so a run of query failures
// could schedule a retry spanning twice the window the loop exists to work
// inside -- and a query that keeps failing is exactly the case where no later
// pass arrives to narrow it.
func TestRefreshBackoffKeepsTheLastKnownTrustingPeriodWhenTheQueryFails(t *testing.T) {
	const trustingPeriod = 10 * time.Minute

	// The rule. A failing query must not fall back to the default margin's cap,
	// because that cap is the larger one for every short trusting period.
	t.Run("the default cap is the wrong size for a short period", func(t *testing.T) {
		margin := refreshMarginFor(trustingPeriod)
		if refreshBackoffMaxFor(0) < margin {
			t.Fatal("the default cap already fits inside this margin; there would be nothing to fix")
		}
		if refreshBackoffMaxFor(trustingPeriod) >= margin {
			t.Fatalf("this client's own cap %s does not fit inside its %s margin",
				refreshBackoffMaxFor(trustingPeriod), margin)
		}
	})

	// The wiring, asserted structurally. Driving the loop to saturation is not
	// practical here: the backoff starts at one tick and doubles, so reaching a
	// 37s cap takes about fifteen failures and the loop SLEEPS each backoff in
	// between -- minutes of wall clock for one assertion. The early spacings are
	// identical whichever period is passed, so a log-reading test passes against
	// the bug; I tried that first and the mutation walked through it.
	//
	// So the assertion is on the call site, the same technique services uses for
	// fsync: the property is invisible from outside at test speed, and a test that
	// cannot see it should say which half it is pinning.
	t.Run("the failure path passes the remembered period, not zero", func(t *testing.T) {
		src, err := os.ReadFile("module.go")
		if err != nil {
			t.Fatalf("read module.go: %v", err)
		}
		body := string(src)
		start := strings.Index(body, "func (m *Module) refreshLoop(")
		if start < 0 {
			t.Fatal("refreshLoop not found")
		}
		end := strings.Index(body[start:], "\nfunc ")
		if end < 0 {
			end = len(body) - start
		}
		fn := body[start : start+end]

		if !strings.Contains(fn, "lastTrustingPeriod = trustingPeriod") {
			t.Error("refreshLoop does not remember the period it read, so a failing query has nothing to size with")
		}
		if !strings.Contains(fn, `refreshFailed(backoff, lastTrustingPeriod, "query client expiry"`) {
			t.Error("the expiry-query failure path does not size its retry with the remembered trusting period; " +
				"falling back to the default lets the backoff span the refresh window on a short-lived client")
		}
	})
}
