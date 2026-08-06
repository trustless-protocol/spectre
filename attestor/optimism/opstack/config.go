package opstack

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Head selects which replica head gates attestation verdicts. Any verdict
// made on a block the finalized head does not yet cover is provisional and is
// automatically re-checked once finalized covers it; the finalized-head result
// is always authoritative.
type Head string

const (
	// HeadFinalized (default) gates on the replica's finalized head — blocks
	// derived from finalized L1 data ("tx finality", ~20-40 min). Irreversible;
	// verdicts are confirmed immediately.
	HeadFinalized Head = "finalized"
	// HeadSafe gates on the safe head — derived from L1 batch data whose L1
	// block may still reorg before finality. Slightly faster than finalized;
	// verdicts are provisional until the finalized recheck.
	HeadSafe Head = "safe"
	// HeadUnsafe gates on the unsafe head — sequencer gossip. Fastest;
	// verdicts carry sequencer trust until the finalized recheck.
	HeadUnsafe Head = "unsafe"
)

const (
	defaultPollInterval      = 30 * time.Second
	defaultBootstrapLookback = 7200 // L1 blocks (~24h at 12s)
	// defaultDerivedGapBlocks is the minimum L2-block gap between derived
	// attestations — ~5 min at 2s blocks. It is also the floor on how long a
	// packet waits to be relayable, since RelayableHeight follows this frontier.
	defaultDerivedGapBlocks = 150
	defaultMaxDerivedRoots  = 1000 // confirmed derived roots kept in the feed
)

// Config configures one OP Stack source attestor.
type Config struct {
	// SrcChain names the source chain, e.g. "op-mainnet". Used as the metrics
	// label and log prefix.
	SrcChain string
	// L1RpcUrl is the Ethereum L1 HTTP RPC endpoint.
	L1RpcUrl string
	// L1WsUrl optionally enables the DisputeGameCreated wake hint. Empty means
	// poll-only.
	L1WsUrl string
	// OpNodeRpcUrl is the RPC endpoint of the op-node driving the verify-mode
	// replica (optimism_syncStatus / optimism_outputAtBlock).
	OpNodeRpcUrl string
	// DisputeGameFactory is the L1 DisputeGameFactory address.
	DisputeGameFactory common.Address
	// RespectedGameType filters ingested games; games of any other type are
	// skipped (counted, never attested).
	RespectedGameType uint32
	// AttestationHead selects the gating head (safe by default).
	AttestationHead Head
	// PollInterval is the runOnce cadence.
	PollInterval time.Duration
	// StatePath is the JSON state file holding the ingest cursor, pending and
	// recheck sets, and the attested-root feed.
	StatePath string
	// BootstrapLookbackBlocks bounds the first-run ingest scan: with no state
	// file, ingestion starts at the first game created inside the last N L1
	// blocks instead of replaying all of history.
	BootstrapLookbackBlocks uint64
	// DisableDerivedRoots turns off self-derived attestations, leaving only
	// the proposal-verification watchdog. The feed then depends on the
	// proposer's cadence again.
	DisableDerivedRoots bool
	// DerivedGapBlocks is the minimum L2-block gap between consecutive
	// self-derived attestations (bounds feed growth).
	DerivedGapBlocks uint64
	// MaxDerivedRoots caps confirmed self-derived entries kept in the feed;
	// the oldest are pruned beyond it. Game-verified entries are never pruned.
	MaxDerivedRoots uint64
}

func (c *Config) applyDefaults() {
	if c.AttestationHead == "" {
		c.AttestationHead = HeadFinalized
	}
	if c.PollInterval == 0 {
		c.PollInterval = defaultPollInterval
	}
	if c.BootstrapLookbackBlocks == 0 {
		c.BootstrapLookbackBlocks = defaultBootstrapLookback
	}
	if c.DerivedGapBlocks == 0 {
		c.DerivedGapBlocks = defaultDerivedGapBlocks
	}
	if c.MaxDerivedRoots == 0 {
		c.MaxDerivedRoots = defaultMaxDerivedRoots
	}
}

func (c *Config) Validate() error {
	c.applyDefaults()
	if c.SrcChain == "" {
		return fmt.Errorf("op_source.src_chain is required")
	}
	for _, u := range []struct{ val, name string }{
		{c.L1RpcUrl, "op_source.l1_rpc_url"},
		{c.OpNodeRpcUrl, "op_source.op_node_rpc_url"},
	} {
		if u.val == "" {
			return fmt.Errorf("%s is required", u.name)
		}
		if _, err := url.Parse(u.val); err != nil {
			return fmt.Errorf("%s is not a valid URL: %w", u.name, err)
		}
	}
	if c.L1WsUrl != "" {
		if _, err := url.Parse(c.L1WsUrl); err != nil {
			return fmt.Errorf("op_source.l1_ws_url is not a valid URL: %w", err)
		}
	}
	if c.DisputeGameFactory == (common.Address{}) {
		return fmt.Errorf("op_source.dispute_game_factory is required")
	}
	switch c.AttestationHead {
	case HeadFinalized:
	case HeadSafe, HeadUnsafe:
		// Verdicts below the finalized head are provisional (L1-reorg risk for
		// safe, sequencer trust for unsafe) until the finalized recheck
		// confirms them; consumers must opt in explicitly.
	default:
		return fmt.Errorf("op_source.attestation_head must be %q, %q or %q, got %q", HeadFinalized, HeadSafe, HeadUnsafe, c.AttestationHead)
	}
	if c.StatePath == "" {
		return fmt.Errorf("op_source.state_path is required")
	}
	return nil
}
