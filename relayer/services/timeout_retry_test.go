package services

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func trackedPacket(seq uint64) channeltypesv2.Packet {
	return channeltypesv2.Packet{SourceClient: "client-0", Sequence: seq}
}

// TestNextTimeoutBackoffGrowsAndCaps pins the schedule: 1m, 2m, 4m … capped at
// 30m. Without it a deterministically-failing timeout was re-attempted on the raw
// 30s scan cadence, and every attempt costs a Groth16 proof and a transaction.
func TestNextTimeoutBackoffGrowsAndCaps(t *testing.T) {
	want := []time.Duration{
		1 * time.Minute, 2 * time.Minute, 4 * time.Minute, 8 * time.Minute,
		16 * time.Minute, 30 * time.Minute, 30 * time.Minute,
	}
	for i, w := range want {
		if got := nextTimeoutBackoff(i + 1); got != w {
			t.Fatalf("attempt %d: backoff = %s, want %s", i+1, got, w)
		}
	}
	// A very large attempt must still clamp rather than overflow the shift.
	if got := nextTimeoutBackoff(64); got != timeoutRetryBackoffMax {
		t.Fatalf("attempt 64: backoff = %s, want the %s cap (shift overflow guard)", got, timeoutRetryBackoffMax)
	}
}

// A failed timeout submission must hold the packet out of the next scans, so the
// scanner stops re-proving it every 30 seconds.
func TestRecordTimeoutFailureHoldsPacketOutOfScan(t *testing.T) {
	tr := NewPendingPacketTracker()
	tr.Add(trackedPacket(1), 100)

	now := time.Now()
	if got := len(tr.GetDue(now)); got != 1 {
		t.Fatalf("a fresh packet must be due, got %d", got)
	}

	if dead := tr.recordTimeoutFailureForTest("client-0", 1, now); dead {
		t.Fatal("one failure must not dead-letter")
	}
	if got := len(tr.GetDue(now)); got != 0 {
		t.Fatalf("packet under backoff must not be re-scanned, got %d due", got)
	}
	// It comes back once the backoff elapses.
	if got := len(tr.GetDue(now.Add(2 * time.Minute))); got != 1 {
		t.Fatalf("packet must be due again after its backoff, got %d", got)
	}
}

// After maxTimeoutAttempts the packet stops being retried and is kept for an
// operator. It must leave the scan set — otherwise the cap achieves nothing.
func TestRecordTimeoutFailureDeadLettersAtCap(t *testing.T) {
	tr := NewPendingPacketTracker()
	tr.Add(trackedPacket(1), 100)

	now := time.Now()
	for i := 1; i <= maxTimeoutAttempts; i++ {
		if dead := tr.recordTimeoutFailureForTest("client-0", 1, now); dead {
			t.Fatalf("attempt %d dead-lettered early (cap is %d)", i, maxTimeoutAttempts)
		}
	}
	if dead := tr.recordTimeoutFailureForTest("client-0", 1, now); !dead {
		t.Fatalf("attempt %d must dead-letter", maxTimeoutAttempts+1)
	}
	if tr.Len() != 0 {
		t.Fatalf("a dead-lettered packet must leave the scan set, %d still tracked", tr.Len())
	}
	if got := tr.DeadLetteredTimeouts(); got != 1 {
		t.Fatalf("DeadLetteredTimeouts = %d, want 1 — escrowed funds must stay visible", got)
	}
}

// Recording against a packet that is no longer tracked (it succeeded, or was
// purged) must be a no-op rather than resurrecting it.
func TestRecordTimeoutFailureIgnoresUnknownPacket(t *testing.T) {
	tr := NewPendingPacketTracker()
	if dead := tr.recordTimeoutFailureForTest("client-0", 99, time.Now()); dead {
		t.Fatal("unknown packet must not be reported as dead-lettered")
	}
	if tr.Len() != 0 || tr.DeadLetteredTimeouts() != 0 {
		t.Fatal("unknown packet must not be added to any set")
	}
}

