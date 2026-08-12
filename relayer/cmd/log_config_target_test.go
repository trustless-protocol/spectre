package main

import (
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// TestLogConfigTarget checks the line that tells an operator which file a
// command actually read. --config defaults to "config.json", so a run meant for
// another file and invoked without the flag reads that default instead; the
// resolved absolute path is what makes the two distinguishable in the log.
func TestLogConfigTarget(t *testing.T) {
	t.Parallel()

	capture := func(configPath string, cfg *appConfig) string {
		core, logs := observer.New(zapcore.InfoLevel)
		logConfigTarget(zap.New(core), "create-clients-eth", configPath, cfg)
		if logs.Len() != 1 {
			t.Fatalf("logged %d lines, want 1", logs.Len())
		}
		return logs.All()[0].Message
	}

	t.Run("names the resolved path and the selected source", func(t *testing.T) {
		t.Parallel()
		cfg := &appConfig{}
		cfg.CosmosToEthConfig.ICS26ClientID = "arb-client-0"

		msg := capture("config-arb.json", cfg)

		want, err := filepath.Abs("config-arb.json")
		if err != nil {
			t.Fatalf("abs: %v", err)
		}
		if !strings.Contains(msg, want) {
			t.Errorf("message %q does not carry the absolute path %q", msg, want)
		}
		if !strings.Contains(msg, "arb-client-0") {
			t.Errorf("message %q does not name the selected source", msg)
		}
	})

	// A relative path alone would not have separated config.json from
	// config-arb.json in the report that prompted this; the absolute form does.
	t.Run("distinguishes two configs in the same directory", func(t *testing.T) {
		t.Parallel()
		cfg := &appConfig{}
		cfg.CosmosToEthConfig.ICS26ClientID = "arb-client-0"

		if a, b := capture("config.json", cfg), capture("config-arb.json", cfg); a == b {
			t.Fatalf("both configs logged the same line: %q", a)
		}
	})

	t.Run("omits the source when none is selected", func(t *testing.T) {
		t.Parallel()
		msg := capture("config.json", &appConfig{})
		if strings.Contains(msg, "source") {
			t.Errorf("message %q mentions a source when none was selected", msg)
		}
	})
}
