package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"attestor/arbitrum"
	boldattestor "attestor/arbitrum/bold"
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

func runAttestor(ctx context.Context, configPath string) error {
	config, err := arbitrum.LoadDaemonConfig(configPath)
	if err != nil {
		return err
	}

	nitroClient, err := ethclient.DialContext(ctx, config.NitroRPCURL)
	if err != nil {
		return fmt.Errorf("connect to Nitro RPC: %w", err)
	}
	defer nitroClient.Close()
	if err := validateNitroChainID(ctx, nitroClient, config.L2ChainID, "RPC"); err != nil {
		return err
	}

	nitroWSClient, err := ethclient.DialContext(ctx, config.NitroWSURL)
	if err != nil {
		return fmt.Errorf("connect to Nitro WebSocket: %w", err)
	}
	defer nitroWSClient.Close()
	if err := validateNitroChainID(ctx, nitroWSClient, config.L2ChainID, "WebSocket"); err != nil {
		return err
	}

	runtimeState, err := arbitrum.NewRuntimeState(nitroClient)
	if err != nil {
		return fmt.Errorf("initialize Nitro runtime state: %w", err)
	}
	runtimePollInterval, err := config.RuntimePollDuration()
	if err != nil {
		return fmt.Errorf("load runtime poll interval: %w", err)
	}
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
	attestationHead, err := config.NormalizedAttestationHead()
	if err != nil {
		return fmt.Errorf("load attestation head: %w", err)
	}
	derivedAttestor, err := arbitrum.NewDerivedRootAttestor(
		runtimeState,
		attestedRootStore,
		arbitrum.DerivedAttestorConfig{
			AttestationHead: attestationHead,
			Disabled:        config.DisableDerivedRoots,
			GapBlocks:       config.DerivedAttestationGap(),
			MaxRoots:        config.DerivedRootLimit(),
		},
	)
	if err != nil {
		return fmt.Errorf("initialize Arbitrum derived-root attestor: %w", err)
	}
	go monitorRuntimeState(
		ctx,
		runtimeAndDerivedRefresher{runtime: runtimeState, derived: derivedAttestor},
		nitroWSClient,
		runtimePollInterval,
	)
	// Assertion backfill can be large when using a rate-limited public L1 RPC.
	// Start it in the retrying background loop so a transient log-query failure
	// does not prevent the gRPC service from becoming available.
	go assertionLoop.Run(ctx)

	grpcService, err := attestorserver.NewAttestorServerWithRuntimeFeedsAndHeads(
		runtimeState,
		map[string]attestorserver.AttestedRootReader{
			config.SrcChain: attestedRootStore,
		},
		map[string]arbitrum.RunMode{
			config.SrcChain: attestationHead,
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
	if err := serveAttestor(ctx, grpcServer, listener); err != nil {
		return fmt.Errorf("serve attestor gRPC: %w", err)
	}
	return nil
}

type chainIDReader interface {
	ChainID(context.Context) (*big.Int, error)
}

func validateNitroChainID(
	ctx context.Context,
	client chainIDReader,
	expected uint64,
	transport string,
) error {
	observed, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("query Nitro %s chain ID: %w", transport, err)
	}
	if observed == nil || !observed.IsUint64() || observed.Uint64() != expected {
		return fmt.Errorf(
			"Nitro %s chain ID is %v, expected %d",
			transport,
			observed,
			expected,
		)
	}
	return nil
}

func newAssertionSource(
	config arbitrum.DaemonConfig,
	client boldattestor.L1Client,
) (boldattestor.AssertionSource, error) {
	address := common.HexToAddress(config.RollupCoreAddress)
	return boldattestor.NewRollupCoreSource(
		client,
		address,
		common.HexToHash(config.AssertionsMappingSlot),
		config.AssertionStatusOffset,
	)
}

func assertionSourceIdentity(config arbitrum.DaemonConfig) string {
	address := common.HexToAddress(config.RollupCoreAddress).Hex()
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

type derivedRootSyncer interface {
	SyncOnce(context.Context) error
}

type runtimeAndDerivedRefresher struct {
	runtime runtimeStateRefresher
	derived derivedRootSyncer
}

func (r runtimeAndDerivedRefresher) Refresh(
	ctx context.Context,
) ([]arbitrum.FinalizedConsistency, error) {
	checks, err := r.runtime.Refresh(ctx)
	if err != nil {
		return nil, err
	}
	if err := r.derived.SyncOnce(ctx); err != nil {
		return nil, err
	}
	return checks, nil
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

func serveAttestor(
	ctx context.Context,
	server *grpc.Server,
	listener net.Listener,
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
