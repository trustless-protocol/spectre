package relay

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"relayer/chain"
)

// fakeClock drives the tracker's notion of time so age assertions are exact
// rather than timing-dependent.
type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time { return c.t }

func newTestTracker() (*waitTracker, *fakeClock) {
	clk := &fakeClock{t: time.Unix(1_700_000_000, 0)}
	w := newWaitTracker()
	w.now = clk.now
	return w, clk
}

func ackEvent(seq, height uint64) chain.Event {
	return chain.Event{Type: chain.AckPacket, Sequence: seq, Height: height}
}

// TestObserveAccumulatesAge is the property the whole change exists for: a
// packet re-observed across flushes must report a GROWING age, because that is
// the only signal separating "waiting for finality" from "stuck forever".
func TestObserveAccumulatesAge(t *testing.T) {
	w, clk := newTestTracker()
	e := ackEvent(1, 21)

	if age := w.observe(e); age != 0 {
		t.Fatalf("first observation age = %s, want 0 (no age can be known yet)", age)
	}

	clk.t = clk.t.Add(90 * time.Second)
	if age := w.observe(e); age != 90*time.Second {
		t.Fatalf("age after 90s = %s, want 1m30s", age)
	}

	clk.t = clk.t.Add(30 * time.Second)
	if age := w.observe(e); age != 2*time.Minute {
		t.Fatalf("age after a further 30s = %s, want 2m0s", age)
	}
}

// TestClearResetsAge covers the settle path: a delivered packet must not leave
// an entry behind that would inflate the age of a later packet reusing the same
// sequence (which happens when a client is re-created and sequences restart).
func TestClearResetsAge(t *testing.T) {
	w, clk := newTestTracker()
	e := ackEvent(1, 21)

	w.observe(e)
	clk.t = clk.t.Add(5 * time.Minute)
	w.clearPacket(chain.RelayPacket{Type: e.Type, Sequence: e.Sequence})

	if age := w.observe(e); age != 0 {
		t.Fatalf("age after clear = %s, want 0 (the entry must be gone)", age)
	}
}

// TestRecvAndAckAreTrackedSeparately: a packet's forward delivery and its
// returning ack share a sequence, so keying on sequence alone would make one
// leg's clear wipe the other's age.
func TestRecvAndAckAreTrackedSeparately(t *testing.T) {
	w, clk := newTestTracker()
	recv := chain.Event{Type: chain.SendPacket, Sequence: 7, Height: 10}
	ack := ackEvent(7, 21)

	w.observe(recv)
	clk.t = clk.t.Add(time.Minute)
	w.observe(ack)
	clk.t = clk.t.Add(time.Minute)

	w.clearPacket(chain.RelayPacket{Type: chain.SendPacket, Sequence: 7})

	if age := w.observe(ack); age != time.Minute {
		t.Fatalf("ack age after clearing the recv = %s, want 1m0s (clears must not cross legs)", age)
	}
}

// TestPurgeStaleBoundsTheTracker: entries for packets that left the queue via a
// path with no clear must not accumulate for the process's lifetime.
func TestPurgeStaleBoundsTheTracker(t *testing.T) {
	w, clk := newTestTracker()
	stale := ackEvent(1, 21)
	fresh := ackEvent(2, 22)

	w.observe(stale)
	clk.t = clk.t.Add(waitEntryTTL + time.Minute)
	w.observe(fresh) // refreshes only fresh's lastSeen

	w.purgeStale()

	if len(w.entries) != 1 {
		t.Fatalf("entries after purge = %d, want 1 (only the fresh one survives)", len(w.entries))
	}
	if age := w.observe(stale); age != 0 {
		t.Fatalf("purged entry still carried an age of %s, want 0", age)
	}
}

// TestPurgeKeepsActivelyWaitingEntries guards the inverse mistake: purging on
// first-seen instead of last-seen would silently reset the age of exactly the
// long-waiting packet the operator is trying to diagnose.
func TestPurgeKeepsActivelyWaitingEntries(t *testing.T) {
	w, clk := newTestTracker()
	e := ackEvent(1, 21)

	w.observe(e)
	// Re-observed every 15s (the waiting backoff cap) for well past the TTL.
	for elapsed := time.Duration(0); elapsed < waitEntryTTL+10*time.Minute; elapsed += 15 * time.Second {
		clk.t = clk.t.Add(15 * time.Second)
		w.observe(e)
		w.purgeStale()
	}

	if age := w.observe(e); age < waitEntryTTL {
		t.Fatalf("actively-waiting packet lost its age: %s, want >= %s", age, waitEntryTTL)
	}
}

