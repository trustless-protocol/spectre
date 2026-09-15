package services

import (
	"errors"
	"fmt"
	"testing"
	"time"

	relayerclient "relayer/client"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// trackerWith returns a tracker already holding one packet. trackedPacket (in
// timeout_retry_test.go) builds the packet itself.
func trackerWith(t *testing.T, seq uint64) (*PendingPacketTracker, channeltypesv2.Packet) {
	t.Helper()
	tracker := NewPendingPacketTracker()
	pkt := trackedPacket(seq)
	tracker.Add(pkt, 1)
	return tracker, pkt
}

// TestDeferTimeoutRetryNeverDeadLetters is the core of the first P1: a packet
// must not be abandoned for an outage it had no part in.
//
// RecordTimeoutFailure was charged on every false return of timeoutEthSend,
// including the two paths that never build or send anything — a failed
// light-client update, and a counterparty whose clock has not reached the
// timeout. A counterparty lagging across the ~90 minutes the backoff sums to
// therefore dead-lettered a healthy packet, and a dead-lettered timeout means
// escrowed funds that nothing will release.
func TestDeferTimeoutRetryNeverDeadLetters(t *testing.T) {
	t.Parallel()

	tracker, pkt := trackerWith(t, 7)
	now := time.Now()

	// Far more deferrals than maxTimeoutAttempts.
	for i := 0; i < maxTimeoutAttempts*3; i++ {
		tracker.deferTimeoutRetryForTest(pkt.SourceClient, pkt.Sequence, now)
	}

	if got := tracker.DeadLetteredTimeouts(); got != 0 {
		t.Fatalf("deferrals dead-lettered %d packet(s); an outage must never abandon a packet", got)
	}
	if tracker.Len() != 1 {
		t.Fatal("the packet must stay tracked so it is retried once the outage clears")
	}
}

// A deferral still backs off — the RLY-06 cost problem is real whoever is at
// fault. Without this the scanner re-runs a client update and a Groth16 proof
// for every tracked packet every 30s for as long as the outage lasts.
func TestDeferTimeoutRetryStillBacksOff(t *testing.T) {
	t.Parallel()

	tracker, pkt := trackerWith(t, 7)
	now := time.Now()

	tracker.deferTimeoutRetryForTest(pkt.SourceClient, pkt.Sequence, now)
	if due := tracker.GetDue(now.Add(30 * time.Second)); len(due) != 0 {
		t.Fatal("a deferred packet must not come due on the next 30s scan")
	}
	if due := tracker.GetDue(now.Add(timeoutRetryBackoffBase + time.Second)); len(due) != 1 {
		t.Fatal("a deferred packet must come due once its backoff elapses")
	}
}

// Deferrals and attempts are counted separately: an outage must not consume the
// budget that exists to catch a genuinely poisoned packet.
func TestDeferralsDoNotConsumeTheAttemptBudget(t *testing.T) {
	t.Parallel()

	tracker, pkt := trackerWith(t, 7)
	now := time.Now()

	for i := 0; i < maxTimeoutAttempts; i++ {
		tracker.deferTimeoutRetryForTest(pkt.SourceClient, pkt.Sequence, now)
	}
	// The full attempt budget must still be available afterwards.
	for i := 0; i < maxTimeoutAttempts; i++ {
		if tracker.recordTimeoutFailureForTest(pkt.SourceClient, pkt.Sequence, now) {
			t.Fatalf("dead-lettered on attributable attempt %d of %d; deferrals ate the budget", i+1, maxTimeoutAttempts)
		}
	}
	if !tracker.recordTimeoutFailureForTest(pkt.SourceClient, pkt.Sequence, now) {
		t.Fatal("the attempt budget must still terminate at maxTimeoutAttempts")
	}
}

// TestUpdateMsgsFailed is the second P1. SendCosmosTxBatch submits
// [updateMsgs..., timeoutMsgs...] and SucceededCount is a prefix length, so
// fewer successes than update messages proves no timeout was ever executed.
func TestUpdateMsgsFailed(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		err            error
		updateMsgCount int
		want           bool
	}{
		{
			name:           "update prefix incomplete: no timeout could have run",
			err:            &BatchPartialError{SucceededCount: 1, Err: errors.New("boom")},
			updateMsgCount: 2,
			want:           true,
		},
		{
			name:           "update prefix complete: the failure is in the timeouts",
			err:            &BatchPartialError{SucceededCount: 2, Err: errors.New("boom")},
			updateMsgCount: 2,
			want:           false,
		},
		{
			name:           "whole tx failed with update msgs present: nothing executed",
			err:            errors.New("broadcast failed"),
			updateMsgCount: 1,
			want:           true,
		},
		{
			name:           "no update msgs: the timeouts are the only suspects",
			err:            errors.New("broadcast failed"),
			updateMsgCount: 0,
			want:           false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := updateMsgsFailed(tc.err, tc.updateMsgCount); got != tc.want {
				t.Fatalf("updateMsgsFailed = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestDeadLetteredTimeoutsAreInspectable covers the RLY-02 half: the set used to
// be countable but not inspectable, so the packet identities behind the count
// existed only in log lines — unlike the batch dead letters, which have
// DrainDeadLetters/OldestDeadLetter.
func TestDeadLetteredTimeoutsAreInspectable(t *testing.T) {
	t.Parallel()

	tracker, pkt := trackerWith(t, 42)
	deadAt := time.Now()
	for i := 0; i <= maxTimeoutAttempts; i++ {
		tracker.recordTimeoutFailureForTest(pkt.SourceClient, pkt.Sequence, deadAt)
	}
	if tracker.DeadLetteredTimeouts() != 1 {
		t.Fatalf("expected the packet to be dead-lettered, got %d", tracker.DeadLetteredTimeouts())
	}

	if age := tracker.OldestDeadLetteredTimeout(deadAt.Add(2 * time.Hour)); age < 2*time.Hour {
		t.Fatalf("age = %s, want at least 2h: a count says a timeout gave up, the age says nobody has looked", age)
	}

	packets := tracker.DeadLetteredTimeoutPackets()
	if len(packets) != 1 || packets[0].Sequence != 42 || packets[0].SourceClient != pkt.SourceClient {
		t.Fatalf("the listing must identify which packet gave up (client and sequence), got %+v", packets)
	}

	if again := tracker.DeadLetteredTimeoutPackets(); len(again) != 1 || again[0].Sequence != pkt.Sequence {
		t.Fatalf("reading must not silence the timeout alarm, got %+v", again)
	}
}

// TestTimeoutOutcomeForErrorDefersTransientFailures is the second half of the
// "only charge what the packet caused" rule, on the per-packet send path rather
// than the shared client-update one.
//
// SendEthTx returns untyped errors for everything before a transaction exists
// (nonce lookup, fee estimation, broadcast) and for a receipt wait that timed
// out after the transaction may already have been accepted. Charging those let
// an RPC outage burn all maxTimeoutAttempts and dead-letter an escrowed packet
// whose timeout transaction was never confirmed — the property this change set
// out to remove, reappearing one layer down.
func TestTimeoutOutcomeForErrorDefersTransientFailures(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		err  error
		want timeoutSendOutcome
	}{
		{"nonce lookup failed", errors.New("failed to get pending nonce"), timeoutDeferred},
		{"broadcast failed", errors.New("dial tcp: connection refused"), timeoutDeferred},
		{"receipt wait timed out", errors.New("transaction wait mined timed out after 5 attempts"), timeoutDeferred},
		{"rate limited", errors.New("429 Too Many Requests"), timeoutDeferred},

		{"on-chain revert is attributable", ErrPermanentRelayFailure, timeoutFailed},
		{"wrapped permanent failure", fmt.Errorf("tx reverted: %w", ErrPermanentRelayFailure), timeoutFailed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := timeoutOutcomeForError(tc.err); got != tc.want {
				t.Fatalf("timeoutOutcomeForError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// An RPC outage lasting the whole retry window must delay a timeout, never
// abandon it — the packet's funds are escrowed and only a timeout releases them.
func TestTransientSendFailuresNeverDeadLetter(t *testing.T) {
	t.Parallel()

	tracker, pkt := trackerWith(t, 7)
	now := time.Now()
	transient := errors.New("dial tcp: connection refused")

	for i := 0; i < maxTimeoutAttempts*3; i++ {
		if timeoutOutcomeForError(transient) == timeoutFailed {
			tracker.recordTimeoutFailureForTest(pkt.SourceClient, pkt.Sequence, now)
		} else {
			tracker.deferTimeoutRetryForTest(pkt.SourceClient, pkt.Sequence, now)
		}
	}

	if got := tracker.DeadLetteredTimeouts(); got != 0 {
		t.Fatalf("an RPC outage dead-lettered %d packet(s); escrowed funds must not be abandoned for someone else's downtime", got)
	}
	if tracker.Len() != 1 {
		t.Fatal("the packet must stay tracked so it is retried once the outage clears")
	}
}

// A proof build fails in several ways, and only one of them is evidence about
// the packet. Getting this split wrong is expensive in both directions: charge
// an endpoint's failure and an escrowed packet is dead-lettered for someone
// else's downtime; defer a packet that was already received and the scanner
// rebuilds a proof that can never succeed, forever, at the 30-minute cap.
func TestTimeoutOutcomeForErrorChargesOnlyPacketAttributableProofFailures(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		err  error
		want timeoutSendOutcome
	}{
		{
			"receipt slot already holds a value",
			fmt.Errorf("failed to get ETH non-membership proof: storage slot not empty at key 0xabc: value=0x1: %w",
				relayerclient.ErrPacketAlreadyReceived),
			// Delivered, not failed: it leaves the tracker without charging.
			timeoutAlreadyReceived,
		},
		{"eth_getProof transport error", errors.New("eth_getProof failed: connection refused"), timeoutDeferred},
		{"answer carried no storage proof", errors.New("eth_getProof returned no storage proofs"), timeoutDeferred},
		{"proof answered at the wrong height", errors.New("proof height mismatch"), timeoutDeferred},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := timeoutOutcomeForError(tc.err); got != tc.want {
				t.Fatalf("timeoutOutcomeForError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// A received packet must never consume the retry budget, on either direction.
//
// This replaces an earlier test that pinned the opposite: the classifier used to
// map ErrPacketAlreadyReceived to timeoutFailed so that any path which did not
// special-case it would at least stop retrying. That fallback charged the packet
// and dead-lettered it, which reports escrowed funds needing rescue for a packet
// whose funds already moved. Both scanners now recognise the sentinel directly,
// so the fallback is gone and charging is simply wrong.
func TestReceivedPacketNeverChargesTheRetryBudget(t *testing.T) {
	t.Parallel()

	tracker, pkt := trackerWith(t, 11)
	now := time.Now()
	received := fmt.Errorf("build timeout proof: %w", relayerclient.ErrPacketAlreadyReceived)

	outcome := timeoutOutcomeForError(received)
	if outcome == timeoutFailed {
		t.Fatal("a received packet must not be charged; dead-lettering it reports stuck funds that are not stuck")
	}
	if !outcome.packetIsDone() {
		t.Fatal("a received packet must leave the tracker rather than be retried")
	}

	// Drive the scanners' action for that outcome many times over: it must never
	// dead-letter, no matter how often it is reached.
	for i := 0; i < maxTimeoutAttempts*2; i++ {
		if outcome.packetIsDone() {
			tracker.removeCurrentForTest(pkt.SourceClient, pkt.Sequence)
			continue
		}
		tracker.recordTimeoutFailureForTest(pkt.SourceClient, pkt.Sequence, now)
	}
	if got := tracker.DeadLetteredTimeouts(); got != 0 {
		t.Fatalf("DeadLetteredTimeouts() = %d; a delivered packet is not an operator problem", got)
	}
	if tracker.Len() != 0 {
		t.Fatal("a received packet must leave the tracker")
	}
}

// A packet whose receipt already exists on the destination is not a failure: it
// was received, its funds moved, and no timeout for it can ever be valid. It
// must leave the tracker quietly.
//
// Dead-lettering it instead — which is what the classifier alone would do —
// would count it under "escrowed funds cannot be refunded without operator
// action", sending someone to investigate a packet that completed normally.
func TestReceivedPacketIsDroppedNotDeadLettered(t *testing.T) {
	t.Parallel()

	tracker, pkt := trackerWith(t, 21)
	received := fmt.Errorf("build timeout proof: %w", relayerclient.ErrPacketAlreadyReceived)

	// The action scanForCosmosTimeouts takes on this error.
	if errors.Is(received, relayerclient.ErrPacketAlreadyReceived) {
		tracker.removeCurrentForTest(pkt.SourceClient, pkt.Sequence)
	}

	if got := tracker.Len(); got != 0 {
		t.Fatalf("tracker still holds %d packet(s); a received packet must not be retried", got)
	}
	if got := tracker.DeadLetteredTimeouts(); got != 0 {
		t.Fatalf("DeadLetteredTimeouts() = %d; a received packet is not an operator problem", got)
	}
}

// Deferrals shape the backoff and feed the stuck gauge, so a packet that
// recovers must not keep paying for an outage that has passed.
func TestClearDeferralsResetsAfterRecovery(t *testing.T) {
	t.Parallel()

	tracker, pkt := trackerWith(t, 23)
	now := time.Now()
	for i := 0; i < deferralAlertThreshold; i++ {
		tracker.deferTimeoutRetryForTest(pkt.SourceClient, pkt.Sequence, now)
	}
	if count, worst := tracker.StuckTimeouts(); count != 1 || worst != deferralAlertThreshold {
		t.Fatalf("StuckTimeouts() = (%d, %d), want (1, %d)", count, worst, deferralAlertThreshold)
	}

	tracker.clearDeferralsForTest(pkt.SourceClient, pkt.Sequence)

	if count, worst := tracker.StuckTimeouts(); count != 0 || worst != 0 {
		t.Fatalf("StuckTimeouts() = (%d, %d) after recovery, want (0, 0)", count, worst)
	}
	if got := tracker.deferTimeoutRetryForTest(pkt.SourceClient, pkt.Sequence, now); got != 1 {
		t.Fatalf("first deferral after a reset reported %d, want 1", got)
	}
	if tracker.Len() != 1 {
		t.Fatal("clearing deferrals must not drop the packet")
	}
}

// The stuck gauge exists because deferrals never dead-letter: without it a
// packet at the backoff cap reads as zero on every other counter.
func TestStuckTimeoutsCountsOnlyPastTheThreshold(t *testing.T) {
	t.Parallel()

	tracker, pkt := trackerWith(t, 29)
	now := time.Now()

	for i := 1; i < deferralAlertThreshold; i++ {
		tracker.deferTimeoutRetryForTest(pkt.SourceClient, pkt.Sequence, now)
		if count, _ := tracker.StuckTimeouts(); count != 0 {
			t.Fatalf("deferral %d already counted as stuck", i)
		}
	}
	tracker.deferTimeoutRetryForTest(pkt.SourceClient, pkt.Sequence, now)
	if count, worst := tracker.StuckTimeouts(); count != 1 || worst != deferralAlertThreshold {
		t.Fatalf("StuckTimeouts() = (%d, %d), want (1, %d)", count, worst, deferralAlertThreshold)
	}
}

// The ETH direction's mirror of TestReceivedPacketIsDroppedNotDeadLettered.
//
// scanForEthTimeouts checks the ETH commitment, then builds a Cosmos
// non-membership proof. When the receipt landed between those two steps the
// proof cannot be built -- and that is not a failure, it is delivery. Reading it
// as a generic proof error defers the packet at the backoff cap forever and ends
// in a STUCK alert about a packet whose funds moved normally.
//
// This was a one-sided fix: the Cosmos scanner recognised the sentinel and the
// ETH scanner did not.
func TestEthTimeoutTreatsAReceivedPacketAsDone(t *testing.T) {
	t.Parallel()

	received := fmt.Errorf("build non-membership proof: %w", relayerclient.ErrPacketAlreadyReceived)

	// Built by the production path, not hand-assembled: this is what catches the
	// sentinel being dropped from cosmosNonMembership's receipt check.
	if !errors.Is(cosmosReceiptPresentError(1234, 32), relayerclient.ErrPacketAlreadyReceived) {
		t.Fatal("a Cosmos receipt at the non-membership path must carry ErrPacketAlreadyReceived; " +
			"without it the packet is deferred forever and raises a false STUCK alert")
	}

	got := timeoutOutcomeForError(received)
	if got != timeoutAlreadyReceived {
		t.Fatalf("timeoutOutcomeForError(already-received) = %v, want timeoutAlreadyReceived", got)
	}
	if !got.packetIsDone() {
		t.Fatal("an already-received packet must leave the tracker, not be retried")
	}
	if got == timeoutFailed {
		t.Fatal("an already-received packet must not charge the retry budget")
	}

	// The outcomes that are genuinely unfinished must not claim to be done.
	for _, o := range []timeoutSendOutcome{timeoutNotDue, timeoutDeferred, timeoutFailed} {
		if o.packetIsDone() {
			t.Fatalf("outcome %v reported done; the packet still needs timeout work", o)
		}
	}
	if !timeoutSent.packetIsDone() {
		t.Fatal("a landed timeout must be done")
	}
}
