//! Avalanche C-Chain ICS-08 Wasm light client.
//!
//! Verifies the C-Chain by its own validator set: every update carries a
//! canonical coreth header, a warp `BitSetSignature` (signer bitset + 96-byte
//! aggregate BLS signature) over the block-hash warp message, and an account
//! proof of the ICS26Router against the header's state root. The client
//! recomputes the coreth block hash, rebuilds the exact unsigned warp message,
//! and requires a ≥ quorum stake-weighted aggregate from its **pinned
//! canonical validator set** before trusting anything the header claims.
//!
//! Avalanche is an L1 — acceptance is finality, so verified headers are final
//! and two verified conflicting headers at one height are misbehaviour that
//! freezes the client. The pinned set rotates only through chain governance
//! (`MigrateMsg::ReplaceValidatorSet`, monotonic in P-Chain height) in this
//! version; protocol-attested rotation is a planned upgrade.
//!
//! Reuses the shared EVM machinery from `l2-client` (coreth canonical header,
//! bounded MPT proofs, packet paths, 08-wasm envelopes) and the pure warp
//! verification from `avalanche-warp`. The runtime/transition layer is a
//! deliberate specialization, not the shared attestor kernel: the trust model
//! (stake quorum vs indexed operator signatures) differs in state, wire, and
//! authentication.

#![deny(clippy::nursery, clippy::pedantic, warnings, unused_crate_dependencies)]
#![allow(clippy::doc_markdown, clippy::missing_errors_doc)]

pub mod bls;
pub mod entrypoints;
pub mod error;
pub mod msg;
pub mod runtime;
pub mod state;
pub mod verification;

#[cfg(test)]
mod runtime_tests;
#[cfg(test)]
mod verification_tests;

/// Generates the three CosmWasm entrypoints for the artifact crate.
#[macro_export]
macro_rules! avalanche_client_entrypoints {
    () => {
        #[cosmwasm_std::entry_point]
        pub fn instantiate(
            deps: cosmwasm_std::DepsMut<$crate::bls::AvalancheCustomQuery>,
            env: cosmwasm_std::Env,
            _info: cosmwasm_std::MessageInfo,
            msg: $crate::msg::InstantiateMsg,
        ) -> Result<cosmwasm_std::Response, $crate::error::Error> {
            $crate::entrypoints::instantiate(deps, &env, msg)
        }

        #[cosmwasm_std::entry_point]
        pub fn migrate(
            deps: cosmwasm_std::DepsMut<$crate::bls::AvalancheCustomQuery>,
            _env: cosmwasm_std::Env,
            msg: $crate::msg::MigrateMsg,
        ) -> Result<cosmwasm_std::Response, $crate::error::Error> {
            $crate::entrypoints::migrate(deps, msg)
        }

        #[cosmwasm_std::entry_point]
        pub fn sudo(
            deps: cosmwasm_std::DepsMut<$crate::bls::AvalancheCustomQuery>,
            env: cosmwasm_std::Env,
            msg: $crate::msg::SudoMsg,
        ) -> Result<cosmwasm_std::Response, $crate::error::Error> {
            $crate::entrypoints::sudo(deps, &env, msg)
        }

        #[cosmwasm_std::entry_point]
        pub fn query(
            deps: cosmwasm_std::Deps<$crate::bls::AvalancheCustomQuery>,
            _env: cosmwasm_std::Env,
            msg: $crate::msg::QueryMsg,
        ) -> Result<cosmwasm_std::Binary, $crate::error::Error> {
            $crate::entrypoints::query(deps, msg)
        }
    };
}
