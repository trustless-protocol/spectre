//!  query helpers.

use crate::{
    error::Error,
    state::{ConsensusState, FinalityPolicy, Header, ProposalStatus},
};

/// Validates a header without mutating state.
pub fn verify_client_message(header: &Header) -> Result<ConsensusState, Error> {
    header.consensus_state()
}

/// Returns the exact nanosecond timestamp recorded in a consensus state.
#[must_use]
pub const fn timestamp_at_height(consensus: &ConsensusState) -> u64 {
    consensus.timestamp_nanos
}

/// Reports whether two valid, membership-trusted headers conflict at the same height.
pub fn check_for_misbehaviour(
    policy: &FinalityPolicy,
    first: &Header,
    second: &Header,
) -> Result<bool, Error> {
    let first_state = first.consensus_state()?;
    let second_state = second.consensus_state()?;
    Ok(first.height == second.height
        && consensus_states_conflict(&first_state, &second_state)
        && consensus_state_is_trusted(policy, &first_state)
        && consensus_state_is_trusted(policy, &second_state))
}

fn consensus_states_conflict(first: &ConsensusState, second: &ConsensusState) -> bool {
    first.l2_block_hash != second.l2_block_hash
        || first.state_root != second.state_root
        || first.ibc_storage_root != second.ibc_storage_root
        || first.parent_hash != second.parent_hash
        || first.l1_origin_number != second.l1_origin_number
        || first.l1_origin_hash != second.l1_origin_hash
        || first.rollup_commitment != second.rollup_commitment
}

fn consensus_state_is_trusted(policy: &FinalityPolicy, state: &ConsensusState) -> bool {
    state.proposal_status != ProposalStatus::ResolvedInvalid
        && state.finality_level >= policy.minimum_membership_level
        && (!policy.require_resolved_proposal
            || state.proposal_status == ProposalStatus::ResolvedValid)
}
