# `CosmWasm` ICS-08 Ethereum Light Client

This program is a `CosmWasm` implementation of an Ethereum light client, designed for integration with [`ibc-go`](https://github.com/cosmos/ibc-go)'s [`08-wasm`](https://github.com/cosmos/ibc-go/tree/main/modules/light-clients/08-wasm) client wrapper. It manages IBC client and consensus states, delegating all light client-specific logic to the core module in [`packages/ethereum/light-client`](../../packages/ethereum/light-client).

## Light Client Functions

### Update Client Logic

The `update_client` functionality is based on the [Ethereum Light Client Sync Protocol](https://github.com/ethereum/consensus-specs/blob/dev/specs/altair/light-client/sync-protocol.md).

Every header must be anchored at the client's latest slot: `trusted_slot` has to equal `latest_slot`, or the update fails with `TrustedSlotNotLatest`. Historical updates are not supported. An older consensus state is refused as an anchor because the keys of a past sync committee carry no stake any more and could forge a root.

Each update stores the execution state root, the timestamp and the slot of the finalized header.

### Membership Proofs

The client verifies Merkle-Patricia Trie (MPT) inclusion and exclusion proofs for IBC commitments held by the `ICS26Router` contract on Ethereum (`ibc_contract_address` in the client state; the commitments live in its `IBCStoreUpgradeable` storage). Each proof carries an account proof of that contract against the stored execution state root, followed by a storage proof of the commitment.

### Misbehaviour Handling

Misbehaviour is the submission of two conflicting light client updates for the same block height, each valid on its own. When the client detects it, it **freezes** and rejects every later update.

A frozen client is unfrozen only by a governance `MsgMigrateContract` with the `reinstantiate` migration, which deletes every stored consensus state and writes a new bootstrap. See [Migrating the Ethereum 08-wasm client on Cosmos](../../docs/runbooks/PRODUCTION_DEPLOYMENT.md#migrating-the-ethereum-08-wasm-client-on-cosmos). The IBC sudo messages `migrate_client_store` (substitute-client recovery) and `verify_upgrade_and_update_state` are not supported and return `UnsupportedLifecycleOperation`.

## Deployment

### Build Requirements

- [Rust](https://rustup.rs/)
- [Just](https://just.systems/man/en/)
- [Protobuf compiler](https://grpc.io/docs/protoc-installation/)
- [Docker](https://docs.docker.com/get-docker/)

We use the [`cosmwasm/optimizer`](https://github.com/CosmWasm/optimizer) Docker image to build the `CosmWasm` contracts. Ensure you have Docker installed and running.

```bash
just build-cw-ics08-wasm-eth
```

### Storing the Wasm Binary

Once built, the wasm binary needs to be stored on the Cosmos chain through governance. On the local devnet, `scripts/local/wasm.sh` submits and passes the proposal and prints the checksum. On any other chain, follow the [IBC-Go documentation](https://ibc.cosmos.network/v10/ibc/light-clients/wasm/governance/#storing-new-wasm-light-client-byte-code).

### Creating the Light Client

The relayer creates the client on the Cosmos chain where the bytecode is stored. `--wasm-checksum` is the checksum of the wasm bytecode in hex:

```bash
./relayer create-clients-cosmos --config config.json --wasm-checksum <hex>
```

The full sequence is in [docs/runbooks/E2E_LOCAL.md](../../docs/runbooks/E2E_LOCAL.md#local-cosmos--ethereum-e2e).

## Acknowledgements

This implementation builds on the Ethereum light client work by [Union Labs](https://github.com/unionlabs/union).
