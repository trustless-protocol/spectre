package l2rollup

import (
	"encoding/json"
	"fmt"
)

// ClientMessageData is the relayer's Go mirror of the UTF-8 JSON the L2 wasm light
// client verifies (the `data` inside the wasm ClientMessage). It is TRUSTLESS —
// NO relayer signature: the L2 contract verifies the L1 rollup proof witnesses
// against the shared cw-ics08-wasm-eth client (referenced by the immutable
// l1_client_id in ClientState, at l1_beacon_slot), then the L2 IBC-handler account
// proof against the resulting L2 root.
//
// Field content is per Dũng's spec; the JSON tags are PROPOSED — confirm the exact
// names byte-for-byte against the cw-ics08 L2 client's deserializer before this is
// the frozen contract. Per Dũng, do NOT add relayer signatures, l1_state_root,
// URLs, block tags, or configured contract addresses/slots (those are derived or
// immutable ClientState config, not part of the message).
type ClientMessageData struct {
	// L1BeaconSlot selects which finalized Ethereum consensus state (in the shared
	// cw-ics08-wasm-eth client) to verify the L1 rollup proof against. The shared
	// ETH client MUST already cover this slot — update it first.
	L1BeaconSlot uint64 `json:"l1_beacon_slot"`

	// L1 MPT witness nodes proving the L2's rollup contract state on L1 (RollupCore
	// for Arbitrum, AnchorStateRegistry for OP). Raw witness bytes, not references.
	L1AccountProof [][]byte `json:"l1_account_proof"`
	L1StorageProof [][]byte `json:"l1_storage_proof"`

	// L2HeaderRLP is the canonical RLP-encoded L2 Ethereum header.
	L2HeaderRLP []byte `json:"l2_header_rlp"`

	// L2IBCAccountProof is the eth_getProof account proof for the L2 IBC-handler
	// account (whose storage_root packet commitments are proven against).
	L2IBCAccountProof [][]byte `json:"l2_ibc_account_proof"`

	// Exactly one rollup-specific input set is present, matching the L2 client kind.
	Arbitrum *ArbitrumInputs `json:"arbitrum,omitempty"`
	Optimism *OptimismInputs `json:"optimism,omitempty"`
}

// ArbitrumInputs are the Nitro assertion inputs the L2 client checks against the
// RollupCore storage witness (in L1StorageProof).
type ArbitrumInputs struct {
	ParentAssertionHash []byte `json:"parent_assertion_hash"`
	AssertionState      []byte `json:"assertion_state"`
	InboxAcc            []byte `json:"inbox_acc"`
}

// OptimismInputs are the OP-Stack output-root inputs (Optimism + Base) the L2
// client checks against the AnchorStateRegistry output-root + L2-height storage
// witnesses (in L1StorageProof).
type OptimismInputs struct {
	OutputRootPreimage []byte `json:"output_root_preimage"`
}

// Encode marshals the client message data to the UTF-8 JSON the wasm client reads.
func (d *ClientMessageData) Encode() ([]byte, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("l2: encode client message data: %w", err)
	}
	return b, nil
}
