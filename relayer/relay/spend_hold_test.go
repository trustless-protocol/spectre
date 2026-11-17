package relay

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"relayer/chain"
	"relayer/services"
)

// capRefusal is what a transaction handler returns when an operator spending cap
// refuses to sign.
var capRefusal = fmt.Errorf("base fee 90 + tip 2 wei per gas is above ETH_MAX_FEE_PER_GAS=50; not signed: %w",
	services.ErrSpendCapExceeded)

// capOnIsolateDest reverts any multi-packet batch and refuses every singleton at
// the spending cap: fees rose while the batch was being isolated.
type capOnIsolateDest struct{ mockDest }

func (d *capOnIsolateDest) RelayPackets(_ context.Context, packets []chain.RelayPacket) (int, error) {
	d.relayCalls++
	if len(packets) > 1 {
		return 0, chain.Permanent(errors.New("batch reverted"))
	}
	return 0, fmt.Errorf("send: %w", capRefusal)
}

// packetSends counts the packet transactions a test destination was asked for.
func packetSends(t *testing.T, d chain.Destination) int {
	t.Helper()
	switch d := d.(type) {
	case *foldingMockDest:
		return d.relayCalls + d.foldCalls
	case *capOnIsolateDest:
		return d.relayCalls
	case *mockDest:
		return d.relayCalls
	}
	t.Fatalf("unknown destination %T", d)
	return 0
}

// #575 M3: the caps are checked when a transaction is signed, after the client
// update and the proofs are built. A refused batch came back within 15 s and
// rebuilt both -- a Groth16 proof on a Cosmos->EVM path -- for as long as the cap
// stayed binding. Every send in a batch that the cap refuses holds the direction.
func TestSpendCapHoldsTheDirection(t *testing.T) {
	unprovable := map[string]bool{"a": true, "b": true, "c": true}
	cases := []struct {
		name string
		dst  chain.Destination
		src  *mockSource
	}{
		{"refused client update", &mockDest{updateErr: capRefusal}, &mockSource{latest: 100}},
		{"refused packet transaction", &mockDest{relayErr: capRefusal}, &mockSource{latest: 100}},
		{
			"refused folded transaction, then its update sent apart",
			&foldingMockDest{enabled: true, foldErr: capRefusal, mockDest: mockDest{updateErr: capRefusal}},
			&mockSource{latest: 100},
		},
		{
			"refused folded transaction, then its packets sent apart",
			&foldingMockDest{enabled: true, foldErr: capRefusal, mockDest: mockDest{relayErr: capRefusal}},
			&mockSource{latest: 100},
		},
		{
			"refused update sent with no provable packet",
			&foldingMockDest{enabled: true, mockDest: mockDest{updateErr: capRefusal}},
			&mockSource{latest: 100, failMembershipOn: unprovable},
		},
		{
			"refused update after a reverted folded transaction",
			&foldingMockDest{enabled: true, foldErr: chain.Permanent(errors.New("reverted")), mockDest: mockDest{updateErr: capRefusal}},
			&mockSource{latest: 100},
		},
		{"refused singleton while isolating", &capOnIsolateDest{}, &mockSource{latest: 100}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := &mockBuilder{}
			m := NewModule("cosmos->eth", "client-0", tc.src, tc.dst, b)

			var rq []int
			out := captureLog(func() { rq = m.handleBatch(context.Background(), threePackets()) })
			if len(rq) != 3 || strings.Contains(out, "DROP") {
				t.Fatalf("a refusal must keep every packet: rq=%v\n%s", rq, out)
			}
			if m.spendProbeAt.IsZero() {
				t.Fatalf("a spending-cap refusal must hold the direction:\n%s", out)
			}
			if !m.sourceProbeAt.IsZero() {
				t.Fatal("a spending cap must not be blamed on the source")
			}

			// Held: the next batch builds nothing, proves nothing, sends nothing.
			built, proofs, sends := b.built, tc.src.membershipCalls, packetSends(t, tc.dst)
			rq = m.handleBatch(context.Background(), threePackets())
			if len(rq) != 3 || b.built != built || tc.src.membershipCalls != proofs || packetSends(t, tc.dst) != sends {
				t.Fatalf("held batch must cost nothing: rq=%v builds %d->%d proofs %d->%d sends %d->%d",
					rq, built, b.built, proofs, tc.src.membershipCalls, sends, packetSends(t, tc.dst))
			}
		})
	}

	t.Run("an RPC failure does not hold", func(t *testing.T) {
		m := NewModule("cosmos->eth", "client-0", &mockSource{latest: 100},
			&mockDest{relayErr: errors.New("connection refused")}, &mockBuilder{})
		if rq := m.handleBatch(context.Background(), threePackets()); len(rq) != 3 {
			t.Fatalf("want all 3 re-queued, got %v", rq)
		}
		if !m.spendProbeAt.IsZero() {
			t.Fatal("only a spending-cap refusal holds")
		}
	})
}

