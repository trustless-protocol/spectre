# UpdateClient benchmark

Current values are the `bucket= ... gas= ...` logs emitted by
`test/solidity-ibc/UpdateClientGasTest.t.sol`.

## Command

```bash
forge test --match-contract UpdateClientGasTest -vv
```

## Results

| bucket | gas |
|---|---:|
| 4 | 289,482 |
| 8 | 496,485 |
| 16 | 854,000 |
| 32 | 1,238,901 |
| 64 | 2,364,856 |
| 128 | 4,708,835 |
