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
| 4 | 254,930 |
| 8 | 407,427 |
| 16 | 669,744 |
| 32 | 959,163 |
| 64 | 1,796,106 |
| 128 | 3,551,731 |
