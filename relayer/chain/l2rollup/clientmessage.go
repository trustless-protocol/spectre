package l2rollup

// This file mirrors, byte-for-byte on the JSON wire, the header types the cw-ics08
// L2 wasm light clients deserialize. Ground truth, and there is no longer a
// per-chain half — the three verifier crates are identical adapters:
//   - packages/l2-client/src/msg.rs              → ClientMessage, AttestedL2Header,
//                                                  EvmAccountProof, EvmStorageProof
//   - packages/l2-client/src/canonical_header.rs → CanonicalEvmHeader
//
// The client verifies the header against itself (the canonical fork layout pinned in
// the profile, and the block hash derived from it) and the L2 router account proof
// against the state root that header declares. It does NOT verify any L1 settlement
// object: there is no dispute game, no BoLD assertion, and no shared
// cw-ics08-wasm-eth client in the picture any more. What makes a header trustworthy
// is the attestor gate on the relayer side, which the client cannot re-check —
// see attestedHeaderBuilder for what that does and does not cover.
//
// Field NAMES and value REPRESENTATIONS are the frozen encoding contract (Dũng, PR
// #245); the Rust structs carry `#[serde(deny_unknown_fields)]`, so emit exactly
// these fields and no extras.
//
// scalar wire types (hexBytes / byteList / byteMatrix / u256) live in wire.go.

// ClientMessage is a header that encodes to the wasm ClientMessage envelope
// {"type":"header","value":<header>}. *AttestedL2Header is the only implementation;
// the interface stays because the envelope is the stable part and the header shape
// is what the signed-attestation wire version will change.
type ClientMessage interface {
	// EncodeClientMessage marshals the header into the tagged ClientMessage JSON.
	EncodeClientMessage() ([]byte, error)
}

// EvmAccountProof mirrors l2-client `EvmAccountProof { proof: Vec<Vec<u8>> }`.
type EvmAccountProof struct {
	Proof byteMatrix `json:"proof"`
}

// EvmStorageProof mirrors l2-client `EvmStorageProof { key: B256, value: Vec<u8>,
// proof: Vec<Vec<u8>> }`. value is minimal big-endian bytes (1→[1], zero→[]).
type EvmStorageProof struct {
	Key   hexBytes   `json:"key"`
	Value byteList   `json:"value"`
	Proof byteMatrix `json:"proof"`
}

// CanonicalEvmHeader mirrors l2-client `CanonicalEvmHeader` — the full EVM execution
// header (l2_header) whose state_root the rollup output/assertion commits to. The
// optional tail is fork-gated on the Rust side (validate_for_fork): a Cancun L2
// carries base_fee..parent_beacon_block_root; Prague adds requests_hash. serde
// treats a missing Option field as None, so nil optionals are omitted here.
type CanonicalEvmHeader struct {
	ParentHash       hexBytes `json:"parent_hash"`
	OmmersHash       hexBytes `json:"ommers_hash"`
	Beneficiary      hexBytes `json:"beneficiary"`
	StateRoot        hexBytes `json:"state_root"`
	TransactionsRoot hexBytes `json:"transactions_root"`
	ReceiptsRoot     hexBytes `json:"receipts_root"`
	LogsBloom        hexBytes `json:"logs_bloom"`
	Difficulty       u256     `json:"difficulty"`
	Number           uint64   `json:"number"`
	GasLimit         uint64   `json:"gas_limit"`
	GasUsed          uint64   `json:"gas_used"`
	Timestamp        uint64   `json:"timestamp"`
	ExtraData        hexBytes `json:"extra_data"`
	MixHash          hexBytes `json:"mix_hash"`
	Nonce            hexBytes `json:"nonce"`

	BaseFeePerGas         *u256     `json:"base_fee_per_gas,omitempty"`
	WithdrawalsRoot       *hexBytes `json:"withdrawals_root,omitempty"`
	BlobGasUsed           *uint64   `json:"blob_gas_used,omitempty"`
	ExcessBlobGas         *uint64   `json:"excess_blob_gas,omitempty"`
	ParentBeaconBlockRoot *hexBytes `json:"parent_beacon_block_root,omitempty"`
	RequestsHash          *hexBytes `json:"requests_hash,omitempty"`
}

// AttestedL2Header mirrors l2-client `AttestedL2Header` — the only update shape the
// attestor-trusted clients accept, and identical for every chain.
//
// The settlement fields the per-chain headers used to carry (beacon_slot,
// l1_state_root, factory/game/assertion proofs, output-root preimages) are gone
// because nothing verifies them any more. So is the attestation metadata that
// briefly replaced them: the client reads no level, no attestor id and no L1
// origin, and a field carried but unchecked reads as load-bearing when it is not.
type AttestedL2Header struct {
	L2Header    CanonicalEvmHeader `json:"l2_header"`
	RouterProof EvmAccountProof    `json:"router_proof"`
}

// EncodeClientMessage marshals the header into the ClientMessage envelope. There is
// no inner envelope any more: the Arbitrum verifier used to wrap its header in a
// second tagged enum (`{"type":"bold_v2",...}`) because it had more than one header
// variant, and it no longer does.
func (h *AttestedL2Header) EncodeClientMessage() ([]byte, error) { return encodeHeaderMessage(h) }

// Header types satisfy ClientMessage.
var _ ClientMessage = (*AttestedL2Header)(nil)
