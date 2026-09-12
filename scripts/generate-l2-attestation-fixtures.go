// Command generate-l2-attestation-fixtures writes deterministic cross-language
// Ed25519 fixtures consumed by the Rust and wasmvm conformance suites.
package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
)

const (
	protocolDomain = "SPECTRE_L2_ATTESTATION_V1"
	chainID        = uint64(11_155_420)
	blockNumber    = uint64(42)
	routerHex      = "645280885749dc97ea461de280eb3273c91d36df"
	blockHashHex   = "1112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f30"
	stateRootHex   = "3132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f50"
)

type vectorFile struct {
	ProtocolDomain string       `json:"protocol_domain"`
	Cases          []vectorCase `json:"cases"`
}

type vectorCase struct {
	Name            string             `json:"name"`
	Threshold       uint16             `json:"threshold"`
	PublicKeys      [][]byte           `json:"public_keys"`
	AttestorSetHash string             `json:"attestor_set_hash"`
	Context         attestationContext `json:"context"`
	Statement       string             `json:"statement"`
	Signatures      []vectorSignature  `json:"signatures"`
}

type attestationContext struct {
	L2ChainID   uint64 `json:"l2_chain_id"`
	Router      string `json:"router_address"`
	BlockNumber uint64 `json:"block_number"`
	BlockHash   string `json:"block_hash"`
	StateRoot   string `json:"state_root"`
}

type vectorSignature struct {
	AttestorIndex uint16 `json:"attestor_index"`
	PublicKey     []byte `json:"public_key"`
	Signature     []byte `json:"signature"`
}

type signer struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func main() {
	output := flag.String(
		"out",
		"test/fixtures/wasm-contracts/l2-attestation-vectors.json",
		"output path, or - for stdout",
	)
	flag.Parse()

	domain := sha256.Sum256([]byte(protocolDomain))
	vectors := vectorFile{
		ProtocolDomain: hex.EncodeToString(domain[:]),
		Cases: []vectorCase{
			buildCase("1-of-1", 1, 1),
			buildCase("2-of-3", 3, 2),
			buildCase("32-of-32", 32, 32),
		},
	}
	encoded, err := json.MarshalIndent(vectors, "", "  ")
	if err != nil {
		panic(err)
	}
	encoded = append(encoded, '\n')
	if *output == "-" {
		if _, err := os.Stdout.Write(encoded); err != nil {
			panic(err)
		}
		return
	}
	if err := os.WriteFile(*output, encoded, 0o644); err != nil {
		panic(err)
	}
}

func buildCase(name string, members int, threshold uint16) vectorCase {
	signers := signerSet(members)
	publicKeys := make([][]byte, len(signers))
	for index, item := range signers {
		publicKeys[index] = append([]byte(nil), item.publicKey...)
	}
	setHash := attestorSetHash(publicKeys, threshold)
	statement := attestationStatement(setHash)
	signatures := make([]vectorSignature, threshold)
	for index := range signatures {
		signatures[index] = vectorSignature{
			AttestorIndex: uint16(index),
			PublicKey:     append([]byte(nil), signers[index].publicKey...),
			Signature:     ed25519.Sign(signers[index].privateKey, statement),
		}
	}
	return vectorCase{
		Name:            name,
		Threshold:       threshold,
		PublicKeys:      publicKeys,
		AttestorSetHash: hex.EncodeToString(setHash[:]),
		Context: attestationContext{
			L2ChainID:   chainID,
			Router:      routerHex,
			BlockNumber: blockNumber,
			BlockHash:   blockHashHex,
			StateRoot:   stateRootHex,
		},
		Statement:  hex.EncodeToString(statement),
		Signatures: signatures,
	}
}

func signerSet(members int) []signer {
	signers := make([]signer, members)
	for member := range signers {
		seed := make([]byte, ed25519.SeedSize)
		for index := range seed {
			seed[index] = byte(member + 1)
		}
		privateKey := ed25519.NewKeyFromSeed(seed)
		signers[member] = signer{
			publicKey:  privateKey.Public().(ed25519.PublicKey),
			privateKey: privateKey,
		}
	}
	sort.Slice(signers, func(left, right int) bool {
		return string(signers[left].publicKey) < string(signers[right].publicKey)
	})
	return signers
}

func attestorSetHash(publicKeys [][]byte, threshold uint16) [32]byte {
	encoded := make([]byte, 4, 4+len(publicKeys)*ed25519.PublicKeySize)
	binary.BigEndian.PutUint16(encoded[0:2], threshold)
	binary.BigEndian.PutUint16(encoded[2:4], uint16(len(publicKeys)))
	for _, publicKey := range publicKeys {
		encoded = append(encoded, publicKey...)
	}
	return sha256.Sum256(encoded)
}

func attestationStatement(setHash [32]byte) []byte {
	domain := sha256.Sum256([]byte(protocolDomain))
	router := mustDecode(routerHex, 20)
	blockHash := mustDecode(blockHashHex, 32)
	stateRoot := mustDecode(stateRootHex, 32)
	statement := make([]byte, 164)
	copy(statement[0:32], domain[:])
	binary.BigEndian.PutUint64(statement[32:40], chainID)
	copy(statement[40:60], router)
	copy(statement[60:92], setHash[:])
	binary.BigEndian.PutUint64(statement[92:100], blockNumber)
	copy(statement[100:132], blockHash)
	copy(statement[132:164], stateRoot)
	return statement
}

func mustDecode(encoded string, expectedLength int) []byte {
	decoded, err := hex.DecodeString(encoded)
	if err != nil || len(decoded) != expectedLength {
		panic(fmt.Sprintf("invalid fixed fixture value %q", encoded))
	}
	return decoded
}
