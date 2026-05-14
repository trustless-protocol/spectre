// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Script } from "forge-std/Script.sol";
import { Misbehaviour } from "../contracts/programs/Misbehaviour.sol";

contract DeployMisbehaviourOnly is Script {
    function run() external returns (address deployed) {
        vm.startBroadcast();
        deployed = address(new Misbehaviour());
        vm.stopBroadcast();
    }
}
