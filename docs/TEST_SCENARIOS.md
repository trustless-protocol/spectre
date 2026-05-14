# Test Scenarios

The 8-step flow in the README covers the happy path: compile circuits → start nodes → create clients → relay → send transfer → verify balance. This document describes additional scenarios worth testing.

```bash
# 1. (One-time) compile per-bucket circuits + emit Groth16Verifier_N{N}.sol.
#    Re-run only when circuit code changes. After this, redeploy contracts.
cd relayer
go run ./prover/cmd ./bin ../contracts/verifiers

# 2. Build the relayer binary
go build -o relayer ./cmd

# 3. Start Ethereum first and wait until the beacon node finalizes.
#    Replace 56246 with your Kurtosis-mapped beacon RPC port.
kurtosis enclave rm -f my-testnet
./run_eth_node.sh        # Kurtosis Ethereum testnet + deploys core contracts
# Poll until finalized.epoch > 0:
curl -s <eth_beacon_api_ur>/eth/v1/beacon/states/head/finality_checkpoints
curl -s http://127.0.0.1:59717/eth/v1/beacon/states/head/finality_checkpoints

# 4. Then start Cosmos and submit the Ethereum LC WASM.
#    This requires a wasm-enabled Cosmos binary (08-wasm), e.g. simd:
#    COSMOS_BIN=simd ./run_cosmos_node.sh
#    COSMOS_BIN=simd ./wasm.sh
#    If you only have stock gaiad, use the container flow instead:
#    ./run_cosmos_node_docker.sh
#    ./wasm_docker.sh
#    The container flow now copies
#    e2e/interchaintestv8/wasm/cw_ics08_wasm_eth.wasm.gz into the simd container
#    and runs `simd tx ibc-wasm store-code` directly, instead of generating an
#    older proposal.json payload.
./run_cosmos_node.sh
./wasm.sh

# 5. Deploy Tendermint light client on Ethereum.rela
#    Copies the ICS07 address back into relayer/config.json automatically.
./relayer create-clients \
  --config config.json \
  --wasm-checksum <hex-from-wasm.sh>

./relayer create-clients  --config config.example.json --wasm-checksum 0xd24688886ed8cec00c667fa69c173fbab9a08c75900ce18afe10517c82e55592

# 6. Start the bi-directional relay loop
./relayer start --config config.example.json

# 7. send tx

ABS_TIMEOUT=$(($(date +%s) + 2000))

gaiad tx ibc-transfer transfer transfer 08-wasm-0 0x8943545177806ed17b9f23f0a21ee5948ecaa776 1000stake \
  --from test1 \
  --home /Users/donglieu/.gaia \
  --chain-id test-ibc-eth \
  --node tcp://127.0.0.1:26657 \
  --keyring-backend test \
  --gas-prices 1stake \
  --absolute-timeouts \
  --packet-timeout-timestamp "$ABS_TIMEOUT" \
  --generate-only \
| jq '.body.messages[0].encoding = "application/x-solidity-abi"' \
| gaiad tx sign /dev/stdin \
    --from test1 \
    --home /Users/donglieu/.gaia \
    --chain-id test-ibc-eth \
    --keyring-backend test \
| gaiad tx broadcast /dev/stdin \
    --node tcp://127.0.0.1:26657 \
    -y

# timeoutTimestamp in this setup is interpreted as unix seconds, not nanoseconds.


# 8. check
cast call 0xee0fcb8e5ccad0b4197baabd633333886f5c364d \
  'ibcERC20Contract(string)(address)' \
  'transfer/cosmoshub-1/stake' \
  --rpc-url http://127.0.0.1:59619

cast call 0x016f5f33DbCb653e6393698Beba9DC19d828D75e \
  'fullDenomPath()(string)' \
  --rpc-url http://127.0.0.1:59619

cast call 0x016f5f33DbCb653e6393698Beba9DC19d828D75e \
  'balanceOf(address)(uint256)' \
  0x8943545177806ed17b9f23f0a21ee5948ecaa776 \
  --rpc-url http://127.0.0.1:59619

cast call 0x016f5f33DbCb653e6393698Beba9DC19d828D75e \
  'escrow()(address)' \
  --rpc-url http://127.0.0.1:59619


gaiad q txs \
  --query "message.action='/ibc.core.channel.v2.MsgAcknowledgement'" \
  --node tcp://127.0.0.1:26657 \
  -o json


cast receipt 0x4d611d65a802bea81865e7f0e1f0413064518a79883b29481968804d0efa1692--rpc-url http://127.0.0.1:56310 | grep -A3 "IBCAppRecvPacket\|topics\|data"


```

