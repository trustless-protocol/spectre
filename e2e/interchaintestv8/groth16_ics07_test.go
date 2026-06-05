package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"

	"github.com/decentrio/fast-ibc/packages/go-abigen/groth16ics07tendermint"
	"github.com/decentrio/fast-ibc/packages/go-abigen/ics26router"

	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/e2esuite"
	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/ethereum"
	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/relayer"
	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/testvalues"
	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/types"
	relayertypes "github.com/srdtrk/solidity-ibc-eureka/e2e/v8/types/relayer"
)

// Groth16ICS07TendermintTestSuite is a suite of tests that wraps TestSuite
// and can provide additional functionality
type Groth16ICS07TendermintTestSuite struct {
	e2esuite.TestSuite

	// Whether to generate fixtures for the solidity tests
	generateFixtures bool

	// Addresses of the deployed contracts
	groth16Ics07Address ethcommon.Address
	ics26Address        ethcommon.Address

	// The private key of a test account
	key *ecdsa.PrivateKey
	// The Groth16ICS07Tendermint contract
	contract *groth16ics07tendermint.Contract
	// The ICS26 router contract, needed for the relayer to pass proofs
	ics26Contract *ics26router.Contract

	// The relayer API (only used for deployment at the moment)
	RelayerClient relayertypes.RelayerServiceClient
}

