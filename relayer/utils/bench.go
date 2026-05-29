package utils

import (
	"os"
	"strings"
	"sync/atomic"
)

// Benchmark mode toggles per-call instrumentation that's useful for performance
// measurement but noisy in production. Two flags:
//
//   - RELAYER_BENCHMARK=1 → enables [bench][prover], [bench][eth], [bench][cosmos]
//     summary lines per tx/proof.
//   - RELAYER_BENCH_INNER_GAS=1 → enables the per-inner-call gas estimate inside
//     SendEthTxBatch (one eth_estimateGas RPC per inner). Defaults to the parent
//     RELAYER_BENCHMARK flag when unset.
const (
	BenchEnvKey         = "RELAYER_BENCHMARK"
	BenchInnerGasEnvKey = "RELAYER_BENCH_INNER_GAS"
)

var benchOverride atomic.Bool

// SetBenchEnabled lets tests force benchmark mode on/off without touching env.
func SetBenchEnabled(enabled bool) {
	benchOverride.Store(enabled)
}

// BenchEnabled reports whether benchmark-mode logging is on.
func BenchEnabled() bool {
	return benchOverride.Load() || envBool(BenchEnvKey)
}

// BenchInnerGasEnabled reports whether per-inner-call gas estimation in
// SendEthTxBatch is on. Falls back to BenchEnabled() when the dedicated env
// var isn't set.
func BenchInnerGasEnabled() bool {
	if raw, ok := os.LookupEnv(BenchInnerGasEnvKey); ok {
		return parseBenchBool(raw)
	}
	return BenchEnabled()
}

func envBool(key string) bool {
	raw, ok := os.LookupEnv(key)
	return ok && parseBenchBool(raw)
}

func parseBenchBool(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "t", "yes", "y", "on", "enable", "enabled":
		return true
	default:
		return false
	}
}
