# Production deployment

`ProductionDeploy.s.sol` bootstraps the IBC proxies under an `AccessManager` and
removes the bootstrap account before returning the deployment addresses.
`ProductionVerify.s.sol` must be run against the same addresses and environment
configuration after deployment.

## Governance requirements

`GOVERNANCE_ADMIN` must be an OpenZeppelin `TimelockController`, or a compatible
contract exposing `getMinDelay()`, with a minimum delay at least equal to
`SECURITY_DELAY` (two days by default). The AccessManager `ADMIN_ROLE` has an
execution delay of zero by design, so the governance timelock is what protects
role and selector changes made by the administrator.

The `UPGRADER_ACCOUNT`, `UNPAUSER_ACCOUNT`, and `ID_CUSTOMIZER_ROLE` use the
AccessManager delay. The relayer, pausers, and misbehaviour watcher are
immediate operational roles.

## Escrow launch gate

This deployment does not claim that `RATE_LIMITER_ROLE` can operate newly
created escrows. Escrows are BeaconProxy instances created per client, and
their `setRateLimit` selector must be mapped after each escrow is created.
Complete the TK-01/#293 escrow launch-gate deployment before using rate limits
in production; that flow pre-creates configured escrows, installs their limits,
maps the selector, and verifies the result.

## Light-client provisioning

The documented `./relayer create-clients-eth` flow submits `addClient` directly
from the relayer key. A production deployment grants `addClient` to governance
with an AccessManager delay, so that command is not a complete production
provisioning flow. Deploy the light-client implementation separately, then
schedule and execute the `addClient` calldata through the governance and
AccessManager timelocks before writing the resulting client address to relayer
configuration.

### Governed client migration

Migration is scoped by `LIGHT_CLIENT_MIGRATOR_ROLE_<clientId>`, which cannot be
granted during bootstrap because the client ID does not yet exist.

1. Read the role ID from `getLightClientMigratorRole(clientId)`.
2. Through governance, grant that role to the approved migrator with a
   **non-zero** AccessManager execution delay. A zero-delay grant cannot create a
   proposal.
3. The granted migrator calls `proposeClientMigration` directly with the exact
   client ID, counterparty information, and replacement light-client address.
4. Verify the emitted proposal or `getClientMigration(clientId)`, then wait
   until `executeAfter`.
5. During the execution window, any account may submit the exact committed
   values to `executeClientMigration`.
6. Revoke the per-client role after the approved migration is complete.

The contract enforces a wait of `max(role execution delay, 48 hours)` and
re-checks at execution that the original proposer still holds the role.

`migrateClient` remains an ABI-compatible alias for
`executeClientMigration`; it executes a mature matching proposal rather than
failing unconditionally. New operator tooling should use the explicit
`executeClientMigration` name.

Cancellation requires `PAUSER_ROLE` with AccessManager execution delay
**exactly zero**. Such a pauser may call `cancelClientMigration` even while the
router target is closed.

An authorized migrator can occupy a client ID with an unwanted proposal for the
full proposal window. Revoking that migrator blocks execution, but does not free
the proposal slot until it expires. For an unwanted proposal, revoke the role
and have a zero-delay pauser cancel it before granting a replacement migrator.
