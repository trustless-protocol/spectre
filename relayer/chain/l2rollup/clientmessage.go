package l2rollup

// This file mirrors, byte-for-byte on the JSON wire, the header types the cw-ics08
// L2 wasm light client deserializes (#245). Ground truth:
//   - packages/l2-client/src/msg.rs            → EvmAccountProof, EvmStorageProof
//   - packages/l2-client/src/canonical_header.rs → CanonicalEvmHeader
//   - packages/op-stack-verifier/src/lib.rs    → OutputRootProof, OpStack Header
//   - packages/arbitrum-verifier/src/header.rs  → GlobalState, AssertionState,
//                                                  AssertionClaim, tagged Header enum
//                                                  (BoldHeader | LegacyHeader)
//
// These messages are TRUSTLESS — there is NO relayer signature. The L2 contract
// verifies the L1 rollup-contract proof against the shared cw-ics08-wasm-eth client
// (at beacon_slot / l1_state_root), derives the L2 state root, then verifies the L2
// router account proof against it. Field NAMES and value REPRESENTATIONS are the
// frozen encoding contract (Dũng, PR #245); the Rust structs carry
// `#[serde(deny_unknown_fields)]`, so emit exactly these fields and no extras.
//
// scalar wire types (hexBytes / byteList / byteMatrix / u256) live in wire.go.

// ClientMessage is a per-L2 header that encodes to the wasm ClientMessage envelope
// {"type":"header","value":<header>}. *OpStackHeader and *ArbitrumHeader implement it.
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

// OutputRootProof mirrors op-stack-verifier `OutputRootProof` — the four V0
// output-root preimage words (keccak of the concatenation is the game root claim).
type OutputRootProof struct {
	Version                  hexBytes `json:"version"`
	StateRoot                hexBytes `json:"state_root"`
	MessagePasserStorageRoot hexBytes `json:"message_passer_storage_root"`
	LatestBlockhash          hexBytes `json:"latest_blockhash"`
}

// OpStackHeader mirrors op-stack-verifier `Header` (Optimism + Base): proves a
// DisputeGameFactory game commitment in finalized L1 state, then binds the L2
// router state to the game's output root.
type OpStackHeader struct {
	BeaconSlot       uint64             `json:"beacon_slot"`
	L1StateRoot      hexBytes           `json:"l1_state_root"`
	FactoryProof     EvmAccountProof    `json:"factory_proof"`
	GameIndex        uint64             `json:"game_index"`
	GameProof        EvmStorageProof    `json:"game_proof"`
	GameAccountProof EvmAccountProof    `json:"game_account_proof"`
	GameRuntime      byteList           `json:"game_runtime"`
	OutputRootProof  OutputRootProof    `json:"output_root_proof"`
	L2Header         CanonicalEvmHeader `json:"l2_header"`
	RouterProof      EvmAccountProof    `json:"router_proof"`
}

// EncodeClientMessage marshals the OP-Stack header into the ClientMessage envelope.
func (h *OpStackHeader) EncodeClientMessage() ([]byte, error) { return encodeHeaderMessage(h) }

// MachineStatus mirrors arbitrum-verifier `MachineStatus` (serde snake_case).
type MachineStatus string

const (
	MachineStatusRunning  MachineStatus = "running"
	MachineStatusFinished MachineStatus = "finished"
	MachineStatusErrored  MachineStatus = "errored"
)

// GlobalState mirrors arbitrum-verifier `GlobalState`: [blockHash, sendRoot] and
// [inboxPosition, positionInMessage].
type GlobalState struct {
	Bytes32Vals [2]hexBytes `json:"bytes32_vals"`
	U64Vals     [2]uint64   `json:"u64_vals"`
}

// AssertionState mirrors arbitrum-verifier `AssertionState`.
type AssertionState struct {
	GlobalState    GlobalState   `json:"global_state"`
	MachineStatus  MachineStatus `json:"machine_status"`
	EndHistoryRoot hexBytes      `json:"end_history_root"`
}

// AssertionClaim mirrors arbitrum-verifier `AssertionClaim` — the AssertionCreated
// event fields sufficient to recompute the BoLD assertion hash.
type AssertionClaim struct {
	ParentAssertionHash hexBytes       `json:"parent_assertion_hash"`
	AfterState          AssertionState `json:"after_state"`
	InboxAcc            hexBytes       `json:"inbox_acc"`
}

// ArbitrumBoldHeader mirrors arbitrum-verifier `BoldHeader` — the `bold_v2` variant of
// the tagged Arbitrum `Header` enum: proves a nonzero-status BoLD assertion in finalized
// L1 state, then binds the L2 router state to the assertion's committed L2 block.
type ArbitrumBoldHeader struct {
	BeaconSlot     uint64             `json:"beacon_slot"`
	L1StateRoot    hexBytes           `json:"l1_state_root"`
	RollupProof    EvmAccountProof    `json:"rollup_proof"`
	AssertionHash  hexBytes           `json:"assertion_hash"`
	AssertionProof EvmStorageProof    `json:"assertion_proof"`
	Assertion      AssertionClaim     `json:"assertion"`
	L2Header       CanonicalEvmHeader `json:"l2_header"`
	RouterProof    EvmAccountProof    `json:"router_proof"`
}

// EncodeClientMessage marshals the header as the bold_v2 variant of the tagged Arbitrum
// Header enum, inside the ClientMessage envelope:
// {"type":"header","value":{"type":"bold_v2","value":<header>}}.
func (h *ArbitrumBoldHeader) EncodeClientMessage() ([]byte, error) {
	return encodeHeaderMessage(arbitrumHeaderEnvelope{Type: "bold_v2", Value: h})
}

// ArbitrumLegacyHeader mirrors arbitrum-verifier `LegacyHeader` — the `legacy_nitro`
// variant: proves a pending or confirmed pre-BoLD Nitro numeric node in finalized L1
// state (the packed node-lifecycle slot + `_nodes[node_number].confirmData`), binds it
// to the committed L2 block via confirmData = keccak(blockHash || send_root), then
// proves the L2 router state. node_lifecycle_proof and confirm_data_proof are two
// storage proofs against the same RollupCore account.
type ArbitrumLegacyHeader struct {
	BeaconSlot         uint64             `json:"beacon_slot"`
	L1StateRoot        hexBytes           `json:"l1_state_root"`
	RollupProof        EvmAccountProof    `json:"rollup_proof"`
	NodeNumber         uint64             `json:"node_number"`
	NodeLifecycleProof EvmStorageProof    `json:"node_lifecycle_proof"`
	ConfirmDataProof   EvmStorageProof    `json:"confirm_data_proof"`
	SendRoot           hexBytes           `json:"send_root"`
	L2Header           CanonicalEvmHeader `json:"l2_header"`
	RouterProof        EvmAccountProof    `json:"router_proof"`
}

// EncodeClientMessage marshals the header as the legacy_nitro variant:
// {"type":"header","value":{"type":"legacy_nitro","value":<header>}}.
func (h *ArbitrumLegacyHeader) EncodeClientMessage() ([]byte, error) {
	return encodeHeaderMessage(arbitrumHeaderEnvelope{Type: "legacy_nitro", Value: h})
}

// Header types satisfy ClientMessage.
var (
	_ ClientMessage = (*OpStackHeader)(nil)
	_ ClientMessage = (*ArbitrumBoldHeader)(nil)
	_ ClientMessage = (*ArbitrumLegacyHeader)(nil)
)
