package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

// signerLockDirName identifies the shared per-user lock directory. It must not
// be derived from a relay's state or working directory: those are deliberately
// per-path values, so a lock located under either would let two processes take
// different files.
const signerLockDirName = "spectre-relayer"

// signerLockRoot is a literal on purpose: every stdlib helper that would name a
// directory for us reads an environment variable (TMPDIR for os.TempDir,
// XDG_CACHE_HOME/HOME for os.UserCacheDir), and an environment variable is
// exactly what must not decide this path.
const signerLockRoot = "/tmp"

// signerLockDir resolves the one directory every relay process of this user, on
// this machine, agrees on. It is a seam for lock tests, which must not write
// outside their temporary directory.
//
// It reads no environment variable. start loads .env before taking the lock;
// using os.UserCacheDir here would let processes with different HOME or
// XDG_CACHE_HOME values acquire different files and both start.
var signerLockDir = func() (string, error) {
	uid := os.Getuid()
	if uid < 0 {
		cache, err := os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("signer lock: no directory to place it in: %w", err)
		}
		return filepath.Join(cache, signerLockDirName), nil
	}
	return filepath.Join(signerLockRoot, fmt.Sprintf("%s-%d", signerLockDirName, uid)), nil
}

// signerLockOwner is the uid the lock directory must belong to. A seam only so a
// test can stand in for a directory another user owns.
var signerLockOwner = os.Getuid

// ensureSignerLockDir creates the lock directory if it is missing and refuses
// one this user cannot trust.
//
// The directory sits in /tmp, which every local user can write. Trusted as found,
// another user could create it first: owned by them it makes start fail with a
// bare permission error; world-writable, or a symlink to a directory they
// control, it lets them remove a held lock file so two relayers on one key both
// start (#575 L6). It must be a real directory, owned by this user, that no one
// else can write.
func ensureSignerLockDir(dir string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("signer lock: create %q: %w", dir, err)
	}
	info, err := os.Lstat(dir)
	if err != nil {
		return fmt.Errorf("signer lock: inspect %q: %w", dir, err)
	}
	switch {
	case info.Mode()&fs.ModeSymlink != 0:
		return fmt.Errorf("signer lock: %q is a symlink; refusing to place locks where it points. Remove it", dir)
	case !info.IsDir():
		return fmt.Errorf("signer lock: %q is not a directory. Remove it", dir)
	case info.Mode().Perm()&0o022 != 0:
		return fmt.Errorf("signer lock: %q has mode %04o; other users can write it and replace a held lock. "+
			"Run chmod 700 on it, or remove it", dir, info.Mode().Perm())
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if owner, want := int(stat.Uid), signerLockOwner(); owner != want {
			return fmt.Errorf("signer lock: %q is owned by uid %d, not this user (uid %d); "+
				"another account created it. Remove it as that user or root", dir, owner, want)
		}
	}
	return nil
}

// openSignerLockFile opens (creating) one lock file without following a symlink
// in its place.
func openSignerLockFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW, 0o600)
	if errors.Is(err, syscall.ELOOP) {
		return nil, fmt.Errorf("%q is a symlink; refusing to lock through it: %w", path, err)
	}
	return file, err
}
