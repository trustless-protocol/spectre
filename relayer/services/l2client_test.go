package services

import (
	"encoding/json"
	"testing"

	relayerclient "relayer/client"

	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	"github.com/ethereum/go-ethereum/common"
)

func testL2Params() L2ClientParams {
	return L2ClientParams{
		WasmChecksum:  "0xabcd",
		RollupProfile: json.RawMessage(`{"common":{"l2_router":"0x1111111111111111111111111111111111111111"},"dispute_game_factory":"0x22"}`),
		Bootstrap: relayerclient.L2BootstrapState{
			Height:            123,
			BlockHash:         common.HexToHash("0xcccc"),
			ParentHash:        common.HexToHash("0xdddd"),
			StateRoot:         common.HexToHash("0xaaaa"),
			RouterStorageRoot: common.HexToHash("0xbbbb"),
			TimestampSeconds:  1700000000,
		},
	}
}

func TestBuildL2WasmClientState_ClientStateShape(t *testing.T) {
	cs, consensus, err := BuildL2WasmClientState(testL2Params())
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	wasmCS, ok := cs.(*ibcwasmtypes.ClientState)
	if !ok {
		t.Fatalf("client state is not a wasm ClientState: %T", cs)
	}
	if wasmCS.LatestHeight.RevisionHeight != 123 {
		t.Fatalf("latest height = %d, want 123", wasmCS.LatestHeight.RevisionHeight)
	}

	// The client-state Data must be ClientState<Profile>: latest_height, frozen_height
	// (null), and the rollup profile embedded verbatim as `profile`.
	var clientState l2ClientStateJSON
	if err := json.Unmarshal(wasmCS.Data, &clientState); err != nil {
		t.Fatalf("client-state Data is not JSON: %v", err)
	}
	if clientState.LatestHeight != 123 {
		t.Fatalf("latest_height = %d, want 123", clientState.LatestHeight)
	}
	if clientState.FrozenHeight != nil {
		t.Fatalf("frozen_height should be null, got %v", *clientState.FrozenHeight)
	}
	if string(clientState.Profile) != string(testL2Params().RollupProfile) {
		t.Fatalf("profile not embedded verbatim: %s", clientState.Profile)
	}
	var clientStateFields map[string]json.RawMessage
	if err := json.Unmarshal(wasmCS.Data, &clientStateFields); err != nil {
		t.Fatalf("client-state Data is not a JSON object: %v", err)
	}
	if _, found := clientStateFields["state_version"]; found {
		t.Fatal("state_version must not be emitted for fresh clients")
	}
	if string(clientState.FinalityPolicy) != string(defaultL2FinalityPolicy) {
		t.Fatalf("unexpected default finality policy: %s", clientState.FinalityPolicy)
	}

	// Consensus from the bootstrap roots: ibc_storage_root = router_storage_root,
	// timestamp_nanos = timestamp_seconds * 1e9.
	wasmCons, ok := consensus.(*ibcwasmtypes.ConsensusState)
	if !ok {
		t.Fatalf("consensus is not a wasm ConsensusState: %T", consensus)
	}
	var cons l2ConsensusStateJSON
	if err := json.Unmarshal(wasmCons.Data, &cons); err != nil {
		t.Fatalf("consensus Data is not JSON: %v", err)
	}
	if cons.TimestampNanos != 1700000000*1_000_000_000 {
		t.Fatalf("timestamp_nanos = %d, want seconds*1e9", cons.TimestampNanos)
	}
	if cons.IBCStorageRoot != common.HexToHash("0xbbbb").Hex() {
		t.Fatalf("ibc_storage_root = %s, want router storage root", cons.IBCStorageRoot)
	}
	if cons.StateRoot != common.HexToHash("0xaaaa").Hex() {
		t.Fatalf("state_root = %s, want bootstrap state root", cons.StateRoot)
	}
	if cons.L2Height != 123 || cons.L2BlockHash != common.HexToHash("0xcccc").Hex() || cons.ParentHash != common.HexToHash("0xdddd").Hex() {
		t.Fatalf("bootstrap block identity missing: %+v", cons)
	}
	if cons.FinalityLevel != "unsafe" || cons.ProposalStatus != "pending" {
		t.Fatalf("bootstrap finality = %s/%s", cons.FinalityLevel, cons.ProposalStatus)
	}
}

func TestBuildL2WasmClientState_RequiresRollupProfile(t *testing.T) {
	p := testL2Params()
	p.RollupProfile = nil
	if _, _, err := BuildL2WasmClientState(p); err == nil {
		t.Fatal("missing rollup profile must error")
	}
}

func TestBuildL2WasmClientState_PreservesConfiguredPoliciesAndRejectsMalformedOnes(t *testing.T) {
	p := testL2Params()
	p.FinalityPolicy = json.RawMessage(`{"minimum_update_level":"finalized"}`)
	p.FreshnessPolicy = json.RawMessage(`{"max_time_without_finalized_update":60}`)
	cs, _, err := BuildL2WasmClientState(p)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	var state l2ClientStateJSON
	if err := json.Unmarshal(cs.(*ibcwasmtypes.ClientState).Data, &state); err != nil {
		t.Fatalf("decode client state: %v", err)
	}
	if string(state.FinalityPolicy) != string(p.FinalityPolicy) || string(state.FreshnessPolicy) != string(p.FreshnessPolicy) {
		t.Fatalf("policies changed: finality=%s freshness=%s", state.FinalityPolicy, state.FreshnessPolicy)
	}

	p.FinalityPolicy = json.RawMessage(`[]`)
	if _, _, err := BuildL2WasmClientState(p); err == nil {
		t.Fatal("array finality policy must be rejected")
	}
}
