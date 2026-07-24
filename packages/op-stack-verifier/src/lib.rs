//! Shared optimistic verification for Base and OP Stack dispute games.

#![deny(
    clippy::nursery,
    clippy::pedantic,
    warnings,
    missing_docs,
    unused_crate_dependencies
)]

use alloy_primitives::{keccak256, Address, B256};
use l2_client::{
    canonical_header::{CanonicalEvmHeader, ExecutionHeaderFork},
    error::Error,
    evm_proof::{verify_bounded_account, ProofLimits},
    msg::{EvmAccountProof, EvmStorageProof},
    state::{CommonProfile, Header as VerifiedHeader, Height, RuntimeProfile},
    verification::{
        address_from_storage_value, storage_word, verify_account_runtime, verify_storage,
    },
};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

/// Conservative bound applied to every proof in an OP Stack update.
pub const PROOF_LIMITS: ProofLimits = ProofLimits {
    max_account_nodes: 64,
    max_storage_nodes: 64,
    max_storage_proofs: 2,
    max_node_bytes: 32 * 1024,
};

/// Output-root commitment format selected by a runtime profile.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum OutputRootFormat {
    /// `keccak256(version || stateRoot || messagePasserStorageRoot || blockHash)`.
    LegacyV0,
}

/// Immutable runtime profile for one OP Stack deployment.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct Profile {
    /// Chain/client/router data shared by all optimistic L2 profiles.
    pub common: CommonProfile,
    /// L1 `DisputeGameFactory` account.
    #[schemars(with = "String")]
    pub dispute_game_factory: Address,
    /// Solidity storage slot of the factory's game list.
    #[schemars(with = "String")]
    pub game_list_slot: B256,
    /// Byte offset of the root claim in authenticated game runtime code.
    pub root_claim_bytecode_offset: u32,
    /// Output-root preimage format committed by games in this profile.
    pub output_root_format: OutputRootFormat,
    /// Canonical L2 execution-header fork encoded by update headers.
    pub l2_header_fork: ExecutionHeaderFork,
}

impl RuntimeProfile for Profile {
    fn common(&self) -> &CommonProfile {
        &self.common
    }
}

/// Supplied preimage of an OP Stack output-root commitment.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct OutputRootProof {
    /// Commitment format version; zero for V0.
    #[schemars(with = "String")]
    pub version: B256,
    /// Canonical L2 execution state root.
    #[schemars(with = "String")]
    pub state_root: B256,
    /// L2-to-L1 message passer storage root.
    #[schemars(with = "String")]
    pub message_passer_storage_root: B256,
    /// Canonical L2 execution block hash.
    #[schemars(with = "String")]
    pub latest_blockhash: B256,
}

impl OutputRootProof {
    /// Hashes the four fixed-width V0 output-root words.
    #[must_use]
    pub fn hash(self) -> B256 {
        let mut preimage = [0_u8; 128];
        preimage[..32].copy_from_slice(self.version.as_slice());
        preimage[32..64].copy_from_slice(self.state_root.as_slice());
        preimage[64..96].copy_from_slice(self.message_passer_storage_root.as_slice());
        preimage[96..].copy_from_slice(self.latest_blockhash.as_slice());
        keccak256(preimage)
    }
}

/// Optimistic update proving that a selected game commitment exists in finalized L1 state.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct Header {
    /// Exact Ethereum beacon slot used by every L1 proof.
    pub beacon_slot: u64,
    /// Ethereum execution state root authenticated at `beacon_slot`.
    #[schemars(with = "String")]
    pub l1_state_root: B256,
    /// Factory account proof.
    pub factory_proof: EvmAccountProof,
    /// Caller-selected index in the factory game list.
    pub game_index: u64,
    /// Proof of the selected game-list element.
    pub game_proof: EvmStorageProof,
    /// Account proof for the game address extracted from `game_proof`.
    pub game_account_proof: EvmAccountProof,
    /// Full game runtime returned by `eth_getCode` at the fixed L1 block.
    pub game_runtime: Vec<u8>,
    /// Preimage whose hash must equal the root claim embedded in `game_runtime`.
    pub output_root_proof: OutputRootProof,
    /// Canonical L2 execution header committed by the output root.
    pub l2_header: CanonicalEvmHeader,
    /// Account proof for the profile's L2 router.
    pub router_proof: EvmAccountProof,
}

