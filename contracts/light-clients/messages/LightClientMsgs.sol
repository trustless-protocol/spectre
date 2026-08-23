// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { ICS02ClientMsgs } from "contracts/core/messages/ICS02ClientMsgs.sol";
import { MembershipMsgs } from "contracts/light-clients/spectre/messages/MembershipMsgs.sol";
import { SpectreMsgs } from "contracts/light-clients/spectre/messages/SpectreMsgs.sol";

/// @title LightClient Messages
/// @notice Interface defining light client messages
interface LightClientMsgs {
    /// @notice Message for querying the membership of a key-value pair in the Merkle root at a given height.
    struct MsgVerifyMembership {
        ICS02ClientMsgs.Height height;
        MembershipMsgs.KVPair[] kvPairs;
        MembershipMsgs.MerkleProof[] merkleProofs;
        bytes32 appHash;
        SpectreMsgs.ConsensusState trustedConsensusState;
        MembershipMsgs.MembershipType membershipType;
        bytes[] path;
        bytes value;
    }

    /// @notice Message for querying the non-membership of a key in the Merkle root at a given height.
    struct MsgVerifyNonMembership {
        ICS02ClientMsgs.Height height;
        MembershipMsgs.KVPair[] kvPairs;
        MembershipMsgs.MerkleProof[] merkleProofs;
        bytes32 appHash;
        SpectreMsgs.ConsensusState trustedConsensusState;
        MembershipMsgs.MembershipType membershipType;
        bytes[] path;
    }

    /// @notice The result of an update operation
    enum UpdateResult {
        /// The update was successful
        Update,
        /// A misbehaviour was detected
        Misbehaviour,
        /// Client is already up to date
        NoOp
    }
}
