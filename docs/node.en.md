## 1 Cosmos escrow stuck forever: does not submit MsgTimeout back to Cosmos

When sending a packet with a short timeout, the relayer may not relay it to Ethereum in time. When the timeout expires, the funds are not refunded automatically.

### solution

Add timeout fallback logic.

New flow for Cosmos → ETH packet timeout:

    - UpdateEthClient
    Update the ETH light client on Cosmos to the latest finalized slot.

    - GetEthereumClientState
    Fetch proofBlockNumber and proofSlot from the updated client state.

    - Build receipt path
    destClientID + [0x02] + seqBytes
    (0x02 is the receipt type byte, symmetric to ack type 0x03 used in WriteAck).

    - GetEthMembershipProof
    Call eth_getProof on Ethereum. Since the receipt was never written, the storage value is 0x0, which serves as the non-membership proof.

    - SendCosmosTx(MsgTimeout)
    Submit the timeout packet to Cosmos to release the escrow.

## 2a The relayer only subscribed from the current block and did not query historical events after restart
test 2a

Summary of what happened:

    -The old relayer relayed seq 6 successfully → ETH received it (balance = 4000)
    -The old relayer was killed (Ctrl+C)
    -A new relayer started 1 minute later → it missed the WriteAcknowledgement event at block 1425
    -MsgAcknowledgement for seq 6 was never submitted → Cosmos escrow still locked the funds
    -The token was minted on ETH but Cosmos did not know about it → inconsistent state

Why it failed: the relayer only subscribed from the current block and did not query historical events after restart.

### solution

On startup, the relayer now:

gets the latest ETH block
rescans WriteAcknowledgement logs in a lookback window
only enqueues events that are still pending on Cosmos by checking that the packet commitment still exists
then opens the live WS watch from latestBlock + 1 so there is no gap between recovery and live subscription

The key point is that the relayer no longer only “listens from now on”. If a WriteAcknowledgement happened while the relayer was down, it can still be recovered and submitted back to Cosmos as MsgAcknowledgement. The new logic also avoids replaying old acknowledgements that were already completed by skipping historical events once the escrow/commitment on Cosmos has already been cleared.

By default, the relayer rescans the last 256 ETH blocks on startup. This can be adjusted with `ETH_STARTUP_LOOKBACK_BLOCKS` if a larger window is needed for longer downtimes.

## 2b

After stopping the ETH node and starting it again, the relayer did not reconnect immediately.

The logs showed that the relayer had already consumed the batch and then lost it while ETH RPC was down.

Specifically in `relayer/services/main.go`:

receive batch

call `UpdateCosmosClient`

if error then `continue`

the old batch disappears, with no retry

### solution

This was fixed so that:

if `UpdateCosmosClient` fails or returns nil, the relayer requeues the batch instead of losing it

The fix is in `relayer/services/main.go` and its helper.

After ETH comes back up, the batch is flushed again a few seconds later and retried.

One remaining note: if both EL and CL are stopped, the ETH WS subscription can still die and fail to resubscribe automatically. This fix handles the “lost batch during ETH RPC outage” part, but if WriteAcknowledgement is still missed afterward, the next step is to add auto-resubscribe logic in `relayer/subscriber/event.go`.

For a more stable test, only stopping/starting `el-1-geth-lighthouse` is still recommended for 2b.

## 2b
failed because:

there was no acknowledgement (`MsgAcknowledgement`) and also no `MsgTimeout` on Cosmos
which means the state was: ETH had already minted, but Cosmos had not received the acknowledgement

## 3
ok

## 4a
ok

## 4b
ok

## 4c

ok

## 5a
pass

## 5b

pass

## 5c
skip

On the current setup, the latest on-chain payload had `bucket=4`, `active_count=3`, `unique_active=3`, so it was not the exact boundary case.

## 5d
skip

On the current setup, `MAX_BUCKET=4`, `min_signers_for_2of3=3`, so it did not hit the “requires more than largest bucket” case.

## 6a

Without a change, this case could not be tested in a practical way.

Reason:

- `./relayer create-clients` originally did not support `--trusting-period`
- the relayer runtime also did not read `trusting_period` from the test config
- because of that, the client always used the long default trusting period from the real chain, which made expiry testing take too long

### solution

Add `trusting_period` support in both client creation flow and relayer runtime config.

Then a dedicated test client can be created with a short trusting period, for example `60s`, and the test can run with this flow:

- start relayer with the test config
- stop relayer
- wait until trusting period has expired
- send a new Cosmos -> ETH packet
- start relayer again

Expected result, which was also observed locally:

- the relayer receives the packet and tries `updateClient`
- ETH reverts with:
  `invalid block: untrusted state is outside of trusting period`
- final logs:
  `[UpdateCosmosClient] SendEthTx failed: tx ... reverted`
  `[StartLoop] Failed to update cosmos light client: tx ... reverted`

pass

## 6b

Without a change, the current misbehaviour path could not freeze the client.

Reason:

