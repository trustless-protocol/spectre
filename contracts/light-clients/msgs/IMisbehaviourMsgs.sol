// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import { ISP1Msgs } from "./ISP1Msgs.sol";
import { IICS07TendermintMsgs } from "./IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../../msgs/IICS02ClientMsgs.sol";

/// @title Misbehaviour Program Messages
/// @author gjermundgaraba
/// @notice Defines shared types for the misbehaviour program.
interface IMisbehaviourMsgs {
    struct MsgSubmitMisbehaviour {
        IICS07TendermintMsgs.ClientState clientState;
        IMisbehaviourMsgs.Misbehaviour misbehaviour;
        IICS07TendermintMsgs.ConsensusState trustedConsensusState1;
        IICS07TendermintMsgs.ConsensusState trustedConsensusState2;
        uint128 time;
    }

    /// @notice The public value output for the sp1 misbehaviour program.
    /// @param clientState The client state that was used to verify the misbehaviour.
    /// @param time The time which the misbehaviour was verified in unix nanoseconds.
    /// @param trustedHeight1 The trusted height of header 1
    /// @param trustedHeight2 The trusted height of header 2
    /// @param trustedConsensusState1 The trusted consensus state of header 1
    /// @param trustedConsensusState2 The trusted consensus state of header 2
    struct MisbehaviourOutput {
        IICS02ClientMsgs.Height trustedHeight1;
        IICS02ClientMsgs.Height trustedHeight2;
    }

    struct Misbehaviour {
        IICS07TendermintMsgs.ChainId client_id;
        IICS07TendermintMsgs.Header header1;
        IICS07TendermintMsgs.Header header2;
    }
}
