package l2rollup

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"

	attestorpb "attestor/types/attestor"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// fakeCosmosReader returns a fixed trusted L1 slot/block (the builder proves the
// factory/game at this L1 block). For the devnet check we point it at the current L1
// head so the recently-proposed game is present.
type fakeCosmosReader struct{ slot, block uint64 }

func (f fakeCosmosReader) EthClientLatestSlotAndBlock(string) (uint64, uint64, error) {
	return f.slot, f.block, nil
}

// TestOPBuildHeader_Devnet exercises the full RPC assembly against a live local OP
// devnet. It is skipped unless OP_DEVNET_L1_RPC is set (the devnet is ephemeral, not
// part of CI). It confirms BuildHeader assembles a well-formed OpStackHeader from real
// factory/game/output-root/L2-header data and that it serializes to the wire envelope.
//
//	Run: OP_DEVNET_L1_RPC=... OP_DEVNET_L2_RPC=... OP_DEVNET_OPNODE=... \
//	     OP_DEVNET_FACTORY=... OP_DEVNET_ROUTER=... go test -run TestOPBuildHeader_Devnet -v ./chain/l2rollup/
func TestOPBuildHeader_Devnet(t *testing.T) {
	l1URL := os.Getenv("OP_DEVNET_L1_RPC")
	if l1URL == "" {
		t.Skip("OP_DEVNET_L1_RPC not set; skipping live OP devnet integration test")
	}
	ctx := context.Background()

	l1, err := ethclient.Dial(l1URL)
	if err != nil {
		t.Fatalf("dial L1: %v", err)
	}
	l2, err := ethclient.Dial(os.Getenv("OP_DEVNET_L2_RPC"))
	if err != nil {
		t.Fatalf("dial L2: %v", err)
	}
	opNode, err := rpc.Dial(os.Getenv("OP_DEVNET_OPNODE"))
	if err != nil {
		t.Fatalf("dial op-node: %v", err)
	}

	l1Head, err := l1.BlockNumber(ctx)
	if err != nil {
		t.Fatalf("L1 head: %v", err)
	}
	// A safe L2 block the op-node has an output for.
	var sync struct {
		SafeL2 struct {
			Number uint64 `json:"number"`
		} `json:"safe_l2"`
	}
	if err := opNode.CallContext(ctx, &sync, "optimism_syncStatus"); err != nil {
		t.Fatalf("optimism_syncStatus: %v", err)
	}
	l2Height := sync.SafeL2.Number
	if l2Height == 0 {
		t.Skip("op-node safe_l2 is 0; devnet not synced yet")
	}

	factory := ethcommon.HexToAddress(os.Getenv("OP_DEVNET_FACTORY"))
	profile := OPProfile{
		DisputeGameFactory: factory,
		GameListSlot:       ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000068"),
		L2Router:           ethcommon.HexToAddress(os.Getenv("OP_DEVNET_ROUTER")), // any deployed L2 account (e.g. a predeploy)
		L1ClientID:         "08-wasm-0",
	}

	// Stand in for the attestor: resolve the newest game via gameCount() (what the
	// attestor's AttestedRootAtOrBelow returns), and pair it with the safe L2 head.
	gameCountSel := crypto.Keccak256([]byte("gameCount()"))[:4]
	out, err := l1.CallContract(ctx, ethereum.CallMsg{To: &factory, Data: gameCountSel}, new(big.Int).SetUint64(l1Head))
	if err != nil {
		t.Fatalf("gameCount(): %v", err)
	}
	count := new(big.Int).SetBytes(out).Uint64()
	if count == 0 {
		t.Skip("factory has no games yet")
	}
	at := &fakeAttestor{found: true, root: &attestorpb.AttestedRoot{
		Source:        "game",
		Provenance:    &attestorpb.AttestedRoot_GameIndex{GameIndex: count - 1},
		L2BlockNumber: l2Height,
	}}
	builder := NewOPStackHeaderBuilder(l1, l2, opNode, fakeCosmosReader{slot: 1, block: l1Head}, at, "op-devnet", false, profile)

	msg, committed, err := builder.BuildHeader(ctx, l2Height)
	if err != nil {
		t.Fatalf("BuildHeader(l2Height=%d, l1Block=%d): %v", l2Height, l1Head, err)
	}
	if committed != l2Height {
		t.Errorf("committed height = %d, want %d", committed, l2Height)
	}
	h := msg.(*OpStackHeader)

	// Sanity: the assembled header must have non-empty proofs + the output-root words.
	if len(h.FactoryProof.Proof) == 0 {
		t.Error("factory proof is empty")
	}
	if len(h.GameProof.Proof) == 0 {
		t.Error("game proof is empty")
	}
	if len(h.GameRuntime) == 0 {
		t.Error("game runtime is empty")
	}
	if len(h.OutputRootProof.StateRoot) != 32 || len(h.OutputRootProof.MessagePasserStorageRoot) != 32 {
		t.Errorf("output root words wrong width: state=%d mp=%d", len(h.OutputRootProof.StateRoot), len(h.OutputRootProof.MessagePasserStorageRoot))
	}
	// The output root's state root must equal the canonical L2 header's state root
	// (the same binding verify() checks) — validates the two RPCs agree.
	if string(h.OutputRootProof.StateRoot) != string(h.L2Header.StateRoot) {
		t.Errorf("output stateRoot != l2_header.state_root:\n out=%x\n hdr=%x", h.OutputRootProof.StateRoot, h.L2Header.StateRoot)
	}

	raw, err := h.EncodeClientMessage()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatalf("envelope not JSON: %v", err)
	}
	t.Logf("assembled OpStackHeader: l2Height=%d game_index=%d l1Block=%d bytes=%d",
		l2Height, h.GameIndex, l1Head, len(raw))
	t.Logf("output_root state=%x mp=%x blockhash=%x",
		h.OutputRootProof.StateRoot, h.OutputRootProof.MessagePasserStorageRoot, h.OutputRootProof.LatestBlockhash)
}
