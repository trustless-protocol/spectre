package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestResumeCursor(t *testing.T) {
	tests := []struct {
		name                        string
		persisted, lookback, latest uint64
		start                       uint64
		resumed                     bool
		behind                      uint64
	}{
		{name: "missing", lookback: 744, latest: 1000, start: 744},
		{name: "older than lookback", persisted: 100, lookback: 744, latest: 1000, start: 100, resumed: true, behind: 900},
		{name: "at lookback", persisted: 744, lookback: 744, latest: 1000, start: 744, resumed: true, behind: 256},
		{name: "inside lookback", persisted: 900, lookback: 744, latest: 1000, start: 744, resumed: true},
		{name: "ahead of head", persisted: 1100, lookback: 744, latest: 1000, start: 744, resumed: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, resumed, behind := ResumeCursor(tt.persisted, tt.lookback, tt.latest)
			if start != tt.start || resumed != tt.resumed || behind != tt.behind {
				t.Fatalf("ResumeCursor() = (%d,%t,%d), want (%d,%t,%d)", start, resumed, behind, tt.start, tt.resumed, tt.behind)
			}
		})
	}
}

func TestRecoveryStateMissingFileAndPartialCheckpoints(t *testing.T) {
	resetLoadedRecoveryStoresForTest()
	path := filepath.Join(t.TempDir(), "state", "cursors.json")
	store, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := store.Get(); ok {
		t.Fatal("missing file unexpectedly contained a source")
	}
	if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 20, "eth_send": 30}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 10, "eth_write_ack": 40}); err != nil {
		t.Fatal(err)
	}
	got, ok := store.Get()
	if !ok || !got.equal(RecoveryCursors{"cosmos": 10, "eth_send": 30, "eth_write_ack": 40}) {
		t.Fatalf("merged cursors = %+v, ok=%t", got, ok)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != stateFilePerm {
		t.Fatalf("state mode = %o, want %o", info.Mode().Perm(), stateFilePerm)
	}

	resetLoadedRecoveryStoresForTest()
	reloaded, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	got, ok = reloaded.Get()
	if !ok || got["cosmos"] != 10 || got["eth_send"] != 30 || got["eth_write_ack"] != 40 {
		t.Fatalf("reloaded cursors = %+v, ok=%t", got, ok)
	}
}

func TestRecoveryStateCheckpointCanMoveBackward(t *testing.T) {
	resetLoadedRecoveryStoresForTest()
	path := filepath.Join(t.TempDir(), "state.json")
	store, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 100, "eth_send": 200}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 50}); err != nil {
		t.Fatal(err)
	}

	got, ok := store.Get()
	want := RecoveryCursors{"cosmos": 50, "eth_send": 200}
	if !ok || !got.equal(want) {
		t.Fatalf("checkpoint cursors = %+v, ok=%t, want %+v", got, ok, want)
	}

	resetLoadedRecoveryStoresForTest()
	reloaded, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	got, ok = reloaded.Get()
	if !ok || !got.equal(want) {
		t.Fatalf("reloaded checkpoint cursors = %+v, ok=%t, want %+v", got, ok, want)
	}
}

func TestRecoveryStateFailedWriteDoesNotAdvanceMemory(t *testing.T) {
	resetLoadedRecoveryStoresForTest()
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	store, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 10}); err != nil {
		t.Fatal(err)
	}

	blockedPath := filepath.Join(dir, "blocked")
	if err := os.Mkdir(blockedPath, 0o700); err != nil {
		t.Fatal(err)
	}
	store.path = blockedPath
	if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 20}); err == nil {
		t.Fatal("Save succeeded when the target path was a directory")
	}
	got, ok := store.Get()
	if !ok || got["cosmos"] != 10 {
		t.Fatalf("cursor after failed write = %+v, ok=%t, want height 10", got, ok)
	}

	store.path = path
	if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 20}); err != nil {
		t.Fatal(err)
	}
	resetLoadedRecoveryStoresForTest()
	reloaded, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	got, ok = reloaded.Get()
	if !ok || got["cosmos"] != 20 {
		t.Fatalf("reloaded cursor after retry = %+v, ok=%t, want height 20", got, ok)
	}
}

func TestRecoveryStateRejectsCorruptAndUnsupportedFiles(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{name: "corrupt", data: []byte("{")},
		{name: "wrong version", data: []byte(`{"version":1,"sources":{}}`)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resetLoadedRecoveryStoresForTest()
			path := filepath.Join(t.TempDir(), "state.json")
			if err := os.WriteFile(path, tc.data, 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := LoadRecoveryState(path); err == nil {
				t.Fatal("LoadRecoveryState succeeded for invalid state")
			}
		})
	}
}

// One process now keeps one record, but two recovery streams write to it: the
// Cosmos subscriber owns CosmosHeight, the ETH subscriber owns EthSendBlock and
// EthWriteAckBlk. They run concurrently and each save is a read-modify-write of
// the whole file, so a lost update here means a stream silently resumes from an
// older cursor and rescans -- or worse, the other stream's newer cursor is
// rolled back and its window is skipped.
//
// Memoization by path is what makes the two share one view; without it each
// would hold its own copy and the last writer would erase the other's field.
func TestRecoveryStateMemoizationSharesOneViewBetweenStreams(t *testing.T) {
	resetLoadedRecoveryStoresForTest()
	path := filepath.Join(t.TempDir(), "state.json")
	cosmosView, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	ethView, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	if cosmosView != ethView {
		t.Fatal("the same path returned two stores; the two streams would each keep their own cursors")
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := cosmosView.SaveCheckpoint(RecoveryCursors{"cosmos": 11}); err != nil {
			t.Errorf("cosmos save: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := ethView.SaveCheckpoint(RecoveryCursors{"eth_send": 22, "eth_write_ack": 33}); err != nil {
			t.Errorf("eth save: %v", err)
		}
	}()
	wg.Wait()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var file recoveryStateFile
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatal(err)
	}
	want := RecoveryCursors{"cosmos": 11, "eth_send": 22, "eth_write_ack": 33}
	if !file.Cursors.equal(want) {
		t.Fatalf("persisted cursors = %+v, want %+v; one stream's write erased the other's", file.Cursors, want)
	}
}

