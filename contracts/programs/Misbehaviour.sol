// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import { IMisbehaviour } from "../interfaces/IMisbehaviour.sol";
import { IMisbehaviourMsgs } from "../light-clients/msgs/IMisbehaviourMsgs.sol";
import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../msgs/IICS02ClientMsgs.sol";
import { ISP1ICS07TendermintErrors } from "../light-clients/errors/ISP1ICS07TendermintErrors.sol";
import { Predicates } from "../utils/Predicates.sol";
import { Header } from "../utils/Header.sol";
import { HeightCmp } from "../utils/HeightCmp.sol";
import { Strings } from "@openzeppelin-contracts/utils/Strings.sol";
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
        IICS07TendermintMsgs.ChainId memory chainId = getChainId(clientState.chainId);

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
                revisionHeight: misbehaviour_.header1.signedHeader.header.height
            }),
            trustedHeight2: IICS02ClientMsgs.Height({
                revisionNumber: chainId.revisionNumber,
                revisionHeight: misbehaviour_.header2.signedHeader.header.height
            })
        });
        return output;
    }

    function validateBasic(IMisbehaviourMsgs.Misbehaviour memory misbehaviour_) internal pure {
        validateHeaderBasic(misbehaviour_.header1);
        validateHeaderBasic(misbehaviour_.header2);

        if (keccak256(abi.encode(misbehaviour_.header1.signedHeader.header.chainId)) != keccak256(abi.encode(misbehaviour_.header2.signedHeader.header.chainId))) {
            revert ISP1ICS07TendermintErrors.ChainIdMismatch ({
                expected: misbehaviour_.header1.signedHeader.header.chainId,
                actual: misbehaviour_.header2.signedHeader.header.chainId
            });
        }

        IICS07TendermintMsgs.ChainId memory header1ChainId = getChainId(misbehaviour_.header1.signedHeader.header.chainId);
        IICS07TendermintMsgs.ChainId memory header2ChainId = getChainId(misbehaviour_.header2.signedHeader.header.chainId);

        IICS02ClientMsgs.Height memory header1Height = IICS02ClientMsgs.Height({
            revisionNumber: header1ChainId.revisionNumber,
            revisionHeight: misbehaviour_.header1.signedHeader.header.height
        });

        IICS02ClientMsgs.Height memory header2Height = IICS02ClientMsgs.Height({
            revisionNumber: header2ChainId.revisionNumber,
            revisionHeight: misbehaviour_.header2.signedHeader.header.height
        });

        if (HeightCmp.lt(header1Height, header2Height)) {
            revert ISP1ICS07TendermintErrors.InsufficientMisbehaviourHeaderHeight ({
                height1: misbehaviour_.header1.signedHeader.header.height,
                height2: misbehaviour_.header2.signedHeader.header.height
            });
        }
    }

    function validateHeaderBasic(
        IICS07TendermintMsgs.Header memory header
    ) internal pure {
        IICS07TendermintMsgs.ChainId memory chainId = getChainId(header.signedHeader.header.chainId);
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
            if (currentTimestamp < trustedTime) {
                revert ISP1ICS07TendermintErrors.InvalidConsensusStateTimestamp({
                    timestamp: trustedTime
                });
            }
            uint128 durationSinceConsensusState =currentTimestamp - trustedTime;
            if (durationSinceConsensusState >= options.trustingPeriod) {
                revert ISP1ICS07TendermintErrors.InsufficientTrustingPeriod ({
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

        Predicates.verifyValSets(untrustedState);
        Predicates.verifyAgainstTrusted(untrustedState, trustedState, options.trustingPeriod, currentTimestamp);
        Predicates.verifyCommitAgainstTrusted(untrustedState, trustedState, options);
    }

    function checkTrustedNextValidatorSet(
        IICS07TendermintMsgs.Header memory header,
        bytes32 trustedNextValidatorHash
    ) internal pure {
        bytes32 validatorsHash = Header.hashValSet(header.trustedNextValidatorSet);

        if (validatorsHash != trustedNextValidatorHash) {
            revert ISP1ICS07TendermintErrors.FailedToVerifyHeader({
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

    function getChainId(string memory id) pure internal returns (IICS07TendermintMsgs.ChainId memory) {
        // TODO: validate id
        bytes memory chainIdBytes = bytes(id);

        // Find the last occurrence of '-'
        int256 lastDashIndex = -1;
        for (uint256 i = chainIdBytes.length-1; i >= 0; i--) {
            if (chainIdBytes[i] == 0x2D) { // '-' character
                lastDashIndex = int256(i);
                break;
            }
        }

        if (lastDashIndex == -1) {
            if (1 <= chainIdBytes.length && chainIdBytes.length < 64) {
                return IICS07TendermintMsgs.ChainId({
                    id: id, 
                    revisionNumber: 0
                });
            } else {
                revert("Invalid chain ID length");
            }
        }

        uint256 dashIndex = uint256(lastDashIndex);
        // Extract revision number string
        bytes memory revisionBytes = new bytes(chainIdBytes.length - dashIndex - 1);
        for (uint256 i = 0; i < revisionBytes.length; i++) {
            revisionBytes[i] = chainIdBytes[dashIndex + 1 + i];
        }

        // Validates the revision number not to start with leading zeros, like "01".
        // Zero is the only allowed revision number with leading zero.
        if (revisionBytes.length == 0 || (revisionBytes[0] == 0x30 && revisionBytes.length > 1)) {
            if (1 <= chainIdBytes.length && chainIdBytes.length < 64) {
                return IICS07TendermintMsgs.ChainId({
                    id: id, 
                    revisionNumber: 0
                });
            } else {
                revert("Invalid chain ID length");
            }
        }

        (bool success, uint256 parsedRevisionNumber) = Strings.tryParseUint(string(revisionBytes));
        if (! success || parsedRevisionNumber > type(uint64).max) {
            if (1 <= chainIdBytes.length && chainIdBytes.length < 64) {
                return IICS07TendermintMsgs.ChainId({
                    id: id, 
                    revisionNumber: 0
                });
            } else {
                revert("Invalid chain ID length");
            }
        }

        uint64 revisionNumber = uint64(parsedRevisionNumber);
        
        // Extract chain name
        bytes memory chainNameBytes = new bytes(dashIndex);
        for (uint256 i = 0; i < dashIndex; i++) {
            chainNameBytes[i] = chainIdBytes[i];
        }

        uint256 min = 1;
        // Prefix must be at most `max_id_length - 21` characters long since the
        // longest identifier we can construct is `{prefix}-{u64::MAX}` which
        // extends prefix by 21 characters.
        uint256 max = 43;
        if (chainNameBytes.length <= min || chainNameBytes.length >= max) {
            revert("Invalid chain prefix length");
        }

        return IICS07TendermintMsgs.ChainId({
            id: id, 
            revisionNumber: revisionNumber
        });
    }
}