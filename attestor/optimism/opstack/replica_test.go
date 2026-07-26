package opstack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newFakeOpNode serves a minimal JSON-RPC endpoint speaking the two op-node
// methods the attestor uses.
func newFakeOpNode(t *testing.T, safe, unsafe uint64, outputRoot string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     json.RawMessage   `json:"id"`
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("bad rpc request: %v", err)
			return
		}
		var result any
		switch req.Method {
		case "optimism_syncStatus":
			result = map[string]any{
				"safe_l2":   map[string]any{"number": safe},
				"unsafe_l2": map[string]any{"number": unsafe},
			}
		case "optimism_outputAtBlock":
			result = map[string]any{"outputRoot": outputRoot}
		default:
			t.Errorf("unexpected rpc method %q", req.Method)
			return
		}
		resp := map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": result}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode rpc response: %v", err)
		}
	}))
}

func TestReplicaClientSyncStatusAndOutput(t *testing.T) {
	wantRoot := "0x" + "ab" + "00000000000000000000000000000000000000000000000000000000000000"[0:62]
	srv := newFakeOpNode(t, 123, 456, wantRoot)
	defer srv.Close()

	ctx := context.Background()
	rc, err := DialReplica(ctx, srv.URL)
	if err != nil {
		t.Fatalf("DialReplica: %v", err)
	}
	defer rc.Close()

	status, err := rc.SyncStatus(ctx)
	if err != nil {
		t.Fatalf("SyncStatus: %v", err)
	}
	if status.SafeL2 != 123 || status.UnsafeL2 != 456 {
		t.Fatalf("SyncStatus = %+v, want safe 123 unsafe 456", status)
	}

	got, err := rc.OutputAtBlock(ctx, 100)
	if err != nil {
		t.Fatalf("OutputAtBlock: %v", err)
	}
	if got[0] != 0xab {
		t.Fatalf("OutputAtBlock root[0] = %#x, want 0xab", got[0])
	}
}

func TestReplicaClientRejectsShortOutputRoot(t *testing.T) {
	srv := newFakeOpNode(t, 1, 2, "0xabcd")
	defer srv.Close()

	ctx := context.Background()
	rc, err := DialReplica(ctx, srv.URL)
	if err != nil {
		t.Fatalf("DialReplica: %v", err)
	}
	defer rc.Close()

	if _, err := rc.OutputAtBlock(ctx, 100); err == nil {
		t.Fatal("a non-32-byte output root must be rejected")
	}
}
