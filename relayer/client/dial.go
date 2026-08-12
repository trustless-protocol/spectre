package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	jsonrpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// DefaultRPCTimeout is the transport-level ceiling on a single JSON-RPC round
// trip. It is a LIVENESS BACKSTOP, not a latency target: neither library we
// dial supplies one of its own, so without it a node that accepts the TCP
// connection and then stops answering blocks the caller forever.
//
// That is not hypothetical here. The relay loops drive one direction on a
// single goroutine, so one unanswered request wedges that whole direction with
// no error, no retry and nothing in the log — the first row of the failure
// table in the local-devnet runbook. See GetEthMembershipProof for the same
// hazard caught at a call site.
//
// It is deliberately much larger than Config.FetchTimeout (15s). FetchTimeout
// bounds calls we expect to be fast; this bounds every call, including ones
// that are legitimately slow — a wide eth_getLogs scan, an eth_getProof against
// an archive node, a paginated TxSearch over a long recovery range. Set it low
// enough to be a useful SLO and it starts cancelling healthy work, which is a
// worse failure than the one it prevents. Two minutes is far beyond any healthy
// response and far below "forever".
const DefaultRPCTimeout = 2 * time.Minute

// DialCosmosRPC connects to a CometBFT node with a bounded HTTP client.
//
// The timeout applies ONLY to the HTTP JSON-RPC calls (Status, TxSearch,
// ABCIQuery, Commit, Validators, Block). It does not reach the WebSocket
// subscription: NewWithClient passes the http.Client to the JSON-RPC caller and
// builds WSEvents separately via newWSEvents, which dials its own long-lived
// connection (cometbft rpc/client/http/http.go:132-152). A timeout that killed
// the subscription would silently undo the live-path fix in #377, so this
// separation is load-bearing rather than incidental.
func DialCosmosRPC(remote, wsEndpoint string, timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := jsonrpcclient.DefaultHTTPClient(remote)
	if err != nil {
		return nil, fmt.Errorf("cosmos rpc %s: build http client: %w", remote, err)
	}
	httpClient.Timeout = timeout
	c, err := rpchttp.NewWithClient(remote, wsEndpoint, httpClient)
	if err != nil {
		return nil, fmt.Errorf("cosmos rpc %s: %w", remote, err)
	}
	return c, nil
}

// DialEthRPC connects to an Ethereum-family JSON-RPC endpoint, applying the
// timeout to HTTP(S) endpoints only.
//
// WebSocket endpoints are returned unbounded on purpose. A ws:// connection is
// long-lived by design — an event subscription is meant to stay open for the
// process's lifetime — so an http.Client deadline is both inapplicable (geth
// routes ws/wss through newClientTransportWS, which never reads cfg.httpClient
// — rpc/client.go:206-214) and wrong in intent. Liveness on the WS path is
// already covered: the subscription surfaces failures on sub.Err(), which
// SubscribeEth turns into a reconnect.
func DialEthRPC(ctx context.Context, rawurl string, timeout time.Duration) (*ethclient.Client, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, fmt.Errorf("eth rpc %s: parse url: %w", rawurl, err)
	}

	var opts []rpc.ClientOption
	if u.Scheme == "http" || u.Scheme == "https" {
		opts = append(opts, rpc.WithHTTPClient(&http.Client{Timeout: timeout}))
	}

	rpcClient, err := rpc.DialOptions(ctx, rawurl, opts...)
	if err != nil {
		return nil, fmt.Errorf("eth rpc %s: %w", rawurl, err)
	}
	return ethclient.NewClient(rpcClient), nil
}
