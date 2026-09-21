//! Avalanche C-Chain warp ICS-08 Wasm client entrypoint crate.

#![deny(
    clippy::nursery,
    clippy::pedantic,
    warnings,
    missing_docs,
    unused_crate_dependencies
)]

use avalanche_light_client as _;
use cosmwasm_std as _;

pub mod contract;
pub mod msg;
