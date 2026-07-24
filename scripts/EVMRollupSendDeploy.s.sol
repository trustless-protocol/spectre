// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { stdJson } from "forge-std/StdJson.sol";
import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";

import { AccessManager } from "@openzeppelin-contracts/access/manager/AccessManager.sol";
import { ERC1967Proxy } from "@openzeppelin-contracts/proxy/ERC1967/ERC1967Proxy.sol";
import { Strings } from "@openzeppelin-contracts/utils/Strings.sol";

import { ICS20Transfer } from "../contracts/ICS20Transfer.sol";
import { ICS26Router } from "../contracts/ICS26Router.sol";
import { IICS02ClientMsgs } from "../contracts/msgs/IICS02ClientMsgs.sol";
import { IBCERC20 } from "../contracts/utils/IBCERC20.sol";
import { Escrow } from "../contracts/utils/Escrow.sol";
import { ICS20Lib } from "../contracts/utils/ICS20Lib.sol";
import { RelayerHelper } from "../contracts/utils/RelayerHelper.sol";
import { ILightClient } from "../contracts/interfaces/ILightClient.sol";
import { ILightClientMsgs } from "../contracts/msgs/ILightClientMsgs.sol";
import { DeployAccessManagerWithRoles } from "./deployments/DeployAccessManagerWithRoles.sol";

/// @notice Placeholder client for an EVM rollup -> Cosmos send-only bootstrap path.
/// @dev This is deliberately unusable for recv/ack/timeout verification.
contract SendOnlyLightClientPlaceholder is ILightClient {
    error SendOnlyClient();

    function updateApplicationState(bytes calldata) external pure returns (ILightClientMsgs.UpdateResult) {
        revert SendOnlyClient();
    }

    function updateConsensusState(bytes calldata) external pure returns (ILightClientMsgs.UpdateResult) {
        revert SendOnlyClient();
    }

    function verifyMembership(ILightClientMsgs.MsgVerifyMembership calldata) external pure returns (uint256) {
        revert SendOnlyClient();
    }

    function verifyNonMembership(ILightClientMsgs.MsgVerifyNonMembership calldata) external pure returns (uint256) {
        revert SendOnlyClient();
    }

    function misbehaviour(bytes calldata) external pure {
        revert SendOnlyClient();
    }

    function unfreeze() external pure {
        revert SendOnlyClient();
    }

    function upgradeClient(bytes calldata) external pure {
        revert SendOnlyClient();
    }

    function getClientState() external pure returns (bytes memory) {
        return bytes("send-only-placeholder");
    }
}

