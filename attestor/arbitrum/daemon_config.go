package arbitrum

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// DaemonConfig contains the gRPC listener, finalized-L1 RollupCore source, and
// managed Nitro process settings.
type DaemonConfig struct {
	GRPCListenAddress      string   `json:"grpc_listen_address"`
	RuntimePollInterval    string   `json:"runtime_poll_interval"`
	SrcChain               string   `json:"src_chain"`
	L1RPCURL               string   `json:"l1_rpc_url"`
	L1ChainID              uint64   `json:"l1_chain_id"`
	L2ChainID              uint64   `json:"l2_chain_id"`
	RollupCoreAddress      string   `json:"rollup_core_address"`
	AssertionsMappingSlot  string   `json:"assertions_mapping_slot"`
	AssertionStatusOffset  uint8    `json:"assertion_status_offset"`
	AssertionStartBlock    uint64   `json:"assertion_start_block"`
	AssertionPollInterval  string   `json:"assertion_poll_interval"`
	AssertionMaxBlockRange uint64   `json:"assertion_max_block_range"`
	AttestorStatePath      string   `json:"attestor_state_path"`
	NitroBinaryPath        string   `json:"nitro_binary_path"`
	NitroBinarySHA256      string   `json:"nitro_binary_sha256"`
	NitroArguments         []string `json:"nitro_arguments"`
	NitroWorkDir           string   `json:"nitro_work_dir"`
	NitroDataDir           string   `json:"nitro_data_dir"`
	NitroIPCPath           string   `json:"nitro_ipc_path"`
	NitroStartupTimeout    string   `json:"nitro_startup_timeout"`
	NitroShutdownTimeout   string   `json:"nitro_shutdown_timeout"`
}

// LoadDaemonConfig reads, validates, and resolves filesystem paths relative to
// the directory containing the configuration file.
func LoadDaemonConfig(path string) (DaemonConfig, error) {
	if path == "" {
		return DaemonConfig{}, errors.New("attestor config path must not be empty")
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return DaemonConfig{}, fmt.Errorf("resolve attestor config path %q: %w", path, err)
	}
	data, err := os.ReadFile(absolutePath)
	if err != nil {
		return DaemonConfig{}, fmt.Errorf("read attestor config %q: %w", absolutePath, err)
	}
	var config DaemonConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return DaemonConfig{}, fmt.Errorf("decode attestor config %q: %w", absolutePath, err)
	}

	base := filepath.Dir(absolutePath)
	config.NitroBinaryPath = resolveCommandPath(base, config.NitroBinaryPath)
	config.NitroDataDir = resolveConfigPath(base, config.NitroDataDir)
	config.NitroIPCPath = resolveConfigPath(base, config.NitroIPCPath)
	config.AttestorStatePath = resolveConfigPath(base, config.AttestorStatePath)
	if config.NitroWorkDir == "" {
		config.NitroWorkDir = base
	} else {
		config.NitroWorkDir = resolveConfigPath(base, config.NitroWorkDir)
	}
	if err := config.Validate(); err != nil {
		return DaemonConfig{}, fmt.Errorf("validate attestor config %q: %w", absolutePath, err)
	}
	return config, nil
}

// Validate checks the gRPC listener and managed Nitro process settings.
func (c DaemonConfig) Validate() error {
	if err := validateListenAddress(c.GRPCListenAddress); err != nil {
		return err
	}
	if _, err := c.RuntimePollDuration(); err != nil {
		return err
	}
	if c.SrcChain == "" {
		return errors.New("src_chain is required")
	}
	if err := validateRPCURL(c.L1RPCURL, "l1_rpc_url"); err != nil {
		return err
	}
	if c.L1ChainID == 0 {
		return errors.New("l1_chain_id must be greater than zero")
	}
	if c.L2ChainID == 0 {
		return errors.New("l2_chain_id must be greater than zero")
	}
	if !common.IsHexAddress(c.RollupCoreAddress) ||
		common.HexToAddress(c.RollupCoreAddress) == (common.Address{}) {
		return errors.New("rollup_core_address must be a non-zero EVM address")
	}
	if _, err := decodeHash32(c.AssertionsMappingSlot, "assertions_mapping_slot"); err != nil {
		return err
	}
	if c.AssertionStatusOffset >= 32 {
		return errors.New("assertion_status_offset must be between 0 and 31")
	}
	if c.AssertionStartBlock == 0 {
		return errors.New("assertion_start_block must be a non-zero RollupCore log scan start block")
	}
	if _, err := c.AssertionPollDuration(); err != nil {
		return err
	}
	if c.AttestorStatePath == "" || !filepath.IsAbs(c.AttestorStatePath) {
		return errors.New("attestor_state_path must resolve to an absolute path")
	}
	if c.NitroBinaryPath == "" {
		return errors.New("nitro_binary_path is required")
	}
	if _, err := c.binaryDigest(); err != nil {
		return err
	}
	if c.NitroWorkDir == "" {
		return errors.New("nitro_work_dir is required")
	}
	if c.NitroDataDir == "" || !filepath.IsAbs(c.NitroDataDir) {
		return errors.New("nitro_data_dir must resolve to an absolute path")
	}
	if c.NitroIPCPath == "" || !filepath.IsAbs(c.NitroIPCPath) {
		return errors.New("nitro_ipc_path must resolve to an absolute path")
	}
	for _, argument := range c.NitroArguments {
		switch {
		case ownsNitroArgument(argument, "--ipc.path"):
			return errors.New("nitro_arguments must not set --ipc.path; the attestor owns the IPC endpoint")
		case ownsNitroArgument(argument, "--persistent.chain"):
			return errors.New("nitro_arguments must not set --persistent.chain; use nitro_data_dir")
		}
	}
	if _, err := parsePositiveDuration(c.NitroStartupTimeout, "nitro_startup_timeout"); err != nil {
		return err
	}
	if _, err := parsePositiveDuration(c.NitroShutdownTimeout, "nitro_shutdown_timeout"); err != nil {
		return err
	}
	return nil
}

