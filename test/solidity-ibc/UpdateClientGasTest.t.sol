// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test, console } from "forge-std/Test.sol";

import { Groth16ICS07Tendermint } from "../../contracts/light-clients/Groth16ICS07Tendermint.sol";
import { WrapperVerifier } from "../../contracts/utils/WrapperVerifier.sol";
import { UpdateClient } from "../../contracts/programs/UpdateClient.sol";
import { Header } from "../../contracts/utils/Header.sol";
import { IUpdateClientMsgs } from "../../contracts/light-clients/msgs/IUpdateClientMsgs.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../../contracts/msgs/IICS02ClientMsgs.sol";
import { ILightClientMsgs } from "../../contracts/msgs/ILightClientMsgs.sol";

/// @dev Stub bucket verifier — always returns true so gas measurement excludes
///      only the ~250K-gas BN254 pairing cost. Everything else (real
///      `UPDATE_CLIENT.updateClient`, `_verifyBatchAndQuorum`, `_hashWitness`,
///      `_voteSignBytes`) runs unmodified.
contract AlwaysTrueVerifier {
    function verifyProof(bytes calldata, uint256[2] calldata) external pure returns (bool) {
        return true;
    }
}

/// @notice Production-realistic gas benchmark for `Groth16ICS07Tendermint.updateClient`.
///         Builds a self-consistent MsgUpdateClient (validator-set hash matches
///         header field, header merkle root matches commit.blockId.hashData,
///         commitSigs cover the validator set) so the real `UPDATE_CLIENT.updateClient`
///         path executes end-to-end. Only the per-bucket Groth16 verifier is
///         stubbed — adding ~250K to each measured number recovers prod gas.
///
///         Run:    forge test --match-contract UpdateClientGasTest -vv
///         Diff:   compare measured value to BENCHMARK_2026-05-19_bucket16.md
///                 (n=16 → 2.79M ≈ measured + ~250K Groth16).
///
///         === Topology (bench-matched) ===
///         For each bucket, validator set + active count just clear the 2/3
///         quorum threshold, matching how the relayer picks signers in prod:
///         n=16 → 20 vals total, 14 active, 2 padding (3·14·100 > 2·20·100 ⇒ pass).
contract UpdateClientGasTest is Test {
    WrapperVerifier wrapper;
    UpdateClient updateClientImpl;
    AlwaysTrueVerifier stubBucket;

    address constant STUB_MEMBERSHIP   = address(0xBABE);
    address constant STUB_MISBEHAVIOUR = address(0xBEEF);

    // chainId revisionNumber must match latestHeight.revisionNumber.
    // "cosmoshub-0" → revisionNumber=0.
    string  constant CHAIN_ID         = "cosmoshub-0";
    uint64  constant TRUSTED_HEIGHT   = 1_000;
    uint64  constant NEW_HEIGHT       = 1_001;
    uint64  constant NEXT_HEIGHT      = 1_002;
    uint128 constant TRUSTED_TS_NS    = 1_700_000_000 * 1e9;
    uint128 constant NEW_TS_NS        = 1_700_000_010 * 1e9;
    uint128 constant NEXT_TS_NS       = 1_700_000_020 * 1e9;
    uint32  constant TRUSTING_PERIOD  = 14 days;
    uint32  constant UNBONDING_PERIOD = 21 days;

    struct BucketConfig {
        uint16 bucket;       // padded slot count (== per-slot array lengths)
        uint16 valCount;     // validators in proposedHeader.validatorSet
        uint16 activeCount;  // real signers; padding = bucket - activeCount
    }

    /// Bench-matched: active just clears 2/3 quorum of valCount.
    function _cfg(uint16 bucket) internal pure returns (BucketConfig memory) {
        if (bucket == 4)   return BucketConfig(4,   4,  3);
        if (bucket == 8)   return BucketConfig(8,  10,  7);
        if (bucket == 16)  return BucketConfig(16, 20, 14);
        if (bucket == 32)  return BucketConfig(32, 30, 21);
        if (bucket == 64)  return BucketConfig(64, 60, 42);
        if (bucket == 128) return BucketConfig(128,120, 84);
        revert("unknown bucket");
    }

    function setUp() public {
        // pin block.timestamp so the clock-drift check in
        // _validateClientStateAndTime accepts NEW_TS_NS.
        vm.warp(1_700_000_020);

        wrapper = new WrapperVerifier(address(this));
        stubBucket = new AlwaysTrueVerifier();
        updateClientImpl = new UpdateClient();

        uint16[6] memory bs = [uint16(4), 8, 16, 32, 64, 128];
        for (uint256 i = 0; i < bs.length; i++) {
            wrapper.setBucket(bs[i], address(stubBucket), AlwaysTrueVerifier.verifyProof.selector);
        }
    }

    function _clientState() internal pure returns (IICS07TendermintMsgs.ClientState memory) {
        return IICS07TendermintMsgs.ClientState({
            chainId: CHAIN_ID,
            trustLevel: IICS07TendermintMsgs.TrustThreshold({numerator: 1, denominator: 3}),
            latestHeight: IICS02ClientMsgs.Height({revisionNumber: 0, revisionHeight: TRUSTED_HEIGHT}),
            trustingPeriod: TRUSTING_PERIOD,
            unbondingPeriod: UNBONDING_PERIOD,
            isFrozen: false,
            zkAlgorithm: IICS07TendermintMsgs.SupportedZkAlgorithm.Groth16
        });
    }

    function _buildValSet(uint16 valCount) internal pure returns (IICS07TendermintMsgs.ValidatorSet memory vs) {
        IICS07TendermintMsgs.ValidatorInfo[] memory vals = new IICS07TendermintMsgs.ValidatorInfo[](valCount);
        uint64 total = 0;
        for (uint256 i = 0; i < valCount; i++) {
            vals[i] = IICS07TendermintMsgs.ValidatorInfo({
                valAddress: abi.encodePacked(uint160(i + 1)),
                pubKey: bytes32(uint256(0xA0000000 + i + 1)),
                votingPower: 100,
                proposerPriority: 0
            });
            total += 100;
        }
        vs = IICS07TendermintMsgs.ValidatorSet({
            validators: vals,
            hasProposer: false,
            proposer: vals[0],
            totalVotingPower: total
        });
    }

    function _emptyValSet() internal pure returns (IICS07TendermintMsgs.ValidatorSet memory vs) {
        vs = IICS07TendermintMsgs.ValidatorSet({
            validators: new IICS07TendermintMsgs.ValidatorInfo[](0),
            hasProposer: false,
            proposer: IICS07TendermintMsgs.ValidatorInfo({
                valAddress: "",
                pubKey: bytes32(0),
                votingPower: 0,
                proposerPriority: 0
            }),
            totalVotingPower: 0
        });
    }

    /// Build commit signatures: first `activeCount` are COMMIT (with matching
    /// validator address), rest are ABSENT. Length == val_count per Tendermint
    /// invariant (`Predicates.validateCommit` enforces commit.length == vals.length).
    function _buildCommitSigs(IICS07TendermintMsgs.ValidatorSet memory vs, uint16 activeCount)
        internal pure returns (IICS07TendermintMsgs.CommitSig[] memory sigs)
    {
        sigs = new IICS07TendermintMsgs.CommitSig[](vs.validators.length);
        for (uint256 i = 0; i < vs.validators.length; i++) {
            if (i < activeCount) {
                sigs[i] = IICS07TendermintMsgs.CommitSig({
                    flag: IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_COMMIT,
                    data: IICS07TendermintMsgs.CommitSigData({
                        validatorAddress: vs.validators[i].valAddress,
                        timestamp: NEW_TS_NS,
                        hasSignature: false,
                        signature: ""
                    })
                });
            } else {
                sigs[i] = IICS07TendermintMsgs.CommitSig({
                    flag: IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_ABSENT,
                    data: IICS07TendermintMsgs.CommitSigData({
                        validatorAddress: "",
                        timestamp: 0,
                        hasSignature: false,
                        signature: ""
                    })
                });
            }
        }
    }

    function _buildHeader(
        uint64 trustedHeight,
        uint64 newHeight,
        uint128 newTimestamp,
        bytes32 appHash,
        IICS07TendermintMsgs.ValidatorSet memory currentValSet,
        IICS07TendermintMsgs.ValidatorSet memory trustedNextValSet,
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
        bh.time = newTimestamp;
        bh.appHash = appHash;
        bh.validatorsHash = currentValSetHash;
        bh.nextValidatorsHash = currentValSetHash;
        bytes32 headerHash = Header.hashHeader(bh);

        IICS07TendermintMsgs.BlockCommit memory bc = IICS07TendermintMsgs.BlockCommit({
            height: newHeight,
            round: 0,
            blockId: IICS07TendermintMsgs.BlockId({
                hashData: headerHash,
                partSetHeader: IICS07TendermintMsgs.PartSetHeader({total: 1, hashData: bytes32(uint256(0x9A57))})
            }),
            commitSigs: _buildCommitSigs(currentValSet, activeCount)
        });

        header = IICS07TendermintMsgs.Header({
            signedHeader: IICS07TendermintMsgs.SignedHeader({header: bh, commit: bc}),
            validatorSet: currentValSet,
            trustedHeight: IICS02ClientMsgs.Height({revisionNumber: 0, revisionHeight: trustedHeight}),
            trustedNextValidatorSet: trustedNextValSet
        });
    }

    /// Build a fully consistent Header: validatorsHash matches hashValSet,
    /// blockId.hashData matches hashHeader, trustedNextValSet hashes back to
    /// the trustedConsensusState.nextValidatorsHash.
    function _buildSelfConsistent(BucketConfig memory cfg)
        internal
        pure
        returns (
            IICS07TendermintMsgs.Header memory header,
            IICS07TendermintMsgs.ConsensusState memory trustedCS
        )
    {
        IICS07TendermintMsgs.ValidatorSet memory vs = _buildValSet(cfg.valCount);
        bytes32 valSetHash = Header.hashValSet(vs);

        IICS07TendermintMsgs.BlockHeader memory bh;
        bh.chainId = CHAIN_ID;
        bh.height = NEW_HEIGHT;
        bh.time = NEW_TS_NS;
        bh.appHash = bytes32(uint256(0xCCC));
        bh.validatorsHash = valSetHash;
        bh.nextValidatorsHash = valSetHash;
        bytes32 headerHash = Header.hashHeader(bh);

        IICS07TendermintMsgs.BlockCommit memory bc = IICS07TendermintMsgs.BlockCommit({
            height: NEW_HEIGHT,
            round: 0,
            blockId: IICS07TendermintMsgs.BlockId({
                hashData: headerHash,
                partSetHeader: IICS07TendermintMsgs.PartSetHeader({total: 1, hashData: bytes32(uint256(0x9A57))})
            }),
            commitSigs: _buildCommitSigs(vs, cfg.activeCount)
        });

        header = IICS07TendermintMsgs.Header({
            signedHeader: IICS07TendermintMsgs.SignedHeader({header: bh, commit: bc}),
            validatorSet: vs,
            trustedHeight: IICS02ClientMsgs.Height({revisionNumber: 0, revisionHeight: TRUSTED_HEIGHT}),
            trustedNextValidatorSet: vs
        });

        // trustedConsensusState.nextValidatorsHash must equal hashValSet(trustedNextValSet)
        // (we reuse the same set, so the hash is identical).
        trustedCS = IICS07TendermintMsgs.ConsensusState({
            timestamp: TRUSTED_TS_NS,
            root: bytes32(uint256(0xAAA)),
            nextValidatorsHash: valSetHash
        });
    }

    function _deployLightClient(
        IICS07TendermintMsgs.ClientState memory cs,
        IICS07TendermintMsgs.ConsensusState memory trustedCS
    )
        internal
        returns (Groth16ICS07Tendermint)
    {
        // Deploy a fresh light client per bucket, pre-seeded with this run's
        // trustedConsensusState hash at TRUSTED_HEIGHT.
        return new Groth16ICS07Tendermint(
            address(wrapper),
            STUB_MEMBERSHIP,
            STUB_MISBEHAVIOUR,
            address(updateClientImpl),
            abi.encode(cs),
            keccak256(abi.encode(trustedCS)),
            address(0) // permissionless
        );
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
        returns (IUpdateClientMsgs.MsgUpdateClient memory)
    {
        // Build per-slot bucket arrays (active = first cfg.activeCount, rest = padding).
        uint32[] memory idx = new uint32[](bucket);
        bytes32[] memory pks = new bytes32[](bucket);
        uint64[] memory tsS = new uint64[](bucket);
        uint32[] memory tsN = new uint32[](bucket);
        bool[] memory act = new bool[](bucket);
        for (uint256 i = 0; i < bucket; i++) {
            if (i < activeCount) {
                idx[i] = uint32(i);
                pks[i] = header.validatorSet.validators[i].pubKey;
                act[i] = true;
            }
            tsS[i] = uint64(1_700_000_000 + i);
            tsN[i] = uint32(i * 1_000_000);
        }

        return IUpdateClientMsgs.MsgUpdateClient({
            clientState: cs,
            trustedConsensusState: trustedCS,
            proposedHeader: header,
            time: header.signedHeader.header.time,
            proof: [uint256(0), 0, 0, 0, 0, 0, 0, 0],
            commitments: [uint256(0), 0],
            commitmentPok: [uint256(0), 0],
            bucket: bucket,
            signerIndices: idx,
            signerPubkeys: pks,
            timestampSeconds: tsS,
            timestampNanos: tsN,
            active: act
        });
    }

    function _measure(uint16 bucket) internal {
        BucketConfig memory cfg = _cfg(bucket);

        (IICS07TendermintMsgs.Header memory header, IICS07TendermintMsgs.ConsensusState memory trustedCS) =
            _buildSelfConsistent(cfg);

        IICS07TendermintMsgs.ClientState memory cs = _clientState();
        Groth16ICS07Tendermint ics07 = _deployLightClient(cs, trustedCS);
        bytes memory encoded = abi.encode(_buildMsg(cs, trustedCS, header, bucket, cfg.activeCount));

        uint256 g0 = gasleft();
        ILightClientMsgs.UpdateResult result = ics07.updateClient(encoded);
        uint256 used = g0 - gasleft();
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.Update), "expected Update");
        console.log("bucket=", bucket, "  gas=", used);
    }

    function _measureCacheHitReplay(uint16 bucket) internal {
        BucketConfig memory cfg = _cfg(bucket);

        (IICS07TendermintMsgs.Header memory header, IICS07TendermintMsgs.ConsensusState memory trustedCS) =
            _buildSelfConsistent(cfg);

        IICS07TendermintMsgs.ClientState memory cs = _clientState();
        Groth16ICS07Tendermint ics07 = _deployLightClient(cs, trustedCS);

        IUpdateClientMsgs.MsgUpdateClient memory fullMsg = _buildMsg(cs, trustedCS, header, bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(fullMsg))), uint8(ILightClientMsgs.UpdateResult.Update), "warmup Update");

        IUpdateClientMsgs.MsgUpdateClient memory cacheMsg = _buildMsg(cs, trustedCS, header, bucket, cfg.activeCount);
        cacheMsg.proposedHeader.validatorSet = _emptyValSet();
        cacheMsg.proposedHeader.trustedNextValidatorSet = _emptyValSet();

        uint256 g0 = gasleft();
        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(cacheMsg));
        uint256 used = g0 - gasleft();
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.NoOp), "expected cache-hit replay NoOp");
        console.log("bucket=", bucket, "  cache-hit replay gas=", used);
    }

    function _measureCacheHitAdjacentUpdate(uint16 bucket) internal {
        BucketConfig memory cfg = _cfg(bucket);
        IICS07TendermintMsgs.ClientState memory cs = _clientState();
        IICS07TendermintMsgs.ValidatorSet memory valSet = _buildValSet(cfg.valCount);
        bytes32 valSetHash = Header.hashValSet(valSet);

        IICS07TendermintMsgs.ConsensusState memory trustedCS0 = IICS07TendermintMsgs.ConsensusState({
            timestamp: TRUSTED_TS_NS,
            root: bytes32(uint256(0xAAA1)),
            nextValidatorsHash: valSetHash
        });

        Groth16ICS07Tendermint ics07 = _deployLightClient(cs, trustedCS0);

        IICS07TendermintMsgs.Header memory header1001 =
            _buildHeader(TRUSTED_HEIGHT, NEW_HEIGHT, NEW_TS_NS, bytes32(uint256(0xCCC1)), valSet, valSet, cfg.activeCount);
        IUpdateClientMsgs.MsgUpdateClient memory msg1001 =
            _buildMsg(cs, trustedCS0, header1001, bucket, cfg.activeCount);
        assertEq(uint8(ics07.updateClient(abi.encode(msg1001))), uint8(ILightClientMsgs.UpdateResult.Update), "warmup Update");

        IICS07TendermintMsgs.ConsensusState memory trustedCS1001 = IICS07TendermintMsgs.ConsensusState({
            timestamp: NEW_TS_NS,
            root: header1001.signedHeader.header.appHash,
            nextValidatorsHash: valSetHash
        });

        IICS07TendermintMsgs.Header memory header1002 = _buildHeader(
            NEW_HEIGHT,
            NEXT_HEIGHT,
            NEXT_TS_NS,
            bytes32(uint256(0xCCC2)),
            valSet,
            _emptyValSet(),
            cfg.activeCount
        );

        IUpdateClientMsgs.MsgUpdateClient memory cacheMsg =
            _buildMsg(cs, trustedCS1001, header1002, bucket, cfg.activeCount);
        cacheMsg.proposedHeader.validatorSet = _emptyValSet();
        cacheMsg.proposedHeader.trustedNextValidatorSet = _emptyValSet();

        uint256 g0 = gasleft();
        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(cacheMsg));
        uint256 used = g0 - gasleft();
        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.Update), "expected cache-hit adjacent Update");
        console.log("bucket=", bucket, "  cache-hit adjacent update gas=", used);
    }

    function test_gas_n4() public { _measure(4); }
    function test_gas_n8() public { _measure(8); }
    function test_gas_n16() public { _measure(16); }
    function test_gas_n16_cacheHitReplay() public { _measureCacheHitReplay(16); }
    function test_gas_n16_cacheHitAdjacentUpdate() public { _measureCacheHitAdjacentUpdate(16); }
    function test_gas_n32() public { _measure(32); }
    function test_gas_n64() public { _measure(64); }
    function test_gas_n128() public { _measure(128); }
}
