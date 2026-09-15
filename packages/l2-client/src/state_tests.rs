use cosmwasm_std::{testing::MockApi, Api, Binary};
use serde::Deserialize;

use crate::{
    error::Error,
    state::{AttestorConfig, MAX_ATTESTORS},
    verification::{
        build_attestation_statement, ATTESTATION_PROTOCOL_DOMAIN, ATTESTATION_STATEMENT_LENGTH,
    },
};

fn key(byte: u8) -> Binary {
    Binary::from(vec![byte; 32])
}

#[test]
fn accepts_canonical_attestor_sets_at_both_bounds() {
    let one = AttestorConfig {
        public_keys: vec![key(1)],
        threshold: 1,
    };
    assert!(one.validate().is_ok());

    let maximum = AttestorConfig {
        public_keys: (0..MAX_ATTESTORS)
            .map(|index| key(u8::try_from(index).unwrap()))
            .collect(),
        threshold: u16::try_from(MAX_ATTESTORS).unwrap(),
    };
    assert!(maximum.validate().is_ok());
}

#[test]
fn rejects_every_noncanonical_attestor_set_shape() {
    let cases = [
        (
            AttestorConfig {
                public_keys: vec![],
                threshold: 1,
            },
            "invalid attestor set",
        ),
        (
            AttestorConfig {
                public_keys: vec![key(1); MAX_ATTESTORS + 1],
                threshold: 1,
            },
            "attestor set exceeds maximum of 32",
        ),
        (
            AttestorConfig {
                public_keys: vec![key(1)],
                threshold: 0,
            },
            "invalid attestor threshold",
        ),
        (
            AttestorConfig {
                public_keys: vec![key(1)],
                threshold: 2,
            },
            "invalid attestor threshold",
        ),
        (
            AttestorConfig {
                public_keys: vec![Binary::from(vec![1; 31])],
                threshold: 1,
            },
            "invalid attestor public key length",
        ),
        (
            AttestorConfig {
                public_keys: vec![key(1), key(1)],
                threshold: 1,
            },
            "invalid attestor set",
        ),
        (
            AttestorConfig {
                public_keys: vec![key(2), key(1)],
                threshold: 1,
            },
            "invalid attestor set",
        ),
    ];

    for (config, expected) in cases {
        assert_eq!(config.validate().unwrap_err().to_string(), expected);
    }
}

#[test]
fn hashes_threshold_count_and_sorted_raw_keys() {
    let config = AttestorConfig {
        public_keys: vec![key(1), key(2), key(3)],
        threshold: 2,
    };
    assert_eq!(
        hex::encode(config.set_hash().unwrap()),
        "0e46b5cb19c2489cffffa98793641e03fc46f7b4668cebf1e8cce47c4935dbe5"
    );
}

#[test]
fn authentication_error_text_is_stable() {
    let cases = [
        (Error::InvalidAttestorSet, "invalid attestor set"),
        (
            Error::TooManyAttestors,
            "attestor set exceeds maximum of 32",
        ),
        (
            Error::InvalidAttestorThreshold,
            "invalid attestor threshold",
        ),
        (
            Error::InvalidAttestorPublicKeyLength,
            "invalid attestor public key length",
        ),
        (
            Error::InvalidAttestorSignatureCount,
            "invalid attestor signature count",
        ),
        (Error::InvalidAttestorIndex, "invalid attestor index"),
        (
            Error::DuplicateOrUnsortedSignatureIndex,
            "duplicate or unsorted attestor signature index",
        ),
        (
            Error::InvalidAttestorSignatureLength,
            "invalid attestor signature length",
        ),
        (
            Error::AttestorSignatureVerificationFailed,
            "attestor signature verification failed",
        ),
        (
            Error::UnsupportedUnsignedHeader,
            "unsigned L2 headers are unsupported",
        ),
        (
            Error::UnsupportedLegacyClientState,
            "legacy L2 client state is unsupported; create a fresh authenticated client",
        ),
    ];

    for (error, expected) in cases {
        assert_eq!(error.to_string(), expected);
    }
}

