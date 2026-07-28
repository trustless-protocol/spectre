//! Strictly decoded shared message and proof containers.

use alloy_primitives::B256;
use cosmwasm_std::Binary;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::state::ProposalStatus;

/// Typed evidence requested by an L2 update. The verifier derives the resulting level.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(
    tag = "type",
    content = "value",
    rename_all = "snake_case",
    deny_unknown_fields
)]
pub enum FinalityEvidence {
    /// Sequencer/proposal evidence; never implies canonical L1 inclusion.
    Unsafe(UnsafeEvidence),
    /// Canonical but non-finalized L1 inclusion evidence.
    Safe(SafeEvidence),
    /// Finalized L1 inclusion plus rollup resolution evidence.
    Finalized(FinalizedEvidence),
}

/// Evidence for an Unsafe update.
#[derive(Clone, Copy, Debug, Default, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct UnsafeEvidence {
    /// Commitment identity authenticated by the rollup verifier.
    #[schemars(with = "String")]
    pub commitment: B256,
}

/// Evidence for Safe canonicality.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct SafeEvidence {
    /// L2 height bound by the attestation.
    pub l2_height: u64,
    /// L2 block hash bound by the attestation.
    #[schemars(with = "String")]
    pub l2_block_hash: B256,
    /// L1 origin hash bound by the attestation.
    #[schemars(with = "String")]
    pub l1_origin_hash: B256,
    /// L1 origin number bound by the attestation.
    pub l1_origin_number: u64,
    /// Authenticated safe-head hash.
    #[schemars(with = "String")]
    pub safe_l1_hash: B256,
    /// Expiry of the attestation in Unix seconds.
    pub valid_until: u64,
}

/// Evidence for finalized canonicality and optimistic proposal resolution.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct FinalizedEvidence {
    /// L1 origin hash bound by the finalized commitment.
    #[schemars(with = "String")]
    pub l1_origin_hash: B256,
    /// L1 origin number bound by the finalized commitment.
    pub l1_origin_number: u64,
    /// Independently proven proposal result.
    pub proposal_status: ProposalStatus,
}

/// Common wire shape for a rollup update before a chain-specific verifier normalizes it.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct L2Header<ExecutionHeader, RollupProof> {
    /// Canonical execution header.
    pub l2_header: ExecutionHeader,
    /// Ethereum consensus height used to authenticate the L1 state.
    pub l1_consensus_height: IbcHeight,
    /// Rollup-specific commitment proof.
    pub rollup_proof: RollupProof,
    /// Typed finality evidence request.
    pub finality_evidence: FinalityEvidence,
}

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
