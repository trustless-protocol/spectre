//! Data-only Avalanche C-Chain profile for the shared attested-header client.
//!
//! Avalanche is an L1, not a rollup — it merely reuses the attested-header
//! client machinery the rollup profiles run on, and fits it exactly:
//! Avalanche acceptance is finality (no reorgs, no challenge window), so the
//! client verifies an attestor-signed canonical coreth header plus a router
//! account proof. The coreth header shape (ext_data_hash + coreth optional
//! tail) is selected by pinning `l2_header_fork: "coreth"` in the profile —
//! the field name is the shared wire schema's, not a statement about the
//! chain.

#![deny(clippy::nursery, clippy::pedantic, warnings, unused_crate_dependencies)]

use l2_client::state::{CommonProfile, RuntimeProfile};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

/// Immutable Avalanche C-Chain artifact profile.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct AvalancheProfile {
    /// Chain-independent verification parameters (shared wire schema).
    pub common: CommonProfile,
}

impl RuntimeProfile for AvalancheProfile {
    fn common(&self) -> &CommonProfile {
        &self.common
    }

    fn expected_profile_version() -> &'static str {
        "avalanche_attestor_v1"
    }
}

#[cfg(test)]
mod tests {
    use super::AvalancheProfile;
    use l2_client::canonical_header::ExecutionHeaderFork;

    const FUJI_PROFILE_BYTES: &[u8] = b"{\"common\":{\"l2_chain_id\":43113,\"l2_router\":\"0x0000000000000000000000000000000000000000\",\"commitment_slot\":\"0x1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600\",\"profile_version\":\"avalanche_attestor_v1\",\"l2_header_fork\":\"coreth\"}}\n";

    #[test]
    fn fuji_profile_preserves_current_payload_bytes() {
        let bytes = include_bytes!("../config/fuji.json");
        assert_eq!(bytes, FUJI_PROFILE_BYTES);
        let profile: AvalancheProfile = serde_json::from_slice(bytes).unwrap();
        assert_eq!(profile.common.profile_version, "avalanche_attestor_v1");
        assert_eq!(profile.common.l2_header_fork, ExecutionHeaderFork::Coreth);
    }
}
