package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestClientCleanupRunsAfterRelayDrain(t *testing.T) {
	var order []string
	stopRelaysAndCleanup(
		func() { order = append(order, "cancel") },
		func() { order = append(order, "drain") },
		[]func(){func() { order = append(order, "cleanup") }},
	)
	want := []string{"cancel", "drain", "cleanup"}
	if fmt.Sprint(order) != fmt.Sprint(want) {
		t.Fatalf("order = %v, want %v", order, want)
	}
}

func TestIsShutdownErr(t *testing.T) {
	for _, tc := range []struct {
		err  error
		want bool
	}{
		{context.Canceled, true},
		{context.DeadlineExceeded, true},
		{fmt.Errorf("wrapped: %w", context.Canceled), true},
		{errors.New("request was cancelled by server"), false},
		{nil, false},
	} {
		if got := isShutdownErr(tc.err); got != tc.want {
			t.Fatalf("isShutdownErr(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}
