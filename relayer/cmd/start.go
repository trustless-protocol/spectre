// The start command: build the configured relay path and run it.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"relayer/prover"
	"relayer/relay"
	"relayer/services"
	"relayer/transaction"
	utils "relayer/utils"
)

// Start starts the relay service loop.
// It loads config from a JSON file, connects to Cosmos and Ethereum,
// and subscribes to send_packet events to relay packets.
func Start(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start the relay loop",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigPath)
			if err != nil {
				return fmt.Errorf("failed to get config path: %w", err)
			}

			// Load .env for prover paths, private keys, etc.
			_ = godotenv.Load()

			benchmarkFlag, err := cmd.Flags().GetBool(flagBenchmark)
			if err != nil {
				return fmt.Errorf("failed to get benchmark flag: %w", err)
			}
			utils.SetBenchEnabled(benchmarkFlag || utils.BenchEnabled())
			if utils.BenchEnabled() {
				log.Printf("[benchmark] enabled: detailed gas/timing logs are active")
			}

			if err := validateStartupKeys(); err != nil {
				return err
			}

			runCtx, stopSignals := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stopSignals()

			// Load JSON config
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			// Checks that only apply to relaying, so they live here rather than in
			// the shared loadConfig. Run before the prover load: a config error
			// should surface immediately, not after the bucket registry is read.
			if err := validateRelayStartupConfig(cfg); err != nil {
				return err
			}

			// Load the prover once and share it across every source loop — the
			// bucket registry (r1cs/pk/vk) is read-only after load, so concurrent
			// GenerateProof calls are safe. On a GPU backend the calls serialize
			// on the device; correctness is unaffected.
			binDir := envOrDefault("PROVER_BIN_DIR", "./bin")
			selectedBackend, hasBackendOverride, err := proofBackendFromFlags(cmd)
			if err != nil {
				return fmt.Errorf("failed to resolve proof backend: %w", err)
			}

			var p *prover.EcipProver
			if hasBackendOverride {
				logger.Sugar().Infof("start: overriding proof backend via flags: %s", selectedBackend.Name())
				p, err = prover.NewProverWithBackend(binDir, selectedBackend)
			} else {
				p, err = prover.NewProver(binDir)
			}
			if err != nil {
				return fmt.Errorf("failed to load prover: %w", err)
			}

			sources := cfg.CosmosToEthConfigs
			l2Dests := cfg.CosmosToL2Configs
			l2Sources := cfg.L2ToCosmosConfigs
			if len(sources) == 0 && len(l2Dests) == 0 && len(l2Sources) == 0 {
				return fmt.Errorf("no relay source configured in %s (need a cosmos_to_eth, cosmos_to_l2, or l2_to_cosmos module)", configPath)
			}
			// validateL2TimeoutReturnPathConfigs is gone with #328: the return path
			// is now resolved per dest by matching the L2 chain id at build time
			// (l2TimeoutReturnPathConfigForDest), so a separate up-front pass would
			// only duplicate it. This check stays because it answers a different
			// question -- whether each l2_to_cosmos source's configured chain id is
			// the chain its RPC actually serves -- and it must run before anything
			// dials.
			// runCtx, not cmd.Context(): the probe wraps whatever it is given in a
			// 10s timeout, so on cmd.Context() a Ctrl-C during a hung L2 RPC waited
			// out 10s per source before the process could exit.
			if err := validateL2ChainIDs(runCtx, l2Sources); err != nil {
				return err
			}
			// EVM nonces are per (chain id, sender), and the Handler mutex only
			// serializes allocations in this process. Resolve every EVM endpoint's
			// actual chain id and hold those cross-process locks before any builder
			// can submit. A non-answer is fatal here: starting without a chain id
			// would silently mean starting without this guard.
			releaseEVMSignerLocks, err := acquireEVMSignerLocks(runCtx, cfg)
			if err != nil {
				return err
			}
			defer releaseEVMSignerLocks()
			// Env overrides (ICS26_CLIENT_ID, COSMOS_WASM_CLIENT_ID, ROLE_MANAGER)
			// name a single source; only honor them when exactly one is
			// configured, otherwise they would wrongly apply to every source.
			allowEnvOverride := len(sources) == 1

			// One shared TransactionHandler across all sources. What sharing buys is
			// the mutex, not the cache: two sources can target the SAME (chain,
			// signer) pair, and only a single h.mu can serialize their nonce
			// allocation. Per-source handlers would each take their own lock and
			// hand the same nonce to both.
			//
			// Sharing the cache is safe because it is keyed by (chainID, address)
			// — see transaction.evmNonceKey — so sources submitting to DIFFERENT
			// EVM chains keep separate nonces for the same signer. Before that
			// keying the cache was a single counter, and a nonce allocated against
			// one L2 was reused against another as a future nonce: accepted into
			// the queued pool, never minable, and never surfaced as an error (#320).
			//
			// The Cosmos side re-queries the account sequence fresh under cosmosMu,
			// so sharing is safe there too.
			txHandler := &transaction.Handler{}

			// One independent relay loop per Cosmos→ETH source. Each has its own
			// Tendermint RPC, SpectreClient and router client id; they share the
			// prover, the TransactionHandler, the ETH beacon endpoint, and the
			// same ICS26Router (ETH events are partitioned by the per-source
			// router client id filter).
			var wg sync.WaitGroup
			total := len(sources) + len(l2Dests) + len(l2Sources)
			loopErrCh := make(chan error, total)
			cleanups := make([]func(), 0, total)
			relayCtx, cancelRelays := context.WithCancel(runCtx)
			defer cancelRelays()
			l2ReturnPaths := make([]l2TimeoutReturnPathConfig, 0, len(l2Dests))
			onceCleanup := func(cleanup func()) func() {
				var once sync.Once
				return func() {
					if cleanup != nil {
						once.Do(cleanup)
					}
				}
			}
			stopAndCleanup := func() {
				stopRelaysAndCleanup(cancelRelays, wg.Wait, cleanups)
			}
			for i := range sources {
				svc, deps, cleanup, err := buildCosmosToEthSource(
					logger, sources[i], cfg.EthToCosmosConfig, cfg.BatchConfig, p, txHandler,
					buildCosmosToEthSourceOptions{
						allowEnvOverride:   allowEnvOverride,
						startSubscriptions: true,
						// start drives the eth->cosmos direction, so a missing
						// cosmos_wasm_client_id must fail here rather than at the
						// first packet that needs it.
						relaysEVMToCosmos: true,
						pendingStateDir:   pendingStateDir(configPath, "cosmos-to-eth", sources[i].ICS26ClientID),
					},
				)
				if err != nil {
					stopAndCleanup()
					return fmt.Errorf("cosmos_to_eth source %q: %w", sources[i].ICS26ClientID, err)
				}
				cleanup = onceCleanup(cleanup)
				cleanups = append(cleanups, cleanup)
				wg.Add(1)
				go func(svc *services.Services, deps services.RelayDeps, cleanup func()) {
					defer wg.Done()
					defer cleanup()
					// Chain-adapter RelayModule engine — the sole relay engine since
					// the legacy services.StartLoop was removed after the cutover.
					if err := runAdapterEngine(relayCtx, svc, deps); err != nil {
						loopErrCh <- fmt.Errorf("cosmos_to_eth source %q: %w", deps.IDs.CosmosOnEVM, err)
					}
				}(svc, deps, cleanup)
			}

			// One independent relay loop per Cosmos→L2 destination: the same groth16
			// pipeline as Cosmos→ETH pointed at the L2's SpectreClient/ICS26Router,
			// with no reverse beacon direction.
			//
			// Built in two passes: first resolve every dest's return path and match
			// every l2_to_cosmos source, then launch the goroutines below. This way
			// a return-path mismatch aborts before any Cosmos→L2 connection opens,
			// instead of after the relay loop is already dialing/subscribing.
			type builtL2Dest struct {
				svc     *services.Services
				deps    services.RelayDeps
				cleanup func()
			}
			builtL2Dests := make([]builtL2Dest, 0, len(l2Dests))
			for i := range l2Dests {
				svc, deps, cleanup, err := buildCosmosToL2Dest(
					logger, l2Dests[i], cfg.BatchConfig, p, txHandler,
					pendingStateDir(configPath, "cosmos-to-l2", l2Dests[i].ICS26ClientID),
				)
				if err != nil {
					stopAndCleanup()
					return fmt.Errorf("cosmos_to_l2 dest %q: %w", l2Dests[i].ICS26ClientID, err)
				}
				cleanup = onceCleanup(cleanup)
				cleanups = append(cleanups, cleanup)
				name := ""
				if i < len(cfg.CosmosToL2Names) {
					name = cfg.CosmosToL2Names[i]
				}
				// Only resolve the return-path chain id when there's an
				// l2_to_cosmos source to match it against; otherwise a
				// forward-only deployment pays an avoidable L2 RPC round
				// trip (and startup-abort risk) for a path nothing consults.
				if len(l2Sources) > 0 {
					returnPath, err := l2TimeoutReturnPathConfigForDest(runCtx, name, l2Dests[i], l2TimeoutReturnPath{svc: svc, deps: deps})
					if err != nil {
						stopAndCleanup()
						return fmt.Errorf("cosmos_to_l2 dest %q: %w", l2Dests[i].ICS26ClientID, err)
					}
					l2ReturnPaths = append(l2ReturnPaths, returnPath)
				}
				// Accumulate rather than launching here: nothing dials, subscribes
				// or submits until every dest and source has built and matched.
				builtL2Dests = append(builtL2Dests, builtL2Dest{svc: svc, deps: deps, cleanup: cleanup})
			}

			// One independent relay module per L2->Cosmos source (opstack/arbitrum).
			// Each dials its own L1/L2/Cosmos clients + attestor sidecar; they share
			// the TransactionHandler (same Cosmos signer -> shared sequence path).
			type builtL2Source struct {
				module   *relay.Module
				cleanup  func()
				srcChain string
			}
			builtL2Sources := make([]builtL2Source, 0, len(l2Sources))
			for i := range l2Sources {
				timeoutReturn, err := findL2TimeoutReturnPath(runCtx, l2Sources[i], l2ReturnPaths)
				if err != nil {
					stopAndCleanup()
					return fmt.Errorf("l2_to_cosmos source %q: %w", l2Sources[i].AttestorSrcChain, err)
				}
				module, cleanup, err := buildL2ToCosmosModule(logger, l2Sources[i], txHandler, timeoutReturn)
				if err != nil {
					stopAndCleanup()
					return fmt.Errorf("l2_to_cosmos source %q: %w", l2Sources[i].AttestorSrcChain, err)
				}
				cleanup = onceCleanup(cleanup)
				cleanups = append(cleanups, cleanup)
				builtL2Sources = append(builtL2Sources, builtL2Source{module: module, cleanup: cleanup, srcChain: l2Sources[i].AttestorSrcChain})
			}

			// Every Cosmos→L2 dest and L2→Cosmos source above built and matched
			// cleanly — only now do we start dialing/subscribing/submitting.
			for _, d := range builtL2Dests {
				wg.Add(1)
				go func(svc *services.Services, deps services.RelayDeps, cleanup func()) {
					defer wg.Done()
					defer cleanup()
					if err := runCosmosToL2Engine(relayCtx, svc, deps); err != nil {
						loopErrCh <- fmt.Errorf("cosmos_to_l2 dest %q: %w", deps.IDs.CosmosOnEVM, err)
					}
				}(d.svc, d.deps, d.cleanup)
			}
			for _, s := range builtL2Sources {
				wg.Add(1)
				go func(module *relay.Module, cleanup func(), srcChain string) {
					defer wg.Done()
					defer cleanup()
					if err := runL2Engine(relayCtx, module); err != nil {
						loopErrCh <- fmt.Errorf("l2_to_cosmos source %q: %w", srcChain, err)
					}
				}(s.module, s.cleanup, s.srcChain)
			}

			logger.Sugar().Infof("Relayer started: relaying %d Cosmos→ETH + %d Cosmos→L2 + %d L2→Cosmos source(s)", len(sources), len(l2Dests), len(l2Sources))
			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			select {
			case err := <-loopErrCh:
				stopAndCleanup()
				return err
			case <-done:
				// A worker sends its error before its deferred wg.Done. If all
				// workers have returned, prefer that buffered error over treating
				// the coincident done signal as a clean shutdown.
				select {
				case err := <-loopErrCh:
					return err
				default:
					return nil
				}
			case <-runCtx.Done():
				logger.Sugar().Infof("Relayer shutdown requested: %v", runCtx.Err())
				stopAndCleanup()
				logger.Sugar().Info("Relayer clients stopped; exiting")
			}

			return nil
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().Bool(flagGPUProve, false, "use the ICICLE GPU backend for proving (or set GPU_PROVE=1); requires an icicle-enabled build")
	cmd.Flags().Bool(flagBenchmark, false, "enable detailed benchmark gas/timing logs (or set RELAYER_BENCHMARK=1)")
	return cmd
}
