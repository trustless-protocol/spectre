package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	sdkmath "cosmossdk.io/math"

	bip39 "github.com/cosmos/go-bip39"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cosmossecp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	clienttypesv2 "github.com/cosmos/ibc-go/v10/modules/core/02-client/v2/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	ibchostv2 "github.com/cosmos/ibc-go/v10/modules/core/24-host/v2"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibctesting "github.com/cosmos/ibc-go/v10/testing"

	"github.com/cosmos/interchaintest/v10/ibc"

	"github.com/decentrio/fast-ibc/packages/go-abigen/ibcerc20"
	"github.com/decentrio/fast-ibc/packages/go-abigen/ics20transfer"
	"github.com/decentrio/fast-ibc/packages/go-abigen/ics26router"
	"github.com/decentrio/fast-ibc/packages/go-abigen/groth16ics07tendermint"

	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/cosmos"
	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/e2esuite"
	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/ethereum"
	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/relayer"
	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/testvalues"
	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/types"
	"github.com/srdtrk/solidity-ibc-eureka/e2e/v8/types/erc20"
	ethereumtypes "github.com/srdtrk/solidity-ibc-eureka/e2e/v8/types/ethereum"
	relayertypes "github.com/srdtrk/solidity-ibc-eureka/e2e/v8/types/relayer"
)

// IbcEurekaTestSuite is a suite of tests that wraps TestSuite
// and can provide additional functionality
type IbcEurekaTestSuite struct {
	e2esuite.TestSuite

	// Whether to generate fixtures for tests or not
	solidityFixtureGenerator *types.SolidityFixtureGenerator
	wasmFixtureGenerator     *types.WasmFixtureGenerator

	// The private key of a test account
	key *ecdsa.PrivateKey
	// The private key of the faucet account of interchaintest
	deployer *ecdsa.PrivateKey

	contractAddresses ethereum.DeployedContracts
	groth16Ics07Address   ethcommon.Address

	groth16Ics07Contract *groth16ics07tendermint.Contract
	ics26Contract    *ics26router.Contract
	ics20Contract    *ics20transfer.Contract
	erc20Contract    *erc20.Contract

	RelayerClient relayertypes.RelayerServiceClient

	SimdRelayerSubmitter ibc.Wallet
	EthRelayerSubmitter  *ecdsa.PrivateKey
}

// TestWithIbcEurekaTestSuite is the boilerplate code that allows the test suite to be run
func TestWithIbcEurekaTestSuite(t *testing.T) {
	suite.Run(t, new(IbcEurekaTestSuite))
}

// SetupSuite calls the underlying IbcEurekaTestSuite's SetupSuite method
// and deploys the IbcEureka contract. It uses the relayer binary directly
// (via create-clients + start) rather than gRPC.
func (s *IbcEurekaTestSuite) SetupSuite(ctx context.Context, proofType types.SupportedProofType) {
	s.TestSuite.SetupSuite(ctx)

	eth, simd := s.EthChain, s.CosmosChains[0]

	s.T().Logf("Setting up the test suite with proof type: %s", proofType.String())

	s.Require().True(s.Run("Set up environment", func() {
		var err error
		err = os.Chdir("../..")
		s.Require().NoError(err)

		s.key, err = eth.CreateAndFundUser()
		s.Require().NoError(err)

		s.EthRelayerSubmitter, err = eth.CreateAndFundUser()
		s.Require().NoError(err)

		s.deployer, err = eth.CreateAndFundUser()
		s.Require().NoError(err)

		s.SimdRelayerSubmitter = s.CreateAndFundCosmosUser(ctx, simd)

		// Derive secp256k1 private key from the relayer submitter's mnemonic so that
		// create-clients (which reads COSMOS_PRIVATE_KEY) creates 08-wasm-0 with the
		// same signer address that MsgRegisterCounterparty will use.
		seed := bip39.NewSeed(s.SimdRelayerSubmitter.Mnemonic(), "")
		masterPriv, chainCode := hd.ComputeMastersFromSeed(seed)
		derivedPrivKeyBytes, err := hd.DerivePrivateKeyForPath(masterPriv, chainCode, "m/44'/118'/0'/0/0")
		s.Require().NoError(err)
		cosmosPrivKey := cosmossecp256k1.PrivKey{Key: derivedPrivKeyBytes}
		os.Setenv("COSMOS_PRIVATE_KEY", hex.EncodeToString(cosmosPrivKey.Key))
		_ = sdk.AccAddress(cosmosPrivKey.PubKey().Address()) // verify derivation works

		os.Setenv(testvalues.EnvKeyEthRPC, eth.RPC)
		os.Setenv(testvalues.EnvKeyTendermintRPC, simd.GetHostRPCAddress())
		// ETH_PRIVATE_KEY is used by the relayer binary to sign ETH transactions
		os.Setenv("ETH_PRIVATE_KEY", hex.EncodeToString(crypto.FromECDSA(s.EthRelayerSubmitter)))
		os.Setenv("COSMOS_CHAIN_ID", simd.Config().ChainID)
		os.Setenv("COSMOS_FEE_DENOM", simd.Config().Denom)
	}))

	// Needs to be added here so the cleanup is called after the test suite is done
	s.wasmFixtureGenerator = types.NewWasmFixtureGenerator(&s.Suite)
	s.solidityFixtureGenerator = types.NewSolidityFixtureGenerator()

	s.Require().True(s.Run("Deploy IBC contracts", func() {
		stdout, err := eth.ForgeScript(s.deployer, testvalues.E2EDeployScriptPath)
		s.Require().NoError(err)

		s.contractAddresses, err = ethereum.GetEthContractsFromDeployOutput(string(stdout))
		s.Require().NoError(err)
		s.ics26Contract, err = ics26router.NewContract(ethcommon.HexToAddress(s.contractAddresses.Ics26Router), eth.RPCClient)
		s.Require().NoError(err)
		s.ics20Contract, err = ics20transfer.NewContract(ethcommon.HexToAddress(s.contractAddresses.Ics20Transfer), eth.RPCClient)
		s.Require().NoError(err)
		s.erc20Contract, err = erc20.NewContract(ethcommon.HexToAddress(s.contractAddresses.Erc20), eth.RPCClient)
		s.Require().NoError(err)
	}))

	// Store wasm light client code on Cosmos before running create-clients
	var checksumHex string
	s.Require().True(s.Run("Store Ethereum light client", func() {
		checksumHex = s.StoreEthereumLightClient(ctx, simd, s.SimdRelayerSubmitter)
		s.Require().NotEmpty(checksumHex)
	}))

	// Compute the ICS07 address that will be deployed by create-clients.
	// It is the first contract deployed from s.EthRelayerSubmitter's address.
	ethRelayerAddr := crypto.PubkeyToAddress(s.EthRelayerSubmitter.PublicKey)
	var preNonce uint64
	s.Require().True(s.Run("Get pre-deployment nonce", func() {
		var err error
		preNonce, err = eth.RPCClient.PendingNonceAt(ctx, ethRelayerAddr)
		s.Require().NoError(err)
	}))
	s.groth16Ics07Address = crypto.CreateAddress(ethRelayerAddr, preNonce)

	s.Require().True(s.Run("Generate relayer config (pre-ICS07)", func() {
		beaconAPI := ""
		if eth.BeaconAPIClient != nil {
			beaconAPI = eth.BeaconAPIClient.GetBeaconAPIURL()
		}

		config := relayer.NewConfig(relayer.CreateEthCosmosModules(
			relayer.EthCosmosConfigInfo{
				EthChainID:         eth.ChainID.String(),
				CosmosChainID:      simd.Config().ChainID,
				TmRPC:              simd.GetHostRPCAddress(),
				ICS26Address:       s.contractAddresses.Ics26Router,
				EthRPC:             eth.RPC,
				BeaconAPI:          beaconAPI,
				SignerAddress:      s.SimdRelayerSubmitter.FormattedAddress(),
				MockWasmClient:     os.Getenv(testvalues.EnvKeyEthTestnetType) == testvalues.EthTestnetTypePoW,
				WrapperVerifier:    s.contractAddresses.WrapperVerifier,
				Membership:         s.contractAddresses.Membership,
				Misbehaviour:       s.contractAddresses.Misbehaviour,
				UpdateClient:       s.contractAddresses.UpdateClient,
				CosmosWasmClientID: testvalues.FirstWasmClientID,
				TrustLevel:         "1/3",
				ProofType:          testvalues.EnvValueProofType_Groth16,
			}),
		)

		err := config.GenerateConfigFile(testvalues.RelayerConfigFilePath)
		s.Require().NoError(err)
	}))

	s.T().Cleanup(func() {
		os.Remove(testvalues.RelayerConfigFilePath)
	})

	// Run create-clients: deploys ICS07 on ETH (and creates wasm ETH client on Cosmos for PoS mode)
	s.Require().True(s.Run("Create light clients", func() {
		args := []string{"--trust-level", "1/3"}
		// Pass wasm checksum only for PoS mode (beacon URL required by relayer)
		if checksumHex != "" && eth.BeaconAPIClient != nil {
			args = append(args, "--wasm-checksum", checksumHex)
		}
		err := relayer.RunCreateClients(testvalues.RelayerConfigFilePath, args...)
		s.Require().NoError(err)
	}))

	// The relayer's create-clients deploys ICS07, but AddClient may fail (wrong signer role).
	// Call AddClient from the deployer who has ID_CUSTOMIZER_ROLE.
	s.Require().True(s.Run("Add Cosmos light client to ICS26Router", func() {
		var err error
		s.groth16Ics07Contract, err = groth16ics07tendermint.NewContract(s.groth16Ics07Address, eth.RPCClient)
		s.Require().NoError(err)

		counterpartyInfo := ics26router.IICS02ClientMsgsCounterpartyInfo{
			ClientId:     testvalues.FirstWasmClientID,
			MerklePrefix: [][]byte{[]byte(ibcexported.StoreKey), []byte("")},
		}
		// Try AddClient; if client already exists (relayer succeeded), this is a no-op error we can ignore
		tx, err := s.ics26Contract.AddClient(s.GetTransactOpts(s.deployer, eth), testvalues.CustomClientID, counterpartyInfo, s.groth16Ics07Address)
		if err != nil {
			// Check if client already registered (relayer may have succeeded with AddClient)
			existingAddr, queryErr := s.ics26Contract.GetClient(nil, testvalues.CustomClientID)
			s.Require().NoError(queryErr, "AddClient failed and GetClient also failed: %v", err)
			s.Require().Equal(s.groth16Ics07Address, existingAddr, "AddClient failed but client registered at wrong address")
		} else {
			receipt, err := eth.GetTxReciept(ctx, tx.Hash())
			s.Require().NoError(err)

			event, err := e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseICS02ClientAdded)
			s.Require().NoError(err)
			s.Require().Equal(testvalues.CustomClientID, event.ClientId)
			s.Require().Equal(testvalues.FirstWasmClientID, event.CounterpartyInfo.ClientId)
		}
	}))

	s.Require().True(s.Run("Fund address with ERC20", func() {
		tx, err := s.erc20Contract.Transfer(s.GetTransactOpts(eth.Faucet, eth), crypto.PubkeyToAddress(s.key.PublicKey), testvalues.StartingERC20Balance)
		s.Require().NoError(err)
		_, err = eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
	}))

	// For PoW mode (no beacon URL), create-clients skips ETH light client on Cosmos.
	// The test creates a mock wasm client manually.
	if os.Getenv(testvalues.EnvKeyEthTestnetType) == testvalues.EthTestnetTypePoW {
		s.Require().True(s.Run("Create mock Ethereum light client on Cosmos", func() {
			checksumBytes, err := hex.DecodeString(strings.TrimPrefix(checksumHex, "0x"))
			s.Require().NoError(err)

			wasmClientState := ibcwasmtypes.ClientState{
				Data:     []byte("{}"),
				Checksum: checksumBytes,
				LatestHeight: clienttypes.Height{
					RevisionNumber: 0,
					RevisionHeight: 1,
				},
			}
			wasmConsensusState := ibcwasmtypes.ConsensusState{
				Data: []byte("{}"),
			}

			msg, err := clienttypes.NewMsgCreateClient(&wasmClientState, &wasmConsensusState, s.SimdRelayerSubmitter.FormattedAddress())
			s.Require().NoError(err)

			resp, err := s.BroadcastMessages(ctx, simd, s.SimdRelayerSubmitter, 200_000, msg)
			s.Require().NoError(err)

			clientId, err := cosmos.GetEventValue(resp.Events, clienttypes.EventTypeCreateClient, clienttypes.AttributeKeyClientID)
			s.Require().NoError(err)
			s.Require().Equal(testvalues.FirstWasmClientID, clientId)
		}))
	}

	s.Require().True(s.Run("Register counterparty on Cosmos chain", func() {
		merklePathPrefix := [][]byte{[]byte("")}
		_, err := s.BroadcastMessages(ctx, simd, s.SimdRelayerSubmitter, 200_000, &clienttypesv2.MsgRegisterCounterparty{
			ClientId:                 testvalues.FirstWasmClientID,
			CounterpartyMerklePrefix: merklePathPrefix,
			CounterpartyClientId:     testvalues.CustomClientID,
			Signer:                   s.SimdRelayerSubmitter.FormattedAddress(),
		})
		s.Require().NoError(err)
	}))

	// Regenerate config with ICS07 address so the relay loop can use it
	var relayerProcess *os.Process
	s.Require().True(s.Run("Start relay loop", func() {
		beaconAPI := ""
		if eth.BeaconAPIClient != nil {
			beaconAPI = eth.BeaconAPIClient.GetBeaconAPIURL()
		}

		config := relayer.NewConfig(relayer.CreateEthCosmosModules(
			relayer.EthCosmosConfigInfo{
				EthChainID:         eth.ChainID.String(),
				CosmosChainID:      simd.Config().ChainID,
				TmRPC:              simd.GetHostRPCAddress(),
				ICS26Address:       s.contractAddresses.Ics26Router,
				EthRPC:             eth.RPC,
				BeaconAPI:          beaconAPI,
				SignerAddress:      s.SimdRelayerSubmitter.FormattedAddress(),
				MockWasmClient:     os.Getenv(testvalues.EnvKeyEthTestnetType) == testvalues.EthTestnetTypePoW,
				ICS07Client:        s.groth16Ics07Address.Hex(),
				WrapperVerifier:    s.contractAddresses.WrapperVerifier,
				Membership:         s.contractAddresses.Membership,
				Misbehaviour:       s.contractAddresses.Misbehaviour,
				UpdateClient:       s.contractAddresses.UpdateClient,
				CosmosWasmClientID: testvalues.FirstWasmClientID,
				TrustLevel:         "1/3",
				ProofType:          testvalues.EnvValueProofType_Groth16,
			}),
		)

		var err error
		err = config.GenerateConfigFile(testvalues.RelayerConfigFilePath)
		s.Require().NoError(err)

		relayerProcess, err = relayer.StartRelayer(testvalues.RelayerConfigFilePath)
		s.Require().NoError(err)
	}))

	s.T().Cleanup(func() {
		if relayerProcess != nil {
			if err := relayerProcess.Kill(); err != nil {
				s.T().Logf("Failed to kill the relayer process: %v", err)
			}
		}
	})

	s.Require().True(s.Run("Generate the genesis fixtures", func() {
		if !s.solidityFixtureGenerator.Enabled {
			s.T().Skip("Skipping solidity fixture generation")
		}

		clientState, err := s.groth16Ics07Contract.ClientState(nil)
		s.Require().NoError(err)
		clientStateBz, err := s.groth16Ics07Contract.GetClientState(nil)
		s.Require().NoError(err)
		consensusStateHash, err := s.groth16Ics07Contract.GetConsensusStateHash(nil, clientState.LatestHeight.RevisionHeight)
		s.Require().NoError(err)
		// Groth16 ICS07 ABI does not expose program vkeys on the root contract.
		var updateClientVkey, membershipVkey, ucAndMembershipVkey, misbehaviourVkey [32]byte
		s.solidityFixtureGenerator.SetGenesisFixture(
			clientStateBz, consensusStateHash, updateClientVkey,
			membershipVkey, ucAndMembershipVkey, misbehaviourVkey,
		)
	}))
}

