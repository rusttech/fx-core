package keeper_test

import (
	_ "embed"
	"github.com/functionx/fx-core/app/helpers"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	app "github.com/functionx/fx-core/app"
	fxtypes "github.com/functionx/fx-core/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/crypto/ethsecp256k1"
	"github.com/functionx/fx-core/tests"
	"github.com/functionx/fx-core/x/feemarket/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.App
	queryClient types.QueryClient
	address     common.Address
	consAddress sdk.ConsAddress

	// for generate test tx
	clientCtx client.Context
	ethSigner ethtypes.Signer

	appCodec codec.Codec
	signer   keyring.Signer
	checkTx  bool
}

/// DoSetupTest setup test environment, it uses`require.TestingT` to support both `testing.T` and `testing.B`.
func (suite *KeeperTestSuite) DoSetupTest(t require.TestingT) {
	fxtypes.ChangeNetworkForTest(fxtypes.NetworkDevnet())

	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.address = common.BytesToAddress(priv.PubKey().Address().Bytes())
	suite.signer = tests.NewSigner(priv)

	// consensus key
	priv, err = ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.consAddress = sdk.ConsAddress(priv.PubKey().Address())

	suite.app = helpers.Setup(suite.T(), false, 1)
	suite.ctx = suite.app.BaseApp.NewContext(suite.checkTx, tmproto.Header{
		Height:          fxtypes.EvmV1SupportBlock(),
		ChainID:         "fxcore",
		Time:            time.Now().UTC(),
		ProposerAddress: suite.consAddress.Bytes(),
		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.FeeMarketKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	acc := authtypes.NewBaseAccount(suite.address.Bytes(), nil, 0, 0)

	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
	suite.app.EvmKeeper.SetAddressCode(suite.ctx, suite.address, common.BytesToHash(crypto.Keccak256(nil)).Bytes())

	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, priv.PubKey(), stakingtypes.Description{})
	suite.Require().NoError(err)
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	suite.Require().NoError(err)
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	suite.Require().NoError(err)
	suite.app.StakingKeeper.SetValidator(suite.ctx, validator)

	encodingConfig := app.MakeEncodingConfig()
	suite.clientCtx = client.Context{}.WithTxConfig(encodingConfig.TxConfig)
	suite.ethSigner = ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())
	suite.appCodec = encodingConfig.Marshaler

	suite.ctx = suite.ctx.WithBlockHeight(fxtypes.EvmV1SupportBlock())
	//forks.UpdateMetadata(suite.ctx, suite.app.BankKeeper)
	//require.NoError(suite.T(), forks.InitSupportEvm(suite.ctx, suite.app.AccountKeeper,
	//	suite.app.FeeMarketKeeper, types.DefaultParams(),
	//	suite.app.EvmKeeper, evmtypes.DefaultParams(),
	//	suite.app.Erc20Keeper, erc20types.DefaultParams()))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetGetBlockGasUsed() {
	testCases := []struct {
		name     string
		malleate func()
		expGas   uint64
	}{
		{
			"with last block given",
			func() {
				suite.app.FeeMarketKeeper.SetBlockGasUsed(suite.ctx, uint64(1000000))
			},
			uint64(1000000),
		},
	}
	for _, tc := range testCases {
		tc.malleate()

		gas := suite.app.FeeMarketKeeper.GetBlockGasUsed(suite.ctx)
		suite.Require().Equal(tc.expGas, gas, tc.name)
	}
}

func (suite *KeeperTestSuite) TestSetGetGasFee() {
	testCases := []struct {
		name     string
		malleate func()
		expFee   *big.Int
	}{
		{
			"with last block given",
			func() {
				suite.app.FeeMarketKeeper.SetBaseFee(suite.ctx, sdk.OneDec().BigInt())
			},
			sdk.OneDec().BigInt(),
		},
	}

	for _, tc := range testCases {
		tc.malleate()

		fee := suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx)
		suite.Require().Equal(tc.expFee, fee, tc.name)
	}
}
