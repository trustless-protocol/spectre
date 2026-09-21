//! Storage-backed client lifecycle: a specialization of the shared attested
//! runtime's transition discipline (validate and encode everything fallible
//! before the first storage write; first write wins; freeze on misbehaviour)
//! for the warp trust model.

use cosmwasm_std::{to_json_binary, Binary, Storage};
use ibc_proto::{
    google::protobuf::Any,
    ibc::{
        core::client::v1::Height,
        lightclients::wasm::v1::{
            ClientState as WasmClientState, ConsensusState as WasmConsensusState,
        },
    },
};
use l2_client::{
    msg::EvmStorageProof,
    packet,
    store::{
        consensus_db_key, get_wasm_client_state, get_wasm_consensus_state,
        may_get_wasm_consensus_state, HOST_CLIENT_STATE_KEY,
    },
};
use prost::Message;

use crate::{
    error::Error,
    msg::{IbcHeight, MigrateMsg, UpdateStateResult},
    state::{ClientState, ConsensusState},
    verification::AuthenticatedHeader,
};

/// Stores explicitly trusted bootstrap state.
pub fn instantiate(
    storage: &mut dyn Storage,
    client: &ClientState,
    consensus: &ConsensusState,
    checksum: Vec<u8>,
    accepted_at: u64,
) -> Result<(), Error> {
    client.validate()?;
    consensus.validate()?;
    if consensus.height != client.latest_height {
        return Err(Error::Invalid(
            "bootstrap consensus height differs from client height",
        ));
    }

    let height = client.latest_height;
    let mut consensus = consensus.clone();
    consensus.first_accepted_at = accepted_at;
    let client_envelope = Any::from_msg(&WasmClientState {
        checksum,
        data: serde_json::to_vec(client)?,
        latest_height: Some(Height {
            revision_number: 0,
            revision_height: height,
        }),
    })
    .map_err(Error::Proto)?
    .encode_to_vec();
    let consensus_envelope = Any::from_msg(&WasmConsensusState {
        data: serde_json::to_vec(&consensus)?,
    })
    .map_err(Error::Proto)?
    .encode_to_vec();

    storage.set(HOST_CLIENT_STATE_KEY.as_bytes(), &client_envelope);
    storage.set(consensus_db_key(height).as_bytes(), &consensus_envelope);
    Ok(())
}

/// Applies the only supported governance migration: keep or atomically replace
/// the pinned canonical validator set. A replacement must strictly advance the
/// snapshot's P-Chain height, so even governance cannot roll the set back.
pub fn migrate(storage: &mut dyn Storage, msg: MigrateMsg) -> Result<(), Error> {
    let mut client = client_state(storage)?;
    client.validate()?;
    let Some(replacement) = msg.replacement() else {
        return Ok(());
    };
    replacement.validate()?;
    if replacement.p_chain_height <= client.validator_set.p_chain_height {
        return Err(Error::NonMonotonicRotation {
            stored: client.validator_set.p_chain_height,
            proposed: replacement.p_chain_height,
        });
    }
    client.validator_set = replacement;
    client.validate()?;
    store_client_state(storage, &client)
}

/// Loads client state from the host envelope.
pub fn client_state(storage: &dyn Storage) -> Result<ClientState, Error> {
    let envelope = get_wasm_client_state(storage).map_err(Error::Evm)?;
    let client: ClientState = serde_json::from_slice(&envelope.data)?;
    let Some(height) = envelope.latest_height else {
        return Err(Error::Invalid("host envelope carries no latest height"));
    };
    if height.revision_number != 0 || height.revision_height != client.latest_height {
        return Err(Error::Invalid(
            "host envelope height differs from client state",
        ));
    }
    Ok(client)
}

/// Loads an exact consensus state.
pub fn consensus_state(storage: &dyn Storage, height: u64) -> Result<ConsensusState, Error> {
    let state: ConsensusState =
        serde_json::from_slice(&get_wasm_consensus_state(storage, height).map_err(Error::Evm)?.data)?;
    state.validate()?;
    if state.height != height {
        return Err(Error::Invalid("stored consensus height mismatch"));
    }
    Ok(state)
}

