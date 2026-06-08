// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { Groth16ICS07Tendermint } from "../../contracts/light-clients/Groth16ICS07Tendermint.sol";
import { IICS02ClientMsgs } from "../../contracts/msgs/IICS02ClientMsgs.sol";
import { ILightClientMsgs } from "../../contracts/msgs/ILightClientMsgs.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { IMembershipMsgs } from "../../contracts/light-clients/msgs/IMembershipMsgs.sol";
import { IGroth16ICS07TendermintErrors } from "../../contracts/light-clients/errors/IGroth16ICS07TendermintErrors.sol";
import { IMembership } from "../../contracts/interfaces/IMembership.sol";

contract DummyMembershipForTrustingPeriod is IMembership {
    function membership(
        bytes32,
        IMembershipMsgs.KVPair[] calldata,
        IMembershipMsgs.MerkleProof[] calldata
    )
        external
        pure { }
}

contract MembershipTrustingPeriodTest is Test {
    uint64 private constant HEIGHT = 10;
    uint32 private constant TRUSTING_PERIOD = 100;
    uint32 private constant UNBONDING_PERIOD = 2 hours;
    uint256 private constant CURRENT_TIME = 1_700_000_000;
    bytes32 private constant APP_HASH = bytes32(uint256(0xA11CE));

    DummyMembershipForTrustingPeriod private membership;

    function setUp() public {
        membership = new DummyMembershipForTrustingPeriod();
        vm.warp(CURRENT_TIME);
    }

    function test_verifyMembershipAcceptsConsensusStateWithinTrustingPeriod() public {
        uint256 consensusTime = block.timestamp - TRUSTING_PERIOD + 1;
        Groth16ICS07Tendermint lightClient = _deploy(_toNanos(consensusTime));

        assertEq(lightClient.verifyMembership(_membershipMsg(_toNanos(consensusTime))), consensusTime);
    }

    function test_verifyMembershipRejectsExpiredConsensusState() public {
        uint256 consensusTime = block.timestamp - TRUSTING_PERIOD;
        Groth16ICS07Tendermint lightClient = _deploy(_toNanos(consensusTime));

        vm.expectRevert(
            abi.encodeWithSelector(
                IGroth16ICS07TendermintErrors.InsufficientTrustingPeriod.selector,
                uint128(TRUSTING_PERIOD),
                uint128(TRUSTING_PERIOD)
            )
        );
        lightClient.verifyMembership(_membershipMsg(_toNanos(consensusTime)));
    }

    function test_verifyMembershipRejectsFutureConsensusState() public {
        uint256 consensusTime = block.timestamp + 1;
        Groth16ICS07Tendermint lightClient = _deploy(_toNanos(consensusTime));

        vm.expectRevert(
            abi.encodeWithSelector(
                IGroth16ICS07TendermintErrors.ProofIsInTheFuture.selector, block.timestamp, consensusTime
            )
        );
        lightClient.verifyMembership(_membershipMsg(_toNanos(consensusTime)));
    }

    function test_verifyNonMembershipRejectsExpiredConsensusState() public {
        uint256 consensusTime = block.timestamp - TRUSTING_PERIOD;
        Groth16ICS07Tendermint lightClient = _deploy(_toNanos(consensusTime));

        vm.expectRevert(
            abi.encodeWithSelector(
                IGroth16ICS07TendermintErrors.InsufficientTrustingPeriod.selector,
                uint128(TRUSTING_PERIOD),
                uint128(TRUSTING_PERIOD)
            )
        );
        lightClient.verifyNonMembership(_nonMembershipMsg(_toNanos(consensusTime)));
    }

    function test_verifyNonMembershipRejectsFutureConsensusState() public {
        uint256 consensusTime = block.timestamp + 1;
        Groth16ICS07Tendermint lightClient = _deploy(_toNanos(consensusTime));

        vm.expectRevert(
            abi.encodeWithSelector(
                IGroth16ICS07TendermintErrors.ProofIsInTheFuture.selector, block.timestamp, consensusTime
            )
        );
        lightClient.verifyNonMembership(_nonMembershipMsg(_toNanos(consensusTime)));
    }

    function _deploy(uint128 consensusTimestamp) private returns (Groth16ICS07Tendermint) {
        IICS07TendermintMsgs.ClientState memory clientState = IICS07TendermintMsgs.ClientState({
            chainId: "test-chain-0",
            trustLevel: IICS07TendermintMsgs.TrustThreshold({ numerator: 1, denominator: 3 }),
            latestHeight: IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: HEIGHT }),
            trustingPeriod: TRUSTING_PERIOD,
            unbondingPeriod: UNBONDING_PERIOD,
            isFrozen: false,
            zkAlgorithm: IICS07TendermintMsgs.SupportedZkAlgorithm.Groth16
        });
        IICS07TendermintMsgs.ConsensusState memory consensusState = _consensusState(consensusTimestamp);

        return new Groth16ICS07Tendermint(
            address(0),
            address(membership),
            address(0),
            address(0),
            abi.encode(clientState),
            keccak256(abi.encode(consensusState)),
            address(0)
        );
    }

    function _membershipMsg(uint128 consensusTimestamp)
        private
        pure
        returns (ILightClientMsgs.MsgVerifyMembership memory msg_)
    {
        bytes memory value = hex"01";
        msg_ = ILightClientMsgs.MsgVerifyMembership({
            height: _height(),
            kvPairs: _kvPairs(value),
            merkleProofs: new IMembershipMsgs.MerkleProof[](0),
            appHash: APP_HASH,
            trustedConsensusState: _consensusState(consensusTimestamp),
            membershipType: IMembershipMsgs.MembershipType.Membership,
            path: _path(),
            value: value
        });
    }

    function _nonMembershipMsg(uint128 consensusTimestamp)
        private
        pure
        returns (ILightClientMsgs.MsgVerifyNonMembership memory msg_)
    {
        msg_ = ILightClientMsgs.MsgVerifyNonMembership({
            height: _height(),
            kvPairs: _kvPairs(bytes("")),
            merkleProofs: new IMembershipMsgs.MerkleProof[](0),
            appHash: APP_HASH,
            trustedConsensusState: _consensusState(consensusTimestamp),
            membershipType: IMembershipMsgs.MembershipType.Membership,
            path: _path()
        });
    }

    function _consensusState(uint128 timestamp) private pure returns (IICS07TendermintMsgs.ConsensusState memory) {
        return IICS07TendermintMsgs.ConsensusState({
            timestamp: timestamp, root: APP_HASH, nextValidatorsHash: bytes32(uint256(0xB0B))
        });
    }

    function _kvPairs(bytes memory value) private pure returns (IMembershipMsgs.KVPair[] memory kvPairs) {
        kvPairs = new IMembershipMsgs.KVPair[](1);
        kvPairs[0] = IMembershipMsgs.KVPair({ path: _path(), value: value });
    }

    function _path() private pure returns (bytes[] memory path) {
        path = new bytes[](1);
        path[0] = bytes("key");
    }

    function _height() private pure returns (IICS02ClientMsgs.Height memory) {
        return IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: HEIGHT });
    }

    function _toNanos(uint256 seconds_) private pure returns (uint128) {
        // forge-lint: disable-next-line(unsafe-typecast)
        return uint128(seconds_ * 1_000_000_000);
    }
}
