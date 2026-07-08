package services

import (
	"time"

	relayerclient "relayer/client"
)

const DEFAULT_INTERVAL = 5
const DEFAULT_TIME_PERIODS = time.Second * 20 //default 20 seconds

// DEFAULT_CLOCK_DRIFT mirrors relayerclient.DefaultClockDrift (the canonical
// source) so config defaults and genesis generation stay in lock-step.
const DEFAULT_CLOCK_DRIFT uint32 = relayerclient.DefaultClockDrift

// DEFAULT_BEACON_FINALITY_RETRIES caps how many 10s polls the ETH→Cosmos relay
// waits for beacon finality to cover an event block before re-queueing. It must
// exceed Ethereum worst-case finality (~2 epochs ≈ 12.8 min); 90 × 10s = 15 min
// leaves headroom so valid packets are not spuriously re-queued.
const DEFAULT_BEACON_FINALITY_RETRIES uint32 = 90

const DEFAULT_COSMOS_APP_HASH_WAIT_RETRIES uint32 = 30
const DEFAULT_COSMOS_APP_HASH_WAIT_INTERVAL = time.Second

type IntervalType uint8

const (
	blockHeight IntervalType = iota
	timestamp
)

type IntervalConfig struct {
	blockHeight uint8
	blockTime   time.Duration
}

type BatchConfig struct {
	BatchSize    uint8
	BatchPeriods time.Duration
}

type Config struct {
	KeyPath               string
	IntervalParams        IntervalConfig
	IntervalType          IntervalType
	BatchConfig           BatchConfig
	TrustingPeriod        uint32
	TrustLevel            string
	ProofType             string
	ClockDrift            uint32
	BeaconFinalityRetries uint32
	AppHashWaitRetries    uint32
	AppHashWaitInterval   time.Duration
	FetchTimeout          time.Duration
}

func NewConfig(KeyPath string, params IntervalConfig, intervalType IntervalType) Config {
	return Config{
		IntervalParams: params,
		IntervalType:   intervalType,
		FetchTimeout:   time.Second * 15,
	}
}

func DefaultConfig() Config {
	return Config{
		IntervalParams: IntervalConfig{
			blockHeight: DEFAULT_INTERVAL,
			blockTime:   DEFAULT_TIME_PERIODS,
		},
		IntervalType: timestamp,
		BatchConfig: BatchConfig{
			BatchPeriods: time.Second * 3, // default each batch waits for 3 seconds
			BatchSize:    5,               // default 5 packets per batch (multicall gas budget: ~5×2M wasm verify ≈ 10M, fits prod block limit)
		},
		TrustLevel:            "2/3",
		ProofType:             "groth16",
		ClockDrift:            DEFAULT_CLOCK_DRIFT,
		BeaconFinalityRetries: DEFAULT_BEACON_FINALITY_RETRIES,
		AppHashWaitRetries:    DEFAULT_COSMOS_APP_HASH_WAIT_RETRIES,
		AppHashWaitInterval:   DEFAULT_COSMOS_APP_HASH_WAIT_INTERVAL,
		FetchTimeout:          time.Second * 15,
	}
}
