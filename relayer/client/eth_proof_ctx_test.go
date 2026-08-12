package client

import (
	"context"
	"math/big"
	"testing"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// The proof helpers are called from goroutines that drive one relay direction
// each, so an eth_getProof that never comes back used to wedge that direction
// with no error and nothing in the log. These tests pin that the caller's
// context can now interrupt the call.
//
// Both run against a server that accepts the request and never answers — the
// shape of the incident, rather than a connection failure, which would prove
// nothing about context handling.

func TestGetEthMembershipProofHonoursContext(t *testing.T) {
	srv := blackHoleServer(t)
	c, err := DialEthRPC(context.Background(), srv.URL, time.Hour)
	if err != nil {
		t.Fatalf("DialEthRPC: %v", err)
	}
	t.Cleanup(c.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		_, err := GetEthMembershipProof(ctx, c, ethcommon.Address{}, []byte("path"), ethcommon.Hash{}, big.NewInt(1))
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected an error from the cancelled call")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("GetEthMembershipProof ignored its context: the relay direction would be wedged")
	}
}

func TestGetEthNonMembershipProofHonoursContext(t *testing.T) {
	srv := blackHoleServer(t)
	c, err := DialEthRPC(context.Background(), srv.URL, time.Hour)
	if err != nil {
		t.Fatalf("DialEthRPC: %v", err)
	}
	t.Cleanup(c.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		_, err := GetEthNonMembershipProof(ctx, c, ethcommon.Address{}, []byte("path"), ethcommon.Hash{}, big.NewInt(1))
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected an error from the cancelled call")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("GetEthNonMembershipProof ignored its context: the timeout scanner would stop for good")
	}
}
