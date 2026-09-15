package attestation

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSignerBindsEveryWasmHeaderIdentityField(t *testing.T) {
	signer, err := NewSigner(8453, "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20")
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	router := make([]byte, RouterLength)
	setHash := make([]byte, HashLength)
	blockHash := make([]byte, HashLength)
	stateRoot := make([]byte, HashLength)
	router[19] = 1
	setHash[31] = 2
	blockHash[31] = 3
	stateRoot[31] = 1
	signature, err := signer.Sign(router, setHash, 123, blockHash, stateRoot)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	const expectedPublicKey = "79b5562e8fe654f94078b112e8a98ba7901f853ae695bed7e0e3910bad049664"
	const expectedSignature = "bbf5d84390b7316391593425428d3e89c44f930795482bb5afaf4aa4b187c13436c0accbbe1ad323b99c39ed67b57894700a44553f62b46afcd3f4c8c882c105"
	if got := hex.EncodeToString(signer.PublicKey()); got != expectedPublicKey {
		t.Fatalf("public key = %s, want %s", got, expectedPublicKey)
	}
	if got := hex.EncodeToString(signature); got != expectedSignature {
		t.Fatalf("signature = %s, want %s", got, expectedSignature)
	}
	message, err := SigningBytes(8453, router, setHash, 123, blockHash, stateRoot)
	if err != nil {
		t.Fatalf("SigningBytes: %v", err)
	}
	if !ed25519.Verify(signer.PublicKey(), message, signature) {
		t.Fatal("valid block identity signature did not verify")
	}
	for _, bad := range [][]byte{
		mustSigningBytes(t, 8454, router, setHash, 123, blockHash, stateRoot),
		mustSigningBytes(t, 8453, append([]byte(nil), router[:19]...), setHash, 123, blockHash, stateRoot),
		mustSigningBytes(t, 8453, router, append([]byte(nil), stateRoot...), 123, blockHash, stateRoot),
		mustSigningBytes(t, 8453, router, setHash, 124, blockHash, stateRoot),
		mustSigningBytes(t, 8453, router, setHash, 123, stateRoot, stateRoot),
	} {
		if ed25519.Verify(signer.PublicKey(), bad, signature) {
			t.Fatal("signature verified for a different L2 block identity")
		}
	}
}

func mustSigningBytes(t *testing.T, chainID uint64, router, setHash []byte, blockNumber uint64, blockHash, stateRoot []byte) []byte {
	t.Helper()
	message, err := SigningBytes(chainID, router, setHash, blockNumber, blockHash, stateRoot)
	if err != nil {
		return []byte("different-invalid-statement")
	}
	return message
}

