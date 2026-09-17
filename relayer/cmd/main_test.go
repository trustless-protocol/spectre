package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"relayer/workers"
	workergroup "relayer/workers"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestValidateStartupKeys_RejectsInvalidCosmosAddressPrefix(t *testing.T) {
	const privateKey = "a81f9eb900c02f35d28ec80d1da79dde52cc4cc37d6fbaef5768ab59f32cbfdd"
	t.Setenv("ETH_PRIVATE_KEY", privateKey)
	t.Setenv("COSMOS_PRIVATE_KEY", privateKey)
	t.Setenv("COSMOS_ADDRESS_PREFIX", "invalid prefix")

	if err := validateStartupKeys(); err == nil {
		t.Fatal("startup accepted a Cosmos signer address with an invalid bech32 prefix")
	}
}

func TestReplaceConfigMember_ICS07Address(t *testing.T) {
	t.Parallel()

	const addr = "0x1234567890abcdef1234567890abcdef12345678"

	tests := []struct {
		name   string
		input  string
		assert func(t *testing.T, out []byte, cosmosToEth map[string]any)
	}{
		{
			name: "replaces existing spectre_client",
			input: `{
				"modules": [
					{"name": "cosmos_to_eth", "config": {"spectre_client": "0xold", "other": "keep"}},
					{"name": "eth_to_cosmos", "config": {"spectre_client": "0xother"}}
				]
			}`,
			assert: func(t *testing.T, out []byte, cosmosToEth map[string]any) {
				if got := cosmosToEth["other"]; got != "keep" {
					t.Fatalf("other field = %v, want keep", got)
				}
				ethToCosmos := moduleConfigByName(t, out, "eth_to_cosmos")
				if got := ethToCosmos["spectre_client"]; got != "0xother" {
					t.Fatalf("eth_to_cosmos spectre_client = %v, want 0xother", got)
				}
			},
		},
		{
			name: "inserts spectre_client when absent",
			input: `{
				"modules": [
					{"name": "cosmos_to_eth", "config": {"eth_rpc_url": "http://localhost", "nested": {"keep": true}}}
				]
			}`,
			assert: func(t *testing.T, _ []byte, cosmosToEth map[string]any) {
				if got := cosmosToEth["eth_rpc_url"]; got != "http://localhost" {
					t.Fatalf("eth_rpc_url = %v, want http://localhost", got)
				}
				nested, ok := cosmosToEth["nested"].(map[string]any)
				if !ok {
					t.Fatalf("nested = %T, want object", cosmosToEth["nested"])
				}
				if got := nested["keep"]; got != true {
					t.Fatalf("nested.keep = %v, want true", got)
				}
			},
		},
		{
			name: "preserves nested config objects",
			input: `{
				"modules": [
					{"name": "cosmos_to_eth", "config": {"nested": {"config": {"spectre_client": "nested"}, "items": [{"k": "v,}]"}]}, "spectre_client": "0xold"}}
				]
			}`,
			assert: func(t *testing.T, _ []byte, cosmosToEth map[string]any) {
				nested, ok := cosmosToEth["nested"].(map[string]any)
				if !ok {
					t.Fatalf("nested = %T, want object", cosmosToEth["nested"])
				}
				inner, ok := nested["config"].(map[string]any)
				if !ok {
					t.Fatalf("nested.config = %T, want object", nested["config"])
				}
				if got := inner["spectre_client"]; got != "nested" {
					t.Fatalf("nested.config.spectre_client = %v, want nested", got)
				}
			},
		},
		{
			name: "handles escaped strings and keys",
			input: `{
				"modules": [
					{"name": "cosmos\u005fto\u005feth", "config": {"note": "quote: \" and slash: \\", "ics07\u005fclient": "0xold"}}
				]
			}`,
			assert: func(t *testing.T, _ []byte, cosmosToEth map[string]any) {
				if got := cosmosToEth["note"]; got != `quote: " and slash: \` {
					t.Fatalf("note = %q, want escaped string preserved", got)
				}
			},
		},
		{
			name:  "handles unusual whitespace",
			input: "{\n\t\"modules\"\t:\t[\n\t { \n \"name\" : \"cosmos_to_eth\" ,\n \"config\" : { \n \"foo\" : [1, {\"bar\" : \"baz\"}] \n } \n }\n]\n}",
			assert: func(t *testing.T, _ []byte, cosmosToEth map[string]any) {
				foo, ok := cosmosToEth["foo"].([]any)
				if !ok {
					t.Fatalf("foo = %T, want array", cosmosToEth["foo"])
				}
				if len(foo) != 2 {
					t.Fatalf("foo length = %d, want 2", len(foo))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := replaceConfigMember([]byte(tc.input), "spectre_client", addr)
			if err != nil {
				t.Fatalf("replaceConfigMember() error = %v", err)
			}

			cosmosToEth := moduleConfigByName(t, out, "cosmos_to_eth")
			if got := cosmosToEth["spectre_client"]; got != addr {
				t.Fatalf("spectre_client = %v, want %s", got, addr)
			}
			tc.assert(t, out, cosmosToEth)
		})
	}
}

func TestReplaceConfigMember_WasmClientID(t *testing.T) {
	t.Parallel()

	const id = "08-wasm-445"

	t.Run("inserts cosmos_wasm_client_id when absent", func(t *testing.T) {
		t.Parallel()
		in := `{
			"modules": [
				{"name": "cosmos_to_eth", "config": {"spectre_client": "0xabc"}}
			]
		}`
		out, err := replaceConfigMember([]byte(in), "cosmos_wasm_client_id", id)
		if err != nil {
			t.Fatalf("replaceConfigMember() error = %v", err)
		}
		cfg := moduleConfigByName(t, out, "cosmos_to_eth")
		if got := cfg["cosmos_wasm_client_id"]; got != id {
			t.Fatalf("cosmos_wasm_client_id = %v, want %s", got, id)
		}
		if got := cfg["spectre_client"]; got != "0xabc" {
			t.Fatalf("spectre_client = %v, want preserved 0xabc", got)
		}
	})

	t.Run("replaces existing cosmos_wasm_client_id", func(t *testing.T) {
		t.Parallel()
		in := `{
			"modules": [
				{"name": "cosmos_to_eth", "config": {"cosmos_wasm_client_id": "08-wasm-0", "keep": "yes"}}
			]
		}`
		out, err := replaceConfigMember([]byte(in), "cosmos_wasm_client_id", id)
		if err != nil {
			t.Fatalf("replaceConfigMember() error = %v", err)
		}
		cfg := moduleConfigByName(t, out, "cosmos_to_eth")
		if got := cfg["cosmos_wasm_client_id"]; got != id {
			t.Fatalf("cosmos_wasm_client_id = %v, want %s", got, id)
		}
		if got := cfg["keep"]; got != "yes" {
			t.Fatalf("keep = %v, want yes", got)
		}
	})
}

func moduleConfigByName(t *testing.T, data []byte, name string) map[string]any {
	t.Helper()

	var cfg struct {
		Modules []struct {
			Name   string         `json:"name"`
			Config map[string]any `json:"config"`
		} `json:"modules"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, string(data))
	}
	for _, m := range cfg.Modules {
		if m.Name == name {
			if m.Config == nil {
				t.Fatalf("module %q config is nil", name)
			}
			return m.Config
		}
	}
	t.Fatalf("module %q not found", name)
	return nil
}

// TestReplaceConfigMemberForSource confirms a write-back targets only the named
// source's module and leaves the other source untouched.
func TestReplaceConfigMemberForSource(t *testing.T) {
	t.Parallel()

	in := `{"modules":[` +
		`{"name":"cosmos_to_eth","config":{"ics26_client_id":"chain-a","ics07_client":"0xAAA"}},` +
		`{"name":"cosmos_to_eth","config":{"ics26_client_id":"chain-b","ics07_client":"0xBBB"}}` +
		`]}`
	out, err := replaceConfigMemberForSource([]byte(in), "chain-b", "ics07_client", "0xNEW")
	if err != nil {
		t.Fatalf("replaceConfigMemberForSource() error = %v", err)
	}
	got := string(out)
	if !strings.Contains(got, `"chain-a","ics07_client":"0xAAA"`) {
		t.Fatalf("chain-a should be untouched, got: %s", got)
	}
	if !strings.Contains(got, `"chain-b","ics07_client":"0xNEW"`) {
		t.Fatalf("chain-b should be updated, got: %s", got)
	}

	// TestReplaceConfigMemberForSource/unknown confirms an unknown source errors
	// rather than silently writing the wrong module.
	t.Run("unknown", func(t *testing.T) {
		t.Parallel()

		in := `{"modules":[{"name":"cosmos_to_eth","config":{"ics26_client_id":"chain-a"}}]}`
		_, err := replaceConfigMemberForSource([]byte(in), "chain-x", "ics07_client", "0x1")
		if err == nil || !strings.Contains(err.Error(), "chain-x") {
			t.Fatalf("error = %v, want unknown-source error", err)
		}
	})

	// TestReplaceConfigMemberForSource/via src chain confirms the module-level src_chain
	// is used as the source id when config.ics26_client_id is absent (matching how
	// loadConfig defaults the id).
	t.Run("via src chain", func(t *testing.T) {
		t.Parallel()

		in := `{"modules":[{"name":"cosmos_to_eth","src_chain":"chain-a","config":{"ics07_client":"0xAAA"}}]}`
		out, err := replaceConfigMemberForSource([]byte(in), "chain-a", "ics07_client", "0xNEW")
		if err != nil {
			t.Fatalf("replaceConfigMemberForSource() error = %v", err)
		}
		if !strings.Contains(string(out), `"ics07_client":"0xNEW"`) {
			t.Fatalf("src_chain match failed, got: %s", string(out))
		}
	})

	// TestReplaceConfigMemberForSource/canonical is the regression for the canonical
	// schema: the module carries a free-form name label ("cosmos-to-eth") and is
	// classified by src_chain/dst_chain, exactly as loadConfig/classifyModule does.
	// The old write-back matched name == "cosmos_to_eth" and so failed to find a
	// canonical module after create-clients mutated chain state. It must also skip the
	// eth->cosmos module even when it shares the source id.
	t.Run("canonical", func(t *testing.T) {
		t.Parallel()

		in := `{"modules":[` +
			`{"name":"eth-to-cosmos","src_chain":"ethereum","dst_chain":"cosmos","config":{"ics26_client_id":"chain-a","ics07_client":"0xETH"}},` +
			`{"name":"cosmos-to-eth","src_chain":"cosmos","dst_chain":"ethereum","config":{"ics26_client_id":"chain-a","ics07_client":"0xAAA"}}` +
			`]}`
		out, err := replaceConfigMemberForSource([]byte(in), "chain-a", "ics07_client", "0xNEW")
		if err != nil {
			t.Fatalf("replaceConfigMemberForSource() error = %v", err)
		}
		got := string(out)
		if !strings.Contains(got, `"ics26_client_id":"chain-a","ics07_client":"0xNEW"`) {
			t.Fatalf("canonical cosmos->eth module not updated, got: %s", got)
		}
		// The eth->cosmos module (same source id, wrong direction) must be untouched.
		if !strings.Contains(got, `"ics07_client":"0xETH"`) {
			t.Fatalf("eth->cosmos module must be skipped, got: %s", got)
		}
	})
}

// selectCosmosSource is the promotion half of what --source used to do: the
// create-clients and submit-misbehaviour flows read CosmosToEthConfig, and a
// cosmos_to_l2 module is the same struct pointed at the L2, so a file declaring
// only that module has it promoted into the field those flows read. There is no
// disambiguation half any more -- loadConfig refuses a file with two paths.
func TestSelectCosmosSource(t *testing.T) {
	t.Parallel()

	t.Run("a cosmos_to_eth module is read as is", func(t *testing.T) {
		single := &appConfig{
			HasCosmosToEth: true, CosmosToEthConfig: cosmosToEthConfig{ICS26ClientID: "solo", TmRpcUrl: "s"},
		}
		if got := selectCosmosSource(single); got.CosmosToEthConfig.TmRpcUrl != "s" {
			t.Fatalf("selected wrong source: %+v", got.CosmosToEthConfig)
		}
	})

	// create-clients-eth also targets cosmos_to_l2 destinations (the L2-side
	// SpectreClient is created by the same flow).
	t.Run("a sole cosmos_to_l2 module is promoted", func(t *testing.T) {
		l2Only := &appConfig{
			HasCosmosToL2: true, CosmosToL2Config: cosmosToEthConfig{ICS26ClientID: "l2-solo", TmRpcUrl: "l2s"},
		}
		got := selectCosmosSource(l2Only)
		if got.CosmosToEthConfig.TmRpcUrl != "l2s" {
			t.Fatalf("selected wrong dest: %+v", got.CosmosToEthConfig)
		}
		// A copy, not a write into the caller's config: the L2 module must not
		// start looking like a cosmos_to_eth one to anything that reads the flag.
		if l2Only.HasCosmosToEth || l2Only.CosmosToEthConfig.TmRpcUrl != "" {
			t.Fatalf("selectCosmosSource mutated its input: %+v", l2Only)
		}
	})

	t.Run("nothing declared is returned unchanged", func(t *testing.T) {
		empty := &appConfig{}
		if got := selectCosmosSource(empty); got != empty {
			t.Fatal("a config with no Cosmos→EVM module must come back as is, for the command's own preflight to report")
		}
	})
}

// An eth_to_cosmos module without a websocket used to start and then run half-dead:
// the forward direction relayed, SubscribeEth bailed out after one log line, and no
// acknowledgement was ever picked up. Refuse to START on that config.
func TestValidateRelayStartupConfig_RequiresEthWs(t *testing.T) {
	const cfg = `{
		"modules": [
			{"name": "cosmos_to_eth", "src_chain": "cosmos", "dst_chain": "ethereum",
			 "config": {"tm_rpc_url": "http://localhost:26657", "eth_rpc_url": "http://localhost:8545",
			            "ics26_address": "0x80741a37e3644612f0465145c9709a90b6d77ee3", "ics26_client_id": "chain-a"}},
			{"name": "eth_to_cosmos", "src_chain": "ethereum", "dst_chain": "cosmos",
			 "config": {"eth_beacon_api_url": "http://localhost:5052"}}
		]
	}`
	appCfg, err := loadConfig(writeTempConfig(t, cfg))
	if err != nil {
		t.Fatalf("loadConfig must accept this — create-clients-* never read eth_ws_url: %v", err)
	}
	err = validateRelayStartupConfig(appCfg)
	if err == nil {
		t.Fatal("start accepted eth_to_cosmos without eth_ws_url; the ETH→Cosmos direction cannot run")
	}
	if !strings.Contains(err.Error(), "eth_ws_url") {
		t.Fatalf("error should name the missing field, got: %v", err)
	}
}

// The prefix has to land AFTER the timestamp. Without log.Lmsgprefix the
// standard logger writes it first, so every line starts with the path and the
// times stop lining up as a column — which is how a merged log is actually read.
func TestStampRelayPath_OnLogsPutsThePathAfterTheTimestamp(t *testing.T) {
	origFlags, origPrefix := log.Flags(), log.Prefix()
	origOut := log.Writer()
	t.Cleanup(func() {
		log.SetFlags(origFlags)
		log.SetPrefix(origPrefix)
		log.SetOutput(origOut)
	})

	var buf bytes.Buffer
	log.SetOutput(&buf)
	stampRelayPath(zap.NewNop(), "cosmos<->eth/cosmoshub-1")
	log.Printf("[Subscribe] send_packet received: seq=7")

	line := buf.String()
	if !strings.Contains(line, "[cosmos<->eth/cosmoshub-1]") {
		t.Fatalf("line %q does not carry the path", line)
	}
	if !strings.Contains(line, "[Subscribe] send_packet received: seq=7") {
		t.Fatalf("line %q lost the original message", line)
	}
	pathAt := strings.Index(line, "[cosmos<->eth/cosmoshub-1]")
	if pathAt == 0 {
		t.Fatalf("the path is written before the timestamp; log.Lmsgprefix is missing: %q", line)
	}
	// The timestamp must be everything before the path, not a substring of it.
	if stamp := strings.TrimSpace(line[:pathAt]); stamp == "" {
		t.Fatalf("nothing precedes the path, so the timestamp was dropped: %q", line)
	}
}

// The other half of the same call: the injected zap logger must carry the path
// too. Found in review -- the first version stamped only log.Default(), so every
// lifecycle line Start emits through zap ("Relayer started", "shutdown
// requested", "clients stopped") stayed unattributable in a merged log, and the
// test then in place only looked at log.Prefix() so it could not see that.
func TestStampRelayPath_ScopesTheZapLoggerToo(t *testing.T) {
	origFlags, origPrefix, origOut := log.Flags(), log.Prefix(), log.Writer()
	t.Cleanup(func() {
		log.SetFlags(origFlags)
		log.SetPrefix(origPrefix)
		log.SetOutput(origOut)
	})
	log.SetOutput(&bytes.Buffer{})

	core, logs := observer.New(zapcore.InfoLevel)
	scoped := stampRelayPath(zap.New(core), "cosmos<->eth/cosmoshub-1")
	scoped.Info("Relayer started")

	entries := logs.All()
	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}
	got, ok := entries[0].ContextMap()[relayPathLogField]
	if !ok {
		t.Fatalf("the zap line carries no %q field: %+v; a merged log cannot attribute it "+
			"to a relay path", relayPathLogField, entries[0].ContextMap())
	}
	if got != "cosmos<->eth/cosmoshub-1" {
		t.Fatalf("%s = %v, want the relay path", relayPathLogField, got)
	}
}

