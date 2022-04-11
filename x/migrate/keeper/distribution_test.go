package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/functionx/fx-core/app/fxcore"
	fxtypes "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/migrate/keeper"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
	"time"
)

func TestMigrateDistributionHandler(t *testing.T) {

	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := fxcore.GenerateGenesisValidator(3,
		sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, initBalances)))
	app := fxcore.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	delegateAddressArr := fxcore.AddTestAddrsIncremental(app, ctx, 4, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10000)))

	valA, valB, valC := validator.Validators[0], validator.Validators[1], validator.Validators[2]

	addA, addB, addC, addD := delegateAddressArr[0], delegateAddressArr[1], delegateAddressArr[2], delegateAddressArr[3]

	ctx = commitBlock(t, ctx, app)

	validatorA, found := app.StakingKeeper.GetValidator(ctx, valA.Address.Bytes())
	require.True(t, found)

	validatorB, found := app.StakingKeeper.GetValidator(ctx, valB.Address.Bytes())
	require.True(t, found)

	validatorC, found := app.StakingKeeper.GetValidator(ctx, valC.Address.Bytes())
	require.True(t, found)

	_, err := app.StakingKeeper.Delegate(ctx, addA, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000)), stakingtypes.Unbonded, validatorA, true)
	require.NoError(t, err)

	_, err = app.StakingKeeper.Delegate(ctx, addC, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000)), stakingtypes.Unbonded, validatorB, true)
	require.NoError(t, err)

	_, err = app.StakingKeeper.Delegate(ctx, addD, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000)), stakingtypes.Unbonded, validatorC, true)
	require.NoError(t, err)

	ctx = commitBlock(t, ctx, app)

	delAAndValAReward := getDelegationRewards(t, ctx, app.StakingKeeper, app.DistrKeeper, valA.Address.Bytes(), addA, true)
	require.False(t, delAAndValAReward.IsZero())

	migrateBeforeVerify(t, ctx, app.DistrKeeper, valA.Address.Bytes(), addA, addB, true)
	migrateBeforeVerify(t, ctx, app.DistrKeeper, valA.Address.Bytes(), addB, addB, false)

	migrateKeeper := app.MigrateKeeper
	m := keeper.NewDistrStakingMigrate(app.GetKey(distrtypes.StoreKey), app.GetKey(stakingtypes.StoreKey), app.StakingKeeper)
	err = m.Validate(ctx, migrateKeeper, addA, addB)
	require.NoError(t, err)
	err = m.Execute(ctx, migrateKeeper, addA, addB)
	require.NoError(t, err)

	migrateAfterVerify(t, ctx, app.DistrKeeper, valA.Address.Bytes(), addA, addB, false, true)

	ctx = commitBlock(t, ctx, app)
	delAAndValAReward = getDelegationRewards(t, ctx, app.StakingKeeper, app.DistrKeeper, valA.Address.Bytes(), addA, false)
	require.True(t, delAAndValAReward.IsZero())

	delAAndValAReward = getDelegationRewards(t, ctx, app.StakingKeeper, app.DistrKeeper, valA.Address.Bytes(), addB, true)
	require.False(t, delAAndValAReward.IsZero())
}

func commitBlock(t *testing.T, ctx sdk.Context, app *fxcore.App) sdk.Context {
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(5 * time.Second))

	staking.EndBlocker(ctx, app.StakingKeeper)
	mint.BeginBlocker(ctx, app.MintKeeper)

	distribution.BeginBlocker(ctx, abcitypes.RequestBeginBlock{
		Hash:   nil,
		Header: tmproto.Header{},
		LastCommitInfo: abcitypes.LastCommitInfo{
			Round: 0,
			Votes: buildCommitVotes(t, ctx, app.StakingKeeper, app.AppCodec()),
		},
		ByzantineValidators: nil,
	}, app.DistrKeeper)

	return ctx
}

func getDelegationRewards(t *testing.T, ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, distrKeeper distrkeeper.Keeper, val sdk.ValAddress, del sdk.AccAddress, expectFoundDelete bool) sdk.DecCoins {
	t.Helper()
	validator, found := stakingKeeper.GetValidator(ctx, val)
	require.True(t, found)

	delegation, found := stakingKeeper.GetDelegation(ctx, del, val)
	require.EqualValues(t, expectFoundDelete, found)

	if !expectFoundDelete {
		return sdk.DecCoins{}
	}
	info := distrKeeper.HasDelegatorStartingInfo(ctx, val, del)

	if !info {
		return sdk.DecCoins{}
	}

	endingPeriod := distrKeeper.IncrementValidatorPeriod(ctx, validator)

	return distrKeeper.CalculateDelegationRewards(ctx, validator, delegation, endingPeriod)
}

func migrateBeforeVerify(t *testing.T, ctx sdk.Context, distrKeeper distrkeeper.Keeper, validator sdk.ValAddress, fromDelegate sdk.AccAddress, toDelegate sdk.AccAddress, fromExists bool) {
	t.Helper()
	require.EqualValues(t, fromExists, distrKeeper.HasDelegatorStartingInfo(ctx, validator, fromDelegate))
	require.EqualValues(t, false, distrKeeper.HasDelegatorStartingInfo(ctx, validator, toDelegate))
}

func migrateAfterVerify(t *testing.T, ctx sdk.Context, distrKeeper distrkeeper.Keeper, validator sdk.ValAddress, fromDelegate sdk.AccAddress, toDelegate sdk.AccAddress, fromExists, toExists bool) {
	t.Helper()
	require.EqualValues(t, fromExists, distrKeeper.HasDelegatorStartingInfo(ctx, validator, fromDelegate))
	require.EqualValues(t, toExists, distrKeeper.HasDelegatorStartingInfo(ctx, validator, toDelegate))
}

func buildCommitVotes(t *testing.T, ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, codec codec.Codec) []abcitypes.VoteInfo {
	t.Helper()
	validators := stakingKeeper.GetAllValidators(ctx)

	var result []abcitypes.VoteInfo
	for _, validator := range validators {
		if !validator.IsBonded() {
			continue
		}

		var pubkey cryptotypes.PubKey
		err := codec.UnpackAny(validator.ConsensusPubkey, &pubkey)
		require.NoError(t, err)
		result = append(result, abcitypes.VoteInfo{
			Validator: abcitypes.Validator{
				Address: pubkey.Address(),
				Power:   validator.GetConsensusPower(),
			},
			SignedLastBlock: true,
		})
	}
	return result
}
