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
