// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { SpectreClientMsgs } from "contracts/light-clients/spectre/messages/SpectreClientMsgs.sol";

/// @title IUpdateClient
/// @notice Interface for the UpdateClient module. The module validates a proposed header and its
///         batched signature proof; it is invoked via delegatecall by SpectreClient.
interface IUpdateClient {
    /// @notice Validates a proposed header (basic checks + trusted-state binding + signature proof).
    /// @param msg_ The application-state update message.
    /// @return The verified header output (trusted/new consensus states and heights).
    function verifyHeader(SpectreClientMsgs.MsgUpdateApplicationState calldata msg_)
        external
        returns (SpectreClientMsgs.VerifyHeaderOutput memory);
}