func TestWasmAttestationVectors(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "test", "fixtures", "wasm-contracts", "l2-attestation-vectors.json"))
	if err != nil {
		t.Fatal(err)
	}
	var vectors struct {
		ProtocolDomain string `json:"protocol_domain"`
		Cases          []struct {
			Name            string   `json:"name"`
			Threshold       uint16   `json:"threshold"`
			PublicKeys      []string `json:"public_keys"`
			AttestorSetHash string   `json:"attestor_set_hash"`
			Statement       string   `json:"statement"`
			Context         struct {
				L2ChainID     uint64 `json:"l2_chain_id"`
				RouterAddress string `json:"router_address"`
				BlockNumber   uint64 `json:"block_number"`
				BlockHash     string `json:"block_hash"`
				StateRoot     string `json:"state_root"`
			} `json:"context"`
			Signatures []struct {
				AttestorIndex uint16 `json:"attestor_index"`
				PublicKey     string `json:"public_key"`
				Signature     string `json:"signature"`
			} `json:"signatures"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatal(err)
	}
	if got := hex.EncodeToString(ProtocolDomain[:]); got != vectors.ProtocolDomain {
		t.Fatalf("protocol domain = %s, want %s", got, vectors.ProtocolDomain)
	}
	for _, vector := range vectors.Cases {
		t.Run(vector.Name, func(t *testing.T) {
			keys := make([][]byte, len(vector.PublicKeys))
			for i, encoded := range vector.PublicKeys {
				keys[i], err = base64.StdEncoding.DecodeString(encoded)
				if err != nil {
					t.Fatal(err)
				}
			}
			config := AttestorConfig{PublicKeys: keys, Threshold: vector.Threshold}
			setHash, err := config.SetHash()
			if err != nil {
				t.Fatal(err)
			}
			if got := hex.EncodeToString(setHash[:]); got != vector.AttestorSetHash {
				t.Fatalf("set hash = %s, want %s", got, vector.AttestorSetHash)
			}
			router, _ := hex.DecodeString(vector.Context.RouterAddress)
			blockHash, _ := hex.DecodeString(vector.Context.BlockHash)
			stateRoot, _ := hex.DecodeString(vector.Context.StateRoot)
			statement, err := SigningBytes(vector.Context.L2ChainID, router, setHash[:], vector.Context.BlockNumber, blockHash, stateRoot)
			if err != nil {
				t.Fatal(err)
			}
			want, _ := hex.DecodeString(vector.Statement)
			if !bytes.Equal(statement, want) {
				t.Fatalf("statement = %x, want %x", statement, want)
			}
			for _, signed := range vector.Signatures {
				key, _ := base64.StdEncoding.DecodeString(signed.PublicKey)
				signature, _ := base64.StdEncoding.DecodeString(signed.Signature)
				if !ed25519.Verify(key, statement, signature) {
					t.Fatalf("signature for index %d did not verify", signed.AttestorIndex)
				}
			}
		})
	}
}

func TestNewSignerValidatesKeyMaterialAndEnvironment(t *testing.T) {
	const seed = "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	t.Setenv("ATTESTATION_TEST_SEED", seed)
	t.Setenv("ATTESTATION_TEST_MISSING", "")
	if _, err := NewSigner(10, "env:ATTESTATION_TEST_SEED"); err != nil {
		t.Fatalf("NewSigner from environment: %v", err)
	}

	for _, tc := range []struct {
		name    string
		chainID uint64
		key     string
		want    string
	}{
		{name: "zero chain ID", chainID: 0, key: seed, want: "chain ID"},
		{name: "invalid hex", chainID: 10, key: "not-hex", want: "decode"},
		{name: "wrong seed length", chainID: 10, key: "01", want: "32-byte"},
		{name: "missing environment key", chainID: 10, key: "env:ATTESTATION_TEST_MISSING", want: "unset"},
		{name: "empty environment variable name", chainID: 10, key: "env:", want: "empty"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewSigner(tc.chainID, tc.key)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("NewSigner(%d, %q) error = %v, want %q", tc.chainID, tc.key, err, tc.want)
			}
		})
	}
}

func TestSigningBytesRejectsInvalidIdentity(t *testing.T) {
	validRouter := make([]byte, RouterLength)
	validHash := make([]byte, HashLength)
	for _, tc := range []struct {
		name      string
		chainID   uint64
		router    []byte
		setHash   []byte
		blockHash []byte
		stateRoot []byte
	}{
		{name: "zero chain ID", router: validRouter, setHash: validHash, blockHash: validHash, stateRoot: validHash},
		{name: "short router", chainID: 10, router: validRouter[:RouterLength-1], setHash: validHash, blockHash: validHash, stateRoot: validHash},
		{name: "short set hash", chainID: 10, router: validRouter, setHash: validHash[:HashLength-1], blockHash: validHash, stateRoot: validHash},
		{name: "short block hash", chainID: 10, router: validRouter, setHash: validHash, blockHash: validHash[:HashLength-1], stateRoot: validHash},
		{name: "short state root", chainID: 10, router: validRouter, setHash: validHash, blockHash: validHash, stateRoot: validHash[:HashLength-1]},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := SigningBytes(tc.chainID, tc.router, tc.setHash, 1, tc.blockHash, tc.stateRoot); err == nil {
				t.Fatal("SigningBytes succeeded for invalid identity")
			}
		})
	}
}
