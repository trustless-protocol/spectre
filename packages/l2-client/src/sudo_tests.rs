//! Update and freeze transition tests.

use alloy_primitives::B256;

use crate::{
    error::Error,
    misbehaviour,
    state::{ClientState, Header, Height},
    sudo::update_state,
};

fn header(height: u64, root: u8) -> Header {
    Header {
        height: Height::new(height).unwrap(),
        state_root: B256::with_last_byte(root),
        router_storage_root: B256::with_last_byte(root + 1),
        timestamp_seconds: 1,
    }
}

#[test]
fn preserves_history_and_rejects_conflicts_or_updates_after_freeze() {
    let mut client = ClientState {
        latest_height: 1,
        frozen_height: None,
        profile: (),
    };
    let initial = header(1, 1).consensus_state().unwrap();
    assert!(update_state(&mut client, &header(1, 1), Some(&initial))
        .unwrap()
        .is_none());
    assert!(matches!(
        update_state(&mut client, &header(1, 3), Some(&initial)),
        Err(Error::Conflict(1))
    ));
    assert!(update_state(&mut client, &header(2, 5), None)
        .unwrap()
        .is_some());
    assert!(misbehaviour::apply(&mut client, &header(2, 5), &header(2, 7)).unwrap());
    assert!(matches!(
        update_state(&mut client, &header(3, 9), None),
        Err(Error::Frozen(2))
    ));
}
