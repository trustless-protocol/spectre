// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import { IUpdateClientMsgs } from "../light-clients/msgs/IUpdateClientMsgs.sol";
import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../msgs/IICS02ClientMsgs.sol";
import { IUpdateClient } from "../interfaces/IUpdateClient.sol";
import { Predicates } from "../utils/Predicates.sol";
import { Header } from "../utils/Header.sol";
import { HeightCmp } from "../utils/HeightCmp.sol";

import { Strings } from "@openzeppelin-contracts/utils/Strings.sol";

contract UpdateClient is IUpdateClient {
    string clientId = "07-tendermint-0";

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
    error HeaderChainIdMismatch(
        string expected,
        string actual
    );
    error FailedToVerifyHeader(
        string reason
    );

    function updateClient(
        IUpdateClientMsgs.MsgUpdateClient calldata msg_
    ) view external returns (IUpdateClientMsgs.UpdateClientOutput memory) {
        IICS07TendermintMsgs.ChainId memory chainId = getChainId(msg_.clientState.chainId);
        IICS07TendermintMsgs.Options memory options = IICS07TendermintMsgs.Options({
            trustThreshold: msg_.clientState.trustLevel,
            trustingPeriod: msg_.clientState.trustingPeriod,
            clockDrift: 15
        });

        IICS07TendermintMsgs.ClientConsensusStatePath memory path = IICS07TendermintMsgs.ClientConsensusStatePath({
            clientId: clientId,
            revisionNumber: msg_.proposedHeader.trustedHeight.revisionNumber,
            revisionHeight: msg_.proposedHeader.trustedHeight.revisionHeight
        });

        verifyHeader(
            msg_.proposedHeader,
            clientId,
            chainId,
            options,
            msg_.time,
            path,
            msg_.trustedConsensusState
        );

        IICS07TendermintMsgs.ConsensusState memory newConsensusState = IICS07TendermintMsgs.ConsensusState({
            timestamp: msg_.proposedHeader.signedHeader.header.time,
            root: msg_.proposedHeader.signedHeader.header.appHash,
            nextValidatorsHash: msg_.proposedHeader.signedHeader.header.nextValidatorsHash
        });

        IICS02ClientMsgs.Height memory newHeight = IICS02ClientMsgs.Height({
            revisionNumber: chainId.revisionNumber,
            revisionHeight: msg_.proposedHeader.signedHeader.header.height
        });
        IUpdateClientMsgs.UpdateClientOutput memory output = IUpdateClientMsgs.UpdateClientOutput({
            clientState: msg_.clientState,
            trustedConsensusState: msg_.trustedConsensusState,
            newConsensusState: newConsensusState,
            time: msg_.time,
            trustedHeight: msg_.proposedHeader.trustedHeight,
            newHeight: newHeight
        });
        return output;
    }

    function verifyHeader(
        IICS07TendermintMsgs.Header memory proposedHeader,
        string memory clientId_,
        IICS07TendermintMsgs.ChainId memory chainId,
        IICS07TendermintMsgs.Options memory options,
        uint128 time,
        IICS07TendermintMsgs.ClientConsensusStatePath memory path_,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState
    ) internal pure {
        // Checks that the header fields are valid.
        validateBasic(proposedHeader);

        // The tendermint-light-client crate though works on heights that are assumed
        // to have the same revision number. We ensure this here.
        verifyChainIdVersion(chainId, proposedHeader);

        // Delegate to tendermint-light-client, which contains the required checks
        // of the new header against the trusted consensus state.
        {
            bytes32 nextValSetHash = Header.hashValSet(proposedHeader.trustedNextValidatorSet);
            if (nextValSetHash != trustedConsensusState.nextValidatorsHash) {
                revert FailedToVerifyHeader(
                    "trusted next validator set hash does not match hash stored on chain"
                );
            }

            IICS07TendermintMsgs.TrustedBlockState memory trustedState = IICS07TendermintMsgs.TrustedBlockState({
                chainId: chainId.id,
                headerTime: trustedConsensusState.timestamp,
                height: proposedHeader.trustedHeight.revisionHeight,
                nextValidatorSet: proposedHeader.trustedNextValidatorSet,
                nextValidatorHash: nextValSetHash
            });

            IICS07TendermintMsgs.UntrustedBlockState memory untrustedState = IICS07TendermintMsgs.UntrustedBlockState({
                signedHeader: proposedHeader.signedHeader,
                validatorSet: proposedHeader.validatorSet
            });

            verifyUpdateHeader(
                untrustedState,
                trustedState,
                options,
                time
            );
        }


    }

    function verifyUpdateHeader(
        IICS07TendermintMsgs.UntrustedBlockState memory untrustedState,
        IICS07TendermintMsgs.TrustedBlockState memory trustedState,
        IICS07TendermintMsgs.Options memory options,
        uint128 time
    ) internal pure {
        Predicates.verifyValSets(untrustedState);
        Predicates.verifyAgainstTrusted(untrustedState, trustedState, options.trustingPeriod, time);
        /// Check that the untrusted header is from past.
        uint128 drifted = time + uint128(options.clockDrift) * 1_000_000_000;
        require(untrustedState.signedHeader.header.time < drifted, "invalid block: header is from the future");
        Predicates.verifyCommitAgainstTrusted(untrustedState, trustedState, options);
    }

    function verifyChainIdVersion(
        IICS07TendermintMsgs.ChainId memory chainId,
        IICS07TendermintMsgs.Header memory header
    ) internal pure {
        IICS07TendermintMsgs.ChainId memory headerChainId = getChainId(header.signedHeader.header.chainId);
        if (chainId.revisionNumber != headerChainId.revisionNumber) {
            revert HeaderChainIdMismatch(
                header.signedHeader.header.chainId,
                chainId.id
            );
        }
    }

    function validateBasic(
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