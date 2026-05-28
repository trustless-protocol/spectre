package benchmark

import (
	"os"
	"strings"
	"sync/atomic"
)

const (
	EnvKey         = "RELAYER_BENCHMARK"
	InnerGasEnvKey = "RELAYER_BENCH_INNER_GAS"
)

var enabledOverride atomic.Bool

func SetEnabled(enabled bool) {
	enabledOverride.Store(enabled)
}

func Enabled() bool {
	return enabledOverride.Load() || envBool(EnvKey)
}

func InnerGasEnabled() bool {
	if raw, ok := os.LookupEnv(InnerGasEnvKey); ok {
		return parseBool(raw)
	}
	return Enabled()
}

func envBool(key string) bool {
	raw, ok := os.LookupEnv(key)
	return ok && parseBool(raw)
}

func parseBool(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "t", "yes", "y", "on", "enable", "enabled":
		return true
	default:
		return false
	}
}