/// Applies a verified update and returns the encoded host response.
///
/// Semantics mirror the shared attested runtime: a conflicting same-height
/// submission is rejected, adjacent stored heights must chain by parent hash,
/// and an identical re-submission is a no-op (first write wins).
pub fn update(
    storage: &mut dyn Storage,
    header: &AuthenticatedHeader,
    accepted_at: u64,
) -> Result<Binary, Error> {
    let mut client = client_state(storage)?;
    client.validate()?;
    if let Some(height) = client.frozen_height {
        return Err(Error::Frozen(height));
    }

    let header = header.header();
    let height = header.height;
    let incoming = header.consensus_state(accepted_at)?;

    let existing = may_load_consensus_state(storage, height)?;
    let previous = if height > 1 {
        may_load_consensus_state(storage, height - 1)?
    } else {
        None
    };
    let next = height
        .checked_add(1)
        .map(|next_height| may_load_consensus_state(storage, next_height))
        .transpose()?
        .flatten();

    if existing
        .as_ref()
        .is_some_and(|stored| stored.conflicts_with(&incoming))
    {
        return Err(Error::Invalid("conflicting state already stored"));
    }
    if previous.is_some_and(|state| header.parent_hash != state.block_hash)
        || next.is_some_and(|state| state.parent_hash != header.block_hash)
    {
        return Err(Error::Invalid("adjacent parent hash mismatch"));
    }

    let response = to_json_binary(&UpdateStateResult {
        heights: vec![IbcHeight {
            revision_number: 0,
            revision_height: height,
        }],
    })?;
    if existing.is_some() {
        // A re-submission of an already-stored identical block: nothing to
        // write, and the stored state must not depend on arrival order.
        return Ok(response);
    }

    // Prepare every fallible encoding before the first storage write.
    client.latest_height = client.latest_height.max(height);
    let client_envelope = prepare_client_state_envelope(storage, &client)?;
    let consensus_envelope = Any::from_msg(&WasmConsensusState {
        data: serde_json::to_vec(&incoming)?,
    })
    .map_err(Error::Proto)?
    .encode_to_vec();

    storage.set(consensus_db_key(height).as_bytes(), &consensus_envelope);
    storage.set(HOST_CLIENT_STATE_KEY.as_bytes(), &client_envelope);
    Ok(response)
}

/// Returns whether two verified headers constitute actionable misbehaviour.
pub fn is_actionable_misbehaviour(
    first: &AuthenticatedHeader,
    second: &AuthenticatedHeader,
) -> Result<bool, Error> {
    first.header().conflicts_with(second.header())
}

/// Applies actionable misbehaviour evidence by freezing the client.
pub fn apply_misbehaviour(
    storage: &mut dyn Storage,
    first: &AuthenticatedHeader,
    second: &AuthenticatedHeader,
) -> Result<bool, Error> {
    let mut client = client_state(storage)?;
    client.validate()?;
    if let Some(height) = client.frozen_height {
        return Err(Error::Frozen(height));
    }
    if !is_actionable_misbehaviour(first, second)? {
        return Ok(false);
    }
    client.frozen_height = Some(first.header().height);
    store_client_state(storage, &client)?;
    Ok(true)
}

/// Verifies packet membership against an exact stored router storage root.
pub fn verify_membership(
    storage: &dyn Storage,
    height: u64,
    paths: &[Binary],
    value: &[u8],
    proof: &[u8],
) -> Result<(), Error> {
    let client = active_client(storage)?;
    let consensus = consensus_state(storage, height)?;
    let path = single_path(paths)?;
    let proof: EvmStorageProof = serde_json::from_slice(proof)?;
    packet::verify_membership(
        consensus.ibc_storage_root,
        client.commitment_slot,
        path,
        value,
        &proof,
    )
    .map_err(Error::Evm)
}

/// Verifies packet non-membership against an exact stored router storage root.
pub fn verify_non_membership(
    storage: &dyn Storage,
    height: u64,
    paths: &[Binary],
    proof: &[u8],
) -> Result<(), Error> {
    let client = active_client(storage)?;
    let consensus = consensus_state(storage, height)?;
    let path = single_path(paths)?;
    let proof: EvmStorageProof = serde_json::from_slice(proof)?;
    packet::verify_non_membership(
        consensus.ibc_storage_root,
        client.commitment_slot,
        path,
        &proof,
    )
    .map_err(Error::Evm)
}

fn single_path(paths: &[Binary]) -> Result<&Binary, Error> {
    let [path] = paths else {
        return Err(Error::Invalid("exactly one IBC commitment path is required"));
    };
    Ok(path)
}

fn active_client(storage: &dyn Storage) -> Result<ClientState, Error> {
    let client = client_state(storage)?;
    client.validate()?;
    if let Some(height) = client.frozen_height {
        return Err(Error::Frozen(height));
    }
    Ok(client)
}

fn may_load_consensus_state(
    storage: &dyn Storage,
    height: u64,
) -> Result<Option<ConsensusState>, Error> {
    may_get_wasm_consensus_state(storage, height)
        .map_err(Error::Evm)?
        .map(|state| {
            let state: ConsensusState = serde_json::from_slice(&state.data)?;
            state.validate()?;
            if state.height != height {
                return Err(Error::Invalid("stored consensus height mismatch"));
            }
            Ok(state)
        })
        .transpose()
}

fn store_client_state(storage: &mut dyn Storage, client: &ClientState) -> Result<(), Error> {
    let envelope = prepare_client_state_envelope(storage, client)?;
    storage.set(HOST_CLIENT_STATE_KEY.as_bytes(), &envelope);
    Ok(())
}

fn prepare_client_state_envelope(
    storage: &dyn Storage,
    client: &ClientState,
) -> Result<Vec<u8>, Error> {
    let mut envelope = get_wasm_client_state(storage).map_err(Error::Evm)?;
    envelope.data = serde_json::to_vec(client)?;
    envelope.latest_height = Some(Height {
        revision_number: 0,
        revision_height: client.latest_height,
    });
    Ok(Any::from_msg(&envelope).map_err(Error::Proto)?.encode_to_vec())
}
