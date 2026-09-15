//! This module contains the instantiate helper functions

use cosmwasm_std::{ensure, Storage};
use ethereum_light_client::{
    client_state::ClientState as EthClientState,
    consensus_state::ConsensusState as EthConsensusState,
};
use ibc_proto::ibc::{
    core::client::v1::Height as IbcProtoHeight,
    lightclients::wasm::v1::{
        ClientState as WasmClientState, ConsensusState as WasmConsensusState,
    },
};

use crate::{
    msg::InstantiateMsg,
    state::{consensus_db_key, encode_client_state, encode_consensus_state, HOST_CLIENT_STATE_KEY},
    ContractError,
};

/// Fully validated and encoded bootstrap writes.
pub struct PreparedClient {
    client: Vec<u8>,
    consensus_key: Vec<u8>,
    consensus: Vec<u8>,
}

impl PreparedClient {
    /// Commits the already validated bootstrap state without fallible work between writes.
    pub fn commit(self, storage: &mut dyn Storage) {
        storage.set(HOST_CLIENT_STATE_KEY.as_bytes(), &self.client);
        storage.set(&self.consensus_key, &self.consensus);
    }
}

/// Initializes the client state and consensus state
/// # Errors
/// Will return an error if the client state or consensus state cannot be deserialized.
/// # Panics
/// Will panic if the client state latest height cannot be unwrapped
#[allow(clippy::needless_pass_by_value)]
pub fn prepare_client(msg: InstantiateMsg) -> Result<PreparedClient, ContractError> {
    let client_state_bz: Vec<u8> = msg.client_state.into();
    let client_state: EthClientState = serde_json::from_slice(&client_state_bz)
        .map_err(ContractError::DeserializeClientStateFailed)?;
    let wasm_client_state = WasmClientState {
        checksum: msg.checksum.into(),
        data: client_state_bz,
        latest_height: Some(IbcProtoHeight {
            revision_number: 0,
            revision_height: client_state.latest_slot,
        }),
    };

    let consensus_state_bz: Vec<u8> = msg.consensus_state.into();
    let consensus_state: EthConsensusState = serde_json::from_slice(&consensus_state_bz)
        .map_err(ContractError::DeserializeConsensusStateFailed)?;
    let wasm_consensus_state = WasmConsensusState {
        data: consensus_state_bz,
    };

    ensure!(
        wasm_client_state.latest_height.unwrap().revision_height == client_state.latest_slot,
        ContractError::ClientStateSlotMismatch
    );

    ensure!(
        client_state.latest_slot == consensus_state.slot,
        ContractError::ClientAndConsensusStateMismatch
    );
    ensure!(client_state.latest_slot != 0, ContractError::ZeroHeight);
    ensure!(
        client_state.slots_per_epoch != 0,
        ContractError::InvalidClientState("slots_per_epoch must be non-zero")
    );
    ensure!(
        client_state.epochs_per_sync_committee_period != 0,
        ContractError::InvalidClientState("epochs_per_sync_committee_period must be non-zero")
    );

    client_state
        .verify_supported_fork_at_epoch(
            client_state.compute_epoch_at_slot(client_state.latest_slot),
        )
        .map_err(ContractError::UnsupportedForkVersion)?;

    Ok(PreparedClient {
        client: encode_client_state(&wasm_client_state)?,
        consensus_key: consensus_db_key(consensus_state.slot).into_bytes(),
        consensus: encode_consensus_state(&wasm_consensus_state)?,
    })
}
