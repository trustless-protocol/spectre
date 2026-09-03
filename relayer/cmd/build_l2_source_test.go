package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"

	"relayer/chain"
	"relayer/chain/l2rollup"
	"relayer/services"
)

func validL2Config() l2ToCosmosConfig {
	return l2ToCosmosConfig{
		L2RpcUrl:         "http://127.0.0.1:9545",
		TmRpcUrl:         "http://127.0.0.1:26657",
		AttestorAddr:     "127.0.0.1:3001",
		AttestorSrcChain: "op-sepolia",
		L2WasmClientID:   "08-wasm-1",
		L2ICS26ClientID:  "client-0",
		HeadKind:         "safe",
		RollupProfile:    json.RawMessage(`{"common":{"l2_chain_id":10,"l2_router":"0x1111111111111111111111111111111111111111","attestor_public_key":"0x1111111111111111111111111111111111111111111111111111111111111111","attestation_head":"safe"}}`),
		kind:             chain.OPStack,
	}
}

func TestL2Config_Validate(t *testing.T) {
	if err := validL2Config().validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}

	cases := map[string]func(*l2ToCosmosConfig){
		"missing l2_rpc_url":     func(c *l2ToCosmosConfig) { c.L2RpcUrl = "" },
		"missing attestor_addr":  func(c *l2ToCosmosConfig) { c.AttestorAddr = "" },
		"missing wasm client id": func(c *l2ToCosmosConfig) { c.L2WasmClientID = "" },
		"empty profile":          func(c *l2ToCosmosConfig) { c.RollupProfile = nil },
		"profile head mismatch": func(c *l2ToCosmosConfig) {
			c.RollupProfile = json.RawMessage(`{"common":{"l2_chain_id":10,"l2_router":"0x1111111111111111111111111111111111111111","attestor_public_key":"0x1111111111111111111111111111111111111111111111111111111111111111","attestation_head":"finalized"}}`)
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			c := validL2Config()
			mutate(&c)
			if err := c.validate(); err == nil {
				t.Fatalf("expected validation error for %s", name)
			}
		})
	}
}

// TestL2Config_ValidateForClientCreation covers #309: create-clients-cosmos is the
// command that CREATES the L2 wasm client and writes its id back into the config,
// so requiring that id to already be present made a fresh config unusable — the
// operator had to invent a plausible one, which the command then overwrote. Only
// that id is exempt; everything else is still needed to reach the chain at all.
func TestL2Config_ValidateForClientCreation(t *testing.T) {
	c := validL2Config()
	c.L2WasmClientID = ""
	if err := c.validateForClientCreation(); err != nil {
		t.Fatalf("empty l2_wasm_client_id rejected during client creation: %v", err)
	}

	// The exemption is narrow: it must not spill onto the other required fields.
	for name, mutate := range map[string]func(*l2ToCosmosConfig){
		"missing l2_rpc_url":      func(c *l2ToCosmosConfig) { c.L2RpcUrl = "" },
		"missing attestor_addr":   func(c *l2ToCosmosConfig) { c.AttestorAddr = "" },
		"missing l2_ics26_client": func(c *l2ToCosmosConfig) { c.L2ICS26ClientID = "" },
		"missing tm_rpc_url":      func(c *l2ToCosmosConfig) { c.TmRpcUrl = "" },
		"missing attestor_src":    func(c *l2ToCosmosConfig) { c.AttestorSrcChain = "" },
		"empty profile":           func(c *l2ToCosmosConfig) { c.RollupProfile = nil },
		"bad head_kind":           func(c *l2ToCosmosConfig) { c.HeadKind = "nonsense" },
	} {
		t.Run(name, func(t *testing.T) {
			c := validL2Config()
			c.L2WasmClientID = ""
			mutate(&c)
			if err := c.validateForClientCreation(); err == nil {
				t.Fatalf("expected validation error for %s", name)
			}
		})
	}

	// A config that is already complete stays valid on this path too.
	if err := validL2Config().validateForClientCreation(); err != nil {
		t.Fatalf("complete config rejected: %v", err)
	}
}

