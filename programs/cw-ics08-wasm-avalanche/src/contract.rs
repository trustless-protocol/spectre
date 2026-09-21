//! Avalanche C-Chain warp ICS-08 Wasm client entrypoints.

#![allow(missing_docs)]

avalanche_light_client::avalanche_client_entrypoints!();

#[cfg(test)]
mod tests {
    use avalanche_light_client::msg::MigrateMsg;
    use cosmwasm_std::{DepsMut, Env, Response};

    #[test]
    fn exposes_the_warp_client_entrypoints() {
        let _: fn(
            DepsMut<avalanche_light_client::bls::AvalancheCustomQuery>,
            Env,
            MigrateMsg,
        ) -> Result<Response, avalanche_light_client::error::Error> = super::migrate;
    }
}