// TestDeadLettersRemainVisibleUntilRetryRetentionExpires ensures a read-only
// inspection does not clear the gauge while the timeout remains suppressed.
func TestDeadLettersRemainVisibleUntilRetryRetentionExpires(t *testing.T) {
	tracker, packet := trackerWith(t, 1)
	deadAt := time.Now()
	for i := 0; i <= maxTimeoutAttempts; i++ {
		tracker.recordTimeoutFailureForTest(packet.SourceClient, packet.Sequence, deadAt)
	}
	if got := tracker.DeadLetteredTimeouts(); got != 1 {
		t.Fatalf("expected one dead-lettered timeout, got %d", got)
	}
	packets := tracker.DeadLetteredTimeoutPackets()
	if len(packets) != 1 {
		t.Fatalf("dead-lettered packets=%d, want 1", len(packets))
	}
	if tracker.deadLettered[0].DeadLetteredAt.IsZero() {
		t.Fatal("DeadLetteredAt must be stamped so the age is observable")
	}
	if got := tracker.DeadLetteredTimeouts(); got != 1 {
		t.Fatalf("inspection must not clear the timeout dead-letter gauge, got %d", got)
	}
}

// OldestDeadLetteredTimeout reports zero with nothing dead-lettered and a real
// age once something is — a count alone cannot distinguish a failure nobody has
// looked at from one that just happened.
func TestOldestDeadLetterReportsAge(t *testing.T) {
	tracker := NewPendingPacketTracker()
	if age := tracker.OldestDeadLetteredTimeout(time.Now()); age != 0 {
		t.Fatalf("empty queue must report zero age, got %s", age)
	}

	packet := trackedPacket(1)
	tracker.Add(packet, 1)
	deadAt := time.Now()
	for i := 0; i <= maxTimeoutAttempts; i++ {
		tracker.recordTimeoutFailureForTest(packet.SourceClient, packet.Sequence, deadAt)
	}
	if age := tracker.OldestDeadLetteredTimeout(deadAt.Add(time.Hour)); age < time.Hour {
		t.Fatalf("timeout dead-letter age = %s, want at least an hour", age)
	}
}

// A sequence is not an identity. packetIdentity pairs it with the source client
// precisely because a sequence is reused across clients, and after A1 a merged
// operator log carries several of them.
//
// The two lines that matter most are the ones telling an operator to act --
// STUCK and ATTENTION, both of which say funds are escrowed. Naming a sequence
// there without its client tells them something is wrong and not which packet
// it is. The rest of the scanner's lines are the trail leading to those two, so
// the rule covers the whole family rather than the two endpoints.
//
// Structural rather than behavioural: asserting on captured log output would pin
// the wording, and the invariant is about what the line identifies, not how it
// reads.
func TestTimeoutScanLogsIdentifyThePacketNotJustTheSequence(t *testing.T) {
	var offenders []string
	for _, file := range []string{"services.go", "timeout_retry.go", "pending.go", "batch.go"} {
		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		for i, line := range strings.Split(string(src), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") || !strings.Contains(line, "TimeoutScan]") {
				continue
			}
			if !strings.Contains(line, "seq=%d") || strings.Contains(line, "src=%s") {
				continue
			}
			offenders = append(offenders, fmt.Sprintf("%s:%d: %s", file, i+1, trimmed))
		}
	}
	if len(offenders) > 0 {
		t.Fatalf("timeout-scan log lines naming a sequence with no source client:\n%s\n\n"+
			"An operator reading a merged log cannot tell which client's packet this is.",
			strings.Join(offenders, "\n"))
	}
}

// A tag is interpolated into "[%sTimeoutScan]", so the value passed must be the
// CHAIN, not a name that already ends in Timeout.
//
// Found while merging main into #436: services.go passed "CosmosTimeout", which
// rendered as [CosmosTimeoutTimeoutScan] -- a fourth name for the scanner this
// branch exists to give one name. The label gate in cmd cannot see it, because
// the format string there is correct and only the argument is wrong.
func TestTimeoutScanTagsAreChainNamesNotWorkNames(t *testing.T) {
	src, err := os.ReadFile("services.go")
	if err != nil {
		t.Fatalf("read services.go: %v", err)
	}
	tagRe := regexp.MustCompile(`(?:tag:\s*|chargeTimeoutFailure\([^"]*)"([A-Za-z0-9]+)"`)
	for i, line := range strings.Split(string(src), "\n") {
		for _, m := range tagRe.FindAllStringSubmatch(line, -1) {
			if strings.Contains(m[1], "Timeout") || strings.Contains(m[1], "Scan") {
				t.Errorf("services.go:%d: tag %q becomes [%sTimeoutScan]; pass the chain alone",
					i+1, m[1], m[1])
			}
		}
	}
}
