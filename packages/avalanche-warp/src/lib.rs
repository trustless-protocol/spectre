//! Avalanche warp aggregate-signature verification.
//!
//! The verifiable primitive Avalanche exposes is a stake-weighted BLS
//! aggregate: validators (BLS keys registered on the P-Chain) sign an unsigned
//! warp message, and a relayer aggregates the signatures into a 96-byte
//! BLS12-381 G2 signature plus a bitset over the *canonical validator set*.
//! This crate implements the verifying side exactly as avalanchego does
//! (`vms/platformvm/warp/signature.go`):
//!
//! 1. parse the signer bitset, rejecting non-minimal encodings;
//! 2. select the signing validators from the stored canonical set, rejecting
//!    any bit beyond it;
//! 3. sum the signers' weights (overflow-checked) and require
//!    `signed * quorum_den >= total * quorum_num`, where `total` includes
//!    validators without BLS keys (they dilute the quorum);
//! 4. verify the aggregate BLS signature over the raw unsigned message bytes.
//!
//! The BLS pairing itself is behind [`BlsVerify`], so the CosmWasm client can
//! plug in the chain's custom querier (the same `aggregate_verify` shape the
//! Ethereum client uses — Avalanche and Ethereum share the BLS ciphersuite and
//! domain-separation tag) or native host functions, and tests plug in a real
//! in-process implementation. Everything else in this crate is pure.
//!
//! Avalanche is an L1; this crate is chain machinery for verifying it, shared
//! by nothing else.

#![deny(
    clippy::nursery,
    clippy::pedantic,
    warnings,
    missing_docs,
    unused_crate_dependencies
)]
#![allow(clippy::missing_errors_doc)]
#![cfg_attr(not(test), no_std)]

extern crate alloc;

use alloc::vec::Vec;

/// Compressed BLS12-381 G1 public key length (the form the P-Chain API serves
/// and the chain-side BLS querier consumes).
pub const PUBLIC_KEY_LEN: usize = 48;
/// BLS12-381 G2 signature length.
pub const SIGNATURE_LEN: usize = 96;
/// Warp codec version prefix (avalanchego codec version 0).
const CODEC_VERSION: [u8; 2] = [0, 0];
/// `Hash` payload type id in the warp payload codec.
const HASH_PAYLOAD_TYPE_ID: [u8; 4] = [0, 0, 0, 0];
/// `BitSetSignature` type id in the warp message codec.
const BIT_SET_SIGNATURE_TYPE_ID: [u8; 4] = [0, 0, 0, 0];

/// One entry of the stored canonical validator set, in canonical order.
///
/// The canonical order (ascending uncompressed-key bytes, duplicate keys
/// merged) is established by the P-Chain when the set is snapshotted; the
/// client stores the set in that order and never re-sorts, because the signer
/// bitset indexes into it.
#[derive(Clone, Debug, PartialEq, Eq)]
pub struct Validator {
    /// Compressed BLS public key.
    pub public_key: [u8; PUBLIC_KEY_LEN],
    /// Stake weight.
    pub weight: u64,
}

