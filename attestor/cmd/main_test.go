package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
