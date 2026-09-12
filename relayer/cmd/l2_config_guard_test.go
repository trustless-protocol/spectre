package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadL2ClientConfigRejectsRelayerConfig covers the mistake of passing
// relayer/config.json to --l2-config.
//
// The two files are both "the config" in conversation, both live in relayer/, and
// the flag takes a path — so the mix-up is easy. Without this guard the failure is
// "l2-config: wasm_checksum is required", which sends the operator looking for a
// missing field in a file that was never the right one.
func TestLoadL2ClientConfigRejectsRelayerConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	relayerConfig := `{
		"modules": [
			{"name": "cosmos_to_eth", "src_chain": "test", "config": {"eth_rpc_url": "http://localhost:8545"}}
		]
	}`
	if err := os.WriteFile(path, []byte(relayerConfig), configFilePerm); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, err := loadL2ClientConfig(path)
	if err == nil {
		t.Fatal("expected the relayer config to be rejected as an l2-config")
	}
	if !strings.Contains(err.Error(), "modules") {
		t.Fatalf("error should say WHY the file is wrong, got: %v", err)
	}
	if !strings.Contains(err.Error(), "example.json") {
		t.Fatalf("error should point at a template, got: %v", err)
	}
}

// A real l2-config must still load — the guard keys off "modules", which a valid
// l2-config never has.
func TestLoadL2ClientConfigAcceptsRealL2Config(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arb.json")
	l2 := `{
		"wasm_checksum": "dab7322721ae62e6b582c373ead06168f1281b696d67651f390922892215ce51",
		"l2_rpc_url": "http://localhost:8547",
		"attestors": {"public_keys": ["gojj3XQJ8ZX9UtstPLpdcspnCb8dlBIb83SIAbQPb1w="], "threshold": 1},
		"rollup_profile": {"common": {"l2_router": "0x9fcf7d13d10dedf17d0f24c62f0cf4ed462f65b7"}}
	}`
	if err := os.WriteFile(path, []byte(l2), configFilePerm); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, err := loadL2ClientConfig(path); err != nil {
		t.Fatalf("a valid l2-config must load, got: %v", err)
	}
}
