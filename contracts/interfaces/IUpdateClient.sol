// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IUpdateClientMsgs } from "../light-clients/msgs/IUpdateClientMsgs.sol";
import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
interface IUpdateClient {
     function updateClient(
        IUpdateClientMsgs.MsgUpdateClient calldata msg
    ) view external returns (IUpdateClientMsgs.UpdateClientOutput calldata);
}