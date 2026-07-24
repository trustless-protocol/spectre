use crate::{error::Error, fixtures::load_fixture};

const FIXTURE: &[u8] = br#"{"height":42}"#;
const DIGEST: &str = "70b6a06775faec6fb8179a3ba5a9d6a5a698ef7504ffb2e1dbbc15081aab9b81";

#[test]
fn loads_a_labeled_fixture() {
    let metadata = format!(
        "{{\"schema_version\":2,\"config_id\":\"base-sepolia\",\"l1_block_number\":1,\"l1_block_hash\":\"0x{}\",\"l2_block_number\":42,\"l2_block_hash\":\"0x{}\",\"l1_beacon_slot\":1,\"ethereum_source_revision\":\"ethereum-client-v1\",\"rollup_source_revision\":\"op-contracts/v1\",\"profile_sha256\":\"{}\",\"fixture_sha256\":\"{DIGEST}\"}}",
        "11".repeat(32), "22".repeat(32), "33".repeat(32),
    );
    let (value, loaded): (serde_json::Value, _) =
        load_fixture(FIXTURE, metadata.as_bytes()).unwrap();
    assert_eq!(value["height"], 42);
    assert_eq!(loaded.config_id, "base-sepolia");
}

#[test]
fn rejects_incomplete_fixture_metadata() {
    let metadata = format!(
        "{{\"schema_version\":2,\"config_id\":\"\",\"l1_block_number\":1,\"l1_block_hash\":\"0x{}\",\"l2_block_number\":42,\"l2_block_hash\":\"0x{}\",\"l1_beacon_slot\":0,\"ethereum_source_revision\":\"eth\",\"rollup_source_revision\":\"rollup\",\"profile_sha256\":\"{}\",\"fixture_sha256\":\"{DIGEST}\"}}",
        "11".repeat(32), "22".repeat(32), "33".repeat(32),
    );
    assert!(matches!(
        load_fixture::<serde_json::Value>(FIXTURE, metadata.as_bytes()),
        Err(Error::InvalidFixtureProvenance(_))
    ));
}
