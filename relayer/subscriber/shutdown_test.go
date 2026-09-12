package subscriber

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	relayerclient "relayer/client"
	"relayer/services"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/websocket"
)

func TestSleepOrDoneCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	if sleepOrDone(ctx, time.Minute) {
		t.Fatal("cancelled sleep reported elapsed")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("cancelled sleep took %s", elapsed)
	}
}

func TestSleepOrDoneCancellationDuringWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)
	go func() { done <- sleepOrDone(ctx, time.Minute) }()
	cancel()
	select {
	case elapsed := <-done:
		if elapsed {
			t.Fatal("cancelled sleep reported elapsed")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("reconnect sleep ignored cancellation")
	}
}

func TestSubscribeCosmosAlreadyCancelledDoesNotDereferenceClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		NewSubscriber().SubscribeCosmos(ctx, services.CosmosEndpoint{}, services.EVMEndpoint{}, services.ClientIDs{}, nil, nil)
	}()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Cosmos subscriber ignored already-cancelled context")
	}
}

func TestSubscribeEthAlreadyCancelledDoesNotDialOrDereferenceClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		NewSubscriber().SubscribeEth(ctx, services.CosmosEndpoint{}, services.EVMEndpoint{WSURL: "ws://127.0.0.1:1"}, services.ClientIDs{}, nil, nil)
	}()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Ethereum subscriber dialled after cancellation")
	}
}

func TestSubscribeEthWatchHonorsCancellation(t *testing.T) {
	var subscriptions atomic.Int32
	ready := make(chan struct{})
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			var req struct {
				JSONRPC string          `json:"jsonrpc"`
				ID      json.RawMessage `json:"id"`
				Method  string          `json:"method"`
			}
			if err := conn.ReadJSON(&req); err != nil {
				return
			}
			var result any = true
			if req.Method == "eth_subscribe" {
				n := subscriptions.Add(1)
				result = fmt.Sprintf("0x%x", n)
				if n == 4 {
					close(ready)
				}
			}
			if err := conn.WriteJSON(map[string]any{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"result":  result,
			}); err != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)

	watchClient, err := ethclient.Dial("ws" + strings.TrimPrefix(server.URL, "http"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(watchClient.Close)
	router := common.Address{}
	deps := ethDeps{
		EVM:    services.EVMEndpoint{Client: watchClient, Contracts: services.EVMContracts{Router: router}},
		Logger: log.New(io.Discard, "", 0),
	}
	recoveryFilterer, err := contractICS26Router.NewContractICS26RouterFilterer(router, watchClient)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- NewSubscriber().subscribeEthOnce(
			ctx, deps, services.NewBatchBuilder(), watchClient, recoveryFilterer,
			1, new(uint64), new(uint64), nil, new(uint64), &relayerclient.LogSpan{}, &blockRate{},
		)
	}()

	select {
	case <-ready:
	case <-time.After(3 * time.Second):
		t.Fatal("Ethereum watch subscriptions did not start")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("watch error = %v, want context.Canceled", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Ethereum watch ignored cancellation")
	}
}
