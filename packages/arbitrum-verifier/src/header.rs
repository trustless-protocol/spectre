//! Optimistic Arbitrum `BoLD` v2 and legacy Nitro verification.

use alloy_primitives::{keccak256, B256};
use l2_client::{
    canonical_header::{CanonicalEvmHeader, ExecutionHeaderFork},
    error::Error,
    evm_proof::{verify_bounded_account, ProofLimits},
    msg::{EvmAccountProof, EvmStorageProof, FinalityEvidence},
    state::{Header as VerifiedHeader, Height},
    verification::{mapping_slot_bytes32, packed_storage_field, storage_word, verify_storage},
};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::config::{BoldProfile, LegacyProfile, Profile, RollupProtocol};

/// Conservative bound for all proofs accepted in one Arbitrum update.
pub const PROOF_LIMITS: ProofLimits = ProofLimits {
    max_account_nodes: 64,
    max_storage_nodes: 64,
    max_storage_proofs: 2,
    max_node_bytes: 32 * 1024,
};

/// Mutually exclusive Arbitrum update-header formats.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(
    tag = "type",
    content = "value",
    rename_all = "snake_case",
    deny_unknown_fields
)]
pub enum Header {
    /// A `BoLD` v2 assertion update.
    BoldV2(BoldHeader),
    /// A pre-BoLD Nitro numeric-node update.
    LegacyNitro(LegacyHeader),
}

impl Header {
    /// Returns the exact Ethereum beacon slot authenticated by this header.
    #[must_use]
    pub const fn beacon_slot(&self) -> u64 {
        match self {
            Self::BoldV2(header) => header.beacon_slot,
            Self::LegacyNitro(header) => header.beacon_slot,
        }
    }
}

/// `BoLD` machine status encoded in `AssertionState`.
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

/// Global state committed by a `BoLD` assertion.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct GlobalState {
    /// L2 block hash followed by the L2 send root.
    #[schemars(with = "[String; 2]")]
    pub bytes32_vals: [B256; 2],
    /// Inbox position followed by position within the message.
    pub u64_vals: [u64; 2],
}

/// `AssertionState` fields emitted by `AssertionCreated` and hashed by `BoLD`.
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

/// Event data sufficient to reconstruct a `BoLD` assertion hash.
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
    /// Reproduces `BoLD` v2 `RollupLib.assertionHash`.
    #[must_use]
    pub fn hash(self) -> B256 {
        let mut preimage = [0_u8; 96];
        preimage[..32].copy_from_slice(self.parent_assertion_hash.as_slice());
        preimage[32..64].copy_from_slice(self.after_state.hash().as_slice());
        preimage[64..].copy_from_slice(self.inbox_acc.as_slice());
        keccak256(preimage)
    }
}

/// Optimistic update proving a nonzero-status `BoLD` assertion in finalized L1 state.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct BoldHeader {
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
    /// L1 origin number committed by the L2 derivation, when available.
    #[serde(default)]
    pub l1_origin_number: u64,
    /// L1 origin hash committed by the L2 derivation, when available.
    #[serde(default)]
    #[schemars(with = "String")]
    pub l1_origin_hash: B256,
    /// Optional typed finality evidence. The shared client independently validates it.
    #[serde(default)]
    pub finality_evidence: Option<FinalityEvidence>,
}

/// Update proving a pending or confirmed pre-BoLD Nitro node.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct LegacyHeader {
    /// Exact Ethereum beacon slot used by the L1 proofs.
    pub beacon_slot: u64,
    /// Ethereum execution state root authenticated at `beacon_slot`.
    #[schemars(with = "String")]
    pub l1_state_root: B256,
    /// Rollup account proof.
    pub rollup_proof: EvmAccountProof,
    /// Numeric legacy `RollupCore` node identifier.
    pub node_number: u64,
    /// Proof of the packed legacy node-lifecycle slot.
    pub node_lifecycle_proof: EvmStorageProof,
    /// Proof of `_nodes[node_number].confirmData`.
    pub confirm_data_proof: EvmStorageProof,
    /// L2 send root paired with the committed block hash by `confirmData`.
    #[schemars(with = "String")]
    pub send_root: B256,
    /// Canonical Arbitrum L2 execution header.
    pub l2_header: CanonicalEvmHeader,
    /// Account proof for the profile's L2 router.
    pub router_proof: EvmAccountProof,
    /// L1 origin number committed by the L2 derivation, when available.
    #[serde(default)]
    pub l1_origin_number: u64,
    /// L1 origin hash committed by the L2 derivation, when available.
    #[serde(default)]
    #[schemars(with = "String")]
    pub l1_origin_hash: B256,
    /// Optional typed finality evidence. The shared client independently validates it.
    #[serde(default)]
    pub finality_evidence: Option<FinalityEvidence>,
}

