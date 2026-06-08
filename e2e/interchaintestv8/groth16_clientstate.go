package main

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/decentrio/fast-ibc/packages/go-abigen/groth16ics07tendermint"
)

// The Groth16ICS07Tendermint contract exposes its client state only as
// ABI-encoded bytes via getClientState(); the struct getter was made private in
// the contract refactor, so there is no generated ClientState() binding method.
// These helpers decode those bytes into a struct so the e2e suites can assert on
// individual fields, mirroring relayer/client.DecodeClientState.

// groth16TrustLevel mirrors IICS07TendermintMsgs.TrustThreshold.
type groth16TrustLevel struct {
	Numerator   uint8
	Denominator uint8
}

// groth16ClientState mirrors the IICS07TendermintMsgs.ClientState tuple returned
// (ABI-encoded) by Groth16ICS07Tendermint.getClientState().
type groth16ClientState struct {
	ChainId         string
	TrustLevel      groth16TrustLevel
	LatestHeight    groth16ics07tendermint.IICS02ClientMsgsHeight
	TrustingPeriod  uint32
	UnbondingPeriod uint32
	IsFrozen        bool
	ZkAlgorithm     uint8
}

var groth16ClientStateABIType abi.Type

func init() {
	var err error
	groth16ClientStateABIType, err = abi.NewType("tuple", "", []abi.ArgumentMarshaling{
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
	})
	if err != nil {
		panic(fmt.Sprintf("build groth16 client state ABI type: %v", err))
	}
}

// decodeGroth16ClientState decodes the ABI-encoded bytes returned by
// Groth16ICS07Tendermint.getClientState() into a struct.
func decodeGroth16ClientState(data []byte) (groth16ClientState, error) {
	args := abi.Arguments{{Type: groth16ClientStateABIType}}
	unpacked, err := args.Unpack(data)
	if err != nil {
		return groth16ClientState{}, fmt.Errorf("unpack client state: %w", err)
	}
	if len(unpacked) == 0 {
		return groth16ClientState{}, fmt.Errorf("no client state data unpacked")
	}
	// unpacked[0] is an anonymous struct matching the tuple; round-trip through
	// JSON to populate the named target type by matching field names.
	jsonBytes, err := json.Marshal(unpacked[0])
	if err != nil {
		return groth16ClientState{}, fmt.Errorf("marshal unpacked tuple: %w", err)
	}
	var cs groth16ClientState
	if err := json.Unmarshal(jsonBytes, &cs); err != nil {
		return groth16ClientState{}, fmt.Errorf("unmarshal client state: %w", err)
	}
	return cs, nil
}

// getGroth16ClientState fetches and decodes the on-chain client state from a
// Groth16ICS07Tendermint contract instance.
func getGroth16ClientState(c *groth16ics07tendermint.Contract) (groth16ClientState, error) {
	bz, err := c.GetClientState(nil)
	if err != nil {
		return groth16ClientState{}, err
	}
	return decodeGroth16ClientState(bz)
}
