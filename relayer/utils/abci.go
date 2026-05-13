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
	return IbcPath(packet.SourceClient, packet.Sequence, appendByte)
}

func IbcPath(clientID string, sequence uint64, appendByte []byte) [][]byte {
	sequenceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sequenceBytes, sequence)
	path := []byte(clientID)
	path = append(path, appendByte...)
	path = append(path, sequenceBytes...)

	return [][]byte{[]byte("ibc"), path}
}
