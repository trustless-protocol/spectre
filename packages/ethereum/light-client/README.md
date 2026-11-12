# `ethereum-light-client`

This is the stateless verification implementation of the Ethereum light client. It contains all the core logic for verifying Ethereum consensus, proving state (membership and non-membership) and the headers submitted to update the light client.
The state is handled by the `CosmWasm` implementation in `programs/cw-ics08-wasm-eth`.

## Acknowledgements

This work is based on the Ethereum light client created by [Union](https://github.com/unionlabs/union/)
