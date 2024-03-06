package keeper

import (
	"math/big"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

func (k Keeper) bridgeCallERC20Handler(
	ctx sdk.Context,
	asset []byte,
	sender common.Address,
	to *common.Address,
	receiver sdk.AccAddress,
	dstChainID, message string,
	value sdkmath.Int,
	gasLimit, eventNonce uint64,
) error {
	tokens, amounts, err := types.UnpackERC20Asset(asset)
	if err != nil {
		return errorsmod.Wrap(types.ErrInvalid, "asset erc20")
	}
	coins, err := k.bridgeCallTransferToSender(ctx, sender.Bytes(), tokens, amounts)
	if err != nil {
		return err
	}

	switch dstChainID {
	case types.FxcoreChainID:
		if err = k.bridgeCallTransferToReceiver(ctx, sender.Bytes(), receiver, coins); err != nil {
			return err
		}
		if len(message) > 0 || to != nil {
			_, err := k.bridgeCallEvmHandler(ctx, sender, to, message, value, gasLimit, eventNonce)
			if err != nil {
				return err
			}
		}
	default:
		// not support chain, refund
	}
	// todo refund asset

	return nil
}

func (k Keeper) bridgeCallTransferToSender(ctx sdk.Context, receiver sdk.AccAddress, tokens [][]byte, amounts []*big.Int) (sdk.Coins, error) {
	tokens, amounts = types.MergeDuplicationERC20(tokens, amounts)

	mintCoins := sdk.NewCoins()
	unlockCoins := sdk.NewCoins()
	for i := 0; i < len(tokens); i++ {
		bridgeToken := k.GetBridgeTokenDenom(ctx, fxtypes.AddressToStr(tokens[i], k.moduleName))
		if bridgeToken == nil {
			return nil, errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
		}
		amount := sdkmath.NewIntFromBigInt(amounts[i])
		if !amount.IsPositive() {
			continue
		}
		coin := sdk.NewCoin(bridgeToken.Denom, amount)
		isOriginOrConverted := k.erc20Keeper.IsOriginOrConvertedDenom(ctx, bridgeToken.Denom)
		if !isOriginOrConverted {
			mintCoins = mintCoins.Add(coin)
		}
		unlockCoins = unlockCoins.Add(coin)
	}
	if mintCoins.IsAllPositive() {
		if err := k.bankKeeper.MintCoins(ctx, k.moduleName, mintCoins); err != nil {
			return nil, errorsmod.Wrapf(err, "mint vouchers coins")
		}
	}
	if unlockCoins.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, receiver, unlockCoins); err != nil {
			return nil, errorsmod.Wrap(err, "transfer vouchers")
		}
	}

	targetCoins := sdk.NewCoins()
	for _, coin := range unlockCoins {
		targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(ctx, receiver, coin, fxtypes.ParseFxTarget(fxtypes.ERC20Target))
		if err != nil {
			return nil, errorsmod.Wrap(err, "convert to target coin")
		}
		targetCoins = targetCoins.Add(targetCoin)
	}
	return targetCoins, nil
}

func (k Keeper) bridgeCallTransferToReceiver(ctx sdk.Context, sender sdk.AccAddress, receiver []byte, coins sdk.Coins) error {
	for _, coin := range coins {
		if coin.Denom == fxtypes.DefaultDenom {
			if err := k.bankKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(coin)); err != nil {
				return err
			}
			continue
		}
		if _, err := k.erc20Keeper.ConvertCoin(sdk.WrapSDKContext(ctx), &erc20types.MsgConvertCoin{
			Coin:     coin,
			Receiver: common.BytesToAddress(receiver).String(),
			Sender:   sender.String(),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) bridgeCallEvmHandler(ctx sdk.Context, sender common.Address, to *common.Address, message string, value sdkmath.Int, gasLimit, eventNonce uint64) (*evmtypes.MsgEthereumTxResponse, error) {
	callErr, callResult := "", false
	defer func() {
		attrs := []sdk.Attribute{
			sdk.NewAttribute(types.AttributeKeyEventNonce, strconv.FormatUint(eventNonce, 10)),
			sdk.NewAttribute(types.AttributeKeyEvmCallResult, strconv.FormatBool(callResult)),
		}
		if len(callErr) > 0 {
			attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyEvmCallError, callErr))
		}
		ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(types.EventTypeBridgeCallEvm, attrs...)})
	}()

	txResp, err := k.evmKeeper.CallEVM(ctx, sender, to, value.BigInt(), gasLimit, types.MustDecodeMessage(message), true)
	if err != nil {
		callErr = err.Error()
		return nil, err
	}

	callResult = !txResp.Failed()
	callErr = txResp.VmError
	return txResp, nil
}