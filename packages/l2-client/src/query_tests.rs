//!  query regression tests.

use alloy_primitives::B256;

use crate::{
    query::{check_for_misbehaviour, timestamp_at_height, verify_client_message},
    state::{Header, Height},
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
fn validates_headers_and_detects_only_same_height_conflicts() {
    let consensus = verify_client_message(&header(4, 1)).unwrap();
    assert_eq!(timestamp_at_height(&consensus), 1_000_000_000);
    assert!(check_for_misbehaviour(&header(4, 1), &header(4, 3)).unwrap());
    assert!(!check_for_misbehaviour(&header(4, 1), &header(5, 3)).unwrap());
}
