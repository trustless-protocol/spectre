//! Finality policy, promotion, continuity, and evidence tests.

use alloy_primitives::B256;

use crate::{
    error::Error,
    msg::{FinalityEvidence, FinalizedEvidence, SafeEvidence, UnsafeEvidence},
    runtime::ensure_consensus_state_is_usable,
    state::{
        ClientState, ConsensusState, FinalityLevel, FinalityPolicy, FreshnessPolicy, Header,
        Height, ProposalStatus, SafeVerificationMode,
    },
    sudo::update_state_with_parent,
    verification::verify_finality_evidence,
};

fn header(level: FinalityLevel, proposal_status: ProposalStatus, block: u8) -> Header {
    Header {
        height: Height::new(7).unwrap(),
        state_root: B256::with_last_byte(10),
        router_storage_root: B256::with_last_byte(11),
        timestamp_seconds: 10,
        l2_block_hash: B256::with_last_byte(block),
        parent_hash: B256::with_last_byte(6),
        l1_origin_number: 9,
        l1_origin_hash: B256::with_last_byte(9),
        finality_level: level,
        proposal_status,
        first_accepted_at: 10,
        finality_reached_at: 10,
        evidence_hash: B256::with_last_byte(block + 20),
        rollup_commitment: B256::with_last_byte(30),
    }
}

fn client(policy: FinalityPolicy) -> ClientState<()> {
    ClientState {
        latest_height: 6,
        frozen_height: None,
        finality_policy: policy,
        freshness_policy: FreshnessPolicy::default(),
        last_finalized_update_at: None,
        profile: (),
    }
}

fn consensus(level: FinalityLevel, status: ProposalStatus) -> ConsensusState {
    header(level, status, 7).consensus_state().unwrap()
}

#[test]
fn finality_levels_have_stable_order_and_safe_is_disabled_by_default() {
    assert!(FinalityLevel::Unsafe < FinalityLevel::Safe);
    assert!(FinalityLevel::Safe < FinalityLevel::Finalized);
    assert!(matches!(
        FinalityPolicy::default().safe_verification_mode,
        SafeVerificationMode::Disabled
    ));
}

#[test]
fn promotes_the_same_block_without_changing_authenticated_roots() {
    let policy = FinalityPolicy {
        minimum_update_level: FinalityLevel::Unsafe,
        minimum_membership_level: FinalityLevel::Finalized,
        require_resolved_proposal: true,
        ..FinalityPolicy::default()
    };
    let mut client = client(policy);
    let unsafe_header = header(FinalityLevel::Unsafe, ProposalStatus::Pending, 7);
    let existing = unsafe_header.consensus_state().unwrap();
    assert!(
        update_state_with_parent(&mut client, &unsafe_header, None, None)
            .unwrap()
            .is_some()
    );

    let finalized = header(FinalityLevel::Finalized, ProposalStatus::ResolvedValid, 7);
    let promoted = update_state_with_parent(&mut client, &finalized, Some(&existing), None)
        .unwrap()
        .unwrap();
    assert_eq!(promoted.l2_block_hash, existing.l2_block_hash);
    assert_eq!(promoted.state_root, existing.state_root);
    assert_eq!(promoted.finality_level, FinalityLevel::Finalized);
    assert!(matches!(
        update_state_with_parent(&mut client, &unsafe_header, Some(&promoted), None),
        Err(Error::FinalityDowngrade)
    ));
}

#[test]
fn trusted_child_must_extend_the_previous_block() {
    let policy = FinalityPolicy {
        minimum_membership_level: FinalityLevel::Finalized,
        require_resolved_proposal: true,
        ..FinalityPolicy::default()
    };
    let mut client = client(policy);
    let previous = ConsensusState {
        l2_height: 6,
        l2_block_hash: B256::with_last_byte(1),
        finality_level: FinalityLevel::Finalized,
        proposal_status: ProposalStatus::ResolvedValid,
        ..consensus(FinalityLevel::Finalized, ProposalStatus::ResolvedValid)
    };
    let child = header(FinalityLevel::Finalized, ProposalStatus::ResolvedValid, 7);
    assert!(matches!(
        update_state_with_parent(&mut client, &child, None, Some(&previous)),
        Err(Error::ParentHashMismatch)
    ));
}

#[test]
fn legacy_consensus_state_does_not_enforce_parent_or_conflict_checks() {
    let mut client = client(FinalityPolicy::default());
    let legacy = ConsensusState {
        l2_height: 0,
        l2_block_hash: B256::ZERO,
        ..consensus(FinalityLevel::Unsafe, ProposalStatus::Pending)
    };
    let incoming = header(FinalityLevel::Unsafe, ProposalStatus::Pending, 8);

    assert!(
        update_state_with_parent(&mut client, &incoming, None, Some(&legacy))
            .unwrap()
            .is_some()
    );
    assert!(
        update_state_with_parent(&mut client, &incoming, Some(&legacy), None)
            .unwrap()
            .is_some()
    );
    assert_eq!(client.frozen_height, None);
}

#[test]
fn membership_policy_distinguishes_finality_maturity_and_staleness() {
    let policy = FinalityPolicy {
        minimum_membership_level: FinalityLevel::Finalized,
        require_resolved_proposal: true,
        maturity_delay_seconds: 5,
        ..FinalityPolicy::default()
    };
    let mut state = consensus(FinalityLevel::Finalized, ProposalStatus::ResolvedValid);
    let mut client = client(policy);
    assert!(matches!(
        ensure_consensus_state_is_usable(&client, &state, 14),
        Err(Error::FinalityMaturityPending { usable_at: 15 })
    ));
    client.finality_policy.maturity_delay_seconds = 0;
    client.freshness_policy.max_time_without_finalized_update = Some(5);
    client.last_finalized_update_at = Some(10);
    assert!(matches!(
        ensure_consensus_state_is_usable(&client, &state, 16),
        Err(Error::ClientStale)
    ));
    state.finality_level = FinalityLevel::Unsafe;
    assert!(matches!(
        ensure_consensus_state_is_usable(&client, &state, 10),
        Err(Error::InsufficientFinality { .. })
    ));
}