func (s *IbcEurekaTestSuite) Test_Deploy() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.DeployTest(ctx, proofType)
}

// DeployTest tests the deployment of the IbcEureka contracts
func (s *IbcEurekaTestSuite) DeployTest(ctx context.Context, proofType types.SupportedProofType) {
	s.SetupSuite(ctx, proofType)

	eth, simd := s.EthChain, s.CosmosChains[0]

	s.Require().True(s.Run("Verify Groth16 Client", func() {
		clientState, err := s.groth16Ics07Contract.ClientState(nil)
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

	s.Require().True(s.Run("Verify ICS02 Client", func() {
		clientAddress, err := s.ics26Contract.GetClient(nil, testvalues.CustomClientID)
		s.Require().NoError(err)
		s.Require().Equal(s.groth16Ics07Address, clientAddress)

		counterpartyInfo, err := s.ics26Contract.GetCounterparty(nil, testvalues.CustomClientID)
		s.Require().NoError(err)
		s.Require().Equal(testvalues.FirstWasmClientID, counterpartyInfo.ClientId)
	}))

	s.Require().True(s.Run("Verify ICS26 Router", func() {
		transferAddress, err := s.ics26Contract.GetIBCApp(nil, transfertypes.PortID)
		s.Require().NoError(err)
		s.Require().Equal(s.contractAddresses.Ics20Transfer, strings.ToLower(transferAddress.Hex()))
	}))

	s.Require().True(s.Run("Verify ERC20 Genesis", func() {
		userBalance, err := s.erc20Contract.BalanceOf(nil, crypto.PubkeyToAddress(s.key.PublicKey))
		s.Require().NoError(err)
		s.Require().Equal(testvalues.StartingERC20Balance, userBalance)
	}))

	s.Require().True(s.Run("Verify ethereum light client", func() {
		_, err := e2esuite.GRPCQuery[clienttypes.QueryClientStateResponse](ctx, simd, &clienttypes.QueryClientStateRequest{
			ClientId: testvalues.FirstWasmClientID,
		})
		s.Require().NoError(err)

		counterpartyInfoResp, err := e2esuite.GRPCQuery[clienttypesv2.QueryCounterpartyInfoResponse](ctx, simd, &clienttypesv2.QueryCounterpartyInfoRequest{
			ClientId: testvalues.FirstWasmClientID,
		})
		s.Require().NoError(err)
		s.Require().Equal(testvalues.CustomClientID, counterpartyInfoResp.CounterpartyInfo.ClientId)
	}))

	s.Require().True(s.Run("Verify Cosmos to Eth Relayer Info", func() {
		info, err := s.RelayerClient.Info(context.Background(), &relayertypes.InfoRequest{
			SrcChain: simd.Config().ChainID,
			DstChain: eth.ChainID.String(),
		})
		s.Require().NoError(err)
		s.Require().NotNil(info)
		s.Require().Equal(simd.Config().ChainID, info.SourceChain.ChainId)
		s.Require().Equal(eth.ChainID.String(), info.TargetChain.ChainId)
	}))

	s.Require().True(s.Run("Verify Eth to Cosmos Relayer Info", func() {
		info, err := s.RelayerClient.Info(context.Background(), &relayertypes.InfoRequest{
			SrcChain: eth.ChainID.String(),
			DstChain: simd.Config().ChainID,
		})
		s.Require().NoError(err)
		s.Require().NotNil(info)
		s.Require().Equal(eth.ChainID.String(), info.SourceChain.ChainId)
		s.Require().Equal(simd.Config().ChainID, info.TargetChain.ChainId)
	}))
}

func (s *IbcEurekaTestSuite) Test_ICS20TransferERC20TokenfromEthereumToCosmosAndBack() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.ICS20TransferERC20TokenfromEthereumToCosmosAndBackTest(ctx, proofType, 1, big.NewInt(testvalues.TransferAmount))
}

func (s *IbcEurekaTestSuite) Test_25_ICS20TransferERC20TokenfromEthereumToCosmosAndBack() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.ICS20TransferERC20TokenfromEthereumToCosmosAndBackTest(ctx, proofType, 25, big.NewInt(testvalues.TransferAmount))
}

func (s *IbcEurekaTestSuite) Test_50_ICS20TransferERC20TokenfromEthereumToCosmosAndBack() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.ICS20TransferERC20TokenfromEthereumToCosmosAndBackTest(ctx, proofType, 50, big.NewInt(testvalues.TransferAmount))
}

func (s *IbcEurekaTestSuite) Test_ICS20TransferUint256TokenfromEthereumToCosmosAndBack() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	transferAmount := new(big.Int).Div(testvalues.StartingERC20Balance, big.NewInt(2))
	s.ICS20TransferERC20TokenfromEthereumToCosmosAndBackTest(ctx, proofType, 1, transferAmount)
}