---

## 1. Packet Timeout

### 1a. Timeout Before Relay

Send a packet with an absolute timeout in the past (or a 1-second window), so the relayer cannot relay it in time.

**Steps:**
```bash
ABS_TIMEOUT=$(($(date +%s) + 1))   # chỉ còn 1 giây

gaiad tx ibc-transfer transfer transfer 08-wasm-0 \
  0x8943545177806ed17b9f23f0a21ee5948ecaa776 100stake \
  --absolute-timeouts \
  --packet-timeout-timestamp "$ABS_TIMEOUT" \
  --from test1 --home ~/.gaia --chain-id test-ibc-eth \
  --node tcp://127.0.0.1:26657 --keyring-backend test \
  --gas-prices 1stake -y

```

**Expected outcome:**
- The relayer refuses to relay (`recvPacket` reverts with timeout error).
- After the relayer submits `timeoutPacket` back to Cosmos, the sender's balance is restored (escrow releases).
- `MsgTimeout` appears in Cosmos tx history:
  ```bash
  gaiad q txs --query "message.action='/ibc.core.channel.v2.MsgTimeout'" \
    --node tcp://127.0.0.1:26657 -o json
  ```

  ```bash
  gaiad q bank balances $(gaiad keys show test1 -a --home ~/.gaia --keyring-backend test) \
  --node tcp://127.0.0.1:26657

  ```

### 1b. Timeout Race — Relayer Wins

Set a 10-second window but delay the relay deliberately (e.g., pause the relayer for 8 s then resume).

```

# Terminal 1: gửi packet với 10-giây window
ABS_TIMEOUT=$(($(date +%s) + 60))
gaiad tx ibc-transfer transfer transfer 08-wasm-0 \
  0x8943545177806ed17b9f23f0a21ee5948ecaa776 100stake \
  --absolute-timeouts --packet-timeout-timestamp "$ABS_TIMEOUT" \
  --from test1 --home ~/.gaia --chain-id test-ibc-eth \
  --node tcp://127.0.0.1:26657 --keyring-backend test --gas-prices 1stake -y

```

**Expected outcome:** Packet is relayed successfully; no timeout.

---

## 2. Relayer Crash and Recovery

### 2a. Crash Between updateClient and recvPacket

1. Start relay loop.
2. Send a transfer (step 7 of the happy path).
3. Immediately kill the relayer process (`Ctrl-C` or `kill`) after `[UpdateCosmosClient]` log appears but before `[RecvPacket]`.
4. Restart the relayer.

**Expected outcome:**
- Relayer detects the pending packet on restart (queries `ICS26Router` for unrelayed packets).
- `recvPacket` is submitted successfully without re-running `updateClient` (client is already up-to-date).
- Cosmos ACK arrives; final balance matches expected.


```bash
cast call 0x016f5f33DbCb653e6393698Beba9DC19d828D75e \
  'balanceOf(address)(uint256)' \
  0x8943545177806ed17b9f23f0a21ee5948ecaa776 \
  --rpc-url http://127.0.0.1:62880


gaiad q txs \
  --query "message.action='/ibc.core.channel.v2.MsgAcknowledgement'" \
  --node tcp://127.0.0.1:26657 -o json | jq '.total_count'

```

### 2b. Network Partition Mid-Relay

Simulate an Ethereum RPC outage (e.g., `iptables` block or stop Kurtosis Ethereum service) while the relay loop is running.

**Expected outcome:**
- Relayer logs an error and retries with backoff.
- After connectivity is restored, the relay resumes and delivers all pending packets.


