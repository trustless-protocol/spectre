//! Local validation shared by every L2 adapter.

use alloy_primitives::B256;

use crate::{
    error::Error,
    evm_proof::{verify_bounded_account, ProofLimits},
    msg::AttestedL2Header,
    state::{Header, Height, RuntimeProfile},
};

/// Conservative bound for a router account proof.
pub const ROUTER_PROOF_LIMITS: ProofLimits = ProofLimits {
    max_account_nodes: 64,
    max_storage_nodes: 64,
    max_node_bytes: 32 * 1024,
};

/// Validates a common attested header and derives every stored root.
pub fn verify_attested_header<P: RuntimeProfile>(
    profile: &P,
    header: &AttestedL2Header,
) -> Result<Header, Error> {
    let common = profile.common();
    if common.profile_version != P::expected_profile_version() {
        return Err(Error::InvalidProfileVersion);
    }

    header.l2_header.validate_for_fork(common.l2_header_fork)?;
    let block_hash = header.l2_header.hash(common.l2_header_fork)?;
    let router = verify_bounded_account(
        ROUTER_PROOF_LIMITS,
        header.l2_header.state_root(),
        common.l2_router,
        &header.router_proof,
    )?;

    Ok(Header {
        height: Height::new(header.l2_header.number())?,
        state_root: header.l2_header.state_root(),
        router_storage_root: B256::from(router.storage_root.0),
        timestamp_seconds: header.l2_header.timestamp(),
        l2_block_hash: block_hash,
        parent_hash: header.l2_header.parent_hash,
    })
}
