use alloy_primitives::{Address, B256};
use cosmwasm_std::{testing::MockStorage, Storage};
use serde::{Deserialize, Serialize, Serializer};

use crate::{
    error::Error,
    runtime,
    state::{ClientState, CommonProfile, ConsensusState, Header, Height, RuntimeProfile},
};

#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize)]
struct Profile(CommonProfile);

#[derive(Clone, Debug, PartialEq, Eq, Deserialize)]
struct FailingSerializeProfile(CommonProfile);

impl Serialize for FailingSerializeProfile {
    fn serialize<SerializerType>(
        &self,
        _serializer: SerializerType,
    ) -> Result<SerializerType::Ok, SerializerType::Error>
    where
        SerializerType: Serializer,
    {
        Err(serde::ser::Error::custom(
            "intentional serialization failure",
        ))
    }
}

impl RuntimeProfile for Profile {
    fn common(&self) -> &CommonProfile {
        &self.0
    }
    fn expected_profile_version() -> &'static str {
        "test_attestor_v1"
    }
}

impl RuntimeProfile for FailingSerializeProfile {
    fn common(&self) -> &CommonProfile {
        &self.0
    }

    fn expected_profile_version() -> &'static str {
        "test_attestor_v1"
    }
}

fn storage_snapshot(storage: &dyn Storage) -> [Option<Vec<u8>>; 3] {
    [
        storage.get(crate::store::HOST_CLIENT_STATE_KEY.as_bytes()),
        storage.get(crate::store::consensus_db_key(5).as_bytes()),
        storage.get(crate::store::consensus_db_key(7).as_bytes()),
    ]
}

fn profile() -> Profile {
    Profile(CommonProfile {
        l2_chain_id: 2,
        l2_router: Address::with_last_byte(1),
        commitment_slot: B256::with_last_byte(2),
        profile_version: "test_attestor_v1".into(),
        l2_header_fork: crate::canonical_header::ExecutionHeaderFork::London,
    })
}

fn header(height: u64, block: u8, parent: B256) -> Header {
    Header {
        height: Height::new(height).unwrap(),
        state_root: B256::with_last_byte(block),
        router_storage_root: B256::with_last_byte(block + 1),
        timestamp_seconds: height,
        l2_block_hash: B256::with_last_byte(block + 2),
        parent_hash: parent,
    }
}

fn bootstrap() -> (ClientState<Profile>, ConsensusState) {
    let header = header(5, 5, B256::with_last_byte(4));
    (
        ClientState {
            latest_height: 5,
            frozen_height: None,
            profile: profile(),
        },
        header.consensus_state(0).unwrap(),
    )
}

#[test]
fn stores_higher_updates_and_historical_backfills_without_rewinding() {
    let mut storage = MockStorage::new();
    let (client, consensus) = bootstrap();
    runtime::instantiate(&mut storage, &client, &consensus, vec![9], 10).unwrap();

    assert_eq!(
        runtime::update::<Profile>(&mut storage, &header(7, 7, B256::ZERO), 11).unwrap(),
        Some(7)
    );
    assert_eq!(
        runtime::update::<Profile>(&mut storage, &header(3, 3, B256::ZERO), 12).unwrap(),
        Some(3)
    );
    assert_eq!(
        runtime::client_state::<Profile>(&storage)
            .unwrap()
            .latest_height,
        7
    );
}

#[test]
fn constructs_explicit_update_effects_before_commit() {
    let (client, _) = bootstrap();
    let transition =
        runtime::build_update_transition(client, &header(7, 7, B256::ZERO), None, None, None, 11)
            .unwrap();

    assert_eq!(transition.client().unwrap().latest_height, 7);
    let (height, consensus) = transition.consensus().unwrap();
    assert_eq!(height, 7);
    assert_eq!(consensus.first_accepted_at, 11);
}

#[test]
fn validates_every_update_effect_before_writing_any_of_them() {
    let mut storage = MockStorage::new();
    let (client, consensus) = bootstrap();
    runtime::instantiate(&mut storage, &client, &consensus, vec![9], 10).unwrap();
    let before = storage_snapshot(&storage);

    assert!(matches!(
        runtime::update::<FailingSerializeProfile>(&mut storage, &header(7, 7, B256::ZERO), 11),
        Err(Error::JsonDecode(_))
    ));
    assert_eq!(storage_snapshot(&storage), before);
}

