//! Deterministic fixture decoding.
use serde::de::DeserializeOwned;
use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};

use crate::error::Error;

/// Metadata adjacent to a fixture.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct FixtureMetadata {
    /// Envelope version.
    pub schema_version: u32,
    /// Runtime profile label used by capture tooling.
    pub config_id: String,
    /// Fixed numeric Ethereum execution block tag.
    pub l1_block_number: u64,
    /// Hash returned for the fixed Ethereum execution block.
    pub l1_block_hash: String,
    /// Fixed numeric L2 execution block tag.
    pub l2_block_number: u64,
    /// Canonical L2 execution block hash.
    pub l2_block_hash: String,
    /// Source revision used for Ethereum proof capture.
    pub ethereum_source_revision: String,
    /// Rollup-contract source revision used to review storage and bytecode offsets.
    pub rollup_source_revision: String,
    /// Beacon slot that authenticates the captured L1 execution state.
    pub l1_beacon_slot: u64,
    /// Lowercase SHA-256 digest of the exact runtime profile JSON.
    pub profile_sha256: String,
    /// Lowercase SHA-256 digest of the fixture JSON.
    pub fixture_sha256: String,
}

impl FixtureMetadata {
    /// Validates required fixture identity and checksum fields.
    pub fn validate(&self) -> Result<(), Error> {
        if self.schema_version != 2
            || self.config_id.is_empty()
            || self.l1_block_number == 0
            || self.l1_beacon_slot == 0
            || self.l2_block_number == 0
            || self.ethereum_source_revision.is_empty()
            || self.rollup_source_revision.is_empty()
            || !valid_hash(&self.l1_block_hash)
            || !valid_hash(&self.l2_block_hash)
            || !valid_digest(&self.profile_sha256)
        {
            return Err(Error::InvalidFixtureProvenance(
                "fixture metadata is incomplete",
            ));
        }
        if !valid_digest(&self.fixture_sha256) {
            return Err(Error::InvalidFixtureProvenance(
                "fixture SHA-256 is invalid",
            ));
        }
        Ok(())
    }
}

fn valid_hash(value: &str) -> bool {
    value.len() == 66
        && value.starts_with("0x")
        && value[2..].bytes().all(|byte| byte.is_ascii_hexdigit())
}

fn valid_digest(value: &str) -> bool {
    value.len() == 64
        && value
            .bytes()
            .all(|byte| byte.is_ascii_digit() || (b'a'..=b'f').contains(&byte))
}

/// Decodes a checked-in fixture after checksum and label validation.
pub fn load_fixture<T: DeserializeOwned>(
    fixture_json: &[u8],
    metadata_json: &[u8],
) -> Result<(T, FixtureMetadata), Error> {
    let metadata: FixtureMetadata = serde_json::from_slice(metadata_json)?;
    metadata.validate()?;
    if format!("{:x}", Sha256::digest(fixture_json)) != metadata.fixture_sha256 {
        return Err(Error::FixtureChecksumMismatch);
    }
    Ok((serde_json::from_slice(fixture_json)?, metadata))
}
