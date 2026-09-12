// Package chain defines the adapter contract that lets one relayer process
// serve many source/destination chains (Cosmos, and EVM L2s such as Arbitrum,
// Optimism, Base) without the core knowing any chain's specifics.
//
// It replaces the hardcoded Cosmos<->ETH coupling in the current core:
//
//	services.EventListener.SubscribeCosmos/SubscribeEth   -> Source.Subscribe
//	services.TransactionHandler (mixed Cosmos+Eth methods) -> Destination.*
//	services.Prover.GenerateProof (Groth16 only)           -> ClientUpdateBuilder.Build
//	the hardcoded finality/refresh in StartLoop            -> Source.LatestHeight + ClientExpiresAt
//
// Method names follow IBC / relayer convention (UpdateClient, RelayPackets,
// HasPacketReceipt, QueryHeader, LatestHeight) so the contract reads the same
// as any Cosmos relayer.
//
// This is the v0 contract for review. Payload types are intentionally opaque
// ([]byte / small structs) so each adapter owns its own encoding; later wiring
// slices refine them. Adapters must not leak chain-specific types across this
// boundary — that is the whole point.
package chain

import (
	"context"
	"fmt"
	"time"
)

// ChainType identifies a chain family so the registry can pick an adapter.
// Instances of the same family differ only by config (chain id, RPC, contracts).
type ChainType string

const (
	Cosmos   ChainType = "cosmos"
	Ethereum ChainType = "ethereum" // L1
	OPStack  ChainType = "opstack"  // Optimism + Base (same stack; differ by config/binary)
	Arbitrum ChainType = "arbitrum" // Nitro
)

// EventType classifies a relayable event by its IBC packet-lifecycle message,
// so the core can route it without decoding chain-specific payloads.
type EventType int

const (
	SendPacket EventType = iota
	AckPacket
	TimeoutPacket
)

// String names the event by the message it becomes on the destination, so logs
// distinguish a forward delivery from a returning acknowledgement without the
// reader decoding an integer.
func (t EventType) String() string {
	switch t {
	case SendPacket:
		return "recv"
	case AckPacket:
		return "ack"
	case TimeoutPacket:
		return "timeout"
	default:
		return fmt.Sprintf("unknown(%d)", int(t))
	}
}

// Event is a relayable packet observed on a source chain. Raw carries the
// adapter-encoded packet; the core only reads the routing fields.
type Event struct {
	Type   EventType
	Height uint64 // source height the event was observed at (gated by LatestHeight)
	// Sequence is the packet's IBC sequence, duplicated out of Raw purely so the
	// core can name a packet in logs. Decoding Raw for that would put a proto
	// dependency (and adapter-specific knowledge) into the generic relay loop,
	// which is the thing this package exists to avoid — so the adapters, which
	// already hold the decoded packet, copy it in.
	Sequence uint64
	ClientID string   // destination client id this event routes to (event-stream partition key)
	Raw      []byte   // proto-marshaled channeltypesv2.Packet
	AckBytes [][]byte // app acknowledgements, populated only for AckPacket events
}

// ClientUpdate is a light-client update for the destination's client-of-source
// (the payload of an IBC MsgUpdateClient).
//
// Height is exposed so the core can enforce the append-only / monotonic
// invariant generically — reject an update at or below the destination client's
// current height before submitting, which catches reorg-attestations and
// replays without decoding Payloads. (Destination.UpdateClient must still reject
// non-monotonic updates itself; this is defense in depth at the core layer.)
//
// Payloads are adapter-owned, destination-specific update bytes. Most adapters
// emit exactly one payload (a Groth16 proof bundle for Cosmos->X or a trustless
// rollup-proof JSON ClientMessage for L2->Cosmos); beacon updates may need an
// ordered sequence when crossing sync-committee periods. Only the destination
// adapter understands their contents. An empty list means the client is already
// current and no transaction is needed.
type ClientUpdate struct {
	Height    uint64
	Payloads  [][]byte
	TrustedAt time.Time
}

// RelayPacket is a source packet together with the source proof the destination
// needs to submit it. The destination adapter turns it into the concrete IBC
// message (MsgRecvPacket / MsgAcknowledgement / MsgTimeout) and submits it.
type RelayPacket struct {
	Type     EventType // recv (SendPacket), ack, or timeout
	Sequence uint64    // carried from the source Event, for logging and wait bookkeeping
	Packet   []byte    // proto-marshaled channeltypesv2.Packet
	Proof    []byte    // from Source.MembershipProof / NonMembershipProof
	Height   uint64    // the source height the proof is against
	AckBytes [][]byte  // app acknowledgements, for AckPacket only
}

