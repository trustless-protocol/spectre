package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// pathExamples maps each runnable per-path example to the module names it must
// carry. config.example.json is the catalogue of every module across every path;
// these are the files an operator actually copies to config.json.
var pathExamples = map[string][]string{
	"config.ethereum.example.json":  {"cosmos-to-eth", "eth-to-cosmos"},
	"config.op.example.json":        {"cosmos-to-op", "op-to-cosmos"},
	"config.arbitrum.example.json":  {"cosmos-to-arbitrum", "arbitrum-to-cosmos"},
	"config.base.example.json":      {"cosmos-to-base", "base-to-cosmos"},
	"config.avalanche.example.json": {"cosmos-to-avalanche", "avalanche-to-cosmos"},
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

// TestPathExamplesMatchCatalogue keeps the per-path examples byte-identical to
// the catalogue they are cut from. Without this they drift: the file these
// replaced (config.arbitrum-only.json) had gone stale on log_scan_chunk and
// include_provisional and carried the addresses of a long-finished devnet run,
// and nothing in the repo referenced it so nothing caught that.
func TestPathExamplesMatchCatalogue(t *testing.T) {
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
}

// TestExampleL2RouterMatchesPair pins the one cross-module invariant the example
// has to demonstrate: a return module's rollup_profile.common.l2_router is the
// SAME contract as its forward pair's ics26_address. deploy_l2_contracts.sh
// writes both from one deploy, so an example where they disagree teaches the
// mistake that used to break the return direction silently.
func TestExampleL2RouterMatchesPair(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("..", "config.example.json"))
	if err != nil {
		t.Fatalf("read catalogue: %v", err)
	}
	var parsed struct {
		Modules []struct {
			Name   string `json:"name"`
			Config struct {
				ICS26Address    string `json:"ics26_address"`
				ICS26ClientID   string `json:"ics26_client_id"`
				L2ICS26ClientID string `json:"l2_ics26_client_id"`
				RollupProfile   *struct {
					Common struct {
						L2Router string `json:"l2_router"`
					} `json:"common"`
				} `json:"rollup_profile"`
			} `json:"config"`
		} `json:"modules"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("parse catalogue: %v", err)
	}

	routerByClientID := map[string]string{}
	for _, m := range parsed.Modules {
		if m.Config.ICS26ClientID != "" && m.Config.ICS26Address != "" {
			routerByClientID[m.Config.ICS26ClientID] = m.Config.ICS26Address
		}
	}

	checked := 0
	for _, m := range parsed.Modules {
		if m.Config.RollupProfile == nil {
			continue
		}
		want, ok := routerByClientID[m.Config.L2ICS26ClientID]
		if !ok {
			t.Errorf("module %q has l2_ics26_client_id %q with no forward pair",
				m.Name, m.Config.L2ICS26ClientID)
			continue
		}
		if got := m.Config.RollupProfile.Common.L2Router; got != want {
			t.Errorf("module %q l2_router = %s, want its pair's ics26_address %s",
				m.Name, got, want)
		}
		checked++
	}
	if checked == 0 {
		t.Fatal("no l2_to_cosmos modules were checked; the catalogue lost its L2 pairs")
	}
}

// TestPathExamplesLoad is the point of these files: unlike the catalogue, each
// one is a config the relayer can actually start from.
func TestPathExamplesLoad(t *testing.T) {
	t.Parallel()

	for file := range pathExamples {
		t.Run(file, func(t *testing.T) {
			t.Parallel()
			if _, err := loadConfig(filepath.Join("..", file)); err != nil {
				t.Fatalf("loadConfig(%s) error = %v", file, err)
			}
		})
	}
}

func TestExampleRollupProfilesUseOnlyFreshV2Schema(t *testing.T) {
	t.Parallel()

	examples := map[string]int{
		"config.example.json":                  4,
		"config.op.example.json":               1,
		"config.base.example.json":             1,
		"config.arbitrum.example.json":         1,
		"config.avalanche.example.json":        1,
		"op-l2-config.example.json":            1,
		"base-l2-config.example.json":          1,
		"arb-l2-config.example.json":           1,
		"avalanche-client-config.example.json": 1,
	}
	wantKeys := map[string]struct{}{
		"l2_chain_id": {}, "l2_router": {}, "commitment_slot": {},
		"profile_version": {}, "l2_header_fork": {},
	}
	allowedVersions := map[string]struct{}{
		"op_attestor_v1": {}, "base_attestor_v1": {}, "arbitrum_attestor_v1": {},
		"avalanche_attestor_v1": {},
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
}
