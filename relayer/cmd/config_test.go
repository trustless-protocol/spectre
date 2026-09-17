package main

import (
	attestorpb "attestor/types/attestor"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"relayer/chain"
	"relayer/rpcmock"
	"strings"
	"testing"

	cmtp2p "github.com/cometbft/cometbft/p2p"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/ethereum/go-ethereum/common"
)

// One relayer process serves exactly one source → destination pair (03-Relayer
// §2.1 requirement 5). Before this guard, `start` looped over every configured
// module and ran them in one process, which is what made a shared nonce cache and
// a shared error path necessary in the first place. A7 moved the guard to parse
// time: loadConfig takes a census of the modules and refuses a second path
// before any module is parsed in full, so nothing downstream ever holds a list.
//
// The census is taken from the raw modules, so the cases below carry only the
// identity each module is named by; that a fully valid single-path file loads
// is TestLoadConfig/shipped examples carry one path each's job.
func TestRejectExtraRelayPaths(t *testing.T) {
	c2e := func(ids ...string) []string {
		out := make([]string, 0, len(ids))
		for _, id := range ids {
			out = append(out, `{"name":"cosmos_to_eth","src_chain":"cosmos","dst_chain":"ethereum","config":{"ics26_client_id":"`+id+`"}}`)
		}
		return out
	}
	c2l2 := func(ids ...string) []string {
		out := make([]string, 0, len(ids))
		for _, id := range ids {
			out = append(out, `{"name":"cosmos_to_l2","src_chain":"cosmos","dst_chain":"opstack","config":{"ics26_client_id":"`+id+`"}}`)
		}
		return out
	}
	l2src := func(chains ...string) []string {
		out := make([]string, 0, len(chains))
		for _, c := range chains {
			out = append(out, `{"name":"l2_to_cosmos","src_chain":"`+c+`","dst_chain":"cosmos","config":{"attestor_src_chain":"`+c+`"}}`)
		}
		return out
	}
	join := func(groups ...[]string) []string {
		var out []string
		for _, g := range groups {
			out = append(out, g...)
		}
		return out
	}

	tests := []struct {
		name    string
		modules []string
		// wantNames are the identities the error must name so the operator knows
		// which paths to split apart; empty means the config must be accepted.
		wantNames []string
	}{
		{
			name:    "one cosmos_to_eth",
			modules: c2e("chain-a"),
		},
		{
			name:    "one cosmos_to_l2 with its return leg is one pair",
			modules: join(c2l2("08-wasm-3"), l2src("opstack")),
		},
		{
			name:    "forward-only cosmos_to_l2",
			modules: c2l2("08-wasm-3"),
		},
		{
			name:    "two cosmos_to_eth sources",
			modules: c2e("chain-a", "chain-b"),
			// Naming both is the point: "config has 2 sources" leaves the operator
			// guessing which two, and file order is not what they edited by.
			wantNames: []string{"chain-a", "chain-b"},
		},
		{
			name:      "two cosmos_to_l2 destinations",
			modules:   c2l2("08-wasm-3", "08-wasm-4"),
			wantNames: []string{"08-wasm-3", "08-wasm-4"},
		},
		{
			name:      "two l2_to_cosmos sources",
			modules:   l2src("opstack", "arbitrum"),
			wantNames: []string{"opstack", "arbitrum"},
		},
		{
			name:      "cosmos_to_eth alongside cosmos_to_l2 is two pairs",
			modules:   join(c2e("chain-a"), c2l2("08-wasm-3")),
			wantNames: []string{"chain-a", "08-wasm-3"},
		},
		{
			name:      "cosmos_to_eth alongside an l2_to_cosmos return leg is two pairs",
			modules:   join(c2e("chain-a"), l2src("opstack")),
			wantNames: []string{"chain-a", "opstack"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jc jsonConfig
			if err := json.Unmarshal([]byte(`{"modules":[`+strings.Join(tt.modules, ",")+`]}`), &jc); err != nil {
				t.Fatalf("parse fixture: %v", err)
			}
			census, err := censusRelayPaths(jc.Modules)
			if err != nil {
				t.Fatalf("census: %v", err)
			}
			err = rejectExtraRelayPaths("config.json", census)
			if len(tt.wantNames) == 0 {
				if err != nil {
					t.Fatalf("single-pair config rejected: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("start accepted more than one relay path in one process")
			}
			for _, name := range tt.wantNames {
				if !strings.Contains(err.Error(), name) {
					t.Errorf("error must name path %q so the operator knows what to split; got: %v", name, err)
				}
			}
			// Telling the operator what is wrong without telling them the way out
			// leaves them editing the config at random.
			if !strings.Contains(err.Error(), "--config") {
				t.Errorf("error must show the split command; got: %v", err)
			}
		})
	}
}

// `start` with nothing to relay is a config mistake, not an idle process.
func TestRequireRelayPath_RejectsEmptyConfig(t *testing.T) {
	err := requireRelayPath(&appConfig{}, "/etc/relayer/config.json")
	if err == nil {
		t.Fatal("start accepted a config with no relay module")
	}
	if !strings.Contains(err.Error(), "/etc/relayer/config.json") {
		t.Errorf("error must name the config file it read; got: %v", err)
	}
}

func TestLoadConfig(t *testing.T) {
	// TestLoadConfig/L2 source confirms an l2_to_cosmos module parses into
	// appConfig.L2ToCosmosConfig (no longer the fail-loud "not wired" error).
	t.Run("L2 source", func(t *testing.T) {
		raw := `{"modules":[{"name":"op","src_chain":"opstack","dst_chain":"cosmos","config":{
		"l1_rpc_url":"http://l1","l2_rpc_url":"http://l2","tm_rpc_url":"http://tm",
		"attestor_endpoints":[{"address":"127.0.0.1:3001","attestor_index":0}],"attestor_src_chain":"op-sepolia",
		"attestors":{"public_keys":["gojj3XQJ8ZX9UtstPLpdcspnCb8dlBIb83SIAbQPb1w="],"threshold":1},
		"eth_beacon_api_url":"http://beacon",
		"l2_wasm_client_id":"08-wasm-1","l2_ics26_client_id":"client-0","head_kind":"safe",
		"rollup_profile":{"common":{"l2_chain_id":10,"l2_router":"0x1111111111111111111111111111111111111111"}}}}]}`
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
	})

	// The shipped examples encode the split the guard now enforces: one runnable file
	// per path, plus a catalogue that carries every module and is not meant to be
	// started. Without this test the two drift silently — the catalogue is the file
	// whose name most invites `--config config.example.json`.
	t.Run("shipped examples carry one path each", func(t *testing.T) {
		t.Parallel()

		for file := range pathExamples {
			t.Run(file, func(t *testing.T) {
				t.Parallel()
				cfg, err := loadConfig(filepath.Join("..", file))
				if err != nil {
					t.Fatalf("loadConfig(%s): %v", file, err)
				}
				if err := requireRelayPath(cfg, file); err != nil {
					t.Fatalf("%s is a runnable example but start rejects it: %v", file, err)
				}
			})
		}

		// Refused at load, not merely at start: every command that reads a config
		// shares loadConfig, and none of them can do anything sensible with a file
		// that names several paths.
		t.Run("catalogue is not loadable", func(t *testing.T) {
			t.Parallel()
			_, err := loadConfig(filepath.Join("..", "config.example.json"))
			if err == nil {
				t.Fatal("the catalogue lists every path; loadConfig must send the operator to a per-path file instead")
			}
			if !strings.Contains(err.Error(), "--config") {
				t.Errorf("the rejection must show the way out; got: %v", err)
			}
		})
	})

	// TestLoadConfig/path examples match catalogue keeps the per-path examples byte-identical to
	// the catalogue they are cut from. Without this they drift: the file these
	// replaced (config.arbitrum-only.json) had gone stale on log_scan_chunk and
	// include_provisional and carried the addresses of a long-finished devnet run,
	// and nothing in the repo referenced it so nothing caught that.
	t.Run("path examples match catalogue", func(t *testing.T) {
		t.Parallel()

		catalogue := readModules(t, "config.example.json")

		for file, wantModules := range pathExamples {
			t.Run(file, func(t *testing.T) {
				t.Parallel()

				got := readModules(t, file)
				if len(got) != len(wantModules) {
					t.Fatalf("%s has %d modules, want %d — it must carry exactly one relay pair",
						file, len(got), len(wantModules))
				}
				for _, name := range wantModules {
					mine, ok := got[name]
					if !ok {
						t.Errorf("%s is missing module %q", file, name)
						continue
					}
					theirs, ok := catalogue[name]
					if !ok {
						t.Errorf("config.example.json has no module %q to compare against", name)
						continue
					}
					var a, b any
					if err := json.Unmarshal(mine, &a); err != nil {
						t.Fatalf("parse %s module %q: %v", file, name, err)
					}
					if err := json.Unmarshal(theirs, &b); err != nil {
						t.Fatalf("parse catalogue module %q: %v", name, err)
					}
					if !reflect.DeepEqual(a, b) {
						t.Errorf("%s module %q has drifted from config.example.json.\n"+
							"Regenerate it rather than hand-editing:\n"+
							"  cd relayer && jq '{modules: [.modules[] | select(.name|test(\"<filter>\"))], batch: .batch}' \\\n"+
							"    config.example.json > %s", file, name, file)
					}
				}
			})
		}
	})

	// TestLoadConfig/path examples load is the point of these files: unlike the catalogue, each
	// one is a config the relayer can actually start from.
	t.Run("path examples load", func(t *testing.T) {
		t.Parallel()

		for file := range pathExamples {
			t.Run(file, func(t *testing.T) {
				t.Parallel()
				if _, err := loadConfig(filepath.Join("..", file)); err != nil {
					t.Fatalf("loadConfig(%s) error = %v", file, err)
				}
			})
		}
	})

	t.Run("example rollup profiles use only fresh V2 schema", func(t *testing.T) {
		t.Parallel()

		examples := map[string]int{
			"config.example.json":          3,
			"config.op.example.json":       1,
			"config.base.example.json":     1,
			"config.arbitrum.example.json": 1,
			"op-l2-config.example.json":    1,
			"base-l2-config.example.json":  1,
			"arb-l2-config.example.json":   1,
		}
		wantKeys := map[string]struct{}{
			"l2_chain_id": {}, "l2_router": {}, "commitment_slot": {},
			"profile_version": {}, "l2_header_fork": {},
		}
		allowedVersions := map[string]struct{}{
			"op_attestor_v1": {}, "base_attestor_v1": {}, "arbitrum_attestor_v1": {},
		}

		var collectProfiles func(any, *[]map[string]any)
		collectProfiles = func(value any, profiles *[]map[string]any) {
			switch value := value.(type) {
			case map[string]any:
				if version, ok := value["profile_version"].(string); ok {
					if _, allowed := allowedVersions[version]; allowed {
						*profiles = append(*profiles, value)
					}
				}
				for _, child := range value {
					collectProfiles(child, profiles)
				}
			case []any:
				for _, child := range value {
					collectProfiles(child, profiles)
				}
			}
		}

		for file, wantCount := range examples {
			file, wantCount := file, wantCount
			t.Run(file, func(t *testing.T) {
				t.Parallel()
				raw, err := os.ReadFile(filepath.Join("..", file))
				if err != nil {
					t.Fatal(err)
				}
				var document any
				if err := json.Unmarshal(raw, &document); err != nil {
					t.Fatal(err)
				}
				var profiles []map[string]any
				collectProfiles(document, &profiles)
				if len(profiles) != wantCount {
					t.Fatalf("found %d authenticated rollup profiles, want %d", len(profiles), wantCount)
				}
				for _, profile := range profiles {
					if len(profile) != len(wantKeys) {
						t.Fatalf("profile has %d fields, want exactly five: %#v", len(profile), profile)
					}
					for key := range profile {
						if _, ok := wantKeys[key]; !ok {
							t.Fatalf("profile contains stale/unknown field %q", key)
						}
					}
				}
			})
		}
	})

	// TestLoadConfig/legacy src chain fallback: a legacy cosmos_to_eth module with no
	// ics26_client_id falls back to src_chain (back-compat).
	t.Run("legacy src chain fallback", func(t *testing.T) {
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
	})

	// TestLoadConfig/canonical no fallback: a canonical module does NOT overload
	// src_chain (a chain kind) as the ics26 fallback.
	t.Run("canonical no fallback", func(t *testing.T) {
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
	})

	// TestLoadConfig/unknown module fails loud: a misspelled module is a hard error, not
	// a silent drop.
	t.Run("unknown module fails loud", func(t *testing.T) {
		t.Parallel()
		path := writeTempConfig(t, `{"modules":[{"name":"cosmos_to_ethh","config":{}}]}`)
		if _, err := loadConfig(path); err == nil {
			t.Fatal("expected loadConfig to fail on an unrecognized module, got nil")
		}
	})

	// TestLoadConfig/L2 to cosmos not wired: the optimistic L2→Cosmos direction is not on
	// this branch, so it still fails loud.
	t.Run("L2 to cosmos not wired", func(t *testing.T) {
		t.Parallel()
		path := writeTempConfig(t, `{"modules":[{"name":"op","src_chain":"opstack","dst_chain":"cosmos","config":{}}]}`)
		if _, err := loadConfig(path); err == nil {
			t.Fatal("expected loadConfig to fail on an unwired l2_to_cosmos direction, got nil")
		}
	})

	// TestLoadConfig/cosmos to L2 parses: a cosmos→L2 module parses into CosmosToL2Config
	// (same schema as cosmos_to_eth), pointed at the L2 endpoint.
	t.Run("cosmos to L2 parses", func(t *testing.T) {
		t.Parallel()
		path := writeTempConfig(t, `{"modules":[
		{"name":"arb-dst","src_chain":"cosmos","dst_chain":"arbitrum","config":{"ics26_client_id":"arb-client-0","tm_rpc_url":"http://localhost:26657","eth_rpc_url":"http://localhost:8547","ics26_address":"0x80741a37e3644612f0465145c9709a90b6d77ee3","spectre_client":"0x80741a37e3644612f0465145c9709a90b6d77ee3"}}
	]}`)
		cfg, err := loadConfig(path)
		if err != nil {
			t.Fatalf("loadConfig: %v", err)
		}
		if !cfg.HasCosmosToL2 {
			t.Fatal("HasCosmosToL2 = false, want the module recorded")
		}
		if got := cfg.CosmosToL2Config.ICS26ClientID; got != "arb-client-0" {
			t.Fatalf("ics26_client_id = %q, want arb-client-0", got)
		}
		if cfg.HasCosmosToEth {
			t.Fatal("HasCosmosToEth should be false for a cosmos_to_l2-only config")
		}
	})

	t.Run("cosmos to L2 validation uses module name", func(t *testing.T) {
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
	})

	// TestLoadConfig/cosmos to L2 duplicate client ID: two cosmos→L2 destinations must not
	// share an ics26_client_id.
	t.Run("cosmos to L2 duplicate client ID", func(t *testing.T) {
		t.Parallel()
		path := writeTempConfig(t, `{"modules":[
		{"name":"arb-dst","src_chain":"cosmos","dst_chain":"arbitrum","config":{"ics26_client_id":"dup","tm_rpc_url":"http://localhost:26657","eth_rpc_url":"http://localhost:8547","ics26_address":"0x80741a37e3644612f0465145c9709a90b6d77ee3","spectre_client":"0x80741a37e3644612f0465145c9709a90b6d77ee3"}},
		{"name":"op-dst","src_chain":"cosmos","dst_chain":"opstack","config":{"ics26_client_id":"dup","tm_rpc_url":"http://localhost:26657","eth_rpc_url":"http://localhost:8548","ics26_address":"0x80741a37e3644612f0465145c9709a90b6d77ee3","spectre_client":"0x80741a37e3644612f0465145c9709a90b6d77ee3"}}
	]}`)
		if _, err := loadConfig(path); err == nil {
			t.Fatal("expected loadConfig to fail on duplicate cosmos_to_l2 ics26_client_id, got nil")
		}
	})

	t.Run("rejects deprecated batch config", func(t *testing.T) {
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
	})

	// TestLoadConfig/ignores legacy server block pins the compatibility half of dropping
	// the dead `server` block: operators' existing config.json files still carry it,
	// and they must keep loading. encoding/json ignores unknown members and nothing
	// here sets DisallowUnknownFields, so this holds — the test exists to keep it that
	// way, because turning on strict decoding later would break every deployed config.
	t.Run("ignores legacy server block", func(t *testing.T) {
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
	})

	t.Run("example", func(t *testing.T) {
		t.Parallel()

		// The per-path Ethereum example, not the catalogue: config.example.json lists
		// every module across every path and loadConfig refuses it for that reason
		// (TestLoadConfig/shipped examples carry one path each pins that).
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
	})

	// TestLoadConfig/ignores pruned fields confirms old configs carrying the now-removed
	// eth_to_cosmos fields and dst_chain still load (unknown JSON keys are ignored),
	// preserving backward compatibility.
	t.Run("ignores pruned fields", func(t *testing.T) {
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
	})

	// A second cosmos_to_eth module is a second relay path, and loadConfig refuses
	// it before parsing either in full -- naming both, so the operator knows which
	// two to split, and showing the split command. This is the loadConfig-level
	// half of TestRejectExtraRelayPaths: that the census is wired in at all.
	t.Run("rejects a second cosmos source", func(t *testing.T) {
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
	})

	// Two cosmos_to_eth modules sharing an ics26_client_id used to be refused as a
	// duplicate id (their ETH event streams, filtered by that id, would cross-feed).
	// They are still refused -- as two paths, which subsumes the duplicate case.
	t.Run("rejects duplicate client ID", func(t *testing.T) {
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
	})

	// The client-creation commands share loadConfig but never read eth_ws_url
	// (build_source.go, reached only from `start`, is its sole consumer). Rejecting a
	// missing websocket at load time blocked create-clients-cosmos on a field it does
	// not use — the config below must load, and only `start` may refuse it.
	t.Run("accepts missing eth ws for client creation", func(t *testing.T) {
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
	})

	// A configured ETH→Cosmos direction needs a beacon REST endpoint for finality
	// data and Ethereum light-client updates. Treating an empty value as if the
	// module did not exist makes the relayer silently omit that direction.
	t.Run("rejects eth to cosmos without beacon URL", func(t *testing.T) {
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
	})

	// The same config with a websocket must load AND start.
	t.Run("accepts eth to cosmos with ws", func(t *testing.T) {
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
	})
}

// --- state directory ---------------------------------------------------------

func TestValidateStateDir_AcceptsAWritableDirAndLeavesNothingBehind(t *testing.T) {
	// The directory the process will really write into, passed in. It used to be
	// read from RELAYER_RECOVERY_STATE_FILE; that variable is now REFUSED at
	// startup, so validating it would judge a path nothing writes.
	dir := t.TempDir()
	stateDir := filepath.Join(dir, "state")

	if findings := validateStateDir(stateDir); len(findings) != 0 {
		t.Fatalf("expected no findings for a writable dir, got %+v", findings)
	}

	// The probe must not survive. A leftover file in the state directory is not
	// merely untidy: the next run's O_EXCL probe would fail on it and report a
	// writable directory as unwritable.
	entries, err := os.ReadDir(filepath.Join(dir, "state"))
	if err != nil {
		t.Fatalf("read state dir: %v", err)
	}
	if len(entries) != 0 {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("probe left files behind: %v", names)
	}
}

func TestValidateStateDir_RejectsADirectoryItCannotWrite(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores the mode bits this test relies on")
	}
	parent := t.TempDir()
	readOnly := filepath.Join(parent, "ro")
	if err := os.Mkdir(readOnly, 0o500); err != nil {
		t.Fatalf("create read-only dir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(readOnly, 0o700) })

	findings := validateStateDir(filepath.Join(readOnly, "nested"))
	if len(findings) != 1 {
		t.Fatalf("expected exactly one finding, got %+v", findings)
	}
	if findings[0].key != "state directory" {
		t.Fatalf("finding names %q, want the thing an operator would fix", findings[0].key)
	}
}

// --- cosmos chain id ---------------------------------------------------------

func TestValidateCosmosChainID(t *testing.T) {
	t.Run("fails before dialing when the variable is unset", func(t *testing.T) {
		node := rpcmock.NewCosmosNode(t)
		node.Handle("status", func(json.RawMessage) (any, error) {
			return statusResult("some-chain"), nil
		})
		t.Setenv("COSMOS_CHAIN_ID", "")

		findings := validateCosmosChainID(context.Background(), &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{TmRpcUrl: node.URL()},
		})

		if len(findings) != 1 || findings[0].key != "COSMOS_CHAIN_ID" {
			t.Fatalf("expected one COSMOS_CHAIN_ID finding, got %+v", findings)
		}
		// Nothing to compare against means nothing to ask, and asking anyway would
		// make an unset variable look like an endpoint problem in the log.
		if calls := node.Calls("status"); calls != 0 {
			t.Fatalf("queried the node %d times with no declared chain id; want 0", calls)
		}
	})

	t.Run("reports a definite disagreement", func(t *testing.T) {
		node := rpcmock.NewCosmosNode(t)
		node.Handle("status", func(json.RawMessage) (any, error) {
			return statusResult("served-chain"), nil
		})
		t.Setenv("COSMOS_CHAIN_ID", "declared-chain")

		findings := validateCosmosChainID(context.Background(), &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{TmRpcUrl: node.URL()},
		})

		if len(findings) != 1 {
			t.Fatalf("expected one finding, got %+v", findings)
		}
		// Both halves must appear: an operator has to know which of the two to fix,
		// and a message naming only one of them cannot tell them.
		for _, want := range []string{"declared-chain", "served-chain", node.URL()} {
			if !strings.Contains(findings[0].detail, want) {
				t.Fatalf("finding %q does not mention %q", findings[0].detail, want)
			}
		}
	})

	t.Run("accepts a match", func(t *testing.T) {
		node := rpcmock.NewCosmosNode(t)
		node.Handle("status", func(json.RawMessage) (any, error) {
			return statusResult("test-ibc-eth"), nil
		})
		t.Setenv("COSMOS_CHAIN_ID", "test-ibc-eth")

		if findings := validateCosmosChainID(context.Background(), &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{TmRpcUrl: node.URL()},
		}); len(findings) != 0 {
			t.Fatalf("expected no findings for a matching chain id, got %+v", findings)
		}
	})

	t.Run("tolerates an endpoint that does not answer", func(t *testing.T) {
		// A node with no "status" handler answers "method not found" — the shape of
		// an endpoint that is reachable but useless. Startup must survive it: an
		// unavailable node is an operations problem, and the relay loops retry.
		node := rpcmock.NewCosmosNode(t)
		t.Setenv("COSMOS_CHAIN_ID", "declared-chain")

		if findings := validateCosmosChainID(context.Background(), &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{TmRpcUrl: node.URL()},
		}); len(findings) != 0 {
			t.Fatalf("a non-answering endpoint must not fail startup, got %+v", findings)
		}
	})
}