```
ABS_TIMEOUT=$(($(date +%s) + 10000))

gaiad tx ibc-transfer transfer transfer 08-wasm-0 0x8943545177806ed17b9f23f0a21ee5948ecaa776 1000stake \
  --from test1 \
  --home /Users/donglieu/.gaia \
  --chain-id test-ibc-eth \
  --node tcp://127.0.0.1:26657 \
  --keyring-backend test \
  --gas-prices 1stake \
  --absolute-timeouts \
  --packet-timeout-timestamp "$ABS_TIMEOUT" \
  --generate-only \
| jq '.body.messages[0].encoding = "application/x-solidity-abi"' \
| gaiad tx sign /dev/stdin \
    --from test1 \
    --home /Users/donglieu/.gaia \
    --chain-id test-ibc-eth \
    --keyring-backend test \
| gaiad tx broadcast /dev/stdin \
    --node tcp://127.0.0.1:26657 \
    -y


kurtosis service stop my-testnet el-1-geth-lighthouse


kurtosis service start my-testnet el-1-geth-lighthouse



kurtosis service stop my-testnet el-1-geth-lighthouse
kurtosis service stop my-testnet cl-1-lighthouse-geth

kurtosis service start my-testnet el-1-geth-lighthouse
kurtosis service start my-testnet cl-1-lighthouse-geth


cast call 0x016f5f33DbCb653e6393698Beba9DC19d828D75e \
  'balanceOf(address)(uint256)' \
  0x8943545177806ed17b9f23f0a21ee5948ecaa776 \
  --rpc-url http://127.0.0.1:59619


gaiad q txs \
  --query "message.action='/ibc.core.channel.v2.MsgAcknowledgement'" \
  --node tcp://127.0.0.1:26657 -o json | jq '.total_count'


gaiad q txs --query "message.action='/ibc.core.channel.v2.MsgTimeout'"  --node tcp://127.0.0.1:26657 -o json | jq '.total_count'

```
---

## 3. Replay / Double-Relay Attack

Attempt to submit the same `recvPacket` call-data twice.

**Steps:**
1. Capture the raw calldata of a successful `recvPacket` transaction.
2. Re-broadcast it (or call `ICS26Router.recvPacket` again with identical arguments).

**Expected outcome:**
- Second call reverts: `ICS24Host` already stores the packet commitment/receipt.
- No duplicate token minting on Ethereum.

---

## 4. Invalid / Forged Proof

### 4a. Tampered Groth16 Proof

Modify one byte of the `proof` field in the `updateClient` calldata before broadcasting.

**Recommended command (safe / `eth_call`):**
```bash
./scripts/test_invalid_updateclient.sh --case 4a \
  --config relayer/config.example.json \
  --env-file relayer/.env
```

**Expected outcome:**
- Script prints `PASS: tampered eth_call reverted as expected`.
- `eth_call` reverts and client state is unchanged.

### 4b. Wrong Public Input (Witness Hash Mismatch)

Submit a valid proof from a previous `updateClient` call paired with a new header (different `appHash`).

**Recommended command (safe / `eth_call`):**
```bash
./scripts/test_invalid_updateclient.sh --case 4b \
  --config relayer/config.example.json \
  --env-file relayer/.env
```

**Expected outcome:**
- Script prints `PASS: tampered eth_call reverted as expected`.
- On-chain verification rejects the forged payload (revert reason may vary by which validation step fails first).

### 4c. Pubkey Swap in Calldata

Supply a real validator pubkey in calldata but submit a proof computed for a different pubkey.

**Recommended command (safe / `eth_call`):**
```bash
./scripts/test_invalid_updateclient.sh --case 4c \
  --config relayer/config.example.json \
  --env-file relayer/.env
```

**Expected outcome:**
- Script prints `PASS: tampered eth_call reverted as expected`.
- Verification fails because pubkeys no longer match the proven witness.

**Optional (broadcast failing tx on-chain):**
```bash
./scripts/test_invalid_updateclient.sh --case 4b --send
```
or
```bash
./scripts/test_invalid_updateclient.sh --case 4c --send
```

---

## 5. Voting-Power / Quorum Edge Cases

### 5a. Below 2/3 Threshold

Construct a proof using validators whose total voting power sums to exactly `2/3 * totalVotingPower` (not strictly greater).

**Expected outcome:**
- `Groth16ICS07Tendermint` checks `power * 3 > totalVotingPower * 2` (strict). Exactly 2/3 fails.

