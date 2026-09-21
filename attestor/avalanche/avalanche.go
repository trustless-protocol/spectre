// Package avalanche is the Avalanche C-Chain attestor plugin. It is the
// simplest plugin in the tree by design: Avalanche acceptance is finality —
// there are no reorgs, no dispute games, no derivation to replay — so the
// "replica" is the plugin's own C-Chain RPC endpoint and every accepted block
// is immediately attestable. The plugin polls the finalized (= last accepted)
// head, answers feed reads from it, and signs block identities it re-fetches
// from its own endpoint at verification time.
//
// Coreth headers carry extra fields go-ethereum's types.Header drops, and a
// dropped field would make types.Header.Hash() compute the wrong block hash.
// The plugin therefore never decodes full headers: it reads number, hash and
// stateRoot straight from the RPC JSON — the node already computed the true
// hash — and leaves hash *recomputation* to the wasm client and the relayer's
// coreth reader, which are the parties that must not trust an RPC.
package avalanche

import (
	"context"
	"fmt"
	"sync"
	"time"

	"attestor/core"
	"attestor/types/attestation"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

const (
	defaultPollInterval = 2 * time.Second
	defaultWatchPoll    = 2 * time.Second
)

// Config is one avalanche_source module.
type Config struct {
	SrcChain     string
	RpcURL       string
	ChainID      uint64 // C-Chain EVM chain id (43114 mainnet / 43113 Fuji / 43112 local)
	PollInterval time.Duration
}

// Validate rejects an unusable source before anything is dialed.
func (c Config) Validate() error {
	if c.SrcChain == "" {
		return fmt.Errorf("avalanche_source src_chain must not be empty")
	}
	if c.RpcURL == "" {
		return fmt.Errorf("avalanche_source c_chain_rpc_url is required")
	}
	if c.ChainID == 0 {
		return fmt.Errorf("avalanche_source chain_id must not be zero")
	}
	return nil
}

// blockIdentity is the canonical identity of one accepted C-Chain block as the
// node reports it.
type blockIdentity struct {
	Number    uint64
	Hash      ethcommon.Hash
	StateRoot ethcommon.Hash
}

// CChainAttestor implements every core port for one C-Chain source.
type CChainAttestor struct {
	cfg    Config
	client *gethrpc.Client
	signer attestation.Signer

	mu   sync.RWMutex
	head blockIdentity
	seen bool
}

// New wires a C-Chain attestor over an already-dialed RPC client.
func New(cfg Config, client *gethrpc.Client, signer attestation.Signer) (*CChainAttestor, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("avalanche_source %q: rpc client must not be nil", cfg.SrcChain)
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = defaultPollInterval
	}
	return &CChainAttestor{cfg: cfg, client: client, signer: signer}, nil
}

// Name identifies the runner in host supervision logs.
func (a *CChainAttestor) Name() string { return "avalanche-" + a.cfg.SrcChain }