/// Verifies a protocol-selected Arbitrum assertion and its L2 router state.
///
/// # Errors
///
/// Returns an error when the header protocol differs from the immutable client
/// profile or any supplied L1, rollup, L2-header, or router proof is invalid.
pub fn verify(
    profile: &Profile,
    authenticated_l1_root: B256,
    authenticated_l1_timestamp: u64,
    header: &Header,
) -> Result<VerifiedHeader, Error> {
    match (&profile.protocol, header) {
        (RollupProtocol::BoldV2(protocol), Header::BoldV2(header)) => verify_bold(
            profile,
            protocol,
            authenticated_l1_root,
            authenticated_l1_timestamp,
            header,
        ),
        (RollupProtocol::LegacyNitro(protocol), Header::LegacyNitro(header)) => verify_legacy(
            profile,
            protocol,
            authenticated_l1_root,
            authenticated_l1_timestamp,
            header,
        ),
        _ => Err(Error::InvalidHeader(
            "Arbitrum header protocol differs from client profile",
        )),
    }
}

fn verify_bold(
    profile: &Profile,
    protocol: &BoldProfile,
    authenticated_l1_root: B256,
    authenticated_l1_timestamp: u64,
    header: &BoldHeader,
) -> Result<VerifiedHeader, Error> {
    verify_l1_root(authenticated_l1_root, header.l1_state_root)?;
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
        mapping_slot_bytes32(header.assertion_hash, protocol.assertions_mapping_slot);
    if header.assertion_proof.key != assertion_slot {
        return Err(Error::Proof(
            "unexpected BoLD assertion storage slot".into(),
        ));
    }
    verify_storage(PROOF_LIMITS, &rollup, &header.assertion_proof)?;
    let status = packed_storage_field(
        &header.assertion_proof.value,
        usize::from(protocol.assertion_status_offset),
        1,
    )?[0];
    if status == 0 {
        return Err(Error::Proof("BoLD assertion does not exist".into()));
    }

    let (finality_level, proposal_status) = if status == protocol.confirmed_status {
        (
            l2_client::state::FinalityLevel::Finalized,
            l2_client::state::ProposalStatus::ResolvedValid,
        )
    } else if status == protocol.rejected_status {
        (
            l2_client::state::FinalityLevel::Unsafe,
            l2_client::state::ProposalStatus::ResolvedInvalid,
        )
    } else {
        (
            l2_client::state::FinalityLevel::Unsafe,
            l2_client::state::ProposalStatus::Pending,
        )
    };

    let block_hash = validate_l2_header(&header.l2_header)?;
    if header.assertion.after_state.global_state.bytes32_vals[0] != block_hash {
        return Err(Error::Proof(
            "BoLD assertion commits to a different L2 block hash".into(),
        ));
    }
    verify_l2_router(
        profile,
        &header.l2_header,
        &header.router_proof,
        header.l1_origin_number,
        header.l1_origin_hash,
        authenticated_l1_timestamp,
        header.assertion_hash,
        finality_level,
        proposal_status,
        header,
    )
}

