package services

import (
	"errors"
	"log"
	"time"

	client "relayer/client"
)

type timeoutSendOutcome int

const (
	timeoutSent timeoutSendOutcome = iota
	timeoutNotDue
	timeoutDeferred
	timeoutFailed
	timeoutAlreadyReceived
)

func (o timeoutSendOutcome) packetIsDone() bool {
	return o == timeoutSent || o == timeoutAlreadyReceived
}

// timeoutOutcomeForError charges only an error that proves this packet's
// transaction was included and deterministically reverted. RPC, broadcast,
// proof-endpoint, and receipt-wait failures are infrastructure deferrals.
func timeoutOutcomeForError(err error) timeoutSendOutcome {
	if errors.Is(err, client.ErrPacketAlreadyReceived) {
		return timeoutAlreadyReceived
	}
	if errors.Is(err, ErrPermanentRelayFailure) {
		return timeoutFailed
	}
	return timeoutDeferred
}

func deferTimeoutRetries(tracker *PendingPacketTracker, pending []pendingPacketInfo, now time.Time, tag string) {
	deferralCounts := tracker.DeferTimeoutRetriesIfCurrent(pending, now)
	for i, info := range pending {
		deferrals := deferralCounts[i]
		if deferralIsStuck(deferrals) {
			log.Printf("[%s][STUCK] seq=%d has deferred timeout %d times; funds remain escrowed", tag, info.Packet.Sequence, deferrals)
		}
	}
}

func chargeTimeoutFailure(tracker *PendingPacketTracker, info pendingPacketInfo, now time.Time, tag string) {
	if tracker.RecordTimeoutFailureIfCurrent(info, now) {
		log.Printf("[%s][ATTENTION] seq=%d exhausted timeout attempts; escrowed funds require operator action", tag, info.Packet.Sequence)
	}
}

func applyTimeoutOutcome(tracker *PendingPacketTracker, info pendingPacketInfo, outcome timeoutSendOutcome, now time.Time, tag string) {
	switch {
	case outcome.packetIsDone():
		tracker.RemoveIfCurrent(info)
	case outcome == timeoutNotDue:
		// The counterparty clock has not reached the packet deadline. This is not
		// chargeable, but leaving it immediately due makes every scan repeat the
		// shared client update and proof preparation while the clock is lagging.
		deferTimeoutRetries(tracker, []pendingPacketInfo{info}, now, tag)
	case outcome == timeoutDeferred:
		deferTimeoutRetries(tracker, []pendingPacketInfo{info}, now, tag)
	case outcome == timeoutFailed:
		chargeTimeoutFailure(tracker, info, now, tag)
	}
}
