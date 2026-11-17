package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Reported by @DongLieu on #467, and it is the same defect as the one this file
// already fixed once, one level further out. The first version keyed the lock off
// the recovery state directory, which the runbook tells operators to vary. The
// second used os.UserCacheDir, which reads XDG_CACHE_HOME on Unix and HOME on
// macOS -- and start loads .env at start.go:39 before taking the lock at :91, so
// a variable set there (or simply two units running under different HOME) hands
// the two processes different lock files and lets both start.
//
// The lock path must therefore be a function of the user and the machine and
// nothing else. This drives the REAL resolver, not the test seam, because the
// property under test is which inputs it reads.
func TestSignerLockPathFor_IgnoresTheEnvironment(t *testing.T) {
	if os.Getuid() < 0 {
		t.Skip("uid-based path is not used on this platform")
	}
	baseline, err := signerLockDir()
	if err != nil {
		t.Fatalf("signerLockDir: %v", err)
	}

	// Every variable any stdlib directory helper consults.
	for _, envVar := range []string{"HOME", "XDG_CACHE_HOME", "TMPDIR", "XDG_RUNTIME_DIR"} {
		t.Run(envVar, func(t *testing.T) {
			t.Setenv(envVar, t.TempDir())
			got, err := signerLockDir()
			if err != nil {
				t.Fatalf("signerLockDir with %s set: %v", envVar, err)
			}
			if got != baseline {
				t.Fatalf("%s moved the lock directory: %q vs %q. Two processes with different "+
					"%s take different lock files and both start, which is the account-sequence "+
					"race this lock exists to prevent", envVar, got, baseline, envVar)
			}
		})
	}

	// And the whole set at once, which is what a .env file actually does.
	t.Run("all of them at once, as .env would", func(t *testing.T) {
		for _, envVar := range []string{"HOME", "XDG_CACHE_HOME", "TMPDIR", "XDG_RUNTIME_DIR"} {
			t.Setenv(envVar, t.TempDir())
		}
		got, err := signerLockDir()
		if err != nil {
			t.Fatalf("signerLockDir: %v", err)
		}
		if got != baseline {
			t.Fatalf("the lock directory moved to %q, want %q", got, baseline)
		}
	})
}

// #575 L6: the lock directory sits in /tmp and was used as found. Another local
// user could pre-create it, open it up, or point it elsewhere.
func TestEnsureSignerLockDir(t *testing.T) {
	t.Run("creates a missing directory closed to everyone else", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "spectre-relayer-test")
		if err := ensureSignerLockDir(dir); err != nil {
			t.Fatal(err)
		}
		info, err := os.Lstat(dir)
		if err != nil || !info.IsDir() || info.Mode().Perm() != 0o700 {
			t.Fatalf("lock directory = %v, %v; want a 0700 directory", info, err)
		}
		// Existing and correct is accepted as is.
		if err := ensureSignerLockDir(dir); err != nil {
			t.Fatalf("an existing 0700 directory was refused: %v", err)
		}
	})

	t.Run("refuses a symlink", func(t *testing.T) {
		base := t.TempDir()
		target := filepath.Join(base, "elsewhere")
		if err := os.Mkdir(target, 0o700); err != nil {
			t.Fatal(err)
		}
		link := filepath.Join(base, "spectre-relayer-test")
		if err := os.Symlink(target, link); err != nil {
			t.Fatal(err)
		}
		if err := ensureSignerLockDir(link); err == nil || !strings.Contains(err.Error(), "symlink") {
			t.Fatalf("err = %v, want the symlink refused", err)
		}
	})

	t.Run("refuses a directory other users can write", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "spectre-relayer-test")
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(dir, 0o777); err != nil {
			t.Fatal(err)
		}
		if err := ensureSignerLockDir(dir); err == nil || !strings.Contains(err.Error(), "0777") {
			t.Fatalf("err = %v, want the world-writable directory refused", err)
		}
	})

	t.Run("refuses a directory another user owns", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "spectre-relayer-test")
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		restore := signerLockOwner
		signerLockOwner = func() int { return os.Getuid() + 1 } // the directory is ours; pretend we are someone else
		t.Cleanup(func() { signerLockOwner = restore })
		if err := ensureSignerLockDir(dir); err == nil || !strings.Contains(err.Error(), "owned by uid") {
			t.Fatalf("err = %v, want the foreign-owned directory refused", err)
		}
	})

	t.Run("refuses a file in its place", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "spectre-relayer-test")
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := ensureSignerLockDir(path); err == nil {
			t.Fatal("a regular file was accepted as the lock directory")
		}
	})
}

func TestOpenSignerLockFile_DoesNotFollowASymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "cosmos-signer-x.lock")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if file, err := openSignerLockFile(link); err == nil {
		file.Close()
		t.Fatal("a symlinked lock file was opened")
	}
	file, err := openSignerLockFile(filepath.Join(dir, "evm-signer-1-0xabc.lock"))
	if err != nil {
		t.Fatalf("a fresh lock file was refused: %v", err)
	}
	file.Close()
}

// The lock functions themselves refuse an untrusted directory.
func TestSignerLocksRefuseAnUntrustedDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "spectre-relayer-test")
	if err := os.Mkdir(dir, 0o777); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o777); err != nil {
		t.Fatal(err)
	}
	restoreDir := signerLockDir
	signerLockDir = func() (string, error) { return dir, nil }
	restoreAddr := cosmosSignerAddressForLock
	cosmosSignerAddressForLock = func() (string, error) { return "cosmos1lockaddr", nil }
	t.Cleanup(func() { signerLockDir, cosmosSignerAddressForLock = restoreDir, restoreAddr })

	if release, err := acquireCosmosSignerLock(); err == nil {
		release()
		t.Fatal("the Cosmos signer lock was taken in a world-writable directory")
	}
	if release, err := acquireEVMSignerLockKeys([]evmSignerLockKey{{chainID: "1"}}); err == nil {
		release()
		t.Fatal("the EVM signer lock was taken in a world-writable directory")
	}
}
