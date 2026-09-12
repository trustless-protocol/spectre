package main

import (
	"attestor/types/attestation"
	"encoding/json"
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
