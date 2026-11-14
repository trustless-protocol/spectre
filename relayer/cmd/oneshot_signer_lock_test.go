package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"relayer/rpcmock"
)

const (
	oneShotTestETHKey          = "a81f9eb900c02f35d28ec80d1da79dde52cc4cc37d6fbaef5768ab59f32cbfdd"
	oneShotTestMisbehaviourKey = "4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"
)

// isolateOneShotRun keeps .env files, operator overrides and the machine's real
// lock directory out of a one-shot command run.
func isolateOneShotRun(t *testing.T) {
	t.Helper()
	t.Chdir(t.TempDir())
	for _, name := range []string{"COSMOS_WASM_CLIENT_ID", "ICS26_CLIENT_ID", "ROLE_MANAGER", "ETH_MISBEHAVIOUR_GAS_LIMIT", "COSMOS_PRIVATE_KEY"} {
		t.Setenv(name, "")
	}
	t.Setenv("ETH_PRIVATE_KEY", oneShotTestETHKey)
	t.Setenv("MISBEHAVIOUR_PRIVATE_KEY", oneShotTestMisbehaviourKey)
	// A prover that cannot load: a command that gets past its lock stops there.
	t.Setenv("PROVER_BIN_DIR", filepath.Join(t.TempDir(), "no-such-bin"))
	dir := t.TempDir()
	restore := signerLockDir
	signerLockDir = func() (string, error) { return dir, nil }
	t.Cleanup(func() { signerLockDir = restore })
}

func addressOfKey(t *testing.T, hexKey string) common.Address {
	t.Helper()
	key, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		t.Fatal(err)
	}
	return crypto.PubkeyToAddress(key.PublicKey)
}

// oneShotTestConfig is a cosmos_to_eth config against stand-in nodes that answer
// the create-clients preflight. The EVM stand-in reports chain id 1.
func oneShotTestConfig(t *testing.T) string {
	t.Helper()
	evm := rpcmock.NewEVMNode(t)
	cosmos := rpcmock.NewCosmosNode(t)
	cosmos.Handle("status", func(json.RawMessage) (any, error) {
		return &coretypes.ResultStatus{SyncInfo: coretypes.SyncInfo{LatestBlockHeight: 10, LatestBlockTime: time.Now()}}, nil
	})
	path := filepath.Join(t.TempDir(), "config.json")
	config := fmt.Sprintf(`{"modules":[{"name":"cosmos_to_eth","src_chain":"cosmoshub-1","config":{`+
		`"tm_rpc_url":%q,"eth_rpc_url":%q,`+
		`"ics26_address":"0x80741a37e3644612f0465145c9709a90b6d77ee3","ics26_client_id":"cosmoshub-1",`+
		`"cosmos_wasm_client_id":"08-wasm-0"}}]}`, cosmos.URL(), evm.URL())
	if err := os.WriteFile(path, []byte(config), configFilePerm); err != nil {
		t.Fatal(err)
	}
	return path
}

