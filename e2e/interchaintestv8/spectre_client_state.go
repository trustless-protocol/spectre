package main

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/decentrio/fast-ibc/packages/go-abigen/spectreclient"
)

// The SpectreClient contract exposes its client state only as
// ABI-encoded bytes via getClientState(); the struct getter was made private in
// the contract refactor, so there is no generated ClientState() binding method.
// These helpers decode those bytes into a struct so the e2e suites can assert on
// individual fields, mirroring relayer/client.DecodeClientState.

// spectreTrustLevel mirrors SpectreMsgs.TrustThreshold.
type spectreTrustLevel struct {
	Numerator   uint8
	Denominator uint8
}

// spectreClientState mirrors the SpectreMsgs.ClientState tuple returned
// (ABI-encoded) by SpectreClient.getClientState().
type spectreClientState struct {
	ChainId         string
	TrustLevel      spectreTrustLevel
	LatestHeight    spectreclient.ICS02ClientMsgsHeight
	TrustingPeriod  uint32
	UnbondingPeriod uint32
	IsFrozen        bool
	ZkAlgorithm     uint8
	ClockDrift      uint32
}

var spectreClientStateABIType abi.Type

func init() {
	var err error
	spectreClientStateABIType, err = abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "chainId", Type: "string"},
		{Name: "trustLevel", Type: "tuple", Components: []abi.ArgumentMarshaling{
			{Name: "numerator", Type: "uint8"},
			{Name: "denominator", Type: "uint8"},
		}},
		{Name: "latestHeight", Type: "tuple", Components: []abi.ArgumentMarshaling{
			{Name: "revisionNumber", Type: "uint64"},
			{Name: "revisionHeight", Type: "uint64"},
		}},
		{Name: "trustingPeriod", Type: "uint32"},
		{Name: "unbondingPeriod", Type: "uint32"},
		{Name: "isFrozen", Type: "bool"},
		{Name: "zkAlgorithm", Type: "uint8"},
		{Name: "clockDrift", Type: "uint32"},
	})
	if err != nil {
		panic(fmt.Sprintf("build Spectre client state ABI type: %v", err))
	}
}

// decodeSpectreClientState decodes the ABI-encoded bytes returned by
// SpectreClient.getClientState() into a struct.
func decodeSpectreClientState(data []byte) (spectreClientState, error) {
	args := abi.Arguments{{Type: spectreClientStateABIType}}
	unpacked, err := args.Unpack(data)
	if err != nil {
		return spectreClientState{}, fmt.Errorf("unpack client state: %w", err)
	}
	if len(unpacked) == 0 {
		return spectreClientState{}, fmt.Errorf("no client state data unpacked")
	}
	// unpacked[0] is an anonymous struct matching the tuple; round-trip through
	// JSON to populate the named target type by matching field names.
	jsonBytes, err := json.Marshal(unpacked[0])
	if err != nil {
		return spectreClientState{}, fmt.Errorf("marshal unpacked tuple: %w", err)
	}
	var cs spectreClientState
	if err := json.Unmarshal(jsonBytes, &cs); err != nil {
		return spectreClientState{}, fmt.Errorf("unmarshal client state: %w", err)
	}
	return cs, nil
}

// getSpectreClientState fetches and decodes the on-chain client state from a
// SpectreClient contract instance.
func getSpectreClientState(c *spectreclient.Contract) (spectreClientState, error) {
	bz, err := c.GetClientState(nil)
	if err != nil {
		return spectreClientState{}, err
	}
	return decodeSpectreClientState(bz)
}