func TestCosmosEndpointsInConfig_AreGatheredFromEveryDirectionWithoutDuplicates(t *testing.T) {
	cfg := &appConfig{
		HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{TmRpcUrl: "http://a"},
		HasCosmosToL2: true, CosmosToL2Config: cosmosToEthConfig{TmRpcUrl: "http://a"},
		HasL2ToCosmos: true, L2ToCosmosConfig: l2ToCosmosConfig{TmRpcUrl: "http://c"},
	}
	got := cosmosEndpointsInConfig(cfg)
	want := []string{"http://a", "http://c"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("gathered %v, want %v", got, want)
	}
}

// --- EVM deployments ---------------------------------------------------------

// deployedCode is any non-empty code answer; the check only asks whether there
// is code at all, never what it is.
const deployedCode = "0x60006000"

func evmConfigFor(node *rpcmock.EVMNode) cosmosToEthConfig {
	return cosmosToEthConfig{
		EthRpcUrl:         node.URL(),
		ICS26Address:      "0x1111111111111111111111111111111111111111",
		ICS26ClientID:     "client-0",
		SpectreClient:     "0x2222222222222222222222222222222222222222",
		SignatureVerifier: "0x3333333333333333333333333333333333333333",
		Membership:        "0x4444444444444444444444444444444444444444",
		Misbehaviour:      "0x5555555555555555555555555555555555555555",
		UpdateClient:      "0x6666666666666666666666666666666666666666",
	}
}

