//! Malformed input regression tests.

use alloy_primitives::B256;

use crate::{
    error::Error,
    msg::EvmStorageProof,
    state::{Header, Height},
};

#[test]
fn rejects_timestamp_overflow_and_malformed_header_json() {
    let header = Header {
        height: Height::new(1).unwrap(),
        state_root: B256::with_last_byte(1),
        router_storage_root: B256::with_last_byte(2),
        timestamp_seconds: u64::MAX,
        l2_block_hash: B256::with_last_byte(3),
        parent_hash: B256::ZERO,
        l1_origin_number: 0,
        l1_origin_hash: B256::ZERO,
        finality_level: crate::state::FinalityLevel::Unsafe,
        proposal_status: crate::state::ProposalStatus::Pending,
        first_accepted_at: 0,
        finality_reached_at: 0,
        evidence_hash: B256::ZERO,
        rollup_commitment: B256::ZERO,
    };
    assert!(matches!(
        header.consensus_state(),
        Err(Error::TimestampOverflow)
    ));
    assert!(serde_json::from_slice::<Header>(
        br#"{"height":{"revision_number":0,"revision_height":1},"state_root":"0x01"}"#
    )
    .is_err());
}

#[test]
fn rejects_trailing_and_oversized_proof_data_before_traversal() {
    let malformed: Result<EvmStorageProof, _> = serde_json::from_slice(br#"{"key":"0x0000000000000000000000000000000000000000000000000000000000000000","value":[],"proof":[],"extra":true}"#);
    assert!(malformed.is_err());
    let proof = EvmStorageProof {
        key: B256::ZERO,
        value: vec![],
        proof: vec![vec![0; 32 * 1024 + 1]],
    };
    assert!(crate::packet::PROOF_LIMITS
        .validate_storage(&proof)
        .is_err());
}
