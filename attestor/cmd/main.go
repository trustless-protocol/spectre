// The attestor sidecar binary: runs the OP Stack source attestors
// (derivation-replay verification of output roots) and serves the attested
// root feed over gRPC for the relayer to consume. Standalone — no packet
// relaying, no on-chain writes. Its own binary for the same reason the
// prover setup tool has one: independent lifecycle from the relay loop.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"attestor/opstack"
	attestorserver "attestor/server"
)

const flagConfigPath = "config"

// opSourceConfig is one op_source module entry from the attestor's own JSON
// config file (see config.example.json; the attestor no longer shares the
// relayer's config).
type opSourceConfig struct {
	SrcChain                string `json:"-"` // taken from the module's src_chain
	L1RpcUrl                string `json:"l1_rpc_url"`
	L1WsUrl                 string `json:"l1_ws_url"`
	OpNodeRpcUrl            string `json:"op_node_rpc_url"`
	DisputeGameFactory      string `json:"dispute_game_factory"`
	RespectedGameType       uint32 `json:"respected_game_type"`
	AttestationHead         string `json:"attestation_head"`
	PollIntervalSeconds     uint64 `json:"poll_interval_seconds"`
	StatePath               string `json:"state_path"`
	L1BootstrapLookbackBlks uint64 `json:"l1_bootstrap_lookback_blocks"`
	DisableDerivedRoots     bool   `json:"disable_derived_roots"`
	DerivedGapBlocks        uint64 `json:"derived_attestation_gap_blocks"`
	MaxDerivedRoots         uint64 `json:"max_derived_roots"`
}

type serverConfig struct {
	LogLevel string `json:"log_level"`
	Address  string `json:"address"`
	Port     uint64 `json:"port"`
	// GrpcPort serves the attestor sidecar gRPC API; 0 disables it.
	GrpcPort uint64 `json:"grpc_port"`
}

type attestorConfig struct {
	Server  serverConfig
	Sources []opSourceConfig
}

func validateHexAddress(addr, fieldName string) error {
	if addr == "" {
		return fmt.Errorf("%s is empty", fieldName)
	}
	if !common.IsHexAddress(addr) {
		return fmt.Errorf("invalid hex address %q for %s", addr, fieldName)
	}
	return nil
}

// loadConfig parses the op_source modules and server block out of the attestor
// config file, ignoring every other module type.
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
	// Two attestors must not share a src_chain (metrics label) or a state
	// file (they would clobber each other's cursors).
	seenChains := make(map[string]struct{})
	seenStatePaths := make(map[string]struct{})
	for _, m := range jc.Modules {
		if m.Name != "op_source" {
			continue
		}
		var op opSourceConfig
		if err := json.Unmarshal(m.Config, &op); err != nil {
			return nil, fmt.Errorf("failed to parse op_source config: %w", err)
		}
		op.SrcChain = m.SrcChain
		if op.SrcChain == "" {
			return nil, fmt.Errorf("op_source module requires a non-empty src_chain")
		}
		if _, dup := seenChains[op.SrcChain]; dup {
			return nil, fmt.Errorf("duplicate op_source src_chain %q; each source needs a distinct name", op.SrcChain)
		}
		seenChains[op.SrcChain] = struct{}{}
		if op.StatePath != "" {
			if _, dup := seenStatePaths[op.StatePath]; dup {
				return nil, fmt.Errorf("duplicate op_source state_path %q; each source needs its own state file", op.StatePath)
			}
			seenStatePaths[op.StatePath] = struct{}{}
		}
		if op.L1WsUrl != "" {
			if _, err := url.Parse(op.L1WsUrl); err != nil {
				return nil, fmt.Errorf("op_source.l1_ws_url is not a valid URL: %w", err)
			}
		}
		cfg.Sources = append(cfg.Sources, op)
	}
	return cfg, nil
}

