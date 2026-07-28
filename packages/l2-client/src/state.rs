//! Strict shared client-state and consensus-state models.

use alloy_primitives::{Address, B256};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::error::Error;

/// L2 canonicality established by authenticated evidence.
#[derive(
    Clone, Copy, Debug, Default, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize, JsonSchema,
)]
#[serde(rename_all = "snake_case")]
pub enum FinalityLevel {
    /// Sequencer or proposal evidence that remains reorg/dispute exposed.
    #[default]
    Unsafe,
    /// Canonical L1 inclusion has been authenticated, but is not irreversible.
    Safe,
    /// The relevant L1 commitment is finalized and irreversible under the configured model.
    Finalized,
}

impl FinalityLevel {
    /// Returns the stable event/telemetry spelling for this level.
    #[must_use]
    pub const fn as_str(self) -> &'static str {
        match self {
            Self::Unsafe => "unsafe",
            Self::Safe => "safe",
            Self::Finalized => "finalized",
        }
    }
}

/// Optimistic proposal resolution, tracked independently from L1 finality.
#[derive(Clone, Copy, Debug, Default, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ProposalStatus {
    /// The proposal is not known to have resolved successfully.
    #[default]
    Pending,
    /// The proposal was resolved successfully by authenticated evidence.
    ResolvedValid,
    /// The proposal was resolved against the accepted commitment.
    ResolvedInvalid,
}

impl ProposalStatus {
    /// Returns the stable event/telemetry spelling for this status.
    #[must_use]
    pub const fn as_str(self) -> &'static str {
        match self {
            Self::Pending => "pending",
            Self::ResolvedValid => "resolved_valid",
            Self::ResolvedInvalid => "resolved_invalid",
        }
    }
}

/// Mechanism accepted for authenticating a non-finalized canonical L1 head.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(
    tag = "type",
    content = "value",
    rename_all = "snake_case",
    deny_unknown_fields
)]
pub enum SafeVerificationMode {
    /// Safe evidence is unavailable. This is the secure default.
    Disabled,
    /// Safe evidence is supplied by an authenticated L1 client identified here.
    AuthenticatedL1Head {
        /// IBC identifier of the authenticated Ethereum client.
        ethereum_client_id: String,
    },
    /// Safe evidence is supplied by a threshold attestation over a committed validator set.
    ThresholdAttestation {
        /// Hash committing to the attestation validator set.
        #[schemars(with = "String")]
        validator_set_hash: B256,
        /// Minimum number of valid attestations required.
        threshold: u32,
    },
}

impl Default for SafeVerificationMode {
    fn default() -> Self {
        Self::Disabled
    }
}

/// Policy separating update acceptance from membership-proof usability.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(default, deny_unknown_fields)]
pub struct FinalityPolicy {
    /// Minimum verified level that may be stored.
    pub minimum_update_level: FinalityLevel,
    /// Minimum verified level usable for IBC membership proofs.
    pub minimum_membership_level: FinalityLevel,
    /// Whether optimistic proposals must have resolved successfully for membership.
    pub require_resolved_proposal: bool,
    /// Delay after finality is reached before membership use.
    pub maturity_delay_seconds: u64,
    /// Whether a conflict involving a trusted state freezes the client.
    pub freeze_on_trusted_conflict: bool,
    /// Authentication mechanism for Safe evidence.
    #[serde(default)]
    pub safe_verification_mode: SafeVerificationMode,
}

impl Default for FinalityPolicy {
    fn default() -> Self {
        Self {
            minimum_update_level: FinalityLevel::Unsafe,
            minimum_membership_level: FinalityLevel::Unsafe,
            require_resolved_proposal: false,
            maturity_delay_seconds: 0,
            freeze_on_trusted_conflict: true,
            safe_verification_mode: SafeVerificationMode::Disabled,
        }
    }
}

impl FinalityPolicy {
    /// Validates policy combinations that would otherwise make the client unusable.
    pub fn validate(&self) -> Result<(), Error> {
        match &self.safe_verification_mode {
            SafeVerificationMode::Disabled => {}
            SafeVerificationMode::AuthenticatedL1Head { ethereum_client_id }
                if ethereum_client_id.is_empty() =>
            {
                return Err(Error::InvalidFinalityPolicy(
                    "authenticated Safe mode requires an Ethereum client id",
                ));
            }
            SafeVerificationMode::AuthenticatedL1Head { .. } => {}
            SafeVerificationMode::ThresholdAttestation {
                validator_set_hash,
                threshold,
            } if validator_set_hash.is_zero() || *threshold == 0 => {
                return Err(Error::InvalidFinalityPolicy(
                    "threshold Safe mode requires a validator set and nonzero threshold",
                ));
            }
            SafeVerificationMode::ThresholdAttestation { .. } => {}
        }
        Ok(())
    }
}

