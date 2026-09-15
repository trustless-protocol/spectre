package main

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TestEVMSignerLock pins the distinction that makes this guard different from
// the Cosmos signer lock: an EVM nonce belongs to (chain id, address), not to
// an address alone. The same key on two chains must therefore take two locks.
func TestEVMSignerLock(t *testing.T) {
	address := common.HexToAddress("0xd58ed941051839DB5ffe8735905881d3C7460cE1")

	setup := func(t *testing.T) {
		t.Helper()
		dir := t.TempDir()
		restoreDir := signerLockDir
		signerLockDir = func() (string, error) { return dir, nil }
		t.Cleanup(func() { signerLockDir = restoreDir })
	}

	t.Run("the second process with the same key on the same chain is refused", func(t *testing.T) {
		setup(t)
		key := evmSignerLockKey{chainID: "1", address: address}

		release, err := acquireEVMSignerLockKeys([]evmSignerLockKey{key})
		if err != nil {
			t.Fatalf("first process could not take the lock: %v", err)
		}
		defer release()

		_, err = acquireEVMSignerLockKeys([]evmSignerLockKey{key})
		if err == nil {
			t.Fatal("a second process took the same EVM nonce domain")
		}
		path, pathErr := evmSignerLockPathFor(key)
		if pathErr != nil {
			t.Fatalf("evmSignerLockPathFor: %v", pathErr)
		}
		for _, want := range []string{address.Hex(), "1", path, "ETH_PRIVATE_KEY"} {
			if !strings.Contains(err.Error(), want) {
				t.Fatalf("error %q does not mention %q", err, want)
			}
		}
	})

	t.Run("the same key on different chains can both start", func(t *testing.T) {
		setup(t)

		releaseOne, err := acquireEVMSignerLockKeys([]evmSignerLockKey{{chainID: "1", address: address}})
		if err != nil {
			t.Fatalf("first chain could not take the lock: %v", err)
		}
		defer releaseOne()

		releaseTwo, err := acquireEVMSignerLockKeys([]evmSignerLockKey{{chainID: "8453", address: address}})
		if err != nil {
			t.Fatalf("a process using the same key on another EVM chain must start: %v", err)
		}
		defer releaseTwo()
	})
}

func TestEVMSignerLockKeysCoverEveryEVMSubmitPath(t *testing.T) {
	address := common.HexToAddress("0xd58ed941051839DB5ffe8735905881d3C7460cE1")
	chainIDs := map[string]string{
		"http://eth":        "1",
		"http://l2":         "8453",   // Base
		"http://return-leg": "421614", // Arbitrum Sepolia
	}

	keys, err := evmSignerLockKeysWithChainID(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{{EthRpcUrl: "http://eth"}},
		CosmosToL2Configs:  []cosmosToEthConfig{{EthRpcUrl: "http://l2"}},
		L2ToCosmosConfigs:  []l2ToCosmosConfig{{L2RpcUrl: "http://return-leg"}},
	}, address, func(_ context.Context, endpoint string) (string, error) {
		return chainIDs[endpoint], nil
	})
	if err != nil {
		t.Fatalf("evmSignerLockKeys: %v", err)
	}
	if len(keys) != 3 {
		t.Fatalf("lock keys = %#v, want the L1, forward L2, and l2_to_cosmos return-leg L2", keys)
	}
	got := make(map[string]bool, len(keys))
	for _, key := range keys {
		if key.address != address {
			t.Fatalf("key address = %s, want %s", key.address.Hex(), address.Hex())
		}
		got[key.chainID] = true
	}
	for _, want := range []string{"1", "8453", "421614"} {
		if !got[want] {
			t.Fatalf("lock keys %v omit chain %s", got, want)
		}
	}
}

func TestEVMSignerLockKeysRefuseAnUnansweredEndpoint(t *testing.T) {
	_, err := evmSignerLockKeysWithChainID(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{{EthRpcUrl: "http://unavailable"}},
	}, common.HexToAddress("0xd58ed941051839DB5ffe8735905881d3C7460cE1"), func(context.Context, string) (string, error) {
		return "", errors.New("endpoint did not answer")
	})
	if err == nil {
		t.Fatal("startup accepted an EVM endpoint without a chain id and would have run without its nonce lock")
	}
	for _, want := range []string{"eth_chainId", "cosmos_to_eth[0]", "without the nonce-domain lock"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not explain %q", err, want)
		}
	}
}
