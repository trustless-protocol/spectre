package main

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"relayer/services"
)

type l2TimeoutReturnPath struct {
	svc  *services.Services
	deps services.RelayDeps
}

func (p l2TimeoutReturnPath) validate() error {
	if p.svc == nil {
		return fmt.Errorf("l2_to_cosmos timeout return path is missing services")
	}
	if p.deps.Cosmos.CosmosClient() == nil {
		return fmt.Errorf("l2_to_cosmos timeout return path has nil Cosmos client")
	}
	if p.deps.EVM.EthClient() == nil {
		return fmt.Errorf("l2_to_cosmos timeout return path has nil L2 client")
	}
	if p.deps.EVM.RouterContract() == nil {
		return fmt.Errorf("l2_to_cosmos timeout return path has nil L2 router address")
	}
	if p.deps.EVM.SpectreClientContract() == nil {
		return fmt.Errorf("l2_to_cosmos timeout return path has nil SpectreClient address")
	}
	return nil
}

type l2TimeoutReturnPathConfig struct {
	cfg  cosmosToEthConfig
	path l2TimeoutReturnPath
}

type l2TimeoutReturnPathKey struct {
	l2RPCURL string
	tmRPCURL string
	router   string
}

func l2TimeoutReturnPathKeyForSource(cfg l2ToCosmosConfig) (l2TimeoutReturnPathKey, error) {
	router, err := l2RouterFromProfile(cfg.RollupProfile)
	if err != nil {
		return l2TimeoutReturnPathKey{}, err
	}
	return l2TimeoutReturnPathKey{
		l2RPCURL: cfg.L2RpcUrl,
		tmRPCURL: cfg.TmRpcUrl,
		router:   strings.ToLower(router.Hex()),
	}, nil
}

func l2TimeoutReturnPathKeyForDest(cfg cosmosToEthConfig) (l2TimeoutReturnPathKey, error) {
	if !common.IsHexAddress(cfg.ICS26Address) {
		return l2TimeoutReturnPathKey{}, fmt.Errorf("cosmos_to_l2.ics26_address %q is not a valid hex address", cfg.ICS26Address)
	}
	return l2TimeoutReturnPathKey{
		l2RPCURL: cfg.EthRpcUrl,
		tmRPCURL: cfg.TmRpcUrl,
		router:   strings.ToLower(common.HexToAddress(cfg.ICS26Address).Hex()),
	}, nil
}

func findL2TimeoutReturnPath(src l2ToCosmosConfig, paths []l2TimeoutReturnPathConfig) (l2TimeoutReturnPath, error) {
	key, err := l2TimeoutReturnPathKeyForSource(src)
	if err != nil {
		return l2TimeoutReturnPath{}, err
	}

	match := -1
	for i := range paths {
		candidate, err := l2TimeoutReturnPathKeyForDest(paths[i].cfg)
		if err != nil {
			return l2TimeoutReturnPath{}, err
		}
		if candidate != key {
			continue
		}
		if match != -1 {
			return l2TimeoutReturnPath{}, fmt.Errorf(
				"multiple cosmos_to_l2 return paths match l2_to_cosmos source (l2_rpc_url=%s tm_rpc_url=%s l2_router=%s)",
				key.l2RPCURL, key.tmRPCURL, key.router,
			)
		}
		match = i
	}
	if match == -1 {
		return l2TimeoutReturnPath{}, fmt.Errorf(
			"no matching cosmos_to_l2 return path for l2_to_cosmos source (l2_rpc_url=%s tm_rpc_url=%s l2_router=%s)",
			key.l2RPCURL, key.tmRPCURL, key.router,
		)
	}
	return paths[match].path, nil
}

func validateL2TimeoutReturnPathConfigs(sources []l2ToCosmosConfig, dests []cosmosToEthConfig) error {
	if len(sources) == 0 {
		return nil
	}
	paths := make([]l2TimeoutReturnPathConfig, 0, len(dests))
	for i := range dests {
		paths = append(paths, l2TimeoutReturnPathConfig{cfg: dests[i]})
	}
	for i := range sources {
		if _, err := findL2TimeoutReturnPath(sources[i], paths); err != nil {
			return err
		}
	}
	return nil
}
