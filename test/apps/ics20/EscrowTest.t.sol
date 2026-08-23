// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

// solhint-disable max-line-length,gas-custom-errors,multiple-sends

import { Test } from "forge-std/Test.sol";

import { RateLimitErrors } from "contracts/apps/ics20/errors/RateLimitErrors.sol";
import { IAccessManaged } from "@openzeppelin-contracts/access/manager/IAccessManaged.sol";
import { IRateLimit } from "contracts/apps/ics20/interfaces/IRateLimit.sol";
import { IERC20 } from "@openzeppelin-contracts/token/ERC20/IERC20.sol";

import { Escrow } from "contracts/apps/ics20/Escrow.sol";
import { UpgradeableBeacon } from "@openzeppelin-contracts/proxy/beacon/UpgradeableBeacon.sol";
import { BeaconProxy } from "@openzeppelin-contracts/proxy/beacon/BeaconProxy.sol";
import { AccessManager } from "@openzeppelin-contracts/access/manager/AccessManager.sol";
import { IBCRolesLib } from "contracts/shared/access/IBCRolesLib.sol";

contract EscrowTest is Test {
    address public rateLimiter = makeAddr("rateLimiter");

    Escrow public escrow;
    AccessManager public accessManager;

    function setUp() public {
        // setup code here
        address _escrowLogic = address(new Escrow());
        address escrowBeacon = address(new UpgradeableBeacon(_escrowLogic, address(this)));
        accessManager = new AccessManager(address(this));

        BeaconProxy escrowProxy =
            new BeaconProxy(escrowBeacon, abi.encodeCall(Escrow.initialize, (address(this), address(accessManager))));
        escrow = Escrow(address(escrowProxy));
        assert(escrow.ics20() == address(this));

        accessManager.grantRole(IBCRolesLib.RATE_LIMITER_ROLE, rateLimiter, 0);
        (bool hasRole, uint32 execDelay) = accessManager.hasRole(IBCRolesLib.RATE_LIMITER_ROLE, rateLimiter);
        assertTrue(hasRole, "Rate limiter role not granted");
        assertEq(execDelay, 0, "Rate limiter role needs 0 delay");
        assertEq(
            accessManager.getTargetFunctionRole(address(escrow), IRateLimit.setRateLimit.selector),
            IBCRolesLib.ADMIN_ROLE,
            "rate limit selector should not need a target mapping"
        );
    }

    function test_success_setRateLimit() public {
        address mockToken = makeAddr("mockToken");
        uint256 rateLimit = 10_000;

        vm.prank(rateLimiter);
        escrow.setRateLimit(mockToken, rateLimit);
        assertEq(escrow.getRateLimit(mockToken), rateLimit);
    }

    function test_failure_setRateLimit() public {
        address unauthorized = makeAddr("unauthorized");
        address mockToken = makeAddr("mockToken");
        uint256 rateLimit = 10_000;

        vm.prank(unauthorized);
        vm.expectRevert(abi.encodeWithSelector(IAccessManaged.AccessManagedUnauthorized.selector, unauthorized));
        escrow.setRateLimit(mockToken, rateLimit);
        assertEq(escrow.getRateLimit(mockToken), 0);
    }

    function test_failure_setRateLimitWithExecutionDelay() public {
        address delayedRateLimiter = makeAddr("delayedRateLimiter");
        address mockToken = makeAddr("mockToken");

        accessManager.grantRole(IBCRolesLib.RATE_LIMITER_ROLE, delayedRateLimiter, 1 days);

        vm.prank(delayedRateLimiter);
        vm.expectRevert(abi.encodeWithSelector(IAccessManaged.AccessManagedUnauthorized.selector, delayedRateLimiter));
        escrow.setRateLimit(mockToken, 10_000);
    }

    function test_failure_setRateLimitWhenTargetClosed() public {
        address mockToken = makeAddr("mockToken");
        accessManager.setTargetClosed(address(escrow), true);

        vm.prank(rateLimiter);
        vm.expectRevert(abi.encodeWithSelector(IAccessManaged.AccessManagedUnauthorized.selector, rateLimiter));
        escrow.setRateLimit(mockToken, 10_000);
    }

    function test_dailyUsage() public {
        address mockToken = makeAddr("mockToken");

        // Daily usage should not be updated if rate limit is 0
        vm.mockCall(mockToken, IERC20.transferFrom.selector, abi.encode(true));
        escrow.send(IERC20(mockToken), address(this), 10_000);
        assertEq(escrow.getDailyUsage(mockToken), 0);

        escrow.recvCallback(mockToken, address(this), 10_000);
        assertEq(escrow.getDailyUsage(mockToken), 0);

        // Set rate limit and check daily usage
        vm.prank(rateLimiter);
        escrow.setRateLimit(mockToken, 100_000);
        assertEq(escrow.getRateLimit(mockToken), 100_000);

        vm.mockCall(mockToken, IERC20.transferFrom.selector, abi.encode(true));
        escrow.send(IERC20(mockToken), address(this), 1020);
        assertEq(escrow.getDailyUsage(mockToken), 1020);

        escrow.recvCallback(mockToken, address(this), 20);
        assertEq(escrow.getDailyUsage(mockToken), 1000);

        // Next day
        vm.warp(block.timestamp + 1 days);
        assertEq(escrow.getDailyUsage(mockToken), 0);

        vm.mockCall(mockToken, IERC20.transferFrom.selector, abi.encode(true));
        escrow.send(IERC20(mockToken), address(this), 100_000);
        assertEq(escrow.getDailyUsage(mockToken), 100_000);

        escrow.recvCallback(mockToken, address(this), 100_000);
        assertEq(escrow.getDailyUsage(mockToken), 0);

        // next day
        vm.warp(block.timestamp + 1 days);
        assertEq(escrow.getDailyUsage(mockToken), 0);

        escrow.recvCallback(mockToken, address(this), 100_000);
        assertEq(escrow.getDailyUsage(mockToken), 0);

        vm.mockCall(mockToken, IERC20.transferFrom.selector, abi.encode(true));
        escrow.send(IERC20(mockToken), address(this), 100_000);
        assertEq(escrow.getDailyUsage(mockToken), 100_000);

        escrow.recvCallback(mockToken, address(this), 150_000);
        assertEq(escrow.getDailyUsage(mockToken), 0);
    }

    /// forge-config: default.fuzz.runs = 256
    function testFuzz_rateLimit(uint8 n) public {
        vm.assume(1 < n);

        address mockToken = makeAddr("mockToken");
        uint256 sendAmount = 10_000;
        uint256 rateLimit = sendAmount * n - 1;

        vm.prank(rateLimiter);
        escrow.setRateLimit(mockToken, rateLimit);

        for (uint256 i = 0; i < n - 1; ++i) {
            vm.mockCall(mockToken, IERC20.transferFrom.selector, abi.encode(true));
            escrow.send(IERC20(mockToken), address(this), sendAmount);
            assertEq(escrow.getDailyUsage(mockToken), sendAmount * (i + 1));
        }

        vm.mockCall(mockToken, IERC20.transferFrom.selector, abi.encode(true));
        vm.expectRevert(abi.encodeWithSelector(RateLimitErrors.RateLimitExceeded.selector, rateLimit, sendAmount * n));
        escrow.send(IERC20(mockToken), address(this), sendAmount);
    }

    /// forge-config: default.fuzz.runs = 256
    function testFuzz_sendBackAndForth(uint8 n) public {
        vm.assume(1 < n);

        address mockToken = makeAddr("mockToken");
        uint256 sendAmount = 10_000;
        uint256 rateLimit = sendAmount + 1;

        vm.prank(rateLimiter);
        escrow.setRateLimit(mockToken, rateLimit);

        for (uint256 i = 0; i < n; ++i) {
            vm.mockCall(mockToken, IERC20.transferFrom.selector, abi.encode(true));
            escrow.send(IERC20(mockToken), address(this), sendAmount);
            assertEq(escrow.getDailyUsage(mockToken), sendAmount);

            escrow.recvCallback(mockToken, address(this), sendAmount);
            assertEq(escrow.getDailyUsage(mockToken), 0);
        }
    }

    function test_noMidnightBypass() public {
        address mockToken = makeAddr("mockToken");
        uint256 rateLimit = 10_000;

        vm.prank(rateLimiter);
        escrow.setRateLimit(mockToken, rateLimit);

        vm.mockCall(mockToken, IERC20.transferFrom.selector, abi.encode(true));

        escrow.send(IERC20(mockToken), address(this), rateLimit);
        assertEq(escrow.getDailyUsage(mockToken), rateLimit);

        // Warp just enough that a calendar-day boundary would be crossed in the old system
        vm.warp(block.timestamp + 30 minutes);

        // Only ~208 has decayed (10000 * 30min / 1day)
        uint256 decayed = rateLimit * 30 minutes / 1 days;

        // Sending more than the decayed portion should revert
        vm.expectRevert(abi.encodeWithSelector(RateLimitErrors.RateLimitExceeded.selector, rateLimit, rateLimit + 1));
        escrow.send(IERC20(mockToken), address(this), decayed + 1);

        // The decayed amount can be sent
        escrow.send(IERC20(mockToken), address(this), decayed);
        assertEq(escrow.getDailyUsage(mockToken), rateLimit);

        // After the full period, the full limit is available again
        vm.warp(block.timestamp + 1 days);
        escrow.send(IERC20(mockToken), address(this), rateLimit);
        assertEq(escrow.getDailyUsage(mockToken), rateLimit);
    }

    function test_refund_succeeds_when_limit_saturated() public {
        address mockToken = makeAddr("mockToken");
        uint256 rateLimit = 10_000;

        vm.prank(rateLimiter);
        escrow.setRateLimit(mockToken, rateLimit);

        // Saturate the limit
        vm.mockCall(mockToken, IERC20.transfer.selector, abi.encode(true));
        vm.mockCall(mockToken, IERC20.transferFrom.selector, abi.encode(true));
        escrow.send(IERC20(mockToken), address(this), rateLimit);
        assertEq(escrow.getDailyUsage(mockToken), rateLimit);

        // Attempting to send another token should fail/revert
        vm.expectRevert(abi.encodeWithSelector(RateLimitErrors.RateLimitExceeded.selector, rateLimit, rateLimit + 1));
        escrow.send(IERC20(mockToken), address(this), 1);

        // Refund should succeed even though the daily limit is saturated
        escrow.sendRefund(IERC20(mockToken), address(this), 5000, 5000);

        // Daily usage increases by 5_000 uncapped (restoring what the deposit removed)
        assertEq(escrow.getDailyUsage(mockToken), rateLimit + 5000);
    }

    function test_refund_symmetry() public {
        address mockToken = makeAddr("mockToken");
        uint256 rateLimit = 10_000;

        vm.prank(rateLimiter);
        escrow.setRateLimit(mockToken, rateLimit);

        // 1. Initial usage = 0
        assertEq(escrow.getDailyUsage(mockToken), 0);

        // 2. Deposit (tokens enter escrow, lowering usage, i.e. usage = usage - amount)
        // Since usage is currently 0, reducing it keeps it at 0 (floored)
        escrow.recvCallback(mockToken, address(this), 2000);
        assertEq(escrow.getDailyUsage(mockToken), 0);

        // Let's first consume some rate limit to avoid floor-at-zero effects
        vm.mockCall(mockToken, IERC20.transfer.selector, abi.encode(true));
        vm.mockCall(mockToken, IERC20.transferFrom.selector, abi.encode(true));
        escrow.send(IERC20(mockToken), address(this), 5000);
        assertEq(escrow.getDailyUsage(mockToken), 5000);

        // Deposit 2_000 tokens
        escrow.recvCallback(mockToken, address(this), 2000);
        assertEq(escrow.getDailyUsage(mockToken), 3000);

        // Refund 2_000 tokens
        escrow.sendRefund(IERC20(mockToken), address(this), 2000, 2000);

        // Daily usage should be restored to 5_000 (pre-deposit value)
        assertEq(escrow.getDailyUsage(mockToken), 5000);
    }

    function test_refund_restores_only_usage_removed_by_deposit() public {
        address mockToken = makeAddr("mockToken");
        uint256 rateLimit = 10_000;

        vm.prank(rateLimiter);
        escrow.setRateLimit(mockToken, rateLimit);

        vm.mockCall(mockToken, IERC20.transfer.selector, abi.encode(true));
        escrow.send(IERC20(mockToken), address(this), 1000);
        assertEq(escrow.getDailyUsage(mockToken), 1000);

        // The deposit is larger than current usage, so only 1,000 units are removed.
        uint256 usageRemoved = escrow.recvCallback(mockToken, address(this), 2000);
        assertEq(usageRemoved, 1000);
        assertEq(escrow.getDailyUsage(mockToken), 0);

        escrow.sendRefund(IERC20(mockToken), address(this), 2000, usageRemoved);
        assertEq(escrow.getDailyUsage(mockToken), 1000);
    }

    function test_refund_credit_is_zero_when_limit_is_disabled() public {
        address mockToken = makeAddr("mockToken");

        // With no configured limit, deposits remove no tracked usage and refunds restore none.
        uint256 usageRemoved = escrow.recvCallback(mockToken, address(this), 2000);
        assertEq(usageRemoved, 0);
        vm.mockCall(mockToken, IERC20.transfer.selector, abi.encode(true));
        escrow.sendRefund(IERC20(mockToken), address(this), 2000, usageRemoved);
        assertEq(escrow.getDailyUsage(mockToken), 0);
    }
}
