package keeper_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	ibctmtypes "github.com/cosmos/cosmos-sdk/x/ibc/light-clients/07-tendermint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/functionx/fx-core/app/fxcore"
	"github.com/functionx/fx-core/crypto/ethsecp256k1"
	"github.com/functionx/fx-core/server/config"
	"github.com/functionx/fx-core/tests"
	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	"github.com/functionx/fx-core/x/crosschain"
	crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	evmkeeper "github.com/functionx/fx-core/x/evm/keeper"
	evm "github.com/functionx/fx-core/x/evm/types"
	"github.com/functionx/fx-core/x/gravity"
	gravitytypes "github.com/functionx/fx-core/x/gravity/types"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"math"
	"math/big"
	"testing"
)

type IBCTransferSimulate struct {
	T *testing.T
}

func (it *IBCTransferSimulate) Transfer(goCtx context.Context, msg *ibctransfertypes.MsgTransfer) (*ibctransfertypes.MsgTransferResponse, error) {
	//it.T.Logf("ibc transfer simulate sender %s, receiver %s, amount %s, fee %s", msg.Sender, msg.Receiver, msg.Token.String(), msg.Fee.String())
	return &ibctransfertypes.MsgTransferResponse{}, nil
}

func (it *IBCTransferSimulate) GetRouter() *ibctransfertypes.Router {
	router := ibctransfertypes.NewRouter()

	return router
}

type IBCChannelSimulate struct {
}

func (ic *IBCChannelSimulate) GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error) {
	return "", &ibctmtypes.ClientState{
		ChainId:         "fxcore",
		TrustLevel:      ibctmtypes.Fraction{},
		TrustingPeriod:  0,
		UnbondingPeriod: 0,
		MaxClockDrift:   0,
		FrozenHeight: ibcclienttypes.Height{
			RevisionHeight: 1000,
			RevisionNumber: 1000,
		},
		LatestHeight: ibcclienttypes.Height{
			RevisionHeight: 10,
			RevisionNumber: 10,
		},
		ProofSpecs:                   nil,
		UpgradePath:                  nil,
		AllowUpdateAfterExpiry:       false,
		AllowUpdateAfterMisbehaviour: false,
	}, nil
}

func (ic *IBCChannelSimulate) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	return 1, true
}

var (
	wfxMetadata = banktypes.Metadata{
		Description: "Wrap Function X",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "FX",
				Exponent: 0,
				Aliases:  nil,
			},
			{
				Denom:    "WFX",
				Exponent: 18,
				Aliases:  nil,
			},
		},
		Base:    "FX",
		Display: "WFX",
	}

	purseMetadata = banktypes.Metadata{
		Description: "Pundi X Purse Token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    purseDenom,
				Exponent: 0,
				Aliases:  nil,
			},
			{
				Denom:    "PURSE",
				Exponent: 18,
				Aliases:  nil,
			},
		},
		Base:    purseDenom,
		Display: "PURSE",
	}

	purseDenom = "ibc/5BAB702195E8411500FC0256EBC717AB9C7DC039D7FFC6E1A471022AA939E600"
)

