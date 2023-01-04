package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

// RegisterCoin deploys an erc20 contract and creates the token pair for the existing cosmos coin
func (k Keeper) RegisterCoin(ctx sdk.Context, coinMetadata banktypes.Metadata) (*types.TokenPair, error) {
	// check if the conversion is globally enabled
	if !k.GetEnableErc20(ctx) {
		return nil, sdkerrors.Wrap(types.ErrERC20Disabled, "registration is currently disabled by governance")
	}

	decimals := getErc20Decimals(coinMetadata)

	// check if the denomination already registered
	if k.IsDenomRegistered(ctx, coinMetadata.Base) {
		return nil, sdkerrors.Wrapf(types.ErrTokenPairAlreadyExists, "coin denomination already registered: %s", coinMetadata.Base)
	}

	//base not register as alias
	if k.IsAliasDenomRegistered(ctx, coinMetadata.Base) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidMetadata, "denom %s already registered", coinMetadata.Base)
	}

	if len(coinMetadata.DenomUnits) > 0 && len(coinMetadata.DenomUnits[0].Aliases) > 0 {
		for _, alias := range coinMetadata.DenomUnits[0].Aliases {
			if alias == coinMetadata.Base || alias == coinMetadata.Display || alias == coinMetadata.Symbol {
				return nil, sdkerrors.Wrap(types.ErrInvalidMetadata, "alias can not equal base, display or symbol")
			}
			// alias not register as base
			if k.IsDenomRegistered(ctx, alias) {
				return nil, sdkerrors.Wrapf(types.ErrInvalidMetadata, "denom %s already registered", alias)
			}
			// alias must not register
			if k.IsAliasDenomRegistered(ctx, alias) {
				return nil, sdkerrors.Wrapf(types.ErrInvalidMetadata, "alias %s already registered", alias)
			}
		}
		k.SetAliasesDenom(ctx, coinMetadata.Base, coinMetadata.DenomUnits[0].Aliases...)
	}

	meta, isExist := k.bankKeeper.GetDenomMetaData(ctx, coinMetadata.Base)
	if isExist {
		if err := types.EqualMetadata(meta, coinMetadata); err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalidMetadata, err.Error())
		}
	} else {
		k.bankKeeper.SetDenomMetaData(ctx, coinMetadata)
	}

	addr, err := k.DeployUpgradableToken(ctx, k.moduleAddress, coinMetadata.Name, coinMetadata.Symbol, decimals)
	if err != nil {
		return nil, err
	}

	pair := types.NewTokenPair(addr, coinMetadata.Base, true, types.OWNER_MODULE)
	k.AddTokenPair(ctx, pair)
	return &pair, nil
}

// RegisterERC20 creates a cosmos coin and registers the token pair between the coin and the ERC20
func (k Keeper) RegisterERC20(ctx sdk.Context, contract common.Address) (*types.TokenPair, error) {
	if !k.GetEnableErc20(ctx) {
		return nil, sdkerrors.Wrap(types.ErrERC20Disabled, "registration is currently disabled by governance")
	}

	if k.IsERC20Registered(ctx, contract) {
		return nil, sdkerrors.Wrapf(types.ErrTokenPairAlreadyExists, "token ERC20 contract already registered: %s", contract.String())
	}

	erc20Data, err := k.QueryERC20(ctx, contract)
	if err != nil {
		return nil, err
	}

	// base denomination
	base := strings.ToLower(erc20Data.Symbol)
	if k.IsDenomRegistered(ctx, base) || erc20Data.Symbol == fxtypes.DefaultDenom {
		return nil, sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin denomination already registered: %s", erc20Data.Name)
	}

	// base not register as alias
	if k.IsAliasDenomRegistered(ctx, base) {
		return nil, sdkerrors.Wrapf(types.ErrInternalTokenPair, "alias %s already registered", base)
	}

	_, isExist := k.bankKeeper.GetDenomMetaData(ctx, base)
	if isExist {
		// metadata already exists; exit
		return nil, sdkerrors.Wrap(types.ErrInternalTokenPair, "denom metadata already registered")
	}

	// create a bank denom metadata based on the ERC20 token ABI details
	// metadata name is should always be the contract since it's the key
	// to the bank store
	metadata := banktypes.Metadata{
		Description: types.CreateDenomDescription(contract.String()),
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    base,
				Exponent: 0,
			},
		},
		Base:    base,
		Display: base,
		Name:    erc20Data.Name,
		Symbol:  erc20Data.Symbol,
	}

	// only append metadata if decimals > 0, otherwise validation fails
	if erc20Data.Decimals > 0 {
		metadata.DenomUnits = append(
			metadata.DenomUnits,
			&banktypes.DenomUnit{
				Denom:    erc20Data.Symbol,
				Exponent: uint32(erc20Data.Decimals),
			},
		)
	}

	if err := metadata.Validate(); err != nil {
		return nil, sdkerrors.Wrapf(err, "FIP20 token data is invalid for contract %s", contract.String())
	}
	k.bankKeeper.SetDenomMetaData(ctx, metadata)

	pair := types.NewTokenPair(contract, metadata.Base, true, types.OWNER_EXTERNAL)
	k.AddTokenPair(ctx, pair)
	return &pair, nil
}