// Source is a chain whose state we prove to a destination.
type Source interface {
	Chain() ChainType

	// Subscribe streams relayable events to handler in BATCHES until ctx is
	// cancelled. It must carry the existing bidirectional gap-recovery guarantees
	// (startup lookback, periodic ticker, reconnect-gap) — reused from the generic
	// subscriber, not reimplemented per chain.
	//
	// CALL ONCE PER PROCESS LIFETIME. Subscribe owns and drains its subscriber and
	// batch-handoff workers before returning; ctx cancellation must interrupt any
	// reconnect wait, live watch, recovery scan, or RPC in flight.
	//
	// Batching (not per-event) is deliberate: it lets the destination fold N
	// packet messages into one multicall (amortizing the ~21k per-tx intrinsic gas
	// and one client-update over the batch) and lets the module perform batch-level
	// waits (source-provability / finality) once instead of per packet — the
	// handleCosmos/handleEth structure, expressed generically.
	//
	// handler returns the indices (into the passed slice) of events that failed
	// TRANSIENTLY and must be re-queued for retry; Subscribe re-queues exactly
	// those. An empty/nil return means every event was handled or permanently
	// dropped (the handler logs permanents). Indices are always valid positions in
	// the slice handler was given. The ctx passed to handler is Subscribe's ctx, so
	// the handler can abort in-flight work (proofs, reads) on shutdown.
	Subscribe(ctx context.Context, handler func(context.Context, []Event) []int) error

	// LatestHeight returns the highest source height that is safe to relay right
	// now. THIS IS THE FINALITY / CONFIRMATION POLICY HOOK — "latest" means the
	// latest height this adapter is willing to prove:
	//   - Cosmos (BFT instant finality): the latest committed height.
	//   - EVM L2, full-finality mode:    the highest L1-finalized L2 height.
	//   - EVM L2, "skip finality" mode:   the latest unsafe/soft head (reorg-exposed).
	// Changing the trust/latency tradeoff = changing only this method.
	LatestHeight(ctx context.Context) (uint64, error)

	// RelayableHeight returns the highest source height whose packets can be PROVEN
	// AND relayed right now. It is the CHEAP precondition the relay module checks
	// before doing any expensive client-update / proof work: a packet above this
	// height is re-queued without burning a proof, instead of proving repeatedly
	// while its source state is not yet available (the legacy waitCosmosAppHash /
	// waitBeaconFinality guards, expressed as a value instead of a blocking wait).
	// It encodes each chain's state-availability lag:
	//   - Cosmos: latest committed height minus the AppHash lag — a commitment
	//     written at H is only reflected in the queryable app state at H+2, so
	//     latest-2.
	//   - Ethereum L1: the execution block number of the latest FINALIZED beacon
	//     header — the beacon light client only proves against finalized state.
	// It is always <= LatestHeight.
	RelayableHeight(ctx context.Context) (uint64, error)

	// QueryHeader returns the adapter-owned data a ClientUpdateBuilder needs to
	// build the update for source state at the given height.
	QueryHeader(ctx context.Context, height uint64) ([]byte, error)

	// MembershipProof proves the packet's commitment EXISTS in the source state at
	// height — for recv/ack. eventType (SendPacket vs AckPacket) selects which
	// membership path is proven (commitment vs ack), which sit at different
	// ICS-24 paths. NonMembershipProof proves the packet receipt is ABSENT — for
	// timeout (a single unambiguous path). Both mirror the existing
	// membershipFn / nonMembershipFn used by the packet relay. packet is the
	// proto-marshaled channeltypesv2.Packet.
	MembershipProof(ctx context.Context, packet []byte, height uint64, eventType EventType) ([]byte, error)
	NonMembershipProof(ctx context.Context, packet []byte, height uint64) ([]byte, error)
}

// Destination is a chain that hosts the light client of some source and where
// we submit updates + packet messages.
type Destination interface {
	Chain() ChainType

	// UpdateClient advances the destination's light-client-of-source. Must be
	// append-only safe: reject a non-monotonic / non-contiguous update rather
	// than corrupt client state.
	UpdateClient(ctx context.Context, clientID string, update ClientUpdate) error

	// RelayPackets builds and submits the concrete IBC packet messages
	// (recv/ack/timeout) from the given source packets + proofs.
	RelayPackets(ctx context.Context, packets []RelayPacket) error

	// HasPacketReceipt reports whether a packet was already delivered on this
	// chain, so timeout scanners can skip already-settled packets (the #232
	// pre-check, generalized). Prevents advancing state on already-resolved
	// packets.
	HasPacketReceipt(ctx context.Context, packet []byte) (bool, error)

	// ClientExpiresAt returns when the light-client-of-source on this chain will
	// expire, so the generic refresh routine can schedule updates ahead of it. A
	// zero time means "no expiry" (e.g. a permissioned client with no trusting
	// period).
	//
	// It also returns the trusting period that expiry was derived from. The
	// caller needs it to size its safety margin: a margin fixed in absolute time
	// is either wasteful against a period measured in days or unreachable against
	// one measured in minutes, and every implementation already has the number --
	// it computes the expiry from it. A zero period means "unknown", and the
	// caller falls back to its own default.
	ClientExpiresAt(ctx context.Context, clientID string) (expiresAt time.Time, trustingPeriod time.Duration, err error)
}

