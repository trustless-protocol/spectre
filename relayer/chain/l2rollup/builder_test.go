package l2rollup

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"strings"
	"testing"

	"relayer/chain"
)

// mockHeaderBuilder is a fake HeaderBuilder so the generic Builder is testable
// without real L1/L2 RPC.
type mockHeaderBuilder struct {
	data *ClientMessageData
	err  error
	last uint64 // records the height BuildHeader was called with
}

func (m *mockHeaderBuilder) Name() string { return "l2-mock" }
func (m *mockHeaderBuilder) BuildHeader(_ context.Context, height uint64) (*ClientMessageData, error) {
	m.last = height
	if m.err != nil {
		return nil, m.err
	}
	d := m.data
	if d == nil {
		d = &ClientMessageData{L1BeaconSlot: 999, L2HeaderRLP: []byte("rlp")}
	}
	return d, nil
}

func height8(h uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, h)
	return b
}

func TestBuild_EncodesHeaderAsPayload(t *testing.T) {
	m := &mockHeaderBuilder{}
	b := NewBuilder(m)

	upd, err := b.Build(context.Background(), height8(4242))
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if m.last != 4242 {
		t.Fatalf("assembled height = %d, want 4242", m.last)
	}
	if upd.Height != 4242 {
		t.Fatalf("update height = %d, want 4242", upd.Height)
	}
	// Payload must be the JSON-encoded client message data.
	var got ClientMessageData
	if err := json.Unmarshal(upd.Payload, &got); err != nil {
		t.Fatalf("payload is not client message JSON: %v", err)
	}
	if got.L1BeaconSlot != 999 {
		t.Fatalf("payload did not carry the assembled data: %+v", got)
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
	_, err := b(m).Build(context.Background(), height8(1))
	if err == nil || !chain.IsRetryable(err) {
		t.Fatalf("assemble failure must be retryable, got %v", err)
	}
}

func b(a HeaderBuilder) *Builder { return NewBuilder(a) }

// TestClientMessageData_JSONShape locks the proposed wire field names (confirm
// against Dũng's cw-ics08 L2 client deserializer before freezing).
func TestClientMessageData_JSONShape(t *testing.T) {
	d := &ClientMessageData{
		L1BeaconSlot:      7,
		L1AccountProof:    [][]byte{{0x1}},
		L2HeaderRLP:       []byte{0x2},
		L2IBCAccountProof: [][]byte{{0x3}},
		Optimism:          &OptimismInputs{OutputRootPreimage: []byte{0x4}},
	}
	raw, err := d.Encode()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	for _, field := range []string{"l1_beacon_slot", "l1_account_proof", "l2_header_rlp", "l2_ibc_account_proof", "optimism", "output_root_preimage"} {
		if !strings.Contains(string(raw), field) {
			t.Fatalf("client message JSON missing %q: %s", field, raw)
		}
	}
	// The unused chain-specific set must be omitted.
	if strings.Contains(string(raw), "arbitrum") {
		t.Fatalf("empty arbitrum set must be omitted: %s", raw)
	}
}
