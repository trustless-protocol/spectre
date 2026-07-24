use std::{fs, path::Path};

use crate::config::Config;
use l2_client::canonical_header::ExecutionHeaderFork;

#[test]
fn tooling_profile_round_trips_without_being_compiled_into_wasm() {
    let path = Path::new(env!("CARGO_MANIFEST_DIR")).join("config/base-sepolia.json");
    let bytes = fs::read(path).unwrap();
    let profile: Config = serde_json::from_slice(&bytes).unwrap();
    assert_eq!(profile.common.l2_chain_id, 84_532);
    assert_eq!(profile.l2_header_fork, ExecutionHeaderFork::Prague);
    assert_eq!(
        serde_json::from_slice::<Config>(&serde_json::to_vec(&profile).unwrap()).unwrap(),
        profile
    );
}
