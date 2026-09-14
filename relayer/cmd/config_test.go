package main

import (
	"path/filepath"
	"strings"
	"testing"
)

// One relayer process serves exactly one source → destination pair (03-Relayer
// §2.1 requirement 5). Before this guard, `start` looped over every configured
// module and ran them in one process, which is what made a shared nonce cache and
// a shared loopErrCh necessary in the first place.
func TestValidateSingleRelayPair(t *testing.T) {
	c2e := func(ids ...string) []cosmosToEthConfig {
		out := make([]cosmosToEthConfig, 0, len(ids))
		for _, id := range ids {
			out = append(out, cosmosToEthConfig{ICS26ClientID: id})
		}
		return out
	}
	l2src := func(chains ...string) []l2ToCosmosConfig {
		out := make([]l2ToCosmosConfig, 0, len(chains))
		for _, c := range chains {
			out = append(out, l2ToCosmosConfig{AttestorSrcChain: c})
		}
		return out
	}

	tests := []struct {
		name string
		cfg  appConfig
		// wantNames are the identities the error must name so the operator knows
		// which paths to split apart; empty means the config must be accepted.
		wantNames []string
	}{
		{
			name: "one cosmos_to_eth",
			cfg:  appConfig{CosmosToEthConfigs: c2e("chain-a")},
		},
		{
			name: "one cosmos_to_l2 with its return leg is one pair",
			cfg: appConfig{
				CosmosToL2Configs: c2e("08-wasm-3"),
				L2ToCosmosConfigs: l2src("opstack"),
			},
		},
		{
			name: "forward-only cosmos_to_l2",
			cfg:  appConfig{CosmosToL2Configs: c2e("08-wasm-3")},
		},
		{
			name: "two cosmos_to_eth sources",
			cfg:  appConfig{CosmosToEthConfigs: c2e("chain-a", "chain-b")},
			// Naming both is the point: "config has 2 sources" leaves the operator
			// guessing which two, and file order is not what they edited by.
			wantNames: []string{"chain-a", "chain-b"},
		},
		{
			name:      "two cosmos_to_l2 destinations",
			cfg:       appConfig{CosmosToL2Configs: c2e("08-wasm-3", "08-wasm-4")},
			wantNames: []string{"08-wasm-3", "08-wasm-4"},
		},
		{
			name:      "two l2_to_cosmos sources",
			cfg:       appConfig{L2ToCosmosConfigs: l2src("opstack", "arbitrum")},
			wantNames: []string{"opstack", "arbitrum"},
		},
		{
			name: "cosmos_to_eth alongside cosmos_to_l2 is two pairs",
			cfg: appConfig{
				CosmosToEthConfigs: c2e("chain-a"),
				CosmosToL2Configs:  c2e("08-wasm-3"),
			},
			wantNames: []string{"chain-a", "08-wasm-3"},
		},
		{
			name: "cosmos_to_eth alongside an l2_to_cosmos return leg is two pairs",
			cfg: appConfig{
				CosmosToEthConfigs: c2e("chain-a"),
				L2ToCosmosConfigs:  l2src("opstack"),
			},
			wantNames: []string{"chain-a", "opstack"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSingleRelayPair(&tt.cfg, "config.json")
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
func TestValidateSingleRelayPairRejectsEmptyConfig(t *testing.T) {
	err := validateSingleRelayPair(&appConfig{}, "/etc/relayer/config.json")
	if err == nil {
		t.Fatal("start accepted a config with no relay module")
	}
	if !strings.Contains(err.Error(), "/etc/relayer/config.json") {
		t.Errorf("error must name the config file it read; got: %v", err)
	}
}

// The shipped examples encode the split the guard now enforces: one runnable file
// per path, plus a catalogue that carries every module and is not meant to be
// started. Without this test the two drift silently — the catalogue is the file
// whose name most invites `--config config.example.json`.
func TestShippedExamplesCarryOnePathEach(t *testing.T) {
	t.Parallel()

	for file := range pathExamples {
		t.Run(file, func(t *testing.T) {
			t.Parallel()
			cfg, err := loadConfig(filepath.Join("..", file))
			if err != nil {
				t.Fatalf("loadConfig(%s): %v", file, err)
			}
			if err := validateSingleRelayPair(cfg, file); err != nil {
				t.Fatalf("%s is a runnable example but start rejects it: %v", file, err)
			}
		})
	}

	t.Run("catalogue is not runnable", func(t *testing.T) {
		t.Parallel()
		cfg, err := loadConfig(filepath.Join("..", "config.example.json"))
		if err != nil {
			t.Fatalf("loadConfig(config.example.json): %v", err)
		}
		err = validateSingleRelayPair(cfg, "config.example.json")
		if err == nil {
			t.Fatal("the catalogue lists every path; start must send the operator to a per-path file instead")
		}
		if !strings.Contains(err.Error(), "--config") {
			t.Errorf("the rejection must show the way out; got: %v", err)
		}
	})
}
