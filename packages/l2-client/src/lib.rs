//! Shared state and proof utilities for attestor-trusted L2 clients.

#![deny(clippy::nursery, clippy::pedantic, warnings, unused_crate_dependencies)]
#![allow(clippy::doc_markdown, clippy::missing_errors_doc)]

pub mod attestation;
pub mod canonical_header;
pub mod entrypoints;
pub mod error;
pub mod evm_proof;
pub mod msg;
pub mod packet;
pub mod runtime;
pub mod state;
pub mod store;
pub mod verification;

use cosmwasm_std::Api;
use serde::{de::DeserializeOwned, Serialize};

use crate::{
    error::Error,
    msg::AttestedL2Header,
    state::{Header, RuntimeProfile},
};

/// Chain-specific adapter around the common attested-header verifier.
pub trait L2LightClient {
    /// Immutable artifact profile.
    type Profile: Clone + RuntimeProfile + DeserializeOwned + Serialize;

    /// Locally validates and normalizes one attested update.
    fn verify(
        api: &dyn Api,
        profile: &Self::Profile,
        header: &AttestedL2Header,
    ) -> Result<Header, Error>;
}

/// Generates the three standard ICS-08 `CosmWasm` entrypoints.
#[macro_export]
macro_rules! l2_client_entrypoints {
    ($adapter:ty) => {
        /// Initializes explicitly trusted client and consensus state.
        #[cosmwasm_std::entry_point]
        pub fn instantiate(
            deps: cosmwasm_std::DepsMut,
            env: cosmwasm_std::Env,
            _info: cosmwasm_std::MessageInfo,
            msg: $crate::msg::InstantiateMsg,
        ) -> Result<cosmwasm_std::Response, l2_client::error::Error> {
            l2_client::entrypoints::instantiate::<$adapter>(deps, &env, msg)
        }

        /// Executes a host state transition or packet-proof check.
        #[cosmwasm_std::entry_point]
        pub fn sudo(
            deps: cosmwasm_std::DepsMut,
            env: cosmwasm_std::Env,
            msg: $crate::msg::SudoMsg,
        ) -> Result<cosmwasm_std::Response, l2_client::error::Error> {
            l2_client::entrypoints::sudo::<$adapter>(deps, &env, msg)
        }

        /// Executes a read-only light-client query.
        #[cosmwasm_std::entry_point]
        pub fn query(
            deps: cosmwasm_std::Deps,
            _env: cosmwasm_std::Env,
            msg: $crate::msg::QueryMsg,
        ) -> Result<cosmwasm_std::Binary, l2_client::error::Error> {
            l2_client::entrypoints::query::<$adapter>(deps, msg)
        }
    };
}

#[cfg(test)]
mod attestation_tests;
#[cfg(test)]
mod canonical_header_tests;
#[cfg(test)]
mod entrypoints_tests;
#[cfg(test)]
mod evm_proof_tests;
#[cfg(test)]
mod msg_tests;
#[cfg(test)]
mod packet_tests;
#[cfg(test)]
mod runtime_tests;
#[cfg(test)]
mod store_tests;
#[cfg(test)]
mod verification_tests;
