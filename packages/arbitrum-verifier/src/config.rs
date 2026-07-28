//! Immutable Arbitrum runtime profiles.

use alloy_primitives::{Address, B256};
use l2_client::state::{CommonProfile, RuntimeProfile};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

/// Runtime-selected Arbitrum deployment data.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct Profile {
    /// Chain/client/router data shared by all optimistic L2 profiles.
    pub common: CommonProfile,
    /// L1 Rollup account containing the configured assertion or node state.
    #[schemars(with = "String")]
    pub rollup: Address,
    /// Rollup protocol and its reviewed storage layout.
    pub protocol: RollupProtocol,
}

/// Mutually exclusive `RollupCore` storage layouts supported by the verifier.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(
    tag = "type",
    content = "value",
    rename_all = "snake_case",
    deny_unknown_fields
)]
pub enum RollupProtocol {
    /// `BoLD` v2 assertion-hash storage.
    BoldV2(BoldProfile),
    /// Pre-BoLD Nitro numeric-node storage.
    LegacyNitro(LegacyProfile),
}

/// Reviewed `BoLD` v2 storage layout.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct BoldProfile {
    /// Solidity storage slot of `mapping(bytes32 => AssertionNode) _assertions`.
    #[schemars(with = "String")]
    pub assertions_mapping_slot: B256,
    /// Packed byte offset of `AssertionNode.status` from the least-significant end of its slot.
    pub assertion_status_offset: u8,
    /// Reviewed `BoLD` contract-layout identifier.
    pub version: String,
    /// Authenticated status value meaning the assertion is confirmed.
    #[serde(default = "default_confirmed_status")]
    pub confirmed_status: u8,
    /// Authenticated status value meaning the assertion was rejected.
    #[serde(default = "default_rejected_status")]
    pub rejected_status: u8,
}

fn default_confirmed_status() -> u8 {
    2
}

fn default_rejected_status() -> u8 {
    3
}

/// Reviewed pre-BoLD Nitro storage layout.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct LegacyProfile {
    /// Packed slot containing `_latestConfirmed`, `_firstUnresolvedNode`, and
    /// `_latestNodeCreated`.
    #[schemars(with = "String")]
    pub node_lifecycle_slot: B256,
    /// Solidity storage slot of `mapping(uint64 => Node) _nodes`.
    #[schemars(with = "String")]
    pub nodes_mapping_slot: B256,
    /// Byte offset of `_latestConfirmed` from the least-significant end of the
    /// lifecycle slot.
    pub latest_confirmed_offset: u8,
    /// Byte offset of `_firstUnresolvedNode` from the least-significant end of
    /// the lifecycle slot.
    pub first_unresolved_offset: u8,
    /// Byte offset of `_latestNodeCreated` from the least-significant end of
    /// the lifecycle slot.
    pub latest_created_offset: u8,
    /// Storage-word offset of `Node.confirmData` from the node mapping base.
    pub confirm_data_offset: u8,
    /// Reviewed legacy Nitro contract-layout identifier.
    pub version: String,
}

impl RuntimeProfile for Profile {
    fn common(&self) -> &CommonProfile {
        &self.common
    }
}

/// Compatibility name used by existing tooling.
pub type Config = Profile;
