# Benchmark Report — Happy Path (Cosmos → ETH transfer + ack back)

**Run**: 2026-05-16 09:55:14 → 09:58:32
**Branch**: `benchmark`
**Bucket**: n=4 (3 active validator sigs)
**Packets**: 1 (seq=1)

## End-to-end wall-clock: **~3m 18s** (198s)

```
Cosmos send_packet ─┐
                    ├── batch flush ............................... 1s
                    ├── wait 2 blocks for AppHash ................. 6s
                    ├── Groth16 prove (sigs=3, bucket=4) .......... 6.20s    ◄ on-chain blocker
                    ├── ETH updateClient tx ....................... 3.02s
                    ├── ETH recvPacket tx ......................... 2.01s
                    │   (WriteAck event emitted)
                    ├── wait beacon finality (16 polls × 10s) ... 160s      ◄ DOMINANT cost
                    ├── ETH→Cosmos updateClient prep + catch-up . 10s
                    ├── Cosmos updateClient batch ................. 3.90s
                    └── Cosmos MsgAcknowledgement ................. 5.28s
```

## Gas usage

| Direction  | Tx                                  | Gas used     | Limit / wanted |
| ---------- | ----------------------------------- | -----------: | -------------: |
| Cosmos→ETH | `ICS26Router.updateClient`          |  **1,155,760** |      3,000,000 |
| Cosmos→ETH | `ICS26Router.recvPacket` (seq=1)    |  **1,033,543** |      3,000,000 |
| ETH→Cosmos | `MsgUpdateClient` (wasm LC verify)  |    **822,282** |      2,000,000 |
| ETH→Cosmos | `MsgAcknowledgement`                |    **156,303** |        200,000 |

## Latency breakdown

### Prover (Groth16, bucket n=4)

| Stage             |     Duration |
| ----------------- | -----------: |
| Build witness     |      3.58 ms |
| `groth16.Prove`   |  **6.193 s** |
| Local verify      |      2.20 ms |
| **Total**         |  **6.199 s** |

Solver: 1.23 s (734,503 constraints) — included inside `Prove`.

### Ethereum tx submission

| Tx              | Submit (RPC send) | Wait (mining) |   Total |
| --------------- | ----------------: | ------------: | ------: |
| `updateClient`  |           10.6 ms |        3.00 s | 3.02 s |
| `recvPacket`    |           4.28 ms |        2.00 s | 2.01 s |

Wait time ≈ 1 block on the PoS dev network.

### Cosmos tx submission

| Tx                              | Broadcast (CheckTx + DeliverTx) |   Total |
| ------------------------------- | ------------------------------: | ------: |
| `MsgUpdateClient` (batch=1)     |                          3.89 s | 3.90 s |
| `MsgAcknowledgement`            |                          5.28 s | 5.28 s |

### Beacon finality wait (ETH side)

- Event block: 400 (emitted at 09:55:32)
- Finalized at: 09:58:13 (block 413 ≥ 400)
- **16 polls × 10s = 160 s** — dominant cost of the full round-trip

### One-time prover load (bucket n=4)

| Artifact   |     Size | Read time |
| ---------- | -------: | --------: |
| `r1cs.bin` |  72.7 MB |     ~1 s |
| `pk.bin`   |  200 MB  |  **~43 s** |
| `vk.bin`   |  1.5 KB  |    <1 ms |
| **Total**  | ~273 MB  | **~44 s** |

## Conclusions

1. **Wall-clock is dominated by beacon finality** (~81% of total). For dev/perf
   testing, consider a shorter epoch on Kurtosis or accepting `attestedSlot`
   instead of `finalized`.
2. **Groth16 proving (6.2 s)** is the critical on-chain prep cost — would scale
   superlinearly with bucket size (current 734k constraints at n=4). Worth
   re-running with n=64/128.
3. **Eth `updateClient` gas (1.16M)** is the most expensive Solidity call and
   matches expectations for Groth16 verify + quorum check at bucket=4. Will
   increase with bigger buckets due to more on-chain pubkey lookups.
4. **Cosmos `MsgUpdateClient` (822k gas)** runs the wasm light client verify —
   well under the 2M limit headroom.
5. **Prover startup (44 s)** is a one-time cost dominated by `pk.bin` read;
   could be amortized by a long-lived process or memory-map.

## Caveats

- Only 1 packet measured; tail / contention behavior not captured.
- Only bucket n=4 exercised (3 sigs). Need runs at n=32/64/128 for the
  "180-validator gaia-shape" benchmark.
- No `BenchGas` Solidity events yet (contracts not redeployed with the bench
  logging). Re-deploy + parse events to break gas down across
  `verifyMembership` / quorum / `_verifyBatchAndQuorum`.
- Cosmos broadcast time (3.9 s + 5.3 s) suggests ~1 block per tx — verify the
  dev chain's block time setting.

## Raw log artifacts captured

- `[bench][prover]`  → per-proof witness / prove / verify / total
- `[bench][eth]`     → per-tx label / gasUsed / submit / wait / total
- `[bench][cosmos]`  → per-tx msg type / gasWanted / gasUsed / broadcast / total

Source instrumentation:
- `relayer/prover/prover.go` (`GenerateProof`)
- `relayer/transaction/handler.go` (`SendEthTx`, `SendCosmosTx`, `SendCosmosTxBatch`)
- `contracts/ICS26Router.sol`, `contracts/light-clients/Groth16ICS07Tendermint.sol` (event `BenchGas`)
