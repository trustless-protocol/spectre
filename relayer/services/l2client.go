package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	relayerclient "relayer/client"

	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/ethereum/go-ethereum/common"
)

// decodeHexPrefixed decodes a 0x-optional hex string into bytes.
func decodeHexPrefixed(s string) ([]byte, error) {
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid hex %q: %w", s, err)
	}
	return b, nil
}

// L2ClientParams is the immutable configuration needed to bootstrap one L2 rollup
// wasm light client on Cosmos (Arbitrum / Base / Optimism). The bootstrap roots are
// read from the L2 chain (see relayerclient.GetL2BootstrapState); the rollup profile
// (which carries the L1 client id + checksum, l2 chain id, router address, commitment
// slot, and the rollup-specific verifier fields) is operator-supplied verbatim.
type L2ClientParams struct {
	// WasmChecksum is the hex checksum of the L2 client's own stored wasm code.
	WasmChecksum string
	// RollupProfile is the full ICS-08 L2 verifier Profile JSON (the `common` block
	// plus the rollup-specific finality fields), supplied verbatim from config and
	// embedded as ClientState.profile. Rollup-specific, so it is opaque here.
	RollupProfile json.RawMessage
	// FinalityPolicy is the optional packages/l2-client FinalityPolicy JSON.
	// Empty selects the contract defaults explicitly.
	FinalityPolicy json.RawMessage
	// FreshnessPolicy is the optional packages/l2-client FreshnessPolicy JSON.
	// Empty selects the contract defaults explicitly.
	FreshnessPolicy json.RawMessage
	// Bootstrap is the trusted initial state read from the L2 chain.
	Bootstrap relayerclient.L2BootstrapState
	// CounterpartyClientID is the L2-side client (on the rollup's ICS26Router) that
	// tracks Cosmos, registered inline as this client's counterparty. Empty when the
	// L2-side client id is not yet known, in which case registration is deferred.
	CounterpartyClientID string
}

// The JSON shapes below mirror the ICS-08 CosmWasm L2 client types
// (packages/l2-client/src/state.rs). serde uses deny_unknown_fields, so field names
// must match exactly. instantiate takes the client state and consensus state
// DIRECTLY (no Bootstrap envelope):
//   - client_state Data = the complete ClientState<Profile>, including
//     finality and freshness policy
//   - consensus_state Data = the complete ConsensusState, including L2 block
//     identity and conservative Unsafe/Pending bootstrap metadata
//
// The B256 fields serialize as 0x-prefixed lowercase hex (alloy default);
// timestamp_nanos = timestamp_seconds * 1e9 (matches Header::consensus_state).

type l2ClientStateJSON struct {
	LatestHeight          uint64          `json:"latest_height"`
	FrozenHeight          *uint64         `json:"frozen_height"`
	FinalityPolicy        json.RawMessage `json:"finality_policy"`
	FreshnessPolicy       json.RawMessage `json:"freshness_policy"`
	LastFinalizedUpdateAt *uint64         `json:"last_finalized_update_at"`
	Profile               json.RawMessage `json:"profile"`
}

type l2ConsensusStateJSON struct {
	StateRoot         string `json:"state_root"`       // B256 0x-hex
	IBCStorageRoot    string `json:"ibc_storage_root"` // B256 0x-hex
	TimestampNanos    uint64 `json:"timestamp_nanos"`
	L2Height          uint64 `json:"l2_height"`
	L2BlockHash       string `json:"l2_block_hash"`
	ParentHash        string `json:"parent_hash"`
	L1OriginNumber    uint64 `json:"l1_origin_number"`
	L1OriginHash      string `json:"l1_origin_hash"`
	FinalityLevel     string `json:"finality_level"`
	ProposalStatus    string `json:"proposal_status"`
	FirstAcceptedAt   uint64 `json:"first_accepted_at"`
	FinalityReachedAt uint64 `json:"finality_reached_at"`
	EvidenceHash      string `json:"evidence_hash"`
	RollupCommitment  string `json:"rollup_commitment"`
}

var (
	defaultL2FinalityPolicy  = json.RawMessage(`{"minimum_update_level":"unsafe","minimum_membership_level":"unsafe","require_resolved_proposal":false,"maturity_delay_seconds":0,"freeze_on_trusted_conflict":true,"safe_verification_mode":{"type":"disabled"}}`)
	defaultL2FreshnessPolicy = json.RawMessage(`{}`)
)

