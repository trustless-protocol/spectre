// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// solhint-disable custom-errors

import { Encode } from "contracts/light-clients/spectre/libraries/Encode.sol";
import { Header } from "contracts/light-clients/spectre/libraries/Header.sol";
import { SpectreMsgs } from "contracts/light-clients/spectre/messages/SpectreMsgs.sol";

/// @title EncodeWrapper
/// @notice Wrapper with primitive parameters for cross-validation via cast call.
///         Used by Go test (operator/cmd/encode_debug/cross_test.go) to call
///         Solidity encoding functions and compare output with Go proto.Marshal().
contract EncodeWrapper {
    function encodeVersion(uint64 blockVersion, uint64 appVersion) external pure returns (bytes memory) {
        return Encode.encodeVersion(SpectreMsgs.Version(blockVersion, appVersion));
    }

    function encodeValidator(bytes32 pubKey, uint64 votingPower) external pure returns (bytes memory) {
        return Encode.encodeValidator(SpectreMsgs.SimpleValidator(pubKey, votingPower));
    }

    function encodePartSetHeader(uint32 total, bytes32 hashData) external pure returns (bytes memory) {
        return Encode.encodePartSetHeader(SpectreMsgs.PartSetHeader(total, hashData));
    }

    function encodeBlockId(bytes32 hashData, uint32 pshTotal, bytes32 pshHash) external pure returns (bytes memory) {
        return Encode.encodeBlockId(SpectreMsgs.BlockId(hashData, SpectreMsgs.PartSetHeader(pshTotal, pshHash)));
    }

    function cdcEncodeString(string calldata value) external pure returns (bytes memory) {
        return Encode.cdcEncodeString(value);
    }

    function cdcEncodeInt64(uint256 value) external pure returns (bytes memory) {
        return Encode.cdcEncodeInt64(value);
    }

    function cdcEncodeBytes32(bytes32 value) external pure returns (bytes memory) {
        return Encode.cdcEncodeBytes32(value);
    }

    function encodeTimestamp(uint128 nanos) external pure returns (bytes memory) {
        return Encode.encodeTimestamp(nanos);
    }

    function voteSignBytes(
        uint64 height,
        uint32 round,
        bytes32 blockIdHash,
        uint32 pshTotal,
        bytes32 pshHash,
        uint8 flag,
        uint128 timestamp,
        string calldata chainId
    )
        external
        pure
        returns (bytes memory)
    {
        SpectreMsgs.CommitSig[] memory sigs = new SpectreMsgs.CommitSig[](1);
        // Canonical-vote encoding does not cover validatorAddress.
        sigs[0] = SpectreMsgs.CommitSig({ flag: SpectreMsgs.CommitSigFlag(flag), validatorAddress: bytes20(0) });

        SpectreMsgs.BlockCommit memory commit = SpectreMsgs.BlockCommit({
            height: height,
            round: round,
            blockId: SpectreMsgs.BlockId(blockIdHash, SpectreMsgs.PartSetHeader(pshTotal, pshHash)),
            commitSigs: sigs
        });

        return Encode.voteSignBytes(commit, chainId, 0, timestamp);
    }

    function hashValSet(bytes32[] calldata pubKeys, uint64[] calldata votingPowers) external pure returns (bytes32) {
        require(pubKeys.length == votingPowers.length, "length mismatch");

        SpectreMsgs.ValidatorInfo[] memory vals = new SpectreMsgs.ValidatorInfo[](pubKeys.length);
        for (uint256 i = 0; i < pubKeys.length; i++) {
            vals[i] = SpectreMsgs.ValidatorInfo({
                valAddress: hex"", pubKey: pubKeys[i], votingPower: votingPowers[i], proposerPriority: 0
            });
        }

        SpectreMsgs.ValidatorSet memory valSet = SpectreMsgs.ValidatorSet({
            validators: vals,
            hasProposer: false,
            proposer: SpectreMsgs.ValidatorInfo({
                valAddress: hex"", pubKey: bytes32(0), votingPower: 0, proposerPriority: 0
            }),
            totalVotingPower: 0
        });

        return Header.hashValSet(valSet);
    }
}