// buildOpAttestor converts one op_source config entry into a running-ready
// attestor: validate, dial L1 + op-node, load the state file.
func buildOpAttestor(ctx context.Context, logger *zap.Logger, op opSourceConfig, metrics *opstack.Metrics) (*opstack.OpStackAttestor, func(), error) {
	if err := validateHexAddress(op.DisputeGameFactory, "op_source.dispute_game_factory"); err != nil {
		return nil, nil, err
	}
	cfg := opstack.Config{
		SrcChain:                op.SrcChain,
		L1RpcUrl:                op.L1RpcUrl,
		L1WsUrl:                 op.L1WsUrl,
		OpNodeRpcUrl:            op.OpNodeRpcUrl,
		DisputeGameFactory:      common.HexToAddress(op.DisputeGameFactory),
		RespectedGameType:       op.RespectedGameType,
		AttestationHead:         opstack.Head(op.AttestationHead),
		PollInterval:            time.Duration(op.PollIntervalSeconds) * time.Second,
		StatePath:               op.StatePath,
		BootstrapLookbackBlocks: op.L1BootstrapLookbackBlks,
		DisableDerivedRoots:     op.DisableDerivedRoots,
		DerivedGapBlocks:        op.DerivedGapBlocks,
		MaxDerivedRoots:         op.MaxDerivedRoots,
	}
	if err := cfg.Validate(); err != nil {
		return nil, nil, err
	}

	l1Client, err := ethclient.DialContext(ctx, cfg.L1RpcUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial L1 rpc %s: %w", cfg.L1RpcUrl, err)
	}
	games, err := opstack.NewFactoryGameSource(ctx, l1Client, cfg.DisputeGameFactory)
	if err != nil {
		l1Client.Close()
		return nil, nil, err
	}
	replica, err := opstack.DialReplica(ctx, cfg.OpNodeRpcUrl)
	if err != nil {
		l1Client.Close()
		return nil, nil, err
	}
	store, err := opstack.LoadStore(cfg.StatePath)
	if err != nil {
		l1Client.Close()
		replica.Close()
		return nil, nil, err
	}
	cleanup := func() {
		l1Client.Close()
		replica.Close()
	}
	hook := &opstack.LogChallengeHook{Logger: logger.Sugar()}
	return opstack.New(cfg, games, replica, store, hook, metrics, logger.Sugar()), cleanup, nil
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
		return fmt.Errorf("no op_source module configured in %s", configPath)
	}

	runCtx, stopSignals := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector())
	metrics := opstack.NewMetrics(registry)

	attestors := make([]*opstack.OpStackAttestor, 0, len(cfg.Sources))
	byChain := make(map[string]*opstack.OpStackAttestor, len(cfg.Sources))
	cleanups := make([]func(), 0, len(cfg.Sources))
	defer func() {
		for _, cleanup := range cleanups {
			cleanup()
		}
	}()
	for _, op := range cfg.Sources {
		a, cleanup, err := buildOpAttestor(runCtx, logger, op, metrics)
		if err != nil {
			return fmt.Errorf("op_source %q: %w", op.SrcChain, err)
		}
		cleanups = append(cleanups, cleanup)
		attestors = append(attestors, a)
		byChain[op.SrcChain] = a
	}

	// Sidecar gRPC API: the relayer consumes the attested-root feed over this
	// endpoint instead of in-process.
	if cfg.Server.GrpcPort != 0 {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.GrpcPort)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen on attestor grpc address %s: %w", addr, err)
		}
		grpcSrv := grpc.NewServer()
		attestorserver.New(byChain).Register(grpcSrv)
		// Reflection lets grpcurl & co discover the service without the proto
		// files — read-only metadata, safe on a local sidecar endpoint.
		reflection.Register(grpcSrv)
		go func() {
			logger.Sugar().Infof("attestor: sidecar gRPC API listening on %s", addr)
			if err := grpcSrv.Serve(lis); err != nil {
				logger.Sugar().Errorf("attestor: grpc server stopped: %v", err)
			}
		}()
		defer grpcSrv.GracefulStop()
	}

	// Prometheus endpoint from the server config block.
	if cfg.Server.Port != 0 {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		metricsSrv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
		go func() {
			logger.Sugar().Infof("attestor: metrics listening on %s/metrics", addr)
			if err := metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Sugar().Errorf("attestor: metrics server failed: %v", err)
			}
		}()
		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = metricsSrv.Shutdown(shutdownCtx)
		}()
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(attestors))
	for _, a := range attestors {
		wg.Add(1)
		go func(a *opstack.OpStackAttestor) {
			defer wg.Done()
			if err := a.Run(runCtx); err != nil {
				errCh <- fmt.Errorf("%s: %w", a.Name(), err)
			}
		}(a)
	}
	logger.Sugar().Infof("attestor: running %d OP Stack attestor(s)", len(attestors))
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case err := <-errCh:
		stopSignals()
		wg.Wait()
		return err
	case <-done:
		return nil
	case <-runCtx.Done():
		logger.Sugar().Infof("attestor: shutdown requested: %v", runCtx.Err())
		wg.Wait()
		return nil
	}
}

func main() {
	zLogger, _ := zap.NewProduction(zap.AddStacktrace(zap.DPanicLevel))
	defer zLogger.Sync()

	rootCmd := &cobra.Command{
		Use:   "attestor",
		Short: "attestor sidecar — replay-verify OP Stack output roots and serve the attested-root feed",
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