/// Verifies an unresolved or resolved OP Stack game's commitment and L2 router state.
///
/// Game type, implementation, pause, blacklist, retirement, result, and resolution time are
/// intentionally outside this optimistic client's security model.
pub fn verify(
    profile: &Profile,
    authenticated_l1_root: B256,
    _authenticated_l1_timestamp: u64,
    header: &Header,
) -> Result<VerifiedHeader, Error> {
    if header.l1_state_root != authenticated_l1_root {
        return Err(Error::Proof(
            "OP Stack header L1 root differs from host consensus state".into(),
        ));
    }
    let factory = verify_bounded_account(
        PROOF_LIMITS,
        header.l1_state_root,
        profile.dispute_game_factory,
        &header.factory_proof,
    )?;
    let expected_game_slot = dynamic_array_element_slot(profile.game_list_slot, header.game_index)?;
    if header.game_proof.key != expected_game_slot {
        return Err(Error::Proof("unexpected factory game-list slot".into()));
    }
    verify_storage(PROOF_LIMITS, &factory, &header.game_proof)?;
    let game = address_from_storage_value(&header.game_proof.value)?;
    if game.is_zero() {
        return Err(Error::Proof("factory game-list entry is empty".into()));
    }
    let _game_account = verify_account_runtime(
        PROOF_LIMITS,
        header.l1_state_root,
        game,
        &header.game_account_proof,
        &header.game_runtime,
    )?;
    let offset = usize::try_from(profile.root_claim_bytecode_offset)
        .map_err(|_| Error::Proof("root-claim offset exceeds address space".into()))?;
    let end = offset
        .checked_add(32)
        .ok_or_else(|| Error::Proof("root-claim offset overflow".into()))?;
    let root_claim = B256::from_slice(
        header
            .game_runtime
            .get(offset..end)
            .ok_or_else(|| Error::Proof("root claim lies outside game runtime".into()))?,
    );

    match profile.output_root_format {
        OutputRootFormat::LegacyV0 if !header.output_root_proof.version.is_zero() => {
            return Err(Error::Proof("V0 output-root version is nonzero".into()));
        }
        OutputRootFormat::LegacyV0 => {}
    }
    header.l2_header.validate_for_fork(profile.l2_header_fork)?;
    let block_hash = header.l2_header.hash(profile.l2_header_fork)?;
    if header.output_root_proof.state_root != header.l2_header.state_root()
        || header.output_root_proof.latest_blockhash != block_hash
        || header.output_root_proof.hash() != root_claim
    {
        return Err(Error::Proof(
            "game root claim does not bind the canonical L2 header".into(),
        ));
    }
    let router = verify_bounded_account(
        PROOF_LIMITS,
        header.l2_header.state_root(),
        profile.common.l2_router,
        &header.router_proof,
    )?;
    Ok(VerifiedHeader {
        height: Height::new(header.l2_header.number())?,
        state_root: header.l2_header.state_root(),
        router_storage_root: B256::from(router.storage_root.0),
        timestamp_seconds: header.l2_header.timestamp(),
    })
}

fn dynamic_array_element_slot(slot: B256, index: u64) -> Result<B256, Error> {
    let mut value: [u8; 32] = keccak256(storage_word(slot.as_slice())?).into();
    let mut carry = index;
    for byte in value.iter_mut().rev() {
        let sum = u64::from(*byte) + (carry & 0xff);
        *byte = sum as u8;
        carry = (carry >> 8) + (sum >> 8);
    }
    if carry != 0 {
        return Err(Error::Proof("game-list slot arithmetic overflow".into()));
    }
    Ok(B256::from(value))
}

#[cfg(test)]
mod tests {
    use alloy_primitives::B256;

    use super::{dynamic_array_element_slot, OutputRootProof};

    #[test]
    fn output_root_and_game_slots_are_deterministic() {
        let proof = OutputRootProof {
            version: B256::ZERO,
            state_root: B256::with_last_byte(1),
            message_passer_storage_root: B256::with_last_byte(2),
            latest_blockhash: B256::with_last_byte(3),
        };
        assert_ne!(proof.hash(), B256::ZERO);
        assert_ne!(
            dynamic_array_element_slot(B256::with_last_byte(7), 4).unwrap(),
            B256::ZERO
        );
    }
}
