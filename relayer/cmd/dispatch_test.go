package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClassifyModule(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		module     configModule
		wantDir    moduleDirection
		wantLegacy bool
		wantErr    bool
	}{
		{
			name:       "legacy cosmos_to_eth by name",
			module:     configModule{Name: "cosmos_to_eth", SrcChain: "cosmoshub-1"},
			wantDir:    dirCosmosToEth,
			wantLegacy: true,
		},
		{
			name:       "legacy eth_to_cosmos by name",
			module:     configModule{Name: "eth_to_cosmos"},
			wantDir:    dirEthToCosmos,
			wantLegacy: true,
		},
		{
			name:    "canonical cosmos->ethereum",
			module:  configModule{Name: "my-label", SrcChain: "cosmos", DstChain: "ethereum"},
			wantDir: dirCosmosToEth,
		},
		{
			name:    "canonical ethereum->cosmos",
			module:  configModule{Name: "eth-side", SrcChain: "ethereum", DstChain: "cosmos"},
			wantDir: dirEthToCosmos,
		},
		{
			name:    "canonical opstack->cosmos",
			module:  configModule{Name: "op", SrcChain: "opstack", DstChain: "cosmos"},
			wantDir: dirL2ToCosmos,
		},
		{
			name:    "canonical arbitrum->cosmos",
			module:  configModule{Name: "arb", SrcChain: "arbitrum", DstChain: "cosmos"},
			wantDir: dirL2ToCosmos,
		},
		{
			name:    "canonical cosmos->opstack",
			module:  configModule{Name: "op-dst", SrcChain: "cosmos", DstChain: "opstack"},
			wantDir: dirCosmosToL2,
		},
		{
			name:    "unrecognized legacy name fails loud",
			module:  configModule{Name: "typo_to_eth"},
			wantErr: true,
		},
		{
			// solana is not a known chain kind, so this is not canonical and the
			// unrecognized name fails loud.
			name:    "unknown chain kind fails loud",
			module:  configModule{Name: "x", SrcChain: "solana", DstChain: "cosmos"},
			wantErr: true,
		},
		{
			// both valid kinds but the pair is unsupported.
			name:    "unsupported canonical pair fails loud",
			module:  configModule{Name: "x", SrcChain: "arbitrum", DstChain: "ethereum"},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir, legacy, err := classifyModule(tc.module)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got dir=%q", dir)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dir != tc.wantDir {
				t.Fatalf("dir = %q, want %q", dir, tc.wantDir)
			}
			if legacy != tc.wantLegacy {
				t.Fatalf("isLegacy = %v, want %v", legacy, tc.wantLegacy)
			}
		})
	}
}

