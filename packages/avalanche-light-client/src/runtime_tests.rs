//! Runtime transition tests: bootstrap, update discipline, freeze, rotation.

use alloy_primitives::B256;
use cosmwasm_std::testing::MockStorage;

use crate::{
    error::Error,
    msg::MigrateMsg,
    runtime,
    state::{ConsensusState, Header, ValidatorSetState, WireValidator},
    verification::AuthenticatedHeader,
    verification_tests::fuji_client_state,
};

fn consensus_at(height: u64, block_byte: u8, parent_byte: u8) -> ConsensusState {
    ConsensusState {
        state_root: B256::with_last_byte(0xaa),
        ibc_storage_root: B256::with_last_byte(0xbb),
        timestamp_nanos: 1,
        height,
        block_hash: B256::with_last_byte(block_byte),
        parent_hash: B256::with_last_byte(parent_byte),
        first_accepted_at: 0,
    }
}

fn header_at(height: u64, block_byte: u8, parent_byte: u8) -> AuthenticatedHeader {
    AuthenticatedHeader::from_test_header(Header {
        height,
        state_root: B256::with_last_byte(0xaa),
        router_storage_root: B256::with_last_byte(0xbb),
        timestamp_seconds: 1,
        block_hash: B256::with_last_byte(block_byte),
        parent_hash: B256::with_last_byte(parent_byte),
    })
}

fn bootstrapped() -> MockStorage {
    let mut storage = MockStorage::new();
    let mut client = fuji_client_state();
    client.latest_height = 10;
    runtime::instantiate(
        &mut storage,
        &client,
        &consensus_at(10, 0x10, 0x0f),
        vec![1, 2, 3],
        7,
    )
    .unwrap();
    storage
}

#[test]
fn updates_advance_height_and_are_idempotent() {
    let mut storage = bootstrapped();

    // Child of the bootstrap block: parent hash links to the stored identity.
    runtime::update(&mut storage, &header_at(11, 0x11, 0x10), 8).unwrap();
    assert_eq!(runtime::client_state(&storage).unwrap().latest_height, 11);
    assert_eq!(runtime::consensus_state(&storage, 11).unwrap().height, 11);

    // Identical re-submission: accepted, nothing rewritten (first write wins).
    runtime::update(&mut storage, &header_at(11, 0x11, 0x10), 99).unwrap();
    assert_eq!(
        runtime::consensus_state(&storage, 11).unwrap().first_accepted_at,
        8
    );

    // Conflicting same-height submission: rejected.
    assert!(matches!(
        runtime::update(&mut storage, &header_at(11, 0x21, 0x10), 9),
        Err(Error::Invalid("conflicting state already stored"))
    ));

    // A child that does not link to the stored parent: rejected.
    assert!(matches!(
        runtime::update(&mut storage, &header_at(12, 0x12, 0x99), 9),
        Err(Error::Invalid("adjacent parent hash mismatch"))
    ));
}

#[test]
fn misbehaviour_freezes_and_blocks_everything() {
    let mut storage = bootstrapped();
    let first = header_at(20, 0x20, 0x1f);
    let second = header_at(20, 0x2f, 0x1f);
    assert!(runtime::is_actionable_misbehaviour(&first, &second).unwrap());
    assert!(runtime::apply_misbehaviour(&mut storage, &first, &second).unwrap());

    assert!(matches!(
        runtime::update(&mut storage, &header_at(11, 0x11, 0x10), 8),
        Err(Error::Frozen(20))
    ));
    assert!(matches!(
        runtime::verify_membership(&mut storage, 10, &[], &[], &[]),
        Err(Error::Frozen(20))
    ));

    // Identical headers are not misbehaviour.
    assert!(!runtime::is_actionable_misbehaviour(&first, &first).unwrap());
}

#[test]
fn rotation_is_monotonic_in_p_chain_height() {
    let mut storage = bootstrapped();
    let current = runtime::client_state(&storage).unwrap();
    let replacement = ValidatorSetState {
        validators: vec![WireValidator {
            public_key: cosmwasm_std::Binary::from(vec![7; avalanche_warp::PUBLIC_KEY_LEN]),
            weight: 5,
        }],
        total_weight: 5,
        p_chain_height: current.validator_set.p_chain_height,
    };

    // Same height: refused — governance cannot roll the set back or sideways.
    assert!(matches!(
        runtime::migrate(
            &mut storage,
            MigrateMsg::ReplaceValidatorSet {
                validators: replacement.validators.clone(),
                total_weight: replacement.total_weight,
                p_chain_height: replacement.p_chain_height,
            }
        ),
        Err(Error::NonMonotonicRotation { .. })
    ));

    // A strictly newer snapshot rotates.
    runtime::migrate(
        &mut storage,
        MigrateMsg::ReplaceValidatorSet {
            validators: replacement.validators,
            total_weight: replacement.total_weight,
            p_chain_height: current.validator_set.p_chain_height + 1,
        },
    )
    .unwrap();
    let rotated = runtime::client_state(&storage).unwrap();
    assert_eq!(rotated.validator_set.validators.len(), 1);
    assert_eq!(
        rotated.validator_set.p_chain_height,
        current.validator_set.p_chain_height + 1
    );

    // Keep is a validated no-op.
    runtime::migrate(&mut storage, MigrateMsg::KeepValidatorSet {}).unwrap();
}

/// Pins the bootstrap wire contract: `tools/warp-spike -dump-bootstrap` emits
/// {client_state, consensus_state} that must deserialize, validate, and
/// instantiate verbatim. The fixture is a live Fuji capture (73 validators,
/// settled-height proof).
#[test]
fn accepts_a_real_dumped_bootstrap() {
    #[derive(serde::Deserialize)]
    struct Bootstrap {
        client_state: crate::state::ClientState,
        consensus_state: ConsensusState,
    }
    let bootstrap: Bootstrap =
        serde_json::from_str(include_str!("../testdata/fuji_bootstrap.json")).unwrap();
    bootstrap.client_state.validate().unwrap();
    bootstrap.consensus_state.validate().unwrap();
    assert_eq!(bootstrap.client_state.network_id, 5);
    assert_eq!(bootstrap.client_state.evm_chain_id, 43113);
    assert_eq!(bootstrap.client_state.validator_set.validators.len(), 73);

    let mut storage = MockStorage::new();
    runtime::instantiate(
        &mut storage,
        &bootstrap.client_state,
        &bootstrap.consensus_state,
        vec![1],
        7,
    )
    .unwrap();
    assert_eq!(
        runtime::client_state(&storage).unwrap().latest_height,
        bootstrap.client_state.latest_height
    );
}
