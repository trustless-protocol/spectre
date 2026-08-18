package services

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestHasEthIBCPathValueHonorsCancellation(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { <-release }))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	client, err := ethclient.Dial(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.Close)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := HasEthIBCPathValue(ctx, EVMEndpoint{Client: client, Contracts: EVMContracts{Router: common.HexToAddress("0x1")}}, []byte("path"))
		done <- err
	}()
	time.Sleep(30 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("storage query ignored cancellation")
	}
}

func TestScanForCosmosTimeoutsHonorsCancellation(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { <-release }))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	client, err := ethclient.Dial(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.Close)

	svc := New(nil, nil, DefaultConfig())
	svc.TrackCosmosPending(channeltypesv2.Packet{
		SourceClient:      "07-tendermint-0",
		DestinationClient: "client-0",
		Sequence:          1,
		TimeoutTimestamp:  1,
	}, 100)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		svc.scanForCosmosTimeouts(ctx, CosmosEndpoint{}, EVMEndpoint{
			Client:    client,
			Contracts: EVMContracts{Router: common.HexToAddress("0x1")},
		}, "08-wasm-0")
		close(done)
	}()
	time.Sleep(30 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Cosmos timeout scan ignored cancellation")
	}
}
