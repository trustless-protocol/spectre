//!  query helpers.

use crate::{
    error::Error,
    state::{ConsensusState, Header},
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

/// Reports whether two valid headers are a same-height consensus conflict.
pub fn check_for_misbehaviour(first: &Header, second: &Header) -> Result<bool, Error> {
    let first_state = first.consensus_state()?;
    let second_state = second.consensus_state()?;
    Ok(first.height == second.height && first_state != second_state)
}