func TestHookChainGravity(t *testing.T) {
	app, validators, _, delegateAddressArr := initTest(t)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{ProposerAddress: validators[0].Address})
	require.NoError(t, InitEvmModuleParams(ctx, app.EvmKeeper, true))

	pair, err := app.IntrarelayerKeeper.RegisterCoin(ctx, wfxMetadata)
	require.NoError(t, err)

	val := validators[0]
	validator := GetValidator(t, app, val)[0]
	del := delegateAddressArr[0]

	ctx = ctx.WithBlockHeight(504000)

	signer1, addr1 := privateSigner()
	_, addr2 := privateSigner()
	amt := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(100))
	err = app.BankKeeper.SendCoins(ctx, del, sdk.AccAddress(addr1.Bytes()), sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, amt)))
	require.NoError(t, err)

	ctx = testInitGravity(t, ctx, app, validator.GetOperator(), addr1.Bytes(), addr2)

	balances := app.BankKeeper.GetAllBalances(ctx, addr1.Bytes())
	t.Log("b", balances.String())

	err = app.IntrarelayerKeeper.ConvertDenomToFIP20(ctx, addr1.Bytes(), addr1, sdk.NewCoin(fxcore.MintDenom, sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(10))))
	require.NoError(t, err)

	balanceOf, err := app.IntrarelayerKeeper.QueryFIP20BalanceOf(ctx, pair.GetFIP20Contract(), addr1)
	require.NoError(t, err)
	t.Log("c", balanceOf.String())

	token := pair.GetFIP20Contract()
	crossChainTarget := fmt.Sprintf("%s%s", contracts.TransferChainPrefix, gravitytypes.ModuleName)
	transferChainData := packTransferCrossData(t, addr2.String(), fxcore.CoinOne, fxcore.CoinOne, crossChainTarget)
	sendEthTx(t, ctx, app, signer1, addr1, token, transferChainData)

	transactions := app.GravityKeeper.GetPoolTransactions(ctx)
	for _, tx := range transactions {
		t.Log("sender", tx.Sender, "dest", tx.DestAddress, "amount", tx.Erc20Token.String())
	}
}

func TestHookChainBSC(t *testing.T) {
	app, validators, genesisAccount, delegateAddressArr := initTest(t)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{ProposerAddress: validators[0].Address})
	require.NoError(t, InitEvmModuleParams(ctx, app.EvmKeeper, true))

	pair, err := app.IntrarelayerKeeper.RegisterCoin(ctx, purseMetadata)
	require.NoError(t, err)

	fip20, err := app.IntrarelayerKeeper.QueryFIP20(ctx, pair.GetFIP20Contract())
	require.NoError(t, err)
	_ = fip20
	//t.Log("fip20", fip20.Name, fip20.Symbol, fip20.Decimals)

	del := delegateAddressArr[0]
	ga := genesisAccount[0]

	ctx = ctx.WithBlockHeight(504000)

	signer1, addr1 := privateSigner()
	_, addr2 := privateSigner()
	amt := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(100))

	err = app.BankKeeper.SendCoins(ctx, ga.GetAddress(), sdk.AccAddress(addr1.Bytes()), sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, amt), sdk.NewCoin(purseDenom, amt)))
	require.NoError(t, err)

	ctx = testInitBscCrossChain(t, ctx, app, del, addr1.Bytes(), addr2)

	balances := app.BankKeeper.GetAllBalances(ctx, addr1.Bytes())
	t.Log("balance", balances.String())

	err = app.IntrarelayerKeeper.ConvertDenomToFIP20(ctx, addr1.Bytes(), addr1, sdk.NewCoin(purseDenom, sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(10))))
	require.NoError(t, err)

	balanceOf, err := app.IntrarelayerKeeper.QueryFIP20BalanceOf(ctx, pair.GetFIP20Contract(), addr1)
	require.NoError(t, err)
	t.Log("balanceOf", balanceOf.String())

	token := pair.GetFIP20Contract()
	crossChainTarget := fmt.Sprintf("%s%s", contracts.TransferChainPrefix, bsctypes.ModuleName)
	transferChainData := packTransferCrossData(t, addr2.String(), fxcore.CoinOne, fxcore.CoinOne, crossChainTarget)
	sendEthTx(t, ctx, app, signer1, addr1, token, transferChainData)

	transactions := app.BscKeeper.GetUnbatchedTransactions(ctx)
	for _, tx := range transactions {
		t.Log("sender", tx.Sender, "dest", tx.DestAddress, "amount", tx.Token.String())
	}
}