// TestLoadConfigForClientCreation_AcceptsUncreatedWasmClient is the end-to-end of
// the above: the whole config file loads, rather than failing at the first module.
func TestLoadConfigForClientCreation_AcceptsUncreatedWasmClient(t *testing.T) {
	raw := `{"modules":[
		{"name":"cosmos-to-arb","src_chain":"cosmos","dst_chain":"arbitrum","config":{
			"cosmos_wasm_client_id":"","eth_rpc_url":"http://l2","tm_rpc_url":"http://tm",
			"ics26_address":"0x1111111111111111111111111111111111111111",
			"ics26_client_id":"arb-client-0"}},
		{"name":"arb-to-cosmos","src_chain":"arbitrum","dst_chain":"cosmos","config":{
			"l2_rpc_url":"http://l2","tm_rpc_url":"http://tm",
			"attestor_addr":"127.0.0.1:3002","attestor_src_chain":"arbitrum-sepolia",
			"l2_wasm_client_id":"","l2_ics26_client_id":"arb-client-0","head_kind":"unsafe",
			"rollup_profile":{"common":{"l2_chain_id":421614,"l2_router":"0x1111111111111111111111111111111111111111","attestor_public_key":"0x1111111111111111111111111111111111111111111111111111111111111111","attestation_head":"unsafe"}}}}]}`

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(raw), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := loadConfig(path); err == nil {
		t.Fatal("loadConfig accepted an empty l2_wasm_client_id; the relay path needs it")
	}

	cfg, err := loadConfigForClientCreation(path)
	if err != nil {
		t.Fatalf("loadConfigForClientCreation rejected a pre-creation config: %v", err)
	}
	if len(cfg.L2ToCosmosConfigs) != 1 {
		t.Fatalf("L2ToCosmosConfigs = %d, want 1", len(cfg.L2ToCosmosConfigs))
	}
	if got := cfg.L2ToCosmosConfigs[0].L2ICS26ClientID; got != "arb-client-0" {
		t.Errorf("l2_ics26_client_id = %q, want arb-client-0", got)
	}

	// Closing the loop: the id the command creates lands in the field that was
	// empty, and the config then loads on the strict path. Without this the
	// relaxation above would just move the failure one step later.
	if err := writeL2SourceClientIDs(path, "arb-client-0", "08-wasm-0"); err != nil {
		t.Fatalf("write back the created client id: %v", err)
	}
	reloaded, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig after write-back: %v", err)
	}
	if got := reloaded.L2ToCosmosConfigs[0].L2WasmClientID; got != "08-wasm-0" {
		t.Errorf("l2_wasm_client_id = %q, want 08-wasm-0", got)
	}
}

func TestParseHeadKind(t *testing.T) {
	cases := map[string]l2rollup.HeadKind{
		"":          l2rollup.Safe,
		"safe":      l2rollup.Safe,
		"unsafe":    l2rollup.Unsafe,
		"finalized": l2rollup.Finalized,
	}
	for in, want := range cases {
		got, err := parseHeadKind(in)
		if err != nil {
			t.Fatalf("parseHeadKind(%q): %v", in, err)
		}
		if got != want {
			t.Fatalf("parseHeadKind(%q) = %v, want %v", in, got, want)
		}
	}
	if _, err := parseHeadKind("nonsense"); err == nil {
		t.Fatal("expected error for invalid head_kind")
	}
}

func TestL2RouterFromProfile(t *testing.T) {
	addr, err := l2RouterFromProfile(json.RawMessage(`{"common":{"l2_router":"0xAbC1111111111111111111111111111111111111"}}`))
	if err != nil {
		t.Fatalf("valid router rejected: %v", err)
	}
	// Address parsing is case-insensitive; the canonical hex is lowercased internally.
	if !strings.EqualFold(addr.Hex(), "0xAbC1111111111111111111111111111111111111") {
		t.Fatalf("router = %s, want 0xAbC1111111111111111111111111111111111111", addr.Hex())
	}

	for _, bad := range []string{`{"common":{"l2_router":"0xnothex"}}`, `{"common":{}}`, `{`} {
		if _, err := l2RouterFromProfile(json.RawMessage(bad)); err == nil {
			t.Fatalf("expected error for %q", bad)
		}
	}
}

// TestLoadConfig_L2Source confirms an l2_to_cosmos module parses into
// appConfig.L2ToCosmosConfigs (no longer the fail-loud "not wired" error).
func TestLoadConfig_L2Source(t *testing.T) {
	raw := `{"modules":[{"name":"op","src_chain":"opstack","dst_chain":"cosmos","config":{
		"l1_rpc_url":"http://l1","l2_rpc_url":"http://l2","tm_rpc_url":"http://tm",
		"attestor_addr":"127.0.0.1:3001","attestor_src_chain":"op-sepolia",
		"eth_beacon_api_url":"http://beacon",
		"l2_wasm_client_id":"08-wasm-1","l2_ics26_client_id":"client-0","head_kind":"safe",
		"rollup_profile":{"common":{"l2_chain_id":10,"l2_router":"0x1111111111111111111111111111111111111111","attestor_public_key":"0x1111111111111111111111111111111111111111111111111111111111111111","attestation_head":"safe"}}}}]}`
	var jc jsonConfig
	if err := json.Unmarshal([]byte(raw), &jc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	m := jc.Modules[0]
	dir, _, err := classifyModule(m)
	if err != nil || dir != dirL2ToCosmos {
		t.Fatalf("classify = %v, %v; want dirL2ToCosmos", dir, err)
	}
	var one l2ToCosmosConfig
	if err := json.Unmarshal(m.Config, &one); err != nil {
		t.Fatalf("parse l2 config: %v", err)
	}
	one.kind = chain.ChainType(m.SrcChain)
	if err := one.validate(); err != nil {
		t.Fatalf("validate parsed l2 config: %v", err)
	}
	if one.kind != chain.OPStack || one.AttestorSrcChain != "op-sepolia" {
		t.Fatalf("parsed config mismatch: %+v", one)
	}
}

func newChainIDRPCServer(t *testing.T, chainIDHex string) *httptest.Server {
	t.Helper()
	type rpcReq struct {
		JSONRPC string `json:"jsonrpc"`
		ID      any    `json:"id"`
		Method  string `json:"method"`
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var req rpcReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode rpc request: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Method != "eth_chainId" {
			t.Errorf("rpc method = %q, want eth_chainId", req.Method)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"error":   map[string]any{"code": -32601, "message": "method not found"},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result":  chainIDHex,
		})
	}))
}

