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
    accepted_at: u64,
) -> Result<(), Error>
where
    Profile: Clone + Serialize,
{
    client.finality_policy.validate()?;
    let mut client = client.clone();
    let height = client.latest_height;
    let mut consensus = consensus.clone();
    if consensus.l2_height == 0 {
        consensus.l2_height = height;
    } else if consensus.l2_height != height {
        return Err(Error::InvalidHeader(
            "bootstrap consensus height differs from client height",
        ));
    }
    consensus.first_accepted_at = accepted_at;
    consensus.finality_reached_at = accepted_at;
    if consensus.finality_level == crate::state::FinalityLevel::Finalized {
        client.last_finalized_update_at = Some(
            client
                .last_finalized_update_at
                .unwrap_or_default()
                .max(accepted_at),
        );
    }
    store_wasm_client_state(
        storage,
        &WasmClientState {
            checksum,
            data: serde_json::to_vec(&client)?,
            latest_height: Some(Height {
                revision_number: 0,
                revision_height: height,
            }),
        },
    )?;
    store_consensus(storage, height, &consensus)
}

/// Loads the compact active state from host-owned storage.
pub fn client_state<Config: DeserializeOwned>(
    storage: &dyn Storage,
) -> Result<ClientState<Config>, Error> {
    serde_json::from_slice(&get_wasm_client_state(storage)?.data).map_err(Into::into)
}

/// Loads the consensus state at exactly `height`.
pub fn consensus_state(storage: &dyn Storage, height: u64) -> Result<ConsensusState, Error> {
    let mut state: ConsensusState =
        serde_json::from_slice(&get_wasm_consensus_state(storage, height)?.data)?;
    // Legacy payloads had no authenticated L2 height or finality metadata. Their defaults are
    // deliberately conservative; the envelope key is the only safe height to retain.
    if state.l2_height == 0 {
        state.l2_height = height;
    }
    Ok(state)
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
    let previous = if height > 1 {
        may_get_wasm_consensus_state(storage, height - 1)?
            .map(|state| serde_json::from_slice(&state.data))
            .transpose()?
    } else {
        None
    };
    let old_latest = client.latest_height;
    let result = match sudo::update_state_with_parent(
        &mut client,
        header,
        existing.as_ref(),
        previous.as_ref(),
    ) {
        Ok(result) => result,
        // A trusted conflict is a successful state transition: the client must be frozen on
        // chain. Returning this error through the CosmWasm entrypoint would revert the write.
        Err(Error::TrustedConflict { .. }) => {
            store_client(storage, &client)?;
            return Ok(None);
        }
        Err(error) => return Err(error),
    };
    let Some(consensus) = result else {
        return Ok(None);
    };
    store_consensus(storage, height, &consensus)?;
    if consensus.finality_level == crate::state::FinalityLevel::Finalized {
        client.last_finalized_update_at = Some(
            client
                .last_finalized_update_at
                .unwrap_or_default()
                .max(consensus.finality_reached_at),
        );
    }
    if client.latest_height != old_latest || client.last_finalized_update_at.is_some() {
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
    verify_membership_at::<Profile>(storage, height, path, value, proof, 0)
}

/// Verifies membership while applying finality, maturity, and freshness policy at `now`.
pub fn verify_membership_at<Profile: DeserializeOwned + RuntimeProfile>(
    storage: &dyn Storage,
    height: u64,
    path: &[u8],
    value: &[u8],
    proof: &EvmStorageProof,
    now: u64,
) -> Result<(), Error> {
    let client = client_state::<Profile>(storage)?;
    if let Some(frozen) = client.frozen_height {
        return Err(Error::Frozen(frozen));
    }
    let consensus = consensus_state(storage, height)?;
    ensure_consensus_state_is_usable(&client, &consensus, now)?;
    packet::verify_membership(
        consensus.ibc_storage_root,
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
    verify_non_membership_at::<Profile>(storage, height, path, proof, 0)
}

/// Verifies non-membership while applying finality, maturity, and freshness policy at `now`.
pub fn verify_non_membership_at<Profile: DeserializeOwned + RuntimeProfile>(
    storage: &dyn Storage,
    height: u64,
    path: &[u8],
    proof: &EvmStorageProof,
    now: u64,
) -> Result<(), Error> {
    let client = client_state::<Profile>(storage)?;
    if let Some(frozen) = client.frozen_height {
        return Err(Error::Frozen(frozen));
    }
    let consensus = consensus_state(storage, height)?;
    ensure_consensus_state_is_usable(&client, &consensus, now)?;
    packet::verify_non_membership(
        consensus.ibc_storage_root,
        client.profile.common().commitment_slot,
        path,
        proof,
    )
}

/// Enforces membership finality, proposal resolution, maturity, and freshness policy.
pub fn ensure_consensus_state_is_usable<Profile>(
    client: &ClientState<Profile>,
    consensus: &ConsensusState,
    now: u64,
) -> Result<(), Error> {
    if consensus.finality_level < client.finality_policy.minimum_membership_level {
        return Err(Error::InsufficientFinality {
            required: client.finality_policy.minimum_membership_level,
            actual: consensus.finality_level,
        });
    }
    match consensus.proposal_status {
        crate::state::ProposalStatus::ResolvedInvalid => {
            return Err(Error::ProposalResolvedInvalid);
        }
        crate::state::ProposalStatus::Pending
            if client.finality_policy.require_resolved_proposal =>
        {
            return Err(Error::ProposalPending);
        }
        crate::state::ProposalStatus::Pending | crate::state::ProposalStatus::ResolvedValid => {}
    }
    let usable_at = consensus
        .finality_reached_at
        .checked_add(client.finality_policy.maturity_delay_seconds)
        .ok_or(Error::TimestampOverflow)?;
    if now < usable_at {
        return Err(Error::FinalityMaturityPending { usable_at });
    }
    if let Some(max_age) = client.freshness_policy.max_time_without_finalized_update {
        let last = client.last_finalized_update_at.ok_or(Error::ClientStale)?;
        if now > last.saturating_add(max_age) {
            return Err(Error::ClientStale);
        }
    }
    if let Some(max_lag) = client.freshness_policy.max_l2_height_lag {
        let lag = client.latest_height.saturating_sub(consensus.l2_height);
        if lag > max_lag {
            return Err(Error::ClientStale);
        }
    }
    Ok(())
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
