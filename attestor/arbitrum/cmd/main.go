package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"attestor/arbitrum"
	boldattestor "attestor/arbitrum/bold"
	legacyattestor "attestor/arbitrum/legacy"
	attestorserver "attestor/arbitrum/server"
	attestorpb "attestor/types/attestor"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	grpcGracefulStopPeriod         = 10 * time.Second
	nitroSubscriptionRetryInterval = time.Second
	nitroHeadBufferSize            = 64
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := execute(ctx, os.Args[1:]); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}

func runAttestor(ctx context.Context, configPath string) (runErr error) {
	config, err := arbitrum.LoadDaemonConfig(configPath)
	if err != nil {
		return err
	}
	nitroConfig, err := config.NitroProcessConfig()
	if err != nil {
		return fmt.Errorf("load managed Nitro process config: %w", err)
	}
	nitroProcess, err := arbitrum.StartManagedNitro(ctx, nitroConfig, os.Stdout, os.Stderr)
	if err != nil {
		return fmt.Errorf("start managed Nitro process: %w", err)
	}
	defer func() {
		runErr = errors.Join(runErr, nitroProcess.Close())
	}()

	nitroClient := ethclient.NewClient(nitroProcess.RPC())
	runtimeState, err := arbitrum.NewRuntimeState(nitroClient)
	if err != nil {
		return fmt.Errorf("initialize Nitro runtime state: %w", err)
	}
	runtimePollInterval, err := config.RuntimePollDuration()
	if err != nil {
		return fmt.Errorf("load runtime poll interval: %w", err)
	}
	go monitorRuntimeState(ctx, runtimeState, nitroClient, runtimePollInterval)

	l1Client, err := ethclient.DialContext(ctx, config.L1RPCURL)
	if err != nil {
		return fmt.Errorf("connect to Ethereum L1 RPC: %w", err)
	}
	defer l1Client.Close()
	assertionSource, err := newAssertionSource(config, l1Client)
	if err != nil {
		return fmt.Errorf("initialize RollupCore source: %w", err)
	}
	attestedRootStore, err := arbitrum.LoadAttestedRootStore(
		config.AttestorStatePath,
		config.SrcChain,
		config.AssertionStartBlock,
	)
	if err != nil {
		return fmt.Errorf("load attested-root state: %w", err)
	}
	assertionPollInterval, err := config.AssertionPollDuration()
	if err != nil {
		return fmt.Errorf("load assertion poll interval: %w", err)
	}
	assertionLoop, err := boldattestor.NewAssertionAttestor(
		assertionSource,
		runtimeState,
		attestedRootStore,
		boldattestor.AssertionAttestorConfig{
			PollInterval:    assertionPollInterval,
			MaxL1BlockRange: config.AssertionBlockRange(),
		},
	)
	if err != nil {
		return fmt.Errorf("initialize Arbitrum assertion attestor: %w", err)
	}
	if err := assertionLoop.ValidateChainIDs(ctx, config.L1ChainID, config.L2ChainID); err != nil {
		return fmt.Errorf("validate Arbitrum deployment identity: %w", err)
	}
	sourceIdentity := assertionSourceIdentity(config)
	if err := attestedRootStore.BindSourceIdentity(sourceIdentity); err != nil {
		return fmt.Errorf("bind attested-root source identity: %w", err)
	}
	if err := attestedRootStore.Save(); err != nil {
		return fmt.Errorf("persist attested-root source identity: %w", err)
	}
	if err := assertionLoop.SyncOnce(ctx); err != nil {
		return fmt.Errorf("perform initial assertion attestation: %w", err)
	}
	go assertionLoop.Run(ctx)

	grpcService, err := attestorserver.NewAttestorServerWithRuntimeAndFeeds(
		runtimeState,
		map[string]attestorserver.AttestedRootReader{
			config.SrcChain: attestedRootStore,
		},
	)
	if err != nil {
		return fmt.Errorf("initialize attestor gRPC service: %w", err)
	}
	listener, err := net.Listen("tcp", config.GRPCListenAddress)
	if err != nil {
		return fmt.Errorf("listen for attestor gRPC at %s: %w", config.GRPCListenAddress, err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	attestorpb.RegisterAttestorServiceServer(grpcServer, grpcService)
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(attestorpb.AttestorService_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)
	healthv1.RegisterHealthServer(grpcServer, healthServer)

	log.Printf("attestor gRPC listening at %s", listener.Addr())
	if err := serveWithManagedNitro(ctx, grpcServer, listener, nitroProcess); err != nil {
		return fmt.Errorf("serve attestor gRPC: %w", err)
	}
	return nil
}

type rollupL1Client interface {
	boldattestor.L1Client
	legacyattestor.L1Client
}

func newAssertionSource(
	config arbitrum.DaemonConfig,
	client rollupL1Client,
) (boldattestor.AssertionSource, error) {
	address := common.HexToAddress(config.RollupCoreAddress)
	switch config.EffectiveRollupProtocol() {
	case arbitrum.RollupProtocolBoLDV2:
		return boldattestor.NewRollupCoreSource(
			client,
			address,
			common.HexToHash(config.AssertionsMappingSlot),
			config.AssertionStatusOffset,
		)
	case arbitrum.RollupProtocolLegacyNitro:
		return legacyattestor.NewRollupCoreSource(client, address)
	default:
		return nil, fmt.Errorf(
			"unsupported rollup protocol %q",
			config.EffectiveRollupProtocol(),
		)
	}
}

func assertionSourceIdentity(config arbitrum.DaemonConfig) string {
	address := common.HexToAddress(config.RollupCoreAddress).Hex()
	if config.EffectiveRollupProtocol() == arbitrum.RollupProtocolLegacyNitro {
		return fmt.Sprintf(
			"l1:%d/l2:%d/rollup:%s/protocol:%s",
			config.L1ChainID,
			config.L2ChainID,
			address,
			config.EffectiveRollupProtocol(),
		)
	}
	// Preserve the existing BoLD identity format so upgrading the daemon does
	// not invalidate a correctly pinned production state file.
	return fmt.Sprintf(
		"l1:%d/l2:%d/rollup:%s/assertions:%s/status-offset:%d",
		config.L1ChainID,
		config.L2ChainID,
		address,
		common.HexToHash(config.AssertionsMappingSlot).Hex(),
		config.AssertionStatusOffset,
	)
}

type runtimeStateRefresher interface {
	Refresh(context.Context) ([]arbitrum.FinalizedConsistency, error)
}

type nitroHeadSubscriber interface {
	SubscribeNewHead(context.Context, chan<- *types.Header) (ethereum.Subscription, error)
}

func monitorRuntimeState(
	ctx context.Context,
	runtimeState runtimeStateRefresher,
	subscriber nitroHeadSubscriber,
	interval time.Duration,
) {
	monitorRuntimeStateWithRetry(
		ctx,
		runtimeState,
		subscriber,
		interval,
		nitroSubscriptionRetryInterval,
	)
}

func monitorRuntimeStateWithRetry(
	ctx context.Context,
	runtimeState runtimeStateRefresher,
	subscriber nitroHeadSubscriber,
	reconcileInterval time.Duration,
	subscriptionRetryInterval time.Duration,
) {
	headUpdates := make(chan struct{}, 1)
	go streamNitroHeadUpdates(ctx, subscriber, headUpdates, subscriptionRetryInterval)

	refresh := func() {
		checks, err := runtimeState.Refresh(ctx)
		if err != nil {
			if ctx.Err() == nil {
				log.Printf("attestor runtime refresh failed: %v", err)
			}
			return
		}
		for _, check := range checks {
			switch {
			case check.FinalizedReorg:
				log.Printf(
					"attestor finalized head changed or regressed: height=%d state_root=%s previous_height=%d previous_state_root=%s",
					check.Finalized.BlockNumber,
					check.Finalized.StateRoot.Hex(),
					check.PreviousFinalized.BlockNumber,
					check.PreviousFinalized.StateRoot.Hex(),
				)
			case check.Consistent():
				log.Printf(
					"attestor finalized consistency verified: height=%d state_root=%s",
					check.Finalized.BlockNumber,
					check.Finalized.StateRoot.Hex(),
				)
			case !check.UnsafeObserved || !check.SafeObserved:
				log.Printf(
					"attestor finalized consistency incomplete: height=%d unsafe_observed=%t safe_observed=%t",
					check.Finalized.BlockNumber,
					check.UnsafeObserved,
					check.SafeObserved,
				)
			default:
				log.Printf(
					"attestor finalized consistency mismatch: height=%d finalized_root=%s unsafe_root=%s safe_root=%s finalized_reorg=%t",
					check.Finalized.BlockNumber,
					check.Finalized.StateRoot.Hex(),
					check.Unsafe.StateRoot.Hex(),
					check.Safe.StateRoot.Hex(),
					check.FinalizedReorg,
				)
			}
		}
	}

	refresh()
	ticker := time.NewTicker(reconcileInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-headUpdates:
			refresh()
		case <-ticker.C:
			refresh()
		}
	}
}

func streamNitroHeadUpdates(
	ctx context.Context,
	subscriber nitroHeadSubscriber,
	updates chan<- struct{},
	retryInterval time.Duration,
) {
	for {
		headers := make(chan *types.Header, nitroHeadBufferSize)
		subscription, err := subscriber.SubscribeNewHead(ctx, headers)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("subscribe to Nitro new heads failed: %v", err)
			if !waitForRetry(ctx, retryInterval) {
				return
			}
			continue
		}

		// Reconcile immediately after subscribing to close the gap between the
		// previous snapshot and the subscription becoming active.
		signalRuntimeUpdate(updates)
		err = consumeNitroHeadSubscription(ctx, subscription, headers, updates)
		subscription.Unsubscribe()
		if ctx.Err() != nil {
			return
		}
		log.Printf("Nitro new-head subscription ended: %v; reconnecting", err)
		if !waitForRetry(ctx, retryInterval) {
			return
		}
	}
}

