//! Shared state and verification utilities for ICS-08 L2 clients.

#![deny(
    clippy::nursery,
    clippy::pedantic,
    warnings,
    missing_docs,
    unused_crate_dependencies
)]

pub mod canonical_header;
pub mod error;
pub mod evm_proof;
/// Deterministic fixture loading.
pub mod fixtures;
/// Stargate adapter for authenticated state of the pinned Ethereum Wasm client.
pub mod host_query;
pub mod misbehaviour;
pub mod msg;
pub mod packet;
pub mod query;
/// Storage-backed lifecycle operations used by wrapper artifacts.
pub mod runtime;
pub mod state;
pub mod store;
pub mod sudo;
/// Shared helpers used by chain-specific rollup verifiers.
pub mod verification;

use alloy_primitives::B256;
use serde::{de::DeserializeOwned, Serialize};

use crate::{error::Error, state::RuntimeProfile};

/// Shared normalization interface implemented by rollup-specific verifiers.
pub trait RollupVerifier {
    /// Immutable rollup profile.
    type Profile: RuntimeProfile;
    /// Chain-specific update header.
    type Header;

    /// Authenticates evidence and returns the normalized L2 state.
    fn verify_header(
        client_state: &state::ClientState<Self::Profile>,
        header: &Self::Header,
        authenticated_l1_state: &state::ConsensusState,
    ) -> Result<state::VerifiedL2State, Error>;
}

/// Chain-specific adapter used by the shared CosmWasm entrypoints.
pub trait L2LightClient {
    /// Immutable runtime profile stored inside client state.
    type Profile: Clone + RuntimeProfile + DeserializeOwned + Serialize;
    /// Chain-specific optimistic update header.
    type Header: DeserializeOwned;

    /// Returns typed finality evidence embedded in the update, if present.
    fn finality_evidence(_header: &Self::Header) -> Option<&crate::msg::FinalityEvidence> {
        None
    }

    /// Returns the exact Ethereum beacon height referenced by a header.
    fn l1_height(header: &Self::Header) -> u64;

    /// Authenticates a rollup proposal and its canonical L2 router state.
    fn verify(
        profile: &Self::Profile,
        authenticated_l1_root: B256,
        authenticated_l1_timestamp: u64,
        header: &Self::Header,
    ) -> Result<state::Header, Error>;

    /// Authenticates a header and applies its typed finality evidence exactly once.
    fn verify_with_finality(
        profile: &Self::Profile,
        policy: &state::FinalityPolicy,
        authenticated_l1_root: B256,
        authenticated_l1_timestamp: u64,
        header: &Self::Header,
    ) -> Result<state::Header, Error> {
        let mut verified = Self::verify(
            profile,
            authenticated_l1_root,
            authenticated_l1_timestamp,
            header,
        )?;
        if let Some(evidence) = Self::finality_evidence(header) {
            crate::verification::apply_finality_evidence(policy, evidence, &mut verified)?;
        }
        Ok(verified)
    }
}

