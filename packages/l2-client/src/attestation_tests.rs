use alloy_primitives::B256;
use cosmwasm_std::{testing::MockApi, Api};

use crate::attestation::{signing_bytes, AttestationHead};

#[test]
fn signing_bytes_are_domain_separated_and_big_endian() {
    let state_root = B256::with_last_byte(1);
    let block_hash = B256::with_last_byte(2);
    let bytes = signing_bytes(8453, AttestationHead::Safe, 123, state_root, block_hash);

    assert_eq!(&bytes[..27], b"fast-ibc/l2-attestation/v2\0");
    assert_eq!(&bytes[27..35], &8453_u64.to_be_bytes());
    assert_eq!(bytes[35], AttestationHead::Safe.signing_byte());
    assert_eq!(&bytes[36..44], &123_u64.to_be_bytes());
    assert_eq!(&bytes[44..76], state_root.as_slice());
    assert_eq!(&bytes[76..], block_hash.as_slice());
}

/// This vector is emitted by `attestor/types/attestation` from a fixed seed.
/// It proves that the CosmWasm verifier accepts the exact bytes and Ed25519
/// signature produced by the Go attestor implementation.
#[test]
fn accepts_the_go_attestor_signature_vector() {
    let public_key = [
        0x79, 0xb5, 0x56, 0x2e, 0x8f, 0xe6, 0x54, 0xf9, 0x40, 0x78, 0xb1, 0x12, 0xe8, 0xa9, 0x8b,
        0xa7, 0x90, 0x1f, 0x85, 0x3a, 0xe6, 0x95, 0xbe, 0xd7, 0xe0, 0xe3, 0x91, 0x0b, 0xad, 0x04,
        0x96, 0x64,
    ];
    let signature = [
        0x22, 0xcc, 0x6c, 0xde, 0xe3, 0xb8, 0xe8, 0x65, 0x94, 0x6a, 0xde, 0x58, 0xd2, 0xed, 0x25,
        0xef, 0xdc, 0x5c, 0xff, 0x7a, 0x91, 0x89, 0x70, 0xc0, 0xbc, 0x03, 0xae, 0xe0, 0x3b, 0x32,
        0xeb, 0xbd, 0x6a, 0x2d, 0x25, 0x36, 0x06, 0xf3, 0x7e, 0xec, 0xd0, 0xb6, 0xae, 0xa6, 0x51,
        0xf8, 0x98, 0xc4, 0x4a, 0xbe, 0xb8, 0xb4, 0x93, 0x49, 0xe2, 0xc4, 0xc8, 0xb4, 0x7b, 0xc4,
        0x3e, 0x1f, 0xe6, 0x09,
    ];
    let valid = MockApi::default()
        .ed25519_verify(
            &signing_bytes(
                8453,
                AttestationHead::Safe,
                123,
                B256::with_last_byte(1),
                B256::with_last_byte(2),
            ),
            &signature,
            &public_key,
        )
        .expect("well-formed Ed25519 inputs");
    assert!(valid);

    let wrong_finality = MockApi::default()
        .ed25519_verify(
            &signing_bytes(
                8453,
                AttestationHead::Unsafe,
                123,
                B256::with_last_byte(1),
                B256::with_last_byte(2),
            ),
            &signature,
            &public_key,
        )
        .expect("well-formed Ed25519 inputs");
    assert!(!wrong_finality);
}
