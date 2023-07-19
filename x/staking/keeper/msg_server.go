package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v5/x/staking/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) GrantPrivilege(goCtx context.Context, msg *types.MsgGrantPrivilege) (*types.MsgGrantPrivilegeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	fromAddr := sdk.MustAccAddressFromBech32(msg.FromAddress)
	pk, err := types.ProtoAnyToAccountPubKey(msg.ToPubkey)
	if err != nil {
		return nil, err
	}
	toAddress := sdk.AccAddress(pk.Address())

	// 1. validator
	if _, found := k.GetValidator(ctx, valAddr); !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "validator %s not found", msg.ValidatorAddress)
	}
	// 2. from authorized
	if !k.HasValidatorGrant(ctx, fromAddr, valAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "from address not authorized")
	}
	// 3. revoke old privilege
	if err = k.RevokeAuthorization(ctx, fromAddr, sdk.AccAddress(valAddr)); err != nil {
		return nil, err
	}
	// 4. grant new privilege
	genericGrant := []authz.Authorization{authz.NewGenericAuthorization(sdk.MsgTypeURL(&authz.MsgGrant{}))}
	if err = k.GrantAuthorization(ctx, toAddress, sdk.AccAddress(valAddr), genericGrant, types.GrantExpirationTime); err != nil {
		return nil, err
	}
	// 5. update validator operator
	k.UpdateValidatorOperator(ctx, valAddr, toAddress)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGrantPrivilege,
		sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
		sdk.NewAttribute(types.AttributeKeyFrom, msg.FromAddress),
		sdk.NewAttribute(types.AttributeKeyTo, toAddress.String()),
	))
	return &types.MsgGrantPrivilegeResponse{}, nil
}

func (k Keeper) EditConsensusPubKey(goCtx context.Context, msg *types.MsgEditConsensusPubKey) (*types.MsgEditConsensusPubKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	fromAddr := sdk.MustAccAddressFromBech32(msg.From)
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "validator %s not found", msg.ValidatorAddress)
	}
	// authorized from address
	if !k.HasValidatorGrant(ctx, fromAddr, valAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "from address not authorized")
	}
	// check validator is updating consensus pubkey
	if k.HasConsensusPubKey(ctx, valAddr) || k.HasConsensusProcess(ctx, valAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "validator %s is updating consensus pubkey", msg.ValidatorAddress)
	}

	// pubkey and address
	newPubKey, err := k.validateAnyPubKey(ctx, msg.Pubkey)
	if err != nil {
		return nil, err
	}
	newConsAddr := sdk.ConsAddress(newPubKey.Address())

	newPkFound := false
	totalUpdatePower := math.NewInt(validator.ConsensusPower(k.PowerReduction(ctx)))
	k.IteratorConsensusPubKey(ctx, func(addr sdk.ValAddress, pkBytes []byte) bool {
		var pk cryptotypes.PubKey
		if err := k.cdc.UnmarshalInterfaceJSON(pkBytes, &pk); err != nil {
			k.Logger(ctx).Error("failed to unmarshal pubKey", "validator", valAddr.String(), "err", err.Error())
			return false
		}
		if newConsAddr.Equals(sdk.ConsAddress(pk.Address())) {
			newPkFound = true
			return true
		}
		power := k.GetLastValidatorPower(ctx, addr)
		totalUpdatePower = totalUpdatePower.Add(math.NewInt(power))
		return false
	})

	if newPkFound { // new pk already exists
		return nil, stakingtypes.ErrValidatorPubKeyExists.Wrapf("new consensus pubkey %s already exists", newConsAddr.String())
	}
	totalPowerOneThird := k.GetLastTotalPower(ctx).QuoRaw(3) // less than 1/3 total power
	if totalUpdatePower.GTE(totalPowerOneThird) {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("total update power %s more than 1/3 total power %s",
			totalUpdatePower.String(), totalPowerOneThird.String())
	}

	// set validator new consensus pubkey
	if err = k.SetConsensusPubKey(ctx, valAddr, newPubKey); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEditConsensusPubKey,
		sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
		sdk.NewAttribute(types.AttributeKeyFrom, msg.From),
		sdk.NewAttribute(types.AttributeKeyPubKey, newPubKey.String()),
	))

	return &types.MsgEditConsensusPubKeyResponse{}, err
}

func (k Keeper) validateAnyPubKey(ctx sdk.Context, pubkey *codectypes.Any) (cryptotypes.PubKey, error) {
	// pubkey type
	pk, ok := pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pubkey.GetCachedValue())
	}
	// pubkey exist
	newConsAddr := sdk.GetConsAddress(pk)
	if _, found := k.GetValidatorByConsAddr(ctx, newConsAddr); found {
		return nil, stakingtypes.ErrValidatorPubKeyExists
	}
	// pubkey type supported
	cp := ctx.ConsensusParams()
	if cp != nil && cp.Validator != nil {
		pkType := pk.Type()
		hasKeyType := false
		for _, keyType := range cp.Validator.PubKeyTypes {
			if pkType == keyType {
				hasKeyType = true
				break
			}
		}
		if !hasKeyType {
			return nil, stakingtypes.ErrValidatorPubKeyTypeNotSupported.Wrapf("got: %s, expected: %s", pk.Type(), cp.Validator.PubKeyTypes)
		}
	}
	return pk, nil
}
