// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { ValidatorSetLib } from "contracts/light-clients/spectre/libraries/ValidatorSetLib.sol";
import { IICS07TendermintMsgs } from "contracts/light-clients/spectre/messages/IICS07TendermintMsgs.sol";
import { ISpectreClientErrors } from "contracts/light-clients/spectre/errors/ISpectreClientErrors.sol";
import { Header } from "contracts/light-clients/spectre/libraries/Header.sol";

contract ValidatorSetLibHarness {
    function buildCache(
        bytes32 validatorsHash,
        IICS07TendermintMsgs.ValidatorSet memory validatorSet
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
        IICS07TendermintMsgs.ValidatorSet memory validatorSet = _validatorSet();
        bytes32 validatorsHash = Header.hashValSet(validatorSet);
        bytes memory cacheData = h.buildCache(validatorsHash, validatorSet);
        bytes32 wrongHash = bytes32(uint256(validatorsHash) ^ uint256(1));

        vm.expectRevert(abi.encodeWithSelector(ISpectreClientErrors.CachedValidatorSetCorrupted.selector, wrongHash));
        h.readValidatorCacheHeader(wrongHash, cacheData);
    }

    function test_readValidatorCacheHeader_acceptsMatchingValidatorsHash() public view {
        IICS07TendermintMsgs.ValidatorSet memory validatorSet = _validatorSet();
        bytes32 validatorsHash = Header.hashValSet(validatorSet);
        bytes memory cacheData = h.buildCache(validatorsHash, validatorSet);

        ValidatorSetLib.ValidatorCacheHeader memory header = h.readValidatorCacheHeader(validatorsHash, cacheData);
        assertEq(header.entryCount, 1, "entry count");
        assertEq(header.totalVotingPower, 100, "total voting power");
    }

    function _validatorSet() internal pure returns (IICS07TendermintMsgs.ValidatorSet memory vs) {
        IICS07TendermintMsgs.ValidatorInfo[] memory vals = new IICS07TendermintMsgs.ValidatorInfo[](1);
        vals[0] = IICS07TendermintMsgs.ValidatorInfo({
            valAddress: bytes("validator"), pubKey: bytes32(uint256(1)), votingPower: 100, proposerPriority: 0
        });
        vs = IICS07TendermintMsgs.ValidatorSet({
            validators: vals, hasProposer: false, proposer: vals[0], totalVotingPower: 100
        });
    }
}
