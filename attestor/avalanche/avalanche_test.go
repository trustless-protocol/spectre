package avalanche

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"attestor/core"
	"attestor/types/attestation"

	ethcommon "github.com/ethereum/go-ethereum/common"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

const testSeedHex = "fac8c6533c38c5670c6eadac6cd5d5722a24d8e9eb5705111950f1e849fe15fb"

// fakeCChain serves eth_getBlockByNumber for a fixed canonical chain.
func fakeCChain(t *testing.T, blocks map[string]map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
			Params []any           `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode rpc request: %v", err)
			return
		}
		if req.Method != "eth_getBlockByNumber" {
			t.Errorf("unexpected method %q", req.Method)
			return
		}
		tag, _ := req.Params[0].(string)
		block, ok := blocks[tag]
		var result any
		if ok {
			result = block
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": result})
	}))
}

func testAttestor(t *testing.T, server *httptest.Server) (*CChainAttestor, attestation.Signer) {
	t.Helper()
	client, err := gethrpc.Dial(server.URL)
	if err != nil {
		t.Fatalf("dial fake rpc: %v", err)
	}
	t.Cleanup(client.Close)
	signer, err := attestation.NewSigner(43112, testSeedHex)
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	a, err := New(Config{SrcChain: "avaxdev", RpcURL: server.URL, ChainID: 43112, PollInterval: time.Hour}, client, signer)
	if err != nil {
		t.Fatalf("new attestor: %v", err)
	}
	return a, signer
}

func block(number uint64, hashByte, rootByte byte) map[string]string {
	return map[string]string{
		"number":    fmt.Sprintf("0x%x", number),
		"hash":      ethcommon.Hash{hashByte}.Hex(),
		"stateRoot": ethcommon.Hash{rootByte}.Hex(),
	}
}

func TestVerifyStateRootSignsOnlyTheCanonicalIdentity(t *testing.T) {
	server := fakeCChain(t, map[string]map[string]string{
		"finalized": block(7, 0xaa, 0xbb),
		"0x5":       block(5, 0xcc, 0xdd),
	})
	defer server.Close()
	a, signer := testAttestor(t, server)

	// Head not polled yet: verification is unavailable, not a mismatch.
	if _, err := a.VerifyStateRoot(context.Background(), core.BlockIdentityRequest{
		BlockNumber: 5, RunMode: core.RunModeFinalized,
	}); err == nil {
		t.Fatal("expected unavailable before the first poll")
	}

	head, err := a.fetchBlock(context.Background(), "finalized")
	if err != nil {
		t.Fatalf("fetch head: %v", err)
	}
	a.mu.Lock()
	a.head, a.seen = head, true
	a.mu.Unlock()

	router := [20]byte{1}
	setHash := [32]byte{2}
	request := core.BlockIdentityRequest{
		BlockNumber:       5,
		ExpectedStateRoot: [32]byte(ethcommon.Hash{0xdd}),
		RunMode:           core.RunModeFinalized,
		L2Router:          router,
		AttestorSetHash:   setHash,
	}
	verdict, err := a.VerifyStateRoot(context.Background(), request)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !verdict.Valid || verdict.BlockNumber != 5 {
		t.Fatalf("verdict = %+v, want valid at height 5", verdict)
	}
	statement, err := attestation.SigningBytes(43112, router[:], setHash[:], 5, verdict.BlockHash[:], verdict.StateRoot[:])
	if err != nil {
		t.Fatalf("statement: %v", err)
	}
	if !ed25519.Verify(ed25519.PublicKey(signer.PublicKey()), statement, verdict.Signature) {
		t.Fatal("signature does not verify over the canonical statement")
	}

	// Wrong expected root: valid=false, no signature, no error.
	request.ExpectedStateRoot = [32]byte(ethcommon.Hash{0xee})
	verdict, err = a.VerifyStateRoot(context.Background(), request)
	if err != nil {
		t.Fatalf("mismatch verify: %v", err)
	}
	if verdict.Valid || len(verdict.Signature) != 0 {
		t.Fatalf("mismatch must not be signed: %+v", verdict)
	}

	// Above the accepted head: precondition failure, retryable later.
	request.BlockNumber = 8
	if _, err := a.VerifyStateRoot(context.Background(), request); err == nil {
		t.Fatal("expected precondition failure above the head")
	}

	// Wrong run mode: this plugin only serves finalized.
	request.BlockNumber = 5
	request.RunMode = core.RunModeSafe
	if _, err := a.VerifyStateRoot(context.Background(), request); err == nil {
		t.Fatal("expected run-mode mismatch to fail")
	}
}

func TestFeedReadsFollowTheAcceptedHead(t *testing.T) {
	server := fakeCChain(t, map[string]map[string]string{
		"finalized": block(7, 0xaa, 0xbb),
		"0x3":       block(3, 0x11, 0x22),
	})
	defer server.Close()
	a, _ := testAttestor(t, server)

	if _, found, err := a.AttestedUpTo(context.Background(), core.AttestationPolicy{}); err != nil || found {
		t.Fatalf("before first poll: found=%v err=%v, want not found", found, err)
	}

	head, err := a.fetchBlock(context.Background(), "finalized")
	if err != nil {
		t.Fatalf("fetch head: %v", err)
	}
	a.mu.Lock()
	a.head, a.seen = head, true
	a.mu.Unlock()

	up, found, err := a.AttestedUpTo(context.Background(), core.AttestationPolicy{})
	if err != nil || !found || up.L2BlockNumber != 7 || up.Provisional {
		t.Fatalf("AttestedUpTo = %+v found=%v err=%v", up, found, err)
	}
	below, found, err := a.AttestedRootAtOrBelow(context.Background(), 3, core.AttestationPolicy{})
	if err != nil || !found || below.L2BlockNumber != 3 || string(below.Root) != string(ethcommon.Hash{0x22}.Bytes()) {
		t.Fatalf("AttestedRootAtOrBelow = %+v found=%v err=%v", below, found, err)
	}
	capped, found, err := a.AttestedRootAtOrBelow(context.Background(), 99, core.AttestationPolicy{})
	if err != nil || !found || capped.L2BlockNumber != 7 {
		t.Fatalf("AttestedRootAtOrBelow above head = %+v found=%v err=%v", capped, found, err)
	}

	status, err := a.Status(context.Background())
	if err != nil || !status.Ready || status.ReplicaFinalized != 7 || status.AttestationHead != core.RunModeFinalized {
		t.Fatalf("Status = %+v err=%v", status, err)
	}
}
