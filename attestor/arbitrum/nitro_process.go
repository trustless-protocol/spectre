package arbitrum

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

const nitroIPCConnectPollInterval = 100 * time.Millisecond

// NitroProcessConfig describes the Nitro execution process owned by the
// attestor. BinarySHA256 pins the exact executable that is allowed to become
// the local execution authority.
type NitroProcessConfig struct {
	BinaryPath      string
	BinarySHA256    [sha256.Size]byte
	Arguments       []string
	WorkDir         string
	DataDir         string
	IPCPath         string
	StartupTimeout  time.Duration
	ShutdownTimeout time.Duration
}

// ManagedNitroProcess owns one locally launched Nitro execution node and its
// private IPC connection. Nitro's chain database remains in DataDir across
// attestor restarts.
type ManagedNitroProcess struct {
	command *exec.Cmd
	rpc     *rpc.Client
	ipcPath string

	done    chan struct{}
	waitMu  sync.Mutex
	waitErr error

	shutdownTimeout time.Duration
	closeOnce       sync.Once
	closeErr        error
}

// StartManagedNitro verifies and launches the pinned Nitro executable, then
// waits for the fresh IPC endpoint to expose the eth API used by the gRPC
// verifier. The IPC path must not exist before the process is started.
func StartManagedNitro(
	ctx context.Context,
	config NitroProcessConfig,
	stdout io.Writer,
	stderr io.Writer,
) (*ManagedNitroProcess, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}
	binaryPath, err := exec.LookPath(config.BinaryPath)
	if err != nil {
		return nil, fmt.Errorf("find managed Nitro executable %q: %w", config.BinaryPath, err)
	}
	binaryPath, err = filepath.Abs(binaryPath)
	if err != nil {
		return nil, fmt.Errorf("resolve managed Nitro executable %q: %w", binaryPath, err)
	}
	if err := verifyFileSHA256(binaryPath, config.BinarySHA256); err != nil {
		return nil, err
	}
	if err := requireFreshIPCPath(config.IPCPath); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(config.DataDir, 0o700); err != nil {
		return nil, fmt.Errorf("create managed Nitro persistent data directory: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(config.IPCPath), 0o700); err != nil {
		return nil, fmt.Errorf("create managed Nitro IPC directory: %w", err)
	}

	arguments := append([]string(nil), config.Arguments...)
	arguments = append(arguments, "--persistent.chain="+config.DataDir)
	arguments = append(arguments, "--ipc.path="+config.IPCPath)
	command := exec.Command(binaryPath, arguments...)
	command.Dir = config.WorkDir
	command.Stdout = stdout
	command.Stderr = stderr
	if err := command.Start(); err != nil {
		return nil, fmt.Errorf("start managed Nitro process: %w", err)
	}

	process := &ManagedNitroProcess{
		command:         command,
		ipcPath:         config.IPCPath,
		done:            make(chan struct{}),
		shutdownTimeout: config.ShutdownTimeout,
	}
	go process.wait()

	startupCtx, cancel := context.WithTimeout(ctx, config.StartupTimeout)
	defer cancel()
	client, err := process.waitForIPC(startupCtx)
	if err != nil {
		closeErr := process.Close()
		return nil, errors.Join(err, closeErr)
	}
	process.rpc = client
	return process, nil
}

func (c NitroProcessConfig) validate() error {
	if c.BinaryPath == "" {
		return errors.New("managed Nitro binary path must not be empty")
	}
	if c.WorkDir == "" {
		return errors.New("managed Nitro work directory must not be empty")
	}
	info, err := os.Stat(c.WorkDir)
	if err != nil {
		return fmt.Errorf("stat managed Nitro work directory %q: %w", c.WorkDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("managed Nitro work directory %q is not a directory", c.WorkDir)
	}
	if c.DataDir == "" || !filepath.IsAbs(c.DataDir) {
		return errors.New("managed Nitro persistent data directory must be absolute")
	}
	if info, err := os.Stat(c.DataDir); err == nil && !info.IsDir() {
		return fmt.Errorf("managed Nitro persistent data path %q is not a directory", c.DataDir)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat managed Nitro persistent data directory %q: %w", c.DataDir, err)
	}
	if c.IPCPath == "" || !filepath.IsAbs(c.IPCPath) {
		return errors.New("managed Nitro IPC path must be absolute")
	}
	if c.StartupTimeout <= 0 {
		return errors.New("managed Nitro startup timeout must be greater than zero")
	}
	if c.ShutdownTimeout <= 0 {
		return errors.New("managed Nitro shutdown timeout must be greater than zero")
	}
	for _, argument := range c.Arguments {
		if ownsNitroArgument(argument, "--ipc.path") || ownsNitroArgument(argument, "--persistent.chain") {
			return fmt.Errorf("managed Nitro argument %q overrides an attestor-owned path", argument)
		}
	}
	return nil
}

func verifyFileSHA256(path string, expected [sha256.Size]byte) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open managed Nitro executable %q: %w", path, err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("hash managed Nitro executable %q: %w", path, err)
	}
	actual := hasher.Sum(nil)
	if subtle.ConstantTimeCompare(actual, expected[:]) != 1 {
		return fmt.Errorf(
			"managed Nitro executable SHA-256 mismatch: path=%q expected=%x actual=%x",
			path,
			expected,
			actual,
		)
	}
	return nil
}

