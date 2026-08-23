// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { ValidatorSetLib } from "contracts/light-clients/spectre/libraries/ValidatorSetLib.sol";
import { SpectreMsgs } from "contracts/light-clients/spectre/messages/SpectreMsgs.sol";
import { SpectreClientErrors } from "contracts/light-clients/spectre/errors/SpectreClientErrors.sol";
import { Header } from "contracts/light-clients/spectre/libraries/Header.sol";

contract ValidatorSetLibHarness {
    function buildCache(
        bytes32 validatorsHash,
        SpectreMsgs.ValidatorSet memory validatorSet
    )
        external
        pure
        returns (bytes memory)
    {
        return ValidatorSetLib.buildCache(validatorsHash, validatorSet);
    }

    function readValidatorCacheHeader(
        bytes32 validatorsHash,
        bytes memory cacheData
    )
        external
        pure
        returns (ValidatorSetLib.ValidatorCacheHeader memory)
    {
        return ValidatorSetLib.readValidatorCacheHeader(validatorsHash, cacheData);
    }
}

contract ValidatorSetLibTest is Test {
    ValidatorSetLibHarness internal h;

    function setUp() public {
        h = new ValidatorSetLibHarness();
    }

    function test_readValidatorCacheHeader_rejectsMismatchedValidatorsHash() public {
        SpectreMsgs.ValidatorSet memory validatorSet = _validatorSet();
        bytes32 validatorsHash = Header.hashValSet(validatorSet);
        bytes memory cacheData = h.buildCache(validatorsHash, validatorSet);
        bytes32 wrongHash = bytes32(uint256(validatorsHash) ^ uint256(1));

        vm.expectRevert(abi.encodeWithSelector(SpectreClientErrors.CachedValidatorSetCorrupted.selector, wrongHash));
        h.readValidatorCacheHeader(wrongHash, cacheData);
    }

    function test_readValidatorCacheHeader_acceptsMatchingValidatorsHash() public view {
        SpectreMsgs.ValidatorSet memory validatorSet = _validatorSet();
        bytes32 validatorsHash = Header.hashValSet(validatorSet);
        bytes memory cacheData = h.buildCache(validatorsHash, validatorSet);

        ValidatorSetLib.ValidatorCacheHeader memory header = h.readValidatorCacheHeader(validatorsHash, cacheData);
        assertEq(header.entryCount, 1, "entry count");
        assertEq(header.totalVotingPower, 100, "total voting power");
    }

    function _validatorSet() internal pure returns (SpectreMsgs.ValidatorSet memory vs) {
        SpectreMsgs.ValidatorInfo[] memory vals = new SpectreMsgs.ValidatorInfo[](1);
        vals[0] = SpectreMsgs.ValidatorInfo({
            valAddress: bytes("validator"), pubKey: bytes32(uint256(1)), votingPower: 100, proposerPriority: 0
        });
        vs = SpectreMsgs.ValidatorSet({
            validators: vals, hasProposer: false, proposer: vals[0], totalVotingPower: 100
        });
    }
}
