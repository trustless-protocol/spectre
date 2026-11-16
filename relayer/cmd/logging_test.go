package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	gnarklogger "github.com/consensys/gnark/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"relayer/logbridge"
)

func envOf(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

// restoreGnarkLogger puts gnark's global logger back when the test ends. stampRelayPath
// installs the bridge, which also takes gnark's logger, and the relayer never undoes it;
// a test that reaches it would otherwise leave its own logger in place for every test
// after it.
func restoreGnarkLogger(t *testing.T) {
	t.Helper()
	previous := gnarklogger.Logger()
	t.Cleanup(func() { gnarklogger.Set(previous) })
}

func TestResolveLogSettings(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		env         map[string]string
		tty         bool
		wantLevel   zapcore.Level
		wantFormat  string
		wantIgnored []string // a substring of each expected report, in order
	}{
		{name: "defaults off a terminal", wantLevel: zapcore.InfoLevel, wantFormat: logFormatJSON},
		{name: "defaults on a terminal", tty: true, wantLevel: zapcore.InfoLevel, wantFormat: logFormatConsole},

		{name: "level from a flag", args: []string{"start", "--log-level", "debug"}, wantLevel: zapcore.DebugLevel, wantFormat: logFormatJSON},
		{name: "level from a flag with =", args: []string{"start", "--log-level=warn"}, wantLevel: zapcore.WarnLevel, wantFormat: logFormatJSON},
		{name: "level is case-insensitive", args: []string{"start", "--log-level", "ERROR"}, wantLevel: zapcore.ErrorLevel, wantFormat: logFormatJSON},
		{name: "level from the environment", env: map[string]string{envLogLevel: "error"}, wantLevel: zapcore.ErrorLevel, wantFormat: logFormatJSON},
		{
			name: "the flag wins over the environment", args: []string{"start", "--log-level", "warn"},
			env: map[string]string{envLogLevel: "debug"}, wantLevel: zapcore.WarnLevel, wantFormat: logFormatJSON,
		},
		{name: "an empty variable is unset", env: map[string]string{envLogLevel: "", envLogFormat: ""}, wantLevel: zapcore.InfoLevel, wantFormat: logFormatJSON},

		{name: "format from a flag", args: []string{"start", "--log-format", "console"}, wantLevel: zapcore.InfoLevel, wantFormat: logFormatConsole},
		{name: "format from the environment", env: map[string]string{envLogFormat: "console"}, wantLevel: zapcore.InfoLevel, wantFormat: logFormatConsole},
		{name: "json stays json on a terminal", tty: true, args: []string{"--log-format=json", "start"}, wantLevel: zapcore.InfoLevel, wantFormat: logFormatJSON},
		{name: "auto follows the terminal", tty: true, args: []string{"start", "--log-format", "auto"}, wantLevel: zapcore.InfoLevel, wantFormat: logFormatConsole},

		// Every other flag belongs to a subcommand this pass has not met, and none of
		// them may be mistaken for ours or swallow ours.
		{
			name:       "other flags, with values, do not disturb it",
			args:       []string{"start", "--config", "c.json", "--gpu-prove", "--log-level", "warn", "--benchmark"},
			wantLevel:  zapcore.WarnLevel,
			wantFormat: logFormatJSON,
		},
		{name: "a flag before the subcommand", args: []string{"--log-level", "debug", "start", "--config", "c.json"}, wantLevel: zapcore.DebugLevel, wantFormat: logFormatJSON},
		// pflag ends the parse at --help, as the command itself ends there by printing
		// its help, so a flag after it is not read and nothing is lost by that.
		{name: "a flag before --help is read", args: []string{"start", "--log-level", "warn", "--help"}, wantLevel: zapcore.WarnLevel, wantFormat: logFormatJSON},
		{name: "--help ends the parse", args: []string{"start", "--help", "--log-level", "warn"}, wantLevel: zapcore.InfoLevel, wantFormat: logFormatJSON},
		{name: "nothing after -- is a flag", args: []string{"start", "--", "--log-level", "debug"}, wantLevel: zapcore.InfoLevel, wantFormat: logFormatJSON},
		{name: "a flag missing its value falls back", args: []string{"start", "--log-level"}, wantLevel: zapcore.InfoLevel, wantFormat: logFormatJSON},

		// Never fatal: an unusable value is dropped, reported, and the default used.
		{
			name: "an invalid flag level is reported", args: []string{"start", "--log-level", "verbose"},
			wantLevel: zapcore.InfoLevel, wantFormat: logFormatJSON, wantIgnored: []string{`--log-level="verbose"`},
		},
		{
			name: "an invalid environment level is reported", env: map[string]string{envLogLevel: "loud"},
			wantLevel: zapcore.InfoLevel, wantFormat: logFormatJSON, wantIgnored: []string{envLogLevel + `="loud"`},
		},
		{
			name: "an invalid format is reported and falls back to auto", tty: true, args: []string{"start", "--log-format", "yaml"},
			wantLevel: zapcore.InfoLevel, wantFormat: logFormatConsole, wantIgnored: []string{`--log-format="yaml"`},
		},
		{
			name: "both invalid, both reported", args: []string{"start", "--log-level=x", "--log-format=y"},
			wantLevel: zapcore.InfoLevel, wantFormat: logFormatJSON, wantIgnored: []string{"--log-level", "--log-format"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveLogSettings(tc.args, envOf(tc.env), tc.tty)
			if got.level != tc.wantLevel || got.format != tc.wantFormat {
				t.Fatalf("level=%v format=%s, want level=%v format=%s", got.level, got.format, tc.wantLevel, tc.wantFormat)
			}
			if len(got.ignored) != len(tc.wantIgnored) {
				t.Fatalf("ignored = %q, want %d report(s) containing %q", got.ignored, len(tc.wantIgnored), tc.wantIgnored)
			}
			for i, want := range tc.wantIgnored {
				if !strings.Contains(got.ignored[i], want) {
					t.Errorf("ignored[%d] = %q, want it to contain %q", i, got.ignored[i], want)
				}
			}
		})
	}
}

