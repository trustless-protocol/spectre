// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IUpdateClientMsgs } from "../light-clients/msgs/IUpdateClientMsgs.sol";
import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../msgs/IICS02ClientMsgs.sol";
import { IUpdateClient } from "../interfaces/IUpdateClient.sol";
import { Header } from "../utils/Header.sol";
import { HeightCmp } from "../utils/HeightCmp.sol";
import { ChainId } from "../utils/ChainId.sol";

contract UpdateClient is IUpdateClient {
    error MismatchedRevisionHeight(uint64 expected, uint64 actual);
    error InvalidHeaderHeight(uint64 height);
    error HeaderChainIdMismatch(string expected, string actual);
    error FailedToVerifyHeader(string reason);

    function updateClient(IUpdateClientMsgs.MsgUpdateClient calldata msg_)
        external
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        IICS07TendermintMsgs.ChainId memory chainId = IICS07TendermintMsgs.ChainId({
            id: msg_.clientState.chainId, revisionNumber: msg_.clientState.latestHeight.revisionNumber
        });
        IICS07TendermintMsgs.Options memory options = IICS07TendermintMsgs.Options({
            trustThreshold: msg_.clientState.trustLevel,
            trustingPeriod: msg_.clientState.trustingPeriod,
            clockDrift: msg_.clientState.clockDrift
        });

        verifyHeader(msg_.proposedHeader, chainId, options, msg_.time, msg_.trustedConsensusState);
        return _buildOutput(msg_, chainId.revisionNumber);
    }

    function verifyHeader(
        IICS07TendermintMsgs.Header memory proposedHeader,
        IICS07TendermintMsgs.ChainId memory chainId,
        IICS07TendermintMsgs.Options memory options,
        uint128 time,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState
    )
        internal
        pure
    {
        IICS07TendermintMsgs.ChainId memory headerChainId = ChainId.get(proposedHeader.signedHeader.header.chainId);
        validateBasic(proposedHeader, headerChainId);
        verifyChainIdVersion(chainId, headerChainId);

        if (Header.hashHeader(proposedHeader.signedHeader.header) != proposedHeader.signedHeader.commit.blockId.hashData) {
            revert FailedToVerifyHeader("invalid block: header hash mismatch");
        }

        _verifyAgainstTrusted(proposedHeader, chainId.id, options, time, trustedConsensusState);
    }

    function _verifyAgainstTrusted(
        IICS07TendermintMsgs.Header memory proposedHeader,
        string memory chainId,
        IICS07TendermintMsgs.Options memory options,
        uint128 time,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState
    )
        private
        pure
    {
        uint128 trustingPeriodNanos = uint128(options.trustingPeriod) * 1_000_000_000;
        if (time < trustedConsensusState.timestamp || time - trustedConsensusState.timestamp > trustingPeriodNanos) {
            revert FailedToVerifyHeader("invalid block: untrusted state is outside of trusting period");
        }
        require(proposedHeader.signedHeader.header.time > trustedConsensusState.timestamp, "invalid block: non monotonic bft time");
        require(
            keccak256(bytes(proposedHeader.signedHeader.header.chainId)) == keccak256(bytes(chainId)),
            "invalid block: chain-id mismatch"
        );
        uint128 drifted = time + uint128(options.clockDrift) * 1_000_000_000;
        require(proposedHeader.signedHeader.header.time < drifted, "invalid block: header is from the future");

        uint64 trustedNextHeight = proposedHeader.trustedHeight.revisionHeight + 1;
        if (proposedHeader.signedHeader.header.height <= trustedNextHeight) {
            require(proposedHeader.signedHeader.header.height == trustedNextHeight, "invalid block: non increasing height");
        }
    }

    function verifyChainIdVersion(
        IICS07TendermintMsgs.ChainId memory chainId,
        IICS07TendermintMsgs.ChainId memory headerChainId
    )
        internal
        pure
    {
        if (chainId.revisionNumber != headerChainId.revisionNumber) {
            revert HeaderChainIdMismatch(headerChainId.id, chainId.id);
        }
    }

    function validateBasic(
        IICS07TendermintMsgs.Header memory header,
        IICS07TendermintMsgs.ChainId memory headerChainId
    )
        internal
        pure
    {
        if (headerChainId.revisionNumber != header.trustedHeight.revisionNumber) {
            revert MismatchedRevisionHeight(headerChainId.revisionNumber, header.trustedHeight.revisionNumber);
        }

        IICS02ClientMsgs.Height memory height = IICS02ClientMsgs.Height({
            revisionNumber: headerChainId.revisionNumber, revisionHeight: header.signedHeader.header.height
        });

        if (HeightCmp.ge(header.trustedHeight, height)) {
            revert InvalidHeaderHeight(height.revisionHeight);
        }
    }

    function _buildOutput(
        IUpdateClientMsgs.MsgUpdateClient calldata msg_,
        uint64 revisionNumber
    )
        private
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        IICS02ClientMsgs.Height memory trustedHeight = IICS02ClientMsgs.Height({
            revisionNumber: revisionNumber, revisionHeight: msg_.proposedHeader.trustedHeight.revisionHeight
        });
        IICS02ClientMsgs.Height memory newHeight = IICS02ClientMsgs.Height({
            revisionNumber: revisionNumber, revisionHeight: msg_.proposedHeader.signedHeader.header.height
        });
        IICS07TendermintMsgs.ConsensusState memory newConsensusState = IICS07TendermintMsgs.ConsensusState({
            timestamp: msg_.proposedHeader.signedHeader.header.time,
            root: msg_.proposedHeader.signedHeader.header.appHash,
            nextValidatorsHash: msg_.proposedHeader.signedHeader.header.nextValidatorsHash
        });
        return IUpdateClientMsgs.UpdateClientOutput({
            clientState: msg_.clientState,
            trustedConsensusState: msg_.trustedConsensusState,
            newConsensusState: newConsensusState,
            time: msg_.time,
            trustedHeight: trustedHeight,
            newHeight: newHeight
        });
    }
}
