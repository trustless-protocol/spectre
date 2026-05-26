# `test_gas_n16` flamegraph notes

## Command

```bash
forge test -vvvv --match-test test_gas_n16 --flamegraph
```

## Artifact

- Flamegraph SVG: `cache/flamegraph_UpdateClientGasTest_test_gas_n16.svg`

## Scope

The numbers below are read from the flamegraph call tree for `UpdateClientGasTest.test_gas_n16()`.

This is not an isolated micro-benchmark for each component:

- percentages are relative to the full test gas
- deployment/setup gas is included
- repeated frames can appear more than once in the same execution tree

## Top-level totals

| frame | gas | share |
|---|---:|---:|
| `all` | 4,701,366 | 100.00% |
| `UpdateClientGasTest.test_gas_n16()` | 4,701,366 | 100.00% |
| `new Groth16ICS07Tendermint` | 3,441,477 | 73.20% |
| `Groth16ICS07Tendermint.updateClient(bytes)` | 900,385 | 19.15% |

## Notable frames inside the executed path

| frame | gas | share |
|---|---:|---:|
| `UpdateClient.updateClient(...)` | 445,804 | 9.48% |
| `WrapperVerifier.verifyBatchProof(...)` | 255,890 | 5.44% |
| `WrapperVerifier::_hashWitness` | 247,622 | 5.27% |
| `WrapperVerifier::_voteSignBytes` | 215,929 | 4.59% |
| `Header.hashValSet(uint8)` | 268,118 | 5.70% |
| `Header.hashValSet(uint8)` | 136,559 | 2.90% |
| `Header.hashHeader(uint8)` | 47,897 | 1.02% |
| `Header.hashHeader(uint8)` | 47,897 | 1.02% |

## Related harness output

From the same `forge test` run, the benchmark log for bucket `16` is:

```text
bucket= 16   gas= 909463
```

This log is usually the better number for comparing update-path gas across buckets. The flamegraph is more useful for seeing where gas concentrates inside the executed call tree.