### 5b. Duplicate Signer Index

Include the same validator index twice in the calldata (active flag set for both occurrences).

**Expected outcome:**
- `seen[idx]` bitmap rejects the duplicate; `updateClient` reverts.

### 5c. Bucket Boundary — Exactly N Active Signers

Force a quorum with exactly 4, 8, or 16 unique signers (matching a bucket boundary) — no padding needed.

**Expected outcome:**
- Relayer picks the exact-fit bucket (no wasted padding slots).
- `[UpdateCosmosClient]` log line reports the correct bucket size.
- Proof verifies; relay succeeds.

### 5d. Quorum Requires More Than Largest Bucket

Configure a local Cosmos chain with 200+ validators all with equal voting power (so ≥134 signers are needed for 2/3 quorum, exceeding the 128-bucket cap).

**Expected outcome:**
- Relayer logs an error: no bucket ≥ required signers.
- `updateClient` is never attempted.
- Operator must compile a new 256-bucket and redeploy before relay is possible.

---

## 6. Light Client Staleness and Misbehaviour
### 6a. Trusting Period Expiry

This repo now supports a custom `trusting_period` via `./relayer create-clients --trusting-period`.

**Setup:**
1. Create a dedicated test client with a short trusting period:
   ```bash
   cd relayer
   cp config.example.json config.expiry-test.json
   ./relayer create-clients \
     --config config.expiry-test.json \
     --wasm-checksum 0xd24688886ed8cec00c667fa69c173fbab9a08c75900ce18afe10517c82e55592 \
     --trust-level 2/3 \
     --trusting-period 60
   ```
2. Start the relayer with that config:
   ```bash
   ./relayer start --config config.expiry-test.json
   ```

**Steps:**
1. Stop the relay loop using `config.expiry-test.json`.
2. Wait longer than the configured trusting period, for example `75-90s` when `--trusting-period 60` was used.
3. Send a new Cosmos -> ETH packet.
4. Start the relayer again with `config.expiry-test.json`.

**Expected outcome:**
- `Groth16ICS07Tendermint` rejects the header: client is expired.
- Relayer logs an `updateClient` failure and does not submit the packet.
- A representative revert reason is:
  ```text
  invalid block: untrusted state is outside of trusting period
  ```

**Observed result on local test:**
- `PASS`
- Final relayer logs included:
  ```text
  [UpdateCosmosClient] SendEthTx failed: tx ... reverted
  [StartLoop] Failed to update cosmos light client: tx ... reverted
  ```

### 6b. Misbehaviour — Equivocation

The repo now includes a helper at `relayer/cmd/misbehaviour_attack` to build and submit `submitMisbehaviour`.

Important note: on the current honest Gaia setup, we do not have a source of two valid conflicting headers at the same height. The exercised path below uses two valid headers from different heights with different `appHash`, which is enough to trigger the current on-chain misbehaviour implementation and freeze the client.

```bash
cd relayer

# 1. Create a dedicated client for the misbehaviour test.
cp config.example.json config.misbehaviour-test.json
./relayer create-clients \
  --config config.misbehaviour-test.json \
  --wasm-checksum 0xd24688886ed8cec00c667fa69c173fbab9a08c75900ce18afe10517c82e55592 \
  --trust-level 2/3 \
  --trusting-period 600

# 2. Start the relayer and send at least 2 Cosmos -> ETH packets
# so the client accumulates 2 updateClient txs on Ethereum.
./relayer start --config config.misbehaviour-test.json

# 3. Run the helper in dry-run mode first.
GOCACHE=/tmp/go-build go run ./cmd/misbehaviour_attack \
  --config config.misbehaviour-test.json \
  --tx1 <newer_updateClient_tx_hash> \
  --tx2 <older_updateClient_tx_hash>

# 4. If simulation succeeds, submit the misbehaviour tx.
GOCACHE=/tmp/go-build go run ./cmd/misbehaviour_attack \
  --config config.misbehaviour-test.json \
  --tx1 <newer_updateClient_tx_hash> \
  --tx2 <older_updateClient_tx_hash> \
  --submit
```