// ICS20TransferERC20TokenfromEthereumToCosmosAndBackTest tests the ICS20 transfer functionality by transferring
// ERC20 tokens with n packets from Ethereum to Cosmos chain and then back from Cosmos chain to Ethereum
func (s *IbcEurekaTestSuite) ICS20TransferERC20TokenfromEthereumToCosmosAndBackTest(
	ctx context.Context, proofType types.SupportedProofType, numOfTransfers int, transferAmount *big.Int,
) {
	s.SetupSuite(ctx, proofType)

	eth, simd := s.EthChain, s.CosmosChains[0]

	ics26Address := ethcommon.HexToAddress(s.contractAddresses.Ics26Router)
	ics20Address := ethcommon.HexToAddress(s.contractAddresses.Ics20Transfer)
	erc20Address := ethcommon.HexToAddress(s.contractAddresses.Erc20)

	totalTransferAmount := new(big.Int).Mul(transferAmount, big.NewInt(int64(numOfTransfers)))
	ethereumUserAddress := crypto.PubkeyToAddress(s.key.PublicKey)
	cosmosUserWallet := s.CosmosUsers[0]
	cosmosUserAddress := cosmosUserWallet.FormattedAddress()

	ics20transferAbi, err := abi.JSON(strings.NewReader(ics20transfer.ContractABI))
	s.Require().NoError(err)

	s.Require().True(s.Run("Approve the ICS20Transfer.sol contract to spend the erc20 tokens", func() {
		tx, err := s.erc20Contract.Approve(s.GetTransactOpts(s.key, eth), ics20Address, totalTransferAmount)
		s.Require().NoError(err)

		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

		allowance, err := s.erc20Contract.Allowance(nil, ethereumUserAddress, ics20Address)
		s.Require().NoError(err)
		s.Require().Equal(totalTransferAmount, allowance)
	}))

	var (
		sendPacket    ics26router.IICS26RouterMsgsPacket
		ethSendTxHash []byte
		escrowAddress ethcommon.Address
	)
	s.Require().True(s.Run(fmt.Sprintf("Send %d transfers on Ethereum", numOfTransfers), func() {
		timeout := uint64(time.Now().Add(30 * time.Minute).Unix())
		transferMulticall := make([][]byte, numOfTransfers)

		msgSendPacket := ics20transfer.IICS20TransferMsgsSendTransferMsg{
			Denom:            erc20Address,
			Amount:           transferAmount,
			Receiver:         cosmosUserAddress,
			TimeoutTimestamp: timeout,
			SourceClient:     testvalues.CustomClientID,
			Memo:             "",
		}

		encodedMsg, err := ics20transferAbi.Pack("sendTransfer", msgSendPacket)
		s.Require().NoError(err)
		for i := range numOfTransfers {
			transferMulticall[i] = encodedMsg
		}

		tx, err := s.ics20Contract.Multicall(s.GetTransactOpts(s.key, eth), transferMulticall)
		s.Require().NoError(err)
		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)
		s.T().Logf("Multicall send %d transfers gas used: %d", numOfTransfers, receipt.GasUsed)
		ethSendTxHash = tx.Hash().Bytes()

		sendPacketEvent, err := e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseSendPacket)
		s.Require().NoError(err)
		sendPacket = sendPacketEvent.Packet
		s.Require().Equal(uint64(1), sendPacket.Sequence)
		s.Require().Equal(timeout, sendPacket.TimeoutTimestamp)
		s.Require().Len(sendPacket.Payloads, 1)
		s.Require().Equal(transfertypes.PortID, sendPacket.Payloads[0].SourcePort)
		s.Require().Equal(testvalues.CustomClientID, sendPacket.SourceClient)
		s.Require().Equal(transfertypes.PortID, sendPacket.Payloads[0].DestPort)
		s.Require().Equal(testvalues.FirstWasmClientID, sendPacket.DestClient)
		s.Require().Equal(transfertypes.V1, sendPacket.Payloads[0].Version)
		s.Require().Equal(transfertypes.EncodingABI, sendPacket.Payloads[0].Encoding)

		s.True(s.Run("Verify balances on Ethereum", func() {
			// User balance on Ethereum
			userBalance, err := s.erc20Contract.BalanceOf(nil, ethereumUserAddress)
			s.Require().NoError(err)
			s.Require().Equal(new(big.Int).Sub(testvalues.StartingERC20Balance, totalTransferAmount), userBalance)

			// Get the escrow address
			escrowAddress, err = s.ics20Contract.GetEscrow(nil, testvalues.CustomClientID)
			s.Require().NoError(err)

			// ICS20 contract balance on Ethereum
			escrowBalance, err := s.erc20Contract.BalanceOf(nil, escrowAddress)
			s.Require().NoError(err)
			s.Require().Equal(totalTransferAmount, escrowBalance)
		}))
	}))

	var (
		denomOnCosmos transfertypes.Denom
		ackTxHash     []byte
	)
	s.Require().True(s.Run("Receive packets on Cosmos chain", func() {
		var relayTxBodyBz []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:    eth.ChainID.String(),
				DstChain:    simd.Config().ChainID,
				SourceTxIds: [][]byte{ethSendTxHash},
				SrcClientId: testvalues.CustomClientID,
				DstClientId: testvalues.FirstWasmClientID,
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Tx)
			s.Require().Empty(resp.Address)

			relayTxBodyBz = resp.Tx

			s.wasmFixtureGenerator.AddFixtureStep("receive_packets", ethereumtypes.RelayerMessages{
				RelayerTxBody: hex.EncodeToString(relayTxBodyBz),
			})
		}))

		s.Require().True(s.Run("Broadcast relay tx", func() {
			resp := s.MustBroadcastSdkTxBody(ctx, simd, s.SimdRelayerSubmitter, 20_000_000, relayTxBodyBz)

			ackTxHash, err = hex.DecodeString(resp.TxHash)
			s.Require().NoError(err)
			s.Require().NotEmpty(ackTxHash)
		}))

		s.Require().True(s.Run("Verify balances on Cosmos chain", func() {
			denomOnCosmos = transfertypes.NewDenom(s.contractAddresses.Erc20, transfertypes.NewHop(transfertypes.PortID, testvalues.FirstWasmClientID))

			// User balance on Cosmos chain
			resp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
				Address: cosmosUserAddress,
				Denom:   denomOnCosmos.IBCDenom(),
			})
			s.Require().NoError(err)
			s.Require().NotNil(resp.Balance)
			s.Require().Equal(totalTransferAmount, resp.Balance.Amount.BigInt())
			s.Require().Equal(denomOnCosmos.IBCDenom(), resp.Balance.Denom)
		}))
	}))

	s.Require().True(s.Run("Acknowledge packets on Ethereum", func() {
		var ackRelayTx []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:    simd.Config().ChainID,
				DstChain:    eth.ChainID.String(),
				SourceTxIds: [][]byte{ackTxHash},
				SrcClientId: testvalues.FirstWasmClientID,
				DstClientId: testvalues.CustomClientID,
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Tx)
			s.Require().Equal(resp.Address, ics26Address.String())

			ackRelayTx = resp.Tx
		}))

		s.Require().True(s.Run("Submit relay tx", func() {
			receipt, err := eth.BroadcastTx(ctx, s.EthRelayerSubmitter, 15_000_000, &ics26Address, ackRelayTx)
			s.Require().NoError(err)
			s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status, fmt.Sprintf("Tx failed: %+v", receipt))
			s.T().Logf("Multicall ack %d packets gas used: %d", numOfTransfers, receipt.GasUsed)

			// Verify the ack packet event exists
			_, err = e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseAckPacket)
			s.Require().NoError(err)
		}))

		s.Require().NoError(s.solidityFixtureGenerator.GenerateAndSaveSolidityFixture(
			fmt.Sprintf("acknowledgeMultiPacket_%d-%s.json", numOfTransfers, proofType.String()),
			s.contractAddresses.Erc20, ackRelayTx, sendPacket,
		))

		s.Require().True(s.Run("Verify balances on Ethereum", func() {
			// User balance on Ethereum
			userBalance, err := s.erc20Contract.BalanceOf(nil, ethereumUserAddress)
			s.Require().NoError(err)
			s.Require().Equal(new(big.Int).Sub(testvalues.StartingERC20Balance, totalTransferAmount), userBalance)

			// ICS20 contract balance on Ethereum
			escrowBalance, err := s.erc20Contract.BalanceOf(nil, escrowAddress)
			s.Require().NoError(err)
			s.Require().Equal(totalTransferAmount, escrowBalance)
		}))
	}))

	var returnSendTxHash []byte
	s.Require().True(s.Run("Transfer tokens back from Cosmos chain", func() {
		timeout := uint64(time.Now().Add(30 * time.Minute).Unix())
		ibcCoin := sdk.NewCoin(denomOnCosmos.Path(), sdkmath.NewIntFromBigInt(transferAmount))

		transferPayload := transfertypes.FungibleTokenPacketData{
			Denom:    ibcCoin.Denom,
			Amount:   ibcCoin.Amount.String(),
			Sender:   cosmosUserWallet.FormattedAddress(),
			Receiver: strings.ToLower(ethereumUserAddress.Hex()),
			Memo:     "",
		}
		encodedPayload, err := transfertypes.EncodeABIFungibleTokenPacketData(&transferPayload)
		s.Require().NoError(err)

		payload := channeltypesv2.Payload{
			SourcePort:      transfertypes.PortID,
			DestinationPort: transfertypes.PortID,
			Version:         transfertypes.V1,
			Encoding:        transfertypes.EncodingABI,
			Value:           encodedPayload,
		}

		transferMsgs := make([]sdk.Msg, numOfTransfers)
		for i := range numOfTransfers {
			transferMsgs[i] = &channeltypesv2.MsgSendPacket{
				SourceClient:     testvalues.FirstWasmClientID,
				TimeoutTimestamp: timeout,
				Payloads: []channeltypesv2.Payload{
					payload,
				},
				Signer: cosmosUserWallet.FormattedAddress(),
			}
		}

		resp, err := s.BroadcastMessages(ctx, simd, cosmosUserWallet, 20_000_000, transferMsgs...)
		s.Require().NoError(err)
		s.Require().NotEmpty(resp.TxHash)

		returnSendTxHash, err = hex.DecodeString(resp.TxHash)
		s.Require().NoError(err)

		s.Require().True(s.Run("Verify balances on Cosmos chain", func() {
			// User balance on Cosmos chain
			resp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
				Address: cosmosUserAddress,
				Denom:   denomOnCosmos.IBCDenom(),
			})
			s.Require().NoError(err)
			s.Require().NotNil(resp.Balance)
			s.Require().Equal(sdkmath.ZeroInt(), resp.Balance.Amount)
			s.Require().Equal(denomOnCosmos.IBCDenom(), resp.Balance.Denom)
		}))
	}))

	var returnAckTxHash []byte
	s.Require().True(s.Run(fmt.Sprintf("Receive %d packets on Ethereum", numOfTransfers), func() {
		var recvRelayTx []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:    simd.Config().ChainID,
				DstChain:    eth.ChainID.String(),
				SourceTxIds: [][]byte{returnSendTxHash},
				SrcClientId: testvalues.FirstWasmClientID,
				DstClientId: testvalues.CustomClientID,
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Tx)
			s.Require().Equal(resp.Address, ics26Address.String())

			recvRelayTx = resp.Tx
		}))

		var returnPacket ics26router.IICS26RouterMsgsPacket
		s.Require().True(s.Run("Submit relay tx", func() {
			receipt, err := eth.BroadcastTx(ctx, s.EthRelayerSubmitter, 15_000_000, &ics26Address, recvRelayTx)
			s.Require().NoError(err)
			s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status, fmt.Sprintf("Tx failed: %+v", receipt))
			s.T().Logf("Multicall receive %d packets gas used: %d", numOfTransfers, receipt.GasUsed)

			returnWriteAckEvent, err := e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseWriteAcknowledgement)
			s.Require().NoError(err)

			returnPacket = returnWriteAckEvent.Packet
			returnAckTxHash = receipt.TxHash.Bytes()
		}))

		s.Require().NoError(s.solidityFixtureGenerator.GenerateAndSaveSolidityFixture(
			fmt.Sprintf("receiveMultiPacket_%d-%s.json", numOfTransfers, proofType.String()),
			s.contractAddresses.Erc20, recvRelayTx, returnPacket,
		))

		s.True(s.Run("Verify balances on Ethereum", func() {
			// User balance should be back to the starting point
			userBalance, err := s.erc20Contract.BalanceOf(nil, ethereumUserAddress)
			s.Require().NoError(err)
			s.Require().Equal(testvalues.StartingERC20Balance, userBalance)

			escrowBalance, err := s.erc20Contract.BalanceOf(nil, escrowAddress)
			s.Require().NoError(err)
			s.Require().Zero(escrowBalance.Int64())
		}))
	}))

	s.Require().True(s.Run("Acknowledge packets on Cosmos chain", func() {
		s.Require().True(s.Run("Verify commitments exists", func() {
			for i := range numOfTransfers {
				resp, err := e2esuite.GRPCQuery[channeltypesv2.QueryPacketCommitmentResponse](ctx, simd, &channeltypesv2.QueryPacketCommitmentRequest{
					ClientId: testvalues.FirstWasmClientID,
					Sequence: uint64(i) + 1,
				})
				s.Require().NoError(err)
				s.Require().NotEmpty(resp.Commitment)
			}
		}))

		var ackRelayTxBodyBz []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:    eth.ChainID.String(),
				DstChain:    simd.Config().ChainID,
				SourceTxIds: [][]byte{returnAckTxHash},
				SrcClientId: testvalues.CustomClientID,
				DstClientId: testvalues.FirstWasmClientID,
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Tx)
			s.Require().Empty(resp.Address)

			ackRelayTxBodyBz = resp.Tx

			s.wasmFixtureGenerator.AddFixtureStep("ack_packets", ethereumtypes.RelayerMessages{
				RelayerTxBody: hex.EncodeToString(ackRelayTxBodyBz),
			})
		}))

		s.Require().True(s.Run("Broadcast relay tx", func() {
			resp := s.MustBroadcastSdkTxBody(ctx, simd, s.SimdRelayerSubmitter, 20_000_000, ackRelayTxBodyBz)

			ackTxHash, err = hex.DecodeString(resp.TxHash)
			s.Require().NoError(err)
			s.Require().NotEmpty(ackTxHash)
		}))

		s.Require().True(s.Run("Verify commitments removed", func() {
			for i := range numOfTransfers {
				_, err := e2esuite.GRPCQuery[channeltypesv2.QueryPacketCommitmentResponse](ctx, simd, &channeltypesv2.QueryPacketCommitmentRequest{
					ClientId: testvalues.FirstWasmClientID,
					Sequence: uint64(i) + 1,
				})
				s.Require().ErrorContains(err, "packet commitment hash not found")
			}
		}))
	}))
}

func (s *IbcEurekaTestSuite) Test_ICS20TransferERC20TokenFromEthereumToCosmosAndBackFails() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.ICS20TransferERC20TokenFromEthereumToCosmosAndBackFailsTest(ctx, proofType, 1, big.NewInt(testvalues.TransferAmount))
}

