// Package redact removes API keys from the text of errors that quote an RPC
// endpoint. It is a leaf package with no third-party imports, so the gRPC
// adapter can use it without taking on the Ethereum client.
package redact

import (
	"net/url"
	"regexp"
	"strings"
)

// redactedURL stands in for a URL that cannot be parsed, so there is no telling
// which part of it is a secret.
const redactedURL = "<redacted url>"

// networkURL matches the URL transports the attestor dials. The match stops at
// whitespace and at the quote characters an error message wraps a URL in.
var networkURL = regexp.MustCompile(`(?i)\b(?:https?|wss?)://[^\s"'<>]+`)

// URL reduces a URL to its scheme and host, dropping userinfo, path,
// query and fragment: that is where RPC providers put the API key
// (https://rpc.example/v2/KEY, wss://rpc.example/ws?apikey=KEY). A filesystem
// path (an IPC socket) carries no key and is returned as is.
func URL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return redactedURL
	}
	switch {
	case u.Host != "":
		return u.Scheme + "://" + u.Host
	case u.Scheme != "":
		return u.Scheme + "://"
	default:
		return raw
	}
}

// URLs applies URL to every http, https, ws or wss URL inside s.
// net/http formats a failed request as `Post "<full URL>": <cause>` and strips
// only the password of the userinfo, so a key in the path or query would
// otherwise travel with the error into logs.
func URLs(s string) string {
	return networkURL.ReplaceAllStringFunc(s, func(match string) string {
		// Punctuation that follows a URL in a sentence is not part of it.
		trimmed := strings.TrimRight(match, ".,;:)]}")
		return URL(trimmed) + match[len(trimmed):]
	})
}

// redactedError is an error whose text has its URLs redacted while the
// original stays reachable through Unwrap, so errors.Is and errors.As work.
type redactedError struct {
	text string
	err  error
}

func (e *redactedError) Error() string { return e.text }
func (e *redactedError) Unwrap() error { return e.err }

// Error returns err with the URLs in its text redacted, for an error
// that is about to be logged or returned from startup. It is nil for nil.
func Error(err error) error {
	if err == nil {
		return nil
	}
	return &redactedError{text: URLs(err.Error()), err: err}
}
