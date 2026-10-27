package redact

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"testing"
)

func TestURLKeepsOnlySchemeAndHost(t *testing.T) {
	for _, tc := range []struct {
		name string
		raw  string
		want string
	}{
		{"key in the path", "https://eth-mainnet.example/v2/SECRET-KEY", "https://eth-mainnet.example"},
		{"key in the query", "http://127.0.0.1:8545/rpc?apikey=SECRET-KEY", "http://127.0.0.1:8545"},
		{"key in userinfo", "https://user:SECRET-KEY@rpc.example:8443/", "https://rpc.example:8443"},
		{"websocket with fragment", "wss://rpc.example/ws/SECRET-KEY#frag", "wss://rpc.example"},
		{"ipv6 host", "ws://[::1]:8546/SECRET-KEY", "ws://[::1]:8546"},
		{"nothing to hide", "https://rpc.example", "https://rpc.example"},
		{"ipc socket path is not a secret", "/var/run/geth.ipc", "/var/run/geth.ipc"},
		{"unparsable", "http://[::1/SECRET-KEY", "<redacted url>"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := URL(tc.raw)
			if got != tc.want {
				t.Fatalf("URL(%q) = %q, want %q", tc.raw, got, tc.want)
			}
			if strings.Contains(got, "SECRET-KEY") {
				t.Fatalf("URL(%q) = %q still carries the key", tc.raw, got)
			}
		})
	}
}

func TestURLsRewritesURLsInsideMessages(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   string
		want string
	}{
		{
			"net/http transport failure",
			`Post "http://127.0.0.1:9545/rpc?apikey=SECRET-KEY": EOF`,
			`Post "http://127.0.0.1:9545": EOF`,
		},
		{
			"url followed by a colon",
			"eth rpc https://rpc.example/v2/SECRET-KEY: parse url failed",
			"eth rpc https://rpc.example: parse url failed",
		},
		{
			"two urls",
			"dial wss://a.example/SECRET-KEY, fell back to https://b.example/SECRET-KEY.",
			"dial wss://a.example, fell back to https://b.example.",
		},
		{
			"upper case scheme",
			"HTTPS://RPC.example/SECRET-KEY failed",
			"https://RPC.example failed",
		},
		{
			"unparsable url",
			`parse "http://[::1/SECRET-KEY": missing ']'`,
			`parse "<redacted url>": missing ']'`,
		},
		{"no url", "context deadline exceeded", "context deadline exceeded"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := URLs(tc.in)
			if got != tc.want {
				t.Fatalf("URLs(%q) = %q, want %q", tc.in, got, tc.want)
			}
			if strings.Contains(got, "SECRET-KEY") {
				t.Fatalf("URLs(%q) = %q still carries the key", tc.in, got)
			}
		})
	}
}

// The failure the audit reproduced: net/http strips only the password of the
// userinfo, so a key in the path or query survives into the error text.
func TestErrorHidesTheKeyOfARealTransportError(t *testing.T) {
	transport := &url.Error{Op: "Post", URL: "http://127.0.0.1:9545/rpc?apikey=SECRET-KEY", Err: io.EOF}
	wrapped := fmt.Errorf("optimism_syncStatus failed: %w", transport)
	if !strings.Contains(wrapped.Error(), "SECRET-KEY") {
		t.Fatalf("premise: the unredacted transport error %q should carry the key", wrapped)
	}

	redacted := Error(wrapped)
	if strings.Contains(redacted.Error(), "SECRET-KEY") {
		t.Fatalf("redacted error %q still carries the key", redacted)
	}
	if want := `optimism_syncStatus failed: Post "http://127.0.0.1:9545": EOF`; redacted.Error() != want {
		t.Fatalf("redacted error = %q, want %q", redacted, want)
	}
	// Redaction changes the text, not the identity of the error.
	if !errors.Is(redacted, io.EOF) {
		t.Fatal("Error broke errors.Is")
	}
	var asURL *url.Error
	if !errors.As(redacted, &asURL) || asURL != transport {
		t.Fatal("Error broke errors.As")
	}
}

func TestErrorOfNilIsNil(t *testing.T) {
	if err := Error(nil); err != nil {
		t.Fatalf("Error(nil) = %v", err)
	}
}
