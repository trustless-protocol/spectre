//! Generic entrypoint implementation used by the chain-specific artifacts.

#![allow(clippy::needless_pass_by_value)]

use cosmwasm_std::{to_json_binary, Binary, Deps, DepsMut, Env, Response};

use crate::{
    error::Error,
    msg::{
        AttestedL2Header, CheckForMisbehaviourResult, ClientMessage, IbcHeight, InstantiateMsg,
        QueryMsg, StatusResult, SudoMsg, TimestampAtHeightResult, UpdateStateResult,
    },
    runtime,
    state::{ClientState, ConsensusState},
    L2LightClient,
};

/// Initializes a chain-specific artifact using explicitly trusted bootstrap state.
pub fn instantiate<Adapter: L2LightClient>(
    deps: DepsMut,
    env: &Env,
    msg: InstantiateMsg,
) -> Result<Response, Error> {
    let client: ClientState<Adapter::Profile> = serde_json::from_slice(&msg.client_state)?;
    let consensus: ConsensusState = serde_json::from_slice(&msg.consensus_state)?;
    runtime::instantiate(
        deps.storage,
        &client,
        &consensus,
        msg.checksum.to_vec(),
        env.block.time.seconds(),
    )?;
    Ok(Response::default())
}

/// Executes a host sudo operation without emitting attributes, events, or messages.
///
/// No attributes, events, or messages may be attached to this response: ibc-go's 08-wasm keeper
/// rejects a light client that returns any of them ("returning attributes from a contract is not
/// allowed"), and it does so by panicking inside the VM call, so the whole tx fails rather than
/// the update. A light client is a pure state transition; observability has to come from the
/// caller or from querying the client state afterwards.
pub fn sudo<Adapter: L2LightClient>(
    deps: DepsMut,
    env: &Env,
    msg: SudoMsg,
) -> Result<Response, Error> {
    let data = match msg {
        SudoMsg::UpdateState { client_message } => {
            let header = decode_header::<Adapter>(&client_message, deps.as_ref())?;
            let updated = runtime::update::<Adapter::Profile>(
                deps.storage,
                &header,
                env.block.time.seconds(),
            )?;
            to_json_binary(&UpdateStateResult {
                heights: updated
                    .into_iter()
                    .map(|revision_height| IbcHeight {
                        revision_number: 0,
                        revision_height,
                    })
                    .collect(),
            })?
        }
        SudoMsg::UpdateStateOnMisbehaviour { client_message } => {
            let (first, second) = decode_misbehaviour::<Adapter>(&client_message, deps.as_ref())?;
            if !runtime::apply_misbehaviour::<Adapter::Profile>(deps.storage, &first, &second)? {
                return Err(Error::InvalidHeader("headers do not prove misbehaviour"));
            }
            Binary::default()
        }
        SudoMsg::VerifyMembership {
            height,
            proof,
            merkle_path,
            value,
            ..
        } => {
            let path = single_path(&merkle_path.key_path)?;
            validate_height(&height)?;
            runtime::verify_membership::<Adapter::Profile>(
                deps.storage,
                height.revision_height,
                path,
                &value,
                &serde_json::from_slice(&proof)?,
            )?;
            Binary::default()
        }
        SudoMsg::VerifyNonMembership {
            height,
            proof,
            merkle_path,
            ..
        } => {
            let path = single_path(&merkle_path.key_path)?;
            validate_height(&height)?;
            runtime::verify_non_membership::<Adapter::Profile>(
                deps.storage,
                height.revision_height,
                path,
                &serde_json::from_slice(&proof)?,
            )?;
            Binary::default()
        }
    };
    Ok(Response::default().set_data(data))
}

/// Executes a read-only host query.
pub fn query<Adapter: L2LightClient>(deps: Deps, msg: QueryMsg) -> Result<Binary, Error> {
    match msg {
        QueryMsg::VerifyClientMessage { client_message } => {
            match decode_client_message::<Adapter>(&client_message, deps)? {
                ClientMessage::Header(header) => {
                    header.consensus_state(0)?;
                }
                ClientMessage::Misbehaviour { header_1, header_2 } => {
                    if !runtime::is_actionable_misbehaviour(&header_1, &header_2)? {
                        return Err(Error::InvalidHeader("headers do not prove misbehaviour"));
                    }
                }
            }
            Ok(Binary::default())
        }
        QueryMsg::CheckForMisbehaviour { client_message } => {
            let (first, second) = decode_misbehaviour::<Adapter>(&client_message, deps)?;
            to_json_binary(&CheckForMisbehaviourResult {
                found_misbehaviour: runtime::is_actionable_misbehaviour(&first, &second)?,
            })
            .map_err(Into::into)
        }
        QueryMsg::TimestampAtHeight { height } => {
            validate_height(&height)?;
            to_json_binary(&TimestampAtHeightResult {
                timestamp: runtime::consensus_state(deps.storage, height.revision_height)?
                    .timestamp_nanos,
            })
            .map_err(Into::into)
        }
        QueryMsg::Status {} => {
            let client = runtime::client_state::<Adapter::Profile>(deps.storage)?;
            to_json_binary(&StatusResult {
                status: if client.frozen_height.is_some() {
                    "Frozen"
                } else {
                    "Active"
                }
                .into(),
            })
            .map_err(Into::into)
        }
    }
}

fn decode_header<Adapter: L2LightClient>(
    data: &Binary,
    deps: Deps,
) -> Result<crate::state::Header, Error> {
    let ClientMessage::Header(header) = decode_client_message::<Adapter>(data, deps)? else {
        return Err(Error::InvalidHeader("header envelope required"));
    };
    Ok(header)
}

fn decode_misbehaviour<Adapter: L2LightClient>(
    data: &Binary,
    deps: Deps,
) -> Result<(crate::state::Header, crate::state::Header), Error> {
    let ClientMessage::Misbehaviour { header_1, header_2 } =
        decode_client_message::<Adapter>(data, deps)?
    else {
        return Err(Error::InvalidHeader("misbehaviour envelope required"));
    };
    Ok((header_1, header_2))
}

fn decode_client_message<Adapter: L2LightClient>(
    data: &Binary,
    deps: Deps,
) -> Result<ClientMessage<crate::state::Header>, Error> {
    let message: ClientMessage<AttestedL2Header> = serde_json::from_slice(data)?;
    let client = runtime::client_state::<Adapter::Profile>(deps.storage)?;
    match message {
        ClientMessage::Header(header) => Ok(ClientMessage::Header(Adapter::verify(
            &client.profile,
            &header,
        )?)),
        ClientMessage::Misbehaviour { header_1, header_2 } => Ok(ClientMessage::Misbehaviour {
            header_1: Adapter::verify(&client.profile, &header_1)?,
            header_2: Adapter::verify(&client.profile, &header_2)?,
        }),
    }
}

const fn validate_height(height: &IbcHeight) -> Result<(), Error> {
    if height.revision_number != 0 {
        return Err(Error::InvalidRevision(height.revision_number));
    }
    if height.revision_height == 0 {
        return Err(Error::ZeroHeight);
    }
    Ok(())
}

fn single_path(paths: &[Binary]) -> Result<&Binary, Error> {
    let [path] = paths else {
        return Err(Error::Proof(
            "exactly one IBC commitment path is required".into(),
        ));
    };
    Ok(path)
}
