// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IMisbehaviour } from "../interfaces/IMisbehaviour.sol";
import { IMisbehaviourMsgs } from "../light-clients/msgs/IMisbehaviourMsgs.sol";
import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../msgs/IICS02ClientMsgs.sol";
import { IGroth16ICS07TendermintErrors } from "../light-clients/errors/IGroth16ICS07TendermintErrors.sol";
import { Predicates } from "../utils/Predicates.sol";
import { Header } from "../utils/Header.sol";
import { HeightCmp } from "../utils/HeightCmp.sol";
import { ChainId } from "../utils/ChainId.sol";
import {Math} from "@openzeppelin-contracts/utils/math/Math.sol";
/**
 * @title Misbehavior
 * @dev Contract to verify misbehavior
 * Converted from Rust zkVM code for Cosmos SDK proof verification
 */
contract Misbehaviour is IMisbehaviour {

    error MismatchedRevisionHeight(
        uint64 expected,
        uint64 actual
    );
    error InvalidHeaderHeight(
        uint64 height
    );
    error ValSetHashMismatch(
        bytes32 expected,
        bytes32 actual
    );

    // Custom errors
    error InvalidClientId();
    error InvalidChainId(string id);
    error ChainIdMismatch();
    error MisbehaviourVerificationFailed();
    error CheckForMisbehaviourFailed();
    error MisbehaviourNotDetected();

    function misbehaviour (
        IICS07TendermintMsgs.ClientState memory clientState,
        IMisbehaviourMsgs.Misbehaviour memory misbehaviour_,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState1,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState2,
        uint128 time
    ) external pure returns (IMisbehaviourMsgs.MisbehaviourOutput memory) {
        // check client state id and misbehaviour chain id
        require(keccak256(abi.encode((clientState.chainId))) == keccak256(abi.encode(misbehaviour_.header1.signedHeader.header.chainId)), ChainIdMismatch());

        validateBasic(misbehaviour_);

        IICS07TendermintMsgs.Options memory options = IICS07TendermintMsgs.Options({
            trustThreshold: clientState.trustLevel,
            trustingPeriod: clientState.trustingPeriod,
            clockDrift: uint32(15)
        });

        // TODO: convert timestamp nanos to secs
        // Reuse cached revisionNumber from clientState.latestHeight (set at
        // client init from the parsed chain ID).
        IICS07TendermintMsgs.ChainId memory chainId = IICS07TendermintMsgs.ChainId({
            id: clientState.chainId,
            revisionNumber: clientState.latestHeight.revisionNumber
        });

        verifyMisbehaviourHeader(
            misbehaviour_.header1,
            chainId,
            options,
            trustedConsensusState1.timestamp,
            trustedConsensusState1.nextValidatorsHash,
            time
        );

        verifyMisbehaviourHeader(
            misbehaviour_.header2,
            chainId,
            options,
            trustedConsensusState2.timestamp,
            trustedConsensusState2.nextValidatorsHash,
            time
        );
        IMisbehaviourMsgs.MisbehaviourOutput memory output = IMisbehaviourMsgs.MisbehaviourOutput({
            trustedHeight1: IICS02ClientMsgs.Height({
                revisionNumber: chainId.revisionNumber,
                revisionHeight: misbehaviour_.header1.trustedHeight.revisionHeight
            }),
            trustedHeight2: IICS02ClientMsgs.Height({
                revisionNumber: chainId.revisionNumber,
                revisionHeight: misbehaviour_.header2.trustedHeight.revisionHeight
            })
        });
        return output;
    }

    function validateBasic(IMisbehaviourMsgs.Misbehaviour memory misbehaviour_) internal pure {
        validateHeaderBasic(misbehaviour_.header1);
        validateHeaderBasic(misbehaviour_.header2);

        if (keccak256(abi.encode(misbehaviour_.header1.signedHeader.header.chainId)) != keccak256(abi.encode(misbehaviour_.header2.signedHeader.header.chainId))) {
            revert IGroth16ICS07TendermintErrors.ChainIdMismatch ({
                expected: misbehaviour_.header1.signedHeader.header.chainId,
                actual: misbehaviour_.header2.signedHeader.header.chainId
            });
        }

        // Both headers proven equal by the keccak check above; parse once.
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

        if (HeightCmp.lt(header1Height, header2Height)) {
            revert IGroth16ICS07TendermintErrors.InsufficientMisbehaviourHeaderHeight ({
                height1: misbehaviour_.header1.signedHeader.header.height,
                height2: misbehaviour_.header2.signedHeader.header.height
            });
        }

        // Tendermint misbehaviour at the same logical height is defined by
        // conflicting signed block IDs, not only by appHash divergence. Use
        // the commit block hash that validators actually signed so conflicts on
        // validatorsHash / nextValidatorsHash / timestamp / lastCommitHash /
        // other header fields are also detected.
        if (misbehaviour_.header1.signedHeader.commit.blockId.hashData == misbehaviour_.header2.signedHeader.commit.blockId.hashData) {
            revert MisbehaviourNotDetected();
        }
    }

    function validateHeaderBasic(
        IICS07TendermintMsgs.Header memory header
    ) internal pure {
        IICS07TendermintMsgs.ChainId memory chainId = ChainId.get(header.signedHeader.header.chainId);
        if (chainId.revisionNumber != header.trustedHeight.revisionNumber) {
            revert MismatchedRevisionHeight(
                chainId.revisionNumber,
                header.trustedHeight.revisionNumber
            );
        }

        IICS02ClientMsgs.Height memory height = IICS02ClientMsgs.Height({
            revisionNumber: chainId.revisionNumber,
            revisionHeight: header.signedHeader.header.height
        });

        // We need to ensure that the trusted height (representing the
        // height of the header already on chain for which this client update is
        // based on) must be smaller than height of the new header that we're
        // installing.
        if (HeightCmp.ge(header.trustedHeight, height)) {
            revert InvalidHeaderHeight(
                height.revisionHeight
            );
        }

        bytes32 valSetHash = Header.hashValSet(header.validatorSet);
        if (valSetHash != header.signedHeader.header.validatorsHash) {
            revert ValSetHashMismatch(
                header.signedHeader.header.validatorsHash,
                valSetHash
            );
        }
    }

    function verifyMisbehaviourHeader (
        IICS07TendermintMsgs.Header memory header,
        IICS07TendermintMsgs.ChainId memory chainId,
        IICS07TendermintMsgs.Options memory options,
        uint128 trustedTime,
        bytes32 trustedNextValidatorHash,
        uint128 currentTimestamp
    ) pure internal {
        // ensure correctness of the trusted next validator set provided by the relayer
        checkTrustedNextValidatorSet(header, trustedNextValidatorHash);

        // ensure trusted consensus state is within trusting period
        {
            uint128 currentTimeInSeconds = nanosToSeconds(currentTimestamp);
            uint128 trustedTimeInSeconds = nanosToSeconds(trustedTime);

            if (currentTimeInSeconds < trustedTimeInSeconds) {
                revert IGroth16ICS07TendermintErrors.InvalidConsensusStateTimestamp({
                    timestamp: trustedTime
                });
            }
            uint128 durationSinceConsensusState = currentTimeInSeconds - trustedTimeInSeconds;
            if (durationSinceConsensusState >= options.trustingPeriod) {
                revert IGroth16ICS07TendermintErrors.InsufficientTrustingPeriod ({
                    durationSinceConsensusState: durationSinceConsensusState,
                    trustingPeriod: options.trustingPeriod
                });
            }
        }

        parseChainId(chainId.id);
        // main header verification, delegated to the tendermint-light-client crate.
        IICS07TendermintMsgs.UntrustedBlockState memory untrustedState = getUntrustedBlockState(header);
        IICS07TendermintMsgs.TrustedBlockState memory trustedState = getTrustedBlockState(
            chainId.id,
            trustedTime,
            trustedNextValidatorHash,
            header
        );

        Predicates.verifyHeaderMatchesCommit(untrustedState);
        Predicates.verifyAgainstTrusted(untrustedState, trustedState, options.trustingPeriod, currentTimestamp);
        Predicates.verifyCommitAgainstTrusted(untrustedState, trustedState, options);
    }

    function checkTrustedNextValidatorSet(
        IICS07TendermintMsgs.Header memory header,
        bytes32 trustedNextValidatorHash
    ) internal pure {
        bytes32 validatorsHash = Header.hashValSet(header.trustedNextValidatorSet);

        if (validatorsHash != trustedNextValidatorHash) {
            revert IGroth16ICS07TendermintErrors.FailedToVerifyHeader({
                description: "trusted next validator set hash does not match hash stored on chain"
            });
        }
    }

    function getUntrustedBlockState(
        IICS07TendermintMsgs.Header memory header
    ) internal pure returns (IICS07TendermintMsgs.UntrustedBlockState memory) {
        return IICS07TendermintMsgs.UntrustedBlockState({
            signedHeader: header.signedHeader,
            validatorSet: header.validatorSet
        });
    }

    function getTrustedBlockState(
        string memory tmChainId,
        uint128 trustedTime,
        bytes32 trustedNextValidatorHash,
        IICS07TendermintMsgs.Header memory header
    ) internal pure returns (IICS07TendermintMsgs.TrustedBlockState memory) {
        return IICS07TendermintMsgs.TrustedBlockState({
            chainId: tmChainId,
            headerTime: trustedTime,
            height: header.trustedHeight.revisionHeight,
            nextValidatorSet: header.trustedNextValidatorSet,
            nextValidatorHash: trustedNextValidatorHash
        });
    }

    function parseChainId(
        string memory chainId
    ) internal pure {
        bytes memory chainIdBytes = bytes(chainId);
        if (chainIdBytes.length == 0 || chainIdBytes.length > 50) {
            revert("Invalid chain id length");
        }

        for (uint256 i = 0; i < chainIdBytes.length; i++) {
            bytes1 b = chainIdBytes[i];
            
            // Check if byte is valid: a-z, A-Z, 0-9, -, _, .
            if (!((b >= 0x61 && b <= 0x7A) ||  // a-z
                  (b >= 0x41 && b <= 0x5A) ||  // A-Z
                  (b >= 0x30 && b <= 0x39) ||  // 0-9
                  (b == 0x2D) ||               // -
                  (b == 0x5F) ||               // _
                  (b == 0x2E))) {              // .
                revert("invalid chain id charset");
            }
        }
    }


    function nanosToSeconds(uint128 timestamp) internal pure returns (uint128) {
        return timestamp / 1e9;
    }
}
