//! Base ICS-08 Wasm client entrypoint crate.

#![deny(
    clippy::nursery,
    clippy::pedantic,
    warnings,
    missing_docs,
    unused_crate_dependencies
)]

use base_verifier as _;
use cosmwasm_std as _;
use l2_client as _;

pub mod contract;
pub mod msg;

#[cfg(test)]
mod contract_tests;
