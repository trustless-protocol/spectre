# Benchmark Report — Bucket 32 Happy Path (Cosmos → ETH transfer + ack back)

**Run**: 2026-05-18 14:23:41 → 14:27:16
**Branch**: `benchmark`
**Bucket**: **n=32** (21 active validator sigs, 2,948,514 constraints)
**Packets**: 1 (seq=1)
**Note**: First successful run after bumping ETH tx gas limit 3M → 16M.

## End-to-end wall-clock: **~3m 28s** (208 s)

```
prover start ───────┐
                    ├── Groth16 prove (sigs=21, bucket=32) ........  9.73 s   ◄ scales with bucket
                    ├── ETH updateClient tx .......................  1.00 s
                    ├── ETH recvPacket tx .........................  2.00 s
                    │   (WriteAcknowledgement event emitted)
                    ├── wait beacon finality (18 polls × 10s) .... 180 s     ◄ DOMINANT cost
                    ├── ETH→Cosmos updateClient prep + catch-up .  15 s
                    ├── Cosmos updateClient batch .................  2.72 s
                    └── Cosmos MsgAcknowledgement .................  5.48 s
```

## Gas usage

| Direction  | Tx                                  |     Gas used | Limit / wanted |
| ---------- | ----------------------------------- | -----------: | -------------: |
| Cosmos→ETH | `ICS26Router.updateClient`          |  **4,140,110** |     16,000,000 |
| Cosmos→ETH | `ICS26Router.recvPacket` (seq=1)    |  **1,033,531** |     16,000,000 |
| ETH→Cosmos | `MsgUpdateClient` (wasm LC verify)  |    **836,060** |      2,000,000 |
| ETH→Cosmos | `MsgAcknowledgement`                |    **173,198** |        200,000 |

## Latency breakdown

### Prover (Groth16, bucket n=32, 2.95M constraints)

| Stage             |     Duration |
| ----------------- | -----------: |
| Build witness     |      4.65 ms |
| `groth16.Prove`   |  **9.730 s** |
| Local verify      |      1.17 ms |
| **Total**         |  **9.735 s** |

gnark internal: hint gen 0.68 s, solver 3.45 s, prover 6.28 s.

### Ethereum tx submission

| Tx              | Submit (RPC send) | Wait (mining) |   Total |
| --------------- | ----------------: | ------------: | ------: |
| `updateClient`  |           2.31 ms |        1.00 s | 1.00 s |
| `recvPacket`    |           1.02 ms |        2.00 s | 2.00 s |

### Cosmos tx submission

| Tx                              | Broadcast (CheckTx + DeliverTx) |   Total |
| ------------------------------- | ------------------------------: | ------: |
| `MsgUpdateClient` (batch=1)     |                          2.71 s | 2.72 s |
| `MsgAcknowledgement`            |                          5.48 s | 5.48 s |

### Beacon finality wait (ETH side)

- Event block: 616 (emitted at 14:23:51)
- Finalized at: 14:26:53 (block 640 ≥ 616)
- **18 polls × 10 s = 180 s** (slightly higher than previous 160 s runs — finality epoch boundary variance).

## Comparison vs prior runs (same code, different bucket)

| Metric                              | n=4 (2026-05-16) | n=4 (2026-05-17) | **n=32 (2026-05-18)** |
| ----------------------------------- | ---------------: | ---------------: | -----------------: |
| Constraints                         |          734,503 |          734,503 |          **2,948,514** |
| Active sigs                         |                3 |                3 |                 **21** |
| `groth16.Prove`                     |          6.193 s |          2.418 s |          **9.730 s** |
| `ICS26Router.updateClient` gas      |        1,155,760 |        1,154,462 |        **4,140,110** |
| `ICS26Router.recvPacket` gas        |        1,033,543 |        1,033,531 |          1,033,531 |
| `MsgUpdateClient` gas (Cosmos)      |          822,282 |          836,061 |            836,060 |
| `MsgAcknowledgement` gas            |          156,303 |          156,303 |            173,198 |
| Beacon finality wait                |            160 s |            160 s |              180 s |
| E2E wall-clock                      |           ~198 s |           ~185 s |             ~208 s |

## Observations

1. **`updateClient` gas scales ~3.6× from n=4 → n=32** (1.15M → 4.14M).
   - Bucket loop in `_verifyBatchAndQuorum` grew (32 vs 4 iters): ~50k of the delta.
   - **`WrapperVerifier._hashWitness` O(N²) accumulator** is the main culprit: ~3M of the 4.1M total.
   - Projection holds: **bucket 64 ≈ 12–14M**, **bucket 128 ≈ 50M+** (won't fit Kurtosis 16M tx cap or mainnet block).
2. **`recvPacket` gas is flat** (1.03M) — independent of bucket, as expected (no Groth16 path).
3. **Prover time scales ~4× with constraints** (734k → 2.95M ≈ 4×; prove 2.4s → 9.7s ≈ 4×). Linear, no surprises.
4. **Beacon finality wait** still dominates wall-clock (~86%). Not affected by bucket size.
5. **Cosmos-side gas stable** between runs — wasm light client doesn't carry bucket logic.

## Implications

- **Bucket 32 is now functional** on the bumped 16M tx limit.
- **Bucket 64/128 unfeasible without contract fix.** The `_hashWitness` O(N²) memory accumulation must drop to O(N) before bucket 64+ will fit even the dev chain's 16M cap.
- **Mainnet target (~500k–1M gas)** still requires both:
  - `_hashWitness` rewrite (assembly mstore at known offsets).
  - Move witness hash assertion off-chain (commit via proof public input).

## Raw log artifacts captured

- `[bench][prover]`  → witness / prove / verify / total
- `[bench][eth]`     → label / gasUsed / submit / wait / total
- `[bench][cosmos]`  → msg type / gasWanted / gasUsed / broadcast / total