func TestRecoveryStateWriteLeavesNoTempFiles(t *testing.T) {
	resetLoadedRecoveryStoresForTest()
	dir := t.TempDir()
	store, err := LoadRecoveryState(filepath.Join(dir, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	for i := range 5 {
		if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": uint64(i + 1)}); err != nil {
			t.Fatalf("save %d: %v", i, err)
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tmp") {
			t.Fatalf("temporary recovery file %q was not removed", entry.Name())
		}
	}
}

// A corrupt cursor file must be fatal, never treated as a first run. Silently
// falling back to the lookback skips every event between the lost cursor and
// the lookback window, with nothing in the log to say so -- which is the outage
// the cursor exists to prevent, recreated by the recovery from losing it.
//
// A missing file is the opposite: that IS a first run, and it must load.
func TestLoadRecoveryStateRefusesDamagedState(t *testing.T) {
	write := func(t *testing.T, content string) string {
		t.Helper()
		path := filepath.Join(t.TempDir(), "cursors.json")
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("seed: %v", err)
		}
		return path
	}
	valid := `{"version":2,"cursors":{"cosmos":10,"eth_send":20,"eth_write_ack":30}}`

	t.Run("a truncated file is fatal", func(t *testing.T) {
		for name, cut := range map[string]int{
			"mid object": len(valid) / 2,
			"mid key":    12,
			"first byte": 1,
		} {
			t.Run(name, func(t *testing.T) {
				resetLoadedRecoveryStoresForTest()
				if _, err := LoadRecoveryState(write(t, valid[:cut])); err == nil {
					t.Fatal("a truncated cursor file loaded as if it were a first run")
				}
			})
		}
	})

	t.Run("an unknown version is fatal rather than partially read", func(t *testing.T) {
		resetLoadedRecoveryStoresForTest()
		_, err := LoadRecoveryState(write(t, `{"version":99,"cursors":{"cosmos":10}}`))
		if err == nil {
			t.Fatal("a future version loaded; a field this build does not understand would be dropped silently")
		}
		if !strings.Contains(err.Error(), "version") {
			t.Errorf("error %q does not say the version is the problem", err)
		}
	})

	t.Run("a missing file is a first run, not an error", func(t *testing.T) {
		resetLoadedRecoveryStoresForTest()
		store, err := LoadRecoveryState(filepath.Join(t.TempDir(), "cursors.json"))
		if err != nil {
			t.Fatalf("a first run was rejected: %v", err)
		}
		if _, ok := store.Get(); ok {
			t.Fatal("a first run reported durable cursors")
		}
	})

	// The control: the file this test truncates must load intact, or the cases
	// above would pass against a loader that rejects everything.
	t.Run("an intact file loads", func(t *testing.T) {
		resetLoadedRecoveryStoresForTest()
		store, err := LoadRecoveryState(write(t, valid))
		if err != nil {
			t.Fatalf("intact state was rejected: %v", err)
		}
		got, ok := store.Get()
		if !ok {
			t.Fatal("intact state reported no cursors")
		}
		if want := (RecoveryCursors{"cosmos": 10, "eth_send": 20, "eth_write_ack": 30}); !got.equal(want) {
			t.Fatalf("cursors = %+v, want %+v", got, want)
		}
	})
}

// Reported by @DongLieu on #440. The record changed shape in this PR -- a
// per-source map of chain-named fields became one map keyed by stream -- while
// the file version stayed 1.
//
// The danger is that the old file is not malformed. It unmarshals cleanly into
// the new struct: "sources" is simply an unknown field, "cursors" is absent, so
// Cursors comes back nil and the version check passes. The store would then
// report itself LOADED with no cursors, every stream would miss its key, and
// every recovery scan would quietly restart from the startup lookback -- the
// exact event loss a durable cursor exists to prevent, arriving at the moment of
// upgrade and saying nothing.
//
// The fixture is the real previous format, not a hand-shortened one, because the
// failure depends on it deserializing successfully.
func TestLoadRecoveryStateRefusesThePreviousRecordShape(t *testing.T) {
	resetLoadedRecoveryStoresForTest()

	legacy := `{"version":1,"sources":{"cosmoshub-1":{"cosmos_height":8123,` +
		`"eth_send_block":4477,"eth_write_ack_block":4470}}}`
	path := filepath.Join(t.TempDir(), RecoveryCursorsFileName)
	if err := os.WriteFile(path, []byte(legacy), 0o600); err != nil {
		t.Fatalf("write legacy state: %v", err)
	}

	store, err := LoadRecoveryState(path)
	if err == nil {
		cursors, loaded := store.Get()
		t.Fatalf("a version-1 file was accepted (loaded=%v, cursors=%v); the old cursors are "+
			"unreadable in the new shape, so every stream would silently fall back to the "+
			"startup lookback and lose the progress this file exists to keep", loaded, cursors)
	}
	// The operator has to learn which versions are involved and what to do.
	for _, want := range []string{"version 1", "want 2", "delete it"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not mention %q", err, want)
		}
	}
}
