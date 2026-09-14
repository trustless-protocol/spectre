package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateStartupKeysRejectsInvalidCosmosAddressPrefix(t *testing.T) {
	const privateKey = "a81f9eb900c02f35d28ec80d1da79dde52cc4cc37d6fbaef5768ab59f32cbfdd"
	t.Setenv("ETH_PRIVATE_KEY", privateKey)
	t.Setenv("COSMOS_PRIVATE_KEY", privateKey)
	t.Setenv("COSMOS_ADDRESS_PREFIX", "invalid prefix")

	if err := validateStartupKeys(); err == nil {
		t.Fatal("startup accepted a Cosmos signer address with an invalid bech32 prefix")
	}
}

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

// TestLoadConfig_IgnoresLegacyServerBlock pins the compatibility half of dropping
// the dead `server` block: operators' existing config.json files still carry it,
// and they must keep loading. encoding/json ignores unknown members and nothing
// here sets DisallowUnknownFields, so this holds — the test exists to keep it that
// way, because turning on strict decoding later would break every deployed config.
func TestLoadConfig_IgnoresLegacyServerBlock(t *testing.T) {
	t.Parallel()

	raw := `{
	  "server": {"address": "127.0.0.1", "log_level": "info", "port": 3000},
	  "batch": {"batch_period_seconds": 3, "batch_size": 5},
	  "modules": [{"name":"c2e","src_chain":"cosmos","dst_chain":"ethereum","config":{
	    "tm_rpc_url":"http://tm","eth_rpc_url":"http://eth",
	    "ics26_address":"0x1111111111111111111111111111111111111111",
	    "ics26_client_id":"cosmoshub-1"}}]}`

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(raw), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig rejected a config carrying the legacy server block: %v", err)
	}
	if cfg.CosmosToEthConfig.ICS26ClientID != "cosmoshub-1" {
		t.Errorf("ics26_client_id = %q, want cosmoshub-1", cfg.CosmosToEthConfig.ICS26ClientID)
	}
	if cfg.BatchConfig.BatchSize != 5 {
		t.Errorf("batch_size = %d, want 5 — the rest of the config must still parse",
			cfg.BatchConfig.BatchSize)
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
			 "config": {"tm_rpc_url": "http://localhost:26657", "eth_rpc_url": "http://localhost:8545", "eth_ws_url": "ws://localhost:8546", "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3"}},
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
			 "config": {"tm_rpc_url": "http://localhost:26657", "eth_rpc_url": "http://localhost:8545", "eth_ws_url": "ws://localhost:8546", "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "chain-a"}},
			{"name": "cosmos_to_eth", "src_chain": "chain-b",
			 "config": {"tm_rpc_url": "http://localhost:36657", "eth_rpc_url": "http://localhost:8545", "eth_ws_url": "ws://localhost:8546", "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "chain-b"}},
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
			{"name": "cosmos_to_eth", "config": {"tm_rpc_url": "http://localhost:26657", "eth_rpc_url": "http://localhost:8545", "eth_ws_url": "ws://localhost:8546", "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "dup"}},
			{"name": "cosmos_to_eth", "config": {"tm_rpc_url": "http://localhost:36657", "eth_rpc_url": "http://localhost:8545", "eth_ws_url": "ws://localhost:8546", "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "dup"}}
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

// TestReplaceConfigMemberForSourceCanonical is the regression for the canonical
// schema: the module carries a free-form name label ("cosmos-to-eth") and is
// classified by src_chain/dst_chain, exactly as loadConfig/classifyModule does.
// The old write-back matched name == "cosmos_to_eth" and so failed to find a
// canonical module after create-clients mutated chain state. It must also skip the
// eth->cosmos module even when it shares the source id.
func TestReplaceConfigMemberForSourceCanonical(t *testing.T) {
	t.Parallel()

	in := `{"modules":[` +
		`{"name":"eth-to-cosmos","src_chain":"ethereum","dst_chain":"cosmos","config":{"ics26_client_id":"chain-a","ics07_client":"0xETH"}},` +
		`{"name":"cosmos-to-eth","src_chain":"cosmos","dst_chain":"ethereum","config":{"ics26_client_id":"chain-a","ics07_client":"0xAAA"}}` +
		`]}`
	out, err := replaceConfigMemberForSource([]byte(in), "chain-a", "ics07_client", "0xNEW")
	if err != nil {
		t.Fatalf("replaceConfigMemberForSource() error = %v", err)
	}
	got := string(out)
	if !strings.Contains(got, `"ics26_client_id":"chain-a","ics07_client":"0xNEW"`) {
		t.Fatalf("canonical cosmos->eth module not updated, got: %s", got)
	}
	// The eth->cosmos module (same source id, wrong direction) must be untouched.
	if !strings.Contains(got, `"ics07_client":"0xETH"`) {
		t.Fatalf("eth->cosmos module must be skipped, got: %s", got)
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

	// create-clients-eth also targets cosmos_to_l2 destinations (the L2-side
	// SpectreClient is created by the same flow).
	t.Run("cosmos_to_l2 selected by id", func(t *testing.T) {
		mixed := &appConfig{
			CosmosToEthConfigs: []cosmosToEthConfig{{ICS26ClientID: "eth-a", TmRpcUrl: "a"}},
			CosmosToL2Configs:  []cosmosToEthConfig{{ICS26ClientID: "l2-arb", TmRpcUrl: "l2"}},
		}
		got, err := selectSource(mixed, "l2-arb")
		if err != nil {
			t.Fatalf("selectSource error = %v", err)
		}
		if got.CosmosToEthConfig.TmRpcUrl != "l2" {
			t.Fatalf("selected wrong dest: %+v", got.CosmosToEthConfig)
		}
	})
	t.Run("sole cosmos_to_l2 selected when empty id", func(t *testing.T) {
		l2Only := &appConfig{
			CosmosToL2Configs: []cosmosToEthConfig{{ICS26ClientID: "l2-solo", TmRpcUrl: "l2s"}},
		}
		got, err := selectSource(l2Only, "")
		if err != nil {
			t.Fatalf("selectSource error = %v", err)
		}
		if got.CosmosToEthConfig.TmRpcUrl != "l2s" {
			t.Fatalf("selected wrong dest: %+v", got.CosmosToEthConfig)
		}
	})
	t.Run("ambiguous across eth+l2 errors", func(t *testing.T) {
		mixed := &appConfig{
			CosmosToEthConfigs: []cosmosToEthConfig{{ICS26ClientID: "eth-a", TmRpcUrl: "a"}},
			CosmosToL2Configs:  []cosmosToEthConfig{{ICS26ClientID: "l2-arb", TmRpcUrl: "l2"}},
		}
		_, err := selectSource(mixed, "")
		if err == nil || !strings.Contains(err.Error(), "--source") {
			t.Fatalf("error = %v, want ambiguous-source error", err)
		}
	})
}

// An eth_to_cosmos module without a websocket used to start and then run half-dead:
// the forward direction relayed, SubscribeEth bailed out after one log line, and no
// acknowledgement was ever picked up. Refuse to START on that config.
func TestValidateRelayStartupConfigRequiresEthWs(t *testing.T) {
	const cfg = `{
		"modules": [
			{"name": "cosmos_to_eth", "src_chain": "cosmos", "dst_chain": "ethereum",
			 "config": {"tm_rpc_url": "http://localhost:26657", "eth_rpc_url": "http://localhost:8545",
			            "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "chain-a"}},
			{"name": "eth_to_cosmos", "src_chain": "ethereum", "dst_chain": "cosmos",
			 "config": {"eth_beacon_api_url": "http://localhost:5052"}}
		]
	}`
	appCfg, err := loadConfig(writeTempConfig(t, cfg))
	if err != nil {
		t.Fatalf("loadConfig must accept this — create-clients-* never read eth_ws_url: %v", err)
	}
	err = validateRelayStartupConfig(appCfg)
	if err == nil {
		t.Fatal("start accepted eth_to_cosmos without eth_ws_url; the ETH→Cosmos direction cannot run")
	}
	if !strings.Contains(err.Error(), "eth_ws_url") {
		t.Fatalf("error should name the missing field, got: %v", err)
	}
}

// A first-entry-only check let every later source reach the same half-dead state
// the guard exists to prevent. `start` now rejects a second source earlier
// (validateSingleRelayPair), so this is the function's own contract rather than a
// reachable deployment: it validates every entry it is handed.
func TestValidateRelayStartupConfigChecksEverySource(t *testing.T) {
	const cfg = `{
		"modules": [
			{"name": "cosmos_to_eth_a", "src_chain": "cosmos", "dst_chain": "ethereum",
			 "config": {"tm_rpc_url": "http://localhost:26657", "eth_rpc_url": "http://localhost:8545",
			            "eth_ws_url": "ws://localhost:8546",
			            "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "chain-a"}},
			{"name": "cosmos_to_eth_b", "src_chain": "cosmos", "dst_chain": "ethereum",
			 "config": {"tm_rpc_url": "http://localhost:36657", "eth_rpc_url": "http://localhost:8545",
			            "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "chain-b"}},
			{"name": "eth_to_cosmos", "src_chain": "ethereum", "dst_chain": "cosmos",
			 "config": {"eth_beacon_api_url": "http://localhost:5052"}}
		]
	}`
	appCfg, err := loadConfig(writeTempConfig(t, cfg))
	if err != nil {
		t.Fatalf("loadConfig must accept this: %v", err)
	}
	err = validateRelayStartupConfig(appCfg)
	if err == nil {
		t.Fatal("start accepted a second source with no eth_ws_url; its ETH→Cosmos leg cannot run")
	}
	if !strings.Contains(err.Error(), "eth_ws_url") {
		t.Fatalf("error should name the missing field, got: %v", err)
	}
	// The operator has to know WHICH source is short, not just that one is.
	if !strings.Contains(err.Error(), "chain-b") {
		t.Fatalf("error should identify the offending source, got: %v", err)
	}
}

// The client-creation commands share loadConfig but never read eth_ws_url
// (build_source.go, reached only from `start`, is its sole consumer). Rejecting a
// missing websocket at load time blocked create-clients-cosmos on a field it does
// not use — the config below must load, and only `start` may refuse it.
func TestLoadConfigAcceptsMissingEthWsForClientCreation(t *testing.T) {
	const cfg = `{
		"modules": [
			{"name": "cosmos_to_eth", "src_chain": "cosmos", "dst_chain": "ethereum",
			 "config": {"tm_rpc_url": "http://localhost:26657", "eth_rpc_url": "http://localhost:8545",
			            "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "chain-a"}},
			{"name": "eth_to_cosmos", "src_chain": "ethereum", "dst_chain": "cosmos",
			 "config": {"eth_beacon_api_url": "http://localhost:5052"}}
		]
	}`
	if _, err := loadConfig(writeTempConfig(t, cfg)); err != nil {
		t.Fatalf("loadConfig rejected a config create-clients-cosmos can run: %v", err)
	}
}

// A configured ETH→Cosmos direction needs a beacon REST endpoint for finality
// data and Ethereum light-client updates. Treating an empty value as if the
// module did not exist makes the relayer silently omit that direction.
func TestLoadConfigRejectsEthToCosmosWithoutBeaconURL(t *testing.T) {
	const cfg = `{
		"modules": [
			{"name": "eth_to_cosmos", "src_chain": "ethereum", "dst_chain": "cosmos", "config": {}}
		]
	}`
	_, err := loadConfig(writeTempConfig(t, cfg))
	if err == nil {
		t.Fatal("loadConfig accepted eth_to_cosmos without eth_beacon_api_url")
	}
	if !strings.Contains(err.Error(), "eth_beacon_api_url") {
		t.Fatalf("error should name the missing field, got: %v", err)
	}
}

// The same config with a websocket must load AND start.
func TestLoadConfigAcceptsEthToCosmosWithWs(t *testing.T) {
	const cfg = `{
		"modules": [
			{"name": "cosmos_to_eth", "src_chain": "cosmos", "dst_chain": "ethereum",
			 "config": {"tm_rpc_url": "http://localhost:26657", "eth_rpc_url": "http://localhost:8545",
			            "eth_ws_url": "ws://localhost:8546",
			            "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "chain-a"}},
			{"name": "eth_to_cosmos", "src_chain": "ethereum", "dst_chain": "cosmos",
			 "config": {"eth_beacon_api_url": "http://localhost:5052"}}
		]
	}`
	appCfg, err := loadConfig(writeTempConfig(t, cfg))
	if err != nil {
		t.Fatalf("loadConfig rejected a complete config: %v", err)
	}
	if err := validateRelayStartupConfig(appCfg); err != nil {
		t.Fatalf("start rejected a complete config: %v", err)
	}
}
