package main

import (
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"

	"relayer/services"
)

// FETCH_TIMEOUT bounds the Cosmos RPC calls. It used to be read as
// `if d, err := strconv.Atoi(envVal); err == nil && d > 0`, which treated BOTH a
// parse failure and a non-positive value as "not set" -- so an operator who wrote
// FETCH_TIMEOUT=30s, the way a Go duration is normally written and the natural
// guess for a field that IS a time.Duration, silently got the default.
//
// An RPC running without the deadline the operator asked for is the failure that
// wedged a whole relay direction before, so the wrong value must be refused
// rather than quietly replaced.
func TestEnvSecondsDuration(t *testing.T) {
	const name = "TEST_RELAYER_TIMEOUT_KNOB"

	t.Run("unset reports not-set without an error", func(t *testing.T) {
		t.Setenv(name, "")
		got, set, err := envSecondsDuration(name)
		if err != nil {
			t.Fatalf("unset must not be an error, got %v", err)
		}
		if set {
			t.Fatalf("reported set with value %s", got)
		}
	})

	t.Run("parses whole seconds", func(t *testing.T) {
		t.Setenv(name, "30")
		got, set, err := envSecondsDuration(name)
		if err != nil || !set || got != 30*time.Second {
			t.Fatalf("envSecondsDuration = %s, %v, %v; want 30s, true, nil", got, set, err)
		}
	})

	t.Run("refuses a value it cannot honour", func(t *testing.T) {
		for label, raw := range map[string]string{
			"go duration syntax": "30s",
			"milliseconds":       "500ms",
			"digits then unit":   "30sec",
			"zero":               "0",
			"negative":           "-1",
			"words":              "none",
			"decimal":            "1.5",
		} {
			t.Run(label, func(t *testing.T) {
				t.Setenv(name, raw)
				got, set, err := envSecondsDuration(name)
				if err == nil {
					t.Fatalf("envSecondsDuration(%q) = %s, set=%v with no error", raw, got, set)
				}
				// Must not look like an unset variable either: the caller applies
				// the default on `set == false`, so reporting set=false here would
				// reintroduce the silent fallback this replaced.
				if set {
					t.Fatalf("envSecondsDuration(%q) reported set alongside an error", raw)
				}
				if !strings.Contains(err.Error(), name) {
					t.Errorf("error %q does not name the variable the operator set", err)
				}
			})
		}
	})

	t.Run("surrounding whitespace is tolerated", func(t *testing.T) {
		// Trimmed on purpose: a trailing newline from `export FETCH_TIMEOUT=$(...)`
		// is an artefact of the shell, not something the operator typed.
		t.Setenv(name, "  30\n")
		got, set, err := envSecondsDuration(name)
		if err != nil || !set || got != 30*time.Second {
			t.Fatalf("envSecondsDuration = %s, %v, %v; want 30s, true, nil", got, set, err)
		}
	})
}

// buildCosmosConfig must propagate the refusal rather than swallow it -- the
// whole point of the change is that a bad override stops the process instead of
// being replaced by a default nobody chose.
func TestBuildCosmosConfig_RejectsAMalformedFetchTimeout(t *testing.T) {
	t.Setenv("FETCH_TIMEOUT", "30s")
	if _, err := buildCosmosConfig(cosmosToEthConfig{}, services.BatchConfig{}); err == nil {
		t.Fatal("a malformed FETCH_TIMEOUT was accepted")
	}

	t.Setenv("FETCH_TIMEOUT", "30")
	cfg, err := buildCosmosConfig(cosmosToEthConfig{}, services.BatchConfig{})
	if err != nil {
		t.Fatalf("a valid FETCH_TIMEOUT was rejected: %v", err)
	}
	if cfg.FetchTimeout != 30*time.Second {
		t.Fatalf("FetchTimeout = %s, want 30s", cfg.FetchTimeout)
	}
}

