pragma solidity ^0.8.0;

import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import { VotingPowerCalculator } from "./VotingPowerCalculator.sol";
import { Header } from "./Header.sol";

/// @title Predicates
library Predicates {
    function verifyValSets(
        IICS07TendermintMsgs.UntrustedBlockState memory untrustedState
    ) internal pure {
        // Ensure the header validator hashes match the given validators
        bytes32 valSetHash = Header.hashValSet(untrustedState.validatorSet);
        if (valSetHash != untrustedState.signedHeader.header.validatorsHash) {
            revert("invalid block: validator set hash mismatch");
        }

        // Ensure the header matches the commit
        bytes32 headerHash = Header.hashHeader(untrustedState.signedHeader.header);
        if (headerHash != untrustedState.signedHeader.commit.blockId.hashData) {
            revert("invalid block: header hash mismatch");
        }

        // Additional implementation specific validation
        validateCommit(untrustedState.signedHeader, untrustedState.validatorSet);
    }

    function verifyAgainstTrusted(
        IICS07TendermintMsgs.UntrustedBlockState memory untrustedState,
        IICS07TendermintMsgs.TrustedBlockState memory trustedState,
        uint64 trustingPeriod,
        uint128 time
    ) internal pure {
        // Ensure the latest trusted header hasn't expired
        // trustingPeriod is in seconds; timestamps are in nanoseconds
        uint128 trustingPeriodNanos = uint128(trustingPeriod) * 1_000_000_000;
        if (time < trustedState.headerTime || time - trustedState.headerTime > trustingPeriodNanos) {
            revert("invalid block: untrusted state is outside of trusting period");
        }

        // Check that the untrusted block is more recent than the trusted state
        require(untrustedState.signedHeader.header.time > trustedState.headerTime, "invalid block: non monotonic bft time");

        // Check that the chain-id of the untrusted block matches that of the trusted state
        require(keccak256(abi.encodePacked(untrustedState.signedHeader.header.chainId)) == keccak256(abi.encodePacked(trustedState.chainId)), "invalid block: chain-id mismatch");

        uint64 trustedNextHeight = trustedState.height + 1;

        if (untrustedState.signedHeader.header.height == trustedNextHeight) {
            // If the untrusted block is the very next block after the trusted block,
            // check that their (next) validator sets hashes match.
            require(untrustedState.signedHeader.header.validatorsHash == trustedState.nextValidatorHash, "invalid block: next validator set hash mismatch");
        } else {
            // Otherwise, ensure that the untrusted block has a greater height than
            // the trusted block.
            require(untrustedState.signedHeader.header.height > trustedNextHeight, "invalid block: non increasing height");
        }
    }

    /// Verify that a) there is enough overlap between the validator sets of the
    /// trusted and untrusted blocks and b) more than 2/3 of the validators
    /// correctly committed the block.
    function verifyCommitAgainstTrusted(
        IICS07TendermintMsgs.UntrustedBlockState memory untrustedState,
        IICS07TendermintMsgs.TrustedBlockState memory trustedState,
        IICS07TendermintMsgs.Options memory options
    ) internal pure {
        // If the trusted validator set has changed we need to check if there’s
        // overlap between the old trusted set and the new untrested header in
        // addition to checking if the new set correctly signed the header.
        uint64 trustedNextHeight = trustedState.height + 1;
        bool needBoth = untrustedState.signedHeader.header.height != trustedNextHeight;

        if (needBoth) {
            // Check trust overlap between trusted validators and untrusted header
            _checkVotingPowerOverlap(
                untrustedState.signedHeader,
                trustedState.nextValidatorSet,
                options.trustThreshold
            );
            // Also check that untrusted validators have enough signers
            IICS07TendermintMsgs.TrustThreshold memory twoThirds = IICS07TendermintMsgs.TrustThreshold({
                numerator: 2,
                denominator: 3
            });
            _checkVotingPowerOverlap(
                untrustedState.signedHeader,
                untrustedState.validatorSet,
                twoThirds
            );
        } else {
            // Check that there is enough signers overlap between the given, untrusted
            // validator set and the untrusted signed header (>= 2/3).
            IICS07TendermintMsgs.TrustThreshold memory trustThreshold = IICS07TendermintMsgs.TrustThreshold({
                numerator: 2,
                denominator: 3
            });
            _checkVotingPowerOverlap(
                untrustedState.signedHeader,
                untrustedState.validatorSet,
                trustThreshold
            );
        }
    }

    function validateCommit(
        IICS07TendermintMsgs.SignedHeader memory signedHeader,
        IICS07TendermintMsgs.ValidatorSet memory validators
    ) internal pure {
        IICS07TendermintMsgs.CommitSig[] memory commitSigs = signedHeader.commit.commitSigs;

        bool hasPresentSignature = false;
        for (uint256 i = 0; i < commitSigs.length; i++) {
            if (commitSigs[i].flag != IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_ABSENT) {
                hasPresentSignature = true;
            }
        }

        if (!hasPresentSignature) {
            revert("invalid commit: no present signatures");
        }

        if (commitSigs.length != validators.validators.length) {
            revert("invalid commit: number of signatures does not match number of validators");
        }

        for (uint256 i = 0; i < commitSigs.length; i++) {
            IICS07TendermintMsgs.CommitSig memory sig = commitSigs[i];
            bytes memory validatorAddress;

            if (sig.flag == IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_ABSENT) {
                continue;
            } else {
                validatorAddress = sig.data.validatorAddress;
            }

            bool found = false;
            for (uint256 j = 0; j < validators.validators.length; j++) {
                if (keccak256(abi.encodePacked(validators.validators[j].valAddress)) == keccak256(abi.encodePacked(validatorAddress))) {
                    found = true;
                    break;
                }
            }
            if (!found) {
                revert("invalid commit: faulty signer");
            }
        }
    }

    /// @notice Check that enough validators from the given set signed the header
    /// to meet the trust threshold. Ed25519 signature verification is delegated
    /// to the ZK proof — this function only checks voting power.
    function _checkVotingPowerOverlap(
        IICS07TendermintMsgs.SignedHeader memory signedHeader,
        IICS07TendermintMsgs.ValidatorSet memory validatorSet,
        IICS07TendermintMsgs.TrustThreshold memory trustThreshold
    ) internal pure {
        IICS07TendermintMsgs.CommitSig[] memory commitSigs = signedHeader.commit.commitSigs;

        // Calculate total voting power
        uint64 totalVotingPower = 0;
        for (uint256 i = 0; i < validatorSet.validators.length; i++) {
            totalVotingPower += validatorSet.validators[i].votingPower;
        }

        // Tally voting power of non-absent signers that match validators
        uint64 talliedPower = 0;
        for (uint256 i = 0; i < commitSigs.length; i++) {
            if (commitSigs[i].flag == IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_ABSENT) {
                continue;
            }

            bytes memory signerAddress = commitSigs[i].data.validatorAddress;

            for (uint256 j = 0; j < validatorSet.validators.length; j++) {
                if (keccak256(abi.encodePacked(validatorSet.validators[j].valAddress)) == keccak256(abi.encodePacked(signerAddress))) {
                    talliedPower += validatorSet.validators[j].votingPower;
                    break;
                }
            }

            // Early exit if threshold already met
            if (talliedPower * trustThreshold.denominator > totalVotingPower * trustThreshold.numerator) {
                return;
            }
        }

        require(
            talliedPower * trustThreshold.denominator > totalVotingPower * trustThreshold.numerator,
            "insufficient voting power overlap"
        );
    }
}