func writeTempConfig(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

// TestLoadConfig_LegacySrcChainFallback: a legacy cosmos_to_eth module with no
// ics26_client_id falls back to src_chain (back-compat).
func TestLoadConfig_LegacySrcChainFallback(t *testing.T) {
	t.Parallel()
	path := writeTempConfig(t, `{"modules":[
		{"name":"cosmos_to_eth","src_chain":"my-client-0","config":{"tm_rpc_url":"http://localhost:26657","eth_rpc_url":"http://localhost:8545","ics26_address":"0x80741a37e3644612f0465145c9709a90b6d77ee3"}}
	]}`)
	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if got := cfg.CosmosToEthConfig.ICS26ClientID; got != "my-client-0" {
		t.Fatalf("ics26_client_id = %q, want the src_chain fallback my-client-0", got)
	}
}

// TestLoadConfig_CanonicalNoFallback: a canonical module does NOT overload
// src_chain (a chain kind) as the ics26 fallback.
func TestLoadConfig_CanonicalNoFallback(t *testing.T) {
	t.Parallel()
	path := writeTempConfig(t, `{"modules":[
		{"name":"c2e","src_chain":"cosmos","dst_chain":"ethereum","config":{"ics26_client_id":"real-id","tm_rpc_url":"http://localhost:26657","eth_rpc_url":"http://localhost:8545","ics26_address":"0x80741a37e3644612f0465145c9709a90b6d77ee3"}}
	]}`)
	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if got := cfg.CosmosToEthConfig.ICS26ClientID; got != "real-id" {
		t.Fatalf("ics26_client_id = %q, want real-id (no src_chain overload)", got)
	}
}

// TestLoadConfig_UnknownModuleFailsLoud: a misspelled module is a hard error, not
// a silent drop.
func TestLoadConfig_UnknownModuleFailsLoud(t *testing.T) {
	t.Parallel()
	path := writeTempConfig(t, `{"modules":[{"name":"cosmos_to_ethh","config":{}}]}`)
	if _, err := loadConfig(path); err == nil {
		t.Fatal("expected loadConfig to fail on an unrecognized module, got nil")
	}
}

// TestLoadConfig_L2ToCosmosNotWired: the optimistic L2→Cosmos direction is not on
// this branch, so it still fails loud.
func TestLoadConfig_L2ToCosmosNotWired(t *testing.T) {
	t.Parallel()
	path := writeTempConfig(t, `{"modules":[{"name":"op","src_chain":"opstack","dst_chain":"cosmos","config":{}}]}`)
	if _, err := loadConfig(path); err == nil {
		t.Fatal("expected loadConfig to fail on an unwired l2_to_cosmos direction, got nil")
	}
}

// TestLoadConfig_CosmosToL2Parses: a cosmos→L2 module parses into CosmosToL2Configs
// (same schema as cosmos_to_eth), pointed at the L2 endpoint.
func TestLoadConfig_CosmosToL2Parses(t *testing.T) {
	t.Parallel()
	path := writeTempConfig(t, `{"modules":[
		{"name":"arb-dst","src_chain":"cosmos","dst_chain":"arbitrum","config":{"ics26_client_id":"arb-client-0","tm_rpc_url":"http://localhost:26657","eth_rpc_url":"http://localhost:8547","ics26_address":"0x80741a37e3644612f0465145c9709a90b6d77ee3","spectre_client":"0x80741a37e3644612f0465145c9709a90b6d77ee3"}}
	]}`)
	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if len(cfg.CosmosToL2Configs) != 1 {
		t.Fatalf("CosmosToL2Configs len = %d, want 1", len(cfg.CosmosToL2Configs))
	}
	if got := cfg.CosmosToL2Configs[0].ICS26ClientID; got != "arb-client-0" {
		t.Fatalf("ics26_client_id = %q, want arb-client-0", got)
	}
	if len(cfg.CosmosToEthConfigs) != 0 {
		t.Fatalf("CosmosToEthConfigs should be empty for a cosmos_to_l2-only config, got %d", len(cfg.CosmosToEthConfigs))
	}
}

func TestLoadConfig_CosmosToL2ValidationUsesModuleName(t *testing.T) {
	t.Parallel()
	path := writeTempConfig(t, `{"modules":[
		{"name":"cosmos-to-op","src_chain":"cosmos","dst_chain":"opstack","config":{"ics26_client_id":"op-client-0","tm_rpc_url":"http://localhost:26657","eth_rpc_url":"not-a-url","ics26_address":"0x80741a37e3644612f0465145c9709a90b6d77ee3"}}
	]}`)
	_, err := loadConfig(path)
	if err == nil {
		t.Fatal("expected invalid cosmos_to_l2 URL to fail")
	}
	if !strings.Contains(err.Error(), "cosmos-to-op.eth_rpc_url") {
		t.Fatalf("loadConfig error = %v, want cosmos-to-op.eth_rpc_url", err)
	}
	if strings.Contains(err.Error(), "cosmos_to_eth.eth_rpc_url") {
		t.Fatalf("loadConfig error = %v, should not label cosmos_to_l2 as cosmos_to_eth", err)
	}
}

// TestLoadConfig_CosmosToL2DuplicateClientID: two cosmos→L2 destinations must not
// share an ics26_client_id.
func TestLoadConfig_CosmosToL2DuplicateClientID(t *testing.T) {
	t.Parallel()
	path := writeTempConfig(t, `{"modules":[
		{"name":"arb-dst","src_chain":"cosmos","dst_chain":"arbitrum","config":{"ics26_client_id":"dup","tm_rpc_url":"http://localhost:26657","eth_rpc_url":"http://localhost:8547","ics26_address":"0x80741a37e3644612f0465145c9709a90b6d77ee3","spectre_client":"0x80741a37e3644612f0465145c9709a90b6d77ee3"}},
		{"name":"op-dst","src_chain":"cosmos","dst_chain":"opstack","config":{"ics26_client_id":"dup","tm_rpc_url":"http://localhost:26657","eth_rpc_url":"http://localhost:8548","ics26_address":"0x80741a37e3644612f0465145c9709a90b6d77ee3","spectre_client":"0x80741a37e3644612f0465145c9709a90b6d77ee3"}}
	]}`)
	if _, err := loadConfig(path); err == nil {
		t.Fatal("expected loadConfig to fail on duplicate cosmos_to_l2 ics26_client_id, got nil")
	}
}
