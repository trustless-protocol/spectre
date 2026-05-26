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
| 4 | 292,912 |
| 8 | 504,175 |
| 16 | 868,891 |
| 32 | 1,260,993 |
| 64 | 2,408,551 |
| 128 | 4,795,736 |
