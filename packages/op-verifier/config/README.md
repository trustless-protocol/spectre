# OP profile data

These JSON files are reviewed tooling/test inputs. They are not loaded with `include_str!` and are
not compiled into the Wasm artifact. Replace `ethereum_client.client_id` and
`ethereum_client.wasm_checksum` with the host deployment's pinned Ethereum client when creating a
client.

`l2_header_fork` pins the canonical execution-header field set accepted for this deployment. A
fork transition requires a reviewed profile and client migration; it does not require rebuilding
the verifier Wasm.
