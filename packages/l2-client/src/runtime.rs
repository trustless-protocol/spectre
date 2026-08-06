//! Storage-backed L2 client lifecycle.

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
    msg::EvmStorageProof,
    packet,
    state::{ClientState, ConsensusState, Header, RuntimeProfile},
    store::{
        get_wasm_client_state, get_wasm_consensus_state, may_get_wasm_consensus_state,
        store_wasm_client_state, store_wasm_consensus_state,
    },
};

/// Stores explicitly trusted bootstrap state.
pub fn instantiate<Profile>(
    storage: &mut dyn Storage,
    client: &ClientState<Profile>,
    consensus: &ConsensusState,
    checksum: Vec<u8>,
    accepted_at: u64,
) -> Result<(), Error>
where
    Profile: Clone + Serialize + RuntimeProfile,
{
    validate_profile::<Profile>(&client.profile)?;
    consensus.validate()?;
    if consensus.l2_height != client.latest_height {
        return Err(Error::InvalidHeader(
            "bootstrap consensus height differs from client height",
        ));
    }

    let mut consensus = consensus.clone();
    consensus.first_accepted_at = accepted_at;
    store_wasm_client_state(
        storage,
        &WasmClientState {
            checksum,
            data: serde_json::to_vec(client)?,
            latest_height: Some(Height {
                revision_number: 0,
                revision_height: client.latest_height,
            }),
        },
    )?;
    store_consensus_state(storage, client.latest_height, &consensus)
}

/// Loads client state from the host envelope.
pub fn client_state<Profile: DeserializeOwned>(
    storage: &dyn Storage,
) -> Result<ClientState<Profile>, Error> {
    Ok(serde_json::from_slice(
        &get_wasm_client_state(storage)?.data,
    )?)
}

/// Loads an exact consensus state.
pub fn consensus_state(storage: &dyn Storage, height: u64) -> Result<ConsensusState, Error> {
    let state: ConsensusState =
        serde_json::from_slice(&get_wasm_consensus_state(storage, height)?.data)?;
    state.validate()?;
    Ok(state)
}

/// Applies a verified update and persists any resulting state transition.
pub fn update<Profile>(
    storage: &mut dyn Storage,
    header: &Header,
    accepted_at: u64,
) -> Result<Option<u64>, Error>
where
    Profile: DeserializeOwned + Serialize + RuntimeProfile,
{
    let mut client = client_state::<Profile>(storage)?;
    validate_profile::<Profile>(&client.profile)?;
    let height = header.height.revision_height;
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

    let transition = apply_update(
        &mut client,
        header,
        existing.as_ref(),
        previous.as_ref(),
        next.as_ref(),
        accepted_at,
    );

    let consensus = match transition {
        Ok(consensus) => consensus,
        // Once signed attestations are implemented, a conflict freezes the client and the freeze
        // must be persisted, so it cannot be reported as an Err — the transaction would revert
        // and discard it. Until then, propagate the conflict error and write nothing.
        Err(Error::StateConflict { .. }) if ATTESTATIONS_ARE_AUTHENTICATED => {
            store_client_state(storage, &client)?;
            return Ok(None);
        }
        Err(error) => return Err(error),
    };

    let Some(consensus) = consensus else {
        return Ok(None);
    };
    store_consensus_state(storage, height, &consensus)?;
    store_client_state(storage, &client)?;
    Ok(Some(height))
}

/// Whether every accepted header carries an attestor signature verified by this client.
///
/// False until the signed-attestation wire format and signature verification are implemented.
/// The current envelope carries no attestation, so anyone can manufacture two contradictory
/// headers and submit them through permissionless `MsgUpdateClient` calls. Freezing on those bare
/// headers would let any party permanently brick a client, and there is no unfreeze path.
///
/// The redesign defines misbehaviour as conflicting signed attestations. Set this to true only in
/// the same change that verifies those signatures and makes conflicts attributable.
const ATTESTATIONS_ARE_AUTHENTICATED: bool = false;

/// Returns whether two headers constitute actionable, authenticated misbehaviour.
pub fn is_actionable_misbehaviour(first: &Header, second: &Header) -> Result<bool, Error> {
    let found = first.conflicts_with(second)?;
    if found && !ATTESTATIONS_ARE_AUTHENTICATED {
        return Err(Error::InvalidHeader(
            "conflicting headers do not carry authenticated attestations",
        ));
    }
    Ok(found)
}

