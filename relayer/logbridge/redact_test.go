package logbridge

import (
	"errors"
	"log"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const testAPIKey = "s3cr3t-api-key"

// #575: RPC providers put the API key in the URL, and URLs reach the log from
// dial lines, config errors and the text of RPC errors.
func TestRedactURLs_ZapMessageAndFields(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	logger := zap.New(RedactURLs(core)).With(zap.String("endpoint", "wss://ws.example.com/v2/"+testAPIKey))

	logger.Info("create-clients: dialing ethereum rpc https://rpc.example.com/v2/"+testAPIKey,
		zap.Error(errors.New(`Post "https://rpc.example.com/v2/`+testAPIKey+`": dial tcp: refused`)),
		zap.Int("attempt", 3))

	entries := logs.All()
	if len(entries) != 1 {
		t.Fatalf("logged %d entries, want 1", len(entries))
	}
	entry := entries[0]
	if strings.Contains(entry.Message, testAPIKey) || !strings.Contains(entry.Message, "https://rpc.example.com") {
		t.Fatalf("message = %q, want the URL reduced to scheme and host", entry.Message)
	}
	fields := entry.ContextMap()
	for key, value := range fields {
		if s, ok := value.(string); ok && strings.Contains(s, testAPIKey) {
			t.Fatalf("field %s = %q carries the API key", key, s)
		}
	}
	if fields["attempt"] != int64(3) {
		t.Fatalf("an unrelated field changed: %v", fields["attempt"])
	}
}

func TestRedactURLs_StandardLoggerLines(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	restore := Install(zap.New(RedactURLs(core)))
	t.Cleanup(restore)

	log.Printf("WARN [eth->cosmos Subscribe] Failed to connect to Ethereum WS at wss://ws.example.com/ws?apikey=%s: refused", testAPIKey)

	entries := logs.All()
	if len(entries) != 1 {
		t.Fatalf("logged %d entries, want 1", len(entries))
	}
	if got := entries[0]; strings.Contains(got.Message, testAPIKey) || got.Level != zapcore.WarnLevel ||
		!strings.Contains(got.Message, "wss://ws.example.com") {
		t.Fatalf("entry = %+v, want a warn line with the URL reduced to scheme and host", got)
	}
}

// The bridge redacts on its own too, for a logger built without the core.
func TestBridgeWriterRedactsWithoutTheCore(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	restore := Install(zap.New(core))
	t.Cleanup(restore)

	log.Printf("[CreateEthClient] starting: beacon=https://beacon.example.com/%s checksum=ab", testAPIKey)
	if got := logs.All()[0].Message; strings.Contains(got, testAPIKey) {
		t.Fatalf("message = %q carries the API key", got)
	}
}