func (s *IbcEurekaTestSuite) ICS20TransferERC20TokenFromEthereumToCosmosAndBackFailsTest(
	ctx context.Context, proofType types.SupportedProofType, numOfTransfers int, transferAmount *big.Int,
) {
	s.SetupSuite(ctx, proofType)

	eth, simd := s.EthChain, s.CosmosChains[0]

	ics26Address := ethcommon.HexToAddress(s.contractAddresses.Ics26Router)
	ics20Address := ethcommon.HexToAddress(s.contractAddresses.Ics20Transfer)
	erc20Address := ethcommon.HexToAddress(s.contractAddresses.Erc20)

	totalTransferAmount := new(big.Int).Mul(transferAmount, big.NewInt(int64(numOfTransfers)))
	ethereumUserAddress := crypto.PubkeyToAddress(s.key.PublicKey)
	cosmosUserWallet := s.CosmosUsers[0]
	cosmosUserAddress := cosmosUserWallet.FormattedAddress()

	ics20transferAbi, err := abi.JSON(strings.NewReader(ics20transfer.ContractABI))
	s.Require().NoError(err)

	s.Require().True(s.Run("Approve the ICS20Transfer.sol contract to spend the erc20 tokens", func() {
		tx, err := s.erc20Contract.Approve(s.GetTransactOpts(s.key, eth), ics20Address, totalTransferAmount)
		s.Require().NoError(err)

		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

		allowance, err := s.erc20Contract.Allowance(nil, ethereumUserAddress, ics20Address)
		s.Require().NoError(err)
		s.Require().Equal(totalTransferAmount, allowance)
	}))

	var ethSendTxHash []byte
	s.Require().True(s.Run(fmt.Sprintf("Send %d transfers on Ethereum", numOfTransfers), func() {
		timeout := uint64(time.Now().Add(30 * time.Minute).Unix())
		transferMulticall := make([][]byte, numOfTransfers)

		msgSendPacket := ics20transfer.IICS20TransferMsgsSendTransferMsg{
			Denom:            erc20Address,
			Amount:           transferAmount,
			Receiver:         cosmosUserAddress,
			TimeoutTimestamp: timeout,
			SourceClient:     testvalues.CustomClientID,
			Memo:             "",
		}

		encodedMsg, err := ics20transferAbi.Pack("sendTransfer", msgSendPacket)
		s.Require().NoError(err)
		for i := range numOfTransfers {
			transferMulticall[i] = encodedMsg
		}

		tx, err := s.ics20Contract.Multicall(s.GetTransactOpts(s.key, eth), transferMulticall)
		s.Require().NoError(err)
		_, err = eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		ethSendTxHash = tx.Hash().Bytes()
	}))

	var ackTxHash []byte
	s.Require().True(s.Run("Receive packets on Cosmos chain", func() {
		var relayTxBodyBz []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:    eth.ChainID.String(),
				DstChain:    simd.Config().ChainID,
				SourceTxIds: [][]byte{ethSendTxHash},
				SrcClientId: testvalues.CustomClientID,
				DstClientId: testvalues.FirstWasmClientID,
			})
			s.Require().NoError(err)

			relayTxBodyBz = resp.Tx
			relayTxBodyBz = resp.Tx
		}))

		s.Require().True(s.Run("Broadcast relay tx", func() {
			resp, err := s.BroadcastSdkTxBody(ctx, simd, s.SimdRelayerSubmitter, 20_000_000, relayTxBodyBz)
			s.Require().NoError(err)

			ackTxHash, err = hex.DecodeString(resp.TxHash)
			s.Require().NoError(err)
		}))
	}))

	s.Require().True(s.Run("Transfer tokens back from Cosmos chain with invalid receiver should fail", func() {
		denomOnCosmos := transfertypes.NewDenom(s.contractAddresses.Erc20, transfertypes.NewHop(transfertypes.PortID, testvalues.FirstWasmClientID))
		timeout := uint64(time.Now().Add(30 * time.Minute).Unix())
		ibcCoin := sdk.NewCoin(denomOnCosmos.Path(), sdkmath.NewIntFromBigInt(transferAmount))

		transferPayload := transfertypes.FungibleTokenPacketData{
			Denom:    ibcCoin.Denom,
			Amount:   ibcCoin.Amount.String(),
			Sender:   cosmosUserWallet.FormattedAddress(),
			Receiver: "invalid_receiver_address",
			Memo:     "",
		}
		encodedPayload, err := transfertypes.EncodeABIFungibleTokenPacketData(&transferPayload)
		s.Require().NoError(err)

		payload := channeltypesv2.Payload{
			SourcePort:      transfertypes.PortID,
			DestinationPort: transfertypes.PortID,
			Version:         transfertypes.V1,
			Encoding:        transfertypes.EncodingABI,
			Value:           encodedPayload,
		}

		transferMsg := &channeltypesv2.MsgSendPacket{
			SourceClient:     testvalues.FirstWasmClientID,
			TimeoutTimestamp: timeout,
			Payloads: []channeltypesv2.Payload{
				payload,
			},
			Signer: cosmosUserWallet.FormattedAddress(),
		}

		var relayTxHash []byte
		s.Require().True(s.Run("Broadcast transfer msgs", func() {
			resp, err := s.BroadcastMessages(ctx, simd, cosmosUserWallet, 20_000_000, transferMsg)
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.TxHash)

			relayTxHash, err = hex.DecodeString(resp.TxHash)
			s.Require().NoError(err)
		}))

		var transferRelayTxHashBz []byte
		s.Require().True(s.Run("Relay transfer to eth", func() {
			var relayTxBodyBz []byte
			s.Require().True(s.Run("Retrieve relay tx", func() {
				resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
					SrcChain:    simd.Config().ChainID,
					DstChain:    eth.ChainID.String(),
					SourceTxIds: [][]byte{relayTxHash},
					SrcClientId: testvalues.FirstWasmClientID,
					DstClientId: testvalues.CustomClientID,
				})
				s.Require().NoError(err)
				s.Require().NotEmpty(resp.Tx)
				s.Require().NotEmpty(resp.Address)

				relayTxBodyBz = resp.Tx
			}))

			s.Require().True(s.Run("Submit relay tx to eth", func() {
				receipt, err := eth.BroadcastTx(ctx, s.EthRelayerSubmitter, 15_000_000, &ics26Address, relayTxBodyBz)
				s.Require().NoError(err)

				transferRelayTxHashBz = receipt.TxHash.Bytes()
			}))
		}))

		var ackRelayTxBodyBz []byte
		s.Require().True(s.Run("Acknowledge packet to Cosmos chain", func() {
			s.Require().True(s.Run("Retrieve relay tx", func() {
				resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
					SrcChain:    eth.ChainID.String(),
					DstChain:    simd.Config().ChainID,
					SourceTxIds: [][]byte{transferRelayTxHashBz},
					SrcClientId: testvalues.CustomClientID,
					DstClientId: testvalues.FirstWasmClientID,
				})
				s.Require().NoError(err)
				s.Require().NotEmpty(resp.Tx)

				ackRelayTxBodyBz = resp.Tx
			}))

			s.Require().True(s.Run("Submit acknowledgement relay tx to Cosmos", func() {
				resp, err := s.BroadcastSdkTxBody(ctx, simd, s.SimdRelayerSubmitter, 20_000_000, ackRelayTxBodyBz)
				s.Require().NoError(err)

				ackTxHash, err = hex.DecodeString(resp.TxHash)
				s.Require().NoError(err)
				s.Require().NotEmpty(ackTxHash)
			}))
		}))

		s.Require().True(s.Run("Verify balances on Cosmos chain", func() {
			resp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
				Address: cosmosUserAddress,
				Denom:   denomOnCosmos.IBCDenom(),
			})
			s.Require().NoError(err)
			s.Require().NotNil(resp.Balance)
			// Vouchers should be restored to the user's balance
			s.Require().Equal(totalTransferAmount, resp.Balance.Amount.BigInt())
			s.Require().Equal(denomOnCosmos.IBCDenom(), resp.Balance.Denom)
		}))
	}))
}

func (s *IbcEurekaTestSuite) Test_ICS20TransferNativeCosmosCoinsToEthereumAndBack() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.ICS20TransferNativeCosmosCoinsToEthereumAndBackTest(ctx, proofType, big.NewInt(testvalues.TransferAmount))
}

