//! Strict shared client-state and consensus-state models.

use alloy_primitives::{Address, B256};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::error::Error;

/// Immutable reference to the only Ethereum client trusted by an L2 client.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct L1Client {
    /// IBC identifier of the underlying Ethereum Wasm client.
    pub client_id: String,
    /// SHA-256 checksum of the expected Ethereum Wasm program.
    pub wasm_checksum: Vec<u8>,
}

/// Deployment data shared by every optimistic L2 verification profile.
///
/// Profiles are client state, rather than contract constants.  This lets one audited Wasm
/// checksum serve multiple deployments while keeping every trust assumption explicit and
/// immutable for the lifetime of a client.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct CommonProfile {
    /// Ethereum execution chain identifier.
    pub l1_chain_id: u64,
    /// Rollup execution chain identifier.
    pub l2_chain_id: u64,
    /// Pinned Ethereum light client and Wasm checksum.
    pub ethereum_client: L1Client,
    /// L2 `ICS26Router` account authenticated by update headers.
    #[schemars(with = "String")]
    pub l2_router: Address,
    /// Solidity mapping slot used for IBC commitments in the router.
    #[schemars(with = "String")]
    pub commitment_slot: B256,
    /// Human-readable rollup protocol/version identifier.
    pub rollup_version: String,
}

/// Access to the common portion of a chain-specific runtime profile.
pub trait RuntimeProfile {
    /// Returns the deployment data shared by all supported optimistic rollups.
    fn common(&self) -> &CommonProfile;
}

/// An IBC height restricted to the L2 revision-zero height space.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct Height {
    /// Revision number; always zero for supported L2 clients.
    pub revision_number: u64,
    /// Authenticated L2 execution block number.
    pub revision_height: u64,
}

impl Height {
    /// Constructs a valid revision-zero L2 height.
    pub const fn new(revision_height: u64) -> Result<Self, Error> {
        if revision_height == 0 {
            return Err(Error::ZeroHeight);
        }

        Ok(Self {
            revision_number: 0,
            revision_height,
        })
    }

    /// Validates that the height is revision zero and non-zero.
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

/// Authenticated state retained for one L2 execution block.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct ConsensusState {
    /// Authenticated L2 execution state root.
    #[schemars(with = "String")]
    pub state_root: B256,
    /// Storage root of the proven L2 IBC contract account.
    #[schemars(with = "String")]
    pub ibc_storage_root: B256,
    /// L2 block timestamp converted to nanoseconds using checked arithmetic.
    pub timestamp_nanos: u64,
}

/// An L2 header that has already passed its chain-specific finality verifier.
///
/// This is deliberately not a relayer wire type. OP, Base, and Arbitrum each decode and
/// authenticate their own public header before constructing this compact shared result.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct Header {
    /// Revision-zero L2 height represented by this header.
    pub height: Height,
    /// L2 execution-state root.
    #[schemars(with = "String")]
    pub state_root: B256,
    /// `ICS26Router` storage root.
    #[schemars(with = "String")]
    pub router_storage_root: B256,
    /// L2 timestamp in seconds, converted with checked arithmetic when stored.
    pub timestamp_seconds: u64,
}

impl Header {
    /// Validates structural invariants and constructs the stored consensus state.
    pub fn consensus_state(&self) -> Result<ConsensusState, Error> {
        self.height.validate()?;
        if self.state_root.is_zero() || self.router_storage_root.is_zero() {
            return Err(Error::InvalidHeader("roots must be non-zero"));
        }
        Ok(ConsensusState {
            state_root: self.state_root,
            ibc_storage_root: self.router_storage_root,
            timestamp_nanos: ConsensusState::timestamp_nanos_from_seconds(self.timestamp_seconds)?,
        })
    }
}

impl ConsensusState {
    /// Converts a seconds timestamp to nanoseconds without overflow.
    pub fn timestamp_nanos_from_seconds(seconds: u64) -> Result<u64, Error> {
        seconds
            .checked_mul(1_000_000_000)
            .ok_or(Error::TimestampOverflow)
    }
}

/// Active immutable and mutable L2 client state.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct ClientState<Profile> {
    /// Most recent authenticated L2 height.
    pub latest_height: u64,
    /// Height where verified misbehaviour froze this client.
    pub frozen_height: Option<u64>,
    /// Immutable deployment and verification profile.
    pub profile: Profile,
}
