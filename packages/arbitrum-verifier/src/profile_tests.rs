use std::{fs, path::Path};

use crate::config::{Config, RollupProtocol};

#[test]
fn tooling_profile_round_trips_without_being_compiled_into_wasm() {
    let path = Path::new(env!("CARGO_MANIFEST_DIR")).join("config/arbitrum-sepolia.json");
    let bytes = fs::read(path).unwrap();
    let profile: Config = serde_json::from_slice(&bytes).unwrap();
    assert_eq!(profile.common.l2_chain_id, 421_614);
    assert!(matches!(&profile.protocol, RollupProtocol::LegacyNitro(_)));
    assert_eq!(
        serde_json::from_slice::<Config>(&serde_json::to_vec(&profile).unwrap()).unwrap(),
        profile
    );
}

#[test]
fn bold_profile_remains_a_supported_tagged_variant() {
    let profile: Config = serde_json::from_str(
        r#"{
            "common": {
                "l1_chain_id": 1,
                "l2_chain_id": 42161,
                "ethereum_client": {
                    "client_id": "08-wasm-ethereum",
                    "wasm_checksum": []
                },
                "l2_router": "0x0000000000000000000000000000000000000001",
                "commitment_slot": "0x0000000000000000000000000000000000000000000000000000000000000002",
                "rollup_version": "arbitrum-bold-v2"
            },
            "rollup": "0x0000000000000000000000000000000000000003",
            "protocol": {
                "type": "bold_v2",
                "value": {
                    "assertions_mapping_slot": "0x0000000000000000000000000000000000000000000000000000000000000076",
                    "assertion_status_offset": 25,
                    "version": "bold-v2"
                }
            }
        }"#,
    )
    .unwrap();

    assert!(matches!(profile.protocol, RollupProtocol::BoldV2(_)));
}
