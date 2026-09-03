# Base attestor-trusted L2 client

`MsgUpdateClient` is permissionless, but every accepted header must carry a valid Ed25519 signature
from the public key pinned in the immutable Base attestor profile. The signature binds the Base chain
ID and complete L2 block identity; the client then verifies the router account proof against that
signed header state root. See [docs/L2_CLIENTS.md](../../docs/L2_CLIENTS.md).