/// Verification errors.
#[derive(Debug, PartialEq, Eq, thiserror::Error)]
pub enum Error {
    /// The bitset is not the minimal encoding of its integer (avalanchego
    /// rejects padded bitsets, so a padded bitset must not verify here either).
    #[error("signer bitset is not minimally encoded")]
    NonMinimalBitSet,
    /// A bitset bit refers past the end of the canonical set.
    #[error("signer index {0} is outside the canonical validator set")]
    UnknownValidator(usize),
    /// No signer bit is set.
    #[error("no signers")]
    NoSigners,
    /// Weight arithmetic overflowed.
    #[error("weight overflow")]
    WeightOverflow,
    /// The signers do not reach the required stake quorum.
    #[error("signed weight {signed} of {total} is below quorum {quorum_num}/{quorum_den}")]
    InsufficientQuorum {
        /// Weight that signed.
        signed: u64,
        /// Total weight of the set, including keyless validators.
        total: u64,
        /// Quorum numerator.
        quorum_num: u64,
        /// Quorum denominator.
        quorum_den: u64,
    },
    /// The aggregate signature did not verify over the message.
    #[error("aggregate BLS signature does not verify")]
    InvalidSignature,
    /// A signed warp message envelope could not be parsed.
    #[error("malformed warp message: {0}")]
    Malformed(&'static str),
}

/// BLS backend: verify an aggregate signature by `public_keys` over `message`.
///
/// Implementations must use the standard BLS proof-of-possession ciphersuite
/// (`BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_`) — the one ciphersuite both
/// Avalanche warp and Ethereum consensus use.
pub trait BlsVerify {
    /// Returns `Ok(())` iff the aggregate of `public_keys` signed `message`.
    fn fast_aggregate_verify(
        &self,
        public_keys: &[&[u8; PUBLIC_KEY_LEN]],
        message: &[u8],
        signature: &[u8; SIGNATURE_LEN],
    ) -> Result<(), Error>;
}

/// Builds the unsigned warp message for a `Hash` payload — the exact bytes
/// Avalanche validators sign when attesting an accepted C-Chain block hash.
///
/// Layout (avalanchego linear codec): `codec(2)=0 ‖ network_id(4,BE) ‖
/// source_chain_id(32) ‖ payload_len(4,BE) ‖ payload`, with the payload being
/// `codec(2)=0 ‖ type_id(4)=0 (Hash) ‖ hash(32)`.
#[must_use]
pub fn build_block_hash_message(
    network_id: u32,
    source_chain_id: &[u8; 32],
    block_hash: &[u8; 32],
) -> Vec<u8> {
    let mut payload = Vec::with_capacity(38);
    payload.extend_from_slice(&CODEC_VERSION);
    payload.extend_from_slice(&HASH_PAYLOAD_TYPE_ID);
    payload.extend_from_slice(block_hash);

    let mut message = Vec::with_capacity(42 + payload.len());
    message.extend_from_slice(&CODEC_VERSION);
    message.extend_from_slice(&network_id.to_be_bytes());
    message.extend_from_slice(source_chain_id);
    message.extend_from_slice(&u32::try_from(payload.len()).expect("38").to_be_bytes());
    message.extend_from_slice(&payload);
    message
}

/// The signer bitset of a `BitSetSignature`, byte-order compatible with Go's
/// `big.Int.Bytes()` (big-endian): bit `i` of the integer is bit `i % 8` of
/// byte `bytes[len - 1 - i / 8]`.
#[derive(Clone, Debug)]
pub struct SignerBitSet<'a>(&'a [u8]);

impl<'a> SignerBitSet<'a> {
    /// Parses a bitset, rejecting non-minimal encodings (a leading zero byte),
    /// mirroring avalanchego's re-encode check.
    pub fn parse(bytes: &'a [u8]) -> Result<Self, Error> {
        if bytes.first() == Some(&0) {
            return Err(Error::NonMinimalBitSet);
        }
        Ok(Self(bytes))
    }

    /// Reports whether bit `i` is set.
    #[must_use]
    pub fn contains(&self, i: usize) -> bool {
        let byte_index = i / 8;
        if byte_index >= self.0.len() {
            return false;
        }
        (self.0[self.0.len() - 1 - byte_index] >> (i % 8)) & 1 == 1
    }

    /// The highest representable bit index + 1.
    #[must_use]
    pub const fn capacity(&self) -> usize {
        self.0.len() * 8
    }
}

/// Selects the signing validators and enforces the stake quorum. Returns the
/// signers (in canonical order) and their summed weight.
///
/// `total_weight` is the stored total of the snapshotted set — including any
/// validators without BLS keys, which are absent from `validators` but still
/// dilute the quorum, exactly as in avalanchego's `FlattenValidatorSet`.
pub fn verify_quorum<'a>(
    validators: &'a [Validator],
    total_weight: u64,
    bit_set: &SignerBitSet<'_>,
    quorum_num: u64,
    quorum_den: u64,
) -> Result<(Vec<&'a Validator>, u64), Error> {
    // Any set bit past the canonical set is an unknown validator, checked over
    // the whole bitset capacity, not only over set members.
    for i in validators.len()..bit_set.capacity() {
        if bit_set.contains(i) {
            return Err(Error::UnknownValidator(i));
        }
    }

    let mut signers = Vec::new();
    let mut signed_weight: u64 = 0;
    for (i, validator) in validators.iter().enumerate() {
        if bit_set.contains(i) {
            signed_weight = signed_weight
                .checked_add(validator.weight)
                .ok_or(Error::WeightOverflow)?;
            signers.push(validator);
        }
    }
    if signers.is_empty() {
        return Err(Error::NoSigners);
    }

    // signed/total >= num/den without division: u128 keeps u64*u64 exact.
    if u128::from(signed_weight) * u128::from(quorum_den)
        < u128::from(total_weight) * u128::from(quorum_num)
    {
        return Err(Error::InsufficientQuorum {
            signed: signed_weight,
            total: total_weight,
            quorum_num,
            quorum_den,
        });
    }
    Ok((signers, signed_weight))
}

