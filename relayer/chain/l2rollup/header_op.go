package l2rollup

import (
	"context"
	"fmt"
	"math/big"

	relayerclient "relayer/client"

	attestorpb "attestor/types/attestor"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// maxGameStepDown bounds how far BuildHeader walks back from the attestor's frontier
// looking for a game that is already visible at the pinned (finalized) L1 block. The
// gap is normally one or two games; the cap only stops a pathological walk when the
// pinned client has fallen far behind.
const maxGameStepDown = 16

// OPProfile is the subset of the rollup profile the OP header builder needs to
// assemble proofs (parsed from rollup_profile by the caller). The verifier-only fields
// (output_root_format, l2_header_fork, commitment_slot, root_claim_bytecode_offset) are
// not needed here — the relayer sends witnesses, the wasm client re-derives + hashes.
type OPProfile struct {
	DisputeGameFactory ethcommon.Address // L1 factory account
	GameListSlot       ethcommon.Hash    // Solidity slot of the factory game list
	L2Router           ethcommon.Address // L2 ICS26Router (router_proof target)
	L1ClientID         string            // shared ETH (08-wasm) client id on Cosmos
}

// opNodeOutput is the subset of optimism_outputAtBlock the builder reads (confirmed
// on the OP devnet). NOTE message_passer_storage_root is the `withdrawalStorageRoot`
// field.
type opNodeOutput struct {
	Version                  ethcommon.Hash `json:"version"`
	StateRoot                ethcommon.Hash `json:"stateRoot"`
	MessagePasserStorageRoot ethcommon.Hash `json:"withdrawalStorageRoot"`
	BlockRef                 struct {
		Hash   ethcommon.Hash `json:"hash"`
		Number uint64         `json:"number"`
	} `json:"blockRef"`
}

// opStackHeaderBuilder assembles an OpStackHeader: it proves a DisputeGameFactory
// game commitment in the L1 state the shared ETH client trusts, binds it to the
// canonical L2 header via the output root, and proves the L2 router state.
type opStackHeaderBuilder struct {
	l1       *ethclient.Client // L1 exec (factory/game eth_getProof + eth_getCode)
	l2       *ethclient.Client // L2 exec (l2_header + router eth_getProof)
	opNode   *rpc.Client       // op-node (optimism_outputAtBlock)
	cosmos   cosmosClientStateReader
	attestor AttestorClient
	srcChain string
	profile  OPProfile
	// includeProvisional must match the Source's: the source decides WHICH heights are
	// relayable, this decides which root is proven for them. Two different answers
	// would gate on one rule and prove with another.
	includeProvisional bool
}

// cosmosClientStateReader reads the shared ETH client's trusted L1 slot/block from
// Cosmos, and can ask for it to be advanced. Defined here so the builder is testable
// without a live Cosmos node.
type cosmosClientStateReader interface {
	EthClientLatestSlotAndBlock(l1ClientID string) (slot, execBlock uint64, err error)
	// UpdateEthClientIfStale advances the pinned ETH client if it is too far behind the
	// L1 head to prove against. The builder can only prove at the block that client
	// trusts, so it asks for freshness at the moment it needs it rather than relying
	// on something else to have kept the client current (#276). Returning an error is
	// not fatal on its own — the caller still tries with whatever height the client
	// has — so a deployment where another direction advances the client keeps working.
	UpdateEthClientIfStale(ctx context.Context, l1ClientID string) error
}

// NewOPStackHeaderBuilder wires the OP-Stack builder to the L1/L2 exec RPCs, the
// op-node RPC (for optimism_outputAtBlock), the Cosmos client-state reader, the
// attestor (whose AttestedRootAtOrBelow selects the game to prove), and the parsed
// profile. BuildHeader derives provisional-vs-finalized selection from each request
// so it proves the same game the source gated on.
func NewOPStackHeaderBuilder(l1, l2 *ethclient.Client, opNode *rpc.Client, cosmos cosmosClientStateReader, attestor AttestorClient, srcChain string, profile OPProfile, includeProvisional bool) *opStackHeaderBuilder {
	return &opStackHeaderBuilder{
		l1: l1, l2: l2, opNode: opNode, cosmos: cosmos,
		attestor: attestor, srcChain: srcChain, profile: profile,
		includeProvisional: includeProvisional,
	}
}

func (a *opStackHeaderBuilder) Name() string { return "l2-opstack" }

// BuildHeader assembles the OpStackHeader for the attested game at or below request.Height,
// and returns the L2 block that game commits (which the header binds to, so the client
// must advance by it, not the request). RPC/availability failures are returned plain
// (the generic Builder wraps them chain.Retryable).
func (a *opStackHeaderBuilder) BuildHeader(ctx context.Context, request HeaderRequest) (ClientMessage, uint64, error) {
	// 1. The L1 block the shared ETH client trusts — prove factory/game there.
	// Ask for the pinned ETH client to be current first: this builder can only
	// prove at the block that client trusts, so its freshness is a precondition,
	// not a background nicety (#276). A failure here is logged by the
	// implementation and does not abort — the client may still be fresh enough,
	// or another direction may be advancing it.
	_ = a.cosmos.UpdateEthClientIfStale(ctx, a.profile.L1ClientID)
	beaconSlot, l1Block, err := a.cosmos.EthClientLatestSlotAndBlock(a.profile.L1ClientID)
	if err != nil {
		return nil, 0, fmt.Errorf("l2-opstack: read shared ETH client height: %w", err)
	}
	l1BlockBig := new(big.Int).SetUint64(l1Block)
	l1Header, err := a.l1.HeaderByNumber(ctx, l1BlockBig)
	if err != nil {
		return nil, 0, fmt.Errorf("l2-opstack: L1 header at %d: %w", l1Block, err)
	}

	// 2. The attestor picks the game to prove: the highest independently re-derived
	//    game at or below l2Height (NOT simply the newest factory game — the newest
	//    can be a game the attestor refused or has not yet verified, whose root claim
	//    would not bind to optimism_outputAtBlock for the same block). It returns both
	//    the game_index and the L2 block that game commits; everything below is proven
	//    at that committed block so the output root, L2 header, and root claim agree.
	// Games are posted at the L1 HEAD, while the pinned Ethereum client only advances
	// to FINALIZED L1 — so the newest attested game is routinely absent from the
	// factory's list at the block this proof pins to. On a chain that posts games at
	// roughly the rate finality advances, the attestor's frontier stays permanently
	// ahead of what is provable, and demanding the frontier game deadlocks: the target
	// rises exactly as fast as the pinned block does, so the update never lands and
	// the client never moves. (Measured on a devnet: target 7202→7239→7276 while the
	// factory held 151→152→153 games at the pinned block, client frozen throughout.)
	//
	// So walk DOWN instead: take the highest attested game that is actually visible at
	// the pinned block. Every candidate is still an attested root — the attestor's
	// verdict is what makes a game provable, and stepping down only ever picks an
	// older one it already approved. The client then advances monotonically and
	// crosses any given packet height once finality carries a game past it.
	gameCount, err := factoryGameCount(ctx, a.l1, a.profile.DisputeGameFactory, l1BlockBig)
	if err != nil {
		return nil, 0, annotatePinnedL1(ctx, a.l1, "l2-opstack: factory gameCount", a.profile.L1ClientID, l1Block, err)
	}

	target := request.Height
	var (
		root            *attestorpb.AttestedRoot
		gameIndex       uint64
		committedHeight uint64
	)
	for attempt := 0; ; attempt++ {
		if attempt == maxGameStepDown {
			return nil, 0, fmt.Errorf(
				"l2-opstack: no attested game at or below L2 height %d is visible at the pinned L1 block %d "+
					"after %d step(s) (factory holds %d game(s) there; Ethereum client %s advances only to "+
					"finalized L1)",
				request.Height, l1Block, maxGameStepDown, gameCount, a.profile.L1ClientID)
		}
		var found bool
		root, found, err = a.attestor.AttestedRootAtOrBelow(ctx, a.srcChain, target, a.includeProvisional)
		if err != nil {
			return nil, 0, err
		}
		if !found {
			return nil, 0, fmt.Errorf("l2-opstack: attestor has no %s game at or below L2 height %d", request.Finality, target)
		}
		if root.GetSource() != "game" {
			return nil, 0, fmt.Errorf("l2-opstack: attested root at height %d is source %q, want game", target, root.GetSource())
		}
		gameIndex = root.GetGameIndex()
		committedHeight = root.GetL2BlockNumber()
		if gameIndex < gameCount {
			break // visible at the pinned block — provable
		}
		if committedHeight == 0 {
			return nil, 0, fmt.Errorf("l2-opstack: attested game %d commits L2 height 0", gameIndex)
		}
		target = committedHeight - 1 // drop below this game and ask again
	}

	gameSlot, err := gameListElementSlot(a.profile.GameListSlot, gameIndex)
	if err != nil {
		return nil, 0, err
	}

	factoryProof, err := relayerclient.EthGetProof(ctx, a.l1, a.profile.DisputeGameFactory, []ethcommon.Hash{gameSlot}, l1BlockBig)
	if err != nil {
		return nil, 0, annotatePinnedL1(ctx, a.l1, "l2-opstack: dispute game factory proof", a.profile.L1ClientID, l1Block, err)
	}
	if len(factoryProof.Storage) == 0 {
		return nil, 0, fmt.Errorf("l2-opstack: factory game-list proof missing storage entry")
	}
	gameAddr, err := addressFromStorageValue(factoryProof.Storage[0].Value)
	if err != nil {
		return nil, 0, err
	}
	if gameAddr == (ethcommon.Address{}) {
		return nil, 0, fmt.Errorf("l2-opstack: factory game-list entry %d is empty", gameIndex)
	}

	gameAccountProof, err := relayerclient.EthGetProof(ctx, a.l1, gameAddr, nil, l1BlockBig)
	if err != nil {
		return nil, 0, annotatePinnedL1(ctx, a.l1, "l2-opstack: game account proof", a.profile.L1ClientID, l1Block, err)
	}
	gameRuntime, err := a.l1.CodeAt(ctx, gameAddr, l1BlockBig)
	if err != nil {
		return nil, 0, annotatePinnedL1(ctx, a.l1,
			fmt.Sprintf("l2-opstack: game runtime (eth_getCode %s)", gameAddr), a.profile.L1ClientID, l1Block, err)
	}

	// 3. Output-root preimage from op-node at the committed block.
	var out opNodeOutput
	if err := a.opNode.CallContext(ctx, &out, "optimism_outputAtBlock", hexUint64(committedHeight)); err != nil {
		return nil, 0, fmt.Errorf("l2-opstack: optimism_outputAtBlock(%d): %w", committedHeight, err)
	}

	// 4. Canonical L2 header + 5. router proof against its state root, at the committed block.
	committedBig := new(big.Int).SetUint64(committedHeight)
	l2Header, err := a.l2.HeaderByNumber(ctx, committedBig)
	if err != nil {
		return nil, 0, fmt.Errorf("l2-opstack: L2 header at %d: %w", committedHeight, err)
	}
	routerProof, err := relayerclient.EthGetProof(ctx, a.l2, a.profile.L2Router, nil, committedBig)
	if err != nil {
		return nil, 0, err
	}

	return &OpStackHeader{
		BeaconSlot:       beaconSlot,
		L1StateRoot:      l1Header.Root.Bytes(),
		FactoryProof:     EvmAccountProof{Proof: factoryProof.AccountProof},
		GameIndex:        gameIndex,
		GameProof:        EvmStorageProof{Key: gameSlot.Bytes(), Value: factoryProof.Storage[0].Value, Proof: factoryProof.Storage[0].Proof},
		GameAccountProof: EvmAccountProof{Proof: gameAccountProof.AccountProof},
		GameRuntime:      gameRuntime,
		OutputRootProof: OutputRootProof{
			Version:                  out.Version.Bytes(),
			StateRoot:                out.StateRoot.Bytes(),
			MessagePasserStorageRoot: out.MessagePasserStorageRoot.Bytes(),
			LatestBlockhash:          out.BlockRef.Hash.Bytes(),
		},
		L2Header:    toCanonicalHeader(l2Header),
		RouterProof: EvmAccountProof{Proof: routerProof.AccountProof},
	}, committedHeight, nil
}

var _ HeaderBuilder = (*opStackHeaderBuilder)(nil)
