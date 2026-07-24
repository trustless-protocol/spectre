//! Tests for 08-wasm host storage envelopes.

use cosmwasm_std::testing::MockStorage;
use ibc_proto::ibc::lightclients::wasm::v1::{
    ClientState as WasmClientState, ConsensusState as WasmConsensusState,
};

use crate::store::{
    consensus_db_key, get_wasm_client_state, get_wasm_consensus_state, store_wasm_client_state,
    store_wasm_consensus_state,
};

#[test]
fn stores_and_loads_revision_zero_envelopes() {
    let mut storage = MockStorage::new();
    let client = WasmClientState {
        data: b"client".to_vec(),
        checksum: b"checksum".to_vec(),
        latest_height: None,
    };
    let consensus = WasmConsensusState {
        data: b"consensus".to_vec(),
    };

    store_wasm_client_state(&mut storage, &client).unwrap();
    store_wasm_consensus_state(&mut storage, 42, &consensus).unwrap();

    assert_eq!(get_wasm_client_state(&storage).unwrap(), client);
    assert_eq!(get_wasm_consensus_state(&storage, 42).unwrap(), consensus);
    assert_eq!(consensus_db_key(42), "consensusStates/0-42");
}