/// Generates the common optimistic L2 CosmWasm entrypoints for one thin adapter crate.
#[macro_export]
macro_rules! l2_client_entrypoints {
    ($adapter:ty) => {
        /// Stores directly supplied client and consensus state without an L1 query.
        #[cosmwasm_std::entry_point]
        pub fn instantiate(
            deps: cosmwasm_std::DepsMut,
            env: cosmwasm_std::Env,
            _info: cosmwasm_std::MessageInfo,
            msg: crate::msg::InstantiateMsg,
        ) -> Result<cosmwasm_std::Response, l2_client::error::Error> {
            let client: l2_client::state::ClientState<
                <$adapter as l2_client::L2LightClient>::Profile,
            > = serde_json::from_slice(&msg.client_state)?;
            let consensus: l2_client::state::ConsensusState =
                serde_json::from_slice(&msg.consensus_state)?;
            l2_client::runtime::instantiate(
                deps.storage,
                &client,
                &consensus,
                msg.checksum.to_vec(),
                env.block.time.seconds(),
            )?;
            Ok(cosmwasm_std::Response::default())
        }

        /// Routes state transitions and packet proof checks.
        #[cosmwasm_std::entry_point]
        pub fn sudo(
            deps: cosmwasm_std::DepsMut,
            env: cosmwasm_std::Env,
            msg: crate::msg::SudoMsg,
        ) -> Result<cosmwasm_std::Response, l2_client::error::Error> {
            use l2_client::state::RuntimeProfile as _;
            // No attributes, events, or messages may be attached to this response:
            // ibc-go's 08-wasm keeper rejects a light client that returns any of them
            // ("returning attributes from a contract is not allowed"), and it does so
            // by panicking inside the VM call, so the whole tx fails rather than the
            // update. A light client is a pure state transition; observability has to
            // come from the caller or from querying the client state afterwards.
            let data = match msg {
                crate::msg::SudoMsg::UpdateState { client_message } => {
                    let header = match serde_json::from_slice::<
                        l2_client::msg::ClientMessage<
                            <$adapter as l2_client::L2LightClient>::Header,
                        >,
                    >(&client_message)
                    {
                        Ok(l2_client::msg::ClientMessage::Header(header)) => header,
                        Ok(l2_client::msg::ClientMessage::Misbehaviour { .. }) => {
                            return Err(l2_client::error::Error::InvalidHeader(
                                "header envelope required",
                            )
                            .into());
                        }
                        Err(_) => serde_json::from_slice(&client_message)?,
                    };
                    let client = l2_client::runtime::client_state::<
                        <$adapter as l2_client::L2LightClient>::Profile,
                    >(deps.storage)?;
                    if let Some(height) = client.frozen_height {
                        return Err(l2_client::error::Error::Frozen(height).into());
                    }
                    let common = client.profile.common();
                    let l1 = l2_client::host_query::consensus_state_at(
                        deps.as_ref(),
                        &common.ethereum_client,
                        <$adapter as l2_client::L2LightClient>::l1_height(&header),
                    )?;
                    let verified = <$adapter as l2_client::L2LightClient>::verify_with_finality(
                        &client.profile,
                        &client.finality_policy,
                        l1.consensus.state_root,
                        l1.consensus.timestamp,
                        &header,
                    )?;
                    let updated_height = l2_client::runtime::update::<
                        <$adapter as l2_client::L2LightClient>::Profile,
                    >(deps.storage, &verified)?;
                    cosmwasm_std::to_json_binary(&l2_client::msg::UpdateStateResult {
                        heights: updated_height
                            .into_iter()
                            .map(|revision_height| l2_client::msg::IbcHeight {
                                revision_number: 0,
                                revision_height,
                            })
                            .collect(),
                    })?
                }
                crate::msg::SudoMsg::UpdateStateOnMisbehaviour { client_message } => {
                    let evidence: l2_client::msg::ClientMessage<
                        <$adapter as l2_client::L2LightClient>::Header,
                    > = serde_json::from_slice(&client_message)?;
                    let l2_client::msg::ClientMessage::Misbehaviour { header_1, header_2 } =
                        evidence
                    else {
                        return Err(l2_client::error::Error::InvalidHeader(
                            "misbehaviour envelope required",
                        )
                        .into());
                    };
                    let client = l2_client::runtime::client_state::<
                        <$adapter as l2_client::L2LightClient>::Profile,
                    >(deps.storage)?;
                    if let Some(height) = client.frozen_height {
                        return Err(l2_client::error::Error::Frozen(height).into());
                    }
                    let common = client.profile.common();
                    let first_l1 = l2_client::host_query::consensus_state_at(
                        deps.as_ref(),
                        &common.ethereum_client,
                        <$adapter as l2_client::L2LightClient>::l1_height(&header_1),
                    )?;
                    let second_l1 = l2_client::host_query::consensus_state_at(
                        deps.as_ref(),
                        &common.ethereum_client,
                        <$adapter as l2_client::L2LightClient>::l1_height(&header_2),
                    )?;
                    let header_1 = <$adapter as l2_client::L2LightClient>::verify_with_finality(
                        &client.profile,
                        &client.finality_policy,
                        first_l1.consensus.state_root,
                        first_l1.consensus.timestamp,
                        &header_1,
                    )?;
                    let header_2 = <$adapter as l2_client::L2LightClient>::verify_with_finality(
                        &client.profile,
                        &client.finality_policy,
                        second_l1.consensus.state_root,
                        second_l1.consensus.timestamp,
                        &header_2,
                    )?;
                    let found = l2_client::runtime::apply_misbehaviour::<
                        <$adapter as l2_client::L2LightClient>::Profile,
                    >(deps.storage, &header_1, &header_2)?;
                    if !found {
                        return Err(l2_client::error::Error::InvalidHeader(
                            "headers do not prove misbehaviour",
                        )
                        .into());
                    }
                    cosmwasm_std::Binary::default()
                }
                crate::msg::SudoMsg::VerifyMembership {
                    height,
                    delay_time_period: _,
                    delay_block_period: _,
                    proof,
                    merkle_path,
                    value,
                } => {
                    if height.revision_number != 0 {
                        return Err(l2_client::error::Error::InvalidRevision(
                            height.revision_number,
                        )
                        .into());
                    }
                    let [path] = merkle_path.key_path.as_slice() else {
                        return Err(l2_client::error::Error::Proof(
                            "exactly one IBC commitment path is required".into(),
                        )
                        .into());
                    };
                    let proof = serde_json::from_slice(&proof)?;
                    l2_client::runtime::verify_membership_at::<
                        <$adapter as l2_client::L2LightClient>::Profile,
                    >(
                        deps.storage,
                        height.revision_height,
                        path,
                        &value,
                        &proof,
                        env.block.time.seconds(),
                    )?;
                    cosmwasm_std::Binary::default()
                }
                crate::msg::SudoMsg::VerifyNonMembership {
                    height,
                    delay_time_period: _,
                    delay_block_period: _,
                    proof,
                    merkle_path,
                } => {
                    if height.revision_number != 0 {
                        return Err(l2_client::error::Error::InvalidRevision(
                            height.revision_number,
                        )
                        .into());
                    }
                    let [path] = merkle_path.key_path.as_slice() else {
                        return Err(l2_client::error::Error::Proof(
                            "exactly one IBC commitment path is required".into(),
                        )
                        .into());
                    };
                    let proof = serde_json::from_slice(&proof)?;
                    l2_client::runtime::verify_non_membership_at::<
                        <$adapter as l2_client::L2LightClient>::Profile,
                    >(
                        deps.storage,
                        height.revision_height,
                        path,
                        &proof,
                        env.block.time.seconds(),
                    )?;
                    cosmwasm_std::Binary::default()
                }
            };
            // ibc-go 08-wasm rejects sudo responses with attributes/events/messages.
            // Return data only; relayer observability comes from host tx logs.
            Ok(cosmwasm_std::Response::default().set_data(data))
        }

        /// Routes strict client queries.
        #[cosmwasm_std::entry_point]
        pub fn query(
            deps: cosmwasm_std::Deps,
            _env: cosmwasm_std::Env,
            msg: crate::msg::QueryMsg,
        ) -> Result<cosmwasm_std::Binary, l2_client::error::Error> {
            use l2_client::state::RuntimeProfile as _;
            match msg {
                crate::msg::QueryMsg::VerifyClientMessage { client_message } => {
                    let message = match serde_json::from_slice::<
                        l2_client::msg::ClientMessage<
                            <$adapter as l2_client::L2LightClient>::Header,
                        >,
                    >(&client_message)
                    {
                        Ok(message) => message,
                        Err(_) => l2_client::msg::ClientMessage::Header(serde_json::from_slice(
                            &client_message,
                        )?),
                    };
                    let client = l2_client::runtime::client_state::<
                        <$adapter as l2_client::L2LightClient>::Profile,
                    >(deps.storage)?;
                    let common = client.profile.common();
                    match message {
                        l2_client::msg::ClientMessage::Header(header) => {
                            let l1 = l2_client::host_query::consensus_state_at(
                                deps,
                                &common.ethereum_client,
                                <$adapter as l2_client::L2LightClient>::l1_height(&header),
                            )?;
                            let verified =
                                <$adapter as l2_client::L2LightClient>::verify_with_finality(
                                    &client.profile,
                                    &client.finality_policy,
                                    l1.consensus.state_root,
                                    l1.consensus.timestamp,
                                    &header,
                                )?;
                            l2_client::query::verify_client_message(&verified)?;
                        }
                        l2_client::msg::ClientMessage::Misbehaviour { header_1, header_2 } => {
                            let first_l1 = l2_client::host_query::consensus_state_at(
                                deps,
                                &common.ethereum_client,
                                <$adapter as l2_client::L2LightClient>::l1_height(&header_1),
                            )?;
                            let second_l1 = l2_client::host_query::consensus_state_at(
                                deps,
                                &common.ethereum_client,
                                <$adapter as l2_client::L2LightClient>::l1_height(&header_2),
                            )?;
                            let header_1 =
                                <$adapter as l2_client::L2LightClient>::verify_with_finality(
                                    &client.profile,
                                    &client.finality_policy,
                                    first_l1.consensus.state_root,
                                    first_l1.consensus.timestamp,
                                    &header_1,
                                )?;
                            let header_2 =
                                <$adapter as l2_client::L2LightClient>::verify_with_finality(
                                    &client.profile,
                                    &client.finality_policy,
                                    second_l1.consensus.state_root,
                                    second_l1.consensus.timestamp,
                                    &header_2,
                                )?;
                            if !l2_client::query::check_for_misbehaviour(
                                &client.finality_policy,
                                &header_1,
                                &header_2,
                            )? {
                                return Err(l2_client::error::Error::InvalidHeader(
                                    "headers do not prove misbehaviour",
                                )
                                .into());
                            }
                        }
                    }
                    Ok(cosmwasm_std::Binary::default())
                }
                crate::msg::QueryMsg::CheckForMisbehaviour { client_message } => {
                    let evidence: l2_client::msg::ClientMessage<
                        <$adapter as l2_client::L2LightClient>::Header,
                    > = serde_json::from_slice(&client_message)?;
                    let l2_client::msg::ClientMessage::Misbehaviour { header_1, header_2 } =
                        evidence
                    else {
                        return Err(l2_client::error::Error::InvalidHeader(
                            "misbehaviour envelope required",
                        )
                        .into());
                    };
                    let client = l2_client::runtime::client_state::<
                        <$adapter as l2_client::L2LightClient>::Profile,
                    >(deps.storage)?;
                    let common = client.profile.common();
                    let first_l1 = l2_client::host_query::consensus_state_at(
                        deps,
                        &common.ethereum_client,
                        <$adapter as l2_client::L2LightClient>::l1_height(&header_1),
                    )?;
                    let second_l1 = l2_client::host_query::consensus_state_at(
                        deps,
                        &common.ethereum_client,
                        <$adapter as l2_client::L2LightClient>::l1_height(&header_2),
                    )?;
                    let header_1 = <$adapter as l2_client::L2LightClient>::verify_with_finality(
                        &client.profile,
                        &client.finality_policy,
                        first_l1.consensus.state_root,
                        first_l1.consensus.timestamp,
                        &header_1,
                    )?;
                    let header_2 = <$adapter as l2_client::L2LightClient>::verify_with_finality(
                        &client.profile,
                        &client.finality_policy,
                        second_l1.consensus.state_root,
                        second_l1.consensus.timestamp,
                        &header_2,
                    )?;
                    Ok(cosmwasm_std::to_json_binary(
                        &l2_client::msg::CheckForMisbehaviourResult {
                            found_misbehaviour: l2_client::query::check_for_misbehaviour(
                                &client.finality_policy,
                                &header_1,
                                &header_2,
                            )?,
                        },
                    )?)
                }
                crate::msg::QueryMsg::TimestampAtHeight { height } => {
                    if height.revision_number != 0 {
                        return Err(l2_client::error::Error::InvalidRevision(
                            height.revision_number,
                        )
                        .into());
                    }
                    Ok(cosmwasm_std::to_json_binary(
                        &l2_client::msg::TimestampAtHeightResult {
                            timestamp: l2_client::query::timestamp_at_height(
                                &l2_client::runtime::consensus_state(
                                    deps.storage,
                                    height.revision_height,
                                )?,
                            ),
                        },
                    )?)
                }
                crate::msg::QueryMsg::Status {} => {
                    let client = l2_client::runtime::client_state::<
                        <$adapter as l2_client::L2LightClient>::Profile,
                    >(deps.storage)?;
                    if client.frozen_height.is_some() {
                        return Ok(cosmwasm_std::to_json_binary(
                            &l2_client::msg::StatusResult {
                                status: "Frozen".into(),
                            },
                        )?);
                    }
                    l2_client::host_query::validate_l1_client(
                        deps,
                        &client.profile.common().ethereum_client,
                    )?;
                    Ok(cosmwasm_std::to_json_binary(
                        &l2_client::msg::StatusResult {
                            status: "Active".into(),
                        },
                    )?)
                }
            }
        }
    };
}

#[cfg(test)]
mod adversarial_tests;
#[cfg(test)]
mod canonical_header_tests;
#[cfg(test)]
#[cfg(test)]
mod evm_proof_tests;
#[cfg(test)]
mod finality_tests;
#[cfg(test)]
mod fixtures_tests;
#[cfg(test)]
mod host_query_tests;
#[cfg(test)]
mod msg_tests;
#[cfg(test)]
mod packet_tests;
#[cfg(test)]
mod query_tests;
#[cfg(test)]
mod runtime_tests;
#[cfg(test)]
mod store_tests;
#[cfg(test)]
mod sudo_tests;