// #575 M4: a Cosmos cap bounds one transaction, and the folded one -- the update
// and every packet -- can be over it while each part is not. The path held and
// retried that same folded batch, each time with an update to a newer height,
// and an L2->Cosmos path, whose client nothing else updates, stopped. The parts
// now go apart: the update alone, then the packets on the path that splits.
func TestSpendCapRefusedFoldedTransactionIsSentApart(t *testing.T) {
	dst := &foldingMockDest{enabled: true, foldErr: capRefusal}
	src := &mockSource{latest: 100}
	trusted := time.Unix(1_700_000_000, 0)
	var observed []time.Time
	m := NewModule("l2->cosmos", "client-0", src, dst, &mockBuilder{trustedAt: trusted},
		WithClientUpdateObserver(func(at time.Time) { observed = append(observed, at) }))

	if rq := m.handleBatch(context.Background(), threePackets()); len(rq) != 0 {
		t.Fatalf("want every packet delivered, got %v re-queued", rq)
	}
	if dst.foldCalls != 1 || len(dst.updates) != 1 || dst.updates[0].Height != 100 {
		t.Fatalf("want the refused folded tx, then the update at 100 alone: folds=%d updates=%v", dst.foldCalls, dst.updates)
	}
	if dst.relayCalls != 1 || len(dst.relayed) != 3 {
		t.Fatalf("want the 3 packets in one packets-only send, got %d call(s), %d relayed", dst.relayCalls, len(dst.relayed))
	}
	for _, h := range src.proofHeights {
		if h != 100 {
			t.Fatalf("proof heights = %v, want every proof at the update's 100", src.proofHeights)
		}
	}
	if m.lastHeight != 100 || len(observed) != 1 || !observed[0].Equal(trusted) {
		t.Fatalf("lastHeight = %d, observed %v: want the update at 100 recorded once", m.lastHeight, observed)
	}
	if !m.spendProbeAt.IsZero() {
		t.Fatal("a refusal its parts get past must not hold the path")
	}

	// Once the update has landed, a packet that reverts is that packet's own
	// failure: it is isolated like any packets-only send, and the update is not
	// submitted a second time to attribute it.
	t.Run("a reverting packet after the update landed is isolated", func(t *testing.T) {
		dst := &foldingMockDest{enabled: true, foldErr: capRefusal, mockDest: mockDest{poison: map[string]bool{"b": true}}}
		m := NewModule("l2->cosmos", "client-0", &mockSource{latest: 100}, dst, &mockBuilder{})

		var rq []int
		out := captureLog(func() { rq = m.handleBatch(context.Background(), threePackets()) })
		if len(dst.updates) != 1 {
			t.Fatalf("the update was submitted %d times, want once", len(dst.updates))
		}
		if len(rq) != 0 || len(dst.relayed) != 2 || !strings.Contains(out, "DROP packet (ack seq=2") {
			t.Fatalf("want a and c delivered and the reverting b dropped: rq=%v relayed=%d\n%s", rq, len(dst.relayed), out)
		}
	})

	// The client froze between the folded refusal and the update sent apart: that
	// arms the client hold, not the spending-cap one, and nothing is dropped.
	t.Run("a client frozen before the update apart holds as frozen", func(t *testing.T) {
		dst := &foldingMockDest{enabled: true, foldErr: capRefusal, mockDest: mockDest{updateErr: frozenRevert}}
		src := &mockSource{latest: 100}
		b := &mockBuilder{}
		m := NewModule("l2->cosmos", "client-0", src, dst, b)

		if rq := m.handleBatch(context.Background(), threePackets()); len(rq) != 3 {
			t.Fatalf("want all 3 kept, got %v", rq)
		}
		if m.clientHold.Load() == nil || !m.spendProbeAt.IsZero() {
			t.Fatalf("want the client hold and no spending-cap hold: client=%v spendProbeAt=%v", m.clientHold.Load(), m.spendProbeAt)
		}
		built, proofs, sends := b.built, src.membershipCalls, packetSends(t, dst)
		if rq := m.handleBatch(context.Background(), threePackets()); len(rq) != 3 ||
			b.built != built || src.membershipCalls != proofs || packetSends(t, dst) != sends {
			t.Fatalf("held batch must cost nothing: rq=%v builds %d->%d proofs %d->%d sends %d->%d",
				rq, built, b.built, proofs, src.membershipCalls, sends, packetSends(t, dst))
		}
	})
}

// The hold is tried again once it expires: a repeat refusal widens it on the
// 1m -> 15m ladder, and a send that lands lifts it.
func TestSpendCapHoldIsProbedAndLifted(t *testing.T) {
	dst := &mockDest{relayErr: capRefusal}
	m := NewModule("cosmos->eth", "client-0", &mockSource{latest: 100}, dst, &mockBuilder{})

	m.handleBatch(context.Background(), threePackets())
	if m.spendBackoff != time.Minute {
		t.Fatalf("first hold = %s, want 1m", m.spendBackoff)
	}

	m.spendProbeAt = time.Now().Add(-time.Second) // the hold expired
	m.handleBatch(context.Background(), threePackets())
	if dst.relayCalls != 2 || m.spendBackoff != 2*time.Minute {
		t.Fatalf("an expired hold must be tried once and widened on a repeat refusal: relayCalls=%d hold=%s",
			dst.relayCalls, m.spendBackoff)
	}

	m.spendProbeAt = time.Now().Add(-time.Second)
	dst.relayErr = nil // fees fell
	var rq []int
	out := captureLog(func() { rq = m.handleBatch(context.Background(), threePackets()) })
	if len(rq) != 0 || len(dst.relayed) != 3 {
		t.Fatalf("want all 3 relayed once fees fall, got rq=%v relayed=%d", rq, len(dst.relayed))
	}
	if m.spendBackoff != 0 || !m.spendProbeAt.IsZero() || !strings.Contains(out, "hold lifted") {
		t.Fatalf("a landed send must lift the hold: hold=%s probeAt=%v\n%s", m.spendBackoff, m.spendProbeAt, out)
	}
}
