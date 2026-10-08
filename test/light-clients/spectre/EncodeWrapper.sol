// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// solhint-disable custom-errors

import { Encode } from "contracts/light-clients/spectre/libraries/Encode.sol";
import { Header } from "contracts/light-clients/spectre/libraries/Header.sol";
import { ISpectreTendermintMsgs } from "contracts/light-clients/spectre/messages/ISpectreTendermintMsgs.sol";

/// @title EncodeWrapper
/// @notice Wrapper with primitive parameters so the Encode library functions can be called through
///         `cast call` and compared with Go proto.Marshal() output. Nothing in this repository
///         calls it today: the Go side of that comparison is tracked as SUP-06 in #288.
contract EncodeWrapper {
    function encodeVersion(uint64 blockVersion, uint64 appVersion) external pure returns (bytes memory) {
        return Encode.encodeVersion(ISpectreTendermintMsgs.Version(blockVersion, appVersion));
    }

    function encodeValidator(bytes32 pubKey, uint64 votingPower) external pure returns (bytes memory) {
        return Encode.encodeValidator(ISpectreTendermintMsgs.SimpleValidator(pubKey, votingPower));
    }

    function encodePartSetHeader(uint32 total, bytes32 hashData) external pure returns (bytes memory) {
        return Encode.encodePartSetHeader(ISpectreTendermintMsgs.PartSetHeader(total, hashData));
    }

    function encodeBlockId(bytes32 hashData, uint32 pshTotal, bytes32 pshHash) external pure returns (bytes memory) {
        return Encode.encodeBlockId(
            ISpectreTendermintMsgs.BlockId(hashData, ISpectreTendermintMsgs.PartSetHeader(pshTotal, pshHash))
        );
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
        ISpectreTendermintMsgs.CommitSig[] memory sigs = new ISpectreTendermintMsgs.CommitSig[](1);
        // Canonical-vote encoding does not cover validatorAddress.
        sigs[0] = ISpectreTendermintMsgs.CommitSig({
            flag: ISpectreTendermintMsgs.CommitSigFlag(flag), validatorAddress: bytes20(0)
        });

        ISpectreTendermintMsgs.BlockCommit memory commit = ISpectreTendermintMsgs.BlockCommit({
            height: height,
            round: round,
            blockId: ISpectreTendermintMsgs.BlockId(
                blockIdHash, ISpectreTendermintMsgs.PartSetHeader(pshTotal, pshHash)
            ),
            commitSigs: sigs
        });

        return Encode.voteSignBytes(commit, chainId, 0, timestamp);
    }

    function hashValSet(bytes32[] calldata pubKeys, uint64[] calldata votingPowers) external pure returns (bytes32) {
        require(pubKeys.length == votingPowers.length, "length mismatch");

        ISpectreTendermintMsgs.ValidatorInfo[] memory vals = new ISpectreTendermintMsgs.ValidatorInfo[](pubKeys.length);
        for (uint256 i = 0; i < pubKeys.length; i++) {
            vals[i] = ISpectreTendermintMsgs.ValidatorInfo({
                valAddress: hex"", pubKey: pubKeys[i], votingPower: votingPowers[i], proposerPriority: 0
            });
        }

        ISpectreTendermintMsgs.ValidatorSet memory valSet = ISpectreTendermintMsgs.ValidatorSet({
            validators: vals,
            hasProposer: false,
            proposer: ISpectreTendermintMsgs.ValidatorInfo({
                valAddress: hex"", pubKey: bytes32(0), votingPower: 0, proposerPriority: 0
            }),
            totalVotingPower: 0
        });

        return Header.hashValSet(valSet);
    }
}
