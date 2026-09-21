package avalanche

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Avalanche warp message primitives, mirrored byte-for-byte with
// packages/avalanche-warp/src/lib.rs (the wasm client's verifying side) and
// pinned to the same live-captured fixture. Only the tiny subset the relayer
// needs is implemented: building the unsigned block-hash message the
// signature-aggregator is asked to sign, and splitting the signed message it
// returns.

// PrimaryNetworkID is the CB58 id of the Avalanche primary network — the
// signing-subnet-id for C-Chain messages verified against the primary
// validator set.
const PrimaryNetworkID = "11111111111111111111111111111111LpoYY"

// warpSignatureLen is the BLS12-381 G2 aggregate signature length.
const warpSignatureLen = 96

// buildBlockHashWarpMessage builds the unsigned warp message for a `Hash`
// payload: `codec(2)=0 ‖ networkID(4,BE) ‖ sourceChainID(32) ‖ payloadLen(4,BE)
// ‖ (codec(2)=0 ‖ typeID(4)=0 ‖ hash(32))`.
func buildBlockHashWarpMessage(networkID uint32, sourceChainID, blockHash [32]byte) []byte {
	message := make([]byte, 0, 42+38)
	message = append(message, 0, 0)
	message = binary.BigEndian.AppendUint32(message, networkID)
	message = append(message, sourceChainID[:]...)
	message = binary.BigEndian.AppendUint32(message, 38)
	message = append(message, 0, 0, 0, 0, 0, 0)
	message = append(message, blockHash[:]...)
	return message
}

// splitSignedWarpMessage splits a signed warp message (`unsigned ‖
// BitSetSignature`) into its parts: the unsigned prefix, the signer bitset,
// and the 96-byte aggregate signature.
func splitSignedWarpMessage(raw []byte) (unsigned, bitSet []byte, signature [warpSignatureLen]byte, err error) {
	if len(raw) < 42 {
		return nil, nil, signature, fmt.Errorf("avalanche: signed warp message shorter than the unsigned envelope")
	}
	if raw[0] != 0 || raw[1] != 0 {
		return nil, nil, signature, fmt.Errorf("avalanche: unsupported warp codec version")
	}
	payloadLen := binary.BigEndian.Uint32(raw[38:42])
	unsignedLen := 42 + int(payloadLen)
	if len(raw) < unsignedLen+8 {
		return nil, nil, signature, fmt.Errorf("avalanche: signed warp message truncated after the unsigned part")
	}
	sig := raw[unsignedLen:]
	if !bytes.Equal(sig[0:4], []byte{0, 0, 0, 0}) {
		return nil, nil, signature, fmt.Errorf("avalanche: unsupported warp signature type")
	}
	bitSetLen := binary.BigEndian.Uint32(sig[4:8])
	expected := 8 + int(bitSetLen) + warpSignatureLen
	if len(sig) != expected {
		return nil, nil, signature, fmt.Errorf("avalanche: warp signature section is %d bytes, want %d", len(sig), expected)
	}
	copy(signature[:], sig[8+bitSetLen:])
	return raw[:unsignedLen], sig[8 : 8+bitSetLen], signature, nil
}
