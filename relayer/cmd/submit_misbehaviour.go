package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	contractICS26Router "relayer/bindings/ICS26Router"
	spectreContract "relayer/bindings/SpectreClient"
	"relayer/prover"
	"relayer/services"
	"relayer/transaction"

	"github.com/cosmos/gogoproto/jsonpb"
	tenderminttypes "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	flagEvidence = "evidence"
	flagDryRun   = "dry-run"
)

type misbehaviourCalldata struct {
	To         string `json:"to"`
	Data       string `json:"data"`
	Value      string `json:"value"`
	ChainID    string `json:"chain_id"`
	Method     string `json:"method"`
	ClientID   string `json:"client_id"`
	Height     uint64 `json:"height"`
	BlockHash1 string `json:"block_hash_1"`
	BlockHash2 string `json:"block_hash_2"`
}

func loadMisbehaviourEvidence(path string) (*tenderminttypes.Misbehaviour, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open evidence: %w", err)
	}
	defer file.Close()

	var evidence tenderminttypes.Misbehaviour
	if err := jsonpb.Unmarshal(file, &evidence); err != nil {
		return nil, fmt.Errorf("decode ibc.lightclients.tendermint.v1.Misbehaviour JSON: %w", err)
	}
	return &evidence, nil
}

func buildMisbehaviourCalldata(
	evm services.EVMEndpoint,
	cosmosRouterClientID string,
	prepared services.PreparedCosmosMisbehaviour,
) (common.Address, string, []byte, error) {
	router := evm.RouterContract()
	if transaction.RouterManagesProofSubmission(evm) {
		parsed, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
		if err != nil {
			return common.Address{}, "", nil, fmt.Errorf("load ICS26Router ABI: %w", err)
		}
		data, err := parsed.Pack("submitMisbehaviour", cosmosRouterClientID, prepared.EncodedMessage)
		if err != nil {
			return common.Address{}, "", nil, fmt.Errorf("pack ICS26Router.submitMisbehaviour: %w", err)
		}
		return *router, "ICS26Router.submitMisbehaviour", data, nil
	}

	spectre := evm.SpectreClientContract()
	if spectre == nil || *spectre == (common.Address{}) {
		return common.Address{}, "", nil, fmt.Errorf("SpectreClient address is not configured")
	}
	parsed, err := spectreContract.ContractSpectreClientMetaData.GetAbi()
	if err != nil {
		return common.Address{}, "", nil, fmt.Errorf("load SpectreClient ABI: %w", err)
	}
	data, err := parsed.Pack("misbehaviour", prepared.EncodedMessage)
	if err != nil {
		return common.Address{}, "", nil, fmt.Errorf("pack SpectreClient.misbehaviour: %w", err)
	}
	return *spectre, "SpectreClient.misbehaviour", data, nil
}

func writeMisbehaviourDryRun(
	cmd *cobra.Command,
	evm services.EVMEndpoint,
	cosmosRouterClientID string,
	prepared services.PreparedCosmosMisbehaviour,
) error {
	to, method, data, err := buildMisbehaviourCalldata(evm, cosmosRouterClientID, prepared)
	if err != nil {
		return err
	}
	chainID, err := evm.EthClient().ChainID(cmd.Context())
	if err != nil {
		return fmt.Errorf("query EVM chain ID: %w", err)
	}
	out := misbehaviourCalldata{
		To:         to.Hex(),
		Data:       "0x" + hex.EncodeToString(data),
		Value:      "0x0",
		ChainID:    chainID.String(),
		Method:     method,
		ClientID:   cosmosRouterClientID,
		Height:     prepared.Height,
		BlockHash1: "0x" + hex.EncodeToString(prepared.BlockHash1[:]),
		BlockHash2: "0x" + hex.EncodeToString(prepared.BlockHash2[:]),
	}
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(out)
}

// allowMisbehaviourEnvOverride preserves the legacy single-source environment
// overrides without letting stale process state replace an explicit --source
// selection from a multi-source config.
func allowMisbehaviourEnvOverride(cfg *appConfig) bool {
	return len(cfg.CosmosToEthConfigs)+len(cfg.CosmosToL2Configs) == 1
}

