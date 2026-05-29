# Benchmark Report — Bucket 16 Happy Path (Cosmos → ETH transfer + ack back, ICICLE/GPU)

**Run**: 2026-05-25 09:56:18 → 10:00:00
**Branch**: `benchmark`
**Bucket**: **n=16** (14 active validator sigs, 1,685,171 constraints)
**Packets**: 1 (seq=1)
**Proof backend**: `icicle` (CUDA)
**Note**: Relayer preloaded prover artifacts for buckets `n=4`, `n=8`, and `n=16` before the measured packet flow.

## End-to-end wall-clock: **~3m 42s** (222 s)

```text
Cosmos send_packet ─┐
                    ├── batch flush ...............................  1 s
                    ├── wait 2 blocks for AppHash .................  6 s
                    ├── Groth16 prove (sigs=14, bucket=16) ........  3.52 s
                    ├── ETH updateClient tx .......................  2.00 s
                    ├── ETH recvPacket tx .........................  2.00 s
                    │   (WriteAcknowledgement event emitted)
                    ├── wait beacon finality (18 polls × 10s) .... 180 s     ◄ DOMINANT cost
                    ├── ETH→Cosmos updateClient prep + catch-up .. 15 s
                    ├── Cosmos updateClient batch .................  4.97 s
                    └── Cosmos MsgAcknowledgement .................  5.38 s
```

## Gas usage

| Direction  | Tx                                  |     Gas used | Limit / wanted |
| ---------- | ----------------------------------- | -----------: | -------------: |
| Cosmos→ETH | `ICS26Router.updateClient`          |  **1,709,704** |     16,000,000 |
| Cosmos→ETH | `ICS26Router.recvPacket` (seq=1)    |  **1,030,578** |     16,000,000 |
| ETH→Cosmos | `MsgUpdateClient` (wasm LC verify)  |    **822,220** |      2,000,000 |
| ETH→Cosmos | `MsgAcknowledgement`                |    **173,268** |      2,000,000 |

## Latency breakdown

### Prover (Groth16, bucket n=16, ICICLE/CUDA)

| Stage             |     Duration |
| ----------------- | -----------: |
| Build witness     |      2.55 ms |
| `groth16.Prove`   |  **3.520 s** |
| Local verify      |      1.19 ms |
| **Total**         |  **3.524 s** |

gnark/icicle internal:
- Hint generation: `272.06 ms`
- Solver: `1.720 s`
- Prover core: `1.278 s`

### Ethereum tx submission

| Tx              | Submit (RPC send) | Wait (mining) |   Total |
| --------------- | ----------------: | ------------: | ------: |
| `updateClient`  |           2.44 ms |        2.00 s | 2.00 s |
| `recvPacket`    |         718.59 µs |        2.00 s | 2.00 s |

### Cosmos tx submission

| Tx                              | Broadcast (CheckTx + DeliverTx) |   Total |
| ------------------------------- | ------------------------------: | ------: |
| `MsgUpdateClient` (batch=1)     |                          4.97 s | 4.97 s |
| `MsgAcknowledgement`            |                          5.38 s | 5.38 s |

### Beacon finality wait (ETH side)

- Event block: 201 (emitted at 09:56:33)
- Finalized at: 09:59:34 (block 224 >= 201)
- **18 polls × 10 s = 180 s**.

### One-time prover preload

| Bucket | Artifact bytes | Load time |
| -----: | -------------: | --------: |
| `n=4`  |       ~273.3 MB |      ~6 s |
| `n=8`  |       ~435.9 MB |     ~10 s |
| `n=16` |       ~701.0 MB |     ~14 s |
| **Total** | **~1.41 GB** | **~30 s** |

Loaded artifacts:
- `n=4`: `r1cs.bin` 72.7 MB, `pk.bin` 200.3 MB, `vk.bin` 560 B
- `n=8`: `r1cs.bin` 118.1 MB, `pk.bin` 317.8 MB, `vk.bin` 560 B
- `n=16`: `r1cs.bin` 215.8 MB, `pk.bin` 485.2 MB, `vk.bin` 560 B

## Comparison vs prior bucket-16 run

| Metric                              | n=16 (2026-05-19) | **n=16 ICICLE (2026-05-25)** |
| ----------------------------------- | ----------------: | ----------------------------: |
| Constraints                         |         1,685,201 |                     1,685,171 |
| `groth16.Prove`                     |           5.178 s |                   **3.520 s** |
| `ICS26Router.updateClient` gas      |         2,792,640 |                 **1,709,704** |
| `ICS26Router.recvPacket` gas        |         1,033,555 |                   1,030,578 |
| `MsgUpdateClient` gas (Cosmos)      |           822,569 |                     822,220 |
| `MsgAcknowledgement` gas            |           173,237 |                     173,268 |
| Beacon finality wait                |             150 s |                       180 s |
| E2E wall-clock                      |            ~177 s |                      ~222 s |

## Observations

1. **Proof generation improved materially**: `3.52 s` vs `5.18 s` on the prior bucket-16 run, consistent with ICICLE/CUDA acceleration.
2. **`updateClient` gas dropped sharply**: `2.79M → 1.71M` at the same bucket size, which points to contract/runtime improvements beyond prover acceleration alone.
3. **`recvPacket` and Cosmos-side gas remained flat**: the bucket-sensitive path is still overwhelmingly `updateClient`.
4. **Wall-clock is still dominated by beacon finality**: even with faster proving and lower ETH gas, `180 s` of finality wait overwhelms all other savings.
5. **Artifact preload is now substantial**: preloading three buckets costs about `1.41 GB` of reads and `~30 s` one-time startup latency, which is acceptable for a long-lived relayer but expensive for short benchmark runs.

## Conclusions

1. **Bucket 16 is comfortably viable** on the current 16M ETH tx cap after the latest changes.
2. **ICICLE/GPU meaningfully reduces prover latency**, but it does not change the dominant end-to-end bottleneck, which remains ETH beacon finality.
3. **The latest code path appears materially more gas-efficient on `updateClient`** than the previous bucket-16 benchmark, and this should be validated again at `n=32` and `n=64`.
4. **For perf testing, startup policy matters**: preloading `n=4/8/16` improves steady-state readiness but skews cold-start measurements.

## Raw log artifacts captured

- `[bench][prover]`  → witness / prove / verify / total
- `[bench][eth]`     → label / gasUsed / submit / wait / total
- `[bench][cosmos]`  → msg type / gasWanted / gasUsed / broadcast / total
- `[NewProver]`      → artifact sizes / bucket preload timings / backend selection
