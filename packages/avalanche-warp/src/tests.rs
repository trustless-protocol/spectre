//! Tests against a REAL verified Fuji aggregate.
//!
//! The fixture was captured live by `tools/warp-spike -dump-fixture` on
//! 2026-09-21: a 67% stake-weighted aggregate over C-Chain block 58534295,
//! independently verified with avalanchego's own `BitSetSignature.Verify`
//! before being written. The canonical validator set is stored exactly as the
//! P-Chain served it (compressed keys, canonical order), so these tests pin
//! this crate to the network's real behavior, not to synthetic keys.

use serde::Deserialize;

use crate::{
    build_block_hash_message, parse_signed_message, verify_warp_aggregate, BlsVerify, Error,
    SignerBitSet, Validator, PUBLIC_KEY_LEN, SIGNATURE_LEN,
};

#[derive(Deserialize)]
struct FixtureValidator {
    public_key_compressed: String,
    weight: u64,
}

#[derive(Deserialize)]
struct Fixture {
    network_id: u32,
    source_chain_id: String,
    block_hash: String,
    unsigned_message: String,
    signed_message: String,
    bit_set: String,
    signature: String,
    quorum_num: u64,
    quorum_den: u64,
    total_weight: u64,
    signed_weight: u64,
    validators: Vec<FixtureValidator>,
}

fn unhex(s: &str) -> Vec<u8> {
    hex::decode(s.trim_start_matches("0x")).unwrap()
}

fn fixture() -> (Fixture, Vec<Validator>) {
    let fix: Fixture =
        serde_json::from_str(include_str!("../testdata/fuji_warp_fixture.json")).unwrap();
    let validators = fix
        .validators
        .iter()
        .map(|v| Validator {
            public_key: unhex(&v.public_key_compressed).try_into().unwrap(),
            weight: v.weight,
        })
        .collect();
    (fix, validators)
}

/// Real BLS via milagro (test-only, same ciphersuite as Avalanche warp).
struct MilagroBls;

impl BlsVerify for MilagroBls {
    fn fast_aggregate_verify(
        &self,
        public_keys: &[&[u8; PUBLIC_KEY_LEN]],
        message: &[u8],
        signature: &[u8; SIGNATURE_LEN],
    ) -> Result<(), Error> {
        let keys = public_keys
            .iter()
            .map(|pk| milagro_bls::PublicKey::from_bytes(pk.as_slice()))
            .collect::<Result<Vec<_>, _>>()
            .map_err(|_| Error::InvalidSignature)?;
        let key_refs: Vec<&milagro_bls::PublicKey> = keys.iter().collect();
        let aggregate_key = milagro_bls::AggregatePublicKey::aggregate(&key_refs)
            .map_err(|_| Error::InvalidSignature)?;
        let signature = milagro_bls::Signature::from_bytes(signature.as_slice())
            .map_err(|_| Error::InvalidSignature)?;
        let aggregate_signature = milagro_bls::AggregateSignature::aggregate(&[&signature]);
        if aggregate_signature.fast_aggregate_verify_pre_aggregated(message, &aggregate_key) {
            Ok(())
        } else {
            Err(Error::InvalidSignature)
        }
    }
}

#[test]
fn verifies_the_real_fuji_aggregate() {
    let (fix, validators) = fixture();
    let signed_weight = verify_warp_aggregate(
        &validators,
        fix.total_weight,
        &unhex(&fix.unsigned_message),
        &unhex(&fix.bit_set),
        &unhex(&fix.signature).try_into().unwrap(),
        fix.quorum_num,
        fix.quorum_den,
        &MilagroBls,
    )
    .unwrap();
    assert_eq!(signed_weight, fix.signed_weight);
}

#[test]
fn rebuilds_the_exact_unsigned_message_from_its_parts() {
    let (fix, _) = fixture();
    let rebuilt = build_block_hash_message(
        fix.network_id,
        &unhex(&fix.source_chain_id).try_into().unwrap(),
        &unhex(&fix.block_hash).try_into().unwrap(),
    );
    assert_eq!(rebuilt, unhex(&fix.unsigned_message));
}