// RuntimePollDuration returns the fallback interval for reconciling Nitro's
// unsafe, safe, and finalized heads when no new-head event is received. An
// omitted interval defaults to 15 seconds.
func (c DaemonConfig) RuntimePollDuration() (time.Duration, error) {
	const defaultRuntimePollInterval = 15 * time.Second
	if c.RuntimePollInterval == "" {
		return defaultRuntimePollInterval, nil
	}
	return parsePositiveDuration(c.RuntimePollInterval, "runtime_poll_interval")
}

// AssertionPollDuration returns the finalized-L1 assertion polling interval.
// An omitted interval defaults to 12 seconds.
func (c DaemonConfig) AssertionPollDuration() (time.Duration, error) {
	const defaultAssertionPollInterval = 12 * time.Second
	if c.AssertionPollInterval == "" {
		return defaultAssertionPollInterval, nil
	}
	return parsePositiveDuration(c.AssertionPollInterval, "assertion_poll_interval")
}

// AssertionBlockRange returns the maximum number of finalized L1 blocks read
// in one log query.
func (c DaemonConfig) AssertionBlockRange() uint64 {
	const defaultAssertionMaxBlockRange = uint64(2_000)
	if c.AssertionMaxBlockRange == 0 {
		return defaultAssertionMaxBlockRange
	}
	return c.AssertionMaxBlockRange
}

// NitroProcessConfig returns the validated child-process settings. The
// attestor injects both --persistent.chain and --ipc.path.
func (c DaemonConfig) NitroProcessConfig() (NitroProcessConfig, error) {
	if err := c.Validate(); err != nil {
		return NitroProcessConfig{}, err
	}
	digest, err := c.binaryDigest()
	if err != nil {
		return NitroProcessConfig{}, err
	}
	startupTimeout, err := parsePositiveDuration(c.NitroStartupTimeout, "nitro_startup_timeout")
	if err != nil {
		return NitroProcessConfig{}, err
	}
	shutdownTimeout, err := parsePositiveDuration(c.NitroShutdownTimeout, "nitro_shutdown_timeout")
	if err != nil {
		return NitroProcessConfig{}, err
	}
	return NitroProcessConfig{
		BinaryPath:      c.NitroBinaryPath,
		BinarySHA256:    digest,
		Arguments:       append([]string(nil), c.NitroArguments...),
		WorkDir:         c.NitroWorkDir,
		DataDir:         c.NitroDataDir,
		IPCPath:         c.NitroIPCPath,
		StartupTimeout:  startupTimeout,
		ShutdownTimeout: shutdownTimeout,
	}, nil
}

func validateListenAddress(address string) error {
	if address == "" {
		return errors.New("grpc_listen_address is required")
	}
	_, rawPort, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("parse grpc_listen_address: %w", err)
	}
	port, err := strconv.ParseUint(rawPort, 10, 16)
	if err != nil || port == 0 {
		return errors.New("grpc_listen_address must contain a port between 1 and 65535")
	}
	return nil
}

func validateRPCURL(raw, field string) error {
	if raw == "" {
		return fmt.Errorf("%s is required", field)
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("parse %s: %w", field, err)
	}
	switch parsed.Scheme {
	case "http", "https", "ws", "wss":
	default:
		return fmt.Errorf("%s must use http, https, ws, or wss", field)
	}
	if parsed.Host == "" {
		return fmt.Errorf("%s must include a host", field)
	}
	return nil
}

func decodeHash32(raw, field string) (common.Hash, error) {
	var value common.Hash
	encoded := strings.TrimPrefix(raw, "0x")
	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		return value, fmt.Errorf("decode %s: %w", field, err)
	}
	if len(decoded) != common.HashLength {
		return value, fmt.Errorf(
			"%s must contain exactly %d bytes, got %d",
			field,
			common.HashLength,
			len(decoded),
		)
	}
	copy(value[:], decoded)
	return value, nil
}

func ownsNitroArgument(argument, name string) bool {
	return argument == name || strings.HasPrefix(argument, name+"=")
}

func (c DaemonConfig) binaryDigest() ([sha256.Size]byte, error) {
	var digest [sha256.Size]byte
	encoded := strings.TrimPrefix(c.NitroBinarySHA256, "0x")
	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		return digest, fmt.Errorf("decode nitro_binary_sha256: %w", err)
	}
	if len(decoded) != len(digest) {
		return digest, fmt.Errorf(
			"nitro_binary_sha256 must contain exactly %d bytes, got %d",
			len(digest),
			len(decoded),
		)
	}
	copy(digest[:], decoded)
	return digest, nil
}

func parsePositiveDuration(raw, field string) (time.Duration, error) {
	if raw == "" {
		return 0, fmt.Errorf("%s is required", field)
	}
	duration, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", field, err)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", field)
	}
	return duration, nil
}

func resolveConfigPath(base, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(filepath.Join(base, path))
}

func resolveCommandPath(base, path string) string {
	if path == "" || filepath.IsAbs(path) || !strings.ContainsRune(path, os.PathSeparator) {
		return path
	}
	return filepath.Clean(filepath.Join(base, path))
}
