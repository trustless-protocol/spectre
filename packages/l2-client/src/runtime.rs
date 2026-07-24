//! Storage-backed operations for the L2 client lifecycle.

use cosmwasm_std::Storage;
use ibc_proto::ibc::{
    core::client::v1::Height,
    lightclients::wasm::v1::{
        ClientState as WasmClientState, ConsensusState as WasmConsensusState,
    },
};
use serde::{de::DeserializeOwned, Serialize};

use crate::{
    error::Error,
    misbehaviour,
    msg::EvmStorageProof,
    packet,
    state::{ClientState, ConsensusState, Header, RuntimeProfile},
    store::{
        get_wasm_client_state, get_wasm_consensus_state, may_get_wasm_consensus_state,
        store_wasm_client_state, store_wasm_consensus_state,
    },
    sudo,
};

/// Stores directly supplied client and consensus state in the 08-wasm host layout.
pub fn instantiate<Profile>(
    storage: &mut dyn Storage,
    client: &ClientState<Profile>,
    consensus: &ConsensusState,
    checksum: Vec<u8>,
) -> Result<(), Error>
where
    Profile: Serialize,
{
    let height = client.latest_height;
    store_wasm_client_state(
        storage,
        &WasmClientState {
            checksum,
            data: serde_json::to_vec(client)?,
            latest_height: Some(Height {
                revision_number: 0,
                revision_height: height,
            }),
        },
    )?;
    store_consensus(storage, height, consensus)
}

/// Loads the compact active state from host-owned storage.
pub fn client_state<Config: DeserializeOwned>(
    storage: &dyn Storage,
) -> Result<ClientState<Config>, Error> {
    Ok(serde_json::from_slice(
        &get_wasm_client_state(storage)?.data,
    )?)
}

/// Loads the consensus state at exactly `height`.
pub fn consensus_state(storage: &dyn Storage, height: u64) -> Result<ConsensusState, Error> {
    Ok(serde_json::from_slice(
        &get_wasm_consensus_state(storage, height)?.data,
    )?)
}

/// Applies a chain-verified header and returns the updated height when state advanced.
pub fn update<Profile>(storage: &mut dyn Storage, header: &Header) -> Result<Option<u64>, Error>
where
    Profile: DeserializeOwned + Serialize,
{
    let mut client = client_state::<Profile>(storage)?;
    let height = header.height.revision_height;
    let existing = may_get_wasm_consensus_state(storage, height)?
        .map(|state| serde_json::from_slice(&state.data))
        .transpose()?;
    let old_latest = client.latest_height;
    let result = sudo::update_state(&mut client, header, existing.as_ref())?;
    let Some(consensus) = result else {
        return Ok(None);
    };
    store_consensus(storage, height, &consensus)?;
    if client.latest_height != old_latest {
        store_client(storage, &client)?;
    }
    Ok(Some(height))
}

/// Detects valid same-height conflict evidence and persists the frozen state.
pub fn apply_misbehaviour<Config>(
    storage: &mut dyn Storage,
    first: &Header,
    second: &Header,
) -> Result<bool, Error>
where
    Config: DeserializeOwned + Serialize,
{
    let mut client = client_state::<Config>(storage)?;
    let found = misbehaviour::apply(&mut client, first, second)?;
    if found {
        store_client(storage, &client)?;
    }
    Ok(found)
}

/// Verifies a router commitment at an exact stored height.
pub fn verify_membership<Profile: DeserializeOwned + RuntimeProfile>(
    storage: &dyn Storage,
    height: u64,
    path: &[u8],
    value: &[u8],
    proof: &EvmStorageProof,
) -> Result<(), Error> {
    let client = client_state::<Profile>(storage)?;
    if let Some(frozen) = client.frozen_height {
        return Err(Error::Frozen(frozen));
    }
    packet::verify_membership(
        consensus_state(storage, height)?.ibc_storage_root,
        client.profile.common().commitment_slot,
        path,
        value,
        proof,
    )
}

/// Verifies router commitment absence at an exact stored height.
pub fn verify_non_membership<Profile: DeserializeOwned + RuntimeProfile>(
    storage: &dyn Storage,
    height: u64,
    path: &[u8],
    proof: &EvmStorageProof,
) -> Result<(), Error> {
    let client = client_state::<Profile>(storage)?;
    if let Some(frozen) = client.frozen_height {
        return Err(Error::Frozen(frozen));
    }
    packet::verify_non_membership(
        consensus_state(storage, height)?.ibc_storage_root,
        client.profile.common().commitment_slot,
        path,
        proof,
    )
}

fn store_consensus(
    storage: &mut dyn Storage,
    height: u64,
    consensus: &ConsensusState,
) -> Result<(), Error> {
    store_wasm_consensus_state(
        storage,
        height,
        &WasmConsensusState {
            data: serde_json::to_vec(consensus)?,
        },
    )
}

fn store_client<Config: Serialize>(
    storage: &mut dyn Storage,
    client: &ClientState<Config>,
) -> Result<(), Error> {
    let mut envelope = get_wasm_client_state(storage)?;
    envelope.data = serde_json::to_vec(client)?;
    envelope.latest_height = Some(Height {
        revision_number: 0,
        revision_height: client.latest_height,
    });
    store_wasm_client_state(storage, &envelope)
}
