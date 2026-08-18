package client

import (
	"context"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func silentEthRPC(t *testing.T) *ethclient.Client {
	t.Helper()
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { <-release }))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	client, err := ethclient.Dial(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.Close)
	return client
}

func TestEthProofRPCsHonorCancellation(t *testing.T) {
	for _, tc := range []struct {
		name string
		call func(context.Context, *ethclient.Client) error
	}{
		{"membership", func(ctx context.Context, client *ethclient.Client) error {
			_, err := GetEthMembershipProof(ctx, client, common.Address{}, []byte("path"), common.Hash{}, big.NewInt(1))
			return err
		}},
		{"non-membership", func(ctx context.Context, client *ethclient.Client) error {
			_, err := GetEthNonMembershipProof(ctx, client, common.Address{}, []byte("path"), common.Hash{}, big.NewInt(1))
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			done := make(chan error, 1)
			client := silentEthRPC(t)
			go func() { done <- tc.call(ctx, client) }()
			time.Sleep(30 * time.Millisecond)
			cancel()
			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("error = %v, want context.Canceled", err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("proof RPC ignored cancellation")
			}
		})
	}
}
