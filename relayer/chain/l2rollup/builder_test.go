package l2rollup

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"testing"

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
	// Payload must be the tagged ClientMessage envelope carrying the assembled header.
	var env clientMessageEnvelope
	if err := json.Unmarshal(upd.Payload, &env); err != nil {
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
