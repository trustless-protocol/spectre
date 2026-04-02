// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";
import { stdJson } from "forge-std/StdJson.sol";
import { IVerifier, IGroth16Verifier } from "../../contracts/interfaces/IVerifier.sol";
import { WrapperVerifier } from "../../contracts/utils/WrapperVerifier.sol";
import { Groth16Verifier } from "../../contracts/utils/Groth16Verifier.sol";

/// @dev Harness to expose WrapperVerifier internal decodePoint for unit testing.
contract WrapperVerifierHarness is WrapperVerifier {
    constructor(IGroth16Verifier v) WrapperVerifier(v) { }

    function exposeDecodePoint(uint256 compressed) external returns (uint256 x, uint256 y) {
        return decodePoint(compressed);
    }
}

struct Groth16Fixture {
    uint256[8] proof;
    uint256[2] commitments;
    uint256[2] commitmentPok;
    bytes32 signatureR;
    bytes32 signatureS;
    bytes32 pubkey;
    bytes32 message;
}

contract Groth16VerifierTest is Test {
    using stdJson for string;

    WrapperVerifier internal wrapper;
    WrapperVerifierHarness internal harness;
    Groth16Fixture internal fixture;

    function setUp() public {
        // Deploy Groth16Verifier and WrapperVerifier
        Groth16Verifier groth16Verifier = new Groth16Verifier();
        wrapper = new WrapperVerifier(IGroth16Verifier(address(groth16Verifier)));
        harness = new WrapperVerifierHarness(IGroth16Verifier(address(groth16Verifier)));

        // Load fixture
        string memory root = vm.projectRoot();
        string memory path = string.concat(root, "/test/groth16-ics07/fixtures/groth16_fixture.json");
        string memory json = vm.readFile(path);

        fixture.proof[0] = json.readUint(".proof[0]");
        fixture.proof[1] = json.readUint(".proof[1]");
        fixture.proof[2] = json.readUint(".proof[2]");
        fixture.proof[3] = json.readUint(".proof[3]");
        fixture.proof[4] = json.readUint(".proof[4]");
        fixture.proof[5] = json.readUint(".proof[5]");
        fixture.proof[6] = json.readUint(".proof[6]");
        fixture.proof[7] = json.readUint(".proof[7]");

        fixture.commitments[0] = json.readUint(".commitments[0]");
        fixture.commitments[1] = json.readUint(".commitments[1]");

        fixture.commitmentPok[0] = json.readUint(".commitmentPok[0]");
        fixture.commitmentPok[1] = json.readUint(".commitmentPok[1]");

        fixture.signatureR = json.readBytes32(".signatureR");
        fixture.signatureS = json.readBytes32(".signatureS");
        fixture.pubkey = json.readBytes32(".pubkey");
        fixture.message = json.readBytes32(".message");
    }

    function testVerifyValidProof() public {
        bytes32[2] memory signature = [fixture.signatureR, fixture.signatureS];

        bool result = wrapper.verifyProof(
            fixture.proof,
            fixture.commitments,
            fixture.commitmentPok,
            signature,
            fixture.pubkey,
            abi.encodePacked(fixture.message)
        );
        assertTrue(result, "valid proof should verify");
    }

    function testVerifyInvalidProofReverts() public {
        // Tamper with proof
        uint256[8] memory badProof = fixture.proof;
        badProof[0] = badProof[0] ^ 1;

        bytes32[2] memory signature = [fixture.signatureR, fixture.signatureS];

        bool result = wrapper.verifyProof(
            badProof,
            fixture.commitments,
            fixture.commitmentPok,
            signature,
            fixture.pubkey,
            abi.encodePacked(fixture.message)
        );
        assertFalse(result, "tampered proof should not verify");
    }

    function testVerifyWrongPubkeyReverts() public {
        bytes32[2] memory signature = [fixture.signatureR, fixture.signatureS];
        bytes32 wrongPubkey = bytes32(uint256(fixture.pubkey) ^ 1);

        bool result = wrapper.verifyProof(
            fixture.proof,
            fixture.commitments,
            fixture.commitmentPok,
            signature,
            wrongPubkey,
            abi.encodePacked(fixture.message)
        );
        assertFalse(result, "wrong pubkey should not verify");
    }

    function testVerifyWrongMessageReverts() public {
        bytes32[2] memory signature = [fixture.signatureR, fixture.signatureS];
        bytes32 wrongMessage = bytes32(uint256(fixture.message) ^ 1);

        bool result = wrapper.verifyProof(
            fixture.proof,
            fixture.commitments,
            fixture.commitmentPok,
            signature,
            fixture.pubkey,
            abi.encodePacked(wrongMessage)
        );
        assertFalse(result, "wrong message should not verify");
    }

    /// @dev Test decodePoint against the Ed25519 base point G (RFC 8032 §5.1).
    ///      Input is the compressed point in LE uint256 form (as passed to decodePoint after reverseBytes).
    ///      y = 4/5 mod p = 0x6666...6658; x is even so sign bit = 0.
    function testDecodePointBasePoint() public {
        // Ed25519 base point G in compressed LE form: bit 255 = sign(x) = 0 (x is even)
        uint256 compressed = 0x6666666666666666666666666666666666666666666666666666666666666658;

        (uint256 x, uint256 y) = harness.exposeDecodePoint(compressed);

        // RFC 8032 §5.1 base point coordinates
        assertEq(y, 0x6666666666666666666666666666666666666666666666666666666666666658, "y mismatch");
        assertEq(x, 0x216936d3cd6e53fec0a4e231fdd6dc5c692cc7609525a7b2c9562d608f25d51a, "x mismatch");
    }

    /// @dev Test decodePoint with the odd-x variant of the base point (sign bit = 1).
    ///      The negation of x mod p should be recovered when sign bit != parity of x.
    function testDecodePointBasePointNegX() public {
        uint256 p = 0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffed;
        // Set sign bit = 1 on the base point y value → should recover p - x (odd x)
        uint256 compressed = 0x6666666666666666666666666666666666666666666666666666666666666658
            | (1 << 255);

        (uint256 x,) = harness.exposeDecodePoint(compressed);

        uint256 expectedX = p - 0x216936d3cd6e53fec0a4e231fdd6dc5c692cc7609525a7b2c9562d608f25d51a;
        assertEq(x, expectedX, "negated x mismatch");
    }
}