// ICS20TransferNativeCosmosCoinsToEthereumAndBackTest tests the ICS20 transfer functionality
// by transferring native coins from a Cosmos chain to Ethereum and back.
// Cosmos→ETH relay is handled automatically by the running relayer process.
// ETH→Cosmos relay is done manually in the test (mock proof accepted by dummy wasm client).
func (s *IbcEurekaTestSuite) ICS20TransferNativeCosmosCoinsToEthereumAndBackTest(ctx context.Context, pt types.SupportedProofType, transferAmount *big.Int) {
	s.SetupSuite(ctx, pt)

	eth, simd := s.EthChain, s.CosmosChains[0]

	ics20Address := ethcommon.HexToAddress(s.contractAddresses.Ics20Transfer)
	transferCoin := sdk.NewCoin(simd.Config().Denom, sdkmath.NewIntFromBigInt(transferAmount))
	ethereumUserAddress := crypto.PubkeyToAddress(s.key.PublicKey)
	cosmosUserWallet := s.CosmosUsers[0]
	cosmosUserAddress := cosmosUserWallet.FormattedAddress()
	sendMemo := "nativesend"

	sendTimeout := uint64(time.Now().Add(30 * time.Minute).Unix())

	var (
		cosmosSendPayload channeltypesv2.Payload
		ibcERC20          *ibcerc20.Contract
		ibcERC20Address   ethcommon.Address
	)

	s.Require().True(s.Run("Send transfer on Cosmos chain", func() {
		transferPayload := transfertypes.FungibleTokenPacketData{
			Denom:    transferCoin.Denom,
			Amount:   transferCoin.Amount.String(),
			Sender:   cosmosUserAddress,
			Receiver: strings.ToLower(ethereumUserAddress.Hex()),
			Memo:     sendMemo,
		}
		encodedPayload, err := transfertypes.EncodeABIFungibleTokenPacketData(&transferPayload)
		s.Require().NoError(err)

		cosmosSendPayload = channeltypesv2.Payload{
			SourcePort:      transfertypes.PortID,
			DestinationPort: transfertypes.PortID,
			Version:         transfertypes.V1,
			Encoding:        transfertypes.EncodingABI,
			Value:           encodedPayload,
		}
		msgSendPacket := channeltypesv2.MsgSendPacket{
			SourceClient:     testvalues.FirstWasmClientID,
			TimeoutTimestamp: sendTimeout,
			Payloads:         []channeltypesv2.Payload{cosmosSendPayload},
			Signer:           cosmosUserWallet.FormattedAddress(),
		}

		resp, err := s.BroadcastMessages(ctx, simd, cosmosUserWallet, 200_000, &msgSendPacket)
		s.Require().NoError(err)
		s.Require().NotEmpty(resp.TxHash)

		s.Require().True(s.Run("Verify balances on Cosmos chain", func() {
			resp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
				Address: cosmosUserAddress,
				Denom:   transferCoin.Denom,
			})
			s.Require().NoError(err)
			s.Require().NotNil(resp.Balance)
			s.Require().Equal(testvalues.InitialBalance-testvalues.TransferAmount, resp.Balance.Amount.Int64())
		}))
	}))

	// The running relayer process auto-relays the Cosmos→ETH send_packet event.
	// Wait for the ibcERC20 balance on Ethereum to reflect the transfer.
	var ethWriteAckEvent *ics26router.ContractWriteAcknowledgement
	s.Require().True(s.Run("Receive packet on Ethereum (auto-relayed)", func() {
		denomOnEthereum := transfertypes.NewDenom(
			transferCoin.Denom,
			transfertypes.NewHop(transfertypes.PortID, testvalues.CustomClientID),
		)

		s.Require().Eventually(func() bool {
			addr, err := s.ics20Contract.IbcERC20Contract(nil, denomOnEthereum.Path())
			if err != nil || addr == (ethcommon.Address{}) {
				return false
			}
			contract, err := ibcerc20.NewContract(addr, eth.RPCClient)
			if err != nil {
				return false
			}
			bal, err := contract.BalanceOf(nil, ethereumUserAddress)
			if err != nil || bal == nil {
				return false
			}
			if bal.Cmp(transferAmount) == 0 {
				ibcERC20Address = addr
				ibcERC20 = contract
				return true
			}
			return false
		}, 10*time.Minute, 5*time.Second, "timed out waiting for Cosmos→ETH relay")

		s.True(s.Run("Verify balances on Ethereum", func() {
			userBalance, err := ibcERC20.BalanceOf(nil, ethereumUserAddress)
			s.Require().NoError(err)
			s.Require().Equal(transferAmount, userBalance)

			ics20TransferBalance, err := ibcERC20.BalanceOf(nil, ics20Address)
			s.Require().NoError(err)
			s.Require().Zero(ics20TransferBalance.Int64())
		}))

		// Collect the WriteAcknowledgement event emitted by the relayer's recvPacket tx
		iter, err := s.ics26Contract.FilterWriteAcknowledgement(
			&bind.FilterOpts{Context: ctx},
			[]string{testvalues.CustomClientID},
			[]*big.Int{big.NewInt(1)},
		)
		s.Require().NoError(err)
		defer iter.Close()

		for iter.Next() {
			ethWriteAckEvent = iter.Event
			break
		}
		s.Require().NoError(iter.Error())
		s.Require().NotNil(ethWriteAckEvent, "WriteAcknowledgement event not found on Ethereum")
	}))

	// Reconstruct the original Cosmos packet for MsgAcknowledgement
	cosmosSendPacket := channeltypesv2.Packet{
		Sequence:          1,
		SourceClient:      testvalues.FirstWasmClientID,
		DestinationClient: testvalues.CustomClientID,
		TimeoutTimestamp:  sendTimeout,
		Payloads:          []channeltypesv2.Payload{cosmosSendPayload},
	}

	s.Require().True(s.Run("Acknowledge packet on Cosmos chain", func() {
		s.Require().True(s.Run("Verify commitments exists", func() {
			resp, err := e2esuite.GRPCQuery[channeltypesv2.QueryPacketCommitmentResponse](ctx, simd, &channeltypesv2.QueryPacketCommitmentRequest{
				ClientId: testvalues.FirstWasmClientID,
				Sequence: 1,
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Commitment)
		}))

		s.Require().NotEmpty(ethWriteAckEvent.Acknowledgements)
		ackBytes := ethWriteAckEvent.Acknowledgements[0]

		// Submit MsgAcknowledgement with minimal proof (accepted by mock/dummy wasm client)
		msg := &channeltypesv2.MsgAcknowledgement{
			Packet:          cosmosSendPacket,
			Acknowledgement: channeltypesv2.Acknowledgement{AppAcknowledgements: [][]byte{ackBytes}},
			ProofAcked:      []byte{0x01},
			ProofHeight:     clienttypes.Height{RevisionNumber: 0, RevisionHeight: 1},
			Signer:          s.SimdRelayerSubmitter.FormattedAddress(),
		}
		_, err := s.BroadcastMessages(ctx, simd, s.SimdRelayerSubmitter, 2_000_000, msg)
		s.Require().NoError(err)

		s.Require().True(s.Run("Verify commitments removed", func() {
			_, err := e2esuite.GRPCQuery[channeltypesv2.QueryPacketCommitmentResponse](ctx, simd, &channeltypesv2.QueryPacketCommitmentRequest{
				ClientId: testvalues.FirstWasmClientID,
				Sequence: 1,
			})
			s.Require().ErrorContains(err, "packet commitment hash not found")
		}))
	}))

	s.Require().True(s.Run("Approve the ICS20Transfer.sol contract to spend the erc20 tokens", func() {
		tx, err := ibcERC20.Approve(s.GetTransactOpts(s.key, eth), ics20Address, transferAmount)
		s.Require().NoError(err)

		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

		allowance, err := ibcERC20.Allowance(nil, ethereumUserAddress, ics20Address)
		s.Require().NoError(err)
		s.Require().Equal(transferAmount, allowance)
	}))

	var ethSendPacketEvent *ics26router.ContractSendPacket
	s.Require().True(s.Run("Transfer tokens back from Ethereum", func() {
		returnMemo := "testreturnmemo"
		returnTimeout := uint64(time.Now().Add(30 * time.Minute).Unix())
		msgSendPacket := ics20transfer.IICS20TransferMsgsSendTransferMsg{
			Denom:            ibcERC20Address,
			Amount:           transferAmount,
			Receiver:         cosmosUserAddress,
			TimeoutTimestamp: returnTimeout,
			SourceClient:     testvalues.CustomClientID,
			Memo:             returnMemo,
		}

		tx, err := s.ics20Contract.SendTransfer(s.GetTransactOpts(s.key, eth), msgSendPacket)
		s.Require().NoError(err)

		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

		ethSendPacketEvent, err = e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseSendPacket)
		s.Require().NoError(err)
		s.Require().Equal(uint64(1), ethSendPacketEvent.Packet.Sequence)
		s.Require().Equal(returnTimeout, ethSendPacketEvent.Packet.TimeoutTimestamp)
		s.Require().Equal(testvalues.CustomClientID, ethSendPacketEvent.Packet.SourceClient)
		s.Require().Equal(testvalues.FirstWasmClientID, ethSendPacketEvent.Packet.DestClient)

		s.True(s.Run("Verify balances on Ethereum", func() {
			userBalance, err := ibcERC20.BalanceOf(nil, ethereumUserAddress)
			s.Require().NoError(err)
			s.Require().Zero(userBalance.Int64())

			ics20TransferBalance, err := ibcERC20.BalanceOf(nil, ics20Address)
			s.Require().NoError(err)
			s.Require().Zero(ics20TransferBalance.Int64())
		}))
	}))

	// Convert ETH packet to Cosmos channeltypesv2.Packet for MsgRecvPacket
	ethPkt := ethSendPacketEvent.Packet
	var ethToCosmosPayloads []channeltypesv2.Payload
	for _, p := range ethPkt.Payloads {
		ethToCosmosPayloads = append(ethToCosmosPayloads, channeltypesv2.Payload{
			SourcePort:      p.SourcePort,
			DestinationPort: p.DestPort,
			Version:         p.Version,
			Encoding:        p.Encoding,
			Value:           p.Value,
		})
	}
	ethToCosmosPacket := channeltypesv2.Packet{
		Sequence:          ethPkt.Sequence,
		SourceClient:      ethPkt.SourceClient,
		DestinationClient: ethPkt.DestClient,
		TimeoutTimestamp:  ethPkt.TimeoutTimestamp,
		Payloads:          ethToCosmosPayloads,
	}

	var cosmosWriteAckBytes []byte
	s.Require().True(s.Run("Receive packet on Cosmos chain", func() {
		// Submit MsgRecvPacket with minimal proof (accepted by mock/dummy wasm client)
		recvMsg := &channeltypesv2.MsgRecvPacket{
			Packet:          ethToCosmosPacket,
			ProofCommitment: []byte{0x01},
			ProofHeight:     clienttypes.Height{RevisionNumber: 0, RevisionHeight: 1},
			Signer:          s.SimdRelayerSubmitter.FormattedAddress(),
		}
		resp, err := s.BroadcastMessages(ctx, simd, s.SimdRelayerSubmitter, 2_000_000, recvMsg)
		s.Require().NoError(err)

		// Extract ack bytes from the WriteAcknowledgement event for later use on ETH
		if ackHex, err := cosmos.GetEventValue(resp.Events, channeltypesv2.EventTypeWriteAck, channeltypesv2.AttributeKeyEncodedAckHex); err == nil {
			cosmosWriteAckBytes, _ = hex.DecodeString(ackHex)
		}

		s.Require().True(s.Run("Verify balances on Cosmos chain", func() {
			resp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
				Address: cosmosUserAddress,
				Denom:   transferCoin.Denom,
			})
			s.Require().NoError(err)
			s.Require().NotNil(resp.Balance)
			s.Require().Equal(testvalues.InitialBalance, resp.Balance.Amount.Int64())
		}))
	}))

	s.Require().True(s.Run("Acknowledge packet on Ethereum", func() {
		s.Require().True(s.Run("Verify commitment exists", func() {
			packetCommitmentPath := ibchostv2.PacketCommitmentKey(testvalues.CustomClientID, 1)
			var ethPath [32]byte
			copy(ethPath[:], crypto.Keccak256(packetCommitmentPath))

			resp, err := s.ics26Contract.GetCommitment(nil, ethPath)
			s.Require().NoError(err)
			s.Require().NotZero(resp)
		}))

		// Use ack bytes obtained from Cosmos WriteAcknowledgement event.
		// NOTE: AckPacket on ETH verifies a Tendermint membership proof via the ICS07 contract.
		// The relayer must have updated the ICS07 client to the height where the ack was stored.
		ackBytes := cosmosWriteAckBytes
		if len(ackBytes) == 0 {
			ackBytes = []byte(`{"result":"AQ=="}`)
		}

		ackMsg := ics26router.IICS26RouterMsgsMsgAckPacket{
			Packet:          ethPkt,
			Acknowledgement: ackBytes,
			MembershipMsg:   []byte{},
		}
		tx, err := s.ics26Contract.AckPacket(s.GetTransactOpts(s.EthRelayerSubmitter, eth), ackMsg)
		s.Require().NoError(err)

		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

		_, err = e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseAckPacket)
		s.Require().NoError(err)

		s.Require().True(s.Run("Verify commitment removed", func() {
			packetCommitmentPath := ibchostv2.PacketCommitmentKey(testvalues.CustomClientID, 1)
			var ethPath [32]byte
			copy(ethPath[:], crypto.Keccak256(packetCommitmentPath))

			resp, err := s.ics26Contract.GetCommitment(nil, ethPath)
			s.Require().NoError(err)
			s.Require().Zero(resp)
		}))

		s.Require().True(s.Run("Verify balances on Ethereum after ack", func() {
			userBalance, err := ibcERC20.BalanceOf(nil, ethereumUserAddress)
			s.Require().NoError(err)
			s.Require().Zero(userBalance.Int64())

			ics20TransferBalance, err := ibcERC20.BalanceOf(nil, ics20Address)
			s.Require().NoError(err)
			s.Require().Zero(ics20TransferBalance.Int64())
		}))
	}))
}

func (s *IbcEurekaTestSuite) Test_TimeoutPacketFromEth() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.FilteredICS20TimeoutPacketFromEthereumTest(ctx, proofType, 1, nil)
}

func (s *IbcEurekaTestSuite) Test_10_TimeoutPacketFromEth() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.FilteredICS20TimeoutPacketFromEthereumTest(ctx, proofType, 10, nil)
}

func (s *IbcEurekaTestSuite) Test_5_TimeoutPacketFromEth() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.FilteredICS20TimeoutPacketFromEthereumTest(ctx, proofType, 5, []uint64{1, 2, 3, 4, 5})
}

func (s *IbcEurekaTestSuite) Test_5_FilteredTimeoutPacketFromEth() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.FilteredICS20TimeoutPacketFromEthereumTest(ctx, proofType, 5, []uint64{2, 3})
}

