package subscriber

import (
	"math/big"
	"testing"

	contractICS26Router "relayer/bindings/ICS26Router"
)

func TestEthPacketToCosmosPacket_SinglePayload(t *testing.T) {
	ethPacket := contractICS26Router.IICS26RouterMsgsPacket{
		SourceClient:     "eth-client-0",
		DestClient:       "cosmos-client-0",
		TimeoutTimestamp: 1700000000,
		Payloads: []contractICS26Router.IICS26RouterMsgsPayload{
			{
				SourcePort: "transfer",
				DestPort:   "transfer",
				Version:    "ics20-1",
				Encoding:   "proto3",
				Value:      []byte("test-value"),
			},
		},
	}

	seq := big.NewInt(42)
	result := EthPacketToCosmosPacket(ethPacket, seq)

	if result.Sequence != 42 {
		t.Errorf("Sequence: got %d, want 42", result.Sequence)
	}
	if result.SourceClient != "eth-client-0" {
		t.Errorf("SourceClient: got %q, want %q", result.SourceClient, "eth-client-0")
	}
	if result.DestinationClient != "cosmos-client-0" {
		t.Errorf("DestinationClient: got %q, want %q", result.DestinationClient, "cosmos-client-0")
	}
	if result.TimeoutTimestamp != 1700000000 {
		t.Errorf("TimeoutTimestamp: got %d, want 1700000000", result.TimeoutTimestamp)
	}
	if len(result.Payloads) != 1 {
		t.Fatalf("Payloads length: got %d, want 1", len(result.Payloads))
	}

	p := result.Payloads[0]
	if p.SourcePort != "transfer" {
		t.Errorf("SourcePort: got %q", p.SourcePort)
	}
	if p.DestinationPort != "transfer" {
		t.Errorf("DestinationPort: got %q", p.DestinationPort)
	}
	if p.Version != "ics20-1" {
		t.Errorf("Version: got %q", p.Version)
	}
	if p.Encoding != "proto3" {
		t.Errorf("Encoding: got %q", p.Encoding)
	}
	if string(p.Value) != "test-value" {
		t.Errorf("Value: got %q", p.Value)
	}
}

func TestEthPacketToCosmosPacket_MultiplePayloads(t *testing.T) {
	ethPacket := contractICS26Router.IICS26RouterMsgsPacket{
		SourceClient:     "src",
		DestClient:       "dst",
		TimeoutTimestamp: 999,
		Payloads: []contractICS26Router.IICS26RouterMsgsPayload{
			{SourcePort: "port1", DestPort: "port1d", Version: "v1", Encoding: "e1", Value: []byte("val1")},
			{SourcePort: "port2", DestPort: "port2d", Version: "v2", Encoding: "e2", Value: []byte("val2")},
			{SourcePort: "port3", DestPort: "port3d", Version: "v3", Encoding: "e3", Value: []byte("val3")},
		},
	}

	result := EthPacketToCosmosPacket(ethPacket, big.NewInt(1))

	if len(result.Payloads) != 3 {
		t.Fatalf("Payloads length: got %d, want 3", len(result.Payloads))
	}
	for i, p := range result.Payloads {
		if p.SourcePort == "" {
			t.Errorf("payload[%d] SourcePort is empty", i)
		}
		if p.DestinationPort == "" {
			t.Errorf("payload[%d] DestinationPort is empty", i)
		}
	}
}

func TestEthPacketToCosmosPacket_EmptyPayloads(t *testing.T) {
	ethPacket := contractICS26Router.IICS26RouterMsgsPacket{
		SourceClient:     "src",
		DestClient:       "dst",
		TimeoutTimestamp: 0,
		Payloads:         []contractICS26Router.IICS26RouterMsgsPayload{},
	}

	result := EthPacketToCosmosPacket(ethPacket, big.NewInt(0))

	if len(result.Payloads) != 0 {
		t.Errorf("expected 0 payloads, got %d", len(result.Payloads))
	}
}

func TestEthPacketToCosmosPacket_LargeSequence(t *testing.T) {
	ethPacket := contractICS26Router.IICS26RouterMsgsPacket{
		SourceClient: "src",
		DestClient:   "dst",
	}

	seq := new(big.Int).SetUint64(^uint64(0)) // max uint64
	result := EthPacketToCosmosPacket(ethPacket, seq)

	if result.Sequence != ^uint64(0) {
		t.Errorf("Sequence: got %d, want max uint64", result.Sequence)
	}
}

func TestEthStartupRecoveryLookbackBlocksFromEnv(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		raw  string
		want uint64
	}{
		{name: "empty uses default", raw: "", want: defaultEthStartupRecoveryLookbackBlocks},
		{name: "invalid uses default", raw: "not-a-number", want: defaultEthStartupRecoveryLookbackBlocks},
		{name: "zero allowed", raw: "0", want: 0},
		{name: "custom value", raw: "1024", want: 1024},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := ethStartupRecoveryLookbackBlocksFromEnv(tc.raw)
			if got != tc.want {
				t.Fatalf("ethStartupRecoveryLookbackBlocksFromEnv(%q) = %d, want %d", tc.raw, got, tc.want)
			}
		})
	}
}

func TestEthStartupRecoveryStartBlock(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		latest    uint64
		lookback  uint64
		wantStart uint64
	}{
		{name: "lookback smaller than latest", latest: 1000, lookback: 256, wantStart: 744},
		{name: "lookback equals latest", latest: 256, lookback: 256, wantStart: 0},
		{name: "lookback greater than latest", latest: 42, lookback: 100, wantStart: 0},
		{name: "zero lookback", latest: 42, lookback: 0, wantStart: 42},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := ethStartupRecoveryStartBlock(tc.latest, tc.lookback)
			if got != tc.wantStart {
				t.Fatalf("ethStartupRecoveryStartBlock(%d, %d) = %d, want %d",
					tc.latest, tc.lookback, got, tc.wantStart)
			}
		})
	}
}
