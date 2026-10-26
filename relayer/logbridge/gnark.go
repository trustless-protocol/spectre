package logbridge

import (
	"bytes"
	"encoding/json"
	"sort"

	gnarklogger "github.com/consensys/gnark/logger"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// gnarkWriter is where gnark's global logger writes once Install has run.
//
// gnark logs through its own zerolog logger, which defaults to a colourised console
// writer on os.Stdout. Every proof prints three of its lines ("constraint system
// solver done", "prover done", "verifier done"), none of them JSON and none with a
// level a viewer can read, so they are what is left of the unlabelled stream once
// the standard logger is bridged. gnark exposes the logger on purpose
// (logger.Set), and every one of its call sites asks for it when it runs rather than
// capturing it at init, so replacing it here reaches all of them.
//
// It implements zerolog.LevelWriter, which is how zerolog passes the level apart
// from the JSON line: nothing has to be guessed from the text.
type gnarkWriter struct{ logger *zap.Logger }

// Write is only reached if something uses the logger itself as an io.Writer, which
// gnark does not; the level is unknown there, so it is info.
func (w gnarkWriter) Write(p []byte) (int, error) { return w.WriteLevel(zerolog.NoLevel, p) }

func (w gnarkWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	zl := zapLevel(level)
	// Skip the parse for a line the logger would drop anyway: the three per-proof
	// lines are debug, and at the default level that is all of them.
	if !w.logger.Core().Enabled(zl) {
		return len(p), nil
	}

	dec := json.NewDecoder(bytes.NewReader(p))
	dec.UseNumber() // keep integers exact: nbConstraints must not come back as a float
	var line map[string]any
	if err := dec.Decode(&line); err != nil {
		// zerolog writes JSON here, so this is not expected. Say so rather than drop
		// a line that may be an error.
		w.logger.Warn("[gnark] unparsable log line", zap.ByteString("raw", bytes.TrimSpace(p)))
		return len(p), nil
	}
	msg, _ := line["message"].(string)
	// zerolog's own keys. The level is already known and zap stamps the time.
	for _, k := range [...]string{"message", "level", "time"} {
		delete(line, k)
	}

	// Everything else is gnark's own fields (curve, nbConstraints, took, ...). They go
	// in a namespace, not at the top level: a field of gnark's named like one of this
	// logger's keys ("level", "msg", "ts") would otherwise write the key twice and a
	// viewer would keep the last, which is how a line ends up at the wrong level.
	keys := make([]string, 0, len(line))
	for k := range line {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fields := make([]zap.Field, 0, len(keys)+1)
	fields = append(fields, zap.Namespace("gnark"))
	for _, k := range keys {
		fields = append(fields, field(k, line[k]))
	}

	if ce := w.logger.Check(zl, "[gnark] "+msg); ce != nil {
		ce.Write(fields...)
	}
	return len(p), nil
}

// zapLevel maps zerolog's levels. Fatal and panic are errors here and nothing more:
// zerolog exits or panics by itself after writing them, and a log writer must not
// make that decision a second time.
func zapLevel(l zerolog.Level) zapcore.Level {
	switch l {
	case zerolog.TraceLevel, zerolog.DebugLevel:
		return zapcore.DebugLevel
	case zerolog.WarnLevel:
		return zapcore.WarnLevel
	case zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel:
		return zapcore.ErrorLevel
	default: // info, and NoLevel for a line logged without one
		return zapcore.InfoLevel
	}
}

// field types a decoded JSON value. Numbers arrive as json.Number so an integer stays
// an integer; zap would otherwise log a json.Number as a quoted string.
func field(key string, v any) zap.Field {
	switch x := v.(type) {
	case json.Number:
		if i, err := x.Int64(); err == nil {
			return zap.Int64(key, i)
		}
		if f, err := x.Float64(); err == nil {
			return zap.Float64(key, f)
		}
		return zap.String(key, x.String())
	case string:
		return zap.String(key, x)
	case bool:
		return zap.Bool(key, x)
	default:
		return zap.Any(key, x)
	}
}

// installGnark points gnark's global logger at logger and returns the function that
// puts the previous one back.
//
// gnark's Set writes a package variable with no lock, so this has to run before any
// goroutine that proves or verifies has started. The relayer calls Install from
// main, and again from stampRelayPath, which `start` reaches before it loads the
// prover or starts an engine.
func installGnark(logger *zap.Logger) (restore func()) {
	previous := gnarklogger.Logger()
	// No caller: the line is written from here, and gnark's own call site is not
	// something zerolog tells a writer, so "caller" would name this file for every
	// one of them.
	scoped := logger.WithOptions(zap.WithCaller(false)).With(zap.String("component", "gnark"))
	gnarklogger.Set(zerolog.New(gnarkWriter{logger: scoped}))
	return func() { gnarklogger.Set(previous) }
}
