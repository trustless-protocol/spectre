//! End-to-end verification tests over REAL Fuji data.
//!
//! Three live captures pin the whole pipeline: the warp fixture (verified 67%
//! aggregate + the exact canonical validator set, captured by
//! `tools/warp-spike -dump-fixture`), the full coreth header of the attested
//! block (its recomputed hash must equal the signed hash), and an
//! `eth_getProof` account proof against that header's state root. Real BLS
//! runs in-process through milagro (test-only, same ciphersuite as warp).

use alloy_primitives::B256;
use cosmwasm_std::Binary;
use serde::Deserialize;

use crate::{
    msg::{EvmAccountProof, WarpSignedHeader},
    state::{ClientState, ValidatorSetState, WireValidator},
    verification::verify_signed_header,
};

#[derive(Deserialize)]
struct WarpFixtureValidator {
    public_key_compressed: String,
    weight: u64,
}

#[derive(Deserialize)]
struct WarpFixture {
    network_id: u32,
    source_chain_id: String,
    block_hash: String,
    bit_set: String,
    signature: String,
    total_weight: u64,
    validators: Vec<WarpFixtureValidator>,
}

#[derive(Deserialize)]
struct AccountProofFixture {
    address: String,
    storage_hash: String,
    account_proof: EvmAccountProof,
}

fn unhex(s: &str) -> Vec<u8> {
    hex::decode(s.trim_start_matches("0x")).unwrap()
}

pub(crate) struct MilagroBls;

impl avalanche_warp::BlsVerify for MilagroBls {
    fn fast_aggregate_verify(
        &self,
        public_keys: &[&[u8; avalanche_warp::PUBLIC_KEY_LEN]],
        message: &[u8],
        signature: &[u8; avalanche_warp::SIGNATURE_LEN],
    ) -> Result<(), avalanche_warp::Error> {
        let keys = public_keys
            .iter()
            .map(|pk| milagro_bls::PublicKey::from_bytes(pk.as_slice()))
            .collect::<Result<Vec<_>, _>>()
            .map_err(|_| avalanche_warp::Error::InvalidSignature)?;
        let key_refs: Vec<&milagro_bls::PublicKey> = keys.iter().collect();
        let aggregate_key = milagro_bls::AggregatePublicKey::aggregate(&key_refs)
            .map_err(|_| avalanche_warp::Error::InvalidSignature)?;
        let signature = milagro_bls::Signature::from_bytes(signature.as_slice())
            .map_err(|_| avalanche_warp::Error::InvalidSignature)?;
        let aggregate_signature = milagro_bls::AggregateSignature::aggregate(&[&signature]);
        if aggregate_signature.fast_aggregate_verify_pre_aggregated(message, &aggregate_key) {
            Ok(())
        } else {
            Err(avalanche_warp::Error::InvalidSignature)
        }
    }
}

fn warp_fixture() -> WarpFixture {
    serde_json::from_str(include_str!(
        "../../avalanche-warp/testdata/fuji_warp_fixture.json"
    ))
    .unwrap()
}

fn account_proof_fixture() -> AccountProofFixture {
    serde_json::from_str(include_str!(
        "../testdata/fuji_account_proof_58534295.json"
    ))
    .unwrap()
}

pub(crate) fn fuji_client_state() -> ClientState {
    let fix = warp_fixture();
    let proof = account_proof_fixture();
    ClientState {
        latest_height: 1,
        frozen_height: None,
        network_id: fix.network_id,
        source_chain_id: Binary::from(unhex(&fix.source_chain_id)),
        evm_chain_id: 43113,
        router: proof.address.parse().unwrap(),
        commitment_slot: B256::with_last_byte(1),
        quorum_num: 67,
        quorum_den: 100,
        validator_set: ValidatorSetState {
            validators: fix
                .validators
                .iter()
                .map(|v| WireValidator {
                    public_key: Binary::from(unhex(&v.public_key_compressed)),
                    weight: v.weight,
                })
                .collect(),
            total_weight: fix.total_weight,
            p_chain_height: 100,
        },
    }
}

pub(crate) fn fuji_signed_header() -> WarpSignedHeader {
    let fix = warp_fixture();
    WarpSignedHeader {
        header: serde_json::from_str(include_str!(
            "../testdata/fuji_wire_header_58534295.json"
        ))
        .unwrap(),
        signer_bit_set: Binary::from(unhex(&fix.bit_set)),
        signature: Binary::from(unhex(&fix.signature)),
        router_proof: account_proof_fixture().account_proof,
    }
}

#[test]
fn verifies_a_real_fuji_header_end_to_end() {
    let client = fuji_client_state();
    let fix = warp_fixture();
    let proof = account_proof_fixture();

    let authenticated = verify_signed_header(&client, &fuji_signed_header(), &MilagroBls).unwrap();
    let header = authenticated.header();
    assert_eq!(header.height, 58_534_295);
    assert_eq!(header.block_hash.to_vec(), unhex(&fix.block_hash));
    assert_eq!(
        header.router_storage_root.to_vec(),
        unhex(&proof.storage_hash)
    );
    header.consensus_state(7).unwrap();
}

#[test]
fn rejects_the_wrong_network_id() {
    // A different network id changes the signed message: the same aggregate
    // must not verify (warp's replay domain).
    let mut client = fuji_client_state();
    client.network_id = 1;
    assert!(matches!(
        verify_signed_header(&client, &fuji_signed_header(), &MilagroBls),
        Err(crate::error::Error::Warp(
            avalanche_warp::Error::InvalidSignature
        ))
    ));
}

#[test]
fn rejects_a_quorum_the_signers_do_not_reach() {
    let mut client = fuji_client_state();
    client.quorum_num = 90;
    assert!(matches!(
        verify_signed_header(&client, &fuji_signed_header(), &MilagroBls),
        Err(crate::error::Error::Warp(
            avalanche_warp::Error::InsufficientQuorum { .. }
        ))
    ));
}

#[test]
fn rejects_a_tampered_header_field() {
    // Changing any header field changes the recomputed coreth hash, so the
    // quorum's signature over the true hash must not verify.
    let mut signed = fuji_signed_header();
    signed.header.gas_used += 1;
    assert!(matches!(
        verify_signed_header(&fuji_client_state(), &signed, &MilagroBls),
        Err(crate::error::Error::Warp(
            avalanche_warp::Error::InvalidSignature
        ))
    ));
}

#[test]
fn rejects_a_malformed_signature_length() {
    let mut signed = fuji_signed_header();
    signed.signature = Binary::from(vec![0_u8; 95]);
    assert!(matches!(
        verify_signed_header(&fuji_client_state(), &signed, &MilagroBls),
        Err(crate::error::Error::Invalid(_))
    ));
}

#[test]
fn rejects_a_router_the_proof_does_not_cover() {
    let mut client = fuji_client_state();
    client.router = "0x00000000000000000000000000000000000000aa".parse().unwrap();
    assert!(verify_signed_header(&client, &fuji_signed_header(), &MilagroBls).is_err());
}