**Expected outcome:**
- Client transitions to a frozen state.
- Subsequent `updateClient` and `verifyMembership` calls revert with `FrozenClientState`.
- Governance / admin must unfreeze or create a new client to resume relaying.

**Observed result on local test:**
- `PASS`
- Pre-submit simulation returned `simulate_ok=true`.
- After submit, rerunning the helper returned:
  ```text
  simulate_error=execution reverted
  simulate_decoded_error=FrozenClientState
  client_frozen_before=true
  ```

---

## 7. Reverse Direction — ETH → Cosmos

Send a token transfer from Ethereum to Cosmos (the reverse of the happy path).

**Steps:**
1. Complete the happy path first so `IBCERC20` tokens exist on Ethereum.
2. Resolve the wrapped token address on Ethereum:
   ```bash
   cast call <ICS20_ADDRESS> 'ibcERC20Contract(string)(address)' \
     'transfer/cosmoshub-1/stake' \
     --rpc-url http://127.0.0.1:<eth_rpc_port>
   ```
3. Approve the `ICS20Transfer` contract to spend the `IBCERC20` tokens:
   ```bash
   cast send <IBCERC20_ADDRESS> 'approve(address,uint256)' <ICS20_ADDRESS> 500 \
     --rpc-url http://127.0.0.1:<eth_rpc_port> --private-key $ETH_PRIVATE_KEY
   ```
4. Call `ICS20Transfer.sendTransfer` on Ethereum:
   ```bash
   ABS_TIMEOUT=$(($(date +%s) + 600))

   cast send <ICS20_ADDRESS> \
     'sendTransfer((address,uint256,string,string,string,uint64,string))' \
     "(<IBCERC20_ADDRESS>,500,'<cosmos_receiver>','cosmoshub-1','transfer',$ABS_TIMEOUT,'')" \
     --rpc-url http://127.0.0.1:<eth_rpc_port> \
     --private-key $ETH_PRIVATE_KEY
   ```
5. Let the relay loop pick it up.
6. Verify that the original Cosmos receiver gets the tokens back.

**Expected outcome:**
- Tokens are burned/escrowed on Ethereum.
- Cosmos chain credits the recipient.
- Ethereum light client on Cosmos is updated via the Go relayer.

**Observed result on local test:**
- `PASS`
- Happy path first minted `1000` wrapped `stake` on ETH token `0x016f5f33DbCb653e6393698Beba9DC19d828D75e`.
- Reverse transfer then sent `500` back to Cosmos from Ethereum tx `0x071aa593c808f88c413e9ce85284689599c72c24a82d64fde0e655942a385cc9`.
- After relay completion, `balanceOf(faucet)` on that wrapped token was `500`.
- Cosmos receive tx was `66DF08567603D5D9E841937E78F97B050655532E4D9FA9B8EA1109B2EF0C31B5`.
- That tx included `coin_received=500stake` and `fungible_token_packet success=true`.

**Notes:**
- In this setup, `timeoutTimestamp` is interpreted as unix seconds.
- If the relayer signer reuses the same Cosmos account as the recipient, wallet balance alone is noisy because relayer fees are paid from that account. For verification, prefer the `coin_received` and `fungible_token_packet` events in the Cosmos recv tx.

---

## 8. Rate Limit Exceeded

Send many transfers in quick succession to hit the `RateLimitUpgradeable` cap.

**Steps:**
1. Complete the happy path first so the wrapped `stake` token already exists on Ethereum.
2. Resolve the wrapped token and escrow:
   ```bash
   cast call <ICS20_ADDRESS> 'ibcERC20Contract(string)(address)' \
     'transfer/cosmoshub-1/stake' \
     --rpc-url http://127.0.0.1:<eth_rpc_port>

   cast call <ICS20_ADDRESS> 'getEscrow(string)(address)' \
     'cosmoshub-1' \
     --rpc-url http://127.0.0.1:<eth_rpc_port>
   ```
3. Configure a rate limit on that escrowed wrapped token.
   In the current local E2E deployment, `RATE_LIMITER_ROLE` is not wired by default, so you must first grant the selector role on the escrow via the deployed `AccessManager`, then call:
   ```bash
   cast send <ESCROW_ADDRESS> 'setRateLimit(address,uint256)' <IBCERC20_ADDRESS> 1500 \
     --rpc-url http://127.0.0.1:<eth_rpc_port> --private-key $ETH_PRIVATE_KEY
   ```
