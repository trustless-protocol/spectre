//! Optimistic Arbitrum BoLD v2 assertion verification.

use alloy_primitives::{keccak256, B256};
use l2_client::{
    canonical_header::{CanonicalEvmHeader, ExecutionHeaderFork},
    error::Error,
    evm_proof::{verify_bounded_account, ProofLimits},
    msg::{EvmAccountProof, EvmStorageProof},
    state::{Header as VerifiedHeader, Height},
    verification::{mapping_slot_bytes32, packed_storage_field, verify_storage},
};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::config::Profile;

/// Conservative bound for all proofs accepted in one Arbitrum update.
pub const PROOF_LIMITS: ProofLimits = ProofLimits {
    max_account_nodes: 64,
    max_storage_nodes: 64,
    max_storage_proofs: 1,
    max_node_bytes: 32 * 1024,
};

/// BoLD machine status encoded in `AssertionState`.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum MachineStatus {
    /// A running machine. Rollup creation normally rejects this, but hashing remains exact.
    Running,
    /// Successful terminal state.
    Finished,
    /// Errored terminal state.
    Errored,
}

impl MachineStatus {
    const fn solidity_value(self) -> u8 {
        match self {
            Self::Running => 0,
            Self::Finished => 1,
            Self::Errored => 2,
        }
    }
}

/// Global state committed by a BoLD assertion.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct GlobalState {
    /// L2 block hash followed by the L2 send root.
    #[schemars(with = "[String; 2]")]
    pub bytes32_vals: [B256; 2],
    /// Inbox position followed by position within the message.
    pub u64_vals: [u64; 2],
}

/// `AssertionState` fields emitted by `AssertionCreated` and hashed by BoLD.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct AssertionState {
    /// Global execution state.
    pub global_state: GlobalState,
    /// Terminal machine status.
    pub machine_status: MachineStatus,
    /// End-history root used by the challenge protocol.
    #[schemars(with = "String")]
    pub end_history_root: B256,
}

impl AssertionState {
    /// Reproduces Solidity `keccak256(abi.encode(state))` for the static nested struct.
    #[must_use]
    pub fn hash(self) -> B256 {
        let mut abi = [0_u8; 32 * 6];
        abi[..32].copy_from_slice(self.global_state.bytes32_vals[0].as_slice());
        abi[32..64].copy_from_slice(self.global_state.bytes32_vals[1].as_slice());
        abi[88..96].copy_from_slice(&self.global_state.u64_vals[0].to_be_bytes());
        abi[120..128].copy_from_slice(&self.global_state.u64_vals[1].to_be_bytes());
        abi[159] = self.machine_status.solidity_value();
        abi[160..192].copy_from_slice(self.end_history_root.as_slice());
        keccak256(abi)
    }
}

/// Event data sufficient to reconstruct a BoLD assertion hash.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct AssertionClaim {
    /// Parent assertion hash.
    #[schemars(with = "String")]
    pub parent_assertion_hash: B256,
    /// State after executing the asserted L2 segment.
    pub after_state: AssertionState,
    /// Sequencer inbox accumulator consumed by the assertion.
    #[schemars(with = "String")]
    pub inbox_acc: B256,
}

impl AssertionClaim {
    /// Reproduces BoLD v2 `RollupLib.assertionHash`.
    #[must_use]
    pub fn hash(self) -> B256 {
        let mut preimage = [0_u8; 96];
        preimage[..32].copy_from_slice(self.parent_assertion_hash.as_slice());
        preimage[32..64].copy_from_slice(self.after_state.hash().as_slice());
        preimage[64..].copy_from_slice(self.inbox_acc.as_slice());
        keccak256(preimage)
    }
}

/// Optimistic update proving a nonzero-status BoLD assertion in finalized L1 state.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct Header {
    /// Exact Ethereum beacon slot used by the L1 proofs.
    pub beacon_slot: u64,
    /// Ethereum execution state root authenticated at `beacon_slot`.
    #[schemars(with = "String")]
    pub l1_state_root: B256,
    /// Rollup account proof.
    pub rollup_proof: EvmAccountProof,
    /// Claimed assertion identifier.
    #[schemars(with = "String")]
    pub assertion_hash: B256,
    /// Proof of the packed assertion node containing `status`.
    pub assertion_proof: EvmStorageProof,
    /// `AssertionCreated` event fields used to recompute `assertion_hash`.
    pub assertion: AssertionClaim,
    /// Canonical Arbitrum L2 execution header.
    pub l2_header: CanonicalEvmHeader,
    /// Account proof for the profile's L2 router.
    pub router_proof: EvmAccountProof,
}

/// Verifies a pending or confirmed BoLD assertion and its L2 router state.
pub fn verify(
    profile: &Profile,
    authenticated_l1_root: B256,
    _authenticated_l1_timestamp: u64,
    header: &Header,
) -> Result<VerifiedHeader, Error> {
    if header.l1_state_root != authenticated_l1_root {
        return Err(Error::Proof(
            "Arbitrum header L1 root differs from host consensus state".into(),
        ));
    }
    if header.assertion_hash.is_zero() || header.assertion.hash() != header.assertion_hash {
        return Err(Error::Proof("BoLD assertion hash mismatch".into()));
    }
    let rollup = verify_bounded_account(
        PROOF_LIMITS,
        header.l1_state_root,
        profile.rollup,
        &header.rollup_proof,
    )?;
    let assertion_slot =
        mapping_slot_bytes32(header.assertion_hash, profile.assertions_mapping_slot);
    if header.assertion_proof.key != assertion_slot {
        return Err(Error::Proof(
            "unexpected BoLD assertion storage slot".into(),
        ));
    }
    verify_storage(PROOF_LIMITS, &rollup, &header.assertion_proof)?;
    let status = packed_storage_field(
        &header.assertion_proof.value,
        usize::from(profile.assertion_status_offset),
        1,
    )?[0];
    if status == 0 {
        return Err(Error::Proof("BoLD assertion does not exist".into()));
    }

    header
        .l2_header
        .validate_for_fork(ExecutionHeaderFork::London)?;
    let block_hash = header.l2_header.hash(ExecutionHeaderFork::London)?;
    if header.assertion.after_state.global_state.bytes32_vals[0] != block_hash {
        return Err(Error::Proof(
            "BoLD assertion commits to a different L2 block hash".into(),
        ));
    }
    let router = verify_bounded_account(
        PROOF_LIMITS,
        header.l2_header.state_root(),
        profile.common.l2_router,
        &header.router_proof,
    )?;
    Ok(VerifiedHeader {
        height: Height::new(header.l2_header.number())?,
        state_root: header.l2_header.state_root(),
        router_storage_root: B256::from(router.storage_root.0),
        timestamp_seconds: header.l2_header.timestamp(),
    })
}

#[cfg(test)]
mod tests {
    use alloy_primitives::B256;

    use super::{AssertionClaim, AssertionState, GlobalState, MachineStatus};

    #[test]
    fn assertion_hash_changes_with_every_event_component() {
        let claim = AssertionClaim {
            parent_assertion_hash: B256::with_last_byte(1),
            after_state: AssertionState {
                global_state: GlobalState {
                    bytes32_vals: [B256::with_last_byte(2), B256::with_last_byte(3)],
                    u64_vals: [4, 5],
                },
                machine_status: MachineStatus::Finished,
                end_history_root: B256::with_last_byte(6),
            },
            inbox_acc: B256::with_last_byte(7),
        };
        let expected = claim.hash();
        let mut changed = claim;
        changed.inbox_acc = B256::with_last_byte(8);
        assert_ne!(expected, changed.hash());
    }
}