func TestHookIBC(t *testing.T) {
	app, validators, _, delegateAddressArr := initTest(t)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{ProposerAddress: validators[0].Address})
	require.NoError(t, InitEvmModuleParams(ctx, app.EvmKeeper, true))

	pair, err := app.IntrarelayerKeeper.RegisterCoin(ctx, wfxMetadata)
	require.NoError(t, err)

	//validator := GetValidator(t, app, val)[0]
	//val := validators[0]
	del := delegateAddressArr[0]

	ctx = ctx.WithBlockHeight(504000)

	signer1, addr1 := privateSigner()
	//_, addr2 := privateSigner()
	amt := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(100))
	err = app.BankKeeper.SendCoins(ctx, del, sdk.AccAddress(addr1.Bytes()), sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, amt)))
	require.NoError(t, err)

	balances := app.BankKeeper.GetAllBalances(ctx, addr1.Bytes())
	t.Log("balance", balances.String())

	err = app.IntrarelayerKeeper.ConvertDenomToFIP20(ctx, addr1.Bytes(), addr1, sdk.NewCoin(fxcore.MintDenom, sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(10))))
	require.NoError(t, err)

	balanceOf, err := app.IntrarelayerKeeper.QueryFIP20BalanceOf(ctx, pair.GetFIP20Contract(), addr1)
	require.NoError(t, err)
	t.Log("balanceOf", balanceOf.String())

	//reset ibc
	app.IntrarelayerKeeper.SetIBCTransferKeeper(&IBCTransferSimulate{T: t})
	app.IntrarelayerKeeper.SetIBCChannelKeeper(&IBCChannelSimulate{})

	evmHooks := evmkeeper.NewMultiEvmHooks(app.IntrarelayerKeeper)
	app.EvmKeeper = app.EvmKeeper.SetHooks(evmHooks)

	token := pair.GetFIP20Contract()
	ibcTarget := fmt.Sprintf("%s%s", contracts.TransferIBCPrefix, "px/transfer/channel-0")
	transferIBCData := packTransferCrossData(t, "px16u6kjunrcxkvaln9aetxwjpruply3sgwpr9z8u", fxcore.CoinOne, big.NewInt(0), ibcTarget)
	sendEthTx(t, ctx, app, signer1, addr1, token, transferIBCData)
}

func packTransferCrossData(t *testing.T, to string, amount, fee *big.Int, target string) []byte {
	pack, err := contracts.FIP20Contract.ABI.Pack("transferCross", to, amount, fee, target)
	require.NoError(t, err)
	return pack
}

func sendEthTx(t *testing.T, ctx sdk.Context, app *fxcore.App,
	signer keyring.Signer, from, contract common.Address, data []byte) {

	chainID := app.EvmKeeper.ChainID()

	args, err := json.Marshal(&evm.TransactionArgs{To: &contract, From: &from, Data: (*hexutil.Bytes)(&data)})
	require.NoError(t, err)
	res, err := app.EvmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &evm.EthCallRequest{
		Args:   args,
		GasCap: uint64(config.DefaultGasCap),
	})
	require.NoError(t, err)

	nonce, err := app.AccountKeeper.GetSequence(ctx, from.Bytes())
	require.NoError(t, err)

	ercTransferTx := evm.NewTx(
		chainID,
		nonce,
		&contract,
		nil,
		res.Gas,
		nil,
		app.FeeMarketKeeper.GetBaseFee(ctx),
		big.NewInt(1),
		data,
		&ethtypes.AccessList{}, // accesses
	)

	ercTransferTx.From = from.String()
	err = ercTransferTx.Sign(ethtypes.LatestSignerForChainID(chainID), signer)
	require.NoError(t, err)
	rsp, err := app.EvmKeeper.EthereumTx(sdk.WrapSDKContext(ctx), ercTransferTx)
	require.NoError(t, err)
	require.Empty(t, rsp.VmError)
}

func privateSigner() (keyring.Signer, common.Address) {
	// account key
	priKey := NewPriKey()
	//ethsecp256k1.GenerateKey()
	ethPriv := &ethsecp256k1.PrivKey{Key: priKey.Bytes()}

	return tests.NewSigner(ethPriv), common.BytesToAddress(ethPriv.PubKey().Address())
}

