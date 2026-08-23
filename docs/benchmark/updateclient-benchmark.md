# UpdateClient benchmark

Current values are the `bucket= ... gas= ...` logs emitted by
`test/light-clients/spectre/UpdateClientGasTest.t.sol`.

## Command

```bash
forge test --match-contract UpdateClientGasTest -vv
```

## Results

| bucket | gas |
|---|---:|
| 4 | 246,062 |
| 8 | 397,663 |
| 16 | 658,934 |
| 32 | 947,139 |
| 64 | 1,780,438 |
| 128 | 3,528,775 |
