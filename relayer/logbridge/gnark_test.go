package logbridge

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	gnarklogger "github.com/consensys/gnark/logger"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// gnarkLog is gnark's global logger as its own code takes it, into a variable: zerolog's
// event methods have pointer receivers and Logger returns a value.
func gnarkLog() *zerolog.Logger {
	l := gnarklogger.Logger()
	return &l
}

// num reads a numeric field however it was typed: zerolog writes a whole number of
// milliseconds without a decimal point, which comes back as an int, and a fractional
// one as a float.
func num(t *testing.T, v any) float64 {
	t.Helper()
	switch x := v.(type) {
	case int64:
		return float64(x)
	case float64:
		return x
	}
	t.Fatalf("%v (%T) is not a number", v, v)
	return 0
}

func TestZapLevel(t *testing.T) {
	cases := []struct {
		in   zerolog.Level
		want zapcore.Level
	}{
		{zerolog.TraceLevel, zapcore.DebugLevel},
		{zerolog.DebugLevel, zapcore.DebugLevel},
		{zerolog.InfoLevel, zapcore.InfoLevel},
		{zerolog.NoLevel, zapcore.InfoLevel},
		{zerolog.WarnLevel, zapcore.WarnLevel},
		{zerolog.ErrorLevel, zapcore.ErrorLevel},
		// zerolog exits or panics by itself after writing these. They are errors here and
		// nothing more; a writer that also exited would decide it twice.
		{zerolog.FatalLevel, zapcore.ErrorLevel},
		{zerolog.PanicLevel, zapcore.ErrorLevel},
	}
	for _, tc := range cases {
		if got := zapLevel(tc.in); got != tc.want {
			t.Errorf("zapLevel(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

// Each level has to come out as itself. Fatal and panic are written straight to the
// writer: sent through a zerolog event they would end the test binary.
func TestGnarkWriter_EveryLevelKeepsItsMessage(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	w := gnarkWriter{logger: zap.New(core)}

	for _, level := range []zerolog.Level{
		zerolog.TraceLevel, zerolog.DebugLevel, zerolog.InfoLevel, zerolog.NoLevel,
		zerolog.WarnLevel, zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel,
	} {
		if _, err := w.WriteLevel(level, []byte(`{"level":"x","message":"m `+level.String()+`"}`+"\n")); err != nil {
			t.Fatal(err)
		}
	}
	got := logs.All()
	want := []zapcore.Level{
		zapcore.DebugLevel, zapcore.DebugLevel, zapcore.InfoLevel, zapcore.InfoLevel,
		zapcore.WarnLevel, zapcore.ErrorLevel, zapcore.ErrorLevel, zapcore.ErrorLevel,
	}
	if len(got) != len(want) {
		t.Fatalf("got %d entries, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Level != want[i] {
			t.Errorf("entry %d (%q) at %v, want %v", i, got[i].Message, got[i].Level, want[i])
		}
		if len(got[i].Message) < len("[gnark] m ") || got[i].Message[:len("[gnark] m ")] != "[gnark] m " {
			t.Errorf("entry %d message = %q, want it to start with the [gnark] tag", i, got[i].Message)
		}
	}
}

// What a viewer needs from the three lines gnark prints per proof: the message, and the
// numbers still numbers so they can be compared and charted.
func TestGnarkWriter_FieldsAreTypedAndNamespaced(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	zl := zerolog.New(gnarkWriter{logger: zap.New(core).With(zap.String("component", "gnark"))})

	zl.Debug().Str("curve", "bn254").Int("nbConstraints", 878769).
		Dur("took", 3228*time.Millisecond).Bool("accelerated", false).Msg("prover done")

	got := logs.All()
	if len(got) != 1 {
		t.Fatalf("got %d entries, want 1", len(got))
	}
	e := got[0]
	if e.Level != zapcore.DebugLevel || e.Message != "[gnark] prover done" {
		t.Fatalf("got (%v, %q), want (debug, \"[gnark] prover done\")", e.Level, e.Message)
	}
	ctx := e.ContextMap()
	if ctx["component"] != "gnark" {
		t.Errorf("component = %v, want gnark", ctx["component"])
	}
	fields, ok := ctx["gnark"].(map[string]any)
	if !ok {
		t.Fatalf("gnark's own fields are not under a \"gnark\" key: %v", ctx)
	}
	if fields["curve"] != "bn254" || fields["accelerated"] != false {
		t.Errorf("curve/accelerated = %v/%v", fields["curve"], fields["accelerated"])
	}
	if n, isInt := fields["nbConstraints"].(int64); !isInt || n != 878769 {
		t.Errorf("nbConstraints = %v (%T), want the integer 878769", fields["nbConstraints"], fields["nbConstraints"])
	}
	if took := num(t, fields["took"]); took != 3228 {
		t.Errorf("took = %v, want 3228 (zerolog writes a duration in milliseconds)", took)
	}
	// zerolog's own keys carry nothing a viewer needs: the level is the entry's, and zap
	// stamps the time.
	for _, k := range []string{"message", "level", "time"} {
		if _, present := fields[k]; present {
			t.Errorf("zerolog's %q key leaked into gnark's fields: %v", k, fields)
		}
	}
}

// The reason fields are namespaced. Written at the top level, a gnark field named like
// one of this logger's own keys would put the key in the line twice, and a viewer reads the
// last one: a debug line would arrive as whatever level gnark's field said.
func TestGnarkWriter_AFieldCannotShadowTheLine(t *testing.T) {
	var buf bytes.Buffer
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	l := zap.New(zapcore.NewCore(enc, zapcore.AddSync(&buf), zapcore.DebugLevel), zap.AddCaller())
	t.Cleanup(Install(l))

	gnarkLog().Debug().
		Str("msg", "evil").Str("ts", "tomorrow").Str("caller", "elsewhere").Str("stacktrace", "x").
		Msg("hello")

	var line map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &line); err != nil {
		t.Fatalf("not one JSON line: %v\n%s", err, buf.String())
	}
	if line["level"] != "debug" || line["msg"] != "[gnark] hello" {
		t.Fatalf("level=%v msg=%v, want debug and the gnark message: a field overwrote them", line["level"], line["msg"])
	}
	if _, isNumber := line["ts"].(float64); !isNumber {
		t.Errorf("ts = %v (%T), want zap's timestamp", line["ts"], line["ts"])
	}
	if _, present := line["caller"]; present {
		t.Errorf("caller = %v on a gnark line: it would name gnark.go, not gnark's call site", line["caller"])
	}
	inner, _ := line["gnark"].(map[string]any)
	if inner["msg"] != "evil" || inner["ts"] != "tomorrow" || inner["caller"] != "elsewhere" {
		t.Errorf("gnark's own msg/ts/caller fields are not preserved under \"gnark\": %v", line["gnark"])
	}
}

func TestGnarkWriter_UnparsableLineIsReportedNotDropped(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	w := gnarkWriter{logger: zap.New(core)}
	if n, err := w.WriteLevel(zerolog.ErrorLevel, []byte("this is not json\n")); err != nil || n != len("this is not json\n") {
		t.Fatalf("WriteLevel = (%d, %v)", n, err)
	}
	got := logs.All()
	if len(got) != 1 || got[0].Level != zapcore.WarnLevel || got[0].Message != "[gnark] unparsable log line" {
		t.Fatalf("got %+v, want one warn naming the unparsable line", got)
	}
	if raw := got[0].ContextMap()["raw"]; raw != "this is not json" {
		t.Errorf("raw = %v, want the line itself so it can still be read", raw)
	}
}

// At the default level the three per-proof lines are all debug, so nothing should be spent
// parsing them. A line that would be dropped is not even looked at.
func TestGnarkWriter_ADisabledLevelIsNotParsed(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	w := gnarkWriter{logger: zap.New(core)}
	if _, err := w.WriteLevel(zerolog.DebugLevel, []byte("not json at all")); err != nil {
		t.Fatal(err)
	}
	if logs.Len() != 0 {
		t.Fatalf("a debug line at info level produced %+v", logs.All())
	}
}

func TestInstall_RoutesGnarksLogger(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	restore := Install(zap.New(core, zap.AddCaller()).With(zap.String("relay_path", "cosmos<->eth/cosmoshub-1")))

	gnarkLog().Warn().Str("circuit", "batch").Msg("circuit has unconstrained inputs")

	got := logs.All()
	if len(got) != 1 {
		t.Fatalf("got %d entries, want 1: gnark's logger is not routed into the installed one", len(got))
	}
	e := got[0]
	if e.Level != zapcore.WarnLevel || e.Message != "[gnark] circuit has unconstrained inputs" {
		t.Errorf("got (%v, %q)", e.Level, e.Message)
	}
	ctx := e.ContextMap()
	// The logger's own fields reach gnark's lines: that is how a relay path is attributed.
	if ctx["relay_path"] != "cosmos<->eth/cosmoshub-1" || ctx["component"] != "gnark" {
		t.Errorf("fields = %v, want relay_path and component", ctx)
	}
	if e.Caller.Defined {
		t.Errorf("caller = %s on a gnark line, want none", e.Caller.String())
	}

	// And it is undone, which is what keeps one test's logger out of the next one's.
	restore()
	gnarkLog().Warn().Msg("after restore")
	if logs.Len() != 1 {
		t.Fatalf("a line logged after restore still reached the installed logger: %+v", logs.All())
	}
}

// main installs the bridge, then `start` installs it again on the path-scoped logger. The
// second must win, and undoing it must give back the first, not the original.
func TestInstall_GnarksLoggerFollowsTheLatestInstall(t *testing.T) {
	coreA, logsA := observer.New(zapcore.DebugLevel)
	coreB, logsB := observer.New(zapcore.DebugLevel)

	restoreA := Install(zap.New(coreA))
	restoreB := Install(zap.New(coreB))

	gnarkLog().Info().Msg("to b")
	if logsA.Len() != 0 || logsB.Len() != 1 {
		t.Fatalf("after two installs: a=%d b=%d entries, want 0 and 1", logsA.Len(), logsB.Len())
	}

	restoreB()
	gnarkLog().Info().Msg("to a")
	if logsA.Len() != 1 || logsB.Len() != 1 {
		t.Fatalf("after undoing the second: a=%d b=%d entries, want 1 and 1", logsA.Len(), logsB.Len())
	}

	restoreA()
	gnarkLog().Info().Msg("to nobody")
	if logsA.Len() != 1 || logsB.Len() != 1 {
		t.Fatalf("after undoing both: a=%d b=%d entries, want 1 and 1", logsA.Len(), logsB.Len())
	}
}

// Provers log from their own goroutines. zerolog needs a writer that is safe to share,
// and this one holds no state. Run with -race.
func TestInstall_GnarkConcurrentWriters(t *testing.T) {
	logs := installObserved(t, zapcore.DebugLevel)

	const writers, each = 8, 100
	var wg sync.WaitGroup
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < each; i++ {
				gnarkLog().Debug().Int("writer", w).Int("i", i).Msg("prover done")
			}
		}(w)
	}
	wg.Wait()

	if got := logs.Len(); got != writers*each {
		t.Fatalf("got %d entries, want %d", got, writers*each)
	}
}

type tinyCircuit struct {
	X, Y frontend.Variable `gnark:",secret"`
	Z    frontend.Variable `gnark:",public"`
}

func (c *tinyCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(api.Mul(c.X, c.Y), c.Z)
	return nil
}

// proveOnce runs gnark for real: compile, setup, prove, verify on a two-constraint circuit.
func proveOnce(t *testing.T) {
	t.Helper()
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &tinyCircuit{})
	if err != nil {
		t.Fatal(err)
	}
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatal(err)
	}
	full, err := frontend.NewWitness(&tinyCircuit{X: 3, Y: 5, Z: 15}, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
	public, err := full.Public()
	if err != nil {
		t.Fatal(err)
	}
	proof, err := groth16.Prove(ccs, pk, full)
	if err != nil {
		t.Fatal(err)
	}
	if err := groth16.Verify(proof, vk, public); err != nil {
		t.Fatal(err)
	}
}

