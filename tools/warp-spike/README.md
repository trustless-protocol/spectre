# warp-spike

Measures the one assumption the Avalanche warp light-client design depends on:
that a ≥67% stake-weighted BLS quorum of the **primary network** can be
aggregated for C-Chain block-hash payloads **from outside Avalanche**. Every
aggregate is independently re-verified with avalanchego's own
`BitSetSignature.Verify` against the canonical validator set fetched via
`platform.getAllValidatorsAt` — a success is a verified aggregate, not a
trusted HTTP response.

Standalone module on purpose: it imports avalanchego, which `relayer/go.mod`
must not.

## Usage

Run Ava Labs' signature-aggregator sidecar first (release binary from
github.com/ava-labs/icm-services, config pointing `p-chain-api`/`info-api` at
the target network), then:

```bash
go run . -base-url https://api.avax-test.network -aggregator http://127.0.0.1:18080 \
  -count 20 -interval 30s -csv spike.csv
```

## Measured results — Fuji, 2026-09-21 (6 attempts, this tool, sandbox host)

- **6/6 attempts produced a verified 67% aggregate.**
- Latency: 1904–1995 ms per aggregate (aggregator warm, ~25 s after boot).
- 12 signers reached 67.13% of stake; canonical set 73 entries.
- Validator-set crawl the same day: Fuji 75 validators / mainnet 597, **all**
  with registered BLS keys (0% keyless weight dilution); 67% needs 13 signers
  on Fuji, 78 on mainnet. Pinned-set state ≈ 597 × 56 B ≈ 33 KB on mainnet.
- Semantics settled: the C-Chain `Hash` payload covers the **EVM block hash**
  (validators signed a payload built from `eth_getBlockByNumber().hash`), and
  public API nodes do NOT serve `warp_*` methods — collection is ACP-118 p2p
  only, via the aggregator.
- The aggregator connected outbound to validator p2p (port 9651) from a plain
  container network without special setup.

Numbers above are single-day, single-host, small-sample — rerun with
`-count 100 -interval 60s` (or longer) before quoting anything externally, and
rerun against mainnet (`-base-url https://api.avax.network`) before a mainnet
commitment.
