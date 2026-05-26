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
| 4 | 264,146 |
| 8 | 423,383 |
| 16 | 697,100 |
| 32 | 997,919 |
| 64 | 1,869,062 |
| 128 | 3,693,087 |
