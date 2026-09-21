//! CosmWasm entrypoint implementations.

#![allow(clippy::needless_pass_by_value)]

use cosmwasm_std::{to_json_binary, Binary, Deps, DepsMut, Env, Response};

use crate::{
    bls::{AvalancheCustomQuery, QuerierBlsVerifier},
    error::Error,
    msg::{
        CheckForMisbehaviourResult, ClientMessage, IbcHeight, InstantiateMsg, MigrateMsg, QueryMsg,
        StatusResult, SudoMsg, TimestampAtHeightResult,
    },
    runtime,
    state::{ClientState, ConsensusState},
    verification::{self, AuthenticatedHeader},
};

/// Initializes the client with explicitly trusted bootstrap state.
pub fn instantiate(
    deps: DepsMut<AvalancheCustomQuery>,
    env: &Env,
    msg: InstantiateMsg,
) -> Result<Response, Error> {
    let client: ClientState = serde_json::from_slice(&msg.client_state)?;
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

/// Validates the current client and optionally rotates the pinned set.
pub fn migrate(deps: DepsMut<AvalancheCustomQuery>, msg: MigrateMsg) -> Result<Response, Error> {
    runtime::migrate(deps.storage, msg)?;
    Ok(Response::default())
}

/// Executes a host sudo operation.
///
/// No attributes, events, or messages may be attached to this response:
/// ibc-go's 08-wasm keeper rejects a light client that returns any of them by
/// panicking inside the VM call, so the whole tx fails rather than the update.
pub fn sudo(
    deps: DepsMut<AvalancheCustomQuery>,
    env: &Env,
    msg: SudoMsg,
) -> Result<Response, Error> {
    let data = match msg {
        SudoMsg::UpdateState { client_message } => {
            let header = decode_header(&client_message, deps.as_ref())?;
            runtime::update(deps.storage, &header, env.block.time.seconds())?
        }
        SudoMsg::UpdateStateOnMisbehaviour { client_message } => {
            let (first, second) = decode_misbehaviour(&client_message, deps.as_ref())?;
            if !runtime::apply_misbehaviour(deps.storage, &first, &second)? {
                return Err(Error::Invalid("headers do not prove misbehaviour"));
            }
            Binary::default()
        }
        SudoMsg::VerifyMembership {
            height,
            delay_time_period,
            delay_block_period,
            proof,
            merkle_path,
            value,
        } => {
            ensure_zero_delay(delay_time_period, delay_block_period)?;
            validate_height(&height)?;
            runtime::verify_membership(
                deps.storage,
                height.revision_height,
                &merkle_path.key_path,
                &value,
                &proof,
            )?;
            Binary::default()
        }
        SudoMsg::VerifyNonMembership {
            height,
            delay_time_period,
            delay_block_period,
            proof,
            merkle_path,
        } => {
            ensure_zero_delay(delay_time_period, delay_block_period)?;
            validate_height(&height)?;
            runtime::verify_non_membership(
                deps.storage,
                height.revision_height,
                &merkle_path.key_path,
                &proof,
            )?;
            Binary::default()
        }
        SudoMsg::VerifyUpgradeAndUpdateState { .. } => {
            return Err(Error::Invalid(
                "verify_upgrade_and_update_state is not supported",
            ));
        }
        SudoMsg::MigrateClientStore {} => {
            return Err(Error::Invalid("migrate_client_store is not supported"));
        }
    };
    Ok(Response::default().set_data(data))
}

const fn ensure_zero_delay(delay_time_period: u64, delay_block_period: u64) -> Result<(), Error> {
    if delay_time_period != 0 || delay_block_period != 0 {
        return Err(Error::Invalid("non-zero proof delays are not supported"));
    }
    Ok(())
}

/// Executes a read-only host query.
pub fn query(deps: Deps<AvalancheCustomQuery>, msg: QueryMsg) -> Result<Binary, Error> {
    match msg {
        QueryMsg::VerifyClientMessage { client_message } => {
            match decode_client_message(&client_message, deps)? {
                l2_client::msg::ClientMessage::Header(header) => {
                    header.header().consensus_state(0)?;
                }
                l2_client::msg::ClientMessage::Misbehaviour { header_1, header_2 } => {
                    if !runtime::is_actionable_misbehaviour(&header_1, &header_2)? {
                        return Err(Error::Invalid("headers do not prove misbehaviour"));
                    }
                }
            }
            Ok(Binary::default())
        }
        QueryMsg::CheckForMisbehaviour { client_message } => {
            let (first, second) = decode_misbehaviour(&client_message, deps)?;
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
            let client = runtime::client_state(deps.storage)?;
            client.validate()?;
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

// ClientMessage decoding always verifies through the warp quorum: an
// AuthenticatedHeader only exists past verification.

fn decode_header(
    data: &Binary,
    deps: Deps<AvalancheCustomQuery>,
) -> Result<AuthenticatedHeader, Error> {
    let ClientMessage::Header(header) = decode_raw(data)? else {
        return Err(Error::Invalid("header envelope required"));
    };
    let client = frozen_checked_client(deps)?;
    let bls = QuerierBlsVerifier {
        querier: deps.querier,
    };
    verification::verify_signed_header(&client, &header, &bls)
}

fn decode_misbehaviour(
    data: &Binary,
    deps: Deps<AvalancheCustomQuery>,
) -> Result<(AuthenticatedHeader, AuthenticatedHeader), Error> {
    let ClientMessage::Misbehaviour { header_1, header_2 } = decode_raw(data)? else {
        return Err(Error::Invalid("misbehaviour envelope required"));
    };
    let client = frozen_checked_client(deps)?;
    let bls = QuerierBlsVerifier {
        querier: deps.querier,
    };
    verification::verify_signed_misbehaviour(&client, &header_1, &header_2, &bls)
}

fn decode_client_message(
    data: &Binary,
    deps: Deps<AvalancheCustomQuery>,
) -> Result<l2_client::msg::ClientMessage<AuthenticatedHeader>, Error> {
    match decode_raw(data)? {
        ClientMessage::Header(header) => {
            let client = frozen_checked_client(deps)?;
            let bls = QuerierBlsVerifier {
                querier: deps.querier,
            };
            Ok(l2_client::msg::ClientMessage::Header(
                verification::verify_signed_header(&client, &header, &bls)?,
            ))
        }
        ClientMessage::Misbehaviour { header_1, header_2 } => {
            let client = frozen_checked_client(deps)?;
            let bls = QuerierBlsVerifier {
                querier: deps.querier,
            };
            let (header_1, header_2) =
                verification::verify_signed_misbehaviour(&client, &header_1, &header_2, &bls)?;
            Ok(l2_client::msg::ClientMessage::Misbehaviour { header_1, header_2 })
        }
    }
}

fn decode_raw(data: &Binary) -> Result<ClientMessage, Error> {
    serde_json::from_slice(data).map_err(Into::into)
}

fn frozen_checked_client(deps: Deps<AvalancheCustomQuery>) -> Result<ClientState, Error> {
    let client = runtime::client_state(deps.storage)?;
    client.validate()?;
    if let Some(height) = client.frozen_height {
        return Err(Error::Frozen(height));
    }
    Ok(client)
}

const fn validate_height(height: &IbcHeight) -> Result<(), Error> {
    if height.revision_number != 0 {
        return Err(Error::Invalid("revision number must be zero"));
    }
    if height.revision_height == 0 {
        return Err(Error::Invalid("height must be non-zero"));
    }
    Ok(())
}