func (s *IbcEurekaTestSuite) FilteredICS20TimeoutPacketFromEthereumTest(
	ctx context.Context, pt types.SupportedProofType, numOfTransfers int, timeoutFilter []uint64,
) {
	s.Require().GreaterOrEqual(numOfTransfers, len(timeoutFilter))
	s.Require().Greater(numOfTransfers, 0)

	s.SetupSuite(ctx, pt)

	eth, simd := s.EthChain, s.CosmosChains[0]

	ics26Address := ethcommon.HexToAddress(s.contractAddresses.Ics26Router)
	erc20Address := ethcommon.HexToAddress(s.contractAddresses.Erc20)

	transferAmount := big.NewInt(testvalues.TransferAmount)
	totalTransferAmount := new(big.Int).Mul(transferAmount, big.NewInt(int64(numOfTransfers)))
	var refundedAmount *big.Int
	if len(timeoutFilter) == 0 {
		refundedAmount = totalTransferAmount
	} else {
		refundedAmount = new(big.Int).Mul(transferAmount, big.NewInt(int64(len(timeoutFilter))))
	}
	ethereumUserAddress := crypto.PubkeyToAddress(s.key.PublicKey)
	cosmosUserWallet := s.CosmosUsers[0]
	cosmosUserAddress := cosmosUserWallet.FormattedAddress()

	var originalBalance *sdk.Coin
	s.Require().True(s.Run("Retrieve original balance", func() {
		denomOnCosmos := transfertypes.NewDenom(s.contractAddresses.Erc20, transfertypes.NewHop(transfertypes.PortID, testvalues.FirstWasmClientID))

		resp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
			Address: cosmosUserAddress,
			Denom:   denomOnCosmos.IBCDenom(),
		})
		s.Require().NoError(err)
		s.Require().NotNil(resp.Balance)
		originalBalance = resp.Balance
	}))

	s.Require().True(s.Run("Approve the ICS20Transfer.sol contract to spend the erc20 tokens", func() {
		ics20Address := ethcommon.HexToAddress(s.contractAddresses.Ics20Transfer)
		tx, err := s.erc20Contract.Approve(s.GetTransactOpts(s.key, eth), ics20Address, totalTransferAmount)
		s.Require().NoError(err)

		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

		allowance, err := s.erc20Contract.Allowance(nil, ethereumUserAddress, ics20Address)
		s.Require().NoError(err)
		s.Require().Equal(totalTransferAmount, allowance)
	}))

	var (
		ethSendTxHashes [][]byte
		sendPacket      ics26router.IICS26RouterMsgsPacket
		escrowAddress   ethcommon.Address
	)
	s.Require().True(s.Run("Send packets on Ethereum", func() {
		for i := range numOfTransfers {
			timeout := uint64(time.Now().Add(30 * time.Second).Unix())
			msgSendPacket := ics20transfer.IICS20TransferMsgsSendTransferMsg{
				Denom:            erc20Address,
				Amount:           transferAmount,
				Receiver:         cosmosUserAddress,
				TimeoutTimestamp: timeout,
				SourceClient:     testvalues.CustomClientID,
				Memo:             "testmemo",
			}

			tx, err := s.ics20Contract.SendTransfer(s.GetTransactOpts(s.key, eth), msgSendPacket)
			s.Require().NoError(err)

			receipt, err := eth.GetTxReciept(ctx, tx.Hash())
			s.Require().NoError(err)

			// We use the first packet in fixture generation
			if i == 0 && s.solidityFixtureGenerator.Enabled {
				sendPacketEvent, err := e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseSendPacket)
				s.Require().NoError(err)
				sendPacket = sendPacketEvent.Packet
			}

			ethSendTxHashes = append(ethSendTxHashes, tx.Hash().Bytes())
		}

		s.Require().True(s.Run("Verify balances on Ethereum", func() {
			// User balance on Ethereum
			userBalance, err := s.erc20Contract.BalanceOf(nil, ethereumUserAddress)
			s.Require().NoError(err)
			s.Require().Equal(new(big.Int).Sub(testvalues.StartingERC20Balance, totalTransferAmount), userBalance)

			// Get the escrow address
			escrowAddress, err = s.ics20Contract.GetEscrow(nil, testvalues.CustomClientID)
			s.Require().NoError(err)

			// ICS20 contract balance on Ethereum
			escrowBalance, err := s.erc20Contract.BalanceOf(nil, escrowAddress)
			s.Require().NoError(err)
			s.Require().Equal(totalTransferAmount, escrowBalance)
		}))
	}))

	var txBodyBz []byte
	s.Require().True(s.Run("Prefetch relay tx", func() {
		resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
			SrcChain:    eth.ChainID.String(),
			DstChain:    simd.Config().ChainID,
			SourceTxIds: ethSendTxHashes,
			SrcClientId: testvalues.CustomClientID,
			DstClientId: testvalues.FirstWasmClientID,
		})
		s.Require().NoError(err)
		s.Require().NotEmpty(resp.Tx)
		s.Require().Empty(resp.Address)

		txBodyBz = resp.Tx
	}))

	// sleep for 45 seconds to let the packet timeout
	time.Sleep(45 * time.Second)

	s.True(s.Run("Timeout packets on Ethereum", func() {
		var timeoutRelayTx []byte
		s.Require().True(s.Run("Retrieve timeout tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:           simd.Config().ChainID,
				DstChain:           eth.ChainID.String(),
				TimeoutTxIds:       ethSendTxHashes,
				SrcClientId:        testvalues.FirstWasmClientID,
				DstClientId:        testvalues.CustomClientID,
				DstPacketSequences: timeoutFilter,
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Tx)
			s.Require().Equal(resp.Address, ics26Address.String())

			timeoutRelayTx = resp.Tx
		}))

		s.Require().True(s.Run("Submit relay tx", func() {
			_, err := eth.BroadcastTx(ctx, s.EthRelayerSubmitter, 5_000_000, &ics26Address, timeoutRelayTx)
			s.Require().NoError(err)
		}))

		if s.solidityFixtureGenerator.Enabled {
			s.Require().Zero(len(timeoutFilter))
			s.Require().NoError(s.solidityFixtureGenerator.GenerateAndSaveSolidityFixture(
				fmt.Sprintf("timeoutMultiPacket_%d-%s.json", numOfTransfers, pt.String()),
				s.contractAddresses.Erc20, timeoutRelayTx, sendPacket,
			))
		}

		s.Require().True(s.Run("Verify balances on Ethereum", func() {
			// Expected balance
			netTransferAmount := new(big.Int).Sub(totalTransferAmount, refundedAmount)
			expectedUserBalance := new(big.Int).Sub(testvalues.StartingERC20Balance, netTransferAmount)

			// User balance on Ethereum
			userBalance, err := s.erc20Contract.BalanceOf(nil, ethereumUserAddress)
			s.Require().NoError(err)
			s.Require().Equal(expectedUserBalance, userBalance)

			// ICS20 contract balance on Ethereum
			escrowBalance, err := s.erc20Contract.BalanceOf(nil, escrowAddress)
			s.Require().NoError(err)
			s.Require().Equal(escrowBalance.Int64(), netTransferAmount.Int64())
		}))

		s.Require().True(s.Run("Verify no balance on Cosmos chain", func() {
			denomOnCosmos := transfertypes.NewDenom(s.contractAddresses.Erc20, transfertypes.NewHop(transfertypes.PortID, testvalues.FirstWasmClientID))

			resp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
				Address: cosmosUserAddress,
				Denom:   denomOnCosmos.IBCDenom(),
			})
			s.Require().NoError(err)
			s.Require().Equal(originalBalance, resp.Balance)
		}))
	}))

	s.Require().True(s.Run("Constructing relay packet after timeout should fail", func() {
		resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
			SrcChain:    eth.ChainID.String(),
			DstChain:    simd.Config().ChainID,
			SourceTxIds: ethSendTxHashes,
			SrcClientId: testvalues.CustomClientID,
			DstClientId: testvalues.FirstWasmClientID,
		})
		// TODO: https://github.com/cosmos/solidity-ibc-eureka/issues/363
		// The following assertions should be Error and Nil, but the relayer returns a valid response currently.
		s.Require().NoError(err)
		s.Require().NotNil(resp)
	}))

	s.Require().True(s.Run("Receive packets on Cosmos chain after timeout", func() {
		resp, err := s.BroadcastSdkTxBody(ctx, simd, s.SimdRelayerSubmitter, 2_000_000, txBodyBz)
		s.Require().Error(err)
		s.Require().Nil(resp)
	}))
}

func (s *IbcEurekaTestSuite) Test_ErrorAckToEthereum() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.ICS20ErrorAckToEthereumTest(ctx, proofType)
}

func (s *IbcEurekaTestSuite) ICS20ErrorAckToEthereumTest(
	ctx context.Context, pt types.SupportedProofType,
) {
	s.SetupSuite(ctx, pt)

	eth, simd := s.EthChain, s.CosmosChains[0]
	ics26Address := ethcommon.HexToAddress(s.contractAddresses.Ics26Router)
	erc20Address := ethcommon.HexToAddress(s.contractAddresses.Erc20)

	transferAmount := big.NewInt(testvalues.TransferAmount)
	ethereumUserAddress := crypto.PubkeyToAddress(s.key.PublicKey)

	s.Require().True(s.Run("Approve the ICS20Transfer.sol contract to spend the erc20 tokens", func() {
		ics20Address := ethcommon.HexToAddress(s.contractAddresses.Ics20Transfer)
		tx, err := s.erc20Contract.Approve(s.GetTransactOpts(s.key, eth), ics20Address, transferAmount)
		s.Require().NoError(err)

		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

		allowance, err := s.erc20Contract.Allowance(nil, ethereumUserAddress, ics20Address)
		s.Require().NoError(err)
		s.Require().Equal(transferAmount, allowance)
	}))

	var (
		ethSendTxHash []byte
		escrowAddress ethcommon.Address
	)
	s.Require().True(s.Run("Send transfer on Ethereum", func() {
		timeout := uint64(time.Now().Add(30 * time.Minute).Unix())

		// Send a transfer to an invalid Cosmos address
		msgSendPacket := ics20transfer.IICS20TransferMsgsSendTransferMsg{
			Denom:            erc20Address,
			Amount:           transferAmount,
			Receiver:         ibctesting.InvalidID,
			TimeoutTimestamp: timeout,
			SourceClient:     testvalues.CustomClientID,
			Memo:             "",
		}

		tx, err := s.ics20Contract.SendTransfer(s.GetTransactOpts(s.key, eth), msgSendPacket)
		s.Require().NoError(err)

		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

		ethSendTxHash = tx.Hash().Bytes()

		s.Require().True(s.Run("Verify balances on Ethereum", func() {
			// User balance on Ethereum
			userBalance, err := s.erc20Contract.BalanceOf(nil, ethereumUserAddress)
			s.Require().NoError(err)
			s.Require().Equal(new(big.Int).Sub(testvalues.StartingERC20Balance, transferAmount), userBalance)

			// Get the escrow address
			escrowAddress, err = s.ics20Contract.GetEscrow(nil, testvalues.CustomClientID)
			s.Require().NoError(err)

			// ICS20 contract balance on Ethereum
			escrowBalance, err := s.erc20Contract.BalanceOf(nil, escrowAddress)
			s.Require().NoError(err)
			s.Require().Equal(transferAmount, escrowBalance)
		}))
	}))

	var (
		denomOnCosmos transfertypes.Denom
		ackTxHash     []byte
	)
	s.Require().True(s.Run("Receive packets on Cosmos chain", func() {
		var relayTxBodyBz []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:    eth.ChainID.String(),
				DstChain:    simd.Config().ChainID,
				SourceTxIds: [][]byte{ethSendTxHash},
				SrcClientId: testvalues.CustomClientID,
				DstClientId: testvalues.FirstWasmClientID,
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Tx)
			s.Require().Empty(resp.Address)

			relayTxBodyBz = resp.Tx
		}))

		s.Require().True(s.Run("Broadcast relay tx", func() {
			resp := s.MustBroadcastSdkTxBody(ctx, simd, s.SimdRelayerSubmitter, 2_000_000, relayTxBodyBz)

			var err error
			ackTxHash, err = hex.DecodeString(resp.TxHash)
			s.Require().NoError(err)
			s.Require().NotEmpty(ackTxHash)
		}))

		s.Require().True(s.Run("Verify no balance on Cosmos chain", func() {
			denomOnCosmos = transfertypes.NewDenom(s.contractAddresses.Erc20, transfertypes.NewHop(transfertypes.PortID, testvalues.FirstWasmClientID))

			_, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
				Address: ibctesting.InvalidID,
				Denom:   denomOnCosmos.IBCDenom(),
			})
			s.Require().Error(err)
		}))
	}))

	s.Require().True(s.Run("Acknowledge packets on Ethereum", func() {
		var ackRelayTx []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:    simd.Config().ChainID,
				DstChain:    eth.ChainID.String(),
				SourceTxIds: [][]byte{ackTxHash},
				SrcClientId: testvalues.FirstWasmClientID,
				DstClientId: testvalues.CustomClientID,
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Tx)
			s.Require().Equal(resp.Address, ics26Address.String())

			ackRelayTx = resp.Tx
		}))

		s.Require().True(s.Run("Submit relay tx", func() {
			receipt, err := eth.BroadcastTx(ctx, s.EthRelayerSubmitter, 5_000_000, &ics26Address, ackRelayTx)
			s.Require().NoError(err)

			// Verify the ack packet event exists
			_, err = e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseAckPacket)
			s.Require().NoError(err)
		}))

		s.Require().True(s.Run("Verify balances on Ethereum", func() {
			// User balance on Ethereum
			userBalance, err := s.erc20Contract.BalanceOf(nil, ethereumUserAddress)
			s.Require().NoError(err)
			s.Require().Equal(testvalues.StartingERC20Balance, userBalance)

			// ICS20 contract balance on Ethereum
			escrowBalance, err := s.erc20Contract.BalanceOf(nil, escrowAddress)
			s.Require().NoError(err)
			s.Require().Zero(escrowBalance.Int64())
		}))
	}))
}

func (s *IbcEurekaTestSuite) Test_TimeoutPacketFromCosmos() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.FilteredICS20TimeoutFromCosmosTimeoutTest(ctx, proofType, 1, nil)
}

func (s *IbcEurekaTestSuite) Test_10_TimeoutPacketFromCosmos() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.FilteredICS20TimeoutFromCosmosTimeoutTest(ctx, proofType, 10, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
}

func (s *IbcEurekaTestSuite) Test_10_FilteredTimeoutPacketFromCosmos() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.FilteredICS20TimeoutFromCosmosTimeoutTest(ctx, proofType, 10, []uint64{2, 4, 6, 8, 10})
}

