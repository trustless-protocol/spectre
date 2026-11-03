// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { SignatureVerifier } from "contracts/light-clients/spectre/SignatureVerifier.sol";
import { ISignatureVerifier } from "contracts/light-clients/spectre/interfaces/ISignatureVerifier.sol";
import { Groth16Verifier_N4 } from "contracts/verifiers/Groth16Verifier_N4.sol";

/// @title Groth16OnchainProofTest
/// @notice Runs a real proof of the batch circuit through `SignatureVerifier.verifyBatchProof` and the
///         generated `Groth16Verifier_N4`, the end-to-end check that the contract side agrees with the
///         prover on the proof layout: `Ar || Bs || Krs (256 B) || C0 || C1 (2 * 64 B) || PoK (64 B)`
///         = 448 bytes, two BSB22 commitments and one folded proof of knowledge.
///
///         The proof comes from the smoke test of `relayer/prover/cmd`, written by
///         `scripts/solidity/build-prover-artifacts.sh` to `.artifacts/solidity/prover/onchain_fixture_n4.json`
///         together with the verifier. Both are git-ignored and belong to one randomized setup, so the
///         fixture carries the SHA-256 of its verifier: with no fixture, or one of another verifier, the
///         tests skip instead of failing. Every tampered-input test must see `false` from a valid-shaped
///         call, not a revert: `verifyBatchProof` returns the verifier's `staticcall` result.
contract Groth16OnchainProofTest is Test {
    string internal constant FIXTURE = ".artifacts/solidity/prover/onchain_fixture_n4.json";
    string internal constant VERIFIER_SOURCE = "contracts/verifiers/Groth16Verifier_N4.sol";
    uint16 internal constant BUCKET = 4;

    SignatureVerifier internal signatureVerifier;

    uint256[8] internal proof;
    uint256[4] internal commitments;
    uint256[2] internal commitmentPok;
    bytes32[] internal pubkeys;
    bool[] internal active;
    ISignatureVerifier.SharedBlock internal shared;

    function setUp() public {
        string memory json;
        try vm.readFile(FIXTURE) returns (string memory contents) {
            json = contents;
        } catch {
            vm.skip(true, "no on-chain fixture: run scripts/solidity/build-prover-artifacts.sh --force");
            return;
        }
        if (vm.parseJsonBytes32(json, ".verifierSha256") != sha256(bytes(vm.readFile(VERIFIER_SOURCE)))) {
            vm.skip(true, "the fixture belongs to another Groth16Verifier_N4: regenerate the artifacts");
            return;
        }

        uint256[] memory words = vm.parseJsonUintArray(json, ".proof");
        assertEq(words.length, 8, "proof words");
        for (uint256 i = 0; i < 8; i++) {
            proof[i] = words[i];
        }
        words = vm.parseJsonUintArray(json, ".commitments");
        assertEq(words.length, 4, "two commitments are four words");
        for (uint256 i = 0; i < 4; i++) {
            commitments[i] = words[i];
        }
        words = vm.parseJsonUintArray(json, ".commitmentPok");
        assertEq(words.length, 2, "one folded proof of knowledge is two words");
        for (uint256 i = 0; i < 2; i++) {
            commitmentPok[i] = words[i];
        }
        pubkeys = vm.parseJsonBytes32Array(json, ".pubkeys");
        active = vm.parseJsonBoolArray(json, ".active");
        assertEq(pubkeys.length, BUCKET, "pubkeys");
        assertEq(active.length, BUCKET, "active");

        // PrefixHead = 0x08 0x02 0x11 || sfixed64(height); Round = sfixed64(round). Both little-endian.
        bytes memory prefixHead = vm.parseJsonBytes(json, ".prefixHead");
        bytes memory round = vm.parseJsonBytes(json, ".round");
        uint64 height;
        uint64 roundValue;
        for (uint256 i = 0; i < 8; i++) {
            height |= uint64(uint8(prefixHead[3 + i])) << uint64(8 * i);
            roundValue |= uint64(uint8(round[i])) << uint64(8 * i);
        }
        shared = ISignatureVerifier.SharedBlock({
            height: height, round: roundValue, blockIDHash: vm.parseJsonBytes32(json, ".blockHash")
        });

        signatureVerifier = new SignatureVerifier(address(this));
        Groth16Verifier_N4 verifier = new Groth16Verifier_N4();
        signatureVerifier.setBucket(BUCKET, address(verifier), verifier.verifyProof.selector);
    }

    function _verify() internal view returns (bool) {
        return signatureVerifier.verifyBatchProof(BUCKET, proof, commitments, commitmentPok, pubkeys, active, shared);
    }

    // ------------------------------------------------------------------
    // The honest proof
    // ------------------------------------------------------------------

    function test_acceptsTheSmokeProof() public view {
        assertTrue(_verify(), "the smoke proof must verify on-chain");
    }

    /// @notice The blob `verifyBatchProof` hands the verifier is 448 bytes: a verifier exported with
    ///         fewer commitments than the prover emits would read a different layout.
    function test_proofBlobIs448Bytes() public view {
        assertEq(abi.encodePacked(proof, commitments, commitmentPok).length, 448);
    }

    // ------------------------------------------------------------------
    // The second commitment is checked, not ignored
    // ------------------------------------------------------------------

    function test_rejectsSwappedCommitments() public {
        (uint256 x0, uint256 y0) = (commitments[0], commitments[1]);
        (commitments[0], commitments[1]) = (commitments[2], commitments[3]);
        (commitments[2], commitments[3]) = (x0, y0);
        assertFalse(_verify(), "C0 and C1 swapped");
    }

    function test_rejectsSecondCommitmentReplacedByTheFirst() public {
        (commitments[2], commitments[3]) = (commitments[0], commitments[1]);
        assertFalse(_verify(), "C1 := C0");
    }

    // ------------------------------------------------------------------
    // Every word of the proof is bound. A tampered value is usually not a curve point, and the precompile
    // then burns the gas it was given: one tampered input per test keeps each call within the test's gas.
    // ------------------------------------------------------------------

    function test_rejectsTamperedCommitmentWord0() public {
        commitments[0] ^= 1;
        assertFalse(_verify(), "tampered commitment word accepted");
    }

    function test_rejectsTamperedCommitmentWord1() public {
        commitments[1] ^= 1;
        assertFalse(_verify(), "tampered commitment word accepted");
    }

    function test_rejectsTamperedCommitmentWord2() public {
        commitments[2] ^= 1;
        assertFalse(_verify(), "tampered commitment word accepted");
    }

    function test_rejectsTamperedCommitmentWord3() public {
        commitments[3] ^= 1;
        assertFalse(_verify(), "tampered commitment word accepted");
    }

    function test_rejectsTamperedCommitmentPokWord0() public {
        commitmentPok[0] ^= 1;
        assertFalse(_verify(), "tampered proof of knowledge accepted");
    }

    function test_rejectsTamperedCommitmentPokWord1() public {
        commitmentPok[1] ^= 1;
        assertFalse(_verify(), "tampered proof of knowledge accepted");
    }

    function test_rejectsTamperedProofWord0() public {
        proof[0] ^= 1;
        assertFalse(_verify(), "tampered proof word accepted");
    }

    function test_rejectsTamperedProofWord1() public {
        proof[1] ^= 1;
        assertFalse(_verify(), "tampered proof word accepted");
    }

    function test_rejectsTamperedProofWord2() public {
        proof[2] ^= 1;
        assertFalse(_verify(), "tampered proof word accepted");
    }

    function test_rejectsTamperedProofWord3() public {
        proof[3] ^= 1;
        assertFalse(_verify(), "tampered proof word accepted");
    }

    function test_rejectsTamperedProofWord4() public {
        proof[4] ^= 1;
        assertFalse(_verify(), "tampered proof word accepted");
    }

    function test_rejectsTamperedProofWord5() public {
        proof[5] ^= 1;
        assertFalse(_verify(), "tampered proof word accepted");
    }

    function test_rejectsTamperedProofWord6() public {
        proof[6] ^= 1;
        assertFalse(_verify(), "tampered proof word accepted");
    }

    function test_rejectsTamperedProofWord7() public {
        proof[7] ^= 1;
        assertFalse(_verify(), "tampered proof word accepted");
    }

    // ------------------------------------------------------------------
    // The public input is bound
    // ------------------------------------------------------------------

    function test_rejectsFlippedActiveFlag() public {
        active[0] = !active[0];
        assertFalse(_verify(), "active flag flipped");
    }

    function test_rejectsSwappedPubkey() public {
        pubkeys[0] = pubkeys[1];
        assertFalse(_verify(), "pubkey swapped");
    }

    function test_rejectsOtherHeight() public {
        shared.height += 1;
        assertFalse(_verify(), "height changed");
    }

    function test_rejectsOtherBlockHash() public {
        shared.blockIDHash = ~shared.blockIDHash;
        assertFalse(_verify(), "block hash changed");
    }
}
