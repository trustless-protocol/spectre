// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IMisbehaviour } from "../interfaces/IMisbehaviour.sol";
import { IMisbehaviourMsgs } from "../light-clients/msgs/IMisbehaviourMsgs.sol";
import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../msgs/IICS02ClientMsgs.sol";
import { IGroth16ICS07TendermintErrors } from "../light-clients/errors/IGroth16ICS07TendermintErrors.sol";
import { Header } from "../utils/Header.sol";
import { HeightCmp } from "../utils/HeightCmp.sol";
import { ChainId } from "../utils/ChainId.sol";

contract Misbehaviour is IMisbehaviour {
    error MismatchedRevisionHeight(uint64 expected, uint64 actual);
    error InvalidHeaderHeight(uint64 height);
    error ChainIdMismatch();
    error MisbehaviourNotDetected();

    function misbehaviour(
        IICS07TendermintMsgs.ClientState memory clientState,
        IMisbehaviourMsgs.Misbehaviour memory misbehaviour_,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState1,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState2,
        uint128 time
    )
        external
        pure
        returns (IMisbehaviourMsgs.MisbehaviourOutput memory)
    {
        require(
            keccak256(bytes(clientState.chainId))
                == keccak256(bytes(misbehaviour_.header1.signedHeader.header.chainId)),
            ChainIdMismatch()
        );

        validateBasic(misbehaviour_);

        IICS07TendermintMsgs.Options memory options = IICS07TendermintMsgs.Options({
            trustThreshold: clientState.trustLevel,
            trustingPeriod: clientState.trustingPeriod,
            clockDrift: clientState.clockDrift
        });
        IICS07TendermintMsgs.ChainId memory chainId = IICS07TendermintMsgs.ChainId({
            id: clientState.chainId, revisionNumber: clientState.latestHeight.revisionNumber
        });

        verifyMisbehaviourHeader(misbehaviour_.header1, chainId, options, trustedConsensusState1.timestamp, time);
        verifyMisbehaviourHeader(misbehaviour_.header2, chainId, options, trustedConsensusState2.timestamp, time);

        return IMisbehaviourMsgs.MisbehaviourOutput({
            trustedHeight1: IICS02ClientMsgs.Height({
                revisionNumber: chainId.revisionNumber,
                revisionHeight: misbehaviour_.header1.trustedHeight.revisionHeight
            }),
            trustedHeight2: IICS02ClientMsgs.Height({
                revisionNumber: chainId.revisionNumber,
                revisionHeight: misbehaviour_.header2.trustedHeight.revisionHeight
            })
        });
    }

    function validateBasic(IMisbehaviourMsgs.Misbehaviour memory misbehaviour_) internal pure {
        validateHeaderBasic(misbehaviour_.header1);
        validateHeaderBasic(misbehaviour_.header2);

        if (
            keccak256(bytes(misbehaviour_.header1.signedHeader.header.chainId))
                != keccak256(bytes(misbehaviour_.header2.signedHeader.header.chainId))
        ) {
            revert IGroth16ICS07TendermintErrors.ChainIdMismatch({
                expected: misbehaviour_.header1.signedHeader.header.chainId,
                actual: misbehaviour_.header2.signedHeader.header.chainId
            });
        }

        IICS07TendermintMsgs.ChainId memory headerChainId =
            ChainId.get(misbehaviour_.header1.signedHeader.header.chainId);

        IICS02ClientMsgs.Height memory header1Height = IICS02ClientMsgs.Height({
            revisionNumber: headerChainId.revisionNumber,
            revisionHeight: misbehaviour_.header1.signedHeader.header.height
        });
        IICS02ClientMsgs.Height memory header2Height = IICS02ClientMsgs.Height({
            revisionNumber: headerChainId.revisionNumber,
            revisionHeight: misbehaviour_.header2.signedHeader.header.height
        });

        if (header1Height.revisionHeight != header2Height.revisionHeight) {
            revert IGroth16ICS07TendermintErrors.MismatchedMisbehaviourHeaderHeights({
                height1: misbehaviour_.header1.signedHeader.header.height,
                height2: misbehaviour_.header2.signedHeader.header.height
            });
        }

        if (
            misbehaviour_.header1.signedHeader.commit.blockId.hashData
                == misbehaviour_.header2.signedHeader.commit.blockId.hashData
        ) {
            revert MisbehaviourNotDetected();
        }
    }

    function validateHeaderBasic(IICS07TendermintMsgs.Header memory header) internal pure {
        IICS07TendermintMsgs.ChainId memory chainId = ChainId.get(header.signedHeader.header.chainId);
        if (chainId.revisionNumber != header.trustedHeight.revisionNumber) {
            revert MismatchedRevisionHeight(chainId.revisionNumber, header.trustedHeight.revisionNumber);
        }

        IICS02ClientMsgs.Height memory height = IICS02ClientMsgs.Height({
            revisionNumber: chainId.revisionNumber, revisionHeight: header.signedHeader.header.height
        });
        if (HeightCmp.ge(header.trustedHeight, height)) {
            revert InvalidHeaderHeight(height.revisionHeight);
        }

        if (Header.hashHeader(header.signedHeader.header) != header.signedHeader.commit.blockId.hashData) {
            revert IGroth16ICS07TendermintErrors.FailedToVerifyHeader({
                description: "invalid block: header hash mismatch"
            });
        }
    }

    function verifyMisbehaviourHeader(
        IICS07TendermintMsgs.Header memory header,
        IICS07TendermintMsgs.ChainId memory chainId,
        IICS07TendermintMsgs.Options memory options,
        uint128 trustedTime,
        uint128 currentTimestamp
    )
        internal
        pure
    {
        uint128 currentTimeInSeconds = nanosToSeconds(currentTimestamp);
        uint128 trustedTimeInSeconds = nanosToSeconds(trustedTime);

        if (currentTimeInSeconds < trustedTimeInSeconds) {
            revert IGroth16ICS07TendermintErrors.InvalidConsensusStateTimestamp({ timestamp: trustedTime });
        }
        uint128 durationSinceConsensusState = currentTimeInSeconds - trustedTimeInSeconds;
        if (durationSinceConsensusState >= options.trustingPeriod) {
            revert IGroth16ICS07TendermintErrors.InsufficientTrustingPeriod({
                durationSinceConsensusState: durationSinceConsensusState, trustingPeriod: options.trustingPeriod
            });
        }

        parseChainId(chainId.id);
        require(
            keccak256(bytes(header.signedHeader.header.chainId)) == keccak256(bytes(chainId.id)),
            ChainIdMismatch()
        );
        require(header.signedHeader.header.time > trustedTime, "invalid block: non monotonic bft time");
        uint128 drifted = currentTimestamp + uint128(options.clockDrift) * 1_000_000_000;
        require(header.signedHeader.header.time < drifted, "invalid block: header is from the future");
    }

    function parseChainId(string memory chainId) internal pure {
        bytes memory chainIdBytes = bytes(chainId);
        if (chainIdBytes.length == 0 || chainIdBytes.length > 50) {
            revert("Invalid chain id length");
        }

        for (uint256 i = 0; i < chainIdBytes.length; i++) {
            bytes1 b = chainIdBytes[i];
            if (!((b >= 0x61 && b <= 0x7A) || (b >= 0x41 && b <= 0x5A) || (b >= 0x30 && b <= 0x39)
                || (b == 0x2D) || (b == 0x5F) || (b == 0x2E))) {
                revert("invalid chain id charset");
            }
        }
    }

    function nanosToSeconds(uint128 timestamp) internal pure returns (uint128) {
        return timestamp / 1e9;
    }
}
