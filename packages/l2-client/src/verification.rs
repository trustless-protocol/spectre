//! Common account and storage checks for optimistic rollup verifiers.

use alloy_primitives::{keccak256, Address, B256};
use ethereum_trie_db::trie_db::Account;

use crate::{
    error::Error,
    evm_proof::{verify_bounded_account, verify_bounded_storage_value, ProofLimits},
    msg::{EvmAccountProof, EvmStorageProof},
};

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
