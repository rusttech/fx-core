package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	fxtype "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/migrate/types"
)

var (
	_ types.MsgServer = &Keeper{}
)

func (k Keeper) MigrateAccount(goCtx context.Context, msg *types.MsgMigrateAccount) (*types.MsgMigrateAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//check module enable
	if ctx.BlockHeight() < fxtype.MigrateSupportBlock() || !k.intrarelayerKeeper.HasInit(ctx) {
		return nil, sdkerrors.Wrap(types.InvalidRequest, "migrate module not enable")
	}

	fromAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	toAddress, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, err
	}
	//migrated
	if k.HasMigrateRecord(ctx, fromAddress) {
		return nil, sdkerrors.Wrapf(types.ErrAlreadyMigrate, "address %s has been migrated", msg.From)
	}
	if k.HasMigrateRecord(ctx, toAddress) {
		return nil, sdkerrors.Wrapf(types.ErrAlreadyMigrate, "address %s has been migrated", msg.To)
	}

	//migrate Validate
	for _, m := range k.GetMigrateI() {
		if err := m.Validate(ctx, k, fromAddress, toAddress); err != nil {
			return nil, sdkerrors.Wrap(types.ErrMigrateValidate, err.Error())
		}
	}
	//migrate Execute
	for _, m := range k.GetMigrateI() {
		err := m.Execute(ctx, k, fromAddress, toAddress)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrMigrateExecute, err.Error())
		}
	}
	//set record
	k.SetMigrateRecord(ctx, fromAddress, toAddress)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMigrate,
			sdk.NewAttribute(types.AttributeKeyFrom, msg.From),
			sdk.NewAttribute(types.AttributeKeyTo, msg.To),
		),
	})
	return &types.MsgMigrateAccountResponse{}, nil
}