package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReplaceICS07Address(t *testing.T) {
	t.Parallel()

	const addr = "0x1234567890abcdef1234567890abcdef12345678"

	tests := []struct {
		name   string
		input  string
		assert func(t *testing.T, out []byte, cosmosToEth map[string]any)
	}{
		{
			name: "replaces existing spectre_client",
			input: `{
				"modules": [
					{"name": "cosmos_to_eth", "config": {"spectre_client": "0xold", "other": "keep"}},
					{"name": "eth_to_cosmos", "config": {"spectre_client": "0xother"}}
				]
			}`,
			assert: func(t *testing.T, out []byte, cosmosToEth map[string]any) {
				if got := cosmosToEth["other"]; got != "keep" {
					t.Fatalf("other field = %v, want keep", got)
				}
				ethToCosmos := moduleConfigByName(t, out, "eth_to_cosmos")
				if got := ethToCosmos["spectre_client"]; got != "0xother" {
					t.Fatalf("eth_to_cosmos spectre_client = %v, want 0xother", got)
				}
			},
		},
		{
			name: "inserts spectre_client when absent",
			input: `{
				"modules": [
					{"name": "cosmos_to_eth", "config": {"eth_rpc_url": "http://localhost", "nested": {"keep": true}}}
				]
			}`,
			assert: func(t *testing.T, _ []byte, cosmosToEth map[string]any) {
				if got := cosmosToEth["eth_rpc_url"]; got != "http://localhost" {
					t.Fatalf("eth_rpc_url = %v, want http://localhost", got)
				}
				nested, ok := cosmosToEth["nested"].(map[string]any)
				if !ok {
					t.Fatalf("nested = %T, want object", cosmosToEth["nested"])
				}
				if got := nested["keep"]; got != true {
					t.Fatalf("nested.keep = %v, want true", got)
				}
			},
		},
		{
			name: "preserves nested config objects",
			input: `{
				"modules": [
					{"name": "cosmos_to_eth", "config": {"nested": {"config": {"spectre_client": "nested"}, "items": [{"k": "v,}]"}]}, "spectre_client": "0xold"}}
				]
			}`,
			assert: func(t *testing.T, _ []byte, cosmosToEth map[string]any) {
				nested, ok := cosmosToEth["nested"].(map[string]any)
				if !ok {
					t.Fatalf("nested = %T, want object", cosmosToEth["nested"])
				}
				inner, ok := nested["config"].(map[string]any)
				if !ok {
					t.Fatalf("nested.config = %T, want object", nested["config"])
				}
				if got := inner["spectre_client"]; got != "nested" {
					t.Fatalf("nested.config.spectre_client = %v, want nested", got)
				}
			},
		},
		{
			name: "handles escaped strings and keys",
			input: `{
				"modules": [
					{"name": "cosmos\u005fto\u005feth", "config": {"note": "quote: \" and slash: \\", "ics07\u005fclient": "0xold"}}
				]
			}`,
			assert: func(t *testing.T, _ []byte, cosmosToEth map[string]any) {
				if got := cosmosToEth["note"]; got != `quote: " and slash: \` {
					t.Fatalf("note = %q, want escaped string preserved", got)
				}
			},
		},
		{
			name:  "handles unusual whitespace",
			input: "{\n\t\"modules\"\t:\t[\n\t { \n \"name\" : \"cosmos_to_eth\" ,\n \"config\" : { \n \"foo\" : [1, {\"bar\" : \"baz\"}] \n } \n }\n]\n}",
			assert: func(t *testing.T, _ []byte, cosmosToEth map[string]any) {
				foo, ok := cosmosToEth["foo"].([]any)
				if !ok {
					t.Fatalf("foo = %T, want array", cosmosToEth["foo"])
				}
				if len(foo) != 2 {
					t.Fatalf("foo length = %d, want 2", len(foo))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := replaceConfigMember([]byte(tc.input), "spectre_client", addr)
			if err != nil {
				t.Fatalf("replaceConfigMember() error = %v", err)
			}

			cosmosToEth := moduleConfigByName(t, out, "cosmos_to_eth")
			if got := cosmosToEth["spectre_client"]; got != addr {
				t.Fatalf("spectre_client = %v, want %s", got, addr)
			}
			tc.assert(t, out, cosmosToEth)
		})
	}
}

func TestReplaceConfigMember_WasmClientID(t *testing.T) {
	t.Parallel()

	const id = "08-wasm-445"

	t.Run("inserts cosmos_wasm_client_id when absent", func(t *testing.T) {
		t.Parallel()
		in := `{
			"modules": [
				{"name": "cosmos_to_eth", "config": {"spectre_client": "0xabc"}}
			]
		}`
		out, err := replaceConfigMember([]byte(in), "cosmos_wasm_client_id", id)
		if err != nil {
			t.Fatalf("replaceConfigMember() error = %v", err)
		}
		cfg := moduleConfigByName(t, out, "cosmos_to_eth")
		if got := cfg["cosmos_wasm_client_id"]; got != id {
			t.Fatalf("cosmos_wasm_client_id = %v, want %s", got, id)
		}
		if got := cfg["spectre_client"]; got != "0xabc" {
			t.Fatalf("spectre_client = %v, want preserved 0xabc", got)
		}
	})

	t.Run("replaces existing cosmos_wasm_client_id", func(t *testing.T) {
		t.Parallel()
		in := `{
			"modules": [
				{"name": "cosmos_to_eth", "config": {"cosmos_wasm_client_id": "08-wasm-0", "keep": "yes"}}
			]
		}`
		out, err := replaceConfigMember([]byte(in), "cosmos_wasm_client_id", id)
		if err != nil {
			t.Fatalf("replaceConfigMember() error = %v", err)
		}
		cfg := moduleConfigByName(t, out, "cosmos_to_eth")
		if got := cfg["cosmos_wasm_client_id"]; got != id {
			t.Fatalf("cosmos_wasm_client_id = %v, want %s", got, id)
		}
		if got := cfg["keep"]; got != "yes" {
			t.Fatalf("keep = %v, want yes", got)
		}
	})
}

func TestLoadConfigRejectsDeprecatedBatchConfig(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(configPath, []byte(`{"batch_config":{"batch_size":1},"modules":[]}`), configFilePerm); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := loadConfig(configPath)
	if err == nil {
		t.Fatal("expected deprecated batch_config error")
	}
	if !strings.Contains(err.Error(), "batch_config is deprecated") {
		t.Fatalf("loadConfig error = %v, want deprecated batch_config error", err)
	}
}

func TestLoadConfigExample(t *testing.T) {
	t.Parallel()

	cfg, err := loadConfig(filepath.Join("..", "config.example.json"))
	if err != nil {
		t.Fatalf("loadConfig(config.example.json) error = %v", err)
	}
	if cfg.CosmosToEthConfig.TmRpcUrl == "" {
		t.Fatal("cosmos_to_eth.tm_rpc_url not parsed")
	}
	if cfg.CosmosToEthConfig.ICS26ClientID == "" {
		t.Fatal("ics26_client_id not parsed")
	}
	if cfg.EthToCosmosConfig.BeaconUrl == "" {
		t.Fatal("eth_to_cosmos.eth_beacon_api_url not parsed")
	}
}

// TestLoadConfigIgnoresPrunedFields confirms old configs carrying the now-removed
// eth_to_cosmos fields and dst_chain still load (unknown JSON keys are ignored),
// preserving backward compatibility.
func TestLoadConfigIgnoresPrunedFields(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.json")
	legacy := `{
		"modules": [
			{"name": "cosmos_to_eth", "src_chain": "test", "dst_chain": "0x1",
			 "config": {"tm_rpc_url": "http://localhost:26657", "eth_rpc_url": "http://localhost:8545", "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3"}},
			{"name": "eth_to_cosmos", "dst_chain": "test",
			 "config": {"eth_beacon_api_url": "http://localhost:5052", "tm_rpc_url": "http://localhost:26657", "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "signer_address": "cosmos1abc"}}
		]
	}`
	if err := os.WriteFile(configPath, []byte(legacy), configFilePerm); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("loadConfig(legacy) error = %v", err)
	}
	if cfg.EthToCosmosConfig.BeaconUrl != "http://localhost:5052" {
		t.Fatalf("beacon url = %q, want parsed", cfg.EthToCosmosConfig.BeaconUrl)
	}
}

func moduleConfigByName(t *testing.T, data []byte, name string) map[string]any {
	t.Helper()

	var cfg struct {
		Modules []struct {
			Name   string         `json:"name"`
			Config map[string]any `json:"config"`
		} `json:"modules"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, string(data))
	}
	for _, m := range cfg.Modules {
		if m.Name == name {
			if m.Config == nil {
				t.Fatalf("module %q config is nil", name)
			}
			return m.Config
		}
	}
	t.Fatalf("module %q not found", name)
	return nil
}

// TestLoadConfigMultipleCosmosSources confirms a config with two cosmos_to_eth
// modules yields both in CosmosToEthConfigs (file order), with the singular
// CosmosToEthConfig aliasing the first for the single-source commands.
func TestLoadConfigMultipleCosmosSources(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.json")
	cfgJSON := `{
		"modules": [
			{"name": "cosmos_to_eth", "src_chain": "chain-a",
			 "config": {"tm_rpc_url": "http://localhost:26657", "eth_rpc_url": "http://localhost:8545", "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "chain-a"}},
			{"name": "cosmos_to_eth", "src_chain": "chain-b",
			 "config": {"tm_rpc_url": "http://localhost:36657", "eth_rpc_url": "http://localhost:8545", "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "chain-b"}},
			{"name": "eth_to_cosmos", "config": {"eth_beacon_api_url": "http://localhost:5052"}}
		]
	}`
	if err := os.WriteFile(configPath, []byte(cfgJSON), configFilePerm); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("loadConfig(multi) error = %v", err)
	}
	if got := len(cfg.CosmosToEthConfigs); got != 2 {
		t.Fatalf("CosmosToEthConfigs len = %d, want 2", got)
	}
	if cfg.CosmosToEthConfigs[0].ICS26ClientID != "chain-a" ||
		cfg.CosmosToEthConfigs[1].ICS26ClientID != "chain-b" {
		t.Fatalf("source order/ids wrong: %q, %q",
			cfg.CosmosToEthConfigs[0].ICS26ClientID, cfg.CosmosToEthConfigs[1].ICS26ClientID)
	}
	if cfg.CosmosToEthConfigs[1].TmRpcUrl != "http://localhost:36657" {
		t.Fatalf("second source tm_rpc_url = %q, want distinct", cfg.CosmosToEthConfigs[1].TmRpcUrl)
	}
	if cfg.CosmosToEthConfig.ICS26ClientID != "chain-a" {
		t.Fatalf("singular CosmosToEthConfig = %q, want alias of first source", cfg.CosmosToEthConfig.ICS26ClientID)
	}
}

// TestLoadConfigRejectsDuplicateClientID confirms two sources sharing an
// ics26_client_id are rejected — their ETH event streams (filtered by that id)
// would otherwise cross-feed.
func TestLoadConfigRejectsDuplicateClientID(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.json")
	cfgJSON := `{
		"modules": [
			{"name": "cosmos_to_eth", "config": {"tm_rpc_url": "http://localhost:26657", "eth_rpc_url": "http://localhost:8545", "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "dup"}},
			{"name": "cosmos_to_eth", "config": {"tm_rpc_url": "http://localhost:36657", "eth_rpc_url": "http://localhost:8545", "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "dup"}}
		]
	}`
	if err := os.WriteFile(configPath, []byte(cfgJSON), configFilePerm); err != nil {
		t.Fatalf("write config: %v", err)
	}
	_, err := loadConfig(configPath)
	if err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("loadConfig error = %v, want duplicate client id error", err)
	}
}

// TestLoadConfigExampleSingleSource confirms the shipped example yields exactly
// one source, so single-source behavior is unchanged.
func TestLoadConfigExampleSingleSource(t *testing.T) {
	t.Parallel()

	cfg, err := loadConfig(filepath.Join("..", "config.example.json"))
	if err != nil {
		t.Fatalf("loadConfig(config.example.json) error = %v", err)
	}
	if got := len(cfg.CosmosToEthConfigs); got != 1 {
		t.Fatalf("CosmosToEthConfigs len = %d, want 1", got)
	}
}

// TestReplaceConfigMemberForSource confirms a write-back targets only the named
// source's module and leaves the other source untouched.
func TestReplaceConfigMemberForSource(t *testing.T) {
	t.Parallel()

	in := `{"modules":[` +
		`{"name":"cosmos_to_eth","config":{"ics26_client_id":"chain-a","ics07_client":"0xAAA"}},` +
		`{"name":"cosmos_to_eth","config":{"ics26_client_id":"chain-b","ics07_client":"0xBBB"}}` +
		`]}`
	out, err := replaceConfigMemberForSource([]byte(in), "chain-b", "ics07_client", "0xNEW")
	if err != nil {
		t.Fatalf("replaceConfigMemberForSource() error = %v", err)
	}
	got := string(out)
	if !strings.Contains(got, `"chain-a","ics07_client":"0xAAA"`) {
		t.Fatalf("chain-a should be untouched, got: %s", got)
	}
	if !strings.Contains(got, `"chain-b","ics07_client":"0xNEW"`) {
		t.Fatalf("chain-b should be updated, got: %s", got)
	}
}

// TestReplaceConfigMemberForSourceUnknown confirms an unknown source errors
// rather than silently writing the wrong module.
func TestReplaceConfigMemberForSourceUnknown(t *testing.T) {
	t.Parallel()

	in := `{"modules":[{"name":"cosmos_to_eth","config":{"ics26_client_id":"chain-a"}}]}`
	_, err := replaceConfigMemberForSource([]byte(in), "chain-x", "ics07_client", "0x1")
	if err == nil || !strings.Contains(err.Error(), "chain-x") {
		t.Fatalf("error = %v, want unknown-source error", err)
	}
}

// TestReplaceConfigMemberForSourceViaSrcChain confirms the module-level src_chain
// is used as the source id when config.ics26_client_id is absent (matching how
// loadConfig defaults the id).
func TestReplaceConfigMemberForSourceViaSrcChain(t *testing.T) {
	t.Parallel()

	in := `{"modules":[{"name":"cosmos_to_eth","src_chain":"chain-a","config":{"ics07_client":"0xAAA"}}]}`
	out, err := replaceConfigMemberForSource([]byte(in), "chain-a", "ics07_client", "0xNEW")
	if err != nil {
		t.Fatalf("replaceConfigMemberForSource() error = %v", err)
	}
	if !strings.Contains(string(out), `"ics07_client":"0xNEW"`) {
		t.Fatalf("src_chain match failed, got: %s", string(out))
	}
}

func TestSelectSource(t *testing.T) {
	t.Parallel()

	cfg := &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{
			{ICS26ClientID: "chain-a", TmRpcUrl: "a"},
			{ICS26ClientID: "chain-b", TmRpcUrl: "b"},
		},
	}
	single := &appConfig{
		CosmosToEthConfig:  cosmosToEthConfig{ICS26ClientID: "solo", TmRpcUrl: "s"},
		CosmosToEthConfigs: []cosmosToEthConfig{{ICS26ClientID: "solo", TmRpcUrl: "s"}},
	}

	t.Run("by id", func(t *testing.T) {
		got, err := selectSource(cfg, "chain-b")
		if err != nil {
			t.Fatalf("selectSource error = %v", err)
		}
		if got.CosmosToEthConfig.TmRpcUrl != "b" {
			t.Fatalf("selected wrong source: %+v", got.CosmosToEthConfig)
		}
	})
	t.Run("empty with multiple errors", func(t *testing.T) {
		_, err := selectSource(cfg, "")
		if err == nil || !strings.Contains(err.Error(), "--source") {
			t.Fatalf("error = %v, want ambiguous-source error", err)
		}
	})
	t.Run("unknown id errors", func(t *testing.T) {
		_, err := selectSource(cfg, "chain-x")
		if err == nil || !strings.Contains(err.Error(), "chain-x") {
			t.Fatalf("error = %v, want unknown-source error", err)
		}
	})
	t.Run("empty with single selects sole", func(t *testing.T) {
		got, err := selectSource(single, "")
		if err != nil {
			t.Fatalf("selectSource error = %v", err)
		}
		if got.CosmosToEthConfig.TmRpcUrl != "s" {
			t.Fatalf("selected wrong source: %+v", got.CosmosToEthConfig)
		}
	})
}
