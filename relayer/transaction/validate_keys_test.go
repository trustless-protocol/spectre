package transaction

import (
	"crypto/ecdsa"
	"errors"
	"testing"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type fakeSigner struct {
	ethCalls    int
	cosmosCalls int
	err         error
}

func (s *fakeSigner) EthKey() (*ecdsa.PrivateKey, error) {
	s.ethCalls++
	return nil, s.err
}

func (s *fakeSigner) CosmosKeyBytes() ([]byte, error) {
	s.cosmosCalls++
	return nil, s.err
}

// Startup validation must go THROUGH the seam.
//
// It used to read ETH_PRIVATE_KEY and COSMOS_PRIVATE_KEY itself and parse them
// with keys.RestoreKey, which left a ninth key source outside the seam that
// RLY-07 introduced -- and one a KMS-backed signer could never satisfy, since it
// does not hand back a raw private key to parse. An installed signer must be the
// only thing consulted, with the environment untouched.
func TestValidateKeysUsesTheInstalledSignerNotTheEnvironment(t *testing.T) {
	t.Setenv("ETH_PRIVATE_KEY", "")
	t.Setenv("COSMOS_PRIVATE_KEY", "")

	want := errors.New("kms unavailable")
	fake := &fakeSigner{err: want}
	h := &Handler{}
	h.SetSigner(fake)

	if err := h.ValidateKeys(); !errors.Is(err, want) {
		t.Fatalf("ValidateKeys() = %v, want the signer's own error %v; a validator that reads the "+
			"environment cannot check a KMS-backed key at all", err, want)
	}
	if fake.ethCalls == 0 {
		t.Fatal("the ETH key was never requested from the signer")
	}
}

// Both keys are checked, not just the first: a relayer that starts with a broken
// Cosmos key fails on its first send instead of at startup, which is the whole
// point of validating.
func TestValidateKeysChecksBothKeys(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	fake := &okEthBadCosmosSigner{eth: priv}
	h := &Handler{}
	h.SetSigner(fake)

	if err := h.ValidateKeys(); err == nil {
		t.Fatal("ValidateKeys() = nil with an unusable Cosmos key; the failure would surface on the first send")
	}
	if !fake.cosmosAsked {
		t.Fatal("the Cosmos key was never requested")
	}
}

type okEthBadCosmosSigner struct {
	eth         *ecdsa.PrivateKey
	cosmosAsked bool
}

func (s *okEthBadCosmosSigner) EthKey() (*ecdsa.PrivateKey, error) { return s.eth, nil }

func (s *okEthBadCosmosSigner) CosmosKeyBytes() ([]byte, error) {
	s.cosmosAsked = true
	return nil, errors.New("COSMOS_PRIVATE_KEY environment variable is required in .env file")
}