func TestFindL2TimeoutReturnPath(t *testing.T) {
	sourceRPC := newChainIDRPCServer(t, "0xa")
	t.Cleanup(sourceRPC.Close)

	src := validL2Config()
	src.L2RpcUrl = sourceRPC.URL
	src.TmRpcUrl = "http://cosmos-a"
	src.RollupProfile = json.RawMessage(`{"common":{"l2_chain_id":10,"l2_router":"0x1111111111111111111111111111111111111111","attestor_public_key":"0x1111111111111111111111111111111111111111111111111111111111111111","attestation_head":"safe"}}`)
	matchingDest := cosmosToEthConfig{
		EthRpcUrl:     "http://write-l2-a",
		TmRpcUrl:      "http://cosmos-a",
		ICS26Address:  "0x1111111111111111111111111111111111111111",
		ICS26ClientID: "cosmos-on-l2",
	}

	tests := []struct {
		name    string
		paths   []l2TimeoutReturnPathConfig
		wantErr string
	}{
		{
			name:  "matches by l2 chain id cosmos rpc and router",
			paths: []l2TimeoutReturnPathConfig{{name: "cosmos-to-op", cfg: matchingDest, l2ChainID: "10"}},
		},
		{
			name: "no match fails",
			paths: []l2TimeoutReturnPathConfig{{cfg: cosmosToEthConfig{
				EthRpcUrl:    "http://write-l2-b",
				TmRpcUrl:     "http://cosmos-a",
				ICS26Address: "0x1111111111111111111111111111111111111111",
			}, name: "cosmos-to-base", l2ChainID: "11"}},
			wantErr: "candidates: cosmos-to-base: l2_chain_id=11",
		},
		{
			name: "duplicate match fails",
			paths: []l2TimeoutReturnPathConfig{
				{name: "cosmos-to-op-a", cfg: matchingDest, l2ChainID: "10"},
				{name: "cosmos-to-op-b", cfg: matchingDest, l2ChainID: "10"},
			},
			wantErr: "multiple cosmos_to_l2 return paths match",
		},
		{
			name: "bad destination router fails",
			paths: []l2TimeoutReturnPathConfig{{cfg: cosmosToEthConfig{
				EthRpcUrl:    "http://write-l2-a",
				TmRpcUrl:     "http://cosmos-a",
				ICS26Address: "0xnothex",
			}, name: "cosmos-to-bad", l2ChainID: "10"}},
			wantErr: "cosmos-to-bad.ics26_address",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := findL2TimeoutReturnPath(context.Background(), src, tc.paths)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("findL2TimeoutReturnPath() error = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("findL2TimeoutReturnPath() error = %v, want substring %q", err, tc.wantErr)
			}
		})
	}
}

func TestL2PendingTrackerHooks(t *testing.T) {
	svc := services.New(nil, nil, services.DefaultConfig())
	track, untrack := l2PendingTrackerHooks(svc)

	packet := channeltypesv2.Packet{SourceClient: "client-0", Sequence: 11}
	raw, err := packet.Marshal()
	if err != nil {
		t.Fatalf("marshal packet: %v", err)
	}

	track(raw, 789)
	pending := svc.BatchBuilder.L2PendingTracker.GetAll()
	if len(pending) != 1 {
		t.Fatalf("pending len = %d, want 1", len(pending))
	}
	if pending[0].BlockNumber != 789 {
		t.Fatalf("pending block = %d, want 789", pending[0].BlockNumber)
	}

	untrack(raw)
	if got := svc.BatchBuilder.L2PendingTracker.Len(); got != 0 {
		t.Fatalf("pending len after untrack = %d, want 0", got)
	}
}