/// @notice Deploys the EVM rollup-side contracts needed for steps 1-3:
///         ICS26/ICS20 deployment, transfer app wiring, and a source client id.
contract EVMRollupSendDeploy is Script, DeployAccessManagerWithRoles {
    using stdJson for string;

    string internal constant DEFAULT_SOURCE_CLIENT_ID = "evm-rollup-cosmos-0";
    string internal constant DEFAULT_COUNTERPARTY_CLIENT_ID = "cosmos-evm-rollup-0";

    error CounterpartyLightClientRequired();

    struct Deployment {
        address accessManager;
        address ics26Router;
        address ics20Transfer;
        address relayerHelper;
        address counterpartyLightClient;
        bool sendOnlyPlaceholder;
    }

    function run() public returns (string memory) {
        string memory sourceClientId = vm.envOr("SOURCE_CLIENT_ID", DEFAULT_SOURCE_CLIENT_ID);
        string memory counterpartyClientId = vm.envOr("COUNTERPARTY_CLIENT_ID", DEFAULT_COUNTERPARTY_CLIENT_ID);
        address counterpartyLightClient = vm.envOr("COUNTERPARTY_LIGHT_CLIENT", address(0));
        bool allowSendOnlyPlaceholder = vm.envOr("ALLOW_SEND_ONLY_PLACEHOLDER", false);
        address permit2 = vm.envOr("PERMIT2", address(0));
        string memory outputPath = vm.envOr("EVM_ROLLUP_DEPLOYMENT_OUT", string(""));

        vm.startBroadcast();

        Deployment memory deployment = _deployCore(permit2, counterpartyLightClient, allowSendOnlyPlaceholder);
        _configureSourceClient(deployment, sourceClientId, counterpartyClientId);

        vm.stopBroadcast();

        string memory finalJson = _deploymentJson(deployment, sourceClientId, counterpartyClientId);

        if (bytes(outputPath).length != 0) {
            vm.writeFile(outputPath, finalJson);
        }

        return finalJson;
    }

    function _deployCore(
        address permit2,
        address counterpartyLightClient,
        bool allowSendOnlyPlaceholder
    )
        private
        returns (Deployment memory deployment)
    {
        AccessManager accessManager = new AccessManager(msg.sender);

        ICS26Router routerLogic = new ICS26Router();
        ERC1967Proxy routerProxy =
            new ERC1967Proxy(address(routerLogic), abi.encodeCall(ICS26Router.initialize, (address(accessManager))));

        ICS20Transfer transferLogic = new ICS20Transfer();
        ERC1967Proxy transferProxy = new ERC1967Proxy(
            address(transferLogic),
            abi.encodeCall(
                ICS20Transfer.initialize,
                (address(routerProxy), address(new Escrow()), address(new IBCERC20()), permit2, address(accessManager))
            )
        );

        RelayerHelper relayerHelper = new RelayerHelper(address(routerProxy));

        bool sendOnlyPlaceholder = counterpartyLightClient == address(0);
        if (sendOnlyPlaceholder) {
            if (!allowSendOnlyPlaceholder) {
                revert CounterpartyLightClientRequired();
            }
            console2.log(
                "WARNING: using SendOnlyLightClientPlaceholder; acks/timeouts cannot be verified and timeout refunds are disabled."
            );
            counterpartyLightClient = address(new SendOnlyLightClientPlaceholder());
        }

        accessManagerSetTargetRoles(accessManager, address(routerProxy), address(transferProxy), false);

        address[] memory relayers = new address[](1);
        relayers[0] = msg.sender;
        accessManagerSetRoles(
            accessManager, relayers, new address[](0), new address[](0), msg.sender, msg.sender, msg.sender
        );

        deployment = Deployment({
            accessManager: address(accessManager),
            ics26Router: address(routerProxy),
            ics20Transfer: address(transferProxy),
            relayerHelper: address(relayerHelper),
            counterpartyLightClient: counterpartyLightClient,
            sendOnlyPlaceholder: sendOnlyPlaceholder
        });
    }

    function _configureSourceClient(
        Deployment memory deployment,
        string memory sourceClientId,
        string memory counterpartyClientId
    )
        private
    {
        bytes[] memory cosmosMerklePrefix = new bytes[](2);
        cosmosMerklePrefix[0] = bytes("ibc");
        cosmosMerklePrefix[1] = bytes("");

        ICS26Router(deployment.ics26Router)
            .addClient(
                sourceClientId,
                IICS02ClientMsgs.CounterpartyInfo({ clientId: counterpartyClientId, merklePrefix: cosmosMerklePrefix }),
                deployment.counterpartyLightClient
            );
        ICS26Router(deployment.ics26Router).addIBCApp(ICS20Lib.DEFAULT_PORT_ID, deployment.ics20Transfer);
    }

    function _deploymentJson(
        Deployment memory deployment,
        string memory sourceClientId,
        string memory counterpartyClientId
    )
        private
        returns (string memory)
    {
        string memory json = "json";
        json.serialize("chainId", block.chainid);
        json.serialize("deployer", Strings.toHexString(msg.sender));
        json.serialize("accessManager", Strings.toHexString(deployment.accessManager));
        json.serialize("ics26Router", Strings.toHexString(deployment.ics26Router));
        json.serialize("ics20Transfer", Strings.toHexString(deployment.ics20Transfer));
        json.serialize("relayerHelper", Strings.toHexString(deployment.relayerHelper));
        json.serialize("sourceClientId", sourceClientId);
        json.serialize("counterpartyClientId", counterpartyClientId);
        json.serialize("sendOnlyPlaceholder", deployment.sendOnlyPlaceholder);
        string memory finalJson =
            json.serialize("counterpartyLightClient", Strings.toHexString(deployment.counterpartyLightClient));

        return finalJson;
    }
}
