// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { ChainId } from "../../contracts/utils/ChainId.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";

/// @notice Unit tests for the extracted `ChainId.get(...)` parser.
/// @dev Mirrors the legacy `getChainId` behaviour that lived inline in
///      `UpdateClient.sol` and `Misbehaviour.sol`, plus exercises the edge
///      cases that previously had zero coverage.
contract ChainIdTest is Test {
    // Helper — exposes `ChainId.get` via the test contract so the call is a
    // single internal-library call that Foundry can measure.
    function _get(string memory id) internal pure returns (IICS07TendermintMsgs.ChainId memory) {
        return ChainId.get(id);
    }

    function test_cosmoshub_4() public pure {
        IICS07TendermintMsgs.ChainId memory c = _get("cosmoshub-4");
        assertEq(c.id, "cosmoshub-4");
        assertEq(uint256(c.revisionNumber), 4);
    }

    function test_largeRevision() public pure {
        // u64 max = 18446744073709551615
        IICS07TendermintMsgs.ChainId memory c = _get("chain-18446744073709551615");
        assertEq(uint256(c.revisionNumber), type(uint64).max);
    }

    function test_overflowRevision_returnsZero() public pure {
        // u64 max + 1
        IICS07TendermintMsgs.ChainId memory c = _get("chain-18446744073709551616");
        assertEq(uint256(c.revisionNumber), 0);
    }

    function test_singleZeroRevision_allowed() public pure {
        IICS07TendermintMsgs.ChainId memory c = _get("chain-0");
        assertEq(uint256(c.revisionNumber), 0);
    }

    function test_leadingZeroRevision_returnsZero() public pure {
        IICS07TendermintMsgs.ChainId memory c = _get("chain-01");
        assertEq(uint256(c.revisionNumber), 0);
    }

    function test_nonDigitSuffix_returnsZero() public pure {
        // "test-ibc-eth" — last segment "eth" is non-digit, treated as no revision.
        IICS07TendermintMsgs.ChainId memory c = _get("test-ibc-eth");
        assertEq(uint256(c.revisionNumber), 0);
    }

    function test_multiDashChain_usesLastDash() public pure {
        // last '-' is before "7"; prefix length = 8 ("test-ibc") which is in range.
        IICS07TendermintMsgs.ChainId memory c = _get("test-ibc-7");
        assertEq(uint256(c.revisionNumber), 7);
    }

    function test_noDash_returnsZero() public pure {
        IICS07TendermintMsgs.ChainId memory c = _get("nodashatall");
        assertEq(uint256(c.revisionNumber), 0);
    }

    function test_emptyString_reverts() public {
        vm.expectRevert(ChainId.InvalidChainIdLength.selector);
        this.externalGet("");
    }

    function test_tooLong_reverts() public {
        // 64 chars — exceeds the < 64 cap.
        string memory s = "0123456789012345678901234567890123456789012345678901234567890123";
        assertEq(bytes(s).length, 64);
        vm.expectRevert(ChainId.InvalidChainIdLength.selector);
        this.externalGet(s);
    }

    function test_prefixTooShort_reverts() public {
        // "a-1" — dashPos == 1, fails the `dashPos > 1` constraint.
        vm.expectRevert(ChainId.InvalidChainPrefixLength.selector);
        this.externalGet("a-1");
    }

    function test_prefixTooLong_reverts() public {
        // 43-char prefix + "-1" — dashPos == 43, fails the `dashPos < 43` constraint.
        string memory s = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQ-1";
        // Sanity: position of '-' is index 43.
        assertEq(bytes(s)[43], bytes1("-"));
        vm.expectRevert(ChainId.InvalidChainPrefixLength.selector);
        this.externalGet(s);
    }

    /// @dev External wrapper so `vm.expectRevert` can intercept the revert
    ///      from the internal library call.
    function externalGet(string memory id) external pure returns (IICS07TendermintMsgs.ChainId memory) {
        return ChainId.get(id);
    }
}
