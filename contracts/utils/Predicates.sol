// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import { Header } from "./Header.sol";

/// @title Predicates
library Predicates {
    /// Verify that the header hashes to the committed block ID and that the
    /// commit shape is consistent with the supplied validator set.
    function verifyHeaderMatchesCommit(
        IICS07TendermintMsgs.UntrustedBlockState memory untrustedState
    ) internal pure {
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
            _checkVotingPowerOverlapByAddress(
                untrustedState.signedHeader,
                trustedState.nextValidatorSet,
                options.trustThreshold
            );
            // Also check that untrusted validators have enough signers
            IICS07TendermintMsgs.TrustThreshold memory twoThirds = IICS07TendermintMsgs.TrustThreshold({
                numerator: 2,
                denominator: 3
            });
            _checkVotingPowerOverlapByIndex(
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
            _checkVotingPowerOverlapByIndex(
                untrustedState.signedHeader,
                untrustedState.validatorSet,
                trustThreshold
            );
        }
    }

    /// @notice Check only the trust-threshold overlap against the trusted next
    /// validator set. Adjacent headers do not need this overlap check because
    /// the trusted next-validator hash already pins the next validator set; the
    /// outer Groth16ICS07 contract enforces the >2/3 proof-backed quorum over
    /// the untrusted validator set for all update results, including NoOp.
    function verifyTrustedCommitOverlap(
        IICS07TendermintMsgs.UntrustedBlockState memory untrustedState,
        IICS07TendermintMsgs.TrustedBlockState memory trustedState,
        IICS07TendermintMsgs.Options memory options
    ) internal pure {
        uint64 trustedNextHeight = trustedState.height + 1;
        if (untrustedState.signedHeader.header.height == trustedNextHeight) {
            return;
        }
        _checkVotingPowerOverlapByAddress(
            untrustedState.signedHeader,
            trustedState.nextValidatorSet,
            options.trustThreshold
        );
    }

    function validateCommit(
        IICS07TendermintMsgs.SignedHeader memory signedHeader,
        IICS07TendermintMsgs.ValidatorSet memory validators
    ) internal pure {
        IICS07TendermintMsgs.CommitSig[] memory commitSigs = signedHeader.commit.commitSigs;

        if (commitSigs.length != validators.validators.length) {
            revert("invalid commit: number of signatures does not match number of validators");
        }

        bool hasPresentSignature = false;
        for (uint256 i = 0; i < commitSigs.length; i++) {
            IICS07TendermintMsgs.CommitSig memory sig = commitSigs[i];
            if (sig.flag == IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_ABSENT) {
                continue;
            }
            hasPresentSignature = true;
            if (
                keccak256(abi.encodePacked(validators.validators[i].valAddress))
                    != keccak256(abi.encodePacked(sig.data.validatorAddress))
            ) {
                revert("invalid commit: faulty signer");
            }
        }

        if (!hasPresentSignature) {
            revert("invalid commit: no present signatures");
        }
    }

    /// @notice Check that enough validators from the given set signed the header
    /// to meet the trust threshold. Ed25519 signature verification is delegated
    /// to the ZK proof — this function only checks voting power.
    function _checkVotingPowerOverlapByAddress(
        IICS07TendermintMsgs.SignedHeader memory signedHeader,
        IICS07TendermintMsgs.ValidatorSet memory validatorSet,
        IICS07TendermintMsgs.TrustThreshold memory trustThreshold
    ) internal pure {
        IICS07TendermintMsgs.CommitSig[] memory commitSigs = signedHeader.commit.commitSigs;
        IICS07TendermintMsgs.ValidatorInfo[] memory validators = validatorSet.validators;

        uint64 totalVotingPower = _sumVotingPower(validators);

        // Tally voting power of non-absent signers that match validators.
        // counted[] prevents a duplicate commitSig address from summing a
        // validator's power more than once (double-count attack on trust overlap).
        uint64 talliedPower = 0;
        bool[] memory counted = new bool[](validators.length);
        for (uint256 i = 0; i < commitSigs.length; i++) {
            IICS07TendermintMsgs.CommitSig memory sig = commitSigs[i];
            if (sig.flag == IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_ABSENT) {
                continue;
            }

            bytes32 signerAddressHash = keccak256(sig.data.validatorAddress);
            for (uint256 j = 0; j < validators.length; j++) {
                if (!counted[j] && keccak256(validators[j].valAddress) == signerAddressHash) {
                    talliedPower += validators[j].votingPower;
                    counted[j] = true;
                    break;
                }
            }

            if (_meetsTrustThreshold(talliedPower, totalVotingPower, trustThreshold)) {
                return;
            }
        }

        require(
            _meetsTrustThreshold(talliedPower, totalVotingPower, trustThreshold),
            "insufficient voting power overlap"
        );
    }

    /// @notice Check that enough voting power is present in the untrusted
    /// validator set itself. This path relies on `validateCommit()` having
    /// already established index alignment between `commitSigs[i]` and
    /// `validatorSet.validators[i]`.
    function _checkVotingPowerOverlapByIndex(
        IICS07TendermintMsgs.SignedHeader memory signedHeader,
        IICS07TendermintMsgs.ValidatorSet memory validatorSet,
        IICS07TendermintMsgs.TrustThreshold memory trustThreshold
    ) internal pure {
        IICS07TendermintMsgs.CommitSig[] memory commitSigs = signedHeader.commit.commitSigs;
        IICS07TendermintMsgs.ValidatorInfo[] memory validators = validatorSet.validators;

        require(
            commitSigs.length == validators.length,
            "invalid commit: number of signatures does not match number of validators"
        );

        uint64 totalVotingPower = _sumVotingPower(validators);
        uint64 talliedPower = 0;

        for (uint256 i = 0; i < commitSigs.length; i++) {
            if (commitSigs[i].flag == IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_ABSENT) {
                continue;
            }

            talliedPower += validators[i].votingPower;
            if (_meetsTrustThreshold(talliedPower, totalVotingPower, trustThreshold)) {
                return;
            }
        }

        require(
            _meetsTrustThreshold(talliedPower, totalVotingPower, trustThreshold),
            "insufficient voting power overlap"
        );
    }

    function _sumVotingPower(
        IICS07TendermintMsgs.ValidatorInfo[] memory validators
    ) private pure returns (uint64 totalVotingPower) {
        for (uint256 i = 0; i < validators.length; i++) {
            totalVotingPower += validators[i].votingPower;
        }
    }

    function _meetsTrustThreshold(
        uint64 talliedPower,
        uint64 totalVotingPower,
        IICS07TendermintMsgs.TrustThreshold memory trustThreshold
    ) private pure returns (bool) {
        return
            uint256(talliedPower) * uint256(trustThreshold.denominator)
                > uint256(totalVotingPower) * uint256(trustThreshold.numerator);
    }
}
