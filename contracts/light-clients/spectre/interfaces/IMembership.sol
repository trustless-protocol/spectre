// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { MembershipMsgs } from "contracts/light-clients/spectre/messages/MembershipMsgs.sol";

/// @title IMembership
/// @notice Interface for the Membership module (ICS-23 Merkle proof verification).
///         Stateless and pure math over calldata; invoked via staticcall by SpectreClient.
interface IMembership {
    /// @notice Verifies ICS-23 (non)membership proofs against the trusted appHash root.
    function verifyMembership(
        bytes32 appHash,
        MembershipMsgs.KVPair[] calldata kvPairs,
        MembershipMsgs.MerkleProof[] calldata merkleProofs
    )
        external
        view;
}