func deployAll(node *rpcmock.EVMNode, c cosmosToEthConfig) {
	for _, a := range []string{
		c.ICS26Address, c.SpectreClient, c.SignatureVerifier,
		c.Membership, c.Misbehaviour, c.UpdateClient,
	} {
		node.SetCode(common.HexToAddress(a), deployedCode)
	}
	// getClient's return value: one ABI-encoded non-zero address.
	node.SetCallResult("0x000000000000000000000000" + strings.Repeat("ab", 20))
}

func TestValidateEVMDeployments(t *testing.T) {
	t.Run("reports every address that holds no code", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		cfg := evmConfigFor(node)
		// Only the router is deployed, so the other five are the findings — and the
		// router being present is what lets the client-id check run at all.
		node.SetCode(common.HexToAddress(cfg.ICS26Address), deployedCode)
		node.SetCallResult("0x000000000000000000000000" + strings.Repeat("ab", 20))

		findings := validateEVMDeployments(context.Background(), &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cfg,
		})

		got := map[string]bool{}
		for _, f := range findings {
			got[f.key] = true
		}
		for _, key := range []string{
			"cosmos_to_eth.spectre_client", "cosmos_to_eth.signature_verifier",
			"cosmos_to_eth.membership", "cosmos_to_eth.misbehaviour",
			"cosmos_to_eth.update_client",
		} {
			if !got[key] {
				t.Fatalf("no finding for %s; got %+v", key, findings)
			}
		}
		if got["cosmos_to_eth.ics26_address"] {
			t.Fatalf("reported the deployed router as code-less: %+v", findings)
		}
		// This test used a presence map, which is why it could not see the router
		// check handing back the findings it was given -- five problems reported as
		// ten. Count as well as check.
		if len(findings) != len(got) {
			t.Fatalf("%d findings for %d distinct keys; some are reported twice: %+v",
				len(findings), len(got), findings)
		}
	})

	// validateRouterClientID receives the findings so far so it can stay quiet when
	// the router address is itself a finding. It must not hand them back: the caller
	// appends its return value, so returning the accumulator reports everything
	// twice. Found in review -- a module with five missing contracts came out as ten
	// problems.
	t.Run("findings are not reported twice", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		cfg := evmConfigFor(node)
		// Router deployed (so the client-id check runs) and nothing else, which is
		// the shape that accumulates findings BEFORE the router check.
		node.SetCode(common.HexToAddress(cfg.ICS26Address), deployedCode)
		node.SetCallResult("0x000000000000000000000000" + strings.Repeat("ab", 20))

		findings := validateEVMDeployments(context.Background(), &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cfg,
		})

		seen := map[string]int{}
		for _, f := range findings {
			seen[f.key]++
		}
		for key, n := range seen {
			if n > 1 {
				t.Fatalf("%q reported %d times (%d findings, %d distinct): an operator "+
					"reading this cannot tell how many problems there actually are",
					key, n, len(findings), len(seen))
			}
		}
	})

	// The mirror of TestValidateEVMDeployments/validates the ID startup will actually use,
	// and the bug is its exact inverse: buildCosmosToL2Dest reads ICS26ClientID
	// straight from the JSON and never consults the environment, so applying the
	// override to a cosmos_to_l2 entry validates an id that path will never use.
	// A stale value in .env then fails a good config, and a correct one hides a bad
	// JSON id. Validation must read whatever the BUILDER reads, per path.
	t.Run("cosmos to L2 validates the config ID not the env override", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		cfg := evmConfigFor(node)
		cfg.ICS26ClientID = "client-from-json"
		deployAll(node, cfg)
		// Every getClient reverts, so whichever id is asked about becomes a finding
		// naming that id -- which is how we see WHICH id was validated.
		node.FailCall(errors.New("execution reverted"))
		t.Setenv("ICS26_CLIENT_ID", "stale-from-dotenv")

		findings := validateEVMDeployments(context.Background(),
			&appConfig{HasCosmosToL2: true, CosmosToL2Config: cfg})

		if len(findings) == 0 {
			t.Fatal("no finding at all; the cosmos_to_l2 entry was not validated")
		}
		for _, f := range findings {
			if strings.Contains(f.detail, "stale-from-dotenv") {
				t.Fatalf("validated the env override on a cosmos_to_l2 entry: %s\n"+
					"buildCosmosToL2Dest never reads it, so this rejects a config the relay "+
					"would have run fine", f.detail)
			}
		}
		for _, f := range findings {
			if strings.Contains(f.detail, "client-from-json") {
				return
			}
		}
		t.Fatalf("no finding names the JSON client id, so nothing proves which id was checked: %+v", findings)
	})

	t.Run("accepts a fully deployed module", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		cfg := evmConfigFor(node)
		deployAll(node, cfg)

		if findings := validateEVMDeployments(context.Background(), &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cfg,
		}); len(findings) != 0 {
			t.Fatalf("expected no findings for a deployed module, got %+v", findings)
		}
	})

	t.Run("reports a client ID the router does not know", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		cfg := evmConfigFor(node)
		deployAll(node, cfg)
		// getClient reverts with IBCClientNotFound for an unregistered id
		// (contracts/utils/ICS02ClientUpgradeable.sol:78-80).
		node.FailCall(errors.New("execution reverted"))

		findings := validateEVMDeployments(context.Background(), &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cfg,
		})

		if len(findings) != 1 || findings[0].key != "cosmos_to_eth.ics26_client_id" {
			t.Fatalf("expected one ics26_client_id finding, got %+v", findings)
		}
		if !strings.Contains(findings[0].detail, "client-0") {
			t.Fatalf("finding %q does not name the client id", findings[0].detail)
		}
	})

	t.Run("tolerates a transport error from the router call", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		cfg := evmConfigFor(node)
		deployAll(node, cfg)
		// Not a revert. The node did not execute the call, so this says nothing
		// about the client id and must not stop startup.
		node.FailCall(errors.New("connection reset by peer"))

		if findings := validateEVMDeployments(context.Background(), &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cfg,
		}); len(findings) != 0 {
			t.Fatalf("a transport error must not fail startup, got %+v", findings)
		}
	})

	t.Run("does not report the client ID when the router itself is missing", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		cfg := cosmosToEthConfig{
			EthRpcUrl:     node.URL(),
			ICS26Address:  "0x1111111111111111111111111111111111111111",
			ICS26ClientID: "client-0",
		}
		// Nothing deployed, and getClient would revert. "The router has no code" and
		// "the router does not know this client" are then the same fact, and
		// reporting both sends an operator looking for a second problem.
		node.FailCall(errors.New("execution reverted"))

		findings := validateEVMDeployments(context.Background(), &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cfg,
		})

		if len(findings) != 1 || findings[0].key != "cosmos_to_eth.ics26_address" {
			t.Fatalf("expected only the missing-router finding, got %+v", findings)
		}
	})

	t.Run("tolerates an unreachable endpoint", func(t *testing.T) {
		cfg := cosmosToEthConfig{
			// Reserved TEST-NET-1 address (RFC 5737): routable syntax, no listener.
			EthRpcUrl:     "http://192.0.2.1:8545",
			ICS26Address:  "0x1111111111111111111111111111111111111111",
			ICS26ClientID: "client-0",
		}
		ctx, cancel := context.WithTimeout(context.Background(), configProbeTimeout)
		defer cancel()

		if findings := validateEVMDeployments(ctx, &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cfg,
		}); len(findings) != 0 {
			t.Fatalf("an unreachable endpoint must not fail startup, got %+v", findings)
		}
	})

	t.Run("covers the L2 return leg too", func(t *testing.T) {
		// The mirror of TestValidateEVMDeployments/reports every address that holds no code. An
		// l2_to_cosmos module reaches a router and a client id just like the outbound
		// legs; it only names them differently, and checking one direction while
		// leaving the other exposed is the asymmetry this repo keeps rediscovering.
		node := rpcmock.NewEVMNode(t)
		router := "0x7777777777777777777777777777777777777777"
		cfg := &appConfig{
			HasL2ToCosmos: true,
			L2ToCosmosConfig: l2ToCosmosConfig{
				L2RpcUrl:        node.URL(),
				L2ICS26ClientID: "l2-client-0",
				RollupProfile: json.RawMessage(
					`{"common":{"l2_chain_id":1,"l2_router":"` + router + `"}}`),
			},
		}

		// Undeployed router AND a reverting getClient. Both facts have the same
		// cause, so only the router may be reported — and the dedupe that ensures
		// it has to match on the L2 leg's own key name, not the outbound one.
		node.FailCall(jsonRPCError{msg: "execution reverted", data: "0x"})
		findings := validateEVMDeployments(context.Background(), cfg)
		if len(findings) != 1 || findings[0].key != "l2_to_cosmos.rollup_profile.common.l2_router" {
			t.Fatalf("expected only the undeployed L2 router to be reported, got %+v", findings)
		}
		node.FailCall(nil)

		// Deployed: the router now resolves and there is nothing to report.
		node.SetCode(common.HexToAddress(router), deployedCode)
		node.SetCallResult("0x000000000000000000000000" + strings.Repeat("ab", 20))
		if findings := validateEVMDeployments(context.Background(), cfg); len(findings) != 0 {
			t.Fatalf("expected no findings once the L2 router is deployed, got %+v", findings)
		}

		// And an unregistered client id on that router is reported against the L2
		// key, not the outbound one.
		node.FailCall(jsonRPCError{msg: "execution reverted", data: "0x"})
		findings = validateEVMDeployments(context.Background(), cfg)
		if len(findings) != 1 || findings[0].key != "l2_to_cosmos.l2_ics26_client_id" {
			t.Fatalf("expected an l2_ics26_client_id finding, got %+v", findings)
		}
	})

	// Reported by @DongLieu on #467. getClient succeeding proves only that the
	// router knows SOME client under that id. It says nothing about whether that
	// client is the spectre_client this config deploys against, nor whether its
	// counterparty is the Cosmos client this config relays to. A wrong-but-
	// registered id passes startup, and then every update and proof targets one
	// client while the router verifies another.
	//
	// routerWiringIsReusable already decides this and is already tested; only the
	// RPC to feed it was missing here. Per-selector call results are what let the
	// stub answer getClient and getCounterparty differently, which SetCallResult
	// alone could not express.
	t.Run("rejects a router wired to another client", func(t *testing.T) {
		const (
			getClientSelector       = "7eb78932"
			getCounterpartySelector = "b0777bfa"
		)
		abiAddress := func(hexAddr string) string {
			return "0x000000000000000000000000" + strings.TrimPrefix(hexAddr, "0x")
		}

		node := rpcmock.NewEVMNode(t)
		cfg := evmConfigFor(node)
		deployAll(node, cfg)
		// The counterparty the router reports, ABI-encoded as (string,bytes[]): a
		// dynamic head, then the Cosmos client id.
		node.SetCallResultFor(getCounterpartySelector, abiCounterparty("08-wasm-0"))

		t.Run("the client id points at a different light client", func(t *testing.T) {
			node.SetCallResultFor(getClientSelector, abiAddress("0x9999999999999999999999999999999999999999"))
			local := cfg
			local.CosmosWasmClientID = "08-wasm-0"
			findings := validateEVMDeployments(context.Background(), &appConfig{HasCosmosToEth: true, CosmosToEthConfig: local})
			if len(findings) == 0 {
				t.Fatal("a router pointing that client id at another light client was accepted; " +
					"updates would target spectre_client while the router verifies something else")
			}
		})

		t.Run("the counterparty is a different Cosmos client", func(t *testing.T) {
			node.SetCallResultFor(getClientSelector, abiAddress(cfg.SpectreClient))
			local := cfg
			local.CosmosWasmClientID = "08-wasm-7" // the router says 08-wasm-0
			findings := validateEVMDeployments(context.Background(), &appConfig{HasCosmosToEth: true, CosmosToEthConfig: local})
			if len(findings) == 0 {
				t.Fatal("a router counterparty-wired to another Cosmos client was accepted; " +
					"recvPacket would reject every packet's counterparty")
			}
		})

		t.Run("correctly wired passes", func(t *testing.T) {
			node.SetCallResultFor(getClientSelector, abiAddress(cfg.SpectreClient))
			local := cfg
			local.CosmosWasmClientID = "08-wasm-0"
			if findings := validateEVMDeployments(context.Background(), &appConfig{HasCosmosToEth: true, CosmosToEthConfig: local}); len(findings) != 0 {
				t.Fatalf("a correctly wired router was rejected: %v", findings)
			}
		})
	})

	// Reported by @DongLieu on #467. build_source.go:102 reads ICS26_CLIENT_ID over
	// the config, so validating the raw config value proves something startup then
	// does not use: a stale id in .env passes here and fails at the first packet,
	// which is exactly the failure this validation exists to move forward.
	t.Run("validates the ID startup will actually use", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		cfg := evmConfigFor(node)
		deployAll(node, cfg)
		// The router knows only the id in the config file: an id it does not know
		// reverts, which is how getClient reports IBCClientNotFound.
		node.FailCall(errors.New("execution reverted"))

		t.Setenv("ICS26_CLIENT_ID", "stale-from-dotenv")
		findings := validateEVMDeployments(context.Background(), &appConfig{HasCosmosToEth: true, CosmosToEthConfig: cfg})
		if len(findings) == 0 {
			t.Fatal("the override was ignored and the config value validated instead; a stale " +
				"ICS26_CLIENT_ID in .env would pass startup and fail at the first packet")
		}
		for _, f := range findings {
			if strings.Contains(f.detail, "stale-from-dotenv") {
				return
			}
		}
		t.Fatalf("no finding names the effective client id: %v", findings)
	})
}

