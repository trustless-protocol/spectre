package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	recoveryStateFileEnv     = "RELAYER_RECOVERY_STATE_FILE"
	defaultRecoveryStateFile = ".relayer-state/recovery-cursors.json"
	recoveryStateFilePerm    = 0o600
	recoveryStateDirPerm     = 0o700
	recoveryStateVersion     = 1
)

// RecoveryCursors stores the next height/block each recovery stream must scan.
// Zero means that stream has no durable cursor yet and must use its lookback.
type RecoveryCursors struct {
	CosmosHeight   uint64 `json:"cosmos_height"`
	EthSendBlock   uint64 `json:"eth_send_block"`
	EthWriteAckBlk uint64 `json:"eth_write_ack_block"`
}

type recoveryStateFile struct {
	Version uint32                     `json:"version"`
	Sources map[string]RecoveryCursors `json:"sources"`
}

// RecoveryStateStore persists per-source recovery cursors in one small JSON
// file. One process shares one store instance per path so multi-source saves do
// not overwrite each other's in-memory view.
type RecoveryStateStore struct {
	path string

	mtx     sync.Mutex
	sources map[string]RecoveryCursors
}

// RecoveryStatePath returns the configured recovery-state path.
func RecoveryStatePath() string {
	if path := os.Getenv(recoveryStateFileEnv); path != "" {
		return path
	}
	return defaultRecoveryStateFile
}

var loadedRecoveryStores = struct {
	mtx    sync.Mutex
	byPath map[string]*RecoveryStateStore
}{byPath: map[string]*RecoveryStateStore{}}

// LoadRecoveryState loads and memoizes a recovery store. A missing file is a
// first run; unreadable, corrupt, or unsupported state is fatal to the caller.
func LoadRecoveryState(path string) (*RecoveryStateStore, error) {
	loadedRecoveryStores.mtx.Lock()
	defer loadedRecoveryStores.mtx.Unlock()
	if existing, ok := loadedRecoveryStores.byPath[path]; ok {
		return existing, nil
	}
	store, err := loadRecoveryStateFile(path)
	if err != nil {
		return nil, err
	}
	loadedRecoveryStores.byPath[path] = store
	return store, nil
}

func resetLoadedRecoveryStoresForTest() {
	loadedRecoveryStores.mtx.Lock()
	defer loadedRecoveryStores.mtx.Unlock()
	loadedRecoveryStores.byPath = map[string]*RecoveryStateStore{}
}

func loadRecoveryStateFile(path string) (*RecoveryStateStore, error) {
	store := &RecoveryStateStore{path: path, sources: map[string]RecoveryCursors{}}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read recovery state %q: %w", path, err)
	}

	var file recoveryStateFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse recovery state %q (delete it to resume from the lookback window, accepting that older events may be lost): %w", path, err)
	}
	if file.Version != recoveryStateVersion {
		return nil, fmt.Errorf("recovery state %q has version %d, want %d", path, file.Version, recoveryStateVersion)
	}
	if file.Sources != nil {
		store.sources = file.Sources
	}
	return store, nil
}

// Get returns one source's cursors.
func (s *RecoveryStateStore) Get(sourceID string) (RecoveryCursors, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	cursors, ok := s.sources[sourceID]
	return cursors, ok
}

// Save merges non-zero fields monotonically and atomically rewrites the file.
func (s *RecoveryStateStore) Save(sourceID string, cursors RecoveryCursors) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	previous := s.sources[sourceID]
	merged := RecoveryCursors{
		CosmosHeight:   maxRecoveryCursor(previous.CosmosHeight, cursors.CosmosHeight),
		EthSendBlock:   maxRecoveryCursor(previous.EthSendBlock, cursors.EthSendBlock),
		EthWriteAckBlk: maxRecoveryCursor(previous.EthWriteAckBlk, cursors.EthWriteAckBlk),
	}
	if merged == previous {
		return nil
	}
	s.sources[sourceID] = merged

	snapshot := recoveryStateFile{Version: recoveryStateVersion, Sources: make(map[string]RecoveryCursors, len(s.sources))}
	for id, value := range s.sources {
		snapshot.Sources[id] = value
	}
	return writeRecoveryState(s.path, snapshot)
}

func maxRecoveryCursor(a, b uint64) uint64 {
	if b > a {
		return b
	}
	return a
}

func writeRecoveryState(path string, file recoveryStateFile) error {
	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return fmt.Errorf("encode recovery state: %w", err)
	}
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, recoveryStateDirPerm); err != nil {
			return fmt.Errorf("create recovery state dir %q: %w", dir, err)
		}
	}
	tmp, err := os.CreateTemp(dir, ".recovery-cursors-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp recovery state: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temp recovery state: %w", err)
	}
	if err := tmp.Chmod(recoveryStateFilePerm); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("chmod temp recovery state: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync temp recovery state: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp recovery state: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replace recovery state %q: %w", path, err)
	}
	return nil
}

// ResumeCursor never lets a persisted cursor shorten the normal lookback. A
// cursor older than the lookback extends recovery backwards; a newer cursor
// still rescans the lookback because scanned events may only have reached RAM.
func ResumeCursor(persisted, lookbackStart, latest uint64) (start uint64, resumed bool, behind uint64) {
	if persisted == 0 {
		return lookbackStart, false, 0
	}
	if persisted > lookbackStart {
		return lookbackStart, true, 0
	}
	if latest > persisted {
		behind = latest - persisted
	}
	return persisted, true, behind
}
