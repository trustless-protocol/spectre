// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IMisbehaviourMsgs } from "../light-clients/msgs/IMisbehaviourMsgs.sol";
import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";

interface IMisbehaviour {
    function misbehaviour(
       IICS07TendermintMsgs.ClientState memory clientState,
        IMisbehaviourMsgs.Misbehaviour memory misbehaviour,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState1,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState2,
        uint128 time
    ) external pure returns (IMisbehaviourMsgs.MisbehaviourOutput memory);
}