package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Found in review: Start named the state directory from the RAW config id while
// buildCosmosToEthSource relayed the EFFECTIVE one. The directory is the lock,
// the recovery cursors and the pending tracker, so the two must be the same
// function of the same inputs or the state describes a path nobody is serving.
func TestEffectiveRelayIDs(t *testing.T) {
	cfg := cosmosToEthConfig{ICS26ClientID: "from-json", CosmosWasmClientID: "wasm-from-json"}

	t.Run("the environment wins, for both ids", func(t *testing.T) {
		t.Setenv("ICS26_CLIENT_ID", "from-env")
		t.Setenv("COSMOS_WASM_CLIENT_ID", "wasm-from-env")
		ids := effectiveRelayIDs(cfg)
		router, wasm := ids.router, ids.wasm
		if router != "from-env" || wasm != "wasm-from-env" {
			t.Fatalf("got (%q, %q), want the environment values", router, wasm)
		}
	})

	t.Run("an unset override leaves the config value", func(t *testing.T) {
		t.Setenv("ICS26_CLIENT_ID", "")
		t.Setenv("COSMOS_WASM_CLIENT_ID", "")
		ids := effectiveRelayIDs(cfg)
		router, wasm := ids.router, ids.wasm
		if router != "from-json" || wasm != "wasm-from-json" {
			t.Fatalf("got (%q, %q); an empty override must not erase the configured id", router, wasm)
		}
	})
}

// The first consequence: changing an override must move the state directory,
// because the process is now relaying a different pair and must not inherit the
// previous pair's cursors and pending work.
func TestStateDirFollowsTheOverride(t *testing.T) {
	cfg := cosmosToEthConfig{ICS26ClientID: "from-json"}
	configPath := filepath.Join(t.TempDir(), "config.json")

	dirFor := func(override string) string {
		t.Setenv("ICS26_CLIENT_ID", override)
		return relayStateDir(configPath, "cosmos-to-eth", effectiveRelayIDs(cfg))
	}

	first := dirFor("client-a")
	second := dirFor("client-b")
	if first == second {
		t.Fatalf("both overrides resolved to %s: a restart under a new ICS26_CLIENT_ID would "+
			"replay the previous pair's pending packets under the wrong client", first)
	}

	// And the raw config value must not be what names it, or the two paths above
	// would have collided on the config's own id.
	raw := relayStateDir(configPath, "cosmos-to-eth", relayClientIDs{router: cfg.ICS26ClientID, wasm: cfg.CosmosWasmClientID})
	if first == raw || second == raw {
		t.Fatal("an overridden run reused the directory named for the raw config id")
	}
}

// The second consequence, and the one that breaks the guard rather than the
// state: two configs whose RAW ids differ but whose EFFECTIVE ids are the same
// serve one relay path, so they must contend for one lock. Naming the directory
// from the raw id gave them two, and both started.
func TestTwoConfigsOnOneEffectivePairShareOneStateDir(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("ICS26_CLIENT_ID", "the-real-client")

	one := effectiveRelayIDs(cosmosToEthConfig{ICS26ClientID: "stale-a"})
	two := effectiveRelayIDs(cosmosToEthConfig{ICS26ClientID: "stale-b"})

	dirOne := relayStateDir(configPath, "cosmos-to-eth", one)
	dirTwo := relayStateDir(configPath, "cosmos-to-eth", two)
	if dirOne != dirTwo {
		t.Fatalf("%s != %s: both processes relay %q, so they must contend for one lock; "+
			"two directories means both start and advance the same cursors",
			dirOne, dirTwo, "the-real-client")
	}
}

// Every value the environment may override for a source is decided in one
// place, beside effectiveRelayIDs -- which the role-manager address was not,
// after the identity fix moved the client ids to the caller and left this one
// reading the environment in the builder.
//
// It is not cosmetic: RouterManagesProofSubmission compares this address with
// the router to choose ICS26Router versus direct SpectreClient submission, so a
// stale value routes incident evidence down the wrong contract path.
func TestEffectiveRoleManager(t *testing.T) {
	cfg := cosmosToEthConfig{ICS26Address: "0xfromconfig"}

	t.Run("the environment wins", func(t *testing.T) {
		t.Setenv("ROLE_MANAGER", "0xfromenv")
		if got := effectiveRoleManager(cfg); got != "0xfromenv" {
			t.Fatalf("role manager = %q, want the override", got)
		}
	})

	t.Run("an unset variable falls back to the config", func(t *testing.T) {
		t.Setenv("ROLE_MANAGER", "")
		if got := effectiveRoleManager(cfg); got != "0xfromconfig" {
			t.Fatalf("role manager = %q, want the config's", got)
		}
	})
}

