package http

import (
	"context"
	"errors"
	"net/http"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

var errNamadaShim = errors.New("tendermint RPC shim is only for compiling unused interchaintest Namada support")

type HTTP struct{}

func NewWithClient(string, string, *http.Client) (*HTTP, error) {
	return &HTTP{}, nil
}

func (h *HTTP) Status(context.Context) (*rpcclient.ResultStatus, error) {
	return nil, errNamadaShim
}

func (h *HTTP) BlockResults(context.Context, *int64) (*rpcclient.ResultBlockResults, error) {
	return nil, errNamadaShim
}
