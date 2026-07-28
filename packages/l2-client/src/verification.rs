//! Common account and storage checks for optimistic rollup verifiers.

use alloy_primitives::{keccak256, Address, B256};
use ethereum_trie_db::trie_db::Account;
use serde::Serialize;
use sha2::{Digest, Sha256};

use crate::{
    error::Error,
    evm_proof::{verify_bounded_account, verify_bounded_storage_value, ProofLimits},
    msg::{EvmAccountProof, EvmStorageProof},
    msg::{FinalityEvidence, SafeEvidence},
    state::{FinalityLevel, FinalityPolicy, ProposalStatus},
};

/// Derives the accepted finality level from typed evidence and verifier-authenticated facts.
pub fn verify_finality_evidence(
    policy: &FinalityPolicy,
    evidence: &FinalityEvidence,
    l2_height: u64,
    l2_block_hash: B256,
    rollup_commitment: B256,
    l1_origin_number: u64,
    l1_origin_hash: B256,
    authenticated_l1_timestamp: u64,
    verified_proposal_status: ProposalStatus,
) -> Result<(FinalityLevel, ProposalStatus), Error> {
    match evidence {
        FinalityEvidence::Unsafe(unsafe_evidence) => {
            if unsafe_evidence.commitment != rollup_commitment {
                return Err(Error::InvalidFinalityEvidence);
            }
            Ok((FinalityLevel::Unsafe, verified_proposal_status))
        }
        FinalityEvidence::Safe(safe) => {
            verify_safe_evidence(
                policy,
                safe,
                l2_height,
                l2_block_hash,
                l1_origin_number,
                l1_origin_hash,
                authenticated_l1_timestamp,
            )?;
            Ok((FinalityLevel::Safe, verified_proposal_status))
        }
        FinalityEvidence::Finalized(finalized) => {
            if finalized.l1_origin_number != l1_origin_number
                || finalized.l1_origin_hash != l1_origin_hash
                || finalized.proposal_status != ProposalStatus::ResolvedValid
                || verified_proposal_status != ProposalStatus::ResolvedValid
            {
                return Err(Error::InvalidFinalityEvidence);
            }
            Ok((FinalityLevel::Finalized, ProposalStatus::ResolvedValid))
        }
    }
}

/// Applies independently verified typed evidence to a normalized rollup result.
pub fn apply_finality_evidence(
    policy: &FinalityPolicy,
    evidence: &FinalityEvidence,
    verified: &mut crate::state::Header,
) -> Result<(), Error> {
    let (level, proposal_status) = verify_finality_evidence(
        policy,
        evidence,
        verified.height.revision_height,
        verified.l2_block_hash,
        verified.rollup_commitment,
        verified.l1_origin_number,
        verified.l1_origin_hash,
        verified.finality_reached_at,
        verified.proposal_status,
    )?;
    verified.finality_level = level;
    verified.proposal_status = proposal_status;
    Ok(())
}

fn verify_safe_evidence(
    policy: &FinalityPolicy,
    evidence: &SafeEvidence,
    l2_height: u64,
    l2_block_hash: B256,
    l1_origin_number: u64,
    l1_origin_hash: B256,
    authenticated_l1_timestamp: u64,
) -> Result<(), Error> {
    if evidence.l2_height != l2_height
        || evidence.l2_block_hash != l2_block_hash
        || evidence.l1_origin_number != l1_origin_number
        || evidence.l1_origin_hash != l1_origin_hash
        || evidence.safe_l1_hash.is_zero()
        || evidence.valid_until < authenticated_l1_timestamp
    {
        return Err(Error::InvalidFinalityEvidence);
    }
    // The configured modes describe future proof systems, not trust assumptions. Neither an
    // authenticated non-finalized Ethereum header nor threshold signatures are verified by this
    // contract yet, so accepting either mode here would turn relayer-provided data into Safe.
    let _ = policy;
    Err(Error::SafeVerificationUnavailable)
}

/// Hashes the canonical JSON representation of authenticated rollup evidence.
pub fn evidence_hash<T: Serialize>(evidence: &T) -> Result<B256, Error> {
    let bytes = serde_json::to_vec(evidence)?;
    let digest = Sha256::digest(bytes);
    Ok(B256::from_slice(digest.as_slice()))
}

/// Maximum runtime bytecode accepted for an authenticated dispute game.
pub const MAX_RUNTIME_BYTES: usize = 16 * 1024;

/// Verifies an account and binds supplied runtime bytecode to its authenticated code hash.
pub fn verify_account_runtime(
    limits: ProofLimits,
    state_root: B256,
    address: Address,
    proof: &EvmAccountProof,
    runtime: &[u8],
) -> Result<Account, Error> {
    if runtime.len() > MAX_RUNTIME_BYTES {
        return Err(Error::Proof(
            "authenticated runtime exceeds byte limit".into(),
        ));
    }
    let account = verify_bounded_account(limits, state_root, address, proof)?;
    if B256::from(account.code_hash.0) != keccak256(runtime) {
        return Err(Error::Proof(
            "supplied runtime differs from authenticated code hash".into(),
        ));
    }
    Ok(account)
}

/// Verifies a storage value beneath an authenticated account.
pub fn verify_storage(
    limits: ProofLimits,
    account: &Account,
    proof: &EvmStorageProof,
) -> Result<(), Error> {
    verify_bounded_storage_value(limits, &account.storage_root.0, proof)
}

/// Left-pads a canonical EVM storage value to one word.
pub fn storage_word(value: &[u8]) -> Result<[u8; 32], Error> {
    if value.len() > 32 {
        return Err(Error::Proof(
            "storage value is wider than one EVM word".into(),
        ));
    }
    let mut word = [0_u8; 32];
    word[32 - value.len()..].copy_from_slice(value);
    Ok(word)
}

/// Extracts packed bytes at a Solidity offset measured from the least-significant byte.
pub fn packed_storage_field(value: &[u8], offset: usize, width: usize) -> Result<Vec<u8>, Error> {
    let word = storage_word(value)?;
    let end = 32_usize
        .checked_sub(offset)
        .ok_or_else(|| Error::Proof("packed storage offset exceeds word".into()))?;
    let start = end
        .checked_sub(width)
        .ok_or_else(|| Error::Proof("packed storage field exceeds word".into()))?;
    Ok(word[start..end].to_vec())
}

/// Derives a Solidity mapping slot for a `bytes32` key.
#[must_use]
pub fn mapping_slot_bytes32(key: B256, slot: B256) -> B256 {
    let mut preimage = [0_u8; 64];
    preimage[..32].copy_from_slice(key.as_slice());
    preimage[32..].copy_from_slice(slot.as_slice());
    keccak256(preimage)
}

/// Decodes a canonical EVM storage word into its low-order address.
pub fn address_from_storage_value(value: &[u8]) -> Result<Address, Error> {
    let word = storage_word(value)?;
    Ok(Address::from_slice(&word[12..]))
}
