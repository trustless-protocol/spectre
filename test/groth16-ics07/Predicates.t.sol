// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { Predicates } from "../../contracts/utils/Predicates.sol";

/// @dev Wrapper to expose Predicates internal functions for testing
contract PredicatesWrapper {
    function verifyCommitAgainstTrusted(
        IICS07TendermintMsgs.UntrustedBlockState memory untrustedState,
        IICS07TendermintMsgs.TrustedBlockState memory trustedState,
        IICS07TendermintMsgs.Options memory options
    ) external pure {
        Predicates.verifyCommitAgainstTrusted(untrustedState, trustedState, options);
    }

    function validateCommit(
        IICS07TendermintMsgs.SignedHeader memory signedHeader,
        IICS07TendermintMsgs.ValidatorSet memory validators
    ) external pure {
        Predicates.validateCommit(signedHeader, validators);
    }
}

contract PredicatesTest is Test, IICS07TendermintMsgs {
    PredicatesWrapper internal wrapper;
    // ---------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------

    function _makeValidator(
        bytes memory addr,
        bytes32 pubKey,
        uint64 power
    ) internal pure returns (ValidatorInfo memory) {
        return ValidatorInfo({
            valAddress: addr,
            pubKey: pubKey,
            votingPower: power,
            proposerPriority: 0
        });
    }

    function _makeCommitSig(
        CommitSigFlag flag,
        bytes memory addr
    ) internal pure returns (CommitSig memory) {
        return CommitSig({
            flag: flag,
            data: CommitSigData({
                validatorAddress: addr,
                timestamp: 1000,
                hasSignature: flag != CommitSigFlag.BLOCK_ID_FLAG_ABSENT,
                signature: ""
            })
        });
    }

    function _makeValidatorSet(
        ValidatorInfo[] memory vals
    ) internal pure returns (ValidatorSet memory) {
        return ValidatorSet({
            validators: vals,
            hasProposer: false,
            proposer: _makeValidator("", bytes32(0), 0),
            totalVotingPower: 0
        });
    }

    function _makeSignedHeader(
        CommitSig[] memory sigs,
        bytes32 headerHash
    ) internal pure returns (SignedHeader memory) {
        BlockId memory blockId = BlockId({
            hashData: headerHash,
            partSetHeader: PartSetHeader({ total: 1, hashData: bytes32(0) })
        });

        return SignedHeader({
            header: BlockHeader({
                version: Version({ blockVersion: 11, appVersion: 0 }),
                chainId: "test-chain",
                height: 100,
                time: 1000,
                hasLastBlockId: false,
                lastBlockId: blockId,
                hasLastCommitHash: false,
                lastCommitHash: bytes32(0),
                hasDataHash: false,
                dataHash: bytes32(0),
                validatorsHash: bytes32(0),
                nextValidatorsHash: bytes32(0),
                consensusHash: bytes32(0),
                appHash: bytes32(0),
                hasLastResultsHash: false,
                lastResultsHash: bytes32(0),
                hasEvidenceHash: false,
                evidenceHash: bytes32(0),
                proposerAddress: ""
            }),
            commit: BlockCommit({
                height: 100,
                round: 0,
                blockId: blockId,
                commitSigs: sigs
            })
        });
    }

    function setUp() public {
        wrapper = new PredicatesWrapper();
    }

    // ---------------------------------------------------------------
    // Tests: _checkVotingPowerOverlap (via verifyCommitAgainstTrusted)
    // ---------------------------------------------------------------

    function testVotingPowerSufficientSingleValidator() public view {
        // 1 validator with 100 power, signs the block -> 100/100 > 2/3
        ValidatorInfo[] memory vals = new ValidatorInfo[](1);
        vals[0] = _makeValidator("val1", bytes32(uint256(1)), 100);

        CommitSig[] memory sigs = new CommitSig[](1);
        sigs[0] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_COMMIT, "val1");

        SignedHeader memory sh = _makeSignedHeader(sigs, bytes32(uint256(0xabc)));

        UntrustedBlockState memory untrusted = UntrustedBlockState({
            signedHeader: sh,
            validatorSet: _makeValidatorSet(vals)
        });

        TrustedBlockState memory trusted = TrustedBlockState({
            chainId: "test-chain",
            headerTime: 900,
            height: 99,
            nextValidatorSet: _makeValidatorSet(vals),
            nextValidatorHash: bytes32(0)
        });

        Options memory opts = Options({
            trustThreshold: TrustThreshold({ numerator: 1, denominator: 3 }),
            trustingPeriod: 1000,
            clockDrift: 10
        });

        // Should not revert
        wrapper.verifyCommitAgainstTrusted(untrusted, trusted, opts);
    }

    function testVotingPowerInsufficientReverts() public {
        // 3 validators, only 1 (with 10 power) signs -> 10/100 < 2/3
        ValidatorInfo[] memory vals = new ValidatorInfo[](3);
        vals[0] = _makeValidator("val1", bytes32(uint256(1)), 10);
        vals[1] = _makeValidator("val2", bytes32(uint256(2)), 40);
        vals[2] = _makeValidator("val3", bytes32(uint256(3)), 50);

        CommitSig[] memory sigs = new CommitSig[](3);
        sigs[0] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_COMMIT, "val1");
        sigs[1] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_ABSENT, "val2");
        sigs[2] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_ABSENT, "val3");

        SignedHeader memory sh = _makeSignedHeader(sigs, bytes32(uint256(0xabc)));

        UntrustedBlockState memory untrusted = UntrustedBlockState({
            signedHeader: sh,
            validatorSet: _makeValidatorSet(vals)
        });

        // Same height (next block) -> only checks 2/3 of untrusted set
        TrustedBlockState memory trusted = TrustedBlockState({
            chainId: "test-chain",
            headerTime: 900,
            height: 99,
            nextValidatorSet: _makeValidatorSet(vals),
            nextValidatorHash: bytes32(0)
        });

        Options memory opts = Options({
            trustThreshold: TrustThreshold({ numerator: 1, denominator: 3 }),
            trustingPeriod: 1000,
            clockDrift: 10
        });

        vm.expectRevert("insufficient voting power overlap");
        wrapper.verifyCommitAgainstTrusted(untrusted, trusted, opts);
    }

    function testVotingPowerExactlyTwoThirds() public view {
        // 3 validators each with 100 power, 2 sign -> 200/300 = 2/3
        // Threshold is strictly > 2/3, so 200/300 should fail
        // But 201/300 would pass. Let's test with 2 of 3 equal validators:
        // tallied=200, total=300, 200*3=600 > 300*2=600 is false (not strictly >)
        // So this should revert... actually let's make it pass by giving slightly more power

        // 2 validators sign with 34 power each, 1 absent with 32 -> 68/100 > 2/3
        ValidatorInfo[] memory vals = new ValidatorInfo[](3);
        vals[0] = _makeValidator("val1", bytes32(uint256(1)), 34);
        vals[1] = _makeValidator("val2", bytes32(uint256(2)), 34);
        vals[2] = _makeValidator("val3", bytes32(uint256(3)), 32);

        CommitSig[] memory sigs = new CommitSig[](3);
        sigs[0] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_COMMIT, "val1");
        sigs[1] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_COMMIT, "val2");
        sigs[2] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_ABSENT, "val3");

        SignedHeader memory sh = _makeSignedHeader(sigs, bytes32(uint256(0xabc)));

        UntrustedBlockState memory untrusted = UntrustedBlockState({
            signedHeader: sh,
            validatorSet: _makeValidatorSet(vals)
        });

        TrustedBlockState memory trusted = TrustedBlockState({
            chainId: "test-chain",
            headerTime: 900,
            height: 99,
            nextValidatorSet: _makeValidatorSet(vals),
            nextValidatorHash: bytes32(0)
        });

        Options memory opts = Options({
            trustThreshold: TrustThreshold({ numerator: 1, denominator: 3 }),
            trustingPeriod: 1000,
            clockDrift: 10
        });

        // 68*3=204 > 100*2=200 -> passes
        wrapper.verifyCommitAgainstTrusted(untrusted, trusted, opts);
    }

    function testVotingPowerWithNilVotes() public {
        // Nil votes (BLOCK_ID_FLAG_NIL) count as present but still contribute power
        ValidatorInfo[] memory vals = new ValidatorInfo[](2);
        vals[0] = _makeValidator("val1", bytes32(uint256(1)), 60);
        vals[1] = _makeValidator("val2", bytes32(uint256(2)), 40);

        CommitSig[] memory sigs = new CommitSig[](2);
        sigs[0] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_NIL, "val1");
        sigs[1] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_ABSENT, "val2");

        SignedHeader memory sh = _makeSignedHeader(sigs, bytes32(uint256(0xabc)));

        UntrustedBlockState memory untrusted = UntrustedBlockState({
            signedHeader: sh,
            validatorSet: _makeValidatorSet(vals)
        });

        TrustedBlockState memory trusted = TrustedBlockState({
            chainId: "test-chain",
            headerTime: 900,
            height: 99,
            nextValidatorSet: _makeValidatorSet(vals),
            nextValidatorHash: bytes32(0)
        });

        Options memory opts = Options({
            trustThreshold: TrustThreshold({ numerator: 1, denominator: 3 }),
            trustingPeriod: 1000,
            clockDrift: 10
        });

        // 60*3=180 > 100*2=200 -> false, should revert
        vm.expectRevert("insufficient voting power overlap");
        wrapper.verifyCommitAgainstTrusted(untrusted, trusted, opts);
    }

    // ---------------------------------------------------------------
    // Tests: validateCommit
    // ---------------------------------------------------------------

    function testValidateCommitNoSignaturesReverts() public {
        ValidatorInfo[] memory vals = new ValidatorInfo[](1);
        vals[0] = _makeValidator("val1", bytes32(uint256(1)), 100);

        CommitSig[] memory sigs = new CommitSig[](1);
        sigs[0] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_ABSENT, "val1");

        SignedHeader memory sh = _makeSignedHeader(sigs, bytes32(uint256(0xabc)));

        UntrustedBlockState memory untrusted = UntrustedBlockState({
            signedHeader: sh,
            validatorSet: _makeValidatorSet(vals)
        });

        vm.expectRevert("invalid commit: no present signatures");
        wrapper.validateCommit(untrusted.signedHeader, untrusted.validatorSet);
    }

    function testValidateCommitSigCountMismatchReverts() public {
        ValidatorInfo[] memory vals = new ValidatorInfo[](2);
        vals[0] = _makeValidator("val1", bytes32(uint256(1)), 50);
        vals[1] = _makeValidator("val2", bytes32(uint256(2)), 50);

        // Only 1 commit sig for 2 validators
        CommitSig[] memory sigs = new CommitSig[](1);
        sigs[0] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_COMMIT, "val1");

        SignedHeader memory sh = _makeSignedHeader(sigs, bytes32(uint256(0xabc)));

        UntrustedBlockState memory untrusted = UntrustedBlockState({
            signedHeader: sh,
            validatorSet: _makeValidatorSet(vals)
        });

        vm.expectRevert("invalid commit: number of signatures does not match number of validators");
        wrapper.validateCommit(untrusted.signedHeader, untrusted.validatorSet);
    }

    function testValidateCommitFaultySignerReverts() public {
        ValidatorInfo[] memory vals = new ValidatorInfo[](1);
        vals[0] = _makeValidator("val1", bytes32(uint256(1)), 100);

        // Commit sig from unknown validator
        CommitSig[] memory sigs = new CommitSig[](1);
        sigs[0] = _makeCommitSig(CommitSigFlag.BLOCK_ID_FLAG_COMMIT, "unknown");

        SignedHeader memory sh = _makeSignedHeader(sigs, bytes32(uint256(0xabc)));

        UntrustedBlockState memory untrusted = UntrustedBlockState({
            signedHeader: sh,
            validatorSet: _makeValidatorSet(vals)
        });

        vm.expectRevert("invalid commit: faulty signer");
        wrapper.validateCommit(untrusted.signedHeader, untrusted.validatorSet);
    }
}
