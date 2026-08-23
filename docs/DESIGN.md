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

- **UUPS proxy**: Core contracts (ICS26Router, ICS20Transfer, SpectreClient)
- **Beacon proxy**: Instance contracts (IBCERC20, Escrow) — ICS20Transfer upgrades all atomically
- **Access control**: OpenZeppelin `AccessManager` with roles in `IBCRolesLib.sol`
- **Error handling**: Custom errors owned by `contracts/core/errors/`, `contracts/apps/ics20/errors/`, and `contracts/light-clients/spectre/errors/`
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
- `modules` array, each entry has `name`, `src_chain`, `dst_chain`, and `config`:
  - `cosmos_to_eth` config: tm_rpc_url, ics26_address, eth_rpc_url, spectre_client, signature_verifier, membership, misbehaviour, update_client, and optional fields: fetch_timeout, trusting_period, trust_level, proof_type. **One entry per Cosmos source** (distinct `ics26_client_id`); `start` runs an independent relay loop for each, and `create-clients-cosmos --source <ics26_client_id>` (then `create-clients-eth --source ...`) targets one.
  - `eth_to_cosmos` config: tm_rpc_url, ics26_address, eth_rpc_url, eth_beacon_api_url, signer_address
  - `cosmos_to_l2` config (dst_chain `opstack`/`arbitrum`): same shape as `cosmos_to_eth`, reusing its struct — **`eth_rpc_url`/`eth_ws_url` here point at the L2's own RPC**, not L1. To pair with an `l2_to_cosmos` module for the timeout return path, its `eth_rpc_url` must resolve to the same chain id as that module's `l2_rpc_url` (they may still be different endpoints).
  - `l2_to_cosmos` config (src_chain `opstack`/`arbitrum`): l2_rpc_url, tm_rpc_url, attestor_addr, attestor_src_chain, l2_wasm_client_id, l2_ics26_client_id, head_kind, rollup_profile. No `l1_rpc_url` or `eth_beacon_api_url`: the L2 client verifies nothing against L1, so this module never dials one (#347, `relayer/cmd/build_l2_source.go`). The return path back to the paired `cosmos_to_l2` dest is resolved at startup by matching L2 chain id + tm_rpc_url + L2 router address (`relayer/cmd/l2_timeout_return_path.go`) — no explicit link field needed. Resolution is skipped entirely when no `l2_to_cosmos` source is configured.

Secrets (private keys, prover paths) stay in `.env` file.

### Error Handling

- Relayer uses `log.Fatal` for unrecoverable errors (process exits)
- `start` command uses batch processing via the chain-adapter engine `runAdapterEngine` (`relayer/cmd/run_adapters.go`, one `relay.Module` per direction) with size/time thresholds
- Per-packet errors logged without crashing the relay loop
- Timeout checking: packets past their timeout are skipped before submission
- Transaction handlers retry with re-queried account sequence on nonce errors

## Rust

### Workspace Structure

Key Rust workspace packages:
- `packages/ethereum/` — Ethereum light client for CosmWasm

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
| Rust crates | kebab-case | `ethereum-light-client` |

## Access Control Roles

Defined in `contracts/shared/access/IBCRolesLib.sol`:

| Role | Value | Purpose |
|------|-------|---------|
| RELAYER_ROLE | 1 | Submit IBC packets |
| PAUSER_ROLE | 2 | Emergency pause |
| UNPAUSER_ROLE | 3 | Resume after pause |
| DELEGATE_SENDER_ROLE | 4 | Send on behalf of others |
| RATE_LIMITER_ROLE | 5 | Set transfer rate limits |
