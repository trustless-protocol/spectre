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
| 4 | 281,427 |
| 8 | 487,703 |
| 16 | 844,007 |
| 32 | 1,227,696 |
| 64 | 2,350,015 |
| 128 | 4,686,722 |