// The structural half. Both regressions on this PR were the same shape: a value
// the environment can override, resolved inside the builder rather than in the
// caller beside the others. The fix moved every one of them to the caller, and
// this is what keeps them there -- a new override added inside the builder makes
// this red instead of silently escaping the rule.
func TestBuilderReadsNoEnvironmentItself(t *testing.T) {
	src, err := os.ReadFile("build_source.go")
	if err != nil {
		t.Fatalf("read build_source.go: %v", err)
	}
	body := string(src)
	start := strings.Index(body, "func buildCosmosToEthSource(")
	if start < 0 {
		t.Fatal("buildCosmosToEthSource not found; this guard is pinned to its name")
	}
	// To the next top-level declaration, which is where the function ends.
	rest := body[start+len("func buildCosmosToEthSource("):]
	end := strings.Index(rest, "\nfunc ")
	if end < 0 {
		end = len(rest)
	}
	if strings.Contains(rest[:end], "envOrDefault(") {
		t.Fatal("buildCosmosToEthSource reads the environment directly; resolve it in the caller " +
			"beside effectiveRelayIDs/effectiveRoleManager, where every override is decided")
	}
}

// The state directory must move when EITHER half of the pair moves.
//
// Reported by @neitdung: it was named from the router id alone while the builder
// received {router, wasm}. Holding the config and the router fixed and changing
// only COSMOS_WASM_CLIENT_ID therefore relayed a different path out of the same
// directory -- the same lock, the same cursors, the same pending trackers. On
// restart the new wasm client inherits a cursor that skips its own older events,
// and pending packets from the previous pair replay under the new destination.
//
// Each half is changed on its own, because a test that changes both at once
// passes against a digest that reads only one of them.
func TestStateDirTracksBothHalvesOfTheEffectivePair(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	cfg := cosmosToEthConfig{ICS26ClientID: "cosmoshub-1", CosmosWasmClientID: "08-wasm-0"}

	dirFor := func(t *testing.T, router, wasm string) string {
		t.Helper()
		t.Setenv("ICS26_CLIENT_ID", router)
		t.Setenv("COSMOS_WASM_CLIENT_ID", wasm)
		return relayStateDir(configPath, "cosmos-to-eth", effectiveRelayIDs(cfg))
	}

	t.Run("the wasm half alone", func(t *testing.T) {
		a := dirFor(t, "cosmoshub-1", "08-wasm-0")
		b := dirFor(t, "cosmoshub-1", "08-wasm-7")
		if a == b {
			t.Fatalf("08-wasm-0 and 08-wasm-7 share state directory %q: one cursor, one lock, two relay paths", a)
		}
	})

	t.Run("the router half alone", func(t *testing.T) {
		a := dirFor(t, "cosmoshub-1", "08-wasm-0")
		b := dirFor(t, "osmosis-1", "08-wasm-0")
		if a == b {
			t.Fatalf("cosmoshub-1 and osmosis-1 share state directory %q", a)
		}
	})

	t.Run("the same pair, twice", func(t *testing.T) {
		a := dirFor(t, "cosmoshub-1", "08-wasm-0")
		b := dirFor(t, "cosmoshub-1", "08-wasm-0")
		if a != b {
			t.Fatalf("%s != %s: one identity must always reach one directory, or two processes never contend", a, b)
		}
	})
}

// Fields are framed by length, not separated by a byte. A separator is only
// unambiguous while no field can contain it, and these come from a JSON config.
// Two distinct identities colliding into one directory is the duplicate-owner
// case the lock exists to prevent, and it would show up as two processes quietly
// sharing one cursor rather than as an error.
func TestStateDirFramingCannotBeForgedByASeparatorInAClientID(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")

	// Under "\x00" separators both of these flatten to <path>\x00<kind>\x00a\x00b\x00c,
	// so they are the pair that actually collides. A shift with an empty field on
	// one side does not: the separator for that field still lands in the string.
	left := relayStateDir(configPath, "cosmos-to-eth", relayClientIDs{router: "a", wasm: "b\x00c"})
	right := relayStateDir(configPath, "cosmos-to-eth", relayClientIDs{router: "a\x00b", wasm: "c"})
	if left == right {
		t.Fatalf("{a, b\\x00c} and {a\\x00b, c} share state directory %q; the framing is forgeable", left)
	}
}
