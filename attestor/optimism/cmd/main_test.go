package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestLoadConfigExample(t *testing.T) {
	t.Parallel()

	cfg, err := loadConfig(filepath.Join("..", "config.example.json"))
	if err != nil {
		t.Fatalf("loadConfig(config.example.json) error = %v", err)
	}
	if got := len(cfg.Sources); got != 1 {
		t.Fatalf("Sources len = %d, want 1", got)
	}
	op := cfg.Sources[0]
	if op.SrcChain != "op-mainnet" {
		t.Fatalf("op_source src_chain = %q, want op-mainnet", op.SrcChain)
	}
	if op.OpNodeRpcUrl == "" || op.DisputeGameFactory == "" || op.StatePath == "" {
		t.Fatalf("op_source fields not parsed: %+v", op)
	}
	if op.AttestationHead != "finalized" {
		t.Fatalf("op_source attestation_head = %q, want finalized", op.AttestationHead)
	}
	if cfg.Server.GrpcPort == 0 {
		t.Fatal("server.grpc_port not parsed")
	}
}

func TestLoadConfigRejectsDuplicateSources(t *testing.T) {
	t.Parallel()

	write := func(t *testing.T, body string) string {
		t.Helper()
		configPath := filepath.Join(t.TempDir(), "config.json")
		if err := os.WriteFile(configPath, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		return configPath
	}

	t.Run("duplicate src_chain", func(t *testing.T) {
		configPath := write(t, `{"modules":[`+
			`{"name":"op_source","src_chain":"op-a","config":{"state_path":"a.json"}},`+
			`{"name":"op_source","src_chain":"op-a","config":{"state_path":"b.json"}}`+
			`]}`)
		if _, err := loadConfig(configPath); err == nil || !strings.Contains(err.Error(), "duplicate op_source src_chain") {
			t.Fatalf("error = %v, want duplicate src_chain error", err)
		}
	})
	t.Run("duplicate state_path", func(t *testing.T) {
		configPath := write(t, `{"modules":[`+
			`{"name":"op_source","src_chain":"op-a","config":{"state_path":"same.json"}},`+
			`{"name":"op_source","src_chain":"op-b","config":{"state_path":"same.json"}}`+
			`]}`)
		if _, err := loadConfig(configPath); err == nil || !strings.Contains(err.Error(), "duplicate op_source state_path") {
			t.Fatalf("error = %v, want duplicate state_path error", err)
		}
	})
	t.Run("missing src_chain", func(t *testing.T) {
		configPath := write(t, `{"modules":[`+
			`{"name":"op_source","config":{"state_path":"a.json"}}`+
			`]}`)
		if _, err := loadConfig(configPath); err == nil || !strings.Contains(err.Error(), "src_chain") {
			t.Fatalf("error = %v, want missing src_chain error", err)
		}
	})
	t.Run("other modules ignored", func(t *testing.T) {
		configPath := write(t, `{"modules":[`+
			`{"name":"cosmos_to_eth","src_chain":"hub","config":{"tm_rpc_url":"http://x"}}`+
			`]}`)
		cfg, err := loadConfig(configPath)
		if err != nil || len(cfg.Sources) != 0 {
			t.Fatalf("(%+v, %v), want zero sources and no error", cfg, err)
		}
	})
}

func TestLoadConfigFailureMatrix(t *testing.T) {
	write := func(t *testing.T, body string) string {
		t.Helper()
		path := filepath.Join(t.TempDir(), "config.json")
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	if _, err := loadConfig(filepath.Join(t.TempDir(), "missing.json")); err == nil || !strings.Contains(err.Error(), "read config") {
		t.Fatalf("missing config error = %v", err)
	}
	if _, err := loadConfig(write(t, "{")); err == nil || !strings.Contains(err.Error(), "parse config") {
		t.Fatalf("malformed JSON error = %v", err)
	}
	if _, err := loadConfig(write(t, `{"modules":[{"name":"op_source","src_chain":"op","config":true}]}`)); err == nil || !strings.Contains(err.Error(), "parse op_source") {
		t.Fatalf("malformed op_source config error = %v", err)
	}
	if _, err := loadConfig(write(t, `{"modules":[{"name":"op_source","src_chain":"op","config":{"l1_ws_url":":"}}]}`)); err == nil || !strings.Contains(err.Error(), "l1_ws_url") {
		t.Fatalf("malformed L1 websocket URL error = %v", err)
	}
}

func TestValidateHexAddressFailureMatrix(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value string
		want  string
	}{
		{name: "empty", want: "empty"},
		{name: "malformed", value: "not-an-address", want: "invalid hex address"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateHexAddress(tc.value, "factory"); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("validateHexAddress error = %v, want %q", err, tc.want)
			}
		})
	}
	if err := validateHexAddress("0x0000000000000000000000000000000000000001", "factory"); err != nil {
		t.Fatalf("validateHexAddress valid address: %v", err)
	}
}

func TestBuildOpAttestorRejectsMissingSigningKeyBeforeNetworkDial(t *testing.T) {
	_, _, _, err := buildOpAttestor(context.Background(), zap.NewNop(), opSourceConfig{
		SrcChain:              "op-mainnet",
		DisputeGameFactory:    "0x0000000000000000000000000000000000000001",
		L2ChainID:             10,
		AttestationSigningKey: "",
	}, nil)
	if err == nil || !strings.Contains(err.Error(), "attestation signer") {
		t.Fatalf("buildOpAttestor missing signing key error = %v", err)
	}
}

func TestBuildOpAttestorRejectsInvalidInputBeforeNetworkDial(t *testing.T) {
	validSigningKey := strings.Repeat("11", 32)
	for _, tc := range []struct {
		name   string
		mutate func(*opSourceConfig)
		want   string
	}{
		{
			name:   "invalid factory address",
			mutate: func(c *opSourceConfig) { c.DisputeGameFactory = "not-an-address" },
			want:   "invalid hex address",
		},
		{
			name:   "missing L1 RPC",
			mutate: func(c *opSourceConfig) { c.L1RpcUrl = "" },
			want:   "l1_rpc_url",
		},
		{
			name:   "missing op-node RPC",
			mutate: func(c *opSourceConfig) { c.OpNodeRpcUrl = "" },
			want:   "op_node_rpc_url",
		},
		{
			name:   "missing state path",
			mutate: func(c *opSourceConfig) { c.StatePath = "" },
			want:   "state_path",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := opSourceConfig{
				SrcChain:              "op-mainnet",
				L1RpcUrl:              "https://ethereum.example",
				OpNodeRpcUrl:          "https://op-node.example",
				DisputeGameFactory:    "0x0000000000000000000000000000000000000001",
				L2ChainID:             10,
				AttestationSigningKey: validSigningKey,
				StatePath:             "attestor-state.json",
			}
			tc.mutate(&config)
			_, _, _, err := buildOpAttestor(context.Background(), zap.NewNop(), config, nil)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("buildOpAttestor error = %v, want %q", err, tc.want)
			}
		})
	}
}
