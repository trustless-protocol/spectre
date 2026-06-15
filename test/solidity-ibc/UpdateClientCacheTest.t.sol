// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test, console } from "forge-std/Test.sol";

import { Groth16ICS07Tendermint } from "../../contracts/light-clients/Groth16ICS07Tendermint.sol";
import { UpdateClient } from "../../contracts/programs/UpdateClient.sol";
import { WrapperVerifier } from "../../contracts/utils/WrapperVerifier.sol";
import { Header } from "../../contracts/utils/Header.sol";
import { IUpdateClientMsgs } from "../../contracts/light-clients/msgs/IUpdateClientMsgs.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../../contracts/msgs/IICS02ClientMsgs.sol";
import { ILightClientMsgs } from "../../contracts/msgs/ILightClientMsgs.sol";
import { IGroth16ICS07TendermintErrors } from "../../contracts/light-clients/errors/IGroth16ICS07TendermintErrors.sol";

contract CacheAlwaysTrueVerifier {
    function verifyProof(bytes calldata, uint256[2] calldata) external pure returns (bool) {
        return true;
    }
}

contract UpdateClientCacheTest is Test {
    WrapperVerifier internal wrapper;
    UpdateClient internal updateClientImpl;
    CacheAlwaysTrueVerifier internal stubBucket;

    address constant STUB_MEMBERSHIP = address(0xBABE);
    address constant STUB_MISBEHAVIOUR = address(0xBEEF);

    string constant CHAIN_ID = "cosmoshub-0";
    uint64 constant TRUSTED_HEIGHT = 1000;
    uint64 constant HEIGHT_1001 = 1001;
    uint64 constant HEIGHT_1002 = 1002;
    uint128 constant TRUSTED_TS_NS = 1_700_000_000 * 1e9;
    uint128 constant TS_1001_NS = 1_700_000_010 * 1e9;
    uint128 constant TS_1002_NS = 1_700_000_020 * 1e9;
    uint32 constant TRUSTING_PERIOD = 14 days;
    uint32 constant UNBONDING_PERIOD = 21 days;

    struct BucketConfig {
        uint16 bucket;
        uint16 valCount;
        uint16 activeCount;
    }

    function setUp() public {
        vm.warp(1_700_000_030);

        wrapper = new WrapperVerifier(address(this));
        stubBucket = new CacheAlwaysTrueVerifier();
        updateClientImpl = new UpdateClient();

        uint16[6] memory bs = [uint16(4), 8, 16, 32, 64, 128];
        for (uint256 i = 0; i < bs.length; i++) {
            wrapper.setBucket(bs[i], address(stubBucket), CacheAlwaysTrueVerifier.verifyProof.selector);
        }
    }

    function test_updateClient_revert_whenCurrentSetIsEmptyAndCacheMisses() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        bytes32 hashA = Header.hashValSet(valA);

        IICS07TendermintMsgs.Header memory header =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, _emptyValidatorSet(), hashA, TS_1001_NS, cfg.activeCount);
        IICS07TendermintMsgs.ConsensusState memory trustedCS =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));

        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ =
            _buildMsg(_clientState(), trustedCS, header, cfg.bucket, cfg.activeCount);
        msg_.proposedHeader.validatorSet = _emptyValidatorSet();

        vm.expectRevert(abi.encodeWithSelector(IGroth16ICS07TendermintErrors.ValidatorSetCacheMiss.selector, hashA));
        ics07.updateClient(abi.encode(msg_));
    }

    function test_updateClient_reverts_whenValidatorSetExceedsLimit() public {
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(181, 0);
        bytes32 hashA = Header.hashValSet(valA);

        IICS07TendermintMsgs.Header memory header =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashA, TS_1001_NS, 121);
        IICS07TendermintMsgs.ConsensusState memory trustedCS =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));

        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ = _buildMsg(_clientState(), trustedCS, header, 128, 121);

        vm.expectRevert(
            abi.encodeWithSelector(
                IGroth16ICS07TendermintErrors.ValidatorCountExceedsLimit.selector, uint256(181), uint256(180)
            )
        );
        ics07.updateClient(abi.encode(msg_));
    }

    function test_updateClient_reverts_whenHeightIsNonMonotonic() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        bytes32 hashA = Header.hashValSet(valA);

        IICS07TendermintMsgs.ConsensusState memory trustedCS =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));

        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS);

        // First update: advance latestHeight from 1000 to 1002.
        IICS07TendermintMsgs.Header memory header1002 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1002, valA, valA, hashA, TS_1002_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1002 =
            _buildMsg(_clientState(), trustedCS, header1002, cfg.bucket, cfg.activeCount);
        ics07.updateClient(abi.encode(msg1002));

        // Second update: try height 1001 < latestHeight 1002.
        IICS07TendermintMsgs.Header memory header1001 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashA, TS_1001_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1001 =
            _buildMsg(_clientState(), trustedCS, header1001, cfg.bucket, cfg.activeCount);

        vm.expectRevert(
            abi.encodeWithSelector(
                IGroth16ICS07TendermintErrors.NonMonotonicHeightUpdate.selector,
                uint64(HEIGHT_1002),
                uint64(HEIGHT_1001)
            )
        );
        ics07.updateClient(abi.encode(msg1001));
    }

    function test_updateClient_reverts_whenClockDriftDoesNotMatchStoredClientState() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        bytes32 hashA = Header.hashValSet(valA);

        IICS07TendermintMsgs.Header memory header =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashA, TS_1001_NS, cfg.activeCount);
        IICS07TendermintMsgs.ConsensusState memory trustedCS =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));

        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ =
            _buildMsg(_clientState(), trustedCS, header, cfg.bucket, cfg.activeCount);
        msg_.clientState.clockDrift = 3600;

        vm.expectRevert(
            abi.encodeWithSelector(
                IGroth16ICS07TendermintErrors.ClockDriftMismatch.selector, uint256(1800), uint256(3600)
            )
        );
        ics07.updateClient(abi.encode(msg_));
    }

    function test_updateClient_reverts_whenProofTimeExceedsStoredClockDrift() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        bytes32 hashA = Header.hashValSet(valA);

        IICS07TendermintMsgs.Header memory header =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashA, TS_1001_NS, cfg.activeCount);
        IICS07TendermintMsgs.ConsensusState memory trustedCS =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));

        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ =
            _buildMsg(_clientState(), trustedCS, header, cfg.bucket, cfg.activeCount);

        uint256 proofTimestamp = uint256(TS_1002_NS / 1e9);
        vm.warp(proofTimestamp + 1801);

        vm.expectRevert(
            abi.encodeWithSelector(
                IGroth16ICS07TendermintErrors.ProofIsTooOld.selector, proofTimestamp + 1801, proofTimestamp
            )
        );
        ics07.updateClient(abi.encode(msg_));
    }

    function test_updateClient_populatesCache_andUsesItOnEmptyAdjacentNoOp() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        bytes32 hashA = Header.hashValSet(valA);

        IICS07TendermintMsgs.Header memory header =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashA, TS_1001_NS, cfg.activeCount);
        IICS07TendermintMsgs.ConsensusState memory trustedCS =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));

        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS);

        IUpdateClientMsgs.MsgUpdateClient memory fullMsg =
            _buildMsg(_clientState(), trustedCS, header, cfg.bucket, cfg.activeCount);

        ILightClientMsgs.UpdateResult first = ics07.updateClient(abi.encode(fullMsg));
        assertEq(uint8(first), uint8(ILightClientMsgs.UpdateResult.Update), "first update should succeed");
        assertTrue(_hasCachedValidatorSet(ics07, hashA), "validator set should be cached");

        IUpdateClientMsgs.MsgUpdateClient memory cacheMsg =
            _buildMsg(_clientState(), trustedCS, header, cfg.bucket, cfg.activeCount);
        cacheMsg.proposedHeader.validatorSet = _emptyValidatorSet();
        cacheMsg.proposedHeader.trustedNextValidatorSet = _emptyValidatorSet();

        ILightClientMsgs.UpdateResult second = ics07.updateClient(abi.encode(cacheMsg));
        assertEq(uint8(second), uint8(ILightClientMsgs.UpdateResult.NoOp), "cache-hit replay should NoOp");
    }

    function test_updateClient_cacheHit_allowsDifferentSignerSubsetForSameValidatorSet() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        bytes32 hashA = Header.hashValSet(valA);

        IICS07TendermintMsgs.ConsensusState memory trustedCS0 =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS0);

        IICS07TendermintMsgs.Header memory header1001 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashA, TS_1001_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1001 =
            _buildMsg(_clientState(), trustedCS0, header1001, cfg.bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(msg1001))), uint8(ILightClientMsgs.UpdateResult.Update));

        IICS07TendermintMsgs.ConsensusState memory trustedCS1001 =
            _consensusState(TS_1001_NS, hashA, header1001.signedHeader.header.appHash);
        IICS07TendermintMsgs.Header memory header1002 =
            _buildHeader(HEIGHT_1001, HEIGHT_1002, valA, _emptyValidatorSet(), hashA, TS_1002_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory cacheMsg =
            _buildMsg(_clientState(), trustedCS1001, header1002, cfg.bucket, cfg.activeCount);
        _setSignerRange(cacheMsg, valA, 6, cfg.activeCount);
        cacheMsg.proposedHeader.validatorSet = _emptyValidatorSet();
        cacheMsg.proposedHeader.trustedNextValidatorSet = _emptyValidatorSet();

        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(cacheMsg));
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.Update), "cache should cover all validator slots");
    }

    function test_updateClient_adjacent_skipsTrustedNextSetEvenBeforeCacheExists() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        bytes32 hashA = Header.hashValSet(valA);

        IICS07TendermintMsgs.Header memory header =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, _emptyValidatorSet(), hashA, TS_1001_NS, cfg.activeCount);
        IICS07TendermintMsgs.ConsensusState memory trustedCS =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));

        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ =
            _buildMsg(_clientState(), trustedCS, header, cfg.bucket, cfg.activeCount);
        msg_.proposedHeader.trustedNextValidatorSet = _emptyValidatorSet();

        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(msg_));
        assertEq(
            uint8(result),
            uint8(ILightClientMsgs.UpdateResult.Update),
            "adjacent update should not need trusted next set"
        );
        assertTrue(_hasCachedValidatorSet(ics07, hashA), "current validator set should still be cached");
    }

    function test_updateClient_nonAdjacent_usesCachedCurrentAndSuppliedTrustedNextSet() public {
        BucketConfig memory cfg = _cfg(16);

        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        IICS07TendermintMsgs.ValidatorSet memory valB = _mutateValidatorSet(valA);
        bytes32 hashA = Header.hashValSet(valA);
        bytes32 hashB = Header.hashValSet(valB);

        IICS07TendermintMsgs.ConsensusState memory trustedCS0 =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));

        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS0);

        IICS07TendermintMsgs.Header memory header1001 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashB, TS_1001_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1001 =
            _buildMsg(_clientState(), trustedCS0, header1001, cfg.bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(msg1001))), uint8(ILightClientMsgs.UpdateResult.Update));
        assertTrue(_hasCachedValidatorSet(ics07, hashA), "validator set A should be cached");

        IICS07TendermintMsgs.ConsensusState memory trustedCS1001 =
            _consensusState(TS_1001_NS, hashB, header1001.signedHeader.header.appHash);
        IICS07TendermintMsgs.Header memory header1002 =
            _buildHeader(HEIGHT_1001, HEIGHT_1002, valB, valB, hashB, TS_1002_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1002 =
            _buildMsg(_clientState(), trustedCS1001, header1002, cfg.bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(msg1002))), uint8(ILightClientMsgs.UpdateResult.Update));
        assertTrue(_hasCachedValidatorSet(ics07, hashB), "validator set B should be cached");

        IICS07TendermintMsgs.Header memory nonAdjacentHeader =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1002, valB, valA, hashB, TS_1002_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory cacheMsg =
            _buildMsg(_clientState(), trustedCS0, nonAdjacentHeader, cfg.bucket, cfg.activeCount);
        cacheMsg.proposedHeader.validatorSet = _emptyValidatorSet();
        cacheMsg.proposedHeader.trustedNextValidatorSet = valA;

        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(cacheMsg));
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.NoOp), "non-adjacent cache-hit replay should NoOp");
    }

    function test_updateClient_nonAdjacent_sameTrustedNextHashUsesCachedValidation() public {
        BucketConfig memory cfg = _cfg(16);

        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        bytes32 hashA = Header.hashValSet(valA);

        IICS07TendermintMsgs.ConsensusState memory trustedCS0 =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS0);

        IICS07TendermintMsgs.Header memory header1001 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashA, TS_1001_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1001 =
            _buildMsg(_clientState(), trustedCS0, header1001, cfg.bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(msg1001))), uint8(ILightClientMsgs.UpdateResult.Update));
        assertTrue(_hasCachedValidatorSet(ics07, hashA), "validator set A should be cached");

        IICS07TendermintMsgs.Header memory nonAdjacentHeader =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1002, valA, valA, hashA, TS_1002_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory cacheMsg =
            _buildMsg(_clientState(), trustedCS0, nonAdjacentHeader, cfg.bucket, cfg.activeCount);
        cacheMsg.proposedHeader.validatorSet = _emptyValidatorSet();

        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(cacheMsg));
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.Update), "same-hash cached trusted next");
    }

    function test_updateClient_nonAdjacent_sameTrustedNextHashRejectsMismatchedTrustedNextSet() public {
        BucketConfig memory cfg = _cfg(16);

        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        bytes32 hashA = Header.hashValSet(valA);
        IICS07TendermintMsgs.ValidatorSet memory badTrustedNext = _buildValSet(cfg.valCount, 0);
        badTrustedNext.validators[0].votingPower = 101;
        badTrustedNext.totalVotingPower = valA.totalVotingPower + 1;

        IICS07TendermintMsgs.ConsensusState memory trustedCS0 =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS0);

        IICS07TendermintMsgs.Header memory header1001 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashA, TS_1001_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1001 =
            _buildMsg(_clientState(), trustedCS0, header1001, cfg.bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(msg1001))), uint8(ILightClientMsgs.UpdateResult.Update));

        IICS07TendermintMsgs.Header memory nonAdjacentHeader =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1002, valA, badTrustedNext, hashA, TS_1002_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory cacheMsg =
            _buildMsg(_clientState(), trustedCS0, nonAdjacentHeader, cfg.bucket, cfg.activeCount);
        cacheMsg.proposedHeader.validatorSet = _emptyValidatorSet();

        vm.expectRevert(
            abi.encodeWithSelector(IGroth16ICS07TendermintErrors.MismatchedValidatorHashes.selector, hashA, bytes32(0))
        );
        ics07.updateClient(abi.encode(cacheMsg));
    }

    function test_updateClient_deltaCache_singleVotingPowerChange() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        IICS07TendermintMsgs.ValidatorSet memory valB = _changeVotingPower(valA, 0, 110);
        bytes32 hashA = Header.hashValSet(valA);
        bytes32 hashB = Header.hashValSet(valB);

        IICS07TendermintMsgs.ConsensusState memory trustedCS0 =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS0);

        IICS07TendermintMsgs.Header memory header1001 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashB, TS_1001_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1001 =
            _buildMsg(_clientState(), trustedCS0, header1001, cfg.bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(msg1001))), uint8(ILightClientMsgs.UpdateResult.Update));
        assertTrue(_hasCachedValidatorSet(ics07, hashA), "validator set A should be cached");

        IICS07TendermintMsgs.ConsensusState memory trustedCS1001 =
            _consensusState(TS_1001_NS, hashB, header1001.signedHeader.header.appHash);
        IICS07TendermintMsgs.Header memory header1002 =
            _buildHeader(HEIGHT_1001, HEIGHT_1002, valB, _emptyValidatorSet(), hashB, TS_1002_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory deltaMsg =
            _buildMsg(_clientState(), trustedCS1001, header1002, cfg.bucket, cfg.activeCount);
        deltaMsg.proposedHeader.validatorSet = _emptyValidatorSet();
        deltaMsg.proposedHeader.trustedNextValidatorSet = _emptyValidatorSet();
        uint32[16] memory deltaIndices;
        bytes32[16] memory deltaPubKeys;
        uint64[16] memory deltaVotingPowers;
        deltaIndices[0] = 0;
        deltaPubKeys[0] = valB.validators[0].pubKey;
        deltaVotingPowers[0] = 110;
        deltaMsg.currentValidatorSetDelta = IUpdateClientMsgs.ValidatorSetDelta({
            baseValidatorsHash: hashA,
            leafCount: 1,
            indices: deltaIndices,
            pubKeys: deltaPubKeys,
            votingPowers: deltaVotingPowers
        });

        uint256 g0 = gasleft();
        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(deltaMsg));
        uint256 used = g0 - gasleft();
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.Update), "delta update should succeed");
        assertTrue(_hasCachedValidatorSet(ics07, hashB), "delta-derived validator set B should be cached");
        console.log("bucket=16 delta-cache adjacent update gas=", used);

        (uint32[] memory indices, bytes32[] memory pubkeys, uint64[] memory votingPowers) =
            ics07.getCachedValidatorSet(hashB);
        assertEq(indices.length, valB.validators.length, "cached B length");
        assertEq(indices[0], 0, "cached changed index");
        assertEq(pubkeys[0], valB.validators[0].pubKey, "cached changed pubkey");
        assertEq(votingPowers[0], 110, "cached changed voting power");
        assertEq(pubkeys[1], valB.validators[1].pubKey, "cached inherited pubkey");
        assertEq(votingPowers[1], valB.validators[1].votingPower, "cached inherited voting power");
    }

    function test_gas_n16_deltaCache_singleVotingPowerChange_nonAdjacent() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        IICS07TendermintMsgs.ValidatorSet memory valB = _changeVotingPower(valA, 0, 110);
        bytes32 hashA = Header.hashValSet(valA);
        bytes32 hashB = Header.hashValSet(valB);

        IICS07TendermintMsgs.ConsensusState memory trustedCS0 =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS0);

        IICS07TendermintMsgs.Header memory header1001 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashB, TS_1001_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1001 =
            _buildMsg(_clientState(), trustedCS0, header1001, cfg.bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(msg1001))), uint8(ILightClientMsgs.UpdateResult.Update));

        IICS07TendermintMsgs.Header memory header1002 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1002, valB, valA, hashB, TS_1002_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory deltaMsg =
            _buildMsg(_clientState(), trustedCS0, header1002, cfg.bucket, cfg.activeCount);
        deltaMsg.proposedHeader.validatorSet = _emptyValidatorSet();

        uint32[16] memory deltaIndices;
        bytes32[16] memory deltaPubKeys;
        uint64[16] memory deltaVotingPowers;
        deltaIndices[0] = 0;
        deltaPubKeys[0] = valB.validators[0].pubKey;
        deltaVotingPowers[0] = 110;
        deltaMsg.currentValidatorSetDelta = IUpdateClientMsgs.ValidatorSetDelta({
            baseValidatorsHash: hashA,
            leafCount: 1,
            indices: deltaIndices,
            pubKeys: deltaPubKeys,
            votingPowers: deltaVotingPowers
        });

        uint256 g0 = gasleft();
        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(deltaMsg));
        uint256 used = g0 - gasleft();
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.Update), "non-adjacent delta update should succeed");
        console.log("bucket=16 delta-cache non-adjacent one-leaf update gas=", used);
    }

    function test_updateClient_deltaCache_twoLeafSwapAndPowerChange() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        IICS07TendermintMsgs.ValidatorSet memory valB = _swapAdjacentAndChangePower(valA, 7, 110);
        bytes32 hashA = Header.hashValSet(valA);
        bytes32 hashB = Header.hashValSet(valB);

        IICS07TendermintMsgs.ConsensusState memory trustedCS0 =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS0);

        IICS07TendermintMsgs.Header memory header1001 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashB, TS_1001_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1001 =
            _buildMsg(_clientState(), trustedCS0, header1001, cfg.bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(msg1001))), uint8(ILightClientMsgs.UpdateResult.Update));

        IICS07TendermintMsgs.ConsensusState memory trustedCS1001 =
            _consensusState(TS_1001_NS, hashB, header1001.signedHeader.header.appHash);
        IICS07TendermintMsgs.Header memory header1002 =
            _buildHeader(HEIGHT_1001, HEIGHT_1002, valB, _emptyValidatorSet(), hashB, TS_1002_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory deltaMsg =
            _buildMsg(_clientState(), trustedCS1001, header1002, cfg.bucket, cfg.activeCount);
        deltaMsg.proposedHeader.validatorSet = _emptyValidatorSet();
        deltaMsg.proposedHeader.trustedNextValidatorSet = _emptyValidatorSet();

        uint32[16] memory deltaIndices;
        bytes32[16] memory deltaPubKeys;
        uint64[16] memory deltaVotingPowers;
        deltaIndices[0] = 7;
        deltaPubKeys[0] = valB.validators[7].pubKey;
        deltaVotingPowers[0] = valB.validators[7].votingPower;
        deltaIndices[1] = 8;
        deltaPubKeys[1] = valB.validators[8].pubKey;
        deltaVotingPowers[1] = valB.validators[8].votingPower;
        deltaMsg.currentValidatorSetDelta = IUpdateClientMsgs.ValidatorSetDelta({
            baseValidatorsHash: hashA,
            leafCount: 2,
            indices: deltaIndices,
            pubKeys: deltaPubKeys,
            votingPowers: deltaVotingPowers
        });

        uint256 g0 = gasleft();
        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(deltaMsg));
        uint256 used = g0 - gasleft();
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.Update), "multi-leaf delta should succeed");
        assertTrue(_hasCachedValidatorSet(ics07, hashB), "multi-leaf delta-derived validator set should be cached");
        console.log("bucket=16 delta-cache two-leaf update gas=", used);

        (uint32[] memory indices, bytes32[] memory pubkeys, uint64[] memory votingPowers) =
            ics07.getCachedValidatorSet(hashB);
        assertEq(indices.length, valB.validators.length, "cached B length");
        assertEq(pubkeys[7], valB.validators[7].pubKey, "cached swapped pubkey at 7");
        assertEq(votingPowers[7], valB.validators[7].votingPower, "cached changed power at 7");
        assertEq(pubkeys[8], valB.validators[8].pubKey, "cached swapped pubkey at 8");
        assertEq(votingPowers[8], valB.validators[8].votingPower, "cached inherited power at 8");

        IICS07TendermintMsgs.ConsensusState memory trustedCS1002 =
            _consensusState(TS_1002_NS, hashB, header1002.signedHeader.header.appHash);
        IICS07TendermintMsgs.Header memory header1003 = _buildHeader(
            HEIGHT_1002, HEIGHT_1002 + 1, valB, _emptyValidatorSet(), hashB, TS_1002_NS + 10, cfg.activeCount
        );
        IUpdateClientMsgs.MsgUpdateClient memory postDeltaCacheMsg =
            _buildMsg(_clientState(), trustedCS1002, header1003, cfg.bucket, cfg.activeCount);
        postDeltaCacheMsg.proposedHeader.validatorSet = _emptyValidatorSet();
        postDeltaCacheMsg.proposedHeader.trustedNextValidatorSet = _emptyValidatorSet();

        g0 = gasleft();
        result = ics07.updateClient(abi.encode(postDeltaCacheMsg));
        used = g0 - gasleft();
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.Update), "post-delta cache hit should succeed");
        console.log("bucket=16 post-delta cache-hit adjacent update gas=", used);
    }

    function test_updateClient_deltaCache_consecutiveSiblingVotingPowerChanges() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        IICS07TendermintMsgs.ValidatorSet memory valB = _changeVotingPower(valA, 0, 110);
        IICS07TendermintMsgs.ValidatorSet memory valC = _changeVotingPower(valB, 1, 90);
        bytes32 hashA = Header.hashValSet(valA);
        bytes32 hashB = Header.hashValSet(valB);
        bytes32 hashC = Header.hashValSet(valC);

        IICS07TendermintMsgs.ConsensusState memory trustedCS0 =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS0);

        IICS07TendermintMsgs.Header memory header1001 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashB, TS_1001_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1001 =
            _buildMsg(_clientState(), trustedCS0, header1001, cfg.bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(msg1001))), uint8(ILightClientMsgs.UpdateResult.Update));

        IICS07TendermintMsgs.ConsensusState memory trustedCS1001 =
            _consensusState(TS_1001_NS, hashB, header1001.signedHeader.header.appHash);
        IICS07TendermintMsgs.Header memory header1002 =
            _buildHeader(HEIGHT_1001, HEIGHT_1002, valB, _emptyValidatorSet(), hashC, TS_1002_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory deltaMsgB =
            _buildMsg(_clientState(), trustedCS1001, header1002, cfg.bucket, cfg.activeCount);
        deltaMsgB.proposedHeader.validatorSet = _emptyValidatorSet();
        deltaMsgB.proposedHeader.trustedNextValidatorSet = _emptyValidatorSet();
        deltaMsgB.currentValidatorSetDelta = _singleLeafDelta(hashA, valB, 0);
        assertEq(uint8(ics07.updateClient(abi.encode(deltaMsgB))), uint8(ILightClientMsgs.UpdateResult.Update));

        IICS07TendermintMsgs.ConsensusState memory trustedCS1002 =
            _consensusState(TS_1002_NS, hashC, header1002.signedHeader.header.appHash);
        IICS07TendermintMsgs.Header memory header1003 = _buildHeader(
            HEIGHT_1002, HEIGHT_1002 + 1, valC, _emptyValidatorSet(), hashC, TS_1002_NS + 10, cfg.activeCount
        );
        IUpdateClientMsgs.MsgUpdateClient memory deltaMsgC =
            _buildMsg(_clientState(), trustedCS1002, header1003, cfg.bucket, cfg.activeCount);
        deltaMsgC.proposedHeader.validatorSet = _emptyValidatorSet();
        deltaMsgC.proposedHeader.trustedNextValidatorSet = _emptyValidatorSet();
        deltaMsgC.currentValidatorSetDelta = _singleLeafDelta(hashB, valC, 1);

        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(deltaMsgC));
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.Update), "second sibling delta should succeed");

        (uint32[] memory indices, bytes32[] memory pubkeys, uint64[] memory votingPowers) =
            ics07.getCachedValidatorSet(hashC);
        assertEq(indices.length, valC.validators.length, "cached C length");
        assertEq(pubkeys[0], valC.validators[0].pubKey, "cached first sibling pubkey");
        assertEq(votingPowers[0], valC.validators[0].votingPower, "cached first sibling power");
        assertEq(pubkeys[1], valC.validators[1].pubKey, "cached second sibling pubkey");
        assertEq(votingPowers[1], valC.validators[1].votingPower, "cached second sibling power");
    }

    function test_gas_n16_cacheHit_emptySets() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        bytes32 hashA = Header.hashValSet(valA);

        IICS07TendermintMsgs.Header memory header =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashA, TS_1001_NS, cfg.activeCount);
        IICS07TendermintMsgs.ConsensusState memory trustedCS =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));

        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS);

        IUpdateClientMsgs.MsgUpdateClient memory fullMsg =
            _buildMsg(_clientState(), trustedCS, header, cfg.bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(fullMsg))), uint8(ILightClientMsgs.UpdateResult.Update));

        IUpdateClientMsgs.MsgUpdateClient memory cacheMsg =
            _buildMsg(_clientState(), trustedCS, header, cfg.bucket, cfg.activeCount);
        cacheMsg.proposedHeader.validatorSet = _emptyValidatorSet();
        cacheMsg.proposedHeader.trustedNextValidatorSet = _emptyValidatorSet();

        uint256 g0 = gasleft();
        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(cacheMsg));
        uint256 used = g0 - gasleft();
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.NoOp));
        console.log("bucket=16 cache-hit empty-sets gas=", used);
    }

    function test_gas_n16_cacheHit_adjacentUpdate() public {
        BucketConfig memory cfg = _cfg(16);
        IICS07TendermintMsgs.ValidatorSet memory valA = _buildValSet(cfg.valCount, 0);
        bytes32 hashA = Header.hashValSet(valA);

        IICS07TendermintMsgs.ConsensusState memory trustedCS0 =
            _consensusState(TRUSTED_TS_NS, hashA, bytes32(uint256(0xAAA1)));
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS0);

        IICS07TendermintMsgs.Header memory header1001 =
            _buildHeader(TRUSTED_HEIGHT, HEIGHT_1001, valA, valA, hashA, TS_1001_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1001 =
            _buildMsg(_clientState(), trustedCS0, header1001, cfg.bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(msg1001))), uint8(ILightClientMsgs.UpdateResult.Update));

        IICS07TendermintMsgs.ConsensusState memory trustedCS1001 =
            _consensusState(TS_1001_NS, hashA, header1001.signedHeader.header.appHash);
        IICS07TendermintMsgs.Header memory header1002 =
            _buildHeader(HEIGHT_1001, HEIGHT_1002, valA, _emptyValidatorSet(), hashA, TS_1002_NS, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory cacheMsg =
            _buildMsg(_clientState(), trustedCS1001, header1002, cfg.bucket, cfg.activeCount);
        cacheMsg.proposedHeader.validatorSet = _emptyValidatorSet();
        cacheMsg.proposedHeader.trustedNextValidatorSet = _emptyValidatorSet();

        uint256 g0 = gasleft();
        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(cacheMsg));
        uint256 used = g0 - gasleft();
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.Update));
        console.log("bucket=16 cache-hit adjacent update gas=", used);
    }

    function _deployLightClient(IICS07TendermintMsgs.ConsensusState memory trustedCS)
        internal
        returns (Groth16ICS07Tendermint)
    {
        return new Groth16ICS07Tendermint(
            address(wrapper),
            STUB_MEMBERSHIP,
            STUB_MISBEHAVIOUR,
            address(updateClientImpl),
            abi.encode(_clientState()),
            keccak256(abi.encode(trustedCS)),
            address(0)
        );
    }

    function _hasCachedValidatorSet(Groth16ICS07Tendermint ics07, bytes32 validatorsHash) internal view returns (bool) {
        (uint32[] memory indices,,) = ics07.getCachedValidatorSet(validatorsHash);
        return indices.length != 0;
    }

    function _buildMsg(
        IICS07TendermintMsgs.ClientState memory cs,
        IICS07TendermintMsgs.ConsensusState memory trustedCS,
        IICS07TendermintMsgs.Header memory header,
        uint16 bucket,
        uint16 activeCount
    )
        internal
        pure
        returns (IUpdateClientMsgs.MsgUpdateClient memory msg_)
    {
        uint32[] memory idx = new uint32[](bucket);
        bytes32[] memory pks = new bytes32[](bucket);
        uint64[] memory tsS = new uint64[](bucket);
        uint32[] memory tsN = new uint32[](bucket);
        bool[] memory act = new bool[](bucket);
        uint32[] memory trustedOverlap = new uint32[](bucket);
        for (uint256 i = 0; i < bucket; i++) {
            trustedOverlap[i] = type(uint32).max;
            if (i < activeCount) {
                idx[i] = uint32(i);
                pks[i] = header.validatorSet.validators[i].pubKey;
                act[i] = true;
                trustedOverlap[i] = uint32(i);
            }
            tsS[i] = uint64(1_700_000_000 + i);
            tsN[i] = uint32(i * 1_000_000);
        }

        msg_.clientState = cs;
        msg_.trustedConsensusState = trustedCS;
        msg_.proposedHeader = header;
        msg_.time = TS_1002_NS;
        msg_.proof = [uint256(0), 0, 0, 0, 0, 0, 0, 0];
        msg_.commitments = [uint256(0), 0];
        msg_.commitmentPok = [uint256(0), 0];
        msg_.bucket = bucket;
        msg_.signerIndices = idx;
        msg_.signerPubkeys = pks;
        msg_.timestampSeconds = tsS;
        msg_.timestampNanos = tsN;
        msg_.active = act;
        msg_.trustedOverlapIndices = trustedOverlap;
    }

    function _setSignerRange(
        IUpdateClientMsgs.MsgUpdateClient memory msg_,
        IICS07TendermintMsgs.ValidatorSet memory validatorSet,
        uint32 start,
        uint16 activeCount
    )
        internal
        pure
    {
        for (uint256 i = 0; i < msg_.proposedHeader.signedHeader.commit.commitSigs.length; i++) {
            msg_.proposedHeader.signedHeader.commit.commitSigs[i] = IICS07TendermintMsgs.CommitSig({
                flag: IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_ABSENT,
                data: IICS07TendermintMsgs.CommitSigData({
                    validatorAddress: "", timestamp: 0, hasSignature: false, signature: ""
                })
            });
        }

        for (uint256 i = 0; i < msg_.signerIndices.length; i++) {
            if (i < activeCount) {
                uint32 idx = start + uint32(i);
                msg_.signerIndices[i] = idx;
                msg_.signerPubkeys[i] = validatorSet.validators[idx].pubKey;
                msg_.active[i] = true;
                msg_.proposedHeader.signedHeader.commit.commitSigs[idx] = IICS07TendermintMsgs.CommitSig({
                    flag: IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_COMMIT,
                    data: IICS07TendermintMsgs.CommitSigData({
                        validatorAddress: validatorSet.validators[idx].valAddress,
                        timestamp: msg_.proposedHeader.signedHeader.header.time,
                        hasSignature: false,
                        signature: ""
                    })
                });
                msg_.trustedOverlapIndices[i] = idx;
            } else {
                msg_.signerIndices[i] = 0;
                msg_.signerPubkeys[i] = bytes32(0);
                msg_.active[i] = false;
                msg_.trustedOverlapIndices[i] = type(uint32).max;
            }
        }
    }

    function _singleLeafDelta(
        bytes32 baseValidatorsHash,
        IICS07TendermintMsgs.ValidatorSet memory validatorSet,
        uint32 changedIndex
    )
        internal
        pure
        returns (IUpdateClientMsgs.ValidatorSetDelta memory delta)
    {
        uint32[16] memory deltaIndices;
        bytes32[16] memory deltaPubKeys;
        uint64[16] memory deltaVotingPowers;
        deltaIndices[0] = changedIndex;
        deltaPubKeys[0] = validatorSet.validators[changedIndex].pubKey;
        deltaVotingPowers[0] = validatorSet.validators[changedIndex].votingPower;
        delta = IUpdateClientMsgs.ValidatorSetDelta({
            baseValidatorsHash: baseValidatorsHash,
            leafCount: 1,
            indices: deltaIndices,
            pubKeys: deltaPubKeys,
            votingPowers: deltaVotingPowers
        });
    }

    function _buildHeader(
        uint64 trustedHeight,
        uint64 newHeight,
        IICS07TendermintMsgs.ValidatorSet memory currentValSet,
        IICS07TendermintMsgs.ValidatorSet memory trustedNextValSet,
        bytes32 nextValidatorsHash,
        uint128 headerTime,
        uint16 activeCount
    )
        internal
        pure
        returns (IICS07TendermintMsgs.Header memory header)
    {
        bytes32 currentValSetHash = Header.hashValSet(currentValSet);

        IICS07TendermintMsgs.BlockHeader memory bh;
        bh.chainId = CHAIN_ID;
        bh.height = newHeight;
        bh.time = headerTime;
        bh.appHash = bytes32(uint256(0xCC0000 + newHeight));
        bh.validatorsHash = currentValSetHash;
        bh.nextValidatorsHash = nextValidatorsHash;
        bytes32 headerHash = Header.hashHeader(bh);

        IICS07TendermintMsgs.BlockCommit memory bc = IICS07TendermintMsgs.BlockCommit({
            height: newHeight,
            round: 0,
            blockId: IICS07TendermintMsgs.BlockId({
                hashData: headerHash,
                partSetHeader: IICS07TendermintMsgs.PartSetHeader({ total: 1, hashData: bytes32(uint256(0x9A57)) })
            }),
            commitSigs: _buildCommitSigs(currentValSet, activeCount, headerTime)
        });

        header = IICS07TendermintMsgs.Header({
            signedHeader: IICS07TendermintMsgs.SignedHeader({ header: bh, commit: bc }),
            validatorSet: currentValSet,
            trustedHeight: IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: trustedHeight }),
            trustedNextValidatorSet: trustedNextValSet
        });
    }

    function _buildCommitSigs(
        IICS07TendermintMsgs.ValidatorSet memory vs,
        uint16 activeCount,
        uint128 headerTime
    )
        internal
        pure
        returns (IICS07TendermintMsgs.CommitSig[] memory sigs)
    {
        sigs = new IICS07TendermintMsgs.CommitSig[](vs.validators.length);
        for (uint256 i = 0; i < vs.validators.length; i++) {
            if (i < activeCount) {
                sigs[i] = IICS07TendermintMsgs.CommitSig({
                    flag: IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_COMMIT,
                    data: IICS07TendermintMsgs.CommitSigData({
                        validatorAddress: vs.validators[i].valAddress,
                        timestamp: headerTime,
                        hasSignature: false,
                        signature: ""
                    })
                });
            } else {
                sigs[i] = IICS07TendermintMsgs.CommitSig({
                    flag: IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_ABSENT,
                    data: IICS07TendermintMsgs.CommitSigData({
                        validatorAddress: "", timestamp: 0, hasSignature: false, signature: ""
                    })
                });
            }
        }
    }

    function _buildValSet(
        uint16 valCount,
        uint256 seed
    )
        internal
        pure
        returns (IICS07TendermintMsgs.ValidatorSet memory vs)
    {
        IICS07TendermintMsgs.ValidatorInfo[] memory vals = new IICS07TendermintMsgs.ValidatorInfo[](valCount);
        uint64 total = 0;
        for (uint256 i = 0; i < valCount; i++) {
            vals[i] = IICS07TendermintMsgs.ValidatorInfo({
                valAddress: abi.encodePacked(uint160(seed + i + 1)),
                pubKey: bytes32(uint256(0xA0000000 + seed + i + 1)),
                votingPower: 100,
                proposerPriority: 0
            });
            total += 100;
        }
        vs = IICS07TendermintMsgs.ValidatorSet({
            validators: vals, hasProposer: false, proposer: vals[0], totalVotingPower: total
        });
    }

    function _mutateValidatorSet(IICS07TendermintMsgs.ValidatorSet memory base)
        internal
        pure
        returns (IICS07TendermintMsgs.ValidatorSet memory mutated)
    {
        uint256 len = base.validators.length;
        IICS07TendermintMsgs.ValidatorInfo[] memory vals = new IICS07TendermintMsgs.ValidatorInfo[](len);
        uint64 total = 0;
        for (uint256 i = 0; i < len; i++) {
            vals[i] = base.validators[i];
        }

        vals[0].votingPower = 110;
        vals[1].votingPower = 90;

        for (uint256 i = 0; i < len; i++) {
            total += vals[i].votingPower;
        }

        mutated = IICS07TendermintMsgs.ValidatorSet({
            validators: vals, hasProposer: false, proposer: vals[0], totalVotingPower: total
        });
    }

    function _changeVotingPower(
        IICS07TendermintMsgs.ValidatorSet memory base,
        uint256 changedIndex,
        uint64 newVotingPower
    )
        internal
        pure
        returns (IICS07TendermintMsgs.ValidatorSet memory mutated)
    {
        uint256 len = base.validators.length;
        IICS07TendermintMsgs.ValidatorInfo[] memory vals = new IICS07TendermintMsgs.ValidatorInfo[](len);
        uint64 total = 0;
        for (uint256 i = 0; i < len; i++) {
            vals[i] = base.validators[i];
            if (i == changedIndex) {
                vals[i].votingPower = newVotingPower;
            }
            total += vals[i].votingPower;
        }

        mutated = IICS07TendermintMsgs.ValidatorSet({
            validators: vals, hasProposer: false, proposer: vals[0], totalVotingPower: total
        });
    }

    function _swapAdjacentAndChangePower(
        IICS07TendermintMsgs.ValidatorSet memory base,
        uint256 firstIndex,
        uint64 newFirstVotingPower
    )
        internal
        pure
        returns (IICS07TendermintMsgs.ValidatorSet memory mutated)
    {
        uint256 len = base.validators.length;
        IICS07TendermintMsgs.ValidatorInfo[] memory vals = new IICS07TendermintMsgs.ValidatorInfo[](len);
        for (uint256 i = 0; i < len; i++) {
            vals[i] = base.validators[i];
        }

        vals[firstIndex] = base.validators[firstIndex + 1];
        vals[firstIndex].votingPower = newFirstVotingPower;
        vals[firstIndex + 1] = base.validators[firstIndex];

        uint64 total = 0;
        for (uint256 i = 0; i < len; i++) {
            total += vals[i].votingPower;
        }

        mutated = IICS07TendermintMsgs.ValidatorSet({
            validators: vals, hasProposer: false, proposer: vals[0], totalVotingPower: total
        });
    }

    function _consensusState(
        uint128 timestamp,
        bytes32 nextValidatorsHash,
        bytes32 root
    )
        internal
        pure
        returns (IICS07TendermintMsgs.ConsensusState memory)
    {
        return IICS07TendermintMsgs.ConsensusState({
            timestamp: timestamp, root: root, nextValidatorsHash: nextValidatorsHash
        });
    }

    function _clientState() internal pure returns (IICS07TendermintMsgs.ClientState memory) {
        return IICS07TendermintMsgs.ClientState({
            chainId: CHAIN_ID,
            trustLevel: IICS07TendermintMsgs.TrustThreshold({ numerator: 1, denominator: 3 }),
            latestHeight: IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: TRUSTED_HEIGHT }),
            trustingPeriod: TRUSTING_PERIOD,
            unbondingPeriod: UNBONDING_PERIOD,
            isFrozen: false,
            zkAlgorithm: IICS07TendermintMsgs.SupportedZkAlgorithm.Groth16,
            clockDrift: 1800
        });
    }

    function _cfg(uint16 bucket) internal pure returns (BucketConfig memory) {
        if (bucket == 4) return BucketConfig(4, 4, 3);
        if (bucket == 8) return BucketConfig(8, 10, 7);
        if (bucket == 16) return BucketConfig(16, 20, 14);
        if (bucket == 32) return BucketConfig(32, 30, 21);
        if (bucket == 64) return BucketConfig(64, 60, 42);
        if (bucket == 128) return BucketConfig(128, 120, 84);
        revert("unknown bucket");
    }

    function _emptyValidatorSet() internal pure returns (IICS07TendermintMsgs.ValidatorSet memory validatorSet) {
        validatorSet.validators = new IICS07TendermintMsgs.ValidatorInfo[](0);
    }
}
