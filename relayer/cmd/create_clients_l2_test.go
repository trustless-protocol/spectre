package main

import (
	"attestor/types/attestation"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validL2Attestors() attestation.AttestorConfig {
	return attestation.AttestorConfig{
		PublicKeys: [][]byte{{0x82, 0x88, 0xe3, 0xdd, 0x74, 0x09, 0xf1, 0x95, 0xfd, 0x52, 0xdb, 0x2d, 0x3c, 0xba, 0x5d, 0x72, 0xca, 0x67, 0x09, 0xbf, 0x1d, 0x94, 0x12, 0x1b, 0xf3, 0x74, 0x88, 0x01, 0xb4, 0x0f, 0x6f, 0x5c}},
		Threshold:  1,
	}
}

func TestRouterAddress(t *testing.T) {
	profile := func(router string) json.RawMessage {
		if router == "__omit__" {
			return json.RawMessage(`{"common":{}}`)
		}
		return json.RawMessage(`{"common":{"l2_router":"` + router + `"}}`)
	}

	cases := []struct {
		name    string
		profile json.RawMessage
		want    string
		wantErr string // substring; empty = expect success
	}{
		{"valid", profile("0x1111111111111111111111111111111111111111"), "0x1111111111111111111111111111111111111111", ""},
		{"missing", profile("__omit__"), "", "is required"},
		{"empty", profile(""), "", "is required"},
		{"malformed hex", profile("0xnothex"), "", "not a valid hex address"},
		{"wrong length", profile("0x1234"), "", "not a valid hex address"},
		{"bad json", json.RawMessage(`{"common": `), "", "parse rollup_profile"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &l2ClientConfig{RollupProfile: tc.profile, Attestors: validL2Attestors()}
			got, err := c.routerAddress()
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("routerAddress() err = %v, want substring %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("routerAddress() unexpected err: %v", err)
			}
			if got != tc.want {
				t.Fatalf("routerAddress() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestLoadL2ClientConfig_RejectsRelayerConfig covers the mistake of passing
// relayer/config.json to --l2-config.
//
// The two files are both "the config" in conversation, both live in relayer/, and
// the flag takes a path — so the mix-up is easy. Without this guard the failure is
// "l2-config: wasm_checksum is required", which sends the operator looking for a
// missing field in a file that was never the right one.
func TestLoadL2ClientConfig_RejectsRelayerConfig(t *testing.T) {
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
func TestLoadL2ClientConfig_AcceptsRealL2Config(t *testing.T) {
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
