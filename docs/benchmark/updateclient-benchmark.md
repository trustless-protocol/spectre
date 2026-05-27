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
| 4 | 263,333 |
| 8 | 422,534 |
| 16 | 696,283 |
| 32 | 997,100 |
| 64 | 1,868,235 |
| 128 | 3,692,244 |
