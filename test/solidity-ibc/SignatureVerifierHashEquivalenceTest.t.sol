// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { SignatureVerifier } from "../../contracts/light-clients/SignatureVerifier.sol";
import { Encode } from "../../contracts/utils/Encode.sol";
import { ISignatureVerifier } from "../../contracts/light-clients/interfaces/ISignatureVerifier.sol";

contract SignatureVerifierHashHarness is SignatureVerifier {
    constructor() SignatureVerifier(address(this)) { }

    function hashWitness(
        bytes32[] calldata pubkeys,
        bool[] calldata active,
        ISignatureVerifier.SharedBlock calldata shared
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
contract SignatureVerifierHashEquivalenceTest is Test {
    SignatureVerifierHashHarness internal harness;

    function setUp() public {
        harness = new SignatureVerifierHashHarness();
    }

    function test_hashWitness_matchesReference_bucket16_mixedSlots() public view {
        uint16 bucket = 16;
        bytes32[] memory pubkeys = new bytes32[](bucket);
        bool[] memory active = new bool[](bucket);
        for (uint256 i = 0; i < bucket; i++) {
            pubkeys[i] = bytes32(uint256(0xA0 + i));
            active[i] = i < 11;
        }

        ISignatureVerifier.SharedBlock memory shared =
            ISignatureVerifier.SharedBlock({ height: 53, round: 0, blockIDHash: bytes32(uint256(0x1234)) });

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

        ISignatureVerifier.SharedBlock memory shared =
            ISignatureVerifier.SharedBlock({ height: 16_384, round: 127, blockIDHash: bytes32(uint256(0xCAFE)) });

        assertEq(harness.hashWitness(pubkeys, active, shared), _hashWitnessReference(pubkeys, active, shared));
    }

    /// @notice ZK-07: cross-checks `_hashWitness` against digests computed by the Go prover.
    /// @dev The witness commitment is computed in Go (prover.ComputeWitnessHash), in-circuit
    ///      (BatchCircuit.Define) and here, and a proof only verifies when all three agree — but
    ///      the tests above compare Solidity to Solidity, so a Go-vs-Solidity drift would have
    ///      surfaced only in E2E. The vectors are produced by
    ///      `go test ./prover -run TestWitnessHashVectors -update` from `relayer/`, which also
    ///      re-verifies them on every run, so the generator cannot rot away from the fixture.
    ///
    ///      A failure here means the two layouts disagree. Read the Go test's output first: if it
    ///      passes and this fails, Solidity drifted; if both fail, the fixture is stale.
    function test_hashWitness_matchesGoVectors() public view {
        string memory raw = vm.readFile("test/fixtures/witness_hash_vectors.json");
        uint256 count = vm.parseJsonUint(raw, ".count");
        assertGt(count, 0, "no cross-check vectors");

        for (uint256 v = 0; v < count; v++) {
            string memory base = string.concat(".vectors[", vm.toString(v), "]");
            string memory name = vm.parseJsonString(raw, string.concat(base, ".name"));

            bytes32[] memory pubkeys = vm.parseJsonBytes32Array(raw, string.concat(base, ".pubkeys"));
            bool[] memory active = vm.parseJsonBoolArray(raw, string.concat(base, ".active"));
            assertEq(pubkeys.length, active.length, name);

            ISignatureVerifier.SharedBlock memory shared = ISignatureVerifier.SharedBlock({
                height: uint64(vm.parseJsonUint(raw, string.concat(base, ".height"))),
                round: uint64(vm.parseJsonUint(raw, string.concat(base, ".round"))),
                blockIDHash: vm.parseJsonBytes32(raw, string.concat(base, ".blockIdHash"))
            });

            assertEq(
                harness.hashWitness(pubkeys, active, shared),
                vm.parseJsonBytes32(raw, string.concat(base, ".witnessHash")),
                name
            );
        }
    }

    /// @dev Independent reference for the #199 layout. PrefixHead = Type(0x08
    ///      0x02) || Height(0x11 || sfixed64); roundPresent = (round > 0);
    ///      BlockHash = 32 bytes; then per slot active(1) || pubkey(32).
    function _hashWitnessReference(
        bytes32[] memory pubkeys,
        bool[] memory active,
        ISignatureVerifier.SharedBlock memory shared
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