// SetupSuite calls the underlying Groth16ICS07TendermintTestSuite's SetupSuite method
// and deploys the Groth16ICS07Tendermint contract
func (s *Groth16ICS07TendermintTestSuite) SetupSuite(ctx context.Context, proofType types.SupportedProofType) {
	s.TestSuite.SetupSuite(ctx)

	eth, simd := s.EthChain, s.CosmosChains[0]

	s.T().Logf("Setting up the test suite with proof type: %s", proofType.String())

	var prover string
	s.Require().True(s.Run("Set up environment", func() {
		err := os.Chdir("../..")
		s.Require().NoError(err)

		s.key, err = eth.CreateAndFundUser()
		s.Require().NoError(err)

		// Use mock verifier in E2E tests (gnark Groth16 prover is the real prover)
		os.Setenv(testvalues.EnvKeyVerifier, testvalues.EnvValueVerifier_Mock)
		os.Setenv(testvalues.EnvKeyEthRPC, eth.RPC)
		os.Setenv(testvalues.EnvKeyTendermintRPC, simd.GetHostRPCAddress())
		os.Setenv(testvalues.EnvKeyOperatorPrivateKey, hex.EncodeToString(crypto.FromECDSA(s.key)))
		s.generateFixtures = os.Getenv(testvalues.EnvKeyGenerateSolidityFixtures) == testvalues.EnvValueGenerateFixtures_True
	}))

	s.Require().True(s.Run("Deploy IBC contracts", func() {
		stdout, err := eth.ForgeScript(s.key, testvalues.E2EDeployScriptPath)
		s.Require().NoError(err)

		contractAddresses, err := ethereum.GetEthContractsFromDeployOutput(string(stdout))
		s.Require().NoError(err)
		s.ics26Address = ethcommon.HexToAddress(contractAddresses.Ics26Router)
		s.ics26Contract, err = ics26router.NewContract(s.ics26Address, eth.RPCClient)
		s.Require().NoError(err)
	}))

	var relayerProcess *os.Process
	s.Require().True(s.Run("Start Relayer", func() {
		beaconAPI := ""
		// The BeaconAPIClient is nil when the testnet is `pow`
		if eth.BeaconAPIClient != nil {
			beaconAPI = eth.BeaconAPIClient.GetBeaconAPIURL()
		}

		config := relayer.NewConfig(relayer.CreateEthCosmosModules(
			relayer.EthCosmosConfigInfo{
				EthChainID:     eth.ChainID.String(),
				CosmosChainID:  simd.Config().ChainID,
				TmRPC:          simd.GetHostRPCAddress(),
				ICS26Address:   s.ics26Address.Hex(),
				EthRPC:         eth.RPC,
				BeaconAPI:      beaconAPI,
				SignerAddress:  "",   // unused
				MockWasmClient: true, // unused
			}),
		)

		err := config.GenerateConfigFile(testvalues.RelayerConfigFilePath)
		s.Require().NoError(err)

		relayerProcess, err = relayer.StartRelayer(testvalues.RelayerConfigFilePath)
		s.Require().NoError(err)

		s.T().Cleanup(func() {
			os.Remove(testvalues.RelayerConfigFilePath)
		})
	}))

	s.T().Cleanup(func() {
		if relayerProcess != nil {
			err := relayerProcess.Kill()
			if err != nil {
				s.T().Logf("Failed to kill the relayer process: %v", err)
			}
		}
	})

	s.Require().True(s.Run("Create Relayer Client", func() {
		var err error
		s.RelayerClient, err = relayer.GetGRPCClient(relayer.DefaultRelayerGRPCAddress())
		s.Require().NoError(err)
	}))

	s.Require().True(s.Run("Deploy Groth16 ICS07 contract", func() {
		stdout, err := eth.ForgeScript(s.key, testvalues.E2EDeployScriptPath)
		s.Require().NoError(err)

		contractAddresses, err := ethereum.GetEthContractsFromDeployOutput(string(stdout))
		s.Require().NoError(err)

		var verfierAddress string
		if prover == testvalues.EnvValueGroth16Prover_Mock {
			verfierAddress = contractAddresses.VerifierMock
		} else {
			switch proofType {
			case types.ProofTypeGroth16:
				verfierAddress = contractAddresses.VerifierGroth16
			case types.ProofTypePlonk:
				verfierAddress = contractAddresses.VerifierPlonk
			default:
				s.Require().Fail("invalid proof type: %s", proofType)
			}
		}

		var createClientTxBz []byte
		s.Require().True(s.Run("Retrieve create client tx", func() {
			resp, err := s.RelayerClient.CreateClient(context.Background(), &relayertypes.CreateClientRequest{
				SrcChain: simd.Config().ChainID,
				DstChain: eth.ChainID.String(),
				Parameters: map[string]string{
					testvalues.ParameterKey_Groth16Verifier: verfierAddress,
					testvalues.ParameterKey_ZkAlgorithm:     proofType.String(),
					testvalues.ParameterKey_RoleManager:     ethcommon.Address{}.Hex(),
				},
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Tx)
			s.Require().Empty(resp.Address)

			createClientTxBz = resp.Tx
		}))

		s.Require().True(s.Run("Broadcast relay tx", func() {
			receipt, err := eth.BroadcastTx(ctx, s.key, 15_000_000, nil, createClientTxBz)
			s.Require().NoError(err)
			s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status, fmt.Sprintf("Tx failed: %+v", receipt))
			s.Require().NotEmpty(receipt.ContractAddress.Hex())
			s.groth16Ics07Address = receipt.ContractAddress

			s.contract, err = groth16ics07tendermint.NewContract(receipt.ContractAddress, eth.RPCClient)
			s.Require().NoError(err)
		}))

		s.Require().True(s.Run("Add client and counterparty on EVM", func() {
			counterpartyInfo := ics26router.IICS02ClientMsgsCounterpartyInfo{
				ClientId:     testvalues.FirstWasmClientID,
				MerklePrefix: [][]byte{[]byte(ibcexported.StoreKey), []byte("")},
			}
			tx, err := s.ics26Contract.AddClient(s.GetTransactOpts(s.key, eth), testvalues.CustomClientID, counterpartyInfo, s.groth16Ics07Address)
			s.Require().NoError(err)

			receipt, err := eth.GetTxReciept(ctx, tx.Hash())
			s.Require().NoError(err)

			event, err := e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseICS02ClientAdded)
			s.Require().NoError(err)
			s.Require().Equal(testvalues.CustomClientID, event.ClientId)
			s.Require().Equal(testvalues.FirstWasmClientID, event.CounterpartyInfo.ClientId)
		}))
	}))
}

// TestWithGroth16ICS07TendermintTestSuite is the boilerplate code that allows the test suite to be run
func TestWithGroth16ICS07TendermintTestSuite(t *testing.T) {
	suite.Run(t, new(Groth16ICS07TendermintTestSuite))
}

func (s *Groth16ICS07TendermintTestSuite) Test_Deploy() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.DeployTest(ctx, proofType)
}

// DeployTest tests the deployment of the Groth16ICS07Tendermint contract with the given arguments
func (s *Groth16ICS07TendermintTestSuite) DeployTest(ctx context.Context, proofType types.SupportedProofType) {
	s.SetupSuite(ctx, proofType)

	_, simd := s.EthChain, s.CosmosChains[0]

	s.Require().True(s.Run("Verify deployment", func() {
		clientState, err := getGroth16ClientState(s.contract)
		s.Require().NoError(err)

		stakingParams, err := simd.StakingQueryParams(ctx)
		s.Require().NoError(err)

		s.Require().Equal(simd.Config().ChainID, clientState.ChainId)
		s.Require().Equal(uint8(testvalues.DefaultTrustLevel.Numerator), clientState.TrustLevel.Numerator)
		s.Require().Equal(uint8(testvalues.DefaultTrustLevel.Denominator), clientState.TrustLevel.Denominator)
		s.Require().Equal(uint32(testvalues.DefaultTrustPeriod), clientState.TrustingPeriod)
		s.Require().Equal(uint32(stakingParams.UnbondingTime.Seconds()), clientState.UnbondingPeriod)
		s.Require().False(clientState.IsFrozen)
		s.Require().Equal(uint64(1), clientState.LatestHeight.RevisionNumber)
		s.Require().Greater(clientState.LatestHeight.RevisionHeight, uint64(0))
	}))
}

