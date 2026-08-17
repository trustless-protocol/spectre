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

`RATE_LIMITER_ROLE` is global: it can configure every escrow, including an
escrow created after deployment. Production deployment enables a permanent
launch gate and, for each configured client, creates its escrow, installs
non-zero limits for each configured ERC-20, then activates the escrow. Both
deployment and verification reject zero-address or non-contract token entries.

For a later client, preserve the same order: `createEscrow(clientId)` through
governance, install non-zero limits from the rate-limiter account, then call
`activateEscrow(clientId, tokens)` through governance. Packet processing
rejects a missing or inactive escrow, so the client cannot operate before its
listed tokens have the configured limits.

Activation is a point-in-time check of only the supplied token list. Removing a
limit later by setting it to `0` does not deactivate the escrow, and tokens not
included in the activation call remain uncapped. Governance and the rate-limiter
must therefore treat later token additions and limit changes as separate safety
decisions.

When upgrading an existing live deployment, activate every existing escrow with
its required token list **before** calling `enableEscrowLaunchGate()`. The active
mapping defaults to `false`, so enabling the permanent gate first would halt
packet processing for those clients. `activateEscrow` is available before the
gate is enabled, which permits this safe migration order.

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
The proposal must preserve the client's existing counterparty client ID and
Merkle prefix; changing either binding is rejected before the delay begins.

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

## Escrow upgrade compatibility

The token-layer correctness upgrade changes the `Escrow.recvCallback` and
`Escrow.sendRefund` ABI. Upgrade the Escrow beacon before, or atomically with,
the ICS20Transfer implementation; activating the new transfer implementation
against the old escrow implementation makes sends revert.

Refund-credit records exist only for packets sent after this upgrade. Refunds
of packets already in flight restore zero rate-limit usage, so operators must
account for that one-time transition when monitoring escrow limits.