func TestDescribeWaitingNamesEachPacket(t *testing.T) {
	got := describeWaiting([]waiting{
		{event: ackEvent(1, 21), age: 4*time.Minute + 12*time.Second},
	})
	want := "ack seq=1 h=21 for 4m12s"
	if got != want {
		t.Fatalf("describeWaiting = %q, want %q", got, want)
	}
}

// TestDescribeWaitingOrdersByAge: with the list capped, the entries shown must
// be the oldest ones — those are the ones worth looking at.
func TestDescribeWaitingOrdersByAge(t *testing.T) {
	got := describeWaiting([]waiting{
		{event: ackEvent(1, 21), age: time.Second},
		{event: ackEvent(2, 22), age: time.Hour},
		{event: ackEvent(3, 23), age: time.Minute},
	})
	if !strings.HasPrefix(got, "ack seq=2 ") {
		t.Fatalf("describeWaiting = %q, want the 1h entry first", got)
	}
	if strings.Index(got, "seq=3") > strings.Index(got, "seq=1") {
		t.Fatalf("describeWaiting = %q, want seq=3 (1m) before seq=1 (1s)", got)
	}
}

func TestDescribeWaitingCapsTheList(t *testing.T) {
	items := make([]waiting, 0, waitListLimit+3)
	for i := range waitListLimit + 3 {
		items = append(items, waiting{
			event: ackEvent(uint64(i+1), 21),
			age:   time.Duration(i) * time.Second,
		})
	}
	got := describeWaiting(items)
	if strings.Count(got, "seq=") != waitListLimit {
		t.Fatalf("describeWaiting listed %d entries, want %d: %q",
			strings.Count(got, "seq="), waitListLimit, got)
	}
	if !strings.HasSuffix(got, "+3 more") {
		t.Fatalf("describeWaiting = %q, want a trailing \"+3 more\"", got)
	}
}

func TestDescribeWaitingEmpty(t *testing.T) {
	if got := describeWaiting(nil); got != "" {
		t.Fatalf("describeWaiting(nil) = %q, want empty", got)
	}
}

func TestEventTypeNames(t *testing.T) {
	cases := map[chain.EventType]string{
		chain.SendPacket:    "recv",
		chain.AckPacket:     "ack",
		chain.TimeoutPacket: "timeout",
		chain.EventType(99): "unknown(99)",
	}
	for typ, want := range cases {
		if got := typ.String(); got != want {
			t.Fatalf("EventType(%d).String() = %q, want %q", int(typ), got, want)
		}
	}
}

// --- wiring: the age must survive across flushes through handleBatch itself ---

// captureLog redirects the standard logger for the duration of fn and returns
// what was written. The waiting output is the deliverable here, so it is asserted
// directly rather than through the tracker's internals.
func captureLog(fn func()) string {
	var buf strings.Builder
	prevOut := log.Writer()
	prevFlags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() {
		log.SetOutput(prevOut)
		log.SetFlags(prevFlags)
	}()
	fn()
	return buf.String()
}

// TestHandleBatch_WaitingLogNamesThePacket: the source-gate wait must identify
// the packet. A bare count cannot tell an operator whether the ack they are
// looking for is in the queue at all.
func TestHandleBatch_WaitingLogNamesThePacket(t *testing.T) {
	src := &mockSource{latest: 100, relayable: 50}
	m := NewModule("arbitrum->cosmos", "client-0", src, &mockDest{}, &mockBuilder{})

	out := captureLog(func() {
		m.handleBatch(context.Background(), []chain.Event{
			{Type: chain.AckPacket, Sequence: 7, Height: 80, Raw: []byte("ack"), AckBytes: [][]byte{{1}}},
		})
	})

	for _, want := range []string{"waiting:", "ack seq=7", "h=80", "source relayable height=50"} {
		if !strings.Contains(out, want) {
			t.Fatalf("waiting log missing %q:\n%s", want, out)
		}
	}
}

