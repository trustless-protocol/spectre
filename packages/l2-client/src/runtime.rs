//! Storage-backed L2 client lifecycle.

use cosmwasm_std::Storage;
use ibc_proto::{
    google::protobuf::Any,
    ibc::{
        core::client::v1::Height,
        lightclients::wasm::v1::{
            ClientState as WasmClientState, ConsensusState as WasmConsensusState,
        },
    },
};
use prost::Message;
use serde::{de::DeserializeOwned, Serialize};

use crate::{
    error::Error,
    msg::EvmStorageProof,
    packet,
    state::{ClientState, ConsensusState, Header, RuntimeProfile},
    store::{
        get_wasm_client_state, get_wasm_consensus_state, may_get_wasm_consensus_state,
        store_wasm_client_state, store_wasm_consensus_state, HOST_CLIENT_STATE_KEY,
    },
};

/// Fully validated logical effects of one L2 update.
#[derive(Debug, PartialEq, Eq)]
pub(crate) struct StateTransition<Profile> {
    client: Option<ClientState<Profile>>,
    consensus: Option<(u64, ConsensusState)>,
}

impl<Profile> StateTransition<Profile> {
    const fn empty() -> Self {
        Self {
            client: None,
            consensus: None,
        }
    }

    #[cfg(test)]
    pub(crate) const fn client(&self) -> Option<&ClientState<Profile>> {
        self.client.as_ref()
    }

    #[cfg(test)]
    pub(crate) fn consensus(&self) -> Option<(u64, &ConsensusState)> {
        self.consensus
            .as_ref()
            .map(|(height, consensus)| (*height, consensus))
    }

    fn reported_height(&self) -> Option<u64> {
        self.consensus.as_ref().map(|(height, _)| *height)
    }
}

struct PreparedTransition {
    client: Option<Vec<u8>>,
    consensus: Option<(Vec<u8>, Vec<u8>)>,
}

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
    let client = client_state::<Profile>(storage)?;
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

    let transition = build_update_transition(
        client,
        header,
        existing.as_ref(),
        previous.as_ref(),
        next.as_ref(),
        accepted_at,
    )?;
    let reported_height = transition.reported_height();
    commit_transition(storage, &transition)?;
    Ok(reported_height)
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

pub(crate) fn build_update_transition<Profile>(
    mut client: ClientState<Profile>,
    header: &Header,
    existing: Option<&ConsensusState>,
    previous: Option<&ConsensusState>,
    next: Option<&ConsensusState>,
    accepted_at: u64,
) -> Result<StateTransition<Profile>, Error> {
    if let Some(height) = client.frozen_height {
        return Err(Error::Frozen(height));
    }

    let incoming = header.consensus_state(accepted_at)?;
    let height = header.height.revision_height;

    if existing.is_some_and(|stored| stored.conflicts_with(&incoming)) {
        if ATTESTATIONS_ARE_AUTHENTICATED {
            client.frozen_height = Some(height);
            return Ok(StateTransition {
                client: Some(client),
                consensus: None,
            });
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
        return Ok(StateTransition::empty());
    }

    client.latest_height = client.latest_height.max(height);
    Ok(StateTransition {
        client: Some(client),
        consensus: Some((height, incoming)),
    })
}

fn commit_transition<Profile: Serialize>(
    storage: &mut dyn Storage,
    transition: &StateTransition<Profile>,
) -> Result<(), Error> {
    let prepared = prepare_transition(storage, transition)?;
    if let Some((key, value)) = prepared.consensus {
        storage.set(&key, &value);
    }
    if let Some(value) = prepared.client {
        storage.set(HOST_CLIENT_STATE_KEY.as_bytes(), &value);
    }
    Ok(())
}

fn prepare_transition<Profile: Serialize>(
    storage: &dyn Storage,
    transition: &StateTransition<Profile>,
) -> Result<PreparedTransition, Error> {
    let client = transition
        .client
        .as_ref()
        .map(|client| -> Result<Vec<u8>, Error> {
            let mut envelope = get_wasm_client_state(storage)?;
            envelope.data = serde_json::to_vec(client)?;
            envelope.latest_height = Some(Height {
                revision_number: 0,
                revision_height: client.latest_height,
            });
            Ok(Any::from_msg(&envelope)?.encode_to_vec())
        })
        .transpose()?;

    let consensus = transition
        .consensus
        .as_ref()
        .map(|(height, consensus)| -> Result<(Vec<u8>, Vec<u8>), Error> {
            let envelope = Any::from_msg(&WasmConsensusState {
                data: serde_json::to_vec(consensus)?,
            })?;
            Ok((
                crate::store::consensus_db_key(*height).into_bytes(),
                envelope.encode_to_vec(),
            ))
        })
        .transpose()?;

    Ok(PreparedTransition { client, consensus })
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
