// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { stdJson } from "forge-std/StdJson.sol";
import { Script } from "forge-std/Script.sol";

import { IERC20 } from "@openzeppelin-contracts/token/ERC20/IERC20.sol";
import { Strings } from "@openzeppelin-contracts/utils/Strings.sol";
import { SafeCast } from "@openzeppelin-contracts/utils/math/SafeCast.sol";

import { ICS20Transfer } from "../contracts/ICS20Transfer.sol";
import { IIBCStore } from "../contracts/interfaces/IIBCStore.sol";
import { IICS20TransferMsgs } from "../contracts/msgs/IICS20TransferMsgs.sol";
import { ICS20Lib } from "../contracts/utils/ICS20Lib.sol";
import { ICS24Host } from "../contracts/utils/ICS24Host.sol";

interface IWETH is IERC20 {
    function deposit() external payable;
}

/// @notice Sends an ERC20/WETH packet through an EVM rollup-side ICS20 contract
///         and verifies that ICS26 stored the packet commitment.
contract EVMRollupSendTransfer is Script {
    using stdJson for string;
    using SafeCast for uint256;

    string internal constant DEFAULT_SOURCE_CLIENT_ID = "evm-rollup-cosmos-0";

    function run() public returns (string memory) {
        ICS20Transfer ics20 = ICS20Transfer(vm.envAddress("ICS20_TRANSFER"));
        address token = vm.envAddress("TOKEN");
        uint256 amount = vm.envUint("AMOUNT");
        string memory receiver = vm.envString("RECEIVER");
        string memory sourceClientId = vm.envOr("SOURCE_CLIENT_ID", DEFAULT_SOURCE_CLIENT_ID);
        string memory memo = vm.envOr("MEMO", string(""));
        uint256 timeoutSeconds = vm.envOr("TIMEOUT_SECONDS", uint256(3600));
        bool wrapEth = vm.envOr("WRAP_ETH", false);
        string memory outputPath = vm.envOr("EVM_ROLLUP_SEND_OUT", string(""));

        require(amount > 0, "AMOUNT must be non-zero");
        require(bytes(receiver).length != 0, "RECEIVER must be non-empty");
        require(timeoutSeconds <= 1 days, "TIMEOUT_SECONDS exceeds router maximum");

        address router = ics20.ics26();
        uint64 timeoutTimestamp = (block.timestamp + timeoutSeconds).toUint64();

        vm.startBroadcast();

        if (wrapEth) {
            IWETH(token).deposit{ value: amount }();
        }

        IERC20(token).approve(address(ics20), amount);

        uint64 sequence = ics20.sendTransfer(
            IICS20TransferMsgs.SendTransferMsg({
                denom: token,
                amount: amount,
                receiver: receiver,
                sourceClient: sourceClientId,
                destPort: ICS20Lib.DEFAULT_PORT_ID,
                timeoutTimestamp: timeoutTimestamp,
                memo: memo
            })
        );

        vm.stopBroadcast();

        bytes32 commitmentKey = ICS24Host.packetCommitmentKeyCalldata(sourceClientId, sequence);
        bytes32 commitment = IIBCStore(router).getCommitment(commitmentKey);
        require(commitment != bytes32(0), "packet commitment was not stored");

        string memory json = "json";
        json.serialize("chainId", block.chainid);
        json.serialize("ics26Router", Strings.toHexString(router));
        json.serialize("ics20Transfer", Strings.toHexString(address(ics20)));
        json.serialize("token", Strings.toHexString(token));
        json.serialize("sourceClientId", sourceClientId);
        json.serialize("sequence", sequence);
        json.serialize("commitmentKey", vm.toString(commitmentKey));
        json.serialize("commitment", vm.toString(commitment));
        string memory finalJson = json.serialize("escrow", Strings.toHexString(ics20.getEscrow(sourceClientId)));

        if (bytes(outputPath).length != 0) {
            vm.writeFile(outputPath, finalJson);
        }

        return finalJson;
    }
}