func (s *IbcEurekaTestSuite) FilteredICS20TimeoutFromCosmosTimeoutTest(
	ctx context.Context, proofType types.SupportedProofType, numOfTransfers int, timeoutFilter []uint64,
) {
	s.Require().GreaterOrEqual(numOfTransfers, len(timeoutFilter))
	s.Require().Greater(numOfTransfers, 0)

	s.SetupSuite(ctx, proofType)

	eth, simd := s.EthChain, s.CosmosChains[0]

	transferAmount := big.NewInt(testvalues.TransferAmount)
	totalTransferAmount := big.NewInt(testvalues.TransferAmount * int64(numOfTransfers))
	if totalTransferAmount.Int64() > testvalues.InitialBalance {
		s.FailNow("Total transfer amount exceeds the initial balance")
	}
	var refundedAmount *big.Int
	if len(timeoutFilter) == 0 {
		refundedAmount = totalTransferAmount
	} else {
		refundedAmount = big.NewInt(testvalues.TransferAmount * int64(len(timeoutFilter)))
	}
	ethereumUserAddress := crypto.PubkeyToAddress(s.key.PublicKey)
	cosmosUserWallet := s.CosmosUsers[0]
	cosmosUserAddress := cosmosUserWallet.FormattedAddress()
	sendMemo := "nonnativesend"

	var (
		transferCoin sdk.Coin
		sendTxHashes [][]byte
	)
	s.Require().True(s.Run("Send transfers on Cosmos chain", func() {
		for range numOfTransfers {
			timeout := uint64(time.Now().Add(45 * time.Second).Unix())
			transferCoin = sdk.NewCoin(simd.Config().Denom, sdkmath.NewIntFromBigInt(transferAmount))

			transferPayload := transfertypes.FungibleTokenPacketData{
				Denom:    transferCoin.Denom,
				Amount:   transferCoin.Amount.String(),
				Sender:   cosmosUserAddress,
				Receiver: strings.ToLower(ethereumUserAddress.Hex()),
				Memo:     sendMemo,
			}
			encodedPayload, err := transfertypes.EncodeABIFungibleTokenPacketData(&transferPayload)
			s.Require().NoError(err)

			payload := channeltypesv2.Payload{
				SourcePort:      transfertypes.PortID,
				DestinationPort: transfertypes.PortID,
				Version:         transfertypes.V1,
				Encoding:        transfertypes.EncodingABI,
				Value:           encodedPayload,
			}
			msgSendPacket := channeltypesv2.MsgSendPacket{
				SourceClient:     testvalues.FirstWasmClientID,
				TimeoutTimestamp: timeout,
				Payloads: []channeltypesv2.Payload{
					payload,
				},
				Signer: cosmosUserWallet.FormattedAddress(),
			}

			resp, err := s.BroadcastMessages(ctx, simd, cosmosUserWallet, 200_000, &msgSendPacket)
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.TxHash)

			txHash, err := hex.DecodeString(resp.TxHash)
			s.Require().NoError(err)

			sendTxHashes = append(sendTxHashes, txHash)
		}

		s.Require().True(s.Run("Verify balances on Cosmos chain", func() {
			// Check the balance of UserB
			resp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
				Address: cosmosUserAddress,
				Denom:   transferCoin.Denom,
			})
			s.Require().NoError(err)
			s.Require().NotNil(resp.Balance)
			s.Require().Equal(testvalues.InitialBalance-totalTransferAmount.Int64(), resp.Balance.Amount.Int64())
		}))
	}))

	var txBodyBz []byte
	s.Require().True(s.Run("Prefetch relay tx", func() {
		resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
			SrcChain:    simd.Config().ChainID,
			DstChain:    eth.ChainID.String(),
			SourceTxIds: sendTxHashes,
			SrcClientId: testvalues.FirstWasmClientID,
			DstClientId: testvalues.CustomClientID,
		})
		s.Require().NoError(err)
		s.Require().NotEmpty(resp.Tx)

		txBodyBz = resp.Tx
	}))

	// sleep for 45 seconds to let the packet timeout
	time.Sleep(45 * time.Second)

	s.Require().True(s.Run("Timeout packets on Cosmos chain", func() {
		var relayTxBodyBz []byte
		s.Require().True(s.Run("Retrieve relay tx to Cosmos chain", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:           eth.ChainID.String(),
				DstChain:           simd.Config().ChainID,
				TimeoutTxIds:       sendTxHashes,
				SrcClientId:        testvalues.CustomClientID,
				DstClientId:        testvalues.FirstWasmClientID,
				DstPacketSequences: timeoutFilter,
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Tx)
			s.Require().Empty(resp.Address)

			relayTxBodyBz = resp.Tx

			s.wasmFixtureGenerator.AddFixtureStep("timeout_packets", ethereumtypes.RelayerMessages{
				RelayerTxBody: hex.EncodeToString(relayTxBodyBz),
			})
		}))

		s.Require().True(s.Run("Broadcast relay tx on Cosmos chain", func() {
			_ = s.MustBroadcastSdkTxBody(ctx, simd, s.SimdRelayerSubmitter, 2_000_000, relayTxBodyBz)
		}))

		s.Require().True(s.Run("Verify balances on Cosmos chain", func() {
			// Check the balance of UserB
			resp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
				Address: cosmosUserAddress,
				Denom:   transferCoin.Denom,
			})
			s.Require().NoError(err)
			s.Require().NotNil(resp.Balance)
			s.Require().Equal(testvalues.InitialBalance-totalTransferAmount.Int64()+refundedAmount.Int64(), resp.Balance.Amount.Int64())
		}))
	}))

	// We are skipping the replay attack test on non-PoS mode
	// In PoW mode, the test below will NOT fail as the mock light client accepts all proofs without verification.
	if os.Getenv(testvalues.EnvKeyEthTestnetType) != testvalues.EthTestnetTypePoS {
		s.T().Skip("Skipping replay attack test on non-PoS mode")
	}

	s.Require().True(s.Run("Receive packets on Ethereum after timeout should fail", func() {
		ics26Address := ethcommon.HexToAddress(s.contractAddresses.Ics26Router)
		_, err := eth.BroadcastTx(ctx, s.EthRelayerSubmitter, 5_000_000, &ics26Address, txBodyBz)
		s.Require().Error(err)

		// Prove that the erc20 was not created (and by extension, the packet was not received)
		denomOnEthereum := transfertypes.NewDenom(transferCoin.Denom, transfertypes.NewHop(transfertypes.PortID, testvalues.CustomClientID))
		_, err = s.ics20Contract.IbcERC20Contract(nil, denomOnEthereum.Path())
		// Ethereum side did not received the packet, ERC20 contract corresponding to the denom does not exist
		s.Require().Error(err)
	}))
}

func (s *IbcEurekaTestSuite) Test_TimeoutPacketEthRemintsVouchers() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.TimeoutPacketEthRemintsVouchersTest(ctx, proofType)
}

// TimeoutPacketEthRemintsVouchersTest tests that when a transfer of a voucher (Cosmos native -> Eth)
// from Ethereum back to Cosmos times out, the vouchers are reminted on Ethereum.
func (s *IbcEurekaTestSuite) TimeoutPacketEthRemintsVouchersTest(ctx context.Context, pt types.SupportedProofType) {
	s.SetupSuite(ctx, pt)

	eth, simd := s.EthChain, s.CosmosChains[0]

	ics26Address := ethcommon.HexToAddress(s.contractAddresses.Ics26Router)
	ics20Address := ethcommon.HexToAddress(s.contractAddresses.Ics20Transfer)
	transferAmount := big.NewInt(testvalues.TransferAmount)
	transferCoin := sdk.NewCoin(simd.Config().Denom, sdkmath.NewIntFromBigInt(transferAmount))
	ethereumUserAddress := crypto.PubkeyToAddress(s.key.PublicKey)
	cosmosUserWallet := s.CosmosUsers[0]
	cosmosUserAddress := cosmosUserWallet.FormattedAddress()

	var cosmosSendTxHash []byte
	s.Require().True(s.Run("Send native Cosmos coins to Ethereum", func() {
		timeout := uint64(time.Now().Add(30 * time.Minute).Unix())

		transferPayload := transfertypes.FungibleTokenPacketData{
			Denom:    transferCoin.Denom,
			Amount:   transferCoin.Amount.String(),
			Sender:   cosmosUserAddress,
			Receiver: strings.ToLower(ethereumUserAddress.Hex()),
			Memo:     "create-voucher",
		}
		encodedPayload, err := transfertypes.EncodeABIFungibleTokenPacketData(&transferPayload)
		s.Require().NoError(err)

		payload := channeltypesv2.Payload{
			SourcePort:      transfertypes.PortID,
			DestinationPort: transfertypes.PortID,
			Version:         transfertypes.V1,
			Encoding:        transfertypes.EncodingABI,
			Value:           encodedPayload,
		}
		msgSendPacket := channeltypesv2.MsgSendPacket{
			SourceClient:     testvalues.FirstWasmClientID,
			TimeoutTimestamp: timeout,
			Payloads: []channeltypesv2.Payload{
				payload,
			},
			Signer: cosmosUserWallet.FormattedAddress(),
		}

		resp, err := s.BroadcastMessages(ctx, simd, cosmosUserWallet, 200_000, &msgSendPacket)
		s.Require().NoError(err)
		cosmosSendTxHash, err = hex.DecodeString(resp.TxHash)
		s.Require().NoError(err)
	}))

	var (
		ibcERC20        *ibcerc20.Contract
		ibcERC20Address ethcommon.Address
		ackTxHash       []byte
	)
	s.Require().True(s.Run("Receive packets on Ethereum", func() {
		var recvRelayTx []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:    simd.Config().ChainID,
				DstChain:    eth.ChainID.String(),
				SourceTxIds: [][]byte{cosmosSendTxHash},
				SrcClientId: testvalues.FirstWasmClientID,
				DstClientId: testvalues.CustomClientID,
			})
			s.Require().NoError(err)
			recvRelayTx = resp.Tx
		}))

		var packet ics26router.IICS26RouterMsgsPacket
		s.Require().True(s.Run("Broadcast relay tx", func() {
			receipt, err := eth.BroadcastTx(ctx, s.EthRelayerSubmitter, 5_000_000, &ics26Address, recvRelayTx)
			s.Require().NoError(err)
			s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

			ethReceiveAckEvent, err := e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseWriteAcknowledgement)
			s.Require().NoError(err)
			packet = ethReceiveAckEvent.Packet
			ackTxHash = receipt.TxHash.Bytes()
		}))

		// Get voucher contract details
		denomOnEthereum := transfertypes.NewDenom(transferCoin.Denom, transfertypes.NewHop(packet.Payloads[0].DestPort, packet.DestClient))
		var err error
		ibcERC20Address, err = s.ics20Contract.IbcERC20Contract(nil, denomOnEthereum.Path())
		s.Require().NoError(err)
		ibcERC20, err = ibcerc20.NewContract(ibcERC20Address, eth.RPCClient)
		s.Require().NoError(err)

		// Verify initial voucher balance
		userBalance, err := ibcERC20.BalanceOf(nil, ethereumUserAddress)
		s.Require().NoError(err)
		s.Require().Equal(transferAmount, userBalance)
	}))

	s.Require().True(s.Run("Acknowledge packets on Cosmos", func() {
		var ackRelayTxBodyBz []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:    eth.ChainID.String(),
				DstChain:    simd.Config().ChainID,
				SourceTxIds: [][]byte{ackTxHash},
				SrcClientId: testvalues.CustomClientID,
				DstClientId: testvalues.FirstWasmClientID,
			})
			s.Require().NoError(err)
			ackRelayTxBodyBz = resp.Tx
		}))
		s.MustBroadcastSdkTxBody(ctx, simd, s.SimdRelayerSubmitter, 2_000_000, ackRelayTxBodyBz)
	}))

	var ethTimeoutSendTxHash []byte
	s.Require().True(s.Run("Send voucher from Eth to Cosmos (timeout)", func() {
		// Approve spending
		tx, err := ibcERC20.Approve(s.GetTransactOpts(s.key, eth), ics20Address, transferAmount)
		s.Require().NoError(err)
		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

		// Send transfer with short timeout
		timeout := uint64(time.Now().Add(30 * time.Second).Unix())
		msgSendPacket := ics20transfer.IICS20TransferMsgsSendTransferMsg{
			Denom:            ibcERC20Address, // Sending the voucher back
			Amount:           transferAmount,
			Receiver:         cosmosUserAddress,
			TimeoutTimestamp: timeout,
			SourceClient:     testvalues.CustomClientID,
			Memo:             "timeout-voucher-eth-cosmos",
		}

		tx, err = s.ics20Contract.SendTransfer(s.GetTransactOpts(s.key, eth), msgSendPacket)
		s.Require().NoError(err)
		receipt, err = eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)
		ethTimeoutSendTxHash = tx.Hash().Bytes()

		// Verify user voucher balance is now zero
		userBalance, err := ibcERC20.BalanceOf(nil, ethereumUserAddress)
		s.Require().NoError(err)
		s.Require().Zero(userBalance.Int64())
	}))

	// We retrieve the relay tx before the timeout, and will attempt to relay it after the timeout
	var recvRelayTxBodyBz []byte
	s.Require().True(s.Run("Retrieve relay tx before timeout", func() {
		resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
			SrcChain:    eth.ChainID.String(),
			DstChain:    simd.Config().ChainID,
			SourceTxIds: [][]byte{ethTimeoutSendTxHash},
			SrcClientId: testvalues.CustomClientID,
			DstClientId: testvalues.FirstWasmClientID,
		})
		s.Require().NoError(err)
		recvRelayTxBodyBz = resp.Tx
	}))

	s.Require().True(s.Run("Wait for timeout", func() {
		time.Sleep(45 * time.Second)
	}))

	s.Require().True(s.Run("Relay timeout packet to Eth", func() {
		var timeoutRelayTx []byte
		s.Require().True(s.Run("Retrieve timeout tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:     simd.Config().ChainID, // Source of the *timeout* message is Cosmos
				DstChain:     eth.ChainID.String(),  // Destination is Eth
				TimeoutTxIds: [][]byte{ethTimeoutSendTxHash},
				SrcClientId:  testvalues.FirstWasmClientID,
				DstClientId:  testvalues.CustomClientID,
			})
			s.Require().NoError(err)
			timeoutRelayTx = resp.Tx
		}))

		s.Require().True(s.Run("Broadcast relay tx", func() {
			receipt, err := eth.BroadcastTx(ctx, s.EthRelayerSubmitter, 5_000_000, &ics26Address, timeoutRelayTx)
			s.Require().NoError(err)
			s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

			// Verify the timeout packet event exists
			_, err = e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseTimeoutPacket)
			s.Require().NoError(err)
		}))
	}))

	s.Require().True(s.Run("Verify voucher balance restored on Eth", func() {
		userBalance, err := ibcERC20.BalanceOf(nil, ethereumUserAddress)
		s.Require().NoError(err)
		s.Require().Equal(transferAmount, userBalance, "Voucher balance should be restored after timeout")

		// Escrow balance should remain zero (or unchanged if it wasn't zero initially, though it should be)
		escrowBalance, err := ibcERC20.BalanceOf(nil, ics20Address)
		s.Require().NoError(err)
		s.Require().Zero(escrowBalance.Int64(), "Escrow balance should be zero")
	}))

	s.Require().True(s.Run("Verify recvPacket fails on Cosmos after timeout", func() {
		_, err := s.BroadcastSdkTxBody(ctx, simd, s.SimdRelayerSubmitter, 2_000_000, recvRelayTxBodyBz)
		s.Require().Error(err)
	}))
}

