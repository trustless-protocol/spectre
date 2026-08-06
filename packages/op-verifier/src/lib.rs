//! Optimism adapter for the shared attestor-trusted L2 client.
#![deny(clippy::nursery, clippy::pedantic, warnings, unused_crate_dependencies)]
use l2_client::{
    msg::AttestedL2Header,
    state::{CommonProfile, Header, RuntimeProfile},
};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct Profile {
    pub common: CommonProfile,
}
impl RuntimeProfile for Profile {
    fn common(&self) -> &CommonProfile {
        &self.common
    }
    fn expected_profile_version() -> &'static str {
        "op_attestor_v1"
    }
}
pub struct Adapter;
impl l2_client::L2LightClient for Adapter {
    type Profile = Profile;
    fn verify(
        profile: &Profile,
        header: &AttestedL2Header,
    ) -> Result<Header, l2_client::error::Error> {
        l2_client::verification::verify_attested_header(profile, header)
    }
}

#[cfg(test)]
mod tests {
    use super::Profile;

    #[test]
    fn example_profile_matches_artifact_version() {
        let profile: Profile =
            serde_json::from_str(include_str!("../config/op-sepolia.json")).unwrap();
        assert_eq!(profile.common.profile_version, "op_attestor_v1");
    }
}