// FoldingDestination is an optional destination capability for submitting a
// client update followed by packet messages in one atomic transaction. The
// relay module uses it only when SupportsUpdatePacketFolding returns true;
// otherwise it retains Destination's update-then-packets fallback.
//
// update.Payloads is non-empty when RelayWithUpdate is called. Implementations
// must apply the update before the packet messages and revert the whole
// transaction if any message fails.
type FoldingDestination interface {
	SupportsUpdatePacketFolding() bool
	RelayWithUpdate(ctx context.Context, clientID string, update ClientUpdate, packets []RelayPacket) error
}

// ClientUpdateBuilder turns a source header into a destination-verifiable
// ClientUpdate. The strategy is chosen by what the DESTINATION's light client
// verifies, not one-size-fits-all:
//   - "groth16":     batched Ed25519 ZK proof — the destination runs a ZK
//     Tendermint client (SpectreClient). Cosmos -> any EVM. Relayer generates
//     the proof.
//   - "beacon":      assemble the Ethereum sync-committee + storage/account
//     proof update — the destination runs the 08-wasm beacon light client.
//     ETH L1 -> Cosmos. NO relayer-side proving or re-execution; the wasm
//     client verifies BLS + Merkle itself. This is the existing eth_to_cosmos
//     path.
//   - "l2-opstack" / "l2-arbitrum": assemble the trustless rollup proof — the L1
//     rollup-contract witnesses (AnchorStateRegistry for OP-Stack, RollupCore for
//     Arbitrum) + the RLP L2 header + L2 IBC-handler account proof, as a JSON
//     ClientMessage. L2 -> Cosmos. NO relayer signature or re-execution; the L2
//     wasm client verifies the rollup proofs against the shared ETH client. The
//     Source.HeadKind selects the requested latency tier; the request is carried
//     to the builder and the wasm client independently authenticates the claimed
//     finality evidence before applying its membership policy.
type ClientUpdateBuilder interface {
	Name() string
	Build(ctx context.Context, header []byte) (ClientUpdate, error)
}

// --- Registry: adding a chain/builder = register + config, zero core edits. ---

type (
	SourceFactory              func(cfg []byte) (Source, error)
	DestinationFactory         func(cfg []byte) (Destination, error)
	ClientUpdateBuilderFactory func(cfg []byte) (ClientUpdateBuilder, error)
)

var (
	sources      = map[ChainType]SourceFactory{}
	destinations = map[ChainType]DestinationFactory{}
	builders     = map[string]ClientUpdateBuilderFactory{}
)

// RegisterSource wires a source-chain family. Call from an adapter's init().
func RegisterSource(chain ChainType, factory SourceFactory) { sources[chain] = factory }

// RegisterDestination wires a destination-chain family.
func RegisterDestination(chain ChainType, factory DestinationFactory) { destinations[chain] = factory }

// RegisterClientUpdateBuilder wires a client-update strategy by name (config field `builder`).
func RegisterClientUpdateBuilder(name string, factory ClientUpdateBuilderFactory) {
	builders[name] = factory
}

// NewSource / NewDestination / NewClientUpdateBuilder build a configured adapter,
// or report an unknown chain/builder so a typo in config fails loudly at startup.
func NewSource(chain ChainType, cfg []byte) (Source, error) {
	factory, ok := sources[chain]
	if !ok {
		return nil, &UnknownAdapterError{Category: "source", Name: string(chain)}
	}
	return factory(cfg)
}

func NewDestination(chain ChainType, cfg []byte) (Destination, error) {
	factory, ok := destinations[chain]
	if !ok {
		return nil, &UnknownAdapterError{Category: "destination", Name: string(chain)}
	}
	return factory(cfg)
}

func NewClientUpdateBuilder(name string, cfg []byte) (ClientUpdateBuilder, error) {
	factory, ok := builders[name]
	if !ok {
		return nil, &UnknownAdapterError{Category: "client-update builder", Name: name}
	}
	return factory(cfg)
}

// UnknownAdapterError is returned when config names an unregistered adapter.
type UnknownAdapterError struct {
	Category string
	Name     string
}

func (e *UnknownAdapterError) Error() string {
	return "chain: no registered " + e.Category + " for " + e.Name
}