func l2PolicyJSON(name string, configured, fallback json.RawMessage) (json.RawMessage, error) {
	value := configured
	if len(value) == 0 {
		value = fallback
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(value, &object); err != nil || object == nil {
		return nil, fmt.Errorf("l2 client: %s must be a JSON object", name)
	}
	return append(json.RawMessage(nil), value...), nil
}

// BuildL2WasmClientState assembles the wasm ClientState + ConsensusState the L2
// MsgCreateClient submits: the client-state Data is ClientState<Profile> (the
// operator's verifier profile + the bootstrap height), the consensus-state Data is
// the bootstrap roots read from the L2 chain (ibc_storage_root = the router account
// storage root; timestamp in nanos).
func BuildL2WasmClientState(p L2ClientParams) (ibcexported.ClientState, ibcexported.ConsensusState, error) {
	checksumBz, err := decodeHexPrefixed(p.WasmChecksum)
	if err != nil {
		return nil, nil, fmt.Errorf("l2 client: decode wasm checksum: %w", err)
	}
	if len(p.RollupProfile) == 0 {
		return nil, nil, fmt.Errorf("l2 client: rollup profile is required")
	}
	if p.Bootstrap.TimestampSeconds > ^uint64(0)/1_000_000_000 {
		return nil, nil, fmt.Errorf("l2 client: bootstrap timestamp overflows nanoseconds")
	}
	finalityPolicy, err := l2PolicyJSON("finality policy", p.FinalityPolicy, defaultL2FinalityPolicy)
	if err != nil {
		return nil, nil, err
	}
	freshnessPolicy, err := l2PolicyJSON("freshness policy", p.FreshnessPolicy, defaultL2FreshnessPolicy)
	if err != nil {
		return nil, nil, err
	}

	clientState := l2ClientStateJSON{
		LatestHeight:          p.Bootstrap.Height,
		FrozenHeight:          nil,
		FinalityPolicy:        finalityPolicy,
		FreshnessPolicy:       freshnessPolicy,
		LastFinalizedUpdateAt: nil,
		Profile:               p.RollupProfile,
	}
	clientStateBz, err := json.Marshal(clientState)
	if err != nil {
		return nil, nil, fmt.Errorf("l2 client: marshal client state: %w", err)
	}

	consensus := l2ConsensusStateJSON{
		StateRoot:      p.Bootstrap.StateRoot.Hex(),
		IBCStorageRoot: p.Bootstrap.RouterStorageRoot.Hex(),
		TimestampNanos: p.Bootstrap.TimestampSeconds * 1_000_000_000,
		L2Height:       p.Bootstrap.Height,
		L2BlockHash:    p.Bootstrap.BlockHash.Hex(),
		ParentHash:     p.Bootstrap.ParentHash.Hex(),
		L1OriginHash:   common.Hash{}.Hex(),
		FinalityLevel:  "unsafe",
		ProposalStatus: "pending",
		// The wasm instantiate entrypoint replaces these with Cosmos block time.
		FirstAcceptedAt:   0,
		FinalityReachedAt: 0,
		EvidenceHash:      common.Hash{}.Hex(),
		RollupCommitment:  common.Hash{}.Hex(),
	}
	consensusBz, err := json.Marshal(consensus)
	if err != nil {
		return nil, nil, fmt.Errorf("l2 client: marshal consensus state: %w", err)
	}

	wasmClientState := &ibcwasmtypes.ClientState{
		Data:     clientStateBz,
		Checksum: checksumBz,
		LatestHeight: clienttypes.Height{
			RevisionNumber: 0,
			RevisionHeight: p.Bootstrap.Height,
		},
	}
	wasmConsensusState := &ibcwasmtypes.ConsensusState{Data: consensusBz}
	return wasmClientState, wasmConsensusState, nil
}

// CreateL2Client bootstraps one L2 rollup wasm light client on Cosmos and returns
// its auto-assigned client id. It goes through the shared generic wasm-create tx path
// (the same one CreateEthClient uses), passing p.CounterpartyClientID (the L2-side
// client that tracks Cosmos) so it is registered inline like the ETH beacon client.
// When that id is not yet configured it is empty, and counterparty registration is
// deferred to a later step (registering the ETH source's router id here would write a
// wrong, durable mapping).
func (w *Worker) CreateL2Client(stdCtx context.Context, ctx Context, p L2ClientParams) (string, error) {
	clientState, consensusState, err := BuildL2WasmClientState(p)
	if err != nil {
		return "", err
	}
	return w.TxHandler.CreateWasmClient(stdCtx, ctx, clientState, consensusState, p.CounterpartyClientID)
}