// ToggleRelay toggles relaying for a given token pair
func (k Keeper) ToggleRelay(ctx sdk.Context, token string) (types.TokenPair, error) {
	id := k.GetTokenPairID(ctx, token)
	if len(id) == 0 {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token '%s' not registered by id", token)
	}

	pair, found := k.GetTokenPair(ctx, id)
	if !found {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token '%s' not registered", token)
	}

	pair.Enabled = !pair.Enabled

	k.SetTokenPair(ctx, pair)
	return pair, nil
}

// UpdateDenomAlias update denom alias
// if alias not registered, add to denom alias
// if alias registered with denom, remove from denom alias
// if alias registered, but not with denom, return error
func (k Keeper) UpdateDenomAlias(ctx sdk.Context, denom, alias string) (bool, error) {
	// check if the denom denomination already registered
	if !k.IsDenomRegistered(ctx, denom) {
		return false, sdkerrors.Wrapf(types.ErrInvalidDenom, "coin denomination not registered: %s", denom)
	}
	// check if the alias not registered
	if k.IsDenomRegistered(ctx, alias) {
		return false, sdkerrors.Wrapf(types.ErrInvalidDenom, "coin denomination already registered: %s", alias)
	}

	md, found := k.HasDenomAlias(ctx, denom)
	if !found {
		return false, sdkerrors.Wrapf(types.ErrInvalidMetadata, "denom %s not support update denom aliases", denom)
	}

	oldAliases := md.DenomUnits[0].Aliases
	newAliases := make([]string, 0, len(oldAliases)+1)

	registeredDenom, found := k.GetAliasDenom(ctx, alias)
	//check if the alias not register denom-alias
	if !found {
		newAliases = append(oldAliases, alias)
		k.SetAliasesDenom(ctx, denom, alias)
	} else if registeredDenom == denom {
		// check if the denom equal alias registered denom
		for _, denomAlias := range oldAliases {
			if denomAlias == alias {
				continue
			}
			newAliases = append(newAliases, denomAlias)
		}
		if len(newAliases) == 0 {
			return false, sdkerrors.Wrapf(types.ErrInvalidDenom, "can not remove, alias %s is the last one", alias)
		}
		k.DeleteAliasesDenom(ctx, alias)
	} else {
		//check if denom not equal alias registered denom, return error
		return false, sdkerrors.Wrapf(types.ErrInvalidDenom,
			"alias %s already registered, but denom expected: %s, actual: %s",
			alias, registeredDenom, denom)
	}

	md.DenomUnits[0].Aliases = newAliases
	k.bankKeeper.SetDenomMetaData(ctx, md)

	addFlag := len(newAliases) > len(oldAliases)
	return addFlag, nil
}

func (k Keeper) DeployUpgradableToken(ctx sdk.Context, from common.Address, name, symbol string, decimals uint8) (common.Address, error) {
	var tokenContract fxtypes.Contract
	if symbol == fxtypes.DefaultDenom {
		tokenContract = fxtypes.GetWFX()
		name = fmt.Sprintf("Wrapped %s", name)
		symbol = fmt.Sprintf("W%s", symbol)
	} else {
		tokenContract = fxtypes.GetERC20()
	}
	k.Logger(ctx).Info("deploy token", "name", name, "symbol", symbol, "decimals", decimals)

	//deploy proxy
	erc1967Proxy := fxtypes.GetERC1967Proxy()
	contract, err := k.DeployContract(ctx, from, erc1967Proxy.ABI, erc1967Proxy.Bin, tokenContract.Address, []byte{})
	if err != nil {
		return common.Address{}, err
	}

	_, err = k.CallEVM(ctx, tokenContract.ABI, from, contract, true, "initialize", name, symbol, decimals, k.moduleAddress)
	if err != nil {
		return common.Address{}, err
	}
	return contract, nil
}

func (k Keeper) DeployContract(ctx sdk.Context, from common.Address, abi abi.ABI, bin []byte, constructorData ...interface{}) (common.Address, error) {
	args, err := abi.Pack("", constructorData...)
	if err != nil {
		return common.Address{}, sdkerrors.Wrap(err, "pack constructor data")
	}
	data := make([]byte, len(bin)+len(args))
	copy(data[:len(bin)], bin)
	copy(data[len(bin):], args)

	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	_, err = k.evmKeeper.CallEVMWithData(ctx, from, nil, data, true)
	if err != nil {
		return common.Address{}, err
	}
	contractAddr := crypto.CreateAddress(from, nonce)
	return contractAddr, nil
}

func getErc20Decimals(md banktypes.Metadata) (decimals uint8) {
	decimals = uint8(0)
	for _, du := range md.DenomUnits {
		if du.Denom == md.Symbol {
			decimals = uint8(du.Exponent)
			break
		}
	}
	if md.Base == fxtypes.DefaultDenom {
		decimals = fxtypes.DenomUnit
	}
	return decimals
}
