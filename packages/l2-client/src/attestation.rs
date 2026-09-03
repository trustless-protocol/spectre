//! Canonical bytes for an attestor's signed L2 block identity.

use alloy_primitives::B256;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

/// Ed25519 domain separator shared with `attestor/types/attestation`.
///
/// This deployed wire-format identifier must not follow a Fast-IBC/Spectre
/// branding rename: changing it invalidates existing client profiles and their
/// signatures.
const DOMAIN: &[u8] = b"fast-ibc/l2-attestation/v2\0";

/// Immutable finality policy for one attestor key.
///
/// Its byte discriminants are included in every signature so a key that also
/// serves a lower-latency relay path cannot authorize that state for this
/// client.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum AttestationHead {
    /// The sequencer's unsafe head.
    Unsafe,
    /// L1-derived but reorgable state.
    Safe,
    /// L1-finalized state.
    Finalized,
}

impl AttestationHead {
    /// Stable byte encoding shared with the Go signer.
    #[must_use]
    pub const fn signing_byte(self) -> u8 {
        match self {
            Self::Unsafe => 1,
            Self::Safe => 2,
            Self::Finalized => 3,
        }
    }
}

/// Produces the exact bytes an attestor signs for an L2 execution block.
#[must_use]
pub fn signing_bytes(
    l2_chain_id: u64,
    attestation_head: AttestationHead,
    l2_block_number: u64,
    state_root: B256,
    block_hash: B256,
) -> Vec<u8> {
    let mut message = Vec::with_capacity(DOMAIN.len() + 8 + 1 + 8 + 32 + 32);
    message.extend_from_slice(DOMAIN);
    message.extend_from_slice(&l2_chain_id.to_be_bytes());
    message.push(attestation_head.signing_byte());
    message.extend_from_slice(&l2_block_number.to_be_bytes());
    message.extend_from_slice(state_root.as_slice());
    message.extend_from_slice(block_hash.as_slice());
    message
}
