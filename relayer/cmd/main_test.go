package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
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

	// The per-path Ethereum example, not the catalogue: config.example.json lists
	// every module across every path and loadConfig refuses it for that reason
	// (TestShippedExamplesCarryOnePathEach pins that).
	cfg, err := loadConfig(filepath.Join("..", "config.ethereum.example.json"))
	if err != nil {
		t.Fatalf("loadConfig(config.ethereum.example.json) error = %v", err)
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

// A second cosmos_to_eth module is a second relay path, and loadConfig refuses
// it before parsing either in full -- naming both, so the operator knows which
// two to split, and showing the split command. This is the loadConfig-level
// half of TestRejectExtraRelayPaths: that the census is wired in at all.
func TestLoadConfigRejectsASecondCosmosSource(t *testing.T) {
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
	_, err := loadConfig(configPath)
	if err == nil {
		t.Fatal("loadConfig accepted two cosmos_to_eth modules; one process serves one path")
	}
	for _, want := range []string{"chain-a", "chain-b", "--config"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error must contain %q so the operator knows what to split and how; got: %v", want, err)
		}
	}
}

// Two cosmos_to_eth modules sharing an ics26_client_id used to be refused as a
// duplicate id (their ETH event streams, filtered by that id, would cross-feed).
// They are still refused -- as two paths, which subsumes the duplicate case.
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
	if err == nil || !strings.Contains(err.Error(), "more than one relay path") {
		t.Fatalf("loadConfig error = %v, want the second path refused", err)
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

// selectCosmosSource is the promotion half of what --source used to do: the
// create-clients and submit-misbehaviour flows read CosmosToEthConfig, and a
// cosmos_to_l2 module is the same struct pointed at the L2, so a file declaring
// only that module has it promoted into the field those flows read. There is no
// disambiguation half any more -- loadConfig refuses a file with two paths.
func TestSelectCosmosSource(t *testing.T) {
	t.Parallel()

	t.Run("a cosmos_to_eth module is read as is", func(t *testing.T) {
		single := &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{ICS26ClientID: "solo", TmRpcUrl: "s"},
		}
		if got := selectCosmosSource(single); got.CosmosToEthConfig.TmRpcUrl != "s" {
			t.Fatalf("selected wrong source: %+v", got.CosmosToEthConfig)
		}
	})

	// create-clients-eth also targets cosmos_to_l2 destinations (the L2-side
	// SpectreClient is created by the same flow).
	t.Run("a sole cosmos_to_l2 module is promoted", func(t *testing.T) {
		l2Only := &appConfig{
			HasCosmosToL2: true, CosmosToL2Config: cosmosToEthConfig{ICS26ClientID: "l2-solo", TmRpcUrl: "l2s"},
		}
		got := selectCosmosSource(l2Only)
		if got.CosmosToEthConfig.TmRpcUrl != "l2s" {
			t.Fatalf("selected wrong dest: %+v", got.CosmosToEthConfig)
		}
		// A copy, not a write into the caller's config: the L2 module must not
		// start looking like a cosmos_to_eth one to anything that reads the flag.
		if l2Only.HasCosmosToEth || l2Only.CosmosToEthConfig.TmRpcUrl != "" {
			t.Fatalf("selectCosmosSource mutated its input: %+v", l2Only)
		}
	})

	t.Run("nothing declared is returned unchanged", func(t *testing.T) {
		empty := &appConfig{}
		if got := selectCosmosSource(empty); got != empty {
			t.Fatal("a config with no Cosmos→EVM module must come back as is, for the command's own preflight to report")
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

// The prefix has to land AFTER the timestamp. Without log.Lmsgprefix the
// standard logger writes it first, so every line starts with the path and the
// times stop lining up as a column — which is how a merged log is actually read.
func TestStampRelayPathOnLogsPutsThePathAfterTheTimestamp(t *testing.T) {
	origFlags, origPrefix := log.Flags(), log.Prefix()
	origOut := log.Writer()
	t.Cleanup(func() {
		log.SetFlags(origFlags)
		log.SetPrefix(origPrefix)
		log.SetOutput(origOut)
	})

	var buf bytes.Buffer
	log.SetOutput(&buf)
	stampRelayPath(zap.NewNop(), "cosmos<->eth/cosmoshub-1")
	log.Printf("[Subscribe] send_packet received: seq=7")

	line := buf.String()
	if !strings.Contains(line, "[cosmos<->eth/cosmoshub-1]") {
		t.Fatalf("line %q does not carry the path", line)
	}
	if !strings.Contains(line, "[Subscribe] send_packet received: seq=7") {
		t.Fatalf("line %q lost the original message", line)
	}
	pathAt := strings.Index(line, "[cosmos<->eth/cosmoshub-1]")
	if pathAt == 0 {
		t.Fatalf("the path is written before the timestamp; log.Lmsgprefix is missing: %q", line)
	}
	// The timestamp must be everything before the path, not a substring of it.
	if stamp := strings.TrimSpace(line[:pathAt]); stamp == "" {
		t.Fatalf("nothing precedes the path, so the timestamp was dropped: %q", line)
	}
}

// The other half of the same call: the injected zap logger must carry the path
// too. Found in review -- the first version stamped only log.Default(), so every
// lifecycle line Start emits through zap ("Relayer started", "shutdown
// requested", "clients stopped") stayed unattributable in a merged log, and the
// test then in place only looked at log.Prefix() so it could not see that.
func TestStampRelayPathScopesTheZapLoggerToo(t *testing.T) {
	origFlags, origPrefix, origOut := log.Flags(), log.Prefix(), log.Writer()
	t.Cleanup(func() {
		log.SetFlags(origFlags)
		log.SetPrefix(origPrefix)
		log.SetOutput(origOut)
	})
	log.SetOutput(&bytes.Buffer{})

	core, logs := observer.New(zapcore.InfoLevel)
	scoped := stampRelayPath(zap.New(core), "cosmos<->eth/cosmoshub-1")
	scoped.Info("Relayer started")

	entries := logs.All()
	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}
	got, ok := entries[0].ContextMap()[relayPathLogField]
	if !ok {
		t.Fatalf("the zap line carries no %q field: %+v; a merged log cannot attribute it "+
			"to a relay path", relayPathLogField, entries[0].ContextMap())
	}
	if got != "cosmos<->eth/cosmoshub-1" {
		t.Fatalf("%s = %v, want the relay path", relayPathLogField, got)
	}
}

// stringLiterals returns the contents of each double-quoted literal on a line,
// so a check meant for format strings cannot read the code around them.
//
// Escapes are honoured only far enough to keep \" from ending a literal early;
// this reads Go source for one narrow property, not to compile it.
func stringLiterals(line string) []string {
	var out []string
	inside, escaped := false, false
	start := 0
	for i, r := range line {
		switch {
		case escaped:
			escaped = false
		case r == '\\' && inside:
			escaped = true
		case r == '"':
			if inside {
				out = append(out, line[start:i])
			} else {
				start = i + 1
			}
			inside = !inside
		}
	}
	return out
}

// Two label families are the target: [<direction> <Work>] and the closed
// lowercase list. Neither can contain "][", which is how a caller fakes a second
// tier when the function it calls is missing a parameter (subscriber/event.go
// passed "SubscribeCosmos][recovery" as a label). This fails on the forgery, not
// on the confusing log line it produces months later.
//
// [ATTENTION] is exempt, and is the one exemption. It is not a forged label: it
// is a fixed severity tier this repository already uses in 12 files, arriving
// over #377, #395, #415 and #457, and it is written the way this test asks for --
// the label is a parameter, the tier a literal after it. Forbidding it would mean
// rewriting a settled convention to satisfy a gate aimed at something else.
func TestNoLogLabelForgesASecondTier(t *testing.T) {
	root := ".."
	var offenders []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if name := info.Name(); name == "bindings" || name == "third_party" || name == "bin" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(src), "\n") {
			trimmed := strings.TrimSpace(line)
			if !strings.Contains(line, "][") || strings.HasPrefix(trimmed, "//") {
				// A comment describing a past forgery is not one.
				continue
			}
			// "][" is ordinary Go outside a string literal (`map[k][]byte`,
			// `x[i][j]`). Only a quoted occurrence is a forged label, so the
			// literals are extracted and nothing else on the line is read.
			//
			// Splitting on the first quote and scanning the remainder is not the
			// same thing, and the difference is a false positive: on
			//   dst.UpdateClient(ctx, "client-0", chain.ClientUpdate{Payloads: [][]byte{{1}}})
			// the [][]byte sits after a string literal and was reported as a
			// forged label.
			for _, literal := range stringLiterals(line) {
				// Drop the one allowed tier before looking, so "[%s][ATTENTION][x]"
				// is still caught: only the exemption itself is removed.
				if rest := strings.ReplaceAll(literal, "][ATTENTION]", "]"); strings.Contains(rest, "][") {
					offenders = append(offenders, fmt.Sprintf("%s:%d: %s", path, i+1, trimmed))
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	if len(offenders) > 0 {
		t.Fatalf("log labels forging a second tier:\n%s\n\nA label that needs a second level needs a "+
			"parameter, not a closing bracket in the middle of the format string.", strings.Join(offenders, "\n"))
	}
}

// The two noise words are gone and must stay gone: the direction is already the
// label (or the process prefix), so "relay" and "adapter" said nothing.
func TestNoLogLabelRepeatsTheLayerName(t *testing.T) {
	for _, banned := range []string{`"[relay `, `"[adapter `, `"\n[relay `, `"\n[adapter `} {
		out, err := exec.Command("grep", "-rl", "--include=*.go", banned, "..").Output()
		if err == nil && len(strings.TrimSpace(string(out))) > 0 {
			t.Errorf("label %q is back in:\n%s", strings.TrimSpace(banned), out)
		}
	}
}

// The point of E3 is one concept, one name. Four names for the timeout scanner
// and three for the Ethereum sender is not a formatting problem: an operator
// grepping [SendEthTx] silently missed a third of the sender's own lines.
//
// This pins the closed set. Adding a label is deliberate and cheap -- add it
// here too. Re-inventing a name for something already named is what this stops.
func TestLogLabelsStayOnTheAgreedList(t *testing.T) {
	allowed := map[string]bool{
		// [<direction> <Work>] -- direction is literal where the code knows it,
		// and "%s" where it comes from the module's own name.
		"eth->cosmos Subscribe": true, "l2->cosmos Subscribe": true,
		"cosmos->eth Relay": true, "cosmos->l2 Relay": true, "l2->cosmos Relay": true,
		"cosmos->eth UpdateClient": true, "cosmos->l2 UpdateClient": true,
		"eth->cosmos UpdateClient": true,
		"%s Relay":                 true, "%s UpdateClient": true,

		// [<Work>] -- the emitting code serves more than one direction, so the
		// chain stays in the name: it is the only thing left that separates them.
		"SubscribeCosmos": true, "SendEthTx": true, "SendCosmosTx": true,
		"CreateEthClient": true, "CreateCosmosClient": true,
		"UpdateCosmosClient": true, "UpdateEthClient": true, "RefreshCosmosClient": true,
		"CosmosTimeoutScan": true, "%sTimeoutScan": true,
		"Misbehaviour": true, "PendingTracker": true,
		// Added on purpose, which is what the message below asks for. #452's
		// enumeration backstop and #454's owed-ack ledger both landed after the
		// E3 renaming pass ran, so their labels were never put through it:
		// FlushCosmos is chain-qualified for the same reason SubscribeCosmos is
		// (cosmos.NewSource serves both cosmos->eth and cosmos->l2), and AckWatch
		// replaces a bare [adapter], which is one of the two noise words E3
		// removed.
		"FlushCosmos": true, "AckWatch": true,

		// The closed lowercase list: shared by every path, so a direction would
		// be a lie.
		// `create-clients` sits here for the same reason as `start`: it names a CLI
		// command that serves both directions, so a direction in the label would be
		// a lie.
		"bench": true, "prover": true, "BatchBuilder": true, "QueueState": true, "start": true,
		// The status endpoint serves the whole process, so a direction would be a
		// lie -- the same reason start and create-clients sit here.
		"status":         true,
		"create-clients": true,

		// Deleted by #417 along with the unbound-header fallback it belongs to;
		// touching it here would only conflict with that PR.
		"%s": true,
	}

	labelRe := regexp.MustCompile(`\.(?:Printf|Println|Fatalf|Infof|Warnf|Errorf|Debugf)\("(?:\\n)?\[([^\]"]+)\]`)
	var offenders []string
	err := filepath.Walk("..", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if n := info.Name(); n == "bindings" || n == "third_party" || n == "bin" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(src), "\n") {
			for _, m := range labelRe.FindAllStringSubmatch(line, -1) {
				if !allowed[m[1]] {
					offenders = append(offenders, fmt.Sprintf("%s:%d: [%s]", path, i+1, m[1]))
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	if len(offenders) > 0 {
		t.Fatalf("log labels outside the agreed set:\n%s\n\nEither reuse the name this concept already "+
			"has, or add the new one to this list on purpose.", strings.Join(offenders, "\n"))
	}
}
