package main

import (
	"encoding/json"
	"strings"
	"testing"
)

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
			c := &l2ClientConfig{RollupProfile: tc.profile}
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

func TestInjectL1ClientID(t *testing.T) {
	t.Run("sets client_id and preserves everything else", func(t *testing.T) {
		in := json.RawMessage(`{
			"common": {
				"l1_chain_id": 1,
				"ethereum_client": {"client_id": "", "wasm_checksum": "0xabc"},
				"l2_router": "0x1111111111111111111111111111111111111111"
			},
			"dispute_game_factory": "0x2222222222222222222222222222222222222222"
		}`)
		out, err := injectL1ClientID(in, "08-wasm-7")
		if err != nil {
			t.Fatalf("injectL1ClientID: %v", err)
		}

		var got map[string]any
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatalf("unmarshal result: %v", err)
		}
		common := got["common"].(map[string]any)
		eth := common["ethereum_client"].(map[string]any)
		if eth["client_id"] != "08-wasm-7" {
			t.Fatalf("client_id = %v, want 08-wasm-7", eth["client_id"])
		}
		// Sibling and parent fields must be preserved verbatim.
		if eth["wasm_checksum"] != "0xabc" {
			t.Fatalf("wasm_checksum lost: %v", eth["wasm_checksum"])
		}
		if common["l2_router"] != "0x1111111111111111111111111111111111111111" {
			t.Fatalf("l2_router lost: %v", common["l2_router"])
		}
		if common["l1_chain_id"].(float64) != 1 {
			t.Fatalf("l1_chain_id lost: %v", common["l1_chain_id"])
		}
		if got["dispute_game_factory"] != "0x2222222222222222222222222222222222222222" {
			t.Fatalf("top-level rollup field lost: %v", got["dispute_game_factory"])
		}
	})

	for _, tc := range []struct {
		name    string
		profile string
		wantErr string
	}{
		{"missing common", `{"foo":1}`, "common is missing"},
		{"missing ethereum_client", `{"common":{"l2_router":"0x1"}}`, "ethereum_client is missing"},
		{"bad json", `{`, "parse rollup_profile"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := injectL1ClientID(json.RawMessage(tc.profile), "08-wasm-0")
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("err = %v, want substring %q", err, tc.wantErr)
			}
		})
	}
}
