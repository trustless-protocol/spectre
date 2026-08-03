package l2rollup

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"testing"
	"time"

	"relayer/chain"
)

// mockHeaderBuilder is a fake HeaderBuilder so the generic Builder is testable
// without real L1/L2 RPC.
type mockHeaderBuilder struct {
	msg       ClientMessage
	err       error
	last      HeaderRequest // records the request BuildHeader was called with
	committed uint64        // committed height to report (0 → echo the requested height)
}

func (m *mockHeaderBuilder) Name() string { return "l2-mock" }
func (m *mockHeaderBuilder) BuildHeader(_ context.Context, request HeaderRequest) (ClientMessage, uint64, error) {
	m.last = request
	if m.err != nil {
		return nil, 0, m.err
	}
	msg := m.msg
	if msg == nil {
		msg = &OpStackHeader{BeaconSlot: 999, GameRuntime: byteList{0xde, 0xad}}
	}
	committed := m.committed
	if committed == 0 {
		committed = request.Height
	}
	return msg, committed, nil
}

func headerRequest(kind HeadKind, h uint64) []byte {
	b := make([]byte, 9)
	b[0] = byte(kind)
	binary.BigEndian.PutUint64(b[1:], h)
	return b
}

func TestBuild_EncodesHeaderAsPayload(t *testing.T) {
	m := &mockHeaderBuilder{}
	b := NewBuilder(m)

	upd, err := b.Build(context.Background(), headerRequest(Safe, 4242))
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if m.last != (HeaderRequest{Height: 4242, Finality: Safe}) {
		t.Fatalf("assembled request = %+v", m.last)
	}
	if upd.Height != 4242 {
		t.Fatalf("update height = %d, want 4242", upd.Height)
	}
	if len(upd.Payloads) != 1 {
		t.Fatalf("payload count = %d, want 1", len(upd.Payloads))
	}
	// The sole payload must be the tagged ClientMessage envelope carrying the assembled header.
	var env clientMessageEnvelope
	if err := json.Unmarshal(upd.Payloads[0], &env); err != nil {
		t.Fatalf("payload is not a client message envelope: %v", err)
	}
	if env.Type != "header" {
		t.Fatalf("envelope type = %q, want header", env.Type)
	}
	var got OpStackHeader
	if err := json.Unmarshal(env.Value, &got); err != nil {
		t.Fatalf("envelope value is not an OP header: %v", err)
	}
	if got.BeaconSlot != 999 {
		t.Fatalf("payload did not carry the assembled header: %+v", got)
	}
}

// TestBuild_AdvancesByCommittedHeight guards the fix for the blocker where the
// ClientUpdate reported the requested height instead of the height the header commits:
// a BoLD/legacy/OP builder can return an attested block below the request, and the
// client must advance by THAT, or proofs are built at an unproven height.
func TestBuild_AdvancesByCommittedHeight(t *testing.T) {
	m := &mockHeaderBuilder{committed: 4200} // requested 4242, committed 4200
	upd, err := NewBuilder(m).Build(context.Background(), headerRequest(Safe, 4242))
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if m.last != (HeaderRequest{Height: 4242, Finality: Safe}) {
		t.Fatalf("assembled request = %+v", m.last)
	}
	if upd.Height != 4200 {
		t.Fatalf("update height = %d, want the committed 4200 (not the requested 4242)", upd.Height)
	}
}

func TestBuild_MalformedHeaderRetryable(t *testing.T) {
	b := NewBuilder(&mockHeaderBuilder{})
	_, err := b.Build(context.Background(), []byte{1, 2, 3})
	if err == nil || !chain.IsRetryable(err) {
		t.Fatalf("malformed header must be retryable, got %v", err)
	}
}

func TestBuild_BuildHeaderErrorRetryable(t *testing.T) {
	m := &mockHeaderBuilder{err: context.DeadlineExceeded}
	_, err := b(m).Build(context.Background(), headerRequest(Finalized, 1))
	if err == nil || !chain.IsRetryable(err) {
		t.Fatalf("assemble failure must be retryable, got %v", err)
	}
}

