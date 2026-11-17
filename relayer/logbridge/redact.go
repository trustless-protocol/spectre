package logbridge

import (
	"attestor/types/redact"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RedactURLs wraps a core so every URL in a message, a string field or an error
// field is reduced to its scheme and host before it is written (#575).
//
// RPC providers put the API key in the URL path or query, and URLs reach the log
// from config lines, dial lines and, most often, the text of an RPC error. Each
// call site can forget; the sink cannot. A URL with nothing secret in it loses
// only its path.
func RedactURLs(core zapcore.Core) zapcore.Core { return redactingCore{Core: core} }

type redactingCore struct{ zapcore.Core }

func (c redactingCore) With(fields []zapcore.Field) zapcore.Core {
	return redactingCore{Core: c.Core.With(redactFields(fields))}
}

// Check adds this core, not the wrapped one, so Write below is what runs.
func (c redactingCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

func (c redactingCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	entry.Message = redact.URLs(entry.Message)
	return c.Core.Write(entry, redactFields(fields))
}

func redactFields(fields []zapcore.Field) []zapcore.Field {
	var out []zapcore.Field
	for i, field := range fields {
		redacted, changed := redactField(field)
		if !changed {
			continue
		}
		if out == nil {
			out = append([]zapcore.Field(nil), fields...)
		}
		out[i] = redacted
	}
	if out == nil {
		return fields
	}
	return out
}

func redactField(field zapcore.Field) (zapcore.Field, bool) {
	switch field.Type {
	case zapcore.StringType:
		if redacted := redact.URLs(field.String); redacted != field.String {
			field.String = redacted
			return field, true
		}
	case zapcore.ErrorType:
		if err, ok := field.Interface.(error); ok && err != nil {
			return zap.NamedError(field.Key, redact.Error(err)), true
		}
	}
	return field, false
}

// redactedLine is what the standard logger's lines are written as.
func redactedLine(msg string) string { return redact.URLs(msg) }
