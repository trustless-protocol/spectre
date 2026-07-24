#[cfg(feature = "rpc")]
alloy_sol_types::sol!(
    #[sol(rpc)]
    #[allow(clippy::nursery, clippy::too_many_arguments)]
    groth16_ics07_tendermint,
    "../../abi/bytecode/Groth16ICS07Tendermint.json"
);

// NOTE: The riscv program won't compile with the `rpc` features.
#[cfg(not(feature = "rpc"))]
alloy_sol_types::sol!(
    groth16_ics07_tendermint,
    "../../abi/Groth16ICS07Tendermint.json"
);