#[test]
fn membership_policy_enforces_authenticated_l2_height_lag() {
    let mut client = client(FinalityPolicy::default());
    client.latest_height = 10;
    client.freshness_policy.max_l2_height_lag = Some(3);
    let state = consensus(FinalityLevel::Unsafe, ProposalStatus::Pending);

    assert!(ensure_consensus_state_is_usable(&client, &state, 10).is_ok());

    client.freshness_policy.max_l2_height_lag = Some(2);
    assert!(matches!(
        ensure_consensus_state_is_usable(&client, &state, 10),
        Err(Error::ClientStale)
    ));
}

#[test]
fn resolved_invalid_state_is_never_usable_even_when_resolution_is_optional() {
    let client = client(FinalityPolicy::default());
    let state = consensus(FinalityLevel::Unsafe, ProposalStatus::ResolvedInvalid);
    assert!(matches!(
        ensure_consensus_state_is_usable(&client, &state, 10),
        Err(Error::ProposalResolvedInvalid)
    ));
}

#[test]
fn finalized_reorg_replaces_an_untrusted_unsafe_candidate() {
    let policy = FinalityPolicy {
        minimum_membership_level: FinalityLevel::Finalized,
        require_resolved_proposal: true,
        ..FinalityPolicy::default()
    };
    let mut client = client(policy);
    let existing = header(FinalityLevel::Unsafe, ProposalStatus::Pending, 7)
        .consensus_state()
        .unwrap();
    let canonical = header(FinalityLevel::Finalized, ProposalStatus::ResolvedValid, 8);
    let replacement = update_state_with_parent(&mut client, &canonical, Some(&existing), None)
        .unwrap()
        .unwrap();
    assert_eq!(replacement.l2_block_hash, canonical.l2_block_hash);
    assert_eq!(client.frozen_height, None);
}

#[test]
fn typed_evidence_cannot_claim_safe_or_finalized_without_authenticated_facts() {
    let default_policy = FinalityPolicy::default();
    let safe = FinalityEvidence::Safe(SafeEvidence {
        l2_height: 7,
        l2_block_hash: B256::with_last_byte(7),
        l1_origin_hash: B256::with_last_byte(9),
        l1_origin_number: 9,
        safe_l1_hash: B256::with_last_byte(8),
        valid_until: 100,
    });
    assert!(matches!(
        verify_finality_evidence(
            &default_policy,
            &safe,
            7,
            B256::with_last_byte(7),
            B256::with_last_byte(30),
            9,
            B256::with_last_byte(9),
            10,
            ProposalStatus::Pending,
        ),
        Err(Error::SafeVerificationUnavailable)
    ));
    let unsafe_evidence = FinalityEvidence::Unsafe(UnsafeEvidence {
        commitment: B256::with_last_byte(99),
    });
    assert!(matches!(
        verify_finality_evidence(
            &default_policy,
            &unsafe_evidence,
            7,
            B256::with_last_byte(7),
            B256::with_last_byte(30),
            9,
            B256::with_last_byte(9),
            10,
            ProposalStatus::Pending,
        ),
        Err(Error::InvalidFinalityEvidence)
    ));
    let unbound_unsafe = FinalityEvidence::Unsafe(UnsafeEvidence {
        commitment: B256::ZERO,
    });
    assert!(matches!(
        verify_finality_evidence(
            &default_policy,
            &unbound_unsafe,
            7,
            B256::with_last_byte(7),
            B256::with_last_byte(30),
            9,
            B256::with_last_byte(9),
            10,
            ProposalStatus::Pending,
        ),
        Err(Error::InvalidFinalityEvidence)
    ));
    let inconsistent_finalized = FinalityEvidence::Finalized(FinalizedEvidence {
        l1_origin_hash: B256::with_last_byte(9),
        l1_origin_number: 9,
        proposal_status: ProposalStatus::Pending,
    });
    assert!(matches!(
        verify_finality_evidence(
            &default_policy,
            &inconsistent_finalized,
            7,
            B256::with_last_byte(7),
            B256::with_last_byte(30),
            9,
            B256::with_last_byte(9),
            10,
            ProposalStatus::ResolvedValid,
        ),
        Err(Error::InvalidFinalityEvidence)
    ));
}

#[test]
fn configured_safe_modes_fail_closed_until_their_proof_verifiers_exist() {
    let policy = FinalityPolicy {
        safe_verification_mode: SafeVerificationMode::AuthenticatedL1Head {
            ethereum_client_id: "eth-client-0".to_owned(),
        },
        ..FinalityPolicy::default()
    };
    let safe = FinalityEvidence::Safe(SafeEvidence {
        l2_height: 7,
        l2_block_hash: B256::with_last_byte(7),
        l1_origin_hash: B256::with_last_byte(9),
        l1_origin_number: 9,
        safe_l1_hash: B256::with_last_byte(8),
        valid_until: 100,
    });

    assert!(matches!(
        verify_finality_evidence(
            &policy,
            &safe,
            7,
            B256::with_last_byte(7),
            B256::with_last_byte(30),
            9,
            B256::with_last_byte(9),
            10,
            ProposalStatus::Pending,
        ),
        Err(Error::SafeVerificationUnavailable)
    ));
}
