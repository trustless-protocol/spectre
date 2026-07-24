//!  state transitions.

use crate::{
    error::Error,
    state::{ClientState, ConsensusState, Header},
};

/// Classifies and applies one update while preserving historical consensus states.
pub fn update_state<Config>(
    client: &mut ClientState<Config>,
    header: &Header,
    existing_same_height: Option<&ConsensusState>,
) -> Result<Option<ConsensusState>, Error> {
    if let Some(height) = client.frozen_height {
        return Err(Error::Frozen(height));
    }
    let consensus = header.consensus_state()?;
    let height = header.height.revision_height;
    if let Some(existing) = existing_same_height {
        // A normal update may be a relayer mistake, so a conflict is rejected but cannot freeze
        // the client. Freezing requires two independently verified headers through the explicit
        // UpdateStateOnMisbehaviour path.
        return if existing == &consensus {
            Ok(None)
        } else {
            Err(Error::Conflict(height))
        };
    }
    if height > client.latest_height {
        client.latest_height = height;
    }
    Ok(Some(consensus))
}

/// Freezes a client only after a checked same-height conflict.
pub fn freeze_on_conflict<Config>(client: &mut ClientState<Config>, height: u64) {
    client.frozen_height = Some(height);
}