func TestEvmDeploymentsInConfig_SkipUnsetAndMalformedAddresses(t *testing.T) {
	deployments := evmDeploymentsInConfig(&appConfig{
		HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{
			EthRpcUrl:     "http://eth",
			ICS26Address:  "0x1111111111111111111111111111111111111111",
			SpectreClient: "",            // unset
			Membership:    "not-an-addr", // malformed
		},
		HasCosmosToL2: true, CosmosToL2Config: cosmosToEthConfig{
			EthRpcUrl:    "http://l2",
			ICS26Address: "0x2222222222222222222222222222222222222222",
		},
	})

	if len(deployments) != 2 {
		t.Fatalf("expected one deployment per EVM endpoint, got %d", len(deployments))
	}
	if len(deployments[0].addresses) != 1 {
		t.Fatalf("unset/malformed addresses must not be probed, got %v", deployments[0].addresses)
	}
	if deployments[0].label != "cosmos_to_eth" || deployments[1].label != "cosmos_to_l2" {
		t.Fatalf("labels are %q and %q; an operator greps the config by these",
			deployments[0].label, deployments[1].label)
	}
}

func TestEvmDeploymentsInConfig_SkipAModuleWithNoRPC(t *testing.T) {
	if d := evmDeploymentsInConfig(&appConfig{
		HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{ICS26Address: "0x1111111111111111111111111111111111111111"},
	}); len(d) != 0 {
		t.Fatalf("a module with no eth_rpc_url has nothing to probe, got %+v", d)
	}
}