/// Optional freshness limits for membership verification.
#[derive(Clone, Debug, Default, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(default, deny_unknown_fields)]
pub struct FreshnessPolicy {
    /// Maximum wall-clock time since the last finalized update.
    pub max_time_without_finalized_update: Option<u64>,
    /// Optional authenticated L2 height lag limit, when external height evidence exists.
    pub max_l2_height_lag: Option<u64>,
}

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
    /// L2 height redundantly retained for audit and conflict checks.
    #[serde(default)]
    pub l2_height: u64,
    /// Authenticated L2 block hash.
    #[serde(default)]
    #[schemars(with = "String")]
    pub l2_block_hash: B256,
    /// Authenticated parent block hash.
    #[serde(default)]
    #[schemars(with = "String")]
    pub parent_hash: B256,
    /// L1 origin number used to derive this L2 block.
    #[serde(default)]
    pub l1_origin_number: u64,
    /// L1 origin hash used to derive this L2 block.
    #[serde(default)]
    #[schemars(with = "String")]
    pub l1_origin_hash: B256,
    /// Verified canonicality level.
    #[serde(default)]
    pub finality_level: FinalityLevel,
    /// Verified optimistic-proposal status.
    #[serde(default)]
    pub proposal_status: ProposalStatus,
    /// First time this state was accepted, in Unix seconds.
    #[serde(default)]
    pub first_accepted_at: u64,
    /// Time at which the current finality level was reached, in Unix seconds.
    #[serde(default)]
    pub finality_reached_at: u64,
    /// Hash of the authenticated evidence used to create or promote the state.
    #[serde(default)]
    #[schemars(with = "String")]
    pub evidence_hash: B256,
    /// Rollup-specific commitment identity used during promotion.
    #[serde(default)]
    #[schemars(with = "String")]
    pub rollup_commitment: B256,
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
    /// Authenticated L2 block hash.
    #[schemars(with = "String")]
    pub l2_block_hash: B256,
    /// Authenticated parent block hash.
    #[schemars(with = "String")]
    pub parent_hash: B256,
    /// L1 origin number used by the L2 derivation.
    pub l1_origin_number: u64,
    /// L1 origin hash used by the L2 derivation.
    #[schemars(with = "String")]
    pub l1_origin_hash: B256,
    /// Independently verified canonicality level.
    pub finality_level: FinalityLevel,
    /// Independently verified proposal status.
    pub proposal_status: ProposalStatus,
    /// Time at which this evidence was accepted, in Unix seconds.
    pub first_accepted_at: u64,
    /// Time at which this evidence established its finality level, in Unix seconds.
    pub finality_reached_at: u64,
    /// Hash committing to the submitted evidence.
    #[schemars(with = "String")]
    pub evidence_hash: B256,
    /// Rollup-specific commitment identity.
    #[schemars(with = "String")]
    pub rollup_commitment: B256,
}

/// Normalized result returned by a rollup-specific verifier.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct VerifiedL2State {
    /// L2 execution height.
    pub l2_height: u64,
    /// L2 execution block hash.
    #[schemars(with = "String")]
    pub l2_block_hash: B256,
    /// L2 parent block hash.
    #[schemars(with = "String")]
    pub parent_hash: B256,
    /// L2 execution state root.
    #[schemars(with = "String")]
    pub state_root: B256,
    /// IBC router storage root.
    #[schemars(with = "String")]
    pub router_storage_root: B256,
    /// L1 origin number.
    pub l1_origin_number: u64,
    /// L1 origin hash.
    #[schemars(with = "String")]
    pub l1_origin_hash: B256,
    /// Verified canonicality level.
    pub finality_level: FinalityLevel,
    /// Verified proposal status.
    pub proposal_status: ProposalStatus,
    /// Hash of the authenticated evidence.
    #[schemars(with = "String")]
    pub evidence_hash: B256,
    /// Rollup-specific commitment identity.
    #[schemars(with = "String")]
    pub rollup_commitment: B256,
    /// L2 timestamp in seconds.
    pub timestamp_seconds: u64,
    /// Time at which the evidence was accepted.
    pub first_accepted_at: u64,
    /// Time at which the verified level was reached.
    pub finality_reached_at: u64,
}

impl VerifiedL2State {
    /// Converts the normalized verifier result to the shared update representation.
    pub fn into_header(self) -> Result<Header, Error> {
        Ok(Header {
            height: Height::new(self.l2_height)?,
            state_root: self.state_root,
            router_storage_root: self.router_storage_root,
            timestamp_seconds: self.timestamp_seconds,
            l2_block_hash: self.l2_block_hash,
            parent_hash: self.parent_hash,
            l1_origin_number: self.l1_origin_number,
            l1_origin_hash: self.l1_origin_hash,
            finality_level: self.finality_level,
            proposal_status: self.proposal_status,
            first_accepted_at: self.first_accepted_at,
            finality_reached_at: self.finality_reached_at,
            evidence_hash: self.evidence_hash,
            rollup_commitment: self.rollup_commitment,
        })
    }
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
            l2_height: self.height.revision_height,
            l2_block_hash: self.l2_block_hash,
            parent_hash: self.parent_hash,
            l1_origin_number: self.l1_origin_number,
            l1_origin_hash: self.l1_origin_hash,
            finality_level: self.finality_level,
            proposal_status: self.proposal_status,
            first_accepted_at: self.first_accepted_at,
            finality_reached_at: self.finality_reached_at,
            evidence_hash: self.evidence_hash,
            rollup_commitment: self.rollup_commitment,
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
    /// Explicit policy governing updates and membership verification.
    #[serde(default)]
    pub finality_policy: FinalityPolicy,
    /// Optional freshness restrictions for membership verification.
    #[serde(default)]
    pub freshness_policy: FreshnessPolicy,
    /// Unix time of the most recent accepted Finalized update.
    #[serde(default)]
    pub last_finalized_update_at: Option<u64>,
    /// Immutable deployment and verification profile.
    pub profile: Profile,
}
