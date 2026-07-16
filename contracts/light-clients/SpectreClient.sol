// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// solhint-disable gas-strict-inequalities

import { ISpectreClientMsgs } from "./msgs/ISpectreClientMsgs.sol";
import { IICS07TendermintMsgs } from "./msgs/IICS07TendermintMsgs.sol";
import { IMembershipMsgs } from "./msgs/IMembershipMsgs.sol";
import { ILightClientMsgs } from "../msgs/ILightClientMsgs.sol";
import { IICS02ClientMsgs } from "../msgs/IICS02ClientMsgs.sol";

import { ISpectreClientErrors } from "./errors/ISpectreClientErrors.sol";
import { ISpectreClient } from "./ISpectreClient.sol";
import { IMembership } from "./interfaces/IMembership.sol";
import { IUpdateClient } from "./interfaces/IUpdateClient.sol";
import { IMisbehaviour } from "./interfaces/IMisbehaviour.sol";
import { ILightClient } from "../interfaces/ILightClient.sol";

import { SpectreStore } from "./store/SpectreStore.sol";
import { ValidatorSetLib } from "./store/ValidatorSetLib.sol";
import { Paths } from "./utils/Paths.sol";
import { Header } from "../utils/Header.sol";
import { ChainId } from "../utils/ChainId.sol";
import { SSTORE2 } from "../utils/SSTORE2.sol";
import { AccessControl } from "@openzeppelin-contracts/access/AccessControl.sol";