// SubmitMisbehaviour creates the explicit incident-response path for bringing
// independently observed conflicting Cosmos branches to SpectreClient.
func SubmitMisbehaviour(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-misbehaviour",
		Short: "prove and submit two independently observed conflicting Cosmos headers",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigPath)
			if err != nil {
				return fmt.Errorf("get config flag: %w", err)
			}
			evidencePath, err := cmd.Flags().GetString(flagEvidence)
			if err != nil {
				return fmt.Errorf("get evidence flag: %w", err)
			}
			source, err := cmd.Flags().GetString(flagSource)
			if err != nil {
				return fmt.Errorf("get source flag: %w", err)
			}
			dryRun, err := cmd.Flags().GetBool(flagDryRun)
			if err != nil {
				return fmt.Errorf("get dry-run flag: %w", err)
			}

			_ = godotenv.Load()
			if _, err := transaction.MisbehaviourGasLimit(); err != nil {
				return err
			}
			if !dryRun {
				if err := transaction.ValidateMisbehaviourPrivateKey(); err != nil {
					return err
				}
			}
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			allowEnvOverride := allowMisbehaviourEnvOverride(cfg)
			cfg, err = selectSource(cfg, source)
			if err != nil {
				return err
			}
			evidence, err := loadMisbehaviourEvidence(evidencePath)
			if err != nil {
				return err
			}

			binDir := envOrDefault("PROVER_BIN_DIR", "./bin")
			selectedBackend, hasBackendOverride, err := proofBackendFromFlags(cmd)
			if err != nil {
				return fmt.Errorf("resolve proof backend: %w", err)
			}
			var proverEngine *prover.EcipProver
			if hasBackendOverride {
				proverEngine, err = prover.NewProverWithBackend(binDir, selectedBackend)
			} else {
				proverEngine, err = prover.NewProver(binDir)
			}
			if err != nil {
				return fmt.Errorf("load prover: %w", err)
			}

			txHandler := &transaction.Handler{}
			svc, evidenceDeps, cleanup, err := buildCosmosToEthSource(
				logger,
				cfg.CosmosToEthConfig,
				cfg.EthToCosmosConfig,
				cfg.BatchConfig,
				proverEngine,
				txHandler,
				buildCosmosToEthSourceOptions{allowEnvOverride: allowEnvOverride},
			)
			if err != nil {
				return fmt.Errorf("build evidence context: %w", err)
			}
			defer cleanup()

			if dryRun {
				prepared, err := svc.PrepareCosmosMisbehaviour(
					cmd.Context(), evidenceDeps.Cosmos, evidenceDeps.EVM, evidenceDeps.IDs.CosmosOnEVM, evidence,
				)
				if err != nil {
					return err
				}
				return writeMisbehaviourDryRun(cmd, evidenceDeps.EVM, evidenceDeps.IDs.CosmosOnEVM, prepared)
			}

			prepared, err := svc.SubmitCosmosMisbehaviourEvidence(
				cmd.Context(), evidenceDeps.Cosmos, evidenceDeps.EVM, evidenceDeps.IDs.CosmosOnEVM, evidence,
			)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "submitted misbehaviour for client %s at height %d (%x != %x)\n",
				evidenceDeps.IDs.CosmosOnEVM, prepared.Height, prepared.BlockHash1, prepared.BlockHash2)
			return nil
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().String(flagEvidence, "", "path to ibc.lightclients.tendermint.v1.Misbehaviour JSON evidence")
	cmd.Flags().String(flagSource, "", "ics26_client_id of the Cosmos source to target (required when several are configured)")
	cmd.Flags().Bool(flagDryRun, false, "generate proofs and print transaction calldata without submitting")
	cmd.Flags().Bool(flagGPUProve, false, "use the ICICLE GPU proof backend")
	_ = cmd.MarkFlagRequired(flagEvidence)
	return cmd
}
