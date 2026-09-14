package services

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBuildCosmosClientUpdateRejectsAlreadyCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	worker := &Worker{}
	_, err := worker.BuildCosmosClientUpdateMsg(ctx, CosmosEndpoint{}, EVMEndpoint{}, time.Second, "2/3", "", 0, "2/3", false, 0)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}

func TestHopSearchRejectsAlreadyCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := findHighestFeasibleHop(ctx, cosmosClientDeps{}, 1, 10, "chain-1", pinnedCosmosValidatorSet{}, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}
