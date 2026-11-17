// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";
import { Groth16Verifier_N4 } from "contracts/verifiers/Groth16Verifier_N4.sol";

/// Replays a historical relayer proof against the locally-compiled verifier. The proof was made for an
/// older circuit: one BSB22 commitment (a 384-byte proof) and another verifying key. The regenerated
/// verifier must reject it, and each test pins the reason, so the proof cannot be rejected by accident:
///   - as recorded, it is refused for its length before any curve arithmetic;
///   - padded to the current 448-byte layout with a second commitment, it reaches the commitment check
///     and fails it.
contract ReplayProofTest is Test {
    Groth16Verifier_N4 internal verifier;

    function setUp() public {
        verifier = new Groth16Verifier_N4();
    }

    function test_LegacyProofIsRejectedForItsLength() public {
        (uint256[8] memory proof, uint256[2] memory commitment, uint256[2] memory commitmentPok) = _legacyProof();
        bytes memory proofBytes = abi.encodePacked(proof, commitment, commitmentPok);
        assertEq(proofBytes.length, 384, "the legacy layout carries one commitment");

        vm.expectRevert(bytes("invalid proof length"));
        verifier.verifyProof(proofBytes, _legacyPublicInputs());
    }

    function test_LegacyProofPaddedToTwoCommitmentsFailsTheCommitmentCheck() public {
        (uint256[8] memory proof, uint256[2] memory commitment, uint256[2] memory commitmentPok) = _legacyProof();
        // The legacy commitment is a curve point, so repeating it as C1 gives a well-formed proof of the
        // current layout.
        bytes memory proofBytes = abi.encodePacked(proof, commitment, commitment, commitmentPok);
        assertEq(proofBytes.length, 448, "the current layout carries two commitments");

        vm.expectRevert(Groth16Verifier_N4.CommitmentInvalid.selector);
        verifier.verifyProof(proofBytes, _legacyPublicInputs());
    }

    function _legacyProof()
        private
        pure
        returns (uint256[8] memory proof, uint256[2] memory commitment, uint256[2] memory commitmentPok)
    {
        proof = [
            uint256(0x21c1dce7b6d9248d4f4a85e2027328e35590fa8dd9bdcdac246b069aa50977ed),
            uint256(0x1c231441c222a22b4108e977f351b4a2edcad024ca67d098df6930d89f914a74),
            uint256(0x017f6268d91c95b3d0fd5cbbe40f24735765865a20472f3984641a55707157ee),
            uint256(0x08a31b2ab391ae33ec5b716ac62dc2d3b9d5a4a900d978c87fae73765e0ea921),
            uint256(0x09c32888067a17882316ed6c46409d695d2887b6f4c86f100b3e181ffd11e99e),
            uint256(0x21d278cf625021fbbe048975127eaea5569a577a4d424754be838f9c9bffb488),
            uint256(0x075420bebef7b38c7fd7a9f57589c49d16edb108f28175b3b979fcdbe82dd8b6),
            uint256(0x08c35e2d388b3951b1271733f0600543b84332c2068a85130863b3edb2bb5fac)
        ];
        commitment = [
            uint256(0x00a1cce8912ac36d680cfc826a1fd960351d012bed1b94aa3d207d52975b9af9),
            uint256(0x25713a539188e279488031f63dba905906a6dc65004d70717f2bee99da3753ec)
        ];
        commitmentPok = [
            uint256(0x2a7ce4360031b6a384d2a31f6ecd7e75eb45403542eb713c10861cea1e559ef3),
            uint256(0x1c14885f5295b8ae68f8f60b1560e344960b8fc1f6e5308ea8732633719a9890)
        ];
    }

    /// @dev The witness digest the legacy proof was made for, split into the two 128-bit public inputs.
    function _legacyPublicInputs() private pure returns (uint256[2] memory publicInputs) {
        uint256 digest = uint256(0x3bbb741e27fe3429a6de507cacefd4f3a4b478995eae56e546c869910c0a6b8d);
        publicInputs[0] = digest >> 128;
        publicInputs[1] = digest & type(uint128).max;
    }
}
