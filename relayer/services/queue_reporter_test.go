package services

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"
)

// newQueueReporterTestServices builds the same Services shape the other
// reporter tests use: no chain clients, because nothing here reaches one.
func newQueueReporterTestServices(t *testing.T) *Services {
	t.Helper()
	return New(nil, nil, DefaultConfig())
}

// enqueueQueueReporterPacket makes one counter move, which is all these tests
// need from a "change".
func enqueueQueueReporterPacket(t *testing.T, svc *Services) {
	t.Helper()
	svc.BatchBuilder.PendingTracker.Add(timeoutRetryPacket(1), 10)
}

// captureLog redirects the standard logger for the duration of f and returns
// what was written. The reporter logs through the package logger, which is the
// surface this behaviour is about, so the test reads it rather than restructuring
// production code to hand back strings.
func captureLog(t *testing.T, f func()) string {
	t.Helper()
	var buf bytes.Buffer
	flags, out := log.Flags(), log.Writer()
	log.SetFlags(0)
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetFlags(flags); log.SetOutput(out) })
	f()
	return buf.String()
}

func countLines(s, substr string) int {
	n := 0
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(line, substr) {
			n++
		}
	}
	return n
}

// An idle relayer printed two [QueueState] lines a minute with nothing to say:
// 748 lines over a six-hour run in which the queue was empty throughout. The
// routine lines now report a CHANGE, plus a heartbeat so an idle relayer still
// proves it is running.
func TestQueueReporterIsQuietWhileNothingChanges(t *testing.T) {
	svc := newQueueReporterTestServices(t)
	now := time.Unix(1_700_000_000, 0)

	out := captureLog(t, func() {
		var quiet uint64
		last := svc.reportQueueState(now, nil, &quiet)
		// One hour of ticks with an empty, unchanging queue.
		for i := 1; i <= 60; i++ {
			last = svc.reportQueueState(now.Add(time.Duration(i)*time.Minute), &last, &quiet)
		}
	})

	// The first call always reports (last == nil): an operator starting the
	// process must see the state it started in.
	// After that: 60 unchanged ticks, one heartbeat per queueReportQuietHeartbeat.
	want := 1 + 60/queueReportQuietHeartbeat
	if got := countLines(out, "[QueueState] queued("); got != want {
		t.Fatalf("queued() lines = %d, want %d (1 initial + %d heartbeats):\n%s",
			got, want, 60/queueReportQuietHeartbeat, out)
	}
	// Before this change the same hour produced 61 of them.
	if got := countLines(out, "[QueueState] queued("); got >= 61 {
		t.Fatalf("the reporter is still printing every tick: %d lines", got)
	}
}

// The client-update ages advance every tick on their own. Including them in the
// comparison is the mistake D6 documents: the key never repeats, so the dedupe
// dedupes nothing. They must print with the routine line and never drive it.
func TestQueueReporterAgesDoNotDefeatTheDedupe(t *testing.T) {
	svc := newQueueReporterTestServices(t)
	now := time.Unix(1_700_000_000, 0)
	svc.ObserveCosmosOnEVMUpdate(now.Add(-time.Hour))
	svc.ObserveEVMOnCosmosUpdate(now.Add(-2 * time.Hour))

	out := captureLog(t, func() {
		var quiet uint64
		last := svc.reportQueueState(now, nil, &quiet)
		for i := 1; i < queueReportQuietHeartbeat; i++ {
			last = svc.reportQueueState(now.Add(time.Duration(i)*time.Minute), &last, &quiet)
		}
	})

	if got := countLines(out, "last client update:"); got != 1 {
		t.Fatalf("client-update line printed %d times; the ages change every tick, so "+
			"driving output from them means never de-duplicating:\n%s", got, out)
	}
	if !strings.Contains(out, "cosmos-on-eth=1h0m0s ago") {
		t.Fatalf("the age must still be printed on the line it accompanies:\n%s", out)
	}
}

// A change must be reported at once. A relayer that goes from idle to busy and
// says nothing for ten minutes is worse than one that repeats itself.
func TestQueueReporterReportsAChangeImmediately(t *testing.T) {
	svc := newQueueReporterTestServices(t)
	now := time.Unix(1_700_000_000, 0)

	out := captureLog(t, func() {
		var quiet uint64
		last := svc.reportQueueState(now, nil, &quiet)
		last = svc.reportQueueState(now.Add(time.Minute), &last, &quiet)
		enqueueQueueReporterPacket(t, svc)
		svc.reportQueueState(now.Add(2*time.Minute), &last, &quiet)
	})

	if got := countLines(out, "[QueueState] queued("); got != 2 {
		t.Fatalf("queued() lines = %d, want 2 (the initial one and the change):\n%s", got, out)
	}
	if !strings.Contains(out, "pending(cosmos=1") {
		t.Fatalf("the change itself was not reported:\n%s", out)
	}
}

// A quiet stretch that follows a change must restart from zero, not inherit the
// old count -- otherwise the heartbeat lands at an arbitrary offset.
func TestQueueReporterResetsTheQuietCountOnAChange(t *testing.T) {
	svc := newQueueReporterTestServices(t)
	now := time.Unix(1_700_000_000, 0)

	var quiet uint64
	last := svc.reportQueueState(now, nil, &quiet)
	for i := 1; i < queueReportQuietHeartbeat; i++ {
		last = svc.reportQueueState(now.Add(time.Duration(i)*time.Minute), &last, &quiet)
	}
	enqueueQueueReporterPacket(t, svc)
	last = svc.reportQueueState(now.Add(time.Hour), &last, &quiet)
	if quiet != 0 {
		t.Fatalf("quiet tick count = %d after a change, want 0", quiet)
	}
}

// The merge with #415 added the paused-retry flags to queueState. They belong in
// counts() by this file's own rule -- they are signal, not something that
// advances on its own -- and leaving them out would make a persistence outage
// ENDING invisible: the ATTENTION line simply stops firing, and with the routine
// lines still suppressed the log would go quiet with nothing saying the state
// directory came back.
func TestQueueReporterTreatsPausedRetriesAsAChange(t *testing.T) {
	base := queueState{cosmosQueued: 1}

	paused := base
	paused.cosmosRetryPaused = true
	if base.counts() == paused.counts() {
		t.Fatal("pausing timeout retries is invisible to the change check; " +
			"neither the outage starting nor it ending would reprint the routine lines")
	}

	// The other half of the same rule, which the merge must not break: an age
	// advances on its own, so it must NOT count, or the dedupe dedupes nothing.
	aged := base
	aged.cosmosTimeoutDeadAge = time.Hour
	if base.counts() != aged.counts() {
		t.Fatal("an age difference reached counts(); every tick would look like a change")
	}
}
