use alloy_primitives::{keccak256, B256};

use crate::packet::{commitment_path_hash, commitment_storage_slot, IBCSTORE_STORAGE_SLOT};

#[test]
fn derives_the_solidity_mapping_slot_from_the_erc7201_namespace() {
    let path_hash = B256::with_last_byte(7);
    let mut preimage = [0_u8; 64];
    preimage[..32].copy_from_slice(path_hash.as_slice());
    preimage[32..].copy_from_slice(IBCSTORE_STORAGE_SLOT.as_slice());
    assert_eq!(commitment_storage_slot(path_hash), keccak256(preimage));
}

#[test]
fn rejects_empty_commitment_paths() {
    assert!(commitment_path_hash(&[]).is_err());
    assert!(commitment_path_hash(b"commitments/test").is_err());
}

#[test]
fn accepts_router_binary_packet_paths() {
    for path_type in 1..=3 {
        let mut path = b"08-wasm-7".to_vec();
        path.push(path_type);
        path.extend_from_slice(&42_u64.to_be_bytes());
        assert_eq!(
            commitment_path_hash(&path).unwrap(),
            keccak256(path.as_slice())
        );
    }

    let mut unsupported = b"08-wasm-7".to_vec();
    unsupported.push(4);
    unsupported.extend_from_slice(&42_u64.to_be_bytes());
    assert!(commitment_path_hash(&unsupported).is_err());
    let mut non_utf8_client = vec![0xff, 1];
    non_utf8_client.extend_from_slice(&42_u64.to_be_bytes());
    assert!(commitment_path_hash(&non_utf8_client).is_err());
}
