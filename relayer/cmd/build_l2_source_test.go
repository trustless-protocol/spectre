package main

import (
	"encoding/json"
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
		RollupProfile:    json.RawMessage(`{"common":{"l2_router":"0x1111111111111111111111111111111111111111"}}`),
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
		"eth_beacon_api_url":"http://beacon",
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

func TestFindL2TimeoutReturnPath(t *testing.T) {
	src := validL2Config()
	src.L2RpcUrl = "http://l2-a"
	src.TmRpcUrl = "http://cosmos-a"
	src.RollupProfile = json.RawMessage(`{"common":{"l2_router":"0x1111111111111111111111111111111111111111"}}`)
	matchingDest := cosmosToEthConfig{
		EthRpcUrl:     "http://l2-a",
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
			name:  "matches by l2 rpc cosmos rpc and router",
			paths: []l2TimeoutReturnPathConfig{{cfg: matchingDest}},
		},
		{
			name: "no match fails",
			paths: []l2TimeoutReturnPathConfig{{cfg: cosmosToEthConfig{
				EthRpcUrl:    "http://l2-b",
				TmRpcUrl:     "http://cosmos-a",
				ICS26Address: "0x1111111111111111111111111111111111111111",
			}}},
			wantErr: "no matching cosmos_to_l2 return path",
		},
		{
			name: "duplicate match fails",
			paths: []l2TimeoutReturnPathConfig{
				{cfg: matchingDest},
				{cfg: matchingDest},
			},
			wantErr: "multiple cosmos_to_l2 return paths match",
		},
		{
			name: "bad destination router fails",
			paths: []l2TimeoutReturnPathConfig{{cfg: cosmosToEthConfig{
				EthRpcUrl:    "http://l2-a",
				TmRpcUrl:     "http://cosmos-a",
				ICS26Address: "0xnothex",
			}}},
			wantErr: "not a valid hex address",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := findL2TimeoutReturnPath(src, tc.paths)
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