// The flags are read ahead of cobra, so what cobra then does with them still has
// to work: every subcommand has to accept them, or a value the logger honoured
// would be rejected one step later.
func TestRegisterLogFlags_EverySubcommandAcceptsThem(t *testing.T) {
	root := &cobra.Command{Use: "relayer", SilenceUsage: true, SilenceErrors: true}
	registerLogFlags(root)
	ran := false
	root.AddCommand(&cobra.Command{Use: "start", Run: func(*cobra.Command, []string) { ran = true }})

	root.SetArgs([]string{"start", "--log-level", "debug", "--log-format=json"})
	if err := root.Execute(); err != nil || !ran {
		t.Fatalf("start rejected the logging flags (ran=%t): %v", ran, err)
	}

	root.SetArgs([]string{"start", "--no-such-flag"})
	if err := root.Execute(); err == nil {
		t.Fatal("an unknown flag was accepted; registering the logging flags must not loosen cobra's own parse")
	}
}

// A failed command's error must reach the operator only through the logger,
// whose core redacts URLs. Printed by cobra first, it went to stderr raw, the API
// key in the RPC URL included (#575).
func TestRootCmd_LeavesFailedCommandErrorsToTheLogger(t *testing.T) {
	const secret = "FAKESECRET123"
	root := newRootCmd(zap.NewNop())
	root.AddCommand(&cobra.Command{
		Use: "fail",
		RunE: func(*cobra.Command, []string) error {
			return fmt.Errorf("read eth_chainId at http://127.0.0.1:1/v3/%s: connection refused", secret)
		},
	})
	var printed bytes.Buffer
	root.SetOut(&printed)
	root.SetErr(&printed)
	root.SetArgs([]string{"fail"})

	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), secret) {
		t.Fatalf("Execute() = %v, want the command's error back for main to log", err)
	}
	if strings.Contains(printed.String(), secret) {
		t.Fatalf("cobra printed the error raw, API key included:\n%s", printed.String())
	}

	// main logs it through the production core (at fatal; Error here, same core).
	var logged bytes.Buffer
	newLogger(logSettings{level: zapcore.InfoLevel, format: logFormatJSON}, zapcore.AddSync(&logged)).Sugar().Error(err)
	if strings.Contains(logged.String(), secret) || !strings.Contains(logged.String(), "http://127.0.0.1:1") {
		t.Fatalf("logged %s, want the URL cut to scheme and host", logged.String())
	}
}