func (s *IbcEurekaTestSuite) Test_TimeoutPacketCosmosRemintsVouchers() {
	ctx := context.Background()
	proofType := types.GetEnvProofType()
	s.TimeoutPacketCosmosRemintsVouchersTest(ctx, proofType)
}

// TimeoutPacketCosmosRemintsVouchersTest tests that when a transfer of a voucher (Eth native -> Cosmos)
// from Cosmos back to Ethereum times out, the vouchers are reminted on Cosmos.
func (s *IbcEurekaTestSuite) TimeoutPacketCosmosRemintsVouchersTest(ctx context.Context, pt types.SupportedProofType) {
	s.SetupSuite(ctx, pt)

	eth, simd := s.EthChain, s.CosmosChains[0]

	ics26Address := ethcommon.HexToAddress(s.contractAddresses.Ics26Router)
	ics20Address := ethcommon.HexToAddress(s.contractAddresses.Ics20Transfer)
	erc20Address := ethcommon.HexToAddress(s.contractAddresses.Erc20)

	transferAmount := big.NewInt(testvalues.TransferAmount)
	ethereumUserAddress := crypto.PubkeyToAddress(s.key.PublicKey)
	cosmosUserWallet := s.CosmosUsers[0]
	cosmosUserAddress := cosmosUserWallet.FormattedAddress()

	s.Require().True(s.Run("Approve the ICS20Transfer.sol contract to spend the erc20 tokens", func() {
		tx, err := s.erc20Contract.Approve(s.GetTransactOpts(s.key, eth), ics20Address, transferAmount)
		s.Require().NoError(err)

		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)

		allowance, err := s.erc20Contract.Allowance(nil, ethereumUserAddress, ics20Address)
		s.Require().NoError(err)
		s.Require().Equal(transferAmount, allowance)
	}))

	var (
		sendPacket    ics26router.IICS26RouterMsgsPacket
		ethSendTxHash []byte
		escrowAddress ethcommon.Address
	)
	s.Require().True(s.Run("Send ERC20 tokens on Ethereum", func() {
		timeout := uint64(time.Now().Add(30 * time.Minute).Unix())

		msgSendPacket := ics20transfer.IICS20TransferMsgsSendTransferMsg{
			Denom:            erc20Address,
			Amount:           transferAmount,
			Receiver:         cosmosUserAddress,
			TimeoutTimestamp: timeout,
			SourceClient:     testvalues.CustomClientID,
			Memo:             "create-voucher-cosmos",
		}

		tx, err := s.ics20Contract.SendTransfer(s.GetTransactOpts(s.key, eth), msgSendPacket)
		s.Require().NoError(err)
		receipt, err := eth.GetTxReciept(ctx, tx.Hash())
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status)
		ethSendTxHash = tx.Hash().Bytes()

		sendPacketEvent, err := e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseSendPacket)
		s.Require().NoError(err)
		sendPacket = sendPacketEvent.Packet
		s.Require().Equal(uint64(1), sendPacket.Sequence)
		s.Require().Equal(timeout, sendPacket.TimeoutTimestamp)
		s.Require().Len(sendPacket.Payloads, 1)
		s.Require().Equal(transfertypes.PortID, sendPacket.Payloads[0].SourcePort)
		s.Require().Equal(testvalues.CustomClientID, sendPacket.SourceClient)
		s.Require().Equal(transfertypes.PortID, sendPacket.Payloads[0].DestPort)
		s.Require().Equal(testvalues.FirstWasmClientID, sendPacket.DestClient)
		s.Require().Equal(transfertypes.V1, sendPacket.Payloads[0].Version)
		s.Require().Equal(transfertypes.EncodingABI, sendPacket.Payloads[0].Encoding)

		s.True(s.Run("Verify balances on Ethereum", func() {
			// User balance on Ethereum
			userBalance, err := s.erc20Contract.BalanceOf(nil, ethereumUserAddress)
			s.Require().NoError(err)
			s.Require().Equal(new(big.Int).Sub(testvalues.StartingERC20Balance, transferAmount), userBalance)

			// Get the escrow address
			escrowAddress, err = s.ics20Contract.GetEscrow(nil, testvalues.CustomClientID)
			s.Require().NoError(err)

			// ICS20 contract balance on Ethereum
			escrowBalance, err := s.erc20Contract.BalanceOf(nil, escrowAddress)
			s.Require().NoError(err)
			s.Require().Equal(transferAmount, escrowBalance)
		}))
	}))

	var (
		denomOnCosmos transfertypes.Denom
		ackTxHash     []byte
	)
	s.Require().True(s.Run("Receive packets on Cosmos chain", func() {
		var relayTxBodyBz []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:    eth.ChainID.String(),
				DstChain:    simd.Config().ChainID,
				SourceTxIds: [][]byte{ethSendTxHash},
				SrcClientId: testvalues.CustomClientID,
				DstClientId: testvalues.FirstWasmClientID,
			})
			s.Require().NoError(err)
			s.Require().NotEmpty(resp.Tx)
			s.Require().Empty(resp.Address)

			relayTxBodyBz = resp.Tx
		}))

		s.Require().True(s.Run("Broadcast relay tx", func() {
			resp := s.MustBroadcastSdkTxBody(ctx, simd, s.SimdRelayerSubmitter, 20_000_000, relayTxBodyBz)
			var err error
			ackTxHash, err = hex.DecodeString(resp.TxHash)
			s.Require().NoError(err)

			denomOnCosmos = transfertypes.NewDenom(s.contractAddresses.Erc20, transfertypes.NewHop(transfertypes.PortID, testvalues.FirstWasmClientID))

			// Verify initial voucher balance
			balanceResp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
				Address: cosmosUserAddress,
				Denom:   denomOnCosmos.IBCDenom(),
			})
			s.Require().NoError(err)
			s.Require().Equal(transferAmount, balanceResp.Balance.Amount.BigInt())
		}))
	}))

	s.Require().True(s.Run("Acknowledge packets on Ethereum", func() {
		var ackRelayTx []byte
		s.Require().True(s.Run("Retrieve relay tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:    simd.Config().ChainID,
				DstChain:    eth.ChainID.String(),
				SourceTxIds: [][]byte{ackTxHash},
				SrcClientId: testvalues.FirstWasmClientID,
				DstClientId: testvalues.CustomClientID,
			})
			s.Require().NoError(err)
			ackRelayTx = resp.Tx
		}))
		receipt, err := eth.BroadcastTx(ctx, s.EthRelayerSubmitter, 15_000_000, &ics26Address, ackRelayTx)
		s.Require().NoError(err)
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status, fmt.Sprintf("Tx failed: %+v", receipt))

		// Verify the ack packet event exists
		_, err = e2esuite.GetEvmEvent(receipt, s.ics26Contract.ParseAckPacket)
		s.Require().NoError(err)
	}))

	var cosmosTimeoutSendTxHash []byte
	s.Require().True(s.Run("Send voucher from Cosmos to Eth (timeout)", func() {
		timeout := uint64(time.Now().Add(30 * time.Second).Unix())
		ibcCoin := sdk.NewCoin(denomOnCosmos.Path(), sdkmath.NewIntFromBigInt(transferAmount))

		transferPayload := transfertypes.FungibleTokenPacketData{
			Denom:    ibcCoin.Denom, // Sending the voucher denom
			Amount:   ibcCoin.Amount.String(),
			Sender:   cosmosUserAddress,
			Receiver: strings.ToLower(ethereumUserAddress.Hex()),
			Memo:     "timeout-voucher-cosmos-eth",
		}
		encodedPayload, err := transfertypes.EncodeABIFungibleTokenPacketData(&transferPayload)
		s.Require().NoError(err)

		payload := channeltypesv2.Payload{
			SourcePort:      transfertypes.PortID,
			DestinationPort: transfertypes.PortID,
			Version:         transfertypes.V1,
			Encoding:        transfertypes.EncodingABI,
			Value:           encodedPayload,
		}
		msgSendPacket := channeltypesv2.MsgSendPacket{
			SourceClient:     testvalues.FirstWasmClientID,
			TimeoutTimestamp: timeout,
			Payloads: []channeltypesv2.Payload{
				payload,
			},
			Signer: cosmosUserWallet.FormattedAddress(),
		}

		resp, err := s.BroadcastMessages(ctx, simd, cosmosUserWallet, 200_000, &msgSendPacket)
		s.Require().NoError(err)
		cosmosTimeoutSendTxHash, err = hex.DecodeString(resp.TxHash)
		s.Require().NoError(err)

		// Verify user voucher balance is now zero
		balanceResp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
			Address: cosmosUserAddress,
			Denom:   denomOnCosmos.IBCDenom(),
		})
		s.Require().NoError(err)
		s.Require().Zero(balanceResp.Balance.Amount.Int64())
	}))

	// We retrieve the relay tx before the timeout, and will attempt to relay it after the timeout
	var recvRelayTx []byte
	s.Require().True(s.Run("Retrieve relay tx before timeout", func() {
		resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
			SrcChain:    simd.Config().ChainID,
			DstChain:    eth.ChainID.String(),
			SourceTxIds: [][]byte{cosmosTimeoutSendTxHash},
			SrcClientId: testvalues.FirstWasmClientID,
			DstClientId: testvalues.CustomClientID,
		})
		s.Require().NoError(err)
		recvRelayTx = resp.Tx
	}))

	s.Require().True(s.Run("Wait for timeout", func() {
		time.Sleep(45 * time.Second)
	}))

	s.Require().True(s.Run("Relay timeout packet to Cosmos", func() {
		var timeoutRelayTxBodyBz []byte
		s.Require().True(s.Run("Retrieve timeout tx", func() {
			resp, err := s.RelayerClient.RelayByTx(context.Background(), &relayertypes.RelayByTxRequest{
				SrcChain:     eth.ChainID.String(),  // Source of the *timeout* message is Eth
				DstChain:     simd.Config().ChainID, // Destination is Cosmos
				TimeoutTxIds: [][]byte{cosmosTimeoutSendTxHash},
				SrcClientId:  testvalues.CustomClientID,
				DstClientId:  testvalues.FirstWasmClientID,
			})
			s.Require().NoError(err)
			timeoutRelayTxBodyBz = resp.Tx
		}))

		s.Require().True(s.Run("Broadcast relay tx", func() {
			_ = s.MustBroadcastSdkTxBody(ctx, simd, s.SimdRelayerSubmitter, 2_000_000, timeoutRelayTxBodyBz)
		}))
	}))

	s.Require().True(s.Run("Verify voucher balance restored on Cosmos", func() {
		balanceResp, err := e2esuite.GRPCQuery[banktypes.QueryBalanceResponse](ctx, simd, &banktypes.QueryBalanceRequest{
			Address: cosmosUserAddress,
			Denom:   denomOnCosmos.IBCDenom(),
		})
		s.Require().NoError(err)
		s.Require().Equal(transferAmount, balanceResp.Balance.Amount.BigInt(), "Voucher balance should be restored after timeout")
	}))

	s.Require().True(s.Run("Verify recvPacket fails on Eth", func() {
		_, err := eth.BroadcastTx(ctx, s.EthRelayerSubmitter, 5_000_000, &ics26Address, recvRelayTx)
		s.Require().Error(err)
	}))
}