// --- revert classification ---------------------------------------------------

func TestIsEVMRevert_SeparatesExecutionFromTransport(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"geth revert text", errors.New("execution reverted"), true},
		{"revert text with a reason", errors.New("execution reverted: IBCClientNotFound"), true},
		{"mixed case", errors.New("Execution Reverted"), true},
		{"wrapped revert", fmt.Errorf("call getClient: %w", errors.New("execution reverted")), true},
		{"dial failure", errors.New("dial tcp 127.0.0.1:8545: connect: connection refused"), false},
		{"timeout", context.DeadlineExceeded, false},
		{"cancelled", context.Canceled, false},

		// The rpc.DataError branch. go-ethereum wraps EVERY JSON-RPC error in a
		// type satisfying it, so the interface alone cannot separate a revert
		// from a node complaint — only the presence of data can.
		{"revert with data", jsonRPCError{msg: "execution reverted", data: "0x08c379a0"}, true},
		{"pruned node", jsonRPCError{msg: "header not found"}, false},
		{"rate limited", jsonRPCError{msg: "your app has exceeded its compute units"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isEVMRevert(tc.err); got != tc.want {
				t.Fatalf("isEVMRevert(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// The path is what the log prefix carries, so it has to be derivable from config
// alone and stable for the life of the process. A1 makes that true by refusing a
// config with more than one path.
func TestRelayPathID(t *testing.T) {
	for _, tc := range []struct {
		name string
		cfg  *appConfig
		want string
	}{
		{
			name: "cosmos to eth carries the source client id",
			cfg:  &appConfig{HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{ICS26ClientID: "cosmoshub-1"}},
			want: "cosmos<->eth/cosmoshub-1",
		},
		{
			// The rollup names itself: an operator running an OP and an Arbitrum
			// path needs to tell the two apart, and "l2" does not.
			name: "cosmos to l2 names the rollup from the return leg",
			cfg: &appConfig{
				HasCosmosToL2: true, CosmosToL2Config: cosmosToEthConfig{ICS26ClientID: "cosmoshub-1"},
				HasL2ToCosmos: true, L2ToCosmosConfig: l2ToCosmosConfig{AttestorSrcChain: "arbitrum"},
			},
			want: "cosmos<->arbitrum/cosmoshub-1",
		},
		{
			name: "cosmos to l2 without a return leg still identifies the path",
			cfg:  &appConfig{HasCosmosToL2: true, CosmosToL2Config: cosmosToEthConfig{ICS26ClientID: "cosmoshub-1"}},
			want: "cosmos<->l2/cosmoshub-1",
		},
		{
			name: "return leg only",
			cfg:  &appConfig{HasL2ToCosmos: true, L2ToCosmosConfig: l2ToCosmosConfig{AttestorSrcChain: "opstack"}},
			want: "opstack->cosmos",
		},
		{
			// Never empty: an empty prefix would silently drop the identity from
			// every line rather than being obviously wrong.
			name: "nothing configured still yields a label",
			cfg:  &appConfig{},
			want: "relayer",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := relayPathID(tc.cfg); got != tc.want {
				t.Fatalf("relayPathID = %q, want %q", got, tc.want)
			}
		})
	}

	// Two processes relaying cosmos->eth from different source clients are a valid
	// deployment, and their lines are merged by journald or Loki with no file
	// boundary left. The prefix is the only thing that tells them apart, so it must
	// actually differ.
	t.Run("distinguishes two cosmos sources", func(t *testing.T) {
		a := relayPathID(&appConfig{HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{ICS26ClientID: "cosmoshub-1"}})
		b := relayPathID(&appConfig{HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{ICS26ClientID: "osmosis-1"}})
		if a == b {
			t.Fatalf("both sources produced %q; merged logs cannot be told apart", a)
		}
	})
}

// --- aggregation -------------------------------------------------------------

func TestValidateStartupConfig(t *testing.T) {
	t.Run("reports every problem in a single message", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("RELAYER_RECOVERY_STATE_FILE", filepath.Join(dir, "recovery-cursors.json"))
		t.Setenv("COSMOS_CHAIN_ID", "declared-chain")

		cosmos := rpcmock.NewCosmosNode(t)
		cosmos.Handle("status", func(json.RawMessage) (any, error) {
			return statusResult("served-chain"), nil
		})
		eth := rpcmock.NewEVMNode(t)
		cfg := evmConfigFor(eth)
		cfg.TmRpcUrl = cosmos.URL()

		err := validateStartupConfig(context.Background(), &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cfg,
		}, "")
		if err == nil {
			t.Fatal("expected config validation to fail")
		}

		// One run must diagnose the whole config. Failing on the first problem turns
		// a three-key mistake into three restarts.
		msg := err.Error()
		for _, want := range []string{
			"COSMOS_CHAIN_ID", "served-chain",
			"cosmos_to_eth.ics26_address", "cosmos_to_eth.update_client",
		} {
			if !strings.Contains(msg, want) {
				t.Fatalf("message does not mention %q:\n%s", want, msg)
			}
		}
	})

	t.Run("passes on a config with nothing to check", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("RELAYER_RECOVERY_STATE_FILE", filepath.Join(dir, "recovery-cursors.json"))

		if err := validateStartupConfig(context.Background(), &appConfig{}, ""); err != nil {
			t.Fatalf("an empty config has nothing to disagree with: %v", err)
		}
	})

	// Over a real connection: the finding reaches validateStartupConfig's single
	// message, which is where an operator reads it.
	t.Run("reports an attestor head mismatch", func(t *testing.T) {
		addr := serveInfo(t, &attestorpb.ChainInfo{SrcChain: "opstack", AttestationHead: "finalized"})
		captureStartupLog(t)

		err := validateStartupConfig(context.Background(), l2ConfigWithAttestor("opstack", "unsafe", addr), "")
		if err == nil {
			t.Fatal("expected config validation to fail")
		}
		for _, want := range []string{"attestor_endpoints[0].address / head_kind", `"unsafe"`, `"finalized"`, addr} {
			if !strings.Contains(err.Error(), want) {
				t.Fatalf("message does not mention %q:\n%s", want, err)
			}
		}
	})
}