#[derive(Deserialize)]
struct VectorFile {
    protocol_domain: String,
    cases: Vec<VectorCase>,
}

#[derive(Deserialize)]
struct VectorCase {
    threshold: u16,
    public_keys: Vec<Binary>,
    attestor_set_hash: String,
    context: VectorContext,
    statement: String,
    signatures: Vec<VectorSignature>,
}

#[derive(Deserialize)]
struct VectorContext {
    l2_chain_id: u64,
    router_address: String,
    block_number: u64,
    block_hash: String,
    state_root: String,
}

#[derive(Deserialize)]
struct VectorSignature {
    attestor_index: u16,
    public_key: Binary,
    signature: Binary,
}

fn fixed<const N: usize>(encoded: &str) -> [u8; N] {
    hex::decode(encoded).unwrap().try_into().unwrap()
}

#[test]
fn v1_domain_and_statement_offsets_are_wire_locked() {
    let router = [0x11; 20];
    let set_hash = [0x22; 32];
    let block_hash = [0x33; 32];
    let state_root = [0x44; 32];
    let statement = build_attestation_statement(
        0x0102_0304_0506_0708,
        router,
        set_hash,
        0x1112_1314_1516_1718,
        block_hash,
        state_root,
    );

    assert_eq!(ATTESTATION_STATEMENT_LENGTH, 164);
    assert_eq!(
        hex::encode(ATTESTATION_PROTOCOL_DOMAIN),
        "7ee55ff0a9e64c455d44d77fa7ee083945f82a38ed7fb638516854726bd440ad"
    );
    assert_eq!(&statement[0..32], &ATTESTATION_PROTOCOL_DOMAIN);
    assert_eq!(&statement[32..40], &0x0102_0304_0506_0708_u64.to_be_bytes());
    assert_eq!(&statement[40..60], &router);
    assert_eq!(&statement[60..92], &set_hash);
    assert_eq!(
        &statement[92..100],
        &0x1112_1314_1516_1718_u64.to_be_bytes()
    );
    assert_eq!(&statement[100..132], &block_hash);
    assert_eq!(&statement[132..164], &state_root);
}

#[test]
fn decodes_every_go_vector_byte_for_byte() {
    let fixture =
        include_bytes!("../../../test/fixtures/wasm-contracts/l2-attestation-vectors.json");
    let vectors: VectorFile = serde_json::from_slice(fixture).unwrap();
    assert_eq!(
        vectors.protocol_domain,
        hex::encode(ATTESTATION_PROTOCOL_DOMAIN)
    );

    for vector in vectors.cases {
        let config = AttestorConfig {
            public_keys: vector.public_keys,
            threshold: vector.threshold,
        };
        let set_hash = config.set_hash().unwrap();
        assert_eq!(hex::encode(set_hash), vector.attestor_set_hash);

        let statement = build_attestation_statement(
            vector.context.l2_chain_id,
            fixed(&vector.context.router_address),
            set_hash,
            vector.context.block_number,
            fixed(&vector.context.block_hash),
            fixed(&vector.context.state_root),
        );
        assert_eq!(statement.len(), ATTESTATION_STATEMENT_LENGTH);
        assert_eq!(hex::encode(statement), vector.statement);
        assert_eq!(vector.signatures.len(), usize::from(vector.threshold));
        for (expected_index, signature) in vector.signatures.into_iter().enumerate() {
            assert_eq!(usize::from(signature.attestor_index), expected_index);
            assert_eq!(signature.public_key, config.public_keys[expected_index]);
            assert_eq!(signature.signature.len(), 64);
            assert!(MockApi::default()
                .ed25519_verify(
                    &statement,
                    signature.signature.as_slice(),
                    signature.public_key.as_slice(),
                )
                .unwrap());
        }
    }
}
