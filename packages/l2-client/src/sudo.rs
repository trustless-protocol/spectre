//!  state transitions.

use crate::{
    error::Error,
    state::{ClientState, ConsensusState, FinalityLevel, Header, ProposalStatus},
};

/// Classifies and applies one update while preserving historical consensus states.
pub fn update_state<Config>(
    client: &mut ClientState<Config>,
    header: &Header,
    existing_same_height: Option<&ConsensusState>,
) -> Result<Option<ConsensusState>, Error> {
    update_state_with_parent(client, header, existing_same_height, None)
}

/// Applies an update with the previous height available for parent continuity checks.
pub fn update_state_with_parent<Config>(
    client: &mut ClientState<Config>,
    header: &Header,
    existing_same_height: Option<&ConsensusState>,
    previous: Option<&ConsensusState>,
) -> Result<Option<ConsensusState>, Error> {
    if let Some(height) = client.frozen_height {
        return Err(Error::Frozen(height));
    }
    client.finality_policy.validate()?;
    let consensus = header.consensus_state()?;
    let height = header.height.revision_height;
    if header.finality_level < client.finality_policy.minimum_update_level {
        return Err(Error::InsufficientFinality {
            required: client.finality_policy.minimum_update_level,
            actual: header.finality_level,
        });
    }

    if height > 1 {
        // Legacy consensus JSON did not include `l2_height` or the block-link fields.  Runtime
        // deserialization retains its zero default as a migration marker, so it must not be used
        // to enforce continuity that was never authenticated.
        if let Some(previous) = previous.filter(|state| state.l2_height != 0) {
            let child_trusted = is_trusted(client, &consensus);
            let parent_trusted = is_trusted(client, previous);
            if child_trusted || parent_trusted {
                if header.parent_hash != previous.l2_block_hash {
                    return Err(Error::ParentHashMismatch);
                }
            }
        }
    }

    // A legacy state has no authenticated block hash and therefore cannot establish a
    // same-height conflict.  Accepting this update upgrades the stored entry in place.
    if let Some(existing) = existing_same_height.filter(|state| state.l2_height != 0) {
        if existing.l2_block_hash != consensus.l2_block_hash {
            let existing_is_trusted = is_trusted(client, existing);
            let incoming_is_trusted = is_trusted(client, &consensus);
            if existing.finality_level == FinalityLevel::Finalized
                || (client.finality_policy.freeze_on_trusted_conflict && existing_is_trusted)
            {
                client.frozen_height = Some(height);
                return Err(Error::TrustedConflict { height });
            }
            // An untrusted optimistic candidate may be replaced by a conflicting state only
            // after the replacement itself crosses the configured trust threshold. This is an
            // expected reorg correction, not evidence that two trusted histories coexist.
            if incoming_is_trusted {
                return Ok(Some(consensus));
            }
            return Err(Error::PromotionBlockMismatch);
        }
        if existing.state_root != consensus.state_root
            || existing.ibc_storage_root != consensus.ibc_storage_root
            || existing.parent_hash != consensus.parent_hash
            || existing.l1_origin_number != consensus.l1_origin_number
            || existing.l1_origin_hash != consensus.l1_origin_hash
            || existing.rollup_commitment != consensus.rollup_commitment
        {
            return Err(Error::PromotionStateRootMismatch);
        }
        if consensus.finality_level < existing.finality_level {
            return Err(Error::FinalityDowngrade);
        }
        if !valid_proposal_transition(existing.proposal_status, consensus.proposal_status) {
            return Err(Error::InvalidFinalityEvidence);
        }
        if consensus.finality_level == existing.finality_level
            && consensus.proposal_status == existing.proposal_status
        {
            return Ok(None);
        }
        let mut promoted = consensus;
        promoted.first_accepted_at = existing.first_accepted_at;
        return Ok(Some(promoted));
    }
    if consensus.proposal_status == ProposalStatus::ResolvedInvalid {
        return Err(Error::ProposalResolvedInvalid);
    }
    if height > client.latest_height {
        client.latest_height = height;
    }
    Ok(Some(consensus))
}

fn valid_proposal_transition(existing: ProposalStatus, incoming: ProposalStatus) -> bool {
    existing == incoming || matches!(existing, ProposalStatus::Pending)
}

fn is_trusted<Config>(client: &ClientState<Config>, state: &ConsensusState) -> bool {
    state.proposal_status != ProposalStatus::ResolvedInvalid
        && state.finality_level >= client.finality_policy.minimum_membership_level
        && (!client.finality_policy.require_resolved_proposal
            || state.proposal_status == ProposalStatus::ResolvedValid)
}

/// Freezes a client only after a checked same-height conflict.
pub fn freeze_on_conflict<Config>(client: &mut ClientState<Config>, height: u64) {
    client.frozen_height = Some(height);
}
