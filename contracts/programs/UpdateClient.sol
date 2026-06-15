// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IUpdateClientMsgs } from "../light-clients/msgs/IUpdateClientMsgs.sol";
import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../msgs/IICS02ClientMsgs.sol";
import { IUpdateClient } from "../interfaces/IUpdateClient.sol";
import { Predicates } from "../utils/Predicates.sol";
import { Header } from "../utils/Header.sol";
import { HeightCmp } from "../utils/HeightCmp.sol";
import { ChainId } from "../utils/ChainId.sol";

contract UpdateClient is IUpdateClient {
    error MismatchedRevisionHeight(uint64 expected, uint64 actual);
    error InvalidHeaderHeight(uint64 height);
    error ValSetHashMismatch(bytes32 expected, bytes32 actual);
    error HeaderChainIdMismatch(string expected, string actual);
    error FailedToVerifyHeader(string reason);

    function updateClient(IUpdateClientMsgs.MsgUpdateClient calldata msg_)
        external
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        return _updateClient(msg_, false);
    }

    function updateClientResolved(IUpdateClientMsgs.MsgUpdateClient calldata msg_)
        external
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        return _updateClient(msg_, true);
    }

    function updateClientCachedCurrent(IUpdateClientMsgs.MsgUpdateClient calldata msg_)
        external
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        return _updateClientCachedCurrent(
            msg_,
            Header.chainIdLeafHash(msg_.clientState.chainId),
            Header.bytes32LeafHash(msg_.proposedHeader.signedHeader.header.validatorsHash),
            false
        );
    }

    function updateClientCachedCurrentWithHeaderCache(
        IUpdateClientMsgs.MsgUpdateClient calldata msg_,
        bytes32 chainIdLeafHash,
        bytes32 validatorsHashLeaf
    )
        external
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        return _updateClientCachedCurrent(msg_, chainIdLeafHash, validatorsHashLeaf, false);
    }

    function updateClientCachedCurrentTrustedNextResolvedWithHeaderCache(
        IUpdateClientMsgs.MsgUpdateClient calldata msg_,
        bytes32 chainIdLeafHash,
        bytes32 validatorsHashLeaf
    )
        external
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        return _updateClientCachedCurrent(msg_, chainIdLeafHash, validatorsHashLeaf, true);
    }

    function _updateClientCachedCurrent(
        IUpdateClientMsgs.MsgUpdateClient calldata msg_,
        bytes32 chainIdLeafHash,
        bytes32 validatorsHashLeaf,
        bool trustedNextResolved
    )
        private
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

        verifyHeaderCachedCurrent(
            msg_.proposedHeader,
            chainId,
            options,
            msg_.time,
            msg_.trustedConsensusState,
            chainIdLeafHash,
            validatorsHashLeaf,
            trustedNextResolved,
            msg_.trustedOverlapIndices,
            msg_.signerPubkeys,
            msg_.active
        );
        return _buildOutput(msg_, chainId.revisionNumber);
    }

    function _updateClient(
        IUpdateClientMsgs.MsgUpdateClient calldata msg_,
        bool assumeResolvedValidatorSets
    )
        internal
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        // Reuse the revision number already cached in clientState.latestHeight
        // (populated from the chain ID at client init time). Parsing the
        // chainId string here would cost ~9K gas per update and never change
        // for the lifetime of the client.
        IICS07TendermintMsgs.ChainId memory chainId = IICS07TendermintMsgs.ChainId({
            id: msg_.clientState.chainId, revisionNumber: msg_.clientState.latestHeight.revisionNumber
        });
        IICS07TendermintMsgs.Options memory options = IICS07TendermintMsgs.Options({
            trustThreshold: msg_.clientState.trustLevel,
            trustingPeriod: msg_.clientState.trustingPeriod,
            clockDrift: msg_.clientState.clockDrift
        });

        if (assumeResolvedValidatorSets) {
            verifyHeaderResolved(
                msg_.proposedHeader,
                chainId,
                options,
                msg_.time,
                msg_.trustedConsensusState,
                msg_.trustedOverlapIndices,
                msg_.signerPubkeys,
                msg_.active
            );
        } else {
            verifyHeader(
                msg_.proposedHeader,
                chainId,
                options,
                msg_.time,
                msg_.trustedConsensusState,
                msg_.trustedOverlapIndices,
                msg_.signerPubkeys,
                msg_.active
            );
        }

        return _buildOutput(msg_, chainId.revisionNumber);
    }

    function verifyHeader(
        IICS07TendermintMsgs.Header memory proposedHeader,
        IICS07TendermintMsgs.ChainId memory chainId,
        IICS07TendermintMsgs.Options memory options,
        uint128 time,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState,
        uint32[] memory trustedOverlapIndices,
        bytes32[] memory signerPubkeys,
        bool[] memory active
    )
        internal
        pure
    {
        IICS07TendermintMsgs.ChainId memory headerChainId = ChainId.get(proposedHeader.signedHeader.header.chainId);
        // Checks that the header fields are valid.
        validateBasic(proposedHeader, headerChainId);

        // The tendermint-light-client crate though works on heights that are assumed
        // to have the same revision number. We ensure this here.
        verifyChainIdVersion(chainId, headerChainId);

        // Delegate to tendermint-light-client, which contains the required checks
        // of the new header against the trusted consensus state.
        {
            bytes32 nextValSetHash = Header.hashValSet(proposedHeader.trustedNextValidatorSet);
            if (nextValSetHash != trustedConsensusState.nextValidatorsHash) {
                revert FailedToVerifyHeader("trusted next validator set hash does not match hash stored on chain");
            }

            IICS07TendermintMsgs.TrustedBlockState memory trustedState = IICS07TendermintMsgs.TrustedBlockState({
                chainId: chainId.id,
                headerTime: trustedConsensusState.timestamp,
                height: proposedHeader.trustedHeight.revisionHeight,
                nextValidatorSet: proposedHeader.trustedNextValidatorSet,
                nextValidatorHash: nextValSetHash
            });

            IICS07TendermintMsgs.UntrustedBlockState memory untrustedState = IICS07TendermintMsgs.UntrustedBlockState({
                signedHeader: proposedHeader.signedHeader, validatorSet: proposedHeader.validatorSet
            });

            verifyUpdateHeader(
                untrustedState, trustedState, options, time, trustedOverlapIndices, signerPubkeys, active
            );
        }
    }

    function verifyHeaderResolved(
        IICS07TendermintMsgs.Header memory proposedHeader,
        IICS07TendermintMsgs.ChainId memory chainId,
        IICS07TendermintMsgs.Options memory options,
        uint128 time,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState,
        uint32[] memory trustedOverlapIndices,
        bytes32[] memory signerPubkeys,
        bool[] memory active
    )
        internal
        pure
    {
        IICS07TendermintMsgs.ChainId memory headerChainId = ChainId.get(proposedHeader.signedHeader.header.chainId);
        validateBasicResolved(proposedHeader, headerChainId);
        verifyChainIdVersion(chainId, headerChainId);

        IICS07TendermintMsgs.TrustedBlockState memory trustedState = IICS07TendermintMsgs.TrustedBlockState({
            chainId: chainId.id,
            headerTime: trustedConsensusState.timestamp,
            height: proposedHeader.trustedHeight.revisionHeight,
            nextValidatorSet: proposedHeader.trustedNextValidatorSet,
            nextValidatorHash: trustedConsensusState.nextValidatorsHash
        });

        IICS07TendermintMsgs.UntrustedBlockState memory untrustedState = IICS07TendermintMsgs.UntrustedBlockState({
            signedHeader: proposedHeader.signedHeader, validatorSet: proposedHeader.validatorSet
        });

        if (proposedHeader.validatorSet.validators.length == 0) {
            revert FailedToVerifyHeader("proposed validator set not resolved");
        }

        uint64 trustedNextHeight = proposedHeader.trustedHeight.revisionHeight + 1;
        if (
            proposedHeader.signedHeader.header.height != trustedNextHeight
                && proposedHeader.trustedNextValidatorSet.validators.length == 0
        ) {
            revert FailedToVerifyHeader("trusted next validator set not resolved");
        }

        verifyUpdateHeader(untrustedState, trustedState, options, time, trustedOverlapIndices, signerPubkeys, active);
    }

    function verifyHeaderCachedCurrent(
        IICS07TendermintMsgs.Header memory proposedHeader,
        IICS07TendermintMsgs.ChainId memory chainId,
        IICS07TendermintMsgs.Options memory options,
        uint128 time,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState,
        bytes32 chainIdLeafHash,
        bytes32 validatorsHashLeaf,
        bool trustedNextResolved,
        uint32[] memory trustedOverlapIndices,
        bytes32[] memory signerPubkeys,
        bool[] memory active
    )
        internal
        pure
    {
        IICS07TendermintMsgs.ChainId memory headerChainId = ChainId.get(proposedHeader.signedHeader.header.chainId);
        validateBasicResolved(proposedHeader, headerChainId);
        verifyChainIdVersion(chainId, headerChainId);

        bytes32 headerHash =
            Header.hashHeaderWithCachedLeaves(proposedHeader.signedHeader.header, chainIdLeafHash, validatorsHashLeaf);
        if (headerHash != proposedHeader.signedHeader.commit.blockId.hashData) {
            revert FailedToVerifyHeader("invalid block: header hash mismatch");
        }

        uint64 trustedNextHeight = proposedHeader.trustedHeight.revisionHeight + 1;
        if (proposedHeader.signedHeader.header.height != trustedNextHeight) {
            if (proposedHeader.trustedNextValidatorSet.validators.length == 0) {
                revert FailedToVerifyHeader("trusted next validator set not resolved");
            }
            if (!trustedNextResolved) {
                bytes32 nextValSetHash = Header.hashValSet(proposedHeader.trustedNextValidatorSet);
                if (nextValSetHash != trustedConsensusState.nextValidatorsHash) {
                    revert FailedToVerifyHeader("trusted next validator set hash does not match hash stored on chain");
                }
            }
        }

        IICS07TendermintMsgs.TrustedBlockState memory trustedState = IICS07TendermintMsgs.TrustedBlockState({
            chainId: chainId.id,
            headerTime: trustedConsensusState.timestamp,
            height: proposedHeader.trustedHeight.revisionHeight,
            nextValidatorSet: proposedHeader.trustedNextValidatorSet,
            nextValidatorHash: trustedConsensusState.nextValidatorsHash
        });

        IICS07TendermintMsgs.UntrustedBlockState memory untrustedState = IICS07TendermintMsgs.UntrustedBlockState({
            signedHeader: proposedHeader.signedHeader, validatorSet: proposedHeader.validatorSet
        });

        Predicates.verifyAgainstTrusted(untrustedState, trustedState, options.trustingPeriod, time);
        uint128 drifted = time + uint128(options.clockDrift) * 1_000_000_000;
        require(untrustedState.signedHeader.header.time < drifted, "invalid block: header is from the future");
        Predicates.verifyTrustedCommitOverlapBySignerPubkey(
            untrustedState, trustedState, options, trustedOverlapIndices, signerPubkeys, active
        );
    }

    function verifyUpdateHeader(
        IICS07TendermintMsgs.UntrustedBlockState memory untrustedState,
        IICS07TendermintMsgs.TrustedBlockState memory trustedState,
        IICS07TendermintMsgs.Options memory options,
        uint128 time,
        uint32[] memory trustedOverlapIndices,
        bytes32[] memory signerPubkeys,
        bool[] memory active
    )
        internal
        pure
    {
        Predicates.verifyHeaderMatchesCommit(untrustedState);
        Predicates.verifyAgainstTrusted(untrustedState, trustedState, options.trustingPeriod, time);
        /// Check that the untrusted header is from past.
        uint128 drifted = time + uint128(options.clockDrift) * 1_000_000_000;
        require(untrustedState.signedHeader.header.time < drifted, "invalid block: header is from the future");
        Predicates.verifyTrustedCommitOverlapBySignerPubkey(
            untrustedState, trustedState, options, trustedOverlapIndices, signerPubkeys, active
        );
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

        // We need to ensure that the trusted height (representing the
        // height of the header already on chain for which this client update is
        // based on) must be smaller than height of the new header that we're
        // installing.
        if (HeightCmp.ge(header.trustedHeight, height)) {
            revert InvalidHeaderHeight(height.revisionHeight);
        }

        bytes32 valSetHash = Header.hashValSet(header.validatorSet);
        if (valSetHash != header.signedHeader.header.validatorsHash) {
            revert ValSetHashMismatch(header.signedHeader.header.validatorsHash, valSetHash);
        }
    }

    function validateBasicResolved(
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
        returns (IUpdateClientMsgs.UpdateClientOutput memory output)
    {
        output.clientState = msg_.clientState;
        output.trustedConsensusState = msg_.trustedConsensusState;
        output.newConsensusState = IICS07TendermintMsgs.ConsensusState({
            timestamp: msg_.proposedHeader.signedHeader.header.time,
            root: msg_.proposedHeader.signedHeader.header.appHash,
            nextValidatorsHash: msg_.proposedHeader.signedHeader.header.nextValidatorsHash
        });
        output.time = msg_.time;
        output.trustedHeight = msg_.proposedHeader.trustedHeight;
        output.newHeight = IICS02ClientMsgs.Height({
            revisionNumber: revisionNumber, revisionHeight: msg_.proposedHeader.signedHeader.header.height
        });
    }
}