func TestSortedKeys_IsStable(t *testing.T) {
	got := sortedKeys(map[string]string{"c": "", "a": "", "b": ""})
	if fmt.Sprint(got) != fmt.Sprint([]string{"a", "b", "c"}) {
		t.Fatalf("sortedKeys returned %v", got)
	}
}

// jsonRPCError mirrors go-ethereum's jsonError: it satisfies rpc.DataError for
// every node-returned error, carrying data only when the call actually reverted
// with a reason.
type jsonRPCError struct {
	msg  string
	data any
}

func (e jsonRPCError) Error() string { return e.msg }

func (e jsonRPCError) ErrorData() any { return e.data }

func (e jsonRPCError) ErrorCode() int { return -32000 }

// statusResult is the minimal /status reply: validation reads only the
// network name, so the stub models only that.
func statusResult(network string) *coretypes.ResultStatus {
	return &coretypes.ResultStatus{
		NodeInfo: cmtp2p.DefaultNodeInfo{Network: network},
	}
}

// abiCounterparty encodes what getCounterparty returns: a single dynamic STRUCT
// {string clientId; bytes[] merklePrefix}. A dynamic struct return is itself
// referenced by an offset, so the head is one word pointing at the tuple, and
// only then do the tuple's own two offsets begin. Getting that outer word wrong
// is why the first version of this test could not decode.
func abiCounterparty(clientID string) string {
	word := func(v int) string { return fmt.Sprintf("%064x", v) }
	padded := hex.EncodeToString([]byte(clientID))
	for len(padded)%64 != 0 {
		padded += "0"
	}
	tuple := word(0x40) + // clientId offset, relative to the tuple
		word(0x60+len(padded)/2) + // merklePrefix offset, past the string
		word(len(clientID)) + padded +
		word(0) // empty bytes[]
	return "0x" + word(0x20) + tuple
}

