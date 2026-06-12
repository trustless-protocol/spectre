package client

import "context"

type Client interface {
	Status(context.Context) (*ResultStatus, error)
	BlockResults(context.Context, *int64) (*ResultBlockResults, error)
}

type ResultStatus struct {
	SyncInfo SyncInfo `json:"sync_info"`
}

type SyncInfo struct {
	LatestBlockHeight int64 `json:"latest_block_height"`
	CatchingUp        bool  `json:"catching_up"`
}

type ResultBlockResults struct {
	EndBlockEvents []Event `json:"end_block_events"`
}

type Event struct {
	Type       string           `json:"type"`
	Attributes []EventAttribute `json:"attributes,omitempty"`
}

type EventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Index bool   `json:"index,omitempty"`
}
