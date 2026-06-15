# Design Conventions

## Solidity

### Protobuf Encoding (`Encode.sol`)

All encoding functions must match CometBFT's `proto.Marshal()` exactly:
- **Proto3 rule**: Skip zero-value fields (varint 0, empty bytes, empty string)
- **Nested messages**: Use proper wire format with tag + length prefix
- **cdcEncode wrappers**: Header fields are wrapped in `StringValue`, `Int64Value`, `BytesValue`, `Timestamp` before hashing
- **Merkle prefix**: Use `bytes1(0x00)` / `bytes1(0x01)` (1-byte), NOT `[0x00]` (32-byte `uint256[1]`)

Cross-validate any encoding changes via:
```bash
forge test --match-contract EncodeTest -vvv  # Solidity output must match the Go-reference hex baked into the fixtures
```

### Contract Patterns

- **UUPS proxy**: Core contracts (ICS26Router, ICS20Transfer, Groth16ICS07Tendermint)
- **Beacon proxy**: Instance contracts (IBCERC20, Escrow) — ICS20Transfer upgrades all atomically
- **Access control**: OpenZeppelin `AccessManager` with roles in `IBCRolesLib.sol`
- **Error handling**: Custom errors defined in `contracts/errors/` interfaces
- **Storage layout**: `IBCStoreUpgradeable` + `ICS24Host` for deterministic storage slots

### Formatting

Configured in `foundry.toml`:
- Line length: 120
- Tab width: 4
- Quote style: double quotes
- Bracket spacing: true

### Linting

- **solhint**: Max cyclomatic complexity 8, function max lines 70 (warning)
- **natlint**: NatSpec documentation required
- **Slither**: Security static analysis, excludes low/informational

## Go (Relayer)

### Module Dependencies

`relayer/go.mod` uses `replace` directives for local paths:
- `ecip-gnark` → `../../ecip-gnark`
- `gnark` → `../../decentrio-gnark`

These must be adjusted per developer's local setup.

### Contract Bindings

After modifying Solidity contracts the relayer depends on:
```bash
bun install && forge build
abigen --abi <abi_json> --pkg <PkgName> --out relayer/bindings/<PkgName>/binding.go
```

Shared bindings also live in `packages/go-abigen/`.

### Configuration

Relayer uses JSON config file (see `relayer/config.example.json`):
- Top-level `server` block: log_level, address, port
- `modules` array, each entry has `name`, `src_chain`, `dst_chain`, and `config`:
  - `cosmos_to_eth` config: tm_rpc_url, ics26_address, eth_rpc_url, ics07_client, wrapper_verifier, membership, misbehaviour, update_client, and optional fields: fetch_timeout, trusting_period, trust_level, proof_type
  - `eth_to_cosmos` config: tm_rpc_url, ics26_address, eth_rpc_url, eth_beacon_api_url, signer_address

Secrets (private keys, prover paths) stay in `.env` file.

### Error Handling

- Relayer uses `log.Fatal` for unrecoverable errors (process exits)
- `start` command uses batch processing via `services.StartLoop()` with size/time thresholds
- Per-packet errors logged without crashing the relay loop
- Timeout checking: packets past their timeout are skipped before submission
- Transaction handlers retry with re-queried account sequence on nonce errors

## Rust

### Workspace Structure

25 crate members in root `Cargo.toml`. Key packages:
- `packages/ethereum/` — Ethereum light client for CosmWasm
- `packages/tendermint-light-client/` — Tendermint proof verification

### Linting

- `cargo fmt` — formatting
- `cargo clippy` — lint with warnings as errors in CI

## Naming Conventions

| Context | Convention | Example |
|---------|-----------|---------|
| Solidity contracts | PascalCase | `ICS20Transfer` |
| Solidity functions | camelCase | `encodeValidator` |
| Solidity errors | PascalCase | `InsufficientSignersOverlap` |
| Solidity events | PascalCase | `SendPacket` |
| Go packages | lowercase | `prover`, `subscriber` |
| Go functions | PascalCase (exported) | `ExtractValidatorSignature` |
| Go test files | `*_test.go` | `prover_test.go` |
| Rust crates | kebab-case | `tendermint-light-client` |

## Access Control Roles

Defined in `contracts/utils/IBCRolesLib.sol`:

| Role | Value | Purpose |
|------|-------|---------|
| RELAYER_ROLE | 1 | Submit IBC packets |
| PAUSER_ROLE | 2 | Emergency pause |
| UNPAUSER_ROLE | 3 | Resume after pause |
| DELEGATE_SENDER_ROLE | 4 | Send on behalf of others |
| RATE_LIMITER_ROLE | 5 | Set transfer rate limits |
