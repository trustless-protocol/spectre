//! Storage-backed advance, backfill, idempotency, and conflict tests.

use alloy_primitives::{Address, B256};
use cosmwasm_std::testing::MockStorage;
use serde::{Deserialize, Serialize};

use crate::{
    error::Error,
    runtime,
    state::{ClientState, CommonProfile, ConsensusState, Header, Height, L1Client, RuntimeProfile},
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
    }
}

#[test]
fn stores_backfill_without_rewinding_and_protects_every_existing_height() {
    let mut storage = MockStorage::new();
    let client = ClientState {
        latest_height: 5,
        frozen_height: None,
        profile: profile(),
    };
    let initial = ConsensusState {
        state_root: B256::with_last_byte(5),
        ibc_storage_root: B256::with_last_byte(6),
        timestamp_nanos: 5,
    };
    runtime::instantiate(&mut storage, &client, &initial, vec![9]).unwrap();

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
        Err(Error::Conflict(3))
    ));
    assert_eq!(
        runtime::consensus_state(&storage, 3).unwrap(),
        header(3, 3).consensus_state().unwrap()
    );
}
