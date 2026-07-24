//! Strictly decoded shared message and proof containers.

use alloy_primitives::B256;
use cosmwasm_std::Binary;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

/// Direct Union-style client creation message.
#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct InstantiateMsg {
    /// JSON-encoded [`crate::state::ClientState`].
    pub client_state: Binary,
    /// JSON-encoded [`crate::state::ConsensusState`].
    pub consensus_state: Binary,
    /// Checksum retained in the outer ICS-08 Wasm client-state envelope.
    pub checksum: Binary,
}

/// The revision-zero height used by optimistic L2 clients.
#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct IbcHeight {
    /// Revision number, which must be zero for proof operations.
    pub revision_number: u64,
    /// L2 execution block height.
    pub revision_height: u64,
}

/// The Merkle path shape passed by the 08-wasm host.
#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct MerklePath {
    /// Ordered binary path components.
    pub key_path: Vec<Binary>,
}

/// Host sudo operations supported by every optimistic L2 client.
#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize, JsonSchema)]
#[serde(rename_all = "snake_case", deny_unknown_fields)]
pub enum SudoMsg {
    /// Applies an authenticated optimistic-rollup header.
    UpdateState {
        /// JSON-encoded chain-specific header.
        client_message: Binary,
    },
    /// Freezes only on valid same-height conflicting headers.
    UpdateStateOnMisbehaviour {
        /// JSON-encoded two-header evidence.
        client_message: Binary,
    },
    /// Verifies a router commitment.
    VerifyMembership {
        /// Stored consensus height used as the proof root.
        height: IbcHeight,
        /// Host delay period in nanoseconds; optimistic L2 clients do not use it.
        delay_time_period: u64,
        /// Host delay period in blocks; optimistic L2 clients do not use it.
        delay_block_period: u64,
        /// JSON-encoded bounded EVM storage proof.
        proof: Binary,
        /// The single ICS-26 commitment path supplied by the host.
        merkle_path: MerklePath,
        /// Expected EVM storage value.
        value: Binary,
    },
    /// Verifies router commitment absence or EVM zero.
    VerifyNonMembership {
        /// Stored consensus height used as the proof root.
        height: IbcHeight,
        /// Host delay period in nanoseconds; optimistic L2 clients do not use it.
        delay_time_period: u64,
        /// Host delay period in blocks; optimistic L2 clients do not use it.
        delay_block_period: u64,
        /// JSON-encoded bounded EVM storage proof.
        proof: Binary,
        /// The single ICS-26 commitment path supplied by the host.
        merkle_path: MerklePath,
    },
}

/// Response returned to the host after an update.
#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct UpdateStateResult {
    /// Consensus heights created by the update.
    pub heights: Vec<IbcHeight>,
}

/// Response returned by the misbehaviour query.
#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct CheckForMisbehaviourResult {
    /// Whether the supplied evidence proves misbehaviour.
    pub found_misbehaviour: bool,
}

/// Response returned by the timestamp query.
#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct TimestampAtHeightResult {
    /// Consensus timestamp in nanoseconds.
    pub timestamp: u64,
}

/// Response returned by the status query.
#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct StatusResult {
    /// Host-facing client status.
    pub status: String,
}

/// Read-only host queries supported by every optimistic L2 client.
#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize, JsonSchema)]
#[serde(rename_all = "snake_case", deny_unknown_fields)]
pub enum QueryMsg {
    /// Strictly validates a header.
    VerifyClientMessage {
        /// JSON-encoded chain-specific header.
        client_message: Binary,
    },
    /// Tests two headers for valid conflicting evidence.
    CheckForMisbehaviour {
        /// JSON-encoded two-header evidence.
        client_message: Binary,
    },
    /// Gets timestamp at a stored exact height.
    TimestampAtHeight {
        /// Exact stored L2 height.
        height: IbcHeight,
    },
    /// Gets status inherited from the pinned Ethereum client and local freeze state.
    Status {},
}

/// Bounded EVM account proof supplied in a typed rollup header role.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct EvmAccountProof {
    /// RLP-encoded trie nodes from the account proof.
    pub proof: Vec<Vec<u8>>,
}

/// Bounded EVM storage proof supplied in a typed rollup header role.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct EvmStorageProof {
    /// Declared storage key, which must equal the verifier-derived key.
    #[schemars(with = "String")]
    pub key: B256,
    /// Decoded Solidity word value represented as canonical bytes.
    pub value: Vec<u8>,
    /// RLP-encoded trie nodes from the storage proof.
    pub proof: Vec<Vec<u8>>,
}

/// One unambiguous client-message envelope.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(
    tag = "type",
    content = "value",
    rename_all = "snake_case",
    deny_unknown_fields
)]
pub enum ClientMessage<Header> {
    /// A normal L2 finality update header.
    Header(Header),
    /// Fully verified evidence containing two candidate headers.
    Misbehaviour {
        /// First candidate header.
        header_1: Header,
        /// Second candidate header.
        header_2: Header,
    },
}