// Config that can be judged without a network must be judged BEFORE any client
// is dialled or started, because the early return has no way to release one:
// the cleanup closure is not built until the end of the builder.
//
// Making buildCosmosConfig fallible introduced exactly such a return between
// cosmosClient.Start() and that closure, leaking a live WebSocket subscription
// on an operator typo. Review caught it.
//
// The assertion is on WHICH error comes back, which is what pins the order.
// Every endpoint below is unreachable, so a builder that dialled first would
// report a connection failure; reporting the FETCH_TIMEOUT failure instead
// proves nothing was opened.
func TestBuildSources_ValidateConfigBeforeOpeningClients(t *testing.T) {
	t.Setenv("FETCH_TIMEOUT", "30s") // seconds are expected; "30s" is malformed

	const unreachable = "http://127.0.0.1:1"
	c2e := cosmosToEthConfig{
		EthRpcUrl:          unreachable,
		TmRpcUrl:           unreachable,
		ICS26ClientID:      "cosmos-0",
		CosmosWasmClientID: "08-wasm-0",
		SpectreClient:      "0x0000000000000000000000000000000000000001",
	}

	t.Run("cosmos_to_eth", func(t *testing.T) {
		_, _, cleanup, err := buildCosmosToEthSource(
			zap.NewNop(), c2e, ethToCosmosConfig{}, services.BatchConfig{}, nil, nil,
			buildCosmosToEthSourceOptions{startSubscriptions: true},
		)
		requireFetchTimeoutRefusal(t, err, cleanup)
	})

	t.Run("cosmos_to_l2", func(t *testing.T) {
		_, _, cleanup, err := buildCosmosToL2Dest(
			zap.NewNop(), c2e, services.BatchConfig{}, nil, nil, t.TempDir(),
		)
		requireFetchTimeoutRefusal(t, err, cleanup)
	})
}

func requireFetchTimeoutRefusal(t *testing.T, err error, cleanup func()) {
	t.Helper()
	if cleanup != nil {
		t.Error("a builder that failed on config returned a cleanup func, so it opened something")
		cleanup()
	}
	if err == nil {
		t.Fatal("a malformed FETCH_TIMEOUT was accepted")
	}
	if !strings.Contains(err.Error(), "FETCH_TIMEOUT") {
		t.Fatalf("error = %v; want the FETCH_TIMEOUT refusal. Any connection error here means a "+
			"client was opened before the config was judged, and the early return cannot close it", err)
	}
}

// strconv accepting a number is not the same as time.Duration holding it.
// time.Duration is an int64 of NANOseconds, so 9223372037 seconds -- a perfectly
// ordinary-looking integer -- wraps to a NEGATIVE duration.
//
// A negative timeout is worse than a wrong one because it means two different
// things depending on who reads it: context.WithTimeout cancels immediately,
// while services.fetchCtx substitutes its 15s default. One accepted setting,
// two behaviours, neither of them what the operator asked for.
func TestSecondsToDurationRefusesWhatDoesNotFit(t *testing.T) {
	const maxSeconds = uint64(math.MaxInt64 / int64(time.Second)) // 9223372036

	t.Run("the largest representable value is accepted", func(t *testing.T) {
		got, err := secondsToDuration("fetch_timeout", maxSeconds)
		if err != nil {
			t.Fatalf("the boundary value was rejected: %v", err)
		}
		if got <= 0 {
			t.Fatalf("duration = %v, want positive", got)
		}
	})

	t.Run("one past the boundary is refused", func(t *testing.T) {
		if _, err := secondsToDuration("fetch_timeout", maxSeconds+1); err == nil {
			t.Fatal("a value one second past the maximum was accepted")
		}
	})

	t.Run("MaxUint64 is refused rather than wrapping", func(t *testing.T) {
		got, err := secondsToDuration("fetch_timeout", math.MaxUint64)
		if err == nil {
			t.Fatalf("MaxUint64 was accepted as %v", got)
		}
		if !strings.Contains(err.Error(), "fetch_timeout") {
			t.Errorf("error %q does not name the setting the operator wrote", err)
		}
	})

	// The property, not just the three points above: nothing this function
	// returns without an error may be non-positive.
	t.Run("no accepted value is ever non-positive", func(t *testing.T) {
		for _, seconds := range []uint64{1, 30, 3600, maxSeconds - 1, maxSeconds, maxSeconds + 1, 1 << 62, math.MaxInt64, math.MaxUint64} {
			got, err := secondsToDuration("fetch_timeout", seconds)
			if err != nil {
				continue
			}
			if got <= 0 {
				t.Errorf("%d seconds was accepted as %v", seconds, got)
			}
		}
	})
}