// TestHandleBatch_AgeGrowsAtDestinationGate is the regression guard for the
// subtle half of this change. A packet that clears the SOURCE gate every flush
// but parks at the DESTINATION-coverage gate is the exact shape of a stuck L2
// ack. Clearing the wait entry when it passes the source gate would reset its age
// to zero on every flush, so the log would report a permanently-stuck packet as
// freshly arrived — hiding the very condition the age exists to reveal.
func TestHandleBatch_AgeGrowsAtDestinationGate(t *testing.T) {
	src := &mockSource{latest: 100}
	// The builder can only ever advance the client to 50, so a packet at 100
	// clears RelayableHeight but never the destination-coverage check.
	m := NewModule("arbitrum->cosmos", "client-0", src, &mockDest{}, &mockBuilder{forceHeight: 50})

	clk := &fakeClock{t: time.Unix(1_700_000_000, 0)}
	m.waits.now = clk.now

	event := chain.Event{Type: chain.AckPacket, Sequence: 7, Height: 100, Raw: []byte("ack"), AckBytes: [][]byte{{1}}}
	batch := []chain.Event{event}

	m.handleBatch(context.Background(), batch) // first sighting: age 0
	clk.t = clk.t.Add(10 * time.Minute)

	out := captureLog(func() { m.handleBatch(context.Background(), batch) })

	if !strings.Contains(out, "waiting 10m0s") {
		t.Fatalf("age did not accumulate across flushes at the destination gate:\n%s", out)
	}
	if !strings.Contains(out, "ack seq=7") {
		t.Fatalf("destination-gate log did not name the packet:\n%s", out)
	}
}

// TestHandleBatch_DeliveryClearsTheAge: once a packet lands, a later packet
// reusing its sequence must not inherit the old age. Sequences do restart —
// create-clients-eth repoints an existing router client with migrateClient, so
// the source begins again at sequence 1 under the SAME client id.
//
// Two independent properties keep that honest, and this asserts both:
//   - a delivered packet clears its own entry (settleDelivered), and
//   - a sequence reused at a different source height starts its own clock rather
//     than inheriting the retired packet's age (see TestReusedSequenceRestartsTheAge).
func TestHandleBatch_DeliveryClearsTheAge(t *testing.T) {
	src := &mockSource{latest: 100, relayable: 50}
	m := NewModule("arbitrum->cosmos", "client-0", src, &mockDest{}, &mockBuilder{})

	clk := &fakeClock{t: time.Unix(1_700_000_000, 0)}
	m.waits.now = clk.now

	// Generation 1, seq 7 at height 80: above the relayable frontier, so it parks.
	old := chain.Event{Type: chain.AckPacket, Sequence: 7, Height: 80, Raw: []byte("ack-old"), AckBytes: [][]byte{{1}}}
	m.handleBatch(context.Background(), []chain.Event{old})
	clk.t = clk.t.Add(30 * time.Minute)

	oldKey := waitKey{typ: chain.AckPacket, seq: 7}
	if _, ok := m.waits.entries[oldKey]; !ok {
		t.Fatal("the parked packet should have a wait entry to inherit from")
	}

	// Generation 2 after a client re-creation: the same sequence, but emitted at a
	// height of its own, and low enough to relay immediately.
	fresh := chain.Event{Type: chain.AckPacket, Sequence: 7, Height: 40, Raw: []byte("ack-new"), AckBytes: [][]byte{{1}}}
	m.handleBatch(context.Background(), []chain.Event{fresh})

	if _, ok := m.waits.entries[oldKey]; ok {
		t.Fatal("a delivered packet must not leave a wait entry behind")
	}
}

// TestReusedSequenceRestartsTheAge: sequences restart when create-clients-eth
// repoints a router client with migrateClient — the client id is unchanged, so
// only the source height tells the two generations apart. The later packet must
// start its own clock instead of reporting the retired one's age.
func TestReusedSequenceRestartsTheAge(t *testing.T) {
	w, clk := newTestTracker()

	retired := ackEvent(1, 21)
	w.observe(retired)
	clk.t = clk.t.Add(45 * time.Minute)
	if age := w.observe(retired); age != 45*time.Minute {
		t.Fatalf("age of the still-waiting packet = %s, want 45m0s", age)
	}

	// Same sequence after the re-creation, emitted at a height of its own.
	fresh := retired
	fresh.Height = retired.Height + 500
	if age := w.observe(fresh); age != 0 {
		t.Fatalf("reused sequence reported age %s; it inherited the retired packet's clock", age)
	}

	// And it now ages on its own.
	clk.t = clk.t.Add(2 * time.Minute)
	if age := w.observe(fresh); age != 2*time.Minute {
		t.Fatalf("age after 2m = %s, want 2m0s", age)
	}
}
