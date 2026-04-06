# Product Sense

## Domain Model

### IBC v2 (Inter-Blockchain Communication)

IBC is a protocol for trustless cross-chain communication. This implementation enables Ethereum ↔ Cosmos token transfers with ZK proof verification.

**Key entities:**
- **Client**: Light client state tracking a remote chain (Tendermint headers on Ethereum, Ethereum headers on Cosmos)
- **Packet**: A unit of cross-chain data (source/destination ports, sequence, payload)
- **Channel**: Logical connection between two chains identified by client IDs
- **Commitment**: On-chain proof that a packet was sent (stored in ICS24Host)

### Token Transfer Flow (ICS-20)

**Send (lock/burn):**
1. User calls `ICS20Transfer.sendTransfer()`
2. Native tokens → locked in `Escrow`; IBC tokens → `IBCERC20.burn()`
3. Packet commitment stored in `ICS26Router`

**Receive (mint/unlock):**
1. Relayer submits `recvPacket()` with membership proof
2. Light client verifies counterparty state
3. New chain tokens → `IBCERC20.mint()`; returning tokens → `Escrow.release()`

**Denomination tracing:**
- Tokens crossing chains get prefixed: `transfer/channel-0/uatom`
- Returning to origin chain removes prefix and unlocks original tokens

### Proof Aggregation

Multiple packets can be batched into a single proof submission:
- `multiRecvPacket` / `multiAckPacket` in ICS20Transfer
- ~90% gas reduction per packet for large batches (25-50 packets)

## API Surface

### On-chain (Solidity)

**ICS26Router** — IBC packet lifecycle:
- `sendPacket()` — record outbound packet commitment
- `recvPacket()` — verify and deliver inbound packet
- `ackPacket()` — process acknowledgement
- `timeoutPacket()` — handle expired packets

**ICS20Transfer** — token bridge:
- `sendTransfer()` — initiate cross-chain transfer
- `sendTransferWithSender()` — send on behalf of another address (requires DELEGATE_SENDER_ROLE)
- `multiRecvPacket()` / `multiAckPacket()` — batched operations

**Groth16ICS07Tendermint** — light client:
- `updateClient()` — update with new Tendermint header + Groth16 proof
- `membership()` / `nonMembership()` — verify ICS-23 Merkle proofs

### Off-chain (Go Relayer)

**CLI commands:**
- `relayer start --config config.json` — relay loop (Cosmos → ETH)
- `relayer create-clients --config config.json` — one-time light client setup
- `relayer genesis` — generate genesis state
- `relayer fixtures membership` — verify membership proof on-chain

**JSON config** (`relayer/config.example.json`):
- `cosmos_to_eth` module: tm_rpc_url, ics26_address, eth_rpc_url, ics07_client, wrapper_verifier, membership, misbehaviour, update_client
- `eth_to_cosmos` module: tm_rpc_url, ics26_address, eth_rpc_url, eth_beacon_api_url, signer_address

### Environment Variables

Key variables from `.env`:
- `ETH_PRIVATE_KEY`: Ethereum signer for relay transactions
- `COSMOS_PRIVATE_KEY`: Cosmos signer for MsgCreateClient
- `PROVER_R1CS_PATH`, `PROVER_PK_PATH`, `PROVER_VK_PATH`: Groth16 circuit artifacts
- `ETH_RPC_URL`: Ethereum RPC for shadowfork tests
- `TENDERMINT_RPC_URL`: CometBFT node endpoint (used by genesis/fixtures commands)

## Business Rules

1. **Trust threshold**: Configurable fraction of validator overlap needed to update header (default 2/3)
2. **Trusting period**: Duration (in seconds) during which submitted headers are valid
3. **Unbonding period**: Staking unbonding period — header updates rejected if trusting period exceeded
4. **Rate limiting**: Per-token transfer rate limits (configurable by RATE_LIMITER_ROLE)
5. **Emergency pause**: PAUSER_ROLE can halt all operations; UNPAUSER_ROLE resumes
6. **Frozen state**: Client freezes on misbehaviour detection (two conflicting headers at same height)
