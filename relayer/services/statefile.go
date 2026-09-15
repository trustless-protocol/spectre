package services

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	stateFilePerm = 0o600
	stateDirPerm  = 0o700
)

// writeStateFile atomically replaces path with data, syncing both the file and
// its parent directory. Syncing the directory makes the rename itself durable:
// without it, a crash after Rename can lose the replacement name even though
// the temporary file contents were synced.
func writeStateFile(path string, data []byte) error {
	return writeStateFileWithDirSync(path, data, fsyncDir)
}

// writeStateFileWithDirSync keeps the filesystem boundary injectable without a
// package-global test hook. The production path always uses fsyncDir.
func writeStateFileWithDirSync(path string, data []byte, syncDir func(string) error) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, stateDirPerm); err != nil {
		return fmt.Errorf("create state directory %q: %w", dir, err)
	}
	tmp, err := os.CreateTemp(dir, ".state-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary state file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if err := tmp.Chmod(stateFilePerm); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("chmod temporary state file: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temporary state file: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync temporary state file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temporary state file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("replace state file %q: %w", path, err)
	}
	if err := syncDir(dir); err != nil {
		return fmt.Errorf("sync state directory %q: %w", dir, err)
	}
	return nil
}

// fsyncDir durably records a completed state-file rename in its parent
// directory. The relayer targets Unix-like hosts, where directories can be
// opened and synced just like regular files.
func fsyncDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	return d.Sync()
}
