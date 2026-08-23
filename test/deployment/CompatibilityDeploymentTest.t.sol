// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";
import { ProductionConfigLib } from "scripts/deployments/ProductionConfigLib.sol";

contract ProductionConfigHarness {
    function verifierList(address n4) external pure returns (address[] memory) {
        return ProductionConfigLib.verifierList(n4);
    }
}

/// @notice Locks the deployment topology that must move atomically with prover support.
contract CompatibilityDeploymentTest is Test {
    function testSupportedVerifierTopologyIsN4Only() public {
        ProductionConfigHarness harness = new ProductionConfigHarness();
        address expected = address(0x4);
        address[] memory verifiers = harness.verifierList(expected);

        assertEq(verifiers.length, 1);
        assertEq(verifiers[0], expected);
    }

    function testCompatibilityAndReleaseManifestsAreReadable() public view {
        assertGt(bytes(vm.readFile("scripts/solidity-refactor/baseline.json")).length, 0);
        assertGt(bytes(vm.readFile("scripts/solidity-refactor/ownership-manifest.json")).length, 0);
        assertGt(bytes(vm.readFile("scripts/solidity-refactor/verifier-manifest.json")).length, 0);
    }
}
