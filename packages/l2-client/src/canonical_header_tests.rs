//! Tests for strict canonical L2 header validation and hashing.

use alloy_primitives::{Address, Bloom, Bytes, FixedBytes, B256, U256};

use crate::{
    canonical_header::{CanonicalEvmHeader, ExecutionHeaderFork},
    error::Error,
};

fn header() -> CanonicalEvmHeader {
    CanonicalEvmHeader {
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
        mix_hash: B256::with_last_byte(8),
        nonce: FixedBytes::ZERO,
        base_fee_per_gas: Some(U256::from(9)),
        withdrawals_root: Some(B256::with_last_byte(10)),
        blob_gas_used: Some(11),
        excess_blob_gas: Some(12),
        parent_beacon_block_root: Some(B256::with_last_byte(13)),
        requests_hash: Some(B256::with_last_byte(14)),
        ..CanonicalEvmHeader::default()
    }
}

#[test]
fn hashes_the_canonical_prague_field_order_deterministically() {
    let header = header();
    let rlp = header.rlp_bytes(ExecutionHeaderFork::Prague).unwrap();
    assert!(rlp.first().is_some_and(|prefix| prefix >= &0xf8));
    // This vector is deliberately fixed so accidental changes to RLP ordering
    // (including optional post-Cancun fields) cannot silently self-confirm.
    assert_eq!(
        header.hash(ExecutionHeaderFork::Prague).unwrap(),
        "0x22ffbd4be33e618a3687c5b7dd88e333b466a0dd8d197528970c4a287ec42f6e"
            .parse::<B256>()
            .unwrap()
    );
}

#[test]
fn rejects_partial_cancun_fields_and_oversized_extra_data() {
    let mut incomplete = header();
    incomplete.excess_blob_gas = None;
    assert!(matches!(
        incomplete.validate_for_fork(ExecutionHeaderFork::Cancun),
        Err(Error::InvalidEvmHeader(
            "optional fields do not match the pinned fork"
        ))
    ));

    let mut oversized = header();
    oversized.extra_data = Bytes::from(vec![0; CanonicalEvmHeader::MAX_EXTRA_DATA_BYTES + 1]);
    assert!(matches!(
        oversized.validate_for_fork(ExecutionHeaderFork::Prague),
        Err(Error::InvalidEvmHeader("extra_data exceeds 32 bytes"))
    ));
}

#[test]
fn rejects_fields_from_a_later_fork() {
    assert!(matches!(
        header().validate_for_fork(ExecutionHeaderFork::London),
        Err(Error::InvalidEvmHeader(
            "optional fields do not match the pinned fork"
        ))
    ));
}

// Real Fuji C-Chain header #58513648 (post-Helicon: full optional tail),
// fetched 2026-09-21 from https://api.avax-test.network/ext/bc/C/rpc.
fn fuji_coreth_header() -> CanonicalEvmHeader {
    CanonicalEvmHeader {
        parent_hash: "0xb325f8030812940061ca0ee57d205a7744fa45af4fe477936a76527fdad9e685".parse().unwrap(),
        ommers_hash: "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347".parse().unwrap(),
        beneficiary: "0x0100000000000000000000000000000000000000".parse().unwrap(),
        state_root: "0xc614444ff2a5113b441c004426baaba133e948a4a72ede8f592e5f994381170f".parse().unwrap(),
        transactions_root: "0xc26bd92d087df0e6adb5e88a14b3804f6e98af13f134d0243188af5948ea7f49".parse().unwrap(),
        receipts_root: "0x35c453ac203b6659d172bcf79bd993d7dd2d2a5a178b2395b0b0f271e69ef3e6".parse().unwrap(),
        logs_bloom: "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000080000000000040000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000".parse().unwrap(),
        difficulty: U256::from(1u64),
        number: 58_513_648,
        gas_limit: 32_000_000,
        gas_used: 54_898,
        timestamp: 1_789_932_885,
        extra_data: "0x000000000000".parse().unwrap(),
        mix_hash: B256::ZERO,
        nonce: FixedBytes::ZERO,
        base_fee_per_gas: Some(U256::from(10u64)),
        blob_gas_used: Some(0),
        excess_blob_gas: Some(0),
        parent_beacon_block_root: Some(B256::ZERO),
        ext_data_hash: Some(
            "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421".parse().unwrap(),
        ),
        ext_data_gas_used: Some(U256::ZERO),
        block_gas_cost: Some(U256::ZERO),
        time_milliseconds: Some(1_789_932_885_320),
        min_delay_excess: Some(7_132_828),
        target_exponent: Some(15_770_705),
        min_price_exponent: Some(957_480_584_338_323_632),
        settled_height: Some(58_513_645),
        settled_gas_unix: Some(1_789_932_879),
        settled_gas_numerator: Some(2_052_024),
        settled_excess: Some(320_574_659),
        ..CanonicalEvmHeader::default()
    }
}

#[test]
fn hashes_a_real_fuji_coreth_header_to_its_reported_block_hash() {
    assert_eq!(
        fuji_coreth_header().hash(ExecutionHeaderFork::Coreth).unwrap(),
        "0x68c074dadf040a70a81825dfdce0581bc604beca6cf8778a6f99c4a5acfb74e5"
            .parse::<B256>()
            .unwrap()
    );
}

#[test]
fn hashes_a_granite_era_coreth_header_without_the_helicon_tail() {
    // Same header with the Helicon fields absent: the cascade must stop at
    // min_delay_excess and still hash deterministically (mainnet-Granite shape).
    let mut granite = fuji_coreth_header();
    granite.target_exponent = None;
    granite.min_price_exponent = None;
    granite.settled_height = None;
    granite.settled_gas_unix = None;
    granite.settled_gas_numerator = None;
    granite.settled_excess = None;
    let rlp = granite.rlp_bytes(ExecutionHeaderFork::Coreth).unwrap();
    let helicon_rlp = fuji_coreth_header().rlp_bytes(ExecutionHeaderFork::Coreth).unwrap();
    assert!(rlp.len() < helicon_rlp.len());
    // No fixed vector exists for this synthetic shape; shortening the tail must
    // still change the hash, which guards against the cascade padding trailing
    // absents instead of omitting them.
    assert_ne!(
        granite.hash(ExecutionHeaderFork::Coreth).unwrap(),
        fuji_coreth_header().hash(ExecutionHeaderFork::Coreth).unwrap()
    );
}

#[test]
fn rejects_coreth_fields_on_ethereum_forks_and_vice_versa() {
    let mut op_with_coreth = header();
    op_with_coreth.ext_data_hash = Some(B256::with_last_byte(15));
    assert!(matches!(
        op_with_coreth.validate_for_fork(ExecutionHeaderFork::Prague),
        Err(Error::InvalidEvmHeader(
            "optional fields do not match the pinned fork"
        ))
    ));

    let mut coreth_with_withdrawals = fuji_coreth_header();
    coreth_with_withdrawals.withdrawals_root = Some(B256::with_last_byte(16));
    assert!(matches!(
        coreth_with_withdrawals.validate_for_fork(ExecutionHeaderFork::Coreth),
        Err(Error::InvalidEvmHeader(
            "optional fields do not match the pinned fork"
        ))
    ));

    let mut missing_ext_data_hash = fuji_coreth_header();
    missing_ext_data_hash.ext_data_hash = None;
    assert!(missing_ext_data_hash
        .validate_for_fork(ExecutionHeaderFork::Coreth)
        .is_err());
}