func initTest(t *testing.T) (*fxcore.App, []*tmtypes.Validator, authtypes.GenesisAccounts, []sdk.AccAddress) {
	initBalances := sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(1,
		sdk.NewCoins(sdk.NewCoin(fxcore.MintDenom, initBalances),
			sdk.NewCoin("eth0xd9EEd31F5731DfC3Ca18f09B487e200F50a6343B", initBalances),
			sdk.NewCoin("ibc/4757BC3AA2C696F7083C825BD3951AE3D1631F2A272EA7AFB9B3E1CCCA8560D4", initBalances),
			sdk.NewCoin(purseDenom, initBalances)))

	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	delegateAddressArr := fxcore.AddTestAddrsIncremental(app, ctx, 1, sdk.NewIntFromBigInt(fxcore.CoinOne).Mul(sdk.NewInt(1000000000)))

	return app, validator.Validators, genesisAccounts, delegateAddressArr
}

func GetValidator(t *testing.T, app *fxcore.App, vals ...*tmtypes.Validator) []stakingtypes.Validator {
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	validators := make([]stakingtypes.Validator, 0, len(vals))
	for _, val := range vals {
		validator, found := app.StakingKeeper.GetValidator(ctx, val.Address.Bytes())
		require.True(t, found)
		validators = append(validators, validator)
	}
	return validators
}

var (
	FxOriginatedTokenContract = common.HexToAddress("0x0000000000000000000000000000000000000000")
	BSCBridgeTokenContract    = common.HexToAddress("0x0000000000000000000000000000000000000001")
)

func TestContractAddr(t *testing.T) {
	t.Log(FxOriginatedTokenContract.String())
	t.Log(BSCBridgeTokenContract.String())
}

