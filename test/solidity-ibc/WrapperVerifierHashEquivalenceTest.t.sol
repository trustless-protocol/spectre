// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { WrapperVerifier } from "../../contracts/utils/WrapperVerifier.sol";
import { Encode } from "../../contracts/utils/Encode.sol";
import { IVerifier } from "../../contracts/interfaces/IVerifier.sol";

contract WrapperVerifierHashHarness is WrapperVerifier {
    constructor() WrapperVerifier(address(this)) { }

    function hashWitness(
        bytes32[] calldata pubkeys,
        bool[] calldata active,
        IVerifier.SharedBlock calldata shared
    )
        external
        pure
        returns (bytes32)
    {
        return _hashWitness(pubkeys, active, shared);
    }
}

/// @notice Locks the #199 witness layout
///   PrefixHead(11) || roundPresent(1) || BlockHash(32) || per slot: active(1) || A(32)
/// against an independent in-test reference. The production `_hashWitness`
/// builds PrefixHead with low-level mstore/mcopy; the reference re-derives the
/// same bytes via `abi.encodePacked`, so any drift in offsets, the round-present
/// byte, or per-slot packing is caught Solidity-vs-Solidity. The off-chain Go
/// path (prover/hash_witness.go) and the in-circuit commit are cross-checked
/// separately via test.IsSolved.
contract WrapperVerifierHashEquivalenceTest is Test {
    WrapperVerifierHashHarness internal harness;

    function setUp() public {
        harness = new WrapperVerifierHashHarness();
    }

    function test_hashWitness_matchesReference_bucket16_mixedSlots() public view {
        uint16 bucket = 16;
        bytes32[] memory pubkeys = new bytes32[](bucket);
        bool[] memory active = new bool[](bucket);
        for (uint256 i = 0; i < bucket; i++) {
            pubkeys[i] = bytes32(uint256(0xA0 + i));
            active[i] = i < 11;
        }

        IVerifier.SharedBlock memory shared =
            IVerifier.SharedBlock({ height: 53, round: 0, blockIDHash: bytes32(uint256(0x1234)) });

        assertEq(harness.hashWitness(pubkeys, active, shared), _hashWitnessReference(pubkeys, active, shared));
    }

    function test_hashWitness_matchesReference_multiByteHeightAndRoundPresent() public view {
        uint16 bucket = 4;
        bytes32[] memory pubkeys = new bytes32[](bucket);
        bool[] memory active = new bool[](bucket);
        pubkeys[0] = bytes32(uint256(1));
        pubkeys[1] = bytes32(uint256(2));
        pubkeys[2] = bytes32(uint256(3));
        pubkeys[3] = bytes32(uint256(4));
        active[0] = true;
        active[1] = false;
        active[2] = true;
        active[3] = false;

        IVerifier.SharedBlock memory shared =
            IVerifier.SharedBlock({ height: 16_384, round: 127, blockIDHash: bytes32(uint256(0xCAFE)) });

        assertEq(harness.hashWitness(pubkeys, active, shared), _hashWitnessReference(pubkeys, active, shared));
    }

    /// @dev Independent reference for the #199 layout. PrefixHead = Type(0x08
    ///      0x02) || Height(0x11 || sfixed64); roundPresent = (round > 0);
    ///      BlockHash = 32 bytes; then per slot active(1) || pubkey(32).
    function _hashWitnessReference(
        bytes32[] memory pubkeys,
        bool[] memory active,
        IVerifier.SharedBlock memory shared
    )
        internal
        pure
        returns (bytes32)
    {
        bytes memory buf = abi.encodePacked(
            hex"080211",
            Encode.encodeSfixed64(int64(uint64(shared.height))),
            shared.round > 0 ? bytes1(0x01) : bytes1(0x00),
            shared.blockIDHash
        );
        for (uint256 i = 0; i < pubkeys.length; i++) {
            buf = abi.encodePacked(buf, active[i] ? bytes1(0x01) : bytes1(0x00), pubkeys[i]);
        }
        return sha256(buf);
    }
}
