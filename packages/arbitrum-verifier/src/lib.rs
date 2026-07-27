//! Arbitrum configuration and proof primitives for the ICS-08 Wasm client.

#![deny(
    clippy::nursery,
    clippy::pedantic,
    warnings,
    missing_docs,
    unused_crate_dependencies
)]

/// Immutable Arbitrum verification configuration.
pub mod config;
/// Protocol-selected update-header verification.
pub mod header;

#[cfg(test)]
mod profile_tests;

/// Thin Arbitrum adapter for the shared L2 contract entrypoints.
pub struct Adapter;

impl l2_client::L2LightClient for Adapter {
    type Header = header::Header;
    type Profile = config::Profile;

    fn l1_height(header: &Self::Header) -> u64 {
        header.beacon_slot()
    }

    fn verify(
        profile: &Self::Profile,
        authenticated_l1_root: alloy_primitives::B256,
        authenticated_l1_timestamp: u64,
        header: &Self::Header,
    ) -> Result<l2_client::state::Header, l2_client::error::Error> {
        header::verify(
            profile,
            authenticated_l1_root,
            authenticated_l1_timestamp,
            header,
        )
    }
}