func testInitGravity(t *testing.T, ctx sdk.Context, app *fxcore.App, val sdk.ValAddress, orch sdk.AccAddress, addr common.Address) sdk.Context {
	app.GravityKeeper.SetOrchestratorValidator(ctx, val, orch)
	app.GravityKeeper.SetEthAddressForValidator(ctx, val, addr.String())

	testValSetUpdateClaim(t, ctx, app, orch, addr)

	testFxOriginatedTokenClaim(t, ctx, app, orch)

	gravity.EndBlocker(ctx, *app.GravityKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	return ctx
}

func testFxOriginatedTokenClaim(t *testing.T, ctx sdk.Context, app *fxcore.App, orch sdk.AccAddress) {
	msg := &gravitytypes.MsgFxOriginatedTokenClaim{
		EventNonce:    2,
		BlockHeight:   uint64(ctx.BlockHeight()),
		TokenContract: FxOriginatedTokenContract.String(),
		Name:          "Function X",
		Symbol:        "FX",
		Decimals:      18,
		Orchestrator:  orch.String(),
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = app.GravityKeeper.Attest(ctx, msg, any)
	require.NoError(t, err)
}

func testValSetUpdateClaim(t *testing.T, ctx sdk.Context, app *fxcore.App, orch sdk.AccAddress, addr common.Address) {
	msg := &gravitytypes.MsgValsetUpdatedClaim{
		EventNonce:  1,
		BlockHeight: uint64(ctx.BlockHeight()),
		ValsetNonce: 0,
		Members: []*gravitytypes.BridgeValidator{
			{
				Power:      uint64(math.MaxUint32),
				EthAddress: addr.String(),
			},
		},
		Orchestrator: orch.String(),
	}

	for _, member := range msg.Members {
		memberVal := app.GravityKeeper.GetValidatorByEthAddress(ctx, member.EthAddress)
		require.NotEmpty(t, memberVal)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = app.GravityKeeper.Attest(ctx, msg, any)
	require.NoError(t, err)
}

func testInitBscCrossChain(t *testing.T, ctx sdk.Context, app *fxcore.App, oracleAddress, orchestratorAddr sdk.AccAddress, externalAddress common.Address) sdk.Context {
	deposit := sdk.NewCoin(fxcore.MintDenom, sdk.NewIntFromBigInt(big.NewInt(0).Mul(big.NewInt(10), fxcore.CoinOne)))
	err := app.BankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddress, app.BscKeeper.GetModuleName(), sdk.NewCoins(deposit))
	require.NoError(t, err)

	testBSCParamsProposal(t, ctx, app, oracleAddress)

	oracle := crosschaintypes.Oracle{
		OracleAddress:       oracleAddress.String(),
		OrchestratorAddress: orchestratorAddr.String(),
		ExternalAddress:     externalAddress.String(),
		DepositAmount:       deposit,
		StartHeight:         ctx.BlockHeight(),
		Jailed:              false,
		JailedHeight:        0,
	}
	// save oracle
	app.BscKeeper.SetOracle(ctx, oracle)

	app.BscKeeper.SetOracleByOrchestrator(ctx, oracleAddress, orchestratorAddr)
	// set the ethereum address
	app.BscKeeper.SetExternalAddressForOracle(ctx, oracleAddress, externalAddress.String())
	// save total deposit amount
	totalDeposit := app.BscKeeper.GetTotalDeposit(ctx)
	app.BscKeeper.SetTotalDeposit(ctx, totalDeposit.Add(deposit))

	app.BscKeeper.CommonSetOracleTotalPower(ctx)

	testBSCOracleSetUpdateClaim(t, ctx, app, orchestratorAddr, externalAddress)

	testBSCBridgeTokenClaim(t, ctx, app, orchestratorAddr)

	crosschain.EndBlocker(ctx, app.BscKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	return ctx
}

func testBSCParamsProposal(t *testing.T, ctx sdk.Context, app *fxcore.App, oracles sdk.AccAddress) {
	slashFraction, _ := sdk.NewDecFromStr("0.001")
	oracleSetUpdatePowerChangePercent, _ := sdk.NewDecFromStr("0.1")
	proposal := &crosschaintypes.InitCrossChainParamsProposal{
		Title:       "bsc cross chain",
		Description: "bsc cross chain init",
		Params: &crosschaintypes.Params{
			GravityId:                         "fx-bsc-bridge",
			SignedWindow:                      20000,
			ExternalBatchTimeout:              86400000,
			AverageBlockTime:                  1000,
			AverageExternalBlockTime:          3000,
			SlashFraction:                     slashFraction,
			OracleSetUpdatePowerChangePercent: oracleSetUpdatePowerChangePercent,
			IbcTransferTimeoutHeight:          20000,
			Oracles:                           []string{oracles.String()},
			DepositThreshold:                  sdk.NewCoin(fxcore.MintDenom, sdk.NewIntFromBigInt(big.NewInt(0).Mul(big.NewInt(10), fxcore.CoinOne))),
		},
		ChainName: app.BscKeeper.GetModuleName(),
	}

	k := crosschainkeeper.EthereumMsgServer{Keeper: app.BscKeeper}
	err := k.HandleInitCrossChainParamsProposal(ctx, proposal)
	require.NoError(t, err)
}

func testBSCBridgeTokenClaim(t *testing.T, ctx sdk.Context, app *fxcore.App, orchAddr sdk.AccAddress) {
	msg := &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:    2,
		BlockHeight:   uint64(ctx.BlockHeight()),
		TokenContract: BSCBridgeTokenContract.String(),
		Name:          "PURSE Token",
		Symbol:        "PURSE",
		Decimals:      18,
		Orchestrator:  orchAddr.String(),
		ChannelIbc:    hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:     app.BscKeeper.GetModuleName(),
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = app.BscKeeper.Attest(ctx, msg, any)
	require.NoError(t, err)
}

func testBSCOracleSetUpdateClaim(t *testing.T, ctx sdk.Context, app *fxcore.App, orch sdk.AccAddress, addr common.Address) {
	msg := &crosschaintypes.MsgOracleSetUpdatedClaim{
		EventNonce:     1,
		BlockHeight:    uint64(ctx.BlockHeight()),
		OracleSetNonce: 0,
		Members: crosschaintypes.BridgeValidators{
			{
				Power:           uint64(math.MaxUint32),
				ExternalAddress: addr.String(),
			},
		},
		Orchestrator: orch.String(),
		ChainName:    app.BscKeeper.GetModuleName(),
	}
	for _, member := range msg.Members {
		_, found := app.BscKeeper.GetOracleByExternalAddress(ctx, member.ExternalAddress)
		require.True(t, found)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = app.BscKeeper.Attest(ctx, msg, any)
	require.NoError(t, err)
}