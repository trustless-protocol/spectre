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
cd operator && go run ./cmd/encode_debug/  # Go reference hex
forge test --match-contract EncodeTest -vvv # Solidity must match
```

### Contract Patterns

- **UUPS proxy**: Core contracts (ICS26Router, ICS20Transfer, SP1ICS07Tendermint)
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

## Go (Operator)

### Module Dependencies

`operator/go.mod` uses `replace` directives for local paths:
- `ecip-gnark` → `../../ecip-gnark`
- `gnark` → `../../decentrio-gnark`

These must be adjusted per developer's local setup.

### Contract Bindings

After modifying Solidity contracts the operator depends on:
```bash
bun install && forge build
abigen --abi <abi_json> --bin <bin_hex> --pkg <PkgName> --out operator/bindings/<PkgName>/binding.go
```

ABI and Bin bytecode must stay in sync.

### Error Handling

- Operator uses `log.Fatal` for unrecoverable errors (process exits)
- Service loops return errors up the call chain
- Transaction handlers retry with re-queried account sequence on nonce errors

## Rust

### Workspace Structure

25 crate members in root `Cargo.toml`. Key packages:
- `packages/relayer/` — multi-chain relay modules
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
