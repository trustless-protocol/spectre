// generate_golden produces VoteSignBytes golden vectors from cometbft and
// prints them as Go constants for pasting into ecip-gnark canonvote tests.
//
// Usage: go run ./prover/generate_golden > /tmp/canonvote_golden.txt
package main

import (
	"encoding/hex"
	"fmt"
	"time"

	cmttypes "github.com/cometbft/cometbft/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

type vector struct {
	name             string
	chainID          string
	height           int64
	round            int32
	blockIDHash      []byte
	partSetTotal    uint32
	partSetHash     []byte
	timestampSec     int64
	timestampNano   int32
}

func main() {
	makeBytes := func(seed byte, n int) []byte {
		b := make([]byte, n)
		for i := range b {
			b[i] = seed + byte(i)
		}
		return b
	}
	vectors := []vector{
		{
			name:          "typical_precommit",
			chainID:       "cosmoshub-4",
			height:        1234567,
			round:         0,
			blockIDHash:   makeBytes(0x11, 32),
			partSetTotal:  1,
			partSetHash:   makeBytes(0x22, 32),
			timestampSec:  1700000000,
			timestampNano: 123456789,
		},
		{
			name:          "zero_timestamp",
			chainID:       "osmosis-1",
			height:        9000000,
			round:         3,
			blockIDHash:   makeBytes(0x33, 32),
			partSetTotal:  100,
			partSetHash:   makeBytes(0x44, 32),
			timestampSec:  0,
			timestampNano: 0,
		},
		{
			name:          "long_chainid",
			chainID:       "my-very-long-chain-identifier-47",
			height:        1,
			round:         0,
			blockIDHash:   makeBytes(0x55, 32),
			partSetTotal:  0xffffffff,
			partSetHash:   makeBytes(0x66, 32),
			timestampSec:  1,
			timestampNano: 1,
		},
	}

	for _, v := range vectors {
		vote := &cmtproto.Vote{
			Type:   cmtproto.PrecommitType,
			Height: v.height,
			Round:  v.round,
			BlockID: cmtproto.BlockID{
				Hash: v.blockIDHash,
				PartSetHeader: cmtproto.PartSetHeader{
					Total: v.partSetTotal,
					Hash:  v.partSetHash,
				},
			},
			Timestamp: time.Unix(v.timestampSec, int64(v.timestampNano)).UTC(),
		}
		signBytes := cmttypes.VoteSignBytes(v.chainID, vote)
		fmt.Printf("// %s\n", v.name)
		fmt.Printf("// chainID=%q height=%d round=%d partSetTotal=%d tsSec=%d tsNano=%d\n",
			v.chainID, v.height, v.round, v.partSetTotal, v.timestampSec, v.timestampNano)
		fmt.Printf("hash=%s partSetHash=%s\n", hex.EncodeToString(v.blockIDHash), hex.EncodeToString(v.partSetHash))
		fmt.Printf("signBytes=%s\n\n", hex.EncodeToString(signBytes))
	}
}
