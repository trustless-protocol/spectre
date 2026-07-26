package arbitrum

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadDaemonConfigResolvesPersistentNitroPaths(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "attestor.json")
	data := `{
		"grpc_listen_address":"127.0.0.1:50051",
		"runtime_poll_interval":"2s",
		"src_chain":"arbitrum-one",
		"l1_rpc_url":"https://ethereum.example",
		"l1_chain_id":1,
		"l2_chain_id":42161,
		"rollup_core_address":"0x0000000000000000000000000000000000000001",
		"assertions_mapping_slot":"0x0000000000000000000000000000000000000000000000000000000000000076",
		"assertion_status_offset":25,
		"assertion_start_block":123,
		"assertion_poll_interval":"3s",
		"assertion_max_block_range":1000,
		"attestor_state_path":"data/attested-roots.json",
		"nitro_binary_path":"bin/nitro",
		"nitro_binary_sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"nitro_arguments":["--conf.file=nitro.json"],
		"nitro_work_dir":"nitro-work",
		"nitro_data_dir":"data/nitro-chain",
		"nitro_ipc_path":"run/nitro.ipc",
		"nitro_startup_timeout":"2m",
		"nitro_shutdown_timeout":"30s"
	}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := LoadDaemonConfig(path)
	if err != nil {
		t.Fatalf("load daemon config: %v", err)
	}
	if config.NitroBinaryPath != filepath.Join(dir, "bin", "nitro") {
		t.Fatalf("Nitro binary path: got %q", config.NitroBinaryPath)
	}
	if config.NitroWorkDir != filepath.Join(dir, "nitro-work") {
		t.Fatalf("Nitro work directory: got %q", config.NitroWorkDir)
	}
	if config.NitroDataDir != filepath.Join(dir, "data", "nitro-chain") {
		t.Fatalf("Nitro data directory: got %q", config.NitroDataDir)
	}
	if config.NitroIPCPath != filepath.Join(dir, "run", "nitro.ipc") {
		t.Fatalf("Nitro IPC path: got %q", config.NitroIPCPath)
	}
	if config.AttestorStatePath != filepath.Join(dir, "data", "attested-roots.json") {
		t.Fatalf("attestor state path: got %q", config.AttestorStatePath)
	}
	runtimePollInterval, err := config.RuntimePollDuration()
	if err != nil {
		t.Fatalf("runtime poll interval: %v", err)
	}
	if runtimePollInterval != 2*time.Second {
		t.Fatalf("runtime poll interval: got %s want %s", runtimePollInterval, 2*time.Second)
	}
	assertionPollInterval, err := config.AssertionPollDuration()
	if err != nil {
		t.Fatalf("assertion poll interval: %v", err)
	}
	if assertionPollInterval != 3*time.Second {
		t.Fatalf("assertion poll interval: got %s want %s", assertionPollInterval, 3*time.Second)
	}
	if config.AssertionBlockRange() != 1000 {
		t.Fatalf("assertion block range: got %d want 1000", config.AssertionBlockRange())
	}
}

func TestDaemonConfigRejectsOwnedNitroArguments(t *testing.T) {
	for _, argument := range []string{
		"--ipc.path=/tmp/unowned.ipc",
		"--persistent.chain=/tmp/unowned-chain",
	} {
		t.Run(argument, func(t *testing.T) {
			config := validDaemonConfig(t)
			config.NitroArguments = []string{argument}
			if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "must not set") {
				t.Fatalf("validate owned Nitro argument: got %v", err)
			}
		})
	}
}

func TestDaemonConfigRejectsInvalidGRPCAddress(t *testing.T) {
	config := validDaemonConfig(t)
	config.GRPCListenAddress = "127.0.0.1"
	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "grpc_listen_address") {
		t.Fatalf("validate gRPC address: got %v", err)
	}
}

func TestDaemonConfigRuntimePollInterval(t *testing.T) {
	config := validDaemonConfig(t)
	interval, err := config.RuntimePollDuration()
	if err != nil {
		t.Fatalf("default runtime poll interval: %v", err)
	}
	if interval != 15*time.Second {
		t.Fatalf("default runtime poll interval: got %s want %s", interval, 15*time.Second)
	}

	config.RuntimePollInterval = "0s"
	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "runtime_poll_interval") {
		t.Fatalf("validate runtime poll interval: got %v", err)
	}
}

func TestDaemonConfigAssertionDefaults(t *testing.T) {
	config := validDaemonConfig(t)
	interval, err := config.AssertionPollDuration()
	if err != nil {
		t.Fatalf("default assertion poll interval: %v", err)
	}
	if interval != 12*time.Second {
		t.Fatalf("default assertion poll interval: got %s want %s", interval, 12*time.Second)
	}
	if config.AssertionBlockRange() != 2_000 {
		t.Fatalf("default assertion block range: got %d want 2000", config.AssertionBlockRange())
	}

	config.AssertionPollInterval = "0s"
	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "assertion_poll_interval") {
		t.Fatalf("validate assertion poll interval: got %v", err)
	}
}

func validDaemonConfig(t *testing.T) DaemonConfig {
	t.Helper()
	dir := t.TempDir()
	return DaemonConfig{
		GRPCListenAddress:     "127.0.0.1:50051",
		SrcChain:              "arbitrum-one",
		L1RPCURL:              "https://ethereum.example",
		L1ChainID:             1,
		L2ChainID:             42161,
		RollupCoreAddress:     "0x0000000000000000000000000000000000000001",
		AssertionsMappingSlot: "0x" + strings.Repeat("0", 62) + "76",
		AssertionStatusOffset: 25,
		AssertionStartBlock:   1,
		AttestorStatePath:     filepath.Join(dir, "attested-roots.json"),
		NitroBinaryPath:       "nitro",
		NitroBinarySHA256:     strings.Repeat("a", 64),
		NitroWorkDir:          dir,
		NitroDataDir:          filepath.Join(dir, "chain"),
		NitroIPCPath:          filepath.Join(dir, "nitro.ipc"),
		NitroStartupTimeout:   "2m",
		NitroShutdownTimeout:  "30s",
	}
}