#[test]
fn splits_the_signed_message_into_the_fixture_parts() {
    let (fix, _) = fixture();
    let signed = unhex(&fix.signed_message);
    let parsed = parse_signed_message(&signed).unwrap();
    assert_eq!(parsed.unsigned_message, unhex(&fix.unsigned_message));
    assert_eq!(parsed.bit_set, unhex(&fix.bit_set));
    assert_eq!(parsed.signature.to_vec(), unhex(&fix.signature));
}

#[test]
fn rejects_a_tampered_message() {
    let (fix, validators) = fixture();
    let mut message = unhex(&fix.unsigned_message);
    let last = message.len() - 1;
    message[last] ^= 1;
    assert_eq!(
        verify_warp_aggregate(
            &validators,
            fix.total_weight,
            &message,
            &unhex(&fix.bit_set),
            &unhex(&fix.signature).try_into().unwrap(),
            fix.quorum_num,
            fix.quorum_den,
            &MilagroBls,
        ),
        Err(Error::InvalidSignature)
    );
}

#[test]
fn rejects_a_padded_bitset() {
    let (fix, validators) = fixture();
    let mut padded = vec![0u8];
    padded.extend_from_slice(&unhex(&fix.bit_set));
    assert_eq!(
        verify_warp_aggregate(
            &validators,
            fix.total_weight,
            &unhex(&fix.unsigned_message),
            &padded,
            &unhex(&fix.signature).try_into().unwrap(),
            fix.quorum_num,
            fix.quorum_den,
            &MilagroBls,
        ),
        Err(Error::NonMinimalBitSet)
    );
}

#[test]
fn rejects_a_bit_beyond_the_canonical_set() {
    let (fix, validators) = fixture();
    // A high bit in a longer bitset points past the 73-entry set.
    let mut bit_set = unhex(&fix.bit_set);
    bit_set.insert(0, 0x80);
    let parsed = SignerBitSet::parse(&bit_set).unwrap();
    let out_of_range = bit_set.len() * 8 - 1;
    assert!(parsed.contains(out_of_range));
    assert_eq!(
        crate::verify_quorum(
            &validators,
            fix.total_weight,
            &parsed,
            fix.quorum_num,
            fix.quorum_den
        ),
        Err(Error::UnknownValidator(out_of_range))
    );
}

#[test]
fn rejects_when_dropping_a_signer_breaks_quorum() {
    let (fix, validators) = fixture();
    // Clear every set bit but the lowest-index signer: far below 67%.
    let bit_set = unhex(&fix.bit_set);
    let parsed = SignerBitSet::parse(&bit_set).unwrap();
    let first_signer = (0..validators.len()).find(|i| parsed.contains(*i)).unwrap();
    let mut lone = vec![0u8; bit_set.len()];
    let len = lone.len();
    lone[len - 1 - first_signer / 8] = 1 << (first_signer % 8);
    let minimal_start = lone.iter().position(|b| *b != 0).unwrap();
    let lone = &lone[minimal_start..];
    match crate::verify_quorum(
        &validators,
        fix.total_weight,
        &SignerBitSet::parse(lone).unwrap(),
        fix.quorum_num,
        fix.quorum_den,
    ) {
        Err(Error::InsufficientQuorum { signed, .. }) => {
            assert_eq!(signed, validators[first_signer].weight);
        }
        other => panic!("expected quorum failure, got {other:?}"),
    }
}

#[test]
fn keyless_weight_dilutes_the_quorum_via_total_weight() {
    // Synthetic 3-validator set: two keyed signers with weight 30 each, plus
    // 40 keyless weight present only in the total. 60/100 < 67/100.
    let validators = vec![
        Validator {
            public_key: [1; PUBLIC_KEY_LEN],
            weight: 30,
        },
        Validator {
            public_key: [2; PUBLIC_KEY_LEN],
            weight: 30,
        },
    ];
    let bit_set = [0b11u8];
    match crate::verify_quorum(
        &validators,
        100,
        &SignerBitSet::parse(&bit_set).unwrap(),
        67,
        100,
    ) {
        Err(Error::InsufficientQuorum { signed, total, .. }) => {
            assert_eq!((signed, total), (60, 100));
        }
        other => panic!("expected dilution to break quorum, got {other:?}"),
    }
    // Without the keyless weight the same signers clear it.
    assert!(crate::verify_quorum(
        &validators,
        60,
        &SignerBitSet::parse(&bit_set).unwrap(),
        67,
        100
    )
    .is_ok());
}
