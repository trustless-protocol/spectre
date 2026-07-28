package l2rollup

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"relayer/chain"
	relayerclient "relayer/client"
	"relayer/services"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
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

func (k HeadKind) validate() error {
	switch k {
	case Unsafe, Safe, Finalized:
		return nil
	default:
		return fmt.Errorf("l2 source: invalid head kind %d", k)
	}
}

func (k HeadKind) String() string {
	switch k {
	case Unsafe:
		return "unsafe"
	case Safe:
		return "safe"
	case Finalized:
		return "finalized"
	default:
		return fmt.Sprintf("unknown(%d)", k)
	}
}

// Source is the EVM-L2 chain.Source (L2->Cosmos). An L2 is EVM, so packet
// membership proofs are the same eth_getProof account+storage proofs the ETH L1
// source builds — verified by the L2 light client against the L2 world state_root
// (account proof) then the IBC-handler storage_root (storage proof), per Dũng P2.
// The L2 difference is entirely the finality policy here.
type Source struct {
	chainType  chain.ChainType // OPStack or Arbitrum (for Chain())
	eth        *ethclient.Client
	headKind   HeadKind
	l2ClientID string            // the ICS26Router client id on the L2 (event partition key)
	router     ethcommon.Address // the L2 ICS26Router address (from rollup_profile.common.l2_router)

	// attestor gates RelayableHeight on independent L1 re-derivation (#240). When
	// nil, RelayableHeight falls back to the raw L2 head — the skip-finality interim
	// (or an explicitly attestor-less deployment). srcChainID is the attestor's
	// src_chain key (distinct from the on-L2 client id).
	attestor   AttestorClient
	srcChainID string
}

// Source satisfies chain.Source.
var _ chain.Source = (*Source)(nil)

// NewSource wires an L2 source. router is the L2 ICS26Router address (the packet
// membership proofs are taken against its storage_root). attestor may be nil
// (raw-head interim); when set, srcChainID identifies this L2 to the attestor.
func NewSource(chainType chain.ChainType, eth *ethclient.Client, headKind HeadKind, l2ClientID, srcChainID string, router ethcommon.Address, attestor AttestorClient) *Source {
	return &Source{
		chainType:  chainType,
		eth:        eth,
		headKind:   headKind,
		l2ClientID: l2ClientID,
		router:     router,
		attestor:   attestor,
		srcChainID: srcChainID,
	}
}

func (s *Source) Chain() chain.ChainType { return s.chainType }

// LatestHeight is the latest L2 block visible at the configured head tag — how far
// the subscriber can see packets. It is NOT the trust gate (that is RelayableHeight).
func (s *Source) LatestHeight(ctx context.Context) (uint64, error) { return s.head(ctx) }

// RelayableHeight is the highest L2 height a packet may be proven at. With an
// attestor it is the attestor's independently re-derived frontier (AttestedUpTo);
// includeProvisional accepts the Safe-but-not-finalized head, so only Finalized
// demands non-provisional roots. Without an attestor it degrades to the raw L2 head
// (skip-finality). A not-yet-attested source returns 0 — nothing is relayable yet,
// so the module waits rather than relaying an unverified height.
func (s *Source) RelayableHeight(ctx context.Context) (uint64, error) {
	if s.attestor == nil {
		return s.head(ctx)
	}
	cctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	root, found, err := s.attestor.AttestedUpTo(cctx, s.srcChainID, s.headKind != Finalized)
	if err != nil {
		return 0, fmt.Errorf("l2 source: attested-up-to (kind=%d): %w", s.headKind, err)
	}
	if !found {
		return 0, nil // nothing attested yet — the module waits
	}
	return root.GetL2BlockNumber(), nil
}

// head reads the L2 head at the configured head-kind tag.
func (s *Source) head(ctx context.Context) (uint64, error) {
	if err := s.headKind.validate(); err != nil {
		return 0, err
	}
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

// QueryHeader returns the selected finality byte plus the target L2 height as an
// 8-byte big-endian integer for the L2 client-update builder.
func (s *Source) QueryHeader(_ context.Context, height uint64) ([]byte, error) {
	if err := s.headKind.validate(); err != nil {
		return nil, err
	}
	b := make([]byte, 9)
	b[0] = byte(s.headKind)
	binary.BigEndian.PutUint64(b[1:], height)
	return b, nil
}

// Subscribe (in subscribe.go) streams L2 ICS26Router packet events to the handler.

// MembershipProof builds the L2 eth_getProof (account proof vs the L2 world
// state_root, storage proof vs the ICS26Router storage_root) that the packet
// commitment (send) or acknowledgement exists at height — the exact mechanics of the
// ETH L1 source, reused against the L2 router. height is the L2 block the module
// advanced the destination L2 wasm client to (proofHeight = m.lastHeight), so no
// Cosmos-side read is needed here. Do NOT use a Cosmos app_hash — the two EVM roots
// are the L2 world state_root and the router account storage_root (Dũng P2).
func (s *Source) MembershipProof(_ context.Context, packet []byte, height uint64, eventType chain.EventType) ([]byte, error) {
	var pkt channeltypesv2.Packet
	if err := pkt.Unmarshal(packet); err != nil {
		return nil, fmt.Errorf("l2 source: decode packet: %w", err)
	}
	var clientID string
	var pathType byte
	switch eventType {
	case chain.SendPacket:
		// An L2->Cosmos send past its timeout can never be received on Cosmos, so
		// report it permanent and let the module DROP it (mirroring the ETH source's
		// dead-send pre-filter). Cosmos has ~wall-clock BFT time, so compare against
		// time.Now like the ETH path.
		//
		// LIMITATION: unlike the Cosmos->ETH path, this L2->Cosmos path wires NO
		// timeout scanner today (buildL2ToCosmosModule adds no WithTimeoutScanner /
		// WithPacketTracker), so a dropped send is NOT yet refunded on the L2 — the L2
		// escrow stays locked. The refund path (a TimeoutPacket back to the L2 rollup)
		// belongs to the Cosmos->L2 return direction and is tracked as a follow-up.
		if pkt.TimeoutTimestamp > 0 && uint64(time.Now().Unix()) >= pkt.TimeoutTimestamp {
			return nil, chain.Permanent(fmt.Errorf("l2 source: send seq=%d timed out and is dropped (no L2 refund scanner yet)", pkt.Sequence))
		}
		clientID, pathType = pkt.SourceClient, 1 // packet commitment
	case chain.AckPacket:
		clientID, pathType = pkt.DestinationClient, 3 // ack
	default:
		return nil, fmt.Errorf("l2 source: MembershipProof: unsupported event type %d", eventType)
	}
	path := services.EthPath(clientID, pkt.Sequence, pathType)
	return relayerclient.GetEthMembershipProof(
		s.eth, s.router, path,
		ethcommon.HexToHash(services.ICS26_IBC_STORAGE_SLOT),
		new(big.Int).SetUint64(height),
	)
}

// NonMembershipProof is unused on the L2 source relay path: L2-origin packet
// timeouts are handled by the async timeout scanner (as on the ETH source), not the
// Subscribe->relay flow, so no TimeoutPacket event reaches here.
func (s *Source) NonMembershipProof(_ context.Context, _ []byte, _ uint64) ([]byte, error) {
	return nil, fmt.Errorf("l2 source: NonMembershipProof unused (timeouts are scanner-handled)")
}
