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
| 4 | 256,091 |
| 8 | 414,601 |
| 16 | 687,107 |
| 32 | 986,714 |
| 64 | 1,854,221 |
| 128 | 3,670,974 |