// stringLiterals returns the contents of each double-quoted literal on a line,
// so a check meant for format strings cannot read the code around them.
//
// Escapes are honoured only far enough to keep \" from ending a literal early;
// this reads Go source for one narrow property, not to compile it.
func stringLiterals(line string) []string {
	var out []string
	inside, escaped := false, false
	start := 0
	for i, r := range line {
		switch {
		case escaped:
			escaped = false
		case r == '\\' && inside:
			escaped = true
		case r == '"':
			if inside {
				out = append(out, line[start:i])
			} else {
				start = i + 1
			}
			inside = !inside
		}
	}
	return out
}

// Two label families are the target: [<direction> <Work>] and the closed
// lowercase list. Neither can contain "][", which is how a caller fakes a second
// tier when the function it calls is missing a parameter (subscriber/event.go
// passed "SubscribeCosmos][recovery" as a label). This fails on the forgery, not
// on the confusing log line it produces months later.
//
// [ATTENTION] is exempt, and is the one exemption. It is not a forged label: it
// is a fixed severity tier this repository already uses in 12 files, arriving
// over #377, #395, #415 and #457, and it is written the way this test asks for --
// the label is a parameter, the tier a literal after it. Forbidding it would mean
// rewriting a settled convention to satisfy a gate aimed at something else.
func TestNoLogLabelForgesASecondTier(t *testing.T) {
	root := ".."
	var offenders []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if name := info.Name(); name == "bindings" || name == "third_party" || name == "bin" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(src), "\n") {
			trimmed := strings.TrimSpace(line)
			if !strings.Contains(line, "][") || strings.HasPrefix(trimmed, "//") {
				// A comment describing a past forgery is not one.
				continue
			}
			// "][" is ordinary Go outside a string literal (`map[k][]byte`,
			// `x[i][j]`). Only a quoted occurrence is a forged label, so the
			// literals are extracted and nothing else on the line is read.
			//
			// Splitting on the first quote and scanning the remainder is not the
			// same thing, and the difference is a false positive: on
			//   dst.UpdateClient(ctx, "client-0", chain.ClientUpdate{Payloads: [][]byte{{1}}})
			// the [][]byte sits after a string literal and was reported as a
			// forged label.
			for _, literal := range stringLiterals(line) {
				// Drop the one allowed tier before looking, so "[%s][ATTENTION][x]"
				// is still caught: only the exemption itself is removed.
				if rest := strings.ReplaceAll(literal, "][ATTENTION]", "]"); strings.Contains(rest, "][") {
					offenders = append(offenders, fmt.Sprintf("%s:%d: %s", path, i+1, trimmed))
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	if len(offenders) > 0 {
		t.Fatalf("log labels forging a second tier:\n%s\n\nA label that needs a second level needs a "+
			"parameter, not a closing bracket in the middle of the format string.", strings.Join(offenders, "\n"))
	}
}

// The two noise words are gone and must stay gone: the direction is already the
// label (or the process prefix), so "relay" and "adapter" said nothing.
func TestNoLogLabelRepeatsTheLayerName(t *testing.T) {
	for _, banned := range []string{`"[relay `, `"[adapter `, `"\n[relay `, `"\n[adapter `} {
		out, err := exec.Command("grep", "-rl", "--include=*.go", banned, "..").Output()
		if err == nil && len(strings.TrimSpace(string(out))) > 0 {
			t.Errorf("label %q is back in:\n%s", strings.TrimSpace(banned), out)
		}
	}
}

// The point of E3 is one concept, one name. Four names for the timeout scanner
// and three for the Ethereum sender is not a formatting problem: an operator
// grepping [SendEthTx] silently missed a third of the sender's own lines.
//
// This pins the closed set. Adding a label is deliberate and cheap -- add it
// here too. Re-inventing a name for something already named is what this stops.
func TestLogLabelsStayOnTheAgreedList(t *testing.T) {
	allowed := map[string]bool{
		// [<direction> <Work>] -- direction is literal where the code knows it,
		// and "%s" where it comes from the module's own name.
		"eth->cosmos Subscribe": true, "l2->cosmos Subscribe": true,
		"cosmos->eth Relay": true, "cosmos->l2 Relay": true, "l2->cosmos Relay": true,
		"cosmos->eth UpdateClient": true, "cosmos->l2 UpdateClient": true,
		"eth->cosmos UpdateClient": true,
		"%s Relay":                 true, "%s UpdateClient": true,

		// [<Work>] -- the emitting code serves more than one direction, so the
		// chain stays in the name: it is the only thing left that separates them.
		"SubscribeCosmos": true, "SendEthTx": true, "SendCosmosTx": true,
		"CreateEthClient": true, "CreateCosmosClient": true,
		"UpdateCosmosClient": true, "UpdateEthClient": true, "RefreshCosmosClient": true,
		"CosmosTimeoutScan": true, "%sTimeoutScan": true,
		"Misbehaviour": true, "PendingTracker": true,
		// Added on purpose, which is what the message below asks for. #452's
		// enumeration backstop and #454's owed-ack ledger both landed after the
		// E3 renaming pass ran, so their labels were never put through it:
		// FlushCosmos is chain-qualified for the same reason SubscribeCosmos is
		// (cosmos.NewSource serves both cosmos->eth and cosmos->l2), and AckWatch
		// replaces a bare [adapter], which is one of the two noise words E3
		// removed.
		"FlushCosmos": true, "AckWatch": true,
		// The enumeration backstop's EVM halves, one per router this process reads:
		// FlushEth for the L1 source, FlushL2 for an L2 source. Chain-qualified for
		// the same reason FlushCosmos is.
		"FlushEth": true, "FlushL2": true,

		// The closed lowercase list: shared by every path, so a direction would
		// be a lie.
		// `create-clients` sits here for the same reason as `start`: it names a CLI
		// command that serves both directions, so a direction in the label would be
		// a lie.
		"bench": true, "prover": true, "BatchBuilder": true, "QueueState": true, "start": true,
		// The status endpoint serves the whole process, so a direction would be a
		// lie -- the same reason start and create-clients sit here.
		"status":         true,
		"create-clients": true,

		// Deleted by #417 along with the unbound-header fallback it belongs to;
		// touching it here would only conflict with that PR.
		"%s": true,
	}

	labelRe := regexp.MustCompile(`\.(?:Printf|Println|Fatalf|Infof|Warnf|Errorf|Debugf)\("(?:\\n)?\[([^\]"]+)\]`)
	var offenders []string
	err := filepath.Walk("..", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if n := info.Name(); n == "bindings" || n == "third_party" || n == "bin" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(src), "\n") {
			for _, m := range labelRe.FindAllStringSubmatch(line, -1) {
				if !allowed[m[1]] {
					offenders = append(offenders, fmt.Sprintf("%s:%d: [%s]", path, i+1, m[1]))
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	if len(offenders) > 0 {
		t.Fatalf("log labels outside the agreed set:\n%s\n\nEither reuse the name this concept already "+
			"has, or add the new one to this list on purpose.", strings.Join(offenders, "\n"))
	}
}

// The per-path half of the same contract: the goroutine's own cleanup must
// stand down too. It used to be an unconditional `defer cleanup()`, which fired
// the instant the engine returned -- and onceCleanup made it the winning call,
// so the central skip never got a say.
func TestReleaseOrHold(t *testing.T) {
	for _, tc := range []struct {
		name       string
		err        error
		wantClosed bool
		wantFlag   bool
	}{
		{name: "a clean stop closes the clients", err: nil, wantClosed: true},
		{name: "an ordinary failure still closes them", err: errors.New("rpc broke"), wantClosed: true},
		{
			name: "a stuck worker holds them open and raises the flag",
			err:  errWrap(workergroup.ErrWorkersStillRunning),
			// no close, and the central cleanups must stand down as well
			wantFlag: true,
		},
		{
			name:     "a stuck worker joined with a real failure still holds them",
			err:      errors.Join(errors.New("proof failed"), errWrap(workergroup.ErrWorkersStillRunning)),
			wantFlag: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			closed := false
			var flag atomic.Bool
			releaseOrHold(zap.NewNop(), func() { closed = true }, &flag, tc.err, "cosmos_to_eth", "client-0")

			if closed != tc.wantClosed {
				t.Fatalf("clients closed = %v, want %v", closed, tc.wantClosed)
			}
			if flag.Load() != tc.wantFlag {
				t.Fatalf("workersStuck = %v, want %v: the central cleanups read this, so a "+
					"missed flag closes the same clients a moment later", flag.Load(), tc.wantFlag)
			}
		})
	}

	// P1c, found in review. stopRelaysAndCleanup already refuses to close clients
	// when its own wait overruns, for a stated reason: "a goroutine that has not
	// returned may still be using them, and closing underneath it trades a slow
	// shutdown for a crash". Its wait cannot see a worker orphaned INSIDE a source,
	// because the engine returns as soon as its own drain times out and the
	// orphaned goroutine is in no WaitGroup. The wait then succeeds and the
	// cleanups run over a live worker -- the same crash, one level down.
	t.Run("cleanup is held when a source could not stop its workers", func(t *testing.T) {
		t.Run("held when a source reports stuck workers", func(t *testing.T) {
			cleanupRan := false
			err := stopRelaysAndCleanup(
				func() {}, func() {},
				[]func(){func() { cleanupRan = true }},
				func() bool { return true },
			)
			if cleanupRan {
				t.Fatal("closed the clients while a worker was still using them; " +
					"that is the crash the drain-overrun branch already refuses to cause")
			}
			if err == nil {
				t.Fatal("a shutdown that left clients open must be reported, or it looks clean")
			}
			if !strings.Contains(err.Error(), "leaving clients open") {
				t.Fatalf("error %q does not say what was skipped or why", err)
			}
		})

		t.Run("runs when every worker stopped", func(t *testing.T) {
			cleanupRan := false
			err := stopRelaysAndCleanup(
				func() {}, func() {},
				[]func(){func() { cleanupRan = true }},
				func() bool { return false },
			)
			if err != nil {
				t.Fatalf("clean shutdown reported an error: %v", err)
			}
			if !cleanupRan {
				t.Fatal("clients were never closed on a clean shutdown; the guard is too eager " +
					"and every ordinary stop now leaks its connections")
			}
		})
	})
}

func errWrap(err error) error { return errors.Join(errors.New("cosmos source"), err) }

// awaitRelayExit decides the process exit code. The SIGTERM branch is the one
// that was wrong: a module reports the workers it could not drain, that error is
// recorded, and the branch returned nil anyway -- a stuck worker looked exactly
// like a clean stop.
func TestAwaitRelayExit(t *testing.T) {
	logger := zap.NewNop()
	noopStop := func() error { return nil }

	t.Run("clean stop returns nil", func(t *testing.T) {
		done := make(chan struct{})
		close(done)
		if err := awaitRelayExit(context.Background(), done, newEngineErrors(func() {}), noopStop, logger); err != nil {
			t.Fatalf("clean stop returned %v", err)
		}
	})

	// The failure reaches this function through done, not through an arm of its
	// own: recording cancels the sibling engines, so every engine returns and the
	// wait group closes done. The assertion is that "all engines returned" no
	// longer reads as a clean exit when one of them returned because it failed.
	t.Run("an engine failure is returned", func(t *testing.T) {
		done := make(chan struct{})
		close(done)
		engineErrs := newEngineErrors(func() {})
		engineErrs.record(errors.New("cosmos_to_eth source \"cosmos-0\": boom"))

		err := awaitRelayExit(context.Background(), done, engineErrs, noopStop, logger)
		if err == nil || !strings.Contains(err.Error(), "boom") {
			t.Fatalf("engine failure = %v, want the engine's own error", err)
		}
	})

	t.Run("SIGTERM surfaces a recorded drain failure", func(t *testing.T) {
		// Exactly the shape relay.Module.Run produces on an overrun.
		const stuck = `relay cosmos->eth: shutdown drain timed out after 30s; workers still running: [cosmos-subscribe]`
		engineErrs := newEngineErrors(func() {})
		engineErrs.record(errors.New(stuck))

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := awaitRelayExit(ctx, make(chan struct{}), engineErrs, noopStop, logger)
		if err == nil {
			t.Fatal("SIGTERM with a stuck worker exited cleanly; that is the bug this branch had")
		}
		if !strings.Contains(err.Error(), "cosmos-subscribe") {
			t.Fatalf("error %q does not name the worker that is stuck", err)
		}
	})

	t.Run("SIGTERM surfaces a shutdown that overran its budget", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		stopErr := errors.New("relay engines did not drain within 30s")

		err := awaitRelayExit(ctx, make(chan struct{}), newEngineErrors(func() {}),
			func() error { return stopErr }, logger)
		if !errors.Is(err, stopErr) {
			t.Fatalf("SIGTERM returned %v, want the shutdown failure", err)
		}
	})

	// The control for the two above: an ordinary SIGTERM with nothing wrong exits
	// 0. Without it, "return the error" is satisfied by always returning one.
	t.Run("an ordinary SIGTERM exits cleanly", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := awaitRelayExit(ctx, make(chan struct{}), newEngineErrors(func() {}), noopStop, logger); err != nil {
			t.Fatalf("ordinary SIGTERM returned %v, want nil", err)
		}
	})

	// P1b: awaitRelayExit existed, was documented as the fix, was tested -- and
	// Start never called it, so the SIGTERM ending still dropped a recorded failure
	// and returned nil. This pins the contract awaitRelayExit implements, so the
	// wiring test below has something to be the wiring OF.
	t.Run("reports a failure recorded during shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // SIGTERM already landed

		engineErrs := newEngineErrors(func() {})
		engineErrs.record(errors.New("cosmos source: workers still running"))

		err := awaitRelayExit(ctx, make(chan struct{}), engineErrs,
			func() error { return nil }, zap.NewNop())
		if err == nil {
			t.Fatal("exit code 0 with a recorded failure: the relayer reports success " +
				"and the operator never learns a worker was stuck")
		}
	})
}

func TestStopRelaysAndCleanup(t *testing.T) {
	t.Run("client cleanup runs after relay drain", func(t *testing.T) {
		var order []string
		stopRelaysAndCleanup(
			func() { order = append(order, "cancel") },
			func() { order = append(order, "drain") },
			[]func(){func() { order = append(order, "cleanup") }},
			nil,
		)
		want := []string{"cancel", "drain", "cleanup"}
		if fmt.Sprint(order) != fmt.Sprint(want) {
			t.Fatalf("order = %v, want %v", order, want)
		}
	})

	// A worker that will not return must not hold the process. Before this, wait()
	// was called bare: the relayer sat holding its Cosmos account sequence and its
	// ETH nonce, and an operator restarting it collided with both.
	t.Run("gives up on a stuck drain", func(t *testing.T) {
		withShortShutdownBudgets(t, 40*time.Millisecond, time.Second)

		release := make(chan struct{})
		t.Cleanup(func() { close(release) })
		cleanupRan := false

		start := time.Now()
		err := stopRelaysAndCleanup(
			func() {},
			func() { <-release }, // the worker that refuses to exit
			[]func(){func() { cleanupRan = true }},
			nil,
		)
		elapsed := time.Since(start)

		if err == nil {
			t.Fatal("a drain that never returned was reported as a clean shutdown")
		}
		if !strings.Contains(err.Error(), "drain") {
			t.Fatalf("error %q does not say which phase overran", err)
		}
		if elapsed > 2*shutdownDrainBudget {
			t.Fatalf("waited %s for a %s budget", elapsed, shutdownDrainBudget)
		}
		// Deliberate: a goroutine that has not returned may still be using these
		// clients, and closing underneath it trades a slow shutdown for a crash.
		if cleanupRan {
			t.Fatal("clients were closed while a worker was still running")
		}
	})

	// The cleanup phase is bounded too -- cosmosClient.Stop() and ethClient.Close()
	// talk to sockets, so neither is guaranteed to return.
	t.Run("gives up on a stuck cleanup", func(t *testing.T) {
		withShortShutdownBudgets(t, time.Second, 40*time.Millisecond)

		release := make(chan struct{})
		t.Cleanup(func() { close(release) })

		err := stopRelaysAndCleanup(func() {}, func() {}, []func(){func() { <-release }}, nil)
		if err == nil {
			t.Fatal("a cleanup that never returned was reported as a clean shutdown")
		}
		if !strings.Contains(err.Error(), "cleanup") {
			t.Fatalf("error %q does not say which phase overran", err)
		}
	})

	// The control case: nothing is stuck, so nothing is reported. Without this the
	// tests above are satisfied by a function that always errors.
	t.Run("reports nothing on a clean stop", func(t *testing.T) {
		withShortShutdownBudgets(t, time.Second, time.Second)

		if err := stopRelaysAndCleanup(func() {}, func() {}, []func(){func() {}}, nil); err != nil {
			t.Fatalf("a clean shutdown reported %v", err)
		}
	})
}

func TestIsShutdownErr(t *testing.T) {
	for _, tc := range []struct {
		err  error
		want bool
	}{
		{context.Canceled, true},
		{context.DeadlineExceeded, true},
		{fmt.Errorf("wrapped: %w", context.Canceled), true},
		{errors.New("request was cancelled by server"), false},
		{nil, false},
	} {
		if got := isShutdownErr(tc.err); got != tc.want {
			t.Fatalf("isShutdownErr(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}

// withShortShutdownBudgets runs the same code path in milliseconds. The budgets
// are the only thing shortened -- every branch under test is the production one.
func withShortShutdownBudgets(t *testing.T, drain, cleanup time.Duration) {
	t.Helper()
	origDrain, origCleanup := shutdownDrainBudget, shutdownCleanupBudget
	shutdownDrainBudget, shutdownCleanupBudget = drain, cleanup
	t.Cleanup(func() { shutdownDrainBudget, shutdownCleanupBudget = origDrain, origCleanup })
}

// The budget an operator is promised. Nothing reads this at runtime, so a
// change to either half would otherwise be invisible until a SIGTERM in
// production took longer than the runbook says it can.
func TestShutdownBudget_IsFortyFiveSeconds(t *testing.T) {
	if shutdownBudget != 45*time.Second {
		t.Fatalf("shutdown budget = %s, want 45s (docs and the runbook state this number)", shutdownBudget)
	}
	if shutdownDrainBudget != workers.ShutdownDrainTimeout {
		t.Fatalf("drain budget = %s, want the per-module timeout %s: modules drain in parallel, so the "+
			"two are the same bound and must not drift", shutdownDrainBudget, workers.ShutdownDrainTimeout)
	}
}