// pathExamples maps each runnable per-path example to the module names it must
// carry. config.example.json is the catalogue of every module across every path;
// these are the files an operator actually copies to config.json.
var pathExamples = map[string][]string{
	"config.ethereum.example.json": {"cosmos-to-eth", "eth-to-cosmos"},
	"config.op.example.json":       {"cosmos-to-op", "op-to-cosmos"},
	"config.arbitrum.example.json": {"cosmos-to-arbitrum", "arbitrum-to-cosmos"},
	"config.base.example.json":     {"cosmos-to-base", "base-to-cosmos"},
}

func readModules(t *testing.T, path string) map[string]json.RawMessage {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", path))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var parsed struct {
		Modules []struct {
			Name string `json:"name"`
		} `json:"modules"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	// Re-unmarshal into raw messages so comparison is on the full module object,
	// not just the fields this test happens to name.
	var full struct {
		Modules []json.RawMessage `json:"modules"`
	}
	if err := json.Unmarshal(raw, &full); err != nil {
		t.Fatalf("parse %s modules: %v", path, err)
	}
	out := make(map[string]json.RawMessage, len(full.Modules))
	for i, m := range full.Modules {
		out[parsed.Modules[i].Name] = m
	}
	return out
}

// TestAssertConfigMemberWritable_AgreesWithTheWrite pins the property that makes the
// precheck trustworthy: it accepts exactly the configs the real write accepts. A
// re-implemented lookup could drift and let the spend through again.
func TestAssertConfigMemberWritable_AgreesWithTheWrite(t *testing.T) {
	for _, tc := range []struct {
		name, body, sourceID string
		wantWritable         bool
	}{
		{"no cosmos_to_eth module", ethToCosmosPlusL2Config, "arb-client-0", false},
		{"unknown source id", validCosmosToEthConfig, "no-such-client", false},
		{"matching cosmos_to_eth module", validCosmosToEthConfig, "chain-a", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeConfigFixture(t, tc.body)
			checkErr := assertConfigMemberWritable(path, tc.sourceID, "cosmos_wasm_client_id", dirCosmosToEth)
			writeErr := writeEthWasmClientID(path, tc.sourceID, "08-wasm-0")
			if (checkErr == nil) != (writeErr == nil) {
				t.Fatalf("precheck and write disagree: check=%v write=%v", checkErr, writeErr)
			}
			if (checkErr == nil) != tc.wantWritable {
				t.Fatalf("writable = %v (%v), want %v", checkErr == nil, checkErr, tc.wantWritable)
			}
		})
	}
}

// TestUnrecordedClientError_NamesTheClientAndTheRecovery covers the case the precheck
// cannot prevent — the destination existed but the write still failed (permissions, a
// concurrent edit). The id only exists in this message, so it has to carry both the id
// and the one recovery that does not create a second client.
func TestUnrecordedClientError_NamesTheClientAndTheRecovery(t *testing.T) {
	cause := errors.New("permission denied")
	err := unrecordedClientError("/etc/relayer/config.json", "08-wasm-7", cause)

	for _, want := range []string{
		"08-wasm-7",                // the id, which exists nowhere else
		"/etc/relayer/config.json", // where it should have gone
		"cosmos_wasm_client_id",    // the field to set by hand
		"create-clients-eth",       // the command that resumes the flow
		"Do NOT re-run",            // the reaction that would orphan it
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error message missing %q:\n%s", want, err)
		}
	}
	if !errors.Is(err, cause) {
		t.Error("wrapped cause is not retrievable with errors.Is")
	}
}

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