func consumeNitroHeadSubscription(
	ctx context.Context,
	subscription ethereum.Subscription,
	headers <-chan *types.Header,
	updates chan<- struct{},
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case header, ok := <-headers:
			if !ok {
				return errors.New("Nitro new-head channel closed")
			}
			if header != nil {
				signalRuntimeUpdate(updates)
			}
		case err, ok := <-subscription.Err():
			if !ok || err == nil {
				return errors.New("Nitro new-head subscription closed")
			}
			return err
		}
	}
}

func signalRuntimeUpdate(updates chan<- struct{}) {
	select {
	case updates <- struct{}{}:
	default:
	}
}

func waitForRetry(ctx context.Context, interval time.Duration) bool {
	timer := time.NewTimer(interval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

type nitroProcessMonitor interface {
	Done() <-chan struct{}
	Err() error
}

func serveWithManagedNitro(
	ctx context.Context,
	server *grpc.Server,
	listener net.Listener,
	nitro nitroProcessMonitor,
) error {
	serveDone := make(chan error, 1)
	go func() {
		serveDone <- server.Serve(listener)
	}()

	select {
	case err := <-serveDone:
		if errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}
		return err
	case <-nitro.Done():
		server.Stop()
		<-serveDone
		if err := nitro.Err(); err != nil {
			return fmt.Errorf("managed Nitro process exited: %w", err)
		}
		return errors.New("managed Nitro process exited unexpectedly")
	case <-ctx.Done():
		gracefulStop(server, grpcGracefulStopPeriod)
		<-serveDone
		return ctx.Err()
	}
}

func gracefulStop(server *grpc.Server, timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		server.Stop()
		<-done
	}
}
