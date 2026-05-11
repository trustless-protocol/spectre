# Test Scenarios

The 8-step flow in the README covers the happy path: compile circuits → start nodes → create clients → relay → send transfer → verify balance. This document describes additional scenarios worth testing.

---

## 1. Packet Timeout

### 1a. Timeout Before Relay

Send a packet with an absolute timeout in the past (or a 1-second window), so the relayer cannot relay it in time.

**Steps:**
```bash
ABS_TIMEOUT=$(($(date +%s) + 1))   # 1-second window — expires immediately

gaiad tx ibc-transfer transfer transfer 08-wasm-0 0x8943545177806ed17b9f23f0a21ee5948ecaa776 100stake \
  --absolute-timeouts --packet-timeout-timestamp "$ABS_TIMEOUT" \
  --from test1 --home ~/.gaia --chain-id test-ibc-eth \
  --node tcp://127.0.0.1:26657 --keyring-backend test --gas-prices 1stake -y
```

**Expected outcome:**
- The relayer refuses to relay (`recvPacket` reverts with timeout error).
- After the relayer submits `timeoutPacket` back to Cosmos, the sender's balance is restored (escrow releases).
- `MsgTimeout` appears in Cosmos tx history:
  ```bash
  gaiad q txs --query "message.action='/ibc.core.channel.v2.MsgTimeout'" \
    --node tcp://127.0.0.1:26657 -o json
  ```

### 1b. Timeout Race — Relayer Wins

Set a 10-second window but delay the relay deliberately (e.g., pause the relayer for 8 s then resume).

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

### 2b. Network Partition Mid-Relay

Simulate an Ethereum RPC outage (e.g., `iptables` block or stop Kurtosis Ethereum service) while the relay loop is running.

**Expected outcome:**
- Relayer logs an error and retries with backoff.
- After connectivity is restored, the relay resumes and delivers all pending packets.

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

**Expected outcome:**
- `Groth16Verifier_N{N}.sol` reverts — pairing check fails.

### 4b. Wrong Public Input (Witness Hash Mismatch)

Submit a valid proof from a previous `updateClient` call paired with a new header (different `appHash`).

**Expected outcome:**
- `WrapperVerifier` recomputes the SHA-256 witness commit from the supplied header and validators; the recomputed hash does not match the proof's public input → revert.

### 4c. Pubkey Swap in Calldata

Supply a real validator pubkey in calldata but submit a proof computed for a different pubkey.

**Expected outcome:**
- Because pubkey `A` is hashed into the witness commit on-chain, the recomputed public input diverges from the proof's → revert.

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

1. Stop the relay loop.
2. Wait until `now > lastUpdateTime + trustingPeriod`.
3. Attempt to relay a new packet.

**Expected outcome:**
- `Groth16ICS07Tendermint` rejects the header: client is expired.
- Relayer logs a "client expired" error and does not submit the packet.

### 6b. Misbehaviour — Equivocation

Submit two conflicting headers at the same height (same `trustedHeight` but different `appHash`) as a misbehaviour report.

```bash
# Pseudo-call via cast
cast send <ICS07_ADDRESS> 'submitMisbehaviour(bytes)' <misbehaviour_abi_encoded> \
  --rpc-url http://127.0.0.1:<eth_rpc_port> --private-key $ETH_PRIVATE_KEY
```

**Expected outcome:**
- Client transitions to a frozen state.
- Subsequent `updateClient` and `verifyMembership` calls revert with "client frozen".
- Governance / admin must unfreeze or create a new client to resume relaying.

---

## 7. Reverse Direction — ETH → Cosmos

Send a token transfer from Ethereum to Cosmos (the reverse of the happy path).

**Steps:**
1. Complete the happy path first so `IBCERC20` tokens exist on Ethereum.
2. Approve the `ICS20Transfer` contract to spend the `IBCERC20` tokens:
   ```bash
   cast send <IBCERC20_ADDRESS> 'approve(address,uint256)' <ICS20_ADDRESS> 500 \
     --rpc-url http://127.0.0.1:<eth_rpc_port> --private-key $ETH_PRIVATE_KEY
   ```
3. Call `ICS20Transfer.sendTransfer` on Ethereum.
4. Let the relay loop pick it up.
5. Verify that the original Cosmos sender receives the tokens back.

**Expected outcome:**
- Tokens are burned/escrowed on Ethereum.
- Cosmos chain credits the recipient.
- Ethereum light client on Cosmos is updated via the Go relayer.

---

## 8. Rate Limit Exceeded

Send many transfers in quick succession to hit the `RateLimitUpgradeable` cap.

**Steps:**
1. Check the configured rate limit for the `stake` denom.
2. Send transfers summing to just below the cap (should succeed).
3. Send one more transfer that crosses the cap.

**Expected outcome:**
- The packet that crosses the cap is accepted on Cosmos but reverted on Ethereum by `ICS20Transfer` with a rate-limit error.
- The relayer submits a write-acknowledgement with an error back to Cosmos so the sender can reclaim funds.

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

---

## 11. Missing or Invalid Environment Variables

Start the relayer with a missing or malformed `ETH_PRIVATE_KEY` or `COSMOS_PRIVATE_KEY`.

**Expected outcome:**
- Relayer fails at startup with a clear key-load error before attempting any RPC calls.

---

## 12. Circuit / Verifier Mismatch After Re-Setup

Re-run `go run ./prover/cmd ./bin ../contracts/verifiers` (which generates fresh `(pk, vk)` with new randomness) **without** redeploying `Groth16Verifier_N{N}.sol`.

**Expected outcome:**
- `WrapperVerifier` dispatches to the old verifier contract whose VK no longer matches the new proving key.
- `updateClient` reverts on every call (pairing check fails).
- Fix: redeploy all `Groth16Verifier_N{N}.sol` contracts atomically and re-register them via `WrapperVerifier.setBucket(...)`.

---

## Summary Table

| # | Scenario | Layer | Expected Result |
|---|----------|-------|-----------------|
| 1a | Packet timeout — expired | Protocol | Timeout ACK, escrow released |
| 1b | Packet timeout — relayer wins race | Protocol | Normal delivery |
| 2a | Relayer crash post-updateClient | Relayer | Resumes on restart, delivers packet |
| 2b | Ethereum RPC outage | Relayer | Retries, recovers automatically |
| 3 | Replay / double recvPacket | Security | Second call reverts |
| 4a | Tampered Groth16 proof | Security | Pairing check revert |
| 4b | Wrong public input | Security | Witness hash mismatch revert |
| 4c | Pubkey swap | Security | Witness hash mismatch revert |
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
