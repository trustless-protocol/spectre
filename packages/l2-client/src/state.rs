//! State models for attestor-trusted L2 clients.

use alloy_primitives::{Address, B256};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::{canonical_header::ExecutionHeaderFork, error::Error};

/// Immutable deployment data needed for local L2 proof verification.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct CommonProfile {
    /// L2 execution chain identifier.
    pub l2_chain_id: u64,
    /// L2 `ICS26Router` address.
    #[schemars(with = "String")]
    pub l2_router: Address,
    /// Solidity commitment mapping slot.
    #[schemars(with = "String")]
    pub commitment_slot: B256,
    /// Artifact-specific profile version.
    pub profile_version: String,
    /// Canonical execution-header fork.
    pub l2_header_fork: ExecutionHeaderFork,
}

/// Access to an artifact-specific runtime profile.
pub trait RuntimeProfile {
    /// Returns shared profile fields.
    fn common(&self) -> &CommonProfile;
    /// Returns the only profile version accepted by this artifact.
    fn expected_profile_version() -> &'static str;
}

/// Revision-zero IBC height.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct Height {
    /// Must be zero.
    pub revision_number: u64,
    /// Nonzero L2 block number.
    pub revision_height: u64,
}

impl Height {
    /// Builds a revision-zero height.
    pub const fn new(revision_height: u64) -> Result<Self, Error> {
        if revision_height == 0 {
            return Err(Error::ZeroHeight);
        }
        Ok(Self {
            revision_number: 0,
            revision_height,
        })
    }

    /// Validates the supported height space.
    pub const fn validate(self) -> Result<(), Error> {
        if self.revision_number != 0 {
            return Err(Error::InvalidRevision(self.revision_number));
        }
        if self.revision_height == 0 {
            return Err(Error::ZeroHeight);
        }
        Ok(())
    }
}

/// State retained for one attested L2 block.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct ConsensusState {
    #[schemars(with = "String")]
    pub state_root: B256,
    #[schemars(with = "String")]
    pub ibc_storage_root: B256,
    pub timestamp_nanos: u64,
    pub l2_height: u64,
    #[schemars(with = "String")]
    pub l2_block_hash: B256,
    #[schemars(with = "String")]
    pub parent_hash: B256,
    pub first_accepted_at: u64,
}

impl ConsensusState {
    /// Converts seconds to nanoseconds with overflow checking.
    pub fn timestamp_nanos_from_seconds(seconds: u64) -> Result<u64, Error> {
        seconds
            .checked_mul(1_000_000_000)
            .ok_or(Error::TimestampOverflow)
    }

    /// Validates bootstrap and stored-state invariants.
    pub fn validate(&self) -> Result<(), Error> {
        Height::new(self.l2_height)?;
        if self.state_root.is_zero()
            || self.ibc_storage_root.is_zero()
            || self.l2_block_hash.is_zero()
        {
            return Err(Error::InvalidHeader(
                "consensus roots and block hash must be non-zero",
            ));
        }
        Ok(())
    }

    /// Returns whether another state disagrees on block identity.
    #[must_use]
    pub fn conflicts_with(&self, other: &Self) -> bool {
        self.l2_block_hash != other.l2_block_hash
            || self.parent_hash != other.parent_hash
            || self.state_root != other.state_root
            || self.ibc_storage_root != other.ibc_storage_root
    }
}

/// Locally verified normalized update.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct Header {
    pub height: Height,
    #[schemars(with = "String")]
    pub state_root: B256,
    #[schemars(with = "String")]
    pub router_storage_root: B256,
    pub timestamp_seconds: u64,
    #[schemars(with = "String")]
    pub l2_block_hash: B256,
    #[schemars(with = "String")]
    pub parent_hash: B256,
}

impl Header {
    /// Converts the verified header into stored consensus state.
    pub fn consensus_state(&self, accepted_at: u64) -> Result<ConsensusState, Error> {
        self.height.validate()?;
        let state = ConsensusState {
            state_root: self.state_root,
            ibc_storage_root: self.router_storage_root,
            timestamp_nanos: ConsensusState::timestamp_nanos_from_seconds(self.timestamp_seconds)?,
            l2_height: self.height.revision_height,
            l2_block_hash: self.l2_block_hash,
            parent_hash: self.parent_hash,
            first_accepted_at: accepted_at,
        };
        state.validate()?;
        Ok(state)
    }

    /// Checks for a well-formed, same-height state conflict.
    pub fn conflicts_with(&self, other: &Self) -> Result<bool, Error> {
        Ok(self.height == other.height
            && self
                .consensus_state(0)?
                .conflicts_with(&other.consensus_state(0)?))
    }
}

/// Active L2 client state.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct ClientState<Profile> {
    pub latest_height: u64,
    pub frozen_height: Option<u64>,
    pub profile: Profile,
}
