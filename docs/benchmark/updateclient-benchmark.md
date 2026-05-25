# UpdateClient benchmark

Current values are the `bucket= ... gas= ...` logs emitted by
`test/solidity-ibc/UpdateClientGasTest.t.sol`.

## Command

```bash
forge test --match-path test/solidity-ibc/UpdateClientGasTest.t.sol --match-test test_gas_n4 -vv
forge test --match-path test/solidity-ibc/UpdateClientGasTest.t.sol --match-test test_gas_n8 -vv
forge test --match-path test/solidity-ibc/UpdateClientGasTest.t.sol --match-test test_gas_n16 -vv
forge test --match-path test/solidity-ibc/UpdateClientGasTest.t.sol --match-test test_gas_n32 -vv
forge test --match-path test/solidity-ibc/UpdateClientGasTest.t.sol --match-test test_gas_n64 -vv
forge test --match-path test/solidity-ibc/UpdateClientGasTest.t.sol --match-test test_gas_n128 -vv
```

## Results

| bucket | gas |
|---|---:|
| 4 | 305,370 |
| 8 | 526,640 |
| 16 | 909,463 |
| 32 | 1,323,663 |
| 64 | 2,536,467 |
| 128 | 5,078,009 |
