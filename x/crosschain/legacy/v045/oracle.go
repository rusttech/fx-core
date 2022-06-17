package v045

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	fxtypes "github.com/functionx/fx-core/types"
	v042 "github.com/functionx/fx-core/x/crosschain/legacy/v042"
	"github.com/functionx/fx-core/x/crosschain/types"
)

func MigrateOracle(ctx sdk.Context, cdc codec.BinaryCodec, storeKey sdk.StoreKey, stakingKeeper StakingKeeper) (types.Oracles, stakingtypes.Validator, error) {

	validatorsByPower := stakingKeeper.GetBondedValidatorsByPower(ctx)
	if len(validatorsByPower) <= 0 {
		panic("no found bonded validator")
	}
	delegateValidator := validatorsByPower[0]

	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleKey)
	defer iterator.Close()

	oracles := types.Oracles{}
	for ; iterator.Valid(); iterator.Next() {
		var legacyOracle v042.LegacyOracle
		cdc.MustUnmarshal(iterator.Value(), &legacyOracle)
		if legacyOracle.DelegateAmount.Denom != fxtypes.DefaultDenom {
			return nil, delegateValidator, sdkerrors.Wrapf(types.ErrInvalid, "delegate denom: %s", legacyOracle.DelegateAmount.Denom)
		}

		oracle := types.Oracle{
			OracleAddress:     legacyOracle.OracleAddress,
			BridgerAddress:    legacyOracle.BridgerAddress,
			ExternalAddress:   legacyOracle.ExternalAddress,
			DelegateAmount:    legacyOracle.DelegateAmount.Amount,
			StartHeight:       legacyOracle.StartHeight,
			Online:            !legacyOracle.Jailed,
			DelegateValidator: delegateValidator.OperatorAddress,
			SlashTimes:        0,
		}
		store.Set(types.GetOracleKey(oracle.GetOracle()), cdc.MustMarshal(&oracle))
		oracles = append(oracles, oracle)
	}
	return oracles, delegateValidator, nil
}
