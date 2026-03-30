// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IMembershipMsgs } from "../light-clients/msgs/IMembershipMsgs.sol";
interface IMembership {
    function membership(
        bytes32 appHash,
        IMembershipMsgs.KVPair[] calldata kvPairs,
        IMembershipMsgs.MerkleProof[] calldata merkleProofs
    ) external returns (IMembershipMsgs.MembershipOutput memory);
}