pragma solidity ^0.8.0;

import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";

/// @title VotingPowerCalculator
library VotingPowerCalculator {

    error InsufficientSignersOverlap();
    error DuplicateValidator(bytes valAddress);

    // /// Check if there is 2/3rd overlap between an untrusted header and untrusted validator set
    // function checkSignersOverlap(
    //     IICS07TendermintMsgs.SignedHeader untrustedHeader,
    //     IICS07TendermintMsgs.ValidatorSet untrustedValidators
    // ) {
    //     let trust_threshold = TrustThreshold::TWO_THIRDS;
    //     self.voting_power_in(untrustedHeader, untrustedValidators, trust_threshold)    
    // }

    // function votingPowerIn(
    //     IICS07TendermintMsgs.SignedHeader signedHeader,
    //     IICS07TendermintMsgs.ValidatorSet validatorSet,
    //     IICS07TendermintMsgs.TrustThreshold trustThreshold,
    //     uint64 totalVotingPower,
    // ) internal pure returns (IICS07TendermintMsgs.VotingPowerTally) {
    //     IICS07TendermintMsgs.VotingPowerTally power = IICS07TendermintMsgs.VotingPowerTally({
    //         total: totalVotingPower,
    //         tallied: 0,
    //         trustThreshold: trustThreshold,
    //     });
    //     uint64[] seenVals = [];

    //     IICS07TendermintMsgs.NonAbsentCommitVotes votes = signedHeader
    //         .commit
    //         .signatures
    //         .iter()
    //         .enumerate()
    //         .flat_map(|(idx, signature)| {
    //             // We never have more than 2³¹ signatures so this always
    //             // succeeds.
    //             let idx = ValidatorIndex::try_from(idx).unwrap();
    //             NonAbsentCommitVote::new(
    //                 signature,
    //                 idx,
    //                 &signed_header.commit,
    //                 &signed_header.header.chain_id,
    //             )
    //         })
    //         .collect::<Result<Vec<_>, VerificationError>>()?;
    //     votes.sort_unstable_by_key(NonAbsentCommitVote::validator_id);

    //     // Check if there are duplicate signatures.  If at least one duplicate
    //     // is found, report it as an error.
    //     let duplicate = votes
    //         .windows(2)
    //         .find(|pair| pair[0].validator_id() == pair[1].validator_id());
    //     if let Some(pair) = duplicate {
    //         Err(VerificationError::duplicate_validator(
    //             pair[0].validator_id(),
    //         ))
    //     } else {
    //         Ok(Self {
    //             votes,
    //             sign_bytes: Vec::with_capacity(Self::SIGN_BYTES_INITIAL_CAPACITY),
    //         })
    //     }

    //     for (uint i = 0; i <  validatorSet.validators.length; i++) {
    //         IICS07TendermintMsgs.ValidatorInfo validator = validatorSet.validators[i];
    //         (IICS07TendermintMsgs.NonAbsentCommitVotes[] votes, uint64 idx, bool found) = hasVoted(votes, validator)
    //         if (found) {
    //             // Check if this validator has already voted.
    //             if (seenVals.indexOf(idx) > 0) {
    //                 revert (DuplicateValidator(validator.address));
    //             }
    //             seenVals.push(idx);

    //             power = tallyVote(power, validator.votingPower);

    //             // Break early if sufficient voting power is reached.
    //             if (checkTally(power)) {
    //                 break;
    //             }
    //         }
    //     }

    //     return power
    // }

    // function hasVoted(
    //     IICS07TendermintMsgs.NonAbsentCommitVotes votes,
    //     IICS07TendermintMsgs.ValidatorInfo validator,
    // ) public pure returns (IICS07TendermintMsgs.NonAbsentCommitVotes, uint64, bool) {
    //     uint64 index = 0
    //     bool found = false

    //     for (uint i = 0; i < votes.length; i++) {
    //         NonAbsentCommitVote vote = votes[i];
    //         if (vote.signedVote.validatorAddress == validator.valAddress) {
    //             index = i;

    //             if (!vote.verified) {
    //                 let signBytes = vote.signBytes;

    //                 // TODO add sig verify here
    //                 // validator
    //                 //     .verify_signature::<V>(sign_bytes, vote.signed_vote.signature())
    //                 //     .map_err(|_| {
    //                 //         VerificationError::invalid_signature(
    //                 //             vote.signed_vote.signature().as_bytes().to_vec(),
    //                 //             Box::new(validator.clone()),
    //                 //             sign_bytes.to_vec(),
    //                 //         )
    //                 //     })?;

    //                 vote.verified = true;
    //             }

    //             found = true;
    //         }
    //     }

    //     return (votes, index, found);
    // }


    function tallyVote(
        IICS07TendermintMsgs.VotingPowerTally memory power,
        uint64 votingPower
    ) pure internal returns (IICS07TendermintMsgs.VotingPowerTally memory) {
        power.tallied += votingPower;
        require(power.tallied <= power.total, "tallied should be less than total voting power");
        return power;
    }

    function checkTally(
        IICS07TendermintMsgs.VotingPowerTally memory power
    ) pure internal returns (bool) {
        return power.tallied * power.trustThreshold.denominator > power.total * power.trustThreshold.numerator;
    }
}
