// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IIBCAppCallbacks } from "../msgs/IIBCAppCallbacks.sol";
import { IIBCSenderCallbacks } from "../interfaces/IIBCSenderCallbacks.sol";

import { ERC165Checker } from "@openzeppelin-contracts/utils/introspection/ERC165Checker.sol";

/// @title IBC Callbacks Library
/// @notice This library provides utility functions for IBC Apps to make callbacks to sender contracts.
library IBCSenderCallbacksLib {
    uint256 internal constant SENDER_CALLBACK_GAS_LIMIT = 100_000;

    /// @notice Emitted when a sender acknowledgement callback reverts or runs out of gas
    /// @param callbackAddress The address of the sender callback contract
    /// @param reason The revert reason, or empty bytes if unavailable
    event IBCSenderAckPacketCallbackError(address indexed callbackAddress, bytes reason);
    /// @notice Emitted when a sender timeout callback reverts or runs out of gas
    /// @param callbackAddress The address of the sender callback contract
    /// @param reason The revert reason, or empty bytes if unavailable
    event IBCSenderTimeoutPacketCallbackError(address indexed callbackAddress, bytes reason);

    /// @notice Checks if the given address implements the IIBCSenderCallbacks interface.
    /// @param sender The address to check
    /// @return bool True if the address implements IIBCSenderCallbacks, false otherwise
    function _supportsCallbacks(address sender) private view returns (bool) {
        return ERC165Checker.supportsInterface(sender, type(IIBCSenderCallbacks).interfaceId);
    }

    /// @notice Make a callback to the sender contract when a packet is acknowledged if it supports IIBCSenderCallbacks.
    /// @param callbackAddress The address of the callback contract
    /// @param success Whether the packet was successfully received by the destination chain
    /// @param msg_ The callback message containing details about the acknowledgement
    function ackPacketCallback(
        address callbackAddress,
        bool success,
        IIBCAppCallbacks.OnAcknowledgementPacketCallback calldata msg_
    )
        internal
    {
        if (_supportsCallbacks(callbackAddress)) {
            try IIBCSenderCallbacks(callbackAddress).onAckPacket{ gas: SENDER_CALLBACK_GAS_LIMIT }(success, msg_) { }
            catch (bytes memory reason) {
                emit IBCSenderAckPacketCallbackError(callbackAddress, reason);
            }
        }
    }

    /// @notice Make a callback to the sender contract when a packet times out if it supports IIBCSenderCallbacks.
    /// @param callbackAddress The address of the callback contract
    /// @param msg_ The callback message containing details about the timeout
    function timeoutPacketCallback(
        address callbackAddress,
        IIBCAppCallbacks.OnTimeoutPacketCallback calldata msg_
    )
        internal
    {
        if (_supportsCallbacks(callbackAddress)) {
            try IIBCSenderCallbacks(callbackAddress).onTimeoutPacket{ gas: SENDER_CALLBACK_GAS_LIMIT }(msg_) { }
            catch (bytes memory reason) {
                emit IBCSenderTimeoutPacketCallbackError(callbackAddress, reason);
            }
        }
    }
}
