# Arbitrum profile data

These JSON files are reviewed tooling/test inputs. They are not loaded with `include_str!` and are
not compiled into the Wasm artifact. Replace `ethereum_client.client_id` and
`ethereum_client.wasm_checksum` with the host deployment's pinned Ethereum client when creating a
client.

`Profile.protocol` retains a tagged shape for stable client configuration:

- `bold_v2` pins the `_assertions` mapping and packed status offset.

The Arbitrum Sepolia profile pins the BoLD RollupCore and its reviewed
`_assertions` layout.