// The thing this exists for, against the real library and not a stand-in for it: the three lines
// gnark prints for every proof arrive as debug entries through zap, and at the default level
// they are not there at all.
func TestInstall_ARealProofLogsThroughTheBridge(t *testing.T) {
	perProof := []string{
		"[gnark] constraint system solver done",
		"[gnark] prover done",
		"[gnark] verifier done",
	}

	debugLogs := installObserved(t, zapcore.DebugLevel)
	proveOnce(t)
	seen := map[string]zapcore.Level{}
	for _, e := range debugLogs.All() {
		seen[e.Message] = e.Level
	}
	for _, msg := range perProof {
		level, ok := seen[msg]
		if !ok {
			t.Errorf("gnark's %q never reached the bridge; saw %v", msg, seen)
		} else if level != zapcore.DebugLevel {
			t.Errorf("%q arrived at %v, want debug", msg, level)
		}
	}

	infoLogs := installObserved(t, zapcore.InfoLevel)
	proveOnce(t)
	for _, e := range infoLogs.All() {
		for _, msg := range perProof {
			if e.Message == msg {
				t.Errorf("%q is visible at the default level; it is per-proof chatter", msg)
			}
		}
	}
	if infoLogs.Len() == 0 {
		t.Error("nothing from gnark at info level, but compiling a circuit logs at info: the router drops too much")
	}
}
