package l2rollup

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

// annotatePinnedL1 explains a failed L1 read in terms of the block it was pinned to.
//
// Every L2 header builder proves the rollup contracts (OP's DisputeGameFactory, or
// Arbitrum's RollupCore) at the L1 block the Cosmos-side Ethereum light client
// currently trusts — it cannot pick another block, because the L2 wasm client
// verifies those proofs against the L1 state root it reads back from that same
// client. So the pinned block only works while that client keeps advancing:
//
//   - once it falls outside the execution node's state window, eth_getProof fails
//     with a bare "historical state ... is not available" that names neither the
//     block nor the client;
//   - and a game/assertion covering a recent L2 block does not exist at an old L1
//     block at all, so even an archive node cannot make the proof meaningful.
//
// Naming the pinned block, the client, and how far behind the head it is turns that
// into a diagnosis. The head lookup only happens on the error path.
func annotatePinnedL1(ctx context.Context, l1 *ethclient.Client, tag, l1ClientID string, l1Block uint64, err error) error {
	if err == nil {
		return nil
	}
	head, headErr := l1.BlockNumber(ctx)
	if headErr != nil || head < l1Block {
		return fmt.Errorf("%s: L1 read at pinned block %d (Ethereum client %s on Cosmos): %w",
			tag, l1Block, l1ClientID, err)
	}
	return fmt.Errorf("%s: L1 read at pinned block %d (Ethereum client %s on Cosmos; L1 head %d, "+
		"%d blocks behind — that client must keep advancing for this proof to be servable): %w",
		tag, l1Block, l1ClientID, head, head-l1Block, err)
}