/// A re-submission of a block already stored is ignored.
///
/// First write wins, so the stored state does not depend on arrival order — an honest relayer
/// re-sends the same block routinely (a retry, a cached header, the same block re-attested at a
/// higher head), and with every root identical there is nothing to update.
#[test]
fn re_submitting_a_stored_block_is_ignored() {
    let mut storage = MockStorage::new();
    let (client, consensus) = bootstrap();
    runtime::instantiate(&mut storage, &client, &consensus, vec![9], 10).unwrap();
    let original = runtime::consensus_state(&storage, 5).unwrap();

    let resubmission = header(5, 5, B256::with_last_byte(4));
    for accepted_at in [20, 30, 40] {
        assert_eq!(
            runtime::update::<Profile>(&mut storage, &resubmission, accepted_at).unwrap(),
            None,
            "a re-submission must not be stored as an update"
        );
        assert_eq!(
            runtime::consensus_state(&storage, 5).unwrap(),
            original,
            "a re-submission rewrote the stored state"
        );
    }
}

#[test]
fn same_height_conflict_is_rejected_without_freezing_or_overwriting() {
    let mut storage = MockStorage::new();
    let (client, consensus) = bootstrap();
    runtime::instantiate(&mut storage, &client, &consensus, vec![9], 10).unwrap();
    let original = runtime::consensus_state(&storage, 5).unwrap();

    // Anyone can author contradictory bare headers, so a conflict without authenticated
    // attestations is rejected rather than acted upon: the update fails and the client stays open.
    assert!(matches!(
        runtime::update::<Profile>(&mut storage, &header(5, 8, B256::with_last_byte(4)), 11),
        Err(Error::StateConflict { height: 5 })
    ));
    assert_eq!(
        runtime::client_state::<Profile>(&storage)
            .unwrap()
            .frozen_height,
        None,
        "a conflict without authenticated attestations must not freeze the client"
    );
    assert_eq!(runtime::consensus_state(&storage, 5).unwrap(), original);

    // And the client keeps working: a later honest update still lands.
    assert_eq!(
        runtime::update::<Profile>(&mut storage, &header(6, 6, original.l2_block_hash), 12)
            .unwrap(),
        Some(6)
    );
}

#[test]
fn same_height_conflict_with_a_bad_parent_is_rejected_before_the_parent_check() {
    let mut storage = MockStorage::new();
    let (client, consensus) = bootstrap();
    runtime::instantiate(&mut storage, &client, &consensus, vec![9], 10).unwrap();
    runtime::update::<Profile>(&mut storage, &header(4, 2, B256::with_last_byte(1)), 11).unwrap();
    let original = runtime::consensus_state(&storage, 5).unwrap();
    let conflicting = header(5, 8, B256::with_last_byte(99));

    // The conflict is detected before the parent-hash check, so the reported error names the
    // conflict rather than the bad parent — the ordering is what makes the message useful.
    assert!(matches!(
        runtime::update::<Profile>(&mut storage, &conflicting, 12),
        Err(Error::StateConflict { height: 5 })
    ));
    assert_eq!(
        runtime::client_state::<Profile>(&storage)
            .unwrap()
            .frozen_height,
        None
    );
    assert_eq!(runtime::consensus_state(&storage, 5).unwrap(), original);
}

#[test]
fn parent_links_are_checked_in_both_directions() {
    let mut storage = MockStorage::new();
    let (client, consensus) = bootstrap();
    runtime::instantiate(&mut storage, &client, &consensus, vec![9], 10).unwrap();
    let block_five = consensus.l2_block_hash;
    let six = header(6, 6, block_five);
    runtime::update::<Profile>(&mut storage, &six, 11).unwrap();

    let bad_seven = header(7, 7, B256::with_last_byte(99));
    assert!(matches!(
        runtime::update::<Profile>(&mut storage, &bad_seven, 12),
        Err(Error::ParentHashMismatch)
    ));

    let mut storage = MockStorage::new();
    let (client, consensus) = bootstrap();
    runtime::instantiate(&mut storage, &client, &consensus, vec![9], 10).unwrap();
    runtime::update::<Profile>(&mut storage, &header(7, 7, B256::with_last_byte(99)), 11).unwrap();
    assert!(matches!(
        runtime::update::<Profile>(&mut storage, &header(6, 6, consensus.l2_block_hash), 12),
        Err(Error::ParentHashMismatch)
    ));
}

#[test]
fn bootstrap_rejects_zero_roots_and_wrong_profile_version() {
    let mut storage = MockStorage::new();
    let (mut client, mut consensus) = bootstrap();
    consensus.state_root = B256::ZERO;
    assert!(matches!(
        runtime::instantiate(&mut storage, &client, &consensus, vec![], 1),
        Err(Error::InvalidHeader(_))
    ));

    let (_, consensus) = bootstrap();
    client.profile.0.profile_version = "wrong".into();
    assert!(matches!(
        runtime::instantiate(&mut storage, &client, &consensus, vec![], 1),
        Err(Error::InvalidProfileVersion)
    ));
}
