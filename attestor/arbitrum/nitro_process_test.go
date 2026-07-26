package arbitrum

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

func TestManagedNitroProcessOwnsFreshIPCAndStops(t *testing.T) {
	t.Setenv("GO_WANT_NITRO_HELPER_PROCESS", "1")
	binary, err := os.Executable()
	if err != nil {
		t.Fatalf("find test executable: %v", err)
	}
	workDir := shortTempDir(t)
	dataDir := filepath.Join(workDir, "chain")
	t.Setenv("GO_NITRO_HELPER_DATA_DIR", dataDir)
	ipcPath := filepath.Join(workDir, "nitro.ipc")
	process, err := StartManagedNitro(
		context.Background(),
		NitroProcessConfig{
			BinaryPath:      binary,
			BinarySHA256:    fileDigest(t, binary),
			Arguments:       []string{"-test.run=^TestManagedNitroHelperProcess$", "--"},
			WorkDir:         workDir,
			DataDir:         dataDir,
			IPCPath:         ipcPath,
			StartupTimeout:  5 * time.Second,
			ShutdownTimeout: 2 * time.Second,
		},
		os.Stdout,
		os.Stderr,
	)
	if err != nil {
		t.Fatalf("start managed Nitro helper: %v", err)
	}
	if process.RPC() == nil {
		t.Fatal("managed Nitro RPC client is nil")
	}
	var modules map[string]string
	if err := process.RPC().CallContext(context.Background(), &modules, "rpc_modules"); err != nil {
		t.Fatalf("query helper modules: %v", err)
	}
	if _, found := modules["eth"]; !found {
		t.Fatalf("eth module missing: %v", modules)
	}
	marker := filepath.Join(dataDir, "persistent-marker")
	if err := os.WriteFile(marker, []byte("keep"), 0o600); err != nil {
		t.Fatalf("write persistent marker: %v", err)
	}
	if err := process.Close(); err != nil {
		t.Fatalf("close managed Nitro helper: %v", err)
	}
	select {
	case <-process.Done():
	case <-time.After(3 * time.Second):
		t.Fatal("managed Nitro helper did not exit")
	}
	if _, err := os.Lstat(ipcPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("IPC path remains after shutdown: %v", err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("Nitro persistent data was not retained: %v", err)
	}
}

func TestManagedNitroProcessRejectsExistingIPC(t *testing.T) {
	binary, err := os.Executable()
	if err != nil {
		t.Fatalf("find test executable: %v", err)
	}
	workDir := t.TempDir()
	dataDir := filepath.Join(workDir, "chain")
	ipcPath := filepath.Join(workDir, "nitro.ipc")
	if err := os.WriteFile(ipcPath, []byte("unowned"), 0o600); err != nil {
		t.Fatalf("create existing IPC path: %v", err)
	}
	_, err = StartManagedNitro(
		context.Background(),
		NitroProcessConfig{
			BinaryPath:      binary,
			BinarySHA256:    fileDigest(t, binary),
			WorkDir:         workDir,
			DataDir:         dataDir,
			IPCPath:         ipcPath,
			StartupTimeout:  time.Second,
			ShutdownTimeout: time.Second,
		},
		nil,
		nil,
	)
	if err == nil || !strings.Contains(err.Error(), "refusing to attach to an unowned endpoint") {
		t.Fatalf("existing IPC error: got %v", err)
	}
}

func TestManagedNitroProcessRejectsWrongBinaryDigest(t *testing.T) {
	binary, err := os.Executable()
	if err != nil {
		t.Fatalf("find test executable: %v", err)
	}
	workDir := t.TempDir()
	_, err = StartManagedNitro(
		context.Background(),
		NitroProcessConfig{
			BinaryPath:      binary,
			BinarySHA256:    [sha256.Size]byte{1},
			WorkDir:         workDir,
			DataDir:         filepath.Join(workDir, "chain"),
			IPCPath:         filepath.Join(workDir, "nitro.ipc"),
			StartupTimeout:  time.Second,
			ShutdownTimeout: time.Second,
		},
		nil,
		nil,
	)
	if err == nil || !strings.Contains(err.Error(), "SHA-256 mismatch") {
		t.Fatalf("binary digest error: got %v", err)
	}
}

func TestManagedNitroHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_NITRO_HELPER_PROCESS") != "1" {
		return
	}
	if os.Getenv("GO_NITRO_HELPER_EXIT") == "1" {
		os.Exit(17)
	}
	var ipcPath string
	var dataDir string
	for _, argument := range os.Args {
		if strings.HasPrefix(argument, "--ipc.path=") {
			ipcPath = strings.TrimPrefix(argument, "--ipc.path=")
		}
		if strings.HasPrefix(argument, "--persistent.chain=") {
			dataDir = strings.TrimPrefix(argument, "--persistent.chain=")
		}
	}
	if ipcPath == "" {
		os.Exit(18)
	}
	if dataDir == "" || dataDir != os.Getenv("GO_NITRO_HELPER_DATA_DIR") {
		os.Exit(20)
	}
	listener, server, err := rpc.StartIPCEndpoint(ipcPath, []rpc.API{{
		Namespace: "eth",
		Service:   &helperEthAPI{},
	}})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(19)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	server.Stop()
	_ = listener.Close()
	os.Exit(0)
}

type helperEthAPI struct{}

func (*helperEthAPI) BlockNumber(context.Context) (uint64, error) {
	return 1, nil
}

func fileDigest(t *testing.T, path string) [sha256.Size]byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	return sha256.Sum256(data)
}

func shortTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("/tmp", "fast-ibc-nitro-")
	if err != nil {
		t.Fatalf("create short temporary directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("remove short temporary directory: %v", err)
		}
	})
	return dir
}
