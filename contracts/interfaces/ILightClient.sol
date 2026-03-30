// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IUpdateClientMsgs } from "../light-clients/msgs/IUpdateClientMsgs.sol";
import { IMisbehaviourMsgs } from "../light-clients/msgs/IMisbehaviourMsgs.sol";
import { ILightClientMsgs } from "../msgs/ILightClientMsgs.sol";

/// @title Light Client Interface
/// @notice Interface for all IBC Eureka light clients to implement.
interface ILightClient {
    /// @notice Updating the client and consensus state
    /// @param updateClientMsg The update client message.
    /// @return The result of the update operation
    function updateClient(bytes calldata updateClientMsg) external returns (ILightClientMsgs.UpdateResult);

    /// @notice Querying the membership of a key-value pair
    /// @dev Notice that this message is not view, as it may update the client state for caching purposes.
    /// @param msg_ The membership message
    /// @return The unix timestamp of the verification height in the counterparty chain in seconds.
    function verifyMembership(ILightClientMsgs.MsgVerifyMembership calldata msg_) external returns (uint256);

    /// @notice Querying the non-membership of a key
    /// @dev Notice that this message is not view, as it may update the client state for caching purposes.
    /// @param msg_ The membership message
    /// @return The unix timestamp of the verification height in the counterparty chain in seconds.
    function verifyNonMembership(ILightClientMsgs.MsgVerifyNonMembership calldata msg_) external returns (uint256);

    /// @notice Misbehaviour handling, moves the light client to the frozen state if misbehaviour is detected
    /// @param misbehaviourMsg The misbehaviour message
    function misbehaviour(bytes calldata misbehaviourMsg) external;

    /// @notice Upgrading the client
    /// @param upgradeMsg The upgrade message
    function upgradeClient(bytes calldata upgradeMsg) external;

    /// @notice Returns the client state.
    /// @return The client state.
    function getClientState() external view returns (bytes memory);
}
