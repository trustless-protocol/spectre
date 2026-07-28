//!  query regression tests.

use alloy_primitives::B256;

use crate::{
    query::{check_for_misbehaviour, timestamp_at_height, verify_client_message},
    state::{FinalityLevel, FinalityPolicy, Header, Height, ProposalStatus},
};

fn header(height: u64, root: u8) -> Header {
    Header {
        height: Height::new(height).unwrap(),
        state_root: B256::with_last_byte(root),
        router_storage_root: B256::with_last_byte(root + 1),
        timestamp_seconds: 1,
        l2_block_hash: B256::with_last_byte(root + 2),
        parent_hash: B256::ZERO,
        l1_origin_number: 0,
        l1_origin_hash: B256::ZERO,
        finality_level: crate::state::FinalityLevel::Unsafe,
        proposal_status: crate::state::ProposalStatus::Pending,
        first_accepted_at: 0,
        finality_reached_at: 0,
        evidence_hash: B256::ZERO,
        rollup_commitment: B256::ZERO,
    }
}

#[test]
fn validates_headers_and_detects_only_same_height_conflicts() {
    let consensus = verify_client_message(&header(4, 1)).unwrap();
    assert_eq!(timestamp_at_height(&consensus), 1_000_000_000);
    let policy = FinalityPolicy::default();
    assert!(check_for_misbehaviour(&policy, &header(4, 1), &header(4, 3)).unwrap());
    assert!(!check_for_misbehaviour(&policy, &header(4, 1), &header(5, 3)).unwrap());
}

#[test]
fn ignores_promotions_and_conflicts_below_the_membership_threshold() {
    let unsafe_header = header(4, 1);
    let mut promotion = unsafe_header.clone();
    promotion.finality_level = FinalityLevel::Finalized;
    promotion.proposal_status = ProposalStatus::ResolvedValid;
    promotion.finality_reached_at = 10;
    promotion.evidence_hash = B256::with_last_byte(99);
    assert!(
        !check_for_misbehaviour(&FinalityPolicy::default(), &unsafe_header, &promotion,).unwrap()
    );

    let finalized_only = FinalityPolicy {
        minimum_membership_level: FinalityLevel::Finalized,
        require_resolved_proposal: true,
        ..FinalityPolicy::default()
    };
    assert!(!check_for_misbehaviour(&finalized_only, &header(4, 1), &header(4, 3),).unwrap());
}
