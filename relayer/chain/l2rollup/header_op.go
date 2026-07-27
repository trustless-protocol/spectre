package l2rollup

import (
	"context"
	"fmt"
	"math/big"

	relayerclient "relayer/client"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

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
	l1                 *ethclient.Client // L1 exec (factory/game eth_getProof + eth_getCode)
	l2                 *ethclient.Client // L2 exec (l2_header + router eth_getProof)
	opNode             *rpc.Client       // op-node (optimism_outputAtBlock)
	cosmos             cosmosClientStateReader
	attestor           AttestorClient
	srcChain           string
	includeProvisional bool
	profile            OPProfile
}

// cosmosClientStateReader reads the shared ETH client's trusted L1 slot/block from
// Cosmos. Defined here so the builder is testable without a live Cosmos node.
type cosmosClientStateReader interface {
	EthClientLatestSlotAndBlock(l1ClientID string) (slot, execBlock uint64, err error)
}

// NewOPStackHeaderBuilder wires the OP-Stack builder to the L1/L2 exec RPCs, the
// op-node RPC (for optimism_outputAtBlock), the Cosmos client-state reader, the
// attestor (whose AttestedRootAtOrBelow selects the game to prove), and the parsed
// profile. includeProvisional must match the source's head policy (Unsafe/Safe accept
// provisional; Finalized does not) so the builder proves the same game the source gated
// on.
func NewOPStackHeaderBuilder(l1, l2 *ethclient.Client, opNode *rpc.Client, cosmos cosmosClientStateReader, attestor AttestorClient, srcChain string, includeProvisional bool, profile OPProfile) *opStackHeaderBuilder {
	return &opStackHeaderBuilder{
		l1: l1, l2: l2, opNode: opNode, cosmos: cosmos,
		attestor: attestor, srcChain: srcChain, includeProvisional: includeProvisional, profile: profile,
	}
}

func (a *opStackHeaderBuilder) Name() string { return "l2-opstack" }

// BuildHeader assembles the OpStackHeader for the attested game at or below l2Height,
// and returns the L2 block that game commits (which the header binds to, so the client
// must advance by it, not the request). RPC/availability failures are returned plain
// (the generic Builder wraps them chain.Retryable).
func (a *opStackHeaderBuilder) BuildHeader(ctx context.Context, l2Height uint64) (ClientMessage, uint64, error) {
	// 1. The L1 block the shared ETH client trusts — prove factory/game there.
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
	root, found, err := a.attestor.AttestedRootAtOrBelow(ctx, a.srcChain, l2Height, a.includeProvisional)
	if err != nil {
		return nil, 0, err
	}
	if !found {
		return nil, 0, fmt.Errorf("l2-opstack: attestor has no game at or below L2 height %d", l2Height)
	}
	if root.GetSource() != "game" {
		return nil, 0, fmt.Errorf("l2-opstack: attested root at height %d is source %q, want game", l2Height, root.GetSource())
	}
	gameIndex := root.GetGameIndex()
	committedHeight := root.GetL2BlockNumber()

	gameSlot, err := gameListElementSlot(a.profile.GameListSlot, gameIndex)
	if err != nil {
		return nil, 0, err
	}

	factoryProof, err := relayerclient.EthGetProof(a.l1, a.profile.DisputeGameFactory, []ethcommon.Hash{gameSlot}, l1BlockBig)
	if err != nil {
		return nil, 0, err
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

	gameAccountProof, err := relayerclient.EthGetProof(a.l1, gameAddr, nil, l1BlockBig)
	if err != nil {
		return nil, 0, err
	}
	gameRuntime, err := a.l1.CodeAt(ctx, gameAddr, l1BlockBig)
	if err != nil {
		return nil, 0, fmt.Errorf("l2-opstack: game runtime (eth_getCode %s): %w", gameAddr, err)
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
	routerProof, err := relayerclient.EthGetProof(a.l2, a.profile.L2Router, nil, committedBig)
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