/// Verifies a warp aggregate: quorum first (cheap), then the BLS pairing over
/// the raw unsigned message bytes. Returns the signed weight.
pub fn verify_warp_aggregate<B: BlsVerify>(
    validators: &[Validator],
    total_weight: u64,
    unsigned_message: &[u8],
    bit_set_bytes: &[u8],
    signature: &[u8; SIGNATURE_LEN],
    quorum_num: u64,
    quorum_den: u64,
    bls: &B,
) -> Result<u64, Error> {
    let bit_set = SignerBitSet::parse(bit_set_bytes)?;
    let (signers, signed_weight) =
        verify_quorum(validators, total_weight, &bit_set, quorum_num, quorum_den)?;
    let keys: Vec<&[u8; PUBLIC_KEY_LEN]> = signers.iter().map(|v| &v.public_key).collect();
    bls.fast_aggregate_verify(&keys, unsigned_message, signature)?;
    Ok(signed_weight)
}

/// A signed warp message split into its three parts.
#[derive(Clone, Debug, PartialEq, Eq)]
pub struct ParsedSignedMessage<'a> {
    /// The unsigned message bytes (what validators signed).
    pub unsigned_message: &'a [u8],
    /// The signer bitset bytes.
    pub bit_set: &'a [u8],
    /// The 96-byte aggregate signature.
    pub signature: [u8; SIGNATURE_LEN],
}

/// Splits a signed warp message (`unsigned ‖ BitSetSignature`) into its parts,
/// mirroring avalanchego's message codec. Used by the relayer and by tests;
/// the wasm client receives the parts pre-split.
pub fn parse_signed_message(bytes: &[u8]) -> Result<ParsedSignedMessage<'_>, Error> {
    // Unsigned envelope: codec(2) ‖ network_id(4) ‖ chain_id(32) ‖ len(4) ‖ payload.
    if bytes.len() < 42 {
        return Err(Error::Malformed("shorter than the unsigned envelope"));
    }
    if bytes[0..2] != CODEC_VERSION {
        return Err(Error::Malformed("unsupported codec version"));
    }
    let payload_len = u32::from_be_bytes(
        bytes[38..42]
            .try_into()
            .map_err(|_| Error::Malformed("payload length"))?,
    ) as usize;
    let unsigned_len = 42 + payload_len;

    // Signature: type_id(4)=0 (BitSetSignature) ‖ bitset_len(4) ‖ bitset ‖ sig(96).
    let sig_part = bytes
        .get(unsigned_len..)
        .ok_or(Error::Malformed("truncated after the unsigned message"))?;
    if sig_part.len() < 8 {
        return Err(Error::Malformed("truncated signature header"));
    }
    if sig_part[0..4] != BIT_SET_SIGNATURE_TYPE_ID {
        return Err(Error::Malformed("unsupported signature type"));
    }
    let bit_set_len = u32::from_be_bytes(
        sig_part[4..8]
            .try_into()
            .map_err(|_| Error::Malformed("bitset length"))?,
    ) as usize;
    let expected = 8 + bit_set_len + SIGNATURE_LEN;
    if sig_part.len() != expected {
        return Err(Error::Malformed("signature section length mismatch"));
    }
    let signature: [u8; SIGNATURE_LEN] = sig_part[8 + bit_set_len..]
        .try_into()
        .map_err(|_| Error::Malformed("signature length"))?;
    Ok(ParsedSignedMessage {
        unsigned_message: &bytes[..unsigned_len],
        bit_set: &sig_part[8..8 + bit_set_len],
        signature,
    })
}

#[cfg(test)]
mod tests;
