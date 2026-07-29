package main

import (
	"encoding/json"
	"os"

	"testing"
)

// l2SourceConfigJSON mirrors the shape config.example.json ships: an l2_to_cosmos
// module carrying placeholder client ids next to a cosmos_to_l2 module that must
// stay untouched.
const l2SourceConfigJSON = `{
  "modules": [
    {
      "name": "cosmos-to-op",
      "src_chain": "cosmos",
      "dst_chain": "opstack",
      "config": {
        "cosmos_wasm_client_id": "08-wasm-0",
        "ics26_client_id": "cosmoshub-2"
      }
    },
    {
      "name": "op-to-cosmos",
      "src_chain": "opstack",
      "dst_chain": "cosmos",
      "config": {
        "l2_wasm_client_id": "08-wasm-1",
        "l2_ics26_client_id": "cosmoshub-2",
        "head_kind": "safe",
        "rollup_profile": {
          "common": {
            "l2_router": "0xrouter",
            "ethereum_client": {
              "client_id": "08-wasm-0"
            }
          },
          "dispute_game_factory": "0xfactory"
        }
      }
    }
  ]
}`

// parsed decodes the written config so assertions read fields rather than bytes.
type parsedConfig struct {
	Modules []struct {
		Name     string `json:"name"`
		SrcChain string `json:"src_chain"`
		DstChain string `json:"dst_chain"`
		Config   struct {
			CosmosWasmClientID string `json:"cosmos_wasm_client_id"`
			L2WasmClientID     string `json:"l2_wasm_client_id"`
			L2ICS26ClientID    string `json:"l2_ics26_client_id"`
			HeadKind           string `json:"head_kind"`
			RollupProfile      struct {
				Common struct {
					L2Router       string `json:"l2_router"`
					EthereumClient struct {
						ClientID string `json:"client_id"`
					} `json:"ethereum_client"`
				} `json:"common"`
				DisputeGameFactory string `json:"dispute_game_factory"`
			} `json:"rollup_profile"`
		} `json:"config"`
	} `json:"modules"`
}

// writeTempConfig lives in dispatch_test.go.

func readParsed(t *testing.T, path string) parsedConfig {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	var cfg parsedConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse written config: %v\n%s", err, data)
	}
	return cfg
}

// The regression this whole file exists for: a stale ethereum_client.client_id left
// in the l2_to_cosmos module makes the relayer pin headers to a different Ethereum
// client than the one the L2 wasm client actually queries.
func TestWriteL2SourceClientIDsOverwritesBothPlaceholders(t *testing.T) {
	path := writeTempConfig(t, l2SourceConfigJSON)

	if err := writeL2SourceClientIDs(path, "cosmoshub-2", "08-wasm-5", "08-wasm-4"); err != nil {
		t.Fatalf("writeL2SourceClientIDs: %v", err)
	}

	cfg := readParsed(t, path)
	l2 := cfg.Modules[1].Config
	if l2.L2WasmClientID != "08-wasm-5" {
		t.Errorf("l2_wasm_client_id = %q, want 08-wasm-5", l2.L2WasmClientID)
	}
	if got := l2.RollupProfile.Common.EthereumClient.ClientID; got != "08-wasm-4" {
		t.Errorf("ethereum_client.client_id = %q, want 08-wasm-4", got)
	}
}

func TestWriteL2SourceClientIDsPreservesSiblings(t *testing.T) {
	path := writeTempConfig(t, l2SourceConfigJSON)

	if err := writeL2SourceClientIDs(path, "cosmoshub-2", "08-wasm-5", "08-wasm-4"); err != nil {
		t.Fatalf("writeL2SourceClientIDs: %v", err)
	}

	cfg := readParsed(t, path)
	// The cosmos_to_l2 module shares the same client id; its Ethereum-side field
	// must not be swept up by the l2_to_cosmos write.
	if got := cfg.Modules[0].Config.CosmosWasmClientID; got != "08-wasm-0" {
		t.Errorf("cosmos_to_l2 cosmos_wasm_client_id = %q, want it untouched (08-wasm-0)", got)
	}
	l2 := cfg.Modules[1].Config
	if l2.HeadKind != "safe" {
		t.Errorf("head_kind = %q, want it untouched (safe)", l2.HeadKind)
	}
	if l2.RollupProfile.Common.L2Router != "0xrouter" {
		t.Errorf("l2_router = %q, want it untouched", l2.RollupProfile.Common.L2Router)
	}
	if l2.RollupProfile.DisputeGameFactory != "0xfactory" {
		t.Errorf("dispute_game_factory = %q, want it untouched", l2.RollupProfile.DisputeGameFactory)
	}
}

