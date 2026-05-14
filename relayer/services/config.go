package services

import "time"

const DEFAULT_INTERVAL = 5
const DEFAULT_TIME_PERIODS = time.Second * 20 //default 20 seconds

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
	KeyPath        string
	IntervalParams IntervalConfig
	IntervalType   IntervalType
	BatchConfig    BatchConfig
	TrustingPeriod uint32
	TrustLevel     string
	ProofType      string
}

func NewConfig(KeyPath string, params IntervalConfig, intervalType IntervalType) Config {
	return Config{
		IntervalParams: params,
		IntervalType:   intervalType,
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
			BatchPeriods: time.Second * 3, // default each batch waits for 10 seconds
			BatchSize:    10,              // default 10 packets per batch
		},
		TrustLevel: "2/3",
		ProofType:  "groth16",
	}
}