4. Send transfers summing to just below the cap, for example one packet of `1000stake` (should succeed).
5. Send one more transfer that crosses the cap, for example `600stake`.

**Expected outcome:**
- The packet that crosses the cap is accepted on Cosmos but reverted on Ethereum by `ICS20Transfer` with a rate-limit error.
- The relayer submits a write-acknowledgement with an error back to Cosmos so the sender can reclaim funds.

**Observed result on local test:**
- `PASS`
- Rate limit was set to `1500`.
- First packet sent `1000stake` and succeeded.
- After the first packet, escrow `dailyUsage=1000` and wrapped balance of the faucet was `1500`.
- Second packet sent `600stake` and was not credited on Ethereum.
- After the second packet, escrow `dailyUsage` remained `1000` and wrapped balance remained `1500`.
- Cosmos acknowledgement tx for the failed packet was `F909DEDDBD75D9310A7E4EF23CFE4C5B356F0251E22960DE53684A1508583831`.
- That tx included `fungible_token_packet acknowledgement = error:"ABCI code: 16: error handling packet: see events for details"`.

**Notes:**
- The current relayer/CLI setup often reuses `test1` both for manual Cosmos sends and for relayer-submitted Cosmos txs. That can cause `account sequence mismatch` while relaying the final acknowledgement; restarting the relayer is enough to recover and clear the pending historical `WriteAcknowledgement`.

---

## 9. Concurrent Packets

Send N packets from Cosmos in rapid succession (in the same block or back-to-back blocks) without waiting for relay acknowledgements.

```bash
for i in $(seq 1 5); do
  gaiad tx ibc-transfer transfer transfer 08-wasm-0 0x8943545177806ed17b9f23f0a21ee5948ecaa776 100stake \
    --from test1 --home ~/.gaia --chain-id test-ibc-eth \
    --node tcp://127.0.0.1:26657 --keyring-backend test --gas-prices 1stake -y
  sleep 0.5
done
```

**Expected outcome:**
- All 5 packets are relayed (possibly batched into fewer `updateClient` calls if they land in the same Cosmos block).
- Final Ethereum balance equals 500 wrapped-stake tokens.
- All 5 ACKs arrive back on Cosmos.

**Observed result on local test:**
- `PASS`
- Local run used an isolated ETH receiver `0x1111111111111111111111111111111111111111`, starting from balance `0`.
- Five packets were created with `seq=4..8`.
- ETH receiver final wrapped balance was `500`.
- Cosmos acknowledgement query ended with `total_count=8`, and the latest five `MsgAcknowledgement` txs corresponded to `seq=4..8`.
- Relayer logs showed mixed batching due timing:
  `seq=4` was relayed first,
  `seq=5` followed,
  and `seq=6..8` were queued together and later flushed as one Cosmos batch of `3` packets.

**Notes:**
- If the relayer signer already uses `test1` on Cosmos, avoid using the same account for the rapid manual sends in this test. Using other funded local accounts such as `test2` and `test3` avoids `account sequence mismatch`.
- In this setup, the `application/x-solidity-abi` encoding patch is still required for the generated Cosmos transfer messages, same as the other Cosmos -> ETH scenarios above.

---

## 10. Wrong WASM Checksum During Client Creation

Pass a checksum that does not match the uploaded WASM binary.

```bash
./relayer create-clients --config config.example.json \
  --wasm-checksum 0x0000000000000000000000000000000000000000000000000000000000000000
```

**Expected outcome:**
- Cosmos rejects the `MsgCreateClient` with a checksum-mismatch error.
- Relayer exits with a descriptive error; no partial state is left on-chain.

**Observed result on local test:**
- `PASS`
- `create-clients` failed early with:
  `wasm checksum 0x0000000000000000000000000000000000000000000000000000000000000000 has not been previously stored on Cosmos`
- The command exited before any new ICS07 deploy or router migration on Ethereum.
- The test config file was not rewritten.
- Cosmos did not create a new wasm client; client-state count remained `1`.