// An empty id means "nothing new to persist" — the existing value must survive
// rather than being blanked.
func TestWriteL2SourceClientIDsSkipsEmptyIDs(t *testing.T) {
	path := writeTempConfig(t, l2SourceConfigJSON)

	if err := writeL2SourceClientIDs(path, "cosmoshub-2", "", "08-wasm-4"); err != nil {
		t.Fatalf("writeL2SourceClientIDs: %v", err)
	}

	l2 := readParsed(t, path).Modules[1].Config
	if l2.L2WasmClientID != "08-wasm-1" {
		t.Errorf("l2_wasm_client_id = %q, want the pre-existing 08-wasm-1", l2.L2WasmClientID)
	}
	if got := l2.RollupProfile.Common.EthereumClient.ClientID; got != "08-wasm-4" {
		t.Errorf("ethereum_client.client_id = %q, want 08-wasm-4", got)
	}
}

// A profile that omits the nested objects entirely still has to end up anchored,
// otherwise the operator is back to hand-copying the id.
func TestWriteL2SourceClientIDsCreatesMissingNesting(t *testing.T) {
	path := writeTempConfig(t, `{
  "modules": [
    {
      "src_chain": "opstack",
      "dst_chain": "cosmos",
      "config": {
        "l2_ics26_client_id": "cosmoshub-2",
        "rollup_profile": {}
      }
    }
  ]
}`)

	if err := writeL2SourceClientIDs(path, "cosmoshub-2", "08-wasm-5", "08-wasm-4"); err != nil {
		t.Fatalf("writeL2SourceClientIDs: %v", err)
	}

	l2 := readParsed(t, path).Modules[0].Config
	if got := l2.RollupProfile.Common.EthereumClient.ClientID; got != "08-wasm-4" {
		t.Errorf("ethereum_client.client_id = %q, want 08-wasm-4", got)
	}
	if l2.L2WasmClientID != "08-wasm-5" {
		t.Errorf("l2_wasm_client_id = %q, want 08-wasm-5", l2.L2WasmClientID)
	}
}

// Several L2 sources share one config; the write must land on the addressed one.
func TestWriteL2SourceClientIDsTargetsMatchingSource(t *testing.T) {
	path := writeTempConfig(t, `{
  "modules": [
    {
      "src_chain": "opstack",
      "dst_chain": "cosmos",
      "config": {
        "l2_ics26_client_id": "cosmoshub-2",
        "rollup_profile": {"common": {"ethereum_client": {"client_id": "08-wasm-0"}}}
      }
    },
    {
      "src_chain": "arbitrum",
      "dst_chain": "cosmos",
      "config": {
        "l2_ics26_client_id": "cosmoshub-3",
        "rollup_profile": {"common": {"ethereum_client": {"client_id": "08-wasm-0"}}}
      }
    }
  ]
}`)

	if err := writeL2SourceClientIDs(path, "cosmoshub-3", "08-wasm-7", "08-wasm-4"); err != nil {
		t.Fatalf("writeL2SourceClientIDs: %v", err)
	}

	cfg := readParsed(t, path)
	if got := cfg.Modules[0].Config.RollupProfile.Common.EthereumClient.ClientID; got != "08-wasm-0" {
		t.Errorf("opstack ethereum_client.client_id = %q, want it untouched", got)
	}
	if got := cfg.Modules[1].Config.RollupProfile.Common.EthereumClient.ClientID; got != "08-wasm-4" {
		t.Errorf("arbitrum ethereum_client.client_id = %q, want 08-wasm-4", got)
	}
}

