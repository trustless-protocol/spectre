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
| 4 | 266,970 |
| 8 | 429,255 |
| 16 | 708,052 |
| 32 | 1,013,951 |
| 64 | 1,900,334 |
| 128 | 3,754,839 |

