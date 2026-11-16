package main

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"relayer/transaction"
)

// The one-shot commands sign with the keys `start` signs with, so they take the
// same locks (#575). Without them, update-client beside a running start on one
// ETH_PRIVATE_KEY hands out a nonce start already used, and update-eth-client
// beside it races the Cosmos account sequence; neither error names the other
// process.
//
// Each command locks only what it signs, after its cheap config checks and
// before any proving or signing. A refusal names the lock file: stop the other
// process, or give the command its own key.
//
//	create-clients-cosmos, update-eth-client  Cosmos signer (acquireCosmosSignerLock)
//	create-clients-eth, update-client         ETH_PRIVATE_KEY on the module's eth_rpc_url
//	submit-misbehaviour                       MISBEHAVIOUR_PRIVATE_KEY on the same endpoint

// acquireOneShotEVMSignerLock locks address's nonce domain on the chain rpcURL
// serves, and holds handler to that chain, as start does.
func acquireOneShotEVMSignerLock(stdCtx context.Context, label, rpcURL string, address common.Address, handler *transaction.Handler) (release func(), err error) {
	keys, err := evmSignerLockKeysForEndpoints(contextOrBackground(stdCtx),
		[]evmSignerLockEndpoint{{label: label, rpcURL: rpcURL}}, address, evmChainIDForSignerLock)
	if err != nil {
		return nil, err
	}
	release, err = acquireEVMSignerLockKeys(keys)
	if err != nil {
		return nil, err
	}
	handler.AllowEVMChainIDs(lockedChainIDs(keys))
	return release, nil
}

// oneShotEVMEndpointLabel names the module whose eth_rpc_url a one-shot command
// submits to: selectCosmosSource moves cosmos_to_l2 into CosmosToEthConfig when
// the file has no cosmos_to_eth module.
func oneShotEVMEndpointLabel(cfg *appConfig) string {
	if !cfg.HasCosmosToEth && cfg.HasCosmosToL2 {
		return "cosmos_to_l2.eth_rpc_url"
	}
	return "cosmos_to_eth.eth_rpc_url"
}
