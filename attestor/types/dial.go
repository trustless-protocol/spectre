package types

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// DefaultRPCTimeout is the transport-level ceiling for one HTTP JSON-RPC
// request. It is a liveness backstop: a reachable node can accept a request
// and never answer, and go-ethereum does not apply a client timeout itself.
const DefaultRPCTimeout = 2 * time.Minute

// DialEthRPC connects to an Ethereum-family JSON-RPC endpoint. HTTP(S)
// endpoints receive a transport timeout; WebSocket endpoints do not, because
// their connections are intentionally long-lived and do not use http.Client.
func DialEthRPC(ctx context.Context, rawURL string, timeout time.Duration) (*ethclient.Client, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("eth rpc %s: parse url: %w", rawURL, err)
	}

	var opts []rpc.ClientOption
	if u.Scheme == "http" || u.Scheme == "https" {
		opts = append(opts, rpc.WithHTTPClient(&http.Client{Timeout: timeout}))
	}

	rpcClient, err := rpc.DialOptions(ctx, rawURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("eth rpc %s: %w", rawURL, err)
	}
	return ethclient.NewClient(rpcClient), nil
}
