package relay

import (
	"context"
	"errors"
	"strings"
	"testing"

	"relayer/chain"
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