/// @title SpectreClient
/// @notice ICS-07 Tendermint light client using gnark Groth16. SpectreClient is the only stateful
///         piece: it owns the Store, performs the >2/3 quorum accounting, and applies every state
///         write. Header validation and signature verification are delegatecalled into the
///         UpdateClient / Misbehaviour modules; ICS-23 membership is staticcalled into the
///         Membership module.
contract SpectreClient is ISpectreClientErrors, ISpectreClient, ILightClient, AccessControl {
    using SpectreStore for SpectreStore.Store;

    /// @notice UpdateClient module — delegatecall target that validates headers + signature proofs.
    address private immutable UPDATE_CLIENT_MODULE;
    /// @notice Membership module — staticcall target for ICS-23 (non)membership proofs.
    address private immutable MEMBERSHIP_MODULE;
    /// @notice Misbehaviour module — delegatecall target that validates conflicting headers.
    address private immutable MISBEHAVIOUR_MODULE;

    bytes32 private constant PROOF_SUBMITTER_ROLE = keccak256("PROOF_SUBMITTER_ROLE");

    /// @inheritdoc ISpectreClient
    bytes32 public immutable MISBEHAVIOUR_SUBMITTER_ROLE = keccak256("MISBEHAVIOUR_SUBMITTER_ROLE");

    /// @notice Sets the modules and the initial client, consensus, and pinned validator states.
    constructor(
        address updateClientModule,
        address membershipModule,
        address misbehaviourModule,
        bytes memory clientState_,
        bytes32 consensusState,
        IICS07TendermintMsgs.ValidatorSet memory initialPinnedValidatorSet,
        address roleManager
    ) {
        SpectreStore.Store storage $ = SpectreStore.load();
        IICS07TendermintMsgs.ClientState memory cs = abi.decode(clientState_, (IICS07TendermintMsgs.ClientState));
        $.clientState = cs;

        uint64 parsedRevision = ChainId.get(cs.chainId).revisionNumber;
        require(
            parsedRevision == cs.latestHeight.revisionNumber,
            MismatchedRevisionHeights(parsedRevision, cs.latestHeight.revisionNumber)
        );
        $.consensusStateHashes[cs.latestHeight.revisionHeight] = consensusState;

        UPDATE_CLIENT_MODULE = updateClientModule;
        MEMBERSHIP_MODULE = membershipModule;
        MISBEHAVIOUR_MODULE = misbehaviourModule;

        require(
            cs.trustingPeriod + cs.clockDrift <= cs.unbondingPeriod,
            TrustingPeriodTooLong(cs.trustingPeriod, cs.unbondingPeriod)
        );

        _setPinnedValidatorSet(initialPinnedValidatorSet);
        _storePinnedValidatorSetSnapshot(cs.latestHeight.revisionHeight);

        if (roleManager == address(0)) {
            _grantRole(PROOF_SUBMITTER_ROLE, address(0));
            _grantRole(MISBEHAVIOUR_SUBMITTER_ROLE, address(0));
        } else {
            _grantRole(DEFAULT_ADMIN_ROLE, roleManager);
            _grantRole(PROOF_SUBMITTER_ROLE, roleManager);
            _grantRole(MISBEHAVIOUR_SUBMITTER_ROLE, roleManager);
        }
    }

    // ============ Views ============

    /// @inheritdoc ILightClient
    function getClientState() external view returns (bytes memory) {
        return abi.encode(SpectreStore.load().clientState);
    }

    /// @inheritdoc ISpectreClient
    function getPinnedValidatorsHash() external view returns (bytes32) {
        return SpectreStore.load().pinnedValidatorsHash;
    }

    /// @inheritdoc ISpectreClient
    function getPinnedValidatorSet()
        external
        view
        returns (uint32[] memory indices, bytes32[] memory pubkeys, uint64[] memory votingPowers)
    {
        SpectreStore.Store storage $ = SpectreStore.load();
        SpectreStore.PinnedValidatorSetSnapshot memory snapshot = $.currentSnapshot();
        (ValidatorSetLib.ValidatorCacheHeader memory cacheHeader, bytes memory cacheData) =
            ValidatorSetLib.readCache(snapshot);

        indices = new uint32[](cacheHeader.entryCount);
        pubkeys = new bytes32[](cacheHeader.entryCount);
        votingPowers = new uint64[](cacheHeader.entryCount);
        for (uint16 i = 0; i < cacheHeader.entryCount; i++) {
            indices[i] = uint32(i);
            (votingPowers[i], pubkeys[i]) =
                ValidatorSetLib.resolveValidator(snapshot.validatorsHash, cacheData, cacheHeader.entryCount, i);
        }
    }

    // ============ Update ============

    /// @inheritdoc ILightClient
    function updateApplicationState(bytes calldata updateMsg)
        external
        notFrozen
        onlyProofSubmitter
        returns (ILightClientMsgs.UpdateResult)
    {
        ISpectreClientMsgs.MsgUpdateApplicationState memory msg_ =
            abi.decode(updateMsg, (ISpectreClientMsgs.MsgUpdateApplicationState));

        SpectreStore.Store storage $ = SpectreStore.load();
        _requireFreshness($, msg_.time);

        // Quorum accounting first: it gates the (expensive) signature proof and must revert before
        // the SignatureVerifier is reached when the proof metadata is sub-quorum or malformed.
        _verifyPinnedQuorum($, msg_.proof, msg_.proposedHeader);

        // Header validation + the batched signature proof run in the delegatecalled module.
        ISpectreClientMsgs.VerifyHeaderOutput memory output = abi.decode(
            _delegate(UPDATE_CLIENT_MODULE, abi.encodeCall(IUpdateClient.verifyHeader, (msg_))),
            (ISpectreClientMsgs.VerifyHeaderOutput)
        );

        ILightClientMsgs.UpdateResult updateResult = _checkUpdateResult($, output);

        if (updateResult == ILightClientMsgs.UpdateResult.Update) {
            require(
                output.newHeight.revisionHeight > $.clientState.latestHeight.revisionHeight,
                NonMonotonicHeightUpdate($.clientState.latestHeight.revisionHeight, output.newHeight.revisionHeight)
            );
            $.clientState.latestHeight = output.newHeight;
            $.consensusStateHashes[output.newHeight.revisionHeight] = keccak256(abi.encode(output.newConsensusState));
            emit ClientUpdated(output.newHeight.revisionHeight);
        } else if (updateResult == ILightClientMsgs.UpdateResult.Misbehaviour) {
            $.clientState.isFrozen = true;
            emit ClientFrozen();
        }
        return updateResult;
    }

    /// @inheritdoc ILightClient
    function updateConsensusState(bytes calldata updateMsg)
        external
        notFrozen
        onlyProofSubmitter
        returns (ILightClientMsgs.UpdateResult)
    {
        ISpectreClientMsgs.MsgUpdateConsensusState memory msg_ =
            abi.decode(updateMsg, (ISpectreClientMsgs.MsgUpdateConsensusState));
        ISpectreClientMsgs.MsgUpdateApplicationState memory app = msg_.update;

        SpectreStore.Store storage $ = SpectreStore.load();
        _requireFreshness($, app.time);

        // Quorum accounting first: it gates the (expensive) signature proof and must revert before
        // the SignatureVerifier is reached when the proof metadata is sub-quorum or malformed.
        _verifyPinnedQuorum($, app.proof, app.proposedHeader);

        // Header validation + the batched signature proof run in the delegatecalled module.
        ISpectreClientMsgs.VerifyHeaderOutput memory output = abi.decode(
            _delegate(UPDATE_CLIENT_MODULE, abi.encodeCall(IUpdateClient.verifyHeader, (app))),
            (ISpectreClientMsgs.VerifyHeaderOutput)
        );

        ILightClientMsgs.UpdateResult updateResult = _checkUpdateResult($, output);

        if (updateResult == ILightClientMsgs.UpdateResult.Misbehaviour) {
            $.clientState.isFrozen = true;
            emit ClientFrozen();
            return updateResult;
        }

        bytes32 newValidatorsHash = Header.hashValSet(msg_.newValidatorSet);
        require(
            newValidatorsHash == app.proposedHeader.signedHeader.header.nextValidatorsHash,
            MismatchedValidatorHashes(app.proposedHeader.signedHeader.header.nextValidatorsHash, newValidatorsHash)
        );

        if (updateResult == ILightClientMsgs.UpdateResult.Update) {
            require(
                output.newHeight.revisionHeight > $.clientState.latestHeight.revisionHeight,
                NonMonotonicHeightUpdate($.clientState.latestHeight.revisionHeight, output.newHeight.revisionHeight)
            );
            _setPinnedValidatorSet(msg_.newValidatorSet);
            $.clientState.latestHeight = output.newHeight;
            $.consensusStateHashes[output.newHeight.revisionHeight] = keccak256(abi.encode(output.newConsensusState));
            _storePinnedValidatorSetSnapshot(output.newHeight.revisionHeight);
            emit ClientUpdated(output.newHeight.revisionHeight);
            emit ConsensusStateUpdated(output.newHeight.revisionHeight, newValidatorsHash);
        } else {
            // NoOp: re-pin the validator set at the already-trusted height.
            require(
                output.newHeight.revisionHeight == $.clientState.latestHeight.revisionHeight,
                NonMonotonicHeightUpdate($.clientState.latestHeight.revisionHeight, output.newHeight.revisionHeight)
            );
            _setPinnedValidatorSet(msg_.newValidatorSet);
            _storePinnedValidatorSetSnapshot(output.newHeight.revisionHeight);
            emit ConsensusStateUpdated(output.newHeight.revisionHeight, newValidatorsHash);
        }
        return updateResult;
    }

    /// @notice Checks whether a verified header advances, no-ops, or reveals misbehaviour.
    /// @dev Misbehaviour if the consensus state at the new height differs from the mapping,
    ///      or the timestamp is not increasing.
    function _checkUpdateResult(
        SpectreStore.Store storage $,
        ISpectreClientMsgs.VerifyHeaderOutput memory output
    )
        private
        view
        returns (ILightClientMsgs.UpdateResult)
    {
        bytes32 consensusStateHash = $.consensusStateHashes[output.newHeight.revisionHeight];
        if (consensusStateHash == bytes32(0)) {
            return ILightClientMsgs.UpdateResult.Update;
        } else if (
            consensusStateHash != keccak256(abi.encode(output.newConsensusState))
                || output.trustedConsensusState.timestamp >= output.newConsensusState.timestamp
        ) {
            return ILightClientMsgs.UpdateResult.Misbehaviour;
        } else {
            return ILightClientMsgs.UpdateResult.NoOp;
        }
    }

    /// @notice Sums the verified signers' voting power against the current pinned set, requiring >2/3.
    function _verifyPinnedQuorum(
        SpectreStore.Store storage $,
        ISpectreClientMsgs.BatchProof memory proof_,
        IICS07TendermintMsgs.Header memory header
    )
        private
        view
    {
        _verifyQuorum(proof_, header.signedHeader.commit.commitSigs, $.currentSnapshot());
    }

    /// @notice Sums the verified signers' voting power against the snapshot trusted at the header's
    ///         trusted height, requiring >2/3.
    function _verifyMisbehaviourQuorum(
        SpectreStore.Store storage $,
        ISpectreClientMsgs.BatchProof memory proof_,
        IICS07TendermintMsgs.Header memory header
    )
        private
        view
    {
        _verifyQuorum(proof_, header.signedHeader.commit.commitSigs, $.snapshotAt(header.trustedHeight.revisionHeight));
    }

    function _verifyQuorum(
        ISpectreClientMsgs.BatchProof memory proof_,
        IICS07TendermintMsgs.CommitSig[] memory commitSigs,
        SpectreStore.PinnedValidatorSetSnapshot memory snapshot
    )
        private
        view
    {
        require(
            proof_.signerIndices.length == proof_.bucket && proof_.pinnedValidatorIndices.length == proof_.bucket
                && proof_.signerPubkeys.length == proof_.bucket && proof_.active.length == proof_.bucket,
            BatchLengthMismatch()
        );

        (ValidatorSetLib.ValidatorCacheHeader memory cacheHeader, bytes memory cacheData) =
            ValidatorSetLib.readCache(snapshot);
        uint64 totalVotingPower = snapshot.totalVotingPower;
        uint64 accumulatedVotingPower = 0;
        uint256 seenPinned = 0;
        bool hasPrevCommitSigner = false;
        uint32 prevCommitSigner = 0;

        for (uint256 i = 0; i < proof_.bucket; i++) {
            if (!proof_.active[i]) {
                continue;
            }

            uint32 commitIdx = proof_.signerIndices[i];
            if (hasPrevCommitSigner) {
                require(commitIdx > prevCommitSigner, DuplicateSigner(commitIdx));
            }
            hasPrevCommitSigner = true;
            prevCommitSigner = commitIdx;

            uint32 pinnedIdx = proof_.pinnedValidatorIndices[i];
            require(pinnedIdx < cacheHeader.entryCount, SignerIndexOutOfRange(pinnedIdx));
            uint256 mask = uint256(1) << pinnedIdx;
            require((seenPinned & mask) == 0, DuplicateSigner(pinnedIdx));
            seenPinned |= mask;

            (uint64 votingPower, bytes32 pubKey) =
                ValidatorSetLib.resolveValidator(snapshot.validatorsHash, cacheData, cacheHeader.entryCount, pinnedIdx);
            require(pubKey == proof_.signerPubkeys[i], PubkeyMismatch(pinnedIdx));
            accumulatedVotingPower += votingPower;
        }

        require(
            uint256(accumulatedVotingPower) * 3 > uint256(totalVotingPower) * 2,
            InsufficientVotingPower(accumulatedVotingPower, totalVotingPower)
        );

        ValidatorSetLib.requireProofSignersCommitSigs(commitSigs, proof_.signerIndices, proof_.active);
    }

    // ============ Membership ============

    /// @inheritdoc ILightClient
    function verifyMembership(ILightClientMsgs.MsgVerifyMembership calldata msg_)
        external
        notFrozen
        onlyProofSubmitter
        returns (uint256)
    {
        require(msg_.value.length > 0, EmptyValue());
        return _membership(
            msg_.height,
            msg_.kvPairs,
            msg_.merkleProofs,
            msg_.appHash,
            msg_.trustedConsensusState,
            msg_.membershipType,
            msg_.path,
            msg_.value
        );
    }

    /// @inheritdoc ILightClient
    function verifyNonMembership(ILightClientMsgs.MsgVerifyNonMembership calldata msg_)
        external
        notFrozen
        onlyProofSubmitter
        returns (uint256)
    {
        return _membership(
            msg_.height,
            msg_.kvPairs,
            msg_.merkleProofs,
            msg_.appHash,
            msg_.trustedConsensusState,
            msg_.membershipType,
            msg_.path,
            bytes("")
        );
    }

    function _membership(
        IICS02ClientMsgs.Height calldata height,
        IMembershipMsgs.KVPair[] calldata kvPairs,
        IMembershipMsgs.MerkleProof[] calldata merkleProofs,
        bytes32 appHash,
        IICS07TendermintMsgs.ConsensusState calldata trustedConsensusState,
        IMembershipMsgs.MembershipType membershipType,
        bytes[] calldata kvPath,
        bytes memory kvValue
    )
        private
        returns (uint256)
    {
        if (membershipType == IMembershipMsgs.MembershipType.Membership) {
            return _handleMembership(height, kvPairs, merkleProofs, appHash, trustedConsensusState, kvPath, kvValue);
        }

        revert UnknownMembershipType(uint8(membershipType));
    }

    function _handleMembership(
        IICS02ClientMsgs.Height calldata height,
        IMembershipMsgs.KVPair[] calldata kvPairs,
        IMembershipMsgs.MerkleProof[] calldata merkleProofs,
        bytes32 appHash,
        IICS07TendermintMsgs.ConsensusState calldata trustedConsensusState,
        bytes[] calldata kvPath,
        bytes memory kvValue
    )
        private
        returns (uint256)
    {
        require(
            kvPairs.length > 0 && kvPairs.length <= type(uint16).max,
            LengthIsOutOfRange(kvPairs.length, 1, type(uint16).max)
        );

        _validateMembershipInput(appHash, height.revisionHeight, trustedConsensusState);

        {
            bool found = false;
            for (uint256 i = 0; i < kvPairs.length; ++i) {
                if (!Paths.equal(kvPairs[i].path, kvPath)) {
                    continue;
                }

                bytes memory value = kvPairs[i].value;
                require(
                    value.length == kvValue.length && keccak256(value) == keccak256(kvValue),
                    MembershipProofValueMismatch(kvValue, value)
                );

                found = true;
                break;
            }
            require(found, MembershipProofKeyNotFound(kvPath));
        }

        IMembership(MEMBERSHIP_MODULE).verifyMembership(appHash, kvPairs, merkleProofs);

        return _nanosToSeconds(trustedConsensusState.timestamp);
    }

    function _validateMembershipInput(
        bytes32 commitmentRoot,
        uint64 proofHeight,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState
    )
        private
        view
    {
        SpectreStore.Store storage $ = SpectreStore.load();
        bytes32 trustedConsensusStateHash = keccak256(abi.encode(trustedConsensusState));
        bytes32 storedConsensusStateHash = $.getConsensusStateHash(proofHeight);
        require(
            trustedConsensusStateHash == storedConsensusStateHash,
            ConsensusStateHashMismatch(storedConsensusStateHash, trustedConsensusStateHash)
        );

        require(
            commitmentRoot == trustedConsensusState.root,
            ConsensusStateRootMismatch(trustedConsensusState.root, commitmentRoot)
        );

        _validateConsensusStateTrustingPeriod(trustedConsensusState.timestamp);
    }

    function _validateConsensusStateTrustingPeriod(uint128 consensusStateTimestamp) private view {
        SpectreStore.Store storage $ = SpectreStore.load();
        uint256 consensusStateTimestampSeconds = _nanosToSeconds(consensusStateTimestamp);
        require(
            // Membership proof freshness is defined against the destination chain clock.
            // forge-lint: disable-next-line(block-timestamp)
            consensusStateTimestampSeconds <= block.timestamp,
            ProofIsInTheFuture(block.timestamp, consensusStateTimestampSeconds)
        );

        // forge-lint: disable-next-line(unsafe-typecast)
        uint128 durationSinceConsensusState = uint128(block.timestamp - consensusStateTimestampSeconds);
        require(
            durationSinceConsensusState < $.clientState.trustingPeriod,
            InsufficientTrustingPeriod(durationSinceConsensusState, uint128($.clientState.trustingPeriod))
        );
    }

    // ============ Misbehaviour ============

    /// @dev Standalone misbehaviour is accepted only when both conflicting headers carry
    ///      proof-backed >2/3 Ed25519 quorums over the pinned validator set.
    /// @inheritdoc ILightClient
    function misbehaviour(bytes calldata misbehaviourMsg) external notFrozen onlyMisbehaviourSubmitter {
        ISpectreClientMsgs.MsgSubmitMisbehaviour memory msg_ =
            abi.decode(misbehaviourMsg, (ISpectreClientMsgs.MsgSubmitMisbehaviour));

        SpectreStore.Store storage $ = SpectreStore.load();
        // Bind the submitter-chosen time to the wall clock. Without this, `time` could be picked
        // inside the trusting period of an arbitrarily old trusted state, letting a long-expired
        // (unbonded) snapshot validator set freeze the client.
        _requireFreshness($, msg_.time);

        // Quorum accounting first: it gates the (expensive) signature proof and must revert before
        // the SignatureVerifier is reached when the proof metadata is sub-quorum or malformed.
        _verifyMisbehaviourQuorum($, msg_.proof1, msg_.misbehaviour.header1);
        _verifyMisbehaviourQuorum($, msg_.proof2, msg_.misbehaviour.header2);

        // Header validation + the two batched signature proofs run in the delegatecalled module.
        _delegate(MISBEHAVIOUR_MODULE, abi.encodeCall(IMisbehaviour.verifyMisbehaviour, (msg_)));

        $.clientState.isFrozen = true;
        emit ClientFrozen();
    }

    /// @inheritdoc ISpectreClient
    function unfreeze() external override(ISpectreClient, ILightClient) onlyRole(DEFAULT_ADMIN_ROLE) {
        SpectreStore.Store storage $ = SpectreStore.load();
        require($.clientState.isFrozen, ClientNotFrozen());
        $.clientState.isFrozen = false;
        emit ClientUnfrozen();
    }

    /// @inheritdoc ILightClient
    function upgradeClient(bytes calldata) external pure {
        // NOTE: This feature will not be supported.
        revert FeatureNotSupported();
    }

    // ============ Internal state writes ============

    function _setPinnedValidatorSet(IICS07TendermintMsgs.ValidatorSet memory validatorSet) private {
        SpectreStore.Store storage $ = SpectreStore.load();
        bytes32 validatorsHash = Header.hashValSet(validatorSet);
        bytes memory cacheData = ValidatorSetLib.buildCache(validatorsHash, validatorSet);
        ValidatorSetLib.ValidatorCacheHeader memory cacheHeader =
            ValidatorSetLib.readValidatorCacheHeader(validatorsHash, cacheData);
        $.pinnedValidatorSetPointer = SSTORE2.write(cacheData);
        $.pinnedValidatorsHash = validatorsHash;
        $.pinnedTotalVotingPower = cacheHeader.totalVotingPower;
        $.pinnedEntryCount = cacheHeader.entryCount;
    }

    function _storePinnedValidatorSetSnapshot(uint64 height) private {
        SpectreStore.Store storage $ = SpectreStore.load();
        bool exists = $.snapshots[height].pointer != address(0);
        $.snapshots[height] = SpectreStore.PinnedValidatorSetSnapshot({
            validatorsHash: $.pinnedValidatorsHash,
            pointer: $.pinnedValidatorSetPointer,
            totalVotingPower: $.pinnedTotalVotingPower,
            entryCount: $.pinnedEntryCount
        });
        if (!exists) {
            $.snapshotHeights.push(height);
        }
    }

    // ============ Helpers ============

    /// @notice Delegatecalls a module, bubbling up any revert data (so custom errors survive).
    function _delegate(address target, bytes memory data) private returns (bytes memory) {
        (bool ok, bytes memory ret) = target.delegatecall(data);
        if (!ok) {
            assembly ("memory-safe") {
                revert(add(ret, 0x20), mload(ret))
            }
        }
        return ret;
    }

    /// @notice Requires the proof time is not in the future and within clock drift of now.
    function _requireFreshness(SpectreStore.Store storage $, uint128 time) private view {
        uint256 timeSeconds = _nanosToSeconds(time);
        // forge-lint: disable-next-line(block-timestamp)
        require(timeSeconds <= block.timestamp, ProofIsInTheFuture(block.timestamp, timeSeconds));
        require(block.timestamp - timeSeconds <= $.clientState.clockDrift, ProofIsTooOld(block.timestamp, timeSeconds));
    }

    function _nanosToSeconds(uint256 nanos) private pure returns (uint256) {
        return nanos / 1e9;
    }

    modifier notFrozen() {
        require(!SpectreStore.load().clientState.isFrozen, FrozenClientState());
        _;
    }

    modifier onlyProofSubmitter() {
        if (!hasRole(PROOF_SUBMITTER_ROLE, address(0))) {
            _checkRole(PROOF_SUBMITTER_ROLE);
        }
        _;
    }

    modifier onlyMisbehaviourSubmitter() {
        if (!hasRole(MISBEHAVIOUR_SUBMITTER_ROLE, address(0))) {
            _checkRole(MISBEHAVIOUR_SUBMITTER_ROLE);
        }
        _;
    }
}