/// Applies actionable misbehaviour evidence.
pub fn apply_misbehaviour<Profile>(
    storage: &mut dyn Storage,
    first: &Header,
    second: &Header,
) -> Result<bool, Error>
where
    Profile: DeserializeOwned + Serialize,
{
    let mut client = client_state::<Profile>(storage)?;
    if let Some(height) = client.frozen_height {
        return Err(Error::Frozen(height));
    }
    let found = is_actionable_misbehaviour(first, second)?;
    if !found {
        return Ok(false);
    }
    client.frozen_height = Some(first.height.revision_height);
    store_client_state(storage, &client)?;
    Ok(true)
}

/// Verifies packet membership against an exact stored root.
pub fn verify_membership<Profile: DeserializeOwned + RuntimeProfile>(
    storage: &dyn Storage,
    height: u64,
    path: &[u8],
    value: &[u8],
    proof: &EvmStorageProof,
) -> Result<(), Error> {
    let client = active_client::<Profile>(storage)?;
    let consensus = consensus_state(storage, height)?;
    packet::verify_membership(
        consensus.ibc_storage_root,
        client.profile.common().commitment_slot,
        path,
        value,
        proof,
    )
}

/// Verifies packet non-membership against an exact stored root.
pub fn verify_non_membership<Profile: DeserializeOwned + RuntimeProfile>(
    storage: &dyn Storage,
    height: u64,
    path: &[u8],
    proof: &EvmStorageProof,
) -> Result<(), Error> {
    let client = active_client::<Profile>(storage)?;
    let consensus = consensus_state(storage, height)?;
    packet::verify_non_membership(
        consensus.ibc_storage_root,
        client.profile.common().commitment_slot,
        path,
        proof,
    )
}

fn active_client<Profile: DeserializeOwned>(
    storage: &dyn Storage,
) -> Result<ClientState<Profile>, Error> {
    let client = client_state(storage)?;
    if let Some(height) = client.frozen_height {
        return Err(Error::Frozen(height));
    }
    Ok(client)
}

fn apply_update<Profile>(
    client: &mut ClientState<Profile>,
    header: &Header,
    existing: Option<&ConsensusState>,
    previous: Option<&ConsensusState>,
    next: Option<&ConsensusState>,
    accepted_at: u64,
) -> Result<Option<ConsensusState>, Error> {
    if let Some(height) = client.frozen_height {
        return Err(Error::Frozen(height));
    }

    let incoming = header.consensus_state(accepted_at)?;
    let height = header.height.revision_height;

    if existing.is_some_and(|stored| stored.conflicts_with(&incoming)) {
        if ATTESTATIONS_ARE_AUTHENTICATED {
            client.frozen_height = Some(height);
        }
        return Err(Error::StateConflict { height });
    }

    if previous.is_some_and(|state| header.parent_hash != state.l2_block_hash)
        || next.is_some_and(|state| state.parent_hash != header.l2_block_hash)
    {
        return Err(Error::ParentHashMismatch);
    }

    if existing.is_some() {
        // conflicts_with() was false, so every persisted block-identity field matches: this is a
        // re-submission of a block already stored, with nothing to update. First write wins, so
        // the stored state does not depend on arrival order. An honest relayer re-sends the same
        // block routinely (a retry or a cached header), and neither should rewrite anything.
        //
        // The signed-attestation wire format grows this branch into the redesign's promotion and
        // provisional-replacement state machine.
        return Ok(None);
    }

    client.latest_height = client.latest_height.max(height);
    Ok(Some(incoming))
}

fn validate_profile<Profile: RuntimeProfile>(profile: &Profile) -> Result<(), Error> {
    if profile.common().profile_version != Profile::expected_profile_version() {
        return Err(Error::InvalidProfileVersion);
    }
    Ok(())
}

fn may_load_consensus_state(
    storage: &dyn Storage,
    height: u64,
) -> Result<Option<ConsensusState>, Error> {
    may_get_wasm_consensus_state(storage, height)?
        .map(|state| serde_json::from_slice(&state.data).map_err(Error::from))
        .transpose()
}

fn store_consensus_state(
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

fn store_client_state<Profile: Serialize>(
    storage: &mut dyn Storage,
    client: &ClientState<Profile>,
) -> Result<(), Error> {
    let mut envelope = get_wasm_client_state(storage)?;
    envelope.data = serde_json::to_vec(client)?;
    envelope.latest_height = Some(Height {
        revision_number: 0,
        revision_height: client.latest_height,
    });
    store_wasm_client_state(storage, &envelope)
}