// The env path is what an operator actually types, so the guard has to hold
// through it and not only in the helper underneath.
func TestEnvSecondsDurationRefusesAnOverflowingValue(t *testing.T) {
	t.Setenv("FETCH_TIMEOUT", "9223372036854775807")
	got, set, err := envSecondsDuration("FETCH_TIMEOUT")
	if err == nil {
		t.Fatalf("an overflowing FETCH_TIMEOUT was accepted: got=%v set=%v", got, set)
	}
	if set {
		t.Error("a rejected value was still reported as set")
	}

	t.Setenv("FETCH_TIMEOUT", "9223372036")
	if _, set, err := envSecondsDuration("FETCH_TIMEOUT"); err != nil || !set {
		t.Fatalf("the largest usable value was rejected: set=%v err=%v", set, err)
	}
}

// cosmos_wasm_client_id is the id of the Ethereum light client ON Cosmos, and
// only the eth->cosmos direction reads it. Requiring it unconditionally made
// submit-misbehaviour unusable for a cosmos_to_l2 source: selectSource offers L2
// sources, cosmos_to_l2 deliberately carries no such id, and the builder refused
// the evidence before looking at it.
//
// So the requirement belongs to the caller, and both directions of that are
// worth pinning -- a fix that simply dropped the check would satisfy the first
// case and silently remove the guard from `start`.
func TestBuildCosmosToEthSourceRequiresTheWasmClientIDOnlyWhenItIsUsed(t *testing.T) {
	// An L2-shaped source: no cosmos_wasm_client_id, unreachable endpoints. The
	// assertion is on WHICH error comes back, so an endpoint failure proves the
	// wasm-id check was passed rather than skipped along with everything else.
	l2Shaped := cosmosToEthConfig{
		ICS26ClientID: "cosmos-l2-0",
		SpectreClient: "0x0000000000000000000000000000000000000001",
		EthRpcUrl:     "http://127.0.0.1:1",
		TmRpcUrl:      "http://127.0.0.1:1",
	}

	t.Run("a one-shot caller may omit it", func(t *testing.T) {
		_, _, cleanup, err := buildCosmosToEthSource(
			zap.NewNop(), l2Shaped, ethToCosmosConfig{}, services.BatchConfig{}, nil, nil,
			buildCosmosToEthSourceOptions{},
		)
		if cleanup != nil {
			cleanup()
		}
		if err != nil && strings.Contains(err.Error(), "cosmos_wasm_client_id") {
			t.Fatalf("an L2 source was rejected for a field it does not have: %v", err)
		}
	})

	t.Run("a caller that relays eth->cosmos still needs it", func(t *testing.T) {
		_, _, cleanup, err := buildCosmosToEthSource(
			zap.NewNop(), l2Shaped, ethToCosmosConfig{}, services.BatchConfig{}, nil, nil,
			buildCosmosToEthSourceOptions{relaysEVMToCosmos: true},
		)
		if cleanup != nil {
			cleanup()
		}
		if err == nil || !strings.Contains(err.Error(), "cosmos_wasm_client_id") {
			t.Fatalf("error = %v; a relayer driving eth->cosmos must be stopped at startup, "+
				"not at the first packet", err)
		}
	})
}

// The guard is only worth having if `start` opts in, and nothing else checks
// that it does. Asserted structurally because the alternative is booting a real
// relayer.
func TestStartOptsIntoTheWasmClientIDRequirement(t *testing.T) {
	// start.go, not main.go: #459 moved the command bodies out of main.go, and
	// this test is about what `start` passes, not about where it lives.
	src, err := os.ReadFile("start.go")
	if err != nil {
		t.Fatalf("read start.go: %v", err)
	}
	body := string(src)
	start := strings.Index(body, "buildCosmosToEthSourceOptions{")
	if start < 0 {
		t.Fatal("start no longer builds a cosmos_to_eth source")
	}
	end := strings.Index(body[start:], "}")
	if !strings.Contains(body[start:start+end], "relaysEVMToCosmos: true") {
		t.Fatal("start does not opt into the cosmos_wasm_client_id requirement; a missing id " +
			"would then surface at the first eth->cosmos packet instead of at startup")
	}
}
