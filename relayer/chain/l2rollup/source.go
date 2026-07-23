package l2rollup

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"relayer/chain"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// HeadKind is the L2 confirmation policy — the single anti-reorg knob — using the
// OP Stack / attestor head vocabulary (unsafe / safe / finalized). The on-chain L2
// light client verifies VALIDITY (L2 state derived from L1 rollup proofs); this
// decides how much reorg risk we accept for latency by choosing which L2 head we
// relay.
type HeadKind int

const (
	// Unsafe relays at the latest L2 block (skip-finality): lowest latency, exposed
	// to L2 reorgs / sequencer equivocation.
	Unsafe HeadKind = iota
	// Safe relays at the sequencer-safe head (batch posted to L1, not yet finalized).
	Safe
	// Finalized relays only at L2 blocks whose L1 batch is finalized: safe, highest
	// latency.
	Finalized
)

// Source is the EVM-L2 chain.Source (L2->Cosmos). An L2 is EVM, so packet
// membership proofs are the same eth_getProof account+storage proofs the ETH L1
// source builds — verified by the L2 light client against the L2 world state_root
// (account proof) then the IBC-handler storage_root (storage proof), per Dũng P2.
// The L2 difference is entirely the finality policy here.
type Source struct {
	chainType  chain.ChainType // OPStack or Arbitrum (for Chain())
	eth        *ethclient.Client
	headKind   HeadKind
	l2ClientID string // the ICS26Router client id on the L2 (event partition key)
}

// Source satisfies chain.Source.
var _ chain.Source = (*Source)(nil)

// NewSource wires an L2 source. TODO(Đức): take the full L2 config (RPC url, router
// address, head kind) once the multi-chain config schema is settled.
func NewSource(chainType chain.ChainType, eth *ethclient.Client, headKind HeadKind, l2ClientID string) *Source {
	return &Source{chainType: chainType, eth: eth, headKind: headKind, l2ClientID: l2ClientID}
}

func (s *Source) Chain() chain.ChainType { return s.chainType }

// LatestHeight applies the head policy: the unsafe head (skip-finality), the safe
// head, or the finalized head. THIS is where the trust/latency tradeoff lives.
func (s *Source) LatestHeight(ctx context.Context) (uint64, error) { return s.head(ctx) }

// RelayableHeight equals LatestHeight for an L2 in the trustless model: a packet at
// height H is provable once the L2 client is advanced to H. (When the attestor
// oracle — PR #240 — gates which height is safe, RelayableHeight will read its
// AttestedRootAtOrBelow instead; wired once that lands.)
func (s *Source) RelayableHeight(ctx context.Context) (uint64, error) { return s.head(ctx) }

// head reads the L2 head at the configured head-kind tag.
func (s *Source) head(ctx context.Context) (uint64, error) {
	tag := big.NewInt(int64(rpc.LatestBlockNumber)) // Unsafe
	switch s.headKind {
	case Safe:
		tag = big.NewInt(int64(rpc.SafeBlockNumber))
	case Finalized:
		tag = big.NewInt(int64(rpc.FinalizedBlockNumber))
	}
	cctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	h, err := s.eth.HeaderByNumber(cctx, tag)
	if err != nil {
		return 0, fmt.Errorf("l2 source: head (kind=%d): %w", s.headKind, err)
	}
	return h.Number.Uint64(), nil
}

// QueryHeader returns the target L2 height for the l2 client-update builder — the
// 8-byte big-endian encoding it decodes to know which L2 block header to assemble.
func (s *Source) QueryHeader(_ context.Context, height uint64) ([]byte, error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, height)
	return b, nil
}

// Subscribe streams L2 ICS26Router packet events to handler in batches.
//
// TODO(Đức): wire an L2 event listener — mirror the ETH subscriber (filtered by the
// per-source router client id, startup lookback + reconnect-gap recovery) against
// the L2 RPC, feeding the same batch-builder / waiting-backoff machinery. Not
// blocked on Dũng.
func (s *Source) Subscribe(_ context.Context, _ func(context.Context, []chain.Event) []int) error {
	return fmt.Errorf("l2 source: Subscribe not yet wired (needs the L2 event listener)")
}

// MembershipProof builds the L2 eth_getProof (account proof vs world state_root,
// storage proof vs the IBC-handler storage_root) that the packet commitment (recv)
// or ack exists at height — the same mechanics as the ETH L1 source
// (client.GetEthMembershipProof).
//
// TODO(Đức): wire client.GetEthMembershipProof once the L2 router address + storage
// slot are in config. Do NOT use a Cosmos app_hash — the two EVM roots are
// state_root (world) and the IBC account storage_root (Dũng P2).
func (s *Source) MembershipProof(_ context.Context, _ []byte, _ uint64, _ chain.EventType) ([]byte, error) {
	return nil, fmt.Errorf("l2 source: MembershipProof not yet wired (needs L2 router config)")
}

// NonMembershipProof proves the packet receipt is absent at height (timeout),
// verified against the same L2 roots. Same shape as MembershipProof.
func (s *Source) NonMembershipProof(_ context.Context, _ []byte, _ uint64) ([]byte, error) {
	return nil, fmt.Errorf("l2 source: NonMembershipProof not yet wired")
}
