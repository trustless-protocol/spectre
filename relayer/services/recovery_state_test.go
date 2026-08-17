package services

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func TestRecoveryStateMissingFileAndMonotonicPartialSaves(t *testing.T) {
	resetLoadedRecoveryStoresForTest()
	path := filepath.Join(t.TempDir(), "state", "cursors.json")
	store, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := store.Get("source-a"); ok {
		t.Fatal("missing file unexpectedly contained a source")
	}
	if err := store.Save("source-a", RecoveryCursors{CosmosHeight: 20, EthSendBlock: 30}); err != nil {
		t.Fatal(err)
	}
	if err := store.Save("source-a", RecoveryCursors{CosmosHeight: 10, EthWriteAckBlk: 40}); err != nil {
		t.Fatal(err)
	}
	got, ok := store.Get("source-a")
	if !ok || got != (RecoveryCursors{CosmosHeight: 20, EthSendBlock: 30, EthWriteAckBlk: 40}) {
		t.Fatalf("merged cursors = %+v, ok=%t", got, ok)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != recoveryStateFilePerm {
		t.Fatalf("state mode = %o, want %o", info.Mode().Perm(), recoveryStateFilePerm)
	}

	resetLoadedRecoveryStoresForTest()
	reloaded, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	got, ok = reloaded.Get("source-a")
	if !ok || got.CosmosHeight != 20 || got.EthSendBlock != 30 || got.EthWriteAckBlk != 40 {
		t.Fatalf("reloaded cursors = %+v, ok=%t", got, ok)
	}
}

func TestRecoveryStateRejectsCorruptAndUnsupportedFiles(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{name: "corrupt", data: []byte("{")},
		{name: "wrong version", data: []byte(`{"version":2,"sources":{}}`)},
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

func TestRecoveryStateMemoizationPreservesConcurrentSources(t *testing.T) {
	resetLoadedRecoveryStoresForTest()
	path := filepath.Join(t.TempDir(), "state.json")
	a, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	b, err := LoadRecoveryState(path)
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Fatal("same path returned different stores")
	}

	var wg sync.WaitGroup
	for i, id := range []string{"source-a", "source-b", "source-c"} {
		wg.Add(1)
		go func(id string, cursor uint64) {
			defer wg.Done()
			if err := a.Save(id, RecoveryCursors{CosmosHeight: cursor}); err != nil {
				t.Errorf("save %s: %v", id, err)
			}
		}(id, uint64(i+1))
	}
	wg.Wait()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var file recoveryStateFile
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatal(err)
	}
	if len(file.Sources) != 3 {
		t.Fatalf("sources = %d, want 3", len(file.Sources))
	}
}
