# Security

## Security Model

### Trust Assumptions

1. **Validator set**: At least 2/3 of voting power is honest (standard BFT assumption)
2. **Groth16 circuit**: PreHashCircuit correctly verifies Ed25519 signatures
3. **Admin keys**: UUPS upgrades require admin access via AccessManager
4. **Relayer**: Untrusted — relayer cannot forge proofs, only relay valid state transitions

### Proof Verification Chain

```
Cosmos validator signs block
    → Ed25519 signature verified via Groth16 ZK proof
    → WrapperVerifier decompresses points + computes SHA512
    → Groth16Verifier checks proof against verification key
    → Groth16ICS07Tendermint updates client state
    → ICS-23 Merkle proofs verify packet commitments
```

### Attack Surfaces

| Surface | Mitigation |
|---------|-----------|
| Forged headers | Groth16 proof verification (cryptographic) |
| Replay attacks | Packet sequencing in ICS26Router |
| Double-spend | Commitment storage in ICS24Host |
| Validator equivocation | Misbehaviour detection → client freeze |
| Unbounded token minting | Rate limiting (RateLimitUpgradeable) |
| Malicious upgrade | AccessManager RBAC + UUPS auth |

## Static Analysis

```bash
just slither    # Run Slither security scanner
```

Config (`.slither.config.json`):
- Framework: Foundry
- Excludes: informational, low severity, dependencies

## Access Control

Roles defined in `IBCRolesLib.sol`:
- Upgrades require admin (AccessManager)
- Emergency controls via PAUSER/UNPAUSER roles
- Rate limits managed by RATE_LIMITER_ROLE

### Proof Submission (Groth16ICS07Tendermint)

`PROOF_SUBMITTER_ROLE` controls who can call `verifyMembership`/`updateClient`:
- **`roleManager = address(0)`**: Anyone can submit proofs (permissionless, suitable for dev/test). Proofs are still cryptographically verified by Groth16.
- **`roleManager != address(0)`**: Only addresses granted `PROOF_SUBMITTER_ROLE` by the role manager can submit. Use for production to restrict relay to authorized operators.

## Known Considerations

- `Groth16Verifier.sol` is auto-generated from circuit VK — do not manually edit
- WrapperVerifier performs Ed25519 point decompression on-chain — gas-intensive but necessary
- Shadowfork tests (`test/shadowfork/`) test against real mainnet/testnet state
- `relayer/go.mod` replace directives point to local paths — verify before building
