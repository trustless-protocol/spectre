package prover

import (
	"strings"
	"testing"
	"time"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttypes "github.com/cometbft/cometbft/types"
)

// worstCaseVoteSignBytes builds the largest VoteSignBytes CometBFT will
// produce: a chain id at MaxChainIDLen, height/round wide enough to force
// full-width sfixed64 encodings, a PartSetHeader.Total at the varint maximum,
// and a far-future timestamp. Anything a real chain emits is shorter.
func worstCaseVoteSignBytes() []byte {
	v := &cmtproto.Vote{
		Type:   cmtproto.PrecommitType,
		Height: 1 << 62,
		Round:  1 << 30,
		BlockID: cmtproto.BlockID{
			Hash: make([]byte, 32),
			PartSetHeader: cmtproto.PartSetHeader{
				Total: ^uint32(0),
				Hash:  make([]byte, 32),
			},
		},
		// Year ~5138; seconds needs a wider varint than any realistic block time.
		Timestamp: time.Unix(1<<37, 999999999).UTC(),
	}
	return cmttypes.VoteSignBytes(strings.Repeat("x", cmttypes.MaxChainIDLen), v)
}

// TestMaxMsgLenBounds pins both bounds that fix MaxMsgLen at 175.
//
// Lower: every canonical vote CometBFT can emit must fit, otherwise
// buildCircuitAssignment rejects the batch and the relayer cannot prove that
// block at all. This fails if a CometBFT upgrade widens a field or raises
// MaxChainIDLen.
//
// Upper: 64 bytes of R||A plus the buffer must still pad into two SHA-512
// blocks. Crossing into a third costs a full extra permutation per slot
// (~17% of the circuit), so this fails if MaxMsgLen is raised past 175.
func TestMaxMsgLenBounds(t *testing.T) {
	worst := len(worstCaseVoteSignBytes())
	if worst > MaxMsgLen {
		t.Fatalf("worst-case VoteSignBytes is %d bytes but MaxMsgLen is %d: "+
			"real votes would be rejected before proving. Raise MaxMsgLen to cover the "+
			"worst case — but past 175 it adds a third SHA-512 block per slot, so "+
			"re-measure the circuit if it must exceed that.",
			worst, MaxMsgLen)
	}
	t.Logf("worst-case VoteSignBytes = %d bytes, MaxMsgLen = %d (%d bytes of slack)",
		worst, MaxMsgLen, MaxMsgLen-worst)

	// Mirror of sha512 FixedLengthSum's padding: at least 17 bytes (0x80 plus
	// the 16-byte length field), rounded up to the 128-byte block size.
	const sha512Block = 128
	hashed := 64 + MaxMsgLen // R || A || msg
	padded := hashed + (sha512Block - hashed%sha512Block)
	if sha512Block-hashed%sha512Block <= 16 {
		padded += sha512Block
	}
	if blocks := padded / sha512Block; blocks != 2 {
		t.Fatalf("MaxMsgLen=%d makes H_RAM span %d SHA-512 blocks, want 2: "+
			"each extra block is a full permutation per slot", MaxMsgLen, blocks)
	}
}
