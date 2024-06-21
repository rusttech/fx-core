package precompile

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	evmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) BridgeCall(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("bridge call method not readonly")
	}
	if c.router == nil {
		return nil, errors.New("bridge call router is empty")
	}

	var args crosschaintypes.BridgeCallArgs
	if err := evmtypes.ParseMethodArgs(crosschaintypes.BridgeCallMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}
	route, has := c.router.GetRoute(args.DstChain)
	if !has {
		return nil, errors.New("invalid dstChain")
	}
	sender := contract.Caller()

	coins := make([]sdk.Coin, 0, len(args.Tokens)+1)
	value := contract.Value()
	if value.Cmp(big.NewInt(0)) == 1 {
		totalCoin, err := c.handlerOriginToken(ctx, evm, sender, value)
		if err != nil {
			return nil, err
		}
		coins = append(coins, totalCoin)
	}
	for i, token := range args.Tokens {
		coin, err := c.handlerERC20Token(ctx, evm, sender, token, args.Amounts[i])
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	nonce, err := route.PrecompileBridgeCall(
		ctx,
		sender,
		args.Refund,
		coins,
		args.To,
		args.Data,
		args.Memo,
	)
	if err != nil {
		return nil, err
	}

	nonceNonce := big.NewInt(0).SetUint64(nonce)
	if err = c.AddLog(evm, crosschaintypes.BridgeCallEvent,
		[]common.Hash{sender.Hash(), args.Refund.Hash(), args.To.Hash()},
		evm.Origin,
		args.Value,
		nonceNonce,
		args.DstChain,
		args.Tokens,
		args.Amounts,
		args.Data,
		args.Memo,
	); err != nil {
		return nil, err
	}
	return crosschaintypes.BridgeCallMethod.Outputs.Pack(nonceNonce)
}