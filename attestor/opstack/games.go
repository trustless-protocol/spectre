package opstack

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	contractDisputeGameFactory "attestor/bindings/DisputeGameFactory"
	contractFaultDisputeGame "attestor/bindings/FaultDisputeGame"
)

// FactoryGameSource ingests output-root proposals by enumerating the L1
// DisputeGameFactory. The factory is an append-only indexed list
// (gameCount/gameAtIndex), so an index cursor recovers gaps exactly by
// construction — stronger than log scanning.
type FactoryGameSource struct {
	l1      *ethclient.Client
	factory *contractDisputeGameFactory.Contract
	address common.Address
}

func NewFactoryGameSource(ctx context.Context, l1 *ethclient.Client, factoryAddr common.Address) (*FactoryGameSource, error) {
	code, err := l1.CodeAt(ctx, factoryAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check DisputeGameFactory code at %s: %w", factoryAddr, err)
	}
	if len(code) == 0 {
		return nil, fmt.Errorf("no contract code at DisputeGameFactory address %s", factoryAddr)
	}
	factory, err := contractDisputeGameFactory.NewContract(factoryAddr, l1)
	if err != nil {
		return nil, fmt.Errorf("failed to bind DisputeGameFactory at %s: %w", factoryAddr, err)
	}
	return &FactoryGameSource{l1: l1, factory: factory, address: factoryAddr}, nil
}

func (g *FactoryGameSource) GameCount(ctx context.Context) (uint64, error) {
	count, err := g.factory.GameCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("DisputeGameFactory.gameCount failed: %w", err)
	}
	return count.Uint64(), nil
}

// GameAtIndex fetches one proposal: the factory entry plus the game proxy's
// root claim and L2 block number.
func (g *FactoryGameSource) GameAtIndex(ctx context.Context, index uint64) (ProposedRoot, error) {
	opts := &bind.CallOpts{Context: ctx}
	entry, err := g.factory.GameAtIndex(opts, new(big.Int).SetUint64(index))
	if err != nil {
		return ProposedRoot{}, fmt.Errorf("DisputeGameFactory.gameAtIndex(%d) failed: %w", index, err)
	}
	game, err := contractFaultDisputeGame.NewContract(entry.Proxy, g.l1)
	if err != nil {
		return ProposedRoot{}, fmt.Errorf("failed to bind dispute game %s (index %d): %w", entry.Proxy, index, err)
	}
	rootClaim, err := game.RootClaim(opts)
	if err != nil {
		return ProposedRoot{}, fmt.Errorf("disputeGame(%d).rootClaim failed: %w", index, err)
	}
	l2Block, err := gameL2Height(opts, game, index)
	if err != nil {
		return ProposedRoot{}, err
	}
	return ProposedRoot{
		GameIndex:     index,
		GameAddress:   entry.Proxy,
		GameType:      entry.GameType,
		RootClaim:     rootClaim,
		L2BlockNumber: l2Block.Uint64(),
		L1Timestamp:   entry.Timestamp,
	}, nil
}

// gameL2Height resolves the L2 height a dispute game covers, across the
// FaultDisputeGame interface rename l2BlockNumber() -> l2SequenceNumber()
// ("Super Roots", ethereum-optimism/optimism develop). Measured live
// 2026-07-24: Base Sepolia games (index 22032, type 621 — found in PR #241)
// answer only l2SequenceNumber() and revert on l2BlockNumber(); OP Mainnet's
// respected type-8 games (index 18579) answer both selectors with the same
// value. Try the forward-going selector first and fall back to the legacy
// one for older deployments; games arrive ~1/hour, so an occasional extra
// reverted eth_call is negligible.
//
// For non-interop chains the sequence number is believed to coincide with
// the block number (the live value tracks the chain's block height), but
// that equivalence is not independently confirmed — revisit if an interop
// chain is ever configured as a source.
func gameL2Height(opts *bind.CallOpts, game *contractFaultDisputeGame.Contract, index uint64) (*big.Int, error) {
	seq, seqErr := game.L2SequenceNumber(opts)
	if seqErr == nil {
		return seq, nil
	}
	block, blockErr := game.L2BlockNumber(opts)
	if blockErr == nil {
		return block, nil
	}
	return nil, fmt.Errorf("disputeGame(%d) exposes neither l2SequenceNumber (%v) nor l2BlockNumber: %w", index, seqErr, blockErr)
}

// bootstrapStartIndex finds the first game index whose L1 creation timestamp
// is >= cutoff, by binary search over the factory's monotonic timestamps. With
// count == 0 or every game older than cutoff it returns count (start at the
// live edge).
func bootstrapStartIndex(ctx context.Context, games gameSource, count uint64, cutoff uint64) (uint64, error) {
	lo, hi := uint64(0), count
	for lo < hi {
		mid := lo + (hi-lo)/2
		g, err := games.GameAtIndex(ctx, mid)
		if err != nil {
			return 0, fmt.Errorf("bootstrap scan at index %d: %w", mid, err)
		}
		if g.L1Timestamp >= cutoff {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo, nil
}

// wakeHintLoop subscribes to DisputeGameCreated over WebSocket and pokes the
// wake channel so runOnce fires ahead of the poll ticker. It is a latency hint
// only: it owns no shared state, never advances any cursor, and the poll loop
// is fully correct without it. Reconnects with a fixed delay on any error,
// mirroring the subscriber reconnect shape.
func wakeHintLoop(ctx context.Context, logger *zap.SugaredLogger, wsURL string, factoryAddr common.Address, wake chan<- struct{}) {
	const reconnectDelay = 10 * time.Second
	for {
		if ctx.Err() != nil {
			return
		}
		err := wakeHintOnce(ctx, wsURL, factoryAddr, wake)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			logger.Warnf("[opstack attestor] DisputeGameCreated hint subscription error (poll loop unaffected): %v; reconnecting in %s", err, reconnectDelay)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(reconnectDelay):
		}
	}
}

func wakeHintOnce(ctx context.Context, wsURL string, factoryAddr common.Address, wake chan<- struct{}) error {
	client, err := ethclient.DialContext(ctx, wsURL)
	if err != nil {
		return fmt.Errorf("failed to dial L1 ws %s: %w", wsURL, err)
	}
	defer client.Close()
	factory, err := contractDisputeGameFactory.NewContract(factoryAddr, client)
	if err != nil {
		return fmt.Errorf("failed to bind DisputeGameFactory for hint subscription: %w", err)
	}
	events := make(chan *contractDisputeGameFactory.ContractDisputeGameCreated, 16)
	sub, err := factory.WatchDisputeGameCreated(&bind.WatchOpts{Context: ctx}, events, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe DisputeGameCreated: %w", err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-sub.Err():
			return fmt.Errorf("DisputeGameCreated subscription dropped: %w", err)
		case <-events:
			select {
			case wake <- struct{}{}:
			default: // a wake is already queued; the next runOnce ingests everything
			}
		}
	}
}
