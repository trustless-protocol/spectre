package subscriber

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/services"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func silentRPCServer(t *testing.T) (*httptest.Server, <-chan struct{}) {
	t.Helper()
	started := make(chan struct{})
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		select {
		case <-started:
		default:
			close(started)
		}
		<-release
	}))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	return server, started
}

func TestCosmosTxSearchHonorsCancellation(t *testing.T) {
	server, started := silentRPCServer(t)
	client, err := rpchttp.New(server.URL, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	deps := cosmosDeps{Cosmos: services.CosmosEndpoint{Client: client}, Logger: log.New(io.Discard, "", 0)}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := recoverCosmosEventsForQuery(ctx, deps, services.NewBatchBuilder(), cometBFTSendPacketTxSearch, 1, 2, nil)
		done <- err
	}()
	<-started
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("TxSearch ignored cancellation")
	}
}

func TestEthFilterLogsHonorsCancellation(t *testing.T) {
	server, started := silentRPCServer(t)
	client, err := ethclient.Dial(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.Close)
	filterer, err := contractICS26Router.NewContractICS26RouterFilterer(common.Address{}, client)
	if err != nil {
		t.Fatal(err)
	}
	deps := ethDeps{EVM: services.EVMEndpoint{Client: client}, Logger: log.New(io.Discard, "", 0)}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := recoverEthSendPackets(ctx, deps, services.NewBatchBuilder(), filterer, 1, 2, nil)
		done <- err
	}()
	<-started
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("FilterLogs ignored cancellation")
	}
}