func lineOf(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var line map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &line); err != nil {
		t.Fatalf("not one JSON line: %v\n%s", err, buf.String())
	}
	return line
}

func TestNewLogger_JSONKeepsTheKeysZapProductionWrote(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(logSettings{level: zapcore.InfoLevel, format: logFormatJSON}, zapcore.AddSync(&buf))
	l.Warn("hello", zap.String(relayPathLogField, "p"))

	line := lineOf(t, &buf)
	for key, want := range map[string]any{"level": "warn", "msg": "hello", relayPathLogField: "p"} {
		if line[key] != want {
			t.Errorf("%s = %v, want %v", key, line[key], want)
		}
	}
	// A viewer already holds lines in this shape; renaming a key would split the stream.
	for _, key := range []string{"ts", "caller"} {
		if _, ok := line[key]; !ok {
			t.Errorf("no %q key in %v", key, line)
		}
	}
}

func TestNewLogger_ConsoleIsForAPerson(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(logSettings{level: zapcore.InfoLevel, format: logFormatConsole}, zapcore.AddSync(&buf))
	l.Warn("hello")

	out := buf.String()
	if strings.HasPrefix(out, "{") || !strings.Contains(out, "WARN") || !strings.Contains(out, "hello") {
		t.Fatalf("console output = %q, want a plain line naming WARN and the message", out)
	}
}

func TestNewLogger_LevelFiltersAndNothingIsSampled(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(logSettings{level: zapcore.WarnLevel, format: logFormatJSON}, zapcore.AddSync(&buf))
	l.Info("hidden")
	l.Debug("hidden too")
	l.Error("shown")
	if got := strings.Count(buf.String(), "\n"); got != 1 || !strings.Contains(buf.String(), `"shown"`) {
		t.Fatalf("at warn the info and debug lines must be dropped and the error kept, got:\n%s", buf.String())
	}

	// zap.NewProduction samples: the first 100 identical messages in a second, then
	// every 100th. The standard logger never dropped a line, and every one of the
	// engine's now goes through this logger.
	buf.Reset()
	l = newLogger(logSettings{level: zapcore.InfoLevel, format: logFormatJSON}, zapcore.AddSync(&buf))
	const n = 1000
	for i := 0; i < n; i++ {
		l.Info("the same line")
	}
	if got := strings.Count(buf.String(), "\n"); got != n {
		t.Fatalf("logged %d of %d identical lines; a sampler is dropping them", got, n)
	}
}

