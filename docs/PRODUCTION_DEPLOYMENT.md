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

`migrateClient` is scoped by
`LIGHT_CLIENT_MIGRATOR_ROLE_<clientId>`. The role cannot be pre-granted during
bootstrap because the client ID does not yet exist; provision it only for the
specific client during an approved migration procedure and revoke it after use.
