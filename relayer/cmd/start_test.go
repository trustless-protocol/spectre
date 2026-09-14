package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// Reported by @DongLieu on #435. stampRelayPathOnLogs and relayPathID were both
// reachable only from their unit tests: the production call was written in
// 74cc0147, survived the first merge, and was lost in 3b75da08 when the base
// branch moved Start out of main.go into start.go and the resolution took the
// base's file. Every relay log kept the default empty prefix, so the entire
// point of the PR was absent while its helper tests stayed green.
//
// This drives the real cobra command instead of the helper, which is the only
// shape that can catch it. Start is expected to FAIL here -- the prover bin
// directory does not exist -- and that is fine: the stamp happens after config
// validation and before the prover load, so the prefix must already be set by
// the time the error comes back.
func TestStartStampsTheRelayPathBeforeItCanFail(t *testing.T) {
	origFlags, origPrefix, origOut := log.Flags(), log.Prefix(), log.Writer()
	t.Cleanup(func() {
		log.SetFlags(origFlags)
		log.SetPrefix(origPrefix)
		log.SetOutput(origOut)
	})
	log.SetOutput(&bytes.Buffer{})
	log.SetPrefix("")

	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	config := `{"modules":[{"name":"cosmos_to_eth","src_chain":"cosmoshub-1","config":{` +
		`"tm_rpc_url":"http://127.0.0.1:26657","eth_rpc_url":"http://127.0.0.1:8545",` +
		`"ics26_address":"0x80741a37e3644612f0465145c9709a90b6d77ee3","ics26_client_id":"cosmoshub-1"}}]}`
	if err := os.WriteFile(configPath, []byte(config), configFilePerm); err != nil {
		t.Fatalf("write config: %v", err)
	}
	// Valid-shaped keys, because startup rejects missing ones BEFORE it loads the
	// config, and this test is about what happens after that.
	const testKey = "a81f9eb900c02f35d28ec80d1da79dde52cc4cc37d6fbaef5768ab59f32cbfdd"
	t.Setenv("ETH_PRIVATE_KEY", testKey)
	t.Setenv("COSMOS_PRIVATE_KEY", testKey)
	// Point the prover somewhere that cannot load, so the command stops right
	// after the stamp rather than trying to reach a chain.
	t.Setenv("PROVER_BIN_DIR", filepath.Join(dir, "no-such-bin"))

	core, logs := observer.New(zapcore.InfoLevel)
	command := Start(zap.New(core))
	// --gpu-prove=false selects the CPU backend, which resolves fine and makes
	// Start log one line through zap between the stamp and the prover failure --
	// otherwise the zap half has nothing observable to assert on here.
	command.SetArgs([]string{"--config", configPath, "--gpu-prove=false"})
	command.SetOut(&bytes.Buffer{})
	command.SetErr(&bytes.Buffer{})
	err := command.Execute()

	if got := log.Prefix(); !strings.Contains(got, "cosmoshub-1") {
		t.Fatalf("log prefix = %q, want it to carry the relay path; "+
			"start never called stampRelayPath, so every relay line is unlabelled "+
			"(command returned %v)", got, err)
	}

	// The zap half, which the first version of this test could not see: Start and
	// the builders emit their lifecycle lines through the injected logger, and a
	// stamp that only touches log.Default() leaves every one of them
	// unattributable in a merged log.
	entries := logs.All()
	if len(entries) == 0 {
		t.Fatalf("start emitted no zap line before failing (%v), so the scoped logger "+
			"is unverified here; the test needs a line between the stamp and the failure", err)
	}
	for _, e := range entries {
		got, ok := e.ContextMap()[relayPathLogField]
		if !ok {
			t.Fatalf("zap line %q carries no %q field: start used the unscoped logger "+
				"(command returned %v)", e.Message, relayPathLogField, err)
		}
		if !strings.Contains(fmt.Sprint(got), "cosmoshub-1") {
			t.Fatalf("zap line %q has %s = %v, want the relay path", e.Message, relayPathLogField, got)
		}
	}
}
