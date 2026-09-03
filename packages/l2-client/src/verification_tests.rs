use alloy_primitives::{Address, Bloom, Bytes, FixedBytes, B256, U256};
use cosmwasm_std::testing::MockApi;

use crate::{
    attestation::AttestationHead,
    canonical_header::{CanonicalEvmHeader, ExecutionHeaderFork},
    error::Error,
    msg::{AttestedL2Header, EvmAccountProof},
    state::{CommonProfile, RuntimeProfile},
    verification::verify_attested_header,
};

struct Profile(CommonProfile);

impl RuntimeProfile for Profile {
    fn common(&self) -> &CommonProfile {
        &self.0
    }

    fn expected_profile_version() -> &'static str {
        "signature_test_v2"
    }
}

fn profile() -> Profile {
    Profile(CommonProfile {
        l2_chain_id: 10,
        l2_router: Address::with_last_byte(1),
        commitment_slot: B256::with_last_byte(2),
        profile_version: "signature_test_v2".into(),
        l2_header_fork: ExecutionHeaderFork::London,
        attestor_public_key: B256::with_last_byte(3),
        attestation_head: AttestationHead::Safe,
    })
}

fn unsigned_header() -> AttestedL2Header {
    AttestedL2Header {
        l2_header: CanonicalEvmHeader {
            parent_hash: B256::with_last_byte(1),
            ommers_hash: B256::with_last_byte(2),
            beneficiary: Address::with_last_byte(3),
            state_root: B256::with_last_byte(4),
            transactions_root: B256::with_last_byte(5),
            receipts_root: B256::with_last_byte(6),
            logs_bloom: Bloom::ZERO,
            difficulty: U256::ZERO,
            number: 7,
            gas_limit: 30_000_000,
            gas_used: 21_000,
            timestamp: 1_700_000_000,
            extra_data: Bytes::new(),
            mix_hash: B256::with_last_byte(7),
            nonce: FixedBytes::ZERO,
            base_fee_per_gas: Some(U256::from(9)),
            withdrawals_root: None,
            blob_gas_used: None,
            excess_blob_gas: None,
            parent_beacon_block_root: None,
            requests_hash: None,
        },
        router_proof: EvmAccountProof { proof: vec![] },
        attestor_signature: vec![],
    }
}

#[test]
fn rejects_an_unsigned_header_before_any_router_proof_is_accepted() {
    assert!(matches!(
        verify_attested_header(&MockApi::default(), &profile(), &unsigned_header()),
        Err(Error::InvalidAttestorSignature)
    ));
}