// Run polls the finalized head until ctx is cancelled. A failed poll keeps the
// previous head — the feed goes stale rather than backwards, and verification
// re-fetches blocks anyway.
func (a *CChainAttestor) Run(ctx context.Context) error {
	ticker := time.NewTicker(a.cfg.PollInterval)
	defer ticker.Stop()
	for {
		if head, err := a.fetchBlock(ctx, "finalized"); err == nil {
			a.mu.Lock()
			a.head = head
			a.seen = true
			a.mu.Unlock()
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// fetchBlock reads the identity of one block from the plugin's own endpoint.
// tagOrNumber is an RPC block tag or a hex-encoded number.
func (a *CChainAttestor) fetchBlock(ctx context.Context, tagOrNumber string) (blockIdentity, error) {
	var block *struct {
		Number    hexutil.Uint64 `json:"number"`
		Hash      ethcommon.Hash `json:"hash"`
		StateRoot ethcommon.Hash `json:"stateRoot"`
	}
	if err := a.client.CallContext(ctx, &block, "eth_getBlockByNumber", tagOrNumber, false); err != nil {
		return blockIdentity{}, err
	}
	if block == nil {
		return blockIdentity{}, fmt.Errorf("block %s not found", tagOrNumber)
	}
	return blockIdentity{Number: uint64(block.Number), Hash: block.Hash, StateRoot: block.StateRoot}, nil
}

// snapshot returns the last polled head.
func (a *CChainAttestor) snapshot() (blockIdentity, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.head, a.seen
}

func toCoreRoot(b blockIdentity) core.AttestedRoot {
	return core.AttestedRoot{
		L2BlockNumber: b.Number,
		Root:          append([]byte(nil), b.StateRoot.Bytes()...),
		Source:        "avalanche-accepted",
		Provisional:   false,
	}
}

// AttestedUpTo returns the accepted head. Nothing on Avalanche is provisional,
// so the policy never changes the answer.
func (a *CChainAttestor) AttestedUpTo(_ context.Context, _ core.AttestationPolicy) (core.AttestedRoot, bool, error) {
	head, seen := a.snapshot()
	if !seen {
		return core.AttestedRoot{}, false, nil
	}
	return toCoreRoot(head), true, nil
}

// AttestedRootAtOrBelow returns the identity at min(height, head). Every
// accepted block is final, so a below-head read is a plain block fetch.
func (a *CChainAttestor) AttestedRootAtOrBelow(ctx context.Context, height uint64, _ core.AttestationPolicy) (core.AttestedRoot, bool, error) {
	head, seen := a.snapshot()
	if !seen {
		return core.AttestedRoot{}, false, nil
	}
	if height >= head.Number {
		return toCoreRoot(head), true, nil
	}
	block, err := a.fetchBlock(ctx, hexutil.EncodeUint64(height))
	if err != nil {
		return core.AttestedRoot{}, false, core.NewError(core.ErrorUnavailable, fmt.Sprintf("query C-Chain block %d", height), err)
	}
	return toCoreRoot(block), true, nil
}

// Status reports the polled head at every finality level — they are the same
// number on Avalanche.
func (a *CChainAttestor) Status(_ context.Context) (core.FeedStatus, error) {
	status := core.FeedStatus{
		SrcChain:        a.cfg.SrcChain,
		AttestationHead: core.RunModeFinalized,
	}
	if head, seen := a.snapshot(); seen {
		status.Ready = true
		status.ReplicaSeen = true
		status.ReplicaUnsafe = head.Number
		status.ReplicaSafe = head.Number
		status.ReplicaFinalized = head.Number
	}
	return status, nil
}

// VerifyStateRoot binds a relayer candidate to this plugin's own C-Chain view
// and signs the canonical identity only when it matches.
func (a *CChainAttestor) VerifyStateRoot(ctx context.Context, request core.BlockIdentityRequest) (core.SignedBlockIdentityVerdict, error) {
	if !request.RunMode.Valid() {
		return core.SignedBlockIdentityVerdict{}, core.NewError(core.ErrorInvalidArgument, "run_mode must be unsafe, safe or finalized", nil)
	}
	if request.RunMode != core.RunModeFinalized {
		return core.SignedBlockIdentityVerdict{}, core.NewError(
			core.ErrorFailedPrecondition,
			fmt.Sprintf("requested run_mode %s does not match the avalanche attestation head finalized", request.RunMode),
			nil,
		)
	}
	head, seen := a.snapshot()
	if !seen {
		return core.SignedBlockIdentityVerdict{}, core.NewError(core.ErrorUnavailable, "C-Chain head not seen yet", nil)
	}
	if request.BlockNumber > head.Number {
		return core.SignedBlockIdentityVerdict{}, core.NewError(
			core.ErrorFailedPrecondition,
			fmt.Sprintf("block %d is above the accepted head %d", request.BlockNumber, head.Number),
			nil,
		)
	}
	block, err := a.fetchBlock(ctx, hexutil.EncodeUint64(request.BlockNumber))
	if err != nil {
		return core.SignedBlockIdentityVerdict{}, core.NewError(core.ErrorUnavailable, fmt.Sprintf("query C-Chain block %d", request.BlockNumber), err)
	}
	verdict := core.SignedBlockIdentityVerdict{
		BlockNumber: block.Number,
		BlockHash:   [32]byte(block.Hash),
		StateRoot:   [32]byte(block.StateRoot),
	}
	verdict.Valid = block.StateRoot == ethcommon.Hash(request.ExpectedStateRoot)
	if request.ExpectedBlockHash != nil {
		verdict.Valid = verdict.Valid && block.Hash == ethcommon.Hash(*request.ExpectedBlockHash)
	}
	if !verdict.Valid {
		return verdict, nil
	}
	if !a.signer.Configured() {
		return core.SignedBlockIdentityVerdict{}, core.NewError(core.ErrorUnimplemented, "attestation signer is not configured", nil)
	}
	signature, err := a.signer.Sign(request.L2Router[:], request.AttestorSetHash[:], block.Number, block.Hash[:], block.StateRoot[:])
	if err != nil {
		return core.SignedBlockIdentityVerdict{}, core.NewError(core.ErrorInternal, fmt.Sprintf("sign canonical C-Chain block %d", block.Number), err)
	}
	verdict.Signature = signature
	return verdict, nil
}

// WatchFrontier emits the accepted head whenever it changes, mirroring the OP
// adapter's polling stream semantics.
func (a *CChainAttestor) WatchFrontier(ctx context.Context, request core.WatchRequest) (<-chan core.FrontierUpdate, error) {
	updates := make(chan core.FrontierUpdate, 1)
	go func() {
		defer close(updates)
		var previous core.AttestedRoot
		sent := false
		emit := func() bool {
			current, found, err := a.AttestedUpTo(ctx, request.Policy)
			if err != nil || !found {
				sent = false
				return true
			}
			if sent && previous.L2BlockNumber == current.L2BlockNumber && string(previous.Root) == string(current.Root) {
				return true
			}
			select {
			case updates <- core.FrontierUpdate{Root: current.Clone(), Found: true}:
				previous = current
				sent = true
				return true
			case <-ctx.Done():
				return false
			}
		}
		if !emit() {
			return
		}
		ticker := time.NewTicker(defaultWatchPoll)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if !emit() {
					return
				}
			}
		}
	}()
	return updates, nil
}
