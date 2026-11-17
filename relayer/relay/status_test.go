package relay

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"relayer/chain"
	"relayer/services"
)

// errSource fails RelayableHeight the way an RPC client does: the error quotes
// the endpoint, API key included. It records whether the probe's context was
// already cancelled when it was called.
type errSource struct {
	mockSource
	err            error
	ctxErrOnCall   error
	permanentClass bool
}

func (s *errSource) RelayableHeight(ctx context.Context) (uint64, error) {
	s.ctxErrOnCall = ctx.Err()
	if s.permanentClass {
		return 0, chain.Permanent(s.err)
	}
	return 0, s.err
}

// #575 L7: GET /status returned the raw RPC error, which quotes the provider URL
// with its API key, and a requester could force one by hanging up.
func TestStatus_ProbeErrorIsClassifiedAndRedacted(t *testing.T) {
	const key = "s3cr3t-api-key"
	rpcErr := errors.New(`Post "https://eth-mainnet.example.com/v2/` + key + `": dial tcp: connection refused`)

	t.Run("transient", func(t *testing.T) {
		m := NewModule("cosmos->eth", "client-0", &errSource{err: rpcErr}, &mockDest{}, &mockBuilder{})
		got := m.Status(context.Background()).RelayableHeightErr
		if strings.Contains(got, key) {
			t.Fatalf("status error carries the API key: %q", got)
		}
		if !strings.HasPrefix(got, "transient: ") || !strings.Contains(got, "https://eth-mainnet.example.com") {
			t.Fatalf("status error = %q, want the class and the endpoint host", got)
		}
	})

	t.Run("permanent", func(t *testing.T) {
		m := NewModule("l2->cosmos", "client-0", &errSource{err: rpcErr, permanentClass: true}, &mockDest{}, &mockBuilder{})
		if got := m.Status(context.Background()).RelayableHeightErr; !strings.HasPrefix(got, "permanent: ") || strings.Contains(got, key) {
			t.Fatalf("status error = %q, want a redacted permanent failure", got)
		}
	})

	t.Run("a requester that hangs up does not cancel the probe", func(t *testing.T) {
		src := &errSource{err: rpcErr}
		m := NewModule("cosmos->eth", "client-0", src, &mockDest{}, &mockBuilder{})
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // the request context of a client that closed its connection
		m.Status(ctx)
		if src.ctxErrOnCall != nil {
			t.Fatalf("the probe ran on a cancelled context (%v); a requester could force an error on demand", src.ctxErrOnCall)
		}
	})
}

// #575 L7, through two more fields /status serves: a drop reason and a loop
// note quote the RPC error behind them, API key included. The ledger file kept
// the drop reason too.
func TestStatus_DropReasonsAndLoopNotesAreRedacted(t *testing.T) {
	const key = "s3cr3t-api-key"
	rpcErr := errors.New(`Post "https://eth-mainnet.example.com/v2/` + key + `": dial tcp: connection refused`)

	t.Run("drop reason", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), services.DropLedgerFileName("cosmos-to-eth"))
		ledger, err := services.NewDropLedger(path)
		if err != nil {
			t.Fatal(err)
		}
		m := NewModule("cosmos->eth", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{}, WithDropLedger(ledger))
		// The unresolved-nonce drop quotes the failed cancel broadcast.
		m.recordDrop(chain.SendPacket, 3, []byte("send"), chain.Permanent(fmt.Errorf("nonce unresolved; cancel: %w", rpcErr)))

		got := m.Status(context.Background()).Dropped
		if len(got) != 1 || strings.Contains(got[0].Reason, key) || !strings.Contains(got[0].Reason, "https://eth-mainnet.example.com") {
			t.Fatalf("status.Dropped = %+v, want one reason with the host and no API key", got)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), key) {
			t.Fatalf("the ledger file keeps the API key:\n%s", data)
		}
	})

	t.Run("loop note", func(t *testing.T) {
		src := &listerSource{mockSource: mockSource{latest: 100, relayable: 100}, listErr: rpcErr}
		m := NewModule("eth->cosmos", "client-0", src, &mockDest{}, &mockBuilder{})
		m.flushOnce(context.Background(), src)
		for _, loop := range m.Status(context.Background()).Loops {
			if loop.Worker != workerPacketFlush {
				continue
			}
			if strings.Contains(loop.Note, key) || !strings.Contains(loop.Note, "https://eth-mainnet.example.com") {
				t.Fatalf("flush note = %q, want the host and no API key", loop.Note)
			}
			return
		}
		t.Fatal("the failed flush recorded no loop note")
	})
}

// #575: the first periodic update is never later than one interval, whatever
// the caller derived its delay from.
func TestWithPeriodicUpdate_ClampsTheFirstDelay(t *testing.T) {
	noop := func(context.Context) error { return nil }
	for _, tc := range []struct {
		delay, want time.Duration
	}{
		{delay: 48 * time.Hour, want: time.Hour},
		{delay: -time.Minute, want: 0},
		{delay: 10 * time.Minute, want: 10 * time.Minute},
	} {
		m := NewModule("cosmos->eth", "client-0", &mockSource{}, &mockDest{}, &mockBuilder{},
			WithPeriodicUpdate(time.Hour, tc.delay, noop))
		if m.periodicUpdateInitialDelay != tc.want {
			t.Errorf("initial delay %s -> %s, want %s", tc.delay, m.periodicUpdateInitialDelay, tc.want)
		}
	}
}
