package arbitrum

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func TestLoadDaemonConfigResolvesStatePathAndLoadsNitroEndpoints(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "attestor.json")
	data := `{
		"grpc_listen_address":"127.0.0.1:50051",
		"runtime_poll_interval":"2s",
		"attestation_head":"unsafe",
		"disable_derived_roots":true,
		"derived_attestation_gap_blocks":12,
		"max_derived_roots":34,
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
		"nitro_rpc_url":"https://arbitrum.example",
		"nitro_ws_url":"wss://arbitrum.example"
	}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := LoadDaemonConfig(path)
	if err != nil {
		t.Fatalf("load daemon config: %v", err)
	}
	if config.AttestorStatePath != filepath.Join(dir, "data", "attested-roots.json") {
		t.Fatalf("attestor state path: got %q", config.AttestorStatePath)
	}
	if config.NitroRPCURL != "https://arbitrum.example" {
		t.Fatalf("Nitro RPC URL: got %q", config.NitroRPCURL)
	}
	if config.NitroWSURL != "wss://arbitrum.example" {
		t.Fatalf("Nitro WebSocket URL: got %q", config.NitroWSURL)
	}
	attestationHead, err := config.NormalizedAttestationHead()
	if err != nil {
		t.Fatalf("attestation head: %v", err)
	}
	if attestationHead != RunModeUnsafe || !config.DisableDerivedRoots || config.DerivedAttestationGap() != 12 || config.DerivedRootLimit() != 34 {
		t.Fatalf("derived attestation config: head=%s gap=%d max=%d", attestationHead, config.DerivedAttestationGap(), config.DerivedRootLimit())
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

func TestDaemonConfigRequiresNitroHTTPAndWebSocketEndpoints(t *testing.T) {
	config := validDaemonConfig(t)
	config.NitroRPCURL = "wss://arbitrum.example"
	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "nitro_rpc_url") {
		t.Fatalf("validate Nitro RPC URL: got %v", err)
	}

	config = validDaemonConfig(t)
	config.NitroWSURL = "https://arbitrum.example"
	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "nitro_ws_url") {
		t.Fatalf("validate Nitro WebSocket URL: got %v", err)
	}
}

func TestLoadDaemonConfigRejectsManagedNitroFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "attestor.json")
	data := `{
		"nitro_rpc_url":"https://arbitrum.example",
		"nitro_ws_url":"wss://arbitrum.example",
		"nitro_binary_path":"./nitro"
	}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	_, err := LoadDaemonConfig(path)
	if err == nil || !strings.Contains(err.Error(), "nitro_binary_path is no longer supported") {
		t.Fatalf("load mixed endpoint/process config: got %v", err)
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

func TestDaemonConfigRuntimeBackfillDefaults(t *testing.T) {
	config := validDaemonConfig(t)
	if config.BackfillMaxBlocks() != 2_048 {
		t.Fatalf("default backfill max blocks: got %d want 2048", config.BackfillMaxBlocks())
	}
	if config.BackfillConcurrency() != 8 {
		t.Fatalf("default backfill concurrency: got %d want 8", config.BackfillConcurrency())
	}
	if err := config.Validate(); err != nil {
		t.Fatalf("validate config without backfill knobs: %v", err)
	}

	config.RuntimeBackfillMaxBlocks = 512
	config.RuntimeBackfillConcurrency = 4
	if config.BackfillMaxBlocks() != 512 || config.BackfillConcurrency() != 4 {
		t.Fatalf(
			"explicit backfill knobs: got %d/%d want 512/4",
			config.BackfillMaxBlocks(),
			config.BackfillConcurrency(),
		)
	}
}

func TestDaemonConfigAttestationHeadDefaultsAndValidation(t *testing.T) {
	config := validDaemonConfig(t)
	head, err := config.NormalizedAttestationHead()
	if err != nil {
		t.Fatalf("default attestation head: %v", err)
	}
	if head != RunModeFinalized || config.DerivedAttestationGap() != 150 || config.DerivedRootLimit() != 1_000 {
		t.Fatalf("derived defaults: head=%s gap=%d max=%d", head, config.DerivedAttestationGap(), config.DerivedRootLimit())
	}

	config.AttestationHead = RunMode("confirmed")
	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "attestation_head") {
		t.Fatalf("validate attestation head: got %v", err)
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

func TestDaemonConfigRequiresBoLDStorageLayout(t *testing.T) {
	config := validDaemonConfig(t)
	config.AssertionsMappingSlot = ""
	if err := config.Validate(); err == nil ||
		!strings.Contains(err.Error(), "assertions_mapping_slot") {
		t.Fatalf("validate missing assertions mapping slot: got %v", err)
	}
}

func TestArbitrumSepoliaConfigPinsBoLDDeployment(t *testing.T) {
	config, err := LoadDaemonConfig("config.arbitrum-sepolia.json")
	if err != nil {
		t.Fatalf("load Arbitrum Sepolia config: %v", err)
	}
	if config.L1ChainID != 11_155_111 ||
		config.L2ChainID != 421_614 ||
		config.AttestationHead != RunModeFinalized ||
		!config.DisableDerivedRoots ||
		config.DerivedAttestationGap() != 150 ||
		config.NitroRPCURL != "https://arbitrum-sepolia-rpc.publicnode.com" ||
		config.NitroWSURL != "wss://arbitrum-sepolia-rpc.publicnode.com" ||
		common.HexToAddress(config.RollupCoreAddress) != common.HexToAddress(
			"0x042B2E6C5E99d4c521bd49beeD5E99651D9B0Cf4",
		) ||
		common.HexToHash(config.AssertionsMappingSlot) != common.HexToHash("0x75") ||
		config.AssertionStatusOffset != 25 ||
		config.AssertionStartBlock != 11_379_731 {
		t.Fatalf("unexpected Arbitrum Sepolia config: %+v", config)
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
		NitroRPCURL:           "https://arbitrum.example",
		NitroWSURL:            "wss://arbitrum.example",
	}
}
