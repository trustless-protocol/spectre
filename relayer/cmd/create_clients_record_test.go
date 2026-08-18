package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
)

// ethToCosmosPlusL2Config is the shape that makes #310 reachable: an eth_to_cosmos
// module supplies the beacon endpoint, an L2 module supplies the ics26_client_id
// selectSource lands in CosmosToEthConfig, and no cosmos_to_eth module exists for
// the write-back to target. Every guard except the write itself passes.
const ethToCosmosPlusL2Config = `{
	"modules": [
		{"name": "eth-to-cosmos", "src_chain": "ethereum", "dst_chain": "cosmos",
		 "config": {"eth_beacon_api_url": "http://127.0.0.1:5052"}},
		{"name": "cosmos-to-arbitrum", "src_chain": "cosmos", "dst_chain": "arbitrum",
		 "config": {"tm_rpc_url": "http://127.0.0.1:26657", "eth_rpc_url": "http://127.0.0.1:8547",
		            "eth_ws_url": "ws://127.0.0.1:8548",
		            "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3",
		            "ics26_client_id": "arb-client-0"}}
	]
}`

func writeConfigFixture(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(body), configFilePerm); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

// TestCreateClientsCosmosRefusesAConfigItCannotRecordInto is the #310 regression.
//
// The RPC endpoints point at closed ports on purpose: the precheck must run BEFORE
// buildCreateClientsDeps dials anything, so a failure here has to be the recording
// destination and not a connection. If the check is ever moved below the dial, this
// test reports a dial error instead and fails.
func TestCreateClientsCosmosRefusesAConfigItCannotRecordInto(t *testing.T) {
	configPath := writeConfigFixture(t, ethToCosmosPlusL2Config)
	cfg, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	sel, err := selectSource(cfg, "")
	if err != nil {
		t.Fatalf("selectSource: %v", err)
	}
	// Precondition of the regression: this config reaches the spend in every other way.
	if sel.CosmosToEthConfig.ICS26ClientID != "arb-client-0" {
		t.Fatalf("selectSource picked %q, want the L2 module — fixture no longer reproduces the shape",
			sel.CosmosToEthConfig.ICS26ClientID)
	}
	if resolveBeaconURL(sel) == "" {
		t.Fatal("beacon URL empty — the earlier guard would abort and this fixture proves nothing")
	}

	_, err = runCreateClientsCosmos(zap.NewNop(), sel, configPath, "deadbeef")
	if err == nil {
		t.Fatal("runCreateClientsCosmos succeeded; want a refusal before any client is created")
	}
	if !strings.Contains(err.Error(), "could not be recorded") {
		t.Fatalf("error = %v; want the pre-spend refusal (a dial error means the check runs too late)", err)
	}
	if !strings.Contains(err.Error(), "arb-client-0") {
		t.Fatalf("error = %v; want it to name the unresolvable source id", err)
	}
}

// TestAssertConfigMemberWritableAgreesWithTheWrite pins the property that makes the
// precheck trustworthy: it accepts exactly the configs the real write accepts. A
// re-implemented lookup could drift and let the spend through again.
func TestAssertConfigMemberWritableAgreesWithTheWrite(t *testing.T) {
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

const validCosmosToEthConfig = `{
	"modules": [
		{"name": "cosmos-to-eth", "src_chain": "cosmos", "dst_chain": "ethereum",
		 "config": {"tm_rpc_url": "http://127.0.0.1:26657", "eth_rpc_url": "http://127.0.0.1:8545",
		            "eth_ws_url": "ws://127.0.0.1:8546",
		            "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3",
		            "ics26_client_id": "chain-a"}},
		{"name": "eth-to-cosmos", "src_chain": "ethereum", "dst_chain": "cosmos",
		 "config": {"eth_beacon_api_url": "http://127.0.0.1:5052"}}
	]
}`

// TestUnrecordedClientErrorNamesTheClientAndTheRecovery covers the case the precheck
// cannot prevent — the destination existed but the write still failed (permissions, a
// concurrent edit). The id only exists in this message, so it has to carry both the id
// and the one recovery that does not create a second client.
func TestUnrecordedClientErrorNamesTheClientAndTheRecovery(t *testing.T) {
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
