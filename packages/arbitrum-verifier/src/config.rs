//! Immutable Arbitrum BoLD runtime profile.

use alloy_primitives::{Address, B256};
use l2_client::state::{CommonProfile, RuntimeProfile};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

/// Runtime-selected BoLD v2 deployment data.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct Profile {
    /// Chain/client/router data shared by all optimistic L2 profiles.
    pub common: CommonProfile,
    /// L1 Rollup account containing the assertion mapping.
    #[schemars(with = "String")]
    pub rollup: Address,
    /// Solidity storage slot of `mapping(bytes32 => AssertionNode) _assertions`.
    #[schemars(with = "String")]
    pub assertions_mapping_slot: B256,
    /// Packed byte offset of `AssertionNode.status` from the least-significant end of its slot.
    pub assertion_status_offset: u8,
    /// Reviewed BoLD v2 layout identifier carried as profile metadata.
    pub bold_version: String,
}

impl RuntimeProfile for Profile {
    fn common(&self) -> &CommonProfile {
        &self.common
    }
}

/// Compatibility name used by existing tooling.
pub type Config = Profile;