func (s *Groth16ICS07TendermintTestSuite) Test_UpdateClient() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.UpdateClientTest(ctx, proofType)
}

// UpdateClientTest tests the update client functionality
func (s *Groth16ICS07TendermintTestSuite) UpdateClientTest(ctx context.Context, proofType types.SupportedProofType) {
	s.SetupSuite(ctx, proofType)

	_, simd := s.EthChain, s.CosmosChains[0]

	if s.generateFixtures {
		s.T().Log("Generate fixtures is set to true, but TestUpdateClient does not support it (yet)")
	}

	s.Require().True(s.Run("Update client", func() {
		clientState, err := getGroth16ClientState(s.contract)
		s.Require().NoError(err)

		initialHeight := clientState.LatestHeight.RevisionHeight

		s.UpdateClient(ctx)

		clientState, err = getGroth16ClientState(s.contract)
		s.Require().NoError(err)

		stakingParams, err := simd.StakingQueryParams(ctx)
		s.Require().NoError(err)

		s.Require().Equal(simd.Config().ChainID, clientState.ChainId)
		s.Require().Equal(uint8(testvalues.DefaultTrustLevel.Numerator), clientState.TrustLevel.Numerator)
		s.Require().Equal(uint8(testvalues.DefaultTrustLevel.Denominator), clientState.TrustLevel.Denominator)
		s.Require().Equal(uint32(testvalues.DefaultTrustPeriod), clientState.TrustingPeriod)
		s.Require().Equal(uint32(stakingParams.UnbondingTime.Seconds()), clientState.UnbondingPeriod)
		s.Require().False(clientState.IsFrozen)
		s.Require().Equal(uint64(1), clientState.LatestHeight.RevisionNumber)
		s.Require().Greater(clientState.LatestHeight.RevisionHeight, initialHeight)
	}))
}

func (s *Groth16ICS07TendermintTestSuite) UpdateClient(ctx context.Context) clienttypes.Height {
	eth, simd := s.EthChain, s.CosmosChains[0]

	var initialHeight uint64
	s.Require().True(s.Run("Get the initial height", func() {
		clientState, err := getGroth16ClientState(s.contract)
		s.Require().NoError(err)
		s.Require().NotZero(clientState.LatestHeight.RevisionHeight)

		initialHeight = clientState.LatestHeight.RevisionHeight
	}))

	var finalHeight groth16ics07tendermint.IICS02ClientMsgsHeight
	s.Require().True(s.Run("Update the client on Ethereum", func() {
		var updateTxBodyBz []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.UpdateClient(context.Background(), &relayertypes.UpdateClientRequest{
				SrcChain:    simd.Config().ChainID,
				DstChain:    eth.ChainID.String(),
				DstClientId: testvalues.CustomClientID,
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Tx)
			s.Require().Equal(s.ics26Address.String(), resp.Address)

			updateTxBodyBz = resp.Tx
		}))

		s.Require().True(s.Run("Broadcast relay tx", func() {
			receipt, err := eth.BroadcastTx(ctx, s.key, 5_000_000, &s.ics26Address, updateTxBodyBz)
			s.Require().NoError(err)
			s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)
		}))

		s.Require().True(s.Run("Verify the client state is updated", func() {
			clientState, err := getGroth16ClientState(s.contract)
			s.Require().NoError(err)
			s.Require().NotZero(clientState.LatestHeight.RevisionHeight)

			finalHeight = clientState.LatestHeight
			s.Require().Greater(finalHeight.RevisionHeight, initialHeight)
		}))
	}))

	return clienttypes.NewHeight(finalHeight.RevisionNumber, finalHeight.RevisionHeight)
}

// TODO: Port Test_Membership, Test_UpdateClientAndMembership, and Test_DoubleSignMisbehaviour
// to work with new Groth16ICS07Tendermint contract ABI (struct fields changed from SP1 upstream).
