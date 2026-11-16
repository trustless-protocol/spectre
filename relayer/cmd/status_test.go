package main

import (
	"bytes"
	"os"
	"path/filepath"
	"relayer/relay"
	"relayer/services"
	"relayer/status"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// What renderStatus prints is what an operator reads under pressure, so it is a
// pure function of the snapshot and covered like one.
func TestRenderStatus(t *testing.T) {
	now := time.Date(2026, 9, 13, 10, 0, 0, 0, time.UTC)
	full := status.Snapshot{
		Now: now,
		Paths: []relay.Status{{
			Path: "cosmos->eth", ClientID: "07-tendermint-0",
			Source: "cosmos", Dest: "ethereum",
			TrustedHeight: 41, RelayableHeight: 44,
			Waiting: []relay.WaitingPacket{
				{Type: "ack", Sequence: 9, Height: 40, AgeText: "12m0s"},
			},
		}},
		Queues: []services.QueueStatus{
			{
				Direction: "cosmos", Queued: 2,
				Pending:      []services.PendingEntry{{SourceClient: "07-tendermint-0", Sequence: 11, AgeText: "3m0s", TimeoutAttempts: 2}},
				DeadLettered: []services.DeadLetterEntry{{SourceClient: "07-tendermint-0", Sequence: 7, DeadLetteredAt: now.Add(-2 * time.Hour), Attempts: 5}},
			},
			{Direction: "l2"}, // idle
		},
	}

	t.Run("shows both halves of the stuck question", func(t *testing.T) {
		out := render(t, full)
		// The waiting set: relayed nothing yet, and here is how long.
		if !strings.Contains(out, "ack seq=9") || !strings.Contains(out, "12m0s") {
			t.Errorf("the waiting packet and its age are missing:\n%s", out)
		}
		// The pending set: relayed, awaiting settlement.
		if !strings.Contains(out, "pending 07-tendermint-0 seq=11") {
			t.Errorf("the pending packet is missing:\n%s", out)
		}
		if !strings.Contains(out, "relayable height") || !strings.Contains(out, "44") {
			t.Errorf("the relayable height is missing; it is the number that answers 'why has nothing moved':\n%s", out)
		}
	})

	// Escrowed funds waiting on a person must not read like a routine line.
	t.Run("a dead letter is loud and dated", func(t *testing.T) {
		out := render(t, full)
		if !strings.Contains(out, "DEAD-LETTERED") {
			t.Errorf("dead-lettered packets are not called out:\n%s", out)
		}
		if !strings.Contains(out, "2h0m0s ago") {
			t.Errorf("no age on the dead letter; a count cannot tell 'a minute ago' from 'since Tuesday':\n%s", out)
		}
		if !strings.Contains(out, "5 attempt(s)") {
			t.Errorf("no attempt count; it is the difference between a deterministic failure and an outage:\n%s", out)
		}
	})

	t.Run("a paused tracker is called out", func(t *testing.T) {
		paused := status.Snapshot{Now: now, Queues: []services.QueueStatus{
			{Direction: "eth", RetryPaused: true, RetryPausedReason: "disk full"},
		}}
		out := render(t, paused)
		if !strings.Contains(out, "ATTENTION") || !strings.Contains(out, "disk full") {
			t.Errorf("a tracker that stopped doing timeout work reads like a normal line:\n%s", out)
		}
	})

	// A dropped ack or timeout has no other exit: its escrow stays locked until a
	// retry lands, so it must be named even when nothing is waiting.
	t.Run("a packet held after a permanent failure is named", func(t *testing.T) {
		dropped := status.Snapshot{Now: now, Paths: []relay.Status{{
			Path: "cosmos->eth", ClientID: "07-tendermint-0", Source: "cosmos", Dest: "ethereum",
			Dropped: []relay.DroppedPacket{{
				Type: "ack", Sequence: 7, Since: now.Add(-time.Hour), Drops: 3,
				NextRetry: now.Add(40 * time.Minute), Reason: "ackPacket seq=7 would revert",
			}},
		}}}
		out := render(t, dropped)
		for _, want := range []string{"ATTENTION", "ack seq=7", "3 drop(s)", "2026-09-13T10:40:00Z", "would revert"} {
			if !strings.Contains(out, want) {
				t.Errorf("dropped packet line is missing %q:\n%s", want, out)
			}
		}
		if strings.Contains(render(t, full), "dropped after a permanent failure") {
			t.Errorf("a path with nothing dropped printed the dropped section")
		}
	})

	// An incident is the wrong time to scroll past idle trackers.
	t.Run("idle trackers are left out", func(t *testing.T) {
		out := render(t, full)
		if strings.Contains(out, "l2 queue") {
			t.Errorf("an idle tracker was printed:\n%s", out)
		}
	})

	t.Run("says so when the process drives nothing", func(t *testing.T) {
		out := render(t, status.Snapshot{Now: now})
		if !strings.Contains(out, "no relay path") {
			t.Errorf("an empty snapshot printed nothing an operator can act on:\n%s", out)
		}
	})
}

func render(t *testing.T, snap status.Snapshot) string {
	t.Helper()
	var buf bytes.Buffer
	if err := renderStatus(&buf, snap); err != nil {
		t.Fatalf("renderStatus: %v", err)
	}
	return buf.String()
}

// The config line is owed only when the file is read. --addr skips the file, so
// a line naming it there would be the wrong-file misdirection the line exists to
// prevent.
func TestStatusTarget_NamesTheConfigOnlyWhenItReadsOne(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"status_address":"127.0.0.1:9"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Run("from the config", func(t *testing.T) {
		t.Parallel()
		core, logs := observer.New(zapcore.InfoLevel)
		target, err := statusTarget(zap.New(core), "status", "", path)
		if err != nil {
			t.Fatal(err)
		}
		if target != "127.0.0.1:9" {
			t.Fatalf("target = %q, want the config's status_address", target)
		}
		if n := logs.FilterMessageSnippet(path).Len(); n != 1 {
			t.Fatalf("config named %d time(s), want 1: %v", n, logs.All())
		}
	})

	t.Run("from --addr", func(t *testing.T) {
		t.Parallel()
		core, logs := observer.New(zapcore.InfoLevel)
		target, err := statusTarget(zap.New(core), "status", "127.0.0.1:8", path)
		if err != nil {
			t.Fatal(err)
		}
		if target != "127.0.0.1:8" {
			t.Fatalf("target = %q, want the flag", target)
		}
		if n := logs.Len(); n != 0 {
			t.Fatalf("logged %d line(s) without reading a config: %v", n, logs.All())
		}
	})
}
