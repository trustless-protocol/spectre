//! Storage-backed advance, backfill, idempotency, and conflict tests.

use alloy_primitives::{Address, B256};
use cosmwasm_std::testing::MockStorage;
use ibc_proto::ibc::lightclients::wasm::v1::ConsensusState as WasmConsensusState;
use serde::{Deserialize, Serialize};

use crate::{
    runtime,
    state::{ClientState, CommonProfile, ConsensusState, Header, Height, L1Client, RuntimeProfile},
    store::store_wasm_consensus_state,
};

#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize)]
struct Profile(CommonProfile);

impl RuntimeProfile for Profile {
    fn common(&self) -> &CommonProfile {
        &self.0
    }
}

fn profile() -> Profile {
    Profile(CommonProfile {
        l1_chain_id: 1,
        l2_chain_id: 2,
        ethereum_client: L1Client {
            client_id: "ethereum".into(),
            wasm_checksum: vec![1],
        },
        l2_router: Address::with_last_byte(1),
        commitment_slot: B256::with_last_byte(2),
        rollup_version: "test".into(),
    })
}

fn header(height: u64, root: u8) -> Header {
    Header {
        height: Height::new(height).unwrap(),
        state_root: B256::with_last_byte(root),
        router_storage_root: B256::with_last_byte(root + 1),
        timestamp_seconds: height,
        l2_block_hash: B256::with_last_byte(root + 2),
        parent_hash: B256::ZERO,
        l1_origin_number: 0,
        l1_origin_hash: B256::ZERO,
        finality_level: crate::state::FinalityLevel::Unsafe,
        proposal_status: crate::state::ProposalStatus::Pending,
        first_accepted_at: height,
        finality_reached_at: height,
        evidence_hash: B256::with_last_byte(root + 3),
        rollup_commitment: B256::with_last_byte(root + 4),
    }
}

#[test]
fn stores_backfill_without_rewinding_and_protects_every_existing_height() {
    let mut storage = MockStorage::new();
    let client = ClientState {
        latest_height: 5,
        frozen_height: None,
        finality_policy: crate::state::FinalityPolicy::default(),
        freshness_policy: crate::state::FreshnessPolicy::default(),
        last_finalized_update_at: None,
        profile: profile(),
    };
    let initial = ConsensusState {
        state_root: B256::with_last_byte(5),
        ibc_storage_root: B256::with_last_byte(6),
        timestamp_nanos: 5,
        l2_height: 5,
        l2_block_hash: B256::with_last_byte(7),
        parent_hash: B256::ZERO,
        l1_origin_number: 0,
        l1_origin_hash: B256::ZERO,
        finality_level: crate::state::FinalityLevel::Unsafe,
        proposal_status: crate::state::ProposalStatus::Pending,
        first_accepted_at: 5,
        finality_reached_at: 5,
        evidence_hash: B256::with_last_byte(8),
        rollup_commitment: B256::with_last_byte(9),
    };
    runtime::instantiate(&mut storage, &client, &initial, vec![9], 5).unwrap();

    assert_eq!(
        runtime::update::<Profile>(&mut storage, &header(7, 7)).unwrap(),
        Some(7)
    );
    assert_eq!(
        runtime::client_state::<Profile>(&storage)
            .unwrap()
            .latest_height,
        7
    );
    assert_eq!(
        runtime::update::<Profile>(&mut storage, &header(3, 3)).unwrap(),
        Some(3)
    );
    assert_eq!(
        runtime::client_state::<Profile>(&storage)
            .unwrap()
            .latest_height,
        7
    );
    assert_eq!(
        runtime::update::<Profile>(&mut storage, &header(3, 3)).unwrap(),
        None
    );
    assert!(matches!(
        runtime::update::<Profile>(&mut storage, &header(3, 4)),
        Ok(None)
    ));
    assert_eq!(
        runtime::client_state::<Profile>(&storage)
            .unwrap()
            .frozen_height,
        Some(3)
    );
    assert_eq!(
        runtime::consensus_state(&storage, 3).unwrap(),
        header(3, 3).consensus_state().unwrap()
    );
}