// The reason this exists. A line the engine writes with log.Printf, which a viewer
// can only label unknown, comes out as JSON with a level, and the text of the
// message is the text the call site wrote.
func TestStdlibLinesCarryALevelInJSON(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(logSettings{level: zapcore.InfoLevel, format: logFormatJSON}, zapcore.AddSync(&buf))
	t.Cleanup(logbridge.Install(l.With(zap.String(relayPathLogField, "cosmos<->eth/cosmoshub-1"))))

	log.Printf("WARN [SendEthTx] nonce %d too low", 5)
	warn := lineOf(t, &buf)
	if warn["level"] != "warn" || warn["msg"] != "[SendEthTx] nonce 5 too low" {
		t.Fatalf("got level=%v msg=%v, want warn and the message without its token", warn["level"], warn["msg"])
	}
	if warn[relayPathLogField] != "cosmos<->eth/cosmoshub-1" {
		t.Errorf("%s = %v, want the relay path as a field", relayPathLogField, warn[relayPathLogField])
	}
	if c, _ := warn["caller"].(string); !strings.HasPrefix(c, "cmd/logging_test.go:") {
		t.Errorf("caller = %q, want the log statement in this file, not the bridge", c)
	}

	buf.Reset()
	log.Printf("[BatchBuilder] Flushing eth batch: %d packets", 1)
	info := lineOf(t, &buf)
	if info["level"] != "info" || info["msg"] != "[BatchBuilder] Flushing eth batch: 1 packets" {
		t.Fatalf("got level=%v msg=%v, want an unmarked line to be info and untouched", info["level"], info["msg"])
	}
}

// gnark prints three console lines on stdout for every proof, outside everything above.
// Through the bridge they are one more JSON line each: a level, the relay path, and gnark's
// own fields kept apart from the keys a viewer reads. At the default level the per-proof
// ones are not there at all.
func TestGnarkLinesCarryALevelInJSON(t *testing.T) {
	restoreGnarkLogger(t)

	var buf bytes.Buffer
	l := newLogger(logSettings{level: zapcore.DebugLevel, format: logFormatJSON}, zapcore.AddSync(&buf))
	t.Cleanup(logbridge.Install(l.With(zap.String(relayPathLogField, "cosmos<->eth/cosmoshub-1"))))

	g := gnarklogger.Logger()
	g.Debug().Str("curve", "bn254").Str("backend", "groth16").Int("nbConstraints", 878769).
		Dur("took", 3228*time.Millisecond).Msg("prover done")

	line := lineOf(t, &buf)
	if line["level"] != "debug" || line["msg"] != "[gnark] prover done" {
		t.Fatalf("got level=%v msg=%v, want debug and the gnark message", line["level"], line["msg"])
	}
	if line[relayPathLogField] != "cosmos<->eth/cosmoshub-1" || line["component"] != "gnark" {
		t.Errorf("relay_path=%v component=%v, want the relay path and gnark", line[relayPathLogField], line["component"])
	}
	if _, present := line["caller"]; present {
		t.Errorf("caller = %v: a gnark line has no call site to name", line["caller"])
	}
	inner, _ := line["gnark"].(map[string]any)
	if inner["curve"] != "bn254" || inner["backend"] != "groth16" ||
		inner["nbConstraints"] != float64(878769) || inner["took"] != float64(3228) {
		t.Errorf("gnark's fields = %v, want curve, backend, nbConstraints and took (ms) under \"gnark\"", line["gnark"])
	}

	// A real gnark message at a level that matters keeps it: the warning is what must not
	// be lost to the filtering that hides the per-proof lines.
	buf.Reset()
	g.Warn().Msg("circuit has unconstrained inputs")
	if warn := lineOf(t, &buf); warn["level"] != "warn" || warn["msg"] != "[gnark] circuit has unconstrained inputs" {
		t.Errorf("a gnark warning arrived as level=%v msg=%v", warn["level"], warn["msg"])
	}

	// At the default level the per-proof chatter is gone and the warning is not.
	buf.Reset()
	quiet := newLogger(logSettings{level: zapcore.InfoLevel, format: logFormatJSON}, zapcore.AddSync(&buf))
	t.Cleanup(logbridge.Install(quiet))
	g = gnarklogger.Logger()
	g.Debug().Msg("prover done")
	g.Warn().Msg("circuit has unconstrained inputs")
	if got := strings.Count(buf.String(), "\n"); got != 1 || !strings.Contains(buf.String(), "unconstrained") {
		t.Fatalf("at info level want only the warning, got:\n%s", buf.String())
	}
}