**Notes:**
- The flow now validates the wasm checksum on Cosmos before mutating ETH state or rewriting config, which aligns the implementation with the intended semantics of this test.

---

## 11. Missing or Invalid Environment Variables

Start the relayer with a missing or malformed `ETH_PRIVATE_KEY` or `COSMOS_PRIVATE_KEY`.

**Expected outcome:**
- Relayer fails at startup with a clear key-load error before attempting any RPC calls.

**Observed result on local test:**
- `PASS`
- `env ETH_PRIVATE_KEY= ./relayer start --config config.example.json`
  failed immediately with:
  `ETH_PRIVATE_KEY environment variable is required in .env file`
- `env COSMOS_PRIVATE_KEY=notvalidhex ./relayer start --config config.example.json`
  failed immediately with:
  `failed to decode COSMOS_PRIVATE_KEY`

**Notes:**
- The startup flow now validates both private keys immediately after `.env` is loaded.
- In both local runs, the command exited before attempting Cosmos or Ethereum RPC connections.

---

## 12. Circuit / Verifier Mismatch After Re-Setup

Re-run `go run ./prover/cmd ./bin ../contracts/verifiers` (which generates fresh `(pk, vk)` with new randomness) **without** redeploying `Groth16Verifier_N{N}.sol`.

**Expected outcome:**
- `WrapperVerifier` dispatches to the old verifier contract whose VK no longer matches the new proving key.
- `updateClient` reverts on every call (pairing check fails).
- Fix: redeploy all `Groth16Verifier_N{N}.sol` contracts atomically and re-register them via `WrapperVerifier.setBucket(...)`.

**Observed result on local test:**
- `PASS`
- I generated fresh prover artifacts into `/private/tmp/test12-bin` and started relayer with
  `PROVER_BIN_DIR=/private/tmp/test12-bin`, while keeping the existing on-chain verifier set unchanged.
- Sending a new Cosmos -> ETH packet produced Cosmos tx
  `2BA3AFC0FF4C9FFDF37B5345360201F55C043B0708D4D3F0943F2B38946E9C74`
  with `seq=9`.
- Relayer built the proof locally, but the on-chain `updateClient` tx reverted:
  `0x9151f794bc6bfaa5171b86b25ff1df0d5976fc7240b7bf402c495d4676cc1314`
  and returned revert data `0xd611c318`.
- The ETH receiver `0x2222222222222222222222222222222222222222` stayed at balance `0`,
  confirming the packet was not credited after the verifier mismatch.

---

## Summary Table

| # | Scenario | Layer | Expected Result |
|---|----------|-------|-----------------|
| 1a | Packet timeout — expired | Protocol | Timeout ACK, escrow released |
| 1b | Packet timeout — relayer wins race | Protocol | Normal delivery |
| 2a | Relayer crash post-updateClient | Relayer | Resumes on restart, delivers packet |
| 2b | Ethereum RPC outage | Relayer | Retries, recovers automatically |
| 3 | Replay / double recvPacket | Security | Second call reverts |
| 4a | Tampered Groth16 proof | Security | `updateClient` revert |
| 4b | Wrong public input | Security | `updateClient` revert |
| 4c | Pubkey swap | Security | `updateClient` revert |
| 5a | Exactly 2/3 voting power | Quorum | Strict-greater check fails |
| 5b | Duplicate signer index | Quorum | `seen[idx]` revert |
| 5c | Exactly N active signers at bucket edge | Quorum | Correct bucket selected, success |
| 5d | Quorum > 128 signers | Quorum | Relayer error, no submission |
| 6a | Trusting period expiry | Light client | Client expired revert |
| 6b | Misbehaviour / equivocation | Light client | Client frozen |
| 7 | ETH → Cosmos reverse transfer | Protocol | Tokens returned to Cosmos sender |
| 8 | Rate limit exceeded | Protocol | Error ACK, funds recoverable |
| 9 | 5 concurrent packets | Protocol | All delivered, correct final balance |
| 10 | Wrong WASM checksum | Setup | Cosmos rejects MsgCreateClient |
| 11 | Missing private key env var | Setup | Relayer exits at startup |
| 12 | Circuit/verifier mismatch after re-setup | Ops | All updateClient calls revert |
