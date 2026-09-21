// The Avalanche attestor sidecar binary: follows one or more C-Chain sources
// through their own RPC endpoints and serves the attested feed + block
// identity signing over gRPC for the relayer to consume. Standalone — no
// packet relaying, no on-chain writes.
//
// Unlike the OP/Arbitrum daemons there is no settlement machinery and no state
// file: Avalanche acceptance is finality, so the entire replica state is the
// last accepted head, re-polled at startup.
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"attestor/avalanche"
	"attestor/core"
	"attestor/host"
	"attestor/types/attestation"
)

const flagConfigPath = "config"

// avalancheSourceConfig is one avalanche_source module entry.
type avalancheSourceConfig struct {
	SrcChain              string `json:"-"` // taken from the module's src_chain
	CChainRpcURL          string `json:"c_chain_rpc_url"`
	ChainID               uint64 `json:"chain_id"`
	PollIntervalSeconds   uint64 `json:"poll_interval_seconds"`
	AttestationSigningKey string `json:"attestation_signing_key"`
}

type serverConfig struct {
	LogLevel string `json:"log_level"`
	Address  string `json:"address"`
	// GrpcPort serves the attestor sidecar gRPC API; 0 disables it.
	GrpcPort uint64 `json:"grpc_port"`
}

type attestorConfig struct {
	Server  serverConfig
	Sources []avalancheSourceConfig
}

// loadConfig parses the avalanche_source modules and server block, ignoring
// every other module type.
func loadConfig(configPath string) (*attestorConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var jc struct {
		Server  serverConfig `json:"server"`
		Modules []struct {
			Name     string          `json:"name"`
			SrcChain string          `json:"src_chain"`
			Config   json.RawMessage `json:"config"`
		} `json:"modules"`
	}
	if err := json.Unmarshal(data, &jc); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	cfg := &attestorConfig{Server: jc.Server}
	seenChains := make(map[string]struct{})
	for _, m := range jc.Modules {
		if m.Name != "avalanche_source" {
			continue
		}
		var src avalancheSourceConfig
		if err := json.Unmarshal(m.Config, &src); err != nil {
			return nil, fmt.Errorf("failed to parse avalanche_source config: %w", err)
		}
		src.SrcChain = m.SrcChain
		if src.SrcChain == "" {
			return nil, fmt.Errorf("avalanche_source module requires a non-empty src_chain")
		}
		if _, dup := seenChains[src.SrcChain]; dup {
			return nil, fmt.Errorf("duplicate avalanche_source src_chain %q; each source needs a distinct name", src.SrcChain)
		}
		seenChains[src.SrcChain] = struct{}{}
		cfg.Sources = append(cfg.Sources, src)
	}
	return cfg, nil
}

// buildAttestor converts one avalanche_source entry into a running-ready
// plugin: build the signer, dial the C-Chain, verify its chain id.
func buildAttestor(ctx context.Context, src avalancheSourceConfig) (*avalanche.CChainAttestor, attestation.Signer, func(), error) {
	signer, err := attestation.NewSigner(src.ChainID, src.AttestationSigningKey)
	if err != nil {
		return nil, attestation.Signer{}, nil, fmt.Errorf("avalanche_source attestation signer: %w", err)
	}
	client, err := gethrpc.DialContext(ctx, src.CChainRpcURL)
	if err != nil {
		return nil, attestation.Signer{}, nil, fmt.Errorf("failed to dial C-Chain rpc %s: %w", src.CChainRpcURL, err)
	}
	var chainID string
	if err := client.CallContext(ctx, &chainID, "eth_chainId"); err != nil {
		client.Close()
		return nil, attestation.Signer{}, nil, fmt.Errorf("query C-Chain chain id: %w", err)
	}
	if chainID != fmt.Sprintf("0x%x", src.ChainID) {
		client.Close()
		return nil, attestation.Signer{}, nil, fmt.Errorf("C-Chain RPC chain id is %s, expected 0x%x", chainID, src.ChainID)
	}
	plugin, err := avalanche.New(avalanche.Config{
		SrcChain:     src.SrcChain,
		RpcURL:       src.CChainRpcURL,
		ChainID:      src.ChainID,
		PollInterval: time.Duration(src.PollIntervalSeconds) * time.Second,
	}, client, signer)
	if err != nil {
		client.Close()
		return nil, attestation.Signer{}, nil, err
	}
	return plugin, signer, client.Close, nil
}

func run(logger *zap.Logger, cmd *cobra.Command) error {
	configPath, err := cmd.Flags().GetString(flagConfigPath)
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}
	cfg, err := loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if len(cfg.Sources) == 0 {
		return fmt.Errorf("no avalanche_source module configured in %s", configPath)
	}

	runCtx, stopSignals := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()

	routes := make(map[string]core.Ports, len(cfg.Sources))
	runners := make([]host.RunnerSpec, 0, len(cfg.Sources))
	cleanups := make([]func(), 0, len(cfg.Sources))
	defer func() {
		for _, cleanup := range cleanups {
			cleanup()
		}
	}()
	for _, src := range cfg.Sources {
		plugin, signer, cleanup, err := buildAttestor(runCtx, src)
		if err != nil {
			return fmt.Errorf("avalanche_source %q: %w", src.SrcChain, err)
		}
		cleanups = append(cleanups, cleanup)
		routes[src.SrcChain] = core.Ports{Feed: plugin, Verifier: plugin, Status: plugin, Watcher: plugin}
		runners = append(runners, host.RunnerSpec{Name: plugin.Name(), Runner: plugin})
		publicKey := signer.PublicKey()
		logger.Sugar().Infof("avalanche_source %s attestation public key: 0x%x", src.SrcChain, publicKey)
		logger.Sugar().Infof(
			"avalanche_source %s attestation public key (base64 for relayer attestors.public_keys): %s",
			src.SrcChain,
			base64.StdEncoding.EncodeToString(publicKey),
		)
	}

	var grpcListener net.Listener
	if cfg.Server.GrpcPort != 0 {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.GrpcPort)
		grpcListener, err = net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen on attestor grpc address %s: %w", addr, err)
		}
	}

	attestorHost, err := host.New(host.Config{
		Routes:           routes,
		Runners:          runners,
		Listener:         grpcListener,
		EnableReflection: true,
	})
	if err != nil {
		return fmt.Errorf("initialize attestor host: %w", err)
	}
	if grpcListener != nil {
		logger.Sugar().Infof("attestor: sidecar gRPC API listening on %s", grpcListener.Addr())
	}
	logger.Sugar().Infof("attestor: running %d Avalanche C-Chain attestor(s)", len(runners))
	if err := attestorHost.Run(runCtx); err != nil {
		return fmt.Errorf("run attestor host: %w", err)
	}
	return nil
}

func main() {
	zLogger, _ := zap.NewProduction(zap.AddStacktrace(zap.DPanicLevel))
	defer zLogger.Sync()

	rootCmd := &cobra.Command{
		Use:   "attestor",
		Short: "attestor sidecar — follow the Avalanche C-Chain and serve the attested-head feed",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(zLogger, cmd)
		},
	}
	rootCmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")

	if err := rootCmd.Execute(); err != nil {
		zLogger.Sugar().Fatal(err)
	}
}