func oneShotTestEvidence(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "evidence.json")
	data := `{"client_id":"cosmoshub-1",` +
		`"header_1":{"trusted_height":{"revision_number":"1","revision_height":"10"}},` +
		`"header_2":{"trusted_height":{"revision_number":"1","revision_height":"10"}}}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func runOneShot(t *testing.T, cmd *cobra.Command, args ...string) error {
	t.Helper()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return cmd.ExecuteContext(ctx)
}

// holdEVMSignerLock stands in for a `start` signing chain 1 with address.
func holdEVMSignerLock(t *testing.T, address common.Address) func() {
	t.Helper()
	release, err := acquireEVMSignerLockKeys([]evmSignerLockKey{{chainID: "1", address: address}})
	if err != nil {
		t.Fatalf("hold the lock a running start would hold: %v", err)
	}
	return release
}

// #575: only start took the signer locks, so a one-shot command on start's key
// raced its nonce or account sequence.
func TestOneShotCommandsTakeTheEVMSignerLock(t *testing.T) {
	ethSigner := addressOfKey(t, oneShotTestETHKey)
	incidentSigner := addressOfKey(t, oneShotTestMisbehaviourKey)
	refusal := func(address common.Address) string {
		return fmt.Sprintf("already signing EVM transactions with %s on chain id 1", address.Hex())
	}

	for _, tc := range []struct {
		name    string
		command func() *cobra.Command
		args    func(config string) []string
		held    common.Address
		refused bool
	}{
		{name: "update-client", command: func() *cobra.Command { return UpdateClient(zap.NewNop()) },
			args: func(c string) []string { return []string{"--config", c} }, held: ethSigner, refused: true},
		{name: "create-clients-eth", command: func() *cobra.Command { return CreateClientsEth(zap.NewNop()) },
			args: func(c string) []string { return []string{"--config", c} }, held: ethSigner, refused: true},
		{name: "submit-misbehaviour on the incident key", command: func() *cobra.Command { return SubmitMisbehaviour(zap.NewNop()) },
			args: func(c string) []string { return []string{"--config", c, "--evidence", oneShotTestEvidence(t)} }, held: incidentSigner, refused: true},
		// The incident key is its own nonce domain: a start on ETH_PRIVATE_KEY does not block it.
		{name: "submit-misbehaviour beside a start on ETH_PRIVATE_KEY", command: func() *cobra.Command { return SubmitMisbehaviour(zap.NewNop()) },
			args: func(c string) []string { return []string{"--config", c, "--evidence", oneShotTestEvidence(t)} }, held: ethSigner},
		// A dry run signs nothing.
		{name: "submit-misbehaviour --dry-run", command: func() *cobra.Command { return SubmitMisbehaviour(zap.NewNop()) },
			args: func(c string) []string {
				return []string{"--dry-run", "--config", c, "--evidence", oneShotTestEvidence(t)}
			}, held: incidentSigner},
	} {
		t.Run(tc.name, func(t *testing.T) {
			isolateOneShotRun(t)
			release := holdEVMSignerLock(t, tc.held)
			err := runOneShot(t, tc.command(), tc.args(oneShotTestConfig(t))...)
			release()
			if err == nil {
				t.Fatal("command succeeded against stand-in nodes; the test expects it to stop at the lock or the prover")
			}
			if tc.refused && !strings.Contains(err.Error(), refusal(tc.held)) {
				t.Fatalf("err = %v, want the signer lock refusal naming the held domain (%s)", err, refusal(tc.held))
			}
			// Past the lock, the next step is the prover, which cannot load here.
			if !tc.refused && !strings.Contains(err.Error(), "failed to load prover") {
				t.Fatalf("err = %v, want the run past the lock, stopped at the prover", err)
			}
			// Whatever stopped it, the command let its own locks go.
			again, lockErr := acquireEVMSignerLockKeys([]evmSignerLockKey{{chainID: "1", address: ethSigner}, {chainID: "1", address: incidentSigner}})
			if lockErr != nil {
				t.Fatalf("a lock is still held after the command returned: %v", lockErr)
			}
			again()
		})
	}
}

func TestOneShotCommandsTakeTheCosmosSignerLock(t *testing.T) {
	const signer = "cosmos1testsigner"
	newRun := func(t *testing.T) {
		isolateOneShotRun(t)
		restore := cosmosSignerAddressForLock
		cosmosSignerAddressForLock = func() (string, error) { return signer, nil }
		t.Cleanup(func() { cosmosSignerAddressForLock = restore })
	}
	const refusal = "already signing Cosmos transactions with " + signer

	t.Run("update-eth-client", func(t *testing.T) {
		newRun(t)
		release, err := acquireCosmosSignerLock()
		if err != nil {
			t.Fatal(err)
		}
		configPath, node := ethClientUpdateTestConfig(t, false)
		tx := &ethClientUpdateTestHandler{}
		err = runOneShot(t, newUpdateEthClientCommand(zap.NewNop(), tx), "--config", configPath)
		if err == nil || !strings.Contains(err.Error(), refusal) {
			t.Fatalf("err = %v, want the signer lock refusal", err)
		}
		if tx.submitCalls != 0 || node.Calls("abci_query") != 0 {
			t.Fatalf("refused run still read (%d) or submitted (%d)", node.Calls("abci_query"), tx.submitCalls)
		}

		// With the other process gone the same run goes through, and lets the lock go.
		release()
		tx = &ethClientUpdateTestHandler{}
		if err := runOneShot(t, newUpdateEthClientCommand(zap.NewNop(), tx), "--config", configPath); err != nil {
			t.Fatalf("run without a competing process: %v", err)
		}
		if tx.submitCalls != 1 {
			t.Fatalf("submit calls = %d, want 1", tx.submitCalls)
		}
		again, err := acquireCosmosSignerLock()
		if err != nil {
			t.Fatalf("lock still held after the command returned: %v", err)
		}
		again()
	})

	t.Run("create-clients-cosmos", func(t *testing.T) {
		newRun(t)
		release, err := acquireCosmosSignerLock()
		if err != nil {
			t.Fatal(err)
		}
		defer release()
		err = runOneShot(t, CreateClientsCosmos(zap.NewNop()), "--config", oneShotTestConfig(t), "--wasm-checksum", "deadbeef")
		if err == nil || !strings.Contains(err.Error(), refusal) {
			t.Fatalf("err = %v, want the signer lock refusal", err)
		}
	})
}