fn verify_legacy(
    profile: &Profile,
    protocol: &LegacyProfile,
    authenticated_l1_root: B256,
    authenticated_l1_timestamp: u64,
    header: &LegacyHeader,
) -> Result<VerifiedHeader, Error> {
    verify_l1_root(authenticated_l1_root, header.l1_state_root)?;
    if header.node_number == 0 {
        return Err(Error::Proof(
            "legacy Nitro node number must not be zero".into(),
        ));
    }
    if header.node_lifecycle_proof.key == header.confirm_data_proof.key {
        return Err(Error::Proof(
            "legacy Nitro storage proofs use the same key".into(),
        ));
    }
    PROOF_LIMITS.validate_storage(&header.node_lifecycle_proof)?;
    PROOF_LIMITS.validate_storage(&header.confirm_data_proof)?;

    let rollup = verify_bounded_account(
        PROOF_LIMITS,
        header.l1_state_root,
        profile.rollup,
        &header.rollup_proof,
    )?;
    if header.node_lifecycle_proof.key != protocol.node_lifecycle_slot {
        return Err(Error::Proof(
            "unexpected legacy Nitro lifecycle storage slot".into(),
        ));
    }
    verify_storage(PROOF_LIMITS, &rollup, &header.node_lifecycle_proof)?;
    let latest_confirmed = packed_u64(
        &header.node_lifecycle_proof.value,
        protocol.latest_confirmed_offset,
    )?;
    let first_unresolved = packed_u64(
        &header.node_lifecycle_proof.value,
        protocol.first_unresolved_offset,
    )?;
    let latest_created = packed_u64(
        &header.node_lifecycle_proof.value,
        protocol.latest_created_offset,
    )?;
    validate_legacy_lifecycle(latest_confirmed, first_unresolved, latest_created)?;
    if !legacy_node_is_active(
        header.node_number,
        latest_confirmed,
        first_unresolved,
        latest_created,
    ) {
        return Err(Error::Proof(
            "legacy Nitro node is neither pending nor latest confirmed".into(),
        ));
    }

    let (finality_level, proposal_status) = if header.node_number == latest_confirmed {
        (
            l2_client::state::FinalityLevel::Finalized,
            l2_client::state::ProposalStatus::ResolvedValid,
        )
    } else {
        (
            l2_client::state::FinalityLevel::Unsafe,
            l2_client::state::ProposalStatus::Pending,
        )
    };

    let node_base = mapping_slot_u64(header.node_number, protocol.nodes_mapping_slot);
    let confirm_data_slot = add_storage_offset(node_base, protocol.confirm_data_offset)?;
    if header.confirm_data_proof.key != confirm_data_slot {
        return Err(Error::Proof(
            "unexpected legacy Nitro confirmData storage slot".into(),
        ));
    }
    verify_storage(PROOF_LIMITS, &rollup, &header.confirm_data_proof)?;

    let block_hash = validate_l2_header(&header.l2_header)?;
    let mut confirm_preimage = [0_u8; 64];
    confirm_preimage[..32].copy_from_slice(block_hash.as_slice());
    confirm_preimage[32..].copy_from_slice(header.send_root.as_slice());
    if B256::from(storage_word(&header.confirm_data_proof.value)?) != keccak256(confirm_preimage) {
        return Err(Error::Proof(
            "legacy Nitro confirmData commits to a different L2 block".into(),
        ));
    }

    verify_l2_router(
        profile,
        &header.l2_header,
        &header.router_proof,
        header.l1_origin_number,
        header.l1_origin_hash,
        authenticated_l1_timestamp,
        B256::from(storage_word(&header.confirm_data_proof.value)?),
        finality_level,
        proposal_status,
        header,
    )
}

fn verify_l1_root(authenticated: B256, supplied: B256) -> Result<(), Error> {
    if supplied != authenticated {
        return Err(Error::Proof(
            "Arbitrum header L1 root differs from host consensus state".into(),
        ));
    }
    Ok(())
}

fn validate_l2_header(header: &CanonicalEvmHeader) -> Result<B256, Error> {
    header.validate_for_fork(ExecutionHeaderFork::London)?;
    header.hash(ExecutionHeaderFork::London)
}

fn verify_l2_router(
    profile: &Profile,
    l2_header: &CanonicalEvmHeader,
    router_proof: &EvmAccountProof,
    l1_origin_number: u64,
    l1_origin_hash: B256,
    authenticated_l1_timestamp: u64,
    rollup_commitment: B256,
    finality_level: l2_client::state::FinalityLevel,
    proposal_status: l2_client::state::ProposalStatus,
    evidence: &impl serde::Serialize,
) -> Result<VerifiedHeader, Error> {
    let router = verify_bounded_account(
        PROOF_LIMITS,
        l2_header.state_root(),
        profile.common.l2_router,
        router_proof,
    )?;
    Ok(VerifiedHeader {
        height: Height::new(l2_header.number())?,
        state_root: l2_header.state_root(),
        router_storage_root: B256::from(router.storage_root.0),
        timestamp_seconds: l2_header.timestamp(),
        l2_block_hash: validate_l2_header(l2_header)?,
        parent_hash: l2_header.parent_hash,
        l1_origin_number,
        l1_origin_hash,
        finality_level,
        proposal_status,
        first_accepted_at: authenticated_l1_timestamp,
        finality_reached_at: authenticated_l1_timestamp,
        evidence_hash: l2_client::verification::evidence_hash(evidence)?,
        rollup_commitment,
    })
}

fn packed_u64(value: &[u8], offset: u8) -> Result<u64, Error> {
    let bytes = packed_storage_field(value, usize::from(offset), 8)?;
    let bytes: [u8; 8] = bytes
        .try_into()
        .map_err(|_| Error::Proof("legacy Nitro uint64 field has invalid width".into()))?;
    Ok(u64::from_be_bytes(bytes))
}