- `Misbehaviour.sol` returned `trustedHeight1/2` using `header.signedHeader.header.height`
- but `Groth16ICS07Tendermint` verified `trustedConsensusState1/2` using `header.trustedHeight`
- so the input/output semantics were inconsistent, which made the misbehaviour path almost impossible to satisfy correctly

There was also another issue:

- the trusting-period check in `Misbehaviour.sol` compared nanosecond timestamps directly against `trustingPeriod`, which is measured in seconds
- because of that, even valid payloads could revert with a false expiry condition

### solution

Fix the misbehaviour logic as follows:

- `trustedHeight1/2` must come from `header.trustedHeight`
- the timestamp must be converted from nanoseconds to seconds before comparing with `trustingPeriod`

After that, the misbehaviour path can freeze the client successfully.

Observed locally:

- misbehaviour submission succeeded
- rechecking afterward showed the client was frozen
- subsequent calls reverted with `FrozenClientState`

Note:

- this case currently passes according to the on-chain implementation
- the local test uses two valid headers at different heights with different `appHash` values, not true same-height equivocation on an honest Gaia chain

pass

## 7

pass

Observed locally:

- the happy path `Cosmos -> ETH` minted `1000` wrapped `stake` on ETH
- then `ETH -> Cosmos` sent back `500`
- the wrapped token balance on ETH became `500`
- Cosmos successfully received `500stake` back through `MsgRecvPacket`

Notes for local runs:

- `timeoutTimestamp` in this setup must use unix seconds
- if account `test1` is also used as relayer signer, Cosmos wallet balances will be noisy because the relayer also spends gas from the same account
- for the `ETH -> Cosmos` case, the cleanest proof is the `coin_received 500stake` event and `fungible_token_packet success=true` in the recv tx on Cosmos

## 8

pass

Observed locally:

- set a rate limit of `1500` for wrapped `stake` on escrow
- the first packet sending `1000stake` succeeded
- after the first packet, `dailyUsage=1000` and the faucet wrapped balance increased to `1500`
- the next packet sending `600stake` was not credited on ETH
- after the second packet, `dailyUsage` was still `1000` and wrapped balance was still `1500`
- Cosmos received an error `MsgAcknowledgement` for `seq=3`

Notes for local runs:

- the E2E deployment does not yet wire `RATE_LIMITER_ROLE` for escrow by default, so that role must be granted in `AccessManager` before calling `setRateLimit(...)`
- if both the relayer and manual txs use account `test1` on Cosmos, `account sequence mismatch` can happen; restarting the relayer is enough for recovery and resubmitting the pending acknowledgement

## 9

pass

Observed locally:

- sent 5 packets in a row with `seq=4..8`
- ETH received all 5 packets
- final balance of the new ETH receiver was `500`
- Cosmos had all 5 matching `MsgAcknowledgement` messages for `seq=4..8`

Notes for local runs:

- if the relayer signer uses `test1` on Cosmos, `test1` should not also be used to send many packets quickly; use another account such as `test2` or `test3` to avoid `account sequence mismatch`
- when reusing an old environment, it is better to use a new ETH receiver with balance `0` so the final `500` result is easier to verify

## 10

pass

Observed locally:

- `create-clients` failed early at the checksum preflight step
- returned error:
  `wasm checksum 0x0000000000000000000000000000000000000000000000000000000000000000 has not been previously stored on Cosmos`
- no new ICS07 contract was deployed on ETH
- config was not rewritten with a temporary `ics07_client`
- Cosmos did not create a new wasm client, and the total client-state count remained `1`

Notes for local runs:

- the flow was fixed to validate the checksum on Cosmos before mutating ETH state or config

## 11

pass

Observed locally:

- starting the relayer with empty `ETH_PRIVATE_KEY` failed immediately with:
  `ETH_PRIVATE_KEY environment variable is required in .env file`
- starting the relayer with `COSMOS_PRIVATE_KEY=notvalidhex` failed immediately with:
  `failed to decode COSMOS_PRIVATE_KEY`

Notes for local runs:

- startup flow was fixed to validate keys immediately after `godotenv.Load()`
- the relayer fails before dialing either Cosmos or Ethereum RPC

## 12

pass

Observed locally:

- generated a fresh prover artifact set in a temporary directory and started the relayer with
  `PROVER_BIN_DIR=/private/tmp/test12-bin`
- did not redeploy `Groth16Verifier_N{N}.sol` and did not update `WrapperVerifier`
- sent one new Cosmos -> ETH packet, with Cosmos tx hash
  `2BA3AFC0FF4C9FFDF37B5345360201F55C043B0708D4D3F0943F2B38946E9C74`,
  `seq=9`
- the relayer built the proof successfully locally, but `updateClient` on ETH reverted with tx hash
  `0x9151f794bc6bfaa5171b86b25ff1df0d5976fc7240b7bf402c495d4676cc1314`
- returned revert data:
  `0xd611c318`
- ETH receiver `0x2222222222222222222222222222222222222222` remained at balance `0`

Notes for local runs:

- this case confirms that a mismatch between the new proving key and the old verifying key on-chain causes every `updateClient` call to fail
- the operational fix is to redeploy all corresponding verifier contracts and re-register the bucket in `WrapperVerifier` before using the new prover artifact set
