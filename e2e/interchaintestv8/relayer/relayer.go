package relayer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"

	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"

	ethereumtypes "github.com/srdtrk/solidity-ibc-eureka/e2e/v8/types/ethereum"
	relayertypes "github.com/srdtrk/solidity-ibc-eureka/e2e/v8/types/relayer"
)

// DefaultRelayerGRPCAddress returns the default gRPC address for the relayer.
func DefaultRelayerGRPCAddress() string {
	return "127.0.0.1:3000"
}

// binaryPath returns the path to the relayer binary.
//
// Resolution order:
//  1. $RELAYER_BINARY (absolute path, used by the justfile test-e2e target)
//  2. `relayer` on $PATH (works after `just install-go-relayer`)
//
// Panics with a clear hint if neither is available — tests cannot proceed without the binary.
func binaryPath() string {
	if p := os.Getenv("RELAYER_BINARY"); p != "" {
		return p
	}
	if p, err := exec.LookPath("relayer"); err == nil {
		return p
	}
	panic("relayer binary not found: set $RELAYER_BINARY or run `just install-go-relayer`")
}

// proverEnv exposes the per-bucket prover artifact directory to the spawned relayer process.
//
// The relayer loads artifacts via prover.NewProver(binDir) and expects a directory layout of
// $PROVER_BIN_DIR/n{N}/{r1cs.bin,pk.bin,vk.bin} for each compiled bucket. Run
// `just build-prover-artifacts` (or `go run ./relayer/prover/cmd`) to populate it.
func proverEnv() []string {
	dir := os.Getenv("PROVER_BIN_DIR")
	if dir == "" {
		panic("PROVER_BIN_DIR not set: must point to relayer/bin with n{N}/ subdirs (run `just build-prover-artifacts`)")
	}
	return []string{"PROVER_BIN_DIR=" + dir}
}

// StartRelayer starts the relayer with the given config file.
func StartRelayer(configPath string) (*os.Process, error) {
	config, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Starting relayer with config:\n%s\n", config)

	cmd := exec.Command(binaryPath(), "start", "--config", configPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), proverEnv()...)

	// run this command in the background
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	exited := make(chan error, 1)
	go func() {
		exited <- cmd.Wait()
	}()

	// Wait long enough for the relayer to finish loading the gnark prover artifacts
	// (per-bucket r1cs/pk ~200MB, takes ~20s on first read) AND set up both Cosmos and
	// Ethereum event subscriptions. If the test sends a SendPacket tx before WatchSendPacket
	// is active, the event is missed (Watch defaults to fromBlock=latest).
	select {
	case err := <-exited:
		if err != nil {
			return nil, fmt.Errorf("relayer exited during startup: %w", err)
		}
		return nil, errors.New("relayer exited during startup")
	case <-time.After(60 * time.Second):
	}

	return cmd.Process, nil
}

// RunCreateClients runs the relayer's create-clients command synchronously (blocks until completion).
func RunCreateClients(configPath string, extraArgs ...string) error {
	config, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	fmt.Printf("Running create-clients with config:\n%s\n", config)

	args := append([]string{"create-clients", "--config", configPath}, extraArgs...)
	cmd := exec.Command(binaryPath(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), proverEnv()...)

	return cmd.Run()
}

// RunCreateClientsEth runs the relayer's create-clients-eth command synchronously.
func RunCreateClientsEth(configPath string, extraArgs ...string) error {
	config, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	fmt.Printf("Running create-clients-eth with config:\n%s\n", config)

	args := append([]string{"create-clients-eth", "--config", configPath}, extraArgs...)
	cmd := exec.Command(binaryPath(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), proverEnv()...)

	return cmd.Run()
}

// RunUpdateClient runs the relayer's update-client command synchronously.
func RunUpdateClient(configPath string, extraArgs ...string) error {
	config, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	fmt.Printf("Running update-client with config:\n%s\n", config)

	args := append([]string{"update-client", "--config", configPath}, extraArgs...)
	cmd := exec.Command(binaryPath(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), proverEnv()...)

	return cmd.Run()
}

// GetGRPCClient returns a gRPC client for the relayer.
func GetGRPCClient(addr string) (relayertypes.RelayerServiceClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return relayertypes.NewRelayerServiceClient(conn), nil
}

// GetRelayUpdateSlot extracts the latest update slot from the relay body's update messages.
func GetRelayUpdateSlotForWasmClient(relayBody []byte) (uint64, error) {
	var txBody txtypes.TxBody
	err := proto.Unmarshal(relayBody, &txBody)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal relay body: %w", err)
	}

	var updateClientMsgsAny []*codectypes.Any
	for _, msg := range txBody.Messages {
		if msg.TypeUrl == "/ibc.core.client.v1.MsgUpdateClient" {
			updateClientMsgsAny = append(updateClientMsgsAny, msg)
		}
	}
	if len(updateClientMsgsAny) == 0 {
		return 0, errors.New("no update client messages found in relay body")
	}

	var headers []ethereumtypes.Header
	for _, updateClientMsgAny := range updateClientMsgsAny {
		var updateClientMsg clienttypes.MsgUpdateClient
		err = proto.Unmarshal(updateClientMsgAny.Value, &updateClientMsg)
		if err != nil {
			return 0, fmt.Errorf("failed to unmarshal MsgUpdateClient: %w", err)
		}
		var clientMessage ibcwasmtypes.ClientMessage
		err = proto.Unmarshal(updateClientMsg.ClientMessage.Value, &clientMessage)
		if err != nil {
			return 0, fmt.Errorf("failed to unmarshal ClientMessage: %w", err)
		}

		var header ethereumtypes.Header
		err = json.Unmarshal(clientMessage.Data, &header)
		if err != nil {
			return 0, fmt.Errorf("failed to unmarshal header: %w", err)
		}

		headers = append(headers, header)
	}

	latestUpdateSlot := uint64(0)
	for _, header := range headers {
		updateSlot, err := strconv.ParseUint(header.ConsensusUpdate.FinalizedHeader.Beacon.Slot, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse slot from header: %w", err)
		}
		latestUpdateSlot = max(latestUpdateSlot, updateSlot)
	}

	return latestUpdateSlot, nil
}