fn validate_legacy_lifecycle(
    latest_confirmed: u64,
    first_unresolved: u64,
    latest_created: u64,
) -> Result<(), Error> {
    if latest_confirmed > latest_created || first_unresolved <= latest_confirmed {
        return Err(Error::Proof(
            "legacy Nitro lifecycle slot is inconsistent".into(),
        ));
    }
    let maximum_first_unresolved = latest_created
        .checked_add(1)
        .ok_or_else(|| Error::Proof("legacy Nitro latest node overflows uint64".into()))?;
    if first_unresolved > maximum_first_unresolved {
        return Err(Error::Proof(
            "legacy Nitro first unresolved node exceeds latest created".into(),
        ));
    }
    Ok(())
}

const fn legacy_node_is_active(
    node_number: u64,
    latest_confirmed: u64,
    first_unresolved: u64,
    latest_created: u64,
) -> bool {
    node_number == latest_confirmed
        || (node_number >= first_unresolved && node_number <= latest_created)
}

fn mapping_slot_u64(key: u64, slot: B256) -> B256 {
    let mut key_word = [0_u8; 32];
    key_word[24..].copy_from_slice(&key.to_be_bytes());
    mapping_slot_bytes32(B256::from(key_word), slot)
}

fn add_storage_offset(base: B256, offset: u8) -> Result<B256, Error> {
    let mut bytes = [0_u8; 32];
    bytes.copy_from_slice(base.as_slice());
    let mut carry = u16::from(offset);
    for byte in bytes.iter_mut().rev() {
        if carry == 0 {
            break;
        }
        let sum = u16::from(*byte) + carry;
        *byte = sum.to_be_bytes()[1];
        carry = sum >> 8;
    }
    if carry != 0 {
        return Err(Error::Proof(
            "legacy Nitro storage-slot offset overflows".into(),
        ));
    }
    Ok(B256::from(bytes))
}

#[cfg(test)]
mod tests {
    use alloy_primitives::{hex, keccak256, B256};

    use super::{
        add_storage_offset, legacy_node_is_active, mapping_slot_u64, packed_u64,
        validate_legacy_lifecycle, AssertionClaim, AssertionState, GlobalState, MachineStatus,
    };

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

    #[test]
    fn canonical_sepolia_legacy_layout_matches_live_node_10764() {
        let lifecycle = hex!("00000000003f5a360000000000002a0c0000000000002a0d0000000000002a0c");
        assert_eq!(packed_u64(&lifecycle, 0).unwrap(), 10_764);
        assert_eq!(packed_u64(&lifecycle, 8).unwrap(), 10_765);
        assert_eq!(packed_u64(&lifecycle, 16).unwrap(), 10_764);
        validate_legacy_lifecycle(10_764, 10_765, 10_764).unwrap();

        let nodes_slot = B256::with_last_byte(0x76);
        let node_base = mapping_slot_u64(10_764, nodes_slot);
        assert_eq!(
            node_base,
            "0xe4251833ee6ae54b9b2f5e13ed8882e17a2845215254f11f2b85a94b2f9b59b7"
                .parse::<B256>()
                .unwrap()
        );
        assert_eq!(
            add_storage_offset(node_base, 2).unwrap(),
            "0xe4251833ee6ae54b9b2f5e13ed8882e17a2845215254f11f2b85a94b2f9b59b9"
                .parse::<B256>()
                .unwrap()
        );

        let block_hash = hex!("fec7821ffcbba26c10f3e60964e2b7288262203c7dfa390285a61dc92d654ee5");
        let send_root = hex!("e51b1f0c5c8fcd2155d335c309a49c82807ecc19e159d3f60eeba1d3e63ecb77");
        let mut preimage = [0_u8; 64];
        preimage[..32].copy_from_slice(&block_hash);
        preimage[32..].copy_from_slice(&send_root);
        assert_eq!(
            keccak256(preimage),
            "0x9e3f9c983ceb4e5683856972d3302811aff89b0a9b5802a173aaa3e3b4700765"
                .parse::<B256>()
                .unwrap()
        );
    }

    #[test]
    fn legacy_lifecycle_rejects_invalid_state_and_resolved_losers() {
        validate_legacy_lifecycle(10, 12, 15).unwrap();
        assert!(legacy_node_is_active(10, 10, 12, 15));
        assert!(legacy_node_is_active(12, 10, 12, 15));
        assert!(legacy_node_is_active(15, 10, 12, 15));
        assert!(!legacy_node_is_active(9, 10, 12, 15));
        assert!(!legacy_node_is_active(11, 10, 12, 15));
        assert!(!legacy_node_is_active(16, 10, 12, 15));
        assert!(validate_legacy_lifecycle(12, 12, 15).is_err());
        assert!(validate_legacy_lifecycle(16, 17, 15).is_err());
        assert!(validate_legacy_lifecycle(10, 17, 15).is_err());
    }
}