// A missing target must fail loudly: silently writing nothing is how the stale id
// survived in the first place.
func TestWriteL2SourceClientIDsFailsOnUnknownSource(t *testing.T) {
	path := writeTempConfig(t, l2SourceConfigJSON)

	err := writeL2SourceClientIDs(path, "cosmoshub-9", "08-wasm-5", "08-wasm-4")
	if err == nil {
		t.Fatal("writeL2SourceClientIDs succeeded for an unknown l2_ics26_client_id, want an error")
	}
}

// Classification must follow src_chain/dst_chain, not the free-form module name.
func TestWriteL2SourceClientIDsIgnoresModuleName(t *testing.T) {
	path := writeTempConfig(t, `{
  "modules": [
    {
      "name": "whatever-label",
      "src_chain": "opstack",
      "dst_chain": "cosmos",
      "config": {
        "l2_ics26_client_id": "cosmoshub-2",
        "rollup_profile": {"common": {"ethereum_client": {"client_id": "08-wasm-0"}}}
      }
    }
  ]
}`)

	if err := writeL2SourceClientIDs(path, "cosmoshub-2", "", "08-wasm-4"); err != nil {
		t.Fatalf("writeL2SourceClientIDs: %v", err)
	}
	got := readParsed(t, path).Modules[0].Config.RollupProfile.Common.EthereumClient.ClientID
	if got != "08-wasm-4" {
		t.Errorf("ethereum_client.client_id = %q, want 08-wasm-4", got)
	}
}

// Regression for a High found in review: a cosmos_to_eth and a cosmos_to_l2 module may
// legitimately share an ics26_client_id (config.example.json ships exactly that), and
// the direction-agnostic writer resolved the ambiguity by file order — writing the L2
// client id into the Ethereum module and leaving the L2 one untouched.
func TestWriteWasmClientIDRespectsModuleDirection(t *testing.T) {
	const shared = `{
  "modules": [
    {"src_chain": "cosmos", "dst_chain": "ethereum",
     "config": {"ics26_client_id": "cosmoshub-1", "cosmos_wasm_client_id": "eth-original"}},
    {"src_chain": "cosmos", "dst_chain": "opstack",
     "config": {"ics26_client_id": "cosmoshub-1", "cosmos_wasm_client_id": "l2-original"}}
  ]
}`

	t.Run("l2 write lands on the cosmos_to_l2 module", func(t *testing.T) {
		path := writeTempConfig(t, shared)
		if err := writeL2WasmClientID(path, "cosmoshub-1", "08-wasm-7"); err != nil {
			t.Fatalf("writeL2WasmClientID: %v", err)
		}
		cfg := readParsed(t, path)
		if got := cfg.Modules[0].Config.CosmosWasmClientID; got != "eth-original" {
			t.Errorf("cosmos_to_eth cosmos_wasm_client_id = %q, want it untouched", got)
		}
		if got := cfg.Modules[1].Config.CosmosWasmClientID; got != "08-wasm-7" {
			t.Errorf("cosmos_to_l2 cosmos_wasm_client_id = %q, want 08-wasm-7", got)
		}
	})

	t.Run("eth write lands on the cosmos_to_eth module", func(t *testing.T) {
		path := writeTempConfig(t, shared)
		if err := writeEthWasmClientID(path, "cosmoshub-1", "08-wasm-4"); err != nil {
			t.Fatalf("writeEthWasmClientID: %v", err)
		}
		cfg := readParsed(t, path)
		if got := cfg.Modules[0].Config.CosmosWasmClientID; got != "08-wasm-4" {
			t.Errorf("cosmos_to_eth cosmos_wasm_client_id = %q, want 08-wasm-4", got)
		}
		if got := cfg.Modules[1].Config.CosmosWasmClientID; got != "l2-original" {
			t.Errorf("cosmos_to_l2 cosmos_wasm_client_id = %q, want it untouched", got)
		}
	})
}
