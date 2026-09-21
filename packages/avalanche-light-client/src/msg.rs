//! Wire messages. The host-facing scalar messages (sudo/query envelopes, IBC
//! height, Merkle path, proofs) are the shared ones from `l2-client`; only the
//! header and migration shapes are Avalanche-specific.

use cosmwasm_std::Binary;
use l2_client::canonical_header::CanonicalEvmHeader;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

pub use l2_client::msg::{
    CheckForMisbehaviourResult, EvmAccountProof, EvmStorageProof, IbcHeight, InstantiateMsg,
    MerklePath, QueryMsg, StatusResult, SudoMsg, TimestampAtHeightResult, UpdateStateResult,
};

use crate::state::{ValidatorSetState, WireValidator};

/// One warp-attested C-Chain update.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct WarpSignedHeader {
    /// Canonical coreth execution header (fork `coreth`).
    pub header: CanonicalEvmHeader,
    /// Signer bitset over the pinned canonical validator set.
    pub signer_bit_set: Binary,
    /// 96-byte aggregate BLS signature over the block-hash warp message.
    pub signature: Binary,
    /// Account proof of the router against `header.state_root`.
    pub router_proof: EvmAccountProof,
}

/// One unambiguous client-message envelope (shared serde shape:
/// `{"type": "header"|"misbehaviour", "value": ...}`).
pub type ClientMessage = l2_client::msg::ClientMessage<WarpSignedHeader>;

/// Governance-hosted migration operations.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(rename_all = "snake_case", deny_unknown_fields)]
pub enum MigrateMsg {
    /// Validate the current client without rewriting it.
    KeepValidatorSet {},
    /// Atomically replace the pinned canonical validator set with a newer
    /// P-Chain snapshot. Governance-gated rotation is this version's trust
    /// model for set changes; the snapshot height must strictly increase, so
    /// governance itself cannot roll the set back.
    ReplaceValidatorSet {
        /// Canonically ordered validators with BLS keys.
        validators: Vec<WireValidator>,
        /// Total snapshot weight, keyless validators included.
        total_weight: u64,
        /// P-Chain height of the snapshot.
        p_chain_height: u64,
    },
}

impl MigrateMsg {
    /// Converts a replacement into stored state.
    #[must_use]
    pub fn replacement(self) -> Option<ValidatorSetState> {
        match self {
            Self::KeepValidatorSet {} => None,
            Self::ReplaceValidatorSet {
                validators,
                total_weight,
                p_chain_height,
            } => Some(ValidatorSetState {
                validators,
                total_weight,
                p_chain_height,
            }),
        }
    }
}