func TestHeaderRequestRejectsUnknownFinality(t *testing.T) {
	_, err := b(&mockHeaderBuilder{}).Build(context.Background(), headerRequest(HeadKind(9), 1))
	if err == nil || !chain.IsRetryable(err) {
		t.Fatalf("unknown finality must be rejected, got %v", err)
	}
}

func TestSourceQueryHeaderCarriesEveryFinalityLevel(t *testing.T) {
	for _, kind := range []HeadKind{Unsafe, Safe, Finalized} {
		source := &Source{headKind: kind}
		raw, err := source.QueryHeader(context.Background(), 77)
		if err != nil {
			t.Fatalf("query %s header: %v", kind, err)
		}
		request, err := decodeHeaderRequest(raw)
		if err != nil {
			t.Fatalf("decode %s header: %v", kind, err)
		}
		if request != (HeaderRequest{Height: 77, Finality: kind}) {
			t.Fatalf("request = %+v, want height 77 / %s", request, kind)
		}
	}

	if _, err := (&Source{headKind: HeadKind(9)}).QueryHeader(context.Background(), 77); err == nil {
		t.Fatal("unknown source finality must be rejected")
	}
}

func b(a HeaderBuilder) *Builder { return NewBuilder(a) }

// blockingHeaderBuilder blocks until its context is done, standing in for an L1/L2
// node that accepts the connection and never answers.
type blockingHeaderBuilder struct{ observed error }

func (b *blockingHeaderBuilder) Name() string { return "l2-blocking" }

func (b *blockingHeaderBuilder) BuildHeader(ctx context.Context, _ HeaderRequest) (ClientMessage, uint64, error) {
	<-ctx.Done()
	b.observed = ctx.Err()
	return nil, 0, ctx.Err()
}

// A hung RPC must surface as a retryable error rather than wedging the caller. The
// relay module drives one direction on a single goroutine, so a build that never
// returns stops that direction silently — no error, no retry, nothing logged.
func TestBuild_HungHeaderBuilderTimesOutRetryable(t *testing.T) {
	mock := &blockingHeaderBuilder{}
	b := NewBuilder(mock)

	// Cancelling the caller's context stands in for the deadline the Builder applies;
	// asserting on the real 90s timeout would make this test take 90s.
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := b.Build(ctx, headerRequest(Safe, 42))
		done <- err
	}()

	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Build returned nil error for a build that never completed")
		}
		if !chain.IsRetryable(err) {
			t.Fatalf("a timed-out proof assembly must be retryable, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Build did not return after its context was cancelled — the direction would wedge")
	}
}

// The Builder must impose a deadline of its own: the relay context is only cancelled
// at shutdown, so it cannot bound a hang during normal operation.
func TestBuild_AppliesItsOwnDeadline(t *testing.T) {
	var deadlineSeen bool
	probe := headerBuilderFunc(func(ctx context.Context, _ HeaderRequest) (ClientMessage, uint64, error) {
		_, deadlineSeen = ctx.Deadline()
		return &OpStackHeader{BeaconSlot: 1}, 42, nil
	})
	if _, err := NewBuilder(probe).Build(context.Background(), headerRequest(Safe, 42)); err != nil {
		t.Fatalf("Build: %v", err)
	}
	if !deadlineSeen {
		t.Fatal("BuildHeader received a context with no deadline; a hung RPC would never be cut off")
	}
}

// headerBuilderFunc adapts a function to HeaderBuilder.
type headerBuilderFunc func(context.Context, HeaderRequest) (ClientMessage, uint64, error)

func (f headerBuilderFunc) Name() string { return "l2-probe" }
func (f headerBuilderFunc) BuildHeader(ctx context.Context, r HeaderRequest) (ClientMessage, uint64, error) {
	return f(ctx, r)
}
