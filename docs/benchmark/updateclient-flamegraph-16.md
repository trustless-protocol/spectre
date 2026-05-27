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
| `all` | 4,353,287 | 100.00% |
| `UpdateClientGasTest.test_gas_n16()` | 4,353,287 | 100.00% |
| `new Groth16ICS07Tendermint` | 3,359,820 | 77.18% |
| `Groth16ICS07Tendermint.updateClient(bytes)` | 649,856 | 14.93% |

## Notable frames inside the executed path

| frame | gas | share |
|---|---:|---:|
| `UpdateClient.updateClient(...)` | 406,888 | 9.35% |
| `WrapperVerifier.verifyBatchProof(...)` | 85,666 | 1.97% |
| `WrapperVerifier::_hashWitness` | 77,398 | 1.78% |
| `WrapperVerifier::_writeVoteSignBytes` | 43,570 | 1.00% |
| `Header.hashValSet(uint8)` | 242,158 | 5.56% |
| `Header.hashValSet(uint8)` | 123,579 | 2.84% |
| `Header.hashHeader(uint8)` | 44,934 | 1.03% |
| `Header.hashHeader(uint8)` | 44,934 | 1.03% |

## Related harness output

From the same `forge test` run, the benchmark log for bucket `16` is:

```text
bucket= 16   gas= 658934
```

This log is usually the better number for comparing update-path gas across buckets. The flamegraph is more useful for seeing where gas concentrates inside the executed call tree.
