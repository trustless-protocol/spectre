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

	"attestor/arbitrum"
	"attestor/arbitrum/adapter"
	boldattestor "attestor/arbitrum/bold"
	"attestor/core"
	"attestor/host"
	"attestor/types/attestation"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
	signer, err := attestation.NewSigner(config.L2ChainID, config.AttestationSigningKey)
	if err != nil {
		return fmt.Errorf("load attestation signing key: %w", err)
	}
	log.Printf("Arbitrum attestation public key: 0x%x", signer.PublicKey())

	runtimeState, err := arbitrum.NewRuntimeStateWithConfig(nitroClient, arbitrum.RuntimeStateConfig{
		BackfillMaxBlocks:   config.BackfillMaxBlocks(),
		BackfillConcurrency: int(config.BackfillConcurrency()),
	})
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
	bridge, err := adapter.New(adapter.Options{
		SrcChain:        config.SrcChain,
		Runtime:         runtimeState,
		Feed:            attestedRootStore,
		AttestationHead: attestationHead,
		Signer:          signer,
	})
	if err != nil {
		return fmt.Errorf("compose Arbitrum attestor bridge: %w", err)
	}
	listener, err := net.Listen("tcp", config.GRPCListenAddress)
	if err != nil {
		return fmt.Errorf("listen for attestor gRPC at %s: %w", config.GRPCListenAddress, err)
	}
	defer listener.Close()

	attestorHost, err := host.New(host.Config{
		Routes: map[string]core.Ports{
			config.SrcChain: {Feed: bridge, Verifier: bridge, Status: bridge},
		},
		Runners: []host.RunnerSpec{
			{
				Name: "nitro-runtime-monitor",
				Runner: host.RunnerFunc(func(runCtx context.Context) error {
					return arbitrum.MonitorRuntimeState(
						runCtx,
						runtimeAndDerivedRefresher{runtime: runtimeState, derived: derivedAttestor},
						nitroWSClient,
						arbitrum.RuntimeMonitorConfig{ReconcileInterval: runtimePollInterval},
					)
				}),
			},
			{
				Name: "arbitrum-assertion-attestor",
				Runner: host.RunnerFunc(func(runCtx context.Context) error {
					// Assertion backfill can be large; its loop retries transient L1
					// failures without taking the sidecar down.
					assertionLoop.Run(runCtx)
					return nil
				}),
			},
		},
		Listener: listener,
	})
	if err != nil {
		return fmt.Errorf("initialize attestor host: %w", err)
	}
	log.Printf("attestor gRPC listening at %s", listener.Addr())
	if err := attestorHost.Run(ctx); err != nil {
		return fmt.Errorf("run attestor host: %w", err)
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
