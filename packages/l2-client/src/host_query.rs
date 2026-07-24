//! Stargate queries for the Ethereum Wasm client bound to an L2 client.

use cosmwasm_std::{to_json_vec, Binary, ContractResult, Deps, Empty, QueryRequest, SystemResult};
use ethereum_light_client::consensus_state::ConsensusState as EthereumConsensusState;
use ibc_proto::ibc::{
    core::client::v1::{
        QueryClientStateRequest, QueryClientStateResponse, QueryClientStatusRequest,
        QueryClientStatusResponse, QueryConsensusStateRequest, QueryConsensusStateResponse,
    },
    lightclients::wasm::v1::{
        ClientState as WasmClientState, ConsensusState as WasmConsensusState,
    },
};
use prost::Message;

use crate::{error::Error, state::L1Client};

const CLIENT_STATE_PATH: &str = "/ibc.core.client.v1.Query/ClientState";
const CLIENT_STATUS_PATH: &str = "/ibc.core.client.v1.Query/ClientStatus";
const CONSENSUS_STATE_PATH: &str = "/ibc.core.client.v1.Query/ConsensusState";

/// Ethereum consensus state authenticated by the IBC host at one beacon slot.
#[derive(Clone, Debug, PartialEq, Eq)]
pub struct AuthenticatedL1State {
    /// The decoded Ethereum consensus state.
    pub consensus: EthereumConsensusState,
}

/// Validates the pinned Ethereum client identity and active status.
pub fn validate_l1_client(deps: Deps, l1: &L1Client) -> Result<(), Error> {
    let state: QueryClientStateResponse = stargate(
        deps,
        CLIENT_STATE_PATH,
        &QueryClientStateRequest {
            client_id: l1.client_id.clone(),
        },
    )?;
    validate_client_state(state, l1)?;

    let status: QueryClientStatusResponse = stargate(
        deps,
        CLIENT_STATUS_PATH,
        &QueryClientStatusRequest {
            client_id: l1.client_id.clone(),
        },
    )?;
    validate_client_status(status)
}

/// Queries and validates the Ethereum consensus state at `beacon_slot`.
pub fn consensus_state_at(
    deps: Deps,
    l1: &L1Client,
    beacon_slot: u64,
) -> Result<AuthenticatedL1State, Error> {
    validate_l1_client(deps, l1)?;
    let response: QueryConsensusStateResponse = stargate(
        deps,
        CONSENSUS_STATE_PATH,
        &QueryConsensusStateRequest {
            client_id: l1.client_id.clone(),
            revision_number: 0,
            revision_height: beacon_slot,
            latest_height: false,
        },
    )?;
    decode_consensus_state(response, beacon_slot)
}

#[allow(deprecated)]
fn stargate<Request, Response>(deps: Deps, path: &str, request: &Request) -> Result<Response, Error>
where
    Request: Message,
    Response: Message + Default,
{
    let request = QueryRequest::<Empty>::Stargate {
        path: path.into(),
        data: Binary::from(request.encode_to_vec()),
    };
    let request = to_json_vec(&request).map_err(|error| Error::HostQuery(error.to_string()))?;
    let raw = match deps.querier.raw_query(&request) {
        SystemResult::Ok(ContractResult::Ok(value)) => value,
        SystemResult::Ok(ContractResult::Err(error)) => return Err(Error::HostQuery(error)),
        SystemResult::Err(error) => return Err(Error::HostQuery(error.to_string())),
    };
    Response::decode(raw.as_slice()).map_err(Into::into)
}

pub(crate) fn validate_client_state(
    response: QueryClientStateResponse,
    l1: &L1Client,
) -> Result<(), Error> {
    let any = response
        .client_state
        .ok_or_else(|| Error::HostQuery("missing client state".into()))?;
    if any.type_url != "/ibc.lightclients.wasm.v1.ClientState" {
        return Err(Error::UnexpectedHostTypeUrl(any.type_url));
    }
    let state = WasmClientState::decode(any.value.as_slice())?;
    if state.checksum != l1.wasm_checksum {
        return Err(Error::L1ChecksumMismatch);
    }
    Ok(())
}

pub(crate) fn validate_client_status(response: QueryClientStatusResponse) -> Result<(), Error> {
    if response.status != "Active" {
        return Err(Error::L1ClientInactive(response.status));
    }
    Ok(())
}

pub(crate) fn decode_consensus_state(
    response: QueryConsensusStateResponse,
    beacon_slot: u64,
) -> Result<AuthenticatedL1State, Error> {
    let any = response
        .consensus_state
        .ok_or_else(|| Error::HostQuery("missing consensus state".into()))?;
    if any.type_url != "/ibc.lightclients.wasm.v1.ConsensusState" {
        return Err(Error::UnexpectedHostTypeUrl(any.type_url));
    }
    let state = WasmConsensusState::decode(any.value.as_slice())?;
    let consensus: EthereumConsensusState = serde_json::from_slice(&state.data)?;
    if consensus.slot != beacon_slot {
        return Err(Error::L1ConsensusSlotMismatch {
            expected: beacon_slot,
            actual: consensus.slot,
        });
    }
    Ok(AuthenticatedL1State { consensus })
}
