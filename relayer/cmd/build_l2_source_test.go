package main

import (
	"encoding/json"
	"strings"
	"testing"

	"relayer/chain"
	"relayer/chain/l2rollup"
)

func validL2Config() l2ToCosmosConfig {
	return l2ToCosmosConfig{
		L1RpcUrl:         "http://127.0.0.1:8545",
		L2RpcUrl:         "http://127.0.0.1:9545",
		TmRpcUrl:         "http://127.0.0.1:26657",
		AttestorAddr:     "127.0.0.1:3001",
		AttestorSrcChain: "op-sepolia",
		L2WasmClientID:   "08-wasm-1",
		L2ICS26ClientID:  "client-0",
		HeadKind:         "safe",
		RollupProfile:    json.RawMessage(`{"common":{"l2_router":"0x1111111111111111111111111111111111111111"}}`),
		kind:             chain.OPStack,
	}
}

func TestL2Config_Validate(t *testing.T) {
	if err := validL2Config().validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}

	cases := map[string]func(*l2ToCosmosConfig){
		"missing l1_rpc_url":     func(c *l2ToCosmosConfig) { c.L1RpcUrl = "" },
		"missing attestor_addr":  func(c *l2ToCosmosConfig) { c.AttestorAddr = "" },
		"missing wasm client id": func(c *l2ToCosmosConfig) { c.L2WasmClientID = "" },
		"empty profile":          func(c *l2ToCosmosConfig) { c.RollupProfile = nil },
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
		"l2_wasm_client_id":"08-wasm-1","l2_ics26_client_id":"client-0","head_kind":"safe",
		"rollup_profile":{"common":{"l2_router":"0x1111111111111111111111111111111111111111"}}}}]}`
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