#[test]
fn finalized_freshness_never_regresses_and_bootstrap_height_must_match() {
    let mut storage = MockStorage::new();
    let client = ClientState {
        latest_height: 5,
        frozen_height: None,
        finality_policy: crate::state::FinalityPolicy::default(),
        freshness_policy: crate::state::FreshnessPolicy::default(),
        last_finalized_update_at: Some(20),
        profile: profile(),
    };
    let mut initial = header(5, 5).consensus_state().unwrap();
    initial.finality_level = crate::state::FinalityLevel::Finalized;
    initial.proposal_status = crate::state::ProposalStatus::ResolvedValid;
    initial.finality_reached_at = 10;
    runtime::instantiate(&mut storage, &client, &initial, vec![9], 25).unwrap();
    let stored = runtime::client_state::<Profile>(&storage).unwrap();
    assert_eq!(stored.last_finalized_update_at, Some(25));
    assert_eq!(
        runtime::consensus_state(&storage, 5)
            .unwrap()
            .first_accepted_at,
        25
    );

    let mut historical = header(3, 3);
    historical.finality_level = crate::state::FinalityLevel::Finalized;
    historical.proposal_status = crate::state::ProposalStatus::ResolvedValid;
    historical.finality_reached_at = 15;
    runtime::update::<Profile>(&mut storage, &historical).unwrap();
    assert_eq!(
        runtime::client_state::<Profile>(&storage)
            .unwrap()
            .last_finalized_update_at,
        Some(25)
    );

    let mut bad_storage = MockStorage::new();
    initial.l2_height = 4;
    assert!(matches!(
        runtime::instantiate(&mut bad_storage, &client, &initial, vec![9], 25),
        Err(crate::error::Error::InvalidHeader(_))
    ));
}

#[test]
fn legacy_consensus_json_can_be_replaced_at_the_same_height() {
    let mut storage = MockStorage::new();
    let client = ClientState {
        latest_height: 5,
        frozen_height: None,
        finality_policy: crate::state::FinalityPolicy::default(),
        freshness_policy: crate::state::FreshnessPolicy::default(),
        last_finalized_update_at: None,
        profile: profile(),
    };
    runtime::instantiate(
        &mut storage,
        &client,
        &header(5, 5).consensus_state().unwrap(),
        vec![9],
        5,
    )
    .unwrap();
    let legacy = serde_json::json!({
        "state_root": B256::with_last_byte(5),
        "ibc_storage_root": B256::with_last_byte(6),
        "timestamp_nanos": 5_000_000_000_u64,
    });
    store_wasm_consensus_state(
        &mut storage,
        5,
        &WasmConsensusState {
            data: serde_json::to_vec(&legacy).unwrap(),
        },
    )
    .unwrap();

    assert_eq!(runtime::consensus_state(&storage, 5).unwrap().l2_height, 5);
    assert_eq!(
        runtime::update::<Profile>(&mut storage, &header(5, 5)).unwrap(),
        Some(5)
    );
    assert_eq!(
        runtime::client_state::<Profile>(&storage)
            .unwrap()
            .frozen_height,
        None
    );
}

#[test]
fn legacy_predecessor_does_not_enforce_parent_linkage() {
    let mut storage = MockStorage::new();
    let client = ClientState {
        latest_height: 5,
        frozen_height: None,
        finality_policy: crate::state::FinalityPolicy::default(),
        freshness_policy: crate::state::FreshnessPolicy::default(),
        last_finalized_update_at: None,
        profile: profile(),
    };
    runtime::instantiate(
        &mut storage,
        &client,
        &header(5, 5).consensus_state().unwrap(),
        vec![9],
        5,
    )
    .unwrap();
    let legacy = serde_json::json!({
        "state_root": B256::with_last_byte(5),
        "ibc_storage_root": B256::with_last_byte(6),
        "timestamp_nanos": 5_000_000_000_u64,
    });
    store_wasm_consensus_state(
        &mut storage,
        5,
        &WasmConsensusState {
            data: serde_json::to_vec(&legacy).unwrap(),
        },
    )
    .unwrap();
    let mut child = header(6, 6);
    child.parent_hash = B256::with_last_byte(42);

    assert_eq!(
        runtime::update::<Profile>(&mut storage, &child).unwrap(),
        Some(6)
    );
}
