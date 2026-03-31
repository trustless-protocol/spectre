package utils

import (
	"encoding/binary"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func BytesToBytes32(data []byte) [32]byte {
	var result [32]byte
	copy(result[:], data)
	return result
}

func IbcCommitmentPath(packet channeltypesv2.Packet, appendByte []byte) [][]byte {
	sequenceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sequenceBytes, packet.Sequence)
	path := []byte(packet.SourceClient)
	path = append(path, appendByte...)
	path = append(path, sequenceBytes...)

	return [][]byte{[]byte("ibc"), path}
}
