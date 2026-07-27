# Arbitrum profile data

These JSON files are reviewed tooling/test inputs. They are not loaded with `include_str!` and are
not compiled into the Wasm artifact. Replace `ethereum_client.client_id` and
`ethereum_client.wasm_checksum` with the host deployment's pinned Ethereum client when creating a
client.

`Profile.protocol` is a tagged union:

- `bold_v2` pins the `_assertions` mapping and packed status offset.
- `legacy_nitro` pins the packed node-lifecycle slot, `_nodes` mapping, lifecycle offsets, and
  `Node.confirmData` offset.

The canonical Arbitrum Sepolia profile uses `legacy_nitro`. Its reviewed layout is lifecycle slot
`0x75`, `_nodes` slot `0x76`, lifecycle offsets `0/8/16`, and `confirmData` offset `2`.