func requireFreshIPCPath(path string) error {
	_, err := os.Lstat(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil
	case err != nil:
		return fmt.Errorf("inspect managed Nitro IPC path %q: %w", path, err)
	default:
		return fmt.Errorf(
			"managed Nitro IPC path %q already exists; refusing to attach to an unowned endpoint",
			path,
		)
	}
}

func (p *ManagedNitroProcess) wait() {
	err := p.command.Wait()
	p.waitMu.Lock()
	p.waitErr = err
	p.waitMu.Unlock()
	close(p.done)
}

func (p *ManagedNitroProcess) waitForIPC(ctx context.Context) (*rpc.Client, error) {
	ticker := time.NewTicker(nitroIPCConnectPollInterval)
	defer ticker.Stop()

	var lastError error
	for {
		client, err := rpc.DialIPC(ctx, p.ipcPath)
		if err == nil {
			var modules map[string]string
			if err := client.CallContext(ctx, &modules, "rpc_modules"); err != nil {
				client.Close()
				lastError = fmt.Errorf("query managed Nitro RPC modules: %w", err)
			} else if _, found := modules["eth"]; !found {
				client.Close()
				return nil, errors.New("managed Nitro IPC does not expose the required eth API")
			} else {
				return client, nil
			}
		} else {
			lastError = err
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("wait for managed Nitro IPC at %q: %w (last error: %v)", p.ipcPath, ctx.Err(), lastError)
		case <-p.done:
			exitErr := p.Err()
			if exitErr == nil {
				exitErr = errors.New("process exited successfully")
			}
			return nil, fmt.Errorf("managed Nitro process exited before IPC became ready: %w", exitErr)
		case <-ticker.C:
		}
	}
}

// RPC returns the owned IPC client. The process retains ownership and closes
// the client during Close.
func (p *ManagedNitroProcess) RPC() *rpc.Client {
	if p == nil {
		return nil
	}
	return p.rpc
}

// Done closes when the managed Nitro process exits.
func (p *ManagedNitroProcess) Done() <-chan struct{} {
	if p == nil {
		closed := make(chan struct{})
		close(closed)
		return closed
	}
	return p.done
}

// Err returns Nitro's process exit result. It returns an error if called while
// the process is still running.
func (p *ManagedNitroProcess) Err() error {
	if p == nil {
		return errors.New("managed Nitro process is nil")
	}
	select {
	case <-p.done:
		p.waitMu.Lock()
		defer p.waitMu.Unlock()
		return p.waitErr
	default:
		return errors.New("managed Nitro process is still running")
	}
}

// Close closes IPC, requests a graceful interrupt, and kills Nitro only if it
// does not stop within the configured timeout. It is idempotent.
func (p *ManagedNitroProcess) Close() error {
	if p == nil {
		return nil
	}
	p.closeOnce.Do(func() {
		if p.rpc != nil {
			p.rpc.Close()
		}

		select {
		case <-p.done:
		default:
			if err := p.command.Process.Signal(os.Interrupt); err != nil && !errors.Is(err, os.ErrProcessDone) {
				p.closeErr = errors.Join(p.closeErr, fmt.Errorf("interrupt managed Nitro process: %w", err))
			}
			select {
			case <-p.done:
			case <-time.After(p.shutdownTimeout):
				if err := p.command.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
					p.closeErr = errors.Join(p.closeErr, fmt.Errorf("kill managed Nitro process: %w", err))
				}
				<-p.done
			}
		}
		p.closeErr = errors.Join(p.closeErr, removeOwnedIPC(p.ipcPath))
	})
	return p.closeErr
}

func removeOwnedIPC(path string) error {
	info, err := os.Lstat(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil
	case err != nil:
		return fmt.Errorf("inspect managed Nitro IPC during cleanup: %w", err)
	case info.Mode()&os.ModeSocket == 0:
		return fmt.Errorf("refusing to remove non-socket at managed Nitro IPC path %q", path)
	default:
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove managed Nitro IPC socket: %w", err)
		}
		return nil
	}
}
