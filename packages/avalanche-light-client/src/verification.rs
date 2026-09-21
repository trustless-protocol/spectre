//! Header verification: coreth hash → warp quorum → router proof.

use alloy_primitives::B256;
use l2_client::{
    canonical_header::ExecutionHeaderFork,
    evm_proof::{verify_bounded_account, ProofLimits},
    state::Height,
};

use crate::{
    error::Error,
    msg::WarpSignedHeader,
    state::{ClientState, Header},
};

/// Conservative bound for a router account proof (same budget as the shared
/// attested clients).
pub const ROUTER_PROOF_LIMITS: ProofLimits = ProofLimits {
    max_account_nodes: 64,
    max_storage_nodes: 64,
    max_node_bytes: 32 * 1024,
};

/// A header that crossed quorum authentication and proof verification.
///
/// The field and constructor are private: runtime transitions cannot be
/// invoked with a deserialized or hand-built header.
#[derive(Clone, Debug, PartialEq, Eq)]
pub struct AuthenticatedHeader {
    header: Header,
}

impl AuthenticatedHeader {
    const fn new(header: Header) -> Self {
        Self { header }
    }

    pub(crate) const fn header(&self) -> &Header {
        &self.header
    }

    #[cfg(test)]
    pub(crate) const fn from_test_header(header: Header) -> Self {
        Self::new(header)
    }
}

/// Verifies one warp-signed header end to end:
///
/// 1. shape: coreth fork validation, height, proof bounds, signature lengths;
/// 2. identity: recompute the coreth block hash from the header fields;
/// 3. authority: rebuild the block-hash warp message and verify the
///    stake-weighted aggregate against the pinned canonical set;
/// 4. state: prove the router account against the header's state root.
pub fn verify_signed_header<B: avalanche_warp::BlsVerify>(
    client: &ClientState,
    signed: &WarpSignedHeader,
    bls: &B,
) -> Result<AuthenticatedHeader, Error> {
    client.validate()?;
    let block_hash = prepare_and_authenticate(client, signed, bls)?;
    verify_router_proof(client, signed, block_hash)
}

/// Authenticates both headers (quorum first) before traversing either proof.
pub fn verify_signed_misbehaviour<B: avalanche_warp::BlsVerify>(
    client: &ClientState,
    first: &WarpSignedHeader,
    second: &WarpSignedHeader,
    bls: &B,
) -> Result<(AuthenticatedHeader, AuthenticatedHeader), Error> {
    client.validate()?;
    let first_hash = prepare_and_authenticate(client, first, bls)?;
    let second_hash = prepare_and_authenticate(client, second, bls)?;
    Ok((
        verify_router_proof(client, first, first_hash)?,
        verify_router_proof(client, second, second_hash)?,
    ))
}

fn prepare_and_authenticate<B: avalanche_warp::BlsVerify>(
    client: &ClientState,
    signed: &WarpSignedHeader,
    bls: &B,
) -> Result<B256, Error> {
    signed
        .header
        .validate_for_fork(ExecutionHeaderFork::Coreth)
        .map_err(Error::Evm)?;
    Height::new(signed.header.number()).map_err(Error::Evm)?;
    ROUTER_PROOF_LIMITS
        .validate_account_bounds(&signed.router_proof)
        .map_err(Error::Evm)?;
    let signature: &[u8; avalanche_warp::SIGNATURE_LEN] = signed
        .signature
        .as_slice()
        .try_into()
        .map_err(|_| Error::Invalid("aggregate signature must be 96 bytes"))?;

    let block_hash = signed
        .header
        .hash(ExecutionHeaderFork::Coreth)
        .map_err(Error::Evm)?;
    let message = avalanche_warp::build_block_hash_message(
        client.network_id,
        &client.source_chain_id_bytes()?,
        &block_hash.into(),
    );
    let validators = client.validator_set.to_warp_validators()?;
    avalanche_warp::verify_warp_aggregate(
        &validators,
        client.validator_set.total_weight,
        &message,
        signed.signer_bit_set.as_slice(),
        signature,
        client.quorum_num,
        client.quorum_den,
        bls,
    )?;
    Ok(block_hash)
}

fn verify_router_proof(
    client: &ClientState,
    signed: &WarpSignedHeader,
    block_hash: B256,
) -> Result<AuthenticatedHeader, Error> {
    let router = verify_bounded_account(
        ROUTER_PROOF_LIMITS,
        signed.header.state_root(),
        client.router,
        &signed.router_proof,
    )
    .map_err(Error::Evm)?;

    Ok(AuthenticatedHeader::new(Header {
        height: signed.header.number(),
        state_root: signed.header.state_root(),
        router_storage_root: B256::from(router.storage_root.0),
        timestamp_seconds: signed.header.timestamp(),
        block_hash,
        parent_hash: signed.header.parent_hash,
    }))
}
