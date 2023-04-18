package crosschain

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v4/contract"
	fxtypes "github.com/functionx/fx-core/v4/types"
)

const (
	FIP20CrossChainGas      = 40000 // 80000 - 160000
	CrossChainGas           = 40000 // 70000 - 155000
	CancelSendToExternalGas = 30000 // 70000 - 126000
	IncreaseBridgeFeeGas    = 40000 // 70000 - 140000
	BridgeCoinFeeGas        = 10000

	FIP20CrossChainMethodName      = "fip20CrossChain"
	CrossChainMethodName           = "crossChain"
	CancelSendToExternalMethodName = "cancelSendToExternal"
	IncreaseBridgeFeeMethodName    = "increaseBridgeFee"
	BridgeCoinMethodName           = "bridgeCoin"

	CrossChainEventName           = "CrossChain"
	CancelSendToExternalEventName = "CancelSendToExternal"
	IncreaseBridgeFeeEventName    = "IncreaseBridgeFee"
)

const (
	// EventTypeRelayTransferCrossChain
	// Deprecated
	EventTypeRelayTransferCrossChain = "relay_transfer_cross_chain"
	// EventTypeCrossChain new cross chain event type
	EventTypeCrossChain = "cross_chain"

	AttributeKeyDenom        = "coin"
	AttributeKeyTokenAddress = "token_address"
	AttributeKeyFrom         = "from"
	AttributeKeyRecipient    = "recipient"
	AttributeKeyTarget       = "target"
	AttributeKeyMemo         = "memo"
)

var (
	crossChainAddress = common.HexToAddress(fxtypes.CrossChainAddress)
	crossChainABI     = fxtypes.MustABIJson(contract.ICrossChainMetaData.ABI)
)

func GetAddress() common.Address {
	return crossChainAddress
}

func GetABI() abi.ABI {
	return crossChainABI
}
