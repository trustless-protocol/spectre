// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { SpectreClient } from "contracts/light-clients/spectre/SpectreClient.sol";
import { ICS02ClientMsgs } from "contracts/core/messages/ICS02ClientMsgs.sol";
import { LightClientMsgs } from "contracts/light-clients/messages/LightClientMsgs.sol";
import { SpectreMsgs } from "contracts/light-clients/spectre/messages/SpectreMsgs.sol";
import { MembershipMsgs } from "contracts/light-clients/spectre/messages/MembershipMsgs.sol";
import { SpectreClientErrors } from "contracts/light-clients/spectre/errors/SpectreClientErrors.sol";
import { IMembership } from "contracts/light-clients/spectre/interfaces/IMembership.sol";
import { Header } from "contracts/light-clients/spectre/libraries/Header.sol";

contract DummyMembershipForTrustingPeriod is IMembership {
    function verifyMembership(
        bytes32,
        MembershipMsgs.KVPair[] calldata,
        MembershipMsgs.MerkleProof[] calldata
    )
        external
        pure { }
}

contract MembershipTrustingPeriodTest is Test {
    uint64 private constant HEIGHT = 10;
    uint32 private constant TRUSTING_PERIOD = 100;
    uint32 private constant UNBONDING_PERIOD = 2 hours;
    uint32 private constant CLOCK_DRIFT = 1800;
    uint256 private constant CURRENT_TIME = 1_700_000_000;
    bytes32 private constant APP_HASH = bytes32(uint256(0xA11CE));

    DummyMembershipForTrustingPeriod private membership;

    function setUp() public {
        membership = new DummyMembershipForTrustingPeriod();
        vm.warp(CURRENT_TIME);
    }

    function test_verifyMembershipAcceptsConsensusStateWithinTrustingPeriod() public {
        uint256 consensusTime = block.timestamp - TRUSTING_PERIOD + 1;
        SpectreClient lightClient = _deploy(_toNanos(consensusTime));

        assertEq(lightClient.verifyMembership(_membershipMsg(_toNanos(consensusTime))), consensusTime);
    }

    function test_verifyMembershipRejectsExpiredConsensusState() public {
        uint256 consensusTime = block.timestamp - TRUSTING_PERIOD;
        SpectreClient lightClient = _deploy(_toNanos(consensusTime));
        // Built before expectRevert: _membershipMsg -> _consensusState -> Header.hashValSet does a
        // sha256 precompile STATICCALL, which vm.expectRevert would otherwise catch as "the next
        // call" instead of the verifyMembership call below.
        LightClientMsgs.MsgVerifyMembership memory msg_ = _membershipMsg(_toNanos(consensusTime));

        vm.expectRevert(
            abi.encodeWithSelector(
                SpectreClientErrors.InsufficientTrustingPeriod.selector,
                uint128(TRUSTING_PERIOD),
                uint128(TRUSTING_PERIOD)
            )
        );
        lightClient.verifyMembership(msg_);
    }

    function test_verifyMembershipAcceptsFutureConsensusStateWithinClockDrift() public {
        uint256 consensusTime = block.timestamp + CLOCK_DRIFT;
        SpectreClient lightClient = _deploy(_toNanos(consensusTime));

        assertEq(lightClient.verifyMembership(_membershipMsg(_toNanos(consensusTime))), consensusTime);
    }

    function test_verifyMembershipRejectsFutureConsensusStateBeyondClockDrift() public {
        uint256 consensusTime = block.timestamp + CLOCK_DRIFT + 1;
        SpectreClient lightClient = _deploy(_toNanos(consensusTime));
        // Built before expectRevert — see test_verifyMembershipRejectsExpiredConsensusState.
        LightClientMsgs.MsgVerifyMembership memory msg_ = _membershipMsg(_toNanos(consensusTime));

        vm.expectRevert(
            abi.encodeWithSelector(SpectreClientErrors.ProofIsInTheFuture.selector, block.timestamp, consensusTime)
        );
        lightClient.verifyMembership(msg_);
    }

    function test_verifyNonMembershipRejectsExpiredConsensusState() public {
        uint256 consensusTime = block.timestamp - TRUSTING_PERIOD;
        SpectreClient lightClient = _deploy(_toNanos(consensusTime));
        // Built before expectRevert — see test_verifyMembershipRejectsExpiredConsensusState.
        LightClientMsgs.MsgVerifyNonMembership memory msg_ = _nonMembershipMsg(_toNanos(consensusTime));

        vm.expectRevert(
            abi.encodeWithSelector(
                SpectreClientErrors.InsufficientTrustingPeriod.selector,
                uint128(TRUSTING_PERIOD),
                uint128(TRUSTING_PERIOD)
            )
        );
        lightClient.verifyNonMembership(msg_);
    }

    function test_verifyNonMembershipAcceptsFutureConsensusStateWithinClockDrift() public {
        uint256 consensusTime = block.timestamp + CLOCK_DRIFT;
        SpectreClient lightClient = _deploy(_toNanos(consensusTime));

        assertEq(lightClient.verifyNonMembership(_nonMembershipMsg(_toNanos(consensusTime))), consensusTime);
    }

    function test_verifyNonMembershipRejectsFutureConsensusStateBeyondClockDrift() public {
        uint256 consensusTime = block.timestamp + CLOCK_DRIFT + 1;
        SpectreClient lightClient = _deploy(_toNanos(consensusTime));
        // Built before expectRevert — see test_verifyMembershipRejectsExpiredConsensusState.
        LightClientMsgs.MsgVerifyNonMembership memory msg_ = _nonMembershipMsg(_toNanos(consensusTime));

        vm.expectRevert(
            abi.encodeWithSelector(SpectreClientErrors.ProofIsInTheFuture.selector, block.timestamp, consensusTime)
        );
        lightClient.verifyNonMembership(msg_);
    }

    function test_constructorRejectsZeroTrustingPeriod() public {
        SpectreMsgs.ClientState memory clientState = _clientState();
        clientState.trustingPeriod = 0;
        // Built before expectRevert — see test_verifyMembershipRejectsExpiredConsensusState: the
        // sha256 precompile call inside _consensusState must not land inside the armed window.
        SpectreMsgs.ConsensusState memory consensusState = _consensusState(_toNanos(block.timestamp));

        vm.expectRevert(
            abi.encodeWithSelector(
                SpectreClientErrors.LengthIsOutOfRange.selector, uint256(0), uint256(1), uint256(type(uint32).max)
            )
        );
        _deployWithConsensusState(consensusState, clientState);
    }

    function test_constructorRejectsZeroClockDrift() public {
        SpectreMsgs.ClientState memory clientState = _clientState();
        clientState.clockDrift = 0;
        // Built before expectRevert — see test_verifyMembershipRejectsExpiredConsensusState.
        SpectreMsgs.ConsensusState memory consensusState = _consensusState(_toNanos(block.timestamp));

        vm.expectRevert(
            abi.encodeWithSelector(
                SpectreClientErrors.LengthIsOutOfRange.selector, uint256(0), uint256(1), uint256(type(uint32).max)
            )
        );
        _deployWithConsensusState(consensusState, clientState);
    }

    function _deploy(uint128 consensusTimestamp) private returns (SpectreClient) {
        return _deployWithClientState(consensusTimestamp, _clientState());
    }

    function _deployWithClientState(
        uint128 consensusTimestamp,
        SpectreMsgs.ClientState memory clientState
    )
        private
        returns (SpectreClient)
    {
        return _deployWithConsensusState(_consensusState(consensusTimestamp), clientState);
    }

    function _deployWithConsensusState(
        SpectreMsgs.ConsensusState memory consensusState,
        SpectreMsgs.ClientState memory clientState
    )
        private
        returns (SpectreClient)
    {
        return new SpectreClient(
            address(0),
            address(membership),
            address(0),
            abi.encode(clientState),
            consensusState,
            _pinnedValidatorSet(),
            address(0)
        );
    }

    function _clientState() private pure returns (SpectreMsgs.ClientState memory) {
        return SpectreMsgs.ClientState({
            chainId: "test-chain-0",
            trustLevel: SpectreMsgs.TrustThreshold({ numerator: 1, denominator: 3 }),
            latestHeight: ICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: HEIGHT }),
            trustingPeriod: TRUSTING_PERIOD,
            unbondingPeriod: UNBONDING_PERIOD,
            isFrozen: false,
            clockDrift: CLOCK_DRIFT
        });
    }

    function _membershipMsg(uint128 consensusTimestamp)
        private
        pure
        returns (LightClientMsgs.MsgVerifyMembership memory msg_)
    {
        bytes memory value = hex"01";
        msg_ = LightClientMsgs.MsgVerifyMembership({
            height: _height(),
            kvPairs: _kvPairs(value),
            merkleProofs: new MembershipMsgs.MerkleProof[](0),
            appHash: APP_HASH,
            trustedConsensusState: _consensusState(consensusTimestamp),
            membershipType: MembershipMsgs.MembershipType.Membership,
            path: _path(),
            value: value
        });
    }

    function _nonMembershipMsg(uint128 consensusTimestamp)
        private
        pure
        returns (LightClientMsgs.MsgVerifyNonMembership memory msg_)
    {
        msg_ = LightClientMsgs.MsgVerifyNonMembership({
            height: _height(),
            kvPairs: _kvPairs(bytes("")),
            merkleProofs: new MembershipMsgs.MerkleProof[](0),
            appHash: APP_HASH,
            trustedConsensusState: _consensusState(consensusTimestamp),
            membershipType: MembershipMsgs.MembershipType.Membership,
            path: _path()
        });
    }

    function _pinnedValidatorSet() private pure returns (SpectreMsgs.ValidatorSet memory vs) {
        SpectreMsgs.ValidatorInfo[] memory vals = new SpectreMsgs.ValidatorInfo[](1);
        vals[0] = SpectreMsgs.ValidatorInfo({
            valAddress: bytes("validator"), pubKey: bytes32(uint256(1)), votingPower: 100, proposerPriority: 0
        });
        vs = SpectreMsgs.ValidatorSet({
            validators: vals, hasProposer: false, proposer: vals[0], totalVotingPower: 100
        });
    }

    function _consensusState(uint128 timestamp) private pure returns (SpectreMsgs.ConsensusState memory) {
        // nextValidatorsHash must match _pinnedValidatorSet()'s hash: the constructor now asserts
        // Header.hashValSet(initialPinnedValidatorSet) == consensusState.nextValidatorsHash (LC-03).
        return SpectreMsgs.ConsensusState({
            timestamp: timestamp, root: APP_HASH, nextValidatorsHash: Header.hashValSet(_pinnedValidatorSet())
        });
    }

    function _kvPairs(bytes memory value) private pure returns (MembershipMsgs.KVPair[] memory kvPairs) {
        kvPairs = new MembershipMsgs.KVPair[](1);
        kvPairs[0] = MembershipMsgs.KVPair({ path: _path(), value: value });
    }

    function _path() private pure returns (bytes[] memory path) {
        path = new bytes[](1);
        path[0] = bytes("key");
    }

    function _height() private pure returns (ICS02ClientMsgs.Height memory) {
        return ICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: HEIGHT });
    }

    function _toNanos(uint256 seconds_) private pure returns (uint128) {
        // forge-lint: disable-next-line(unsafe-typecast)
        return uint128(seconds_ * 1_000_000_000);
    }
}
